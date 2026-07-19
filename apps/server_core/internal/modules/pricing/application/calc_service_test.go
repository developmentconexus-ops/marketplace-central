package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/pricing/application"
	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// fakeCalcRepo is a configurable in-memory CalcRepository for service tests.
type fakeCalcRepo struct {
	profile      *domain.CalcProfile
	rates        map[string]domain.DifalRate
	scenarios    []domain.Scenario
	upserted     *domain.CalcProfile
	overrideCall *struct {
		uf string
		ov *domain.DifalOverride
	}
}

func newFakeRepo() *fakeCalcRepo {
	return &fakeCalcRepo{rates: map[string]domain.DifalRate{}}
}

func (f *fakeCalcRepo) GetProfile(_ context.Context, _ string) (domain.CalcProfile, error) {
	if f.profile != nil {
		return *f.profile, nil
	}
	return domain.NewDefaultCalcProfile(), nil
}

func (f *fakeCalcRepo) UpsertProfile(_ context.Context, _ string, p domain.CalcProfile) error {
	f.upserted = &p
	f.profile = &p
	return nil
}

func (f *fakeCalcRepo) RateForUF(_ context.Context, _, uf string) (domain.DifalRate, error) {
	r, ok := f.rates[uf]
	if !ok {
		return domain.DifalRate{}, ports.ErrDifalUFNotFound
	}
	return r, nil
}

func (f *fakeCalcRepo) ListDifalRates(_ context.Context, _ string) ([]domain.DifalRate, error) {
	out := make([]domain.DifalRate, 0, len(f.rates))
	for _, r := range f.rates {
		out = append(out, r)
	}
	return out, nil
}

func (f *fakeCalcRepo) UpsertDifalOverride(_ context.Context, _, uf string, ov *domain.DifalOverride) error {
	f.overrideCall = &struct {
		uf string
		ov *domain.DifalOverride
	}{uf: uf, ov: ov}
	r := f.rates[uf]
	r.Override = ov
	f.rates[uf] = r
	return nil
}

func (f *fakeCalcRepo) ListScenarios(_ context.Context, _ string) ([]domain.Scenario, error) {
	return f.scenarios, nil
}
func (f *fakeCalcRepo) SaveScenario(_ context.Context, _ string, s domain.Scenario) error {
	f.scenarios = append(f.scenarios, s)
	return nil
}
func (f *fakeCalcRepo) DeleteScenario(_ context.Context, _, id string) error {
	for i, s := range f.scenarios {
		if s.ID == id {
			f.scenarios = append(f.scenarios[:i], f.scenarios[i+1:]...)
			return nil
		}
	}
	return ports.ErrScenarioNotFound
}

type fakeCost struct {
	money *domain.Money
	err   error
}

func (f fakeCost) CostForProduct(_ context.Context, _ int, _ time.Time) (*domain.Money, error) {
	return f.money, f.err
}

type fakeProducts struct{ exists map[int]bool }

func (f fakeProducts) Exists(_ context.Context, id int) (bool, error) {
	return f.exists[id], nil
}

func newService(repo ports.CalcRepository, cost ports.CostReader, products ports.ProductChecker) application.CalcService {
	return application.NewCalcService(repo, cost, products, "t1")
}

func TestPutProfileInvalidRate(t *testing.T) {
	svc := newService(newFakeRepo(), fakeCost{}, nil)
	_, err := svc.PutProfile(context.Background(), application.ProfileUpdate{
		Regime: domain.RegimeSimples, AliquotaPct: "40", LimiarVerdePct: "18", LimiarAmareloPct: "10",
	})
	if !errors.Is(err, application.ErrInvalidRate) {
		t.Fatalf("err = %v, want ErrInvalidRate", err)
	}
}

