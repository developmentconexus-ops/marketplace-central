# NOTIF-01 D5-F4 — Recipient Discovery Global-Maximum Correction

> **Status:** ACCEPTED / OPERATOR-RATIFIED 2026-08-23
> **Trigger:** post-D5-F3 adversarial check of the complete human routing-settings consumer
> **Accepted inputs:** D0-R + D1-R + D2-R/R2/R3/R4 + D3-R/R1 + D5-F3 findings other than its four-operation count conclusion
> **Reviewed candidate:** [D5-R3 Product Operation Admission Table](D6-R2-NOTIF-01-D5-R3-OPERATION-ADMISSION-TABLE.md)
> **Operator direction:** Global Maximum correctness outranks minimizing diff or preserving any operation census
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Falsifier found after D5-F3

D5-R3 says the Notification-routing Settings client can compose configured recipient IDs with the already-admitted `ListOrganizationMembers` directory for human labels.

That composition is not complete under canonical W4 access semantics:

```text
ListOrganizationMembers
→ access.read
→ returns principal_id + principal_kind + display_name + role_keys

ListNotificationRoutes / SetNotificationRoute
→ notifications.manage
```

W4 explicitly defines Permissions as exact and non-hierarchical:

```text
notifications.manage != access.read
*.manage != *.read
```

Therefore a legitimate human actor with `notifications.manage` but not `access.read` can administer Notification routes semantically yet cannot discover human recipient identities through the existing directory. Requiring opaque Principal IDs is not a human-operable Product contract.

The defect is structural, not cosmetic:

> the accepted route write references Identity/Access-owned Principal identities, but the same admitted human job lacks the least-privilege read surface needed to discover selectable human identities.

This violates the frontend planning law that an opaque-ID parameter does not prove a coherent selector and violates the Engineering Method if solved only by coupling unrelated Permissions for convenience.

---

## 2. Alternatives challenged

### A. Bundle `access.read` with `notifications.manage`

**REJECTED — Local Maximum.**

It avoids a new operation but couples two distinct Product capabilities and grants broader member/access disclosure than the routing job requires. AccessRole bundling is allowed, but using it as a hidden prerequisite would make a current incidental role composition part of the Product contract.

### B. Let `notifications.manage` invoke `ListOrganizationMembers`

**REJECTED.**

The existing operation exposes `role_keys` and is owned as the Organization access directory. Broadening its access merely for one routing consumer either over-discloses access administration state or forces Permission-dependent response shapes.

### C. Embed recipient candidates in `ListNotificationRoutes`

**REJECTED.**

Route configuration state and the Organization human directory have different ownership/lifecycle/pagination. Embedding a directory inside route state would make Personal Notifications carry Identity/Access projection semantics and would make one settings screen shape the owner contract.

### D. Add one purpose-bounded Identity/Access recipient-discovery query

**ACCEPTED — GLOBAL MAXIMUM.**

The operation exposes only the smallest current human identity projection needed by the accepted route-administration job while Identity/Access remains the authority for Principal/Membership/eligibility and Personal Notifications remains the authority for route state.

---

## 3. Fifth Product operation — ratified admission direction

The final D5-R3 table must add one human-only Product Q, provisionally named:

```text
ListNotificationRouteRecipientCandidates
```

Exact path/schema spelling remains final D5 wire work, but the semantic contract is fixed:

```text
semantic owner:
  IdentityAccess

client:
  H only

ordinary access:
  notifications.manage

scope:
  exact Organization

meaning:
  list current human Principals who are currently Product-access eligible
  and current members of that exact Organization, for human recipient selection

minimum projection:
  principal_id
  display_name

collection mechanics:
  limit
  cursor
```

The projection does **not** expose by baseline:

```text
role_keys
Permissions
OIDC identity
email / username
Membership internals
eligibility_epoch
access administration state
```

It is not a generic Principal search API and does not replace `ListOrganizationMembers`.

---

## 4. Authority boundaries

Identity/Access owns:

```text
Principal identity
Principal kind
Product-access eligibility
Organization Membership
human presentation identity
```

Personal Notifications owns:

```text
NotificationRouteKey
route desired state
route revisions
selected recipient bindings
```

