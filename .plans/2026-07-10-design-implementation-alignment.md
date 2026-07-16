# Beacon 设计稿实现对齐实施计划

**目标：** 修复后台工作区可达性与容器编辑链路，并将官网落地到 `.pen` 新版设计。
**需求来源：** `.specs/tasks/design-implementation-alignment.md`
**影响范围：** `common/rpc/schema/`、`collia/service/rpc/`、`beacon/service/`、`beacon/web/src/`、`website/web/src/`

## 技术上下文

- 后台为 Vue 3 + Element Plus + Vite；官网为 Nuxt 4 SSR/预渲染站点，使用自有设计令牌与轻量站内组件。
- 后台新版顶栏仅展示一级路由，旧路由仍拆成 children，造成多个页面无可见入口。
- Server 已具备 `ContainerUpdate` service/repository 契约，但缺 HTTP route/API；Agent RPC 当前明确返回未实现。
- 依赖的 `github.com/amuluze/docker` 把 create restart policy 写死为 `always`，Agent 层需增加本地 runtime adapter，直接调用 Docker SDK 更新策略，并用可替换接口测试重建/回滚。
- 审计列表已支持 Agent 过滤，但 count 未复用同一过滤条件；前端请求拦截器也尚未把 `/api/v1/audit/` 视为 Agent scoped。
- `.pen` 文件为可解析 JSON；本计划不需要数据库迁移，跳过 L1，保留 L2、L3、L4。

## 风险等级说明

| 等级 | 判定标准 |
|------|----------|
| 低 | 样式、路由别名和静态页面调整 |
| 中 | 多组件聚合、状态同步与 API 接入 |
| 高 | Docker 容器替换、跨模块控制契约和视觉验收 |

## Layer 2: 后端实现层

### 任务 2.1: 补齐容器重启策略契约与 HTTP 更新入口
- **依赖**: 无
- **涉及文件**: `common/rpc/schema/container.go`、`beacon/service/schema/container.go`、`beacon/service/container/{api,service}/`、`beacon/service/router.go`
- **风险**: 中 — 共享 RPC 与 HTTP 契约需要保持字段一致
- [ ] 为 create/update 增加 restart policy，并限制 `no/always/unless-stopped/on-failure`
- [ ] 暴露 `/api/v1/container/container_update` 并保持显式 Agent 中间件语义
- [ ] 扩展 service 映射与测试，验证 create/update 字段完整透传
- [ ] 在 container repository/API 链路覆盖缺失与非法 Agent，分别断言 `ErrMissingAgentID`、`ErrInvalidAgentID` 且 RPC caller 未被调用
- **验收标准**:
  - [ ] create/update service 测试断言 restart policy 与显式 Agent 上下文；缺失/非法 Agent 被拒绝
  - [ ] Server container 包测试和构建通过
- **验证**: `go -C beacon test ./service/container/... && go -C beacon build ./...`

### 任务 2.2: 增加可替换 Docker runtime adapter 并实现 Agent 容器编辑
- **依赖**: 2.1
- **涉及文件**: `collia/service/rpc/container_runtime.go`、`collia/service/rpc/container.go`、`collia/service/rpc/container_test.go`
- **风险**: 高 — Docker 不支持原地修改镜像/端口，需安全重建
- [ ] 定义最小 `containerMutationManager` 与 `restartPolicyUpdater` 接口，生产适配器包装现有 Manager 与 Docker SDK
- [ ] RED：覆盖成功、缺失容器、创建失败回滚、策略更新失败、原 stopped/running 状态恢复
- [ ] GREEN：按“旧容器停机/临时改名 → 新容器创建 → 策略更新 → 恢复运行状态 → 删除备份”执行
- [ ] 创建容器后应用显式策略；失败时清理新容器并恢复旧名称与原状态
- [ ] 保持所有错误可见，不返回成功空结果
- **验收标准**:
  - [ ] `ContainerUpdate` 不再返回未实现，回滚测试通过
  - [ ] Agent RPC 包测试与构建通过
- **验证**: `go -C collia test ./service/rpc/... && go -C collia build ./...`

### 任务 2.3: 修复审计筛选与分页 total 一致性 `[可与 2.1 并行]`
- **依赖**: 无
- **涉及文件**: `beacon/service/audit/{repository,service}/`、`beacon/service/testutil/fakes.go`
- **风险**: 中 — 错误 total 会导致设置页分页越界
- [ ] RED：增加 Agent/type 条件下 count 与列表一致的 repository/service 测试
- [ ] 将 `AuditCount` 改为接收 `AuditQueryArgs` 并复用列表过滤条件
- [ ] count 错误必须透传，禁止静默使用零值或全局总数
- **验收标准**:
  - [ ] agent-a / agent-b / 空结果的列表数与 total 一致
