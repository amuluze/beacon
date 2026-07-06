/**
 * @Author     : Amu
 * @Date       : 2024/10/17 16:48:05
 * @Description:
 */
import request from '@/api'
import type { AlarmThresholdQueryResult, AlarmThresholdUpdateArgs } from '@/interface/alarm'

export async function queryAlarmThreshold() {
    return request.get<AlarmThresholdQueryResult>('/api/v1/alarm/alarm_query')
}

export async function updateAlarmThreshold(params: AlarmThresholdUpdateArgs) {
    return request.post('/api/v1/alarm/alarm_update', params)
}
