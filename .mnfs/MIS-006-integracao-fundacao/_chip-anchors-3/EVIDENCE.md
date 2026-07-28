# CHIP-ANCHORS-3 — EVIDENCE

```yaml
chip: CHIP-ANCHORS-3
branch: chip/anchors-3
base_sha: 5441fe18f64171ef61cb03b51b5bf66e2922e4eb
head: 590efdc8
gate_round: 3 fechado — REFUTED dos DOIS lados (A 2 blocking, B 5 blocking)
gate_round_3_reparo: 2 blocking + 2 nao-bloqueantes de autoria do chip corrigidos; 5 universais pre-existentes viram REPORT
gate_round_3_artefatos: os DOIS persistidos pelo orquestrador (A nao tem Write; B teve o apply_patch recusado pelo sandbox)
gate_round_2: SPLIT — lado A REFUTED (4 blocking), lado B CONFIRMED (Findings None)
gate_round_1: REFUTED dos DOIS lados
sweep_do_autor: rodada sobre o EVIDENCE inteiro (745 linhas, 71 hits); 4 alegacoes minhas corrigidas, 2 delas eram o universal do codigo escrito em prosa
r4_execucao: must-fail da segunda metade REPRODUZIDO pelo assento executor do hub em 0264ba84, contra Postgres real (failure_token=test=... blocked → status=passed). Unico criterio que dependia so da minha palavra.
finding_1: ESTREITADO pelo hub — falhou-vs-verde E distinguivel via failure_token; so pulado-vs-verde e byte-identico. Ratificado em ea919c06.
status: NAO FECHADO. Sem AGREEMENT. Todo achado dos dois lados verificado por STRING pelo chip, um recusado com motivo.
authority: .mnfs/MIS-006-integracao-fundacao/_hub-gate-anchors-2/p6-reconciliation-r1.md
contract: .mnfs/MIS-006-integracao-fundacao/_chip-anchors-3/validation-contract.md
```

**Este pack NÃO fecha o chip.** Três rounds do gate P6: round 1 REFUTED dos dois lados, round 2
SPLIT, round 3 REFUTED dos dois lados. Nunca houve `AGREEMENT`, e a linha `P6-DUAL-GATE:` é do hub.
Cada achado dos três rounds foi conferido por STRING por mim antes de virar conserto ou REPORT — um
foi recusado com o motivo escrito. As seções "Gate P6" trazem os verdicts e a conta do que sobrou.

Os achados do gate do CHIP-ANCHORS-2 **não** estão recopiados aqui. Eles vivem em
`_hub-gate-anchors-2/p6-reconciliation-r1.md` e são citados por nome — B-01, G1/B-04, G2/B-05,
B-03, B-07, B-09 (R-14/A13). Uma segunda cópia divergiria.

## Commits

| SHA | Corretivo | Achado que discharge |
|---|---|---|
| `9555ad30` | CORR-1 + CORR-4 + CORR-6b | B-01, B-03, B-09 |
| `055d1705` | CORR-3 (comportamento + OpenAPI + SDK no mesmo commit) | G2/B-05 |
| `bba08b41` | CORR-2 + CORR-6a | G1/B-04, B-07 |
| `54342331` | R2 — apaga o universal falso que **este chip** escreveu | gate round 1, blocking 1 e 2 |
| `2bed7d9d` | R4 — uma só regra de identidade para os três contadores + o que o guard nil PRODUZ | gate round 1, não-bloqueantes 6 e 4 |
| `0264ba84` | R7 — varre a **classe**, não o sítio que o gate nomeou: segundo universal falso + o teste passa a asserir 95/ALTA/ACCEPT | gate round 2, blocking 1 e 2 |
| `590efdc8` | Estreita os três comentários que a varredura refutou + **entra com o pack no git** (ver FINDING 8) | gate round 3, blocking A1/A2/B1 e não-bloqueantes A1/A2 |

CORR-5 é do hub, e não está aqui.

---

## Ledger de dispatch

Obrigatório por HARNESS-CORE §4 e pelo item (d) do complemento do hub. Papel, modelo, esforço,
caminho do artefato persistido.

| # | Papel | Modelo / esforço | Fatia | Artefato persistido |
|---|---|---|---|---|
| 1 | implementer | **`inline (DESVIO §4.2)`** — sessão orquestradora, Opus | CORR-1, CORR-2, CORR-3, CORR-4, CORR-6a, CORR-6b: as dez edições originais | — (sem artefato de worker; é o desvio, ver abaixo) |
| 2 | implementer | sonnet, dispatched subagent | A8 — teste de handler `{id}` malformado nas duas rotas | `dispatches/impl-a8-handler-malformed-id.md` (8.9 KB) |
| 3 | implementer | sonnet, dispatched subagent | A5/A6/A7/A10 — testes de integração contra Postgres real | `dispatches/impl-integration-a5-a7-a10.md` (323 linhas) |
| 4 | reviewer adversarial (§4.3) | sonnet, dispatched subagent, read-only | CORR-1/3/4/6b — implementer ≠ reviewer | `dispatches/review-adversarial-r1.md` — **CONFIRMED**, zero blocking. Ver a ressalva abaixo. |
| 5 | gate P6 lado A, round 1 | Opus, `harness:gate-reviewer`, fisicamente read-only | diff congelado inteiro | `dispatches/p6-opus-gate-r1.md` — **VERDICT: REFUTED**. **ARTEFATO SALVO DE TRANSCRIPT**, não escrito na hora — ver gate round 2, blocking 4 |
| 6 | gate P6 lado B, round 1 | `gpt-5.6-sol` / medium, sandbox read-only | diff congelado inteiro | `dispatches/p6-sol-gate-r1.md` (7648 B) — **VERDICT: REFUTED** |
| 7 | implementer | sonnet, dispatched subagent | R2 — narrow do comentário falso + cobertura de `seller_sku`/`side=erp` | `dispatches/impl-r2-false-universal.md` (19.767 B) |
| 8 | implementer | sonnet, dispatched subagent | R4 — canonicalização do `queued_products` + o que o guard nil produz | `dispatches/impl-r4-queued-canonicalization.md` (24 KB) |
| 9 | gate P6 lado A, round 2 | Opus, `harness:gate-reviewer`, read-only | **delta** `bba08b41..HEAD`, alvos (a)/(b)/(c) | `dispatches/p6-opus-gate-r2.md` |
| 10 | gate P6 lado B, round 2 | `gpt-5.6-sol` / medium, read-only | mesmo delta, sem critério de execução | `dispatches/p6-sol-gate-r2.md` |
| 11 | implementer | sonnet, dispatched subagent | R7 — segundo universal falso + fortalecer o teste para asserir 95/ALTA/ACCEPT | `dispatches/impl-r7-second-false-universal.md` (12606 B) |
| 12 | gate P6 lado A, round 3 | Opus, `harness:gate-reviewer`, read-only | delta `2bed7d9d..0264ba84` + **SWEEP obrigatória** dos dois arquivos inteiros | `dispatches/p6-opus-gate-r3.md` — **VERDICT: REFUTED**, 2 blocking. **Persistido por mim**: assento não tem Write. Cabeçalho de proveniência no arquivo |
| 13 | gate P6 lado B, round 3 | `gpt-5.6-sol` / medium, sandbox read-only | mesma delta, mesma SWEEP | `dispatches/p6-sol-gate-r3.md` — **VERDICT: REFUTED**, 5 blocking. **Persistido por mim**: o assento tentou escrever e o sandbox recusou o `apply_patch`. O corpo é o payload literal da escrita recusada (208 linhas, 16438 chars) |
| — | **assento de EXECUÇÃO** | **do HUB**, não deste chip | ladder, must-fail, fatos de git, sha256 | do hub, em worktree limpo próprio (`hub-exec-anchors3`). `go build ./...` e `go vet ./...` EXIT 0 em `2bed7d9d`; sha256 dos congelados r3 conferidos contra os meus. Registrado aqui para que a cobertura seja legível inteira. |

Brief do round 2: `dispatches/p6-gate-brief-r2.md`. Congelados do round 2: `p6-input-r3.patch`
(sha256 `b419ce47…`, 11 arquivos, +569/-43) e `p6-delta-r1-to-r3.patch` (sha256 `a1f58ee4…`,
4 arquivos, +135/-7). **Ambos obsoletos desde `0264ba84`** — os vigentes são `p6-input-r4.patch`
(`085882a6…`) e `p6-delta-r3-to-r4.patch` (`4943e3c3…`).

Brief idêntico entregue aos dois lados do gate: `dispatches/p6-gate-brief.md`.

**Ressalva sobre a linha 4.** O reviewer de fatia devolveu CONFIRMED, e o gate mostrou que ele
errou. Ele repetiu o mesmo hand-trace errado que eu sobre quais subtestes discriminam
(`review-adversarial-r1.md:19`, "verified by hand-trace"), e não pegou nenhum dos três blocking.
Um CONFIRMED de reviewer que erra o trace não é corroboração; conta menos, não mais. Registro isso
em vez de continuar somando o CONFIRMED dele como se fosse peso.

### O desvio, declarado e não lavado

A linha 1 do ledger diz `inline (DESVIO §4.2)` porque é o que aconteceu: as dez edições originais
foram escritas pela sessão orquestradora, não por worker despachado, antes do complemento do hub
chegar. Foi divulgado ao hub **antes** de continuar.

Ruling R1 do hub: aceito com condições, não refaça por worker. Condições, todas cumpridas:
a string `inline (DESVIO §4.2)` fica no ledger; esta sessão está desqualificada como reviewer
desses hunks e não assinou nenhum lado do gate; os prompts dos dois lados NOMEARAM os hunks
inline como alvo prioritário (`p6-gate-brief.md`, seção "Priority target"); o resto por worker.

**A condição 2 provou seu valor de forma cara.** Os três achados blocking do round 1 estão, os
três, em texto escrito inline. O review que foi pulado teria pegado. O gate pegou no lugar dele.

