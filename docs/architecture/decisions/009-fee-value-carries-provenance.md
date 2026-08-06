# ADR-009: Every fee value carries provenance (layer, origem, coletado_em)

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed MIS-007's `channelfees` module from its first
citation but no document was ever written under a stable number. It is reconstructed here
from Assertion A1 of the 26 live citations of `ADR-09` gathered in
`docs/architecture/decisions/_citations/adr-009-citations.md`. Every clause below is
traceable to code, migration, or mission text that already asserts it.

## Context

`channel_fees` is a shared ledger written by two independent producers in different
milestones — listing ingest writes layer-2 rows (commission/freight as published by
Mercado Livre against a live listing), order ingest writes layer-3 rows (the commission
actually realized on a sale) — and read by a third, later consumer, `pricing`, which
resolves one fee value per `fee_kind` to show on `/precos`. Without a decided contract,
each producer would invent its own representation (percent vs. amount, per-unit vs.
total, listing vs. variation), and the pricing consumer would have no way to tell a
resolved-from-ledger value from a fallback default.

The specific failure this decision prevents is a number on screen that looks authoritative
but cannot be traced to where it came from or when it was fetched — the same defect
`ADR-017` names for absence, applied here to attribution. A fee value with no visible
origin is indistinguishable from one the frontend invented.

## Decision

**Every fee value the `channelfees` module resolves is returned together with its layer,
its `origem`, and `coletado_em` (our fetch clock). A consumer that renders the number
without also rendering its provenance fails the milestone's acceptance criteria.**

**§1 — Provenance travels with the value, not beside it.** `ResolveListingFees` returns,
per `fee_kind`, a single struct carrying `Value`, `ValueType`, `Currency`, `Layer`,
`Detail`, `Origem`, and `ColetadoEm` together — there is no code path that returns the
value without the fields that explain it.
> `apps/server_core/internal/modules/channelfees/domain/fee.go:78-87` — `ResolvedFee`
> struct: `Found`, `Layer`, `ValueType`, `Value`, `Currency`, `Detail`, `Origem`,
> `ColetadoEm` declared as one return type.
> `.mnfs/MIS-007-ml-sync/research/channel-fees-interface-contract.md:71-73` —
> "`ResolveListingFees` retorna por fee_kind: `{value, value_type, currency, layer,
> detail, origem, coletado_em}` — proveniência SEMPRE junto do número; consumidor que
> exibe número sem proveniência reprova milestone (ADR-09)."

**§2 — A screen without provenance is a blocking failure, not a style nit.** The
milestone's own validation contract names the defect this decision exists to catch and
treats it as the thing ADR-09 was written to kill.
> `.mnfs/MIS-007-ml-sync/M-07-pricing-fee-read/validation-contract.md:142` — "Blocking
> failure: número de tarifa sem origem na tela (a doença que ADR-09 mata)."

**§3 — `origem` is a closed, named enum, not a free-text label.** The source of a fee
value is one of five declared constants; `coletado_em` (our fetch clock) is kept distinct
from `source_time` (the provider's own timestamp, nullable, never fabricated when
absent — the same non-fabrication rule as `ADR-017`).
> `apps/server_core/internal/modules/channelfees/domain/fee.go:44-49` — `OrigemAPIListingPrices`,
> `OrigemAPIShippingOptions`, `OrigemAPIOrder`, `OrigemAPIShipment`, `OrigemConfig`.
> `apps/server_core/internal/modules/channelfees/domain/fee.go:68-69` — `ColetadoEm time.Time
> // our fetch clock — never equated with SourceTime for convenience` and `SourceTime *time.Time`.

**§4 — The resolution ladder is layer 2 → layer 1 → absent; layer 3 never enters it.**
Layer 3 rows (realized, per-order fees) are written to the same table but are structurally
excluded from what `pricing` resolves — the ladder answers "what does this listing publish
now", never "what did one order actually pay".
> `apps/server_core/internal/modules/channelfees/domain/fee.go:92` — "Layer 3 never
> participates in either ladder (ADR-09 — layer 3 rows may exist, written by other
> producers, but this reader ignores them)."
> `apps/server_core/internal/modules/channelfees/adapters/postgres/reader.go:47` —
> "Layer 3 is never queried by either ladder (ADR-09)."

