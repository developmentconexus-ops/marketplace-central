# D5-B2 — Operation Admission Matrix

> **Status:** OPEN / ACTIVE — Blocks 1 and 2 ACCEPTED IN-STAGE; Block 3 next  
> **Parent B2:** `D5-B2-PRODUCT-OPERATION-SURFACE.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + Decision Reconciliation Baseline + B2-A  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-18

## 1. Purpose

This matrix derives the smallest Product 1.0 operation surface from real clients/actors and accepted semantic owners before final path/schema spelling.

For every candidate operation the matrix tests, proportionately:

- real consumer/use;
- allowed client class: human, machine/automation/system, or both;
- exactly one accepted semantic owner or D2 substrate authority;
- Q / C / P class;
- explicit Organization scope;
- ordinary Permission, distinct from business disposition/Governance;
- canonical/source-qualified subject identity;
- read knowledge/freshness/provenance;
- consequential Intent/outcome/idempotency/precondition/concurrency;
- provider enrichment only for a named consumer/correctness property;
- pagination/filter/sort only for a real collection consumer;
- bulk only when a real workflow requires member-level semantics.

A candidate is rejected when its only justification is API symmetry, current code, provider endpoint shape, debug convenience or hypothetical future need.

---

# 2. Block 1 — Identity/Access + Marketplace Portfolio + Product & Channel Readiness — ACCEPTED IN-STAGE

## 2.1 D2 Identity / ordinary access substrate

The Product API does not become an IAM platform. It exposes only the access interactions real Product 1.0 clients need after authentication.

| Candidate operation | Consumer / use | Client | Class | Permission | Admission |
|---|---|---|---|---|---|
| `GetCurrentAccessContext` | Product client resolves MPC Principal, accessible Organizations and effective ordinary Permissions | both | Q | authenticated | **ADMIT** |
| `ListOrganizationMembers` | organization admin; assignee selection where needed | human | Q | `access.read` | **ADMIT** |
| `ListAccessRoles` | admin discovers product-defined role bundles | human | Q | `access.read` | **ADMIT** |
| `AssignAccessRole` | admin grants one existing role to one Membership | human | C | `access.manage` | **ADMIT** |
| `RevokeAccessRole` | admin revokes one role assignment | human | C | `access.manage` | **ADMIT** |
| invitation/onboarding by email | self-service/SaaS identity onboarding | human | C | — | **DEFER** |
| create/delete Organization | SaaS tenant provisioning | human | C | — | **DEFER** |
| create Keycloak/OAuth client or service account | authentication infrastructure provisioning | human | — | — | **REJECT FROM PRODUCT API** |
| issue MPC API key/token | duplicate credential authority | both | — | — | **REJECT** |

Binding decisions:

- one current access Q is preferred over fragmented `/me`/roles/permissions/session surfaces;
- IdP authentication does not carry MPC Organization/Permission authority by implication;
- role assignment/revocation is one-relation-at-a-time, avoiding unsafe whole-set replacement by default;
- repeated identical assignment/revocation may use structural idempotency where uniqueness makes duplication harmless;
- no custom-role designer, generic ACL/ReBAC engine or IdP-role-as-business-authority is admitted.

Permission floor:

- `access.read`
- `access.manage`

## 2.2 Marketplace Portfolio

| Candidate operation | Consumer / use | Client | Class | Permission | Admission |
|---|---|---|---|---|---|
| `ListMarketplaceInstallations` | discover Organization marketplace participation | both | Q | `portfolio.read` | **ADMIT** |
| `GetMarketplaceInstallation` | operational context, authoring and qualification | both | Q | `portfolio.read` | **ADMIT** |
| `CreateMarketplaceInstallation` | establish MPC marketplace participation identity/config | human | C | `portfolio.manage` | **ADMIT** |
| `UpdateMarketplaceInstallationConfiguration` | change Portfolio-owned participation/config | human | C | `portfolio.manage` | **ADMIT** |
| `DeactivateMarketplaceInstallation` | stop MPC participation without deleting provider account | human | C | `portfolio.manage` | **ADMIT** |
| `ListSellingEntities` | config/attribution/eligibility selection | both | Q | `portfolio.read` | **ADMIT** |
| arbitrary Selling Entity create/edit | broader multi-entity administration | human | C | — | **DEFER** |
| provider OAuth begin/callback/refresh | provider credential protocol | human/provider | — | — | **NOT PRODUCT API — D4** |
| generic provider/integration catalog | speculative integration platform | both | Q | — | **REJECT** |

Binding decisions:

- `CreateMarketplaceInstallation` requires client idempotency by default because a lost response could mint duplicate MPC identities;
- materially unsafe configuration updates require an opaque MPC concurrency/precondition token;
- deactivation preserves identity/history and may be structurally idempotent;
- provider OAuth remains protocol/integration surface, never Product business authority.

Permission floor:

- `portfolio.read`
- `portfolio.manage`

## 2.3 Product & Channel Readiness

| Candidate operation | Consumer / use | Client | Class | Permission | Admission |
|---|---|---|---|---|---|
| `SearchSourceProductsForMarketplace` | find source Product for marketplace preparation/publication | both | Q | `readiness.read` | **ADMIT** |
| `GetProductChannelReadiness` | inspect correspondence/readiness/missing/conflicting state | both | Q | `readiness.read` | **ADMIT** |
| `GetPublicationRequirements` | editor/automation needs applicable publication requirements + provenance | both | Q | `readiness.read` | **ADMIT** |
| `ResolveProductChannelCorrespondence` | explicitly establish/replace Product↔channel correspondence | both | C | `readiness.manage` | **ADMIT** |
| `ClearProductChannelCorrespondence` | explicitly remove current correspondence when valid | both | C | `readiness.manage` | **ADMIT — SAME LIFECYCLE** |
| `RefreshReadiness` / `RecalculateNow` | force internal evaluation mechanism | both | C | — | **REJECT PRODUCT API** |
| `SyncProducts` | source-acquisition/runtime mechanism | machine | C | — | **REJECT PRODUCT API** |
| `CreateProduct` / `UpdateProduct` | PIM mastery | both | C | — | **REJECT** |
| bulk correspondence mutation | no proven bulk workflow | both | C | — | **DEFER** |

Binding decisions:

- Product remains `SourceInstance + native key`; no MPC Product resource/master is created;
- source-product search is admitted only in marketplace-operating Readiness context and returns source-qualified evidence with honest provenance/freshness;
- readiness and publication-requirements are distinct Q meanings;
- correspondence replacement/clearing requires current-state concurrency when stale overwrite could defeat a newer human decision;
- automation never silently supersedes standing human correspondence from stale state;
- exact duplicate correspondence may use structural idempotency when harmless.

Permission floor:

- `readiness.read`
- `readiness.manage`

## 2.4 Block 1 result

```text
D2 access
  → current access context
  → minimal Membership/RoleAssignment administration

