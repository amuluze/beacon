<template>
    <!-- 表格主体 -->
    <el-card shadow="never">
        <!-- https://blog.csdn.net/qq_24950043/article/details/114292940 -->
        <el-table :data="data" :key="containerKey" border stripe ref="multipleTable" v-loading="loading">
            <el-table-column prop="id" label="容器 ID" align="center" fixed />
            <el-table-column prop="name" label="容器名称" align="center" min-width="100" show-overflow-tooltip fixed />
            <el-table-column prop="image" label="镜像名称" align="center" min-width="100" show-overflow-tooltip />
            <el-table-column prop="ip" label="容器 IP" align="center" min-width="100" show-overflow-tooltip />
            <el-table-column prop="state" label="运行状态" align="center" min-width="90" show-overflow-tooltip>
                <template #default="scope">
                    <el-tag :type="scope.row.state === 'running' ? 'success' : 'danger'">{{ scope.row.state }}</el-tag>
                </template>
            </el-table-column>
            <el-table-column prop="uptime" label="启动时间" align="center" min-width="100" show-overflow-tooltip />
            <el-table-column
                prop="cpu_percent"
                label="CPU使用率"
                align="center"
                min-width="100"
                show-overflow-tooltip
            />
            <el-table-column
                prop="memory_percent"
                label="内存使用率"
                align="center"
                min-width="100"
                show-overflow-tooltip
            />
            <el-table-column
                prop="memory_usage"
                label="内存使用量"
                align="center"
                min-width="100"
                show-overflow-tooltip
            />
            <el-table-column
                prop="memory_limit"
                label="内存限制"
                align="center"
                min-width="100"
                show-overflow-tooltip
            />
            <el-table-column label="操作" width="320" fixed="right" align="center">
                <template #default="scope">
                    <el-button type="primary" size="small" plain @click="viewLog(scope.row.id)"> 日志 </el-button>
                    <el-button type="primary" size="small" plain @click="startContainerByID(scope.row.id)">
                        启动
                    </el-button>
                    <el-button type="warning" size="small" plain @click="stopContainerByID(scope.row.id)">
                        停止
                    </el-button>
                    <el-button type="warning" size="small" plain @click="restartContainerByID(scope.row.id)">
                        重启
                    </el-button>
                    <el-button type="danger" size="small" plain @click="deleteContainerByID(scope.row.id)">
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
</template>

<script setup lang="ts">
import { queryContainers, removeContainer, restartContainer, startContainer, stopContainer } from '@/api/container'
import { useTable } from '@/hooks/useTable'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import { Websocket } from '@/components/Websocket'
import {
    RemoveContainerArgs,
    RestartContainerArgs,
    StartContainerArgs,
    StopContainerArgs
} from '@/interface/container.ts'

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
    restartContainer(params)
}

const deleteContainerByID = (container_id: string) => {
    console.log('start container', container_id)
    const params: RemoveContainerArgs = {
        container_id: container_id
    }
    removeContainer(params).finally(() => {
        refresh()
        containerKey.value += 1
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
</script>

<style scoped lang="scss">
@include b(pagination) {
    margin-top: 10px;
    display: flex;
    justify-content: flex-end;
}
</style>
