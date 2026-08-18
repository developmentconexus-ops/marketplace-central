# D5-B2 — Product Operation / Resource Surface — REVIEW CANDIDATE

> **Status:** NON-AUTHORITATIVE REVIEW CANDIDATE  
> **Stage:** D5 — API  
> **Batch:** B2 — Product Operation / Resource Surface  
> **Parent authority:** accepted D0 → D4 + D5-B1  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Prepared at:** `ba7e54450a6dd9b51d65780cbe0ade3012c46bad` on `docs/global-methodology-alignment`  
> **Purpose:** bounded package for independent adversarial review. This file does not change D5 authority, router status, ADR status, D0–D4, implementation permission or the current OpenAPI.

---

## 1. Question this batch must answer

D5-B1 already decided the Product API laws. B2 must now answer:

> **What is the smallest external Product 1.0 operation/resource surface that real MPC actors need in order to operate the accepted marketplace-control-plane lifecycle, while keeping internal orchestration, acquisition, provider protocol and runtime mechanics out of the Product API?**

B2 starts from accepted product outcomes and owners, never from the legacy route list.

The resulting surface must let Product 1.0 actors:

- understand operational attention;
- establish marketplace participation/configuration;
- prepare Product↔channel readiness/correspondence;
- inspect/control marketplace offerings and price intent;
- configure and observe sellable availability;
- inspect competitive position and economics;
- perform/authorize controlled actions;
- observe marketplace sales and business-system materialization;
- execute the accepted fulfillment path;
- coordinate essential post-sale resolution;
- progress explicit Operational Work;
- manage ordinary product access needed by the Owner/Administrator actor.

It must **not** turn every internal D1/D3/D4 operation into a public endpoint.

---

## 2. Operation-admission predicate

A Product API operation is admitted only when a named Product 1.0 actor/client must do at least one of:

1. **Q — observe owner meaning** needed for a user-visible decision, operation, explanation or inspection;
2. **C — ask the owner to create/change owner-owned MPC state or perform owner-owned work** as part of a Product 1.0 actor flow;
3. **P — consume an explicitly justified cross-owner read projection** required by a D0 user-observable outcome.

An operation is **not** admitted merely because:

- an adapter needs it;
- a worker/job needs it;
- current code exposes it;
- a domain has an internal method for it;
- an external provider has a resource for it;
- a current UI happened to call it;
- it is useful for debugging/probing;
- it could make a future provider easier;
- a batch/import/refresh/sync mechanism exists today.

### Internal-reaction rule

A D3 `E` edge or autonomous owner reaction does not become Product API merely because it performs important work.

Examples:

- Sale commit waking Materialization;
- routine Availability synchronization;
- conference evidence making Materialization eligible to create/advance Invoicing Intent;
- provider notification acquisition;
- Payment/market evidence acquisition;
- reconciliation/recovery sweeps.

These remain owner/runtime behavior unless a real external Product 1.0 actor has an independent Q/C/P need.

---

## 3. Evidence classification

### 3.1 KNOWN from authority

1. Product 1.0 is **Marketplace Operations + Commercial Intelligence**, not an integration console or generic ERP/marketplace platform.
2. D0 names four operational actor classes: Marketplace Operations Operator; Fulfillment/Dispatch Operator; Commercial/Marketplace Manager; Owner/Administrator/Policy Approver.
3. D0 completion requires portfolio-driven attention, readiness/linkage, offer operation, availability convergence, competitive/economic analysis, controlled action, sale→business-order→invoicing→fulfillment, shipment visibility, essential post-sale, explicit Work, settlement/economic explainability and operable MPC-owned governance.
4. D1 owns twelve business boundaries and rejects a generic Mutation/Action domain.
5. Product master remains external; no MPC Product Master API is justified.
6. Marketplace Listing/Sale/Shipment and native financial movements remain source-qualified external identities; no synthetic aliases are required merely for normalization.
7. Material Listing/Price/Availability/Business Order/Invoicing/Fulfillment intents are owner-local when durable; there is no universal Intent resource.
8. Operational Work, Post-Sale Resolution and Authorization Decision are accepted MPC-owned canonical identities.
9. Availability routine synchronization is automatic on the normal sufficiently-known policy-valid path; per-change human approval is not required.
10. Sales fan-out to Materialization/Fulfillment/Economics is internal D3 progression, not an API command from the client.
11. Materialization alone owns Business Order/Invoicing Intent; Fulfillment conference evidence may make invoicing eligible but Fulfillment never creates Invoicing Intent.
12. Provider/business-system protocol, notifications, OAuth handshakes, raw provider status/resources and source mechanics remain D4.
13. D4 provider-rich evidence may be exposed only through the owning semantic contract for a named consumer/correctness need.
14. D5-B1 fixes Organization Product API scope at `/organizations/{organization_id}/...`, source-qualified external identity, honest knowledge/freshness, fail-closed consequential idempotency, RFC 9457 API problems, hard cutover and one OpenAPI wire authority.
15. Exact Permission→operation mapping is a D5 responsibility.

### 3.2 INFERRED

1. The smallest sustainable surface is **actor/use-driven and owner-oriented**, not “one CRUD set per D1 domain”.
2. Internal acquisition/reconciliation/run-control endpoints are accidental Product API complexity unless a named actor needs them independently.
3. Several externally authoritative objects are necessary **references/evidence**, but do not need top-level Product API resources.
4. D0 operational attention justifies one cross-owner read projection even before D6 chooses the final screen topology.
5. Successful and failed domain Intents must be individually inspectable when they are stable historical/correlation anchors, but internal creation need not always be externally invokable.

### 3.3 UNKNOWN

B2 intentionally does not settle yet:

- final request/response field schemas;
- exact status-code mapping per operation beyond D5-B1 laws;
- exact ETag/precondition token encoding;
- exact representation/addressing of relationship resources that D2 did not give canonical IDs, especially Product↔channel correspondence and authorization delegation/grant state;
- exact provider-enrichment union members per operation;
- exact technical connection/setup surface for SourceInstance/provider credentials outside the Product API;
- final role bundles over the Permission catalog;
- concrete OpenAPI generation/server-conformance technology;
- D6 screen composition;
- D7 runtime mechanism.

Unknown remains Unknown. None of these authorizes a generic resource/ID/framework by convenience.

---

## 4. Root cause / target invariant

### Root cause

The legacy external surface exposes historical module boundaries and technical mechanisms as first-class API concepts. That makes the client depend on implementation choreography instead of product semantics.

Examples in current-state evidence include generic `/mutations`, import/refresh/sync/probe/run endpoints, provider/integration nouns, provider category resources, manual profitability calculation/import routes and Sankhya-specific linkage operations.

### B2 target invariant

