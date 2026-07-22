<script setup lang="ts">
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import Terminal, { type TerminalStatus } from '@/components/Terminal/index.vue'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { computed, onMounted, shallowRef, useTemplateRef } from 'vue'

const route = useRoute()
const { t } = useI18n()
const { agentList, selectedAgentID, loading, isAgentEmpty, loadAgents, ensureSelectedAgent } = useAgentSelection({ immediate: false })

const routeAgentID = computed(() => route.query.agent_id as string | undefined)
const agentId = computed(() => selectedAgentID.value)

const terminalRef = useTemplateRef<InstanceType<typeof Terminal>>('terminal')
const status = shallowRef<TerminalStatus>('idle')
const cols = shallowRef(0)
const rows = shallowRef(0)

const selectedAgent = computed({
  get: () => selectedAgentID.value,
  set: (value: string) => selectedAgentID.value = value,
})

const selectedHostname = computed(() => {
  const match = agentList.value.find(a => a.agent_id === selectedAgentID.value)
  return match?.hostname || selectedAgentID.value || '—'
})

const termTitle = computed(() => {
  const dims = cols.value && rows.value ? `${cols.value}x${rows.value}` : '—'
  return `${selectedHostname.value} — bash — ${dims}`
})

const statusMeta = computed(() => {
  switch (status.value) {
  case 'connected':
    return { dot: 'var(--am-accent-success)', text: t('agent.terminalStatusConnected'), tone: 'success' }
  case 'connecting':
    return { dot: 'var(--am-accent-warning, #f59e0b)', text: t('agent.terminalStatusConnecting'), tone: 'warning' }
  case 'error':
    return { dot: 'var(--am-accent-danger, #ef4444)', text: t('agent.terminalStatusError'), tone: 'danger' }
  default:
    return { dot: 'var(--am-foreground-muted)', text: t('agent.terminalStatusDisconnected'), tone: 'muted' }
  }
})

const actionable = computed(() => Boolean(agentId.value) && !isAgentEmpty.value)

function handleClear() {
  terminalRef.value?.clear()
}

function handleNewSession() {
  terminalRef.value?.newSession()
}

function handleStatusChange(nextStatus: TerminalStatus) {
  status.value = nextStatus
}

function handleTerminalResize(size: { rows: number, cols: number }) {
  rows.value = size.rows
  cols.value = size.cols
}

onMounted(async () => {
  await loadAgents()
  if (routeAgentID.value) {
    selectedAgent.value = routeAgentID.value
  }
  else {
    await ensureSelectedAgent()
  }
})
</script>

<template>
    <div class="am-terminal-page">
        <header class="am-terminal-page__header">
            <p class="am-terminal-page__eyebrow">
                {{ $t('agent.terminalEyebrow') }}
            </p>
            <h1 class="am-terminal-page__title">
                {{ $t('agent.terminalTitle') }}
            </h1>
            <p class="am-terminal-page__hint">
                {{ $t('agent.terminalHint') }}
            </p>
        </header>

        <div class="am-terminal-page__toolbar">
            <div class="am-terminal-page__toolbar-left">
                <ElSelect
                    v-model="selectedAgent"
                    :disabled="isAgentEmpty"
                    :loading="loading"
                    :placeholder="t('agent.selectAgent')"
                    :no-data-text="t('agent.noData')"
                    filterable
                    class="am-terminal-page__host"
                >
                    <template #prefix>
                        <i-lucide-server class="am-terminal-page__host-icon" />
                    </template>
                    <ElOption
                        v-for="item in agentList"
                        :key="item.agent_id"
                        :label="item.hostname || item.agent_id"
                        :value="item.agent_id"
                    />
                </ElSelect>
                <span class="am-terminal-page__status" :data-tone="statusMeta.tone">
                    <span class="am-terminal-page__status-dot" :style="{ background: statusMeta.dot }" />
                    {{ statusMeta.text }}
                </span>
            </div>
            <div class="am-terminal-page__toolbar-right">
                <button
                    class="am-terminal-page__btn am-terminal-page__btn--ghost"
                    type="button"
                    :disabled="!actionable"
                    @click="handleClear"
                >
                    <i-lucide-eraser class="am-terminal-page__btn-icon" />
                    <span>{{ $t('agent.terminalClear') }}</span>
                </button>
                <button
                    class="am-terminal-page__btn am-terminal-page__btn--primary"
                    type="button"
                    :disabled="!actionable"
                    @click="handleNewSession"
                >
                    <i-lucide-plus class="am-terminal-page__btn-icon" />
                    <span>{{ $t('agent.terminalNewSession') }}</span>
                </button>
            </div>
        </div>

        <div class="am-terminal-page__panel">
            <div class="am-terminal-page__panel-titlebar">
                <span class="am-terminal-page__dot am-terminal-page__dot--red" />
                <span class="am-terminal-page__dot am-terminal-page__dot--yellow" />
                <span class="am-terminal-page__dot am-terminal-page__dot--green" />
                <span class="am-terminal-page__panel-title">{{ termTitle }}</span>
            </div>
            <div class="am-terminal-page__panel-body">
                <Terminal
                    v-if="agentId && !isAgentEmpty"
                    ref="terminal"
                    :agent-id="agentId"
                    @status-change="handleStatusChange"
                    @resize="handleTerminalResize"
                />
                <AgentEmptyState v-else-if="isAgentEmpty" min-height="360px" @refresh="loadAgents" />
                <div v-else class="am-terminal-page__empty">
                    {{ $t('agent.terminalSelectRequired') }}
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">
.am-terminal-page {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-md);
  height: 100%;
  min-height: 0;
  padding: var(--am-spacing-lg);
  box-sizing: border-box;
}

