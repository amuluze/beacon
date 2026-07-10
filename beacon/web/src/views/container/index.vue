<script setup lang="ts">
import ContainerManager from './container/index.vue'
import ImageManager from './image/index.vue'
import NetworkManager from './network/index.vue'
</script>

<template>
    <main class="workspace">
        <header class="workspace__header">
            <div>
                <p class="workspace__eyebrow">
                    RUNTIME
                </p>
                <h1 class="workspace__title">
                    {{ $t('menu.container') }}
                </h1>
            </div>
            <span class="workspace__hint">Containers · Images · Networks</span>
        </header>
        <section class="workspace__panel">
            <h2 class="workspace__section-title">
                {{ $t('menu.containerManager') }}
            </h2>
            <div class="workspace__content">
                <ContainerManager />
            </div>
        </section>
        <section class="workspace__panel">
            <h2 class="workspace__section-title">
                {{ $t('menu.imageManager') }}
            </h2>
            <div class="workspace__content">
                <ImageManager />
            </div>
        </section>
        <section class="workspace__panel">
            <h2 class="workspace__section-title">
                {{ $t('menu.networkManager') }}
            </h2>
            <div class="workspace__content">
                <NetworkManager />
            </div>
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
  min-height: 520px;
  overflow: hidden;
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-md);
  box-shadow: var(--am-shadow-subtle);
}

.workspace__section-title {
  padding: var(--am-spacing-md) var(--am-spacing-lg) 0;
  font-size: var(--am-font-md);
}

.workspace__content {
  height: 470px;
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
