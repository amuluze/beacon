<template>
    <footer class="site-footer">
        <div class="site-container">
            <div class="site-footer__columns">
                <div class="site-footer__brand">
                    <strong>Beacon</strong>
                    <p>轻量级主机监控及 Docker 容器管理工具，实时掌控服务器与容器资源。</p>
                </div>
                <div class="site-footer__column">
                    <strong>关于我们</strong>
                    <NuxtLink to="/about">
                        团队故事
                    </NuxtLink>
                    <NuxtLink to="/document">
                        使用手册
                    </NuxtLink>
                    <NuxtLink to="/changelog">
                        更新日志
                    </NuxtLink>
                </div>
                <div class="site-footer__column">
                    <strong>关注我们</strong>
                    <NuxtLink to="/wechat">
                        微信公众号
                    </NuxtLink>
                    <a href="https://github.com/amuluze/beacon" target="_blank" rel="noopener noreferrer">GitHub</a>
                </div>
            </div>
            <div class="site-footer__links">
                <span>友情链接：</span>
                <a href="https://www.amuluze.com" target="_blank" rel="noopener noreferrer">amuluze</a>
                <span>·</span>
                <a href="https://www.coverpage.one" target="_blank" rel="noopener noreferrer">coverpage</a>
            </div>
            <p class="site-footer__legal">
                <span>Copyright © 2024-present 阿慕 | 京ICP备2023024485号 | 京公网安备11011402053925号</span>
                <NuxtLink to="/privacy">
                    隐私说明
                </NuxtLink>
            </p>
        </div>
    </footer>
</template>

<style scoped lang="scss">
.site-footer {
  padding: 48px 0 28px;
  color: var(--color-text-secondary);
  background: var(--background);
  border-top: 1px solid var(--border);
}

.site-footer__columns {
  display: grid;
  grid-template-columns: minmax(280px, 2fr) repeat(2, minmax(140px, 1fr));
  gap: 64px;
  padding-bottom: var(--space-8);
}

.site-footer__brand strong,
.site-footer__column strong {
  color: var(--foreground);
  font-size: var(--font-size-md);
}

.site-footer__brand p {
  max-width: 360px;
  margin: var(--space-2) 0 0;
}

.site-footer__column {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.site-footer__column a:hover,
.site-footer__links a:hover {
  color: var(--primary);
}

.site-footer__links,
.site-footer__legal {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding-top: var(--space-4);
  border-top: 1px solid var(--border);
  font-size: var(--font-size-xs);
}

.site-footer__legal {
  justify-content: space-between;
  margin: var(--space-4) 0 0;
}

@media (max-width: 640px) {
  .site-footer {
    padding-top: var(--space-8);
  }

  .site-footer__columns {
    grid-template-columns: 1fr;
    gap: var(--space-6);
  }

  .site-footer__links {
    flex-wrap: wrap;
  }

  .site-footer__legal {
    display: block;
  }
}
</style>
