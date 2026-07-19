package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

var _ ports.CalcRepository = (*CalcRepository)(nil)
var _ ports.TariffDefaultsStore = (*CalcRepository)(nil)

// CalcRepository persists the IC-04 calculator state: the tenant CalcProfile,
// the DIFAL rate table (with overrides), and saved scenarios. Numeric columns
// are read back as ::text so decimal money/percent values never round-trip
// through float64. tenantID is passed per call (not stored) so one pool-backed
// repo serves every tenant.
type CalcRepository struct {
	pool *pgxpool.Pool
}

func NewCalcRepository(pool *pgxpool.Pool) *CalcRepository {
	return &CalcRepository{pool: pool}
}

// GetProfile returns the tenant's stored CalcProfile, or the built-in default
// (SIMPLES 4%, origem "default") when no row exists yet — the initial state is
// materialized on first read, not seeded by migration.
func (r *CalcRepository) GetProfile(ctx context.Context, tenantID string) (domain.CalcProfile, error) {
	var (
		regime     string
		aliquota   string
		limVerde   string
		limAmarelo string
		tarifaFull *string
		difalOn    bool
		difalUF    *string
		origem     string
	)
	err := r.pool.QueryRow(ctx, `
		SELECT regime, aliquota_pct::text, limiar_verde_pct::text, limiar_amarelo_pct::text,
		       tarifa_full::text, difal_enabled, difal_destino_uf, origem
		FROM pricing_calc_profiles
		WHERE tenant_id = $1
	`, tenantID).Scan(&regime, &aliquota, &limVerde, &limAmarelo, &tarifaFull, &difalOn, &difalUF, &origem)
	if err == pgx.ErrNoRows {
		return domain.NewDefaultCalcProfile(), nil
	}
	if err != nil {
		return domain.CalcProfile{}, err
	}
	profile := domain.CalcProfile{
		Regime:           domain.Regime(regime),
		AliquotaPct:      aliquota,
		LimiarVerdePct:   limVerde,
		LimiarAmareloPct: limAmarelo,
		DifalEnabled:     difalOn,
		DifalDestinoUF:   difalUF,
		Origem:           origem,
	}
	if tarifaFull != nil {
		profile.TarifaFull = &domain.Money{Amount: *tarifaFull, Currency: "BRL"}
	}
	return profile, nil
}

