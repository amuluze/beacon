# Beacon Web 模块

## 模块职责

`beacon-web` 的主要角色是：frontend experience module: Vue/Vite 页面、路由、API client、Pinia 状态管理和用户交互。

- 承载 Vue 页面、路由、API client、状态管理和用户可见交互。
- 相关源码包：`beacon/web`。

## 职责边界

| 职责 | 说明 | 是否负责 |
|------|------|----------|
| 模块内实现 | frontend experience module: Vue/Vite 页面、路由、API client、Pinia 状态管理和用户交互 | yes |
| 跨模块协调 | 仅通过公开接口或项目约定完成 | conditional |
| 其他模块内部状态 | 不直接读写其他模块私有实现 | no |

## 模块事实

- 前端包：`beacon/web/package.json`
- package name：`beacon-web`
- 前端源码文件：173（`src` 下 `.vue` + `.ts`，含 46 个测试；非测试 127）
- scripts：`build`, `dev`, `lint`, `lint:fix`, `preview`, `test`, `test:coverage`, `test:run`, `ts`
- 前端信号：19 个
- 推断角色：frontend experience module: Vue/Vite 页面、路由、API client、Pinia 状态管理和用户交互

## 接口与符号信号

本节记录源码扫描得到的导出符号或前端入口信号；它们是候选接口线索，不等同于跨模块公开 API。跨模块调用仍以 facade、模块文档和 Domain Spec 约束为准。

- `beacon/web/src/api/index.ts`：axios/API 调用
- `beacon/web/src/api`：API client
- `beacon/web/src/components/Terminal/index.vue`：WebSocket 实时通道
- `beacon/web/src/components`：共享组件
- `beacon/web/src/hooks`：组合式逻辑/hooks
- `beacon/web/src/interface`：前端接口类型
- `beacon/web/src/languages/modules/en.ts`：界面文案（英文语言包）
- `beacon/web/src/languages/modules/zh.ts`：界面文案（中文语言包）
- `beacon/web/src/router/dynamic.ts`：Vue Router 入口
- `beacon/web/src/router/index.ts`：Vue Router 入口
- `beacon/web/src/router/routes.ts`：Vue Router 入口
- `beacon/web/src/router`：路由定义
- `beacon/web/src/store/modules/agent.ts`：Pinia/store 状态
- `beacon/web/src/store/modules/app.ts`：Pinia/store 状态
- `beacon/web/src/store/modules/theme.ts`：Pinia/store 状态
- `beacon/web/src/store/modules/user.ts`：Pinia/store 状态
- `beacon/web/src/store/setup.ts`：Pinia/store 状态
- `beacon/web/src/store`：前端状态管理
- `beacon/web/src/views`：页面视图

## 依赖关系

- `@antfu/eslint-config` `^3.14.0`
- `@element-plus/icons-vue` `^2.3.1`
- `@iconify-json/ep` `^1.2.2`
- `@iconify-json/lucide` `^1.2.116`
- `@types/codemirror` `^5.60.15`
- `@types/lodash-es` `^4.17.12`
- `@types/node` `^20.11.30`
- `@unocss/eslint-plugin` `^65.4.2`
- `@vitejs/plugin-vue` `^5.2.1`
- `@vitest/coverage-v8` `^2.1.9`
- `@vue/test-utils` `^2.4.11`
- `@xterm/addon-fit` `^0.10.0`
- `@xterm/addon-web-links` `^0.11.0`
- `@xterm/xterm` `^5.5.0`
- `axios` `^1.7.9`
- `codemirror` `^5`
- `codemirror-editor-vue3` `^2.8.0`
- `echarts` `^5.5.0`
- `element-plus` `^2.9.3`
- `lodash-es` `^4.17.21`
- `path-browserify` `^1.0.1`
- `pinia` `^2.3.0`
- `pinia-plugin-persistedstate` `^4.2.0`
- `vue` `^3.5.13`
- `vue-i18n` `^11.0.1`
- `vue-router` `^4.5.0`

另有 17 个未列出的 `devDependencies`（构建、测试与 lint 工具链；`devDependencies` 共 27 个，上文已列 10 个）见 `beacon/web/package.json`，不在本文件逐一展开。

## 配置与入口信号

