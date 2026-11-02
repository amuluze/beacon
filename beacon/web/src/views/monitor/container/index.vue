<script setup lang="ts">
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import MonitorContextBar from '@/components/Agent/MonitorContextBar.vue'
import AgentSelect from '@/layout/navbar/AgentSelect.vue'
import { containerCpuOption, containerMemOption } from '@/config/echarts.ts'
import type { Usage } from '@/interface/host.ts'
import { dayjs } from 'element-plus'
import { queryContainersUsage } from '@/api/container'
import { cloneDeep, set } from 'lodash-es'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { useI18n } from 'vue-i18n'
import IconCpu from '~icons/lucide/cpu'
import IconMemoryStick from '~icons/lucide/memory-stick'
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
const lastUpdated = ref(0)
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

// §7.2 容器指标以「容器名 + 百分比」chip 列表呈现
const containerCpuMetrics = ref<{ name: string, percent: number }[]>([])
const containerMemMetrics = ref<{ name: string, percent: number }[]>([])

// §5.1 图表序列色取自 chart-1..chart-6 token，不写死色值
const CHART_SERIES_TOKENS = ['--am-chart-1', '--am-chart-2', '--am-chart-3', '--am-chart-4', '--am-chart-5', '--am-chart-6']
const containerPalette = ref<string[]>([])

function resolveChartPalette(): string[] {
  if (typeof window === 'undefined')
    return CHART_SERIES_TOKENS.map(() => '')
  const styles = getComputedStyle(document.documentElement)
  return CHART_SERIES_TOKENS.map(token => styles.getPropertyValue(token).trim())
}

let latestRenderID = 0

function clearCharts(): void {
  containerCpuMetrics.value = []
  containerMemMetrics.value = []
  containerPalette.value = []
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

  const cpuData = new Map<string, Usage[]>(Object.entries(data.cpu_usage))
  const memData = new Map<string, Usage[]>(Object.entries(data.mem_usage))
  const containerPaletteColors = resolveChartPalette()
  containerPalette.value = containerPaletteColors

  // 容器名 chip 上的百分比取各容器最新采样值
  function lastPercent(series: Usage[] | undefined): number {
    const last = series?.[series.length - 1]
    return last ? Number(last.value.toFixed(1)) : 0
  }
  containerCpuMetrics.value = data.names.map(name => ({ name, percent: lastPercent(cpuData.get(name)) }))
  containerMemMetrics.value = data.names.map(name => ({ name, percent: lastPercent(memData.get(name)) }))

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
        lineStyle: { width: 1.5, color: containerPaletteColors[i % containerPaletteColors.length] },
        itemStyle: { color: containerPaletteColors[i % containerPaletteColors.length] },
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
        lineStyle: { width: 1.5, color: containerPaletteColors[i % containerPaletteColors.length] },
        itemStyle: { color: containerPaletteColors[i % containerPaletteColors.length] },
      })
    }
  })
  set(memOption, 'series', memSeries)
  lastUpdated.value = Date.now()
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
                <AgentSelect width="180px" />
            </div>
            <MonitorContextBar v-model:density="timeDensity" :density-options="options" :last-updated="lastUpdated" />
        </div>

        <AgentEmptyState v-if="isAgentEmpty" @refresh="refreshAgents" />
        <div v-else class="am-chart-grid">
            <div class="am-chart-row">
                <div class="am-chart-card">
                    <div class="am-chart-card-header">
                        <span class="am-chart-card-icon am-chart-card-icon--primary">
                            <IconCpu />
                        </span>
                        <span class="am-chart-card-title">CPU 使用率</span>
                    </div>
                    <div class="am-chart-card-legend am-chart-card-legend--scrollable" role="list" aria-label="CPU 容器图例">
                        <div v-for="(item, i) in containerCpuMetrics" :key="item.name" class="am-legend-item" role="listitem" :title="item.name">
                            <span class="am-legend-dot" :style="{ background: containerPalette[i % containerPalette.length] }" />
                            <span class="am-legend-label">{{ item.name }}</span>
                            <span class="am-legend-value">{{ item.percent }}%</span>
                        </div>
                    </div>
                    <div class="am-chart-area">
                        <echarts :option="cpuOption" />
                    </div>
                </div>
                <div class="am-chart-card">
                    <div class="am-chart-card-header">
                        <span class="am-chart-card-icon am-chart-card-icon--warning">
                            <IconMemoryStick />
                        </span>
                        <span class="am-chart-card-title">内存使用率</span>
                    </div>
                    <div class="am-chart-card-legend am-chart-card-legend--scrollable" role="list" aria-label="内存容器图例">
                        <div v-for="(item, i) in containerMemMetrics" :key="item.name" class="am-legend-item" role="listitem" :title="item.name">
                            <span class="am-legend-dot" :style="{ background: containerPalette[i % containerPalette.length] }" />
                            <span class="am-legend-label">{{ item.name }}</span>
                            <span class="am-legend-value">{{ item.percent }}%</span>
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
  flex-wrap: wrap;
  align-content: flex-start;
  gap: var(--am-spacing-xs) var(--am-spacing-md);
  min-width: 0;
}

.am-chart-card-legend--scrollable {
  max-height: 72px;
  overflow-x: hidden;
  overflow-y: auto;
  padding-right: var(--am-spacing-xs);
  scrollbar-gutter: stable;
}

.am-legend-item {
  min-width: 0;
  max-width: 180px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.am-legend-dot {
  flex: 0 0 auto;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  display: inline-block;
}
.am-legend-value {
  color: var(--am-foreground-secondary);
  font-family: 'Geist Mono', 'SFMono-Regular', Consolas, monospace;
  font-size: var(--am-font-xs);
}
.am-legend-label {
  min-width: 0;
  overflow: hidden;
  color: var(--am-foreground-secondary);
  font-size: var(--am-font-xs);
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
