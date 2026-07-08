<script setup lang="ts">
import { querySystemAudit } from '@/api/audit'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import { useI18n } from 'vue-i18n'
import useStore from '@/store'
import en from 'element-plus/es/locale/lang/en'
import { getBrowserLanguage } from '@/utils'
import type { AuditQueryResult } from '@/interface/audit.ts'

const tableData = ref<AuditQueryResult['data']>([])
const pageable = reactive({ page: 1, size: 10, total: 0, options: [10, 20, 50, 100, 200] })
const loading = ref(false)
const selectedAgentID = ref<string>('')

async function reload() {
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      type: 'system',
      page: pageable.page,
      size: pageable.size,
    }
    if (selectedAgentID.value) {
      params.agent_id = selectedAgentID.value
    }
    const { data } = await querySystemAudit(params as any)
    tableData.value = data.data
    pageable.total = data.total
  }
  finally {
    loading.value = false
  }
}

function handleSizeChange(size: number) {
  pageable.size = size
  pageable.page = 1
  reload()
}
function handleCurrentChange(page: number) {
  pageable.page = page
  reload()
}

watch(selectedAgentID, () => {
  pageable.page = 1
  reload()
})

onMounted(reload)

const { t } = useI18n()
const store = useStore()
const locale = computed(() => {
  if (store.app.language === 'zh')
    return zhCn
  if (store.app.language === 'en')
    return en
  return getBrowserLanguage() === 'zh' ? zhCn : en
})
</script>

<template>
    <div class="am-container">
        <div class="am-table-operator">
            <el-select
                v-model="selectedAgentID"
                clearable
                :placeholder="t('audit.filterByAgent') || '按 Agent 过滤'"
                size="small"
                style="width: 240px"
            >
                <el-option
                    v-for="agent in store.agent.list"
                    :key="agent.agent_id"
                    :label="agent.hostname || agent.agent_id"
                    :value="agent.agent_id"
                />
            </el-select>
        </div>
        <div class="am-table">
            <el-table
                v-loading="loading"
                :data="tableData"
                :header-cell-style="{ height: '45px', fontSize: '14px', color: '#000', background: '#fafafa' }"
                height="100%"
                border
            >
                <el-table-column prop="id" label="ID" align="center" min-width="150" />
                <el-table-column prop="username" :label="t('audit.username')" align="center" min-width="150" />
                <el-table-column prop="agent_id" :label="t('audit.agentID') || 'Agent'" align="center" min-width="160" />
                <el-table-column prop="operate" :label="t('audit.operate')" align="center" show-overflow-tooltip min-width="150" />
                <el-table-column prop="created" :label="t('audit.operateTime')" align="center" min-width="200" />
            </el-table>
        </div>
        <div class="am-pagination">
            <el-config-provider :locale="locale">
                <el-pagination
                    v-model:current-page="pageable.page"
                    :page-size="pageable.size"
                    layout="total, sizes, prev, pager, next, jumper"
                    :page-sizes="pageable.options"
                    :total="pageable.total"
                    @size-change="handleSizeChange"
                    @current-change="handleCurrentChange"
                />
            </el-config-provider>
        </div>
    </div>
</template>

<style scoped lang="scss">
</style>
