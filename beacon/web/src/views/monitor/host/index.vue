<script setup lang="ts">
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import type { DiskIO, DiskUsage, Freshness, NetIO, NetUsage } from '@/interface/host.ts'

import { queryContainersUsage } from '@/api/container'
import { queryCPUInfo, queryCPUUsage, queryDiskInfo, queryDiskUsage, queryMemInfo, queryMemUsage, queryNetworkUsage } from '@/api/host'
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import DataStaleTag from '@/components/Agent/DataStaleTag.vue'
import MonitorContextBar from '@/components/Agent/MonitorContextBar.vue'
import MonitorStalePanel from '@/components/Agent/MonitorStalePanel.vue'
import type { StaleDetailRow, StaleKpi } from '@/components/Agent/MonitorStalePanel.vue'
import AgentSelect from '@/layout/navbar/AgentSelect.vue'
import { cpuTrendingOption, diskTrendingOption, formatBytesPerSecond, memTrendingOption, netTrendingOption } from '@/config/echarts.ts'
import { convertBytesToReadable } from '@/utils/convert.ts'
import { dayjs } from 'element-plus'
import { cloneDeep, set } from 'lodash-es'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { useI18n } from 'vue-i18n'
import IconCpu from '~icons/lucide/cpu'
import IconHardDrive from '~icons/lucide/hard-drive'
import IconMemoryStick from '~icons/lucide/memory-stick'
import IconNetwork from '~icons/lucide/network'
import { applyMonitorTimeRange, toMonitorSeriesPoint } from '@/utils/monitorChart'
import type { MonitorTimeRange } from '@/utils/monitorChart'
import { formatSampleClock } from '@/utils/monitorStale'
import { resolveChartLegendTextColor, resolveChartPalette, resolveChartTokenValue } from '@/utils/chartTokens'

// Agent switcher
const { selectedAgentID: currentAgent, isAgentEmpty, agentParams, ensureSelectedAgent, loadAgents } = useAgentSelection({ immediate: false })
watch(currentAgent, () => {
  void refreshAll()
})

// Time density
const timeDensity = ref(600)
const lastUpdated = ref(0)
const options = [
  { value: 120, label: '2分钟' },
  { value: 300, label: '5分钟' },
  { value: 600, label: '10分钟' },
  { value: 1800, label: '30分钟' },
  { value: 3600, label: '1小时' },
  { value: 43200, label: '12小时' },
  { value: 86400, label: '24小时' },
]
watch(timeDensity, () => {
  void refreshAll()
})

let latestRefreshID = 0

function isLatestRefresh(refreshID: number): boolean {
  return refreshID === latestRefreshID
}

// §7.1 序列色统一从设计 token 解析，色值缺失时回退 ECharts 默认调色板
function seriesColorStyle(color: string) {
  return color
    ? { lineStyle: { width: 2, color }, itemStyle: { color } }
    : { lineStyle: { width: 2 } }
}

interface DiskSeriesName {
  read: string
  write: string
}

interface NetSeriesName {
  receive: string
  send: string
}

// §7.1 多序列图例按稿成文：磁盘「读 / 写」、网络「接收 / 发送」；多设备时补设备标识
function diskSeriesNames(devices: string[]): DiskSeriesName[] {
  const multi = devices.length > 1
  return devices.map(device => ({
    read: multi ? `${device} 读` : '读',
    write: multi ? `${device} 写` : '写',
  }))
}

function netSeriesNames(ethernets: string[]): NetSeriesName[] {
  const multi = ethernets.length > 1
  return ethernets.map(ethernet => ({
    receive: multi ? `${ethernet} 接收` : '接收',
    send: multi ? `${ethernet} 发送` : '发送',
  }))
}

// CPU
const cpuPercent = ref('0.0%')
const cpuCores = ref(0)
const cpuStale = ref(false)
const cpuTimestamp = ref(0)
const cpuFreshness = ref<Freshness>({ collected_at: 0, age_seconds: 0, stale: false, degraded: false })
const cpuOption = reactive<EChartsOption>(cloneDeep(cpuTrendingOption))
async function renderCPUPercent(refreshID: number) {
  const { data } = await queryCPUInfo({ ...agentParams.value } as any)
  if (!isLatestRefresh(refreshID))
    return
  cpuPercent.value = `${data.percent.toFixed(1)}%`
  cpuCores.value = Number(data.cores) || 0
  cpuStale.value = Boolean(data.stale)
  cpuTimestamp.value = data.timestamp
  if (data.freshness)
    cpuFreshness.value = data.freshness
}
async function renderCPU(param: MonitorTimeRange, refreshID: number) {
  const { data } = await queryCPUUsage(param as any)
  if (!isLatestRefresh(refreshID))
    return
  const cpuData = data.data
  // §7.1 CPU 序列色取 chart-1 token
  const color = resolveChartPalette()[0]
  set(cpuOption, 'legend.data', ['CPU 使用率'])
  set(cpuOption, 'legend.textStyle.color', resolveChartLegendTextColor())
  set(cpuOption, 'series', [{
    name: 'CPU 使用率',
    data: cpuData.map((item: any) => toMonitorSeriesPoint(item.timestamp, Number(item.value.toFixed(1)))),
    type: 'line',
    smooth: true,
    showSymbol: false,
    ...seriesColorStyle(color),
  }])
}

