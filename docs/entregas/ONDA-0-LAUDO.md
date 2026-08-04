# Onda 0 — Laudo técnico

**Escopo auditado:** F-A1 (token de ML invisível), F-A2 (lote abortado), F-A3 (coleta de
mercado periódica + idade visível), F-00 (scheduler de pedidos), mais o P2.b
(imposto ex-ante) que entrou na mesma janela de merge.

**Base de comparação:** `bf2ccf3e` (tip da `main` antes dos merges desta onda).
**Ponta auditada:** `be6db3f` (com os merges `e66ce013` e anteriores).
**Data da auditoria:** 2026-08-03/04.
**Companheiro operacional:** [`ONDA-0-ENTREGA.md`](ONDA-0-ENTREGA.md) — a Ficha, campo a
campo, na língua do operador. Este laudo é o veredito técnico sobre os seis eixos do
[método](../METODO-DE-REVISAO.md).

**Veredito de uma linha:** a onda entregou o que prometeu no encanamento e **não** entregou
valor visível na tela. Quatro dos cinco jobs de sync rodam e aparecem; o quinto
(`market_queue`) nunca rodou. Nenhum pedido tem margem, nenhum anúncio tem margem, e 29 de
34 células de preço não dizem nada — nem que não sabem.

---

## Eixo 1 — Arquitetura e fronteiras

**Camadas.** Nenhuma violação nova de `GOV_MODULE_LAYER` ou `GOV_MODULE_DEPENDENCY` no diff
de conjunto (abaixo). A onda **removeu** três: `listings-sync` deixou de ser importada por
`listings/application/backfill.go` e `listings/composition/scheduler.go`, e o
`GOV_MODULE_LAYER` do scheduler saiu junto. Esse é o sinal certo — o encanamento novo
entrou pela composição, não por atalho lateral.

**Payload de provedor.** Os DTOs do ML continuam parando no adaptador. O `pricing/domain`
recebe grandezas nossas, não `sale_fee` cru.

**`tenant_id`.** Presente nas consultas novas conferidas neste laudo — inclusive na de
resumo de anúncios (`listings/adapters/postgres/repository.go:336`, `WHERE
l.tenant_id=$1 AND l.installation_id=$2`).

**Oracle somente leitura, Postgres canônico.** Respeitado. O `ICMSMatrixReader` lê Oracle e
escreve espelho; não há caminho de escrita para o `METALPRD`.

**Atomicidade de contrato.** `GOV_API_SDK_SPLIT` não disparou. Mas a checagem que existe
pergunta "mexeram nos dois arquivos no mesmo diff?", não "os dois dizem a mesma coisa?" — e
a busca 3 do eixo 3 mostra que **13 `operationId` do OpenAPI não têm homônimo no
`sdk-runtime`**, incluindo seis renomeações silenciosas. Isso passa pelo gate hoje.

### Governança por diff de conjunto

O comando declarado no método (`harness.ps1 -Command governance -BaseSha ...`) **não computa
diff nenhum** — `-BaseSha` alimenta exatamente um check, o `GOV_API_SDK_SPLIT`
(`Policy.psm1:446-461`); todo o resto é a lista absoluta do HEAD. Registrado como **D-25**.
O diff desta onda foi feito à mão: worktree destacado na base, **o mesmo instrumento
copiado nos dois lados**, e diff das triplas `(error_code, id, path)`.

```
HEAD  59 issues   55 únicos
BASE  62 issues   56 únicos

NOVOS (HEAD \ BASE): 2
  GOV_MIGRATION_PREFIX  0093   0093_icms_matrix.sql, 0093_orders_status_details_nullable.sql,
                               0093_sync_state_market_queue_entity_split.sql
  GOV_PRODUCTION_PANIC  panic  apps/server_core/internal/modules/orders/adapters/pricingtax/reader.go

SUMIRAM (BASE \ HEAD): 3
  GOV_MODULE_DEPENDENCY listings-sync   listings/application/backfill.go
  GOV_MODULE_DEPENDENCY listings-sync   listings/composition/scheduler.go
  GOV_MODULE_LAYER      listings-sync-adapters  listings/composition/scheduler.go
```

