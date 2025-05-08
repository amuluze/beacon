<script setup lang="ts">
import zhCn from 'element-plus/es/locale/lang/zh-cn'

import type { Image } from '@/interface/container.ts'
import type { TableInstance } from 'element-plus'

import useCommandComponent from '@/hooks/useCommandComponent.ts'
import DeleteImage from '@/views/container/image/components/DeleteImage.vue'
import ImportImage from '@/views/container/image/components/ImportImage.vue'
import PruneImage from '@/views/container/image/components/PruneImage.vue'
import PullImage from '@/views/container/image/components/PullImage.vue'

import { queryImages } from '@/api/container'
import { useTable } from '@/hooks/useTable'
import useStore from '@/store'
import { getBrowserLanguage } from '@/utils'
import en from 'element-plus/es/locale/lang/en'
import { useI18n } from 'vue-i18n'

const { tableData, pageable, loading, search, handleSizeChange, handleCurrentChange } = useTable(queryImages)
onMounted(async () => {
  await search()
})

const tableRef = ref<TableInstance>()
const tableSelection = ref<Image[]>([])
const selectable = (row: Image) => !['1', '2'].includes(row.id)
function handleSelectionChange(val: Image[]) {
  console.log('selection: ', val)
  tableSelection.value = val
}

const pullImage = useCommandComponent(PullImage)
const importImage = useCommandComponent(ImportImage)
const pruneImage = useCommandComponent(PruneImage)
const deleteImage = useCommandComponent(DeleteImage)

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
            <div class="am-table-operator__left">
                <el-button type="primary" plain size="small" @click="pullImage({ title: 'image.pullImage', imageName: '', update: search })">
                    <svg-icon icon-class="download" />
                    {{ t('image.pullImage') }}
                </el-button>
                <el-button type="primary" plain size="small" @click="importImage({ title: 'image.importImage', update: importImage })">
                    <svg-icon icon-class="upload" />
                    {{ t('image.importImage') }}
                </el-button>
                <el-button type="warning" plain size="small" @click="pruneImage({ title: 'image.pruneImage', update: pruneImage })">
                    <svg-icon icon-class="delete" />
                    {{ t('image.pruneImage') }}
                </el-button>
            </div>
        </div>
        <div class="am-table">
            <el-table
                ref="tableRef"
                v-loading="loading"
                :data="tableData as Image[]"
                :header-cell-style="{ height: '45px', fontSize: '14px', color: '#000' }"
                height="100%"
                border
                @selection-change="handleSelectionChange"
            >
                <el-table-column type="selection" :selection="selectable" width="55" />
                <el-table-column prop="name" :label="t('image.imageName')" align="center" min-width="200" show-overflow-tooltip fixed sortable />
                <el-table-column prop="number" :label="t('image.containerNum')" align="center" show-overflow-tooltip min-width="120" sortable />
                <el-table-column prop="tag" :label="t('image.imageTag')" align="center" show-overflow-tooltip min-width="120" />
                <el-table-column prop="created" :label="t('image.createTime')" align="center" min-width="180" />
                <el-table-column prop="size" :label="t('image.imageSize')" align="center" min-width="120" />
                <el-table-column :label="t('image.operator')" width="160" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="danger" plain size="small" @click="deleteImage({ title: 'image.deleteImage', id: scope.row.id, update: search })">
                            <svg-icon icon-class="delete" />
                            {{ t('image.delete') }}
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
