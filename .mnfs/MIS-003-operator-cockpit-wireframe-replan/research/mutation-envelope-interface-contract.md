# Interface Contract

```yaml
id: IC-03
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

Single durable provider-write path ("fila de sync com preview e protocolo"): `mutations` capability inside a new `mutations` module (or `listings`-adjacent module — final module home decided in ADR-13, not per feature) → HTTP transport → SDK → all write UIs (Anúncios bulk bar, Estoque corrigir, Vínculos aplicar, corrigir-atributo flow).

## Why This Contract Exists

No queue/outbox exists; one implemented StockWriter is synchronous and never live-run. The wireframe shows FIVE write surfaces all promising "preview e protocolo". Without one envelope, workers build N ad-hoc write paths (the redundancy the operator forbade). This contract also implements the seven provider-write gates from `contracts/governance/execution-lanes.json`: actor, idempotency, execute, resolved-link, policy, source-timestamp, before-after-audit.

## Resources Or Entities

**MutationProtocol** (a.k.a. protocolo):

| Field | Type | Null | Notes |
| --- | --- | --- | --- |
| `protocol_id` | string | no | `MP-` + zero-padded sequence, e.g. `MP-000042`; user-visible |
| `installation_id` | string | no | |
| `type` | enum | no | see Enums |
| `state` | enum | no | see lifecycle |
| `actor` | string | no | actor label, e.g. `operator_supplied_unverified` (ADR-009) |
| `intent` | object | no | type-specific payload (see per-type intents) |
| `selection` | object | no | `{mode: "explicit"\|"filter", listing_ids?: string[], filter?: <IC-02 filter object>, q?: string}`; filter selections are SNAPSHOT to explicit item rows at preview time |
| `totals` | object | no | `{items, previewed, applied, failed, skipped}` |
| `source_as_of` | string | no | RFC3339; freshest source-fact time used in preview (source-timestamp gate) |
| `retried_from` | string | yes | protocol_id of the terminal protocol this one was cloned from; set only by `retryMutationFailures`, null otherwise |
| `created_at/previewed_at/approved_at/finished_at` | string | yes | RFC3339 UTC |

**MutationItem** (one row per target; the audit record):

| Field | Type | Null | Notes |
| --- | --- | --- | --- |
| `item_id` | string | no | `protocol_id` + seq |
| `listing_id` | string | no | IC-02 canonical id (or product/link target id for `link_apply`) |
| `idempotency_key` | string | no | `{protocol_id}:{listing_id}` — provider adapter must send it |
| `before` | object | yes | captured at preview; e.g. `{price: {"amount":"54.90"}}`; null only for `link_apply` create-link |
| `after` | object | no | intended value |
| `state` | enum | no | `previewed \| approved \| applying \| applied \| failed \| skipped` |
| `failure` | object | yes | `{code, message_pt, message_provider, retryable}` |
| `applied_at` | string | yes | |

### Lifecycle (protocol)

`draft → previewed → approved → applying → applied | partially_failed | failed_preserved | cancelled`

- `draft`: created with selection + intent. No provider contact.
- `previewed`: items materialized (filter snapshot → explicit rows), `before` captured from read models, per-item validation run (resolved-link gate, policy gate). Preview is idempotent-recomputable while in `draft/previewed`.
- `approved`: explicit operator approval carrying `execute: true` (execute gate). Terminal for editing.
- `applying`: in-process poller picks up approved protocols, applies items in chunks (chunk size 20, sequential per installation) through provider capability adapters; rate-limit aware.
- Terminal: `applied` (all ok), `partially_failed` (mixed; failed items preserved), `failed_preserved` (all failed; nothing hidden, no auto-retry), `cancelled` (only from `draft/previewed`).
- Restart safety: poller resumes `approved`/`applying` protocols after process restart; per-item idempotency keys make re-application safe.

### Per-type intents

| `type` | `intent` | Write target | This mission |
| --- | --- | --- | --- |
| `price_update` | `{new_price: {amount, currency}}` or `{price_expr: null}` — absolute values only, no formulas | ML item price via new PriceWriter capability | implemented |
| `stock_correct` | `{publish_quantity: int}` (deveria-publicar value) | ML available_quantity via existing StockWriter | implemented (folds `inventory_stock_actions` semantics; that service delegates here) |
| `link_apply` | `{action: "approve_candidate"\|"manual_resolve"\|"reject_listing", candidate_id?, product_id?}` | product_links resolutions (internal write) | implemented |
| `listing_pause` | `{}` | ML item status | implemented (bulk bar) |
| `listing_resync` | `{}` | re-fetch item + upsert read model | implemented |
| `listing_edit` | `{attributes: [{id, value_name}]}` | ML item attributes | implemented (corrigir-atributo, M-06) |
| `listing_create` | `{product_id, listing_type_code, price, category_id, attributes[]}` | ML item create | **contract-only this mission** — enum value + intent shape reserved; any runtime call → 422 `type_not_enabled` |

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| `createMutation` `POST /mutations` | write action in any UI | `{installation_id, type, selection, intent, actor}` | `201 MutationProtocol` (state `draft`) | validates intent shape only; `actor` is client-supplied in the POST body and recorded verbatim with provenance label `operator_supplied_unverified` (ADR-009) — never derived from auth context this mission; missing/empty → 400 `actor_required` |
| `previewMutation` `POST /mutations/{id}/preview` | preview step | — | `MutationProtocol` (state `previewed`) + first item page | snapshots selection; caps: max 2000 items → 422 `selection_too_large` |
| `approveMutation` `POST /mutations/{id}/approve` | operator confirm | `{execute: true}` | `MutationProtocol` (state `approved`) | literal `execute:true` required (execute gate); 409 if preview stale (>15 min) → re-preview |
| `cancelMutation` `POST /mutations/{id}/cancel` | operator cancel | — | state `cancelled` | only from draft/previewed |
| `retryMutationFailures` `POST /mutations/{id}/retry` | operator retry | — | new protocol cloned from failed retryable items | never mutates the terminal protocol (audit immutability) |
| `getMutation` `GET /mutations/{id}` | protocolo detail / polling | — | `MutationProtocol` | poll while `applying` |
| `listMutationItems` `GET /mutations/{id}/items` | detail table | `cursor, limit, state?` | item page | sort: seq ASC |
| `listMutations` `GET /mutations` | central de sync | `installation_id, cursor, limit, state?, type?` | page | sort: `created_at DESC` (newest-first) |

## Enums And Statuses

- Protocol `state`: `draft, previewed, approved, applying, applied, partially_failed, failed_preserved, cancelled`.
- Item `state`: `previewed, approved, applying, applied, failed, skipped`.
- `failure.code` (stable taxonomy; UI binds to codes, never message strings — review finding 9d):
  `provider_validation` (retryable=false) · `provider_rate_limited` (true) · `provider_unavailable` (true) · `provider_auth` (false; installation reauth needed) · `listing_paused_remote` (false) · `link_unresolved` (false) · `policy_missing` (false; policy gate — required pricing/stock policy absent, never defaulted) · `sku_invariant_violation` (false; ADR-16 — SELLER_SKU ≠ linked CODPROD, rejected at preview) · `stale_source` (true; re-preview) · `conflict_remote_changed` (false; before-value no longer matches remote) · `type_not_enabled` (false) · `internal` (false; unknown provider errors are never auto-retried — operator re-runs via retry clone after inspection).

## Error Cases

See Error Matrix.

## Persistence Expectations

Tables `mutation_protocols` + `mutation_items` (tenant-scoped). Items are the durable audit: never deleted, never rewritten after terminal state (retry clones). The applying poller is in-process (goroutine + DB claim with `FOR UPDATE SKIP LOCKED`); no external queue infrastructure this mission (ADR-13). `inventory_stock_actions` remains readable for history; new stock corrections route through the envelope.

## Canonical Examples

Success — preview response (truncated):

```json
{
  "protocol_id": "MP-000042",
  "installation_id": "inst_1",
  "type": "price_update",
  "state": "previewed",
  "actor": "operator_supplied_unverified",
  "intent": {"new_price": {"amount": "49.90", "currency": "BRL"}},
  "selection": {"mode": "filter", "filter": {"exception": "below_margin"}},
  "totals": {"items": 6, "previewed": 6, "applied": 0, "failed": 0, "skipped": 0},
  "source_as_of": "2026-07-14T17:05:00Z",
  "created_at": "2026-07-14T17:06:01Z",
  "previewed_at": "2026-07-14T17:06:03Z"
}
```

Rejection — approve without execute flag:

```json
{"error": {"code": "execute_required", "message": "aprovação exige execute=true explícito"}}
```

Failed item (preserved):

```json
{
  "item_id": "MP-000042-003",
  "listing_id": "inst_1~MLB3456790~-",
  "idempotency_key": "MP-000042:inst_1~MLB3456790~-",
  "before": {"price": {"amount": "89.00", "currency": "BRL"}},
  "after": {"price": {"amount": "75.40", "currency": "BRL"}},
  "state": "failed",
  "failure": {"code": "provider_validation", "message_pt": "A marca é obrigatória nesta categoria.", "message_provider": "Attribute [BRAND] is required", "retryable": false}
}
```

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| missing `installation_id` | 400 | `installation_required` | createMutation / listMutations |
| unknown installation | 404 | `installation_not_found` | createMutation |
| missing/empty `actor` | 400 | `actor_required` | createMutation; actor gate — never defaulted server-side |
| preview cannot establish `source_as_of` (no source-fact time for any selected item) | 422 | `source_time_unavailable` | source-timestamp gate; unknown source time is never defaulted to now |
| unknown protocol | 404 | `protocol_not_found` | |
| invalid intent shape for type | 400 | `invalid_intent` | |
| selection resolves to 0 items | 422 | `empty_selection` | at preview |
| selection > 2000 items | 422 | `selection_too_large` | narrow the filter |
| approve without `execute:true` | 400 | `execute_required` | execute gate |
| approve with stale preview (>15 min) | 409 | `preview_stale` | re-preview required |
| state transition not allowed | 409 | `invalid_state` | e.g. cancel after approved |
| unresolved link on price/stock/edit item | — (item-level) | `link_unresolved` | item fails/skips at preview, protocol proceeds with the rest |
| required policy absent on price/stock item | — (item-level) | `policy_missing` | policy gate; item fails at preview, never a silent default |
| SELLER_SKU ≠ linked CODPROD on `listing_edit`/`listing_create` item | — (item-level) | `sku_invariant_violation` | ADR-16; rejected at preview |
| `link_apply` targeting an already-resolved link with same outcome | — (item-level) | — | item ends `skipped` (no-op, not a failure) |
| conflicting re-resolution of an already-resolved link; remote state changed since preview snapshot | — (item-level) | `conflict_remote_changed` | not retryable; operator re-previews |
| provider rejects item payload (business validation) | — (item-level) | `provider_validation` | not retryable; `message_provider` preserved verbatim |
| provider returns 429 / rate limit | — (item-level) | `provider_rate_limited` | retryable via retry clone |
| provider timeout / 5xx / unreachable | — (item-level) | `provider_unavailable` | retryable via retry clone |
| provider returns 401/403 (token invalid/expired) | — (item-level) | `provider_auth` | not retryable; operator reconnects installation first |
| listing paused/closed remotely since preview | — (item-level) | `listing_paused_remote` | not retryable |
| source data older than tolerance at apply time | — (item-level) | `stale_source` | retryable via retry clone (re-preview refreshes source) |
| unknown/unmapped provider error | — (item-level) | `internal` | not retryable; never auto-retried, operator inspects then retry-clones |
| `listing_create` runtime call | 422 | `type_not_enabled` | contract-only this mission |
| retry on protocol with no retryable failures | 422 | `nothing_to_retry` | |

## Database Shape

- `mutation_protocols`: PK `(tenant_id, protocol_id)`; columns per entity above; `state` check constraint; `selection JSONB`, `intent JSONB`, `totals JSONB`; `retried_from TEXT NULL` (set only by `retryMutationFailures`).
- `mutation_items`: PK `(tenant_id, protocol_id, seq)`; `listing_id`, `idempotency_key UNIQUE (tenant_id, idempotency_key)`, `before JSONB NULL`, `after JSONB`, `state` check, `failure JSONB NULL`.
- Timestamps TIMESTAMPTZ UTC.

## Seed Data

- Integration fixture: one protocol per lifecycle terminal state over the IC-02 seed listings, including one `partially_failed` with a `provider_validation` failed item. Ephemeral-postgres lane.

## Timestamp And ID Semantics

- `protocol_id` sequence is tenant-scoped, monotonic, user-visible (wireframe "protocolo #12-B" → `MP-000012` style).
- `source_as_of` must be ≤ 15 min old at approve time (staleness window = preview TTL).

## Compatibility Rules

- New write types extend the `type` enum + one intent schema; lifecycle, gates, and item shape never change per type.
- Real queue infrastructure may later replace the in-process poller behind the same tables + API (the table IS the contract).

## Route Namespace

- Server: `/mutations` prefix, mounted by the module that ADR-13 designates (M-03 owns).
- Client: preview/confirm modal shared component (M-03 owns); protocolo detail page `/protocolos/:protocolId` (IC-05).

## Transport And Integration

- Vite dev proxy: add `/mutations`.
- Polling: UI polls `GET /mutations/{id}` every 2s while `applying`; no websockets this mission.
- Provider-write lane evidence class `production-like` applies to any live execution (governance).

## Must Preserve

- Every provider write routes through this envelope — no direct capability calls from transport/UI.
- Item audit immutability after terminal state; retry = new protocol.
- All seven provider-write gates enforced structurally (schema + state machine), not by convention.
- Failed results preserved honestly (`failed_preserved`, M-06 lesson) — never silently retried into success.
- SKU invariant (ADR-16): `listing_edit`/`listing_create` items require `link.state=resolved`; enforced at preview.

## Must Not Decide In Feature Execution

- Lifecycle states, failure codes, idempotency-key format, selection semantics, caps, polling cadence.

## Validation Impact

M-03 validation contract cites: lifecycle proof per state, restart-resume proof, idempotent re-apply proof, per-gate negative proofs (one criterion per gate), failure-taxonomy rendering proof in UI.
