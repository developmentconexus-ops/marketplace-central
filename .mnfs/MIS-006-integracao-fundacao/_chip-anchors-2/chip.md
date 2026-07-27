# CHIP-ANCHORS-2 — âncoras honestas + leitura da cadeia (backend)

```yaml
chip: CHIP-ANCHORS-2
mission: MIS-006-integracao-fundacao
branch: chip/anchors-2
base_sha: c7f1c2e90371ddaeb9ca55f893d3dd24fd8d037a
wave: 1 (única — os dois chips de FE são a onda 2 e dependem deste)
authority: .mnfs/MIS-006-integracao-fundacao/DECISOES-D122-anchors-telas.md
```

Este chip é o **dono do contrato** desta onda. Comportamento, OpenAPI e SDK saem no MESMO commit
(profile §7). Os dois chips de FE da onda 2 consomem o que sai daqui e não editam nada disto.

As decisões D-A, D-B, D-C e D-D **não estão recopiadas aqui**. A autoridade sobre elas é
[`DECISOES-D122-anchors-telas.md`](../DECISOES-D122-anchors-telas.md), com o motivo de cada uma e o
contrato congelado no fim do arquivo. Leia esse arquivo INTEIRO antes de escrever a primeira linha —
quem copia um fato derivado passa a ter duas versões dele (R-14). Este pack diz apenas o que o chip
tem de FAZER, onde, e o que prova que fez.

---

## R1 do CHIP-ANCHORS está REVERTIDA por D-B. Leia antes de F-02.

O chip anterior recebeu uma ruling explícita, R1, que dizia o oposto do que D-B agora manda:

> o wire não muda de forma. `reasons[]` continua `{anchor, direction, detail}` e `direction`
> continua o enum `FOR|AGAINST|UNAVAILABLE` (…) A distinção vai no `detail` (texto), **nunca** num
> quarto valor de `direction`.

Isso não é o hub mudando de ideia por gosto. R1 protegia uma colisão concreta: o contrato do
chain-read estava previsto para o M-06 e um quarto valor no enum congelaria a superfície de outro
dono no meio do voo. **D-D dá o chain-read a ESTE chip**, então a colisão que R1 protegia não
existe mais — este chip é o único dono das duas superfícies ao mesmo tempo.

E a alternativa que R1 prescrevia — pôr a distinção no `detail` — é exatamente a opção que o
operador rejeitou em D-B, com a razão registrada lá: obrigaria o FE a fazer parsing de frase em
português, que quebra em silêncio quando alguém reescreve o texto.

Registro isto por escrito porque um revisor que conheça o pack anterior vai encontrar R1 e ler o
`INCOMPARABLE` como violação de ruling. Não é. R1 está superada por D-B, com esta razão, nesta data.

---

## Escopo — quatro features, uma por decisão

### F-01 — `refforn` sai do vocabulário cross-side (D-A)

Remove `IdentityAnchorRefforn` de `knownIdentityAnchors` em
[`marketplace_capability.go:22-42`](../../../apps/server_core/internal/modules/connectors/ports/marketplace_capability.go:22).

O que NÃO muda: `refforn` continua existindo como campo do lado ERP (`erp_import_products.refforn`,
`products_mirror`), continua sendo importado e continua aparecendo onde já aparece. Sai **apenas** da
lista de âncoras que o gerador pergunta ao provider.

Nada de migração de dados. Motivos já persistidos em `link_candidates.reasons` ficam como estão; o
`reasons = EXCLUDED.reasons` do repo regenera os candidatos abertos no próximo ciclo, e candidatos
já decididos guardam o motivo que existia no momento da decisão — que é o comportamento correto de
um registro de auditoria, não um bug a corrigir.

### F-02 — `INCOMPARABLE` + `side`, no dado (D-B)

Três camadas, um commit.

**Domínio.** `LinkCandidateReasonDirection` em
[`link_candidate.go:52-65`](../../../apps/server_core/internal/modules/product_links/domain/link_candidate.go:52)
ganha `LinkCandidateReasonDirectionIncomparable = "INCOMPARABLE"`, e `LinkCandidateReason` ganha
`Side` com `json:"side,omitempty"`. `side ∈ provider | erp | both` — tipo próprio, não `string`
solta, porque um enum de três valores digitado à mão em dois lugares vira quatro valores em seis
meses.

**Geração.** [`generation_service.go:639-642`](../../../apps/server_core/internal/modules/product_links/application/generation_service.go:639)
hoje sai calado quando a âncora é declarada mas não tem sinal:

```go
for _, anchor := range identityAnchors {
    if anchor.Supplied || hasSignal[anchor.Anchor] {
        continue          // <- a âncora some inteira da tela
    }
```

O `continue` cobre dois casos diferentes com o mesmo silêncio. Separe-os:

| Situação | `direction` | `side` |
|---|---|---|
| provider não declara a âncora (`Supplied == false`) | `UNAVAILABLE` | ausente |
| provider declara, valor do anúncio vazio | `INCOMPARABLE` | `provider` |
| provider declara e tem valor, produto ERP sem o dado | `INCOMPARABLE` | `erp` |
| declarada, faltando dos dois lados | `INCOMPARABLE` | `both` |

