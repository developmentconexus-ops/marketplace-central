# MIS-004-mvp-demo

```yaml
id: MIS-004
type: mission
status: planned
owner: Mission Strategist
parent: none
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: mission
planning_phase: readiness
```

## Objective

Demo ao cliente segunda 2026-07-20, rodando no docker dev stack local: história completa "do estoque ERP (planilha .xlsx do cliente) ao veredicto honesto 'vale a pena vender?' com preço vs mercado ML, e simulador de margem real", no design novo (papel+verde), com Pedidos funcional, dados live read-only e ZERO writes no Mercado Livre.

## Outcome

Operador importa a planilha .xlsx do cliente (protocolo + relatório de rejeição), vê produtos com identidade correta (CODPROD/EAN/REFFORN), vincula aos anúncios ML da conta conectada, enxerga sinais competitivos honestos (nosso preço, vencedor, alvo, ou estados NO_PRICE_EVIDENCE/INSUFFICIENT_MARKET), simula margem real (comissão viva, frete, imposto, DIFAL por destino real, custo ERP) com veredicto, e opera Pedidos (KPIs, Fila, Lista, Kanban read-only, drawer com decomposição de retorno) — tudo no shell retematizado. Mutação com alvo PROVIDER ("aplicar preço") entra na fila M-03 com preview+protocolo, teto `previewed`, execução contra ML desabilitada; aprovação em massa de vínculos é estado LOCAL aplicado dentro do product_links com preview dry-run + auditoria própria (exceção ratificada P5-F-12 — o envelope não suporta payload por item e o stub de writes desliga o writer local junto).

## Scope

Brief §4 itens 1–9 ratificados + import .xlsx (novo, obrigatório) + Dashboard cortável + DIFAL seed mínimo. Detalhe por milestone abaixo.

## Domain Scope

- **Identidade & ERP**: correção do contrato CODPROD/REFERENCIA(=EAN)/REFFORN; módulo `erp_import` (.xlsx → snapshot + Reader adapter + protocolo); Oracle/Sankhya permanece caminho alternativo.
- **Preço & mercado**: persistência de evidência no `market` (research §6); extensões adapter ML oficial read-only; resolver determinístico; veredicto honesto pré/pós-anúncio.
- **FE**: retheme/shell primeiro; telas Vínculos, Anúncios+sinais, Produto Detalhe (header+veredicto+abas Anúncios vinculados/Estoque), Simulador completo, Pedidos (Fila/Lista/Kanban+drawer), Dashboard MPC (cortável).
- **Fiscal mínimo**: DIFAL seed 27 UFs + overrides + toggle Simulador + chip Pedidos.
- **Mutações**: fila M-03 preview+protocolo p/ "aplicar preço" (teto `previewed`); vínculos em massa = batch LOCAL no product_links (preview dry-run + apply + auditoria — exceção P5-F-12); execução provider OFF.

## Non-Scope

- Webhooks (polling cobre; MIS-005 M-06). Writes ML live (MIS-005 M-10). Scheduler diário de coleta (MIS-005 M-01).
- Scraping/busca pública/proxies/SearchAPI/Gecko/Pricefy/provedor externo não homologado (NO-GO fechado, research §3).
- DIFAL agendamento/lembrete/marcar-pago; Configurações completa (MIS-005 M-03).
- Cancelados/Devoluções/claims/NF-e produção (MIS-005 M-04). Abas Concorrência/Pedidos/Histórico do Produto Detalhe (MIS-005 M-05).
- Auth/multi-tenant real (MIS-005 M-09). Deploy cloud (MIS-005 M-11). Repasses (MIS-005 M-02).
- UI de upload .xlsx dedicada — import via endpoint (runbook); tela de import ERP = MIS-005.

## Current State

