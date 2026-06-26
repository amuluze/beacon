# Task: Agent 远程更新推送与自更新

> 依赖: Domain Spec [.specs/domain/agent-lifecycle-update.md](../domain/agent-lifecycle-update.md), Task [.specs/tasks/agent-version-reporting.md](agent-version-reporting.md)
> 状态: pending
> 优先级: P1 — 核心自更新能力

## 设计意图

让 Server 能通过反向 tunnel 主动推送更新指令给指定 Agent，Agent 能安全完成自更新（下载 → 校验 → 原子替换 → 重启）。

## 可验证行为约束

### 必须满足

| 编号 | 约束 | 验收方式 |
|------|------|----------|
| T2-01 | Server 必须提供 RPC 方法让管理员触发指定 Agent 的更新 | Server 通过 tunnel 调用 Agent 的 UpgradeAgent RPC，Agent 开始下载 |
| T2-02 | Agent 必须从 Server 的安装包端点下载新二进制 | 下载 URL 使用已有 `/api/v1/host/install/package` 端点，携带 Install Token |
| T2-03 | Agent 必须对下载的二进制做校验（SHA256），校验失败必须拒绝替换 | 提供错误二进制，Agent 上报校验失败且不替换 |
| T2-04 | Agent 替换必须原子：新二进制写入临时路径 → rename 到目标路径 → 重启 | 替换中断（杀进程）后旧二进制仍可用 |
| T2-05 | Agent 更新完成后必须以新版本重新注册到 Server | 更新后 Agent 列表显示新版本号 |
| T2-06 | 并发更新请求必须被拒绝或排队，禁止并行替换 | 同时发送两个更新请求，第二个返回"更新进行中"错误 |
| T2-07 | 更新进度必须可查询（Server 侧） | 更新期间 Agent 上报进度帧，Server 可查看 |

### 禁止发生

| 编号 | 约束 | 验收方式 |
|------|------|----------|
| T2-08 | 禁止更新导致 Agent 处于中间不可用状态（旧二进制已删除、新二进制未就位） | 断电后 Agent 进程仍可启动 |
| T2-09 | 禁止未认证的更新请求被执行 | 不携带 Install Token 的下载请求被拒绝 |

## 输入 / 输出

- **输入**：Server 向 Agent 发送 `UpgradeAgent` RPC，包含 `{download_url, sha256, version, install_token}`
- **输出**：Agent 返回更新结果 `{success, version, error}`；更新过程上报进度帧

## 状态影响

- Agent 进入 `updating` 状态
- 成功后进入 `registered`（新版本）
- 失败后进入 `update-failed`，回退到旧版本

## 错误语义

| 条件 | 结果 | 幂等 |
|------|------|------|
| 下载失败（网络中断） | Agent 保留旧版本，上报下载失败 | 是，可重试 |
| 校验失败（SHA256 不匹配） | Agent 保留旧版本，上报校验失败 | 是，可重试 |
| Agent 离线 | Server 返回 Agent 不可达错误 | 是，Agent 上线后重试 |
| rename 失败（权限问题） | Agent 保留旧版本，上报替换失败 | 否，需人工干预 |
| 新版本启动崩溃 | 根据回退策略：保留旧版本或标记失败 | 视策略 |

## 回退策略

更新失败后 Agent 必须保留旧二进制继续运行。替换采用两步 rename：
1. `collia_new` → `collia`（原子 rename）
2. 旧二进制在 rename 前先备份为 `collia.bak`
3. 新版本启动成功后删除 `collia.bak`
4. 新版本启动失败时，从 `collia.bak` 回退
