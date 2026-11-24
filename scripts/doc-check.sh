#!/usr/bin/env bash
# 文档系统完整性检查。
#
# 覆盖：MANIFEST 解析与内容契约、质量规则、登记项双向覆盖、Domain Spec
# 结构/编号/交叉引用与元数据、文档间链接、入口文件一致性与模块入口发现、
# 健康计数一致性、开发命令有效性（task 与脚本两类）、设计 token 覆盖、
# OpenAPI 与运行时路由对账、源码库存与模块文档信号节对账、部署配置清单覆盖。
#
# 用法：bash scripts/doc-check.sh        （或 task doc:check）
# 退出码：0 全部通过；1 存在失败项。

set -uo pipefail
cd "$(dirname "$0")/.." || exit 1

python3 - <<'PYEOF'
import os, re, json, sys, glob, difflib, subprocess, fnmatch

# Ruby reads UTF-8 documents; force the encoding so results never depend on the
# caller's locale (a C locale raises "invalid byte sequence in US-ASCII").
RUBY_ENV = {**os.environ, 'RUBYOPT': ('-EUTF-8 ' + os.environ.get('RUBYOPT', '')).strip()}

FAIL = []
def ok(msg):  print(f"  ✅ {msg}")
def bad(msg): print(f"  ❌ {msg}"); FAIL.append(msg)
def info(msg): print(f"     {msg}")

MANIFEST = '.docs/MANIFEST.yml'
resolve = lambda p: os.path.normpath(os.path.join('.docs', p))

# Ruby/Psych performs full YAML parsing; absence or invalid YAML is fatal.
def load_yaml(path):
    result = subprocess.run(['ruby', '-ryaml', '-rjson', '-e',
        'puts JSON.generate(YAML.load_file(ARGV[0]))', path],
        capture_output=True, text=True, env=RUBY_ENV)
    if result.returncode:
        raise RuntimeError(result.stderr.strip())
    return json.loads(result.stdout)

def parse_manifest(path):
    data = load_yaml(path)
    entries = []
    for bucket, groups in data.items():
        if bucket in ('version', 'quality'):
            continue
        groups = {'index': groups} if bucket == 'index' else groups
        for sub, items in groups.items():
            for item in items:
                entries.append(dict(min_lines=None, must_contain=[], pattern=False,
                    description='', **{k:v for k,v in item.items()
                    if k not in ('min_lines','must_contain','pattern','description')},
                    bucket=bucket, sub=sub))
                entries[-1].update(item)
    return entries, data.get('quality', {}).get('warning_rules', [])

ENTRIES, RULES = parse_manifest(MANIFEST)
required = [e for e in ENTRIES if e['sub'] == 'required' or e['bucket'] == 'index']
index_entries = [e for e in ENTRIES if e['bucket'] == 'index']
registered = {resolve(e['path']) for e in ENTRIES if e['path'] and not e['pattern']}
pattern_groups = [e for e in ENTRIES if e['pattern']]

print("== 文档系统检查 ==")
print(f"\n[1] MANIFEST 解析：登记 {len(registered)} 项 + {len(pattern_groups)} 个通配规则，质量规则 {len(RULES)} 条")
if registered:
    ok("MANIFEST 可解析")
else:
    bad("MANIFEST 解析为空，路径格式可能已变更")

# ── 2. MANIFEST YAML 语法 ────────────────────────────────────────────────
print("\n[2] MANIFEST YAML 语法")
ok("完整 YAML 解析成功（Ruby/Psych 为必需依赖）")

# ── 3. 登记项双向覆盖 ────────────────────────────────────────────────────
print("\n[3] 登记项双向覆盖")
dangling = sorted(p for p in registered if not os.path.exists(p))
if dangling:
    bad(f"登记但文件不存在 {len(dangling)} 项"); [info(p) for p in dangling]
else:
    ok("无悬空登记")

actual = set()
for pat in ('.docs/**/*.md', '.docs/**/*.yml', '.specs/domain/*.md'):
    for f in glob.glob(pat, recursive=True):
        if os.path.isfile(f) and f != MANIFEST:
            actual.add(f)
