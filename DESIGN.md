# Beacon 设计意图文档

> 本文是 Beacon **设计意图的 SSOT**。可视化载体是一份 Axiom 设计稿（格式 `ax@1.6`）：控制台稿 `.pen/beacon.ax` — 41 页。
>
> 该设计稿维护在本仓库之外，不随本仓库分发；本文以纯文本引用稿名与其页数，稿内字面量仍以设计源为准。
>
> 稿件承载「长什么样」，本文承载「为什么这样」「哪些不可变」。实现侧（`beacon/web`）以本文与稿件为准绳，但实现代码本身是「是否已落地」的事实源。

---

## 1. 文档定位与事实边界

| 载体 | 权威范围 | 不负责 |
|------|----------|--------|
| `.pen/beacon.ax`（在本仓库之外维护） | 版式、间距、色值字面量、文案原文、组件几何 | 意图解释、状态机、错误语义 |
| `DESIGN.md`（本文） | 设计原则、布局与组件规范、状态模型、文案语气、无障碍与响应式要求 | 具体像素字面量（引用 token 即可）；与实现的差异见 §1「维护约定」 |
| `.specs/domain/*.md` | 长期行为约束与不变量 | 视觉 |
| `beacon/web` | 是否已实现、路由与接口事实 | 规范来源 |

**冲突裁决顺序**：Domain Spec → 本文 → 设计稿 → 实现代码。出现分歧时先按此顺序定位「谁需要改」，再改，不要让实现反向定义设计。

**维护约定**

- 改稿后必须回写本文对应章节；本文新增规则必须在稿中有对应落点。
- 色值只允许来自稿中 token；实现里禁止写字面色值。
- 本文描述「应然」；稿件与实现之间的差异记在各模块文档（`.docs/modules/`）与 Domain Specs（`.specs/domain/`）中，不要散落在各章节。

---

## 2. 产品面

Beacon 是 **Server-Agent 架构**的主机与 Docker 容器监控平台，由 Server、Agent 与管理控制台三部分组成：Server 汇聚数据并提供 API，Agent（`collia`）采集主机与容器数据并执行 Docker 操作，用户只在**管理控制台**这一个产品面上完成全部操作。

| 产品面 | 实现位置 | 用户 | 设计目标 |
|--------|----------|------|----------|
| 管理控制台 | `beacon/web`（Vue 3 + Element Plus + Pinia） | 已认证的运维 / 管理员 | 密集、可信、可操作：让运维在几秒内判断「这台机器现在是否健康」，并在 Agent 上直接执行 Docker 操作 |

控制台是本产品唯一的界面产品面，其视觉语言定位为密集的技术工具：密度优先、装饰克制、把判断成本压到最低。品牌内核（开源、轻量、Server-Agent、MIT License）在产品各环节保持一致。

---

## 3. 页面地图

### 3.1 控制台（`.pen/beacon.ax`，41 页）

稿件以 **Group 索引页** 组织：每个 Group 页是「分组标题 + 描述 + 子页清单」的目录页，本身不是产品页面。

| Group | 页 | 覆盖的界面 |
|-------|----|-----------|
| 监控 | `Group · 监控` → `Monitor`、`Monitor · 空状态`、`Monitor · 数据过期` | 主机 4 图 + 容器 2 图的趋势工作区；无 Agent 引导；上报中断的过期态 |
| 容器 | `Group · 容器` → `Container Center`、`Container · 无 Agent 空态` | 容器 / 镜像 / 网络三 Tab 工作区；未接入 Agent 时三个 Tab 共用一个空态卡片 |
| 设置 | `Group · 设置` → `Settings Page`、`Settings · 探针管理`、`Settings · 用户管理`、`Settings · 角色管理`、`Settings · 接口管理`、`Settings · 审计日志` | 告警阈值与通道、Docker 镜像仓库、已安装探针清单、账户与角色、接口清单、审计日志 |
| 认证与异常态 | `Group · 认证与异常态` → `Login`、`Change Password`、`404` | 登录、强制改密、兜底错误页 |
| 弹层 | `Group · 弹层 · Modal / Drawer` → 18 个弹层页 | 安装/升级 Agent、容器增删改、镜像拉取/导入/清理、网络新建、镜像仓库、告警邮箱/阈值/标题、用户增改、改密、新建角色 |
| 组件基线 | `Group · 组件基线` → `Components Showcase`、`Navbar · 顶部栏下拉` | 组件几何基线；顶栏展开态 |

---

## 4. 视觉基调

**控制台**：冷静、克制、密度优先。主色是低饱和青绿（`accent-primary`），界面几乎无装饰性图形，靠 1px 边框和表面层级差分区。留白集中在区块之间而非区块内部，不使用装饰性渐变。**唯一例外是登录页左栏的品牌面板**：它承载品牌渐变色块与两处装饰圆形（`accent-primary` / `accent-primary-bright`），文案为浅色（见 §10.1）。

---

## 5. Token 体系

### 5.0 品牌色族与设计系统边界

控制台是本产品唯一的界面设计系统，token 命名与取值都在同一套体系内维护：

- **唯一色相族**：主色落在绿—青色相（`accent-primary` `#0A7C79` 偏青），是品牌色的唯一来源，禁止再引入第二种品牌绿。
- **中性色服务于信息密度**：控制台的 `surface-*/foreground-*/border-*` 面向数据密集的工作面，不做跨场景的浅色化改写。

命名体系统一为 `<族>-<语义>`（`surface-primary`、`accent-primary`），不引入扁平语义名。

| 角色 | 控制台 | 说明 |
|------|--------|------|
| 主色 | `accent-primary` `#0A7C79` / `#0FC6C2` | 品牌主色，明暗两档 |
| 主色 hover | `accent-primary-hover` `#0C918D` / `#3DD2CF` | 主色朝更暗/更亮一档推进 |
| 主色浅底 | `accent-primary-soft` `#E4F6F6` / `#0F3A39` | 用于选中行、标签底色 |

**预留 token（显式声明，非死 token）**：控制台当前没有保留中的预留 token——`brand-deep` 与 `accent-primary-bright` 已由登录页左栏品牌面板消费（渐变深端与装饰圆形，见 §10.1），不再是预留位。稿中新增 token 若实现尚未引用，须在此显式登记为预留位；实现侧不得据预留 token 反推新增样式，若要清理，须先在本节登记删除并同步稿中 tokens 段。`accent-primary-hover` 与 `accent-primary-soft` 已按 §12.3 落地到按钮 / 导航项 / 表格行 / 链接的 hover 全局规则（`styles/element.scss`）；`surface-elevated`、`surface-inset`、`foreground-placeholder`、`accent-success-soft` / `accent-warning-soft` / `accent-danger-soft`、`chart-1`–`chart-6`、`font-sans` / `font-mono` 已随稿面 token 集在实现侧定义并消费（`config/theme.ts` + `styles/element.scss`）。

