/**
 * @Author     : Amu
 * @Date       : 2024/06/26 00:04:34
 * @Description:
 */

export interface GetDNSResult {
    dns: string[]
}

export interface SetDNSArgs {
    dns: string[]
}

export interface GetSystemTimeResult {
    system_time: string
}

export interface SetSystemTimeArgs {
    system_time: string
}

export interface GetSystemTimezoneResult {
    system_timezone: string
}

export interface SetSystemTimezoneArgs {
    system_timezone: string
}
