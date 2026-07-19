# CHIP-FACT-FE — P6 dual-gate result

Reviewed tip: `390b35bcd647f261c36a8df92033032346bd88f3` (base 910d8196)
Mode: dual-gate AGREEMENT (codex quota DEAD til 2026-07-25 — Claude-only lane: cold Opus + sonnet adversarial refute).

## Gate A — cold Opus, independent review
VERDICT: **PASS / NO BLOCKER**.
Confirmed explicitly: number→string membership seam consistent end-to-end (single `factId`
stringification point; every membership/mutation/key/aria path receives stringified id); tests
genuinely exercise the seam (arrayContaining checks fail if membership breaks); cursor-drain
correct; ADR-17 honest "—" rendering; scope compliance (no backend/OpenAPI/SDK, no brand_name).
3 MINOR non-blocking: formatCurrency empty-string→"R$ 0.00" edge (contract says decimal-string|null,
"" not a contract value); multi-page-drain test-coverage gap; simulator deletion residue (expected,
hub-barred).

## Gate B — sonnet adversarial (default-skeptic, attack to break)
VERDICT: **CANNOT-REFUTE** (survived every attack vector).
Attacks + results:
1. id seam leak → NONE. `factId` is the sole stringification point; every path (toggle, select-all,
   clear, create-from-first-check, checked-state, key, aria-label) gets the stringified id. Edge
   `internal_product_id: 0` → `String(0)="0"` truthy — no falsy-id bug.
2. `loadAllFacts` cursor loop → terminates correctly; `next_cursor` `null` AND `""` both fall out of
   `while(cursor)` truthy check; no off-by-one, no dropped last page. (No iteration cap — coverage
   gap only, see caveat.)
3. Currency/stock → `formatCurrency` parses decimal strings correctly (`"1234.5"→"R$ 1234.50"`),
   em-dash only on null/NaN. Stock `?? "—"` (not `||`) → quantity `0` renders `0`, not `"—"`.
4. Search filter → `?.…includes() ?? false` null-safe, no crash on null description/reference.
5. Test weakening → assertions mechanically re-stringified but retain strength; `arrayContaining` /
   `not.arrayContaining` still fail if membership logic breaks. No vacuous test.

Caveat (not a defect): multi-page cursor-drain branch untested (all mocks single-page); logic
verified correct; demo dataset single-page so branch inert tomorrow. Post-demo hardening candidate.
Simulator residue (index.css:6 @source, apps/web/package.json dep) = expected debris from hub-barred
package deletion; harmless (Tailwind glob on missing path → zero matches, no error).

## AGREEMENT
Gate A PASS ∧ Gate B CANNOT-REFUTE → **P6 GREEN**. No defect surfaced by either lane.
Both flagged the SAME two low/cosmetic items (untested drain branch; simulator residue) — neither
blocking; residue routed to hub as cross-seam cleanup REQUEST.
