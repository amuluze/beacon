<script setup lang="ts">
const faq = [
    {
        question: '如何备份和恢复配置？',
        answer: '备份 /data/beacon 目录（含 .env、compose.yaml 与数据），执行 tar -czf beacon-backup.tar.gz /data/beacon；恢复时解压到原目录并执行 docker compose up -d。',
    },
    {
        question: '如何更新到最新版本？',
        answer: '进入安装目录，建议先备份配置与数据，再执行 docker compose pull && docker compose up -d，即可拉取并应用最新镜像。',
    },
    {
        question: '初始化用户有哪些？密码是什么？',
        answer: '管理员 admin / admin123，普通用户 beacon / 123456；管理员可管理普通用户，请上线后及时修改默认密码。',
    },
    {
        question: '如何获取技术支持？',
        answer: '通过 GitHub Issues 反馈需求与问题，或邮件联系 314901758@qq.com，也可关注公众号获取更新动态。',
    },
]

useHead({ title: 'Beacon 使用手册' })
</script>

<template>
    <div>
        <header class="site-page-header">
            <div class="site-container">
                <p class="site-overline">使用手册</p>
                <h1 class="site-page-title">Beacon 使用手册</h1>
                <p class="site-page-description">基于 Docker Compose，几分钟完成部署，开箱即用</p>
            </div>
        </header>

        <main class="site-container docs">
            <aside class="site-card docs__toc">
                <strong>目录</strong>
                <a href="#install">安装使用</a>
                <a href="#faq">常见问题</a>
                <a href="#support">技术支持</a>
            </aside>

            <div class="docs__content">
                <section id="install" class="docs__section">
                    <header class="docs__section-header">
                        <span>01</span>
                        <div>
                            <h2>安装使用</h2>
                            <p>通过官方安装脚本一键部署 Beacon，或手动使用 Docker Compose 启动。</p>
                        </div>
                    </header>

                    <article class="site-card docs__card">
                        <h3>一键安装</h3>
                        <pre class="site-code"><code># 下载并执行安装脚本
curl -fsSL https://official.beacon.amuluze.com/download/install.sh -o install.sh
sh install.sh

# 非交互安装示例
BEACON_HTTP_PORT=1443 sh install.sh</code></pre>

                        <h3>手动安装</h3>
                        <pre class="site-code"><code>mkdir -p /data/beacon && cd /data/beacon
curl -fsSL https://official.beacon.amuluze.com/download/compose.yaml -o compose.yaml
# 编辑 .env 配置端口 / Token / 密钥
docker compose up -d</code></pre>

                        <h3>运维命令</h3>
                        <pre class="site-code"><code>cd /data/beacon
docker compose up -d        # 启动 / 更新
docker compose ps           # 查看状态
docker compose logs -f      # 查看日志
docker compose down         # 停止服务</code></pre>

                        <p class="docs__callout">启动 beacon 容器后，浏览器访问 <code>http://服务器IP:1443</code> 即可进入管理面板。初始账号 admin / admin123。</p>
                    </article>
                </section>

                <section id="faq" class="docs__section">
                    <header class="docs__section-header">
                        <span>02</span>
                        <div>
                            <h2>常见问题</h2>
                            <p>部署与使用中的高频问题解答。</p>
                        </div>
                    </header>
                    <div class="docs__faq">
                        <article v-for="item in faq" :key="item.question" class="site-card docs__faq-item">
                            <div class="docs__qa">Q</div>
                            <div>
                                <h3>{{ item.question }}</h3>
                                <p>{{ item.answer }}</p>
                            </div>
                        </article>
                    </div>
                </section>

                <section id="support" class="site-card docs__support">
                    <Icon name="mdi:lifebuoy" />
                    <div>
                        <h2>技术支持</h2>
                        <p>在 GitHub Issues 描述问题、版本与复现步骤，我们会尽快跟进。</p>
                    </div>
                    <a href="https://github.com/amuluze/amprobe/issues" target="_blank" rel="noopener noreferrer">提交 Issue</a>
                </section>
            </div>
        </main>
    </div>
