# API 路由索引

## 范围

本文档汇总 `beacon` Server HTTP 服务在源码中静态注册的全部路由信号，包含路径、HTTP 方法、处理器、用途描述与保护策略。机器可读契约以 `.docs/api/openapi.yml` 为准，运行时行为仍以 handler、schema 与测试证据为准。

总览：

- `beacon` Server：77 条业务路由（`/api/v1/*`）+ `/health`、`/ready` 探针 + 2 个 WebSocket 端点（`/ws/terminal`、`/ws/:id`）

## 探针与中间件

### `beacon` 服务（`beacon/service/router.go`、`beacon/service/app.go`）

| 方法 | 路径 | 处理器 | 保护 | 用途 |
|------|------|--------|------|------|
| GET | `/health` | `healthProbe.Liveness` | 无（commonAuthSkipperPaths） | 进程存活探针 |
| GET | `/ready` | `healthProbe.Readiness` | 无 | 就绪探针（依赖 DB ping） |
| GET | `/ws/terminal` | `websocket.New(terminalHandler.Handle)` | ws upgrade + 同源校验 + HttpOnly access Cookie（`WebSocketUserAuthMiddleware`，Auth.Enable 时；未启用时回退到无 Cookie 认证的注册分支） | Agent PTY 终端会话（Web → Server → Agent） |
| GET | `/ws/:id` | `websocket.New(loggerHandler.Handler)` | ws upgrade + 同源校验 + HttpOnly access Cookie（`WebSocketUserAuthMiddleware`，Auth.Enable 时；未启用时回退到无 Cookie 认证的注册分支） | 容器日志流（Web → Server → Agent） |
| USE | `/` | `filesystem.New` | 全局 | 托管 `beacon/web/dist` 静态资源 |
| USE | `/api` | `middleware.CSRFMiddleware` | Auth.Enable 时 | 状态变更请求双提交 CSRF 校验 |
| USE | `/api/v1/auth/login` | `limiter.New(Max=10/60s, Key=IP)` | 限流 | 登录限流（防爆破） |
| USE | `/api/v1/host/report` | `limiter.New(Max=60/60s, Key=Agent-ID，缺失时回退 IP)` | 限流 | Agent 上报限流 |
| USE | `/ws` | ws upgrade guard | 全局 | 仅放行 WebSocket Upgrade 请求 |
| USE | 全局 | `securityHeadersMiddleware` | 全局 | 安全响应头（CSP 等） |
| USE | 全局 | `cors.New` | 全局 | 跨域（按 `CORS` 配置白名单） |
| USE | 全局 | 强制 HTTPS 中间件 | 仅生产环境 | 拒绝非 HTTPS Web 请求（`/health`、`/ready` 豁免） |
| USE | 全局 | `compress.New` | 全局 | 响应压缩 |
| USE | `/debug/pprof/*` | `pprof.New` | 仅 development | 性能剖析端点（空值/未知环境默认禁用） |
| USE | 全局 | Agent-ID 注入中间件 | 全局 | 把 `X-Agent-ID` / `agent_id` 写入 `contextx` |
| USE | 全局 | `middleware.PanicMiddleware` | 全局 | panic 恢复 |
| USE | 全局 | `middleware.StackMiddleware` | 全局 | 调用栈日志 |
| USE | 全局 | `middleware.GlobalRateLimitMiddleware` | 可选（RateLimit.Enable） | 全局限流（仅 `/api/` 路由，按 IP） |
| USE | 全局 | `middleware.UserAuthMiddleware` | 可选（Auth.Enable） | JWT 鉴权 |
| USE | 全局 | `middleware.CasbinMiddleware` | 可选（Casbin.Enable） | RBAC 鉴权 |

