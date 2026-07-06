/**
 * Operation type aliases extracted from `api-generated.d.ts`.
 *
 * Why this file exists:
 *   The auto-generated types use deeply-nested paths like
 *     `paths['/api/v1/host/host_info']['get']['responses']['200']['content']['application/json']`
 *   which are hard to read at call sites. We re-export them as friendly
 *   aliases (`HostInfoOkResponse`, `ContainerListOkResponse`, ...) so the
 *   underlying openapi-typescript output stays the single source of truth.
 *
 * Add new aliases here as concrete schemas are folded into openapi.yml.
 */
import type { paths } from './api-generated.d'

// ----- Host info / info endpoints -----
export type HostInfoOkResponse = paths['/api/v1/host/host_info']['get']['responses']['200']['content']['application/json']
export type CpuInfoOkResponse = paths['/api/v1/host/cpu_info']['get']['responses']['200']['content']['application/json']
export type MemInfoOkResponse = paths['/api/v1/host/mem_info']['get']['responses']['200']['content']['application/json']
export type DiskInfoOkResponse = paths['/api/v1/host/disk_info']['get']['responses']['200']['content']['application/json']

// ----- Trending endpoints -----
export type CpuTrendingOkResponse = paths['/api/v1/host/cpu_trending']['get']['responses']['200']['content']['application/json']
export type MemTrendingOkResponse = paths['/api/v1/host/mem_trending']['get']['responses']['200']['content']['application/json']
export type DiskTrendingOkResponse = paths['/api/v1/host/disk_trending']['get']['responses']['200']['content']['application/json']
export type NetTrendingOkResponse = paths['/api/v1/host/net_trending']['get']['responses']['200']['content']['application/json']

// ----- Container list endpoints -----
export type ContainerListOkResponse = paths['/api/v1/container/containers']['get']['responses']['200']['content']['application/json']
export type ImageListOkResponse = paths['/api/v1/container/images']['get']['responses']['200']['content']['application/json']
export type NetworkListOkResponse = paths['/api/v1/container/networks']['get']['responses']['200']['content']['application/json']
export type ContainerUsageOkResponse = paths['/api/v1/container/usage']['get']['responses']['200']['content']['application/json']
export type DockerVersionOkResponse = paths['/api/v1/container/version']['get']['responses']['200']['content']['application/json']