.am-terminal-page__header {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-xs);
  margin: 0;
}

.am-terminal-page__eyebrow {
  margin: 0;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
  font-weight: 600;
  letter-spacing: 0.16em;
}

.am-terminal-page__title {
  margin: 0;
  font-size: var(--am-font-2xl);
  font-weight: 700;
  color: var(--am-foreground-primary);
}

.am-terminal-page__hint {
  margin: 0;
  width: 100%;
  max-width: 560px;
  color: var(--am-foreground-muted);
  font-size: var(--am-font-sm);
  line-height: 1.5;
}

.am-terminal-page__toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--am-spacing-md);
  flex-wrap: wrap;
}

.am-terminal-page__toolbar-left,
.am-terminal-page__toolbar-right {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
}

.am-terminal-page__host {
  width: 220px;
}

.am-terminal-page__host :deep(.el-input__wrapper) {
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-subtle);
  border-radius: 6px;
  box-shadow: none;
  height: 32px;
  padding: 0 10px;
}

.am-terminal-page__host :deep(.el-input__wrapper):hover {
  border-color: var(--am-border-primary);
}

.am-terminal-page__host-icon {
  width: 14px;
  height: 14px;
  color: var(--am-foreground-secondary);
}

.am-terminal-page__btn-icon {
  width: 14px;
  height: 14px;
  flex: 0 0 auto;
}

.am-terminal-page__status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 6px;
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-subtle);
  font-size: var(--am-font-xs);
  color: var(--am-foreground-secondary);
}

.am-terminal-page__status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.am-terminal-page__btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 32px;
  padding: 0 12px;
  border-radius: 4px;
  font-size: var(--am-font-sm);
  font-weight: 600;
  cursor: pointer;
  transition:
    opacity 0.2s ease,
    background-color 0.2s ease;

  &:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }
}

.am-terminal-page__btn--ghost {
  color: var(--am-foreground-primary);
  background: var(--am-surface-card);
  border: 1px solid var(--am-border-primary);

  &:not(:disabled):hover {
    background: var(--am-surface-secondary);
  }
}

.am-terminal-page__btn--primary {
  color: var(--am-foreground-on-accent);
  background: var(--am-accent-primary);
  border: 1px solid var(--am-accent-primary);

  &:not(:disabled):hover {
    opacity: 0.9;
  }
}

.am-terminal-page__panel {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 20px;
  border-radius: 8px;
  overflow: hidden;
  background: #1e1e1e;
  border: 1px solid var(--am-border-subtle);
}

.am-terminal-page__panel-titlebar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0;
}

.am-terminal-page__dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex: 0 0 auto;

  &--red {
    background: #ff5f56;
  }
  &--yellow {
    background: #ffbd2e;
  }
  &--green {
    background: #27c93f;
  }
}

.am-terminal-page__panel-title {
  margin-left: 8px;
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: 12px;
  font-weight: 500;
  color: #6e6e6e;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.am-terminal-page__panel-body {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.am-terminal-page__empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--am-foreground-muted);
}

@media (max-width: 720px) {
  .am-terminal-page__toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .am-terminal-page__host {
    width: 100%;
  }
}
</style>