# optional 通配规则视为「已纳管」：匹配到的新文件不算孤儿，避免与 MANIFEST 的 optional 声明自相矛盾。
pattern_managed = set()
for e in pattern_groups:
    pattern_managed.update(glob.glob(resolve(e['path'])))
orphan = sorted(actual - registered - pattern_managed)
if orphan:
    bad(f"文档未登记 {len(orphan)} 项"); [info(p) for p in orphan]
else:
    ok("所有 .docs 与 Domain Spec 均已登记")

# ── 4. 内容契约：min_lines 与 must_contain ───────────────────────────────
print("\n[4] 内容契约（min_lines / must_contain）")
contract_bad = []
checked_files = set()
for e in ENTRIES:
    targets = []
    if e['pattern']:
        targets = sorted(glob.glob(resolve(e['path'])))
    elif e['path']:
        targets = [resolve(e['path'])]
    for t in targets:
        if not os.path.isfile(t):
            continue
        checked_files.add(t)
        body = open(t, encoding='utf-8', errors='ignore').read()
        if e['min_lines'] is not None:
            n = body.count('\n') + (1 if body and not body.endswith('\n') else 0)
            if n < e['min_lines']:
                contract_bad.append(f"{t} 仅 {n} 行，低于下限 {e['min_lines']}")
        for needle in e.get('must_contain') or []:
            if needle not in body:
                contract_bad.append(f"{t} 缺少必备内容「{needle}」")
if contract_bad:
    bad(f"内容契约未满足 {len(contract_bad)} 项"); [info(x) for x in contract_bad]
else:
    ok(f"{len(checked_files)} 份文档的行数下限与必备章节全部满足")

# ── 5. 质量规则 ──────────────────────────────────────────────────────────
print("\n[5] 质量规则（quality.warning_rules）")
rule_codes = {r['code'] for r in RULES}

if 'low_signal_description' in rule_codes:
    groups = {}
    for e in ENTRIES:
        if e['description']:
            groups.setdefault(e['bucket'], []).append(e)
    low = []
    for bucket, items in groups.items():
        for i in range(len(items)):
            for j in range(i + 1, len(items)):
                a, b = items[i]['description'], items[j]['description']
                if not a or not b:
                    continue
                ratio = difflib.SequenceMatcher(None, a, b).ratio()
                if ratio >= 0.85:
                    low.append(f"{items[i]['path']} 与 {items[j]['path']} 的 description 相似度 {ratio:.0%}（模板化）")
    if low:
        bad(f"description 低信号 {len(low)} 处"); [info(x) for x in low]
    else:
        ok("description 均为具体描述，无模板化重复")

if 'unmanaged_ai_entry' in rule_codes:
    agents = {}
    for f in glob.glob('**/AGENTS.md', recursive=True):
        if 'node_modules' in f or '/dist/' in f:
            continue
        agents[f] = open(f, 'rb').read()
    unmanaged = []
    for f in glob.glob('**/CLAUDE*.md', recursive=True):
        if 'node_modules' in f or '/dist/' in f:
            continue
        body = open(f, 'rb').read()
        if not any(body == v for v in agents.values()):
            unmanaged.append(f)
    if unmanaged:
        bad(f"存在未纳入 AGENTS.md 的 CLAUDE 指引 {len(unmanaged)} 个"); [info(x) for x in unmanaged]
    else:
        ok("所有 CLAUDE 指引均已纳入对应 AGENTS.md")

if 'stale_generated_doc' in rule_codes:
    stale = []
    for f in sorted(glob.glob('.docs/_archive/*.md')):
        base = os.path.basename(f)
        for other in glob.glob('.docs/**/*.md', recursive=True):
            if other == f:
                continue
            # 只判定链接式引用；清单式列出（文档目录）不算陈旧引用。
            body = open(other, encoding='utf-8', errors='ignore').read()
            if re.search(r'\]\([^)]*' + re.escape(base) + r'[^)]*\)', body):
                stale.append(f"归档文档 {f} 仍被 {other} 链接引用")
    docs_basenames = {os.path.basename(f) for f in glob.glob('.docs/**/*.md', recursive=True)}
    for f in ('README.md', 'README.en.md', 'architecture.md', 'deployment.md'):
        if f in docs_basenames and os.path.isfile(f):
            stale.append(f"仓库根存在 {f}，与 .docs 内同名文档重复")
    if stale:
        bad(f"陈旧或重复的生成文档 {len(stale)} 处"); [info(x) for x in stale]
    else:
        ok("无陈旧引用与重复文档")

