package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/pricing/application"
	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// --- fakes ---

type calcRepoStub struct {
	profile   *domain.CalcProfile
	rates     map[string]domain.DifalRate
	scenarios map[string]domain.Scenario
}

func newRepoStub() *calcRepoStub {
	return &calcRepoStub{rates: map[string]domain.DifalRate{}, scenarios: map[string]domain.Scenario{}}
}

func (f *calcRepoStub) GetProfile(context.Context, string) (domain.CalcProfile, error) {
	if f.profile != nil {
		return *f.profile, nil
	}
	return domain.NewDefaultCalcProfile(), nil
}
func (f *calcRepoStub) UpsertProfile(_ context.Context, _ string, p domain.CalcProfile) error {
	f.profile = &p
	return nil
}
func (f *calcRepoStub) RateForUF(_ context.Context, _, uf string) (domain.DifalRate, error) {
	r, ok := f.rates[uf]
	if !ok {
		return domain.DifalRate{}, ports.ErrDifalUFNotFound
	}
	return r, nil
}
func (f *calcRepoStub) ListDifalRates(context.Context, string) ([]domain.DifalRate, error) {
	out := make([]domain.DifalRate, 0, len(f.rates))
	for _, r := range f.rates {
		out = append(out, r)
	}
	return out, nil
}
func (f *calcRepoStub) UpsertDifalOverride(_ context.Context, _, uf string, ov *domain.DifalOverride) error {
	r := f.rates[uf]
	r.Override = ov
	f.rates[uf] = r
	return nil
}
func (f *calcRepoStub) ListScenarios(context.Context, string) ([]domain.Scenario, error) {
	out := make([]domain.Scenario, 0, len(f.scenarios))
	for _, s := range f.scenarios {
		out = append(out, s)
	}
	return out, nil
}
func (f *calcRepoStub) SaveScenario(_ context.Context, _ string, s domain.Scenario) error {
	f.scenarios[s.ID] = s
	return nil
}
func (f *calcRepoStub) DeleteScenario(_ context.Context, _, id string) error {
	if _, ok := f.scenarios[id]; !ok {
		return ports.ErrScenarioNotFound
	}
	delete(f.scenarios, id)
	return nil
}

type costStub struct{ money *domain.Money }

func (c costStub) CostForProduct(context.Context, int, time.Time) (*domain.Money, error) {
	return c.money, nil
}

type productStub struct{ exists map[int]bool }

func (p productStub) Exists(_ context.Context, id int) (bool, error) { return p.exists[id], nil }

func newCalcMux(repo ports.CalcRepository, cost ports.CostReader, products ports.ProductChecker) *http.ServeMux {
	svc := application.NewCalcService(repo, cost, products, "t1")
	h := NewHandler(application.Service{}, nil).WithCalc(svc)
	mux := http.NewServeMux()
	h.Register(mux)
	return mux
}

func do(t *testing.T, mux *http.ServeMux, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == "" {
		r = httptest.NewRequest(method, path, nil)
	} else {
		r = httptest.NewRequest(method, path, strings.NewReader(body))
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w
}

func errCode(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode error body %q: %v", w.Body.String(), err)
	}
	return body.Error.Code
}

// --- error matrix ---

func TestPutProfileInvalidRate422(t *testing.T) {
	mux := newCalcMux(newRepoStub(), costStub{}, nil)
	w := do(t, mux, http.MethodPut, "/pricing/profile",
		`{"regime":"SIMPLES","aliquota_pct":"40","limiar_verde_pct":"18","limiar_amarelo_pct":"10"}`)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", w.Code)
	}
	if c := errCode(t, w); c != "INVALID_RATE" {
		t.Fatalf("code = %q, want INVALID_RATE", c)
	}
}

