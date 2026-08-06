# ADR-022: Provider-write SKU invariant — SELLER_SKU must equal linked CODPROD

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed MIS-003's mutation-envelope writes but no
document was ever written under its original number (ADR-16). It is reconstructed here
from the 8 live citations of `ADR-16`, harvested at
`docs/architecture/decisions/_citations/adr-016-citations.md`. Every clause below is
traceable to code or to the mission's binding contracts that already assert it.

## Context

Provider writes that edit or create a Mercado Livre listing (`listing_edit`,
`listing_create`) carry a `SELLER_SKU` attribute. The listing is linked internally to an
ERP product by `CODPROD`. If a write is allowed to set `SELLER_SKU` to any value other
than the linked `CODPROD`, the live listing on the marketplace becomes mislinked from the
ERP product it is supposed to represent — a defect that only surfaces later, on the
provider's side, once it is already applied.

## Decision

**A `listing_edit` or `listing_create` write item is only valid if its `SELLER_SKU`
equals the internally linked `CODPROD`. A mismatch is rejected before the write reaches
the provider adapter, with the stable failure code `sku_invariant_violation`. It is never
applied.**

**§1 — The invariant.** `SELLER_SKU` on a `listing_edit` intent must equal the linked
`CODPROD`; a mismatch produces a typed failure and the write is not dispatched to any
adapter.
> `apps/server_core/internal/modules/mutations/application/writer.go:82-90` — `itemGate` calls `sellerSKUMismatch(item.After)` for `ProtocolTypeListingEdit` and, on mismatch, returns `failed(domain.FailureCodeSKUInvariantViolation, "SELLER_SKU deve ser igual ao CODPROD vinculado.")`.

**§2 — The comparison is `SELLER_SKU` attribute value vs. `product_id`.** The mismatch
check reads the item's `product_id` as the linked product and looks for a `SELLER_SKU`
attribute in the intent's attribute list; any other value is a mismatch.
> `apps/server_core/internal/modules/mutations/application/writer.go:175-190` (`sellerSKUMismatch`) — compares `strings.TrimSpace(attribute.ValueName) != codprod` where `codprod` is `string(intent.ProductID)`.

**§3 — The failure code is stable and typed.** `sku_invariant_violation` is a named
constant in the failure-code taxonomy, not an ad-hoc string.
> `apps/server_core/internal/modules/mutations/domain/failure.go:13` — `FailureCodeSKUInvariantViolation FailureCode = "sku_invariant_violation"`.

**§4 — Rejection happens before any adapter call; the provider is never invoked.** The
gate runs in `itemGate`, which is checked before `ApplyItem` dispatches to a provider
writer port. A test asserts the adapter's call counter stays at zero on rejection.
> `apps/server_core/internal/modules/mutations/application/writer_test.go:60-67` (`TestWriterRouterRejectsSKUInvariantBeforeAdapter`) — asserts `out.Failure.Code == domain.FailureCodeSKUInvariantViolation && listing.calls == 0`.

**§5 — The code has a pt-BR translation in the shared failure-copy table.** The failure
surfaces to the operator with dedicated copy, not a raw code.
> `packages/web-query/src/failureCopy.ts:9,24` — `"sku_invariant_violation"` maps to `"Violação de invariante de SKU."`.

## Rationale

Rejecting the mismatch before the provider call is the only point at which the invariant
can be enforced cheaply: once a mislinked `SELLER_SKU` reaches Mercado Livre, correcting
it requires a second write and a window in which the live listing and the ERP product
disagree. Keeping the failure code stable and named (rather than a generic validation
error) lets the operator UI and any downstream audit distinguish this specific defect
from other write failures.

## Consequences

- Every `listing_edit`/`listing_create` write item must carry an internally consistent
  `product_id`/`SELLER_SKU` pair before it reaches the writer router; upstream code that
  builds these intents must not allow the two to drift.
- Any new provider-write path added later must run its item through the same
  `itemGate`-style check, or it silently reopens the mislinking defect this rule exists
  to prevent.
- The mission's own citations describe this rejection as happening "at preview" or
  "pre-apply" (see Unverified claims): the verified code path enforces it in the item
  gate that runs immediately before dispatch to the provider adapter, inside the same
  write-application flow the in-process poller drives — not inside the separate
  `Preview()` HTTP handler (`apps/server_core/internal/modules/mutations/application/preview.go`),
  which builds the before/after snapshot but does not call `sellerSKUMismatch` or
  `itemGate`. The substance of the citations — rejected before the provider is ever
  called, never applied — is confirmed; the specific claim that the *preview* HTTP
  response itself carries the failure code is not.

## Alternatives Considered

**Validate only after the provider write fails.** Rejected: this lets a mislinked SKU
reach the marketplace at least once, and requires a corrective second write instead of
preventing the defect at zero cost before dispatch.

**A generic validation failure code instead of a dedicated one.** Rejected: an unnamed
validation error is not distinguishable from any other rejected write in the operator UI
or in audit review, which defeats the purpose of a typed failure-code taxonomy.

## Unverified claims

- The citations (`mutation-envelope-interface-contract.md:162,212`,
  `M-03-mutation-envelope-writes/validation-contract.md:104-112`) state the rejection
  happens "at preview" / "at validation pre-apply," and one describes the precondition as
  `link.state=resolved` rather than a direct `SELLER_SKU == CODPROD` string comparison.
  The code found (`writer.go:82-90,175-190`) enforces the invariant in the write-apply
  gate (`itemGate`, called from the poller-driven `Apply`/`ApplyItem` path), not inside
  `Service.Preview` (`preview.go`), and compares `SELLER_SKU` directly against
  `product_id` rather than checking a link-resolution state. No occurrence of a
  SKU-invariant check inside `preview.go` was found by searching the `mutations` module.
