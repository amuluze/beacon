<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import IconServerOff from '~icons/lucide/server-off'
import InstallAgentDialog from '@/components/Agent/InstallAgentDialog.vue'

interface Props {
  title?: string
  description?: string
  minHeight?: string
  showRefresh?: boolean
  showInstall?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  minHeight: '320px',
  showRefresh: true,
  showInstall: true,
})

const emit = defineEmits<{
  refresh: []
}>()

const { t } = useI18n()
const emptyTitle = computed(() => props.title || t('agent.emptyTitle'))
const emptyDescription = computed(() => props.description || t('agent.emptyDescription'))
const installVisible = ref(false)
</script>

<template>
    <div class="agent-empty-state" :style="{ minHeight }">
        <div class="agent-empty-state__inner">
            <span class="agent-empty-state__icon">
                <IconServerOff />
            </span>
            <p class="agent-empty-state__title">
                {{ emptyTitle }}
            </p>
            <p class="agent-empty-state__description">
                {{ emptyDescription }}
            </p>
            <div class="agent-empty-state__actions">
                <el-button v-if="showInstall" type="primary" @click="installVisible = true">
                    <svg-icon icon-class="download" />
                    {{ t('agent.installAction') }}
                </el-button>
                <el-button v-if="showRefresh" type="primary" @click="emit('refresh')">
                    <svg-icon icon-class="update" />
                    {{ t('agent.refresh') }}
                </el-button>
            </div>
        </div>
        <InstallAgentDialog v-model:visible="installVisible" title="agent.installTitle" />
    </div>
</template>

<style scoped lang="scss">
.agent-empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  padding: var(--am-spacing-xl);
  color: var(--am-foreground-secondary);
  background: var(--am-surface-card);
  border: 1px dashed var(--am-border-subtle);
  border-radius: var(--am-radius-md);
  box-shadow: var(--am-shadow-subtle);
}

.agent-empty-state__inner {
  width: 100%;
  max-width: 480px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--am-spacing-sm);
  text-align: center;
}

.agent-empty-state__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--am-foreground-muted);

  svg {
    width: 28px;
    height: 28px;
  }
}

.agent-empty-state__title {
  margin: 0;
  color: var(--am-foreground-primary);
  font-size: var(--am-font-lg);
  font-weight: 600;
}

.agent-empty-state__description {
  max-width: 480px;
  margin: 0;
  color: var(--am-foreground-secondary);
  font-size: var(--am-font-sm);
  line-height: 1.5;
}

.agent-empty-state__actions {
  display: flex;
  align-items: center;
  gap: var(--am-spacing-sm);
}
</style>
