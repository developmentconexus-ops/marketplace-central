# MIS-006-integracao-fundacao

```yaml
id: MIS-006
type: mission
status: planned
owner: Mission Strategist
parent: none
created: 2026-07-20
updated: 2026-07-20
validation_level: QA-0
lifecycle_scope: mission
planning_phase: readiness
readiness: Ready (joint fold, readiness-review.md; Sol rebound per codex-quota-wall, operator-authorized)
base_sha: 138aac3dff20438d8ddc509daf4171d82e5e45f6
```

## Origin

Demo cliente D-120 (2026-07-20) reprovou: sistema julgado não-funcional. Postmortem
(`docs/audit/D120-POSTMORTEM.md`) + auditoria file:line (`memory/main-code-audit-d120.md`)
provaram que a BASE de integração não existe de fato:

- Oportunidades = 100% cliente, sem backend (`oportunidades.ts:64-81`); escopo "monitorado"
  inexistente; exclusão "não vendemos" ausente.
- **CHICKEN-EGG:** `Collect(codprod)` resolve mercado via `LinkedListings(codprod)`
  (`collection_pipeline_service.go:269`) → produto SEM vínculo não tem caminho de mercado. Por
  isso xlsx nunca apareceu em Oportunidades mesmo com dados.
- Sem `products_mirror` canônico → sem JOIN SQL produto×listing×pedido; consumidores leem
  fontes divergentes.
- Import xlsx = snapshot por protocolo, sem estado corrente consolidado; troca de source cega
  (`LatestCompletedSnapshot` max-imported_at).

MIS-005-produto-completo pressupõe base MIS-004 funcionando — premissa falsa. Esta missão é a
FUNDAÇÃO que faltava: torna a integração de produto real antes de qualquer sync ML avançado.

## Objective

Fundação de integração de PRODUTO: xlsx e Sankhya como dois adapters equivalentes do mesmo port,
convergindo num espelho canônico (`products_mirror`) que habilita JOIN SQL e mata o chicken-egg —
todo produto do ERP tem caminho para vínculo e mercado, independente de já estar anunciado.

## Outcome

- `products_mirror` canônico (contrato E2.1), alimentado por XlsxAdapter e SankhyaAdapter via um
  único `ProductSourceAdapter` port; upsert-merge keep-absent (F-XLSX-1), NULL=honesto-desconhecido.
- Fonte ativa por tenant no banco (mata `MC_ERP_SOURCE` env / seleção cega).
- Cadeia de vínculo automática: mirror × listings por EAN → candidato → EAN-exato-único
  auto-aprovado + audit.
- `sync_state` + esqueleto de scheduler (fundação p/ todo sync futuro; nada de ML-sync pesado
  ainda).
- Tela /importacoes mostra a cadeia real: N importados → N vinculados → N coletas **enfileiradas**.
- Chicken-egg quebrado: produto sem vínculo tem rota de descoberta (F3.7
  `descobrir_produto_catalogo(ean)` — **UNPROVEN, exige rodada live T13-T16, ver Decisão F3.7**).

## Scope (SÓ isto)

Produto / ERP / xlsx. Especificamente:
- `products_mirror` + migração + upsert-merge + stale flag.
- `ProductSourceAdapter` port + XlsxAdapter (reestrutura sobre o parser leniente existente) +
  SankhyaAdapter (`[TESTAR-SKW]`: mapear TGFPRO/TGFEST/TGFCUS na sessão do especialista Oracle
  ANTES de codar).
- Config fonte-ativa por tenant.
- Cadeia de vínculo automática (mirror×listings EAN) + auto-aprovação EAN-exato-único + audit.
- `sync_state` + scheduler esqueleto (só o suficiente p/ rastrear import/rebuild/vínculo/coleta).
- Tela /importacoes com a cadeia visível + tela /integracoes fonte-ativa.
- Rota de descoberta F3.7 (se T13-T16 provar viável) OU decisão de removê-la (honest-unknown).

## Non-Scope (missões seguintes)