</template>

<style scoped lang="scss">
.docs {
  display: grid;
  grid-template-columns: 200px minmax(0, 1fr);
  gap: 48px;
  align-items: start;
  padding-top: 64px;
  padding-bottom: 80px;
}

.docs__toc {
  position: sticky;
  top: 88px;
  display: flex;
  flex-direction: column;
  gap: var(--site-space-sm);
  padding: var(--site-space-md);
}

.docs__toc strong {
  margin-bottom: var(--site-space-xs);
}

.docs__toc a {
  color: var(--site-foreground-secondary);
}

.docs__toc a:hover {
  color: var(--site-accent);
}

.docs__content,
.docs__section,
.docs__faq {
  display: flex;
  flex-direction: column;
}

.docs__content {
  gap: 72px;
}

.docs__section {
  gap: var(--site-space-lg);
  scroll-margin-top: 88px;
}

.docs__section-header {
  display: flex;
  align-items: flex-start;
  gap: var(--site-space-md);
}

.docs__section-header > span {
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  width: 48px;
  height: 48px;
  color: var(--site-on-accent);
  background: var(--site-accent);
  border-radius: 50%;
  font-family: var(--site-font-mono);
  font-weight: 700;
}

.docs__section-header h2,
.docs__section-header p {
  margin: 0;
}

.docs__section-header h2 {
  font-size: 28px;
}

.docs__section-header p {
  color: var(--site-foreground-secondary);
}

.docs__card {
  padding: var(--site-space-xl);
}

.docs__card h3 {
  margin: 32px 0 12px;
}

.docs__card h3:first-child {
  margin-top: 0;
}

.docs__card pre {
  margin: 0;
  padding: var(--site-space-md);
  font-size: 12px;
  line-height: 1.7;
}

.docs__callout {
  margin: var(--site-space-lg) 0 0;
  padding: var(--site-space-md);
  color: var(--site-foreground-secondary);
  background: var(--site-accent-soft);
  border-left: 3px solid var(--site-accent);
}

.docs__faq {
  gap: var(--site-space-md);
}

.docs__faq-item {
  display: grid;
  grid-template-columns: 36px minmax(0, 1fr);
  gap: var(--site-space-md);
  padding: var(--site-space-lg);
}

.docs__qa {
  display: grid;
  place-items: center;
  width: 32px;
  height: 32px;
  color: var(--site-accent);
  background: var(--site-accent-soft);
  border-radius: 50%;
  font-weight: 700;
}

.docs__faq-item h3,
.docs__faq-item p {
  margin: 0;
}

.docs__faq-item h3 {
  font-size: 15px;
}

.docs__faq-item p {
  margin-top: var(--site-space-sm);
  color: var(--site-foreground-secondary);
}

.docs__support {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  gap: var(--site-space-md);
  align-items: center;
  padding: var(--site-space-lg);
  scroll-margin-top: 88px;
}

.docs__support > :deep(svg) {
  color: var(--site-accent);
  font-size: 28px;
}

.docs__support h2,
.docs__support p {
  margin: 0;
}

.docs__support p {
  color: var(--site-foreground-secondary);
}

.docs__support a {
  padding: 8px 12px;
  color: var(--site-on-accent);
  background: var(--site-accent);
  border-radius: var(--site-radius-sm);
}

@media (max-width: 800px) {
  .docs {
    grid-template-columns: 1fr;
    gap: var(--site-space-xl);
    padding-top: 32px;
    padding-bottom: 48px;
  }

  .docs__toc {
    position: static;
    flex-direction: row;
    flex-wrap: wrap;
  }

  .docs__toc strong {
    width: 100%;
  }
}

@media (max-width: 520px) {
  .docs__content {
    gap: 48px;
  }

  .docs__card {
    padding: var(--site-space-md);
  }

  .docs__support {
    grid-template-columns: auto minmax(0, 1fr);
  }

  .docs__support a {
    grid-column: 1 / -1;
    text-align: center;
  }
}
</style>
