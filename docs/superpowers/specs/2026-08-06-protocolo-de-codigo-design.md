# Protocolo de código — desenho

**Data:** 2026-08-06 · **Estado:** PROPOSTO, aguarda revisão do operador
**Substitui:** as Ondas 1–4 de `docs/superpowers/plans/2026-08-05-arquitetura-protocolo-de-modulo-plan.md`
**Emenda:** ADR-023 (protocolo de módulo) — ver §14

---

## 0. Porquê este documento

A ADR-023 já tentou ser este documento e não pegou. Não pegou porque as regras dela viviam
em prosa e o instrumento vinha depois. Medido: **128 violações da regra de fronteira** em
código de produção, mais **100 em ficheiros de teste** que o instrumento nem olhava, mais
**140 imports** da raiz de composição sob um carve-out declarado. A ADR dizia "35".

O catálogo de classes de defeito da casa nomeia a causa: *o sinal é mais barato de fabricar
que o trabalho*. Uma regra que precisa de um humano ou de um teste verde para ser cumprida
é uma regra que não se cumpre.

Este documento corrige isso com uma única disciplina, aplicada a toda a regra sem exceção:

> **Toda regra leva um nível de imposição. Regra que não sobe de nível não entra.**

### Os três níveis

| Nível | Como é imposta | Como se falha |
|---|---|---|
| **1 — impossível** | O tipo ou o compilador recusa. | O código não compila. |
| **2 — verificado** | Uma máquina refaz o trabalho e compara, ou analisa o AST. Corre no `go vet` ou na lane. | Vermelho determinístico, sem baseline. |
| **3 — convenção** | Só revisão humana. | Não falha. |

**Nível 3 não é regra, é intenção.** Ou sobe para 1 ou 2, ou sai deste documento. Este
critério é o que separa este protocolo da ADR-023.

---

## 1. A forma

### 1.1 Contexto

Um contexto é uma pasta sob `apps/server_core/internal/contexts/`:

```
contexts/<nome>/
  contracts/                   # o que este contexto PUBLICA
  port/                        # o que este contexto EXIGE dos outros
  internal/
    domain/                    # modelo e invariantes
    application/               # casos de uso
    postgres/                  # persistência
  module.go                    # construtor único, consumido pela raiz
```

`contracts` e `port` são públicos. Tudo sob `internal/` é privado ao contexto, imposto pela
regra de subpacote `internal` do toolchain Go.

**Regra 1.1 — Nível 1.** Um contexto não importa `internal/` de outro contexto.
*Instrumento:* a regra de prefixo de import do Go.
*Medido:* provado contra um esqueleto de 13 contextos. `pricing/internal/application` a
importar `catalog/internal/domain` produz
`use of internal package .../catalog/internal/domain not allowed`. Removida a linha, compila.
*Defeito que teria apanhado:* as 128 violações, incluindo
`listings/adapters/connectors/backfill.go:9` a ler `connectors/adapters/mercado_livre.ItemMultigetDTO`.

**Regra 1.2 — Nível 1.** A separação entre publicar e exigir é estrutural: `contracts` é o
que oferecemos, `port` é o que precisamos. Um não pode conter o outro porque são pacotes
distintos.
*Defeito que teria apanhado:* `connectors/ports/` de hoje mistura capacidade de provider,
repositório e consumível inter-contexto no mesmo pacote — três coisas com donos e direções
diferentes.

**Regra 1.3 — Nível 2.** `contracts` não importa `port`, e nenhum dos dois importa
`internal/`.
*Instrumento:* analisador `go/analysis`.

**Emenda 2026-08-07** (paga por `apps/server_core/internal/arch/scan.go:131-133`): o
instrumento é o detector de cross-context (`ScanCrossContextInternal`) **sem** o salto de
ficheiros que vivem fora de `contexts/`. O comentário no sítio da medição:

> "here is "" for a file that lives outside any context — the composition root, an
> adapter, a cmd. Those are NOT skipped: the one import that ever broke this rule came
> from exactly there, and a detector that starts by skipping them can only ever report
> zero."

O único sítio que violou esta classe de regra estava fora de qualquer `contexts/<nome>/`
(a raiz de composição), e um detector que pula ficheiros "fora de contexto" por
construção fica cego a exatamente esse caso — fica sempre verde, não porque não haja
violação, mas porque nunca olhou para onde ela vive.

### 1.2 Os treze contextos

