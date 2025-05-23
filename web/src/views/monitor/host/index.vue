<script setup lang="ts">
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import type { CPUTrendingArgs, DiskIO, DiskTrendingArgs, DiskUsage, MemTrendingArgs, NetIO, NetTrendingArgs, NetUsage } from '@/interface/host.ts'
import { useRoute, useRouter } from 'vue-router'
const router = useRouter()
const route = useRoute()

const activeTab = computed(() => {
  return route.path.startsWith('/monitor/host') ? 'host' : 'container'
})

import { queryCPUInfo, queryCPUUsage, queryDiskInfo, queryDiskUsage, queryMemInfo, queryMemUsage, queryNetworkUsage } from '@/api/host'
import { cpuTrendingOption, diskTrendingOption, memTrendingOption, netTrendingOption } from '@/config/echarts.ts'
import { convertBytesToReadable } from '@/utils/convert.ts'
import { dayjs } from 'element-plus'
import { set } from 'lodash-es'
import { useI18n } from 'vue-i18n'

const loading = ref(false)

// 时间密度下拉框
const timeDensity = ref(600)
const options = [
  {
    value: 120,
    label: '2分钟',
  },
  {
    value: 300,
    label: '5分钟',
  },
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
    renderCPUPercent()
    renderCPU()
    renderMemInfo()
    renderMem()
    renderDiskInfo()
    renderDisk()
    renderNet()
  },
)

// CPU 使用率
const cpuPercent = ref('0.0%')
async function renderCPUPercent() {
  const { data } = await queryCPUInfo()
  cpuPercent.value = `${data.percent.toFixed(2)}%`
}

const cpuOption = reactive<EChartsOption>(cpuTrendingOption as EChartsOption)
async function renderCPU() {
  const param: CPUTrendingArgs = {
    start_time: dayjs().unix() - timeDensity.value,
    end_time: dayjs().unix(),
  }

  console.log(param)
  const { data } = await queryCPUUsage(param)
  const cpuData = data.data
  console.log('cpu response:', cpuData)
  // set(cpuOption, 'title', { text: 'CPU使用率' });
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
  console.log('cpu: ', cpuOption)
}

// 内存使用率
const memInfo = ref({
  percent: '0%',
  total: '0',
  used: '0',
})
async function renderMemInfo() {
  const { data } = await queryMemInfo()
  memInfo.value.percent = `${data.percent.toFixed(2)}%`
  memInfo.value.total = convertBytesToReadable(data.total)
  memInfo.value.used = convertBytesToReadable(data.used)
}
const memOption = reactive<EChartsOption>(memTrendingOption as EChartsOption)
async function renderMem() {
  const param: MemTrendingArgs = {
    start_time: dayjs().unix() - timeDensity.value,
    end_time: dayjs().unix(),
  }
  console.log(param)
  const { data } = await queryMemUsage(param)
  const memData = data.data
  console.log('mem response: ', memData)
  // set(memOption, 'title', { text: '内存使用率' });
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
  console.log('mem: ', memOption)
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
  const { data } = await queryDiskInfo()
  console.log(data.info)
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
          hoverAnimation: true,
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
          hoverAnimation: true,
        })
      }
    })
  })
  return series
}
async function renderDisk() {
  const param: DiskTrendingArgs = {
    start_time: dayjs().unix() - timeDensity.value,
    end_time: dayjs().unix(),
  }
  console.log(param)
  const { data } = await queryDiskUsage(param)
  const diskData = data
  console.log('********<<<<>>>>>>*****', diskData, diskData.usage)
  // set(diskOption, 'title', { text: '磁盘使用率' });
  set(
    diskOption,
    'xAxis.data',
    diskData.usage[0].data.map((item: DiskIO) => dayjs(item.timestamp * 1000).format("HH:mm:ss")),
  )
  set(diskOption, 'legend.data', generateDiskLegendData(diskData.usage))
  set(diskOption, 'series', generateDiskSeriesData(diskData.usage))
  console.log('disk options: ', diskOption)
}

// 流量使用率
const netInfo = ref<
  {
    ethernet: string
    read: string
    write: string
  }[]
