# D5-B2 — Operation Admission Matrix

> **Status:** OPEN / ACTIVE — Blocks 1–4 ACCEPTED IN-STAGE; Block 5 next  
> **Parent B2:** `D5-B2-PRODUCT-OPERATION-SURFACE.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + Decision Reconciliation Baseline + B2-A  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-18

## 1. Governing admission rule

The Product API contains the smallest operation surface real Product 1.0 clients need before final path/schema spelling.

Every admitted operation must identify proportionately:

- real consumer/use;
- allowed client class: human, machine/automation/system, or both;
- exactly one accepted semantic owner or D2 substrate authority;
- Q / C / P class;
- explicit Organization scope;
- ordinary Permission, distinct from business disposition/Governance;
- canonical/source-qualified subject identity;
- honest knowledge/freshness/provenance when read semantics require it;
- consequential Intent/outcome/idempotency/precondition/concurrency where applicable;
- provider enrichment only for a named consumer/correctness property;
- pagination/filter/sort only for a real collection consumer;
- bulk only when a real workflow requires member-level semantics.

Reject a candidate whose only justification is API symmetry, current code, provider endpoint shape, debug convenience or hypothetical future need.

---

# 2. Block 1 — Identity/Access + Portfolio + Readiness — ACCEPTED IN-STAGE

## 2.1 D2 identity / ordinary access substrate

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `GetCurrentAccessContext` | both | Q | authenticated | **ADMIT** |
| `ListOrganizationMembers` | human | Q | `access.read` | **ADMIT** |
| `ListAccessRoles` | human | Q | `access.read` | **ADMIT** |
| `AssignAccessRole` | human | C | `access.manage` | **ADMIT** |
| `RevokeAccessRole` | human | C | `access.manage` | **ADMIT** |
| invitation/onboarding by email | human | C | — | **DEFER** |
| create/delete Organization | human | C | — | **DEFER** |
| create Keycloak/OAuth client/service account | human | — | — | **REJECT FROM PRODUCT API** |
| issue MPC API key/token | both | — | — | **REJECT** |

Binding decisions:

- one D2-owned access-context Q is preferred over fragmented `/me`/roles/permissions/session surfaces;
- IdP authentication does not carry MPC Organization/Permission authority by implication;
- role assignment/revocation operates on one explicit Membership + AccessRole relation rather than replacing whole role sets;
- repeated identical assignment/revocation may use structural idempotency where uniqueness makes duplication harmless;
- no custom-role designer, generic ACL/ReBAC engine or IdP-role business authority is admitted.

Permission floor: `access.read`, `access.manage`.

## 2.2 Marketplace Portfolio

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListMarketplaceInstallations` | both | Q | `portfolio.read` | **ADMIT** |
| `GetMarketplaceInstallation` | both | Q | `portfolio.read` | **ADMIT** |
| `CreateMarketplaceInstallation` | human | C | `portfolio.manage` | **ADMIT** |
| `UpdateMarketplaceInstallationConfiguration` | human | C | `portfolio.manage` | **ADMIT** |
| `DeactivateMarketplaceInstallation` | human | C | `portfolio.manage` | **ADMIT** |
| `ListSellingEntities` | both | Q | `portfolio.read` | **ADMIT** |
| arbitrary Selling Entity create/edit | human | C | — | **DEFER** |
| provider OAuth begin/callback/refresh | human/provider | — | — | **NOT PRODUCT API — D4** |
| generic provider/integration catalog | both | Q | — | **REJECT** |

Binding decisions:

- `CreateMarketplaceInstallation` requires client idempotency by default because a lost response could mint duplicate MPC identities;
- materially unsafe configuration updates require an opaque MPC concurrency/precondition token;
- deactivation preserves identity/history and may be structurally idempotent;
- provider OAuth remains D4 protocol, never Portfolio business meaning.

Permission floor: `portfolio.read`, `portfolio.manage`.

