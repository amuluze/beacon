# Beacon 模块

<!-- beaver:managed-project-doc -->

## 模块职责

`beacon` 的主要角色是：HTTP/API service module: Web/API 接入、路由注册、请求校验和服务协调。

- 相关源码包：`beacon/cmd/beacon`, `beacon/pkg/auth`, `beacon/pkg/auth/jwtauth`, `beacon/pkg/contextx`, `beacon/pkg/errors`, `beacon/pkg/fiberx`。

## 职责边界

| 职责 | 说明 | 是否负责 |
|------|------|----------|
| 模块内实现 | HTTP/API service module: Web/API 接入、路由注册、请求校验和服务协调 | yes |
| 跨模块协调 | 仅通过公开接口或项目约定完成 | conditional |
| 其他模块内部状态 | 不直接读写其他模块私有实现 | no |

## 模块事实

- 模块路径：`beacon`
- Go 版本：`1.25.0`
- Go 源码文件：194
- Go 测试文件：60
- 包名：`account`, `agent`, `alarm`, `api`, `audit`, `auth`, `authcookie`, `client`, `container`, `contextx`, `dingtalk`, `errors`, `fiberx`, `hash`, `health`, `host`, `jwtauth`, `mail`, `main`, `middleware`, `model`, `report`, `repository`, `rpc`, `schema`, `secret`, `service`, `task`, `terminal`, `testutil`, `title`, `utils`, `uuid`, `validatex`, `web`
- 推断角色：HTTP/API service module: Web/API 接入、路由注册、请求校验和服务协调

## 接口与符号信号

本节记录源码扫描得到的导出符号或前端入口信号；它们是候选接口线索，不等同于跨模块公开 API。跨模块调用仍以 facade、模块文档和 Domain Spec 约束为准。

- `interface Auther (beacon/pkg/auth/auth.go)`
- `interface Caller (beacon/pkg/rpc/rpc.go)`
- `interface Connection (beacon/service/terminal/bridge.go)`
- `interface DBPinger (beacon/service/health/api/health.go)`
- `interface IAccountRepository (beacon/service/account/repository/account.go)`
- `interface IAccountService (beacon/service/account/service/account.go)`
- `interface IAlarmRepository (beacon/service/alarm/repository/alarm.go)`
- `interface IAlarmService (beacon/service/alarm/service/alarm.go)`
- `interface IAuditRepo (beacon/service/audit/repository/audit.go)`
- `interface IAuditService (beacon/service/audit/service/audit.go)`
- `interface IAuthRepository (beacon/service/auth/repository/auth.go)`
- `interface IAuthService (beacon/service/auth/service/auth.go)`
- `interface IContainerRepo (beacon/service/container/repository/container.go)`
- `interface IContainerService (beacon/service/container/service/container.go)`
- `interface IDingTalkRepository (beacon/service/dingtalk/repository/dingtalk.go)`
- `interface IDingTalkService (beacon/service/dingtalk/service/dingtalk.go)`
- `interface IHostRepo (beacon/service/host/repository/host.go)`
- `interface IHostService (beacon/service/host/service/host.go)`
- 另有其他导出符号未展开；请用 `rg`、codegraph 或 IDE 按需查看完整清单。

## 依赖关系

- `github.com/amuluze/amutool/logger` `v0.0.0-20240821104128-caed9cc0d402`
- `github.com/amuluze/amutool/timex` `v0.0.0-20250508153823-fe9a5de55958`
- `github.com/casbin/casbin/v2` `v2.98.0`
- `github.com/casbin/gorm-adapter/v3` `v3.28.0`
- `github.com/go-playground/validator/v10` `v10.22.0`
- `github.com/gofiber/contrib/websocket` `v1.3.2`
- `github.com/gofiber/fiber/v2` `v2.52.5`
- `github.com/golang-jwt/jwt/v5` `v5.0.0`
- `github.com/google/uuid` `v1.6.0`
- `github.com/google/wire` `v0.6.0`
- `github.com/patrickmn/go-cache` `v2.1.0+incompatible`
- `github.com/pkg/errors` `v0.9.1`
- `github.com/spf13/viper` `v1.19.0`
- `github.com/stretchr/testify` `v1.9.0`
- `github.com/urfave/cli/v2` `v2.27.4`
- `golang.org/x/crypto` `v0.54.0`
- `gopkg.in/gomail.v2` `v2.0.0-20160411212932-81ebce5c23df`
- `gorm.io/gorm` `v1.25.12`
- 间接依赖 70 个，完整列表以模块 `go.mod` 为准。

