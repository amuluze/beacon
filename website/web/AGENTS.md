# Website Web

`website/web` 模块入口文档，由 `/doc update` 根据当前 workspace 事实重写。

该模块当前角色：frontend experience module: Vue/Vite 页面、路由、API client、状态管理和用户交互。

## 文档

模块协作入口以本 AGENTS.md 为准；项目级文档见根 AGENTS 文档地图。

- 模块实现事实见 [`.docs/modules/website-web.md`](../../.docs/modules/website-web.md)。

## 开发导航

- 先读本文件确认模块边界，再读对应 `.docs/modules/` 文档获取当前实现事实；导出符号只作为入口线索，不自动等同跨模块公开 API。
- 涉及长期行为、不变量、状态或错误语义时，必须回到下方相关 Domain Spec；若现有 Domain Spec 不覆盖，应先补可验证约束。
- 代码变更后按本文件“开发命令”执行最小验证；跨模块、配置、接口或副作用变更还要运行项目级质量门禁或 `task sdd:refs`。
- 更新公开 API、配置键、事件、持久化格式或用户可见工作流时，重新运行 `/doc update` 并检查 `AGENTS.md` / `CLAUDE.md` 同步。
- 前端相关修改优先检查 API client、store 状态、组件交互和构建/测试脚本。

## 模块路径

`website/web/package.json`

## 关键目录

| 目录/文件 | 职责 |
|-----------|------|
| `website/web/` | frontend experience module: Vue/Vite 页面、路由、API client、状态管理和用户交互 |
| `website/web/public/` | static assets or embedded resources |
| `website/web/src/` | application source code |
| `website/server/cmd/server/` | Go package `main`，源码 1，测试 0 |
| `website/server/pkg/database/` | Go package `database`，源码 3，测试 1 |
| `website/server/pkg/errors/` | Go package `errors`，源码 1，测试 1 |
| `website/server/pkg/fiberx/` | Go package `fiberx`，源码 1，测试 1 |
| `website/server/pkg/utils/` | Go package `utils`，源码 1，测试 1 |
| `website/server/pkg/utils/hash/` | Go package `hash`，源码 1，测试 0 |

## 依赖

- `@antfu/eslint-config`
- `@iconify-json/mdi`
- `@nuxt/icon`
- `@nuxt/test-utils`
- `@vue/test-utils`
- `eslint`
- `eslint-plugin-format`
- `happy-dom`
- `nuxt`
- `sass`
- `typescript`
- `vitest`
- `vue`
- `vue-router`
- `vue-tsc`

## 模块约束

- 仅通过公开接口与其他模块协作，不依赖其他模块内部实现细节。
- 修改公开 API、配置或副作用边界时，同步更新 `.docs/modules/` 中对应文档。
- 若模块承载长期领域语义，相关约束应在 `.specs/domain/` 中可追踪。

## 开发命令

```bash
cd website/web && npm run lint
cd website/web && npm run test
cd website/web && npm run build
cd website/web && npm run dev
```
