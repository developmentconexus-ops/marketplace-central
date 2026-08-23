# NOTIF-01 D2-R — Notification Identity & Data Ownership

> **Status:** OPEN — D2-R CANDIDATE / OPERATOR ADJUDICATION REQUIRED
> **Parent accepted authority:** [D2 — Identity / Tenant / Data Ownership](D2-IDENTITY-TENANT-DATA-OWNERSHIP.md) + [D1-R Semantic Closure](D6-R2-NOTIF-01-D1-R-SEMANTIC-CLOSURE.md)
> **Semantic contract set:** [Notification Family Semantic Contracts](D6-R2-NOTIF-01-NOTIFICATION-FAMILY-SEMANTIC-CONTRACTS.md) — OPERATOR-APPROVED
> **Supersedes:** the Work-only D2 candidate embedded in [NOTIF-01 Authority Amendment](D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md)
> **Current Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Why D2 reopens

The original NOTIF-01 D2 candidate correctly established opaque Notification identity, Organization scope, exact-human recipient history and read/archive state, but its `source_work_ref` model assumed Operational Work was the only producer.

Accepted D0-R/D1-R now establish:

```text
10 explicit source-owner edges
→ 14 semantic NotificationKinds
→ DIRECT_SOURCE / OWNER_DERIVED / ORG_ROUTED audience strategies
→ bounded Organization routing configuration
→ two bounded per-recipient awareness-suppression laws
```

D2-R therefore rederives identity, references, current durable state and correlation semantics only. It does not choose Q/C/E/P, event names, queue/job mechanics, PostgreSQL schema/RLS, HTTP paths, Permissions or frontend realization.

---

## 2. Governing D2-R laws — CANDIDATE

1. **One Notification is one durable personal-awareness occurrence.** It is not the source business occurrence and not an audit/event record.
2. **Organization remains the isolation root.** Every Notification and routing configuration belongs to exactly one Organization.
3. **Recipient is one exact canonical human Principal.** Email, username, role name and Permission are never recipient identity.
4. **Source references are closed and typed.** Multiple accepted producers justify one bounded closed union; no universal entity graph is admitted.
5. **Source occurrence identity remains producer-owned.** Personal Notifications receives a stable owner-local occurrence discriminator; it does not create a universal EventID.
6. **Historical awareness is stable.** Later routing, Membership, source-state or source-access changes do not retarget an existing Notification.
7. **Routing responsibility is not access authority.** Configuration chooses who should receive awareness; ordinary access still decides what that Principal may read/do.
8. **No source-state duplication by convenience.** Notification stores correlation needed for awareness, deduplication and navigation, not a second current copy of Sale/Fulfillment/Work/etc.
9. **Suppression correlation is bounded.** Only the two operator-approved overlap cases gain explicit typed correlation; no generic related-entity/cross-kind graph.
10. **Absence remains explicit.** Unconfigured routing is not silently replaced by all admins/members or by Permission-derived recipients.

---

## 3. Canonical Notification identity — CANDIDATE

`Notification` is canonical durable MPC state owned by **Personal Notifications** and has one stable opaque `NotificationID`.

Binding laws:

- opaque, non-business-semantic and non-reusable;
- does not encode Organization, Principal, `NotificationKind`, source identity, occurrence time or lifecycle state;
- read/unread and archive/unarchive never change Notification identity;
- later source changes never change Notification identity;
- routing changes never move an existing Notification to another recipient;
- replay of the same semantic awareness occurrence does not mint another NotificationID.

Exact UUID/ULID/database encoding remains D7 realization.

---

## 4. `NotificationKind` identity semantics — CANDIDATE

The fourteen operator-approved values are stable Product semantic keys owned by Personal Notifications:

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

A kind is not silently repurposed. If a future family splits, new semantic values are added through the smallest D0/D1 reopen while historical Notifications retain the meaning of the value under which they were created.

Fourteen is not a protected count.

---

## 5. Canonical Notification state — CANDIDATE

The smallest durable state required by accepted consumers is:

```text
notification_id
organization_id
recipient_principal_id
kind
source_ref
source_occurrence_key
source_occurred_at
created_at
read_at?
archived_at?
revision
```

