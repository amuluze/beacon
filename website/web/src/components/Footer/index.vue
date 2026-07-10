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
                    <NuxtLink to="/about">团队故事</NuxtLink>
                    <NuxtLink to="/document">使用手册</NuxtLink>
                </div>
                <div class="site-footer__column">
                    <strong>关注我们</strong>
                    <NuxtLink to="/wechat">微信公众号</NuxtLink>
                    <a href="https://github.com/amuluze/amprobe" target="_blank" rel="noopener noreferrer">GitHub</a>
                </div>
            </div>
            <div class="site-footer__links">
                <span>友情链接：</span>
                <a href="https://www.amuluze.com" target="_blank" rel="noopener noreferrer">amuluze</a>
                <span>·</span>
                <a href="https://www.coverpage.one" target="_blank" rel="noopener noreferrer">coverpage</a>
            </div>
            <p class="site-footer__legal">Copyright © 2024-present 阿慕 | 京ICP备2023024485号 | 京公网安备11011402053925号</p>
        </div>
    </footer>
</template>

<style scoped lang="scss">
.site-footer {
  padding: 56px 0 24px;
  color: var(--site-foreground-secondary);
  background: var(--site-surface-secondary);
  border-top: 1px solid var(--site-border-subtle);
}

.site-footer__columns {
  display: grid;
  grid-template-columns: minmax(280px, 2fr) repeat(2, minmax(140px, 1fr));
  gap: 64px;
  padding-bottom: var(--site-space-xl);
}

.site-footer__brand strong,
.site-footer__column strong {
  color: var(--site-foreground-primary);
  font-size: 14px;
}

.site-footer__brand p {
  max-width: 360px;
  margin: var(--site-space-sm) 0 0;
}

.site-footer__column {
  display: flex;
  flex-direction: column;
  gap: var(--site-space-sm);
}

.site-footer__column a:hover,
.site-footer__links a:hover {
  color: var(--site-accent);
}

.site-footer__links,
.site-footer__legal {
  display: flex;
  align-items: center;
  gap: var(--site-space-sm);
  padding-top: var(--site-space-md);
  border-top: 1px solid var(--site-border-primary);
  font-size: 12px;
}

.site-footer__legal {
  margin: var(--site-space-md) 0 0;
}

@media (max-width: 640px) {
  .site-footer {
    padding-top: var(--site-space-xl);
  }

  .site-footer__columns {
    grid-template-columns: 1fr;
    gap: var(--site-space-lg);
  }

  .site-footer__links {
    flex-wrap: wrap;
  }

  .site-footer__legal {
    display: block;
  }
}
</style>
