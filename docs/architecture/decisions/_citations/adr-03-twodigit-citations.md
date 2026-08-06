# ADR-03 (two-digit citations) — what they assert
**Harvested:** 2026-08-05 · **Two-digit citations:** 55 (live code: 10, .mnfs only: 45)
**Three-digit `ADR-003` citations:** 2 (1 is the document's own title line `003-integration-spec-split-and-sequencing.md:1`; 1 real external citation in `docs/superpowers/plans/2026-08-05-arquitetura-protocolo-de-modulo-plan.md`)
**Existing document `003-integration-spec-split-and-sequencing.md`:** "Split integrations implementation into operational specs after foundation" (OAuth → fee sync → FE UX, accepted) — DOES NOT MATCH: none of the four live assertions (raw-selective/PII, TTL cache, identity correction, active-source config) is about sequencing integration delivery specs.

## Assertion A1 — Raw is selective, PII never: `buyer.billing_info` is stripped from `/orders/{id}` before persistence; `raw jsonb` exists only on `listings`, never on `order_shipments` (mission: MIS-007, ratified — dominant, code-enforced)
- Citations: 35 (live code: 10)
- Verbatim: "buyer.billing_info stripped (ADR-03) — both from a SINGLE GET"
- Anchors:
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/order_ingest_reader.go:94`
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/order_ingest_reader.go:164`
  - `apps/server_core/internal/modules/connectors/domain/order_detail.go:59`
  - `apps/server_core/migrations/0088_order_shipments.sql:5`
  - `apps/server_core/migrations/0089_orders_marketplace_orders_sync_fields.sql:29`
  - `apps/server_core/migrations/order_shipments_test.go:59`
  - `.mnfs/MIS-007-ml-sync/mission.md:155`
  - `.mnfs/MIS-007-ml-sync/research/orders-persistence-interface-contract.md:55`
  - `.mnfs/MIS-007-ml-sync/M-01-ml-client-hardening/_chip-m01/EVIDENCE.md:53`

## Assertion A2 — Ingest único idempotente; webhook and scheduler are two ports of the same path (mission: MIS-007, pre-ratification candidate — superseded)
This is the same mission as A1, but a candidate-stage assertion that was renumbered before ratification (see Contradictions).
- Citations: 3 (live code: 0)
- Verbatim: "### ADR-03: Ingest único idempotente; webhook e scheduler = 2 portas do mesmo caminho"
- Anchors:
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p3-claude-candidate-r01.md:31`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p3-claude-candidate-r01.md:94`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p3-reconciliation-r01.md:16`

## Assertion A3 — Identity: `CODPROD` is the canonical SKU, `REFERENCIA`=EAN/GTIN (checksum+uniqueness), `REFFORN`=manufacturer; `seller_sku` resolves only to `CODPROD`; fix the Oracle reader (mission: MIS-004, ratified)
- Citations: 8 (live code: 0)
- Verbatim: "ADR-03 Identidade: CODPROD=SKU canônico, REFERENCIA=EAN/GTIN (checksum+unicidade), REFFORN=fabricante"
- Anchors:
  - `.mnfs/MIS-004-mvp-demo/mission.md:87`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:26`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:10`
  - `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/PLAN.md:392`
  - `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/PLAN.md:441`
  - `.mnfs/MIS-004-mvp-demo/DESIGN-TARIFAS-ML.md:229`

## Assertion A4 — In-memory TTL cache + singleflight honoring `FreshnessPolicy.MaxAge`; linkage is never cached; evict-on-mutation (mission: MIS-002, ratified)
- Citations: 4 (live code: 0 — no `ADR-03` comment found in the M-04 server-cache implementation code)
- Verbatim: "ADR-03 In-memory TTL cache + singleflight honoring `FreshnessPolicy.MaxAge`; linkage never cached; evict-on-mutation"
- Anchors:
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:113`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:118`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:133`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/readiness-review.md:32`

## Assertion A5 — Active ERP source is per-tenant DB config; the `MC_ERP_SOURCE` env var dies (mission: MIS-006, ratified)
- Citations: 4 (live code: 0)
- Verbatim: "### ADR-03: fonte ativa = config em banco por tenant; MC_ERP_SOURCE morre"
- Anchors:
  - `.mnfs/MIS-006-integracao-fundacao/mission.md:121`
  - `.mnfs/MIS-006-integracao-fundacao/mission.md:316`
  - `.mnfs/MIS-006-integracao-fundacao/M-04-sankhya-adapter/milestone.md:270`
  - `.mnfs/MIS-006-integracao-fundacao/M-02-mirror-port-active-source/milestone.md:45`

## Contradictions
- **Five distinct meanings under one number.** Four missions ratified unrelated ADR-03s (raw/PII, identity, TTL cache, active-source config), and MIS-007 additionally burned the number once at the candidate stage on a fifth, unrelated subject (ingest idempotency) before reassigning it — see Amendments.
- **Intra-mission candidate churn (MIS-004).** Sol's counterproposal used `ADR-03` for yet another subject: `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-sol-counterproposal-r01.md:33` — "### ADR-03: Links and imports are workflows, not catalog identity" — superseded by the identity-correction ratification (A3).
- **Amended cap size, not silently shipped wrong.** `.mnfs/MIS-007-ml-sync/M-01-ml-client-hardening/F-02-items-multiget-raw-dto/feature.md:95` records: "este brief e ADR-03 citam \"1MB\" para o cap de Raw, mas a [implementação real é 256KiB]" — classified and fixed, per `_chip-m01/EVIDENCE.md:53`.

## Amendments
- MIS-007's candidate-stage `ADR-03` ("Ingest único idempotente", A2) was renumbered to ratified `ADR-04` ("Ingest resource-addressed, caminho único" — see `adr-04-twodigit-citations.md`), freeing the number for the raw/PII rule (A1) at ratification.
- MIS-007 raw-cap size: "1MB" (original `feature.md`/mission ADR-03 text) amended to the as-implemented "256KiB, per-item" — `.mnfs/MIS-007-ml-sync/M-01-ml-client-hardening/validation-contract.md:119`.
