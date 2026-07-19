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

// fakeTariffResolver is a configurable stub TariffResolver for CHIP-T1
// Slice D CalcService wiring tests.
type fakeTariffResolver struct {
	res domain.TariffResolution
	err error
	// gotReq, when non-nil, captures the last TariffRequest the service built
	// (degrau-3 threading assertions). The value receiver writes through it.
	gotReq *ports.TariffRequest
}

func (f fakeTariffResolver) Resolve(_ context.Context, req ports.TariffRequest) (domain.TariffResolution, error) {
	if f.gotReq != nil {
		*f.gotReq = req
	}
	return f.res, f.err
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

// TestDecomposeInvalidDecimalInputsReturnErrNotPanic asserts that empty or
// unparseable client decimal strings are rejected as ErrInvalidPrice (422)
// before the frozen no-error domain ports parse them — never a panic. Regression
// for the M-07 post-merge RED: POST /pricing/decompose {"preco":"100.00"} (empty
// commission) and POST /pricing/solve {} panicked in domain mustRat.
func TestDecomposeInvalidDecimalInputsReturnErrNotPanic(t *testing.T) {
	cases := []struct {
		name string
		req  application.DecomposeRequest
	}{
		{"empty preco", application.DecomposeRequest{Preco: "", ComissaoPct: "12", Modalidade: domain.ModalidadeClassico}},
		{"empty commission, valid price", application.DecomposeRequest{Preco: "100.00", ComissaoPct: "", Modalidade: domain.ModalidadeClassico}},
		{"empty preco and empty commission", application.DecomposeRequest{Preco: "", ComissaoPct: "", Modalidade: domain.ModalidadeClassico}},
		{"unparseable commission", application.DecomposeRequest{Preco: "100.00", ComissaoPct: "abc", Modalidade: domain.ModalidadeClassico}},
		{"empty tarifa_full override", application.DecomposeRequest{Preco: "100.00", ComissaoPct: "12", TarifaFull: strptr(""), Modalidade: domain.ModalidadeClassico}},
		{"empty custo override", application.DecomposeRequest{Preco: "100.00", ComissaoPct: "12", Custo: strptr(""), Modalidade: domain.ModalidadeClassico}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := newService(newFakeRepo(), fakeCost{}, nil)
			err := recoverErr(func() error {
				_, e := svc.Decompose(context.Background(), tc.req)
				return e
			})
			if !errors.Is(err, application.ErrInvalidPrice) {
				t.Fatalf("err = %v, want ErrInvalidPrice (never panic)", err)
			}
		})
	}
}

func TestSolveTargetInvalidDecimalInputsReturnErrNotPanic(t *testing.T) {
	cases := []struct {
		name string
		req  application.SolveRequest
	}{
		{"empty everything (POST /pricing/solve {})", application.SolveRequest{Modalidade: domain.ModalidadeClassico}},
		{"empty target margin, valid commission", application.SolveRequest{TargetMargemPct: "", ComissaoPct: "16", Modalidade: domain.ModalidadeClassico}},
		{"valid target, empty commission", application.SolveRequest{TargetMargemPct: "20", ComissaoPct: "", Modalidade: domain.ModalidadeClassico}},
		{"unparseable target", application.SolveRequest{TargetMargemPct: "x", ComissaoPct: "16", Modalidade: domain.ModalidadeClassico}},
		{"empty frete override", application.SolveRequest{TargetMargemPct: "20", ComissaoPct: "16", FreteProduto: strptr(""), Modalidade: domain.ModalidadeClassico}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := newService(newFakeRepo(), fakeCost{}, nil)
			err := recoverErr(func() error {
				_, e := svc.SolveTarget(context.Background(), tc.req)
				return e
			})
			if !errors.Is(err, application.ErrInvalidPrice) {
				t.Fatalf("err = %v, want ErrInvalidPrice (never panic)", err)
			}
		})
	}
}

