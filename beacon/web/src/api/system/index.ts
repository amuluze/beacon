import request from '@/api'
import type { SystemVersion } from '@/interface/system.ts'

// 查询当前 Beacon Server 的构建版本与最近一次更新检查状态。
export async function querySystemVersion() {
    return request.get<SystemVersion>('/api/v1/system/version', {})
}
