<script setup lang="ts">
import AlarmSettings from './alarm/index.vue'
import AuditLog from './components/AuditLog.vue'
import UserSettings from './components/UserSettings.vue'
import DockerSettings from './docker/index.vue'
import HostSettings from './host/index.vue'

import useStore from '@/store'

const store = useStore()
const isAdmin = computed(() => store.user.userInfo.name === 'admin')
</script>

<template>
    <main class="workspace">
        <header class="workspace__header">
            <div>
                <p class="workspace__eyebrow">
                    CONFIGURATION
                </p>
                <h1 class="workspace__title">
                    {{ $t('menu.setting') }}
                </h1>
            </div>
            <span class="workspace__hint">
                System · Alerts · Docker<span v-if="isAdmin"> · Users</span> · Audit
            </span>
        </header>
        <section class="workspace__panel workspace__panel--flush">
            <h2 class="workspace__section-title">
                {{ $t('menu.systemSetting') }}
            </h2>
            <HostSettings />
        </section>
        <section class="workspace__panel workspace__panel--flush">
            <h2 class="workspace__section-title">
                {{ $t('menu.alarmSetting') }}
            </h2>
            <AlarmSettings />
        </section>
        <section class="workspace__panel workspace__panel--flush">
            <h2 class="workspace__section-title">
                {{ $t('menu.systemDocker') }}
            </h2>
            <DockerSettings />
        </section>
        <section v-if="isAdmin" class="workspace__panel">
            <h2 class="workspace__section-title">
                {{ $t('menu.accountSetting') }}
            </h2>
            <UserSettings />
        </section>
        <section class="workspace__panel">
            <h2 class="workspace__section-title">
                {{ $t('audit.title') }}
            </h2>
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
  gap: var(--am-spacing-md);
}

.workspace__eyebrow,
.workspace__title,
.workspace__section-title {
  margin: 0;
}

.workspace__eyebrow {
  color: var(--am-accent-primary);
  font-size: var(--am-font-xs);
  font-weight: 700;
  letter-spacing: 0.12em;
}

.workspace__title {
  margin-top: var(--am-spacing-xs);
  font-size: var(--am-font-xl);
}

.workspace__hint {
  color: var(--am-foreground-muted);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: var(--am-font-xs);
}

.workspace__panel {
  padding: var(--am-spacing-lg);
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-md);
  box-shadow: var(--am-shadow-subtle);
}

.workspace__panel--flush {
  padding: 0;
  background: transparent;
  border: 0;
  border-radius: 0;
  box-shadow: none;
}

.workspace__section-title {
  margin-bottom: var(--am-spacing-md);
  font-size: var(--am-font-md);
}

:deep(.settings-grid) {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--am-spacing-md);
}

:deep(.settings-card) {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-lg);
  padding: var(--am-spacing-lg);
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-lg);
  box-shadow: var(--am-shadow-subtle);
}

:deep(.settings-card--wide) {
  grid-column: 1 / -1;
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
  border-radius: var(--am-radius-md);
  font-size: 19px;
}

:deep(.settings-card__icon--success) {
  color: var(--am-accent-success);
  background: color-mix(in srgb, var(--am-accent-success) 12%, transparent);
}

:deep(.settings-card__icon--warning) {
  color: var(--am-accent-warning);
  background: color-mix(in srgb, var(--am-accent-warning) 12%, transparent);
}

:deep(.settings-card__heading),
:deep(.settings-card__meta),
:deep(.settings-threshold__main) {
  min-width: 0;
}

:deep(.settings-card__title),
:deep(.settings-card__description),
:deep(.settings-card__tip) {
  margin: 0;
}

:deep(.settings-card__title) {
  color: var(--am-foreground-primary);
  font-size: var(--am-font-md);
  font-weight: 650;
}

:deep(.settings-card__description) {
  margin-top: var(--am-spacing-xs);
  color: var(--am-foreground-muted);
  font-size: var(--am-font-sm);
  line-height: 1.55;
}

:deep(.settings-card__value-row) {
  min-width: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--am-spacing-md);
  margin-top: auto;
  padding: var(--am-spacing-md);
  background: var(--am-surface-primary);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-md);
}

:deep(.settings-card__meta) {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-xs);
}

:deep(.settings-card__label) {
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
  font-weight: 600;
  letter-spacing: 0.02em;
}

:deep(.settings-card__value) {
  min-width: 0;
  color: var(--am-foreground-primary);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: var(--am-font-sm);
}

:deep(.settings-card__value--plain) {
  font-family: inherit;
  font-weight: 600;
}

:deep(.settings-card__actions) {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: var(--am-spacing-sm);
}

:deep(.settings-card__actions--end) {
  justify-content: flex-end;
  margin-top: auto;
}

:deep(.settings-card__form) {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-sm);
}

:deep(.settings-card__controls) {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
}

:deep(.settings-card__textarea .el-textarea__inner) {
  min-height: 140px !important;
  padding: var(--am-spacing-md);
  color: var(--am-foreground-primary);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  line-height: 1.65;
  background: var(--am-surface-primary);
}

:deep(.settings-card__footer) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--am-spacing-md);
}

:deep(.settings-card__tip) {
  max-width: 70ch;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
  line-height: 1.5;
}

:deep(.settings-thresholds) {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-sm);
}

:deep(.settings-threshold) {
  min-width: 0;
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  align-items: center;
  gap: var(--am-spacing-sm);
  padding: var(--am-spacing-sm) var(--am-spacing-md);
  background: var(--am-surface-primary);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-md);
}

:deep(.settings-threshold__icon) {
  width: 30px;
  height: 30px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--am-accent-primary);
  background: color-mix(in srgb, var(--am-accent-primary) 10%, transparent);
  border-radius: var(--am-radius-sm);
}

:deep(.settings-threshold__icon--success) {
  color: var(--am-accent-success);
  background: color-mix(in srgb, var(--am-accent-success) 10%, transparent);
}

:deep(.settings-threshold__icon--danger) {
  color: var(--am-accent-danger);
  background: color-mix(in srgb, var(--am-accent-danger) 10%, transparent);
}

:deep(.settings-threshold__main) {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

:deep(.settings-threshold__title) {
  color: var(--am-foreground-primary);
  font-size: var(--am-font-sm);
}

:deep(.settings-threshold__value) {
  overflow: hidden;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
  text-overflow: ellipsis;
  white-space: nowrap;
}

@media (max-width: 900px) {
  .workspace {
    padding: var(--am-spacing-md);
  }

  .workspace__header {
    align-items: flex-start;
    flex-direction: column;
  }

  .workspace__panel:not(.workspace__panel--flush) {
    padding: var(--am-spacing-md);
  }
}

@media (max-width: 760px) {
  :deep(.settings-grid) {
    grid-template-columns: minmax(0, 1fr);
  }

  :deep(.settings-card) {
    padding: var(--am-spacing-md);
  }

  :deep(.settings-card__controls),
  :deep(.settings-card__footer) {
    align-items: stretch;
    flex-direction: column;
  }

  :deep(.settings-card__footer .el-button) {
    align-self: flex-start;
  }
}
</style>