- Sync de anúncios (colunas E3, backfill scan/multiget) — missão ML-sync.
- Sync de pedidos + decomposição + Fila/SLA — missão ML-sync.
- Mercado/concorrência + tarifas live + `ml_tariffs` — missão mercado.
- **Execução** da coleta de mercado (competitor_offers, market_aggregates) — missão mercado.
  MIS-006 só ENFILEIRA (E8 sync_state), não executa.
- Onboarding saga completa — depois que os adapters existirem.
- Writes ML, webhooks, auth/multi-tenant, repasses — MIS-005.
- Os 3 buracos do adapter (F-ADAPTER-1: backoff/429, tarifa live, Raw DTO) → ver ADR-06.

## Domain Scope (P1a)

Incluído (`lean-core` + fundação): espelho canônico de produto, port+adapters, config fonte-ativa,
cadeia de vínculo automática, rastreio de sync (skeleton), observabilidade de import (telas).
Excluído (com razão) — ver Non-Scope: qualquer sync ML pesado, execução de coleta, tarifas,
onboarding saga, writes/webhooks/auth.

## Clarified Decisions (interview + forced-assumption ledger)

O scaffold já pré-decidiu escopo (operador autor). O planning resolveu as decisões pendentes
via evidência (não houve STOP interativo — decisões abaixo são reversíveis ou já ratificadas em
`INTEGRATION-DATA-CONTRACT.md §3`). Owner-authority pendente: só F3.7 (exige live + hub).

Accepted assumptions (reversíveis):
- Migração sync_state (bloco A) < mirror (bloco B), Fase-0-first — tabelas independentes, só o
  número ordena (reversível: hub aloca blocos).
- root.go editado por additive-lock (M-02 dono, M-01 append) — reversível para serial se hub
  julgar risco de merge alto.
- `/importacoes` vira rota própria + `ImportacaoSection` promovido a módulo compartilhado
  (hoje duplicado em /vinculos e /integracoes) — reversível: alternativa = manter dentro de
  /integracoes (registrado como open question de M-06, não bloqueia decomposição).

## ADRs (cross-cutting, ≥2 milestones)

### ADR-01: products_mirror = estado corrente materializado; snapshots viram history
- Decisão: `products_mirror` (tabela) é a fonte de JOIN SQL. `erp_import_products` continua só
  como history imutável por protocolo. Leituras migram do rescan de snapshot para o mirror.
- Prevents: padrão "read = recompute" (`internalread/reader.go:84-107`) que impede JOIN e
  viola SYSTEM-BLUEPRINT (enrichment no ingest).
- Must preserve: `readports.Reader` shape (consumidores inalterados); history por protocolo.
- Trade-off: dupla escrita (history + mirror) no ingest; aceito (custo O(1) no ingest << O(n) rescan no read).
- Validation impact: mirror populado + JOIN SQL provado com IO pairs A/B.

### ADR-02: ProductSourceAdapter port com read-side preservado + write/Sync-side novo
- Decisão: formalizar o port de-facto (`readports.Reader`) adicionando `Sync()` + `Kind()`;
  xlsx=upload_snapshot, sankhya=live_read_through-com-snapshot-no-mirror.
- Prevents: dois readers divergentes independentes; escolha "priorizar xlsx vs Sankhya".
- Must preserve: interface de leitura dos consumidores.
- Trade-off: none (refactor aditivo).
- Validation impact: ambos adapters escrevem no MESMO mirror; consumidor não sabe a origem.

### ADR-03: fonte ativa = config em banco por tenant; MC_ERP_SOURCE morre
- Decisão: `active_source` (E9) em banco por tenant; ambos adapters construtíveis; resolução
  por-request. `erpSource()` boot-env (`root.go:772`) REMOVIDO.
- Prevents: backend inteiro travado numa fonte no boot (D120 I1); dois tenants incompatíveis.
- Must preserve: fail-closed `ErrUnknownActiveSource`; **`active_source` na chave de cache**
  (lição `chip-import-fix-closed` — poluição cross-source).
