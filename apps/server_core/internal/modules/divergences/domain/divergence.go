// Package domain holds the divergences ledger's provider-agnostic types and
// pure evaluation rules (IC-02, ADR-10). It must never import a Mercado
// Livre / provider-specific type (Q6 architectural boundary, proven by
// boundary_test.go in the parent package).
package domain

import (
	"errors"
	"time"
)

// EntityType values (divergences.entity_type CHECK — IC-02 Fields).
const (
	EntityTypeListing   = "listing"
	EntityTypeOrderLine = "order_line"
)

// Kind values (divergences.kind CHECK — IC-02 Fields), extensible
// additively per IC-02 Compatibility Rules.
const (
	KindEstoque = "estoque"
	KindTarifa  = "tarifa"
)

// ErrExpectedObservedAtRequired is returned by Evaluate when
// expected_observed_at is missing. Named so callers assert the exact
// sentinel — never a silent fallback to now() (IC-02 Error Cases).
var ErrExpectedObservedAtRequired = errors.New("divergences: expected_observed_at is required")

// ErrObservedObservedAtRequired is returned by Evaluate when
// observed_observed_at is missing. Same rationale as
// ErrExpectedObservedAtRequired — both sides are always dated (IC-02 R-5).
var ErrObservedObservedAtRequired = errors.New("divergences: observed_observed_at is required")

// EvaluateInput is the input to DivergenceRecorder.Evaluate (IC-02
// Operations). Linked must be resolved by the caller (M-05/M-06 — a
// producer-side linkage check, not this port's concern): when false,
// Evaluate is a no-op skip, never an evaluation against a fabricated
// "expected" (IC-02 Persistence Expectations — "sem vínculo aprovado → não
// avaliar estoque").
type EvaluateInput struct {
	TenantID           string
	Provider           string
	EntityType         string
	EntityID           string
	Kind               string
	Linked             bool
	ExpectedValue      string
	ObservedValue      string
	ExpectedSource     string
	ObservedSource     string
	ExpectedObservedAt *time.Time
	ObservedObservedAt *time.Time
}

// EvaluateOutcome reports what Evaluate did, in mutually exclusive terms:
// exactly one of Skipped, Opened, Updated, Resolved, or a plain
// no-op (all false — convergent input with no open row to resolve) holds.
type EvaluateOutcome struct {
	// Skipped is true when Linked was false: Evaluate did not run at all,
	// no row was read or written (IC-02 Persistence Expectations).
	Skipped bool
	// Diverged reports whether this evaluation's expected/observed values
	// were outside tolerance for Kind.
	Diverged bool
	// Opened is true when a brand-new open row was created (no open row
	// existed yet for this entity+kind, and the values diverged).
	Opened bool
	// Updated is true when an already-open row's values/sources/observed_at
	// were refreshed by a divergent re-evaluation — detected_at unchanged
	// (IC-02 Persistence Expectations).
	Updated bool
	// Resolved is true when this evaluation set resolved_at on a
	// previously-open row (auto-resolve — the values converged within
	// tolerance).
	Resolved bool
}

// Validate checks the two mandatory observation timestamps (IC-02 Error
// Cases — no fallback to now(), ever). It does not check Linked; that is a
// skip, not a validation failure.
func (in EvaluateInput) Validate() error {
	if in.ExpectedObservedAt == nil {
		return ErrExpectedObservedAtRequired
	}
	if in.ObservedObservedAt == nil {
		return ErrObservedObservedAtRequired
	}
	return nil
}

// OpenDivergence is one open (resolved_at IS NULL) divergences row, as
// returned by DivergenceReader.ListOpenByEntity (IC-02 Required Outputs).
type OpenDivergence struct {
	Kind               string
	ExpectedValue      string
	ObservedValue      string
	ExpectedSource     string
	ObservedSource     string
	ExpectedObservedAt time.Time
	ObservedObservedAt time.Time
	DetectedAt         time.Time
}
