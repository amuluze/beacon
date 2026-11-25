# 部署指南

## 环境要求

- Go 运行环境需满足各模块 `go.mod` 声明的版本（Go 1.25）。
- 控制台前端需安装 Node.js 与 pnpm 11.9 后运行对应 package scripts。
- Docker 用于构建与运行镜像；生产环境需要可信 HTTPS 对外地址。
- 运行前需要确认配置文件中的监听、数据库、Agent、凭据和路径语义。
- 容器内 Nginx 只监听 80，Web TLS 终止由上游反向代理承担，并透传 `X-Forwarded-Proto`（`beacon/nginx/nginx.conf` 的 `$beacon_forwarded_proto` 映射保留受信上游的原始协议）。

## 快速启动

```bash
# 构建控制台静态资源（Server 以 go:embed 内嵌）
task beacon-web:install
task beacon-web:build

# 开发模式直跑：development 环境不强制密钥（仅告警），本地数据与日志写入 beacon/data/
task beacon:dev

# 或构建镜像运行（生产密钥经环境变量注入）
docker build -f beacon/Dockerfile -t beacon:latest .
docker run -d --name beacon -p 8000:8000 beacon:latest

# 验证服务
curl http://127.0.0.1:8000/health
```

## 构建与镜像发布

- `beacon/Dockerfile` 提供镜像构建上下文，仓库根为构建上下文（需要 `go.work`、`common/`、`beacon/`、`collia/` 与 `beacon/web/`）。
- Dockerfile 在隔离的 Node 阶段执行控制台类型检查与生产构建；`.dockerignore` 排除本机 `beacon/web/dist`，因此镜像不会嵌入历史或未提交的本地资源。
- Vite 生产构建每次清空 `dist`，避免累积历史哈希文件。文件数量与体积随源码和工具链变化；发布验收时记录对应提交、平台与实际构建结果，不把历史数字写作当前固定值。
- Beacon 与 amd64/arm64 Collia 均以 `CGO_ENABLED=0`、`-trimpath`、`-ldflags="-s -w"` 构建，并分别注入 `beacon/VERSION` 与 `collia/VERSION`；镜像继续携带两个 Agent 架构。
- `COLLIA_VERSION` 同时写入运行时 `BEACON_COLLIA_VERSION`，使 Server 的目标版本与镜像内 `/app/downloads/collia/<arch>/collia` 始终来自同一次构建。
- 运行时基于 Alpine，保留 Nginx、Supervisor、CA 证书、`wget` 健康检查和 `/app/downloads/collia/<arch>/collia` 安装路径。
- 镜像入口 `/usr/local/bin/beacon-entrypoint`（来自 `beacon/scripts/docker-entrypoint.sh`）承担密钥校验，默认执行 `supervisord`；健康检查与进程编排保持由镜像负责。
- 构建参数 `NODE_IMAGE` 可替换为兼容 Node 22 的可信镜像代理，`RUNTIME_IMAGE` 可替换为兼容 Alpine 3.22 的可信基础镜像；`APK_REPOSITORY` 可替换 Alpine 软件源根地址。默认值均指向官方镜像/软件源，替换值必须由部署方验证可信性。
- 升级包签名：`beacon/Dockerfile` 通过 BuildKit secret 注入签名私钥（`docker build --secret id=collia_signing_key,src=<私钥路径> ...`），构建期 `collia/script/sign_release_packages.sh` 为 amd64/arm64 生成侧车签名并随镜像落到 `/app/downloads/collia/<arch>/collia.sig`；未挂载私钥的构建跳过签名（兼容模式），已固定公钥的 Agent 会拒绝未签名升级包。
- 镜像体积随源码、工具链与前端产物变化，不写作固定值；每次发布在同一平台与同一 Docker 引擎上用 `docker image inspect .Size` 记录实际数值。