>([])
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
          showSymbol: true,
          symbolSize: 2,
          hoverAnimation: true,
          smooth: true,
        })
      }
      else {
        series.push({
          name: `${i.ethernet}_${item}`,
          data: i.data.map((val: NetIO) => val.bytes_sent),
          type: 'line',
          showSymbol: true,
          symbolSize: 2,
          hoverAnimation: true,
          smooth: true,
        })
      }
    })
  })
  return series
}
async function renderNet() {
  const param: NetTrendingArgs = {
    start_time: dayjs().unix() - timeDensity.value,
    end_time: dayjs().unix(),
  }
  console.log(param)
  const { data } = await queryNetworkUsage(param)
  const netData = data
  netInfo.value = []
  netData.usage.map((item) => {
    netInfo.value.push({
      ethernet: item.ethernet,
      read: convertBytesToReadable(item.data[item.data.length - 1].bytes_recv),
      write: convertBytesToReadable(item.data[item.data.length - 1].bytes_sent),
    })
    return netInfo
  })

  console.log('net response: ', netData.usage)
  set(
    netOption,
    'xAxis.data',
    netData.usage[0].data.map((item: NetIO) => dayjs(item.timestamp * 1000).format("HH:mm:ss")),
  )
  set(netOption, 'legend.data', generateNetLegendData(netData.usage))
  set(netOption, 'series', generateNetSeriesData(netData.usage))
  console.log('net options: ', netOption)
}

// 定时器
const timer = ref()
onMounted(() => {
  console.log('mounted')
  renderCPUPercent()
  renderCPU()
  renderMemInfo()
  renderMem()
  renderDiskInfo()
  renderDisk()
  renderNet()
  timer.value = setInterval(() => {
    console.log('start interval')
    renderCPUPercent()
    renderCPU()
    renderMemInfo()
    renderMem()
    renderDiskInfo()
    renderDisk()
    renderNet()
  }, 5000)
  console.log('timer: ', timer.value)
})

onUnmounted(() => {
  clearInterval(timer.value)
})

const { t } = useI18n()
</script>

<template>
    <div class="am-header-container">
        <div class="am-button-group">
            <div class="tab-nav-container">
                <div class="tab-nav-item" :class="{ active: activeTab === 'host' }" @click="router.push('/monitor/host')">
                    主机监控
                </div>
                <div class="tab-nav-item" :class="{ active: activeTab === 'container' }" @click="router.push('/monitor/container')">
                    容器监控
                </div>
            </div>
        </div>

        <div class="am-density-selector">
            <div class="tab-nav-container">
                <div class="tab-nav-item density-item">
                    <el-select
                        v-model="timeDensity"
                        placeholder="Select"
                        size="default"
                        class="density-select"
                        :popper-append-to-body="false"
                        :popper-style="{ minWidth: '200px' }"
                        auto-width>
                        <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
                    </el-select>
                </div>
            </div>
        </div>
    </div>
    <div class="am-column">
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
<style lang="scss" scoped>
// 在部分添加
.chart-card {
  transition: all 0.3s;
  margin-bottom: 16px;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  :deep(.el-card__header) {
    padding: 12px 16px;
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  }

  :deep(.el-card__body) {
    padding: 16px;
  }
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.chart-info {
  margin-bottom: 8px;
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

@include b(column-content) {
  height: 320px;
  width: 100%;
  padding: 8px;
  border-radius: 4px;
}

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

  .tab-nav-container {
    display: flex;
    background-color: #f5f7fa;
    border-radius: 8px;
    padding: 4px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  .tab-nav-item {
    padding: 8px 16px;
    cursor: pointer;
    font-size: 14px;
    color: #606266;
    border-radius: 6px;
    transition: all 0.3s;
    position: relative;
    white-space: nowrap;

    &:hover:not(.active) {
      color: var(--el-color-primary);
      background-color: rgba(var(--el-color-primary-rgb), 0.05);
    }

    &.active {
      color: var(--el-color-primary);
      background-color: #fff;
      box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
      font-weight: 500;
    }
  }
}

@include b(monitor-switch) {
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

@include b(density-selector) {
  background: var(--el-bg-color-page);
  border-radius: 8px;
  padding: 4px;

  .tab-nav-container {
    display: flex;
    background-color: #f5f7fa;
    border-radius: 8px;
    padding: 4px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  .tab-nav-item.density-item {
    padding: 4px 8px;
    background-color: #fff;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
    border-radius: 6px;
  }

  .selector-wrapper {
    display: flex;
    align-items: center;
    gap: 8px;
    height: 100%;
  }
}

.density-select {
  width: auto;
  min-width: 160px;

  :deep(.el-input__wrapper) {
    box-shadow: none !important;
    padding: 0 8px;
  }

  :deep(.el-input__inner) {
    min-width: 160px;
    height: 32px;
    line-height: 32px;
  }
}

@include b(column) {
  height: 100%;
  width: 100%;
  font-size: 14px;
  overflow-y: auto;
}

@include b(column-content) {
  height: 260px;
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

    :deep(.el-card__body) {
        padding: 20px;
    }
}
</style>
