# Project Analysis

## Workspace Modules

| Path | Unit | Runtime |
|------|------|---------|
| `beacon` | `beacon` | Go `1.25.0` |
| `collia` | `collia` | Go `1.25.0` |
| `common` | `common` | Go `1.25.0` |
| `common/docker` | `github.com/amuluze/docker` | Go `1.25.0` |
| `common/logger` | `github.com/amuluze/amutool/logger` | Go `1.25.0` |
| `common/timex` | `github.com/amuluze/amutool/timex` | Go `1.21.4` |
| `scripts/doc-schema` | `beacon.local/docs/schema` | Go `1.25.0` |
| `beacon/web/package.json` | `beacon-web` | pinia/typescript/vite/vue |

## Top-Level Directories

| Directory | Inferred Role |
|-----------|---------------|
| `.docs/` | implementation documentation |
| `.specs/` | SDD 规格目录：`domain/` 存长期领域约束，`tasks/` 存过程产物（不随仓库分发） |
| `scripts/` | 文档门禁与校验脚本：Ruby 校验器（证据台账、OpenAPI 契约、门禁自测）与 Go 字段形状导出器 `scripts/doc-schema` |
| `beacon/` | Beacon Server：HTTP/API service module — Web/API 接入、路由注册、请求校验和服务协调 |
| `collia/` | Collia Agent：monitoring agent module — 主机与 Docker 指标采集、反向 gRPC tunnel 服务端、周期上报 |
| `common/` | shared contract library — 跨 Server/Agent 复用的 schema、数据库封装、RPC 参数/返回值、tunnel 协议 |
| `assets/` | 品牌图标资源 |

## Taskfile Commands

- `task beacon`
- `task beacon-web`
- `task collia`
- `task default`
- `task doc:check`
- `task doc:verify`

## Package Scripts

- `beacon/web/package.json (beacon-web)`: `build`, `dev`, `lint`, `lint:fix`, `preview`, `test`, `test:coverage`, `test:run`, `ts`

## Documentation Inventory

### 当前受管 `.docs/` targets

- `.docs/doc-system.md`
- `.docs/api/routes.md`
- `.docs/api/openapi.yml`
- `.docs/architecture.md`
- `.docs/concepts/data-flow.md`
- `.docs/deployment.md`
- `.docs/modules/beacon-web.md`
- `.docs/modules/beacon.md`
- `.docs/modules/collia.md`
- `.docs/modules/common-docker.md`
- `.docs/modules/common-logger.md`
- `.docs/modules/common-timex.md`
- `.docs/modules/common.md`
- `.docs/project-analysis.md`

## Documentation Signals

- 当前未发现未被纳管的生成文档，也未发现仅存在于 CLAUDE 的本地 AI 条目。

## Validation Evidence Inventory

测试库存按Git可见且实际存在的 `*_test.go`、`*.test.ts`、`*.spec.ts` 文件统计并完整展开；门禁自身的Ruby回归由 `ruby scripts/doc-check-test.rb` 独立运行，不计入此Go/TypeScript库存。Domain Spec 中的验收方式不代表测试已经证明整个约束成立。

- Test files discovered: 155

