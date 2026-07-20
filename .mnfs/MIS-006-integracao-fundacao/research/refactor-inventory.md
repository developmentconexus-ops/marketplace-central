# Inventário de Refactor — MIS-006-integracao-fundacao (CONSOLIDADO)

Base: `main` @138aac3d. Pedido central do operador (verbatim): *"já saber o que tem que
refatorar para atingir o objetivo, e o que for impossível vamos remover."*

Este é o índice consolidado. Detalhe file:line vive em 3 sub-arquivos verificados (arquivos
abertos, não grep-adivinhados):
- [`refactor-inventory-backend.md`](refactor-inventory-backend.md) — Go (`apps/server_core`), 9 subsistemas.
- [`refactor-inventory-frontend.md`](refactor-inventory-frontend.md) — React (`apps/web`) + SDK.
- [`contracts-decisions-scenario.md`](contracts-decisions-scenario.md) — E1–E7, decisões, IO pairs, flaws.

Legenda de ação: **CREATE** net-new · **REFACTOR** existe, muda de forma/wiring · **KEEP**
maduro, não tocar · **REMOVE** impossível de aproximar do alvo, ou stub morto.

---

## 1. Veredicto por subsistema (resumo executivo)

| # | Subsistema | Veredicto dominante | 1 linha |
|---|------------|---------------------|---------|
| S1 | `products_mirror` | **CREATE** (tabela+migração) | 0 hits repo; "verdade de produto" hoje é recomputada por-request de 2 readers divergentes |
| S2 | `ProductSourceAdapter` port | **REFACTOR** | port de-facto JÁ existe (`readports.Reader` satisfeito por xlsx + oracle); falta nome + metade write/sync |
| S3 | XlsxAdapter | **KEEP parser + REFACTOR service** | `parser.go` lenient provado live = KEEP; `import_service.go` para no snapshot, nunca dispara cadeia = REFACTOR |
| S4 | SankhyaAdapter | **REUSE oracle reader + CREATE sync path** | `oracle/reader.go` maduro = base; falta entrypoint que grava mirror; `[TESTAR-SKW]` bloqueia código |
| S5 | Chicken-egg | **REFACTOR (trigger), NÃO reescrita** | metade catalog-evidence já é link-free; defeito real = nada chama `Collect` em lote |
| S6 | Auto-vínculo | **REFACTOR (wiring) + REUSE sinal EAN** | `ApproveCandidate` 100% manual; sinal de unicidade EAN JÁ existe nos 2 readers |
| S7 | sync_state + scheduler | **CREATE tabela + REFACTOR ticker-template** | 0 tabela; ticker `fee_sync` (root.go:577) = molde reusável do loop |
| S8 | fee_sync / Raw DTO / backoff | **DEFER (fora de escopo) + Raw opcional** | 3 buracos F-ADAPTER-1; pertencem à missão que os exercita (backfill) |
| S9 | Active-source por tenant | **REMOVE env-branch + REFACTOR ctx-toggle** | `erpSource()` boot env = REMOVE; ctx-toggle já é a forma certa = estender |

---

## 2. Lista REMOVE (o "impossível / morto" que o operador pediu para nomear)

| path:line | o que é | por que REMOVE |
|---|---|---|
| `internal/composition/root.go:772` `erpSource(getenv)` (branch `MC_ERP_SOURCE` de boot) | escolhe UMA fonte para o backend inteiro no boot | contrato §1b "`MC_ERP_SOURCE` morre"; incompatível com fonte-ativa-por-tenant (dois tenants, duas fontes, mesmo processo). D120-POSTMORTEM I1 = "root cause confusão". Substituir por config em banco. |
| padrão "read = recompute" em `erp_import/adapters/internalread/reader.go:84-107` (`snapshot` rescan de `AcceptedRows` por query) | reconstrói estado corrente varrendo o snapshot inteiro a cada leitura | viola SYSTEM-BLUEPRINT linha 8-9 ("enrichment no ingest, não no read"); impossível de manter com JOIN SQL. O rescan sai; vira leitura fina do mirror. |
| — (nada mais qualifica como stub morto REMOVÍVEL) | — | orphan endpoints (`product_links/transport/http_handler.go:89-90`) NÃO são REMOVE: são wiring faltante (REFACTOR — adicionar chamador interno). Handlers ficam. |

