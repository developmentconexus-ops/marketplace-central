# D5-B2 — Operation Admission Matrix

> **Status:** OPEN / ACTIVE — Blocks 1–5 ACCEPTED IN-STAGE; Whole-Matrix Global Coherence review next  
> **Parent B2:** `D5-B2-PRODUCT-OPERATION-SURFACE.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + Decision Reconciliation Baseline + B2-A  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-18  
> **Block 5 accepted:** 2026-08-18

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

Listing identity remains Marketplace Installation + provider-native Listing key. Read contracts expose Offering interpretation + source freshness/convergence, never raw provider DTOs.

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
- authorized scope never widens intended scope;
- `CreateAuthorizationDecision` requires idempotency and current Intent revision/context preconditions;
- approval never executes, mutates the owner Intent, waives domain validity or becomes eternal authorization;
- actionable approval responsibility belongs to Work/projection, not a second Governance queue authority;
- bounded delegation semantics are admitted without prematurely forcing a universal physical Grant identity/cardinality.

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

Sale identity remains Marketplace Installation + provider-native Sale/Order key. Sales is externally originated and read-centric; only genuine Sales-owned ambiguity such as transaction-specific Selling Entity attribution admits explicit human resolution.

Permission floor: `sales.read`, `sales.manage`.

## 5.4 Business Order / Party / Destination

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListBusinessOrderIntents` | both | Q | `materialization.read` | **ADMIT** |
| `GetBusinessOrderIntent` | both | Q | `materialization.read` | **ADMIT** |
| `CreateBusinessOrderIntent` | both | C | — | **REJECT BASELINE** |
| `RetryBusinessOrderMaterialization` | both | C | — | **REJECT** |
| direct `CreateSankhyaOrder` / `ConfirmSankhyaOrder` | both | C | — | **REJECT** |
| `GetBusinessSystemPartyResolution` | both | Q | `materialization.read` | **ADMIT** |
| `ResolveBusinessSystemPartyResolution` | human baseline | C | `materialization.resolve` | **ADMIT** |
| `GetDestinationRealization` | both | Q | `materialization.read` | **ADMIT** |
| `ResolveDestinationRealization` | human | C | `materialization.resolve` | **DEFER / CONDITIONED ON D8 PROOF** |
| generic Customer/Party/Address/Contact CRUD | both | C | — | **REJECT** |

Normal-path BusinessOrderIntent creation belongs to Materialization reacting to committed Sales meaning, not to React/agent commands. Duplicate/replayed Sale occurrences must not create duplicate semantic BusinessOrderIntents; D7 chooses enforcement. Party Resolution is bounded Materialization correctness state, not CRM mastery. Destination realization remains distinct and write capability is not claimed until D8 proves the safe selected Sankhya lane.

