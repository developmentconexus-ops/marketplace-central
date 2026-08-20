# D2 — Identity / Tenant / Data Ownership

> **Status:** CLOSED / ACCEPTED AS A WHOLE — B1+B2 consolidated; independent challenge and final Global Coherence completed; explicitly ratified by operator  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authorities:** `D0-PRODUCT-SYSTEM-DEFINITION.md`, `D1-DOMAINS-BOUNDARIES.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-16  
> **Accepted:** 2026-08-16  
> **Whole-W4 Principal-access-eligibility confirmation incorporated:** 2026-08-19

## 1. Purpose and boundary

D2 defines the target identity model, Organization/tenant isolation semantics, persistent-state ownership and shared/domain value/knowledge/time representations required by accepted D0/D1 authority.

D2 does **not** choose D3 communication/events, D4 provider/ERP transport contracts and credentials, D5 HTTP/OpenAPI, D6 frontend topology, or D7 runtime/process/transaction/RLS/deployment mechanisms.

Current code, tables, migrations and pre-rebaseline ADR structures are evidence only unless the active authority path independently carries their meaning forward.

## 2. Governing invariants

1. **Canonical identity follows semantic authority.** MPC-owned concepts use MPC-owned identities; externally authoritative objects retain explicit source-qualified identity. Correspondence never collapses two authorities into one identity.
2. **Organization is the canonical MPC tenant/isolation root.** No duplicate canonical `Tenant` identity exists without independently proven meaning.
3. **One meaning, one write authority.** Canonical business state is written by the D1 authority that owns its semantics. References, projections and historical snapshots do not transfer current authority.
4. **Mechanism does not become authority.** Shared identity/access/value/runtime machinery stays bounded to its accepted responsibility.
5. **Unknown remains unknown.** Unknown, zero, empty, absent and not-applicable remain distinguishable whenever material to correctness.
6. **Historical claims remain explainable.** Material MPC decisions/intents/work/economic conclusions retain sufficient historical lineage without requiring universal event sourcing.
7. **Identity is stable.** MPC-owned canonical identifiers are opaque, non-business-semantic and non-reusable; changing attributes/correspondences does not change identity.
8. **Legacy persistence does not constrain the target.** The pre-rebaseline MPC database has no required business-history or compatibility value and is not migrated or archived for the target.
9. **Pre-rebaseline ADRs are not target authority by inheritance.** Still-needed meaning must be explicitly rehomed before legacy ADR files leave the active tree.

---

## 3. Canonical MPC identities — APPROVED

### 3.1 Identity representation rule

MPC-owned canonical identities use stable opaque identifiers. The identifier itself does not embed or derive from Organization, CNPJ, CODPROD, provider account/resource ID, status, timestamp or other mutable/business semantics.

Once an identifier has denoted one canonical entity/occurrence, it is not recycled to denote another. D2 does not choose UUIDv4/UUIDv7/ULID/database encoding; later realization may choose any encoding that preserves the invariant.

### 3.2 Organization

`Organization` is an MPC-owned canonical identity and the root of organization-scoped isolation.

- The first Product 1.0 proof may use one Organization, but Organization is not hard-coded as a singleton.
- No separate target `Tenant` identity exists for Product 1.0.
- Organization ownership of historical/canonical organization-owned state is stable; changing ownership is not an ordinary foreign-key edit.
- Cross-Organization references between organization-owned business state are denied by default unless a later material decision explicitly defines a legitimate cross-Organization relationship.

### 3.3 Marketplace Installation

`Marketplace Installation` is MPC-owned canonical identity representing one Organization's participation/configuration in a marketplace.

- Provider seller/account identity remains external.
- Credentials/auth connections are D4/D7 mechanisms/configuration, not Installation identity.
- One Organization may have multiple Marketplace Installations.

### 3.4 Selling Entity

`Selling Entity` is an Organization-scoped MPC-owned canonical identity for the legal/fiscal actor when that distinction matters to marketplace operation.

- **Marketplace Portfolio owns the Selling Entity registry/lifecycle** as organization-facing operational configuration.
- External legal/company records such as ERP-native company identities remain source-qualified external references; MPC does not become legal-entity master.
- Marketplace Portfolio owns installation↔eligible-Selling-Entity participation/configuration.
- Marketplace Sales owns transaction-specific Selling Entity attribution; downstream domains consume that attribution rather than re-infer it.
- If Selling Entity later gains an independent fiscal/compliance business lifecycle beyond Portfolio's charter, D1 reopens.

### 3.5 Inventory Source

`Inventory Source` is MPC-owned canonical identity for a business-recognized inventory pool/source used by Availability Control.

- Native ERP/WMS/provider stock locations remain external facts/references.
- Mapping may be 1:1 or composed when real business semantics require it.
- Inventory Source does not collapse into Fulfillment Node.

### 3.6 Fulfillment Node

`Fulfillment Node` is MPC-owned canonical identity for a business-recognized physical execution point/capability.

- It is not an ERP company/location code, provider warehouse identifier or Inventory Source by convenience.
- Initial deployment may map Inventory Source and Fulfillment Node to the same physical place without collapsing their meanings.

### 3.7 Source Instance

`SourceInstance` is the MPC-owned reference identity for one logical namespace of an externally authoritative business-system/source when Marketplace Installation is not already the correct qualifier.

```text
Source Instance
  + external object kind
  + native key
