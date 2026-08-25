# D6-R2 P8 — Current Structural Wireframe LOCK Registry

> **Status:** CURRENT ACCEPTANCE EVIDENCE  
> **Mutable program/block sequencing:** `docs/roadmap.md`  
> **Method:** Frontend Product Experience Planning Method v2.3  
> **Product implementation:** blocked until D9

## 1. Registry role

This file records stable operator `LOCK` dispositions for D6-R2 structural low-fidelity HTML. It does **not** own current next action, stage status or upstream Product sufficiency; the roadmap and block-specific P9/current owners do.

P8 HTML decides structural hierarchy/placement/density/navigation/state placement/responsive transformation/interaction proof, not final visual design. Reviewed HTML may still self-label `candidate`; operator `LOCK` lives in durable ratification/this registry.

## 2. Current locked blocks

| Block | Operator disposition | Current evidence |
| --- | --- | --- |
| **B00 App Shell + global IA** | **LOCKED** 2026-08-22 | `qualification/d6-r2-wireframes/b00-app-shell.html` |
| **B01 Overview** | **LOCKED** 2026-08-22 | `qualification/d6-r2-wireframes/b01-overview.html` |
| **B00-R2 Notification utility** | **LOCKED** 2026-08-23 | `qualification/d6-r2-wireframes/b00-r2-notifications.html` + Notifications P8 ratification |
| **B11 Personal Inbox** | **LOCKED** 2026-08-23 | `qualification/d6-r2-wireframes/b11-notifications-inbox.html` + Notifications P8 ratification |
| **B12 Notification Routing Settings** | **LOCKED** 2026-08-23 | `qualification/d6-r2-wireframes/b12-notification-routing-settings.html` + Notifications P8 ratification |
| **B110 Approvals** | **LOCKED** 2026-08-24 | `qualification/d6-r2-wireframes/b110-approvals.html` + B110 P8 ratification |
| **B10 Preparação** | **LOCKED main structure / correspondence CANDIDATE** | `qualification/d6-r2-wireframes/b10-preparation.html` + B10 P8 ratification + bounded correspondence revalidation |

No other D6-R2 block inherits a LOCK by proximity. Only the operator can add/change a LOCK.

## 3. Global shell / IA laws

B00 current locked baseline:

```text
Organization = only global workspace
page-local exact Marketplace Installation host when required
no hidden/default Installation
desktop persistent navigation
tablet collapsible navigation
mobile drawer + stacked local context
responsive transformation never changes authority meaning
```

Current primary IA:

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
  Configurações
```

Notifications is a shell utility/Inbox plus Settings sub-surface, not a new global sidebar business mass.

## 4. B01 Overview

Locked contextual priority:

```text
Work known + actionable  → attention expands/leads
Work known-empty         → attention collapses; never proves overall health
Work unknown/unavailable → uncertainty remains visible
all cases                → marketplace/account + Performance + Economics orientation remains visible
```

## 5. Personal Notifications locks

### B00-R2 utility

```text
Organization-scoped bell
→ unread known-present / known-empty / unavailable
→ bounded recent preview
→ source continuation separate from awareness mutation
→ full R128 Inbox continuation
```

No unread numeric total authority and no `notifications.manage` requirement for personal awareness.

### B11 Inbox

```text
structured vertical awareness list
Ativas | Arquivadas
Todas | Não lidas | Lidas
closed NotificationKind filter
cursor continuation
source continuation
read/unread + archive/restore awareness mutation
```

No text search, total count, mark-all-read, bulk archive, severity/priority or saved views.

### B12 routing settings

```text
Configurações > Notificações
exactly 10 ORG_ROUTED rows
inline affected-row editor
candidate = principal_id + display_name only
cursor candidate continuation
explicit CONFIGURED(one+) / UNCONFIGURED
SetNotificationRoute with If-Match
```

No implicit `access.read`, role/Permission disclosure, configured-empty state, generic subscription DSL or preventive search API.

## 6. B110 Approvals

Locked human structure:

```text
Aprovações
├─ Para decidir   governance.decide
└─ Histórico      governance.read

Para decidir
→ actionable AuthorizationRequest list
→ request-detail route
→ one typed review-basis family
→ inline approve/reject confirmation with evidence visible
→ current Request revision + semantic-attempt idempotency

Histórico
→ immutable AuthorizationDecision list/detail
→ no decision controls
```

Carrier-specific current Product contract is owned by P9/W1/OAD and currently uses `body.etag + outcome` + `Idempotency-Key`, not historical `If-Match`. This transport correction did not reopen the locked human structure.

## 7. B10 Preparação

Accepted main structure:

```text
exact Organization + Marketplace Installation
→ bounded source Product search
→ selected source-qualified subject
→ marketplace requirements + source values/evidence
→ correspondence region
→ downstream ListingIntent authoring navigation
```

PR #70 resolved the upstream human-operability prerequisite without revoking the full B10 LOCK. The correspondence region now has a bounded functional candidate that renders real key+label selection and known-empty/unknown/unavailable population states. It remains **CANDIDATE** until operator walkthrough; P9 rerun follows only after explicit re-LOCK. Main hierarchy/table/handoff remains protected unless new evidence proves otherwise.

## 8. LOCK law

A later backend/Product correction may trigger one of:

```text
UNAFFECTED
bounded region revalidation
full block reopen
```

Choose the smallest proven impact. Do not rewrite accepted HTML merely to update an invisible transport carrier or historical status label.

Current stage progression belongs only to `docs/roadmap.md`.
