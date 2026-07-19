package tariffcomposite

import (
	"context"
	"errors"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

func strp(s string) *string { return &s }

// fakeBase is the degrau-4 resolver: a fixed config-default resolution or an
// error.
type fakeBase struct {
	res domain.TariffResolution
	err error
}

func (f fakeBase) Resolve(context.Context, ports.TariffRequest) (domain.TariffResolution, error) {
	return f.res, f.err
}

// fakeLive is the degrau-3 commission assembler.
type fakeLive struct {
	res    domain.ComponentResolution
	ok     bool
	err    error
	called bool
}

func (f *fakeLive) ResolveCommission(context.Context, ports.TariffRequest) (domain.ComponentResolution, bool, error) {
	f.called = true
	return f.res, f.ok, f.err
}

func degrau4Base() domain.TariffResolution {
	c := "12.00"
	frete := "18.50"
	return domain.TariffResolution{
		Comissao: domain.ComponentResolution{Valor: &c, Fonte: domain.FontePadrao, Degrau: 4, Estimativa: true},
		Frete:    domain.ComponentResolution{Valor: &frete, Fonte: domain.FontePadrao, Degrau: 4, Estimativa: true},
	}
}

var _ ports.TariffResolver = (*Resolver)(nil)

// When degrau 3 resolves a commission, it replaces the degrau-4 commission;
// frete stays degrau 4 (degrau 3 quotes commission only in the MIN scope).
func TestResolve_Degrau3ReplacesCommissionKeepsFrete(t *testing.T) {
	base := fakeBase{res: degrau4Base()}
	live := &fakeLive{res: domain.ComponentResolution{Valor: strp("11.50"), Fonte: domain.FonteCotacao, Degrau: 3, Data: strp("2026-07-19T12:00:00Z")}, ok: true}
	r := NewResolver(base, live)

	got, err := r.Resolve(context.Background(), ports.TariffRequest{Modalidade: domain.ModalidadeClassico})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Comissao.Fonte != domain.FonteCotacao || got.Comissao.Degrau != 3 || got.Comissao.Valor == nil || *got.Comissao.Valor != "11.50" {
		t.Fatalf("commission not upgraded to degrau 3: %+v", got.Comissao)
	}
	if got.Comissao.Data == nil || *got.Comissao.Data != "2026-07-19T12:00:00Z" {
		t.Fatalf("commission lost its cotação timestamp: %+v", got.Comissao)
	}
	if got.Frete.Degrau != 4 || got.Frete.Valor == nil || *got.Frete.Valor != "18.50" {
		t.Fatalf("frete must stay degrau 4: %+v", got.Frete)
	}
}

// A degrau-3 miss keeps the honest degrau-4 config default unchanged.
func TestResolve_Degrau3MissKeepsDegrau4(t *testing.T) {
	base := fakeBase{res: degrau4Base()}
	live := &fakeLive{ok: false}
	r := NewResolver(base, live)

	got, err := r.Resolve(context.Background(), ports.TariffRequest{Modalidade: domain.ModalidadeClassico})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Comissao.Fonte != domain.FontePadrao || got.Comissao.Degrau != 4 {
		t.Fatalf("commission should stay degrau 4 on miss: %+v", got.Comissao)
	}
}

// A degrau-3 live-quote error falls SILENTLY to degrau 4 (a live probe failure
// never breaks a decompose and never fabricates a fact — it degrades to the
// lower-provenance config default). The error is not surfaced.
func TestResolve_Degrau3ErrorFallsSilentlyToDegrau4(t *testing.T) {
	base := fakeBase{res: degrau4Base()}
	live := &fakeLive{err: errors.New("ml timeout")}
	r := NewResolver(base, live)

	got, err := r.Resolve(context.Background(), ports.TariffRequest{Modalidade: domain.ModalidadeClassico})
	if err != nil {
		t.Fatalf("degrau-3 error must not surface: %v", err)
	}
	if got.Comissao.Degrau != 4 {
		t.Fatalf("commission should fall to degrau 4: %+v", got.Comissao)
	}
}

// A degrau-4 (base) error is a real failure and propagates; degrau 3 is not
// consulted (there is no base to upgrade).
func TestResolve_Degrau4ErrorPropagates(t *testing.T) {
	sentinel := errors.New("no defaults")
	base := fakeBase{err: sentinel}
	live := &fakeLive{res: domain.ComponentResolution{Valor: strp("11.50"), Fonte: domain.FonteCotacao, Degrau: 3}, ok: true}
	r := NewResolver(base, live)

	_, err := r.Resolve(context.Background(), ports.TariffRequest{Modalidade: domain.ModalidadeClassico})
	if !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want sentinel", err)
	}
	if live.called {
		t.Fatalf("degrau 3 must not run when base degrau 4 failed")
	}
}
