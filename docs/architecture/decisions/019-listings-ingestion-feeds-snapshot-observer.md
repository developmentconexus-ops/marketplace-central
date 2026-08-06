# ADR-019: Listings ingestion must keep feeding the snapshot observer, one row per item

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed MIS-007 (`ml-sync`) under the label "ADR-13" —
a different decision than the MIS-003 "ADR-13" reconstructed separately as ADR-018. No
standalone document existed for either sense; both lived only as cited rows across
mission, milestone, interface-contract, and Go source files. This document reconstructs
the MIS-007 sense from the ~28 live citations harvested at
`docs/architecture/decisions/_citations/adr-013-citations.md`, Assertion A2, and receives
the new global number 019 per
`docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md`.

## Context

`listings` ingestion in MIS-007 had two independent consumers of the same hydrated
provider data: the canonical `listings`/`listing_variations` read model, and
`product_links`' re-linking matcher, which watches every ingested item for EAN/SKU
anchors via `AbsorbProviderSnapshots` (`connectors/source.go`). The existing single-item
`ReadListing` path fed both. MIS-007 replaced single-item hydration with a
batch/multiget-based backfill and sweep path for throughput. A batch rewrite that only
checked its own read-model output against the old path's read-model output could pass
every listings test while silently starving the snapshot observer — the matcher would
degrade with no error, no failing test, and no visible signal, because "the listings
table looks right" is not the same claim as "the matcher is still being fed."

The other half of the same decision is a modeling question the migration to multiget
hydration forced: an item with variations was, before this decision, at risk of being
represented as multiple rows in `listings` (one per variation) rather than one row with
children. Flattening would break the parent-table PK invariant (`listings` keeps a
sentinel `'-'` in its `variation_id` column for items with no variation, established in
migration 0036) and would radiate a destructive PK change into modules 0022 and 0025
that already depend on it.

## Decision

**Any new or rewritten listings ingestion path must (a) go on feeding the product-links
snapshot observer non-regressively, at the same anchor coverage as the path it replaces,
and (b) emit exactly one `listings` row per item, with variations as child rows in
`listing_variations` rather than a flattened row per variation. The `listings` primary
key, including its `'-'` sentinel for the no-variation case, is not touched by this
decision.**

**§1 — One row per item, never a flattened row per variation.** A hydrated item becomes
one canonical `listings` row (`variation_id = NoVariationID`, the sentinel), with its
variations attached as children, not additional top-level rows.
> `apps/server_core/internal/modules/listings/adapters/connectors/multiget_mapper.go:34-36` —
> "UMA linha por item (VariationID = NoVariationID, ADR-13), com os dados de variação
> como filhos em Variations, e não uma linha achatada por variação."

**§2 — New hydration must keep feeding the snapshot observer.** Every batch-hydrated
item is converted to a `ListingSnapshot` and handed to the same `SnapshotObserver` the
old single-item `ReadListing` path feeds, before the item is mapped to its canonical
`listings` row — a batch path that produced correct listing rows but stopped calling the
observer would still be a regression under this rule.
> `apps/server_core/internal/modules/listings/adapters/connectors/backfill.go:63` — "ADR-13
> non-regression: BEFORE mapping to canonical listing rows, every successfully-hydrated
> item is ALSO converted to a connectorsdomain.ListingSnapshot... and handed to the SAME
> SnapshotObserver the existing single-item ReadListing path feeds."

**§3 — Observer failure fails the whole hydration call.** A batch that hydrated
correctly but could not report to the observer does not report a partial success; the
same fail-closed contract as the single-item path applies, because a pull that silently
skipped the linker would report a complete catalog while the matcher stayed blind to
part of it.
> `apps/server_core/internal/modules/listings/adapters/connectors/backfill.go:70-72` (same
> comment block as §2) — "Observer failure fails the whole HydrateBatch call, mirroring
> Source.observe's own fail-closed contract."