main ≥ `f4612be3` (W1 completo: m-02 FE platform @ 79d6787f, m-03 mutation envelope @ f4612be3, SAT). 15 módulos governance (mutations incluso). Migrações: 41 arquivos, topo 0044, fixture `runner_test.go` = 41. FE: rotas PT-BR + placeholders + InstallationContext + AnunciosPage + ProtocoloPage; SEM tema; nav sidebar divergente da canônica. Páginas legacy Pedidos/Vínculos/Marketplaces/Integrações deletadas. Evidência: `research/repo-baseline-2026-07-17.md` (R-01, baseline pré-merge) + `research/w1-merge-addendum-2026-07-17.md` (R-03, estado corrente).

## Clarified Decisions

- Resolved: ver `research/p1-clarified-decisions-2026-07-17.md` (P1a/P1b/P1c completos) + STOP P3 2026-07-17: grão 9 milestones full-stack aprovado; MIS-005 11 headlines aprovado.
- Accepted assumptions: ledger completo em `research/p1-clarified-decisions-2026-07-17.md` §Accepted assumptions (conta ML existente; coleta on-demand; flag products/{id}/items default OFF; colunas xlsx delegadas ao planning → IC-02; pills "em breve"; Dashboard recriado MPC; MIS-005 grão milestone).
- Accepted assumption (P5 ruling 2026-07-17, sem input do operador — P5-F-12): batch de vínculos é estado LOCAL do módulo dono, FORA do envelope de mutação M-03; `selection_resolver` opera sem payload por item. Exceção nomeada no ADR-08.
- Accepted assumption (verify-at-execution — R-04): alíquota interna MS registrada como 17,0% mas DISPUTADA (fontes citam 19,0%); seed usa 17,0% e o override do operador (`PUT /pricing/difal/{uf}`) é esperado ANTES de uso real de destino MS. Verificação na execução, não no planning.
- Accepted assumption (R-04): AL em transição de alíquota (mudança 01/04/2026); seed usa o valor pós-aumento documentado no R-04 — destino AL antes/perto da virada ⇒ mesmo mecanismo de override do operador que MS.
- Owner decisions still open: None — as 5 decisões do brief §4 foram respondidas (P1 + P3).
- Blocked items: None — precondição merge W1 satisfeita.

### Clarification Interview

Registrada em `research/p1-clarified-decisions-2026-07-17.md` (P1b tabela completa). Sem ambiguidade bloqueante restante.

## Architecture Spine

### System Shape

Monolito Go modular (`apps/server_core`, 15+2 módulos governance) + SPA React (`apps/web`) + Postgres + adapters (Oracle/Sankhya, Mercado Livre OAuth). Domínio MPC-native: payloads de provider morrem nos adapters. Novo módulo `erp_import` atrás do port `internal_read/ports.Reader`.

### Runtime Topology

Docker dev stack local (`npm run docker:dev`): server :8080, web :5174, Postgres local. Installation ML conectada preexistente no DB (credenciais AES via `MPC_ENCRYPTION_KEY`). Coleta de sinais on-demand. Compose monta `.:/workspace` — hub re-aponta pós-merge (memória operacional).

### Runtime Contract

- Tenant fixo `tenant_default` (MC_DEFAULT_TENANT_ID); toda query nova escopa `tenant_id`.
- Escrita ML: NENHUMA. Único write path com alvo PROVIDER = envelope M-03 (`/mutations` → poller → adapter), dispatcher provider desabilitado por configuração na demo. Writes de estado LOCAL (resoluções de vínculo, incl. batch M-04) vivem nos módulos donos com auditoria própria — nunca tocam provider.
- OpenAPI (`contracts/api/marketplace-central.openapi.yaml`) + `packages/sdk-runtime` (client TS manual) mudam no MESMO commit.
- Migrações: prefixo 4 dígitos, blocos pré-alocados por milestone, fixture count bumpa junto.
- Fatos operacionais desconhecidos NUNCA viram zero/default (ADR-17); FE usa `UnknownValue`/`FreshnessIndicator` (packages/ui).

### Cross-Cutting Decisions

