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
        <div v-else class="am-layout__loading" role="status" aria-live="polite">
            <span class="am-layout__loading-mark" aria-hidden="true" />
            <span>{{ $t('agent.loadingWorkspace') }}</span>
        </div>
    </section>
</template>

<style scoped lang="scss">
@include b(layout) {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  color: var(--am-foreground-primary);
  background-color: var(--am-surface-primary);
  overflow: hidden;

  @include e(loading) {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--am-spacing-sm);
    color: var(--am-foreground-muted);
    font-size: var(--am-font-sm);
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

  @include e(loading-mark) {
    width: 16px;
    height: 16px;
    border: 2px solid var(--am-border-primary);
    border-top-color: var(--am-accent-primary);
    border-radius: 50%;
    animation: workspace-loading 0.8s linear infinite;
  }
}

@keyframes workspace-loading {
  to {
    transform: rotate(360deg);
  }
}
</style>