# ── 6. Domain Spec 结构与编号 ────────────────────────────────────────────
print("\n[6] Domain Spec 结构与编号")
REQUIRED_SECS = ('## 约束空间', '## 状态模型', '## 错误条件与失败语义')
ids_all, struct_bad = {}, []
spec_files = sorted(glob.glob('.specs/domain/*.md'))
for f in spec_files:
    txt = open(f, encoding='utf-8').read()
    miss = [s for s in REQUIRED_SECS if s not in txt]
    if miss:
        struct_bad.append(f"{f} 缺 {miss}")
    for i in re.findall(r'^\| ([A-Z]{1,3}\d{3}) \|', txt, re.M):
        ids_all.setdefault(i, []).append(os.path.basename(f))
if struct_bad:
    bad(f"Spec 结构不完整 {len(struct_bad)} 份"); [info(x) for x in struct_bad]
else:
    ok(f"{len(spec_files)} 份 Spec 三段结构齐全")

dups = {k: v for k, v in ids_all.items() if len(v) > 1}
if dups:
    bad(f"编号重复 {len(dups)} 个"); [info(f"{k} -> {v}") for k, v in dups.items()]
else:
    ok(f"编号全局唯一，共 {len(ids_all)} 条")

# ── 7. 跨 Spec 编号交叉引用 ─────────────────────────────────────────────
print("\n[7] 跨 Spec 编号引用")
dangling_refs = []
for f in spec_files:
    txt = open(f, encoding='utf-8').read()
    body = re.sub(r'^\| ([A-Z]{1,3}\d{3}) \|', '|', txt, flags=re.M)
    for i in re.findall(r'\b([A-Z]{1,3}\d{3})\b', body):
        if i not in ids_all:
            dangling_refs.append(f"{os.path.basename(f)} -> {i}")
if dangling_refs:
    bad(f"悬空编号引用 {len(dangling_refs)} 处"); [info(x) for x in dangling_refs]
else:
    ok("跨 Spec 编号引用全部指向已定义约束")

# ── 8. 文档间链接 ────────────────────────────────────────────────────────
print("\n[8] 文档间内部链接")
targets = sorted(registered | pattern_managed | ({'DESIGN.md'} if os.path.exists('DESIGN.md') else set()))
broken, nlink = [], 0
for f in targets:
    if not os.path.isfile(f) or not f.endswith('.md'):
        continue
    base = os.path.dirname(f)
    for m in re.finditer(r'\[[^\]]*\]\(([^)#\s]+)(?:#[^)]*)?\)', open(f, encoding='utf-8').read()):
        link = m.group(1)
        if link.startswith(('http://', 'https://', 'mailto:')):
            continue
        nlink += 1
        if not os.path.exists(os.path.normpath(os.path.join(base, link))):
            broken.append(f"{f} -> {link}")
if broken:
    bad(f"断链 {len(broken)} 个"); [info(x) for x in broken]
else:
    ok(f"{nlink} 个内部链接全部有效")

# ── 9. 入口文件一致性与模块入口发现 ─────────────────────────────────────
print("\n[9] 入口文件一致性")
if os.path.exists('AGENTS.md') and os.path.exists('CLAUDE.md'):
    same = open('AGENTS.md', 'rb').read() == open('CLAUDE.md', 'rb').read()
    ok("AGENTS.md 与 CLAUDE.md 内容一致") if same else bad("AGENTS.md 与 CLAUDE.md 不一致（CLAUDE.md 需同步）")
else:
    bad("缺少 AGENTS.md 或 CLAUDE.md")

