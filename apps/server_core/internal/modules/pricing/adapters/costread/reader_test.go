package costread

import (
	"context"
	"errors"
	"testing"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type fakeFactReader struct {
	cost internalreaddomain.CostAsOf
	err  error

	gotProductID   int
	gotEffectiveAt time.Time
}

func (f *fakeFactReader) GetCostAsOf(_ context.Context, productID int, effectiveAt time.Time) (internalreaddomain.CostAsOf, error) {
	f.gotProductID = productID
	f.gotEffectiveAt = effectiveAt
	return f.cost, f.err
}

func floatptr(v float64) *float64 { return &v }

func TestCostForProduct_PresentAmountMapsToDecimalMoney(t *testing.T) {
	asOf := time.Date(2026, 7, 18, 0, 0, 0, 0, time.UTC)
	fake := &fakeFactReader{cost: internalreaddomain.CostAsOf{Amount: floatptr(12.5)}}

	money, err := NewReader(fake).CostForProduct(context.Background(), 90001, asOf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if money == nil {
		t.Fatalf("expected money, got nil")
	}
	if money.Amount != "12.50" {
		t.Fatalf("amount = %q, want %q", money.Amount, "12.50")
	}
	if money.Currency != "BRL" {
		t.Fatalf("currency = %q, want BRL", money.Currency)
	}
	if fake.gotProductID != 90001 {
		t.Fatalf("productID passed through = %d, want 90001", fake.gotProductID)
	}
	if !fake.gotEffectiveAt.Equal(asOf) {
		t.Fatalf("effectiveAt passed through = %v, want %v", fake.gotEffectiveAt, asOf)
	}
}

// ADR-17: an absent cost fact must propagate as nil Money, never as 0.00.
func TestCostForProduct_AbsentAmountReturnsNilNotZero(t *testing.T) {
	fake := &fakeFactReader{cost: internalreaddomain.CostAsOf{Amount: nil}}

	money, err := NewReader(fake).CostForProduct(context.Background(), 90001, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if money != nil {
		t.Fatalf("expected nil money for absent cost, got %+v", money)
	}
}

func TestCostForProduct_PropagatesReaderError(t *testing.T) {
	wantErr := errors.New("oracle unavailable")
	fake := &fakeFactReader{err: wantErr}

	money, err := NewReader(fake).CostForProduct(context.Background(), 90001, time.Now())
	if !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want %v", err, wantErr)
	}
	if money != nil {
		t.Fatalf("expected nil money on error, got %+v", money)
	}
}

func TestCostForProduct_RoundsToTwoDecimals(t *testing.T) {
	fake := &fakeFactReader{cost: internalreaddomain.CostAsOf{Amount: floatptr(9.005)}}

	money, err := NewReader(fake).CostForProduct(context.Background(), 1, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if money == nil || money.Amount != "9.00" && money.Amount != "9.01" {
		t.Fatalf("amount = %v, want 2dp fixed string", money)
	}
	if len(money.Amount) < 4 || money.Amount[len(money.Amount)-3] != '.' {
		t.Fatalf("amount = %q, want exactly 2 decimal places", money.Amount)
	}
}