## 2.3 Product & Channel Readiness

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `SearchSourceProductsForMarketplace` | both | Q | `readiness.read` | **ADMIT** |
| `GetProductChannelReadiness` | both | Q | `readiness.read` | **ADMIT** |
| `GetPublicationRequirements` | both | Q | `readiness.read` | **ADMIT** |
| `ResolveProductChannelCorrespondence` | both | C | `readiness.manage` | **ADMIT** |
| `ClearProductChannelCorrespondence` | both | C | `readiness.manage` | **ADMIT — SAME LIFECYCLE** |
| `RefreshReadiness` / `RecalculateNow` | both | C | — | **REJECT PRODUCT API** |
| `SyncProducts` | machine | C | — | **REJECT PRODUCT API** |
| `CreateProduct` / `UpdateProduct` | both | C | — | **REJECT** |
| bulk correspondence mutation | both | C | — | **DEFER** |

Binding decisions:

- Product remains `SourceInstance + native key`; no MPC Product resource/master;
- source-product search is admitted only in marketplace-operating Readiness context and returns source-qualified evidence with honest provenance/freshness;
- readiness and publication-requirements are distinct Q meanings;
- correspondence replacement/clearing requires current-state concurrency when stale overwrite could defeat a newer decision;
- automation never silently supersedes standing human correspondence from stale state;
- exact duplicate correspondence may use structural idempotency when harmless.

Permission floor: `readiness.read`, `readiness.manage`.

**Block 1 outcome:** `CURRENT PARENT STRUCTURE CONFIRMED; ADMIT ONLY CONSUMER-PROVEN OPERATIONS`.

---

# 3. Block 2 — Offering + Price + Availability — ACCEPTED IN-STAGE

## 3.1 Governing invariant

> **Marketplace Listing observation, Offering-owned authoring/price intent, Availability-owned sellable quantity/configuration and provider execution remain distinct meanings. A provider request may physically compose owner-issued inputs, but Product API contracts never collapse those authorities into one giant Listing mutation or generic asynchronous Operation.**

