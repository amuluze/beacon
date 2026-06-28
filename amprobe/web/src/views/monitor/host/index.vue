<script setup lang="ts">
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import type { DiskIO, DiskUsage, NetIO, NetUsage } from '@/interface/host.ts'

import { queryCPUInfo, queryCPUUsage, queryDiskInfo, queryDiskUsage, queryMemInfo, queryMemUsage, queryNetworkUsage, queryAgentList } from '@/api/host'
import { cpuTrendingOption, memTrendingOption, diskTrendingOption, netTrendingOption } from '@/config/echarts.ts'
import { convertBytesToReadable } from '@/utils/convert.ts'
import useStore from '@/store'
import { dayjs } from 'element-plus'
import { set } from 'lodash-es'
import { useI18n } from 'vue-i18n'

// Agent switcher
const store = useStore()
const currentAgent = computed({
  get: () => store.agent.currentAgentID,
  set: (value: string) => store.agent.setCurrentAgent(value),
})
async function loadAgents() {
  try {
    const { data } = await queryAgentList()
    store.agent.setAgents(data || [])
  }
  catch {
    store.agent.setAgents([])
  }
}
watch(currentAgent, () => { refreshAll() })

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
watch(timeDensity, () => { refreshAll() })

// CPU
const cpuPercent = ref('0.0%')
const cpuOption = reactive<EChartsOption>(JSON.parse(JSON.stringify(cpuTrendingOption)))
async function renderCPUPercent() {
  const { data } = await queryCPUInfo()
  cpuPercent.value = `${data.percent.toFixed(1)}%`
}
async function renderCPU() {
  const param = { start_time: dayjs().unix() - timeDensity.value, end_time: dayjs().unix() }
  const { data } = await queryCPUUsage(param as any)
  const cpuData = data.data
  const labels = cpuData.map((item: any) => `${dayjs(item.timestamp * 1000).format('HH:mm')}`)
  set(cpuOption, 'xAxis.data', labels)
  set(cpuOption, 'legend.data', ['CPU 使用率'])
  set(cpuOption, 'series', [{ name: 'CPU 使用率', data: cpuData.map((item: any) => item.value.toFixed(1)), type: 'line', smooth: true, showSymbol: false }])
}

// Memory
const memInfo = ref({ percent: '0%', total: '0', used: '0' })
const memOption = reactive<EChartsOption>(JSON.parse(JSON.stringify(memTrendingOption)))
async function renderMemInfo() {
  const { data } = await queryMemInfo()
  memInfo.value.percent = `${data.percent.toFixed(1)}%`
  memInfo.value.total = convertBytesToReadable(data.total)
  memInfo.value.used = convertBytesToReadable(data.used)
}
async function renderMem() {
  const param = { start_time: dayjs().unix() - timeDensity.value, end_time: dayjs().unix() }
  const { data } = await queryMemUsage(param as any)
  const memData = data.data
  const labels = memData.map((item: any) => `${dayjs(item.timestamp * 1000).format('HH:mm')}`)
  set(memOption, 'xAxis.data', labels)
  set(memOption, 'legend.data', ['内存使用率'])
  set(memOption, 'series', [{ name: '内存使用率', data: memData.map((item: any) => item.value.toFixed(1)), type: 'line', smooth: true, showSymbol: false }])
}

// Disk
const diskInfo = ref<{ device: string; total: string; used: string; percent: string }[]>([])
const diskOption = reactive<EChartsOption>(JSON.parse(JSON.stringify(diskTrendingOption)))
async function renderDiskInfo() {
  const { data } = await queryDiskInfo()
  diskInfo.value = (data.info || []).map((item: any) => ({
    device: item.device,
    total: convertBytesToReadable(item.total),
    used: convertBytesToReadable(item.used),
    percent: `${item.percent.toFixed(1)}%`,
  }))
}
function generateDiskSeriesData(data: DiskUsage[]) {
  const series: any[] = []
  data.forEach((i: DiskUsage) => {
    series.push({ name: `${i.device}_Read`, data: i.data.map((v: DiskIO) => v.io_read), type: 'line', smooth: true, showSymbol: false })
    series.push({ name: `${i.device}_Write`, data: i.data.map((v: DiskIO) => v.io_write), type: 'line', smooth: true, showSymbol: false })
  })
  return series
}
async function renderDisk() {
  const param = { start_time: dayjs().unix() - timeDensity.value, end_time: dayjs().unix() }
  const { data } = await queryDiskUsage(param as any)
  if (!data.usage || data.usage.length === 0) return
  const labels = data.usage[0].data.map((item: DiskIO) => `${dayjs(item.timestamp * 1000).format('HH:mm')}`)
  set(diskOption, 'xAxis.data', labels)
  set(diskOption, 'series', generateDiskSeriesData(data.usage))
}

