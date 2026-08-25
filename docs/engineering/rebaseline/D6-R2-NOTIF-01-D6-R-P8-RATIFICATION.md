# D6-R2 P8 — Personal Notifications Operator Ratification

> **Status:** ACCEPTED / OPERATOR-RATIFIED / CURRENT EVIDENCE  
> **Operator adjudication:** 2026-08-23  
> **Current Product:** 106 operations / 31 ordinary Permissions / H-A-S  
> **Implementation:** blocked until D9

## 1. Operator LOCK

The operator reviewed the executable low-fidelity artifacts and locked:

```text
B00-R2  shell bell + bounded preview       LOCKED
B11     full personal Inbox                LOCKED
B12     Configurações > Notificações       LOCKED
```

Immutable reviewed HTML evidence:

- `qualification/d6-r2-wireframes/b00-r2-notifications.html`
- `qualification/d6-r2-wireframes/b11-notifications-inbox.html`
- `qualification/d6-r2-wireframes/b12-notification-routing-settings.html`

HTML may continue to self-report `candidate`; the operator LOCK lives here/P8 registry rather than mutating reviewed evidence after adjudication.

## 2. B00-R2 / B11 locked laws

Personal Notifications is exact-human awareness, not a business navigation mass or source authority.

```text
shell bell
→ unread presence: known-present | known-empty | unavailable
→ bounded preview
→ R128 full Inbox
```

Inbox structure:

```text
Ativas | Arquivadas
Todas | Não lidas | Lidas
closed NotificationKind filter
structured list
cursor continuation
Abrir origem
read/unread + archive/restore
```

Locked rules:

- bell/Inbox self-awareness does not require `notifications.manage`;
- no unread numeric count, total count, text search, mark-all-read, bulk archive, severity/priority or saved-view platform;
- source continuation and Notification awareness mutation remain separate;
- source continuation reauthorizes current source owner and never inherits capability from Notification;
- known empty != unavailable;
- stale awareness write and source denied/not-found remain distinct recoverable states;
- Organization switch invalidates scoped preview/Inbox state before presenting another Organization.

## 3. B12 routing Settings locked laws

Exactly ten fixed Product-defined `ORG_ROUTED` kinds appear under existing `Configurações > Notificações`.

For each row, selected people are exact future human recipients for that kind/current Organization. This is awareness routing only—not Permission/role/Work assignment, source authorization, ownership transfer or subscription DSL.

Locked rules:

```text
ListNotificationRoutes                current route truth
ListNotificationRouteRecipientCandidates
  → principal_id + display_name only
SetNotificationRoute                  current route write
```

- B12 requires `notifications.manage` without implicit `access.read`;
- one inline editor per affected row;
- candidate discovery is usability, never client authorization;
- configured route requires one-or-more recipients;
- explicit `UNCONFIGURED` is shown as `Sem configuração de destinatários`;
- `Remover configuração` writes UNCONFIGURED, not history deletion;
- historical recipient absent from current candidates displays `Destinatário não elegível · configuração anterior`, not opaque ID/invented name;
- current candidate cursor continuation exists; no speculative candidate-search Product API;
- `SetNotificationRoute` uses current route `If-Match` and server save-time eligibility validation;
- stale/rejected save preserves the local draft for correction;
- unavailable route read != ten unconfigured routes;
- Organization switch invalidates editor/candidate/route state;
- no role/Permission/IdP disclosure enters this surface.

## 4. Open evidence assumptions

Two bounded assumptions remain challengeable during later frontend work:

1. eligible routing candidates remain practically enumerable by cursor at real Organization scale without search;
2. the bounded recent preview remains sufficient for quick awareness while full triage belongs R128.

They are not Product facts and may reopen only the smallest affected UX/query decision if real evidence falsifies them.

## 5. Current binding

Exact current route/state/operation/access binding is the retained current P9 contract:

- shell bell / U01: no content-route identity;
- R128: `/org/:organizationId/notificacoes`;
- R129: `/org/:organizationId/configuracoes/notificacoes`.

This ratification owns structural LOCK only; mutable stage status belongs `docs/roadmap.md`.
