# NOTIF-01 D0-R — Personal Notification Trigger-Scope Correction

> **Status:** OPEN — PRODUCT-SCOPE CORRECTION CANDIDATE / OPERATOR ADJUDICATION REQUIRED
> **Trigger:** operator-approved H3 direction + operator-approved Trigger/Audience Census in [Global Notification Reference Study](D6-R2-NOTIF-01-REFERENCE-STUDY.md)
> **Parent accepted authority:** [NOTIF-01 Personal Notifications Authority Amendment](D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md) — original D0/D1 bounded baseline remains historical authority except where this correction is later operator-ratified
> **Current Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Why D0 reopens

The original NOTIF-01 D0 amendment proved that Personal Notification Inbox is a real Product 1.0 capability but used one initial trigger — `Operational Work` assignment/reassignment — as the only admitted Launch-V1 origin.

Frontend/user review supplied a material falsifier: users also need awareness when committed marketplace-operating facts become relevant while they are outside the source workspace, for example a new confirmed sale, fulfillment becoming actionable, shipment exceptions and post-sale attention.

The approved reference study rejected both local maxima:

```text
Work-only notifications
→ too narrow; forces routine operational facts into Work

AnyDomain/event-per-CRUD notifications
→ too broad; creates noise and a generic event hub
```

D0-R therefore corrects only Product scope. It does not yet choose producer semantic edges, data identity/schema, event contracts, API operations, Permissions, runtime mechanisms or frontend layout.

---

## 2. Corrected Product 1.0 responsibility — CANDIDATE

Product 1.0 includes **Personal Notifications** as a curated cross-product awareness capability for exact human Principals inside one Organization.

Purpose:

> allow a human to discover material committed MPC occurrences they specifically need to notice, even when they are outside the owning workspace, while the originating domain remains the truth and execution authority.

Launch V1 therefore requires all of the following Product outcomes:

```text
curated Product-defined notification kinds
+ exact human recipient resolution
+ durable personal Inbox state
+ unread/read
+ active/archived
+ truthful source/deep-link correlation
+ bounded Organization routing configuration where no exact recipient exists naturally
```

Notification remains awareness only. It does not become source truth, Work, authorization, Audit, acknowledgement, source access or cross-domain workflow progression.

---

## 3. Approved Launch-V1 awareness families — CANDIDATE D0 scope

The operator-approved Trigger/Audience Census admits these fourteen Product awareness families as Launch-V1 scope:

```text
1.  marketplace installation / channel operability attention
2.  human-initiated Offering consequential-action result
3.  Availability attention
4.  economic reconciliation attention
5.  new confirmed marketplace sale
6.  material Sale attention / safe-handling stop
7.  business-system materialization attention
8.  Fulfillment becomes actionable
9.  material Fulfillment attention
10. Shipment exception
11. material post-sale attention
12. Work assignment / reassignment
13. authorization action required
14. authorization decision result
```

D0 admits the **human awareness need**, not one notification per state transition. Later D1/D3 authority must derive the exact owner occurrence for every admitted family.

No provider webhook/topic automatically becomes one of these Product kinds. Provider evidence must first become committed meaning under the responsible MPC owner.

---

## 4. Attention-transition law — CANDIDATE

The Product notification model is based on **attention transitions**, not activity logging.

Examples:

```text
new Sale first becomes authoritative in MPC
→ awareness candidate

Fulfillment first becomes legitimately actionable
→ awareness candidate

same Fulfillment records ordinary separation/conference/packing progression
→ no new Notification merely because a checkpoint changed

Shipment first enters one material exception occurrence
→ awareness candidate

polling/rereading the same unresolved exception
→ no new Notification
```

A later materially different attention transition may create a distinct Notification when it represents a new human need, e.g. `FULFILLMENT_ACTIONABLE` followed later by a dispatch-risk attention transition.

---

## 5. Exact human audience requirement — CANDIDATE

Every personal Notification ultimately belongs to one exact human Principal.

D0-R recognizes three Product-level audience situations from the approved census, without assigning their semantic implementation owner yet:

```text
DIRECT_SOURCE
committed occurrence already identifies the exact human

OWNER_DERIVED
source authority already owns the responsibility/authorization semantics needed to resolve exact humans

ORG_ROUTED
material Organization occurrence has no natural exact recipient and requires bounded configured recipients
```

Binding Product laws:

- ordinary Permission may be required for eligibility but never implies notification responsibility by itself;
- no implicit broadcast to every Organization member, admin or Permission holder;
- recipient resolution never grants source access;
- current source authorization is re-evaluated when the user follows the Notification;
- route/recipient changes affect future awareness and do not rewrite historical recipients.

Exact ownership of DIRECT_SOURCE/OWNER_DERIVED/ORG_ROUTED semantics is D1-R work.

---

## 6. Organization notification routing becomes a Launch-V1 supporting capability — CANDIDATE

The approved census proves a real Product consumer for bounded routing configuration: events such as a new Sale or Shipment exception are important but do not naturally identify one exact human.

Product 1.0 therefore admits an **in-app Organization notification-routing capability** whose bounded user outcome is:

```text
Organization
+ Product-defined NotificationKind
→ configured exact human Principal recipients
```