实现侧切换品牌色时必须同步 `beacon/web/src/config/theme.ts`（明暗两档，含 Element Plus 的 `--el-color-primary-light-3/5/7/8/9` 与 `--el-color-primary-dark-2` 派生色）与 `beacon/web/src/styles/element.scss` 的 `:root` 默认值；只改 `--el-color-primary` 会让 hover 与浅底仍停留在旧色。

### 5.1 控制台 token（48 个，取值来自控制台设计稿 `.pen/beacon.ax`；该稿维护在本仓库之外，不随本仓库分发）

命名规则：`<族>-<语义>`，族为 `surface / foreground / border / accent / chart / font / spacing / radius / brand / overlay`。

**表面与前景**

| Token | light | dark | 语义 |
|-------|-------|------|------|
| `surface-primary` | `#FFFFFF` | `#0D0D0E` | 页面底色、弹层外的卡片区 |
| `surface-elevated` | `#FFFFFF` | `#131316` | 顶栏、Modal/Drawer 壳体 |
| `surface-card` | `#F8FAFC` | `#131316` | 卡片、表格容器、图标按钮底 |
| `surface-secondary` | `#F1F5F9` | `#1C1C20` | 表头、状态栏底、分组容器 |
| `surface-inset` | `#FFFFFF` | `#1C1C20` | 输入框内陷面、图区内衬、chip 底 |
| `foreground-primary` | `#0F172A` | `#F8FAFC` | 标题与主文本 |
| `foreground-secondary` | `#475569` | `#CBD5E1` | 次级说明、表头文字 |
| `foreground-muted` | `#7C8CA3` | `#94A3B8` | 弱化说明、轴标签、分隔符 `·` |
| `foreground-placeholder` | `#94A3B8` | `#7C8CA3` | 占位符、禁用按钮文字 |
| `foreground-on-accent` | `#FFFFFF` | `#0F172A` | 主色实心面上的文字 |

**边框与遮罩**

| Token | light | dark | 语义 |
|-------|-------|------|------|
| `border-primary` | `#E2E8F0` | `#27272A` | 控件描边（输入框、按钮、表格线） |
| `border-subtle` | `#EEF2F6` | `#1E293B` | 卡片分隔、区块边界 |
| `overlay-scrim` | `#0F172A66` | `#000000B3` | Modal/Drawer 遮罩 |

**状态色（每色三档：实心 / soft 底 / 弱化文字）**

| 语义 | 实心 | soft 底 | 用途 |
|------|------|---------|------|
| 主色 | `accent-primary` `#0A7C79` / `#0FC6C2` | `accent-primary-soft` `#E4F6F6` / `#0F3A39` | 主按钮、激活导航、eyebrow、CPU 指标 |
| 成功 | `accent-success` `#047857` / `#34D399` | `accent-success-soft` `#ECFDF5` / `#052E24` | 运行中、已启用、正常状态 |
| 警告 | `accent-warning` `#B45309` / `#FBBF24` | `accent-warning-soft` `#FFFBEB` / `#3A2E10` | 数据过期、内存越线、虚悬镜像 |
| 危险 | `accent-danger` `#DC2626` / `#EF4444` | `accent-danger-soft` `#FEF2F2` / `#3B1416` | 已停止、删除操作、校验失败 |
| 品牌深色 | `brand-deep` `#0F3A39` / `#0B3B3A` | — | 登录页左栏品牌面板的渐变深端 |

另有 `accent-primary-hover`（`#0C918D` / `#3DD2CF`）与 `accent-primary-bright`（`#7FC9C6` / `#5FE0DD`）用于 hover 与强强调；`accent-primary-hover` 已按 §12.3 落地到实现的全局 hover 规则，`accent-primary-bright` 用于登录页左栏的下左装饰圆形（§10.1）。

**图表序列色**

| Token | light | dark | 绑定指标 |
|-------|-------|------|---------|
| `chart-1` | `#0A7C79` | `#0FC6C2` | CPU（与主色同源） |
| `chart-2` | `#4FB8B4` | `#5FE0DD` | 磁盘 IO 备用 |
| `chart-3` | `#7C5CE0` | `#A78BFA` | 网络流量（紫，唯一与主色不同族） |
| `chart-4` | `#B45309` | `#FBBF24` | 告警线 / 阈值 |
| `chart-5` | `#047857` | `#34D399` | 成功态曲线 |
| `chart-6` | `#E11D48` | `#FB7185` | 危险态曲线 |
| `chart-3-soft` | `#F1ECFE` | `#2B2140` | 网络类卡片的图标底 |

**排版、间距与圆角**

| Token | 取值 | 用途 |
|-------|------|------|
| `font-sans` | Inter | 界面正文字体 |
| `font-mono` | Geist Mono | 指标值、ID、时间戳、镜像 Tag、大小、版本号 |
| `font-size-xs` | 11px | 时间戳、轴标签、辅助说明 |
| `font-size-sm` | 13px | 控件正文、按钮与表单标签 |
| `font-size-md` | 14px | 正文、表头文字 |
| `font-size-lg` | 18px | 卡片标题 |
| `font-size-xl` | 24px | 区块标题 |
| `font-size-2xl` | 32px | 工作区主标题与 KPI 大数字 |
| `spacing-xs` | 4px | 图标与文字、标签内间距 |
| `spacing-sm` | 8px | 控件内边距、标签间距 |
| `spacing-md` | 16px | 卡片内边距、弹层主体内边距 |
| `spacing-lg` | 24px | 区块间距、弹层头部与底部内边距 |
| `spacing-xl` | 32px | 页面主区内边距 |
| `radius-xs` | 4px | 标签与徽标 |
| `radius-sm` | 6px | 输入框、按钮 |
| `radius-md` | 8px | 卡片、弹层壳体 |
| `radius-lg` | 10px | 登录卡、预览卡 |

### 5.2 语义使用规则

1. **主色只做「可点」和「当前」**：主按钮、激活态导航、当前 Tab、主色实心图标块。不用主色做大面积背景。
2. **状态色三件套**：`实心`（文字/图标/描边）+ `soft`（标签底、图标底）+ `muted`（弱化说明）。同一状态在控制台内必须落到同一档。
3. **数字一律等宽**：所有指标值、ID、IP、端口、时间戳、镜像 Tag、大小、版本号使用 `font-mono`。这是控制台最硬的一条排版规则（稿中 `font-mono` 出现 306 次，仅次于 `font-sans`）。
4. **字号承载层级，不承载情绪**：`2xl` 只用于页面主标题与 KPI 大数字；`xl` 用于区块标题；`md` 用于卡片标题与表头；`sm` 是正文与控件文字的默认档；`xs` 用于时间戳、辅助说明、标签。
5. **字重分工**：`700` = 页面主标题、KPI 数字、eyebrow；`600` = 卡片标题、导航激活文字、主按钮文字、状态标签；`500` = 次级强调（按钮、表格主字段）；`normal` = 正文与说明。
6. **弱化分隔符**：同一行内的元信息用 `·`（`foreground-muted`）分隔，而不是竖线或空格。
7. **等宽字体不用于中文文案**，只用于数据。

### 5.3 Token 治理规则