## 5.5 Invoicing Intent

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListInvoicingIntents` | both | Q | `materialization.read` | **ADMIT** |
| `GetInvoicingIntent` | both | Q | `materialization.read` | **ADMIT** |
| `CreateInvoicingIntent` | both | C | — | **REJECT BASELINE** |
| `RequestInvoice` / direct `SankhyaFaturar` | human/both | C | — | **REJECT** |
| `RetryInvoice` | both | C | — | **REJECT** |

Fulfillment owns physical readiness/conference. Its committed checkpoint makes Materialization eligible to create/block/advance InvoicingIntent. Product clients do not create the fiscal intent directly. Materialization revalidates all current prerequisites and no ambiguous provider effect is blindly retried.

Permission floor: `materialization.read`, `materialization.resolve`.

**Block 4 outcome:** `CURRENT D1/D2/D3/D4-B3 STRUCTURE CONFIRMED; EXPOSE GOVERNANCE/SALES/MATERIALIZATION MEANING WITHOUT CLIENT-COMMANDED OWNER REACTIONS OR SANKHYA CHOREOGRAPHY`.

---

# 6. Block 5 — Fulfillment + Post-Sale + Operational Work + P compositions — ACCEPTED IN-STAGE

## 6.1 Governing invariant

> **Fulfillment owns physical execution and provider-readiness closure, Post-Sale owns scoped consequence coordination, Work owns responsibility/assignment/escalation, and P compositions are read-only convenience views. Physical facts, consequence closure, actionable-work lifecycle and projection convenience never collapse into one generic OrderWorkflow/Task/Case/status authority.**

## 6.2 Fulfillment Lifecycle

| Candidate operation | Client | Class | Permission | Idempotency / concurrency | Admission |
|---|---|---|---|---|---|
| `ListFulfillmentStates` | both | Q | `fulfillment.read` | — | **ADMIT** |
| `GetFulfillmentState` | both | Q | `fulfillment.read` | — | **ADMIT** |
| `ListFulfillmentNodes` | both | Q | `fulfillment.read` | — | **ADMIT** |
| `GetFulfillmentNode` | both | Q | `fulfillment.read` | — | **ADMIT** |
| `CreateFulfillmentNode` | human | C | `fulfillment.manage` | idempotency by default | **ADMIT** |
| `UpdateFulfillmentNode` | human | C | `fulfillment.manage` | concurrency where lost-update material | **ADMIT** |
| `DeactivateFulfillmentNode` | human | C | `fulfillment.manage` | structural + current state | **ADMIT** |
| manual `SelectFulfillmentNode` | human | C | — | — | **DEFER** |
| `RecordSeparation` | human baseline | C | `fulfillment.execute` | idempotency + current-state precondition by default | **ADMIT** |
| `RecordPhysicalConference` | human or explicitly proven physical-system Principal | C | `fulfillment.execute` | idempotency + current-state precondition by default | **ADMIT** |
| `RecordPacking` | human baseline | C | `fulfillment.execute` | idempotency + current-state precondition by default | **ADMIT** |
| `RecordDispatchHandoff` | human or explicitly proven physical-system Principal | C | `fulfillment.execute` | idempotency + current-state precondition by default | **ADMIT** |
| generic `AdvanceFulfillmentStatus` | both | C | — | — | **REJECT** |

Binding decisions:

- physical checkpoints are owner-specific facts/capabilities, not `PATCH status` workflow steps;
- `RecordPhysicalConference` is the material Fulfillment fact that may awaken Materialization for invoicing under D3; the client never bypasses it with a direct invoice command;
- an ordinary automation token cannot fabricate physical evidence merely because it has API access; machine establishment of physical facts requires a separately proven system Principal/source capable of establishing that fact;
- repeat/lost-response safety must prevent duplicate material occurrences and duplicate downstream reactions;
- no company-wide WMS/TMS model is introduced.

## 6.3 Provider requirements and Fulfillment artifacts

`GetFulfillmentState` may expose Fulfillment-owned closure/readiness plus source-qualified provider requirement/deadline evidence sufficient for the claimed path. Provider requirement/artifact truth remains D4/external authority.

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListFulfillmentArtifacts` | human/both read | Q | `fulfillment.execute` | **ADMIT** |
| `GetFulfillmentArtifact` | human/both read | Q | `fulfillment.execute` | **ADMIT** |
| generic `CreateArtifact` | both | C | — | **REJECT** |
| generic provider report/artifact generation | both | C | — | **REJECT BASELINE** |

Artifacts such as labels/handoff documents remain Fulfillment-local/source-qualified and PII-minimized. A future provider operation that must actively generate an artifact is consequential and requires its own owner/effect adjudication; it is never disguised as a read.

## 6.4 Shipment / delivery observation

