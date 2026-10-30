<script setup lang="ts">
import AlarmSettings from './alarm/index.vue'
import AuditLog from './components/AuditLog.vue'
import UserSettings from './components/UserSettings.vue'
import DockerSettings from './docker/index.vue'

import useStore from '@/store'

type SettingTab = 'config' | 'users' | 'audit'

const store = useStore()
const route = useRoute()
const router = useRouter()

const isAdmin = computed(() => store.user.userInfo.name === 'admin')

// 一级 Tab 与设计稿一致：设置 / 用户管理 / 审计日志，用户管理仅管理员可见。
const tabs = computed<{ name: SettingTab, label: string }[]>(() => {
  const items: { name: SettingTab, label: string }[] = [
    { name: 'config', label: 'menu.setting' },
  ]
  if (isAdmin.value)
    items.push({ name: 'users', label: 'menu.userManager' })
  items.push({ name: 'audit', label: 'audit.title' })
  return items
})

function resolveInitialTab(): SettingTab {
  const requested = route?.query.tab
  return tabs.value.some(tab => tab.name === requested) ? requested as SettingTab : 'config'
}

const activeTab = ref<SettingTab>(resolveInitialTab())

watch(activeTab, (tab) => {
  if (!router)
    return
  void router.replace({ query: { ...route.query, tab } })
})

// 账号权限变化后回落到当前账号仍可访问的 Tab。
watch(isAdmin, (admin) => {
  if (!admin && activeTab.value === 'users')
    activeTab.value = 'config'
})
</script>

<template>
    <main class="workspace">
        <header class="workspace__header">
            <div class="workspace__intro">
                <p class="workspace__eyebrow">
                    CONFIGURATION
                </p>
                <h1 class="workspace__title">
                    {{ $t('menu.setting') }}
                </h1>
                <p class="workspace__description">
                    {{ $t('setting.workspaceDescription') }}
                </p>
            </div>
            <span class="workspace__hint">
                {{ isAdmin ? $t('setting.workspaceHint') : $t('setting.workspaceHintRestricted') }}
            </span>
        </header>

        <div class="workspace__tabs" role="tablist">
            <button
                v-for="tab in tabs"
                :key="tab.name"
                type="button"
                role="tab"
                class="workspace__tab"
                :class="{ 'workspace__tab--active': activeTab === tab.name }"
                :aria-selected="activeTab === tab.name"
                @click="activeTab = tab.name"
            >
                {{ $t(tab.label) }}
            </button>
        </div>

        <template v-if="activeTab === 'config'">
            <section class="workspace__section">
                <header class="workspace__section-header">
                    <span class="workspace__section-bar" />
                    <h2 class="workspace__section-title">
                        {{ $t('menu.alarmSetting') }}
                    </h2>
                    <p class="workspace__section-note">
                        {{ $t('setting.alarmThresholdSettingTips') }}
                    </p>
                </header>
                <AlarmSettings />
            </section>
            <section class="workspace__section">
                <header class="workspace__section-header">
                    <span class="workspace__section-bar" />
                    <h2 class="workspace__section-title">
                        {{ $t('menu.systemDocker') }}
                    </h2>
                    <p class="workspace__section-note">
                        {{ $t('setting.sectionDockerNote') }}
                    </p>
                </header>
                <DockerSettings />
            </section>
        </template>
        <section v-else-if="activeTab === 'users'" class="workspace__section">
            <UserSettings />
        </section>
        <section v-else class="workspace__section">
            <header class="workspace__section-header">
                <span class="workspace__section-bar" />
                <h2 class="workspace__section-title">
                    {{ $t('setting.auditTitle') }}
                </h2>
                <p class="workspace__section-note">
                    {{ $t('setting.auditNote') }}
                </p>
            </header>
            <AuditLog />
        </section>
    </main>
</template>

<style scoped lang="scss">
.workspace {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-lg);
  width: 100%;
  padding: var(--am-spacing-lg);
}

.workspace__header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: var(--am-spacing-lg);
}

.workspace__intro {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-xs);
}

.workspace__eyebrow,
.workspace__title,
.workspace__description,
.workspace__section-title,
.workspace__section-note {
  margin: 0;
}

