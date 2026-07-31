package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"marketplace-central/apps/server_core/internal/modules/pricing/application"
	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
	"marketplace-central/apps/server_core/internal/platform/apierror"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

// WithCalc attaches the IC-04 calculator service, enabling the /pricing/profile,
// /pricing/difal, /pricing/scenarios, /pricing/decompose and /pricing/solve
// routes. Returns a copy so the base Handler (W1 simulations) is unchanged when
// the calc service is not wired.
func (h Handler) WithCalc(calc application.CalcService) Handler {
	h.calc = &calc
	return h
}

// registerCalc adds the IC-04 calculator routes. Called from Register only when
// a calc service is attached (WithCalc). W1 /pricing/simulations is untouched.
func (h Handler) registerCalc(mux httpx.RouteRegistrar) {
	mux.HandleFunc("GET /pricing/profile", h.handleGetProfile)
	mux.HandleFunc("PUT /pricing/profile", h.handlePutProfile)
	mux.HandleFunc("GET /pricing/difal", h.handleListDifal)
	mux.HandleFunc("PUT /pricing/difal/{uf}", h.handlePutDifal)
	mux.HandleFunc("GET /pricing/scenarios", h.handleListScenarios)
	mux.HandleFunc("POST /pricing/scenarios", h.handleCreateScenario)
	mux.HandleFunc("DELETE /pricing/scenarios/{id}", h.handleDeleteScenario)
	mux.HandleFunc("POST /pricing/decompose", h.handleDecompose)
	mux.HandleFunc("POST /pricing/solve", h.handleSolve)
	mux.HandleFunc("GET /pricing/tariff-defaults", h.handleGetTariffDefaults)
	mux.HandleFunc("PUT /pricing/tariff-defaults", h.handlePutTariffDefaults)
}

// mapCalcError maps an IC-04 service error to its exact HTTP status + code.
func mapCalcError(err error) (int, string) {
	switch {
	case errors.Is(err, application.ErrInvalidRate):
		return http.StatusUnprocessableEntity, "INVALID_RATE"
	case errors.Is(err, application.ErrInvalidPrice):
		return http.StatusUnprocessableEntity, "INVALID_PRICE"
	case errors.Is(err, application.ErrItemNotFound):
		return http.StatusNotFound, "ITEM_NOT_FOUND"
	case errors.Is(err, ports.ErrDifalUFNotFound):
		return http.StatusNotFound, "UF_NOT_FOUND"
	case errors.Is(err, ports.ErrScenarioNotFound):
		return http.StatusNotFound, "SCENARIO_NOT_FOUND"
	case errors.Is(err, domain.ErrInvalidFretePolicy):
		return http.StatusUnprocessableEntity, "INVALID_FRETE_POLICY"
	case errors.Is(err, application.ErrTariffStoreNotConfigured):
		return http.StatusServiceUnavailable, "PRICING_TARIFF_STORE_NOT_CONFIGURED"
	default:
		return http.StatusInternalServerError, "PRICING_INTERNAL_ERROR"
	}
}

func (h Handler) writeCalcError(w http.ResponseWriter, err error) {
	status, code := mapCalcError(err)
	msg := err.Error()
	if status >= 500 {
		msg = "internal error"
	}
	apierror.Write(w, status, code, msg, nil)
}

// --- Profile ---

type profileDTO struct {
	Regime           string  `json:"regime"`
	AliquotaPct      string  `json:"aliquota_pct"`
	LimiarVerdePct   string  `json:"limiar_verde_pct"`
	LimiarAmareloPct string  `json:"limiar_amarelo_pct"`
	TarifaFull       *string `json:"tarifa_full"`
	DifalEnabled     bool    `json:"difal_enabled"`
	DifalDestinoUF   *string `json:"difal_destino_uf"`
	Origem           string  `json:"origem"`
}

func toProfileDTO(p domain.CalcProfile) profileDTO {
	dto := profileDTO{
		Regime:           string(p.Regime),
		AliquotaPct:      p.AliquotaPct,
		LimiarVerdePct:   p.LimiarVerdePct,
		LimiarAmareloPct: p.LimiarAmareloPct,
		DifalEnabled:     p.DifalEnabled,
		DifalDestinoUF:   p.DifalDestinoUF,
		Origem:           p.Origem,
	}
	if p.TarifaFull != nil {
		dto.TarifaFull = &p.TarifaFull.Amount
	}
	return dto
}

func (h Handler) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := h.calc.GetProfile(r.Context())
	if err != nil {
		h.writeCalcError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toProfileDTO(profile))
}

