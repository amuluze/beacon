import request from '@/api'
import {
    ContainerQueryResult,
    ImageQueryResult,
    RemoveContainerArgs,
    RemoveImageArgs,
    RestartContainerArgs,
    StartContainerArgs,
    StopContainerArgs
} from '@/interface/container.ts'
import { Pagination } from '@/interface/pagination'

export function queryContainers(params: Pagination) {
    return request.get<ContainerQueryResult>('/api/v1/container/containers', params)
}

export function startContainer(params: StartContainerArgs) {
    return request.post('/api/v1/container/container_start', params)
}

export function stopContainer(params: StopContainerArgs) {
    return request.post('/api/v1/container/container_stop', params)
}

export function removeContainer(params: RemoveContainerArgs) {
    return request.post('/api/v1/container/container_remove', params)
}

export function restartContainer(params: RestartContainerArgs) {
    return request.post('/api/v1/container/container_restart', params)
}

export function removeImage(params: RemoveImageArgs) {
    return request.post('/api/v1/container/image_remove', params)
}

export function pruneImages() {
    return request.post('/api/v1/container/images_prune', {})
}

export function queryImages(params: Pagination) {
    return request.get<ImageQueryResult>('/api/v1/container/images', params)
}

export function queryDockerInfo() {
    return request.get<any>('/api/v1/container/version', {})
}
