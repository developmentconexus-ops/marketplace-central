# F-02 Oracle Catalog Cutover — Validation

## Verdict

Corrected after fixed-SHA review; pending independent review/QA.

## Changed Paths And Outcomes

- Catalog canonical read adapter/application/transport and internal-read Oracle identity.
- `TGFPRO.REFERENCIA` populates only the manufacturer/reference field; EAN
  remains null until a governed barcode source exists, and EAN-only lookup is
  rejected as unsupported rather than guessed.
- Catalog source-fact projection maps missing/stale/conflict flags explicitly
  and returns constructor failures instead of silently emitting a zero value.
- Product-link internal IDs are positive at domain, application, HTTP,
  OpenAPI, and SDK request boundaries.
- Server composition/config remains free of active MSDB wiring.

## Proof

- Absolute `GOCACHE=C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache` targeted and broader Go suites over catalog, internal_read, product_links, classifications, pricing, composition, unit tests, migrations, and migration runner — PASS.
- `npm test --workspace @marketplace-central/sdk-runtime` — PASS (39 tests).
- OpenAPI/SDK positive-product-ID parity assertion — PASS.
- Governed runner Pester contract — PASS (14/14).
- Active runtime residue scan for `platform/msdb`, `MS_DATABASE_URL`,
  `MS_TENANT_ID`, and MetalShopping — PASS (no matches).
- Post-commit governed command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/run-live-oracle-docker.ps1`.
- Named sanitized external evidence target:
  `F-02-oracle-catalog-cutover/_fixed-sha-oracle-evidence.md`. It contains only
  frozen SHA, source, observed time, `read_only=true`, positive-CODPROD
  observation, command, and exit status; no identifier value, raw row,
  credential, or PII.

## Limitations

Catalog read unavailability returns `source_unavailable` without an MSDB fallback.
The named Oracle evidence is intentionally produced after the correction commit
so its `frozen_sha` equals the reviewed implementation SHA; it is external,
uncommitted execution evidence.
