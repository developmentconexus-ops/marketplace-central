# DRIVE-EVIDENCE F-5 — três telas medidas pelo HUB no tip `57666417`

Medido pelo hub, stack real, dado real. Mecanismo §6: checkout dos arquivos concedidos do chip na
árvore do hub → dirige → reverte, `git status --porcelain` verificado dos dois lados (limpo antes,
limpo depois, `git diff --quiet HEAD` EXIT 0).

**Declaração de escopo, e ela limita tudo abaixo.** O empréstimo é PARCIAL por construção: o
frontend e o SDK vêm de `57666417`, o **backend rodando é o binário da `main`**. O diff do chip
toca `erp_import/transport/http_handler.go` (+24), que não está no container. Todo veredito abaixo
que depende de resposta do servidor vale contra a `main`, não contra o tip.

## Tela 1 — caminho feliz, `/importacoes/eac3ac9e-b87b-43be-959c-7bf1f11a07f1`

```
Detalhe da importação
Voltar para importações
Estado da importação

Três medidas independentes, lidas do servidor — não são etapas de um funil.

Protocolo #001-E

Linhas importadas
55
Produtos vinculados
0
Linhas na fila de sync
55

Fila lida em: 28/07/2026, 18:47
```

`GET /erp/imports/…/chain → 200`, zero erro de console. Byte-idêntico ao drive anterior de
`b91c7507` (o único delta renderizável entre os dois commits é `errorDetail`, ramo de erro).

**Correção do hub:** a captura que o hub mandou ao chip na mensagem anterior estava **incompleta** —
omitia a última linha do painel, `Fila lida em: 28/07/2026, 18:47`. O chip colou aquele bloco no
pack como verbatim e conferiu as seis cadeias por string; as seis existem, então nenhuma conclusão
muda. Mas o bloco era o painel MENOS uma linha renderizada, e foi rotulado "o que a tela renderiza".
Esta é a captura completa e substitui aquela.

## Tela 2 — 404, `/importacoes/00000000-0000-4000-8000-000000000000`

```
Estado da importação
Três medidas independentes, lidas do servidor — não são etapas de um funil.

Erro ao carregar. Importação não encontrada.

Tentar novamente
```

`GET …/chain → 404 {"error":"import_not_found"}`. Ramo `isImportNotFoundError` verdadeiro.

**Isto NÃO é a linha que o chip mudou.** O delta de `b91c7507..57666417` é o `errorDetail` de
FALLBACK (`"Não foi possível carregar o estado da importação."`), que só é alcançado quando
`error !== "import_not_found"`. O drive que o chip propôs para ver o ramo de erro — uuid válido
inexistente — bate no 404 e renderiza a outra string. Teria confirmado uma linha que ele não mexeu.

## Tela 3 — 5xx, `/importacoes/<id-malformado>`

`GET /erp/imports/nao-e-uuid/chain → 500 {"error":"internal_error"}`

```
Estado da importação
Três medidas independentes, lidas do servidor — não são etapas de um funil.

Carregando…
```

**O painel nunca sai de `Carregando…`.** Três ids distintos (`nao-e-uuid`,
`tambem-nao-e-uuid`, `terceiro-nao-uuid`), o último observado por ~40s. Em cada um: **uma** única
requisição ao backend, nenhum retry, nenhum `ErrorState`, nenhum botão "Tentar novamente", zero erro
de console. O usuário fica com spinner eterno e sem ação.

Atribuição, por string: o ramo de render é **idêntico na `main`** — `ImportChainPanel.tsx:46-50`
(`isPending ? LoadingState : isError || !data ? ErrorState`) e `useErpImportChain.ts:23-28`
(`retry` devolve `failureCount < 1` para 5xx) são os mesmos nas duas revisões. O diff do chip não
toca nenhum dos dois. **Defeito pré-existente da `main`, não deste chip** — registrado aqui para
que nenhum assento o cobre do diff.

Consequência para o critério do chip: no tip, um id malformado devolve `400 invalid_import_id`
(handler no diff dele), 4xx faz `retry` devolver `false`, a query liquida em erro e o fallback
renderiza. A string que ele mudou é **alcançável no tip** — mas só porque a mudança de backend dele
existe, e este drive não a exercitou. Cobertura dessa linha continua sendo `vitest`
(`ImportChainPanel.test.tsx:173,186`), declarada como tal, mais um drive futuro contra o backend do
tip.

## Achado lateral, fora de qualquer chip em voo

Na `main`, id malformado devolve `500 internal_error`. O comentário do SDK em
`erpImport.ts:76-80` já descreve `400 invalid_import_id` como "a resposta para `{id}` malformado …
existe para o chamador distinguir requisição ruim de falha de servidor, que um 500 cru escondia".
Essa prosa descreve o mundo do tip do chip, não o da `main`. Enquanto ANCHORS-3 não mergear, o
contrato afirma um 400 que o servidor não emite.
