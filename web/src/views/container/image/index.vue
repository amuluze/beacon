<template>
    <div class="am-image-operator">
        <el-card>
            <el-button type="primary" @click="pruneImagesForce">清理虚悬镜像</el-button>
        </el-card>
    </div>
    <el-card shadow="never">
        <el-table
            :data="data"
            :key="imageKey"
            max-height="600"
            highlight-current-row
            border
            stripe
            v-loading="loading"
            style="width: 100%"
        >
            <el-table-column prop="id" label="镜像 ID" align="center" width="150" fixed sortable />
            <el-table-column prop="name" label="镜像名称" align="center" min-width="100" show-overflow-tooltip fixed sortable />
            <el-table-column prop="number" label="容器数量" align="center" show-overflow-tooltip width="100" sortable />
            <el-table-column prop="tag" label="镜像 Tag" align="center" show-overflow-tooltip width="100" />
            <el-table-column prop="created" label="创建时间" align="center" width="200" sortable />
            <el-table-column prop="size" label="镜像大小" align="center" width="100" sortable />
            <el-table-column label="操作" width="160" fixed="right" align="center">
                <template #default="scope">
                    <el-button type="danger" plain size="small" @click="deleteImageByID(scope.row.id)">
                        删除
                    </el-button>
                </template>
            </el-table-column>
        </el-table>
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
import {pruneImages, queryImages, removeImage} from '@/api/container'
import { useTable } from '@/hooks/useTable'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import {RemoveImageArgs} from '@/interface/container.ts'

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
@include b(image-operator) {
    height: 48px;
    width: 100%;
    margin-bottom: 4px;
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

.el-card {
    width: 100%;
}

@include b(pagination) {
    margin-top: 10px;
    display: flex;
    justify-content: flex-end;
}
</style>
