# ADR-10 — what the citations assert
**Harvested:** 2026-08-05 · **Total citations (excl. scripts/.runs):** 27
**Spellings found:** `ADR-10` (MIS-004/MIS-007) and `ADR-010` (MIS-001/MIS-003) — same number, zero-padding differs by mission convention.

## Assertion A1 — Divergences are one-open-row-per-(entity,kind), detected at INGEST time (never on read), with both sides' observation timestamps NOT NULL
MIS-007's dominant, code-enforced meaning.
- Citations: 12
- Verbatim: "row aberta com timestamps dos 2 lados NOT NULL (ADR-10)"
- Anchors:
  - `apps/server_core/migrations/0087_divergences.sql:4`
  - `apps/server_core/internal/modules/divergences/ports/ports.go:13`
  - `apps/server_core/internal/modules/divergences/domain/divergence.go:2`
  - `apps/server_core/internal/modules/divergences/adapters/postgres/recorder.go:2`
  - `apps/server_core/internal/modules/divergences/adapters/postgres/recorder.go:84`
  - `.mnfs/MIS-007-ml-sync/research/divergences-interface-contract.md:23`
  - `.mnfs/MIS-007-ml-sync/research/divergences-interface-contract.md:37`
  - `.mnfs/MIS-007-ml-sync/mission.md:218`
  - `.mnfs/MIS-007-ml-sync/M-05-listings-fees-divergence/validation-contract.md:64`

## Assertion A2 — (MIS-004 mission, unrelated subject) DIFAL has a single source of truth inside the `pricing` module (seeded 27 UFs + sparse overrides); Simulator writes the real destination, Pedidos consumes a read-only chip
- Citations: 9
- Verbatim: "ADR-10 DIFAL fonte única no `pricing`: seed 27 UFs \"padrão 2026\" + overrides esparsos"
- Anchors:
  - `.mnfs/MIS-004-mvp-demo/mission.md:94`
  - `.mnfs/MIS-004-mvp-demo/mission.md:98`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:74`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:103`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p7-sol-readiness-r05.md:28`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p7-claude-readiness-r02.md:44`

## Assertion A3 — (MIS-001/MIS-003 lineage) Mocks/test doubles never claim to be a live integration
- Citations: 4
- Verbatim: "Mocks never claim live integration (ADR-010 carried forward)."
- Anchors:
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:227`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/market-data-interface-contract.md:53`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/market-intelligence-digest.md:35`
  - `.mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md:198` (this MIS-001 row is actually about registered Go/web command cwd/GOCACHE compatibility, a fourth distinct reading — see Contradictions)

## Contradictions
- **Number collision across missions, three times over.** ADR-10 means divergence-ledger-detected-at-ingest (MIS-007), DIFAL-single-source (MIS-004), and mocks-never-claim-integration (MIS-003) depending on which mission.md is read — three unrelated engineering rules under one label.
- **MIS-001 itself is internally split.** `mission.md:198` cites "ADR-010 Vertical validation" (isolated feature tests being misreported as an MVP journey) while `mission.md:243` cites the same `ADR-010` for a narrower compatibility rule (absolute `GOCACHE`, valid package cwd, web commands from repo root). These are related-but-distinct claims bundled under one MIS-001 ADR row.
- **Reconciliation-stage renumbering (MIS-004 P3).** `p3-reconciliation-r01.md:17` shows "polling/GET only" entering reconciliation as candidate `ADR-09` on one side and `ADR-10` on the other, reconciled as "idêntico" — the DIFAL assertion (A2 above) is a *different* row (`:18`) that also touches ADR-10/ADR-11 boundary renumbering.

## Exceptions / carve-outs
- None found beyond the general MIS-007 rule; no carve-out language ("except", "nunca" narrowing) attached to ADR-10 citations besides the ingest-vs-read timing constraint itself (detection must happen at ingest, "nunca no read").
