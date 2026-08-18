# D5-B2 — Operation Admission Matrix

> **Status:** OPEN / ACTIVE — Blocks 1, 2 and 3 ACCEPTED IN-STAGE; Block 4 next  
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

Rejected local maximum:

```text
PATCH Listing
  = content
  + price
  + stock
  + fulfillment
  + provider protocol
```

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

- create/edit use the same ListingIntent identity (`target = none | existing Listing`);
- draft updates are declarative desired-state changes, not a mutation-command DSL;
- only `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` baseline value modes are admitted;
- `GetListingIntent` can expose resolved values, requirements/provenance, dispatchability/blockers, lifecycle/convergence and concurrency token;
- stale automation cannot overwrite or submit against a newer human decision;
- `SubmitListingIntent` is the consequential client boundary; clients do not separately invoke technical freeze/authorize/execute/reconcile steps;
- `submitted != authorized != externally applied != converged` remains binding;
- no generic LongRunningOperation is admitted because owner-local Intents already provide the durable tracking identity.

Listing pause/reactivate/close/edit remain ListingIntent semantics when admitted; provider-shaped direct mutations do not become a parallel architecture.

## 3.4 PriceIntent

| Candidate operation | Client | Class | Permission | Idempotency | Admission |
|---|---|---|---|---|---|
| `ListPriceIntents` | both | Q | `offering.read` | — | **ADMIT** |
| `GetPriceIntent` | both | Q | `offering.read` | — | **ADMIT** |
| `CreatePriceIntent` | both | C | `price.manage` | **mandatory** | **ADMIT** |
| direct `SetListingPrice` | both | C | — | — | **REJECT** |
| public mutable PriceDraft | both | C | — | — | **REJECT BASELINE** |
| withdraw/cancel pending PriceIntent | both | C | — | — | **DEFER** |

`PriceIntent` expresses desired exact Money target. Economics owns reasoning/simulation; Offering owns actuation. No `decrease 5%`/`match competitor` economic logic is moved into Offering.

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

Availability Q preserves desired Sellable Availability, provider actual evidence, knowledge/freshness/provenance and owner-specific convergence. Inventory Source and allocation/scope policy remain distinct meanings. Creation of Inventory Source requires idempotency unless a stronger semantic anchor is proven; unsafe configuration update uses concurrency.

## 3.6 Joint Offering × Availability realization

D4-R1 R1-G1 remains binding:

```text
ListingIntent / Offering-issued meaning
            +
Availability-issued sellable-quantity meaning
            ↓
D4/D7 execution mechanism
            ↓
one provider call when required
            ↓
authoritative rereads
   ├─ Offering convergence
   └─ Availability convergence
```

`SubmitListingIntent` never accepts quantity as Offering-owned content. Missing required Availability input fails closed before provider dispatch. A submitted Offering intent may wait on another owner-issued prerequisite only with later execution-time revalidation; submission is not eternal authorization. No cross-owner atomicity is claimed.

Material multi-step/blast-radius effects preserve intended scope, authorized scope and actual attempted/provider-affected outcomes plus member/aspect confirmed/rejected/pending/ambiguous/not-executed distinctions where required.

Permission floor: `offering.read`, `listing.manage`, `price.manage`, `availability.read`, `availability.manage`.

**Block 2 outcome:** `CURRENT D1/D2/D3/D4-R1 STRUCTURE CONFIRMED; ADMIT OWNER-SPECIFIC RESOURCES/INTENTS, REJECT GIANT LISTING CRUD AND GENERIC ASYNC OPERATION`.

---

# 4. Block 3 — Market Intelligence + Commercial Economics — ACCEPTED IN-STAGE

## 4.1 Governing invariant

> **Market Intelligence owns competitive interpretation/evidence sufficiency; Commercial Economics owns economic meaning, simulation and L0/L1/L2 lineage; Offering alone owns PriceIntent. Hypothetical exploration does not gain durable identity by default, while economic evidence/decisions that materially participate in later explanation or reconciliation remain durable under Economics.**

Rejected extremes:

```text
persist every simulation/recommendation
→ duplicate lifecycle/authority + accidental complexity

calculate everything transiently
→ cannot explain expected vs order vs realized outcomes later
```

Selected structure:

```text
Market Intelligence Q
        ↓ interpreted source-qualified evidence
Commercial Economics stateless scenario evaluation
        ↓ decision support
human / automation decision
        ↓
Offering PriceIntent
        ↓
material decision-time Economic basis retained when required
        ↓
L1 Order Economics
        ↓
L2 Realized Economics
        ↓
R1 / R2 reconciliation + variance/calibration
```

## 4.2 Market Intelligence

