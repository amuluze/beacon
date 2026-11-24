# Beacon Web

`beacon/web` 模块协作入口，汇总当前 workspace 的模块事实。

该模块当前角色：frontend experience module: Vue/Vite 页面、路由、API client、Pinia 状态管理和用户交互。

## 文档

模块协作入口以本 AGENTS.md 为准；本仓库不携带项目级文档系统，实现细节以源码与根 README 为准。


## 开发导航

- 先读本文件确认模块边界，再读对应 `.docs/modules/` 文档获取当前实现事实；导出符号只作为入口线索，不自动等同跨模块公开 API。
- 代码变更后按本文件“开发命令”执行最小验证。
- 更新公开 API、配置键、事件、持久化格式或用户可见工作流时，同步维护相关项目文档并检查 `AGENTS.md` / `CLAUDE.md` 一致性。
- 前端相关修改优先检查 API client、store 状态、组件交互和构建/测试脚本。

## 模块路径

`beacon/web/package.json`

## 关键目录

| 目录/文件 | 职责 |
|-----------|------|
| `beacon/web/` | frontend experience module: Vue/Vite 页面、路由、API client、Pinia 状态管理和用户交互 |
| `beacon/web/public/` | static assets or embedded resources |
| `beacon/web/src/` | application source code |
| `beacon/web/types/` | unplugin-auto-import 与 unplugin-vue-components 生成的类型声明（`auto-imports.d.ts`、`components.d.ts` 等） |

## 依赖

- `@antfu/eslint-config`
- `@element-plus/icons-vue`
- `@iconify-json/ep`
- `@iconify-json/lucide`
- `@types/codemirror`
- `@types/lodash-es`
- `@types/node`
- `@unocss/eslint-plugin`
- `@vitejs/plugin-vue`
- `@vitest/coverage-v8`
- `@vue/test-utils`
- `@xterm/addon-fit`
- `@xterm/addon-web-links`
- `@xterm/xterm`
- `axios`
- `codemirror`
- `codemirror-editor-vue3`
- `echarts`
- `element-plus`
- `eslint`
- `eslint-plugin-format`
- `happy-dom`
- `lint-staged`
- `lodash-es`
- `path-browserify`
- `pinia`
- `pinia-plugin-persistedstate`
- `postcss`
- `sass`
- `simple-git-hooks`
- `terser`
- `typescript`
- `unocss`
- `unplugin-auto-import`
- `unplugin-icons`
- `unplugin-vue-components`
- `vite`
- `vite-plugin-svg-icons`
- `vitest`
- `vue`
- `vue-i18n`
- `vue-router`
- `vue-tsc`

## 模块约束

- 仅通过公开接口与其他模块协作，不依赖其他模块内部实现细节。
- 修改公开 API、配置或副作用边界时，同步更新 `.docs/modules/` 中对应文档。

## 开发命令

```bash
cd beacon/web && pnpm run build
cd beacon/web && pnpm run dev
cd beacon/web && pnpm run lint
cd beacon/web && pnpm run preview
cd beacon/web && pnpm run test
cd beacon/web && pnpm run test:coverage
cd beacon/web && pnpm run test:run
cd beacon/web && pnpm run ts
```
