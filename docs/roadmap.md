# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| D6-R2 | [Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) + [P4-R1](engineering/rebaseline/D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md) + [P8](engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md) + [B110 LOCK](engineering/rebaseline/D6-R2-P8-B110-APPROVALS-RATIFICATION.md) + [B10 P6](engineering/rebaseline/D6-R2-P6-B10-PREPARATION-REFERENCE-STUDY.md) — **B00/B01/B00-R2/B11/B12/B110 LOCKED; B10 SUSPENDED; P9 NEXT** |
| NOTIF-01 | [D2-R3](engineering/rebaseline/D6-R2-NOTIF-01-D2-R3-RATIFICATION.md) + [D2-R4](engineering/rebaseline/D6-R2-NOTIF-01-D2-R4-RATIFICATION.md) + [D2-R5](engineering/rebaseline/D6-R2-NOTIF-01-D2-R5-RATIFICATION.md) **ACCEPTED** · [D3-R](engineering/rebaseline/D6-R2-NOTIF-01-D3-R-RATIFICATION.md) **ACCEPTED** · [D3-R1](engineering/rebaseline/D6-R2-NOTIF-01-D3-R1-PRESENTATION-FEED-FORWARD-REVALIDATION.md) + [D3-R2](engineering/rebaseline/D6-R2-NOTIF-01-D3-R2-TYPED-RESULT-CONTINUATION-FEED-FORWARD.md) **PASS** · [D5-F4](engineering/rebaseline/D6-R2-NOTIF-01-D5-F4-RECIPIENT-DISCOVERY-GLOBAL-MAXIMUM.md) + [D5-R3](engineering/rebaseline/D6-R2-NOTIF-01-D5-R3-RATIFICATION.md) **RATIFIED** · [D5-R4](engineering/rebaseline/D6-R2-NOTIF-01-D5-R4-OAD-WIRE-PROOF.md) **PROVED** · [D6-R](engineering/rebaseline/D6-R2-NOTIF-01-D6-R-FRONTEND-FEED-FORWARD.md) + [P8](engineering/rebaseline/D6-R2-NOTIF-01-D6-R-P8-RATIFICATION.md) **LOCKED** · [P9 supersession](engineering/rebaseline/D6-R2-NOTIF-01-D6-R-P9-F1-SUPERSESSION-RATIFICATION.md) **RATIFIED** · [D2-R6](engineering/rebaseline/D6-R2-NOTIF-01-D2-R6-RATIFICATION.md) + [D3-R3](engineering/rebaseline/D6-R2-NOTIF-01-D3-R3-RATIFICATION.md) **ACCEPTED** · [D5-R6](engineering/rebaseline/D6-R2-NOTIF-01-D5-R6-AUTHORIZATION-REQUEST-OAD-WIRE-PROOF.md) **PROVED / CANONICAL** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Derive/prove final P9 Screen Contracts + bidirectional trace; then Fable review/adjudication. B10/D7-R/D8-R remain blocked.** |
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
| D6 — Frontend | **ACCEPTED / CLOSED baseline; B110 operator-LOCKED; final P9 next** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED baseline; NOTIF-01 D7-R BLOCKED BY P9 + FABLE** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; NOTIF-01 D8-R BLOCKED BY D7-R** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — P9 NEXT; FABLE AFTER P9** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- Canonical Product OAD is **106/31**; historical 95/29 + 99/30 proof remains intact.
- AuthorizationRequest D2-R6/D3-R3/D5-R6 is accepted/proved; B110 is operator-`LOCKED`.
- B00-R2/B11/B12 remain operator-`LOCKED`; B10 remains suspended.
- **Fable review + finding adjudication is mandatory after final P9 before Global-Maximum closure / D7-R.**

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.
