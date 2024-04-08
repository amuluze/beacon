/**
 * @Author     : Amu
 * @Date       : 2024/4/8 23:41
 * @Description:
 */

export interface Audit {
    id: number
    username: string
    operate: string
    created: string
}

export interface AuditQueryResult {
    data: Audit[]
    total: number
    page: number
    size: number
}
