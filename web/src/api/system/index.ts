/**
 * @Author     : Amu
 * @Date       : 2024/06/25 23:56:50
 * @Description:
 */
import request from '@/api'
import {
    GetDNSResult,
    GetSystemTimeResult,
    GetSystemTimezoneListResult,
    GetSystemTimezoneResult,
    SetDNSArgs,
    SetSystemTimeArgs,
    SetSystemTimezoneArgs
} from '@/interface/system'

export function reboot() {
    return request.post('/api/v1/host/reboot', {})
}

export function shutdown() {
    return request.post('/api/v1/host/shutdown', {})
}

export function getDNS() {
    return request.get<GetDNSResult>('/api/v1/host/get_dns_settings', {})
}

export function setDNS(params: SetDNSArgs) {
    return request.post('/api/v1/host/set_dns_settings', params)
}

export function getSystemTime() {
    return request.get<GetSystemTimeResult>('/api/v1/host/get_system_time', {})
}

export function setSystemTime(params: SetSystemTimeArgs) {
    return request.post('/api/v1/host/set_system_time', params)
}

export function getSystemTimezoneList() {
    return request.get<GetSystemTimezoneListResult>('/api/v1/host/get_system_timezone_list', {})
}

export function getSystemTimezone() {
    return request.get<GetSystemTimezoneResult>('/api/v1/host/get_system_timezone', {})
}

export function setSystemTimezone(params: SetSystemTimezoneArgs) {
    return request.post('/api/v1/host/set_system_timezone', params)
}