| Decision | Status | Prevents | Must preserve | Validation impact |
| --- | --- | --- | --- | --- |
| ADR-01 Base = main pós-merge W1 (≥ f4612be3); base SHA 40-hex registrado por chip | satisfeita | retrabalho contra shell antigo; duplicar envelope | M-03 envelope = único write path com alvo PROVIDER (estado LOCAL = módulo dono, ADR-08) | base SHA no ledger de cada chip |
| ADR-02 Fonte ERP dual: módulo `erp_import` (.xlsx → snapshot + adapter do subset Reader) ao lado do Oracle; import em lote com protocolo, hash do arquivo, source, import time, relatório de rejeição por linha; runtime independe do workbook | ratificada | demo refém de VPN Sankhya; stub em live path | IC-02; opcional ausente ⇒ estado honesto | import real de planilha exemplo; linha inválida rejeitada com motivo |
| ADR-03 Identidade: CODPROD=SKU canônico, REFERENCIA=EAN/GTIN (checksum+unicidade), REFFORN=fabricante; seller_sku resolve SÓ p/ CODPROD; corrigir reader Oracle | ratificada | colisões EAN viram vínculos errados | gate 2 âncoras; contradição vence EAN; hard negatives; sem auto-ACCEPT por título | fixtures de colisão (Doka/Menegotti, Doka/VW) |
| ADR-04 Evidência de preço persiste no módulo `market` (Snapshot/Signal/ValidatedOffer/Aggregate); pricing consome contrato IC-03, não persiste market_price próprio | ratificada | "market_price" sem origem/evidência; tabelas duplicadas | ADR-17; fonte/fetched_at/idade/n_offers/n_sellers em toda leitura | teste negativo: falha de coleta não zera snapshot |
| ADR-05 Adapter ML: M-02 dono ÚNICO de TODAS as extensões (sale_price, price_to_win, products search/detail, products/{id}/items FLAG, shipments, shipping_options); consumidores usam ports normalizados IC-06 | ratificada | edição concorrente do adapter; DTO ML vazando | flag OFF default + paginação completa + telemetria + fallback; NO-GOs research §3 | lane live-provider-read; telemetria da rota flag visível |
| ADR-06 Veredicto ≠ identidade: `match_status` ACCEPT/REVIEW/REJECT/NO_CANDIDATE (IC-01) + `blocking_state` NO_CANDIDATE/NO_PRICE_EVIDENCE/INSUFFICIENT_MARKET (<5 sellers válidos)/SEM_CUSTO (evidência de mercado sem custo ERP) — enums separados, semântica: IC-03; pós-anúncio usa sale_price+price_to_win; copy nunca promete preço automático pré-anúncio | ratificada | reabrir NO-GOs; buy_box null virar zero | evidência de coleta visível em toda UI de preço | QA renderiza estados negativos; nunca R$0/verde enganoso |
| ADR-07 FE retheme-first: M-03 entrega tokens papel+verde, data-theme light/dark, fontes, pills canônicas (HANDOFF; Mercado/Repasses "em breve"; Vínculos fora da nav), indireção de rotas por área; telas novas em `apps/web/src/pages/<área>/` | ratificada | dois visuais; colisão no seam AppRouter/Layout | decisões HANDOFF (sem inline edit, kanban read-only, config 3 camadas); defeitos de mock não se reproduzem | QA visual light/dark por tela |
| ADR-08 Zero writes ML: toda mutação com alvo PROVIDER via fila M-03 preview+protocolo (dispatcher provider OFF por config na demo); estado local (ex.: batch de vínculos M-04) aplicado no módulo dono com preview+auditoria, fora do envelope — exceção P5-F-12 | ratificada | write acidental na frente do cliente | 7 gates provider-write intactos p/ MIS-005 | prova de rede/audit: nenhuma escrita saiu |
| ADR-09 Polling/GET only | ratificada | endpoint público/tunnel em 3 dias | idade da coleta visível; falha de refresh = estado stale honesto | refresh visível na UI |
| ADR-10 DIFAL fonte única no `pricing`: seed 27 UFs "padrão 2026" + overrides esparsos; Simulador aplica destino REAL; Pedidos consome chip read-only; rotular "seed — não é orientação fiscal" | ratificada | 3 implementações divergentes; SP hardcoded | interna−interestadual; origem SC; 12% MG/PR/RJ/RS/SC/SP, 7% demais (lista exata: IC-04); taxa desconhecida ⇒ unknown, não 0% | UF de exceção ajustada reflete em Simulador E Pedidos |
| ADR-11 Coleta on-demand + runtime docker local | ratificada | infra nova no prazo | sem retroativo: "sem histórico ainda" honesto | rehearsal completo no stack local |
| ADR-12 sdk-runtime manual: OpenAPI+SDK mesmo commit; arquivo de client por milestone + barrel hub-adjudicado | ratificada | "regenerar" fantasma; colisão no client | parity OpenAPI↔SDK↔handler | parity check no smoke L2 |

