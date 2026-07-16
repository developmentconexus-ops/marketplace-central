# M-03 P2 batch plan — mutation-envelope-writes

Plan only; no files are written. Execution order is strictly:

`F-01 → F-02 → F-03 → F-04`

Every slice targets approximately 300 changed lines or fewer, starts with a failing test, receives independent review before a dependent slice begins, and uses no new dependency.

## Plan status and stop conditions

The backend plan is ready subject to one contract correction. F-04 remains gated on CHIP-M02.

| ID | Planning defect | Required resolution |
| --- | --- | --- |
| PD-01 | IC-03 declares `source_as_of` non-null, but `POST /mutations` persists a `draft` before preview can establish source time. Defaulting it violates ADR-17 and `source_time_unavailable`. | Before F01-S1, contract owner must amend IC-03/OpenAPI so `source_as_of` is nullable in `draft` and required non-null from `previewed` onward. No implementation may invent `now()` or an empty string. |
| PD-02 | CHIP-M02 has not merged, so the concrete Anúncios package/file carrying selection state does not exist in this checkout. Its feature brief does not fix that filename. | Before F-04, hub confirms M-02 F-03 merged, rebases CHIP-M03, and identifies the actual Anúncios page/export seam. The F-04 write-set below assumes the expected `packages/feature-listings/src/AnunciosPage.tsx`; a different result requires a delta-plan, not an improvised stand-in. |

---

# F-01 — protocolo-core

## Slice cards

### F01-S1 — IC-03 persistence shape

- Goal: introduce the two tenant-scoped envelope tables, checks, claim indexes, and nullable fact columns. Use migration `0038` only.
- Files:
  - `apps/server_core/migrations/0038_mutation_protocols.sql`
  - `apps/server_core/migrations/mutations_test.go`
- Failing test first: migration-shape test asserts the exact IC-03 columns, enums, PKs, unique `(tenant_id,idempotency_key)`, nullable fact/audit fields, and claim index.
- Done:
  - `mutation_protocols` PK `(tenant_id,protocol_id)`.
  - `mutation_items` PK `(tenant_id,protocol_id,seq)`.
  - Every table/index begins with `tenant_id`.
  - No fact column receives a zero/empty/default.
  - `source_as_of` follows the PD-01 resolution.
  - No FK or provider payload column is added speculatively.
- Complexity: `complex`
- Serves: M03-C01, C02, C03, C12
- Budget: ~220 lines.

### F01-S2 — lifecycle and terminal-state domain

- Goal: implement IC-03 enums, failures, pure table-driven lifecycle transitions, and terminal-state computation.
- Files:
  - `apps/server_core/internal/modules/mutations/domain/protocol.go`
  - `apps/server_core/internal/modules/mutations/domain/protocol_test.go`
  - `apps/server_core/internal/modules/mutations/domain/failure.go`
- Failing test first: exhaustive matrix of every protocol-state pair, including draft→approved and applying→cancelled rejection.
- Done:
  - States and failure codes are byte-exact IC-03 values.
  - Only `draft|previewed → cancelled`.
  - Item outcome matrices produce `applied`, `partially_failed`, or `failed_preserved`.
  - Terminal protocols reject all further mutation.
- Complexity: `complex`
- Serves: M03-C01, C03
- Budget: ~280 lines.

### F01-S3 — protocol IDs and repository creation/read paths

- Goal: create tenant-isolated protocols and allocate monotonic `MP-%06d` IDs without a global sequence.
- Files:
  - `apps/server_core/internal/modules/mutations/ports/repository.go`
  - `apps/server_core/internal/modules/mutations/adapters/postgres/repository.go`
  - `apps/server_core/internal/modules/mutations/adapters/postgres/repository_integration_test.go`
- Failing test first: concurrent creation under two tenants proves monotonic per-tenant IDs and permits both tenants to receive `MP-000001`.
- Done:
  - Allocation uses a transaction-scoped tenant lock plus the tenant’s current maximum; no third sequence table.
  - Every statement includes `tenant_id`.
  - Cross-tenant get/update returns not found.
  - Draft creation is atomic with actor/type/intent/selection.
- Complexity: `complex`
- Serves: M03-C01, C09, C12
- Budget: ~290 lines.

### F01-S4 — item snapshot and claim repository operations

- Goal: persist/rebuild snapshots, approve items, claim protocols with `FOR UPDATE SKIP LOCKED`, and write immutable outcomes.
- Files:
  - `apps/server_core/internal/modules/mutations/ports/repository.go`
  - `apps/server_core/internal/modules/mutations/adapters/postgres/apply_repository.go`
  - `apps/server_core/internal/modules/mutations/adapters/postgres/apply_repository_integration_test.go`
