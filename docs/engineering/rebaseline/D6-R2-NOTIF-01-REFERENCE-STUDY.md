# NOTIF-01 — Global Notification Reference Study

> **Status:** REFERENCE STUDY + TRIGGER/AUDIENCE CENSUS CANDIDATE / NO PRODUCT AUTHORITY
> **Trigger:** operator falsified the Work-only notification-origin assumption during D2 review
> **Method:** Frontend Product Experience Planning Method v2.1 + Global Maximum + YAGNI
> **Current authority impact:** D0/D1 accepted bounded baseline remains historical authority; D0/D1 trigger scope is reopened for bounded correction; D2 candidate is suspended pending re-derivation
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Question

What is the smallest notification architecture that gives Marketplace Central users the awareness they actually need across the whole Product 1.0 lifecycle without turning Notifications into a generic event bus, duplicating domain truth, or forcing routine operational work into `Operational Work` merely so it can notify someone?

The operator supplied the material falsifier:

> Notifications must cover real consumer needs such as a new sale, something that must be shipped, or another materially relevant event; Work assignment cannot be the only origin merely because it is the easiest exact-recipient case.

## 2. External reference evidence

This study extracts patterns only. No external product is authority for MPC semantics.

### Shopify Admin

Sources:

- Alerts feed / topbar notification icon: https://help.shopify.com/en/manual/shopify-admin/shopify-admin-overview?links=false
- Staff notifications: https://help.shopify.com/en/manual/fulfillment/setup/notifications/staff-notifications
- Order management: https://help.shopify.com/en/manual/fulfillment/managing-orders

Observed patterns:

1. Shopify has a topbar Alerts feed for important/time-sensitive information and required actions.
2. Alerts are delivered to the owner and staff with relevant access for those scenarios.
3. Separate staff notifications explicitly cover business events such as a **new order** and **new return request**.
4. The alert/feed is an awareness surface; the underlying order/return/admin page remains the operational source.
5. Shopify exposes read/unread feed behavior and an unread indicator, while the source object remains separate.
6. Notification recipients can be configured per staff notification type; recipient selection and ordinary Product access are distinct concerns.

### Amazon Seller Central

Sources:

- Notification preferences overview: https://sellercentral.amazon.com/seller-forums/discussions/t/6fe572a2-9524-4731-9b84-58ae7bb57132
- Sales/fulfillment notification example: https://sellercentral.amazon.com/seller-forums/discussions/t/f760499d-2c92-4444-baf8-93fe285e67a7

Observed patterns:

1. Notification categories include **Orders**, fulfillment/inbound shipment-related alerts, **Returns and Claims**, account/business alerts and emergency alerts.
2. Different notification types can be directed to different people/addresses; Amazon explicitly gives the example of order notifications going to fulfillment while messaging goes to customer service.
3. The platform therefore separates **event category** from **recipient/audience choice**.
4. Users can suppress categories that become noise; the existence of preferences is evidence that broad event coverage without audience/noise control degrades usefulness.

### Mercado Livre

Sources:

- Notification topics: https://developers.mercadolivre.com.br/produto-receba-notificacoes
- Claims notifications: https://developers.mercadolivre.com.br/pt_br/atributos/gerenciar-reclamacoes
- Sales management / fraud-stop evidence: https://developers.mercadolivre.com.br/pt_br/gerenciamento-de-vendas

Observed patterns:

1. Mercado Livre exposes distinct topics such as `orders_v2`, shipments, messages and post-purchase claims/actions.
2. `orders_v2` covers creation/changes to confirmed sales; provider sale evidence can also carry high-impact changes such as a fraud-risk stop that must prevent shipment.
3. Claims can notify immediately on claim creation/action and expose due dates/responsibility that materially affect seller action.
4. The notification payload points to a resource that the consumer must `GET` for authoritative detail.
5. This strongly supports MPC's existing law that a notification/signal must not become the source-domain truth.
6. Provider notification topics are evidence acquisition, not MPC Personal Notification kinds by automatic 1:1 mapping.

