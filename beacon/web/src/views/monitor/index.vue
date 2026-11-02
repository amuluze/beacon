<script setup lang="ts">
import ContainerMonitor from './container/index.vue'
import HostMonitor from './host/index.vue'
</script>

<template>
    <main class="workspace">
        <header class="workspace__header">
            <div class="workspace__intro">
                <p class="workspace__eyebrow">
                    OBSERVABILITY
                </p>
                <h1 class="workspace__title">
                    {{ $t('menu.monitor') }}
                </h1>
                <p class="workspace__description">
                    {{ $t('monitor.workspaceDescription') }}
                </p>
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

.workspace__intro {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-xs);
}

.workspace__eyebrow {
  color: var(--am-accent-primary);
  font-size: var(--am-font-xs);
  font-weight: 700;
  letter-spacing: 2px;
}

.workspace__title {
  font-size: var(--am-font-2xl);
  font-weight: 700;
}

.workspace__description {
  max-width: 620px;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-sm);
  line-height: 1.55;
}

.workspace__hint {
  height: 28px;
  padding: 4px 10px;
  display: inline-flex;
  align-items: center;
  color: var(--am-foreground-secondary);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: var(--am-font-xs);
  background: var(--am-surface-inset);
  border: 1px solid var(--am-border-primary);
  border-radius: var(--am-radius-sm);
}

.workspace__panel {
  // 高度由内容自适应：主机监控包含 2 个图表 row，容器监控包含 1 个 row，
  // 二者图表高度在全局 .am-chart-area 已统一，无需再以 min-height 撑开容器面板，
  // 否则容器监控会在单 row 时底部留下大面积空白。
  min-height: 0;
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
