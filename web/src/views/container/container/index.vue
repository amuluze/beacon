<template>
    <div class="am-container-title">
        <span @click="$router.push('/docker/container')">容器</span>
        <span @click="$router.push('/docker/image')">镜像</span>
        <span @click="$router.push('/docker/network')">网络</span>
        <span @click="$router.push('/docker/settings')">设置</span>
    </div>
    <div class="am-container-operator">
        <el-card shadow="never">
            <el-button type="primary" @click="drawer = true">创建容器</el-button>
            <!-- <el-button type="warning">清理容器</el-button> -->
        </el-card>
    </div>
    <!-- 表格主体 -->
    <el-card shadow="never">
        <!-- https://blog.csdn.net/qq_24950043/article/details/114292940 -->
        <div class="am-table">
            <el-table :data="data" :key="containerKey" stripe ref="multipleTable" v-loading="loading">
                <el-table-column prop="id" label="容器 ID" min-width="100" align="center" fixed sortable />
                <el-table-column
                    prop="name"
                    label="容器名称"
                    align="center"
                    min-width="120"
                    show-overflow-tooltip
                    fixed
                    sortable
                />
                <el-table-column
                    prop="image"
                    label="镜像名称"
                    align="center"
                    min-width="120"
                    show-overflow-tooltip
                    sortable
                />
                <el-table-column
                    prop="ip"
                    label="容器 IP"
                    align="center"
                    min-width="100"
                    show-overflow-tooltip
                    sortable
                />
                <el-table-column prop="state" label="运行状态" align="center" min-width="120" show-overflow-tooltip>
                    <template #default="scope">
                        <el-tag :type="scope.row.state === 'running' ? 'success' : 'danger'">{{
                            scope.row.state
                        }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column
                    prop="uptime"
                    label="启动时间"
                    align="center"
                    min-width="120"
                    show-overflow-tooltip
                    sortable
                />
                <el-table-column
                    prop="cpu_percent"
                    label="CPU使用率"
                    align="center"
                    min-width="130"
                    show-overflow-tooltip
                    sortable
                />
                <el-table-column
                    prop="memory_percent"
                    label="内存使用率"
                    align="center"
                    min-width="130"
                    show-overflow-tooltip
                    sortable
                />
                <el-table-column
                    prop="memory_usage"
                    label="内存使用量"
                    align="center"
                    min-width="130"
                    show-overflow-tooltip
                    sortable
                />
                <el-table-column
                    prop="memory_limit"
                    label="内存限制"
                    align="center"
                    min-width="120"
                    show-overflow-tooltip
                />
                <el-table-column label="操作" width="200" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="primary" size="small" text @click="viewLog(scope.row.id)"> 日志 </el-button>
                        <el-button type="primary" size="small" text @click="startContainerByID(scope.row.id)">
                            启动
                        </el-button>
                        <el-dropdown>
                            <el-button type="primary" size="small" text>更多</el-button>
                            <template #dropdown>
                                <el-dropdown-menu>
                                    <el-dropdown-item>
                                        <el-button
                                            type="warning"
                                            size="small"
                                            text
                                            @click="stopContainerByID(scope.row.id)"
                                        >
                                            停止
                                        </el-button>
                                    </el-dropdown-item>
                                    <el-dropdown-item>
                                        <el-button
                                            type="warning"
                                            size="small"
                                            text
                                            @click="restartContainerByID(scope.row.id)"
                                        >
                                            重启
                                        </el-button>
                                    </el-dropdown-item>
                                    <el-dropdown-item>
                                        <el-button
                                            type="danger"
                                            size="small"
                                            text
                                            @click="deleteContainerByID(scope.row.id)"
                                        >
                                            删除
                                        </el-button>
                                    </el-dropdown-item>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
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
                    @size-change="(size: number) => pagination.onSizeChange(size, params)"
                    @current-change="(page: number) => pagination.onPageChange(page, params)"
                />
            </el-config-provider>
        </div>
    </el-card>

    <!--  查看日志弹窗  -->
    <el-dialog v-model="dialogVisible" title="日志" width="50%" :before-close="handleClose">
        <el-input v-model="logData" rows="20" type="textarea" />
        <template #footer>
            <div class="dialog-footer">
                <el-button size="small" type="primary" plain @click="downloadLog">下载</el-button>
                <el-button size="small" type="info" plain @click="stopLogView">停止</el-button>
                <el-button size="small" type="success" plain @click="handleClose">关闭</el-button>
            </div>
        </template>
    </el-dialog>
    <!-- 创建容器 -->
    <div class="am-container-create">
        <el-drawer v-model="drawer" size="540" title="创建容器">
            <div class="am-container-create__operator">
                <el-button type="default" size="default" plain>取消</el-button>
                <el-button type="primary" size="default" plain>确定</el-button>
            </div>
        </el-drawer>
    </div>
</template>

<script setup lang="ts">
import { queryContainers, removeContainer, restartContainer, startContainer, stopContainer } from '@/api/container'
import { Websocket } from '@/components/Websocket'
import { useTable } from '@/hooks/useTable'
import {
    RemoveContainerArgs,
    RestartContainerArgs,
    StartContainerArgs,
    StopContainerArgs
} from '@/interface/container.ts'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

onMounted(() => {
    refresh()
})

const params = {}

console.log('.....', params)
const { data, refresh, loading, pagination } = useTable(queryContainers, {}, {})

const dialogVisible = ref(false)
const containerKey = ref(0)
const logData = ref('')

let ws: Websocket

const viewLog = (container_id: string) => {
    logData.value = ''
    dialogVisible.value = true
    console.log('container_id', container_id)
    console.log('host', location.host)
    console.log('port', location.port)

    const onOpen = (ws: Websocket, ev: Event) => {
        console.log(ev)
        ws.send(container_id)
    }

    const onMessage = (ws: Websocket, ev: MessageEvent) => {
        console.log(ws)
        logData.value = logData.value + '\n' + ev.data
    }

    ws = new Websocket('ws/' + container_id, onOpen, onMessage)
}

const startContainerByID = (container_id: string) => {
    console.log('start container', container_id)
    const params: StartContainerArgs = {
        container_id: container_id
    }
    startContainer(params).finally(() => {
        refresh()
        containerKey.value += 1
    })
}

const stopContainerByID = (container_id: string) => {
    console.log('start container', container_id)
    const params: StopContainerArgs = {
        container_id: container_id
    }
    stopContainer(params).finally(() => {
        refresh()
        containerKey.value += 1
    })
}

const restartContainerByID = (container_id: string) => {
    console.log('start container', container_id)
    const params: RestartContainerArgs = {
        container_id: container_id
    }
    restartContainer(params).finally(() => {
        refresh()
        containerKey.value += 1
    })
}

const deleteContainerByID = (container_id: string) => {
    ElMessageBox.confirm('该操作会强制删除该容器. 继续吗?', '警告', {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel',
        type: 'warning'
    })
        .then(() => {
            const params: RemoveContainerArgs = {
                container_id: container_id
            }
            removeContainer(params)
                .then(() => {
                    ElMessage({
                        type: 'success',
                        message: '删除完成'
                    })
                })
                .finally(() => {
                    refresh()
                    containerKey.value += 1
                })
        })
        .catch(() => {
            ElMessage({
                type: 'info',
                message: '删除取消'
            })
        })
}

const handleClose = () => {
    dialogVisible.value = false
    ws.close()
}

const stopLogView = () => {
    ws.close()
}

const downloadLog = () => {
    console.log(`${logData.value}`)
    const a = document.createElement('a')
    a.setAttribute('href', `data:text/plain;charset=utf-8,${encodeURIComponent(logData.value)}`)
    a.setAttribute('download', 'log.txt')
    a.style.display = 'none'
    a.click()
}

const drawer = ref(false)
</script>

<style scoped lang="scss">
@include b(container-title) {
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
@include b(container-operator) {
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
