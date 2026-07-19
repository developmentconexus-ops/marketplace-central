package tarifflive

import (
	"context"
	"errors"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// --- recording fakes -------------------------------------------------------

type fakeIdentity struct {
	id  ports.ProductIdentity
	ok  bool
	err error
	got *int
}

func (f *fakeIdentity) IdentityForProduct(_ context.Context, productID int) (ports.ProductIdentity, bool, error) {
	f.got = &productID
	return f.id, f.ok, f.err
}

type fakeCategory struct {
	res   ports.CategoryResolution
	ok    bool
	err   error
	gotQ  *ports.CategoryQuery
	calls int
}

func (f *fakeCategory) ResolveCategory(_ context.Context, q ports.CategoryQuery) (ports.CategoryResolution, bool, error) {
	f.calls++
	f.gotQ = &q
	return f.res, f.ok, f.err
}

type fakeCommission struct {
	quote ports.CommissionQuote
	ok    bool
	err   error
	gotQ  *ports.CommissionQuery
	calls int
}

func (f *fakeCommission) QuoteCommission(_ context.Context, q ports.CommissionQuery) (ports.CommissionQuote, bool, error) {
	f.calls++
	f.gotQ = &q
	return f.quote, f.ok, f.err
}

func intp(i int) *int       { return &i }
func strp(s string) *string { return &s }

func newResolver(id *fakeIdentity, cat *fakeCategory, com *fakeCommission) *Resolver {
	return NewResolver(id, cat, com)
}

// --- happy path ------------------------------------------------------------

// Full degrau-3 resolution: identity -> category (D-83 sends EAN+titulo) ->
// commission, stamped {COTACAO, degrau 3, source timestamp}. The commission is
// quoted at the request price basis under the modalidade's listing type.
func TestResolveCommission_Happy(t *testing.T) {
	id := &fakeIdentity{id: ports.ProductIdentity{EAN: "789", Titulo: "Tenis X"}, ok: true}
	cat := &fakeCategory{res: ports.CategoryResolution{CategoriaID: "MLB1234", Fonte: ports.FonteCategoriaEANCatalogo, Data: "2026-07-19T10:00:00Z"}, ok: true}
	com := &fakeCommission{quote: ports.CommissionQuote{ComissaoPct: "13.00", Data: "2026-07-19T12:00:00Z"}, ok: true}
	r := newResolver(id, cat, com)

	res, ok, err := r.ResolveCommission(context.Background(), ports.TariffRequest{
		Modalidade: domain.ModalidadeClassico, ProductID: intp(7), PriceBasis: strp("199.90"),
	})
	if err != nil || !ok {
		t.Fatalf("want found no-error, got ok=%v err=%v", ok, err)
	}
	if res.Valor == nil || *res.Valor != "13.00" {
		t.Fatalf("Valor = %v, want 13.00", res.Valor)
	}
	if res.Fonte != domain.FonteCotacao || res.Degrau != 3 {
		t.Fatalf("stamp = {%q,%d}, want {COTACAO,3}", res.Fonte, res.Degrau)
	}
	if res.Data == nil || *res.Data != "2026-07-19T12:00:00Z" {
		t.Fatalf("Data = %v, want fee-read timestamp", res.Data)
	}
	// D-83: category probe carries EAN and titulo together.
	if cat.gotQ == nil || cat.gotQ.EAN != "789" || cat.gotQ.Titulo != "Tenis X" {
		t.Fatalf("category query = %+v, want EAN+titulo together (D-83)", cat.gotQ)
	}
	// Commission quoted at the request price under clássico listing type.
	if com.gotQ == nil || com.gotQ.CategoriaID != "MLB1234" || com.gotQ.Preco.Amount != "199.90" || com.gotQ.ListingTypeID != listingTypeClassico {
		t.Fatalf("commission query = %+v", com.gotQ)
	}
}

// --- listing type mapping --------------------------------------------------

func TestResolveCommission_ListingTypeByModalidade(t *testing.T) {
	cases := map[domain.Modalidade]string{
		domain.ModalidadeClassico: listingTypeClassico,
		domain.ModalidadePremium:  listingTypePremium,
		domain.ModalidadeFull:     listingTypePremium,
	}
	for mod, want := range cases {
		id := &fakeIdentity{id: ports.ProductIdentity{EAN: "789", Titulo: "X"}, ok: true}
		cat := &fakeCategory{res: ports.CategoryResolution{CategoriaID: "MLB1", Data: "t"}, ok: true}
		com := &fakeCommission{quote: ports.CommissionQuote{ComissaoPct: "10.00", Data: "t"}, ok: true}
		r := newResolver(id, cat, com)
		_, ok, err := r.ResolveCommission(context.Background(), ports.TariffRequest{Modalidade: mod, ProductID: intp(1), PriceBasis: strp("10.00")})
		if err != nil || !ok {
			t.Fatalf("mod %q: ok=%v err=%v", mod, ok, err)
		}
		if com.gotQ.ListingTypeID != want {
			t.Fatalf("mod %q -> listing %q, want %q", mod, com.gotQ.ListingTypeID, want)
		}
	}
}

// --- price basis fallback --------------------------------------------------

// With no request price basis, the commission is quoted at the product's
// catalog price.
func TestResolveCommission_PriceFallsBackToCatalog(t *testing.T) {
	price := domain.Money{Amount: "88.00", Currency: "BRL"}
	id := &fakeIdentity{id: ports.ProductIdentity{EAN: "789", Titulo: "X", Preco: &price}, ok: true}
	cat := &fakeCategory{res: ports.CategoryResolution{CategoriaID: "MLB1", Data: "t"}, ok: true}
	com := &fakeCommission{quote: ports.CommissionQuote{ComissaoPct: "10.00", Data: "t"}, ok: true}
	r := newResolver(id, cat, com)

	_, ok, err := r.ResolveCommission(context.Background(), ports.TariffRequest{Modalidade: domain.ModalidadeClassico, ProductID: intp(1)})
	if err != nil || !ok {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	if com.gotQ.Preco.Amount != "88.00" {
		t.Fatalf("quoted at %q, want catalog 88.00", com.gotQ.Preco.Amount)
	}
}

// --- silent-miss matrix ----------------------------------------------------

func TestResolveCommission_Misses(t *testing.T) {
	full := func() (*fakeIdentity, *fakeCategory, *fakeCommission) {
		return &fakeIdentity{id: ports.ProductIdentity{EAN: "789", Titulo: "X"}, ok: true},
			&fakeCategory{res: ports.CategoryResolution{CategoriaID: "MLB1", Data: "t"}, ok: true},
			&fakeCommission{quote: ports.CommissionQuote{ComissaoPct: "10.00", Data: "t"}, ok: true}
	}

	t.Run("nil product id => degrau 3 not applicable, no reads", func(t *testing.T) {
		id, cat, com := full()
		r := newResolver(id, cat, com)
		_, ok, err := r.ResolveCommission(context.Background(), ports.TariffRequest{Modalidade: domain.ModalidadeClassico, PriceBasis: strp("10.00")})
		if ok || err != nil {
			t.Fatalf("ok=%v err=%v, want silent miss", ok, err)
		}
		if id.got != nil || cat.calls != 0 || com.calls != 0 {
			t.Fatalf("nil product id must not read anything")
		}
	})

	t.Run("identity miss => no category probe", func(t *testing.T) {
		id, cat, com := full()
		id.ok = false
		r := newResolver(id, cat, com)
		_, ok, err := r.ResolveCommission(context.Background(), ports.TariffRequest{Modalidade: domain.ModalidadeClassico, ProductID: intp(1), PriceBasis: strp("10.00")})
		if ok || err != nil {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if cat.calls != 0 {
			t.Fatalf("identity miss must not probe category")
		}
	})

	t.Run("category miss => no commission quote", func(t *testing.T) {
		id, cat, com := full()
		cat.ok = false
		r := newResolver(id, cat, com)
		_, ok, err := r.ResolveCommission(context.Background(), ports.TariffRequest{Modalidade: domain.ModalidadeClassico, ProductID: intp(1), PriceBasis: strp("10.00")})
		if ok || err != nil {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if com.calls != 0 {
			t.Fatalf("category miss must not quote commission")
		}
	})

	t.Run("no price basis and no catalog price => no commission quote", func(t *testing.T) {
		id, cat, com := full()
		r := newResolver(id, cat, com)
		_, ok, err := r.ResolveCommission(context.Background(), ports.TariffRequest{Modalidade: domain.ModalidadeClassico, ProductID: intp(1)})
		if ok || err != nil {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if com.calls != 0 {
			t.Fatalf("no price must not quote commission")
		}
	})

	t.Run("commission miss => silent miss", func(t *testing.T) {
		id, cat, com := full()
		com.ok = false
		r := newResolver(id, cat, com)
		_, ok, err := r.ResolveCommission(context.Background(), ports.TariffRequest{Modalidade: domain.ModalidadeClassico, ProductID: intp(1), PriceBasis: strp("10.00")})
		if ok || err != nil {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
	})

	t.Run("unknown modalidade => silent miss, no commission quote", func(t *testing.T) {
		id, cat, com := full()
		r := newResolver(id, cat, com)
		_, ok, err := r.ResolveCommission(context.Background(), ports.TariffRequest{Modalidade: domain.Modalidade("bogus"), ProductID: intp(1), PriceBasis: strp("10.00")})
		if ok || err != nil {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if com.calls != 0 {
			t.Fatalf("unknown modalidade must not quote commission")
		}
	})
}

// --- error propagation -----------------------------------------------------

// A read error propagates (the composite decides the fallback policy); the
// assembler never silently maps an error to a miss.
func TestResolveCommission_ErrorPropagates(t *testing.T) {
	sentinel := errors.New("boom")
	id := &fakeIdentity{id: ports.ProductIdentity{EAN: "789", Titulo: "X"}, ok: true}
	cat := &fakeCategory{err: sentinel}
	com := &fakeCommission{}
	r := newResolver(id, cat, com)

	_, ok, err := r.ResolveCommission(context.Background(), ports.TariffRequest{Modalidade: domain.ModalidadeClassico, ProductID: intp(1), PriceBasis: strp("10.00")})
	if ok {
		t.Fatalf("want not found on error")
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want sentinel", err)
	}
}
