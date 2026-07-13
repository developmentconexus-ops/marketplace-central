# F-02 Oracle Catalog Cutover — Validation

## Verdict

Implemented, pending independent review/QA.

## Changed Paths

- Catalog canonical read adapter/application/transport and internal-read query identity.
- Server composition/config removal of active MSDB wiring.
- Catalog OpenAPI endpoint responses and SDK method types in lockstep.

## Proof

- `GOCACHE=...\.gocache go test ./internal/modules/catalog/... ./internal/modules/internal_read/... ./cmd/server ./internal/composition ./tests/unit ./internal/modules/pricing/adapters/catalog` — PASS.
- `npm test --workspace @marketplace-central/sdk-runtime` — PASS (38 tests).
- Active runtime residue scan over server cmd/composition/platform, governance, and dev composition for `platform/msdb`, `MS_DATABASE_URL`, `MS_TENANT_ID`, and the MetalShopping catalog adapter — PASS (no matches).
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/run-live-oracle-docker.ps1` — PASS against the governed read-only container lane. Output is intentionally suppressed by the runner; no credentials, raw rows, or PII were recorded.

## Limitations

The compatibility migration for persisted consumer references remains F-03 work.
Catalog read unavailability returns `source_unavailable` without an MSDB fallback.
