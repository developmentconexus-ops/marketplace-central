# ADR-06 (two-digit citations) — what they assert
**Harvested:** 2026-08-05 · **Two-digit citations:** 60 (live code: 6, .mnfs only: 54) · **Three-digit `ADR-006` citations:** 14
**Existing document `006-oracle-internal-read-owned-by-mpc.md`:** "MPC owns Oracle internal-read adapters inside `apps/server_core`" — DOES NOT MATCH: none of the four assertions below is about who owns the Oracle read adapter; MIS-002's own ADR-06 (below) is a harness/process rule, not this architectural ownership decision, and the document text never appears in any two-digit citation.

## Assertion A1 — Listings MASS-CLOSURE dies; absence is never inferred mid-run: a listing's status is only ever marked from a COMPLETE run (`last_seen_at`/`absent_since`, keep-absent lifecycle); `ApplyCompletedPull` stays the single writer but stops emitting the unconditional "close everything" UPDATE (mission: MIS-007-ml-sync)
- Citations: 26 (live code: 6)
- Verbatim: "ADR-06, audit D-120 F1, risk R-B: a run truncated by a 429/deadline/kill must never wipe listings it simply never got to"
- Anchors:
  - `apps/server_core/migrations/0090_listings_e3_fields_status_relax.sql:48`
  - `apps/server_core/internal/modules/listings/ports/ingestion.go:50`
  - `apps/server_core/internal/modules/listings/application/backfill_test.go:59`
  - `apps/server_core/internal/modules/listings/adapters/postgres/repository_integration_test.go:24`
  - `apps/server_core/internal/modules/listings/adapters/postgres/repository.go:391`
  - `apps/server_core/internal/modules/listings/adapters/postgres/repository.go:503`
  - `.mnfs/MIS-007-ml-sync/mission.md:177`
  - `.mnfs/MIS-007-ml-sync/mission.md:278`
  - `.mnfs/MIS-007-ml-sync/research/listings-sync-interface-contract.md:38`
  - `.mnfs/MIS-007-ml-sync/M-04-listings-backfill-ingest/F-02-mass-closure-replacement/feature.md:25`
  - `.mnfs/MIS-007-ml-sync/M-04-listings-backfill-ingest/_chip-m04/EVIDENCE.md:203`

Note: this is the final ratified reading. Earlier in the same mission's P3 planning, the candidate token `ADR-06` instead named "client resiliency centralized before backfill" (backoff/jitter/token-bucket) — see Contradictions; that candidate was reconciled into shape `A-01` and does not survive under `ADR-06`.

## Assertion A2 — Identity state and commercial verdict stay distinct: pre-listing verdict emits one of ACCEPT/REVIEW/REJECT/NO_CANDIDATE/NO_PRICE_EVIDENCE/INSUFFICIENT_MARKET; a null `buy_box` never becomes a fabricated zero; copy never promises an automatic "best price" pre-listing (mission: MIS-004-mvp-demo)
- Citations: 19 (live code: 0)
- Verbatim: "Estados honestos ACCEPT/REVIEW/REJECT/NO_CANDIDATE/NO_PRICE_EVIDENCE/INSUFFICIENT_MARKET; veredicto ≠ identidade; buy_box null nunca vira zero"
- Anchors:
  - `.mnfs/MIS-004-mvp-demo/mission.md:156`
  - `.mnfs/MIS-004-mvp-demo/mission.md:177`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:47`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:13`
  - `.mnfs/MIS-004-mvp-demo/validation-contract.md:37`
  - `.mnfs/MIS-004-mvp-demo/M-06-produto-detalhe/validation-contract.md:52`
  - `.mnfs/MIS-004-mvp-demo/M-05-anuncios-sinais/validation-contract.md:67`
  - `.mnfs/MIS-004-mvp-demo/readiness-review.md:27`
  - `.mnfs/MIS-004-mvp-demo/DESIGN-TARIFAS-ML.md:481`

## Assertion A3 — Three named holes in the Mercado Livre adapter (backoff/429, live tariff, raw DTO persistence) are explicitly deferred out of the foundation mission; the mature adapter is reused, never rewritten (mission: MIS-006-integracao-fundacao)
- Citations: 9 (live code: 0)
- Verbatim: "Os 3 buracos do adapter (F-ADAPTER-1: backoff/429, tarifa live, Raw DTO) → ver ADR-06" / "reuso, não reescrita"
- Anchors:
  - `.mnfs/MIS-006-integracao-fundacao/mission.md:78`
  - `.mnfs/MIS-006-integracao-fundacao/mission.md:194`
  - `.mnfs/MIS-006-integracao-fundacao/mission.md:217`
  - `.mnfs/MIS-006-integracao-fundacao/M-07-f37-discovery/milestone.md:121`
  - `.mnfs/MIS-006-integracao-fundacao/M-04-sankhya-adapter/milestone.md:72`
  - `.mnfs/MIS-006-integracao-fundacao/M-03-xlsx-adapter/milestone.md:19`
  - `.mnfs/MIS-006-integracao-fundacao/interface-contracts-mis006.md:179`

## Assertion A4 — Harness/process decision: MIS-002 execution runs under MNFS + hub-style acceptance (gpt-5.6-luna high workers), gated by an EXPLAIN PLAN review before the M-02 Oracle cutover, and the live-Oracle lane runs read-only through a governed Docker runner (godror/CGO/Instant Client compatibility) (mission: MIS-002-oracle-read-rearchitecture)
- Citations: 2 (live code: 0)
- Verbatim: "Execution: MNFS + hub-style acceptance; workers gpt-5.6-luna high; EXPLAIN PLAN gate before M-02 cutover"
- Anchors:
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:116`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:144`

Note: this is a process/harness ratification, not an architecture decision — the same category the operator's renumbering registry (`RENUMBERING-REGISTRY.md`) moves out of the ADR series into `docs/HARNESS-PROFILE.md`.

## Contradictions
- Four unrelated rules share the token `ADR-06` across four missions (MIS-002 process gate, MIS-004 identity/verdict separation, MIS-006 deferred adapter gaps, MIS-007 MASS-CLOSURE death). No cross-mission registry disambiguates the number.
- Internal MIS-007 churn: P3 candidate `ADR-06` = "client resiliency centralized before backfill" (`.mnfs/MIS-007-ml-sync/planning-reviews/p3-claude-candidate-r01.md:63`), reconciled to shape `A-01` (`p3-reconciliation-r01.md:19`) and dropped from the number; the ratified `mission.md:177` instead assigns `ADR-06` to MASS-CLOSURE death. A reader of the P3 candidate doc alone would get a rule that never shipped under this number.

## Amendments
None found specific to ADR-06 itself (A1's lifecycle rule is stated as design-from-inception in `mission.md`, not as an amendment of an earlier ADR-06 reading). MIS-004's mission.md:98 lists ADR-06 among "Accepted trade-offs" ("demo pode exibir NO_PRICE_EVIDENCE em produtos do cliente") but this is a scoping note, not a textual amendment.