- Trade-off: config em banco + resolução por-request (leve).
- Validation impact: toggle por tenant provado nas 2 direções + cache não vaza snapshot errado.

### ADR-04: upsert-merge keep-absent (F-XLSX-1), nunca delete físico
- Decisão: rebuild de mirror faz merge; ausente → `absent_in_last_snapshot=true` + `stale_since`.
- Prevents: cascata de wipe destruindo `product_links`; perda silenciosa.
- Must preserve: ADR-17 (NULL honesto, nunca 0). 
- Trade-off: linhas stale acumulam (mitigado por flag + filtro na UI).
- Validation impact: produto some do snapshot → mirror mantém row flagada (IO Fase 5).

### ADR-05: auto-aprovar vínculo em âncora não-ambígua (CODPROD e/ou EAN)
`AMENDADO 2026-07-25 (D-121) · RATIFIED-BY-OPERATOR` — supersede a redação original
("só EAN-exato-ÚNICO; reusar `validEANCounts`"), que contradizia o motor de match já
ratificado (IC-01 A2) e assumia uma semântica errada de `seller_sku`. Entrevista de
alinhamento com o operador, evidência real do mirror Sankhya (10529 produtos, 7361 com EAN,
91 EANs colidentes cobrindo 188 produtos).

- **Semântica de `seller_sku` (D-121-1):** no Mercado Livre o SKU do anúncio É o CODPROD do
  ERP — o operador cadastra assim, e todos os anúncios já carregam o código. O matcher hoje
  compara `seller_sku` com `p.REFFORN` (código do FABRICANTE, ex. `L.87.22`) —
  `internal_read/adapters/oracle/reader.go:451`. Isso está errado para este tenant e passa a
  casar `p.CODPROD`. REFFORN sai como âncora de SKU (decisão do operador: "esquece REFFORN").
  CODPROD é a chave primária do ERP: não colide por definição, e é a âncora mais forte
  disponível — mais forte que o EAN, que é digitado pelo vendedor e tem 91 colisões reais.
- **Decisão (política de auto-approve):** auto-aprova quando a âncora é não-ambígua —
  (a) CODPROD resolve para 1 produto, com ou sem EAN presente; (b) EAN resolve para 1
  produto; (c) ambos concordando no mesmo produto. Qualquer outro caminho é REVIEW.
- **Conflito CODPROD≠EAN → REVIEW** (nunca auto-resolve por precedência): sinais
  contraditórios são sintoma de cadastro errado no anúncio; auto-resolver esconderia isso.
- **Contradição de título vence tudo:** `detectHardNegative` (kit/combo/cor/voltagem
  divergente entre título do anúncio e nome do produto) rebaixa a REJECT/REVIEW mesmo com
  CODPROD ou EAN batendo. Caso motivador: anúncio-kit cadastrado com o EAN da peça avulsa →
  vínculo silenciosamente errado que distorce estoque e margem.
- **Um produto ↔ N anúncios:** sem limite e sem sinalização. Mesmo codprod anunciado várias
  vezes é operação normal.
- Prevents: auto-aprovação em EAN ambíguo (91 EANs, 188 produtos); vínculo de anúncio-kit a
  peça avulsa; automatizar o sinal fraco enquanto o forte espera na fila manual.
- Must preserve: override manual do operador vence auto-aprovação e nunca é revertido pelo
  automático; título nunca aceita sozinho.
- Trade-off: CODPROD digitado errado no anúncio que caia em outro código VÁLIDO vincula
  errado sem revisão (aceito — o operador é dono do cadastro nos dois lados).
- Validation impact: IO A auto-aprova; EAN colidente, conflito CODPROD≠EAN e hard-negative
  ficam REVIEW.

**Defeitos do plano original corrigidos junto (achados na entrevista, não são decisão de
negócio):**
- *A8*: o plano pedia unique em `(internal_product_id, provider_listing_id)`. Essa coluna não
  existe em `product_links` (é `provider_item_id`), e a PK atual
  `(tenant_id, installation_id, provider_item_id, provider_variation_id)` já garante
  idempotência — e respeita variação, que a chave proposta perderia. A8 = satisfeita pela PK
  existente; nenhum índice novo.
