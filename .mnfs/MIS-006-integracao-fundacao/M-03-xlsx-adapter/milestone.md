# M-03 — XlsxAdapter

```yaml
id: M-03
type: milestone
mission: MIS-006
status: draft
depends_on: [M-01, M-02]
base_sha: 138aac3d
validation_level: QA-0
```

## Objective

Fazer o xlsx virar o primeiro `ProductSourceAdapter` REAL, ligando o parser leniente já
provado (`parser.go`, KEEP) ao `products_mirror` que M-02 criou. Hoje `import_service.go` para
em "snapshot persistido" — M-03 estende a orquestração para: upsert-merge no mirror (keep-absent,
ADR-04), disparar geração de candidato de vínculo, e enfileirar coleta em `sync_state` (E8, só
enfileira — não executa, ADR-06/boundary MC-11). Também move a leitura de `read = recompute`
(`internalread/reader.go:84-107`, rescan de `snapshot.AcceptedRows` a cada request) para leitura
fina do mirror (ADR-01), e troca `activeSourceFromContext` de default-de-request para lookup na
config `active_source` (E9, M-02).

Ver `mission.md` §Milestone Strategy (linha M-03), ADR-01/02/04, `architecture-map.md` (M-03 é
folha de M-01+M-02, raiz de M-05), `research/refactor-inventory-backend.md` §1/§3.

## Scope

- REFACTOR `import_service.go`: hook pós-completion de protocolo → (a) upsert-merge das rows
  aceitas em `products_mirror` (keep-absent), (b) trigger interno de geração de candidato de
  vínculo (chamada de função, não endpoint HTTP — endpoint órfão `product_links/transport/
  http_handler.go:89-90` fica KEEP, ver M-05), (c) enqueue de coleta em `sync_state`
  (`entity=market` ou equivalente, só grava cursor/estado — NUNCA chama `Collect()`).
- REFACTOR `internalread/reader.go:84-107` (`Reader.snapshot` e os métodos que rescaneiam
  `snapshot.AcceptedRows` — `FindProductsForLinking`/`GetSellableStock`/`GetCostAsOf`/
  `catalogPage`): leitura passa a ser feita do `products_mirror` (thin read), não mais recompute
  do snapshot completo a cada chamada.
- REFACTOR `internalread/reader.go:24-70` (`WithActiveSource`/`activeSourceFromContext`/
  `ParseActiveSource`): fonte de verdade passa da ctx-default para lookup em `active_source`
  (tabela de M-02); MANTÉM o shape ctx-scoped (não regressa a env/global).
- IMPLEMENTA o port `ProductSourceAdapter` (formalizado em M-02) para xlsx: `Sync()` = ponto de
  entrada do upsert-merge pós-completion; `Kind()` retorna `upload_snapshot`.

## Non-Scope

- `parser.go` — ZERO edição (KEEP absoluto; já provado live #004-E 2012 produtos,
  `lenient-xlsx-import-two-source`).
- SankhyaAdapter (M-04) — pacote disjunto, ∥.
- Lógica de auto-aprovação EAN-exato-único em si (`generation_service.go`/
  `resolution_service.go`) — M-03 só CHAMA a geração de candidato após o upsert-merge; a regra
  de auto-approve é M-05.
- Execução real da coleta de mercado (`Collect(codprod)`, `market_aggregates`,
  `competitor_offers` writes) — M-03 só ENFILEIRA em `sync_state` (boundary MC-11, ADR-06's
  "adapter maduro" não é tocado aqui).
- Telas /importacoes (M-06) — M-03 só garante que os dados existem para M-06 consumir.
- Extensão de `ImportSource` para incluir `sankhya` como 3º valor — isso é atômico com M-04
  (research §9, "naive enum extension"); M-03 não mexe no enum além do que já usa hoje
  (`xlsx`/`catalogo_cliente`).

## Feature Briefs

### F-01 — Hook pós-completion em `import_service.go`: upsert-merge keep-absent no mirror

