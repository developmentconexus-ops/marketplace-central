# F-01 identity-semantics-fix — validation evidence

Feature = canonical product-identity fix (IC-01): REFERENCIA → GTIN-validate → `ean` or null;
`refforn` EXCLUSIVELY from REFFORN (never derived from EAN); `+ncm`; EAN collision signalled.
Maps to validation-contract **C03**.

Branch `chip/m01-erp-xlsx-identity` · base `59d0e62f`. Every slice: failing test first →
independent Claude-sonnet review BEFORE next slice → commit per green slice.

| Slice | What | Dispatch | Commit | Review verdict |
|---|---|---|---|---|
| F01-S1 | GTIN mod-10 canonicalization validator (catalog/domain/gtin.go) | D03 (luna/high) | `362107ec` | PASS — reviewer recomputed 4 checksums independently |
| F01-S2 | internal_read domain + oracle identity remap (REFERENCIA→ean-or-null, REFFORN→refforn, +ncm) | D05 (sol/low) | `aa9ef075` | PASS — column-align, collision, REFFORN separation traced |
| F01-S3 | catalog-page identity + catalog DTOs (2nd read path) | D07 (sol/low) | `6977a56f` | PASS — 8/8 clean; no REFERENCIA leak, mirror-fidelity vs reader.go, null-not-empty, both collision CODPRODs readable |
| F01-S4 | OpenAPI catalog schemas + manual SDK parity (ADR-12 same commit) | D10 (luna/high) | `a2c06985` | 4🔴 (quality_flags enum under-wide vs unfiltered passthrough) FIXED pre-commit → shared QualityFlag closed 11-value component |

## Contract criteria satisfied
- **C03 identity (IC-01):** GTIN checksum canonicalization (gtin.go, F01-S1); REFFORN sourced only from
  REFFORN column, never from EAN (oracle reader.go / catalog_page.go identity path, F01-S2/S3); invalid
  EAN ⇒ `ean == nil` (not empty string) + INVALID_EAN warning; EAN collision ⇒ both products readable +
  collision QualityFlag signalled. Pinned end-to-end in source_contract_test.go (F03-S3).

## Notes / carried findings
- SEMANTIC TENSION resolved truth-order (IC-01 > pre-existing code): oracle reader previously hardcoded
  EAN=nil and mapped REFERENCIA→ReferenceCode; F-01 re-maps per IC-01 and flips the old tests. (ledger §Field findings)
- F01-S4 enum superset chosen to avoid reopening merged handler; if IC-01 wants product `quality_flags`
  scoped identity-only, that is a hub handler-filter follow-up (would let the enum re-narrow).

**Status: COMPLETE** (4/4 slices merged, all reviewed green). Aggregate ladder green @ db91f385.
