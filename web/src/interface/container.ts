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
    cpu_percent: number
    memory_percent: number
    memory_usage: number
    memory_limit: number
}

export interface CreateContainerArgs {
    container_name: string
    image_name: string
    network_id: string
    network_mode: string
    network_name: string
    ports: string[]
    volumes: string[]
    environments: string[]
    labels: Map<string, string>
}

export interface CreateContainerResult {
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
    network_sebnet: string
    network_gateway: string
    labels: Map<string, string>
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
