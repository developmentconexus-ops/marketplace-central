package connectors

import (
	"context"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	connectorsports "marketplace-central/apps/server_core/internal/modules/connectors/ports"
)

// recordingSnapshotObserver and recordingListingReader are shared test
// doubles for this package's tests (backfill_test.go's F-03 coverage of
// BackfillSource/MultigetHydrator). They originally backed the now-deleted
// page-based Source's own tests too; kept here since backfill_test.go
// depends on them.
type recordingSnapshotObserver struct {
	installationID string
	snapshots      []connectorsdomain.ListingSnapshot
	err            error
}

func (o *recordingSnapshotObserver) AbsorbProviderSnapshots(_ context.Context, installationID string, snapshots []connectorsdomain.ListingSnapshot) error {
	o.installationID = installationID
	o.snapshots = append(o.snapshots, snapshots...)
	return o.err
}

var _ SnapshotObserver = (*recordingSnapshotObserver)(nil)

type recordingListingReader struct {
	input     connectorsdomain.ListListingsInput
	snapshots []connectorsdomain.ListingSnapshot
}

func (r *recordingListingReader) ListListings(_ context.Context, input connectorsdomain.ListListingsInput) ([]connectorsdomain.ListingSnapshot, error) {
	r.input = input
	return r.snapshots, nil
}
func (r *recordingListingReader) ReadListing(context.Context, connectorsdomain.ProviderListingRef) (connectorsdomain.ListingSnapshot, error) {
	return connectorsdomain.ListingSnapshot{}, nil
}

var _ connectorsports.ListingReader = (*recordingListingReader)(nil)
