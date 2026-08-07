# Architecture Decision Records

## The numbering rule

**Three digits, global, one number per decision, forever.** `ADR-017`, never `ADR-17`.
A number is never reused, never renumbered, and never scoped to a mission.

This rule exists because the opposite happened. Until 2026-08-05 there was no global
registry: each mission numbered its own decisions from ~ADR-01, so a single number named up
to four unrelated rules, and two-digit citations in code landed on three-digit documents
about entirely different subjects. `ADR-04` had 94 citations and four live meanings.
`ADR-17` had 1.378 citations and no document at all.

The full measurement, the old-to-new mapping, and the decisions that deliberately received
no document are in [`_citations/RENUMBERING-REGISTRY.md`](_citations/RENUMBERING-REGISTRY.md).
`_citations/` also holds, per number, the harvested citations each document was
reconstructed from — with `file:line` anchors, contradictions, and amendments.

**Citing:** `ADR-0NN`. A citation to a number with no document is a defect, not a
placeholder.

**Mission records under `.mnfs/` are frozen.** They still carry mission-local numbering and
are deliberately not rewritten — a mission record states what was decided at the time, with
the number of the time. Use the registry's crosswalk to read them.

## Product architecture

| # | Decision | Status |
|---|---|---|
| [001](001-metalshopping-direct-read.md) | MPC reads products directly from MetalShopping Postgres | superseded 2026-07-07 |
| [002](002-mpc-schema-same-cluster.md) | MPC tables live in the `mpc` schema on MetalShopping's cluster | superseded 2026-07-07 |
| [003](003-integration-spec-split-and-sequencing.md) | Integrations split into operational specs after foundation | accepted |
| [004](004-integration-catalog-plugin-framework.md) | Integration catalog plugin framework | accepted |
| [005](005-mercado-livre-first-control-plane.md) | Mercado Livre first control plane | accepted |
| [006](006-oracle-internal-read-owned-by-mpc.md) | MPC-owned Oracle internal reads | accepted |
| [007](007-godror-oci-oracle-runtime.md) | Godror with OCI is the canonical Oracle runtime | accepted |
| [008](008-production-deploy-topology.md) | Production deploy topology | accepted |
| [009](009-fee-value-carries-provenance.md) | Every fee value carries provenance (layer, origin, collected-at) | accepted |
| [010](010-mercado-livre-polling-visible-refresh.md) | Polling/GET only against ML; live data needs a visible refresh | accepted |
| [011](011-divergences-one-open-row-per-entity-kind.md) | Divergences: one open row per (entity, kind), detected at ingest | accepted |
| [012](012-difal-single-source-in-pricing.md) | DIFAL has a single source of truth inside `pricing` | accepted |
| [013](013-webhook-is-pointer-never-data.md) | A webhook body is a pointer, never domain data; always 200 | accepted |
| [014](014-market-collection-on-demand.md) | Market reference collection is on-demand | accepted |
| [015](015-listings-module-is-read-only.md) | The canonical `listings` module is read-only | accepted |
| [016](016-sdk-runtime-manual-same-commit.md) | `sdk-runtime` is hand-written; OpenAPI and SDK land in one commit | accepted |
| [017](017-unknown-is-never-zero.md) | **Unknown is never zero — honest absence end to end** | superseded 2026-08-07 by ADR-034 |
| [018](018-mutation-envelope-table-and-poller.md) | Mutation envelope is a protocol table plus an in-process poller | accepted |
| [019](019-listings-ingestion-feeds-snapshot-observer.md) | Listings ingestion feeds the snapshot observer, one row per item | accepted |
| [020](020-market-data-only-via-collectorport.md) | Market data only via `CollectorPort`; no scraping; honest-empty | accepted |
| [021](021-frontend-platform-seam-tanstack.md) | Frontend seam: one shell, TanStack Query exclusive for server state | accepted |
| [022](022-provider-write-sku-invariant.md) | Provider writes require `SELLER_SKU == CODPROD` | accepted |
| [023](023-module-protocol.md) | **The module protocol — private internals, one published consumable** | accepted |
| [024](024-single-writer-order-ingest.md) | `IngestOrder` is the single write path for orders | accepted |
| [025](025-raw-payload-selective-never-pii.md) | Raw provider payloads are stored selectively; PII never | accepted |
| [026](026-scheduler-phase-vocabulary.md) | Scheduler phase vocabulary: `backfill` / `incremental` / `sweep` | accepted, amended 2026-08-01 |
| [027](027-absent-is-not-closed.md) | Absent from a partial pull is not closed | accepted |
| [028](028-auto-link-only-on-concordant-anchors.md) | Auto-link only on concordant CODPROD + EAN anchors | accepted, amended D-121-2 |
| [029](029-resilience-decorator-no-retry-on-writes.md) | Resilience decorator: writes opt out of retry | accepted |
| [030](030-scheduler-second-instance-per-installation.md) | A second Scheduler instance per installation | accepted |
| [031](031-keep-absent-merge.md) | Products-mirror upsert keep-absent merge | accepted |
| [032](032-ml-catalog-offers-read-flag-defaults-off.md) | ML catalog-offers read gated by a flag that defaults off | accepted |
| [033](033-integracoes-entram-por-adapters.md) | External marketplace integrations enter through `adapters/marketplace/<vendor>`, not `connectors` | accepted |
| [034](034-fact-substitui-adr-017.md) | `internal/kernel/fact` supersedes ADR-017 | accepted |

## Harness process decisions

Decisions about how work is dispatched, gated and evidenced are **not** ADRs. Mixing them
into this series is what produced the numbering collision. They live in
[`docs/HARNESS-PROFILE.md` §13](../../HARNESS-PROFILE.md), keyed `P-1`…`P-5`, and are cited
as `HARNESS-PROFILE §13 P-N`.

## Reconstructed documents

Documents 009–022 and 024–030 were written on 2026-08-05 from the citations that already
asserted them — the rules governed the code long before anyone wrote them down. Each carries
a `**Reconstructed:**` header naming its harvest file, and every `file:line` quoted in them
was read and confirmed before being written.

Where a citation asserted something the code does not show, the document says so under
`## Unverified claims` rather than stating it as a rule. Those sections are the honest edge
of this reconstruction and are the first thing to resolve when the relevant code is next
touched.
