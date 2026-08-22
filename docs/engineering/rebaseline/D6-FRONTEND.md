# D6 — Frontend

> **Status:** ACCEPTED / CLOSED — OPERATOR-RATIFIED / FINAL INDEPENDENT CHALLENGE ADJUDICATED
> **Program:** Architecture Rebaseline / Technical System Design
> **Opened:** 2026-08-21
> **Closed:** 2026-08-21
> **Parent authorities:** `ARCHITECTURE.md`, accepted D0–D5 semantics, proved bounded D6-R1 amendment, proved bounded D5-R1 human-browser authentication correction, ratified D6-B1 interaction authority, and canonical Product OAD at `contracts/api/product/openapi.yaml`
> **Method:** DevelopmentConexus Engineering Method v1.0.0

## 1. Purpose and boundary

D6 defines the target Marketplace Central frontend interaction model and frontend realization topology without creating a second business authority.

React + TypeScript are the accepted browser-client technologies. TanStack Query owns server state. D6-B2 selects TanStack Router for navigation/URL state, `openapi-typescript` for generated Product wire shapes and `openapi-fetch` as the thin typed browser Product transport. D6 consumes Product authority; it does not move business policy into React or generated client code.

D6 does **not** choose the server HTTP router/mux, database/RLS, transaction, worker/queue, Keycloak realm/deployment, runtime/process/deployment mechanics, D8 golden-flow execution choreography or Product implementation. D7 is NEXT / NOT STARTED, D8–D9 remain blocked, and Product implementation remains blocked until accepted D9.

## 2. Imported invariants

1. **Frontend is a Product API client, never a business authority.**
2. **Canonical OpenAPI is the single Product wire authority.** No hand-written DTO/API shadow contract.
3. **Organization is explicit isolation scope.** No browser/provider/default tenant authority.
4. **Permission != disposition != provider capability != Governance authorization.** Route/button visibility is usability only.
5. **Knowledge stays honest.** Known zero/empty, unknown, unavailable, partial, unsupported and stale never collapse.
6. **Consequential outcomes stay honest.** Accepted, pending, rejected and ambiguous remain distinct.
7. **No blind replay.** Idempotency and concurrency semantics survive UX retries.
8. **External identities stay qualified.** Provider/native key never becomes global identity by client convenience.
9. **Technical Ingress stays separate.** OAuth/provider callback/Retail Media advertiser binding are not fake Product operations.
10. **Read composition never becomes write authority.**
11. **No screen-shaped API.** D6 does not reshape D5 merely to simplify a page.
12. **Historical Performance evidence remains source-qualified.** MPC custody does not convert provider-reported evidence into MPC-authored fact.
13. **Human browser tokens remain server-side.** D5-R1 requires same-origin application-session authentication plus CSRF request trust for unsafe human requests; browser JavaScript does not own the OIDC access/refresh token.

## 3. D6-B1 target invariant

> **Every material frontend observation or user-initiated Product action derives from an accepted semantic owner and canonical Product operation; the UI preserves Organization, identity, Permission, knowledge/freshness, attribution, concurrency, idempotency and consequential-outcome semantics without duplicating business truth.**

### 3.1 State ownership

| State class | Owner | Law |
| --- | --- | --- |
| Server state | Product/server authority via TanStack Query | never duplicate to a second global server-state store by convenience |
| URL/navigation state | TanStack Router / browser navigation contract | may select scope/filter/view, never business truth |
| Form draft | user editing session | unsent values never impersonate accepted server meaning |
| Ephemeral UI | frontend | panel/focus/temporary selection only |

A fifth durable/global state class requires a concrete consumer and proof these four are insufficient.

### 3.2 Interaction laws