## 3.2 Marketplace Listing observation

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListMarketplaceListings` | both | Q | `offering.read` | **ADMIT** |
| `GetMarketplaceListing` | both | Q | `offering.read` | **ADMIT** |
| direct `CreateMarketplaceListing` | both | C | — | **REJECT** |
| direct `UpdateMarketplaceListing` | both | C | — | **REJECT** |
| direct `Pause/Reactivate/CloseListing` APIs | both | C | — | **REJECT AS PARALLEL BASELINE AUTHORITY** |

Listing identity remains Marketplace Installation + provider-native Listing key. Collection pagination/filtering is justified by real population use. Read contracts expose Offering interpretation + source freshness/convergence, never raw provider DTOs.

## 3.3 ListingIntent

| Candidate operation | Client | Class | Permission | Idempotency | Concurrency | Admission |
|---|---|---|---|---|---|---|
| `ListListingIntents` | both | Q | `offering.read` | — | — | **ADMIT** |
| `GetListingIntent` | both | Q | `offering.read` | — | — | **ADMIT** |
| `CreateListingIntentDraft` | both | C/create | `listing.manage` | **mandatory** | — | **ADMIT** |
| `UpdateListingIntentDraft` | both | C/update | `listing.manage` | structural | **required** | **ADMIT** |
| `DiscardListingIntentDraft` | both | C | `listing.manage` | structural | current state/version | **ADMIT** |
| `SubmitListingIntent` | both | C | `listing.manage` | structural by Intent identity/revision | **required** | **ADMIT** |
| separate `ValidateListingIntent` | both | Q/C | — | — | — | **REJECT BASELINE** |
| separate `PreviewListing` | both | Q | — | — | — | **REJECT BASELINE** |

Binding decisions:

- create/edit use one ListingIntent identity (`target = none | existing Listing`);
- draft updates are declarative desired-state changes, not a mutation DSL;
- only `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` baseline value modes are admitted;
- stale automation cannot overwrite or submit against a newer human decision;
- `SubmitListingIntent` is the consequential client boundary; technical freeze/authorize/execute/reconcile are not separate client operations;
- `submitted != authorized != externally applied != converged` remains binding;
- no generic LongRunningOperation is admitted because owner-local Intents already provide durable tracking identity.

## 3.4 PriceIntent

| Candidate operation | Client | Class | Permission | Idempotency | Admission |
|---|---|---|---|---|---|
| `ListPriceIntents` | both | Q | `offering.read` | — | **ADMIT** |
| `GetPriceIntent` | both | Q | `offering.read` | — | **ADMIT** |
| `CreatePriceIntent` | both | C | `price.manage` | **mandatory** | **ADMIT** |
| direct `SetListingPrice` | both | C | — | — | **REJECT** |
| public mutable PriceDraft | both | C | — | — | **REJECT BASELINE** |
| withdraw/cancel pending PriceIntent | both | C | — | — | **DEFER** |

`PriceIntent` expresses desired exact Money target. Economics owns reasoning/simulation; Offering owns actuation.

## 3.5 Availability

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListSellableAvailability` | both | Q | `availability.read` | **ADMIT** |
| `GetSellableAvailability` | both | Q | `availability.read` | **ADMIT** |
| `ListInventorySources` | both | Q | `availability.read` | **ADMIT** |
| `GetInventorySource` | both | Q | `availability.read` | **ADMIT** |
| `CreateInventorySource` | human | C | `availability.manage` | **ADMIT** |
| `UpdateInventorySource` | human | C | `availability.manage` | **ADMIT** |
| `DeactivateInventorySource` | human | C | `availability.manage` | **ADMIT** |
| get effective allocation/scope policy | both | Q | `availability.read` | **ADMIT** |
| update allocation/scope policy | human | C | `availability.manage` | **ADMIT** |
| public `CreateAvailabilityIntent` | both | C | — | **REJECT BASELINE** |
| `SetAvailableQuantity` | both | C | — | **REJECT** |
| `SyncAvailability` / `RefreshAvailability` | machine/both | C | — | **REJECT PRODUCT API** |

Availability Q preserves desired value, provider actual evidence, knowledge/freshness/provenance and owner-specific convergence. Inventory Source and allocation/scope policy remain distinct meanings.

## 3.6 Joint Offering × Availability realization

D4-R1 R1-G1 remains binding. `SubmitListingIntent` never accepts quantity as Offering-owned content. D4/D7 may jointly serialize owner-issued Offering + Availability meanings; missing required Availability input fails closed before dispatch; no cross-owner atomicity is claimed; provider blast radius never silently widens business intent.

Permission floor: `offering.read`, `listing.manage`, `price.manage`, `availability.read`, `availability.manage`.

**Block 2 outcome:** `CURRENT D1/D2/D3/D4-R1 STRUCTURE CONFIRMED; ADMIT OWNER-SPECIFIC RESOURCES/INTENTS, REJECT GIANT LISTING CRUD AND GENERIC ASYNC OPERATION`.

---

# 4. Block 3 — Market Intelligence + Commercial Economics — ACCEPTED IN-STAGE

## 4.1 Governing invariant

> **Market Intelligence owns competitive interpretation/evidence sufficiency; Commercial Economics owns economic meaning, simulation and L0/L1/L2 lineage; Offering alone owns PriceIntent. Hypothetical exploration does not gain durable identity by default, while economic evidence/decisions that materially participate in later explanation or reconciliation remain durable under Economics.**

