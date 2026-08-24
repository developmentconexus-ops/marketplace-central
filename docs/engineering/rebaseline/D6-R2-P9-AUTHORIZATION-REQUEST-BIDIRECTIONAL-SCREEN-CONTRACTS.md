# D6-R2 P9 — AuthorizationRequest + Notifications Bidirectional Screen Contracts

> **Status:** PROVED / CANONICAL — FINAL P9 BIDIRECTIONAL TRACE PASS
> **Stage:** D6-R2 Complete Frontend Realization Closure
> **Canonical Product:** 106 Product operations · 31 ordinary Permissions · Principal kinds H/A/S
> **Inputs:** [D5-R6 AuthorizationRequest OAD proof](D6-R2-NOTIF-01-D5-R6-AUTHORIZATION-REQUEST-OAD-WIRE-PROOF.md) · [Notification P8 ratification](D6-R2-NOTIF-01-D6-R-P8-RATIFICATION.md) · [B110 P8 ratification](D6-R2-P8-B110-APPROVALS-RATIFICATION.md) · [P9-F1 supersession](D6-R2-NOTIF-01-D6-R-P9-F1-SUPERSESSION-RATIFICATION.md)
> **Implementation:** BLOCKED UNTIL accepted D9
> **Next gate:** independent Fable review + finding adjudication before Global-Maximum closure or D7-R

## 1. Purpose and boundary

This is the final P9 implementation contract for the frontend surfaces materially affected by Personal Notifications and the AuthorizationRequest repair. It binds the operator-locked P8 experience to accepted Product authority in both directions:

1. every material frontend element has one exact owner/read/write authority; and
2. every human-consumer Product operation introduced or materially affected by this redesign has one frontend home or an explicit non-UI disposition.

P9 adds **no Product operation, Permission, identity, owner or screen**. The canonical OAD remains **106 Product operations** and **31 ordinary Permissions**. Runtime/router/database/jobs mechanics remain D7 work.

The former query-only P9 remedy remains **SUPERSEDED**. The original factual gap is closed by D1-R2 → D2-R6 → D3-R3 → D5-R6 → operator-locked B110.

```text
P9-F1 factual falsifier = RESOLVED
```

## 2. Final route and client-state contract

### 2.1 Exact frontend route identities

```text
G00-E / U01  global shell bell + bounded preview; no content-route identity
R128         /org/:organizationId/notificacoes
R129         /org/:organizationId/configuracoes/notificacoes
R110         /org/:organizationId/aprovacoes
R111 request /org/:organizationId/aprovacoes/solicitacoes/:authorizationRequestId
R111 history /org/:organizationId/aprovacoes/decisoes/:authorizationDecisionId
Work         /org/:organizationId/trabalho and /org/:organizationId/trabalho/:workId
```

R128 durable URL state:

```text
archive=active|archived
read=all|unread|read
kind=all|<NotificationKind>
```

`all` is frontend navigation state. The API adapter omits the corresponding Product filter instead of sending an unsupported `all` value. Cursor continuation is query/server state, not durable URL identity.

R110 durable local-lens state:

```text
view=para-decidir|historico
```

The history lens additionally carries admitted Product filters as URL state:

```text
decided_from=<RFC3339 instant>
decided_before=<RFC3339 instant>
```

The request and decision detail IDs remain path identity. Inline route editing in R129 and decision confirmation in R111 are not route identities.

### 2.2 State classes

```text
P9_STATE_CLASS:SERVER_STATE
```
TanStack Query/server state owns Notification lists/probes/routes/candidates, actionable AuthorizationRequests, AuthorizationDecisions, Work reads and all owner-source truth. Server responses, ETags and cursors are never copied into a competing durable client authority.

```text
P9_STATE_CLASS:URL_NAVIGATION_STATE
```
Organization-qualified content routes, R128 filters, R110 lens, history date filters, request IDs and decision IDs are navigation state so refresh/back/deep-link preserve the human task without inventing business state.

```text
P9_STATE_CLASS:LOCAL_EPHEMERAL
```
Bell/preview open state, B12 route-editor draft + selected candidate IDs, R111 pending confirmation outcome and a decision-attempt Idempotency-Key are local ephemeral UI state. Closing/changing Organization discards incompatible draft UI state; server truth is re-read.