```bash
# 国内或隔离网络可显式指定已审核的软件源；不指定时使用 Alpine 官方源。
docker buildx build --platform linux/arm64 -f beacon/Dockerfile \
  --build-arg APK_REPOSITORY=<trusted-alpine-mirror> \
  --build-arg BEACON_VERSION="$(cat beacon/VERSION)" \
  --build-arg COLLIA_VERSION="$(cat collia/VERSION)" \
  -t beacon:size-check -o type=docker .

# Taskfile 也允许通过同名变量传入已审核的软件源。
task beacon:arm64 APK_REPOSITORY=<trusted-alpine-mirror>
```

产品镜像每次发布同时推送到阿里云镜像仓库与 Docker Hub，同 tag 同内容（一次 buildx 构建挂两个 `-t`）：

- amd64：`registry.cn-hangzhou.aliyuncs.com/amuluze/beacon:<tag>` 与 `docker.io/amuluze/beacon:<tag>`
- arm64：`registry.cn-hangzhou.aliyuncs.com/amuluze/beacon-arm:<tag>` 与 `docker.io/amuluze/beacon-arm:<tag>`

执行 `task beacon:build` 完成两个镜像仓库的推送；Docker Hub 侧的 `amuluze/beacon` 与 `amuluze/beacon-arm` 仓库需预先创建，脚本不会自动建仓。

## 配置说明

| 配置文件 | 语义 |
|----------|------|
| `beacon/configs/config.dev.toml` | 开发环境配置，包含本地监听、数据库、Agent 安装与更新参数；`AgentUpdate.Enable=false` 避免本地连接触发自更新；`[CORS] Enable=true` 显式开启跨域（`AllowOrigins` 留空回退本地开发域白名单）；控制通道在 `data/control.sock`（Unix Socket）上监听 |
| `beacon/configs/config.toml` | 生产运行配置；`AgentUpdate` 默认启用，并配置升级超时、全局并发和失败冷却；默认不挂载 CORS（控制台与 API 同源部署），跨域访问可经环境变量 `BEACON_CORS_ENABLE` / `BEACON_CORS_ALLOW_ORIGINS`（逗号分隔白名单）或显式 `[CORS] Enable=true` 启用；其余可经环境变量覆盖的运行键：`BEACON_UPDATE_ENABLE` / `BEACON_UPDATE_URL` / `BEACON_UPDATE_CHECK_INTERVAL`（Server 自身更新检查）、`BEACON_CONTROL_SOCKET`（控制通道监听路径，Unix Socket 唯一承载）、`BEACON_AGENT_INSTALL_CONTROL_SOCKET`（下发 Agent 的宿主 Socket 路径）、`BEACON_DB_NAME`（数据库名） |
| `beacon/configs/model.conf` | Casbin RBAC 策略模型：定义请求方法与路径的匹配规则及内置角色（`admin`）的资源授权范围；Server 启动时按此模型初始化策略，Supervisor 以 `-m` 参数注入 |
| `beacon/nginx/nginx.conf` | 容器内 Nginx 主配置：加载 `mime.types` 与 `conf.d/*.conf` 站点片段，保留上游 `X-Forwarded-Proto`，设置访问日志与代理超时 |
| `beacon/nginx/conf.d/beacon.conf` | Web 入口站点配置：监听 80，`/` 与 `/ws/` 反代到 `127.0.0.1:8000`；`/ws/` 段透传 `Upgrade`/`Connection` 以支持 WebSocket 升级 |
| `beacon/nginx/mime.types` | Nginx MIME 类型表，保证静态资源 Content-Type 正确 |
| `beacon/supervisor/supervisord.conf` | 容器内进程编排入口：unix socket RPC、pid 与日志路径，并 include `conf.d/*.conf` |
| `beacon/supervisor/server.conf` | Beacon Server 进程定义：`/app/beacon run -c /app/configs/config.toml -m /app/configs/model.conf`，工作目录 `/app`，失败自动重启 |
| `beacon/supervisor/nginx.conf` | Supervisor 托管的 Nginx 进程定义：`/usr/sbin/nginx -c /etc/nginx/nginx.conf` |
| `beacon/scripts/docker-entrypoint.sh` | 镜像入口脚本：默认 `supervisord` 启动路径先校验关键密钥（弱值、长度、十六进制），通过后 `exec` CMD；显式覆盖为其他命令时直接执行，便于调试 |
| `beacon/web/tsconfig.json` | TypeScript 编译配置，用于类型检查或前端构建；检测到键/段：`@/*`, `allowImportingTsExtensions`, `baseUrl`, `compilerOptions`, `jsx`, `lib`, `module`, `moduleResolution` |
| `collia/config.yml` | Agent 运行时配置；键与段：`control.server`（必须为 `unix://` 的本机 Socket 路径）、`control.agent_id`、`control.join_token`、`log.output`/`log.level`/`log.rotation`/`log.max_age`、`task.interval`（采集间隔秒）、`task.disk.devices`/`task.ethernet.names`（采集范围）、`task.report.url`/`task.report.token`/`task.report.agent_id`（周期上报地址与共享 Token）、`db.dbtype`/`db.host`/`db.port`/`db.user`/`db.password`/`db.dbname`/`db.sslmode`（兼容段：旧版 Collia 本地库开关，当前版本解析后忽略）、`variables.image_tag`/`variables.host_prefix`/`variables.container_prefix` |