- Failing test first: two transactions claim one approved protocol; exactly one obtains it, while another installation remains independently claimable.
- Done:
  - Chunk query returns at most 20 items in sequence order.
  - One installation is processed sequentially.
  - `approved` and crash-resumed `applying` protocols are eligible.
  - Applied/failed/skipped items cannot be rewritten after terminalization.
  - Duplicate idempotency key is returned as already handled, not resent.
- Complexity: `complex`
- Serves: M03-C02, C03, C04
- Budget: ~300 lines.

### F01-S5 — WriterPort, programmable stub, and poller loop

- Goal: apply claimed items sequentially through the single WriterPort consumer and compute terminal state.
- Files:
  - `apps/server_core/internal/modules/mutations/ports/writer.go`
  - `apps/server_core/internal/modules/mutations/application/poller.go`
  - `apps/server_core/internal/modules/mutations/application/poller_test.go`
  - `apps/server_core/internal/modules/mutations/adapters/stub/writer.go`
  - `apps/server_core/internal/modules/mutations/adapters/stub/writer_test.go`
- Failing test first: programmable outcomes cover all-success, mixed, all-fail, duplicate, and failure-persistence runs.
- Done:
  - WriterPort has the named poller consumer and six F-02 intent implementations.
  - Idempotency key is exactly `{protocol_id}:{listing_id}`.
  - Unknown writer error becomes `internal`, `retryable=false`, with sanitized provider message.
  - No automatic retry converts a failure into success.
- Complexity: `complex`
- Serves: M03-C02, C03, C05
- Budget: ~290 lines.

### F01-S6 — crash resume and temporary composition registration

- Goal: prove restart behavior and make the poller runnable from the real composition root pending F-02 live wiring.
- Files:
  - `apps/server_core/internal/modules/mutations/application/poller_integration_test.go`
  - `apps/server_core/internal/modules/mutations/background/poller.go`
  - `apps/server_core/internal/composition/root.go`
  - `apps/server_core/internal/composition/root_test.go`
- Failing test first: stop after the first provider outcome, restart, and assert the first idempotency key is not resent and remaining items finish.
- Done:
  - Poller receives an injectable interval/tick in tests.
  - Root creates repository/service/poller.
  - Temporary stub wiring carries a dated deferral: `2026-07-16, replaced by F02-S8`.
  - No HTTP routes exist yet.
  - No bus, queue, scheduler framework, or outbox.
- Complexity: `complex`
- Serves: M03-C02, C03
- Budget: ~260 lines.

## F-01 WRITE-SET / write-DAG

```text
S1 migration shape
 ├─> S3 protocol repository
 └─> S4 apply repository
S2 lifecycle ─> S3 ─> S4 ─> S5 poller ─> S6 crash/root
```

Serialization points:

- `ports/repository.go`: S3 then S4.
- `repository_integration_test.go` family: S3 before S4.
- `apps/server_core/internal/composition/root.go`: F01-S6, then F02-S8, then F03-S9.
- `root_test.go`: same order.

---

# F-02 — write-types-adapters

> **Load-bearing (F01-S5b re-review, D-13 cond 2):** the crash window between provider send
> and pool-committed outcome is closed ONLY by provider-side idempotency dedup. Every F-02
> writer adapter MUST forward the item idempotency key `{protocol_id}:{listing_id}` on every
> provider write; an adapter that drops the key reintroduces double-apply. Items found in
> `applying` state on resume are resent under the same key by design.

## Slice cards

### F02-S1 — seven application-layer gates

- Goal: centralize all provider-write gates in the mutations application layer.
- Files:
  - `apps/server_core/internal/modules/mutations/application/gates.go`
  - `apps/server_core/internal/modules/mutations/application/gates_test.go`
  - `apps/server_core/internal/modules/mutations/application/errors.go`
- Failing test first: one dedicated negative test for each governance gate.
- Done:
  - Actor: missing/blank rejects before insert.
  - Idempotency: duplicate skips without writer call.
  - Execute: anything except literal `true` rejects.
  - Resolved link: price/stock/edit fail `link_unresolved`.
  - Policy: price/stock fail `policy_missing`; no fallback policy.
  - Source timestamp: absent rejects `source_time_unavailable`; stale fails `stale_source`.
  - Before/after audit: audit persistence must succeed before any provider call.
- Complexity: `complex`
- Serves: M03-C04, C12
- Budget: ~280 lines.

### F02-S2 — connector capability contracts