Portfolio
  → Marketplace Installation lifecycle/configuration
  → Selling Entity discovery
  → provider auth stays D4 protocol

Readiness
  → marketplace-context source Product discovery
  → readiness
  → publication requirements
  → explicit correspondence lifecycle
```

Explicitly absent: generic IAM platform, Organization SaaS provisioning, MPC API-key/application platform, generic Integration CRUD/catalog, Product/PIM CRUD, provider sync/refresh commands and speculative bulk.

**Method outcome:** `CURRENT PARENT STRUCTURE CONFIRMED; ADMIT ONLY CONSUMER-PROVEN OPERATIONS`.

---

# 3. Block 2 — Marketplace Offering Operations + Availability Control — ACCEPTED IN-STAGE

## 3.1 Governing Block 2 invariant

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

Accepted structure:

```text
Marketplace Listing      → Offering Q of provider-actual interpreted state
ListingIntent            → Offering create/edit authoring + tracking
PriceIntent              → Offering price-actuation intent + tracking
Sellable Availability    → Availability Q / convergence
Inventory Source/policy  → Availability configuration
D4/D7 effect mechanism   → may jointly serialize owner-issued meanings
```

## 3.2 Marketplace Listing observation

| Candidate operation | Consumer / use | Client | Class | Permission | Admission |
|---|---|---|---|---|---|
| `ListMarketplaceListings` | operator/automation observes marketplace listing population | both | Q | `offering.read` | **ADMIT** |
| `GetMarketplaceListing` | inspect one source-qualified Listing state/convergence | both | Q | `offering.read` | **ADMIT** |
| direct `CreateMarketplaceListing` | bypass ListingIntent/domain control | both | C | — | **REJECT** |
| direct `UpdateMarketplaceListing` | direct provider-shaped mutation | both | C | — | **REJECT** |
| direct `Pause/Reactivate/CloseListing` APIs | parallel mutation architecture | both | C | — | **REJECT AS SEPARATE BASELINE AUTHORITY** |

Binding decisions:

- Listing identity is Marketplace Installation + provider-native Listing key; no synthetic MPC mirror ID exists merely for convenience;
- listing collection pagination/filtering is justified by real population use; exact filter/cursor spelling remains later B2 work;
- read contracts expose Offering-owned interpretation plus source freshness/convergence, not raw Item/UserProduct/provider DTOs;
- provider-specific topology appears only as bounded source-qualified enrichment when needed for capability/blast-radius/correctness.

## 3.3 ListingIntent authoring/tracking resource

`ListingIntent` remains the single create/edit authoring identity accepted by D4-R1.

| Candidate operation | Consumer / use | Client | Class | Permission | Idempotency | Concurrency | Admission |
|---|---|---|---|---|---|---|---|
| `ListListingIntents` | review drafts/submitted intents/history | both | Q | `offering.read` | — | — | **ADMIT** |
| `GetListingIntent` | inspect desired state, dispatchability, lifecycle/convergence | both | Q | `offering.read` | — | — | **ADMIT** |
| `CreateListingIntentDraft` | begin one explicit create/edit authoring decision | both | C/resource create | `listing.manage` | **mandatory** | — | **ADMIT** |
| `UpdateListingIntentDraft` | declaratively change current draft desired state | both | C/resource update | `listing.manage` | structural where same revision/value | **required** | **ADMIT** |
| `DiscardListingIntentDraft` | explicitly abandon mutable draft | both | C | `listing.manage` | structural | current state/version | **ADMIT** |
| `SubmitListingIntent` | freeze/submit current draft into controlled consequential lifecycle | both | C | `listing.manage` | structural via Intent identity + current revision | **required** | **ADMIT** |
| separate `ValidateListingIntent` | duplicate Offering truth calculation | both | Q/C | — | — | — | **REJECT BASELINE** |
| separate `PreviewListing` | duplicate draft/dispatchability surface | both | Q | — | — | — | **REJECT BASELINE** |

Binding decisions:

1. **Create/edit share one architecture.** `target = none` means create; `target = existing source-qualified Listing` means edit.
2. **Draft updates are declarative desired-state changes**, not a generic mutation-command DSL.
3. Only `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` baseline value modes are admitted; no DERIVED/rule engine.
4. `GetListingIntent` may expose current resolved values, requirement revision/provenance, dispatchability/blockers, lifecycle/convergence and concurrency token; separate validate/preview operations are not admitted without a distinct consumer/failure class.
5. `CreateListingIntentDraft` requires idempotency because a lost create response could mint duplicate Business Intents.
6. Draft update/submit requires concurrency so stale automation cannot overwrite or submit against a newer human-authored decision.
7. Standing human override/supersession safety from D2/D4-R1 remains binding.

## 3.4 SubmitListingIntent is the consequential client boundary

The Product client asks Offering to submit the current authoring decision; it does not invoke technical freeze/authorize/execute/reconcile stages independently.

Conceptual lifecycle:

```text
mutable draft
  ↓ submit current revision
