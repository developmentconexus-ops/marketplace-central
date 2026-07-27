# Decisões D-122 — âncoras honestas + telas neutras de provider

Entrevista com o operador em 2026-07-27, após o merge do CHIP-ANCHORS (`main @89ddf22e`). Seis
decisões ratificadas. Este arquivo é a autoridade sobre elas; os packs dos chips apontam para cá em
vez de recopiar (R-14 — quem copia um fato derivado passa a ter duas versões dele).

---

## D-A — `refforn` sai do vocabulário cross-side

**Decisão:** `refforn` deixa de ser âncora de identidade cross-side. Sai de `knownIdentityAnchors`
([marketplace_capability.go:32-38](../../apps/server_core/internal/modules/connectors/ports/marketplace_capability.go:32))
e continua existindo apenas como campo do lado ERP.

**Por quê.** `refforn` é a referência do fornecedor dentro do ERP (`ZP1704.1.`). Nenhum marketplace
fornece isso — a pergunta "você fornece refforn?" tem resposta `não` para sempre, em todo provider
presente e futuro. Mantê-la na lista gera, em toda linha de todo vínculo, um motivo
`provider não fornece a âncora refforn` que nunca muda de valor. Ruído honesto continua sendo ruído.

**Se um dia um marketplace expuser referência de fabricante**, ela entra como âncora nova com nome
próprio — não reaproveitando um termo que significa outra coisa do lado ERP.

**Histórico:** motivos já persistidos em `link_candidates.reasons` não são migrados. Candidatos
abertos regeneram no próximo ciclo (o repo faz `reasons = EXCLUDED.reasons`); candidatos já
decididos guardam o motivo que existia no momento da decisão, que é o comportamento correto para um
registro de auditoria.

## D-B — terceiro estado: `INCOMPARABLE` + `side`, no DADO

**Decisão:** `direction` ganha o valor `INCOMPARABLE`, e o motivo ganha o campo
`side ∈ provider | erp | both`.

Hoje `direction` é `FOR | AGAINST | UNAVAILABLE`, e `UNAVAILABLE` carrega um significado só: *o
provider não fornece esta âncora, nunca vai fornecer*. O caso que não tem representação é o outro:
**o provider fornece a âncora, mas o valor faltou neste caso** — anúncio com o campo vazio, ou
produto ERP sem o dado cadastrado. Hoje esse caso não aparece: a âncora some inteira da tela
([generation_service.go:639-642](../../apps/server_core/internal/modules/product_links/application/generation_service.go:639)
só emite motivo quando `Supplied == false`), e o operador não distingue "bateu", "não bateu" e "não
deu para comparar".

**Por que no dado e não só na tela.** As duas situações pedem AÇÕES OPOSTAS do operador: limitação
permanente do provider = não fazer nada; valor faltando = ir cadastrar no ERP ou corrigir o anúncio.
Colapsar as duas no mesmo `direction` obriga qualquer leitor do dado cru — banco, log, suporte — a
conhecer uma regra lateral para saber qual é. Com valor próprio no enum, o `Record<Direction, …>` do
FE deixa de compilar até o caso novo ser tratado, então esquecer passa a ser impossível em vez de
improvável.

Rejeitado: reusar `UNAVAILABLE` com `side` (a tela fica igual, o dado cru fica ambíguo) e distinguir
só pelo texto de `detail` (obrigaria o FE a fazer parsing de frase em português, que quebra em
silêncio quando alguém reescreve o texto).

**`INCOMPARABLE` nunca é evidência contra.** Ausência de dado não é discordância. Não rebaixa
confiança e não muda política de auto-aprovação — a regra D-121 continua intacta (só CODPROD+EAN
concordantes auto-aprovam; âncora única vai para confirmação).

## D-C — display de dimensão mostra a equivalência: `50cm ≡ 500MM`

**Decisão:** quando duas grafias normalizam para a mesma chave, a tela mostra as duas, e a ordenação
vira estável.

[generation_service.go:764-769](../../apps/server_core/internal/modules/product_links/application/generation_service.go:764)
usa `slices.SortFunc` (não estável) seguido de `CompactFunc` por chave. `50cm` e `500MM` colapsam na
mesma chave `500mm` com `display` diferente, e **qual dos dois sobrevive é indeterminado** — a mesma
comparação pode mostrar um numa geração e o outro na seguinte, sem nada ter mudado. A decisão é
idêntica nos dois casos; só o texto dança. Operador que vê texto mudar sozinho para de confiar na
tela, e aí perde-se o valor da tela inteira.

