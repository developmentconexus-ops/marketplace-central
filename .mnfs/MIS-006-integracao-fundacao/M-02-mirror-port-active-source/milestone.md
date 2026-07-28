# M-02 — products_mirror + ProductSourceAdapter port + active-source config

```yaml
id: M-02
type: milestone
mission: MIS-006
status: done   # merged main @49ab3bdd (D-120); dual-gate AGREEMENT + spot-check 7/7 + P7 dev-stack live-drive GREEN @b70fe1b8
depends_on: [M-01]   # SÓ o seam composition/root.go (additive-lock, bloco de tickers) — resto independente
base_sha: 138aac3d
merged_sha: 49ab3bdd
validation_level: QA-0
```

## Objective

Construir a FUNDAÇÃO de dados e contrato que todo o resto da missão consome: o espelho canônico
de produto (`products_mirror`, E2.1), o port formal que unifica xlsx e Sankhya (`ProductSourceAdapter`,
Eport.1), e a config de fonte-ativa por tenant em banco (E9.0) que mata `MC_ERP_SOURCE`. Nenhum
adapter concreto é implementado aqui (M-03/M-04) — M-02 entrega tabela + interface + config +
wiring de composição para que ambos adapters sejam construtíveis lado a lado.

Ver `mission.md` §Milestone Strategy (linha M-02), ADR-01/02/03, `architecture-map.md` (M-02 é
raiz do DAG — M-03,M-04,M-05,M-06,M-07 dependem daqui).

## Scope

- CREATE `products_mirror` + `products_mirror_stock_locations` (migração bloco B, E2.1).
- CREATE `active_source` config table por tenant (migração bloco B, E9.0).
- REMOVE `erpSource(getenv)` env-branch (`composition/root.go:772`, `MC_ERP_SOURCE`).
- REFACTOR `composition/root.go` para construir AMBOS adapters (xlsx + oracle) e resolver a
  fonte ativa por-tenant/por-request (additive-lock com M-01 no bloco de tickers).
- FORMALIZE port `ProductSourceAdapter` (`readports.Reader` preservado + `Sync()` + `Kind()`).
- CREATE modelo `SourceKind` (`upload_snapshot` | `live_read_through`) — substitui a ideia de
  estender `ImportSource` ingenuamente (não serve, ver refactor-inventory-backend.md §9).
- ESTENDER domínio E2 (`NormalizedRow`) para os 10 campos do contrato (`local`, `preco_venda`,
  `grupo` formalizado, além dos 6 já existentes).
- CREATE endpoint OpenAPI de active-source (GET/PUT) + `sdk-runtime` regenerado no MESMO commit.

## Non-Scope

- Implementação concreta de `Sync()` para xlsx (M-03) ou Sankhya (M-04) — M-02 só define a
  assinatura e a tabela-alvo.
- Escrita real de rows em `products_mirror` a partir de um import/sync real (prova plena de
  MC-01 é conjunta M-02+M-03+M-04; M-02 prova só shape+schema, não fluxo end-to-end).
- Extensão da cache-key para a 3ª fonte (sankhya) — isso é atômico com M-04 (ADR-03, lição
  `chip-import-fix-closed`). M-02 só garante que `active_source` está disponível no `ctx` para
  quem for estender a key depois.
- Auto-vínculo (M-05), telas (M-06), discovery F3.7 (M-07).
- Qualquer job de sync pesado, execução de coleta de mercado.

## Feature Briefs

### F-01 — `products_mirror` + `products_mirror_stock_locations` (tabela + migração)

**EARS:**
- WHEN uma linha é upserted em `products_mirror` para `(tenant_id, codigo_produto)` THEN os
  campos ausentes na fonte permanecem `NULL` (nunca `0`/default) — ADR-17.
- WHEN um adapter faz Sync() e um produto do snapshot anterior está ausente do novo THEN a row
  correspondente NÃO é deletada; fica marcada `absent_in_last_snapshot=true` + `stale_since=now()`
  (ADR-04 — a lógica de merge em si é M-03/M-04, mas a COLUNA e o default existem aqui).
- WHEN o produto volta a aparecer num snapshot subsequente THEN a flag `absent_in_last_snapshot`
  é limpável (coluna suporta `false` de volta, não é write-once).
- IF `codigo_produto` está ausente da linha de origem THEN a linha é rejeitada no parse (D3-
  refinada) — não vira row no mirror.

**Inputs/Outputs (MUST have):**
- Tabela `products_mirror`: `tenant_id, source (sankhya|xlsx|catalogo_cliente), codigo_produto`
  PK `(tenant_id, codigo_produto)`, `descricao, referencia, ean, marca, grupo_codigo,
  grupo_descricao, ncm NULL, custo NUMERIC NULL, preco_venda NUMERIC NULL, estoque_total NUMERIC
  NULL, protocol_id NULL, absent_in_last_snapshot BOOL DEFAULT false, stale_since timestamptz
  NULL, updated_at timestamptz` (E2.1, `interface-contracts-mis006.md` §E2).