- 本节只列运行/构建配置文件；源码中的 Args/Reply/schema 类型不作为运行时配置项输出。
- 表中路径为仓库内路径；容器内对应位置为 `/app/configs`、`/etc/nginx`、`/etc/supervisor`、`/usr/local/bin/beacon-entrypoint` 与 `/app/downloads/collia/<arch>`。
- 数据库通过 `[Gorm] DBType` 选择、通过 `[DB]` 配置连接。Beacon 默认使用 SQLite；共享库中的其他驱动分支不等于相应数据库已完成部署验收。PostgreSQL 的 `SSLMode` 由 Beacon 传入共享数据库层；空值保留兼容默认 `disable`。远程数据库须显式选择所需的 TLS 校验模式，例如 `verify-full` 并部署可信根证书；不能把空值理解为已开启 TLS。SQLite 的 WAL、auto_vacuum、synchronous 与文件权限初始化只在 SQLite 分支执行；新库先设置auto_vacuum再开启WAL，已有库不自动执行VACUUM迁移，不能据此宣称历史库的auto_vacuum已改变。真实凭据不得入库，应通过已实现的环境变量绑定或部署期受保护配置文件注入。
- 容器运行时操作由 Collia 经 `common/docker` 调用 Docker Engine API（`client.FromEnv` + 版本协商），连接目标取自 Docker 客户端环境变量；所需权限由 Agent 运行账号授予，运行时删除失败时保留操作前视图并向上返回可定位错误，不伪造成功结果。
- 配置包含 `[Retention]` 段（`Days=7`、`AuditDays=90`）；`CleanupTask` 在取值为 0 时兜底为 7 天，即监控时序数据（CPU / 内存 / 磁盘 / 网络 / 容器）与已关闭会话默认保留 7 天。审计表 `s_audit` 按 `AuditDays`（默认 90 天）独立清理。

PostgreSQL专项测试位于 `beacon/service/config_test.go`：设置 `BEACON_TEST_POSTGRES_PORT` 和 `BEACON_TEST_POSTGRES_PASSWORD` 连接本机临时库（固定测试用户及库名为beacon_contract）。TLS用例另需 `BEACON_TEST_POSTGRES_ROOT_CERT`，可用 `BEACON_TEST_POSTGRES_WRONG_ROOT_CERT` 验证错误根证书被拒绝；测试根证书路径经PGSSLROOTCERT传给驱动。未设置专项环境时这些测试显式跳过，默认doc:verify通过不代表执行了真实PostgreSQL验收。

## 环境变量与敏感信息

| 类别 | 来源 | 文档策略 |
|------|------|----------|
| 运行时配置 | 配置文件、命令行参数、环境变量 | 记录名称和语义，不记录真实密钥 |
| 凭据/Token | 本地配置或外部平台 | 不写入生成文档；只描述加载边界 |
| 构建产物 | Makefile、Taskfile、Go build、前端 scripts | 记录命令和输出职责，不提交临时产物 |