func TestPutDifalUFNotFound404(t *testing.T) {
	mux := newCalcMux(newRepoStub(), costStub{}, nil)
	w := do(t, mux, http.MethodPut, "/pricing/difal/ZZ", `{"interna_pct":"19","actor":"op"}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if c := errCode(t, w); c != "UF_NOT_FOUND" {
		t.Fatalf("code = %q, want UF_NOT_FOUND", c)
	}
}

func TestPutDifalInvalidRate422(t *testing.T) {
	repo := newRepoStub()
	repo.rates["SP"] = domain.DifalRate{UF: "SP", InternaPct: "18.00", InterestadualPct: "12", EfetivoPct: "6.00", OrigemVersao: "padrao-2026"}
	mux := newCalcMux(repo, costStub{}, nil)
	w := do(t, mux, http.MethodPut, "/pricing/difal/SP", `{"interna_pct":"40","actor":"op"}`)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", w.Code)
	}
	if c := errCode(t, w); c != "INVALID_RATE" {
		t.Fatalf("code = %q, want INVALID_RATE", c)
	}
}

func TestPutDifalBelowThreshold200NoPersist(t *testing.T) {
	repo := newRepoStub()
	repo.rates["SP"] = domain.DifalRate{UF: "SP", InternaPct: "18.00", InterestadualPct: "12", EfetivoPct: "6.00", OrigemVersao: "padrao-2026"}
	mux := newCalcMux(repo, costStub{}, nil)
	w := do(t, mux, http.MethodPut, "/pricing/difal/SP", `{"interna_pct":"18.04","actor":"op"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var body struct {
		Persisted bool `json:"persisted"`
	}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body.Persisted {
		t.Fatalf("persisted = true, want false for Δ≤0.049")
	}
	if repo.rates["SP"].Override != nil {
		t.Fatalf("override written, want no-persist")
	}
}

func TestDecomposeInvalidPrice422(t *testing.T) {
	mux := newCalcMux(newRepoStub(), costStub{}, nil)
	w := do(t, mux, http.MethodPost, "/pricing/decompose",
		`{"preco":"0","comissao_pct":"12","modalidade":"classico"}`)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", w.Code)
	}
	if c := errCode(t, w); c != "INVALID_PRICE" {
		t.Fatalf("code = %q, want INVALID_PRICE", c)
	}
}

func TestDecomposeItemNotFound404(t *testing.T) {
	mux := newCalcMux(newRepoStub(), costStub{}, productStub{exists: map[int]bool{}})
	w := do(t, mux, http.MethodPost, "/pricing/decompose",
		`{"preco":"100.00","comissao_pct":"12","modalidade":"classico","product_id":90001}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if c := errCode(t, w); c != "ITEM_NOT_FOUND" {
		t.Fatalf("code = %q, want ITEM_NOT_FOUND", c)
	}
}

func TestDecomposeNoCost200SemCusto(t *testing.T) {
	mux := newCalcMux(newRepoStub(), costStub{money: nil}, productStub{exists: map[int]bool{90001: true}})
	w := do(t, mux, http.MethodPost, "/pricing/decompose",
		`{"preco":"100.00","comissao_pct":"12","modalidade":"classico","product_id":90001}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (sem custo is not an error)", w.Code)
	}
	var body struct {
		BlockingState *string `json:"blocking_state"`
		Decomposition struct {
			Custo                    *string  `json:"custo"`
			ComponentesDesconhecidos []string `json:"componentes_desconhecidos"`
		} `json:"decomposition"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.BlockingState == nil || *body.BlockingState != "SEM_CUSTO" {
		t.Fatalf("blocking_state = %v, want SEM_CUSTO", body.BlockingState)
	}
	if body.Decomposition.Custo != nil {
		t.Fatalf("custo = %v, want null (never fabricated 0)", *body.Decomposition.Custo)
	}
}

func TestSolveUnreachable200Code(t *testing.T) {
	mux := newCalcMux(newRepoStub(), costStub{}, nil)
	// Default profile: DIFAL off, aliquota 4. Ceiling = 100-16-4 = 80; 95% > ceiling.
	w := do(t, mux, http.MethodPost, "/pricing/solve",
		`{"margem_alvo_pct":"95","comissao_pct":"16","modalidade":"classico","custo":"10.00","frete_produto":"15.00"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var body struct {
		Reached bool   `json:"reached"`
		Code    string `json:"code"`
	}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body.Reached {
		t.Fatalf("reached = true, want false")
	}
	if body.Code != "UNREACHABLE_TARGET" {
		t.Fatalf("code = %q, want UNREACHABLE_TARGET", body.Code)
	}
}

func TestDeleteScenarioNotFound404(t *testing.T) {
	mux := newCalcMux(newRepoStub(), costStub{}, nil)
	w := do(t, mux, http.MethodDelete, "/pricing/scenarios/missing", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if c := errCode(t, w); c != "SCENARIO_NOT_FOUND" {
		t.Fatalf("code = %q, want SCENARIO_NOT_FOUND", c)
	}
}

func TestListDifalHasDisclaimer(t *testing.T) {
	repo := newRepoStub()
	repo.rates["SP"] = domain.DifalRate{UF: "SP", InternaPct: "18.0", InterestadualPct: "12", EfetivoPct: "6.00", OrigemVersao: "padrao-2026"}
	mux := newCalcMux(repo, costStub{}, nil)
	w := do(t, mux, http.MethodGet, "/pricing/difal", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var body struct {
		Disclaimer string           `json:"disclaimer"`
		Items      []map[string]any `json:"items"`
	}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body.Disclaimer != domain.DisclaimerSeedPadrao2026 {
		t.Fatalf("disclaimer = %q, want mandatory seed disclaimer", body.Disclaimer)
	}
	if len(body.Items) != 1 || body.Items[0]["efetivo_pct"] != "6.00" {
		t.Fatalf("items = %#v", body.Items)
	}
}

// W1 regression: the pre-existing /pricing/simulations route is preserved
// (calc routes are additive). A wrong method still hits the W1 handler's 405.
func TestW1SimulationsRoutePreserved(t *testing.T) {
	mux := newCalcMux(newRepoStub(), costStub{}, nil)
	w := do(t, mux, http.MethodDelete, "/pricing/simulations", "")
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405 (W1 route still present)", w.Code)
	}
}