for p in sorted(registered):
    if not p.endswith('AGENTS.md'):
        continue
    if p == 'AGENTS.md':   # 根入口已在上文显式比对，跳过避免重复输出
        continue
    sib = os.path.join(os.path.dirname(p), 'CLAUDE.md')
    if not os.path.exists(sib):
        bad(f"缺少 {sib}")
        continue
    same = open(p, 'rb').read() == open(sib, 'rb').read()
    ok(f"{os.path.dirname(p) or '.'} AGENTS/CLAUDE 一致") if same else bad(f"{sib} 与 {p} 不一致（需同步）")

ag = open('AGENTS.md', encoding='utf-8').read()
listed = set(re.findall(r'\]\((\.docs/[^)]+)\)', ag))
unlisted = sorted(f for f in actual if f.startswith('.docs/') and f not in listed)
if unlisted:
    bad(f"文档地图漏列 {len(unlisted)} 项"); [info(p) for p in unlisted]
else:
    ok("文档地图覆盖全部 .docs 文档")

# 模块入口发现：模块文档对应的目录（或其祖先目录）必须有 AGENTS.md
MODULE_DIRS = {
    'beacon': 'beacon', 'beacon-web': 'beacon/web', 'collia': 'collia',
    'common': 'common', 'common-docker': 'common/docker',
    'common-logger': 'common/logger', 'common-timex': 'common/timex',
}
def has_entry(d):
    cur = d
    while cur and cur not in ('.', '/'):
        if os.path.isfile(os.path.join(cur, 'AGENTS.md')):
            return os.path.join(cur, 'AGENTS.md')
        cur = os.path.dirname(cur)
    return None
missing_entry = []
for name, d in sorted(MODULE_DIRS.items()):
    if not os.path.isdir(d):
        continue
    if not has_entry(d):
        missing_entry.append(f"模块 {name}（{d}/）缺少 AGENTS.md 入口")
if missing_entry:
    bad(f"模块入口缺失 {len(missing_entry)} 处"); [info(x) for x in missing_entry]
else:
    ok("全部模块均有 AI 入口")

# ── 10. 健康计数一致性 ───────────────────────────────────────────────────
print("\n[10] 健康计数一致性")
prefix_actual = {}
for k in ids_all:
    pre = re.match(r'([A-Z]+)', k).group(1)
    prefix_actual[pre] = prefix_actual.get(pre, 0) + 1
for path, pattern, actual_n, label in (
        ('AGENTS.md', r'检查目标：(\d+)', len(registered), '检查目标'),
        ('.docs/project-analysis.md', r'Checked targets: (\d+)', len(registered), '检查目标'),
        ('AGENTS.md', r'编号约束 (\d+) 条', len(ids_all), '编号约束'),
        ('.docs/project-analysis.md', r'Domain constraint IDs: (\d+)', len(ids_all), '编号约束')):
    if not os.path.exists(path):
        continue
    txt = open(path, encoding='utf-8').read()
    m = re.search(pattern, txt)
    if not m:
        info(f"{path} 未声明「{label}」计数，跳过")
    elif int(m.group(1)) != actual_n:
        bad(f"{path} 的{label} {m.group(1)} 与实际 {actual_n} 不符")
    else:
        ok(f"{path} 的{label}与实际一致（{actual_n}）")
# 前缀分解校验：project-analysis 若声明了各前缀条数，须与实际一致（防止把文件总数贴成单前缀数）。
pa_txt = open('.docs/project-analysis.md', encoding='utf-8').read() if os.path.exists('.docs/project-analysis.md') else ''
m = re.search(r'Domain constraint IDs: \d+（(.+?)）', pa_txt)
if m:
    declared_pairs = dict((a, int(b)) for a, b in re.findall(r'([A-Z]+)\s*(\d+)', m.group(1)))
    wrong = {p: (n, prefix_actual.get(p)) for p, n in declared_pairs.items() if prefix_actual.get(p) != n}
    missing = sorted(set(prefix_actual) - set(declared_pairs))
    if wrong or missing:
        for p, (n, act) in wrong.items():
            bad(f"约束前缀 {p} 声明 {n} 条，实际 {act} 条")
        if missing:
            bad(f"约束前缀分解漏列：{missing}")
    else:
        ok(f"约束前缀分解与实际一致（{len(declared_pairs)} 个前缀）")
