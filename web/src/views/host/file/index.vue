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
                <el-button type="primary" plain> 新建 <svg-icon icon-class="down" /> </el-button>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item>
                            <el-button type="primary" size="small" text @click="createFolder">
                                <svg-icon icon-class="folder" style="color: #105eeb; margin-right: 4px" />
                                文件夹
                            </el-button>
                        </el-dropdown-item>
                        <el-dropdown-item>
                            <el-button type="primary" size="small" text @click="createFile">文件</el-button>
                        </el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
            <el-button-group class="ml-4" style="margin-right: 16px">
                <el-button type="primary" @click="goBack">
                    <svg-icon icon-class="back"></svg-icon>
                </el-button>
                <el-button type="primary" @click="goNext">
                    <svg-icon icon-class="next"></svg-icon>
                </el-button>
            </el-button-group>
            <el-button-group class="ml-4">
                <el-button type="primary" plain @click="uploadFile"> 上传 </el-button>
                <el-button type="primary" plain @click="downloadFile"> 下载 </el-button>
            </el-button-group>
        </el-card>
    </div>
    <el-card shadow="never">
        <!-- https://blog.csdn.net/qq_24950043/article/details/114292940 -->
        <div class="am-table">
            <el-table :data="filesData" :key="containerKey" stripe ref="multipleTable" v-loading="loading">
                <el-table-column type="selection" width="55" />
                <el-table-column prop="name" label="名称" min-width="150" align="center" fixed>
                    <template #default="scope">
                        <div style="display: flex; align-items: center">
                            <svg-icon
                                v-if="scope.row.is_dir == true"
                                icon-class="folder"
                                style="color: #105eeb; margin-right: 4px"
                            />
                            <span @click="queryFilesByPath(currentPath + '/' + scope.row.name, scope.row.is_dir)">{{
                                scope.row.name
                            }}</span>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column prop="size" label="大小" align="center" min-width="100" sortable>
                    <template #default="scope">
                        <span>{{ convertBytesToReadable(scope.row.size) }}</span>
                    </template>
                </el-table-column>
                <el-table-column prop="mode" label="权限" align="center" min-width="100" />
                <el-table-column prop="mod_time" label="修改时间" align="center" min-width="100">
                    <template #default="scope">
                        <span>{{ dayjs(scope.row.mod_time * 1000).format('YYYY-MM-DD HH:mm:ss') }}</span>
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="200" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="danger" size="small" text @click="deleteFile(scope.row.name)">
                            删除
                        </el-button>
                        <el-button type="primary" size="small" text @click="downloadFile"> 下载 </el-button>
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
import { error, warning } from '@/components/Message/message'
import { FileInfo, FilesSearchArgs } from '@/interface/host'
import { convertBytesToReadable } from '@/utils/convert'
import { dayjs } from 'element-plus'

const loading = ref(false)

const currentPath = ref('/')
const nextPath = ref('')
const filesData = ref<FileInfo[]>([])
onMounted(() => {
    queryFilesByPath(currentPath.value, true)
})

const goBack = () => {
    if (currentPath.value == '/') {
        return
    }
    nextPath.value = currentPath.value
    let path = currentPath.value.split('/')
    path.pop()
    queryFilesByPath(path.join('/'), true)
}

const goNext = () => {
    if (currentPath.value == '/') {
        return
    }
    queryFilesByPath(nextPath.value, true)
}

const queryFilesByPath = async (path: string, isDir: boolean) => {
    if (!isDir) {
        error('该路径不是文件夹')
        return
    }
    loading.value = true
    currentPath.value = path
    let params: FilesSearchArgs = {
        path: path
    }
    const { data } = await queryFiles(params)
    console.log(data.files)
    filesData.value = data.files
    // 以 is_dir 字段进行排序
    filesData.value.sort((a, b) => {
        if (a.is_dir && !b.is_dir) {
            return -1
        } else if (!a.is_dir && b.is_dir) {
            return 1
        } else {
            return 0
        }
    })
    loading.value = false
}

const containerKey = ref(0)

const deleteFile = (name: string) => {
    console.log(name)
    warning('该功能尚未实现')
}

const uploadFile = () => {
    warning('该功能尚未实现')
}

const downloadFile = () => {
    warning('该功能尚未实现')
}

const createFolder = () => {
    warning('该功能尚未实现')
}

const createFile = () => {
    warning('该功能尚未实现')
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
    height: calc(100vh - 188px);
    overflow-y: auto;
}
</style>
