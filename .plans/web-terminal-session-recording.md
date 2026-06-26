# Plan: Web 终端与会话录制

## 目标

实现基于浏览器的 Web 终端，支持选择目标 Agent、双向 PTY 流、窗口大小调整，并在 Server 侧按 asciinema v2 格式录制会话。

## 技术上下文

- 反向 gRPC tunnel 已落地：`common/rpc/tunnel` 提供 `ServerTunnel.StreamCall` 和 Agent 端 `dispatch`。
- WebSocket 基础设施已存在：`amprobe/service/websocket.go` 中 `/ws` 当前为 `TermHandler` stub。
- 前端已有 WebSocket 封装 `amprobe/web/src/components/Websocket/index.ts` 和日志查看组件 `ViewLog.vue`。
- Agent RPC dispatcher 在 `collia/service/rpc/dispatcher.go`，支持流式响应。

## 选定层级

本功能为**全栈新功能**，但数据库模型轻量，因此合并数据模型到后端层，保留 **L2（后端）、L3（前端）、L4（验证）** 三层。

## 任务分解

### L2 后端

#### 2.1 会话元数据模型
- **文件**: `amprobe/service/model/session.go`
- **说明**: 定义 `Session` GORM 模型，记录会话元数据。
- **字段**:
  - `SessionID` (string, uniqueIndex)
  - `AgentID` (string, index)
  - `UserID` (string)
  - `StartedAt` (time.Time)
  - `EndedAt` (time.Time, nullable)
  - `FilePath` (string)
  - `Status` (string: active / closed / failed)
  - `Width`, `Height` (int)
- **依赖**: 无
- **验证**: `cd amprobe && go build ./service/model/...`

#### 2.2 注册模型与配置
- **文件**: `amprobe/service/model/model.go`, `amprobe/configs/config.toml`
- **说明**: 在 GORM 自动迁移中注册 `Session`；新增 `[Session]` 配置段（录制目录、是否启用录制）。
- **依赖**: 2.1
- **验证**: `cd amprobe && go build ./service/model/...`

#### 2.3 通用终端 RPC Schema
- **文件**: `common/rpc/schema/terminal.go`
- **说明**: 定义 `TerminalSessionArgs`、`ResizeTerminalArgs`、`ResizeTerminalReply`。
- **依赖**: 无
- **验证**: `cd common && go build ./rpc/schema/...`

#### 2.4 Agent PTY 执行器
- **文件**: `collia/service/rpc/terminal.go`
- **说明**:
  - 使用 `github.com/creack/pty` 启动 `/bin/bash`。
  - 实现 `TerminalSessionStream(ctx, args, streamSender)`，持续读取 PTY 输出并通过 `streamSender` 发送。
  - 实现 `ResizeTerminal(args)` 调用 `pty.Setsize`。
  - 启动 goroutine 等待 bash 进程退出，退出时发送 `Eos=true`。
  - **清理保证**: stream 断开或 context 取消时，必须调用 `process.Kill()` 并 `Wait()`，防止僵尸进程。
- **依赖**: 2.3
- **验证**: `cd collia && go build ./service/rpc/...`

#### 2.5 Agent RPC 分发注册
- **文件**: `collia/service/rpc/dispatcher.go`
- **说明**: 在 `Call` switch 中新增 `TerminalSession` 和 `ResizeTerminal` 分支。
- **依赖**: 2.3, 2.4
- **验证**: `cd collia && go build ./... && go test ./service/rpc/...`

#### 2.6 Server 录制文件写入器
- **文件**: `amprobe/service/terminal/recorder.go`
- **说明**:
  - 实现 `asciinema v2` 写入器：写入 header、`[ts, "o", data]` 行。
  - 提供 `WriteOutput(data []byte)`、`Resize(rows, cols)`、`Close()` 方法。
  - 内部加锁，支持并发写入。
  - 文件路径通过 `filepath.Join(SessionDir, sessionID+".cast")` 生成，并用 `filepath.Clean` 限制在配置目录内。
- **依赖**: 无
- **验证**: `cd amprobe && go test ./service/terminal/...`