| Contexto | Possui |
|---|---|
| `account` | Tenants, utilizadores, contas de canal, ligações ao ERP, referências de credencial |
| `catalog` | Identidade canónica de produto, variantes, atributos, observações de fonte, importações |
| `listings` | Anúncios observados no canal, variações, estado reportado pelo marketplace |
| `linking` | Candidatos e vínculos produto↔variação, com evidência e histórico de decisão |
| `inventory` | Observações de stock do ERP, elegibilidade, buffers, stock vendável, disponibilidade desejada |
| `costing` | Custo com vigência e método — médio, reposição, padrão, por empresa e local |
| `tax` | Classificação fiscal, matrizes com vigência, fonte legal, regras e determinações |
| `pricing` | Cenários, resolução de margem alvo, comissão e frete previstos, políticas, recomendações |
| `orders` | Pedidos como observados: itens, pagamentos, taxas, envios, ciclo de vida |
| `reconciliation` | Documentos do ERP e o seu emparelhamento com pedidos e itens; divergências |
| `profitability` | Cálculo de lucro realizado versionado, e os factos de input exatos de cada cálculo |
| `intelligence` | Observações de concorrência, confiança de emparelhamento, histórico, sinais |
| `changecontrol` | Preview, aprovação, execução, retentativa, verificação e auditoria de toda a escrita |

Três destes são novos e nenhum é invenção: `costing`, `tax` e `reconciliation` existem hoje
como matéria dispersa, e é exatamente onde temos quatro fórmulas de margem, três cópias da
tabela de ICMS e a reconciliação Sankhya enterrada dentro de `orders`.

Oito módulos atuais deixam de ser módulos:

| Módulo atual | Destino |
|---|---|
| `connectors` | `adapters/marketplace/{mercadolivre,shopee,amazon}/` |
| `internal_read` | reparte-se por `catalog`, `inventory`, `costing`, `tax`, `adapters/erp` |
| `sync` | `platform/scheduler` + cursores por contexto |
| `dashboard` | `views/` |
| `sourcekind` | `kernel/provenance`, decomposto |
| `erp_import` | `adapters/erp/sankhyaoracle` + ingestão em `catalog` |
| `marketplaces` | `account` + `pricing` |
| `divergences`, `channelfees`, `classifications`, `tenant_config` | dissolvem-se nos donos |

### 1.3 Direção de dependência

**Regra 1.4 — Nível 1.** Não há ciclo de pacote.
*Instrumento:* o compilador Go. `go build ./...` falha com ciclo.
*Medido:* os 13 contextos compõem sem ciclo no esqueleto de referência.
*Defeito que teria apanhado:* `internal_read ↔ erp_import`, 7 imports em cada sentido — a
maior aresta do grafo, nos dois sentidos. E `catalog ↔ internal_read`.

**Regra 1.5 — Nível 2.** As arestas permitidas entre contextos são declaradas e verificadas.
Aresta não declarada é vermelho.
*Instrumento:* analisador sobre os imports de `contracts` e `port`.

---

## 2. Adapters e sistemas externos

```
adapters/
  marketplace/<vendor>/
    <vendor>.go        # ÚNICA superfície importável: New() Bundle
    internal/api/      # HTTP, auth, paginação, DTOs de fio, erros crus
    listings/          # implementa contexts/listings/port
    orders/            # implementa contexts/orders/port
    pricing/           # implementa as portas de cobrança de pricing
    writes/            # implementa contexts/changecontrol/port
  erp/sankhyaoracle/
    oracle/       # SQL, tipos Oracle, conexão
    catalogfeed/  inventoryfeed/  costingfeed/  taxfeed/  documents/
  spreadsheet/catalogimport/
```

**Regra 2.1 — Nível 1.** A porta pertence ao contexto que consome, nunca ao adapter. O
adapter implementa; não define.
*Instrumento:* a porta vive em `contexts/<n>/port/`, que o adapter importa. A inversão é
estrutural.

**Regra 2.2 — Nível 1.** O DTO de fio vive em `adapters/marketplace/<v>/internal/api` e é
inalcançável de fora da árvore de `<v>`: outro vendor, um contexto, a raiz de composição.
*Instrumento:* a regra de subpacote `internal` do Go, aplicada à raiz do vendor e não à raiz
do módulo.
*Medido, 2026-08-06, contra o esqueleto de 13 contextos:* com `api/` no sítio antigo,
`adapters/marketplace/shopee` importa `adapters/marketplace/mercadolivre/api` e **compila
limpo** — o irmão passa. Movido para `mercadolivre/internal/api`, o mesmo import produz
`use of internal package .../mercadolivre/internal/api not allowed`. Removido, `go build ./...`
e `go vet ./...` voltam a exit 0.
*Defeito que teria apanhado:* os **70 tokens de vendor fora de adapters contra 54 dentro**.
Há mais conhecimento de Mercado Livre fora dos adapters do que dentro.

