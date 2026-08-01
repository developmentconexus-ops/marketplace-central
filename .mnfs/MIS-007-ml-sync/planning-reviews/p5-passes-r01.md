# P5 Mandatory Passes — MIS-007 r01

Autor: Mission Strategist (sessão de planning). Data: 2026-07-31.
Insumos: 9 milestone.md + 25 feature.md + mission.md `## Parallel Execution Plan` +
IC-01..IC-07 + `research/p5-prerequisites.md` (investigator, 12 itens, tip dd89d4b3).

## Pass 1 — Feature-level write-DAG

Método: Ownership de cada brief (owned/forbidden/parallel-safe) cruzado par a par dentro
do milestone e entre milestones da mesma lane.

Intra-milestone: todo DAG interno serializa os overlaps declarados
(M-01 F-01∥F-02 disjuntos por arquivo; M-02 F-01→F-02, F-03∥F-04∥F-01 disjuntos;
M-03/M-04 seriais; M-05 F-01∥F-02 disjuntos (channel_fees vs divergences) → F-03;
M-06/M-07/M-08/M-09 seriais). Sem overlap sem edge.

Cross-milestone (todos com edge/lock nomeado na matriz da missão):
1. M-05 → superfícies de listings do M-04: **lock aditivo** registrado (application/
   adapters de listings, additive-only, M-04 fechado quando lane C abre).
2. M-06 → `orders/**` do M-03: **herança de lane** (M-03 fechado; nunca simultâneos).
3. M-07 → região pricing root.go `:828-858`: única edição de bloco existente — **merge
   arbitrado pelo hub** (R-7).
4. M-08 ↔ M-09 → `WebhookStatsReader`: M-09 publica porta + impl nula (F-01); M-08 fornece
   impl real + troca fiação NA REGIÃO DELE (F-02). Assinatura pinada em IC-05. Sem edit
   cruzado de arquivo.
5. OpenAPI + SDK (`index.ts` literal `:2113-2446`): milestones FE-contract = M-03 (aditivo
   /orders), M-05, M-06, M-07, M-09. Lanes A e B têm ≤1 cada (M-09; M-03). Lane C tem 3 →
   **hub serializa COMMITS de contrato dentro da lane** (ADR-14, registrado na matriz).
6. root.go: 1 linha ancorada por milestone (M-03/M-04/M-06/M-08/M-09), região própria;
   exceção M-07 (item 3).
7. Allowlist do guard (M-02 F-04): encolhida por M-03 (A/B) e M-07 (C/D) — lanes
   diferentes (B e C), sem janela simultânea; shrink-only (ADR-05).

Veredito: PASS — nenhum write-set overlap sem edge serial ou lock registrado.

## Pass 2 — Contract-satisfiability

Promessas de brief checadas contra o conjunto ADR/IC ratificado:

1. M-05 F-03 `filter.divergentes=true` + cursor de paginação existente
   (`DecodeListingCursor`, sort `l.title,l.provider_listing_id,l.variation_id`): filtro é
   predicado WHERE adicional — página única por statement preservada. SATISFIABLE.
   (Armadilha ORDER BY output-column citada no brief.)
2. M-06 decomposição: chaves ausentes ≠ zero (IC-03) × invariante soma-das-partes — braço
   incompleto[] definido: soma só é exigida quando componentes presentes; margem NULL.
   SATISFIABLE.
3. M-07 cascata 2→1→config × "camada 1 sem escritor nesta missão": braço 1 implementado e
   testado por fixture, adormecido — IC-01 permite (resolução é ordem de leitura, não
   exige fonte viva). Registrado no brief como decisão. SATISFIABLE.
4. M-09 payload IC-05 × schema 0075: VERIFICADO nesta rodada
   (`0075_sync_sync_state.sql`) — cursor/last_full/last_incremental/last_error/
   consecutive_failures existem; `phase` deriva do cursor IC-06; `last_success_at` :=
   last_full_sync_at. SATISFIABLE (brief atualizado com colunas verbatim).
