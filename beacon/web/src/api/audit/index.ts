import request from '@/api'
import type { AuditQueryParams, AuditQueryResult } from '@/interface/audit.ts'

export async function queryAudit(params: AuditQueryParams) {
    return request.get<AuditQueryResult>('/api/v1/audit/query', params)
}