---

## Critérios

### A1 — o lado ERP do `seller_sku` é o CODPROD canônico — **PASS** (os dois lados do gate concordam)

As duas metades, verificadas por string:

- `grep "ReferenceCode"` **dentro** de `identityAnchorValues` = **0** ocorrências. As três
  ocorrências restantes em `generation_service.go` (linhas 403, 404, 416) estão em `newCandidate`,
  projetando o campo para dentro do registro do candidato — não são comparação cross-side.
- o caminho novo lê o campo canônico: `canonicalProductID(*product)` → `strconv.Itoa(productID)`.

O lado A foi além e confirmou a canonicidade rio acima: `internalread/reader.go:162`
`strconv.ParseInt(strings.TrimSpace(row.CodigoProduto)…)` → `:418` `InternalProductID: &canonical`.

### A2 — o teste de fixação é alcançável em produção — **FAIL no round 1, corrigido em `54342331`**

**O que aconteceu com `TestNamedMissingAnchorSitesAreIncomparableWithCorrectSide`, explicitamente,
como o contrato exige: foi CORRIGIDO, um caso foi DELETADO, e a justificativa da deleção estava
ERRADA. A correção da correção é o commit `54342331`.**

Mecânica, que passou: todos os fixtures de produto passaram a carregar
`InternalProductID: canonicalIDPtr(9xx)` — o fixture antigo era `&ProductCandidate{}`, sem id
canônico, forma que `findProducts` descarta antes de qualquer candidato existir. Teste novo no
nível do gerador: `TestSellerSKUAnchorReadsCanonicalCodprodNotSupplierReference` chama
`generateSingle`, que passa por `GenerateLinkCandidates` num service real.

**O que reprovou.** O caso `"exact EAN ERP seller SKU empty"` foi deletado sob a alegação de que
`side=erp` teria ficado INALCANÇÁVEL para `seller_sku`. **A alegação é falsa**, e o comentário que
a carregava foi commitado dentro do arquivo que o critério existe para proteger:

> ~~The ERP side of seller_sku is the CODPROD, which a resolved product always has — so side=erp is
> reachable for seller_sku only through the listing, never through the product.~~

As duas metades são falsas, verificadas no código por mim antes de qualquer conserto:

- "only through the listing" — invertido. `missingMatchedAnchorReason:640` `case listingValue == "":`
  põe em `SideProvider`, nunca `SideERP`.
- "never through the product" — `:643` `case product == nil || productValue == "":` põe em
  `SideERP`, e `product == nil` **é caminho de produção**: `generation_service.go:216`
  `applyUnresolvedScore(&unresolved, newProviderIdentityAnchorComparison(snapshot, identityAnchors, nil))`.
  Um anúncio que **traga** um `seller_sku` não vazio e cujo seller_sku, EAN e título não casem nada
  emite `seller_sku`/`INCOMPARABLE`/`side=erp` hoje, em produção. Com `seller_sku` vazio o ramo é
  outro — `case listingValue == "":` (`:645-647`) → `SideProvider` — e a frase não o cobre.
  (Estreitado na reconciliação da rodada 4: dizia "Qualquer anúncio", com maiúscula, e por isso
  escapou do padrão case-sensitive da varredura da rodada 3.)

Efeito líquido antes do conserto: **um caso de teste foi apagado justificado por uma frase que o
código refuta.** **Um chip cujo CORR-4 existe para apagar um universal falso tinha escrito um
universal falso novo.** É a forma R-24 que o próprio contrato nomeia: total na redação, parcial no
código.

Correção da rodada 3 a esta própria linha: eu escrevia aqui que "a única asserção de
`seller_sku`/`side=erp` da suíte tinha sido removida". **Isso era falso, e era o MESMO universal
sob outra roupa** — medir a perda como total sem varrer o arquivo.
`TestCase8NoAnchorResolvedYieldsZeroConfidenceNoCandidateWithReasons`
(`generation_service_test.go:1823-1826`) assere `seller_sku` `INCOMPARABLE` com `Side=erp`, é
anterior a este chip (`grep "MLB-FX8" p6-input-r4.patch` = 0 ocorrências, nem como contexto) e
nunca foi tocado. O que foi apagado foi um caso de tabela, não a cobertura da direção.

Conserto em `54342331`, por worker despachado (ledger linha 7), sob R-25 — a metade falsa foi
**deletada**, não anotada:

- comentário estreitado para o que se sustenta: com produto presente, `side=erp` não surge para
  `seller_sku`; o caso nil-produto é apontado para onde agora está fixado;
- teste separado, não caso de tabela: `TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide`
  dirige `GenerateLinkCandidates` ponta a ponta, então o produto nil vem do código de produção e não
  do fixture. O runner da tabela chama o scorer interno e fixa `declarations` nil, então não provaria
  que o `side=erp` semeado SOBREVIVE ao passo de âncora declarada. **Não chame isso de "cobertura
  restaurada"** — a direção já era coberta pelo Case 8 pelo mesmo caminho de produção. O que este
  teste acrescenta sobre o Case 8 é o `Detail` `"seller_sku sem correspondência"`, que o Case 8
  cita só na mensagem do `t.Fatalf` e nunca assere; o resto é corroboração. Rodada 3 corrigiu o
  comentário do teste, que carregava a mesma alegação de exclusividade;
- must-fail: `SideERP` mutado para `SideBoth` produz
  `Side:"both" … want direction="INCOMPARABLE" side="erp"`; restaurado por edição para a frente, com
  `git diff HEAD -- generation_service.go` vazio.

### A3 — must-fail do A1 — **PASS**

Revertido **só** o `productValue` do case `"seller_sku"` para `product.ReferenceCode`. Bruto em
`dispatches/a3-mustfail-raw.txt`. Assinaturas de falha, com o valor que apareceu no lugar de qual:

```
seller_sku reason=domain.LinkCandidateReason{Anchor:"seller_sku", Direction:"INCOMPARABLE", Side:"erp",
  Detail:"sem CODPROD para corroborar o EAN: o seller_sku do anúncio não casa nenhum produto"},
  want direction="UNAVAILABLE" side=""
```

```
seller_sku reason=domain.LinkCandidateReason{Anchor:"seller_sku", Direction:"INCOMPARABLE", Side:"both",
  Detail:"sem CODPROD para corroborar o EAN"}, want direction="INCOMPARABLE" side="provider"
```

A segunda é a mais dura: com `refforn`, o classificador afirma `side=both` — que o **anúncio E o
produto ERP** não têm `seller_sku` — sobre um anúncio que tem `"SKU-D1"` e um produto com CODPROD.

Restaurado. `git diff HEAD -- generation_service.go` = **vazio**.

**Correção sobre quantos subtestes discriminam.** A primeira versão deste EVIDENCE afirmou que o
discriminante era só `"listing has a seller_sku, ERP product has no refforn"` e que os outros dois
passavam sob os dois códigos. **Falso, e refutado pelo meu próprio raw:** `a3-mustfail-raw.txt:18`
mostra `listing_has_no_seller_sku` também FAIL. Discriminam **dois** dos três; só
`…has_a_refforn` passa nos dois. O erro era conservador — eu subestimei minha própria cobertura —
mas pack e raw divergiam, e isso conta.

### A4 — a âncora não some mais por comparando errado — **PASS, com a direção declarada**

A direção que sai agora é `UNAVAILABLE`, `Side` vazio: **continua não emitindo, mas agora por
comparar as coisas certas.** A2-R2 nega o `FOR` e continua valendo; o que este critério cobre é o
comparando. Antes o par era `seller_sku` contra `refforn`; agora é CODPROD contra CODPROD.

Asserção antes e depois, verbatim: `dispatches/a2-assertion-before-after.md` (HUB ruling R2).

### A5 — `vinculados` conta CODPROD com zero à esquerda — **PASS**

`TestGetImportChainCountsLeadingZeroCodprod`, contra Postgres real: `codprod = '00101'`, link
resolvido `internal_product_id = 101`, produto **DENTRO** de `vinculados`.

```
--- PASS: TestGetImportChainCountsLeadingZeroCodprod (0.03s)
```

Não fechado pelo `summary.txt` da lane — ver FINDING 1. **Nenhum dos dois lados do gate conseguiu
executar isto**; os dois graduaram derivado/NOT-PROVEN. Ver "O que o gate não pôde verificar".

### A6 — must-fail do A5 — **PASS**

Junção revertida para a forma antiga, por edição para a frente e restauração para a frente:

```
chain=domain.ImportChain{Protocol:"#640-E", Importados:1, Vinculados:0, Enfileirados:0, ...}
--- FAIL: TestGetImportChainCountsLeadingZeroCodprod (0.10s)
```

O número exato: **`Vinculados:0`** com `Importados:1` ao lado.

### A7 — o `DISTINCT` do `vinculados` continua vivo — **PASS na asserção que o critério pede, REPORT anexado**

Re-rodado 9 vezes. A asserção de contagens — `Protocol`/`Importados`/`Vinculados`/`Enfileirados`,
que é o que o guard DISTINCT fixa — passou **9/9**. Uma asserção SEPARADA no mesmo teste, a janela
de `QueueReadAt`, flakou **4/9**.

**Correção sobre o root-cause.** A primeira versão atribuiu o flake a skew de relógio de ~600ms do
container. O lado A do gate apontou que os números não fecham: um skew constante de 600 ms contra
um bracket host-side falharia 9/9 por ~600 ms, não 4/9 por ~1 ms. **Ele está certo.** O flake é
real e o grau REPORT continua honesto; o diagnóstico anexado **não está provado** e sai daqui como
não provado. Ninguém rodou esse teste em `5441fe18`, então nem eu nem o gate podemos chamá-lo de
pré-existente — está não provado nas duas direções. O que continua sustentado: a mudança de junção
não pode ser a causa, porque `statement_timestamp()` não é tocado pelo diff.

