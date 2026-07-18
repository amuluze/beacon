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
                <div class="site-header__links">
                    <NuxtLink v-for="link in links" :key="link.to" :to="link.to" class="site-header__link">
                        {{ link.label }}
                    </NuxtLink>
                    <a class="site-header__link" href="https://github.com/amuluze/beacon" target="_blank" rel="noopener noreferrer">
                        GitHub
                    </a>
                </div>
                <button type="button" class="site-header__toggle" :aria-label="isDark ? '切换到浅色模式' : '切换到暗色模式'" @click="toggleTheme">
                    <Icon :name="isDark ? 'lucide:sun' : 'lucide:moon'" />
                </button>
                <NuxtLink to="/document" class="site-button site-button--primary site-header__cta">
                    <Icon name="lucide:rocket" />
                    <span>立即体验</span>
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
                <Icon :name="isMenuOpen ? 'lucide:x' : 'lucide:menu'" />
            </button>
        </div>
        <nav v-if="isMenuOpen" id="mobile-navigation" class="site-header__mobile-menu" aria-label="移动端主导航">
            <NuxtLink v-for="link in mobileLinks" :key="link.to" :to="link.to" @click="closeMenu">
                {{ link.label }}
            </NuxtLink>
            <a href="https://github.com/amuluze/beacon" target="_blank" rel="noopener noreferrer" @click="closeMenu">GitHub</a>
            <button type="button" @click="toggleTheme">
                <Icon :name="isDark ? 'lucide:sun' : 'lucide:moon'" />
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
  gap: 10px;
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
}

.site-header__mark {
  width: 26px;
  height: 26px;
  object-fit: contain;
}

.site-header__nav {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.site-header__links {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  margin-right: var(--space-2);
}

.site-header__link {
  padding: 8px 4px;
  color: var(--color-text-secondary);
  font-size: var(--font-size-md);
  font-weight: var(--font-weight-medium);
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
  font-size: 16px;
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

@media (max-width: 640px) {
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
    font-size: 20px;
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