```

- Changing credentials/token/connection mechanics does not itself create a new SourceInstance.
- Pointing at a materially different authoritative namespace/environment does.
- Marketplace integrations use Marketplace Installation where it already supplies the correct namespace; no generic SourceInstance wrapper is added merely for uniformity.
- SourceInstance registry administration is operator/administrative integration setup, not a D1 business-domain workflow and carries no business decision authority.
- D4 owns how concrete source identity, capabilities, protocol, credentials and configuration are verified.

No generic integration/source/entity registry is introduced.

### 3.8 Principal

`Principal` is the single MPC-owned canonical actor identity used for accountable attribution in MPC state/history.

Principal semantics distinguish at least:

- **human** actor;
- **automation** actor executing an explicitly authorized/policy-owned automated path;
- **system** actor for bounded machine/system actions.

Exact enum/subtype encoding is implementation detail. Non-interactive automation/system Principals do not require OIDC binding, do not authenticate through the interactive-human path, and must not be used to impersonate a human.

When automation is delegated/authorized by a human or grant, history preserves the effective non-human Principal **and** the applicable authorization/delegation context rather than attributing the machine action directly to the human.

Any organization-owned action by a non-human Principal remains explicitly Organization-scoped and subject to the same domain disposition/Governance boundaries as human action.

---

## 4. External/source-qualified identities — APPROVED

Externally authoritative identifiers are references, not MPC-global identities. EAN/GTIN, seller SKU, provider IDs, ERP document numbers and native codes remain identifiers/evidence within their proper source namespace.

### 4.1 Business-system Product

Product master remains externally authoritative. MPC references Product by:

```text
authoritative business-system SourceInstance
+ native Product key
```

MPC does not create a synthetic Product master or mirror identity merely for normalization. Exact Sankhya native key/contract is reverified in D4.

### 4.2 Marketplace Listing / Variation

Provider Listing/Offer is external source-qualified identity scoped through Marketplace Installation plus provider-native resource key.

Provider Variation is external identity scoped to its parent Listing plus provider-native variation key when applicable.

MPC-owned listing/change intents, desired-state semantics and convergence remain Marketplace Offering Operations state; no synthetic MPC ListingID is required merely as a mirror alias.

### 4.3 Marketplace Sale / Provider Order

Provider Sale/Order is external source-qualified identity scoped by Marketplace Installation plus provider-native order key.

Marketplace Sales owns MPC interpretation/context/correlation and transaction-specific Selling Entity attribution. No synthetic MPC sale ID is introduced merely as an alias for one provider order.

### 4.4 Shipment / Delivery

Provider Shipment is external source-qualified identity scoped by Marketplace Installation plus provider-native shipment key.

- Shipment does not collapse into Order.
- Provider Pack remains provider-native when relevant; no universal MPC Pack entity is introduced.
- Delivery remains provider lifecycle/outcome unless later evidence proves an independently meaningful MPC entity.

### 4.5 Native financial movements

Provider/payment-native Payment, Refund, Fee, Adjustment, Settlement, Payout and equivalent financial movements remain external/source-qualified identities.

No synthetic MPC Payment/Refund/Settlement identity is created merely for normalization. Source account/instance participates in the namespace when required.

---

## 5. MPC-owned intents, resolution, economics, work and governance — APPROVED

### 5.0 Material domain-owned Business Intent

A durable domain-owned intent that may be authorized, cause an external side effect, require convergence/reconciliation, or participate in material historical attribution has a stable **domain-local MPC-owned identity**.

This includes, where material, Offering-owned Listing/Price Intents, Availability-owned Availability Intents, Materialization-owned Business Order/Invoicing Intents and durable Fulfillment routing/dispatch intents.

The target does **not** introduce a generic `Action`, `Mutation`, `Command` or universal BusinessIntent business owner. Each action-owning domain owns its own intent meaning/lifecycle. Ephemeral calculations/recommendations do not gain canonical identity merely because they exist.

### 5.1 Business Order Intent / Invoicing Intent

`Business Order Intent` and `Invoicing Intent` are MPC-owned canonical identities under Business-System Materialization.

- Native ERP/business-system order/fiscal/document results remain external/source-qualified.
- Intent→native-result correlation is explicit.
- MPC does not create duplicate `MPCOrder` / `MPCInvoice` aliases over external results.

### 5.2 OperationalStage

`OperationalStage` may exist as a normalized **derived projection** for operator UX, e.g. waiting payment → ERP order generated → separation → ready invoice → invoiced → shipped → delivered.

It is not external truth, business authority, authorization, command, retry trigger or replacement for underlying facts/intents/outcomes.

### 5.3 Post-Sale Resolution

`Post-Sale Resolution` is MPC-owned canonical identity representing one explicitly scoped material post-sale resolution obligation.

- One Sale may have 0..N resolutions.
- Resolution scope may address lines/items/quantities for partial outcomes.
- Cancellation, Return and Refund are not one mutually-exclusive canonical enum; separate/concurrent/causally related consequences may coexist.
- Provider Claim/Case/Return/reverse-shipment/refund/payment resources and ERP/fiscal results remain external/source-qualified.
- Resolution closes only when applicable consequences for its explicit scope are sufficiently evidenced; one external terminal state does not imply global closure.

### 5.4 Economic Attribution

`Economic Attribution` is persistent MPC-owned semantic state under Commercial Economics.

It records what an externally authoritative financial movement means economically and to what MPC scope it is attributable, preserving provenance and whether interpretation is observed, modeled/configured or derived.

Attribution may be exact, partial, ambiguous or unresolved; missing correlation is not fabricated.

Its attributable scope is a legitimate **domain-local polymorphic subject** (for example sale, resolution, installation-level or period aggregate) and does not justify a platform-wide entity graph.

### 5.5 Economic Reconciliation

Commercial Economics owns persistent reconciliation semantics across the accepted lineage:

```text
R1: Expected Economics ↔ Order Economics
R2: Order Economics ↔ authoritative marketplace/payment settlement/realized evidence
R3: provider payout/settlement ↔ Bank Cash Receipt, only when an accepted bank source exists
```

Reconciliation may progress incrementally as evidence arrives. Month-end/period close is an aggregation/control view, not the sole moment reconciliation becomes true. ERP receivables/baixas remain ERP-native facts and do not become authority for marketplace settlement or bank cash.

No universal `ReconciliationID` is required merely for abstraction.

### 5.6 Operational Work

Every material actionable-work **obligation** has exactly one MPC-owned canonical Work identity.

Work owns responsibility, optional assignment, escalation, work-state and work-item closure. Its subject may be object, relationship, absence, population/coverage scope or another sufficiently explicit source-domain subject; an existing source row is not required.

This is a legitimate domain-local polymorphic subject, not permission to create a universal entity/reference graph.

A Work responsibility/queue **role** is operational-work semantics, not an ordinary access `AccessRole`. Any future mapping between them must be explicit and does not transfer authority between Work and identity/access.

Work references originating condition/evidence but does not replace or mutate source-domain truth. Closing Work alone does not declare the originating condition resolved.

Read-only attention indicators are projections. Responsibility/assignment/hold/dismissal/escalation or independent work lifecycle belongs to Work.

### 5.7 Authorization Decision

`Authorization Decision` is MPC-owned canonical decision occurrence under Controlled Action Governance.

It is distinct from Authorization Grant/Delegation, domain Business Intent, Operational Work and Execution Outcome. A Decision preserves sufficient authority context, decision Principal, outcome and exact authorized target-scope snapshot for the concrete case.

Later revocation/reapproval/rejection/invalidation does not erase prior decision history. Approval does not mutate Intent, prove execution, waive safety invariants or guarantee execution-time validity.

### 5.8 Authorization Grant / Delegation

Controlled Action Governance owns consequential-action authorization Grant/Delegation semantics: who/what may authorize which action class/scope.

Grant is distinct from ordinary product access permission, business disposition, concrete Authorization Decision and execution. Exact physical Grant ID/cardinality is not frozen merely for convenience; target state/provenance must prove applicable authority at decision time and support revocation without rewriting history.

---

## 6. Authentication and ordinary access — APPROVED

### 6.1 Identity/access state owner

The non-domain D2 identity/access substrate owns the canonical state/lifecycle for:

- Organization identity registry (without introducing SaaS tenant provisioning as a Product 1.0 requirement);
- Principal, including **current Principal access eligibility** as Principal-scoped revocable identity/access state;
- external OIDC identity bindings;
- Organization Membership;
- product-defined ordinary-access `AccessRole` / Permission definitions;
- RoleAssignment.

**Principal access eligibility confirmation.** Disabling/revoking a Principal's current Product-access eligibility blocks future Product access, including platform-scoped `/access-context`, without deleting Organization Membership/RoleAssignment or rewriting historical actor attribution. Exact lifecycle representation, storage, token/cache propagation and revocation mechanism remain later D7 realization; this confirmation creates no new Permission, AccessRole, domain or generic IAM lifecycle framework.

It does **not** own Marketplace Installation/Selling Entity/Inventory Source/Fulfillment Node business lifecycles, action disposition, consequential authorization grants/decisions or marketplace business policy.

### 6.2 Authentication boundary

MPC does **not** own end-user credentials or interactive authentication-session authority.

Interactive human AuthN is delegated through a standards-based OpenID Connect boundary to an external Identity Provider.

- Architecture depends on OIDC, not a vendor-specific user model.
- Mutable claims such as email/username are not canonical identity.
- External human identity binding uses stable OIDC `(issuer, subject)`.
- Keycloak remains the preferred first self-hosted candidate, while provider/deployment/realm topology are later technical choices.

### 6.3 External identity binding

One OIDC `(issuer, subject)` binding maps to at most one MPC Principal. A Principal may hold multiple bindings only through explicit identity administration, for example during IdP replacement.

Email/username matching never auto-merges Principals. Replacing the IdP must not rewrite historical Principal identity or past Work/Authorization/audit attribution.

### 6.4 Membership, RoleAssignment, AccessRole and Permission

The minimal ordinary-access kernel is:

```text
Principal
  → Organization Membership
  → RoleAssignment
  → product-defined AccessRole
  → product-defined Permission
