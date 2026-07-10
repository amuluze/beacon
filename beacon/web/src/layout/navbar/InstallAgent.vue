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
  </button>
</template>

<style scoped lang="scss">
.am-install-trigger {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
  width: 32px;
  height: 32px;
  padding: 0;
  color: var(--am-foreground-on-accent);
  background: var(--am-accent-primary);
  border: 0;
  border-radius: 4px;
  cursor: pointer;
  transition: opacity 0.2s ease;

  &:hover {
    opacity: 0.9;
  }
}
</style>