else:
    info("project-analysis 未声明约束前缀分解，跳过")

# ── 11. 开发命令有效性 ──────────────────────────────────────────────────
print("\n[11] 开发命令有效性")
ag_txt = ag
blk = re.search(r'## 开发命令.*?```bash\n(.*?)```', ag_txt, re.S)
cmds = [c.strip() for c in blk.group(1).strip().split('\n')] if blk else []
if not cmds:
    bad("AGENTS.md 未声明开发命令块")
else:
    task_ok, script_ok, script_bad, unchecked = 0, 0, [], []
    task_available = subprocess.run(['bash', '-c', 'command -v task >/dev/null'], capture_output=True).returncode == 0
    if task_available:
        out = subprocess.run(['task', '--list'], capture_output=True, text=True).stdout
        avail = set()
        for line in out.split('\n'):
            line = line.strip()
            if not line.startswith('*'):
                continue
            m = re.match(r'^\*\s+(\S+):\s+(.*)$', line)
            if not m:
                continue
            avail.add(m.group(1))
            al = re.search(r'\(aliases: ([^)]+)\)', m.group(2))
            if al:
                avail.update(a.strip() for a in al.group(1).split(','))
        task_cmds = [c for c in cmds if c.startswith('task ')]
        for c in task_cmds:
            if c.split()[1] not in avail:
                script_bad.append(f"无效 task 目标：{c}")
        task_ok = len(task_cmds) - sum(1 for c in task_cmds if c.split()[1] not in avail)
    else:
        bad("缺少必需依赖 task")
    for c in cmds:
        if c.startswith('task '):
            continue
        m = re.match(r'^cd (\S+) && (?:npm|pnpm|yarn) run (\S+)$', c)
        if m:
            d, s = m.group(1), m.group(2)
            pj = os.path.join(d, 'package.json')
            if not os.path.isdir(d):
                script_bad.append(f"目录不存在：{d}")
            elif not os.path.isfile(pj):
                script_bad.append(f"缺少 package.json：{d}")
            else:
                try:
                    scripts = json.load(open(pj, encoding='utf-8')).get('scripts', {})
                except Exception:
                    scripts = {}
                if s in scripts:
                    script_ok += 1
                else:
                    script_bad.append(f"package.json 无该 script：{d} -> {s}")
            continue
        unchecked.append(c)
    if script_bad:
        bad(f"无效开发命令 {len(script_bad)} 条"); [info(x) for x in script_bad]
    elif not task_available:
        ok(f"{script_ok} 条脚本命令有效（task 未安装，目标校验跳过）")
    else:
        ok(f"{task_ok} 条 task 目标与 {script_ok} 条脚本命令全部有效（目录与 script 已核对）")
    if unchecked:
        info(f"{len(unchecked)} 条命令未纳入校验：{unchecked}")

# Module verification blocks: maintain cwd and validate package-manager shorthand.
print("\n[11b] 模块验证命令")
# Logical shell commands (including continuations) and their execution snapshots
# are checked separately; this does not execute dev servers or Docker builds.
result = subprocess.run(['ruby', 'scripts/doc-evidence-check.rb'], capture_output=True, text=True, env=RUBY_ENV)
if result.returncode:
    bad("证据映射、完整测试清单或模块命令验收记录不一致")
    info(result.stdout + result.stderr)
else:
    ok("证据台账结构与源码双向一致（不代表行为全部验收）")
    info(result.stdout.strip())