Accepted trade-offs: ADR-02 caminho de ingestão precisa endurecer lifecycle no MIS-005 · ADR-05 M-07/M-08 esperam ports do M-02 · ADR-06 demo pode exibir NO_PRICE_EVIDENCE em produtos do cliente (mitigação no runbook) · ADR-07 telas esperam shell (~½ wave A) · ADR-08 "aplicar" não conclui de verdade (protocolo "na fila") · ADR-09 dado "ao vivo" exige refresh · ADR-10 edição DIFAL só via drawer Simulador · ADR-12 sync manual propenso a drift, mitigado por parity L2.

## Shared Contracts

| Contract | Boundary | Path | Why it exists |
| --- | --- | --- | --- |
| IC-01 identidade & matching | catalog ↔ erp_import ↔ product_links ↔ market ↔ FE | `research/identity-matching-interface-contract.md` | dois workers não podem divergir em CODPROD/EAN/estados |
| IC-02 import .xlsx ERP | arquivo ↔ erp_import ↔ Reader port | `research/erp-xlsx-import-interface-contract.md` | contrato de colunas + degradação honesta + protocolo |
| IC-03 leitura de evidência de mercado | market ↔ pricing ↔ FE (M-05/M-06/M-07) | `research/market-evidence-read-interface-contract.md` | campos de evidência obrigatórios em toda UI de preço |
| IC-04 cálculo & DIFAL | pricing ↔ FE Simulador ↔ orders | `research/pricing-difal-interface-contract.md` | decomposição única; tabela UF única |
| IC-05 seams FE | shell M-03 ↔ telas M-04..M-09 | `research/fe-shell-seams-interface-contract.md` | tokens/nav/rotas/primitivas sem colisão |
| IC-06 ports de leitura ML | connectors(M-02) ↔ market/pricing/orders | `research/ml-read-ports-interface-contract.md` | shapes normalizados; DTO morre no adapter |

## Milestone Strategy

| ID | Name | System change | Why this order | Path |
| --- | --- | --- | --- | --- |
| M-01 | erp-xlsx-identity | identidade corrigida + módulo erp_import + Reader adapter | fundação de dados; wave A | `M-01-erp-xlsx-identity/` |
| M-02 | price-intel-core | persistência market + extensões adapter ML + resolver + veredicto | fundação de evidência; wave A | `M-02-price-intel-core/` |
| M-03 | shell-retheme | tokens/tema/fontes/pills canônicas/indireção de rotas/primitivas | todas as telas herdam; wave A | `M-03-shell-retheme/` |
| M-04 | vinculos-import-ui | gaps product_links + tela Vínculos & Importação | consome M-01+M-03; wave B | `M-04-vinculos-import-ui/` |
| M-05 | anuncios-sinais | sinais competitivos em listings + extensão AnunciosPage | consome M-02+M-03; wave B | `M-05-anuncios-sinais/` |
| M-06 | produto-detalhe | header+veredicto+abas Anúncios vinculados/Estoque | consome M-01/M-02/M-04; wave C | `M-06-produto-detalhe/` |
| M-07 | simulador | serviço de cálculo real + DIFAL seed + tela Simulador | consome M-01+M-02; wave B | `M-07-simulador/` |
| M-08 | pedidos | projeção orders + tela Pedidos completa | consome M-02+M-07(read); wave B | `M-08-pedidos/` |
| M-09 | dashboard-demo | Dashboard MPC de agregações reais — CORTÁVEL | última; só com jornada central verde | `M-09-dashboard-demo/` |

