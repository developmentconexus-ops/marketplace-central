# D6 — Frontend

> **Status:** OPEN / ACTIVE — D6-B1 Frontend Interaction & Authority Model
> **Program:** Architecture Rebaseline / Technical System Design
> **Opened:** 2026-08-21
> **Parent authorities:** `ARCHITECTURE.md`, accepted D0–D5 semantics, proved bounded D6-R1 amendment, and canonical Product OAD at `contracts/api/product/openapi.yaml`
> **Method:** DevelopmentConexus Engineering Method v1.0.0

## 1. Purpose and boundary

D6 defines the target Marketplace Central frontend interaction model and frontend topology without creating a second business authority.

React is the accepted client technology. TanStack Query remains the accepted server-state client unless material evidence explicitly reopens it. D6 consumes Product authority; it does not move business policy into React.

D6 does **not** choose server/router/database/transaction/worker/deployment mechanics, D8 golden-flow execution choreography or Product implementation. D7–D9 remain blocked and Product implementation remains blocked until accepted D9.

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

## 3. D6-B1 target invariant

> **Every material frontend observation or user-initiated Product action derives from an accepted semantic owner and canonical Product operation; the UI preserves Organization, identity, Permission, knowledge/freshness, attribution, concurrency, idempotency and consequential-outcome semantics without duplicating business truth.**

### 3.1 State ownership

| State class | Owner | Law |
| --- | --- | --- |
| Server state | Product/server authority via TanStack Query | never duplicate to a second global server-state store by convenience |
| URL/navigation state | browser/navigation contract | may select scope/filter/view, never business truth |
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

Without new evidence D6-B1 does not introduce:

- second global server-state store;
- generic workflow/action/command frontend layer;
- BFF/screen-shaped API;
- offline-first sync;
- websocket/event-stream architecture;
- microfrontends/plugins;
- universal design-system platform work;
- SSR/meta-framework selection by fashion;
- router/form/state library merely by preference;
- generic Analytics/Metric/Strategy client vocabulary;
- client-owned opportunity score/recommendation/AI explanation.

## 4. Current proved Product input

The current D6 candidate Product authority is mechanically proved at:

```text
99 Product operations
30 ordinary Permissions
Principal kinds H / A / S
28 List/Search operations
```

The proof also re-executes the accepted **95/29 baseline byte-for-byte** before proving the bounded D6-R1 extension. Performance adds exactly four Qs and `performance.read`; it adds zero C/P operations.

### 4.1 Bounded D6 findings already adjudicated

**Channel onboarding**
- currently connectable marketplace kind remains Mercado Livre only;
- available marketplace kind != connected Marketplace Installation != navigation context != Settings;
- no generic provider/integration catalog Product API;
- OAuth remains Technical Non-Product Ingress.

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

## 5. App Shell / information architecture — OPERATOR APPROVED, D6-R1 REFINED

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

Page modes:

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

The strategic experience may compose:

```text
Performance
+ Market Intelligence
+ Commercial Economics
+ Offering
+ Sales
+ Availability
```

under each owner's Permission. There is no `Strategy` domain, Product endpoint, store or write authority.

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
- empty navigation groups may be hidden when the current access context exposes no relevant read/action capability;
- no badge/count is inferred from one paginated page;
- `performance.read` grants Performance reads only and does not imply `market.read`, `economics.read`, Offering/Sales/Availability reads or any write Permission.

## 7. Responsive/context panels

- desktop: persistent primary nav + top Organization + page-local Installation/subnav;
- tablet: collapsible nav;
- mobile: navigation drawer + stacked content;
- responsive layout never removes semantic state/qualification;
- no permanent global right rail; page-owned drawer may show evidence/history/related Work while remaining ephemeral UI state.

## 8. D6-B1 revised interaction proof — CANDIDATE / INTERNAL CHALLENGE PASS

The revised [D6-B1 Frontend Interaction Map](D6-B1-INTERACTION-MAP.md) is written in Portuguese UX language and maps all **99/99** admitted Product operations to coherent interaction homes while preserving canonical operation names.

The revised low-fidelity prototype is [`qualification/d6-wireframes/index.html`](../../../qualification/d6-wireframes/index.html). It is also Portuguese and now exercises the approved shell plus:

- Preparação and ListingIntent authoring;
- Performance / Resumo;
- Performance / Publicações and Listing Performance detail;
- Performance / Mídia with campaign/listing/catalog/family scopes;
- Mercado and Economia as distinct strategic meanings;
- Disponibilidade;
- Venda cross-owner composition;
- Expedição physical checkpoints;
- Pós-venda;
- Trabalho;
- Configurações / Canais, Acesso and Aprovações.

Internal challenge against the interaction-map falsifiers found no additional Product/API prerequisite. Materially, the prototype makes visible or explicitly blocks:

- hidden Organization/Installation/SourceInstance authority;
- speculative current Amazon/Shopee connectability;
- generic Dashboard/Strategy/Analytics/Metric Product API;
- provider CVR/ROAS reconstruction in the client;
- campaign/family/catalog evidence collapsed into one Listing;
- partial/unknown evidence presented as complete/zero;
- historical provider evidence presented as MPC-authored truth;
- Ads management/optimization/AI controls;
- Market/Economics acquiring Price write;
- SetAvailableQuantity/Sync/Refresh shortcuts;
- Sale composition acquiring cross-owner write authority;
- physical checkpoint authority inferred from the button/client;
- Work/PostSale generic close;
- Governance decision acquiring target execution/Permission;
- blind retry of ambiguous effects.

This remains **interaction-level design proof**, not browser/runtime/provider/persistence proof and not D6-B1 ratification by itself.

## 9. Exact next D6 work

1. operator-review/adjudicate the revised [99-operation interaction map](D6-B1-INTERACTION-MAP.md) and [Portuguese low-fidelity wireframe proof](../../../qualification/d6-wireframes/index.html);
2. only if that interaction proof is approved, evaluate the smallest frontend feature/package topology and exact dependency needs from current official evidence rather than convention;
3. frontend topology research must not choose D7 server/runtime/router/database/deployment mechanics;
4. once one coherent interaction + topology candidate exists, submit it to the repository's independent milestone/final challenge path rather than using Fable/Claude as iterative co-authors;
5. do not begin D7–D9, implement Product code, ratify D6-B1 or merge PR #54 without explicit operator authorization.
