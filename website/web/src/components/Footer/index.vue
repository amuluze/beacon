<script setup lang="ts">
const logoSrc = '/beacon.svg'
</script>

<template>
    <footer class="site-footer">
        <div class="site-container site-footer__inner">
            <div class="site-footer__columns">
                <div class="site-footer__brand">
                    <div class="site-footer__brand-heading">
                        <img class="site-footer__brand-logo" :src="logoSrc" alt="" aria-hidden="true">
                        <strong>Beacon</strong>
                    </div>
                    <p>轻量级主机监控及 Docker 容器管理工具，实时掌控服务器与容器资源。</p>
                </div>
                <nav class="site-footer__column" aria-label="产品资源">
                    <strong>产品资源</strong>
                    <NuxtLink to="/document">
                        使用手册
                    </NuxtLink>
                    <NuxtLink to="/changelog">
                        更新日志
                    </NuxtLink>
                </nav>
                <nav class="site-footer__column" aria-label="相关作品">
                    <strong>相关作品</strong>
                    <a href="https://www.amuluze.com" target="_blank" rel="noopener noreferrer">amuluze</a>
                    <a href="https://amicon.amuluze.com" target="_blank" rel="noopener noreferrer">amicon</a>
                    <a href="https://www.coverpage.one" target="_blank" rel="noopener noreferrer">coverpage</a>
                </nav>
                <nav class="site-footer__column" aria-label="关注我们">
                    <strong>关注我们</strong>
                    <NuxtLink to="/wechat">
                        微信公众号
                    </NuxtLink>
                    <a href="https://github.com/amuluze/beacon" target="_blank" rel="noopener noreferrer">GitHub</a>
                </nav>
            </div>
            <div class="site-footer__divider" />
            <div class="site-footer__bottom">
                <span>Copyright © 2024-present 阿慕 | 京ICP备2023024485号 | 京公网安备11011402053925号</span>
                <NuxtLink to="/privacy">
                    隐私说明
                </NuxtLink>
            </div>
        </div>
    </footer>
</template>

<style scoped lang="scss">
.site-footer {
  padding: 48px 0 28px;
  background: var(--background);
  border-top: 1px solid var(--border);
}

.site-footer__inner {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.site-footer__columns {
  display: flex;
  justify-content: space-between;
  gap: var(--spacing-xl);
}

.site-footer__brand {
  display: flex;
  flex: 0 0 320px;
  flex-direction: column;
  gap: 8px;
}

.site-footer__brand-heading {
  display: flex;
  align-items: center;
  gap: 10px;
}

.site-footer__brand-logo {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

.site-footer__brand-heading strong {
  color: var(--foreground);
  font-size: 24px;
  font-weight: 700;
}

.site-footer__brand p {
  max-width: 320px;
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: 1.6;
}

.site-footer__column {
  display: flex;
  flex: 0 0 160px;
  flex-direction: column;
  gap: 14px;
}

.site-footer__column strong {
  color: var(--foreground);
  font-size: var(--font-size-md);
  font-weight: 600;
}

.site-footer__column a {
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  transition: color 0.2s ease;
}

.site-footer__column a:hover {
  color: var(--primary);
}

.site-footer__divider {
  height: 1px;
  background: var(--border);
}

.site-footer__bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-2);
  color: var(--muted-foreground);
  font-size: var(--font-size-xs);
}

.site-footer__bottom a {
  color: var(--color-text-secondary);
  transition: color 0.2s ease;
}

.site-footer__bottom a:hover {
  color: var(--primary);
}

@media (max-width: 960px) {
  .site-footer__columns {
    flex-wrap: wrap;
    gap: var(--space-8);
  }

  .site-footer__brand {
    flex: 1 1 100%;
  }

  .site-footer__column {
    flex: 1 1 160px;
  }
}

@media (max-width: 640px) {
  .site-footer__columns {
    gap: var(--space-6);
  }

  .site-footer__bottom {
    align-items: flex-start;
    flex-direction: column;
    gap: var(--space-3);
  }
}
</style>
