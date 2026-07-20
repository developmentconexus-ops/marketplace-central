# P7 Claude Readiness Fold — round 01

```yaml
round: 01
manifest: planning-reviews/p7-input-r01.sha256
manifest_top_digest: a9ba45f1dbc2e6255770a4bcab14c6e0508d23dbf717a1129e4a328eedf9d92d
crew: 5 cold mission-reviewer (parallel, Claude — Sol rebind per codex quota-wall)
claude_side_verdict: Needs revision
```

## Per-★ fold (computed, union of findings; a ★ FAILS if ANY covering reviewer FAILs)

| ★ | Criterion | Verdict | Source |
|---|-----------|---------|--------|
| ★1 | Completeness | PASS | R1 |
| ★2 | Consistency | **FAIL** | R2, R5 |
| ★3 | Seam Ownership | **FAIL** | R2 |
| ★4 | Verifiability | PASS | R3 |
| ★5 | Traceability | PASS | R1 |
| ★6 | Evidence Honesty | PASS | R3 |
| ★7 | Security Posture | PASS | R4, R5 |

**Claude-side joint: Needs revision** (★2 + ★3 FAIL). Sol NOT dispatched this round (protocol: Sol
runs only after Claude-side Ready).

## Blocking findings (all valid; all repaired in-place — see readiness-review.md disposition)

1. **★2 Div1 — sync_state key drift (M-07).** M-07 keyed `sync_state` by `(tenant_id,
   codigo_produto, entity)`; E8/M-01 key is `(tenant_id, installation_id, entity)` (one row per
   installation/entity, per-product state in `cursor JSONB`). Loci: `M-07/milestone.md:154`,
   `M-07/validation-contract.md:38`. FIXED → key restated to E8 shape, codigo_produto accreted in cursor.
2. **★2 Div2 — ownership matrix vs M-07 conditional migration.** `mission.md:224` M-07 Migração
   `—` contradicted M-07's conditional identity migration. FIXED → matrix cell = bloco C conditional.
3. **★2 R5 locus1 — no Error Matrix in IC.** Feature-returned typed errors (400 invalid
   active_source, ErrUnknownActiveSource, 403 PolicyAgent, 429, chain-read→"—", missing
   codigo_produto) had no consolidated IC row. FIXED → `## Error Matrix` added to
   interface-contracts-mis006.md.
4. **★2 R5 locus2 — undeclared list ordering.** `listErpImports` / Fila-Resolvidos declared no
   sort order. FIXED → M-06 F-01 EARS declares `imported_at DESC` / `created_at DESC`.
5. **★3 — M-07 conditional migration on M-02-owned `products_mirror` without grant/block.** No
   named additive-lock grant / unallocated block. FIXED → identity persisted in NEW M-07-owned
   `product_catalog_identity` table (never ALTER mirror) + bloco C allocated in mission matrix +
   architecture-map migration row.

## Advisories folded (non-blocking, repaired to reduce round-02 noise)

- MC-10 evidence types missing `could-not-run` (mission VC) → added.
- No `## Quality Attributes`/`## Non-Functional Scope` in mission.md (R4/R7 flagged for ★11) → added.
- `owner: TBD` only on M-01; status enum drift (planned vs draft) → normalized to draft, owner dropped.
- IC ImportSource 3rd-value naming hedge (R5 advisory) → covered by Error Matrix + M02-C8 (ImportSource
  keeps 2 values; source-identity is the `active_source`/`source` enum). No further change needed.

## Reviewer artifacts

R1 ★1+★5, R2 ★2+★3, R3 ★4+★6, R4 ★7, R5 adversarial ★2+★7 — full verbatim outputs captured in the
dispatch transcript; findings folded above. Read-only crew; no reviewer edited artifacts.
