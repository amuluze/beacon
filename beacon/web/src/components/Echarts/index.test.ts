import { shallowMount } from '@vue/test-utils'
import { nextTick, reactive } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { EChartsOption } from './echarts'
import Echarts from './index.vue'

const echartsMocks = vi.hoisted(() => ({
    setOptions: vi.fn(),
    initCharts: vi.fn(),
    echartsResize: vi.fn(),
}))

vi.mock('@/hooks/useEcharts', () => ({
    useEcharts: () => echartsMocks,
}))

vi.mock('@/store', () => ({
    default: () => reactive({
        echarts: { currentColorArray: [] },
        app: { isCollapse: false },
    }),
}))

beforeEach(() => {
    vi.clearAllMocks()
})

describe('echarts option updates', () => {
    it('replaces series so an empty monitoring window removes old curves', async () => {
        const option = reactive<EChartsOption>({
            xAxis: { type: 'time' },
            series: [{ type: 'line', data: [[1_000, 10]] }],
        })
        const wrapper = shallowMount(Echarts, {
            props: { option },
        })

        option.series = []
        await nextTick()

        expect(echartsMocks.setOptions).toHaveBeenLastCalledWith(
            expect.objectContaining({ series: [] }),
            { replaceMerge: ['series'] },
        )
        wrapper.unmount()
    })
})
