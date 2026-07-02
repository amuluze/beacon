import request from '@/api'
import type {
    ContainerQueryResult,
    ContainerTrending,
    ContainerTrendingArgs,
    CreateContainerArgs,
    CreateContainerResult,
    GetDockerRegistryMirrorsResult,
    ImageQueryResult,
    NetworkCreateArgs,
    NetworkCreateResult,
    NetworkDeleteArgs,
    NetworkDeleteResult,
    NetworkQueryResult,
    PullImageArgs,
    PullImageResult,
    RemoveContainerArgs,
    RemoveImageArgs,
    RestartContainerArgs,
    SetDockerRegistryMirrorsArgs,
    StartContainerArgs,
    StopContainerArgs,
} from '@/interface/container.ts'
import type { Pagination } from '@/interface/pagination'

export async function queryContainers(params: Pagination) {
    return request.get<ContainerQueryResult>('/api/v1/container/containers', params)
}

export async function queryContainersUsage(param: ContainerTrendingArgs) {
    return request.get<ContainerTrending>('/api/v1/container/usage', param)
}

export async function createContainer(params: CreateContainerArgs) {
    return request.post<CreateContainerResult>('/api/v1/container/container_create', params)
}

export async function startContainer(params: StartContainerArgs) {
    return request.post('/api/v1/container/container_start', params)
}

export async function stopContainer(params: StopContainerArgs) {
    return request.post('/api/v1/container/container_stop', params)
}

export async function removeContainer(params: RemoveContainerArgs) {
    return request.post('/api/v1/container/container_remove', params)
}

export async function restartContainer(params: RestartContainerArgs) {
    return request.post('/api/v1/container/container_restart', params)
}

export async function removeImage(params: RemoveImageArgs) {
    return request.post('/api/v1/container/image_remove', params)
}

export async function pruneImages() {
    return request.post('/api/v1/container/images_prune', {})
}

export async function queryImages(params: Pagination) {
    return request.get<ImageQueryResult>('/api/v1/container/images', params)
}

export async function pullImage(params: PullImageArgs) {
    return request.post<PullImageResult>('/api/v1/container/image_pull', params)
}

export async function queryNetworks(params: Pagination) {
    return request.get<NetworkQueryResult>('/api/v1/container/networks', params)
}

export async function createNetwork(params: NetworkCreateArgs) {
    return request.post<NetworkCreateResult>('/api/v1/container/network_create', params)
}

export async function removeNetwork(params: NetworkDeleteArgs) {
    return request.post<NetworkDeleteResult>('/api/v1/container/network_delete', params)
}

export async function queryDockerInfo() {
    return request.get<any>('/api/v1/container/version', {})
}

export async function getDockerRegistryMirrors() {
    return request.get<GetDockerRegistryMirrorsResult>('/api/v1/container/get_docker_registry_mirrors', {})
}

export async function SetDockerRegistryMirrors(params: SetDockerRegistryMirrorsArgs) {
    return request.post('/api/v1/container/set_docker_registry_mirrors', params)
}
