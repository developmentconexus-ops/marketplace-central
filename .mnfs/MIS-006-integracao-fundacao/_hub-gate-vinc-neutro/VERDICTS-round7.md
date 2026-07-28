# VEREDITOS — CHIP-VINC-NEUTRO, round 7 (gate diff-only)

`alvo` = `main` @ `94374920` · `chip` = `chip/vinc-neutro` @ `8271991b`
Entrada dos dois assentos: `code-diff.patch` (2.662 linhas, sha256 `511a9b04…e3e6bb3a`) + `GATE-P6.md`.
Assentos CEGOS entre si. Nenhum recebeu o pack do chip.

Resultado bruto: **Opus APROVA (9 REPORT, zero bloqueante) / Sol REFUTED (1 bloqueante)** — split.
O veredito é do hub, na seção `## VEREDITO DO HUB`.

---

## ASSENTO 1 — Opus, adversarial — TEXTO CRU **PERDIDO**, e isto é um achado

O retorno integral deste assento **não existe mais em custódia do hub**. Medido, não suposto:

```
tasks/<id>.output                    -> 0 bytes  (quatro arquivos de assento, todos 0)
transcript jsonl, varredura de bloco -> nenhum bloco >1.200 chars com o texto do assento
                                        na vizinhanca do turno que o recebeu (linhas 3400-3460)
```

O que sobrevive é **paráfrase do hub**, escrita no turno do retorno, e está marcada como tal.
Não é verbatim e **não vale como verbatim** para nenhum efeito:

> APROVA. 9 REPORTs, zero bloqueante — e todos sobre código.
> Entre eles, no mesmo lugar em que o Sol bloqueou (`QueueRow.tsx` `reasonSideLabel`, `patch:208`):
> o gate é o literal `reason.direction !== "INCOMPARABLE"`; o comentário em `patch:211-214`
> afirma isso totalmente ("Unknown falls through verbatim — never hidden"); e
> **"No observable today"** — o assento verificou que o backend só seta `Side` nos ramos
> INCOMPARABLE (`generation_service.go:636-640` e `:711,715,723,726,728`).

Por que se perdeu: o veredito do Sol coube inline na notificação (5.625 chars) e sobreviveu; o do
Opus era grande demais, foi para arquivo, e o arquivo esvaziou. A janela entre **retornar** e
**ser colado no pack** atravessou uma compactação. O §11 dizia "verbatim antes da análise" e eu li
isso como "antes de escrever o ruling" — tarde demais. Ratificado como emenda: **cola no retorno**.

Consequência desta rodada, declarada e não contornada: dos 9 REPORTs do Opus, só viram ordem os
que eu consiga re-derivar por string hoje. Um REPORT que existia apenas na paráfrase **não vira
ordem** — e eu digo isso em vez de reconstruir de memória o que o assento teria dito.

Registro de honestidade: numa primeira tentativa colei aqui, como se fosse este assento, o veredito
do gate do **CHIP-ANCHORS-2** (B-01 `refforn`, B-04 `codprod`, B-08 `RouteClassMux`) — outro chip,
outra rodada. Peguei por `grep` de tamanho e não conferi o cabeçalho. Corrigido antes de qualquer
ruling; fica escrito porque é exatamente a classe de erro que este gate bloqueia nos assentos.

---

## ASSENTO 2 — Sol medium — VERBATIM (task `a1d714197daabb7d7`, recuperado do transcript)

<task-notification>
<task-id>a1d714197daabb7d7</task-id>
<tool-use-id>toolu_01XFJ6Y4TkiFLgEMsj91Mztu</tool-use-id>
<output-file>C:\Users\LEANDR~1.THE\AppData\Local\Temp\claude\C--Users-leandro-theodoro-Documents-marketplace-central--claude-worktrees-epic-lehmann-4ffbad\5690900d-57ed-4100-aa83-26e768290dee\tasks\a1d714197daabb7d7.output</output-file>
<status>completed</status>
<summary>Agent "Sol medium seat, vinc-neutro diff" finished</summary>
<note>A task-notification fires each time this agent stops with no live background children of its own. The user can send it another message and resume it, so the same task-id may notify more than once.</note>
<result>GPT-5.6 Sol, effort medium

VERDICT: REFUTED