- 匿名只读与全局只读模式的保护语义（`beacon/service/router.go:147/:156-172`）：
  - `NoAuth=true` 时改挂 `OptionalUserAuthMiddleware` / `OptionalCasbinMiddleware`：匿名会话不要求 JWT，表内 `JWT + Casbin` 退化为「匿名可访问但只读」；已登录会话仍走完整 JWT + Casbin 并恢复正常权限（管理员可写）。末尾的 `ReadOnlyMiddleware(true)` 把匿名写请求统一拒绝为 `403 read-only mode`，登录会话不受影响。
  - `NoAuth=false` 且 `ReadOnly=true` 时，`ReadOnlyMiddleware(false)` 先于认证/CSRF 挂载：所有业务写请求（无论是否登录）稳定得到 `403 read-only mode`，读路径不受影响；登录、刷新、改密和 Agent 上报按各自认证规则放行。下表保护列的 `JWT + Casbin` 在该模式下对写路由实际表现为「JWT + Casbin + 只读拒绝」。
  - 两种模式都不改变 `commonAuthSkipperPaths` 与安装/上报端点的既有令牌规则。

`commonAuthSkipperPaths`（不要求登录/Casbin）：`/health`、`/ready`、`/api/v1/index/index`、`/api/v1/auth/csrf`、`/api/v1/auth/login`、`/api/v1/auth/token_update`、`/api/v1/host/install`、`/api/v1/host/report`。

## `beacon` 业务路由（`/api/v1/*`）

### Index（健康测试）

| 方法 | 路径 | 处理器 | 保护 | 用途 |
|------|------|--------|------|------|
| GET | `/api/v1/index/index` | inline `SendString("hello world")` | 无（commonAuthSkipperPaths） | 路由连通性测试 |

### 账户管理（`user` / `role` / `resource`）

| 方法 | 路径 | 处理器 | 保护 | 用途 |
|------|------|--------|------|------|
| GET | `/api/v1/user/user_query` | `accountAPI.UserQuery` | JWT + Casbin | 查询用户列表 |
| POST | `/api/v1/user/user_create` | `accountAPI.UserCreate` | JWT + Casbin | 创建用户 |
| POST | `/api/v1/user/user_update` | `accountAPI.UserUpdate` | JWT + Casbin | 更新用户 |
| POST | `/api/v1/user/user_delete` | `accountAPI.UserDelete` | JWT + Casbin | 删除用户 |
| GET | `/api/v1/role/role_query` | `accountAPI.RoleQuery` | JWT + Casbin | 查询角色列表 |
| POST | `/api/v1/role/role_create` | `accountAPI.RoleCreate` | JWT + Casbin | 创建角色 |
| POST | `/api/v1/role/role_update` | `accountAPI.RoleUpdate` | JWT + Casbin | 更新角色 |
| POST | `/api/v1/role/role_delete` | `accountAPI.RoleDelete` | JWT + Casbin | 删除角色 |
| GET | `/api/v1/resource/resource_query` | `accountAPI.ResourceQuery` | JWT + Casbin | 查询资源（菜单/接口）列表 |

### Agent 列表与更新

| 方法 | 路径 | 处理器 | 保护 | 用途 |
|------|------|--------|------|------|
| GET | `/api/v1/system/version` | `systemAPI.Version` | JWT + Casbin | 查询 Beacon Server 当前构建版本与最近一次更新检查状态（`version` / `update_available` / `latest_version`）；Server 不自更新，是否升级由管理员手动执行 |
| GET | `/api/v1/agent/list` | `agentAPI.List` | JWT + Casbin | 查询可用 Agent 列表；`status` 按当前 Server 进程的实时 tunnel 注册表投影，不直接信任历史持久化连接标记 |
| GET | `/api/v1/agent/update_setting` | `agentAPI.UpdateSetting` | JWT + Casbin | 查询控制台的 Agent 自动更新开关；记录缺失时按产品默认值 `true` 返回，同时返回部署级门禁 `deployment_enabled` |
| POST | `/api/v1/agent/update_setting` | `agentAPI.SetUpdateSetting` | JWT + CSRF + Casbin（管理员专属资源） | 写入 `{auto_update}` 开关；只影响 Agent 注册后的自动调度，关闭后手动升级与部署级门禁不受影响 |
| POST | `/api/v1/agent/upgrade` | `agentAPI.Upgrade` | JWT + CSRF + Casbin（管理员专属资源） | 对请求体 `{agent_id}` 指定的 Agent 手动触发镜像内置 Collia 更新；可绕过自动失败冷却，不绕过单 Agent 互斥和全局并发限制 |

### 认证与权限

