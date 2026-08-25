# D6-R2 P9 — Personal Notifications + Approvals Current Screen Contracts

> **Status:** PROVED / CURRENT CONSOLIDATED P9 AUTHORITY  
> **Product:** 106 operations / 31 ordinary Permissions / H-A-S  
> **P8 inputs:** current Notifications P8 ratification + B110 P8 ratification  
> **Implementation:** blocked until D9

## 1. Route / state identities

```text
G00-E / U01  shell bell + bounded preview; no content-route identity
R128         /org/:organizationId/notificacoes
R129         /org/:organizationId/configuracoes/notificacoes
R110         /org/:organizationId/aprovacoes
R111 request /org/:organizationId/aprovacoes/solicitacoes/:authorizationRequestId
R111 history /org/:organizationId/aprovacoes/decisoes/:authorizationDecisionId
Work         /org/:organizationId/trabalho[/ :workId]
```

Durable URL state:

```text
R128: archive=active|archived, read=all|unread|read, kind=all|<NotificationKind>
R110: view=para-decidir|historico
history lens: decided_from?, decided_before?
```

Frontend `all` maps to omission of the Product filter. Cursors remain query/server continuation state.

State ownership:

- TanStack Query/server state = owner reads, ETags/revisions, cursors;
- URL = Organization/route/lens/filters/IDs;
- local ephemeral = preview/editor/confirmation drafts and same-attempt idempotency key;
- Organization = only global workspace; switch invalidates incompatible state.

## 2. G00-E / U01 — personal awareness utility

Owner: PersonalNotifications.

Unread-presence probe:

```text
ListMyNotifications
archive_state=active
read_state=unread
limit=1
```

Non-empty = known unread present; empty = known no unread; unavailable remains unavailable. No count is inferred from pagination.

U01 shows a bounded recent awareness subset + source continuation + `Ver todas` to R128. Source open reauthorizes current source owner and never marks Notification read implicitly.

## 3. R128 — Personal Inbox

Owner/read: `ListMyNotifications`, H-only authenticated-self + current Organization Membership.

Structure:

```text
Ativas | Arquivadas
Todas | Não lidas | Lidas
NotificationKind filter
structured list
cursor continuation
```

Only Notification write: `UpdateMyNotificationAwarenessState` for read/unread/archive/restore under Notification-local stale-write protection.

Source continuation is navigation to current source owner; it is not a Notification mutation and possession of Notification ID/ref grants no source capability.

Distinct failures:

- known empty;
- read unavailable;
- stale awareness write → reread Notification;
- current source access denied/not-found → Notification remains intact.

Forbidden: text search, total/unread count, mark-all-read, bulk archive, priority/severity, saved views.

## 4. R129 — Notification Routing Settings

Route truth/write owner: PersonalNotifications. Candidate discovery owner: IdentityAccess.

Reads:

```text
ListNotificationRoutes
ListNotificationRouteRecipientCandidates → principal_id + display_name only
```

UI has exactly ten fixed ORG_ROUTED rows; one affected row expands inline. Current route/ETag remains server state; selected candidate IDs are local draft.

Write:

```text
SetNotificationRoute + If-Match
CONFIGURED → one-or-more recipients
UNCONFIGURED → explicit remove configuration
```

`CONFIGURED([])` is forbidden.

Failures:

- 412 stale route → reread truth while preserving/reconciling local draft;
- 422 recipient invalid → keep editor open and identify rejected selection without role/Permission disclosure;
- historical configured ID not in current candidate set → `Destinatário não elegível · configuração anterior`;
- unavailable route read != ten unconfigured routes.

`notifications.manage` does not imply `access.read`; candidate presence is usability, not authorization.

## 5. R110 — Approvals

Two independent lenses:

```text
Para decidir → ListMyActionableAuthorizationRequests / governance.decide / H only
Histórico    → ListAuthorizationDecisions / governance.read
```

Neither Permission implies the other.

Actionable list contains only Requests the exact current human may decide. No total/search/status/approver/bulk filters. Selecting a Request goes to R111 request detail. F13 may deep-link to the same route but grants no capability.

