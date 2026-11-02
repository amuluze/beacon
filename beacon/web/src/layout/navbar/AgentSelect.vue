<script setup lang="ts">
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { useI18n } from 'vue-i18n'

// 顶栏使用 220px 宽，工作区分组内的上下文条使用紧凑宽度。
const props = withDefaults(defineProps<{ width?: string }>(), {
  width: '220px',
})
const { agentList, selectedAgentID, loading, isAgentEmpty } = useAgentSelection({ immediate: false })
const { t } = useI18n()
</script>

<template>
    <el-select
        v-model="selectedAgentID"
        :loading="loading"
        :disabled="isAgentEmpty"
        size="small"
        :style="{ width: props.width }"
        :placeholder="t('agent.selectAgent')"
        :no-data-text="t('agent.noData')"
        filterable
    >
        <el-option
            v-for="item in agentList"
            :key="item.agent_id"
            :label="item.hostname || item.agent_id"
            :value="item.agent_id"
        >
            <span>{{ item.hostname || item.agent_id }}</span>
            <span style="float: right; color: var(--el-text-color-secondary); font-size: 12px">{{ item.status || item.version }}</span>
        </el-option>
    </el-select>
</template>