.workspace__eyebrow {
  color: var(--am-accent-primary);
  font-size: var(--am-font-xs);
  font-weight: 700;
  letter-spacing: 0.12em;
}

.workspace__title {
  font-size: var(--am-font-xl);
  font-weight: 700;
}

.workspace__description {
  max-width: 620px;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-sm);
  line-height: 1.55;
}

.workspace__hint {
  padding: 4px 10px;
  color: var(--am-foreground-secondary);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: var(--am-font-xs);
  background: var(--am-surface-primary);
  border: 1px solid var(--am-border-primary);
  border-radius: var(--am-radius-sm);
}

.workspace__tabs {
  display: inline-flex;
  align-self: flex-start;
  gap: var(--am-spacing-xs);
  padding: 4px;
  background: var(--am-surface-secondary);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-md);
}

.workspace__tab {
  padding: 7px 14px;
  color: var(--am-foreground-secondary);
  font-family: inherit;
  font-size: var(--am-font-sm);
  font-weight: 500;
  background: transparent;
  border: 1px solid transparent;
  border-radius: var(--am-radius-sm);
  cursor: pointer;
  transition:
    color 0.15s ease,
    background-color 0.15s ease,
    border-color 0.15s ease;
}

.workspace__tab:hover:not(.workspace__tab--active) {
  color: var(--am-foreground-primary);
}

.workspace__tab--active {
  color: var(--am-accent-primary);
  font-weight: 600;
  background: var(--am-surface-card);
  border-color: var(--am-border-primary);
}

.workspace__section {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 0;
}

.workspace__section-header {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.workspace__section-bar {
  flex: 0 0 auto;
  width: 3px;
  height: 14px;
  background: var(--am-accent-primary);
  border-radius: 2px;
}

.workspace__section-title {
  color: var(--am-foreground-primary);
  font-size: var(--am-font-md);
  font-weight: 600;
  white-space: nowrap;
}

.workspace__section-note {
  overflow: hidden;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
  text-overflow: ellipsis;
  white-space: nowrap;
}

:deep(.settings-grid) {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-md);
  min-width: 0;
}

:deep(.settings-card) {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-md);
  padding: var(--am-spacing-lg);
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-md);
  box-shadow: var(--am-shadow-subtle);
}

:deep(.settings-card--warning) {
  border-left: 3px solid var(--am-accent-warning);
}

:deep(.settings-card__header) {
  min-width: 0;
  display: flex;
  align-items: flex-start;
  gap: var(--am-spacing-md);
}

:deep(.settings-card__icon) {
  flex: 0 0 auto;
  width: 40px;
  height: 40px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--am-accent-primary);
  background: color-mix(in srgb, var(--am-accent-primary) 12%, transparent);
  border-radius: var(--am-radius-sm);
  font-size: 18px;
}

:deep(.settings-card__icon--warning) {
  color: var(--am-accent-warning);
  background: color-mix(in srgb, var(--am-accent-warning) 12%, transparent);
}

:deep(.settings-card__heading) {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-xs);
}

:deep(.settings-card__title),
:deep(.settings-card__description) {
  margin: 0;
}

:deep(.settings-card__title) {
  color: var(--am-foreground-primary);
  font-size: var(--am-font-md);
  font-weight: 600;
}

:deep(.settings-card__description) {
  color: var(--am-foreground-muted);
  font-size: var(--am-font-sm);
  line-height: 1.55;
}

:deep(.settings-card__badge) {
  flex: 0 0 auto;
  padding: 2px 8px;
  font-size: var(--am-font-xs);
  font-weight: 500;
  border-radius: var(--am-radius-sm);
}

:deep(.settings-card__badge--success) {
  color: var(--am-accent-success);
  background: color-mix(in srgb, var(--am-accent-success) 12%, transparent);
}

:deep(.settings-card__badge--muted) {
  color: var(--am-foreground-muted);
  background: var(--am-surface-secondary);
}

:deep(.settings-card__hint) {
  padding-top: 10px;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
  line-height: 1.5;
  border-top: 1px solid var(--am-border-subtle);
}

:deep(.settings-card__textarea) {
  min-width: 0;
}