**§5 — Fee data lives only in `channel_fees`; no dual-write elsewhere.** Columns that
would duplicate commission or free-shipping figures on the `listings` table were
deliberately left out of that table's own migration, naming this decision as the reason.
> `apps/server_core/migrations/0090_listings_e3_fields_status_relax.sql:19-22` —
> "Deliberately EXCLUDED (IC-07 decision, planning design §7 item 6): commission_amount,
> commission_pct, free_shipping_cost — fee data lives only in channel_fees (IC-01,
> ADR-09 provenance); dual-write here would be a drift seam, not a feature."

## Amendments

One candidate reading of ADR-09 was ratified and then abandoned without being struck from
the citing text: an enumerated, standalone "fee ledger" concept (seeded 16/22 rows,
`registry_test.go:90-103`) plus a live "degrau-3" (step-3) resolver
(`root.go:845-851`, `tarifflive/resolver.go:43-69`) that looked up `channel_fees` outside
the ladder described above. MIS-007's own mission log records this as dead — no milestone
claims it — and the live degrau-3 path is retired in favor of the ladder in §4.
`baseline_commission_percent: 0.16` (`auth_adapter.go:42-48`) survives as unrelated,
published provider catalog metadata with no call site in pricing; it is explicitly carved
out of this decision rather than folded into `origem='config'`.
> `.mnfs/MIS-007-ml-sync/mission.md:205-212` — "ADR-09 Ledger de fee enumerado. JÁ MORTO
> (nenhum milestone reivindica) ... EMENDA (auditoria P5 r02 N-2): `baseline_commission_percent:
> 0.16` ... é METADATA de catálogo do provider ... fica INTOCADA (a disjunção anterior
> 'morre ou vira row `origem='config'`' caiu — premissa de fallback silencioso era falsa)."

## Rationale

A fee value with no attached source is a claim nobody can check. `pricing` composes
values from at least three producers writing at different cadences (listing sync, order
sync, static config fallback); once a number reaches the screen, the layer it came from
and the moment it was fetched are the only way to tell "this is what Mercado Livre
publishes today" from "this is our stale config default" from "this is what one order
actually paid". Stripping provenance at any hop makes every later hop equally
untrustworthy, which is why the rule is enforced at the type that leaves the module
(`ResolvedFee`), not as a rendering convention in the frontend.

## Consequences

- Every `channelfees` reader function returns a composite type, never a bare value — this
  is more ceremony per call site and is accepted deliberately.
- The `/precos` screen must render `origem` (and a warning when `origem=config`) next to
  every fee figure; a passing test that only asserts the number, not its provenance
  badge, is not sufficient evidence the decision is honored.
- Layer 1 (category-level fee) has no producer in this mission; the ladder code path for
  it is a pinned placeholder, not dead code to delete — its order is contractually
  decided even though nothing populates it yet (`reader.go:41-45`).
- The abandoned "enumerated fee ledger" reading (see Amendments) is retired terminology;
  a citation still referencing it should be read against the ladder in §4, not against a
  separate ledger concept.

## Alternatives Considered

**A single flattened fee amount per listing, without layer/origin.** Rejected: this is
exactly the shape that produced the abandoned degrau-3 lookup and the enumerated-ledger
reading in Amendments — a value with no way to distinguish a live listing fee from a
realized order fee from a config default.

**Provenance as a side-channel (separate endpoint or log line) instead of on the value.**
Rejected: provenance that is not carried with the value is provenance that is easy to
drop at the next hop — the same reasoning `ADR-017` applies to a "known" flag kept next to
a zero.

## Unverified claims

None. All anchors cited for Assertion A1 were read and confirmed to match verbatim (or in
substance, for the mission-log Amendments text) at the cited line.
