# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D8 — Golden Flows — OPEN / ACTIVE** |
| Accepted baseline | **D0–D7 ACCEPTED / CLOSED** |
| D8 authority | [D8 Golden Flows](engineering/rebaseline/D8-GOLDEN-FLOWS.md) — **OPEN / ACTIVE** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Derive and adjudicate the smallest D8 golden-flow set and falsifiable acceptance matrix. Do not begin D9 or Product implementation.** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage | Status |
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
| D8 — Golden Flows | **OPEN / ACTIVE** |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## D8 boundary

D8 composes accepted D0–D7 authority into a **small representative flow set**, not an exhaustive 99-operation test catalog and not Product implementation. Each selected flow must earn its place by falsifying a material cross-stage invariant: ownership, Organization isolation, auth/access, wire/front-end composition, durable effects, ambiguity/reconciliation, knowledge honesty or recovery.

Start from [D8 Golden Flows](engineering/rebaseline/D8-GOLDEN-FLOWS.md), then switch only to the exact accepted owner or canonical OAD needed for the candidate flow under review.

D9 remains blocked until D8 closes. Product implementation remains blocked until D9. Reopen D0–D7 only for a material falsifier at the smallest owning authority.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).
