<script setup lang="ts">
import { useEcharts } from '@/hooks/useEcharts'
import type { EChartsOption } from './echarts'

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

const { setOptions, initCharts } = useEcharts(chartRef as Ref<HTMLDivElement>, currentOptions.value)

watch(
  () => props.option,
  (newVal) => {
    let targetOptions: EChartsOption = {}
    targetOptions = { ...newVal }
    setOptions(targetOptions, { replaceMerge: ['series'] })
  },
  {
    deep: true,
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
