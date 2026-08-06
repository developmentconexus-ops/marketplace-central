# ADR-12 — what the citations assert
**Harvested:** 2026-08-05 · **Total citations (excl. scripts/.runs):** 69 raw matches. Note: 19 of those are the *same* content repeated inside one evidence transcript (`M-01-listings-read-spine/F-01-listings-module-ingestion/evidence/planner-sol-medium.log`, an agent log that re-quotes mission.md multiple times during its own investigation). Effective distinct source citations ≈ 46.
**Spellings found:** `ADR-12` (all missions); no `ADR-012`, `ADR-9`-style bare-digit, or `ADR 12` hits found.

## Assertion A1 — (MIS-003) Canonical `listings` module is read-only; composite `listing_id`, nullable-unknowns, ingestion only via the connectors capability, freshness = `fetched_at` + manual refresh
- Citations: 15 distinct (plus the log-file repeats)
- Verbatim: "ADR-12 Canonical `listings` module (read-only) ... canonical fields, composite listing_id, nullable unknowns; ingestion via connectors capability only; freshness = `fetched_at` + manual refresh"
- Anchors:
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:111`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:207`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-01-listings-read-spine/F-01-listings-module-ingestion/feature.md:17`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-01-listings-read-spine/milestone.md:17`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-02-frontend-platform-anuncios/F-02-web-query-state-components/feature.md:17`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-06-corrigir-atributo-market-contracts/milestone.md:17`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/backend-surface.md:60`

## Assertion A2 — (MIS-007) Migrations 0086+ are pre-allocated in disjoint number ranges by the dispatch hub, declared in each brief, to avoid merge collisions
- Citations: 4
- Verbatim: "ADR-12 Migrações 0086+ pré-alocadas pelo hub. Range disjunto explícito no brief"
- Anchors:
  - `.mnfs/MIS-007-ml-sync/mission.md:237`
  - `.mnfs/MIS-007-ml-sync/mission.md:403`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p7-seat1-star1-star5-r01.md:114`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p5-reconciliation-r06.md:77`

## Assertion A3 — (MIS-004) `sdk-runtime` is a hand-written manual client; every OpenAPI change lands in the same commit as the matching SDK types
- Citations: 12
- Verbatim: "ADR-12 sdk-runtime manual: OpenAPI+SDK mesmo commit; arquivo de client por milestone + barrel hub-adjudicado"
- Anchors:
  - `.mnfs/MIS-004-mvp-demo/mission.md:96`
  - `.mnfs/MIS-004-mvp-demo/mission.md:98`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:87`
  - `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/PLAN.md:98`
  - `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/F-01-identity-semantics-fix/validation.md:15`
  - `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/F-02-erp-import-module/validation.md:17`
  - `.mnfs/MIS-004-mvp-demo/M-04-vinculos-import-ui/evidence/d02-investigator-gaps.md:42`
  - `.mnfs/MIS-004-mvp-demo/M-02-price-intel-core/evidence/p2-batch-plan.md:422`
  - `.mnfs/MIS-004-mvp-demo/validation-contract.md:141`

## Contradictions
- **Number collision across all three missions.** ADR-12 names three unrelated rules: listings-module-is-read-only (MIS-003), migration-range-preallocation (MIS-007), and sdk-runtime-same-commit-discipline (MIS-004) — none aware of the others.
- **This is the number the reconstruction task itself was triggered by: confirmed non-existence of a formal record.** Multiple MIS-003 artifacts explicitly flag that ADR-12 (and ADR-17) have **no formal document** under `docs/architecture/decisions/` — only a "decided" row inside `mission.md` — and treat this as a standing documentation-governance gap, not a behavioral one:
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-01-listings-read-spine/F-01-listings-module-ingestion/plan.md:110`: "ADR-12 and ADR-17 exist only as decided rows in the MIS-003 mission ... the missing formal ADR records should be repaired by the architecture owner, not by F-01."
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-01-listings-read-spine/DECISIONS.md:140`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-01-listings-read-spine/DISPATCH-LEDGER.md:19` (logged as an **ESCALATION**, out-of-scope for the feature, never resolved by an owner in-repo)
- **Reconciliation-stage renumbering (MIS-004 P3).** `p3-reconciliation-r01.md:19` and `:20` show the on-demand-collection assertion and the sdk-runtime-manual assertion shifting between candidate numbers `ADR-11`/`ADR-12`/`ADR-13` across the two reviewers before ratification — i.e. ADR-12's MIS-004 meaning itself wasn't stable during planning.

## Exceptions / carve-outs
- A3 (sdk-runtime) carries an explicit additive-only carve-out: `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/F-01-identity-semantics-fix/feature.md:55` grants F-01 an **additive-lock grant** on `sdk-runtime/src/index.ts` catalog types — "aditivo apenas, ADR-12 mesmo commit" — i.e. the same-commit rule is scoped down to additive-only changes for that owner, not free rewrite rights.
- A1 (listings read-only) carries no live-write exception found in any citation; every anchor reaffirms "read-only" / "ingestion via connectors capability only" without qualification.
