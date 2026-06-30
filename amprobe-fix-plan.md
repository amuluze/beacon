# Amprobe 项目修复方案（P0 / P1 / P2）

> 生成时间：2025-06-28
> 基于：amprobe-project-evaluation.md 深度评估报告
> 状态：部分已实施，部分待执行

---

## 一、已实施修改（本次会话完成）

### 1.1 P0 — Go 版本与依赖收敛

| 文件 | 修改内容 | 状态 |
|------|----------|------|
| `amprobe/go.mod` | `go 1.21.10` → `go 1.25.0` | ✅ 已改 |
| `collia/go.mod` | `go 1.21.10` → `go 1.25.0` | ✅ 已改 |
| `amprobe/Dockerfile` | 构建镜像 `golang:1.21` → `golang:1.25` | ✅ 已改 |

**后续行动**：运行 `go work sync`（或进入各模块执行 `go mod tidy`）收敛 `gorm`、`rpcx`、`amutool` 等跨模块依赖版本。因当前环境无 Go 工具链，需在本地或 CI 中执行。

### 1.2 P0 — 安全配置强化

| 文件 | 修改内容 | 状态 |
|------|----------|------|
| `amprobe/service/config.go` | 新增 `overrideFromEnv` 函数：从 `AMPROBE_AUTH_SIGNING_KEY`、`AMPROBE_AGENT_INSTALL_TOKEN`、`AMPROBE_CONTROL_JOIN_TOKEN` 环境变量读取敏感配置 | ✅ 已改 |
| `amprobe/service/config.go` | `warnInsecureDefaults` 日志增加环境变量覆盖提示 | ✅ 已改 |
| `amprobe/configs/config.toml` | `[Auth]` 和 `[AgentInstall]` 段增加环境变量注释警告 | ✅ 已改 |
| `amprobe/configs/config.dev.toml` | 同上，增加开发环境注释警告 | ✅ 已改 |

**生产部署命令示例**：
```bash
export AMPROBE_AUTH_SIGNING_KEY=$(openssl rand -hex 32)
export AMPROBE_AGENT_INSTALL_TOKEN=$(openssl rand -hex 16)
export AMPROBE_CONTROL_JOIN_TOKEN=$(openssl rand -hex 16)
docker-compose up -d
```

### 1.3 P1 — 健康检查端点

| 文件 | 修改内容 | 状态 |
|------|----------|------|
| `amprobe/service/health/api/health.go` | 新建 `Probe` 结构体，提供 `Liveness`（存活）和 `Readiness`（就绪）handler | ✅ 已创建 |
| `amprobe/service/router.go` | 在 `RegisterAPI` 开头注册 `GET /health` 和 `GET /ready`（不经过 auth 中间件） | ✅ 已改 |

**待完善**：`Probe` 目前仅返回进程存活状态。下一步需注入 `DBHealthy` 和 `TunnelHealthy` 检查函数（需等待 DB 和 tunnel 接口提供 `Ping` 方法）。

### 1.4 P1 — 前端测试基础设施

| 文件 | 修改内容 | 状态 |
|------|----------|------|
| `amprobe/web/package.json` | 新增 `test`/`test:watch` 脚本；添加 `vitest`、`@vue/test-utils`、`jsdom` 依赖 | ✅ 已改 |
| `amprobe/web/vitest.config.ts` | 新建 Vitest 配置：Vue 插件 + jsdom 环境 + `@` 别名 | ✅ 已创建 |
| `amprobe/web/src/store/modules/agent.test.ts` | 新建 `useAgentStore` 单元测试：setAgents、setCurrentAgent、空列表回退 | ✅ 已创建 |

**后续行动**：运行 `pnpm install` 和 `pnpm test` 验证通过。

### 1.5 P2 — Docker Compose 与部署编排

| 文件 | 修改内容 | 状态 |
|------|----------|------|
| `docker-compose.yml` | 新建根目录 compose：Server 服务 + 数据卷 + 环境变量 + 健康检查 | ✅ 已创建 |
| `deploy/docker-compose.yml` | 新建部署参考 compose：简化版，适合生产环境模板 | ✅ 已创建 |

