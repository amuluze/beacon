<template>
    <!-- 表格主体 -->
    <el-card shadow="never">
        <el-table :data="tableData" style="width: 100%" ref="multipleTable" v-loading="loading">
            <el-table-column prop="id" label="容器 ID" align="center" fixed />
            <el-table-column prop="name" label="容器名称" align="center" min-width="100" show-overflow-tooltip fixed />
            <el-table-column prop="image" label="镜像名称" align="center" min-width="150" show-overflow-tooltip />
            <el-table-column prop="ip" label="容器 IP" align="center" min-width="100" show-overflow-tooltip />
            <el-table-column prop="state" label="运行状态" align="center" min-width="90" show-overflow-tooltip>
                <template #default="scope">
                    <el-tag :type="scope.row.state === 'running' ? 'success' : 'danger'">{{ scope.row.state }}</el-tag>
                </template>
            </el-table-column>
            <el-table-column prop="uptime" label="启动时长" align="center" min-width="100" show-overflow-tooltip />
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
            <el-table-column label="操作" width="180" fixed="right" align="center">
                <template #default="scope">
                    <el-button type="primary" size="default" @click="viewLog(scope.row.id)">
                        <i-ep-view />查看日志
                    </el-button>
                </template>
            </el-table-column>
        </el-table>
    </el-card>

    <!--  查看日志弹窗  -->
    <el-dialog v-model="dialogVisible" title="日志" width="50%" :before-close="handleClose">
        <el-input v-model="logData" rows="20" type="textarea" />
        <template #footer>
            <div class="dialog-footer">
                <el-button type="primary" plain @click="downloadLog">下载</el-button>
                <el-button type="info" plain @click="stopLogView">停止</el-button>
                <el-button type="success" plain @click="handleClose">关闭</el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
// 表格
import { queryContainers } from '@/api/container'
import { Container } from '@/interface/container.ts'
import { convertBytesToReadable } from '@/utils/convert'

const loading = ref(true)
type tableDataType = {
    id: string
    name: string
    state: string
    image: string
    ip: string
    uptime: string
    cpu_percent: string
    memory_percent: string
    memory_usage: string
    memory_limit: string
}
const tableData = ref<tableDataType[]>([])
onMounted(() => {
    getData()
})
const getData = async () => {
    loading.value = true
    const { data } = await queryContainers()
    console.log('res', data.containers)
    data.containers.map((item: Container) => {
        const tableDataItem: tableDataType = {
            id: item.id,
            name: item.name,
            state: item.state,
            image: item.image,
            ip: item.ip,
            uptime: item.uptime,
            cpu_percent: item.cpu_percent ? item.cpu_percent.toFixed(2) + '%' : '0.00%',
            memory_percent: item.memory_percent ? item.memory_percent.toFixed(2) + '%' : '0.00%',
            memory_usage: convertBytesToReadable(item.memory_usage),
            memory_limit: convertBytesToReadable(item.memory_limit)
        }
        tableData.value.push(tableDataItem)
    })
    loading.value = false
}

const dialogVisible = ref(false)
const logData = ref('')

let ws: WebSocket

const viewLog = (container_id: string) => {
    dialogVisible.value = true
    console.log('container_id', container_id)

    // 建立 WebSocket 连接
    ws = new WebSocket('wss://' + location.host + '/ws/' + container_id)
    ws.onopen = () => {
        ws.send(container_id)
    }

    // 监听消息事件
    ws.addEventListener('message', (event) => {
        logData.value = logData.value + '\n' + event.data
    })

    // 监听错误事件
    ws.addEventListener('error', (error) => {
        console.error(error)
    })

    // 监听连接关闭事件
    ws.addEventListener('close', () => {
        console.log('WebSocket closed')
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

<style scoped lang="scss"></style>
