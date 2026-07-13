# F-01 Canonical Product Identity — Specification

## Scope

Establish the typed contract that later catalog cutover work will populate from
Sankhya. `internal_product_id` is a positive integer and is exactly `CODPROD`.
EAN, manufacturer reference, seller SKU, and provider listing identifiers remain
separate optional fields. This feature does not map legacy catalog strings to
CODPROD or change the active legacy reader.

## Contract

- `catalog/domain` provides a positive `InternalProductID` value object and a
  `CanonicalProduct` DTO.
- `internal_read/domain.ProductCandidate` carries the same typed positive
  identity rather than an unqualified integer.
- A source fact has nullable numeric value, server-owned
  `current|stale|unknown|conflict` quality, source, optional observed time, and
  a nonblank reason for `unknown` or `stale`.
- The public OpenAPI `CatalogProduct` and SDK model expose this canonical shape
  in lockstep. Numeric facts are nullable so known zero remains distinguishable
  from unknown.

## Rejections

- `CODPROD <= 0` is rejected as `invalid_identity`.
- EAN/reference/seller SKU alone cannot construct `InternalProductID`.
- Conflicting eligible identities are represented as `identity_conflict`; no
  selection is made by this feature.
- A missing numeric fact must serialize as `null`, `quality=unknown`, and a
  nonblank `quality_reason`.

## Boundaries

No Oracle, provider, MSDB, composition, database, auth, or web changes are in
scope. F-02 owns wiring the active catalog read path to the canonical contract;
F-03 owns deterministic legacy compatibility/cutover. Existing legacy catalog
routes are therefore not reinterpreted as CODPROD in F-01.

## Acceptance

Deterministic Go and SDK tests prove positive identity validation, identifier
separation, known-zero versus unknown serialization, and exact OpenAPI/SDK
field/type parity.

### F01-AC01 — Canonical identity

Only a positive CODPROD constructs `InternalProductID`; separate identifiers
cannot populate it.

### F01-AC02 — Honest source facts

Unknown numeric facts remain null with a reason, while an observed numeric zero
is retained as zero.

### F01-AC03 — Atomic public contract

The canonical OpenAPI schema and SDK interfaces contain the same identity,
identifier, source, quality, and nullable numeric-fact fields.
