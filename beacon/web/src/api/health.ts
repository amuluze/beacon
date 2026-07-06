/**
 * Health probe endpoints — typed via the auto-generated OpenAPI surface.
 *
 * This module is the first consumer of `@/types/operations` to demonstrate
 * the end-to-end flow: openapi.yml -> pnpm generate:types -> consumer.
 */
import request from '@/api'
import type { HealthResponse } from '@/types/api'

export async function fetchLiveness(): Promise<HealthResponse> {
    const { data } = await request.get<HealthResponse>('/health')
    return data
}

export async function fetchReadiness(): Promise<HealthResponse> {
    const { data } = await request.get<HealthResponse>('/ready')
    return data
}