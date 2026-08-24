# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| D6-R2 | [Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) + [authority route](engineering/rebaseline/D6-R2-AUTHORITY-ROUTE.md) + [AuthorizationRequest closure](engineering/rebaseline/D6-R2-AUTHORIZATION-REQUEST-FABLE-RATIFICATION.md) + [D7-R](engineering/rebaseline/D6-R2-NOTIF-01-D7-R-AUTHORIZATION-REQUEST-RUNTIME-REPAIR.md) + [D8-R](engineering/rebaseline/D6-R2-NOTIF-01-D8-R-AUTHORIZATION-REQUEST-GOLDEN-FLOW-REVALIDATION.md) + [D5-R7 W1 repair](engineering/rebaseline/D6-R2-FABLE-R1-D5-R7-AUTHORIZATION-DECISION-W1-CARRIER-REPAIR.md) — **FABLE WHOLE-REPO REVIEW OPEN REMEDIATION; R-1 ACCEPTED/REPAIRED; R-2 NEXT; B10 PAUSED** |
| NOTIF-01 | D2-R6 + D3-R3 + D5-R6 + final P9 + Fable closure + D7-R + D8-R **ACCEPTED / PROVED**; D5-R7 supersedes only `CreateAuthorizationDecision` carrier-specific wire meaning |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Decision carrier | **`CreateAuthorizationDecision`: `Idempotency-Key` header + typed body `etag`/`outcome`; missing/invalid revision proof 422, stale revision proof 409** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Close Fable R-2 OAD source-orphan hygiene with a bounded source-reachability proof and pruning only of proven-unreachable superseded source symbols. Preserve the canonical 106/31 bundle. Do not resume B10, merge PR #61 or implement Product code first.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; NOTIF-01 D0-R ACCEPTED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; D1-R2 PASS |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED baseline; **D2-R6 OPERATOR-RATIFIED / ACCEPTED** |
| D3 — Communication / Events | ACCEPTED / CLOSED baseline; **D3-R3 OPERATOR-RATIFIED / ACCEPTED** |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; **D5-R6 106/31 PROVED; D5-R7 W1 CARRIER REPAIR OPERATOR-RATIFIED / ACCEPTED** |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110 LOCKED; final P9 PROVED; B110 REVALIDATED / STRUCTURE UNAFFECTED by D5-R7** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED baseline; D7-R ACCEPTED; carrier-specific wording superseded only by D5-R7** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; D8-R ACCEPTED; carrier-specific wording superseded only by D5-R7** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — WHOLE-REPO FABLE VERDICT ACCEPT WITH BOUNDED FIXES; R-1 CLOSED; R-2 NEXT; B10 PAUSED** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- Product remains **106/31**; historical 95/29 + 99/30 proof remains required.
- AuthorizationRequest model C, D7-R and D8-R remain accepted Global-Maximum authority; D5-R7 changes **only** the `:decide` revision carrier/status grammar to canonical W1 typed-request semantics.
- Exact committed replay remains before current revision-precondition evaluation; Principal-scoped idempotency, typed semantic 503, F13/F14, zero-decider Work, invalidation recovery and PITR laws are unchanged.
- The D5-R7 LOCK impact sweep marks B00/B01/B00-R2/B11/B12 `UNAFFECTED` and B110 `REVALIDATE → STRUCTURE UNAFFECTED`; no P8 reopen is required.
- Independent whole-repository Fable review verdict is **ACCEPT WITH BOUNDED FIXES**. R-2 source-orphan hygiene is the next bounded repository repair. Global methodology evolution is handled separately and must be reconciled before B10 resumes.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.
