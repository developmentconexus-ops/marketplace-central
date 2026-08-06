# ADR-013: Webhook is a pointer, never data — always-200, dedupe amended to `_id`-only

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed MIS-007's webhook ingest from its first
citation and was enforced in review, but no document was ever written under this number.
It is reconstructed here from the 22 live citations of `ADR-11`/`ADR-011` that assert
this meaning, harvested at
`docs/architecture/decisions/_citations/adr-011-citations.md` (Assertion A1). The same
`ADR-11` label was also used, unrelated, for two other decisions in two other missions —
see the registry at `docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md`.
Every clause below is traceable to the mission record or the interface contract that
already asserts it.

## Context

Mercado Livre calls `POST /webhooks/{provider}` to notify MIS-007's system of order
events. That endpoint is the first public, unauthenticated write surface in the system:
anyone who finds the URL can send it a body. Two temptations follow from that shape.
The first is to trust the body — order id, status, amounts — and write it straight into
domain tables, which lets a forged or stale notification inject data nobody fetched from
the authenticated API. The second is to reject anything that looks wrong with a 4xx,
which trains Mercado Livre's retry logic to storm the endpoint (documented risk: 8
retries within an hour) the moment a topic is unrecognized or a body is malformed.

The system also needs to avoid processing the same notification twice. The natural key
looked like the full identifying tuple Mercado Livre sends — `(provider, topic,
resource, notification_id)` — but that tuple, matched with `COALESCE` against nulls,
turned out to block a legitimate second notification about the same resource forever.
The dedupe rule was amended mid-mission (P5 round 2, negotiation N-4) to narrow the key,
and the amendment is part of what this ADR records.

## Decision

**The webhook body is never treated as domain data. It is a pointer used to schedule an
authenticated refetch, and the endpoint always answers `200` regardless of content.**

### The clauses

**§1 — The body is a pointer, not a fact.** Only `resource` and `topic` are read from
the notification, and only to schedule an authenticated refetch (ADR-04 lineage). No
field from the webhook body is written into a domain table.
> `.mnfs/MIS-007-ml-sync/mission.md:224` — "ADR-11 Webhook: ponteiro, nunca dado;
> always-200." ... "Só `resource`+`topic` usados, só como ponteiro p/ refetch AUTENTICADO
> via ADR-04."

**§2 — `resource` is pinned before use.** `resource` is attacker-controlled input. It is
only usable as a pointer when it matches `^/orders/[0-9]+$`; outside that pin the row is
recorded with terminal status `malformed` and is never concatenated into an authenticated
call URL.
> `.mnfs/MIS-007-ml-sync/research/webhook-inbox-interface-contract.md:53-57` — "`resource`
> só é usável como ponteiro quando casa `^/orders/[0-9]+$` (ancorado); fora do pin →
> status terminal `malformed` ... O `resource` é attacker-controlled — NUNCA concatenado
> em URL de chamada ML autenticada sem casar o pin."

**§3 — The endpoint always returns 200.** Every processable case — including an unknown
topic, a malformed body, or an unmapped `user_id` — receives an empty `200`. The only
exception is a failure to insert into the inbox itself, which is accepted as rare and
covered by the 5-minute reconciliation cycle.
> `.mnfs/MIS-007-ml-sync/research/webhook-inbox-interface-contract.md:42` — "**200 vazio
> SEMPRE** — inclusive topic desconhecido, body malformado, user_id desconhecido" ... "nunca
> 4xx/5xx p/ conteúdo (evita storm 8×/1h)".
> `.mnfs/MIS-007-ml-sync/M-08-webhook-ingest/validation-contract.md:54` — "webhook NUNCA
> vira dado de domínio (ADR-11)".

