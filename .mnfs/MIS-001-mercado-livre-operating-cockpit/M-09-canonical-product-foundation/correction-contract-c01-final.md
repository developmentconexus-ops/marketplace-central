# M-09 Final C01 Correction Contract

```yaml
id: M-09-CORR-03
type: portfolio-correction-contract
status: authorized
owner: M-09 Milestone Orchestrator
parent: M-09
authorized_by: Portfolio Hub
base_sha: 230dc78306d3775894a00b5424238529382cc9b0
attempts_allowed: 1
created: 2026-07-13
```

## Reason For Exception

The ordinary correction budget is exhausted. Fixed-SHA review nevertheless proves
M-09-C02 through M-09-C05 and isolates two mechanical M-09-C01 defects. Portfolio
authorizes one final attempt because the scope is closed, does not change product
intent, and can be proved deterministically without another live-source run.

## Required Corrections

1. Product-link candidate generation must require a non-nil positive
   `ProductCandidate.InternalProductID`. Filtering, equality/conflict detection,
   deduplication, candidate ID construction, and persisted `InternalProductID` must
   use that canonical value. Legacy `ProductCandidate.ProductID` remains source
   compatibility metadata only and can never be promoted to the link identity.
2. `CanonicalCatalogProduct.required` in OpenAPI must include nullable
   `brand_name` and `product_group_name`, matching Go's always-emitted JSON keys and
   the SDK's required `string | null` fields.

## Allowed Paths

- `apps/server_core/internal/modules/product_links/application/generation_service.go`
- `apps/server_core/internal/modules/product_links/application/generation_service_test.go`
- `apps/server_core/internal/modules/product_links/application/generation_integration_test.go`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.test.ts`
- M-09 feature `validation.md`, `checkpoint.md`, fixed-SHA review/QA evidence, and
  final `validation-result.md` required by the harness

No other production, migration, Oracle runner, SDK runtime, M-06, mission, or
roadmap path is authorized. A proved need outside this list stops the attempt.

## Required Proof

- Product-link generation tests prove canonical IDs are used throughout and legacy
  positive `ProductID` with nil/invalid `InternalProductID` creates no resolved or
  persistable canonical candidate.
- Conflict, deduplication, candidate-ID stability, and existing generation
  integration tests pass using positive canonical IDs.
- OpenAPI/SDK required-nullable parity passes for all canonical product fields.
- Targeted product_links Go tests and SDK tests pass at the correction SHA.
- One intentional correction commit, followed by `mpc-verifier` fixed-SHA review.
  After review Pass, run proportional QA; QA alone may pass M-09.

## Terminal Rule

This exception has no retry. Any remaining C01 failure, new architecture/contract
conflict, or path expansion returns terminal `failed` to Portfolio. C02-C05 evidence
may be reused only if review confirms this correction did not invalidate it.
