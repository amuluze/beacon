<script setup lang="ts">
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import type { DiskIO, DiskUsage, NetIO, NetUsage } from '@/interface/host.ts'

import { queryCPUInfo, queryCPUUsage, queryDiskInfo, queryDiskUsage, queryMemInfo, queryMemUsage, queryNetworkUsage } from '@/api/host'
import AgentEmptyState from '@/components/Agent/AgentEmptyState.vue'
import DataStaleTag from '@/components/Agent/DataStaleTag.vue'
import { cpuTrendingOption, diskTrendingOption, formatBytesPerSecond, memTrendingOption, netTrendingOption } from '@/config/echarts.ts'
import { convertBytesToReadable } from '@/utils/convert.ts'
import { dayjs } from 'element-plus'
import { cloneDeep, set } from 'lodash-es'
import { useAgentSelection } from '@/hooks/useAgentSelection'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { applyMonitorTimeRange, toMonitorSeriesPoint } from '@/utils/monitorChart'
import type { MonitorTimeRange } from '@/utils/monitorChart'

const router = useRouter()

// Agent switcher
const { selectedAgentID: currentAgent, isAgentEmpty, agentParams, ensureSelectedAgent, loadAgents } = useAgentSelection({ immediate: false })
function openTerminal(): void {
  if (!currentAgent.value)
    return
  router.push({ path: '/terminal', query: { agent_id: currentAgent.value } })
}
watch(currentAgent, () => {
  void refreshAll()
})

// Time density
const timeDensity = ref(600)
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
    return
  }
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
</script>

<template>
    <!-- Host Section -->
    <div class="am-section host-monitor">
        <div class="am-section-header">
            <div class="am-section-title-group">
                <span class="am-section-title">{{ t('menu.hostMonitor') }}</span>
                <el-button size="small" type="primary" plain :disabled="!currentAgent" @click="openTerminal">
                    {{ t('agent.openTerminal') }}
                </el-button>
            </div>
            <div class="am-density-group">
                <span class="am-density-label">{{ t('monitor.timeDensity') }}：</span>
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
                        <div class="am-chart-card-title-row">
                            <span class="am-chart-card-title">CPU 使用率</span>
                            <DataStaleTag :stale="cpuStale" :timestamp="cpuTimestamp" />
                        </div>
                        <span class="am-chart-card-percent accent-primary">{{ cpuPercent }}</span>
                    </div>
                    <div class="am-chart-area">
                        <echarts :option="cpuOption" />
                    </div>
                </div>
                <div class="am-chart-card">
                    <div class="am-chart-card-header">
                        <div class="am-chart-card-title-row">
                            <span class="am-chart-card-title">内存使用率</span>
                            <DataStaleTag :stale="memInfo.stale" :timestamp="memInfo.timestamp" />
                        </div>
                        <span class="am-chart-card-percent accent-warning">{{ memInfo.percent }}</span>
                    </div>
                    <div class="am-chart-area">
                        <echarts :option="memOption" />
                    </div>
                </div>
            </div>
            <div class="am-chart-row">
                <div class="am-chart-card">
                    <div class="am-chart-card-header">
                        <div class="am-chart-card-title-row">
                            <span class="am-chart-card-title">磁盘 IO</span>
                            <DataStaleTag :stale="diskStale" :timestamp="diskTimestamp" />
                        </div>
                        <span v-if="diskInfo.length > 0" class="am-chart-card-percent accent-success">{{ diskInfo[0].percent }}</span>
                    </div>
                    <div class="am-chart-area">
                        <echarts :option="diskOption" />
                    </div>
                </div>
                <div class="am-chart-card">
                    <div class="am-chart-card-header">
                        <span class="am-chart-card-title">网络流量</span>
                        <span v-if="netInfo.length > 0" class="am-chart-card-percent accent-primary">{{ netInfo[0].read }} / {{ netInfo[0].write }}</span>
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
</style>
