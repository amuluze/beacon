<script setup lang="ts">
import { queryContainersUsage } from '@/api/container'
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import { containerCpuOption, containerMemOption } from '@/config/echarts.ts'
import type { CPUTrendingArgs, Usage } from '@/interface/host.ts'
import { dayjs } from 'element-plus'
import { set } from 'lodash-es'
import { useI18n } from 'vue-i18n'

const loading = ref(false)
const { t } = useI18n()

// 时间密度下拉框
const timeDensity = ref(600)
const options = [
    {
        value: 600,
        label: '10分钟',
    },
    {
        value: 1800,
        label: '30分钟',
    },
    {
        value: 3600,
        label: '1 小时',
    },
    {
        value: 43200,
        label: '12小时',
    },
    {
        value: 86400,
        label: '24小时',
    },
]
watch(
    () => timeDensity.value,
    () => {
        render()
    },
)

// CPU 使用率
const cpuOption = reactive<EChartsOption>(containerCpuOption)
const memOption = reactive<EChartsOption>(containerMemOption)
interface Series {
    name: string
    type: string
    smooth: boolean
    data: string[]
}
async function render() {
    const param: CPUTrendingArgs = {
        start_time: dayjs().unix() - timeDensity.value,
        end_time: dayjs().unix(),
    }
    const { data } = await queryContainersUsage(param)
    const cpuData = new Map<string, Usage[]>(Object.entries(data.cpu_usage))
    const memData = new Map<string, Usage[]>(Object.entries(data.mem_usage))
    set(cpuOption, 'legend.data', data.names)
    set(memOption, 'legend.data', data.names)

    const couIterator = cpuData.keys()
    const cpuDataFirstKey = couIterator.next().value
    const memIterator = memData.keys()
    const memDataFirstKey = memIterator.next().value
    set(
        cpuOption,
        'xAxis.data',
        cpuData.get(cpuDataFirstKey as string)?.map(item => `${dayjs(item.timestamp * 1000).hour()}:${dayjs(item.timestamp * 1000).minute()}`),
    )
    set(
        memOption,
        'xAxis.data',
        memData.get(memDataFirstKey as string)?.map(item => `${dayjs(item.timestamp * 1000).hour()}:${dayjs(item.timestamp * 1000).minute()}`),
    )
    const cpuSeries = reactive<Series[]>([])
    cpuData.forEach((value, key) => {
        cpuSeries.push({
            name: key,
            type: 'line',
            smooth: true,
            data: value.map(item => item.value.toFixed(2)),
        })
    })
    set(cpuOption, 'series', cpuSeries)
    const memSeries = reactive<Series[]>([])
    memData.forEach((value, key) => {
        memSeries.push({
            name: key,
            type: 'line',
            smooth: true,
            data: value.map(item => item.value.toFixed(2)),
        })
    })
    set(memOption, 'series', memSeries)
    console.log('cpuOption: ', cpuOption)
    console.log('memOption: ', memOption)
}

// WebSocket连接
let ws: WebSocket | null = null
const containerStats = ref<any[]>([])

// 初始化WebSocket连接
function initWebSocket() {
    loading.value = true

    // 关闭已存在的连接
    if (ws) {
        ws.close()
    }

    // 创建新的WebSocket连接
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws/monitor`

    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
        console.log('WebSocket连接已建立')
        loading.value = false
    }

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data)
        if (data.type === 'containerStats') {
            containerStats.value = data.stats
            updateCharts()
        }
        else if (data.error) {
            console.error('WebSocket错误:', data.error)
        }
    }

    ws.onerror = (error) => {
        console.error('WebSocket错误:', error)
        loading.value = false
    }

    ws.onclose = () => {
        console.log('WebSocket连接已关闭')
        loading.value = false
    }
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
    const currentTime = `${dayjs().hour()}:${dayjs().minute()}`

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
            existingSeries.data.push(stat.mem.toFixed(2))
            if (existingSeries.data.length > 60) {
                existingSeries.data.shift()
            }
        }
        else {
            (memSeries as Series[]).push({
                name: stat.name,
                type: 'line',
                smooth: true,
                data: [stat.mem.toFixed(2)],
            })
        }
    })
    set(memOption, 'series', memSeries)
}

// 定时器
const timer = ref()
onMounted(() => {
    console.log('mounted')
    // 初始化WebSocket连接
    initWebSocket()

    // 如果WebSocket连接失败，则回退到HTTP API
    timer.value = setInterval(() => {
        if (!ws || ws.readyState !== WebSocket.OPEN) {
            console.log('WebSocket未连接，使用HTTP API')
            render()
        }
    }, 5000)
})
onUnmounted(() => {
    clearInterval(timer.value)

    // 关闭WebSocket连接
    if (ws) {
        ws.close()
        ws = null
    }
})
</script>

<template>
    <div class="am-density">
        <el-card shadow="never">
            <span>{{ t('monitor.timeDensity') }}:</span>
            <el-select v-model="timeDensity" placeholder="Select" size="default" style="width: 120px">
                <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
        </el-card>
    </div>
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