### A8 — `{id}` malformado responde 4xx nas DUAS rotas — **PASS**

`TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery`, tabela sobre
`GET /erp/imports/not-a-uuid` e `…/chain`. Status 400, corpo `{"error":"invalid_import_id"}`, e —
o ponto do teste — **assere que o querier não rodou**. Uma asserção só de status passaria também
num handler que consultasse primeiro e validasse depois.

OpenAPI e `packages/sdk-runtime` mudaram **no mesmo commit** que o comportamento (`055d1705`,
profile §7). `vitest` 4/4. Grant (g) exercido só no bloco `/erp/imports*`.

Metade que continua **NOT-PROVEN pelo gate**: nenhum dos dois lados tem git, então a atomicidade do
commit `055d1705` não foi verificada por eles — só por mim, que sou parte interessada.

### A9 — o comentário não afirma mais um universal falso — **PASS**

Texto novo inteiro, citado como o critério exige (é sobre a frase):

```go
// knownIdentityAnchors is the identity vocabulary THIS file governs: the
// anchors a connector declares for the product_links candidate generator. D-A
// (D-122) removed `refforn` from it — it is the supplier reference INSIDE the
// ERP (`ZP1704.1.`), no connector declares it, and keeping it here minted one
// permanent never-changing reason per candidate row. It stays a field on the
// ERP side (`erp_import_products.refforn`); what left is this list, not the
// datum.
//
// Deliberately NOT a claim that no marketplace datum ever compares against
// refforn. The catalog matcher in `market/domain/identity_resolver.go` scores
// the ERP `RefForn` against a candidate's `MODEL` attribute under the anchor
// name "refforn" today, and that resolver is wired; it runs on the market
// module's own vocabulary and never reads this list. A marketplace field that
// one day belongs HERE enters as a NEW anchor with its own name, not by
// reusing a term that means something else on the ERP side.
```

O lado A verificou as duas afirmações internas em vez de aceitá-las: "never reads this list" vale
(só `identity_anchor_adapter.go:28-29` e `marketplace_capability_service.go:138` chamam
`KnownIdentityAnchors()`), e "that resolver is wired" vale
(`collection_pipeline_service.go:148,168` → `composition/root.go:683`).

Vale registrar a ironia sem suavizá-la: **o critério que apaga um universal falso passou, e o chip
escreveu outro universal falso três arquivos ao lado.** Ver A2.

### A10 — `pending` não-array não derruba o endpoint — **PASS**

Dois fixtures contra Postgres real — `cursor -> 'pending'` como objeto e como escalar — ambos
devolvem resposta válida com `enfileirados = 0`. Must-fail contra a forma antiga (`COALESCE` só):

```
ERROR: cannot extract elements from an object (SQLSTATE 22023)
```

O guard é `CASE`, não `AND`: predicado de junção/filtro não dá garantia de ordem de avaliação.

### A11 — `*comparison.product` não deref nil — **NOT-PROVEN em produção**, guard com must-fail

Os dois lados do gate chegaram independentemente ao mesmo grau. **Um** call site em produção,
`generation_service.go:303`, passando `&product` onde `product := skuMatches.Products[0]` —
endereço de local, nunca nil. Função não exportada, sem despacho dinâmico.

Varredura da rodada 3 sobre esta frase: ela dizia "sem segundo caminho", e isso é falso lido ao pé
da letra — `buildConcordantCandidate` tem **dois** call sites, `:303` (produção, `&product`) e
`generation_service_test.go:580` (teste, `nil` de propósito), além da definição em `:489` e de duas
menções em comentário (`generation_service.go:231`, `auto_link_policy_test.go:290`). O segundo call
site é o guard, não um caminho de produção; a frase agora diz o que se sustenta. Linha reconferida
na reconciliação da rodada 4 — era `:568` e virou `:580` por causa dos meus próprios comentários em
`590efdc8`.

O guard **tem** teste direto: `TestConcordantCandidateDoesNotDerefNilProduct` chama
`buildConcordantCandidate(...)` sem service, com `newProviderIdentityAnchorComparison(snapshot, nil, nil)`.
Mesmo pacote, função não exportada alcançável do teste. Então o guard não é defesa não-testada; o
que continua NOT-PROVEN é só a alcançabilidade em produção.

**A invariante que torna o nil inalcançável NÃO é imposta pelo compilador** (exigência do hub). O
campo é ponteiro; o tipo admite nil. A invariante é o único call site passar endereço de local, e
ela vive em código, não no sistema de tipos. Quem apagar o guard por parecer morto reintroduz o
deref no dia em que aparecer um segundo call site.

**Achado do lado A que eu não tinha visto, e que é pior que o meu grau sugere:** o guard degrada
para um candidato **corroborado**. Com produto nil, a função segue com `ProductCandidate{}` e
constrói `seller_sku`/`FOR` e `ean`/`FOR` com `Confidence: 95`, banda `ALTA`, status `ACCEPT` —
duas âncoras alegando resolver para um codprod ausente.

**Aqui eu aleguei cobertura que não existia, por duas versões deste pack.** Escrevi que
`TestConcordantCandidateDoesNotDerefNilProduct` "fixa exatamente isso" quando o corpo dele só
asseria `InternalProductID != nil` e a presença de um motivo `seller_sku` — nada de `Confidence`,
banda ou status. O round 2 pegou. **Desde `0264ba84` a frase é verdadeira:** o teste assere `ean`
FOR, `Confidence == 95`, `ConfidenceBand` ALTA e `MatchStatus` ACCEPT, por constantes de domínio,
cada uma com must-fail próprio.

Sobre o comentário: ele afirma paridade no **nil-check** com as duas irmãs — e isso é verdade,
`applySingleAnchorScore` tem o mesmo `if comparison.product != nil`. O que ele não dizia, e um
leitor inferia, é que a **degradação** também teria paridade: as irmãs degradam para motivos de
ausência, este degrada para corroboração. Duas frases somadas em `2bed7d9d` no lado da produção,
sem apagar o texto existente (ele não é falso — ver R5); a mesma negação explícita entrou no
comentário do teste em `0264ba84`. É a mesma classe do B-01. Inalcançável hoje; sai como REPORT R5,
não corrigido, porque mexer nisso é redesenhar o contrato da função e está fora do meu contrato.

### A12 — nenhum dano colateral — **PASS, com a citação de artefato corrigida**

**O que reprovou no round 1:** a primeira versão escreveu "Bruto: `dispatches/ladder-l0-l1-raw.txt`"
e pendurou nele `go build`, `go vet` e os guards `-count=10`. **Esse arquivo não contém nenhum dos
três** — `grep -rn "count=10|go vet|go build" dispatches/` devolvia zero hits para `count=10`. E a
linha 3 dele estava corrompida (`[no test fi107`), então nem a lista que ele de fato tinha era
certificável. Os números eram verdadeiros; **a citação era falsa**. As duas coisas contam sob R-24.

Recapturado íntegro em `dispatches/ladder-l0-l1-raw-r2.txt`, com o comando acima de cada saída:

```
$ go build ./...                                              EXIT=0
$ go vet ./...                                                EXIT=0
$ go vet -tags=integration ./internal/modules/erp_import/...   EXIT=0
$ go test -count=10 -run 'TestProviderUnavailableReasonOrderingIsStable|TestHardNegativeDimensionStableEquivalentDisplay' -v
     10 --- PASS: TestHardNegativeDimensionStableEquivalentDisplay
     10 --- PASS: TestProviderUnavailableReasonOrderingIsStable
```

20/20, zero FAIL. `go test -count=1 ./...`: **107 pacotes `ok`**.

**Um FAIL apareceu na recaptura e não é regressão** — registro porque omitir seria a mesma doença:
`dashboard/transport [build failed]` com `fatal error: runtime: cannot allocate memory`. OOM do host
linkando o binário de teste com os dois agentes de gate rodando junto. `go build ./...` passou
porque não compila arquivos de teste. Re-rodado isolado: `ok …/dashboard/transport 1.969s`, EXIT=0.
106 + 1 = 107. O pacote não está no meu diff.

Suíte de política D-121 depois do R2: **17 top-level PASS, 0 FAIL** na família consultada pelo
worker (16 antes, +1 do teste novo). `vitest` do `sdk-runtime`: 4/4.

Escopo, contra o `main` de verdade — `base_sha` é piso:

```
git diff --name-only 5441fe18… HEAD -- apps contracts packages   → 11 arquivos
```

**Zero** sob `apps/web/`, **zero** migrations, **zero** `platform/httpx`. `main` moveu para
`81457e4f` e o que ganhou toca só `.mnfs/`.

#### Lane de governança — vermelha, com delta ZERO deste chip

`status=failed`, exit 1, **53** achados no HEAD e **53** no `base_sha`, rodados em worktree limpo e
destacado. Brutos: `governance-raw.txt`, `governance-baseline-raw.txt`. Achados em `sourcekind`,
`tenant_config`, `catalog`, `internal_read` — nenhum no meu diff. O lado A comparou contagem,
primeiras 20 linhas e linhas 125-179 dos dois e bateu; não diffou linha a linha por não ter `diff`.
Worktree temporário removido.

### A13 — o EVIDENCE aponta, não recopia — **PASS** (os dois lados concordam)

---

## Gate P6 — round 1: **os dois lados REFUTARAM**

Input congelado do round 1: `p6-input-r1.patch`, sha256 `b928bcef…`, 11 arquivos, +440/-42.
Brief idêntico: `dispatches/p6-gate-brief.md`. Os dois rodaram concorrentes e cegos um ao outro, e
os dois confirmaram por escrito que não abriram a saída do outro. Esta sessão não assinou nenhum
lado (HUB R1 condição 2).

Input do round 2, se o hub autorizar: `p6-input-r2.patch`, sha256
`d06f98814868e8fb17fe9238569f15cca4d3648f5824254309e1862f947031b2`, 11 arquivos, +481/-42.

