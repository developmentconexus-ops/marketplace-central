# Milestone Validation Contract — M-03-xlsx-adapter

```yaml
id: M-03-VC
type: milestone-validation-contract
mission: MIS-006
milestone: M-03
validation_level: QA-0
base_sha: 138aac3d
```

Verdicts binários. Evidência = caminho inspecionável concreto (core §5, ladder L0-L4). Tipos:
`ran` (executado, output salvo), `assumed` (design, não executado), `could-not-run` (bloqueado —
nomear). Nenhum seam contra dependência real provado por stub sem autorização.

## Critérios

| ID | Critério | Prova mínima inspecionável | Ladder | Feature dona |
|----|----------|----------------------------|--------|--------------|
| M03-C1 | `parser.go` não é editado por esta milestone | `git diff` do commit/chip mostra ZERO hits em `erp_import/adapters/xlsx/parser.go` | L0 ran (diff) | F-01 |
| M03-C2 | Import xlsx completo dispara upsert-merge no `products_mirror` | teste de integração: import de snapshot com N rows aceitas → `SELECT count(*) FROM products_mirror WHERE tenant_id=T AND source='xlsx'` retorna N rows correspondentes pós-import | L2 ran | F-01 |
| M03-C3 | Keep-absent: produto ausente do snapshot novo mantém row, nunca delete físico | IO Fase 5 (mission.md/scenario): import 1 com produto X presente, import 2 sem X → `SELECT absent_in_last_snapshot, stale_since FROM products_mirror WHERE codigo_produto=X` retorna `true`/timestamp setado, row PRESENTE (não 0 rows) | L2 ran | F-01 |
| M03-C4 | `product_links` de produto que ficou stale continua intacto (F-XLSX-1 protection) | mesma rodada de M03-C3: `SELECT * FROM product_links WHERE internal_product_id=(mirror de X)` retorna a MESMA row pré e pós import 2 (nenhuma cascata de delete/orphan) | L2 ran | F-01 |
| M03-C5 | Flag `absent_in_last_snapshot` é limpável quando produto reaparece | import 3 traz X de volta → `SELECT absent_in_last_snapshot, stale_since FROM products_mirror WHERE codigo_produto=X` retorna `false`/`NULL` | L2 ran | F-01 |
| M03-C6 | NULL honesto: prospect #004-E (custo/estoque ausentes na planilha) grava `NULL`, nunca `0` | import real ou fixture equivalente ao #004-E → `SELECT custo, estoque_total FROM products_mirror WHERE source='xlsx' AND codigo_produto IN (...)` retorna `NULL` nessas rows | L2 ran | F-01 |
| M03-C7 | Linha sem `codigo_produto` é rejeitada no parse, nunca chega ao mirror como row órfã | fixture com linha sem `codigo_produto` → row não aparece em `products_mirror` (contagem do mirror = contagem de rows aceitas pelo parser, não pelo total de linhas da planilha) | L1 ran | F-01 |
| M03-C8 | Hook pós-completion dispara geração de candidato de vínculo (chamada interna, não HTTP) | diff mostra chamada de função a `generation_service` (ou equivalente) dentro de `import_service.go`; teste de integração: pós-import, candidatos existem para os produtos com EAN sem precisar de chamada HTTP externa | L2 ran | F-01 |
| M03-C9 | Hook pós-completion enfileira coleta em `sync_state`, NUNCA executa | pós-import, `SELECT * FROM sync_state WHERE entity=...` mostra cursor/estado atualizado; grep no diff do hook = 0 hits de `Collect(`/write em `market_aggregates`/`competitor_offers` | L1 ran + L0 ran (grep) | F-01 |
| M03-C10 | Falha no upsert-merge não deixa o protocolo "completo" com mirror desatualizado silenciosamente | teste: força erro no passo de upsert (ex. conexão simulada) → erro é propagado/logado, protocolo não fica marcado como sucesso pleno sem o mirror refletir | L1 ran | F-01 |
| M03-C11 | Leitura migrou de rescan-de-snapshot para mirror (`internalread/reader.go:84-107`) | diff mostra `FindProductsForLinking`/`GetSellableStock`/`GetCostAsOf`/`catalogPage` consultando `products_mirror`, não mais `LatestCompletedSnapshot`+rescan de `AcceptedRows` | L1 ran | F-02 |
| M03-C11b | `catalogPage` mirror-backed declara sort estável (paginação determinística) | diff da query `catalogPage` inclui `ORDER BY codigo_produto` (ou chave estável equivalente); teste de paginação: 2 páginas consecutivas não duplicam/pulam row | L1 ran | F-02 |
| M03-C12 | Consumidores (pricing, vínculo) compilam sem alteração de assinatura pós-refactor de leitura | `go build ./...` verde; `git diff` dos pacotes consumidores (fora de `erp_import`) = 0 hits | L1 ran | F-02 |
| M03-C13 | `activeSourceFromContext` resolve via `active_source` (config de M-02), não via default de ctx | diff mostra a função consultando o repo/lookup de `active_source` (M-02) em vez de `ParseActiveSource(defaultParam)` hardcoded | L1 ran | F-03 |
| M03-C14 | Fail-closed preservado: tenant sem `active_source` não cai em default silencioso de xlsx | teste: tenant sem row em `active_source` → leitura retorna `ErrUnknownActiveSource`, nunca uma lista vazia disfarçada de "sem produtos" nem fallback mudo a `xlsx` | L1 ran | F-03 |
| M03-C15 | XlsxAdapter implementa `ProductSourceAdapter` (port de M-02) e compila | `go build ./...` verde; tipo concreto satisfaz a interface (read-side de F-02 + `Sync()`/`Kind()`) | L1 ran | F-04 |
| M03-C16 | `Kind()` do XlsxAdapter retorna exatamente `upload_snapshot` | leitura direta do código/teste unitário: `adapter.Kind() == SourceKind("upload_snapshot")` | L1 ran | F-04 |
| M03-C17 | `Sync()` é idempotente: 2 chamadas com o mesmo snapshot não duplicam rows no mirror | teste: `Sync()` chamado 2× sobre o mesmo protocolo → `SELECT count(*) FROM products_mirror WHERE codigo_produto IN (...)` inalterado entre a 1ª e a 2ª chamada (upsert por PK, não insert) | L1 ran | F-04 |