// Memory
const memInfo = ref({ percent: '0%', total: '0', used: '0', stale: false, timestamp: 0 })
const memFreshness = ref<Freshness>({ collected_at: 0, age_seconds: 0, stale: false, degraded: false })
const memOption = reactive<EChartsOption>(cloneDeep(memTrendingOption))
async function renderMemInfo(refreshID: number) {
  const { data } = await queryMemInfo({ ...agentParams.value } as any)
  if (!isLatestRefresh(refreshID))
    return
  memInfo.value.percent = `${data.percent.toFixed(1)}%`
  memInfo.value.total = convertBytesToReadable(data.total)
  memInfo.value.used = convertBytesToReadable(data.used)
  memInfo.value.stale = Boolean(data.stale)
  memInfo.value.timestamp = data.timestamp
  if (data.freshness)
    memFreshness.value = data.freshness
}
async function renderMem(param: MonitorTimeRange, refreshID: number) {
  const { data } = await queryMemUsage(param as any)
  if (!isLatestRefresh(refreshID))
    return
  const memData = data.data
  // §7.1 内存序列色取 accent-warning 语义色
  const color = resolveChartTokenValue('--am-accent-warning')
  set(memOption, 'legend.data', ['内存使用率'])
  set(memOption, 'legend.textStyle.color', resolveChartLegendTextColor())
  set(memOption, 'series', [{
    name: '内存使用率',
    data: memData.map((item: any) => toMonitorSeriesPoint(item.timestamp, Number(item.value.toFixed(1)))),
    type: 'line',
    smooth: true,
    showSymbol: false,
    ...seriesColorStyle(color),
  }])
}

// Disk
const diskInfo = ref<{ device: string, total: string, used: string, percent: string, stale?: boolean, timestamp?: number }[]>([])
const diskIo = ref({ read: '-', write: '-' })
const diskFreshness = ref<Freshness>({ collected_at: 0, age_seconds: 0, stale: false, degraded: false })
const diskStale = computed(() => diskInfo.value.some(item => item.stale) || diskFreshness.value.stale)
const diskTimestamp = computed(() => diskInfo.value.find(item => item.stale)?.timestamp || diskInfo.value[0]?.timestamp || 0)
const diskOption = reactive<EChartsOption>(cloneDeep(diskTrendingOption))
async function renderDiskInfo(refreshID: number) {
  const { data } = await queryDiskInfo({ ...agentParams.value } as any)
  if (!isLatestRefresh(refreshID))
    return
  diskInfo.value = (data.info || []).map((item: any) => ({
    device: item.device,
    total: convertBytesToReadable(item.total),
    used: convertBytesToReadable(item.used),
    percent: `${item.percent.toFixed(1)}%`,
    stale: Boolean(item.stale),
    timestamp: item.timestamp,
  }))
  if (data.freshness)
    diskFreshness.value = data.freshness
}
function generateDiskSeriesData(data: DiskUsage[], names: DiskSeriesName[]) {
  const series: any[] = []
  // §7.1 读序列取 chart-1、写序列取 chart-5
  const palette = resolveChartPalette()
  data.forEach((i: DiskUsage, index: number) => {
    const name = names[index]
    series.push({ name: name.read, data: i.data.map((v: DiskIO) => toMonitorSeriesPoint(v.timestamp, v.io_read)), type: 'line', smooth: true, showSymbol: false, ...seriesColorStyle(palette[0]) })
    series.push({ name: name.write, data: i.data.map((v: DiskIO) => toMonitorSeriesPoint(v.timestamp, v.io_write)), type: 'line', smooth: true, showSymbol: false, ...seriesColorStyle(palette[4]) })
  })
  return series
}
async function renderDisk(param: MonitorTimeRange, refreshID: number) {
  const { data } = await queryDiskUsage(param as any)
  if (!isLatestRefresh(refreshID))
    return
  if (!data.usage || data.usage.length === 0) {
    set(diskOption, 'series', [])
    set(diskOption, 'legend.data', [])
    diskIo.value = { read: '-', write: '-' }
    return
  }
  const lastIo = data.usage[0]?.data?.[data.usage[0].data.length - 1]
  diskIo.value = lastIo
    ? { read: formatBytesPerSecond(lastIo.io_read), write: formatBytesPerSecond(lastIo.io_write) }
    : { read: '-', write: '-' }
  const names = diskSeriesNames(data.usage.map((item: DiskUsage) => item.device))
  set(diskOption, 'legend.data', names.flatMap(name => [name.read, name.write]))
  set(diskOption, 'legend.textStyle.color', resolveChartLegendTextColor())
  set(diskOption, 'series', generateDiskSeriesData(data.usage, names))
}

