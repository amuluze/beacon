/**
 * @Author     : Amu
 * @Date       : 2024/3/7 14:48
 * @Description:
 */

export interface HostInfo {
    timestamp: number
    uptime: string
    hostname: string
    os: string
    platform: string
    platform_version: string
    kernel_version: string
    kernel_arch: string
}

export interface Usage {
    timestamp: number
    value: number
}

export interface DiskIO {
    timestamp: number
    io_read: number
    io_write: number
}

export interface DiskUsage {
    device: string
    data: DiskIO[]
}

export interface DiskUsageResult {
    usage: DiskUsage[]
}

export interface NetIO {
    timestamp: number
    bytes_sent: number
    bytes_recv: number
}

export interface NetUsage {
    ethernet: string
    data: NetIO[]
}

export interface NetUsageResult {
    usage: NetUsage[]
}

export interface CPUInfo {
    percent: number
}

export interface MemInfo {
    percent: number
    total: number
    used: number
}
export interface DiskInfo {
    device: string
    percent: number
    total: number
    used: number
}
export interface DiskInfoResult {
    info: DiskInfo[]
}

export interface CPUTrendingArgs {
    start_time: number
    end_time: number
}

export interface CPUTrending {
    data: Usage[]
}

export interface MemTrendingArgs {
    start_time: number
    end_time: number
}

export interface MemTrending {
    data: Usage[]
}

export interface DiskTrendingArgs {
    start_time: number
    end_time: number
}

export interface NetTrendingArgs {
    start_time: number
    end_time: number
}

export interface FilesSearchArgs {
    path: string
}

export interface FileInfo {
    name: string
    size: number
    mode: string
    mod_time: number
    is_dir: boolean
}

export interface FilesSearchResult {
    files: FileInfo[]
}

export interface FileUploadArgs {
    prefix: string
}

export interface FileUploadResult {}

export interface FileDownloadArgs {
    filepath: string
}

export interface FileDownloadResult {}

export interface FileDeleteArgs {
    filepath: string
}

export interface FileDeleteResult {}

export interface FileCreateArgs {
    path: string
    file_name: string
}
export interface FileCreateResult {}

export interface FolderCreateArgs {
    path: string
    folder_name: string
}

export interface FolderCreateResult {}

export interface FolderDeleteArgs {
    path: string
}

export interface FolderDeleteResult {}
