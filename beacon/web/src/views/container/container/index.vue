<script setup lang="ts">
import type { Container } from '@/interface/container.ts'
import type { TableInstance } from 'element-plus'

import { queryContainers, removeContainer, startContainer, stopContainer } from '@/api/container'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { useConfirmCommand } from '@/hooks/useConfirmCommand'

import useCommandComponent from '@/hooks/useCommandComponent.ts'
import { useTable } from '@/hooks/useTable.ts'
import AddContainer from './components/AddContainer.vue'
import DeleteContainer from './components/DeleteContainer.vue'
import EditContainer from './components/EditContainer.vue'
import RestartContainer from './components/RestartContainer.vue'
import StartContainer from './components/StartContainer.vue'
import StopContainer from './components/StopContainer.vue'
import { isProtectedContainerName } from './containerProtection'
import BatchActionBar from '@/views/container/components/BatchActionBar.vue'
import TablePagination from '@/views/container/components/TablePagination.vue'

import { useI18n } from 'vue-i18n'

// 父级工作区用总数渲染 Tab 计数徽标
const emit = defineEmits<{
  total: [number]
}>()

const { tableData, pageable, loading, search, handleSizeChange, handleCurrentChange, updateTotalParam } = useTable(queryContainers)
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

// 搜索与筛选在服务端生效：total 与列表使用同一组过滤条件
const keyword = ref('')
const stateFilter = ref('')

async function applyFilters() {
  pageable.value.page = 1
  await updateTotalParam({
    name: keyword.value.trim() || undefined,
    state: stateFilter.value || undefined,
  })
}

const tableRef = ref<TableInstance>()
const tableSelection = ref<Container[]>([])
// 受保护容器（beacon / amprobe）不允许批量操作，也不进入勾选集合
const selectable = (row: Container) => !isProtectedContainerName(row.name)

function handleSelectionChange(val: Container[]) {
  tableSelection.value = val
}

const selectedIDs = computed(() => tableSelection.value.map(row => row.id))

const addContainer = useCommandComponent(AddContainer)
const startContainerCommand = useCommandComponent(StartContainer)
const stopContainerCommand = useCommandComponent(StopContainer)
const restartContainer = useCommandComponent(RestartContainer)
const deleteContainer = useCommandComponent(DeleteContainer)
const editContainer = useCommandComponent(EditContainer)

async function runBatch(ids: string[], run: (containerID: string) => Promise<unknown>) {
  for (const id of ids)
    await run(id)
}

// 确认弹层在打开时读取 params，因此在触发前回填当前选中数量。
// 批量操作按单容器接口串行执行，失败即中断，避免在半途静默跳过目标。
const batchParams = reactive({ count: 0 })

function createBatchCommand(title: string, message: string, run: (containerID: string) => Promise<unknown>) {
  return useConfirmCommand({
    title,
    message,
    params: batchParams,
    i18nPrefix: 'container',
    action: async () => {
      await runBatch(selectedIDs.value, run)
      tableRef.value?.clearSelection()
    },
    onResolved: () => {
      void search()
    },
  })
}

const batchStart = createBatchCommand('container.batchStart', 'container.confirmBatchStart', id => startContainer({ container_id: id }))
const batchStop = createBatchCommand('container.batchStop', 'container.confirmBatchStop', id => stopContainer({ container_id: id }))
const batchDelete = createBatchCommand('container.batchDelete', 'container.confirmBatchDelete', id => removeContainer({ container_id: id }))

const batchActions = [
  { key: 'start', label: 'container.batchStart' },
  { key: 'stop', label: 'container.batchStop' },
  { key: 'delete', label: 'container.batchDelete' },
]

function handleBatchAction(key: string) {
  if (selectedIDs.value.length === 0)
    return
  batchParams.count = selectedIDs.value.length
  if (key === 'start')
    batchStart()
  else if (key === 'stop')
    batchStop()
  else if (key === 'delete')
    batchDelete()
}

const { t } = useI18n()
</script>

