# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Integrated checkpoints | **PR #61 Product/frontend checkpoint → `b54d17bfe6d794645d198a9160f4a2a1c63647e8`; PR #62 methodology adoption → `689dab34b0a756cbd7c790a6c5277d887ced0b4c`; PR #65 proportional CI verification → `a2aeac19816c90ee30bf373cef0448d52a486c7e`** |
| Method profile | **`developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535`** |
| Current prerequisite increment | **PR #68 — PublicationRequirements wire truth repair — READY / OPEN; R1 findings corrected; full reproof + R2 confirmation pending** |
| B10 Product candidate | **PR #64 — Preparação — PAUSED at P8 / NOT LOCKED pending prerequisite integration** |
| B10 finding | **UPSTREAM FINDING — accepted D4-R1/W2 publication-requirement meaning is richer than the current machine-readable OAD realization; smallest owner is the Product OAD wire realization** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Run fresh full aggregate proof and R2 independent confirmation on the corrected PR #68 candidate, then STOP for fresh explicit operator integration authorization. Do not resume PR #64/B10 P8, begin P9, Pre-D9/D9, or Product implementation before the prerequisite is integrated.** |
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
| D6 — Frontend | **ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110 LOCKED; B110 REVALIDATED / UNAFFECTED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED baseline; D7-R ACCEPTED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; D8-R ACCEPTED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B10 P8 paused on bounded PublicationRequirements wire prerequisite PR #68** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- PR #65 is integrated at `a2aeac19816c90ee30bf373cef0448d52a486c7e`. Draft PRs publish advisory `quick`; Ready candidates and `main` publish ruleset-required `required` and run the full aggregate gate.
- B10 P8 revalidation exposed a bounded downstream-falsification result: provider/context-specific publication requirements and source states are accepted semantics, but the current Product OAD realization cannot carry every material distinction the operated frontend must preserve.
- PR #68 corrects only that machine-readable wire realization: exact publication context, requirement class, bounded safe-authoring value constraints, and owner-local source evidence states `known / missing / conflicting / unknown / unavailable / unsupported`.
- R1 challenge found three valid proof/contract defects now corrected in the candidate: source candidate identity is structurally unique by opaque candidate key using a closed typed-key map that preserves the existing OAD closed-object invariant; each requirement value-spec family is mechanically coupled to the matching source-value family while `not_applicable` remains override-only; and semantic proof runs against the parsed Redocly bundle rather than raw YAML substrings.
- The repair does **not** add a Product/PIM master, provider field bag, rule DSL, new Product operation, Permission, Principal kind, runtime, ListingIntent override editor, or provider-specific business authority.
- `missing in source` remains distinct from `publication impossible`; `FOLLOW_SOURCE` / `EXPLICIT_OVERRIDE` remain Offering/ListingIntent resolution meaning, outside B10 editing.
- PR #64 remains paused at B10 P8. Existing operator LOCKs, Product **106/31/H-A-S**, runtime NONE, Pre-D9/D9 blocks and implementation block remain unchanged.

```text
PR #65 INTEGRATED
→ PR #68 PublicationRequirements wire prerequisite: fresh full proof + R2 confirmation + explicit integration authorization
→ rebase/revalidate PR #64 B10 P8
→ operator walkthrough / REVISE | UPSTREAM FINDING | explicit LOCK
→ P9 only after LOCK
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