## 4.2 Market Intelligence

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListCompetitivePositions` | both | Q | `market.read` | **ADMIT** |
| `GetCompetitivePosition` | both | Q | `market.read` | **ADMIT** |
| `ListComparableOffers` | both | Q | `market.read` | **ADMIT** |
| generic `Get/CreateMarketObservation` | both | Q/C | — | **REJECT** |
| `CollectMarketNow` / `RefreshCompetitivePosition` | both | C | — | **REJECT PRODUCT API** |
| manual `SetComparableOffer` | both | C | — | **DEFER** |
| generic scraper/collector configuration | human | C | — | **REJECT** |

Provider-rich evidence remains source-qualified/optional and enters only for explanation/correctness. Source/provider population never implies universal market completeness.

Permission floor: `market.read`.

## 4.3 Expected Economics and scenario evaluation

| Candidate operation | Client | Class | Permission | Idempotency | Admission |
|---|---|---|---|---|---|
| `ListExpectedEconomics` | both | Q | `economics.read` | — | **ADMIT** |
| `GetExpectedEconomics` | both | Q | `economics.read` | — | **ADMIT** |
| `EvaluatePriceScenario` | both | C — stateless/side-effect-free | `economics.read` | not required | **ADMIT** |
| persistent `Simulation` resource | both | C | — | — | **REJECT BASELINE** |
| persistent `Recommendation` resource | both | C | — | — | **REJECT BASELINE** |
| generic mutable `Profitability` blob | both | Q/C | — | — | **REJECT** |

Scenario input may state legitimate hypothetical variables but cannot impersonate authoritative fee/tax/cost/market evidence. No SimulationID exists merely for history. When a consequential PriceIntent is later created, Economics preserves/re-establishes the material decision-time L0 basis required for future explanation without trusting client-computed evidence.

## 4.4 Sale Economics and lineage

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListSaleEconomics` | both | Q | `economics.read` | **ADMIT** |
| `GetSaleEconomics` | both | Q | `economics.read` | **ADMIT** |
| `GetEconomicPerformanceSummary` | both | Q | `economics.read` | **ADMIT** |
| separate mutable Order/Realized Economics CRUD | both | Q/C | — | **REJECT** |
| universal `Reconciliation` resource | both | Q/C | — | **REJECT** |

`SaleEconomics` preserves distinct L0 expected basis, L1 order economics, L2 realized/settlement evidence and R1/R2 reconciliation. Later refunds/reversals append/reinterpret lineage rather than erasing prior material occurrences. Period summaries preserve explicit scope/coverage/partiality. R3 bank-side reconciliation remains deferred.

## 4.5 Commercial policy and Economic Attribution

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `GetCommercialPolicy` | both | Q | `economics.read` | **ADMIT** |
| `UpdateCommercialPolicy` | human | C/update | `economics.policy.manage` | **ADMIT** |
| generic Policy CRUD / rules DSL | both/human | C | — | **REJECT** |
| `ListEconomicAttributions` | both | Q | `economics.read` | **ADMIT** |
| `GetEconomicAttribution` | both | Q | `economics.read` | **ADMIT** |
| `ResolveEconomicAttribution` | human baseline | C | `economics.reconcile` | **ADMIT** |
| generic `MarkReconciled` | both | C | — | **REJECT** |
| `ReconcileNow` | both | C | — | **REJECT PRODUCT API** |

Commercial policy remains Economics-owned; Governance does not acquire business thresholds. Policy updates use concurrency where lost-update is material. Economic Attribution is persistent Economics state; explicit ambiguity resolution is human baseline and requires current-state concurrency where necessary.

Permission floor: `economics.read`, `economics.policy.manage`, `economics.reconcile`.

**Block 3 outcome:** `CURRENT D1/D2/D4-B4 STRUCTURE CONFIRMED; USE STATELESS HYPOTHETICAL ANALYSIS + DURABLE MATERIAL ECONOMIC LINEAGE, KEEP PRICE ACTUATION IN OFFERING`.

---

# 5. Block 4 — Controlled Action Governance + Marketplace Sales + Business-System Materialization — ACCEPTED IN-STAGE

## 5.1 Governing invariant