This capability is intentionally narrower than a notification-preferences/subscription platform.

Launch V1 does not require:

```text
custom notification kinds
custom boolean/rule expressions
routing DSL
nested groups
role-derived automatic subscription
per-user arbitrary subscriptions
e-mail routing
push routing
digests
```

Later D1/D2/D5/D6 authority must derive ownership, persistence, API and settings UX before implementation.

---

## 7. Cross-product experience result — CANDIDATE

The Product promise is now:

```text
source workspace remains primary truth/work surface
                    │
material attention occurrence
                    │
                    ▼
               personal Inbox
                    │
                    ▼
                 topbar bell
                    │
                    ▼
source-specific deep link + current reread/auth
```

The Inbox is not a replacement for `Vendas`, `Expedição`, `Pós-venda`, `Trabalho`, `Aprovações`, Configurações or strategy surfaces.

Examples of legitimate different awareness moments for one sale lifecycle:

```text
NEW_MARKETPLACE_SALE
→ marketplace operations notices new demand

FULFILLMENT_ACTIONABLE
→ fulfillment team notices physical execution can begin

FULFILLMENT_ATTENTION / SHIPMENT_EXCEPTION
→ responsible operators notice a later material risk/divergence
```

These are distinct human handoffs; they do not imply one cross-owner workflow state.

---

## 8. Noise and duplicate Product laws — CANDIDATE

Launch V1 requires the later design to preserve these outcomes:

1. the same committed source occurrence does not create repeated semantic Notifications merely because acquisition/recovery/polling repeats;
2. routine state progression is not a notification feed;
3. a source occurrence that immediately causes Work already assigned to the same Principal may suppress the generic routed source alert for that recipient while preserving the mandatory `WORK_ASSIGNMENT` awareness;
4. suppression is per-recipient, so other configured humans may still receive the source-owner awareness;
5. read/archive does not reset because source state changed; a genuinely new source occurrence creates a new awareness item;
6. routing changes do not retroactively create/move/delete historical Notifications by default.

Exact occurrence IDs, suppression keys and delivery mechanics remain D2/D3/D7 work.

---

## 9. Explicit D0-R non-goals — CANDIDATE

D0-R does **not** make any of the following Product 1.0 requirements:

```text
notify every Product state transition
notify every provider webhook/topic
notify every analytics/market/economic metric movement
notify every Fulfillment checkpoint
normal delivered-shipment completion alert by default
buyer Q&A/chat/general marketplace messaging
seen state
exact unread numeric aggregate
mark-all-read
bulk archive
per-user arbitrary subscriptions/preferences
custom routing expressions / DSL
e-mail notifications
mobile/web push
digest delivery
generic notification template platform
generic EventStore/pub-sub platform
external broker
infrastructure/on-call incidents in the user Inbox
A/S human-style Inbox
Notification-triggered source mutation
```

These require separate evidence and owning-authority reopen if they later become real consumers.

---

## 10. D0-R negative controls — CANDIDATE

The corrected Product scope must fail review if it attempts to:

1. return to Work-only origin and thereby force ordinary sales/fulfillment/shipment facts into Work;
2. create `AnyDomain → Notifications` generic fan-out;
3. equate provider notification topics with MPC Product NotificationKinds;
4. infer recipients from ordinary Permissions alone;
5. broadcast ORG_ROUTED kinds to every member/admin when routing is absent;
6. make Notification the source truth or execution surface;
7. make Notification read/archive acknowledge or resolve the source;
8. let routing grant access to source objects;
9. create a generic subscription/rule engine before a proved consumer requires it;
10. preserve the current 99-operation/30-Permission census by force-fitting future Notification API into unrelated capabilities.

---

## 11. Coherence with accepted Product definition — CANDIDATE

This correction remains coherent with Marketplace Central as **Marketplace Operations Control Plane + Commercial Intelligence** because:

- the admitted awareness families correspond to already accepted operating jobs and owners rather than adding a new marketplace lifecycle;
- Notifications makes those accepted jobs discoverable across screens without absorbing their meaning;
- routing is supporting Product configuration, not a new business-policy engine;
- the product remains neither marketplace dashboard replacement nor ERP/WMS/TMS/CRM;
- Global Maximum serves all proved consumers, while YAGNI still excludes generic messaging/subscription/event infrastructure.

The Product wire remains unchanged until a later accepted D5 reopen legitimately derives operations/Permissions from corrected authority.

---

## 12. Gate

```text
Reference study / H3          OPERATOR-APPROVED
Trigger + Audience Census     OPERATOR-APPROVED
D0-R trigger-scope correction READY FOR OPERATOR REVIEW
D1-R producer/routing         BLOCKED
D2-R identity/data ownership  BLOCKED / prior candidate suspended
D3 communication/events       BLOCKED
D5 Product/OAD                BLOCKED
D6 bell/Inbox/settings        BLOCKED
D7 runtime                    BLOCKED
D8 proof                      BLOCKED
Product implementation        BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates only this D0-R Product-scope correction. If approved, close D0-R and open the bounded D1-R producer-edge + notification-routing ownership gate. Do not rederive D2 or open D3 before D1-R is ratified.
