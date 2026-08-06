// Package postgres is the Postgres-backed implementation of the divergences
// core ports (IC-02, ADR-011). The one-open-row invariant is enforced by the
// partial unique index divergences_one_open_row (0087_divergences.sql); this
// adapter's INSERT ... ON CONFLICT ... WHERE resolved_at IS NULL DO UPDATE
// targets exactly that index, following the same upsert idiom proven at
// internal_read/adapters/mirror/writer.go.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/modules/divergences/domain"
	"marketplace-central/apps/server_core/internal/modules/divergences/ports"
)

// Recorder is the Postgres-backed ports.DivergenceRecorder.
type Recorder struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

// RecorderOption configures a Recorder at construction.
type RecorderOption func(*Recorder)

// WithClock overrides the Recorder's clock (used to stamp detected_at /
// last_evaluated_at / resolved_at) — tests use this to prove detected_at
// stays immutable across a later re-evaluation without racing the wall
// clock.
func WithClock(now func() time.Time) RecorderOption {
	return func(r *Recorder) { r.now = now }
}

// NewRecorder builds a Recorder over pool. A nil pool is valid as long as
// every call is rejected or skipped before reaching the pool (see
// recorder_test.go) — the zero-DB validation/skip paths never dereference it.
func NewRecorder(pool *pgxpool.Pool, opts ...RecorderOption) *Recorder {
	r := &Recorder{pool: pool, now: time.Now}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

var _ ports.DivergenceRecorder = (*Recorder)(nil)

// divergentUpsertSQL opens a new row or refreshes an already-open one in a
// single statement. The `xmax = 0` trick distinguishes a genuine INSERT
// (xmax 0) from an UPDATE via the ON CONFLICT arm (xmax non-zero) so the
// caller can report Opened vs Updated without a second round trip.
// detected_at is only ever written by the INSERT arm — the DO UPDATE SET
// list deliberately omits it, so a later divergent re-evaluation of the same
// open row can never move it (IC-02 Persistence Expectations).
const divergentUpsertSQL = `
INSERT INTO divergences (
	tenant_id, provider, entity_type, entity_id, kind,
	expected_value, observed_value, expected_source, observed_source,
	expected_observed_at, observed_observed_at, detected_at, last_evaluated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$12)
ON CONFLICT (tenant_id, provider, entity_type, entity_id, kind) WHERE resolved_at IS NULL
DO UPDATE SET
	expected_value       = EXCLUDED.expected_value,
	observed_value       = EXCLUDED.observed_value,
	expected_source      = EXCLUDED.expected_source,
	observed_source      = EXCLUDED.observed_source,
	expected_observed_at = EXCLUDED.expected_observed_at,
	observed_observed_at = EXCLUDED.observed_observed_at,
	last_evaluated_at    = EXCLUDED.last_evaluated_at
RETURNING (xmax = 0) AS inserted`

// resolveOpenRowSQL auto-resolves an open row when the new observation
// converges. If no row is open, RowsAffected is 0 — a legitimate no-op, not
// an error (nothing to resolve, nothing to open).
const resolveOpenRowSQL = `
UPDATE divergences
SET resolved_at = $6, last_evaluated_at = $6
WHERE tenant_id = $1 AND provider = $2 AND entity_type = $3 AND entity_id = $4 AND kind = $5
  AND resolved_at IS NULL`

// Evaluate implements IC-02's one-open-row-per-(entity,kind) decision
// (ADR-011). Order of checks: (1) no linkage → skip, before any other
// validation — there is nothing to evaluate; (2) both observation
// timestamps are mandatory, named-error rejected, never defaulted; (3)
// compare expected vs observed within Kind's tolerance and either
// auto-resolve an open row (convergent) or open/refresh one (divergent).
func (r *Recorder) Evaluate(ctx context.Context, input domain.EvaluateInput) (domain.EvaluateOutcome, error) {
	if !input.Linked {
		return domain.EvaluateOutcome{Skipped: true}, nil
	}
	if err := input.Validate(); err != nil {
		return domain.EvaluateOutcome{}, err
	}

	diverged, err := domain.Diverges(input.Kind, input.ExpectedValue, input.ObservedValue)
	if err != nil {
		return domain.EvaluateOutcome{}, fmt.Errorf("divergences: evaluate: %w", err)
	}

	now := r.now()

	if !diverged {
		tag, err := r.pool.Exec(ctx, resolveOpenRowSQL,
			input.TenantID, input.Provider, input.EntityType, input.EntityID, input.Kind, now)
		if err != nil {
			return domain.EvaluateOutcome{}, fmt.Errorf("divergences: resolve open row: %w", err)
		}
		return domain.EvaluateOutcome{Diverged: false, Resolved: tag.RowsAffected() > 0}, nil
	}

	var inserted bool
	err = r.pool.QueryRow(ctx, divergentUpsertSQL,
		input.TenantID, input.Provider, input.EntityType, input.EntityID, input.Kind,
		input.ExpectedValue, input.ObservedValue, input.ExpectedSource, input.ObservedSource,
		input.ExpectedObservedAt, input.ObservedObservedAt, now,
	).Scan(&inserted)
	if err != nil {
		return domain.EvaluateOutcome{}, fmt.Errorf("divergences: upsert open row: %w", err)
	}
	return domain.EvaluateOutcome{Diverged: true, Opened: inserted, Updated: !inserted}, nil
}