- every screen/action maps to owner + Product operation/capability + Permission;
- current access context conditions navigation but never authorizes a server action;
- exact Installation is required whenever the Product operation requires it;
- Organization switch invalidates incompatible server/navigation context;
- consequential client retries preserve the same intake idempotency identity only for the same semantic request;
- concurrency/precondition failure stays distinct from business/provider rejection;
- frontend never manufactures source completeness, authoritative counts, provider metric definitions or causal claims;
- no generic `Refresh`, `Sync`, `SetPrice`, `SetAvailableQuantity`, `CloseWork`, `ClosePostSale`, `ConnectMarketplace`, Ads-write or AI-write shortcut exists unless Product authority later admits it.

### 3.3 YAGNI exclusions

Without new evidence D6 does not introduce:

- second global server-state store;
- generic workflow/action/command frontend layer;
- BFF/screen-shaped business API;
- offline-first sync;
- websocket/event-stream architecture;
- microfrontends/plugins;
- universal design-system platform work;
- SSR/meta-framework selection by fashion;
- form/schema/UI-component framework merely by preference;
- generated generic TanStack Query hooks/query-key policy from the OAD;
- generic Analytics/Metric/Strategy client vocabulary;
- client-owned opportunity score/recommendation/AI explanation.

## 4. Accepted Product input

The accepted D6 Product authority is mechanically proved at:

```text
99 Product operations
30 ordinary Permissions
Principal kinds H / A / S
28 List/Search operations
```

The proof re-executes the accepted **95/29 baseline**, the pre-auth **99/30 D6-R1 surface**, the current D5-R1 authentication profile and deterministic TypeScript/Go generated projections. Performance adds exactly four Qs and `performance.read`; it adds zero C/P operations.

Current authentication carriers are:

```text
H browser -> server-side OIDC login -> Secure HttpOnly application session + CSRF on unsafe methods
A / S     -> Client Credentials -> audience-bound bearer
```

Owning detail: [D5-R1 Human Browser Authentication Correction](D5-R1-HUMAN-BROWSER-AUTHENTICATION.md).

### 4.1 Bounded D6 findings already adjudicated

**Channel onboarding**
- currently connectable marketplace kind remains Mercado Livre only;
- available marketplace kind != connected Marketplace Installation != navigation context != Settings;
- no generic provider/integration catalog Product API;
- marketplace OAuth remains Technical Non-Product Ingress.

**D2-R1 Presentation Identity**
- current Principal, accessible Organizations, members and AccessRoles have human-readable `display_name` presentation metadata;
- canonical IDs/keys remain identity, scope, correlation and authorization authority.

**Readiness SourceInstance discovery**
- `SearchSourceProductsForMarketplace.source_instance_id` is an optional narrowing filter only;
- omission means bounded multi-source Readiness search, never hidden/default source selection;
- every result remains SourceInstance-qualified and point readiness/requirements reads require the selected exact source.

**D6-R1 Marketplace Performance Intelligence**
- the original frontend exposed a material strategic-analysis gap rather than a mere layout omission;
- bounded D0–D5 repair admits Marketplace Performance Intelligence as read/derive authority;
- D2-R2 preserves only source-qualified Performance evidence required for history/comparison/explanation when provider retention would otherwise erase it;
- Mercado Livre baseline evidence is Visits + current Product Ads performance; Ads management remains deferred;
- exact Retail Media advertiser binding remains Technical Non-Product configuration under current human + `portfolio.manage` + exact Installation;
- no generic Metric/Analytics API, no time-series baseline, no `signals[]`, no Ads mutation and no AI/MCP authority.

Owning detail: [D6-R1 Marketplace Performance Intelligence](D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md).

## 5. App Shell / information architecture — OPERATOR-RATIFIED

Stable shell model:

```text
Marketplace Central
Organization: <display_name>

VISÃO GERAL

OPERAÇÕES
  Preparação
  Publicações
  Disponibilidade
  Vendas
  Expedição
  Pós-venda

ESTRATÉGIA E INTELIGÊNCIA
  Performance
    Resumo
    Publicações
    Mídia
  Mercado
  Economia

CONTROLE
  Trabalho
  Aprovações

CONFIGURAÇÕES
```

