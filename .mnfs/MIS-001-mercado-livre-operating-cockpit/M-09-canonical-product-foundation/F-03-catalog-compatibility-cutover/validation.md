# F-03 Catalog Compatibility and Cutover — Validation

## Verdict

Corrected after fixed-SHA review; pending independent review/QA.

## Compatibility Outcome

Historic classification, enrichment, and pricing strings remain immutable
evidence. Forward migration `0035_catalog_codprod_compatibility.sql` maps only
when exactly one governed, read-only Oracle/Sankhya candidate proves exact
positive `legacy_product_id = internal_product_id::text` equality. Zero
candidates persist `not_found`; duplicates persist `identity_conflict`; both
keep `internal_product_id` null. Active repositories read only `mapped` rows and
reject non-positive or noncanonical string inputs.

## Proof

- Targeted and broader Go tests pass with the absolute repository `.gocache`,
  including mapped, unmapped, duplicate/ambiguous, source-quality, and positive
  product-ID cases.
- Full forward migration chain applied 35 migrations and was idempotent on an
  isolated harness database — PASS.
- SQL readback across classification, enrichment, and pricing produced exactly
  one `mapped`, one `not_found`, and one `identity_conflict` per consumer while
  preserving every legacy value — PASS. The fixture database was removed.
- Embedded migration contract and application resolver tests — PASS.
- Active MSDB residue scan — PASS.
- SDK runtime — PASS (39/39); OpenAPI/SDK parity — PASS.