# ── 12. 设计 token 登记一致性（DESIGN.md 内部对账）─────────────────────
# 控制台设计稿维护在本仓之外，因此本检查只核对 DESIGN.md 自身的 token 登记：
# §5.1 段内出现的 token 名去重后必须与该段声明的数量一致。
print("\n[12] 设计 token 登记一致性")
token_bad = []
try:
    design_lines = open('DESIGN.md', encoding='utf-8').read().split('\n')
    start = next((i for i, l in enumerate(design_lines) if l.startswith('### 5.1')), None)
    end = next((i for i, l in enumerate(design_lines[start + 1:], start + 1) if l.startswith('### 5.2')), None) if start is not None else None
    if start is None or end is None:
        token_bad.append("DESIGN.md 缺少控制台 token 段（### 5.1 ~ ### 5.2）")
    else:
        seg = '\n'.join(design_lines[start:end])
        families = ('surface', 'foreground', 'border', 'accent', 'chart', 'font', 'spacing', 'radius', 'brand', 'overlay')
        names = set(n for n in re.findall(r'`([a-z][a-z0-9]*(?:-[a-z0-9]+)+)`', seg)
                    if n.split('-')[0] in families)
        declared = re.search(r'### 5\.1[^\n]*（(\d+)\s*个', seg)
        if not declared:
            token_bad.append("DESIGN.md §5.1 未声明 token 数量")
        elif int(declared.group(1)) != len(names):
            token_bad.append(f"DESIGN.md §5.1 token 计数声明 {declared.group(1)} 与段内实际 {len(names)} 不符")
except Exception as exc:
    token_bad.append(f"解析 DESIGN.md 失败：{exc}")
if token_bad:
    bad(f"设计 token 登记不一致 {len(token_bad)} 项"); [info(x) for x in token_bad]
else:
    ok("DESIGN.md §5.1 token 登记与声明数量一致")

# ── 13. OpenAPI 契约与路由注册对账 ──────────────────────────────────────
print("\n[13] OpenAPI 契约与路由注册对账")
HTTP_METHODS = {'get', 'post', 'put', 'patch', 'delete', 'head', 'options'}
oapi = load_yaml('.docs/api/openapi.yml')
declared = {(method, path) for path, item in oapi['paths'].items()
            for method in item if method in HTTP_METHODS}
import tempfile
with tempfile.TemporaryDirectory() as directory:
    output = os.path.join(directory, 'routes.json')
    env = dict(os.environ, BEACON_DOC_ROUTES=output)
    result = subprocess.run(['go', 'test', './service', '-run',
        '^TestDocumentationRoutes$', '-count=1'], cwd='beacon', env=env,
        capture_output=True, text=True)
    if result.returncode:
        bad("运行时路由导出失败"); info(result.stdout + result.stderr)
        implemented = set()
    else:
        implemented = {(method.lower(), re.sub(r':(\w+)', r'{\1}', path))
                       for method, path in json.load(open(output))}
for method, route in sorted(implemented - declared):
    bad(f"实现已注册但契约未声明：{method.upper()} {route}")
for method, route in sorted(declared - implemented):
    bad(f"契约已声明但实现未注册：{method.upper()} {route}")
if declared == implemented:
    ok(f"{len(implemented)} 个运行时 HTTP 方法与路径组合双向一致")

result = subprocess.run(['ruby', 'scripts/doc-contract-check.rb'], capture_output=True, text=True, env=RUBY_ENV)
if result.returncode:
    bad("OpenAPI 请求、响应或 Go 字段模型不一致")
    info(result.stdout + result.stderr)
else:
    ok("OpenAPI 请求/响应建模检查通过")
    info(result.stdout.strip())

print("\n[14] 源码库存事实")

tracked = subprocess.check_output(['git', 'ls-files', '--cached', '--others', '--exclude-standard'], text=True).splitlines()
tracked = sorted({f for f in tracked if os.path.isfile(f)})
tests = [f for f in tracked if f.endswith(('_test.go', '.test.ts', '.spec.ts'))]
pa = open('.docs/project-analysis.md', encoding='utf-8').read()
m = re.search(r'Test files discovered: (\d+)', pa)
if not m or int(m[1]) != len(tests):
    bad(f"测试库存应为 {len(tests)}")
else:
    ok(f"测试文件库存 {len(tests)} 与 project-analysis 声明一致")

