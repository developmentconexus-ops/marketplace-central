# CHIP-IMPORT-CHAIN — `/importacoes`: rota própria e a cadeia de import de verdade

```yaml
chip: CHIP-IMPORT-CHAIN
branch: chip/import-chain
base_sha: 5441fe18f64171ef61cb03b51b5bf66e2922e4eb
wave: 2 (paralelo a CHIP-VINC-NEUTRO e CHIP-ANCHORS-3 — write-sets disjuntos)
milestone: M-06 (F-01 + F-03, fatia por TELA conforme D-E)
design_ref: .mnfs/MIS-004-mvp-demo/DESIGN-REFERENCE.md
```

## Autoridade

- `.mnfs/MIS-006-integracao-fundacao/DECISOES-D122-anchors-telas.md` — **leia inteiro**. D-D (o
  chain-read é backend e fecha a lacuna de decomposição do M-06) e D-E (M-06 fatiado por TELA)
  governam este chip.
- `.mnfs/MIS-006-integracao-fundacao/M-06-telas-sdk/milestone.md` — F-01 e F-03. **F-02 está
  descarregado, veja abaixo.** F-04 é do CHIP-VINC-NEUTRO.

Aponte, não recopie (R-14).

## O que já existe — verificado no repo, não suposto

O `milestone.md` é de D-120 e o repo andou desde então. Confira você mesmo antes de construir; estes
são os fatos que o hub verificou em `5441fe18`:

| Coisa | Estado |
|---|---|
| `GET /erp/imports/{id}/chain` | **existe** — landou no CHIP-ANCHORS-2 (F-04) com OpenAPI e SDK no mesmo commit |
| `client.getErpImportChain(id)` | **existe** — `packages/sdk-runtime/src/index.ts:1901`, tipo `ErpImportChain` |
| consumo do chain no FE | **não existe** — `pages/vinculos/useErpImports.ts` só chama `listErpImports` (`:19`) e `getErpImport` (`:29`) |
| rota `/importacoes` | **não existe** — `app/AppRouter.tsx` não tem esse path |
| rota `/integracoes` | **existe** — `AppRouter.tsx:66` |
| `ImportacaoSection` | **existe**, 200 linhas, **dentro de** `pages/vinculos/`, renderizada em `VinculosPage.tsx:159` |

### F-02 (`ActiveSourceCard` do DB) está DESCARREGADO — não refaça

O `milestone.md` traz uma decisão de operador D-120 dizendo que o rádio de active-source está morto
desde o M-02 e que a fiação fecha aqui. **Já fechou**, antes deste chip:

- `packages/web-query/src/activeSource.ts` — `useActiveSourceQuery` / `useSetActiveSourceMutation`
  sobre `GET/PUT /config/active-source`, com invalidação global no sucesso.
- `apps/web/src/pages/integracoes/IntegracoesPage.tsx:297-346` — consome o hook; `active_source` sai
  do servidor; sem default no cliente; tenant sem linha falha fechado.

O comentário no topo do `activeSource.ts` registra o defeito exato que o brief descreve, já corrigido.

**Sua obrigação aqui é verificar e declarar, não construir.** Rode a tela, confirme que o toggle
persiste no DB e que a leitura reflete o servidor, e registre no EVIDENCE que o F-02 já estava feito e
por qual commit. Se a verificação **falhar**, aí sim vira escopo seu — e você diz o que falhou.

Também note a correção de drift já registrada no `milestone.md`: o brief cita
`GET/PUT /tenants/{tenant_id}/active-source`, mas o endpoint que landou é `/config/active-source`.
O landado é o que vale.

## F-01 — rota `/importacoes` com a cadeia de verdade

Três partes.

**1. Promover o `ImportacaoSection` para fora de `pages/vinculos/`.** Ele é de importação, não de
vínculos, e hoje mora na tela errada. Mova o componente e seu teste para o diretório da tela nova.
Em `VinculosPage.tsx` isso são exatamente **duas linhas**: o `import` em `:8` e o `<ImportacaoSection />`
em `:159`. **Nada mais nesse arquivo é seu** — o CHIP-VINC-NEUTRO está editando `/vinculos` em
paralelo, e seu diff em `VinculosPage.tsx` tem que caber nessas duas linhas.

Decida e declare o que acontece com o lugar antigo: `/vinculos` perde a seção (e então diga para onde
o operador vai), ou ela aparece nos dois lugares. Sumir em silêncio com uma seção que o operador usa
é regressão, não limpeza.

