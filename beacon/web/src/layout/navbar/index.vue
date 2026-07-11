<script setup lang="ts">
import Avatar from '@/layout/navbar/Avatar.vue'
import AgentSelect from '@/layout/navbar/AgentSelect.vue'
import InstallAgent from '@/layout/navbar/InstallAgent.vue'
import Language from '@/layout/navbar/Language.vue'
import ThemeChange from '@/layout/navbar/ThemeChange.vue'
import { dynamicRoutes } from '@/router/dynamic.ts'
import type { RouteRecordRaw } from 'vue-router'
import { useI18n } from 'vue-i18n'
import type { Component } from 'vue'
import IconActivity from '~icons/lucide/activity'
import IconPackage from '~icons/lucide/package'
import IconSettings from '~icons/lucide/settings'
import IconTerminal from '~icons/lucide/terminal'
import IconUsers from '~icons/lucide/users'
import IconMenu from '~icons/lucide/menu'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const iconMap: Record<string, Component> = {
  activity: IconActivity,
  package: IconPackage,
  settings: IconSettings,
  terminal: IconTerminal,
  users: IconUsers,
  menu: IconMenu,
}

const visibleRoutes = computed(() => dynamicRoutes.filter(item => item.meta?.show))

function routeTarget(item: RouteRecordRaw): string {
  return typeof item.redirect === 'string' ? item.redirect : item.path
}

function isRouteActive(item: RouteRecordRaw): boolean {
  return route.path === item.path || route.path.startsWith(`${item.path}/`)
}

function menuIcon(item: RouteRecordRaw): Component {
  const name = String(item.meta?.icon || 'menu')
  return iconMap[name] || IconMenu
}

function menuTitle(item: RouteRecordRaw): string {
  return t(String(item.meta?.title || item.name || ''))
}

function goRoute(item: RouteRecordRaw): void {
  router.push(routeTarget(item))
}
</script>

<template>
    <header class="am-navbar">
        <button class="am-navbar__brand" type="button" @click="router.push('/monitor')">
            <span class="am-navbar__brand-mark" />
            <span class="am-navbar__brand-text">Beacon</span>
        </button>

        <nav class="am-navbar__menu" aria-label="Primary">
            <button
                v-for="item in visibleRoutes"
                :key="String(item.name)"
                class="am-navbar__menu-item"
                :class="{ 'am-navbar__menu-item--active': isRouteActive(item) }"
                type="button"
                @click="goRoute(item)"
            >
                <component :is="menuIcon(item)" class="am-navbar__menu-icon" />
                <span>{{ menuTitle(item) }}</span>
            </button>
        </nav>

        <div class="am-navbar__right">
            <AgentSelect />
            <InstallAgent />
            <Language />
            <ThemeChange />
            <Avatar />
        </div>
    </header>
</template>

<style scoped lang="scss">
@include b(navbar) {
  flex: 0 0 var(--am-topbar-height);
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: var(--am-topbar-height);
  width: 100%;
  padding: 0 var(--am-spacing-lg);
  gap: var(--am-spacing-md);
  color: var(--am-foreground-primary);
  background-color: var(--am-surface-secondary);
  border-bottom: 1px solid var(--am-border-primary);

  @include e(brand) {
    min-width: 180px;
    padding: 0;
    display: inline-flex;
    align-items: center;
    gap: var(--am-spacing-sm);
    color: inherit;
    background: transparent;
    border: 0;
    cursor: pointer;
  }

  @include e(brand-mark) {
    width: 24px;
    height: 24px;
    display: inline-block;
    background: var(--am-accent-primary);
    border-radius: 5px;
  }

  @include e(brand-text) {
    font-size: var(--am-font-lg);
    font-weight: 700;
  }

  @include e(menu) {
    min-width: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--am-spacing-xs);
    overflow-x: auto;
    scrollbar-width: none;

    &::-webkit-scrollbar {
      display: none;
    }
  }

  @include e(menu-icon) {
    width: 15px;
    height: 15px;
    flex: 0 0 auto;
  }

  @include e(menu-item) {
    height: 32px;
    padding: 8px 12px;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    color: var(--am-foreground-secondary);
    background: transparent;
    border: 0;
    border-radius: 6px;
    font-size: var(--am-font-sm);
    font-weight: 500;
    white-space: nowrap;
    cursor: pointer;
    transition:
      color 0.2s ease,
      background-color 0.2s ease;

    &:hover,
    &--active {
      color: var(--am-surface-primary);
      background: var(--am-accent-primary);
    }
  }

  @include e(right) {
    min-width: 180px;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: var(--am-spacing-sm);
  }
}

@media (max-width: 900px) {
  @include b(navbar) {
    padding: 0 var(--am-spacing-md);

    @include e(brand) {
      min-width: auto;
      flex: 0 0 auto;
    }

    @include e(brand-text) {
      display: none;
    }

    @include e(menu) {
      justify-content: flex-start;
    }

    @include e(right) {
      min-width: auto;
      flex: 0 0 auto;
    }
  }
}

@media (max-width: 560px) {
  @include b(navbar) {
    gap: var(--am-spacing-sm);

    @include e(menu-item) {
      padding-inline: var(--am-spacing-sm);
    }
  }
}
</style>
