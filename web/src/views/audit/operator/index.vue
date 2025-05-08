<script setup lang="ts">
import { queryOperateAudit } from '@/api/audit'
import { useTable } from '@/hooks/useTable.ts'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import { useI18n } from 'vue-i18n'
import useStore from '@/store'
import en from 'element-plus/es/locale/lang/en'
import { getBrowserLanguage } from '@/utils'

const { tableData, pageable, loading, search, handleSizeChange, handleCurrentChange } = useTable(queryOperateAudit, {}, true)

onMounted(() => {
  search()
})

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
