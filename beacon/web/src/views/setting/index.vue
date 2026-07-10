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
        <section class="workspace__panel">
            <h2 class="workspace__section-title">
                {{ $t('menu.systemSetting') }}
            </h2>
            <HostSettings />
        </section>
        <section class="workspace__panel">
            <h2 class="workspace__section-title">
                {{ $t('menu.alarmSetting') }}
            </h2>
            <AlarmSettings />
        </section>
        <section class="workspace__panel">
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

.workspace__section-title {
  margin-bottom: var(--am-spacing-md);
  font-size: var(--am-font-md);
}

@media (max-width: 900px) {
  .workspace {
    padding: var(--am-spacing-md);
  }

  .workspace__header {
    align-items: flex-start;
    flex-direction: column;
  }

  .workspace__panel {
    padding: var(--am-spacing-md);
  }
}
</style>
