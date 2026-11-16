<script setup lang="ts">
import { queryAudit } from '@/api/audit'
import type { Audit } from '@/interface/audit'

// 审计页展示一条合并时间线：既包含 Agent 系统记录，也包含登录、用户、
// 网络、镜像与设置变更等操作员记录。后端 type 留空即返回两类记录。
const rows = ref<Audit[]>([])
const loading = shallowRef(false)
const page = shallowRef(1)
const size = shallowRef(10)
const total = shallowRef(0)

async function loadAudit() {
  loading.value = true
  try {
    const { data } = await queryAudit({
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

async function handleSizeChange(value: number) {
  size.value = value
  page.value = 1
  await loadAudit()
}

async function handlePageChange(value: number) {
  page.value = value
  await loadAudit()
}

onMounted(loadAudit)
</script>

<template>
    <div class="audit-log">
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
