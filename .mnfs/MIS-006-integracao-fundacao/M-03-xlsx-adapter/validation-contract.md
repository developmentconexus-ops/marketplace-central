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
| M03-C13 | `activeSourceFromContext` não carrega mais default silencioso de fonte; a resolução de-record vem do `active_source` via routing.Reader de M-02 (RECONCILIADO pelo hub D-120 — ver nota abaixo da tabela) | diff mostra (a) o default silencioso `return erpdomain.SourceXLSX` de `reader.go:47` REMOVIDO e nenhuma resolução `ParseActiveSource(default)` hardcoded restante, (b) evidência de que o lookup-de-record é `routing/reader.go:46` (`tenant_config.Repository.Get`) + pin em `:53` (`WithActiveSource`), já landado por M-02; grep prova que nenhum caller que contorna o routing dependia do default xlsx antes da remoção | L1 ran | F-03 |
| M03-C14 | Fail-closed preservado: tenant sem `active_source` não cai em default silencioso de xlsx | teste: tenant sem row em `active_source` → leitura retorna `ErrUnknownActiveSource`, nunca uma lista vazia disfarçada de "sem produtos" nem fallback mudo a `xlsx` | L1 ran | F-03 |
| M03-C15 | XlsxAdapter implementa `ProductSourceAdapter` (port de M-02) e compila | `go build ./...` verde; tipo concreto satisfaz a interface (read-side de F-02 + `Sync()`/`Kind()`) | L1 ran | F-04 |
| M03-C16 | `Kind()` do XlsxAdapter retorna exatamente `upload_snapshot` | leitura direta do código/teste unitário: `adapter.Kind() == SourceKind("upload_snapshot")` | L1 ran | F-04 |
| M03-C17 | `Sync()` é idempotente: 2 chamadas com o mesmo snapshot não duplicam rows no mirror | teste: `Sync()` chamado 2× sobre o mesmo protocolo → `SELECT count(*) FROM products_mirror WHERE codigo_produto IN (...)` inalterado entre a 1ª e a 2ª chamada (upsert por PK, não insert) | L1 ran | F-04 |

## Nota de reconciliação C13 (hub ruling D-120, 2026-07-22)

Defeito de plano corrigido no artefato (mission-planning reconciliation rule — conflito
resolvido ratificando NO plano, não empurrado à execução). C13 na redação original
("`activeSourceFromContext` consultando o lookup de `active_source` DENTRO de internalread") é
**arquiteturalmente impossível E redundante**:
- **Impossível**: `tenant_config` já importa `internalread` (`tenant_config/active_source.go:11`,
  alias de `ErrUnknownActiveSource`). `internalread` importar `tenant_config` = ciclo de import.
- **Redundante**: M-02 já faz o lookup em `routing/reader.go:46` (`tenant_config.Repository.Get`,
  fail-closed) e fixa o valor no ctx via `WithActiveSource` (`routing/reader.go:53`). O
  `activeSourceFromContext` de internalread apenas LÊ o valor já fixado — o routing.Reader é o
  lookup-de-record.

Leitura ratificada de F-03: satisfeita por (a) o valor observado por `activeSourceFromContext`
ser alimentado pelo `active_source` via pinning do routing de M-02 (JÁ verdadeiro), MAIS (b)
remover o default silencioso `SourceXLSX` (`reader.go:47`) para que um ctx não-fixado falhe
fechado (intenção real do EARS de F-03 + C14). Constraint da remoção: o chip DEVE provar por grep
que nenhum caller que contorna o routing dependia do default xlsx antes de removê-lo (o routing
sempre fixa o ctx; se existir caller bypass, surfar ao hub antes de remover). Atualizar também o
comentário de doc de `WithActiveSource` (`reader.go:25-26`) que hoje diz "defaults to xlsx so the
demo opens on real data" — passa a refletir fail-closed.

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

## Critérios de user-drive (AMENDMENT D-120 — obrigatório, ratificado pelo operador)

Origem: a validação de onda-1 "como usuário" pegou uma regressão que TODOS os gates de código
perderam (/catalogo 503 — capability opcional apagada por wrapper; hub-fix @2567eb44). Regra
ratificada pelo operador (2026-07-22): todo contrato de validação passa a exigir, além dos
critérios acima, dirigir o programa NO DEV STACK COMO USUÁRIO (browser real + HTTP real, perfil
limpo) cobrindo as telas EXISTENTES que tocam o seam da milestone — nunca só a superfície nova
do changeset. Verdict binário; evidência = página renderizada (screenshot/page-text) + status
HTTP + amostra de body. could-not-run só com bloqueio nomeado.

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M03-U1 | /catalogo renderiza produtos reais mirror-backed nas DUAS fontes upload: flip xlsx/catalogo_cliente via PUT /config/active-source muda o dataset na tela (as_of/conteúdo), sem 503 e sem lista vazia disfarçada | browser drive + 2 GETs com as_of distintos |
| M03-U2 | /integracoes: upload de xlsx real cria protocolo novo SEM apagar os anteriores (lista completa visível); re-upload do mesmo arquivo mostra 409 visível ao usuário | browser drive + lista de protocolos antes/depois |
| M03-U3 | Jornada pós-import: candidatos de vínculo aparecem na tela de vínculos existente sem nenhuma chamada manual extra (hook interno provado na UI, não só no DB) | browser drive da tela de vínculos pós-import |
| M03-U4 | Zero regressão nas telas existentes que leem produto: /catalogo, /precos, /anuncios, /pedidos carregam sem erro novo com a leitura mirror-backed | browser drive (4 telas) + console/network sem erro novo |