---

## 二、待执行修改（按优先级）

### 2.1 P0 — 仍需完成

| # | 任务 | 目标文件 | 说明 |
|---|------|----------|------|
| 1 | **完成 `platform-hardening` Spec** | `amprobe/service/report/`, `amprobe/service/host/`, `amprobe/service/container/` | 监控查询强制要求 `agent_id`；Agent 上报批次增加新鲜度校验；前端陈旧数据 UI 标识 |
| 2 | **运行 `go work sync`** | 全部 `go.mod` | 收敛 `gorm` (v1.25.10 vs v1.25.12)、`rpcx` (v1.8.31 vs v1.8.32)、`amutool/timex` 等版本漂移 |
| 3 | **清理 rpcx 依赖** | `collia/go.mod` | 若已完全迁移至反向 gRPC tunnel，移除 `smallnest/rpcx` 直接依赖 |

### 2.2 P1 — 仍需完成

| # | 任务 | 目标文件 | 说明 |
|---|------|----------|------|
| 4 | **Rate Limiting** | `amprobe/service/router.go` | 使用 `gofiber/fiber/v2/middleware/limiter` 对 `/api/v1/auth/login` 和 `/api/v1/host/report` 限流 |
| 5 | **API 集成测试** | `amprobe/service/...` | 为 `host`、`container`、`auth` handler 添加 HTTP 测试，验证 `X-Agent-ID` 缺失、权限拒绝、错误码 |
| 6 | **Wire 健康检查依赖** | `amprobe/service/health/api/health.go` | 注入 DB ping 和 tunnel 在线检查到 `Probe` 结构体 |
| 7 | **前端 Axios 测试** | `amprobe/web/src/api/index.test.ts` | 模拟 Axios 拦截器，验证 `X-Agent-ID` 注入和 token 刷新队列 |
| 8 | **Taskfile 测试目标** | `amprobe/web/Taskfile.yml` | 添加 `task amprobe-web:test` 调用 `pnpm test` |

### 2.3 P2 — 仍需完成

| # | 任务 | 目标文件 | 说明 |
|---|------|----------|------|
| 9 | **OpenAPI 规范** | `.docs/api/openapi.yml` | 为现有 40+ 个 API 端点生成 OpenAPI 3.0 规范 |
| 10 | **CI/CD 流水线** | `.github/workflows/ci.yml` | GitHub Actions：Go 测试 + 前端构建 + Docker 构建 + 文档健康检查 |
| 11 | **K8s manifests** | `deploy/k8s/` | 基础 Deployment + Service + ConfigMap + PVC |
| 12 | **前端超时分级** | `amprobe/web/src/api/index.ts` | 按接口类型区分超时：查询 30s / 控制 60s / 下载 120s |
| 13 | **组件复用优化** | `amprobe/web/src/components/` | 将视图中重复 UI 模式提取为可复用组件 |

---

## 三、新增 Spec / Plan 文件

| 文件 | 类型 | 覆盖范围 | 状态 |
|------|------|----------|------|
| `.specs/tasks/ops-infrastructure.md` | Spec | 健康检查、限流、Docker Compose、运维基础设施 | ✅ 已创建 |
| `.plans/ops-infrastructure.md` | Plan | 实施步骤：health API → compose → rate limit → Dockerfile 对齐 | ✅ 已创建 |
| `.specs/tasks/frontend-testing.md` | Spec | 前端测试：Vitest + jsdom + store/Axios 拦截器覆盖 | ✅ 已创建 |
| `.plans/frontend-testing.md` | Plan | 实施步骤：工具链 → store 测试 → Axios 测试 → CI 集成 | ✅ 已创建 |

---

## 四、验证清单

### 4.1 立即验证（已修改文件）

