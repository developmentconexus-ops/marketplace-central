# ADR-024: Single writer for order ingest

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed order ingest under the two-digit citation
`ADR-04` inside MIS-007, and was never given a document of its own. The two-digit number
collided with three unrelated decisions from other missions and with the pre-existing
`004-integration-catalog-plugin-framework.md`, which is about a different subject
entirely (provider plugin self-registration). This document reconstructs the MIS-007
single-writer rule under its own global number, from the 13 live-code citations harvested
at `docs/architecture/decisions/_citations/adr-04-twodigit-citations.md` (Assertion A1).
Every clause below is traceable to code that already asserts it. Nothing here is new
policy.

## Context

Before this rule, an order could in principle be written to the database from more than
one place: the batch import path that walks `ListOrders`, and — per the milestone's own
roadmap — a future backfill job, an incremental sync job, and a webhook worker, each
capable of independently deciding what an order's current state is and persisting it.
Multiple independent writers for the same resource-addressed entity is a race: two
triggers can observe the provider at different moments, compute different derived state
(bucket, shipment status, faturado handling), and the one that commits last silently wins
regardless of which observation was actually newer.

## Decision

**Every trigger that can mutate an order — import enumeration today, and backfill,
incremental sync, and the webhook worker as they are added — converges on exactly one
write path: `IngestService.IngestOrder`. No other code path calls the order store's write
methods.**

**§1 — `IngestOrder` is the one write path.** One order, one atomic transaction: every
fetch (order detail, shipment detail, buyer-fiscal) runs before the repository call, so
any real error from any of them returns without touching the store at all.
> `apps/server_core/internal/modules/orders/application/ingest_service.go:65-66` — "IngestOrder is
> IC-06's single write path (ADR-04). One order, one atomic transaction, no partial
> writes"

**§2 — The port is named for this, not for CRUD.** `OrderIngestor` exposes exactly one
method, `IngestOrder(ctx, installationID, providerOrderID) error` — it takes an
identifier, not a payload, forcing every caller to go through the fetch-then-write
sequence above rather than constructing state itself and handing it to the store.
> `apps/server_core/internal/modules/orders/ports/order_ingestor.go:5-9` — "OrderIngestor is
> IC-06's single write-path port (ADR-04): every trigger that mutates an order ...
> converges on this ONE call per order."

**§3 — The store-level write transaction is likewise singular.** `OrderIngestStore`
persists exactly one order (header, items, payments, optional shipment row)
atomically per call; there is no batch-write method on this port.
> `apps/server_core/internal/modules/orders/ports/order_store.go:67-77` — "OrderIngestStore is
> IC-06's single-transaction write port (ADR-04): IngestOrder persists exactly one order
> ... atomically"

**§4 — The repository implementation is the literal convergence point.** The Postgres
adapter's `IngestOrder` reuses the same `upsertOrder`/`replaceItems`/`replacePayments`
freshness-guard and line-identity-reconciliation invariants as the older batch
`UpsertOrders` path, so a concurrent ingest of the same order from two different triggers
cannot regress a newer snapshot with an older one.
> `apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go:809-815` —
> "IngestOrder persists ONE order ... atomically — IC-06's single write path (ADR-04):
> every ingest trigger ... converges here."

**§5 — The old batch importer no longer writes.** `ImportService.Import` was the
pre-existing write path (`ListOrders` → `normalizeOrders` → `UpsertOrders`). It is now
reduced to enumeration: it walks `ListOrders` for `provider_order_id`s and calls
`IngestOrder` once per id. It holds an `OrderIngestor`, not an `OrderStore`.
> `apps/server_core/internal/composition/root.go:591-592` — "F-02: IngestOrder is the
> single write path (ADR-04) — Import now only enumerates provider_order_ids from
> ListOrders and delegates each one here."
> `apps/server_core/internal/modules/orders/application/import_service.go:50,56,92` — `Import`
> enumerates `ListOrders`' snapshot for `provider_order_id`s, then calls
> `s.ingestor.IngestOrder(ctx, installationID, providerOrderID)` once per id.

## Rationale

A resource-addressed write (one order, identified by `installationID` +
`providerOrderID`) has exactly one correct current state at any moment. Letting several
independent code paths each compute and persist that state means the correctness of the
data depends on which trigger happened to run last — an accident of scheduling, not a
fact about the order. Collapsing every trigger onto one function that always re-fetches
before writing removes the possibility of two different derivations disagreeing: there is
only one derivation, called from every entry point.

## Consequences

- New order-mutating triggers (backfill, incremental sync, webhook worker) must call
  `IngestOrder` per order; adding a second write path for any of them silently reopens
  this hazard and is not caught by any test that only exercises `IngestOrder` itself.
- `BuyerFiscal` is fetched on every `IngestOrder` call, unlike the list-path
  `EnrichService.Enrich`, which skips it for a request-deadline reason that does not apply
  to ingest running off the interactive path. This is a deliberate asymmetry between the
  two, not an oversight.
- The batch `OrderStore.UpsertOrders` path still exists in the repository and is reused
  internally by `IngestOrder` for its upsert/freshness logic, but `ImportService` no
  longer calls it directly — it is no longer a second external write path.

## Alternatives Considered

**Let each trigger write independently, guarded by a `date_last_updated` freshness
check at the row level.** Rejected in the code as implemented: a freshness guard at the
row level still allows two triggers to run the full derivation (bucket, shipment mapping)
independently and only reconciles at the final write, which does not prevent divergent
reads of ephemeral state (e.g., in-flight shipment lookups) between the two derivations.
The chosen design removes the derivation duplication itself, not just its output.

## Unverified claims

None — every clause above matches a verified anchor in code.
