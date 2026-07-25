package productlinks

import (
	"context"
	"testing"

	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type stubWorkflowStore struct {
	links []productlinksdomain.ProductLink
}

func (s stubWorkflowStore) GetProductLink(context.Context, productlinksdomain.ListingIdentity) (productlinksdomain.ProductLink, bool, error) {
	return productlinksdomain.ProductLink{}, false, nil
}

func (s stubWorkflowStore) ApplyProductLinkTransition(context.Context, productlinksdomain.ProductLinkTransition) error {
	return nil
}

func (s stubWorkflowStore) ListProductLinks(context.Context, string, int) ([]productlinksdomain.ProductLink, error) {
	return s.links, nil
}

func (s stubWorkflowStore) ListProductLinkAuditEntries(context.Context, string, int) ([]productlinksdomain.ProductLinkAuditEntry, error) {
	return nil, nil
}

func (s stubWorkflowStore) InsertBatch(context.Context, productlinksdomain.ProductLinkBatch) error {
	return nil
}

func (s stubWorkflowStore) GetAuditEntry(context.Context, string) (productlinksdomain.ProductLinkAuditEntry, bool, error) {
	return productlinksdomain.ProductLinkAuditEntry{}, false, nil
}

func (s stubWorkflowStore) LatestAuditForIdentity(context.Context, productlinksdomain.ListingIdentity) (productlinksdomain.ProductLinkAuditEntry, bool, error) {
	return productlinksdomain.ProductLinkAuditEntry{}, false, nil
}

func (s stubWorkflowStore) ListAuditByBatch(context.Context, string) ([]productlinksdomain.ProductLinkAuditEntry, error) {
	return nil, nil
}

func (s stubWorkflowStore) ListDecisionsForLink(context.Context, productlinksdomain.ListingIdentity) ([]productlinksdomain.ProductLinkDecision, error) {
	return nil, nil
}

type stubCandidateStore struct {
	candidates []productlinksdomain.LinkCandidate
}

func (s stubCandidateStore) ReplaceLinkCandidates(context.Context, string, []productlinksdomain.ListingIdentity, []productlinksdomain.LinkCandidate) error {
	return nil
}

func (s stubCandidateStore) ListLinkCandidates(context.Context, string, int) ([]productlinksdomain.LinkCandidate, error) {
	return s.candidates, nil
}

func (s stubCandidateStore) GetLinkCandidate(context.Context, string) (productlinksdomain.LinkCandidate, bool, error) {
	return productlinksdomain.LinkCandidate{}, false, nil
}

func TestResolveLinksExposesInternalProductOnlyForResolvedPersistedLink(t *testing.T) {
	productID := 20303
	reader := NewLinkReader(
		stubWorkflowStore{links: []productlinksdomain.ProductLink{
			{InstallationID: "inst-1", ProviderItemID: "resolved-link", State: productlinksdomain.ProductLinkStateResolved, InternalProductID: &productID},
			{InstallationID: "inst-1", ProviderItemID: "rejected-link", State: productlinksdomain.ProductLinkStateRejected, InternalProductID: &productID},
			{InstallationID: "inst-1", ProviderItemID: "conflict-link", State: productlinksdomain.ProductLinkStateConflict, InternalProductID: &productID},
			{InstallationID: "inst-1", ProviderItemID: "unresolved-link", State: productlinksdomain.ProductLinkStateUnresolved, InternalProductID: &productID},
		}},
		stubCandidateStore{candidates: []productlinksdomain.LinkCandidate{
			{InstallationID: "inst-1", ProviderItemID: "exact-ean-candidate", State: productlinksdomain.LinkCandidateStateExactEAN, InternalProductID: &productID},
			{InstallationID: "inst-1", ProviderItemID: "conflict-candidate", State: productlinksdomain.LinkCandidateStateConflict, InternalProductID: &productID},
			{InstallationID: "inst-1", ProviderItemID: "unresolved-candidate", State: productlinksdomain.LinkCandidateStateUnresolved, InternalProductID: &productID},
		}},
	)

	links, err := reader.ResolveLinks(context.Background(), "inst-1", nil)
	if err != nil {
		t.Fatalf("ResolveLinks() error = %v", err)
	}

	tests := []struct {
		providerItemID string
		wantQuality    ordersdomain.LinkQuality
		wantProductID  bool
	}{
		{providerItemID: "exact-ean-candidate", wantQuality: ordersdomain.LinkQualityUnresolved},
		{providerItemID: "conflict-candidate", wantQuality: ordersdomain.LinkQualityConflict},
		{providerItemID: "unresolved-candidate", wantQuality: ordersdomain.LinkQualityUnresolved},
		{providerItemID: "resolved-link", wantQuality: ordersdomain.LinkQualityResolved, wantProductID: true},
		{providerItemID: "rejected-link", wantQuality: ordersdomain.LinkQualityRejected},
		{providerItemID: "conflict-link", wantQuality: ordersdomain.LinkQualityConflict},
		{providerItemID: "unresolved-link", wantQuality: ordersdomain.LinkQualityUnresolved},
	}
	for _, tt := range tests {
		t.Run(tt.providerItemID, func(t *testing.T) {
			link := links[linkKey("inst-1", tt.providerItemID, "")]
			if link.Quality != tt.wantQuality {
				t.Fatalf("quality = %q, want %q", link.Quality, tt.wantQuality)
			}
			if tt.wantProductID {
				if link.InternalProductID == nil || *link.InternalProductID != productID {
					t.Fatalf("InternalProductID = %v, want %d", link.InternalProductID, productID)
				}
				return
			}
			if link.InternalProductID != nil {
				t.Fatalf("InternalProductID = %d, want nil for %s link", *link.InternalProductID, link.Quality)
			}
		})
	}
}