// Network
const netInfo = ref<{ ethernet: string; read: string; write: string }[]>([])
const netOption = reactive<EChartsOption>(JSON.parse(JSON.stringify(netTrendingOption)))
async function renderNet() {
  const param = { start_time: dayjs().unix() - timeDensity.value, end_time: dayjs().unix() }
  const { data } = await queryNetworkUsage(param as any)
  if (!data.usage || data.usage.length === 0) return
  netInfo.value = data.usage.map((item: NetUsage) => ({
    ethernet: item.ethernet,
    read: convertBytesToReadable(item.data[item.data.length - 1].bytes_recv),
    write: convertBytesToReadable(item.data[item.data.length - 1].bytes_sent),
  }))
  const labels = data.usage[0].data.map((item: NetIO) => `${dayjs(item.timestamp * 1000).format('HH:mm')}`)
  set(netOption, 'xAxis.data', labels)
  const series: any[] = []
  data.usage.forEach((i: NetUsage) => {
    series.push({ name: `${i.ethernet}_Receive`, data: i.data.map((v: NetIO) => v.bytes_recv), type: 'line', smooth: true, showSymbol: false })
    series.push({ name: `${i.ethernet}_Send`, data: i.data.map((v: NetIO) => v.bytes_sent), type: 'line', smooth: true, showSymbol: false })
  })
  set(netOption, 'series', series)
}

function refreshAll() {
  renderCPUPercent()
  renderCPU()
  renderMemInfo()
  renderMem()
  renderDiskInfo()
  renderDisk()
  renderNet()
}

const timer = ref()
onMounted(() => {
  loadAgents()
  refreshAll()
  timer.value = setInterval(refreshAll, 10000)
})
onUnmounted(() => { clearInterval(timer.value) })

const { t } = useI18n()
</script>

<template>
  <!-- Host Section -->
  <div class="am-section">
    <div class="am-section-header">
      <div class="am-section-title-group">
        <span class="am-section-title">{{ t('monitor.hostMonitor') }}</span>
        <el-select v-model="currentAgent" size="small" style="width: 160px" placeholder="选择主机">
          <el-option
            v-for="item in store.agent.agents"
            :key="item.agent_id"
            :label="item.hostname || item.agent_id"
            :value="item.agent_id"
          />
        </el-select>
      </div>
      <div class="am-density-group">
        <span class="am-density-label">{{ t('monitor.timeDensity') }}：</span>
        <el-select v-model="timeDensity" size="small" style="width: 110px">
          <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </div>
    </div>

    <div class="am-chart-grid">
      <div class="am-chart-row">
        <div class="am-chart-card">
          <div class="am-chart-card-header">
            <span class="am-chart-card-title">CPU 使用率</span>
            <span class="am-chart-card-percent accent-primary">{{ cpuPercent }}</span>
          </div>
          <div class="am-chart-area">
            <echarts :option="cpuOption" />
          </div>
        </div>
        <div class="am-chart-card">
          <div class="am-chart-card-header">
            <span class="am-chart-card-title">内存使用率</span>
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
            <span class="am-chart-card-title">磁盘 IO</span>
            <span class="am-chart-card-percent accent-success" v-if="diskInfo.length > 0">{{ diskInfo[0].percent }}</span>
          </div>
          <div class="am-chart-area">
            <echarts :option="diskOption" />
          </div>
        </div>
        <div class="am-chart-card">
          <div class="am-chart-card-header">
            <span class="am-chart-card-title">网络流量</span>
            <span class="am-chart-card-percent accent-primary" v-if="netInfo.length > 0">{{ netInfo[0].read }} / {{ netInfo[0].write }}</span>
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
.am-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
}
.am-section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
}
.am-section-title-group {
  display: flex;
  align-items: center;
  gap: 10px;
}
.am-section-title {
  font-size: 15px;
  font-weight: 600;
  color: #1a1a2e;
}
.am-density-group {
  display: flex;
  align-items: center;
  gap: 6px;
}
.am-density-label {
  font-size: 13px;
  color: #666;
}
.am-chart-grid {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 0 12px 12px;
}
.am-chart-row {
  flex: 1;
  display: flex;
  gap: 12px;
}
.am-chart-card {
  flex: 1;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.06);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.am-chart-card-header {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.am-chart-card-title {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a2e;
}
.am-chart-card-percent {
  font-size: 12px;
  font-family: 'Geist Mono', monospace;
  font-weight: 500;
}
.accent-primary { color: #4f7cff; }
.accent-warning { color: #f5a623; }
.accent-success { color: #52c41a; }
.am-chart-area {
  flex: 1;
  min-height: 180px;
}
</style>