### ShipStation

Source:

- Order Alerts: https://help.shipstation.com/hc/en-us/articles/360025869112-Order-Alerts

Observed patterns:

1. ShipStation exposes specific operational alert classes, such as combine-order opportunities, low inventory thresholds and automation-rule alerts.
2. Alerts are visible through a badge near the profile/Orders UI.
3. It does not model every order/shipment state transition as an alert; alert taxonomy is curated for operator relevance.

### Jira

Sources:

- Notification schemes: https://support.atlassian.com/jira-cloud-administration/docs/configure-notification-schemes/
- Project/space notification configuration: https://support.atlassian.com/jira-cloud-administration/docs/configure-a-project/

Observed patterns:

1. Jira binds **specific events** to **specific recipient rules** instead of treating all changes as one generic notification stream.
2. Recipients may come from source-owned semantics such as assignee/reporter/current user, or from configured users/roles/groups.
3. Notification recipients must still be allowed to view the underlying work item; being selected for a notification does not replace source authorization.
4. This is strong evidence for separating responsibility/audience selection from ordinary read permission while still validating source access eligibility.

### GitHub

Sources:

- Notifications overview/configuration: https://docs.github.com/en/subscriptions-and-notifications
- Inbox filters/reasons: https://docs.github.com/en/subscriptions-and-notifications/reference/inbox-filters

Observed patterns:

1. GitHub records a **reason** for why the user received an Inbox item: assignment, mention, review request, security alert, state change, CI activity, etc.
2. Users can triage and filter by reason/event family rather than receiving one undifferentiated activity feed.
3. The header Inbox uses an unread indicator, reinforcing that awareness presence and exact global count are separate UX choices.
4. Subscription customization exists because notification relevance is contextual; this supports preserving a future preference seam without making a generic subscription engine a Launch-V1 requirement.

## 3. Cross-platform convergence

The useful common pattern is not “notify every event” and not “notify only assigned tasks”. It is:

```text
material source occurrence
  → product-defined notification kind
  → explicit audience/recipient resolution
  → durable personal awareness item
  → deep link / reread of current source truth
```

Converged laws:

1. **Multiple producer families are legitimate.** Sales, offering/effect outcomes, availability, materialization, fulfillment/shipping, post-sale/claims, assigned work, governance and material channel-health outcomes can all have real consumers.
2. **Curated trigger taxonomy beats event-per-CRUD.** Normal low-value transitions must not flood the Inbox.
3. **Audience is part of the product design.** Permission alone is not responsibility; new-sale and shipment notifications need an explicit recipient rule.
4. **Notification is not source truth.** Opening the item must re-read/re-authorize the source.
5. **Attention surfaces and operational workspaces coexist.** The notification feed tells the user what happened; Sales/Fulfillment/Post-Sale/Work/etc. remain the places where domain work is understood/executed.
6. **Noise control is a first-class concern.** External products expose per-category routing/preferences because one undifferentiated feed degrades at scale.
7. **Source-owned responsibility should be reused when it exists.** Assignment/authorization semantics should not be duplicated into a second Notification responsibility model.

## 4. Hypotheses

### H1 — Work-only personal Inbox

```text
Operational Work → Personal Notifications
```

**FALSIFIED.** It handles exact assignment elegantly but cannot represent legitimate awareness such as a new confirmed sale or a material shipment exception unless MPC fabricates Work for routine operational events. That would distort `Operational Work` authority.

### H2 — Generic AnyDomain event fan-out

```text
AnyDomain event → generic Notifications router → users
```

**REJECT.** It encourages event-per-CRUD, generic entity references, hidden cross-domain dependencies, alert fatigue and a new hub that starts to acquire semantic authority.

### H3 — Curated Product notification kinds + explicit producer edges + bounded audience routing