**Regra 2.2-a — Nível 1.** **Regra do molde, não regra de adaptador:** vale para qualquer
árvore que contenha um `internal/` — vendor sob `adapters/marketplace/<v>` ou contexto sob
`contexts/<nome>`, indistintamente. Essa árvore expõe um único pacote raiz, com um
construtor exportado cujos parâmetros são só tipos que um estranho pode nomear — um pool,
uma DSN, um `*sql.DB`, as portas dos contextos consumidores — nunca um tipo declarado sob
o `internal/` da própria árvore. Tudo o resto é montado lá dentro.
*Porquê é forçado e não escolhido:* movido o `api` para `internal/`, a mesma medição
mostrou `bootstrap/modules.go:8` a ser rejeitado junto com o irmão — a raiz construía
`mlapi.StubClient{}` e injetava-o em quatro adapters de capacidade. E os construtores desses
adapters têm assinatura `New(client api.Client)`: **um construtor exportado cujo parâmetro é
um tipo interno é inchamável de fora da árvore.** O compilador não permite outra topologia
senão a fachada. Com a fachada, o build volta a exit 0 e a raiz deixa de nomear um único tipo
de vendor.
*Defeito que teria apanhado:* exatamente as 10 declarações de adapter de ML na nossa raiz
(§2.5) — é a mesma doença, reproduzida em miniatura e curada pelo compilador.

**Emenda 2026-08-07** (paga por `apps/server_core/internal/contexts/catalog/module.go:23-38`):
a regra estava escrita só para `adapters/marketplace/<v>`, e escrita assim não alcançou
`contexts/`. Medido: `catalog.New(store, ids, reader)`, com `store` e `reader` tipados por
`catalog/internal/application`/`catalog/internal/postgres`, nasceu com **zero chamadores
legais** — a raiz de composição só o conseguia satisfazer importando o pacote interno, e o
compilador recusou. A raiz não estava errada; a assinatura é que era. O construtor atual —

```go
// New assembles the context from the ONLY thing an outsider may legitimately
// name: a connection pool.
func New(pool *pgxpool.Pool) *Module {
	repo := postgres.NewRepository(pool)
	return &Module{
		service: application.NewService(repo, postgres.NewULIDFactory()),
		reader:  postgres.NewSummaryReader(repo),
	}
}
```
(`module.go:32-38`) — é a fachada de raiz única aplicada a um contexto: um parâmetro
(`*pgxpool.Pool`), tipos internos montados dentro. A regra existia; o escopo é que era
estreito demais para a apanhar.

**Regra 2.3 — Nível 2.** Nome de vendor não aparece em `contexts/`. Código de canal é dado
em runtime, nunca `enum` fechado em Go nem literal comparada.
*Instrumento:* analisador, lista de tokens de vendor.
*Defeito que teria apanhado:* `orders/application/ingest_service.go:334` e
`assisted_sankhya_linkage_service.go:296` comparam a string `"mercado_livre"` na camada de
aplicação. Adicionar Shopee toca ali.

**Emenda 2026-08-07** (paga por `apps/server_core/internal/arch/scan.go:34-37` vs
`scan.go:217-226`; a citação original desta emenda, `:33-36` vs `:203`, apontava para as
linhas certas antes das dez tarefas anteriores deste plano moverem código — recontada
contra o ficheiro atual): o escopo (`adapters/` isento, `internal/arch/` isento) vive **no
próprio detector**, na função `vendorRuleApplies` (`scan.go:217-226`), não numa convenção à
parte:

> "The rule is 'a vendor name does not appear OUTSIDE adapters/', so adapters are exempt
> by definition. The scanner package is exempt too: its closed token list is the
> instrument, not a violation, and a detector that permanently accuses itself is a
> detector nobody can ever act on."

Sem esse filtro de caminho, `ScanVendorTokens` (que percorre `.go` por identificador e
literal de string, `scan.go:228-275`) acusaria a sua própria `VendorTokens`
(`scan.go:34-37`, os literais `"mercado_livre"`, `"shopee"` etc.) e a sua própria
implementação nos adapters, e nunca ficaria verde — a lista de tokens é o instrumento, não
uma violação, e um detector que se acusa permanentemente é um detector que ninguém
consegue agir sobre.

**Regra 2.4 — Nível 1.** Não existe interface `Marketplace` única. Cada capacidade é uma
porta pequena e independente.
*Instrumento:* não há tipo onde escrever a interface grande.
*Porquê:* com uma interface única, a semântica do primeiro provider define o "genérico" e o
segundo entra com campos condicionais. É como `connectors/domain/capability.go` chegou a ser
um modelo de comércio universal que é o Mercado Livre com nomes neutros.

**Regra 2.5 — Nível 1.** A raiz de composição não declara adapters.
*Instrumento:* sob `internal/`, a raiz não consegue nomear os tipos internos necessários para
os construir.
*Defeito que teria apanhado:* **10 das 16 declarações de tipo da raiz** são adapters de
Mercado Livre — `market_adapters.go:228`, `orders_ingest_adapters.go:25` e `:59`,
`orders_adapters.go:88`, `pricing_adapters.go:50` e `:115`. A raiz é, de facto, o adapter do
Mercado Livre.