> **Every admitted Product API operation exists for a named Product 1.0 actor/use, maps to exactly one accepted owner/substrate authority and Q/C/P class, exposes only the owner meaning required by that use, and leaves provider acquisition, autonomous owner reactions and runtime orchestration internal unless a separate real Product API consumer is proven.**

Corollaries:

1. Important internal work does not imply an endpoint.
2. Externally authoritative resources may be source-qualified references/evidence without becoming top-level Product API resources.
3. A Product API list/view does not become a second owner merely because it aggregates or serializes owner meaning.
4. A client never orchestrates the accepted cross-owner workflow by manually calling each internal D3 edge.
5. Exception repair happens through the owning semantic capability and/or Operational Work, never through a generic retry/sync/mutation endpoint.

---

## 5. Credible alternatives

### Alternative A — Preserve legacy route families and rename them

Examples: `/mutations` → `/actions`, `/integrations` → `/marketplace-installations`, `/orders/*/sankhya-linkage` → `/materialization/*` while retaining the same underlying API choreography.

**Rejected — Local Maximum.** Renaming leaves client-visible mechanism/orchestration as product contract.

### Alternative B — Expose every public domain/internal method over HTTP

Each D1 owner gets broad CRUD + explicit internal progression/reconciliation operations.

**Rejected — overexposure.** It turns internal semantic edges/runtime mechanics into a second orchestration protocol and forces clients to know workflow choreography.

### Alternative C — One generic control surface

Generic resources/actions/workflows/entities with a common command/result envelope.

**Rejected — contradicts D1/D2/D5-B1.** This recreates Mutation/Workflow/Provider authority through the API.

### Alternative D — Actor/use-driven semantic owner surface; internal mechanics remain internal

Expose only real Q/C/P Product 1.0 needs. MPC-owned configuration/resources are ordinary resources where truthful. Consequential actions use owner-specific Intents/capabilities. Provider-rich evidence remains qualified inside owner responses. Internal autonomous progression, acquisition and reconciliation are not client choreography.

**Recommended Global Maximum.** It is the smallest surface that satisfies D0 outcomes without losing D1–D5-B1 correctness.

---

# 6. Candidate Product API resource families

The following are **semantic resource families**, not final OpenAPI schemas. Exact tuple addressing for source-qualified/relationship subjects may be sharpened later without changing the admitted operation set.

All Organization-owned families live under:

```text
/organizations/{organization_id}/...
```

One platform-scoped read family for product-defined AccessRole definitions is admitted because those definitions have no Organization-specific business meaning under D2.

## 6.1 Identity/access substrate

### Resources / operations

- `AccessRole` definitions — **list/read only**, platform-scoped.
- Organization Membership — list/get existing memberships.
- RoleAssignment — assign/revoke a product-defined role for an existing Membership.

### Explicit non-surface

- end-user credentials/session management;
- custom role designer;
- generic ACL/ReBAC;
- IdP user directory mirroring;
- arbitrary Principal CRUD without a proven onboarding consumer.

Membership provisioning/binding mechanics not needed by the first Product 1.0 actor flow remain a later auth/setup concern.

---

## 6.2 Marketplace Portfolio

### `MarketplaceInstallation`

Admit:

- list installations;
- get one installation;
- create Organization marketplace participation/configuration;
- update/deactivate MPC-owned participation/configuration.

Installation reads may expose Portfolio-owned operational posture/attention derived from D4 evidence, but never credential/token DTOs.

### `SellingEntity`

Admit:

- list/get;
- create;
- update/deactivate organization-facing Selling Entity configuration.

External company/legal records remain source-qualified references and are not mastered by this resource.

### Installation ↔ eligible Selling Entity participation

Admit:

- list eligible participation;
- add/remove one explicit Selling Entity participation relation.

No synthetic relationship entity/ID is introduced unless later evidence proves independent lifecycle meaning.

### Protocol fence

OAuth callbacks, credential submission/rotation and provider-native auth state remain protocol/technical surface. A Product API Installation may expose semantic connection posture without exposing the auth machinery itself.

---

## 6.3 Product & Channel Readiness

No standalone MPC Product Master resource is admitted.

### `ProductReadiness`

Subject is explicit:

```text
Organization
+ Marketplace Installation
+ SourceInstance-qualified Product
```

Admit:

- list/search source products with Readiness-owned conclusion/context;
- get exact Product Readiness;
- read missing/conflicting requirements and evidence sufficiency;
- read bounded correspondence candidates/evidence needed for operator adjudication.

Product descriptors required for the readiness use may be returned from external source evidence; that does not create Product master authority.

### Product↔channel correspondence

Admit owner-specific capabilities to:

- establish one explicit correspondence;
- remove/supersede one explicit correspondence.

The semantic key is the source-qualified Product + source-qualified Offering/Variation relationship. B2 does not invent a canonical `CorrespondenceID` merely to make REST addressing convenient.

No `generate candidates` job endpoint is admitted; candidate evidence is a Q result over current owner/source meaning.

---

## 6.4 Marketplace Offering Operations

### `Offering`

Externally authoritative Listing identity remains:

```text
Marketplace Installation
+ provider-native Listing key
(+ Variation key where applicable)
```

Admit:

- list Offerings;
- get one Offering current semantic view;
- expose Listing/Price convergence and current provider-enriched evidence only where Offering needs it.

### `ListingRequirements`

Admit a Q surface for information a real operator needs before a Listing Intent, such as source-qualified provider category/attribute/constraint evidence translated into Offering-owned listing requirements.

Do **not** expose provider category/resource taxonomy as an independent Product API ontology.

### `ListingIntent`

Admit:

- create an owner-specific Listing Intent for a material create/change/close operation supported by Product 1.0;
- get one Listing Intent;
- list Listing Intents for operator history/pending/ambiguous/convergence inspection.

### `PriceIntent`

Admit:

- create Price Intent;
- get one Price Intent;
- list Price Intents.

Price Intent remains distinct from Economics analysis/recommendation. Economics never creates the marketplace write.

No generic preview/approve/retry Mutation endpoints survive. Analysis/readiness/economics are owner Qs; authorization is Governance; retry after ambiguity is reconciliation, not a generic user command.

---

## 6.5 Availability Control

### `InventorySource`

Admit:

- list/get;
- create;
- update/deactivate.

Source-native company/location/stock selectors remain source-qualified configuration/evidence and never become Inventory Source identity by convenience.

Inventory Scope eligibility may be represented inside the Availability-owned configuration surface; B2 does not force a separate canonical `InventoryScopeID`.

### `AvailabilityPolicy`

Admit:

- read effective Organization baseline policy/provenance;
- update the MPC-owned Organization baseline allocation/automation policy.

No generic policy-scope DSL is admitted. More-specific overrides become new API surface only when a proven business scope requires them, preserving D0's default/inheritance/override seam without implementing speculative override dimensions now.

