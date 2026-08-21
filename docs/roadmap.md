# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6 — Frontend — OPEN / ACTIVE** |
| Accepted baseline | **D0–D5 ACCEPTED / CLOSED** |
| D6-R1 | **Marketplace Performance Intelligence — OPERATOR-APPROVED / FABLE REVIEWED / BOUNDED FIXES ADJUDICATED** |
| D6-B1 | **OPERATOR-RATIFIED — corrected Portuguese interaction map + wireframes** |
| Cross-repo review | **Marketplace Central ↔ MetalDocs adjudicated by operator; bounded D5 human-browser auth correction approved** |
| D5-R1 auth correction | **OPERATOR-APPROVED / EXECUTABLE PROOF IN PROGRESS** — [Human Browser Authentication Correction](engineering/rebaseline/D5-R1-HUMAN-BROWSER-AUTHENTICATION.md) |
| Exact next action | **Complete the bounded D5-R1 RED→GREEN auth/OAD proof; if green, resume D6-B2 frontend topology/dependency adjudication** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
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
| D5 — API | ACCEPTED / CLOSED — bounded D5-R1 human-browser auth carrier correction active in current D6 candidate |
| D6 — Frontend | OPEN / ACTIVE — B1 ratified; B2 paused only for bounded D5-R1 proof |
| D7 — Runtime / Jobs / Transactions | BLOCKED |
| D8 — Golden Flows | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Current accepted Product authority

`main @ 9d2c81e175bc39ac388c9d8924ddad21f2a86480` remains the accepted D5 closeout baseline at 95 operations / 29 Permissions. The bounded [D6-R1 Marketplace Performance Intelligence](engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md) repair produces the proved 99/30 Product candidate with four Performance Qs and `performance.read`.

The operator-approved [D5-R1 Human Browser Authentication Correction](engineering/rebaseline/D5-R1-HUMAN-BROWSER-AUTHENTICATION.md) changes **authentication carriers only**:

```text
H browser  -> server-side OIDC login -> Secure HttpOnly application session + CSRF on unsafe requests
A / S      -> Client Credentials -> audience-bound bearer
```

It does not add/remove operations, Permissions, Principal kinds, owners or business semantics. Keycloak remains the preferred first OIDC provider candidate; exact session/CSRF/Keycloak/runtime realization remains D7.

Executable proof must preserve:

```text
accepted D5 baseline             95/95 operations · 29/29 Permissions
pre-auth D6-R1 surface           99/99 operations · 30/30 Permissions
current Product                  99 Product operations · 30 ordinary Permissions · H/A/S only
Performance                      exact four Qs · performance.read
legacy runtime population        0
```

## D6-B1 ratified interaction authority

[D6 Frontend](engineering/rebaseline/D6-FRONTEND.md) owns the detailed frontend laws. The ratified interaction model preserves explicit Organization/Installation context, task-oriented Portuguese navigation, Product/Technical-Ingress separation, read composition without write authority and honest knowledge/outcome states.

Strategic IA remains:

```text
ESTRATÉGIA E INTELIGÊNCIA
  Performance — Resumo / Publicações / Mídia
  Mercado
  Economia
```

The [interaction map](engineering/rebaseline/D6-B1-INTERACTION-MAP.md) covers all 99 operations and the Portuguese [wireframes](../qualification/d6-wireframes/index.html) are operator-ratified.

## D6-B2 bounded continuation

After D5-R1 is green, D6-B2 resumes from the operator-approved cross-repository adjudication:

- React + TypeScript;
- TanStack Query remains server-state authority;
- TanStack Router is the current common router candidate;
- `openapi-typescript` remains generated wire-shape authority and `openapi-fetch` is the current Marketplace low-level transport candidate;
- UI packages follow stable human lenses/flows while stateless Product API adapters may group reusable owner/operation-family consumption;
- URL/navigation, form draft and ephemeral UI remain separate from server state;
- no second global server-state store, generic workflow/action/query layer, microfrontends, SSR-by-fashion, offline-first, realtime or universal design-system platform without a proven consumer.

D6-B2 does **not** select D7 server runtime, HTTP mux, database/RLS, transaction, River, deployment or Keycloak realm topology.

## Blocked / deferred boundary

D7 inherits pre-vetted candidates from the cross-repo review but none are selected by D6: modular-monolith class, pgx/pgxpool, structural tenant isolation/RLS proof, River-first durable-work falsification, OpenTelemetry/OTLP/slog, sqlc, tern and real-dependency test tooling.

Do not begin D7–D9, restore retired runtime/manual SDK authority, add Ads management/autonomous optimization/AI authority, or merge PR #54 without explicit operator authorization. Reopen accepted authority only for a material falsifier at the smallest owning stage.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).
