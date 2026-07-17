import type { EChartsOption } from '@/components/Echarts/echarts'
import { describe, expect, it } from 'vitest'
import {
    cpuTrendingOption,
    diskTrendingOption,
    formatBytesPerSecond,
    memTrendingOption,
    netTrendingOption,
} from './echarts'

function firstYAxis(option: EChartsOption): Record<string, any> {
    const yAxis = option.yAxis
    return (Array.isArray(yAxis) ? yAxis[0] : yAxis) as Record<string, any>
}

function axisFormatter(option: EChartsOption): (value: number) => string {
    return firstYAxis(option).axisLabel.formatter
}

describe('host monitoring chart axes', () => {
    it.each([
        ['CPU', cpuTrendingOption],
        ['memory', memTrendingOption],
    ])('keeps the %s percentage axis between 0 and 100', (_name, option) => {
        expect(firstYAxis(option)).toMatchObject({
            min: 0,
            max: 100,
        })
    })

    it('formats disk and network rates as bytes per second', () => {
        expect(formatBytesPerSecond(0)).toBe('0.00 B/s')
        expect(formatBytesPerSecond(1024)).toBe('1.00 KB/s')
        expect(formatBytesPerSecond(1024 ** 2)).toBe('1.00 MB/s')

        expect(axisFormatter(diskTrendingOption)(1024 ** 2)).toBe('1.00 MB/s')
        expect(axisFormatter(netTrendingOption)(1024 ** 3)).toBe('1.00 GB/s')
    })

    it.each([
        ['disk', diskTrendingOption],
        ['network', netTrendingOption],
    ])('keeps the %s rate axis anchored at zero', (_name, option) => {
        expect(firstYAxis(option).min).toBe(0)
    })

    it('does not mutate ECharts tooltip data while formatting rates', () => {
        const params = [{ seriesName: 'eth0_Receive', value: 1024 ** 2 }]
        const formatter = (netTrendingOption.tooltip as any).formatter as (items: typeof params) => string

        expect(formatter(params)).toContain('1.00 MB/s')
        expect(params[0].value).toBe(1024 ** 2)
    })
})