**Nota honesta:** o inventário achou POUCO código genuinamente morto/removível. O gargalo não é
lixo acumulado — é **cadeia nunca conectada** (triggers órfãos, fonte mal selecionada, estado
nunca materializado). O trabalho é 80% REFACTOR/wiring + CREATE de fundação, ~5% REMOVE.

---

## 3. Lista CREATE (fundação net-new)

| item | onde | migração/módulo |
|---|---|---|
| `products_mirror` + `products_mirror_stock_locations` | nova tabela + writer | migração bloco B |
| `sync_state` | nova tabela + engine skeleton | migração bloco A |
| active-source config (por tenant) | nova tabela + repo | migração bloco B (junto do mirror) |
| source-KIND model (upload-snapshot vs live-read-through) | `erp_import/domain` + novo tipo | — (a extensão naïve do enum `ImportSource` NÃO serve, ver S9) |
| SankhyaAdapter sync entrypoint | sobre `oracle/reader.go` | módulo `internal_read` |
| E8/E9/E10 contratos | doc | `interface-contracts-mis006.md` |
| F3.7 `descobrir_produto_catalogo(ean)` | `market/application` (CONDICIONAL live) | — |

---

## 4. Sequenciamento por superfície de colisão (entrada do Parallel Plan)

Arquivos/tabelas que >1 stream toca — ordenar ou travar (fonte: §"Collision surfaces" do backend):

- **`internal/composition/root.go`** — hotspot único: S7 (registro scheduler ~:575-577) + S9
  (source-wiring ~:520,:772). Editar seções DISJUNTAS. → **additive-lock** (M-02 dono, M-01
  append no bloco de tickers) OU serial M-01→M-02. Fallback serial se hub julgar risco alto.
- **`internal/modules/erp_import/`** — maior densidade (4 subsistemas: S1 mirror-write, S3
  post-import trigger, S6 link-gen trigger, S9 source-domain). Dono único = trilha xlsx (M-03),
  com M-02 congelando o domínio de source antes.
- **`internal/modules/internal_read/` (ports, oracle, cache)** — S2 (port), S4 (Sankhya), S9
  (cache key). **Cache-key extend DEVE landar no MESMO change da adição da 3ª source** (lição
  `chip-import-fix-closed`: split causou poluição de cache cross-source uma vez). M-04 atômico.
- **Migrações novas** — todas 9xx-series, sem conflito com existentes; ordenar sync_state (Fase 0)
  antes de mirror (Fase 1) por intenção de design — tabelas são independentes, só o NÚMERO ordena.
- **`product_links/`** (S6) depende de S1 (mirror) + S3 (hook de import) landarem antes.
- **`market/`** (S5) depende de S1 (mirror = join) + S6 (EANs auto-vinculados).
- **`connectors/adapters/mercado_livre/`** (S8) — ISOLADO, paralelizável, sem tabela/arquivo
  compartilhado → mas majoritariamente FORA de escopo (defer ML-sync).

---

## 5. Mapa subsistema → milestone (ver mission.md §Milestone Strategy)

| Subsistema | Milestone |
|---|---|
| S7 sync_state + scheduler skeleton | **M-01** |
| S1 mirror + S2 port + S9 active-source/config | **M-02** |
| S3 XlsxAdapter (parser KEEP + service REFACTOR + upsert-merge) | **M-03** |
| S4 SankhyaAdapter (reuse oracle + sync path + `[TESTAR-SKW]`) | **M-04** |
| S6 auto-vínculo (trigger + reuse EAN signal + auto-approve + audit) | **M-05** |
| FE /importacoes + /integracoes + SDK + S5 chain-viz | **M-06** |
| S5 chicken-egg break + F3.7 discovery (CONDICIONAL live) | **M-07** |
| S8 fee_sync/Raw/backoff | **DEFERRED** (ML-sync; Raw = candidato barato, hub decide) |