The candidate query supplies current human-readable references only. It does not create a second directory inside Personal Notifications.

`SetNotificationRoute(CONFIGURED)` remains authoritative for the write and MUST independently revalidate every submitted recipient against current state, including:

```text
same Organization
human Principal
current Membership
current Product-access eligibility
current source-read eligibility required by the selected ORG_ROUTED NotificationKind
```

The accepted D2-R2 eligibility-continuity epoch is captured server-side only after those checks. A candidate result observed earlier never authorizes a later write.

---

## 5. Source-read eligibility remains write-time authority

D5-F3's existing mapping from ORG_ROUTED kind to required source-read Permission remains accepted input for the final table.

The new Identity/Access candidate query does not attempt to encode NotificationKind-specific source-read policy into Identity/Access and does not return effective Permissions. Doing so would move Notification semantics into the identity substrate or over-disclose access state.

Therefore the baseline experience is:

```text
candidate discovery
→ current eligible human members

route write
→ exact selected NotificationKind
→ revalidate source-read eligibility
→ accept or reject safely
```

If later user evidence proves that route administrators need pre-save filtering by kind at scale, reopen the smallest owner contract then. Do not generalize now.

---

## 6. Product surface consequence

This correction supersedes only D5-F3's conclusion that four public operations are sufficient. D5-F3's other findings survive, including the bounded `NotificationKind` Inbox filter and the D2-R5 typed-result/requester-continuation repair.

Final NOTIF-01 operation direction becomes:

```text
ListMyNotifications
UpdateMyNotificationAwarenessState
ListNotificationRoutes
ListNotificationRouteRecipientCandidates
SetNotificationRoute
```

Consequential future wire census, if the final table is ratified with all current findings:

```text
99 + 5 = 104 Product operations
30 + 1 = 31 ordinary Permissions
Principal kinds remain H/A/S
```

Counts are consequences only.

---

## 7. Negative controls

This correction fails if it:

1. implies `notifications.manage -> access.read`;
2. grants `access.read` merely as a routing UI prerequisite;
3. exposes `role_keys`/Permissions/IdP data through the candidate projection;
4. makes Personal Notifications own or persist a duplicate Organization member directory;
5. treats the candidate list as authorization for `SetNotificationRoute`;
6. moves NotificationKind-to-source-Permission mapping into generic Identity/Access policy;
7. adds generic Principal search/filter DSL without a proved consumer;
8. adds a second routing Permission or `notifications.read` by symmetry;
9. changes accepted D2-R5, D3-R1, route-history or source-read eligibility findings except where later evidence independently falsifies them.

---

## 8. Proof contract for later D5 wire work

The final OAD/proof must falsify at least:

```text
notifications.manage without access.read
→ can list recipient candidates

candidate response
→ contains principal_id + display_name only from the admitted identity projection
→ does not disclose role_keys / Permissions

non-human Principal
→ absent

other-Organization Principal
→ absent

currently ineligible/non-member human
→ absent

candidate becomes invalid after list
→ SetNotificationRoute rejects under current authority

source-read Permission missing for selected NotificationKind
→ SetNotificationRoute rejects

candidate query
→ does not grant source access or route authority
```

A string-presence check is not sufficient proof; executable contract/auth negative controls belong to the later OAD authoring/proof step.

---

## 9. Gate

```text
D5-F3 bounded kind/result/continuation findings    RETAINED
D5-F3 four-operation conclusion                    SUPERSEDED BY D5-F4
D5-F4 recipient-discovery Global Maximum           ACCEPTED / OPERATOR-RATIFIED
D2-R5 typed result/requester continuation          CANDIDATE / NEXT OPERATOR GATE
D3-R2 bounded feed-forward                         BLOCKED BY D2-R5
D5-R3 final operation table                        BLOCKED BY D2-R5/D3-R2
canonical Product OAD                              UNCHANGED — 99/30
final expected NOTIF-01 wire consequence           104 operations / 31 Permissions
D6 / D7 / D8                                       BLOCKED for NOTIF-01
Product implementation                             BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates only the already-open D2-R5 typed-result/requester-continuation repair. After D2-R5 and its bounded D3 feed-forward close, rewrite/finally ratify D5-R3 with all five operations before any canonical OpenAPI modification.