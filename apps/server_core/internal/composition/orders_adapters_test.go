package composition

import (
	"context"
	"errors"
	"testing"
	"time"

	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	internalreadapp "marketplace-central/apps/server_core/internal/modules/internal_read/application"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

// stubCostReader is a minimal internal_read ports.Reader used to drive
// ordersCostReaderAdapter's honest branches (source-unavailable, other
// error, success passthrough) without booting a real Oracle/xlsx source.
type stubCostReader struct {
	cost internalreaddomain.CostAsOf
	err  error
}

func (s stubCostReader) FindProductsForLinking(context.Context, internalreadports.FindProductsInput) ([]internalreaddomain.ProductCandidate, error) {
	return nil, errors.New("not implemented")
}

func (s stubCostReader) GetSellableStock(context.Context, internalreadports.SellableStockInput) (internalreaddomain.SellableStock, error) {
	return internalreaddomain.SellableStock{}, errors.New("not implemented")
}

func (s stubCostReader) GetCurrentPrice(context.Context, internalreadports.CurrentPriceInput) (internalreaddomain.CurrentPrice, error) {
	return internalreaddomain.CurrentPrice{}, errors.New("not implemented")
}

func (s stubCostReader) GetCostAsOf(context.Context, internalreadports.CostAsOfInput) (internalreaddomain.CostAsOf, error) {
	return s.cost, s.err
}

func (s stubCostReader) GetSalesHistory(context.Context, internalreadports.SalesHistoryInput) (internalreaddomain.SalesHistory, error) {
	return internalreaddomain.SalesHistory{}, errors.New("not implemented")
}

func (s stubCostReader) GetTaxInputs(context.Context, internalreadports.TaxInput) (internalreaddomain.TaxInputs, error) {
	return internalreaddomain.TaxInputs{}, errors.New("not implemented")
}

func TestOrdersCostReaderAdapter(t *testing.T) {
	t.Run("unavailable source degrades to honest nil, no fabricated cost", func(t *testing.T) {
		adapter := newOrdersCostReaderAdapter(internalreadapp.Service{}, false)
		got, err := adapter.GetCostAsOf(context.Background(), 42, time.Now())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Amount != nil {
			t.Fatalf("expected nil Amount when source unavailable, got %v", *got.Amount)
		}
	})

	t.Run("source-unavailable read error degrades to honest nil, no error surfaced", func(t *testing.T) {
		service := internalreadapp.NewService(stubCostReader{
			err: internalreaddomain.NewReadError(internalreaddomain.ReadErrorSourceUnavailable, "oracle not configured", nil),
		})
		adapter := newOrdersCostReaderAdapter(service, true)
		got, err := adapter.GetCostAsOf(context.Background(), 42, time.Now())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Amount != nil {
			t.Fatalf("expected nil Amount on source-unavailable, got %v", *got.Amount)
		}
	})

	t.Run("non-source-unavailable error propagates, never swallowed", func(t *testing.T) {
		wantErr := errors.New("boom")
		service := internalreadapp.NewService(stubCostReader{err: wantErr})
		adapter := newOrdersCostReaderAdapter(service, true)
		_, err := adapter.GetCostAsOf(context.Background(), 42, time.Now())
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected propagated error, got %v", err)
		}
	})

	t.Run("successful read passes CostAsOf through unmapped", func(t *testing.T) {
		amount := 12.5
		service := internalreadapp.NewService(stubCostReader{
			cost: internalreaddomain.CostAsOf{ProductID: 42, Amount: &amount},
		})
		adapter := newOrdersCostReaderAdapter(service, true)
		got, err := adapter.GetCostAsOf(context.Background(), 42, time.Now())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Amount == nil || *got.Amount != amount {
			t.Fatalf("expected Amount %v passed through, got %v", amount, got.Amount)
		}
	})
}

func TestOrdersShipmentReaderAdapter(t *testing.T) {
	t.Run("installation not found returns honest error, never a zeroed ShipmentInfo", func(t *testing.T) {
		installations := integrationsapp.NewInstallationService(fakeInstallationRepo{}, "tenant_default")
		adapter := newOrdersShipmentReaderAdapter(nil, installations, "tenant_default")
		_, err := adapter.GetShipment(context.Background(), "inst-missing", "ship-1")
		if err == nil {
			t.Fatal("expected an error for an unresolved installation, got nil")
		}
	})

	t.Run("installation lookup error propagates", func(t *testing.T) {
		wantErr := errors.New("repo down")
		installations := integrationsapp.NewInstallationService(erroringInstallationRepo{err: wantErr}, "tenant_default")
		adapter := newOrdersShipmentReaderAdapter(nil, installations, "tenant_default")
		_, err := adapter.GetShipment(context.Background(), "inst-1", "ship-1")
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected propagated error, got %v", err)
		}
	})
}

// erroringInstallationRepo drives InstallationService.Get's error path;
// same-package fakeInstallationRepo (market_adapters_test.go:121) only
// covers the not-found path since it never errors.
type erroringInstallationRepo struct {
	err error
}

func (e erroringInstallationRepo) CreateInstallation(context.Context, integrationsdomain.Installation) error {
	return errors.New("not implemented")
}

func (e erroringInstallationRepo) GetInstallation(context.Context, string) (integrationsdomain.Installation, bool, error) {
	return integrationsdomain.Installation{}, false, e.err
}

func (e erroringInstallationRepo) ListInstallations(context.Context) ([]integrationsdomain.Installation, error) {
	return nil, e.err
}

func (e erroringInstallationRepo) UpdateInstallationStatus(context.Context, string, integrationsdomain.InstallationStatus, integrationsdomain.HealthStatus) error {
	return errors.New("not implemented")
}

func (e erroringInstallationRepo) ApplyConnectionSnapshot(context.Context, string, integrationsdomain.ConnectionSnapshot, string) error {
	return errors.New("not implemented")
}
