# NOTIF-01 D0-R — Personal Notification Trigger-Scope Correction

> **Status:** ACCEPTED / OPERATOR-RATIFIED
> **Trigger:** operator-approved H3 direction + operator-approved Trigger/Audience Census in [Global Notification Reference Study](D6-R2-NOTIF-01-REFERENCE-STUDY.md)
> **Parent accepted authority:** [NOTIF-01 Personal Notifications Authority Amendment](D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md)
> **Current Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Correction

The original NOTIF-01 D0 amendment correctly admitted a Personal Notification Inbox but was too narrow when it treated `Operational Work` assignment/reassignment as the only Launch-V1 notification origin.

Operator review proved real consumers across the wider Product lifecycle: new confirmed Sales, Fulfillment handoffs/attention, Shipment exceptions, Post-Sale attention and other material occurrences.

The accepted Product direction rejects both local maxima:

```text
Work-only notifications
→ too narrow; would force routine Product facts into Work

AnyDomain / event-per-CRUD notifications
→ too broad; would create noise and a generic event hub
```

D0-R therefore corrects Product scope only. D1 owns semantic boundaries; D2 identity/data; D3 communication; D5 wire; D6 UX; D7 runtime.

## 2. Accepted Product 1.0 responsibility

Product 1.0 includes **Personal Notifications** as a curated cross-product awareness capability for exact human Principals inside one Organization.

Purpose:

> allow a human to discover material committed MPC occurrences they specifically need to notice, even while outside the owning workspace, while the originating domain remains truth and execution authority.

Launch V1 requires:

```text
curated Product-defined notification kinds
+ exact human recipient resolution
+ durable personal Inbox state
+ unread/read
+ active/archived
+ truthful source/deep-link correlation
+ bounded Organization in-app routing where no exact recipient exists naturally
```

Notification remains awareness only. It is not source truth, Work, authorization, Audit, acknowledgement, source access or workflow progression.

## 3. Accepted Launch-V1 awareness families

The approved census admits fourteen awareness families:

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

D0 admits the human awareness need, not one Notification per state transition. Provider webhook/topic arrival never becomes a Product NotificationKind automatically; provider evidence must first become committed meaning under the responsible MPC owner.

## 4. Attention-transition law

Notifications represent **attention transitions**, not an activity log.

```text
first authoritative new Sale
→ awareness

Fulfillment first becomes legitimately actionable
→ awareness

ordinary separation / conference / packing progression
→ no Notification merely because a checkpoint changed

Shipment first enters a material exception occurrence
→ awareness

reread/poll of the same unresolved exception
→ no new awareness
```

A later materially different attention transition may create another Notification when it represents a genuinely new human need.

## 5. Exact-human audience law

Every personal Notification ultimately belongs to one exact human Principal.

D0 recognizes three audience situations:

```text
DIRECT_SOURCE
source occurrence already identifies the exact human

OWNER_DERIVED
source authority owns responsibility/authorization semantics that resolve exact humans

ORG_ROUTED
material Organization occurrence has no natural exact recipient and uses bounded configured recipients
```

Binding laws:

- ordinary Permission may be necessary eligibility but never implies notification responsibility by itself;
- no implicit broadcast to all members, admins or Permission holders;
- routing/recipient selection never grants source access;
- following a Notification re-enters current source authorization;
- routing changes affect future awareness and never rewrite historical recipients.

Semantic ownership of these strategies is D1-R authority.

## 6. Accepted Organization routing capability

Launch V1 includes a bounded **in-app Organization notification-routing capability**:

```text
Organization
+ Product-defined NotificationKind
→ configured exact human Principal recipients
```

It is not a generic notification-preferences/subscription platform.

Launch V1 does not require custom kinds, boolean routing expressions, routing DSL, nested groups, role-derived automatic subscriptions, arbitrary per-user subscriptions, e-mail routing, push routing or digests.

## 7. Cross-product experience law

```text
source workspace remains primary truth/work surface
                    │
material attention occurrence
                    ▼
               personal Inbox
                    ▼
                 topbar bell
                    ▼
source-specific deep link + current reread/auth
```

The Inbox does not replace `Vendas`, `Expedição`, `Pós-venda`, `Trabalho`, `Aprovações`, Configurações or strategy surfaces.

A single Sale lifecycle may legitimately create different human handoffs:

```text
NEW_MARKETPLACE_SALE
→ marketplace operations notices new demand

FULFILLMENT_ACTIONABLE
→ fulfillment notices physical execution can begin

FULFILLMENT_ATTENTION / SHIPMENT_EXCEPTION
→ responsible operators notice later material risk/divergence
```

These are separate awareness moments, not one cross-owner workflow state.

## 8. Noise / duplicate laws

Launch V1 requires later authority to preserve:

1. same committed source occurrence never creates repeated semantic Notifications from polling/recovery/replay;
2. routine state progression is not a notification feed;
3. when the same occurrence immediately creates Work already assigned to Principal P, the generic routed source alert may be suppressed for P while `WORK_ASSIGNMENT` remains;
4. suppression is per recipient; other configured recipients may still receive source awareness;
5. read/archive never resets because source state changed; a new source occurrence creates a new item;
6. routing changes do not retroactively create, move or delete historical Notifications by default.

Exact occurrence identities, correlation and delivery mechanics remain downstream.

## 9. Accepted non-goals

D0-R does not require:

```text
notification for every Product transition / provider topic / metric movement
notification for every Fulfillment checkpoint
normal delivered-shipment completion alert by default
buyer Q&A/chat/general marketplace messaging
seen
exact unread numeric aggregate
mark-all-read
bulk archive
arbitrary subscriptions/preferences
routing DSL
e-mail / push / digest
generic template platform
generic EventStore/pub-sub
external broker
infrastructure/on-call incidents in user Inbox
A/S human-style Inbox
Notification-triggered source mutation
```

## 10. Negative controls

The Product scope is invalid if it:

1. returns to Work-only origin;
2. creates generic `AnyDomain → Notifications` fan-out;
3. equates provider topics with MPC NotificationKinds;
4. infers recipients from ordinary Permissions alone;
5. broadcasts ORG_ROUTED kinds to all members/admins when routing is absent;
6. makes Notification source truth or execution surface;
7. lets read/archive acknowledge or resolve source truth;
8. lets routing grant source access;
9. creates a generic subscription/rule engine without a proved consumer;
10. preserves 99 operations / 30 Permissions by force-fitting later Notification wire into unrelated capabilities.

## 11. Coherence result

D0-R remains coherent with Marketplace Central as **Marketplace Operations Control Plane + Commercial Intelligence**: it makes already accepted operating jobs discoverable across screens without adding a new marketplace lifecycle or absorbing existing owners.

The Product wire remains unchanged until later accepted D5 authority derives the exact operations/Permissions.

## 12. Gate

```text
Reference study / H3          OPERATOR-APPROVED
Trigger + Audience Census     OPERATOR-APPROVED
D0-R trigger-scope correction ACCEPTED / OPERATOR-RATIFIED
D1-R producer/routing         OPEN / CANDIDATE
D2-R identity/data ownership  BLOCKED / prior candidate suspended
D3 communication/events       BLOCKED
D5 Product/OAD                BLOCKED
D6 bell/Inbox/settings        BLOCKED
D7 runtime                    BLOCKED
D8 proof                      BLOCKED
Product implementation        BLOCKED UNTIL D9
```

**Exact next action:** adjudicate only D1-R producer-edge and routing ownership. Do not rederive D2 or open D3 before D1-R is ratified.
