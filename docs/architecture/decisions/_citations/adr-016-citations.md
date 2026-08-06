# ADR-016 — what the citations assert
**Harvested:** 2026-08-05 · **Total citations (excl. scripts/.runs):** 8
**Spellings found:** ADR-16 only.

All citations live in a single mission (`MIS-003-operator-cockpit-wireframe-replan`); no cross-mission collision found (this is a low-citation-count ADR — reported plainly per instruction 7, not padded).

## Assertion A1 — Provider-write SKU invariant: a `listing_edit`/`listing_create` write item is only valid if `SELLER_SKU` equals the linked `CODPROD`; violations are rejected at preview/pre-apply with a stable `sku_invariant_violation` failure code, never applied
- Citations: 8
- Verbatim: "SKU invariant | decided | mislinked live listings | `listing_edit`/`listing_create` require resolved link; SELLER_SKU = CODPROD | M-03/M-06 negative criteria"
- Anchors:
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:115` (ADR table row)
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/mutation-envelope-interface-contract.md:162` (error matrix: `sku_invariant_violation` — ADR-16; rejected at preview)
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/mutation-envelope-interface-contract.md:212` ("SKU invariant (ADR-16): `listing_edit`/`listing_create` items require `link.state=resolved`; enforced at preview")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-03-mutation-envelope-writes/validation-contract.md:104` (criterion M03-C06: "ADR-16 SKU invariant")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-03-mutation-envelope-writes/F-02-write-types-adapters/feature.md:17,31,55,61` (module binding, EARS requirement, negative-scenario test, validation-expectation test)
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-04-read-workspaces/F-01-produto-detalhe/feature.md:17` (binding contract restatement, "ADR-16/17")

## Contradictions
None found. All 8 citations state the identical rule (SELLER_SKU must equal linked CODPROD; enforced pre-apply/at preview; failure code `sku_invariant_violation`) with no divergent wording.

## Exceptions / carve-outs
None found.