| Severity | Criterion | File:line | Defect | Concrete failure scenario |
|---|---|---|---|---|
| BLOCKING | NEW | `apps/web/src/pages/vinculos/QueueRow.tsx:229` | `reasonSideLabel` enumerates `INCOMPARABLE` by literal and discards `side` for every future direction, despite the drift policy requiring unknown wire data to degrade verbatim. | API adds `{direction:"PARTIAL", side:"provider"}` before SDK regeneration; the user sees `PARTIAL EAN` in both row and drawer instead of the backend-sent side. **Discriminator:** the operator observes less reason data than the wire supplied and loses where the datum applies. |
| REPORT | NEW | `apps/web/src/pages/vinculos/QueueTab.test.tsx:92` | The `exact_ean`/`CONFIRM`/60 fixture contains `title AGAINST`, an impossible tuple because the sole producing path for that reason applies the hard-negative override to `REJECT`/25. | The test remains green using a row that `generation_service.go:560-568` cannot emit, so it cannot fail against the real rejected-row state. |
| REPORT | NEW | `apps/web/src/pages/vinculos/QueueTab.test.tsx:111` | The `exact_sku` fixture encodes a missing EAN as `UNAVAILABLE`, although an empty listing EAN is reclassified as `INCOMPARABLE` with side `provider`. | With `detail: "sem EAN para corroborar..."`, `generation_service.go:702-728` produces `INCOMPARABLE/provider`; the test instead validates a non-production chip. |
| REPORT | NEW | `apps/web/src/pages/vinculos/wireFixtures.ts:293` | The total guarantee "backend can actually emit" is partial because reason/state consistency is explicitly not validated, as the two accepted impossible fixtures demonstrate. | A test author calls `wireCandidate` with either impossible reason set above and receives a candidate instead of the promised exception. |
| REPORT | NEW | `apps/web/src/pages/vinculos/wireFixtures.ts:108` | `DECLARED_PROVIDER_CAPABILITIES` claims to be the complete adapter set, but its guard inspects only the hard-coded Mercado Livre declaration and cannot detect a newly added provider adapter. | After a Shopee capability adapter lands, the guard can remain green while `wireCandidate({provider_code:"shopee"})` still rejects a now-producible row. |

Criteria:

- DISCHARGED: V-1's narrow ranking invariant—the all-`INCOMPARABLE` row renders motivos, every reason participates in the sort, future directions are retained, and glyph/class lookups degrade safely. The blocking finding concerns associated `side` data, not loss of the reason itself.
- DISCHARGED: V-2. `refforn` is deliberately retained for persisted historical reasons.
- DISCHARGED: V-3. Structural labels are neutral, `mercado_livre` renders as `Mercado Livre`, and unknown providers degrade without blanking.
- DISCHARGED: V-4. The mapping mirrors `decisionRuleForCandidate`; title-only `REVIEW` does not become an identifying anchor.
- DISCHARGED: V-5. The badge uses the newest resolving audit entry's exact `actor.actor_type === "system"` predicate; absent/manual entries receive no badge.
- DISCHARGED: V-6. All changed paths are within the granted frontend directory; forbidden production seams are untouched and the golden assertions were not loosened.
- Could not judge from diff alone: none after the bounded producer/type reads.

Paths directly read:

- `.mnfs\MIS-006-integracao-fundacao\_hub-gate-vinc-neutro\GATE-P6.md`
- `.mnfs\MIS-006-integracao-fundacao\_hub-gate-vinc-neutro\code-diff.patch`
- `packages\sdk-runtime\src\index.ts`
- `apps\server_core\internal\modules\product_links\application\resolution_service.go`
- `apps\server_core\internal\modules\product_links\adapters\postgres\link_candidate_repo.go`
- `apps\server_core\internal\modules\product_links\application\generation_service.go`

A bounded symbol scan also returned individual matching lines from:

- `packages\sdk-runtime\src\index.test.ts`
- `apps\server_core\tests\integration\product_link_decisions_test.go`
- `apps\server_core\internal\modules\product_links\transport\http_handler_test.go`
- `apps\server_core\internal\modules\product_links\transport\http_handler.go`
- `apps\server_core\internal\modules\product_links\ports\workflow_store.go`
- `apps\server_core\internal\modules\product_links\domain\product_link.go`
- `apps\server_core\internal\modules\orders\adapters\productlinks\link_reader_test.go`
- `apps\server_core\internal\modules\product_links\application\resolution_service_test.go`
- `apps\server_core\internal\modules\product_links\application\auto_link_policy_test.go`
- `apps\server_core\internal\modules\mutations\ports\linkage.go`
- `apps\server_core\internal\modules\mutations\application\writer_test.go`
- `apps\server_core\internal\modules\mutations\adapters\productlinks\writer.go`
- `apps\web\src\pages\mutations\MutationIntentForm.tsx`</result>
<usage><subagent_tokens>57661</subagent_tokens><tool_uses>23</tool_uses><duration_ms>760542</duration_ms></usage>
</task-notification>

---

## VEREDITO DO HUB — round 7

Split: Opus APROVA, Sol REFUTED. Reclassificar é do hub, e só por string.

### 1. O bloqueante do Sol (`QueueRow.tsx:229`) vira **REPORT**

