/**
 * @Author     : Amu
 * @Date       : 2024/3/7 14:47
 * @Description:
 */

import request from '@/api'
import type {
    CPUInfo,
    CPUTrending,
    CPUTrendingArgs,
    DiskInfoResult,
    DiskTrendingArgs,
    DiskUsageResult,
    FileCreateArgs,
    FileCreateResult,
    FileDeleteArgs,
    FileDeleteResult,
    FileDownloadResult,
    FilesSearchArgs,
    FilesSearchResult,
    FolderCreateArgs,
    FolderCreateResult,
    FolderDeleteArgs,
    FolderDeleteResult,
    HostInfo,
    MemInfo,
    MemTrending,
    MemTrendingArgs,
    NetTrendingArgs,
    NetUsageResult,
} from '@/interface/host.ts'

export async function queryHostInfo() {
    return request.get<HostInfo>('/api/v1/host/host_info', {})
}

export async function queryCPUInfo() {
    return request.get<CPUInfo>('/api/v1/host/cpu_info', {})
}
export async function queryCPUUsage(param: CPUTrendingArgs) {
    return request.get<CPUTrending>('/api/v1/host/cpu_trending', param)
}

export async function queryMemInfo() {
    return request.get<MemInfo>('/api/v1/host/mem_info', {})
}
export async function queryMemUsage(param: MemTrendingArgs) {
    return request.get<MemTrending>('/api/v1/host/mem_trending', param)
}

export async function queryDiskInfo() {
    return request.get<DiskInfoResult>('/api/v1/host/disk_info', {})
}

export async function queryDiskUsage(param: DiskTrendingArgs) {
    return request.get<DiskUsageResult>('/api/v1/host/disk_trending', param)
}

export async function queryNetworkUsage(param: NetTrendingArgs) {
    return request.get<NetUsageResult>('/api/v1/host/net_trending', param)
}

export async function queryFiles(params: FilesSearchArgs) {
    return request.get<FilesSearchResult>('/api/v1/host/file_search', params)
}

export async function createFile(params: FileCreateArgs) {
    return request.post<FileCreateResult>('/api/v1/host/file_create', params)
}

export async function deleteFile(params: FileDeleteArgs) {
    return request.post<FileDeleteResult>('/api/v1/host/file_delete', params)
}

export async function createFolder(params: FolderCreateArgs) {
    return request.post<FolderCreateResult>('/api/v1/host/folder_create', params)
}

export async function deleteFolder(params: FolderDeleteArgs) {
    return request.post<FolderDeleteResult>('/api/v1/host/folder_delete', params)
}

export async function downloadFile(params: FileDeleteArgs) {
    return request.post<FileDownloadResult>('/api/v1/host/file_download', params)
}
