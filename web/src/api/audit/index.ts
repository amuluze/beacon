/**
 * @Author     : Amu
 * @Date       : 2024/4/8 23:39
 * @Description:
 */
import request from '@/api'
import { AuditQueryResult } from '@/interface/audit.ts'
import { Pagination } from '@/interface/pagination.ts'

export function queryOperateAudit(params: Pagination) {
    return request.get<AuditQueryResult>('/api/v1/audit/query', params)
}

export function querySystemAudit(params: Pagination) {
    return request.get<AuditQueryResult>('/api/v1/audit/query', params)
}
