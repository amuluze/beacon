<script setup lang="ts">
import { queryAgentList } from '@/api/agent'
import useStore from '@/store'
import { Refresh } from '@element-plus/icons-vue'

const store = useStore()
const loading = ref(false)

const currentAgent = computed({
  get: () => store.agent.currentAgentID,
  set: (value: string) => {
    store.agent.setCurrentAgent(value)
  },
})

async function loadAgents() {
  loading.value = true
  try {
    const { data } = await queryAgentList()
    store.agent.setAgents(data || [])
  }
  catch {
    store.agent.setAgents([])
  }
  finally {
    loading.value = false
  }
}

onMounted(() => {
  loadAgents()
})
</script>

<template>
  <div class="am-agent-select">
    <el-select
      v-model="currentAgent"
      :loading="loading"
      size="small"
      placeholder="Agent"
      style="width: 180px"
      filterable
    >
      <el-option
        v-for="item in store.agent.agents"
        :key="item.agent_id"
        :label="item.hostname || item.agent_id"
        :value="item.agent_id"
      >
        <span>{{ item.hostname || item.agent_id }}</span>
        <span class="am-agent-select__status">{{ item.status }}</span>
      </el-option>
    </el-select>
    <el-button :icon="Refresh" :loading="loading" size="small" text @click="loadAgents" />
  </div>
</template>

<style scoped lang="scss">
@include b(agent-select) {
  display: flex;
  align-items: center;
  gap: 4px;

  @include e(status) {
    float: right;
    margin-left: 12px;
    color: var(--el-text-color-secondary);
    font-size: 12px;
  }
}
</style>