- Goal: add canonical price and listing-write capabilities and wire the already implemented StockWriter capability into the registry.
- Files:
  - `apps/server_core/internal/modules/connectors/domain/capability.go`
  - `apps/server_core/internal/modules/connectors/ports/marketplace_capability.go`
  - `apps/server_core/internal/modules/connectors/application/marketplace_capability_service.go`
  - `apps/server_core/internal/modules/connectors/application/marketplace_capability_service_test.go`
- Failing test first: capability lookup tests fail for missing `PriceWrites`, `ListingWrites`, and `StockWrites`.
- Done:
  - New `PriceWriter` has the mutations adapter as named consumer.
  - A single `ListingWriter` covers pause/edit; it has the same named consumer.
  - Existing StockWriter is exposed live rather than merely implemented.
  - Canonical request/result types contain no ML DTO.
- Complexity: `standard`
- Serves: M03-C04, C05
- Budget: ~260 lines.

### F02-S3 — Mercado Livre PriceWriter

- Goal: implement absolute price update in the ML adapter.
- Files:
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/price_writer.go`
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/price_writer_test.go`
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go`
- Failing test first: local `httptest` verifies method/path/body/idempotency forwarding and rejects malformed canonical input before network.
- Done:
  - Only absolute `{amount,currency}` is accepted.
  - ML request/response structs remain in `adapters/mercado_livre`.
  - 401/403, 429, validation, timeout/5xx, and unknown failures preserve sanitized provider message.
  - No live ML request occurs in tests.
- Complexity: `standard`
- Serves: M03-C04, C05
- Budget: ~280 lines.

### F02-S4 — ML pause/edit and read-backed resync

- Goal: implement listing pause/edit provider writes and the existing-reader resync path.
- Files:
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/listing_writer.go`
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/listing_writer_test.go`
  - `apps/server_core/internal/modules/mutations/adapters/listings/resync_writer.go`
  - `apps/server_core/internal/modules/mutations/adapters/listings/resync_writer_test.go`
- Failing test first: pause/edit payload tests plus resync test proving provider read and canonical listing upsert.
- Done:
  - Paused/closed remote state maps honestly.
  - Edit accepts only canonical attribute inputs.
  - Resync reuses existing listing ingestion/repository behavior; no duplicate provider parser.
  - Provider DTOs do not escape adapters.
- Complexity: `complex`
- Serves: M03-C05, C06
- Budget: ~300 lines.

### F02-S5 — product-link and policy adapters

- Goal: expose product_links resolution and pricing-policy reads to mutations through local adapters.
- Files:
  - `apps/server_core/internal/modules/mutations/ports/linkage.go`
  - `apps/server_core/internal/modules/mutations/ports/policy.go`
  - `apps/server_core/internal/modules/mutations/adapters/productlinks/writer.go`
  - `apps/server_core/internal/modules/mutations/adapters/productlinks/writer_test.go`
  - `apps/server_core/internal/modules/mutations/adapters/marketplaces/policy_reader.go`
  - `apps/server_core/internal/modules/mutations/adapters/marketplaces/policy_reader_test.go`
- Failing test first: unresolved link and missing policy produce the IC-03 codes with zero provider calls.
- Done:
  - `approve_candidate`, `manual_resolve`, and `reject_listing` delegate to the published product_links application API.
  - Repeating an identical link outcome returns `skipped`.
  - Conflicting re-resolution becomes `conflict_remote_changed`.
  - Policy absence remains absence.
- Complexity: `standard`
- Serves: M03-C04, C06
- Budget: ~280 lines.
- Design note (post-implementation, G2): the adapter checks resolution state twice —
  preflight via ListLinkWorkflows (catches `skipped`/`conflict_remote_changed` WITHOUT
  invoking the resolution API, keeping repeated approvals side-effect-free) and post-call
  via the resolution audit entry (source of truth when the write actually ran).
  Alternative considered: audit-only single check — rejected because idempotent re-sends
  would re-invoke the resolution API for outcomes already known identical, producing
  spurious audit rows. Preflight-only was rejected because the list read races the write;
  the audit entry is authoritative for what the API actually did.

### F02-S6 — six-intent writer router and failure mapping

- Goal: dispatch the six enabled intents and map connector errors to IC-03 failures.
- Files:
  - `apps/server_core/internal/modules/mutations/application/writer.go`
  - `apps/server_core/internal/modules/mutations/application/writer_test.go`
  - `apps/server_core/internal/modules/mutations/adapters/connectors/failure_mapping.go`
  - `apps/server_core/internal/modules/mutations/adapters/connectors/failure_mapping_test.go`
- Failing test first: table covers every intent plus ML error mapping for 429, 401/403, validation, paused, timeout/5xx, and unknown.
- Done:
  - `provider_rate_limited` and `provider_unavailable` are retryable.
  - Unknown maps to `internal`, non-retryable.
  - `message_provider` is retained as sanitized text only.
  - Listing edit changing SELLER_SKU away from linked CODPROD fails `sku_invariant_violation` before adapter call.
  - `listing_create` validates its IC-03 shape then returns typed `type_not_enabled`.
- Complexity: `complex`
- Serves: M03-C04, C05, C06, C07, C12
- Budget: ~300 lines.

### F02-S7 — stock_correct and StockActionService fold

- Goal: route stock correction through the envelope while preserving legacy stock-action history.
- Files:
  - `apps/server_core/migrations/0039_inventory_stock_actions_mutation_protocol.sql`
  - `apps/server_core/internal/modules/inventory/domain/action.go`
  - `apps/server_core/internal/modules/inventory/ports/stock_action.go`
  - `apps/server_core/internal/modules/inventory/adapters/postgres/stock_action_repo.go`
  - `apps/server_core/internal/modules/inventory/application/stock_action_service.go`
  - `apps/server_core/internal/modules/inventory/application/stock_action_service_test.go`
  - `apps/server_core/internal/modules/inventory/application/manual_action_facade.go`
  - `apps/server_core/internal/modules/inventory/application/manual_action_facade_test.go`
  - `apps/server_core/internal/modules/inventory/adapters/mutations/envelope.go`
- Failing test first: existing manual stock endpoint creates one envelope protocol and one linked legacy history record, with no direct provider call.
- Done:
  - Migration adds nullable `mutation_protocol_id`; it never defaults.
  - Existing history reads remain valid.
  - StockActionService no longer calls StockWriter directly.
  - Approved stock action delegates through mutations and preserves prior risk/policy validation.
  - Existing regression suite remains green.
- Complexity: `complex`
- Serves: M03-C04, C12
- Budget: split implementation into domain/repository and service commits if either exceeds ~300 lines; both remain within this slice contract.

### F02-S8 — live composition wiring and F-01 stub removal

- Goal: replace the dated F-01 stub with all real application adapters.
- Files:
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go`
  - `apps/server_core/internal/composition/root.go`
  - `apps/server_core/internal/composition/root_test.go`
