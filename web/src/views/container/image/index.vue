<template>
    <div class="am-image-title">
        <span @click="$router.push('/docker/container')">容器</span>
        <span @click="$router.push('/docker/image')">镜像</span>
        <span @click="$router.push('/docker/network')">网络</span>
        <span @click="$router.push('/docker/settings')">设置</span>
    </div>
    <div class="am-image-operator">
        <el-card shadow="never">
            <el-button type="primary" plain @click="imagePullDialog = true">拉取镜像</el-button>
            <el-button type="primary" plain @click="imageImportDialog = true">导入镜像</el-button>
            <el-button type="warning" plain @click="pruneImagesForce">清理虚悬镜像</el-button>
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
                        <!-- <el-button type="danger" plain size="small"> 导出 </el-button> -->
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
    <!-- 拉取镜像 -->
    <div class="am-image-pull">
        <el-dialog v-model="imagePullDialog" title="拉取镜像" width="50%">
            <el-input v-model="imageNameForPull" placeholder="请输入镜像名称" />
            <el-button size="default" type="info" plain @click="confirmImagePull" v-loading="imagePullLoading">
                确定
            </el-button>
        </el-dialog>
    </div>

    <!-- 导入镜像 -->
    <div class="am-image-import">
        <el-dialog v-model="imageImportDialog" title="导入镜像" width="50%">
            <el-upload drag action="/app/api/v1/container/image_import" multiple>
                <svg-icon icon-class="upload" />
                <div class="el-upload__text">Drop file here or <em>click to upload</em></div>
                <template #tip>
                    <div class="el-upload__tip">tar/tar.gz files</div>
                </template>
            </el-upload>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import { pruneImages, pullImage, queryImages, removeImage } from '@/api/container'
import { useTable } from '@/hooks/useTable'
import { PullImageArgs, RemoveImageArgs } from '@/interface/container.ts'
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
        type: 'warning',
        beforeClose: (action, instance, done) => {
            if (action === 'confirm') {
                instance.confirmButtonLoading = true
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
                        instance.confirmButtonLoading = false
                        done()
                        refresh()
                        imageKey.value += 1
                    })
            } else {
                done()
            }
        }
    })
        .then(() => {})
        .catch(() => {
            ElMessage({
                type: 'info',
                message: '删除取消'
            })
        })
    refresh()
}

const pruneImagesForce = () => {
    pruneImages().then(() => {
        ElMessage({
            type: 'success',
            message: '清理完成'
        })
    })
}

const imagePullLoading = ref(false)
const imagePullDialog = ref(false)
const imageNameForPull = ref('')
const confirmImagePull = async () => {
    imagePullLoading.value = true
    let params: PullImageArgs = {
        image_name: imageNameForPull.value
    }
    const { data } = await pullImage(params)
    console.log(data)
    imagePullLoading.value = false
    imagePullDialog.value = false
    refresh()
}

const imageImportDialog = ref(false)
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

@include b(image-pull) {
    :deep(.el-dialog) {
        .el-dialog__body {
            display: flex;
            flex-direction: row;
            align-items: center;
            justify-content: space-between;
            padding: 10px 20px;
            .el-input {
                width: 90%;
            }
        }
    }
}

@include b(image-import) {
    :deep(.el-upload) {
        .el-upload-dragger {
            width: 100%;
            height: 100%;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            .svg-icon {
                font-size: 67px;
                // color: #c0c4cc;
            }
            .el-upload__text {
                margin-top: 16px;
                color: #606266;
            }
        }
    }
}
</style>
