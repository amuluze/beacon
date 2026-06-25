# 控制通道重构分析

> 参考项目：gravitational/teleport、dushixiang/next-terminal、moul/sshportal
> 项目：Amprobe（Server-Agent 监控/探测平台）

---

## 一、当前控制通道全景

### 1.1 现有能力

控制通道（rpcx RPC）当前包含 30+ 个 RPC 方法，通过 amprobe/pkg/rpc 的 Client 多 Agent 连接池路由到目标 Agent。

### 1.2 通信拓扑

Server 主动连接 Agent（Unix socket 本地或 TCP 跨机），Agent 列表静态配置。

### 1.3 主要局限

| 问题 | 表现 | 影响 |
|------|------|------|
| Agent 发现静态 | 硬编码配置 | 扩缩容需重启 Server |
| 无心跳检测 | 无法知在线状态 | 超时才知离线 |
| 终端禁用 | TermHandler 是 stub | 无远程 Web 终端 |
| 日志非流式 | 一次性全量读取 | 大日志超时 |
| 无审计 | 无操作记录 | 不可追溯 |
| 无会话录制 | 无录制 | 不可回放 |
| Server 前置连接 | Server 主动连接 Agent | NAT 防火墙困难 |

## 二、参考项目的关键模式

### 2.1 三者的共性模式

```
1. 反向连接（Agent → Server 建立连接）
   - 解决 NAT/防火墙问题，无须 Server 主动连接 Agent

2. 动态注册（Join Token / Cert Auth）
   - 节点通过 Token 加入集群，自动签发短期证书

3. 会话录制与回放
   - 所有交互操作录制为可回放文件，审计追溯的基础

4. 网关代理（Gateway/Proxy）
   - Server 作为网关/代理，浏览器 → Server → Agent 的完整链路

5. Web 终端
   - 浏览器 WebSocket → Server → Agent，xterm.js + PTY

6. 资产与凭证分离
   - 目标资产动态管理，凭证加密存储

7. 多租户 + RBAC
   - 用户/角色/资源三层模型，粒度权限控制
```

## 三、重构规划

### 3.1 总体架构演进

```
当前（阶段 0）
Server ── rpcx（主动连接）──→ Agent

阶段一（反向隧道）
Server ←── gRPC Stream ──── Agent（反向连接，双向流）

阶段二（注册 + 心跳）
Server ←── 心跳 + 负载 ──── Agent
Agent 通过 Join Token 注册，状态动态管理

阶段三（终端 + 录制）
Browser ── WebSocket ──→ Server ── gRPC Stream ──→ Agent ──→ PTY Shell
                                                  └──→ 会话录制文件
```

### 3.2 阶段一：反向隧道

Agent 启动后主动连接 Server 控制端口，Server 复用此连接下发指令。

Agent 配置新增：
```yaml
control:
  server: "amprobe.example.com:17000"
  join_token: "xxxx"
  agent_id: "host-a"
```

涉及范围：
- collia/service/rpc.go: rpcx Server → gRPC Client
- amprobe/pkg/rpc/rpc.go: rpcx Client → gRPC Server
- amprobe/configs/config.toml: 移除 [[Rpc.Agents]]
- collia/config.yml: 新增 control 段

### 3.3 阶段二：Agent 动态注册 + 心跳

流程：
1. 管理员在 Web UI 生成 Join Token
2. 安装 Agent 时携带 Token
3. Agent 启动后反向连接 Server
4. Server 验证 Token → 注册 Agent → 返回配置
5. 心跳维持在线状态

新增模型：
```go
type Agent struct {
    AgentID  string    `gorm:"uniqueIndex"`
    Hostname string
    OS       string
    Arch     string
    Version  string
    LastSeen time.Time
    Status   string    // online / offline
}
```

新增端点：
```
GET /api/v1/agent/list       → Agent 列表
GET /api/v1/agent/detail     → Agent 详情
POST /api/v1/agent/token     → 创建 Join Token
DELETE /api/v1/agent/token   → 吊销 Token
```

### 3.4 阶段三：Web 终端 + 会话录制

Browser → WebSocket → Server → gRPC Stream → Agent → PTY → /bin/bash

关键组件：
- collia 侧: PTY 创建 + 双向转发 + 窗口大小调整
- Server 侧: 会话管理 + WebSocket 代理
- 录制: asciinema 格式写入 Session 文件

### 3.5 阶段四：流式日志 + Docker Exec

- ContainerLogs: 从一次性读取改为 gRPC Server Stream 流式推送
- Docker Exec: 复用终端链路，远端替换为 docker exec -it

## 四、模式对照与采纳策略

| 模式 | 来源 | 采纳时机 | 说明 |
|------|------|---------|------|
| 反向隧道 | Teleport / sshportal | 阶段一 | Agent 建立长连接，复用下发指令 |
| Join Token 注册 | Teleport | 阶段二 | 安装时带 Token，免静态配置 |
| Web 终端 | Next Terminal / Teleport | 阶段三 | xterm.js + PTY |
| 会话录制 | 三个项目均有 | 阶段三 | asciinema 格式 |
| 操作审计 | Teleport / Next Terminal | 阶段三 | 记录所有控制操作 |
| 资产分组 | Next Terminal | 后续 | 标签系统 |
| RBAC | 三个项目均有 | 已有 | Casbin 已实现 |

## 五、实施建议

### 优先级路线

```
P0 阶段一：反向隧道
   └── 替代 rpcx，解决网络穿透，阻断性前提

P0 阶段二：动态注册 + 心跳
   └── 替代静态配置，依赖阶段一

P1 阶段三：Web 终端 + 会话录制
   └── 解锁核心交互，依赖阶段一 + 二

P1 阶段四：流式日志 + Docker Exec
   └── 优化已有功能，依赖阶段一

P2 资产分组 + 标签系统
   └── 提升多机体验，依赖阶段二
```

### 技术选型

| 项目 | 建议 | 理由 |
|------|------|------|
| 隧道传输 | gRPC Bidirectional Stream | 多路复用，流式原生，强类型 |
| 序列化 | Protocol Buffers (proto3) | 替换 Gob，跨语言 |
| 会话录制 | asciinema v2 (.cast) | 标准格式 |
| Web 终端 | xterm.js | 成熟稳定 |
| Agent 注册 | Join Token + 短期证书 | Teleport 验证过的模式 |

## 六、与监控通道的关系

控制通道重构不影响监控数据通道（HTTP POST /report），两者继续保持解耦：

- 监控通道：Agent → Server，HTTP JSON Push，**保持不变**
- 控制通道：双向 gRPC Stream，反向连接，**重构主体**
