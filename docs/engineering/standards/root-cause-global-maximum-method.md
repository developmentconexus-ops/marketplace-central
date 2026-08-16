# Root-Cause / Global-Maximum Engineering Method

> **Status:** canonical engineering method  
> **Primary reader:** LLM-assisted engineering  
> **Applies to:** architecture, design, debugging, refactoring, implementation planning and review

## Objective

> **Find the smallest sustainable solution that preserves essential complexity, removes accidental complexity, resolves the root cause and avoids foreseeable structural dead ends.**

Do not optimize for agreement, familiarity, minimum code, maximum abstraction, framework preference or preservation of the current implementation.

## Required decision flow

For every material decision:

```text
Evidence
→ Known / Inferred / Unknown / Deferred
→ Root Cause
→ Target Invariant
→ Constraints
→ Credible Alternatives
→ Local Maximum vs Global Maximum
→ Essential vs Accidental Complexity
→ YAGNI / Overengineering
→ Future-Cost Test
→ Authority / Boundary when relevant
→ Strongest Reasonable Enforcement
→ Proof
→ Adversarial Review
→ Decision
→ Reopen Triggers
```

After a meaningful group of decisions, perform a **Global Coherence Review**.

## 1. Evidence first

Start from what is known, not from a preferred solution.

Classify material information as:

- **Known:** directly supported by reliable evidence.
- **Inferred:** reasoned from known evidence.
- **Unknown:** material but unresolved.
- **Deferred:** intentionally left to later work.

Unknown must not become a convenient default.

Current code, schemas, package layout, APIs, tests, runtime and historical decisions are evidence, not target authority merely because they exist.

Use the **Structural Inversion Test**:

> If the current implementation were completely different, which parts of this conclusion would still be true?

Prefer primary/official sources for unstable external or technical facts.

## 2. Root cause before correction

Do not patch symptoms before identifying the structural condition that made the defect possible.

Ask:

- What structural fact made this possible?
- Can it produce other failures?
- Would a local fix leave the defect class reachable?

Then state an implementation-independent **target invariant**: what must remain true in every valid implementation.

Bad: `Create ValidationService.`  
Good: `Invalid state cannot enter committed state through any supported path.`

The invariant is durable; the implementation is replaceable.

## 3. Local Maximum vs Global Maximum

**Local Maximum:** best solution inside the current structure.

**Global Maximum:** best sustainable structure for the real constraints, even when the current structure must change.

Prefer the Global Maximum when the current structure preserves the root cause.

Global Maximum does **not** mean maximum abstraction, infrastructure, generality or future-proofing.

## 4. Essential vs accidental complexity

**Essential complexity** comes from the real problem: uncertainty, concurrency, authorization, isolation, auditability, partial failure, idempotency, temporal correctness, external dependencies, contract evolution and similar constraints.

**Accidental complexity** comes from the chosen solution: duplicate authorities, hand-synced models, unnecessary indirection, speculative frameworks, repeated enforcement and compatibility without consumers.

> **Remove accidental complexity. Never simplify correctness by flattening essential complexity.**

## 5. YAGNI, overengineering and future cost

YAGNI means **do not build capability without justified need**. It may remove speculative frameworks, hypothetical integrations, unused extensibility, generic engines and abstractions with no real responsibility.

YAGNI must not remove a known invariant, safety property, required isolation/recoverability/auditability, evidence/provenance, or a seam already justified by evidenced evolution.

Before adding abstraction, ask:

1. What concrete problem does it solve?
2. What defect class does it eliminate?
3. Is that defect reachable or evidenced?
4. Is more than one real use case justifying the abstraction?
5. Are we generalizing from one example?
6. Could it be added later without dismantling existing authority?
7. Does it reduce total complexity or only move it?
8. Is a simpler existing boundary sufficient?

Before deliberately choosing a simpler structure, run the **Future-Cost Test**:

> If foreseeable evolution already supported by evidence occurs, can this design extend additively, or will it require dismantling authority, rewriting core contracts or duplicating semantics?