> **Governance decides consequential authorization, Sales owns marketplace-sale meaning, and Materialization owns business-system intents/results. Client access, approval, upstream domain facts and runtime execution are distinct. A durable domain Intent may be publicly readable without being client-created when its legitimate cause is an accepted owner reaction.**

Rejected local maxima:

```text
POST /actions
PATCH /workflow/status
POST /erp/orders
POST /invoice
POST /retry
```

Selected structure:

```text
action-owner Intent
  → Governance AuthorizationDecision when required
  → owner revalidation/execution

provider sale evidence
  → Marketplace Sales commits Sale meaning
  → E to Materialization/Fulfillment/Economics

SaleCommitted
  → Materialization creates/advances BusinessOrderIntent

Fulfillment physical-readiness checkpoint
  → Materialization creates/advances InvoicingIntent
```

## 5.2 Controlled Action Governance

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListAuthorizationDecisions` | both | Q | `governance.read` | **ADMIT** |
| `GetAuthorizationDecision` | both | Q | `governance.read` | **ADMIT** |
| `CreateAuthorizationDecision` | human baseline | C/create | `governance.decide` | **ADMIT** |
| `ListAuthorizationDelegations` | human | Q | `governance.manage` | **ADMIT** |
| set/update bounded Authorization Delegation | human | C | `governance.manage` | **ADMIT** |
| revoke Authorization Delegation | human | C | `governance.manage` | **ADMIT** |
| separate `ListPendingApprovals` queue authority | human | Q | — | **REJECT BASELINE** |
| domain-specific `ApproveListingIntent` / `ApprovePriceIntent` | human | C | — | **REJECT** |
| `ExecuteApprovedAction` | both | C | — | **REJECT** |
| generic approval/policy/rules engine | both | C | — | **REJECT** |

Binding decisions:

- Authorization Decision is a durable Governance-owned occurrence, not an `approved=true` field or Intent-status mutation;
- decision targets one concrete owner Intent/material revision/context and preserves decision Principal, authority context, outcome and exact authorized target-scope snapshot;
- authorized scope is constrained by intended scope and never widens it;
- `CreateAuthorizationDecision` requires client idempotency because a lost response must not create duplicate decision occurrences;
- current Intent revision/context is a material precondition; stale approval is rejected rather than silently applied to newer meaning;
- approval never executes an effect, mutates owner Intent, waives business validity or becomes eternal authorization; the action owner revalidates at execution time;
- approval queue/work responsibility is not a second Governance queue authority; Work/projection handles actionable attention later;
- D2 intentionally did not freeze physical Grant/Delegation identity/cardinality, so B2 admits bounded list/set/update/revoke semantics without inventing a mandatory universal `GrantID` model prematurely.

Permission floor: `governance.read`, `governance.decide`, `governance.manage`.

## 5.3 Marketplace Sales

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListMarketplaceSales` | both | Q | `sales.read` | **ADMIT** |
| `GetMarketplaceSale` | both | Q | `sales.read` | **ADMIT** |
| `ResolveSaleSellingEntityAttribution` | human baseline | C | `sales.manage` | **ADMIT** |
| `CreateSale` / generic `UpdateSale` | both | C | — | **REJECT** |
| `CancelSale` / `RefundSale` | both | C | — | **REJECT — POST-SALE CONCERN** |
| `SyncSales` / `RefreshOrders` | both | C | — | **REJECT PRODUCT API** |
| provider Pack/Shipment absorbed as Sales entities | both | Q | — | **REJECT** |

Binding decisions:

- Sale identity remains Marketplace Installation + provider-native Sale/Order key; no synthetic MPC Sale alias exists merely for normalization;
- Sale reads preserve source observation/freshness and Sales-owned interpretation/context/correlation/transaction-specific Selling Entity attribution without absorbing downstream Materialization/Fulfillment/Economics/Post-Sale state as mutable ownership;
- Sales is externally originated; Product clients do not create/update provider sales;
- transaction-specific Selling Entity attribution may normally be established automatically by Sales, but a genuinely ambiguous case admits explicit human resolution;
- ambiguous attribution resolution requires current-state concurrency when stale overwrite could replace a newer decision; exact repeat may be structurally idempotent.

