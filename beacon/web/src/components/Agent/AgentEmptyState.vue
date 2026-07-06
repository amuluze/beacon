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
}

.agent-empty-state__description {
  margin-bottom: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
</style>
