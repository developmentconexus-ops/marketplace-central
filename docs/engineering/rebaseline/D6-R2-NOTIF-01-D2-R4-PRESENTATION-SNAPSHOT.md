# NOTIF-01 D2-R4 — Notification Presentation Snapshot

> **Status:** OPEN — TARGETED D2 REOPEN / OPERATOR ADJUDICATION REQUIRED
> **Trigger:** [D5-F2 Inbox Presentation Context Finding](D6-R2-NOTIF-01-D5-F2-INBOX-PRESENTATION-CONTEXT-FINDING.md)
> **Accepted base:** D2-R + D2-R2 + D2-R3 remain ACCEPTED / OPERATOR-RATIFIED; all unaffected clauses remain unchanged
> **Scope law:** add only the smallest immutable human-presentation state required for a useful personal Inbox without source-per-row rereads or a generic payload
> **Product wire:** unchanged — 99 Product operations · 30 ordinary Permissions · Principal kinds H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Why D2 reopens narrowly

Canonical Notification identity/state is already sufficient for correctness, routing, replay, read/archive and bounded supersession. The accepted frontend consumer now proves one presentation gap: opaque source identities alone cannot make several same-kind Inbox rows human-distinguishable without N source rereads.

The Notification Architecture Design already permits a bounded title/summary for Inbox usability, while accepted D2-R explicitly requires a targeted D2 reopen rather than hiding source data in free-form payload/template state.

This repair adds one presentation atom only.

---

## 2. Canonical Notification state gains one immutable field — CANDIDATE

Canonical Notification state adds:

```text
subject_display_label
```

The revised relevant state is conceptually:

```text
notification_id
organization_id
recipient_principal_id
kind
source_ref
source_occurrence_key
source_occurred_at
source_committed_at
subject_display_label       // D2-R4
created_at
read_at?
archived_at?
superseded_at?
superseded_by_notification_id?
supersession_reason?
work_replacement_basis?
sale_attention_replacement_basis?
revision
```

Exact physical length/encoding is later wire/storage work. Semantically the label is non-empty and bounded.

---

## 3. Meaning and ownership — CANDIDATE

`subject_display_label` is:

> a small immutable, human-readable label for the exact source occurrence/subject, retained only so the personal Inbox can distinguish neighboring awareness items without rereading every source.

Ownership laws:

- the **source owner** owns the business-derived meaning used to produce the label for its occurrence;
- **Personal Notifications** owns the retained immutable presentation snapshot as part of the Notification after materialization;
- the label is not source business truth, current source state, identity or authorization;
- Personal Notifications may not invent or reinterpret source business meaning merely to produce a nicer label;
- exact D3 transport/hand-off of the source-owned presentation atom is not selected by D2-R4 and must be revalidated after acceptance.

Historical law:

```text
source label/title later changes
→ existing Notification subject_display_label does not change
```

A later source reread may show different current presentation and remains authoritative for the source.

---

## 4. Presentation law — CANDIDATE

The final Notification title/sentence is **not** persisted by this repair.

Baseline presentation model:

```text
NotificationKind
→ client-localized Product title/copy

subject_display_label
→ immutable human context line
```

Examples only:

```text
kind = NEW_MARKETPLACE_SALE
subject_display_label = "Mercado Livre · Pedido #200391"

kind = FULFILLMENT_ATTENTION
subject_display_label = "Venda #200388"

kind = WORK_ASSIGNMENT
subject_display_label = "Revisar materialização da venda #200340"

kind = AUTHORIZATION_ACTION_REQUIRED
subject_display_label = "Alteração de preço · Anúncio MLB-..."
```

The examples do not freeze localization/copy.

No second `summary`, `detail`, `subtitle`, template or arbitrary payload field is admitted now.

---

## 5. Safety / data-minimization law — CANDIDATE

Because a Notification may remain after current source access is revoked, `subject_display_label` must be safe to retain as personal-awareness history.

It must not contain by default:

```text
buyer/customer name
address
phone/e-mail
payment detail
fiscal payload/document detail
full product/customer free text
credential/secret/token
provider raw body
arbitrary source JSON
```

A source owner must choose a notification-safe label using the smallest non-sensitive identifiers/display context required by the human job.

If a source cannot produce a safe specific label, it must fall back to a bounded generic safe label under source semantics rather than copying sensitive current detail.

---

## 6. Non-authority / non-machine-use law — CANDIDATE

`subject_display_label` must never be used to:

- identify or join a source object;
- build a deep link;
- select a recipient;
- decide routing or Permission;
- deduplicate/replay an occurrence;
- correlate Work/Post-Sale replacement;
- infer source state/currentness;
- authorize source access;
- mutate or acknowledge a source;
- drive machine automation.

All such behavior continues to use canonical typed IDs/refs, occurrence keys and owner authority.

Equal labels never imply equal source identity.

---

## 7. Human Inbox consequence — CANDIDATE

A `ListMyNotifications` representation may later expose:

```text
notification_id
kind
subject_display_label
source_ref
source_occurred_at
created_at
read/archive state
```

proportionately, allowing a bounded bell preview/full Inbox to remain comprehensible without one current source fetch per row.

Opening/navigating the Notification still performs a current authorized source read. The retained label never guarantees that the source remains accessible or that its current presentation/state matches the historical label.

---

## 8. Negative controls — CANDIDATE

This repair fails if it introduces:

1. free-form JSON payload/metadata;
2. generic template engine or template variables;
3. stored final localized notification sentence as business authority;
4. mutable presentation that tracks current source state;
5. sensitive source detail retained merely for convenience;
6. source identity/navigation/dedup based on the label;
7. a second summary/detail field without a proved consumer;
8. frontend-local aliases becoming apparent Product truth;
9. D3 event spelling, D5 schema/path, D7 storage or D6 component mechanics inside D2.

---

## 9. Feed-forward obligation — CANDIDATE

If D2-R4 is accepted, the already-accepted D3-R communication contract must be revalidated only for this new immutable presentation atom:

> source-owner committed-fact communication must convey enough source-owned presentation meaning to materialize the accepted `subject_display_label` without making a source reread-per-Inbox-row the correctness baseline.

If that cannot be done without changing communication ownership or introducing a generic payload, stop and reopen only the implicated D3 clause.

---

## 10. Gate

```text
D2-R / D2-R2 / D2-R3              ACCEPTED / OPERATOR-RATIFIED
D3-R                               ACCEPTED
D5-F2 Inbox presentation finding   PROVED CONSUMER
D2-R4 presentation snapshot        READY FOR OPERATOR REVIEW
D3 presentation feed-forward       BLOCKED BY D2-R4
D5-R3 operation table              BLOCKED BY D2-R4/D3 feed-forward
canonical Product OAD              UNCHANGED
D6 / D7 / D8                       BLOCKED for NOTIF-01
Product implementation             BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates only this bounded D2-R4 presentation-snapshot repair. Do not alter D3/OAD/frontend/runtime first.