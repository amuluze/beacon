<script setup lang="ts">
import { queryResource } from '@/api/account'
import { useTable } from '@/hooks/useTable.ts'
import type { TableInstance } from 'element-plus'
import { useI18n } from 'vue-i18n'

const tableRef = ref<TableInstance>()
const { tableData, loading, search } = useTable(queryResource, {}, false)

onMounted(async () => {
  await search()
})
const { t } = useI18n()
</script>

<template>
    <div class="am-container">
        <div class="am-table">
            <el-table
                ref="tableRef"
                v-loading="loading"
                :data="tableData"
                height="100%"
                :header-cell-style="{ height: '45px', fontSize: '14px', color: '#000'}"
                border
            >
                <el-table-column prop="name" label="接口名称" min-width="100" />
                <el-table-column prop="path" label="URL" min-width="200">
                    <template #default="scope">
                        <el-tag>{{ scope.row.path }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="method" label="请求方式" min-width="100" sortable />
                <el-table-column prop="status" label="状态" wmin-idth="100" sortable>
                    <template #default="scope">
                        <el-tag v-if="scope.row.status === 1" type="success">
                            {{ t('account.enable') }}
                        </el-tag>
                        <el-tag v-else type="danger">
                            {{ t('account.disable') }}
                        </el-tag>
                    </template>
                </el-table-column>
            </el-table>
        </div>
    </div>
</template>

<style scoped lang="scss">
</style>
