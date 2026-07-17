# P3 Reconciliation — r01 (2026-07-17)

Inputs: `p3-claude-candidate-r01.md` × `p3-sol-counterproposal-r01.md` (Sol BLIND, manifest `p3-input-r01.sha256`, top-digest `c95678eb…`). Regra: dual-model agreement só quando semântica material coincide; diferenças respondíveis por evidência resolvidas aqui com evidência citada; ao STOP sobem só decisões de autoridade do operador.

## 1. Dual-model AGREEMENT (semântica material coincide)

| Tópico | Claude | Sol | Nota |
|---|---|---|---|
| Base execução = main pós-merge W1; nenhum chip parte de cd74b401 | ADR-01 | precondition + Wave 0 | idêntico incl. "não inferir diff dos chips" |
| Identidade: CODPROD canônico, REFERENCIA=EAN (corrigir reader), REFFORN=fabricante; catalog dono; gate 2 âncoras, contradição vence EAN, sem auto-ACCEPT por título | ADR-03 | ADR-02 | idêntico; Sol acrescenta "seller_sku resolve só p/ CODPROD" — absorvido no IC-01 |
| Modelo de preço persiste no módulo `market` existente (research §6: Snapshot/Signal/ValidatedOffer/Aggregate); pricing consome contrato, não persiste market_price próprio | ADR-04 | ADR-04 | idêntico |
| Adapter ML: dono único das extensões; flag+paginação+telemetria+fallback em products/{id}/items; NO-GOs fechados | ADR-05 | ADR-05 | ver Δ3 (escopo shipments) |
| Estados honestos ACCEPT/REVIEW/REJECT/NO_CANDIDATE/NO_PRICE_EVIDENCE/INSUFFICIENT_MARKET; veredicto ≠ identidade; buy_box null nunca vira zero | ADR-06 | ADR-06 | Sol explicita "<5 sellers válidos = INSUFFICIENT_MARKET" (research §7) — absorvido |
| ADR-17: snapshot válido nunca sobrescrito por zero/falha; observed/fetched/expiry distinguíveis | ADR-04 | ADR-07 | teste negativo obrigatório em ambos |
| Retheme-first; nav canônica; Mercado/Repasses "em breve"; Vínculos fora da nav; defeitos de mock não se reproduzem | ADR-07 | ADR-08 | idêntico |
| Zero writes ML; toda mutação via envelope M-03 preview+protocolo; execução provider inalcançável na demo; prova de rede/audit no QA | ADR-08 | ADR-09 | idêntico |
| Polling/GET only; sem webhooks/scheduler no MIS-004 | ADR-09 | ADR-10 | idêntico |
| DIFAL fonte única no `pricing` (seed 27 UFs + overrides), destino real, toggle no Simulador, chip read-only em Pedidos; sem agendar/pagar | ADR-10 | ADR-11 | idêntico |
| Coleta on-demand; runtime = docker dev stack local | ADR-11 | ADR-12 | idêntico |
| sdk-runtime é manual; OpenAPI+SDK mesmo commit | ADR-12 | ADR-13 | ver Δ2 (mecânica do lock) |

## 2. Diferenças resolvidas por evidência (sem autoridade do operador)

### Δ1 — Onde vive o import .xlsx: `erp_import` (Claude) vs workflow do `product_links` (Sol)
**Resolução: Claude ADR-02 mantido (módulo `erp_import` implementando subset do Reader port).**
Evidência: baseline — `internal_read` É o boundary ERP com port `ports/reader.go:48-55` e injeção `Unavailable*` (root.go:327-362); `product_links` possui workflows de VÍNCULO anúncio-ML→produto (imports/candidates/resolutions), não master-data ERP. Colocar planilha ERP dentro de product_links mistura domínios e quebra o seam existente; adapter novo no port existente preserva. **Absorvido do Sol:** hash do arquivo + source + import time persistidos; runtime não depende do workbook montado; rejeição por linha inspecionável (entra no IC-02).

### Δ2 — Contract lock: milestone dedicado dono de TODO OpenAPI+SDK (Sol M004-01) vs seções disjuntas por milestone sob lock do hub (Claude)
**Resolução: seções OpenAPI disjuntas por milestone (Claude), reforçado com pré-atribuição de arquivos no sdk-runtime (do Sol).**
Evidência: harness profile §5 instancia eixo "contract artifacts" com lock aditivo — mecanismo existe exatamente p/ seções disjuntas em paralelo; milestone-contrato único serializa a autoria de TODOS os schemas antes de qualquer consumidor (custo serial fatal em 3 dias) e produz contrato especulativo sem feedback de implementação. Risco real que Sol aponta (colisão no client manual) fecha assim: cada milestone possui SEU arquivo de client no sdk-runtime; barrel/index = seam aditivo hub-adjudicado. Vai ao Parallel Execution Plan.

