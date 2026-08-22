# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D8 — Golden Flows — OPEN / ACTIVE — DERIVED CANDIDATE** |
| Accepted baseline | **D0–D7 ACCEPTED / CLOSED** |
| D8 authority | [D8 Golden Flows](engineering/rebaseline/D8-GOLDEN-FLOWS.md) — **OPEN / ACTIVE — DERIVED CANDIDATE** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Independently challenge the derived D8 candidate, then adjudicate only material findings. Do not begin D6-R2, D9 or Product implementation; no live Mercado Livre or irreversible Sankhya write without separate explicit operator authorization.** |
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
| D8 — Golden Flows | **OPEN / ACTIVE — DERIVED CANDIDATE** |
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

Start from [D8 Golden Flows](engineering/rebaseline/D8-GOLDEN-FLOWS.md), then switch only to the exact accepted owner or canonical OAD needed for the candidate flow under review.

D4-deferred real Mercado Livre/Sankhya consequential probes remain D8 proof obligations where the accepted owner requires them, but any live/irreversible write requires separate explicit operator authorization. D8 architecture proof does not claim an implemented Product runtime.

## Approved pre-D9 realization sequence

D8 surfaced that **architecture closure is not by itself implementation-readiness closure**. After D8 closes, the smallest accepted sequence is:

```text
D8 close
→ D6-R2 — Complete Frontend Realization Closure
→ Pre-D9 Implementation Readiness Contract
→ D9 — Adversarial Architecture Review
→ Product implementation only after accepted D9
```

D6-R2 will complete material route/screen/state/action/API/Permission/wireframe realization against current 99/30 + D7 authority; it does not reopen accepted frontend technologies/topology by preference. The Pre-D9 Implementation Readiness Contract will project accepted architecture into per-slice outcomes, dependency/import rules, negative controls, executable acceptance and explicit exit states so coding agents do not decide material behavior during implementation.

D9 remains blocked until these preconditions close. Product implementation remains blocked until accepted D9. Reopen D0–D7 only for a material falsifier at the smallest owning authority.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).
