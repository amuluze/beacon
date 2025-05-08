<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import { containerCpuOption, containerMemOption } from '@/config/echarts.ts'
import type { CPUTrendingArgs, Usage } from '@/interface/host.ts'
import { dayjs } from 'element-plus'
import { queryContainersUsage } from '@/api/container'
import { set } from 'lodash-es'

const loading = ref(false)

// 时间密度下拉框
const timeDensity = ref(600)
const options = [
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
    render()
  },
)

// CPU 使用率
const cpuOption = reactive<EChartsOption>(containerCpuOption)
const memOption = reactive<EChartsOption>(containerMemOption)
interface Series {
  name: string
  type: string
  smooth: boolean
  data: string[]
}
async function render() {
  const param: CPUTrendingArgs = {
    start_time: dayjs().unix() - timeDensity.value,
    end_time: dayjs().unix(),
  }
  const { data } = await queryContainersUsage(param)
  const cpuData = new Map<string, Usage[]>(Object.entries(data.cpu_usage))
  const memData = new Map<string, Usage[]>(Object.entries(data.mem_usage))
  set(cpuOption, 'legend.data', data.names)
  set(memOption, 'legend.data', data.names)

  const couIterator = cpuData.keys()
  const cpuDataFirstKey = couIterator.next().value
  const memIterator = memData.keys()
  const memDataFirstKey = memIterator.next().value
  set(
    cpuOption,
    'xAxis.data',
    cpuData.get(cpuDataFirstKey as string)?.map(item => `${dayjs(item.timestamp * 1000).hour()}:${dayjs(item.timestamp * 1000).minute()}`),
  )
  set(
    memOption,
    'xAxis.data',
    memData.get(memDataFirstKey as string)?.map(item => `${dayjs(item.timestamp * 1000).hour()}:${dayjs(item.timestamp * 1000).minute()}`),
  )
  const cpuSeries = reactive<Series[]>([])
  cpuData.forEach((value, key) => {
    cpuSeries.push({
      name: key,
      type: 'line',
      smooth: true,
      data: value.map(item => item.value.toFixed(2)),
    })
  })
  set(cpuOption, 'series', cpuSeries)
  const memSeries = reactive<Series[]>([])
  memData.forEach((value, key) => {
    memSeries.push({
      name: key,
      type: 'line',
      smooth: true,
      data: value.map(item => item.value.toFixed(2)),
    })
  })
  set(memOption, 'series', memSeries)
  console.log('cpuOption: ', cpuOption)
  console.log('memOption: ', memOption)
}

// 定时器
const timer = ref()
onMounted(() => {
  console.log('mounted')
  render()
  timer.value = setInterval(() => {
    console.log('start interval')
    render()
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
                        <div class="am-column-content">
                            <echarts :option="cpuOption" />
                        </div>
                    </el-skeleton>
                </el-card>
            </el-col>
            <el-col :lg="12" :md="12" :sm="12" :xs="24">
                <el-card shadow="never">
                    <el-skeleton :loading="loading" animated>
                        <div class="am-column-content">
                            <echarts :option="memOption" />
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
  height: 320px;
  width: 100%;
}
</style>
