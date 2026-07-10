<script setup lang="ts">
const links = [
    { label: '首页', to: '/' },
    { label: '使用手册', to: '/document' },
]
</script>

<template>
    <header class="site-header">
        <div class="site-header__inner site-container">
            <NuxtLink to="/" class="site-header__brand" aria-label="Beacon 首页">
                <span class="site-header__mark" />
                <span>Beacon</span>
            </NuxtLink>
            <nav class="site-header__nav" aria-label="主导航">
                <NuxtLink v-for="link in links" :key="link.to" :to="link.to" class="site-header__link">
                    {{ link.label }}
                </NuxtLink>
                <a class="site-header__link" href="https://github.com/amuluze/amprobe" target="_blank" rel="noopener noreferrer">
                    GitHub
                </a>
            </nav>
        </div>
    </header>
</template>

<style scoped lang="scss">
.site-header {
  position: fixed;
  z-index: 100;
  inset: 0 0 auto;
  height: 64px;
  background: rgba(251, 252, 254, 0.94);
  border-bottom: 1px solid var(--site-border-subtle);
  backdrop-filter: blur(12px);
}

.site-header__inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
}

.site-header__brand {
  display: inline-flex;
  align-items: center;
  gap: var(--site-space-sm);
  font-size: 18px;
  font-weight: 700;
}

.site-header__mark {
  width: 24px;
  height: 24px;
  background: var(--site-accent);
  border-radius: 5px;
}

.site-header__nav {
  display: flex;
  align-items: center;
  gap: var(--site-space-xl);
}

.site-header__link {
  color: var(--site-foreground-secondary);
  font-size: 13px;
  font-weight: 500;
  transition: color 0.2s ease;
}

.site-header__link:hover,
.site-header__link.router-link-active {
  color: var(--site-accent);
}

@media (max-width: 640px) {
  .site-header__nav {
    display: none;
  }
}
</style>
