<script setup lang="ts">
import type { ContainerSummary } from '@/interface/container.ts'

import { queryContainerSummary } from '@/api/container'
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { convertBytesToReadable } from '@/utils/convert'

import ContainerManager from './container/index.vue'
import ImageManager from './image/index.vue'
import NetworkManager from './network/index.vue'

type ContainerTab = 'containers' | 'images' | 'networks'

const tabs: { name: ContainerTab, label: string }[] = [
  { name: 'containers', label: 'container.tabContainers' },
  { name: 'images', label: 'container.tabImages' },
  { name: 'networks', label: 'container.tabNetworks' },
]

// §8.5 无 Agent 时三段各自渲染空态卡片，不做整页拦截
const emptySections = [
  { title: 'container.tabContainers', description: 'container.emptyContainersDesc' },
  { title: 'container.tabImages', description: 'container.emptyImagesDesc' },
  { title: 'container.tabNetworks', description: 'container.emptyNetworksDesc' },
]

const route = useRoute()
const router = useRouter()

function resolveInitialTab(): ContainerTab {
  const requested = route?.query.tab
  return tabs.some(tab => tab.name === requested) ? requested as ContainerTab : 'containers'
}

const activeTab = ref<ContainerTab>(resolveInitialTab())
// 三个 Tab 的总数在各自首次加载后回填，未加载的 Tab 不展示计数。
const counts = reactive<Record<ContainerTab, number | null>>({
  containers: null,
  images: null,
  networks: null,
})

watch(activeTab, (tab) => {
  if (!router)
    return
  void router.replace({ query: { ...route.query, tab } })
})

const { isAgentEmpty, loadAgents, selectedAgentID } = useAgentSelection({ immediate: false })

const summary = ref<ContainerSummary | null>(null)

// 标题区的「刷新」按钮驱动三个 Tab 的列表重新拉取，并同步刷新概览统计。
const refreshToken = ref(0)

async function loadSummary() {
  if (!selectedAgentID.value) {
    summary.value = null
    return
  }
  try {
    const { data } = await queryContainerSummary()
    summary.value = data ?? null
  }
  catch {
    // 概览失败不影响三个 Tab 的列表与操作
    summary.value = null
  }
}

onMounted(loadSummary)
watch(selectedAgentID, loadSummary)

async function refreshAgents() {
  // Agent 列表刷新后由各 Tab 自身的 selectedAgentID 监听触发重新查询。
  await loadAgents()
}

function refreshWorkspace() {
  refreshToken.value += 1
  void loadSummary()
}
</script>

<template>
    <main class="container-workspace">
        <header class="container-workspace__header">
            <div class="container-workspace__intro">
                <p class="container-workspace__eyebrow">
                    RUNTIME
                </p>
                <h1 class="container-workspace__title">
                    {{ $t('menu.container') }}
                </h1>
                <p class="container-workspace__description">
                    {{ $t('container.workspaceDescription') }}
                </p>
            </div>
            <div class="container-workspace__aside">
                <div v-if="summary" class="container-workspace__summary">
                    <span class="container-workspace__metric">
                        {{ $t('container.tabContainers') }}
                        <strong>{{ summary.containers }}</strong>
                        · {{ $t('container.enable') }}
                        <strong>{{ summary.containers_running }}</strong>
                    </span>
                    <span class="container-workspace__metric">
                        {{ $t('container.tabImages') }}
                        <strong>{{ summary.images }}</strong>
                        · {{ convertBytesToReadable(summary.image_bytes) }}
                    </span>
                    <span class="container-workspace__metric">
                        {{ $t('container.tabNetworks') }}
                        <strong>{{ summary.networks }}</strong>
                    </span>
                </div>
                <button class="container-workspace__refresh" type="button" @click="refreshWorkspace">
                    <svg-icon icon-class="update" />
                    {{ $t('agent.refresh') }}
                </button>
            </div>
        </header>
        <!-- §8.5 无可用 Agent 时容器 / 镜像 / 网络三段各自渲染空态卡片与刷新入口 -->
        <div v-if="isAgentEmpty" class="container-workspace__empty">
            <section v-for="section in emptySections" :key="section.title" class="container-workspace__empty-section">
                <h2 class="container-workspace__empty-title">
                    {{ $t(section.title) }}
                </h2>
                <AgentEmptyState
                    :description="$t(section.description)"
                    min-height="240px"
                    :show-install="false"
                    @refresh="refreshAgents"
                />
            </section>
        </div>
        <el-tabs v-else v-model="activeTab" class="container-workspace__tabs">
            <el-tab-pane v-for="tab in tabs" :key="tab.name" :name="tab.name" lazy>
                <template #label>
                    <span class="container-workspace__tab-label">
                        {{ $t(tab.label) }}
                        <span v-if="counts[tab.name] !== null" class="container-workspace__tab-count">{{ counts[tab.name] }}</span>
                    </span>
                </template>
                <div class="container-workspace__panel">
                    <ContainerManager
                        v-if="tab.name === 'containers'"
                        :refresh-token="refreshToken"
                        @total="counts.containers = $event"
                    />
                    <ImageManager
                        v-else-if="tab.name === 'images'"
                        :refresh-token="refreshToken"
                        @total="counts.images = $event"
                    />
                    <NetworkManager
                        v-else
                        :refresh-token="refreshToken"
                        @total="counts.networks = $event"
                    />
                </div>
            </el-tab-pane>
        </el-tabs>
    </main>
