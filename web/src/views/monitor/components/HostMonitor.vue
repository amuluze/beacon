<script setup lang="ts">
import { queryCPUInfo, queryCPUUsage, queryDiskInfo, queryDiskUsage, queryMemInfo, queryMemUsage, queryNetworkUsage } from '@/api/host'
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import { cpuTrendingOption, diskTrendingOption, memTrendingOption, netTrendingOption } from '@/config/echarts.ts'
import type { CPUTrendingArgs, DiskIO, DiskTrendingArgs, DiskUsage, MemTrendingArgs, NetIO, NetTrendingArgs, NetUsage } from '@/interface/host.ts'
import { convertBytesToReadable } from '@/utils/convert.ts'
import { dayjs } from 'element-plus'
import { set } from 'lodash-es'

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
  const param: DiskTrendingArgs = {
    start_time: dayjs().unix() - timeDensity.value,
    end_time: dayjs().unix(),
  }
  const { data } = await queryDiskUsage(param)
  const diskData = data.data
  set(
    diskOption,
    'xAxis.data',
    diskData[0].data.map(item => dayjs(item.timestamp * 1000).format("HH:mm:ss")),
  )
  set(diskOption, 'legend.data', generateDiskLegendData(diskData))
  set(diskOption, 'series', generateDiskSeriesData(diskData))
}

// 网络使用率
const netInfo = ref<NetUsage[]>([])
async function renderNetInfo() {
  const { data } = await queryNetworkUsage()
  netInfo.value = data.info
}
const netOption = reactive<EChartsOption>(netTrendingOption as EChartsOption)
function generateNetLegendData(data: NetUsage[]) {
  const options: string[] = ['Receive', 'Send']
  const res: string[] = []
  options.forEach((item: string) => {
    data.forEach((i: NetUsage) => {
      res.push(`${i.device}_${item}`)
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
          name: `${i.device}_${item}`,
          data: i.data.map((val: NetIO) => val.io_receive),
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
          data: i.data.map((val: NetIO) => val.io_send),
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
  const param: NetTrendingArgs = {
    start_time: dayjs().unix() - timeDensity.value,
    end_time: dayjs().unix(),
  }
  const { data } = await queryNetworkUsage(param)
  const netData = data.data
  set(
    netOption,
    'xAxis.data',
    netData[0].data.map(item => dayjs(item.timestamp * 1000).format("HH:mm:ss")),
  )
  set(netOption, 'legend.data', generateNetLegendData(netData))
  set(netOption, 'series', generateNetSeriesData(netData))
}

onMounted(() => {
  renderCPUPercent()
  renderCPU()
  renderMemInfo()
  renderMem()
  renderDiskInfo()
  renderDisk()
  renderNetInfo()
  renderNet()
})
</script>

<template>
  <div class="monitor-container">
    <div class="monitor-header">
      <el-select v-model="timeDensity" class="time-density">
        <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
    </div>
    <div class="monitor-content">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card class="monitor-card">
            <template #header>
              <div class="card-header">
                <span>CPU 使用率</span>
                <span class="percent">{{ cpuPercent }}</span>
              </div>
            </template>
            <Echarts :option="cpuOption" />
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="monitor-card">
            <template #header>
              <div class="card-header">
                <span>内存使用率</span>
                <div class="info">
                  <span>总内存: {{ memInfo.total }}</span>
                  <span>已用: {{ memInfo.used }}</span>
                  <span class="percent">{{ memInfo.percent }}</span>
                </div>
              </div>
            </template>
            <Echarts :option="memOption" />
          </el-card>
        </el-col>
      </el-row>
      <el-row :gutter="20" class="mt-4">
        <el-col :span="24">
          <el-card class="monitor-card">
            <template #header>
              <div class="card-header">
                <span>磁盘使用率</span>
              </div>
            </template>
            <el-table :data="diskInfo" style="width: 100%">
              <el-table-column prop="device" label="设备" />
              <el-table-column prop="total" label="总容量" />
              <el-table-column prop="used" label="已用" />
              <el-table-column prop="percent" label="使用率" />
            </el-table>
            <Echarts :option="diskOption" class="mt-4" />
          </el-card>
        </el-col>
      </el-row>
      <el-row :gutter="20" class="mt-4">
        <el-col :span="24">
          <el-card class="monitor-card">
            <template #header>
              <div class="card-header">
                <span>网络使用率</span>
              </div>
            </template>
            <el-table :data="netInfo" style="width: 100%">
              <el-table-column prop="device" label="设备" />
              <el-table-column prop="ip" label="IP地址" />
              <el-table-column prop="mac" label="MAC地址" />
            </el-table>
            <Echarts :option="netOption" class="mt-4" />
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<style scoped lang="scss">
.monitor-container {
  padding: 20px;

  .monitor-header {
    margin-bottom: 20px;

    .time-density {
      width: 120px;
    }
  }

  .monitor-content {
    .monitor-card {
      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .info {
          display: flex;
          gap: 16px;
        }

        .percent {
          font-weight: bold;
          color: var(--el-color-primary);
        }
      }
    }
  }

  .mt-4 {
    margin-top: 16px;
  }
}
</style>
