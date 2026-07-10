<script setup lang="ts">
import ContainerMonitor from './container/index.vue'
import HostMonitor from './host/index.vue'
</script>

<template>
    <main class="workspace">
        <header class="workspace__header">
            <div>
                <p class="workspace__eyebrow">
                    OBSERVABILITY
                </p>
                <h1 class="workspace__title">
                    {{ $t('menu.monitor') }}
                </h1>
            </div>
            <span class="workspace__hint">Host & Container</span>
        </header>
        <section class="workspace__panel workspace__panel--monitor">
            <HostMonitor />
        </section>
        <section class="workspace__panel workspace__panel--monitor">
            <ContainerMonitor />
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
.workspace__title {
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
  min-height: 640px;
  overflow: hidden;
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-md);
  box-shadow: var(--am-shadow-subtle);
}

@media (max-width: 900px) {
  .workspace {
    padding: var(--am-spacing-md);
  }

  .workspace__header {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