// TestProfileSourcedUnparseableTarifaFullReturnsErrNotPanic covers the resolved
// (profile-sourced) tarifa_full path: PutProfile persists tarifa_full without a
// parse check, so a malformed stored value reaches the domain via
// resolveTarifaFull even when the request carries no override. Both Decompose
// and SolveTarget (modalidade=full) must reject it as ErrInvalidPrice, not panic.
func TestProfileSourcedUnparseableTarifaFullReturnsErrNotPanic(t *testing.T) {
	newRepoWithBadTarifa := func() *fakeCalcRepo {
		repo := newFakeRepo()
		repo.profile = &domain.CalcProfile{
			Regime: domain.RegimeSimples, AliquotaPct: "4", LimiarVerdePct: "18", LimiarAmareloPct: "10",
			TarifaFull: &domain.Money{Amount: "NaN", Currency: "BRL"}, Origem: "operator",
		}
		return repo
	}
	custo := "40.00"

	t.Run("decompose", func(t *testing.T) {
		svc := newService(newRepoWithBadTarifa(), fakeCost{}, nil)
		err := recoverErr(func() error {
			_, e := svc.Decompose(context.Background(), application.DecomposeRequest{
				Preco: "100.00", ComissaoPct: "12", Modalidade: domain.ModalidadeFull, Custo: &custo,
			})
			return e
		})
		if !errors.Is(err, application.ErrInvalidPrice) {
			t.Fatalf("err = %v, want ErrInvalidPrice (never panic)", err)
		}
	})

	t.Run("solve", func(t *testing.T) {
		svc := newService(newRepoWithBadTarifa(), fakeCost{}, nil)
		err := recoverErr(func() error {
			_, e := svc.SolveTarget(context.Background(), application.SolveRequest{
				TargetMargemPct: "15", ComissaoPct: "12", Modalidade: domain.ModalidadeFull, Custo: &custo,
			})
			return e
		})
		if !errors.Is(err, application.ErrInvalidPrice) {
			t.Fatalf("err = %v, want ErrInvalidPrice (never panic)", err)
		}
	})
}

// recoverErr runs fn and converts a panic into an error so the test asserts the
// no-panic contract explicitly rather than crashing the run.
func recoverErr(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("panicked")
		}
	}()
	return fn()
}

func strptr(s string) *string { return &s }

// --- CHIP-T1 Slice D: resolveTariff wiring ---