<template>
    <div class="am-container">
        <div class="am-table-operator">
            <div class="am-table-operator__left">
                <el-button type="primary" size="small" :disabled="!selectedAgentID" @click="addContainer({ title: 'container.addContainer', update: search })">
                    <svg-icon icon-class="add" />
                    {{ t('container.addContainer') }}
                </el-button>
                <el-input
                    v-model="keyword"
                    class="am-table-operator__search"
                    size="small"
                    clearable
                    :placeholder="t('container.searchPlaceholder')"
                    @keyup.enter="applyFilters"
                    @clear="applyFilters"
                >
                    <template #prefix>
                        <svg-icon icon-class="search" />
                    </template>
                </el-input>
                <el-select v-model="stateFilter" class="am-table-operator__select" size="small" @change="applyFilters">
                    <el-option :label="t('container.stateAll')" value="" />
                    <el-option :label="t('container.enable')" value="running" />
                    <el-option :label="t('container.disable')" value="exited" />
                </el-select>
            </div>
            <BatchActionBar
                :count="selectedIDs.length"
                :actions="batchActions"
                @select="handleBatchAction"
            />
        </div>
        <!-- 表格主体 -->
        <div class="am-table">
            <el-table
                ref="tableRef"
                v-loading="loading"
                :data="tableData as Container[]"
                :header-cell-style="{ height: '44px', fontSize: '13px', color: 'var(--am-foreground-primary)', background: 'var(--am-surface-secondary)' }"
                height="100%"
                border
                @selection-change="handleSelectionChange"
            >
                <el-table-column type="selection" :selectable="selectable" width="48" />
                <el-table-column prop="name" :label="t('container.containerName')" align="center" min-width="180" show-overflow-tooltip sortable fixed />
                <el-table-column prop="image" :label="t('container.imageName')" align="center" min-width="180" show-overflow-tooltip sortable />
                <el-table-column prop="cpu_percent" :label="t('container.cpuPercent')" align="center" min-width="110" show-overflow-tooltip />
                <el-table-column prop="memory_usage" :label="t('container.memUsed')" align="center" min-width="120" show-overflow-tooltip />
                <el-table-column prop="ip" :label="t('container.containerIP')" align="center" min-width="120" show-overflow-tooltip />
                <el-table-column prop="ports" :label="t('container.containerPort')" align="center" min-width="140" show-overflow-tooltip />
                <el-table-column prop="state" :label="t('container.state')" align="center" min-width="110" show-overflow-tooltip>
                    <template #default="scope">
                        <el-tag v-if="scope.row.state === 'running'" type="success">
                            {{ t('container.enable') }}
                        </el-tag>
                        <el-tag v-else type="danger">
                            {{ t('container.disable') }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="uptime" :label="t('container.uptime')" align="center" min-width="150" show-overflow-tooltip />
                <el-table-column :label="t('container.operator')" width="190" fixed="right" align="center">
                    <template #default="scope">
                        <span v-if="isProtectedContainerName(scope.row.name)" class="am-row-note">
                            {{ t('container.protected') }}
                        </span>
                        <template v-else>
                            <!-- 主操作跟随运行状态：运行中只能停止，已停止只能启动 -->
                            <el-button
                                v-if="scope.row.state === 'running'"
                                type="warning"
                                size="small"
                                text
                                :disabled="!selectedAgentID"
                                @click="stopContainerCommand({ title: 'container.stopContainer', id: scope.row.id, update: search })"
                            >
                                <svg-icon icon-class="stop" />
                                {{ t('container.stop') }}
                            </el-button>
                            <el-button
                                v-else
                                type="primary"
                                size="small"
                                text
                                :disabled="!selectedAgentID"
                                @click="startContainerCommand({ title: 'container.startContainer', id: scope.row.id, update: search })"
                            >
                                <svg-icon icon-class="start" />
                                {{ t('container.start') }}
                            </el-button>
                            <el-dropdown>
                                <el-button type="primary" size="small" text :disabled="!selectedAgentID">
                                    <svg-icon icon-class="more" />
                                    {{ t('container.more') }}
                                </el-button>
                                <template #dropdown>
                                    <el-dropdown-menu>
                                        <el-dropdown-item>
                                            <el-button type="primary" size="small" text :disabled="!selectedAgentID" @click="editContainer({ title: 'container.editContainer', container: scope.row, update: search })">
                                                <svg-icon icon-class="edit" />
                                                {{ t('container.edit') }}
                                            </el-button>
                                        </el-dropdown-item>
                                        <el-dropdown-item>
                                            <el-button type="warning" size="small" text :disabled="!selectedAgentID" @click="restartContainer({ title: 'container.restartContainer', id: scope.row.id, update: search })">
                                                <svg-icon icon-class="update" />
                                                {{ t('container.restart') }}
                                            </el-button>
                                        </el-dropdown-item>
                                        <el-dropdown-item>
                                            <el-button type="danger" size="small" text :disabled="!selectedAgentID" @click="deleteContainer({ title: 'container.deleteContainer', id: scope.row.id, update: search })">
                                                <svg-icon icon-class="delete" />
                                                {{ t('container.delete') }}
                                            </el-button>
                                        </el-dropdown-item>
                                    </el-dropdown-menu>
                                </template>
                            </el-dropdown>
                        </template>
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