- 新增 token 必须先在稿中定义、再在本文登记、最后在实现中消费；顺序不可颠倒。
- 同一语义不得出现两个 token（如 `surface-elevated` 与 `surface-primary` 在 light 下同值是允许的，但语义必须不同）。

---

## 6. 控制台布局系统

所有工作区页共用同一套骨架：**顶栏 / 内容区 / 底部状态栏**。画布基准宽 1440px。

### 6.1 顶栏（高 56）

从左到右固定顺序：

| 序 | 元素 | 规格 |
|----|------|------|
| 1 | 品牌 | 28×28 品牌图片（`object-fit: contain`，稿件资源 `beacon.svg`，随设计稿在本仓库之外）+ 「Beacon」`font-size-lg` 700 |
| 2 | 主导航 | `BeaconNavItem` 三项：Monitor(activity) / Container(package) / Settings(settings) |
| 3 | Agent 选择器 | 220×32，`surface-inset` + `border-primary` + `radius-sm`；左 server 图标 + hostname + chevron-down |
| 4 | Beacon 版本提示 | circle-arrow-up 图标常驻；无新版本时置灰不可点，检测到新版本时用 `accent-primary-soft` 底 + `accent-primary` 图标与文字「有新版本」，点击弹层只读展示当前版本 / 最新版本与手动升级说明 |
| 5 | 安装 Agent | primary 按钮（`accent-primary`），download 图标 |
| 6 | 语言切换 | 32×32 图标按钮（languages） |
| 7 | 主题切换 | 32×32 图标按钮（moon） |
| 8 | 账户 | 30×30 主色圆形头像 + 用户名 |

顶栏底色 `surface-elevated`，仅底部 1px `border-subtle`，左右内边距 24（`$spacing-lg`）。

**导航状态**：`BeaconNavItemActive` = 主色实心胶囊 + `foreground-on-accent` 文字（600）；`BeaconNavItem` = 无底色 + `foreground-muted` 文字（500）。图标 15px，间距 6，内边距 8/12，`radius-sm`。

### 6.2 标题区

四个区块：eyebrow → 主标题 → 副标题 → 右侧元信息。

- **eyebrow**：`font-size-xs` 700、`accent-primary`、`letterSpacing: 2`、全大写英文。每个工作区一个：监控 `OBSERVABILITY`、容器 `RUNTIME`、设置 `CONFIGURATION`、顶栏专项 `NAVIGATION`。这是控制台唯一的英文装饰性元素，用于区分工作区。
- **主标题**：`font-size-2xl` 700。
- **副标题**：`font-size-sm`、`foreground-muted`、`lineHeight` 1.55、限宽 620。
- **右侧元信息**：mono 小字标签，`surface-inset` 底 + `border-primary` + `radius-sm`，内边距 4/10。

### 6.3 底部状态栏

`surface-secondary` 底 + 顶部 1px `border-primary`，内边距 12/24。

- 左组（4 项）：主机 / 系统 / 内核 / OS
- 中位独立项（1 项）：运行时间（左右两组之间居中，不归属任何一组）
- 右组（4 项）：Docker / API / OS / Beacon 版本

两个「OS」是有意设计：左组 OS 是主机系统（`os/kernel_arch`），右组 OS 是 Docker Engine 平台（`docker os/arch`），数据源不同。

右组末项是 Beacon 服务端版本：`label` 固定「Beacon」，`value` 为服务端构建版本（如 `v3.0.15`，取自 `beacon/VERSION`），构建期未注入版本时显示 `dev`。该项固定在状态栏最右端，与同组其他项一样以「·」分隔，不与主机侧信息混排。

每项为 `BeaconStatusItem`：`label` 11px `foreground-secondary` + `value` 11px 500 `font-mono`。

顶栏与状态栏的底色分层：顶栏 `surface-elevated` + 1px `border-subtle`（浮于内容面之上），状态栏 `surface-secondary` + 1px `border-primary`（沉于内容面之下）；右侧末项 Beacon 版本项的实现落点是 `beacon/web/src/layout/statusbar/index.vue`。

### 6.4 内容区

内容区左右内边距 `$spacing-lg`（24），区块间距 `$spacing-lg`。页面结构固定为「标题区 → 工具区 → 主体区 → 分页/尾部」。

### 6.5 只读模式提示条

部署处于只读模式（`read_only`，含匿名只读会话）时，顶栏之下、内容区之上插入一条常驻提示条：`accent-warning-soft` 底 + 1px `accent-warning` 下边线，内边距 6/24，居中排布，12px 锁图标 + `font-size-xs` `accent-warning` 文案「当前部署为只读模式：可浏览页面内容，所有修改操作已被禁用。」（`common.readOnlyBanner`）。

提示条不改变导航与状态栏，也不隐藏页面：只读态下页面结构完整，写操作按钮由实现置灰。稿中 `Monitor` 页是该全局条的落点样本；实现落点为 `beacon/web/src/layout/index.vue`。

---

## 7. 监控工作区

**意图**：一页回答「这台主机现在健康吗」。主机与容器两组趋势并列，采样口径、刷新节奏、Agent 选择器共用同一套上下文。

### 7.1 主机监控（4 张图卡）

| 卡片 | 标题 | 指标值示例 | 语义色 |
|------|------|-----------|--------|
| 1 | CPU 使用率 | `当前 23.5% · 12 核` | `accent-primary` |
| 2 | 内存使用率 | `3.2 GiB / 7.8 GiB · 72.1%` | `accent-warning` |
| 3 | 磁盘 IO | `读 5.6 / 写 2.1 MB/s · 最高 68.4%` | `accent-success` |
| 4 | 网络流量 | `接收 1.2 MB/s · 发送 3.4 MB/s` | `chart-3` |

**图卡几何**：`surface-card` 底 + `border-subtle` + `radius-md` + 内边距 24，高 260。标题行为 28×28 soft 色图标块 + 标题（14/600）+ 指标值（mono 13/600，语义色），左对齐紧邻（间距 8）。图区内衬 `surface-inset` + `border-subtle` + `radius-sm`，高 161：顶部图例行高 18（线段 + 圆点 + 序列名，13px `foreground-secondary`，右对齐；四卡都有，CPU / 内存各 1 项，磁盘「读 / 写」、网络「接收 / 发送」各 2 项），绘图区高 82，底部轴标签行高 15（11px mono muted，与绘图区间距 6）。

**图内构成**：左列纵轴刻度（右对齐，11px mono muted；百分比卡 6 档 0–100%，速率卡 6 / 8 档并带单位后缀）→ 绘图区（横向网格线 1px `border-primary`，条数与刻度一致）→ 底部时间轴标签（等距，10 个采样点）。曲线为单条平滑折线（单调三次插值、无数据点符号），线色取该卡语义色（CPU `chart-1` / 内存 `accent-warning` / 磁盘 `chart-1` 读 + `chart-5` 写 / 网络 `chart-1` 收 + `chart-5` 发）；实现侧对应 ECharts `smooth: true` + `showSymbol: false` + `lineStyle.width: 2` 的折线，磁盘 / 网络为双序列叠加。

