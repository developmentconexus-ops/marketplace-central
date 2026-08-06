# ADR-05 (two-digit citations) — what they assert
**Harvested:** 2026-08-05 · **Two-digit citations:** 49 (live code: 4, .mnfs only: 45) · **Three-digit `ADR-005` citations:** 18
**Existing document `005-mercado-livre-first-control-plane.md`:** "Mercado Livre First Control Plane" (VTEX removal, Mercado Livre becomes the first marketplace control plane) — DOES NOT MATCH: none of the four two-digit assertions below is "drop VTEX, Mercado Livre first"; the document is a repo-wide platform choice, the two-digit citations are all mission-scoped engineering rules.

## Assertion A1 — Oracle read-adapter batch/import size ceilings: IN-list chunking at 500, `ImportMarginInputs.Limit` ceiling 200 (→422), `GetSalesHistory` row cap 5000 (→explicit `truncated=true`, never silent) (mission: MIS-002-oracle-read-rearchitecture)
- Citations: 9 (live code: 0)
- Verbatim: "ADR-05 IN-list chunking at 500; `ImportMarginInputs.Limit` ceiling 200; `GetSalesHistory` row cap 5000"
- Anchors:
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:115`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:118`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:165`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/readiness-review.md:35`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/readiness-review.md:43`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/research/oracle-runtime-dimensioning.md:32`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/M-04-server-cache/milestone.md:25`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/M-04-server-cache/F-01-freshness-cache/feature.md:17`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/M-03-batch-inventory-profitability-sankhya/F-02-cost-tax-batch/feature.md:25`

## Assertion A2 — Mercado Livre adapter has one owner: a single milestone extends ALL `connectors/mercado_livre` read surfaces (sale_price, price_to_win, catalog search/detail, `products/{id}/items` flag-gated, shipments); the ML catalog-offers read flag defaults OFF (mission: MIS-004-mvp-demo)
- Citations: 13 (live code: 2)
- Verbatim: "Adapter ML: dono ÚNICO de TODAS as extensões ... consumidores usam ports normalizados IC-06"
- Anchors:
  - `.mnfs/MIS-004-mvp-demo/mission.md:89`
  - `.mnfs/MIS-004-mvp-demo/mission.md:180`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:40`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:12`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:51`
  - `.mnfs/MIS-004-mvp-demo/research/ml-read-ports-interface-contract.md:21`
  - `apps/server_core/internal/composition/root.go:1002`
  - `docker-compose.yml:30`

## Assertion A3 — Replacement-before-deletion + shrinking read-guard allowlist: `GET /orders`/`GET /orders/{id}` enrich stops calling live Mercado Livre readers and switches to Postgres rows written by the ingester (zero ML calls on the read path); a live site can only be removed from the allowlist in the same commit that its Postgres replacement lands (mission: MIS-007-ml-sync)
- Citations: 15 (live code: 0)
- Verbatim: "Replacement-before-deletion (ADR-05): `GET /orders/{id}`'s enrich path stops calling the live Mercado Livre shipment/buyer-fiscal readers"
- Anchors:
  - `.mnfs/MIS-007-ml-sync/mission.md:44`
  - `.mnfs/MIS-007-ml-sync/mission.md:170`
  - `.mnfs/MIS-007-ml-sync/mission.md:276`
  - `.mnfs/MIS-007-ml-sync/M-02-sync-core-seam/F-04-read-guard-allowlist/feature.md:25`
  - `.mnfs/MIS-007-ml-sync/M-03-orders-shipment-persist/F-03-read-path-switch/spec.md:5`
  - `.mnfs/MIS-007-ml-sync/M-03-orders-shipment-persist/F-03-read-path-switch/feature.md:25`
  - `.mnfs/MIS-007-ml-sync/M-03-orders-shipment-persist/_chip-m03/EVIDENCE.md:188`
  - `.mnfs/MIS-007-ml-sync/research/orders-persistence-interface-contract.md:40`
  - `.mnfs/MIS-007-ml-sync/research/sync-health-interface-contract.md:137`

Note: the same MIS-007 planning wave, before reconciliation, used the token `ADR-05` for a DIFFERENT candidate rule ("Divergences: dedicated table, detected at ingest, auto-resolve" — see Amendments); that candidate reading was superseded before ratification and does not appear in the final `mission.md`.

## Assertion A4 — Auto-approve a product link only on an unambiguous anchor: CODPROD **and** EAN must both resolve to the same product; a single matching anchor becomes a one-click "confirmation" state, never silent auto-approve; conflicting/colliding/hard-negative anchors stay in REVIEW (mission: MIS-006-integracao-fundacao; AMENDADO 2026-07-25 D-121 → D-121-2, RATIFIED-BY-OPERATOR)
- Citations: 12 (live code: 2)
- Verbatim: "auto-aprova SÓ o concordante (CODPROD e EAN no mesmo produto); âncora única vira CONFIRMAÇÃO"
- Anchors:
  - `.mnfs/MIS-006-integracao-fundacao/mission.md:137`
  - `.mnfs/MIS-006-integracao-fundacao/mission.md:211`
  - `.mnfs/MIS-006-integracao-fundacao/M-05-auto-vinculo/milestone.md:24`
  - `.mnfs/MIS-006-integracao-fundacao/M-05-auto-vinculo/milestone.md:75`
  - `.mnfs/MIS-006-integracao-fundacao/M-05-auto-vinculo/milestone.md:179`
  - `.mnfs/MIS-006-integracao-fundacao/M-05-auto-vinculo/EVIDENCE.md:28`
  - `apps/server_core/internal/modules/product_links/application/resolution_service.go:217`
  - `apps/server_core/internal/modules/product_links/application/auto_link_policy_test.go:14`

## Contradictions
- Four unrelated decisions share the token `ADR-05`, one per mission (MIS-002 Oracle batch ceilings, MIS-004 ML-adapter single-owner, MIS-007 replacement-before-deletion, MIS-006 auto-approve anchor policy). Nothing in the repo disambiguates the number across missions; each mission's `mission.md` renumbers its own ratified-decisions table starting near ADR-01, matching the pattern already recorded for ADR-09 in `adr-009-citations.md`.
- Internal MIS-007 churn: during P3 planning the candidate `ADR-05` named "Divergences: dedicated table, ingest-time, auto-resolve, two-sided timestamps" (`.mnfs/MIS-007-ml-sync/planning-reviews/p3-claude-candidate-r01.md:53`, `p3-reconciliation-r01.md:18` — reconciled to shape `A-10`). The ratified `mission.md` instead assigns `ADR-05` to "replacement-before-deletion + shrinking allowlist" (a different rule); the divergences-table rule survives in the mission under a different number. A reader following only the candidate doc would get the wrong rule for the final code.

## Amendments
- MIS-006's A4 is itself an amendment: `mission.md:138` records "AMENDADO 2026-07-25 (D-121) · RATIFIED-BY-OPERATOR — supersede a redação original ('só EAN-exato-ÚNICO')"; the live code (`resolution_service.go:217`, `auto_link_policy_test.go:14`) implements the amended (D-121-2) rule only, not the original EAN-exact-unique wording.