**OPERATOR-APPROVED DIRECTION.** Each notification kind is Product-defined, owned by Personal Notifications as awareness semantics, backed by an explicit source-owner occurrence and admitted only when a user need is proven.

The trigger/audience census below now derives the bounded Launch-V1 candidate from the whole accepted Product lifecycle.

---

# Trigger + Audience Census — CANDIDATE

## 5. Admission law

A candidate occurrence earns a Launch-V1 Personal Notification only when **all** of these hold:

1. the occurrence is committed meaning owned by an accepted source owner;
2. a proved human job benefits from noticing it while the user may be outside the source workspace;
3. the awareness has independent value from merely rereading a queue later;
4. exact human recipients can be resolved without treating ordinary Permission as responsibility;
5. a truthful deep link/source reread exists;
6. repeat delivery can be deduplicated by source occurrence identity;
7. the trigger has an explicit noise/suppression rule;
8. Notification can remain awareness only and does not acquire source mutation/workflow authority.

A candidate is rejected/deferred when it is merely:

- every state/CRUD transition;
- routine progression already visible in the specialized queue;
- an analytical metric change with no attention boundary;
- a provider webhook/topic copied 1:1 into Product semantics;
- a condition with no legitimate exact human audience;
- an infrastructure/on-call incident rather than Product-user awareness.

## 6. Audience strategies

The Global-Maximum model needs **three** bounded strategies; fewer would duplicate existing responsibility or leave real consumers unsupported.

### A1 — DIRECT_SOURCE

The committed source occurrence already identifies the exact human Principal.

Examples:

```text
Work assigned/reassigned → assignee Principal
user-initiated async effect result → initiating Principal, when lineage is exact
```

No Notification routing configuration selects another recipient.

### A2 — OWNER_DERIVED

The source owner already owns the responsibility/authorization semantics needed to derive exact current human Principals.

Candidate example:

```text
Governance authorization becomes actionable
→ Governance determines exact currently eligible decision Principals
→ Notifications receives the resolved Principal set
```

Notifications must not recreate Governance delegation/authority merely to route an alert.

### A3 — ORG_ROUTED

The source occurrence has real organization-wide operational value but no source-owned exact human recipient.

Examples:

```text
new confirmed sale
shipment exception
fulfillment becomes actionable
```

Personal Notifications owns a small Organization routing configuration:

```text
(Organization, Product-defined NotificationKind)
→ configured exact human Principals
```

Candidate laws:

- exact human Principals only in Launch V1;
- no arbitrary expression language, nested groups or routing DSL;
- ordinary Permission is an **eligibility/access check**, never the routing selector by itself;
- source-read/access eligibility must be revalidated when a Notification is created/opened; routing never grants source access;
- a missing route is explicit configuration absence, never implicit broadcast to all members or admins;
- route changes affect future occurrences and never rewrite historical recipients.

## 7. Launch-V1 candidate kind census

The table covers the accepted N01–N16 jobs and D8 golden flows. `ADMIT` here is a **census candidate**; D0/D1 authority is not corrected until the operator ratifies this census.

