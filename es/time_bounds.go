package es

// TimeBounds represent transaction time bounds
type TimeBounds struct {
	MinTime int64 `json:"min_time"`
	MaxTime int64 `json:"max_time"`
}
