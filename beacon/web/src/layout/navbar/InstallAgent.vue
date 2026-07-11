<script setup lang="ts">
import InstallAgentDialog from '@/components/Agent/InstallAgentDialog.vue'
import useCommandComponent from '@/hooks/useCommandComponent.ts'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const installAgent = useCommandComponent(InstallAgentDialog)

function openInstall() {
  installAgent({ title: 'agent.installTitle' })
}
</script>

<template>
  <button class="am-install-trigger" type="button" :title="t('agent.installTitle')" @click="openInstall">
    <svg-icon icon-class="download" size="14px" />
    <span class="am-install-trigger__label">{{ t('agent.installAction') }}</span>
  </button>
</template>

<style scoped lang="scss">
.am-install-trigger {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
  gap: var(--am-spacing-xs);
  height: 32px;
  padding: 0 12px;
  color: var(--am-foreground-on-accent);
  background: var(--am-accent-primary);
  border: 0;
  border-radius: 4px;
  font-size: var(--am-font-sm);
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s ease;

  &:hover {
    opacity: 0.9;
  }

  &__label {
    white-space: nowrap;
  }
}

@media (max-width: 720px) {
  .am-install-trigger {
    padding: 0;
    width: 32px;

    &__label {
      display: none;
    }
  }
}
</style>
