<script setup lang="ts">
import { usePageSeo } from '~/composables/usePageSeo'

const faq = [
    {
        question: '安装前需要准备什么环境？',
        answer: '操作系统需为 Linux（x86_64 或 arm64），Docker 版本 20.10.14 以上，Docker Compose 2.0.0 以上，最低 1 核 CPU / 2 GB 内存 / 5 GB 磁盘。',
    },
    {
        question: '如何更新到最新版本？',
        answer: '执行一键升级命令：bash -c "$(curl -fsSLk https://help.beacon.amuluze.com/release/latest/manager.sh)" -- upgrade。升级前会自动备份 .env 与 compose.yaml，健康检查失败会自动回滚。',
    },
    {
        question: '初始账号密码是什么？',
        answer: '初始管理员账号 admin / admin123。默认密码已对外公开、极不安全，首次登录后请立即在「设置」中修改。',
    },
    {
        question: '如何获取技术支持？',
        answer: '通过 GitHub Issues 反馈需求与问题，或邮件联系 314901758@qq.com，也可关注公众号获取更新动态。',
    },
]

const envRequirements = [
    '操作系统：Linux',
    'CPU 指令架构：x86_64, arm64',
    '软件依赖：Docker 20.10.14 版本以上',
    '软件依赖：Docker Compose 2.0.0 版本以上',
    '最低资源需求：1 核 CPU / 2 GB 内存 / 5 GB 磁盘',
]

const quickInstallLines = [
    '# 一键安装（推荐新手使用）',
    'bash -c "$(curl -fsSLk https://help.beacon.amuluze.com/release/latest/manager.sh)"',
    '',
    '# 非交互安装示例',
    'bash -c "$(curl -fsSLk https://help.beacon.amuluze.com/release/latest/manager.sh)" -- install BEACON_HTTP_PORT=1443 BEACON_PUBLIC_BASE_URL=https://beacon.example.com',
    '',
    '# 一键升级',
    'bash -c "$(curl -fsSLk https://help.beacon.amuluze.com/release/latest/manager.sh)" -- upgrade',
    '',
    '# 一键卸载',
    'bash -c "$(curl -fsSLk https://help.beacon.amuluze.com/release/latest/manager.sh)" -- uninstall',
]

const manualInstallLines = [
    'mkdir -p /data/beacon && cd /data/beacon',
    'curl -fsSLO https://help.beacon.amuluze.com/release/latest/compose.yaml',
    'curl -fsSLO https://help.beacon.amuluze.com/release/latest/SHA256SUMS',
    '# 使用 openssl rand -hex 32 生成三个独立随机密钥',
    'cat > .env <<EOF',
    'BEACON_IMAGE=registry.cn-hangzhou.aliyuncs.com/amuluze/beacon:v3.0.4',
    'BEACON_VERSION=v3.0.4',
    'BEACON_CONTAINER_NAME=beacon',
    'BEACON_HTTP_PORT=1443',
    'BEACON_CONTROL_PORT=17000',
    'BEACON_DATA_DIR=./data',
    'BEACON_LOG_DIR=./logs',
    'BEACON_DB_NAME=/app/data/beacon',
    'BEACON_PUBLIC_BASE_URL=http://127.0.0.1:1443',
    'BEACON_AGENT_INSTALL_TOKEN=<your-token>',
    'BEACON_AUTH_SIGNING_KEY=<your-key>',
    'BEACON_CONTROL_JOIN_TOKEN=<your-token>',
    'EOF',
]

const opsLines = [
    'cd /data/beacon',
    'docker compose up -d        # 启动 / 更新',
    'docker compose ps           # 查看状态',
    'docker compose logs -f      # 查看日志',
    'docker compose down         # 停止服务',
]

usePageSeo({
    title: 'Beacon 使用手册',
    description: '使用安全的一键安装流程部署 Beacon，并了解备份、升级与日常运维方式。',
    path: '/document',
})
</script>

