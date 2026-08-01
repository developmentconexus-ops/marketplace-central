# P7 Claude readiness fold — round r02

```yaml
round: r02
manifest: planning-reviews/p7-input-r02.sha256
manifest_digest: e238f447823ab01eaa1780deb84d0db3280488ab4eda864033d51134f8a4ebaa
manifest_verified: 2026-08-01 (recompute 65/65 match, top digest identical — no drift)
crew: 5/5 seats returned (cold mission-reviewer, parallel Task dispatch)
verdict: "Needs revision"
sol_dispatched: false  # Claude-side not Ready ⇒ Sol NOT dispatched (skill §2); round still counts (2/3 of cap)
```

## Seat artifacts

| Seat | Scope | Artifact | Per-★ result |
|---|---|---|---|
| 1 | ★1 + ★5 | `p7-seat1-star1-star5-r02.md` | ★1 PASS; ★5 PASS |
| 2 | ★2 + ★3 | `p7-seat2-star2-star3-r02.md` | ★2 PASS; ★3 PASS |
| 3 | ★4 + ★6 | `p7-seat3-star4-star6-r02.md` | ★4 PASS; ★6 PASS |
| 4 | ★7 adversarial | `p7-seat4-star7-r02.md` | ★7 PASS |
| 5 | ★2 + ★7 double-pass | `p7-seat5-doublepass-r02.md` | ★2 FAIL; ★7 PASS |

## Computed fold (union; a ★ FAILS when ANY covering seat returns a valid FAIL)

| ★ | Verdict | Basis |
|---|---|---|
| ★1 Completeness | PASS | seat 1 |
| ★2 Consistency | **FAIL** | seat 5 FAIL prevails over seat 2 PASS — never downgrade a valid FAIL. Seat 5's findings name criterion, excerpt, exact locus, offending token, yes-if: VALID. |
| ★3 Seam Ownership | PASS | seat 2 |
| ★4 Verifiability | PASS | seat 3 |
| ★5 Traceability | PASS | seat 1 |
| ★6 Evidence Honesty | PASS | seat 3 |
| ★7 Security Posture | PASS | seat 4 PASS + seat 5 PASS (dual coverage, both PASS) |

**Verdict: Needs revision** (1 ★ FAIL).

## Blocking findings (union)

**★2-A (seat 5) — `listing_variations` PK: 5 colunas vs 4.**
- `mission.md:224-225` ADR-13: tuple `(tenant_id, installation_id, provider, provider_listing_id, variation_id)` (5 col).
- `research/listings-sync-interface-contract.md:65` (IC-07) e `M-04-listings-backfill-ingest/F-01-listings-ddl/feature.md:31`: `(tenant_id, provider, provider_listing_id, variation_id)` (4 col).
- Load-bearing p/ DDL 0091; A-13 ratificada em `p3-reconciliation-r01.md:37` com leitura do spine (5 col). Nenhuma emenda registrada suporta o drop de `installation_id`.
- Yes-if: pinar UM tuple nos 3 loci, com razão registrada.

**★2-B (seat 5) — slug de webhook `mercadolivre` ≠ provider_code real `mercado_livre`.**
- `research/webhook-inbox-interface-contract.md:31` (IC-04) e `M-08-webhook-ingest/F-01-inbox-endpoint/feature.md:39` afirmam slug `mercadolivre` = "mesmo vocabulário provider_code das 4 superfícies".
- Fato manifestado contradiz: `research/p5-prerequisites.md:113` `provider_code: "mercado_livre"`; repo confirma (`market_adapters.go:239`, `IntegracoesPage.tsx:514`, `PedidosTable.tsx:99`). `mercadolivre` é nome de PACKAGE Go, não provider_code (`M-01/milestone.md:48` distingue).
- Blast radius: coluna `provider` do inbox + UNIQUE dedupe + join de health IC-05 + URL de callback M08-C6.
- Yes-if: trocar token p/ `mercado_livre` nos 2 loci, OU deletar a alegação de vocabulário e pinar mapping slug→provider_code explícito no IC-04.

## Advisories (union, nunca mudam veredito)

| # | Seat | Advisory | Disposição |
|---|---|---|---|
| A-1 | 1, 2, 5(dup) | M-08 VC:175 Handoff "M-09 já fechado na prática" — claim de estado de execução em artefato de planejamento | auto-fix r03 |
| A-2 | 1, 3 | fato #11 (`assumed` ~1500 req/min) não ecoado em Accepted assumptions do mission.md | auto-fix r03 |
| A-3 | 1 | PII scrub `docs/design/evidence/ml-api/` — housekeeping de execução, não defeito de planejamento | permanece declarado no VC |
| A-4 | 2, 5(dup) | M07-C2 "camada 1 → origem api_listing_prices" — origem é fixture-defined (camada 1 sem produtor nesta missão) | auto-fix r03 |
| A-5 | 2 | IC-05 cita `root.go:856`; p5-prerequisites coloca chain em `:853-855` — drift de precisão de linha em citação | auto-fix r03 |
| A-6 | 2 | M-09/F-01:94 owned paths "sync/application/" mais largo que grant `health_*` da matriz | auto-fix r03 (restate prefixo) |
| A-7 | 3 | M-07 VC "N simulações" — N não pinado | auto-fix r03 (N≥3) |
| A-8 | 3 | M-09 VC "recente" não pinado numericamente (×2) | auto-fix r03 (citar threshold IC-05) |
| A-9 | 4 | Inbox: sem retention/prune + amplificação de refetch (forge flood drena token-bucket) | NÃO auto-fixável — precisa valor decidido; registrado como dívida declarada (ver Disposition) |
| A-10 | 4 | IC-05 agregado de webhook sem predicado de tenant (porta é tenant-parametrizada) — pinar semântica | auto-fix r03 (1 frase) |
| A-11 | 4 | M-01 baseline pré-merge sem regra de scrub/no-auth-headers | auto-fix r03 (copiar cláusula M-03) |
| A-12 | 5 | M-04/F-04:58 "enfileira/rejeita 409" hedge + sem row na Error Matrix IC-07 | auto-fix r03 |
| A-13 | 5 | M-09/F-01:87 "entities todas null-honestas" ambíguo vs `entities: []` | auto-fix r03 |
| A-14 | 5 | `validation-contract.md:124` grep `0086*..0095*` — sintaxe de glob inválida, comando não roda | auto-fix r03 |

## Repair disposition

- ★2-A e ★2-B: reparos locais dentro da repair-authority (preservam outcome/escopo/arquitetura aprovados — ★2-A restaura a leitura RATIFICADA da A-13; ★2-B alinha ao vocabulário provider_code já decidido nas 4 superfícies). Aplicar, depois r03.
- A-9 (retention/refetch-amplification): exige valor decidido pelo dono — NÃO inventar. Fica como dívida nomeada no Handoff da missão; decisão do operador pode entrar em qualquer momento sem reabrir P7 (não muda critério existente).
- Advisories auto-fixáveis: aplicar no mesmo passe de reparo (baratos, nenhum muda contrato).

## Round accounting

- r01: Needs revision (9 loci B-1..B-9) — consumida.
- r02: Needs revision (★2-A, ★2-B) — consumida.
- **r03 = ÚLTIMA rodada do cap 3.** Não-Ready em r03 ⇒ `status: blocked` + escalar todos os yes-ifs restantes.

Next action: aplicar reparos ★2-A/★2-B + advisories auto-fixáveis → author pre-check → freeze `p7-input-r03.sha256` → crew fria ×5 (r03).
