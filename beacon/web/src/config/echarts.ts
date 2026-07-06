/**
 * @Author     : Amu
 * @Date       : 2024/11/26 12:00
 * @Description:
 */
import type { EChartsOption } from '@/components/Echarts/echarts.ts'

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
                color: ['#fff'],
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
                color: ['#fff'],
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
                color: ['#fff'],
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
        formatter(params: any) {
            params = params[0]
            return `${params.value}%`
        },
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
    xAxis: {
        type: 'category',
        boundaryGap: false,
        data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    },
    yAxis: [
        {
            type: 'value',
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
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
        },
    ],
} as EChartsOption

export const memTrendingOption = {
    tooltip: {
        trigger: 'axis',
        formatter(params: any) {
            params = params[0]
            return `${params.value}%`
        },
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
    xAxis: {
        type: 'category',
        boundaryGap: false,
        data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    },
    yAxis: [
        {
            type: 'value',
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
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
        },
    ],
} as EChartsOption

export const diskTrendingOption = {
    tooltip: {
        trigger: 'axis',
        formatter(params: any): string {
            let res = ''
            params.forEach((item: any) => {
                const units = ['bps', 'Kbps', 'Mbps', 'Gbps']
                let unitIndex = 0
                while (item.value >= 1024 && unitIndex < units.length - 1) {
                    item.value /= 1024
                    unitIndex++
                }
                res += `${item.seriesName}: ${item.value.toFixed(2)} ${units[unitIndex]}<br/>`
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
    xAxis: {
        type: 'category',
        boundaryGap: false,
        data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    },
    yAxis: [
        {
            type: 'value',
            axisLabel: {
                show: true,
                formatter: function yAxisLabelFormatter(value: number): string {
                    const units = ['bps', 'Kbps', 'Mbps', 'Gbps']
                    let unitIndex = 0
                    while (value >= 1024 && unitIndex < units.length - 1) {
                        value /= 1024
                        unitIndex++
                    }
                    return `${value.toFixed(2)} ${units[unitIndex]}`
                },
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
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
        },
        {
            name: 'Write',
            type: 'line',
            smooth: true,
            lineStyle: {
                width: 2,
            },
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
        },
    ],
} as EChartsOption

export const netTrendingOption = {
    tooltip: {
        trigger: 'axis',
        formatter(params: any): string {
            let res = ''
            params.forEach((item: any) => {
                const units = ['bps', 'Kbps', 'Mbps', 'Gbps']
                let unitIndex = 0
                while (item.value >= 1024 && unitIndex < units.length - 1) {
                    item.value /= 1024
                    unitIndex++
                }
                res += `${item.seriesName}: ${item.value.toFixed(2)} ${units[unitIndex]}<br/>`
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
    xAxis: {
        type: 'category',
        boundaryGap: false,
        data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    },
    yAxis: [
        {
            type: 'value',
            axisLabel: {
                show: true,
                formatter: function yAxisLabelFormatter(value: number): string {
                    const units = ['bps', 'Kbps', 'Mbps', 'Gbps']
                    let unitIndex = 0
                    while (value >= 1024 && unitIndex < units.length - 1) {
                        value /= 1024
                        unitIndex++
                    }
                    return `${value.toFixed(2)} ${units[unitIndex]}`
                },
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
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
        },
        {
            name: 'Send',
            type: 'line',
            smooth: true,
            lineStyle: {
                width: 2,
            },
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
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
        formatter(params: any) {
            params = params[0]
            return `${params.value}%`
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
    xAxis: {
        type: 'category',
        boundaryGap: false,
        data: [],
    },
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
                const units = ['B', 'KB', 'MB', 'GB']
                let unitIndex = 0
                while (item.value >= 1024 && unitIndex < units.length - 1) {
                    item.value /= 1024
                    unitIndex++
                }
                res += `${item.seriesName}: ${item.value.toFixed(2)} ${units[unitIndex]}<br/>`
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
    xAxis: {
        type: 'category',
        boundaryGap: false,
        data: [],
    },
    yAxis: {
        type: 'value',
        axisLabel: {
            show: true,
            formatter: function yAxisLabelFormatter(value: number): string {
                const units = ['B', 'KB', 'MB', 'GB']
                let unitIndex = 0
                while (value >= 1024 && unitIndex < units.length - 1) {
                    value /= 1024
                    unitIndex++
                }
                return `${value.toFixed(2)} ${units[unitIndex]}`
            },
        },
    },
    series: [],
} as EChartsOption