| Critério | Lado A (Opus) | Lado B (Sol medium) |
|---|---|---|
| A1 | PASS | PASS |
| A2 | **FAIL** | PASS |
| A3 | PASS (derivado) | NOT-PROVEN |
| A4 | PASS | PASS |
| A5 | PASS (derivado) | NOT-PROVEN |
| A6 | PASS (derivado) | NOT-PROVEN |
| A7 | REPORT | REPORT |
| A8 | PASS (metade same-commit NOT-PROVEN) | PASS |
| A9 | PASS | PASS |
| A10 | PASS (derivado) | NOT-PROVEN |
| A11 | NOT-PROVEN | NOT-PROVEN |
| A12 | **NOT-PROVEN** | NOT-PROVEN |
| A13 | PASS | PASS |

### A pergunta nomeada (HUB R2)

> "a asserção nova codifica o comportamento correto, ou codifica o código que foi escrito?"

**Os dois lados responderam REPORT, e responderam a mesma coisa.** Nenhum dos dois conseguiu
afirmar a primeira opção por evidência, e é isso que a instrução mandava fazer nesse caso.

O lado A verificou por conta própria e concluiu: codifica o código que foi escrito, **mas o código
que foi escrito é a regra pré-existente e independentemente fixada da âncora irmã**, então não é
regressão enterrada. Os três fatos que ele reproduziu sozinho: a asserção antiga só era verdadeira
sob o defeito (fixava a frase "o produto ERP não tem CODPROD" sobre um produto com CODPROD 300); a
asserção nova é **estritamente mais forte** e falha quando o fix é removido, logo é guard e não
carimbo; e a simetria com `TestExactSKUWithUnmatchedListingEANKeepsSeededEANReason` é **real** — ele
confirmou que aquele teste é anterior ao chip lendo o primeiro hunk do patch naquele arquivo
(`:397`), acima do teste em `:367-388`. Regra "criado vs. tornado frequente": ele decidiu
**tornado frequente**, no código — o ramo, a constante de direção e o `Side` vazio estão todos
inalterados no diff.

Por que ainda assim REPORT, nas palavras dele: `UNAVAILABLE` está documentado em
`generation_service.go:709` como *"provider não fornece a âncora %s"*, e neste fixture o provider
**fornece**. A mudança troca uma afirmação falsa por uma imprecisa que o módulo já fazia em todo
lugar. Não é o mesmo que correta.

O lado B chegou ao mesmo lugar por caminho próprio.

### Os três blocking do round 1 — todos confirmados por mim, nenhum contestado

1. **`generation_service_test.go`, o comentário do universal falso.** Confirmado lendo
   `missingMatchedAnchorReason:640-645` e o call site `:216`. Corrigido em `54342331`. Ver A2.
2. **`EVIDENCE.md`, o mesmo universal falso, usado para justificar a deleção de um caso.**
   Confirmado. Esta versão do documento é a correção. Ver A2.
3. **A12, artefato citado que não continha três das quatro coisas citadas.** Confirmado por
   `grep`: zero hits de `count=10` em qualquer raw. Recapturado. Ver A12.

Os dois blocking do lado B são de outra natureza e **não apontam defeito no código**: ele não
conseguiu executar nada. Verbatim: *"the read-only environment denied creation of Go's temporary
work directory"*, *"I could not provision or write to a PostgreSQL test database"*. Todos os
NOT-PROVEN dele são execution-shaped. Ele inspecionou SQL, fixtures e diff e não achou erro em
nenhum. Isso não é atenuante do meu lado — é FINDING 4 sobre a forma do gate.

### O que o gate não pôde verificar, e que continua não verificado por ninguém independente

Registro porque um gate REFUTED cuja cobertura fica implícita vira alegação inflada ao contrário.

- **Execução de qualquer teste.** O lado A é fisicamente read-only sem Bash; o lado B teve o
  diretório de trabalho do Go negado. Os **quatro** `PASS (derived)` do lado A (contagem por
  `grep -o "PASS ([a-z]*)" dispatches/p6-opus-gate-r1.md`, e o round 3 não emitiu nenhum) são
  hand-trace de semântica
  determinística de Go/SQL mais leitura do fonte, cruzado com stdout colado — não é prova
  executada. A5, A6, os 9 runs do A7, A10, a suíte D-121, os 107 pacotes e o `vitest` 4/4
  descansam todos em saída que **eu** produzi.
- **Fatos de git.** Nenhum dos dois tem git. HEAD, os file sets dos commits, a atomicidade de
  `055d1705` (obrigação do A8), o diff vazio após os round-trips de must-fail e o conteúdo do
  `81457e4f` foram verificados só por mim. O lado A registrou que o snapshot da sessão dele mostrava
  a branch em `bcab8269` com os 11 arquivos **não commitados** — snapshot velho, não evidência
  contra o chip, mas significa que ele certificou a **árvore de trabalho**, não os commits.
- **O sha256 do patch congelado.** O lado A leu o arquivo e não pôde hasheá-lo.
- **Byte-identidade das duas rodadas de governança.** Comparadas por contagem e amostras, não linha
  a linha.
- **Se o flake do A7 reproduz em `5441fe18`.** Ninguém rodou lá. Não provado nas duas direções.
- **Consumidores em `apps/web` de `GET /erp/imports/{id}`.** Fora de escopo; ninguém checou se algum
  chamador manda id não-canônico e agora recebe 400 em vez de 200.

---

## Gate P6 — round 2: **SPLIT**, e o lado que reprovou está certo por STRING

Forma de três assentos (profile §11, ratificada em `5ec83724`): dois assentos de leitura + o
assento executável, que é do hub. Critério de execução saiu de escopo dos dois lados por brief.

| Lado | Agente / modelo | Veredito | Artefato |
|---|---|---|---|
| A | `harness:gate-reviewer`, Opus, read-only | **REFUTED** — 4 blocking, 3 não-bloqueantes | `dispatches/p6-opus-gate-r2.md` (14910 B) |
| B | `codex:codex-rescue` → `gpt-5.6-sol` medium | **CONFIRMED** — "Findings: None" | `dispatches/p6-sol-gate-r2.md` (6396 B) |

**Não há AGREEMENT. Não fecho.** E o split não é empate: verifiquei os quatro blocking do lado A
por string, sozinho, e três procedem. O CONFIRMED do lado B não os contradiz — ele **não varreu**.

| Blocking do lado A | Minha verificação | Estado |
|---|---|---|
| 1 — o universal falso sobrevive num **segundo** sítio, `generation_service_test.go:505-506` | `grep "unreachable for seller_sku"` → **1 hit**, verbatim `// side=erp is unreachable for seller_sku now, so the only honest` | **PROCEDE — CORRIGIDO `0264ba84`**, grep = 0 |
| 2 — `:541` `// that this scorer now degrades like both of its siblings instead of panicking.` é a paridade de **degradação**, que é falsa; e o teste não fixa 95/ALTA/ACCEPT | `grep "degrades like both of its siblings"` → **1 hit**; corpo do teste só assere `InternalProductID != nil` e a presença do motivo `seller_sku` | **PROCEDE — CORRIGIDO `0264ba84`**, grep = 0 e o teste passa a asserir |
| 3 — o pack não descreve o HEAD contra o qual foi submetido; o R4 reporta como aberto um defeito que o delta fechou | **OBSOLETO por edição minha**, não por refutação: `head: 2bed7d9d`, linha de commit presente, `grep "continua juntando cru"` → **0 hits**. O lado A leu a versão anterior. Está declarado na nota de fase da tabela de correções. |
| 4 — `dispatches/p6-opus-gate-r1.md` citado e **inexistente** | `glob p6-*gate-r*.md` confirmou: só existia em `_hub-gate-anchors-2/`, que é outro chip | **PROCEDE — e é meu** |

**B4 é a falha que eu mesmo nomeei como doutrina e depois cometi.** "Streamar não é persistir." O
assento read-only do round 1 não tem Write; devolveu o veredito como mensagem final e eu **não
persisti**. Seis afirmações do pack atribuídas a "o lado A" ficaram irrecuperáveis por um leitor.

Recuperado agora do transcript da sessão, verbatim, casando o task-id `ad20d064691a9fcb2`: 21711
chars, primeira linha `VERDICT: REFUTED`, contém a string conhecida `never through the product`.
Gravado em `dispatches/p6-opus-gate-r1.md` com cabeçalho de proveniência que diz o que é —
**recuperação de transcript, não escrita contemporânea**. Vale menos que um artefato escrito na
hora, e o cabeçalho diz isso.

**O que o split ensina sobre o lado B.** Ele deu PASS em (a) checando **o sítio nomeado no brief** e
parando ali. O brief dizia "decida por si se é verdade do código" e ele decidiu — sobre uma linha.
Duas frases falsas no mesmo arquivo, ambas achaveis por `grep`, passaram. `Findings: None` num
arquivo que tem duas é a forma que a doutrina chama de falso-negativo por cobertura: **o veredito
não estava errado sobre o que leu; estava errado sobre o que não leu.** Um CONFIRMED sem varredura
declarada não corrobora o REFUTED do outro lado — não o toca.

### Correção dos blocking 1 e 2 — `0264ba84`

Worker despachado (ledger linha 11), só `generation_service_test.go`, 24 inserções / 3 deleções.
`generation_service.go` **intocado**. Relatório bruto: `dispatches/impl-r7-second-false-universal.md`
(12606 B).

**Blocking 1.** O comentário sobrevivente passou a dizer só o que a tabela cobre — "Every case in
this table runs with a product present" — e aponta para onde o caso nil-product está fixado. Não é
mais universal: é limitado ao escopo em que é verdadeiro. `grep "unreachable for seller_sku"` = 0.

