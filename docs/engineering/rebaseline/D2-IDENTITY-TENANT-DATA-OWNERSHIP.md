# D2 — Identity / Tenant / Data Ownership

> **Status:** CLOSED / ACCEPTED AS A WHOLE — current consolidated authority  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authorities:** `D0-PRODUCT-SYSTEM-DEFINITION.md`, `D1-DOMAINS-BOUNDARIES.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Original acceptance:** 2026-08-16  
> **Later accepted amendments consolidated:** Principal-access eligibility + NOTIF-01 D2-R→R6 / AuthorizationRequest

## 1. Purpose and governing invariants

D2 owns target identity, Organization isolation semantics, persistent-state ownership and shared value/knowledge/time semantics. D2 does not choose D3 communication, D4 provider transport, D5 wire, D6 UX or D7 storage/process realization.

Binding invariants:

1. **Canonical identity follows semantic authority.** MPC-owned concepts use opaque MPC-owned identities; external objects remain source-qualified references.
2. **Organization is the canonical MPC isolation root.** No duplicate `Tenant` identity exists without independent meaning.
3. **One meaning, one write authority.** References, projections and historical snapshots never become current producer authority.
4. **Mechanism ≠ authority.** Identity/access/value/runtime mechanisms do not gain business decision authority.
5. **Unknown remains unknown.** Unknown, zero, empty, absent and not-applicable remain distinct when material.
6. **Material history remains explainable** without requiring universal event sourcing.
7. **MPC-owned IDs are stable, opaque, non-reusable and non-business-semantic.**
8. **External native identity is never promoted to MPC-global identity by convenience.**
9. **Git/history and legacy persistence are evidence, not target authority.** Still-valid meaning is rehomed before legacy artifacts retire.

---

## 2. Canonical MPC identities

### 2.1 Organization

`Organization` is MPC-owned and is the isolation root. A first deployment may use one Organization, but the model never hard-codes a singleton. Cross-Organization references between organization-owned state fail closed unless a later accepted relationship explicitly permits them.

### 2.2 Marketplace Installation

One MPC-owned Organization-scoped participation/configuration identity for one marketplace namespace. Provider account/seller identity and credentials remain external/D4 concerns. One Organization may have many Installations.

### 2.3 Selling Entity

Organization-scoped MPC identity for the legal/fiscal actor when material. Marketplace Portfolio owns registry/lifecycle and installation eligibility; Marketplace Sales owns transaction-specific attribution. ERP/legal records remain external.

### 2.4 Inventory Source / Fulfillment Node

Both are MPC-owned and distinct:

- `InventorySource` = business-recognized inventory pool/source under Availability;
- `FulfillmentNode` = business-recognized physical execution point/capability under Fulfillment.

They may map to the same physical place without collapsing semantics.

### 2.5 SourceInstance

MPC-owned namespace qualifier for an externally authoritative business-system/source when Marketplace Installation is not already the correct namespace.

```text
SourceInstance + external object kind + native key
```

Credential rotation does not create a new SourceInstance; moving to a materially different authoritative namespace does. No generic integration/entity registry is admitted.

### 2.6 Principal

One MPC-owned accountable actor identity supporting human, automation and bounded system actors. Non-human actors never impersonate humans. Historical attribution preserves the effective Principal and applicable authorization/delegation context.

Interactive human binding uses stable OIDC `(issuer, subject)`; email/username is presentation, not identity. One binding maps to at most one Principal. IdP replacement does not rewrite Principal history.

---

## 3. External/source-qualified identities

Externally authoritative objects remain qualified references:

- **Business-system Product:** `SourceInstance + native Product key`; no MPC Product master.
- **Marketplace Listing/Variation:** Marketplace Installation + provider-native key; no synthetic mirror ID merely for normalization.
- **Marketplace Sale/Order:** Marketplace Installation + provider-native sale/order key; Marketplace Sales owns interpretation, not source identity.
- **Shipment:** Marketplace Installation + provider-native shipment key; Shipment != Order; provider Pack remains provider-native.
- **Native financial movements:** Payment, Refund, Fee, Adjustment, Settlement, Payout and equivalents remain provider/payment-source identities.

Correspondence/linkage relates authorities; it never merges identities.

---

## 4. Domain-owned durable identities and historical state

A durable domain-owned intent that may be authorized, cause an external effect, require convergence/reconciliation or participate in material history has a stable domain-local identity. There is no generic `Action`, `Mutation`, `Command` or universal BusinessIntent owner.

Accepted durable MPC identities/occurrences include proportionately:

- Offering `ListingIntent` / `PriceIntent`;
- Materialization `BusinessOrderIntent` / `InvoicingIntent`;
- `PostSaleResolution`;
- `EconomicAttribution` and domain-local reconciliation state where persistent identity is actually needed;
- `Work`;
- `AuthorizationRequest`;
- `AuthorizationDecision`;
- bounded domain-specific fulfillment/availability intents when their lifecycle is material.

### 4.1 PostSaleResolution

One explicit material post-sale obligation/scope. One Sale may have many. Provider claim/return/refund/reverse-shipment/native fiscal results remain external. Resolution closes only when applicable consequences are sufficiently evidenced.

### 4.2 Economic Attribution / reconciliation

Economics owns persistent attribution meaning for external financial evidence. Attribution may be exact, partial, ambiguous or unresolved. Domain-local polymorphic subject is allowed; no universal entity graph or universal Reconciliation ID is created by abstraction alone.

### 4.3 Operational Work

Every material actionable-work obligation has one MPC-owned Work identity. Work owns responsibility, optional assignment, hold/resume/escalation/work lifecycle and evidence-backed closure of the Work item. It never becomes the source-domain truth or authorization authority; closing Work alone does not close the originating business condition.

### 4.4 AuthorizationDecision

Governance-owned immutable consequential decision occurrence. It is distinct from ordinary access, delegation/grant, Business Intent, AuthorizationRequest, Work and execution outcome. Later revocation/reapproval/execution invalidation never rewrites the historical decision.

### 4.5 Authorization delegation/grant

Governance owns who/what may authorize which action class/scope. Grant/delegation is distinct from ordinary Product Permission, business disposition, concrete Decision and execution.

---

## 5. AuthorizationRequest — canonical pre-decision authorization episode

`AuthorizationRequest` is MPC-owned canonical state under **Controlled Action Governance**, with one opaque stable `AuthorizationRequestID`.

### 5.1 Identity and target

A Request identifies one concrete authorization episode and is distinct from:

```text
governed Business Intent
AuthorizationDecision
Work
Notification
execution outcome
```

The Launch-V1 governed-target union remains closed:

```text
ListingIntent
PriceIntent
BusinessOrderIntent
InvoicingIntent
```

A later reauthorization episode receives a new Request ID even when target/revision is unchanged. Retry/recovery of the same semantic episode resolves to the same Request rather than minting duplicates.

### 5.2 Immutable authorization review/basis snapshot

Each Request retains the smallest typed immutable snapshot needed to explain **what materially was being authorized and under which relevant basis**. The action owner supplies source meaning; Governance owns retention for authorization history.

It is not a generic payload, current target projection or source-authority transfer. A coarse target ETag/revision may be evidence, but **is not the business validity oracle**: irrelevant target drift must not force needless reapproval while material governing drift may invalidate the episode.

### 5.3 Request revision vs material validity

AuthorizationRequest has owner-local revision for concurrent decision/stale-write protection. Before committing a decision Governance revalidates:

```text
Request is PENDING
+ Principal is currently allowed to decide
+ authorization episode remains materially valid enough to decide
```

Current source/business validity remains evaluated through the action-owning semantic boundary where needed.

### 5.4 Lifecycle

Closed baseline lifecycle:

```text
PENDING
DECIDED        // exactly 0..1 terminal AuthorizationDecision
INVALIDATED    // terminal without fabricated reject decision
```

No generic expiry state exists without a real governing rule. Invalidated Requests are not reopened; legitimate later authorization is a new Request. Optional `predecessor_authorization_request_id` may preserve explicit reauthorization lineage for the same governed action.

Requester/initiator Principal lineage is retained when materially available for accountability/result awareness; email/role/Permission is never requester identity.

### 5.5 Notification consequence

`AUTHORIZATION_ACTION_REQUIRED` points to the exact pending `AuthorizationRequest`. `AUTHORIZATION_DECISION_RESULT` remains requester-oriented toward the governed target; requester awareness never grants `governance.read` merely to inspect Governance history.

---

## 6. Personal Notification identity and owned state

`Notification` is one durable MPC-owned **personal-awareness occurrence** under the Personal Notifications supporting owner. It has one opaque stable `NotificationID`, one Organization and one exact historical human Principal recipient.

Notification is not the source occurrence, a capability token, Work, authorization, audit event or source-state mirror.

### 6.1 NotificationKind

Accepted stable Product keys:

```text
MARKETPLACE_INSTALLATION_ATTENTION
OFFERING_ASYNC_ACTION_RESULT
AVAILABILITY_ATTENTION
ECONOMIC_RECONCILIATION_ATTENTION
NEW_MARKETPLACE_SALE
SALE_ATTENTION
MATERIALIZATION_ATTENTION
FULFILLMENT_ACTIONABLE
FULFILLMENT_ATTENTION
SHIPMENT_EXCEPTION
POST_SALE_ATTENTION
WORK_ASSIGNMENT
AUTHORIZATION_ACTION_REQUIRED
AUTHORIZATION_DECISION_RESULT
```

The count is derived, not a protected platform constant. Historical Notifications keep the meaning of the kind under which they were created.

### 6.2 Current canonical Notification state

Conceptually, current accepted state includes:

```text
notification_id
organization_id
recipient_principal_id
kind
source_ref
source_occurrence_key
source_occurred_at
source_committed_at
subject_display_label
offering_async_result_outcome?      // F02 only
authorization_decision_outcome?     // F14 only
created_at
read_at?
archived_at?
superseded_at?
superseded_by_notification_id?
supersession_reason?
work_replacement_basis?
sale_attention_replacement_basis?
revision
```

`read_at` and `archived_at` are orthogonal. Read/archive/supersession changes Personal Notifications state only and never acknowledges/resolves/mutates the source.

### 6.3 Closed typed source references

Notification uses only the smallest already-meaningful typed source references required by accepted kinds; no universal `{entity_type, entity_id}` graph.

Current compatibility is conceptually:

```text
MARKETPLACE_INSTALLATION_ATTENTION   → MarketplaceInstallationRef
OFFERING_ASYNC_ACTION_RESULT         → ListingIntentRef | PriceIntentRef
AVAILABILITY_ATTENTION               → AvailabilityTargetRef
ECONOMIC_RECONCILIATION_ATTENTION    → EconomicAttributionRef | MarketplaceSaleRef
NEW_MARKETPLACE_SALE                 → MarketplaceSaleRef
SALE_ATTENTION                       → MarketplaceSaleRef
MATERIALIZATION_ATTENTION            → BusinessOrderIntentRef | InvoicingIntentRef
FULFILLMENT_ACTIONABLE               → FulfillmentExecutionRef
FULFILLMENT_ATTENTION                → FulfillmentExecutionRef
SHIPMENT_EXCEPTION                   → MarketplaceShipmentRef
POST_SALE_ATTENTION                  → PostSaleResolutionRef
WORK_ASSIGNMENT                      → WorkRef
AUTHORIZATION_ACTION_REQUIRED        → AuthorizationRequestRef
AUTHORIZATION_DECISION_RESULT        → AuthorizationTargetRef
```

A mismatched kind/source pair is invalid state. No synthetic Sale/Shipment/Availability/Reconciliation identity is introduced solely for Notification symmetry.

### 6.4 Source occurrence and deduplication

Each source owner supplies a stable owner-local `source_occurrence_key`. Personal-awareness semantic uniqueness is:

```text
Organization
+ recipient Principal
+ NotificationKind
+ typed source_ref
+ source_occurrence_key
```

Same occurrence replay/recovery does not mint duplicate awareness. Distinct episodes on the same subject may legitimately create distinct Notifications.

`source_occurred_at`, `source_committed_at` and Notification `created_at` remain distinct. Routing cutover uses the source-owner durable commit lineage, not provider webhook arrival or later transport delivery.

---

## 7. Notification audience and Organization routing

Audience semantics are the D1 closed strategies:

- **DIRECT_SOURCE:** exact source-owned human lineage;
- **OWNER_DERIVED:** source owner (currently Governance for action-required) resolves the exact current eligible human set;
- **ORG_ROUTED:** Personal Notifications uses explicit Organization routing state.

Permission/role names never choose recipients.

### 7.1 Natural route identity and revisions

```text
NotificationRouteKey = (organization_id, notification_kind)
```

No synthetic RoutingConfigID is required. ORG_ROUTED route history is immutable revision lineage so delayed/replayed occurrences use the route revision that was current when the source occurrence committed.

Each revision is exactly:

```text
CONFIGURED(one-or-more recipient bindings)
UNCONFIGURED
```

Configured-empty is invalid. Explicit unconfigure is a new revision, not deletion of prior history. A later reconfiguration never backfills occurrences committed while unconfigured or resurrects old bindings by convenience.

### 7.2 Recipient eligibility continuity

A route binding preserves exact Principal identity plus an opaque Organization/Principal access-eligibility continuity epoch sufficient to distinguish uninterrupted eligibility from revoke→re-enable.

An old binding does not silently reactivate after access/Membership continuity is broken; an explicit new routing decision is required. Exact physical epoch realization is D7.

### 7.3 Historical recipient stability

Existing Notifications never retarget when routing, Membership, assignment, display name or access later changes. Access loss blocks current source continuation and future materialization as applicable without rewriting historical ownership. Notification possession never grants source access.

---

## 8. Bounded Notification presentation/result/supersession

### 8.1 `subject_display_label`

A small immutable non-sensitive human label retained solely to distinguish Inbox items without N source rereads.

- source owner derives the business-meaningful safe label;
- Personal Notifications owns the retained historical snapshot;
- it is not current source truth, identity, navigation key, dedup key, authorization, routing evidence or automation input;
- later source presentation changes do not rewrite it;
- no provider DTO, buyer/customer PII, address, payment/fiscal payload, credentials or arbitrary JSON is retained for convenience.

Client-localized Product title/copy comes from NotificationKind; no generic template/payload engine is admitted.

### 8.2 Typed result atoms

Only two result-bearing families retain typed immutable occurrence outcomes:

```text
OFFERING_ASYNC_ACTION_RESULT:
  converged | rejected | ambiguous | divergent

