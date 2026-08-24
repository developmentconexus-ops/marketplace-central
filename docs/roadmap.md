# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| D6-R2 | [Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) + [AuthorizationRequest closure](engineering/rebaseline/D6-R2-AUTHORIZATION-REQUEST-FABLE-RATIFICATION.md) + [D7-R](engineering/rebaseline/D6-R2-NOTIF-01-D7-R-AUTHORIZATION-REQUEST-RUNTIME-REPAIR.md) + [B10 preparation](engineering/rebaseline/D6-R2-P6-B10-PREPARATION-REFERENCE-STUDY.md) — **AUTH GLOBAL-MAXIMUM CLOSED; D7-R ACCEPTED; D8-R NEXT; B10 SUSPENDED** |
| NOTIF-01 | D2-R6 + D3-R3 + D5-R6 + final P9 + Fable closure **ACCEPTED / PROVED**; [D7-R runtime repair](engineering/rebaseline/D6-R2-NOTIF-01-D7-R-AUTHORIZATION-REQUEST-RUNTIME-REPAIR.md) **OPERATOR-RATIFIED / ACCEPTED** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Open bounded NOTIF-01 D8-R against accepted D7-R and revalidate affected golden-flow/control obligations. D8-R is NEXT / NOT STARTED; do not resume B10, merge PR #61 or implement Product code first.** |
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
| D5 — API | ACCEPTED / CLOSED; **D5-R6 AUTHORIZATIONREQUEST OAD 106/31 PROVED / CANONICAL** |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED baseline; B110 LOCKED; final P9 PROVED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED baseline; NOTIF-01 D7-R OPERATOR-RATIFIED / ACCEPTED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; NOTIF-01 D8-R NEXT / NOT STARTED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — AUTH GLOBAL-MAXIMUM CLOSED; D7-R ACCEPTED; D8-R NEXT; B10 SUSPENDED** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- Product remains **106/31**; historical 95/29 + 99/30 proof remains green.
- AuthorizationRequest redesign + independent Fable closure is operator-ratified as Global Maximum under present requirements.
- D7-R is **OPERATOR-RATIFIED / ACCEPTED** without new Product surface or infrastructure: `CreateAuthorizationDecision` idempotency is specialized to `Organization + effective Principal + operation + key`; exact committed replay precedes stale `If-Match`; current eligibility/validity, typed semantic 503, F13/F14, zero-decider Work and invalidation recovery are explicit runtime laws.
- The D7-R repository verifier proves structural authority only. Real PostgreSQL/River/runtime proof remains mandatory after the implementation gate opens.
- D8-R is **NEXT / NOT STARTED**. B10 remains suspended; Product implementation remains blocked until accepted D9.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.