---

## 3. Kernel

```
kernel/
  tenant/      # tenant.ID
  channel/     # Code, AccountRef
  exact/       # Decimal, Money
  fact/        # Knowledge, Fact[T]
  provenance/  # Evidence
  period/      # EffectivePeriod
```

**Regra 3.1 — Nível 2.** O kernel tem seis membros. Membro novo exige emenda a este
documento com justificação escrita.
*Critério de admissão:* mesmo significado e mesmas invariantes em pelo menos três contextos;
objeto de valor imutável; sem dependência de vendor ou de contexto; dono nomeado.
*Explicitamente fora:* `Product`, `Listing`, `Order`, `Address`, `Warehouse`, `TaxRate`, e
qualquer `Source` genérico. São palavras com modelos diferentes em contextos diferentes.

**Regra 3.2 — Nível 1.** Dinheiro é `exact.Money`. Não existe construtor a partir de
`float64`.
*Instrumento:* o construtor não aceita `float64`; o campo é privado.
*Defeito que teria apanhado:* **quatro structs `Money` independentes** —
`listings/domain/read_model.go:83`, `connectors/domain/money.go:10`,
`pricing/domain/decimal.go:12`, `market/domain/market.go:13` — mais um gémeo
`connectors/domain/capability.go:200`, mais `erp_import/domain/import.go:6` que é `string`.
E **51 campos monetários em `float64`** num sistema fiscal.

**Regra 3.3 — Nível 2.** `float64` proibido em qualquer campo de `contracts` ou de tabela.
*Instrumento:* analisador.

---

## 4. Factos, e o fim do zero silencioso

```go
type Knowledge uint8

const (
    Unknown Knowledge = iota
    Known
    Estimated
    NotApplicable
)

type Fact[T any] struct {
    state    Knowledge
    value    *T
    reason   string
    evidence provenance.Evidence
}
```

**Regra 4.1 — Nível 1.** Combinação inválida é inexprimível. O construtor recusa `Unknown`
com valor, `Known` sem valor, `Unknown` ou `Estimated` sem código de razão, e qualquer facto
sem evidência. Campos privados.

Zero conhecido é `Known` com ponteiro para zero. Desconhecido não tem valor.

*Isto substitui o ADR-017* (ver ADR-034). A regra mais citada da casa — 1.378 citações no
harvest de `docs/architecture/decisions/_citations/adr-017-citations.md` — passa de doutrina
aplicada à mão a tipo cujo mau uso não compila.

*Defeito que teria apanhado:* `sourcekind/sourcekind.go:47`, onde qualquer sistema
desconhecido vira `LiveReadThrough` no ramo default. O núcleo partilhado viola a regra que o
resto do repo cita 1.378 vezes.

**Emenda 2026-08-07** (paga por `apps/server_core/internal/kernel/fact/knowledge.go:28`):
ratifica-se o código, não o texto acima como estava antes desta emenda — o valor zero de
`Knowledge` é `Unknown = iota`, nunca `Known = iota + 1`. Razão medida: `Fact[string]{}`
compila em qualquer pacote — um literal de struct vazio não nomeia campos, logo "não há
literal de struct" era falso nas duas variantes, com `Known` zero ou com `Unknown` zero. A
guarda real não é a inexistência do literal; é `e.IsZero()`, chamada dentro de todo
construtor (`NewKnown`, `NewEstimated`, `NewUnknown`, `NewNotApplicable`,
`knowledge.go:74,86,99,111`). Com `Unknown = iota`, o zero-value que esse literal produz é
o estado seguro, não o confiante — por isso o bloco de código acima já foi corrigido para
`Unknown Knowledge = iota` e o campo `reasonCode` renomeado para `reason`, como o código.

**Regra 4.2 — Nível 2.** Grandeza de negócio em `contracts` é `Fact[T]`, nunca número nu.
*Instrumento:* analisador sobre os tipos de campo em `contracts/`.

**Regra 4.3 — Nível 2.** `COALESCE(x, 0)` proibido em query de domínio financeiro.
*Instrumento:* analisador sobre SQL embebido.

**Regra 4.4 — Nível 1.** Cálculo com input `Unknown` produz resultado `Unknown`.
*Instrumento:* a aritmética vive em **funções genéricas de pacote**, `Map`/`Combine2`; não
há operador `+` sobre factos e não há aritmética em métodos de `Fact[T]`.
*Porquê, com o número real:* na simulação MG→BA com frete desconhecido `U`, a equação fecha
em `P = 178,69 + 2,02·U`. R$ 178,69 é o resultado contrafactual se `U = 0`. **Cada real
ignorado move o preço dois reais.** Um sistema que trate desconhecido como zero não erra um
pouco — erra o dobro do que ignorou.

