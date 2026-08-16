# D2 — Identity / Tenant / Data Ownership

> **Status:** OPEN / IN PROGRESS — operator-approved decisions below are binding within D2; D2 is not yet closed as a whole  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authorities:** `D0-PRODUCT-SYSTEM-DEFINITION.md`, `D1-DOMAINS-BOUNDARIES.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-16

## 1. Purpose and boundary

D2 defines the target identity model, Organization/tenant isolation semantics, persistent-state ownership, and shared/domain value/knowledge/time representations required by accepted D0/D1 authority.

D2 does **not** choose D3 communication/events, D4 provider/ERP transport contracts and credentials, D5 HTTP/OpenAPI, D6 frontend topology, or D7 runtime/process/transaction/RLS/deployment mechanisms.

Current code, tables, migrations and historical ADR structures remain evidence only unless the active authority path marks a constraint binding.

## 2. Governing invariants

1. **Canonical identity follows semantic authority.** MPC-owned concepts use MPC-owned identities; externally authoritative objects retain explicit source-qualified identity. A correspondence never collapses two authorities into one identity.
2. **Organization is the canonical MPC tenant/isolation root.** No duplicate canonical `Tenant` identity exists without independently proven meaning.
3. **One meaning, one write authority.** Canonical business state is written by the D1 boundary that owns its semantics. References, projections and historical snapshots do not transfer current authority.
4. **Mechanism does not become authority.** Shared identity/access/value/runtime machinery stays bounded to its accepted responsibility.
5. **Unknown remains unknown.** Unknown, zero, empty, absent and not-applicable remain distinguishable whenever the distinction is material to correctness.
6. **Historical claims are explainable from historical evidence.** Material MPC decisions/intents/work/economic conclusions retain sufficient lineage without requiring universal event sourcing.
7. **Legacy persistence does not constrain the target.** The pre-rebaseline MPC database has no required business-history or compatibility value and is not migrated into the target by default.

---

## 3. Canonical identity model — APPROVED

### 3.1 Identity authority rule

- MPC-owned concepts receive MPC-owned canonical identity.
- Externally authoritative objects retain external identity through an explicit source namespace/instance plus native key sufficient to make the reference unambiguous.
- External identifiers such as ERP product codes, EAN/GTIN, seller SKU, provider IDs or document numbers do not become MPC-global identity by convenience.
- No global MDM, universal entity table, generic entity graph or universal identity service is introduced.

### 3.2 Organization

`Organization` is an MPC-owned canonical identity and the root of organization-scoped isolation.

- The first Product 1.0 proof may use one Organization, but Organization is not hard-coded as a singleton.
- No separate target `Tenant` identity exists for Product 1.0.
- Organization ownership of historical/canonical organization-owned state is stable; changing ownership is not an ordinary foreign-key edit.
- Cross-Organization references between organization-owned business state are denied by default unless a later material decision explicitly defines a legitimate cross-Organization relationship.

### 3.3 Marketplace Installation

`Marketplace Installation` is MPC-owned canonical identity representing one Organization's participation/configuration in a marketplace.

- Provider seller/account identity remains external.
- Credentials/auth connections are mechanisms/integration configuration, not Installation identity.
- One Organization may have multiple Marketplace Installations.

### 3.4 Selling Entity

`Selling Entity` is an Organization-scoped MPC-owned canonical identity for the legal/fiscal actor when that distinction matters to marketplace operation.

- **Marketplace Portfolio owns the Selling Entity registry/lifecycle** as organization-facing operational configuration.
- External legal/company records such as ERP-native company identities remain source-qualified external references; MPC does not become legal-entity master.
- Marketplace Portfolio owns installation↔eligible-Selling-Entity participation/configuration.
- Marketplace Sales owns transaction-specific Selling Entity attribution; downstream domains consume that attribution rather than re-infer it.
- If Selling Entity later gains an independent fiscal/compliance business lifecycle beyond Portfolio's charter, D1 must be reopened.

### 3.5 Inventory Source

`Inventory Source` is MPC-owned canonical identity for a business-recognized inventory pool/source used by Availability Control.

- Native ERP/WMS/provider stock locations remain external facts/references.
- Mapping may be 1:1 or composed when real business semantics require it.
- Inventory Source does not collapse into Fulfillment Node.

### 3.6 Fulfillment Node

`Fulfillment Node` is MPC-owned canonical identity for a business-recognized physical execution point/capability.

- It is not an ERP company/location code, provider warehouse identifier or Inventory Source by convenience.
- Initial deployment may map Inventory Source and Fulfillment Node to the same physical place without collapsing their meanings.

### 3.7 Minimal Source Instance

D2 defines a **minimal stable Source Instance identity** sufficient to make source-qualified references unambiguous and Organization-attributable when Marketplace Installation is not already the appropriate qualifier.

Example semantic shape:

```text
Source Instance
  + external object kind
  + provider/native key
