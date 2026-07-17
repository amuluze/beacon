<script setup lang="ts">
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import { containerCpuOption, containerMemOption } from '@/config/echarts.ts'
import type { Usage } from '@/interface/host.ts'
import { dayjs } from 'element-plus'
import { queryContainersUsage } from '@/api/container'
import { cloneDeep, set } from 'lodash-es'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { useI18n } from 'vue-i18n'
import { applyMonitorTimeRange, toMonitorSeriesPoint } from '@/utils/monitorChart'
import type { MonitorTimeRange } from '@/utils/monitorChart'

// Agent switcher
const { selectedAgentID: currentAgent, isAgentEmpty, agentParams, ensureSelectedAgent, loadAgents } = useAgentSelection({ immediate: false })
const { t } = useI18n()
watch(currentAgent, () => {
  void render()
})

// Time density
const timeDensity = ref(600)
const options = [
  { value: 300, label: '5分钟' },
  { value: 600, label: '10分钟' },
  { value: 1800, label: '30分钟' },
  { value: 3600, label: '1小时' },
  { value: 43200, label: '12小时' },
  { value: 86400, label: '24小时' },
]
watch(timeDensity, () => {
  void render()
})

// Charts
const cpuOption = reactive<EChartsOption>(cloneDeep(containerCpuOption))
const memOption = reactive<EChartsOption>(cloneDeep(containerMemOption))

const containerNames = ref<string[]>([])
const containerPalette = ['#409EFF', '#569A2E', '#C28014', '#DB5050']
let latestRenderID = 0

function clearCharts(): void {
  containerNames.value = []
  set(cpuOption, 'legend.data', [])
  set(memOption, 'legend.data', [])
  set(cpuOption, 'series', [])
  set(memOption, 'series', [])
}

async function render() {
  const renderID = ++latestRenderID
  if (!currentAgent.value)
    return
  const endTime = dayjs().unix()
  const param: MonitorTimeRange = { start_time: endTime - timeDensity.value, end_time: endTime, ...agentParams.value }
  applyMonitorTimeRange(cpuOption, param)
  applyMonitorTimeRange(memOption, param)
  const { data } = await queryContainersUsage(param as any)
  if (renderID !== latestRenderID)
    return
  if (!data.names || data.names.length === 0) {
    clearCharts()
    return
  }

  containerNames.value = data.names

  const cpuData = new Map<string, Usage[]>(Object.entries(data.cpu_usage))
  const memData = new Map<string, Usage[]>(Object.entries(data.mem_usage))

  // Legend
  set(cpuOption, 'legend.data', data.names.map((n: string, i: number) => ({ name: n, textStyle: { color: containerPalette[i % containerPalette.length] } })))
  set(memOption, 'legend.data', data.names.map((n: string, i: number) => ({ name: n, textStyle: { color: containerPalette[i % containerPalette.length] } })))

  // Series
  const cpuSeries: any[] = []
  data.names.forEach((name: string, i: number) => {
    const values = cpuData.get(name)
    if (values) {
      cpuSeries.push({
        name,
        data: values.map(item => toMonitorSeriesPoint(item.timestamp, Number(item.value.toFixed(1)))),
        type: 'line',
        smooth: true,
        showSymbol: false,
        lineStyle: { width: 1.5, color: containerPalette[i % containerPalette.length] },
        itemStyle: { color: containerPalette[i % containerPalette.length] },
      })
    }
  })
  set(cpuOption, 'series', cpuSeries)

  const memSeries: any[] = []
  data.names.forEach((name: string, i: number) => {
    const values = memData.get(name)
    if (values) {
      memSeries.push({
        name,
        data: values.map(item => toMonitorSeriesPoint(item.timestamp, Number(item.value.toFixed(1)))),
        type: 'line',
        smooth: true,
        showSymbol: false,
        lineStyle: { width: 1.5, color: containerPalette[i % containerPalette.length] },
        itemStyle: { color: containerPalette[i % containerPalette.length] },
      })
    }
  })
  set(memOption, 'series', memSeries)
}

async function refreshAgents() {
  await loadAgents()
  await render()
}

const timer = ref()
onMounted(async () => {
  await ensureSelectedAgent()
  void render()
  timer.value = setInterval(() => void render(), 10000)
})
onUnmounted(() => {
  clearInterval(timer.value)
})
</script>

<template>
    <div class="am-section">
        <div class="am-section-header">
            <div class="am-section-title-group">
                <span class="am-section-title">{{ t('menu.containerMonitor') }}</span>
            </div>
            <div class="am-density-group">
                <span class="am-density-label">时间密度：</span>
                <el-select v-model="timeDensity" size="small" style="width: 110px">
                    <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </div>
        </div>

        <AgentEmptyState v-if="isAgentEmpty" @refresh="refreshAgents" />
        <div v-else class="am-chart-grid">
            <div class="am-chart-row">
                <div class="am-chart-card">
                    <div class="am-chart-card-header">
                        <span class="am-chart-card-title">CPU 使用率</span>
                        <div class="am-chart-card-legend">
                            <div v-for="(name, i) in containerNames" :key="name" class="am-legend-item">
                                <span class="am-legend-dot" :style="{ background: containerPalette[i % containerPalette.length] }" />
                                <span class="am-legend-label">{{ name }}</span>
                            </div>
                        </div>
                    </div>
                    <div class="am-chart-area">
                        <echarts :option="cpuOption" />
                    </div>
                </div>
                <div class="am-chart-card">
                    <div class="am-chart-card-header">
                        <span class="am-chart-card-title">内存使用率</span>
                        <div class="am-chart-card-legend">
                            <div v-for="(name, i) in containerNames" :key="name" class="am-legend-item">
                                <span class="am-legend-dot" :style="{ background: containerPalette[i % containerPalette.length] }" />
                                <span class="am-legend-label">{{ name }}</span>
                            </div>
                        </div>
                    </div>
                    <div class="am-chart-area">
                        <echarts :option="memOption" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">
.am-chart-card-legend {
  display: flex;
  gap: 14px;
  flex-wrap: wrap;
}
.am-legend-item {
  display: flex;
  align-items: center;
  gap: 4px;
}
.am-legend-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  display: inline-block;
}
.am-legend-label {
  color: var(--am-foreground-muted);
  font-size: var(--am-font-xs);
}
</style>
