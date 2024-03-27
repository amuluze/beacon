<template>
    <!-- 表格主体 -->
    <el-card shadow="never">
        <el-table :data="data" border stripe ref="multipleTable" v-loading="loading">
            <el-table-column prop="id" label="容器 ID" align="center" fixed />
            <el-table-column prop="name" label="容器名称" align="center" min-width="100" show-overflow-tooltip fixed />
            <el-table-column prop="image" label="镜像名称" align="center" min-width="100" show-overflow-tooltip />
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
                    <el-button type="primary" size="small" @click="viewLog(scope.row.id)">
                        <i-ep-view />查看日志
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
                    @size-change="(size) => pagination.onSizeChange(size, params)"
                    @current-change="(page) => pagination.onPageChange(page, params)"
                />
            </el-config-provider>
        </div>
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
import { queryContainers } from '@/api/container'
import { useTable } from '@/hooks/useTable'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

onMounted(() => {
    refresh()
})

const params = {}

console.log('.....', params)
const { data, refresh, loading, pagination } = useTable(queryContainers, {}, {})

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

<style scoped lang="scss">
@include b(pagination) {
    margin-top: 10px;
    display: flex;
    justify-content: flex-end;
}
</style>