```

D2 does not create a generic integration/source/entity registry. D4 owns source contracts, capabilities, protocol, credentials and concrete integration configuration.

---

## 4. External/source-qualified identities — APPROVED

### 4.1 Business-system Product

The Product master remains externally authoritative. MPC references Product by source-qualified identity, conceptually:

```text
authoritative business-system Source Instance
+ native Product key
```

MPC does not create a synthetic Product master merely to mirror the source. Exact Sankhya native key/contract is reverified in D4.

### 4.2 Marketplace Listing / Variation

Provider Listing/Offer is external source-qualified identity scoped through the relevant Marketplace Installation plus provider-native resource key.

Provider Variation is external identity scoped to its parent Listing plus provider-native variation key when applicable.

MPC-owned listing/change intents, desired-state semantics and convergence remain owned by Marketplace Offering Operations; no synthetic MPC ListingID is required merely as a mirror alias.

### 4.3 Marketplace Sale / Provider Order

Provider Sale/Order is external source-qualified identity scoped by Marketplace Installation plus provider-native order key.

Marketplace Sales owns the MPC interpretation/context/correlation and transaction-specific Selling Entity attribution. No synthetic MPC sale ID is introduced merely as an alias for one provider order.

### 4.4 Shipment / Delivery

Provider Shipment is external source-qualified identity scoped by the relevant Installation plus provider-native shipment key.

- Shipment does not collapse into Order.
- Provider Pack remains provider-native when relevant; no universal MPC Pack entity is introduced.
- Delivery remains provider lifecycle/outcome unless later evidence proves an independently meaningful MPC entity.

### 4.5 Native financial movements

Provider/payment-native Payment, Refund, Fee, Adjustment, Settlement, Payout and equivalent financial movements remain external/source-qualified identities.

No synthetic MPC Payment/Refund/Settlement identity is created merely for normalization. Source account/instance is part of the identity namespace when required.

---

## 5. MPC-owned intents, resolution and economic/work/governance state — APPROVED

### 5.1 Business Order Intent / Invoicing Intent

`Business Order Intent` and `Invoicing Intent` are MPC-owned canonical identities.

- Native ERP/business-system order/fiscal/document results remain external/source-qualified.
- Intent→native-result correlation is explicit.
- MPC does not create duplicate `MPCOrder` / `MPCInvoice` aliases over external native results.

### 5.2 OperationalStage

`OperationalStage` may exist as an MPC normalized **derived projection** for operator UX, for example waiting payment → ERP order generated → separation → ready invoice → invoiced → shipped → delivered.

It is not external truth, business authority, authorization, command, retry trigger or replacement for underlying facts/intents/outcomes.

### 5.3 Post-Sale Resolution

`Post-Sale Resolution` is MPC-owned canonical identity representing one explicitly scoped material post-sale resolution obligation.

- One Sale may have 0..N resolutions.
- Resolution scope may address lines/items/quantities where partial outcomes exist.
- Cancellation, Return and Refund are not a single mutually-exclusive canonical enum; they may be separate/concurrent/causally related consequences.
- Provider Claim/Case/Return/reverse-shipment/refund/payment resources and ERP/fiscal results remain external/source-qualified.
- Resolution closes only when the applicable consequences for its explicit scope are sufficiently evidenced; one external terminal state does not imply global closure.

### 5.4 Economic Attribution

`Economic Attribution` is persistent MPC-owned semantic state under Commercial Economics.

It records what an externally authoritative financial movement means economically and to what MPC scope it is attributable while preserving provenance and whether the interpretation is observed, modeled/configured or derived.

Attribution may be exact, partial, ambiguous or unresolved; missing correlation is not fabricated.

### 5.5 Economic Reconciliation

Commercial Economics owns persistent reconciliation semantics across the accepted financial lineage:

```text
R1: Expected Economics ↔ Order Economics
R2: Order Economics ↔ authoritative marketplace/payment settlement/realized evidence
R3: provider payout/settlement ↔ Bank Cash Receipt, only when an accepted bank source exists
```

Reconciliation may progress incrementally as evidence arrives. Month-end/period close is an aggregation/control view, not the sole moment reconciliation becomes true. ERP receivables/baixas remain ERP-native facts and do not become authority for marketplace settlement or bank cash.

No universal `ReconciliationID` representation is required merely for abstraction; exact physical representation remains a D2 closure/data-model detail only if independently needed.

### 5.6 Operational Work

Every material actionable-work **obligation** has exactly one MPC-owned canonical Work identity.

Work owns responsibility, optional assignment, escalation, work-state and work-item closure. Its subject may be an object, relationship, absence, population/coverage scope or another sufficiently explicit source-domain subject; an existing source row is not required.

Work references originating condition/evidence but does not replace or mutate source-domain truth. Closing Work alone does not declare the originating condition resolved.

Read-only attention indicators are projections. Responsibility/assignment/hold/dismissal/escalation or independent work lifecycle belongs to Work.

### 5.7 Authorization Decision

`Authorization Decision` is MPC-owned canonical decision occurrence under Controlled Action Governance.

It is distinct from Authorization Grant/Delegation, Business Intent, Operational Work and Execution Outcome. A Decision preserves sufficient authority context, decision actor, outcome and exact authorized target-scope snapshot for the concrete case.

Later revocation/reapproval/rejection/invalidation does not erase prior decision history. Approval does not mutate Intent, prove execution, waive safety invariants or guarantee execution-time validity.

### 5.8 Authorization Grant / Delegation

Controlled Action Governance owns consequential-action authorization Grant/Delegation semantics: who/what may authorize which action class/scope.

Grant is distinct from ordinary product access permission, business disposition, concrete Authorization Decision and execution. Exact physical Grant ID/cardinality is not frozen merely for convenience; the target must preserve enough state/provenance to prove applicable authority at decision time and support revocation without rewriting history.

---

## 6. Authentication, Principal and ordinary access model — APPROVED

### 6.1 Authentication boundary

MPC does **not** own end-user credentials or authentication-session authority.

Interactive AuthN is delegated through a standards-based OpenID Connect boundary to an external Identity Provider.

- The architectural contract is OIDC, not a vendor-specific user model.
- Mutable claims such as email/username are not canonical identity.
- External identity binding uses the stable OIDC namespace `issuer + subject`.
- Keycloak is a preferred self-hosted implementation candidate, but provider/deployment/realm topology are not D2 authority and remain later technical choices.

### 6.2 MPC Principal

`Principal` is MPC-owned stable canonical identity for an actor participating in MPC state/history.

An external authenticated identity binds explicitly to Principal through `(issuer, subject)`. Replacing the IdP must not rewrite historical Principal identity or past Work/Authorization/audit attribution.

Principal may be platform-scoped; participation in an Organization is established through explicit Membership rather than embedding one immutable Organization directly into Principal identity.

### 6.3 Membership, RoleAssignment, Role and Permission

The minimal ordinary access kernel is:

```text
Principal
  → Organization Membership
  → RoleAssignment
  → product-defined Role
  → product-defined Permission