| 方法 | 路径 | 处理器 | 保护 | 用途 |
|------|------|--------|------|------|
| GET | `/api/v1/auth/csrf` | `authAPI.CSRF` | 无（设置 CSRF Cookie） | 为首次登录提供双提交 CSRF Cookie |
| POST | `/api/v1/auth/login` | `authAPI.Login` | CSRF + 登录限流 | 登录并设置 HttpOnly access/refresh Cookie，只返回强制改密状态 |
| POST | `/api/v1/auth/logout` | `authAPI.Logout` | Cookie JWT + CSRF | 注销会话并清除 Cookie |
| POST | `/api/v1/auth/pass_update` | `authAPI.PassUpdate` | Cookie JWT + CSRF + Casbin | 修改当前用户密码并使旧会话失效 |
| POST | `/api/v1/auth/token_update` | `authAPI.TokenUpdate` | refresh Cookie + CSRF | 轮换会话 Cookie，响应体不暴露 JWT |
| GET | `/api/v1/auth/user_info` | `authAPI.UserInfo` | Cookie JWT + Casbin | 查询当前用户、权限与强制改密状态 |

### 容器与镜像（`container`）

> 本节全部路由均要求目标 Agent 标识（`X-Agent-ID` header 或 `agent_id` query）：读路径（列表/查询/统计）由 repository 层 `RequireAgentID` 强制，控制写路径（启停/删除/拉取等经 tunnel 转发的操作）由 `pkg/rpc` 的 `ResolveAgentID` 强制。下表保护列不再重复标注。

| 方法 | 路径 | 处理器 | 保护 | 用途 |
|------|------|--------|------|------|
| GET | `/api/v1/container/version` | `containerAPI.Version` | JWT + Casbin | 获取 Docker 版本信息 |
| GET | `/api/v1/container/containers` | `containerAPI.ContainerList` | JWT + Casbin | 获取容器列表 |
| GET | `/api/v1/container/usage` | `containerAPI.Usage` | JWT + Casbin | 获取容器资源使用情况 |
| GET | `/api/v1/container/summary` | `containerAPI.ContainerSummary` | JWT + Casbin | 获取容器工作区概览 |
| POST | `/api/v1/container/container_create` | `containerAPI.ContainerCreate` | JWT + Casbin | 创建容器 |
| POST | `/api/v1/container/container_update` | `containerAPI.ContainerUpdate` | JWT + Casbin | 编辑容器 |
| POST | `/api/v1/container/container_start` | `containerAPI.ContainerStart` | JWT + Casbin | 启动容器 |
| POST | `/api/v1/container/container_stop` | `containerAPI.ContainerStop` | JWT + Casbin | 停止容器 |
| POST | `/api/v1/container/container_restart` | `containerAPI.ContainerRestart` | JWT + Casbin | 重启容器 |
| POST | `/api/v1/container/container_remove` | `containerAPI.ContainerRemove` | JWT + Casbin | 删除容器 |
| GET | `/api/v1/container/images` | `containerAPI.ImageList` | JWT + Casbin | 获取镜像列表 |
| POST | `/api/v1/container/image_remove` | `containerAPI.ImageRemove` | JWT + Casbin | 删除镜像 |
| POST | `/api/v1/container/images_prune` | `containerAPI.ImagesPrune` | JWT + Casbin | 清理虚悬镜像 |
| POST | `/api/v1/container/image_pull` | `containerAPI.ImagePull` | JWT + Casbin | 拉取镜像 |
| POST | `/api/v1/container/image_import` | `containerAPI.ImageImport` | JWT + Casbin | 导入镜像 |
| POST | `/api/v1/container/network_create` | `containerAPI.NetworkCreate` | JWT + Casbin | 创建网络 |
| POST | `/api/v1/container/network_delete` | `containerAPI.NetworkDelete` | JWT + Casbin | 删除网络 |
| GET | `/api/v1/container/networks` | `containerAPI.NetworkList` | JWT + Casbin | 获取网络列表 |
| GET | `/api/v1/container/get_docker_registry_mirrors` | `containerAPI.GetDockerRegistryMirrors` | JWT + Casbin | 获取 Docker 镜像加速设置 |
| POST | `/api/v1/container/set_docker_registry_mirrors` | `containerAPI.SetDockerRegistryMirrors` | JWT + Casbin | 更新 Docker 镜像加速设置 |