**Emenda 2026-08-07** (paga por `apps/server_core/internal/kernel/fact/combine.go:10-15`):
"a aritmética vive em métodos de `Fact[T]`" era inimplementável — Go não permite que um
método declare parâmetros de tipo próprios (`func (f Fact[A]) Map[B any](...)` não
compila). `Map[A, B any](f Fact[A], method string, fn func(A) (B, error)) (Fact[B], error)`
e `Combine2[A, B, C any](...)` são funções de pacote em `internal/kernel/fact/combine.go:16,45`;
`Combine2` só chama `fn` quando os dois inputs são usáveis (`combine.go:46-49`), o que é a
garantia que a Regra 4.4 pede, só que expressa como função e não como método.

---

## 5. Identidade

**Regra 5.1 — Nível 1.** Identificador canónico é opaco e gerado por nós. Não deriva de dado
de fonte.

```go
type ProductID string  // ULID gerado; sem semântica de fonte

type SourceProductKey struct {
    Tenant         tenant.ID
    SourceSystem   string  // "sankhya", "spreadsheet"
    SourceInstance string  // ligação/importação
    ObjectKind     string  // "product", "variant"
    ExternalKey    string  // o código do Sankhya vive AQUI
}
```

*Instrumento:* `ProductID` é tipo próprio; não recebe `int`.
*Defeito que teria apanhado — e é o maior:*
`catalog/domain/canonical_product.go:12` declara `type InternalProductID int`, que **é o
`CODPROD` do Sankhya**, descrito como canónico. Enquanto isso for verdade, adicionar um
segundo ERP não é aditivo: é colisão de identidade a atravessar todos os contratos.

**Regra 5.2 — Nível 2.** EAN, SKU e código de fonte são *identificadores*, nunca identidade.
Podem faltar, repetir, ser reciclados ou mudar. Concordância de identificador é evidência
para uma decisão de vínculo, nunca a decisão.

**Regra 5.3 — Nível 1.** Identidade de item externo é tipo, não `string`.
*Defeito que teria apanhado:* `ProviderItemID` aparece **186 vezes** redeclarado como
`string` nu em seis módulos.

**Regra 5.4 — Nível 2.** Fusão de produtos canónicos é não-destrutiva: substitui-se com
histórico, nunca se apaga.

---

## 6. Tempo

**Regra 6.1 — Nível 1.** Três carimbos distintos, tipos distintos:
`SourceUpdatedAt` (quando a fonte diz que o facto passou a valer), `ObservedAt` (quando nós o
recolhemos), `EffectivePeriod` (validade legal ou de negócio). Um não é atribuível ao outro.

**Regra 6.2 — Nível 2.** Facto com natureza temporal — custo, alíquota, comissão,
classificação fiscal — carrega `EffectivePeriod` obrigatório.
*Porquê:* retroajustar vigência depois de as colunas serem números nus é reconstrução de
semântica. Seis meses depois já não se sabe o que era desconhecido.

---

## 7. Posse de dados

**Regra 7.1 — Nível 2.** Um esquema Postgres por contexto. **Um escritor por tabela** — só a
camada de aplicação e o repositório do contexto dono executam DML.
*Instrumento:* verificação sobre o SQL por pacote, mais `GRANT` por papel.
*Defeito que teria apanhado:* `products_mirror` tem hoje **dois escritores**, `erp_import` e
`internal_read`.

**Regra 7.2 — Nível 2.** Sem foreign key entre esquemas de contextos diferentes. A
integridade referencial cruzada valida-se por contrato e por reconciliação.
*Porquê:* FK cruzada é acoplamento que o compilador não vê. Sem ela, congelar pacotes sem
resolver posse de tabela deixa o `contracts` virar rota legal de lavagem — fronteira verde,
domínio acoplado pelo Postgres.

**Regra 7.3 — Nível 2.** Toda a chave primária e única de linha multi-tenant começa em
`tenant_id`. Toda a referência é composta: `(tenant_id, product_id)`, nunca `product_id` só.

**Regra 7.4 — Nível 2.** RLS ligado, com `FORCE ROW LEVEL SECURITY`, e papel de aplicação
que não é dono das tabelas e não tem `BYPASSRLS`.
*Porquê:* `WHERE tenant_id = ?` na aplicação é disciplina, não garantia.

---

## 8. Escrita

**Regra 8.1 — Nível 1.** Toda a escrita para sistema externo passa por `changecontrol`.
*Instrumento:* a porta de execução é de `changecontrol`; o adapter de escrita implementa-a e
nenhum contexto a alcança.

Máquina de estados:

```
proposed → previewed → awaiting_approval → approved → executing
         → succeeded | failed | outcome_unknown
         → verified | diverged
```

