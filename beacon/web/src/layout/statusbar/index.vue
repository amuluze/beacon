<script setup lang="ts">
import type { DockerInfo } from '@/interface/container.ts'
import type { HostInfo } from '@/interface/host.ts'
import { queryDockerInfo } from '@/api/container'
import { queryHostInfo } from '@/api/host'
import { querySystemVersion } from '@/api/system'
import { useAgentSelection } from '@/hooks/useAgentSelection'

const { agentList, selectedAgentID, ensureSelectedAgent } = useAgentSelection({ immediate: false })

const hostInfo = shallowRef<HostInfo>()
const dockerInfo = shallowRef<DockerInfo>()
// §6.3 右组末项固定展示 Beacon 服务端构建版本，不随 Agent 切换变化
const serverVersion = shallowRef('-')
const loading = shallowRef(false)

const selectedAgent = computed(() => agentList.value.find(item => item.agent_id === selectedAgentID.value))

const hostItems = computed(() => [
  { label: '主机', value: hostInfo.value?.hostname || selectedAgent.value?.hostname || '-' },
  { label: '系统', value: hostInfo.value?.platform_version || selectedAgent.value?.os || '-' },
  { label: '内核', value: hostInfo.value?.kernel_version || '-' },
  { label: 'OS', value: hostInfo.value ? `${hostInfo.value.os}/${hostInfo.value.kernel_arch}` : selectedAgent.value ? `${selectedAgent.value.os}/${selectedAgent.value.arch}` : '-' },
])

const dockerItems = computed(() => [
  { label: 'Docker', value: dockerInfo.value?.docker_version || '-' },
  { label: 'API', value: dockerInfo.value ? `${dockerInfo.value.min_api_version}-${dockerInfo.value.api_version}` : '-' },
  { label: 'OS', value: dockerInfo.value ? `${dockerInfo.value.os}/${dockerInfo.value.arch}` : '-' },
  { label: 'Beacon', value: serverVersion.value || '-' },
])

async function refreshStatus() {
  const agentID = await ensureSelectedAgent()
  if (!agentID)
    return

  loading.value = true
  try {
    const [hostResult, dockerResult, versionResult] = await Promise.allSettled([
      queryHostInfo(),
      queryDockerInfo(),
      querySystemVersion(),
    ])
    if (hostResult.status === 'fulfilled')
      hostInfo.value = hostResult.value.data
    if (dockerResult.status === 'fulfilled')
      dockerInfo.value = dockerResult.value.data
    if (versionResult.status === 'fulfilled')
      serverVersion.value = versionResult.value.data?.version || '-'
  }
  finally {
    loading.value = false
  }
}

watch(selectedAgentID, () => {
  void refreshStatus()
}, { immediate: true })
</script>

<template>
    <footer v-if="selectedAgentID || agentList.length > 0" class="am-statusbar" :class="{ 'am-statusbar--loading': loading }">
        <div class="am-statusbar__group">
            <template v-for="(item, index) in hostItems" :key="item.label">
                <span v-if="index > 0" class="am-statusbar__sep">·</span>
                <span class="am-statusbar__item">
                    <span class="am-statusbar__label">{{ item.label }}</span>
                    <span class="am-statusbar__value">{{ item.value }}</span>
                </span>
            </template>
        </div>

        <div class="am-statusbar__uptime">
            <span class="am-statusbar__label">运行时间</span>
            <span class="am-statusbar__value">{{ hostInfo?.uptime || '-' }}</span>
        </div>

        <div class="am-statusbar__group">
            <template v-for="(item, index) in dockerItems" :key="item.label">
                <span v-if="index > 0" class="am-statusbar__sep">·</span>
                <span class="am-statusbar__item">
                    <span class="am-statusbar__label">{{ item.label }}</span>
                    <span class="am-statusbar__value">{{ item.value }}</span>
                </span>
            </template>
        </div>
    </footer>
</template>

<style scoped lang="scss">
@include b(statusbar) {
  flex: 0 0 var(--am-statusbar-height);
  height: var(--am-statusbar-height);
  min-height: var(--am-statusbar-height);
  padding: 0 var(--am-spacing-lg);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--am-spacing-md);
  overflow-x: auto;
  color: var(--am-foreground-secondary);
  background: var(--am-surface-secondary);
  border-top: 1px solid var(--am-border-primary);
  scrollbar-width: none;

  &::-webkit-scrollbar {
    display: none;
  }

  @include m(loading) {
    opacity: 0.88;
  }

  @include e(group) {
    display: flex;
    align-items: center;
    gap: var(--am-spacing-md);
    min-width: 0;
    flex: 1 1 auto;
    white-space: nowrap;
  }

  @include e(item) {
    display: inline-flex;
    align-items: center;
    gap: var(--am-spacing-xs);
    min-width: 0;
    flex: 0 1 auto;
  }

  @include e(label) {
    color: var(--am-foreground-secondary);
    font-size: var(--am-font-xs);
    white-space: nowrap;
  }

  @include e(value) {
    max-width: 180px;
    overflow: hidden;
    color: var(--am-foreground-secondary);
    font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
    font-size: var(--am-font-xs);
    font-weight: 500;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  @include e(sep) {
    color: var(--am-foreground-secondary);
    font-size: var(--am-font-xs);
  }

  @include e(uptime) {
    flex: 0 0 auto;
    display: inline-flex;
    align-items: center;
    gap: var(--am-spacing-xs);
    white-space: nowrap;
  }
}

@media (max-width: 980px) {
  @include b(statusbar) {
    justify-content: flex-start;
    padding: 0 var(--am-spacing-md);

    @include e(group) {
      flex: 0 0 auto;
    }
  }
}
</style>