Semantics:

```text
read_at = null       → unread
read_at != null      → read

archived_at = null   → active Inbox
archived_at != null  → archived
```

`read_at` and `archived_at` are orthogonal. Archiving does not imply read, and marking unread does not unarchive.

`source_occurred_at` is immutable source-owned occurrence time preserved for truthful chronology. `created_at` is when Personal Notifications durably created the personal awareness item. Propagation delay must not rewrite one into the other.

`revision` is owner-local mutation lineage for stale-write distinguishability. Exact ETag/header/storage representation belongs D5/D7.

Baseline canonical Notification state does **not** add:

```text
seen
delivered
acknowledged
resolved
dismissed
severity
priority
generic status
free-form payload
arbitrary metadata
template variables
stored source-business-state copy
```

If a later D5/D6 consumer proves that an immutable presentation snapshot is required, D2 must reopen for explicit typed fields rather than hiding business data in a generic payload.

---

## 6. Organization and historical recipient — CANDIDATE

Every Notification belongs to exactly one `Organization` and exactly one canonical human `Principal` recipient.

At Notification materialization time the recipient must be a human Principal valid for that Organization under current identity/access eligibility. Exact Permission compatibility for each Product read surface is D5 authority; recipient selection itself is never inferred from Permission.

Historical laws:

- recipient identity never changes because display name/e-mail/OIDC binding changes;
- reassignment of Work does not retarget an earlier `WORK_ASSIGNMENT` Notification;
- routing edits affect future occurrences only;
- Membership/Principal-access revocation blocks future personal Notification materialization and Inbox access but does not rewrite historical Notification ownership;
- re-enabling/re-adding a Principal does not silently backfill missed Notifications or resurrect prior ORG_ROUTED responsibility;
- possession of a Notification never grants source access.

---

## 7. Closed typed `NotificationSourceRef` — CANDIDATE

Multiple accepted producers now justify a **closed typed union**. This is explicitly not a platform-wide `{entity_type, entity_id}` graph.

The allowed source-reference variants are the smallest references already meaningful under accepted owner identity:

```text
MarketplaceInstallationRef
ListingIntentRef
PriceIntentRef
AvailabilityTargetRef
EconomicAttributionRef
MarketplaceSaleRef
BusinessOrderIntentRef
InvoicingIntentRef
FulfillmentExecutionRef
MarketplaceShipmentRef
PostSaleResolutionRef
WorkRef
AuthorizationTargetRef
AuthorizationDecisionRef
```

Reference semantics:

- `MarketplaceInstallationRef` → canonical Marketplace Installation ID.
- `ListingIntentRef` / `PriceIntentRef` → owner-local MPC Intent IDs.
- `AvailabilityTargetRef` → the accepted Availability-owned closed target identity: pre-creation ListingIntent target **or** existing source-qualified Marketplace Listing target. No public AvailabilityIntent ID is invented merely for Notification symmetry.
- `EconomicAttributionRef` → canonical Economic Attribution ID.
- `MarketplaceSaleRef` → Marketplace Installation + provider-native Sale key; no synthetic Sale ID.
- `BusinessOrderIntentRef` / `InvoicingIntentRef` → owner-local MPC Intent IDs.
- `FulfillmentExecutionRef` → canonical Fulfillment Execution ID already exposed by accepted Product authority.
- `MarketplaceShipmentRef` → Marketplace Installation + provider-native Shipment key; no synthetic Shipment ID.
- `PostSaleResolutionRef` → canonical Post-Sale Resolution ID.
- `WorkRef` → canonical Work ID.
- `AuthorizationTargetRef` → the accepted closed governed-target union: ListingIntent, PriceIntent, BusinessOrderIntent or InvoicingIntent identity. It identifies a pending governed target without pretending a completed Authorization Decision already exists.
- `AuthorizationDecisionRef` → canonical Authorization Decision ID after a decision occurrence exists.

No arbitrary source-ref variant may be added by UI/provider convenience.

### 7.1 Kind ↔ source compatibility matrix — CANDIDATE

