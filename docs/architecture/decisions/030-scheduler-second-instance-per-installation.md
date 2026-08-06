# ADR-030: A second Scheduler instance is created per Mercado Livre installation

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed MIS-007's listings and orders sync
composition but was only ever cited by its local two-digit label `ADR-08`, a number
that collides with an unrelated MIS-004 decision (zero writes to Mercado Livre in the
demo) and with a pre-existing three-digit document about production deploy topology.
It is reconstructed here from the 2 live-code citations harvested at
`docs/architecture/decisions/_citations/adr-08-twodigit-citations.md`, Assertion A1.
The renumbering to ADR-030 is recorded in
`docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md:89`.

## Context

`sync/composition.NewProductsScheduler` established a pattern: one
`sync/application.Scheduler` wired tenant-wide, driving the products sync job. When
listings backfill/sweep and orders backfill/incremental needed their own scheduled
jobs, the tenant-wide pattern did not carry over as-is, because a listings or orders
sync run is not tenant-wide — it is scoped to one specific Mercado Livre
`InstallationAccount` — and a tenant can have more than one ML installation.

## Decision

**A separate `sync/application.Scheduler` instance is constructed per active
Mercado Livre installation for listings backfill/sweep and for orders
backfill/incremental, instead of reusing one tenant-wide scheduler instance.**

**§1 — The constraint comes from the Scheduler's own construction contract,
not from a stylistic preference.** `sync/application.Scheduler` is bound to a single
`installationID` at construction, and `RegisterJob` allows only one registration per
entity per instance. A single shared scheduler instance therefore cannot serve N
installations of the same entity.
> `apps/server_core/internal/modules/listings/composition/scheduler.go:9-13` —
> "sync/application.Scheduler is bound to a single installationID at construction
> and RegisterJob allows only one registration per entity per instance, so a single
> shared scheduler cannot serve N installations of the same entity — this package
> builds one Scheduler per active ML installation instead."

**§2 — The listings scheduler package is explicit that it cannot reuse the products
pattern verbatim, for the same reason.** Products data is tenant-wide (one active ERP
source per tenant); a listings backfill run is scoped to one installation.
> `apps/server_core/internal/modules/listings/composition/scheduler.go:1-8` —
> "the ADR-08 \"segunda instância\" of the sync scheduler pattern
> sync/composition.NewProductsScheduler established. It cannot reuse that pattern
> verbatim: products data is tenant-wide ..., but a listings backfill run is scoped
> to ONE specific Mercado Livre InstallationAccount, and a tenant can have several
> ML installations."

**§3 — The wiring site names the same rule at the point where both schedulers are
started.** `root.go` starts the listings schedulers as one instance per active
installation, and comments the reason inline.
> `apps/server_core/internal/composition/root.go:837-840` — "Daily listings
> backfill/sweep scheduler (ADR-08 second Scheduler instance, mirroring
> synccomposition.NewProductsScheduler): one instance per active Mercado Livre
> installation, since (unlike products) a listings backfill run is scoped to one
> specific InstallationAccount."

The orders scheduler package follows the identical reasoning and cites the same
constraint against the listings package rather than restating it independently:
"o fan-out por instalação (e não um scheduler compartilhado) é imposto pelo próprio
sync ... É a mesma razão que listings/composition documenta"
(`apps/server_core/internal/modules/orders/composition/scheduler.go:6-9`). `root.go`
wires `orderscomposition.NewOrdersSchedulers` immediately after the listings
schedulers, on a 15-minute cadence versus listings' 24-hour cadence
(`apps/server_core/internal/composition/root.go:844-856`).

## Rationale

Reusing a single tenant-wide scheduler instance for a per-installation job would
require either binding that one instance to a single installation (silently
dropping every other installation's runs) or fighting the `RegisterJob`
one-registration-per-entity limit by registering the same entity multiple times
against one instance, which the Scheduler's contract does not support. Building one
instance per installation instead keeps each installation's job independent: one
installation's failure, cadence, or backoff does not touch another's.

## Consequences

- The number of running Scheduler instances scales with the number of active ML
  installations per tenant, not with the number of entities (products, listings,
  orders) alone.
- Reusing the Scheduler for job scheduling (instead of a bespoke ticker) is what
  makes a stuck installation observable: failures record into `sync_state` and
  surface through `/sync/health`, which a parallel ad-hoc ticker would not provide
  for free (`apps/server_core/internal/modules/orders/composition/scheduler.go:11-15`).
- Any future entity whose sync scope is per-installation (rather than tenant-wide)
  needs to follow the same one-instance-per-installation pattern rather than the
  tenant-wide `NewProductsScheduler` pattern.

## Alternatives Considered

**One tenant-wide Scheduler instance for listings/orders, mirroring products.**
Rejected: the Scheduler's construction contract binds one instance to one
`installationID` and allows only one registration per entity per instance, so a
single instance cannot register and run the same job across a tenant's multiple ML
installations — a second installation's job would either fail to register or would
silently never run.

**A bespoke ticker per installation instead of reusing `sync/application.Scheduler`.**
Rejected by the orders package's own reasoning: a parallel ticker would cost the same
amount of code as reusing the Scheduler while losing the failure visibility
(`sync_state.consecutive_failures`, `/sync/health`) that reuse gets for free.

## Unverified claims

None. Both live-code anchors in the harvest were read and confirmed to state the
clauses above; the orders-scheduler cross-reference and the cadence values are
additional context read directly from `orders/composition/scheduler.go` and
`root.go` and are reported as observed, not asserted as part of the ADR-08 citation
itself.