func TestDecomposeEmptyComissaoUsesResolverPadrao(t *testing.T) {
	resolver := fakeTariffResolver{res: domain.TariffResolution{
		Comissao: domain.ComponentResolution{Valor: strptr("13.00"), Fonte: domain.FontePadrao, Degrau: 4, Estimativa: true},
		Frete:    domain.ComponentResolution{Fonte: domain.FontePadrao, Degrau: 4, Estimativa: true},
	}}
	svc := newService(newFakeRepo(), fakeCost{}, nil).WithTariffResolver(resolver)
	custo := "40.00"
	res, err := svc.Decompose(context.Background(), application.DecomposeRequest{
		Preco: "100.00", ComissaoPct: "", Modalidade: domain.ModalidadeClassico, Custo: &custo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Tarifa == nil {
		t.Fatalf("Tarifa = nil, want resolved tariff")
	}
	if res.Tarifa.Comissao.Valor == nil || *res.Tarifa.Comissao.Valor != "13.00" {
		t.Fatalf("Tarifa.Comissao.Valor = %v, want 13.00", res.Tarifa.Comissao.Valor)
	}
	if res.Tarifa.Comissao.Fonte != domain.FontePadrao {
		t.Fatalf("Tarifa.Comissao.Fonte = %v, want FontePadrao", res.Tarifa.Comissao.Fonte)
	}
	if !res.Tarifa.Comissao.Estimativa {
		t.Fatalf("Tarifa.Comissao.Estimativa = false, want true")
	}
	// preço 100.00 @ comissão 13.00% ⇒ comissão amount 13.00; confirms the
	// resolver's pct (not an empty/zero) reached the domain engine.
	if res.Decomposition.Comissao != "13.00" {
		t.Fatalf("decomposition Comissao = %v, want 13.00 (resolver pct applied)", res.Decomposition.Comissao)
	}
}

func TestDecomposeManualComissaoOverridesResolver(t *testing.T) {
	resolver := fakeTariffResolver{res: domain.TariffResolution{
		Comissao: domain.ComponentResolution{Valor: strptr("13.00"), Fonte: domain.FontePadrao, Degrau: 4, Estimativa: true},
		Frete:    domain.ComponentResolution{Fonte: domain.FontePadrao, Degrau: 4, Estimativa: true},
	}}
	svc := newService(newFakeRepo(), fakeCost{}, nil).WithTariffResolver(resolver)
	custo := "40.00"
	res, err := svc.Decompose(context.Background(), application.DecomposeRequest{
		Preco: "100.00", ComissaoPct: "10", Modalidade: domain.ModalidadeClassico, Custo: &custo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Tarifa == nil {
		t.Fatalf("Tarifa = nil, want resolved tariff")
	}
	if res.Tarifa.Comissao.Valor == nil || *res.Tarifa.Comissao.Valor != "10" {
		t.Fatalf("Tarifa.Comissao.Valor = %v, want 10 (manual)", res.Tarifa.Comissao.Valor)
	}
	if res.Tarifa.Comissao.Fonte != domain.FonteManual {
		t.Fatalf("Tarifa.Comissao.Fonte = %v, want FonteManual", res.Tarifa.Comissao.Fonte)
	}
	// preço 100.00 @ comissão 10% ⇒ comissão amount 10.00; confirms the
	// manual override (not the resolver's 13.00) reached the domain engine.
	if res.Decomposition.Comissao != "10.00" {
		t.Fatalf("decomposition Comissao = %v, want 10.00 (manual pct applied)", res.Decomposition.Comissao)
	}
}

func TestDecomposeResolverFreteSemDadosStaysNil(t *testing.T) {
	resolver := fakeTariffResolver{res: domain.TariffResolution{
		Comissao: domain.ComponentResolution{Valor: strptr("13.00"), Fonte: domain.FontePadrao, Degrau: 4},
		Frete:    domain.ComponentResolution{Valor: nil, Fonte: domain.FontePadrao, Degrau: 4},
	}}
	svc := newService(newFakeRepo(), fakeCost{}, nil).WithTariffResolver(resolver)
	custo := "40.00"
	res, err := svc.Decompose(context.Background(), application.DecomposeRequest{
		Preco: "100.00", ComissaoPct: "", Modalidade: domain.ModalidadeClassico, Custo: &custo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Tarifa == nil || res.Tarifa.Frete.Valor != nil {
		t.Fatalf("Tarifa.Frete.Valor = %v, want nil (sem_dados, ADR-17: never 0)", res.Tarifa.Frete.Valor)
	}
	if res.Decomposition.Frete != nil {
		t.Fatalf("Decomposition.Frete = %v, want nil (sem_dados, unknown, never 0)", res.Decomposition.Frete)
	}
}

func TestSolveTargetResolverFreteSemDadosSurfacesFreteDesconhecido(t *testing.T) {
	resolver := fakeTariffResolver{res: domain.TariffResolution{
		Comissao: domain.ComponentResolution{Valor: strptr("12"), Fonte: domain.FontePadrao, Degrau: 4},
		Frete:    domain.ComponentResolution{Valor: nil, Fonte: domain.FontePadrao, Degrau: 4},
	}}
	svc := newService(newFakeRepo(), fakeCost{}, nil).WithTariffResolver(resolver)
	custo := "40.00"
	// High target margem (above what the low segment, frete=0, can reach given
	// custo 40.00) is only reachable in the ≥limiar (produto frete required)
	// segment; with frete unknown from the resolver and no request override,
	// SolveTarget must surface FreteDesconhecido rather than fabricate 0.
	out, err := svc.SolveTarget(context.Background(), application.SolveRequest{
		TargetMargemPct: "50", Modalidade: domain.ModalidadeClassico, Custo: &custo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Tarifa == nil || out.Tarifa.Frete.Valor != nil {
		t.Fatalf("Tarifa.Frete.Valor = %v, want nil (sem_dados)", out.Tarifa.Frete.Valor)
	}
	if !out.Result.FreteDesconhecido {
		t.Fatalf("Result.FreteDesconhecido = false, want true (frete required, sem_dados)")
	}
}

func TestDecomposeNoResolverEmptyComissaoStaysInvalidPrice(t *testing.T) {
	svc := newService(newFakeRepo(), fakeCost{}, nil) // no WithTariffResolver
	custo := "40.00"
	res, err := svc.Decompose(context.Background(), application.DecomposeRequest{
		Preco: "100.00", ComissaoPct: "", Modalidade: domain.ModalidadeClassico, Custo: &custo,
	})
	if !errors.Is(err, application.ErrInvalidPrice) {
		t.Fatalf("err = %v, want ErrInvalidPrice (back-compat: no resolver, empty comissao)", err)
	}
	if res.Tarifa != nil {
		t.Fatalf("Tarifa = %v, want nil on error", res.Tarifa)
	}
}

func intptr(i int) *int { return &i }

// Decompose threads the product id and the decompose price into the tariff
// request so degrau 3 can quote a live commission at that price.
func TestDecomposeThreadsDegrau3Inputs(t *testing.T) {
	var got ports.TariffRequest
	resolver := fakeTariffResolver{res: domain.TariffResolution{Comissao: domain.ComponentResolution{Valor: strptr("12.00")}}, gotReq: &got}
	svc := newService(newFakeRepo(), fakeCost{}, nil).WithTariffResolver(resolver)
	custo := "40.00"
	_, err := svc.Decompose(context.Background(), application.DecomposeRequest{
		Preco: "150.00", ComissaoPct: "", Modalidade: domain.ModalidadeClassico, ProductID: intptr(7), Custo: &custo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ProductID == nil || *got.ProductID != 7 {
		t.Fatalf("ProductID not threaded: %v", got.ProductID)
	}
	if got.PriceBasis == nil || *got.PriceBasis != "150.00" {
		t.Fatalf("PriceBasis not threaded: %v", got.PriceBasis)
	}
}

// A manual commission override wins outright, so degrau 3 is gated off: the
// service must not offer a product id/price to the resolver (no wasted live
// probe, and the override cannot be second-guessed by a quote).
func TestDecomposeManualCommissionGatesDegrau3(t *testing.T) {
	var got ports.TariffRequest
	resolver := fakeTariffResolver{res: domain.TariffResolution{Comissao: domain.ComponentResolution{Valor: strptr("12.00")}}, gotReq: &got}
	svc := newService(newFakeRepo(), fakeCost{}, nil).WithTariffResolver(resolver)
	custo := "40.00"
	_, err := svc.Decompose(context.Background(), application.DecomposeRequest{
		Preco: "150.00", ComissaoPct: "9.50", Modalidade: domain.ModalidadeClassico, ProductID: intptr(7), Custo: &custo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ProductID != nil {
		t.Fatalf("manual commission must gate degrau 3: ProductID=%v", got.ProductID)
	}
	if got.PriceBasis != nil {
		t.Fatalf("manual commission must gate degrau 3: PriceBasis=%v", got.PriceBasis)
	}
}

// Solve carries the product id but no price basis (it is solving FOR the
// price); degrau 3 falls back to the product's catalog price downstream.
func TestSolveThreadsProductIDNoPriceBasis(t *testing.T) {
	var got ports.TariffRequest
	resolver := fakeTariffResolver{res: domain.TariffResolution{Comissao: domain.ComponentResolution{Valor: strptr("12.00")}}, gotReq: &got}
	svc := newService(newFakeRepo(), fakeCost{}, nil).WithTariffResolver(resolver)
	custo := "40.00"
	_, err := svc.SolveTarget(context.Background(), application.SolveRequest{
		TargetMargemPct: "50", ComissaoPct: "", Modalidade: domain.ModalidadeClassico, ProductID: intptr(7), Custo: &custo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ProductID == nil || *got.ProductID != 7 {
		t.Fatalf("ProductID not threaded: %v", got.ProductID)
	}
	if got.PriceBasis != nil {
		t.Fatalf("Solve has no target price basis: %v", got.PriceBasis)
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