```

- Roles are product-defined bundles used for convenient assignment.
- Backend authorization checks semantic Permissions rather than hard-coded role-name branching.
- Product 1.0 does not require custom-role designer, nested groups, generic ACL/ReBAC graph, explicit-deny engine, OpenFGA/SpiceDB or a generic IAM platform.
- Exact permission catalog/operation mapping may be refined with D5 when the final API operation set exists, without changing D2 ownership.

### 6.4 Identity/access substrate fence

The D2 identity/access substrate is an important **non-domain D2 authority**, not a 13th D1 business domain.

It may answer identity/membership/role-holding questions such as:

- "Is Principal P a member of Organization O?"
- "Does Principal P hold Role R?"
- "What ordinary product Permissions follow from this assignment?"

It MUST NOT answer a marketplace/business action's substantive permissibility, approval or execution validity.

- Action-owning domains own disposition/business validity.
- Controlled Action Governance owns consequential authorization grants/delegations and Authorization Decisions.

If this substrate begins making independent marketplace/business decisions, D1 must be reopened rather than silently growing a hidden domain.

---

## 7. Organization isolation semantics — APPROVED

Every organization-owned persistent business state or persisted external observation/evidence inside MPC belongs to exactly one Organization isolation scope.

Organization scope is explicit and is not inferred from Marketplace Installation, Selling Entity, external account, IdP organization, source key or process-global default.

Platform-owned definitions with no tenant-specific business meaning may remain platform-scoped, for example product Permission/Role definitions or provider descriptors where justified.

A realization where isolation depends only on developers remembering `WHERE organization_id = ?` is invalid. Exact enforcement mechanism — schema constraints, transaction context, RLS/runtime checks or combination — remains D7/runtime design.

---

## 8. Persistence and data ownership — APPROVED

### 8.1 Write authority follows D1 semantic authority

Each canonical MPC business fact has one semantic write authority: the D1 boundary that owns its meaning/lifecycle.

Consumers may retain explicit references to producer-owned identity/state, but they do not mutate or maintain a second current authority for that meaning.

Immutable historical snapshots are allowed when required to explain a past decision/event (for example Governance's authorized target-scope snapshot). A historical snapshot is historical context, never current producer authority.

### 8.2 No shared mutable business model

Shared/kernel state is limited to genuine identity/value/knowledge primitives and technical mechanisms. The target does not introduce canonical shared mutable `SharedOrder`, `SharedProduct`, `SharedListing`, generic status entity or equivalent cross-domain business owner.

### 8.3 Persistent-state classes

D2 recognizes four semantic state classes; these do not imply four schemas or physical storage technologies:

1. **Canonical durable MPC state** — MPC-owned meaning that cannot be reconstructed merely by re-reading a provider.
2. **External observation/evidence** — provider/business-system truth preserved/translated only to the extent required by MPC correctness/history, with source authority explicit and provider PII minimized.
3. **Derived projection/read model** — rebuildable state for reading/UX/analytics; never write authority.
4. **Technical ephemeral state** — cache, cursor, lease or processing state that gains no business authority merely by being persisted.

Pending business intents are canonical durable state; transport/outbox/queue mechanics remain later D3/D7 concerns.

### 8.4 Historical lineage without universal event sourcing

Material Authorization Decisions, Business/Invoicing Intents, Post-Sale Resolutions, Economic Attribution/Reconciliation conclusions, Operational Work lifecycle and controlled-action lineage preserve enough historical state/evidence to explain material past decisions and outcomes.

The target does **not** require a universal event-sourced database or event store merely to satisfy this property.

### 8.5 Clean baseline / no legacy-data migration

**Operator ruling:** the current pre-rebaseline Marketplace Central database contains no business history or compatibility state that must be preserved.

Therefore:

- target persistence starts from a clean baseline;
- no legacy-data inventory, archival dump, semantic migration, backfill, dual-read or compatibility layer is required;
- existing legacy table contents SHALL NOT constrain target identities, schemas, ownership or lifecycle;
- external facts needed by the new target are reacquired from their authoritative systems after integrations are configured;
- required MPC-owned configuration is recreated explicitly under the target model rather than migrated merely because a legacy representation exists;
- possession of credentials/access needed to reconnect an external system is operational setup, not a requirement to migrate legacy integration-state tables;
- Git history remains the archive for historical implementation structure; the old database is not target authority or required historical evidence.

This ruling concerns the **pre-rebaseline legacy database**. Once the new target operates, its own accepted durable history is subject to the lineage rules above.

---

## 9. Shared/domain value, knowledge, provenance and time semantics — APPROVED

### 9.1 Money and rates

`Money` semantically carries:

```text
exact decimal amount
+ explicit currency
```

Rates/percentages used in canonical persisted state or consequential decisions are exact decimal dimensionless values unless a more specific domain concept says otherwise.

Binary floating point must not become authoritative persisted/decision state where its approximation can affect monetary/tax/cost/pricing correctness. Non-authoritative derived analytics may use floating point when exactness loss cannot affect canonical state or consequential decisions.

Rounding is owned by the applicable domain rule or external provider/business-system contract when material; there is no universal `round(2)` convenience rule.

### 9.2 `Fact<T>` scope

ADR-034's `Fact<T>` primitive is used where knowledge state, provenance or temporal validity are materially part of correctness. It is not forced onto every value.

When material:

```text
Unknown ≠ Known(0) ≠ Empty ≠ Absent ≠ NotApplicable
```

`Fact<T>` combine/unknown-propagation semantics are adopted per consuming domain where proven necessary. Kernel `Map`/`Combine2` helpers are mechanism without cross-domain authority; D2 mandates no universal Fact algebra.

### 9.3 Provenance

Material externally observed or derived facts preserve sufficient provenance to support the claim MPC makes from them.

A modeled fee, provider-observed fee and realized settlement movement may carry the same numeric amount while remaining semantically different evidence. Provenance is not a separate D1 business domain.

### 9.4 Time

When material to correctness, distinct temporal meanings remain distinct, including:

- source/effective/event time;
- `observed_at` / acquisition time;
- `recorded_at` / MPC persistence time;
- decision time;
- externally authoritative deadline/window;
- relative-target anchor.

Internal instants represent an unambiguous instant rather than an ambiguous local wall-clock value. Timezone/local-calendar semantics are represented only where the business/external rule actually depends on them.

Shared clock/timer/scheduler machinery remains D7 mechanism and does not own obligation/deadline meaning.

---

## 10. Explicit defers

D2 intentionally does not yet decide:

- exact UUID/ULID/database encoding for every MPC-owned identity unless the remaining D2 closure proves a material need;
- exact table names, schemas, indexes, FK strategy or repository/package layout merely from legacy structure;
- exact D3 sync/event/projection/outbox communication contracts;
- exact D4 Mercado Livre/Sankhya/payment DTOs, credentials, source capabilities, API/polling/webhook contracts;
- exact D5 endpoint/permission-operation catalog;
- D6 UI navigation, work inbox, approval screens or projection topology;
- D7 transaction/RLS/session context, process/worker/scheduler/queue/deployment topology;
- generic event sourcing, generic integration registry, universal evidence graph, generic policy engine, generic IAM platform or generic workflow engine.

Later stages may choose mechanisms only within the identity/authority/ownership meanings fixed by accepted D2 decisions.

---

## 11. Proof strategy for the accepted D2 set

The target realization must be able to falsify/prove at least these cases:

1. Same native external key in two Source Instances remains unambiguous.
2. Same provider/native identifiers in different Organizations cannot collide or cross tenant isolation.
3. User/Principal authenticates successfully but lacks MPC Membership/Permission → access is denied.
4. Revoking ordinary access or consequential authorization does not rewrite historical Authorization Decisions.
5. A consumer cannot become current write authority merely by persisting a producer reference/snapshot/projection.
6. Unknown source fact cannot silently become a confident zero/default.
7. Monetary/rate canonical state does not lose exactness through binary floating-point convenience.
8. Event/source time, observation time and decision time are not silently substituted for one another.
9. Rebuilding a projection cannot mutate canonical source-domain state.
10. Starting from an empty clean target database is sufficient; required external facts/configuration can be reacquired/recreated without semantic dependence on legacy tables.

---

## 12. Reopen triggers

Reopen only the implicated D2 decision when material evidence shows, for example:

- an MPC-owned concept gains an independently meaningful lifecycle that the assigned owner cannot hold without distortion;
- Selling Entity gains independent fiscal/compliance decision authority beyond Marketplace Portfolio's operational-configuration charter;
- a second business-system/source instance cannot be represented unambiguously by the minimal Source Instance identity;
- a real cross-domain consumer requires shared Fact arithmetic not safely expressible per consuming domain;
- a real requirement proves Organization is not the correct isolation root or requires an independently meaningful Tenant identity;
- Identity/access substrate starts owning business action permissibility, consequential grants/decisions or marketplace policy;
- D7 cannot realize fail-closed Organization isolation within the accepted D2 semantics;
- an actual post-cutover durable-history requirement is discovered that invalidates the operator's clean-baseline ruling before implementation begins.

Hypothetical future providers, SaaS packaging or preference changes alone do not reopen accepted D2 identity/ownership decisions.

---

## 13. Current D2 state / exact next action

D2 is **OPEN / IN PROGRESS**. The decisions above are operator-approved and binding within D2, but the stage is not yet accepted as a whole.

Next: **D2 Batch B2 — representation/persistence closure and Global Coherence preparation**. It must identify only the remaining D2 decisions necessary to make identity/value/persistent-state semantics implementation-ready, explicitly defer D3–D7 mechanics, and prepare the final D2 Global Coherence + YAGNI review.

Do not start D3 or product implementation until D2 is closed and the router advances.