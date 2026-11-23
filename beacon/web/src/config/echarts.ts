/**
 * @Author     : Amu
 * @Date       : 2024/11/26 12:00
 * @Description:
 */
import type { EChartsOption } from '@/components/Echarts/echarts.ts'

const bytesPerSecondUnits = ['B/s', 'KB/s', 'MB/s', 'GB/s', 'TB/s']
const byteUnits = ['B', 'KB', 'MB', 'GB', 'TB']

function seriesMetricValue(value: unknown): number {
    if (Array.isArray(value))
        return Number(value[1] ?? 0)
    return Number(value)
}

function formatBytes(bytes: number): string {
    let value = bytes
    let unitIndex = 0
    while (Math.abs(value) >= 1024 && unitIndex < byteUnits.length - 1) {
        value /= 1024
        unitIndex++
    }
    return `${value.toFixed(2)} ${byteUnits[unitIndex]}`
}

function formatPercentageTooltip(params: any[]): string {
    return `${seriesMetricValue(params[0]?.value)}%`
}

function timeAxis() {
    return {
        type: 'time',
        boundaryGap: false,
    }
}

export function formatBytesPerSecond(bytesPerSecond: number): string {
    let value = bytesPerSecond
    let unitIndex = 0
    while (Math.abs(value) >= 1024 && unitIndex < bytesPerSecondUnits.length - 1) {
        value /= 1024
        unitIndex++
    }
    return `${value.toFixed(2)} ${bytesPerSecondUnits[unitIndex]}`
}

export const cpuTrendingOption = {
    tooltip: {
        trigger: 'axis',
        formatter: formatPercentageTooltip,
        axisPointer: {
            type: 'cross',
            label: {
                backgroundColor: '#6a7985',
            },
        },
    },
    legend: {
        data: [{ name: 'CPU 使用率' }],
        left: 'right',
    },
    grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
    },
    xAxis: timeAxis(),
    yAxis: [
        {
            type: 'value',
            min: 0,
            max: 100,
            interval: 20,
            axisLabel: {
                show: true,
                formatter: '{value} %',
            },
        },
    ],
    series: [
        {
            name: 'CPU 使用率',
            type: 'line',
            smooth: true,
            lineStyle: {
                width: 2,
            },
            data: [],
        },
    ],
} as EChartsOption

export const memTrendingOption = {
    tooltip: {
        trigger: 'axis',
        formatter: formatPercentageTooltip,
        axisPointer: {
            type: 'cross',
            label: {
                backgroundColor: '#6a7985',
            },
        },
    },
    legend: {
        data: ['内存使用率'],
        left: 'right',
    },
    grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
    },
    xAxis: timeAxis(),
    yAxis: [
        {
            type: 'value',
            min: 0,
            max: 100,
            interval: 20,
            axisLabel: {
                show: true,
                formatter: '{value} %',
            },
        },
    ],
    series: [
        {
            name: '内存使用率',
            type: 'line',
            smooth: true,
            lineStyle: {
                width: 2,
            },
            data: [],
        },
    ],
} as EChartsOption

export const diskTrendingOption = {
    tooltip: {
        trigger: 'axis',
        formatter(params: any): string {
            let res = ''
            params.forEach((item: any) => {
                res += `${item.seriesName}: ${formatBytesPerSecond(seriesMetricValue(item.value))}<br/>`
            })
            return res
        },
        axisPointer: {
            type: 'cross',
            label: {
                backgroundColor: '#6a7985',
            },
        },
    },
    legend: {
        data: ['读', '写'],
        left: 'right',
    },
    grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
    },
    xAxis: timeAxis(),
    yAxis: [
        {
            type: 'value',
            min: 0,
            axisLabel: {
                show: true,
                formatter: formatBytesPerSecond,
            },
        },
    ],
    series: [
        {
            name: 'Read',
            type: 'line',
            smooth: true,
            lineStyle: {
                width: 2,
            },
            data: [],
        },
        {
            name: 'Write',
            type: 'line',
            smooth: true,
            lineStyle: {
                width: 2,
            },
            data: [],
        },
    ],
} as EChartsOption

export const netTrendingOption = {
    tooltip: {
        trigger: 'axis',
        formatter(params: any): string {
            let res = ''
            params.forEach((item: any) => {
                res += `${item.seriesName}: ${formatBytesPerSecond(seriesMetricValue(item.value))}<br/>`
            })
            return res
        },
        axisPointer: {
            type: 'cross',
            label: {
                backgroundColor: '#6a7985',
            },
        },
    },
    legend: {
        data: ['接收', '发送'],
        left: 'right',
    },
    grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
    },
    xAxis: timeAxis(),
    yAxis: [
        {
            type: 'value',
            min: 0,
            axisLabel: {
                show: true,
                formatter: formatBytesPerSecond,
            },
        },
    ],
    series: [
        {
            name: 'Receive',
            type: 'line',
            smooth: true,
            lineStyle: {
                width: 2,
            },
            data: [],
        },
        {
            name: 'Send',
            type: 'line',
            smooth: true,
            lineStyle: {
                width: 2,
            },
            data: [],
        },
    ],
} as EChartsOption

export const containerCpuOption = {
    title: {
        text: 'CPU',
        textStyle: {
            fontSize: '15px',
        },
    },
    tooltip: {
        trigger: 'axis',
        formatter: formatPercentageTooltip,
        axisPointer: {
            type: 'cross',
            label: {
                backgroundColor: '#6a7985',
            },
        },
    },
    legend: {
        show: false,
        data: [],
    },
    grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
    },
    xAxis: timeAxis(),
    yAxis: {
        type: 'value',
    },
    series: [],
} as EChartsOption

export const containerMemOption = {
    title: {
        text: '内存',
        textStyle: {
            fontSize: '15px',
        },
    },
    tooltip: {
        trigger: 'axis',
        formatter(params: any): string {
            let res = ''
            params.forEach((item: any) => {
                res += `${item.seriesName}: ${formatBytes(seriesMetricValue(item.value))}<br/>`
            })
            return res
        },
        axisPointer: {
            type: 'cross',
            label: {
                backgroundColor: '#6a7985',
            },
        },
    },
    legend: {
        show: false,
        data: [],
    },
    grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
    },
    xAxis: timeAxis(),
    yAxis: {
        type: 'value',
        axisLabel: {
            show: true,
            formatter: formatBytes,
        },
    },
    series: [],
} as EChartsOption
