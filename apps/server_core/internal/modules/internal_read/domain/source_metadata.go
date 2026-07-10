package domain

import "time"

type SourceMetadata struct {
	System     string
	FetchedAt  time.Time
	ObservedAt *time.Time
}

type FreshnessPolicy struct {
	MaxAge time.Duration
}

func (p FreshnessPolicy) IsZero() bool {
	return p.MaxAge <= 0
}