**图例与序列对应**：图例项与折线一一对应（图例名 = 序列名、图例色 = 序列色），序列色一律从 `chart-1`–`chart-6` token 解析、不写死色值；多磁盘 / 多网卡时序列名补设备标识（如 `disk0 读`），单设备时用「读 / 写」「接收 / 发送」。图例文字色取 `foreground-secondary`，暗色主题下同样可读。

**指标值口径**：指标值按「当前值 + 上下文」成文（CPU `当前 X% · N 核`、内存 `used / total · percent`、磁盘 `读 r / 写 w · 最高 p%`、网络 `接收 r · 发送 w`）。CPU 核数由 collia 采集（`cpu.Counts`）随批次上报，落库 `m_cpu.cpu_count`，`GET /api/v1/host/cpu_info` 以 `cores` 返回；核数缺失（历史数据或采集失败）时只渲染 `当前 X%`，不补占位。

### 7.2 容器监控（2 张图卡）

容器 CPU 使用率 / 内存使用率。与主机卡同外框（`surface-card` 底 + `border-subtle` + `radius-md` + 内边距 24），高 372。

**图卡几何**：标题行为 28×28 soft 色图标块 + 标题（14/600；CPU `accent-primary` / 内存 `accent-warning`），左对齐紧邻（间距 8），无副标题行。标题行下方是图例 chip 区：每 chip = 7×7 序列色圆点 + 容器名（11px `foreground-secondary`）+ 百分比（11px mono `foreground-secondary`），4 列换行排列，列距 16 / 行距 4，超出 72 高纵向滚动。图表区内衬 `surface-inset` + `border-subtle` + `radius-sm`，高 204：左上角图内名（14/600，「CPU」/「内存」），左侧纵轴刻度（11px mono muted，右对齐；CPU 7 档 1.8–0，内存 7 档 171.66 MB–0.00 B），绘图区高 104（横向网格线 1px `border-primary`，条数与刻度一致），底部时间轴行高 15（11px mono muted，10 个采样点）。

**序列表达**：每个容器一条平滑折线（单调三次插值、无数据点符号），序列色按 `chart-1`–`chart-6` 循环取用（容器数超过 6 时复用色板），实现侧对应 ECharts `containerCpuOption` / `containerMemOption` 的 `smooth: true` + `showSymbol: false` + `lineStyle.width: 1.5`。图例 chip 的圆点色与折线色一一对应；ECharts 内置图例关闭（`legend.show: false`），图例职责由卡内 chip 区承担。

**图例百分比口径**：chip 上的百分比必须是 0–100 的使用率（CPU 取容器 CPU 百分比，内存取容器内存百分比），不得直接渲染字节数等非百分比原始值。实现落点：内存百分比取 `ContainerUsageReply.mem_percent`（0–100），纵轴与折线仍用 `mem_usage` 的字节口径。

### 7.3 Agent 上下文条

每个分组（主机监控 / 容器监控）顶部一行：`Agent 切换器` + 绿色状态点 + 「实时」 + `·` + 「最后更新 10:24:08」+ `·` + 「时间密度：」+ 「10分钟」下拉。**状态、时间戳、密度三者的呈现方式必须在两分组间完全一致**。分组数据过期时（`freshness.stale`）状态点与文案改用 `accent-warning`「数据过期」，不得继续显示「实时」（domain I002）。

### 7.4 空状态（`Monitor · 空状态`）

标题「暂无可用 Agent」+ 说明「在本机执行安装命令接入 collia Agent，连接成功后会自动展示当前节点的监控数据。」+ 两个主按钮「安装 Agent」「刷新」。

规则：**空态不展示任何指标骨架**，不给用户「有数据但没画出来」的错觉；必须给出可执行的接入入口 + 刷新入口。

### 7.5 数据过期（`Monitor · 数据过期`）

Agent 上报中断后的**独立态**，不是错误态：趋势图让位给「警示条 + KPI 大数字 + 采样明细」，KPI 保留最后一次采样值。

- **判定与取值**：用各指标响应里的 `freshness`（`collected_at` / `age_seconds` / `stale`）判定——CPU、内存、磁盘任一 `stale` 即进入过期态；「最后一次采样」取各指标 `collected_at` 的最大值，中断时长按 `now - collected_at` 实时计算（<60 秒按秒、<60 分钟按分钟、更长按「N 小时 M 分钟」）。
- 顶部警示条：`accent-warning-soft` 底 + `accent-warning` 文案「Agent 上报已中断 12 分钟」+ `·` + 「最后一次采样 2026-09-14 10:24:08」+ muted「数据恢复后自动续传，无需手动刷新」。
- KPI 卡（4 张，`surface-card` + `border-subtle` + `radius-md` + 内边距 24）：标题（13/600）+「数据过期」徽标（`accent-warning` + `accent-warning-soft` + `radius-xs`）+ 右对齐「最后采样 10:24:08」→ 2xl 大数字 → 11px muted 脚注。脚注按卡成文：CPU「N 核 · 最近一次采样来自 HH:mm:ss」、内存「used / total · 采样停在上报中断前」、磁盘「device · 数据过期期间不再判断告警阈值」、网络「eth · 收/发采样均停在上报中断前」。
- 「采样明细」表（列：指标 / 最后采样值 150 / 过期时长 120 / 状态 180；表头 `surface-secondary`，行 10/16 + 底部 `border-subtle`）：只列平台真实采集的 **7 项**指标——CPU 使用率、内存使用率、磁盘 IO、网络流量、容器 CPU 使用率（各容器最后采样的平均值）、容器内存使用率（各容器合计）、磁盘使用率；**不列未采集的主机负载**。过期时长用 `accent-warning`，状态固定「等待 Agent 续传」。
- **过期期间不判断告警阈值**，恢复上报后趋势线自动续传（正常态下四张图卡的 `数据过期` 标签也按同一 `freshness` 显示）。

---

## 8. 容器工作区

**意图**：以当前 Agent 的 Docker Engine 为唯一操作目标，完成容器 / 镜像 / 网络三类资源的管理。

### 8.1 页面结构

标题区（eyebrow `RUNTIME`）→ 右侧 3 个 mono 统计标签（容器 7 · 运行中 5 / 镜像 6 · 4.2 GB / 网络 6）+「刷新」→ 三 Tab → 工具栏 → 表格 → 分页。

Tab 变体：激活项 = `accent-primary` 600 文字，右侧 mono 计数；未激活 = `foreground-secondary` 500。

### 8.2 工具栏

左：主操作按钮（新建容器 / 拉取镜像，primary，h32，`radius-sm`）+ 搜索框（260×32，`surface-inset`）+ 筛选下拉（容器 Tab：状态 / 网络；镜像 Tab：使用中；网络 Tab：驱动）+「已选 N 项」+「批量操作」（无选中时为 `foreground-placeholder` 禁用态）。筛选值 «全部» 必须用显式值渲染（空字符串会被 Element 回落到「请选择」占位符），筛选在服务端生效、`total` 与列表同条件。

### 8.3 表格规格

