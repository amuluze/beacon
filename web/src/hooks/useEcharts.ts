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
    const isActive = ref(true)

    const setOptions = (options: EChartsOption): void => {
        charts.value && charts.value.setOption(options)
    }

    const initCharts = (): void => {
        charts.value = echarts.init(elRef.value, null, { renderer: 'canvas' })
        setOptions(options)
    }

    const echartsResize = (): void => {
        if (isActive.value && charts.value) {
            charts.value.resize()
        }
    }

    // 统一处理事件监听器的添加和移除
    const handleResize = () => {
        if (isActive.value) {
            window.addEventListener('resize', echartsResize)
        }
    }

    const removeResize = () => {
        window.removeEventListener('resize', echartsResize)
    }

    onMounted(() => {
        handleResize()
    })

    onActivated(() => {
        isActive.value = true
        handleResize()
    })

    onDeactivated(() => {
        isActive.value = false
        removeResize()
    })

    onBeforeUnmount(() => {
        isActive.value = false
        removeResize()
        charts.value && charts.value.dispose()
    })

    return { initCharts, setOptions, echartsResize }
}

export { useEcharts }