- `BEACON_DATA_ENCRYPTION_KEY` 必须是解码后 32 字节的持久密钥（base64 或 64 位十六进制）；镜像入口脚本按 64 位十六进制校验。丢失或静默轮换会导致已加密 SMTP/钉钉/飞书/企微凭据无法解密。
- `BEACON_AUTH_SIGNING_KEY`、`BEACON_AGENT_INSTALL_TOKEN` 和 `BEACON_CONTROL_JOIN_TOKEN` 在生产模式必须为独立高熵值（镜像入口脚本要求至少 48 字符）；升级时保留已有值。
- `BEACON_AGENT_INSTALL_TOKEN_TTL`（配置文件 `[AgentInstall] TokenTTL`，秒）控制控制台申请安装令牌的有效期：0 或缺省为默认 1800 秒短时令牌，响应携带 `expires_at`；负值显式回退为返回静态令牌的兼容模式。短时令牌只可用于安装期端点（script/config/package），HTTP 上报回退端点只接受静态安装令牌（令牌分域）。
- `BEACON_AGENT_UPDATE_PUBLIC_KEY`（配置文件 `[AgentUpdate] PublicKey`）配置 Collia 发布签名公钥（base64 Ed25519）：配置后随安装 config 写入 Agent 的 `control.upgrade_public_key` 并在本地固定，此后升级包必须携带匹配的侧车签名（`collia.sig`），缺失或不匹配即拒绝替换（IU008）。发布方用 `collia/cmd/sign -gen <私钥路径>` 生成密钥对（私钥 0600、离线保管、不进镜像），用 `-key <私钥> <collia 二进制>` 生成 `<二进制>.sig`；Server 只透传签名，不持有私钥。公钥轮换使用逗号分隔的双钥（或 `upgrade_public_keys` 列表）形成重叠窗口：窗口期内新旧签名均通过，旧钥退役后从配置移除；升级 RPC 永远不下发公钥（RU013）。
- `BEACON_PUBLIC_BASE_URL` 在生产模式必须是无用户信息的 HTTPS URL；HTTP 或空值会在 Server 启动前被拒绝。
- `BEACON_COLLIA_VERSION` 由镜像构建注入，不应由部署方手工覆盖；`BEACON_AGENT_UPDATE_ENABLE` 是部署级门禁（关闭后控制台开关也无效），控制台「探针管理」页另有数据库级自动更新开关（默认开启，仅停自动调度），`BEACON_AGENT_UPDATE_TIMEOUT`、`BEACON_AGENT_UPDATE_MAX_CONCURRENT`、`BEACON_AGENT_UPDATE_RETRY_INTERVAL` 分别控制单次超时、全局并发和失败冷却。关闭自动调度不影响管理员手动触发；升级器仅在启用时校验目标版本。
- `BEACON_RETENTION_DAYS`、`BEACON_RETENTION_AUDIT_DAYS` 分别覆盖监控时序数据与审计日志的保留天数；未配置时取配置文件 `[Retention]` 段或内置默认（7 / 90 天）。
- `[App] ReadOnly` / `[App] NoAuth`（环境变量 `BEACON_READ_ONLY` / `BEACON_NO_AUTH`）是可选运行模式：
  - `ReadOnly=true` 且 `NoAuth=false`：保留登录，所有业务写请求被拒绝。
  - `NoAuth=true`：允许匿名访问；匿名会话只读，但可手动进入登录页并登录。登录后的会话按正常权限执行（Casbin 生效，管理员可写），退出登录后回到匿名只读。匿名容器终端保持关闭，登录后恢复；容器日志匿名可读。
  - 两者同时为 `true` 时，`ReadOnly` 只约束匿名会话，登录会话仍按正常权限使用。未登录的写请求统一返回 `403 read-only mode`；登录、刷新、改密和 Agent 上报按各自认证规则放行。
  - 两个开关默认关闭，由部署方按需通过环境变量注入。