容器 `surface-card` + `radius-md`；表头 `surface-secondary` 高 44；数据行内边距 10/16、行间 `border-subtle`；分页条高 48、`surface-card` 底、顶部 `border-subtle`。

**容器表列宽**

| 列 | 宽 |
|----|----|
| 选择框 | 44 |
| 容器名称 | `fill_container` |
| 镜像 | 180 |
| CPU | 80 |
| 内存 | 96 |
| 容器 IP | 120 |
| 端口 | 120 |
| 状态 | 96 |
| 启动时间 | 116 |
| 操作 | 180 |

**镜像表列宽**：镜像名称 `fill_container` / 容器数量 100 / Tag 140 / 创建时间 200 / 镜像大小 120 / 操作 120。

### 8.4 状态与保护

- 状态标签：`运行中` = `accent-success`，`已停止` = `accent-danger`，Tag 值 = `accent-primary`，`受保护` = `foreground-muted`。
- **受保护容器**（容器名 `beacon`）的启停与删除按钮在列表中禁用。删除镜像时，容器数量 > 0 的条目同样禁用删除。
- 操作列按状态给出不同动词：运行中给「停止」，已停止给「启动」，另有「更多」。

### 8.5 无 Agent 空态

三段各自渲染空态卡片与刷新入口，不做整页拦截——保持结构可预期，用户能看到「有哪些资源类别」。

---

## 9. 设置工作区

**意图**：一处收敛「告警怎么发、镜像怎么拉、谁能做什么、做过什么」。

### 9.1 分区

一级 Tab：设置（`config`）/ 探针管理（`probe`）/ 用户管理 / 角色管理 / 接口管理 / 审计日志；`config` Tab 内分**告警设置**与 **Docker 设置**两块。探针管理对所有用户可见（Agent 列表本身是只读资源），行内升级入口仍限管理员。

**稿页 ↔ 实现映射**：稿把设置拆成 6 个独立页（9–14 页），实现是 `/setting` 一个路由 + query Tab 切换。二者的对应关系固定如下，改稿或加 Tab 时必须同步本表：

| 稿页 | 实现 | 深链 |
|------|------|------|
| Settings（告警设置 + Docker 设置） | 一级 Tab `tab=config`（默认） | `/setting` |
| 探针管理 | 一级 Tab `tab=probe` | `/setting?tab=probe`（行内升级仅 `admin`） |
| 用户管理 | 一级 Tab `tab=users` | `/setting?tab=users`（仅 `admin`） |
| 角色管理 | 一级 Tab `tab=roles` | `/setting?tab=roles`（仅 `admin`） |
| 接口管理 | 一级 Tab `tab=apis` | `/setting?tab=apis`（仅 `admin`） |
| 审计日志 | 一级 Tab `tab=audit` | `/setting?tab=audit` |

**一级 Tab 共 6 项**，与稿完全对齐：`config` 与 `audit` 对全部登录用户可见；`users` / `roles` / `apis` 都涉及权限编排，统一仅 `admin` 可见且**同处一级**——不把角色与接口藏进用户管理内部。

Tab 切换不整页刷新，标题区、eyebrow 与右侧 mono 标签（`Alerts · Docker · Users · Roles · APIs · Audit`，非管理员为 `Alerts · Docker · Audit`）保持不变。

一级 Tab 条为分段控件：条体 `surface-secondary` 底、`border-subtle` 描边、`radius-md`、内边距 4；每个 Tab **等宽**（`min-width: 96px`，内边距 7/14，选中态 `surface-card` 底 + `accent-primary` 文字 + `border-primary` 描边），条宽只由 Tab 数量决定（管理员 6 项约 606px、非管理员 3 项约 306px），左对齐，不随标签字数或选中状态抖动。`tab` 与 URL 双向同步，无权访问的 Tab 不渲染，深链带无权 Tab 时回落 `config`。列表页沿用固定高度宿主（560），让表格 `height="100%"` 有确定参照；用户管理的新建弹层由工具栏按钮的票据驱动。

后端接口：`/user/user_*`、`/role/role_query|create|update|delete`、`/resource/resource_query`、`/audit/*`，均已注册。标题区 eyebrow 为 `CONFIGURATION`，副标题「告警通道、Docker 镜像加速、账户权限与审计日志的统一配置入口。」（右侧 mono 标签沿用上一段统一口径，不在此重复）。

### 9.2 告警设置

| 卡片 | 内容 | 交互 |
|------|------|------|
| 告警阈值配置 | CPU / 内存 / 磁盘三条，各带「编辑」 | 打开对应阈值 Modal，修改后立即生效 |
| 邮箱服务配置 | 服务器、端口、发信账号、密码（`••••` + 「已保存」）、收信账号、最近测试结果 | 「发信测试」按钮 |
| 群机器人告警 | 钉钉 / 飞书 / 企业微信三个渠道：启用开关、Webhook（`access_token=••••` 遮罩展示）、加签 / 签名密钥、提醒范围、保存与测试 | 卡内顶部三渠道 Tab 切换，一次只渲染当前渠道面板以压缩卡片高度；Tab 项等宽 132px，选中态 `surface-card` 底 + `accent-primary` 文字；行右侧徽标展示该渠道已保存的启用状态，切换不触发保存，各渠道未保存的本地编辑互不影响 |
| 告警消息标题 | 当前标题「Beacon 主机告警」 | 打开标题 Modal，≤128 字符、禁止换行 |

阈值语义：**指标连续超过阈值才触发**，CPU / 内存 / 磁盘统一为「连续时长 + 阈值」判定——连续窗口内所有样本严格高于阈值才进入告警，单点越限、样本缺失或间断属于未知状态（与 monitoring-platform I013 一致）。修改后立即对邮件与群机器人（钉钉 / 飞书 / 企业微信）各通道同步生效。任一通道发送失败不影响其他通道。

### 9.3 Docker 设置

镜像仓库配置卡：textarea 每行一个 URL，保存后由 Agent 写入 `daemon.json` 并重启 Docker。留空则取消加速并恢复官方仓库；写入失败时**保留原有配置**并在弹层内提示原因。

### 9.4 五个列表页

| 页 | 列 | 备注 |
|----|----|------|
| 探针管理 | Agent 名称 / 状态 / 版本 / 操作 | 名称列主行为 Agent 名称、次行为 mono 的 `os/arch`；状态用 `BeaconStatusTagSuccess`（在线）/ `BeaconStatusTagDanger`（离线）；版本列为 mono 版本号，落后于最新版本时加「可升级」软標签（`accent-primary-soft` 底 + `accent-primary` 文字），与最新一致时标「最新」；操作列为 `circle-arrow-up` + 「升级」，仅管理员可见且离线时置灰；工具条左侧为名称搜索与状态分段筛选（全部 / 在线 / 离线带计数），右侧为「自动更新」开关（仅管理员可见，部署级门禁关闭时置灰并给出原因）、「共 N 个探针，M 个在线」与刷新 |
| 用户管理 | 用户名 / 角色 / 状态 / 创建时间 / 操作 | 状态用 `BeaconStatusTagSuccess/Danger/Neutral`；分页「共 5 条 / 10 条/页」 |
| 角色管理 | 角色名 / 权限 / 状态 / 创建时间 | 权限列以**资源树**呈现（分组名 + 接口名），三类角色统一为同一形态；角色分 `viewer` / `operator` / `admin`；区段页头右侧「新建角色」主按钮（`plus` 图标）打开 `Create Role Modal`，弹层内定义角色名与备注、按接口分组勾选可授权接口（带「已选 N / M」、「全选」「清空」）、状态默认启用 |
| 接口管理 | 接口名称 / URL / 请求方式 / 状态 | URL 用 mono；请求方式用 tag（GET / POST）；清单一次性加载，不提供分页 |
| 审计日志 | ID / 用户 / Agent / 操作 / 操作时间 | 记录用户与 Agent 的关键操作，保留 90 天；「共 128 条」 |

