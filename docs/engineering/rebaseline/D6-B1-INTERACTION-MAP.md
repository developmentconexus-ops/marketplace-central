# D6-B1 — Frontend Interaction Map

> **Status:** DERIVED CANDIDATE — App Shell / IA operator-approved; user-flow and screen inventory derived for D6-B1 proof, not yet a D6-B1 ratification
> **Parent:** `D6-FRONTEND.md`
> **Wire authority:** `contracts/api/product/openapi.yaml`
> **Scope:** Product 1.0 frontend only; no D7 runtime mechanics and no Product implementation
> **Derived:** 2026-08-21

## 1. Purpose

This artifact is the D6-B1 proof map from user-visible frontend interaction to accepted semantic authority.

It answers:

```text
user need
  → navigation/screen state
  → semantic owner
  → canonical Product operation/capability
  → ordinary Permission
  → safety/state treatment
```

It is deliberately **not** a component library, design system, router selection, frontend package tree, BFF contract or implementation plan.

The admitted 95-operation Product surface does not imply 95 screens. Multiple owner operations may be presented inside one coherent owner workspace, and a read-only screen may compose several owners without becoming write authority.

---

## 2. Approved global information architecture

The operator approved the following D6 App Shell / IA direction on 2026-08-21:

```text
Organization-global
+ task-oriented primary navigation
+ explicit contextual Marketplace Installation selection where semantically valid
+ low-frequency Settings separated from routine operation
+ read-only Overview composition
```

### 2.1 Primary navigation

```text
Overview

OPERATIONS
  Readiness
  Listings
  Availability
  Sales
  Fulfillment
  Post-sale

INTELLIGENCE
  Market
  Economics

CONTROL
  Work
  Approvals

Settings
```

Primary navigation is a usability model, not the D1 bounded-context taxonomy.

### 2.2 Organization context

Organization is the one global business scope in the shell.

- canonical scope is always `organization_id`;
- `display_name` is presentation metadata only under D2-R1;
- switching Organization clears incompatible navigation/channel context and causes Organization-scoped server reads to be re-resolved;
- no provider identity, Marketplace Installation, browser default or display label substitutes for Organization scope.

### 2.3 Channel context

Marketplace Installation is **contextual navigation state**, never ambient business authority.

Three page modes are allowed:

| Mode | Meaning |
| --- | --- |
| `organization-wide` | operation has no marketplace-installation dimension in its public contract |
| `all-or-exact` | admitted read supports an optional Installation filter; `All accounts` is honest |
| `exact-required` | Product contract requires an exact Marketplace Installation; no synthetic cross-Installation merge is presented |

`All accounts` is forbidden for source-qualified collections whose Product operation exists only under one Installation and whose independent cursors/coverage cannot be honestly merged. Product 1.0 therefore requires exact Installation context for Marketplace Listing observations, Marketplace Sales and Shipments.

### 2.4 Available channel versus connected account

```text
available marketplace kind
!= Marketplace Installation
!= channel navigation context
!= low-frequency configuration
```

Product 1.0 admits **Mercado Livre only** as currently connectable. The Add-channel interaction architecture is stable for future providers, but D6 does not display Amazon/Shopee as connectable until responsible D4/D5 authority admits them.

No generic Integration/Channel Catalog Product operation is introduced.

### 2.5 Overview

Overview is a read-only composition and owns no business meaning.

Baseline eligible panels are deliberately small:

- Marketplace Installation posture via `ListMarketplaceInstallations` when `portfolio.read` is present;
- economic period summary via `GetEconomicPerformanceSummary` when `economics.read` is present;
- open Work preview via `ListWork` when `work.read` is present.

A collection page is never converted into an authoritative global count merely because the frontend can count the current page. No new `/dashboard` Product operation is required.

### 2.6 Responsive structure

- wide desktop: persistent primary sidebar + top Organization context + page-local context/subnavigation;
- constrained desktop/tablet: collapsible sidebar/rail;
- narrow viewport: navigation drawer and stacked content;
- responsive changes presentation only; Permission, Organization, source identity and business state semantics do not change by viewport.

---

## 3. Frontend navigation identity grammar

These are frontend navigation identities, not Product API paths and not a router-library commitment.

