<template>
    <div class="am-box">
        <!-- 插槽 -->
        <slot></slot>
        <div ref="chartRef" :style="{ width: props.width, height: props.height }" />
    </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { EChartsOption } from './echarts'
import echarts from './echarts'

interface Props {
    option: EChartsOption
    width?: string
    height?: string
}

const props = withDefaults(defineProps<Props>(), {
    width: '100%',
    height: '100%',
    option: () => ({})
})

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts

const options = computed(() => {
    return props.option
})

watch(
    () => options.value,
    (options) => {
        console.log('>>>>>>><<<<<<<<', options)
        if (chart) {
            chart.setOption(options, true)
        }
    },
    {
        deep: true
    }
)

const resizeHandler = () => {
    chart?.resize()
}

const initChart = () => {
    chart = echarts.init(chartRef.value as HTMLDivElement)
    console.log('>>>>>', props)
    chart?.setOption(options.value as EChartsOption, true)
}

onMounted(() => {
    initChart()
    window.addEventListener('resize', resizeHandler)
})
onDeactivated(() => {
    console.log('deactivated')
    window.removeEventListener('resize', resizeHandler)
})
onBeforeUnmount(() => {
    console.log('before unmount')
    window.removeEventListener('resize', resizeHandler)
    chart?.dispose()
})
onActivated(() => {
    console.log('activated')
    window.addEventListener('resize', resizeHandler)
})
</script>

<style scoped lang="scss">
@include b(box) {
    width: 100%;
    height: 90%;
}
</style>
