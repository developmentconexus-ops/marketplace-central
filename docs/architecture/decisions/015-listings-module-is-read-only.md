# ADR-015: Canonical `listings` module is read-only

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed MIS-003's listings read spine and was enforced
across every feature and milestone that touched it, but no document was ever written
under this number — MIS-003's own artifacts flag this as a standing documentation gap
(see below). It is reconstructed here from the 15 live citations of `ADR-12` that assert
this meaning, harvested at
`docs/architecture/decisions/_citations/adr-012-citations.md` (Assertion A1). The same
`ADR-12` label was also used, unrelated, for two other decisions in two other missions —
see the registry at `docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md`.

## Context

Before this decision, listing data (Mercado Livre-shaped: `MLB` ids, provider-specific
fields) leaked into `product_links`, a module that is supposed to be provider-agnostic.
The mission needed a canonical, marketplace-agnostic read model for listings — serving
the "Anúncios" screen, the "Anúncios vinculados" tab, dashboard counters, and the
selection input for a separate mutation contract (IC-03) — without turning that read
model into a second place where writes could happen, and without letting the provider's
transport shape dictate the domain shape.

## Decision

**The `listings` module is a read model only. It never accepts writes, mutations, or
edits — those belong to a separate mutation contract (IC-03) — and its own state comes
exclusively from provider ingestion via the connectors capability, on a canonical,
provider-agnostic shape.**

### The clauses

**§1 — The module is read-only.** `listings` is a Go module that serves reads; edits and
mutations are IC-03's concern, not this module's.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-01-listings-read-spine/F-01-listings-module-ingestion/feature.md:17`
> — "Canonical marketplace-agnostic listing read model: new `listings` Go module
> (read-only) → HTTP transport → sdk-runtime → web UI."

**§2 — Ingestion is the only way data enters the module, and it runs through the
connectors capability.** Rows are upserted by a full-page pull via the connectors
capability, keyed on `(tenant_id, installation_id, provider_listing_id, variation_id)`.
Rows absent from a completed pull are marked `status=closed`, never deleted.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/listings-read-interface-contract.md:93`
> — "New table `listings` (module `listings`), tenant-scoped, upserted by ingestion keyed
> on `(tenant_id, installation_id, provider_listing_id, variation_id)`. Ingestion is
> full-page pull via connectors capability; rows absent from a completed pull are marked
> `status=closed` (never deleted)."

**§3 — Identity is a canonical composite `listing_id`, opaque to clients.** The id is
built as `installation~provider_listing_id~variation`, with a literal `-` sentinel for a
null variation. It is never parsed on the client side.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/listings-read-interface-contract.md:29`
> — "`listing_id` | string | no | canonical composite
> `installation~provider_listing_id~variation`, literal `-` for null variation (carried
> from MIS-001 M-13)"
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/listings-read-interface-contract.md:164`
> — "`listing_id` composite uses `~` separators, `-` for null variation; opaque to
> clients (never parsed client-side)."

**§4 — Facts the module cannot establish are nullable, never defaulted.** Price,
published quantity, quality score, and sales are all nullable columns; a null cost is
never coerced into "not below margin" — it is surfaced as `unknown` through product
completeness instead.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/listings-read-interface-contract.md:188`
> — "Unknown never becomes zero/default (price, quantity, quality, sales all nullable)."
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/listings-read-interface-contract.md:93`
> — "`below_margin` exception computed at read time from pricing policy min-margin vs
> latest cost fact (null cost → NOT below_margin; it is `unknown`, surfaced via product
> completeness, not a false alarm)."

**§5 — Freshness is `fetched_at`, refreshed manually, and distinct from serve time.**
`fetched_at` records the provider's capture time and only moves on a successful provider
fetch; `as_of` is the separate read-model serve time. The module does not poll on a
schedule inside this mission — refresh is manual.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:111` — "canonical fields,
> composite listing_id, nullable unknowns; ingestion via connectors capability only;
> freshness = `fetched_at` + manual refresh (review 9c)"
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/listings-read-interface-contract.md:165-166`
> — "`as_of` = read-model serve time; `fetched_at` = provider capture time. Distinct, both
> surfaced." ... "`fetched_at` only on successful provider fetch."

## Rationale

Splitting listings into a dedicated, read-only, canonical module keeps two concerns from
colliding: what the provider says about a listing (ingested, canonical, versioned by
`fetched_at`) and what an operator wants to change about it (a separate mutation
contract, IC-03, with its own write path and its own validation). If the read model
accepted writes directly, every consumer of `listings` would have to reason about
in-flight edits racing against the next ingestion pull. Keeping ingestion as the only
writer, gated through the connectors capability, means the module's state is always
explainable as "what the last completed pull said," never "what the last pull said, plus
whatever a screen wrote on top."

## Consequences

- Any feature that wants to change listing data (price override, manual close, edit) must
  go through IC-03's mutation contract, not through the `listings` module or its table
  directly.
- Screens reading `listings` must render `fetched_at` and treat a stale `fetched_at` as a
  freshness signal, not silently show data as current.
- A listing whose provider row disappears from a completed pull is marked `closed`, never
  deleted — history of what was once listed is preserved.
- Extension is additive-only: new `filter` keys and new nullable `Listing` fields are
  allowed; changing the semantics of an existing key is not.
  > `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/listings-read-interface-contract.md:171`
  > — "`filter` extends by new keys only. `Listing` extends by new nullable fields only."
- This decision had no formal ADR document while MIS-003 was running. Two of the
  mission's own artifacts logged this as an unresolved documentation-governance
  escalation, explicitly out of scope for the feature that hit it:
  > `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-01-listings-read-spine/DECISIONS.md:140`
  > — "ADR-12 / ADR-17 have no formal record under docs/architecture/decisions (behavior
  > unambiguous in mission.md). Architecture owner repairs; F-01 proceeds on the
  > mission-fixed behavior."
  > `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-01-listings-read-spine/DISPATCH-LEDGER.md:19`
  > — "**ESCALATION (out-of-scope):** ADR-12/ADR-17 have no formal record under
  > docs/architecture/decisions (behavior unambiguous; architecture owner repairs)."
  This document is that repair.

## Alternatives Considered

**Let `listings` accept writes directly (e.g. manual price edits) alongside ingestion.**
Rejected: a module that is both the ingestion target and a live write target has no way
to say, at read time, whether a field reflects the provider or a local override, without
extra bookkeeping this mission did not build. Keeping edits in a separate mutation
contract (IC-03) avoids that ambiguity entirely.

**Reuse `product_links`' existing shape for listings.** Rejected: that module's fields
were already leaking Mercado Livre-specific shape (`MLB` ids and provider-specific
columns) into what was meant to be a provider-agnostic surface — the exact defect this
module was created to fix.

**Poll on a schedule instead of manual refresh.** Not adopted within this mission's
scope: freshness is `fetched_at` plus a manual trigger; a scheduled poller was not part
of what MIS-003 built for this module.