#### 2.7 Server 终端连接管理
- **文件**: `amprobe/service/terminal/handler.go`
- **说明**:
  - 实现 `TerminalHandler` struct，依赖 `rpc.Caller` 和 DB。
  - `Handler(conn *websocket.Conn)` 完成：
    - 读取 `agent_id`，校验非空。
    - 生成 `session_id`，写入 `Session` 记录（status=active）。
    - 校验 Agent 在线，调用 `StreamCall(ctx, "TerminalSession", args)` 建立 tunnel stream。
    - 启动清理 goroutine：监听 WebSocket 关闭/tunnel 断开/context 取消，统一进入关闭流程。
    - 关闭时更新 `Session` 记录（status=closed/failed, EndedAt）。
- **依赖**: 2.1, 2.2, 2.3, 2.6
- **验证**: `cd amprobe && go build ./service/terminal/...`

#### 2.8 Server 终端数据桥接与录制
- **文件**: `amprobe/service/terminal/bridge.go`
- **说明**:
  - 实现 `bridge(ctx, conn, stream, recorder, sessionID)`：
    - WebSocket 读 goroutine：解析 JSON 消息，input 写入 tunnel；resize 调用 `Call(ctx, "ResizeTerminal", args)`。
    - tunnel 读 goroutine：读取输出帧，base64 编码后写入 WebSocket；同时写入 recorder。
    - 任一 goroutine 出错时取消 context，触发另一方退出。
  - 保证录制器 `Close()` 在所有 goroutine 退出后被调用。
- **依赖**: 2.3, 2.6, 2.7
- **验证**: `cd amprobe && go build ./service/terminal/... && go test ./service/terminal/...`

#### 2.9 Server 路由与依赖注入
- **文件**: `amprobe/service/router.go`, `amprobe/service/injector.go`, `amprobe/service/wire_gen.go`, `amprobe/service/websocket.go`
- **说明**:
  - 将 `TermHandler` 改造为依赖 `rpcClient` 和 DB 的 `TerminalHandler`。
  - 在 `/ws/terminal` 注册终端 handler（保留 `/ws` 作为兼容别名或移除 stub）。
  - 更新 wire 注入。
- **依赖**: 2.7, 2.8
- **验证**: `cd amprobe && go build ./... && go test ./...`

---

### L3 前端

#### 3.1 安装 xterm.js 依赖
- **文件**: `amprobe/web/package.json`
- **说明**: 添加 `xterm`、`xterm-addon-fit`、`xterm-addon-web-links` 依赖并安装。
- **依赖**: 无
- **验证**: `cd amprobe/web && npm install`

#### 3.2 xterm.js 终端组件
- **文件**: `amprobe/web/src/components/Terminal/index.vue`
- **说明**:
  - 引入 `xterm` 及相关 addon。
  - 组件接收 `agentId` prop，创建 `Terminal` 实例并连接到 `/ws/terminal?agent_id=xxx`。
  - 键盘输入编码为 base64 发送 `input` 消息。
  - 收到 `output` 消息解码后写入 terminal。
  - 收到 `error` 消息时显示错误并关闭连接。
  - 窗口 resize 时通过 `fit` addon 计算 cols/rows，发送 `resize` 消息。
- **依赖**: 3.1
- **验证**: `cd amprobe/web && npm run build`

#### 3.3 终端页面与路由
- **文件**: `amprobe/web/src/views/terminal/index.vue`, `amprobe/web/src/router/dynamic.ts`
- **说明**:
  - 新增终端页面，包含 Agent 选择器 + 终端组件。
  - 在动态路由中注册 `/terminal`。
- **依赖**: 3.2
- **验证**: `cd amprobe/web && npm run build`

#### 3.4 Agent 列表入口
- **文件**: `amprobe/web/src/views/agent/index.vue`（或现有 Agent 列表页）
- **说明**: 在 Agent 列表增加"打开终端"按钮，跳转到 `/terminal?agent_id=xxx`。
- **依赖**: 3.3
- **验证**: 手动验证页面跳转

---

### L4 集成验证

