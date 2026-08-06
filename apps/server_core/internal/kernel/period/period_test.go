package period_test

import (
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/kernel/period"
)

func d(day int) time.Time { return time.Date(2026, 8, day, 0, 0, 0, 0, time.UTC) }

func TestNewRejectsZeroFrom(t *testing.T) {
	if _, err := period.New(time.Time{}, d(10)); err == nil {
		t.Fatal("New with zero from returned no error")
	}
}

func TestNewRejectsUntilBeforeFrom(t *testing.T) {
	if _, err := period.New(d(10), d(5)); err == nil {
		t.Fatal("New with until before from returned no error")
	}
}

func TestNewRejectsEmptyInterval(t *testing.T) {
	if _, err := period.New(d(10), d(10)); err == nil {
		t.Fatal("New with until == from returned no error; a half-open [t,t) holds nothing")
	}
}

// Half-open [from, until). The instant `until` belongs to the NEXT period, so
// two adjacent periods never both claim the same instant.
func TestContainsIsHalfOpen(t *testing.T) {
	p, err := period.New(d(5), d(10))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if !p.Contains(d(5)) {
		t.Fatal("Contains(from) = false; from is inside")
	}
	if p.Contains(d(10)) {
		t.Fatal("Contains(until) = true; until is outside")
	}
	if p.Contains(d(4)) || p.Contains(d(11)) {
		t.Fatal("Contains returned true outside the interval")
	}
}

func TestOpenPeriodHasNoUntil(t *testing.T) {
	p, err := period.From(d(5))
	if err != nil {
		t.Fatalf("From: %v", err)
	}
	if !p.IsOpen() {
		t.Fatal("IsOpen() = false for an open period")
	}
	if _, ok := p.Until(); ok {
		t.Fatal("Until() reported a bound for an open period")
	}
	if !p.Contains(d(9999 % 28)) && !p.Contains(d(28)) {
		t.Fatal("open period does not contain a later instant")
	}
}

func TestOverlapsIsSymmetric(t *testing.T) {
	a, _ := period.New(d(1), d(10))
	b, _ := period.New(d(5), d(15))
	c, _ := period.New(d(10), d(20)) // adjacent to a, not overlapping
	if !a.Overlaps(b) || !b.Overlaps(a) {
		t.Fatal("overlapping periods reported as disjoint")
	}
	if a.Overlaps(c) || c.Overlaps(a) {
		t.Fatal("adjacent half-open periods reported as overlapping")
	}
}

func TestZeroPeriodIsZero(t *testing.T) {
	var p period.EffectivePeriod
	if !p.IsZero() {
		t.Fatal("zero EffectivePeriod.IsZero() = false")
	}
	if p.Contains(d(5)) {
		t.Fatal("zero period contains an instant")
	}
}
