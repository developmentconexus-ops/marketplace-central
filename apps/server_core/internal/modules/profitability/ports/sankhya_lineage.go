package ports

import "context"

type SankhyaLineageState string

const (
	SankhyaLineageNone        SankhyaLineageState = "none"
	SankhyaLineagePartial     SankhyaLineageState = "partial"
	SankhyaLineageComplete    SankhyaLineageState = "complete"
	SankhyaLineageConflict    SankhyaLineageState = "conflict"
	SankhyaLineageUnavailable SankhyaLineageState = "unavailable"
)

type SankhyaTaxSource struct {
	DocumentID int64
	LineNumber int
}

type SankhyaLineage struct {
	State       SankhyaLineageState
	Descendants []SankhyaTaxSource
}

type SankhyaLineageReader interface {
	ResolveCurrentLineage(ctx context.Context, installationID, providerOrderID, mpcLineID string) (SankhyaLineage, error)
}
