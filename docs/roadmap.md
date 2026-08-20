# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-program stage/status/allowed-work/next-action authority. Detailed Product and architecture semantics remain in their owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D5 — API — OPEN / ACTIVE** |
| Accepted through | D0–D4, D4-R1, D5-B1, D5-B2 W1/W2/W3/W4, Technical Ingress, final Problem/media consistency, and OpenAPI authority/tooling |
| Exact next action | **Author and prove the canonical Product OpenAPI Description** |
| Entry document to create | `contracts/api/product/openapi.yaml` |
| Product surface | **95 Product operations · 29 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Product Problem namespace | `https://conexus.fun/marketplace-central/problems/product/{slug}` |
| Technical Problem namespace | `https://conexus.fun/marketplace-central/problems/technical/{surface}/{slug}` when independently contracted |
| Ngrok | preview/tunnel only; never canonical `Problem.type` authority |
| Active runtime baseline | **NONE** |
| Implementation | **BLOCKED UNTIL D9** |

The admitted Product surface is contract scope, not an instruction to build 95 runtime handlers at once.

## Stage progression

| Stage | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED |
| D3 — Communication / Events | ACCEPTED / CLOSED |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | OPEN / ACTIVE |
| D6 — Frontend | BLOCKED BY D5 |
| D7 — Runtime / Jobs / Transactions | BLOCKED |
| D8 — Golden Flows | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Current D5 work boundary

The current sub-batch may:

- author `contracts/api/product/openapi.yaml` and its repository-local source closure;
- encode exactly the accepted 95 operations and 29-Permission projection;
- use the accepted Product Problem constants under `https://conexus.fun`;
- add the accepted OAD-specific lint/bundle/generation proof;
- generate/compile the accepted derived TypeScript and Go contract projections;
- preserve zero population of the retired legacy OpenAPI/manual SDK;
- stop and reopen only the smallest accepted authority if executable proof finds a material contradiction.

It may not:

- begin D6–D9 or Product implementation;
- add a 96th Product operation, Permission #30, or a fourth Principal kind;
- select runtime/router/database/deployment mechanics;
- restore the removed legacy runtime or a parallel OpenAPI/manual SDK authority;
- place provider/OAuth ingress or authored-media delivery into the Product OAD by convenience.

## Progression law

One coherent stage/gate lands before the next begins. A material contradiction reopens only the smallest implicated decision. Product implementation becomes eligible only after accepted D9.

For task-specific reading, return to [`index.md`](index.md).
