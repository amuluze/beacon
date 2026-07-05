# Amprobe

`amprobe` 模块入口文档，由 `/sdd doc update` 根据当前 workspace 事实重写。

该模块当前角色：Server control plane: Web/API 接入、认证授权、Agent 生命周期、监控批次落库、目标选择和反向 tunnel client。

## 文档

模块协作入口以本 AGENTS.md 为准；项目级文档见根 AGENTS 文档地图。

- 模块实现事实见 [`.docs/modules/amprobe.md`](../.docs/modules/amprobe.md)。

## 开发导航

- 先读本文件确认模块边界，再读对应 `.docs/modules/` 文档获取当前实现事实；导出符号只作为入口线索，不自动等同跨模块公开 API。
- 涉及长期行为、不变量、状态或错误语义时，必须回到下方相关 Domain Spec；若现有 Domain Spec 不覆盖，应先补可验证约束。
- 代码变更后按本文件“开发命令”执行最小验证；跨模块、配置、接口或副作用变更还要运行项目级质量门禁或 `task sdd:refs`。
- 更新公开 API、配置键、事件、持久化格式或用户可见工作流时，重新运行 `/sdd doc update` 并检查 `AGENTS.md` / `CLAUDE.md` 同步。
- CLI/TUI 相关修改优先检查命令入口、交互状态、工具授权、会话协调和用户可见输出。

## 模块路径

`amprobe`

## 关键目录

| 目录/文件 | 职责 |
|-----------|------|
| `amprobe/` | Server control plane: Web/API 接入、认证授权、Agent 生命周期、监控批次落库、目标选择和反向 tunnel client |
| `amprobe/cmd/` | 命令行或进程入口 |
| `amprobe/configs/` | 运行时配置文件 |
| `amprobe/nginx/` | Nginx 部署或反向代理配置 |
| `amprobe/pkg/` | 可复用公共包集合 |
| `amprobe/pkg/auth/` | 认证授权与访问控制 |
| `amprobe/pkg/contextx/` | 请求上下文辅助封装 |
| `amprobe/pkg/errors/` | 错误类型与错误响应封装 |
| `amprobe/pkg/fiberx/` | HTTP/Fiber 响应或中间件封装 |
| `amprobe/pkg/psutil/` | supporting project directory |
| `amprobe/pkg/rpc/` | RPC client/server 封装 |
| `amprobe/pkg/utils/` | 通用辅助函数 |
| `amprobe/pkg/validatex/` | 输入校验封装 |
| `amprobe/service/` | 业务核心层：路由、认证、Server/RPC、数据库或领域服务 |
| `amprobe/supervisor/` | Supervisor 进程管理配置 |
| `amprobe/web/` | 前端应用或静态资源目录 |
| `amprobe/cmd/amprobe/` | Go package `main`，源码 1，测试 0 |
| `amprobe/pkg/auth/jwtauth/` | Go package `jwtauth`，源码 4，测试 0 |

## 依赖

- `github.com/amuluze/amutool/logger` `v0.0.0-20240821104128-caed9cc0d402`
- `github.com/amuluze/amutool/timex` `v0.0.0-20250508153823-fe9a5de55958`
- `github.com/casbin/casbin/v2` `v2.98.0`
- `github.com/casbin/gorm-adapter/v3` `v3.28.0`
- `github.com/go-playground/validator/v10` `v10.22.0`
- `github.com/gofiber/contrib/websocket` `v1.3.2`
- `github.com/gofiber/fiber/v2` `v2.52.5`
- `github.com/golang-jwt/jwt` `v3.2.2+incompatible`
- `github.com/google/uuid` `v1.6.0`
- `github.com/google/wire` `v0.6.0`
- `github.com/patrickmn/go-cache` `v2.1.0+incompatible`
- `github.com/pkg/errors` `v0.9.1`
- `github.com/spf13/viper` `v1.19.0`
- `github.com/urfave/cli/v2` `v2.27.4`
- `gopkg.in/gomail.v2` `v2.0.0-20160411212932-81ebce5c23df`
- `gorm.io/gorm` `v1.25.12`

## 模块约束

- 仅通过公开接口与其他模块协作，不依赖其他模块内部实现细节。
- 修改公开 API、配置或副作用边界时，同步更新 `.docs/modules/` 中对应文档。
- 若模块承载长期领域语义，相关约束应在 `.specs/domain/` 中可追踪。
- 监控查询读取 Server 本地监控表；Agent 上报通过 HTTP report 入口写入；控制操作通过反向 tunnel 调用 Agent。
- 新增控制调用时必须明确 Agent 选择来源，并禁止未实现操作返回成功空结果。

## 开发命令

```bash
cd amprobe && go test ./...
cd amprobe && go build ./...
```