---

## 10. 认证与异常态

### 10.1 登录（`Login`）

左右分栏，各 720 宽（1440 视口），不提供重排稿：低于 §14 声明的 1280 最小支持宽度时出现横向滚动。

- **左栏**：品牌面板，`accent-primary` 与 `brand-deep` 的 225° 线性渐变（`color-mix(accent-primary 40%, brand-deep)` → `brand-deep`，保证两档主题下前景为浅色且可读），越界裁切两处装饰圆形：上右 `accent-primary` 420、下左 `accent-primary-bright` 300。内容上下分置（`space-between`，内边距 56）：顶部品牌行（32 品牌标图片 + 「Beacon」`font-size-xl`/700），底部居中文案（「面向 Server-Agent 的可观测性平台」`font-size-lg`/500 + 「实时监控 · 容器编排 · 安全隧道」`font-size-md`），文案用浅色（`#FFF` 与 72% 白），不使用 `foreground-*` 系列（它们随主题翻转，会在深色品牌面板上失效）。
- **右栏**：`surface-primary` 底，居中 400 宽卡片（`surface-card` + `border-subtle` + `radius-lg` + 内边距 32 + 区块间距 24，无投影）。卡片内自上而下：语言选择器（chip：`surface-inset` + `border-primary` + `radius-sm`，内边距 6/10，13px `foreground-secondary` + `chevron-down`）右对齐、标题「登录」`font-size-xl`/700 居中、用户名 / 密码两个输入框、主按钮「登 录」（`accent-primary` 实心、`radius-sm`、高 48、14px/600 `foreground-on-accent`）。
- **输入框**：字段容器自身承担描边——`surface-inset` 底 + `border-primary` 1px + `radius-sm` + 内边距 12/16，左侧 16px `foreground-muted` 图标（`user` / `lock`，必须 `currentColor`，不得写死颜色），占位符 14px `foreground-placeholder`；聚焦转为 `accent-primary` 描边 + 3px 主色光晕，校验失败转为 `accent-danger` 描边。密码框保留可见性切换。
- **校验与反馈**：登录只校验必填，密码策略由服务端裁决（客户端不设长度门槛）；认证失败与网络异常由全局响应拦截器统一提示，页面保持原状并结束提交态，不做二次跳转。

### 10.2 强制改密（`Change Password`）

独立页，居中卡片。字段：用户名（只读，填 `admin`）、原密码、新密码。

**硬约束：强制改密期间抽屉不提供关闭入口，遮罩点击与 Esc 均被禁用。** 文案必须写明「首次登录或密码已过期」与「密码至少 12 位，需包含大小写字母、数字、符号中的至少三类，修改成功后需重新登录」（与 account-administration IA006 对齐）。

### 10.3 非强制改密（`Update Password Drawer`）

从顶栏账户菜单进入的 480 宽抽屉。与强制改密的区别：**有取消按钮、可关闭**。

### 10.4 404

**兜底路由渲染独立错误页，不套用工作区导航栏与状态栏。** `surface-primary` 整页居中一列（间距 8）：`file-question` 图标（`foreground-muted`）→ 「404!」（`font-size-2xl`/700）→ 说明句（`font-size-md`，`foreground-secondary`）→ 「返回首页」chip 主按钮（`accent-primary` + `radius-sm` + 内边距 8/16 + `arrow-left` 14 + 13/600 `foreground-on-accent` + `0F172A1A` 阴影），点击回到监控工作区。文案走 `error.notFound` / `error.backHome`，不使用插画与 emoji。

---

## 11. 弹层体系（Modal / Drawer）

### 11.1 壳体规格

| 项 | Modal | Drawer |
|----|-------|--------|
| 宽度 | **480** 短表单（多数）/ **560** 多字段长表单（安装 Agent、升级 Agent、告警邮箱） | **480**（容器、改密）/ **540** 多行 textarea（镜像仓库）/ **720** 双列（用户表单） |
| 高度 | 内容自适应 | 全高 |
| 底色 | `surface-elevated` | `surface-elevated` |
| 描边 / 圆角 | `border-subtle` 1px / `radius-md` | `border-subtle` 1px / 左缘贴屏 |
| 阴影 | `0 8px 24px #00000066` | 无 |
| 遮罩 | `overlay-scrim` | `overlay-scrim` |
| 头部 | 内边距 24，下边框 | 内边距 16/24，下边框 |
| 主体 | 内边距 24，间距 16 | 内边距 16（用户表单 24），间距 8 |
| 底部 | 内边距 24，右对齐，上边框 | 内边距 16/24，右对齐，上边框 |

**宽度阶梯（唯一口径）**：Modal 三档 `480 / 520 / 560`，Drawer 三档 `480 / 540 / 720`。档位由内容形态决定且只有这 6 个值——短表单用 480，单组凭据用 520，多字段长表单用 560；容器与改密用 480，多行文本用 540，双列表单用 720。禁止按页面另起宽度。历史验收口径「480px Modal / 540px Drawer」的单一档位记录**作废**，以本节为准。

**头部**（`BeaconDialogHeader`）：28×28 `accent-primary-soft` 图标块 + 标题（14/600）+ 右侧关闭图标（16，`foreground-muted`）。图标枚举共 17 个：`bell-ring / circle-arrow-up / circle-check / circle-x / delete / download / gauge / mail-check / network / package-plus / server / shield-plus / triangle-alert / upload / user / webhook / x`。三处弹层稿面不设图标块、只保留标题与关闭：`Install Agent Modal`、`Confirm Delete Dialog`、`Update Password Drawer`。

### 11.2 按钮

| 变体 | 底 | 字 | 边框 | 几何 |
|------|----|----|------|------|
| Primary | `accent-primary` | `foreground-on-accent` 500 | — | 内边距 8/16，`radius-sm`，13px |
| Ghost | `surface-inset` | `foreground-primary` 500 | `border-primary` 1px | 同上 |

底部按钮顺序恒为「Ghost 在左、Primary 在右」。设计稿另定义了 Secondary / Danger 变体（见 `Components Showcase`），但弹层稿中只使用 Ghost 与 Primary。

### 11.3 表单字段（`BeaconFormField`）