Primary navigation is a user mental model, not D1 package taxonomy.

### 5.1 Organization context

- `organization_id` is always canonical scope;
- `display_name` is presentation only;
- switch Organization re-resolves server state and clears incompatible Installation context;
- a remembered Organization preference is usable only after current access context confirms it.

### 5.2 Marketplace Installation context

Installation is contextual navigation state, never ambient authority.

| Mode | Meaning |
| --- | --- |
| organization-wide | Product operation has no Installation dimension |
| all-or-exact | admitted read supports optional Installation narrowing |
| exact-required | Product contract requires one exact Installation; no synthetic merge |

Source-qualified Listing/Sale/Shipment collections and all Performance reads require exact Installation where their Product paths require it.

### 5.3 Performance workspace

Performance answers **how our participation is performing**, while Market answers **what the competitive market is doing** and Economics answers **what the economics mean**.

Binding UX laws:

- Performance always selects one exact Marketplace Installation;
- `Resumo` consumes `GetMarketplacePerformanceSummary` plus separate owner reads only when the UI needs them;
- `Publicações` consumes `ListMarketplaceListingPerformance`; every known Listing remains visible even when Performance evidence is unknown/unavailable;
- Listing detail gains a **Performance** context using `GetMarketplaceListingPerformance` without moving Listing ownership from Offering;
- `Mídia` consumes `ListRetailMediaPerformance` and preserves campaign/listing/catalog-group/family-group scope rather than pretending every metric belongs to one Listing;
- provider-reported percentages/multiples retain their measurement basis; frontend never reconstructs provider CVR/ROAS;
- current/comparison periods are explicit Product requests; presets such as “últimos 30 dias” remain navigation convenience;
- comparison may display a numeric change only when Product Performance declares the evidence comparable;
- historical values may come from preserved source evidence and must display coverage/provenance honestly;
- no global “Todos os marketplaces” KPI aggregate is baseline. Future multi-provider views present per-Installation results unless equivalence is explicitly proven.

### 5.4 Strategy Workspace is composition only

The strategic experience may compose Performance + Market Intelligence + Commercial Economics + Offering + Sales + Availability under each owner's Permission. There is no `Strategy` domain, Product endpoint, store or write authority.

### 5.5 Overview

Overview remains a small read-only composition. It may show installation posture, bounded economic summary, bounded Work preview and a per-Installation Performance entry point when respective Permissions are present. It never creates `/dashboard`, hidden global counts or cross-provider metric aggregation.

### 5.6 Settings

Settings remains low-frequency UX grouping only:

- Canais / Marketplace Installations → Portfolio + Technical Ingress ceremony;
- Entidades vendedoras → Portfolio;
- Acesso → IdentityAccess;
- Fontes de estoque / política de disponibilidade → Availability;
- Nós/metas de expedição → Fulfillment;
- Política comercial → Economics;
- Delegações → Governance.

Settings never becomes semantic owner.

## 6. Permission-conditioned visibility

- menu/action hiding reduces noise only;
- server authorization remains authoritative on deep links and stale access context;
- empty navigation groups may be hidden when current access context exposes no relevant capability;
- no badge/count is inferred from one paginated page;
- `performance.read` grants Performance reads only and does not imply `market.read`, `economics.read`, Offering/Sales/Availability reads or any write Permission.

## 7. Responsive/context panels

- desktop: persistent primary nav + top Organization + page-local Installation/subnav;
- tablet: collapsible nav;
- mobile: navigation drawer + stacked content;
- responsive layout never removes semantic state/qualification;
- no permanent global right rail; page-owned drawer may show evidence/history/related Work while remaining ephemeral UI state.

## 8. D6-B1 revised interaction proof — OPERATOR-RATIFIED

