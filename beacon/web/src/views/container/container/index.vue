<script setup lang="ts">
import type { Container } from '@/interface/container.ts'
import type { TableInstance } from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

import { queryContainers } from '@/api/container'
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import { useTable } from '@/hooks/useTable.ts'
import { useAgentSelection } from '@/hooks/useAgentSelection'

import useCommandComponent from '@/hooks/useCommandComponent.ts'
import AddContainer from '@/views/container/container/components/AddContainer.vue'
import DeleteContainer from '@/views/container/container/components/DeleteContainer.vue'
import EditContainer from '@/views/container/container/components/EditContainer.vue'
import RestartContainer from '@/views/container/container/components/RestartContainer.vue'
import StartContainer from '@/views/container/container/components/StartContainer.vue'
import StopContainer from '@/views/container/container/components/StopContainer.vue'
import { isProtectedContainerName } from '@/views/container/container/containerProtection'

import useStore from '@/store'
import { getBrowserLanguage } from '@/utils'
import en from 'element-plus/es/locale/lang/en'
import { useI18n } from 'vue-i18n'

const { tableData, pageable, loading, search, handleSizeChange, handleCurrentChange } = useTable(queryContainers)
const { selectedAgentID, isAgentEmpty, ensureSelectedAgent, loadAgents } = useAgentSelection({ immediate: false })
const initialized = ref(false)
onMounted(async () => {
  if (await ensureSelectedAgent())
    await search()
  initialized.value = true
})
watch(selectedAgentID, async () => {
  if (initialized.value && selectedAgentID.value) {
    await search()
  }
})

const tableRef = ref<TableInstance>()
const tableSelection = ref<Container[]>([])
const selectable = (row: Container) => !['1', '2'].includes(row.id)

function handleSelectionChange(val: Container[]) {
  tableSelection.value = val
}
async function refreshAgents() {
  await loadAgents()
  if (selectedAgentID.value)
    await search()
}

const addContainer = useCommandComponent(AddContainer)
const startContainer = useCommandComponent(StartContainer)
const stopContainer = useCommandComponent(StopContainer)
const restartContainer = useCommandComponent(RestartContainer)
const deleteContainer = useCommandComponent(DeleteContainer)
const editContainer = useCommandComponent(EditContainer)

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
            <el-button type="primary" plain size="small" :disabled="!selectedAgentID" @click="addContainer({ title: 'container.addContainer', update: search })">
                <svg-icon icon-class="add" />
                {{ t('container.addContainer') }}
            </el-button>
        </div>
        <!-- 表格主体 -->
        <AgentEmptyState v-if="isAgentEmpty" @refresh="refreshAgents" />
        <div v-else class="am-table">
            <el-table
                ref="tableRef"
                v-loading="loading"
                :data="tableData as Container[]"
                :header-cell-style="{ height: '44px', fontSize: '13px', color: 'var(--am-foreground-primary)', background: 'var(--am-surface-secondary)' }"
                height="100%"
                border
                @selection-change="handleSelectionChange"
            >
                <el-table-column type="selection" :selection="selectable" width="55" />
                <el-table-column prop="name" :label="t('container.containerName')" align="center" min-width="180" show-overflow-tooltip sortable fixed />
                <el-table-column prop="image" :label="t('container.imageName')" align="center" min-width="180" show-overflow-tooltip sortable />
                <el-table-column prop="ip" :label="t('container.containerIP')" align="center" min-width="120" show-overflow-tooltip />
                <el-table-column prop="ports" :label="t('container.containerPort')" align="center" min-width="140" show-overflow-tooltip />
                <el-table-column prop="state" :label="t('container.state')" align="center" min-width="120" show-overflow-tooltip>
                    <template #default="scope">
                        <el-tag v-if="scope.row.state === 'running'" type="success">
                            {{ t('container.enable') }}
                        </el-tag>
                        <el-tag v-else type="danger">
                            {{ t('container.disable') }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="uptime" :label="t('container.uptime')" align="center" min-width="160" show-overflow-tooltip />
                <!--                <el-table-column prop="cpu_percent" :label="t('container.cpuPercent')" align="center" min-width="140" show-overflow-tooltip /> -->
                <!--                <el-table-column prop="memory_percent" :label="t('container.menPercent')" align="center" min-width="180" show-overflow-tooltip /> -->
                <!--                <el-table-column prop="memory_usage" :label="t('container.memUsed')" align="center" min-width="140" show-overflow-tooltip /> -->
                <!--                <el-table-column prop="memory_limit" :label="t('container.memLimited')" align="center" min-width="160" show-overflow-tooltip /> -->
                <el-table-column :label="t('container.operator')" width="190" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="primary" size="small" text :disabled="!selectedAgentID || isProtectedContainerName(scope.row.name)" @click="startContainer({ title: 'container.startContainer', id: scope.row.id, update: search })">
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
                                        <el-button type="primary" size="small" text :disabled="!selectedAgentID || isProtectedContainerName(scope.row.name)" @click="editContainer({ title: 'container.editContainer', container: scope.row, update: search })">
                                            <svg-icon icon-class="edit" />
                                            {{ t('container.edit') }}
                                        </el-button>
                                    </el-dropdown-item>
                                    <el-dropdown-item>
                                        <el-button type="warning" size="small" text :disabled="!selectedAgentID || isProtectedContainerName(scope.row.name)" @click="stopContainer({ title: 'container.stopContainer', id: scope.row.id, update: search })">
                                            <svg-icon icon-class="stop" />
                                            {{ t('container.stop') }}
                                        </el-button>
                                    </el-dropdown-item>
                                    <el-dropdown-item>
                                        <el-button type="warning" size="small" text :disabled="!selectedAgentID || isProtectedContainerName(scope.row.name)" @click="restartContainer({ title: 'container.restartContainer', id: scope.row.id, update: search })">
                                            <svg-icon icon-class="update" />
                                            {{ t('container.restart') }}
                                        </el-button>
                                    </el-dropdown-item>
                                    <el-dropdown-item>
                                        <el-button type="danger" size="small" text :disabled="!selectedAgentID || isProtectedContainerName(scope.row.name)" @click="deleteContainer({ title: 'container.deleteContainer', id: scope.row.id, update: search })">
                                            <svg-icon icon-class="delete" />
                                            {{ t('container.delete') }}
                                        </el-button>
                                    </el-dropdown-item>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <div v-if="!isAgentEmpty" class="am-pagination">
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