| NotificationKind | Allowed `source_ref` |
| --- | --- |
| `MARKETPLACE_INSTALLATION_ATTENTION` | `MarketplaceInstallationRef` |
| `OFFERING_ASYNC_ACTION_RESULT` | `ListingIntentRef` or `PriceIntentRef` |
| `AVAILABILITY_ATTENTION` | `AvailabilityTargetRef` |
| `ECONOMIC_RECONCILIATION_ATTENTION` | `EconomicAttributionRef` or `MarketplaceSaleRef` |
| `NEW_MARKETPLACE_SALE` | `MarketplaceSaleRef` |
| `SALE_ATTENTION` | `MarketplaceSaleRef` |
| `MATERIALIZATION_ATTENTION` | `BusinessOrderIntentRef` or `InvoicingIntentRef` |
| `FULFILLMENT_ACTIONABLE` | `FulfillmentExecutionRef` |
| `FULFILLMENT_ATTENTION` | `FulfillmentExecutionRef` |
| `SHIPMENT_EXCEPTION` | `MarketplaceShipmentRef` |
| `POST_SALE_ATTENTION` | `PostSaleResolutionRef` |
| `WORK_ASSIGNMENT` | `WorkRef` |
| `AUTHORIZATION_ACTION_REQUIRED` | `AuthorizationTargetRef` |
| `AUTHORIZATION_DECISION_RESULT` | `AuthorizationDecisionRef` |

A mismatched kind/source pair is invalid state.

### 7.2 Economic reconciliation does not gain a universal ID — CANDIDATE

Current accepted consumers can identify economic attention through either a durable `EconomicAttributionID` or a source-qualified Sale whose `SaleEconomics` reconciliation requires attention.

D2-R therefore does **not** introduce a universal `EconomicReconciliationID`. If a future R3/bank-cash reconciliation workflow requires an independently addressable persistent case beyond these accepted subjects, D2 reopens then.

---

## 8. Source occurrence identity and semantic deduplication — CANDIDATE

Every admitted source owner must expose one stable opaque **owner-local `source_occurrence_key`** for each committed attention occurrence.

It is scoped by the source owner/source reference and is not a universal EventID, event-store identity or cross-domain sequence.

The semantic uniqueness of one personal Notification is:

```text
Organization
+ recipient Principal
+ NotificationKind
+ typed source_ref
+ source_occurrence_key
```

Therefore:

```text
same source occurrence replayed/recovered twice for P
→ one semantic Notification

distinct attention occurrence on same source for P
→ distinct Notification

Work A→B→A
→ two distinct WORK_ASSIGNMENT occurrences for A

same unresolved Shipment exception re-polled
→ same occurrence, no new Notification

resolved Shipment later enters a genuinely new exception episode
→ new occurrence, new Notification
```

The source owner owns what constitutes a distinct occurrence. Personal Notifications owns deduplication of personal awareness from that occurrence.

No global monotonic event sequence is required.

---

## 9. Bounded Organization routing state — CANDIDATE

`ORG_ROUTED` configuration is canonical durable state owned by Personal Notifications.

It does **not** require a synthetic RoutingConfigID. Its semantic key is naturally:

```text
NotificationRouteKey = (organization_id, notification_kind)
```

Current state is:

```text
organization_id
notification_kind
recipient_principal_ids  // explicit set, one or more
revision
```

Only these ten accepted ORG_ROUTED kinds may have baseline routing configuration:

```text
MARKETPLACE_INSTALLATION_ATTENTION
AVAILABILITY_ATTENTION
ECONOMIC_RECONCILIATION_ATTENTION
NEW_MARKETPLACE_SALE
SALE_ATTENTION
MATERIALIZATION_ATTENTION
FULFILLMENT_ACTIONABLE
FULFILLMENT_ATTENTION
SHIPMENT_EXCEPTION
POST_SALE_ATTENTION
```

Routing state for DIRECT_SOURCE or OWNER_DERIVED kinds is invalid because configuration must not override source-owned recipient meaning.

### 9.1 Unconfigured vs empty — CANDIDATE

Launch V1 has no proved consumer for “configured intentionally to nobody”. Therefore the baseline state is:

```text
no route record / no configured recipient set
→ UNCONFIGURED

route present
→ one-or-more explicit human Principal recipients
```