<template>
    <div class="docs-page">
        <header class="docs-page__header">
            <div class="site-container docs-page__header-inner">
                <p class="docs-page__overline">
                    <Icon name="lucide:book-open" />
                    <span>使用手册</span>
                </p>
                <h1>Beacon 使用手册</h1>
                <p>支持一键安装、手动安装与离线安装，几分钟完成部署，开箱即用。</p>
            </div>
        </header>

        <div class="site-container docs-page__main">
            <aside class="docs-page__toc">
                <strong>目录</strong>
                <a href="#install">安装使用</a>
                <a href="#faq">常见问题</a>
                <a href="#support">技术支持</a>
            </aside>

            <div class="docs-page__content">
                <section id="install" class="docs-page__section">
                    <header class="docs-page__sec-head">
                        <span class="docs-page__sec-num">01</span>
                        <h2>安装使用</h2>
                    </header>

                    <div class="docs-page__group">
                        <h3 class="docs-page__h3">
                            环境依赖
                        </h3>
                        <p class="docs-page__text">
                            安装 Beacon 前请确保你的系统环境符合以下要求：
                        </p>
                        <ul class="docs-page__list">
                            <li v-for="item in envRequirements" :key="item">
                                {{ item }}
                            </li>
                        </ul>
                    </div>

                    <p class="docs-page__text">
                        Beacon 支持一键安装与手动安装。安装前请确认系统满足环境依赖要求。
                    </p>

                    <div class="docs-page__group">
                        <h3 class="docs-page__h3">
                            一键安装
                        </h3>
                        <div class="docs-page__code">
                            <header class="docs-page__code-head">
                                <span class="docs-page__code-dot" style="background: var(--color-error)" />
                                <span class="docs-page__code-dot" style="background: var(--color-warning)" />
                                <span class="docs-page__code-dot" style="background: var(--color-success)" />
                                <span class="docs-page__code-name">bash</span>
                            </header>
                            <pre class="docs-page__code-body"><code>{{ quickInstallLines.join('\n') }}</code></pre>
                        </div>
                        <p class="docs-page__text">
                            离线环境安装：如果你的服务器无法连接互联网，可以下载离线安装包（镜像包 + compose.yaml + .env 模板）后手动加载并启动。
                        </p>
                    </div>

                    <div class="docs-page__group">
                        <h3 class="docs-page__h3">
                            手动安装
                        </h3>
                        <div class="docs-page__code">
                            <header class="docs-page__code-head">
                                <span class="docs-page__code-dot" style="background: var(--color-error)" />
                                <span class="docs-page__code-dot" style="background: var(--color-warning)" />
                                <span class="docs-page__code-dot" style="background: var(--color-success)" />
                                <span class="docs-page__code-name">bash</span>
                            </header>
                            <pre class="docs-page__code-body"><code>{{ manualInstallLines.join('\n') }}</code></pre>
                        </div>
                    </div>

                    <div class="docs-page__group">
                        <h3 class="docs-page__h3">
                            运维命令
                        </h3>
                        <div class="docs-page__code">
                            <header class="docs-page__code-head">
                                <span class="docs-page__code-dot" style="background: var(--color-error)" />
                                <span class="docs-page__code-dot" style="background: var(--color-warning)" />
                                <span class="docs-page__code-dot" style="background: var(--color-success)" />
                                <span class="docs-page__code-name">bash</span>
                            </header>
                            <pre class="docs-page__code-body"><code>{{ opsLines.join('\n') }}</code></pre>
                        </div>
                    </div>

                    <div class="docs-page__callout">
                        <Icon name="lucide:globe" />
                        <p>启动 Beacon 容器后，浏览器访问 http://服务器IP:1443 即可进入管理面板。初始管理员账号 admin / admin123；默认密码已公开，首次登录后请立即在「设置」中修改。</p>
                    </div>
                </section>

                <section id="faq" class="docs-page__section">
                    <header class="docs-page__sec-head">
                        <span class="docs-page__sec-num">02</span>
                        <h2>常见问题</h2>
                    </header>
                    <p class="docs-page__text">
                        部署与使用中的高频问题解答。
                    </p>
                    <div class="docs-page__faq">
                        <article v-for="item in faq" :key="item.question" class="docs-page__faq-item">
                            <div class="docs-page__qa-row">
                                <span class="docs-page__badge docs-page__badge--q">Q</span>
                                <h3>{{ item.question }}</h3>
                            </div>
                            <div class="docs-page__qa-row">
                                <span class="docs-page__badge docs-page__badge--a">A</span>
                                <p>{{ item.answer }}</p>
                            </div>
                        </article>
                    </div>
                </section>

                <section id="support" class="docs-page__support">
                    <Icon name="lucide:life-buoy" />
                    <div>
                        <h2>技术支持</h2>
                        <p>在 GitHub Issues 描述问题、版本与复现步骤，我们会尽快跟进。</p>
                    </div>
                    <a href="https://github.com/amuluze/beacon/issues" target="_blank" rel="noopener noreferrer">提交 Issue</a>
                </section>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">
