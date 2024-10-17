<template>
    <div class="am-audit-title">
        <span @click="$router.push('/audit/operator')">操作日志</span>
        <span @click="$router.push('/audit/system')">系统日志</span>
    </div>
    <el-card shadow="never">
        <div class="am-table">
            <el-table :data="data" highlight-current-row stripe v-loading="loading" style="width: 100%">
                <el-table-column prop="id" label="ID" align="center" min-width="150" />
                <el-table-column prop="username" label="用户名" align="center" min-width="150" />
                <el-table-column prop="operate" label="操作" align="center" show-overflow-tooltip min-width="150" />
                <el-table-column prop="created" label="操作时间" align="center" min-width="200" />
            </el-table>
        </div>
        <div class="am-pagination">
            <el-config-provider :locale="zhCn">
                <el-pagination
                    v-model:current-page="pagination.page"
                    :page-size="pagination.size"
                    layout="total, sizes, prev, pager, next, jumper"
                    :page-sizes="pagination.sizeOption"
                    :total="pagination.total"
                    @size-change="pagination.onSizeChange"
                    @current-change="pagination.onPageChange"
                />
            </el-config-provider>
        </div>
    </el-card>
</template>

<script setup lang="ts">
import { queryAudit } from '@/api/audit';
import { useTable } from '@/hooks/useTable.ts';
import zhCn from 'element-plus/es/locale/lang/zh-cn';

onMounted(() => {
    refresh()
})
const { data, refresh, loading, pagination } = useTable(queryAudit, {}, {})
</script>

<style scoped lang="scss">
@include b(audit-title) {
    display: flex;
    align-items: center;
    justify-content: space-start;
    height: 48px;
    width: 100%;
    background-color: #ffffff;
    // box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);

    border-radius: 4px;
    margin-bottom: 8px;
    padding: 0 16px;

    span {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-start;
        font-size: 16px;
        font-weight: 600;
        margin-left: 16px;
        margin-right: 16px;
        color: #424244;
        cursor: pointer;

        &:first-child {
            color: #2f7bff;
            &::before {
                content: '';
                position: absolute;
                left: 20px;
                width: 4px;
                height: 16px;
                text-align: center;
                background-color: #2f7bff;
                border-radius: 2px;
            }
        }
    }
}
@include b(pagination) {
    margin-top: 10px;
    display: flex;
    justify-content: flex-end;
}
@include b(table) {
    width: 100%;
    height: calc(100vh - 174px);
    overflow-y: auto;
}
</style>
