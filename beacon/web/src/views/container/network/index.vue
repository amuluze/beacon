<script setup lang="ts">
import type { Network } from '@/interface/container.ts'
import type { TableInstance } from 'element-plus'

import { queryNetworks, removeNetwork } from '@/api/container'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { useConfirmCommand } from '@/hooks/useConfirmCommand'

import useCommandComponent from '@/hooks/useCommandComponent.ts'
import { useTable } from '@/hooks/useTable'
import AddNetwork from './components/AddNetwork.vue'
import DeleteNetwork from './components/DeleteNetwork.vue'
import BatchActionBar from '@/views/container/components/BatchActionBar.vue'
import TablePagination from '@/views/container/components/TablePagination.vue'

import { useI18n } from 'vue-i18n'

// 父级工作区用总数渲染 Tab 计数徽标
const emit = defineEmits<{
  total: [number]
}>()

const protectedNetworkNames = new Set(['bridge', 'host', 'none'])
function isProtectedNetwork(networkName: string) {
  return protectedNetworkNames.has(networkName)
}

const { tableData, pageable, loading, search, handleSizeChange, handleCurrentChange, updateTotalParam } = useTable(queryNetworks)
const { selectedAgentID, ensureSelectedAgent } = useAgentSelection({ immediate: false })
const initialized = ref(false)
onMounted(async () => {
  if (await ensureSelectedAgent())
    await search()
  initialized.value = true
  // 首次加载完成后才回填 Tab 计数，避免未查询时先报 0
  emit('total', pageable.value.total)
})
watch(selectedAgentID, async () => {
  if (initialized.value && selectedAgentID.value) {
    await search()
  }
})
watch(() => pageable.value.total, (total) => {
  if (initialized.value)
    emit('total', total)
})

// 搜索与驱动筛选在服务端生效：total 与列表使用同一组过滤条件
const keyword = ref('')
const driverFilter = ref('')

async function applyFilters() {
  pageable.value.page = 1
  await updateTotalParam({
    name: keyword.value.trim() || undefined,
    driver: driverFilter.value || undefined,
  })
}

const tableRef = ref<TableInstance>()
const tableSelection = ref<Network[]>([])
// Docker 内置网络不可删除，因此也不进入勾选集合
const selectable = (row: Network) => !isProtectedNetwork(row.name)
function handleSelectionChange(val: Network[]) {
  tableSelection.value = val
}

const selectedIDs = computed(() => tableSelection.value.map(row => row.id))

const addNetwork = useCommandComponent(AddNetwork)
const deleteNetwork = useCommandComponent(DeleteNetwork)

// 确认弹层在打开时读取 params，因此在触发前回填当前选中数量。
const batchParams = reactive({ count: 0 })

async function removeSelected(ids: string[]) {
  for (const id of ids)
    await removeNetwork({ network_id: id })
}

const batchDelete = useConfirmCommand({
  title: 'network.batchDelete',
  message: 'network.confirmBatchDelete',
  params: batchParams,
  i18nPrefix: 'network',
  action: async () => {
    await removeSelected(selectedIDs.value)
    tableRef.value?.clearSelection()
  },
  onResolved: () => {
    void search()
  },
})

const batchActions = [{ key: 'delete', label: 'network.batchDelete' }]

function handleBatchAction() {
  if (selectedIDs.value.length === 0)
    return
  batchParams.count = selectedIDs.value.length
  batchDelete()
}

const { t } = useI18n()
</script>

<template>
    <div class="am-container">
        <div class="am-table-operator">
            <div class="am-table-operator__left">
                <el-button type="primary" size="small" :disabled="!selectedAgentID" @click="addNetwork({ title: 'network.newNetwork', update: search })">
                    <svg-icon icon-class="add" />
                    {{ t('network.newNetwork') }}
                </el-button>
                <el-input
                    v-model="keyword"
                    class="am-table-operator__search"
                    size="small"
                    clearable
                    :placeholder="t('network.searchPlaceholder')"
                    @keyup.enter="applyFilters"
                    @clear="applyFilters"
                >
                    <template #prefix>
                        <svg-icon icon-class="search" />
                    </template>
                </el-input>
                <el-select v-model="driverFilter" class="am-table-operator__select" size="small" @change="applyFilters">
                    <el-option :label="t('network.driverAll')" value="" />
                    <el-option label="bridge" value="bridge" />
                    <el-option label="host" value="host" />
                    <el-option label="none" value="none" />
                </el-select>
            </div>
            <BatchActionBar
                :count="selectedIDs.length"
                :actions="batchActions"
                @select="handleBatchAction"
            />
        </div>
        <div class="am-table">
            <el-table
                ref="tableRef"
                v-loading="loading"
                :data="tableData as Network[]"
                :header-cell-style="{ height: '44px', fontSize: '13px', color: 'var(--am-foreground-primary)', background: 'var(--am-surface-secondary)' }"
                height="100%"
                border
                @selection-change="handleSelectionChange"
            >
                <el-table-column type="selection" :selectable="selectable" width="48" />
                <el-table-column prop="name" :label="t('network.networkName')" align="center" min-width="160" fixed />
                <el-table-column prop="driver" :label="t('network.networkMode')" align="center" min-width="120" />
                <el-table-column prop="subnet" :label="t('network.subNetwork')" align="center" show-overflow-tooltip min-width="140" />
                <el-table-column prop="gateway" :label="t('network.networkGateway')" align="center" show-overflow-tooltip min-width="160" />
                <el-table-column prop="created" :label="t('network.createTime')" align="center" min-width="200" />
                <el-table-column :label="t('network.operator')" width="160" fixed="right" align="center">
                    <template #default="scope">
                        <span v-if="isProtectedNetwork(scope.row.name)" class="am-row-note">
                            {{ t('container.builtinNetwork') }}
                        </span>
                        <el-button v-else type="danger" size="small" text :disabled="!selectedAgentID" @click="deleteNetwork({ title: 'network.deleteNetwork', id: scope.row.id, update: search })">
                            <svg-icon icon-class="delete" />
                            {{ t('network.delete') }}
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <TablePagination
            :pageable="pageable"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
        />
    </div>
</template>

<style scoped lang="scss">
.am-row-note {
  color: var(--am-foreground-muted);
  font-size: var(--am-font-sm);
}

.am-table-operator__search {
  width: 220px;
}

.am-table-operator__select {
  width: 140px;
}
</style>
