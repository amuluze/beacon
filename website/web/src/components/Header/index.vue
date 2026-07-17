<script setup lang="ts">
import { useTheme } from '~/composables/useTheme'

const logoSrc = '/beacon.svg'

const links = [
    { label: '首页', to: '/' },
    { label: '使用手册', to: '/document' },
    { label: '更新日志', to: '/changelog' },
]

const mobileLinks = [
    ...links,
    { label: '微信公众号', to: '/wechat' },
]

const { isDark, toggleTheme } = useTheme()
const isMenuOpen = shallowRef(false)

function closeMenu() {
    isMenuOpen.value = false
}

function toggleMenu() {
    isMenuOpen.value = !isMenuOpen.value
}
</script>

<template>
    <header class="site-header">
        <div class="site-header__inner site-container">
            <NuxtLink to="/" class="site-header__brand" aria-label="Beacon 首页">
                <img class="site-header__mark" :src="logoSrc" alt="" aria-hidden="true">
                <span>Beacon</span>
            </NuxtLink>
            <nav class="site-header__nav site-header__nav--desktop" aria-label="主导航">
                <NuxtLink v-for="link in links" :key="link.to" :to="link.to" class="site-header__link">
                    {{ link.label }}
                </NuxtLink>
                <a class="site-header__link" href="https://github.com/amuluze/beacon" target="_blank" rel="noopener noreferrer">
                    GitHub
                </a>
                <button type="button" class="site-header__toggle" :aria-label="isDark ? '切换到浅色模式' : '切换到暗色模式'" @click="toggleTheme">
                    <Icon :name="isDark ? 'mdi:white-balance-sunny' : 'mdi:weather-night'" />
                </button>
                <NuxtLink to="/document" class="site-button site-button--primary site-header__cta">
                    <Icon name="mdi:rocket-launch-outline" />
                    <span>开始安装</span>
                </NuxtLink>
            </nav>
            <button
                type="button"
                class="site-header__menu-toggle"
                aria-label="打开主导航"
                aria-controls="mobile-navigation"
                :aria-expanded="isMenuOpen"
                @click="toggleMenu"
            >
                <Icon :name="isMenuOpen ? 'mdi:close' : 'mdi:menu'" />
            </button>
        </div>
        <nav v-if="isMenuOpen" id="mobile-navigation" class="site-header__mobile-menu" aria-label="移动端主导航">
            <NuxtLink v-for="link in mobileLinks" :key="link.to" :to="link.to" @click="closeMenu">
                {{ link.label }}
            </NuxtLink>
            <a href="https://github.com/amuluze/beacon" target="_blank" rel="noopener noreferrer" @click="closeMenu">GitHub</a>
            <button type="button" @click="toggleTheme">
                <Icon :name="isDark ? 'mdi:white-balance-sunny' : 'mdi:weather-night'" />
                <span>{{ isDark ? '浅色模式' : '暗色模式' }}</span>
            </button>
        </nav>
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
  width: 30px;
  height: 30px;
  object-fit: contain;
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

.site-header__menu-toggle,
.site-header__mobile-menu {
  display: none;
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
  .site-header__nav--desktop {
    display: none;
  }

  .site-header__menu-toggle {
    display: inline-grid;
    place-items: center;
    width: 40px;
    height: 40px;
    padding: 0;
    color: var(--foreground);
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: 22px;
  }

  .site-header__mobile-menu {
    position: absolute;
    inset: 64px 0 auto;
    display: grid;
    gap: var(--space-1);
    padding: var(--space-3) 16px var(--space-4);
    background: var(--background);
    border-bottom: 1px solid var(--border);
    box-shadow: var(--shadow-card);
  }

  .site-header__mobile-menu a,
  .site-header__mobile-menu button {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    min-height: 44px;
    padding: 0 var(--space-3);
    color: var(--color-text-secondary);
    background: transparent;
    border: 0;
    border-radius: var(--radius-sm);
    font: inherit;
    text-align: left;
  }

  .site-header__mobile-menu a:hover,
  .site-header__mobile-menu a.router-link-active,
  .site-header__mobile-menu button:hover {
    color: var(--primary);
    background: var(--color-bg-hover);
  }
}
</style>