```text
P9_STATE_CLASS:GLOBAL_WORKSPACE_CONTEXT
```
Current Organization is the only global workspace identity. Organization switch invalidates Organization-scoped query caches, cursors, preview/editor/confirmation state and incompatible deep links before new Organization truth is shown.

```text
P9_ORG_SWITCH:INVALIDATE_SCOPED_STATE
```

## 3. Screen Contract — G00-E bell

<!-- P9_SURFACE:G00-E -->

**Goal / flow.** Give any currently eligible authenticated human a bounded signal that personal Organization awareness may need attention without turning Notifications into a global navigation mass.

**Owner / read truth.** `PersonalNotifications` through `ListMyNotifications`; self/current Organization eligibility, H only, no ordinary `notifications.read` Permission and no `notifications.manage` prerequisite.

**Read mechanics.** The unread-presence probe is exact:

```text
archive_state=active
read_state=unread
limit=1
```

Successful non-empty means unread presence; successful empty means known no unread item; request unavailable is explicit knowledge unavailable. It never derives a count from pagination.

**Writes.** None. Bell activation, preview activation and source continuation do not mutate Notification awareness.

**Failure / recovery.** Unavailable != empty. Organization switch closes the utility surface and invalidates the previous Organization probe.

**Success / continuation.** Bell opens U01. Full triage continues to R128.

## 4. Screen Contract — U01 bounded preview

<!-- P9_SURFACE:U01 -->

**Goal / flow.** Show a small recent active-awareness projection for fast triage while keeping full filtering/pagination in R128.

**Owner / identity.** `PersonalNotifications`; Notification ID and source ref remain Notification-owned projections. The preview does not fetch broad source-owner data merely to enrich itself.

**Read mechanics.** `ListMyNotifications` with bounded recent-page behavior. Preview size is UX-bounded, not a total-count contract.

**Controls.** Bounded item source continuation + `Ver todas` to R128. No mark-all-read, bulk, search or independent filter platform.

**Source continuation.** Current source owner is reauthorized at navigation/read time. Open-source does not mark Notification read.

```text
P9_SOURCE_OPEN:REAUTHORIZE_CURRENT_OWNER
```

## 5. Screen Contract — R128 full personal Inbox

<!-- P9_SURFACE:R128 -->

**Goal / hierarchy.** Durable personal awareness triage in a structured vertical list: archive lens → read lens → bounded NotificationKind filter → Notification items → explicit source/awareness actions → cursor continuation.

**Owner / reads.** `PersonalNotifications` through `ListMyNotifications`. The list supports all 14 accepted NotificationKinds. F02/F14 typed outcomes are presentation-safe Notification projection; source-owner current truth is not copied into Inbox state.

**Writes.** `UpdateMyNotificationAwarenessState` is the only Notification mutation and controls only read/unread/archive/restore. It uses Notification-local stale-write protection. `Abrir origem` is a separate navigation event and cannot be coupled to awareness mutation.

**URL/state mapping.** `archive`, `read` and `kind` are URL navigation state; API queries map them to `archive_state`, `read_state`, `notification_kind`, omitting filters for frontend `all`. Cursor pages remain server/query state.

**Failures.** Known empty, request unavailable, stale awareness write and current source access denied/not-found are distinct recoverable states. A stale awareness write rereads Notification truth; a denied source continuation leaves the Notification intact.

**Forbidden convenience.** No text search, unread count, total count, mark-all-read, bulk archive, priority/severity, saved views or server-side DSL is introduced.

## 6. Screen Contract — R129 Notification routing Settings

<!-- P9_SURFACE:R129 -->

**Goal / hierarchy.** Under the existing Settings mass, manage the ten fixed ORG_ROUTED slots without exposing access-administration internals.

**Owners.** Current route truth + write owner: `PersonalNotifications`. Candidate discovery owner: `IdentityAccess`.

**Reads.** `ListNotificationRoutes` renders exactly ten rows. `ListNotificationRouteRecipientCandidates` supplies only `principal_id + display_name`, cursor-paginated. Candidate presence is discovery/usability only, never authorization.