// UpsertProfile writes the tenant's CalcProfile. A nil TarifaFull persists as
// SQL NULL (unset/unknown), never 0 (ADR-17).
func (r *CalcRepository) UpsertProfile(ctx context.Context, tenantID string, profile domain.CalcProfile) error {
	var tarifaFull *string
	if profile.TarifaFull != nil {
		tarifaFull = &profile.TarifaFull.Amount
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO pricing_calc_profiles (
			tenant_id, regime, aliquota_pct, limiar_verde_pct, limiar_amarelo_pct,
			tarifa_full, difal_enabled, difal_destino_uf, origem, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now())
		ON CONFLICT (tenant_id) DO UPDATE SET
			regime = EXCLUDED.regime,
			aliquota_pct = EXCLUDED.aliquota_pct,
			limiar_verde_pct = EXCLUDED.limiar_verde_pct,
			limiar_amarelo_pct = EXCLUDED.limiar_amarelo_pct,
			tarifa_full = EXCLUDED.tarifa_full,
			difal_enabled = EXCLUDED.difal_enabled,
			difal_destino_uf = EXCLUDED.difal_destino_uf,
			origem = EXCLUDED.origem,
			updated_at = now()
	`, tenantID, string(profile.Regime), profile.AliquotaPct, profile.LimiarVerdePct, profile.LimiarAmareloPct,
		tarifaFull, profile.DifalEnabled, profile.DifalDestinoUF, profile.Origem)
	return err
}

// RateForUF satisfies ports.DifalReader (and the frozen ports.DifalForUF read):
// it returns the tenant's DIFAL row for uf, reflecting any active override.
// ErrDifalUFNotFound when the tenant has no row for that UF.
func (r *CalcRepository) RateForUF(ctx context.Context, tenantID, uf string) (domain.DifalRate, error) {
	rate, err := scanDifalRate(r.pool.QueryRow(ctx, `
		SELECT uf, interna_pct::text, interestadual_pct::text, efetivo_pct::text, origem_versao,
		       override_interna_pct::text, override_updated_at, override_actor
		FROM pricing_difal_rates
		WHERE tenant_id = $1 AND uf = $2
	`, tenantID, uf))
	if err == pgx.ErrNoRows {
		return domain.DifalRate{}, ports.ErrDifalUFNotFound
	}
	if err != nil {
		return domain.DifalRate{}, err
	}
	return rate, nil
}

// ListDifalRates returns every DIFAL row for the tenant ordered by UF, each
// with its active override (if any) reflected.
func (r *CalcRepository) ListDifalRates(ctx context.Context, tenantID string) ([]domain.DifalRate, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT uf, interna_pct::text, interestadual_pct::text, efetivo_pct::text, origem_versao,
		       override_interna_pct::text, override_updated_at, override_actor
		FROM pricing_difal_rates
		WHERE tenant_id = $1
		ORDER BY uf
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rates := make([]domain.DifalRate, 0, 27)
	for rows.Next() {
		rate, err := scanDifalRate(rows)
		if err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}
	return rates, rows.Err()
}

// UpsertDifalOverride sets (override != nil) or clears (override == nil) the
// override_* columns for one UF. The |Δ| > 0,049pp audit gate is applied by the
// caller before this write. ErrDifalUFNotFound when the UF row is absent.
func (r *CalcRepository) UpsertDifalOverride(ctx context.Context, tenantID, uf string, override *domain.DifalOverride) error {
	if override == nil {
		tag, err := r.pool.Exec(ctx, `
			UPDATE pricing_difal_rates
			SET override_interna_pct = NULL, override_updated_at = NULL, override_actor = NULL
			WHERE tenant_id = $1 AND uf = $2
		`, tenantID, uf)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return ports.ErrDifalUFNotFound
		}
		return nil
	}
	updatedAt, err := time.Parse(time.RFC3339, override.UpdatedAt)
	if err != nil {
		return err
	}
	tag, err := r.pool.Exec(ctx, `
		UPDATE pricing_difal_rates
		SET override_interna_pct = $3, override_updated_at = $4, override_actor = $5
		WHERE tenant_id = $1 AND uf = $2
	`, tenantID, uf, override.InternaPct, updatedAt, override.Actor)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ports.ErrDifalUFNotFound
	}
	return nil
}

// ListScenarios returns the tenant's saved scenarios newest-first.
func (r *CalcRepository) ListScenarios(ctx context.Context, tenantID string) ([]domain.Scenario, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, payload::text, created_at
		FROM pricing_scenarios
		WHERE tenant_id = $1
		ORDER BY created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scenarios := make([]domain.Scenario, 0)
	for rows.Next() {
		var s domain.Scenario
		if err := rows.Scan(&s.ID, &s.Name, &s.Payload, &s.CreatedAt); err != nil {
			return nil, err
		}
		scenarios = append(scenarios, s)
	}
	return scenarios, rows.Err()
}

// SaveScenario inserts or updates a scenario for the tenant. A zero CreatedAt
// falls back to now() so callers may leave it unset on create.
func (r *CalcRepository) SaveScenario(ctx context.Context, tenantID string, scenario domain.Scenario) error {
	var createdAt *time.Time
	if !scenario.CreatedAt.IsZero() {
		createdAt = &scenario.CreatedAt
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO pricing_scenarios (id, tenant_id, name, payload, created_at)
		VALUES ($1, $2, $3, $4::jsonb, COALESCE($5, now()))
		ON CONFLICT (tenant_id, id) DO UPDATE SET
			name = EXCLUDED.name,
			payload = EXCLUDED.payload
	`, scenario.ID, tenantID, scenario.Name, scenario.Payload, createdAt)
	return err
}

