import request from '@/api'
import { ContainerQueryResult, ImageQueryResult } from '@/interface/container.ts'
import { Pagination } from '@/interface/pagination'

export function queryContainers(params: Pagination) {
    return request.get<ContainerQueryResult>('/api/v1/container/containers', params)
}

export function queryImages(params: Pagination) {
    return request.get<ImageQueryResult>('/api/v1/container/images', params)
}

export function queryDockerInfo() {
    return request.get<any>('/api/v1/container/version', {})
}
