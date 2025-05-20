<script setup lang="ts">
import type { EChartsOption } from '@/components/Echarts/echarts.ts';
import { containerCpuOption, containerMemOption } from '@/config/echarts.ts';
import { dayjs } from 'element-plus';
import { set } from 'lodash-es';
// import { useI18n } from 'vue-i18n'
import { Websocket } from '@/components/Websocket';

const loading = ref(false)
// const { t } = useI18n()

// CPU 使用率
const cpuOption = reactive<EChartsOption>(containerCpuOption)
const memOption = reactive<EChartsOption>(containerMemOption)
interface Series {
    name: string
    type: string
    smooth: boolean
    data: string[]
}
// WebSocket连接
let ws: Websocket | null = null
const containerStats = ref<any[]>([])

// 初始化WebSocket连接
function initWebSocket() {
    loading.value = true

    // 清空历史数据
    containerStats.value = []
    set(cpuOption, 'xAxis.data', [])
    set(memOption, 'xAxis.data', [])
    set(cpuOption, 'series', [])
    set(memOption, 'series', [])
    set(cpuOption, 'legend.data', [])
    set(memOption, 'legend.data', [])
    
    // 添加初始化数据，避免图表为空
    const initialTime = dayjs().format('HH:mm:ss')
    set(cpuOption, 'xAxis.data', [initialTime])
    set(memOption, 'xAxis.data', [initialTime])
    
    // 设置初始化的空系列，保持图表结构
    set(cpuOption, 'series', [{
        name: '等待数据...',
        type: 'line',
        smooth: true,
        data: [0],
    }])
    set(memOption, 'series', [{
        name: '等待数据...',
        type: 'line',
        smooth: true,
        data: [0],
    }])
    
    // 设置初始化的图例
    set(cpuOption, 'legend.data', ['等待数据...'])
    set(memOption, 'legend.data', ['等待数据...'])

    loading.value = false

    // 关闭已存在的连接
    if (ws) {
        ws.close()
    }

    // 创建新的WebSocket连接
    // const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `ws/monitor`

    const onopen = (ws: Websocket, ev: Event) => {
        console.log('WebSocket连接已建立', ws, ev)
    }

    const onmessage = (ws: Websocket, ev: MessageEvent) => {
        console.log('WebSocket接收到消息', ws, ev)
        const data = JSON.parse(ev.data)
        console.log('WebSocket接收到消息', data)
        if (data.type === 'containerStats') {
            containerStats.value = data.stats
            updateCharts()
        }
        else if (data.error) {
            console.error('WebSocket错误:', data.error)
        }
    }

    ws = new Websocket(wsUrl, onopen, onmessage)

    // ws.onerror = (error) => {
    //     console.error('WebSocket错误:', error)
    //     loading.value = false
    // }

    // ws.onclose = () => {
    //     console.log('WebSocket连接已关闭')
    //     loading.value = false
    // }
}

// 更新图表数据
function updateCharts() {
    if (!containerStats.value || containerStats.value.length === 0)
        return

    // 提取容器名称列表
    const names = containerStats.value.map(item => item.name)
    set(cpuOption, 'legend.data', names)
    set(memOption, 'legend.data', names)

    // 获取当前时间
    // const currentTime = `${dayjs().hour()}:${dayjs().minute()}`
    const currentTime = dayjs().format('HH:mm:ss')

    // 更新X轴数据（时间）
    const xAxisData = (cpuOption.xAxis as any)?.data || []
    if (Array.isArray(xAxisData)) {
        xAxisData.push(currentTime)
        // 保持数据点数量在合理范围内
        if (xAxisData.length > 60) {
            xAxisData.shift()
        }
        set(cpuOption, 'xAxis.data', xAxisData)
        set(memOption, 'xAxis.data', xAxisData)
    }
    else {
        set(cpuOption, 'xAxis.data', [currentTime])
        set(memOption, 'xAxis.data', [currentTime])
    }

    // 更新CPU数据系列
    const cpuSeries = cpuOption.series || []
    containerStats.value.forEach((stat) => {
        const existingSeries = (cpuSeries as Series[]).find(s => s.name === stat.name)
        if (existingSeries) {
            existingSeries.data.push(stat.cpu.toFixed(2))
            if (existingSeries.data.length > 60) {
                existingSeries.data.shift()
            }
        }
        else {
            (cpuSeries as Series[]).push({
                name: stat.name,
                type: 'line',
                smooth: true,
                data: [stat.cpu.toFixed(2)],
            })
        }
    })
    set(cpuOption, 'series', cpuSeries)

    // 更新内存数据系列
    const memSeries = memOption.series || []
    containerStats.value.forEach((stat) => {
        const existingSeries = (memSeries as Series[]).find(s => s.name === stat.name)
        if (existingSeries) {
            existingSeries.data.push(stat.memUsed)
            if (existingSeries.data.length > 60) {
                existingSeries.data.shift()
            }
        }
        else {
            (memSeries as Series[]).push({
                name: stat.name,
                type: 'line',
                smooth: true,
                data: [stat.memUsed],
            })
        }
    })
    set(memOption, 'series', memSeries)
}

onMounted(() => {
    console.log('mounted')
    // 初始化WebSocket连接
    initWebSocket()
})
onUnmounted(() => {
    // 关闭WebSocket连接
    if (ws) {
        ws.close()
        ws = null
    }
})
</script>

<template>
    <div class="am-column">
        <el-row>
            <el-col :lg="12" :md="12" :sm="12" :xs="24">
                <el-card shadow="never">
                    <el-skeleton :loading="loading" animated>
                        <div class="am-column-content">
                            <echarts :option="cpuOption" />
                        </div>
                    </el-skeleton>
                </el-card>
            </el-col>
            <el-col :lg="12" :md="12" :sm="12" :xs="24">
                <el-card shadow="never">
                    <el-skeleton :loading="loading" animated>
                        <div class="am-column-content">
                            <echarts :option="memOption" />
                        </div>
                    </el-skeleton>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<style scoped lang="scss">
@include b(density) {
    height: 48px;
    width: 100%;

    span {
        font-size: 14px;
    }

    .el-card {
        height: 100%;

        :deep(.el-card__body) {
            height: 100% !important;
            padding: 0 16px 0 16px;
            display: flex;
            flex-direction: row;
            align-items: center;
            justify-content: flex-end;
        }
    }

    border-radius: 4px;
    margin-bottom: 4px;
}

@include b(column) {
    height: 100%;
    width: 100%;
    font-size: 14px;
    overflow-y: auto;
}

@include b(column-content) {
    height: 320px;
    width: 100%;
}
</style>