GO_MODULES = {
    'beacon': 'beacon', 'collia': 'collia', 'common': 'common',
    'common-docker': 'common/docker', 'common-logger': 'common/logger',
    'common-timex': 'common/timex',
}
for name, module in sorted(GO_MODULES.items()):
    body = open(f'.docs/modules/{name}.md', encoding='utf-8').read()
    files = [f for f in tracked if f.startswith(module + '/') and f.endswith('.go')
             and '/web/' not in f and not any(f.startswith(module + '/' + child + '/')
             for child in ('docker', 'logger', 'timex'))]
    for label, count in [('Go 源码文件', sum(not f.endswith('_test.go') for f in files)),
                         ('Go 测试文件', sum(f.endswith('_test.go') for f in files))]:
        match = re.search(label + r'：(\d+)', body)
        if not match:
            bad(f".docs/modules/{name}.md 未声明「{label}」")
        elif int(match[1]) != count:
            bad(f".docs/modules/{name}.md {label} 应为 {count}")
        else:
            ok(f"{name} {label} {count} 与实际一致")

FRONTEND_MODULES = {'beacon-web': 'beacon/web/src'}
for name, root in sorted(FRONTEND_MODULES.items()):
    body = open(f'.docs/modules/{name}.md', encoding='utf-8').read()
    files = [f for f in tracked if f.startswith(root + '/') and f.endswith(('.vue', '.ts'))]
    tcount = sum(f.endswith(('.test.ts', '.spec.ts')) for f in files)
    claims = re.search(r'前端源码文件：(\d+)（[^）]*?含 (\d+) 个[^）]*?非测试 (\d+)）', body)
    if not claims:
        bad(f".docs/modules/{name}.md 未按「前端源码文件：N（…含 M 个…非测试 K）」声明计数")
    elif tuple(int(x) for x in claims.groups()) != (len(files), tcount, len(files) - tcount):
        bad(f".docs/modules/{name}.md 前端源码计数应为 {len(files)}/{tcount}/{len(files) - tcount}")
    else:
        ok(f"{name} 前端源码 {len(files)}（测试 {tcount}）与实际一致")

for pkg, name in (('beacon/web', 'beacon-web'),):
    pj = json.load(open(os.path.join(pkg, 'package.json'), encoding='utf-8'))
    body = open(f'.docs/modules/{name}.md', encoding='utf-8').read()
    dep_sec = re.search(r'## 依赖关系\n(.*?)(?=\n## |\Z)', body, re.S)
    listed = set(re.findall(r'^- `([^`]+)`', dep_sec[1], re.M)) if dep_sec else set()
    missing = sorted(set(pj['dependencies']) - listed)
    if missing:
        bad(f".docs/modules/{name}.md 依赖清单遗漏运行时依赖：{missing}")
    else:
        ok(f"{name} 依赖清单覆盖全部 {len(pj['dependencies'])} 个运行时依赖")

print("\n[15] 部署配置清单覆盖")
dep = open('.docs/deployment.md', encoding='utf-8').read()
dep_sec = re.search(r'## 配置说明\n(.*?)(?=\n## |\Z)', dep, re.S)
dep_sec = dep_sec[1] if dep_sec else ''
CONFIG_SCAN = [
    ('beacon/configs', ('*.toml', '*.conf')),
    ('beacon/nginx', ('*.conf', '*.types')),
    ('beacon/nginx/conf.d', ('*.conf',)),
    ('beacon/supervisor', ('*.conf',)),
    ('collia', ('config.yml',)),
    ('beacon/web', ('tsconfig.json',)),
]
unlisted = []
for d, pats in CONFIG_SCAN:
    if not os.path.isdir(d):
        continue
    for f in sorted(os.listdir(d)):
        p = os.path.join(d, f)
        if not os.path.isfile(p) or not any(fnmatch.fnmatch(f, pat) for pat in pats):
            continue
        if p not in dep_sec:
            unlisted.append(p)
if unlisted:
    bad(f"运行时配置文件未在 deployment.md 配置说明中登记 {len(unlisted)} 项")
    for p in unlisted:
        info(p)
else:
    ok("全部运行时配置文件已在 deployment.md 配置说明中登记")
phantom = [p for p in re.findall(r'^\| `([^`]+)` \|', dep_sec, re.M)
           if '/' in p and not os.path.exists(p)]
if phantom:
    bad(f"deployment.md 登记了不存在的配置文件 {len(phantom)} 项")
    for p in phantom:
        info(p)
else:
    ok("deployment.md 配置清单中的文件均真实存在")

