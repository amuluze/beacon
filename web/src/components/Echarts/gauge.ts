import { EChartsOption } from './echarts'

export const cpuGaugeOptions: EChartsOption = {
    series: [
        {
            type: 'gauge',
            startAngle: 90,
            endAngle: -270,
            pointer: {
                show: false
            },
            progress: {
                show: true,
                overlap: false,
                roundCap: true,
                clip: false,
                itemStyle: {
                    color: '#FAC858'
                    // borderWidth: 1,
                    // borderColor: '#464646'
                }
            },
            axisLine: {
                lineStyle: {
                    width: 8
                }
            },
            splitLine: {
                show: false,
                distance: 0,
                length: 10
            },
            axisTick: {
                show: false
            },
            axisLabel: {
                show: false,
                distance: 50
            },
            data: [],
            title: {
                fontSize: 14
            },
            detail: {
                width: 50,
                height: 14,
                fontSize: 14,
                color: 'inherit',
                borderColor: 'inherit',
                borderRadius: 20,
                borderWidth: 1,
                formatter: '{value}%'
            }
        }
    ]
}

export const memGaugeOptions: EChartsOption = {
    series: [
        {
            type: 'gauge',
            startAngle: 90,
            endAngle: -270,
            pointer: {
                show: false
            },
            progress: {
                show: true,
                overlap: false,
                roundCap: true,
                clip: false,
                itemStyle: {
                    color: '#92CC76'
                    // borderWidth: 1,
                    // borderColor: '#464646'
                }
            },
            axisLine: {
                lineStyle: {
                    width: 8
                }
            },
            splitLine: {
                show: false,
                distance: 0,
                length: 10
            },
            axisTick: {
                show: false
            },
            axisLabel: {
                show: false,
                distance: 50
            },
            data: [],
            title: {
                fontSize: 14
            },
            detail: {
                width: 50,
                height: 14,
                fontSize: 14,
                color: 'inherit',
                borderColor: 'inherit',
                borderRadius: 20,
                borderWidth: 1,
                formatter: '{value}%'
            }
        }
    ]
}

export const diskGaugeOptions: EChartsOption = {
    series: [
        {
            type: 'gauge',
            startAngle: 90,
            endAngle: -270,
            pointer: {
                show: false
            },
            progress: {
                show: true,
                overlap: false,
                roundCap: true,
                clip: false,
                itemStyle: {
                    color: '#5470C6'
                    // borderWidth: 1,
                    // borderColor: '#464646'
                }
            },
            axisLine: {
                lineStyle: {
                    width: 8
                }
            },
            splitLine: {
                show: false,
                distance: 0,
                length: 10
            },
            axisTick: {
                show: false
            },
            axisLabel: {
                show: false,
                distance: 50
            },
            data: [],
            title: {
                fontSize: 14
            },
            detail: {
                width: 50,
                height: 14,
                fontSize: 14,
                color: 'inherit',
                borderColor: 'inherit',
                borderRadius: 20,
                borderWidth: 1,
                formatter: '{value}%'
            }
        }
    ]
}
