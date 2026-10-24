/**
 * @Author     : Amu
 * @Date       : 2024/10/17 16:48:05
 * @Description:
 */
import request from '@/api'
import type {
    AlarmMessageSetting,
    AlarmMessageSettingUpdateArgs,
    AlarmThresholdQueryResult,
    AlarmThresholdUpdateArgs,
} from '@/interface/alarm'

export async function queryAlarmThreshold() {
    return request.get<AlarmThresholdQueryResult>('/api/v1/alarm/alarm_query')
}

export async function updateAlarmThreshold(params: AlarmThresholdUpdateArgs) {
    return request.post('/api/v1/alarm/alarm_update', params)
}

export async function queryAlarmMessageSetting() {
    return request.get<AlarmMessageSetting>('/api/v1/alarm/message_setting')
}

export async function updateAlarmMessageSetting(params: AlarmMessageSettingUpdateArgs) {
    return request.post('/api/v1/alarm/message_setting_update', params)
}
