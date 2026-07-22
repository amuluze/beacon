/**
 * @Author     : Amu
 * @Date       : 2024/3/7 14:48
 * @Description:
 */

export interface InstallTokenResult {
    token: string
}

export interface HostInfo {
    timestamp: number
    stale?: boolean
    uptime: string
    hostname: string
    os: string
    platform: string
    platform_version: string
    kernel_version: string
    kernel_arch: string
    freshness: Freshness
}

export interface Freshness {
    collected_at: number
    age_seconds: number
    stale: boolean
    degraded: boolean
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
    timestamp: number
    stale?: boolean
    percent: number
    freshness: Freshness
}

export interface MemInfo {
    timestamp: number
    stale?: boolean
    percent: number
    total: number
    used: number
    freshness: Freshness
}
export interface DiskInfo {
    timestamp: number
    stale?: boolean
    device: string
    percent: number
    total: number
    used: number
}
export interface DiskInfoResult {
    info: DiskInfo[]
    freshness: Freshness
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