- Tabela `products_mirror_stock_locations`: `tenant_id, codigo_produto, local_codigo,
  local_descricao, quantidade` (FK lógica ao mirror via `(tenant_id, codigo_produto)`).
- `NormalizedRow` (`erp_import/domain/import.go:10-22`) estendido para os 10 campos E2 (add
  `local`, `preco_venda`, formaliza `grupo_codigo`/`grupo_descricao`).

**Negative Scenarios:**
- Migração roda em banco com `erp_import_products` populado → não quebra (tabelas independentes,
  sem FK física obrigatória cruzando o histórico).
- `custo`/`preco_venda`/`estoque_total` NUNCA têm `DEFAULT 0` ou `NOT NULL` na DDL — grep na
  migração deve dar zero hits de `DEFAULT 0` nessas 3 colunas (AC-03).

**Write-set:** `apps/server_core/migrations/09xx_products_mirror.{up,down}.sql` (bloco B, após
bloco A de M-01/`sync_state`); `internal/modules/erp_import/domain/import.go` (`NormalizedRow`
+10 campos); modelo de linha do mirror (novo tipo, provavelmente `internal_read/domain` ou
`erp_import/domain` — decisão de pacote fica com quem implementa, não fixado aqui).

---

### F-02 — `ProductSourceAdapter` port (formalização Eport.1)

**EARS:**
- WHEN o port `ProductSourceAdapter` é definido THEN ele embute `readports.Reader` sem alterar
  nenhuma assinatura existente (`FindProductsForLinking`/`GetSellableStock`/`GetCostAsOf`/
  `GetTaxInputs`) — consumidores (pricing, vínculo) não recompilam com mudança de comportamento.
- WHEN o port declara o lado novo THEN adiciona `Sync(ctx) (SyncResult, error)` e `Kind()
  SourceKind` — SEM implementação concreta nesta milestone (M-03/M-04 implementam).
- WHEN `SourceKind` é modelado THEN é um tipo NOVO e dedicado (`upload_snapshot` |
  `live_read_through`), NÃO uma extensão ingênua do enum `ImportSource` existente (que só cobre
  fontes upload-shaped) — ver refactor-inventory-backend.md §9 "naive enum extension would
  misfit the domain model".

**Inputs/Outputs (MUST have):**
- `internal_read/ports` (ou pacote equivalente): interface `ProductSourceAdapter` compilável,
  embutindo `readports.Reader` existente.
- Tipo `SourceKind` com exatamente 2 valores no E-port.1: `upload_snapshot`, `live_read_through`.
- `SyncResult` (shape mínima: contagem processada/erros — usada por M-03/M-04, não precisa ser
  rica aqui, só existir para o port compilar).

**Negative Scenarios:**
- Alterar qualquer assinatura de `readports.Reader` existente = falha de contrato (consumidores
  quebram) — o port SÓ adiciona, nunca modifica o read-side.
- `SourceKind` com valor único ou reaproveitando `ImportSource` diretamente = falha (viola a nota
  do inventário — Sankhya não tem "protocolo"/"snapshot" e não pode ser modelado como se tivesse).

**Write-set:** `internal/modules/internal_read/ports/*` (novo arquivo ou extensão do existente),
`internal/modules/erp_import/domain/import.go` ou novo módulo de domínio compartilhado para
`SourceKind` (ver nota de posicionamento no ADR-02 — decisão de pacote exato não é bloqueante,
mas DEVE ser visível a ambos `erp_import` e `internal_read` sem import cycle).

---

### F-03 — `active_source` config (tabela + repo + resolução por-tenant) + REMOVE `MC_ERP_SOURCE`

**EARS:**
- WHEN a composição (`composition/root.go`) inicializa THEN constrói AMBOS adapters (xlsx +
  oracle/sankhya) simultaneamente — nenhum boot-time branch escolhe só um.
- WHEN uma requisição/contexto de tenant precisa resolver a fonte ativa THEN consulta
  `active_source` (tabela, por `tenant_id`), NÃO uma env var.
- IF não existe config `active_source` para o tenant THEN a resolução falha fechado com
  `ErrUnknownActiveSource` (já existe em `reader.go:53`, reusar — não reimplementar) — NUNCA
  cai em fallback silencioso para uma fonte default.
- WHEN `MC_ERP_SOURCE` é removido THEN nenhum código lê essa env var em lugar nenhum do repo.

**Inputs/Outputs (MUST have):**
- Tabela `active_source`: `tenant_id [PK], active_source (sankhya|xlsx|catalogo_cliente),
  source_kind (upload_snapshot|live_read_through), set_at, set_by` (E9.0).
- Repo/lookup: método que recebe `tenant_id` e retorna `(active_source, source_kind)` ou
  `ErrUnknownActiveSource`.
- `composition/root.go`: remoção de `erpSource(getenv)` (linha ~772) + wiring que registra AMBOS
  readers/adapters na composição, resolução movida para request-time via o repo acima.
