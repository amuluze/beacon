package schema

import "time"

// DefaultFreshnessStaleAfter 是默认的数据新鲜度过期时间。
// 超过此时间的数据标记为 stale / degraded。
const DefaultFreshnessStaleAfter = 2 * time.Minute

// Freshness describes whether monitoring data is still recent enough to be
// shown as live state.
type Freshness struct {
	CollectedAt int64 `json:"collected_at"`
	AgeSeconds  int64 `json:"age_seconds"`
	Stale       bool  `json:"stale"`
	Degraded    bool  `json:"degraded"`
}

// ComputeFreshness 根据采集时间戳计算数据新鲜度。
// ts 为零值时返回 Stale+Degraded 标记。
func ComputeFreshness(ts time.Time) Freshness {
	if ts.IsZero() {
		return Freshness{Stale: true, Degraded: true}
	}
	age := time.Since(ts)
	if age < 0 {
		age = 0
	}
	stale := age > DefaultFreshnessStaleAfter
	return Freshness{
		CollectedAt: ts.Unix(),
		AgeSeconds:  int64(age.Seconds()),
		Stale:       stale,
		Degraded:    stale,
	}
}