**EARS:**
- WHEN um protocolo de import xlsx completa (rows aceitas persistidas em `erp_import_products`,
  padrão atual preservado) THEN o hook pós-completion faz upsert-merge das rows aceitas em
  `products_mirror` — cada row aceita vira `INSERT ... ON CONFLICT (tenant_id, codigo_produto)
  DO UPDATE` (source=`xlsx`, campos do E2, `absent_in_last_snapshot=false`).
- WHEN um `codigo_produto` que existia no snapshot anterior do MESMO tenant está ausente do
  snapshot novo THEN a row correspondente em `products_mirror` NÃO é deletada; fica
  `absent_in_last_snapshot=true` + `stale_since=now()` (ADR-04, F-XLSX-1) — `product_links`
  apontando para esse `codigo_produto` continuam intactos (nada em cascata os apaga).
- WHEN um `codigo_produto` que estava `absent_in_last_snapshot=true` volta a aparecer num
  snapshot subsequente THEN a flag limpa (`absent_in_last_snapshot=false`, `stale_since=NULL`).
- IF a linha de origem não tem `codigo_produto` THEN ela já é rejeitada no parse (KEEP,
  comportamento de `parser.go` preservado) — nunca chega ao hook de upsert.
- WHEN um campo diferente de `codigo_produto` está ausente na linha de origem (custo, estoque,
  etc.) THEN o upsert grava `NULL` nesse campo do mirror — nunca `0`/default (ADR-17, prova viva
  = prospect #004-E).
- WHEN o upsert-merge termina com sucesso THEN o hook dispara (a) geração de candidato de
  vínculo (chamada interna à função já existente de `generation_service.go`, não uma nova
  implementação) e (b) grava/atualiza uma row em `sync_state` (`entity=market`, chave E8
  `(tenant_id, installation_id, entity)`, cursor aponta para os `codigo_produto`s recém-tocados
  via append atômico M01-C11) — SEM chamar `Collect()` (enfileira, não executa).

**Inputs/Outputs (MUST have):**
- `import_service.go`: novo passo pós-`erp_import_products` write, antes de retornar sucesso ao
  chamador do import.
- Upsert usa os 10 campos E2 de `NormalizedRow` (já estendido por M-02).
- Chamada interna (função Go, não HTTP) à geração de candidato — reusa
  `generation_service.go:60-78` como está hoje (M-05 é quem refina a regra de auto-approve
  dentro dela; M-03 só garante que ela É chamada).
- `sync_state` upsert: `entity=market`, `cursor` (jsonb com os `codigo_produto`s do batch via
  append atômico), `last_incremental_at=now()` — shape de M-01/E8, M-03 só popula.

**Negative Scenarios:**
- Import de snapshot SEM um produto que existia antes → `SELECT absent_in_last_snapshot, stale_since
  FROM products_mirror WHERE codigo_produto=X` retorna `true`/timestamp setado, row PRESENTE
  (não 0 rows).
- `product_links` de um produto que ficou `absent_in_last_snapshot=true` → `SELECT * FROM
  product_links WHERE internal_product_id=...` continua retornando a mesma row (nenhuma
  cascata de delete/orphan).
- Import do prospect #004-E (custo/estoque ausentes na planilha) → `SELECT custo, estoque_total
  FROM products_mirror WHERE source='xlsx' AND codigo_produto IN (...)` retorna `NULL`, nunca
  `0`.
- Diff do hook NÃO contém nenhuma chamada a `Collect(`/write em `market_aggregates`/
  `competitor_offers` (grep no diff = 0 hits) — boundary MC-11.
- Falha no upsert-merge (ex. erro de conexão) NÃO deve deixar o protocolo em estado "completo"
  sem o mirror atualizado silenciosamente — erro é propagado/logado, não engolido (ladder de
  observabilidade mínima; não exige transação distribuída cross-DB se mirror e
  `erp_import_products` estiverem no MESMO Postgres — usar transação local se viável).

**Write-set:** `internal/modules/erp_import/application/import_service.go` (hook pós-completion);
`internal/modules/erp_import/adapters/xlsx/` (SÓ o ponto de entrada `Sync()`/`Kind()` do port —
NUNCA `parser.go`); possivelmente um novo repo/query de upsert em
`internal/modules/erp_import/adapters/postgres/` (mirror write) ou pacote equivalente decidido
por quem implementa.

---

### F-02 — Migra leitura de rescan-de-snapshot para mirror (`internalread/reader.go:84-107`)

**EARS:**
- WHEN um consumidor chama `FindProductsForLinking`/`GetSellableStock`/`GetCostAsOf`/
  `catalogPage` no reader xlsx THEN a implementação lê de `products_mirror` (filtrado por
  `tenant_id` + `source='xlsx'`), NÃO mais recomputa `snapshot.AcceptedRows` do zero a cada
  chamada.
- WHEN a assinatura pública desses métodos é preservada (mesmo shape que consumidores
  (`pricing`, `vínculo`) já usam) THEN nenhum consumidor recompila com mudança de
  comportamento — só a implementação interna muda de fonte.
- WHEN `catalogPage` (leitura paginada) passa a ler de `products_mirror` THEN a query declara
  `ORDER BY` estável (`codigo_produto ASC`) — a ordem física arbitrária do Postgres NUNCA é
  aceita para paginação (sem sort declarado, páginas duplicam/pulam rows). Substitui a ordem
  implícita de `AcceptedRows` do rescan antigo por ordem explícita e determinística.

**Inputs/Outputs (MUST have):**
- `internalread/reader.go:84-107` reescrito para consultar `products_mirror` em vez de
  `LatestCompletedSnapshot` + rescan. NOTA: `FindProductsForLinking` é bulk não-paginado (input
  do matcher por-EAN, ordem irrelevante) — só `catalogPage` (paginado) exige `ORDER BY`.
- `readports.Reader` (E-port, formalizado em M-02) continua satisfeito pela nova implementação.

**Negative Scenarios:**
- Chamada aos 4 métodos citados após o refactor NÃO dispara nenhuma query a
  `erp_import_protocols`/`erp_import_products` para servir leitura corrente (essas tabelas
  continuam existindo só como history, `LatestCompletedSnapshot` fica reservado para
  auditoria/protocolo, não para servir leitura de produto corrente).
- Consumidor existente (pricing, vínculo) compila sem alteração de import/assinatura
  (`go build ./...` verde, zero diff nos pacotes consumidores).

**Write-set:** `internal/modules/erp_import/adapters/internalread/reader.go`.

---

### F-03 — `activeSourceFromContext` → lookup em `active_source` (config de M-02)

**EARS:**
- WHEN a resolução de fonte ativa para uma requisição de leitura xlsx acontece THEN ela consulta
  a tabela `active_source` (M-02, por `tenant_id`), não mais um valor default de ctx.
- WHEN o mecanismo ctx-scoped (`WithActiveSource`) é preservado como SHAPE (não regressa a env
  global) THEN o VALOR que ele carrega passa a ser alimentado pelo lookup em `active_source` em
  vez do default hardcoded do request.
- IF não existe row em `active_source` para o tenant THEN a resolução falha fechado
  (`ErrUnknownActiveSource`, já existe `reader.go:53`, reusar) — nunca cai num default
  silencioso de `xlsx`.

**Inputs/Outputs (MUST have):**
- `internalread/reader.go:24-70` (`activeSourceKey`, `WithActiveSource`, `ParseActiveSource`)
  refatorado para ler de `active_source` (repo/lookup de M-02) em vez de
  `ParseActiveSource(defaultParam)`.
- Fail-closed preservado.

**Negative Scenarios:**
- `grep -n "default.*xlsx\|xlsx.*default"` na função de resolução pós-refactor = 0 hits de
  fallback silencioso a xlsx (era o comportamento antigo pré-D119 fix, não pode voltar).
- Tenant sem `active_source` configurada → leitura retorna erro tipado, não uma lista vazia
  disfarçada de "sem produtos".

**Write-set:** `internal/modules/erp_import/adapters/internalread/reader.go` (mesmas linhas de
F-02, mudança relacionada — pode ser o mesmo commit).

---

### F-04 — XlsxAdapter implementa o port `ProductSourceAdapter` (`upload_snapshot`)

**EARS:**
- WHEN o `ProductSourceAdapter` (interface formalizada em M-02) é implementado para xlsx THEN
  `Kind()` retorna `upload_snapshot` (SourceKind de M-02, sem valor novo).
- WHEN `Sync(ctx)` é chamado no adapter xlsx THEN ele é o mesmo caminho de código do hook F-01
  (upsert-merge pós-completion) — `Sync()` não é uma segunda implementação paralela, é o ponto
  de entrada nomeado do port sobre a MESMA lógica.
- WHEN o adapter compila contra o port de M-02 THEN `readports.Reader` (read-side) continua
  satisfeito pela implementação de F-02, sem alteração de assinatura.

**Inputs/Outputs (MUST have):**
- Tipo concreto em `erp_import/adapters/xlsx/` (ou pacote decidido por quem implementa)
  satisfazendo `ProductSourceAdapter` (read-side de F-02 + `Sync()`/`Kind()` de F-01).
- `SyncResult` (shape mínima definida por M-02) retornado por `Sync()` com contagem
  processada/erros do batch de upsert.

**Negative Scenarios:**
- `Kind()` retornando qualquer valor fora de `{upload_snapshot, live_read_through}` = falha de
  contrato (M02-C8 upstream).
- `Sync()` chamado duas vezes com o MESMO protocolo/snapshot não duplica rows no mirror (upsert
  é idempotente por `(tenant_id, codigo_produto)` — mesma garantia de F-01, testada aqui como
  prova de port).

**Write-set:** `internal/modules/erp_import/adapters/xlsx/` (novo arquivo de adapter/port
binding — NÃO `parser.go`).

## Ownership & Concurrency (six-axis)

| Eixo | M-03 |
|------|------|
| Migração | nenhuma (M-03 só escreve dados em tabelas já criadas por M-01/M-02) |
| DB shape | nenhuma tabela nova; escritor de `products_mirror` (M-02 é dona do shape) e leitor/escritor de `sync_state` (M-01 é dono do shape) |
| Módulo Go | `erp_import/application`, `erp_import/adapters/{internalread,xlsx}` — dono; `parser.go` dentro de `adapters/xlsx` é EXCEÇÃO explícita (KEEP, zero edição) |
| `root.go` | nenhum toque (composition-wiring é M-02) |
| Contrato/SDK | nenhum (M-03 não expõe endpoint novo; M-06 é quem expõe leitura de cadeia) |
| FE surface | nenhuma |

Regra de dependência de pacote: M-03 consome o port formalizado por M-02
(`internal_read/ports`) e a config `active_source`/tabela `products_mirror` de M-02 — M-02
precisa estar mergeado (schema+port existindo) antes de M-03 poder compilar contra eles; M-03
também consome a tabela `sync_state` de M-01 (só o shape, não o scheduler). ∥ M-04
(`internal_read/adapters/oracle` é pacote disjunto de `erp_import/*` — zero overlap de arquivo).

## Dependencies

- **M-01** — shape de `sync_state` (E8) precisa existir para o enqueue do hook F-01 escrever
  nela; M-03 não depende do scheduler/ticker em si, só da tabela.
- **M-02** — `products_mirror` (shape+PK+colunas), `active_source` (config), e o port
  `ProductSourceAdapter` formalizado precisam existir ANTES de M-03 implementar contra eles.
  M-02 é a "fundação — tudo depende daqui" (mission.md).
- Nenhuma dependência de M-04/M-05/M-06/M-07 — M-03 é consumido por M-05 (auto-vínculo depende
  do hook de trigger existir), mas não depende de M-05.

## Validation

Critérios de missão que M-03 é dono ou co-dono (`validation-contract.md` da missão):

- **MC-01** (parcial — junto com M-02,M-04): mirror alimentado pelo lado xlsx via `Sync()`.
- **MC-02**: upsert-merge keep-absent — M-03 é a ÚNICA dona (lógica de merge é implementada
  aqui, M-02 só deu a coluna).
- **MC-03** (parcial — junto com M-02): NULL honesto provado com dado real (#004-E) fluindo
  pelo hook de M-03.
- **MC-11** (parcial — todos): boundary enfileira-não-executa, provado no diff do hook F-01.

Detalhe binário → `M-03-xlsx-adapter/validation-contract.md`.
