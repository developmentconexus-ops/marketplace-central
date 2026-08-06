// Package period carries validity in time. Costs, tax rates, commissions and
// fiscal classifications change value without changing identity; a number
// without its period is a number whose meaning has to be reconstructed later,
// and by then nobody remembers.
//
// The interval is half-open, [from, until). Two adjacent periods therefore
// never both claim the same instant, which is the property that makes
// "the cost as of T" a question with exactly one answer.
package period

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrNoFrom is returned when the period has no start.
	ErrNoFrom = errors.New("period: from is the zero time")
	// ErrNotOrdered is returned when until is not strictly after from.
	ErrNotOrdered = errors.New("period: until must be strictly after from")
)

// EffectivePeriod is the half-open interval a fact is valid for.
type EffectivePeriod struct {
	from  time.Time
	until time.Time // zero means open-ended
}

// New builds a bounded period.
func New(from, until time.Time) (EffectivePeriod, error) {
	if from.IsZero() {
		return EffectivePeriod{}, ErrNoFrom
	}
	if until.IsZero() {
		return EffectivePeriod{}, fmt.Errorf("%w: use From for an open period", ErrNotOrdered)
	}
	if !until.After(from) {
		return EffectivePeriod{}, fmt.Errorf("%w: from=%s until=%s", ErrNotOrdered, from, until)
	}
	return EffectivePeriod{from: from.UTC(), until: until.UTC()}, nil
}

// From builds an open-ended period that starts at from and has no end yet.
func From(from time.Time) (EffectivePeriod, error) {
	if from.IsZero() {
		return EffectivePeriod{}, ErrNoFrom
	}
	return EffectivePeriod{from: from.UTC()}, nil
}

// Start returns the first instant the fact holds.
func (p EffectivePeriod) Start() time.Time { return p.from }

// Until returns the first instant the fact no longer holds, and whether the
// period is bounded at all.
func (p EffectivePeriod) Until() (time.Time, bool) {
	if p.until.IsZero() {
		return time.Time{}, false
	}
	return p.until, true
}

// IsOpen reports whether the period has no end yet.
func (p EffectivePeriod) IsOpen() bool { return !p.from.IsZero() && p.until.IsZero() }

// Contains reports whether t falls in [from, until).
func (p EffectivePeriod) Contains(t time.Time) bool {
	if p.IsZero() {
		return false
	}
	if t.Before(p.from) {
		return false
	}
	if p.until.IsZero() {
		return true
	}
	return t.Before(p.until)
}

// Overlaps reports whether the two periods share at least one instant.
func (p EffectivePeriod) Overlaps(o EffectivePeriod) bool {
	if p.IsZero() || o.IsZero() {
		return false
	}
	if !p.until.IsZero() && !o.from.Before(p.until) {
		return false
	}
	if !o.until.IsZero() && !p.from.Before(o.until) {
		return false
	}
	return true
}

// IsZero reports whether this is the zero value rather than a built period.
func (p EffectivePeriod) IsZero() bool { return p.from.IsZero() }
