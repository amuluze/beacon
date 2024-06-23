<template>
    <div class="am-image-title">
        <span @click="$router.push('/docker/container')">容器</span>
        <span @click="$router.push('/docker/image')">镜像</span>
        <span @click="$router.push('/docker/network')">网络</span>
        <span @click="$router.push('/docker/settings')">配置</span>
    </div>
    <div class="am-image-operator">
        <el-card shadow="never">
            <el-button type="primary">拉取镜像</el-button>
            <el-button type="primary">导入镜像</el-button>
            <el-button type="warning" @click="pruneImagesForce">清理虚悬镜像</el-button>
        </el-card>
    </div>
    <el-card shadow="never">
        <div class="am-table">
            <el-table :data="data" :key="imageKey" highlight-current-row height="100%" stripe v-loading="loading">
                <el-table-column prop="id" label="镜像 ID" align="center" width="120" fixed sortable />
                <el-table-column
                    prop="name"
                    label="镜像名称"
                    align="center"
                    min-width="150"
                    show-overflow-tooltip
                    fixed
                    sortable
                />
                <el-table-column
                    prop="number"
                    label="容器数量"
                    align="center"
                    show-overflow-tooltip
                    width="120"
                    sortable
                />
                <el-table-column prop="tag" label="镜像 Tag" align="center" show-overflow-tooltip width="120" />
                <el-table-column prop="created" label="创建时间" align="center" width="200" sortable />
                <el-table-column prop="size" label="镜像大小" align="center" width="120" sortable />
                <el-table-column label="操作" width="160" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="danger" plain size="small"> 导出 </el-button>
                        <el-button type="danger" plain size="small" @click="deleteImageByID(scope.row.id)">
                            删除
                        </el-button>
                    </template>
                </el-table-column>
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
import { pruneImages, queryImages, removeImage } from '@/api/container'
import { useTable } from '@/hooks/useTable'
import { RemoveImageArgs } from '@/interface/container.ts'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

onMounted(() => {
    refresh()
})

const imageKey = ref(0)

const { data, refresh, loading, pagination } = useTable(queryImages, {}, {})

const deleteImageByID = (id: string) => {
    ElMessageBox.confirm('该操作会删除该镜像. 继续吗?', '警告', {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel',
        type: 'warning'
    })
        .then(() => {
            const params: RemoveImageArgs = {
                image_id: id
            }
            removeImage(params)
                .then(() => {
                    ElMessage({
                        type: 'success',
                        message: '删除完成'
                    })
                })
                .finally(() => {
                    refresh()
                    imageKey.value += 1
                })
        })
        .catch(() => {
            ElMessage({
                type: 'info',
                message: '删除取消'
            })
        })
}

const pruneImagesForce = () => {
    pruneImages().then(() => {
        ElMessage({
            type: 'success',
            message: '清理完成'
        })
    })
}
</script>

<style scoped lang="scss">
@include b(image-title) {
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
        align-items: center;
        font-size: 16px;
        font-weight: 600;
        color: #424244;
        margin-left: 16px;
        margin-right: 16px;
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
}
@include b(image-operator) {
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
}

.el-card {
    width: 100%;
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