```text
/org/:organizationId/overview
/org/:organizationId/readiness
/org/:organizationId/listings/*
/org/:organizationId/availability
/org/:organizationId/sales/*
/org/:organizationId/fulfillment/*
/org/:organizationId/post-sale/*
/org/:organizationId/market
/org/:organizationId/economics/*
/org/:organizationId/work/*
/org/:organizationId/approvals/*
/org/:organizationId/settings/*
```

MPC-owned resources may use their own opaque ID in frontend detail routes. Externally authoritative resources must keep source qualification in navigation identity, for example an Installation plus provider-native Listing/Sale/Shipment key.

A navigation query parameter or selected row is never sufficient business correlation unless the Product request also carries the canonical required identity/scope.

---

## 4. User flows

### F1 — Add and authorize a marketplace account

```text
Settings / Channels
  → ListMarketplaceInstallations
  → Add channel (currently Mercado Livre only)
  → CreateMarketplaceInstallation [Idempotency-Key]
  → Installation account_posture may require authorization
  → D5 Technical Ingress OAuth begin for exact existing Installation
  → provider ceremony/callback
  → GetMarketplaceInstallation
  → account_posture = bound | authorization_required | unavailable
```

Owner split:

- Marketplace Portfolio owns Installation business participation/configuration;
- D5 Technical Ingress owns authorization ceremony mechanics;
- frontend owns neither credentials nor provider account identity.

There is no Product `ConnectMarketplace` operation.

### F2 — Prepare a source Product and publish a ListingIntent

```text
Readiness / exact Marketplace Installation
  → SearchSourceProductsForMarketplace
       source_instance_id optional narrowing filter
       every result remains SourceInstance + native_product_key
  → select exact source-qualified Product
  → GetProductChannelReadiness
  → GetPublicationRequirements
  → Resolve/Clear correspondence when required
  → CreateListingIntentDraft [Idempotency-Key]
  → Get/Update ListingIntent [ETag / If-Match]
  → CreateListingIntentMedia when required [Idempotency-Key + typed parent etag]
  → CreatePriceIntent separately when a desired price is required
  → GetSellableAvailability for pre-creation target
  → optionally EvaluatePriceScenario
  → SubmitListingIntent with current typed etag
  → observe owner external-effect/convergence state; never infer convergence from transport success
```

The source-search omission law is explicit: no `source_instance_id` means bounded multi-source search across current Organization-scoped Readiness-admitted sources; it never means “pick the first/default source.”

### F3 — Analyze market/economics and request a price change

```text
Listing/source-product subject
  → Get/ListCompetitivePosition
  → ListComparableOffers
  → Get/ListExpectedEconomics
  → EvaluatePriceScenario (stateless)
  → CreatePriceIntent [Idempotency-Key]
  → GetPriceIntent / convergence observation
```

Market evidence and Economics never acquire price-write authority. `CreatePriceIntent` remains Offering-owned under `price.manage`.

### F4 — Observe and configure Availability

Routine operation:

```text
Availability
  → ListSellableAvailability [All accounts or exact Installation]
  → GetSellableAvailability
  → observe desired/provider quantity, knowledge/provenance and convergence
```

Low-frequency configuration:

```text
Settings / Inventory Sources
  → List/Get/Create/Update/DeactivateInventorySource

Settings / Availability Policy
  → GetEffectiveAvailabilityAllocationScopePolicy
  → UpdateAvailabilityAllocationScopePolicy
```

There is no frontend `SetAvailableQuantity`, `SyncAvailability` or manual AvailabilityIntent screen.

### F5 — Operate a marketplace Sale through cross-owner lifecycle

```text
Sales / exact Marketplace Installation
  → ListMarketplaceSales
  → GetMarketplaceSale
  → ResolveSaleSellingEntityAttribution when needed

Sale detail read-only composition
  → GetSaleEconomics
  → List/GetBusinessOrderIntent
  → Get/ResolveBusinessSystemPartyResolution when needed
  → GetDestinationRealization
  → List/GetInvoicingIntent
  → List/GetFulfillmentExecution
  → List/GetPostSaleResolution
```

The Sale detail is a **P-like client composition**, not a Sale-owned workflow aggregate. Every write remains on the owning operation.

There is no client `CreateBusinessOrderIntent`, `CreateInvoicingIntent`, direct Sankhya command or generic retry operation.

### F6 — Execute physical fulfillment and observe Shipment

