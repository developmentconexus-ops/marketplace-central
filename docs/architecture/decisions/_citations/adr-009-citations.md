# ADR-09 — what the citations assert
**Harvested:** 2026-08-05 · **Total citations (excl. scripts/.runs):** 37
**Spellings found:** `ADR-09` (dominant); no `ADR-009`, `ADR-9`, or `ADR 09` hits found.

## Assertion A1 — Every fee value carries provenance (layer, origem, coletado_em); a display without provenance fails the milestone
This is the dominant, code-enforced meaning (MIS-007 `channelfees`/`orders` domain).
- Citations: 26
- Verbatim: "proveniência SEMPRE junto do número; consumidor que exibe número sem proveniência reprova milestone (ADR-09)."
- Anchors:
  - `apps/server_core/internal/modules/channelfees/domain/fee.go:2`
  - `apps/server_core/internal/modules/channelfees/domain/fee.go:92`
  - `apps/server_core/internal/modules/channelfees/adapters/postgres/reader.go:47`
  - `apps/server_core/migrations/0090_listings_e3_fields_status_relax.sql:21`
  - `.mnfs/MIS-007-ml-sync/research/channel-fees-interface-contract.md:25`
  - `.mnfs/MIS-007-ml-sync/research/channel-fees-interface-contract.md:73`
  - `.mnfs/MIS-007-ml-sync/research/orders-persistence-interface-contract.md:82`
  - `.mnfs/MIS-007-ml-sync/M-07-pricing-fee-read/validation-contract.md:142`
  - `.mnfs/MIS-007-ml-sync/mission.md:205`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p5-claude-decomposition-audit-r06.md:164`

## Assertion A2 — (MIS-004 mission, unrelated subject) Polling/GET only against Mercado Livre; "live" data always requires a visible refresh, never silent auto-freshness
Same token, different mission, different rule — see Contradictions.
- Citations: 6
- Verbatim: "ADR-09 dado \"ao vivo\" exige refresh" / "ADR-09 Polling/GET only"
- Anchors:
  - `.mnfs/MIS-004-mvp-demo/mission.md:93`
  - `.mnfs/MIS-004-mvp-demo/mission.md:98`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:68`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:17`

## Assertion A3 — (MIS-001/MIS-003 lineage) Proportional security: writes off by default, no PII/secrets in evidence
Only reachable via one summary line, not a primary mission table row for MIS-007/004; listed as a "still-binding MIS-001 ADR" in MIS-003 research.
- Citations: 1
- Verbatim: "ADR-009 proportional security (writes off by default, no PII/secrets in evidence)"
- Anchors:
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/market-intelligence-digest.md:35`

## Contradictions
- **Number collision across missions.** ADR-09 names three unrelated rules depending on which mission's `mission.md` is read: fee-provenance-mandatory (MIS-007, code-enforced), polling/GET-only-with-visible-refresh (MIS-004), and proportional-security/no-PII (MIS-001, only echoed once in MIS-003). No cross-mission ADR registry disambiguates the number; each mission renumbers its own ratified-decisions table starting near ADR-01.
- **Numbering churn within MIS-004 planning.** `p3-reconciliation-r01.md:16` shows the "zero writes ML" assertion entering reconciliation as candidate `ADR-08` reconciled to ratified `ADR-09` — i.e. even inside one mission the number was reassigned mid-planning (P3 candidate vs. reconciled).

## Exceptions / carve-outs
- `apps/server_core/internal/modules/channelfees/domain/fee.go:92` carves out Layer 3 fee rows from the provenance ladder entirely: "Layer 3 never participates in either ladder (ADR-09 — layer 3 rows may ...)" — an explicit exemption from the general A1 rule for one fee layer.
- `.mnfs/MIS-007-ml-sync/mission.md:205` records the "fee ledger enumerado" reading of ADR-09 as **already dead** ("JÁ MORTO — nenhum milestone reivindica"), i.e. one candidate assertion under this number was ratified then abandoned without being struck from the citing text.
