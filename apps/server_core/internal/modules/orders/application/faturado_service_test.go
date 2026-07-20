package application

import (
	"context"
	"errors"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

// fakeOrderFaturadoStore mirrors fakeOrderSummaryStore's idiom for
// ports.OrderFaturadoStore.
type fakeOrderFaturadoStore struct {
	err             error
	called          bool
	installationID  string
	providerOrderID string
}

func (f *fakeOrderFaturadoStore) MarkOrderFaturado(_ context.Context, installationID, providerOrderID string) error {
	f.called = true
	f.installationID = installationID
	f.providerOrderID = providerOrderID
	return f.err
}

func TestFaturadoServiceMarkFaturadoDelegatesToStore(t *testing.T) {
	store := &fakeOrderFaturadoStore{}
	service := NewFaturadoService(store)

	err := service.MarkFaturado(context.Background(), " installation-1 ", " order-1 ")
	if err != nil {
		t.Fatalf("MarkFaturado() error = %v", err)
	}
	if !store.called {
		t.Fatal("MarkFaturado() did not call the store")
	}
	if store.installationID != "installation-1" {
		t.Fatalf("installation ID = %q, want installation-1", store.installationID)
	}
	if store.providerOrderID != "order-1" {
		t.Fatalf("provider order ID = %q, want order-1", store.providerOrderID)
	}
}

func TestFaturadoServiceMarkFaturadoMissingInstallationID(t *testing.T) {
	store := &fakeOrderFaturadoStore{}
	service := NewFaturadoService(store)

	err := service.MarkFaturado(context.Background(), "   ", "order-1")
	if !errors.Is(err, ErrFaturadoInvalidInput) {
		t.Fatalf("MarkFaturado() error = %v, want ErrFaturadoInvalidInput", err)
	}
	if store.called {
		t.Fatal("MarkFaturado() called the store on invalid input")
	}
}

func TestFaturadoServiceMarkFaturadoMissingProviderOrderID(t *testing.T) {
	store := &fakeOrderFaturadoStore{}
	service := NewFaturadoService(store)

	err := service.MarkFaturado(context.Background(), "installation-1", "")
	if !errors.Is(err, ErrFaturadoInvalidInput) {
		t.Fatalf("MarkFaturado() error = %v, want ErrFaturadoInvalidInput", err)
	}
	if store.called {
		t.Fatal("MarkFaturado() called the store on invalid input")
	}
}

func TestFaturadoServiceMarkFaturadoPropagatesStoreError(t *testing.T) {
	wantErr := errors.New("faturado store unavailable")
	service := NewFaturadoService(&fakeOrderFaturadoStore{err: wantErr})

	err := service.MarkFaturado(context.Background(), "installation-1", "order-1")
	if !errors.Is(err, wantErr) {
		t.Fatalf("MarkFaturado() error = %v, want %v", err, wantErr)
	}
}

func TestFaturadoServiceMarkFaturadoPropagatesNotFound(t *testing.T) {
	notFound := &ports.OrderNotFoundError{InstallationID: "installation-1", ProviderOrderID: "order-1"}
	service := NewFaturadoService(&fakeOrderFaturadoStore{err: notFound})

	err := service.MarkFaturado(context.Background(), "installation-1", "order-1")
	if !errors.Is(err, ports.ErrOrderNotFound) {
		t.Fatalf("MarkFaturado() error = %v, want ports.ErrOrderNotFound", err)
	}
}
