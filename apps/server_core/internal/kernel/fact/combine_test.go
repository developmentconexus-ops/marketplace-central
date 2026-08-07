package fact_test

import (
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
)

// TestCombine2PropagatesTheWorstState pins the rule that makes the kernel worth
// having: a calculation with one unknown input has an unknown output. Solving a
// target-margin price with one unknown component U gives P = 178.69 + 2.02*U —
// an ignored cost does not err slightly, it errs by twice what was ignored.
func TestCombine2PropagatesTheWorstState(t *testing.T) {
	add := func(a, b exact.Decimal) (exact.Decimal, error) { return a.Add(b), nil }

	t.Run("known + known -> known", func(t *testing.T) {
		a, err := fact.NewKnown(exact.FromInt(2), ev(t))
		if err != nil {
			t.Fatalf("NewKnown a: %v", err)
		}
		b, err := fact.NewKnown(exact.FromInt(3), ev(t))
		if err != nil {
			t.Fatalf("NewKnown b: %v", err)
		}

		out, err := fact.Combine2(a, b, "sum", add)
		if err != nil {
			t.Fatalf("Combine2: %v", err)
		}
		if out.State() != fact.Known {
			t.Fatalf("State() = %v, want Known", out.State())
		}
		v, ok := out.Value()
		if !ok {
			t.Fatal("Value() ok = false for a Known result")
		}
		if v.Cmp(exact.FromInt(5)) != 0 {
			t.Fatalf("Value() = %s, want 5", v.StringFixed(0))
		}
	})

	t.Run("known + estimated -> estimated, carrying the estimate's reason", func(t *testing.T) {
		a, err := fact.NewKnown(exact.FromInt(2), ev(t))
		if err != nil {
			t.Fatalf("NewKnown: %v", err)
		}
		b, err := fact.NewEstimated(exact.FromInt(3), "modelled from last 30 days", ev(t))
		if err != nil {
			t.Fatalf("NewEstimated: %v", err)
		}

		out, err := fact.Combine2(a, b, "sum", add)
		if err != nil {
			t.Fatalf("Combine2: %v", err)
		}
		if out.State() != fact.Estimated {
			t.Fatalf("State() = %v, want Estimated", out.State())
		}
		if !strings.Contains(out.Reason(), "modelled from last 30 days") {
			t.Fatalf("Reason() = %q, does not carry the estimated input's reason", out.Reason())
		}
	})

	t.Run("known + unknown -> unknown, fn never called", func(t *testing.T) {
		a, err := fact.NewKnown(exact.FromInt(2), ev(t))
		if err != nil {
			t.Fatalf("NewKnown: %v", err)
		}
		b, err := fact.NewUnknown[exact.Decimal]("marketplace did not return it", ev(t))
		if err != nil {
			t.Fatalf("NewUnknown: %v", err)
		}
		ran := false
		fn := func(x, y exact.Decimal) (exact.Decimal, error) {
			ran = true
			return x.Add(y), nil
		}

		out, err := fact.Combine2(a, b, "sum", fn)
		if err != nil {
			t.Fatalf("Combine2: %v", err)
		}
		if out.State() != fact.Unknown {
			t.Fatalf("State() = %v, want Unknown", out.State())
		}
		if ran {
			t.Fatal("fn ran even though one input was unknown")
		}
		if _, ok := out.Value(); ok {
			t.Fatal("Unknown result returned a value")
		}
	})

	t.Run("known + not_applicable -> not_applicable", func(t *testing.T) {
		a, err := fact.NewKnown(exact.FromInt(2), ev(t))
		if err != nil {
			t.Fatalf("NewKnown: %v", err)
		}
		b, err := fact.NewNotApplicable[exact.Decimal]("CST 60: no own ICMS on this operation", ev(t))
		if err != nil {
			t.Fatalf("NewNotApplicable: %v", err)
		}

		out, err := fact.Combine2(a, b, "sum", add)
		if err != nil {
			t.Fatalf("Combine2: %v", err)
		}
		if out.State() != fact.NotApplicable {
			t.Fatalf("State() = %v, want NotApplicable", out.State())
		}
		if _, ok := out.Value(); ok {
			t.Fatal("NotApplicable result returned a value")
		}
	})
}

// TestCombine2NeverCallsFnOnAnUnusableInput is the negative control: a fn that
// records it ran must not have run.
func TestCombine2NeverCallsFnOnAnUnusableInput(t *testing.T) {
	a, err := fact.NewUnknown[exact.Decimal]("no data", ev(t))
	if err != nil {
		t.Fatalf("NewUnknown: %v", err)
	}
	b, err := fact.NewKnown(exact.FromInt(5), ev(t))
	if err != nil {
		t.Fatalf("NewKnown: %v", err)
	}

	ran := false
	fn := func(x, y exact.Decimal) (exact.Decimal, error) {
		ran = true
		return x.Add(y), nil
	}

	if _, err := fact.Combine2(a, b, "sum", fn); err != nil {
		t.Fatalf("Combine2: %v", err)
	}
	if ran {
		t.Fatal("fn ran even though one input was unusable")
	}
}