<!-- test-inventory:start -->
- `beacon/pkg/auth/jwtauth/auth_test.go`
- `beacon/pkg/auth/jwtauth/options_test.go`
- `beacon/pkg/auth/jwtauth/store_test.go`
- `beacon/pkg/auth/password_test.go`
- `beacon/pkg/authcookie/manager_test.go`
- `beacon/pkg/contextx/contextx_test.go`
- `beacon/pkg/errors/errors_test.go`
- `beacon/pkg/fiberx/fiberx_test.go`
- `beacon/pkg/fiberx/transport_test.go`
- `beacon/pkg/rpc/disabled_test.go`
- `beacon/pkg/rpc/rpc_test.go`
- `beacon/pkg/secret/cipher_test.go`
- `beacon/pkg/utils/hash/hash_test.go`
- `beacon/pkg/validatex/validatex_test.go`
- `beacon/service/account/repository/account_test.go`
- `beacon/service/account/service/account_test.go`
- `beacon/service/agent/agent_test.go`
- `beacon/service/agent/updater_test.go`
- `beacon/service/agent_install_test.go`
- `beacon/service/alarm/repository/alarm_test.go`
- `beacon/service/alarm/service/alarm_test.go`
- `beacon/service/alarm/title/title_test.go`
- `beacon/service/app_test.go`
- `beacon/service/audit/repository/audit_test.go`
- `beacon/service/audit/service/audit_test.go`
- `beacon/service/auth/api/auth_test.go`
- `beacon/service/auth/repository/auth_test.go`
- `beacon/service/auth/service/auth_test.go`
- `beacon/service/config_test.go`
- `beacon/service/container/api/container_test.go`
- `beacon/service/container/repository/container_test.go`
- `beacon/service/container/service/container_test.go`
- `beacon/service/data_encryption_test.go`
- `beacon/service/dingtalk/client/client_test.go`
- `beacon/service/dingtalk/repository/dingtalk_test.go`
- `beacon/service/dingtalk/service/dingtalk_test.go`
- `beacon/service/doc_routes_test.go`
- `beacon/service/feishu/client/client_test.go`
- `beacon/service/feishu/repository/feishu_test.go`
- `beacon/service/feishu/service/feishu_test.go`
- `beacon/service/health/api/health_test.go`
- `beacon/service/host/repository/host_test.go`
- `beacon/service/host/service/host_test.go`
- `beacon/service/install_report_test.go`
- `beacon/service/mail/repository/mail_test.go`
- `beacon/service/mail/service/mail_test.go`
- `beacon/service/middleware/auth_test.go`
- `beacon/service/middleware/panic_test.go`
- `beacon/service/middleware/ratelimit_test.go`
- `beacon/service/middleware/readonly_test.go`
- `beacon/service/middleware/skip_test.go`
- `beacon/service/model/monitor_test.go`
- `beacon/service/model/session_test.go`
- `beacon/service/prepare_test.go`
- `beacon/service/report/handler_test.go`
- `beacon/service/report/report_test.go`
- `beacon/service/reset_admin_test.go`
- `beacon/service/router_test.go`
- `beacon/service/rpc_test.go`
- `beacon/service/schema/account_test.go`
- `beacon/service/schema/alarm_test.go`
- `beacon/service/schema/container_test.go`
- `beacon/service/system_version_test.go`
- `beacon/service/task/cleanup_test.go`
- `beacon/service/task/task_test.go`
- `beacon/service/terminal/bridge_test.go`
- `beacon/service/terminal/handler_test.go`
- `beacon/service/transport_security_test.go`
- `beacon/service/tunnel_report_test.go`
- `beacon/service/version_check_test.go`
- `beacon/service/wecom/client/client_test.go`
- `beacon/service/wecom/repository/wecom_test.go`
- `beacon/service/wecom/service/wecom_test.go`
- `beacon/web/src/api/agent/index.test.ts`
- `beacon/web/src/api/auth/index.test.ts`
- `beacon/web/src/api/index.test.ts`
- `beacon/web/src/components/Agent/AgentEmptyState.test.ts`
- `beacon/web/src/components/Agent/DataStaleTag.test.ts`
- `beacon/web/src/components/Agent/MonitorContextBar.test.ts`
- `beacon/web/src/components/Agent/MonitorStalePanel.test.ts`
- `beacon/web/src/components/Agent/UpgradeAgentDialog.test.ts`
- `beacon/web/src/components/ConfirmDialog/index.test.ts`
- `beacon/web/src/components/Echarts/index.test.ts`
- `beacon/web/src/components/Message/message.test.ts`
- `beacon/web/src/components/Terminal/index.test.ts`
- `beacon/web/src/config/echarts.test.ts`
- `beacon/web/src/hooks/useAgentSelection.test.ts`
- `beacon/web/src/hooks/useCommandComponent.test.ts`
- `beacon/web/src/hooks/useConfirmCommand.test.ts`
- `beacon/web/src/layout/index.test.ts`
- `beacon/web/src/layout/navbar/Avatar.test.ts`
- `beacon/web/src/layout/navbar/BeaconVersion.test.ts`
- `beacon/web/src/layout/navbar/components/UpdatePassword.test.ts`
- `beacon/web/src/layout/navbar/index.test.ts`
- `beacon/web/src/layout/statusbar/index.test.ts`
- `beacon/web/src/router/dynamic.test.ts`
- `beacon/web/src/store/modules/agent.test.ts`
- `beacon/web/src/store/modules/user.test.ts`
- `beacon/web/src/utils/chartTokens.test.ts`
- `beacon/web/src/utils/convert.test.ts`
- `beacon/web/src/utils/csrf.test.ts`
- `beacon/web/src/utils/monitorChart.test.ts`
- `beacon/web/src/utils/monitorStale.test.ts`
- `beacon/web/src/views/account/role/components/CreateRole.test.ts`
- `beacon/web/src/views/account/user/components/DeleteUser.test.ts`
- `beacon/web/src/views/container/commandWrappers.test.ts`
- `beacon/web/src/views/container/components/BatchActionBar.test.ts`
- `beacon/web/src/views/container/container/containerForm.test.ts`
- `beacon/web/src/views/container/container/containerProtection.test.ts`
- `beacon/web/src/views/container/index.test.ts`
- `beacon/web/src/views/login/index.test.ts`
- `beacon/web/src/views/monitor/container/index.test.ts`
- `beacon/web/src/views/monitor/host/index.test.ts`
- `beacon/web/src/views/monitor/index.test.ts`
- `beacon/web/src/views/setting/alarm/components/settingsDialogs.test.ts`
- `beacon/web/src/views/setting/alarm/index.test.ts`
- `beacon/web/src/views/setting/components/AuditLog.test.ts`
- `beacon/web/src/views/setting/index.test.ts`
- `beacon/web/src/views/setting/probe/index.test.ts`
- `collia/cmd/main_test.go`
- `collia/pkg/psutil/psutil_test.go`
- `collia/pkg/utils/file_test.go`
- `collia/service/config_test.go`
- `collia/service/logger_test.go`
- `collia/service/preflight_test.go`
- `collia/service/report/client_test.go`
- `collia/service/report_channel_test.go`
- `collia/service/report_tls_test.go`
- `collia/service/rpc/container_test.go`
- `collia/service/rpc/dispatcher_test.go`
- `collia/service/rpc/docker_test.go`
- `collia/service/rpc/file_test.go`
- `collia/service/rpc/lifecycle_test.go`
- `collia/service/rpc/terminal_test.go`
- `collia/service/task_test.go`
- `collia/service/transport_security_test.go`
- `common/database/db_log_test.go`
- `common/database/db_test.go`
- `common/database/permissions_test.go`
- `common/database/queryoption_test.go`
- `common/docker/conflict_test.go`
- `common/docker/container_test.go`
- `common/docker/image_test.go`
- `common/docker/manager_test.go`
- `common/docker/network_test.go`
- `common/docker/version_test.go`
- `common/logger/slog_test.go`
- `common/rpc/schema/terminal_test.go`
- `common/rpc/tunnel/agent_test.go`
- `common/rpc/tunnel/data_stream_test.go`
- `common/rpc/tunnel/server_test.go`
- `common/rpc/tunnel/unix_socket_test.go`
- `common/timex/ticker_test.go`
- `common/timex/timex_test.go`
- `scripts/doc-schema/main_test.go`
<!-- test-inventory:end -->

### Domain Spec Evidence Notes

逐条映射覆盖全部85个编号；仅记录下列断言范围，不把相关测试通过解释为整个约束已满足。
覆盖状态分为「部分覆盖」与「待专项验收」；缺口列是尚须完成的行为验收。

### Evidence Backlink Table