`UNAVAILABLE` mantém um significado só: **o provider não fornece esta âncora**. `side` só é emitido
quando `direction = INCOMPARABLE`.

`INCOMPARABLE` **nunca é evidência contra**. Ausência de dado não é discordância: não rebaixa
confiança, não entra em nenhum somatório de score e não toca a política de auto-aprovação. A regra
D-121 fica intacta — só CODPROD+EAN concordantes auto-aprovam, âncora única vai para confirmação.
Um teste tem de provar isso explicitamente (C4), porque é o jeito mais fácil de este chip quebrar
algo que já funciona.

**Contrato.** `ProductLinkCandidateReason` em
`contracts/api/marketplace-central.openapi.yaml:6514-6526` (o `enum: [FOR, AGAINST, UNAVAILABLE]`
está em `:6525`) ganha o quarto valor e a propriedade `side`, e
`packages/sdk-runtime/src/index.ts:1060-1064` acompanha:

```ts
export type ProductLinkReasonDirection = "FOR" | "AGAINST" | "UNAVAILABLE" | "INCOMPARABLE";
export type ProductLinkReasonSide = "provider" | "erp" | "both";
```

Isso é o que faz a onda 2 compilar. `apps/web/src/pages/vinculos/QueueRow.tsx` tem três
`Record<ProductLinkReasonDirection, …>` (`:34`, `:75`, e o `byDirection` em `:157`) — quando o tipo
ganhar o quarto valor, o `tsc` de `apps/web` **vai ficar vermelho**. Isso é o mecanismo funcionando,
não um defeito deste chip: é a garantia de que esquecer o caso novo na tela passa a ser impossível
em vez de improvável. **Este chip não edita `apps/web/`.** Deixe vermelho e registre no EVIDENCE
que deixou, dizendo qual chip da onda 2 fecha (CHIP-VINC-NEUTRO). Se você "consertar" a tela aqui,
quebra a disjunção de write-set da onda 2 e o hub rejeita o merge.

### F-03 — display de dimensão mostra a equivalência, ordenação estável (D-C)

[`generation_service.go:764-769`](../../../apps/server_core/internal/modules/product_links/application/generation_service.go:764):

```go
slices.SortFunc(pairs, func(a, b dimensionPair) int {   // NÃO estável
    return strings.Compare(a.key, b.key)
})
pairs = slices.CompactFunc(pairs, func(a, b dimensionPair) bool {
    return a.key == b.key                                // descarta o display do perdedor
})
```

`50cm` e `500MM` colapsam na mesma chave `500mm` com `display` diferente, e qual dos dois sobrevive
é **indeterminado**. Duas correções, e as duas são necessárias:

1. `SortStableFunc` no lugar de `SortFunc`.
2. Ao compactar, não descarte o display do perdedor: quando duas grafias distintas caem na mesma
   chave, o texto exibido passa a mostrar as duas — `50cm ≡ 500MM`. Grafias idênticas continuam
   aparecendo uma vez só (`50cm` e `50cm` é `50cm`, não `50cm ≡ 50cm`). A ordem dos dois lados do
   `≡` é a ordem de entrada, que agora é determinística porque o sort é estável.

Guard obrigatório em `-count=10`. Um guard de ORDEM que roda uma vez não distingue estável de
sortudo — e `-count=1` passando é exatamente o que deixaria este bug entrar de novo. O guard tem de
ser **provado load-bearing**: reverta o `SortStableFunc` para `SortFunc`, mostre o guard falhando
dentro dos 10, restaure. Cole a saída das duas execuções no EVIDENCE (R5). Um guard que passa nas
duas versões não guarda nada.

### F-04 — endpoint de leitura da cadeia (D-D)

`GET /erp/imports/{id}/chain` — o `{id}` é o mesmo do `GET /erp/imports/{id}` que já existe
(`http_handler.go:42`).

Responde os três números que o M-06 F-01 promete, para um protocolo:

```json
{
  "protocol":    "#003-E",
  "importados":  1845,
  "vinculados":  31,
  "enfileirados": 1802,
  "queue_read_at": "2026-07-27T12:00:00Z"
}
```

As três contagens, com as tabelas que já existem — **este chip não cria tabela nem migration**:

- **importados** — `count(*)` de `erp_import_products` em `(tenant_id, protocol_id)`. É a linha do
  produto aceito, não o `accepted_count` denormalizado do protocolo: contar a linha que existe é
  mais honesto que ler um contador que alguém escreveu.
- **vinculados** — `count(DISTINCT eip.codprod)` juntando `erp_import_products` com `product_links`
  por `pl.internal_product_id::text = eip.codprod`, filtrando
  `pl.state = 'resolved'` (`ProductLinkStateResolved`, `product_link.go:12`). `internal_product_id`
  é `integer` e `codprod` é `text` — o cast é obrigatório e vai no JOIN, **nunca** num `ORDER BY`
  sem alias (armadilha conhecida: um `::text` sem alias no SELECT sequestra o `ORDER BY` bare).