- Failing test first: root test proves mutations receives price/stock/listing/link/resync writers and legacy inventory receives the envelope adapter.
- Done:
  - `ProviderCapabilitySet()` adds `PriceWrites: a`, `StockWrites: a`, `ListingWrites: a`.
  - Root removes the programmable writer from the live path.
  - Inventory direct `inventoryconnectors.NewStockWriter(...)` wiring is removed from StockActionService.
  - Stub remains test-only.
  - Live ML execution remains disabled unless the operator authorizes the provider-write lane.
- Complexity: `complex`
- Serves: M03-C04, C05, C06
- Budget: ~220 lines.

## F-02 WRITE-SET / write-DAG

```text
S1 gates ────────────────┐
S2 connector contracts ─> S3 price ─┐
                         └> S4 listing/resync ─┐
S5 links/policy ────────────────────────────────┤
                                               ├─> S6 router
F01 repository ────────────────────────────────┘
S6 router ─> S7 stock fold ─> S8 live root wiring
```

Serialization points:

- `connectors/.../capability_adapter.go`: S3, S4, then S8.
- `connectors/domain/capability.go`, ports, capability service: S2 owns the shape; later slices consume only.
- `inventory/application/stock_action_service.go`: F02-S7 exclusive.
- `composition/root.go`: F01-S6 → F02-S8 → F03-S9.
- Migration order: `0038` → `0039`.

---

# F-03 — selection-preview-api

## Slice cards

### F03-S1 — shared IC-02 filter grammar and selection resolver

- Goal: extract semantic filter validation once and reuse it for URL listings and JSON mutation selection.
- Files:
  - `apps/server_core/internal/modules/listings/domain/filter.go`
  - `apps/server_core/internal/modules/listings/domain/filter_test.go`
  - `apps/server_core/internal/modules/listings/transport/query.go`
  - `apps/server_core/internal/modules/listings/transport/query_test.go`
  - `apps/server_core/internal/modules/mutations/adapters/listings/selection_resolver.go`
  - `apps/server_core/internal/modules/mutations/adapters/listings/selection_resolver_test.go`
- Failing test first: the same table of IC-02 keys/values is run through both URL and mutation JSON entry points.
- Done:
  - Exact grammar: status, sync_state, link_state, exception, has_exception, listing_type_code, product_id, plus separate `q`.
  - Unknown keys and arrays/operators are rejected.
  - Selection iteration uses the existing listings read service/cursors and stops after item 2001.
  - No duplicate parser or SQL filter implementation.
- Complexity: `complex`
- Serves: M03-C08
- Budget: ~300 lines.

### F03-S2 — create and preview application services

