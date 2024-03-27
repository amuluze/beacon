/**
 * @Author     : Amu
 * @Date       : 2024/03/27 10:27:52
 * @Description:
 */
import echarts, { EChartsOption } from '@/components/Echarts/echarts'

const useEcharts = (elRef: Ref<HTMLDivElement>, options: EChartsOption) => {
    const charts = shallowRef<echarts.ECharts>()

    const initCharts = () => {
        charts.value = echarts.init(elRef.value)
        setOptions(options)
    }

    const setOptions = (options: EChartsOption) => {
        charts.value && charts.value.setOption(options, true)
    }

    const echartsResize = () => {
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
