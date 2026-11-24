# 文档系统

## 所有权
- AGENTS 是协作入口，CLAUDE 为逐字同步副本。
- Domain Spec 只记录长期行为、状态和失败语义。
- .docs 记录当前实现；DESIGN 记录设计意图。
- 项目模块清单以 `.docs/project-analysis.md` 的 Workspace Modules 为统一来源。
- 目录表是各文档按职责选择的导航，不是完整磁盘清单。
- .beaver、.claude、.omo、.zcode 和 .codegraph 是被忽略的本地运维或工具状态。
- 本地工具不属于仓库交付物；如转为共享工具，应先审查凭据与环境依赖再纳管。

## 门禁
运行 task doc:check，退出码为当前检出状态的结果；文档不保存未经执行即可宣称的 pass。
必需依赖：Python 3、Ruby/Psych、Go 和 Task；缺失依赖不得视为验证成功。
Ruby 脚本固定 UTF-8 编码并以 `RUBYOPT=-EUTF-8` 调用，结果不依赖调用方的 locale。
门禁自测（`scripts/doc-check-test.rb`）随每次 `task doc:check` 执行，证明门禁能拒绝坏输入。
`scripts/doc-schema` 是独立 Go module：字段形状导出与OpenAPI标准校验器纳入 `task doc:verify` 的 build/test/vet。
OpenAPI通过固定版本 `github.com/getkin/kin-openapi v0.133.0` 校验标准结构，并额外检查根级和operation级security方案引用；禁止外部引用加载。依赖记录在工具模块的go.mod/go.sum，首次运行需要下载依赖，不要求额外安装Python/Node验证库。
`.docs/MANIFEST.yml` 用完整 YAML 解析，required 为必需集合，optional 对已存在的匹配文件执行规则。
路由从 `beacon/service/router.go` 的 Fiber 注册结果导出，与 `.docs/api/openapi.yml` 对账 HTTP 方法和路径；不发送业务请求，自动 HEAD 不重复计数。
路由测试不证明鉴权和业务语义正确。另由 `scripts/doc-contract-check.rb` 检查全部operation的成功响应建模、204无正文、引用有效性、handler请求类型与Go字段模型；不把字段对账视为真实HTTP响应验收。
普通JSON接口的响应引用沿handler返回变量绑定service方法的返回类型；无法解析的返回形式必须失败。成功状态集合精确对账，204接口不得额外声明200。升级、上报、安装下载、探针和WebSocket采用显式人工维护的绑定表，不能把这些绑定视为已完成运行时响应验证。升级请求与结果、Agent自有字段和统一错误响应也参与Go字段对账；已删除的派生类型不得保留陈旧组件。GORMModel的第三方嵌入字段、HealthResponse等inline map响应仍需人工审查与运行时测试。
Query参数使用Go的query标签，分别核对名称、字段形状、边界和必填状态；不从JSON字段名推断Query别名。空schema不得作为响应模型通过门禁。
Go字段变化后用 `ruby scripts/doc-contract-check.rb --update-go-schemas` 刷新派生组件及已识别handler的Query参数，再评审handler与响应类型绑定。验证标签保存在 `x-go-validate`；description与pattern标签也进入派生schema。Agent上报标识的必填、长度和字符范围，以及新密码至少12字符，参与机器校验；密码Unicode字符类别与新旧密码不同等规则写入描述，仍以service测试为准，标准schema校验不替代这些业务规则。
required 标量按实际 Go validator 非零语义导出：字符串非空、数值非零、布尔值为真；切片/映射的 required 仅排除 null，不额外禁止非 nil 空集合。用户/角色状态与模型及前端统一为1启用、2停用；创建必须显式提供，编辑的0或省略表示不修改。刷新接口独立声明 RefreshCookieAuth，门禁固定检查刷新 Cookie 和 CSRF Cookie/请求头声明；该绑定不是自动推导全部中间件语义。
门禁自测包含实际JSON样本校验，覆盖短密码、缺失/空/非法/超长Agent ID，以及合法密码与Agent标识。
Token 校验只检查名称和数量，不证明视觉一致。
源码库存按Git可见且磁盘实际存在的文件统计；已删除但尚未暂存的跟踪文件不计入库存，不包含本地生成文件。
源码库存同时对账模块文档声明的 Go 源码/测试文件数、前端源码文件数与运行时依赖清单；数字必须与仓库实际一致。
部署配置清单必须覆盖仓库中全部运行时配置文件，且登记项必须真实存在。
模块文档的「接口与符号信号」与「配置与入口信号」不得重复承载同一批对象，也不得把测试文件当作入口信号。
Domain Spec 必须声明非占位符的反推来源与最后更新时间；跨 Spec 引用必须双向对称。
Domain Spec 编号全局唯一，不得被另一份 Spec 重复占用；跨 Spec 编号引用必须指向已定义约束。
文档之间的相对链接必须解析到仓库中真实存在的文件；文件链接存在不代表标题锚点有效。
入口文件一致性单独核对：根 `AGENTS.md` 与 `CLAUDE.md` 逐字一致，各模块 `AGENTS.md` 与 `CLAUDE.md` 成对存在且内容一致，CLAUDE 不得出现未纳管的独立条目。

## 证据与维护
Domain 约束须关联具体测试及其断言；仅列测试名不能声明整个约束已满足。
全部Domain约束已在 `.docs/project-analysis.md` 的 Evidence Backlink Table 逐条登记：具名测试仅声明部分覆盖；只有源码审查入口的条目保持待专项验收；不声明行为全部验收通过。
门禁对账编号集合、证据文件、测试声明、非空断言范围与残余缺口，不能自动证明断言足够覆盖约束。
测试库存完整展开所有路径，并与Git可见测试文件集合逐项对账，拒绝重复、遗漏或失效路径。
模块命令按逻辑续行及工作目录规范化；验收记录与命令集合双向对账，记录通过、失败、未执行原因。记录为历史快照，不证明当前源码仍通过。
`task doc:check` 校验证据记录完整性，允许显式保留行为缺口和历史失败；`task doc:verify` 实际运行六个工作区 Go 模块与独立的 `scripts/doc-schema` 工具模块的build、race测试和vet，并在任一失败时返回非零。常驻服务、依赖改写、发布镜像和远端部署验收不由该安全子集执行。
PostgreSQL数据库专项测试以BEACON_TEST_POSTGRES_PORT等测试环境变量显式启用，默认跳过；具体运行条件和已验收范围见 `.docs/deployment.md` 与 `.docs/project-analysis.md`。
依赖 Docker 守护进程或镜像仓库的模块测试由 `scripts/doc-commands.rb` 分类为 `requires=docker`，`task doc:verify` 默认跳过并打印原因；设置 `BEACON_DOC_VERIFY_DOCKER=1` 开启。测试镜像默认alpine:3.20，可用 `BEACON_DOCKER_TEST_IMAGE` 指定完整且具有sh命令的测试镜像；通过指定镜像的专项测试不等于默认镜像缓存故障已修复。
模块命令的本次执行快照与原因详见 `.docs/project-analysis.md`；新增命令须补齐对应记录。
新增统计须有自动对账或明确历史快照范围，避免人工复制计数。
Go、前端及Ruby具名测试可回挂约束证据。
变更后运行门禁；任务和计划属于过程材料，完成后清理。
