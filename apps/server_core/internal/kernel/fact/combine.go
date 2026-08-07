package fact

import (
	"fmt"
	"strings"

	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

// Map applies fn to a usable fact and propagates everything else untouched.
//
// This is a package-level function and not a method because a Go method cannot
// declare type parameters of its own: `func (f Fact[A]) Map[B any](...)` does
// not compile. The design said the arithmetic lives in methods of Fact[T]; the
// language says otherwise, and the language wins.
func Map[A, B any](f Fact[A], method string, fn func(A) (B, error)) (Fact[B], error) {
	v, ok := f.Value()
	if !ok {
		return propagate[A, B](f, method)
	}
	out, err := fn(v)
	if err != nil {
		return Fact[B]{}, err
	}
	e, err := provenance.Derived(method, f.Evidence())
	if err != nil {
		return Fact[B]{}, err
	}
	if f.State() == Estimated {
		return NewEstimated(out, f.Reason(), e)
	}
	return NewKnown(out, e)
}

// Combine2 folds two facts into one.
//
// fn is NOT called unless BOTH inputs are usable. That is the whole guarantee:
// there is no code path in which a calculation sees the zero value of an unknown
// input, because the calculation is never reached.
//
// The output state is the worst of the two inputs, in this order:
// NotApplicable < Unknown < Estimated < Known. NotApplicable dominates because
// "this quantity does not exist here" makes the derived quantity not exist
// either, which is a stronger statement than "we failed to learn it".
func Combine2[A, B, C any](a Fact[A], b Fact[B], method string, fn func(A, B) (C, error)) (Fact[C], error) {
	av, aok := a.Value()
	bv, bok := b.Value()
	if !aok || !bok {
		return propagate2[A, B, C](a, b, method)
	}
	out, err := fn(av, bv)
	if err != nil {
		return Fact[C]{}, err
	}
	e, err := provenance.Derived(method, a.Evidence(), b.Evidence())
	if err != nil {
		return Fact[C]{}, err
	}
	if a.State() == Estimated || b.State() == Estimated {
		return NewEstimated(out, joinReasons(method, a, b), e)
	}
	return NewKnown(out, e)
}

func propagate[A, B any](f Fact[A], method string) (Fact[B], error) {
	e, err := provenance.Derived(method, f.Evidence())
	if err != nil {
		return Fact[B]{}, err
	}
	reason := fmt.Sprintf("%s: input is %s (%s)", method, f.State(), f.Reason())
	if f.State() == NotApplicable {
		return NewNotApplicable[B](reason, e)
	}
	return NewUnknown[B](reason, e)
}

func propagate2[A, B, C any](a Fact[A], b Fact[B], method string) (Fact[C], error) {
	e, err := provenance.Derived(method, a.Evidence(), b.Evidence())
	if err != nil {
		return Fact[C]{}, err
	}
	reason := joinReasons(method, a, b)
	if a.State() == NotApplicable || b.State() == NotApplicable {
		return NewNotApplicable[C](reason, e)
	}
	return NewUnknown[C](reason, e)
}

// joinReasons keeps BOTH inputs' reasons. An operator told only the first reason
// fixes one source and sees the value stay unknown with no idea why.
func joinReasons[A, B any](method string, a Fact[A], b Fact[B]) string {
	parts := []string{method + ":"}
	for _, s := range []struct {
		state  Knowledge
		reason string
	}{{a.State(), a.Reason()}, {b.State(), b.Reason()}} {
		if s.state != Known {
			parts = append(parts, fmt.Sprintf("%s (%s)", s.state, s.reason))
		}
	}
	return strings.Join(parts, " ")
}