func (h Handler) handlePutProfile(w http.ResponseWriter, r *http.Request) {
	var req profileDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierror.Write(w, http.StatusBadRequest, "PRICING_REQUEST_INVALID", "malformed request body", nil)
		return
	}
	profile, err := h.calc.PutProfile(r.Context(), application.ProfileUpdate{
		Regime:           domain.Regime(req.Regime),
		AliquotaPct:      req.AliquotaPct,
		LimiarVerdePct:   req.LimiarVerdePct,
		LimiarAmareloPct: req.LimiarAmareloPct,
		TarifaFull:       req.TarifaFull,
		DifalEnabled:     req.DifalEnabled,
		DifalDestinoUF:   req.DifalDestinoUF,
	})
	if err != nil {
		h.writeCalcError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toProfileDTO(profile))
}

// --- DIFAL ---

type overrideDTO struct {
	InternaPct string `json:"interna_pct"`
	UpdatedAt  string `json:"updated_at"`
	Actor      string `json:"actor"`
}

type difalRateDTO struct {
	UF               string       `json:"uf"`
	InternaPct       string       `json:"interna_pct"`
	InterestadualPct string       `json:"interestadual_pct"`
	EfetivoPct       string       `json:"efetivo_pct"`
	OrigemVersao     string       `json:"origem_versao"`
	Override         *overrideDTO `json:"override"`
}

// toDifalRateDTO renders a rate with efetivo_pct reflecting any active override
// (via domain.DifalForUF) — never the raw stored seed efetivo when overridden.
func toDifalRateDTO(rate domain.DifalRate) difalRateDTO {
	dto := difalRateDTO{
		UF:               rate.UF,
		InternaPct:       rate.InternaPct,
		InterestadualPct: rate.InterestadualPct,
		EfetivoPct:       domain.DifalForUF(rate).EfetivoPct,
		OrigemVersao:     rate.OrigemVersao,
	}
	if rate.Override != nil {
		dto.Override = &overrideDTO{InternaPct: rate.Override.InternaPct, UpdatedAt: rate.Override.UpdatedAt, Actor: rate.Override.Actor}
	}
	return dto
}

func (h Handler) handleListDifal(w http.ResponseWriter, r *http.Request) {
	rates, err := h.calc.ListDifal(r.Context())
	if err != nil {
		h.writeCalcError(w, err)
		return
	}
	items := make([]difalRateDTO, 0, len(rates))
	for _, rate := range rates {
		items = append(items, toDifalRateDTO(rate))
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"items":      items,
		"disclaimer": domain.DisclaimerSeedPadrao2026,
	})
}

