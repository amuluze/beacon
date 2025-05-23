<script setup lang="ts">
import { queryCPUInfo, queryCPUUsage, queryDiskInfo, queryDiskUsage, queryMemInfo, queryMemUsage, queryNetworkUsage } from '@/api/host'
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import { cpuTrendingOption, diskTrendingOption, memTrendingOption, netTrendingOption } from '@/config/echarts.ts'
import type { CPUTrendingArgs, DiskIO, DiskTrendingArgs, DiskUsage, MemTrendingArgs, NetIO, NetTrendingArgs, NetUsage } from '@/interface/host.ts'
import { convertBytesToReadable } from '@/utils/convert.ts'
import { dayjs } from 'element-plus'
import { set } from 'lodash-es'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const loading = ref(false)

// 定义 props
const props = defineProps<{
  timeDensity: number
}>()

// 定义数据缓存
const dataCache = reactive({
  cpu: {
    lastUpdate: 0,
    data: null as any,
    interval: 5000, // CPU数据每5秒更新一次
  },
  mem: {
    lastUpdate: 0,
    data: null as any,
    interval: 5000, // 内存数据每5秒更新一次
  },
  disk: {
    lastUpdate: 0,
    data: null as any,
    interval: 5000, // 磁盘数据每10秒更新一次
  },
  net: {
    lastUpdate: 0,
    data: null as any,
    interval: 5000, // 网络数据每5秒更新一次
  },
})

// 检查数据是否需要更新
function shouldUpdate(type: keyof typeof dataCache): boolean {
  const now = Date.now()
  return !dataCache[type].data || (now - dataCache[type].lastUpdate) >= dataCache[type].interval
}

// CPU 使用率
const cpuPercent = ref('0.0%')
async function renderCPUPercent() {
  if (!shouldUpdate('cpu')) return
  const { data } = await queryCPUInfo()
  cpuPercent.value = `${data.percent.toFixed(2)}%`
  dataCache.cpu.lastUpdate = Date.now()
}

const cpuOption = reactive<EChartsOption>(cpuTrendingOption as EChartsOption)
async function renderCPU() {
  if (!shouldUpdate('cpu')) return
  const param: CPUTrendingArgs = {
    start_time: dayjs().unix() - props.timeDensity,
    end_time: dayjs().unix(),
  }

  const { data } = await queryCPUUsage(param)
  const cpuData = data.data
  set(
    cpuOption,
    'xAxis.data',
    cpuData.map(item => dayjs(item.timestamp * 1000).format("HH:mm:ss")),
  )
  set(cpuOption, 'legend.data', ['CPU使用率'])
  set(cpuOption, 'series', [
    {
      name: 'CPU使用率',
      data: cpuData.map(item => item.value.toFixed(2)),
      type: 'line',
      smooth: true,
      showSymbol: false,
    },
  ])
  dataCache.cpu.data = cpuData
  dataCache.cpu.lastUpdate = Date.now()
}

// 内存使用率
const memInfo = ref({
  percent: '0%',
  total: '0',
  used: '0',
})
async function renderMemInfo() {
  if (!shouldUpdate('mem')) return
  const { data } = await queryMemInfo()
  memInfo.value.percent = `${data.percent.toFixed(2)}%`
  memInfo.value.total = convertBytesToReadable(data.total)
  memInfo.value.used = convertBytesToReadable(data.used)
  dataCache.mem.lastUpdate = Date.now()
}
const memOption = reactive<EChartsOption>(memTrendingOption as EChartsOption)
async function renderMem() {
  if (!shouldUpdate('mem')) return
  const param: MemTrendingArgs = {
    start_time: dayjs().unix() - props.timeDensity,
    end_time: dayjs().unix(),
  }
  const { data } = await queryMemUsage(param)
  const memData = data.data
  set(
    memOption,
    'xAxis.data',
    memData.map(item => dayjs(item.timestamp * 1000).format("HH:mm:ss")),
  )
  set(memOption, 'legend.data', ['内存使用率'])
  set(memOption, 'series', [
    {
      name: '内存使用率',
      data: memData.map(item => item.value.toFixed(2)),
      type: 'line',
      smooth: true,
      showSymbol: false,
    },
  ])
  dataCache.mem.data = memData
  dataCache.mem.lastUpdate = Date.now()
}

// 磁盘使用率
const diskInfo = ref<
  {
    device: string
    total: string
    used: string
    percent: string
  }[]