### `SellableAvailability`

Admit:

- list/query current sellable availability by real operational subjects;
- get exact current sellable availability/convergence for a source-qualified Product/Offering target;
- preserve knowledge/freshness/source/convergence semantics.

### `AvailabilityIntent`

Routine Availability Intent creation is **internal owner behavior** for the accepted automatic path, not a public command.

Admit read-by-id (and only if later required, bounded history listing) so Work/attention/audit references can explain a material synchronization attempt. No public `create availability mutation` / `manual apply` baseline operation is admitted.

---

## 6.6 Market Intelligence

### `CompetitivePosition`

Admit:

- list/query competitive positions for Product/Offering subjects;
- get one current competitive interpretation including sufficiency/knowledge/freshness.

### `ComparableOffer`

Admit a bounded read collection when the operator needs to inspect comparable-market evidence behind the competitive interpretation.

Provider-rich evidence such as Mercado Livre `price_to_win`, buyer-facing shipping/free-shipping, competition status/boosts/reasons remains source-qualified enrichment inside these Market Intelligence responses.

No public `collect market data`, scraper, provider offer probe or acquisition-run operation is admitted.

---

## 6.7 Commercial Economics

### Price simulation / Expected Economics

Admit a non-durable owner Q operation that evaluates a candidate price/context and returns explainable Expected Economics under an explicit Cost Basis.

The transport may use POST when request-body complexity requires it; semantic class remains Q. The calculation does not gain canonical identity merely because POST is used.

Admit a Q for current Expected Economics where a Product/Offering context already exists.

### `SaleEconomics`

Admit:

- list sale-level economic summaries for operational/commercial review;
- get one sale's explainable L1 Order Economics + L2 realized/settlement/reconciliation meaning, preserving the rungs rather than flattening them.

Payment/refund/fee/shipment/native financial objects remain source-qualified evidence inside Economics; no top-level MPC Payment/Refund/Settlement resource is created.

### Commercial Economics policy

Admit Organization-baseline read/update of MPC-owned:

- margin floors;
- price boundaries;
- economic approval thresholds/other proven commercial thresholds;
- explainable provenance.

### Cost Basis configuration

Admit read/update of the current Organization baseline Cost Basis selection/interpretation configuration.

Do not create a universal policy engine or provider cost-type aliases. More-specific policy overrides remain unadmitted until a proven business scope requires them.

Legacy import/calculate/manual-adjustment job endpoints do not survive. Ambiguous economic attribution requiring human adjudication is progressed through explicit Work/source-owner resolution, not a generic manual-adjustment API.

---

## 6.8 Controlled Action Governance

### `AuthorizationDecision`

Admit:

- list decisions, especially pending/current review work;
- get one decision with authority context and referenced owner Intent;
- record an authorized/rejected decision outcome when the Principal has ordinary API access **and** applicable Governance authority.

Ordinary Permission never substitutes for Grant/Delegation/business validity.

### Authorization Grant / Delegation configuration

Admit:

- list current delegation/grant configuration;
- create/update/revoke bounded authorization delegation semantics required for organization-operable Governance.

D2 intentionally did not freeze a canonical Grant ID. B2 admits the operation family but leaves exact REST addressing/cardinality to a later contract-shape decision; no canonical identity is invented merely for CRUD convenience.

---

## 6.9 Marketplace Sales

Marketplace Sales is external-observation/canonical-interpretation authority, not an externally creatable MPC order resource.

### `Sale`

Identity remains:

```text
Marketplace Installation
+ provider-native Sale/Order key
```

Admit:

- list Sales with real operational filters;
- get one Sale semantic interpretation/context/correlation including transaction-specific Selling Entity attribution.

No Product API `import orders`, `mark faturado`, provider status patch or manual sale-ingest command is admitted.

---

## 6.10 Business-System Materialization

### `BusinessOrderIntent`

Admit:

- get one Business Order Intent by its canonical MPC identity;
- resolve from Sale detail/Work/attention to that Intent;
- optionally list only if independent operator/support evidence proves a real queue/history need beyond Sales + Work.

**Do not admit public creation on the baseline.** Sales→Materialization autonomous progression creates/advances the owner Intent under D3; the client does not orchestrate that edge.

The read exposes source-qualified native business-order result/correlation, acceptance/convergence and relevant prerequisite status without exposing TOP/NUNOTA choreography as MPC semantics.

### `InvoicingIntent`

Admit:

- get one Invoicing Intent by canonical identity;
- optionally list only on the same real-consumer predicate.

**Do not admit public creation on the baseline.** Fulfillment conference/readiness makes Materialization eligible; Materialization alone creates/advances the Intent.

### Party Resolution / Destination Realization

Do not create Product API Customer/Party/Address resources.

Ambiguous/unsupported cases become source-domain conditions + Operational Work. Human resolution evidence is submitted through Work and adjudicated by Materialization under its own semantics.

No Sankhya-linkage candidate/confirm surface survives as provider-shaped Product API.

---

## 6.11 Fulfillment Lifecycle

No new universal `FulfillmentID` is introduced merely for API addressing. Baseline fulfillment operations are scoped from the source-qualified Sale and, where needed, explicit Fulfillment-owned scope.

### Fulfillment current state

Admit:

- list sale fulfillments/work queue for the Fulfillment Operator;
- get current Fulfillment state for an exact Sale/scope;
- expose selected/eligible Fulfillment Node, physical state, provider-requirement closure, material external/internal deadlines, current Shipment/delivery observation and explicit blockers.

Provider Shipment identity remains source-qualified evidence; a separate synthetic MPC Shipment resource is not required.

### Fulfillment capabilities

Admit owner-specific C operations for the accepted physical path:

- select/confirm Fulfillment Node/routing where operator choice is required;
- record separation completion/evidence;
- record physical conference result/evidence;
- prepare provider dispatch requirements/artifacts when an explicit operator action is required;
- read dispatch readiness + required operator artifacts;
- record packing completion;
- confirm verified dispatch handoff.

These are Fulfillment-owned business meanings, not a generic workflow engine. Provider-native intermediate steps remain D4.

Conference does not create Invoicing Intent in the API. It changes Fulfillment-owned meaning; Materialization reacts internally under D3.

Shipment/delivery remains visible through this Fulfillment read surface until relevant terminal outcome. Material exceptions become Work.

### Fulfillment policy

Admit Organization-baseline read/update for MPC-owned internal operational targets applicable to the accepted fulfillment path. External provider deadlines remain read-only external evidence and are never editable through this policy.

---

## 6.12 Post-Sale Resolution

### `PostSaleResolution`

Admit:

- list Resolutions;
- get one Resolution;
- initiate one explicitly scoped material Resolution when an authorized Product 1.0 actor starts a cancellation/return/refund consequence flow.