:deep(.settings-card__textarea .el-textarea__inner) {
  min-height: 110px !important;
  padding: var(--am-spacing-md);
  color: var(--am-foreground-primary);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: var(--am-font-sm);
  line-height: 1.7;
  background: var(--am-surface-primary);
  border: 1px solid var(--am-border-primary);
  border-radius: var(--am-radius-sm);
  box-shadow: none;
}

:deep(.settings-card__actions) {
  display: flex;
  justify-content: flex-end;
  gap: var(--am-spacing-sm);
}

:deep(.settings-panels) {
  display: flex;
  align-items: stretch;
  gap: var(--am-spacing-lg);
  min-width: 0;
}

:deep(.settings-panel) {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-sm);
  padding: var(--am-spacing-md);
  background: var(--am-surface-secondary);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-sm);
}

:deep(.settings-panel__header) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--am-spacing-sm);
}

:deep(.settings-panel__label) {
  color: var(--am-foreground-secondary);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: var(--am-font-xs);
  font-weight: 600;
  letter-spacing: 0.04em;
}

:deep(.settings-panel__title) {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--am-foreground-primary);
  font-size: var(--am-font-sm);
  font-weight: 600;
}

:deep(.settings-panel__mark) {
  width: 24px;
  height: 24px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--am-accent-primary);
  background: color-mix(in srgb, var(--am-accent-primary) 12%, transparent);
  border-radius: var(--am-radius-sm);
  font-size: 15px;
}

:deep(.settings-panel__footer) {
  display: flex;
  justify-content: flex-end;
  gap: var(--am-spacing-sm);
}

:deep(.settings-panel__cols) {
  display: flex;
  align-items: flex-start;
  gap: var(--am-spacing-lg);
  min-width: 0;
}

:deep(.settings-panel__col) {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-sm);
}

:deep(.settings-row) {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
  min-width: 0;
}

:deep(.settings-row__label) {
  flex: 0 0 90px;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
}

:deep(.settings-row__value) {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  color: var(--am-foreground-primary);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: var(--am-font-sm);
  text-overflow: ellipsis;
  white-space: nowrap;
}

:deep(.settings-row__value--text) {
  font-family: inherit;
}

:deep(.settings-row__value--ok) {
  color: var(--am-accent-success);
  font-family: inherit;
}

:deep(.settings-row__hint) {
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
}

:deep(.settings-thresholds) {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

:deep(.settings-threshold) {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: var(--am-spacing-md);
  padding: 8px 0;
}

:deep(.settings-threshold + .settings-threshold) {
  border-top: 1px solid var(--am-border-subtle);
}

:deep(.settings-threshold__label) {
  flex: 0 0 120px;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-sm);
}

:deep(.settings-threshold__value) {
  flex: 1;
  min-width: 0;
  color: var(--am-foreground-primary);
  font-size: var(--am-font-sm);
}

:deep(.settings-chip.el-button) {
  height: 26px;
  padding: 0 12px;
  color: var(--am-accent-primary);
  font-size: var(--am-font-sm);
  font-weight: 500;
  background: var(--am-surface-primary);
  border: 1px solid var(--am-border-primary);
  border-radius: var(--am-radius-sm);
}

:deep(.settings-chip.el-button:hover:not(.is-disabled)) {
  color: var(--am-accent-primary);
  background: color-mix(in srgb, var(--am-accent-primary) 10%, transparent);
  border-color: var(--am-accent-primary);
}

:deep(.settings-chip--sm.el-button) {
  height: 24px;
  padding: 0 10px;
  font-size: var(--am-font-xs);
}

:deep(.settings-chip--sm.el-button > span) {
  display: inline-flex;
  align-items: center;
  gap: 5px;
}

@media (max-width: 900px) {
  .workspace {
    padding: var(--am-spacing-md);
  }

  .workspace__header {
    align-items: flex-start;
    flex-direction: column;
  }

  .workspace__section-header {
    align-items: flex-start;
    flex-direction: column;
    gap: var(--am-spacing-xs);
  }

  .workspace__section-note {
    white-space: normal;
  }

  :deep(.settings-panels),
  :deep(.settings-panel__cols) {
    flex-direction: column;
  }
}
</style>
