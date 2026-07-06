package schema

type Freshness struct {
	CollectedAt int64 `json:"collected_at"`
	AgeSeconds  int64 `json:"age_seconds"`
	Stale       bool  `json:"stale"`
	Degraded    bool  `json:"degraded"`
}