AUTHORIZATION_DECISION_RESULT:
  authorize | reject
```

They are historical presentation/result meaning only, not current source state or execution authority.

### 8.3 Two bounded supersession rules

Personal Notifications may converge duplicate awareness only for the accepted cases:

```text
work_assignment_replacement
post_sale_precedence
```

A source alert can yield to the causally related richer Work-assignment Notification for the same human; Sale attention can yield to the same exact PostSaleResolution's richer Post-Sale attention. Correlation uses exact typed identities/occurrence lineage, never title/time similarity or a generic related-entity graph.

Late arrival may mark the less-specific Notification superseded through owner-local awareness state rather than deleting/reinterpreting source truth. Read/archive state remains orthogonal.

---

## 9. Identity/access substrate and ordinary access

The non-domain D2 identity/access substrate owns:

- Organization registry;
- Principal and current Principal Product-access eligibility;
- OIDC bindings;
- Organization Membership;
- product-defined ordinary `AccessRole` / `Permission` definitions;
- RoleAssignment;
- the bounded access-eligibility continuity fact needed by consumers such as historical Notification routing.

It does **not** own business-action disposition, consequential Governance authority, marketplace policy or source business truth.

### 9.1 Ordinary access kernel

```text
Principal access eligibility
+ Organization Membership
+ RoleAssignment
→ product-defined AccessRole
→ product-defined Permission
```

Disabling Principal access blocks future Product access, including `/access-context`, without erasing Membership/RoleAssignment/history. Revocation changes future authority, not past actor attribution.

Possessing an ordinary Permission does not prove business validity, automation eligibility or consequential authorization.

No custom-role designer, generic ACL/ReBAC engine, OpenFGA/SpiceDB-like platform or generic IAM framework is required for Product 1.0.

---

## 10. Organization isolation, persistence and reference laws

Every organization-owned canonical state and persisted external observation/evidence belongs to one explicit Organization scope. Scope is never inferred from Installation, Selling Entity, source key, IdP organization or process-global default.

A target where isolation relies only on developers remembering query predicates is invalid; exact constraints/RLS/runtime enforcement belong D7.

Persistent semantic classes remain:

1. canonical durable MPC state;
2. external observation/evidence retained only as correctness/history requires;
3. rebuildable derived projection/read model;
4. technical ephemeral state with no business authority.

Known relationships use typed references. There is no universal entity/reference graph, shared mutable Product/Order/Listing model or generic status owner.

Historical snapshots are allowed only for their named purpose and never silently refresh as current producer truth.

Pre-rebaseline legacy DB state is disposable by operator ruling; target starts clean and reacquires external facts/recreates required configuration from authority. Once the target runs, target-owned durable history follows these current lineage laws.

---

## 11. Value, knowledge, provenance and time

### Money/rates/quantities

Canonical monetary/rate/material quantity state uses exact representations where approximation can change correctness. `Money = exact decimal amount + currency`. No universal rounding or UoM framework is implied.

### Fact / absence

Use a bounded `Fact<T>`-style knowledge shape only where knowledge/provenance/temporal validity is materially part of correctness:

```text
Unknown != Known(0) != Empty != Absent != NotApplicable
```

No universal EvidenceID/ObservationID/evidence graph exists. Provenance remains sufficient for the claim and source authority remains explicit.

ADR017/034 remain current evidence until their still-valid domain-judgment clauses are rehomed under their separate retirement condition.

### Time

Keep materially distinct time meanings distinct, including source/effective occurrence time, source commit time, observation/acquisition time, MPC record/materialization time, decision time and external deadline/relative-target anchor. Internal instants are unambiguous instants; scheduler/clock mechanics remain D7.

---

## 12. Correspondence and automation safety

- identity-bearing outbound provider fields must be consistent with accepted Readiness correspondence **before** external dispatch;
- unattended Product↔channel correspondence requires sufficient independent corroboration and no material contradiction; one matching identifier is insufficient by itself;
- recurring automation never silently overrides a standing human decision in the same semantic scope;
- automation/system Principals use explicit authorized paths and never masquerade as humans.

Concrete provider mapping/wire/runtime enforcement belongs D4/D5/D7.

---

## 13. Legacy/ADR transition

Pre-rebaseline ADRs remain history/evidence, not inherited target authority. Still-needed meaning must be present in accepted D-stage/Architecture authority before a legacy ADR retires.

Current gates:

- ADR017/034 remain until the separate Fact-clause rehome condition closes;
- ADR035 remains while the D0–D9 transition authority is needed;
- D7-conditioned legacy ADR residue retires after accepted D7 has demonstrably rehomed its surviving meaning;
- new target ADR numbers are never reused.

Git history is the archive; no legacy database or `docs/archive/` tree is target authority.

---

## 14. Proof obligations

A conforming later realization must be able to prove/falsify at least:

- equal native keys in different sources/Organizations do not collide or cross isolation;
- IdP/email changes do not rewrite Principal identity/history;
- machine actions remain distinct from human attribution;
- disabling Principal eligibility blocks future Product access while preserving history;
- Permission != business disposition/authorization;
- projections/snapshots cannot mutate producer truth;
- unknown never becomes confident zero/default;
- exact financial/material values do not lose correctness through floating-point convenience;
- correspondence mismatch fails before provider effect;
- one identifier cannot establish unattended correspondence alone;
- recurring automation cannot reverse standing human intent silently;
- same Notification occurrence replay does not duplicate personal awareness;
- route edits/unconfigure are temporally correct for delayed/replayed source occurrences;
- revoke→re-enable does not silently reactivate old Notification routing;
- Notifications never grant current source access and do not mutate source truth;
- exactly one terminal AuthorizationDecision can belong to one AuthorizationRequest;
- irrelevant drift does not force needless reapproval, while material authorization-context drift cannot execute from stale permission;
- a later reauthorization episode is historically distinct from its predecessor.

---

## 15. Explicit defers and reopen triggers

D2 does not choose physical ID encoding, tables/indexes/FKs, OIDC vendor/deployment, event/outbox contracts, provider DTO/auth/paging, exact HTTP/Permission operation encoding, frontend layout, queues/workers/RLS/transactions/deployment or generic frameworks rejected above.

Reopen only the smallest implicated D2 meaning when material evidence proves, for example:

- the current owner cannot honestly hold a new independent identity/lifecycle;
- SourceInstance cannot qualify a real second source namespace;
- identity/access begins making business decisions;
- Organization is no longer the correct isolation root;
- domain-local typed references cannot express a real relation without distortion;
- current Notification routing/identity cannot satisfy a real accepted audience/recovery job;
- AuthorizationRequest lifecycle cannot represent a materially new authorization episode class;
- D7 cannot realize fail-closed isolation/temporal routing/access continuity within these semantics.

Preference, hypothetical SaaS/provider scope or framework fashion alone does not reopen D2.