| Candidate operation | Consumer / use | Client | Class | Permission | Admission |
|---|---|---|---|---|---|
| `ListCompetitivePositions` | commercial/operator/agent scans market position | both | Q | `market.read` | **ADMIT** |
| `GetCompetitivePosition` | inspect one interpreted competitive position/explanation | both | Q | `market.read` | **ADMIT** |
| `ListComparableOffers` | explain/evidence the comparison behind Market Intelligence conclusion | both | Q | `market.read` | **ADMIT** |
| generic `GetMarketObservation` | provider/evidence-object exposure without consumer semantic need | both | Q | — | **REJECT** |
| `CreateMarketObservation` | client authors external evidence | both | C | — | **REJECT** |
| `CollectMarketNow` / `RefreshCompetitivePosition` | D4/D7 acquisition/runtime mechanism | both | C | — | **REJECT PRODUCT API** |
| manual `SetComparableOffer` | manually alter comparability absent proven workflow | both | C | — | **DEFER** |
| generic scraper/collector configuration | speculative source platform | human | C | — | **REJECT** |

Binding decisions:

- Competitive Position is semantic context, not necessarily a synthetic identity/resource ID;
- provider-rich evidence such as buyer shipping, free-shipping state, provider competition reasons or provider price guidance may be exposed only as bounded source-qualified enrichment needed for explanation/correctness;
- another provider need not fabricate equivalent fields;
- comparable-offer pagination is legitimate where source population is larger, but source/provider population never implies universal market completeness;
- stale/insufficient evidence is an honest Q result, not permission to expose a generic refresh mechanism.

Permission floor: `market.read`.

## 4.3 Expected Economics current meaning

| Candidate operation | Client | Class | Permission | Admission |
|---|---|---|---|---|
| `ListExpectedEconomics` | both | Q | `economics.read` | **ADMIT** |
| `GetExpectedEconomics` | both | Q | `economics.read` | **ADMIT** |
| generic mutable `Profitability` blob | both | Q/C | — | **REJECT** |

Expected Economics remains L0 meaning composed from material evidence such as current marketplace context, Cost Basis, Expected Tax, expected marketplace fee, expected seller shipping and material promotion/discount evidence.

The response must preserve which components are known/modeled/derived/unknown/unavailable/partial/stale. It cannot produce plausible precision from missing material components.

Collection filtering/pagination is justified by commercial-intelligence needs such as low/negative margin or insufficient/stale evidence; exact spelling is later B2 schema work.

## 4.4 Stateless price-scenario evaluation

| Candidate operation | Consumer / use | Client | Class | Permission | Idempotency | Admission |
|---|---|---|---|---|---|---|
| `EvaluatePriceScenario` | operator/agent asks Economics to evaluate a hypothetical candidate price/context | both | C — stateless, side-effect-free owner capability | `economics.read` | not required | **ADMIT** |
| create persistent `Simulation` resource for every evaluation | both | C | — | — | **REJECT BASELINE** |
| persistent `Recommendation` resource | both | C | — | — | **REJECT BASELINE** |

Binding decisions:

- hypothetical analysis is an Economics capability, not current-owner Q truth;
- scenario input may state legitimate hypothetical variables such as candidate price/context but cannot impersonate authoritative evidence by supplying arbitrary fake fee/tax/cost/competitor facts;
- no SimulationID is created merely for history;
- no Recommendation business authority is created; Economics may return conclusions/ranges/reasons when evidence supports them;
- when a consequential PriceIntent is later created, Economics must be able to preserve/re-establish the material decision-time L0 basis required for future expected↔order↔realized explanation without trusting client-computed economic evidence.

## 4.5 Sale Economics and economic lineage

| Candidate operation | Consumer / use | Client | Class | Permission | Admission |
|---|---|---|---|---|---|
| `ListSaleEconomics` | operator/manager/agent scans sales profitability/reconciliation | both | Q | `economics.read` | **ADMIT** |
| `GetSaleEconomics` | inspect one Sale's L0/L1/L2 lineage and reconciliation | both | Q | `economics.read` | **ADMIT** |
| `GetEconomicPerformanceSummary` | period/portfolio commercial-intelligence aggregate | both | Q | `economics.read` | **ADMIT** |
| separate mutable `OrderEconomics` / `RealizedEconomics` CRUD | both | Q/C | — | **REJECT** |
| universal `Reconciliation` resource with generic ID | both | Q/C | — | **REJECT** |

`SaleEconomics` is an Economics-owned read surface keyed by the source-qualified Sale context, not a single mutable row that overwrites history. It preserves distinct rungs proportionately:

```text
expected economic basis (L0 reference when material)
order economics (L1)
realized/settlement economics (L2)
reconciliation
  R1 expected ↔ order
  R2 order ↔ realized
```

Later refund/reversal/settlement evidence appends/reinterprets realized lineage; it does not erase an earlier approval/release/order fact merely because the final economic result changed.

Performance summaries are derived Qs with explicit period/scope/coverage. Partial reconciliation or missing material evidence must remain visible; a period aggregate cannot present itself as complete merely because an arithmetic value is available.

R3 bank-side reconciliation remains deferred until an accepted bank source exists.

## 4.6 Commercial Economics policy