**Saldo: −1.** Os dois novos são F-4 (o `panic` em `reader.go`) e F-11 (três migrações
disputando o prefixo `0093`). Nenhum dos dois estava na base; ambos entraram nesta onda e
ambos são reais.

Duas coisas tiveram que ser consertadas para que essa medição existisse, e as duas são de
classe:

- **D-24 — a varredura lia checkouts de outros branches.** `Get-SourceFiles`
  (`Policy.psm1:193-204`) excluía diretórios com regex ancorada na raiz, então
  `.worktrees/<branch>/` e `.claude/worktrees/<branch>/` entravam inteiros. Antes da
  correção o mesmo diff dava **266 "novos"**, dos quais 264 vinham de
  `.claude/worktrees/f00-scheduler-pedidos/`. O gate estava **vermelho** — o lado que parece
  seguro — e mesmo assim inútil: ninguém acha 2 achados verdadeiros dentro de 266.
- **D-23 — o `.dockerignore` mandava 939,9 MB por rebuild.** A linha 2 era `.gocache`
  ancorada na raiz, e não pegava `apps/server_core/.gocache` (15.511 arquivos). Corrigida
  para `**/.gocache`, como os vizinhos `**/.gomodcache` e `**/node_modules` já eram.

A lição comum: **exclusão de instrumento não se ancora na raiz.** O que se exclui é *tipo de
diretório* (derivado, vendorizado, checkout alheio), e tipo aparece em qualquer profundidade.

---

## Eixo 2 — Verdade dos dados (ADR-17)

Nenhum `coalesce` novo sobre coluna de fato de provedor, nenhum `?? 0` novo no FE sobre
campo do SDK. Nesse sentido a onda está limpa: **ela não fabricou número.**

O que ela deixou é o outro lado do mesmo eixo — o desconhecido que não se anuncia:

| Achado | Veredito | O que a tela faz com o desconhecido |
|---|---|---|
| F-13 | **MUDO** | 29 de 34 células de PREÇO ficam em branco. `NoPriceEvidence` e `Stale` renderizam nada, e nada é indistinguível de "está tudo bem" |
| F-7 | **MENTIRA** | "Com erro de sync **0**" e "Desatualizados **0**" com `sync_state` constante em `synced`: a tela afirma ausência de erro quando o que ela sabe é que não mede erro |
| F-14 | **MUDO** | 7 anúncios `under_review` no `Total` e em nenhum contador; 34 = 9 + 18 + 7 invisíveis |
| F-5 | **MENTIRA** | cabeçalho do `/mercado` diz "coletado 03/08, 21:50" sobre uma tabela em que 33 de 34 linhas dizem "idade desconhecida" |
| F-9 | ÓRFÃO | `net_amount`, `margin_pct` e `decomposition` nulos em 39 de 39 |

`0` suprimido não foi encontrado — o outro lado do ADR-17 está respeitado.

**A distinção que este eixo cobra:** F-7 e F-5 são piores que F-13 e F-14. Célula em branco
não afirma nada; "0" e "coletado às 21:50" afirmam. Uma tela que cala sobre o que não sabe é
uma dívida de usabilidade; uma tela que **afirma** o que não sabe é um defeito de verdade.
Os dois primeiros viram defeito aberto; os dois últimos, dívida nomeada.

---

## Eixo 3 — Máximo local, redundância e código morto

As cinco buscas rodaram com número (detalhe em `ONDA-0-ENTREGA.md` §6):

| Busca | Resultado |
|---|---|
| 1 — motor duplicado | **0.** O risco histórico (segundo motor de preço ao lado de `pricing/domain`) morreu no re-escopo do P2.b, 17 tasks → 7 |
| 2 — fórmula sem consumidor | **3** de 98 funções exportadas de `*/domain/` sem chamador |
| 3 — operação de contrato sem consumidor | **50** de 111 `operationId` sem consumidor no `apps/web`; **13** sem homônimo no `sdk-runtime`; **6** renomeações silenciosas (F-10) |
| 4 — campo sem produtor | F-2, F-3, F-7, F-9 — já contabilizados por tela |
| 5 — abstração que não abstrai | **0 cerimônias.** Uma exceção nomeada: `ICMSMatrixReader`/`ICMSMatrixWriter` (F-8) tem porta, implementação e teste e **nenhum chamador de produção** — não é abstração vazia, é abstração **desligada** |

