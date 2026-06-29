<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const selectedAgent = ref((route.query.agent_id as string) || '')

const agentId = computed(() => selectedAgent.value)

function handleAgentChange(value: string): void {
  selectedAgent.value = value
}
</script>

<template>
    <div class="am-terminal-page">
        <ContentWrap title="Web Terminal" message="选择 Agent 并打开远程终端">
            <div class="am-terminal-page__toolbar">
                <span class="am-terminal-page__label">Agent:</span>
                <ElInput
                    v-model="selectedAgent"
                    placeholder="请输入 Agent ID"
                    clearable
                    style="width: 240px"
                    @change="handleAgentChange"
                />
            </div>
            <div class="am-terminal-page__container">
                <Terminal v-if="agentId" :agent-id="agentId" />
                <div v-else class="am-terminal-page__empty">
                    请选择或输入 Agent ID
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
