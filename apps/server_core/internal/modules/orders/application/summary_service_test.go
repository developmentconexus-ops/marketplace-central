package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

type fakeOrderSummaryStore struct {
	summary        ports.OrderSummary
	err            error
	installationID string
	referenceTime  time.Time
}

func (f *fakeOrderSummaryStore) GetOrderSummary(_ context.Context, installationID string, referenceTime time.Time) (ports.OrderSummary, error) {
	f.installationID = installationID
	f.referenceTime = referenceTime
	return f.summary, f.err
}

func TestOrderSummaryServiceDelegatesToStore(t *testing.T) {
	referenceTime := time.Date(2026, 7, 16, 12, 30, 0, 0, time.FixedZone("BRT", -3*60*60))
	store := &fakeOrderSummaryStore{summary: ports.OrderSummary{Today: 3, SevenDays: 11}}
	service := NewSummaryService(store)

	got, err := service.Summary(context.Background(), " installation-1 ", referenceTime)
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if got != store.summary {
		t.Fatalf("Summary() = %+v, want %+v", got, store.summary)
	}
	if store.installationID != "installation-1" {
		t.Fatalf("installation ID = %q, want installation-1", store.installationID)
	}
	if !store.referenceTime.Equal(referenceTime) {
		t.Fatalf("reference time = %v, want %v", store.referenceTime, referenceTime)
	}
}

func TestOrderSummaryServicePropagatesStoreError(t *testing.T) {
	wantErr := errors.New("summary source unavailable")
	service := NewSummaryService(&fakeOrderSummaryStore{err: wantErr})

	got, err := service.Summary(context.Background(), "installation-1", time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC))
	if !errors.Is(err, wantErr) {
		t.Fatalf("Summary() error = %v, want %v", err, wantErr)
	}
	if got != (ports.OrderSummary{}) {
		t.Fatalf("Summary() = %+v, want zero result on error", got)
	}
}