### 主机与 Agent 控制（`host`）

| 方法 | 路径 | 处理器 | 保护 | 用途 |
|------|------|--------|------|------|
| GET | `/api/v1/host/install` | `Router.AgentInstallScript` | 无 JWT/Casbin（仅受 `AgentInstall.Enable` 门禁） | 获取 Collia 安装脚本；`node` 为安全命名字符串（字母/数字开头，允许 `.` `_` `-`），安装 Token 由执行参数 `--token` 传入 |
| GET | `/api/v1/host/install/package` | `Router.AgentInstallPackage` | `X-Install-Token`（或 Bearer；JWT/Casbin 前缀放行） | 按 `arch=amd64|arm64` 下载 Collia 安装包 |
| GET | `/api/v1/host/install/config` | `Router.AgentInstallConfig` | `X-Install-Token`（或 Bearer；JWT/Casbin 前缀放行） | 下载指定 `node` 的 Collia 配置 |
| POST | `/api/v1/host/get_install_token` | `Router.AgentInstallToken` | JWT + Casbin | 申请 Collia 安装令牌：默认签发短时令牌（`AgentInstall.TokenTTL`，0=默认 1800 秒，负值回退静态令牌），响应含 `expires_at` |
| POST | `/api/v1/host/report` | `reportSvc.HandleReport` | `X-Install-Token` + 限流（60/60s，Key 为 Agent-ID、缺失时回退 IP） | Agent 上报监控数据（按批次原子写入）；Token 未配置返回 503，不接受 query token。**已弃用**：主上报通道改为 Unix Socket 控制连接的数据流（同连接独立流），本接口作为无法建立数据流的旧 Agent 兼容回落，保留若干版本后移除 |
| GET | `/api/v1/host/host_info` | `hostAPI.HostInfo` | JWT + Casbin + Agent-ID | 主机信息（host/host_info via tunnel） |
| GET | `/api/v1/host/cpu_info` | `hostAPI.CPUInfo` | JWT + Casbin + Agent-ID | CPU 信息 |
| GET | `/api/v1/host/mem_info` | `hostAPI.MemInfo` | JWT + Casbin + Agent-ID | 内存信息 |
| GET | `/api/v1/host/disk_info` | `hostAPI.DiskInfo` | JWT + Casbin + Agent-ID | 磁盘信息 |
| GET | `/api/v1/host/cpu_trending` | `hostAPI.CPUUsage` | JWT + Casbin + Agent-ID | CPU 使用率趋势（5/10/30min） |
| GET | `/api/v1/host/mem_trending` | `hostAPI.MemUsage` | JWT + Casbin + Agent-ID | 内存使用率趋势 |
| GET | `/api/v1/host/disk_trending` | `hostAPI.DiskUsage` | JWT + Casbin + Agent-ID | 磁盘使用率趋势 |
| GET | `/api/v1/host/net_trending` | `hostAPI.NetUsage` | JWT + Casbin + Agent-ID | 网络使用率趋势 |
| GET | `/api/v1/host/file_search` | `hostAPI.FilesSearch` | JWT + Casbin + Agent-ID | 文件搜索 |
| POST | `/api/v1/host/file_download` | `hostAPI.FileDownload` | JWT + Casbin + Agent-ID | 文件下载 |
| POST | `/api/v1/host/file_delete` | `hostAPI.FileDelete` | JWT + Casbin + Agent-ID | 文件删除 |
| POST | `/api/v1/host/file_create` | `hostAPI.FileCreate` | JWT + Casbin + Agent-ID | 文件创建 |
| POST | `/api/v1/host/folder_create` | `hostAPI.FolderCreate` | JWT + Casbin + Agent-ID | 文件夹创建 |
### 审计

| 方法 | 路径 | 处理器 | 保护 | 用途 |
|------|------|--------|------|------|
| GET | `/api/v1/audit/query` | `auditAPI.AuditQuery` | JWT + Casbin | 查询审计日志 |

### 邮件告警配置

