import type { EChartsOption } from '@/components/Echarts/echarts'
import { cloneDeep } from 'lodash-es'
import { describe, expect, it } from 'vitest'
import { cpuTrendingOption, diskTrendingOption } from '@/config/echarts'
import { applyMonitorTimeRange, toMonitorSeriesPoint } from './monitorChart'

function xAxisOf(option: EChartsOption): Record<string, any> {
    const axis = option.xAxis
    return (Array.isArray(axis) ? axis[0] : axis) as Record<string, any>
}

describe('applyMonitorTimeRange', () => {
    it('maps a second-based window onto millisecond time axis bounds', () => {
        const option = cloneDeep(cpuTrendingOption)

        applyMonitorTimeRange(option, { start_time: 2_700, end_time: 3_000, agent_id: 'agent-a' })

        expect(xAxisOf(option)).toMatchObject({
            type: 'time',
            boundaryGap: false,
            min: 2_700_000,
            max: 3_000_000,
        })
    })

    it('keeps the chart axis configuration and leaves the range untouched', () => {
        const option = cloneDeep(diskTrendingOption)
        const range = { start_time: 0, end_time: 60 }

        applyMonitorTimeRange(option, range)

        expect(xAxisOf(option).min).toBe(0)
        expect(xAxisOf(option).max).toBe(60_000)
        expect(range).toEqual({ start_time: 0, end_time: 60 })
    })

    it('overwrites the previous window when the density changes', () => {
        const option = cloneDeep(cpuTrendingOption)

        applyMonitorTimeRange(option, { start_time: 1_000, end_time: 1_600 })
        applyMonitorTimeRange(option, { start_time: 1_200, end_time: 3_000 })

        expect(xAxisOf(option).min).toBe(1_200_000)
        expect(xAxisOf(option).max).toBe(3_000_000)
    })
})

describe('toMonitorSeriesPoint', () => {
    it('converts a second timestamp into an ECharts time axis tuple', () => {
        expect(toMonitorSeriesPoint(900, 12.3)).toEqual([900_000, 12.3])
    })

    it('keeps zero and negative values without rounding', () => {
        expect(toMonitorSeriesPoint(1, 0)).toEqual([1_000, 0])
        expect(toMonitorSeriesPoint(1, -2.5)).toEqual([1_000, -2.5])
    })
})
