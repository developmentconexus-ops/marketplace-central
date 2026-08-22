# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D8 — Golden Flows — OPEN / ACTIVE — FABLE ADJUDICATED / LIVE-PROBE DISPOSITION REQUIRED** |
| Accepted baseline | **D0–D7 ACCEPTED / CLOSED** |
| D8 authority | [D8 Golden Flows](engineering/rebaseline/D8-GOLDEN-FLOWS.md) + [D8-R1 Proof Closure & Implementation-Readiness Coherence](engineering/rebaseline/D8-R1-PROOF-CLOSURE-COHERENCE.md) — **OPEN / ACTIVE** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Obtain explicit operator disposition for unresolved D4-deferred D8 probes: authorize bounded real execution/recording or explicitly ratify re-deferral to a named later gate. Do not begin D6-R2, D9 or Product implementation until D8 closes.** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED |
| D3 — Communication / Events | ACCEPTED / CLOSED |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED |
| D6 — Frontend | **ACCEPTED / CLOSED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED** |
| D8 — Golden Flows | **OPEN / ACTIVE — FABLE ADJUDICATED / LIVE-PROBE DISPOSITION REQUIRED** |
| D6-R2 — Complete Frontend Realization Closure | **BLOCKED UNTIL D8 CLOSE** |
| Pre-D9 Implementation Readiness Contract | **BLOCKED UNTIL D6-R2** |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## D8 boundary

D8 composes accepted D0–D7 authority into the smallest representative falsifiable set:

```text
GF-01  Publication & Marketplace Convergence
GF-02  Sale → Business System → Fiscal → Fulfillment → Outcome
GF-03  Performance Evidence Honesty
SR-01  PITR / Timeline Continuity Recovery
```

This is **3 business golden flows + 1 systemic recovery falsifier**, not an exhaustive 99-operation test catalog and not Product implementation. Cross-cutting Organization isolation, auth/access, wire/frontend composition, idempotency/concurrency, durable-effect safety, Governance/Work and knowledge honesty are exercised through these flows rather than becoming extra business journeys by symmetry.

Start from [D8 Golden Flows](engineering/rebaseline/D8-GOLDEN-FLOWS.md). [D8-R1](engineering/rebaseline/D8-R1-PROOF-CLOSURE-COHERENCE.md) supersedes only its bounded proof-closure/revalidation/Governance/dependency-direction seams after independent Fable review and GPT adjudication. Switch to exact prior authority or the canonical OAD only when the candidate flow requires it.

D4-deferred real Mercado Livre/Sankhya consequential probes remain D8 proof obligations under the D8-R1 ledger. Live/irreversible execution requires separate explicit operator authorization; unconditional rows must be executed/recorded or explicitly operator-redeferred before D8 closes. Conditional rows may be recorded `NOT_TRIGGERED` only when their D4 condition is genuinely absent.

## Approved pre-D9 realization sequence

D8 surfaced that **architecture closure is not by itself implementation-readiness closure**. After D8 closes, the smallest accepted sequence is:

```text
D8 close
→ D6-R2 — Complete Frontend Realization Closure
→ Pre-D9 Implementation Readiness Contract
→ D9 — Adversarial Architecture Review
→ Product implementation only after accepted D9
```

D6-R2 will complete material route/screen/state/action/API/Permission/wireframe realization against current accepted Product + D7 authority; it does not reopen accepted frontend technologies/topology by preference. The Pre-D9 Implementation Readiness Contract will project accepted architecture into per-slice outcomes, the complete D6 dependency/import law, negative controls, executable acceptance and explicit exit states so coding agents do not decide material behavior during implementation.

If D6-R2 or the readiness contract legitimately changes an accepted operation, Permission, semantic owner or another material invariant through the smallest owning authority, affected D8 flows/controls must be revalidated before D9 opens. D9 never reviews golden flows derived from superseded authority.

D9 remains blocked until these preconditions close. Product implementation remains blocked until accepted D9. Reopen D0–D7 only for a material falsifier at the smallest owning authority.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).
