import type { EChartsOption } from '@/components/Echarts/echarts'
import { describe, expect, it } from 'vitest'
import {
    containerCpuOption,
    containerMemOption,
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

function xAxis(option: EChartsOption): Record<string, any> {
    const axis = option.xAxis
    return (Array.isArray(axis) ? axis[0] : axis) as Record<string, any>
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
        const params = [{ seriesName: 'eth0_Receive', value: [1_000_000, 1024 ** 2] }]
        const formatter = (netTrendingOption.tooltip as any).formatter as (items: typeof params) => string

        expect(formatter(params)).toContain('1.00 MB/s')
        expect(params[0].value).toEqual([1_000_000, 1024 ** 2])
    })

    it.each([
        ['host CPU', cpuTrendingOption],
        ['host memory', memTrendingOption],
        ['host disk', diskTrendingOption],
        ['host network', netTrendingOption],
        ['container CPU', containerCpuOption],
        ['container memory', containerMemOption],
    ])('uses a real time axis for %s', (_name, option) => {
        expect(xAxis(option)).toMatchObject({
            type: 'time',
            boundaryGap: false,
        })
    })

    it('formats tuple series values in percentage and byte tooltips', () => {
        const cpuFormatter = (cpuTrendingOption.tooltip as any).formatter as (items: any[]) => string
        const memFormatter = (containerMemOption.tooltip as any).formatter as (items: any[]) => string

        expect(cpuFormatter([{ value: [1_000_000, 12.3] }])).toBe('12.3%')
        expect(memFormatter([{ seriesName: 'app', value: [1_000_000, 1024 ** 2] }])).toContain('app: 1.00 MB')
    })
})