>([])
async function renderDiskInfo() {
  if (!shouldUpdate('disk')) return
  const { data } = await queryDiskInfo()
  diskInfo.value = []
  data.info.map((item) => {
    diskInfo.value.push({
      device: item.device,
      total: convertBytesToReadable(item.total),
      used: convertBytesToReadable(item.used),
      percent: `${item.percent.toFixed(2)}%`,
    })
    return diskInfo
  })
  dataCache.disk.lastUpdate = Date.now()
}
const diskOption = reactive<EChartsOption>(diskTrendingOption as EChartsOption)
function generateDiskLegendData(data: DiskUsage[]) {
  const options: string[] = ['Read', 'Write']
  const res: string[] = []
  options.forEach((item: string) => {
    data.forEach((i: DiskUsage) => {
      res.push(`${i.device}_${item}`)
    })
  })
  return res
}
interface DiskSeriesData {
  name: string
  data: number[]
  type: string
  showSymbol: boolean
  symbolSize: number
  hoverAnimation: boolean
  smooth: boolean
}
function generateDiskSeriesData(data: DiskUsage[]) {
  const options: string[] = ['Read', 'Write']
  const series: DiskSeriesData[] = []
  options.forEach((item: string) => {
    data.forEach((i: DiskUsage) => {
      if (item === 'Read') {
        series.push({
          name: `${i.device}_${item}`,
          data: i.data.map((val: DiskIO) => val.io_read),
          type: 'line',
          smooth: true,
          showSymbol: true,
          symbolSize: 2,
          hoverAnimation: false,
        })
      }
      else {
        series.push({
          name: `${i.device}_${item}`,
          data: i.data.map((val: DiskIO) => val.io_write),
          type: 'line',
          smooth: true,
          showSymbol: true,
          symbolSize: 2,
          hoverAnimation: false,
        })
      }
    })
  })
  return series
}
async function renderDisk() {
  if (!shouldUpdate('disk')) return
  const param: DiskTrendingArgs = {
    start_time: dayjs().unix() - props.timeDensity,
    end_time: dayjs().unix(),
  }
  const { data } = await queryDiskUsage(param)
  const diskData = data.usage
  set(
    diskOption,
    'xAxis.data',
    diskData[0].data.map((item: DiskIO) => dayjs(item.timestamp * 1000).format("HH:mm:ss")),
  )
  set(diskOption, 'legend.data', generateDiskLegendData(diskData))
  set(diskOption, 'series', generateDiskSeriesData(diskData))
  dataCache.disk.data = diskData
  dataCache.disk.lastUpdate = Date.now()
}

// 网络使用率
const netInfo = ref<{
  ethernet: string
  read: string
  write: string
}[]>([])
async function renderNetInfo() {
  if (!shouldUpdate('net')) return
  const param: NetTrendingArgs = {
    start_time: dayjs().unix() - props.timeDensity,
    end_time: dayjs().unix(),
  }
  const { data } = await queryNetworkUsage(param)
  netInfo.value = data.usage.map(item => ({
    ethernet: item.ethernet,
    read: convertBytesToReadable(item.data[item.data.length - 1].bytes_recv),
    write: convertBytesToReadable(item.data[item.data.length - 1].bytes_sent)
  }))
  dataCache.net.lastUpdate = Date.now()
}
const netOption = reactive<EChartsOption>(netTrendingOption as EChartsOption)
function generateNetLegendData(data: NetUsage[]) {
  const options: string[] = ['Receive', 'Send']
  const res: string[] = []
  options.forEach((item: string) => {
    data.forEach((i: NetUsage) => {
      res.push(`${i.ethernet}_${item}`)
    })
  })
  return res
}
interface NetSeriesData {
  name: string
  data: number[]
  type: string
  showSymbol: boolean
  symbolSize: number
  hoverAnimation: boolean
  smooth: boolean
}
function generateNetSeriesData(data: NetUsage[]) {
  const options: string[] = ['Receive', 'Send']
  const series: NetSeriesData[] = []
  options.forEach((item: string) => {
    data.forEach((i: NetUsage) => {
      if (item === 'Receive') {
        series.push({
          name: `${i.ethernet}_${item}`,
          data: i.data.map((val: NetIO) => val.bytes_recv),
          type: 'line',
          smooth: true,
          showSymbol: true,
          symbolSize: 2,
          hoverAnimation: false,
        })
      }
      else {
        series.push({
          name: `${i.ethernet}_${item}`,
          data: i.data.map((val: NetIO) => val.bytes_sent),
          type: 'line',
          smooth: true,
          showSymbol: true,
          symbolSize: 2,
          hoverAnimation: false,
        })
      }
    })
  })
  return series
}
async function renderNet() {
  if (!shouldUpdate('net')) return
  const param: NetTrendingArgs = {
    start_time: dayjs().unix() - props.timeDensity,
    end_time: dayjs().unix(),
  }
  const { data } = await queryNetworkUsage(param)
  const netData = data.usage
  set(
    netOption,
    'xAxis.data',
    netData[0].data.map((item: NetIO) => dayjs(item.timestamp * 1000).format("HH:mm:ss")),
  )
  set(netOption, 'legend.data', generateNetLegendData(netData))
  set(netOption, 'series', generateNetSeriesData(netData))
  dataCache.net.data = netData
  dataCache.net.lastUpdate = Date.now()
}