- 前端：`beacon/web/package.json`：依赖清单与 scripts（`build`/`dev`/`lint`/`test`/`ts` 等）
- 前端：`beacon/web/vite.config.ts`：Vite 构建配置（插件、别名、产物目录）
- 前端：`beacon/web/vitest.config.ts`：Vitest 测试配置（环境、setup 文件、覆盖率）
- 前端：`beacon/web/uno.config.ts`：UnoCSS 原子化样式配置
- 前端：`beacon/web/eslint.config.js`：ESLint 配置与自动导入规则
- 前端：`beacon/web/tsconfig.json`：TypeScript 编译配置与路径别名
- 前端：`beacon/web/.eslintrc-auto-import.json`：unplugin-auto-import 生成物，声明自动导入的 API 白名单
- 前端：`beacon/web/web.go`：Go 侧 Web 产物嵌入入口（`go:embed` 打包 `dist`）

## 数据模型与状态

- 本节记录当前实现事实，可追溯到源码与测试断言；长期行为约束以 `.specs/domain/` 为准，本节的「必须 / 禁止」表述不构成约束来源。
- Agent 列表、当前 Agent、加载中和已加载状态由 `store/modules/agent.ts` 统一持有；受保护工作区在 `ensureSelectedAgent()` 完成后才挂载业务内容和状态栏，避免首次进入页面时发出缺少目标 Agent 的数据请求。
- 登录成功、主动退出和认证强制失效都会清空 Agent store，下一认证会话必须重新解析目标 Agent；Agent 相关 API 请求由请求拦截器附加 `X-Agent-ID`。401 刷新队列会把重放结果传回原 Promise，刷新失败时拒绝全部等待请求，避免共享 Agent 加载状态永久 pending。
- axios 始终启用 `withCredentials`，不保存或注入 Authorization Token；JWT 由浏览器 HttpOnly Cookie 承载。首次登录前获取可读 CSRF Cookie，所有 POST/PUT/PATCH/DELETE 请求自动附加 `X-CSRF-Token`。WebSocket 使用同源 Cookie，URL 只携带 Agent 与终端尺寸参数。
- 匿名只读部署（`BEACON_NO_AUTH=true`）把匿名会话视为已浏览会话，允许访问首页与只读页面；路由守卫不把 `/login` 重定向到监控页，顶栏头像菜单提供「登录」入口。登录成功后服务端返回真实用户信息，前端恢复常规权限菜单；主动退出会重新查询匿名只读会话，而不是停留在登录页。
- Pinia 只持久化非敏感用户展示状态；应用启动会清理旧 `sessionStorage` 用户快照并重新查询服务端会话。强制改密用户被路由守卫限制到独立改密页，成功后旧 Cookie 会话失效并返回登录页。
- 容器工作区把容器、镜像与网络拆为页内 Tab（`?tab=containers|images|networks` 深链同步，懒加载：只有当前 Tab 在首次激活时发起查询），三个 Tab 分别只有一个实心主操作，行内主操作跟随容器运行状态（运行中只能停止、已停止只能启动）；标题区提供「刷新」入口并用 `refreshToken` 驱动当前 Tab 列表重新拉取；未选中 Agent 时不渲染任何 Tab，而是按容器 / 镜像 / 网络三段各自渲染空态卡片与刷新入口。
- 受保护容器名单（`beacon`）由 `containerProtection.ts` 维护：前端禁止批量勾选并在名称旁渲染「受保护」标记（行内按钮不禁用）；判定与服务端 `assertNotProtected` 对齐（先 trim 再转小写），因此 `Beacon` / ` beacon ` 这类变体同样受保护；服务端同样强制——`assertNotProtected` 对受保护容器的 start/stop/restart/remove 返回 `ErrProtectedContainer` 拒绝执行（`beacon/service/container/service/container.go`），响应状态为 403 Forbidden，因此保护是服务端约束，界面提示只是其中一层。
- 容器、镜像和网络页面在没有当前 Agent 时禁用运行时操作；确认型命令包装组件首次以 `visible=true` 挂载时立即打开确认框。每次命令调用使用新 vnode，外挂容器在命令关闭或宿主组件卸载时移除。
- 设置工作区提供「设置 / 探针管理 / 用户管理 / 角色管理 / 接口管理 / 审计日志」六个等宽一级 Tab（分段控件，每项 `min-width: 96px`，条宽随 Tab 数量稳定、不随标签字数或选中态抖动；`?tab=config|probe|users|roles|apis|audit` 深链同步，Tab 切换不整页刷新）；`users`、`roles`、`apis` 三个 Tab 涉及权限编排，统一仅 `admin` 可见，`probe` 对所有用户可见（Agent 列表是只读资源，行内升级入口仍限 `admin`，离线探针禁用升级），无权深链回落 `config`，账号权限变化后也会回落到仍可访问的 Tab；新建用户入口在 Tab 工具栏并由票据 prop 驱动列表弹层。设置 Tab 分「告警设置」与「Docker 设置」两个分区，卡片覆盖告警阈值、邮箱通道、群机器人三通道（钉钉 / 飞书 / 企业微信）、告警消息标题与镜像仓库加速；群机器人卡片顶部为三渠道 Tab（等宽 132px，`role="tablist"` + `role="tabpanel"`），一次只渲染当前渠道面板以压缩卡片高度，Tab 行右侧徽标展示当前渠道已保存的启用状态，切换 Tab 不触发保存、各渠道表单的未保存编辑互不影响，面板内展示掩码凭据与本次会话的最近一次测试结果，邮箱通道卡片在卡头展示启用状态并平铺展示配置行。旧深链 `/setting/alarm`、`/setting/container`、`/account/user` 重定向到对应 Tab；`/account/role`、`/account/api` 重定向到 `tab=roles`、`tab=apis`。不再提供 DNS、系统时间、时区、重启或关机等主机系统设置入口与 API client。告警页通过独立编辑组件管理钉钉启用状态、凭据保留/清除、提醒范围和测试发送，页面按字段白名单只接收掩码与配置状态；掩码是只读展示值，不会回填替换输入或作为更新参数发送。
- 容器监控关闭 ECharts 内置图例，只保留卡片级图例；高基数图例使用固定最大高度和纵向滚动，长名称单行省略并通过原始 `title` 提示完整值，CPU 与内存图表仍保留全部容器 series；图例 chip 同时展示容器名与最新采样百分比。
- 控制台视觉以 `.pen/beacon.ax` 设计稿为准：token 集在 `config/theme.ts`（运行时写入 `documentElement` 内联变量）与 `styles/element.scss`（`:root` 与 `.dark` 兜底）两处保持同源；主机 / 容器 / 设置三个工作区共用同一套标题区、图卡与状态栏几何，监控分组的 Agent 上下文条由 `components/Agent/MonitorContextBar.vue` 单点实现（主机与容器分组呈现必须一致）。
- 主机监控按 `freshness` 派生「数据过期」独立态：CPU、内存、磁盘任一 `stale` 时，四张图卡整体让位给 `components/Agent/MonitorStalePanel.vue`（警示条 + 4 张 KPI 卡 + 7 行采样明细），KPI 保留最后采样值并按 `now - collected_at` 实时计算中断时长；容器两行明细只在过期态查询 `/api/v1/container/usage`（取各容器最后采样的平均值 / 合计），正常态不额外请求；正常态四卡（含网络卡）都按同一 `freshness` 显示 `DataStaleTag`。中断时长与时刻格式化的单一来源是 `utils/monitorStale.ts`。
- CPU 指标值按稿成文「当前 X% · N 核」：核数由 collia 用 `cpu.Counts` 采集并随批次上报（`cpu_count`），Beacon 落库 `m_cpu.cpu_count`，`GET /api/v1/host/cpu_info` 以 `cores` 返回；缺失时只渲染「当前 X%」。
- 容器工作区三个 Tab 的表头与工具栏文案以 `.pen/beacon.ax` `Container Center` 为准：容器列 容器名称 / 镜像 / CPU / 内存 / 容器 IP / 端口 / 状态 / 启动时间 / 操作，镜像列 镜像名称 / 容器数量 / Tag / 创建时间 / 镜像大小 / 操作，网络列 网络名称 / 驱动 / 子网 / 网关 / 创建时间 / 操作；筛选的 «全部» 必须用显式值（`all` → 请求时不传该参数），空字符串会被 Element 回落成「请选择」。三个 Tab 的筛选（容器：状态 + 网络；镜像：使用中；网络：驱动）全部在服务端生效、与 `total` 同条件；网络筛选项取当前 Agent 的网络清单，镜像「使用中」数量由上报批次按容器镜像引用计算。
- 监控上下文条（`components/Agent/MonitorContextBar.vue`）按 `stale` prop 切换状态表达：正常为成功色点 +「实时」，过期改为 `accent-warning` 点 +「数据过期」；主机分组传 CPU/内存/磁盘任一的 `freshness.stale`，容器分组传容器用量的 `freshness.stale`。
- 弹层头部统一由 `components/DialogHeader/index.vue` 承担（28×28 `accent-primary-soft` 图标块 + 14/600 标题），通过 `#header` 插槽挂载，关闭按钮沿用 Element Plus 的 header 关闭入口；三处稿面无图标块的弹层（安装 Agent、删除确认、修改密码）只保留默认标题。弹层壳体（`surface-elevated` 底、`border-subtle` 描边、Modal `radius-md` + `0 8px 24px #00000066`、Drawer 左缘贴屏无阴影、头部 16/24、底部上下边框）由 `styles/element.scss` 的 `.el-overlay .el-dialog` / `.el-drawer` 规则单点定义，共享确认弹窗可经 `ConfirmDialog` / `useConfirmCommand` 的 `badge` 参数追加图标块（清理虚悬镜像用）。**Element Plus 样式只允许从 `main.ts` 的全量 `element-plus/dist/index.css` 统一加载**（`vite.config.ts` 的 `ElementPlusResolver` 统一 `importStyle: false`）：组件按需懒加载样式会在路由 chunk 里晚于应用覆盖样式注入，反压弹层 padding / 圆角、抽屉与输入框底色等同特异性规则；`.el-overlay .el-dialog` 前缀则保证壳体规格不被 Element 默认值反压。
- 用户列表以 `ids: string[]` typed props 打开删除确认组件；删除成功才关闭弹窗并刷新列表，请求失败时保留目标与输入状态。告警设置统一使用 `memory` 类型，内存与磁盘编辑组件提交规范类型并仅在成功后关闭、刷新。
- 告警设置页加载并展示全局告警消息标题，通过独立 `EditAlarmMessageTitle` 组件提交裁剪后的单行 1–128 字符标题；成功后关闭并刷新，失败时保留弹窗与本地输入。
- 管理控制台不再注册 Web Terminal 路由、顶层导航或主机监控跳转入口；底层 Terminal 组件及 Server `/ws/terminal` 通道仍保留，但不会由当前页面树挂载。
- 探针管理页在工具栏为管理员提供「自动更新」开关：读取 `GET /api/v1/agent/update_setting`（记录缺失时后端返回 `auto_update=true`），切换即保存 `POST /api/v1/agent/update_setting`，失败回滚并保留成功提示；`deployment_enabled=false` 时开关置灰并提示部署配置已关闭自动调度，非管理员不渲染该开关。
- 顶栏的 circle-arrow-up 按钮表达的是 Beacon 服务端自身的版本更新（`beacon/web/src/layout/navbar/BeaconVersion.vue`）：启动时读取只读端点 `GET /api/v1/system/version`（后端 `beacon/service/system/api/system.go`，返回 `BuildVersion` 与最近一次更新检查状态），无新版本时按钮置灰不可点，检测到新版本时高亮为「有新版本」，点击弹层只读展示当前版本 / 最新版本与手动升级说明；Beacon 不会自升级，是否升级由管理员在服务器上手动执行升级命令决定。Collia 升级入口原在顶栏，现迁入设置工作区「探针管理」Tab 的行内操作列。
- Vite 生产构建会先清空 `dist`，避免历史哈希资源被 `go:embed` 写入 Beacon；Docker 镜像在隔离 Node 阶段构建前端，并通过 `.dockerignore` 拒绝本机 `dist` 进入构建上下文。

## 使用注意事项

- 不跨越模块边界直接依赖内部路径或未导出符号。
- 修改跨模块公开 API、导出符号信号、配置或副作用边界时，同步更新本文件、根入口文档和相关 Domain Spec。
- 配置、凭据和外部副作用只记录语义边界，不记录真实敏感值。

## 验证命令

```bash
cd beacon/web && pnpm run build
cd beacon/web && pnpm run lint
cd beacon/web && pnpm run test
cd beacon/web && pnpm run test:coverage
cd beacon/web && pnpm run test:run
cd beacon/web && pnpm run ts
```

## Taskfile Commands

- `task beacon`
- `task beacon-web`
- `task collia`
- `task default`