Permission floor: `sales.read`, `sales.manage`.

## 5.4 Business Order Intent — owner-triggered, client-readable

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListBusinessOrderIntents` | both | Q | `materialization.read` | **ADMIT** |
| `GetBusinessOrderIntent` | both | Q | `materialization.read` | **ADMIT** |
| `CreateBusinessOrderIntent` | both | C | — | **REJECT BASELINE** |
| `RetryBusinessOrderMaterialization` | both | C | — | **REJECT** |
| direct `CreateSankhyaOrder` / `ConfirmSankhyaOrder` | both | C | — | **REJECT** |

Normal-path BusinessOrderIntent creation belongs to Materialization reacting to committed Sales meaning, not to React/agent commands. Duplicate/replayed Sale occurrences must not create duplicate semantic BusinessOrderIntents; D7 chooses the enforcement mechanism for this owner-level idempotency property.

BusinessOrderIntent Q may expose owner meaning/prerequisites, Party/Destination resolution state, source-qualified native result references, accepted/rejected/pending/ambiguous outcome and owner-specific convergence without exposing TOP/NUNOTA/CACSP/status choreography as Product semantics.

## 5.5 Party Resolution and Destination Realization

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `GetBusinessSystemPartyResolution` | both | Q | `materialization.read` | **ADMIT** |
| `ResolveBusinessSystemPartyResolution` | human baseline | C | `materialization.resolve` | **ADMIT** |
| `GetDestinationRealization` | both | Q | `materialization.read` | **ADMIT** |
| `ResolveDestinationRealization` | human | C | `materialization.resolve` | **DEFER / CONDITIONED ON D8 PROOF** |
| generic Customer/Party CRUD | both | C | — | **REJECT** |
| generic Address/Contact CRUD | both | C | — | **REJECT** |
| raw `SelectCODPARC` / mutate Partner address | human/both | C | — | **REJECT PROVIDER VOCABULARY / UNSAFE MASTER MUTATION** |

Binding decisions:

- Party Resolution is a bounded Materialization prerequisite, not Customer/CRM authority;
- a human resolution chooses a compatible source-qualified native party meaning, never a raw provider code as MPC business ontology;
- zero-match native creation may occur only under owner rules and sufficiently known legitimate transaction evidence; the client does not author arbitrary Customer master data;
- resolution targets one existing PartyResolution/current version; structural idempotency may be used only if later proof shows duplicate consequential native creation cannot occur, otherwise B1's idempotency-key default applies;
- Destination Realization remains distinct from Party Resolution and never overwrites registered/master address by convenience;
- contact-based alternate destination is the strongest current Sankhya candidate but its write capability remains conditioned on D8 controlled proof, so B2 does not claim a Product C operation before that evidence exists;
- no safe destination realization remains explicit `external-required` / Work rather than fabricated equivalence.

## 5.6 Invoicing Intent — owner-triggered from Fulfillment readiness

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListInvoicingIntents` | both | Q | `materialization.read` | **ADMIT** |
| `GetInvoicingIntent` | both | Q | `materialization.read` | **ADMIT** |
| `CreateInvoicingIntent` | both | C | — | **REJECT BASELINE** |
| `RequestInvoice` / direct `SankhyaFaturar` | human/both | C | — | **REJECT** |
| `RetryInvoice` | both | C | — | **REJECT** |

Fulfillment owns physical readiness/conference. Its committed checkpoint makes Materialization eligible to create/block/advance InvoicingIntent; Product clients do not create the fiscal intent directly.

Materialization must revalidate current business-order state, Party/Destination prerequisites, Fulfillment readiness, applicable Governance, source binding/capability and execution-time validity before any irreversible fiscal effect. A prior event/approval is never sufficient by itself.