- *`collisions[ean]`*: `len(eanMatches.Products)==1` no gerador JÁ é a contagem de colisão do
  lado ERP. `validEANCounts` mede duplicidade dentro do arquivo xlsx — outro fato. O
  auto-approve usa a contagem do gerador; `collisions_at_decision` em E10 grava esse número.

### ADR-06: 3 buracos do adapter ML (F-ADAPTER-1) fora da fundação
- Decisão: backoff/429 + tarifa live entram na missão que os exercita (backfill em lote).
  `Raw json.RawMessage` = candidato barato (aditivo, 1 assign/unmarshal) — **DEFAULT: diferido**;
  hub pode puxar para um feature isolado (módulo `connectors` não colide com nada). Ver Decisão Raw.
- Prevents: escopo inchar com trabalho que só paga em ML-sync.
- Must preserve: adapter ML maduro intocado (REUSAR, não reescrever).
- Trade-off: risco de perda silenciosa de campo ML até Raw landar (aceito na fundação — sem
  ingest ML pesado ainda).
- Validation impact: n/a nesta missão (não exercita backfill).

## Decisions Resolution

| ID | Decisão | Resolução MIS-006 | Fonte |
|----|---------|-------------------|-------|
| D1 | Arquitetura de fonte | **RATIFICADO** (adapters convergentes, ADR-02) | contrato §3 |
| D2 | Mecanismo incremental | **RATIFICADO** scheduler-first; MIS-006 = skeleton cadence-agnostic | contrato §3 |
| D3 | Campos fiscais obrigatórios | **RESOLVIDO = N/A à E2.** Tabela 10-campos ratificada NÃO tem campo fiscal; fiscal (CPF/CNPJ/IE) vive em E5/Pedido (non-scope). `ncm` fica como passthrough opcional (honesto-desconhecido) pois já existe em `NormalizedRow`. Nenhum campo fiscal novo em products_mirror. | contrato §E2/§3; research |
| D4 | Auto-aprovar vínculo em âncora não-ambígua | **RATIFICADO + AMENDADO 2026-07-25 (D-121)** — CODPROD-único, EAN-único ou ambos concordantes auto-aprovam; conflito/colisão/hard-negative ficam REVIEW; `seller_sku`→CODPROD (não REFFORN). Ver ADR-05 amendado. | contrato §3; entrevista operador D-121 |
| D5 | Horizonte backfill pedidos | RATIFICADO 12 meses — **non-scope** (pedidos) | contrato §3 |
| D6 | Cadência coleta | **RESOLVIDO = sync_state cadence-agnostic.** MIS-006 guarda `schedule` genérico por entity; não hardcoda "diário". Valor real ratifica na missão mercado sem mudar shape E8. | research; SYSTEM-BLUEPRINT §4 |
| D7 | Snapshot leve p/ fonte-banco | **RESOLVIDO = SIM.** SankhyaAdapter escreve snapshot no mirror (mirror é o join target compartilhado). Não é read-through puro. | STORAGE-SCHEMA §products_mirror |
| F3.7 | Rota de descoberta EAN→catálogo | **PENDENTE — owner-authority.** UNPROVEN. Exige rodada live **T13-T16** (mlprobe, EANs #004-E, read-only, credencial ativa do DB NUNCA exposta). REQUEST ao hub. Se provada → M-07 constrói discovery+enqueue. Se disprovada → REMOVE F3.7, honest-unknown, produto sem anúncio recebe só caminho de vínculo. | SCENARIO-WALKTHROUGH Adendo A1; ML-API-QUERY-CATALOG |
| stale | Política de ausência (F-XLSX-1) | **RATIFICADO** keep-absent + `absent_in_last_snapshot`, nunca delete físico (ADR-04) | mission.md; scenario Fase 5; contrato §1b |
| Raw | Raw json.RawMessage nos DTOs ML | **DIFERIDO (default)** — ADR-06; hub pode puxar como feature isolado barato | F-ADAPTER-1; research §8 |

## Milestone Strategy

7 milestones. Detalhe (write-set, EARS, IO, validação) em cada `M-*/milestone.md` +
`M-*/validation-contract.md`. Specs-fonte abaixo (uma linha densa por milestone; o milestone.md
elabora):

- **M-01 — sync_state + scheduler skeleton (Fase 0).** CREATE tabela `sync_state` (E8, bloco A);
  scheduler skeleton reusando o padrão ticker (`root.go:577` template) — loop + read/write cursor
  + last_error + seam de registro, cadence-agnostic; NENHUM job ML pesado. Novo módulo sync +
  entry `modules.json`. Owns: sync_state table, módulo sync, root.go bloco de tickers
  (additive-lock). Dep: — (Fase 0). Write-set: `internal/modules/sync/*`, migração bloco A,
  `contracts/governance/modules.json`, `composition/root.go` (append ticker).
- **M-02 — products_mirror + ProductSourceAdapter port + active-source config (Fase 1).**
  CREATE `products_mirror`+`_stock_locations` (E2.1, bloco B) com `absent_in_last_snapshot`/
  `stale_since`; CREATE `active_source` config (E9, bloco B); REMOVE `erpSource()` env-branch;
  REFACTOR composition p/ construir ambos adapters + resolver por-tenant; FORMALIZE port
  (`readports.Reader` + `Sync`+`Kind`, Eport); CREATE source-KIND model; estender domínio E2 a
  10 campos. Owns: mirror/config tables, `internal_read/ports`, `erp_import/domain` (source-kind),
  `composition/root.go` (source-wiring), OpenAPI active-source endpoint (contract-lock). Dep:
  M-01 (só seam root.go). **Fundação — tudo depende daqui.**
- **M-03 — XlsxAdapter.** KEEP `parser.go`; REFACTOR `import_service.go` → hook pós-completion:
  upsert-merge no mirror (keep-absent) + trigger link-gen + enqueue coleta (sync_state);
  REFACTOR `activeSourceFromContext` → lookup config; implementa port (upload_snapshot).
  Owns: `erp_import/application`, `erp_import/adapters/{internalread,xlsx}`. Dep: M-01+M-02.
- **M-04 — SankhyaAdapter.** REUSE `oracle/reader.go` como core; ADD sync entrypoint → escreve
  mirror; implementa port (live_read_through); cache-key extend ATÔMICO com 3ª source.
  **BLOQUEADO por `[TESTAR-SKW]`** (db-consult REQUEST antes de codar). Owns:
  `internal_read/adapters/{oracle,cache}`. Dep: M-02. ∥ M-03 (pacotes disjuntos).
- **M-05 — auto-vínculo.** REFACTOR `generation_service.go` (trigger interno) + REUSE sinal EAN
  (`collisions==1`) + REFACTOR `resolution_service.go` (auto-approve reusando audit) + E10 trail;
  KEEP handlers órfãos, add invocação interna. Owns: `product_links/*`. Dep: M-02+M-03.
- **M-06 — telas + SDK + chain-viz.** CREATE rota `/importacoes` + promover `ImportacaoSection`
  compartilhado + REFACTOR p/ N-imported→N-linked→N-enqueued (join sync_state+links); REFACTOR
  `ActiveSourceCard` localStorage→DB-por-tenant (`useActiveErpSource` + endpoint M-02); REFACTOR
  `listErpImports` SDK (shape rica) + CREATE métodos active-source/chain; KEEP /integracoes 4
  cards; VinculosPage badge auto-aprovado. Owns: `apps/web` AppRouter/nav/ImportacaoSection/
  integracoes, `web-query` useActiveErpSource, `sdk-runtime`+OpenAPI (contract-lock, seção
  disjunta de M-02). Dep: M-01+M-02+M-05.
- **M-07 — chicken-egg break + F3.7 discovery (CONDICIONAL).** Gate: rodada live T13-T16 (REQUEST
  hub). Se PROVADO: CREATE `descobrir_produto_catalogo(ean)` (EAN→catalog_product_id, read-only)
  + persist identidade + enqueue coleta p/ produto sem anúncio (o "caminho de mercado" que quebra
  o chicken-egg). Se DISPROVADO: REMOVE F3.7, registra honest-unknown. NÃO executa coleta
  (missão mercado). Owns: `market/application` (discovery/enqueue half), `cmd/mlprobe` (evidência).
  Dep: M-02 + prova live. ∥ ondas 2-4.

## Parallel Execution Plan

DAG e ownership em `architecture-map.md`. Resumo:

**Dependency DAG (edges):** M-01→M-03 · M-02→{M-03,M-04,M-05,M-06,M-07} · M-03→M-05 · M-05→M-06 ·
M-01→M-06 · (live T13-T16)→M-07. M-01∥M-02 (additive-lock root.go). M-03∥M-04.

**Ownership matrix (eixo de colisão → dono):**
| Eixo | M-01 | M-02 | M-03 | M-04 | M-05 | M-06 | M-07 |
|------|------|------|------|------|------|------|------|
| Migração | bloco A | bloco B | — | — | (E10 audit, bloco B+) | — | bloco C (condicional, só gate PASS: `product_catalog_identity` NOVA, nunca ALTER mirror) |
| DB shape | sync_state | mirror, active_source | — | — | product_links(+audit) | — | — |
| Módulo Go | sync(novo) | internal_read/ports, erp_import/domain | erp_import/app+adapters | internal_read/oracle+cache | product_links | — | market/app |
| root.go | tickers (add-lock) | source-wiring (owner) | — | — | — | — | — |
| Contrato/SDK | — | active-source (lock) | — | — | — | chain-read (lock, disjunto) | — |
| FE surface | — | — | — | — | — | AppRouter/nav/pages | — |

**Ondas:** (1) M-01∥M-02 · (2) M-03∥M-04 · (3) M-05 · (4) M-06 · M-07 paralelo pós-M-02+live.
Caminho crítico: M-02→M-03→M-05→M-06.

## Real-integration bindings (up-front, core §5 · no-stub)

Seams que EXIGEM prova contra dependência REAL (nunca stub sem autorização):

1. **Sankhya Oracle `[TESTAR-SKW]`** — mapear `TGFPRO.CODPROD/DESCRPROD/MARCA/REFERENCIA`,
   `TGFCUS` (qual custo: gerencial/reposição/média), `TGFEST.ESTOQUE/CODLOCAL` (−reservado?),
   `TGFBAR.CODBARRA` (EAN), `TGFTAB/TGFEXC` (qual NUTAB p/ preço), `TGFGRU`. **Bloqueia M-04.**
   Via: chip manda `REQUEST db-consult` ao hub → hub relaia à sessão MNOS especialista
   (`local_ec787804`, profile §6). M-04 valida contra Oracle REAL (env dev stack), nunca stub.
2. **ML EAN-discovery (F3.7)** — rodada live **T13-T16** read-only via `apps/server_core/cmd/mlprobe`
   (untracked no worktree hub-erp-main que `/workspace` monta): T13 `/products/search?product_identifier=EAN`
   hit-rate; T14 `/products/{id}/items` demanda; T15 fallback `/sites/MLB/search` (checar 403 PolicyAgent);
   T16 simulação margem produto B. EANs de `erp_import_products` protocolo **#004-E** (2012 prod).
   Credencial = da conta ativa no DB, **NUNCA exposta/impressa**. Evidência →
   `docs/design/evidence/ml-api/`. **Gate de M-07.** REQUEST ao hub (decisão + env).
3. **Dev stack (docker)** para L2 smoke (mirror end-to-end, toggle de source, cadeia pós-import) —
   hub-owned, chip manda `REQUEST dev-stack`.

## Quality Attributes (Non-Functional Scope)

Bar não-funcional além do baseline, com mitigação + critério de validação nomeado (não só prosa):

- **Segurança / segredos:** credenciais ML e Oracle NUNCA lidas de `.env*`, impressas, ou logadas
  (AC-05). Mitigação: `mlprobe` lê do token-storage existente; db-consult Oracle via hub REQUEST,
  nunca direto à sessão MNOS. Validação: `M07-C3` (grep 0-hits token na evidência) + `M04-C14`
  (creds Oracle nunca dumpadas).
- **Isolamento multi-tenant:** `tenant_id` em toda tabela/query nova (AC-01). Validação: `M02` AC-01
  nas 3 tabelas + `M04-C13` grep de filtro tenant. Auth/multi-tenant COMPLETO = MIS-005; o modelo
  de dados já nasce tenant-scoped e o guard de tenant do app existente cobre o PUT active-source.
- **Isolamento de cache cross-source:** `active_source` do ctx na chave de cache; extensão para a
  3ª fonte é atômica (ADR-03). Validação: `MC-04` + must-fail test de M-04 (lição
  `chip-import-fix-closed`).
- **Honestidade de dado (ADR-17):** unknown → NULL/`—`, nunca 0 fabricado. Validação: `MC-03`
  (`SELECT ... WHERE custo IS NULL` retorna rows do prospect) + `M02-C3` (grep 0 `DEFAULT 0`/
  `NOT NULL` na migração) + `M06-C5` (`—` na UI).
- **Boundary enqueue-não-executa:** MIS-006 enfileira mercado, nunca escreve `market_aggregates`/
  `competitor_offers` (MC-11). Validação: `M03-C9` + `M07-C8` (git diff 0 writes).
- **Confiabilidade de sync:** `sync_state.last_error` genérico, nunca vaza credencial; `scheduler`
  skeleton não dispara job ML pesado. Validação: `M01` (mensagem genérica + zero-ML-dispatch).

## Risks

| Risco | Prob | Impacto | Mitigação | Trigger | Owner |
|-------|------|---------|-----------|---------|-------|
| F3.7 disprovada em T13-T16 | Média | Alto (objetivo "caminho de mercado" fica parcial) | M-07 condicional; se disprova, entrega só caminho de vínculo + honest-unknown; mercado-path vira missão mercado | T15 retorna 403 PolicyAgent | hub |
| `[TESTAR-SKW]` demora / especialista indisponível | Média | Alto (bloqueia M-04) | M-04 é folha do DAG; ondas 1-3 seguem sem ele; db-consult cedo | REQUEST sem resposta | hub |
| root.go merge-conflict M-01×M-02 | Baixa | Médio | additive-lock com seções disjuntas; fallback serial | conflito no merge | hub |
| cache-key não estendida com sankhya | Baixa | Alto (poluição cross-source, bug já visto) | M-04 atômico + must-fail test (lição chip-import-fix) | review | M-04 |
| escopo vaza p/ execução de coleta | Média | Médio | boundary explícito: MIS-006 ENFILEIRA, não executa | PR toca market_aggregates write | reviewer |

## Current State

- Status: draft. planning_phase: validation (P4/P5 feitos; P6 VC + P7 readiness em curso).
- Base: main @138aac3dff20438d8ddc509daf4171d82e5e45f6.
- Adapter ML: reusar, não reescrever (F-ADAPTER-1).
- Profile §1 diz default branch `master` mas origin default = `main` — DRIFT a ratificar (não bloqueia).
- Inventário de refactor: `research/refactor-inventory.md` (+ 3 sub-arquivos file:line).

## Handoff

- Next owner: hub (execução por milestone) após P7 Ready.
- Required artifact paths: este `mission.md`; `interface-contracts-mis006.md`; `architecture-map.md`;
  `research/` (inventário + 3 investigações + contratos/decisões); `validation-contract.md`;
  `M-01..M-07/{milestone.md,validation-contract.md}`; `planning-reviews/` (P7).
- Pré-execução: 2 REQUESTs ao hub — db-consult Sankhya (M-04) + rodada live T13-T16 (M-07 gate).