// Network
const netInfo = ref<{ ethernet: string, read: string, write: string }[]>([])
const netTimestamp = ref(0)
const netOption = reactive<EChartsOption>(cloneDeep(netTrendingOption))
async function renderNet(param: MonitorTimeRange, refreshID: number) {
  const { data } = await queryNetworkUsage(param as any)
  if (!isLatestRefresh(refreshID))
    return
  if (!data.usage || data.usage.length === 0) {
    netInfo.value = []
    netTimestamp.value = 0
    set(netOption, 'series', [])
    set(netOption, 'legend.data', [])
    return
  }
  netTimestamp.value = data.usage[0]?.data?.[data.usage[0].data.length - 1]?.timestamp || 0
  netInfo.value = data.usage.map((item: NetUsage) => ({
    ethernet: item.ethernet,
    read: formatBytesPerSecond(item.data[item.data.length - 1]?.bytes_recv || 0),
    write: formatBytesPerSecond(item.data[item.data.length - 1]?.bytes_sent || 0),
  }))
  const series: any[] = []
  // §7.1 接收序列取 chart-1、发送序列取 chart-5
  const palette = resolveChartPalette()
  const names = netSeriesNames(data.usage.map((item: NetUsage) => item.ethernet))
  data.usage.forEach((i: NetUsage, index: number) => {
    const name = names[index]
    series.push({ name: name.receive, data: i.data.map((v: NetIO) => toMonitorSeriesPoint(v.timestamp, v.bytes_recv)), type: 'line', smooth: true, showSymbol: false, ...seriesColorStyle(palette[0]) })
    series.push({ name: name.send, data: i.data.map((v: NetIO) => toMonitorSeriesPoint(v.timestamp, v.bytes_sent)), type: 'line', smooth: true, showSymbol: false, ...seriesColorStyle(palette[4]) })
  })
  set(netOption, 'legend.data', names.flatMap(name => [name.receive, name.send]))
  set(netOption, 'legend.textStyle.color', resolveChartLegendTextColor())
  set(netOption, 'series', series)
}

const { t } = useI18n()

// §7.1 图卡指标行：标题下方的 mono 指标值（语义色）；核数是稿面的「· N 核」上下文
const cpuMetric = computed(() => cpuCores.value > 0 ? `当前 ${cpuPercent.value} · ${cpuCores.value} 核` : `当前 ${cpuPercent.value}`)
const memMetric = computed(() => `${memInfo.value.used} / ${memInfo.value.total} · ${memInfo.value.percent}`)
const diskMaxPercent = computed(() => {
  if (!diskInfo.value.length)
    return '-'
  const max = Math.max(...diskInfo.value.map(item => Number.parseFloat(item.percent) || 0))
  return `${max.toFixed(1)}%`
})
const diskMetric = computed(() => `读 ${diskIo.value.read} / 写 ${diskIo.value.write} · 最高 ${diskMaxPercent.value}`)
const netMetric = computed(() => `接收 ${netInfo.value[0]?.read ?? '-'} · 发送 ${netInfo.value[0]?.write ?? '-'}`)

// §7.5 数据过期：任一主机指标过期即进入过期态；卡片保留最后采样值，不判断阈值。
const containerStale = ref({ avgCpu: '-', memSum: '-', lastSampleAt: 0 })