**Local edit state.** One affected row expands inline. Current route/ETag remain server state; selected candidate IDs and validation messages are a local draft until save/cancel. Candidate cursor continuation is server/query state; no speculative candidate search operation is invented.

**Write.** `SetNotificationRoute` uses `If-Match`. Save requires CONFIGURED with one-or-more recipients; `Remover configuração` writes UNCONFIGURED. `CONFIGURED([])` is impossible in the frontend contract.

**Failures.** 412 preserves/reconciles draft after rereading route truth; 422 recipient validation keeps the row/draft open and identifies the rejected selection without Permission/role disclosure. Historical configured IDs absent from the current candidate set render `Destinatário não elegível · configuração anterior`, never an opaque ID or invented name.

## 7. Screen Contract — R110 actionable Approvals queue

<!-- P9_SURFACE:R110 -->

**Goal / hierarchy.** Answer one question: “what can this exact human legitimately decide now?” independently from Governance history.

**Owner / read.** `ControlledActionGovernance` through `ListMyActionableAuthorizationRequests`; H only + `governance.decide` + current Organization + current exact-human decision eligibility. Collection items are only PENDING/currently actionable requests.

**Representation.** Structured list + cursor continuation. No total/search/status/approver filters or bulk decision model.

**Permission separation.** `Para decidir` exists with `governance.decide`. `Histórico` exists with `governance.read`. One does not imply the other.

```text
P9_PERMISSION_LAW:GOVERNANCE_DECIDE_READ_INDEPENDENT
```

**Navigation.** `view=para-decidir|historico` is URL state. Selecting an actionable item navigates to the request-detail R111 route. F13 may deep-link directly to the same route, but possession of that ID grants nothing.

## 8. Screen Contract — R111 actionable request + immutable history

<!-- P9_SURFACE:R111 -->

### 8.1 Actionable request detail

**Read / owner.** `GetMyActionableAuthorizationRequest` revalidates `governance.decide`, current Membership and exact current decision eligibility. 200 returns one typed review-basis view and an **AuthorizationRequest ETag**. ID possession or F13 does not bypass this read.

**Typed evidence.** Exactly one of:

```text
P9_REVIEW_BASIS:listing_intent
P9_REVIEW_BASIS:price_intent
P9_REVIEW_BASIS:business_order_intent
P9_REVIEW_BASIS:invoicing_intent
```

The basis is immutable decision-purpose evidence, sufficient to understand the requested authorization without granting broad source-owner read. It is not current source truth and does not become a generic payload.

**Decision UX.** Approve/Reject enters inline confirmation while evidence stays visible. The pending outcome is local ephemeral state. Explicit confirmation creates one decision attempt.

**Write mechanics.** `CreateAuthorizationDecision` targets the exact `authorization_request_id`, sends outcome only, carries `If-Match` from the AuthorizationRequest ETag and an Idempotency-Key for the same semantic attempt.

```text
P9_DECISION_CONCURRENCY:IF_MATCH_AUTHORIZATION_REQUEST
P9_DECISION_RETRY:IDEMPOTENCY_KEY_SAME_SEMANTIC_ATTEMPT
P9_DECISION_AUTO_RETRY:FORBIDDEN
```

The Idempotency-Key is reused only when transport/acceptance outcome is ambiguous for the same request revision + same outcome attempt. A known no-effect 503 followed by explicit revalidation/confirm is a new human attempt. No generic consequential auto-retry exists.

**Failures.** 412 means request-local stale concurrency: reread Governance truth; if no longer actionable, show the bounded unavailable-for-decision state and do not leak Decision history unless separately authorized. 409 means current state/eligibility no longer admits the decision and is recovered through current Governance truth. 503 means material basis validity cannot be established now: no Decision was recorded; request remains pending; reread/revalidation is safe, blind decision retry is not. 403/404 fail closed.

**Success.** 201 shows Decision recorded. Governance still does not execute the underlying action. `Abrir origem` navigates to the current source owner, which reauthorizes its own read/action surface. `governance.decide` never implies source read.

### 8.2 Historical lens/detail

