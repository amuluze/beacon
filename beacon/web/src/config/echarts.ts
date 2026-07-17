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

export const cpuOption = {
    title: {
        text: 'CPU 使用率',
        x: '50%',
        y: 30,
        textAlign: 'center',
        textStyle: {
            color: '#363535',
            fontSize: '16px',
            fontWeight: 'bold',
            textAlign: 'center',
        },
    },
    series: [{
        type: 'liquidFill',
        radius: '50%',
        center: ['50%', '65%'], // 分别是 x、y 轴的便宜
        data: [0.5],
        label: {
            normal: {
                color: '#045cc0',
                insideColor: '#045cc0',
                textStyle: {
                    fontSize: '24px',
                    fontWeight: 'bold',
                },
            },
        },
        color: [{
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [{
                offset: 1,
                color: ['#fbfcfe'],
            }, {
                offset: 0,
                color: ['#6a7feb'],
            }],
            global: false,
        }],
        backgroundStyle: {
            borderWidth: 1,
            color: 'transparent',
        },
        outline: {
            show: true,
            borderDistance: 8, // 内层白圈的宽度
            itemStyle: { // 最外层圈的颜色的宽度
                borderColor: '#6a7feb',
                borderWidth: 4,
            },
        },
    }],
} as EChartsOption

export const memOption = {
    title: {
        text: '内存使用率',
        x: '50%',
        y: 30,
        textAlign: 'center',
        textStyle: {
            color: '#363535',
            fontSize: '16px',
            fontWeight: 'bold',
            textAlign: 'center',
        },
    },
    series: [{
        type: 'liquidFill',
        radius: '50%',
        center: ['50%', '65%'], // 分别是 x、y 轴的便宜
        data: [0.5],
        label: {
            normal: {
                color: '#c06504',
                insideColor: '#c06504',
                textStyle: {
                    fontSize: '24px',
                    fontWeight: 'bold',
                },
            },
        },
        color: [{
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [{
                offset: 1,
                color: ['#fbfcfe'],
            }, {
                offset: 0,
                color: ['#c06504'],
            }],
            global: false,
        }],
        backgroundStyle: {
            borderWidth: 1,
            color: 'transparent',
        },
        outline: {
            show: true,
            borderDistance: 8, // 内层白圈的宽度
            itemStyle: { // 最外层圈的颜色的宽度
                borderColor: '#c06504',
                borderWidth: 4,
            },
        },
    }],
} as EChartsOption

export const diskOption = {
    title: {
        text: '磁盘使用率',
        x: '50%',
        y: 30,
        textAlign: 'center',
        textStyle: {
            color: '#363535',
            fontSize: '16px',
            fontWeight: 'bold',
            textAlign: 'center',
        },
    },
    series: [{
        type: 'liquidFill',
        radius: '50%',
        center: ['50%', '65%'], // 分别是 x、y 轴的便宜
        data: [0.5],
        label: {
            normal: {
                color: '#5f7906',
                insideColor: '#5f7906',
                textStyle: {
                    fontSize: '24px',
                    fontWeight: 'bold',
                },
            },
        },
        color: [{
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [{
                offset: 1,
                color: ['#fbfcfe'],
            }, {
                offset: 0,
                color: ['#5f7906'],
            }],
            global: false,
        }],
        backgroundStyle: {
            borderWidth: 1,
            color: 'transparent',
        },
        outline: {
            show: true,
            borderDistance: 8, // 内层白圈的宽度
            itemStyle: { // 最外层圈的颜色的宽度
                borderColor: '#5f7906',
                borderWidth: 4,
            },
        },
    }],
} as EChartsOption

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
        data: ['Read', 'Write'],
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
        data: ['Receive', 'Send'],
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
