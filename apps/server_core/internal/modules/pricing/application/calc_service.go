package application

import (
	"context"
	"errors"
	"math/big"
	"time"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// IC-04 service-layer error sentinels. Transport maps each to its exact HTTP
// code (INVALID_RATE/INVALID_PRICE ⇒ 422; ITEM_NOT_FOUND ⇒ 404). UF and
// scenario not-found reuse the ports sentinels.
var (
	ErrInvalidRate  = errors.New("INVALID_RATE")
	ErrInvalidPrice = errors.New("INVALID_PRICE")
	ErrItemNotFound = errors.New("ITEM_NOT_FOUND")
)

// BlockingStateSemCusto is the IC-03/IC-04 blocking_state emitted when a
// decomposition is computed for a known product whose custo_erp is absent —
// the response is a 200 partial decomposition, never an error (ADR-17).
const BlockingStateSemCusto = "SEM_CUSTO"

// overrideThresholdPP is the DIFAL override audit gate: an override persists
// only when |Δ interna_pct| exceeds 0,049 percentage points. At or below the
// threshold the write is skipped and the endpoint returns 200 without
// persisting.
var overrideThresholdPP = big.NewRat(49, 1000) // 0.049

// CalcService orchestrates the IC-04 calculator: the pure domain engine
// (Decompose/SolveTargetPrice) over the persistence repository, the decimal
// cost port, and the product-existence check. All money/percent values are
// decimal strings; the service never introduces float64.
type CalcService struct {
	repo     ports.CalcRepository
	cost     ports.CostReader
	products ports.ProductChecker
	tenantID string
}

func NewCalcService(repo ports.CalcRepository, cost ports.CostReader, products ports.ProductChecker, tenantID string) CalcService {
	return CalcService{repo: repo, cost: cost, products: products, tenantID: tenantID}
}

// --- Profile ---

func (s CalcService) GetProfile(ctx context.Context) (domain.CalcProfile, error) {
	return s.repo.GetProfile(ctx, s.tenantID)
}

// ProfileUpdate is a validated profile edit. Percent/money fields are decimal
// strings; TarifaFull nil clears tarifa_full (unset, never 0).
type ProfileUpdate struct {
	Regime           domain.Regime
	AliquotaPct      string
	LimiarVerdePct   string
	LimiarAmareloPct string
	TarifaFull       *string
	DifalEnabled     bool
	DifalDestinoUF   *string
}

// PutProfile validates and persists a profile edit. aliquota_pct outside
// 0..35 ⇒ ErrInvalidRate (422). Origem is stamped "operator" (an edit is by
// definition a tenant override, never the built-in default).
func (s CalcService) PutProfile(ctx context.Context, in ProfileUpdate) (domain.CalcProfile, error) {
	if !rateInRange(in.AliquotaPct) {
		return domain.CalcProfile{}, ErrInvalidRate
	}
	profile := domain.CalcProfile{
		Regime:           in.Regime,
		AliquotaPct:      in.AliquotaPct,
		LimiarVerdePct:   in.LimiarVerdePct,
		LimiarAmareloPct: in.LimiarAmareloPct,
		DifalEnabled:     in.DifalEnabled,
		DifalDestinoUF:   in.DifalDestinoUF,
		Origem:           "operator",
	}
	if in.TarifaFull != nil {
		profile.TarifaFull = &domain.Money{Amount: *in.TarifaFull, Currency: "BRL"}
	}
	if err := s.repo.UpsertProfile(ctx, s.tenantID, profile); err != nil {
		return domain.CalcProfile{}, err
	}
	return profile, nil
}

// --- DIFAL ---

func (s CalcService) ListDifal(ctx context.Context) ([]domain.DifalRate, error) {
	return s.repo.ListDifalRates(ctx, s.tenantID)
}

// OverrideResult reports whether a DIFAL override was persisted and the
// resulting rate as it now reads (with any active override applied).
type OverrideResult struct {
	Persisted bool
	Rate      domain.DifalRate
}

// PutDifalOverride applies an operator override to one UF's interna_pct.
// internaPct outside 0..35 ⇒ ErrInvalidRate (422). Absent UF ⇒
// ports.ErrDifalUFNotFound (404). |Δ| ≤ 0,049pp vs the current interna ⇒ no
// persist, Persisted=false, 200. |Δ| > 0,049pp ⇒ persisted + audited.
func (s CalcService) PutDifalOverride(ctx context.Context, uf, internaPct, actor string, at time.Time) (OverrideResult, error) {
	if !rateInRange(internaPct) {
		return OverrideResult{}, ErrInvalidRate
	}
	current, err := s.repo.RateForUF(ctx, s.tenantID, uf)
	if err != nil {
		return OverrideResult{}, err
	}
	currentInterna := current.InternaPct
	if current.Override != nil {
		currentInterna = current.Override.InternaPct
	}
	within, err := deltaWithinThreshold(currentInterna, internaPct)
	if err != nil {
		return OverrideResult{}, err
	}
	if within {
		return OverrideResult{Persisted: false, Rate: current}, nil
	}
	override := &domain.DifalOverride{
		InternaPct: internaPct,
		UpdatedAt:  at.UTC().Format(time.RFC3339),
		Actor:      actor,
	}
	if err := s.repo.UpsertDifalOverride(ctx, s.tenantID, uf, override); err != nil {
		return OverrideResult{}, err
	}
	updated, err := s.repo.RateForUF(ctx, s.tenantID, uf)
	if err != nil {
		return OverrideResult{}, err
	}
	return OverrideResult{Persisted: true, Rate: updated}, nil
}

// --- Scenarios ---

func (s CalcService) ListScenarios(ctx context.Context) ([]domain.Scenario, error) {
	return s.repo.ListScenarios(ctx, s.tenantID)
}

func (s CalcService) SaveScenario(ctx context.Context, scenario domain.Scenario) error {
	return s.repo.SaveScenario(ctx, s.tenantID, scenario)
}

func (s CalcService) DeleteScenario(ctx context.Context, id string) error {
	return s.repo.DeleteScenario(ctx, s.tenantID, id)
}

// --- Decompose / Solve ---

// DecomposeRequest carries the simulator inputs. Preco/ComissaoPct/Modalidade
// are always caller-supplied; TarifaFull/FreteProduto/Custo are optional
// overrides (nil ⇒ resolve from profile/cost port). When ProductID is set the
// service checks existence (ITEM_NOT_FOUND) and resolves custo_erp.
type DecomposeRequest struct {
	Preco        string
	ComissaoPct  string
	Modalidade   domain.Modalidade
	TarifaFull   *string
	FreteProduto *string
	Custo        *string
	ProductID    *int
	AsOf         time.Time
}

// DecomposeResult is a decomposition plus an optional blocking_state. The
// blocking_state is SEM_CUSTO when a known product's custo_erp is absent; nil
// otherwise. UNREACHABLE has no equivalent here (decompose always resolves).
type DecomposeResult struct {
	Decomposition domain.Decomposition
	BlockingState *string
}

// Decompose runs the single IC-04 formula over resolved inputs. Preco ≤ 0 ⇒
// ErrInvalidPrice (422). Unknown ProductID ⇒ ErrItemNotFound (404). A known
// product with no custo ⇒ 200 with blocking_state SEM_CUSTO and custo as an
// unknown component (never 0).
func (s CalcService) Decompose(ctx context.Context, req DecomposeRequest) (DecomposeResult, error) {
	if !priceIsPositive(req.Preco) {
		return DecomposeResult{}, ErrInvalidPrice
	}
	profile, err := s.repo.GetProfile(ctx, s.tenantID)
	if err != nil {
		return DecomposeResult{}, err
	}
	custo, blocking, err := s.resolveCusto(ctx, req)
	if err != nil {
		return DecomposeResult{}, err
	}
	in := domain.DecomposeInput{
		Preco:        req.Preco,
		ComissaoPct:  req.ComissaoPct,
		AliquotaPct:  profile.AliquotaPct,
		Modalidade:   req.Modalidade,
		TarifaFull:   resolveTarifaFull(req.TarifaFull, profile),
		FreteProduto: optionalMoney(req.FreteProduto),
		Custo:        custo,
	}
	s.applyDifal(ctx, profile, &in.DifalEnabled, &in.DestinoUF, &in.EfetivoPct)
	return DecomposeResult{Decomposition: domain.Decompose(in), BlockingState: blocking}, nil
}

// SolveRequest carries the bidirectional solver inputs (margem-alvo → preço).
type SolveRequest struct {
	TargetMargemPct string
	ComissaoPct     string
	Modalidade      domain.Modalidade
	TarifaFull      *string
	FreteProduto    *string
	Custo           *string
	ProductID       *int
	AsOf            time.Time
}

// SolveOutput is the solver result plus an optional blocking_state (SEM_CUSTO
// when a known product's custo is absent — the solver cannot converge without
// custo and reports it as a blocked, not errored, state).
type SolveOutput struct {
	Result        domain.SolveResult
	BlockingState *string
}

// SolveTarget resolves the price achieving a target margin. An unreachable
// target is NOT an error: Result.Reached=false + CeilingPct is returned as a
// 200 UNREACHABLE_TARGET by transport. Unknown ProductID ⇒ ErrItemNotFound.
func (s CalcService) SolveTarget(ctx context.Context, req SolveRequest) (SolveOutput, error) {
	profile, err := s.repo.GetProfile(ctx, s.tenantID)
	if err != nil {
		return SolveOutput{}, err
	}
	custo, blocking, err := s.resolveCusto(ctx, DecomposeRequest{ProductID: req.ProductID, Custo: req.Custo, AsOf: req.AsOf})
	if err != nil {
		return SolveOutput{}, err
	}
	in := domain.SolveInput{
		TargetMargemPct: req.TargetMargemPct,
		ComissaoPct:     req.ComissaoPct,
		AliquotaPct:     profile.AliquotaPct,
		Modalidade:      req.Modalidade,
		TarifaFull:      resolveTarifaFull(req.TarifaFull, profile),
		FreteProduto:    optionalMoney(req.FreteProduto),
		Custo:           custo,
	}
	s.applyDifal(ctx, profile, &in.DifalEnabled, &in.DestinoUF, &in.EfetivoPct)
	return SolveOutput{Result: domain.SolveTargetPrice(in), BlockingState: blocking}, nil
}

// resolveCusto returns the custo Money for the request: an explicit req.Custo
// wins; otherwise a ProductID is checked for existence (ItemNotFound) and its
// custo_erp resolved via the cost port (nil ⇒ SEM_CUSTO blocking_state). No
// product and no explicit custo ⇒ nil custo, no blocking (custo unknown).
func (s CalcService) resolveCusto(ctx context.Context, req DecomposeRequest) (*domain.Money, *string, error) {
	if req.Custo != nil {
		return &domain.Money{Amount: *req.Custo, Currency: "BRL"}, nil, nil
	}
	if req.ProductID == nil {
		return nil, nil, nil
	}
	if s.products != nil {
		exists, err := s.products.Exists(ctx, *req.ProductID)
		if err != nil {
			return nil, nil, err
		}
		if !exists {
			return nil, nil, ErrItemNotFound
		}
	}
	custo, err := s.cost.CostForProduct(ctx, *req.ProductID, req.AsOf)
	if err != nil {
		return nil, nil, err
	}
	if custo == nil {
		blocking := BlockingStateSemCusto
		return nil, &blocking, nil
	}
	return custo, nil, nil
}

// applyDifal resolves the profile's DIFAL destino to an efetivo_pct via the
// repository. Enabled with a resolvable destino ⇒ difal applied; a destino
// that is unset or absent from the table ⇒ difal treated as unknown
// (ADR-17), not an error — the 404 UF path lives on the DIFAL endpoints.
func (s CalcService) applyDifal(ctx context.Context, profile domain.CalcProfile, enabled *bool, destino, efetivo *string) {
	if !profile.DifalEnabled || profile.DifalDestinoUF == nil {
		*enabled = profile.DifalEnabled
		return
	}
	*enabled = true
	rate, err := s.repo.RateForUF(ctx, s.tenantID, *profile.DifalDestinoUF)
	if err != nil {
		return // destino unknown ⇒ leave DestinoUF/EfetivoPct empty ⇒ unknown component
	}
	res := domain.DifalForUF(rate)
	*destino = rate.UF
	*efetivo = res.EfetivoPct
}

// --- helpers ---

func resolveTarifaFull(override *string, profile domain.CalcProfile) *domain.Money {
	if override != nil {
		return &domain.Money{Amount: *override, Currency: "BRL"}
	}
	return profile.TarifaFull
}

func optionalMoney(amount *string) *domain.Money {
	if amount == nil {
		return nil
	}
	return &domain.Money{Amount: *amount, Currency: "BRL"}
}

func rateInRange(pct string) bool {
	r, err := domain.ParseRat(pct)
	if err != nil {
		return false
	}
	return r.Cmp(big.NewRat(0, 1)) >= 0 && r.Cmp(big.NewRat(35, 1)) <= 0
}

func priceIsPositive(preco string) bool {
	r, err := domain.ParseRat(preco)
	if err != nil {
		return false
	}
	return r.Sign() > 0
}

// deltaWithinThreshold reports whether |a - b| ≤ 0,049 percentage points.
func deltaWithinThreshold(a, b string) (bool, error) {
	ar, err := domain.ParseRat(a)
	if err != nil {
		return false, err
	}
	br, err := domain.ParseRat(b)
	if err != nil {
		return false, err
	}
	diff := new(big.Rat).Sub(ar, br)
	diff.Abs(diff)
	return diff.Cmp(overrideThresholdPP) <= 0, nil
}
