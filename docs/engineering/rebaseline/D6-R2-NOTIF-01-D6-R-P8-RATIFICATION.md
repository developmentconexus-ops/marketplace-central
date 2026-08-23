# NOTIF-01 D6-R — P8 Operator Ratification

> **Status:** ACCEPTED / OPERATOR-RATIFIED
> **Operator adjudication:** 2026-08-23
> **Parent:** [D6-R Frontend Feed-Forward](D6-R2-NOTIF-01-D6-R-FRONTEND-FEED-FORWARD.md)
> **P8 ledger:** [D6-R2 P8 Block Ledger](D6-R2-P8-BLOCK-LEDGER.md)
> **Canonical Product wire:** 104 Product operations · 31 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Ratified P8 result

The operator reviewed the executable low-fidelity HTML artifacts and ratified the structural Notifications experience:

```text
B00-R2  topbar bell + bounded preview       LOCKED
B11     full personal Inbox                LOCKED
B12     Configurações > Notificações       LOCKED
```

Reviewed artifacts remain unchanged evidence:

- `qualification/d6-r2-wireframes/b00-r2-notifications.html`
- `qualification/d6-r2-wireframes/b11-notifications-inbox.html`
- `qualification/d6-r2-wireframes/b12-notification-routing-settings.html`

This record supersedes only the earlier **B12 RENDERED CANDIDATE / visual adjudication pending** status and corresponding next-action statements in the D6-R feed-forward and P8 ledger. It does not reopen or replace their structural content.

## 2. Operator-confirmed meaning of B12

Each of the ten rows is one fixed Product-defined `ORG_ROUTED` Notification kind.

For one row, the selected people are the exact human recipients who should receive **future personal Notifications of that kind in the current Organization**.

This is awareness routing only. It is **not**:

```text
Permission assignment
role assignment
Work assignment
source authorization
ownership transfer
subscription DSL
```

The selector is a human discovery aid. Candidate presence is not authorization.

`SetNotificationRoute` remains authoritative at save time and revalidates the submitted recipients under current Organization membership, Product eligibility and the kind-specific source-read eligibility contract before accepting the new route revision.

## 3. Locked B12 laws

The operator-approved B12 structure locks:

- exactly ten fixed `ORG_ROUTED` route rows;
- local placement under `Configurações > Notificações`, never a new global sidebar mass;
- access conditioned by `notifications.manage` without an implicit `access.read` grant;
- one inline editor per affected route row;
- current candidate discovery through `ListNotificationRouteRecipientCandidates` with only `principal_id + display_name`;
- `ListNotificationRoutes` for current route truth and `SetNotificationRoute` for configure/unconfigure writes;
- candidate discovery as usability only, never client-side authorization;
- configured routes require one or more recipients; `CONFIGURED([])` remains forbidden;
- **Sem configuração de destinatários** represents explicit `UNCONFIGURED`;
- **Remover configuração** writes `UNCONFIGURED`; it does not delete route history;
- historical recipients absent from current candidates render only as `Destinatário não elegível · configuração anterior` and are never silently reactivated;
- stale write and save-time recipient rejection preserve the local draft for correction;
- route/candidate request unavailable remains distinct from ten unconfigured routes;
- Organization switching invalidates editor, candidate pagination and route state before presenting new Organization truth;
- no role, Permission, IdP or access-administration disclosure is introduced;
- no speculative recipient-search API is added before the existing P12 scale assumption is falsified.

## 4. Proof preserved

B12 TDD evidence remains:

```text
RED  87ed8a1a12461cb348188a7eaee52c801148caca
     failed only because the B12 artifact did not yet exist

GREEN ae18ba0a9aa0faa4ea866e035c312e08b1e9c111
      d6_r_b12_route_slots=10/10
      d6_r_b12_editor=INLINE_ROW
      d6_r_b12_permission=notifications.manage
      d6_r_b12_access_read=FORBIDDEN
      d6_r_b12_candidate_disclosure=MINIMAL
      d6_r_b12_configured_empty=FORBIDDEN
      d6_r_b12_wireframe=PASS
```

Pre-ratification routed checkpoint `6232636a6986a761dad6c05afec0e8d4d9c58a5e` passed CI #567, pr-title #629 and the repository full gate.

The HTML snapshots intentionally continue to self-report `CANDIDATE`; the human `LOCK` lives in durable authority rather than rewriting reviewed evidence after adjudication.

## 5. Remaining assumptions

The D6-R assumption register remains open for P12 challenge:

1. current human routing candidates remain practically enumerable through cursor pagination without server-side search at realistic Organization scale;
2. the bounded recent preview remains sufficient for quick awareness while full triage belongs to R128.

Neither assumption is silently promoted to fact by this ratification.

## 6. Exact next action

P8 Notifications is complete. Open **P9 — exact Screen Contracts / bidirectional wire trace** for the locked Notifications surfaces:

```text
G00-E / U01  bell + bounded preview
R128         full personal Inbox
R129         Configurações > Notificações
```

P9 must bind each material read/write/state to exact frontend route/context, semantic owner, `operationId`, ordinary Permission or authenticated-self rule, client class, identity/scope requirements, stale/error behavior and source/effect trace.

Do not begin NOTIF-01 D7-R, D8-R, Product implementation or merge PR #61 before D6-R P9 closes unless the operator explicitly changes the program gate.