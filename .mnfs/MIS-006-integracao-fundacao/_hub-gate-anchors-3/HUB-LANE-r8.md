# ADENDO ao GATE-P6-r8 — a lane é do HUB, e o tip mudou sem mudar o patch

`tip agora: 1542839a` · `patch: r8-code-diff.patch, INALTERADO` ·
`sha256 362fa10c…ef957fd1`, 1227 linhas, **byte-idêntico** ao corte original

## Por que o tip mudou e por que a rodada fica de pé

O gate foi cortado em `c3acf62b`. O chip commitou `1542839a` depois, a pedido do hub, para
versionar a saída CRUA do par must-fail. Medido:

```
git diff --stat c3acf62b 1542839a -- ':!.mnfs'   ->  VAZIO
```

Delta é só pack. Não toca código de produção nem nada em que um veredito devolvido se apoie, e
nenhum veredito havia sido devolvido — os assentos não tinham sido despachados. Rodada de pé,
delta registrado como HUB-verificado.

O patch foi RE-CORTADO contra o tip atual da main e comparado byte a byte com o corte original:
idêntico, mesmo sha256. A main avançou só em `docs/HARNESS-PROFILE.md`, que é arquivo do hub e
está fora do patch.

## As lanes que viajam com o patch são do HUB, não do chip

O chip entregou os runs crus, que é o que foi pedido. O hub **refez todos** no
`1542839a`, em worktree própria, com o EXIT do PROCESSO (não o do `sed` do pipe):

```
vitest importacoes+integracoes    5 files, 29 tests, 0 failed     EXIT(vitest)=0
web tsc -p tsconfig --noEmit      15 errors, 0 em importacoes, 0 em integracoes   EXIT(tsc)=2
sdk-runtime tsc --noEmit          0 errors                        EXIT(tsc)=0
sdk-runtime vitest                5 files, 77 tests, 0 failed     EXIT(vitest)=0
```

Todo número que o chip alegou reproduz na casa.

## F8-6 — atribuição do must-fail, refeita pelo hub

Mutação isolada em `apps/web/src/pages/importacoes/ImportacaoSection.tsx:151`,
`Ver estado` → `Ver cadeia`, nada mais tocado (diff da mutação em `lane-r8/hub-02-mustfail.log`):

```
× ImportacaoSection > renders import rows with protocol, status, and counts
  -> Unable to find an accessible element with the role "link" and name "Ver estado"

Test Files  1 failed | 4 passed (5)
     Tests  1 failed | 28 passed (29)
```

**Exatamente 1 vermelho**, e ele falha pelo NOME ACESSÍVEL — role `link`, name `Ver estado` —
não por `data-testid`. É isso que o F8-6 pede: o arm prova que o rótulo que o operador lê é o
que está afixado, não que um atributo de teste existe.

Restauração conferida por md5 nos dois lados (`d240bc9a…`), `git diff --quiet HEAD` = 0, e o
braço VERDE refeito: 5 files / 29 tests, EXIT 0.

## O que o hub NÃO mediu

Nada da lane de integração Go (`chain_query_repository_integration_test.go`) foi executado aqui.
O que essa lane precisa está declarado no GATE-P6-r8 e continua critério do assento sobre o
DESENHO do teste, não sobre a execução dele.

## Arquivos

`lane-r8/hub-01-vitest-importacoes-integracoes.log`, `hub-02-mustfail.log` (mutação + vermelho +
restauração + verde), `hub-03-web-tsc.log`, `hub-04-sdk-tsc.log`, `hub-05-sdk-vitest.log`.