Provider-originated post-sale facts may create/reconcile a Resolution internally without a client POST.

The creation request identifies the material scope/desired resolution obligation without collapsing Cancellation, Return and Refund into one mutually-exclusive universal enum. The exact consequence schema remains later contract work.

Post-Sale coordinates consequence owners internally. No generic `refund`, `cancel everything` or cross-owner workflow command is exposed.

---

## 6.13 Operational Work

### `Work`

Admit:

- list Work by real responsibility/assignment/state/urgency/origin filters;
- get one Work item;
- assign/reassign/unassign responsibility where legitimate;
- escalate;
- submit/point to resolution evidence for source-owner adjudication.

No arbitrary `dismiss` / `force close` operation is admitted.

Submitting resolution evidence does not declare source truth resolved. The source owner evaluates resolved/unresolved/unknown-or-pending; Work reconciles/ closes its own lifecycle accordingly.

Work detail may carry typed references to originating conditions/evidence; it never becomes a universal entity graph or source-domain truth owner.

---

## 6.14 Operational Attention projection

D0 completion outcome 1 requires portfolio-driven attention before D6 chooses a screen.

Admit one read-only `OperationalAttention` P surface that may compose minimally sufficient attention from accepted owners, including where material:

- installation posture;
- readiness/offering/availability divergence;
- approval-required/pending action attention;
- material evidence freshness/coverage degradation;
- fulfillment/deadline attention;
- post-sale exceptions;
- unresolved Work.

Rules:

- no canonical `AttentionID` business entity is created;
- attention items use typed owner/subject references;
- the projection cannot mutate source state;
- it cannot be consequential write/concurrency authority;
- it must preserve coverage/freshness honesty for the universe it claims;
- its ordinary Permission authorizes the bounded summary only; following an item into owner detail still requires that owner's Permission.

No generic `/dashboard` aggregate or hidden owner status is authority.

---

# 7. Candidate ordinary Permission catalog

Permissions are product-defined ordinary API access capabilities under D2. They do **not** grant business disposition, Governance authority or execution-time validity.

Proposed stable catalog:

| Permission | Ordinary API access |
|---|---|
| `operations.attention.read` | OperationalAttention projection |
| `portfolio.read` | Installation / Selling Entity reads |
| `portfolio.manage` | Installation / Selling Entity / participation configuration |
| `readiness.read` | ProductReadiness / correspondence evidence reads |
| `readiness.manage` | establish/remove correspondence |
| `offering.read` | Offering, requirements and Intent reads |
| `offering.operate` | create Listing/Price Intent |
| `availability.read` | Inventory Source / Sellable Availability / Availability Intent reads |
| `availability.manage` | Inventory Source + Availability policy configuration |
| `market_intelligence.read` | competitive position/comparable evidence |
| `economics.read` | simulation/expected/sale economics |
| `economics.manage` | Cost Basis + Economics policy configuration |
| `governance.read` | Authorization Decision/delegation reads |
| `governance.decide` | record Authorization Decision outcome; authority still independently required |
| `governance.manage` | configure Grant/Delegation semantics; authority still independently required |
| `sales.read` | Sale list/detail |
| `materialization.read` | BusinessOrder/Invoicing Intent reads |
| `fulfillment.read` | Fulfillment queue/state/readiness/artifacts/Shipment observations |
| `fulfillment.operate` | physical progression / routing / dispatch capabilities |
| `fulfillment.manage` | internal operational-target policy configuration |
| `post_sale.read` | Post-Sale Resolution reads |
| `post_sale.operate` | initiate a Resolution |
| `work.read` | Work list/detail |
| `work.operate` | assignment/escalation/resolution submission |
| `access.read` | AccessRole / Membership reads |
| `access.manage` | RoleAssignment changes |

### Permission YAGNI rule

B2 does not create one Permission per endpoint. The permission boundary follows a stable user capability whose operations should normally travel together. If review shows any pair has materially different least-privilege risk, split only that pair.

B2 does not freeze AccessRole bundles. D0 actors inform later product-defined bundles; an operational/business responsibility name does not automatically become an AccessRole.

---

# 8. Operation / owner matrix

This is the B2 admission ledger. `Key` means B1 `Idempotency-Key` is required by default. `Structural` means the operation must prove a stable owner-semantic duplicate predicate before it may omit the key.