**Regra 8.2 — Nível 1.** O resultado da execução tem três desfechos, não dois:
`ProviderSucceeded`, `ProviderFailed`, `ProviderUnknown`. Timeout depois de enviar é
`Unknown`.
*Instrumento:* o tipo tem três variantes; não há booleano de sucesso.
*Porquê:* retentar às cegas uma escrita não idempotente pode sobrepor ação mais recente do
utilizador ou duplicar. **HTTP 200 significa aceite, não convergido.**

**Regra 8.3 — Nível 1.** `desired`, `last_requested` e `last_observed` são três valores
separados e persistidos. A convergência prova-se por leitura posterior, nunca pela resposta
da escrita.

---

## 9. Ingestão

**Regra 9.1 — Nível 1.** O agendador é plataforma e só possui: registo de job, partição por
tenant e conta, leases com token de cerca, coordenação de rate limit, backoff, telemetria.
*Instrumento:* `platform/scheduler` não importa nenhum contexto.

**Regra 9.2 — Nível 2.** Cursor, janela de sobreposição, chave de deduplicação, semântica de
frescura e política de varredura completa pertencem ao contexto.
*Defeito que teria apanhado:* `sync/domain/sync_state.go:15` é um enum central que cobre
produtos, anúncios, pedidos, mercado, tarifas e ICMS. Capacidade de negócio nova obriga a
editar um enum técnico.

**Regra 9.3 — Nível 2.** Observação crua é retida e imutável, com hash. Ingestão é idempotente
por hash. Cursor só avança na mesma transação em que a observação é aceite.

**Regra 9.4 — Nível 2.** Publicação de evento na mesma transação da mudança de estado, por
outbox. Nunca publicar depois de fazer commit.

---

## 10. Leitura

**Regra 10.1 — Nível 1.** Pergunta com semântica de domínio vai ao contexto dono, pela porta.
"Stock vendável agora" é `inventory`. "Custo à data X pelo método Y" é `costing`. Nunca é
uma query.
*Instrumento:* sem acesso ao esquema alheio e sem FK cruzada, não há query a escrever.

**Regra 10.2 — Nível 2.** Painel usa projeção descartável em `views/`, com um caso de uso
nomeado, um escritor, alimentada por eventos versionados. Expõe `as_of` e completude. Não é
autoritativa, não emite comando, não tem regra de negócio.
*Defeito que teria apanhado:* `dashboard/ports/sources.go:15` declara seis portas para seis
contextos — um módulo-deus de query com licença para juntar tudo.

**Regra 10.3 — Nível 2.** Projeção que precisa de cinco contextos nomeia os cinco contratos
de evento. Não existe `query` genérico.

---

## 11. Contrato externo e SDK

**Regra 11.1 — Nível 2.** O OpenAPI é a fonte. Os tipos do servidor Go e o SDK TypeScript são
**gerados** a partir dele.
*Instrumento:* `go generate` mais gerador de SDK, com o gate do §13.
*Defeito que teria apanhado:* SDK com **2.595 linhas e 172 interfaces à mão**, a mesma forma
copiada em quatro sítios, e um gate de governança que só exige *mesmo commit* — nunca
concordância.

**Regra 11.2 — Nível 2.** O dialeto do OpenAPI é normalizado antes de qualquer geração.
*Medido:* `contracts/api/marketplace-central.openapi.yaml:1` declara 3.1.0 e usa
`nullable: true` em estilo 3.0, inclusive em campo obrigatório na linha 3536. Em 3.1
`nullable` não é normativo. Geradores divergem no tratamento.

---

## 12. Evidência

**Regra 12.1 — Nível 2.** Aceite é medição, nunca teste verde. Critério de aceite escreve-se
como número que se conta.

**Regra 12.2 — Nível 2.** RED antes do código, contra árvore limpa. Todos os RED primeiro.
Itera-se até a implementação condizer com o que é para ser, nunca até o teste condizer com a
implementação.

**Regra 12.3 — Nível 2.** Lane hermética não é superconjunto da de unidade. Teste pulado e
teste verde são byte-idênticos sem `failure_token`.

**Regra 12.4 — Nível 3 → cortada.** "Revisão de arquitetura" como critério não sobrevive à
regra do §0. Ou vira analisador, ou sai.

---

## 13. O gate de geração

```
checkout limpo
versões de gerador pinadas (Go e Node)
go generate ./...        a partir de apps/server_core
gerador de SDK --check
git diff --exit-code
git diff --cached --exit-code
git status --porcelain --untracked-files=all   vazio
go build ./... && go vet ./... && tsc
lint do OpenAPI
```

**Porquê tudo isto e não só `git diff`:** `git diff --exit-code` **ignora ficheiro novo não
rastreado**. Um gerador que cria um ficheiro novo passa verde. As três verificações juntas
fecham o buraco.

O que este gate prova: que o gerado está em sincronia com a fonte, refazendo-o. O que não
prova: correção semântica, validade de contrato, ou concordância com o comportamento em
runtime.