func (h Handler) handlePutDifal(w http.ResponseWriter, r *http.Request) {
	uf := r.PathValue("uf")
	var req struct {
		InternaPct string `json:"interna_pct"`
		Actor      string `json:"actor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierror.Write(w, http.StatusBadRequest, "PRICING_REQUEST_INVALID", "malformed request body", nil)
		return
	}
	res, err := h.calc.PutDifalOverride(r.Context(), uf, req.InternaPct, req.Actor, time.Now())
	if err != nil {
		h.writeCalcError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"persisted": res.Persisted,
		"rate":      toDifalRateDTO(res.Rate),
	})
}

// --- Scenarios ---

type scenarioDTO struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt string          `json:"created_at"`
}

func toScenarioDTO(s domain.Scenario) scenarioDTO {
	return scenarioDTO{
		ID:        s.ID,
		Name:      s.Name,
		Payload:   json.RawMessage(s.Payload),
		CreatedAt: s.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func (h Handler) handleListScenarios(w http.ResponseWriter, r *http.Request) {
	scenarios, err := h.calc.ListScenarios(r.Context())
	if err != nil {
		h.writeCalcError(w, err)
		return
	}
	items := make([]scenarioDTO, 0, len(scenarios))
	for _, s := range scenarios {
		items = append(items, toScenarioDTO(s))
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h Handler) handleCreateScenario(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID      string          `json:"id"`
		Name    string          `json:"name"`
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierror.Write(w, http.StatusBadRequest, "PRICING_REQUEST_INVALID", "malformed request body", nil)
		return
	}
	scenario := domain.Scenario{ID: req.ID, Name: req.Name, Payload: string(req.Payload)}
	if err := h.calc.SaveScenario(r.Context(), scenario); err != nil {
		h.writeCalcError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, toScenarioDTO(scenario))
}

func (h Handler) handleDeleteScenario(w http.ResponseWriter, r *http.Request) {
	if err := h.calc.DeleteScenario(r.Context(), r.PathValue("id")); err != nil {
		h.writeCalcError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Decompose / Solve ---

type decompositionDTO struct {
	Preco                    string   `json:"preco"`
	Comissao                 string   `json:"comissao"`
	TaxaFixa                 string   `json:"taxa_fixa"`
	Frete                    *string  `json:"frete"`
	Imposto                  string   `json:"imposto"`
	Difal                    *string  `json:"difal"`
	TarifaFull               *string  `json:"tarifa_full"`
	Custo                    *string  `json:"custo"`
	MargemValor              *string  `json:"margem_valor"`
	MargemPct                *string  `json:"margem_pct"`
	ComponentesDesconhecidos []string `json:"componentes_desconhecidos"`
}

func toDecompositionDTO(d domain.Decomposition) decompositionDTO {
	unknown := d.ComponentesDesconhecidos
	if unknown == nil {
		unknown = []string{}
	}
	return decompositionDTO{
		Preco: d.Preco, Comissao: d.Comissao, TaxaFixa: d.TaxaFixa, Frete: d.Frete,
		Imposto: d.Imposto, Difal: d.Difal, TarifaFull: d.TarifaFull, Custo: d.Custo,
		MargemValor: d.MargemValor, MargemPct: d.MargemPct, ComponentesDesconhecidos: unknown,
	}
}

type calcInputDTO struct {
	Preco         string  `json:"preco"`
	MargemAlvoPct string  `json:"margem_alvo_pct"`
	ComissaoPct   string  `json:"comissao_pct"`
	Modalidade    string  `json:"modalidade"`
	TarifaFull    *string `json:"tarifa_full"`
	FreteProduto  *string `json:"frete_produto"`
	Custo         *string `json:"custo"`
	ProductID     *int    `json:"product_id"`
}

func (h Handler) handleDecompose(w http.ResponseWriter, r *http.Request) {
	var req calcInputDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierror.Write(w, http.StatusBadRequest, "PRICING_REQUEST_INVALID", "malformed request body", nil)
		return
	}
	res, err := h.calc.Decompose(r.Context(), application.DecomposeRequest{
		Preco:        req.Preco,
		ComissaoPct:  req.ComissaoPct,
		Modalidade:   domain.Modalidade(req.Modalidade),
		TarifaFull:   req.TarifaFull,
		FreteProduto: req.FreteProduto,
		Custo:        req.Custo,
		ProductID:    req.ProductID,
		AsOf:         time.Now(),
	})
	if err != nil {
		h.writeCalcError(w, err)
		return
	}
	body := map[string]any{
		"decomposition":  toDecompositionDTO(res.Decomposition),
		"blocking_state": res.BlockingState,
	}
	if tarifa := toTarifaDTO(res.Tarifa); tarifa != nil {
		body["tarifa"] = tarifa
	}
	httpx.WriteJSON(w, http.StatusOK, body)
}

func (h Handler) handleSolve(w http.ResponseWriter, r *http.Request) {
	var req calcInputDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierror.Write(w, http.StatusBadRequest, "PRICING_REQUEST_INVALID", "malformed request body", nil)
		return
	}
	out, err := h.calc.SolveTarget(r.Context(), application.SolveRequest{
		TargetMargemPct: req.MargemAlvoPct,
		ComissaoPct:     req.ComissaoPct,
		Modalidade:      domain.Modalidade(req.Modalidade),
		TarifaFull:      req.TarifaFull,
		FreteProduto:    req.FreteProduto,
		Custo:           req.Custo,
		ProductID:       req.ProductID,
		AsOf:            time.Now(),
	})
	if err != nil {
		h.writeCalcError(w, err)
		return
	}
	desconhecidos := out.Result.Desconhecidos
	if desconhecidos == nil {
		desconhecidos = []string{}
	}
	body := map[string]any{
		"reached":            out.Result.Reached,
		"preco":              out.Result.Preco,
		"ceiling_pct":        out.Result.CeilingPct,
		"desconhecidos":      desconhecidos,
		"frete_desconhecido": out.Result.FreteDesconhecido,
		"blocking_state":     out.BlockingState,
	}
	if tarifa := toTarifaDTO(out.Tarifa); tarifa != nil {
		body["tarifa"] = tarifa
	}
	if code := solveCode(out.Result); code != "" {
		body["code"] = code
	}
	httpx.WriteJSON(w, http.StatusOK, body)
}

type tarifaComissaoDTO struct {
	Valor      *string `json:"valor"`
	Fonte      string  `json:"fonte"`
	Degrau     int     `json:"degrau"`
	Data       *string `json:"data"`
	Estimativa bool    `json:"estimativa"`
}

type tarifaFreteDTO struct {
	Valor      *string `json:"valor"`
	Fonte      string  `json:"fonte"`
	Degrau     int     `json:"degrau"`
	Data       *string `json:"data"`
	Estimativa bool    `json:"estimativa"`
	SemDados   bool    `json:"sem_dados"`
}

type tarifaDTO struct {
	Comissao tarifaComissaoDTO `json:"comissao"`
	Frete    tarifaFreteDTO    `json:"frete"`
}

// toTarifaDTO renders the resolved tariff stamps. `data` (source timestamp) is
// null at degrau 4 — the config default has no authoritative per-resolve
// timestamp; fabricating one would violate ADR-17. Degrau 3 (COTACAO) carries a
// real timestamp in ComponentResolution.Data, mapped straight through here.
func toTarifaDTO(t *domain.TariffResolution) *tarifaDTO {
	if t == nil {
		return nil
	}
	return &tarifaDTO{
		Comissao: tarifaComissaoDTO{
			Valor: t.Comissao.Valor, Fonte: string(t.Comissao.Fonte),
			Degrau: t.Comissao.Degrau, Data: t.Comissao.Data, Estimativa: t.Comissao.Estimativa,
		},
		Frete: tarifaFreteDTO{
			Valor: t.Frete.Valor, Fonte: string(t.Frete.Fonte),
			Degrau: t.Frete.Degrau, Data: t.Frete.Data, Estimativa: t.Frete.Estimativa,
			SemDados: t.Frete.Valor == nil,
		},
	}
}

// solveCode maps a solve outcome to its IC-04 code by CAUSE — never a blanket
// UNREACHABLE_TARGET. Mutually exclusive in the current domain: structural
// unknowns, then missing frete (high segment), then true unreachability.
func solveCode(r domain.SolveResult) string {
	switch {
	case len(r.Desconhecidos) > 0:
		return "DADOS_INCOMPLETOS"
	case r.FreteDesconhecido:
		return "SEM_FRETE"
	case !r.Reached && r.CeilingPct != "":
		return "UNREACHABLE_TARGET"
	default:
		return ""
	}
}

// --- Tariff Defaults (CHIP-T1 Slice A) ---

type tariffDefaultsDTO struct {
	ComissaoClassicoPct   string  `json:"comissao_classico_pct"`
	ComissaoPremiumPct    string  `json:"comissao_premium_pct"`
	FreteEstimativaAmount *string `json:"frete_estimativa_amount"`
	FretePolicy           string  `json:"frete_policy"`
}

func toTariffDefaultsDTO(t domain.TariffDefaults) tariffDefaultsDTO {
	return tariffDefaultsDTO{
		ComissaoClassicoPct:   t.ComissaoClassicoPct,
		ComissaoPremiumPct:    t.ComissaoPremiumPct,
		FreteEstimativaAmount: t.FreteEstimativaAmount,
		FretePolicy:           t.FretePolicy,
	}
}

// tariffInstallationID reads the optional installation_id query parameter,
// defaulting to "" (the single-installation sentinel).
func tariffInstallationID(r *http.Request) string {
	return r.URL.Query().Get("installation_id")
}

func (h Handler) handleGetTariffDefaults(w http.ResponseWriter, r *http.Request) {
	defaults, err := h.calc.GetTariffDefaults(r.Context(), tariffInstallationID(r))
	if err != nil {
		h.writeCalcError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toTariffDefaultsDTO(defaults))
}

func (h Handler) handlePutTariffDefaults(w http.ResponseWriter, r *http.Request) {
	var req tariffDefaultsDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierror.Write(w, http.StatusBadRequest, "PRICING_REQUEST_INVALID", "malformed request body", nil)
		return
	}
	defaults, err := h.calc.PutTariffDefaults(r.Context(), tariffInstallationID(r), domain.TariffDefaults{
		ComissaoClassicoPct:   req.ComissaoClassicoPct,
		ComissaoPremiumPct:    req.ComissaoPremiumPct,
		FreteEstimativaAmount: req.FreteEstimativaAmount,
		FretePolicy:           req.FretePolicy,
	})
	if err != nil {
		h.writeCalcError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toTariffDefaultsDTO(defaults))
}
