# NOTIF-01 — Global Notification Reference Study

> **Status:** REFERENCE STUDY / NO PRODUCT AUTHORITY
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
- Shipping management: https://developers.mercadolivre.com.br/pt_br/registre-o-seu-aplicativo/gerenciamento-de-envios

Observed patterns:

1. Mercado Livre exposes distinct topics such as `orders_v2`, `shipments`, messages and post-purchase claims/actions.
2. `orders_v2` covers creation/changes to confirmed sales; shipments covers shipment changes; claims can notify immediately on claim creation/action.
3. The notification payload points to a resource that the consumer must `GET` for authoritative detail.
4. This strongly supports MPC's existing law that a notification/signal must not become the source-domain truth.
5. Provider notification topics are evidence acquisition, not MPC Personal Notification kinds by automatic 1:1 mapping.

### ShipStation

Source:

- Order Alerts: https://help.shipstation.com/hc/en-us/articles/360025869112-Order-Alerts

Observed patterns:

1. ShipStation exposes specific operational alert classes, such as combine-order opportunities, low inventory thresholds and automation-rule alerts.
2. Alerts are visible through a badge near the profile/Orders UI.
3. It does not model every order/shipment state transition as an alert; alert taxonomy is curated for operator relevance.

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

1. **Multiple producer families are legitimate.** Sales, fulfillment/shipping, post-sale/claims, assigned work and system/action outcomes can all have real consumers.
2. **Curated trigger taxonomy beats event-per-CRUD.** Normal low-value transitions must not flood the Inbox.
3. **Audience is part of the product design.** Permission alone is not responsibility; new-sale and shipment notifications need an explicit recipient rule.
4. **Notification is not source truth.** Opening the item must re-read/re-authorize the source.
5. **Attention surfaces and operational workspaces coexist.** The notification feed tells the user what happened; Sales/Fulfillment/Post-Sale/Work remain the places where domain work is understood/executed.
6. **Noise control is a first-class concern.** External products expose per-category routing/preferences because one undifferentiated feed degrades at scale.

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

**RECOMMENDED.** Each notification kind is Product-defined, owned by Personal Notifications as awareness semantics, backed by an explicit source-owner occurrence and admitted only when a user need is proven. Direct-recipient occurrences can target a Principal directly; organization-wide operational occurrences require a bounded Organization notification-routing rule that resolves exact human Principals.

This preserves semantic owners while serving the human Product globally.

## 5. Candidate Notification taxonomy for Launch V1

This is a study candidate, not accepted authority.

| Candidate kind | Source owner | Human value | Candidate disposition |
|---|---|---|---|
| New confirmed marketplace sale | Marketplace Sales | operator knows new commercial/operational work entered MPC | **ADMIT candidate** |
| Fulfillment becomes actionable / enters physical execution | Fulfillment Lifecycle | fulfillment team knows there is work ready to start | **ADMIT candidate**, but one meaningful entry event rather than every checkpoint |
| Dispatch deadline at risk / dispatch action materially due | Fulfillment Lifecycle | avoids missed provider obligation | **ADMIT candidate** |
| Physical conference discrepancy / fulfillment blocked | Fulfillment Lifecycle | immediate exception awareness | **ADMIT candidate** |
| Shipment enters material exception | Fulfillment Lifecycle / shipment observation authority | operations knows delivery diverged | **ADMIT candidate** |
| Material post-sale resolution opens/requires attention | Post-Sale Resolution | cancellation/return/refund consequence needs awareness | **ADMIT candidate** |
| Work assigned/reassigned to exact human | Operational Work | direct personal responsibility | **ADMIT** |
| Approval/authorization becomes personally actionable or a decision result matters to the initiator | Controlled Action Governance | approver/initiator awareness | **STUDY / likely admit if exact recipient is provable** |
| Long-running controlled action completes, rejects, becomes ambiguous/divergent | action-owning domain | initiator/responsible operator needs result awareness | **STUDY / likely admit where lineage gives exact recipient or configured audience** |
| Marketplace Installation health/reputation materially degrades | Marketplace Portfolio | operator/admin may need containment | **STUDY** |
| Normal delivered shipment | Fulfillment/shipment observation | low marginal action value | **DEFER by default** |
| Every fulfillment checkpoint | Fulfillment Lifecycle | redundant with workspace, high noise | **REJECT by default** |
| Every listing/availability/sync change | owning domain | high-volume low-value feed | **REJECT by default** |
| Every analytics/market/economic metric change | intelligence/economics | feed noise, no direct action boundary | **REJECT by default** |

Buyer Q&A/chat/message management remains outside current Product 1.0 scope unless D0 is explicitly reopened for that separate capability; provider messaging topics do not silently add it.

## 6. Audience model finding

The multi-producer requirement exposes a real second gap: not every occurrence naturally carries one exact recipient.

Examples:

```text
Work assigned to Principal P
→ direct recipient P

New sale confirmed
→ no single Principal is naturally encoded in Sale truth

Shipment exception
→ may concern Marketplace Operations, Fulfillment, or a configured responsible set
```

Therefore the recommended Global-Maximum model is **not** to infer recipients from ordinary Permissions and not to broadcast to every Organization member.

Candidate bounded model:

```text
NotificationKind
  ├─ direct recipient rule
  │    source occurrence carries exact human Principal
  │
  └─ organization routing rule
       Product-defined kind
       → configured exact human Principal recipients
```

Candidate YAGNI limits:

- Product-defined kinds only; no custom expression language.
- Exact Principal recipients only in baseline; no arbitrary groups/nested routing DSL.
- Ordinary `Permission` does not automatically imply notification responsibility.
- No user-authored generic subscriptions yet.
- No e-mail/push channel preferences; routing is for the in-app Inbox only.
- A future per-user opt-out/preferences feature requires separate evidence, but the data model must not make it impossible.

This likely expands the Personal Notifications owner from awareness lifecycle only to also own **bounded Notification routing configuration** for Product-defined kinds. That requires D1 correction before D2 is rederived.

## 7. Source-reference implication

The suspended D2 candidate's Work-only `source_work_ref` is too narrow if multiple producer owners are admitted.

The likely later D2 direction is a **closed typed source union**, derived only from accepted notification kinds, for example:

```text
NotificationSourceRef =
    MarketplaceSaleRef
  | FulfillmentExecutionRef
  | ShipmentRef
  | PostSaleResolutionRef
  | WorkRef
  | ...only later accepted producers
```

This is not a generic `{entity_type, entity_id}` graph. The union must remain closed and source-qualified where the underlying identity requires it.

The final union must not be frozen until the trigger/audience census is operator-ratified.

## 8. Recommended bounded correction sequence

```text
Reference study — this document
  ↓
Trigger + Audience Census
  ↓
D0-R trigger-scope correction
  ↓
D1-R explicit producer edges + bounded routing ownership
  ↓
rederive D2 source/occurrence/recipient model
  ↓
D3 communication/recovery
  ↓
D5 wire
  ↓
D6 bell/Inbox
```

Do not continue the current D2 candidate as if Work were the only producer.

## 9. Decision candidate

> **Personal Notifications should be a curated cross-product awareness capability, not a Work-only feature and not a generic event hub. Launch V1 should admit a bounded set of Product-defined notification kinds from multiple semantic owners, with explicit trigger contracts and explicit audience resolution. Direct personal occurrences target their Principal directly; operational broadcast occurrences use a small Organization routing configuration owned by Personal Notifications to resolve exact human recipients. Notifications remain durable personal awareness state and deep-link back to current authorized source truth.**

**Next gate:** operator reviews this reference study and the H3 direction. Authority documents are not corrected until that design direction is approved.