---

## 14. Emendas à ADR-023

A coluna **Aplicado** diz se a emenda já está escrita no ficheiro do ADR — as três
primeiras entraram em `ccdc49c`. Enquanto disser
*pendente*, a ordem de verdade do repositório põe o ADR acima deste documento e é o ADR que
vale — uma emenda que vive só aqui não emendou nada.

| § | Estado | Aplicado |
|---|---|---|
| §1 (o que é um módulo) | Substituído por §1.1 e §1.2 deste documento. `contexts/`, não `modules/`. Chave do registo é o par `(kind, id)`; `kind` ausente = `"module"`. Pasta não inscrita dá `GOV_CONTEXT_UNREGISTERED`, medido percorrendo a árvore e não o registo. | **sim** |
| §2 (só `X/ports` importável) | Substituído por §1.1. **Nível 1, não nível 3** — a frase *"one line, and a grep can check it"* foi **apagada**, não anotada. A regra `internal/` do Go vale para qualquer `internal/` da árvore. Os dois carve-outs continuam dispositivos em `internal/modules/`, cujas camadas estão em diretórios comuns; sob `internal/contexts/` o carve-out 1 (raiz) **não é exercível** — o compilador recusa-o, medido em `catalog_wiring.go:9`. Carve-out 2 (`sourcekind`) extingue-se quando o módulo desaparecer para `kernel/provenance`. | **sim** |
| §2-a (a fachada) | **Novo.** Não existia no ADR. Qualquer árvore com `internal/` força construtor que monta os próprios internos; parâmetro de tipo interno é inchamável de fora. Escrita antes só para vendors, foi por isso que não alcançou `contexts/`. | **sim** |
| §3 (porta basta-se) | Mantido, absorvido em §2.1. | n/a — sem alteração |
| §4 (quando existe porta) | Mantido. | n/a — sem alteração |
| §5 (núcleo partilhado) | Substituído por §3. Um membro passa a seis, com critério de admissão escrito. | pendente |
| §6 (um tipo de dinheiro) | Substituído por §3.2, agora nível 1. | pendente |
| §7 (vocabulário de camadas fechado) | Substituído por §1.1. *Só 5 dos 21 módulos cumpriam a lista fechada.* | pendente |
| §8 (módulos fora do molde) | Extinto: o molde deixa de ter exceção. | pendente |

E o ADR-017 ("desconhecido nunca é zero") é substituído pela §4.1: deixa de ser regra citada
e passa a tipo cujo mau uso não compila.

---

## 15. Sequência

Sem strangler. Não há produção e os dados são re-deriváveis do Sankhya e do Mercado Livre.
A camada de compatibilidade que torna o strangler caro existe para não derrubar produção —
não é o nosso caso. `internal/contexts/` nasce ao lado de `internal/modules/`, e cada módulo
velho é apagado quando o contexto que o substitui aterra.

1. **Kernel.** Os seis tipos, com os construtores da §4.1 e §3.2. Não parte nada.
2. **Analisadores.** As regras de nível 2 que o resto depende: DTO de vendor, token de vendor,
   `float64`, aresta de contexto.
3. **Fatia vertical completa.** Um produto do ERP → `catalog` → `listings` → `linking` →
   `inventory` → `pricing` com `tax` e `costing` → `changecontrol` → `orders` →
   `reconciliation` → `profitability`. Com payloads reais gravados, reconciliada contra os
   sistemas externos.
4. **Ratificação.** A fatia confirma ou refuta este documento. Emenda-se antes de escalar.
5. **Migração contexto a contexto**, com o compilador a enumerar o que falta. `internal/`
   fecha atrás de cada contexto que aterra, não de todos ao mesmo tempo.
6. **Contrato e SDK** por geração, depois de o dialeto do OpenAPI estar normalizado.

A fatia vertical é deliberadamente estreita e funda, e vem antes do molde universal. Foi a
recomendação independente das duas consultas: o molde sai das fatias, não o contrário.

---

## 16. O que este documento não decide

- As fórmulas fiscais. Este documento diz onde vivem e que forma têm os inputs e outputs;
  não diz que estão certas. Isso é matéria de especialista fiscal e de casos dourados
  revistos.
- Se `goverter` entra. Mapeamento gerado é atrativo mas
  `useZeroValueOnPointerInconsistency` e `ignoreMissing` restauram o zero silencioso. Se
  entrar, entra com essas opções proibidas e é dependência nova — `REQUEST` ao hub.
- O gerador de OpenAPI. `oapi-codegen` estável é 3.0; `ogen` precisa de spike. Decide-se
  depois da normalização do dialeto (§11.2).
- Que contextos podem ser fundidos na prática. `costing` separado de `catalog` é a separação
  mais discutível deste documento e a fatia vertical é que decide.

---

## 17. Log de emendas a este documento