```

- Membership and RoleAssignment are explicit durable/revocable MPC state.
- Current Principal access eligibility is separately Principal-scoped; it is a global Product-access gate and does not replace Organization-specific Membership/RoleAssignment.
- Revocation changes future access but does not rewrite historical actor attribution.
- AccessRoles are product-defined bundles; Permissions are stable ordinary-access capabilities.
- Organization-specific mutable ordinary-access state is Membership/RoleAssignment, not a custom role designer.
- Backend/business entry access checks consume semantic Permissions rather than hard-coded role-name branching.
- An ordinary Permission may allow a Principal to invoke/view a product capability; it does **not** itself prove business-action disposition, automation eligibility or consequential authorization authority. Those remain with the action-owning domain and Controlled Action Governance.
- Product 1.0 does not require custom-role design, nested groups, generic ACL/ReBAC, explicit-deny policy engines, OpenFGA/SpiceDB or a generic IAM platform.
- Exact permission→API-operation catalog may be completed in D5 without changing D2 ownership.

### 6.5 Identity/access substrate fence

The identity/access substrate is an important **non-domain D2 authority**, not a 13th D1 business domain.

It may answer identity/membership/role-holding questions such as:

- is Principal P currently eligible for Product access?;
- is Principal P a member of Organization O?;
- does Principal P hold AccessRole R?;
- which ordinary product Permissions follow from an assignment?

It MUST NOT answer a marketplace/business action's substantive permissibility, approval or execution validity.

- Action-owning domains own disposition/business validity.
- Controlled Action Governance owns consequential authorization grants/delegations and Authorization Decisions.

If the substrate begins making independent marketplace/business decisions, D1 reopens rather than silently growing a hidden domain.

---

## 7. Organization isolation semantics — APPROVED

Every organization-owned persistent business state or persisted external observation/evidence inside MPC belongs to exactly one Organization isolation scope.

Organization scope is explicit and is not inferred from Marketplace Installation, Selling Entity, external account, IdP organization, source key or process-global default.

Platform-owned definitions with no tenant-specific business meaning may remain platform-scoped, e.g. product AccessRole/Permission definitions or provider descriptors where justified.

A realization where isolation depends only on developers remembering `WHERE organization_id = ?` is invalid. Exact enforcement mechanism — schema constraints, transaction context, RLS/runtime checks or combination — remains D7/runtime design.

---

## 8. Persistence, references and lifecycle — APPROVED

### 8.1 Write authority follows D1 semantic authority

Each canonical MPC business fact has one semantic write authority: the D1 boundary that owns its meaning/lifecycle.

Consumers may retain explicit references to producer-owned identity/state, but they do not mutate or maintain a second current authority for that meaning.

Immutable historical snapshots are allowed when required to explain a past decision/event. A historical snapshot is historical context, never current producer authority.

### 8.2 Typed references; no universal entity graph

Known semantic relationships use explicit typed references/contracts owned by their domains, e.g. BusinessOrderIntent→Sale reference.

The target does not introduce a universal `{entity_type, entity_id}` system, universal entity table or generic relationship graph.

Domain-local polymorphic subjects are allowed only where the owning domain genuinely requires them. Operational Work and Economic Attribution are the currently accepted examples.

### 8.3 No shared mutable business model

Shared/kernel state is limited to genuine identity/value/knowledge primitives and technical mechanisms. No canonical shared mutable `SharedOrder`, `SharedProduct`, `SharedListing`, generic status entity or equivalent cross-domain business owner is introduced.

### 8.4 Persistent-state classes

D2 recognizes four semantic state classes; these do not imply four schemas or storage technologies:

1. **Canonical durable MPC state** — MPC-owned meaning that cannot be reconstructed merely by re-reading a provider.
2. **External observation/evidence** — provider/business-system truth preserved/translated only as required by MPC correctness/history, with source authority explicit and provider PII minimized.
3. **Derived projection/read model** — rebuildable state for reading/UX/analytics; never write authority.
4. **Technical ephemeral state** — cache, cursor, lease or processing state that gains no business authority merely by being persisted.

Pending business intents are canonical durable state; transport/outbox/queue mechanics remain D3/D7 concerns.

### 8.5 Historical lineage without universal event sourcing

Material Authorization Decisions, domain Business Intents, Post-Sale Resolutions, Economic Attribution/Reconciliation conclusions, Operational Work lifecycle and controlled-action lineage preserve enough historical state/evidence to explain material past decisions and outcomes.

The target does **not** require a universal event-sourced database/event store to satisfy this property.

### 8.6 Canonical identity lifetime

Canonical identities already referenced by material history are not silently deleted, reassigned or recycled such that past meaning changes.

Deactivation/terminal lifecycle is owned by the applicable domain and represented only where semantically needed. D2 does not require `deleted_at`, `is_active`, generic status or soft-delete columns on every table.

Historical occurrences such as Authorization Decisions are not rewritten into a different past decision; later revocation/supersession creates later state/history.

### 8.7 Clean baseline / no legacy-data migration

**Operator ruling:** the pre-rebaseline Marketplace Central database contains no business history or compatibility state that must be preserved.

Therefore:

- target persistence starts from a clean baseline;
- no legacy-data inventory, archival dump, semantic migration, backfill, dual-read or compatibility layer is required;
- legacy table contents do not constrain target identities, schemas, ownership or lifecycle;
- external facts needed by the target are reacquired from authoritative systems after integrations are configured;
- required MPC-owned configuration is recreated explicitly under the target model;
- possession of credentials/access needed to reconnect an external system is operational setup, not a requirement to migrate legacy integration-state tables;
- Git history is the archive for historical implementation structure; the old database is not target authority or required evidence.

This ruling concerns the **pre-rebaseline legacy database**. Once the new target operates, its own durable history follows the lineage rules above.

---

## 9. Value, knowledge, provenance and time — APPROVED

### 9.1 Money and rates

`Money` semantically carries:

```text
exact decimal amount
+ explicit currency
```

Rates/percentages used in canonical persisted state or consequential decisions are exact decimal dimensionless values unless a more specific domain concept says otherwise.

Binary floating point must not become authoritative persisted/decision state where approximation can affect monetary/tax/cost/pricing correctness. Non-authoritative derived analytics may use floating point when exactness loss cannot affect canonical state or consequential decisions.

Rounding is owned by the applicable domain rule or external provider/business-system contract when material; there is no universal `round(2)` rule.

### 9.2 Material quantities

Quantities used in canonical persisted state or consequential decisions preserve exactness where approximation could change business correctness.

Domains/sources determine whether quantity is integer/fractional and what unit/precision is meaningful. D2 does not create a generic Unit-of-Measure/conversion framework.

### 9.3 `Fact<T>` scope

ADR-034's `Fact<T>` primitive remains the current implementation/evidence anchor for knowledge shape during the rebaseline, but D2 owns its target scope.

Use `Fact<T>` where knowledge state, provenance or temporal validity are materially part of correctness; do not force it onto every value.

When material:

```text
Unknown ≠ Known(0) ≠ Empty ≠ Absent ≠ NotApplicable
```

`Fact<T>` combine/unknown-propagation semantics are adopted per consuming domain where proven necessary. Kernel `Map`/`Combine2` helpers are mechanism without cross-domain authority; no universal Fact algebra is mandated.

Before ADR-017/034 are removed from the active tree, every still-valid domain-judgment clause preserved by ADR-034 (including named-unknown components, opaque-stays-opaque, no silent cross-source fallback, lenient-ingestion boundaries and equivalent safety reasoning) must be explicitly rehomed in a new target Fact ADR and/or the owning D-stage artifacts.

### 9.4 Provenance

Material externally observed or derived facts preserve sufficient provenance to support the claim MPC makes from them.

A modeled fee, provider-observed fee and realized settlement movement may carry the same numeric amount while remaining semantically different evidence. Provenance is not a separate D1 business domain.

No universal `EvidenceID`, `ObservationID`, evidence graph or generic observation authority is introduced. Domains persist domain-specific evidence/observations as their semantics require.

### 9.5 Time

When material to correctness, distinct temporal meanings remain distinct, including:

- source/effective/event time;
- `observed_at` / acquisition time;
- `recorded_at` / MPC persistence time;
- decision time;
- externally authoritative deadline/window;
- relative-target anchor.

Internal instants represent an unambiguous instant rather than ambiguous local wall-clock values. Timezone/local-calendar semantics are represented only where the business/external rule actually depends on them.

Shared clock/timer/scheduler machinery remains D7 mechanism and does not own obligation/deadline meaning.

---

## 10. Correspondence and automation safety carried forward — APPROVED

### 10.1 Pre-dispatch correspondence consistency

Identity-bearing fields in an outbound provider write must be consistent with the accepted Product & Channel Readiness correspondence before the external side effect is dispatched.

A material mismatch fails closed **before** the external side effect; it is not intentionally sent and then repaired afterward. The concrete Mercado Livre fields/mappings, validation mechanism and external/API error representation belong to D4/D5/D7.

This preserves ADR-022's safety principle without carrying `SELLER_SKU == CODPROD` as a universal identity law.

### 10.2 Corroboration bar for unattended correspondence

For Product 1.0 automatic product↔channel correspondence, one matching identifier alone is not sufficient evidence for unattended establishment of correspondence.

Product & Channel Readiness owns the domain policy defining sufficient independent corroboration and material contradictory signals. D4 supplies current provider/source identifier evidence.

The old `CODPROD + EAN` formula is not the target identity law; the safety property is that unattended correspondence requires corroborating evidence and no material contradiction.

### 10.3 Automation does not silently override standing human decisions

A later automatic run must not silently reopen or reverse a standing human decision within the same semantic scope.

If a domain needs expiration, supersession or re-evaluation, that must be explicit domain semantics with historical attribution; automation does not acquire override authority by recurrence.

For Readiness matching specifically, machine actors use only the domain-owned automation-eligible path and do not masquerade through human/operator entry semantics. Other automated capabilities remain governed by their action-owning domain disposition and Controlled Action Governance where applicable.

---

## 11. D2 disposition of legacy ADRs — APPROVED

D2 adjudicates the following pre-rebaseline ADRs as historical evidence, not target structure:

| Legacy ADR | D2 target disposition |
|---|---|
| ADR-011 — shared Divergences ledger | **SUPERSEDED as target structure.** Divergence/business correctness stays with the originating D1 authority; actionable work lifecycle belongs to Operational Work. Honest evidence/history remains domain-owned. |
| ADR-012 — DIFAL table in legacy `pricing` | **SUPERSEDED as target structure.** Retain one-authority/no-fabricated-rule/provenance principles. Commercial Economics owns economic interpretation; external fiscal facts/rules remain external and D4 re-verifies contracts. |
| ADR-022 — `SELLER_SKU == CODPROD` | **SUPERSEDED as identity law.** Preserve pre-dispatch correspondence consistency (§10.1); D4 may establish a current Mercado Livre-specific mapping if evidence requires it. |
| ADR-028 — CODPROD+EAN auto-link | **SUPERSEDED as identity law.** Preserve unattended-corroboration and no-silent-human-override safety (§10.2–10.3); Readiness owns matching policy. |
| ADR-031 — `products_mirror` keep-absent | **SUPERSEDED as target mechanism.** Preserve honest absence/partial-observation semantics; no canonical Product mirror survives merely for this implementation. |

Legacy ADR text remains available through Git history after cleanup.

---

## 12. Legacy ADR transition / new target ADR baseline — APPROVED

**Operator direction:** pre-rebaseline ADRs are architecture history of the old system, not a catalog to carry into the rebuilt target.

### 12.1 Target authority during D0–D9

During the rebaseline, target authority is the accepted D-stage artifacts + stable `ARCHITECTURE.md` constraints + any explicitly created new target ADRs. No old ADR becomes target architecture merely because its legacy status once said accepted/binding.

### 12.2 Rehoming gates before legacy ADR deletion

1. **Reopened-ADR gate:** a legacy ADR awaiting D3/D4/D5/D7/D9 adjudication stays available as evidence until its owning stage has adjudicated the relevant meaning.
2. **Binding-constraint gate:** a currently carried safety/product constraint is not deleted from its only active authority location; it must first exist in `ARCHITECTURE.md`, an accepted D-stage artifact or a new target ADR.
3. **Fact-clause gate:** ADR-017/034 are not removed until the still-true domain-judgment clauses preserved by ADR-034 are explicitly rehomed (§9.3).
4. **Program gate:** ADR-035 remains until the D0–D9 rebaseline program closes because it governs the rebaseline authority transition itself.

### 12.3 New ADR series

New target ADRs use the next unused numeric sequence beginning at **ADR-036** (or later as consumed). Historical ADR numbers are never reused; deleting legacy files does not reset numbering because Git/history citations must remain unambiguous.

Only accepted target decisions that genuinely benefit from durable ADR treatment become new ADRs; D-stage artifacts are not mechanically exploded into one ADR per bullet.

### 12.4 Eventual cleanup

After the relevant stages have adjudicated/re-homed every still-needed constraint, legacy ADR files, the legacy registry and citation-archeology material may be removed from the active tree. Git history remains the archive.

This cleanup is not semantic migration and does not preserve old implementation authority.

---

## 13. Explicit defers

D2 intentionally does not decide:

- exact UUID/ULID/database encoding for MPC-owned identities;
- exact table names, schemas, indexes, FK strategy or repository/package layout;
- machine/service credential mechanism for non-interactive Principals;
- exact D3 sync/event/projection/outbox communication contracts;
- exact D4 Mercado Livre/Sankhya/payment DTOs, credentials, source capabilities, API/polling/webhook contracts;
- exact D5 endpoint/permission-operation/error catalog;
- D6 UI navigation, work inbox, approval screens or projection topology;
- D7 transaction/RLS/session context, process/worker/scheduler/queue/deployment topology;
- generic event sourcing, generic integration registry, universal entity/evidence graph, generic policy engine, generic IAM platform, generic Unit-of-Measure engine or generic workflow engine.

Later stages choose mechanisms only within accepted identity/authority/ownership meanings.

---

## 14. Proof strategy

A valid target realization must be able to falsify/prove at least these cases:

1. Same native external key in two SourceInstances remains unambiguous.
2. Same provider/native identifiers in different Organizations cannot collide/cross tenant isolation.
3. Changing email/username or IdP does not change Principal identity; one `(issuer,subject)` cannot map to two Principals.
4. Human and system/automation actions remain distinguishable; a machine action cannot be historically attributed as though the human executed it directly.
5. Authenticated Principal without MPC Membership/Permission is denied ordinary access.
6. Revoking Principal access eligibility, Membership, RoleAssignment or Grant changes future authority without rewriting prior Work/Authorization history.
7. Consumer persistence of producer refs/snapshots/projections cannot become current write authority.
8. Unknown source fact cannot silently become confident zero/default.
9. Monetary/rate/material-quantity canonical state cannot lose correctness through binary floating-point convenience.
10. Event/source time, observation time and decision time are not silently substituted.
11. Rebuilding a projection cannot mutate canonical source-domain state.
12. Product provider-write identity fields inconsistent with accepted correspondence are rejected before external side effect.
13. One identifier match alone cannot silently create unattended Product 1.0 correspondence.
14. Recurring automation cannot silently reverse a standing human decision in the same scope.
15. Starting from an empty clean target database is sufficient; required external facts/configuration can be reacquired/recreated without legacy-table dependency.
16. Removing an old ADR cannot remove the only active copy of a still-required invariant.
17. Disabling Principal access eligibility blocks `/access-context` and Organization-scoped Product calls even while historical Membership/RoleAssignment remains preserved.

---

## 15. Reopen triggers

Reopen only the implicated D2 decision when material evidence shows, for example:

- an MPC-owned concept gains an independently meaningful lifecycle that the assigned owner cannot hold without distortion;
- Selling Entity gains independent fiscal/compliance business authority beyond Marketplace Portfolio's charter;
- a second business-system/source namespace cannot be represented unambiguously by the minimal SourceInstance identity;
- multiple IdPs/identity lifecycle prove explicit binding administration insufficient;
- a real non-human actor cannot be represented/accounted for by Principal semantics without conflating human and machine attribution;
- a real cross-domain consumer requires shared Fact arithmetic not safely expressible per consuming domain;
- a legitimate relation cannot be represented without a broader reference model and domain-local polymorphism is insufficient;
- a real requirement proves Organization is not the correct isolation root or requires an independently meaningful Tenant identity;
- Identity/access substrate starts owning business action permissibility, consequential grants/decisions or marketplace policy;
- D7 cannot realize fail-closed Organization isolation within accepted D2 semantics;
- an actual pre-implementation durable-history requirement invalidates the clean-baseline ruling;
- legacy ADR cleanup would delete evidence/authority before its responsible D-stage adjudicates/re-homes it;
- new ADR numbering would collide with archived ADR identity.

Hypothetical future providers, SaaS packaging or preference changes alone do not reopen accepted D2 identity/ownership decisions.

---

## 16. Final Global Coherence + YAGNI / Overengineering / Future-Cost review — COMPLETED

**Outcome: CURRENT STRUCTURE CONFIRMED with bounded consolidation corrections. No B3, D0 reopen or D1 reopen is required.**

The review checked the complete B1+B2 target against D0/D1/stable architecture and the DevelopmentConexus Engineering Method.

### 16.1 Bounded corrections applied

1. **Material intent identity completeness.** D1 already assigns listing/price/availability/materialization intent meaning to action-owning domains. D2 now states that a durable material intent participating in authorization/external effect/convergence/history has stable domain-local identity (§5.0), preventing Governance/D3 from referencing an unnamed generic action.
2. **Identity/access ownership completeness.** The non-domain D2 substrate now explicitly owns Organization registry, Principal/bindings, Membership, ordinary AccessRole/Permission definitions and RoleAssignment (§6.1), while business identities stay with their D1 authorities.
3. **Role/authority vocabulary fence.** Ordinary `AccessRole`/Permission is separated from Operational Work responsibility role and from consequential authorization Grant/Decision (§5.6, §6.4), preventing one IAM concept from silently becoming work-routing or business-approval authority.
4. **Principal access-eligibility confirmation.** Whole-W4 review confirmed that current Principal access eligibility is revocable D2 identity/access state required to stop future Product access independently from Organization Membership while preserving historical attribution (§6.1–6.4).

These corrections make already-accepted authority explicit; they do not add a business domain or new product capability.

### 16.2 Duplicate/missing authority

**PASS.**

- Organization/Principal ordinary access → D2 identity/access substrate.
- Marketplace Installation/Selling Entity → Marketplace Portfolio.
- Inventory Source/Sellable Availability → Availability Control.
- Fulfillment Node/physical execution → Fulfillment Lifecycle.
- Product↔channel correspondence → Product & Channel Readiness.
- Listing/Price intents → Marketplace Offering Operations.
- Sale interpretation/entity attribution → Marketplace Sales.
- Business/Invoicing Intents → Business-System Materialization.
- Economics/Attribution/Reconciliation → Commercial Economics.
- Post-Sale obligation → Post-Sale Resolution.
- Consequential authorization Grant/Decision → Controlled Action Governance.
- Actionable-work lifecycle → Operational Work.
- External base facts/resources remain external/source-qualified.

No current business meaning has two write authorities and no material D2 identity/state is left ownerless.

### 16.3 Circular ownership / communication leakage

**PASS.** D2 adds typed references and ownership but does not choose cross-boundary communication. The accepted Materialization⇄Fulfillment business workflow cycle remains semantic only; D3 must realize it without private-code/write-authority cycles.

### 16.4 Trust / authorization layering

**PASS.**

```text
Principal access eligibility
  → global Product-access gate

