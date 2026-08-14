# ADR-034: `internal/kernel/fact` supersedes ADR-017

**Date:** 2026-08-07
**Status:** accepted, supersedes ADR-017 (`docs/architecture/decisions/017-unknown-is-never-zero.md`)

## Context

ADR-017 recorded "unknown is never zero" as doctrine, reconstructed on 2026-08-05 from
1.378 live citations of the tags `ADR-17`/`ADR-017`, harvested at
`docs/architecture/decisions/_citations/adr-017-citations.md` (whose own header states
"Total citations (excl. scripts/.runs): 1378"). Enforcement was entirely by citation and
review: `017-unknown-is-never-zero.md` itself records that the rule "governed the
codebase from its first citation ... but no document was ever written," and that it was
"the single most-cited rule in the repository precisely because it was violated
repeatedly and caught in review each time."

This plan built `apps/server_core/internal/kernel/fact` (package `fact`). Its package doc
states the same failure mode the citation count already showed:

> "The repository cites that rule 1378 times and still broke it in its own shared core,
> because a rule enforced by remembering is a rule that gets forgotten."
> — `apps/server_core/internal/kernel/fact/knowledge.go:3-5`

`Knowledge` is a typed enum whose zero value is `Unknown`:

```go
type Knowledge uint8

const (
	// Unknown means we have no value. It is the zero value.
	Unknown Knowledge = iota
	Known
	Estimated
	NotApplicable
)
```
(`knowledge.go:24-36`)

Every constructor validates before a `Fact[T]` can exist: `NewKnown` and `NewEstimated`
reject a zero `provenance.Evidence` (`knowledge.go:73-90`); `NewUnknown` and
`NewNotApplicable` reject both a missing reason and a zero evidence
(`knowledge.go:94-115`). `Fact.Value()` returns `(T, bool)`, and reading only the first
result hands back the zero value of `T` on an `Unknown` fact — exactly the mistake the
package exists to prevent (`knowledge.go:117-127`); a syntactic detector for that call
shape (`ScanFactValueDiscard`) lives in `internal/arch/scan.go:285-322`.

## Decision

**ADR-017 is superseded by `internal/kernel/fact`.** The rule survives — every clause
§1 through §13 of ADR-017 remains a true statement about what the system must do — and
what changes is the level of imposition: doctrine cited 1.378 times, requiring a reviewer
to remember and catch each violation, becomes a type whose invalid construction returns
an error the caller must handle, not a convention a comment asserts.

**§1 — The predicate does not change; the mechanism does.** "An operational fact that
could not be established is represented as absent ... never substituted with `0`, an
empty string, a default" (ADR-017's Decision) is unchanged. What moved is enforcement:
from citation-and-review to `fact.Fact[T]`'s constructors, which cannot produce the
invalid combinations ADR-017 prohibited by prose (`Unknown` with a value, `Known` with
none, a non-Known state with no reason, any fact with no evidence).

**§2 — `Unknown` is the zero value, deliberately, and the real guard is evidence, not the
absence of a literal.** `Fact[T]{}` — the zero value of the generic struct — is `Unknown`
because `Unknown Knowledge = iota` (`knowledge.go:28`), not `Known`. This is a narrower
and more precise claim than "there is no struct literal that can construct a `Fact`": a
zero-valued struct literal, `Fact[string]{}`, compiles in any package, because an empty
struct literal names no fields and therefore violates no `internal`/unexported-field
rule. The actual guard against a confidently-wrong zero is `e.IsZero()`, checked inside
every one of the four constructors (`knowledge.go:74, 86, 99, 111`) — a `Fact` built any
other way than through them either fails validation or is the unusable zero value,
`Unknown`, by construction. Placing `Unknown` rather than `Known` at zero is what makes
that fallback the safe one instead of the confident one.

**§3 — ADR-017 is marked superseded, not deleted or rewritten.** The `Status` field of
`docs/architecture/decisions/017-unknown-is-never-zero.md` changes to `Superseded by
ADR-034`. Its clauses remain the readable statement of the rule for readers who want the
predicate; this document is the record of where enforcement now lives.

## Rationale

A rule enforced only by citation degrades exactly the way ADR-017 itself documents: cited
1.378 times and still broken in the repository's own shared core (`sourcekind.go:47-52`,
per ADR-017's own text and the design spec's Regra 4.1 defect note). Encoding the
predicate in a type does not make the domain knowledge in ADR-017's thirteen clauses
(no fabricated value, no suppressed zero, no silent cross-source fallback, no fabricated
time, and so on) obsolete — those clauses describe *what counts as* an operational fact
being absent across a stock read, a fee, a tax rate. `fact.Fact[T]` enforces the
*shape* — you cannot construct the invalid combination — but does not by itself know
that, say, a Sankhya stock row missing `PERCREDBASEDEST` should be `Unknown` rather than
`Known` with a fabricated zero; that judgment still needs the domain-specific clauses
ADR-017 wrote down. Superseding the document, not deleting it, keeps that knowledge
available while being honest that a compiler-checked type is a stronger guarantee than a
comment citing a decision number.

## Consequences

- `docs/architecture/decisions/017-unknown-is-never-zero.md`'s `Status` line changes to
  `Superseded by ADR-034`; its Context/Decision/clauses are left as written, since they
  remain accurate descriptions of the predicate.
- New code representing an operational fact that might be absent uses
  `apps/server_core/internal/kernel/fact.Fact[T]` (or, where a call site cannot yet
  migrate, states explicitly why not — the type is not retrofitted automatically onto
  every existing `float64`/`*float64` field by this ADR).
- The citation count — 1.378, as harvested for ADR-017 on 2026-08-05 — is not
  re-harvested by this ADR; it is quoted as the historical figure that motivated moving
  enforcement into a type, not as a current live count.
- `fact.Fact[T]`'s arithmetic (`Map`, `Combine2` in `internal/kernel/fact/combine.go`) is not
  addressed by this ADR. The design spec that was going to carry it has been retired with the
  rest of the pre-rebaseline documentary tree, so the question is open and belongs to D2 —
  which ADR-035 and the ADR registry already name as the owner of this primitive's
  application scope.

## Alternatives Considered

**Leave ADR-017 as the sole authority and treat `fact.Fact[T]` as an implementation
detail that satisfies it.** Rejected: the design spec's Regra 4.1 states in writing that
`fact.Fact[T]` "substitui o ADR-017," and the repository's order of truth requires an ADR
before a design claim about superseding a frozen/reconstructed decision is authoritative.
Leaving ADR-017 unmarked would let two documents both claim to be the live rule.

**Delete ADR-017 outright once `fact.Fact[T]` exists.** Rejected: ADR-017's thirteen
clauses encode domain-specific judgments (§3 named-unknown-components, §10 opaque stays
opaque, §13 lenient ingestion) that `fact.Fact[T]`'s generic type does not restate. Losing
that text loses the reasoning, not just the enforcement mechanism.