</template>

<style scoped lang="scss">
.container-workspace {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-md);
  width: 100%;
  height: 100%;
  min-height: 0;
  padding: var(--am-spacing-lg);
}

.container-workspace__header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: var(--am-spacing-lg);
}

.container-workspace__intro {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-xs);
}

.container-workspace__aside {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--am-spacing-sm);
  flex-wrap: wrap;
}

.container-workspace__eyebrow,
.container-workspace__title,
.container-workspace__description {
  margin: 0;
}

.container-workspace__eyebrow {
  color: var(--am-accent-primary);
  font-size: var(--am-font-xs);
  font-weight: 700;
  letter-spacing: 2px;
}

.container-workspace__title {
  font-size: var(--am-font-2xl);
  font-weight: 700;
}

.container-workspace__description {
  max-width: 620px;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-sm);
  line-height: 1.55;
}

.container-workspace__summary {
  display: flex;
  flex-wrap: wrap;
  gap: var(--am-spacing-sm);
}

.container-workspace__metric {
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

.container-workspace__refresh {
  height: 32px;
  padding: 6px 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--am-spacing-sm);
  color: var(--am-accent-primary);
  font-family: inherit;
  font-size: var(--am-font-sm);
  font-weight: 500;
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-primary);
  border-radius: var(--am-radius-sm);
  cursor: pointer;
  transition:
    color 0.2s ease,
    background-color 0.2s ease,
    border-color 0.2s ease;

  &:hover,
  &:focus-visible {
    color: var(--am-accent-primary-hover);
    border-color: var(--am-accent-primary-hover);
  }
}

.container-workspace__metric strong {
  color: var(--am-foreground-primary);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-weight: 600;
}

.container-workspace__empty {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-sm);
}

.container-workspace__empty-section {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-sm);
}

.container-workspace__empty-title {
  margin: 0;
  color: var(--am-foreground-primary);
  font-size: var(--am-font-md);
  font-weight: 600;
}

.container-workspace__tabs {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.container-workspace__tabs :deep(.el-tabs__header) {
  margin-bottom: var(--am-spacing-md);
}

.container-workspace__tabs :deep(.el-tabs__content) {
  flex: 1;
  min-height: 0;
}

.container-workspace__tabs :deep(.el-tab-pane) {
  height: 100%;
}

.container-workspace__tab-label {
  display: inline-flex;
  align-items: center;
  gap: var(--am-spacing-xs);
  font-size: var(--am-font-sm);
  font-weight: 600;
}

.container-workspace__tab-count {
  padding: 0 6px;
  color: var(--am-foreground-secondary);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: var(--am-font-xs);
  background: var(--am-surface-secondary);
  border: 1px solid var(--am-border-subtle);
  border-radius: var(--am-radius-sm);
}

.container-workspace__panel {
  height: 100%;
  min-height: 0;
}

.container-workspace__panel :deep(.am-container) {
  padding: 0;
  background: transparent;
}

@media (max-width: 900px) {
  .container-workspace {
    padding: var(--am-spacing-md);
  }
}
</style>
