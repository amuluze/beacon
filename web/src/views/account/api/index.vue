<template>
    <div class="am-account-header">
        <span @click="$router.push('/account/user')">用户</span>
        <span @click="$router.push('/account/role')">角色</span>
        <span @click="$router.push('/account/api')">接口</span>
    </div>

    <!--  表格主体  -->
    <el-card shadow="never">
        <div class="am-table">
            <el-table :data="data" height="100%" :key="resourceKey" stripe ref="multipleTable" v-loading="loading">
                <el-table-column prop="name" label="接口名称" min-width="100"></el-table-column>
                <el-table-column prop="path" label="URL" min-width="200">
                    <template #default="scope">
                        <el-tag>{{ scope.row.path }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="method" label="请求方式" min-width="100" sortable></el-table-column>
                <el-table-column prop="status" label="状态" wmin-idth="100" sortable>
                    <template #default="scope">
                        <el-tag v-if="scope.row.status === 1" type="success">正常</el-tag>
                        <el-tag v-else type="danger">禁用</el-tag>
                    </template>
                </el-table-column>
            </el-table>
        </div>
    </el-card>
</template>
<script setup lang="ts">
import { queryResource } from '@/api/account'
import { useTable } from '@/hooks/useTable'

onMounted(() => {
    refresh()
})
const resourceKey = ref(0)
const params = {}
const { data, refresh, loading, pagination } = useTable(queryResource, {}, {})
</script>

<style scoped lang="scss">
@include b(account-header) {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: flex-start;
    height: 48px;
    width: 100%;
    background-color: #fff;

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

        &:nth-child(3) {
            color: #2f7bff;
            &::before {
                content: '';
                position: absolute;
                left: 148px;
                width: 4px;
                height: 16px;
                text-align: center;
                background-color: #2f7bff;
                border-radius: 2px;
            }
        }
    }

    .el-card {
        height: 100%;
        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 8px 0 0;
            display: flex;
            align-items: center;
            justify-content: flex-end;
        }
    }
}

@include b(table) {
    width: 100%;
    height: calc(100vh - 136px);
    overflow-y: auto;
}
</style>
