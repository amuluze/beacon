# Task: Agent 版本上报与感知

> 依赖: Domain Spec [.specs/domain/agent-lifecycle-update.md](../domain/agent-lifecycle-update.md)
> 状态: done
> 优先级: P0 — 是后续更新/卸载功能的基础

## 设计意图

让 Server 能感知已部署 Agent 的版本信息，为后续版本对比、更新推送提供数据基础。

## 可验证行为约束

### 必须满足

| 编号 | 约束 | 验收方式 |
|------|------|----------|
| T1-01 | Agent 注册帧必须携带 version、os、arch 信息 | Agent 连接后，`s_agent` 表的 Version、OS、Arch 字段非空 |
| T1-02 | Server 侧 `OnAgentConnect` 必须将注册帧中的 version 等信息写入 Agent 记录 | 注册后查询 Agent 列表 API 返回完整 version/os/arch |
| T1-03 | Agent 版本必须是编译时注入的常量，禁止从配置文件读取 | 修改二进制不修改配置，版本仍正确反映 |
| T1-04 | Agent 心跳帧必须携带 version（或 Server 从注册帧缓存，心跳不携带也可，但注册帧必须有） | Agent 上线后 Server 侧版本信息准确 |

### 禁止发生

| 编号 | 约束 | 验收方式 |
|------|------|----------|
| T1-05 | 禁止 Agent 注册帧只携带 agentID 和 joinToken 而不携带版本信息 | 注册帧 Payload 包含可解析的 version 字段 |

## 输入 / 输出

- **输入**：Agent 注册帧 Payload 为 JSON 格式，包含 `agent_id`、`version`、`os`、`arch`、`join_token`
- **输出**：`GET /api/v1/agent/list` 返回的每个 Agent 对象包含 `version`、`os`、`arch` 字段

## 状态影响

- `s_agent` 表新增数据写入：Version、OS、Arch 从注册帧填充
- Agent 注册帧格式变更（从纯 joinToken 变为 JSON）— **breaking change**，旧版 Agent 无法注册

## 错误语义

| 条件 | 结果 | 幂等 |
|------|------|------|
| 注册帧 Payload 不是合法 JSON | Server 返回 FRAME_REGISTER_REJECTED | 是 |
| Version 字段为空字符串 | Server 存储 Version 为空，Agent 列表显示为 "unknown" | 是 |
