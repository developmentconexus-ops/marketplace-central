# D5-B2 — Operation Admission Matrix

> **Status:** OPEN / ACTIVE — B2-A + Blocks 1–5 + Whole-Matrix Global Coherence **ACCEPTED IN-STAGE / OPERATOR-RATIFIED**; Wire Contract next  
> **Parent B2:** `D5-B2-PRODUCT-OPERATION-SURFACE.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + Decision Reconciliation Baseline + B2-A  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-18  
> **Whole-Matrix ratified:** 2026-08-18

## 1. Governing admission and safety rules

The Product API contains the smallest operation surface real Product 1.0 clients need before final path/schema spelling.

Every admitted operation must identify proportionately:

- real consumer/use;
- allowed client class: human, machine/automation/system, or both;
- exactly one accepted semantic owner or D2 substrate authority;
- Q / C / P class;
- explicit Organization scope, except the bounded self-only access-context discovery named below;
- ordinary Permission, distinct from business disposition/Governance;
- canonical/source-qualified subject identity;
- honest knowledge/freshness/provenance when read semantics require it;
- outcome/idempotency/precondition/concurrency semantics when applicable;
- provider enrichment only for a named consumer/correctness property;
- pagination/filter/sort only for a real collection consumer;
- bulk only when a real workflow requires member-level semantics.

Reject a candidate whose only justification is API symmetry, current code, provider endpoint shape, debug convenience or hypothetical future need.

### 1.1 Complete C-operation safety tuple

Every **admitted C operation** declares all three fields below. Silence is non-conformant.

```text
consequence class:
  consequential | non-consequential / side-effect-free

idempotency:
  mandatory client key
  | named structural owner anchor / exemption
  | N/A with explicit reason

concurrency / precondition:
  required
  | not material with explicit reason
```

For consequential C, D5-B1's fail-closed default remains **mandatory client idempotency** unless this matrix names a structural owner anchor making duplicate intake unreachable or harmless.

Idempotency and concurrency solve different failure classes and never substitute for one another. A client idempotency key never authorizes blind replay of an external effect whose acceptance is ambiguous.

The complete ratified sweep is in §7.

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

- `GetCurrentAccessContext` is the one bounded exception to Organization-owned routing: it is a **platform-scoped, self-only D2 discovery Q** for the authenticated Principal resolved from trusted token context; it accepts no Principal identifier and returns only Organizations where that Principal currently has Membership;
- cross-Principal access enumeration remains Organization-path-scoped `ListOrganizationMembers` under `access.read`;
- IdP authentication does not carry MPC Organization/Permission authority by implication;
- assignment/revocation acts on one explicit Membership + AccessRole relation, never a whole-role-set replacement;
- access-role revocation is fail-safe and monotonic for the targeted standing authority: stale snapshots do not keep the grant alive merely because the revoker did not re-read; later re-grant is a new explicit assignment;
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
| explicit Marketplace Installation reactivation | human | C | — | **DEFER** |

Binding decisions:

- Installation creation is duplicate-sensitive and uses client idempotency;
- configuration updates are desired-state writes protected from stale overwrite;
- deactivation preserves identity/history and is not provider-account deletion;
- reactivation is not inferred from deactivation and remains deferred until a real workflow requires it;
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
- correspondence replacement/clearing uses current-state preconditions so stale automation cannot defeat a newer decision;
- automation cannot silently supersede a standing human correspondence decision; explicit supersession remains attributed;
- exact repeated correspondence meaning is structurally idempotent when the current version still matches.

Permission floor: `readiness.read`, `readiness.manage`.

**Block 1 outcome:** `CURRENT PARENT STRUCTURE CONFIRMED; ADMIT ONLY CONSUMER-PROVEN OPERATIONS`.

---

# 3. Block 2 — Offering + Price + Availability — ACCEPTED IN-STAGE

## 3.1 Governing invariant

> **Marketplace Listing observation, ListingIntent authoring, PriceIntent actuation, Availability-owned sellable quantity/configuration and provider execution remain distinct meanings. Provider protocol may physically serialize them together, but Product API contracts never collapse those meanings into one giant Listing mutation or generic asynchronous Operation.**

## 3.2 Marketplace Listing observation

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListMarketplaceListings` | both | Q | `offering.read` | **ADMIT** |
| `GetMarketplaceListing` | both | Q | `offering.read` | **ADMIT** |
| direct `CreateMarketplaceListing` | both | C | — | **REJECT** |
| direct `UpdateMarketplaceListing` | both | C | — | **REJECT** |
| direct `Pause/Reactivate/CloseListing` APIs | both | C | — | **REJECT AS PARALLEL BASELINE AUTHORITY** |

