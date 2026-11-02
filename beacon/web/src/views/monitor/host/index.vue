<script setup lang="ts">
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import type { DiskIO, DiskUsage, NetIO, NetUsage } from '@/interface/host.ts'

import { queryCPUInfo, queryCPUUsage, queryDiskInfo, queryDiskUsage, queryMemInfo, queryMemUsage, queryNetworkUsage } from '@/api/host'
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import DataStaleTag from '@/components/Agent/DataStaleTag.vue'
import MonitorContextBar from '@/components/Agent/MonitorContextBar.vue'
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

// CPU
const cpuPercent = ref('0.0%')
const cpuStale = ref(false)
const cpuTimestamp = ref(0)
const cpuOption = reactive<EChartsOption>(cloneDeep(cpuTrendingOption))
async function renderCPUPercent(refreshID: number) {
  const { data } = await queryCPUInfo({ ...agentParams.value } as any)
  if (!isLatestRefresh(refreshID))
    return
  cpuPercent.value = `${data.percent.toFixed(1)}%`
  cpuStale.value = Boolean(data.stale)
  cpuTimestamp.value = data.timestamp
}
async function renderCPU(param: MonitorTimeRange, refreshID: number) {
  const { data } = await queryCPUUsage(param as any)
  if (!isLatestRefresh(refreshID))
    return
  const cpuData = data.data
  set(cpuOption, 'legend.data', ['CPU 使用率'])
  set(cpuOption, 'series', [{
    name: 'CPU 使用率',
    data: cpuData.map((item: any) => toMonitorSeriesPoint(item.timestamp, Number(item.value.toFixed(1)))),
    type: 'line',
    smooth: true,
    showSymbol: false,
  }])
}

// Memory
const memInfo = ref({ percent: '0%', total: '0', used: '0', stale: false, timestamp: 0 })
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
}
async function renderMem(param: MonitorTimeRange, refreshID: number) {
  const { data } = await queryMemUsage(param as any)
  if (!isLatestRefresh(refreshID))
    return
  const memData = data.data
  set(memOption, 'legend.data', ['内存使用率'])
  set(memOption, 'series', [{
    name: '内存使用率',
    data: memData.map((item: any) => toMonitorSeriesPoint(item.timestamp, Number(item.value.toFixed(1)))),
    type: 'line',
    smooth: true,
    showSymbol: false,
  }])
}

// Disk
const diskInfo = ref<{ device: string, total: string, used: string, percent: string, stale?: boolean, timestamp?: number }[]>([])
const diskIo = ref({ read: '-', write: '-' })
const diskStale = computed(() => diskInfo.value.some(item => item.stale))
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
}
function generateDiskSeriesData(data: DiskUsage[]) {
  const series: any[] = []
  data.forEach((i: DiskUsage) => {
    series.push({ name: `${i.device}_Read`, data: i.data.map((v: DiskIO) => toMonitorSeriesPoint(v.timestamp, v.io_read)), type: 'line', smooth: true, showSymbol: false })
    series.push({ name: `${i.device}_Write`, data: i.data.map((v: DiskIO) => toMonitorSeriesPoint(v.timestamp, v.io_write)), type: 'line', smooth: true, showSymbol: false })
  })
  return series
}
async function renderDisk(param: MonitorTimeRange, refreshID: number) {
  const { data } = await queryDiskUsage(param as any)
  if (!isLatestRefresh(refreshID))
    return
  if (!data.usage || data.usage.length === 0) {
    set(diskOption, 'series', [])
    diskIo.value = { read: '-', write: '-' }
    return
  }
  const lastIo = data.usage[0]?.data?.[data.usage[0].data.length - 1]
  diskIo.value = lastIo
    ? { read: formatBytesPerSecond(lastIo.io_read), write: formatBytesPerSecond(lastIo.io_write) }
    : { read: '-', write: '-' }
  set(diskOption, 'series', generateDiskSeriesData(data.usage))
}

// Network
const netInfo = ref<{ ethernet: string, read: string, write: string }[]>([])
const netOption = reactive<EChartsOption>(cloneDeep(netTrendingOption))
async function renderNet(param: MonitorTimeRange, refreshID: number) {
  const { data } = await queryNetworkUsage(param as any)
  if (!isLatestRefresh(refreshID))
    return
  if (!data.usage || data.usage.length === 0) {
    netInfo.value = []
    set(netOption, 'series', [])
    return
  }
  netInfo.value = data.usage.map((item: NetUsage) => ({
    ethernet: item.ethernet,
    read: formatBytesPerSecond(item.data[item.data.length - 1]?.bytes_recv || 0),
    write: formatBytesPerSecond(item.data[item.data.length - 1]?.bytes_sent || 0),
  }))
  const series: any[] = []
  data.usage.forEach((i: NetUsage) => {
    series.push({ name: `${i.ethernet}_Receive`, data: i.data.map((v: NetIO) => toMonitorSeriesPoint(v.timestamp, v.bytes_recv)), type: 'line', smooth: true, showSymbol: false })
    series.push({ name: `${i.ethernet}_Send`, data: i.data.map((v: NetIO) => toMonitorSeriesPoint(v.timestamp, v.bytes_sent)), type: 'line', smooth: true, showSymbol: false })
  })
  set(netOption, 'series', series)
}

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

const { t } = useI18n()

// §7.1 图卡指标行：标题下方的 mono 指标值（语义色）
const cpuMetric = computed(() => `当前 ${cpuPercent.value}`)
const memMetric = computed(() => `${memInfo.value.used} / ${memInfo.value.total} · ${memInfo.value.percent}`)
const diskMaxPercent = computed(() => {
  if (!diskInfo.value.length)
    return '-'
  const max = Math.max(...diskInfo.value.map(item => Number.parseFloat(item.percent) || 0))
  return `${max.toFixed(1)}%`
})
const diskMetric = computed(() => `读 ${diskIo.value.read} / 写 ${diskIo.value.write} · 最高 ${diskMaxPercent.value}`)
const netMetric = computed(() => `接收 ${netInfo.value[0]?.read ?? '-'} · 发送 ${netInfo.value[0]?.write ?? '-'}`)
</script>

<template>
    <!-- Host Section -->
    <div class="am-section host-monitor">
        <div class="am-section-header">
            <div class="am-section-title-group">
                <span class="am-section-title">{{ t('menu.hostMonitor') }}</span>
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
