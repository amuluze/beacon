# Beacon

`beacon` 模块入口文档，由 `/doc update` 根据当前 workspace 事实重写。

该模块当前角色：HTTP/API service module: Web/API 接入、路由注册、请求校验和服务协调。

## 文档

模块协作入口以本 AGENTS.md 为准；项目级文档见根 AGENTS 文档地图。

- 模块实现事实见 [`.docs/modules/beacon.md`](../.docs/modules/beacon.md)。

## 开发导航

- 先读本文件确认模块边界，再读对应 `.docs/modules/` 文档获取当前实现事实；导出符号只作为入口线索，不自动等同跨模块公开 API。
- 涉及长期行为、不变量、状态或错误语义时，必须回到下方相关 Domain Spec；若现有 Domain Spec 不覆盖，应先补可验证约束。
- 代码变更后按本文件“开发命令”执行最小验证；跨模块、配置、接口或副作用变更还要运行项目级质量门禁或 `task sdd:refs`。
- 更新公开 API、配置键、事件、持久化格式或用户可见工作流时，重新运行 `/doc update` 并检查 `AGENTS.md` / `CLAUDE.md` 同步。

## 模块路径

`beacon`

## 关键目录

| 目录/文件 | 职责 |
|-----------|------|
| `beacon/` | HTTP/API service module: Web/API 接入、路由注册、请求校验和服务协调 |
| `beacon/cmd/` | 命令行或进程入口 |
| `beacon/configs/` | 运行时配置文件 |
| `beacon/nginx/` | Nginx 部署或反向代理配置 |
| `beacon/pkg/` | 可复用公共包集合 |
| `beacon/pkg/auth/` | 认证授权与访问控制 |
| `beacon/pkg/contextx/` | 请求上下文辅助封装 |
| `beacon/pkg/errors/` | 错误类型与错误响应封装 |
| `beacon/pkg/fiberx/` | HTTP/Fiber 响应或中间件封装 |
| `beacon/pkg/rpc/` | RPC client/server 封装 |
| `beacon/pkg/utils/` | 通用辅助函数 |
| `beacon/pkg/validatex/` | 输入校验封装 |
| `beacon/service/` | 业务核心层：路由、认证、Server/RPC、数据库或领域服务 |
| `beacon/supervisor/` | Supervisor 进程管理配置 |
| `beacon/web/` | 前端应用或静态资源目录 |
| `beacon/cmd/beacon/` | Go package `main`，源码 1，测试 0 |
| `beacon/pkg/auth/jwtauth/` | Go package `jwtauth`，源码 4，测试 3 |

## 依赖

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

## 模块约束

- 仅通过公开接口与其他模块协作，不依赖其他模块内部实现细节。
- 修改公开 API、配置或副作用边界时，同步更新 `.docs/modules/` 中对应文档。
- 若模块承载长期领域语义，相关约束应在 `.specs/domain/` 中可追踪。

## 开发命令

```bash
# 未检测到该模块来自 CI、Taskfile 或 Makefile 的开发/验证命令。
```
