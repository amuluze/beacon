<script setup lang="ts">
import Avatar from '@/layout/navbar/Avatar.vue'
import Language from '@/layout/navbar/Language.vue'
import ThemeChange from '@/layout/navbar/ThemeChange.vue'
import { dynamicRoutes } from '@/router/dynamic.ts'
import useStore from '@/store'
import type { RouteRecordRaw } from 'vue-router'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const router = useRouter()
const store = useStore()
const { t } = useI18n()

const primaryRouteNames = new Set(['monitor', 'container', 'setting'])

const visibleRoutes = computed(() => {
  return dynamicRoutes.filter((item) => {
    if (!item.meta?.show)
      return false
    if (store.user.userInfo.name !== 'admin' && item.name === 'account')
      return false
    return true
  })
})

const primaryMenus = computed(() => visibleRoutes.value.filter(item => primaryRouteNames.has(String(item.name))))
const secondaryMenus = computed(() => visibleRoutes.value.filter(item => !primaryRouteNames.has(String(item.name))))

function routeTarget(item: RouteRecordRaw): string {
  return typeof item.redirect === 'string' ? item.redirect : item.path
}

function isRouteActive(item: RouteRecordRaw): boolean {
  return route.path === item.path || route.path.startsWith(`${item.path}/`)
}

function menuIcon(item: RouteRecordRaw): string {
  return String(item.meta?.icon || 'menu')
}

function menuTitle(item: RouteRecordRaw): string {
  return t(String(item.meta?.title || item.name || ''))
}

function goRoute(item: RouteRecordRaw): void {
  router.push(routeTarget(item))
}

function goSecondary(path: string): void {
  router.push(path)
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
                v-for="item in primaryMenus"
                :key="String(item.name)"
                class="am-navbar__menu-item"
                :class="{ 'am-navbar__menu-item--active': isRouteActive(item) }"
                type="button"
                @click="goRoute(item)"
            >
                <svg-icon :icon-class="menuIcon(item)" size="15px" />
                <span>{{ menuTitle(item) }}</span>
            </button>
            <el-dropdown v-if="secondaryMenus.length > 0" trigger="click" @command="goSecondary">
                <button class="am-navbar__menu-item" type="button">
                    <svg-icon icon-class="more" size="15px" />
                    <span>{{ t('container.more') }}</span>
                </button>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item
                            v-for="item in secondaryMenus"
                            :key="String(item.name)"
                            :command="routeTarget(item)"
                        >
                            <svg-icon :icon-class="menuIcon(item)" size="14px" />
                            <span class="am-navbar__dropdown-label">{{ menuTitle(item) }}</span>
                        </el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </nav>

        <div class="am-navbar__right">
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

  @include e(dropdown-label) {
    margin-left: var(--am-spacing-xs);
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