| Kind family | Exact source occurrence / admission boundary | Source owner | Main human consumer | Audience | Deep-link home | Disposition |
| --- | --- | --- | --- | --- | --- | --- |
| `MARKETPLACE_INSTALLATION_ATTENTION` | Installation becomes materially non-operable / external authorization or effective capability blocks normal operation; or provider-authoritative health/reputation crosses an owner-declared material attention boundary | Marketplace Portfolio | marketplace operator / admin | A3 ORG_ROUTED | `/configuracoes/canais/:marketplaceInstallationId` | **ADMIT** |
| `OFFERING_ASYNC_ACTION_RESULT` | a human-initiated ListingIntent/PriceIntent consequential effect was not terminal in the initiating request and later becomes converged, rejected, ambiguous or divergent | Marketplace Offering Operations | initiating operator/manager | A1 DIRECT_SOURCE via exact initiator lineage | Listing/Intent/Price workspace | **ADMIT** |
| `AVAILABILITY_ATTENTION` | automatic or authorized Availability synchronization becomes materially blocked, ambiguous, divergent or unable to maintain accepted sellable availability | Availability Control | marketplace operator | A3 ORG_ROUTED | `/disponibilidade` | **ADMIT** |
| `ECONOMIC_RECONCILIATION_ATTENTION` | an attribution/reconciliation enters an owner-declared state requiring bounded human resolution; ordinary variance/metric movement alone does not qualify | Commercial Economics | commercial/marketplace manager | A3 ORG_ROUTED | `/economia/reconciliacao` | **ADMIT** |
| `NEW_MARKETPLACE_SALE` | first authoritative confirmation of a new Marketplace Sale in MPC; provider change/update noise does not retrigger it | Marketplace Sales | marketplace operations; optionally fulfillment via configured routing | A3 ORG_ROUTED | `/vendas/:nativeSaleKey` | **ADMIT** |
| `SALE_ATTENTION` | a committed Sale-owned condition materially changes safe handling, e.g. cancellation/hold/fraud-risk-like stop or other owner-declared exception; not every order update | Marketplace Sales | marketplace operator | A3 ORG_ROUTED | `/vendas/:nativeSaleKey` | **ADMIT** |
| `MATERIALIZATION_ATTENTION` | Business Order or Invoicing materialization becomes ambiguous/rejected/divergent/external-required in a way that needs human attention; normal successful progression does not qualify | Business-System Materialization | marketplace operations / fulfillment depending on source condition | A3 ORG_ROUTED, subject to Work suppression | `/vendas/:nativeSaleKey` materialization region | **ADMIT** |
| `FULFILLMENT_ACTIONABLE` | FulfillmentExecution crosses once into the first owner-declared state where physical execution can legitimately begin; not every subsequent checkpoint | Fulfillment Lifecycle | fulfillment/dispatch operator | A3 ORG_ROUTED | `/expedicao/execucoes/:fulfillmentExecutionId` | **ADMIT** |
| `FULFILLMENT_ATTENTION` | owner-declared material attention transition such as dispatch becoming due/at-risk, physical conference discrepancy, blocked physical readiness or provider-requirement closure failure | Fulfillment Lifecycle | fulfillment/dispatch operator | A3 ORG_ROUTED, subject to Work suppression | `/expedicao/execucoes/:fulfillmentExecutionId` | **ADMIT** |
| `SHIPMENT_EXCEPTION` | Shipment enters a material exception/outcome requiring operator awareness; repeated observation of the same exception is not a new occurrence | Fulfillment Lifecycle shipment observation | marketplace/fulfillment operations | A3 ORG_ROUTED, subject to Work suppression | `/expedicao/envios/:nativeShipmentKey` | **ADMIT** |
| `POST_SALE_ATTENTION` | PostSaleResolution opens or provider evidence creates a new seller-relevant cancellation/return/refund/claim consequence, action requirement or due-date attention transition | Post-Sale Resolution | marketplace/post-sale operations | A3 ORG_ROUTED, subject to Work suppression | `/pos-venda/:postSaleResolutionId` | **ADMIT** |
| `WORK_ASSIGNMENT` | Work assignment/reassignment commits to an exact human Principal | Operational Work | exact assignee | A1 DIRECT_SOURCE | `/trabalho/:workId` | **ADMIT — mandatory personal** |
| `AUTHORIZATION_ACTION_REQUIRED` | an accepted governed action becomes personally actionable to exact currently valid decision Principal(s) according to Governance authority | Controlled Action Governance | approver/policy authority | A2 OWNER_DERIVED | `/aprovacoes/:authorizationDecisionId` or exact governed target | **ADMIT** |
| `AUTHORIZATION_DECISION_RESULT` | a Governance decision is committed and exact requester/initiator lineage exists for a human who needs the result | Controlled Action Governance | requester/initiating operator | A1 DIRECT_SOURCE | `/aprovacoes/:authorizationDecisionId` and governed target | **ADMIT** |

