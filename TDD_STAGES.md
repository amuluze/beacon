# TDD 阶段记录

## RED 阶段

- 状态：✅ DONE
- 目标：为容器 restart policy 透传、Agent 容器编辑回滚、审计 filtered total 与前端 Agent scoped 请求建立失败测试。
- 计划测试：
  - `TestContainerService_ContainerCreatePassesRestartPolicy`
  - `TestContainerService_ContainerUpdatePassesRestartPolicy`
  - `TestContainerUpdateRecreatesAndRestoresRunningState`
  - `TestContainerUpdateRollsBackWhenCreateFails`
  - `TestContainerUpdateRollsBackWhenRestartPolicyFails`
  - `TestAuditCountUsesSameFiltersAsQuery`
  - `api/index` audit prefix Agent header test
  - 后台一级工作区与旧深链路由测试
  - 容器 create/update payload 序列化测试
  - 邮箱取消不污染 props、阈值失败不提示成功测试
  - Agent 列表并发请求去重测试
  - 容器更新未提供运行配置时继承旧值测试
- 运行结果：✅ RED（预期失败）
  - Beacon：`RestartPolicy` 字段不存在；`AuditCount` 尚未接收筛选参数。
  - Collia：`recreateContainer` 尚不存在；RPC schema 尚无 `RestartPolicy`。

## GREEN 阶段

- 状态：✅ DONE
- 运行结果：
  - Common：`go test ./... && go build ./...` 通过。
  - Beacon：`go test ./service/container/... ./service/audit/... && go build ./...` 通过。
  - Collia：`go test ./service/rpc/... && go build ./...` 通过（仅保留已有 `go-m1cpu` clang warning）。

## REFACTOR 阶段

- 状态：✅ DONE
- 完成内容：
  - `fakes.go` 已收敛为仅保留审计接口签名修改。
  - 容器重建对 nil 运行配置执行安全继承，显式空数组仍保留清空语义。
  - 聚合工作区共享 Agent 列表请求，避免首次加载请求风暴。
  - 官网移除无引用 Carousel、旧发光/渐变组件和失效 SvgIcon 插件。
  - 后台与官网编译、相关/全量测试全部通过；视觉截图因无浏览器单独记录为阻塞。
