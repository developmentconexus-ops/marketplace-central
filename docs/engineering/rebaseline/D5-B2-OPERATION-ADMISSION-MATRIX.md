# D5-B2 — Operation Admission Matrix

> **Status:** OPEN / ACTIVE — Block 1 accepted in-stage; Block 2 next  
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

External engineering references may inform HTTP/mechanism choices but never override MPC semantic authority. In particular, RFC 9110's idempotent resource-replacement semantics, RFC 5789 PATCH semantics and resource-oriented API guidance support choosing resource update versus owner-specific capability according to actual semantics rather than CRUD uniformity.

---

# 2. Block 1 — Identity/Access + Marketplace Portfolio + Product & Channel Readiness — ACCEPTED IN-STAGE

## 2.1 D2 Identity / ordinary access substrate

The Product API does not become an IAM platform. It exposes only the access interactions real Product 1.0 clients need after authentication.

| Candidate operation | Consumer / use | Client | Class | Permission | Admission |
|---|---|---|---|---|---|
| `GetCurrentAccessContext` | React/automation establishes resolved MPC Principal, accessible Organization memberships and effective ordinary Permissions | both | Q | authenticated | **ADMIT** |
| `ListOrganizationMembers` | organization administrator; Work assignee selection where needed | human | Q | `access.read` | **ADMIT** |
| `ListAccessRoles` | administrator discovers product-defined access bundles assignable in MPC | human | Q | `access.read` | **ADMIT** |
| `AssignAccessRole` | owner/admin grants one existing product-defined role to one Membership | human | C | `access.manage` | **ADMIT** |
| `RevokeAccessRole` | owner/admin revokes one existing role assignment | human | C | `access.manage` | **ADMIT** |
| invitation/onboarding by email | SaaS/self-service identity onboarding | human | C | — | **DEFER** |
| create/delete Organization | tenant/SaaS provisioning | human | C | — | **DEFER** |
| create Keycloak/OAuth client or service account | authentication infrastructure provisioning | human | — | — | **REJECT FROM PRODUCT API** |
| issue MPC API key/token | duplicate credential authority | both | — | — | **REJECT** |

### Current access context

Do not create a fragmented `/me`, `/me/roles`, `/me/permissions`, `/session` family merely for convention. One D2-owned Q may expose the smallest resolved access context required by Product clients:

```text
authenticated Principal
+ accessible Organization memberships
+ effective ordinary MPC Permissions
```

The external IdP proves authentication. It does not need to carry MPC Organization or MPC business Permission as authority inside the token.

### Access-role mutations

`AssignAccessRole` and `RevokeAccessRole` operate on one explicit Membership + product-defined AccessRole relation rather than replacing an entire role set by default.

Where the semantic relation is unique, repeating the same assignment/revocation may be structurally idempotent and therefore can qualify for the D5-B1 operation-local idempotency-key exemption. Stale whole-set replacement is not introduced by convenience.

No custom role designer, generic ACL/ReBAC engine or IdP-role-as-business-authority is admitted.

### Initial Permission floor

The operation inventory currently justifies only:

- `access.read`
- `access.manage`

This is ordinary access only; it never grants business disposition or Governance authority.

---

## 2.2 Marketplace Portfolio

Marketplace Portfolio owns marketplace participation/configuration. Provider authentication remains D4 protocol.

| Candidate operation | Consumer / use | Client | Class | Permission | Admission |
|---|---|---|---|---|---|
| `ListMarketplaceInstallations` | human/automation discovers Organization marketplace participation available for operations | both | Q | `portfolio.read` | **ADMIT** |
| `GetMarketplaceInstallation` | operational context, authoring, diagnostics and source qualification | both | Q | `portfolio.read` | **ADMIT** |
| `CreateMarketplaceInstallation` | owner/admin establishes a new MPC marketplace participation identity/configuration | human | C | `portfolio.manage` | **ADMIT** |
| `UpdateMarketplaceInstallationConfiguration` | update Portfolio-owned participation/configuration such as eligible Selling Entity participation | human | C | `portfolio.manage` | **ADMIT** |
| `DeactivateMarketplaceInstallation` | owner/admin stops MPC participation without pretending provider account deletion | human | C | `portfolio.manage` | **ADMIT** |
| `ListSellingEntities` | configuration, attribution display and eligible-participation selection | both | Q | `portfolio.read` | **ADMIT** |
| arbitrary Selling Entity create/edit | multi-legal-entity operational administration beyond first proven need | human | C | — | **DEFER** |
| provider OAuth begin/callback/refresh | establish/maintain Mercado Livre credential protocol | human/provider | — | — | **NOT PRODUCT API — D4 TECHNICAL SURFACE** |
| generic provider/integration catalog | speculative integration-platform discovery | both | Q | — | **REJECT** |

### Creation idempotency

`CreateMarketplaceInstallation` is consequential and can mint duplicate MPC identities after a lost response. No safe natural duplicate predicate is assumed. **Client idempotency key is mandatory** under D5-B1 unless a later concrete contract proves a stronger structural anchor.

### Configuration concurrency

`UpdateMarketplaceInstallationConfiguration` may overwrite a newer participation/configuration decision if based on stale client state. The admitted contract therefore requires an **opaque MPC concurrency/precondition token** when changing state for which lost-update would be material.

Provider-native version tokens remain adapter-local unless the Portfolio-owned semantic precondition genuinely needs qualified exposure.

### Deactivation