frozen/submitted Offering intent
  ↓ owner disposition / Governance when required
execution-time revalidation
  ↓
D4/D7 realization
  ↓
authoritative reread
  ↓
Offering-specific convergence/divergence/ambiguity
```

Binding distinctions survive:

```text
submitted != authorized
accepted != completed
completed != externally applied
externally applied != converged
ambiguous != safe-to-retry
```

No boolean `createListing = success` contract is admitted.

### No generic LongRunningOperation

A generic `Operation` resource is rejected for Product 1.0 baseline because the durable owner-local Intents already provide the correct pollable business tracking identity.

Do not create duplicate tracking identities such as:

```text
ListingIntent ID
+ MPC Operation ID
+ provider task/operation ID
```

unless a future concrete technical operation lacks any legitimate owner-local resource and creates a separate real consumer need.

Provider task/import IDs remain source-qualified reconciliation evidence, not MPC business identity.

## 3.5 Listing lifecycle actions remain Offering intent semantics

Pause/reactivate/close/edit of an existing Listing are Offering-owned desired lifecycle meanings represented through the existing ListingIntent architecture when admitted by current owner/provider capability.

Do not create independent provider-shaped mutation authorities merely because an HTTP endpoint could be named `/pause` or `/close`.

Final HTTP spelling remains later B2 work and may use explicit owner-specific methods where resource CRUD would lie; semantic identity remains ListingIntent.

## 3.6 PriceIntent

Price actuation remains Offering authority and distinct from Commercial Economics calculation/simulation.

| Candidate operation | Consumer / use | Client | Class | Permission | Idempotency | Admission |
|---|---|---|---|---|---|---|
| `ListPriceIntents` | review price-actuation history/current outcomes | both | Q | `offering.read` | — | **ADMIT** |
| `GetPriceIntent` | inspect desired price + lifecycle/convergence | both | Q | `offering.read` | — | **ADMIT** |
| `CreatePriceIntent` | make one explicit desired-price actuation decision | both | C | `price.manage` | **mandatory** | **ADMIT** |
| direct `SetListingPrice` | bypass durable intent/control | both | C | — | — | **REJECT** |
| public mutable `PriceDraft` lifecycle | duplicates Economics simulation / adds unused lifecycle | both | C | — | — | **REJECT BASELINE** |
| withdraw/cancel pending PriceIntent | no proven workflow yet | both | C | — | — | **DEFER** |

Binding decisions:

- `PriceIntent` expresses the desired exact Money target; Offering does not own the economic formula that produced it;
- command semantics such as `decrease 5%`, `match competitor`, `increase R$ 20` belong to economic reasoning/client decision, not Offering actuation authority;
- no public PriceDraft is created because Commercial Economics is the proper place for simulation/candidate analysis;
- creation requires idempotency because duplicate intake would create duplicate durable price decisions;
- owner/provider preconditions and convergence remain distinct from successful intake.

## 3.7 Availability current meaning and convergence

D0 already establishes automatic normal-path availability synchronization when sufficiently known and policy-valid. Therefore internal `AvailabilityIntent` identity does **not** create a public authoring API by symmetry.

| Candidate operation | Consumer / use | Client | Class | Permission | Admission |
|---|---|---|---|---|---|
| `ListSellableAvailability` | portfolio/operations inspection across products/listings | both | Q | `availability.read` | **ADMIT** |
| `GetSellableAvailability` | inspect current desired/actual/convergence for one subject | both | Q | `availability.read` | **ADMIT** |
| `ListInventorySources` | discover recognized inventory pools/config | both | Q | `availability.read` | **ADMIT** |
| `GetInventorySource` | inspect one source/config | both | Q | `availability.read` | **ADMIT** |
| `CreateInventorySource` | establish one business-recognized Inventory Source | human | C | `availability.manage` | **ADMIT** |
| `UpdateInventorySource` | change Availability-owned source config | human | C | `availability.manage` | **ADMIT** |
| `DeactivateInventorySource` | stop source participation without erasing history | human | C | `availability.manage` | **ADMIT** |
| get effective allocation/scope policy | inspect current Availability-owned operating policy | both | Q | `availability.read` | **ADMIT** |
| update allocation/scope policy | change MPC-owned Availability policy | human | C | `availability.manage` | **ADMIT** |
| public `CreateAvailabilityIntent` | direct intent authoring absent real Product consumer | both | C | — | **REJECT BASELINE** |
| `SetAvailableQuantity` | direct stock-style mutation bypassing derivation/policy | both | C | — | **REJECT** |
| `SyncAvailability` / `RefreshAvailability` | runtime/acquisition mechanism | machine/both | C | — | **REJECT PRODUCT API** |

### Availability Q semantics

Where material, the read contract preserves distinct owner meaning such as:

```text
Sellable Availability desired value
provider actual availability evidence
knowledge / unknown / unavailable / partial
freshness/provenance
owner-specific convergence / pending / divergence / ambiguity
```

Exact wire names remain later B2 schema work.

### Configuration semantics

Inventory Source and allocation/scope policy remain separate meanings even though both are Availability-owned. Do not collapse them prematurely into one generic `AvailabilityConfig` object merely for API convenience.

Material configuration updates require concurrency where stale overwrite would be unsafe. Creation of a new Inventory Source requires idempotency unless a stronger unique semantic anchor is explicitly proven.

## 3.8 Joint Offering × Availability realization

D4-R1 R1-G1 remains binding.

When provider protocol requires Listing representation plus quantity in one physical request:

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

Product API consequences:

- `SubmitListingIntent` does not accept quantity as Offering-owned content;
- the client does not create a fake combined `PublishListingRequest { listing, price, stock }` authority;
- missing required Availability-owned value fails closed before external dispatch;
- an Offering intent may be valid/submitted while execution remains pending on another owner-issued prerequisite, provided later execution revalidates authorization/validity and does not treat submission as eternal permission;
- no cross-owner atomicity is claimed;
- partial/asynchronous provider outcomes remain owner-specific.

## 3.9 Multi-step / partial / blast-radius effects

One provider request or early `2xx/201/202` never becomes whole-operation convergence.

Where material, Product-facing owner state preserves enough distinction for:

- confirmed / rejected / pending / ambiguous / not-executed aspects or members;
- authoritative reread per material owner concern;
- intended target scope;
- authorized target scope;
- actual attempted/provider-affected outcome scope;
- no blind replay of already-confirmed or ambiguous members/steps.

Provider shared-resource blast radius is source-qualified evidence and cannot silently widen the business intent.

## 3.10 Permission floor

Block 2 currently justifies only:

- `offering.read`
- `listing.manage`
- `price.manage`
- `availability.read`
- `availability.manage`

`listing.manage` and `price.manage` remain separate because least privilege is materially useful: listing-content automation need not inherit commercial repricing authority.

## 3.11 Explicit Block 2 exclusions

Not admitted:

- giant Listing CRUD owning content + price + stock + fulfillment;
- direct provider Listing create/update operations;
- generic Mutation/Command/Action APIs;
- generic `LongRunningOperation` business resource;
- separate PublicationPreparation;
- separate generic validate/preview APIs duplicating Offering truth;
- direct `SetListingPrice`;
- public PriceDraft baseline;
- public `CreateAvailabilityIntent` without a real client authoring use;
- direct `SetAvailableQuantity`;
- public sync/refresh mechanism;
- provider wire/task identifiers as MPC business identity.

## 3.12 Block 2 method outcome

**Parent structure:** `CURRENT D1/D2/D3/D4-R1 STRUCTURE CONFIRMED`.

**B2 outcome:**

> **ADMIT owner-specific Listing/Intent/Availability resource interactions; reject giant Listing CRUD, direct provider-shaped mutations and generic asynchronous-operation abstractions.**

---

# 4. Exact next matrix work

**Block 3 — Market Intelligence + Commercial Economics.**

Derive the minimum Product API surface while preserving these accepted fences:

- Market Intelligence owns comparable-market observation, comparability, competitive position/change and evidence sufficiency;
- Commercial Economics owns Cost Basis, economic interpretation, simulation, L0 Expected Economics, L1 Order Economics, L2 realized/settlement economics, variance/calibration;
- Economics does not write marketplace price; `PriceIntent` stays Offering-owned;
- market evidence/source-specific richness remains source-qualified and admitted only for named consumer/correctness needs;
- simulation/recommendation must not silently become authoritative action or Listing mutation;
- expected/modelled/order/realized economics remain distinct evidence rungs;
- absent/partial/stale market or cost evidence must not produce plausible precision;
- no generic financial ledger or generic market-collector API;
- no campaign/discount authoring surface merely because observed promotions affect economics;
- no public refresh/collection command unless a real Product client requires an owner-level capability rather than a D4/D7 mechanism.

The block must decide whether price simulation is a stateless calculation, a durable owner-owned economic object, or a mix by consumer need — without creating speculative `Simulation` persistence merely for history. It must also decide the minimum read surfaces for expected/order/realized economics and reconciliation without merging them into one mutable profitability row.

Do not spell final HTTP paths/schemas until the operation inventory is coherent.

Implementation remains blocked until D9.