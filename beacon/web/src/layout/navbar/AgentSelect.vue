<script setup lang="ts">
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { useI18n } from 'vue-i18n'

const { agentList, selectedAgentID, loading, isAgentEmpty } = useAgentSelection({ immediate: false })
const { t } = useI18n()
</script>

<template>
    <el-select
        v-model="selectedAgentID"
        :loading="loading"
        :disabled="isAgentEmpty"
        size="small"
        style="width: 220px"
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
