# ADR-013 — what the citations assert
**Harvested:** 2026-08-05 · **Total citations (excl. scripts/.runs):** ~40 (mission.md `research/*.md` restatements, milestone/feature briefs, validation artifacts, and Go source/tests)
**Spellings found:** ADR-13 (only variant observed; no `ADR-013`, no `ADR 13`, no lowercase hits found). A shorthand `A-13` appears once inside `.mnfs/MIS-007-ml-sync/planning-reviews/p7-claude-readiness-r02.md:42` referring to the same number.

**IMPORTANT — this number is used for two unrelated rules in two different missions.** MIS-003 (`operator-cockpit-wireframe-replan`) and MIS-007 (`ml-sync`) each independently minted an "ADR-13" for a different decision. Nothing in the repo cross-references or reconciles the two; they simply collide on the same number.

## Assertion A1 — Mutation envelope is a durable protocol table plus an in-process poller, never an external queue/bus (MIS-003 usage)
- Citations: 10
- Verbatim: "no external queue infrastructure this mission (ADR-13)" / "No outbox, no message bus, no background framework — plain goroutine + ticker per ADR-13."
- Anchors:
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:112` (ADR table row: "Mutation envelope = protocolo table + in-process poller")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:95` ("mutation envelope (module home fixed in ADR-13)")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:209` (Q3 Reliability: poller resumes approved/applying after restart — ADR-13)
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/mutation-envelope-interface-contract.md:17` (module home "decided in ADR-13, not per feature")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/mutation-envelope-interface-contract.md:103` (`mutation_protocols`/`mutation_items`, in-process poller with `FOR UPDATE SKIP LOCKED`, "no external queue infrastructure this mission (ADR-13)")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/mutation-envelope-interface-contract.md:197` (`/mutations` prefix "mounted by the module that ADR-13 designates")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-03-mutation-envelope-writes/F-01-protocolo-core/feature.md:17` ("ADR-13 (table + in-process poller, NO outbox/bus)")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-03-mutation-envelope-writes/F-01-protocolo-core/feature.md:48` ("No outbox, no message bus, no background framework — plain goroutine + ticker per ADR-13")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-03-mutation-envelope-writes/milestone.md:17`
  - `.mnfs/MIS-004-mvp-demo/research/repo-baseline-2026-07-17.md:66` ("ADR-13 protocolo table + in-process poller; lifecycle 8 estados")

## Assertion A2 — Any new ingestion/hydration path for listings must keep feeding the product-links snapshot observer non-regressively, and must emit one row per item (no per-variation flattening); the `listings` PK sentinel is untouched (MIS-007 usage)
- Citations: ~28
- Verbatim: "acoplamento nomeado ADR-13; snapshots/âncoras não-regressivos vs pull pré-mudança" / "UMA linha por item (VariationID = NoVariationID, ADR-13)"
- Anchors (docs):
  - `.mnfs/MIS-007-ml-sync/mission.md:242` ("ADR-13 `listing_variations` aditiva; PK de `listings` NÃO muda")
  - `.mnfs/MIS-007-ml-sync/research/listings-sync-interface-contract.md:36,66,111,146,172` (PK tuple, snapshot re-vínculo coupling, must-fail "snapshot observer starved")
  - `.mnfs/MIS-007-ml-sync/M-04-listings-backfill-ingest/milestone.md:57,67`
  - `.mnfs/MIS-007-ml-sync/M-04-listings-backfill-ingest/validation-contract.md:88` (criterion "Âncoras de snapshots não-regressivas (ADR-13)")
  - `.mnfs/MIS-007-ml-sync/M-04-listings-backfill-ingest/_chip-m04/EVIDENCE.md:24,58,329` (regression found and fixed against this ADR-13 must-fail)
  - `docs/superpowers/plans/2026-08-01-p1-espinha-preco-estoque.md:934`
- Anchors (Go code, different subsystem than docs):
  - `apps/server_core/internal/modules/listings/ports/ingestion.go:80` ("feeding the product-links SnapshotObserver (ADR-13)")
  - `apps/server_core/internal/modules/listings/adapters/connectors/backfill.go:63` ("ADR-13 non-regression: BEFORE mapping to canonical listing rows...")
  - `apps/server_core/internal/modules/listings/adapters/connectors/backfill_test.go:75` ("TestMultigetHydratorFeedsObserverBeforeMappingRows is the ADR-13 must-fail proof")
  - `apps/server_core/internal/modules/listings/adapters/connectors/multiget_mapper.go:34,167,173` (one row per item; observer-feed parity)
  - `apps/server_core/internal/modules/listings/adapters/connectors/multiget_mapper_test.go:202,209`
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/items_multiget_reader.go:67` ("leaving that ADR-13 seam honest-empty for lack of a typed field")
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/items_multiget_reader_test.go:441`

## Contradictions
- **Cross-mission number collision.** ADR-13 in MIS-003 means "mutation envelope = table + in-process poller" (a write-path architecture decision). ADR-13 in MIS-007 means "listings snapshot-observer non-regression + one-row-per-item + PK sentinel stability" (a read/ingestion-path decision). These are unrelated decisions sharing one number; no repo artifact acknowledges or reconciles the collision.
- **Internal drift inside Assertion A2, self-detected and repaired.** The `listing_variations` PK tuple was stated two incompatible ways under the ADR-13 label: 5 columns including `installation_id` (`.mnfs/MIS-007-ml-sync/mission.md:224-225,242-243`) vs. 4 columns omitting it (`research/listings-sync-interface-contract.md:65` pre-fix, `M-04-listings-backfill-ingest/F-01-listings-ddl/feature.md:31` pre-fix). Flagged in `planning-reviews/p7-seat5-doublepass-r02.md:27-32` and `p7-claude-readiness-r02.md:39-42`, then repaired to the uniform 5-column tuple per `planning-reviews/p7-seat-star2-r03.md:29-32` — the fix is recorded, not the original defect state.

## Exceptions / carve-outs
- None found beyond the PK-consolidation deferral already noted in A2: "Consolidação de PK (tirar sentinela `'-'`) = missão futura nomeada, NUNCA nesta (ADR-13)" (`.mnfs/MIS-007-ml-sync/research/listings-sync-interface-contract.md:146`) — i.e. ADR-13 (MIS-007 sense) explicitly forbids touching the PK sentinel in the current mission and pushes that to a future one.
