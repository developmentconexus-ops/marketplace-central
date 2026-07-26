# CHIP-ANCHORS — contrato de validação

```yaml
chip: CHIP-ANCHORS
base_sha: 917f7bb58e385847fba5612201823f9db48791c6
level: QA-0
```

Critério só fecha com evidência `file:line` ou saída de comando citada. "Rodei e passou" sem
artefato = não fechado. Critério que o chip não conseguir provar volta como `BLOCKED`/defer
explícito, nunca como PASS otimista.

## C — código (chip prova)

| ID | Critério | Prova mínima |
|----|----------|--------------|
| C1 | O lado provider **declara** as âncoras de identidade que fornece, pendurado no registro de capability existente (`connectors/ports/marketplace_capability.go:9-20`), e `product_links` consome por **port próprio** com wiring na composition | `file:line` da declaração, do port, do adapter e da linha de wiring |
| C2 | `mandatoryUnavailableReasons()` não existe mais e **nenhum nome de âncora** está hardcoded no gerador | grep do símbolo = 0 hits; grep de `"marca"`/`"refforn"` em `product_links/application` = 0 hits em código de produção |
| C3 | R2 provada: zero ramificação por provider dentro de `product_links` | grep por `mercado_livre` / `ProviderCode ==` em `product_links/**` (produção) = 0 hits, com a lista dos arquivos varridos |
| C4 | `UNAVAILABLE` distingue "provider não fornece essa âncora" de "provider fornece e este anúncio não tem valor", **pelo `detail`** | teste com dois casos, um de cada, asserindo os dois `detail` distintos |
| C5 | R1 provada: OpenAPI e `packages/sdk-runtime` byte-idênticos ao BASE-SHA | `git diff --stat <base>..HEAD -- <caminhos>` vazio, colado na evidência |
| C6 | **D-A**: `50cm` vs `500MM` não é contradição; `50cm` vs `40cm` continua sendo | teste com os dois pares + **must-fail** (reverter a canonicalização deixa o teste vermelho) |
| C7 | **D-A regressão**: o hard-negative continua bloqueando o caminho corroborado (kit/combo/cor/voltagem seguem `REJECT` com CODPROD+EAN concordantes) | teste do caminho corroborado + `REJECT` preservado |
| C8 | **D-B**: candidatos, vínculos e audits têm limites independentes; fixture com 29 vínculos devolve 29 | teste com fixture ≥29 asserindo `len(links)==29` + **must-fail** (religar o limite compartilhado deixa vermelho) |
| C9 | R3 provada: nenhum `UPDATE`/backfill de motivo já persistido, nenhuma migração | grep por `UPDATE` em motivos = 0; `git diff --stat` sem arquivo em `migrations/` |
| C10 | Non-scope respeitado: **zero** arquivo `apps/web/` no diff | `git diff --name-only <base>..HEAD` colado inteiro |
| C11 | Ladder L0+L1 por profile §2 (GOCACHE absoluto) + governance de worktree limpo detached com BaseSha 40-hex | saídas citadas; falha pré-existente = citar allowlist, não re-provar |
| C12 | Grant de `root.go` foi **additive-only** | diff da região citado no payload do CLOSED |

## U — user-drive (rodado pelo HUB na stack, conta ML conectada)

O chip não dirige browser e não sobe stack. Estes critérios são a condição de aceite do hub;
o chip só precisa deixá-los alcançáveis.

| ID | Critério | Como o hub prova |
|----|----------|------------------|
| U1 | Depois de regerar candidatos, **nenhuma** linha da fila `/vinculos` mostra motivo `marca` ou `refforn` afirmando que o provider não fornece, quando isso é suposição do núcleo e não declaração do adapter | fila real na tela + inspeção do payload |
| U2 | `SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO` deixa de ser reprovado por contradição de dimensão contra o produto ERP em `cm` | linha na tela antes/depois + confiança/banda resultante |
| U3 | `/vinculos` → Resolvidos mostra os **29** vínculos (não 20), e a KPI concorda com o `count(*)` do banco | tela + `SELECT count(*)` do Postgres do dev stack |

Nota de honestidade sobre U3: o guard de truncamento da lista (dizer "200+" quando a página está
cheia, como a fila já faz) é **do M-06**. Este chip faz o número certo chegar; a tela que nunca
mente sobre estar truncada é a milestone seguinte.
