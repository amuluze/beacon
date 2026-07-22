# Collia

`collia` 模块协作入口，汇总当前 workspace 的模块事实。

该模块当前角色：**monitoring agent module** — Beacon Server-Agent 架构中的 Agent 端，承担主机与 Docker 指标采集、本机执行能力，并通过反向 `gRPC tunnel` 主动连接 `beacon` Server，按 `agent_id` 让 Server 反向调用本机资源。

## 文档

模块协作入口以本 AGENTS.md 为准；项目级文档见根 AGENTS 文档地图。

- 模块实现事实见 [`.docs/modules/collia.md`](../.docs/modules/collia.md)。

## 开发导航

- 先读本文件确认模块边界，再读对应 `.docs/modules/` 文档获取当前实现事实；导出符号只作为入口线索，不自动等同跨模块公开 API。
- 涉及长期行为、不变量、状态或错误语义时，必须回到下方相关 Domain Spec；若现有 Domain Spec 不覆盖，应先补可验证约束。
- 代码变更后按本文件“开发命令”执行最小验证；跨模块、配置、接口或副作用变更还要运行项目级质量门禁或 `task sdd:refs`。
- 更新公开 API、配置键、事件、持久化格式或用户可见工作流时，同步维护相关项目文档并检查 `AGENTS.md` / `CLAUDE.md` 一致性。

## 模块路径

`collia`

## 关键目录

| 目录/文件 | 职责 |
|-----------|------|
| `collia/` | monitoring agent module: 主机与 Docker 指标采集、反向 gRPC tunnel 服务端、周期上报 |
| `collia/cmd/` | 命令行或进程入口（install / start / stop / status / remove） |
| `collia/config.yml` | Agent 配置文件（Server 地址、Agent ID、采集间隔、磁盘/网卡列表、上报 URL） |
| `collia/pkg/` | 可复用公共包集合 |
| `collia/pkg/psutil/` | gopsutil 指标采集辅助 |
| `collia/pkg/timectl/` | 时间控制辅助 |
| `collia/pkg/utils/` | 通用辅助函数（文件、字符串等） |
| `collia/script/` | 安装/卸载脚本 |
| `collia/service/` | 业务核心层：RPC 服务、采集任务、预检、报告客户端 |
| `collia/service/rpc/` | 反向 gRPC tunnel 服务端、调度器、Agent 生命周期 RPC |
| `collia/service/task/` | 采集任务调度（CPU/内存/磁盘/网络/Docker 指标） |
| `collia/service/report/` | 上报到 beacon Server 的 HTTP 客户端 |
| `collia/service/model/` | Go package `model`，源码 2，测试 0 |

## 依赖

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
- `gorm.io/gorm` `v1.25.12`

## 模块约束

- 仅通过公开接口与其他模块协作，不依赖其他模块内部实现细节。
- 修改公开 API、配置或副作用边界时，同步更新 `.docs/modules/` 中对应文档。
- 若模块承载长期领域语义，相关约束应在 `.specs/domain/` 中可追踪。

## 开发命令

```bash
# 列出全部任务
task --list

# 生成 Wire 依赖注入代码
task wire

# 构建 linux/amd64 二进制
task amd64

# 生成 Wire 代码并构建 linux/arm64 二进制
task arm64

# 直接使用 Go 工具链
cd collia
go test -race ./...
go vet ./...

# 安装到目标主机（需先有 Server 颁发的 install token）
collia install --token=<install-token>

# 启动/停止/查看状态/卸载
collia start
collia stop
collia status
collia remove
```