Principal + Membership + AccessRole/Permission
  → ordinary Organization-scoped access capability

action-owning domain
  → business disposition / automation eligibility

Controlled Action Governance
  → consequential authorization Grant/Decision when required

execution/runtime
  → later revalidation + actual attempt/outcome
```

Possessing an ordinary Permission is never sufficient proof that an action is business-valid, automation-eligible or consequentially authorized. Human and non-human attribution remain distinguishable.

### 16.5 External-authority preservation

**PASS.** Product, provider Listing/Variation, Sale/Order, Shipment/Pack, native financial movements and ERP/fiscal results remain source-qualified external identities/facts. MPC owns correspondence, intents, interpretation, work, governance and reconciliation around them without becoming Product Master, provider truth or ERP replacement.

### 16.6 Tenant/isolation coherence

**PASS.** Organization is the only canonical business tenant root; Principal can participate through explicit Membership; external/source state persisted by MPC remains Organization-scoped. D7 still owes fail-closed enforcement, so D2 does not pretend a schema/runtime mechanism has already been chosen.

### 16.7 Persistence/history coherence

**PASS.** Clean reset of the pre-rebaseline database does not conflict with target historical-lineage requirements: old data is explicitly disposable; new target canonical decisions/intents/work/economic history becomes durable once the target operates. Projections remain rebuildable and never write authority.

### 16.8 Value/evidence/time coherence

**PASS.** Exact Money/rates/material quantities, bounded `Fact<T>`, provenance and distinct time meanings compose without requiring a universal evidence/time/unit framework. Unknown/absence semantics remain fail-honest.

### 16.9 YAGNI / overengineering

**PASS.** D2 explicitly rejects unsupported generic capability:

- no MDM/Product Master;
- no generic integration/source/entity registry;
- no universal entity/reference graph;
- no universal Evidence/Observation graph;
- no event-sourcing requirement;
- no generic IAM/ReBAC engine;
- no custom-role designer;
- no generic Unit-of-Measure engine;
- no generic policy/action/workflow engine;
- no speculative physical table/schema/process topology.

The retained abstractions all have a concrete current consumer or protect an accepted invariant.

### 16.10 Future-cost / seam review

**PASS.** D2 prepares only justified seams:

- opaque MPC identity decouples business identity from provider/ERP codes;
- SourceInstance/Marketplace Installation keep external namespaces explicit;
- OIDC binding decouples Principal from one IdP vendor;
- Principal supports human/automation/system attribution without a second actor graph;
- typed/domain-local references prevent a generic entity graph while allowing real polymorphic domains;
- exact values/provenance/time prevent irreversible correctness loss;
- clean legacy reset removes compatibility tax;
- legacy ADR rehoming preserves useful reasoning without carrying obsolete architecture.

### 16.11 Later-stage leakage

**PASS.** Exact event/outbox contracts remain D3; provider/ERP contracts/credentials remain D4; endpoint/permission/error catalog remains D5; UX topology remains D6; RLS/transactions/workers/queues/deployment/identity encoding remain D7/implementation as applicable.

### 16.12 Strongest counterexamples checked

- Same `15723` in two business-system namespaces does not collide.
- Same provider resource ID in two Organizations does not cross isolation.
- Keycloak/email changes do not rewrite Principal/history.
- Machine automation cannot impersonate a human actor.
- Disabled Principal cannot use `/access-context` as a bootstrap disclosure bypass.
- AccessRole permitting an approval screen does not itself authorize the approval.
- Price/availability/listing intent referenced by Governance has a domain-owned stable identity, not a generic Mutation row.
- Closing Work does not close the source condition by itself.
- Projection `OperationalStage` cannot trigger business writes by authority.
- Provider write identity mismatch fails before external side effect.
- Automatic matching cannot rest on one identifier or silently reverse a standing human resolution.
- Legacy DB deletion does not remove required target history because the operator explicitly declared it disposable.
- Legacy ADR deletion cannot remove the only active copy of a still-needed invariant because rehoming gates block that cleanup.

**Conclusion:** no material contradiction, missing authority, hidden God component, speculative framework or later-stage mechanism leak remains in D2.

---

## 17. D2 closure

D2 is **CLOSED / ACCEPTED AS A WHOLE** by explicit operator ratification on 2026-08-16, with the bounded Whole-W4 identity/access authority confirmation incorporated on 2026-08-19.

B1+B2 are consolidated above, independently challenged, and passed the final Global Coherence + YAGNI / Overengineering / Future-Cost review. Reopen only the specific D2 decision implicated by new material evidence under §15; do not re-run the stage for preference or hypothetical future scope.

**Next stage:** D3 — Communication / Events.

D3 owns the target communication/event/projection matrix and related semantics only within the D1 authority edges and D2 identity/ownership invariants. Product implementation remains blocked until D9 is accepted.