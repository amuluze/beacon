<script setup lang="ts">
import { useTheme } from '~/composables/useTheme'

const links = [
    { label: '首页', to: '/' },
    { label: '使用手册', to: '/document' },
]

const { isDark, toggleTheme } = useTheme()
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
                <a class="site-header__link" href="https://github.com/amuluze/beacon" target="_blank" rel="noopener noreferrer">
                    GitHub
                </a>
                <button type="button" class="site-header__toggle" :aria-label="isDark ? '切换到浅色模式' : '切换到暗色模式'" @click="toggleTheme">
                    <Icon :name="isDark ? 'mdi:white-sunny' : 'mdi:weather-night'" />
                </button>
                <NuxtLink to="/document" class="site-button site-button--primary site-header__cta">
                    <Icon name="mdi:rocket-launch-outline" />
                    <span>立即体验</span>
                </NuxtLink>
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
  background: color-mix(in srgb, var(--background) 92%, transparent);
  border-bottom: 1px solid var(--border);
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
  gap: var(--space-2);
  font-size: var(--font-size-lg);
  font-weight: 700;
}

.site-header__mark {
  width: 26px;
  height: 26px;
  background: var(--primary);
  border-radius: var(--radius-sm);
}

.site-header__nav {
  display: flex;
  align-items: center;
  gap: var(--space-8);
}

.site-header__link {
  color: var(--color-text-secondary);
  font-size: var(--font-size-md);
  font-weight: 500;
  transition: color 0.2s ease;
}

.site-header__link:hover,
.site-header__link.router-link-active {
  color: var(--primary);
}

.site-header__toggle {
  display: inline-grid;
  place-items: center;
  width: 36px;
  height: 36px;
  padding: 0;
  color: var(--color-text-secondary);
  background: transparent;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 18px;
  transition:
    color 0.2s ease,
    border-color 0.2s ease;
}

.site-header__toggle:hover {
  color: var(--primary);
  border-color: var(--primary);
}

.site-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: 0 20px;
  font-size: var(--font-size-md);
  font-weight: 600;
  border-radius: var(--radius-sm);
  cursor: pointer;
}

.site-button--primary {
  height: 40px;
  color: var(--color-text-inverse);
  background: var(--primary);
  transition: background 0.2s ease;
}

.site-button--primary:hover {
  background: var(--color-brand-hover);
}

@media (max-width: 720px) {
  .site-header__nav {
    gap: var(--space-4);
  }

  .site-header__link {
    display: none;
  }

  .site-header__cta span {
    display: none;
  }

  .site-header__cta {
    width: 40px;
    padding: 0;
  }
}
</style>