### Δ3 — Shipments no adapter ML: additive lock p/ M-08 (Claude) vs dono único total M-02 (Sol)
**Resolução: adotado Sol.** M-02 price-intel-core possui TODAS as extensões do capability_adapter (incl. /shipments/*, shipping_options); M-08 consome ports normalizados. Remove a contingência de lock aditivo (retrospectiva M-01: locks descobertos mid-flight custam escalação). Custo: M-08 espera ports do M-02 — já coberto pelo edge M-02→M-08 existente.

### Δ4 — Blocos de migração: por wave (Claude) vs por milestone (Sol)
**Resolução: adotado Sol (defeito no candidato Claude).** Wave A tinha M-01 E M-02 precisando migração no MESMO bloco 0045-0049 = colisão interna. Novo: M-01: 0045-0049 · M-02: 0050-0054 · M-07: 0055-0059 · M-08: 0060-0064 (se precisar projeção) · 0065-0069 reserva integração (só correção aprovada). Toda migração bumpa o literal em `runner_test.go:25-27`.

### Δ5 — Milestone de integração dedicado + composition root exclusivo na integração (Sol M004-11)
**Resolução: parcial.** (a) Wiring exclusivo do composition root na integração — REJEITADO: prática harness comprovada = cada chip registra seu módulo no root ao mergear (senão o chip não roda a própria lane de integração; memória W1: registro governance aterrissa via merge do chip). Root/modules.json continuam seams de merge do hub. (b) Rehearsal/runbook NÃO podem morar num milestone cortável — ACEITO (defeito no candidato: estavam dentro do M-09 cortável). Novo: rehearsal + runbook + prova zero-writes = conteúdo do QA live-drive de FECHAMENTO DA MISSÃO (obrigatório, incortável); M-09 fica só dashboard, cortável.

### Δ6 — MIS-005: 9 headlines (Claude) vs 12 (Sol)
**Resolução: união = 11.** Sol expôs 2 lacunas reais no candidato: habilitação de writes de produção (Claude ADR-08 prometia "gated pós-demo" sem milestone dono) e deployment/hardening de produção. Adicionados M-10 writes-producao e M-11 producao-hardening. Inventory/fulfillment depth (Sol M005-11) funde no M-08 full-visits-benchmarks. Restante mapeia 1:1.

## 3. Diferença que sobe ao STOP (autoridade do operador — é a decisão 5 do brief)

**Grão MIS-004: 9 milestones full-stack por tela (Claude) vs 11 com BE/FE separados + contract-milestone (Sol).**
Recomendação: **9 full-stack**. Fundamento: cada close = dual gate (Opus+Sol) + QA live-drive — 11 closes > 9 closes em custo de gate no prazo de 3 dias; no próprio plano do Sol os milestones FE "podem começar" mas só FECHAM contra o backend real, ou seja o split dobra gates sem destravar paralelismo além do que as waves já dão; vertical slice por tela alinha rota FE + seção OpenAPI + módulo no MESMO dono (colisão zero por construção). Estrutura Sol fica registrada como alternativa.

## 4. Plano reconciliado (delta sobre o candidato)

- ADRs 1–12 mantidos com absorções: IC-01 += seller_sku→CODPROD only; ADR-06 += INSUFFICIENT_MARKET <5 sellers; ADR-02 += hash/source/import-time + independência do workbook; ADR-05 = dono único TOTAL incl. shipments (Δ3).
- Milestones M-01…M-09 mantidos; M-08 perde lock aditivo (consome ports); M-09 = dashboard only (cortável); rehearsal/runbook/zero-writes-proof → QA de fechamento da missão.
- Migrações por milestone (Δ4). sdk-runtime: arquivo por milestone + barrel hub (Δ2).
- MIS-005 = 11 headlines (Δ6).
- Riscos absorvidos do Sol: preflight de permissões da conta ML cedo (runbook); dry-run da planilha real do cliente pré-demo; rotulagem DIFAL "seed padrão 2026 — não é orientação fiscal"; guarda anti-time-sink no M-09 (só inicia com jornada central verde).

**Verdict: dual-model agreement no spine (13/13 tópicos após absorções); 1 item de autoridade do operador (grão) sobe ao STOP.**

## 5. Addendum pós-STOP — verificação W1 mergeado (pedido do operador, 2026-07-17)

Operador pediu no STOP: verificar alinhamento do M-02 (W1) recém-mergeado com o plano. Evidência: `research/w1-merge-addendum-2026-07-17.md` (R-03). R-01 e manifest P3 intactos — round P3 permanece válido sobre seus inputs congelados; R-03 entra como evidência nova para P4/P5.

**Resultado: spine e milestones INALTERADOS; nenhuma conclusão P3 derrubada.** Deltas absorvidos (nível de superfície, entram no mission.md/briefs em P4/P5):

- ADR-01: precondição merge W1 SATISFEITA (m-02 @ 79d6787f, m-03 @ f4612be3). Base chips = main ≥ f4612be3.
- Blocos de migração 0045+ confirmados livres (M-03 usou gaps 0038/0039; topo 0044; fixture 41).
- M-03 shell-retheme: rotas PT-BR + placeholders + LegacyRedirects + InstallationContext JÁ existem (W1). Escopo restante: tokens papel+verde, data-theme, fontes, sidebar→pills canônicas (HANDOFF), Vínculos fora da nav, indireção de rotas por área (telas trocam placeholder sem tocar seam AppRouter).
- M-04/M-08: páginas legacy DELETADAS no W1 (feature-orders/product-links/marketplaces/integrations) — telas novas em `apps/web/src/pages/<área>/` (padrão M-02 W1), não retheme.
- M-05: estende AnunciosPage existente. M-06: rota placeholder já registrada. M-07: superfície = /precos + packages/feature-simulator (sobrevivente).
- Risco novo no runbook: UI de conexão OAuth removida (/integracoes placeholder) — installation ML deve estar conectada no DB antes da demo (server OAuth endpoints continuam).