**Blocking 2, metade prosa.** A paridade de nil-check ficou (é verdade, é o seu ruling do R5). A
paridade de degradação foi negada explicitamente: as irmãs degradam para ABSENCE, esta degrada para
CORROBORATION sobre um `ProductCandidate{}` zerado. `grep "degrades like both of its siblings"` = 0.

**Blocking 2, metade código — esta é a que importa.** Eu tinha alegado que o teste fixava
95/ALTA/ACCEPT. Não fixava. Agora fixa, e por constantes de domínio, não literais: motivo `ean` FOR,
`Confidence == 95`, `ConfidenceBand == …BandAlta`, `MatchStatus == …StatusAccept`. **A correção
certa era fortalecer o teste, não abrandar a frase** — abrandar teria deixado o REPORT R5 sem nada
segurando, que era exatamente a acusação.

Must-fail por classe de asserção, contra código de produção mutado, cada um com actual-vs-want
próprio (bruto no relatório, linhas 128–186):

```
generation_service_test.go:574: reasons=[]domain.LinkCandidateReason{...Anchor:"seller_sku"...}, want the concordant ean FOR reason
generation_service_test.go:577: confidence=90, want 95 for a nil-product concordant candidate
generation_service_test.go:580: confidence_band="MEDIA", want "ALTA"
generation_service_test.go:583: match_status="CONFIRM", want "ACCEPT"
```

Nenhuma mutação foi no-op. Restauração por edição para a frente; `git diff` de
`generation_service.go` vazio. Ladder do worker: `go build ./...` limpo, `go vet` do módulo limpo,
`-count=1 -v` PASS nos dois testes tocados, `-count=10` `ok` nos seis pacotes.

**Estado da árvore, dito com precisão em vez de "limpa".** `git status` mostra
` M generation_service.go`, e **não é uma mudança**: 41615 bytes dos dois lados, zero CR dos dois
lados, `diff` após normalizar = idêntico, `git diff` vazio. É stat-cache sob `autocrlf`, que
`git update-index --refresh` não limpa. Registrado porque "árvore limpa" seria alegação falsa e
`git checkout --` é proibido.

**Patches recongelados** (os do r3 ficam obsoletos junto com os do r2):

| Patch | sha256 |
|---|---|
| `p6-input-r4.patch` (cumulativo, 11 arquivos) | `085882a6085546666a2aa56adfb683ca3d5c551f73c00c3d406d62b6be2a7444` |
| `p6-delta-r3-to-r4.patch` (1 arquivo) | `4943e3c3f9c7e4972641b79c90e2345775a17016b35c68b102cb3f91cd37cbd8` |

Escopo do cumulativo reconferido em `0264ba84`: 11 arquivos, zero `apps/web`, zero migrações, zero
`platform/httpx`.

---

## REPORTs — o que este chip NÃO fechou

### R1 — a direção de `both-present-and-DIFFERENT` continua errada, e o CORR-1 a torna frequente

`UNAVAILABLE` sob D-B significa "o provider não fornece essa âncora". No caso que este chip torna
comum, o provider **fornece** — o anúncio tem `seller_sku`, ele só aponta para outro CODPROD.

**O CORR-1 não CRIA esse defeito; torna-o frequente**, e os dois lados do gate ratificaram essa
distinção examinando o diff. Antes, `seller_sku` só caía no `default` quando o produto ERP tinha
`refforn` preenchido; agora cai sempre que o anúncio traz um `seller_sku` **e há produto resolvido**
— porque aí o `productValue` é o CODPROD, que o produto resolvido sempre tem, então o `switch` de
`missingMatchedAnchorReason:644-652` chega ao `default`. Sem produto resolvido o ramo é o
`case product == nil`, que dá `INCOMPARABLE`/`erp`, não `UNAVAILABLE`.

**O balde errado não é do `seller_sku`.** `ean` já cai no mesmo `default` quando os dois lados
divergem, fixado por `TestExactSKUWithUnmatchedListingEANKeepsSeededEANReason` desde antes deste
chip — o lado A confirmou essa anterioridade pelo próprio patch. A decisão A2-R1 governa as **duas**
âncoras primárias, e uma delas já se comporta assim em produção sem nenhum chip ter tocado.

### R2 — A7: flake de janela de tempo, com diagnóstico NÃO provado

4/9 falhas na asserção de bracket de `QueueReadAt`. O root-cause de skew de container que eu anexei
não fecha com os números (ver A7) e sai como não provado. O grau REPORT do resultado continua.

### R3 — a lane de governança está vermelha no `base_sha`

Delta zero deste chip. Dono: o hub / os chips dos módulos citados.

### R4 — primeira metade REJEITADA contra dado real; segunda metade CORRIGIDA em `2bed7d9d`

**A supercontagem que eu reportei não procede, e a razão da rejeição é melhor que a minha alegação.**
O hub mediu contra o banco real (10584 linhas de mirror): `mirror_leading_zero`, `mirror_ltrim_collisions`,
`mirror_nonnumeric`, `links_leading_zero`, `import_leading_zero`, `collisions_same_import` — **zero
em todos**. E, mais importante que o número: `importados` é `count(*)` de linhas **cruas** de
`import_products`, então `SELECT DISTINCT products.codprod` cru é a chave **consistente com a
população que o operador vê na tela**. Um import com `'00101'` e `'101'` dá importados=2,
vinculados=2 — coerente. Canonicalizar a chave é que seria inconsistente: daria vinculados=1 sobre
importados=2 e o operador leria "um não vinculou". Apliquei a intuição certa ao contador errado.

**A segunda metade era real e é minha.** `queued_products` juntava CRU
(`ON products.codprod = pending.codprod`) enquanto `resolved_products` juntava canonicalizado.
Duas regras de identidade na mesma tela: com um codprod carregando zero à esquerda, `vinculados`
conta e `enfileirados` não, e os três números deixam de ser uma decomposição da mesma população.

Corrigido em `2bed7d9d`, por worker despachado (ledger linha 8): `queued_products` passa a
`ON ltrim(products.codprod, '0') = ltrim(pending.codprod, '0')`. **A chave contada continua CRUA nos
dois CTEs** — canonicaliza o predicado, não a chave, que é exatamente a distinção que a rejeição
acima estabelece. Guard `jsonb_typeof` intocado.

Regressão contra Postgres real: `TestGetImportChainCountsLeadingZeroCodprodInQueue`, `'00101'`
importado contra cursor carregando `"101"`, mais uma linha simétrica (`'102'` vs `"00102"`) para que
um fix meio-aplicado também falhe. Must-fail com a junção crua restaurada:

```
chain=domain.ImportChain{Protocol:"#660-E", Importados:3, Vinculados:1, Enfileirados:0, ...}
```

`Vinculados:1` e `Enfileirados:0` **para o mesmo CODPROD** — a decomposição quebrando, com número.
Restaurado por edição para a frente.

**DESCARREGADO POR EXECUÇÃO INDEPENDENTE — assento executor do hub, em `0264ba84`.** Este era o
único critério do pack que dependia inteiramente da minha palavra: eu produzi a mutação, rodei e
relatei. O hub repetiu em worktree destacado dele, `git diff` = 0 bytes, contra Postgres real, 69
migrações aplicadas. Junção crua restaurada:

```
migrations_first=69
failure_token=test=TestGetImportChainCountsLeadingZeroCodprodInQueue
status=blocked
postgres lifecycle failed reasons=HPG_TEST_FAILED exit_code=1
```

Restaurada a canonicalização, mesma lane:

```
migrations_first=69
status=passed
run_id=3191c4b7e8664a42971f61e3b8013bde
```

Vermelho com a junção crua, verde com a canonicalizada, **medido por quem não é o autor**. O ladder
do mesmo `0264ba84` no worktree dele: `go build` EXIT 0, `go vet ./...` EXIT 0, `go test ./...` 153
pacotes com `ok=107` / `no-test-files=46` / `FAIL=0` — o `ok=107` bate com o meu número, contado e
não lido no rabo da saída.

**Achado colateral, do hub e não meu:** a run mutada derrubou quatro pacotes
(`erp_import/adapters/postgres`, `market/adapters/postgres`, `mutations/application`,
`orders/adapters/postgres`), e a mutação é uma linha num CTE do `GetImportChain`, sem alcance sobre
`market`, `mutations` ou `orders`. `Postgres.psm1:209` dá **um** banco por run
(`mpc_test_$RunId`) e os pacotes de integração rodam em paralelo contra ele, então contaminação
cruzada é estruturalmente possível. Qual dos dois — contaminação ou flake independente — não foi
isolado, e o hub não inventou. Consequência que já vale para mim: **`failure_token=package=` não é
atribuível sozinho** (o pacote listado pode ser vítima); citar sempre `failure_token=test=`, que é
específico.

Latente nos dois sentidos, com a razão achada e não suposta: o hook de import
(`import_service.go:153-158` → `enqueuer.go:44`) escreve o codprod **cru** de volta, então aquele
produtor não consegue discordar de si mesmo. A exposição é um produtor que sirva do lado
inteiro — a mesma perda que o `resolved_products` já teve que reparar. Latente é esperando o dado.

**Sargabilidade → G4, com a correção do hub:** `ltrim(...)` dos dois lados é não-sargable, mas o
lado `links.internal_product_id::text` **já era** não-sargable na base, antes deste chip. Eu
dupliquei, não introduzi.

**Dívida aberta, não consertada:** `ltrim(x,'0')` colide `'0'`/`'00'`/`'000'`. Pré-existente no
`resolved_products` desde o CORR-2; o worker estendeu ao `queued_products` por consistência em vez
de criar uma terceira regra de identidade. Não testado. Fechar isso é fatia nova e precisa de grant.

### R5 — o comentário NÃO é falso; o que faltava é dizer o que a degradação PRODUZ

