/**
 * @Author     : Amu
 * @Date       : 2024/3/10 17:42
 * @Description:
 */
import { EChartsOption } from '@/components/Echarts/echarts.ts'

export const cpuOptions: EChartsOption = {
    // title: {
    //     text: 'CPU 使用率',
    // },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross',
            label: {
                backgroundColor: '#6a7985'
            }
        }
    },
    legend: {
        data: [{ name: 'CPU 使用率' }],
        left: 'right'
    },
    grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
    },
    xAxis: {
        type: 'category',
        boundaryGap: false,
        data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
    },
    yAxis: [
        {
            type: 'value',
            axisLabel: {
                show: true,
                formatter: '{value} %'
            }
        }
    ],
    series: [
        {
            name: 'CPU 使用率',
            type: 'line',
            stack: 'Total',
            areaStyle: {},
            emphasis: {
                focus: 'series'
            },
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
        }
    ]
}

export const memOptions: EChartsOption = {
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross',
            label: {
                backgroundColor: '#6a7985'
            }
        }
    },
    legend: {
        data: ['内存使用率'],
        left: 'right'
    },
    grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
    },
    xAxis: {
        type: 'category',
        boundaryGap: false,
        data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
    },
    yAxis: [
        {
            type: 'value',
            axisLabel: {
                show: true,
                formatter: '{value} %'
            }
        }
    ],
    series: [
        {
            name: '内存使用率',
            type: 'line',
            stack: 'Total',
            areaStyle: {},
            emphasis: {
                focus: 'series'
            },
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
        }
    ]
}

export const diskOptions: EChartsOption = {
    // title: {
    //     text: '磁盘使用率',
    // },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross',
            label: {
                backgroundColor: '#6a7985'
            }
        }
    },
    legend: {
        data: ['磁盘'],
        left: 'right'
    },
    grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
    },
    xAxis: {
        type: 'category',
        boundaryGap: false,
        data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
    },
    yAxis: [
        {
            type: 'value',
            axisLabel: {
                show: true,
                formatter: function yAxisLabelFormatter(value: number): string {
                    const units = ['B', 'KB', 'MB', 'GB', 'TB']
                    let unitIndex = 0
                    while (value >= 1024 && unitIndex < units.length - 1) {
                        value /= 1024
                        unitIndex++
                    }
                    return value.toFixed(2) + ' ' + units[unitIndex]
                }
            }
        }
    ],
    series: [
        {
            name: 'Read',
            type: 'line',
            stack: 'Total',
            areaStyle: {},
            emphasis: {
                focus: 'series'
            },
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
        },
        {
            name: 'Write',
            type: 'line',
            stack: 'Total',
            areaStyle: {},
            emphasis: {
                focus: 'series'
            },
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
        }
    ]
}

export const netOptions: EChartsOption = {
    // title: {
    //     text: '网络流量',
    // },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross',
            label: {
                backgroundColor: '#6a7985'
            }
        }
    },
    legend: {
        data: ['Receive', 'Send'],
        left: 'right'
    },
    grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
    },
    xAxis: {
        type: 'category',
        boundaryGap: false,
        data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
    },
    yAxis: [
        {
            type: 'value',
            axisLabel: {
                show: true,
                formatter: function yAxisLabelFormatter(value: number): string {
                    const units = ['B', 'KB', 'MB', 'GB', 'TB']
                    let unitIndex = 0
                    while (value >= 1024 && unitIndex < units.length - 1) {
                        value /= 1024
                        unitIndex++
                    }
                    return value.toFixed(2) + ' ' + units[unitIndex]
                }
            }
        }
    ],
    series: [
        {
            name: 'Receive',
            type: 'line',
            stack: 'Total',
            areaStyle: {},
            emphasis: {
                focus: 'series'
            },
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
        },
        {
            name: 'Send',
            type: 'line',
            stack: 'Total',
            areaStyle: {},
            emphasis: {
                focus: 'series'
            },
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
        }
    ]
}
