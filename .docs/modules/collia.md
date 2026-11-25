# Collia 模块

## 模块职责

`collia` 是 Beacon Server-Agent 架构中的 **Agent 端**，承担以下职责：

- 主机与容器指标采集：CPU、内存、磁盘、网络使用率与趋势（基于 `shirou/gopsutil/v3`）。
- 本机 Docker 操作执行：依赖 `common/docker`（`github.com/amuluze/docker`），由 Server 通过反向 gRPC tunnel 调用。
- 反向 gRPC tunnel 的 Agent 端：启动后主动连接 Beacon Server 的本机 Unix Socket，按 `agent_id` 接收 Server 的反向调用；控制面不使用 mTLS，注册仍校验 JoinToken。
- 监控数据上报：优先经控制连接的数据流推送批次；HTTP `POST /api/v1/host/report` 保留为兼容回退（`task.report.url`）。
- 本地配置与服务注册：通过 `takama/daemon` 注册 systemd/init 服务，支持 `install/start/stop/status/remove` 子命令。重复执行 `install` 时若单元文件已存在，保留现有注册并按幂等成功返回（输出中性提示，不再打印 `[FAILED]`）；`remove` 先停止服务再注销单元，并在删除单元后执行 `systemctl daemon-reload`；对未被 systemd 跟踪的残留 `collia daemon` 进程，按 `/proc` 扫描做 SIGTERM→SIGKILL 兜底清理（避免旧版本卸载后以已删除二进制继续运行）。

`collia` **不**负责：

- 业务持久化：监控数据写入 Server；Agent 不维护本地业务库（历史 `s_container/s_image/s_network` 本地表与本地 SQLite 装配已移除，容器以 Docker 运行时为唯一事实源，见 IC007）。
- 用户鉴权与 RBAC：所有鉴权由 `beacon` Server 统一处理。
- 跨 Agent 协调：每个 collia 独立运行，仅与配置中 `control.server` 指定的 Server 通信。

- 相关源码包：`collia/cmd`、`collia/pkg/psutil`、`collia/pkg/utils`、`collia/service`、`collia/service/rpc`、`collia/service/task`、`collia/service/report`。

## 职责边界

| 职责 | 说明 | 是否负责 |
|------|------|----------|
| 主机与容器指标采集 | 基于 gopsutil 与 Docker Engine API | yes |
| 反向 gRPC tunnel | Unix Socket 客户端，主动连接 Server 接收反向调用 | yes |
| 监控数据上报 | 控制连接数据流优先，HTTP POST 兼容回退 | yes |
| 业务持久化 | 无本地业务库，容器事实源为 Docker 运行时 | no |
| 用户鉴权与权限 | 由 beacon Server 统一处理 | no |
| HTTP API 注册 | 不暴露 HTTP 入口（仅 Server 端） | no |
| 跨模块协调 | 仅通过 `common/*` 公开接口协作 | conditional |
| 其他模块内部状态 | 不直接读写其他模块私有实现 | no |

## 模块事实

- 模块路径：`collia`
- Go 版本：`1.25.0`
- Go 源码文件：28
- Go 测试文件：17
- 包名：`main`, `psutil`, `report`, `rpc`, `service`, `task`, `utils`
- 推断角色：monitoring agent module — 主机与 Docker 指标采集、反向 gRPC tunnel 客户端、周期上报到 beacon Server

## 接口与符号信号

本节记录源码扫描得到的导出符号或前端入口信号；它们是候选接口线索，不等同于跨模块公开 API。跨模块调用仍以 facade、模块文档和 Domain Spec 约束为准。

