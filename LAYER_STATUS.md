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

## L3 - 前端

- 状态：✅ DONE
- 前置条件：相关 L2 契约完成；工作区聚合与官网 Shell 可并行准备
- 已完成：
  - Monitor / Container / Settings 一级工作区与旧深链兼容。
  - 按 Agent 查询的 Settings 审计表格与请求头注入。
  - 容器创建/编辑、重启策略与运行配置安全继承。
  - 480px Modal/Drawer、540px Registry Drawer 与失败反馈修复。
  - 官网 Landing / About / Changelog / Docs、设计令牌与 375px 响应式布局。

## L4 - 集成验证

- 状态：⚠️ PARTIAL（工程验证完成，视觉截图阻塞）
- 前置条件：L2、L3 完成
- 已完成：三 Go module 全量测试与构建、后台 53 个 Vitest、类型检查与 Vite build、`task website-web:build`、源码反模式扫描与 `git diff --check`。
- 阻塞：内置浏览器返回 `No browser is available`，无法完成 1440px/375px 截图对照及统计 500 的浏览器拦截验收。
