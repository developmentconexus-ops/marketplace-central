# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| D6-R2 | [Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) + [P4-R1](engineering/rebaseline/D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md) + [P8](engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md) — **B00/B01 LOCKED; B10 SUSPENDED** |
| NOTIF-01 | D0-R + D1-R + D2-R + D2-R2 + [D2-R3](engineering/rebaseline/D6-R2-NOTIF-01-D2-R3-RATIFICATION.md) + [D3-R](engineering/rebaseline/D6-R2-NOTIF-01-D3-R-RATIFICATION.md) **ACCEPTED** · [D5-F2](engineering/rebaseline/D6-R2-NOTIF-01-D5-F2-INBOX-PRESENTATION-CONTEXT-FINDING.md) proved Inbox presentation gap · [D2-R4](engineering/rebaseline/D6-R2-NOTIF-01-D2-R4-PRESENTATION-SNAPSHOT.md) **CANDIDATE** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicates only NOTIF-01 D2-R4 presentation snapshot. D3 feed-forward + D5-R3 table/OAD remain blocked.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; **NOTIF-01 D0-R ACCEPTED** |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; **NOTIF-01 D1-R ACCEPTED** |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; **NOTIF-01 D2-R + D2-R2 + D2-R3 ACCEPTED / D2-R4 CANDIDATE** |
| D3 — Communication / Events | ACCEPTED / CLOSED; **NOTIF-01 D3-R ACCEPTED / presentation feed-forward BLOCKED BY D2-R4** |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; **NOTIF-01 D5-R3 analysis OPEN / table not ratified** |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B00/B01 LOCKED; B10 SUSPENDED; NOTIF-01 D2-R4 gate** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- D0-R, D1-R, 14 family contracts, D2-R/R2/R3 and D3-R are accepted.
- D5-R3 still derives four Product operations and one candidate `notifications.manage` Permission; OAD remains unchanged.
- D5-F2 proves same-kind Inbox rows need one immutable notification-safe `subject_display_label` to remain human-distinguishable without source-per-row N+1 reads.
- D2-R4 candidate adds only that label: source-owner-derived, immutable, non-sensitive, non-authoritative, never used for identity/routing/dedup/navigation/automation; final title remains `NotificationKind`-derived/localized.
- No generic payload/template, second summary field, source-state copy, broker, runtime or Product implementation.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.