O achado: `reasonSideLabel` abre com `if (reason.direction !== "INCOMPARABLE") return undefined;`
e portanto descarta `side` de qualquer direção futura.

**Verificado por string, na `main` de hoje** — todo `return` do classificador que carrega `Side`
não-vazio é `Incomparable`, sem exceção:

```
grep -n 'return domain.LinkCandidateReasonDirection' .../generation_service.go
704: ...DirectionUnavailable,  ""                  <- sem side
711: ...DirectionIncomparable, ""                  <- sem side
715: ...DirectionIncomparable, ...SideProvider
723: ...DirectionIncomparable, ...SideBoth
726: ...DirectionIncomparable, ...SideProvider
728: ...DirectionIncomparable, ...SideERP
```

e `generation_service.go:684` é o **único** sítio do gerador que atribui `Side:`, alimentado por
esse classificador. Hoje `Side != ""` implica `direction == INCOMPARABLE`.

Discriminador: sobe o código com o achado intacto — nenhum usuário, operador, chamador ou linha
gravada faz nada diferente, porque o wire não produz a entrada que dispara o defeito. **REPORT.**

**O que os DOIS assentos leram pela metade.** Ambos citaram o comentário do `?? reason.side`
("Unknown falls through verbatim — never hidden") como frase total sob guarda parcial — a forma
R-24. Não é. Essa frase governa o **VALOR** de `side` (o fallback quando o SDK não foi
regenerado), e o código faz exatamente o que ela diz. A política de **DIREÇÃO** está declarada no
bloco de comentário imediatamente acima da função, e é verdadeira:

> The direction is checked FIRST: the frozen D-122 contract emits `side` only for INCOMPARABLE,
> so a stray value on any other direction is wire noise, not a fact to render.

Alegação com escopo nomeado, e o escopo confere em código. É a primeira vez nesta missão que uma
frase de aparência total passa no teste do R-24 em vez de falhar nele. Fica registrado como o
contra-exemplo — R-24 mata universal falso, não mata alegação escopada.

### 2. O resíduo real existe, e o dono é o HUB — não o chip

A invariante *"`Side` só em INCOMPARABLE"* **não é guardada em lugar nenhum do Go**. Eu a verifiquei
lendo os seis `return` de UMA função; nada impede o sétimo. O literal do front é consumidor
correto de uma invariante de produtor não-imposta — e é no produtor que ela se impõe barato.

Backend Go está FORA do write-set do VINC-NEUTRO (V-6, verbatim: "Backend Go: nenhum"). Logo isto
não é ordem para este chip. Vai para o território do ANCHORS-3 (mesmo arquivo, mesma função).

### 3. Duas fixtures improdutíveis — **ORDENADAS**, verificadas por string

Ambas do Sol, ambas confirmadas por mim antes de virar ordem:

**F-A — `QueueTab.test.tsx:92`**: `{ anchor: "title", direction: "AGAINST" }` sob
`confidence: 60, confidence_band: "MEDIA", match_status: "CONFIRM"`. Em
`generation_service.go:560-562` o `title/AGAINST` e o score saem da **mesma atribuição**:

```go
if hardNeg, detail := detectHardNegative(snapshot.Title, product.Name); hardNeg {
    confidence, band, status = 25, ...BandBaixa, ...MatchStatusReject
    reasons = append(reasons, ...{Anchor: "title", Direction: ...DirectionAgainst, ...})
}
```

Existir o motivo obriga `25/BAIXA/REJECT`. A tupla da fixture é improdutível.

**F-B — `QueueTab.test.tsx:111`**: `{ anchor: "ean", direction: "UNAVAILABLE", detail: "sem EAN" }`.
`UNAVAILABLE` é só o ramo de **capability** (`:704`, "provider não fornece a âncora"). Anúncio sem
o VALOR é `INCOMPARABLE` + `side: "provider"` (`:715`, "anúncio sem %s"). A fixture funde
ausência-de-capacidade com ausência-de-valor — a distinção exata que o ADR-17 e o V-1 existem para
manter. Um teste que fixa o mundo errado é pior que teste nenhum; já foi o meu ruling no CORR-6 do
ANCHORS-3 e vale aqui igual.

Classe: REPORT pela régua (nenhum observável muda). **Ordenadas assim mesmo** — REPORT não segura
merge, não significa que não se conserta.

### 4. Mecanismo, não remendo nas duas fixtures

`PRODUCIBLE_TUPLES` valida a tupla de score por sítio produtor. Não restringe `reasons[]` contra
essa tupla — por isso F-A e F-B são **legais** sob o mecanismo que o round 6 entregou. O defeito
não está nas duas fixtures; está no escopo do guard.

