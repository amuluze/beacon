/**
 * @Author     : Amu
 * @Date       : 2024/4/8 23:39
 * @Description:
 */
import request from '@/api'
import type { AuditQueryResult } from '@/interface/audit.ts'
import type { Pagination } from '@/interface/pagination.ts'

export async function queryOperateAudit(params: Pagination) {
    return request.get<AuditQueryResult>('/api/v1/audit/query', params)
}

export async function querySystemAudit(params: Pagination) {
    return request.get<AuditQueryResult>('/api/v1/audit/query', params)
}
