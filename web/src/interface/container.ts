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

export interface Image {
    id: string
    name: string
    tag: string
    created: string
    size: string
}

export interface ImageQueryResult {
    images: Image[]
}

export interface DockerInfo {
    docker_version: string
    api_version: string
    min_api_version: string
    git_commit: string
    go_version: string
    os: string
    arch: string
}