.docs-page {
  background: var(--background);
}

.docs-page__header {
  padding: 88px 0 40px;
  background: var(--background);
}

.docs-page__header-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-4);
  text-align: center;
}

.docs-page__overline {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  margin: 0;
  color: var(--primary);
  font-size: var(--font-size-sm);
  font-weight: 600;
  letter-spacing: 1px;
}

.docs-page__overline :deep(svg) {
  width: 16px;
  height: 16px;
}

.docs-page__header h1 {
  margin: 0;
  color: var(--foreground);
  font-size: 44px;
  font-weight: 800;
  letter-spacing: -1px;
  line-height: 1.1;
}

.docs-page__header p {
  max-width: 520px;
  margin: 0;
  color: var(--muted-foreground);
  font-size: var(--font-size-md);
  line-height: 1.6;
}

.docs-page__main {
  display: flex;
  justify-content: center;
  gap: var(--space-12);
  padding: var(--space-8) 0 96px;
}

.docs-page__toc {
  position: sticky;
  top: 88px;
  display: flex;
  flex: 0 0 220px;
  flex-direction: column;
  gap: 6px;
  align-self: flex-start;
}

.docs-page__toc strong {
  margin: 0 0 var(--space-2);
  color: var(--muted-foreground);
  font-size: var(--font-size-xs);
  font-weight: 600;
  letter-spacing: 1px;
}

.docs-page__toc a {
  display: flex;
  align-items: center;
  padding: 10px 14px;
  color: var(--color-text-secondary);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
  font-weight: 500;
  transition:
    color 0.2s ease,
    background 0.2s ease;
}

.docs-page__toc a:hover {
  color: var(--primary);
  background: var(--color-bg-hover);
}

.docs-page__toc a.router-link-active,
.docs-page__toc a[aria-current='true'] {
  color: var(--primary);
  background: var(--color-surface-muted);
  font-weight: 600;
}

.docs-page__content {
  display: flex;
  flex: 0 1 800px;
  flex-direction: column;
  gap: var(--space-12);
  min-width: 0;
}

.docs-page__section {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
  scroll-margin-top: 88px;
}

.docs-page__sec-head {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.docs-page__sec-num {
  display: grid;
  flex: 0 0 auto;
  place-items: center;
  width: 30px;
  height: 30px;
  color: var(--foreground);
  background: var(--primary);
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
  font-size: var(--font-size-sm);
  font-weight: 700;
}

.docs-page__sec-head h2 {
  margin: 0;
  color: var(--foreground);
  font-size: var(--font-size-xl);
  font-weight: 700;
  line-height: 1.2;
}

.docs-page__group {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.docs-page__h3 {
  margin: var(--space-4) 0 0;
  color: var(--foreground);
  font-size: 17px;
  font-weight: 600;
}

.docs-page__section > .docs-page__group:first-child .docs-page__h3,
.docs-page__group:first-child .docs-page__h3 {
  margin-top: 0;
}

.docs-page__text {
  margin: 0;
  color: var(--muted-foreground);
  font-size: var(--font-size-sm);
  line-height: 1.7;
}

.docs-page__list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin: 0;
  padding: 0;
  color: var(--muted-foreground);
  font-size: var(--font-size-sm);
  line-height: 1.7;
  list-style: none;
}

.docs-page__list li {
  padding-left: 12px;
  border-left: 2px solid var(--border);
}

.docs-page__code {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--color-surface-muted);
  border: 1px solid var(--border);
  border-radius: var(--radius-m);
}

