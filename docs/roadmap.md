# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Integrated checkpoints | **PR #61 → `b54d17bfe6d794645d198a9160f4a2a1c63647e8`; PR #62 methodology → `689dab34b0a756cbd7c790a6c5277d887ced0b4c`; PR #65 proportional CI → `a2aeac19816c90ee30bf373cef0448d52a486c7e`; PR #68 PublicationRequirements wire → `ed3d164b0574b7950c2c7467d150c89576bba1ec`** |
| Method profile | **local [`engineering-method.md`](development/engineering-method.md) + [`frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md) v2.3** |
| Current acceptance increment | **PR #64 / B10 — P8 OPERATOR-RATIFIED / LOCKED; A01 `ACCEPT_FOR_LOCK_WITH_LATER_PROBE`; P9 PASS / CLOSED; P10 COMPLETE / NO NEW SHARED PATTERN; integration pending** |
| B10 Global Maximum | **requirements + source values + downstream authoring/provider validation; `source_sufficiency` REJECTED; NO NEW UPSTREAM WIRE FIELD** |
| B10 evidence | **[P6 / Global Maximum revalidation](engineering/rebaseline/D6-R2-P6-B10-PREPARATION-REFERENCE-STUDY.md) + [P8 operator LOCK](engineering/rebaseline/D6-R2-P8-B10-PREPARATION-RATIFICATION.md) + [P9 Screen Contract PASS](engineering/rebaseline/D6-R2-P9-B10-PREPARATION-SCREEN-CONTRACT.md)** |
| Resolved upstream finding | **PublicationRequirements wire gap remains RESOLVED by integrated PR #68; B10 P9 found NO new upstream gap** |
| LOCK impact | **B10 re-LOCKED on the simplified requirement/value projection; B00 / B01 / B00-R2 / B11 / B12 / B110 remain unaffected** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Aggregate CI | **one required check / one `npm run gate`; closure commit must be green before integration** |
| Exact next action | **Verify the `required` check on the B10 closure commit, then obtain explicit operator authorization before integrating PR #64. Do not begin the next dependent B20 increment before PR #64 lands.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; NOTIF-01 D0-R ACCEPTED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; D1-R2 PASS |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; D2-R6 ACCEPTED |
| D3 — Communication / Events | ACCEPTED / CLOSED; D3-R3 ACCEPTED |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; D5-R6 106/31 PROVED; D5-R7 W1 REPAIR ACCEPTED |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110 LOCKED; B10 re-LOCKED in D6-R2** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED baseline; D7-R ACCEPTED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; D8-R ACCEPTED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B10 P8 LOCKED; P9 PASS/CLOSED; P10 COMPLETE; PR #64 integration pending** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- PR #68 remains integrated; Product **106/31/H-A-S**, runtime NONE and unrelated LOCKs are unchanged.
- P9 finding `F-P9-B10-01` correctly falsified the earlier `Atendido` projection and triggered the operator-requested Global Maximum re-evaluation.
- Fresh AnyMarket, Hub2b, Magis5 and Mercado Livre evidence supports the simpler hub model: marketplace requirements + available source values + downstream listing authoring/provider validation.
- Generic `source_sufficiency` remains rejected as accidental complexity. No Product operation, Permission, Principal kind or new Product wire field was added.
- N03/UF03/P3/P4/P5 were boundedly aligned with the simpler job: prepare exact source data for marketplace authoring rather than decide synthetic publishability/readiness per field.
- The operator approved the browser-operable simplified B10 candidate on 2026-08-25; **P8 is OPERATOR-RATIFIED / LOCKED** and A01 remains `ACCEPT_FOR_LOCK_WITH_LATER_PROBE`.
- P9 now binds route/state/identity, owners, reads/writes, failures, access/disclosure and both trace directions. **BACKEND SUFFICIENT / UPSTREAM FINDING NONE / P9 PASS.**
- P10 compared B10 with prior LOCKED evidence. It reuses the B00 exact-context/invalidation model, the existing known-empty-vs-unavailable honesty law and the existing navigation-handoff-vs-mutation law. Search→detail, requirement table and correspondence remain B10-local. **No new shared component/pattern authority is justified.**
- CI remains one objective `required` aggregate check; no planning-string or ratification meta-tests were reintroduced.

```text
F-P9-B10-01
→ GLOBAL MAXIMUM REVALIDATED
→ simplified B10 P8
→ operator LOCK
→ P9 bidirectional Screen Contract PASS
→ P10 NO NEW SHARED PATTERN
→ required CI on closure commit
→ operator-authorized integration of PR #64
→ only then next dependent B20 increment
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.