**Reads.** `ListAuthorizationDecisions` and `GetAuthorizationDecision` use independent `governance.read` authority (H/A/S). List filters are only admitted `decided_from`/`decided_before` + cursor. Historical detail renders immutable Decision + preserved review basis and has no Approve/Reject controls.

**History is not actionability.** A human who can decide one request does not gain Governance history; a reader of Governance history does not gain decision authority.

## 9. F13 and F14 continuation contract

### F13 — action required

```text
P9_F13:AUTHORIZATION_REQUEST_REF
P9_F13_NOTIFICATION_AUTHORITY:AWARENESS_NOT_CAPABILITY
```

`AUTHORIZATION_ACTION_REQUIRED` points to `AuthorizationRequestRef`. The deep link lands on R111 request detail and must succeed again through `GetMyActionableAuthorizationRequest`. Notification does not contain decision authority, request ETag, current basis validity or Permission proof. If another person already decided, request was invalidated, or eligibility changed, the deep link fails closed/reconciles to current actionable truth.

### F14 — decision result

```text
P9_F14:TARGET_ORIENTED
```

`AUTHORIZATION_DECISION_RESULT` remains target-oriented requester continuation. It does **not** point to `AuthorizationDecisionID` and does not require `governance.read` merely to continue work. The frontend resolves the target through its current semantic owner under that owner’s current read authority, then navigates to the accepted source route/workspace. For BusinessOrder/Invoicing targets, owner reads provide the sale/correlation required to reach the current Sale composition. If source read is not currently authorized, continuation fails closed without changing the Notification.

## 10. Zero-decider Operational Work disposition

Known PENDING + approval-required + zero currently eligible human deciders is an operational blocking condition, not authorization authority.

```text
P9_ZERO_DECIDER:WORK_NOT_AUTHORITY
P9_WORK_HOME:ListWork+GetWork
```

Governance may materialize/reconcile the accepted Operational Work condition whose origin is the AuthorizationRequest. Existing R100/R101 use `ListWork`/`GetWork` (plus already accepted Work controls) to make the blocked condition operable. Work may assign/escalate responsibility for fixing the authority/delegation gap; assignment never grants permission to decide and Work never becomes a fallback approver. When a valid decider exists again, the Work obligation reconciles/closes according to owner truth.

## 11. Frontend → backend bidirectional trace

| Surface | Material frontend responsibility | Exact owner / wire |
| --- | --- | --- |
| G00-E | unread presence | PersonalNotifications / `ListMyNotifications(active, unread, limit=1)` |
| U01 | bounded recent awareness + full-Inbox continuation | PersonalNotifications / `ListMyNotifications` |
| R128 | filters/list/cursor | PersonalNotifications / `ListMyNotifications` |
| R128 | read/unread/archive/restore | PersonalNotifications / `UpdateMyNotificationAwarenessState` |
| R128/U01 | source continuation | current source owner reauthorization; no Notification mutation |
| R129 | fixed route truth | PersonalNotifications / `ListNotificationRoutes` |
| R129 | recipient discovery | IdentityAccess / `ListNotificationRouteRecipientCandidates` |
| R129 | configure/unconfigure | PersonalNotifications / `SetNotificationRoute` + If-Match |
| R110 | exact-human actionable queue | ControlledActionGovernance / `ListMyActionableAuthorizationRequests` |
| R111 request | current actionable detail + typed basis | ControlledActionGovernance / `GetMyActionableAuthorizationRequest` + ETag |
| R111 request | authorize/reject | ControlledActionGovernance / `CreateAuthorizationDecision` + If-Match + Idempotency-Key |
| R110 history | immutable history list | ControlledActionGovernance / `ListAuthorizationDecisions` |
| R111 history | immutable history detail | ControlledActionGovernance / `GetAuthorizationDecision` |
| R100/R101 | zero-decider operational obligation | OperationalWork / existing Work reads/controls; never authority |

No frontend region owns copied source business truth merely because it presents a bounded snapshot.

## 12. Backend → frontend human-operation homes

The redesign has exactly ten human-facing Notification/Governance operations and no orphan:

