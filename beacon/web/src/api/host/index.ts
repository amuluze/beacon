/**
 * @Author     : Amu
 * @Date       : 2024/3/7 14:47
 * @Description:
 */

import request from '@/api'
export { queryAgentList } from '@/api/agent'
import type {
    CPUInfo,
    CPUTrending,
    CPUTrendingArgs,
    DiskInfoResult,
    DiskTrendingArgs,
    DiskUsageResult,
    HostInfo,
    InstallTokenResult,
    MemInfo,
    MemTrending,
    MemTrendingArgs,
    NetTrendingArgs,
    NetUsageResult,
} from '@/interface/host.ts'

export async function queryHostInfo(params?: object) {
    return request.get<HostInfo>('/api/v1/host/host_info', params || {})
}

export async function queryCPUInfo(params?: object) {
    return request.get<CPUInfo>('/api/v1/host/cpu_info', params || {})
}
export async function queryCPUUsage(param: CPUTrendingArgs) {
    return request.get<CPUTrending>('/api/v1/host/cpu_trending', param)
}

export async function queryMemInfo(params?: object) {
    return request.get<MemInfo>('/api/v1/host/mem_info', params || {})
}
export async function queryMemUsage(param: MemTrendingArgs) {
    return request.get<MemTrending>('/api/v1/host/mem_trending', param)
}

export async function queryDiskInfo(params?: object) {
    return request.get<DiskInfoResult>('/api/v1/host/disk_info', params || {})
}

export async function queryDiskUsage(param: DiskTrendingArgs) {
    return request.get<DiskUsageResult>('/api/v1/host/disk_trending', param)
}

export async function queryNetworkUsage(param: NetTrendingArgs) {
    return request.get<NetUsageResult>('/api/v1/host/net_trending', param)
}

export async function getInstallToken() {
    return request.post<InstallTokenResult>('/api/v1/host/get_install_token', {})
}