The revised [D6-B1 Frontend Interaction Map](D6-B1-INTERACTION-MAP.md) is written in Portuguese UX language and maps all **99/99** admitted Product operations to coherent interaction homes while preserving canonical operation names.

The revised low-fidelity prototype is [`qualification/d6-wireframes/index.html`](../../../qualification/d6-wireframes/index.html). It exercises the approved shell plus Preparação, ListingIntent, Performance, Mercado, Economia, Disponibilidade, Venda, Expedição, Pós-venda, Trabalho, Canais, Acesso and Aprovações.

The proof visibly preserves or blocks hidden Organization/Installation/SourceInstance authority, fake provider connectability, generic Dashboard/Strategy/Analytics APIs, provider metric reconstruction, scope collapse, partial/unknown evidence as complete/zero, historical provider evidence as MPC-authored truth, Ads/AI controls, Market/Economics price mutation, Availability sync shortcuts, cross-owner write leakage, client-inferred physical authority, generic Work/PostSale close, Governance-as-execution and blind retry.

After independent Fable review, bounded corrections and final executable proof, the operator ratified D6-B1 on 2026-08-21. This remains interaction-level frontend authority, not D7 runtime/persistence proof.

## 9. D6-B2 frontend realization — OPERATOR-RATIFIED

### 9.1 Target invariant

> **Frontend packages follow stable human lenses/flows for UI composition while reusable Product consumption follows owner/operation-family adapters; navigation owns URL state, TanStack Query owns server state, generated OpenAPI shapes own the wire, and no layer acquires business authority by convenience.**

This bounded hybrid was selected after the Marketplace Central ↔ MetalDocs reciprocal review. Strict owner-first UI mirroring and strict lens-only API duplication are both rejected.

### 9.2 Selected dependency profile

| Concern | D6-B2 decision | Why now |
| --- | --- | --- |
| UI runtime | React + TypeScript strict | accepted client technology; typed composition without new authority |
| Server state | TanStack Query | already accepted; cache/query lifecycle stays one server-state authority |
| Navigation/URL state | `@tanstack/react-router` | typed routes plus runtime-validated/inherited search params protect Organization/Installation/period/filter deep links |
| Wire type generation | `openapi-typescript` | generated `paths/components` consume the canonical OAD without hand-authored DTOs |
| Low-level Product HTTP client | `openapi-fetch` | thin OAD-bound Fetch client removes manual path/param/body drift while preserving headers, middleware, custom serializers and raw/stream response handling |
| Browser network primitive | native Fetch via `openapi-fetch` | no Axios/general HTTP abstraction needed |

Current official evidence checked at ratification:

- TanStack Router search params: <https://tanstack.com/router/latest/docs/guide/search-params>
- TanStack Router type safety: <https://tanstack.com/router/latest/docs/guide/type-safety>
- `openapi-fetch`: <https://openapi-ts.dev/openapi-fetch/>
- `openapi-fetch` API/serializers: <https://openapi-ts.dev/openapi-fetch/api>
- `openapi-fetch` middleware: <https://openapi-ts.dev/openapi-fetch/middleware-auth>

Exact dependency versions remain implementation-manifest concerns unless a version-specific property later becomes architectural.

Not selected now:

- Orval / Hey API / generated generic query SDK;
- `openapi-react-query` as architectural layer;
- Axios;
- Redux/Zustand as server-state mirror;
- React Hook Form or another form framework;
- Zod/Valibot/ArkType as a universal Product schema authority;
- MUI/Radix/shadcn/Storybook or universal design-system platform;
- SSR/meta-framework.

A local validation/component/form dependency may be admitted later only for a concrete repeated UI failure class. Product schemas remain OAD-owned regardless.

### 9.3 Target source topology

Exact folder spelling may vary mechanically, but these architectural classes and allowed responsibilities are binding:

