// Beacon Server 版本与更新状态契约，与后端 /api/v1/system/version 返回结构对齐。
export interface SystemVersion {
    version: string
    update_available: boolean
    latest_version?: string
    release_notes?: string
    checked_at?: string
}
