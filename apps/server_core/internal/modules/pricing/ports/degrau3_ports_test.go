package ports

import (
	"context"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
)

// TariffRequest grows two OPTIONAL degrau-3 inputs (per-product id + the price
// basis to quote against) without breaking degrau 4, which reads only
// Modalidade. Both nil => degrau 3 is skipped, degrau 4 answers (ADR-17: a
// missing input is never coerced to a live quote).
func TestTariffRequestCarriesDegrau3Inputs(t *testing.T) {
	pid := 42
	basis := "199.90"
	req := TariffRequest{Modalidade: domain.ModalidadeClassico, ProductID: &pid, PriceBasis: &basis}
	if req.ProductID == nil || *req.ProductID != 42 {
		t.Fatalf("ProductID not carried: %v", req.ProductID)
	}
	if req.PriceBasis == nil || *req.PriceBasis != basis {
		t.Fatalf("PriceBasis not carried: %v", req.PriceBasis)
	}
}

// fakeProductIdentityReader / fakeCategoryResolver / fakeCommissionQuoter prove
// the three new degrau-3 ports are satisfiable and pin their neutral I/O
// shapes (no ML payload leaks across the port — those are strings/domain.Money).
type fakeProductIdentityReader struct {
	id  ProductIdentity
	ok  bool
	err error
}

func (f fakeProductIdentityReader) IdentityForProduct(context.Context, int) (ProductIdentity, bool, error) {
	return f.id, f.ok, f.err
}

type fakeCategoryResolver struct {
	res CategoryResolution
	ok  bool
	err error
}

func (f fakeCategoryResolver) ResolveCategory(context.Context, CategoryQuery) (CategoryResolution, bool, error) {
	return f.res, f.ok, f.err
}

type fakeCommissionQuoter struct {
	quote CommissionQuote
	ok    bool
	err   error
}

func (f fakeCommissionQuoter) QuoteCommission(context.Context, CommissionQuery) (CommissionQuote, bool, error) {
	return f.quote, f.ok, f.err
}

var (
	_ ProductIdentityReader = fakeProductIdentityReader{}
	_ CategoryResolver      = fakeCategoryResolver{}
	_ CommissionQuoter      = fakeCommissionQuoter{}
)

// TestDegrau3PortResultsCarryProvenance pins the honest degrau-3 outputs: a
// resolved category stamps its provenance fonte + source timestamp, and a
// commission quote carries the percent + the quote's source timestamp. Absent
// (ok=false) results are the ADR-17 miss path, verified in the assembler slice.
func TestDegrau3PortResultsCarryProvenance(t *testing.T) {
	price := domain.Money{Amount: "199.90", Currency: "BRL"}
	id := ProductIdentity{EAN: "7891234567890", Titulo: "Tenis X", Preco: &price}
	if id.Preco == nil || id.Preco.Amount != "199.90" {
		t.Fatalf("ProductIdentity.Preco not carried")
	}
	cat := CategoryResolution{CategoriaID: "MLB1234", Fonte: FonteCategoriaEANCatalogo, Data: "2026-07-19T12:00:00Z"}
	if cat.CategoriaID != "MLB1234" || cat.Fonte != FonteCategoriaEANCatalogo || cat.Data == "" {
		t.Fatalf("CategoryResolution not stamped: %+v", cat)
	}
	q := CommissionQuote{ComissaoPct: "13.00", Data: "2026-07-19T12:00:00Z"}
	if q.ComissaoPct != "13.00" || q.Data == "" {
		t.Fatalf("CommissionQuote not stamped: %+v", q)
	}
}