## 配置与入口信号

- 配置：`BASE_URL` at `beacon/service/agent_install.go:255`
- 配置：`BEACON_AGENT_INSTALL_TOKEN` at `beacon/service/auth.go:148`
- 配置：`BEACON_AGENT_INSTALL_TOKEN` at `beacon/service/config.go:84`
- 配置：`BEACON_AUTH_SIGNING_KEY` at `beacon/service/auth.go:93`
- 配置：`BEACON_AUTH_SIGNING_KEY` at `beacon/service/config.go:81`
- 配置：`BEACON_BOOTSTRAP_ADMIN_PASSWORD` at `beacon/service/prepare.go`
- 配置：`BEACON_DATA_ENCRYPTION_KEY` at `beacon/service/data_encryption.go`
- 配置：`BEACON_CONTAINER_NAME` at `beacon/service/install_report.go:97`
- 配置：`BEACON_CONTROL_JOIN_TOKEN` at `beacon/service/auth.go:123`
- 配置：`BEACON_CONTROL_JOIN_TOKEN` at `beacon/service/config.go:87`
- 配置：`BEACON_CONTROL_PORT` at `beacon/service/install_report.go:96`
- 配置：`BEACON_HTTP_PORT` at `beacon/service/install_report.go:95`
- 配置：`BEACON_IMAGE` at `beacon/service/install_report.go:91`
- 配置：`BEACON_PUBLIC_BASE_URL` at `beacon/service/install_report.go:93`
- 配置：`BEACON_VERSION` at `beacon/service/install_report.go:92`
- 配置：`CORS` at `beacon/service/app.go:32`
- 配置：`CORS` at `beacon/service/app.go:75`
- 配置：`Config` in `beacon/service/config.go`
- 配置：`Option` in `beacon/pkg/auth/jwtauth/options.go`
- 入口：`api := app.Group("/api")` at `beacon/service/router.go:70`
- 入口：`app.Get("/health", a.healthProbe.Liveness)` at `beacon/service/router.go:76`
- 入口：`app.Get("/ready", a.healthProbe.Readiness)` at `beacon/service/router.go:77`
- 入口：`app.Get("/ws/terminal", webSocketAuth, websocket.New(a.terminalHandler.Handle))` at `beacon/service/router.go:154`
- 入口：`app.Get("/ws/:id", webSocketAuth, websocket.New(a.loggerHandler.Handler))` at `beacon/service/router.go:155`
- 入口：`app.Use("/", filesystem.New(filesystem.Config{` at `beacon/service/app.go:46`
- 入口：`app.Use("/api/v1/auth/login", limiter.New(limiter.Config{` at `beacon/service/router.go:87`
- 入口：`app.Use("/api/v1/host/report", limiter.New(limiter.Config{` at `beacon/service/router.go:97`
- 入口：`app.Use("/ws", func(c *fiber.Ctx) error {` at `beacon/service/router.go:144`
- 入口：`app.Use(func(c *fiber.Ctx) error {` at `beacon/service/router.go:115`
- 入口：`app.Use(middleware.CasbinMiddleware(` at `beacon/service/router.go:136`
- 入口：`app.Use(middleware.GlobalRateLimitMiddleware(config.RateLimit.GlobalMax))` at `beacon/service/app.go:35`
- 入口：`app.Use(middleware.PanicMiddleware())` at `beacon/service/app.go:43`
- 入口：`app.Use(middleware.StackMiddleware)` at `beacon/service/app.go:44`
- 入口：`app.Use(middleware.UserAuthMiddleware(` at `beacon/service/router.go:129`
- 入口：`g.Get("/alarm_query", a.alarmAPI.AlarmQuery).Name("查询告警阈值")` at `beacon/service/router.go:294`
- 入口：`g.Get("/message_setting", a.alarmAPI.AlarmMessageSettingQuery).Name("查询告警消息标题")` at `beacon/service/router.go:295`
- 入口：`g.Post("/message_setting_update", a.alarmAPI.AlarmMessageSettingUpdate).Name("更新告警消息标题")` at `beacon/service/router.go:296`
- 入口：`g.Get("/query", a.dingTalkAPI.Query).Name("查询钉钉告警配置")` at `beacon/service/router.go:286`
- 入口：`g.Post("/update", a.dingTalkAPI.Update).Name("更新钉钉告警配置")` at `beacon/service/router.go:287`
- 入口：`g.Post("/test", a.dingTalkAPI.Test).Name("测试钉钉告警")` at `beacon/service/router.go:288`