```text
src/
  app/
    providers/       # Query/Router/browser composition only
    router/          # route tree wiring/navigation policy
    shell/           # app shell/workspace chrome

  routes/            # thin TanStack Router route definitions

  features/
    <lens-or-flow>/  # user-facing composition + form draft + ephemeral UI

  api/
    generated/       # openapi-typescript output; never hand edited
    transport/       # one thin openapi-fetch Product transport policy
    <owner-family>/  # stateless Product operation/query/command adapters

  ui/                # bounded reusable presentation primitives only
```

Representative **feature/lens** homes follow the ratified UX rather than D1 package names:

```text
overview
preparation
publications
availability
sales
fulfillment
post-sale
performance
market
economics
work
approvals
settings
```

Representative `api/<owner-family>` homes may reflect Product operation ownership/families such as Access, Portfolio, Readiness, Offering, Availability, Market, Performance, Economics, Governance, Sales, Materialization, Fulfillment, Post-Sale and Work. These are **stateless API-consumption modules**, not client-side domain models.

### 9.4 Dependency direction

Allowed direction is closed-world:

```text
app/routes
    ↓
features
    ↓
api/<owner-family>
    ↓
api/transport
    ↓
api/generated

features ──→ ui
app/routes ─→ ui
```

Binding prohibitions:

- one feature does not import another feature's private internals;
- one `api/<owner-family>` does not import another owner family's internals;
- `routes` do not fetch Product data directly around TanStack Query;
- `ui` imports neither features nor Product API modules;
- `generated` imports no app/feature code and is never hand edited;
- `transport` contains HTTP/session/CSRF/Problem/header mechanics only, never business decisions;
- there is no `features/strategy`, `features/dashboard`, `api/analytics`, generic action/workflow bus or catch-all business `shared/` package.

Cross-owner screen composition belongs in the user lens/feature or route composition boundary. It never moves operations or semantics into a new owner.

When implementation becomes authorized, first-party import edges must be mechanically checked with a default-deny rule. The exact lint/plugin mechanism is intentionally not frozen before a real source tree exists.

### 9.5 Router / URL contract

D6-B1 route meanings remain canonical navigation identity:

```text
/org/:organizationId/visao-geral
/org/:organizationId/preparacao
/org/:organizationId/publicacoes/*
/org/:organizationId/disponibilidade
/org/:organizationId/performance/*
/org/:organizationId/mercado
/org/:organizationId/economia/*
/org/:organizationId/vendas/*
/org/:organizationId/expedicao/*
/org/:organizationId/pos-venda/*
/org/:organizationId/trabalho/*
/org/:organizationId/aprovacoes/*
/org/:organizationId/configuracoes/*
```

TanStack Router realizes, but does not redefine, those meanings.

Laws:

- `organizationId` is explicit path navigation state and is revalidated against current access context;
- Marketplace Installation is typed/validated URL state on all-or-exact/exact-required lenses; exact-required routes render an explicit chooser/blocked state when missing or invalid and never silently choose a first/default Installation;
- Performance current/comparison periods are shareable typed URL state and translate to the explicit Product period parameters;
- filters/view selection may live in URL when sharing/back-forward behavior is a real user property;
- opaque cursor state is not promoted to durable business identity;
- malformed URL input is validated/fails to a safe navigation state rather than becoming Product truth;
- router loaders do not become a second server-state cache. If route-level prefetch is later useful, it may only prewarm the same TanStack Query contract (`ensureQueryData`/equivalent) used by the feature.

No router loader/action becomes Product command/business authority.

### 9.6 Product client / transport contract

Pipeline:

```text
canonical Product OAD
  -> openapi-typescript generated paths/components
  -> one openapi-fetch client
  -> stateless owner/operation adapters
  -> TanStack Query query/mutation contracts
  -> feature/lens composition
```

The transport/adapters must preserve:

