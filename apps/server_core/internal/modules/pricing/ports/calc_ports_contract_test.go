package ports

import (
	"context"
	"reflect"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
)

// This file FREEZES the IC-04 port surface consumed by M-08. A change here is
// a contract break: it must be a deliberate, coordinated amendment, never an
// incidental edit.

// Compile-time freeze of the single decomposition function's signature.
var _ func(domain.DecomposeInput) domain.Decomposition = domain.Decompose

// engineAdapter proves CalcEngine is satisfiable by the pure formula without a
// speculative production type (real wiring lands where the engine is
// constructed, S6/S7).
type engineAdapter struct{}

func (engineAdapter) Decompose(in domain.DecomposeInput) domain.Decomposition {
	return domain.Decompose(in)
}

var _ CalcEngine = engineAdapter{}

// TestDecompositionShapeFrozen pins the exact output fields M-08 reads. The
// nullable components (unknown ≠ zero) MUST stay *string, the always-known
// ones plain string, and margem both nullable.
//
// P2.b Task 4 amendment: ICMSSaida, PisCofins, RestituicaoST added (D-41 ex
// ante tax components) — deliberate, coordinated per the plan's Task 4 scope.
//
// Task A7 amendment: Imposto is now *string. Task A6 made the legacy D-38
// regime-aliquota field a named-absent nil whenever ICMSCell is present (the
// real per-item tax path — ICMSSaida/Difal/PisCofins/RestituicaoST — wins and
// the legacy flat rate is never parsed/shown alongside it); the port surface
// has to carry that nullability through, deliberate and coordinated per this
// task's contract-seam scope, never an incidental widen.
func TestDecompositionShapeFrozen(t *testing.T) {
	want := map[string]string{
		"Preco":                    "string",
		"Comissao":                 "string",
		"TaxaFixa":                 "string",
		"Frete":                    "*string",
		"Imposto":                  "*string",
		"Difal":                    "*string",
		"TarifaFull":               "*string",
		"Custo":                    "*string",
		"ICMSSaida":                "*string",
		"PisCofins":                "*string",
		"RestituicaoST":            "*string",
		"MargemValor":              "*string",
		"MargemPct":                "*string",
		"ComponentesDesconhecidos": "[]string",
	}
	assertFields(t, reflect.TypeOf(domain.Decomposition{}), want)
}

// TestDifalForUFResultShapeFrozen pins DifalForUF(uf) → {efetivo_pct, versao}.
func TestDifalForUFResultShapeFrozen(t *testing.T) {
	want := map[string]string{
		"EfetivoPct": "string",
		"Versao":     "string",
	}
	assertFields(t, reflect.TypeOf(domain.DifalForUFResult{}), want)
}

// TestDecomposeInputShapeFrozen pins the resolved-input contract adapters fill.
//
// P2.b Task 4 amendment: ICMSCell added (D-41 pre-resolved tax-matrix cell +
// product fiscal facts) — deliberate, coordinated per the plan's Task 4 scope.
func TestDecomposeInputShapeFrozen(t *testing.T) {
	want := map[string]string{
		"Preco":        "string",
		"ComissaoPct":  "string",
		"AliquotaPct":  "string",
		"Modalidade":   "domain.Modalidade",
		"TarifaFull":   "*domain.Money",
		"FreteProduto": "*domain.Money",
		"Custo":        "*domain.Money",
		"DifalEnabled": "bool",
		"DestinoUF":    "string",
		"EfetivoPct":   "string",
		"ICMSCell":     "*domain.ICMSCell",
	}
	assertFields(t, reflect.TypeOf(domain.DecomposeInput{}), want)
}

// fakeDifalReader records the tenant/uf it was asked for and returns a fixed
// rate carrying an active override.
type fakeDifalReader struct {
	gotTenant string
	gotUF     string
	rate      domain.DifalRate
	err       error
}

func (f *fakeDifalReader) RateForUF(_ context.Context, tenantID, uf string) (domain.DifalRate, error) {
	f.gotTenant = tenantID
	f.gotUF = uf
	return f.rate, f.err
}

// TestDifalForUFPortReflectsOverrideAndScopesTenant freezes the port behavior
// M-08 relies on: the override wins over the seed, and the read is tenant
// scoped.
func TestDifalForUFPortReflectsOverrideAndScopesTenant(t *testing.T) {
	reader := &fakeDifalReader{
		rate: domain.DifalRate{
			UF:               "BA",
			InternaPct:       "20.5",
			InterestadualPct: "7",
			EfetivoPct:       "13.50",
			OrigemVersao:     "padrao-2026",
			Override: &domain.DifalOverride{
				InternaPct: "21.0",
				UpdatedAt:  "2026-07-18T00:00:00Z",
				Actor:      "operator-1",
			},
		},
	}
	got, err := DifalForUF(context.Background(), reader, "tenant-xyz", "BA")
	if err != nil {
		t.Fatalf("DifalForUF: unexpected error %v", err)
	}
	// override 21.0 − interestadual 7 = 14.00 (seed efetivo 13.50 is overridden).
	if got.EfetivoPct != "14.00" {
		t.Fatalf("EfetivoPct = %q, want 14.00 (override must win)", got.EfetivoPct)
	}
	if got.Versao != "padrao-2026" {
		t.Fatalf("Versao = %q, want padrao-2026", got.Versao)
	}
	if reader.gotTenant != "tenant-xyz" {
		t.Fatalf("reader tenant = %q, want tenant-xyz (read must be tenant scoped)", reader.gotTenant)
	}
	if reader.gotUF != "BA" {
		t.Fatalf("reader uf = %q, want BA", reader.gotUF)
	}
}

func assertFields(t *testing.T, typ reflect.Type, want map[string]string) {
	t.Helper()
	if typ.NumField() != len(want) {
		t.Fatalf("%s has %d fields, want %d (contract change)", typ.Name(), typ.NumField(), len(want))
	}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		w, ok := want[f.Name]
		if !ok {
			t.Fatalf("%s has unexpected field %q (contract change)", typ.Name(), f.Name)
		}
		if got := f.Type.String(); got != w {
			t.Fatalf("%s.%s type = %q, want %q (contract change)", typ.Name(), f.Name, got, w)
		}
	}
}
