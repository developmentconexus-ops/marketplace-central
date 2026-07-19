package catalog

import (
	"context"
	"errors"
	"testing"

	catalogdomain "marketplace-central/apps/server_core/internal/modules/catalog/domain"
	pricingports "marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

var _ pricingports.ProductIdentityReader = (*ProductIdentityReader)(nil)

// A known product yields its identity: EAN + title feed the catalog match
// (D-83), and a positive catalog price becomes the fallback quote basis. The
// int product id is passed to the string-keyed catalog service as its decimal.
func TestProductIdentityReader_KnownProduct(t *testing.T) {
	svc := &fakeCatalogService{productsByID: map[string]catalogdomain.Product{
		"7": {ProductID: "7", Name: "Tenis X", EAN: "7891234567890", PriceAmount: 199.9},
	}}
	r := NewProductIdentityReader(svc)

	id, ok, err := r.IdentityForProduct(context.Background(), 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("want found=true for known product")
	}
	if id.EAN != "7891234567890" || id.Titulo != "Tenis X" {
		t.Fatalf("identity mismatch: %+v", id)
	}
	if id.Preco == nil || id.Preco.Amount != "199.90" || id.Preco.Currency != "BRL" {
		t.Fatalf("Preco mismatch: %+v", id.Preco)
	}
	if got := svc.listProductsByIDsIDs; len(got) != 1 || len(got[0]) != 1 || got[0][0] != "7" {
		t.Fatalf("service queried with %#v, want [[7]]", got)
	}
}

// A non-positive catalog price is absent, never a fabricated 0 quote basis
// (ADR-17). The product is still found (its EAN/title can drive a match).
func TestProductIdentityReader_ZeroPriceIsAbsent(t *testing.T) {
	svc := &fakeCatalogService{productsByID: map[string]catalogdomain.Product{
		"7": {ProductID: "7", Name: "Tenis X", EAN: "789", PriceAmount: 0},
	}}
	r := NewProductIdentityReader(svc)

	id, ok, err := r.IdentityForProduct(context.Background(), 7)
	if err != nil || !ok {
		t.Fatalf("want found no-error, got ok=%v err=%v", ok, err)
	}
	if id.Preco != nil {
		t.Fatalf("Preco = %+v, want nil for zero price", id.Preco)
	}
}

// An unknown product id is a miss, not an error: the caller degrades to
// degrau 4.
func TestProductIdentityReader_UnknownProductIsMiss(t *testing.T) {
	svc := &fakeCatalogService{productsByID: map[string]catalogdomain.Product{}}
	r := NewProductIdentityReader(svc)

	_, ok, err := r.IdentityForProduct(context.Background(), 99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("want found=false for unknown product")
	}
}

// A product with neither EAN nor title carries no identity to match on: a
// miss, so no pointless downstream catalog probe.
func TestProductIdentityReader_EmptyIdentityIsMiss(t *testing.T) {
	svc := &fakeCatalogService{productsByID: map[string]catalogdomain.Product{
		"7": {ProductID: "7", Name: "  ", EAN: ""},
	}}
	r := NewProductIdentityReader(svc)

	_, ok, err := r.IdentityForProduct(context.Background(), 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("want found=false for empty identity")
	}
}

// A service error propagates — an unknown read is not silently a miss.
func TestProductIdentityReader_ServiceErrorPropagates(t *testing.T) {
	sentinel := errors.New("boom")
	svc := &fakeCatalogService{err: sentinel}
	r := NewProductIdentityReader(svc)

	_, _, err := r.IdentityForProduct(context.Background(), 7)
	if !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want sentinel", err)
	}
}