| Operation family | Real consumer | Owner | Class | Ordinary Permission | Identity / scope | Write safety |
|---|---|---|---|---|---|---|
| list/get AccessRole | Owner/Admin | D2 access substrate | Q | `access.read` | platform definition | read |
| list/get Membership | Owner/Admin | D2 access substrate | Q | `access.read` | Organization + Principal | read |
| assign/revoke RoleAssignment | Owner/Admin | D2 access substrate | C | `access.manage` | Organization + Principal + AccessRole | structurally idempotent PUT/DELETE candidate; no business authorization semantics |
| list/get Installation | Marketplace Ops / Admin | Portfolio | Q | `portfolio.read` | Organization + Installation | honest posture evidence |
| create Installation | Admin | Portfolio | C | `portfolio.manage` | server-owned Installation identity | durable creation: key unless exact creation contract proves duplicate-safe |
| update/deactivate Installation | Admin | Portfolio | C | `portfolio.manage` | Installation | conditional update where stale config is material |
| list/get SellingEntity | Admin/Manager | Portfolio | Q | `portfolio.read` | Organization + SellingEntity | read |
| create/update SellingEntity | Admin | Portfolio | C | `portfolio.manage` | SellingEntity + qualified external legal refs | create key if POST; stale update conditional where material |
| add/remove Installation↔SellingEntity participation | Admin | Portfolio | C | `portfolio.manage` | semantic relation tuple | Structural candidate; no synthetic relation ID |
| list/search ProductReadiness | Marketplace Ops | Readiness | Q | `readiness.read` | Installation + SourceInstance Product | cursor/search; honest source state |
| get ProductReadiness | Marketplace Ops | Readiness | Q | `readiness.read` | exact Product + Installation | knowledge/freshness |
| read correspondence candidates | Marketplace Ops | Readiness | Q | `readiness.read` | Product + Installation | no generate-job command |
| establish/remove correspondence | Marketplace Ops | Readiness | C | `readiness.manage` | Product + source-qualified Offering/Variation | Structural candidate; current evidence/precondition preserved |
| list/get Offering | Marketplace Ops | Offering | Q | `offering.read` | Installation + native Listing/Variation | cursor; source-qualified external identity |
| get ListingRequirements | Marketplace Ops | Offering | Q | `offering.read` | Installation + Product/category context | provider-enriched requirements only |
| create ListingIntent | Marketplace Ops | Offering | C | `offering.operate` | canonical ListingIntent + target scope | **Key**; disposition/Governance/execution validity separate |
| list/get ListingIntent | Marketplace Ops / Approver | Offering | Q | `offering.read` | ListingIntent | cursor for list |
| create PriceIntent | Marketplace Ops | Offering | C | `offering.operate` | canonical PriceIntent + target scope | **Key**; Economics informs but does not write |
| list/get PriceIntent | Marketplace Ops / Approver | Offering | Q | `offering.read` | PriceIntent | cursor for list |
| list/get InventorySource | Manager / Marketplace Ops | Availability | Q | `availability.read` | InventorySource | bounded config read |
| create/update InventorySource | Manager | Availability | C | `availability.manage` | InventorySource + source-qualified selectors | create key if POST; conditional update where material |
| get/update AvailabilityPolicy baseline | Manager | Availability | Q/C | `availability.read` / `availability.manage` | Organization baseline | conditional update; no generic scope DSL |
| list/get SellableAvailability | Marketplace Ops | Availability | Q | `availability.read` | Product/Offering target | cursor; knowledge/freshness/convergence |
| get AvailabilityIntent | Marketplace Ops / Work | Availability | Q | `availability.read` | canonical AvailabilityIntent | read; no baseline public create |
| list/get CompetitivePosition | Marketplace Ops / Manager | Market Intelligence | Q | `market_intelligence.read` | Product/Offering subject | cursor; sufficiency/freshness |
| list ComparableOffer evidence | Marketplace Ops / Manager | Market Intelligence | Q | `market_intelligence.read` | competitive subject | bounded/cursor when needed; provider-rich evidence |
| simulate price | Marketplace Ops / Manager | Commercial Economics | Q | `economics.read` | explicit Product/Offering + candidate context | non-durable calculation; POST allowed without creating identity |
| get Expected Economics | Marketplace Ops / Manager | Commercial Economics | Q | `economics.read` | Product/Offering context | exact/provenance/knowledge |
| list/get SaleEconomics | Manager / Marketplace Ops | Commercial Economics | Q | `economics.read` | source-qualified Sale | cursor on list; L1/L2/rungs stay distinct |
| get/update Economics policy baseline | Manager | Commercial Economics | Q/C | `economics.read` / `economics.manage` | Organization baseline | conditional update |
| get/update Cost Basis baseline | Manager | Commercial Economics | Q/C | `economics.read` / `economics.manage` | Organization baseline | conditional update; source variants never policy aliases |
| list/get AuthorizationDecision | Approver | Governance | Q | `governance.read` | AuthorizationDecision | cursor/pending filter |
| record authorization outcome | Manager/Owner when authorized | Governance | C | `governance.decide` | Decision + referenced Intent/scope | **Key** by default; Governance authority separately proven |
| list delegation/grant config | Owner/Admin | Governance | Q | `governance.read` | accepted delegation semantics | read |
| configure/revoke delegation | Owner/Admin | Governance | C | `governance.manage` | no invented canonical GrantID | security/consequential; conditional/idempotent contract required |
| list/get Sale | Marketplace Ops / Fulfillment / Manager | Marketplace Sales | Q | `sales.read` | Installation + native Sale key | cursor; source-qualified identity |
| get BusinessOrderIntent | Marketplace Ops / Fulfillment / Work | Materialization | Q | `materialization.read` | canonical Intent + Sale/native result refs | no baseline public create |
| get InvoicingIntent | Fulfillment / Marketplace Ops / Work | Materialization | Q | `materialization.read` | canonical Intent + native result refs | no baseline public create |
| list Fulfillment work | Fulfillment Operator | Fulfillment | Q | `fulfillment.read` | source-qualified Sale + owner scope | cursor; stage/deadline/exception filters |
| get Fulfillment state | Fulfillment Operator | Fulfillment | Q | `fulfillment.read` | Sale/scope | provider requirements + shipment observation source-qualified |
| select/confirm Fulfillment Node/routing | Fulfillment Operator | Fulfillment | C | `fulfillment.operate` | Sale/scope + FulfillmentNode | **Key** unless Structural proof |
| record separation | Fulfillment Operator | Fulfillment | C | `fulfillment.operate` | Sale/scope | **Key** unless owner checkpoint identity proves Structural safety |
| record conference | Fulfillment Operator | Fulfillment | C | `fulfillment.operate` | Sale/scope | **Key** unless owner checkpoint identity proves Structural safety; may wake Materialization internally |
| prepare dispatch requirements/artifacts | Fulfillment Operator | Fulfillment | C | `fulfillment.operate` | Sale/scope | **Key** where external/provider effect reachable |
| get dispatch readiness/artifacts | Fulfillment Operator | Fulfillment | Q | `fulfillment.read` | Sale/scope | provider evidence translated; no raw provider artifact ontology |
| record packing | Fulfillment Operator | Fulfillment | C | `fulfillment.operate` | Sale/scope | **Key/Structural** adjudicated by operation semantics |
| confirm dispatch handoff | Fulfillment Operator | Fulfillment | C | `fulfillment.operate` | Sale/scope | **Key**; provider applied/converged distinction preserved |
| get/update Fulfillment internal-target policy | Manager | Fulfillment | Q/C | `fulfillment.read` / `fulfillment.manage` | Organization baseline | external deadlines read-only; conditional update |
| list/get PostSaleResolution | Marketplace Ops / Manager | Post-Sale | Q | `post_sale.read` | canonical Resolution + Sale scope | cursor on list |
| initiate PostSaleResolution | Marketplace Ops | Post-Sale | C | `post_sale.operate` | canonical Resolution + explicit consequence scope | **Key**; consequence owners react internally |
| list/get Work | operational actors | Operational Work | Q | `work.read` | Work identity | cursor; responsibility/assignment/state/origin filters |
| assign/reassign/unassign Work | responsible actor/manager | Operational Work | C | `work.operate` | Work + Principal | conditional update / structural idempotency |
| escalate Work | responsible actor/manager | Operational Work | C | `work.operate` | Work | key/conditional semantics according to occurrence model |
| submit resolution evidence | responsible actor | Operational Work + source-owner C | C | `work.operate` | Work + typed source reference/evidence | **Key** where duplicate submission could create consequential reaction; Work cannot self-declare source resolution |
| list OperationalAttention | Marketplace Ops / Manager / Admin | P only | P | `operations.attention.read` | typed owner/subject refs | cursor; coverage/freshness; summary-only authority |

### Optional-list rule for internal-created Intents

B2 baseline requires `get` for any durable Intent that can be referenced by Work/Authorization/history. A collection/list operation is admitted only when the named actor has a real independent queue/history use.

Therefore:

- ListingIntent / PriceIntent lists — **admitted** because actors create/review them directly.
- AvailabilityIntent list — **not baseline**; get-by-id is sufficient until a real user queue/history need appears.
- BusinessOrderIntent / InvoicingIntent list — **not baseline**; Sales/Fulfillment/Work provide the operational entry points. Add only if D6 later proves independent materialization-history/queue need without turning the list into a convenience duplicate.