Prefer:

> **Prepare the seam, not the entire future capability.**

## 6. Mechanism is not authority

A shared mechanism does not automatically own the meaning it supports.

Retries, scheduling, idempotency, policy lookup, validation primitives, observability, audit transport, caching, serialization and event transport may be centralized when that reduces repeated correctness work.

> **Centralize repeated mechanisms when justified. Keep meaning and decisions with the component that actually understands them.**

When ownership matters, explicitly state who owns the meaning/lifecycle, what remains external, what callers may depend on and what the boundary does not own. Two authorities for the same meaning are presumed wrong until justified.

## 7. Strongest reasonable enforcement

After defining the invariant, prefer enforcement in this order:

1. structure / public boundary;
2. type system;
3. schema / database constraint;
4. runtime fail-closed boundary;
5. executable test;
6. static guard / lint;
7. documentation / convention.

Use a weaker layer when a stronger one creates disproportionate complexity or encodes assumptions that are not universally valid.

## 8. Proof before implementation

Define proof before implementation. Artifact existence is not proof:

- test exists ≠ relevant path executed;
- mock passes ≠ integration works;
- success response ≠ convergence;
- guard exists ≠ guard can fail;
- configuration exists ≠ enforcement works.

Prefer proof that can fail for the exact protected property: compile failure, type/schema rejection, negative fixture, real integration test, restart/recovery test, concurrency test, contract diff or end-to-end failure-path exercise.

> **A control that cannot be shown to fire is not proven.**

## 9. Adversarial and global review

Before accepting a material decision, ask:

- What is the strongest argument against it?
- Which assumption would invalidate it?
- Is there hidden duplicate authority?
- Are we fitting the current implementation instead of the real problem?
- Are we overfitting one framework, provider, database or use case?
- Are we using YAGNI to justify underengineering?
- Are we using future-proofing to justify overengineering?
- What happens under uncertainty, partial failure, concurrency or restart?
- What will be hardest to change later?

Periodically review the whole system for duplicate/missing/circular ownership, contradictory assumptions, repeated mechanisms, God components, excessive fragmentation, abstractions caused only by other abstractions, missing seams and speculative extensibility.

> **Local correctness does not guarantee global coherence.**

## 10. Decision, reopen and stop

Every material decision ends as one of:

- **RESTRUCTURE NOW** — implement the Global Maximum now.
- **CURRENT STRUCTURE CONFIRMED** — structure is sound; make the bounded correction.
- **TRANSITIONAL SOLUTION** — temporary local maximum with named successor and deletion condition.
- **STOP / SPLIT PREREQUISITE** — another unresolved decision blocks correctness.
- **DEFER SAFELY** — current work can proceed without deciding this detail.

Do not repeatedly reopen accepted decisions for preference or hypothetical futures. Record concrete **reopen triggers**: new real use case, changed scale/ownership, newly reachable failure mode, or external requirement that invalidates the assumptions.

Stop when evidence is sufficient; root cause and invariant are explicit; credible alternatives were compared; Global Maximum, YAGNI, overengineering and future cost were checked; ownership is clear where relevant; proof is defined; strongest objections were addressed; no material contradiction remains; and reopen triggers are recorded.

## Compact decision record

```markdown
## Decision: <name>
Status: PROPOSED | ACCEPTED | DEFERRED | REOPENED

### Evidence
### Known / Inferred / Unknown / Deferred
### Root Cause
### Target Invariant
### Alternatives
### Global Maximum
### YAGNI / Overengineering / Future Cost
- Remove:
- Preserve:
- Prepare seam:
- Do not build yet:
### Authority / Boundary
Only when relevant.
### Enforcement & Proof
### Adversarial Findings
### Decision
RESTRUCTURE NOW | CURRENT STRUCTURE CONFIRMED | TRANSITIONAL SOLUTION | STOP / SPLIT PREREQUISITE | DEFER SAFELY
### Reopen Triggers
```

Keep the record proportional. Do not add ceremony that does not improve the decision.