```text
Fulfillment / Executions
  → ListFulfillmentExecutions
  → GetFulfillmentExecution
  → RecordSeparation [Idempotency-Key + current etag]
  → RecordPhysicalConference [Idempotency-Key + current etag]
  → RecordPacking [Idempotency-Key + current etag]
  → RecordDispatchHandoff [Idempotency-Key + current etag]
  → List/GetFulfillmentArtifacts when caller has fulfillment.execute

Fulfillment / Shipments / exact Marketplace Installation
  → ListShipments
  → GetShipment
```

Physical checkpoint buttons appear only for usability when the current access context permits them; backend Principal-kind and physical qualification enforcement remains authoritative.

### F7 — Coordinate post-sale consequences

```text
Post-sale
  → ListPostSaleResolutions [optionally narrowed to Sale/Installation]
  → GetPostSaleResolution
  → CreatePostSaleResolution [Idempotency-Key]
```

There is no direct frontend close, refund, cancel-sale, provider-claim action or generic resolution-status update. Closure comes from sufficient owner evidence.

### F8 — Handle Operational Work

```text
Work
  → ListWork [lifecycle/responsibility/assignment/origin filters]
  → GetWork
  → AssignWork
  → ClearWorkAssignment
  → HoldWork
  → ResumeWork
  → EscalateWork
```

There is no `CreateWork`, `CloseWork`, `DismissWork` or generic `SubmitWorkResolution` baseline. Source-owner resolution controls the originating condition; Work does not become a command bus.

### F9 — Make and administer authorization decisions

Contextual decision path:

```text
action-owner target requiring authorization
  → exact target/revision context
  → CreateAuthorizationDecision [Idempotency-Key]
  → GetAuthorizationDecision
```

Control/history:

```text
Approvals
  → ListAuthorizationDecisions
  → GetAuthorizationDecision

Settings / Authorization Delegations
  → ListAuthorizationDelegations
  → EstablishAuthorizationDelegation [Idempotency-Key]
  → UpdateAuthorizationDelegation [If-Match]
  → RevokeAuthorizationDelegation
```

Governance never grants target-domain write Permission and never executes the target action itself.

### F10 — Administer ordinary access

```text
Shell
  → GetCurrentAccessContext

Settings / Access
  → ListOrganizationMembers
  → ListAccessRoles
  → AssignAccessRole
  → RevokeAccessRole
```

No invitation flow, custom-role designer, IdP-role authority or frontend-local identity registry is admitted.

---

## 5. Screen inventory

`Channel mode` refers only to navigation/filter semantics described in §2.3.