---

# 9. Collection / pagination / filter laws

B2 does not introduce a universal arbitrary query DSL.

## 9.1 Collections that are materially unbounded

Use opaque cursor pagination for collections whose Product 1.0 cardinality/history can grow without a small business bound, including baseline:

- ProductReadiness search/list;
- Offerings;
- Listing/Price Intent histories;
- SellableAvailability population;
- CompetitivePosition / ComparableOffer collections where material;
- SaleEconomics;
- Authorization Decisions;
- Sales;
- Fulfillment work queue;
- Post-Sale Resolutions;
- Operational Work;
- OperationalAttention.

## 9.2 Bounded configuration registries

Marketplace Installations, Selling Entities, Inventory Sources, Organization-baseline policies and product AccessRole definitions may initially return the bounded current set without pagination. Add cursor pagination only when real scale makes the current assumption false.

Membership collection may use cursor pagination because Organization member count is not semantically bounded even though Product 1.0 begins small.

## 9.3 Real baseline filters

Only filters serving named operations are admitted. Baseline semantic filters:

- Readiness — Installation, SourceInstance, readiness state, search text;
- Offering — Installation, lifecycle/convergence state, Product reference, search text;
- Listing/Price Intent — state, target/Offering, Installation, time window when history requires it;
- Availability — Installation, knowledge/convergence/attention state, Product/Offering subject;
- Market Intelligence — Installation, evidence sufficiency/competitive-position state, Product/Offering subject;
- SaleEconomics — Installation, Sale time window, reconciliation/variance state;
- Authorization Decision — pending/final state, owning action domain/class;
- Sales — Installation, business lifecycle state, time window, search text;
- Fulfillment — Fulfillment Node, operational stage, external/internal deadline attention, exception/block state;
- Post-Sale — Sale, Resolution state;
- Work — state, responsibility, assignee, originating owner/type, urgency/deadline attention;
- OperationalAttention — owning domain/category, Installation, urgency/attention class.

Raw provider statuses, ERP TOP/status codes, provider fulfillment-mode codes and technical run/job states are not Product API filters.

## 9.4 Sorting

Each collection owns one stable default sort coherent with its user use and cursor. B2 does not add arbitrary `sort=<field>` support without a real consumer.

Queue/attention semantic ordering may use owner-defined urgency/deadline meaning; B2 does not freeze D6 visual ranking or a generic global priority model.

---

# 10. Bulk / multi-target posture

**No Product API bulk action endpoint is admitted in the B2 baseline.**

D0 multi-target correctness laws do not themselves prove a Product API bulk consumer. Internal policy-driven Availability convergence and provider bulk mechanics may process multiple targets without exposing a generic batch API.

If a real actor workflow later proves that individual operations are materially insufficient, admit a **domain/operation-local** bulk surface and preserve:

- intended target scope;
- authorized scope snapshot;
- member-level confirmed/rejected/ambiguous/not-executed outcomes;
- no blind whole-batch replay.

No `/bulk`, `/mutations`, comma-ID action list or generic Batch identity is created by symmetry.

---

# 11. Idempotency / concurrency posture by operation class

D5-B1 remains authoritative.

### 11.1 Consequential owner capabilities

`Idempotency-Key` is mandatory/fail-closed by default for:

- Listing Intent creation;
- Price Intent creation;
- authorization-decision outcome intake;
- Post-Sale Resolution initiation;
- Fulfillment progression operations able to trigger durable/downstream/external effects;
- any later admitted external-effect capability.

An operation may omit the key only if B2/B3 records a structural owner-semantic duplicate proof.

### 11.2 Durable resource creation

For POST operations creating a canonical MPC identity (e.g. Marketplace Installation, Selling Entity, Inventory Source), the concrete contract must prevent network retry from creating duplicate semantic resources. Prefer:

- natural idempotent addressing when a true semantic key already exists; otherwise
- the same fail-closed idempotency-key mechanism.

Do not invent a business-natural key merely to avoid request idempotency.

### 11.3 Configuration/resource updates

Use optimistic concurrency only where stale overwrite is materially unsafe, especially access/security, policies and lifecycle/configuration with competing editors.

A projection never supplies the concurrency token for owner writes.

### 11.4 No retry endpoint

No generic `retry`, `rerun`, `replay`, `refresh` or `force apply` Product API operation is admitted.

After definitive rejection, a new owner-specific Intent may be created if valid. After possible external acceptance, reconciliation must establish the prior result before another effect is dispatched.

---

# 12. Provider-rich evidence placement

Provider richness remains inside the owner surface that needs it.

Baseline examples:

- Portfolio — provider account/reputation/restriction posture translated to installation attention;
- Readiness — seller SKU/GTIN/provider catalog/attribute relationship evidence;
- Offering — source-qualified category/attribute/listing capability/restriction evidence;
- Availability — provider stock/writability/blast-radius evidence only as needed for current convergence explanation;
- Market Intelligence — `price_to_win`, winner/offer delivered-price/shipping/free-shipping/boost evidence;
- Economics — source-specific expected fee/shipping, Order fee, Payment charge/release/refund decomposition;
- Fulfillment — provider-effective requirements, deadline/readiness/artifact and Shipment observations;
- Post-Sale — provider Claim/Return/Refund/available-action evidence required for Resolution closure.

No top-level ProviderResource/Payment/Refund/Fee/Category/StockLocation/Claim/Return Product API collection is created from these evidence classes.

---

# 13. Legacy OpenAPI semantic disposition

Current routes are evidence only. B2 carries the useful requirement and discards accidental shape.

