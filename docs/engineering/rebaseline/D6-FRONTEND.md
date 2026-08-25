# D6 — Frontend

> **Status:** ACCEPTED / CLOSED BASELINE — current consolidated authority  
> **Current realization stage:** D6-R2 remains OPEN; see `docs/roadmap.md`  
> **Product input:** canonical OAD **106 operations / 31 ordinary Permissions / H-A-S**  
> **Frontend method:** `docs/development/frontend-product-experience-planning-method.md` v2.3

## 1. Purpose

D6 defines Marketplace Central's browser interaction/realization architecture. The frontend is a human-operable projection/client of accepted Product authority, never a second business authority.

Accepted implementation direction after D9:

```text
React + TypeScript
TanStack Router    → navigation / URL state
TanStack Query     → server state
openapi-typescript → generated Product wire types
openapi-fetch      → thin typed Product transport
```

No frontend Product implementation begins until D9 allows it.

## 2. Governing frontend invariants

1. Product/OpenAPI is the server business-wire authority; no hand-written shadow DTO/API.
2. Organization is explicit global workspace scope; no hidden/default tenant.
3. Marketplace Installation is explicit page/local context only where the Product job requires it; no ambient/default account.
4. Permission-conditioned visibility is usability only; server authorization is authoritative.
5. `Permission != business disposition != provider capability != Governance authorization`.
6. Known empty/zero, unknown, unavailable, partial, unsupported and stale remain distinguishable.
7. Accepted/pending/rejected/ambiguous/divergent/converged remain distinct where the owner can distinguish them.
8. Idempotency/precondition semantics survive retries; no blind replay after potentially accepted consequential writes.
9. External IDs remain source-qualified.
10. Read composition never becomes write authority.
11. Frontend never calls provider/business-system protocol directly to fill a Product-contract gap.
12. No screen-shaped Product API/BFF by convenience.
13. Human browser authentication uses server-side session + CSRF; browser JavaScript owns no OIDC access/refresh token.
14. Frontend evidence may falsify an upstream Product/contract assumption; then reopen the **smallest owning authority**, find Global Maximum, repair, and return to the block.

## 3. Frontend state ownership

| State class | Owner |
| --- | --- |
| server state | Product/server via TanStack Query |
| URL/navigation state | TanStack Router/browser |
| unsent form draft | current user editing session |
| ephemeral UI state | frontend only |

A fifth durable/global state class requires concrete evidence. No normalized global business-entity mirror, generic action bus or duplicate domain store is admitted.

Organization switch invalidates incompatible Installation/subject server + navigation state. URL/router state never becomes business authority.

## 4. Current Product/client surface

Current canonical Product input is:

```text
106 Product operations
31 ordinary Permissions
Principal kinds H / A / S
```

Important later accepted additions over the earlier D6 baseline are:

- Marketplace Performance Intelligence read surface (`performance.read`);
- Personal Notifications self Inbox + Organization routing (`notifications.manage` only for route admin);
- actionable AuthorizationRequest reads and request-scoped decision flow under `governance.decide`.

No generic Analytics/Strategy API, Notification delivery-channel platform or approval-workflow engine was added.

## 5. Information architecture / shell

The current operator-ratified primary mental model is:

```text
VISÃO GERAL
  Visão geral

OFERTA
  Preparação
  Anúncios
  Preços
  Disponibilidade

OPERAÇÃO
  Visão operacional
  Vendas
  Expedição
  Pós-venda

ESTRATÉGIA E INTELIGÊNCIA
  Performance
    Resumo
    Anúncios
    Mídia
  Mercado
  Economia

CONTROLE
  Trabalho
  Aprovações

CONFIGURAÇÕES
  ...
```

Personal Notifications is a shell utility/Inbox concern rather than a new primary business area. Notification routing settings live under appropriate configuration/admin UX. IA is a human mental model, not D1 package taxonomy.

### 5.1 Shell context

Accepted structural baseline from B00:

- Organization = only global workspace;
- page-local exact Marketplace Installation selector where material;
- no hidden first/default Installation;
- changing Organization clears incompatible Installation/subject state;
- desktop persistent navigation, tablet collapse, mobile drawer/stacked local context;
- responsive transformation never changes authority/scope meaning.