| ID | Screen / route state | Channel mode | Owner(s) read or acted on | Product operations / capability home | Ordinary Permission(s) |
| --- | --- | --- | --- | --- | --- |
| S00 | App Shell / current access | organization-wide | IdentityAccess | `GetCurrentAccessContext` | authenticated special condition |
| S01 | Overview | organization-wide | Portfolio + Economics + Work, read-only composition | `ListMarketplaceInstallations`, `GetEconomicPerformanceSummary`, `ListWork` as available | `portfolio.read`, `economics.read`, `work.read` independently |
| S10 | Readiness workspace | exact-required | ProductChannelReadiness | Search + readiness + requirements + resolve/clear correspondence | `readiness.read`, `readiness.manage` |
| S20 | Marketplace Listings | exact-required | Offering | `ListMarketplaceListings` | `offering.read` |
| S21 | Marketplace Listing detail | exact-required | Offering; optional contextual reads from other owners | `GetMarketplaceListing` | `offering.read` plus permissions of optional composed panels |
| S22 | Listing Intents | all-or-exact | Offering | `ListListingIntents` | `offering.read` |
| S23 | ListingIntent editor/detail | target-explicit | Offering + contextual Readiness/Availability/Economics | `GetListingIntent`, create/update/discard/submit/media; contextual owner Qs | `offering.read`, `listing.manage`, plus contextual read Permissions |
| S24 | Price Intents | all-or-exact | Offering | `ListPriceIntents`, `GetPriceIntent`, contextual `CreatePriceIntent` | `offering.read`, `price.manage` |
| S30 | Availability | all-or-exact | Availability | `ListSellableAvailability`, `GetSellableAvailability` | `availability.read` |
| S40 | Market Intelligence | all-or-exact for list; subject-explicit for detail | MarketIntelligence | `ListCompetitivePositions`, `GetCompetitivePosition`, `ListComparableOffers` | `market.read` |
| S50 | Economics / Expected | all-or-exact | CommercialEconomics | `ListExpectedEconomics`, `GetExpectedEconomics`, `EvaluatePriceScenario` | `economics.read` |
| S51 | Economics / Realized | all-or-exact | CommercialEconomics | `ListSaleEconomics`, `GetSaleEconomics`, `GetEconomicPerformanceSummary` | `economics.read` |
| S52 | Economics / Reconciliation | organization-wide | CommercialEconomics | `ListEconomicAttributions`, `GetEconomicAttribution`, `ResolveEconomicAttribution` | `economics.read`, `economics.reconcile` |
| S60 | Marketplace Sales | exact-required | MarketplaceSales | `ListMarketplaceSales` | `sales.read` |
| S61 | Sale detail / lifecycle composition | exact-required subject | Sales + Economics + Materialization + Fulfillment + PostSale, read-only composition | `GetMarketplaceSale`, `GetSaleEconomics`, related owner reads; owner-specific actions only | component Permissions independently |
| S62 | Sales / ERP Orders | all-or-exact via sale filter | BusinessSystemMaterialization | `ListBusinessOrderIntents`, `GetBusinessOrderIntent`, party/destination reads and party resolution | `materialization.read`, `materialization.resolve` |
| S63 | Sales / Invoicing | organization-wide/context filtered | BusinessSystemMaterialization | `ListInvoicingIntents`, `GetInvoicingIntent` | `materialization.read` |
| S70 | Fulfillment / Executions | all-or-exact via sale filter | Fulfillment | `ListFulfillmentExecutions` | `fulfillment.read` |
| S71 | Fulfillment Execution detail | execution-explicit | Fulfillment | `GetFulfillmentExecution`, checkpoint capabilities, artifact reads | `fulfillment.read`, `fulfillment.execute` |
| S72 | Fulfillment / Shipments | exact-required | Fulfillment | `ListShipments` | `fulfillment.read` |
| S73 | Shipment detail | exact-required | Fulfillment | `GetShipment` | `fulfillment.read` |
| S80 | Post-sale Resolutions | all-or-exact via sale filter | PostSaleResolution | `ListPostSaleResolutions`, `CreatePostSaleResolution` | `post_sale.read`, `post_sale.manage` |
| S81 | Post-sale Resolution detail | organization-wide ID; source Sale remains qualified | PostSaleResolution | `GetPostSaleResolution` | `post_sale.read` |
| S90 | Work Inbox | organization-wide | OperationalWork | `ListWork` | `work.read` |
| S91 | Work detail | organization-wide Work ID | OperationalWork + source link | `GetWork`, assign/clear/hold/resume/escalate | `work.read`, `work.manage` |
| S100 | Approvals / Decisions | organization-wide | ControlledActionGovernance | `ListAuthorizationDecisions` | `governance.read` |
| S101 | Authorization Decision detail/contextual decision | organization-wide target with exact target revision | ControlledActionGovernance | `GetAuthorizationDecision`, contextual `CreateAuthorizationDecision` | `governance.read`, `governance.decide` |
| S110 | Settings / Channels | organization-wide | MarketplacePortfolio | `ListMarketplaceInstallations`, `CreateMarketplaceInstallation` + non-Product OAuth begin after create | `portfolio.read`, `portfolio.manage` |
| S111 | Settings / Channel Installation | exact Installation | MarketplacePortfolio | `GetMarketplaceInstallation`, update config, deactivate | `portfolio.read`, `portfolio.manage` |
| S112 | Settings / Selling Entities | organization-wide | MarketplacePortfolio | `ListSellingEntities` only | `portfolio.read` |
| S113 | Settings / Access | organization-wide | IdentityAccess | member/role lists + assign/revoke | `access.read`, `access.manage` |
| S114 | Settings / Inventory Sources | organization-wide | Availability | list/get/create/update/deactivate inventory source | `availability.read`, `availability.manage` |
| S115 | Settings / Availability Policy | organization-wide | Availability | get/update effective allocation scope policy | `availability.read`, `availability.manage` |
| S116 | Settings / Fulfillment Nodes | organization-wide | Fulfillment | list/get/create/update/deactivate nodes | `fulfillment.read`, `fulfillment.manage` |
| S117 | Settings / Fulfillment Targets | organization-wide | Fulfillment | get/update operating targets | `fulfillment.read`, `fulfillment.manage` |
| S118 | Settings / Commercial Policy | organization-wide | CommercialEconomics | get/update policy | `economics.read`, `economics.policy.manage` |
| S119 | Settings / Authorization Delegations | organization-wide | ControlledActionGovernance | list/create/update/revoke delegations | `governance.manage` |

