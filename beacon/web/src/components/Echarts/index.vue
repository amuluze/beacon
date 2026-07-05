<script setup lang="ts">
import { useEcharts } from '@/hooks/useEcharts'
import type { EChartsOption } from './echarts'
import useStore from '@/store'

interface Props {
  width?: string
  height?: string
  option: EChartsOption
}

const props = withDefaults(defineProps<Props>(), {
  width: '100%',
  height: '100%',
  option: () => ({}),
})

const chartRef = shallowRef<HTMLDivElement>()
const currentOptions = shallowRef<EChartsOption>(props.option)

const { setOptions, initCharts, echartsResize } = useEcharts(chartRef as Ref<HTMLDivElement>, currentOptions.value)

const store = useStore()

watch(
  () => props.option,
  (newVal) => {
    console.log('echarts new val: ', newVal)
    let targetOptions: EChartsOption = {}
    targetOptions = { ...newVal }
    setOptions(targetOptions)
  },
  {
    deep: true,
  },
)

watch(
  () => store.echarts.currentColorArray,
  (newValue) => {
    currentOptions.value.color = newValue
    setOptions(currentOptions.value)
  },
)

watch(
  () => store.app.isCollapse,
  () => {
    echartsResize()
  },
)

onMounted(() => {
  initCharts()
})
</script>

<template>
    <div class="am-echarts">
        <slot />
        <div ref="chartRef" :style="{ height: props.height, width: props.width }" />
    </div>
</template>

<style lang="scss" scoped>
@include b(echarts) {
  height: 100%;
  width: 100%;
}
</style>
