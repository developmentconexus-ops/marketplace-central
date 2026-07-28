SLICE CARD — S4 · as três condições do dual gate P6

## Por que este slice existe

O gate P6 rodou com dois revisores independentes e cegos um ao outro (Opus read-only + GPT-5.6 Sol
medium). Os dois APROVARAM, sem bloqueante, e os dois escreveram a MESMA ressalva por conta própria.
Registro do hub em `.mnfs/MIS-006-integracao-fundacao/_hub-gate-import-chain/GATE-P6.md` @ `43d7ee5c`.

Três condições. Elas são delta pequeno e cirúrgico. Não reabra nada além delas.

## write_set (nada mais)

- `apps/web/src/pages/importacoes/ImportChainPanel.tsx`
- `apps/web/src/pages/importacoes/ImportChainPanel.test.tsx`

NÃO toque em `useErpImportChain.ts`, nem em `ImportacaoDetailPage.tsx`, nem em nada sob
`pages/vinculos/`, nem em `packages/**`. Em especial: **não conserte `formatDateTime`** — ele é
compartilhado por outras telas e outro dono; o guard é do lado do CONSUMIDOR.

## Condição 1 — `queue_read_at` recebe guard de TIPO

**O defeito, provado por execução pelo lado GPT do gate** (não é argumento, é medição):

```
formatDateTime(["2026-07-18T12:00:00Z"]) → 2026-07-18T12:00:00.000Z
formatDateTime(true)                     → 1970-01-01T00:00:00.001Z
formatDateTime(1)                        → 1970-01-01T00:00:00.001Z
```

`formatDateTime` (`packages/web-query/src/index.ts:114`) guarda `!value` e
`Number.isNaN(date.getTime())`. **Não checa `typeof`.** `new Date(true)` é uma data VÁLIDA, então o
`?? <UnknownValue/>` na linha 70 nunca dispara e a tela mostra `01/01/1970` como se fosse fato do
servidor. Um campo em tipo errado tem que virar `—`, não uma data que parece boa.

Isto é assimetria dentro do MESMO card: `renderCounter` checa `typeof value === "number"`,
`renderProtocol` checa `typeof value === "string"`, e `queue_read_at` não checa nada. Dois dos três
mecanismos checam tipo. Conserte o terceiro, no mesmo formato dos outros dois.

Escreva um helper local nomeado ao lado de `renderProtocol`, na mesma forma: valor é conhecido
quando é `string` não-vazia (depois de `trim`); aí sim passa por `formatDateTime`, e se `formatDateTime`
devolver `null` (data impossível de parsear) também é `<UnknownValue/>`. Qualquer outra coisa —
`boolean`, `number`, array, objeto, `null`, ausente — é `<UnknownValue/>` direto, sem passar pelo
`formatDateTime`. Mantenha o hint pt-BR que já existe na linha 70.

Não generalize `renderCounter`/`renderProtocol`/o novo helper num renderizador polimórfico. Três
helpers pequenos que dizem o que fazem valem mais que um esperto.

## Condição 2 — trava unitária do ZERO CONHECIDO

Nenhum dos 8 testes do painel assere um `0`. O comportamento está certo hoje e o hub observou ao vivo
(`vinculados = 0` renderizou `0`), mas nada no repositório trava isso: uma "simplificação" futura de
`renderCounter` para `value ? value : <UnknownValue/>` quebraria a metade do ADR-17 que o operador
mais lê — zero conhecido virando `—` diz "não sabemos" quando o servidor sabia e disse zero — e os
521 testes continuariam verdes.

Adicione UM teste: payload com `enfileirados: 0` (os outros contadores com valores diferentes de zero
e diferentes entre si) → `toHaveTextContent("0")` **e** `not.toHaveTextContent("—")` no
`erp-import-chain-enfileirados`. As duas asserções: a primeira sozinha passaria contra um `—0`
qualquer, a segunda é a que prova que não virou desconhecido.

## Condição 3 — 404 deixa de ser atribuído por STATUS

`ImportChainPanel.tsx:12` hoje:

```tsx
return candidate.status === 404 || candidate.error === "import_not_found";
```