## 数据模型与状态

- 数据模型以源码结构、导出类型信号和测试断言为准；本节只记录静态发现结果。
- `MonitorCPU.CPUPercent`、`MonitorMemory.MemPercent` 和 `MonitorDisk.DiskPercent` 均沿用 Collia/gopsutil 的 `0–100` 百分数单位，不做二次乘以 100 的换算。
- CPU、内存和磁盘告警使用 GORM `CreatedAt` 作为 Server 接收时间。正持续时间要求样本覆盖完整窗口、最新样本足够新、相邻样本间隔不超过 90 秒且所有样本严格超过阈值；窗口不完整、陈旧或间断时保持未知态。磁盘按 Agent + 设备独立应用同一窗口。
- 资源告警按 Agent（磁盘额外按设备）构造状态与去重键；首次触发立即通知，持续异常每 10 分钟最多重复一次，获得正常证据后写一次恢复审计/通知并清除去重状态。告警检查、审计写入和缓存更新由互斥区保护，避免重叠定时任务重复发送。
- 容器状态任务只读取最近 2 分钟内写入的每容器最新记录，陈旧历史状态不触发变化告警。
- 钉钉告警配置以 `s_dingtalk` 单例记录持久化；查询 API 只返回 Webhook 掩码和凭据配置状态，更新时空白凭据保留原值，签名密钥可显式清除。
- 钉钉 HTTP sender 只接受 `https://oapi.dingtalk.com/robot/send` 与唯一非空 `access_token`，使用 5 秒超时并禁止重定向；可选签名通过 HMAC-SHA256 生成。
- 告警任务先持久化审计，再独立尝试邮件和钉钉；任一通道失败只写告警日志，不回滚审计或阻断另一通道。
- 告警消息标题以 `s_alarm_message_setting` 单例持久化，缺失时初始化为 `Beacon 主机告警` 且重复启动不覆盖用户值；实际邮件 Subject、邮件测试 Subject、实际钉钉首行和钉钉测试首行共享该标题。标题按单次通知只读取一次，缺失或非法存量值安全回退默认标题，审计正文不拼接标题。
- SMTP 密码、钉钉 Webhook 和签名密钥使用 `BEACON_DATA_ENCRYPTION_KEY` 对应的 AES-256-GCM 密钥加密存储；启动时事务迁移历史明文，读取时兼容并回写遗留值。该密钥缺失或非法时拒绝启动。
- 浏览器 JWT 仅写入 HttpOnly、SameSite=Strict Cookie，CSRF Cookie 可读并要求 `X-CSRF-Token` 双提交；登录响应只返回 `must_change_password`。新库由 `BEACON_BOOTSTRAP_ADMIN_PASSWORD` 创建唯一管理员并强制改密，历史公开默认密码也被迁移到强制改密状态。
- 生产环境拒绝非 HTTPS Web 请求并要求 Agent 控制通道双向 TLS；WebSocket 仅使用同源 access Cookie，不接受 Authorization 或 query token。
- 控制通道 CA 与服务端证书运行于 `/etc/beacon/control-certs`；官方发布 Compose 通过 `BEACON_CONTROL_CERT_DIR` 持久化该目录，避免容器升级静默轮换 CA。
- Beacon 发布镜像在隔离 Node 阶段生成干净 Web 产物，再以静态、路径裁剪和符号裁剪参数构建 Server；Alpine 运行时仍通过 Nginx + Supervisor 暴露 80 端口，并同时携带 amd64/arm64 Collia 安装包。
- `/api/v1/host/report` 只接受 `X-Install-Token`。Server Token 未配置返回 503，缺失/错误请求头返回 401，query token 不参与认证。
- Host HTTP API 不注册 DNS、系统时间、时区、主机重启或关机设置端点；主机监控查询、文件操作与 Docker 控制能力保持独立。
- Agent `online/offline` 表示当前 Beacon 进程的反向控制隧道可达性：绑定新的 tunnel 管理器时先把上一进程遗留状态归零，列表响应再按实时连接注册表投影；`last_seen` 只由真实连接/心跳更新，HTTPS 监控上报新鲜度不伪装成控制通道在线。
- Collia 注册携带 version/os/arch；Agent 持久化完成后异步进入 `Updater`。Updater 只把空版本或低版本更新到镜像注入的目标版本，按 `amd64/arm64` 选择与安装端点相同的文件并计算 SHA256；相同版本静默跳过，高版本、未知架构和显式非 Linux 目标拒绝更新。
- Collia 自动更新使用单 Agent 执行集合、全局并发槽和失败冷却；管理员专属 `POST /api/v1/agent/upgrade` 可绕过冷却手动重试。开始审计写入失败时不执行 RPC，开始/成功/失败记录均使用 `system` 用户并携带 Agent ID；Agent 状态不复用为更新进度。
- 用户删除 API 要求 1–32 个格式有效的明确 UUID；Service 会在持久化前确认全部目标存在，删除成功后同步移除对应 Casbin grouping policy，不存在目标返回 not found。
- Server 启动在事务中维护 `cpu`、`memory`、`disk` 各一条活动阈值：旧 `mem` 按最早记录迁移并保留用户值，同类型重复记录软删除，缺失类型才创建默认值。阈值更新按 `id + canonical type` 定位且不允许修改类型，零值会真实写入，目标不存在返回 not found。
- 容器创建的可选 Docker 网络 ID 接受 1–64 字符；该字段按 Docker Engine 可解析标识处理，不假设 Collia 上报值为固定长度 UUID。
- 容器指标历史按批次追加以支持趋势查询；容器管理列表、总数和镜像占用数只读取目标 Agent 的最新容器批次，包含未发布端口的容器，并排除已从最新批次消失的历史容器。
- `/ws/terminal` 在动态容器日志路由 `/ws/:id` 之前精确注册，并在 WebSocket upgrade 前读取 HttpOnly access Cookie；真实 `terminal.Handler` 等待 Agent PTY 可接受 resize 后发送 ready，以 Ping 控制帧保持空闲代理链路活跃，关闭时携带 Agent ID 调用 `TerminalClose`。
- WebSocket 终端与容器日志链路均有读写 deadline 兜底，半开 TCP（连接还在但客户端不读不写）不会让 Server 端 goroutine 永久阻塞：终端读 deadline 10 分钟、单次写 deadline 10 秒，活跃会话每次读到消息即重置读 deadline；容器日志每次写前设置 10 秒写 deadline。
- 流式控制调用（`TunnelClient.StreamCall`）派生最长 2 小时 `maxStreamDuration` 的截止 ctx，即便 Agent PTY hang 且客户端半开导致调用方 ctx 永不取消，超时也会触发 Server 端 `streams` map 与缓冲的回收，避免 entry + 64 帧缓冲 + goroutine 永驻。
- HTTP server 强制启用读/写/空闲超时（默认 30s/30s/120s，缺值兜底）并将请求体上限从 1GB 收紧到 32MB，慢速请求与异常上报不再能逐步占满容器内存。
- 反向 gRPC tunnel server 启用 keepalive（`MaxConnectionAge` 30 分钟、`MaxConnectionIdle` 5 分钟），周期性回收超龄与僵尸 Agent 连接；Agent 端已有重连循环，周期性重连无副作用。
- SQLite 启用 `auto_vacuum=INCREMENTAL` 与 `synchronous=NORMAL`；清理任务每小时执行 `wal_checkpoint(PASSIVE)` 与 `incremental_vacuum` 增量回收已删除页面，避免 WAL 与主库文件只增不减导致 RSS 随 mmap 缓慢上涨。非 SQLite 后端跳过该 PRAGMA 路径。
- 清理任务除 5 张时序表（CPU/内存/磁盘/网络/容器）外，还删除已关闭（`ended_at` 非空）且超过保留期的 `s_session` 记录，保留活跃与异常会话供排查。
- Collia 自动更新的 `failedAt` 冷却表在每次 Agent 注册（`Schedule`）时清理超过 24 小时的过期条目，避免长期离线 Agent 的 entry 在内存 map 中永久累积。
- 镜像删除先通过反向 tunnel 调用目标 Agent；Agent 返回成功后，Server 按 Agent ID 与镜像 ID 失效 `m_image` 缓存，使同一请求后的镜像列表刷新立即反映删除结果。RPC 失败时保留缓存。
- `interface IMailRepository (beacon/service/mail/repository/mail.go)`
- `interface IMailService (beacon/service/mail/service/mail.go)`
- `interface IDingTalkRepository (beacon/service/dingtalk/repository/dingtalk.go)`
- `interface IDingTalkService (beacon/service/dingtalk/service/dingtalk.go)`
- `interface IRouter (beacon/service/router.go)`
- `interface ITask (beacon/service/task/task.go)`
- `interface MailSender (beacon/service/task/task.go)`
- `interface Sender (beacon/service/dingtalk/client/client.go)`
- `interface Storer (beacon/pkg/auth/jwtauth/store.go)`
- `interface TokenInfo (beacon/pkg/auth/auth.go)`
- `interface TunnelMonitor (beacon/service/health/api/health.go)`
- `struct API (beacon/service/agent/agent.go)`
- `struct AccountAPI (beacon/service/account/api/account.go)`
- 另有其他导出模型/接口符号未展开；完整清单以源码为准。