## Anti-critérios (falha se presente)

- AC-01: query tenant-scoped sem `tenant_id` (mirror read/write, `sync_state` write — ambas
  tocadas por esta milestone).
- AC-02: payload de provider (ML) dentro do adapter xlsx — fora de escopo, xlsx não fala com ML;
  qualquer hit disso é vazamento de responsabilidade.
- AC-03: campo ausente (custo/estoque/etc.) gravado como `0`/default em vez de `NULL` (M03-C6).
- AC-04: `parser.go` editado (M03-C1) — mesmo edição "cosmética"/refactor menor conta como
  violação; KEEP é absoluto nesta milestone.
- AC-05: cascata de delete físico em `products_mirror` ou `product_links` a partir do
  upsert-merge (M03-C3/C4) — qualquer `DELETE FROM products_mirror` fora de rollback de migração
  é falha automática.
- AC-06: `Collect(`/execução de coleta de mercado chamada pelo hook (M03-C9) — boundary MC-11 da
  missão; violação aqui reprova a missão inteira, não só M-03.
- AC-07: push para remote (profile §9).
- AC-08: `.env*` lido/impresso; credencial exposta (não aplicável a fluxo xlsx puro, mas mantido
  por consistência do gate).

## Nota de escopo (não confundir com MC-01 pleno / M-05)

M03-C8 prova que o hook DISPARA a geração de candidato — a REGRA de auto-aprovação
EAN-exato-único (o que acontece dentro do gerador/resolver) é M-05, não M-03. Se o teste de
M03-C8 encontrar candidatos gerados mas nenhum auto-aprovado, isso é esperado e correto nesta
milestone — `could-not-run`/`assumed` não se aplica aqui, é escopo explicitamente de outra
milestone (ver `mission.md` §Milestone Strategy M-05, Dep: M-02+M-03).

MC-01 pleno ("alimentado por AMBOS adapters") só fecha com M-04 também mergeado — M-03 sozinho
prova só o lado xlsx (M03-C2), o `WHERE source IN ('xlsx','sankhya')` completo é
`could-not-run` nomeado aqui até M-04 landar, igual à nota equivalente em M-02.