Deactivation is an MPC Portfolio lifecycle operation, not provider account deletion. Repeating an already-effective deactivation may be structurally idempotent. The operation must still preserve actor/history and must not erase the Installation identity or source-qualified historical references.

### Permission floor

- `portfolio.read`
- `portfolio.manage`

---

## 2.3 Product & Channel Readiness

There is **no Product/PIM CRUD**. Product remains externally authoritative and source-qualified. The Product API exposes source Product facts only through the marketplace-operating Readiness use case.

| Candidate operation | Consumer / use | Client | Class | Permission | Admission |
|---|---|---|---|---|---|
| `SearchSourceProductsForMarketplace` | operator/agent finds a source Product to prepare, link or publish in a specific marketplace context | both | Q | `readiness.read` | **ADMIT** |
| `GetProductChannelReadiness` | authoring/operations determines correspondence/readiness/missing/conflicting state | both | Q | `readiness.read` | **ADMIT** |
| `GetPublicationRequirements` | ListingIntent editor/automation needs current applicable publication requirement contract and provenance | both | Q | `readiness.read` | **ADMIT** |
| `ResolveProductChannelCorrespondence` | human/automation explicitly establishes/replaces the Product↔channel correspondence decision | both | C | `readiness.manage` | **ADMIT** |
| `ClearProductChannelCorrespondence` | human/automation explicitly removes a current correspondence when semantics allow it | both | C | `readiness.manage` | **ADMIT — SAME LIFECYCLE** |
| `RefreshReadiness` / `RecalculateNow` | force an internal acquisition/evaluation mechanism | both | C | — | **REJECT PRODUCT API** |
| `SyncProducts` | source-acquisition/runtime mechanism | machine | C | — | **REJECT PRODUCT API** |
| `CreateProduct` / `UpdateProduct` | Product/PIM mastery | both | C | — | **REJECT** |
| bulk correspondence mutation | mass operation without a proven workflow/member semantics | both | C | — | **DEFER** |

### Marketplace-context Product search

`SearchSourceProductsForMarketplace` returns source-qualified Product candidates and only facts needed for Readiness/authoring use, with honest provenance/freshness. It never creates an MPC Product identity/master.

The subject is proportionately:

```text
Organization
+ Marketplace Installation
+ SourceInstance-qualified Product candidate
```

Real product population size justifies search/pagination here. Exact query fields and cursor/filter spelling remain later B2 schema/path work; common business identifiers such as source native key, seller-SKU/GTIN/name evidence may be searchable only where the source/Readiness contract can establish them honestly.

### Readiness versus publication requirements

`GetProductChannelReadiness` and `GetPublicationRequirements` remain separate Q semantics:

- readiness answers whether the source-qualified Product/channel relation is sufficiently known/ready and what source-level conditions are missing/conflicting/unsupported;
- publication requirements answer which current marketplace/category/product-type requirements the authoring client must satisfy, including material requirement identity/revision/provenance.

Neither operation makes Readiness own ListingIntent overrides or draft dispatchability.

### Correspondence concurrency and automation safety

Correspondence is a Readiness-owned semantic decision, not a matching formula. Replacing/clearing a current correspondence can conflict with a newer human decision.

The mutation therefore requires a current-state precondition/concurrency control whenever replacing material current meaning. An automation Principal cannot silently supersede a standing human decision from stale state; D2 historical attribution and explicit supersession semantics remain binding.

Where repeating the exact same correspondence resolution against the same current semantic subject is provably harmless, structural idempotency may exempt the operation from a generic idempotency key; concurrency/precondition still protects stale replacement.

### Permission floor

- `readiness.read`
- `readiness.manage`

---

## 2.4 Block 1 result

Accepted Product API surface so far is intentionally small:

```text
D2 access
  → current access context
  → minimal membership/role-assignment administration

Portfolio
  → Marketplace Installation lifecycle/configuration
  → Selling Entity discovery
  → provider auth remains D4 protocol

Readiness
  → marketplace-context source Product discovery
  → readiness
  → publication requirements
  → explicit correspondence lifecycle
```

Explicitly still absent:

- generic user/IAM platform;
- Organization SaaS provisioning;
- MPC API-key/application platform;
- generic Integration CRUD/catalog;
- Product/PIM CRUD;
- provider sync/refresh commands;
- generic bulk.

**Method outcome for Block 1:** `CURRENT PARENT STRUCTURE CONFIRMED; ADMIT ONLY CONSUMER-PROVEN OPERATIONS`.

---

# 3. Exact next matrix work

**Block 2 — Marketplace Offering Operations + Availability Control.**

Derive the minimum operations for:

- provider-actual Listing observation without making provider topology Product ontology;
- `ListingIntent` draft/create/edit lifecycle;
- draft dispatchability / freeze / consequential submission without `createListing = success` semantics;
- `PriceIntent` as Offering-owned price actuation meaning;
- Sellable Availability and Availability convergence;
- Inventory Source/Scope and allocation-policy configuration only where a real Product client needs it;
- the R1-G1 joint technical realization seam where one provider request needs owner-issued Offering + Availability inputs;
- multi-step/partial/asynchronous provider outcomes without introducing a generic LongRunningOperation/Mutation business owner.

Before admitting a public AvailabilityIntent creation endpoint, prove a real Product client needs to author that intent directly; internal automatic availability synchronization does not create a public API operation by symmetry.

Do not spell final HTTP paths/schemas until this block's ownership and operation inventory is coherent.

Implementation remains blocked until D9.