| 方法 | 路径 | 处理器 | 保护 | 用途 |
|------|------|--------|------|------|
| POST | `/api/v1/mail/mail_create` | `mailAPI.MailCreate` | JWT + Casbin | 创建邮件告警配置 |
| POST | `/api/v1/mail/mail_update` | `mailAPI.MailUpdate` | JWT + Casbin | 更新邮件告警配置 |
| GET | `/api/v1/mail/mail_query` | `mailAPI.MailQuery` | JWT + Casbin | 查询邮件告警配置 |
| POST | `/api/v1/mail/mail_test` | `mailAPI.MailTest` | JWT + Casbin | 测试邮件告警通道 |

### 钉钉告警配置

| 方法 | 路径 | 处理器 | 保护 | 用途 |
|------|------|--------|------|------|
| GET | `/api/v1/dingtalk/query` | `dingTalkAPI.Query` | JWT + Casbin（管理员专属资源） | 查询钉钉告警启用状态、凭据掩码和配置状态 |
| POST | `/api/v1/dingtalk/update` | `dingTalkAPI.Update` | JWT + Casbin（管理员专属资源） | 更新官方群机器人 Webhook、可选签名密钥和提醒范围 |
| POST | `/api/v1/dingtalk/test` | `dingTalkAPI.Test` | JWT + Casbin（管理员专属资源） | 使用已保存配置发送钉钉测试消息 |

### 飞书告警配置

| 方法 | 路径 | 处理器 | 保护 | 用途 |
|------|------|--------|------|------|
| GET | `/api/v1/feishu/query` | `feishuAPI.Query` | JWT + Casbin（管理员专属资源） | 查询飞书告警启用状态、凭据掩码和配置状态 |
| POST | `/api/v1/feishu/update` | `feishuAPI.Update` | JWT + Casbin（管理员专属资源） | 更新官方群机器人 Webhook、可选签名密钥 |
| POST | `/api/v1/feishu/test` | `feishuAPI.Test` | JWT + Casbin（管理员专属资源） | 使用已保存配置发送飞书测试消息 |

### 企业微信告警配置

| 方法 | 路径 | 处理器 | 保护 | 用途 |
|------|------|--------|------|------|
| GET | `/api/v1/wecom/query` | `wecomAPI.Query` | JWT + Casbin（管理员专属资源） | 查询企业微信告警启用状态、凭据掩码和配置状态 |
| POST | `/api/v1/wecom/update` | `wecomAPI.Update` | JWT + Casbin（管理员专属资源） | 更新官方群机器人 Webhook 和提醒范围 |
| POST | `/api/v1/wecom/test` | `wecomAPI.Test` | JWT + Casbin（管理员专属资源） | 使用已保存配置发送企业微信测试消息 |

### 告警阈值与消息标题

| 方法 | 路径 | 处理器 | 保护 | 用途 |
|------|------|--------|------|------|
| POST | `/api/v1/alarm/alarm_update` | `alarmAPI.AlarmUpdate` | JWT + Casbin | 更新告警阈值配置 |
| GET | `/api/v1/alarm/alarm_query` | `alarmAPI.AlarmQuery` | JWT + Casbin | 查询告警阈值配置 |
| GET | `/api/v1/alarm/message_setting` | `alarmAPI.AlarmMessageSettingQuery` | JWT + Casbin | 查询邮件与钉钉共用的全局告警消息标题 |
| POST | `/api/v1/alarm/message_setting_update` | `alarmAPI.AlarmMessageSettingUpdate` | JWT + Casbin | 更新全局告警消息标题，拒绝空白、换行和超过 128 字符的输入 |

## 维护约束

- 新增或修改路由时，同步检查请求参数、响应 schema、权限和错误语义。
- 该文件是 OpenAPI 的人类可读索引，不替代 `.docs/api/openapi.yml`；新增或变更 API 时必须同步维护两者及对应契约测试。
- 容器与主机控制类路由必须携带 Agent-ID（`X-Agent-ID` header 或 `agent_id` query），由全局中间件注入到 `contextx`，下游必须显式校验。这两个入口构成稳定协议表面：第三方客户端以任一形式提交均须被接受，移除任一入口属于破坏性变更。
- 路由新增、删除或修改方法时，同步更新本文件与 OpenAPI，并运行 `task doc:check`；仅当调用边界变化时更新架构文档。
