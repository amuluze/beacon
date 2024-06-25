/**
 * @Author     : Amu
 * @Date       : 2024/3/7 14:47
 * @Description:
 */

import request from '@/api'
import {
    CPUInfo,
    CPUTrending,
    CPUTrendingArgs,
    DiskInfoResult,
    DiskTrendingArgs,
    DiskUsageResult,
    FilesSearchArgs,
    FilesSearchResult,
    HostInfo,
    MemInfo,
    MemTrending,
    MemTrendingArgs,
    NetTrendingArgs,
    NetUsageResult
} from '@/interface/host.ts'

export function queryHostInfo() {
    return request.get<HostInfo>('/api/v1/host/host_info', {})
}

export function queryCPUInfo() {
    return request.get<CPUInfo>('/api/v1/host/cpu_info', {})
}
export function queryCPUUsage(param: CPUTrendingArgs) {
    return request.get<CPUTrending>('/api/v1/host/cpu_trending', param)
}

export function queryMemInfo() {
    return request.get<MemInfo>('/api/v1/host/mem_info', {})
}
export function queryMemUsage(param: MemTrendingArgs) {
    return request.get<MemTrending>('/api/v1/host/mem_trending', param)
}

export function queryDiskInfo() {
    return request.get<DiskInfoResult>('/api/v1/host/disk_info', {})
}

export function queryDiskUsage(param: DiskTrendingArgs) {
    return request.get<DiskUsageResult>('/api/v1/host/disk_trending', param)
}

export function queryNetworkUsage(param: NetTrendingArgs) {
    return request.get<NetUsageResult>('/api/v1/host/net_trending', param)
}

export function queryFiles(params: FilesSearchArgs) {
    return request.get<FilesSearchResult>('/api/v1/host/file_search', params)
}