Mostrar as duas mata o nondeterminismo e ganha explicabilidade no mesmo movimento: o operador vê a
conversão de unidade que o sistema aplicou em vez de ter que confiar nela. Foi exatamente a
conversão que F-02 do CHIP-ANCHORS entregou e que o U2 provou live
(`MLB4735326915` × `sku:33698`).

Guard obrigatório: `-count=10` no pacote, porque um guard de ORDEM que roda uma vez só não distingue
estável de sortudo.

## D-D — chain-read é backend, e fecha o gap de decomposição do M-06

**Decisão:** o endpoint de chain-read entra no chip de backend.

M-06 F-01 promete `N importados → N vinculados → N enfileirados` vindos de um join real no servidor
(`sync_state` + `product_links` por protocolo), e M-06 é FE+SDK puro. **Nenhuma milestone anterior
implementa esse endpoint** — a própria tabela de risco do M-06 registra isso com dono = hub. O gap
existia desde o planejamento e é fechado aqui, não descoberto no meio da implementação.

## D-E — M-06 fatia por TELA, não por camada

**Decisão:** dois chips de FE com superfícies disjuntas.

- **`/vinculos`** — badge auto-aprovado (F-04), vocabulário neutro de provider (F-05), render do
  estado `INCOMPARABLE`.
- **`/importacoes` + `/integracoes`** — rota nova, `ImportacaoSection` promovido, cadeia real,
  `ActiveSourceCard` do banco (F-01, F-02).

Com F-05 e o render do terceiro estado, o M-06 chegaria a seis features num chip só. Escopo grande
foi o que produziu seis rounds de gate no CHIP-ANCHORS. O corte por tela dá write-sets naturalmente
disjuntos — os dois chips não compartilham um único arquivo — enquanto o corte por camada
serializaria tudo e fecharia o chip de SDK sem QA de tela nenhum.

## D-F — sequenciamento: backend é onda 1; os dois FE são onda 2, paralelos entre si

**Decisão do operador foi "paralelo".** Ao detalhar, o paralelo real é entre os dois chips de FE, e
o backend precede. Registro a correção com a razão, porque a razão é técnica e tem resposta única:

Ambos os chips de FE consomem tipos que só existem depois do backend — `INCOMPARABLE`/`side` para
`/vinculos`, o endpoint de chain-read para `/importacoes`. Um chip de FE despachado antes disso não
falha no gate, **falha no L0**: não compila. E o P7 de tela precisa de dado real, que só existe com
o backend no ar.

Então: **onda 1 = chip backend** (dono do contrato: comportamento + OpenAPI + SDK no mesmo commit,
profile §7). **Onda 2 = os dois chips de FE em paralelo**, ambos consumindo, nenhum editando o
contrato.

O paralelismo perdido é uma onda. O paralelismo forçado custaria um contrato congelado por
antecipação e um round de retrabalho em dois chips se ele estivesse errado.

---

## Contrato congelado (o que a onda 2 pode assumir)

Motivo de vínculo, forma final:

```json
{
  "anchor":    "marca",
  "direction": "FOR" | "AGAINST" | "UNAVAILABLE" | "INCOMPARABLE",
  "side":      "provider" | "erp" | "both",
  "detail":    "produto ERP sem marca cadastrada"
}
```

- `side` só é emitido quando `direction = INCOMPARABLE`; ausente nos demais.
- `UNAVAILABLE` mantém o significado atual e apenas ele: **o provider não fornece esta âncora**.
- `INCOMPARABLE` significa: **o provider fornece, o valor faltou neste caso**, e `side` diz de que
  lado.

Coluna "Identificado por" (`/vinculos`, substitui "SKU ML"/"GTIN"):

- mostra **todas as âncoras que decidiram**, unidas por ` + ` — ex. `CODPROD + EAN`;
- é diferente da coluna Motivo, que mostra tudo o que opinou. Uma âncora `UNAVAILABLE` ou
  `INCOMPARABLE` aparece em Motivo e **não** aparece em Identificado por;
- a corroboração dupla é justamente o que separa auto-aprovado de mandado para confirmação
  (D-121), então esconder a segunda âncora apagaria da tela a explicação do que a tela mostra.
