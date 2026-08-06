package fact_test

import (
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

func ev(t *testing.T) provenance.Evidence {
	t.Helper()
	e, err := provenance.NewEvidence("sankhya", "product", "10529",
		time.Date(2026, 8, 6, 12, 15, 0, 0, time.UTC), "sha256:ab91")
	if err != nil {
		t.Fatalf("NewEvidence: %v", err)
	}
	return e
}

func TestKnownCarriesItsValue(t *testing.T) {
	f, err := fact.NewKnown(exact.FromInt(18), ev(t))
	if err != nil {
		t.Fatalf("NewKnown: %v", err)
	}
	v, ok := f.Value()
	if !ok {
		t.Fatal("Value() ok = false for a Known fact")
	}
	if v.Cmp(exact.FromInt(18)) != 0 {
		t.Fatalf("Value() = %s, want 18", v.StringFixed(0))
	}
	if f.State() != fact.Known {
		t.Fatalf("State() = %v, want Known", f.State())
	}
}

// Known zero is a value. It is NOT the same thing as Unknown, and this test is
// the one that stops the two from ever collapsing into each other again.
func TestKnownZeroIsNotUnknown(t *testing.T) {
	known, err := fact.NewKnown(exact.FromInt(0), ev(t))
	if err != nil {
		t.Fatalf("NewKnown: %v", err)
	}
	v, ok := known.Value()
	if !ok {
		t.Fatal("Value() ok = false for Known zero")
	}
	if !v.IsZero() {
		t.Fatal("Known zero did not carry zero")
	}
	if !known.IsUsable() {
		t.Fatal("Known zero is not usable; a measured zero is a fact")
	}

	unknown, err := fact.NewUnknown[exact.Decimal]("marketplace did not return it", ev(t))
	if err != nil {
		t.Fatalf("NewUnknown: %v", err)
	}
	if _, ok := unknown.Value(); ok {
		t.Fatal("Unknown returned a value")
	}
	if unknown.IsUsable() {
		t.Fatal("Unknown is usable; nothing may compute with it")
	}
	if unknown.State() == known.State() {
		t.Fatal("Unknown and Known collapsed into the same state")
	}
}

func TestUnknownRequiresAReason(t *testing.T) {
	if _, err := fact.NewUnknown[exact.Decimal]("", ev(t)); err == nil {
		t.Fatal("NewUnknown with empty reason returned no error")
	}
	if _, err := fact.NewUnknown[exact.Decimal]("   ", ev(t)); err == nil {
		t.Fatal("NewUnknown with blank reason returned no error")
	}
}

func TestEstimatedRequiresAReason(t *testing.T) {
	if _, err := fact.NewEstimated(exact.FromInt(5), "", ev(t)); err == nil {
		t.Fatal("NewEstimated with empty reason returned no error")
	}
}

func TestEveryStateRequiresEvidence(t *testing.T) {
	var none provenance.Evidence
	if _, err := fact.NewKnown(exact.FromInt(1), none); err == nil {
		t.Fatal("NewKnown without evidence returned no error")
	}
	if _, err := fact.NewUnknown[exact.Decimal]("no data", none); err == nil {
		t.Fatal("NewUnknown without evidence returned no error")
	}
	if _, err := fact.NewEstimated(exact.FromInt(1), "modelled", none); err == nil {
		t.Fatal("NewEstimated without evidence returned no error")
	}
	if _, err := fact.NewNotApplicable[exact.Decimal]("CST 60 has no own ICMS", none); err == nil {
		t.Fatal("NewNotApplicable without evidence returned no error")
	}
}

func TestNotApplicableIsNotUnknown(t *testing.T) {
	na, err := fact.NewNotApplicable[exact.Decimal]("CST 60: no own ICMS on this operation", ev(t))
	if err != nil {
		t.Fatalf("NewNotApplicable: %v", err)
	}
	if na.State() != fact.NotApplicable {
		t.Fatalf("State() = %v, want NotApplicable", na.State())
	}
	if _, ok := na.Value(); ok {
		t.Fatal("NotApplicable returned a value")
	}
	if na.IsUsable() {
		t.Fatal("NotApplicable is usable")
	}
}

func TestZeroFactIsUnknownAndUnusable(t *testing.T) {
	var f fact.Fact[exact.Decimal]
	if f.State() != fact.Unknown {
		t.Fatalf("zero Fact.State() = %v, want Unknown; the zero value must be the safe one", f.State())
	}
	if _, ok := f.Value(); ok {
		t.Fatal("zero Fact returned a value")
	}
	if f.IsUsable() {
		t.Fatal("zero Fact is usable")
	}
}

func TestMustValuePanicsOnUnknown(t *testing.T) {
	f, err := fact.NewUnknown[exact.Decimal]("no data", ev(t))
	if err != nil {
		t.Fatalf("NewUnknown: %v", err)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("MustValue on Unknown did not panic")
		}
	}()
	_ = f.MustValue()
}

func TestEstimatedIsUsableButDistinct(t *testing.T) {
	f, err := fact.NewEstimated(exact.FromInt(7), "modelled from last 30 days", ev(t))
	if err != nil {
		t.Fatalf("NewEstimated: %v", err)
	}
	if !f.IsUsable() {
		t.Fatal("Estimated is not usable; an estimate is a value with a caveat, not an absence")
	}
	if f.State() == fact.Known {
		t.Fatal("Estimated collapsed into Known")
	}
	if f.Reason() == "" {
		t.Fatal("Estimated lost its reason")
	}
}
