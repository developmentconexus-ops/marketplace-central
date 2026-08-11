package fact_test

import (
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

func evidenceAt(t *testing.T, ref string, at time.Time) provenance.Evidence {
	t.Helper()
	e, err := provenance.NewEvidence("test_system", "item", ref, at, "hash-"+ref)
	if err != nil {
		t.Fatalf("evidence: %v", err)
	}
	return e
}

// TestEqualIgnoresEvidence is the whole point of this comparison. A mirror asks
// "did what I store change?", and the answer must not be "yes" merely because
// the second poll happened a second later. State, value and reason are what a
// reader of the row sees; when those three match, the knowledge is the same
// knowledge however many times it was observed.
func TestEqualIgnoresEvidence(t *testing.T) {
	early := evidenceAt(t, "MLB1", time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC))
	late := evidenceAt(t, "MLB2", time.Date(2026, 8, 11, 18, 30, 0, 0, time.UTC))

	a, err := fact.NewKnown("SKU-1", early)
	if err != nil {
		t.Fatalf("a: %v", err)
	}
	b, err := fact.NewKnown("SKU-1", late)
	if err != nil {
		t.Fatalf("b: %v", err)
	}
	if !fact.Equal(a, b) {
		t.Fatal("same value observed twice compared unequal — every poll would read as a change")
	}
}

func TestEqualSeparatesStateValueAndReason(t *testing.T) {
	e := evidenceAt(t, "MLB1", time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC))

	known1, _ := fact.NewKnown("A", e)
	known2, _ := fact.NewKnown("B", e)
	unknown1, _ := fact.NewUnknown[string]("ml omitted it", e)
	unknown2, _ := fact.NewUnknown[string]("ml sent it blank", e)
	notApplicable, _ := fact.NewNotApplicable[string]("ml omitted it", e)
	estimated, _ := fact.NewEstimated("A", "derived", e)

	cases := []struct {
		name  string
		a, b  fact.Fact[string]
		equal bool
	}{
		{"same value", known1, known1, true},
		{"different value", known1, known2, false},
		// The defect this whole change exists to catch: the mapper was fixed,
		// so the fact went from Unknown to Known while the channel's bytes
		// stayed identical. If state did not count, the row would never heal.
		{"known versus unknown", known1, unknown1, false},
		// Two absences for different reasons are different knowledge: one says
		// the channel omitted the field, the other says it sent it empty.
		{"different reason", unknown1, unknown2, false},
		// Same reason, different state — Unknown means "we do not know" and
		// NotApplicable means "there is nothing to know".
		{"unknown versus not applicable", unknown1, notApplicable, false},
		// Same value, but one is measured and the other is derived.
		{"known versus estimated", known1, estimated, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := fact.Equal(c.a, c.b); got != c.equal {
				t.Fatalf("Equal = %v, want %v", got, c.equal)
			}
			if got := fact.Equal(c.b, c.a); got != c.equal {
				t.Fatalf("Equal reversed = %v, want %v (equality must be symmetric)", got, c.equal)
			}
		})
	}
}

// TestEqualFuncComparesValuesThatCannotUseOperatorEqual covers the types a
// listing actually carries: exact.Money holds a decimal with a pointer inside,
// so == on it compares pointers, not amounts.
func TestEqualFuncComparesValuesThatCannotUseOperatorEqual(t *testing.T) {
	e := evidenceAt(t, "MLB1", time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC))
	type box struct{ digits []int }
	sameDigits := func(x, y box) bool {
		if len(x.digits) != len(y.digits) {
			return false
		}
		for i := range x.digits {
			if x.digits[i] != y.digits[i] {
				return false
			}
		}
		return true
	}

	a, _ := fact.NewKnown(box{digits: []int{1, 9, 9}}, e)
	b, _ := fact.NewKnown(box{digits: []int{1, 9, 9}}, e)
	c, _ := fact.NewKnown(box{digits: []int{2, 0, 0}}, e)

	if !fact.EqualFunc(a, b, sameDigits) {
		t.Fatal("equal values compared unequal")
	}
	if fact.EqualFunc(a, c, sameDigits) {
		t.Fatal("different values compared equal")
	}

	unknown, _ := fact.NewUnknown[box]("absent", e)
	// A stateless comparison would dereference a nil value here.
	if fact.EqualFunc(a, unknown, sameDigits) {
		t.Fatal("a value compared equal to an absence")
	}
	if !fact.EqualFunc(unknown, unknown, sameDigits) {
		t.Fatal("an absence compared unequal to itself")
	}
}
