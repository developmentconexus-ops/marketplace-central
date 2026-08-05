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

// tariffStoreStub is an in-memory ports.TariffDefaultsStore for handler
// tests. A nil rows map means "not yet materialized" — Get/Upsert seed the
// DB-default row (13.00/16.00/sem_dados), mirroring the real repo's
// materialize-on-read behavior.
type tariffStoreStub struct {
	rows map[string]domain.TariffDefaults
}

func newTariffStoreStub() *tariffStoreStub {
	return &tariffStoreStub{rows: map[string]domain.TariffDefaults{}}
}

func tariffKey(tenantID, installationID string) string { return tenantID + "/" + installationID }

func (f *tariffStoreStub) GetTariffDefaults(_ context.Context, tenantID, installationID string) (domain.TariffDefaults, error) {
	key := tariffKey(tenantID, installationID)
	if row, ok := f.rows[key]; ok {
		return row, nil
	}
	row := domain.TariffDefaults{ComissaoClassicoPct: "13.00", ComissaoPremiumPct: "16.00", FretePolicy: domain.FretePolicySemDados}
	f.rows[key] = row
	return row, nil
}

func (f *tariffStoreStub) UpsertTariffDefaults(_ context.Context, tenantID, installationID string, in domain.TariffDefaults) (domain.TariffDefaults, error) {
	if !domain.ValidFretePolicy(in.FretePolicy) {
		return domain.TariffDefaults{}, domain.ErrInvalidFretePolicy
	}
	f.rows[tariffKey(tenantID, installationID)] = in
	return in, nil
}

func newCalcMuxWithTariffStore(repo ports.CalcRepository, tariffStore ports.TariffDefaultsStore) *http.ServeMux {
	svc := application.NewCalcService(repo, costStub{}, nil, "t1").WithTariffStore(tariffStore)
	h := NewHandler(application.Service{}, nil).WithCalc(svc)
	mux := http.NewServeMux()
	h.Register(mux)
	return mux
}

// fakeTariffResolver is a handler-test-local ports.TariffResolver stub —
// mirrors application.fakeTariffResolver (unexported there, not reusable
// across packages).
type fakeTariffResolver struct {
	res domain.TariffResolution
}

func (f fakeTariffResolver) Resolve(context.Context, ports.TariffRequest) (domain.TariffResolution, error) {
	return f.res, nil
}

func newCalcMuxWithResolver(repo ports.CalcRepository, cost ports.CostReader, products ports.ProductChecker, resolver ports.TariffResolver) *http.ServeMux {
	svc := application.NewCalcService(repo, cost, products, "t1").WithTariffResolver(resolver)
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

// --- tariff defaults (CHIP-T1 Slice A) ---

func TestGetTariffDefaultsNotConfigured503(t *testing.T) {
	mux := newCalcMux(newRepoStub(), costStub{}, nil)
	w := do(t, mux, http.MethodGet, "/pricing/tariff-defaults", "")
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", w.Code)
	}
	if c := errCode(t, w); c != "PRICING_TARIFF_STORE_NOT_CONFIGURED" {
		t.Fatalf("code = %q, want PRICING_TARIFF_STORE_NOT_CONFIGURED", c)
	}
}

func TestGetTariffDefaultsMaterializesDBDefaults(t *testing.T) {
	mux := newCalcMuxWithTariffStore(newRepoStub(), newTariffStoreStub())
	w := do(t, mux, http.MethodGet, "/pricing/tariff-defaults", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", w.Code, w.Body.String())
	}
	var got tariffDefaultsDTO
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.ComissaoClassicoPct != "13.00" || got.ComissaoPremiumPct != "16.00" || got.FretePolicy != "sem_dados" {
		t.Fatalf("defaults = %#v", got)
	}
	if got.FreteEstimativaAmount != nil {
		t.Fatalf("frete_estimativa_amount = %v, want nil (ADR-17)", *got.FreteEstimativaAmount)
	}
}

func TestPutTariffDefaultsInvalidFretePolicy422(t *testing.T) {
	mux := newCalcMuxWithTariffStore(newRepoStub(), newTariffStoreStub())
	w := do(t, mux, http.MethodPut, "/pricing/tariff-defaults",
		`{"comissao_classico_pct":"13.00","comissao_premium_pct":"16.00","frete_policy":"invalida"}`)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422, body=%s", w.Code, w.Body.String())
	}
	if c := errCode(t, w); c != "INVALID_FRETE_POLICY" {
		t.Fatalf("code = %q, want INVALID_FRETE_POLICY", c)
	}
}

