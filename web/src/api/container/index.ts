import request from '@/api'
import { ContainerQueryResult } from '@/interface/container.ts'

export function queryContainers() {
    return request.get<ContainerQueryResult>('/api/v1/container/containers', {})
}
