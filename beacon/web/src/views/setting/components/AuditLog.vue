<script setup lang="ts">
import { queryAudit } from '@/api/audit'
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import type { Audit } from '@/interface/audit'

const { selectedAgentID, isAgentEmpty, ensureSelectedAgent, loadAgents } = useAgentSelection({ immediate: false })
const rows = ref<Audit[]>([])
const loading = shallowRef(false)
const initialized = shallowRef(false)
const page = shallowRef(1)
const size = shallowRef(10)
const total = shallowRef(0)

async function loadAudit() {
  if (!selectedAgentID.value) {
    rows.value = []
    total.value = 0
    return
  }

  loading.value = true
  try {
    const { data } = await queryAudit({
      type: 'system',
      agent_id: selectedAgentID.value,
      page: page.value,
      size: size.value,
    })
    rows.value = data.data ?? []
    total.value = data.total
  }
  finally {
    loading.value = false
  }
}

async function refreshAgents() {
  await loadAgents()
  page.value = 1
  await loadAudit()
}

async function handleSizeChange(value: number) {
  size.value = value
  page.value = 1
  await loadAudit()
}

async function handlePageChange(value: number) {
  page.value = value
  await loadAudit()
}

watch(selectedAgentID, async () => {
  if (!initialized.value)
    return
  page.value = 1
  await loadAudit()
})

onMounted(async () => {
  await ensureSelectedAgent()
  await loadAudit()
  initialized.value = true
})
</script>

<template>
    <AgentEmptyState v-if="isAgentEmpty" min-height="280px" @refresh="refreshAgents" />
    <div v-else class="audit-log">
        <el-table
            v-loading="loading"
            :data="rows"
            :header-cell-style="{ height: '44px', fontSize: '13px', color: 'var(--am-foreground-primary)', background: 'var(--am-surface-secondary)' }"
            height="360px"
            border
        >
            <el-table-column prop="id" label="ID" width="90" />
            <el-table-column prop="username" :label="$t('audit.username')" min-width="140" />
            <el-table-column prop="agent_id" :label="$t('audit.agentID')" min-width="160" show-overflow-tooltip />
            <el-table-column prop="operate" :label="$t('audit.operate')" min-width="280" show-overflow-tooltip />
            <el-table-column prop="created" :label="$t('audit.operateTime')" min-width="180" />
        </el-table>
        <div class="audit-log__pagination">
            <el-pagination
                v-model:current-page="page"
                :page-size="size"
                :page-sizes="[10, 20, 50, 100]"
                :total="total"
                layout="total, sizes, prev, pager, next"
                @size-change="handleSizeChange"
                @current-change="handlePageChange"
            />
        </div>
    </div>
</template>

<style scoped lang="scss">
.audit-log {
  display: flex;
  flex-direction: column;
  gap: var(--am-spacing-md);
}

.audit-log__pagination {
  display: flex;
  justify-content: flex-end;
  overflow-x: auto;
}
</style>
