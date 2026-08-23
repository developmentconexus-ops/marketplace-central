# NOTIF-01 D1-R — Producer Edges & Notification Routing Boundary Correction

> **Status:** OPEN — D1 BOUNDARY CORRECTION CANDIDATE / OPERATOR ADJUDICATION REQUIRED
> **Trigger:** operator-approved D0-R Product scope + operator-approved Trigger/Audience Census
> **Parent accepted authority:** [D1 — Domains / Boundaries](D1-DOMAINS-BOUNDARIES.md) + [NOTIF-01 original bounded amendment](D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md)
> **Evidence:** [Global Notification Reference Study / Trigger + Audience Census](D6-R2-NOTIF-01-REFERENCE-STUDY.md)
> **Current Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Why D1 reopens

The original NOTIF-01 D1 amendment correctly created **Personal Notifications** as a supporting semantic owner, but admitted only:

```text
Operational Work → Personal Notifications
```

D0-R now admits fourteen curated Launch-V1 awareness families from multiple accepted owners. Keeping a Work-only edge would either leave proved Product consumers unowned or force ordinary Sales/Fulfillment/Post-Sale facts into Work.

D1-R therefore corrects only semantic ownership and allowed dependencies. It does not select event names, Q/C/E/P, identity/schema, Product API, Permissions, River, PostgreSQL layout, SSE, topbar behavior or implementation.

---

## 2. Personal Notifications remains one supporting semantic owner — CANDIDATE

No new business workflow domain is created.

**Personal Notifications** owns only:

```text
bounded Product NotificationKind vocabulary
personal awareness lifecycle
exact historical Notification recipient
ORG_ROUTED notification-routing configuration
recipient resolution from that routing configuration
per-recipient awareness deduplication/suppression semantics
Notification-local source correlation needed for explanation/navigation
```

It explicitly does **not** own:

```text
whether a Sale is confirmed or unsafe
whether Fulfillment is actionable/blocked/due
whether Shipment is exceptional
whether materialization converged/diverged
whether Economics requires reconciliation
whether Governance authorizes an action
who is assigned a Work
source-domain deadlines/policies/business truth
ordinary Product access / Permission semantics
source mutation or cross-domain workflow progression
provider protocol or webhook meaning
runtime delivery/retry/queue mechanics
```

The source owner must commit the source meaning first. Personal Notifications may represent that meaning as awareness but may not independently reinterpret source state to manufacture an attention occurrence.

---

## 3. NotificationKind ownership law — CANDIDATE

`NotificationKind` is a **bounded awareness vocabulary owned by Personal Notifications**, derived from the operator-approved Product census.

This does not transfer source semantics.

Example:

```text
Marketplace Sales owns:
  “this is the first authoritative confirmation of Sale S”

Personal Notifications owns:
  awareness classification NEW_MARKETPLACE_SALE
```

Personal Notifications cannot inspect arbitrary Sales state and decide by preference that another transition is `NEW_MARKETPLACE_SALE`. The accepted source edge must expose the already-qualified source occurrence.

Likewise:

```text
Fulfillment owns:
  actionable / blocked / dispatch-attention meaning

Personal Notifications owns:
  FULFILLMENT_ACTIONABLE / FULFILLMENT_ATTENTION awareness classes
```

Adding a new `NotificationKind` whose source meaning is not already owned by an accepted domain requires the smallest responsible D0/D1 reopen; a provider topic or UI preference cannot create it silently.

---

## 4. Explicit producer-edge set — CANDIDATE

The fourteen admitted awareness families require **ten explicit producer-owner edges** into Personal Notifications.

| Source owner → Personal Notifications | Admitted awareness family/families | Source owner retains authority for |
| --- | --- | --- |
| **Marketplace Portfolio → Personal Notifications** | `MARKETPLACE_INSTALLATION_ATTENTION` | Installation lifecycle, effective operability/posture and material channel-attention boundary |
| **Marketplace Offering Operations → Personal Notifications** | `OFFERING_ASYNC_ACTION_RESULT` | ListingIntent/PriceIntent intent, effect/convergence/result semantics and exact initiator lineage when present |
| **Availability Control → Personal Notifications** | `AVAILABILITY_ATTENTION` | Sellable Availability, synchronization/convergence and material attention boundary |
| **Commercial Economics → Personal Notifications** | `ECONOMIC_RECONCILIATION_ATTENTION` | attribution/reconciliation state and when bounded human resolution is required |
| **Marketplace Sales → Personal Notifications** | `NEW_MARKETPLACE_SALE`, `SALE_ATTENTION` | canonical Sale confirmation, sale-safe-handling/stop meaning and source-qualified Sale identity |
| **Business-System Materialization → Personal Notifications** | `MATERIALIZATION_ATTENTION` | Business Order/Invoicing materialization, ambiguity/rejection/divergence/external-required meaning |
| **Fulfillment Lifecycle → Personal Notifications** | `FULFILLMENT_ACTIONABLE`, `FULFILLMENT_ATTENTION`, `SHIPMENT_EXCEPTION` | physical readiness/execution, dispatch attention, provider-requirement closure and Shipment observation/exception meaning |
| **Post-Sale Resolution → Personal Notifications** | `POST_SALE_ATTENTION` | cancellation/return/refund consequence coordination and seller-relevant attention transition |
| **Operational Work → Personal Notifications** | `WORK_ASSIGNMENT` | Work responsibility, assignment/reassignment, escalation and Work lifecycle |
| **Controlled Action Governance → Personal Notifications** | `AUTHORIZATION_ACTION_REQUIRED`, `AUTHORIZATION_DECISION_RESULT` | delegation/grant semantics, exact decision authority, authorization decision/context and requester/initiator lineage where owned |

