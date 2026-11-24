# Common Docker 模块

## 模块职责

`common-docker` 是独立的 Go module（`github.com/amuluze/docker`），提供 Docker Engine API 的薄封装，供 `collia` Agent 直接使用、`beacon` Server 通过反向 tunnel 间接复用。承担以下职责：

- Docker 管理抽象：`Manager` 结构 + `IManager` 接口，统一容器、镜像、网络三类资源的操作入口（避免各业务侧直接依赖 `docker/docker/client`）。
- 容器摘要结构：`ContainerSummary`、`PortMapping` 等，用于跨模块传递容器视图；屏蔽 Engine API 的私有字段。
- 镜像摘要结构：`ImageSummary`，按需提供仓库、标签、大小、创建时间等。
- 网络摘要结构：`NetworkSummary`、`SubNetworkConfig`，承载网络拓扑与子网/IPAM 配置。
- 版本探测：`Version` 结构 + `Version()` 方法，用于在控制台展示 Docker 版本。
- Docker Engine SDK 兼容性：当前固定 `docker/docker v27.2.1+incompatible`，依赖 `go-connections`、`libcompose`、`tidwall/gjson`。

`common-docker` **不**负责：

- 业务持久化：Docker 状态变更不入库，由 `beacon`/`collia` 的业务层落库。
- 鉴权与权限：所有操作由上层按 Agent-ID 校验权限。
- HTTP 入口：纯 Go 库，不暴露任何路由。

相关源码包：`common/docker`（全部代码位于该目录）。

## 职责边界

| 职责 | 说明 | 是否负责 |
|------|------|----------|
| Docker Engine API 薄封装 | `IManager` 接口 + `Manager` 实现 | yes |
| 容器/镜像/网络摘要结构 | `ContainerSummary`、`ImageSummary`、`NetworkSummary` | yes |
| 子网/IPAM 配置 | `SubNetworkConfig` | yes |
| Docker 版本探测 | `Version` 结构 + `Version()` | yes |
| 业务持久化 | 不入库 | no |
| 用户鉴权与权限 | 由上层校验 | no |
| HTTP 入口 | 纯 Go 库，不暴露路由 | no |
| 跨模块协调 | 仅作为 Go library 被引用 | conditional |
| 其他模块内部状态 | 不直接读写其他模块私有实现 | no |

## 模块事实

- 模块路径：`github.com/amuluze/docker`（独立 Go module，通过 monorepo 路径 `common/docker/` 维护）
- Go 版本：`1.25.0`
- Go 源码文件：7
- Go 测试文件：6
- 包名：`docker`
- 推断角色：Docker Engine SDK 薄封装，为 collia 与 beacon 提供统一的容器/镜像/网络操作抽象

## 接口与符号信号

本节记录源码扫描得到的导出符号或前端入口信号；它们是候选接口线索，不等同于跨模块公开 API。跨模块调用仍以 facade、模块文档和 Domain Spec 约束为准。

- `interface IManager (common/docker/manager.go)`
- `struct SubNetworkConfig (common/docker/network.go)`
- `func NewManager (common/docker/manager.go)`
- `func IsNameConflict (common/docker/conflict.go)`
- `struct ContainerSummary (common/docker/container.go)`
- `struct ImageSummary (common/docker/image.go)`
- `struct Manager (common/docker/manager.go)`
- `struct NetworkSummary (common/docker/network.go)`
- `struct PortMapping (common/docker/container.go)`
- `struct Version (common/docker/version.go)`

## 依赖关系

- `github.com/docker/docker` `v27.2.1+incompatible`
- `github.com/docker/go-connections` `v0.5.0`
- `github.com/docker/libcompose` `v0.4.0`
- `github.com/tidwall/gjson` `v1.17.1`
- `gopkg.in/yaml.v3` `v3.0.1`
- 间接依赖 27 个，完整列表以模块 `go.mod` 为准。

## 配置与入口信号

- 配置：`IPAM` at `common/docker/network.go:49`
- 配置：`IPAM` at `common/docker/network.go:91`
- 配置：`IPAM` at `common/docker/network.go:133`
- 配置：`SubNetworkConfig` in `common/docker/network.go`

## 数据模型与状态

- 数据模型以源码结构、导出类型信号和测试断言为准；本节只记录静态发现结果。
- 容器 CPU/内存百分比仅在分子、分母均为正且计算结果有限时返回计算值；Docker 对未启动容器返回的零值 stats 会归一为 `0`，禁止向上层泄漏 `NaN`/无穷值。
- `ContainerSummary` 携带创建时配置的启动命令（`Commands`，取自 Docker `Config.Cmd`）与环境变量：`ListContainer` 对每个容器只做一次 `ContainerInspect`，同时用于健康态判定、环境变量与启动命令，避免同一容器重复请求 Engine API。
- `interface IManager (common/docker/manager.go)`
- `struct ContainerSummary (common/docker/container.go)`
- `struct ImageSummary (common/docker/image.go)`
- `struct Manager (common/docker/manager.go)`
- `struct NetworkSummary (common/docker/network.go)`
- `struct PortMapping (common/docker/container.go)`
- `struct SubNetworkConfig (common/docker/network.go)`
- `struct Version (common/docker/version.go)`

## 使用注意事项

- 不跨越模块边界直接依赖内部路径或未导出符号。
- 修改跨模块公开 API、导出符号信号、配置或副作用边界时，同步更新本文件、根入口文档和相关 Domain Spec。
- 配置、凭据和外部副作用只记录语义边界，不记录真实敏感值。

## 验证命令

```bash
# common/docker 子模块无独立 Taskfile.yml，验证通过 Go 标准工具链完成
cd common/docker

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
