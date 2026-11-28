<script setup lang="ts">
import Content from '@/layout/content/index.vue'
import Navbar from '@/layout/navbar/index.vue'
import StatusBar from '@/layout/statusbar/index.vue'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import useStore from '@/store'

const store = useStore()
const readOnly = computed(() => store.user.userInfo.readOnly === true)
const { ensureSelectedAgent } = useAgentSelection({ immediate: false })
const workspaceReady = shallowRef(false)

onBeforeMount(async () => {
  try {
    await ensureSelectedAgent()
  }
  catch {
    // 请求层已经展示具体错误；允许工作区继续挂载，以便空状态提供重试入口。
  }
  finally {
    workspaceReady.value = true
  }
})
</script>

<template>
    <section class="am-layout">
        <Navbar />
        <p v-if="readOnly" class="am-layout__readonly" role="status">
            {{ $t('common.readOnlyBanner') }}
        </p>
        <template v-if="workspaceReady">
            <Content />
            <StatusBar />
        </template>
        <div v-else class="am-layout__loading" role="status" aria-live="polite" :aria-label="$t('agent.loadingWorkspace')">
            <span class="am-layout__loading-line am-layout__loading-line--eyebrow" aria-hidden="true" />
            <span class="am-layout__loading-line am-layout__loading-line--title" aria-hidden="true" />
            <span class="am-layout__loading-line am-layout__loading-line--copy" aria-hidden="true" />
            <div class="am-layout__loading-grid" aria-hidden="true">
                <span v-for="index in 4" :key="index" class="am-layout__loading-card" />
            </div>
        </div>
    </section>
</template>

<style scoped lang="scss">
@include b(layout) {
  display: flex;
  flex-direction: column;
  width: 100%;
  min-width: 1280px;
  height: 100%;
  color: var(--am-foreground-primary);
  background-color: var(--am-surface-primary);
  overflow: hidden;

  @include e(loading) {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: var(--am-spacing-md);
    padding: var(--am-spacing-xl) var(--am-spacing-lg);
    overflow: hidden;
  }

  @include e(loading-line) {
    display: block;
    border-radius: var(--am-radius-sm);
    background: var(--am-surface-secondary);
    animation: am-skeleton-pulse 1.25s ease-in-out infinite;

    &--eyebrow {
      width: 76px;
      height: 11px;
    }

    &--title {
      width: 190px;
      height: 32px;
    }

    &--copy {
      width: min(560px, 70%);
      height: 14px;
    }
  }

  @include e(loading-grid) {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: var(--am-spacing-lg);
    margin-top: var(--am-spacing-sm);
  }

  @include e(loading-card) {
    height: 260px;
    border: 1px solid var(--am-border-subtle);
    border-radius: var(--am-radius-md);
    background:
      linear-gradient(var(--am-surface-secondary), var(--am-surface-secondary)) 24px 24px / 35% 14px no-repeat,
      linear-gradient(var(--am-surface-secondary), var(--am-surface-secondary)) 24px 56px / 55% 28px no-repeat,
      linear-gradient(var(--am-surface-secondary), var(--am-surface-secondary)) 24px 110px / calc(100% - 48px) 140px
        no-repeat,
      var(--am-surface-card);
    animation: am-skeleton-pulse 1.25s ease-in-out infinite;
  }

  @include e(readonly) {
    flex: 0 0 auto;
    margin: 0;
    padding: 6px var(--am-spacing-lg);
    color: var(--am-accent-warning);
    background: var(--am-accent-warning-soft);
    border-bottom: 1px solid color-mix(in srgb, var(--am-accent-warning) 28%, transparent);
    font-size: var(--am-font-xs);
    text-align: center;
  }

}
</style>
