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

// fakeOrderBucketStore mirrors fakeOrderSummaryStore's idiom for
// ports.OrderBucketStore (F01-A).
type fakeOrderBucketStore struct {
	counts         ports.OrderBucketCounts
	err            error
	installationID string
}

func (f *fakeOrderBucketStore) GetOrderBucketCounts(_ context.Context, installationID string) (ports.OrderBucketCounts, error) {
	f.installationID = installationID
	return f.counts, f.err
}

func TestOrderSummaryServiceBucketCountsDelegatesToStore(t *testing.T) {
	bucketStore := &fakeOrderBucketStore{counts: ports.OrderBucketCounts{Novo: 1, Faturar: 2, Enviar: 3, Enviado: 4}}
	service := NewSummaryServiceWithBuckets(&fakeOrderSummaryStore{}, bucketStore)

	got, err := service.BucketCounts(context.Background(), " installation-1 ")
	if err != nil {
		t.Fatalf("BucketCounts() error = %v", err)
	}
	if got != bucketStore.counts {
		t.Fatalf("BucketCounts() = %+v, want %+v", got, bucketStore.counts)
	}
	if bucketStore.installationID != "installation-1" {
		t.Fatalf("installation ID = %q, want installation-1", bucketStore.installationID)
	}
}

func TestOrderSummaryServiceBucketCountsMissingInstallation(t *testing.T) {
	service := NewSummaryServiceWithBuckets(&fakeOrderSummaryStore{}, &fakeOrderBucketStore{})

	_, err := service.BucketCounts(context.Background(), "   ")
	if !errors.Is(err, ErrSummaryInstallationRequired) {
		t.Fatalf("BucketCounts() error = %v, want ErrSummaryInstallationRequired", err)
	}
}

func TestOrderSummaryServiceBucketCountsStoreNotConfigured(t *testing.T) {
	service := NewSummaryService(&fakeOrderSummaryStore{})

	_, err := service.BucketCounts(context.Background(), "installation-1")
	if !errors.Is(err, ErrSummaryStoreNotConfigured) {
		t.Fatalf("BucketCounts() error = %v, want ErrSummaryStoreNotConfigured", err)
	}
}

func TestOrderSummaryServiceBucketCountsPropagatesStoreError(t *testing.T) {
	wantErr := errors.New("bucket source unavailable")
	service := NewSummaryServiceWithBuckets(&fakeOrderSummaryStore{}, &fakeOrderBucketStore{err: wantErr})

	got, err := service.BucketCounts(context.Background(), "installation-1")
	if !errors.Is(err, wantErr) {
		t.Fatalf("BucketCounts() error = %v, want %v", err, wantErr)
	}
	if got != (ports.OrderBucketCounts{}) {
		t.Fatalf("BucketCounts() = %+v, want zero result on error", got)
	}
}