Listing identity remains Marketplace Installation + provider-native Listing key. Read contracts expose Offering interpretation + source freshness/convergence, never raw provider DTOs.

## 3.3 ListingIntent authoring and listing-scoped media

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListListingIntents` | both | Q | `offering.read` | **ADMIT** |
| `GetListingIntent` | both | Q | `offering.read` | **ADMIT** |
| `CreateListingIntentDraft` | both | C/create | `listing.manage` | **ADMIT** |
| `UpdateListingIntentDraft` | both | C/update | `listing.manage` | **ADMIT** |
| `DiscardListingIntentDraft` | both | C | `listing.manage` | **ADMIT** |
| `SubmitListingIntent` | both | C | `listing.manage` | **ADMIT** |
| `CreateListingIntentMedia` | both | C/create | `listing.manage` | **ADMIT** |
| separate `ValidateListingIntent` | both | Q/C | — | **REJECT BASELINE** |
| separate `PreviewListing` | both | Q | — | **REJECT BASELINE** |
| generic ProductAsset/media-library CRUD | both | Q/C | — | **REJECT** |

Binding decisions:

- create/edit use one ListingIntent identity (`target = none | existing source-qualified Listing`);
- draft updates are declarative desired-state changes, not a mutation DSL;
- only `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` baseline value modes are admitted;
- stale automation cannot overwrite or submit against a newer human decision;
- `SubmitListingIntent` is the consequential client boundary; technical freeze/authorize/execute/reconcile are not separate Product operations;
- `submitted != authorized != externally applied != converged` remains binding;
- no generic `LongRunningOperation` is admitted because owner-local Intents already provide durable tracking identity;
- `CreateListingIntentMedia` is scoped to an exact mutable ListingIntent context and does **not** create Product/PIM/media-master authority;
- arbitrary client-supplied external URLs are not trusted authored media; external/source media remains a D4 evidence/acquisition concern;
- media reference/selection/order/role is read through `GetListingIntent`; blob upload/storage/hash/resizing/binary delivery is D7 mechanism and does not justify a separate Product media read resource by symmetry.

## 3.4 PriceIntent — including initial publication

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListPriceIntents` | both | Q | `offering.read` | **ADMIT** |
| `GetPriceIntent` | both | Q | `offering.read` | **ADMIT** |
| `CreatePriceIntent` | both | C | `price.manage` | **ADMIT** |
| direct `SetListingPrice` | both | C | — | **REJECT** |
| public mutable PriceDraft | both | C | — | **REJECT BASELINE** |
| withdraw/cancel pending PriceIntent | both | C | — | **DEFER** |

Binding decisions:

- price is **never ListingIntent-owned content**, including initial publication;
- `PriceIntent` expresses desired exact Money target and remains distinct from Economics reasoning/simulation;
- target duality is an attribute of the same PriceIntent identity: `existing source-qualified Listing | pre-creation ListingIntent context`; no new intent class is created;
- changing the intended pre-dispatch price creates a newer PriceIntent that explicitly supersedes the pending one with preserved attribution/lineage; public mutable PriceDraft remains rejected;
- automation cannot silently supersede a standing human-authored pending PriceIntent;
- `GetListingIntent` may expose the typed correlated PriceIntent identity, never an embedded ListingIntent-owned price value;
- after provider dispatch, representation convergence belongs to ListingIntent and price convergence belongs to PriceIntent independently.

Permission split remains material: `listing.manage != price.manage`.

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