async function renderContainerStale(param: MonitorTimeRange, refreshID: number) {
  try {
    const { data } = await queryContainersUsage(param as any)
    if (!isLatestRefresh(refreshID))
      return
    const cpuMap = (data.cpu_usage || {}) as unknown as Record<string, { value: number }[]>
    const memMap = (data.mem_usage || {}) as unknown as Record<string, { value: number }[]>
    const cpuLast: number[] = []
    const memLast: number[] = []
    Object.values(cpuMap).forEach((series) => {
      const point = series?.[series.length - 1]
      if (point && Number.isFinite(point.value))
        cpuLast.push(point.value)
    })
    Object.values(memMap).forEach((series) => {
      const point = series?.[series.length - 1]
      if (point && Number.isFinite(point.value))
        memLast.push(point.value)
    })
    const avgCpu = cpuLast.length
      ? `平均 ${(cpuLast.reduce((sum, value) => sum + value, 0) / cpuLast.length).toFixed(1)}%`
      : '-'
    const memSum = memLast.length
      ? `合计 ${convertBytesToReadable(memLast.reduce((sum, value) => sum + value, 0))}`
      : '-'
    containerStale.value = {
      avgCpu,
      memSum,
      lastSampleAt: data.freshness?.collected_at || 0,
    }
  }
  catch {
    containerStale.value = { avgCpu: '-', memSum: '-', lastSampleAt: 0 }
  }
}

const staleActive = computed(() => cpuFreshness.value.stale || memFreshness.value.stale || diskFreshness.value.stale)

const staleLastSampleAt = computed(() => Math.max(
  cpuFreshness.value.collected_at,
  memFreshness.value.collected_at,
  diskFreshness.value.collected_at,
  netTimestamp.value,
  containerStale.value.lastSampleAt,
))

const staleKpis = computed<StaleKpi[]>(() => {
  const cores = cpuCores.value > 0 ? `${cpuCores.value} 核 · ` : ''
  return [
    {
      key: 'cpu',
      label: 'CPU 使用率',
      value: cpuPercent.value,
      note: `${cores}最近一次采样来自 ${formatSampleClock(cpuTimestamp.value)}`,
      lastSampleAt: cpuTimestamp.value,
    },
    {
      key: 'memory',
      label: '内存使用率',
      value: memInfo.value.percent,
      note: `${memInfo.value.used} / ${memInfo.value.total} · 采样停在上报中断前`,
      lastSampleAt: memInfo.value.timestamp,
    },
    {
      key: 'disk',
      label: '磁盘使用率',
      value: diskMaxPercent.value,
      note: `${diskInfo.value[0]?.device ?? '-'} · 数据过期期间不再判断告警阈值`,
      lastSampleAt: diskTimestamp.value || diskFreshness.value.collected_at,
    },
    {
      key: 'network',
      label: '网络流量',
      value: `收 ${netInfo.value[0]?.read ?? '-'}`,
      note: `${netInfo.value[0]?.ethernet ?? '-'} · 收/发采样均停在上报中断前`,
      lastSampleAt: netTimestamp.value,
    },
  ]
})

const staleRows = computed<StaleDetailRow[]>(() => {
  const rows: StaleDetailRow[] = [
    { key: 'cpu', label: 'CPU 使用率', value: cpuPercent.value, lastSampleAt: cpuTimestamp.value },
    { key: 'memory', label: '内存使用率', value: memInfo.value.percent, lastSampleAt: memInfo.value.timestamp },
    { key: 'disk_io', label: '磁盘 IO', value: `读 ${diskIo.value.read} / 写 ${diskIo.value.write}`, lastSampleAt: diskTimestamp.value || diskFreshness.value.collected_at },
    { key: 'net', label: '网络流量', value: `收 ${netInfo.value[0]?.read ?? '-'} / 发 ${netInfo.value[0]?.write ?? '-'}`, lastSampleAt: netTimestamp.value },
  ]
  if (containerStale.value.avgCpu !== '-')
    rows.push({ key: 'container_cpu', label: '容器 CPU 使用率', value: containerStale.value.avgCpu, lastSampleAt: containerStale.value.lastSampleAt })
  if (containerStale.value.memSum !== '-')
    rows.push({ key: 'container_mem', label: '容器内存使用率', value: containerStale.value.memSum, lastSampleAt: containerStale.value.lastSampleAt })
  rows.push({ key: 'disk_usage', label: '磁盘使用率', value: `${diskInfo.value[0]?.device ?? '-'} ${diskMaxPercent.value}`, lastSampleAt: diskTimestamp.value || diskFreshness.value.collected_at })
  return rows
})