**F-8 é o achado mais caro da onda.** A cadeia:

```
Task 7 do P2.b nunca executada
  -> icms_matrix_mirror vazia
  -> MatrixCellFor não acha célula para nenhum (UF, grupo)
  -> reader.go all-or-nothing: ICMS saída, DIFAL, PIS/COFINS e restituição ST = nil no pedido inteiro
  -> BuildProfitability não fecha os sete componentes
  -> Margem e Retorno líquido nil em 39 de 39 pedidos
```

O código do imposto ex-ante está escrito, testado e correto. Ele só nunca foi ligado. **A
onda entregou T1–T6 de sete tasks e a que faltou é a que produzia o valor.** Isso não é
código ruim — é entrega incompleta declarada como completa, que é pior de achar depois.

**Corte YAGNI aplicado:** nada foi apagado nesta onda. As 50 operações sem consumidor e as 3
fórmulas mortas viram dívida nomeada, não deleção — YAGNI vale contra o que ainda não tem
consumidor, e algumas dessas operações são consumo previsto da Onda 1. A regra é decidir na
abertura da Onda 1, com a Ficha prospectiva na mão.

---

## Eixo 4 — Execução real (live drive)

Executado nas quatro telas. Detalhe completo em `ONDA-0-ENTREGA.md` §7; aqui só o que o
eixo cobra.

**Proveniência do binário** (a regra que existe porque já falhou aqui):

```
container   0dc9cb1f1db2   marketplace-central-backend-1
Created     2026-08-03 21:35:23 -03   (= 2026-08-04T00:35:23Z)
último commit em apps/server_core e apps/web   e66ce013   2026-08-03 19:29:22 -03
```

Container **2h06 mais novo** que o código. A primeira medição da noite pegou o oposto — a
stack de pé estava 2h29 **atrás** do merge, com `schema_migrations` em 82 contra 83 arquivos
e a faltante sendo exatamente `0093_icms_matrix.sql`. Sem esse passo o drive teria medido a
onda anterior e assinado como onda 0.

**Vermelho antes do verde, com controle negativo nomeado.** O controle negativo desta onda é
o próprio par de medições do binário: a mesma tela, sob o binário velho, não tinha a linha
`Orders (incremental)` no `/sync/health`; sob o binário novo ela aparece com
`last_incremental_at=2026-08-04 00:52:05Z`, **17 minutos depois** do `Created` do container.
O verde só existe no mundo com a mudança. **F-00 está provado vivo.**

**Tela verde não prova visibilidade de falha.** Confirmado por medição: `sync_state` é
constante `synced` em 34 de 34 anúncios, então os contadores de exceção da tela nunca foram
exercidos contra uma falha real. Nada nesta onda provou que uma falha de sync apareceria.

**Um defeito falso foi evitado, e a regra que o evitou vale registrar.** A linha de Orders
dizia "há 2 h" na tela enquanto o `/sync/health` consultado direto no mesmo minuto já
devolvia `00:52:05Z`. A tela mostrava a idade do momento do fetch. **Antes de acusar a tela,
pergunte à API que a alimenta.**

---

## Eixo 5 — Lanes e evidência

Quatro lanes na ordem do método (§6), todas com contagem por linha.

**1. `go build`**

```bash
cd apps/server_core && GOCACHE="$PWD/.gocache" go build ./...
```
```
BUILD_EXIT=0
```

**2. `go test`**

```bash
cd apps/server_core && GOCACHE="$PWD/.gocache" go test ./...
```
```
ok            120
no test files  48
FAIL            0
```

**3. `npm run test --workspace @marketplace-central/web`**

```
Test Files  72 passed (72)
     Tests  601 passed (601)
  Duration  24.59s
EXIT=0
```

Este número precisa de uma nota, porque a primeira execução **reprovou** e a diferença não
foi o código:

