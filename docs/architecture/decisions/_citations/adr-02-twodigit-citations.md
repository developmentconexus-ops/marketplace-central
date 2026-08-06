# ADR-02 (two-digit citations) — what they assert
**Harvested:** 2026-08-05 · **Two-digit citations:** 36 (live code: 3, .mnfs only: 33)
**Three-digit `ADR-002` citations:** 2 (1 is the document's own title line `002-mpc-schema-same-cluster.md:1`; 1 real external citation in `docs/superpowers/plans/2026-08-05-arquitetura-protocolo-de-modulo-plan.md`)
**Existing document `002-mpc-schema-same-cluster.md`:** "MPC own tables live in `mpc` schema on MetalShopping's Postgres cluster" (superseded 2026-07-07) — DOES NOT MATCH: none of the live assertions (dual ERP source, resilience decorator, SourceKind/ProductSourceAdapter port, catalog cursor envelope) concerns schema/cluster topology.

## Assertion A1 — Resilience is a decorator at the ML HTTP choke point: token-bucket rate limit, exponential backoff+jitter on 429/`Retry-After`, opt-out no-retry for writes (mission: MIS-007, ratified — code-enforced)
- Citations: ~16 (live code: 2)
- Verbatim: "(ADR-02: \"opt-out no-retry para writes\")."
- Anchors:
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/resilience_decorator.go:36`
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:748`
  - `.mnfs/MIS-007-ml-sync/mission.md:148`
  - `.mnfs/MIS-007-ml-sync/mission.md:275`
  - `.mnfs/MIS-007-ml-sync/M-01-ml-client-hardening/F-01-resilience-decorator/feature.md:44`
  - `.mnfs/MIS-007-ml-sync/M-01-ml-client-hardening/milestone.md:32`
  - `.mnfs/MIS-007-ml-sync/M-01-ml-client-hardening/F-02-items-multiget-raw-dto/feature.md:27`
  - `docs/superpowers/specs/2026-08-03-mis008-operacao-diaria-design.md:192`

## Assertion A2 — Fonte ERP dual: `erp_import` module (.xlsx → snapshot + subset-Reader adapter) alongside Oracle; batched import with protocol, file hash, source, import time, per-row rejection report (mission: MIS-004, ratified)
- Citations: ~9 (live code: 0)
- Verbatim: "ADR-02 Fonte ERP dual: módulo `erp_import` (.xlsx → snapshot + adapter do subset Reader) ao lado do Oracle"
- Anchors:
  - `.mnfs/MIS-004-mvp-demo/mission.md:86`
  - `.mnfs/MIS-004-mvp-demo/mission.md:98`
  - `.mnfs/MIS-004-mvp-demo/HUB-LEDGER.md:8`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:19`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:10`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:25`
  - `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/PLAN.md:3`

## Assertion A3 — `ProductSourceAdapter` port: read-side preserved, new write/Sync side; `SourceKind` lives in a dependency-free package to avoid an import cycle (mission: MIS-006, ratified)
- Citations: ~7 (live code: 1)
- Verbatim: "reference the type without an import cycle (ADR-02)."
- Anchors:
  - `apps/server_core/internal/modules/sourcekind/sourcekind.go:4`
  - `.mnfs/MIS-006-integracao-fundacao/mission.md:113`
  - `.mnfs/MIS-006-integracao-fundacao/mission.md:208`
  - `.mnfs/MIS-006-integracao-fundacao/M-02-mirror-port-active-source/milestone.md:118`
  - `.mnfs/MIS-006-integracao-fundacao/M-04-sankhya-adapter/milestone.md:74`
  - `.mnfs/MIS-006-integracao-fundacao/_chip-m02/PLAN.md:7`
  - `docs/superpowers/plans/2026-08-05-arquitetura-protocolo-de-modulo-plan.md:98`

## Assertion A4 — Catalog API cursor envelope `{items, next_cursor, page_size}` + `as_of` (mission: MIS-002, ratified)
- Citations: 2 (live code: 0)
- Verbatim: "ADR-02 Catalog API cursor envelope `{items, next_cursor, page_size}` + `as_of`"
- Anchors:
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:112`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:118`

## Contradictions
- **Four unrelated decisions, one per mission** — same pattern as ADR-01: resilience decorator (MIS-007), dual ERP source (MIS-004), SourceKind/ProductSourceAdapter (MIS-006), catalog cursor envelope (MIS-002).
- **Intra-mission candidate churn (MIS-004).** Sol's counterproposal used `ADR-02` for a different subject: `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-sol-counterproposal-r01.md:25` — "### ADR-02: Canonical identity belongs to `catalog`". Reconciliation explicitly kept Claude's version instead: `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:25` — "Resolução: Claude ADR-02 mantido (módulo `erp_import` implementando subset do Reader port)."
- **Intra-mission candidate churn (MIS-007).** Before ratification, `ADR-02` briefly meant "Read = Postgres, ponto; enriquecimento persiste no INGEST" (`.mnfs/MIS-007-ml-sync/planning-reviews/p3-claude-candidate-r01.md:17`, echoed at `p3-reconciliation-r01.md:15` and `:49`). The ratified `mission.md:148` repointed the same number to the resilience-decorator rule instead; the "read=Postgres" content survives unnumbered inside the ratified text (§ enrichment-at-ingest) rather than under a stable ADR number.

## Amendments
None found beyond the candidate-to-ratified renumbering already listed under Contradictions.
