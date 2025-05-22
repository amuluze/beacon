<script setup lang="ts">
import type { EChartsOption } from '@/components/Echarts/echarts.ts';
import { containerCpuOption, containerMemOption } from '@/config/echarts.ts';
import { dayjs } from 'element-plus';
import { set } from 'lodash-es';
// import { useI18n } from 'vue-i18n'
import { Websocket } from '@/components/Websocket';

import { useRoute, useRouter } from 'vue-router';
const router = useRouter()
const route = useRoute()

const activeTab = computed(() => {
  return route.path.startsWith('/monitor/host') ? 'host' : 'container'
})

const loading = ref(false)
// const { t } = useI18n()

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

    // 检查是否需要清除初始化数据
    const needClearInitialData = (cpuOption.legend as { data: string[] })?.data &&
        (cpuOption.legend as { data: string[] }).data.length === 1 &&
        (cpuOption.legend as { data: string[] }).data[0] === '等待数据...';

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
    console.log('cpuSeries', cpuSeries)
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
    <div class="am-header-container">
        <div class="am-button-group">
            <el-button-group class="monitor-switch">
                <el-button
                    type="primary"
                    class="am-monitor-btn"
                    :class="{ active: activeTab === 'host' }"
                    @click="router.push('/monitor/host')">
                    主机监控
                </el-button>
                <el-button
                    type="primary"
                    class="am-monitor-btn"
                    :class="{ active: activeTab === 'container' }"
                    @click="router.push('/monitor/container')">
                    容器监控
                </el-button>
            </el-button-group>
        </div>
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
@include b(header-container) {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  gap: 16px;

  .el-card {
    border-radius: 8px;
    background: var(--el-bg-color-page);

    :deep(.el-card__body) {
      padding: 8px 16px !important;
    }
  }
}

@include b(button-group) {
  background: var(--el-bg-color-page);
  border-radius: 8px;
  padding: 4px;

  .monitor-switch {
    display: flex;
    gap: 4px;

    .am-monitor-btn {
      border: 1px solid var(--el-border-color);
      background: transparent;
      color: var(--el-text-color-regular);
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

      &.active {
        background: var(--el-color-primary);
        border-color: var(--el-color-primary);
        color: var(--el-color-white);
        box-shadow: 0 2px 4px rgba(var(--el-color-primary-rgb), 0.2);
      }

      &:hover:not(.active) {
        border-color: var(--el-color-primary-light-5);
        color: var(--el-color-primary);
        transform: none;
      }

      &:first-child {
        border-radius: 6px 0 0 6px;
      }

      &:last-child {
        border-radius: 0 6px 6px 0;
      }
    }
  }
}

@include b(density-selector) {
  background: var(--el-bg-color-page);
  border-radius: 8px;
  padding: 8px 16px;

  .selector-wrapper {
    display: flex;
    align-items: center;
    gap: 8px;
    height: 100%;
  }
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