Shipment remains Marketplace Installation + provider-native Shipment key and never collapses into Sale/Order.

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListShipments` | both | Q | `fulfillment.read` | **ADMIT** |
| `GetShipment` | both | Q | `fulfillment.read` | **ADMIT** |
| `RefreshShipment` | both | C | — | **REJECT PRODUCT API** |
| generic `UpdateShipmentStatus` | both | C | — | **REJECT** |
| baseline `MarkDelivered` | both | C | — | **REJECT SELECTED LANE** |
| provider Shipment/status DTO mirror | both | Q/C | — | **REJECT** |

Reads preserve interpreted Shipment meaning, relevant deadline/SLA/delivery outcome and source freshness/provenance without importing provider status/substatus ontology wholesale.

## 6.5 Post-Sale Resolution

| Candidate operation | Client | Class | Permission | Idempotency / concurrency | Admission |
|---|---|---|---|---|---|
| `ListPostSaleResolutions` | both | Q | `post_sale.read` | — | **ADMIT** |
| `GetPostSaleResolution` | both | Q | `post_sale.read` | — | **ADMIT** |
| `CreatePostSaleResolution` | both | C/create | `post_sale.manage` | **mandatory idempotency** | **ADMIT** |
| generic `UpdateResolutionStatus` | both | C | — | — | **REJECT** |
| direct `ClosePostSaleResolution` | both | C | — | — | **REJECT** |
| direct `CancelSale` / `RefundSale` | both | C | — | — | **REJECT** |
| provider-shaped `AcceptClaim` / `ReviewReturn` | both | C | — | — | **REJECT BASELINE** |
| generic `SubmitPostSaleAction` | both | C | — | — | **REJECT BASELINE** |
| concrete later post-sale decisions | human/both | C | — | — | **DEFER UNTIL SELECTED ACTION IS PROVEN** |

Binding decisions:

- a Resolution may be provider-originated or explicitly MPC-initiated; Product creation creates the scoped coordination obligation, not all underlying cancel/return/refund/ERP/physical effects;
- one Sale may have 0..N Resolutions and scope may be line/item/quantity-specific;
- cancellation, return, reverse shipment, refund and business-system/economic consequences may coexist and remain separately evidenced;
- provider `available_actions` is capability evidence, not MPC permission or Product API operation vocabulary;
- Post-Sale closes only when applicable consequence owners provide sufficient evidence; no single provider terminal flag or direct `close` write can fabricate closure;
- future concrete human post-sale decisions are admitted only when a selected flow proves their semantics, avoiding a speculative CRM/claims/reverse-logistics platform.

## 6.6 Operational Work

| Candidate operation | Client | Class | Permission | Idempotency / concurrency | Admission |
|---|---|---|---|---|---|
| `ListWork` | both | Q | `work.read` | — | **ADMIT** |
| `GetWork` | both | Q | `work.read` | — | **ADMIT** |
| `AssignWork` | both | C | `work.manage` | structural + concurrency | **ADMIT** |
| `ClearWorkAssignment` | both | C | `work.manage` | structural + concurrency | **ADMIT** |
| `HoldWork` | both | C | `work.manage` | structural + concurrency | **ADMIT** |
| `ResumeWork` | both | C | `work.manage` | structural + concurrency | **ADMIT** |
| `EscalateWork` | both | C | `work.manage` | structural/current-state | **ADMIT** |
| `SubmitWorkResolution` | both | C | `work.manage` | **idempotency + current-state precondition** | **ADMIT** |
| direct `CreateWork` | both | C | — | — | **REJECT BASELINE** |
| direct `CloseWork` | both | C | — | — | **REJECT** |
| generic `DismissWork` | both | C | — | — | **REJECT BASELINE** |
| generic Task/Case CRUD | both | Q/C | — | — | **REJECT** |

Binding decisions:

- baseline Work is born from a source-domain committed actionable condition, not arbitrary user task creation;
- Work responsibility/assignment is distinct from D2 AccessRole/Permission;
- `SubmitWorkResolution` submits/points to evidence for source-owner evaluation; it does not set the source condition or Work to closed by fiat;
- closing Work alone never declares source truth resolved; source resolution/reconciliation governs legitimate closure;
- no-ownerless-work remains binding and duplicate source-condition delivery must reconcile to one obligation rather than create silent duplicate work.

Permission floor: `work.read`, `work.manage`.

## 6.7 Read-only P composition

One baseline composition is admitted:

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `GetSaleOperationalView` | both | P | component read permissions | **ADMIT** |
| global cockpit / Product360 / Marketplace360 | both | P | — | **DEFER** |
| generic cross-owner query surface | both | P | — | **REJECT** |
| P mutation / projection concurrency authority | both | C/P | — | **REJECT** |

`GetSaleOperationalView` may compose Sales + Materialization + Fulfillment + Shipment + Post-Sale summary + Work attention references and may expose a derived `OperationalStage` for convenience.

Binding rules:

- P is explicitly read-only and never owns business meaning, authorization, retry or write preconditions;
- component freshness/partiality remains visible where material; no global snapshot time is fabricated;
- any concurrency token/version from the projection cannot authorize owner writes;
- baseline access requires the relevant component read permissions rather than a broad projection permission that bypasses owners;
- Economics is not included by default because it has its own `GetSaleEconomics` surface and may carry different access sensitivity;
- broader cockpit/D6-shaped compositions wait for a real repeated consumer need.

## 6.8 Permission floor

- `fulfillment.read`
- `fulfillment.execute`
- `fulfillment.manage`
- `post_sale.read`
- `post_sale.manage`
- `work.read`
- `work.manage`

`fulfillment.execute` is separate from `fulfillment.manage`: operators who perform physical work do not automatically administer Fulfillment Node configuration.

## 6.9 Explicit Block 5 exclusions

Not admitted:

- generic `OrderWorkflow`, `Task`, `Case`, `AdvanceStatus` or mutable OperationalStage authority;
- arbitrary user-created Work/tasks;
- direct Work close as substitute for source-domain resolution;
- provider Shipment/Claim/Return/refund protocol mirror;
- generic fulfillment/provider requirement/artifact framework;
- direct client invoice command from Fulfillment;
- generic sync/refresh/retry commands;
- speculative support for every fulfillment mode or post-sale action;
- global cockpit / Product360 / broad cross-owner query language before a real consumer proves it;
- any write, authorization or concurrency authority through P.

## 6.10 Block 5 method outcome

**Parent structure:** `CURRENT D1/D2/D3/D4 STRUCTURE CONFIRMED`.

**B2 outcome:**

> **Expose owner-specific physical Fulfillment checkpoints, source-qualified Shipment observation, canonical Post-Sale Resolution and canonical Operational Work; admit one minimal read-only Sale operational composition; reject generic status/workflow/task/provider-action surfaces.**

---

# 7. Whole-Matrix review — EXACT NEXT WORK

Blocks 1–5 now complete the first operation-admission pass. **Do not move directly to path/schema spelling.**

Run a D5-B2 Whole-Matrix Global Coherence review using the Method across the admitted surface as one system. The review must challenge at least:

1. **Duplicate / missing authority:** every Product 1.0 actor outcome has a legitimate surface or owner reaction, and no operation creates a second authority.
2. **Client-class correctness:** human vs automation/system access is no broader than the facts/actions each class can legitimately establish.
3. **Permission coherence / least privilege:** permission families are materially useful without one-per-endpoint fragmentation or broad bypasses.
4. **Q/C/P honesty:** reads, stateless capabilities, durable intent/resource writes and projections use the correct interaction class.
5. **Idempotency / concurrency:** default consequential safety is retained; every structural exemption has a real owner anchor; lost-update and duplicate-create failure classes are not confused.
6. **Owner-triggered vs client-triggered behavior:** D3 reactions such as BusinessOrder/Invoicing creation are not accidentally reintroduced as Product commands.
7. **Knowledge/freshness/evidence:** unknown/unavailable/partial/stale/provider-enriched evidence remains honest in every relevant read family.
8. **External identity / Organization:** no bare provider/native IDs, hidden Organization inference or cross-Organization secondary reference.
9. **Generic abstraction pressure:** no Product/PIM, generic Integration/Mutation/Workflow/Rule/Finance/Task/Operation/Provider graph emerges from repeated mechanics.
10. **Pagination/filter/bulk:** admitted only for real collections/workflows and not by API symmetry.
11. **P projection safety:** projections cannot bypass permissions, become snapshot/concurrency authority or absorb business writes.
12. **Future-cost/YAGNI:** the seams needed for second marketplace/business system, stronger automation and later concrete post-sale/fulfillment modes remain possible without implementing them today.
13. **Structural inversion:** accepted surface must remain correct if legacy routes/OpenAPI/controllers are opposite in every relevant respect.
14. **Missing-operation challenge:** explicitly look for a Product 1.0 requirement that is currently unreachable from a legitimate client or owner-triggered reaction.

Outcomes remain:

- `RESTRUCTURE NOW`
- `CURRENT STRUCTURE CONFIRMED`
- `STOP / SPLIT PREREQUISITE`

Only after the Whole-Matrix review is ratified should B2 move to resource/path grammar, HTTP methods/custom methods, request/response schemas, statuses, precondition/idempotency headers, pagination/filter grammar and OpenAPI spelling.

Implementation remains blocked until D9.