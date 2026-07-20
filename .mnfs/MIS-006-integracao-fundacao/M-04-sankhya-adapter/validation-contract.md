# Milestone Validation Contract — M-04-sankhya-adapter

```yaml
id: M-04-VC
type: milestone-validation-contract
mission: MIS-006
milestone: M-04
validation_level: QA-0
base_sha: 138aac3d
```

Verdicts binários. Evidência = caminho inspecionável concreto (core §5, ladder L0-L4). Tipos:
`ran` (executado, output salvo), `assumed` (design, não executado), `could-not-run` (bloqueado —
nomear). Nenhum seam contra dependência real (Oracle Sankhya) provado por stub sem autorização
explícita (AC-04).

## Critérios

| ID | Critério | Prova mínima inspecionável | Ladder | Feature dona |
|----|----------|----------------------------|--------|--------------|
| M04-C1 | `[TESTAR-SKW]` mapping confirmado via db-consult ANTES de qualquer código de sync escrito | documento de mapeamento (coluna Sankhya → campo E2.1) anexado como evidência, com resposta do especialista Oracle (relaiada via hub REQUEST) para os 6 itens do checklist de `milestone.md` §Prerequisites; timestamp/commit do sync entrypoint POSTERIOR ao timestamp do mapeamento | L1 ran / **could-not-run nomeado** se db-consult não retornou | F-01 |
| M04-C2 | Se `[TESTAR-SKW]` bloqueado, milestone inteira marcada `could-not-run` nomeada — nenhum código de sync escrito sobre suposição | se M04-C1 = could-not-run, todos os critérios subsequentes (C3-C9) herdam `could-not-run`, nomeando o bloqueio explicitamente (nunca implementação silenciosa sobre hipótese não ratificada) | L0 (nomeação honesta) | F-01 |
| M04-C3 | `SankhyaAdapter` reusa `oracle/reader.go` como core (não reescreve) | diff mostra `FindProductsForLinking`/`GetSellableStock`/`GetCostAsOf`/`GetTaxInputs` intactos (zero mudança de assinatura ou lógica interna); novo método de sync é ADIÇÃO, não substituição | L1 ran | F-01 |
| M04-C4 | Sync entrypoint escreve `products_mirror` com `source='sankhya'`, validado contra Oracle REAL (dev stack), nunca stub | `SELECT * FROM products_mirror WHERE source='sankhya'` executado pós-sync contra dev stack real retorna rows; nenhuma execução contra double/fixture Oracle nesta prova (AC-04) | L2 ran / could-not-run nomeado se dev stack indisponível | F-01 |
| M04-C5 | Upsert-merge keep-absent replicado no lado Sankhya (mesma regra ADR-04 do xlsx) | 2 rodadas de sync, produto ausente na 2ª → row permanece, `absent_in_last_snapshot=true`+`stale_since` set, `product_links` associado intacto (se existir) | L2 ran | F-01 |
| M04-C6 | NULL honesto no lado Sankhya: campo não resolvível pelo mapeamento grava `NULL`, nunca `0` | produto com `TGFCUS`/`TGFTAB`/`TGFBAR` ausente/não-resolúvel → `custo`/`preco_venda`/`estoque_total`/`ean` `NULL` no mirror; grep de hardcode `0`/default nessas colunas no código de sync = zero hits | L2 ran | F-01 |
| M04-C7 | `Kind()` retorna `live_read_through` para `SankhyaAdapter`, port compila | `go build ./...` verde; chamada a `Kind()` no `SankhyaAdapter` retorna `live_read_through` (teste unitário ou leitura direta do valor constante) | L1 ran | F-02 |
| M04-C8 | `Sync()` implementado retornando `SyncResult` refletindo contagem processada/erros do F-01 | teste unitário/integração: `Sync(ctx)` retorna `SyncResult{Processed: N, Errors: M}` compatível com o resultado real da rodada | L1-L2 ran | F-02 |
| M04-C9 | Read-side do `SankhyaAdapter` inalterado (assinaturas de `readports.Reader`) | diff mostra zero mudança de assinatura em `FindProductsForLinking`/`GetSellableStock`/`GetCostAsOf`/`GetTaxInputs`; consumidores existentes (pricing/vínculo) compilam sem alteração | L1 ran | F-02 |
| M04-C10 | `canonicalKey()` estendida para incluir `source` (3 valores: xlsx, catalogo_cliente, sankhya) NO MESMO commit que adiciona sankhya ao domínio | `git show --stat <commit>` mostra `cache.go` (`canonicalKey`) E a extensão do domínio de source no MESMO commit; diff mostra `source` presente na composição da chave | L1 ran | F-03 |
| M04-C11 | Must-fail test prova isolamento cross-source (xlsx vs sankhya) na cache | teste: cache populado para `(tenant, xlsx)` NÃO retorna para `(tenant, sankhya)`; teste FALHA se `source` for removido da key manualmente (prova de que o teste é load-bearing, não decorativo — lição `chip-import-fix-closed`) | L1 ran | F-03 |
| M04-C12 | Cache-key extension nunca split de outra mudança — landam juntas | mesma evidência de M04-C10 (commit único); se landarem separados, critério FALHA mesmo que ambas as metades existam eventualmente | L0 ran (diff) | F-03 |
| M04-C13 | Nenhuma query tenant-scoped sem `tenant_id` no sync entrypoint | grep/leitura do SQL do sync entrypoint mostra filtro `tenant_id` presente em toda query nova | L0 ran (grep/leitura) | F-01 |
| M04-C14 | Credenciais Oracle nunca lidas/impressas durante a prova (`.env*`, secrets) | transcript/evidência da rodada de validação não contém valor de credencial Oracle; leitura de var de ambiente feita via `printenv VAR` pontual se necessário, nunca dump completo | L0 (revisão de evidência) | F-01 |

## Anti-critérios (falha se presente)

- AC-01: query nesta milestone sem `tenant_id` (M04-C13).
- AC-02: payload/schema Oracle vazando fora do adapter (`internal_read/adapters/oracle`) — outros
  módulos não devem importar tipos internos do reader Oracle.
- AC-03: campo não resolvível pelo mapeamento Sankhya virando `0`/default em vez de `NULL`
  (M04-C6).
- AC-04: sync entrypoint validado contra stub/fixture Oracle em vez de dev stack REAL, sem
  autorização explícita do hub (M04-C4) — se dev stack indisponível, `could-not-run` nomeado é
  o único caminho honesto, nunca um stub silencioso apresentado como prova.
- AC-05: `.env*` lido/impresso; credencial Oracle exposta em qualquer evidência (M04-C14).
- AC-06: push para remote (profile §9).
- AC-07 (específico M-04): código de sync escrito ANTES do `[TESTAR-SKW]` db-consult retornar —
  viola o bloqueio explícito de `milestone.md`; se detectado, milestone inteira é rejeitada
  independente de os testes passarem (mapeamento não-ratificado pode estar silenciosamente
  errado e nunca seria pego por teste unitário que também assume o mapeamento errado).

## Nota de escopo (não confundir com MC-01/MC-09 plenos da missão)

M-04 sozinho prova o lado Sankhya de MC-01 (feeding real) — o critério de missão pleno
("AMBOS adapters") só fecha em conjunto com M-03 (xlsx). MC-09 ("validado contra Oracle REAL")
é majoritariamente owned aqui: se `[TESTAR-SKW]` nunca for confirmado dentro do ciclo da missão,
MC-09 fecha como `could-not-run` nomeado — isso é um resultado válido e honesto do processo, não
um gap a esconder ou contornar com stub (AC-04, AC-07).