- Goal: create drafts, resolve explicit/filter selections, and atomically replace preview snapshots.
- Files:
  - `apps/server_core/internal/modules/mutations/application/service.go`
  - `apps/server_core/internal/modules/mutations/application/preview.go`
  - `apps/server_core/internal/modules/mutations/application/preview_test.go`
  - `apps/server_core/internal/modules/mutations/adapters/postgres/preview_repository.go`
  - `apps/server_core/internal/modules/mutations/adapters/postgres/preview_repository_integration_test.go`
- Failing test first: preview, underlying listing mutation, re-preview, and source-time-unavailable cases.
- Done:
  - Empty → `empty_selection`.
  - 2001st result → `selection_too_large`.
  - Snapshot captures canonical before/after and freshest known source time.
  - Re-preview replaces old items transactionally only in draft/previewed.
  - No provider call occurs during preview.
- Complexity: `complex`
- Serves: M03-C04, C07, C08
- Budget: ~300 lines.

### F03-S3 — approve and cancel application services

- Goal: enforce execute, TTL, and lifecycle transitions below transport.
- Files:
  - `apps/server_core/internal/modules/mutations/application/approval.go`
  - `apps/server_core/internal/modules/mutations/application/approval_test.go`
- Failing test first: execute absent/false, exactly 15 minutes, over 15 minutes, and cancel-after-approved.
- Done:
  - Only literal `execute:true` approves.
  - Preview older than 15 minutes returns `preview_stale`.
  - Approve promotes the existing item snapshot; it never re-resolves listings.
  - Cancel only works from draft/previewed.
- Complexity: `complex`
- Serves: M03-C01, C04, C08
- Budget: ~220 lines.

### F03-S4 — retry clone and read application services

- Goal: clone only retryable failures and provide stable list/detail/item pagination.
- Files:
  - `apps/server_core/internal/modules/mutations/application/retry.go`
  - `apps/server_core/internal/modules/mutations/application/retry_test.go`
  - `apps/server_core/internal/modules/mutations/application/read.go`
  - `apps/server_core/internal/modules/mutations/application/read_test.go`
  - `apps/server_core/internal/modules/mutations/adapters/postgres/read_repository.go`
  - `apps/server_core/internal/modules/mutations/adapters/postgres/read_repository_integration_test.go`
- Failing test first: retry a mixed terminal protocol and compare the original rows byte-for-byte afterward.
- Done:
  - New tenant-scoped `MP-…` ID.
  - `retried_from` set.
  - Only retryable failed items cloned.
  - Zero eligible failures returns `nothing_to_retry`.
  - Items sort by `seq ASC`; protocols by `created_at DESC`.
- Complexity: `complex`
- Serves: M03-C09, C12
- Budget: ~300 lines.

### F03-S5 — command HTTP endpoints

- Goal: expose create, preview, approve, cancel, and retry with thin handlers.
- Files:
  - `apps/server_core/internal/modules/mutations/transport/command_handler.go`
  - `apps/server_core/internal/modules/mutations/transport/command_handler_test.go`
  - `apps/server_core/internal/modules/mutations/transport/errors.go`
- Failing test first: status/code table for invalid body, actor, intent, type_not_enabled, empty/large selection, execute, stale preview, illegal state, unknown protocol, and nothing-to-retry.
- Done:
  - Strict JSON decoding with unknown-field rejection.
  - Actor is taken only from the body; never synthesized from auth.
  - `listing_create` returns 422 and inserts zero rows.
  - No combined preview+approve route.
- Complexity: `standard`
- Serves: M03-C01, C04, C07, C08, C09
- Budget: ~300 lines.

### F03-S6 — query HTTP endpoints

- Goal: expose list, detail, and cursor-paginated items.
- Files:
  - `apps/server_core/internal/modules/mutations/transport/query_handler.go`
  - `apps/server_core/internal/modules/mutations/transport/query_handler_test.go`
- Failing test first: tenant isolation, invalid cursor/filter, unknown protocol, and stable ordering.
- Done:
  - `GET /mutations` requires installation and accepts IC-03 state/type filters.
  - `GET /mutations/{id}` and `/items` never cross tenants.
  - Limits/cursors use existing SDK error envelope conventions.
  - No secret, token, header, or buyer data is serialized.
- Complexity: `standard`
- Serves: M03-C09, C11, C12
- Budget: ~260 lines.

### F03-S7 — command OpenAPI and SDK, atomic commit

- Goal: add command-side contract and matching SDK methods in one commit.
- Files:
  - `contracts/api/marketplace-central.openapi.yaml`
  - `packages/sdk-runtime/src/index.ts`
  - `packages/sdk-runtime/src/index.test.ts`
