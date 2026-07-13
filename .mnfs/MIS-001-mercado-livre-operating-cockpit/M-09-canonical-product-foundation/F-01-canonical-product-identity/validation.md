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