// DeleteScenario removes a scenario. ErrScenarioNotFound when the id is absent
// for the tenant.
func (r *CalcRepository) DeleteScenario(ctx context.Context, tenantID, id string) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM pricing_scenarios WHERE tenant_id = $1 AND id = $2
	`, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ports.ErrScenarioNotFound
	}
	return nil
}

// GetTariffDefaults returns the tenant/installation's stored TariffDefaults,
// materializing the row on first read: INSERT ... ON CONFLICT DO NOTHING lets
// the DB column DEFAULTs (13.00/16.00/sem_dados) fill a new row, then it is
// read back. frete_estimativa_amount scans into *string — NULL stays nil
// (ADR-17), never 0.
func (r *CalcRepository) GetTariffDefaults(ctx context.Context, tenantID, installationID string) (domain.TariffDefaults, error) {
	if _, err := r.pool.Exec(ctx, `
		INSERT INTO pricing_tariff_defaults (tenant_id, installation_id)
		VALUES ($1, $2)
		ON CONFLICT (tenant_id, installation_id) DO NOTHING
	`, tenantID, installationID); err != nil {
		return domain.TariffDefaults{}, err
	}

	var (
		comissaoClassico string
		comissaoPremium  string
		freteEstimativa  *string
		fretePolicy      string
	)
	err := r.pool.QueryRow(ctx, `
		SELECT comissao_classico_pct::text, comissao_premium_pct::text,
		       frete_estimativa_amount::text, frete_policy
		FROM pricing_tariff_defaults
		WHERE tenant_id = $1 AND installation_id = $2
	`, tenantID, installationID).Scan(&comissaoClassico, &comissaoPremium, &freteEstimativa, &fretePolicy)
	if err != nil {
		return domain.TariffDefaults{}, err
	}
	return domain.TariffDefaults{
		ComissaoClassicoPct:   comissaoClassico,
		ComissaoPremiumPct:    comissaoPremium,
		FreteEstimativaAmount: freteEstimativa,
		FretePolicy:           fretePolicy,
	}, nil
}

// UpsertTariffDefaults validates in.FretePolicy (domain.ErrInvalidFretePolicy
// otherwise, no write attempted) and writes the full row, returning it as
// stored. tenant_id + installation_id predicate is mandatory on conflict.
func (r *CalcRepository) UpsertTariffDefaults(ctx context.Context, tenantID, installationID string, in domain.TariffDefaults) (domain.TariffDefaults, error) {
	if !domain.ValidFretePolicy(in.FretePolicy) {
		return domain.TariffDefaults{}, domain.ErrInvalidFretePolicy
	}

	var (
		comissaoClassico string
		comissaoPremium  string
		freteEstimativa  *string
		fretePolicy      string
	)
	err := r.pool.QueryRow(ctx, `
		INSERT INTO pricing_tariff_defaults (
			tenant_id, installation_id, comissao_classico_pct, comissao_premium_pct,
			frete_estimativa_amount, frete_policy, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, now())
		ON CONFLICT (tenant_id, installation_id) DO UPDATE SET
			comissao_classico_pct = EXCLUDED.comissao_classico_pct,
			comissao_premium_pct = EXCLUDED.comissao_premium_pct,
			frete_estimativa_amount = EXCLUDED.frete_estimativa_amount,
			frete_policy = EXCLUDED.frete_policy,
			updated_at = now()
		RETURNING comissao_classico_pct::text, comissao_premium_pct::text,
		          frete_estimativa_amount::text, frete_policy
	`, tenantID, installationID, in.ComissaoClassicoPct, in.ComissaoPremiumPct,
		in.FreteEstimativaAmount, in.FretePolicy,
	).Scan(&comissaoClassico, &comissaoPremium, &freteEstimativa, &fretePolicy)
	if err != nil {
		return domain.TariffDefaults{}, err
	}
	return domain.TariffDefaults{
		ComissaoClassicoPct:   comissaoClassico,
		ComissaoPremiumPct:    comissaoPremium,
		FreteEstimativaAmount: freteEstimativa,
		FretePolicy:           fretePolicy,
	}, nil
}

// rowScanner is the read surface shared by pgx.Row (QueryRow) and pgx.Rows
// (Query loop), letting one scanDifalRate serve both RateForUF and
// ListDifalRates.
type rowScanner interface {
	Scan(dest ...any) error
}

// scanDifalRate reads one DIFAL row and attaches an active override when the
// override_* columns are populated (the paired CHECK guarantees all-or-none).
func scanDifalRate(row rowScanner) (domain.DifalRate, error) {
	var (
		rate        domain.DifalRate
		ovInterna   *string
		ovUpdatedAt *time.Time
		ovActor     *string
	)
	if err := row.Scan(
		&rate.UF, &rate.InternaPct, &rate.InterestadualPct, &rate.EfetivoPct, &rate.OrigemVersao,
		&ovInterna, &ovUpdatedAt, &ovActor,
	); err != nil {
		return domain.DifalRate{}, err
	}
	if ovInterna != nil && ovUpdatedAt != nil && ovActor != nil {
		rate.Override = &domain.DifalOverride{
			InternaPct: *ovInterna,
			UpdatedAt:  ovUpdatedAt.UTC().Format(time.RFC3339),
			Actor:      *ovActor,
		}
	}
	return rate, nil
}