## 3.6 Joint active-publication realization

D4-R1 R1-G1 remains binding and is extended without ownership merge for PriceIntent.

For a provider lane that physically requires representation + price + quantity in one create request:

```text
ListingIntent
  + correlated current PriceIntent
  + Availability-issued current meaning / intent
        ↓
D4/D7 validate + correlate
        ↓
one provider request when required
        ↓ authoritative rereads
ListingIntent convergence
PriceIntent convergence
Availability convergence
```

Rules:

- `SubmitListingIntent` never owns or embeds price or Availability quantity;
- before **provider dispatch**, the current required PriceIntent and Availability-issued input must exist and remain valid/current; missing required input is fail-closed/non-dispatchable, never defaulted by adapter convenience;
- submission may remain accepted/pending under owner semantics while execution prerequisites are unresolved, but it is never eternal authorization;
- no cross-owner atomicity is claimed;
- provider blast radius never silently widens intended/authorized scope.

Permission floor: `offering.read`, `listing.manage`, `price.manage`, `availability.read`, `availability.manage`.

**Block 2 outcome:** `CURRENT D1/D2/D3/D4-R1 STRUCTURE CONFIRMED; KEEP LISTING/PRICE/AVAILABILITY MEANING DISTINCT THROUGH INITIAL PUBLICATION`.

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

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListExpectedEconomics` | both | Q | `economics.read` | **ADMIT** |
| `GetExpectedEconomics` | both | Q | `economics.read` | **ADMIT** |
| `EvaluatePriceScenario` | both | C — stateless/side-effect-free | `economics.read` | **ADMIT** |
| persistent `Simulation` resource | both | C | — | **REJECT BASELINE** |
| persistent `Recommendation` resource | both | C | — | **REJECT BASELINE** |
| generic mutable `Profitability` blob | both | Q/C | — | **REJECT** |

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

Commercial policy remains Economics-owned; Governance does not acquire business thresholds. Economic Attribution is persistent Economics state; explicit ambiguity resolution is human baseline.

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
| `EstablishAuthorizationDelegation` | human | C/create | `governance.manage` | **ADMIT** |
| `UpdateAuthorizationDelegation` | human | C/update | `governance.manage` | **ADMIT** |
| `RevokeAuthorizationDelegation` | human | C | `governance.manage` | **ADMIT** |
| separate `ListPendingApprovals` queue authority | human | Q | — | **REJECT BASELINE** |
| domain-specific `ApproveListingIntent` / `ApprovePriceIntent` | human | C | — | **REJECT** |
| `ExecuteApprovedAction` | both | C | — | **REJECT** |
| generic approval/policy/rules engine | both | C | — | **REJECT** |

Binding decisions:

- Authorization Decision is a durable Governance-owned occurrence, not an `approved=true` field or Intent-status mutation;
- a Decision targets one concrete owner Intent/material revision/context and preserves decision Principal, authority context, outcome and exact authorized target-scope snapshot;
- authorized scope never widens intended scope;
- approval never executes, mutates the owner Intent, waives domain validity or becomes eternal authorization;
- actionable approval responsibility belongs to Work/owner attention, not a second Governance queue authority;
- delegation establishment is duplicate-sensitive unless later Wire Contract proves a unique semantic-keyed upsert with a safe structural anchor;
- delegation modification is a stale-update problem and requires current-state protection;
- delegation **revocation is monotonic/fail-safe** for the targeted standing authority: stale snapshots do not keep broader authority alive; later re-grant/re-establishment is a new explicit authority action and history is never rewritten.

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

`ResolveBusinessSystemPartyResolution` is duplicate-sensitive because zero-match resolution can lead to a native-party effect; it always uses a client idempotency key and a current resolution/candidate-set precondition.

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

> **Fulfillment owns physical execution and provider-readiness closure, Post-Sale owns scoped consequence coordination, Work owns responsibility/assignment/escalation, and any future P composition is read-only convenience. Physical facts, consequence closure, actionable-work lifecycle and projection convenience never collapse into one generic OrderWorkflow/Task/Case/status authority.**

## 6.2 Fulfillment Lifecycle

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListFulfillmentStates` | both | Q | `fulfillment.read` | **ADMIT** |
| `GetFulfillmentState` | both | Q | `fulfillment.read` | **ADMIT** |
| `ListFulfillmentNodes` | both | Q | `fulfillment.read` | **ADMIT** |
| `GetFulfillmentNode` | both | Q | `fulfillment.read` | **ADMIT** |
| `CreateFulfillmentNode` | human | C | `fulfillment.manage` | **ADMIT** |
| `UpdateFulfillmentNode` | human | C | `fulfillment.manage` | **ADMIT** |
| `DeactivateFulfillmentNode` | human | C | `fulfillment.manage` | **ADMIT** |
| `GetFulfillmentOperatingTargets` | both | Q | `fulfillment.read` | **ADMIT** |
| `UpdateFulfillmentOperatingTargets` | human | C/update | `fulfillment.manage` | **ADMIT** |
| manual `SelectFulfillmentNode` | human | C | — | **DEFER** |
| `RecordSeparation` | human baseline | C | `fulfillment.execute` | **ADMIT** |
| `RecordPhysicalConference` | human or explicitly proven physical-system Principal | C | `fulfillment.execute` | **ADMIT** |
| `RecordPacking` | human baseline | C | `fulfillment.execute` | **ADMIT** |
| `RecordDispatchHandoff` | human or explicitly proven physical-system Principal | C | `fulfillment.execute` | **ADMIT** |
| generic `AdvanceFulfillmentStatus` | both | C | — | **REJECT** |

