<script setup lang="ts">
import { queryNetworks } from '@/api/container'
import { useTable } from '@/hooks/useTable'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

import type { Network } from '@/interface/container.ts'
import type { TableInstance } from 'element-plus'

import useCommandComponent from '@/hooks/useCommandComponent.ts'
import AddNetwork from '@/views/container/network/components/AddNetwork.vue'
import DeleteNetwork from '@/views/container/network/components/DeleteNetwork.vue'

import useStore from '@/store'
import { getBrowserLanguage } from '@/utils'
import en from 'element-plus/es/locale/lang/en'
import { useI18n } from 'vue-i18n'

const { tableData, pageable, loading, search, handleSizeChange, handleCurrentChange } = useTable(queryNetworks)
onMounted(async () => {
  await search()
})

const tableRef = ref<TableInstance>()
const tableSelection = ref<Network[]>([])
const selectable = (row: Network) => !['1', '2'].includes(row.id)
function handleSelectionChange(val: Network[]) {
  console.log('selection: ', val)
  tableSelection.value = val
}

const addNetwork = useCommandComponent(AddNetwork)
const deleteNetwork = useCommandComponent(DeleteNetwork)

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
            <el-button type="primary" plain size="small" @click="addNetwork({ title: 'network.newNetwork', update: search })">
                <svg-icon icon-class="add" />
                {{ t('network.newNetwork') }}
            </el-button>
        </div>
        <div class="am-table">
            <el-table
                ref="tableRef"
                v-loading="loading"
                :data="tableData as Network[]"
                :header-cell-style="{ height: '45px', fontSize: '14px', color: '#000' }"
                height="100%"
                border
                @selection-change="handleSelectionChange"
            >
                <el-table-column type="selection" :selection="selectable" width="55" />
                <el-table-column prop="name" :label="t('network.networkName')" align="center" min-width="140" fixed />
                <el-table-column prop="driver" :label="t('network.networkMode')" align="center" min-width="120" />
                <el-table-column prop="subnet" :label="t('network.subNetwork')" align="center" show-overflow-tooltip min-width="120" />
                <el-table-column prop="gateway" :label="t('network.networkGateway')" align="center" show-overflow-tooltip min-width="160" />
                <el-table-column prop="created" :label="t('network.createTime')" align="center" min-width="200" />
                <el-table-column :label="t('network.operator')" width="160" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="danger" plain size="small" @click="deleteNetwork({ title: 'network.deleteNetwork', id: scope.row.id, update: search })">
                            <svg-icon icon-class="delete" />
                            {{ t('network.delete') }}
                        </el-button>
                    </template>
                </el-table-column>
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
