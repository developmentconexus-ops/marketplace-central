# CHIP-ANCHORS-2 — rulings do hub

Numeração `A2-Rn`, separada da série `R-1…R-26` do CHIP-ANCHORS, que continua valendo e é citada
por lá.

---

## A2-R1 — GRANT (a), com uma exclusão nomeada

**Pedido:** `UNAVAILABLE` carrega hoje um segundo sentido — "procurei e não achei" — herdado da
base. Reclassificar (a), ou reconhecer o segundo sentido no contrato (b)?

**Ruling: (a) CONCEDIDO.** O grant amplia o write-set para os sítios de semeadura de motivo em
`generation_service.go` nomeados abaixo. **(b) está recusado**, e a razão é doutrina, não gosto: o
contrato congelado afirma que `UNAVAILABLE` tem "o significado atual e apenas ele". Se o código tem
dois, essa frase é FALSA. R-25 — frase falsa se DELETA, não se anota. A opção (b) é anotar uma
falsidade e chamar isso de contrato.

**O achado é mais forte do que o pedido descreve, e é isso que decide.** Em `applySingleAnchorScore`,
o código **já computa a distinção** que D-B quer:

```go
eanDetail := "sem EAN para corroborar o CODPROD"
if strings.TrimSpace(snapshot.EAN) != "" {
    eanDetail = "sem EAN para corroborar o CODPROD: o EAN do anúncio não casa nenhum produto"
}
```

O ramo existe. Ele só está sendo achatado em `UNAVAILABLE` mais uma frase em português — que é
exatamente a alternativa que D-B examinou e rejeitou por escrito ("distinguir só pelo texto de
`detail` obrigaria o FE a fazer parsing de frase em português, que quebra em silêncio quando alguém
reescreve o texto"). Não é escopo novo: é D-B aplicado onde o dado já sabe a resposta. Deixar de
fora entregaria D-B meio aplicado, e meio aplicado justamente nos caminhos que rodam — `ExactSKU` e
`ExactEAN` são a fila de confirmação.

**Mapeamento.** Localize os sítios por STRING, nunca por linha:
`"sem EAN para corroborar o CODPROD"`, `"sem CODPROD para corroborar o EAN"`,
`"seller_sku sem correspondência"`, `"ean sem correspondência"`, e as sementes de
`applyUnresolvedScore`.

| Situação no sítio | `direction` | `side` |
|---|---|---|
| anúncio sem valor da âncora (`snapshot.X == ""`) | `INCOMPARABLE` | `provider` |
| anúncio com valor, e o produto ERP em mãos tem o campo VAZIO | `INCOMPARABLE` | `erp` |
| anúncio com valor e nenhum produto ERP casado (`applyUnresolvedScore`) | `INCOMPARABLE` | `erp` |

O `detail` de cada sítio **fica como está**. A direção dá a classe legível por máquina; o `detail` dá
a especificidade humana. É a mesma divisão de trabalho que D-B escolheu, e reescrever as frases agora
só arriscaria perder informação que já está certa.

**Exclusão nomeada — não toque, e não afirme nada sobre ela.** Existe um ramo que este mapeamento
NÃO cobre: anúncio com valor, e o produto ERP em mãos tem valor **não-vazio e diferente**. Exemplo
alcançável hoje: `seller_sku` resolve o produto A, cujo EAN é `111`; o anúncio traz EAN `999`, que
não casa produto nenhum, então não há conflito CODPROD≠EAN e a execução cai em `ExactSKU` com um
`UNAVAILABLE` para `ean`. A verdade ali é discordância — `AGAINST` — e não incomparabilidade.

Isso **não é deste chip**. Virar `AGAINST` mexe em confiança, banda e status, e portanto na política
D-121, que é decisão do operador. Deixe o ramo exatamente como está, **não o classifique como
`INCOMPARABLE`** (seria trocar uma afirmação errada por outra), e registre-o no EVIDENCE como
FINDING com a evidência por string. O hub leva ao operador.

Um teste tem de provar que o ramo excluído não foi tocado: produto com valor diferente do anúncio
continua produzindo o motivo de hoje. Sem esse teste, a exclusão é uma intenção, não um fato.

---

## A2-R2 — NEGADO o emit; C2 é que está errado, e o hub corrige C2

**Pedido:** âncora declarada, com valor dos dois lados, não emite nada (`emit == false` em
`classifyProviderIdentityAnchor`) e some da tela. Fica fora de escopo (a), ou D-B ganha uma quinta
linha (b)?

**Ruling: (a), e a razão não é escopo — é que não existe valor honesto para emitir.**

O classificador compara **presença**, não **valor**:

```go
if listingValue != "" && productValue != "" {
    return "", "", "", false
}
```

Ele nunca estabeleceu que os dois lados CONCORDAM. Emitir `FOR` ali afirmaria uma corroboração que
não foi verificada — a forma exata do defeito que custou seis rounds ao chip anterior e produziu
R-24: **total na redação, parcial no código**. E o único outro sinal que existe sobre título é
`detectHardNegative`, que reporta contradição e cala na ausência dela; "não achei contradição" não é
"concorda", e vender um pelo outro é a mesma alegação inflada.

Uma quinta direção para "comparado, sem contradição, não decide" está negada por R-24 pelo outro
lado: mais aparato para carregar uma distinção que nenhuma tela foi desenhada para mostrar, num chip
que já tem quatro features.

**A preocupação com "Identificado por" está correta mas mal endereçada.** `title FOR` **já existe na
base**, semeado no estado `TitleMatch` com o detail `"match por título (ranking-only, nunca ACCEPT)"`.
Ou seja: um FE que derive "Identificado por" de `direction == FOR` já está errado hoje, antes deste
chip. Por D-122, "Identificado por" é o conjunto das âncoras que **DECIDIRAM**, e `direction` não diz
isso — `FOR` é opinião, não decisão. Isso é uma restrição de contrato para a onda 2, e o hub a leva
para o pack do CHIP-VINC-NEUTRO. Não é trabalho deste chip.

**O que muda por causa deste achado é o C2, não o código.** O título do critério — *"nenhuma âncora
declarada some calada"* — alega mais do que a tabela dele verifica, e alegação maior que a
verificação é precisamente o que R-24 manda consertar **estreitando a alegação**, não alargando o
código. O hub reescreveu C2 no contrato de validação. Leia a versão nova antes de fechar C2.

Registre a lacuna no EVIDENCE: âncora comparável e presente dos dois lados não aparece em
`reasons[]`, por desenho, porque a comparação de valor não é feita. A onda 2 sabe que `title` pode
faltar da lista.

---

## Nota de método

O desvio declarado — F-01/F-02 saíram por feature inteira, não por slice — está **aceito**. Foi
declarado por quem o cometeu, antes de o hub perguntar, e o ladder foi verificado pelo próprio chip
em vez de auto-reportado por worker. Registre no EVIDENCE como está.

A procedência do achado A2-R2 (reviewer adversarial independente, cego às conclusões do chip) foi
registrada sem ser perguntada. É a prática certa: um achado vale o que vale sua origem, e uma origem
omitida vira crédito do relator.