label 13px `foreground-secondary`（与输入框间距 4）→ 输入框高 36、内边距 0/12、`surface-inset` 底、`border-primary` 1px、`radius-sm`、文字 13px `foreground-primary`。多行输入以独立 textarea 呈现。

凭据类字段一律遮蔽显示（`••••`）；已保存的凭据旁标「已保存」；留空表示保留已保存值。

### 11.4 危险操作统一模式

所有删除类操作走共享 `ConfirmDialog`：

- 文案必须包含**操作对象 + 执行通道 + 不可恢复后果**。范例：「确认删除容器? 删除操作经反向隧道下发到目标节点执行，容器内未持久化的数据将无法恢复。」
- 状态机：`等待确认 → 确认中（按钮 loading，弹窗保持打开，禁用重复提交）→ 失败反馈（保留弹窗与本地输入，由全局请求拦截器统一提示）`。
- 取消不触发任何请求。

### 11.5 弹层清单与关键约束

| 弹层 | 关键约束 |
|------|---------|
| 安装 Collia Agent | Server 地址 / 节点标识 / 安装令牌（`••••` 遮蔽展示；令牌是部署期静态持久值，无 TTL 语义，见 release-install IR008/IR014）；命令块 mono；安装命令中的 token 必须遮蔽；需 root 或 sudo；自动识别 amd64 / arm64，落地 `/usr/sbin/collia` |
| 升级 Collia Agent | 目标节点 `beacon-server · ag_7f2a91c4`、当前版本、目标版本为「Server 镜像内置版本」；**同一 Agent 不支持并发升级**，全局并发受限，排队节点轮到时自动开始；升级期间该节点数据短暂中断 |
| 新建 / 编辑容器 | 容器名、镜像、网络、端口、重启策略（`unless-stopped`）、环境变量、目录挂载、标签（可选）；多值字段每行一条 |
| 拉取镜像 | 镜像名支持 `name:tag`，省略 tag 默认 latest；携带当前 Agent 标识 |
| 导入镜像 | 拖拽上传 `.tar` / `.tar.gz`，仅支持单个；经 `X-Agent-ID` 提交后转发给对应节点 |
| 清理虚悬镜像 | 仅删除 `<none>:<none>`；被容器引用的镜像不受影响；清理后镜像层无法恢复 |
| 新建网络 | 名称、模式（bridge）、子网、网关、标签；新建后需在容器编辑时选择才生效 |
| 镜像仓库配置 | 每行一个 URL；写入失败保留原配置并提示 |
| 邮箱服务配置 | 服务器 / 端口 / 发信账号 / 密码 / 收信账号（英文逗号分隔）/ 发信测试 |
| 告警阈值（CPU / 内存 / 磁盘） | 连续时长 + 阈值 + 当前阈值对照；修改后立即生效并同步两条告警通道 |
| 告警消息标题 | ≤128 字符、不支持换行；校验文案必须可见 |
| 创建 / 编辑用户 | 用户名唯一、角色、备注、状态（禁用即不可登录）；角色决定可见菜单与接口权限 |
| 修改密码 | 强制版无关闭入口；非强制版可取消；改密后需重新登录 |

---

## 12. 组件基线与状态模型

### 12.1 设计局部符号

稿内以 symbol 形式复用组件（`beacon.ax` 10 个）：

| 符号 | props | 实例数 |
|------|-------|--------|
| `BeaconNavItem` / `BeaconNavItemActive` | `label`, `icon`(activity/package/settings) | 24 / 12 |
| `BeaconStatusItem` | `label`, `value` | 63 |
| `BeaconStatusTagSuccess` / `Danger` / `Neutral` | `label` | 16 / 3 / 3 |
| `BeaconDialogHeader` | `title`, `badge`(13 枚举) | 0（预留） |
| `BeaconDialogButtonPrimary` / `Ghost` | `label` | 16 / 39 |
| `BeaconFormField` | `label`, `value` | 0（预留） |

> 这些是**设计稿内部的可复用构件**，不是前端组件名；实现侧不要求同名文件，但几何与语义必须等价。

**符号 → 实现映射**（`design_query --mode=reconcile` 按 JSX 开标签匹配，对 Vue SFC 会全部报 `present: false`；对账时以下表为准，不以工具的默认匹配结果判定缺失）：

| 设计符号 | 实现落点 | 说明 |
|---------|---------|------|
| `BeaconNavItem` / `BeaconNavItemActive` | `beacon/web/src/layout/navbar/index.vue` + Element Plus `el-menu` 项 | 激活态由 `--el-menu-active-bg-color` 与 `--el-menu-active-color` 驱动 |
| `BeaconStatusItem` | 各工作区页底部状态栏组件（9 项） | 稿中每页固定 9 个实例（左组 4 + 运行时间 1 + 右组 4） |
| 底部状态栏 Beacon 版本项（右组末项） | `beacon/web/src/layout/statusbar/index.vue` | 服务端构建版本，label「Beacon」+ mono 值，固定在状态栏最右端 |
| 顶栏 Beacon 版本提示（升级图标） | `beacon/web/src/layout/navbar/BeaconVersion.vue` + `beacon/web/src/api/system/index.ts` | circle-arrow-up 图标常驻；有新版本时高亮「有新版本」，弹层只读展示当前 / 最新版本与手动升级说明 |
| 探针管理 Tab（Agent 列表） | `beacon/web/src/views/setting/probe/index.vue` | 一级 Tab `tab=probe`；列表列为 Agent 名称（附 `os/arch`）/ 状态 / 版本 / 操作，升级操作仅管理员可见且离线置灰；名称搜索与状态筛选在前端完成（Agent 量级小），版本列对比 `GET /api/v1/system/version` 的最新版本标记「可升级 / 最新」 |
| `BeaconStatusTagSuccess` / `Danger` / `Neutral` | Element Plus `el-tag` + `--el-color-success/danger` 与 `foreground-muted` | 状态必须带文字，不只靠颜色 |
| `BeaconDialogHeader` | `beacon/web/src/components/DialogHeader/index.vue` | 28×28 `accent-primary-soft` 图标块 + 14/600 标题，经 `#header` 插槽挂载；关闭入口为 Element Plus header 关闭按钮（16，`foreground-muted`） |
| `BeaconDialogButtonPrimary` / `Ghost` | `el-button` 的 `type=primary` / 默认 | 底部顺序 Ghost 左、Primary 右 |
| `BeaconFormField` | `el-form-item` + `el-input` | label 13px、输入高 36 |
| 监控图卡图标块（28×28 soft 色） | `beacon/web/src/styles/element.scss` `.am-chart-card-icon` | CPU / 内存 / 磁盘 / 网络四类，soft 底取语义色 |
| 监控图卡平滑曲线区（纵轴刻度 + 网格线 + 折线） | `beacon/web/src/components/Echarts/index.vue` + `beacon/web/src/config/echarts.ts` | CPU / 内存单序列，磁盘 / 网络双序列叠加；`smooth: true` + `showSymbol: false` |
| 监控图卡序列色与图例 | `beacon/web/src/utils/chartTokens.ts` + `beacon/web/src/config/echarts.ts` | 序列色按 `chart-1`–`chart-6` token 解析（图例色 = 序列色），多序列卡 `legend.data` 与实际序列名一致 |
| 容器监控图例 chip 区（容器名 + 百分比） | `beacon/web/src/views/monitor/container/index.vue` `.am-chart-card-legend` / `.am-legend-item` | 4 列换行 chip，7×7 圆点色与折线同源（`chart-1`–`chart-6`），百分比用 mono 小字 |
| 容器监控多序列曲线区（每容器一条线 + 图内名） | `beacon/web/src/components/Echarts/index.vue` + `beacon/web/src/config/echarts.ts`（`containerCpuOption` / `containerMemOption`） | ECharts `title` 给图内名（「CPU」/「内存」），`legend.show: false`，数值纵轴 + 时间轴 |
| Agent 上下文条 | `beacon/web/src/components/Agent/MonitorContextBar.vue` | 主机与容器分组共用同一组件，保证 §7.3 要求的呈现一致 |

