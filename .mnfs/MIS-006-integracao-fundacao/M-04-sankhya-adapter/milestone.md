# M-04 — SankhyaAdapter

```yaml
id: M-04
type: milestone
mission: MIS-006
status: draft
depends_on: [M-02]
blocked_by: [TESTAR-SKW db-consult]
base_sha: 138aac3d
validation_level: QA-0
```

## Objective

Formalizar o Oracle reader existente (`internal_read/adapters/oracle/reader.go`) como o
`SankhyaAdapter` do port `ProductSourceAdapter` (Eport.1, M-02): reusar o core de query
live-read-through, ADICIONAR um entrypoint de sync que escreve em `products_mirror` (D7 —
snapshot no mirror, não é read-through puro), implementar `Kind()=live_read_through`, e estender
a chave de cache (`canonicalKey`) de forma ATÔMICA para incluir a 3ª fonte. M-04 NÃO reescreve o
adapter Oracle maduro — reuso, não reescrita (mission.md ADR-06, F-ADAPTER-1 verdict).

Ver `mission.md` §Milestone Strategy (linha M-04), §Real-integration bindings item 1, ADR-02/03,
`architecture-map.md` (M-04 folha de M-02, ∥ M-03), `research/refactor-inventory-backend.md` §4
(SankhyaAdapter) e §9 (cache key / source-kind).

## BLOQUEIO — `[TESTAR-SKW]` (ler antes de qualquer código)