## 使用注意事项

- 不跨越模块边界直接依赖内部路径或未导出符号。
- 修改跨模块公开 API、导出符号信号、配置或副作用边界时，同步更新本文件、根入口文档和相关 Domain Spec。
- 配置、凭据和外部副作用只记录语义边界，不记录真实敏感值。

## 验证命令

```bash
# beacon 模块提供完整 Taskfile.yml，命令清单：
cd beacon
task --list

# 以开发配置运行（http://127.0.0.1:8000）
task dev

# 生成 Wire 依赖注入代码（service/wire_gen.go）
task wire

# 本地二进制构建
task bin

# 构建 linux/amd64 Docker 镜像（写入本地 Docker）
task amd64

# 构建 linux/arm64 Docker 镜像（registry.cn-hangzhou.aliyuncs.com/amuluze/beacon-arm）
task arm64

# 隔离网络可直接覆盖可信基础镜像或 Alpine 软件源
docker buildx build --platform linux/arm64 -f Dockerfile \
  --build-arg APK_REPOSITORY=<trusted-alpine-mirror> \
  --build-arg BEACON_VERSION="$(cat VERSION)" \
  --build-arg COLLIA_VERSION="$(cat ../collia/VERSION)" \
  -t beacon:size-check -o type=docker ..

# 直接使用 Go 工具链
go test -race ./...
go vet ./...
```

## Taskfile Commands

- `task beacon`
- `task beacon-web`
- `task collia`
- `task default`
- `task website`
- `task website-server`
- `task website-web`