History list uses only admitted date/cursor filters and immutable Decision truth.

## 6. R111 request — actionable decision

Read:

```text
GetMyActionableAuthorizationRequest
→ current exact-human eligibility revalidation
→ typed ActionableAuthorizationRequestView
→ current AuthorizationRequest ETag
```

Exactly one review basis:

```text
listing_intent
price_intent
business_order_intent
invoicing_intent
```

Review basis is immutable decision-purpose evidence. It is not current source truth or generic payload.

Approve/Reject enters inline confirmation while evidence remains visible.

### Current write carrier

```http
POST /organizations/{organization_id}/authorization-requests/{authorization_request_id}:decide
Idempotency-Key: <same semantic attempt>

{
  "etag": "<current AuthorizationRequest StrongETag>",
  "outcome": "authorize|reject"
}
```

There is **no If-Match** and no target/target ETag in the body.

Failure/recovery:

```text
missing/invalid etag      → 422 validation error
stale Request revision    → 409 revision conflict → reread Governance truth
current state conflict    → 409 → reread current actionability
known validity unavailable→ typed 503, known no Decision recorded
403/404                   → fail closed
```

Idempotency-Key is reused only for the same semantic attempt when acceptance is ambiguous. Exact committed replay resolves before current revision-precondition evaluation. Consequential auto-retry is forbidden.

201 = AuthorizationDecision recorded. It does **not** mean source action executed. `Abrir origem` reauthorizes current source owner.

## 7. R111 history

`GetAuthorizationDecision` under independent `governance.read` renders immutable Decision + preserved review basis. No approve/reject controls.

A human allowed to decide one Request does not gain Governance history. A Governance-history reader does not gain decision authority.

## 8. F13 / F14

### F13 `AUTHORIZATION_ACTION_REQUIRED`

Source ref = `AuthorizationRequestRef`. Deep-link lands on actionable R111 Request and revalidates via Governance. Notification contains no Request ETag/current eligibility/material-validity/Permission proof.

### F14 `AUTHORIZATION_DECISION_RESULT`

Source ref = `AuthorizationTargetRef`, target-oriented requester continuation. It does not point to AuthorizationDecisionID and does not require `governance.read` merely to continue the user's source job. Current source owner reauthorizes its own read.

## 9. Zero-decider Work

Known pending approval-required + zero currently eligible human deciders may produce/reconcile explicit Operational Work originated from AuthorizationRequest.

Work owns assignment/escalation/work state only. Assignment never grants decision authority; Work never becomes fallback approver. When valid decision authority exists again, Work reconciles according to Governance/source truth.

## 10. Bidirectional home summary

| Human surface | Product owner / operations |
| --- | --- |
| G00-E/U01/R128 | PersonalNotifications / `ListMyNotifications`; R128 also `UpdateMyNotificationAwarenessState` |
| R129 | PersonalNotifications / `ListNotificationRoutes`, `SetNotificationRoute`; IdentityAccess / `ListNotificationRouteRecipientCandidates` |
| R110 Para decidir | ControlledActionGovernance / `ListMyActionableAuthorizationRequests` |
| R111 request | ControlledActionGovernance / `GetMyActionableAuthorizationRequest`, `CreateAuthorizationDecision` |
| R110/R111 history | ControlledActionGovernance / `ListAuthorizationDecisions`, `GetAuthorizationDecision` |
| R100/R101 zero-decider condition | OperationalWork existing reads/controls; never authorization |

No Product operation introduced or materially changed by this accepted redesign is orphaned from its human home; non-UI producer/runtime operations remain outside frontend by design.

## 11. Forbidden frontend authority

- Notification never grants source access/decision capability;
- frontend does not reconstruct actionable approval population from general Governance history;
- client does not infer recipients from roles/Permissions;
- decision confirmation does not locally assume source execution/convergence;
- no generic approval/workflow payload or Notification action bus;
- no blind consequential retry;
- no durable client copy of server ETags/cursors/business truth outside server/query state.

This document is the current P9 composition authority for Personal Notifications + Approvals. Earlier carrier-specific P9/Fable repair chronology remains Git history rather than a required read chain.