There are **32 screen/route states** above including shell/overview. Detail drawers may realize some of these states without creating a separate full-page visual, but every state keeps a stable deep-link/navigation identity where source qualification or user recovery requires it.

---

## 6. Exact 95-operation coverage check

The W4 enforcement groups map to the screen homes above as follows.

| W4 group | Count | Frontend home(s) |
| --- | ---: | --- |
| Identity / ordinary access | 5 | S00, S113 |
| Marketplace Portfolio | 6 | S110, S111, S112 |
| Product & Channel Readiness | 5 | S10 |
| Marketplace Listing observation | 2 | S20, S21 |
| ListingIntent authoring/media | 7 | S22, S23 |
| PriceIntent | 3 | S23, S24 |
| Availability | 9 | S30, S114, S115 |
| Market Intelligence | 3 | S40 |
| Expected Economics / scenario | 3 | S50, contextual S23/S21 |
| Sale Economics | 3 | S51, contextual S61 |
| Commercial Policy / Economic Attribution | 5 | S52, S118 |
| Controlled Action Governance | 7 | S100, S101, S119 |
| Marketplace Sales | 3 | S60, S61 |
| Business-System Materialization | 5 | S61, S62 |
| InvoicingIntent | 2 | S61, S63 |
| Fulfillment lifecycle / nodes / targets / checkpoints | 13 | S70, S71, S116, S117 |
| Fulfillment artifacts | 2 | S71 |
| Shipment observation | 2 | S72, S73 |
| Post-Sale | 3 | S61, S80, S81 |
| Operational Work | 7 | S90, S91 |
| **Total** | **95** | **all admitted operations have a frontend interaction home** |

An interaction home does not mean every Principal sees or can invoke every operation. W4 client-class and Permission rules remain authoritative.

---

## 7. State and safety treatment

### 7.1 Server state versus browser state

- TanStack Query-managed server state remains a cache of Product reads, never Product truth.
- Organization/channel/tab/filter/search values are URL/navigation state.
- unsent edits are form-draft state.
- drawer/open/selection/focus are ephemeral UI state.
- no second durable/global business store is admitted.

### 7.2 Knowledge / coverage / freshness

Where the Product contract distinguishes them, UI must visibly distinguish:

```text
known value
known empty
unknown / insufficiently known
unavailable
partial / incomplete
materially stale
```

Rules:

- unavailable/unknown never renders as zero, false, empty success or “nothing found”;
- incomplete collection coverage never produces a global completeness/count claim;
- owner/source observation/evaluation time is shown when material; browser fetch time is never substituted;
- source-product search results always expose sufficient source qualification to retain `SourceInstance + native_product_key` after selection.

### 7.3 Consequential outcomes

Consequential interactions preserve:

```text
accepted
pending
rejected
ambiguous
```

where those owner semantics are reachable.

The UI never equates:

```text
HTTP 2xx = converged
accepted = completed
pending = failed
ambiguous = failed
ambiguous = safe to retry
```

Ambiguous possible external acceptance removes generic automatic retry. The user is directed to owner state, reconciliation evidence or Work as applicable.

### 7.4 Idempotency

When the OAD requires `Idempotency-Key`:

1. a key is created for one semantic intake;
2. it survives a network/lost-response retry of that same semantic request;
3. editing the semantic request creates a new key;
4. the key never becomes Intent identity or provider-retry authorization;
5. retry UI first resolves current owner/intake state rather than blindly redispatching an ambiguous effect.

### 7.5 Concurrency

Standard resource updates using `If-Match`:

- `412` means the client draft was based on stale owner state;
- UI reloads current owner representation and asks the user to re-decide rather than auto-overwrite.

Owner `:verb` capabilities carrying typed `etag`:

- stale revision conflict remains distinct from business/provider rejection;
- UI rereads current resource/meaning before enabling a new semantic decision.

### 7.6 Governance and Work

- Governance is shown only when authorization meaning actually applies; it does not become a generic modal around every mutation.
- Work surfaces responsibility/assignment/escalation but never owns source-domain closure.
- a source condition may link to Work; Work may link back to its source subject; neither may mutate the other through frontend-local inference.