- exact Product paths and path/query/body serialization;
- `Idempotency-Key`, `If-Match`, ETag and other admitted headers without generic mutation abstractions;
- Problem Details as typed Product errors rather than thrown provider/runtime strings;
- multipart/custom serializer needs such as ListingIntent media;
- blob/stream/arrayBuffer/raw-response needs where the Product OAD admits them;
- no hidden retry;
- no provider DTO or Technical Ingress operation leaking into the Product client.

For the human React client:

- use the D5-R1 same-origin application-session carrier;
- ordinary browser code does **not** inject/store an OIDC bearer or refresh token;
- CSRF is a request-trust mechanism for unsafe Product calls, not user identity or Permission;
- exact CSRF bootstrap/storage/session realization remains D7.

The OAD still admits `MpcMachineBearerAuth` for non-browser A/S clients; D6 does not turn the human SPA into that client class.

### 9.7 TanStack Query contract

Query identity derives from:

```text
stable Product operation identity
+ exact canonical semantic inputs
```

Inputs include `organization_id`, exact Marketplace Installation/source-qualified identity, Product filters and Performance periods whenever the operation semantics include them. Query-key spelling is implementation detail; omitting a semantic scope input is non-conformant.

Further laws:

- no normalized universal entity store alongside TanStack Query;
- generated OAD types are the wire representation; owner adapters may expose bounded query/command functions but not shadow DTO/domain models;
- presentation-only derived values/view models remain feature-local and never become durable server truth;
- mutation success invalidates/refetches only reads whose accepted semantics may have changed; no generic client event bus is introduced;
- no global retry policy hides Product knowledge/outcome states;
- consequential mutations are not automatically replayed; ambiguous potentially accepted effects follow accepted idempotency/reconciliation semantics;
- freshness/staleness is operation/source-aware; the client does not invent provider freshness by a universal `staleTime` constant;
- Organization switch invalidates/re-scopes incompatible queries and navigation state.

### 9.8 Reopen triggers

Reopen only the smallest D6-B2 decision when a real proof establishes one of these:

- TanStack Router creates greater measured type/build/route complexity than it removes, or cannot preserve the accepted URL-state contract without competing with TanStack Query;
- `openapi-fetch` cannot represent an admitted Product path/param/header/multipart/stream behavior without unsafe casts or a larger custom bypass than a thin native-Fetch wrapper;
- the lens + stateless owner-adapter split creates unavoidable circular dependencies or duplicate semantic logic;
- the four accepted frontend state classes cannot represent a concrete real interaction;
- a concrete repeated form/accessibility/component failure class proves a missing UI dependency.

Preference or technology fashion is not a reopen trigger.

## 10. D6 closeout — ACCEPTED / CLOSED

1. D6-B1 interaction authority is operator-ratified.
2. D6-B2 topology/dependency authority is operator-ratified.
3. The isolated final Fable challenge returned **ACCEPT WITH BOUNDED FIXES** and found no Critical finding or material reason to reopen D5-R1, D6-R1, D6-B1, D6-B2, the 99/30 Product surface, the authentication split or the selected frontend profile.
4. GPT independently adjudicated every material finding. The accepted fixes repaired the default gate entry point, D6-B1 status/count precision and fixed-metric Traffic/Sales comparison-state precision; Retail Media remains intentionally sparse by scope/basis rather than acquiring false metric uniformity. Repository-only reporting observations were not promoted into architecture work.
5. Fresh executable proof after those bounded fixes confirms the accepted 95/29 D5 baseline, current 99/30 Product surface, H/A/S principal kinds, 28/28 List/Search operations, 7/7 Performance controls, 5/5 authentication controls, deterministic/compilable TypeScript and Go projections, 2/2 Performance-knowledge controls, zero legacy runtime population and the repository negative control.
6. No material contradiction survives adjudication; a second Fable round is not warranted.
7. The operator explicitly authorized D6 closeout and PR #54 merge on 2026-08-21.
8. D7 — Runtime / Jobs / Transactions is **NEXT / NOT STARTED**. Do not begin D8–D9 or implement Product code before accepted D9.