Uma linha por emenda, com data e o `file:line` medido que a pagou. Distinto do quadro do
§14, que regista o estado das emendas à ADR-023; este regista emendas às regras do
protocolo em si. `ADR-033` e `ADR-034` (`docs/architecture/decisions/`) são os documentos
que dão a estas emendas prioridade sobre decisões congeladas, quando aplicável.

| Data | Regra | Emenda | `file:line` que pagou |
|---|---|---|---|
| 2026-08-07 | 4.1 (valor zero de `Knowledge`) | Ratifica-se o código: `Unknown = iota`, não `Known = iota + 1`; a guarda real é `Evidence().IsZero()`, não a ausência de literal de struct. | `apps/server_core/internal/kernel/fact/knowledge.go:24-36,74,86,99,111` |
| 2026-08-07 | 4.4 (aritmética) | De "métodos de `Fact[T]`" para funções genéricas de pacote `Map`/`Combine2` — um método Go não pode declarar parâmetros de tipo próprios. | `apps/server_core/internal/kernel/fact/combine.go:10-16,45` |
| 2026-08-07 | 1.3 (instrumento cross-context) | O instrumento é o detector `ScanCrossContextInternal` **sem** o salto de ficheiros fora de `contexts/` — o único sítio que violou a regra estava fora, e um detector que pula "fora de contexto" é cego a exatamente esse caso. | `apps/server_core/internal/arch/scan.go:131-133` |
| 2026-08-07 | 2.2-a (fachada de raiz única) | Sobe de regra de adaptador a regra do molde: vale para qualquer árvore com `internal/`, contexto incluído. Escrita só para vendors, não alcançava `contexts/`; `catalog.New(store, ids, reader)` nasceu com zero chamadores legais pela mesma razão. | `apps/server_core/internal/contexts/catalog/module.go:23-38` |
| 2026-08-07 | 2.3 (tokens de vendor) | O escopo (`adapters/` isento, `internal/arch/` isento) vive no detector, na função `vendorRuleApplies` — sem filtro de caminho o detector acusaria a sua própria lista de tokens e nunca ficaria verde. | `apps/server_core/internal/arch/scan.go:34-37,217-226` |
| 2026-08-10 | 2 (molde de vendor) | **Sem emenda.** A perna `listings` construiu o segundo vendor do molde e nada foi refutado: fachada de raiz única, DTO de fio confinado em `internal/api`, porta declarada pelo contexto e implementada pelo adapter — compilou na primeira composição, sem exceção, sem entrada de governança nova além do registo do contexto. | `apps/server_core/internal/adapters/marketplace/mercadolivre/mercadolivre.go:32-46`, `.../internal/api/items.go:14-112`, `apps/server_core/internal/contexts/listings/port/feed.go:32-34` |
| 2026-08-10 | 2.2-a (fachada de raiz única) | Confirma-se pelo custo, não pela leitura: `Config.Token` é `func(ctx) (string, error)` e **não** `api.TokenSource`, porque um tipo de `internal/` na fachada não compila para chamador de fora da árvore. A regra não é estilística — o compilador é quem a aplica. | `apps/server_core/internal/adapters/marketplace/mercadolivre/mercadolivre.go:20-26` |
| 2026-08-10 | 9 (desconhecido nunca vira default) — **regra nova** | O canal decide, não a guarda. Parâmetro de invocação (que tenant, que instância, que tamanho de página) viaja como **argumento**; configuração de deployment (URL, chave) viaja como ambiente e é lida por **um** carregador tipado. Medido: os três pontos de entrada liam `MC_DEFAULT_TENANT_ID` cru para distinguir "não definido" de "definido como `tenant_default`" — `pgdb.LoadConfig` substitui um pelo outro. Essa guarda (`requireTenantConfigured`) era sintoma do canal errado: uma flag obrigatória não tem default com que se confundir, e o fail-closed passa a ser propriedade da construção. Ratchet `governance-drift` **encolheu** 58 → 53. | `apps/server_core/cmd/catalogingest/main.go:82-109`, `apps/server_core/cmd/listingsingest/main.go:78-100`, `apps/server_core/cmd/mlprobe/main.go:74-96`, `contracts/gate/baselines.json` |
| 2026-08-10 | 2.3 (tokens de vendor) | Emenda de escopo medida: `cmd/` **não** é isento. O ponto de entrada de operador não pode nomear o próprio vendor — o host da API e o `provider_code` em `cmd/listingsingest` levaram `adapters/vendor-token-outside-adapters` de 4 para 9. A identidade do vendor passa a viver como constante da raiz de composição, e o `cmd/` referencia por identificador; o `provider_code` entra no SQL como parâmetro, não como literal. | `apps/server_core/internal/composition/listings_wiring.go:19-30`, `apps/server_core/cmd/listingsingest/main.go:166-175` |