// A non-decimal numeric field must be a 422 INVALID_PRICE (guarded before the
// DB numeric cast), never a 500 PRICING_INTERNAL_ERROR.
func TestPutTariffDefaultsNonNumeric422(t *testing.T) {
	mux := newCalcMuxWithTariffStore(newRepoStub(), newTariffStoreStub())
	w := do(t, mux, http.MethodPut, "/pricing/tariff-defaults",
		`{"comissao_classico_pct":"abc","comissao_premium_pct":"16.00","frete_policy":"sem_dados"}`)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422, body=%s", w.Code, w.Body.String())
	}
	if c := errCode(t, w); c != "INVALID_PRICE" {
		t.Fatalf("code = %q, want INVALID_PRICE", c)
	}
}

func TestPutTariffDefaultsRoundTrip(t *testing.T) {
	mux := newCalcMuxWithTariffStore(newRepoStub(), newTariffStoreStub())
	w := do(t, mux, http.MethodPut, "/pricing/tariff-defaults",
		`{"comissao_classico_pct":"14.00","comissao_premium_pct":"17.00","frete_estimativa_amount":"22.50","frete_policy":"estimativa"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", w.Code, w.Body.String())
	}
	var got tariffDefaultsDTO
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.ComissaoClassicoPct != "14.00" || got.ComissaoPremiumPct != "17.00" || got.FretePolicy != "estimativa" {
		t.Fatalf("put response = %#v", got)
	}
	if got.FreteEstimativaAmount == nil || *got.FreteEstimativaAmount != "22.50" {
		t.Fatalf("frete_estimativa_amount = %v, want 22.50", got.FreteEstimativaAmount)
	}

	w2 := do(t, mux, http.MethodGet, "/pricing/tariff-defaults", "")
	var got2 tariffDefaultsDTO
	if err := json.Unmarshal(w2.Body.Bytes(), &got2); err != nil {
		t.Fatal(err)
	}
	if got2.ComissaoClassicoPct != got.ComissaoClassicoPct ||
		got2.ComissaoPremiumPct != got.ComissaoPremiumPct ||
		got2.FretePolicy != got.FretePolicy ||
		got2.FreteEstimativaAmount == nil || got.FreteEstimativaAmount == nil ||
		*got2.FreteEstimativaAmount != *got.FreteEstimativaAmount {
		t.Fatalf("GET after PUT = %#v, want %#v", got2, got)
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

// Task A7: imposto's *string widen must not touch the legacy value on the
// path that still produces one. application.CalcService does not wire
// ICMSCell yet (that wiring is Fatia B, out of scope here — DecomposeInput
// literal in calc_service.go never sets it), so every real HTTP decompose
// today takes the ICMSCell==nil branch and imposto stays the legacy
// regime-aliquota number. This is the negative control demanded by the task
// brief: if it breaks, the A7 cross-cut killed the legacy path wholesale
// instead of just widening the pointer.
func TestDecomposeImpostoLegacyPresentNoCell(t *testing.T) {
	mux := newCalcMux(newRepoStub(), costStub{}, nil)
	w := do(t, mux, http.MethodPost, "/pricing/decompose",
		`{"preco":"100.00","comissao_pct":"12","modalidade":"classico","custo":"10.00","frete_produto":"15.00"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Decomposition struct {
			Imposto *string `json:"imposto"`
		} `json:"decomposition"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Decomposition.Imposto == nil || *body.Decomposition.Imposto != "4.00" {
		t.Fatalf("imposto = %v, want \"4.00\" (default SIMPLES aliquota — legacy path unaffected)", body.Decomposition.Imposto)
	}
}

// TestDecompositionDTOImpostoNullWithICMSCell proves the transport DTO/JSON
// contract site (calc_handler.go,282): when domain.Decompose runs the
// D-41 ICMSCell path, Imposto is the named-absent nil (Task A6), and
// toDecompositionDTO/json.Marshal must carry that through as JSON null —
// never "0.00" and never "". application.CalcService does not resolve/pass
// an ICMSCell yet (Fatia B wiring, explicitly out of this task's scope per
// the brief), so this drives domain.Decompose directly with the same fixture
// as domain/decompose_icms_golden_test.go's "ICMS-1" case (read-only reuse —
// no edits to pricing/domain in this task) rather than the full HTTP path.
func TestDecompositionDTOImpostoNullWithICMSCell(t *testing.T) {
	in := domain.DecomposeInput{
		Preco: "1000.00", ComissaoPct: "12", AliquotaPct: "4",
		Modalidade:   domain.ModalidadeClassico,
		Custo:        &domain.Money{Amount: "200.00", Currency: "BRL"},
		FreteProduto: &domain.Money{Amount: "30.00", Currency: "BRL"},
		ICMSCell: &domain.ICMSCell{
			UFDestino: "MG", CodTrib: intPtr(0), Ambiguo: false,
			Origprod: intPtr(0), AliquotaInterna: strPtr("18"),
			FcpEmbutido:     strPtr("0"),
			RestituicaoUnit: strPtr("50.00"),
		},
	}
	decomp := domain.Decompose(in)
	if decomp.Imposto != nil {
		t.Fatalf("precondition: domain.Decompose returned Imposto = %q, want nil (Task A6 ICMSCell path)", *decomp.Imposto)
	}

	raw, err := json.Marshal(toDecompositionDTO(decomp))
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	val, present := got["imposto"]
	if !present {
		t.Fatalf("imposto key absent from JSON, want present with null value")
	}
	if val != nil {
		t.Fatalf("imposto = %v (%T), want JSON null (ICMSCell path, D-41: tax is apurado per fiscal cell; this DTO publishes only difal of the successors, B4 adds the rest)", val, val)
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

// --- CHIP-T1 Slice E: solve code CAUSE-branching (never a blanket
// UNREACHABLE_TARGET) + tarifa stamp surfacing ---

func TestSolveDadosIncompletos200Code(t *testing.T) {
	mux := newCalcMux(newRepoStub(), costStub{}, nil)
	// No custo, no product_id: custo_erp is a structural unknown ⇒
	// Desconhecidos non-empty ⇒ DADOS_INCOMPLETOS, never UNREACHABLE_TARGET.
	w := do(t, mux, http.MethodPost, "/pricing/solve",
		`{"margem_alvo_pct":"50","comissao_pct":"16","modalidade":"classico"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Code              string   `json:"code"`
		Desconhecidos     []string `json:"desconhecidos"`
		FreteDesconhecido bool     `json:"frete_desconhecido"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Code != "DADOS_INCOMPLETOS" {
		t.Fatalf("code = %q, want DADOS_INCOMPLETOS", body.Code)
	}
	if len(body.Desconhecidos) == 0 {
		t.Fatalf("desconhecidos = %v, want non-empty", body.Desconhecidos)
	}
	if body.FreteDesconhecido {
		t.Fatalf("frete_desconhecido = true, want false")
	}
}

func TestSolveSemFrete200Code(t *testing.T) {
	mux := newCalcMux(newRepoStub(), costStub{}, nil)
	// custo known, no frete_produto: target margem (65%) exceeds what the low
	// segment (frete=0, custo=10.00) can reach, forcing the high (>=limiar)
	// segment where produto frete is required — but unknown here (no resolver,
	// no request override) ⇒ FreteDesconhecido ⇒ SEM_FRETE.
	w := do(t, mux, http.MethodPost, "/pricing/solve",
		`{"margem_alvo_pct":"65","comissao_pct":"16","modalidade":"classico","custo":"10.00"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Code              string `json:"code"`
		FreteDesconhecido bool   `json:"frete_desconhecido"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Code != "SEM_FRETE" {
		t.Fatalf("code = %q, want SEM_FRETE", body.Code)
	}
	if !body.FreteDesconhecido {
		t.Fatalf("frete_desconhecido = false, want true")
	}
}

func TestSolveReached200NoCode(t *testing.T) {
	mux := newCalcMux(newRepoStub(), costStub{}, nil)
	w := do(t, mux, http.MethodPost, "/pricing/solve",
		`{"margem_alvo_pct":"20","comissao_pct":"16","modalidade":"classico","custo":"10.00","frete_produto":"15.00"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if reached, _ := body["reached"].(bool); !reached {
		t.Fatalf("reached = %v, want true, body=%s", body["reached"], w.Body.String())
	}
	if _, hasCode := body["code"]; hasCode {
		t.Fatalf("code key present = %v, want absent when reached", body["code"])
	}
}

func TestSolveWithResolverSurfacesTarifa(t *testing.T) {
	resolver := fakeTariffResolver{res: domain.TariffResolution{
		Comissao: domain.ComponentResolution{Valor: strPtr("13.00"), Fonte: domain.FontePadrao, Degrau: 4, Estimativa: true},
		Frete:    domain.ComponentResolution{Valor: nil, Fonte: domain.FontePadrao, Degrau: 4, Estimativa: true},
	}}
	custo := "10.00"
	mux := newCalcMuxWithResolver(newRepoStub(), costStub{money: &domain.Money{Amount: custo, Currency: "BRL"}}, nil, resolver)
	w := do(t, mux, http.MethodPost, "/pricing/solve",
		`{"margem_alvo_pct":"20","modalidade":"classico","custo":"10.00"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Tarifa *tarifaDTO `json:"tarifa"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Tarifa == nil {
		t.Fatalf("tarifa = nil, want resolved tarifa object, body=%s", w.Body.String())
	}
	if body.Tarifa.Comissao.Valor == nil || *body.Tarifa.Comissao.Valor != "13.00" {
		t.Fatalf("tarifa.comissao.valor = %v, want 13.00", body.Tarifa.Comissao.Valor)
	}
	if body.Tarifa.Comissao.Data != nil {
		t.Fatalf("tarifa.comissao.data = %v, want null (degrau 4 has no source timestamp)", *body.Tarifa.Comissao.Data)
	}
	if !body.Tarifa.Frete.SemDados {
		t.Fatalf("tarifa.frete.sem_dados = false, want true (resolver frete unknown)")
	}
}

func TestDecomposeDegrau3StampsCotacaoData(t *testing.T) {
	ts := "2026-07-19T12:00:00Z"
	resolver := fakeTariffResolver{res: domain.TariffResolution{
		Comissao: domain.ComponentResolution{Valor: strPtr("11.50"), Fonte: domain.FonteCotacao, Degrau: 3, Data: strPtr(ts)},
		Frete:    domain.ComponentResolution{Valor: nil, Fonte: domain.FontePadrao, Degrau: 4, Estimativa: true},
	}}
	mux := newCalcMuxWithResolver(newRepoStub(), costStub{money: &domain.Money{Amount: "10.00", Currency: "BRL"}}, nil, resolver)
	w := do(t, mux, http.MethodPost, "/pricing/decompose",
		`{"preco":"100.00","modalidade":"classico","custo":"10.00"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Tarifa *tarifaDTO `json:"tarifa"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Tarifa == nil {
		t.Fatalf("tarifa = nil, want resolved tarifa object, body=%s", w.Body.String())
	}
	if body.Tarifa.Comissao.Fonte != "COTACAO" {
		t.Fatalf("tarifa.comissao.fonte = %q, want COTACAO", body.Tarifa.Comissao.Fonte)
	}
	if body.Tarifa.Comissao.Degrau != 3 {
		t.Fatalf("tarifa.comissao.degrau = %d, want 3", body.Tarifa.Comissao.Degrau)
	}
	if body.Tarifa.Comissao.Data == nil || *body.Tarifa.Comissao.Data != ts {
		t.Fatalf("tarifa.comissao.data = %v, want %q (degrau 3 stamps the cotação timestamp)", body.Tarifa.Comissao.Data, ts)
	}
}

func TestDecomposeWithResolverSurfacesTarifa(t *testing.T) {
	resolver := fakeTariffResolver{res: domain.TariffResolution{
		Comissao: domain.ComponentResolution{Valor: strPtr("13.00"), Fonte: domain.FontePadrao, Degrau: 4},
		Frete:    domain.ComponentResolution{Valor: strPtr("22.50"), Fonte: domain.FontePadrao, Degrau: 4},
	}}
	custo := "10.00"
	mux := newCalcMuxWithResolver(newRepoStub(), costStub{money: &domain.Money{Amount: custo, Currency: "BRL"}}, nil, resolver)
	w := do(t, mux, http.MethodPost, "/pricing/decompose",
		`{"preco":"100.00","modalidade":"classico","custo":"10.00"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Tarifa *tarifaDTO `json:"tarifa"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Tarifa == nil {
		t.Fatalf("tarifa = nil, want resolved tarifa object, body=%s", w.Body.String())
	}
	if body.Tarifa.Comissao.Valor == nil || *body.Tarifa.Comissao.Valor != "13.00" {
		t.Fatalf("tarifa.comissao.valor = %v, want 13.00", body.Tarifa.Comissao.Valor)
	}
	if body.Tarifa.Frete.SemDados {
		t.Fatalf("tarifa.frete.sem_dados = true, want false (resolver frete known)")
	}
}

func strPtr(s string) *string { return &s }

func intPtr(i int) *int { return &i }

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
