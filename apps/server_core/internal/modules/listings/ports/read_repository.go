package ports

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/listings/domain"
)

type ListingQuery struct {
	InstallationID string
	Filter         domain.ListingFilter
	Q              string
	Cursor         ListingCursor
	Limit          int
}

type ListingGroupQuery struct {
	InstallationID string
	Filter         domain.ListingFilter
	Q              string
	Cursor         GroupCursor
	Limit          int
}

type SummaryQuery struct{ InstallationID string }

type ListingRowPage struct {
	Items      []domain.ListingReadModel
	NextCursor *ListingCursor
	AsOf       time.Time
}

type ListingGroupRowPage struct {
	Groups     []domain.ListingGroup
	NextCursor *GroupCursor
	AsOf       time.Time
}

type ListingSummaryRow struct {
	Total                int
	Active               int
	Paused               int
	SyncError            int
	Stale                int
	Unlinked             int
	BelowMarginWorstCase *int
	MarginUnknown        *int
	AsOf                 time.Time
	Linked               []SummaryLinkedRow
}

type SummaryLinkedRow struct {
	CostID int64
	Price  *domain.Money
}

type ListingReadRepository interface {
	ListListingRows(context.Context, ListingQuery) (ListingRowPage, error)
	ListListingGroupRows(context.Context, ListingGroupQuery) (ListingGroupRowPage, error)
	GetListingRow(context.Context, domain.ListingKey) (domain.ListingReadModel, bool, error)
	GetListingsSummary(context.Context, SummaryQuery) (ListingSummaryRow, error)
	ListListingTimeline(context.Context, domain.ListingKey, int) ([]domain.TimelineEvent, error)
}
