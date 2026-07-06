# P0 修复总结

## ✅ 编译状态
| 模块 | 构建 | 测试 |
|------|------|------|
| `beacon` | ✅ 通过 | 25/25 通过 |
| `collia` | ✅ 通过 | 5/5 通过 |
| `common` | ✅ 通过 | 4/4 通过 |

## 变更清单

### 模块路径统一
- `beacon/go.mod`: `module amprobe` → `module beacon`
- 16 个源文件：`"amprobe/` → `"beacon/`

### Collia 构建修复
- 新增 `collia/service/version.go`（`NewVersion`/`Version`）
- `NewRPCServer` 接受 `Version` 参数并传递 `rootDir` 给 `rpc.NewService`

### API 表面补齐
- `contextx`: `ResolveAgentID`, `ErrMissingAgentID`, `ErrInvalidAgentID`, `IsValidAgentID`, `maxAgentIDLen`
- `hash`: `SHA1String`
- `errors`: `New401Error`, `New409Error`, `Error()` 方法，`FromError` 增强（404/504/500）
- `fiberx`: `ServiceError` 类型断言辅助函数
- `config`: 新增 `CORS`/`RateLimit`/`App`/`Retention` 字段
- `agent`: `NewAgentService` 分离 tunnel，新增 `SetTunnel` 方法
- `rpc`: 修复 `NewTunnelClient` 调用（移除已弃用的 `defaultID` 参数）

### Tunnel 协议增强
- `AgentInfo` + `RegistrationPayload` 结构化注册
- `OnAgentConnect(agentID)` → `OnAgentConnect(info AgentInfo)`
- `ServerOption`, `WithJoinToken`, `WithAgentLifecycle`
- `FRAME_STREAM_END` 路由修复（stream 而非 pending call）
- `StreamCall` ctx 取消 goroutine 清理

### Schema & Service 增强
- `CPUInfoReply`, `MemoryInfoReply`, `DiskInfo`: 新增 `Timestamp`, `Stale` 字段
- `NewHostService`: 支持过期阈值参数
- `isStale()`: 时间戳过期判断

### 测试基础设施
- 任务测试 DB: 新增 `model.Agent` 迁移、Mail 种子数据
- `sendMail`: 失败时仅警告日志，不阻断事务
- Tunnel 测试: 超时模式避免挂起
