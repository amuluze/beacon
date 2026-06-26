# Task: Agent 远程卸载与本地卸载完善

> 依赖: Domain Spec [.specs/domain/agent-lifecycle-update.md](../domain/agent-lifecycle-update.md)
> 状态: pending
> 优先级: P2 — 运维便利性

## 设计意图

让 Server 能远程触发指定 Agent 的自卸载，同时完善本地卸载命令，确保卸载后不残留文件和服务注册。

## 可验证行为约束

### 必须满足

| 编号 | 约束 | 验收方式 |
|------|------|----------|
| T3-01 | Server 必须提供 RPC 方法让管理员触发指定 Agent 的卸载 | Server 通过 tunnel 调用 Agent 的 UninstallAgent RPC |
| T3-02 | Agent 卸载必须清理：系统服务注册、二进制文件、配置文件、数据目录、日志目录 | 卸载后验证 `/usr/sbin/collia`、`/etc/collia/`、`/data/amprobe/resources/collia/`、`/data/amprobe/logs/collia/` 均不存在 |
| T3-03 | Agent 卸载完成后必须断开 tunnel 连接 | 卸载后 Agent 列表中该 Agent 状态变为 offline |
| T3-04 | 本地 `collia remove` 命令必须执行完整清理（当前只移除服务注册） | 本地 remove 后验证文件均已清理 |
| T3-05 | 卸载操作必须有确认机制，禁止误操作 | 远程卸载需要二次确认或特定权限 |

### 禁止发生

| 编号 | 约束 | 验收方式 |
|------|------|----------|
| T3-06 | 禁止卸载后残留可执行文件或可自动重启的服务 | 卸载后验证无 collia 进程、无服务注册 |

## 输入 / 输出

- **输入**：Server 向 Agent 发送 `UninstallAgent` RPC，包含 `{force: bool}`；或本地执行 `collia remove`
- **输出**：Agent 返回卸载结果 `{success, error}`

## 错误语义

| 条件 | 结果 | 幂等 |
|------|------|------|
| 文件删除失败（权限不足） | Agent 上报部分卸载失败，列出残留项 | 否，需人工处理 |
| 服务移除失败 | Agent 上报卸载失败 | 是，可重试 |
| Agent 离线 | Server 返回 Agent 不可达错误 | 是，Agent 上线后重试 |