**§4 — The rule has a named must-fail test.** `TestMultigetHydratorFeedsObserverBeforeMappingRows`
exists specifically to catch the regression this decision forbids: every hydrated item,
including its per-variation EAN/SKU, must reach the observer exactly once.
> `apps/server_core/internal/modules/listings/adapters/connectors/backfill_test.go:75` —
> "TestMultigetHydratorFeedsObserverBeforeMappingRows is the ADR-13 must-fail proof
> (required test 6)."

**§5 — Content parity, not just count parity, with the old path.** It is not enough for
the new path to call the observer the same number of times; the values it sends
(top-level SellerSKU/EAN derived from the item's own fields) must match what the old
single-item path derived, including for items with no variations.
> `apps/server_core/internal/modules/listings/adapters/connectors/multiget_mapper_test.go:201-203` —
> "TestMapMultigetItemToListingSnapshotFeedsObserverHonestly proves the ADR-13
> observer-feed shape now has CONTENT parity with the old single-item ReadListing path...
> not just count parity."

**§6 — A missing typed field is not license to leave the observer seam honest-empty.**
When the multiget DTO lacked a typed field for an item's own SellerSKU/attributes, the
fix was to type the field, not to accept a degraded, empty snapshot for the sake of
shipping the batch path.
> `apps/server_core/internal/modules/connectors/adapters/mercado_livre/items_multiget_reader.go:67` —
> "instead of leaving that ADR-13 seam honest-empty for lack of a typed field."

**§7 — `listings` keeps its existing primary key; only the child table is new.** The
`listing_variations` table is purely additive; the `listings` table's own PK, including
the `'-'` sentinel it uses for items with no variation, does not change under this
decision.
> `.mnfs/MIS-007-ml-sync/mission.md:242` — "ADR-13 `listing_variations` aditiva; PK de
> `listings` NÃO muda... `listings` mantém PK com sentinela `'-'` (0036)."

**§8 — PK consolidation is explicitly out of scope.** Removing the `'-'` sentinel from
`listings` is named as future, separate-mission work; this decision forbids doing it
inside MIS-007.
> `.mnfs/MIS-007-ml-sync/research/listings-sync-interface-contract.md:146` —
> "Consolidação de PK (tirar sentinela `'-'`) = missão futura nomeada, NUNCA nesta
> (ADR-13)."

**§9 — Non-regression is measured, not asserted.** The validation obligation is stated
as a measured comparison — row counts and anchor coverage of the new path against a
pull performed before the change — not a claim that the new path "should" behave the
same.
> `.mnfs/MIS-007-ml-sync/M-04-listings-backfill-ingest/validation-contract.md:88` —
> criterion "Âncoras de snapshots não-regressivas (ADR-13)."

**§10 — A regression against this rule was found and fixed in review, not just
specified.** The M-04 chip's evidence records an actual instance of the observer-feed
regression being caught and closed against this rule, not merely a theoretical risk.
> `.mnfs/MIS-007-ml-sync/M-04-listings-backfill-ingest/_chip-m04/EVIDENCE.md:24` (regression
> found and fixed against the ADR-13 must-fail; see also lines 58, 329 of the same file).

## Amendments

The harvest that produced this document flagged a live contradiction inside its own
source citations: the `listing_variations` primary key was stated two incompatible ways
under the same "ADR-13" label.

- **5-column form**, including `installation_id`: `.mnfs/MIS-007-ml-sync/mission.md:242-243,
  245`. The mission.md rationale is explicit — `installation_id` is required in the
  tuple "because the real PK of the parent `listings` is `(tenant_id, installation_id,
  provider_listing_id, variation_id)`... and the child sem `installation_id` não endereça
  rows do parent por instalação."
- **4-column form**, omitting `installation_id`: recorded as the pre-fix state of
  `research/listings-sync-interface-contract.md:65` and
  `M-04-listings-backfill-ingest/F-01-listings-ddl/feature.md:31`, flagged in
  `planning-reviews/p7-seat5-doublepass-r02.md:27-32` and
  `p7-claude-readiness-r02.md:39-42`.
- The corpus **self-corrected**: `planning-reviews/p7-seat-star2-r03.md:29-32` records the
  repair to the uniform 5-column tuple, and the interface contract as it stands today
  already carries the corrected form: `.mnfs/MIS-007-ml-sync/research/listings-sync-interface-contract.md:65-70`
  states the PK as `(tenant_id, installation_id, provider, provider_listing_id,
  variation_id)`.

**This document verifies the resolution directly against the migration that ships it,
rather than trusting either side of the citation record.** The actual `listing_variations`
DDL declares:

> `apps/server_core/migrations/0091_listing_variations.sql:31` — `PRIMARY KEY (tenant_id,
> installation_id, provider, provider_listing_id, variation_id)`.

That is **5 columns**, matching the corrected form. The migration's own header comment
states the reasoning verbatim:

> `apps/server_core/migrations/0091_listing_variations.sql:3-8` — "PK is the full parent
> tuple plus variation_id: (tenant_id, installation_id, provider, provider_listing_id,
> variation_id). installation_id is required in the tuple because the real PK of the
> parent `listings` table is (tenant_id, installation_id, provider_listing_id,
> variation_id) — 0036_listings.sql:2-31 — a child key without installation_id could not
> address parent rows per installation."

The parent table's own PK was independently verified against its migration:

> `apps/server_core/migrations/0036_listings.sql:21` — `PRIMARY KEY (tenant_id,
> installation_id, provider_listing_id, variation_id)` (4 columns — `provider` is a
> column on `listings` but is not part of its own PK; `provider` enters the child's PK
> only because the child needs to disambiguate its own rows, per the migration comment
> above).

**Resolution: the 5-column form is correct and is what shipped.** The 4-column citations
were a pre-fix draft state, not a surviving alternate design; §7 and §8 above are written
against the 5-column tuple accordingly.

## Rationale

The observer-feed rule exists because "the read model looks right" and "the matcher is
still fed" are two different claims that a rewritten ingestion path can satisfy or
violate independently — a regression in the second is invisible to every test that only
checks the first. Requiring the new path to call the same observer, with the same
per-item content, before it even finishes mapping its own rows, converts a silent
degradation into a build-time and test-time obligation.

The one-row-per-item rule exists because the alternative — a flattened row per variation
— would multiply `listings` rows against the parent PK sentinel invariant that migrations
0022 and 0025 already depend on, forcing a destructive PK migration as a side effect of
what should be a purely additive throughput change.

## Consequences

- Every future rewrite of listings ingestion inherits the same non-regression
  obligation: it must be measured against the snapshot observer's anchor coverage, not
  just against `listings` row shape.
- `listing_variations` carries `installation_id` and `provider` in its PK even though
  neither is strictly needed to identify a variation within one provider account — the
  cost is a wider key, accepted so the child table can address parent rows per
  installation without a join.
- PK consolidation (removing the `'-'` sentinel from `listings`) remains explicitly
  deferred; this document does not authorize it.

## Unverified claims

- The harvest's anchors at `research/listings-sync-interface-contract.md:36,66,111,172`
  and `M-04-listings-backfill-ingest/milestone.md:57,67` were read in the course of this
  reconstruction and are consistent with the clauses above, but are not separately
  quoted here because they restate content already anchored elsewhere in this document
  (the same PK tuple, the same non-regression must-fail) rather than asserting anything
  additional.
- The pre-fix draft citations (`research/listings-sync-interface-contract.md:65` pre-fix,
  `F-01-listings-ddl/feature.md:31` pre-fix, and the `planning-reviews/p7-seat5-doublepass-r02.md`
  / `p7-claude-readiness-r02.md` / `p7-seat-star2-r03.md` review files) were not
  independently re-read line-by-line for this document; their existence and resolution
  is taken from the citation harvest's own account of the review trail, since the
  question that matters — which PK actually shipped — was settled directly against the
  migration file instead.

## Alternatives Considered

**Flattening variations into separate top-level `listings` rows.** Rejected: it would
require a destructive PK change radiating into two already-dependent modules, in
exchange for no benefit the child-table design does not already provide.

**Skipping the observer call on the batch path and re-syncing anchors separately
later.** Not recorded as a considered alternative in the citations — the decision treats
observer-feed parity as a hard requirement of the ingestion path itself, not a
reconcilable-later gap.