- `BEACON_BOOTSTRAP_ADMIN_PASSWORD` 仅在首次建库（路由/角色表为空）时用于创建强制改密的 `admin`；至少 12 字符且覆盖四类字符中的至少三类，首次登录后强制改密。已初始化实例启动时若 `admin` 缺失，Server 以可定位错误拒绝初始化而非静默重建——恢复必须来自数据备份，或显式清空数据目录后重新安装。
- `BEACON_APP_ENV` 决定运行模式：`production` 强制 HTTPS 地址与强密钥，`development` 只告警不拒绝。
- `BEACON_DB_NAME` 覆盖数据库名（默认 SQLite 文件路径由 `[DB] DBName` 决定）；`BEACON_CORS_ENABLE` 与 `BEACON_CORS_ALLOW_ORIGINS` 控制跨域中间件与白名单。
- `BEACON_UPDATE_ENABLE` / `BEACON_UPDATE_URL` / `BEACON_UPDATE_CHECK_INTERVAL` 控制 Server 自身的版本检查（`[Update]` 段，默认 6h、最低 1h）；关闭后不执行任何远程版本检查。

忘记 `admin` 密码时使用内置命令（可在服务运行时执行；写入失败请检查数据库锁后重试）：

```bash
# 容器内路径 /app/configs/config.toml 对应仓库的 beacon/configs/config.toml
docker exec -it beacon ./beacon reset-admin -c /app/configs/config.toml
# 默认生成临时密码并打印；也可显式指定：--password '<新密码>' 或 --password-stdin < 文件
```

命令校验密码策略（≥12 字符、至少三类字符），bcrypt 入库后置 `must_change_password=1`（下次登录强制改密），并写入一条 `system` 审计；`admin` 不存在时拒绝执行，不会静默重建管理员。

## 单机部署：Unix Socket 控制通道

Beacon 仅支持单机部署：Beacon 与 Collia 同机运行，控制通道走共享数据目录内的 Unix Socket，不发布任何控制端口。

- 成对声明：`BEACON_CONTROL_SOCKET`（容器/进程内监听路径，`[Control] UnixSocket`）与 `BEACON_AGENT_INSTALL_CONTROL_SOCKET`（下发给 Agent 的宿主可见路径，`[AgentInstall] ControlSocket`）必须指向同一个 socket。`Control.Enable` 而缺 `UnixSocket`、`AgentInstall.Enable` 而缺 `ControlSocket`，或两者只配其一，Server 启动即拒绝（`beacon/service/transport_security.go`）。
- Socket 监听目录权限为 0700、Socket 文件 0600（`common/rpc/tunnel/server.go`）；仅本机 root 可连接。
- JoinToken 校验保留：Agent 注册仍需携带令牌（IU009）。
- Server 下发的 Agent 安装 config 使用 `control.server: unix://<宿主路径>`；Agent 侧对非 `unix://` 控制目标拒绝启动（`collia/service/transport_security.go`）。
- 从历史 TCP 版本（17000 + mTLS）升级时：更新服务端配置为上述成对声明，并**重新安装该机 Agent（或手工把 `/etc/collia/config.yml` 的 `control.server` 改为 `unix://…` 后重启 Agent 服务）**，否则旧 Agent 仍指向已停止监听的 17000 而离线。
- 安全边界：UDS 依赖文件系统权限（仅 root 可连接）与 JoinToken；只适用于同机部署。
- HTTP 上报回落：回落请求按系统信任库校验 Web 证书；自签 Web 部署如需保留回落能力，需把 Web CA 导入宿主信任库（或改用受信 CA 签发的 Web 证书）。主通道为 UDS 数据流，回落仅在其不可用时使用。

## 版本发布节奏

