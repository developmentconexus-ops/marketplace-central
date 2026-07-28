# CHIP-IMPORT-CHAIN — contrato de validação

PASS/FAIL/NOT-PROVEN por critério. Um critério que você não consegue fazer **totalmente** vira
REPORT e o chip fecha sem ele (R-24). Verifique por **string**, nunca por linha.

---

## I1 — a rota `/importacoes` existe e monta

`AppRouter.tsx` registra o path; a tela renderiza. Declare a decisão de gate (dentro ou fora de
`InstallationGatedRoutes`) **com a justificativa**, e prove o comportamento que você escolheu: se
ficou fora do gate, mostre a tela montando sem instalação de marketplace.

## I2 — o `ImportacaoSection` mudou de casa sem perder função

Componente e teste no diretório novo; `pages/vinculos/` não o contém mais.

O teste que existe hoje continua verde no lugar novo. Se você o reescreveu, diga o que mudou de
asserção — um teste "adaptado" que perdeu cobertura na mudança é regressão disfarçada de refactor.

Declare o que `/vinculos` mostra agora no lugar dele.

## I3 — `VinculosPage.tsx` foi tocado em duas linhas e nada mais

```bash
git diff 5441fe18f64171ef61cb03b51b5bf66e2922e4eb HEAD -- apps/web/src/pages/vinculos/VinculosPage.tsx
```

Esperado: só a remoção do `import` e do `<ImportacaoSection />`. Qualquer outra hunk é colisão com o
CHIP-VINC-NEUTRO, que está editando `/vinculos` em paralelo — e colisão é FAIL, não detalhe de merge.

## I4 — a tela consome `getErpImportChain`, não os números velhos

Prove o consumo do endpoint novo. Um contador que você calculou no cliente a partir de
`listErpImports` **não é** a cadeia — é a aparência dela, e é exatamente a gambiarra que o D-D
existe para não deixar acontecer.

## I5 — ADR-17 na cadeia: campo ausente é `—`, nunca zero

Teste com payload de cadeia com campo ausente/nulo. Esperado: `—` (ou o `UnknownValue` do design
system), **nunca** `0`.

Zero fabricado num contador de decomposição diz ao operador "não há nada vinculado" quando a verdade é
"não sabemos". Essa é a distinção que a missão inteira existe para preservar.

## I6 — erro do chain-read é honesto

Teste com a chamada falhando. A tela mostra estado de erro, não uma cadeia vazia que parece um import
sem produtos. Diga qual estado você renderiza.

## I7 — F-02 verificado, não reconstruído

O EVIDENCE declara:

- que `packages/web-query/src/activeSource.ts` e `IntegracoesPage.tsx:297-346` já implementam o F-02
  contra `GET/PUT /config/active-source` — com o commit que os trouxe;
- o resultado da sua verificação: o toggle persiste no DB e a leitura reflete o servidor (**ou** o
  que falhou, e aí vira escopo seu);
- que o `milestone.md` cita `/tenants/{tenant_id}/active-source`, que não é o endpoint landado.

FAIL se você reconstruiu o que já existia. FAIL se você declarou "já existe" sem ter verificado.

## I8 — F-03 resolvido de um jeito ou de outro

Ou **satisfeito** (diga qual campo a tela precisava e onde `listErpImports` já o entrega), ou
**estendido** — e aí OpenAPI + `packages/sdk-runtime` + consumo no **mesmo commit** (profile §7), com
a justificativa do campo. Nunca "consumi um campo que o contrato não declara".

## I9 — sem dano colateral

- `tsc -p apps/web --noEmit` da raiz do repo principal. Os 15 erros de `main` continuam **15 ou
  menos** e nenhum novo em arquivo seu. Liste-os; os 3 de `/vinculos` são do CHIP-VINC-NEUTRO e os 12
  são baseline.
- `vitest` com contagem citada antes e depois. Verde que virou vermelha é FAIL.
- `git diff --name-only 5441fe18f64171ef61cb03b51b5bf66e2922e4eb HEAD` — **zero** arquivos Go, **zero**
  migrations, e de `pages/vinculos/` apenas `ImportacaoSection.tsx` (+ teste, deletados) e
  `VinculosPage.tsx`.

`base_sha` é PISO, não ponto fixo. `main` se move enquanto você trabalha; rode o comando acima contra
o `main` de verdade antes do fecho e declare o que apareceu que não é seu.

---

## L2 — o drive ao vivo, e por que ele é deste chip

**Este é o L2 mais importante da onda 2, e ele é do HUB.** O `GET /erp/imports/{id}/chain` landou no
CHIP-ANCHORS-2 e **nunca foi dirigido ao vivo por ninguém**: o operador deu waiver em 2026-07-28,
com o risco nomeado antes do aceite, precisamente porque a prova viva acontece nesta tela. O risco
aceito foi: **um defeito de FIAÇÃO — composition root, decorator perdido — fica escondido até agora.**
É a classe exata do catalog-503 do M-02, que teste nenhum pegou.

Então L2 aqui não é formalidade. É a cobrança do waiver.

Roteiro, executado pelo hub:

1. `/importacoes` carrega com dado real de import.
2. A cadeia de um import de verdade renderiza, e os números batem com o DB — **conferidos contra
   consulta direta**, não contra a própria tela.
3. Um id que não existe → erro honesto na tela, não cadeia vazia.
4. `/integracoes`: virar o active-source persiste no DB e a app inteira passa a ler a outra fonte
   (é o que a invalidação global do `activeSource.ts` promete).
5. Claro e escuro.

Se a cadeia responder mas com número visivelmente baixo, suspeite do defeito de subcontagem que o
CHIP-ANCHORS-3 está consertando (`vinculados` perde CODPROD com zero à esquerda) antes de culpar o
FE — e registre a suspeita com o número observado.

Envie **REQUEST** ao hub para o L2. Não tente subir stack no worktree: o compose resolve
`env_file: .env` relativo ao diretório do compose, então `--env-file` apontando para outro lugar
**não funciona**, e dirigir daqui exigiria copiar credenciais para um segundo diretório — que é
chamada do operador, não sua.

---

## O que este contrato NÃO cobre

- **`/vinculos`** (badge, `INCOMPARABLE`, vocabulário) é do CHIP-VINC-NEUTRO.
- **Os 15 erros de `tsc`** não são deste chip; nenhum deles.
- **Qualquer correção de backend**, inclusive a subcontagem de `vinculados` — é do CHIP-ANCHORS-3.
- **`G4`** (índice) e **`B-08`** (deadline de rota) têm dono nomeado na reconciliação do gate.