**Correção de uma autodenúncia que passou do ponto.** Eu tinha marcado o comentário do guard nil
como falso e ia deletá-lo. O hub checou e está certo: as duas metades são verdadeiras — as irmãs
nil-checam mesmo (`applySingleAnchorScore` tem o mesmo `if comparison.product != nil`) e este sítio
de fato dava panic. Deletar teria sido o erro espelhado do round 1: apagar verdade em vez de
apagar falsidade. R-25 não se aplica.

O que é verdade é mais estreito: o comentário afirma paridade no **nil-check**, e um leitor lê
paridade na **degradação**, que não existe. Corrigido em `2bed7d9d` com duas frases a mais, texto
existente intocado: a degradação aqui emite `seller_sku` FOR e `ean` FOR a `Confidence: 95` /
`ALTA` / `ACCEPT` sobre um `ProductCandidate{}` zerado, e `autoApprovals` ainda lê esse ACCEPT como
auto-aprovável.

Caminho segue inalcançável **em produção** hoje (o único call site de produção é `:303`; o teste é
o outro call site e passa nil de propósito); REPORT com yes-if nomeado — se um segundo call site de
produção aparecer, isto vira defeito de operador na mesma classe do B-01.

**Uma alegação minha aqui era falsa e o round 2 pegou.** Eu escrevia que o teste "fixa exatamente
isso" (95/ALTA/ACCEPT) enquanto ele só asseria `InternalProductID != nil` e a presença de um motivo
`seller_sku`. Alegação total, cobertura parcial — R-24. Consertado em `0264ba84` **pelo lado do
teste**: agora ele assere `ean` FOR, `Confidence == 95`, banda ALTA e status ACCEPT por constantes
de domínio, com must-fail por asserção. A forma degradada está fixada de fato, e a frase virou
verdadeira em vez de ser abrandada.

### R6 — a premissa estava errada; o achado, corrigido, é MAIOR

**Minha formulação estava mal colocada.** Eu escrevi "se um provider declarar `seller_sku` com
`Supplied: false`". `Supplied` não é algo que um provider declara: o adapter enumera **todas** as
`KnownIdentityAnchors()` e marca `Supplied = (anchor ∈ declaredSet)`
(`identity_anchor_adapter.go:28-35`). `Supplied:false` é a representação **normal** de "este
provider não fornece esta âncora".

E o ramo **dispara hoje**, todo dia: `knownIdentityAnchors` tem quatro entradas — `seller_sku`,
`ean`, `title`, `marca` (`marketplace_capability.go:40-45`) — e o Mercado Livre declara três
(`capability_adapter.go:90`). Então `marca` entra em `classifyProviderIdentityAnchor` com
`!anchor.Supplied` em produção.

Para `seller_sku` especificamente segue latente, porque o ML declara. **O gatilho tem nome:** o dia
em que a lista de capacidade do ML perder `seller_sku`, ou um segundo marketplace não declarar —
momento em que a promoção substitui o motivo semeado por inteiro e derruba o `side=erp` que o
`54342331` acabou de restaurar. Varredura feita pelo hub: um único adapter de marketplace
registrado (`melhorenvio` é frete, `postgres` é persistência). REPORT com gatilho nomeado.

### R8 — exposição de credencial em transcript de worker

O worker do R4 imprimiu a senha do container Postgres de sessão no transcript dele enquanto
localizava o alvo. Container local descartável em tmpfs, senha de teste, morre com o container — mas
a regra de não expor segredo não abre exceção por ser throwaway. Reportado ao hub sem minimizar;
rotação ou expurgo é seam dele.

### R7 — `IsValid` é mais estrito que `uuid_in`, e o texto do OpenAPI é mais largo que o guard

`IsValid` exige a forma canônica de 36 chars; `uuid_in` também aceita 32 hex sem hífen e formas com
chaves. Ids de produção são canônicos (`import_service.go:210`), a direção da falha é conservadora,
e os dois reviewers recomendam deixar. Mas o OpenAPI descreve o 400 como *"Malformed ERP import id
(not a UUID)"*, mais largo do que o guard significa. Não corrigido no round 1 para não invalidar o
patch congelado enquanto o gate rodava; agora é dívida nomeada.

### Fora de contrato, com dono nomeado

**B-02** (`apps/web` `QueueRow.tsx:159`) → CHIP-VINC-NEUTRO. **B-08** (deadline de `platform/httpx`)
e **G4** (índice) → chips próprios. Erros de `tsc` do `apps/web` não são meus. **Nenhum L2** — este
chip não subiu servidor, não ligou em `:8080`, não leu `.env*`.

O lado A registrou que `LIVE-VERIFIED:` / `LIVE-WAIVED-BY-OPERATOR:` e qualquer `EXEMPLO-IO` estão
**ausentes do pack inteiro**. Ele não cobra isso do chip porque o contrato exclui L2, mas anota que
a milestone não fecha sobre este pack como está. Repasso sem atenuar.

---

## Gate P6 — round 3: os DOIS lados REFUTARAM

| Lado | Assento | Veredito | Artefato |
|---|---|---|---|
| A | `harness:gate-reviewer`, Opus, read-only | **REFUTED** — 2 blocking, 3 não-bloqueantes | `dispatches/p6-opus-gate-r3.md` |
| B | `gpt-5.6-sol` / medium, sandbox read-only | **REFUTED** — 5 blocking | `dispatches/p6-sol-gate-r3.md` |

Não há `AGREEMENT`. O close continua BLOQUEADO e a linha `P6-DUAL-GATE:` continua sendo do hub.

**Os DOIS artefatos foram persistidos por MIM, não pelo assento** — e por motivos diferentes, o que
é o achado desta rodada. O lado A não tem Write (assento fisicamente read-only) e devolve o veredito
como mensagem final. O lado B **tentou** escrever o caminho que o brief mandou, com
`apply_patch`/`*** Add File:`, e o runtime recusou: `patch rejected: writing is blocked by read-only
sandbox; rejected by user approval settings` (rollout `019fa9da-…`, linha 127). O corpo do lado B
salvo em disco é o **payload literal daquela escrita recusada** — o hunk `+`-prefixado, 208 linhas,
16438 chars, prefixo removido, zero linhas sem prefixo — e não uma prosa recuperada. Os dois
cabeçalhos dizem o que são. O brief não pode descarregar persistência mandando o assento persistir.

### O que eu verifiquei por STRING, achado a achado

Não arbitrei entre os dois lados. Conferi cada achado no código.

| Achado | Lado | Verificação minha | Veredito | Autoria |
|---|---|---|---|---|
| `generation_service_test.go:642-643` "the one direction no other seller_sku assertion in this file covers" | A (blocking 1) | `sed -n '1803,1832p'` mostra `TestCase8…` asserindo `findReason(…,"seller_sku",…Incomparable)` + `Side != …SideERP` → `t.Fatalf`. E `grep -c "MLB-FX8" p6-input-r4.patch` = **0** → o Case 8 é anterior a este chip e nunca foi tocado | **PROCEDE** | **minha** (`p6-input-r4.patch:747`) |
| `generation_service_test.go:400-403` "a `ProductCandidate{}` is unreachable in production" | A (blocking 2) + B (blocking 1) | `grep "internalreaddomain.ProductCandidate{}" generation_service.go` → **5 hits**: `:215`, `:340`, `:379` passam o zero literal a `newCandidate` em caminho de produção; `:497` e `:523` o constroem | **PROCEDE** | **minha** |
| `:552-559` "the siblings degrade into ABSENCE reasons" / "nil-checks like both of its siblings" | A (não-bloqueantes 1 e 2) | `sed '519,556p'`: `applySingleAnchorScore` no `StateExactSKU` emite `{seller_sku, For, "seller_sku resolve exato para codprod"}` a 70/MEDIA/CONFIRM mesmo com produto nil — degrada só o motivo que LÊ o produto. `sed '626,640p'`: `applyUnresolvedScore` não checa nada, passa `nil` fixo | **PROCEDE** | **minha** (delta `0264ba84`) |
| `generation_service.go:85-89` + `generation_service_test.go:222-224` "a capped run leaves every uncapped listing without a candidate" | B (blocking 2) | `link_candidate_repo.go:38-48`: o `DELETE` roda **por identity fornecida** (`provider_item_id = $3 AND provider_variation_id = $4`), então um anúncio fora do cap **retém** o candidato anterior | **PROCEDE** | pré-existente (`grep "capped" p6-input-r4.patch` = 0) |
| `generation_service_test.go:1636-1637` "title … can never grant ACCEPT/REVIEW-grade confidence" | A (não-bloq. 3) + B (blocking 3) | `generation_service.go:559`: `case StateTitleMatch: … MatchStatusReview`. Concede REVIEW | **PROCEDE** | pré-existente (`grep "never grant"` = 0) |
| `generation_service.go:697-699` "promote only its classification and retain that sentence" | B (blocking 5) | `sed '693,706p'`: a retenção do `Detail` está **dentro** de `if direction == …Incomparable`. Em `UNAVAILABLE` a sentença semeada é substituída | **PROCEDE** | pré-existente |
| `generation_service.go:858-861` "normal titles without measurements are never flagged" | A (não-bloq. 3) | `sed '875,878p'`: os matches de `hardNegativeSizePattern` entram em `pairs`, então título sem medida nenhuma é sinalizado — e `generation_service_test.go:1517` assere `detectHardNegative("Camisa M","Camisa G")` | **PROCEDE** | pré-existente |
| `generation_service.go:473` "seller_sku and ean are the only cross-side anchors" | B (blocking 4) | A frase **continua** na mesma sentença: "title ranks only, never accepts". O parágrafo nomeia o título; "only cross-side anchors" lido no contexto é "as únicas âncoras que corroboram". O lado A leu a mesma linha e a graduou TRUE | **NÃO ACEITO** — registro a divergência em vez de somar ao placar | pré-existente |