## 5.7 No blind retry / no workflow command surface

A provider timeout or connection loss after possible Sankhya acceptance is reconciled by authoritative reread/correlation before any further effect decision. The Product API does not expose generic `RetryBusinessOrder`, `RetryInvoice`, `AdvanceWorkflow`, TOP progression or provider-command replay.

If reconciliation later establishes a new semantically safe action, the owning domain progresses its existing intent or creates a new owner-local intent according to its lifecycle; it never blindly replays an ambiguous external request.

## 5.8 Permission floor

- `governance.read`
- `governance.decide`
- `governance.manage`
- `sales.read`
- `sales.manage`
- `materialization.read`
- `materialization.resolve`

The split preserves materially useful least privilege without creating one permission per endpoint.

## 5.9 Explicit Block 4 exclusions

Not admitted:

- generic `/actions`, workflow engine or `PATCH status` orchestration;
- domain-specific approve endpoints duplicating Governance;
- approval queue as a second Work authority;
- `ExecuteApprovedAction`;
- client-created Sale;
- Sales-owned cancellation/refund;
- generic Sales sync/refresh;
- client-created BusinessOrderIntent/InvoicingIntent on normal path;
- direct Sankhya order/confirmation/invoice/TOP/NUNOTA APIs;
- Customer/Address master CRUD;
- unsafe Partner-address mutation;
- generic Retry/replay operations after ambiguous external effects.

## 5.10 Block 4 method outcome

**Parent structure:** `CURRENT D1/D2/D3/D4-B3 STRUCTURE CONFIRMED`.

**B2 outcome:**

> **Expose Governance decisions/delegation and owner-tracking state; keep externally-originated Sales read-centric except explicit attribution resolution; let Materialization create BusinessOrder/Invoicing intents from accepted upstream facts rather than client commands; expose human resolution only where ambiguity is a genuine business decision; never leak Sankhya choreography or generic retry/workflow commands into Product API.**

---

# 6. Exact next matrix work

**Block 5 — Fulfillment Lifecycle + Post-Sale Resolution + Operational Work + justified read-only P compositions.**

Derive the smallest Product API surface while preserving these fences:

- Fulfillment owns physical readiness/execution, Fulfillment Node eligibility/selection, separation/conference/packing/dispatch and provider-requirement closure for claimed fulfillment paths;
- Shipment remains source-qualified external identity and does not collapse into Sale/Order;
- provider-native requirements/artifacts remain D4 evidence; Fulfillment owns business closure/readiness meaning;
- physical conference/readiness is the legitimate owner signal that can enable Materialization invoicing; clients must not bypass it with a direct invoice command;
- Post-Sale Resolution owns coordinated cancellation/return/refund consequence scope/closure without becoming provider Claim/Return ontology, CRM or generic reverse logistics;
- one Sale may have 0..N scoped Post-Sale Resolutions; cancellation/return/refund are not one mutually exclusive status;
- Operational Work owns responsibility/assignment/escalation/work state, never source truth or domain closure;
- no-ownerless-work and source-domain closure evaluation remain binding;
- read-only P compositions may summarize several owners for operator attention but never become write/concurrency/business truth authority;
- no generic workflow/Task/Case/OperationalStage mutation surface;
- no direct provider shipment/claim/refund protocol API;
- no generic retry/sync/refresh commands by convenience;
- pagination/filter/sort only where queues/populations create a real consumer need;
- bulk only if a concrete fulfillment/work workflow proves member-level semantics are necessary.

The block must decide which physical fulfillment checkpoints are true client capabilities, how Work mutations are scoped without stealing source-domain closure, how Post-Sale is initiated/coordinated, and which cross-owner cockpit/attention reads are justified P rather than D6-only view composition.

Do not spell final HTTP paths/schemas until Block 5 admission is coherent.

Implementation remains blocked until D9.