- **enfileirados** — interseção dos `codprod` do protocolo com a fila **atual**:
  `sync_state.cursor -> 'pending'`, `entity = 'market'`, via
  `jsonb_array_elements_text`, `DISTINCT` e agregando todas as `installation_id` (a fila é por
  instalação, o protocolo não é).

**A semântica de `enfileirados` tem de aparecer no contrato, não só no código.** A fila drena: o
scheduler consome `pending` e o número CAI com o tempo, sem nada ter dado errado. Então este campo
é *quantos produtos deste protocolo ainda estão na fila agora*, não *quantos foram enfileirados na
importação*. Escreva isso no `description` do OpenAPI e devolva `queue_read_at` junto, para que a
tela possa dizer "na fila agora" com uma hora ao lado. Um número que encolhe sozinho e se chama
"enfileirados" é lido como perda pelo operador — e aí a tela inteira perde o crédito.

Se o protocolo não existe: `404`. Se existe e nunca teve fila (nenhuma linha em `sync_state`):
`enfileirados: 0`, que é a verdade, não `null` — a pergunta "quantos estão na fila" tem resposta
conhecida e é zero. Isso é diferente de não saber (ADR-17); não há caso de não-saber aqui.

O endpoint entra no mesmo commit de OpenAPI + SDK que F-02 (profile §7). Método SDK:
`getErpImportChain(id)` em `packages/sdk-runtime/src/erpImport.ts`.

---

## Propriedade — matriz de seis eixos

| Eixo | Deste chip | Proibido |
|---|---|---|
| Arquivos Go | `connectors/ports/marketplace_capability.go`, `product_links/domain/link_candidate.go`, `product_links/application/generation_service.go`, `erp_import/**` (transport + application + adapters do chain-read) | qualquer outro módulo |
| OpenAPI | `ProductLinkCandidateReason`, path novo `/erp/imports/{id}/chain` | qualquer outra seção |
| SDK | `packages/sdk-runtime/src/index.ts` (tipos de reason), `src/erpImport.ts` (método novo) | resto do pacote |
| Migrations | **nenhuma.** Este chip não cria migration. Se você achar que precisa de uma, pare e mande `REQUEST` ao hub | todas |
| FE | **nenhum.** `apps/web/` está fora de escopo, inclusive para consertar o `tsc` vermelho que F-02 causa | `apps/web/**` |
| Infra / stack | nada | `docker-compose*`, `.env*`, portas |

`packages/sdk-runtime` e o OpenAPI são exclusivos desta onda: nenhum outro chip está no ar.

---

## Ladder e regras de execução

L0: `go build ./...`, `go vet ./...`, lane de governança. L1: `go test ./...` mais o guard de F-03
em `-count=10`. `apps/web` `tsc` fica vermelho de propósito por F-02 — declare no EVIDENCE, não
conserte.

`GOFLAGS` **não pode** estar setado como `-mod=mod` (o repo está em workspace mode; a mensagem é
`go: -mod may only be set to readonly or vendor when in workspace mode`). Rode `go` de dentro de
`apps/server_core`, nunca da raiz do worktree — só `apps/server_core/.gomodcache/` está no
gitignore, e um `go` rodado da raiz polui o `git add -A`. Confira `gomodcache = 0` no `git status`
antes de commitar.

O chip **nunca sobe servidor**, nunca liga em `:8080`, nunca lê `.env*`. Precisa de stack de pé ou
de dado real → `REQUEST` ao hub. O hub roda o live-drive.

Nada de push. O merge é do hub.

## P6 — portão duplo

Você **não** escreve a linha `P6-DUAL-GATE:`. Esse marcador é autoridade do hub, e um chip que o
carimba está atestando o próprio trabalho com a assinatura de outro. Feche com
`AGREEMENT — P6 discharged` e o ledger de discharge ao lado, como o CHIP-ANCHORS fez.

O que o EVIDENCE tem de carregar, e que custou seis rounds ao chip anterior:

- **Cite por STRING, não por linha.** Uma coordenada é uma alegação como outra qualquer, e um
  corretivo que insere linhas acima dela deixa a citação *legível como prova* apontando para código
  sem relação.
- **R-24.** Um artefato de verificação guarda só a alegação que consegue fazer TOTALMENTE. Se a
  redação diz "todos" e o código reconhece "alguns", a saída é DELETAR a alegação, não alargar o
  reconhecedor. Alargar foi o que produziu os rounds 3–6.
- **R-25.** Frase falsa se DELETA. Anotar uma frase falsa como falsa é desculpá-la, e quem encontra
  a falsa primeiro não tem sinal nenhum para continuar lendo. Honest-unknown é para LACUNA.
- **R-26.** VERBATIM é alegação sobre FORMA. Se reproduzir fielmente não cabe no container que você
  escolheu (célula de tabela, por exemplo), troque o CONTAINER — nunca o texto. Comentário dentro de
  uma cerca de citação é conteúdo autoral vestido de citação.

Critérios C1–C10 e as condições de merge: [`validation-contract.md`](validation-contract.md).
