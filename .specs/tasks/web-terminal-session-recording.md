# Task Spec: Web 终端与会话录制

## 实现状态

- [x] L2 后端：Session 模型、终端 RPC schema、Agent PTY 执行器、asciinema 录制器、WebSocket 连接管理与数据桥接、路由注册
- [x] L3 前端：xterm.js 终端组件、终端页面与路由、Agent 列表「打开终端」入口
- [x] L4 测试：recorder / bridge / PTY 执行器单元测试
- [ ] 手动 E2E：需启动 amprobe + collia 进行端到端验证（待用户在运行环境执行）

## 目标

在 Amprobe 控制通道上实现基于浏览器的 Web 终端，并将会话内容按 asciinema v2 格式录制到 Server 本地。用户可在 Web UI 中选择目标 Agent，打开一个交互式 shell；Server 负责把浏览器输入转发给 Agent，并把 Agent 输出回写给浏览器，同时持久化会话记录以便回放。

## 对外可观测行为

1. **终端入口**
   - Web UI 提供入口，用户选择已在线的 Agent 后打开终端页面。
   - 终端页面加载成功后，用户应能立即看到远程 shell 提示符。

2. **输入/输出双向流**
   - 用户在键盘上的输入通过 WebSocket 实时发送到 Server，Server 立即通过反向 gRPC tunnel 转发到目标 Agent。
   - Agent 上 PTY 的输出实时通过反向 gRPC tunnel → Server → WebSocket 回写到浏览器终端。
   - 用户输入特殊按键（如方向键、Tab、Ctrl+C、Ctrl+D）必须按原始字节流透传，不得被浏览器或 Server 解释或截断。

3. **窗口大小调整**
   - 浏览器终端窗口大小变化时，应通知 Server，Server 再通过反向 tunnel 通知 Agent 调整 PTY 的 `rows` 和 `cols`。
   - 调整成功后，终端输出应自动换行或重排以适应新尺寸。

4. **会话录制**
   - 从 WebSocket 建立到关闭的整个会话，Server 必须按 asciinema v2 `.cast` 格式写入文件。
   - 录制文件头必须包含 `version: 2`、`width`、`height`、`timestamp`。
   - 录制数据行格式为 `[timestamp_seconds, "o", "base64_or_plain_output"]`，其中 `timestamp_seconds` 为相对会话开始的时间（秒，浮点数）。
   - 会话关闭时，录制文件必须正确 flush 并关闭。

5. **录制文件兼容性（可观测结果）**
   - 录制文件可通过标准 asciinema 播放器（如 `asciinema play`）正常播放。
   - 本 Task 不实现 Server 侧回放 API 或前端回放组件；回放能力由外部 asciinema 工具保证。

## 输入 / 输出约束

### WebSocket 连接

- **URL**: `GET /ws/terminal`
- **必需 Query/Header**:
  - `agent_id` (string, required): 目标 Agent 标识，从 query 或 `X-Agent-ID` header 获取。
  - WebSocket 升级前必须经过身份认证和授权中间件（复用现有 JWT/Casbin）。
- **WebSocket 消息类型**: 仅使用 `TextMessage` 和 `BinaryMessage`；控制帧（如 Close、Ping/Pong）由框架处理。

### WebSocket 消息格式

消息以 JSON 对象传输，字段如下：

| 字段 | 类型 | 含义 |
|------|------|------|
| `type` | string | `input` / `resize` / `output` / `error` |
| `data` | string | 类型为 `input` 或 `output` 时，为 base64 编码的字节数据 |
| `rows` | int | 类型为 `resize` 时，PTY 行数 |
| `cols` | int | 类型为 `resize` 时，PTY 列数 |
| `msg` | string | 类型为 `error` 时，可读错误信息 |

- `input`: 浏览器 → Server → Agent，数据为 base64 编码的用户键盘输入。
- `output`: Agent → Server → 浏览器，数据为 base64 编码的 PTY 输出。
- `resize`: 浏览器 → Server → Agent，携带 `rows` 和 `cols`。
- `error`: Server → 浏览器，发生不可恢复错误时发送，随后关闭连接。

