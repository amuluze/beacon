# Common 模块

## 模块职责

`common` 是 `beacon` Server 与 `collia` Agent **共用的 Go module 顶层包**（与 `common/docker`、`common/logger`、`common/timex` 是不同的独立 Go module，但代码路径都在 `common/` 目录下），承载以下跨模块契约：

- 数据库选项与查询选项：`database/options.go`、`database/queryoption.go`，提供 GORM 初始化、连接池配置、按列过滤的 `QueryOption` 链。
- RPC schema 与反向 tunnel 协议：`rpc/schema/` 定义 Server-Agent 通信的消息结构（含 `freshness.go` 用于新鲜度判断），`rpc/tunnel/` 定义反向 gRPC tunnel 协议（`ReverseTunnelServer` 接口、`ServerTunnel` 实现、`AgentLifecycle` 钩子），承载方式仅为 Unix Socket（单机部署，`listenControl` 拒绝非 `unix://` 目标）。
- 数据库驱动列表：间接依赖 `glebarez/sqlite`、`gorm.io/driver/{mysql,postgres,clickhouse}`，由上层按需选择。

注意：`common`（顶级 Go module）与 `common/docker`、`common/logger`、`common/timex` **是平级的不同 Go module**，后者分别在独立的 `github.com/amuluze/docker`、`github.com/amuluze/amutool/logger`、`github.com/amuluze/amutool/timex` 中维护。当前仓库通过 `go.work` 的 `use` 指令把它们放在同一工作区。

相关源码包：`common/database`、`common/rpc`、`common/rpc/schema`、`common/rpc/tunnel`。

## 职责边界

| 职责 | 说明 | 是否负责 |
|------|------|----------|
| 数据库选项与查询抽象 | GORM 初始化、连接池、QueryOption 链 | yes |
| RPC schema 与 tunnel 协议 | Server-Agent 消息结构、反向 gRPC tunnel | yes |
| 数据库驱动选型 | 提供多种驱动间接依赖，由上层选用 | yes |
| 业务逻辑 | 不实现任何业务规则 | no |
| UI / HTTP 入口 | 不暴露任何用户接口 | no |
| 跨模块协调 | 仅提供跨模块共享类型与配置 | conditional |
| 其他模块内部状态 | 不直接读写其他模块私有实现 | no |

## 模块事实

- 模块路径：`common`
- Go 版本：`1.25.0`
- Go 源码文件：16（顶级 common module 非测试 `.go`，不含 common/docker、common/logger、common/timex 子模块）
- Go 测试文件：9
- 包名：`database`、`rpc`、`schema`、`tunnel`
- 推断角色：跨 Server/Agent 复用的 Go module 顶级包，承载数据库选项、RPC schema 与 Unix Socket tunnel 协议

## 接口与符号信号

本节记录源码扫描得到的导出符号或前端入口信号；它们是候选接口线索，不等同于跨模块公开 API。跨模块调用仍以 facade、模块文档和 Domain Spec 约束为准。

- `interface AgentLifecycle (common/rpc/tunnel/server.go)`
- `interface IService (common/rpc/service.go)`
- `interface ReverseTunnelClient (common/rpc/tunnel/tunnel_grpc.pb.go)`
- `interface ReverseTunnelServer (common/rpc/tunnel/tunnel_grpc.pb.go)`
- `interface UnsafeReverseTunnelServer (common/rpc/tunnel/tunnel_grpc.pb.go)`
- `type Handler (common/rpc/tunnel/agent.go)`
- `func NewServerTunnel (common/rpc/tunnel/server.go)`
- `func RegisterReverseTunnelServer (common/rpc/tunnel/tunnel_grpc.pb.go)`
- `func WithAgentLifecycle (common/rpc/tunnel/server.go)`
- `func WithJoinToken (common/rpc/tunnel/server.go)`
- `struct AgentInfo (common/rpc/tunnel/server.go)`
- `struct AgentOfflineError (common/rpc/tunnel/server.go)`
- 另有其他导出符号未展开；请用 `rg`、codegraph 或 IDE 按需查看完整清单。

## 依赖关系

- `github.com/glebarez/sqlite` `v1.11.0`
- `google.golang.org/grpc` `v1.81.1`
- `google.golang.org/protobuf` `v1.36.11`
- `gorm.io/driver/clickhouse` `v0.6.1`
- `gorm.io/driver/mysql` `v1.5.7`
- `gorm.io/driver/postgres` `v1.5.9`
- `gorm.io/gorm` `v1.25.12`
- 间接依赖 38 个，完整列表以模块 `go.mod` 为准。

## 配置与入口信号

- 配置：`Config` in `common/logger/config.go`
- 配置：`IPAM` at `common/docker/network.go:49`
- 配置：`IPAM` at `common/docker/network.go:91`
- 配置：`IPAM` at `common/docker/network.go:133`
- 配置：`Option` in `common/database/options.go`
- 配置：`Option` in `common/logger/config.go`
- 配置：`QueryOption` in `common/database/queryoption.go`
- 配置：`ServerOption` in `common/rpc/tunnel/server.go`
- 配置：`SubNetworkConfig` in `common/docker/network.go`

## 数据模型与状态

- 数据模型以源码结构、导出类型信号和测试断言为准；本节只记录静态发现结果。
- 共享 RPC 契约不定义 DNS、系统时间、时区、主机重启或关机消息；Docker 镜像源配置使用独立 Docker schema。
- `struct AgentUnauthorizedError (common/rpc/tunnel/server.go)`
- `struct DuplicateAgentError (common/rpc/tunnel/server.go)`
- `struct InvalidAgentIDError (common/rpc/tunnel/server.go)`
- `struct RPCError (common/rpc/tunnel/server.go)`
- `struct RegistrationPayload (common/rpc/tunnel/server.go)`
- `struct ServerTunnel (common/rpc/tunnel/server.go)`
- `struct UnimplementedReverseTunnelServer (common/rpc/tunnel/tunnel_grpc.pb.go)`
- `type ReverseTunnel_TunnelServer (common/rpc/tunnel/tunnel_grpc.pb.go)`
- `type ServerOption (common/rpc/tunnel/server.go)`
- `struct Config (common/logger/config.go)`
- 另有其他导出模型/接口符号未展开；完整清单以源码为准。

- PostgreSQL 连接选项支持 `WithSSLMode`：空值沿用disable；DSN中的连接参数使用libpq转义，包含空格、引号或反斜线的值不会被解释为额外选项。

## 使用注意事项

- 不跨越模块边界直接依赖内部路径或未导出符号。
- 修改跨模块公开 API、导出符号信号、配置或副作用边界时，同步更新本文件、根入口文档和相关 Domain Spec。
- 配置、凭据和外部副作用只记录语义边界，不记录真实敏感值。

## 验证命令

```bash
# common 顶级 module 无独立 Taskfile.yml，验证通过 Go 标准工具链完成
cd common

# 编译检查
go build ./...

# 单元测试（带 race 检测）
go test -race ./...

# 静态分析
go vet ./...

```

## Taskfile Commands

- `task beacon`
- `task beacon-web`
- `task collia`
- `task default`