```
1ª execução (concorrente com a lane de integração)
  Test Files  3 failed | 69 passed (72)
       Tests  3 failed | 598 passed (601)
    Duration  54.36s
  ListingDetailPanel > loads one detail without refetching the listing page  — timed out in 5000ms
  ImportChainPanel > renders a known protocol verbatim                       — timed out in 5000ms
  SyncHealthCard.realNetwork > reaches the named ErrorState ...              — expected 0 to be >= 1

re-execução dos MESMOS 3 arquivos, sozinhos
  Test Files  3 passed (3)
       Tests  26 passed (26)
    Duration  5.07s   (o SyncHealthCard levou 1114ms contra o timeout de 5000ms)

re-execução da suíte INTEIRA, sozinha
  601 passed (601)   24.59s
```

Os três eram *timeout* sob contenção de máquina, não defeito. **A suíte de front é sensível
a carga**: 54s contra 25s para o mesmo trabalho, com o timeout padrão de 5000 ms parado a
poucos milissegundos das asserções mais lentas. Isso é fragilidade de lane, não da onda, e
está registrado como **D-26**. Foi resolvido por medição — rodar isolado — e não por
suposição de flake, que é o mesmo erro na direção confortável.

**4. `npm run harness:integration`**

A lane emite `status=passed` e **nada mais**:

```
target=ephemeral-postgres    migrations_first=83
key=MPC_TEST_DATABASE_URL    migrations_second=0
migrations=embedded          resource_count=0
container=ephemeral          status=passed        EXIT=0
```

`status=passed` **não prova execução** — `Invoke-HarnessPostgresLifecycle` só guarda a saída
do `go test` quando ele falha (`Postgres.psm1:418-423`), e o objeto de retorno não a expõe.
Uma corrida em que todo teste fosse pulado sairia byte-idêntica a esta. Registrado como
**D-26** junto com a fragilidade do front.

A contagem por linha teve que ser obtida replicando o ciclo da lane à mão — sessão do
harness de pé, banco `mpc_test_<32 hex>` (o formato que `testsupport/postgres/target.go:19`
exige; qualquer outro nome devolve `HPG_TARGET_INVALID`), `./cmd/testdb migrate`, e então
`go test -tags=integration -v`:

```
applied 83 migration(s)
MIGRATE_EXIT=0

=== RUN     55
--- PASS    55
--- SKIP     0
--- FAIL     0
ok  marketplace-central/apps/server_core/tests/integration  3.908s
TEST_EXIT=0
```

**As 83 migrações aplicam limpo em banco novo, e a segunda passada aplica 0** — idempotência
provada. É a resposta operacional ao F-11: a colisão de prefixo `0093` não quebra um banco
novo hoje. O que ela quebra é a garantia de ordem para a próxima colisão que tiver
dependência.

**5. Governança** — no Eixo 1 acima, por diff de conjunto (o comando declarado não computa
diff; ver D-25).

---

## Eixo 6 — Ordem de verdade e conflitos

Nenhum conflito de arquitetura ou de contrato precisou parar a linha. Três classificações
foram necessárias:

1. **Runtime vs. código.** A stack de pé discordava do repo. Ordem de verdade não decide
   isso — é medição. Resolvido reconstituindo o runtime, não reinterpretando o código.
2. **Nome de comando vs. comportamento.** `governance-drift` promete diff e entrega
   snapshot absoluto. O nome é a alegação, o código é o fato: **D-25**, aberta. Conserto
   proposto: ou o comando roda os dois lados e emite só `HEAD \ BASE`, ou vira
   `governance-snapshot` e o diff ganha comando próprio.
3. **Comentário vs. dado.** `pedidosTabs.ts:17-19` justifica a aba placeholder com "o
   dataset de demo é limitado a REVIEW, não há cancelados". O dado ao vivo tem 7 cancelados.
   Comentário é a camada mais fraca da ordem de verdade e apodreceu primeiro — mesma classe
   do "brief de milestone é alegação sobre o repo e apodrece".