5. M-04 F-02 replacement de MASS-CLOSURE × IC-07 PK congelada: absent-marking usa colunas
   novas (last_seen_at/absent_since) sem ALTER de PK. SATISFIABLE.
6. M-08 always-200 × envelope de erro universal (CHIP-ERROR-UNIFY): única exceção 500 de
   INSERT usa envelope vigente; 200 vazio não conflita (endpoint sem contrato SDK, IC-04).
   SATISFIABLE.
7. IC-06 "ports definidos no M-02" × decomposição (nenhum brief do M-02 os cria) —
   **CONFLITO ACHADO E RESOLVIDO NO PLANO** (regra P5: ratificar a resolução, nunca
   empurrar p/ implementação): IC-06 emendado — port mora no módulo dono (`OrderIngestor`
   M-03, `ListingIngestor` M-04); tenant pode vir por scoping de repo (idioma vigente);
   semântica resource-addressed/idempotente continua binding. Edit aplicado no IC nesta
   rodada.

Veredito: PASS (1 conflito achado, resolvido por emenda ratificada no IC-06).

## Pass 3 — Prerequisite-existence

Fonte: `research/p5-prerequisites.md` (12 itens verificados por leitura direta) + 2
verificações desta rodada.

| Assunção de brief | Evidência |
| --- | --- |
| DeriveOrderBucket assinatura/literais | fato #1 (`order_bucket.go:48`) |
| BuyerFiscalInfo DTO + honest-absence 404 | fato #2 (`buyer_fiscal_reader.go:101-103`) |
| /listings envelope/params/sort + colunas FE | fato #3 |
| IntegracoesPage seções + idioma data-fetching | fato #4 |
| ImportService.Import/UpsertOrders baseline | fato #5 |
| Cadeia tarifflive/composite + root.go :837-856 + 0.16 | fato #6 + audit D-120 |
| root.go anchors (orders/listings/schedulers/pricing/batch) | fato #7 |
| ApplyCompletedPull/Ingestion.Pull/AbsorbProviderSnapshots | fato #8 |
| SDK literal :2113-2446 + métodos por módulo | fato #9 |
| JobFunc/RegisterJob/RecordSuccess(incremental=false) | fato #10 |
| Highest migration 0085 → ranges 0086+ livres | fato #11 |
| /sync/runs idioma (referência do M-09) | fato #12 |
| sync_state colunas p/ IC-05 | verificado r01 (`0075_sync_sync_state.sql`) |
| Enum entity aceita listings/orders SEM migração | comentário do 0075 ("SEMANTIC enum... validated in application layer, NOT a DB CHECK") + IC-06 |
| products_mirror/vínculo (M-05 F-02, M-06 F-02 custo) | MIS-006 M-02/0076 + backfill real @080851c0 (memória; tabelas vivas em produção de tela) |
| route class idiom (`route_deadline.go:23-28`) | research MIS-006/CHIP-ERROR (citado M-08 F-01) |

Assunção ABERTA (não-órfã — ambos desfechos especificados no brief):
- M-05 F-01: multiget de items carrega sale_price com `?context=channel_marketplace`?
  Spec do F-01 verifica contra o DTO do M-01 F-02; caminho A (vem no multiget: zero GET
  extra) e caminho B (GET /items/{id}/prices dedicado) ambos escritos no brief. Não é ★5:
  nenhum caminho depende de símbolo inexistente.

Símbolos a-criar: todos nomeados com feature criadora (channel_fees/divergences/
order_shipments → M-02 F-01; IngestOrder → M-03 F-02; IngestListing → M-04 F-03;
WebhookStatsReader → M-09 F-01 porta / M-08 F-02 impl; feeledger resolver → M-07 F-01).

Veredito: PASS — zero assunção não-verificada sem dono.

## Disposição

3/3 passes PASS. Próximo: auditoria de decomposição (touchpoint Sol P5 — quota Codex
esgotada até 2026-08-05; waiver ratificado pelo operador: crew Claude fria substitui,
Sol retroativo obrigatório antes de `status: planned`). Segue despacho de auditor frio
Claude (Opus) contra o checklist do skill; artefato p5-claude-decomposition-audit-r01.md.