**Consertado nesta rodada (o que é meu):** os dois blocking e os dois não-bloqueantes de autoria do
chip, todos por estreitamento sob R-25 — a metade falsa saiu, não foi anotada. Nenhum teste foi
apagado desta vez; o comentário de `:642` agora **nomeia o Case 8** como cobertura anterior e diz o
que este teste acrescenta sobre ele (o `Detail`), em vez de reivindicar exclusividade.

**REPORT ao hub (não é meu, não está no meu diff):** os quatro universais pré-existentes
(`generation_service.go:85-89`, `:697-699`, `:858-861`, `generation_service_test.go:222-224`,
`:1636-1637`). São da mesma classe e são reais, mas consertá-los é ampliar escopo num arquivo cujo
texto eu não escrevi.

**DONO: o hub** (RULING 2). Vão para um commit dele na `main` **depois** do meu merge, e estão
**FORA do escopo do veredito** do round 4 — declarado nessa forma no brief dos dois assentos, para
que um assento frio não os leia como achado aberto. O hub mediu dois deles antes de assumir e
confirmou as duas refutações: em `:85-89` o `DELETE` de `link_candidate_repo.go:38-48` é escopado às
identidades fornecidas, e o que fica no lugar da frase falsa é pior — anúncio fora do cap **retém o
candidato antigo**, em silêncio, e isso virou tarefa de produto na fila dele; em `:697-699` a
promessa incondicional é condicional no código. Os dois são defeito de **prosa**, não de código.

### SWEEP do autor sobre o próprio pack — obrigação imposta pelo hub

> "a SWEEP não é só dos revisores. Você roda a mesma varredura contra o pack inteiro, com o
> resultado no EVIDENCE. … Varredura que o autor roda contra a própria prosa antes de publicar é o
> que estava faltando — e é barata, é um grep."

Varrido: `EVIDENCE.md` inteiro (745 linhas). Padrão:
`nunca|sempre|apenas|inalcanç|não pode|impossív|todo |toda |todos|todas|nenhum|qualquer|única|único|só o|só a|never|always|only|unreachable|cannot|every|no longer`.
**71 linhas com hit.** Triagem: 38 são citação de código/veredito alheio ou texto já tachado (não
são alegação minha), 19 são fato de `grep`/`git` já reportado com o número, 10 são verdadeiras no
escopo que a própria frase declara. **Quatro eram minhas e carregavam a classe**:

| Linha | O que dizia | Verificação | Ação |
|---|---|---|---|
| A2 "efeito líquido" | "a **única** asserção de `seller_sku`/`side=erp` da suíte tinha sido removida" | Falso pelo mesmo `TestCase8` — era o universal do `:642` vestido de prosa | **Corrigido no lugar**, com a anterioridade do Case 8 provada por `grep MLB-FX8` |
| A2 conserto | "cobertura **restaurada** como teste separado" | Não foi restauração: a direção nunca deixou de ser coberta | **Corrigido**; agora diz corroboração + o `Detail` que o Case 8 não assere |
| A11 | "sem segundo caminho" | `buildConcordantCandidate` → **2** call sites: `:303` (produção) e `generation_service_test.go:580` (teste, nil de propósito) | **Estreitado** para "call site de produção" |
| A2-R1 | "agora cai **sempre** que o anúncio traz um `seller_sku`" | Só com produto resolvido; sem produto o ramo é `case product == nil` → `INCOMPARABLE`/`erp`, não `UNAVAILABLE` | **Estreitado**, com o `switch` citado |

Duas outras frases minhas foram conferidas e **ficam**: "põe em `SideProvider`, nunca `SideERP`" é
verdade dentro do `case listingValue == ""` que ela cita (`:645-647`); e "`statement_timestamp()`
não é tocado pelo diff" é verdade — `git diff 5441fe18… HEAD -- apps contracts packages | grep -c
statement_timestamp` = **0**.

O que essa varredura ensina, e que eu não teria aprendido consertando sítios: **das quatro, duas
eram a MESMA alegação do código, escrita em português no pack.** Varrer só o código teria deixado a
falsidade viva no documento que o gate lê.

### Reconciliação de contagem da varredura — exigida pelo RULING 1, item 3

A pergunta do hub é a certa: **71 de quantos?** Um número sem população, sem o padrão e sem o
momento não é auditável. As duas contagens, com o buraco que elas revelam.

**População.** Os cinco arquivos do pack que **eu** escrevi — o resto de `dispatches/` é retorno de
worker ou de assento, texto que não é meu para varrer:

| Arquivo | Linhas | P1 como rodado (case-sensitive) | P1 + maiúscula | P1 + maiúscula + sem acento | P2 (tokens extras pt/en) |
|---|---|---|---|---|---|
| `EVIDENCE.md` | 916 | 103 | 116 | 120 | 136 |
| `dispatches/p6-gate-brief.md` | 131 | 16 | 16 | 16 | 19 |
| `dispatches/p6-gate-brief-r2.md` | 136 | 14 | 14 | 14 | 14 |
| `dispatches/p6-gate-brief-r3.md` | 143 | 16 | 16 | 16 | 17 |
| `dispatches/a2-assertion-before-after.md` | 72 | 4 | 4 | 4 | 4 |
| **TOTAL** | **1398** | **153** | **166** | **170** | **190** |

Ocorrências do token (não linhas): **220** em P1, **254** em P2.

**Onde foram os 71.** A varredura da rodada 3 rodou sobre o `EVIDENCE.md` de então, 745 linhas,
antes de eu escrever as seções da rodada 3. Essa região ainda existe e ainda são exatamente 745
linhas: hoje ela dá **65** hits com o mesmo padrão. 71 → 65 porque **os consertos da própria
varredura apagaram seis linhas de hit** — a varredura derruba o próprio número, que é o
comportamento esperado. Da linha 746 ao fim (as seções da rodada 3, que citam verbatim cada
universal ofensivo) vêm mais 38. 65 + 38 = 103, o total de hoje do `EVIDENCE.md` em P1.
Consequência de método: **contagem de varredura só significa alguma coisa carimbada com arquivo,
faixa de linhas, padrão e commit.** A minha estava carimbada com um número solto.

**O padrão cobria o português — mas não cobria MAIÚSCULA.** O padrão registrado já trazia
`nunca|sempre|apenas|todo |nenhum|única|único|…` ao lado de `never|always|only|…`, então a metade
pt estava dentro da população (foi ela que pegou as duas alegações em português). O buraco era
outro e era meu: **rodei `grep -E` sem `-i`**. Diferença medida: **13 linhas** só no `EVIDENCE.md`,
mais 4 por acento e 20 por token que eu não tinha listado.

**Triagem das 48 linhas de delta, e o que mudou por causa delas:**

| Bucket | N | Triagem | Ação |
|---|---|---|---|
| Só maiúscula | 24 | 10 são imperativo dos briefs para o assento ("Only this brief binds you", "NEVER pass `-mod=mod`") — instrução, não alegação sobre o código. 2 citam o universal falso já tachado. 9 são alegação minha verdadeira e verificável (`Nunca houve AGREEMENT`, `Nenhum dos dois tem git`, `Nenhuma mutação foi no-op`). **3 exigiram ação** | ver abaixo |
| Só acento | 4 | Todas prosa de cabeçalho/seção (`não pôde`, `Unico criterio` no yaml) | nenhuma |
| Só token novo | 20 | 13 são `zero` como número contado (`zero hits`, `zero migrações`, `zero CR`) — fato de `grep`, não universal. 5 são `ninguém rodou`, que é honest-unknown declarado como tal. 2 são `exclusividade` já corrigida | nenhuma |

As três que exigiram ação, todas minhas, todas verificadas por STRING antes de mexer:

1. **`EVIDENCE.md:137`** — *"**Qualquer** anúncio cujo seller_sku, EAN e título não casem nada emite
   `seller_sku`/`INCOMPARABLE`/`side=erp`"*. Falso na borda: com `seller_sku` vazio o ramo é
   `case listingValue == "":` (`:645-647`) → `SideProvider`. **Estreitado** para "anúncio que traga
   um `seller_sku` não vazio", com o outro ramo nomeado. É irmã exata da frase A2-R1 que a rodada 3
   já tinha estreitado — escapou porque começa a sentença, com `Q` maiúsculo.
2. **`EVIDENCE.md:460`** — *"**Todo** `PASS (derivado)` do lado A é hand-trace"*. Universal sobre
   nota alheia. **Estreitado** para os **quatro** `PASS (derived)`, com a contagem e o comando que a
   produz, e com o registro de que o round 3 não emitiu nenhum.
3. **`EVIDENCE.md:305` e `:805`** — citação `generation_service_test.go:568` **envelhecida para
   `:580`** pelos meus próprios comentários em `590efdc8`. Reconferida por string e corrigida, com
   os dois call sites, a definição e as duas menções em comentário separados. É a prova de que
   `file:line` de teste dentro do pack apodrece a cada commit de comentário do próprio chip.

**`EVIDENCE.md:335` conferida e MANTIDA:** *"Inalcançável hoje"* para o cenário nil-product da
concordante. A enumeração é total — `buildConcordantCandidate` tem dois call sites, o de produção
(`:303`) passa `&product`, endereço de local; o outro é o guard. A frase se sustenta, e a distinção
com o comentário do teste (que diz "reachability NOT proven", mais fraca) é deliberada: aqui eu
afirmo a enumeração, lá eu não afirmo alcance.

**Pack limpo (RULING 1, item 2):** `git status --porcelain .mnfs/MIS-006-integracao-fundacao/_chip-anchors-3`
sai **vazio** em `de4e940c` e sai vazio de novo no tip congelado deste round — nenhum arquivo do
pack fora do git.

### Ladder depois destes consertos (`0264ba84` + delta de comentários)

Comentários e pack; nenhuma linha de comportamento tocada. Mesmo assim rodado, não alegado:

```
go build ./...                                  → EXIT 0
go vet ./internal/modules/product_links/...     → EXIT 0
go test ./internal/modules/product_links/... -count=1
  connectors 2.820s ok · application 2.796s ok · composition 2.323s ok
  domain 2.361s ok · transport 3.148s ok      → EXIT 0
