# ADR-025: Raw payload capture is selective — never PII

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed order ingest under the two-digit citation
`ADR-03` inside MIS-007, and was never given a document of its own. The two-digit number
collided with three unrelated decisions from other missions and with the pre-existing
`003-integration-spec-split-and-sequencing.md`, which is about a different subject
(sequencing integration delivery specs). This document reconstructs the MIS-007
raw-selective / PII-never rule under its own global number, from the 10 live-code
citations harvested at
`docs/architecture/decisions/_citations/adr-03-twodigit-citations.md` (Assertion A1).
Every clause below is traceable to code that already asserts it. Nothing here is new
policy.

## Context

The `/orders/{id}` response from Mercado Livre carries `buyer.billing_info` — the buyer's
fiscal identity (CPF/CNPJ and related fields). Order ingest also wants a raw capture of
the provider payload alongside the typed DTO, so a single GET can serve both without a
second round trip. Holding the raw payload verbatim would mean holding PII in memory (and
potentially at rest) for a resource where nothing downstream needs it. Shipment payloads
carry a related but distinct class of PII: the delivery receiver's name and address.

## Decision

**A raw provider payload may be captured for ingest, but only after the specific fields
that carry PII or fiscal data for that resource are stripped. The exemption is
per-resource, not blanket: `orders`/`order_shipments` get no raw column in the database at
all — the redacted raw payload is a request-scoped, in-memory value only, used solely so
one GET can serve both the typed DTO and the redaction check. `listings` is the only table
with a persisted `raw jsonb` column.**

**§1 — `buyer.billing_info` is stripped before the raw order payload leaves the adapter.**
`redactBillingInfo` operates on `map[string]json.RawMessage`, deleting `billing_info`
inside `buyer` only; every other field (order_items, payments, large numeric ids) is
re-emitted byte-for-byte from its original `RawMessage`, so there is no float64
round-trip or precision loss anywhere else in the payload.
> `apps/server_core/internal/modules/connectors/adapters/mercado_livre/order_ingest_reader.go:164-171`
> — "redactBillingInfo returns the order payload with `buyer.billing_info` removed
> (ADR-03: raw order capture may be held in memory MINUS billing_info)."

**§2 — Redaction failure drops the raw capture, it never falls back to unredacted raw.**
If the `buyer` object cannot be safely re-decoded, or if `redactBillingInfo` itself
errors, the raw capture is discarded entirely (`rawOrder = nil`) rather than risking a
leak; the typed DTO is unaffected and remains the load-bearing data.
> `apps/server_core/internal/modules/connectors/adapters/mercado_livre/order_ingest_reader.go:121-127`
> — "Cannot safely guarantee the billing_info redaction succeeded — drop the raw capture
> entirely rather than risk ADR-03 leaking a PII/fiscal fragment."

**§3 — The redacted raw payload is a transport-layer value, not a persisted column.**
`domain.OrderDetail.RawOrder` is documented as held "for the ingest layer only"; neither
`orders_marketplace_orders` (migration 0089) nor `order_shipments` (migration 0088)
carries a raw column, and the ingest code (F-02) must not persist this field.
> `apps/server_core/internal/modules/connectors/domain/order_detail.go:59-64` — "RawOrder is
> the /orders/{id} response body MINUS buyer.billing_info (ADR-03, feature.md
> Negative/Constraints): held in memory for the ingest layer only. ... F-02 must NOT
> persist this field."

**§4 — Shipment gets no raw column at all, by design, not by omission.** A shipment
payload carries a different class of PII (receiver name, street address, ZIP) than
`billing_info` redaction covers, so rather than write a second redaction function,
`order_shipments` is typed-columns-only. `raw jsonb` stays scoped to `listings`.
> `apps/server_core/migrations/0088_order_shipments.sql:1-8` — "Typed shipment facts only.
> No `raw` column, by design (IC-03, P7 r01 B-7, ADR-03 amended): a shipment payload
> carries PII of delivery (receiver name, street address, ZIP) ... `raw jsonb` stays
> scoped to `listings` only."

**§5 — Buyer-fiscal data is typed-only; the raw exemption does not extend to it either.**
The dedicated buyer-fiscal fetch persists exactly the typed fields the comprador-fiscal
drawer renders, and explicitly forbids a raw billing_info capture on that path.
> `apps/server_core/migrations/0089_orders_marketplace_orders_sync_fields.sql:29-30` — "Buyer
> fiscal, typed only (ADR-03/R-6: raw billing_info payload is forbidden)."

## Rationale

A raw capture exists for one operational reason here: to serve a typed DTO and a
diagnostic payload from a single provider round trip. That reason does not require, and
does not justify, holding a buyer's fiscal identity or a recipient's home address in a
form nothing in the ingest path reads back out. Stripping the specific known-sensitive
fields at the point of capture — rather than at some later serialization boundary — means
there is no window in which the unredacted payload exists past the adapter call that
fetched it, and a redaction failure fails toward dropping the capture, not toward leaking
it.

## Consequences

- Adding a new provider resource to the raw-capture pattern requires naming which fields
  of *that* resource are PII/fiscal and writing a redaction step for them; there is no
  generic redaction utility to reuse; the shipment case (§4) shows the fields differ per
  resource.
- `orders_marketplace_orders` and `order_shipments` have no raw column to backfill,
  audit, or migrate — any future need for a persisted raw order/shipment capture is a new
  decision, not an extension of an existing column.
- The stripped `buyer.billing_info` is unrecoverable from the raw capture once redacted;
  only the typed buyer-fiscal fetch (§5) is a source of that fact going forward.

## Alternatives Considered

**Persist the raw payload and redact at read time.** Rejected as implemented: this would
require holding the unredacted payload at rest between write and read, during which any
bug in the read-time redaction, or any direct query against the table, exposes PII. The
chosen design never lets the unredacted payload reach persistence.

**One shared redaction function covering both `billing_info` and shipment PII fields.**
Not what the code does: shipment does not get a redaction function at all — it gets no
raw column, because its PII surface (receiver name/address) differs from the order's
(fiscal identity) and the simplest correct fix per §4 was to drop the raw capture for that
resource entirely rather than generalize the redaction.

## Unverified claims

None — every clause above matches a verified anchor in code.
