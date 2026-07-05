import type { Freshness } from '@/interface/host'

export function freshnessText(freshness?: Freshness | null) {
    if (!freshness) {
        return '未知'
    }
    if (freshness.collected_at === 0) {
        return '无数据'
    }
    if (freshness.stale || freshness.degraded) {
        return `降级 ${freshness.age_seconds}s`
    }
    return `实时 ${freshness.age_seconds}s`
}

export function freshnessTagType(freshness?: Freshness | null) {
    if (!freshness || freshness.collected_at === 0) {
        return 'info'
    }
    if (freshness.stale || freshness.degraded) {
        return 'warning'
    }
    return 'success'
}

export function worstFreshness(items: Array<Freshness | null | undefined>) {
    const present = items.filter(Boolean) as Freshness[]
    if (present.length === 0) {
        return null
    }
    const empty = present.find(item => item.collected_at === 0)
    if (empty) {
        return empty
    }
    const degraded = present.find(item => item.stale || item.degraded)
    if (degraded) {
        return degraded
    }
    return present[0]
}
