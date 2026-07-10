package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type ListingSnapshotReader interface {
	ListListingSnapshots(ctx context.Context, installationID string, limit int) ([]domain.ListingSnapshot, error)
}

type LinkCandidateStore interface {
	ReplaceLinkCandidates(ctx context.Context, installationID string, identities []domain.ListingIdentity, candidates []domain.LinkCandidate) error
	ListLinkCandidates(ctx context.Context, installationID string, limit int) ([]domain.LinkCandidate, error)
	GetLinkCandidate(ctx context.Context, candidateID string) (domain.LinkCandidate, bool, error)
}