**Defeito reincidente desta onda, com parada:** exclusão de instrumento ancorada na raiz
apareceu duas vezes no mesmo dia, em dois arquivos sem relação (`.dockerignore` e
`Policy.psm1`), custando uma hora de build e vinte e cinco minutos de varredura. Não virou
remendo local: as duas foram corrigidas para casar em qualquer profundidade e a classe está
registrada em D-23/D-24.

---

## Inventário de achados

| # | Achado | Severidade | Estado |
|---|---|---|---|
| F-00 | Scheduler de pedidos | — | **Provado vivo** no live drive |
| F-4 | `panic` em `orders/adapters/pricingtax/reader.go` | defeito | aberto, confirmado por diff de conjunto |
| F-5 | Cabeçalho do `/mercado` afirma coleta recente sobre 33/34 "idade desconhecida" | defeito (MENTIRA) | aberto |
| F-6 | Teto aritmético de renovação (100/h contra validade de 1 h) | propriedade | não morde hoje (4 produtos); morde acima de ~100 |
| F-7 | "Com erro de sync 0" / "Desatualizados 0" com `sync_state` constante | defeito (MENTIRA) | aberto |
| F-8 | Matriz de ICMS com leitor, escritor, teste e nenhum chamador; Task 7 do P2.b não executada | defeito | aberto — causa raiz de F-9 |
| F-9 | `net_amount` / `margin_pct` / `decomposition` nulos em 39/39 | ÓRFÃO | aberto, consequência de F-8 |
| F-10 | 50/111 operações de contrato sem consumidor; 13 sem homônimo no SDK; 6 renomeações silenciosas | dívida | nomeada |
| F-11 | Três migrações no prefixo `0093`; sem `0094` nem `0095` | defeito | aberto |
| F-12 | 7 pedidos cancelados na fila de trabalho, sem contador, com aba placeholder vazia | defeito | aberto |
| F-13 | Coluna PREÇO muda em 29/34 | dívida (MUDO) | nomeada |
| F-14 | 7 anúncios `under_review` fora de todo contador | dívida (MUDO) | nomeada |
| F-15 | `AppendPendingCodigos` concatena sem dedup (110 ids = 55 duas vezes) e nada drena `pending`; a Saúde do sync mostra fila e stream na mesma tabela | dívida | nomeada — **corrigido**: ver 7.4 da Ficha. Não explica o 1-de-34, que é D-53 |

Dívidas de harness abertas nesta auditoria: **D-23** (`.dockerignore`, resolvida), **D-24**
(varredura cross-branch, resolvida), **D-25** (`governance-drift` sem diff, aberta),
**D-26** (lanes sem contagem por linha; front sensível a carga, aberta) — todas em
`.mnfs/HARNESS-DEBTS.md`.

---

## O que o operador ganha hoje, em uma frase por tela

- **`/anuncios`** — vê os 34 anúncios e o preço de cada um. **Não** vê margem (0 de 30
  vinculados), não vê competitividade (1 de 34), não vê os 7 em revisão.
- **`/pedidos`** — vê os 39 pedidos e sabe em que etapa cada um está. **Não** vê retorno
  líquido nem margem em nenhum, e trabalha 7 cancelados como se fossem vivos.
- **`/mercado`** — vê a lista. **Não** tem com o que decidir: margem atual `—` em 34 de 34 e
  `Aplicar` desabilitado em 34 de 34. O 1-de-34 com evidência é D-53, não avaria.
- **`/integracoes`** — **esta funciona.** Diz honestamente quando cada job rodou. Única
  ressalva: põe a fila `market_queue` na mesma tabela dos streams, e o "nunca" dela — que é
  de projeto — se lê como avaria.

O caminho mais curto para o operador sentir a onda é **uma** coisa que já existe e está
desligada: a **Task 7 do P2.b** (F-8, destrava margem em 34 pedidos). Não é código novo de
regra — é composição, e é escopo declarado da própria fatia que entrou nesta onda.

Competitividade em escala **não** é dívida da Onda 0. É D-53, decisão medida e ratificada:
o job periódico renova, não descobre; a primeira evidência de cada produto nasce do clique
do operador. Descoberta periódica precisa de orçamento de rate limit que o bucket
compartilhado não tem hoje, e entra numa onda futura com fila priorizada.
