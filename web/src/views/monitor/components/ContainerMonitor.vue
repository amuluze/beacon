<script setup lang="ts">
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import { Websocket } from '@/components/Websocket'
import { containerCpuOption, containerMemOption } from '@/config/echarts.ts'
import { dayjs } from 'element-plus'
import { set } from 'lodash-es'

const loading = ref(false)

// CPU 使用率
const cpuOption = reactive<EChartsOption>(containerCpuOption)
const memOption = reactive<EChartsOption>(containerMemOption)
interface Series {
    name: string
    type: string
    smooth: boolean
    data: number[]
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
}

// 更新图表数据
function updateCharts() {
    if (!containerStats.value || containerStats.value.length === 0)
        return

    // 检查是否需要清除初始化数据
    const needClearInitialData = (cpuOption.legend as { data: string[] })?.data &&
        (cpuOption.legend as { data: string[] }).data.length === 1 &&
        (cpuOption.legend as { data: string[] }).data[0] === '等待数据...'

    if (needClearInitialData) {
        // 清除初始化的系列数据，保留当前时间点
        const currentTime = dayjs().format('HH:mm:ss')
        set(cpuOption, 'xAxis.data', [currentTime])
        set(memOption, 'xAxis.data', [currentTime])
        set(cpuOption, 'series', [])
        set(memOption, 'series', [])
    }

    // 提取容器名称列表
    const names = containerStats.value.map(item => item.name)
    set(cpuOption, 'legend.data', names)
    set(memOption, 'legend.data', names)

    // 获取当前时间
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
            existingSeries.data.push(stat.cpu)
            if (existingSeries.data.length > 60) {
                existingSeries.data.shift()
            }
        }
        else {
            (cpuSeries as Series[]).push({
                name: stat.name,
                type: 'line',
                smooth: true,
                data: [stat.cpu],
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
    initWebSocket()
})

onUnmounted(() => {
    if (ws) {
        ws.close()
    }
})
</script>

<template>
    <div class="am-container-container">
        <el-card shadow="hover" class="am-chart-card">
            <el-skeleton :loading="loading" animated>
                <div class="am-column-content">
                    <div class="am-chart-container">
                        <div class="am-chart-container__chart-item">
                            <echarts :option="cpuOption" />
                        </div>
                        <div class="am-chart-container__chart-item">
                            <echarts :option="memOption" />
                        </div>
                    </div>
                </div>
            </el-skeleton>
        </el-card>
    </div>
</template>

<style scoped lang="scss">
:deep(.el-skeleton) {
    width: 100%;
    height: 100%;
    padding: 20px;
}

:deep(.el-skeleton__item) {
    height: 300px;
    margin-bottom: 16px;
}

@include b(container-container) {
    height: 100%;
    width: 100%;
}

@include b(chart-card) {
    transition: all 0.3s;
    border-radius: 16px;

    &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }

    :deep(.el-card__body) {
        padding: 20px;
    }
}

@include b(column-content) {
    height: 440px;
    width: 100%;
    padding: 16px;
    border-radius: 4px;
}

@include b(chart-container) {
    display: flex;
    flex-direction: row;
    gap: 24px;
    height: 100%;

    @include e(chart-item) {
        flex: 1;
        min-height: 400px;
        position: relative;
        background: var(--el-bg-color-page);
        border-radius: 8px;
        padding: 16px;
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    }
}
</style>