No other producer edge is admitted by symmetry.

In particular, Launch V1 does **not** admit direct Personal Notification edges from:

```text
Product & Channel Readiness
Market Intelligence
Marketplace Performance Intelligence / read projections
raw D4 provider adapters
D7 runtime/observability mechanics
```

Those areas may produce source-owned Work or future evidence, but they do not gain personal-notification edges without a new proved consumer and bounded reopen.

---

## 5. Producer occurrence law — CANDIDATE

Each producer edge carries only a **committed source-owned attention occurrence** sufficient for downstream awareness.

D1-R requires:

1. the producer owns the predicate that makes the occurrence material;
2. repeat reads/polls of the same source condition do not become new semantic occurrences;
3. the occurrence preserves enough stable source correlation for later D2/D3 deduplication;
4. provider-native webhook/topic arrival is not itself the Product occurrence;
5. Personal Notifications never calls into producer private state to reverse-engineer whether the event “should count”.

D3 later decides whether the public semantic boundary is materialized as `E`, `Q`, `P`, or another accepted communication form. D1-R only makes the dependency legal.

---

## 6. Audience ownership — CANDIDATE

The three operator-approved audience strategies have different semantic owners.

### 6.1 DIRECT_SOURCE

When the source occurrence already owns exact human lineage, **the producer owns recipient meaning** and exposes the exact Principal reference with the occurrence.

Baseline families:

```text
OFFERING_ASYNC_ACTION_RESULT
→ exact initiating human when authoritative initiator lineage exists

WORK_ASSIGNMENT
→ exact assigned/reassigned human Principal

AUTHORIZATION_DECISION_RESULT
→ exact requester/initiating human when Governance owns that lineage
```

Personal Notifications consumes that recipient; it does not replace it with routing configuration.

### 6.2 OWNER_DERIVED

When the source owner itself owns exact responsibility/authority semantics, **the source owner resolves the exact Principal set**.

Baseline family:

```text
AUTHORIZATION_ACTION_REQUIRED
→ Controlled Action Governance resolves exact currently valid decision Principal(s)
```

Personal Notifications must not reproduce Governance delegation/grant logic or infer approvers from ordinary Permissions.

### 6.3 ORG_ROUTED

When a material occurrence has no source-owned exact recipient, **Personal Notifications owns bounded Organization routing semantics**:

```text
(Organization, Product-defined NotificationKind)
→ configured exact human Principal set
```

A3 applies to the admitted operational families whose source truth has organization relevance but no natural personal addressee.

Routing configuration means only **who should receive in-app awareness for that Product-defined kind**. It is not:

```text
business responsibility assignment
source-domain policy
ordinary access control
role membership
workflow routing
provider subscription configuration
e-mail/push delivery preference
```

No source owner may write Personal Notifications routing configuration merely because it emits one of the kinds.

---

## 7. Access eligibility remains separate authority — CANDIDATE

Notification recipient selection and Product authorization are separate meanings.

Binding law:

```text
recipient selection
≠
source access grant
```

Personal Notifications may consume canonical current Organization/Principal access eligibility needed to avoid creating or disclosing an invalid routed item, but it does not own Membership, AccessRole, Permission or source-object authorization semantics.

Therefore:

- ordinary Permission can be a necessary eligibility condition without becoming the recipient selector;
- ORG_ROUTED never means “all Principals with Permission X”;
- no missing routing configuration falls back to all admins/members;
- opening a Notification re-enters normal source authorization;
- current access changes do not rewrite historical recipient identity.

Exact identity/access carriers and persistence enforcement remain D2/D5/D7 work.

---

## 8. ORG_ROUTED configuration semantics — CANDIDATE