### 7.7 Permission-conditioned visibility

Current `AccessContext.permissions` may hide irrelevant navigation/actions for usability.

It is never authorization authority. Deep links and all Product requests remain server-enforced. A hidden menu/button cannot turn a server `403/404` into a frontend business rejection.

---

## 8. Explicit negative controls / forbidden UX

The D6 candidate is invalid if any wireframe or later implementation introduces these behaviors:

1. hardcoded/default Organization used as Product scope;
2. `display_name` used as identity/cache/security key;
3. Marketplace Installation selection treated as authentication/authorization;
4. `All accounts` synthesized by merging independently paginated Listing/Sales/Shipment collections;
5. Amazon/Shopee shown as currently connectable before responsible support is admitted;
6. generic Integration/Channel Catalog Product API invented for UI symmetry;
7. frontend `ConnectMarketplace` business operation invented around OAuth;
8. SourceInstance ID hardcoded or stored as hidden frontend authority;
9. omitted source-search filter interpreted as one ambient/default source;
10. Product master/product CRUD invented from source-product search;
11. direct Listing create/update/pause/close bypassing ListingIntent;
12. ListingIntent owning price or Sellable Availability;
13. direct `SetPrice`, `SetAvailableQuantity`, `Sync*`, `Refresh*`, `CollectMarketNow` Product commands invented;
14. client-side Economics or Readiness conclusion becoming write authority;
15. generic “success” toast for pending/ambiguous consequential state;
16. generic retry on ambiguous external effect;
17. stale `If-Match`/etag conflict silently overwritten;
18. direct BusinessOrder/Invoicing/Sankhya command from React;
19. frontend-calculated OperationalStage used to authorize or trigger work;
20. physical checkpoint enabled for caller-controlled “qualified device” claims;
21. artifact visibility widened from `fulfillment.execute` merely for UI convenience;
22. direct PostSale close/refund/cancel/provider-action invented;
23. direct Work create/close/dismiss/resolution command invented;
24. Governance decision used as target-domain Permission or execution authority;
25. invitation/custom-role/source-registry/admin screens invented without admitted Product capability;
26. reactivation controls shown for Marketplace Installation, Inventory Source or Fulfillment Node when explicit reactivation remains deferred;
27. Overview/card counts inferred from one collection page or partial coverage;
28. read-only composition becoming concurrency/write authority;
29. provider/native vocabulary or IDs losing their accepted source qualifier;
30. browser/request time shown as source/owner freshness.

---

## 9. Low-fidelity wireframe proof set

The smallest HTML wireframe set capable of falsifying this interaction model is:

1. **Shell + Overview** — Organization context, permission-conditioned nav, read-only composition;
2. **Readiness workspace** — exact channel context, source-product search without default SourceInstance, knowledge states, correspondence action;
3. **ListingIntent editor** — requirements, independent price/availability/economics panels, ETag, submission/outcome states;
4. **Availability workspace** — all/exact channel filter, honest desired/provider/convergence state, no manual quantity mutation;
5. **Sale detail** — read-only cross-owner lifecycle composition with owner-local actions;
6. **Fulfillment execution** — physical checkpoints, artifacts, idempotency/current revision, stalled/exception visibility;
7. **Work inbox/detail** — assignment/hold/resume/escalation with source-owner closure fence;
8. **Channels settings** — available kind versus connected Installation, OAuth boundary, deactivation/no reactivation;
9. **Economics / Reconciliation** — expected/realized/attribution distinction and scenario evaluation;
10. **Access / Governance settings** — presentation identity, role assignment/revocation, delegations without custom-role IAM platform.

One self-contained HTML prototype may contain these ten navigable wireframe states; it does not need ten independent frontend applications or final visual styling.

---

## 10. Current decision and next proof

The App Shell / IA is operator-approved. This interaction map is the derived D6-B1 screen/flow candidate.

Next D6 work:

1. produce the low-fidelity HTML proof set in §9;
2. attack the prototype against the negative controls in §8;
3. adjudicate any newly exposed Product/API gap before continuing that interaction;
4. only after the interaction model and wireframes cohere, evaluate the smallest frontend feature/package topology and concrete dependencies required by those accepted properties.

D7–D9 and Product implementation remain blocked.