- **验证**: `go -C beacon test ./service/audit/...`

## Layer 3: 前端实现层

### 任务 3.1: 聚合后台三个设计工作区并兼容旧深链 `[可与 L2 并行]`
- **依赖**: 无
- **涉及文件**: `beacon/web/src/router/dynamic.ts`、`beacon/web/src/views/{monitor,container,setting}/index.vue`
- **风险**: 中 — 聚合组件需避免重复高度与刷新冲突
- [ ] 新增薄工作区组件，分别组合监控、容器和设置现有页面
- [ ] 将旧 `/monitor/*`、`/container/*`、`/setting/*` 配置为 alias/兼容入口
- [ ] 工作区标题与区域顺序对应 `.pen/beacon.pen`
- [ ] 增加路由结构测试，断言所有旧深链仍可解析
- **验收标准**:
  - [ ] 一级入口可见完整功能且旧 URL 不返回 404
- **验证**: `pnpm --dir beacon/web test:run -- src/router && pnpm --dir beacon/web ts`

### 任务 3.2: 恢复按 Agent 查询的审计区
- **依赖**: 2.3, 3.1
- **涉及文件**: `beacon/web/src/api/audit/`、`beacon/web/src/interface/audit.ts`、`beacon/web/src/views/setting/components/AuditLog.vue`、`beacon/web/src/api/index.ts`
- **风险**: 中 — 必须显式携带 Agent 且分页与切换同步
- [ ] 将 audit prefix 纳入 Agent scoped interceptor，并先加失败测试
- [ ] 恢复审计 API/类型，在 Settings 中增加分页表格
- [ ] 切换 Agent 时重置到第一页并刷新；缺失 Agent 时显示空态
- **验收标准**:
  - [ ] 审计请求注入 `X-Agent-ID`，total 与筛选一致
- **验证**: `pnpm --dir beacon/web test:run -- src/api/index.test.ts && pnpm --dir beacon/web ts`

### 任务 3.3: 实现容器创建/编辑与重启策略
- **依赖**: 2.1, 2.2
- **涉及文件**: `beacon/web/src/api/container/`、`beacon/web/src/interface/container.ts`、`beacon/web/src/views/container/container/`
- **风险**: 中 — 表单需要本地副本与完整请求转换
- [ ] 提取创建/编辑共享的配置类型与序列化函数并先写单测
- [ ] 创建表单增加 restart policy；新增 480px 编辑 Drawer 和操作入口
- [ ] update 请求显式包含容器 ID、完整配置与 restart policy
- [ ] loading/error/success 分支分别可见，取消不修改表格行数据
- **验收标准**:
  - [ ] create/update payload 测试覆盖合法与空数组输入
- **验证**: `pnpm --dir beacon/web test:run -- src/views/container && pnpm --dir beacon/web ts`

### 任务 3.4: 统一弹层规格与修复设置交互 `[可与 3.3 并行]`
- **依赖**: 无
- **涉及文件**: `beacon/web/src/views/container/{image,network}/components/`、`beacon/web/src/views/setting/{alarm,docker}/components/`
- **风险**: 中 — 涉及多个独立交互组件
- [ ] Pull/Import Image、New Network、CPU/Mem/Disk Threshold、Email 改为 480px Modal
- [ ] Registry Mirror 保持 540px Drawer；Add/Edit Container 由任务 3.3 统一为 480px Drawer
- [ ] 邮件表单从 props 深拷贝本地状态，保存成功后再通知父级
- [ ] 阈值/邮件/网络请求只在成功分支提示成功，失败保留弹层与输入
- [ ] 增加邮件取消不污染 props、阈值失败不提示成功的组件测试
- **验收标准**:
  - [ ] 每类弹层宽度映射与 DAI-05 一致，状态反馈测试通过
- **验证**: `pnpm --dir beacon/web test:run -- src/views/setting src/views/container && pnpm --dir beacon/web ts`

### 任务 3.5: 建立官网设计令牌与 Header/Footer Shell `[可与后台并行]`
- **依赖**: 无
- **涉及文件**: `website/web/src/styles/`、`website/web/src/components/{Header,Footer}/`、`website/web/src/layouts/default.vue`
- **风险**: 中 — 全站视觉基线变化
- [ ] 建立与 `.pen/beacon-website.pen` 对应的 surface/foreground/border/accent/spacing 字段
- [ ] Header 增加首页、使用手册、GitHub；Footer 增加团队故事、使用手册、微信公众号、GitHub
- [ ] 移除紫蓝渐变、渐变文本和旧玻璃发光 mixin 的页面依赖
- [ ] 布局改用标准 slot，保持固定导航与内容顶部间距一致
- **验收标准**:
  - [ ] Shell 关键链接与设计一致，源码扫描无旧主视觉反模式
