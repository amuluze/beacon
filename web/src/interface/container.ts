/**
 * @Author     : Amu
 * @Date       : 2024/3/7 13:38
 * @Description:
 */
export interface Container {
    id: string
    name: string
    image: string
    ip: string
    state: string
    uptime: string
    cpu_percent: number
    memory_percent: number
    memory_usage: number
    memory_limit: number
}

export interface ContainerQueryResult {
    containers: Container[]
}