```

---

## Correções feitas neste EVIDENCE depois do gate

Listadas para que ninguém precise diffar duas versões do documento para saber o que era falso.

| Onde | O que a primeira versão dizia | O que é verdade |
|---|---|---|
| A2 | `side=erp` ficou INALCANÇÁVEL para `seller_sku` | Alcançável em produção pelo caminho unresolved (`generation_service.go:216`). Comentário e cobertura corrigidos em `54342331`. |
| A3 | um subteste discrimina, os outros dois passam nos dois códigos | **Dois** discriminam; `a3-mustfail-raw.txt:18` já mostrava |
| A7 | flake root-caused a skew de ~600ms do container | Diagnóstico não fecha com 4/9 por ~1ms; sai como **não provado** |
| A12 | "Bruto: `ladder-l0-l1-raw.txt`" cobrindo build/vet/count=10 | Aquele arquivo não continha nenhum dos três, e estava corrompido. Recapturado em `ladder-l0-l1-raw-r2.txt` |
| A11 | o guard "degrada como as duas irmãs" | Paridade no **nil-check** é verdade; paridade na **degradação** não. Frase não é falsa, é estreita — duas frases somadas em `2bed7d9d`, texto existente mantido. Ver R5 |
| ledger | reviewer de fatia CONFIRMED contado como corroboração | Ele repetiu meu hand-trace errado e não pegou nenhum blocking; conta menos |
| R4 | `ltrim` nos dois lados abre **super**contagem de `vinculados` | **Rejeitado contra dado real** (10584 linhas, zero colisão / zero zero-à-esquerda). `importados` é `count(*)` sobre linhas cruas, então a chave crua é a consistente. Intuição certa, contador errado. Ver R4 |
| R5 | o comentário do guard nil é falso e eu ia deletá-lo | As duas metades são verdadeiras. Deletar teria sido apagar verdade — o erro espelhado do round 1. Ver R5 |
| R6 | "se um provider **declarar** `seller_sku` com `Supplied: false`" | `Supplied` não é declarado: o adapter marca `Supplied = (anchor ∈ declaredSet)`. E o ramo **dispara hoje**, por `marca`. Ver R6 |
| A11/R5 | `TestConcordantCandidateDoesNotDerefNilProduct` "fixa exatamente" 95/ALTA/ACCEPT | Não fixava: o corpo só asseria `InternalProductID != nil` e um motivo `seller_sku`. Achado pelo round 2. **Consertado pelo lado do teste** em `0264ba84`, não abrandando a frase |
| A2 (2ª vez) | "comentário e cobertura corrigidos em `54342331`" | O mesmo universal falso sobrevivia num **segundo** sítio, `generation_service_test.go:505`. Corrigi o sítio nomeado pelo gate e não varri a classe. Corrigido em `0264ba84` |
| gate | `dispatches/p6-opus-gate-r1.md` citado como artefato | **Nunca existiu.** Assento read-only não tem Write e eu não persisti o retorno. Salvo do transcript em `0264ba84`-time, com cabeçalho dizendo que é recuperação |
| A2 (3ª vez) | "a **única** asserção de `seller_sku`/`side=erp` da suíte tinha sido removida" | `TestCase8NoAnchorResolvedYieldsZeroConfidenceNoCandidateWithReasons:1823-1826` sempre asseriu essa direção, é anterior ao chip (`grep MLB-FX8 p6-input-r4.patch` = 0) e nunca foi tocado. O que caiu foi um caso de tabela, não a cobertura. **O mesmo universal, agora achado na minha prosa pela varredura que eu mesmo rodei** |
| A11 | "sem segundo caminho" para `buildConcordantCandidate` | Há **dois** call sites; o segundo é o teste, que passa nil de propósito. Estreitado para "call site de produção" |
| A2-R1 | `seller_sku` "cai **sempre** que o anúncio traz um `seller_sku`" | Só com produto resolvido. Sem produto o ramo é `case product == nil` → `INCOMPARABLE`/`erp` |
| A11 (delta) | o guard "degrada em CORROBORAÇÃO, as irmãs em ABSENCE" e "nil-checks como as duas irmãs" | `applySingleAnchorScore` ainda emite o próprio motivo FOR a 70/MEDIA/CONFIRM com produto nil, e `applyUnresolvedScore` não checa nada — passa `nil` fixo. Texto meu do round 3, corrigido no round 3 |

**Duas correções desta tabela (R4, R5) foram achadas pelo hub contra dado que eu não tinha, e as
duas me devolveram crédito, não tiraram.** Ficam listadas assim mesmo: uma autodenúncia errada é uma
alegação falsa como qualquer outra.

**Nota de fase:** as seções R4, R5, R6 e a linha A11 acima foram reescritas **durante** o round 2,
depois que os dois assentos de leitura já estavam em voo contra `2bed7d9d`. O código sob revisão não
mudou (árvore limpa fora do `.mnfs`); mudou o texto do pack. Se um veredito do round 2 citar a
redação anterior de R4/R5/R6/A11, é por isso — e a redação anterior era a errada.

---

## FINDINGS de campo, para ratificação (core §0)

1. **A lane de integração não distingue PULADO de VERDE** — ratificado nesta redação estreita
   (profile `ea919c06`), depois que o hub refutou a minha, que era mais larga. O que eu escrevi era
   "a lane não consegue provar que rodou": falso na metade que importa para must-fail. O `summary.txt`
   grava mesmo só `target`/`status`/`run_id`, mas o **stdout da lane** emite em falha
   `failure_token=package=<pacote>` e `failure_token=test=<teste>`, então **falhou-vs-verde é
   distinguível e citável** — foi com esse token que o hub descarregou o must-fail do R4. O que
   continua byte-idêntico é **pulado-vs-verde**: `go test -tags=integration` roda sem `-v` e todo
   teste de DB abre com `SkipWithoutTarget`. Consequência prática, e ela me dá ferramenta a mais:
   must-fail de integração **não precisa** de run direta com `-v` fora da lane — roda a lane com a
   mutação e cita o `failure_token=test=`, que é mais barato e mais auditável do que o que eu fiz em
   A5/A6/A7/A10. Citar `failure_token=package=` sozinho, não: ver o achado colateral do hub em R4, o
   pacote listado pode ser vítima de contaminação por banco compartilhado.
2. **Asserção de bracket temporal contra relógio de container é frágil** e o diagnóstico fácil
   ("skew") não sobrevive à aritmética. Ver A7/R2.
3. **`git commit -F` com `$TMPDIR`** falha silenciosamente neste ambiente — `$TMPDIR` vem vazio no
   Git Bash e o caminho resolve para `C:/Program Files/Git/…`. Use caminho absoluto do scratchpad.
4. **Um gate P6 fisicamente read-only não descarrega critério de execução, por construção.** O lado
   A não tem Bash nem git; o sandbox do lado B negou o diretório de trabalho do Go e o Postgres. A
   doutrina exige must-fail e ladder; nenhum dos dois pode reproduzir nenhum dos dois, então toda
   obrigação de execução acaba certificada apenas pelo próprio chip — exatamente o que o brief manda
   não aceitar. Isso não invalida o round 1 (os três blocking do lado A são de leitura, e são
   reais), mas significa que a metade executável da doutrina hoje **não tem gate**. Candidato a
   emenda: um terceiro verificador com Bash, independente do implementador, ou uma lane de CI cujo
   artefato o gate possa citar.
5. **`go build ./...` verde não implica pacotes de teste linkáveis**, e OOM de host aparece como
   `[build failed]` — indistinguível de erro de compilação numa leitura rápida do resumo. Ver A12.
6. **Mandar o assento read-only escrever o próprio veredito não persiste nada, nos DOIS lados e por
   motivos diferentes.** `harness:gate-reviewer` não tem a ferramenta Write; o sandbox read-only do
   codex recusa o `apply_patch` (`patch rejected: writing is blocked by read-only sandbox`). O round
   1 perdeu o artefato do lado A assim; o round 3 teria perdido o do lado B. Persistir é passo do
   **orquestrador**, na mesma mensagem em que o veredito chega — e o payload da escrita recusada é
   recuperável do rollout do codex, o que vale mais que prosa de transcript. Candidato a emenda do
   §11: o brief não delega persistência ao assento.
7. **A varredura de classe tem que ser rodada pelo AUTOR contra o próprio pack, não só pelo
   revisor.** Das quatro alegações minhas que a varredura desta rodada pegou, **duas eram a mesma
   alegação do código escrita em português no EVIDENCE**. Um gate que varre só o código deixa a
   falsidade viva no documento que o gate lê. Custo: um `grep`. **RATIFICADO** no profile em
   `ea919c06`, na minha redação: o brief de delta manda varrer a CLASSE, e "veredito sem seção de
   SWEEP é incompleto na cara".
8. **O pack de evidência ficou UNTRACKED por seis commits deste chip.** `chip.md` e
   `validation-contract.md` entraram pela mão do hub em `af61c5e8`; `EVIDENCE.md`, `dispatches/` e os
   seis `p6-*.patch` nunca entraram em commit nenhum — `git check-ignore` sai 1 (nenhuma regra os
   ignora), então não foi política, foi omissão minha: commitei só o código a cada slice verde. Por
   seis commits a doutrina "unwritten = didn't happen" estava satisfeita só no sistema de arquivos de
   um worktree descartável, e os dois vereditos do gate — os artefatos que mais custaram a existir,
   um deles recuperado de um `apply_patch` recusado — teriam morrido com ele. Corrigido agora, no
   commit do round 3. Candidato a emenda: o slice verde inclui o pack, não só o código; e uma
   verificação barata (`git status --porcelain .mnfs/<pack>` limpo) cabe no próprio brief do gate,
   porque o assento lê o pack pelo disco e não percebe a diferença.
