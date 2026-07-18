//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	pricingpostgres "marketplace-central/apps/server_core/internal/modules/pricing/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestCalcRepositoryProfileRoundTripAndDefault(t *testing.T) {
	testpostgres.SkipWithoutTarget(t)
	ctx := context.Background()
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	pool, _ := testpostgres.OpenPool(t, "pricing_calc_"+token)
	repo := pricingpostgres.NewCalcRepository(pool)
	tenant := "calc-profile-" + token

	// Absent profile materializes the built-in default, never a zero row.
	got, err := repo.GetProfile(ctx, tenant)
	if err != nil {
		t.Fatal(err)
	}
	want := domain.NewDefaultCalcProfile()
	if got != want {
		t.Fatalf("default profile = %#v, want %#v", got, want)
	}

	// Upsert with a set TarifaFull, then a clear (nil) — both must round-trip
	// exactly, and nil must stay nil (ADR-17), not become 0.
	edited := domain.CalcProfile{
		Regime:           domain.RegimePresumido,
		AliquotaPct:      "9.25",
		LimiarVerdePct:   "20",
		LimiarAmareloPct: "11",
		TarifaFull:       &domain.Money{Amount: "6.50", Currency: "BRL"},
		DifalEnabled:     true,
		DifalDestinoUF:   strptr("BA"),
		Origem:           "operator",
	}
	if err := repo.UpsertProfile(ctx, tenant, edited); err != nil {
		t.Fatal(err)
	}
	got, err = repo.GetProfile(ctx, tenant)
	if err != nil {
		t.Fatal(err)
	}
	if got.TarifaFull == nil || got.TarifaFull.Amount != "6.50" {
		t.Fatalf("tarifa_full = %#v, want 6.50", got.TarifaFull)
	}
	if got.AliquotaPct != "9.25" || got.Regime != domain.RegimePresumido || got.Origem != "operator" {
		t.Fatalf("profile round-trip = %#v", got)
	}
	if got.DifalDestinoUF == nil || *got.DifalDestinoUF != "BA" {
		t.Fatalf("difal_destino_uf = %#v, want BA", got.DifalDestinoUF)
	}

	edited.TarifaFull = nil
	if err := repo.UpsertProfile(ctx, tenant, edited); err != nil {
		t.Fatal(err)
	}
	got, err = repo.GetProfile(ctx, tenant)
	if err != nil {
		t.Fatal(err)
	}
	if got.TarifaFull != nil {
		t.Fatalf("cleared tarifa_full = %#v, want nil", got.TarifaFull)
	}
}

func TestCalcRepositoryDifalReadOverrideAndNotFound(t *testing.T) {
	testpostgres.SkipWithoutTarget(t)
	ctx := context.Background()
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	pool, _ := testpostgres.OpenPool(t, "pricing_difal_"+token)
	repo := pricingpostgres.NewCalcRepository(pool)
	tenant := "calc-difal-" + token

	// Seed two rows for this tenant (0057 only seeds tenant_default).
	seedDifal(t, ctx, pool, tenant, "SP", "18.0", "12", "6.00")
	seedDifal(t, ctx, pool, tenant, "BA", "20.5", "7", "13.50")

	// Absent UF → sentinel.
	if _, err := repo.RateForUF(ctx, tenant, "ZZ"); !errors.Is(err, ports.ErrDifalUFNotFound) {
		t.Fatalf("RateForUF(ZZ) err = %v, want ErrDifalUFNotFound", err)
	}

	// Seed read: efetivo via domain.DifalForUF = max(interna-inter,0).
	rate, err := repo.RateForUF(ctx, tenant, "BA")
	if err != nil {
		t.Fatal(err)
	}
	if rate.Override != nil {
		t.Fatalf("BA override = %#v, want nil", rate.Override)
	}
	if res := domain.DifalForUF(rate); res.EfetivoPct != "13.50" {
		t.Fatalf("BA efetivo = %q, want 13.50", res.EfetivoPct)
	}

	// Apply an override at 2dp; DifalForUF must use the override interna.
	ov := &domain.DifalOverride{InternaPct: "19.05", UpdatedAt: "2026-07-18T10:00:00Z", Actor: "op-1"}
	if err := repo.UpsertDifalOverride(ctx, tenant, "BA", ov); err != nil {
		t.Fatal(err)
	}
	rate, err = repo.RateForUF(ctx, tenant, "BA")
	if err != nil {
		t.Fatal(err)
	}
	if rate.Override == nil || rate.Override.InternaPct != "19.05" || rate.Override.Actor != "op-1" {
		t.Fatalf("override round-trip = %#v", rate.Override)
	}
	if res := domain.DifalForUF(rate); res.EfetivoPct != "12.05" { // 19.05 - 7
		t.Fatalf("overridden efetivo = %q, want 12.05", res.EfetivoPct)
	}

	// ListDifalRates ordered by UF (BA before SP), reflecting the override.
	rates, err := repo.ListDifalRates(ctx, tenant)
	if err != nil {
		t.Fatal(err)
	}
	if len(rates) != 2 || rates[0].UF != "BA" || rates[1].UF != "SP" {
		t.Fatalf("list = %#v, want [BA, SP]", rates)
	}
	if rates[0].Override == nil {
		t.Fatalf("BA in list missing override")
	}

	// Clear the override.
	if err := repo.UpsertDifalOverride(ctx, tenant, "BA", nil); err != nil {
		t.Fatal(err)
	}
	rate, err = repo.RateForUF(ctx, tenant, "BA")
	if err != nil {
		t.Fatal(err)
	}
	if rate.Override != nil {
		t.Fatalf("cleared override = %#v, want nil", rate.Override)
	}

	// Override on an absent UF → sentinel.
	if err := repo.UpsertDifalOverride(ctx, tenant, "ZZ", ov); !errors.Is(err, ports.ErrDifalUFNotFound) {
		t.Fatalf("override(ZZ) err = %v, want ErrDifalUFNotFound", err)
	}
}