// 分离定时器
const timers = ref<{ [key: string]: NodeJS.Timeout }>({})

onMounted(() => {
  // 初始化数据
  renderCPUPercent()
  renderCPU()
  renderMemInfo()
  renderMem()
  renderDiskInfo()
  renderDisk()
  renderNetInfo()
  renderNet()

  // 为每个指标设置独立的定时器
  timers.value.cpu = setInterval(() => {
    renderCPUPercent()
    renderCPU()
  }, dataCache.cpu.interval)

  timers.value.mem = setInterval(() => {
    renderMemInfo()
    renderMem()
  }, dataCache.mem.interval)

  timers.value.disk = setInterval(() => {
    renderDiskInfo()
    renderDisk()
  }, dataCache.disk.interval)

  timers.value.net = setInterval(() => {
    renderNetInfo()
    renderNet()
  }, dataCache.net.interval)
})

onUnmounted(() => {
  // 清除所有定时器
  Object.values(timers.value).forEach(timer => clearInterval(timer))
})
</script>

<template>
  <div class="monitor-container">
    <el-card shadow="hover" class="chart-card">
      <el-skeleton :loading="loading" animated>
        <div class="am-column-content">
          <div class="chart-container">
            <div class="chart-item">
              <div class="chart-info">{{ t('monitor.percent') }}： {{ cpuPercent }}</div>
              <echarts :option="cpuOption" />
            </div>
            <div class="chart-item">
              <div class="chart-info">{{ t('monitor.total') }}：{{ memInfo.total }} {{ t('monitor.used') }}：{{ memInfo.used }} {{ t('monitor.percent') }}： {{ memInfo.percent }}</div>
              <echarts :option="memOption" />
            </div>
            <div class="chart-item">
              <div class="chart-info" v-for="(item, index) in diskInfo" :key="index">
                {{ item.device }} {{ t('monitor.total') }}：{{ item.total }} {{ t('monitor.used') }}：{{ item.used }} {{ t('monitor.percent') }}：{{ item.percent }}
              </div>
              <echarts :option="diskOption" />
            </div>
            <div class="chart-item">
              <div class="chart-info" v-for="(item, index) in netInfo" :key="index">
                {{ item.ethernet }} {{ t('monitor.receive') }}：{{ item.read }} {{ t('monitor.send') }}：{{ item.write }}
              </div>
              <echarts :option="netOption" />
            </div>
          </div>
        </div>
      </el-skeleton>
    </el-card>
  </div>
</template>

<style scoped lang="scss">
.monitor-container {
  height: 100%;
  width: 100%;
}

.chart-container {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 24px;
  height: 100%;
}

.chart-item {
  flex: 1;
  min-width: calc(50% - 12px);
  min-height: 300px;
  position: relative;
  background: var(--el-bg-color-page);
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.chart-info {
  margin-bottom: 8px;
  color: var(--el-text-color-secondary);
  font-size: 14px;
  line-height: 1.5;
}

@include b(column-content) {
  height: auto;
  min-height: 680px;
  width: 100%;
  padding: 16px;
  border-radius: 4px;
}

.chart-card {
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

:deep(.el-skeleton) {
  width: 100%;
  height: 100%;
  padding: 20px;
}

:deep(.el-skeleton__item) {
  height: 300px;
  margin-bottom: 16px;
}
</style>
