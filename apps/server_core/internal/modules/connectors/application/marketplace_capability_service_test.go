package application

import (
	"context"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/connectors/ports"
)

type fakeListingReader struct {
	snapshot domain.ListingSnapshot
}

func (f fakeListingReader) ListListings(_ context.Context, _ domain.ListListingsInput) ([]domain.ListingSnapshot, error) {
	return []domain.ListingSnapshot{f.snapshot}, nil
}

func (f fakeListingReader) ReadListing(_ context.Context, _ domain.ProviderListingRef) (domain.ListingSnapshot, error) {
	return f.snapshot, nil
}

type fakeStockManager struct{}

type fakeAccountProber struct{}

func (fakeAccountProber) ProbeAccount(_ context.Context, ref domain.ProviderAccountRef) (domain.AccountSnapshot, error) {
	return domain.AccountSnapshot{
		ProviderCode:      ref.ProviderCode,
		ProviderAccountID: ref.ProviderAccountID,
		FetchedAt:         time.Unix(400, 0).UTC(),
	}, nil
}

type fakeFeeQuoteReader struct{}

func (fakeFeeQuoteReader) ReadFeeQuote(_ context.Context, input domain.FeeQuoteInput) (domain.FeeQuoteSnapshot, error) {
	return domain.FeeQuoteSnapshot{
		ProviderCode:  input.AccountRef.ProviderCode,
		ListingTypeID: input.ListingTypeID,
		PriceAmount:   input.PriceAmount,
		CurrencyID:    input.CurrencyID,
		FetchedAt:     time.Unix(500, 0).UTC(),
	}, nil
}

func (fakeStockManager) ReadStock(_ context.Context, ref domain.ProviderListingRef) (domain.StockSnapshot, error) {
	return domain.StockSnapshot{
		ProviderCode:        ref.AccountRef.ProviderCode,
		ProviderItemID:      ref.ProviderItemID,
		ProviderVariationID: ref.ProviderVariationID,
		AvailableQuantity:   4,
		ProviderStatus:      "active",
		FetchedAt:           time.Unix(300, 0).UTC(),
		Scope:               domain.StockScopeItem,
	}, nil
}

func (fakeStockManager) UpdateAvailableQuantity(_ context.Context, request domain.StockWriteRequest) (domain.StockWriteResult, error) {
	return domain.StockWriteResult{
		IDempotencyKey:      request.IdempotencyKey,
		Result:              domain.StockWriteResultApplied,
		ProviderCode:        request.ProviderCode,
		ProviderItemID:      request.ProviderItemID,
		ProviderVariationID: request.ProviderVariationID,
	}, nil
}

func TestMarketplaceCapabilityServiceUsesFakeListingReader(t *testing.T) {
	t.Parallel()

	now := time.Unix(200, 0).UTC()
	svc := NewMarketplaceCapabilityService([]ProviderCapabilitySet{
		{
			ProviderCode: "mercado_livre",
			Listings: fakeListingReader{
				snapshot: domain.ListingSnapshot{
					ProviderCode:      "mercado_livre",
					ProviderItemID:    "MLB123",
					ProviderStatus:    "active",
					Title:             "Produto teste",
					SellerSKU:         "SKU-001",
					AvailableQuantity: intPtr(3),
					FetchedAt:         now,
				},
			},
		},
	})

	reader, err := svc.ListingReader("mercado_livre")
	if err != nil {
		t.Fatalf("ListingReader() error = %v", err)
	}

	snapshots, err := reader.ListListings(context.Background(), domain.ListListingsInput{
		AccountRef: domain.ProviderAccountRef{
			TenantID:       "tenant-default",
			InstallationID: "inst_001",
			ProviderCode:   "mercado_livre",
		},
	})
	if err != nil {
		t.Fatalf("ListListings() error = %v", err)
	}

	if got, want := len(snapshots), 1; got != want {
		t.Fatalf("len(ListListings()) = %d, want %d", got, want)
	}
	if got, want := snapshots[0].ProviderItemID, "MLB123"; got != want {
		t.Fatalf("ProviderItemID = %q, want %q", got, want)
	}
	if got, want := snapshots[0].SellerSKU, "SKU-001"; got != want {
		t.Fatalf("SellerSKU = %q, want %q", got, want)
	}
}

func TestMarketplaceCapabilityServiceReturnsUnsupportedForMissingOperation(t *testing.T) {
	t.Parallel()

	svc := NewMarketplaceCapabilityService([]ProviderCapabilitySet{
		{
			ProviderCode: "mercado_livre",
			Listings:     fakeListingReader{},
		},
	})

	_, err := svc.StockWriter("mercado_livre")
	if err == nil {
		t.Fatal("StockWriter() error = nil, want unsupported error")
	}
	if got, want := domain.ErrorCodeOf(err), domain.ErrCodeProviderUnsupportedShape; got != want {
		t.Fatalf("ErrorCodeOf() = %q, want %q", got, want)
	}
}

func TestMarketplaceCapabilityServiceReturnsStockManagerWhenRegistered(t *testing.T) {
	t.Parallel()

	svc := NewMarketplaceCapabilityService([]ProviderCapabilitySet{
		{
			ProviderCode: "mercado_livre",
			StockReads:   fakeStockManager{},
		},
	})

	manager, err := svc.StockReader("mercado_livre")
	if err != nil {
		t.Fatalf("StockReader() error = %v", err)
	}

	snapshot, err := manager.ReadStock(context.Background(), domain.ProviderListingRef{
		AccountRef: domain.ProviderAccountRef{
			TenantID:       "tenant-default",
			InstallationID: "inst_001",
			ProviderCode:   "mercado_livre",
		},
		ProviderItemID: "MLB456",
	})
	if err != nil {
		t.Fatalf("ReadStock() error = %v", err)
	}

	if got, want := snapshot.ProviderItemID, "MLB456"; got != want {
		t.Fatalf("ProviderItemID = %q, want %q", got, want)
	}
	if got, want := snapshot.Scope, domain.StockScopeItem; got != want {
		t.Fatalf("Scope = %q, want %q", got, want)
	}
}

func TestMarketplaceCapabilityServiceAccountProber(t *testing.T) {
	t.Parallel()

	svc := NewMarketplaceCapabilityService([]ProviderCapabilitySet{{
		ProviderCode:  "mercado_livre",
		AccountProbes: fakeAccountProber{},
	}})

	got, err := svc.AccountProber("mercado_livre")
	if err != nil {
		t.Fatalf("AccountProber() error = %v", err)
	}
	if got == nil {
		t.Fatalf("expected account prober")
	}
}

func TestMarketplaceCapabilityServiceFeeQuoteReaderUnsupported(t *testing.T) {
	t.Parallel()

	svc := NewMarketplaceCapabilityService([]ProviderCapabilitySet{{ProviderCode: "mercado_livre"}})

	_, err := svc.FeeQuoteReader("mercado_livre")
	if err == nil {
		t.Fatalf("expected unsupported fee quote reader")
	}
}

func intPtr(value int) *int {
	return &value
}

var _ ports.ListingReader = fakeListingReader{}
var _ ports.AccountProber = fakeAccountProber{}
var _ ports.FeeQuoteReader = fakeFeeQuoteReader{}
var _ ports.StockReader = fakeStockManager{}
var _ ports.StockWriter = fakeStockManager{}
