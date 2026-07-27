package connectors

import (
	"testing"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsports "marketplace-central/apps/server_core/internal/modules/connectors/ports"
	productlinksports "marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

func TestIdentityAnchorAdapterProjectsCompleteVocabulary(t *testing.T) {
	t.Parallel()

	service := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{
		ProviderCode: "provider",
		IdentityAnchors: []connectorsports.IdentityAnchor{
			connectorsports.IdentityAnchorSellerSKU,
			connectorsports.IdentityAnchorEAN,
			connectorsports.IdentityAnchorTitle,
		},
	}})
	adapter := NewIdentityAnchorAdapter(service)

	got, err := adapter.ProviderIdentityAnchors("provider")
	if err != nil {
		t.Fatalf("ProviderIdentityAnchors() error = %v", err)
	}

	want := []productlinksports.ProviderIdentityAnchor{
		{Anchor: "seller_sku", Supplied: true},
		{Anchor: "ean", Supplied: true},
		{Anchor: "title", Supplied: true},
		{Anchor: "marca", Supplied: false},
		{Anchor: "refforn", Supplied: false},
	}
	if len(got) != len(want) {
		t.Fatalf("len(ProviderIdentityAnchors()) = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ProviderIdentityAnchors()[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}