A zero-recipient configured state is not admitted by symmetry. If the Product later needs an explicit “disable this kind for everyone” control distinct from unconfigured, D2/D5 reopen.

Unconfigured never means all admins, all members, all Permission holders or a default role.

### 9.2 Recipient-reference validity — CANDIDATE

A configured route recipient must be an exact human Principal belonging to the same Organization at configuration time.

Future resolution also fails closed for a Principal whose current Product-access eligibility/Membership is no longer valid. Revocation does not silently preserve latent routing responsibility such that later re-eligibility resumes Notifications without an explicit routing decision.

Exact cross-owner cleanup/reaction mechanics are D3/D7; D2 fixes the semantic outcome only.

Permission may be an additional eligibility condition later, but it never selects the recipient set.

---

## 10. Audience resolution and historical stability — CANDIDATE

### 10.1 DIRECT_SOURCE

The source occurrence supplies the exact human Principal identity already owned by that producer. Personal Notifications records that Principal as the historical recipient.

Baseline families:

```text
OFFERING_ASYNC_ACTION_RESULT
WORK_ASSIGNMENT
AUTHORIZATION_DECISION_RESULT
```

The applicable source-owned historical lineage must be sufficient to resolve the exact human; it may not fall back to current role holders or e-mail matching.

For `AUTHORIZATION_DECISION_RESULT`, Governance must preserve sufficient requester/initiator lineage in the authorization context to identify the exact human when that family applies. This is historical correlation, not transfer of the action-owning domain's business-intent authority.

### 10.2 OWNER_DERIVED

For `AUTHORIZATION_ACTION_REQUIRED`, Controlled Action Governance resolves the exact currently valid decision Principal set from Governance authority/delegation semantics.

The Notification source is the pending `AuthorizationTargetRef`, not an `AuthorizationDecisionID`: under current accepted D2/D5 meaning the canonical Authorization Decision exists only when the consequential decision occurrence is created.

Each distinct Principal eligibility episode must have a distinct source occurrence key if it is to create new awareness.

### 10.3 ORG_ROUTED

Personal Notifications resolves recipients from the current `(Organization, NotificationKind)` routing state after applying current identity/access eligibility.

The resolved recipient set is copied only into individual Notification recipient identities; existing Notifications are never dynamically rebound to the current route.

Routing changes do not backfill prior source occurrences by default.

---

## 11. Two bounded suppression correlations — CANDIDATE

D1-R admits exactly two duplicate-awareness precedence rules. D2-R adds only the typed correlations required to make them falsifiable; it does not create a generic relation graph.

### 11.1 Source alert → assigned Work replacement

A source-owned attention occurrence may carry an explicit causal `WorkRef` when that exact occurrence produced Operational Work.

For recipient `P`, source-family awareness may be suppressed only when the correlated Work has a distinct assignment occurrence to `P` that qualifies for `WORK_ASSIGNMENT`.

Correlation is by exact source occurrence + exact Work identity/assignment occurrence, never title/timestamp/entity-name matching.

### 11.2 Sale attention → richer Post-Sale continuation

A `SALE_ATTENTION` occurrence may carry an explicit correlated `PostSaleResolutionRef` when the exact Sale-owned consequence is the precursor of that resolution.

For recipient `P`, `SALE_ATTENTION` may be suppressed only when the same exact `PostSaleResolutionRef` also produces `POST_SALE_ATTENTION` for `P` and Post-Sale is the richer continuation.

No generic `related_entities[]`, `correlation_type:string`, causal graph or universal consequence ID is introduced.

### 11.3 Ordering remains D3 proof — CANDIDATE

D2-R does not add `superseded`, `replaced_by` or another generic Notification lifecycle state merely to guess future delivery ordering.

D3 must prove that duplicate/reordered/replayed propagation can honor the two bounded suppression outcomes without rewriting source truth or fabricating a generic dedup engine. If that cannot be proved without durable late-supersession state, D2-R must reopen for that **specific** requirement before D3 closes.

---

## 12. Source access and data-minimization law — CANDIDATE

A Notification is not a capability token.

```text
Notification
→ typed source_ref
→ current Product read
→ current source access/Permission checked again
```

Binding laws:

- source access loss does not rewrite historical Notification recipient identity;
- source click fails closed under current access;
- Notification canonical state does not copy mutable current source business state merely to avoid source reads;
- no provider DTO, buyer PII, fiscal payload, address, payment detail or arbitrary source body is stored in Notification by default;
- no free-form template/payload field becomes a backdoor source mirror.

The Inbox may later expose safe typed presentation projections through D5, but those projections do not become second current authority.

---

## 13. Write authority and lifecycle — CANDIDATE

Personal Notifications is the sole write authority for:

```text
Notification identity
historical recipient
read/unread
active/archived
ORG_ROUTED configuration
personal-awareness deduplication/suppression result
```

Source owners remain sole authority for source facts and source occurrence qualification.

Identity/access remains authority for Principal, Membership, access eligibility, AccessRole and Permission.

Read/archive mutations change Personal Notifications state only. They do not acknowledge, resolve, authorize, dismiss or mutate the source object.

Deletion/retention policy remains deferred; source deletion/revocation/reassignment must not silently rewrite historical Notification meaning.

---

## 14. D2-R negative controls — CANDIDATE

D2-R fails review if it:

1. retains Work as the only possible source;
2. creates a universal `{entity_type, entity_id}` source graph;
3. introduces a universal EventID/event-store identity for source occurrences;
4. invents synthetic Sale/Shipment IDs over already source-qualified identities;
5. invents a public AvailabilityIntent ID solely for Notification symmetry;
6. invents a universal EconomicReconciliationID without a proved independent case consumer;
7. uses `AuthorizationDecisionID` for `AUTHORIZATION_ACTION_REQUIRED` before a decision exists;
8. allows kind/source-ref combinations outside the closed compatibility matrix;
9. derives recipients from role/Permission/all-admin membership;
10. lets DIRECT_SOURCE/OWNER_DERIVED kinds be overridden by ORG_ROUTED configuration;
11. preserves ORG_ROUTED responsibility through membership/access revocation such that later reactivation silently resumes it;
12. collapses distinct source occurrences because they share subject+recipient;
13. uses free-form payload/metadata/template state to smuggle source truth into Notification;
14. adds generic `seen`, severity, priority, delivered/acknowledged/resolved state;
15. creates generic related-entity/cross-kind dedup graphs for the two bounded suppression cases;
16. chooses table/RLS/event/River/API/frontend mechanics inside D2-R.

---

## 15. D2-R coherence result — CANDIDATE

The rederived model preserves accepted D2 laws:

```text
one opaque NotificationID
+ one Organization
+ one exact historical human Principal recipient
+ one stable NotificationKind
+ one closed typed source_ref
+ one source-owner-local occurrence key
+ truthful source occurrence time
+ read/archive owner-local lifecycle
+ bounded semantic routing key (Organization, NotificationKind)
+ two explicit suppression correlations only
```

It also resolves three identity questions exposed by the semantic pass without adding generic machinery:

1. `AUTHORIZATION_ACTION_REQUIRED` references the pending governed target; `AUTHORIZATION_DECISION_RESULT` references the later canonical decision.
2. Availability attention uses the already-accepted closed `AvailabilityTarget` identity rather than inventing a public intent resource.
3. Economic attention uses existing Economic Attribution or Sale identity; no universal Reconciliation entity is created.

No D3 communication form, D5 wire/Permission, D6 UI, D7 physical persistence/runtime or D8 proof becomes accepted through this candidate.

---

## 16. Gate

```text
D0-R Product scope                    ACCEPTED
D1-R boundary authority               ACCEPTED / OPERATOR-RATIFIED
14 family semantic contracts          OPERATOR-APPROVED
D2-R identity/data ownership          READY FOR OPERATOR REVIEW
D3 communication/events               BLOCKED
D5 Product/OAD                        BLOCKED
D6 bell/Inbox/settings                BLOCKED
D7 runtime                            BLOCKED
D8 proof                              BLOCKED
Product implementation                BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates only this D2-R identity/data-ownership candidate. If approved, close D2-R and open only the bounded D3 communication/occurrence-propagation gate. Do not edit the OAD or begin D6/D7 before D3/D5 authority is ratified.
