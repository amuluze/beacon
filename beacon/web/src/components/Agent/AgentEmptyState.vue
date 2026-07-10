<script setup lang="ts">
import { useI18n } from 'vue-i18n'

interface Props {
  title?: string
  description?: string
  minHeight?: string
  showRefresh?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  minHeight: '320px',
  showRefresh: true,
})

const emit = defineEmits<{
  refresh: []
}>()

const { t } = useI18n()
const emptyTitle = computed(() => props.title || t('agent.emptyTitle'))
const emptyDescription = computed(() => props.description || t('agent.emptyDescription'))
</script>

<template>
    <div class="agent-empty-state" :style="{ minHeight }">
        <el-empty :description="emptyTitle">
            <div class="agent-empty-state__description">
                {{ emptyDescription }}
            </div>
            <el-button v-if="showRefresh" type="primary" plain size="small" @click="emit('refresh')">
                <svg-icon icon-class="update" />
                {{ t('agent.refresh') }}
            </el-button>
        </el-empty>
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

.agent-empty-state__description {
  margin-bottom: var(--am-spacing-md);
  color: var(--am-foreground-muted);
  font-size: var(--am-font-sm);
}
</style>
