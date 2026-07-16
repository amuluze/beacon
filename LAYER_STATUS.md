# 分层执行状态

## L1 - 基础设施

- 状态：✅ NOT REQUIRED
- 原因：本次不新增数据库表或迁移；复用现有监控、审计和容器控制模型。

## L2 - 后端

- 状态：✅ DONE
- 前置条件：L1 无需执行
- 已完成：2.1 容器策略契约、2.2 Agent 容器编辑、2.3 审计分页一致性
- 基线验证：

  ```text
  (cd common && go test ./...)
  ok common/database
  ok common/rpc/schema
  ok common/rpc/tunnel

  (cd beacon && go test ./service/container/... ./service/audit/...)
  ok beacon/service/container/repository
  ok beacon/service/container/service
  ok beacon/service/audit/repository
  ok beacon/service/audit/service

  (cd collia && go test ./service/rpc/...)
  ok collia/service/rpc
  ```
- 2026-07-16 官网修复追加：
  - 统计接口信封契约、非法响应和版本比较错误路径已补测试。
  - Fiber 仅信任环回 Nitro 代理提供的转发 IP，并限制超时与请求体大小。
  - 写限流通过真实 loopback listener 验证不同转发客户端 IP 的额度彼此独立。
  - Fiber、pgx、ch-go/clickhouse-go 升级到兼容的安全版本，官网后端 `govulncheck ./...` 无可达漏洞。
  - 官网后端 `GOWORK=off go test -race ./...`、`go vet ./...`、`go build ./...` 全部通过。

## L3 - 前端

- 状态：✅ DONE
- 前置条件：相关 L2 契约完成；工作区聚合与官网 Shell 可并行准备
- 已完成：
  - Monitor / Container / Settings 一级工作区与旧深链兼容。
  - 按 Agent 查询的 Settings 审计表格与请求头注入。
  - 容器创建/编辑、重启策略与运行配置安全继承。
  - 480px Modal/Drawer、540px Registry Drawer 与失败反馈修复。
  - 官网 Landing / About / Changelog / Docs、设计令牌与 375px 响应式布局。
- 2026-07-16 官网修复追加：
  - 自研 Toast、统计静默降级、移动 Menu、SEO/隐私页、可访问性和安全头已完成。
  - `pnpm test`（5 文件、13 测试）、typecheck、ESLint、Nuxt build 全部通过。

## L4 - 集成验证

- 状态：✅ DONE
- 前置条件：L2、L3 完成
- 已完成：
  - 四个 Go module 全量 test/vet/build；官网后端额外通过 race 检查。
  - 后台既有 53 个 Vitest、类型检查与 Vite build；官网 13 个 Vitest、类型检查、ESLint 与 Nuxt build。
  - `task website:verify`、SHA-256、Compose 配置、`pnpm audit`、官网后端 `govulncheck` 与 `git diff --check`。
  - 375px 首页菜单与统计失败降级、375px 文档页溢出/语义/SEO、1440px 桌面布局均已在最新构建上实测。
  - 两种视口的控制台均无 hydration、图标或其他 warning/error。