- `active_source` do `ctx` disponível para downstream (pré-requisito da chave de cache — a
  EXTENSÃO da `canonicalKey` para incluir a 3ª fonte é ATÔMICA com M-04, não entra aqui; M-02
  só garante que o valor está carregável do ctx onde M-04 vai consumi-lo).

**Negative Scenarios:**
- `grep -r MC_ERP_SOURCE apps/server_core` retorna 0 hits pós-M-02.
- Tenant sem row em `active_source` → chamada de resolução retorna erro tipado, nunca um valor
  default silencioso (AC-03 adjacente — "unknown vira default" é proibido também para source).
- Dois tenants com fontes diferentes no MESMO processo → cada resolução por-tenant retorna o
  valor correto (não há estado global compartilhado sobrevivendo do boot-env antigo).

**Write-set:** `internal/composition/root.go` (source-wiring, owner — seções disjuntas do bloco
de tickers de M-01, additive-lock); nova tabela + repo (`internal/modules/erp_import` ou novo
pacote `internal/modules/tenant_config`, decisão de pacote não bloqueante); migração bloco B
(mesma leva de F-01, `apps/server_core/migrations/09xx_active_source.{up,down}.sql` ou
consolidada na mesma migração do mirror).

---

### F-04 — OpenAPI active-source endpoint + `sdk-runtime`

**EARS:**
- WHEN o endpoint de active-source é publicado (GET/PUT config por tenant) THEN o OpenAPI spec
  e o `sdk-runtime` regenerado landam no MESMO commit da implementação (profile §7,
  contract-first — não é "planning propõe depois implementação materializa noutro commit").
- WHEN M-06 (chain-read) publicar sua própria seção de contrato THEN as duas seções (active-source
  de M-02, chain-read de M-06) são DISJUNTAS no OpenAPI — sem overlap de path/schema
  (`architecture-map.md` §Superfícies compartilhadas).

**Inputs/Outputs (MUST have):**
- OpenAPI: path novo (ex. `GET/PUT /tenants/{tenant_id}/active-source` ou equivalente já
  ratificado no schema real do repo) com schema `{active_source, source_kind, set_at, set_by}`.
  **Landou como `GET/PUT /config/active-source`** (single-tenant, tenant fixo server-side) —
  o exemplo acima nunca existiu e não deve ser citado como path.
- `sdk-runtime`: cliente gerado/atualizado no mesmo diff.

**Negative Scenarios:**
- Commit que só toca `sdk-runtime` sem o OpenAPI correspondente (ou vice-versa) = falha de
  MC-12 (contract-first quebrado).
- Endpoint aceita `active_source` fora do enum ratificado (`sankhya|xlsx|catalogo_cliente`) →
  rejeitado (400), nunca aceito e gravado cegamente.

**Write-set:** spec OpenAPI do repo (path exato a confirmar contra `contracts/` do repo real),
`sdk-runtime` (pacote gerado consumido por `apps/web`).

## Ownership & Concurrency (six-axis)

| Eixo | M-02 |
|------|------|
| Migração | bloco B (`products_mirror`, `products_mirror_stock_locations`, `active_source`) — após bloco A (`sync_state`, M-01) |
| DB shape | `products_mirror`, `products_mirror_stock_locations`, `active_source` — dono |
| Módulo Go | `internal_read/ports` (formaliza port), `erp_import/domain` (source-kind + NormalizedRow 10 campos) — dono |
| `root.go` | source-wiring — **dono**; M-01 tem additive-lock só no bloco de tickers (~:575-577), seções disjuntas |
| Contrato/SDK | seção active-source — **contract-lock**; disjunta da seção chain-read de M-06 |
| FE surface | nenhuma (M-06 consome o endpoint, não M-02) |

Regra de merge `root.go`: M-02 edita o bloco de source-wiring (hoje ~:520,:772); M-01 edita só o
bloco de tickers (~:575-577). Se o hub julgar risco de conflito alto mesmo com seções disjuntas,
fallback = serial M-01→M-02 (ver mission.md §Accepted assumptions).

## Dependencies

- **M-01** (só o seam `root.go`, bloco de tickers já registrado antes de M-02 tocar source-wiring
  no mesmo arquivo — reduz risco de merge, não dependência funcional real).
- Nenhuma dependência de M-03/M-04/M-05/M-06/M-07 — M-02 é raiz do DAG (todos dependem daqui,
  `architecture-map.md`).

## Validation

Critérios de missão que M-02 é dono ou co-dono (`validation-contract.md` da missão):

- **MC-01** (parcial — junto com M-03,M-04): schema+port existem; feeding real por AMBOS
  adapters é prova conjunta pós-M-03+M-04.
- **MC-03**: NULL honesto — schema não tem `DEFAULT 0`/`NOT NULL` em `custo`/`preco_venda`/
  `estoque_total`.
- **MC-04**: fonte ativa por tenant em banco; `MC_ERP_SOURCE` removido; base para não-poluição
  de cache (extensão real da key é M-04).
- **MC-12**: migração aditiva (só `ADD`); OpenAPI+`sdk-runtime` no mesmo commit.

Detalhe binário → `M-02-mirror-port-active-source/validation-contract.md`.