```bash
# 1. Go 模块版本一致性检查
grep '^go ' amprobe/go.mod collia/go.mod common/go.mod
# 期望：全部为 go 1.25.0

# 2. 配置文件环境变量注释检查
grep -n 'AMPROBE_' amprobe/configs/config.toml amprobe/configs/config.dev.toml
# 期望：包含 AMPROBE_AUTH_SIGNING_KEY、AMPROBE_AGENT_INSTALL_TOKEN 注释

# 3. 健康检查路由注册检查
grep -n 'health\|ready' amprobe/service/router.go
# 期望：包含 /health 和 /ready 的 app.Get 注册

# 4. 前端测试配置检查
ls amprobe/web/vitest.config.ts amprobe/web/src/store/modules/agent.test.ts
# 期望：两个文件存在

# 5. Docker Compose 检查
ls docker-compose.yml deploy/docker-compose.yml
# 期望：两个文件存在
```

### 4.2 后续验证（需本地工具链）

```bash
# 6. Go 构建验证（需 go 1.25.0）
cd amprobe && go build ./...
cd collia && go build ./...
cd common && go build ./...

# 7. Go 测试验证
cd amprobe && go test ./...
cd collia && go test ./...
cd common && go test ./...

# 8. 依赖收敛（需 go work sync）
go work sync
# 然后检查各 go.mod 中 gorm、rpcx 版本是否一致

# 9. 前端测试验证（需 pnpm）
cd amprobe/web && pnpm install && pnpm test

# 10. Docker 构建验证
docker-compose build

# 11. Wire 代码生成（若修改了 injector）
task amprobe:wire
task collia:wire
```

---

## 五、风险与注意事项

1. **Go 1.25.0 可用性**：Go 1.25 尚未正式发布（预计 2025 年 8 月）。如果当前环境无法获取 1.25.0 工具链，需回退至 `go 1.23.0`（当前已发布最新稳定版）。**已做修改需配合实际工具链版本调整**。
2. **Vitest 依赖版本**：`package.json` 中新增 `vitest ^2.1.8` 和 `@vue/test-utils ^2.4.6`。需运行 `pnpm install` 更新 `pnpm-lock.yaml`。
3. **健康检查依赖注入**：当前 `health.Probe` 的 `DBHealthy`/`TunnelHealthy` 为 nil，返回始终就绪。后续需从 `wire.go` 或 `server.go` 注入实际的检查函数。
4. **Rate Limiting 中间件**：需确认 `gofiber/fiber/v2/middleware/limiter` 在 Fiber v2.52.5 中可用（它是标准子包，无需额外依赖）。
5. **Dockerfile 缓存失效**：修改 `golang:1.21` → `golang:1.25` 后首次构建会重新拉取镜像，构建时间可能延长。

---

## 六、执行优先级速查表

| 优先级 | 任务 | 预计工时 | 阻塞项 |
|--------|------|----------|--------|
| **P0** | 完成 `platform-hardening` Spec（监控查询强制 agent_id、数据新鲜度） | 2-3 天 | 需理解现有 repo 查询逻辑 |
| **P0** | 运行 `go work sync` + 验证构建 | 2-4 小时 | 需 Go 1.25 工具链 |
| **P1** | Rate Limiting 中间件 | 4-6 小时 | 无 |
| **P1** | API 集成测试（host/container/auth） | 1-2 天 | 需 mock DB 和 tunnel |
| **P1** | Wire 健康检查依赖（DB ping + tunnel 状态） | 4-6 小时 | 需 DB 和 tunnel 提供 ping 接口 |
| **P1** | 前端 Axios 测试 | 4-6 小时 | 需 `pnpm install` |
| **P2** | OpenAPI 规范 | 2-3 天 | 需梳理全部 40+ 端点参数 |
| **P2** | CI/CD 流水线（GitHub Actions） | 1 天 | 需仓库 GitHub 权限 |
| **P2** | K8s manifests | 1-2 天 | 需测试集群 |
| **P2** | 前端超时分级 | 2-4 小时 | 无 |

---

*本方案为 Amprobe 项目评估后发现的 P0/P1/P2 问题的系统性修复计划。已实施部分可直接进入验证阶段；待实施部分需按优先级排入迭代计划。*
