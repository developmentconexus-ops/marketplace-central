# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| D6-R2 | [Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) + [P4-R1](engineering/rebaseline/D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md) + [P8](engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md) — **B00/B01 LOCKED; B00-R2 RENDERED CANDIDATE / NOT LOCKED; B10 SUSPENDED** |
| NOTIF-01 | D0-R + D1-R + D2-R/R2 + [D2-R3](engineering/rebaseline/D6-R2-NOTIF-01-D2-R3-RATIFICATION.md) + [D2-R4](engineering/rebaseline/D6-R2-NOTIF-01-D2-R4-RATIFICATION.md) + [D2-R5](engineering/rebaseline/D6-R2-NOTIF-01-D2-R5-RATIFICATION.md) **ACCEPTED** · [D3-R](engineering/rebaseline/D6-R2-NOTIF-01-D3-R-RATIFICATION.md) **ACCEPTED** · [D3-R1](engineering/rebaseline/D6-R2-NOTIF-01-D3-R1-PRESENTATION-FEED-FORWARD-REVALIDATION.md) + [D3-R2](engineering/rebaseline/D6-R2-NOTIF-01-D3-R2-TYPED-RESULT-CONTINUATION-FEED-FORWARD.md) **PASS** · [D5-F4](engineering/rebaseline/D6-R2-NOTIF-01-D5-F4-RECIPIENT-DISCOVERY-GLOBAL-MAXIMUM.md) + [D5-R3](engineering/rebaseline/D6-R2-NOTIF-01-D5-R3-RATIFICATION.md) **OPERATOR-RATIFIED** · [D5-R4](engineering/rebaseline/D6-R2-NOTIF-01-D5-R4-OAD-WIRE-PROOF.md) **PROVED / CANONICAL** · [D6-R](engineering/rebaseline/D6-R2-NOTIF-01-D6-R-FRONTEND-FEED-FORWARD.md) **SPEC OPERATOR-APPROVED; B00-R2 RENDERED CANDIDATE / VISUAL ADJUDICATION NEXT** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **104 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator visually adjudicates only P8 B00-R2 — topbar bell + bounded U01 recent-Inbox preview — from `qualification/d6-r2-wireframes/b00-r2-notifications.html`. Do not render B11/B12 as baseline, begin D7-R/D8-R, resume B10 or implement Product code first.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| ---| --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; **NOTIF-01 D0-R ACCEPTED** |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; **NOTIF-01 D1-R ACCEPTED** |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; **NOTIF-01 D2-R/R2/R3/R4/R5 ACCEPTED** |
| D3 — Communication / Events | ACCEPTED / CLOSED; **NOTIF-01 D3-R ACCEPTED / D3-R1+R2 PASS** |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; **NOTIF-01 D5-R3 ACCEPTED / D5-R4 OAD 104/31 PROVED** |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED; NOTIF-01 D6-R SPEC APPROVED / B00-R2 VISUAL REVIEW ACTIVE** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED; NOTIF-01 D7-R BLOCKED BY D6-R** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; NOTIF-01 D8-R BLOCKED BY D7-R** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B00/B01 LOCKED; B00-R2 RENDERED CANDIDATE / NOT LOCKED; B10 SUSPENDED** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- Canonical OAD is **104/31**: historical 95/29 and 99/30 non-regression proofs survive; current TypeScript/Go projections and split auth proof pass.
- NOTIF-01 focused wire proof passes **104/104 operations · 31/31 Permissions · 8/8 negative controls**.
- D6-R spec is operator-approved. The first sequential P8 block, **B00-R2**, is rendered as executable low-fi HTML; structural proof passes with global sidebar delta `0`, numeric unread count forbidden and B11/B12 not rendered.
- Final visual-gate checkpoint `80f897f2eae593e513cf8292086685095d07c497`: **CI #543 SUCCESS · pr-title #602 SUCCESS · repository full gate PASS**.
- B00-R2 remains **CANDIDATE / NOT LOCKED** until explicit operator visual adjudication of the rendered HTML.
- D6-R→D7-R→D8-R→resume dependent D6-R2. No Product implementation before accepted D9.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.