| Domain Ref | Evidence | 已有断言或审查入口 | 覆盖状态 | 尚缺证据 |
|------------|----------|-------------------|----------|----------|
| IA001 | `beacon/service/schema/account_test.go::TestUserDeleteArgsRequiresExplicitUUIDTargets` | 空数组、空元素、非法 UUID 与合法目标 | 部分覆盖 | 持久化层未被调用需集成验证 |
| IA002 | `beacon/service/account/service/account_test.go::TestAccountService_UserDelete`；`beacon/service/account/service/account_test.go::TestAccountService_UserDeleteRejectsMissingTarget` | 删除存在目标、拒绝缺失目标及授权清理 | 部分覆盖 | 跨数据库事务失败回滚需专项验证 |
| IA003 | `beacon/web/src/views/account/user/components/DeleteUser.test.ts::keeps the dialog and list state when deletion fails` | 失败后弹窗保留且不刷新列表 | 部分覆盖 | 真实浏览器错误提示需端到端验证 |
| IA004 | `beacon/service/auth/api/auth_test.go::TestLoginSetsHttpOnlyCookiesWithoutReturningTokens`；`beacon/service/auth/api/auth_test.go::TestRefreshWithoutAccessCookieStillRequiresCSRF`；`beacon/service/middleware/auth_test.go::TestCSRFMiddlewareProtectsCookieAuthenticatedMutations` | Cookie 属性、无正文Token、无Access Cookie刷新及双提交拒绝 | 部分覆盖 | 所有入口和浏览器存储需专项审计 |
| IA005 | `beacon/service/prepare_test.go::TestInitAccountCreatesOnlyMandatoryChangeAdministrator`；`beacon/service/prepare_test.go::TestInitAccountMarksHistoricalDefaultForMandatoryChange` | 只创建管理员、强制改密及旧默认密码迁移 | 部分覆盖 | 完整升级安装链路待验收 |
| IA006 | `beacon/service/prepare_test.go::TestInitAccountRejectsWeakBootstrapPassword`；`beacon/pkg/auth/password_test.go::TestValidatePassword`；`beacon/service/schema/account_test.go::TestPasswordUpdateArgsRejectsShortNewPassword` | 弱bootstrap密码、密码字符类别及更新长度边界 | 部分覆盖 | 完整密码更新HTTP回归与Unicode类别边界需专项验证 |
| IA007 | `beacon/service/reset_admin_test.go::TestResetAdminPasswordGeneratesTemporaryPassword`；`beacon/service/reset_admin_test.go::TestResetAdminPasswordExplicitPassword`；`beacon/service/reset_admin_test.go::TestResetAdminPasswordRejectsWeakPassword`；`beacon/service/reset_admin_test.go::TestResetAdminPasswordMissingAdmin` | 临时/显式密码、bcrypt、强制改密、系统审计、弱密码不改原凭据与管理员缺失拒绝 | 部分覆盖 | 真实 CLI 在服务运行期间的数据库锁竞争与写失败场景需集成验收 |
| RA001 | `beacon/web/src/views/account/user/components/DeleteUser.test.ts::submits the selected ids and refreshes only after deletion succeeds` | 传递选中ID且仅成功后刷新 | 部分覆盖 | 权限范围需API集成验证 |
| RA002 | `beacon/service/account/service/account_test.go::TestAccountService_UserDeleteRejectsMissingTarget`；`beacon/service/schema/account_test.go::TestUserDeleteArgsRequiresExplicitUUIDTargets` | 无效输入与缺失目标拒绝 | 部分覆盖 | 持久化失败错误码需补集成场景 |
| RA003 | `beacon/web/src/views/account/user/components/DeleteUser.test.ts::submits the selected ids and refreshes only after deletion succeeds` | 确认目标与请求ID一致 | 部分覆盖 | 多选及取消并发需端到端验证 |
| RA004 | `beacon/service/middleware/auth_test.go::TestUserAuthMiddlewareRestrictsMandatoryPasswordChangeSession` | 强制改密仅允许受限路由 | 部分覆盖 | 旧会话失效及浏览器跳转需联测 |
| IU001 | `beacon/service/agent/agent_test.go::TestOnAgentConnectPersistsRegistrationMetadataOnExistingRow`；`common/rpc/tunnel/agent_test.go::TestAgentTunnelRegistrationPayloadIncludesMetadata` | 注册帧元数据和现有记录版本持久化 | 部分覆盖 | 非法版本注册拒绝需补测试 |
| IU002 | `beacon/service/agent/updater_test.go::TestCompareSemVersions` | 版本比较序关系 | 部分覆盖 | 非语义化构建序号不在该测试范围 |
| IU003 | `collia/service/rpc/lifecycle_test.go::TestReplaceBinaryAtomicWithBackup`；`collia/service/rpc/lifecycle_test.go::TestReplaceBinaryFailureKeepsCurrentPath` | 原子替换、备份及rename失败保留旧路径 | 部分覆盖 | 进程被杀与断网后的重启需系统验收 |
| IU004 | `beacon/service/agent/updater_test.go::TestUpdaterTriggerComputesChecksumCallsAgentAndAudits`；`beacon/service/agent/updater_test.go::TestUpdaterRPCFailureWritesTruncatedFailureAudit` | 更新成功失败均有审计 | 部分覆盖 | 卸载审计链路需专项验证 |
| IU005 | `beacon/service/agent/updater_test.go::TestUpdaterTriggerComputesChecksumCallsAgentAndAudits`；`collia/service/rpc/lifecycle_test.go::TestHandleUpgradeSHA256MismatchKeepsOld` | 发布文件校验值传递及篡改拒绝 | 部分覆盖 | 发布源信任链需部署验收 |
| IU006 | `beacon/service/agent/agent_test.go::TestServiceListDoesNotExposePersistedOnlineWithoutLiveTunnel`；`beacon/service/agent/agent_test.go::TestServiceListProjectsLiveTunnelAsOnline` | 只按实时tunnel投影在线状态 | 部分覆盖 | 更新窗口中的列表时序需集成验证 |
| IU007 | `beacon/service/router_test.go::TestAgentUpgradeRouteIsAdminOnly` | 升级路由管理员资源校验 | 部分覆盖 | 非管理员界面隐藏需组件证据 |
| IU008 | `collia/service/rpc/lifecycle_test.go::TestVerifyUpgradeSignatureMatrix`；`collia/service/rpc/lifecycle_test.go::TestHandleUpgradeWithPinnedKeyRequiresValidSignature`；`collia/service/rpc/lifecycle_test.go::TestSetUpgradePublicKeys` | 固定公钥缺签/错签拒绝、未固定有签拒绝、双方齐备验签放行且失败保留旧二进制 | 部分覆盖 | 真实跨版本升级与公钥轮换流程待联调 |
| IU009 | `common/rpc/tunnel/unix_socket_test.go::TestTunnelOverUnixSocket`；`collia/service/transport_security_test.go::TestValidateTransportSecurity`；`beacon/service/transport_security_test.go::TestValidateTransportSecurityProduction` | Unix Socket 承载控制通道：注册与双向 RPC 可用、socket 0600、Stop 清理；缺监听/下发路径或非 unix 目标均拒绝启动 | 部分覆盖 | 单机真实部署的断线重连与长期运行待观察 |
| RU001 | `beacon/service/agent/updater_test.go::TestUpdaterTriggerComputesChecksumCallsAgentAndAudits`；`collia/service/rpc/lifecycle_test.go::TestHandleUpgradeSuccess` | 服务端调用及Agent更新成功 | 部分覆盖 | 真实反向tunnel下载联测待验收 |
| RU002 | `collia/service/rpc/lifecycle_test.go::TestHandleUpgradeSHA256MismatchKeepsOld`；`collia/service/rpc/lifecycle_test.go::TestVerifySHA256RejectsEmpty` | 校验错误或缺失时拒绝替换 | 部分覆盖 | 运行进程连续性需系统验证 |
| RU003 | `beacon/service/router_test.go::TestAgentUpgradeRouteIsAdminOnly`；`common/rpc/tunnel/server_test.go::TestE2E_RegistrationRejectsUnauthorized` | 管理员路由及未授权注册拒绝 | 部分覆盖 | 所有更新入口权限需专项审计 |
| RU004 | `collia/service/rpc/lifecycle_test.go::TestHandleUninstallReplyConsistency`；`collia/cmd/main_test.go::TestIsColliaDaemonCommand` | 卸载应答一致性与残留 daemon 进程识别 | 部分覆盖 | 需在隔离 Linux 环境验收 systemd 停止、daemon-reload 与文件清理 |
| RU005 | `collia/service/rpc/lifecycle_test.go::TestScheduleRestartLeavesReplyWindow`；`collia/service/rpc/terminal_test.go::TestTerminalSessionStream_ContextCancel` | 应答窗口与终端取消 | 部分覆盖 | 更新中已有RPC优雅完成需系统联测 |
| RU006 | `collia/service/rpc/lifecycle_test.go::TestHandleUpgradeSerialRejection` | 第二个升级并发请求拒绝 | 部分覆盖 | 多进程并发不在该测试范围 |
| RU007 | `beacon/service/agent/updater_test.go::TestUpdaterVersionDecisions` | 空、低、同、高版本决策 | 部分覆盖 | 实际版本发布矩阵需验收 |
| RU008 | `beacon/service/agent/updater_test.go::TestUpdaterAutomaticScheduleDeduplicatesAndCoolsDownFailures`；`beacon/service/agent/updater_test.go::TestUpdaterGlobalConcurrencyLimit` | 单目标去重、失败冷却和全局并发上限 | 部分覆盖 | 重连风暴负载需专项验收 |
| RU009 | `beacon/service/agent/updater_test.go::TestUpdaterUnknownArchitectureAuditsAndSkips` | 未知架构不调用RPC且写审计 | 部分覆盖 | 所有支持架构发布物需验收 |
| RU010 | `beacon/service/agent/agent_test.go::TestUpgradeAPIValidatesAgentID` | API拒绝非法或空目标 | 部分覆盖 | 未选节点的界面禁用需组件证据 |
| RU011 | `beacon/service/agent/updater_test.go::TestRepositoryAutoUpdateEnabledDefaultsTrueAndPersists`；`beacon/service/agent/updater_test.go::TestUpdaterScheduleHonorsConsoleAutoUpdateSetting`；`beacon/service/agent/agent_test.go::TestUpdateSettingAPIRejectsMissingFieldAndPersistsToggle`；`beacon/web/src/views/setting/probe/index.test.ts::persists the auto-update switch when toggled off` | 默认开启、关闭后不自动调度、手动升级不受影响 | 部分覆盖 | 部署级门禁与控制台开关的浏览器联测需专项验证 |
| RU012 | `beacon/service/agent/updater_test.go::TestUpdaterPassesSidecarSignatureToAgent`；`beacon/service/agent_install_test.go::TestBuildColliaConfigEmbedsUpgradePublicKey` | Server 透传侧车签名、公钥随安装 config 下发并由 Agent 固定 | 部分覆盖 | 发布管线生成侧车签名与密钥轮换待发布环境验收 |
| RU013 | `collia/service/rpc/lifecycle_test.go::TestVerifyUpgradeSignatureAcceptsRotatedKeyRing`；`collia/service/rpc/lifecycle_test.go::TestSetUpgradePublicKeys` | 轮换窗口内新旧签名均通过、集合外公钥拒绝、公钥只经 config 注入 | 部分覆盖 | 升级 RPC 无公钥字段为结构检查；真实轮换流程待部署验收 |
| IC001 | `collia/service/rpc/dispatcher_test.go::TestNewDispatcher_RegistersStandardMethods` | 调度能力快照 | 部分覆盖 | 显式接口与动态绕过需代码审计 |
| IC002 | `beacon/service/container/repository/container_test.go::TestImageDeleteInvalidatesSelectedAgentCacheAfterRPCSuccess`；`beacon/service/container/repository/container_test.go::TestImageDeleteKeepsCacheWhenAgentRPCFails` | 成功后失效目标缓存、失败保留缓存 | 部分覆盖 | 真实Docker删除联测待验收 |
| IC003 | `beacon/service/container/repository/container_test.go::TestContainerListReturnsLatestBatchForSelectedAgent` | 列表仅含目标最新批次 | 部分覆盖 | 镜像占用计数需另行验证 |
| IC004 | `beacon/service/report/report_test.go::TestStorePersistsReportBatch` | 监控批次存储 | 部分覆盖 | 多批历史容器趋势保留需专项测试 |
| IC005 | `beacon/service/container/service/container_test.go::TestContainerService_ContainerCreatePassesRestartPolicy`；`beacon/service/container/service/container_test.go::TestContainerService_ContainerUpdatePassesRestartPolicy` | 创建和更新透传重启策略 | 部分覆盖 | 非法策略拒绝需schema测试 |
| IC006 | `beacon/service/container/service/container_test.go::TestContainerService_RejectsProtectedContainerOperations`；`beacon/service/container/service/container_test.go::TestContainerService_FailsClosedOnLookupError`；`beacon/service/container/api/container_test.go::TestContainerRemoveMapsProtectedRejectionToForbidden` | 保护容器启停删除拒绝、查询失败关闭，且拒绝语义经 HTTP 层映射为 403 | 部分覆盖 | 真实 Agent 目标与批量选择需端到端验证 |
| IC007 | `collia/service/rpc/container_test.go::TestContainerCreateErrorMapsNameConflict`；`common/docker/conflict_test.go::TestIsNameConflict` | 容器删除以运行时为唯一事实源：重名归一为可定位错误、无本地记录判定 | 部分覆盖 | 删除失败端到端错误传播与下一批次列表收敛待部署验收 |
| RC001 | `common/docker/image.go` | 悬空镜像清理实现 | 待专项验收 | 清理悬空镜像排除被引用镜像尚无直接测试 |
| RC002 | `beacon/web/src/views/container/index.vue` | 容器破坏性操作界面入口 | 待专项验收 | 容器及镜像破坏性操作须单独组件验收 |
| I001 | `beacon/service/host/repository/host_test.go::TestHostInfoRequiresAgentID`；`beacon/service/container/repository/container_test.go::TestContainerUpdateRejectsMissingAgentIDBeforeRPC` | 缺目标时拒绝查询和控制 | 部分覆盖 | 所有Agent相关接口需逐路由集成验证 |
| I002 | `beacon/service/host/repository/host_test.go::TestHostInfoMarksStaleData`；`beacon/service/task/task_test.go::TestCPUAlarmIgnoresStaleLatestSample` | 主机陈旧标识和告警跳过陈旧值 | 部分覆盖 | 所有页面渲染需浏览器验证 |
| I003 | `common/rpc/tunnel/server_test.go::TestE2E_MultiAgentIsolation`；`beacon/service/container/repository/container_test.go::TestContainerUpdateCallsAgentRPC` | tunnel多节点隔离及更新RPC目标 | 部分覆盖 | 各破坏性操作及审计需逐项集成验证 |
| I004 | `beacon/service/task/task_test.go::TestCPUAlarmRequiresCompleteContinuousWindow`；`beacon/service/task/task_test.go::TestCPUAlarmRecoversAfterMetricDrops` | 阈值连续窗口和恢复判断 | 部分覆盖 | 全资源矩阵需专项验收 |
| I005 | `beacon/service/report/report_test.go::TestStorePersistsReportBatch`；`beacon/service/report/report_test.go::TestStoreRejectsMissingAgentID` | 合法批次持久化与缺目标拒绝 | 部分覆盖 | 中间表写失败事务回滚尚需故障注入 |
| I006 | `beacon/service/agent/agent_test.go::TestServiceListDoesNotExposePersistedOnlineWithoutLiveTunnel`；`beacon/service/agent/agent_test.go::TestSetTunnelResetsPersistedConnectionStatus` | 进程状态投影及旧在线标记清理 | 部分覆盖 | 前端空列表需联合验证 |
| I007 | `beacon/service/host/repository/host_test.go::TestUsageQueriesRespectClosedTimeWindow`；`beacon/service/container/repository/container_test.go::TestContainerUsageRespectsClosedTimeWindow` | 查询开始结束边界 | 部分覆盖 | 5/10/30分钟图表横轴需组件证据 |
| I008 | `beacon/service/dingtalk/service/dingtalk_test.go::TestDingTalkServiceQueryMasksCredentials`；`beacon/service/dingtalk/client/client_test.go::TestWebhookSenderTimeoutDoesNotLeakWebhook` | 查询掩码和超时不泄露Webhook | 部分覆盖 | 全部错误日志需专项审计 |
| I009 | `beacon/service/prepare_test.go::TestInitAccountSyncsNewRoutesToExistingAdminRole` | 已有角色路由资源同步和旧策略清理 | 部分覆盖 | 生产升级权限矩阵需验收 |
| I010 | `beacon/service/prepare_test.go::TestInitAlarmThresholdNormalizesRecordsAndPreservesOverrides` | 别名、重复规范化且保留覆盖值 | 部分覆盖 | 跨数据库迁移需验收 |
| I011 | `beacon/service/prepare_test.go::TestInitAlarmMessageSettingCreatesDefaultAndPreservesOverride`；`beacon/service/task/task_test.go::TestNotificationFallsBackToDefaultAlarmMessageTitle` | 默认标题及覆盖保留、通知回退 | 部分覆盖 | 双通道自定义标题与重启需集成验证 |
| I012 | `beacon/service/task/task_test.go::TestMemoryAlarmUsesZeroToOneHundredPercentScale`；`beacon/service/task/task_test.go::TestDiskAlarmUsesPercentScaleAndDurationWindow` | 内存磁盘百分数无重复乘100 | 部分覆盖 | CPU采样到显示链路需集成验证 |
| I013 | `beacon/service/task/task_test.go::TestCPUAlarmRequiresCompleteContinuousWindow`；`beacon/service/task/task_test.go::TestCPUAlarmIgnoresStaleLatestSample` | 完整连续窗口、间断及陈旧样本 | 部分覆盖 | 全资源窗口矩阵需验收 |
| I014 | `beacon/service/task/task_test.go::TestAlarmRecoveryNotifiesOnceAndAllowsRetrigger`；`beacon/service/task/task_test.go::TestCPUAlarmRecoversAfterMetricDrops`；`beacon/service/task/task_test.go::TestEvaluateAlarmWindowStrictRecovery` | 单次恢复、清去重后可重触发、完整正常窗口才恢复（混合窗口不判恢复） | 部分覆盖 | 全资源恢复矩阵需验收 |
| I015 | `beacon/service/report/handler_test.go::TestHandleReport_RejectsWhenServerTokenIsNotConfigured`；`beacon/service/report/handler_test.go::TestHandleReport_RejectsQueryTokenFallback`；`beacon/service/report/handler_test.go::TestHandleReport_AcceptsWithValidToken` | 空配置拒绝、拒绝query及接受专用头 | 部分覆盖 | 代理转发配置需部署验证 |
| I016 | `beacon/service/mail/repository/mail_test.go::TestMailRepositoryEncryptsAndDecryptsPassword`；`beacon/service/dingtalk/repository/dingtalk_test.go::TestDingTalkRepositoryTransparentlyMigratesPlaintext`；`beacon/service/data_encryption_test.go::TestMigrateSensitiveDataEncryptsLegacyRowsAcrossChannels` | 加密存储、读取时旧明文迁移与四渠道启动统一迁移 | 部分覆盖 | 版本控制及所有日志明文审计待验收 |
| I017 | `beacon/service/router_test.go::TestHostRoutesExcludeRemovedSystemSettingsAPI`；`collia/service/rpc/dispatcher_test.go::TestNewDispatcher_RegistersStandardMethods` | HTTP和RPC无已移除系统设置能力 | 部分覆盖 | 前端路由需联合验证 |
| I018 | `beacon/web/src/views/monitor/container/index.test.ts::keeps high-cardinality legends scrollable without duplicating them inside charts` | 高基数图例滚动、内部图例关闭及系列保留 | 部分覆盖 | 真实浏览器长名称提示及小屏布局需视觉验收 |
| I019 | `beacon/service/report/handler_test.go::TestHandleReport_RejectsInvalidToken` | 错误Token拒绝 | 部分覆盖 | 恒定时间原语审查及时间差实验未覆盖 |
| I020 | `collia/service/report/client_test.go::TestPushRejectsTokenOverHTTPBeforeNetwork` | HTTP目标发送前拒绝 | 部分覆盖 | userinfo地址需独立用例 |
| I021 | `collia/service/report_tls_test.go::TestBuildReportTLSConfigUsesSystemTrustOnly`；`collia/service/report_tls_test.go::TestReportClientRejectsUntrustedSelfSignedCert` | 上报 TLS 使用系统信任库、不跳过校验且拒绝不受信自签证书 | 部分覆盖 | 受信私有 CA 导入系统信任库后的端到端上报需部署验收 |
| I022 | `beacon/service/health/api/health_test.go::TestLiveness_Always200`；`beacon/service/health/api/health_test.go::TestReadiness_DBError_Returns503`；`beacon/service/health/api/health_test.go::TestReadiness_DBSuccess_TunnelQueried` | 存活独立于DB、就绪反映DB | 部分覆盖 | 数据库故障恢复时序需集成验证 |
| I023 | `beacon/service/middleware/ratelimit_test.go::TestGlobalRateLimitMiddleware_SkipsStaticAndAuth` | 全局限流跳过登录类端点（登录限流由 router 内联承担） | 部分覆盖 | 登录超限与上报按 Agent 限流需独立集成测试 |
| I024 | `common/rpc/tunnel/server_test.go::TestServerTunnel_StreamCall_CtxCancelReleasesStream`；`beacon/service/terminal/bridge_test.go::TestBridgeWriteTimeoutDefaultsToGuard` | 流取消释放及写超时兜底 | 部分覆盖 | 半开TCP、超龄与空闲回收需负载实验 |
| I025 | `beacon/service/app_test.go::TestFiberTimeoutsRespectConfig`；`beacon/service/app_test.go::TestFiberBodyLimitReduced` | 连接超时与请求体上限配置 | 部分覆盖 | 真实慢速请求及超限请求需网络测试 |
| I026 | `beacon/service/audit/repository/audit_test.go::TestAuditQueryFiltersByAgentID`；`beacon/service/task/task_test.go::TestCPUAlarmRecoversAfterMetricDrops` | 审计按 Agent 筛选；告警系统审计携带 AgentID 与粒度键 Object | 部分覆盖 | 启停删除审计的对象及 Agent 内容需端到端验证 |
| I027 | `beacon/service/task/cleanup_test.go::TestCleanupDeletesExpiredAudit`；`beacon/service/task/cleanup_test.go::TestCleanupUsesSeparateRetentionForAuditAndTimeSeries` | 审计按独立保留期清理 | 部分覆盖 | 定时任务调度需部署验收 |
| I028 | `collia/service/report_tls_test.go::TestBuildReportTLSConfigUsesSystemTrustOnly` | 空控制配置下上报 TLS 使用系统信任库，无控制证书或 ServerName 耦合 | 部分覆盖 | 真实安装配置中的控制面与 Web TLS 解耦需持续代码审查 |
| I029 | `common/rpc/tunnel/data_stream_test.go::TestServerTunnel_DataStream_DoesNotRegisterAgent`；`common/rpc/tunnel/data_stream_test.go::TestServerTunnel_DataStream_RejectsBadToken`；`common/rpc/tunnel/data_stream_test.go::TestServerTunnel_DataRegistrationPayloadCarriesRole` | 数据流角色注册、JoinToken 拒绝、不写在线表且不触发生命周期回调 | 部分覆盖 | 真实 Agent 与 Server 同连接双流的端到端联测待执行 |
| I030 | `beacon/service/tunnel_report_test.go::TestTunnelReportHandler_StoresBatch`；`beacon/service/tunnel_report_test.go::TestTunnelReportHandler_RejectsAgentIDMismatch` | 隧道批次入库与身份不一致拒绝 | 部分覆盖 | 真实隧道上的批次投递待部署环境验收 |
| I031 | `collia/service/report_channel_test.go::TestPushReport_UsesTunnelDataStream`；`collia/service/report_channel_test.go::TestPushReport_FallsBackToHTTPWhenDataStreamUnavailable`；`collia/service/report_channel_test.go::TestPushReport_FallsBackToHTTPWhenTunnelPushFails`；`collia/service/report_channel_test.go::TestPushReport_DropsBatchWhenNoChannelAvailable` | 隧道优先、失败回落、丢弃不排队 | 部分覆盖 | 断线期间的真实数据缺口与恢复待部署观察 |
| R001 | `beacon/pkg/fiberx/fiberx_test.go::TestServiceError_Mapping` | 缺目标、离线、超时映射不同HTTP错误 | 部分覆盖 | 真实API权限不足需集成测试 |
| R002 | `beacon/service/container/repository/container_test.go::TestContainerUpdateRejectsInvalidAgentIDBeforeRPC`；`common/rpc/tunnel/server_test.go::TestE2E_MultiAgentIsolation` | 非法目标拒绝及tunnel隔离 | 部分覆盖 | 全部破坏性操作权限上下文需验收 |
| R003 | `beacon/service/host/repository/host_test.go::TestHostInfoMarksStaleData` | 陈旧主机数据标记 | 部分覆盖 | 采集失败后的页面降级需联测 |
| R004 | `beacon/service/task/task_test.go::TestCPUAlarmTaskScopesDataByAgent` | 按目标Agent读取告警指标 | 部分覆盖 | 通知内容与审计关联需专项验证 |
| R005 | `beacon/service/report/handler_test.go::TestHandleReport_RejectsMissingAgentID`；`beacon/service/report/report_test.go::TestStoreRejectsMissingAgentID` | handler及存储均拒绝缺目标 | 部分覆盖 | 真实上报鉴权链路需部署验证 |
| R006 | `collia/service/rpc/dispatcher_test.go::TestDispatcher_UnknownMethod` | 未知控制方法返回错误 | 部分覆盖 | 全部未实现方法需契约审计 |
| R007 | `beacon/service/task/task_test.go::TestNotificationChannelsDoNotBlockEachOther` | 单通道失败不阻断另一通道 | 部分覆盖 | 外部SMTP和钉钉真实故障待验收 |
| R008 | `beacon/service/dingtalk/client/client_test.go::TestValidateWebhookURL`；`beacon/service/dingtalk/service/dingtalk_test.go::TestDingTalkServiceRejectsSSRFWebhook` | URL允许集与SSRF拒绝 | 部分覆盖 | 重定向网络行为需专项核对 |
| R009 | `beacon/service/alarm/repository/alarm_test.go::TestAlarmUpdateUsesIDAndCanonicalType` | ID及规范类型匹配才更新 | 部分覆盖 | 跨数据库零值语义需验收 |
| R010 | `beacon/service/schema/alarm_test.go::TestAlarmMessageSettingUpdateArgsValidation`；`beacon/service/alarm/service/alarm_test.go::TestAlarmService_AlarmMessageSettingUpdateRejectsInvalidTitle` | 空白、换行和长度拒绝 | 部分覆盖 | 展示层转义需浏览器验证 |
| R011 | `beacon/service/task/task_test.go::TestCPUAlarmTaskDeduplicatesRepeatedHighSamples`；`beacon/service/task/task_test.go::TestAlarmRecoveryNotifiesOnceAndAllowsRetrigger` | 异常去重及恢复后重触发 | 部分覆盖 | 真实10分钟时钟边界需专项验证 |
| R012 | `beacon/service/task/task_test.go::TestServiceTaskIgnoresStaleContainerState` | 陈旧容器状态不触发通知 | 部分覆盖 | 恰好2分钟边界需专项验证 |
| R013 | `.docs/api/openapi.yml`；`beacon/service/report/report.go` | 兼容上报接口标记弃用且保留可用 | 部分覆盖 | 旧 Agent 实际回落到 HTTP 的端到端联测待执行 |