- Failing test first: SDK request tests for `createMutation`, `previewMutation`, `approveMutation`, `cancelMutation`, `retryMutationFailures`.
- Done:
  - Add POST operations on `/mutations`, `/mutations/{id}/preview`, `/approve`, `/cancel`, `/retry`.
  - Add base Mutation schemas and command request/response schemas.
  - Every IC-03 status/code is represented.
  - No dashboard, orders, sync-run, market, or category-attribute section touched.
- Complexity: `complex`
- Serves: M03-C07, C08, C09
- Budget: ~300 lines.

### F03-S8 — query OpenAPI and SDK, atomic commit

- Goal: add the three read operations and cursor types in a second OpenAPI+SDK atomic commit.
- Files:
  - `contracts/api/marketplace-central.openapi.yaml`
  - `packages/sdk-runtime/src/index.ts`
  - `packages/sdk-runtime/src/index.test.ts`
- Failing test first: URL encoding and response typing for `listMutations`, `getMutation`, and `listMutationItems`.
- Done:
  - Add GET on `/mutations`, GET `/mutations/{id}`, GET `/mutations/{id}/items`.
  - Exactly eight IC-03 operations exist after S8.
  - Mutation item order and protocol list order are documented.
  - OpenAPI and SDK remain handler-parity green.
- Complexity: `standard`
- Serves: M03-C09, C11, C12
- Budget: ~260 lines.

### F03-S9 — route registration and full lifecycle integration

- Goal: mount the eight endpoints and prove the complete stub-lane lifecycle/error matrix.
- Files:
  - `apps/server_core/internal/composition/root.go`
  - `apps/server_core/internal/composition/root_test.go`
  - `apps/server_core/internal/modules/mutations/integration/lifecycle_test.go`
  - `apps/server_core/internal/modules/mutations/integration/error_matrix_test.go`
- Failing test first: root route smoke and create→preview→approve→poll→terminal transcript.
- Done:
  - All eight routes are reachable.
  - Every transport error row and every applicable item-level error is asserted.
  - Snapshot survives underlying listing change.
  - Retry creates a new immutable protocol.
  - Integration uses ephemeral PostgreSQL plus the planned stub adapter; it makes no live-integration claim.
- Complexity: `complex`
- Serves: M03-C01–C09, C12
- Budget: ~300 lines.

## F-03 WRITE-SET / write-DAG

```text
S1 shared filter ─> S2 preview ─> S3 approve
F02 writer/gates ────────────────┘
F01 repository ─> S4 retry/read
S2/S3/S4 ─> S5 command HTTP
S4 ──────> S6 query HTTP
S5 ─> S7 command contract
S6 ─> S8 query contract
S7 + S8 ─> S9 registration/integration
```

Serialization points:

- `listings/transport/query.go`: F03-S1 only; existing behavior must remain byte-compatible.
- `application/service.go`: S2 owns constructor composition; S3/S4 use separate files.
- OpenAPI and SDK files: S7 then S8, each as its own atomic OpenAPI+SDK commit.
- `composition/root.go`: F03-S9 after F02-S8.

---

# F-04 — preview-confirm-ui

Implementation begins only after the hub confirms M-02 F-03 merged and triggers the rebase.

## Slice cards

### F04-S1 — mutation query hooks and terminal observer

- Goal: bind SDK methods to exact IC-05 mutation keys and encapsulate 2-second terminal polling.
- Files:
  - `packages/feature-listings/src/mutations/useMutationProtocol.ts`
  - `packages/feature-listings/src/mutations/useMutationProtocol.test.tsx`
- Failing test first: applying polls every 2000 ms, terminal stops polling, reconnect resumes.
- Done:
  - Only `mutationsQueryKeys.detail/items` are used.
  - No direct fetch.
  - Protocol truth always comes from the server.
  - Terminal callback is guarded by protocol ID so invalidation occurs exactly once.
- Complexity: `complex`
- Serves: M03-C10, C11
- Budget: ~220 lines.

### F04-S2 — intent forms and preview state machine

- Goal: implement six intent inputs and the exact modal transition graph.
- Files:
  - `packages/feature-listings/src/mutations/MutationIntentForm.tsx`
  - `packages/feature-listings/src/mutations/MutationIntentForm.test.tsx`
  - `packages/feature-listings/src/mutations/useMutationFlow.ts`
  - `packages/feature-listings/src/mutations/useMutationFlow.test.tsx`
- Failing test first: intent→previewing→preview-shown plus stale-preview recovery and double-submit prevention.
- Done:
  - Supports price_update, stock_correct, link_apply, listing_pause, listing_resync, listing_edit.
  - Sends `operator_supplied_unverified` explicitly.
  - 409 `preview_stale` returns to preview with exact copy: “Prévia expirada. Gere novamente.”
  - Double click produces one approve request.
  - Cancel calls cancelMutation only from allowed UI states.