- `interface ITask (collia/service/task/task.go)`
- `func BuildInjector (collia/service/wire.go)`
- `func BuildInjector (collia/service/wire_gen.go)`
- `func CleanupBackup (collia/service/rpc/lifecycle.go)`
- `func DefaultRestartFn (collia/service/rpc/lifecycle.go)`
- `func LoadConfig (collia/service/preflight.go)`
- `func MergeProbeIntoConfig (collia/service/preflight.go)`
- `func NewClient (collia/service/report/client.go)`
- `func NewConfig (collia/service/config.go)`
- `func NewDispatcher (collia/service/rpc/dispatcher.go)`
- `func NewInjector (collia/service/injector.go)`
- `func NewLogger (collia/service/logger.go)`
- `func NewRPCServer (collia/service/rpc.go)`
- `func NewService (collia/service/rpc/rpc.go)`
- 另有 63 个导出符号未展开；请用 `rg`、codegraph 或 IDE 按需查看完整清单。

## 依赖关系

- `github.com/amuluze/amutool/logger` `v0.0.0-20240821104128-caed9cc0d402`
- `github.com/amuluze/amutool/timex` `v0.0.0-20250508153823-fe9a5de55958`
- `github.com/amuluze/docker` `v0.0.0-20240822095446-429928f7463e`
- `github.com/creack/pty` `v1.1.24`
- `github.com/docker/docker` `v27.2.1+incompatible`
- `github.com/google/wire` `v0.6.0`
- `github.com/patrickmn/go-cache` `v2.1.0+incompatible`
- `github.com/shirou/gopsutil/v3` `v3.24.5`
- `github.com/spf13/viper` `v1.19.0`
- `github.com/takama/daemon` `v1.0.0`
- `google.golang.org/grpc` `v1.81.1`
- 间接依赖 60 个，完整列表以模块 `go.mod` 为准。

## 配置与入口信号

- 配置：`control.server`（必须为 `unix://` 目标）/ `control.agent_id` / `control.join_token` at `collia/config.yml`
- 配置：`log.output` / `log.level` / `log.rotation` / `log.max_age` at `collia/config.yml`
- 配置：`task.interval` / `task.disk.devices` / `task.ethernet.names` at `collia/config.yml`
- 配置：`task.report.url` / `task.report.token` / `task.report.agent_id` at `collia/config.yml`
- 配置：`db.dbtype` / `db.host` / `db.port` / `db.user` / `db.password` / `db.dbname` / `db.sslmode` at `collia/config.yml`（兼容段：供旧版 Collia 打开本地库，当前版本解析后忽略）
- 配置：`variables.image_tag` / `variables.host_prefix` / `variables.container_prefix` at `collia/config.yml`
- 配置：`Config` in `collia/service/config.go`

## 数据模型与状态

- 本节记录当前实现事实，可追溯到源码与测试断言；长期行为约束以 `.specs/domain/` 为准，本节的「必须 / 禁止」表述不构成约束来源。
- 容器清单使用 Docker `All: true` 语义，`created`、`exited` 等非运行状态仍会进入上报；仅 `running` 容器读取实时 CPU/内存 stats，其他状态使用可 JSON 序列化的零值指标，避免单个容器阻断整批报告。
- 定时上报在主机与 Docker 指标全部采集完成后确定批次时间戳，使上报时间接近真实可发送时刻；HTTP client 只通过 `X-Install-Token` 发送共享 Token，不把凭据放入 query。
- CPU 批次同时上报百分比与逻辑核数（`psutil.GetCPUCount()` → `cpu.Counts(true)`，字段 `cpu_count`）；核数采集失败只丢弃该字段（记 0），不影响百分比上报。
- 采集在途限制为单槽（`TimedTask.tryExecute`）：上一轮未完成时新 tick 非阻塞丢弃，不叠加 goroutine 与 Docker API 压力（`collia/service/task.go`）。
- 控制通道目标仅支持 `unix://`（单机部署）：`validateTransportSecurity` 拒绝任何非 Unix Socket 目标，装配层以 insecure 凭据直连本机 socket 并记录日志；JoinToken 校验保留（`collia/service/transport_security.go`、`collia/service/rpc.go`）。
- 升级验签 fail-closed（`collia/service/rpc/lifecycle.go`）：`control.upgrade_public_key` 固定后缺签或错签即拒绝替换并保留旧二进制；未固定公钥但收到签名同样拒绝；仅双方都未启用签名链路时按 SHA-256 兼容放行。公钥支持集合（`control.upgrade_public_keys`）形成轮换重叠窗口：任一公钥匹配即放行，非法编码或长度导致启动失败。
- 发布签名工具 `collia/cmd/sign`：`-gen` 生成 Ed25519 密钥对并打印公钥（私钥 0600、离线保管），`-key` 为任意文件生成 `<file>.sig` 侧车签名；镜像构建期由 `collia/script/sign_release_packages.sh` 用 BuildKit secret 私钥自动为 amd64/arm64 生成侧车签名，未挂载私钥时显式跳过。Server 只透传签名，镜像内不含私钥。
- 容器批次上报每个容器实际挂载的网络名（Docker `NetworkSettings.Networks` 中首个带 IP 的网络，字段 `network`），供 Server 侧按网络筛选；「使用中」镜像数量由 Server 按同一批次容器引用计算，Collia 不上报该字段。

