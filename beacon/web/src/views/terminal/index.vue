<script setup lang="ts">
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import { computed, onMounted } from 'vue'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const { agentList, selectedAgentID, loading, isAgentEmpty, loadAgents, ensureSelectedAgent } = useAgentSelection({ immediate: false })
const routeAgentID = computed(() => route.query.agent_id as string | undefined)
const { t } = useI18n()

const selectedAgent = computed({
  get: () => selectedAgentID.value,
  set: (value: string) => selectedAgentID.value = value,
})
const agentId = computed(() => selectedAgent.value)

function handleAgentChange(value: string): void {
  selectedAgent.value = value
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
        <ContentWrap :title="t('agent.terminalTitle')" :message="t('agent.terminalMessage')">
            <div class="am-terminal-page__toolbar">
                <span class="am-terminal-page__label">Agent:</span>
                <ElSelect
                    v-model="selectedAgent"
                    :disabled="isAgentEmpty"
                    :loading="loading"
                    :placeholder="t('agent.selectAgent')"
                    :no-data-text="t('agent.noData')"
                    filterable
                    style="width: 240px"
                    @change="handleAgentChange"
                >
                    <ElOption
                        v-for="item in agentList"
                        :key="item.agent_id"
                        :label="item.hostname || item.agent_id"
                        :value="item.agent_id"
                    />
                </ElSelect>
            </div>
            <div class="am-terminal-page__container">
                <Terminal v-if="agentId && !isAgentEmpty" :agent-id="agentId" />
                <AgentEmptyState v-else-if="isAgentEmpty" min-height="360px" @refresh="loadAgents" />
                <div v-else class="am-terminal-page__empty">
                    {{ t('agent.terminalSelectRequired') }}
                </div>
            </div>
        </ContentWrap>
    </div>
</template>

<style scoped lang="scss">
.am-terminal-page {
  display: flex;
  flex-direction: column;
  height: 100%;

  &__toolbar {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
  }

  &__label {
    font-weight: 500;
  }

  &__container {
    flex: 1;
    min-height: 0;
  }

  &__empty {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: var(--el-text-color-secondary);
  }
}
</style>