M-04 **NÃO PODE INICIAR CODIFICAÇÃO** até o mapeamento de colunas Sankhya Oracle ser confirmado
pelo especialista via `REQUEST db-consult` ao hub (que relaia à sessão MNOS,
`local_ec787804-f8e9-4981-9c12-7d3f45292294`, profile §6). Chip NUNCA fala direto com a sessão
MNOS. Ver `mission.md` §Real-integration bindings item 1 e §Risks ("`[TESTAR-SKW]` demora /
especialista indisponível").

Isto é um pré-requisito bloqueante, não um nice-to-have: `reader.go:22-116` já executa queries
contra TGFPRO/TGFEST/TGFCUS/TGFBAR/TGFTAB/TGFGRU hoje, mas o MAPEAMENTO de negócio (qual custo,
qual NUTAB, estoque líquido de reservado) nunca foi ratificado por um especialista — só
hipotetizado em `INTEGRATION-DATA-CONTRACT.md §E2` (marcado `[TESTAR-SKW]` linha a linha). Codar
o entrypoint de sync sobre um mapeamento não confirmado arrisca materializar `custo`/`preco_venda`
errados no mirror compartilhado que M-05 (vínculo) e M-06 (telas) consomem — silenciosamente.

## Pré-ativação — F1 (achado do refuter M-02, D-120 — VINCULANTE)

Registrado pelo hub a partir do EVENT CLOSED de M-02 (`_chip-m02`, refuter verbatim). **Antes de
M-04 tornar o oracle/sankhya reader efetivamente vivo em produção**, é OBRIGATÓRIO rotear pela
`active_source` do tenant (ou selecionar equivalentes mirror-backed) os seguintes consumidores que
hoje leem `oracleDB` DIRETO e **contornam** o `routing.Reader` de M-02:

- `listingCostReader` (`root.go:602`)
- `inventoryStockReader` (`root.go:410` / `:435`)
- `profitabilityCfg.Internal` — batch cost/tax facts (`root.go:560`)
- assisted-linkage gate (`root.go:521`)

**Falha silenciosa evitada:** uma vez que o Oracle esteja alcançável enquanto o `active_source` de
um tenant é `xlsx`/`catalogo_cliente`, esses consumidores servem dados Oracle LIVE para custo de
anúncio / risco de inventário / rentabilidade daquele tenant, ENQUANTO o catálogo corretamente
serve dados de upload — mistura cross-source, sem erro. INERTE na fundação M-02 (oracle não
configurado); o `routing.Reader` core de M-02 É fail-closed. É seam de batch-reader, não do routing
core. **Consequência para M-04:** a linha da matriz Ownership abaixo (`root.go` = "nenhuma edição")
só vale para a implementação do adapter em si; a ATIVAÇÃO live de M-04 exige tocar esses 4 pontos
de wiring (roteá-los por `active_source` ou apontá-los ao mirror) — a ser confirmado no P5 de M-04
como escopo de wiring, não descoberto mid-flight. Referência cruzada: `mission.md` collision
matrix / `interface-contracts-mis006.md` (E9 active_source consumers).

## Prerequisites

Tabelas/colunas Sankhya a confirmar via db-consult ANTES de codar (checklist — cada item precisa
de resposta do especialista, não de suposição do chip):

1. **`TGFPRO`** — `CODPROD`, `DESCRPROD`, `MARCA`, `REFERENCIA`: confirmar nomes exatos de coluna
   e se `MARCA`/`REFERENCIA` vêm de `TGFPRO` direto ou de tabela relacionada (ex. `TGFMAR`).
2. **`TGFCUS`** — qual custo usar: gerencial, reposição, ou média? `TGFCUS` tem múltiplas colunas
   de custo (`CUSGER`, `CUSREP`, `CUSMED` — nomes hipotéticos, confirmar reais); mapear qual
   corresponde ao `custo` de E2.1.
3. **`TGFEST`** — `ESTOQUE`/`CODLOCAL`: estoque bruto ou já líquido de reservado? Se bruto,
   existe coluna de reservado a subtrair (`RESERVADO`, `QTBLOQUEADA`?) para chegar em
   `estoque_total` honesto (E2.1)? Multi-local (`CODLOCAL`) mapeia para
   `products_mirror_stock_locations` (M-02 F-01) — confirmar granularidade.
4. **`TGFBAR`** — `CODBARRA`: confirmar que é o EAN usado para vínculo (E4.1 auto-approve
   depende de EAN exato); produto pode ter múltiplos códigos de barra (variação de embalagem)?
   Se sim, qual regra de escolha (principal/primeiro/mais recente)?
5. **`TGFTAB`/`TGFEXC`** — qual `NUTAB` (tabela de preço) usar para `preco_venda` (E2.1, campo
   NOVO vs modelo atual)? Existe tabela de preço por tenant/filial que precise ser fixada, ou é
   configurável por instalação?
6. **`TGFGRU`** — mapear `grupo_codigo`/`grupo_descricao` (E2.1) — confirmar se `TGFGRU` é a
   fonte correta e como resolve hierarquia (grupo pode ter subgrupo?).

Resultado esperado do db-consult: documento de mapeamento (coluna Sankhya → campo E2.1,
ratificado, não hipotético) que vira a base literal do SQL em `reader.go`/query builder do sync
entrypoint. Se o db-consult não retornar a tempo, M-04 permanece bloqueado — ondas 1-3 da missão
seguem sem ele (M-04 é folha do DAG, mission.md §Risks).

## Scope

- ADD sync entrypoint no `SankhyaAdapter` (sobre `internal_read/adapters/oracle/reader.go`) que
  executa a query live-read-through e escreve/upserta em `products_mirror` (upsert-merge
  keep-absent, ADR-04 — mesma regra do xlsx, M-03).
- IMPLEMENT `Kind() SourceKind` retornando `live_read_through` (port formalizado em M-02).
- IMPLEMENT `Sync(ctx) (SyncResult, error)` no `SankhyaAdapter` usando o mapeamento confirmado
  pelo db-consult.
- EXTEND `canonicalKey()` (`internal_read/adapters/cache/cache.go`) para incluir a 3ª fonte
  (sankhya) na chave de cache — MESMO commit que adiciona o `SourceKind`/valor de source sankhya
  ao domínio (nunca split, lição `chip-import-fix-closed`).
- CREATE must-fail test provando isolamento cross-source (xlsx vs sankhya) na cache.

## Non-Scope

- Mapeamento de colunas Sankhya sem confirmação do especialista — NENHUMA linha de SQL do
  entrypoint de sync é escrita antes do db-consult retornar (ou ser nomeado `could-not-run`).
- Reescrita de `reader.go:22-116` (`FindProductsForLinking`, `GetSellableStock`, `GetCostAsOf`,
  `GetTaxInputs`) — reuso, o core de query permanece intocado; só ADD o entrypoint de sync.
- `sankhya_linkage_reader.go`/`sankhya_linkage_config.go`/`sankhya_linkage_service.go`
  (feature separada de conciliação fiscal de pedidos) — bounded context diferente, não tocar
  (`research/refactor-inventory-backend.md` §4).
- Auto-vínculo (M-05), telas (M-06), discovery F3.7 (M-07), qualquer job de sync pesado/cadência
  real (fica cadence-agnostic via E8, D6 — non-scope desta milestone fixar cadência).
- Extensão de `ImportSource` para incluir `sankhya` de forma ingênua — `SourceKind` (M-02) já
  resolve isto; se `ImportSource` (upload-shaped) precisar de um 3º valor discriminando sankhya
  como fonte selecionável em `active_source`, isso é E9 (M-02), não redefinição do enum aqui.

## Feature Briefs

### F-01 — Sync entrypoint reusando `oracle/reader.go` como core → escreve mirror snapshot

**EARS:**
- WHEN `[TESTAR-SKW]` é confirmado via db-consult THEN o entrypoint de sync usa o mapeamento
  ratificado (não hipótese) para montar a query sobre TGFPRO/TGFEST/TGFCUS/TGFBAR/TGFTAB/TGFGRU.
- WHEN o sync roda THEN cada linha lida via Oracle é upserted em `products_mirror` com
  `source='sankhya'`, seguindo a MESMA regra de merge keep-absent do xlsx (ADR-04) — produto que
  sai do resultado da query não é deletado fisicamente, fica `absent_in_last_snapshot=true` +
  `stale_since=now()`.
- WHEN um campo Sankhya não resolve (ex. produto sem `TGFBAR` associado) THEN o campo grava
  `NULL` no mirror, nunca `0`/default (ADR-17, AC-03).
- IF o db-consult ainda não retornou THEN nenhuma query de sync é escrita — a feature fica
  `could-not-run` nomeada (nunca substituída por suposição do chip).

**Inputs/Outputs (MUST have):**
- `internal_read/adapters/oracle/reader.go`: novo método (ex. `Sync(ctx) (SyncResult, error)`)
  que reusa a infraestrutura de query existente (`buildFindProductsQuery`-style) para o
  mapeamento ratificado, e escreve no repo de `products_mirror` (M-02 F-01).
- Cada row escrita: `tenant_id, source='sankhya', codigo_produto, descricao, referencia, ean,
  marca, grupo_codigo, grupo_descricao, ncm NULL, custo, preco_venda, estoque_total,
  protocol_id=NULL (sankhya não tem protocolo), absent_in_last_snapshot, stale_since,
  updated_at` (E2.1).
- `products_mirror_stock_locations` populado se o mapeamento de `TGFEST.CODLOCAL` resolver
  multi-local (Prerequisite item 3).

**Negative Scenarios:**
- Sync roda contra Oracle dev stack (REAL, nunca stub — AC-04) e um produto sai da query entre
  duas rodadas → row permanece no mirror com `absent_in_last_snapshot=true`, `product_links`
  associado intacto (mesmo teste de M-03 F-XLSX-1, fonte diferente).
- Custo/estoque/preço não resolvíveis pelo mapeamento (produto sem `TGFCUS`/`TGFTAB` associado)
  → `NULL`, grep de `DEFAULT 0`/hardcode 0 nesses campos = zero hits.
- Query de sync roda sem `tenant_id` no filtro → rejeitada em code review (AC-01, todas as
  queries desta milestone são tenant-scoped).

**Write-set:** `internal/modules/internal_read/adapters/oracle/reader.go` (ADD método de sync,
reader core intocado); repo/writer de `products_mirror` (consumido, não recriado — dono é M-02);
possível novo arquivo `internal_read/adapters/oracle/sync.go` se o entrypoint for grande o
suficiente para separar do reader (decisão de arquivo não bloqueante).

---

### F-02 — Port implementation: `Kind()=live_read_through` + `Sync()`

**EARS:**
- WHEN `SankhyaAdapter` implementa `ProductSourceAdapter` (port formalizado em M-02 F-02) THEN
  `Kind()` retorna `live_read_through` (não `upload_snapshot` — Sankhya não tem protocolo, D1).
- WHEN consumidores chamam o read-side (`FindProductsForLinking`/`GetSellableStock`/
  `GetCostAsOf`/`GetTaxInputs`) THEN o comportamento é IDÊNTICO ao reader atual — nenhuma
  assinatura muda (M-02 F-02 já garante isto no port; M-04 só CUMPRE a garantia na implementação
  concreta).
- WHEN `Sync()` é chamado THEN retorna `SyncResult` (shape mínima de M-02 F-02: contagem
  processada/erros) refletindo o resultado do entrypoint F-01.

**Inputs/Outputs (MUST have):**
- `SankhyaAdapter` (tipo concreto, provavelmente o próprio `oracle.Reader` estendido ou um
  wrapper fino) satisfaz `ProductSourceAdapter` — `go build ./...` verde comprovando compilação
  contra a interface de M-02.
- `Kind()` hardcoded `live_read_through` (constante, não configurável — Sankhya é estruturalmente
  read-through, D1 ratificado).

**Negative Scenarios:**
- `Kind()` retornando `upload_snapshot` para sankhya = falha de contrato (mistura fonte
  read-through com fonte upload-shaped, quebra a distinção que M-02 F-02 introduziu
  especificamente para evitar isto).
- Qualquer assinatura de `FindProductsForLinking`/`GetSellableStock`/`GetCostAsOf`/`GetTaxInputs`
  alterada nesta milestone = falha (consumidores de pricing/vínculo não podem recompilar com
  comportamento diferente).

**Write-set:** `internal/modules/internal_read/adapters/oracle/reader.go` (implements do port);
nenhuma mudança em `internal_read/ports/*` (M-02 já formalizou, M-04 só implementa).

---

### F-03 — Cache-key ATÔMICA extension para 3ª source (sankhya)

**EARS:**
- WHEN a fonte `sankhya` é adicionada ao domínio de `active_source` (E9, M-02) THEN a mesma
  mudança (MESMO commit/PR) estende `canonicalKey()` (`internal_read/adapters/cache/cache.go`)
  para incluir a fonte na composição da chave — nunca em commits separados.
- WHEN dois tenants (ou o mesmo tenant em momentos diferentes) usam fontes diferentes (xlsx vs
  sankhya) THEN o cache NUNCA retorna dado da fonte errada — chave inclui `tenant_id` E `source`.
- IF a chave de cache for estendida sem o must-fail test correspondente THEN a feature é
  considerada incompleta (a lição de `chip-import-fix-closed` foi exatamente uma extensão de
  escopo sem key-fix simultâneo causando poluição cross-source silenciosa).

**Inputs/Outputs (MUST have):**
- `canonicalKey()` inclui `source` (agora com 3 valores possíveis: `xlsx`, `catalogo_cliente`,
  `sankhya`) na composição, ao lado de `tenant_id` já presente (fix anterior, D-119).
- Must-fail test: cache populado para `(tenant, xlsx)` NÃO É retornado para `(tenant, sankhya)` —
  teste que FALHA se a key não incluir a fonte (prova viva do isolamento, não uma asserção
  cosmética).

**Negative Scenarios:**
- Reverter a inclusão de `source` na key (simular a regressão) → o must-fail test PEGA a
  regressão (falha), provando que o teste é load-bearing, não decorativo.
- Cache-key estendida num commit separado do que adiciona `sankhya` ao domínio de fonte → viola
  a regra atômica (ADR-03, `architecture-map.md` §Superfícies compartilhadas linha
  `internal_read/adapters/cache` key: "extend ATÔMICO com adição de sankhya, nunca split").

**Write-set:** `internal/modules/internal_read/adapters/cache/cache.go` (`canonicalKey()`);
teste correspondente (`cache_test.go` ou equivalente) na MESMA mudança.

## Ownership & Concurrency (six-axis)

| Eixo | M-04 |
|------|------|
| Migração | nenhuma (M-04 não cria tabela — escreve em `products_mirror`, tabela de M-02) |
| DB shape | nenhuma tabela nova; consumidor de `products_mirror`/`products_mirror_stock_locations` (M-02) |
| Módulo Go | `internal_read/adapters/{oracle,cache}` — **dono** |
| `root.go` | nenhuma edição direta nesta milestone (wiring de ambos adapters já é escopo de M-02 F-03; M-04 só precisa que o `SankhyaAdapter` exista para M-02 construí-lo) |
| Contrato/SDK | nenhum (M-04 não expõe endpoint novo) |
| FE surface | nenhuma |

∥ **M-03** (`erp_import/*` vs `internal_read/adapters/oracle`) — pacotes disjuntos, sem colisão
de arquivo. Regra não-negociável: F-03 (cache-key) nunca split de outra milestone/chip —
`internal_read/adapters/cache` é dono exclusivo de M-04 (`architecture-map.md` §Superfícies).

## Dependencies

- **M-02** — port `ProductSourceAdapter` formalizado (F-02), `products_mirror` existe (F-01),
  `active_source`/`SourceKind` existem (F-03) — M-04 implementa contra essas formas já landadas.
- **`[TESTAR-SKW]` db-consult** — bloqueio hard, não uma dependência de milestone (é um
  REQUEST ao hub, relaiado ao especialista Oracle). Sem confirmação, M-04 não inicia
  codificação (Prerequisites acima).
- Nenhuma dependência de M-03/M-05/M-06/M-07.

## Validation

Critérios de missão que M-04 é dono ou co-dono (`validation-contract.md` da missão):

- **MC-01** (parcial — junto com M-02,M-03): schema+port já existem (M-02); M-04 prova o lado
  Sankhya do "alimentado por AMBOS adapters" — `SELECT ... FROM products_mirror WHERE
  source='sankhya'` retorna rows após sync contra Oracle REAL (dev stack).
- **MC-03** (parcial — junto com M-02): NULL honesto no lado Sankhya — custo/preço/estoque não
  resolvíveis gravam `NULL`, nunca `0`.
- **MC-09**: SankhyaAdapter validado contra Oracle REAL, não stub. Se `[TESTAR-SKW]` bloqueado →
  `could-not-run` nomeado, nunca substituído por suposição/stub (AC-04).

Detalhe binário → `M-04-sankhya-adapter/validation-contract.md`.