O `status === 404` curto-circuita ANTES do segundo termo, então QUALQUER 404 — rota não montada,
`baseUrl` errado, proxy mal configurado — renderiza "Importação não encontrada.", afirmando um fato
sobre o DADO quando a verdade é que a ROTA não está lá. É a classe do catalog-503 falando com
confiança errada: o pior 404 possível (fiação quebrada) é justamente o que essa linha disfarça de
"importação inexistente".

**Remova o termo `candidate.status === 404 ||`.** Fica só a atribuição pelo corpo.

Isso é seguro e está verificado ponta a ponta: o servidor emite corpo achatado
`{"error":"import_not_found"}` (`apps/server_core/internal/modules/erp_import/transport/http_handler.go:114`
e `:127`, com o teste Go `http_handler_test.go:344` assertando o corpo exato), e o SDK lança
`{ status, error }` (`packages/sdk-runtime/src/index.ts:1715`). Um 404 de importação inexistente
SEMPRE carrega o campo `error`. Um 404 de rota inexistente não carrega.

Teste discriminante obrigatório: rejeição com `{ status: 404 }` e **sem** campo `error` → a mensagem
GENÉRICA ("Não foi possível carregar a cadeia da importação."), NÃO "não encontrada". O teste de 404
existente (com `error: "import_not_found"`) continua valendo e não pode ser enfraquecido — os dois
juntos são o que prova que a atribuição passou a ser pelo corpo.

ATENÇÃO ao tempo: o hook (`useErpImportChain.ts:23-28`) não faz retry em 4xx, então esse caso assenta
na hora, sem timeout estendido. Não copie o `{ timeout: 5000 }` que só o teste de 5xx precisa.

## O que o hub REJEITOU — não implemente

O lado GPT sugeriu apertar `Number.isFinite` para inteiro não-negativo (recusar `-1`, `1.5`). **O hub
rejeitou e você não deve fazer.** Inventaria regra de domínio que o contrato não tem, e faria um
servidor que mandou número real virar `—` — que é a OUTRA mentira do ADR-17. Contadores vêm de
`count(*)`; se vierem negativos o defeito é do backend e o FE mostra o que chegou.

## Testes: falhando PRIMEIRO, depois verde

Escreva cada teste antes do respectivo conserto e confirme que ele falha pela razão certa. Em
especial o da condição 1: contra o código atual ele deve falhar mostrando **1970**, não mostrando
vazio. Se falhar por outra razão, você escreveu o teste errado.

Não enfraqueça, não reescreva e não remova nenhum dos 8 testes existentes. Se algum quebrar, isso é
um ACHADO — reporte, não adapte a asserção.

## Verificação (rode cada uma UMA vez)

O worktree tem `node_modules` real (`npm ci`). Use o `apps/web/vitest.config.ts` de estoque.

- Typecheck, da raiz: `npx --no-install tsc --noEmit -p apps/web/tsconfig.json`
  **15 erros PRÉ-EXISTENTES** são esperados e NÃO são seus (`pages/produto/*`, `pages/anuncios*`,
  `pages/mutations/*`, `pages/vinculos/QueueRow.tsx`, `pages/vinculos/VinculoDrawer.tsx`). Só erro
  nos DOIS arquivos que você tocou é seu.
- Unitários, de `apps/web`: `npx --no-install vitest run src/pages/importacoes`
- Baseline a preservar: **65 arquivos / 521 testes** hoje. Você adiciona ~4 testes; nenhum arquivo
  que passava pode ficar vermelho. Reporte o número final que VOCÊ observou, não o previsto.

Se um comando genuinamente não rodar no seu sandbox, registre `could-not-run (sandbox)` e pare de
tentar — o chip re-roda as duas lanes fora do sandbox e essa re-corrida é a verificação de registro.
NUNCA reporte Pass em comando que você não rodou.

## Commit

Um commit, `fix(web):` na cabeça, no branch `chip/import-chain`. NUNCA `push`, NUNCA `merge`, NUNCA
`reset`/`revert`/`stash`/`clean`.

## G-questions para o seu relatório

- G1: o novo helper de data deveria chamar `formatDateTime` depois do guard de tipo, ou reimplementar
  o parse? Responda pelo código (quem mais consome `formatDateTime`, e de quem é aquele arquivo).
- G2: uma a três linhas sobre qualquer coisa que o card deixou em aberto e você decidiu.
- G3: alguma das três condições encosta numa costura de outro chip?
