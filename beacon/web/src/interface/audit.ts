export interface Audit {
    id: number
    username: string
    agent_id?: string
    operate: string
    created: string
}

export interface AuditQueryResult {
    data: Audit[]
    total: number
    page: number
    size: number
}

export interface AuditQueryParams {
    type?: string
    agent_id?: string
    page: number
    size: number
}