### 7.1 Why these 14 families are not “too simple”

They cover all proved cross-screen awareness classes without forcing one kind per operation or per provider topic:

```text
channel operability / health
consequential action completion
availability safety
commercial reconciliation
new demand
sale safety
business-system materialization
physical execution entry
physical execution attention
shipment divergence
post-sale consequences
personal Work responsibility
governance action
governance result
```

Several source occurrence subtypes can map to one Product kind only when they share the same human job, audience strategy and source/deep-link semantics. If later proof shows two subtypes require materially different audience or interaction, split the kind rather than add a generic `reason` DSL.

## 8. Explicit DEFER / REJECT census

| Candidate occurrence | Disposition | Reason |
| --- | --- | --- |
| every Product/Channel readiness missing/conflict change | **REJECT baseline** | Preparation is a deliberate workspace; per-product feed would be high-volume and duplicates the job surface. Material exceptions may become Work through accepted owner semantics. |
| every Listing/Price/Availability successful sync | **REJECT** | routine convergence noise; only human-initiated async result or material attention transition qualifies. |
| every Market/Performance/Economics metric change | **REJECT** | analytical evidence belongs in strategy surfaces; no generic threshold-alert engine is justified. |
| normal Business Order/Invoicing convergence success | **REJECT** | downstream Fulfillment actionable transition is the meaningful human attention boundary. |
| every separation/conference/packing/dispatch checkpoint | **REJECT** | specialized Fulfillment queue owns routine progression; alerting each step causes duplicate noise. |
| normal dispatch handoff / normal delivered shipment | **DEFER** | useful as optional future completion feedback but not required to keep Launch-V1 operations safe/actionable. |
| Work hold/resume/close without a new exact responsibility occurrence | **REJECT default** | Work workspace/history already carries it; future user subscriptions are not admitted. |
| every access Membership/RoleAssignment change | **DEFER** | in-app delivery may be unavailable precisely when access is removed; access administration/history remains authoritative. |
| buyer Q&A/chat/general marketplace messages | **OUT OF SCOPE** | D0 explicitly excludes buyer-conversation management; provider message topics do not silently reopen Product scope. |
| provider webhook/topic arrival itself | **REJECT as Product kind** | ingress evidence must be translated/reread by the owning domain before Product awareness is created. |
| runtime outage, deployment failure, worker backlog, database incident | **REJECT from Personal Notifications** | these are observability/on-call concerns, not end-user Product awareness semantics. |
| infrastructure retry/recovery progress | **REJECT** | technical mechanism does not become Product meaning. |

## 9. Per-recipient duplicate / suppression law

The Inbox must not punish users for the same underlying attention transition being represented in source truth **and** Work.

Candidate law:

```text
source occurrence S
→ routed source Notification to Principal P
→ same S immediately/causally creates Work W already assigned to P
→ suppress the routed source Notification for P
→ deliver WORK_ASSIGNMENT for P
```

This suppression is **per recipient**, not global. Other configured recipients who are not assigned the Work may still need the source-owner alert.

If Work is assigned later as a genuinely new responsibility occurrence, that later `WORK_ASSIGNMENT` is a new Notification even if the user saw the earlier source alert.

Other noise laws:

- replay/duplicate acquisition of the same source occurrence never creates another semantic Notification;
- polling the same unresolved state never retriggers awareness;
- a materially different later attention transition may notify again (`FULFILLMENT_ACTIONABLE` → later `FULFILLMENT_ATTENTION` is legitimate);
- read/archive does not reset merely because the source later changes; a new source occurrence creates a new Notification;
- changing routing recipients never retroactively creates, moves or deletes historical Notifications by default.

## 10. Routing-configuration ownership candidate

The census confirms a real D1 responsibility that the Work-only design did not have:

> **Personal Notifications must own bounded Organization routing configuration for Product-defined A3 kinds.**

