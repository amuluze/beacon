<template>
    <div class="am-host-header">
        <span @click="$router.push('/host/monitor')">监控</span>
        <span @click="$router.push('/host/file')">文件</span>
        <span @click="$router.push('/host/terminal')">终端</span>
        <span @click="$router.push('/host/settings')">设置</span>
    </div>
    <div class="am-host-operator">
        <el-card shadow="never">
            <el-dropdown>
                <el-button type="info"> 新建 <svg-icon icon-class="down" /> </el-button>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item>文件夹</el-dropdown-item>
                        <el-dropdown-item>文件</el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
            <el-button-group class="ml-4">
                <el-button type="primary"> 上传 </el-button>
                <el-button type="primary"> 下载 </el-button>
                <el-button type="primary"> 删除 </el-button>
            </el-button-group>
        </el-card>
    </div>
    <el-card shadow="never">
        <!-- https://blog.csdn.net/qq_24950043/article/details/114292940 -->
        <div class="am-table">
            <el-table :data="data" :key="containerKey" stripe ref="multipleTable" v-loading="loading">
                <el-table-column type="selection" width="55" />
                <el-table-column prop="name" label="名称" min-width="100" align="center" fixed sortable />
                <el-table-column prop="size" label="大小" align="center" min-width="100" sortable />
                <el-table-column prop="mode" label="权限" align="center" min-width="100" />
                <el-table-column prop="mod_time" label="修改时间" align="center" min-width="160" sortable />
                <el-table-column label="操作" width="200" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="primary" size="small" @click="deleteFile(scope.row.id)"> 删除 </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <!-- <div class="am-pagination">
            <el-config-provider :locale="zhCn">
                <el-pagination
                    v-model:current-page="pagination.page"
                    :page-size="pagination.size"
                    layout="total, sizes, prev, pager, next, jumper"
                    :page-sizes="pagination.sizeOption"
                    :total="pagination.total"
                    @size-change="(size: number) => pagination.onSizeChange(size, params)"
                    @current-change="(page: number) => pagination.onPageChange(page, params)"
                />
            </el-config-provider>
        </div> -->
    </el-card>
</template>
<script setup lang="ts">
import { queryFiles } from '@/api/host'
import { useTable } from '@/hooks/useTable'

onMounted(() => {
    refresh()
})

const params = {}

console.log('.....', params)
const { data, refresh, loading } = useTable(queryFiles, {}, {})

const containerKey = ref(0)

const deleteFile = (id: string) => {
    console.log(id)
}
</script>
<style scoped lang="scss">
@include b(host-header) {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: flex-start;
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
        &:nth-child(2) {
            color: #2f7bff;
            &::before {
                content: '';
                position: absolute;
                left: 84px;
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

@include b(host-operator) {
    height: 48px;
    width: 100%;
    margin-bottom: 4px;
    .el-card {
        height: 100%;
        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 0 0 16px;
            display: flex;
            align-items: center;
            justify-content: flex-start;
        }
    }
    .el-dropdown {
        margin-right: 16px;
    }
}

@include b(pagination) {
    margin-top: 10px;
    display: flex;
    justify-content: flex-end;
}

@include b(table) {
    width: 100%;
    height: calc(100vh - 230px);
    overflow-y: auto;
}
</style>
