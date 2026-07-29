# QA final da MIS-006 — dirigida como operador, na `main` fechada

Base: `main` @ `0851f26c` (missão fechada `975ac82d`). Stack do hub servindo o checkout
primário (`docker inspect` do backend confirma `…/marketplace-central -> /workspace`). Dado
real, sem fixture, sem mock.

Ordem do operador: *"no final fazer um full QA e validar realmente como usuário"*.

## O que esta rodada TEM que provar

Não é re-rodar teste. É abrir a tela e ver o que o software faz com o dado que existe. Cinco
alvos, cada um ligado a algo que esta missão mudou:

| # | tela / caminho | o que se prova | vindo de |
|---|---|---|---|
| Q1 | `/integracoes` | fonte ativa vem do BANCO por tenant (não de localStorage) e o toggle persiste após reload | M-02 |
| Q2 | `/importacoes` | rota existe; a cadeia importados/vinculados/enfileirados renderiza com número real; a legenda NÃO promete funil | M-06 + ANCHORS-3 |
| Q3 | `/vinculos` fila | **nenhum par (anúncio, produto) aparece duas vezes** — o defeito F-1 do live drive | CHIP-FIM item 2 |
| Q4 | `/vinculos` motivos | `side` separa `both` de `provider` na tela; nenhum detail alega ausência de valor que existe | VINC-NEUTRO + CHIP-FIM item 5 |
| Q5 | upload xlsx em `/importacoes` ou `/integracoes` | o POST roda sob o deadline de lote (120s), não os 15s interativos | CHIP-FIM item 1 |

## Regras da rodada

- Zero erro de console é critério, não observação.
- Toda contagem lida na tela é confrontada com SQL no mesmo banco antes de virar PASS.
- Achado novo NÃO reabre a missão: vira task com file:line, e o corte YAGNI continua valendo —
  só entra em chip o que o operador alcança.
- Nada de `.env`, nada de credencial, nada de push.

## Resultado — 5 PASS, medidos

### Q3 — o duplicado entre âncoras morreu (CHIP-FIM item 2)

O banco ainda carregava a geração ANTERIOR ao chip, então o defeito estava na tela antes de
qualquer ação: **9 linhas para 7 pares**, com `MLB4735378521 / 22467` a 40% e a 20% e
`MLB6896001442 / 42519` idem — exatamente o F-1 do live drive. Rodar a geração no binário novo
devolveu `generated_count=36`, e o banco fechou em **36 candidatos / 36 pares distintos**,
`GROUP BY … HAVING count(*)>1` vazio. Na tela, contado pelo próprio DOM: 7 linhas, 7 pares. Os
20% que sobraram são de produtos DIFERENTES (44975, 43534) — pares legítimos, não resíduo.

O sweep STALE (item 4) provou-se duas vezes no mesmo movimento: 38 → 36 sem sobra da rodada
velha, e mais tarde limpou 34 candidatos de uma rodada feita sob fonte errada.

### Q4 — os motivos separam os lados (VINC-NEUTRO + CHIP-FIM item 5)

Lidos dos `title` renderizados, em dado real: `ean (falta nos dois lados)` e
`ean (falta no anúncio)` aparecem como rótulos distintos na mesma fila. O `side` chega ao
operador; `both` e `provider` não colapsam mais numa frase só. Nenhum `detail` alega ausência de
valor que existe nesta base.

### Q2 — a cadeia de importação (M-06 + ANCHORS-3)

A legenda na tela é **"Três medidas da mesma importação, em duas unidades — nenhuma é etapa da
outra."** — a alegação falsa de independência saiu e a unidade está nomeada. Os três contadores
(55 · 0 · 55) foram reproduzidos por SQL rodando a mesma lógica do `GetImportChain`
(comparação numérica do codprod, `DISTINCT` no `internal_product_id`, expansão do
`cursor->'pending'`): bate exato.