### Module Command Evidence

<!-- command-evidence:start -->
| ID | 来源 | 工作目录 | 命令 | 结果 | 执行证据或未执行原因 |
|----|------|----------|------|------|--------------------|
| ca1b1e43b37a | `.docs/modules/beacon.md` | `beacon` | task dev | 未执行 | 常驻服务会使用本地配置并写数据库；仅 task --dry dev 目标解析成功，不冒充运行验收 |
| bfd4449277bb | `.docs/modules/beacon.md` | `beacon` | task wire | 通过 | 2026-11-24，退出码0 |
| f9e43e24f520 | `.docs/modules/beacon.md` | `beacon` | task bin | 通过 | 2026-11-24，退出码0 |
| 63920059decf | `.docs/modules/beacon.md` | `beacon` | task amd64 | 未执行 | 多阶段跨架构 Docker 镜像构建留待发布环境；本次 task --dry 目标解析成功 |
| 9fd60073104f | `.docs/modules/beacon.md` | `beacon` | task arm64 | 未执行 | 多阶段跨架构 Docker 镜像构建留待发布环境；本次 task --dry 目标解析成功 |
| 855eec3e44f8 | `.docs/modules/beacon.md` | `beacon` | docker buildx build --platform linux/arm64 -f Dockerfile  --build-arg APK_REPOSITORY=<trusted-alpine-mirror>  --build-arg BEACON_VERSION="$(cat VERSION)"  --build-arg COLLIA_VERSION="$(cat ../collia/VERSION)"  -t beacon:size-check -o type=docker .. | 未执行 | 多阶段跨架构 Docker 镜像构建需要发布环境与可信镜像源；本次仅 task --dry 目标解析成功 |
| 6158a2401d66 | `.docs/modules/beacon.md` | `beacon` | go test -race ./... | 通过 | 2026-11-24，退出码0 |
| 84285f5ff989 | `.docs/modules/beacon.md` | `beacon` | go vet ./... | 通过 | 2026-11-24，退出码0 |
| c6883cb6e39a | `.docs/modules/collia.md` | `collia` | task wire | 通过 | 2026-11-24，退出码0 |
| 8631159c516b | `.docs/modules/collia.md` | `collia` | task amd64 | 通过 | 2026-11-24，退出码0 |
| 9993e4db33d6 | `.docs/modules/collia.md` | `collia` | task arm64 | 通过 | 2026-11-24，退出码0 |
| e346b91b0b48 | `.docs/modules/collia.md` | `collia` | CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=$(cat VERSION)" -o collia ./cmd | 通过 | 2026-11-24，退出码0（collia amd64 二进制构建） |
| 319eae8bb309 | `.docs/modules/common-docker.md` | `common/docker` | go build ./... | 通过 | 2026-11-24，退出码0 |
| 1f1bb851ce6c | `.docs/modules/common-docker.md` | `common/docker` | go test -race ./... | 通过 | 2026-11-24，退出码0 |
| 5363a3a11def | `.docs/modules/common-docker.md` | `common/docker` | go vet ./... | 通过 | 2026-11-24，退出码0 |
| e7f3227f65e7 | `.docs/modules/common-logger.md` | `common/logger` | go build ./... | 通过 | 2026-11-24，退出码0 |
| b415c412e825 | `.docs/modules/common-logger.md` | `common/logger` | go test -race ./... | 通过 | 2026-11-24，退出码0 |
| 428d6d571720 | `.docs/modules/common-logger.md` | `common/logger` | go vet ./... | 通过 | 2026-11-24，退出码0 |
| f3de9362b869 | `.docs/modules/common-timex.md` | `common/timex` | go build ./... | 通过 | 2026-11-24，退出码0 |
| 75ca410dee8b | `.docs/modules/common-timex.md` | `common/timex` | go test -race ./... | 通过 | 2026-11-24，退出码0 |
| 1fba994dee3d | `.docs/modules/common-timex.md` | `common/timex` | go vet ./... | 通过 | 2026-11-24，退出码0 |
| 449c3a070636 | `.docs/modules/common.md` | `common` | go build ./... | 通过 | 2026-11-24，退出码0 |
| e3bff37b6538 | `.docs/modules/common.md` | `common` | go test -race ./... | 通过 | 2026-11-24，退出码0 |
| 25c5581545e4 | `.docs/modules/common.md` | `common` | go vet ./... | 通过 | 2026-11-24，退出码0 |
<!-- command-evidence:end -->