## Parallel Execution Plan

### Dependency DAG

- Wave A: `M-01 ∥ M-02 ∥ M-03` (superfícies disjuntas; sem edges entre si — M-02 resolver depende do IC-01 RATIFICADO, artefato de planning, não de código M-01).
- Wave B (após A): `M-04 ∥ M-05 ∥ M-07 ∥ M-08`.
- Wave C: `M-06 ∥ M-09` (M-09 só inicia com jornada central verde; cortável).
- Edges (artefato forçante): M-01→M-04 (produtos importados via `FindProductsForLinking`/IC-02) · M-01→M-07 (custo/estoque via Reader/IC-02) · M-01→M-06 (header produto/estoque) · **M-01→M-08** (custo por pedido via `GetCostAsOf`/IC-02 — edge de dado) · M-02→M-05 (port de leitura market IC-03, publicado por M-02 F-04) · M-02→M-07 (ports comissão/frete IC-06) · M-02→M-08 (ports shipments IC-06) · M-02→M-06 (veredicto/sinais IC-03) · M-03→{M-04,M-05,M-06,M-07,M-08,M-09} (tokens+primitivas+indireção de rotas IC-05) · M-04→M-06 (aba Anúncios vinculados consome links resolvidos) · **M-05→M-06** (`listings.ts` + campos de sinal por anúncio — aba Anúncios vinculados) · M-07→M-08 (decomposição/DIFAL read IC-04) · **{M-01, M-03, M-04, M-05, M-08}→M-09** (fontes reais das agregações: erp_import último-import, tokens/rotas, product_links pendentes, listings summary, orders summary — M-02/M-06/M-07 NÃO produzem artefato consumido pelo M-09).
- **Ordem intra-wave B (P5-F-01)**: M-07∥M-08 vale para o trabalho independente de M-08 (projeção orders, enriquecimento shipment via IC-06, tela). O trecho de M-08 F-01 que consome `Decompose`/`DifalForUF` SÓ inicia após M-07 F-01 publicar os ports (assinatura congelada IC-04) — gate operado pelo hub no dispatch/merge; edge de código `M-07 F-01 → M-08 F-01 (trecho decomposição)`.
- M-08 FE pode começar contract-first após M-03; close espera M-02/M-07 (dados reais). Mesma regra p/ M-05/M-06/M-07-FE.

### Ownership matrix

