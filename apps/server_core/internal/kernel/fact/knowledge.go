// Package fact makes "unknown is never zero" structural instead of doctrinal.
//
// The repository cites that rule 1378 times and still broke it in its own
// shared core, because a rule enforced by remembering is a rule that gets
// forgotten. Here the invalid combinations are unbuildable: Unknown never
// carries a value, Known always does, and neither exists without evidence.
//
// The cost of getting this wrong is not cosmetic. Solving a target-margin
// price with one unknown component U gives P = 178.69 + 2.02*U: every unit of
// ignored cost moves the answer by two. A system that treats unknown as zero
// does not err slightly, it errs by twice what it ignored.
package fact

import (
	"errors"
	"fmt"
	"strings"

	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

// Knowledge is what we know about a quantity. The zero value is Unknown, on
// purpose: a Fact that nobody built must be the safe one, not the confident one.
type Knowledge uint8

const (
	// Unknown means we have no value. It is the zero value.
	Unknown Knowledge = iota
	// Known means we observed the value, including an observed zero.
	Known
	// Estimated means we derived the value and it is usable with a caveat.
	Estimated
	// NotApplicable means the quantity does not exist for this case — which is
	// a different statement from not knowing it.
	NotApplicable
)

// String renders the state for logs and error messages.
func (k Knowledge) String() string {
	switch k {
	case Known:
		return "known"
	case Estimated:
		return "estimated"
	case NotApplicable:
		return "not_applicable"
	default:
		return "unknown"
	}
}

var (
	// ErrReasonRequired is returned when a non-Known state carries no reason.
	ErrReasonRequired = errors.New("fact: state requires a reason code")
	// ErrEvidenceRequired is returned when a fact is built without evidence.
	ErrEvidenceRequired = errors.New("fact: evidence is required")
	// ErrNotUsable is the panic value of MustValue on a fact with no value.
	ErrNotUsable = errors.New("fact: no value in this state")
)

// Fact carries a value together with what we know about it and how we know it.
// Every field is unexported, so the only way to obtain one is a constructor,
// and every constructor validates.
type Fact[T any] struct {
	state    Knowledge
	value    *T
	reason   string
	evidence provenance.Evidence
}

// NewKnown records an observed value. A zero value here is a measured zero and
// is a fact.
func NewKnown[T any](v T, e provenance.Evidence) (Fact[T], error) {
	if e.IsZero() {
		return Fact[T]{}, fmt.Errorf("%w: known value %v", ErrEvidenceRequired, v)
	}
	return Fact[T]{state: Known, value: &v, evidence: e}, nil
}

// NewEstimated records a derived value with the reason it is an estimate.
func NewEstimated[T any](v T, reason string, e provenance.Evidence) (Fact[T], error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return Fact[T]{}, fmt.Errorf("%w: estimated", ErrReasonRequired)
	}
	if e.IsZero() {
		return Fact[T]{}, fmt.Errorf("%w: estimated %v", ErrEvidenceRequired, v)
	}
	return Fact[T]{state: Estimated, value: &v, reason: reason, evidence: e}, nil
}

// NewUnknown records an absence. It takes no value — there is no parameter to
// pass a zero into — and demands the reason we do not have one.
func NewUnknown[T any](reason string, e provenance.Evidence) (Fact[T], error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return Fact[T]{}, fmt.Errorf("%w: unknown", ErrReasonRequired)
	}
	if e.IsZero() {
		return Fact[T]{}, fmt.Errorf("%w: unknown", ErrEvidenceRequired)
	}
	return Fact[T]{state: Unknown, reason: reason, evidence: e}, nil
}

// NewNotApplicable records that the quantity does not exist for this case.
func NewNotApplicable[T any](reason string, e provenance.Evidence) (Fact[T], error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return Fact[T]{}, fmt.Errorf("%w: not applicable", ErrReasonRequired)
	}
	if e.IsZero() {
		return Fact[T]{}, fmt.Errorf("%w: not applicable", ErrEvidenceRequired)
	}
	return Fact[T]{state: NotApplicable, reason: reason, evidence: e}, nil
}

// Value returns the value and whether there is one. The bool is not optional
// to read: `v, _ := f.Value()` on an Unknown fact hands back the zero value of
// T, which is exactly the mistake this package exists to prevent. The vet-level
// check for that pattern lives in internal/arch.
func (f Fact[T]) Value() (T, bool) {
	if f.value == nil {
		var zero T
		return zero, false
	}
	return *f.value, true
}

// MustValue returns the value or panics. It is for call sites that have already
// checked IsUsable in the same function.
func (f Fact[T]) MustValue() T {
	v, ok := f.Value()
	if !ok {
		panic(fmt.Errorf("%w: state=%s reason=%q evidence=%s",
			ErrNotUsable, f.state, f.reason, f.evidence.Ref()))
	}
	return v
}

// State returns what we know.
func (f Fact[T]) State() Knowledge { return f.state }

// Reason returns why, for every state except Known.
func (f Fact[T]) Reason() string { return f.reason }

// Evidence returns how we know.
func (f Fact[T]) Evidence() provenance.Evidence { return f.evidence }

// IsUsable reports whether a calculation may consume this fact. Known and
// Estimated are usable; Unknown and NotApplicable are not.
func (f Fact[T]) IsUsable() bool { return f.value != nil }