**Ordem:** estender o mecanismo aos motivos — por sítio produtor, o conjunto de `reasons` que ele
pode emitir, com direção e side. Consertar só as duas fixtures deixa a terceira livre.

Não invoco a regra da terceira rodada por contagem: os rounds 1-6 correram na forma antiga, os
vereditos estão em custódia do CHIP e não do hub, e eu não consigo verificar a contagem por string.
A ordem não depende dela.

### 5. Varredura parcial — o mesmo mock morto sobreviveu no gate de design

`2b956e19` apagou o mock `listErpImports` de `VinculosPage.test.tsx` (restou só o comentário
`:34`). Na mesma tip `8271991b`, `VinculosDesign.golden.test.tsx` ainda carrega os quatro:

```
25:  const listErpImports = vi.fn();
31:    listErpImports: (...a: unknown[]) => listErpImports(...a),
120:    listErpImports.mockReset();
121:    listErpImports.mockResolvedValue({ items: [] });
```

Apagar mock morto não é afrouxar asserção — V-6 não protege isto. Ordenado.

### 6. Ordem de merge, travada

**ANCHORS-3 entra antes.** O rótulo `["CODPROD"]` do front só é verdadeiro depois que o CORR-1
troca `ReferenceCode` por CODPROD canônico no backend. Merge invertido publica um rótulo falso.

### 7. Emenda ratificada por esta rodada

O texto cru do assento Opus se perdeu (seção 1 acima). O §11 pedia verbatim "antes da análise"; a
janela que matou o registro foi entre **retornar** e **ser colado**, atravessando uma compactação.
Emenda: veredito de assento é colado no pack **no turno em que retorna**, antes do hub o ler
inteiro. Um veredito que só existe no contexto do hub não existe.

**Veredito: sem bloqueante sobrevivente. Não merge ainda** — trava de ordem (item 6) e as três
ordens dos itens 3/4/5 entram antes, sem abrir rodada nova de gate.

---

## RESÍDUO PÓS-RULING — o patch de uma rodada FECHADA não se re-corta

A rodada foi julgada com o chip em `8271991b`. Depois disso o chip andou para `fa040538`.
O chip perguntou explicitamente: re-cortar o `code-diff.patch` ou registrar o resíduo — e disse
que não trataria silêncio como nenhuma das duas. Certo, e a resposta é a segunda.

**Decisão: não re-corto.** O `code-diff.patch` de `293c1485` é artefato de uma rodada **decidida**.
Defasagem de patch de rodada fechada não é defeito, é história — o que faltava era alguém dizer
por escrito de quanto ela é. Aplicado o discriminador de `979be178` (toca código de produção, ou
algo em que um veredito devolvido se apoia?):

| delta pós-corte | arquivo | linhas de produção | veredito que se apoia nele |
|---|---|---|---|
| `2b956e19` | `VinculosPage.test.tsx` | 0 | nenhum |
| `1fcf7f1a` | `VinculosDesign.golden.test.tsx` | 0 | nenhum |

Os dois são deleção do mesmo mock morto (`listErpImports`), morto pela mesma causa (`4b76a287`
tirando `ImportacaoSection` de `pages/vinculos/`). Nenhum assento apoiou veredito em qualquer dos
dois arquivos: os REPORTs do Sol citam `QueueTab.test.tsx` e `wireFixtures.ts`; a varredura parcial
do golden test foi achado **meu**, por string, não de assento.

Ambos ficam registrados nestas palavras: **verificado-pelo-HUB, não verificado-por-assento.**

Re-medido pelo hub em `fa040538`, contra a `main` de agora (`0bda36bb`):

```
git diff --numstat 0bda36bb fa040538 -- apps contracts packages
  11 arquivos, TODOS em apps/web/src/pages/vinculos/
  linhas com 0 inserções (revert)                    -> nenhuma
git show fa040538:.../VinculosDesign.golden.test.tsx | grep -n listErpImports
  -> só o comentário :26; zero mock
  beforeEach preservado (:124) com os dois resets vivos
```

**O-3 estava descarregado antes da ordem chegar** (`1fcf7f1a`, mensagens cruzadas). Terceira
convergência cega desta onda. Restam **O-1** (duas fixtures improdutíveis) e **O-2** (estender o
mecanismo aos motivos).

### O que o chip achou sobre a própria varredura, e vale mais que o mock

O chip nomeou por que a primeira varredura não pegou o segundo mock: o padrão foi
**"o arquivo que eu tinha aberto"**. A população virou a pegada da EDIÇÃO em vez da pegada do
FATO — e o fato era um merge que moveu um componente para fora do diretório, cuja população
correta é *todo arquivo que mocava a porta*, na árvore inteira. Instância do §11
("A sweep is only as wide as its pattern"), com o gatilho específico nomeado. Levado para lá.