```text
P9_OP_HOME:ListMyNotifications
P9_OP_HOME:UpdateMyNotificationAwarenessState
P9_OP_HOME:ListNotificationRoutes
P9_OP_HOME:ListNotificationRouteRecipientCandidates
P9_OP_HOME:SetNotificationRoute
P9_OP_HOME:ListMyActionableAuthorizationRequests
P9_OP_HOME:GetMyActionableAuthorizationRequest
P9_OP_HOME:CreateAuthorizationDecision
P9_OP_HOME:ListAuthorizationDecisions
P9_OP_HOME:GetAuthorizationDecision
```

Homes:

- `ListMyNotifications` → G00-E/U01/R128;
- `UpdateMyNotificationAwarenessState` → R128 item actions;
- routing reads/write → R129;
- actionable request list → R110 Para decidir;
- actionable request detail/write → R111 request route;
- decision history reads → R110 Histórico/R111 history route.

Internal D3 owner→Governance request creation/invalidation/revalidation is deliberately **not** a frontend operation. Zero-decider Work uses existing Work operations rather than adding a screen-shaped Governance/Work API.

## 13. Original P9-F1 closure

The old factual defect was: a legitimate `governance.decide` human could receive action-required awareness yet lacked a purpose-bounded current approval object/revision needed to review and decide without broad source read.

Resolution is structural, not query-only:

```text
D2-R6  AuthorizationRequest = canonical pre-decision episode identity/lifecycle
D3-R3  duplicate-safe intake + current eligibility + material validity + invalidation/recovery
D5-R6  actionable list/detail + request-local decision wire + four typed review bases
P8     B110 actionable/history experience operator-LOCKED
P9     exact route/state/owner/write/failure bindings in this contract
```

Therefore the factual falsifier is resolved while the query-only remedy remains **SUPERSEDED**. No `ListMyActionableAuthorizations` Product operation exists or is required.

## 14. Negative controls / falsifiers

P9 fails if implementation introduces any of the following:

```text
P9_FORBID:TARGET_ETAG_AS_DECISION_AUTHORITY
P9_FORBID:GENERIC_REVIEW_PAYLOAD
P9_FORBID:NOTIFICATION_CAPABILITY_TOKEN
P9_FORBID:WORK_FALLBACK_APPROVER
P9_FORBID:F14_AUTHORIZATION_DECISION_ID
P9_FORBID:SOURCE_OPEN_MARKS_NOTIFICATION_READ
P9_FORBID:INBOX_TOTAL_OR_UNREAD_COUNT
P9_FORBID:INBOX_BULK_OR_SEARCH
P9_FORBID:CANDIDATE_PRESENCE_IS_AUTHORIZATION
P9_FORBID:CONFIGURED_EMPTY_NOTIFICATION_ROUTE
P9_FORBID:SECOND_GLOBAL_APPROVAL_HISTORY_DESTINATION
P9_FORBID:SCREEN_SHAPED_PARALLEL_OWNER_AUTHORITY
```

Additional implementation laws:

- hidden navigation is UX only; server authorization remains authoritative;
- unavailable/unknown never becomes empty/zero;
- client does not infer actionability from Notification/candidate/list presence after time passes;
- target drift validity is semantic owner revalidation, not generic ETag equality;
- `AuthorizationDecision` is terminal history, not the pending request lifecycle;
- Fable may challenge this contract, but no review finding becomes authority without adjudication.

## 15. Proof and next gate

The executable verifier is `scripts/verify-d6-r-p9-screen-contracts.mjs` and is part of both repository gate lanes. It verifies six affected surfaces, ten human-operation homes, four typed review-basis families, F13/F14 continuation, zero-decider Work disposition, client-state laws and twelve negative controls while the existing OAD/P8 proofs independently preserve 95/29, 99/30 and current 106/31 authority.

```text
P9_FABLE_GATE:REQUIRED_AFTER_P9
```

**D7-R remains BLOCKED**. The exact next gate is independent Fable review of the complete AuthorizationRequest package plus explicit finding adjudication. P9 proof alone does not declare the redesign Global-Maximum closed, reopen D7-R/D8-R, resume B10, authorize merge or authorize Product implementation.