#### 4.1 单元测试
- **文件**: `amprobe/service/terminal/recorder_test.go`, `amprobe/service/terminal/bridge_test.go`, `collia/service/rpc/terminal_test.go`
- **说明**:
  - 录制器：验证 header、resize、输出数据行格式正确，文件可播放。
  - PTY 执行器：mock `streamSender`，验证启动和退出时发送 `Eos`。
  - 桥接：使用 fake WebSocket 和 fake tunnel stream 验证双向转发。
- **依赖**: 2.4, 2.6, 2.8
- **验证**:
  - `cd amprobe && go test ./service/terminal/...`
  - `cd collia && go test ./service/rpc/...`

#### 4.2 全量构建验证
- **说明**:
  - `task amprobe:build`
  - `task collia:build`
  - `task web:build`
- **依赖**: L2, L3 全部完成
- **验证**: 三个构建命令均成功

#### 4.3 端到端手动测试
- **步骤**:
  1. 启动 amprobe 和 collia。
  2. 登录 Web UI，进入 Agent 列表，点击"打开终端"。
  3. 执行 `ls`、`top`、`exit` 等命令，验证输出正确。
  4. 调整浏览器窗口大小，验证 `stty size` 输出变化。
  5. 检查 Server 录制目录下生成 `.cast` 文件，使用 `asciinema play` 播放验证。
  6. 关闭终端页面后，验证 `s_session` 表状态为 `closed`。
  7. 强制刷新页面或断开 Agent 网络，验证无僵尸进程和 goroutine 泄漏。
- **依赖**: L2, L3 全部完成
- **验证**: 手动 checklist 全部通过

## 依赖图

```
2.1 ──→ 2.2
2.3 ──→ 2.4 ──→ 2.5
2.3, 2.6 ──→ 2.8
2.1, 2.2, 2.3, 2.6 ──→ 2.7
2.7, 2.8 ──→ 2.9
3.1 ──→ 3.2 ──→ 3.3 ──→ 3.4
2.9, 3.4 ──→ 4.2
2.4, 2.6, 2.8 ──→ 4.1
L2, L3 ──→ 4.3
```

## 可并行项

- L2 后端任务与 L3 前端任务可高度并行。
- 2.1/2.2（模型/配置）与 2.3/2.6（schema/recorder）可并行。
- 3.1 安装依赖后可立即开始 3.2。

## 风险与回滚

| 风险 | 等级 | 缓解措施 |
|------|------|----------|
| PTY 子进程泄漏 | 高 | 2.4 中明确要求 stream 断开/context 取消时 kill + wait；4.3 中验证 |
| 大量并发会话导致 goroutine 暴涨 | 中 | 2.7 中每个会话限制 3 个 goroutine，必要时后续迭代添加并发限制 |
| 录制文件路径遍历 | 中 | 2.6 中使用 `filepath.Clean` 限制在配置目录 |
| xterm.js 打包体积增大 | 低 | 3.2 中按需引入 addon |
| macOS 与 Linux PTY 行为差异 | 低 | 使用 creack/pty 跨平台封装 |

## 提交策略建议

按任务原子提交，建议顺序：

1. `feat: 添加会话录制数据模型`
2. `feat: 定义终端 RPC schema`
3. `feat: 实现 Agent 侧 PTY 执行器`
4. `feat: 注册终端 RPC 分发`
5. `feat: 实现 Server 侧 asciinema 录制器`
6. `feat: 实现 WebSocket 终端连接管理`
7. `feat: 实现 WebSocket 与 tunnel 数据桥接`
8. `feat: 注册终端 WebSocket 路由`
9. `feat: 前端安装 xterm 依赖`
10. `feat: 前端 xterm 终端组件与页面`
11. `feat: Agent 列表添加终端入口`
12. `test: 添加终端与会话录制单元测试`
13. `docs: 更新终端功能配置说明`

## 验证清单

- [ ] `cd amprobe && go test ./...` 通过
- [ ] `cd collia && go test ./...` 通过
- [ ] `cd common && go test ./...` 通过
- [ ] `task amprobe:build` 成功
- [ ] `task collia:build` 成功
- [ ] `task web:build` 成功
- [ ] 手动 E2E 测试通过
