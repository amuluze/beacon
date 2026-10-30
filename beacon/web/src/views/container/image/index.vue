<script setup lang="ts">
import type { Image } from '@/interface/container.ts'
import type { TableInstance } from 'element-plus'

import { queryImages, removeImage } from '@/api/container'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { useConfirmCommand } from '@/hooks/useConfirmCommand'

import useCommandComponent from '@/hooks/useCommandComponent.ts'
import { useTable } from '@/hooks/useTable'
import DeleteImage from './components/DeleteImage.vue'
import ImportImage from './components/ImportImage.vue'
import PruneImage from './components/PruneImage.vue'
import PullImage from './components/PullImage.vue'
import BatchActionBar from '@/views/container/components/BatchActionBar.vue'
import TablePagination from '@/views/container/components/TablePagination.vue'

import { useI18n } from 'vue-i18n'

// 父级工作区用总数渲染 Tab 计数徽标
const emit = defineEmits<{
  total: [number]
}>()

const { tableData, pageable, loading, search, handleSizeChange, handleCurrentChange, updateTotalParam } = useTable(queryImages)
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

// 搜索在服务端生效：total 与列表使用同一组过滤条件
const keyword = ref('')

async function applyFilters() {
  pageable.value.page = 1
  await updateTotalParam({ name: keyword.value.trim() || undefined })
}

const tableRef = ref<TableInstance>()
const tableSelection = ref<Image[]>([])
// 仍被容器引用的镜像不允许删除，因此也不进入勾选集合
const selectable = (row: Image) => row.number === 0
function handleSelectionChange(val: Image[]) {
  tableSelection.value = val
}

const selectedIDs = computed(() => tableSelection.value.map(row => row.id))

const pullImage = useCommandComponent(PullImage)
const importImage = useCommandComponent(ImportImage)
const pruneImage = useCommandComponent(PruneImage)
const deleteImage = useCommandComponent(DeleteImage)

// 确认弹层在打开时读取 params，因此在触发前回填当前选中数量。
const batchParams = reactive({ count: 0 })

async function removeSelected(ids: string[]) {
  for (const id of ids)
    await removeImage({ image_id: id })
}

const batchDelete = useConfirmCommand({
  title: 'image.batchDelete',
  message: 'image.confirmBatchDelete',
  params: batchParams,
  i18nPrefix: 'image',
  action: async () => {
    await removeSelected(selectedIDs.value)
    tableRef.value?.clearSelection()
  },
  onResolved: () => {
    void search()
  },
})

const batchActions = [{ key: 'delete', label: 'image.batchDelete' }]

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
            <!-- 每个 Tab 只有一个实心主操作，其余降级为普通/警示按需出现 -->
            <div class="am-table-operator__left">
                <el-button type="primary" size="small" :disabled="!selectedAgentID" @click="pullImage({ title: 'image.pullImage', imageName: '', update: search })">
                    <svg-icon icon-class="download" />
                    {{ t('image.pullImage') }}
                </el-button>
                <el-button plain size="small" :disabled="!selectedAgentID" @click="importImage({ title: 'image.importImage', update: search })">
                    <svg-icon icon-class="upload" />
                    {{ t('image.importImage') }}
                </el-button>
                <el-button type="warning" plain size="small" :disabled="!selectedAgentID" @click="pruneImage({ title: 'image.pruneImage', update: search })">
                    <svg-icon icon-class="delete" />
                    {{ t('image.pruneImage') }}
                </el-button>
                <el-input
                    v-model="keyword"
                    class="am-table-operator__search"
                    size="small"
                    clearable
                    :placeholder="t('image.searchPlaceholder')"
                    @keyup.enter="applyFilters"
                    @clear="applyFilters"
                >
                    <template #prefix>
                        <svg-icon icon-class="search" />
                    </template>
                </el-input>
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
                :data="tableData as Image[]"
                :header-cell-style="{ height: '44px', fontSize: '13px', color: 'var(--am-foreground-primary)', background: 'var(--am-surface-secondary)' }"
                height="100%"
                border
                @selection-change="handleSelectionChange"
            >
                <el-table-column type="selection" :selectable="selectable" width="48" />
                <el-table-column prop="name" :label="t('image.imageName')" align="center" min-width="260" show-overflow-tooltip fixed sortable />
                <el-table-column prop="number" :label="t('image.containerNum')" align="center" show-overflow-tooltip min-width="120" sortable />
                <el-table-column prop="tag" :label="t('image.imageTag')" align="center" show-overflow-tooltip min-width="120" />
                <el-table-column prop="created" :label="t('image.createTime')" align="center" min-width="180" />
                <el-table-column prop="size" :label="t('image.imageSize')" align="center" min-width="120" />
                <el-table-column :label="t('image.operator')" width="160" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="danger" size="small" text :disabled="!selectedAgentID || scope.row.number > 0" @click="deleteImage({ title: 'image.deleteImage', id: scope.row.id, update: search })">
                            <svg-icon icon-class="delete" />
                            {{ t('image.delete') }}
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
.am-table-operator__search {
  width: 240px;
}
</style>