Personal Notifications owns one bounded configuration meaning for Launch V1:

```text
Organization
+ admitted ORG_ROUTED NotificationKind
→ explicit set of exact human Principal recipients
```

Semantic laws:

```text
unconfigured
≠ configured empty
≠ configured recipient set
```

- **unconfigured** means routing has not been established and must remain explicit; no implicit broadcast occurs;
- **configured empty** may be admitted later only if D2/D5 prove a legitimate “intentionally nobody” consumer; D1-R does not require it by symmetry;
- configured recipients affect future occurrences only;
- routing changes never retarget historical Notifications;
- custom kinds, arbitrary expressions, nested groups, role-derived dynamic routing and subscription DSL remain outside Launch V1.

D2-R must choose the smallest persistent model that preserves these semantic distinctions.

---

## 9. Per-recipient awareness suppression ownership — CANDIDATE

The approved census proved one cross-owner duplicate-control rule:

```text
source occurrence S
→ source-family awareness would target Principal P
→ same S causally produces Work W already assigned to P
→ Personal Notifications suppresses the generic source-family Notification for P
→ WORK_ASSIGNMENT remains the personal awareness item for P
```

**Personal Notifications owns this suppression decision** because it is a duplicate-awareness concern, not source truth and not Work lifecycle.

It may consume only explicit source-occurrence/Work correlation supplied through accepted public boundaries. It must not infer equivalence by matching titles, timestamps or arbitrary entity IDs.

Suppression:

- is per recipient, never global;
- does not cancel, close or mutate Work;
- does not erase the source occurrence;
- does not prevent other configured recipients from receiving the source-family Notification;
- does not imply a generic cross-kind deduplication engine.

A broader cross-kind suppression framework requires separate evidence.

---

## 10. One-way correctness dependency — CANDIDATE

Personal Notifications is **never on the correctness path of the source business transition**.

```text
source owner commits valid source meaning
→ source meaning remains valid even if personal awareness is delayed
```

Conversely:

```text
Notification read/archive/routing state
→ never gates or mutates source-domain progress
```

D3/D7 later must make required Notification propagation recoverable, but D1-R does not turn Personal Notifications into a cross-owner transaction coordinator or workflow authority.

No source domain may require a Notification row to define whether its own business fact exists.

---

## 11. D1-R negative controls — CANDIDATE

The boundary correction must fail review if it:

1. creates `AnyDomain → Personal Notifications` instead of the explicit ten-edge set;
2. lets Personal Notifications independently determine Sale/Fulfillment/Shipment/Post-Sale source truth;
3. makes a producer write Personal Notifications private state or routing configuration;
4. makes Personal Notifications own Work assignment or Governance authority merely to resolve recipients;
5. derives ORG_ROUTED recipients from ordinary Permission holders or all admins/members;
6. uses a generic role/group/rule DSL as baseline routing authority;
7. adds Readiness/Market/Performance/runtime notification edges without a proved D0 consumer;
8. treats raw provider topics/webhooks as Product notification occurrences;
9. makes source-domain correctness depend on Notification read/archive/routing state;
10. turns suppression into source mutation or a generic cross-domain workflow engine;
11. silently adds e-mail/push/digest delivery responsibility to Personal Notifications;
12. selects event names, database tables, River jobs, API operations or frontend layout inside D1-R.

---

## 12. D1-R coherence result — CANDIDATE

The corrected model preserves accepted D1 law:

> one material meaning has one semantic authority; consumers use producer meaning through explicit public boundaries; mechanism does not acquire business authority.

Result:

```text
10 explicit source-owner edges
→ 14 Product-defined awareness families
→ 1 Personal Notifications supporting semantic owner
→ 3 bounded audience strategies
→ 1 bounded ORG_ROUTED configuration responsibility
→ 1 proved per-recipient Work-replacement suppression law
```

No new business lifecycle owner, generic event hub, generic subscription platform or deployment boundary is created.

---

## 13. Gate

```text
Reference study / H3          OPERATOR-APPROVED
Trigger + Audience Census     OPERATOR-APPROVED
D0-R Product scope            OPERATOR-APPROVED / ACCEPTED
D1-R producer/routing         READY FOR OPERATOR REVIEW
D2-R identity/data ownership  BLOCKED / prior Work-only candidate suspended
D3 communication/events       BLOCKED
D5 Product/OAD                BLOCKED
D6 bell/Inbox/settings        BLOCKED
D7 runtime                    BLOCKED
D8 proof                      BLOCKED
Product implementation        BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates only this D1-R boundary correction. If approved, close D1-R and rederive D2-R against the accepted ten-edge / fourteen-family / three-audience-strategy model. Do not open D3 or edit the OAD before D2-R is ratified.
