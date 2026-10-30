import type { Freshness, Usage } from '@/interface/host.ts'

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
    ports: string
    state: string
    uptime: string
    // Server 侧（beacon/service/schema/container.go）已格式化为可展示字符串，
    // 例如 "0.42 %"、"128.00 MB"，前端直接渲染，不再二次换算。
    cpu_percent: string
    memory_percent: string
    memory_usage: string
    memory_limit: string
    network?: string
    volumes?: string
    environments?: string
    labels?: Record<string, string>
}

export type RestartPolicy = 'no' | 'always' | 'unless-stopped' | 'on-failure'

export interface CreateContainerArgs {
    container_name: string
    image_name: string
    network_id: string
    network_mode: string
    network_name: string
    ports: string[]
    volumes: string[]
    environments: string[]
    labels: Record<string, string>
    restart_policy: RestartPolicy
}

export interface CreateContainerResult {
    container_id: string
}

export interface UpdateContainerArgs {
    container_id: string
    container_name: string
    image_name: string
    network_name: string
    ports: string[]
    volumes?: string[]
    environments?: string[]
    labels?: Record<string, string>
    restart_policy: RestartPolicy
}

export interface UpdateContainerResult {
    container_id: string
}

export interface StartContainerArgs {
    container_id: string
}
export interface StopContainerArgs {
    container_id: string
}
export interface RemoveContainerArgs {
    container_id: string
}
export interface RestartContainerArgs {
    container_id: string
}

export interface ContainerQueryResult {
    data: Container[]
    total: number
    page: number
    size: number
}

// 容器 / 镜像 / 网络列表的过滤参数，服务端按同一条件计算 total。
export interface ContainerListParams {
    page: number
    size: number
    name?: string
    state?: string
}

export interface ImageListParams {
    page: number
    size: number
    name?: string
}

export interface NetworkListParams {
    page: number
    size: number
    name?: string
    driver?: string
}

// 工作区概览：容器总数 / 运行中数量、镜像数量与总体积、网络总数。
export interface ContainerSummary {
    containers: number
    containers_running: number
    images: number
    image_bytes: number
    networks: number
}

export interface Image {
    id: string
    name: string
    tag: string
    created: string
    size: string
    number: number
}

export interface PullImageArgs {
    image_name: string
}

export interface PullImageResult {}

export interface RemoveImageArgs {
    image_id: string
}

export interface ImageQueryResult {
    data: Image[]
    total: number
    page: number
    size: number
}

export interface Network {
    id: string
    name: string
    driver: string
    created: string
    subnet: string
    gateway: string
    // labels: Map<string, string>
}

export interface NetworkQueryResult {
    data: Network[]
    total: number
    page: number
    size: number
}

export interface NetworkCreateArgs {
    name: string
    driver: string
    subnet: string
    gateway: string
    labels: Record<string, string>
}

export interface NetworkCreateResult {
    network_id: string
}

export interface NetworkDeleteArgs {
    network_id: string
}

export interface NetworkDeleteResult {}

export interface DockerInfo {
    docker_version: string
    api_version: string
    min_api_version: string
    git_commit: string
    go_version: string
    os: string
    arch: string
}

export interface GetDockerRegistryMirrorsResult {
    registry_mirrors: string[]
}

export interface SetDockerRegistryMirrorsArgs {
    registry_mirrors: string[]
}

export interface ContainerTrendingArgs {
    start_time: number
    end_time: number
}

export interface ContainerTrending {
    names: string[]
    cpu_usage: Map<string, Usage[]>
    mem_usage: Map<string, Usage[]>
    freshness: Freshness
}