| Milestone | Files/modules (exclusive) | OpenAPI sections | Migration block | FE surface | DB tables/shape | Shared-seam locks predicted |
| --- | --- | --- | --- | --- | --- | --- |
| M-01 | `modules/erp_import/**` (novo), `modules/internal_read/**` (semântica REFERENCIA), `modules/catalog/**` (identity), `sdk-runtime/src/erpImport.ts` | `/erp/imports*` + schema de produto do catalog (aditivo: ean/refforn/marca/ncm) | **0045–0049** (em `apps/server_core/migrations/`) | none | `erp_import_*` (snapshot produtos/estoque/custo, import_protocols) | governance modules.json (+erp_import, via merge); composition root (registro próprio); **grant aditivo: tipos catalog em `sdk-runtime/src/index.ts` (campos identity, ADR-12 OpenAPI+SDK mesmo commit; hub adjudica)** |
| M-02 | `modules/market/**`, `modules/connectors/**` (capability_adapter ML + ports), `sdk-runtime/src/market.ts` | `/market/*` | **0050–0054** (em `apps/server_core/migrations/`) | none | `market_price_snapshots`, `market_validated_offers`, `market_aggregates`, `market_competitive_signals`, `market_match_decisions` | composition root (registro próprio); publica port Go de leitura `market.EvidenceReader` (IC-03) consumido por M-05 |
| M-03 | `apps/web/src/app/**` (Layout/AppRouter/theme), `apps/web/src/routes/**` (novo, indireção), `packages/ui/src/**` (primitivas novas), tokens/fontes/index.css | none | none | shell global + nav | none | nenhum |
| M-04 | `modules/product_links/**`, `sdk-runtime/src/productLinks.ts`, `apps/web/src/pages/vinculos/**` | `/product-links/*` | **0065–0069** (em `apps/server_core/migrations/`) | `/vinculos` (via arquivo de rota próprio) | ALTERs `product_links_*` (confiança/motivo/lote/undo/auditoria) | nenhum (batch preview/apply local DENTRO de product_links — P5-F-12: envelope /mutations NÃO comporta intents por item de candidato e o stub de writes-off desativa o linkage writer; mutations fora do MVP p/ vínculos) |
| M-05 | `modules/listings/**` (sinais), `sdk-runtime/src/listings.ts` (aditivo), `apps/web/src/pages/` Anuncios*/Listing* existentes | `/listings*` (aditivo) | none | `/anuncios` | none (join read) | nenhum |
| M-06 | `apps/web/src/pages/produto/**` | none | none | `/catalogo/produtos/:productId` | none | leitura composta client-side (catalog/inventory/listings/market via SDK) |
| M-07 | `modules/pricing/**`, `sdk-runtime/src/pricing.ts`, `apps/web/src/pages/precos/**` + `packages/feature-simulator/**` | `/pricing/*` | **0055–0059** (em `apps/server_core/migrations/`) | `/precos` | `pricing_calc_profiles`, `pricing_difal_rates`, `pricing_scenarios` | fila M-03 (consome /mutations, tipo `price_update`, teto `previewed`); publica ports `Decompose`/`DifalForUF` (IC-04, assinatura congelada) consumidos por M-08 |
| M-08 | `modules/orders/**`, `sdk-runtime/src/orders.ts` (aditivo), `apps/web/src/pages/pedidos/**` | `/orders*` (aditivo) | **0060–0064** (em `apps/server_core/migrations/`) | `/pedidos` | ALTERs/projeções `orders_*` | consumo read-only dos ports M-07 (IC-04) e M-02 (IC-06) |
| M-09 | `modules/dashboard/**`, `sdk-runtime/src/dashboard.ts`, `apps/web/src/pages/dashboard/**` (rebuild) | `/dashboard/*` (aditivo) | none | `/` | none (agregação read) | nenhum |

Reserva integração: **0070–0074** (só correção aprovada pelo hub, nunca auto-atribuída). Seams de merge do hub: `contracts/governance/modules.json`, composition root `root.go` (cada chip adiciona registro próprio; conflitos = hub), barrel `sdk-runtime/src/index.ts` (export aditivo 1 linha por milestone), swap de placeholder em arquivo de rota da área (IC-05). Toda migração nova bumpa fixture `runner_test.go` (conflito trivial adjudicado no merge).

## Quality Attributes

| Attribute | Target (concrete) | Owner (ADR/seam) | Validation criterion |
| --- | --- | --- | --- |
| Observabilidade de coleta | fonte, fetched_at, idade, n_offers/n_sellers visíveis em TODA UI de preço; telemetria na rota flag products/{id}/items | ADR-04/ADR-06/IC-03 | QA inspeciona cada tela de preço; log de telemetria da rota flag |
| Durabilidade ADR-17 | snapshot válido nunca sobrescrito por zero/falha | ADR-04/IC-03 | teste negativo no VC do M-02 |
| Security (baseline) | sem superfície nova de auth; tenancy explícito; secrets só env; docs/env nunca commitados | profile §7 | review de PR por chip; sem .env em diff |
| Maintainability (baseline) | L0–L2 verdes por chip | profile §5 | ladder por feature/milestone |

## Non-Functional Scope

| Declined attribute | Reason |
| --- | --- |
| Performance (p95) | volume single-seller, demo local |
| A11y (WCAG) | herda design system; prazo |

## Validation Strategy

