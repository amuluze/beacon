/**
 * Re-exports of the auto-generated OpenAPI types with project-specific aliases
 * so call-sites get clean names like:
 *
 *   import type { paths, operations } from '@/types/api'
 *
 * under the hood they read from `api-generated.d.ts` (produced by
 * `pnpm run generate:types` from `.docs/api/openapi.yml`).
 *
 * `HealthResponse` is also surfaced as a typed alias for callers that want
 * to assert on the liveness/readiness payload.
 */
import type { components } from './api-generated.d'

export type {
    paths,
    operations,
    components,
    webhooks,
    $defs,
} from './api-generated.d'

export type HealthResponse = components['schemas']['HealthResponse']