Binding decisions:

- physical checkpoints are owner-specific facts/capabilities, not `PATCH status` workflow steps;
- `RecordPhysicalConference` is the material Fulfillment fact that may awaken Materialization for invoicing under D3; clients never bypass it with a direct invoice command;
- an ordinary automation token cannot fabricate physical evidence merely because it has API access; machine establishment requires a separately proven system Principal/source capable of establishing the fact;
- repeat/lost-response safety prevents duplicate material occurrences and duplicate downstream reactions;
- Fulfillment operating targets are MPC-owned internal target policy, distinct from provider deadlines/windows;
- `GetFulfillmentOperatingTargets` exposes effective value plus material provenance/default-vs-override explanation;
- future internal targets remain owner-local and consumer-proven; this does not create a generic SLA/target/rules platform;
- no company-wide WMS/TMS model is introduced.

## 6.3 Provider requirements and Fulfillment artifacts

`GetFulfillmentState` may expose Fulfillment-owned closure/readiness plus source-qualified provider requirement/deadline evidence sufficient for the claimed path. Provider requirement/artifact truth remains D4/external authority.

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListFulfillmentArtifacts` | human/both read | Q | `fulfillment.execute` | **ADMIT** |
| `GetFulfillmentArtifact` | human/both read | Q | `fulfillment.execute` | **ADMIT** |
| generic `CreateArtifact` | both | C | — | **REJECT** |
| generic provider report/artifact generation | both | C | — | **REJECT BASELINE** |

`fulfillment.execute` protecting artifact reads is intentional least privilege because labels/handoff documents may contain operational/PII-sensitive material. Artifacts remain Fulfillment-local/source-qualified and PII-minimized. A future provider action that generates an artifact is consequential and requires its own adjudication; it is never disguised as a read.

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

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListPostSaleResolutions` | both | Q | `post_sale.read` | **ADMIT** |
| `GetPostSaleResolution` | both | Q | `post_sale.read` | **ADMIT** |
| `CreatePostSaleResolution` | both | C/create | `post_sale.manage` | **ADMIT** |
| generic `UpdateResolutionStatus` | both | C | — | **REJECT** |
| direct `ClosePostSaleResolution` | both | C | — | **REJECT** |
| direct `CancelSale` / `RefundSale` | both | C | — | **REJECT** |
| provider-shaped `AcceptClaim` / `ReviewReturn` | both | C | — | **REJECT BASELINE** |
| generic `SubmitPostSaleAction` | both | C | — | **REJECT BASELINE** |
| concrete later post-sale decisions | human/both | C | — | **DEFER UNTIL SELECTED ACTION IS PROVEN** |

