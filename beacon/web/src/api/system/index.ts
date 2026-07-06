/**
 * @Author     : Amu
 * @Date       : 2024/06/25 23:56:50
 * @Description:
 */
import request from '@/api'
import type { GetDNSResult, GetSystemTimeResult, GetSystemTimezoneListResult, GetSystemTimezoneResult, SetDNSArgs, SetSystemTimeArgs, SetSystemTimezoneArgs } from '@/interface/system'

export async function reboot() {
    return request.post('/api/v1/host/reboot', {})
}

export async function shutdown() {
    return request.post('/api/v1/host/shutdown', {})
}

export async function getDNS() {
    return request.get<GetDNSResult>('/api/v1/host/get_dns_settings', {})
}

export async function setDNS(params: SetDNSArgs) {
    return request.post('/api/v1/host/set_dns_settings', params)
}

export async function getSystemTime() {
    return request.get<GetSystemTimeResult>('/api/v1/host/get_system_time', {})
}

export async function setSystemTime(params: SetSystemTimeArgs) {
    return request.post('/api/v1/host/set_system_time', params)
}

export async function getSystemTimezoneList() {
    return request.get<GetSystemTimezoneListResult>('/api/v1/host/get_system_timezone_list', {})
}

export async function getSystemTimezone() {
    return request.get<GetSystemTimezoneResult>('/api/v1/host/get_system_timezone', {})
}

export async function setSystemTimezone(params: SetSystemTimezoneArgs) {
    return request.post('/api/v1/host/set_system_timezone', params)
}
