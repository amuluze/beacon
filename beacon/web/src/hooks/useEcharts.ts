/**
 * @Author     : Amu
 * @Date       : 2024/03/27 10:27:52
 * @Description:
 */

import type { EChartsOption } from '@/components/Echarts/echarts'
import echarts from '@/components/Echarts/echarts'

function useEcharts(elRef: Ref<HTMLDivElement>, options: EChartsOption): {
    initCharts: () => void
    setOptions: (options: EChartsOption) => void
    echartsResize: () => void
} {
    const charts = shallowRef<echarts.ECharts>()

    const setOptions = (options: EChartsOption): void => {
        charts.value && charts.value.setOption(options)
    }

    const initCharts = (): void => {
        charts.value = echarts.init(elRef.value, null, { renderer: 'canvas' })
        setOptions(options)
    }

    const echartsResize = (): void => {
        charts.value && charts.value.resize()
    }

    onMounted(() => {
        window.addEventListener('resize', echartsResize)
    })

    // 防止 echarts 页面 keep-alive 时，还在继续监听页面
    onDeactivated(() => {
        window.removeEventListener('resize', echartsResize)
    })

    onBeforeUnmount(() => {
        window.removeEventListener('resize', echartsResize)
        charts.value && charts.value.dispose()
    })

    onActivated(() => {
        window.addEventListener('resize', echartsResize)
    })
    return { initCharts, setOptions, echartsResize }
}

export { useEcharts }