**2. Registrar a rota** em `AppRouter.tsx` + a entrada de navegação. Decida se `/importacoes` fica
dentro do `InstallationGatedRoutes` (`:75`) ou fora, como `/integracoes` (`:66`), e **justifique**:
importação de ERP não depende de instalação de marketplace, então o gate pode estar errado para esta
tela. Uma linha de justificativa, não uma escolha silenciosa.

**3. A cadeia de verdade.** Consuma `client.getErpImportChain(id)`. É por isso que este chip existe:
o endpoint landou no CHIP-ANCHORS-2 e **nunca foi dirigido ao vivo por ninguém** — o operador
concedeu o waiver justamente porque a prova viva acontece aqui, na tela que o consome.

O `ErpImportChain` carrega os números de decomposição (produtos do import, vinculados, enfileirados).
Renderize o que o payload traz, honesto por ADR-17: campo ausente é `—`, nunca zero fabricado. Zero
inventado num contador de decomposição é a classe de bug que a missão inteira existe para matar.

E lembre do que a reconciliação do gate registrou: `vinculados` **subconta** CODPROD com zero à
esquerda. O CHIP-ANCHORS-3 está consertando isso em paralelo, em `main`. Não conserte backend; se seu
drive ao vivo mostrar um número que parece baixo, é possivelmente esse defeito e o registro é REPORT.

## F-03 — SDK

`getErpImportChain` já existe. O que sobra do F-03 é o que `listErpImports` devolve hoje contra o que
a tela precisa. Verifique primeiro. Se faltar campo, é mudança de contrato: OpenAPI +
`packages/sdk-runtime` no **mesmo commit** que o consumo (profile §7), e diga no EVIDENCE por que o
campo é necessário para a tela.

Se não faltar nada, **declare F-03 já satisfeito** em vez de inventar trabalho.

## Matriz de propriedade

| Eixo | CHIP-IMPORT-CHAIN |
|---|---|
| Migração | **nenhuma** |
| Backend Go | **nenhum** |
| Contrato/SDK | só se o F-03 exigir campo novo; nesse caso **só o path `/erp/imports*`** do OpenAPI e o módulo de erp-import do SDK. Os paths `/product-links*` são de outro dono |
| FE — seu | diretório novo da tela `/importacoes`; `app/AppRouter.tsx`; a superfície de navegação; `pages/integracoes/` (verificação do F-02); `pages/vinculos/ImportacaoSection.tsx` + seu teste (**movidos para fora**); `pages/vinculos/VinculosPage.tsx` **restrito às 2 linhas do `ImportacaoSection`** |
| FE — **NÃO seu** | `pages/vinculos/QueueRow.tsx`, `VinculoDrawer.tsx`, `QueueTab.tsx`, `ResolvidosTab.tsx`, `useVinculosResolved.ts` — do CHIP-VINC-NEUTRO |
| `tsc` | os 15 erros de `main` **não são seus** (3 do CHIP-ANCHORS-2 em `/vinculos`, 12 baseline). Declare-os pré-existentes; não os conserte |

## Ladder

- L0: `tsc` no write-set + lint.
- L1: `vitest` — baseline citada antes e depois.
- L2: **do hub**, e este é o L2 mais importante da onda. Ver o contrato de validação.

O worktree do chip não tem `node_modules`. Rode **`npm ci` na raiz do worktree** antes de qualquer
lane de `tsc`/`vitest` — profile §3, `ratified`: install fiel ao lockfile é preparo de ambiente, não
mudança de dependência, e não pede REQUEST. **Nunca reuse o `node_modules` de outro checkout**: os
symlinks de workspace do npm resolveriam `packages/*` para a OUTRA árvore e a suíte validaria o
código errado em silêncio. Se o `npm ci` falhar, é REQUEST — não workaround.

`npx --no-install tsc` no worktree passa vazio — pass vacuoso.

E cuidado com o observável: contar 15 erros de `tsc` com a composição esperada NÃO prova que você leu
a árvore certa, porque os 3 erros de `/vinculos` também existem em `main` desde `dbdcdfb1`. Ler a
árvore errada dá a mesma contagem.

Este chip não sobe servidor, não liga em `:8080`, não lê `.env*`. O drive ao vivo é REQUEST ao hub.
