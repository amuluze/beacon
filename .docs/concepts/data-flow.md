# 数据流

## 请求生命周期

本文聚焦输入、模块协调、外部副作用和用户可见结果之间的数据流；项目级摘要见根 `AGENTS.md` 与 `.docs/architecture.md`。

| 阶段 | 责任 | 可观测信号 |
|------|------|------------|
| 输入接收 | Beacon CLI、HTTP/WebSocket API 或 Agent gRPC 注册与数据流接收外部输入 | 路由注册、命令定义、tunnel 帧与导出函数信号 |
| 协调调度 | Fiber handler、service 与 RPC dispatcher 按路由或显式方法名协调模块 | service 调用、RPC 能力面快照与错误映射 |
| 领域处理 | 模块内执行核心规则并维护状态 | Domain Spec、不变量、测试断言 |
| 外部副作用 | 文件、网络、浏览器、数据库或进程操作 | 工具边界、配置项、错误处理 |

## 数据边界

| 边界 | 输入 | 输出 | 约束 |
|------|------|------|------|
| `beacon`（Server） | HTTP/WebSocket 请求、HttpOnly 会话 Cookie、CSRF 双提交值、目标资源或节点标识（`X-Agent-ID` / `agent_id`）、运行时配置、管理员提交的告警通道配置 | JSON/API 响应、Set-Cookie、实时事件、RPC 调用结果、审计或持久化副作用、邮件/钉钉通知 | Fiber 路由 + Cookie JWT/Casbin 中间件 + GORM 持久化 + 反向 gRPC tunnel 客户端；所有容器/主机控制类路由必须携带 Agent-ID；钉钉查询只返回凭据掩码和配置状态 |
| `collia`（Agent） | `config.yml` 中的本机 Unix Socket 控制目标、Agent ID、采集间隔、磁盘/网卡列表和安装 Token | 主机/CPU/内存/磁盘/网络指标、Docker 操作执行结果、反向 gRPC tunnel 服务端会话 | gopsutil 采集 + Docker Engine SDK + 反向 tunnel 服务端；完整采集结束后确定批次时间，**主通道**经控制连接（Unix Socket）上的数据流按批次上报到 Server，数据流不可用或推送失败时回落到已弃用的 HTTP `POST /api/v1/host/report`（兼容期保留），两条通道皆不可用时丢弃本批 |
| `common` | 调用方传入的 schema、查询参数、数据库/传输配置 | 跨模块复用的结构体、RPC 参数/返回值、tunnel 协议封装 | 顶级 Go module（与 `common/docker`、`common/logger`、`common/timex` 平级），承载 database options、RPC schema 与 Unix Socket tunnel 协议 |
| `common/docker` | `IManager` 操作上下文、容器/镜像/网络标识、网络配置（`SubNetworkConfig`） | `ContainerSummary`、`ImageSummary`、`NetworkSummary`、`Version` 等摘要结构 | 独立 Go module `github.com/amuluze/docker`，封装 docker/docker v27.2.1+incompatible；不实现业务持久化 |
| `common/logger` | 日志级别、格式、输出位置（stdout / 文件）、文件轮转配置（`SetLogLevel/Format/Output/File`） | 标准日志输出到控制台或轮转文件；JSON 格式可选 | 独立 Go module `github.com/amuluze/amutool/logger`，基于 zap + file-rotatelogs；不实现业务事件追踪或远程日志聚合 |
| `common/timex` | 无显式输入，调用 `CurrentTime/CurrentDate/CurrentTimestamp/Int64ToTime` 即取当前时间或转换 | 当前时间/日期/时间戳字符串、`time.Time`、`Ticker` 实例 | 独立 Go module `github.com/amuluze/amutool/timex`；不实现 cron 调度或时区业务逻辑 |
| `beacon-web`（Admin SPA） | 路由参数、表单输入、用户操作、API 响应、实时事件、WebSocket 升级 | 页面状态、API 请求、组件事件、Pinia store 更新、实时通道订阅状态、HTTP client 调用 | Vue 3 + Vite + Pinia + axios + xterm.js + Codemirror + UnoCSS + ECharts；调用 `beacon` Server 的 `/api/v1/*` |

## 事件传播

- 事件传播从 Fiber handler 或 tunnel 注册入口携带请求上下文、Agent 标识和权限信息，下游 repository 与 RPC client 不重新猜测目标。
- 跨模块调用应保留错误分类和关联 ID，使 API 响应、日志、审计记录或实时事件可以互相对齐。
- 实时事件、后台任务和外部副作用应有可重试或可降级语义，避免部分成功被展示为完整成功。
- 告警事实以审计写入为准；邮件和钉钉是相互隔离的外部通知副作用，单通道失败只记录降级证据，不回滚审计或阻断另一通道。
- Collia 注册元数据先写入 Agent 记录，再异步进入更新判断，避免在 tunnel 生命周期回调中同步等待反向 RPC。Server 只对低版本或版本为空的旧 Agent 下发与实际发布文件 SHA256 绑定的更新；同版本静默跳过，高版本、未知架构或非 Linux 目标拒绝更新。
- Collia 更新开始审计是下发 RPC 的前置条件；成功和失败均追加系统审计并携带 Agent ID。自动失败进入冷却，管理员手动触发只绕过冷却，不绕过单 Agent 去重、全局并发或完整性校验。
- 资源告警以 Server 写入时间形成完整连续窗口；单点、陈旧、缺失或间断样本进入未知态，不能被解释为触发或恢复证据。CPU、内存、磁盘上报值均为 `0–100` 百分数。
- SMTP 密码、钉钉 Webhook 和签名密钥在持久化前使用部署期数据密钥加密；数据密钥只存在于运行环境，不进入数据库、日志或 API 响应。
- 错误传播保留可定位的错误分类；可重试、幂等或降级语义按具体操作定义，避免把部分成功展示为完整成功。

## 失败语义

| 失败位置 | 表现 | 处理原则 |
|----------|------|----------|
| 输入校验 | 拒绝或返回错误 | 错误信息应可定位问题，不泄露敏感值 |
| 模块边界 | 调用失败或状态不一致 | 保持边界内状态可回滚或可重试 |
| 外部副作用 | I/O、网络或进程失败 | 保留错误证据并避免部分写入造成漂移 |
| 钉钉 Webhook | 地址校验、超时、HTTP 或业务错误 | 只允许官方 HTTPS 地址且禁止重定向；错误不得包含 access token、签名密钥或上游响应正文 |