补充执行：六个Go模块（`beacon`、`collia`、`common`、`common/docker`、`common/logger`、`common/timex`）均执行 `go build ./...`、`go test -race ./...`、`go vet ./...`，退出码0；`scripts/doc-schema` 以 `GOWORK=off` 执行 build/test/vet，退出码0。`beacon` 的 `task wire` 与 `task bin`、`collia` 的 `task wire`/`task amd64`/`task arm64` 均退出码0。多阶段跨架构 Docker 镜像构建与常驻服务未在本轮执行，保留未执行事实。

## Domain Specs

- `.specs/domain/account-administration.md`
- `.specs/domain/agent-lifecycle-update.md`
- `.specs/domain/container-runtime.md`
- `.specs/domain/monitoring-platform.md`

## Documentation Coverage Model

项目文档按以下所有权模型组织：

| Layer | Path | Purpose |
|-------|------|---------|
| Domain constraints | `.specs/domain/` | Domain Specs describe durable constraints: what must be true, not how to implement it. |
| Design intent | `DESIGN.md` | 设计意图 SSOT：视觉基调、Token 体系、布局与组件规范、变更治理与视觉验收。 |
| Implementation description | `.docs/` | Project docs describe the current implementation, architecture, modules, deployment, APIs and concepts. |
| AI collaboration entry | `AGENTS.md` + `CLAUDE.md` | Entry files provide navigation, commands and collaboration constraints; `CLAUDE.md` is synchronized from `AGENTS.md`. |