Doutrina completa por milestone (P1b): dual gate frio (Opus + Sol medium, concordância) + QA live-drive browser fresh. Fechamento da missão: QA live-drive integrado = rehearsal completo da demo no stack local limpo (import xlsx real → vínculos → sinais → simulador → pedidos), prova zero-writes (log/telemetria adapter sem PUT/POST ML), runbook executado. Contracts em `validation-contract.md` (mission) + por milestone.

## Risks

| id | risk | likelihood | impact | mitigation | trigger | owner |
| --- | --- | --- | --- | --- | --- | --- |
| R1 | 3 dias vs 9 milestones full-doctrine | H | H | waves paralelas; M-09 cortável; M-06 abas mínimas; gates são horas; contratos congelados cedo | wave A não fecha até sexta à noite | hub |
| R2 | NO_PRICE_EVIDENCE em massa nos produtos do cliente (buy_box null 22/22) | H | M | runbook mistura anúncios próprios ativos (sale_price+price_to_win fortes) + pré-anúncio honesto; copy ADR-06; preflight coleta na véspera | preflight véspera devolve <N sinais | operador+hub |
| R3 | Planilha real do cliente foge do IC-02 | M | M | template exemplo + dry-run pré-demo + relatório de rejeição por linha | dry-run rejeita >10% linhas | operador |
| R4 | Permissões/dados da conta ML não suportam leitura necessária | M | H | preflight read-only cedo (wave A, M-02) contra installation real; estados unavailable honestos | preflight falha em rota necessária | M-02 chip |
| R5 | Rota products/{id}/items desliga/instável | L | M | flag OFF default + paginação + telemetria + fallback explícito (ADR-05) | telemetria acusa 4xx/5xx | M-02 chip |
| R6 | Installation ML não conectada no DB da demo (UI OAuth removida no W1) | M | H | runbook: verificar installation antes; conectar via endpoints server se preciso | stack limpo sem installation | operador+hub |
| R7 | Colisão em seams compartilhados (root.go, modules.json, barrel SDK, fixture migrações) | M | M | matriz de ownership; seams de merge = hub; edições aditivas 1-linha | conflito de merge não trivial | hub |
| R8 | DIFAL seed lido como orientação fiscal | M | M | rótulo "seed padrão 2026 — não é orientação fiscal" em toda superfície DIFAL | — | M-07 chip |
| R9 | M-09 consome tempo com jornada central instável | M | M | M-09 só inicia com waves A+B verdes; corte sem reabrir escopo | sábado sem wave B verde | hub |
| R10 | Executor ML de `price_update` já habilitado; "aplicar preço" (M-07) cria protocolo que EXECUTARIA write se aprovado | L | H | teto `previewed` na UI (sem approve), dispatcher gated por config `MPC_PROVIDER_WRITES_ENABLED` (OFF na demo, ADR-08), runbook proíbe approve na demo; trilha de auditoria do protocolo obrigatória (prova MIS-004-C02); habilitação governada = MIS-005 M-10 | protocolo price_update aprovado fora do runbook | M-07 chip + operador |

## Handoff

- Current status: **PLANNED** (2026-07-17). Gate lean redefinido pelo operador pós-cap: r04 Claude targeted (H1 PASS + 1 FAIL reparado) → r05 Sol HIGH targeted (gpt-5.6-sol/high, 4/4 checks PASS, zero advisories, verdict Ready scoped). Fold completo: `readiness-review.md`; Sol verbatim: `planning-reviews/p7-sol-readiness-r05.md`; manifest r05 digest `6ee15d83…`.
- Current owner: hub (harness-hub) — autoria dos chips de milestone.
- Next owner: chips M-01/M-02/M-03 (wave A) per Parallel Execution Plan.
- Next action: hub boot → spawn chips wave A (M-01, M-02, M-03 — paralelos, superfícies disjuntas per matriz).
- Required artifact paths: `M-0X-*/milestone.md`, `M-0X-*/F-0X-*/feature.md`, `validation-contract.md`, ICs em `research/*-interface-contract.md`.
- Required evidence paths: `validation-result.md` (mission), `M-0X-*/validation-result.md`, prova zero-writes no rollup da missão.
- Blocked decisions: None — decisões abertas do brief todas respondidas.