func TestCalcRepositoryScenariosCRUDNewestFirst(t *testing.T) {
	testpostgres.SkipWithoutTarget(t)
	ctx := context.Background()
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	pool, _ := testpostgres.OpenPool(t, "pricing_scn_"+token)
	repo := pricingpostgres.NewCalcRepository(pool)
	tenant := "calc-scn-" + token
	other := "calc-scn-other-" + token

	older := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC)
	if err := repo.SaveScenario(ctx, tenant, domain.Scenario{ID: "s1", Name: "primeiro", Payload: `{"preco":"100.00"}`, CreatedAt: older}); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveScenario(ctx, tenant, domain.Scenario{ID: "s2", Name: "segundo", Payload: `{"preco":"200.00"}`, CreatedAt: newer}); err != nil {
		t.Fatal(err)
	}
	// Other tenant's scenario must not leak.
	if err := repo.SaveScenario(ctx, other, domain.Scenario{ID: "s3", Name: "outro", Payload: `{}`, CreatedAt: newer}); err != nil {
		t.Fatal(err)
	}

	got, err := repo.ListScenarios(ctx, tenant)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].ID != "s2" || got[1].ID != "s1" {
		t.Fatalf("list = %#v, want newest-first [s2, s1]", got)
	}

	// Upsert (same id) updates name/payload, not a duplicate row.
	if err := repo.SaveScenario(ctx, tenant, domain.Scenario{ID: "s1", Name: "primeiro-v2", Payload: `{"preco":"150.00"}`, CreatedAt: older}); err != nil {
		t.Fatal(err)
	}
	got, err = repo.ListScenarios(ctx, tenant)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("after upsert len = %d, want 2", len(got))
	}

	// Delete then re-delete → sentinel.
	if err := repo.DeleteScenario(ctx, tenant, "s1"); err != nil {
		t.Fatal(err)
	}
	if err := repo.DeleteScenario(ctx, tenant, "s1"); !errors.Is(err, ports.ErrScenarioNotFound) {
		t.Fatalf("re-delete err = %v, want ErrScenarioNotFound", err)
	}
	got, err = repo.ListScenarios(ctx, tenant)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "s2" {
		t.Fatalf("after delete = %#v, want [s2]", got)
	}
}

func seedDifal(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tenant, uf, interna, inter, efetivo string) {
	t.Helper()
	_, err := pool.Exec(ctx, `
		INSERT INTO pricing_difal_rates (tenant_id, uf, interna_pct, interestadual_pct, efetivo_pct, origem_versao)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (tenant_id, uf) DO NOTHING
	`, tenant, uf, interna, inter, efetivo, domain.OrigemVersaoPadrao2026)
	if err != nil {
		t.Fatal(err)
	}
}

func strptr(s string) *string { return &s }