- Complexity: `complex`
- Serves: M03-C10
- Budget: ~300 lines.

### F04-S3 — preview/confirm modal rendering

- Goal: render totals and before→after rows before enabling the two-step confirmation.
- Files:
  - `packages/feature-listings/src/mutations/MutationPreviewModal.tsx`
  - `packages/feature-listings/src/mutations/MutationPreviewModal.test.tsx`
  - `packages/feature-listings/src/mutations/MutationItemTable.tsx`
- Failing test first: confirm remains disabled until preview rows render and checkbox is checked.
- Done:
  - Explicit checkbox plus button; no one-click apply.
  - Progress/result views follow the server lifecycle.
  - `selection_too_large` gives narrowing guidance and cancels the draft.
  - Failures use `failureCopy`; provider text appears only behind “▸ técnico”.
  - Terminal calls `invalidateAfterMutation` once.
- Complexity: `complex`
- Serves: M03-C10, C12
- Budget: ~300 lines.

### F04-S4 — Anúncios integration

- Goal: replace M-02’s disabled bulk actions with modal launch using its persisted selection IDs.
- Expected post-M02 files:
  - `packages/feature-listings/src/AnunciosPage.tsx`
  - `packages/feature-listings/src/AnunciosPage.test.tsx`
  - `packages/feature-listings/src/index.ts`
- Failing test first: selection across pages opens the requested intent with the exact composite listing IDs.
- Done:
  - Six buttons map to the six enabled IC-03 types.
  - Selection mode is explicit for checked IDs; no client parsing of listing IDs.
  - Modal close does not cancel applying work.
  - No web-query, context, nav, Vite, or failureCopy file is modified.
- Complexity: `standard`
- Serves: M03-C10
- Budget: ~230 lines.
- Gate: blocked by PD-02 until actual M-02 paths are confirmed.

### F04-S5 — protocolo detail page and assigned route row

- Goal: add the IC-05 detail route, polling, failure inspection, and retry navigation.
- Files:
  - `packages/feature-listings/src/ProtocolDetailPage.tsx`
  - `packages/feature-listings/src/ProtocolDetailPage.test.tsx`
  - `packages/feature-listings/src/index.ts`
  - `apps/web/src/app/AppRouter.tsx`
  - `apps/web/src/app/AppRouter.test.tsx`
- Failing test first: deep-link reload in `applying`, terminal convergence, exactly-once invalidation, and retry navigation with backlink.
- Done:
  - Route is exactly `/protocolos/:protocolId`.
  - F5 reconstructs state from server and resumes 2-second polling.
  - Header shows type, actor, timestamps, state.
  - Each item shows before/after and fixed pt-BR failure copy.
  - Retry navigates to the new protocol and renders `retried_from`.
  - Only the M-03-assigned route row is added.
- Complexity: `complex`
- Serves: M03-C10, C11, C12
- Budget: ~300 lines.

## F-04 WRITE-SET / write-DAG

```text
M-02 F-03 merge + rebase
 ├─> S1 polling hook ─┐
 └─> S2 flow/forms ───┤
                      ├─> S3 modal ─> S4 Anúncios integration
S1 + SDK read methods ──────────────> S5 detail/route
```

Serialization points:

- `packages/feature-listings/src/index.ts`: S4 then S5.
- `AnunciosPage.tsx`: F04-S4 only, after M-02 has released it.
- `apps/web/src/app/AppRouter.tsx` and test: F04-S5 only after CHIP-M02 closes/releases the frontend seam.
- `packages/web-query/src/index.ts` and `failureCopy.ts`: read-only M-02 outputs; M-03 must not edit them.
- `apps/web/vite.config.ts`: no M-03 write; `/mutations` is M-02 F-02’s assigned row.

---

# Contract-satisfiability check

## Current OpenAPI state

Verified current state:

- No `/mutations` path exists.
- No `Mutation*` schema or mutation operationId exists.
- Current mutation grep result is empty.
- Existing path area is `/listings/refresh` followed by `/catalog/taxonomy`.
- Components begin at `components.schemas`; listing schemas currently end with `RefreshListingsAccepted`, followed by `CanonicalNumericSourceFact`.

Exact additive insertion:

- Paths: insert the mutation block immediately after `/listings/refresh` and before `/catalog/taxonomy`.
- Schemas: insert all `Mutation*` schemas immediately after `RefreshListingsAccepted` and before `CanonicalNumericSourceFact`.
- SDK types: insert after `ListingSummary` and before `IntegrationConnectionSnapshot`.
- SDK methods: insert immediately after `refreshListings` and before deprecated catalog aliases.
- SDK tests: mutation-only describe block adjacent to existing listings tests.