.docs-page__code-head {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 10px 14px;
  border-bottom: 1px solid var(--border);
}

.docs-page__code-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.docs-page__code-name {
  margin-left: 6px;
  color: var(--muted-foreground);
  font-family: var(--font-mono);
  font-size: var(--font-size-xs);
}

.docs-page__code-body {
  margin: 0;
  padding: 12px 14px;
  overflow-x: auto;
}

.docs-page__code-body code {
  color: var(--color-text-secondary);
  font-family: var(--font-mono);
  font-size: var(--font-size-sm);
  line-height: 1.7;
  white-space: pre;
}

.docs-page__callout {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  padding: var(--space-4);
  background: var(--color-primary-soft);
  border-left: 3px solid var(--primary);
  border-radius: var(--radius-m);
}

.docs-page__callout :deep(svg) {
  flex: 0 0 auto;
  width: 18px;
  height: 18px;
  color: var(--primary);
  margin-top: 2px;
}

.docs-page__callout p {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: 1.6;
}

.docs-page__faq {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.docs-page__faq-item {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: var(--space-5);
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: var(--radius-m);
}

.docs-page__qa-row {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.docs-page__qa-row h3 {
  margin: 0;
  color: var(--foreground);
  font-size: var(--font-size-md);
  font-weight: 600;
}

.docs-page__qa-row p {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: 1.7;
}

.docs-page__badge {
  display: grid;
  flex: 0 0 auto;
  place-items: center;
  width: 26px;
  height: 26px;
  border-radius: var(--radius-sm);
  font-family: var(--font-primary);
  font-size: var(--font-size-sm);
  font-weight: 700;
}

.docs-page__badge--q {
  color: var(--foreground);
  background: var(--primary);
}

.docs-page__badge--a {
  color: var(--primary);
  background: var(--color-surface-muted);
  border: 1px solid var(--border);
}

.docs-page__support {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  gap: var(--space-4);
  align-items: center;
  padding: var(--space-6);
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  scroll-margin-top: 88px;
}

.docs-page__support > :deep(svg) {
  width: 28px;
  height: 28px;
  color: var(--primary);
}

.docs-page__support h2 {
  margin: 0;
  color: var(--foreground);
  font-size: 18px;
  font-weight: 600;
}

.docs-page__support p {
  margin: var(--space-1) 0 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
}

.docs-page__support a {
  display: inline-flex;
  align-items: center;
  padding: 10px 16px;
  color: var(--color-text-inverse);
  background: var(--primary);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
  font-weight: 600;
  transition: background 0.2s ease;
}

.docs-page__support a:hover {
  background: var(--color-brand-hover);
}

@media (max-width: 960px) {
  .docs-page__main {
    flex-direction: column;
    gap: var(--space-8);
  }

  .docs-page__toc {
    position: static;
    flex: 1 1 auto;
    flex-direction: row;
    flex-wrap: wrap;
    align-items: center;
  }

  .docs-page__toc strong {
    width: 100%;
    margin: 0;
  }

  .docs-page__toc a {
    flex: 1 0 auto;
  }
}

@media (max-width: 640px) {
  .docs-page__header {
    padding: 56px 0 var(--space-8);
  }

  .docs-page__header h1 {
    font-size: 32px;
  }

  .docs-page__main {
    padding: var(--space-6) 0 64px;
  }

  .docs-page__support {
    grid-template-columns: auto minmax(0, 1fr);
  }

  .docs-page__support a {
    grid-column: 1 / -1;
    justify-self: stretch;
    justify-content: center;
  }
}
</style>
