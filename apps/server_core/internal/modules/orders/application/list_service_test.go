package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

type readStoreStub struct {
	page       ports.OrderPage
	listErr    error
	model      domain.OrderReadModel
	getErr     error
	gotQuery   ports.OrderListQuery
	gotInstall string
	gotOrder   string
}

func (s *readStoreStub) ListOrders(_ context.Context, query ports.OrderListQuery) (ports.OrderPage, error) {
	s.gotQuery = query
	return s.page, s.listErr
}

func (s *readStoreStub) GetOrder(_ context.Context, installationID, providerOrderID string) (domain.OrderReadModel, error) {
	s.gotInstall, s.gotOrder = installationID, providerOrderID
	return s.model, s.getErr
}

func TestListServiceDelegatesCanonicalListAndDetailReads(t *testing.T) {
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	store := &readStoreStub{page: ports.OrderPage{
		Items: []domain.OrderReadModel{{ProviderOrderID: "order-1", CreatedAt: &now}},
	}}
	service := NewReadService(store)
	query := ports.OrderListQuery{InstallationID: "inst-1", Limit: 10, Status: "paid"}

	page, err := service.List(context.Background(), query)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(page.Items) != 1 || page.Items[0].ProviderOrderID != "order-1" {
		t.Fatalf("page = %+v", page)
	}
	if store.gotQuery != query {
		t.Fatalf("query = %+v, want %+v", store.gotQuery, query)
	}

	model, err := service.Get(context.Background(), "inst-1", "order-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if model.ProviderOrderID != "" || store.gotInstall != "inst-1" || store.gotOrder != "order-1" {
		t.Fatalf("model=%+v install=%q order=%q", model, store.gotInstall, store.gotOrder)
	}
}

func TestListServicePropagatesReadErrors(t *testing.T) {
	want := errors.New("read failed")
	store := &readStoreStub{listErr: want, getErr: want}
	service := NewReadService(store)

	if _, err := service.List(context.Background(), ports.OrderListQuery{InstallationID: "inst-1"}); !errors.Is(err, want) {
		t.Fatalf("List() error = %v, want %v", err, want)
	}
	if _, err := service.Get(context.Background(), "inst-1", "order-1"); !errors.Is(err, want) {
		t.Fatalf("Get() error = %v, want %v", err, want)
	}
}
