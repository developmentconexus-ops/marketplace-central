package fact

// Equal reports whether two facts carry the same knowledge: the same state,
// the same reason, and the same value or the same absence of one.
//
// Evidence is deliberately excluded. Evidence answers "how do we know", and it
// differs on every observation — a second poll of an unchanged listing carries
// a later timestamp and a different reference by construction. A mirror that
// folded evidence into this comparison could never answer the question it
// actually asks, which is whether the row a reader sees has changed.
//
// This is what lets a mirror heal itself. When a mapper is corrected, the
// channel's bytes do not move but the state does — Unknown becomes Known — and
// a fold that compares facts sees that, where a fold that compares payload
// hashes cannot.
func Equal[T comparable](a, b Fact[T]) bool {
	return EqualFunc(a, b, func(x, y T) bool { return x == y })
}

// EqualFunc is Equal for value types that == cannot compare correctly. Money
// is the case that forces it: a decimal amount holds a pointer, so == would
// compare identity rather than amount and two facts stating the same price
// would read as different.
//
// eq is consulted only when both facts hold a value, so it never has to defend
// against an absent one.
func EqualFunc[T any](a, b Fact[T], eq func(T, T) bool) bool {
	if a.state != b.state || a.reason != b.reason {
		return false
	}
	if (a.value == nil) != (b.value == nil) {
		return false
	}
	if a.value == nil {
		return true
	}
	return eq(*a.value, *b.value)
}