- 容器名唯一性与删除以 Docker 运行时为唯一事实源：创建重名由 errdefs 冲突判定并归一为 `container <name> already exists`（`common/docker.IsNameConflict` + `collia/service/rpc/container.go`），Agent 不维护本地容器记录；冲突判定有单测锁定（`common/docker/conflict_test.go`）。
- 更新容器采用「改名备份 → 新建 → 回滚」流程：请求未显式传入 `commands` 时沿用旧容器的 `Config.Cmd`（与环境变量、端口、卷、标签同口径），避免更新把容器启动命令重置为镜像默认命令；显式传入时以请求为准。
- 反向 tunnel 注册帧携带编译时版本以及运行时 OS/架构；这些元数据在重连时保持不变，供 Server 判断是否需要版本收敛。
- `UpgradeAgent` 从 Server 安装包端点下载目标二进制，要求非空 SHA256 并在校验后执行带备份的原子替换；同进程并发升级由互斥拒绝，替换完成后通过服务管理器重启并以新版本重新注册。
- 反向控制 RPC 不注册 DNS、系统时间、时区、主机重启或关机操作；Docker 镜像源配置继续作为独立 Docker 控制能力提供。
- 生产配置要求 HTTPS 上报地址与本机 Unix Socket 控制目标：控制通道无 TLS 配置面，非 `unix://` 目标启动即拒绝；上报地址非 HTTPS 时拒绝携带 Token（`service/transport_security.go`）。启动预检（`preflight.go`）覆盖目录可写性、控制通道 socket 可达性与探测充分性。
- amd64/arm64 发布二进制统一关闭 CGO，并使用 `-trimpath -ldflags "-s -w"` 裁剪本机路径、符号表和 DWARF；版本仍由 `collia/VERSION` 注入。2026-07-24 验证产物约为 amd64 30 MiB、arm64 28 MiB，均为静态、已剥离 Linux ELF。
- `struct CPUReport (collia/service/task/task.go)`
- `struct Client (collia/service/report/client.go)`
- `struct Config (collia/service/config.go)`
- `struct ContainerReport (collia/service/task/task.go)`
- `struct Control (collia/service/config.go)`
- `struct Disk (collia/service/config.go)`
- `struct DiskReport (collia/service/task/task.go)`
- `struct Dispatcher (collia/service/rpc/dispatcher.go)`
- 另有 29 个导出模型/接口符号未展开；完整清单以源码为准。

## 使用注意事项

- 不跨越模块边界直接依赖内部路径或未导出符号。
- 修改跨模块公开 API、导出符号信号、配置或副作用边界时，同步更新本文件、根入口文档和相关 Domain Spec。
- 配置、凭据和外部副作用只记录语义边界，不记录真实敏感值。

## 验证命令

```bash
cd collia
# 列出全部任务
task --list

# 生成 Wire 依赖注入代码
task wire

# 构建 linux/amd64 二进制
task amd64

# 生成 Wire 代码并构建 linux/arm64 二进制
task arm64

# 直接构建（不走 Taskfile）
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=$(cat VERSION)" -o collia ./cmd
```

## Taskfile Commands

- `task beacon`
- `task beacon-web`
- `task collia`
- `task default`