| Candidate operation | Client | Class | Permission | Idempotency | Concurrency | Admission |
|---|---|---|---|---|---|---|
| `GetCommercialPolicy` | both | Q | `economics.read` | — | — | **ADMIT** |
| `UpdateCommercialPolicy` | human | C/resource update | `economics.policy.manage` | structural desired-state | **required where lost-update material** | **ADMIT** |
| generic Policy CRUD | both | C | — | — | — | **REJECT** |
| rule/expression/condition DSL editor | human | C | — | — | — | **REJECT** |

Commercial Economics owns only its policy meaning, such as material Cost Basis policy/selection, margin floors, price boundaries and economic approval-trigger thresholds. Governance does not own those thresholds merely because they may cause an approval requirement.

Policy update is a current-state lost-update problem, not a duplicate-create problem; concurrency is therefore primary and generic idempotency key is not baseline.

## 4.7 Economic Attribution / exception reconciliation

| Candidate operation | Consumer / use | Client | Class | Permission | Admission |
|---|---|---|---|---|---|
| `ListEconomicAttributions` | inspect exact/partial/ambiguous/unresolved economic attribution work | both | Q | `economics.read` | **ADMIT** |
| `GetEconomicAttribution` | inspect one persistent attribution + provenance | both | Q | `economics.read` | **ADMIT** |
| `ResolveEconomicAttribution` | human resolves a materially ambiguous/unresolved economic meaning/scope | human baseline | C | `economics.reconcile` | **ADMIT** |
| generic `MarkReconciled` | client sets status without owner semantics | both | C | — | **REJECT** |
| `ReconcileNow` | run internal reconciliation mechanism/job | both | C | — | **REJECT PRODUCT API** |

Economic Attribution is persistent MPC-owned Economics state because D2 already established its semantic identity. Resolution requires current-state concurrency where stale overwrite could replace a newer decision. Exact repeat may be structurally idempotent.

Human is the baseline client for explicit ambiguity resolution. Automatic sufficiently-known attribution may occur inside Economics without creating a client operation; broader automated exception resolution requires explicit policy/evidence before widening client class.

Permission floor: `economics.read`, `economics.policy.manage`, `economics.reconcile`.

## 4.8 Explicit Block 3 exclusions

Not admitted:

- generic MarketObservation CRUD;
- public market collection/refresh mechanism;
- generic market collector/scraper platform;
- persistent Simulation resource by default;
- Recommendation authority/resource;
- client-provided fake economic evidence/overrides;
- mutable all-in-one Profitability object;
- separate API authorities for expected/order/realized rungs merely for symmetry;
- universal ReconciliationID/resource;
- `ReconcileNow` runtime command;
- generic financial ledger;
- bank/Open Finance/R3 API before accepted source/consumer;
- campaign/discount authoring merely because observed promotion evidence affects economics;
- price actuation inside Economics.

## 4.9 Block 3 method outcome

**Parent structure:** `CURRENT D1/D2/D4-B4 STRUCTURE CONFIRMED`.

**B2 outcome:**

> **Use stateless Economics capabilities for hypothetical analysis; persist material economic authority/history only where future explanation/reconciliation requires it; preserve L0/L1/L2 and financial occurrences instead of overwriting them; expose Market Intelligence as interpreted source-qualified evidence; keep actual Price actuation exclusively in Offering.**

---

# 5. Exact next matrix work

**Block 4 — Controlled Action Governance + Marketplace Sales + Business-System Materialization.**

Derive the minimum Product API surface while preserving these fences:

- Governance owns authorization delegation/grant semantics and concrete Authorization Decision/context; it does not own business thresholds, domain validity, intent, execution or provider protocol;
- ordinary Permission remains distinct from domain disposition and Governance authorization;
- a concrete Authorization Decision must preserve exact authorized target scope and historical actor/authority context without mutating the underlying domain Intent;
- approval does not prove execution and cannot waive execution-time revalidation;
- Marketplace Sales owns marketplace-sale interpretation/context/correlation and transaction-specific Selling Entity attribution, not Materialization/Fulfillment/Economics/Post-Sale;
- provider Sale/Order identity remains Marketplace Installation + native key; no synthetic MPC sale alias exists by default;
- Business-System Materialization owns Business Order Intent + native-order convergence and Invoicing Intent + fiscal convergence;
- Sankhya TOP/NUNOTA/status/protocol remains D4 adapter-local and never enters Product API semantics;
- Party Resolution and Destination Realization are bounded Materialization prerequisites, not Customer/Address master CRUD;
- material consequential creates require idempotency unless stronger owner anchors prove structural safety;
- no generic workflow/advance-status/ERP-command API;
- no direct `/sankhya/orders` or `/sankhya/invoices` Product surface;
- material ambiguous/duplicate/native-correlation failures become explicit owner state/Work rather than blind retries.

The block must decide which authorization/grant operations have real Product clients, whether Business Order/Invoicing intents are directly client-created or owner-triggered from Sales/Fulfillment conditions, and which Materialization prerequisite resolutions need Product operations versus automatic/internal resolution.

Do not spell final HTTP paths/schemas until the operation inventory is coherent.

Implementation remains blocked until D9.