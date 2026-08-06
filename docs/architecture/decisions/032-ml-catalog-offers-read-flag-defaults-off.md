# ADR-032: ML catalog-offers read gated by a flag that defaults off

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** thin evidence. Only 2 live-code citations of the two-digit token
`ADR-05` assert this decision (Assertion A2 of
`docs/architecture/decisions/_citations/adr-05-twodigit-citations.md`, whose four
mission-scoped meanings for `ADR-05` all collide with each other and with the
pre-existing `005-mercado-livre-first-control-plane.md`, which is about a different
subject entirely). This document is deliberately short: it states only what the two
anchors and the code around them show.

## Context

`ListCatalogOffers` on the Mercado Livre capability adapter reads competitor price
evidence from Mercado Livre's catalog-offers endpoint. The two live citations gate this
specific read behind an environment flag, `MPC_ML_CATALOG_OFFERS_ENABLED`, that defaults
to off.

## Decision

**The `mercado_livre` catalog-offers read route is enabled only when
`MPC_ML_CATALOG_OFFERS_ENABLED` is explicitly set to `"true"`; any other value,
including unset, keeps it disabled.** When disabled, `ListCatalogOffers` returns
`domain.ErrCatalogOffersUnavailable` without calling Mercado Livre.

**§1 — The flag is read-only and independent of the write gate.** It is unrelated to
`MPC_PROVIDER_WRITES_ENABLED`: catalog-offers is a read against Mercado Livre, not a
write, so toggling it never changes whether MPC writes to a provider.
> `apps/server_core/internal/composition/root.go:1001-1006` — "`mlCatalogOffersEnabled`
> gates the mercado_livre catalog-offers READ route (ADR-05: flag defaults OFF). This is
> a read-only provider flag, unrelated to `MPC_PROVIDER_WRITES_ENABLED` and the
> mutations dispatcher."

**§2 — The dev-stack compose file turns it on explicitly, naming the same rule.**
`docker-compose.yml` sets the flag to `"true"` for the backend service, with a comment
citing the same decision and default.
> `docker-compose.yml:30-32` — "ADR-05 read-only flag: ML catalog-offers reads
> (competitor price evidence). Unrelated to MPC_PROVIDER_WRITES_ENABLED (writes stay
> OFF)."

"Read-only provider flag" means concretely: the flag governs a boolean passed into
`CapabilityAdapterConfig.CatalogOffersEnabled`, checked at the top of
`ListCatalogOffers` before any Mercado Livre call is attempted
(`apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:169`).
No write path reads this flag.

## Rationale

The two anchors state a default-off posture for a provider read but do not state, in
live code, *why* the default is off — no anchor gives a cost, rate-limit, or readiness
reason. Recording that omission here rather than inventing a justification is what "thin
evidence" means for this document.

## Consequences

- A tenant with the flag unset gets `ErrCatalogOffersUnavailable` from
  `ListCatalogOffers`, not an error surfaced from a failed Mercado Livre call.
- The dev-stack compose file diverges from the production default by turning the flag on
  explicitly.

## Alternatives Considered

None recorded in the live-code anchors.

## Unverified claims

- **Why the default is off.** Neither anchor states a reason (cost, rate limit,
  readiness). Not asserted as a clause above.