This is not ordinary access control and not business-domain responsibility.

Candidate configuration meaning:

```text
Organization
+ Product-defined NotificationKind
→ exact configured human Principal recipients
→ revision / historical configuration lineage as needed
```

The source owner still decides **that the material occurrence exists**. Personal Notifications decides only **which configured Principals should receive awareness for A3 kinds**.

The configuration must preserve these distinctions:

```text
unconfigured route ≠ configured recipient set
routing recipient ≠ ordinary Permission holder
routing recipient ≠ source owner
routing change ≠ historical Notification rewrite
```

A later D2 correction must decide the minimum persistent representation. D5/D6 later decide configuration operations/surface. No arbitrary subscriptions/preferences DSL is admitted now.

## 11. Access / disclosure candidate law

Audience selection and source authorization stay separate:

```text
resolve exact intended Principal(s)
→ validate current Organization/Product eligibility required for the kind
→ create personal Notification with bounded source-safe presentation
→ on deep link, re-authorize/reread source normally
```

Ordinary Permission may therefore be a **necessary eligibility condition** for a routed kind without being a sufficient reason to receive it.

If a historical recipient later loses source access, the Notification history is not rewritten; the source deep link fails/limits disclosure according to current authority.

## 12. Source-reference implication

The suspended D2 candidate's Work-only `source_work_ref` is falsified.

The later D2 model should be a **closed typed source union derived from this accepted census**, likely including only the source families actually needed by admitted kinds:

```text
NotificationSourceRef =
    MarketplaceInstallationRef
  | ListingIntentRef / PriceIntentRef
  | AvailabilitySubjectRef
  | EconomicAttributionOrReconciliationRef
  | MarketplaceSaleRef
  | BusinessOrderOrInvoicingIntentRef
  | FulfillmentExecutionRef
  | ShipmentRef
  | PostSaleResolutionRef
  | WorkRef
  | AuthorizationDecisionOrGovernedTargetRef
```

This is not a generic `{entity_type, entity_id}` graph. The final closed variants and source-qualified identity carriers are D2 work after the census is ratified.

## 13. Product / frontend implications after later authority repair

The census proves at least these future consumers, without authorizing them yet:

```text
topbar bell / unread presence
personal Inbox
Notification detail/deep-link behavior
Organization notification-routing configuration for A3 kinds
filters/grouping by Product-defined kind/source family
```

P6/P8 must later test notification density and whether visual grouping is needed. This census does **not** yet admit exact unread count, mark-all-read, bulk archive, per-user opt-out, e-mail/push or digest.

## 14. Correction sequence

```text
Reference study + H3              DONE / direction approved
Trigger + Audience Census         CANDIDATE — current gate
  ↓ operator ratification
D0-R trigger-scope correction
  ↓
D1-R explicit producer edges + routing ownership
  ↓
D2-R identity/source/routing re-derivation
  ↓
D3 communication/recovery
  ↓
D5 Product wire / Permissions
  ↓
D6 bell + Inbox + routing configuration
  ↓
D7 PostgreSQL/River/realtime realization
  ↓
D8 composed proof
```

## 15. Census decision candidate

> **Launch V1 Personal Notifications should cover 14 curated awareness families spanning channel operability, human consequential-action outcomes, availability safety, economic reconciliation attention, new sales, sale safety, business-system materialization, fulfillment entry/attention, shipment exceptions, post-sale attention, direct Work responsibility and Governance action/result. Audience resolution uses DIRECT_SOURCE when the occurrence already owns an exact human, OWNER_DERIVED when the source owner owns exact responsibility/authority semantics, and ORG_ROUTED only when the Organization must explicitly configure exact human recipients. Permissions never imply responsibility by themselves. Routine progression, every metric/state change, raw provider topics and infrastructure incidents remain outside the personal Inbox.**

**Next gate:** operator adjudicates this Trigger + Audience Census. Do not correct D0/D1/D2 or open D3 until this census is ratified.