# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| D6-R2 | [Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) + [P4-R1](engineering/rebaseline/D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md) + [P8](engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md) + [B10 P6](engineering/rebaseline/D6-R2-P6-B10-PREPARATION-REFERENCE-STUDY.md) — **B00/B01/B00-R2/B11/B12 LOCKED; B10 SUSPENDED; NOTIF-01 P9 OPEN / BLOCKED BY P9-F1** |
| NOTIF-01 | [D2-R3](engineering/rebaseline/D6-R2-NOTIF-01-D2-R3-RATIFICATION.md) + [D2-R4](engineering/rebaseline/D6-R2-NOTIF-01-D2-R4-RATIFICATION.md) + [D2-R5](engineering/rebaseline/D6-R2-NOTIF-01-D2-R5-RATIFICATION.md) **ACCEPTED** · [D3-R](engineering/rebaseline/D6-R2-NOTIF-01-D3-R-RATIFICATION.md) **ACCEPTED** · [D3-R1](engineering/rebaseline/D6-R2-NOTIF-01-D3-R1-PRESENTATION-FEED-FORWARD-REVALIDATION.md) + [D3-R2](engineering/rebaseline/D6-R2-NOTIF-01-D3-R2-TYPED-RESULT-CONTINUATION-FEED-FORWARD.md) **PASS** · [D5-F4](engineering/rebaseline/D6-R2-NOTIF-01-D5-F4-RECIPIENT-DISCOVERY-GLOBAL-MAXIMUM.md) + [D5-R3](engineering/rebaseline/D6-R2-NOTIF-01-D5-R3-RATIFICATION.md) **OPERATOR-RATIFIED** · [D5-R4](engineering/rebaseline/D6-R2-NOTIF-01-D5-R4-OAD-WIRE-PROOF.md) **PROVED** · [D6-R](engineering/rebaseline/D6-R2-NOTIF-01-D6-R-FRONTEND-FEED-FORWARD.md) + [P8 ratification](engineering/rebaseline/D6-R2-NOTIF-01-D6-R-P8-RATIFICATION.md) **APPROVED/LOCKED** · [P9-F1](engineering/rebaseline/D6-R2-NOTIF-01-D6-R-P9-F1-ACTIONABLE-GOVERNANCE-CONTEXT.md) **MATERIAL FALSIFIER / OPERATOR ADJUDICATION REQUIRED** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **104 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicates only NOTIF-01 P9-F1: one purpose-bounded Governance actionable-authorization projection for current human decision principals, reused by `/aprovacoes` and F13. OAD remains 104/31 until ratified. Do not freeze P9, begin D7-R/D8-R, resume B10 or implement Product code first.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | **ACCEPTED / CLOSED; NOTIF-01 D0-R ACCEPTED** |
| D1 — Domains / Boundaries | **ACCEPTED / CLOSED; NOTIF-01 D1-R ACCEPTED** |
| D2 — Identity / Tenant / Data Ownership | **ACCEPTED / CLOSED; NOTIF-01 D2-R/R2/R3/R4/R5 ACCEPTED** |
| D3 — Communication / Events | **ACCEPTED / CLOSED; NOTIF-01 D3-R ACCEPTED / D3-R1+R2 PASS** |
| D4 — External Integrations | **ACCEPTED / CLOSED** |
| D4-R1 — Publication Input / Listing Authoring | **ACCEPTED / CANONICAL** |
| D5 — API | **ACCEPTED / CLOSED; NOTIF-01 D5-R3 ACCEPTED / D5-R4 OAD 104/31 PROVED** |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED; NOTIF-01 P8 LOCKED / P9-F1 OPEN** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED; NOTIF-01 D7-R BLOCKED BY D6-R** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; NOTIF-01 D8-R BLOCKED BY D7-R** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — NOTIF-01 P9 BLOCKED BY P9-F1** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- Canonical OAD remains **104/31**; historical 95/29 + 99/30 proof remains intact.
- B00-R2, B11 and B12 are operator-`LOCKED`; reviewed HTML snapshots remain unchanged evidence.
- P9 found a Governance operability gap: F13 carries target identity, while `CreateAuthorizationDecision` requires exact current target revision and no current actionable-decision read exists.
- W4 keeps `governance.decide`, `governance.read` and target-owner reads independent; hidden Permission implication is forbidden.
- P9-F1 recommends one H-only `governance.decide` purpose-bounded projection, provisionally `ListMyActionableAuthorizations`; if ratified, consequence is **105 operations / 31 Permissions**.
- No P8 layout reopen is required. D7-R/D8-R and Product implementation remain blocked.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.