## Update Guidance

- `.docs/` 由 LLM 生成与人工维护共同产出：生成负责从源码、配置与脚本提取结构与索引初稿，人工负责事实校正、判断性补充与边界确认。
- 代码行为变更时，同步更新 `.docs/` 与对应模块或概念文档；已有正文必须基于其当前内容与实现证据增量修订，不得用通用模板整体覆盖。
- 人工补充内容（与实现不符的修正、待裁决项、分层说明）优先于生成初稿；两者冲突时以可追溯的源码、配置与测试证据为准。
- `.specs/domain/` 的语义由人工维护；长期业务或架构约束变化时，以可验证规则的形式手工更新。
- 目录、模块、命令或文档结构变化时，更新 `AGENTS.md` 并同步 `CLAUDE.md`。
- `.specs/tasks/` 与 `.plans/` 是过程产物，不要求长久保留：任务收口后可整体清理，长期语义只沉淀在 `.specs/domain/`。

## SDD Doc Health

- Health status: determined by `task doc:check` for the current checkout
- Checked targets: 31（`.docs` 14 + Domain Spec 4 + 设计意图 1 + AI 入口与索引 10 + 项目入口说明 2）
- Domain specs audited: 4
- Domain constraint IDs: 86（account IA 7 + RA 4、agent IU 9 + RU 13、container IC 7 + RC 2、monitoring I 31 + R 13）
- Task Domain refs: 0
- Orphan CLAUDE entries: 0

