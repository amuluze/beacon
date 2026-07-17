import type { EChartsOption } from '@/components/Echarts/echarts.ts'
import { set } from 'lodash-es'

export interface MonitorTimeRange {
    start_time: number
    end_time: number
    agent_id?: string
}

export type MonitorSeriesPoint = [timestamp: number, value: number]

export function applyMonitorTimeRange(option: EChartsOption, range: MonitorTimeRange): void {
    set(option, 'xAxis.min', range.start_time * 1000)
    set(option, 'xAxis.max', range.end_time * 1000)
}

export function toMonitorSeriesPoint(timestamp: number, value: number): MonitorSeriesPoint {
    return [timestamp * 1000, value]
}
