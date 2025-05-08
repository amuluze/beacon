<script setup lang="ts">
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import type { CPUTrendingArgs, DiskIO, DiskTrendingArgs, DiskUsage, MemTrendingArgs, NetIO, NetTrendingArgs, NetUsage } from '@/interface/host.ts'

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
    cpuData.map(item => `${dayjs(item.timestamp * 1000).hour()}:${dayjs(item.timestamp * 1000).minute()}`),
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
    memData.map(item => `${dayjs(item.timestamp * 1000).hour()}:${dayjs(item.timestamp * 1000).minute()}`),
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
    diskData.usage[0].data.map((item: DiskIO) => `${dayjs(item.timestamp * 1000).hour()}:${dayjs(item.timestamp * 1000).minute()}`),
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
    netData.usage[0].data.map((item: NetIO) => `${dayjs(item.timestamp * 1000).hour()}:${dayjs(item.timestamp * 1000).minute()}`),
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
    <div class="am-density">
        <el-card shadow="never">
            <span>{{ t('monitor.timeDensity') }}：</span>
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
                        <div>{{ t('monitor.cpuPercent') }}</div>
                        <div>{{ t('monitor.percent') }}： {{ cpuPercent }}</div>
                        <div class="am-column-content">
                            <echarts :option="cpuOption" />
                        </div>
                    </el-skeleton>
                </el-card>
            </el-col>
            <el-col :lg="12" :md="12" :sm="12" :xs="24">
                <el-card shadow="never">
                    <el-skeleton :loading="loading" animated>
                        <div>{{ t('monitor.memPercent') }}</div>
                        <div>{{ t('monitor.total') }}：{{ memInfo.total }} {{ t('monitor.used') }}：{{ memInfo.used }} {{ t('monitor.percent') }}： {{ memInfo.percent }}</div>
                        <div class="am-column-content">
                            <echarts :option="memOption" />
                        </div>
                    </el-skeleton>
                </el-card>
            </el-col>
            <el-col :lg="12" :md="12" :sm="12" :xs="24">
                <el-card shadow="never">
                    <el-skeleton :loading="loading" animated>
                        <div>{{ t('monitor.diskPercent') }}</div>
                        <div v-for="(item, index) in diskInfo" :key="index">
                            {{ item.device }} {{ t('monitor.total') }}：{{ item.total }} {{ t('monitor.used') }}：{{ item.used }} {{ t('monitor.percent') }}：{{ item.percent }}
                        </div>
                        <div class="am-column-content">
                            <echarts :option="diskOption" />
                        </div>
                    </el-skeleton>
                </el-card>
            </el-col>
            <el-col :lg="12" :md="12" :sm="12" :xs="24">
                <el-card shadow="never">
                    <el-skeleton :loading="loading" animated>
                        <div>{{ t('monitor.netLine') }}</div>
                        <div v-for="(item, index) in netInfo" :key="index">
                            {{ item.ethernet }} {{ t('monitor.receive') }}：{{ item.read }} {{ t('monitor.send') }}：{{ item.write }}
                        </div>
                        <div class="am-column-content">
                            <echarts :option="netOption" />
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
  height: 260px;
  width: 100%;
}
</style>
