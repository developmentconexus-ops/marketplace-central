package mutations

import (
	"context"
	"testing"
	"time"

	inventorydomain "marketplace-central/apps/server_core/internal/modules/inventory/domain"
	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	mutationsports "marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type fakeProtocolRepository struct {
	created  mutationsports.CreateProtocolInput
	replaced []mutationsports.ReplaceItemInput
	approved bool
}

func (f *fakeProtocolRepository) CreateProtocol(_ context.Context, input mutationsports.CreateProtocolInput) (mutationsports.Protocol, error) {
	f.created = input
	return mutationsports.Protocol{ProtocolID: "prot-1", InstallationID: input.InstallationID, Type: input.Type}, nil
}

func (f *fakeProtocolRepository) GetProtocol(context.Context, string) (mutationsports.Protocol, bool, error) {
	return mutationsports.Protocol{}, false, nil
}

func (f *fakeProtocolRepository) ReplaceItems(_ context.Context, _ string, items []mutationsports.ReplaceItemInput) ([]mutationsports.MutationItem, error) {
	f.replaced = items
	return nil, nil
}

func (f *fakeProtocolRepository) ApproveItems(context.Context, string, time.Time) error {
	f.approved = true
	return nil
}

func (f *fakeProtocolRepository) ClaimProtocol(context.Context, string) (mutationsports.ProtocolClaim, bool, error) {
	return nil, false, nil
}

func stockAction(variationID string) inventorydomain.StockAction {
	before := 3
	return inventorydomain.StockAction{
		ActionID: "act-1",
		ProviderRef: inventorydomain.ProviderStockRef{
			InstallationID:      "inst-1",
			ProviderItemID:      "MLB-1",
			ProviderVariationID: variationID,
		},
		BeforeQuantity:    &before,
		RequestedQuantity: 7,
		Operator:          inventorydomain.StockActionOperator{ActorID: "op-1"},
		CreatedAt:         time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC),
	}
}

func TestCreateStockCorrectionUsesCanonicalNoVariationSentinel(t *testing.T) {
	repo := &fakeProtocolRepository{}
	envelope := NewEnvelope(repo)

	protocolID, err := envelope.CreateStockCorrection(context.Background(), stockAction(""))
	if err != nil {
		t.Fatalf("CreateStockCorrection() error = %v", err)
	}
	if protocolID != "prot-1" || !repo.approved {
		t.Fatalf("protocolID = %q, approved = %v", protocolID, repo.approved)
	}
	if len(repo.replaced) != 1 {
		t.Fatalf("replaced items = %d, want 1", len(repo.replaced))
	}
	listingID := repo.replaced[0].ListingID
	if listingID != "inst-1~MLB-1~-" {
		t.Fatalf("listing id = %q, want inst-1~MLB-1~-", listingID)
	}
	parsed, err := listingsdomain.ParseListingID(listingID)
	if err != nil {
		t.Fatalf("canonical parser rejected envelope listing id %q: %v", listingID, err)
	}
	if parsed.VariationID != listingsdomain.NoVariationID {
		t.Fatalf("parsed variation = %q, want %q", parsed.VariationID, listingsdomain.NoVariationID)
	}
}

func TestCreateStockCorrectionPreservesRealVariationID(t *testing.T) {
	repo := &fakeProtocolRepository{}
	envelope := NewEnvelope(repo)

	if _, err := envelope.CreateStockCorrection(context.Background(), stockAction("VAR-9")); err != nil {
		t.Fatalf("CreateStockCorrection() error = %v", err)
	}
	listingID := repo.replaced[0].ListingID
	if listingID != "inst-1~MLB-1~VAR-9" {
		t.Fatalf("listing id = %q, want inst-1~MLB-1~VAR-9", listingID)
	}
	if _, err := listingsdomain.ParseListingID(listingID); err != nil {
		t.Fatalf("canonical parser rejected %q: %v", listingID, err)
	}
}
