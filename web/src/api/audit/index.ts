/**
 * @Author     : Amu
 * @Date       : 2024/4/8 23:39
 * @Description:
 */
import { Pagination } from '@/interface/pagination.ts'
import request from '@/api'
import { AuditQueryResult } from '@/interface/audit.ts'

export function queryAudit(params: Pagination) {
    return request.get<AuditQueryResult>('/api/v1/audit/query', params)
}