Exact locked low-fi HTML remains current D6-R2 evidence under `qualification/d6-r2-wireframes/`.

## 6. Frontend Product Experience realization

D6-R2 follows the local Frontend Product Experience Planning Method:

```text
P0–P4 recover/revalidate needs + IA
P5 complete screen/material-surface inventory
P6/P7 only when material ambiguity triggers them
P8 executable low-fi structural HTML by block; operator-only LOCK
P9 exact frontend ↔ Product/backend contract
P10+ only as method requires
```

A LOCK is operator authority over accepted human structure, not proof that the underlying Product wire can never be falsified later.

When P8/P9 evidence exposes an upstream problem:

```text
screen need / human job
→ prove material backend insufficiency
→ reopen smallest Product/D4/D5 owner
→ Global Maximum
→ repair/prove
→ bounded revalidation of affected frontend region only
```

Do not restart unrelated locked blocks or D0–D8 by default.

## 7. Current accepted D6-R2 evidence

Current operator-LOCK evidence includes the accepted shell/overview and the later bounded blocks recorded in the P8 registry, including:

- B00 App Shell / global IA;
- B01 Overview;
- B00-R2 Notifications utility;
- B11 Personal Inbox;
- B12 Notification Routing Settings;
- B110 Approvals;
- B10 Preparação main structure.

The exact current lock/evidence registry is `D6-R2-P8-BLOCK-LEDGER.md`; mutable overall stage/next-action status belongs only to `docs/roadmap.md`.

### 7.1 B10 Preparação

Accepted human structure remains:

```text
exact Organization + Marketplace Installation
→ bounded multi-source product search
→ selected source-qualified Product
→ marketplace requirements + source evidence
→ correspondence correction when needed
→ continue to downstream ListingIntent authoring
```

B10 does not create ListingIntent merely as a navigation side effect and does not publish directly to provider.

A later B20 planning falsifier found the current read wire insufficient for some human recognition/selection jobs. B10 main structure remains protected, while its correspondence region/P9 awaits bounded revalidation after the paused human-operable-read-projection prerequisite (#70).

### 7.2 Personal Notifications

Accepted UX:

- compact shell utility can expose personal attention without becoming a global unread-count business authority;
- Inbox is exact-self awareness, with read/archive personal state only;
- Notification source continuation always reauthorizes current source access;
- source truth/action remains at source owner;
- Organization routing settings select only bounded eligible Principal candidates;
- no generic subscriptions/preferences/delivery channels.

### 7.3 Approvals / AuthorizationRequest

Approvals presents exact-human actionable AuthorizationRequests plus legitimate Governance history/context without collapsing:

```text
ordinary Permission
pending request
AuthorizationDecision
source target/action
execution outcome
```

Decision UI acts on the exact Request/revision; Notification merely navigates/alerts and never authorizes or decides.

## 8. Frontend package/dependency direction

Accepted target dependency direction remains:

```text
app/routes
    ↓
features/<flow-or-lens>
    ↓
api/<owner-family>
    ↓
api/transport
    ↓
api/generated

features ─→ ui
app/routes ─→ ui
```

Feature grouping follows human jobs/owner boundaries, not a generated endpoint-per-file framework. Generated OAD types stay at transport boundary; frontend state/UX models do not become a new business schema authority.

## 9. YAGNI exclusions

Without new evidence D6 does not introduce:

- second global server-state store;
- generic workflow/action/command frontend layer;
- BFF/screen API;
- offline-first sync;
- websocket/event-stream architecture by fashion;
- microfrontends/plugins;
- universal design-system platform work;
- SSR/meta-framework without a concrete need;
- generic generated Query-hook/business-store layer;
- client-owned metrics/recommendation/opportunity-score business semantics;
- frontend dictionaries that invent dynamic provider/source Product meaning.

## 10. Reopen triggers

Reopen only the implicated D6 decision when real operator evidence, current Product semantics or implementation proof shows the current interaction/state/topology cannot satisfy the human job safely. A material frontend need may reopen upstream Product authority; visual preference, framework fashion or historical document shape alone may not.
