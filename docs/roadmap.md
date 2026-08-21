# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-program stage/status/allowed-work/next-action authority. Detailed semantics remain in their owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6 — Frontend — OPEN / ACTIVE** |
| Accepted baseline | **D0–D5 ACCEPTED / CLOSED** |
| D6-discovered bounded amendment | **D6-R1 Marketplace Performance Intelligence — OPERATOR-APPROVED / EXECUTABLE OAD PROOF PASS** |
| Exact next action | **D6-B1 — re-derive the interaction map and low-fidelity wireframes in Portuguese against the proved 99-operation / 30-Permission Product surface, including Strategy & Intelligence → Performance; then operator-review before frontend topology/dependency adjudication** |
| Canonical Product OAD candidate | `contracts/api/product/openapi.yaml` |
| Product surface in current D6 candidate | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Implementation | **BLOCKED UNTIL D9** |

The Product operation count is contract scope, not a handler/screen count.

## Stage progression

| Stage | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; bounded D6-R1 amendment approved in current candidate |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; bounded D6-R1 amendment approved in current candidate |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; D2-R1 + D2-R2 bounded amendments approved in current candidate |
| D3 — Communication / Events | ACCEPTED / CLOSED; bounded D6-R1 Q/P amendment approved in current candidate |
| D4 — External Integrations | ACCEPTED / CLOSED; bounded Mercado Livre Performance evidence amendment approved in current candidate |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; bounded D6-R1 99/30 wire amendment proved in current candidate |
| D6 — Frontend | **OPEN / ACTIVE** |
| D7 — Runtime / Jobs / Transactions | BLOCKED |
| D8 — Golden Flows | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## D5 baseline and D6-R1 amendment

`main @ 9d2c81e175bc39ac388c9d8924ddad21f2a86480` remains the accepted D5 closeout baseline at 95 operations / 29 ordinary Permissions. D6 exposed a material strategic-analysis gap and the operator approved the smallest bounded parent repair recorded in [D6-R1 Marketplace Performance Intelligence](engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md).

The current D6 candidate adds only:

- one read/derive boundary, **Marketplace Performance Intelligence**;
- one ordinary Permission, `performance.read`;
- four H/A/S Product Qs: marketplace summary, Listing-performance list, Listing-performance point read and Retail Media performance list;
- D2-R2 historical Performance-evidence custody so provider retention does not become MPC historical retention;
- bounded Mercado Livre Visits + current Product Ads evidence semantics and a Technical Non-Product advertiser-binding ceremony under `portfolio.manage`;
- no Performance mutation, Ads management, generic Analytics/Metric API, Product P operation, AI/MCP authority or D7 mechanism.

Executable proof on the repaired OAD establishes:

```text
accepted D5 baseline non-regression  PASS  95/95 ops · 29/29 Permissions · 26/26 List/Search · 12/12 negative controls
current D6-R1 Product candidate      PASS  99/99 ops · 30/30 Permissions · 28/28 List/Search · 7/7 Performance negative controls
Principal kinds                      H / A / S
TypeScript + Go projections          deterministic / compilable
legacy runtime population            0
runtime schema enforcement           NOT_CLAIMED_D7
router selection                     NONE_D7
```

The amendment is bounded; it is not permission to generally reopen D0–D5.

## D6-B1 current authority

[D6 Frontend](engineering/rebaseline/D6-FRONTEND.md) remains the active frontend authority. Binding shell decisions already approved:

- Organization is the global workspace/isolation context;
- primary navigation is task-oriented, not D1 package taxonomy;
- Marketplace Installation is explicit contextual navigation state and never ambient authority;
- Settings is low-frequency grouping, not a business owner;
- read-only screen composition never becomes write authority;
- route/button visibility is usability, never authorization;
- server state, URL/navigation state, form draft and ephemeral UI state remain distinct;
- unknown/partial/unavailable/stale and accepted/pending/rejected/ambiguous remain honest;
- Product/Technical Ingress separation remains binding.

D6-R1 requires the Portuguese strategic group:

```text
ESTRATÉGIA E INTELIGÊNCIA
  Performance
    Resumo
    Publicações
    Mídia
  Mercado
  Economia
```

Performance requires an exact Marketplace Installation. Future multi-marketplace views may compose per-Installation results but may not add provider metrics as if their measurement bases were automatically equivalent.

## D6 blocked/deferred boundary

While D6 is active:

- do not begin D7–D9;
- do not select server runtime/router/database/transaction/worker/deployment mechanics;
- do not implement Product code;
- do not restore retired runtime/frontend/manual SDK authority;
- do not add Ads management, autonomous commercial optimization or AI/MCP authority;
- reopen an accepted earlier decision only for a material falsifier and only at the smallest responsible authority.

## Progression law

One coherent stage/gate lands before the next begins. Product implementation becomes eligible only after accepted D9. PR #54 remains Draft and is not authorized for merge without explicit operator authorization.

For task-specific reading, return to [`index.md`](index.md).