- Checked targets:
  - `.docs/doc-system.md`
  - `.docs/api/routes.md`
  - `.docs/api/openapi.yml`
  - `.docs/architecture.md`
  - `.docs/concepts/data-flow.md`
  - `.docs/deployment.md`
  - `.docs/modules/beacon-web.md`
  - `.docs/modules/beacon.md`
  - `.docs/modules/collia.md`
  - `.docs/modules/common-docker.md`
  - `.docs/modules/common-logger.md`
  - `.docs/modules/common-timex.md`
  - `.docs/modules/common.md`
  - `.docs/project-analysis.md`
  - `.specs/domain/account-administration.md`
  - `.specs/domain/agent-lifecycle-update.md`
  - `.specs/domain/container-runtime.md`
  - `.specs/domain/monitoring-platform.md`
  - `DESIGN.md`
  - `AGENTS.md`
  - `CLAUDE.md`
  - `beacon/AGENTS.md`
  - `beacon/CLAUDE.md`
  - `beacon/web/AGENTS.md`
  - `beacon/web/CLAUDE.md`
  - `collia/AGENTS.md`
  - `collia/CLAUDE.md`
  - `common/AGENTS.md`
  - `common/CLAUDE.md`
  - `README.md`
  - `README.en.md`

## 文档门禁的能力边界

统一规则与验证限制见 [文档系统](doc-system.md)。
