# Common

`common` 模块协作入口，汇总当前 workspace 的模块事实。

该模块当前角色：shared contract library: 复用 schema、数据库封装、RPC 参数/返回值和跨模块类型。

## 文档

模块协作入口以本 AGENTS.md 为准；项目级文档见根 AGENTS 文档地图。

- 模块实现事实见 [`.docs/modules/common.md`](../.docs/modules/common.md)。

## 开发导航

- 先读本文件确认模块边界，再读对应 `.docs/modules/` 文档获取当前实现事实；导出符号只作为入口线索，不自动等同跨模块公开 API。
- 涉及长期行为、不变量、状态或错误语义时，必须回到下方相关 Domain Spec；若现有 Domain Spec 不覆盖，应先补可验证约束。
- 代码变更后按本文件“开发命令”执行最小验证；跨模块、配置、接口或副作用变更还要运行项目级质量门禁或 `task sdd:refs`。
- 更新公开 API、配置键、事件、持久化格式或用户可见工作流时，同步维护相关项目文档并检查 `AGENTS.md` / `CLAUDE.md` 一致性。

## 模块路径

`common`

## 关键目录

| 目录/文件 | 职责 |
|-----------|------|
| `common/` | shared contract library: 复用 schema、数据库封装、RPC 参数/返回值和跨模块类型 |
| `common/database/` | supporting project directory |
| `common/docker/` | supporting project directory |
| `common/logger/` | supporting project directory |
| `common/rpc/` | RPC client/server 封装 |
| `common/timex/` | supporting project directory |
| `common/transport/` | supporting project directory |
| `common/rpc/schema/` | Go package `schema`，源码 8，测试 1 |
| `common/rpc/tunnel/` | Go package `tunnel`，源码 4，测试 1 |

## 依赖

- `github.com/glebarez/sqlite` `v1.11.0`
- `google.golang.org/grpc` `v1.81.1`
- `google.golang.org/protobuf` `v1.36.11`
- `gorm.io/driver/clickhouse` `v0.6.1`
- `gorm.io/driver/mysql` `v1.5.7`
- `gorm.io/driver/postgres` `v1.5.9`
- `gorm.io/gorm` `v1.25.12`

## 模块约束

- 仅通过公开接口与其他模块协作，不依赖其他模块内部实现细节。
- 修改公开 API、配置或副作用边界时，同步更新 `.docs/modules/` 中对应文档。
- 若模块承载长期领域语义，相关约束应在 `.specs/domain/` 中可追踪。

## 开发命令

```bash
# 未检测到该模块来自 CI、Taskfile 或 Makefile 的开发/验证命令。
```
