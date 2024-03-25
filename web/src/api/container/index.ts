import request from '@/api'
import { ContainerQueryResult, ImageQueryResult } from '@/interface/container.ts'

export function queryContainers() {
    return request.get<ContainerQueryResult>('/api/v1/container/containers', {})
}

export function queryImages() {
    return request.get<ImageQueryResult>('/api/v1/container/images', {})
}

export function queryDockerInfo() {
    return request.get<any>('/api/v1/container/version', {})
}