- **验证**: `task website-web:build`（从仓库根目录执行，任务会先安装依赖）

### 任务 3.6: 重做官网首页与统计降级
- **依赖**: 3.5
- **涉及文件**: `website/web/src/pages/index.vue`、可复用首页子组件目录
- **风险**: 中 — 首页结构整体替换且依赖统计接口
- [ ] 按设计实现 Hero、安装命令、累计获取、容器/主机/用户功能区、技术栈和 CTA
- [ ] 统计接口失败时捕获错误并展示 `--`，主体继续渲染
- [ ] 删除旧 Carousel 与科技蓝渐变模板
- [ ] 用 375px 断点折叠为设计稿精简卡片顺序
- **验收标准**:
  - [ ] 首页关键区块与文案存在，统计失败测试/逻辑不抛未捕获异常
- **验证**: `task website-web:build`（从仓库根目录执行，任务会先安装依赖）

### 任务 3.7: 重做关于、更新日志与文档页面 `[可与 3.6 并行]`
- **依赖**: 3.5
- **涉及文件**: `website/web/src/pages/about.vue`、`website/web/src/pages/changelog.vue`、`website/web/src/pages/document.vue`
- **风险**: 中 — 三个独立内容页面需共享排版系统
- [ ] About 增加团队信件、联系卡片与公众号区域
- [ ] 新增 Changelog 并实现 v3.0.4/v3.0.0/v2.0.0 时间线
- [ ] Document 使用目录、安装、常见问题与技术支持结构，修正旧 `Beacon 2.0` 标题
- [ ] 三页共享 section/card/code 样式且兼容移动端
- **验收标准**:
  - [ ] `/about`、`/changelog`、`/document` 文件路由与关键内容完整
- **验证**: `task website-web:build`（从仓库根目录执行，任务会先安装依赖）

### 任务 3.8: 官网 375px 响应式收口
- **依赖**: 3.5, 3.6, 3.7
- **涉及文件**: 官网 Header/Footer 与各页面 scoped styles
- **风险**: 中 — 移动端需要跨页面检查
- [ ] 导航折叠、Hero 字阶、卡片单列、代码横向滚动、Footer 单列
- [ ] 添加结构性 CSS 检查，禁止固定 1440px 内容宽度和页面水平溢出
- **验收标准**:
  - [ ] 375px 截图无横向溢出，区块顺序对应 Landing Mobile
- **验证**: `task website-web:build` + 375px 浏览器截图

## Layer 4: 集成验证层

### 任务 4.1: 后端与前端编译/相关测试汇总
- **依赖**: 2.1, 2.2, 2.3, 3.1-3.8
- **涉及文件**: 全部本次变更
- **风险**: 高 — 涉及 Docker 替换链路、两个前端和共享契约
- [ ] `common`、`collia`、`beacon` 分 module 编译与相关测试
- [ ] 后台类型检查、Vitest 与生产构建
- [ ] 官网安装依赖后执行 Nuxt build
- [ ] 执行 `git diff --check`、`git status`、`git diff --stat` 与完整 diff 审查
- **验收标准**:
  - [ ] 所有可用编译与相关测试通过，无敏感信息与非预期文件改动
- **验证**: `go -C common test ./... && go -C collia test ./service/rpc/... && go -C beacon test ./service/container/... ./service/audit/... && pnpm --dir beacon/web test:run && pnpm --dir beacon/web build && task website-web:build`

### 任务 4.2: 桌面与移动视觉对照
- **依赖**: 4.1
- **涉及文件**: `.pen/beacon.pen`、`.pen/beacon-website.pen` 与运行页面
- **风险**: 高 — 视觉验收依赖可用浏览器环境
- [ ] 1440px 截图对照后台 Monitor/Container/Settings 与官网 Landing/About/Changelog/Docs
- [ ] 375px 截图对照 Landing Mobile 并检查无水平溢出
- [ ] 浏览器拦截官网统计查询返回 HTTP 500，断言累计获取显示 `--` 且 Hero、功能区与 CTA 仍存在
- [ ] 检查 console error 与关键交互状态
- **验收标准**:
  - [ ] 视觉结构、颜色、间距与响应式达到设计验收；若浏览器不可用则明确标记阻塞，不声称完成视觉验收
- **验证**: 可用浏览器中的桌面/移动截图与 console 记录