### gRPC Tunnel RPC

新增两个 RPC 方法，通过 `common/rpc/tunnel` 的 `Frame` 双向流实现：

1. `TerminalSession`
   - Args: `TerminalSessionArgs`
     - `shell` (string, default: `/bin/bash`): 要执行的 shell。
     - `rows` (int, required): 初始行数。
     - `cols` (int, required): 初始列数。
   - Reply: 无单次 reply，通过 stream 持续返回输出帧，直到 `Eos=true`。

2. `ResizeTerminal`
   - Args: `ResizeTerminalArgs`
     - `rows` (int, required)
     - `cols` (int, required)
   - Reply: `ResizeTerminalReply`（空对象或错误）。

### 录制文件输出

- **路径**: 由 Server 配置指定，如 `data/sessions/`。
- **文件名**: `<session_id>.cast`，`session_id` 由 Server 生成（UUID 或 snowflake）。
- **内容**: 符合 asciinema v2 规范。

## 状态模型与迁移规则

### 会话状态（Server 侧）

```
[Idle] -- WebSocket 建立, agent_id 有效 --> [Connecting]
[Connecting] -- Agent 接受 TerminalSession 请求 --> [Active]
[Connecting] -- Agent 离线或拒绝 --> [Failed] -- 发送 error 并关闭 WebSocket --> [Closed]
[Active] -- 浏览器关闭 / Agent 断开 / 异常错误 --> [Closing]
[Closing] -- 录制文件 flush 完成, goroutine 退出 --> [Closed]
```

**禁止路径**:
- `[Active]` 不能直接回到 `[Connecting]`；异常后必须进入 `[Closing]` 并关闭。
- `[Closed]` 之后不得再写入录制文件或发送 WebSocket 消息。

### Agent 侧 PTY 状态

```
[NotStarted] -- 收到 TerminalSession 请求 --> [Running]
[Running] -- PTY 进程退出 或 收到 EOF --> [Exited]
[Running] -- 收到 ResizeTerminal 请求 --> [Running]
```

## 错误条件与失败语义

1. **缺少 `agent_id`**
   - Server 立即返回 WebSocket close frame，状态码 1008，消息 `missing agent_id`。

2. **目标 Agent 离线**
   - Server 通过 WebSocket 发送 `{"type":"error","msg":"agent offline"}`，然后关闭连接。

3. **Agent 创建 PTY 失败**
   - Agent 返回 RPC error，Server 转发为 WebSocket error 消息后关闭连接。

4. **WebSocket 写入失败**
   - Server 进入关闭流程，释放 tunnel stream 和 PTY。

5. **gRPC tunnel 流异常断开**
   - Server 检测到流断开后，关闭 WebSocket（code 1011），停止录制。

6. **录制文件写入失败**
   - 记录错误日志，但**不**中断终端会话；会话可继续使用，只是不被录制。

7. **权限不足**
   - 复用现有 Casbin 中间件，返回 403（WebSocket 升级前）或关闭连接。

## 结果约束（非功能性）

1. **时延**
   - 单次按键输入到终端回显的中位时延应 < 100ms（局域网环境下）。

2. **并发**
   - Server 应支持至少 100 个并发终端会话（每个会话 2 个 goroutine）。

3. **资源释放**
   - 任何连接关闭路径都必须确保：
     - PTY 子进程被 kill（防止僵尸进程）。
     - gRPC stream 被关闭。
     - WebSocket 被关闭。
     - 录制文件被 flush 并关闭。

4. **安全性**
   - 终端输入/输出必须仅在认证用户与授权 Agent 之间转发。
   - 禁止 Server 执行本地 shell；所有 shell 执行必须发生在 Agent 侧。
   - 录制文件路径必须限制在配置目录内，防止路径遍历。

5. **兼容性**
   - 支持 Linux 和 macOS 的 Agent（`creack/pty` 已支持）。
   - 浏览器端支持现代 Chrome/Firefox/Safari。

6. **可观测性**
   - Server 记录会话开始/结束日志，包含 `session_id`、`agent_id`、`user`、`duration`。
   - 错误场景记录结构化日志。