**§4 — The raw body never reaches domain tables.** The received payload is persisted
verbatim (capped at 64KB, truncated with a marker) for audit only, and is never parsed
into a domain fact.
> `.mnfs/MIS-007-ml-sync/research/webhook-inbox-interface-contract.md:61-62` — "`raw_body
> text` NOT NULL (cap 64KB, truncado com marcador — auditoria; NUNCA parseado para dado de
> domínio)."

**§5 — The source IP is informative only, never a gate.** `source_ip` and
`ip_official` (matched against Mercado Livre's published allowlist) are recorded for log
and audit purposes. They never cause a request to be rejected — this was an explicit
operator decision (P1), because the header used to derive the IP is forgeable outside
the tunnel.
> `.mnfs/MIS-007-ml-sync/mission.md:224` (Non-Functional Scope row) — "`ip_official` é
> INFORMATIVO/log-only, NUNCA gate de aceitação".

**§6 — Dedupe key was amended: narrowed to `_id`-only.** The original ratified dedupe key
was the full tuple `(provider, topic, resource, notification_id)`. It was amended (P5
round 2, negotiation N-4) to `UNIQUE (provider, notification_id) WHERE notification_id IS
NOT NULL`, because the full tuple under `COALESCE` would block a legitimate
re-notification about the same resource indefinitely. **The amended, narrowed rule is the
rule.** Citations after the amendment increasingly drop the "emendado" qualifier and cite
"ADR-11" bare — this document names the amendment explicitly so that drift does not
reintroduce the original, broader tuple as the assumed key.
> `.mnfs/MIS-007-ml-sync/research/webhook-inbox-interface-contract.md:101-105` — "Dedupe:
> `UNIQUE (provider, notification_id) WHERE notification_id IS NOT NULL`; duplicata →
> upsert `attempts_provider`/`received_at` (não cria row). Sem `notification_id`: sem
> dedupe de transport ... (Chave estreitada da tupla original de ADR-11 — emenda registrada
> em `mission.md` ADR-11, P5 r02 N-4: tupla cheia com COALESCE bloquearia re-notificação
> legítima do mesmo resource p/ sempre.)"
> `.mnfs/MIS-007-ml-sync/mission.md:373` — "reprocessar notificação → zero duplicata de
> DOMÍNIO (IngestOrder idempotente; dedupe de inbox só com `_id` — ADR-11 emendado; ver
> M-08 Done Means)"

**§7 — Without a `notification_id`, transport dedupe does not apply; domain idempotence
must carry the guarantee.** A notification without `_id` produces extra, harmless inbox
rows on replay. Zero-duplicate is proven at the domain layer (`IngestOrder` idempotence),
never at the transport layer, for that path.
> `.mnfs/MIS-007-ml-sync/M-08-webhook-ingest/milestone.md:77-80` — "Replay com `_id`
> presente → zero row nova (upsert attempts_provider). Sem `_id` NÃO há dedupe de
> transporte (IC-04/ADR-11 emendado — P5 r02 N-4): replay vira rows extras inofensivas; a
> prova de zero-duplicata é de DOMÍNIO — IngestOrder idempotente, zero efeito duplicado."

## Rationale

Treating the webhook as a pointer rather than as data means a forged or malformed
notification can, at worst, trigger a wasted authenticated refetch — it can never inject
a fact into the system that nobody fetched from the source of truth. Always answering
`200` removes the only lever an attacker or a flaky network has to make Mercado Livre
retry-storm the endpoint; the cost of that leniency is paid entirely inside the inbox
table (a row with a terminal status), never by the caller or by domain data.

The dedupe amendment exists because correctness of the *narrow* case turned out to
conflict with correctness of the *common* case: a full-tuple key with `COALESCE`
semantics permanently blocked legitimate re-notifications about the same order, which is
a routine and expected event (status changes, retries) rather than a duplicate. Narrowing
to `_id`-only trades away transport-level dedupe on the no-`_id` path in exchange for not
silently dropping real notifications — but that trade is only sound because a second
mechanism, domain-level idempotence in `IngestOrder`, is required to carry the
zero-duplicate guarantee on exactly the path transport dedupe no longer covers.

## Consequences

- No webhook body field is ever a source of domain truth; every domain write that
  originates from a notification goes through an authenticated refetch first.
- The endpoint's error surface is deliberately flat: external callers never see 4xx/5xx
  for content reasons. Diagnosis of malformed/unmapped/ignored traffic happens by reading
  `notifications_inbox` status, not by inspecting HTTP responses.
- `IngestOrder` must be idempotent independent of transport dedupe — it cannot assume the
  inbox has already deduplicated a given resource.
- A notification without `_id` will always produce a new inbox row on replay. This is
  accepted as harmless noise in the inbox, not as a domain risk.

## Open questions

- **No-`_id` flood/retention gap.** The amendment that narrowed the dedupe key to
  `_id`-only means a notification stream without `_id` bypasses transport dedupe entirely
  and grows `notifications_inbox` without bound. No retention or purge policy for this
  table has been declared anywhere in the interface contract or the milestone record.
  Inbox rows themselves are inert once written (proven by `M08-C3`), so the risk is
  storage exhaustion on a public, unauthenticated surface, not domain corruption. This gap
  was flagged during planning and never closed.
  > `.mnfs/MIS-007-ml-sync/planning-reviews/p7-seat5-doublepass-r02.md:55` — "Flood without
  > `_id` bypasses dedupe (ADR-11 emendado) and grows `notifications_inbox` unboundedly; no
  > retention/purge policy is declared anywhere in IC-04/M-08. Rows are inert (M08-C3), but
  > storage exhaustion on the public surface is unaddressed."

## Alternatives Considered

**Verify a cryptographic signature or token on the webhook before trusting it.**
Rejected by explicit operator decision (P1): the webhook is treated as an untrusted hint
regardless of any signature, and the real defense is that no webhook field ever becomes
domain data without an authenticated refetch. Signature verification would add complexity
without changing what the system is willing to trust.

**Reject unrecognized topics or malformed bodies with 4xx.** Rejected: Mercado Livre's
retry behavior turns any 4xx into a retry storm (documented at up to 8 retries within an
hour). Recording the row as terminal and answering `200` costs nothing at the domain
layer and avoids the storm entirely.

**Keep the original full-tuple dedupe key.** Rejected after measurement: the full tuple
under `COALESCE`-based matching blocked legitimate re-notifications about the same
resource permanently, which is a more frequent and more damaging failure than the
narrowed key's no-`_id` flood risk.