- 版本号遵循语义化版本：正式版 `vX.Y.Z`，测试版 `vX.Y.Z-beta.N`（N 从 1 递增）。
- 节奏：先发 `vX.Y.Z-beta.N`，功能稳定后发布同核正式版 `vX.Y.Z`。测试版必须基于下一个尚未发布的正式版本号——同核正式版已存在时，`beta.N` 会被判定为降级而拒绝升级。
- `beacon/VERSION` 与 `collia/VERSION` 同步推进，镜像 tag 取自 `beacon/VERSION`，`BEACON_COLLIA_VERSION` 取自 `collia/VERSION`。
- 测试版发布同样要求本仓打 tag 并创建 Release，且 Release 标记为 pre-release；正式版标记为 Latest。版本提示通道只接受正式版，写入测试版会被拒绝，正式版实例不会被提示升级到测试版。
- 升级语义：显式传入的 KEY=VAL 与环境变量覆盖优先生效，不被既有部署环境的旧值吞掉；升级前自动备份部署配置（`*.bak.<时间戳>`），健康检查未通过自动回滚并清理备份，中途失败保留备份供手动恢复；显式指定低于当前的版本被降级守卫拒绝。

## Agent 安装与自更新

- 安装入口：管理面板「安装 Agent」对话框生成 `curl -fsSL '<public-base-url>/api/v1/host/install?node=<id>' | sudo bash -s -- --token=<token>`；Token 由管理员接口 `POST /api/v1/host/get_install_token` 获取，脚本内下载安装包与配置的请求同样携带 `X-Install-Token`。
- 脚本行为：识别 amd64/arm64 → 下载对应安装包与节点配置（瞬时传输中断自动重试，最多 3 次）→ `collia install` 执行 preflight（探测磁盘与网卡并回写配置）→ 注册并启动系统服务；重复安装时保留已注册的 systemd 单元并按幂等成功返回；未检测到 systemd 时提示改用 `setsid nohup collia run`。
- 身份语义：`node` 为安全命名字符串（字母/数字开头，允许 `.` `_` `-`）；同一主机重复安装会覆盖配置并以新标识注册，旧标识记录保留在 Server，需在「设置 → 探针管理」核对。
- 自更新：控制台开关存 `s_agent_setting`，记录缺失默认开启；关闭只停自动调度，管理员仍可在「探针管理」行内手动升级；升级目标版本来自镜像注入的 `BEACON_COLLIA_VERSION`。
- 卸载：`collia remove`（别名 `uninstall`）先停止服务，再注销单元并清理二进制、配置、本地数据、日志；删除单元后执行 `systemctl daemon-reload`，未被 systemd 跟踪的残留 daemon 进程会被终止。

## 发布前检查

- 确认根 `AGENTS.md` 与 `CLAUDE.md` 内容一致。
- 确认 `.docs/MANIFEST.yml` 中 required 文档存在且章节完整。
- 确认 Domain Spec 没有混入文件路径、技术选型或伪代码。
- 确认构建、测试和质量门禁命令在目标环境可重复执行。
- `task doc:check` 检查文档、完整测试库存、约束证据台账和OpenAPI字段模型；历史验收失败不等于文档结构失败。
- `task doc:verify` 实际执行六个工作区 Go 模块与 `scripts/doc-schema` 工具模块的build、race测试、vet及门禁回归测试，任一失败返回非零；前端、镜像发布和安装端到端验收仍须按模块分别执行。
- GitHub Actions 的 `Documentation and frontend` workflow 在 PR 上分别执行文档/Go 验证与前端类型检查、测试和构建；Docker 与 PostgreSQL 专项仍需显式环境或本地任务。
- 比较镜像体积时必须在同一平台和同一 Docker 引擎上使用 `docker image inspect .Size`；同时检查两个 Collia 的版本、目标架构与静态链接属性。

## Taskfile Commands

- `task beacon`
- `task beacon-web`
- `task collia`
- `task default`
- `task doc:check`
- `task doc:verify`

## Package Scripts

- `beacon/web/package.json (beacon-web)`: `build`, `dev`, `lint`, `lint:fix`, `preview`, `test`, `test:coverage`, `test:run`, `ts`

## Workspace Modules
模块清单见 [项目盘点](project-analysis.md#workspace-modules)。
