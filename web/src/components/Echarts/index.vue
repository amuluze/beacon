<template>
    <div class="am-box">
        <!-- 插槽 -->
        <slot></slot>
        <div ref="chartRef" :style="{ width: props.width, height: props.height }" />
    </div>
</template>

<script setup lang="ts">
import { useEcharts } from '@/hooks/useEcharts'
import { onMounted } from 'vue'
import { EChartsOption } from './echarts'

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

const chartRef = shallowRef<HTMLElement>()
const currentOptions = shallowRef<EChartsOption>(props.option)
const { setOptions, initCharts } = useEcharts(chartRef as Ref<HTMLDivElement>, currentOptions.value)

watch(
    () => currentOptions,
    (currentOptions) => {
        console.log('====================>', currentOptions)
        setOptions(currentOptions.value)
    },
    {
        deep: true
    }
)
onMounted(() => {
    initCharts()
})
</script>

<style scoped lang="scss">
@include b(box) {
    width: 100%;
    height: 90%;
}
</style>