print("\n[16] 模块文档信号节")
def section_lines(body, name):
    mm = re.search(r'^## ' + re.escape(name) + r'\n(.*?)(?=^## |\Z)', body, re.M | re.S)
    return [l for l in mm[1].split('\n') if l.startswith('- ')] if mm else []
signal_issues = []
for doc in sorted(glob.glob('.docs/modules/*.md')):
    body = open(doc, encoding='utf-8').read()
    sig, cfg = section_lines(body, '接口与符号信号'), section_lines(body, '配置与入口信号')
    def keys(items):
        out = set()
        for l in items:
            mm = re.search(r'`([^`]+)`', l)
            if mm:
                out.add(mm.group(1))
        return out
    ks, kc = keys(sig), keys(cfg)
    if ks and kc and len(ks & kc) / max(len(ks), len(kc)) > 0.5:
        signal_issues.append(f"{doc}「接口与符号信号」与「配置与入口信号」重复承载 {len(ks & kc)} 个相同对象")
    for l in sig + cfg:
        mm = re.search(r'`([^`]+)`', l)
        if mm and mm.group(1).endswith(('.test.ts', '.spec.ts', '.test.go')):
            signal_issues.append(f"{doc} 把测试文件当作入口信号：{mm.group(1)}")
if signal_issues:
    bad(f"模块文档信号节问题 {len(signal_issues)} 处")
    for x in signal_issues:
        info(x)
else:
    ok("模块文档信号节无重复承载，且不含测试文件")

print("\n[17] Domain Spec 元数据与关系图")
PLACEHOLDERS = {'local workspace facts', 'n/a', 'tbd', 'todo', 'unknown', '-', 'none'}
spec_files = sorted(glob.glob('.specs/domain/*.md'))
meta_issues = []
relations = {}
for f in spec_files:
    txt = open(f, encoding='utf-8').read()
    base = os.path.basename(f)
    src = re.search(r'^> 反推来源:\s*(.+)$', txt, re.M)
    if not src:
        meta_issues.append(f"{base} 缺少「反推来源」元数据")
    elif src[1].strip().lower() in PLACEHOLDERS or len(src[1].strip()) < 12:
        meta_issues.append(f"{base} 反推来源为占位符：{src[1].strip()}")
    if not re.search(r'^> 最后更新:\s*\S+', txt, re.M):
        meta_issues.append(f"{base} 缺少「最后更新」元数据")
    sec = re.search(r'## 与其他领域概念的关系\n(.*?)(?=\n## |\Z)', txt, re.S)
    relations[base] = set(re.findall(r'\]\(([^)]+\.md)\)', sec[1])) if sec else set()
if meta_issues:
    bad(f"Domain Spec 元数据问题 {len(meta_issues)} 处")
    for x in meta_issues:
        info(x)
else:
    ok(f"{len(spec_files)} 份 Spec 元数据齐全（反推来源非占位符、含最后更新）")
asym = []
for a, links in sorted(relations.items()):
    for b in sorted(links):
        if a not in relations.get(b, set()):
            asym.append(f"{a} 引用 {b}，但 {b} 未回链 {a}")
if asym:
    bad(f"Domain Spec 关系图单向引用 {len(asym)} 处")
    for x in asym:
        info(x)
else:
    ok("Domain Spec 关系图引用双向对称")

print("\n[18] 门禁自测")
result = subprocess.run(['ruby', 'scripts/doc-check-test.rb'], capture_output=True, text=True, env=RUBY_ENV)
if result.returncode:
    bad("门禁自测未通过：门禁拒绝坏输入的能力未被证明")
    info(result.stdout + result.stderr)
else:
    ok("门禁自测通过（门禁对坏输入的拒绝能力已验证）")
    tail = [line for line in result.stdout.strip().splitlines() if line.strip()]
    if tail:
        info(tail[-1])

# ── 汇总 ────────────────────────────────────────────────────────────────
print("\n" + "=" * 52)
if FAIL:
    print(f"结果：FAIL（{len(FAIL)} 项失败）")
    sys.exit(1)
print("结果：PASS（全部检查通过）")
PYEOF

exit $?