| Legacy surface family | B2 disposition |
|---|---|
| `/mutations*` | **RETIRE** — Listing/Price owner Intents, Governance, Work and execution safety keep the real meanings; no generic Mutation API |
| `/catalog/products*` | **RETIRE AS PRODUCT MASTER** — source Product descriptors enter ProductReadiness/Economics/Availability contracts as qualified evidence; no MPC Product Master resource |
| `/listings*` read | **REHOME** → Offering Q surface |
| `/listings/refresh` | **RETIRE FROM PRODUCT API** → D7 acquisition/recovery mechanism |
| provider category/attribute endpoints | **REHOME** → Offering `ListingRequirements` provider-enriched Q; no provider taxonomy ontology |
| `/classifications*`, legacy catalog taxonomy | **RETIRE** unless a later accepted owner requirement independently re-admits meaning |
| `/integrations/installations*` | **SPLIT** → Marketplace Installation Portfolio resource + D4 protocol/technical auth surface |
| `/integrations/providers`, fee-sync, operations, probes | **RETIRE FROM PRODUCT API** → D4/D7 descriptor/acquisition/support mechanics; provider support does not become business resource |
| `/product-links/*/imports` / candidate generations | **RETIRE MECHANISM** → Readiness Q candidates + correspondence C; no import/generation command |
| legacy inventory risk/manual-apply | **REHOME MEANING** → SellableAvailability/Work + automatic Availability owner path; no baseline manual stock mutation endpoint |
| `/orders/import` | **RETIRE** → D4/D7 acquisition; Sales commit is owner meaning |
| `/orders` / `/orders/{provider_order_id}` | **REHOME** → source-qualified Marketplace Sales list/get |
| `/orders/summary` | **RETIRE DUPLICATE** → OperationalAttention / Sales query as appropriate |
| `/orders/*/faturado` | **RETIRE** — local marker cannot replace Materialization/Invoicing authoritative semantics |
| `/orders/*/sankhya-linkage*` | **RETIRE PROVIDER SHAPE** → Materialization correlation + Work resolution; no Sankhya Product API noun |
| `/profitability/*/import`, `/calculate` | **RETIRE MECHANISM** → Economics Qs are derived by owner/runtime; no calculation job API |
| `/profitability/manual-adjustments` | **RETIRE GENERIC MANUAL PATCH** → Economics owns attribution; ambiguous cases close via Work/source-owner evidence adjudication |
| `/market/collections` | **RETIRE** → D4/D7 acquisition; client reads Market Intelligence meaning |
| `/market/references`, signals, aggregates, verdicts | **REHOME / CONSOLIDATE** → CompetitivePosition/ComparableOffer owner semantics |
| `/sync/*` | **RETIRE FROM PRODUCT API** → D7 mechanics; owner/OperationalAttention expose honest freshness/coverage/divergence |
| `/erp/imports*` | **RETIRE TARGET PRODUCT SURFACE** → historical acquisition mechanism; target Sankhya source is D4 Gateway |
| `/config/active-source` | **RETIRE FROM PRODUCT BUSINESS SURFACE** → SourceInstance technical setup/binding; not business Product configuration |
| `/config/sellable-assortment` | **REHOME MEANING** → Availability-owned Inventory Scope/Policy configuration |
| marketplace definitions/fee schedules/admin seed | **RETIRE** → no generic provider registry/seed Product API; Economics consumes current provider evidence through D4 |

### Negative legacy test

A legacy route survives only if the same operation is independently re-derived from D0 actor use + accepted owner semantics. Existing frontend/code usage alone is not admission evidence.

---

# 14. Actor coverage proof

### Marketplace Operations Operator

Covered by:

- OperationalAttention;
- Portfolio reads;
- ProductReadiness/correspondence;
- Offering + Listing/Price Intent;
- SellableAvailability;
- Market Intelligence;
- Economics read/simulation;
- Sales;
- Post-Sale;
- Work.

### Fulfillment / Dispatch Operator

Covered by:

- Sales read context;
- Fulfillment queue/detail/progression/readiness/artifacts;
- Materialization Intent reads when referenced;
- shipment/delivery observation in Fulfillment;
- Work.

### Commercial / Marketplace Manager

Covered by:

- OperationalAttention;
- Portfolio/readiness/offering/availability/market/economics reads;
- Availability/Economics/Fulfillment internal-target policy configuration;
- Governance decision operations when separately authorized;
- Work.

### Owner / Administrator / Policy Approver

Covered by:

- Portfolio configuration;
- ordinary access RoleAssignment management;
- Governance delegation/decision operations subject to actual authority;
- policy surfaces where their ordinary Permission/authority applies;
- OperationalAttention.

Provider credential/session mechanics are not forced into these Product API resources; a later D5 technical/protocol surface may support the setup interaction without becoming business ontology.

---

# 15. Explicitly rejected Product API resources/operations

Unless new material evidence reopens B2, do not introduce:

- generic `Mutation`, `Action`, `Command`, `Operation`, `Execution`, `SyncRun` business resource;
- generic `Integration`, `Provider`, `ProviderResource`, `Capability` business resource;
- MPC Product Master / Catalog Product mirror;
- MPC Customer/Party/Address master;
- synthetic MPC Payment/Refund/Settlement alias;
- synthetic MPC Shipment alias merely for provider Shipment;
- provider category/taxonomy as a top-level MPC business resource;
- Sankhya TOP/NUNOTA/Parceiro/Contato/fee schedule/stock location business APIs;
- manual `refresh`, `import`, `sync`, `probe`, `seed`, `calculate`, `retry`, `force apply` Product operations;
- generic bulk command endpoint;
- arbitrary provider-status/filter vocabulary;
- generic policy engine/configuration bag;
- screen-specific BFF endpoints without an independently named Product 1.0 read need.

---

# 16. Essential vs accidental complexity

## Essential preserved

- real actor access/operability;
- explicit Organization scope;
- source-qualified Product/Offering/Sale/native-result identities;
- owner-local Intents;
- policy/config authority by owner;
- ordinary access distinct from Governance/business authorization;
- Q knowledge/freshness/coverage;
- consequential acceptance/ambiguity/convergence;
- fulfillment physical progression and provider closure;
- post-sale scoped resolution;
- Operational Work;
- provider-rich evidence;
- explainable Economics;
- portfolio attention projection.

## Accidental removed

- technical run/job resources in Product API;
- acquisition triggers as user commands;
- provider/business-system resource hierarchy;
- client-orchestrated cross-owner workflow;
- duplicate summary/dashboard authorities;
- generic Mutation/retry mechanics;
- import/calculate/probe/manual-sync APIs;
- speculative Product/Customer/Payment/Shipment mirrors;
- generic provider/policy/bulk frameworks.

---

# 17. Adversarial counterexamples / review questions

Fable should reconstruct authority first and attack at least:

