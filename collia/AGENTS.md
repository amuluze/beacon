# Collia

`collia` 模块入口文档，由 `/doc update` 根据当前 workspace 事实重写。

该模块当前角色：persistence-aware service module: 数据库模型、仓储和事务边界。

## 文档

模块协作入口以本 AGENTS.md 为准；项目级文档见根 AGENTS 文档地图。

- 模块实现事实见 [`.docs/modules/collia.md`](../.docs/modules/collia.md)。

## 开发导航

- 先读本文件确认模块边界，再读对应 `.docs/modules/` 文档获取当前实现事实；导出符号只作为入口线索，不自动等同跨模块公开 API。
- 涉及长期行为、不变量、状态或错误语义时，必须回到下方相关 Domain Spec；若现有 Domain Spec 不覆盖，应先补可验证约束。
- 代码变更后按本文件“开发命令”执行最小验证；跨模块、配置、接口或副作用变更还要运行项目级质量门禁或 `task sdd:refs`。
- 更新公开 API、配置键、事件、持久化格式或用户可见工作流时，重新运行 `/doc update` 并检查 `AGENTS.md` / `CLAUDE.md` 同步。

## 模块路径

`collia`

## 关键目录

| 目录/文件 | 职责 |
|-----------|------|
| `collia/` | persistence-aware service module: 数据库模型、仓储和事务边界 |
| `collia/cmd/` | 命令行或进程入口 |
| `collia/pkg/` | 可复用公共包集合 |
| `collia/pkg/conn/` | supporting project directory |
| `collia/pkg/psutil/` | supporting project directory |
| `collia/pkg/resources/` | supporting project directory |
| `collia/pkg/timectl/` | supporting project directory |
| `collia/pkg/utils/` | 通用辅助函数 |
| `collia/script/` | supporting project directory |
| `collia/service/` | 业务核心层：路由、认证、Server/RPC、数据库或领域服务 |
| `collia/cmd/collia/` | Go package `main`，源码 2，测试 0 |

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
# 未检测到该模块来自 CI、Taskfile 或 Makefile 的开发/验证命令。
```