Binding decisions:

- a Resolution may be provider-originated or explicitly MPC-initiated; Product creation creates the scoped coordination obligation, not all underlying effects;
- one Sale may have 0..N Resolutions and scope may be line/item/quantity-specific;
- cancellation, return, reverse shipment, refund and business-system/economic consequences may coexist and remain separately evidenced;
- provider `available_actions` is capability evidence, not MPC Permission or Product operation vocabulary;
- Post-Sale closes only when applicable consequence owners provide sufficient evidence; no direct close or single provider terminal flag fabricates closure;
- future concrete post-sale decisions are admitted only when a selected flow proves their semantics.

## 6.6 Operational Work

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListWork` | both | Q | `work.read` | **ADMIT** |
| `GetWork` | both | Q | `work.read` | **ADMIT** |
| `AssignWork` | both | C | `work.manage` | **ADMIT** |
| `ClearWorkAssignment` | both | C | `work.manage` | **ADMIT** |
| `HoldWork` | both | C | `work.manage` | **ADMIT** |
| `ResumeWork` | both | C | `work.manage` | **ADMIT** |
| `EscalateWork` | both | C | `work.manage` | **ADMIT — DECLARATIVE TARGET/STATE ONLY** |
| generic `SubmitWorkResolution` | both | C | `work.manage` | **DEFER** |
| direct `CreateWork` | both | C | — | **REJECT BASELINE** |
| direct `CloseWork` | both | C | — | **REJECT** |
| generic `DismissWork` | both | C | — | **REJECT BASELINE** |
| generic Task/Case CRUD | both | Q/C | — | **REJECT** |

Binding decisions:

- baseline Work is born from a source-domain committed actionable condition, not arbitrary user task creation;
- Work responsibility/assignment is distinct from D2 AccessRole/Permission;
- generic `SubmitWorkResolution` is deferred because source-owner-specific operations already own concrete semantic resolutions; Work must not become a command bus;
- closing Work alone never declares source truth resolved; source resolution/reconciliation governs legitimate closure;
- `EscalateWork` is admitted only as a declarative desired escalation target/state. If a future wire contract instead needs increment/occurrence semantics, it becomes duplicate-sensitive and requires client idempotency;
- no-ownerless-work remains binding and duplicate source-condition delivery must reconcile to one obligation rather than create silent duplicate work.

### Work closure-path obligation before Wire Contract closure

For every Product 1.0 condition class that can produce Work — including missing evidence, staleness, deadline breach, delivery exception, ambiguous external effect and divergence — Wire Contract review must prove at least one legitimate closure path:

1. source-owner automatic resolution/reconciliation propagated under D3; or
2. an admitted owner-specific human capability.

A human-only evidence class with no legitimate source-owner closure path is the concrete trigger for the smallest bounded Work→source evidence-submission operation. It never transfers source truth to Work.

Permission floor: `work.read`, `work.manage`.

## 6.7 Read-only P composition — ZERO-P BASELINE

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `GetSaleOperationalView` | both | P | component read permissions | **DEFER — D6 CONSUMER PROOF** |
| global cockpit / Product360 / Marketplace360 | both | P | — | **DEFER** |
| generic cross-owner query surface | both | P | — | **REJECT** |
| P mutation / projection concurrency authority | both | C/P | — | **REJECT** |

No P operation is required merely because D5-B1 permits P. Client-side composition from owner Qs is the baseline until D6 proves repeated cross-owner composition pain.

`OperationalStage` may be derived client-side for UX. Materially divergent/repeated derivation across consumers is evidence to reconsider a bounded read-only P; any future P remains non-authoritative for writes, authorization, retry and concurrency.

## 6.8 Permission floor

- `fulfillment.read`
- `fulfillment.execute`
- `fulfillment.manage`
- `post_sale.read`
- `post_sale.manage`
- `work.read`
- `work.manage`

`fulfillment.execute` remains separate from `fulfillment.manage`: physical operators do not automatically administer Fulfillment Node/target policy.

**Block 5 outcome:** `CURRENT D1/D2/D3/D4 STRUCTURE CONFIRMED; OWNER-SPECIFIC FULFILLMENT/POST-SALE/WORK, ZERO-P UNTIL CONSUMER PROOF`.

---

# 7. Whole-Matrix complete C-operation safety sweep — ACCEPTED / RATIFIED

The tables below are the reviewable D5-B1 safety declaration for **every admitted C operation**. Rejected/deferred operations have no Product runtime contract yet and are intentionally absent.

## 7.1 Identity / Portfolio / Readiness

| Operation | Consequence class | Idempotency disposition | Concurrency / precondition disposition |
|---|---|---|---|
| `AssignAccessRole` | consequential access-state mutation | structural anchor: unique current Membership + AccessRole relation; exact repeat harmless | optimistic resource version not material for one-relation add; current same-Organization Membership/authority validity required |
| `RevokeAccessRole` | consequential authority removal | structural anchor: removal of targeted standing relation; monotonic exact-repeat safe | **no stale-snapshot block**; fail-safe removal wins, later re-grant is explicit |
| `CreateMarketplaceInstallation` | consequential durable create | **mandatory client key** | create has no prior resource-version axis; request-time business/auth validity still required |
| `UpdateMarketplaceInstallationConfiguration` | consequential desired-state mutation | structural desired-state repeat | **current Installation configuration precondition required** |
| `DeactivateMarketplaceInstallation` | consequential lifecycle mutation | structural lifecycle anchor | **current Installation state/configuration precondition required**; stale deactivation can be unsafe |
| `ResolveProductChannelCorrespondence` | consequential semantic resolution | structural anchor: one current correspondence meaning; exact repeat harmless | **current correspondence revision required**; human supersession safety preserved |
| `ClearProductChannelCorrespondence` | consequential semantic lifecycle mutation | structural current-correspondence clear | **current correspondence revision required** |

## 7.2 Offering / Availability

| Operation | Consequence class | Idempotency disposition | Concurrency / precondition disposition |
|---|---|---|---|
| `CreateListingIntentDraft` | consequential durable intent create | **mandatory client key** | no stale prior Intent axis; target identity/currentness is validated semantically |
| `UpdateListingIntentDraft` | consequential desired-state mutation | structural desired-state repeat | **current draft revision required** |
| `DiscardListingIntentDraft` | consequential lifecycle mutation | structural Intent-state anchor | **current draft revision/state required** |
| `SubmitListingIntent` | consequential Intent lifecycle transition | structural anchor: exact Intent + submitted revision; repeat of same accepted transition harmless | **current draft revision required**; later external dispatch revalidates all current prerequisites |
| `CreateListingIntentMedia` | consequential listing-context media create | **mandatory client key** | **current mutable ListingIntent revision/state required** for association context |
| `CreatePriceIntent` | consequential durable price intent create | **mandatory client key** | new target intake has no stale resource axis; **explicit current superseded-PriceIntent precondition required when replacing a pending pre-creation intent** |
| `CreateInventorySource` | consequential durable create | **mandatory client key** | create has no prior resource-version axis |
| `UpdateInventorySource` | consequential desired-state mutation | structural desired-state repeat | **current Inventory Source revision required** |
| `DeactivateInventorySource` | consequential lifecycle mutation | structural lifecycle anchor | **current Inventory Source revision required** |
| update allocation/scope policy | consequential policy mutation | structural desired-state repeat | **current policy revision required** |

## 7.3 Commercial Economics

| Operation | Consequence class | Idempotency disposition | Concurrency / precondition disposition |
|---|---|---|---|
| `EvaluatePriceScenario` | **non-consequential / side-effect-free** | **N/A — no durable intake/effect** | **N/A — hypothetical evaluation; evidence freshness is owner semantics, not optimistic concurrency** |
| `UpdateCommercialPolicy` | consequential policy mutation | structural desired-state repeat | **current Commercial Policy revision required** |
| `ResolveEconomicAttribution` | consequential semantic resolution | structural anchor: one current EconomicAttribution meaning; exact repeat harmless | **current attribution revision required** |

## 7.4 Governance / Sales / Materialization

| Operation | Consequence class | Idempotency disposition | Concurrency / precondition disposition |
|---|---|---|---|
| `CreateAuthorizationDecision` | consequential durable decision occurrence | **mandatory client key** | **exact target Intent/material revision/context precondition required** |
| `EstablishAuthorizationDelegation` | consequential standing-authority create | **mandatory client key by default**; Wire Contract may claim structural exemption only if a unique semantic-keyed upsert is proven unable to mint parallel grants | create has no prior delegation-version axis unless a semantic-keyed upsert is later proven; current actor authority required |
| `UpdateAuthorizationDelegation` | consequential authority mutation | structural desired-state repeat | **current delegation revision required** |
| `RevokeAuthorizationDelegation` | consequential authority removal | structural monotonic removal | **no stale-snapshot block**; revoke current targeted standing authority, later re-grant is a new explicit action |
| `ResolveSaleSellingEntityAttribution` | consequential Sales-owned semantic resolution | structural anchor: one current attribution for the source-qualified Sale | **current Sale/attribution revision required** |
| `ResolveBusinessSystemPartyResolution` | consequential semantic resolution that may cause native-party effect | **mandatory client key** | **current Party Resolution + candidate-set revision required by default** |

## 7.5 Fulfillment / Post-Sale / Work

| Operation | Consequence class | Idempotency disposition | Concurrency / precondition disposition |
|---|---|---|---|
| `CreateFulfillmentNode` | consequential durable create | **mandatory client key** | create has no prior resource-version axis |
| `UpdateFulfillmentNode` | consequential desired-state mutation | structural desired-state repeat | **current Fulfillment Node revision required** |
| `DeactivateFulfillmentNode` | consequential lifecycle mutation | structural lifecycle anchor | **current Fulfillment Node revision required** |
| `UpdateFulfillmentOperatingTargets` | consequential policy mutation | structural desired-state repeat | **current operating-target policy revision required** |
| `RecordSeparation` | consequential physical checkpoint occurrence | **mandatory client key** | **current Fulfillment state/precondition required** |
| `RecordPhysicalConference` | consequential physical checkpoint occurrence | **mandatory client key** | **current Fulfillment state/precondition required** |
| `RecordPacking` | consequential physical checkpoint occurrence | **mandatory client key** | **current Fulfillment state/precondition required** |
| `RecordDispatchHandoff` | consequential physical checkpoint occurrence | **mandatory client key** | **current Fulfillment state/precondition required** |
| `CreatePostSaleResolution` | consequential durable Resolution create | **mandatory client key** | create has no prior Resolution-version axis; scope/current Sale validity still checked |
| `AssignWork` | consequential Work desired-state mutation | structural desired assignment | **current Work revision required** |
| `ClearWorkAssignment` | consequential Work desired-state mutation | structural clear | **current Work revision required** |
| `HoldWork` | consequential Work lifecycle mutation | structural desired hold state | **current Work revision required** |
| `ResumeWork` | consequential Work lifecycle mutation | structural desired resumed state | **current Work revision required** |
| `EscalateWork` | consequential Work mutation | structural **only because admitted semantics are declarative to an explicit escalation target/state** | **current Work revision required**; increment/occurrence semantics would require a mandatory client key and are not baseline |

### 7.6 Safety-sweep result

- admitted C operations with silent idempotency disposition: **0**;
- admitted C operations with silent concurrency/precondition disposition: **0**;
- idempotency used as concurrency authority: **0**;
- concurrency token used as duplicate-intake identity: **0**;
- generic replay operation after ambiguous external effect: **0**.

Exact header/mechanism spelling remains Wire Contract/D7 as appropriate; the protected semantic properties above are already B2 authority.

---

# 8. Whole-Matrix Global Coherence — ACCEPTED / OPERATOR-RATIFIED

Independent Fable challenge + GPT adjudication converged after one focused Round 2. Reviewer output was evidence, never authority. The operator ratified the converged package on 2026-08-18.

## 8.1 Applied B2-local corrections

1. **ADD** ListingIntent-scoped authored-media intake; no Product/media master.
2. **ADD** Fulfillment-owned internal operating-target Q/C with effective-value provenance; no generic SLA/rules authority.
3. **DEFER** generic `SubmitWorkResolution`; source-owner closure remains authoritative and Wire Contract must audit closure-path coverage.
4. **DEFER** `GetSaleOperationalView`; zero-P baseline until D6 proves repeated composition need.
5. **COMPLETE** all admitted C safety declarations using §7; silence is prohibited.
6. **HARDEN** Party Resolution with mandatory idempotency + current candidate-set/resolution precondition.
7. **HARDEN** current access discovery as platform-scoped self-only.
8. **HARDEN** authority revocation as fail-safe/monotonic for AccessRole and Authorization Delegation.
9. **CONFIRM** B2-A OIDC/OAuth + MPC-owned Principal/Membership/Permission as the Global Maximum at this stage.
10. **CONFIRM** initial price always uses PriceIntent; ListingIntent may correlate but never owns price.

## 8.2 Global checks

- duplicate/missing semantic authority: **PASS**;
- Product 1.0 lifecycle reachability: **PASS** after media + fulfillment-target corrections;
- client-class correctness: **PASS**;
- Permission / least privilege: **PASS**;
- Q/C/P honesty: **PASS**, with zero-P baseline;
- idempotency/concurrency completeness: **PASS**;
- owner-trigger vs client-trigger: **PASS**;
- Organization/source-qualified identity: **PASS**;
- provider richness without DTO mirroring: **PASS**;
- generic abstraction pressure / YAGNI: **PASS** after Work/P defers;
- bulk by symmetry: **REJECTED / PASS**;
- Structural Inversion against legacy routes/OpenAPI: **PASS**;
- parent-stage reopen: **NONE**.

### Whole-Matrix disposition

```text
D0 / D1 / D2 / D3 / D4 / D4-R1 / D5-B1     CURRENT STRUCTURE CONFIRMED
D5-B2 operation inventory                    RESTRUCTURE NOW corrections APPLIED
D5-B2 Whole-Matrix Global Coherence          ACCEPTED / RATIFIED
```

---

# 9. Exact next B2 work — Wire Contract / Resource-Path-Schema Grammar

The semantic operation inventory is now stable enough to crystallize wire shape. Next work must derive the concrete contract from this matrix, not from legacy OpenAPI/routes/controllers.

The next sub-batch must decide:

1. resource/path hierarchy, keeping Organization-owned business operations under `/organizations/{organization_id}/...` and giving `GetCurrentAccessContext` only the bounded self-only platform discovery shape accepted above;
2. standard HTTP resource methods versus owner-specific methods where CRUD would lie;
3. exact request/response schema families and source-qualified identity representation;
4. read knowledge semantics (`known`, known-empty, unknown, unavailable, partial, freshness/provenance where material) without a universal Fact wrapper by habit;
5. consequential owner outcomes (`accepted`, `rejected`, `pending`, `ambiguous` and later applied/converged distinctions where applicable);
6. RFC 9457 Problem Details for API/transport/access/precondition/idempotency/server problems without turning valid business rejection into HTTP access failure;
7. exact `Idempotency-Key` placement/validation and the opaque MPC concurrency/precondition mechanism for every §7 row;
8. pagination/filter/search/cursor grammar only for admitted collection consumers;
9. exact Permission→wire-operation mapping and client-class restrictions;
10. media wire seam without choosing D7 blob/storage/CDN realization;
11. technical provider OAuth/webhook/external-connector ingress classification outside the Product API, as required by D5-B1/D4;
12. the Work closure-path audit required by §6.6 before Wire Contract closure;
13. OpenAPI operation naming/spelling and the path toward one machine-readable Product API wire authority.

Do not introduce D6 screen/BFF topology, D7 queues/workers/storage/transactions/Keycloak deployment, D8 effect proofs or implementation.

Implementation remains blocked until D9.
