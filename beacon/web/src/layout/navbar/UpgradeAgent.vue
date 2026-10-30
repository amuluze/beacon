<script setup lang="ts">
import UpgradeAgentDialog from '@/components/Agent/UpgradeAgentDialog.vue'
import useCommandComponent from '@/hooks/useCommandComponent.ts'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { useI18n } from 'vue-i18n'
import useStore from '@/store'

const { t } = useI18n()
const store = useStore()
const { agentList, selectedAgentID, loadAgents } = useAgentSelection({ immediate: false })

// 升级 Collia 是管理员专属资源，普通用户不暴露入口，避免点击后触发 403。
const isAdmin = computed(() => store.user.userInfo.isAdmin === 1)
const selectedAgent = computed(() => agentList.value.find(item => item.agent_id === selectedAgentID.value))
const hasTarget = computed(() => Boolean(selectedAgentID.value))

const upgradeAgentDialog = useCommandComponent(UpgradeAgentDialog)

function openUpgrade() {
  if (!selectedAgentID.value) {
    return
  }
  upgradeAgentDialog({
    title: 'agent.upgradeTitle',
    agentID: selectedAgentID.value,
    hostname: selectedAgent.value?.hostname || '',
    version: selectedAgent.value?.version || '',
    onSuccess: () => {
      // 升级后 Agent 版本会变化，刷新列表让顶栏选择器同步最新版本。
      void loadAgents()
    },
  })
}
</script>

<template>
    <button
        v-if="isAdmin"
        class="am-upgrade-trigger"
        type="button"
        :disabled="!hasTarget"
        :title="hasTarget ? t('agent.upgradeTitle') : t('agent.upgradeNoAgent')"
        @click="openUpgrade"
    >
        <svg-icon icon-class="update" size="14px" />
        <span class="am-upgrade-trigger__label">{{ t('agent.upgradeAction') }}</span>
    </button>
</template>

<style scoped lang="scss">
.am-upgrade-trigger {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
  gap: var(--am-spacing-xs);
  height: 32px;
  padding: 0 12px;
  color: var(--am-foreground-secondary);
  background: transparent;
  border: 1px solid var(--am-border-primary);
  border-radius: 4px;
  font-size: var(--am-font-sm);
  font-weight: 600;
  white-space: nowrap;
  cursor: pointer;
  transition:
    color 0.2s ease,
    border-color 0.2s ease,
    opacity 0.2s ease;

  &:hover:not(:disabled) {
    color: var(--am-foreground-primary);
    border-color: var(--am-accent-primary);
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &__label {
    white-space: nowrap;
  }
}

@media (max-width: 720px) {
  .am-upgrade-trigger {
    padding: 0;
    width: 32px;

    &__label {
      display: none;
    }
  }
}
</style>