async function refreshAll(): Promise<void> {
  const refreshID = ++latestRefreshID
  if (!currentAgent.value)
    return
  const endTime = dayjs().unix()
  const param: MonitorTimeRange = {
    start_time: endTime - timeDensity.value,
    end_time: endTime,
    ...agentParams.value,
  }
  applyMonitorTimeRange(cpuOption, param)
  applyMonitorTimeRange(memOption, param)
  applyMonitorTimeRange(diskOption, param)
  applyMonitorTimeRange(netOption, param)
  await Promise.allSettled([
    renderCPUPercent(refreshID),
    renderCPU(param, refreshID),
    renderMemInfo(refreshID),
    renderMem(param, refreshID),
    renderDiskInfo(refreshID),
    renderDisk(param, refreshID),
    renderNet(param, refreshID),
  ])
  // 过期态需要容器侧最后采样值（明细表的两行），正常态不额外请求。
  if (isLatestRefresh(refreshID)) {
    if (staleActive.value)
      await renderContainerStale(param, refreshID)
    else
      containerStale.value = { avgCpu: '-', memSum: '-', lastSampleAt: 0 }
  }
  lastUpdated.value = Date.now()
}

async function refreshAgents() {
  await loadAgents()
  await refreshAll()
}

const timer = ref()
onMounted(async () => {
  await ensureSelectedAgent()
  void refreshAll()
  timer.value = setInterval(() => void refreshAll(), 10000)
})
onUnmounted(() => {
  clearInterval(timer.value)
})
</script>

<template>
    <!-- Host Section -->
    <div class="am-section host-monitor">
        <div class="am-section-header">
            <div class="am-section-title-group">
                <span class="am-section-title">{{ t('menu.hostMonitor') }}</span>
                <AgentSelect width="180px" />
            </div>
            <MonitorContextBar v-model:density="timeDensity" :density-options="options" :last-updated="lastUpdated" :stale="staleActive" />
        </div>

        <AgentEmptyState v-if="isAgentEmpty" @refresh="refreshAgents" />
        <MonitorStalePanel
            v-else-if="staleActive"
            :last-sample-at="staleLastSampleAt"
            :kpis="staleKpis"
            :rows="staleRows"
        />
        <div v-else class="am-chart-grid">
            <div class="am-chart-row">
                <div class="am-chart-card">
                    <div class="am-chart-card-header">
                        <span class="am-chart-card-icon am-chart-card-icon--primary">
                            <IconCpu />
                        </span>
                        <div class="am-chart-card-title-group">
                            <span class="am-chart-card-title">CPU 使用率</span>
                            <span class="am-chart-card-percent accent-primary">{{ cpuMetric }}</span>
                        </div>
                        <DataStaleTag :stale="cpuStale" :timestamp="cpuTimestamp" />
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
                        <div class="am-chart-card-title-group">
                            <span class="am-chart-card-title">内存使用率</span>
                            <span class="am-chart-card-percent accent-warning">{{ memMetric }}</span>
                        </div>
                        <DataStaleTag :stale="memInfo.stale" :timestamp="memInfo.timestamp" />
                    </div>
                    <div class="am-chart-area">
                        <echarts :option="memOption" />
                    </div>
                </div>
            </div>
            <div class="am-chart-row">
                <div class="am-chart-card">
                    <div class="am-chart-card-header">
                        <span class="am-chart-card-icon am-chart-card-icon--success">
                            <IconHardDrive />
                        </span>
                        <div class="am-chart-card-title-group">
                            <span class="am-chart-card-title">磁盘 IO</span>
                            <span class="am-chart-card-percent accent-success">{{ diskMetric }}</span>
                        </div>
                        <DataStaleTag :stale="diskStale" :timestamp="diskTimestamp" />
                    </div>
                    <div class="am-chart-area">
                        <echarts :option="diskOption" />
                    </div>
                </div>
                <div class="am-chart-card">
                    <div class="am-chart-card-header">
                        <span class="am-chart-card-icon am-chart-card-icon--network">
                            <IconNetwork />
                        </span>
                        <div class="am-chart-card-title-group">
                            <span class="am-chart-card-title">网络流量</span>
                            <span class="am-chart-card-percent accent-chart-3">{{ netMetric }}</span>
                        </div>
                        <DataStaleTag :stale="staleActive" :timestamp="netTimestamp" />
                    </div>
                    <div class="am-chart-area">
                        <echarts :option="netOption" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">
// 主机监控的高度约束已统一在全局 styles/element.scss 的 .am-chart-area 中，
// 此处仅保留局部强调色样式，避免重复覆盖 row/area 高度。

.accent-primary {
  color: var(--am-accent-primary);
}
.accent-warning {
  color: var(--am-accent-warning);
}
.accent-success {
  color: var(--am-accent-success);
}

.accent-chart-3 {
  color: var(--am-chart-3);
}
</style>