O estado de erro: forçar a rota do detalhe contra um id inexistente renderiza
**"Erro ao carregar. Importação não encontrada."** com botão *Tentar novamente* — nada de
`Carregando…` eterno. 404 e 5xx caem na mesma branch (`chainQuery.isError`,
[ImportChainPanel.tsx:52](../../../apps/web/src/pages/importacoes/ImportChainPanel.tsx#L52)); o
`isImportNotFoundError` só escolhe a frase.

### Q1 — fonte ativa mora no banco (M-02)

Trocar o rádio em `/integracoes` gravou `active_source = xlsx / upload_snapshot` na tabela.
Depois de **limpar `localStorage` e `sessionStorage`** e recarregar, a tela voltou marcando
`xlsx` com `localStorage.length === 0`: o valor vem do servidor, não do navegador.

A prova mais forte veio de graça: com a fonte em `xlsx`, a geração de vínculos passou a devolver
34 candidatos `unresolved / NO_CANDIDATE / 0`, porque o lookup passou a ler o espelho xlsx de 55
linhas em vez das 10.529 do Sankhya. A fonte ativa comanda a leitura do app inteiro, não só o que
a tela desenha. Restaurada para `sankhya` e regenerado: 36/36, 32 `exact_sku` + 4 `conflict`.

### Q5 — `POST /erp/imports` roda sob o deadline de lote (CHIP-FIM item 1)

Medido no servidor em execução, não lido. O middleware de classe responde no `ctx.Done()`
qualquer que seja o bloqueio do handler, então um corpo enviado devagar discrimina as duas
classes de forma limpa:

| rota | classe | corpo em ~25s | resultado |
|---|---|---|---|
| `POST /erp/imports` | batch declarada | 55 linhas reais em .xlsx | **201 Created em 25,9s** |
| `POST /pricing/decompose` | interativa (default) | JSON pequeno | **conexão morta em 15,7s** |

Mesmo servidor, mesmo ritmo de envio, mesmo middleware. Nenhuma rota interativa sobrevive a 25s;
essa sobreviveu. O 15,7s do controle é o que torna o 201 atribuível à classe batch em vez de a um
middleware ausente — a primeira tentativa de controle usou
`/product-links/link-candidates/generations`, que [root.go:263](../../../apps/server_core/internal/composition/root.go#L263)
declara como batch, e teria produzido a leitura errada.

O upload entrou pelo caminho de verdade: `#002-E`, 55 aceitos / 0 rejeitados, visível em
`/importacoes` e no card ÚLTIMO IMPORT ERP da visão geral. O clique no seletor de arquivo abre
diálogo do sistema operacional, fora do alcance da ferramenta; o request que o clique produz foi
dirigido inteiro e o resultado conferido na tela.

### Varredura

`/`, `/anuncios`, `/catalogo`, `/estoque`, `/pedidos`, `/precos`, `/mercado`, `/integracoes`,
`/importacoes`, `/vinculos` e o detalhe da importação: **zero erro de console, zero request 4xx
ou 5xx**. Catálogo lê Sankhya ao vivo e mostra honest-unknown (`— (missing_stock)`) em vez de
zero inventado.

O `ANÚNCIOS ATIVOS 0` da visão geral foi conferido antes de virar achado e está **correto**:
`/listings/summary` responde `active:0, paused:27, total:34` — a conta está com os anúncios
pausados.

## Achado (1) — não bloqueia, decisão do operador

**`(unproved)` é anotação de dev em frase que o operador lê.**
[generation_service.go:524](../../../apps/server_core/internal/modules/product_links/application/generation_service.go#L524)
e [:577](../../../apps/server_core/internal/modules/product_links/application/generation_service.go#L577)
emitem `"ean corrobora o mesmo codprod (unproved)"` e `"ean corrobora codprod (unproved)"`. O
`detail` vira `title` do chip de motivo e texto do drawer — o golden test já assere
`title="ean: ean corrobora codprod (unproved)"`, ou seja, é string de tela pelo alcance novo da
R-24. A frase ainda se contradiz: afirma que corrobora e em seguida marca que não está provado.

Não aparece nesta base porque o tenant não tem candidato `exact_ean`, mas o caminho renderiza: o
`:577` é CONFIRM/60%, que fica na fila por construção.

Correção é de duas linhas e sem risco. Fica como decisão porque a missão está fechada e o corte
YAGNI vale: string na tela dispara a regra, mas quem decide se abre chip para duas frases é o
operador.