Claimed paths:

1. `POST /mutations`
2. `GET /mutations`
3. `POST /mutations/{id}/preview`
4. `POST /mutations/{id}/approve`
5. `POST /mutations/{id}/cancel`
6. `POST /mutations/{id}/retry`
7. `GET /mutations/{id}`
8. `GET /mutations/{id}/items`

No path is already occupied.

## Sibling-track collision check

M-03 does not claim or touch:

- `/dashboard*`
- `/orders*`
- `/sync*`
- `/market*`
- `/listings/categories/{category_id}/attributes`
- their schemas, SDK methods, handlers, migrations, or module registrations.

The OpenAPI logical sections are disjoint from CHIP-SAT. `packages/sdk-runtime/src/index.ts` is physically shared, but the W1 contract pre-assigns disjoint mutation versus SAT insertion anchors; hub must serialize/reconcile regeneration/merge if sibling edits move those anchors.

F-04 does not claim CHIP-M02’s:

- InstallationContext
- Layout/sidebar/nav
- mutation query-key definitions
- invalidation crosswalk implementation
- failureCopy implementation
- state components
- Vite `/mutations` proxy row

It consumes those outputs after rebase and writes only the assigned protocolo route plus M-03-owned mutation UI files.

## Migration allocation

Current maximum is `0037`.

Planned use:

- `0038_mutation_protocols.sql`
- `0039_inventory_stock_actions_mutation_protocol.sql`

Therefore `0038–0042` is sufficient, leaving `0040–0042` unused. No slice may take `0043+`. An unforeseen third or later migration may use only the remaining reserved numbers; exhausting `0042` requires a hub `REQUEST`.

---

# Pre-identified additive contract-locks

## Server composition root lock

File: `apps/server_core/internal/composition/root.go`

M-03 adds only these registration/wiring hunks:

1. Mutation module imports:
   - postgres repository
   - mutations application
   - poller/background
   - mutation transport
   - local adapters for connectors, listings, product_links, and marketplaces

2. After listings repository/services are available:
   - construct mutation repository with `pool` and `cfg.DefaultTenantID`
   - construct selection/link/policy/resync adapters from existing published services
   - construct mutation application service/writer router

3. Background registration adjacent to existing background starters:
   - start the in-process mutation poller
   - F01 temporary stub carries dated F02 replacement
   - F02 replaces it with live adapters

4. HTTP registration after F03:
   - `mutationstransport.NewHandler(...).Register(mux)`

5. Inventory fold:
   - replace direct StockWriter injection into StockActionService with the mutations envelope adapter

No dashboard/orders/sync/market module registration line is claimed.

## Connector wiring lock — F-02 only

Files:

- `apps/server_core/internal/modules/connectors/application/marketplace_capability_service.go`
- `apps/server_core/internal/modules/connectors/ports/marketplace_capability.go`
- `apps/server_core/internal/modules/connectors/domain/capability.go`
- `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go`

Exact additive wiring:

- `ProviderCapabilitySet.PriceWrites`
- `ProviderCapabilitySet.StockWrites` populated live
- `ProviderCapabilitySet.ListingWrites`
- `MarketplaceCapabilityService.PriceWriter(providerCode)`
- existing `StockWriter(providerCode)` now succeeds for ML
- `MarketplaceCapabilityService.ListingWriter(providerCode)`
- ML `ProviderCapabilitySet()` fields:
  - `PriceWrites: a`
  - `StockWrites: a`
  - `ListingWrites: a`

These lines are owned exclusively by F-02 and released only at milestone CLOSED.

---

# Verification ladder

Each slice first runs its targeted tests. Feature acceptance then follows the profile:

- L0:
  - absolute workspace `GOCACHE`
  - `go build ./...`
  - web typecheck/build when applicable
  - `npm run harness:governance -- -BaseSha <full accepted 40-hex SHA>` from a clean detached worktree

- L1:
  - full `go test ./...` because migrations and composition are touched
  - mutation integration packages through `npm run harness:integration`
  - web Vitest for F-04
  - unchanged `TestPhase1SmokeFlow` failure is cited only through the ratified allowlist

- L2:
  - hub-owned dev stack only
  - route/error-shape/OpenAPI↔SDK↔handler parity
  - stub-adapter lifecycle smoke
  - no chip-owned process on ports 8080/5174

- L3/L4:
  - browser walkthrough for M03-C10/C11
  - fresh QA against M03-C01–C12
  - live ML write remains optional and requires explicit operator authorization; stub evidence never claims live provider integration.