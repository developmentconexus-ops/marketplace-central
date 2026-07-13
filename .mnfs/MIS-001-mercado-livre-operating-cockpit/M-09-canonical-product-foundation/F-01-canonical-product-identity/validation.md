# F-01 Canonical Product Identity — Validation

## Outcome

Implemented a canonical contract without mapping any legacy catalog string to
CODPROD. `InternalProductID` accepts only positive values; source facts preserve
known zero separately from unknown `null`; and OpenAPI/SDK definitions are
updated together.

## Context

- `context.json` compiled and validated at base
  `d69490a0aef4ba2c3198e3e39fb20477e811e276`.
- Context risk: L2 (`api-sdk` shared seam).

## Proof

- `GOCACHE=C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache; go test ./internal/modules/catalog/domain ./internal/modules/internal_read/domain` from `apps/server_core` — PASS.
- `npm test --workspace @marketplace-central/sdk-runtime` — PASS.
- OpenAPI/SDK parity assertion for canonical identity, identifiers, nullable
  source facts, and quality reason — PASS (`PASS openapi-sdk-parity`).
- `context-validate -ContextPath .../F-01-canonical-product-identity/context.json -RequireCurrentBase` — PASS.

## Changed Paths

- `apps/server_core/internal/modules/catalog/domain/canonical_product.go`
- `apps/server_core/internal/modules/catalog/domain/canonical_product_test.go`
- `apps/server_core/internal/modules/internal_read/domain/canonical_identity.go`
- `apps/server_core/internal/modules/internal_read/domain/internal_product.go`
- `apps/server_core/internal/modules/internal_read/domain/contract_test.go`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.ts`
- `packages/sdk-runtime/src/index.test.ts`
- This feature's `spec.md`, `plan.md`, `context.json`, and `validation.md`.

## Limitations And Handoff

The active legacy catalog read path remains unchanged by design: neither its
string identity nor an EAN/reference/seller SKU is coerced to CODPROD. F-02
must populate the canonical contract from Oracle/internal_read; F-03 owns any
deterministic compatibility removal. No Oracle/provider calls, writes, auth, or
dependency installation occurred.

## Required Nullable-Field Correction

- Removed Go `omitempty` from `quality_reason`; canonical nullable identifiers
  already used non-omitting JSON tags at correction base
  `954b88c7fc97fe3063ccec8a68f12caf12732b55`.
- OpenAPI now requires `ean`, `manufacturer_reference`, `seller_sku`, and
  `quality_reason` while retaining `nullable: true`; SDK fields remain required
  `T | null`.
- Go serialization proof covers absent identifiers and current
  `quality_reason` emitted as explicit `null`; unknown/stale blank reasons are
  rejected.
- `GOCACHE=C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache; go test ./internal/modules/catalog/domain ./internal/modules/catalog/adapters/internalread`
  from `apps/server_core` — PASS.
- `npm test --workspace @marketplace-central/sdk-runtime` — PASS (40/40),
  including static OpenAPI/SDK required-nullable parity.

## M-09-CORR-03 Final C01 Correction

- Product-link generation now filters, compares, deduplicates, constructs
  candidate IDs, and persists identity only from a non-nil positive canonical
  `ProductCandidate.InternalProductID`. Legacy `ProductID` remains metadata and
  is never promoted.
- Unit proof covers nil, zero, and negative canonical IDs paired with positive
  legacy IDs; each produces only an unresolved candidate with no persistable
  canonical identity. Conflict, deduplication, and candidate-ID stability use
  deliberately divergent legacy metadata and positive canonical IDs.
- `CanonicalCatalogProduct.required` now includes nullable `brand_name` and
  `product_group_name`, matching the existing required SDK `string | null`
  fields.
- `GOCACHE=C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache; go test ./internal/modules/product_links/... -count=1`
  from `apps/server_core` — PASS.
- `GOCACHE=C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache; go test -tags integration ./internal/modules/product_links/application -run '^$' -count=1`
  from `apps/server_core` — PASS (integration fixture compiles with canonical
  ID; it was not executed because this correction forbids database writes).
- `npm test --workspace @marketplace-central/sdk-runtime -- --run` — PASS
  (40/40).
- No Oracle, provider, network, or database action occurred. Independent
  fixed-SHA review and proportional QA remain Milestone-owned.