### 12.2 组件基线清单（`Components Showcase`）

Button（Primary / Secondary / Danger / Ghost，MD / SM）、FormField（标签 / 输入 / 选项 / 备注）、SettingsCard、ChartCard、ModalShell / DrawerShell、状态标签、反馈提示（Toast）、数据表格、空状态。

基线页同时声明：**加载态保留骨架，错误态由全局请求拦截器统一提示，空态给出引导与刷新入口。**

**基线页的基准范围**：`Components Showcase` 中的 ModalShell / DrawerShell 只提供**最小档 480 的几何与壳体规格**（头部、主体、底部内边距与描边）作为基准；宽度不是基准项，完整宽度阶梯见 §11.1。页面实例使用 520 / 540 / 560 / 720 不构成对基线页的偏离。

**规范页归属**：分组索引中不单设「Design Tokens」与「States」两页——色板与 token 表的规范来源是 §5，状态模型与表达是 §12.3。改稿时不得再引用这两页；若未来需要可视化色板，按「新增页面的最低要求」（§15）正式建页并在此登记。

### 12.3 状态模型

| 状态 | 表达 | 适用 |
|------|------|------|
| 默认 | 内容直接呈现 | 全局 |
| Hover | `accent-primary-hover`（按钮 / 导航项 / 链接与分页、下拉项）；表格行悬停用 `accent-primary-soft` 浅底。已落地为 `styles/element.scss` 全局规则 | 按钮、导航项、表格行、链接 |
| Focus | `outline: 2px solid accent-primary` + `offset 2px`，跟随元素圆角；仅 `:focus-visible` 触发，鼠标点击不呈现；遮罩内的首个可聚焦元素默认获得焦点 | 全局（按钮、导航项、输入框、表格行、链接、弹层内控件） |
| Disabled | `foreground-placeholder` 文字 / 禁用描边 / 不响应 | 受保护容器的启停与删除、批量操作未选时 |
| Loading | 骨架保留布局尺寸 | 列表、图表卡、工作区首次挂载 |
| Empty | 标题 + 一句因果说明 + 可执行入口 | 无 Agent、无搜索结果 |
| Stale / 过期 | `accent-warning` 标注 + 「最后采样」时间 + 明示不判断阈值 | 监控采样中断 |
| Error | 全局请求拦截器统一提示；不内嵌到卡片 | 请求失败 |
| Success | 同上，但成功反馈与失败分开 | 请求成功 |

---

## 13. 文案与语气

**控制台**：中文为主，技术名词保留英文（Agent、Docker、CPU、Webhook、token、daemon.json）。空态写法固定为「发生了什么 + 为什么 + 怎么做」，例如「在目标主机执行安装命令接入 collia Agent，连接成功后会自动展示当前节点的监控数据。」危险操作必须写明后果与不可恢复性。状态提示用短句：「运行中」「已停止」「实时」「数据过期」「等待 Agent 续传」。

**通用禁令**：不写真实密钥、Token、内部凭据或不可公开的环境值；稿中出现凭据时必须遮蔽（`••••`、`[REDACTED]`）。

**可公开标识（白名单）**：以下真实标识属开源公开信息，**必须原样保留**，不遮蔽、不复写成占位符：

| 标识 | 出现位置 | 性质 |
|------|---------|------|
| `registry.cn-hangzhou.aliyuncs.com/amuluze/beacon` | 容器工作区镜像说明 | 公开镜像仓库 |
| `github.com/amuluze/beacon` | 仓库 README 与开源说明 | 开源公开仓库 |

白名单只覆盖上表两项。**新增任何真实地址、账号、编号前必须先登记到本表**，未登记的一律按通用禁令处理。凭据类字段（安装令牌、Webhook `access_token`、加签密钥、邮箱密码）永不允许出现在稿中明文。

---

## 14. 无障碍与响应式

- 状态**不只用颜色**表达：始终有文字（「运行中」而非仅绿点）与图标配合。
- 正文最小字号 11px（仅用于时间戳与轴标签），控件正文 13px 起。
- 控件最小可点区域 32px（顶栏图标按钮），表单控件高 36px。
- 控制台**只设计 1440px 一档版面**，并声明最小支持宽度 **1280px**：1280–1439 之间保持同一版面结构，由栅格与表格列按比例收缩，不重排、不隐藏列；顶栏与底部状态栏固定高度不压缩；低于 1280 出现横向滚动条，禁止折叠导航或改为移动布局。控制台是数据密集型工作面，不提供平板/移动重排稿。
- **键盘焦点环**：2px `accent-primary` + 2px offset，跟随元素圆角，只在 `:focus-visible` 下呈现（鼠标点击不出环）。实现落在 `beacon/web/src/styles/element.scss` 的全局 `:focus-visible` 规则，不逐组件重复定义。焦点环不得被 `overflow: hidden` 容器裁切，弹层内的首个可聚焦元素默认获得焦点。
- 键盘可达性底线：所有交互元素可 Tab 到达、Enter / Space 触发、Esc 关闭弹层；焦点不得陷入遮罩后方内容。

---

## 15. 变更治理与视觉验收

**改稿流程**：改控制台设计稿（`.pen/beacon.ax`，在本仓库之外）→ 回写本文对应章节 → 实现对齐 → 视觉验收 → 若行为变更则同步 `.specs/domain/`。

**新增页面的最低要求**：必须落在既有分组内；必须复用既有 token 与组件几何；必须补齐空态、加载态、错误态；必须在本文的页面地图中登记。

**视觉验收口径**

- 桌面端以 1440px 视口对照控制台 Monitor / Container / Settings 三工作区。
- 对照对象是**本文与稿件**，不是实现现状。
- 若执行环境没有可用浏览器，必须把视觉验收标记为**未完成并说明阻塞**，不得以源码检查冒充截图验收。

**实现侧的硬约束**

- 色值、间距、圆角、字号一律取 token，不写字面量。
- 新增 token 必须先在稿中定义并在本文章节登记。
- 改动影响弹层、表单、状态标签的规格时，必须回到 §11、§12 的规格条目重新对照，而不是各页面自由发挥。
