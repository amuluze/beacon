<script setup lang="ts">
const logoSrc = '/beacon.svg'
</script>

<template>
    <footer class="site-footer">
        <div class="site-container">
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
  color: var(--color-text-secondary);
  background: var(--background);
  border-top: 1px solid var(--border);
}

.site-footer__columns {
  display: grid;
  grid-template-columns: minmax(300px, 2.2fr) repeat(3, minmax(120px, 1fr));
  gap: clamp(32px, 5vw, 72px);
  padding-bottom: var(--space-8);
}

.site-footer__brand strong,
.site-footer__column strong {
  color: var(--foreground);
  font-size: var(--font-size-md);
}

.site-footer__brand-heading {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.site-footer__brand-logo {
  width: 36px;
  height: 36px;
  object-fit: contain;
}

.site-footer__brand p {
  max-width: 400px;
  margin: var(--space-2) 0 0;
}

.site-footer__column {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 10px;
}

.site-footer__column a,
.site-footer__bottom a {
  transition: color 0.2s ease;
}

.site-footer__column a:hover,
.site-footer__bottom a:hover {
  color: var(--primary);
}

.site-footer__bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-2);
  padding-top: var(--space-5);
  border-top: 1px solid var(--border);
  font-size: var(--font-size-xs);
}

@media (max-width: 900px) {
  .site-footer__columns {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .site-footer__brand {
    grid-column: 1 / -1;
  }
}

@media (max-width: 640px) {
  .site-footer {
    padding-top: var(--space-8);
  }

  .site-footer__columns {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: var(--space-6);
  }

  .site-footer__brand {
    grid-column: 1 / -1;
  }

  .site-footer__bottom {
    align-items: flex-start;
    flex-direction: column;
    gap: var(--space-3);
  }
}

@media (max-width: 420px) {
  .site-footer__columns {
    grid-template-columns: 1fr;
  }

  .site-footer__brand {
    grid-column: auto;
  }
}
</style>