func TestPutProfileValidStampsOperator(t *testing.T) {
	repo := newFakeRepo()
	svc := newService(repo, fakeCost{}, nil)
	got, err := svc.PutProfile(context.Background(), application.ProfileUpdate{
		Regime: domain.RegimeSimples, AliquotaPct: "4", LimiarVerdePct: "18", LimiarAmareloPct: "10",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Origem != "operator" {
		t.Fatalf("origem = %q, want operator", got.Origem)
	}
	if repo.upserted == nil || repo.upserted.AliquotaPct != "4" {
		t.Fatalf("not persisted: %#v", repo.upserted)
	}
}

func TestPutDifalOverrideInvalidRate(t *testing.T) {
	svc := newService(newFakeRepo(), fakeCost{}, nil)
	_, err := svc.PutDifalOverride(context.Background(), "SP", "40", "op", time.Now())
	if !errors.Is(err, application.ErrInvalidRate) {
		t.Fatalf("err = %v, want ErrInvalidRate", err)
	}
}

func TestPutDifalOverrideUFNotFound(t *testing.T) {
	svc := newService(newFakeRepo(), fakeCost{}, nil)
	_, err := svc.PutDifalOverride(context.Background(), "ZZ", "19", "op", time.Now())
	if !errors.Is(err, ports.ErrDifalUFNotFound) {
		t.Fatalf("err = %v, want ErrDifalUFNotFound", err)
	}
}

func TestPutDifalOverrideBelowThresholdNoPersist(t *testing.T) {
	repo := newFakeRepo()
	repo.rates["SP"] = domain.DifalRate{UF: "SP", InternaPct: "18.00", InterestadualPct: "12", EfetivoPct: "6.00", OrigemVersao: "padrao-2026"}
	svc := newService(repo, fakeCost{}, nil)
	// Δ = 0.04 ≤ 0.049 ⇒ no persist.
	res, err := svc.PutDifalOverride(context.Background(), "SP", "18.04", "op", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if res.Persisted {
		t.Fatalf("Persisted = true, want false for Δ≤0.049")
	}
	if repo.overrideCall != nil {
		t.Fatalf("override was written, want no-persist")
	}
}

func TestPutDifalOverrideAboveThresholdPersists(t *testing.T) {
	repo := newFakeRepo()
	repo.rates["SP"] = domain.DifalRate{UF: "SP", InternaPct: "18.00", InterestadualPct: "12", EfetivoPct: "6.00", OrigemVersao: "padrao-2026"}
	svc := newService(repo, fakeCost{}, nil)
	// Δ = 0.05 > 0.049 ⇒ persist.
	res, err := svc.PutDifalOverride(context.Background(), "SP", "18.05", "op", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if !res.Persisted {
		t.Fatalf("Persisted = false, want true for Δ>0.049")
	}
	if repo.overrideCall == nil || repo.overrideCall.ov.InternaPct != "18.05" {
		t.Fatalf("override not written correctly: %#v", repo.overrideCall)
	}
}

func TestDecomposeInvalidPrice(t *testing.T) {
	svc := newService(newFakeRepo(), fakeCost{}, nil)
	_, err := svc.Decompose(context.Background(), application.DecomposeRequest{
		Preco: "0", ComissaoPct: "12", Modalidade: domain.ModalidadeClassico,
	})
	if !errors.Is(err, application.ErrInvalidPrice) {
		t.Fatalf("err = %v, want ErrInvalidPrice", err)
	}
}

func TestDecomposeItemNotFound(t *testing.T) {
	products := fakeProducts{exists: map[int]bool{}} // 90001 absent
	svc := newService(newFakeRepo(), fakeCost{}, products)
	id := 90001
	_, err := svc.Decompose(context.Background(), application.DecomposeRequest{
		Preco: "100.00", ComissaoPct: "12", Modalidade: domain.ModalidadeClassico, ProductID: &id,
	})
	if !errors.Is(err, application.ErrItemNotFound) {
		t.Fatalf("err = %v, want ErrItemNotFound", err)
	}
}

func TestDecomposeKnownProductNoCostBlocksSemCusto(t *testing.T) {
	products := fakeProducts{exists: map[int]bool{90001: true}}
	svc := newService(newFakeRepo(), fakeCost{money: nil}, products) // cost absent
	id := 90001
	res, err := svc.Decompose(context.Background(), application.DecomposeRequest{
		Preco: "100.00", ComissaoPct: "12", Modalidade: domain.ModalidadeClassico, ProductID: &id,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.BlockingState == nil || *res.BlockingState != application.BlockingStateSemCusto {
		t.Fatalf("blocking = %v, want SEM_CUSTO", res.BlockingState)
	}
	if res.Decomposition.Custo != nil {
		t.Fatalf("custo = %v, want nil (unknown, never 0)", res.Decomposition.Custo)
	}
	if !containsStr(res.Decomposition.ComponentesDesconhecidos, "custo_erp") {
		t.Fatalf("desconhecidos = %v, want custo_erp", res.Decomposition.ComponentesDesconhecidos)
	}
}

func TestDecomposeExplicitCostKnown(t *testing.T) {
	svc := newService(newFakeRepo(), fakeCost{}, nil)
	custo := "40.00"
	res, err := svc.Decompose(context.Background(), application.DecomposeRequest{
		Preco: "100.00", ComissaoPct: "12", Modalidade: domain.ModalidadeClassico, Custo: &custo,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.BlockingState != nil {
		t.Fatalf("blocking = %v, want nil", res.BlockingState)
	}
	if res.Decomposition.Custo == nil || *res.Decomposition.Custo != "40.00" {
		t.Fatalf("custo = %v, want 40.00", res.Decomposition.Custo)
	}
}

func TestSolveTargetUnreachableIsNotError(t *testing.T) {
	svc := newService(newFakeRepo(), fakeCost{}, nil)
	custo, frete := "10.00", "15.00"
	// Default profile: DIFAL off, aliquota 4. Ceiling = 100 - 16 - 4 = 80;
	// target 95% is above the ceiling ⇒ UNREACHABLE (200, not error).
	out, err := svc.SolveTarget(context.Background(), application.SolveRequest{
		TargetMargemPct: "95", ComissaoPct: "16", Modalidade: domain.ModalidadeClassico, Custo: &custo, FreteProduto: &frete,
	})
	if err != nil {
		t.Fatalf("unreachable must not error: %v", err)
	}
	if out.Result.Reached {
		t.Fatalf("Reached = true, want false for 95%%@comissão16 (ceiling 80)")
	}
	if out.Result.CeilingPct == "" {
		t.Fatalf("CeilingPct empty, want attainable ceiling cited")
	}
}

func TestSolveTargetReachable(t *testing.T) {
	svc := newService(newFakeRepo(), fakeCost{}, nil)
	custo, frete := "40.00", "15.00"
	out, err := svc.SolveTarget(context.Background(), application.SolveRequest{
		TargetMargemPct: "15", ComissaoPct: "12", Modalidade: domain.ModalidadeClassico, Custo: &custo, FreteProduto: &frete,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !out.Result.Reached || out.Result.Preco == nil {
		t.Fatalf("expected reachable price, got %#v", out.Result)
	}
}

func containsStr(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}