1. **Admission completeness:** Is any D0 completion outcome impossible using only the admitted surface + internal accepted D3 reactions?
2. **Overexposure:** Can any admitted operation be deleted because the same actor outcome is already reachable through another owner surface without losing correctness?
3. **Internal choreography leak:** Do Fulfillment or Materialization operations make the client orchestrate a D3 edge that should remain internal?
4. **Fulfillment semantics:** Are separation/conference/packing/dispatch genuine D1 owner capabilities, or has B2 accidentally created a generic workflow-state machine?
5. **Availability:** Is no public AvailabilityIntent creation sufficient for D0's automatic normal path and operator exception handling?
6. **Post-Sale:** Is initiating PostSaleResolution a real external actor capability, or should all Resolutions originate from observed provider facts/Work?
7. **Access surface:** Does Product 1.0 genuinely require Membership/RoleAssignment API now, or is this premature D2 substrate administration?
8. **Governance Grant addressing:** Can the admitted manage-delegation surface avoid inventing a canonical Grant identity while still being usable/revocable?
9. **OperationalAttention:** Does D0 independently justify this P surface now, or is it D6 BFF leakage? If admitted, does `operations.attention.read` leak data beyond underlying owner permissions?
10. **Permission fragmentation:** Are the proposed 26 Permissions the smallest stable least-privilege catalog, or can some safely combine without forcing role-name checks/business authority into access?
11. **Policy configuration:** Are only Organization-baseline Availability/Economics/Fulfillment policies enough to honor D0's future override seam without implementing a speculative scope DSL?
12. **SourceInstance/integration admin:** Does the Product API need a D2/D4 technical SourceInstance configuration operation for the real Owner/Admin actor now, or is keeping it outside B2 the cleaner B1 boundary?
13. **Listing requirements:** Does the Offering-owned requirements Q preserve provider-rich category/attribute capability without making provider taxonomy Product ontology?
14. **Intent observability:** Are Listing/Price lists necessary while Availability/Materialization lists are correctly absent?
15. **External identity addressing:** Can every exact Q path remain unambiguous without synthetic IDs, especially ProductReadiness, Offering, Sale and fulfillment scope?
16. **No Product master:** Can readiness/search and downstream flows operate without reviving `/products` as an MPC-owned resource?
17. **Economics:** Is one SaleEconomics view preserving L1/L2 distinctions sufficient, or does it accidentally merge rungs that should be independently addressable?
18. **Work:** Does submitting resolution evidence through Work preserve source-owner closure authority, or make Work a generic command gateway?
19. **No bulk:** Does any current Product 1.0 operator workflow materially require a bulk endpoint despite B1/D0 multi-target semantics?
20. **Pagination/filtering:** Are any admitted filters technical/provider vocabulary or speculative UI controls rather than real Product use?
21. **Second marketplace:** Does an Amazon/Mirakl-class provider fit without adding provider resources or changing Product operations?
22. **Second business system:** Does Bling-like realization fit Materialization/Economics/Availability semantics without Sankhya nouns leaking?
23. **Hardest legacy requirement:** Is any useful legacy endpoint meaning lost rather than deliberately rehomed/retired?
24. **D6 leakage:** Does any operation exist only because an imagined screen might need it?
25. **D7 leakage:** Does any endpoint expose scheduling/retry/acquisition/execution mechanism rather than business meaning?

Fable returns `APPROVE`, `REVISE` or `REJECT` with material findings, corrected invariants/proposed surface corrections, proof strategy and reopen triggers only.

---

# 18. Proof strategy before implementation

## P1 — Actor/outcome coverage

Build a ledger:

```text
D0 actor/outcome
→ Product API operation OR accepted internal owner reaction
→ D1/D2 owner
→ Q/C/P
→ ordinary Permission
```

Any D0 outcome with no path is a B2 defect. Any operation with no D0 actor/outcome is presumptively accidental.

## P2 — Owner uniqueness

Every admitted operation has exactly one owner/substrate authority. A route that needs two write authorities is invalid; use P for read composition or separate owner capabilities.

## P3 — Internal-mechanism negative set

Demonstrate that removing Product API access to:

- refresh/sync/import/probe/fee-sync/calculate/seed;
- generic retry;
- provider callback payloads;
- current Sankhya-linkage operations

does not remove a required Product 1.0 actor capability because the meaning is rehomed to an admitted owner surface/internal progression/Work.

## P4 — Source identity collision

Two equal native Product/Listing/Sale IDs under different SourceInstances/Installations must remain distinguishable in every affected operation/SDK contract.

## P5 — Workflow inversion

A client must not need to know:

- provider notification ordering;
- Sales fan-out;
- Materialization↔Fulfillment event choreography;
- Sankhya confirmation/TOP sequence;
- provider label/status sequence;
- Payment acquisition sequence

to perform the normal Product 1.0 flow.

If it does, B2 leaks internal choreography.

## P6 — Permission/business-authority split

Positive/negative fixtures at contract altitude:

- ordinary `offering.operate` permits PriceIntent invocation but Governance/business disposition can still reject/pending;
- `governance.decide` without applicable Governance authority cannot approve;
- `work.operate` cannot declare source truth resolved;
- `operations.attention.read` cannot be used as authority to reach owner write/detail without owner permission.

## P7 — Legacy structural inversion

Delete the current OpenAPI mentally. The B2 operation set must still follow from D0–D5-B1. Then invert the test: add every legacy route back; each must independently satisfy the admission predicate or remain retired.

---

# 19. Reopen / stop triggers

Revisit only on material evidence:

1. a D0 Product 1.0 actor/outcome cannot be completed without an omitted external operation;
2. an admitted operation requires a semantic owner absent from D1/D2;
3. an operation requires a new canonical identity not supported by D2;
4. D3 accepted owner reaction cannot replace an external client orchestration without losing a Product requirement;
5. a provider-required normal-path action cannot fit an accepted owner operation without provider ontology leakage → targeted D4 review;
6. a real role/least-privilege case proves the Permission catalog materially too coarse/fine;
7. a real Product 1.0 workflow proves bulk is materially required;
8. a real policy override scope appears and Organization-baseline configuration is insufficient;
9. a real technical/admin consumer proves SourceInstance/credential setup needs a separately named D5 technical surface;
10. a second provider/business system exposes a genuinely new business meaning rather than different protocol/evidence;
11. D6 later proves a cross-owner read need not representable by owner Qs without unacceptable coupling → admit the smallest P surface, never a write authority;
12. operation-addressing cannot remain source-qualified/unambiguous without a new identity → targeted D2 review.

Current-code convenience, provider symmetry, screen preference and hypothetical future workflows are not reopen evidence.

---

# 20. Candidate outcome

**Proposed outcome:** `RESTRUCTURE NOW` for the legacy Product API surface; **CURRENT D0–D5-B1 STRUCTURE CONFIRMED** for owner/product semantics.

If B2 survives independent review and operator ratification, canonical D5 should establish:

- actor/use-driven Product API admission;
- the semantic resource/operation families in §6;
- ordinary Permission→operation mapping in §7/§8;
- cursor/filter posture in §9;
- no baseline Product API bulk action;
- no public client orchestration of internal D3/D4/D7 work;
- no generic import/refresh/sync/probe/calculate/retry/mutation APIs;
- no Product/Customer/Payment/Shipment/provider mirror ontology;
- source-qualified external identities and provider-rich owner evidence;
- explicit retirement/rehoming of the legacy operation families in §13;
- no D0–D4 or D5-B1 reopen unless independent review finds a material contradiction.

This candidate deliberately leaves final per-operation request/response schemas, status mappings, relationship addressing details and OpenAPI tooling to later D5 work after the operation surface is accepted.

**This file remains NON-AUTHORITATIVE until independent Fable review, GPT adjudication and explicit operator ratification. Implementation remains blocked until D9.**
