package ports

import (
	"context"
	"errors"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
)

// ErrDifalUFNotFound is returned by RateForUF when the tenant has no DIFAL row
// for the requested UF (transport maps it to 404 UF_NOT_FOUND).
var ErrDifalUFNotFound = errors.New("PRICING_UF_NOT_FOUND")

// ErrScenarioNotFound is returned when a scenario id is absent for the tenant
// (transport maps it to 404 SCENARIO_NOT_FOUND).
var ErrScenarioNotFound = errors.New("PRICING_SCENARIO_NOT_FOUND")

// CalcRepository is the persistence port for the pricing calculator: the
// tenant CalcProfile, the DIFAL rate table (with overrides), and saved
// scenarios. Every method is tenant-scoped. It composes DifalReader so the
// same store backs the frozen ports.DifalForUF read.
type CalcRepository interface {
	DifalReader

	// GetProfile returns the tenant's CalcProfile, or NewDefaultCalcProfile
	// (SIMPLES 4%, origem "default") when none is stored yet.
	GetProfile(ctx context.Context, tenantID string) (domain.CalcProfile, error)
	UpsertProfile(ctx context.Context, tenantID string, profile domain.CalcProfile) error

	// ListDifalRates returns all 27 rates ordered by UF; an active override is
	// reflected on each row's Override.
	ListDifalRates(ctx context.Context, tenantID string) ([]domain.DifalRate, error)
	// UpsertDifalOverride sets (override != nil) or clears (override == nil) the
	// override columns for one UF. The Δ>0,049pp audit gate is applied by the
	// caller (transport) before persisting.
	UpsertDifalOverride(ctx context.Context, tenantID, uf string, override *domain.DifalOverride) error

	// ListScenarios returns the tenant's saved scenarios newest-first.
	ListScenarios(ctx context.Context, tenantID string) ([]domain.Scenario, error)
	SaveScenario(ctx context.Context, tenantID string, scenario domain.Scenario) error
	DeleteScenario(ctx context.Context, tenantID, id string) error
}
