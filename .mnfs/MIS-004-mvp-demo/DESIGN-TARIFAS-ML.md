# DESIGN — Resolução de Tarifas ML (comissão + frete por produto)

```yaml
id: DESIGN-TARIFAS-ML
mission: MIS-004-mvp-demo
type: design-doc
author: cold-planner (Opus, contingency lane §12)
status: ratificado (operator, 2026-07-19 — ver HUB-LEDGER D-78)
input_brief: scratchpad/design-brief-tarifas.md (HUB v2, 2026-07-19)
base_sha: ca959ddaa37c258109600604de13fe862e8d8e6a
head_verified: ca959ddaa37c258109600604de13fe862e8d8e6a (== base, sem delta)
demo: 2026-07-20 (T-1)
```

> **Natureza deste documento.** Design-only. Nenhum código é escrito ou editado aqui.
> Todos os endpoints ML são READ (zero-write preservado, `MPC_PROVIDER_WRITES_ENABLED`
> unset). Fatos de API são estritamente os do brief §3 (doc-grounded); nada é fabricado
> além disto. Fatos de código foram re-verificados contra o repo em `ca959dd` — deltas
> anotados inline onde a realidade diverge do brief.

---

## 1. Contexto & problema

O solver M-07 ("margem alvo → preço", `pricing/domain/solve.go`) **bloqueia** com
`desconhecidos:["frete"]` para qualquer produto sem override de frete, porque o frete
**não tem fonte de resolução** no cálculo single-item — só o override manual do cliente
(`calc_service.go:213,276` passam `optionalMoney(req.FreteProduto)`, nil por padrão). A
comissão vem de `marketplace_fee_schedules` por **categoria** (seeded/manual, migration
0011), não do produto real.

O operador exige (brief §1) que cada produto **tenha comissão e frete mapeados**, com duas
fontes de verdade:
- **(a) VENDA real** — registrada por produto: data, comissão, frete, valor do pedido →
  tabela histórica para acompanhar variação no tempo;
- **(b) COTAÇÃO na API ML** no início de vida do produto (mesma fonte do simulador do
  site deles).

O cliente entra via Excel (xlsx ERP), que **nunca traz peso/dimensões nem preço**
(`erp_import` parser: CODPROD/DESCRPROD/CUSTO/ESTOQUE_FISICO obrigatórios; EAN/REFFORN/
MARCA/NCM opcionais — brief §4). O `R$79` hardcoded no domain (`decompose.go:20`,
`solve.go:13`) e o `R$150` de probe (`solve.go:101`) são achismos **banidos** — a regra ML
muda 2026-03-02 e o custo de frete-grátis é variável por item (sinal =
`mandatory_free_shipping`).

Padrão **Global Maximum** em código e planejamento. Este design formaliza: a escada de 4
degraus, o modelo de dados histórico, os ports/adapters, o fluxo de início-de-vida, a UI,
os 12 cenários de falha, o faseamento vs demo T-1, e o fatiamento em chips.

---

## 2. Arquitetura da escada de 4 degraus

### 2.1 Princípio de resolução

A tarifa de um produto tem **dois componentes independentes** — comissão e frete — que
**caminham a escada separadamente**. Um produto pode ter comissão VENDA (degrau 1) mas
frete só por estimativa (degrau 4); o cenário 11 (quote parcial) é a norma, não a exceção.
Cada componente resolvido carrega seu **próprio carimbo** (`fonte`, `degrau`, `data`,
flags). A resolução é **determinística e preferencial**: tenta o degrau mais alto viável,
degrada rotulando.

Value object retornado (domínio pricing, ver §4):

```
TariffResolution
├── Comissao  ComponentResolution   // pct + fixed_fee opcional
├── Frete     ComponentResolution   // Money, pode ser NO-DATA (nil, nunca 0)
└── (cada ComponentResolution)
    ├── Valor        *Money / pct string   // nil = NO-DATA honesto (ADR-17)
    ├── Fonte        VENDA | COTACAO | PADRAO
    ├── Degrau       1 | 2 | 3 | 4
    ├── Data         time.Time              // data do evento (venda ou cotação)
    ├── CategoryID   string
    ├── ListingType  gold_special | gold_pro
    ├── Estimativa   bool                   // true no degrau 4
    ├── Predita      bool                   // categoria via domain_discovery não confirmada
    └── Staleness    { FetchedAt, TTL, Stale bool }
```

### 2.2 Degrau 1 — VENDA (fonte: pedido real)

**Pré-requisito:** produto vendido ≥1x (existe linha `product_tariff_history` com
`fonte=VENDA`).

```
[order ingestion / projeção M-08]
  order_items[].sale_fee ─────────────► comissao_amount + comissao_pct (derivado / valor_pedido)
  GET /shipments/{id}/costs             frete_amount ◄── senders[0].cost
    header x-format-new: true           (HOJE live-only, shipping_reader.go:81; PRECISA persistir)
  order total ──────────────────────► valor_pedido
        │
        ▼  (grava 1 linha VENDA por venda, com observed_at)
  product_tariff_history (nova tabela §3)
        │
        ▼  (resolver lê a linha VENDA mais recente não-stale)
  ComponentResolution{ fonte:VENDA, degrau:1, data:observed_at }
```

**Campos exatos:** comissão de `order_items[].sale_fee` (já persistido, `orders/domain/
order.go:72` `SaleFeeAmount *float64` nullable — **verificado**, migration 0027); frete de
`senders[0].cost` (mapeado em `shipping_reader.go:167`, **SEM persistência** — 0 matches
`shipment_cost*` em migrations, **verificado**). `sale_fee` é o valor absoluto da venda; o
`comissao_pct` é derivado `sale_fee / valor_pedido × 100` e **carimbado como derivado**
(não é a taxa nominal ML, é a taxa efetiva daquela venda — honestidade de fonte).

**Source labeling:** `fonte=VENDA`, `degrau=1`, `data=observed_at` (data do pedido).
**ADR-17:** se `sale_fee` nil na venda → comissão VENDA indisponível para aquele produto
(cai para degrau inferior); se `senders[].cost` ausente/decode falhou → frete VENDA
indisponível, componente frete cai (nunca 0).

### 2.3 Degrau 2 — COTAÇÃO-ANÚNCIO (fonte: quote por item vinculado)

**Pré-requisito:** vínculo **RESOLVED** em `product_links` (não REVIEW — cenário 9). Só um
vínculo resolvido dá o `item_id` MLB confiável.

```
[product_links: link RESOLVED] → item MLB (category_id do próprio item)
        │
        ├─ comissão: GET /sites/MLB/listing_prices?price=&category_id=&listing_type_id=gold_special|gold_pro
        │             &currency_id=BRL   (FeeQuoteReader JÁ existe, capability_adapter.go:499)
        │             → sale_fee_details.{percentage_fee, fixed_fee}
        │
        └─ frete:    GET /users/{uid}/shipping_options/free?item_id={MLB}
                     (GetFreeShippingCost JÁ existe p/ item_id, shipping_reader.go:102-121)
                     → cost / coverage.all_country.list_cost
        │
        ▼
  ComponentResolution{ fonte:COTACAO, degrau:2, data:fetched_at }
```

**Endpoints/campos:** `FeeQuoteInput{SiteID,CategoryID,ListingTypeID,PriceAmount,CurrencyID}`
(`capability.go:280` — **verificado**) → `FeeQuoteSnapshot.CommissionPercent/FixedFeeAmount`
(`*float64` nullable — **verificado**). Frete: `FreeShippingQuery{ItemID}` (`shipping_read.go:22`
— **verificado**, hoje só `ItemID`). **Nenhuma extensão de adapter necessária no degrau 2**
para o caminho por `item_id` — os dois readers já existem e estão orquestrados
(`ProviderOperationService.ReadFeeQuote` :149, rota probe `probes/fee-quote`
`auth_handler.go:309`). O que falta é **ligá-los ao pricing calc** (não estão — brief §4).

**Preconditions/labeling:** vínculo RESOLVED obrigatório; `fonte=COTACAO`, `degrau=2`,
`data=fetched_at`; `category_id` real do item (não predita → `Predita=false`).
**ADR-17:** `list_cost=0` por call under-spec (item sem `mandatory_free_shipping`) → tratar
como NO-DATA, nunca frete 0 (cenário 6).

### 2.4 Degrau 3 — COTAÇÃO-MATCH (fonte: quote por match de catálogo ML)

**Pré-requisito:** EAN no xlsx **OU** descrição pt-BR (para `domain_discovery`).

```
EAN → GET /products/search?site_id=MLB&product_identifier={EAN}   (prioridade sobre q)
        │  resultado NÃO traz category_id, só domain_id
        ├─(a) GET /products/{id} → buy_box_winner.category_id   (pode ser null — cenário 3)
        └─(b) GET /sites/MLB/domain_discovery/search?q={nome}&limit=3 → category_id ranqueado
              (SEM score de confiança → top-3 p/ humano, cenário 4; categoria = PREDITA)
        │
        ├─ comissão: listing_prices com a category_id resolvida/predita (mesmo endpoint degrau 2)
        │
        └─ frete:    GET /users/{uid}/shipping_options/free?dimensions=HxWxL,gramas&item_price=&free_shipping=
                     (SÓ se dims existem — product.HeightCM/WidthCM/LengthCM/WeightG,
                      catalog/domain/product.go:18-21, *float64 nullable — VERIFICADO)
        │
        ▼
  ComponentResolution{ fonte:COTACAO, degrau:3, data:fetched_at, Predita:(via domain_discovery) }
```

**DELTA de código (extensão necessária, dono = chip connectors):**
1. `FreeShippingQuery` (`connectors/domain/shipping_read.go:22`) hoje **só tem `ItemID`** →
   estender com `Dimensions string` (`HxWxL,gramas`, cm/g inteiros), `ItemPrice float64`,
   `FreeShipping *bool`. O adapter `getShipmentInfo`/`getFreeShippingCost`
   (`shipping_reader.go:102`) monta hoje só `item_id` → adicionar o ramo `dimensions=`.
2. Caminho EAN→catálogo→categoria **não existe** como reader — novo método no
   `CapabilityAdapter` (`/products/search`, `/products/{id}`, `/domain_discovery/search`),
   payloads morrem no adapter, port normalizado no domínio connectors (IC-06). Categoria
   sem `buy_box_winner` → `domain_discovery` no nome (cenário 3).

**Labeling:** `fonte=COTACAO`, `degrau=3`; categoria via `domain_discovery` ⇒ `Predita=true`
até confirmação humana (**nunca auto-gravar como fato** — cenário 4). Sem dims ⇒ frete degrau
3 indisponível → frete cai para degrau 4 política "sem dados" (cenário 5), **nunca 0**.

### 2.5 Degrau 4 — PADRÃO (fonte: default configurável)

**Pré-requisito:** nenhum — **sempre disponível**, é o piso da escada.

```
config pricing_tariff_defaults (nova, §6.4) por tenant/installation:
  comissao: Clássico 13% / Premium 16%  (valores iniciais do operator, EDITÁVEIS)
  frete:    estimativa configurável OU política "sem dados"
        │
        ▼
  ComponentResolution{ fonte:PADRAO, degrau:4, Estimativa:true }
```

**Labeling obrigatório:** degrau 4 **SEMPRE** rotulado `ESTIMATIVA` na UI (§6). Comissão
degrau 4 vem da config (nunca hardcoded no domain — os 13%/16% são seed inicial editável,
**não** constante de código; o `R$79`/`R$150` do domain são removidos, §4.4). Frete degrau 4:
se o operador configurou uma estimativa → usa rotulada ESTIMATIVA; se não → política "sem
dados de frete" (componente frete = NO-DATA, `Frete.Valor=nil`, UI mostra "—" + instrução).

**ADR-17 no degrau 4:** o degrau 4 nunca inventa frete real; ou é uma estimativa
**explicitamente rotulada** ou é NO-DATA rotulado. A comissão degrau 4 é sempre disponível
(default configurável) — é o único componente garantido, permitindo o solver produzir um
preço (com margem rotulada ESTIMATIVA).

### 2.6 Tabela-resumo da escada (fontes de dados)

| Degrau | Fonte | Comissão (endpoint/campo) | Frete (endpoint/campo) | Pré-requisito | Reuso de código |
|---|---|---|---|---|---|
| 1 VENDA | pedido real | `order_items[].sale_fee` (persistido 0027) | `shipments/{id}/costs`→`senders[0].cost` (x-format-new; **persistir**) | vendido ≥1x | order.go:72; shipping_reader.go:167 |
| 2 COTAÇÃO-ANÚNCIO | quote item vinculado | `listing_prices` cat. do item | `shipping_options/free?item_id` | link RESOLVED | FeeQuoteReader :499; GetFreeShippingCost :102 (prontos) |
| 3 COTAÇÃO-MATCH | quote match catálogo | EAN→search→cat→`listing_prices` | `shipping_options/free?dimensions` | EAN ou descrição pt-BR + dims p/ frete | **extensão** FreeShippingQuery+search |
| 4 PADRÃO | default config | Clássico 13%/Premium 16% editável | estimativa config OU "sem dados" | nenhum | nova config; rótulo ESTIMATIVA |

---

## 3. Modelo de dados

### 3.1 Tabela `product_tariff_history` (nova, migration 0068)

Colunas mínimas operator-stated: produto, installation, data, comissão, frete,
valor_pedido, fonte `VENDA|COTACAO`. Adicionadas pelo planner: degrau, listing_type,
category_id, campos de qualidade/staleness, tenant scoping, índices.

```sql
-- migration 0068 (bloco a alocar pelo hub; 0068-0069 livres, 0070-0074 = reserva hub)
CREATE TABLE product_tariff_history (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           TEXT        NOT NULL,               -- ADR: escopo tenant em toda linha (cenário 12)
    installation_id     TEXT        NOT NULL,               -- multi-installation scoping
    product_id          TEXT        NOT NULL,               -- CODPROD (SKU canônico, ADR-03)

    fonte               TEXT        NOT NULL CHECK (fonte IN ('VENDA','COTACAO')),
    degrau              SMALLINT    NOT NULL CHECK (degrau BETWEEN 1 AND 4),
    observed_at         TIMESTAMPTZ NOT NULL,               -- data do evento: venda OU cotação

    -- comissão (nullable = desconhecido, NUNCA 0 default — ADR-17)
    comissao_pct        NUMERIC(7,4) NULL,                  -- taxa (nominal p/ COTACAO, efetiva/derivada p/ VENDA)
    comissao_pct_source TEXT NULL CHECK (comissao_pct_source IN ('nominal','derivado')),
    comissao_fixed_amount NUMERIC(14,2) NULL,               -- sale_fee_details.fixed_fee (degrau 2/3)

    -- frete (nullable = NO-DATA honesto — ADR-17; jamais 0 por under-spec)
    frete_amount        NUMERIC(14,2) NULL,
    frete_source_label  TEXT NULL,                          -- 'senders_cost' | 'shipping_options_item' | 'shipping_options_dims' | 'estimativa'

    valor_pedido        NUMERIC(14,2) NULL,                 -- order total (VENDA) ou item_price base (COTACAO)
    currency_id         TEXT        NOT NULL DEFAULT 'BRL',

    listing_type_id     TEXT NULL,                          -- gold_special (Clássico) | gold_pro (Premium)
    category_id         TEXT NULL,
    category_predita    BOOLEAN     NOT NULL DEFAULT FALSE, -- via domain_discovery (cenário 4), não confirmada

    -- qualidade / staleness
    fetched_at          TIMESTAMPTZ NULL,                   -- quando a cotação foi buscada (COTACAO)
    source_updated_at   TIMESTAMPTZ NULL,                   -- se ML expõe (FeeQuoteSnapshot.SourceUpdatedAt)
    is_estimativa       BOOLEAN     NOT NULL DEFAULT FALSE, -- degrau 4
    quote_ttl_seconds   INTEGER NULL,                       -- TTL configurado no momento da cotação

    provider_ref        TEXT NULL,                          -- order_id (VENDA) ou endpoint (COTACAO), rastreio
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- resolução: linha mais recente por (tenant, installation, produto), filtrada por fonte/staleness
CREATE INDEX idx_pth_resolve
    ON product_tariff_history (tenant_id, installation_id, product_id, observed_at DESC);
-- caminho VENDA (degrau 1) rápido
CREATE INDEX idx_pth_venda
    ON product_tariff_history (tenant_id, installation_id, product_id, observed_at DESC)
    WHERE fonte = 'VENDA';
```

**Notas de schema:**
- `comissao_pct` e `frete_amount` **nullable** — um quote parcial (cenário 11: comissão OK,
  frete falhou) grava a linha com `frete_amount = NULL` + `frete_source_label = NULL`; a UI
  mostra frete "—". Nunca se grava 0 por ausência.
- `comissao_pct_source`: distingue a taxa **nominal** ML (COTAÇÃO, `listing_prices`) da taxa
  **efetiva derivada** de uma venda (`sale_fee / valor_pedido`). Honestidade de fonte —
  degrau 1 não finge ser a taxa de tabela.
- Uma linha por evento (append-only histórico). O resolver lê a **mais recente viável** por
  degrau; o histórico completo alimenta o acompanhamento de variação no tempo (operator §1a).

### 3.2 Relação com `marketplace_fee_schedules` (categoria)

`marketplace_fee_schedules` (migration 0011) é **por categoria** (marketplace_code,
category_id, listing_type → commission_percent + fixed_fee_amount; lookup 2-nível com
fallback `default`, `fee_schedule_service.go:45`). **Coexistem**; a precedência é:

```
produto (product_tariff_history, degraus 1-3)  ─── sobrepõe ───►  categoria (marketplace_fee_schedules)  ─── fallback ───►  default configurável (degrau 4)
```

`marketplace_fee_schedules` **não** é um degrau da escada de produto — é a fonte de
comissão **por categoria** que hoje o pricing usa via
`pricing/adapters/feeschedule/adapter.go:23`. Ela permanece como **fallback intermediário**
opcional entre degrau 3 e degrau 4 (comissão de categoria seeded quando não há cotação de
produto mas há a categoria). Decisão de ratificação (§10): manter fee_schedules como
degrau "3.5" fallback de comissão, ou colapsar tudo em degrau 4? **Recomendação:** manter
como 3.5 (já ligado ao pricing; barato), rotulado `fonte=CATEGORIA` para honestidade.

### 3.3 Persistência do frete real (`senders[].cost`)

Hoje `senders[0].cost` é **live-only** (mapeado em `shipping_reader.go:167`, sem tabela).
**Onde gravar:** na **escrita/projeção do pedido** (M-08 orders). Quando um pedido é
ingerido e seus custos de shipment resolvem (`ShipmentCosts.SenderCost`,
`shipping_read.go:19`), grava-se **uma linha VENDA** em `product_tariff_history`:
`comissao` de `sale_fee`, `frete_amount` de `senders[0].cost`, `valor_pedido` do total,
`observed_at` da data do pedido, `provider_ref` do order_id.

Isto **acopla** ao caminho de ingestão de pedidos do M-08 — que está sendo retrabalhado no
`CHIP-M08-SHIPFIX` (decode de shipment + bucket). **Colisão** (ver §9): a escrita VENDA
deve entrar **após** o shipfix estabilizar o decode de `senders[].cost`, senão persiste
frete corrompido. Alternativa desacoplada (pós-demo): job de backfill que varre pedidos já
projetados e emite linhas VENDA — evita tocar o write path quente durante a demo.

---

## 4. Ports & wiring

### 4.1 Port `TariffResolver` (domain-facing pricing)

Novo port no `pricing/ports` consumido pelo `CalcService`. Substitui a resolução implícita
(hoje inexistente) de comissão/frete.

```go
// pricing/ports/tariff.go (novo)
type TariffResolver interface {
    // Resolve caminha a escada e devolve comissão + frete carimbados.
    // NUNCA erra por dado ausente — devolve ComponentResolution com Valor nil (ADR-17).
    Resolve(ctx context.Context, req TariffRequest) (TariffResolution, error)
}

type TariffRequest struct {
    TenantID       string
    InstallationID string
    ProductID      int
    Modalidade     domain.Modalidade   // classico | premium | full → listing_type
    PriceBasis     *string             // preço p/ listing_prices (decompose) ou nil (solver resolve depois)
    AsOf           time.Time
}
```

`TariffResolution`/`ComponentResolution` conforme §2.1. O resolver é **aplicação**
(orquestra adapters por degrau); o **domínio** (`Decompose`/`SolveTargetPrice`) permanece
puro e recebe `ComissaoPct`/`FreteProduto` **já resolvidos** (nota G1 de `decompose.go:26`
preservada — a engine continua hermética).

### 4.2 Adapters por degrau

| Degrau | Adapter | Fonte de dados | Dono (chip) |
|---|---|---|---|
| 1 VENDA | `pricing/adapters/tariffhistory` (novo) | lê `product_tariff_history` fonte=VENDA | pricing |
| 2 COTAÇÃO-ANÚNCIO | reusa `FeeQuoteReader` + `ShipmentReader.GetFreeShippingCost(item_id)` via IC-06 | ML live (readers prontos) | pricing consome, connectors publica |
| 3 COTAÇÃO-MATCH | `FeeQuoteReader` + `GetFreeShippingCost(dimensions)` **estendido** + novo catalog-search reader | ML live (extensão) | connectors |
| 4 PADRÃO | `pricing/adapters/tariffdefaults` (novo, lê config) | `pricing_tariff_defaults` | pricing |

O resolver tenta 1→2→3→3.5(categoria)→4 por componente, aplicando staleness (§5). Cotações
(degraus 2/3) são **cacheadas** gravando linha `fonte=COTACAO` em `product_tariff_history`
com `fetched_at`+`quote_ttl_seconds` — a próxima resolução reusa dentro do TTL sem re-bater
a ML (rate-limit friendly, cenário 7).

### 4.3 Substituição do probe R$150 do solver

Hoje `solve.go:100-102` chama `margemDecompose(15000)` (R$150 hardcoded) para detectar
unknowns estruturais, e como não há resolver de frete, `FreteProduto` é sempre nil ⇒
`desconhecidos:["frete"]` sempre. **Fix:** o `CalcService.SolveTarget`/`Decompose` chama
`TariffResolver.Resolve` **antes** do domínio e injeta `FreteProduto` resolvido (degrau 1-3
real, ou degrau 4 estimativa rotulada, ou NO-DATA explícito). Com frete resolvido, o probe
não força mais falso-desconhecido. O `R$150` e o `R$79` (`thresholdCents`, `taxaFixaLimiar`)
saem do domínio — o limiar de frete-grátis passa a ser **sinal por item**
(`mandatory_free_shipping`) resolvido no adapter, não constante (cenário 8; ver §4.4).

### 4.4 Fix dos 3 layers do FINDING-P7-SOLVER (ledger D-74)

**Layer 1 — domínio (`solve.go` / `decompose.go`): "blocking" vs "unreachable" vs
"segment-conditional".** Hoje `structuralUnknowns()` avalia em R$150 (≥79) e reporta frete
como desconhecido estrutural mesmo quando o preço atingível seria <79 (onde frete nem se
aplica). **Design:**
- Frete resolvido upstream (§4.3) elimina o falso-desconhecido na maioria dos casos.
- Quando frete permanece NO-DATA (degrau 4 "sem dados", cenário 5): separar o estado. Custo
  nil = **BLOCKING estrutural** (margem incognoscível a qualquer preço — mantém). Frete nil =
  **condicional ao segmento** — o segmento <79 resolve normalmente (frete não consultado);
  só o segmento ≥79 fica bloqueado. Novo campo em `SolveResult`:
  `FreteDesconhecido bool` distinto de `Desconhecidos` (custo). O `R$79` deixa de ser
  constante hardcoded: o limiar de segmento vem do sinal `mandatory_free_shipping` do item
  (resolver), com fallback política configurável — **nunca** literal `7900`.
- Remover `const thresholdCents = 7900` e `taxaFixaLimiar = 79`; o limiar torna-se parâmetro
  resolvido (por item ou política), passado ao domínio.

**Layer 2 — transport (`calc_handler.go:344-347`): rótulo BLOCKING dressed as
UNREACHABLE_TARGET.** Hoje `if !out.Result.Reached { body["code"]="UNREACHABLE_TARGET" }`
independentemente da causa. **Fix:** ramificar na causa antes de rotular:

```
if len(desconhecidos) > 0            → code = "DADOS_INCOMPLETOS" (ou "SEM_FRETE"/"SEM_CUSTO" específico)
else if !reached && ceiling_pct != "" → code = "UNREACHABLE_TARGET"   (legítimo, com teto)
else if reached                       → preço (com labels de fonte/degrau)
```

`desconhecidos` não-vazio **nunca** vira UNREACHABLE. O payload passa a carregar os
carimbos de fonte/degrau/ESTIMATIVA da `TariffResolution` (novo bloco `tarifa: {comissao,
frete}` com fonte/degrau/data/flags) — a UI precisa deles para os rótulos (§6).
**OpenAPI + sdk-runtime** atualizados no mesmo commit (o response de `/pricing/solve` e
`/pricing/decompose` ganha o bloco `tarifa`).

**Layer 3 — FE (`SolverPanel.tsx:44-46,80-88`): só trata SEM_CUSTO, ceiling em branco.**
Hoje: `semCusto` (SEM_CUSTO) e `!reached` → "Alvo inatingível… ceiling_pct%". Um
`desconhecidos:["frete"]` cai no ramo `!reached` e renderiza "inatingível" com **ceiling em
branco** (o domínio devolve `CeilingPct=""` no caso blocking). **Fix:** novo ramo para
`code=DADOS_INCOMPLETOS`/`SEM_FRETE` com **guidance acionável**, distinto do teto:

```
SEM_FRETE / frete desconhecido →
  "Sem dados de frete para este produto. Cadastre dimensões (peso, altura, largura,
   comprimento) OU vincule um anúncio ML para cotar o frete."  + link p/ ação
DADOS_INCOMPLETOS (custo) → banner SEM_CUSTO existente
UNREACHABLE_TARGET (com ceiling) → "melhor margem possível: X%"  (só quando ceiling != "")
```

Renderizar os badges de fonte/degrau/data e a pill ESTIMATIVA (degrau 4) no resultado do
preço. Usar as primitivas `UnknownValue`/`FreshnessIndicator` (`packages/ui`,
mission.md §Runtime Contract) para o estado desconhecido — nunca R$0/verde enganoso.

---

## 5. Fluxo início-de-vida do produto

### 5.1 Quando cotar

| Gatilho | Degrau alvo | Persiste? |
|---|---|---|
| Import xlsx | **NÃO cota automaticamente** (pós-demo: opcional em lote) | — |
| Botão "cotar agora" (UI produto/simulador) | 2 ou 3 conforme vínculo/EAN/dims | sim, linha COTACAO |
| Abertura do Simulador com produto | resolve da escada; cota on-demand se cache stale | reusa cache ou re-cota |
| Ingestão de pedido (M-08) | **grava VENDA** (não é cotação) | sim, linha VENDA |

**Decisão:** cotação é **on-demand** (não no import — evita rajada de N chamadas ML no
upload, rate-limit, e produtos que nunca serão vendidos). O import só popula identidade/dims
se vierem no xlsx (não vêm — só EAN opcional). O "início de vida" efetivo é a **primeira
abertura no Simulador** ou o **botão explícito "cotar agora"**.

### 5.2 Cache & TTL

Cotações são cacheadas como linhas `fonte=COTACAO` em `product_tariff_history`
(`fetched_at` + `quote_ttl_seconds`). Resolução reusa a cotação mais recente **dentro do
TTL**; expirada → re-cota on-demand (ou degrada rotulando se ML indisponível — cenário 7).
**TTL proposto (ratificação §10):** 7 dias para comissão (taxa muda raramente), 24h para
frete (mais volátil). Valores configuráveis por tenant.

### 5.3 Staleness: VENDA-velha vs COTAÇÃO-fresca (cenário 10)

O degrau é **preferência**, mas staleness **rebaixa**:
- VENDA dentro de `STALE_VENDA_DAYS` (proposto: 90 dias) → degrau 1 autoritativo.
- VENDA mais velha que o limiar → **flag `Stale=true`**; o resolver **prefere** uma COTAÇÃO
  fresca (degrau 2/3) se disponível, e a UI mostra **ambas** com data ("Última venda:
  R$X em 12/2025 · Cotação atual: R$Y"). Nunca se descarta silenciosamente a VENDA velha —
  ela vira contexto histórico rotulado, não a fonte ativa.
- Se não há cotação fresca possível (sem vínculo, sem EAN, sem dims) → VENDA velha
  permanece, **rotulada stale** com data, para o operador julgar.

**Regra:** "fonte mais fresca e autoritativa vence, com data sempre visível". A escolha do
resolver é determinística (limiares configuráveis), a exibição é honesta (mostra o conflito).

---

## 6. UI surface

### 6.1 Onde comissão/frete aparecem

- **Simulador (M-07, `/precos`, `SolverPanel.tsx` + painel decompose):** resultado do
  cálculo mostra comissão e frete resolvidos, cada um com **badge de fonte/degrau/data** e
  pill **ESTIMATIVA** quando degrau 4. Fonte primária da demo.
- **Produto Detalhe (M-06, `/catalogo/produtos/:id`):** card "Tarifas" com comissão + frete
  atuais, fonte/degrau/data, histórico resumido (variação no tempo — operator §1a), botão
  "cotar agora".
- **Vínculos (M-04):** não exibe tarifa (fora de escopo); mas um vínculo RESOLVED
  **habilita** degrau 2 para o produto.

### 6.2 Rótulos de fonte

Cada valor de tarifa na UI carrega: `fonte` (VENDA/COTAÇÃO/CATEGORIA/PADRÃO), `degrau`
(1-4), `data`, e badge de estado: `ESTIMATIVA` (degrau 4), `PREDITA` (categoria via
domain_discovery não confirmada), `STALE` (VENDA velha). Frete NO-DATA → `UnknownValue`
("—") + instrução acionável. Copy nunca promete dado que não tem (ADR-06/ADR-17).

### 6.3 Desambiguação de categoria top-3 (cenário 4)

Quando degrau 3 usa `domain_discovery` e retorna **top-3 categorias sem score**, a UI
apresenta as 3 opções para **confirmação humana** (produto detalhe ou um passo do "cotar
agora"). Até confirmar: categoria = `PREDITA`, cotação rotulada. Ao confirmar: grava a
escolha (não como fato ML, mas como **decisão do operador** — flag própria), `Predita=false`
dali em diante. **Nunca auto-grava** a predição como categoria confirmada.

### 6.4 Config degrau 4 (defaults editáveis)

Nova config **`pricing_tariff_defaults`** (por tenant/installation), **não** no CalcProfile
(que é regime fiscal, migration 0055, sem comissão). Campos:
`comissao_classico_pct` (seed 13,00), `comissao_premium_pct` (seed 16,00),
`frete_estimativa_amount` (nullable — null = política "sem dados"),
`frete_policy` ('estimativa'|'sem_dados'). Editável via nova tela/endpoint
`PUT /pricing/tariff-defaults`. Os 13%/16% são **seed inicial editável**, jamais constante
de código (restrição §7). Migration para a config: 0069 (ou junto do bloco alocado).

---

## 7. Cenários de falha (12 do brief, comportamento desenhado)

| # | Cenário | Comportamento desenhado |
|---|---|---|
| 1 | xlsx sem EAN | degrau 3 via descrição pt-BR (`domain_discovery`) → categoria PREDITA, rotulada; sem descrição útil → degrau 4 |
| 2 | EAN não acha catálogo ML | fallback descrição (`domain_discovery`) ou degrau 4; nunca erro mudo |
| 3 | catalog product sem `buy_box_winner` (null) | `domain_discovery` no nome do catálogo → categoria PREDITA |
| 4 | `domain_discovery` prediz categoria errada | top-3 na UI p/ humano (§6.3); nunca auto-grava; `PREDITA` até confirmar |
| 5 | sem dims/peso | frete degrau 3 impossível → frete = NO-DATA; UI "sem dados de frete" + instrução (cadastrar dims OU vincular anúncio); solver reporta `SEM_FRETE` segment-aware (não erro mudo, §4.4 L3) |
| 6 | `list_cost=0` por call under-spec | tratar como NO-DATA (`frete_amount=NULL`), nunca custo 0 (cenário do `mandatory_free_shipping` ausente) |
| 7 | token 401 / 429 / expirado | backoff exponencial + jitter, cache de cotações (TTL §5.2), degrada p/ degrau inferior **rotulado**; render nunca quebra |
| 8 | preço < limiar frete-grátis | **não** assume frete 0 por R$79 hardcoded (removido); sinal por item = `mandatory_free_shipping`; sem item → política explícita configurável rotulada |
| 9 | vínculo REVIEW (não RESOLVED) | degrau 2 **bloqueado**; usa degrau 3/4; nunca cota por item de vínculo não resolvido |
| 10 | VENDA velha (regra ML mudou) | histórico com data; staleness policy (§5.3): VENDA >90d flag STALE, prefere COTAÇÃO fresca, mostra ambas datadas |
| 11 | quote parcial (comissão OK, frete falhou) | grava linha com `frete_amount=NULL`; UI frete "—" honesto; comissão resolvida normalmente |
| 12 | multi-installation | `tenant_id` + `installation_id` em **toda** linha (§3.1) e toda query do resolver |

---

## 8. Faseamento vs demo 2026-07-20 (T-1)

### 8.1 DEMO-CRÍTICO (mínimo para segunda)

**Objetivo honesto:** o solver **produz um preço** usando degrau 4 (defaults rotulados
ESTIMATIVA) em vez de bloquear com `desconhecidos:["frete"]`; e não mente rótulos.

Escopo mínimo viável (realista para T-1, **sem** a escada completa):
1. **Layer 2 (transport) + Layer 3 (FE)** do FINDING-P7-SOLVER — ramificar rótulo
   BLOCKING≠UNREACHABLE + FE com guidance acionável. **Barato, alto valor** (o solver deixa
   de parecer quebrado).
2. **Degrau 4 mínimo:** config `pricing_tariff_defaults` (comissão 13%/16% + frete política)
   + adapter tier-4 + `TariffResolver` que resolve **só** degrau 4 (comissão default + frete
   estimativa/NO-DATA rotulado). Isso já faz o solver produzir preço com margem rotulada.
3. **Layer 1 (domínio) parcial:** injetar frete resolvido (degrau 4) upstream para o probe
   não forçar falso-desconhecido; remover `R$150`. Remoção do `R$79`/segment-refactor pode
   ser parcial (rotular, não necessariamente refatorar o limiar por-item na demo).

**Corte honesto T-1:** a escada completa (degraus 1-3, tabela histórica, persistência VENDA,
cotação ML live) **não** cabe até segunda com Global Maximum. O demo-crítico entrega "solver
funciona com estimativa rotulada e honesta", não "tarifa real por produto". Runbook deve
dizer isso ao cliente (o degrau 4 é claramente ESTIMATIVA).

### 8.2 PÓS-DEMO (CORTÁVEL da demo)

- **[CORTÁVEL]** Degrau 2 (COTAÇÃO-ANÚNCIO) — ligar FeeQuoteReader+GetFreeShippingCost(item)
  ao resolver (readers prontos; só wiring + cache).
- **[CORTÁVEL]** Degrau 3 (COTAÇÃO-MATCH) — extensão adapter (dimensions + EAN→catálogo→
  categoria), domain_discovery, desambiguação top-3.
- **[CORTÁVEL]** Tabela `product_tariff_history` + persistência VENDA na ingestão de pedidos
  (acoplado ao M-08 shipfix).
- **[CORTÁVEL]** Staleness policy completa, TTL configurável, histórico de variação na UI.
- **[CORTÁVEL]** Refactor completo do limiar de segmento por `mandatory_free_shipping`
  por-item (Layer 1 total).

### 8.3 Realismo T-1

Mesmo o demo-crítico exige: 1 migration (config defaults), 1 port + 2 adapters (tier-4 +
resolver mínimo), mudança de transport (OpenAPI+SDK juntos), mudança de FE, e o fix de
domínio. É **1 chip apertado** (pricing) + toque cirúrgico de transport/FE. Se o tempo
comer, o **piso absoluto** é Layer 2+3 (rótulos honestos) + degrau-4-comissão-só (frete
NO-DATA rotulado) — o solver produz preço com margem, frete "—" explicado. Isso é
defensável na demo e é o corte de emergência.

---

## 9. Fatiamento em chips & collision matrix

### 9.1 Chips propostos

| Chip | Escopo | Seam/dono exclusivo | Fase |
|---|---|---|---|
| **CHIP-T1 pricing-solver-tier4** | TariffResolver port + adapter tier-4 + config defaults + fix Layer 1/2 domínio+transport + OpenAPI/SDK `tarifa` block | `modules/pricing/**`, `sdk-runtime/src/pricing.ts`, migration config | DEMO-CRÍTICO |
| **CHIP-T1-FE solver-labels** | SolverPanel guidance acionável + badges fonte/degrau/ESTIMATIVA | `apps/web/src/pages/precos/**` | DEMO-CRÍTICO |
| **CHIP-T2 connectors-quote-ext** | estende FreeShippingQuery(dimensions) + EAN→catálogo→categoria reader + domain_discovery | `modules/connectors/adapters/mercado_livre/**` + ports | PÓS-DEMO |
| **CHIP-T3 tariff-history** | migration `product_tariff_history` + repo + adapter tier-1 (VENDA read) + persistência VENDA na ingestão | `modules/pricing/adapters/tariffhistory` + write hook orders | PÓS-DEMO |
| **CHIP-T4 tariff-ui** | Produto Detalhe card Tarifas + histórico + desambiguação top-3 + config defaults UI | `apps/web/src/pages/produto/**`, `precos` config | PÓS-DEMO |

Degraus 2/3 no resolver ligam-se em CHIP-T3/T2 (wiring do resolver aos readers).

### 9.2 Collision matrix vs trabalho aberto

| Seam | CHIP-T* que toca | Trabalho aberto que colide | Resolução |
|---|---|---|---|
| `connectors/adapters/mercado_livre/shipping_reader.go` | CHIP-T2 (estende FreeShippingQuery + getFreeShippingCost) | **M-08 CHIP-M08-SHIPFIX** (mexe em `getShipmentInfo`/decode shipment + header x-format-new, mesmo arquivo) | **Serializar** — CHIP-T2 só após M-08 shipfix merged (evita colisão no mesmo arquivo) |
| orders write/ingestion path | CHIP-T3 (grava linha VENDA na ingestão) | **M-08 SHIPFIX** (bucket pedidos + decode senders[].cost) | **Serializar** — persistência VENDA depende do decode senders estável do shipfix; ou desacoplar via job backfill |
| `apps/web/src/pages` anuncios | — | **M-05 anuncios-sinais** (FE `/anuncios` + `listings.ts`) | **Sem colisão** — CHIP-T1-FE/T4 tocam `/precos` e `/produto`, disjuntos de `/anuncios` |
| `sdk-runtime/src/pricing.ts` | CHIP-T1 | M-07 já é dono de `pricing.ts` | Sem colisão externa (mesmo milestone lineage); aditivo |
| `sdk-runtime/src/index.ts` barrel | CHIP-T1 (se novo arquivo SDK) | seam hub-owned | Export aditivo 1 linha, hub adjudica |
| migrations | CHIP-T1 (0068/0069), CHIP-T3 (0068 tabela) | reserva hub 0070-0074 | **Hub aloca bloco** — 0068-0069 livres (topo atual 0067 verificado); NÃO auto-atribuir |
| `contracts/api/*.openapi.yaml` | CHIP-T1 (response tarifa block) | contract-lock hub-owned | REQUEST contract-lock ao hub |
| `root.go` composition | CHIP-T1/T3 (registro adapters/resolver) | seam hub-owned | Registro próprio aditivo, conflito = hub |

**Regra de ordem:** CHIP-T1 + CHIP-T1-FE (demo-crítico) são **independentes** do M-08 shipfix
(tocam pricing/precos, não connectors/orders) → podem correr **em paralelo** com o que
estiver aberto. CHIP-T2 e CHIP-T3 **serializam após M-08 shipfix** por colisão de arquivo/
dado. CHIP-T4 é FE disjunto (paralelo).

---

## 10. Riscos & aberto

### 10.1 Não verificável sem live-drive / doc adicional

- **Comportamento exato de `listing_prices` sem params logísticos:** o `fixed_fee` exato
  exige `logistic_type`/`shipping_mode`/`billable_weight` (brief §3), que o
  `FeeQuoteInput` atual **não carrega** (`capability.go:280` — verificado). O
  `percentage_fee` (comissão %) funciona sem eles; o `fixed_fee` pode vir aproximado.
  Impacto: comissão % confiável; taxa fixa aproximada até estender o input. Verificar em
  live-drive (preflight M-02).
- **`domain_discovery` sem score de confiança** (brief §3): desambiguação humana é
  **obrigatória**; não há como auto-rankear com confiança. Assumido, não contornável.
- **Hook de escrita na ingestão de pedidos (M-08):** o caminho de projeção de pedidos está
  em fluxo (CHIP-M08-SHIPFIX). Não verifiquei um ponto de extensão estável para emitir a
  linha VENDA — depende do estado pós-shipfix. Coordenar com M-08 ou usar job backfill
  desacoplado.
- **`senders[0].cost` = frete do vendedor?** Assumido pelo brief (degrau 1 frete). O
  mapeamento existe (`shipping_reader.go:167`) mas a semântica exata (custo do vendedor vs
  bruto vs desconto) merece confirmação em live-drive antes de persistir como "frete real".

### 10.2 Decisões que precisam de ruling do operador

1. **`marketplace_fee_schedules` como degrau 3.5?** Recomendação: manter como fallback de
   comissão por categoria (rotulado `CATEGORIA`) entre degrau 3 e 4. Ratificar ou colapsar.
2. **TTL das cotações:** proposto 7d comissão / 24h frete. Confirmar valores.
3. **`STALE_VENDA_DAYS`:** proposto 90 dias. Confirmar.
4. **Política de frete degrau 4:** estimativa numérica configurável **ou** "sem dados"?
   Recomendação: default "sem dados" (mais honesto); operador pode configurar estimativa.
5. **Frete degrau 4 na demo:** mostrar "—" (sem dados) é aceitável na frente do cliente, ou
   precisa de um número estimado rotulado? Afeta o corte de emergência §8.3.
6. **Bloco de migração 0068-0069:** hub deve alocar formalmente (reserva 0070-0074 é hub-only).
7. **Corte T-1:** confirmar que o demo-crítico §8.1 (solver com degrau 4 rotulado) é
   aceitável como história da demo, com a escada real posicionada como pós-demo no runbook.

### 10.3 Restrições honradas (checklist)

- ZERO escrita ML — todos os endpoints da escada são GET (`listing_prices`,
  `shipping_options/free`, `products/search`, `products/{id}`, `domain_discovery`,
  `shipments/{id}/costs`). `MPC_PROVIDER_WRITES_ENABLED` unset. ✔
- ADR-17 — desconhecido nunca vira 0/default; frete NO-DATA = nil + "—"; ESTIMATIVA e
  PREDITA rotulados. ✔
- Payloads de provider morrem nos adapters (connectors); ports normalizados IC-06. ✔
- `tenant_id` (+ `installation_id`) em toda linha e toda query. ✔
- OpenAPI + `sdk-runtime` no mesmo commit (response `tarifa` block). ✔
- Nunca hardcodar tiers de taxa fixa nem limiar R$79 — 13%/16% viram config seed editável;
  R$79/R$150 saem do domínio. ✔
- Codex morto até 2026-07-25 → lane Claude-only (planner Opus, implementers sonnet). ✔
- Chips nunca sobem servidor; dep change = REQUEST ao hub. ✔
```

## 11. ADENDO (ratificado D-84, 2026-07-19) — categoria, catmap e T2-MIN

**Contexto:** live-drive da sonda catalog-match (D-83) provou comissão = f(categoria, tipo, faixa-preço)
com valores reais ≠ default tier-4 (11%/14% vs 13%/16% em Kit de Chaves @R$699). Operator elevou a
prioridade: feature principal da demo = preço concorrente + margem com comissão real + sugestão de
produto. Pesquisa de mercado (Anymarket/Hub2b/uappi, category-mapping-research.md) confirmou padrão
universal de-para categoria-fonte→categoria-ML com cascata por grupo.

**Rulings:**
1. **CHIP-T2-MIN (PRÉ-DEMO):** degrau 3 mínimo — adapter fino ligando a máquina catalog-match já
   mergeada (EAN+título→categoria→listing_prices) ao port TariffResolver do CHIP-T1. Mesmo port,
   mesma cadeia, carimbo fonte=COTACAO degrau=3. Serializa após merge do T1 (colisão pricing).
   SEM: persistência history (fica T3), frete por dims, fix buy_box (fica T2 full; "mojibake"
   RETRATADO D-85 — falso positivo de decode cp1252 no console Windows, bytes limpos).
2. **CategoryResolver = port próprio** desde o T2-MIN: fonte 1 = EAN+título (catalog-match).
   Catmap entra pós-demo como fonte 2 atrás do MESMO port. Consumidor nunca vê a origem, só o
   carimbo {categoria, fonte: EAN-CATALOGO|TITULO|GRUPO-DEPARA|MANUAL, confianca, data}.
3. **Import captura grupo/descrgrupo JÁ** (colunas nullable + parser; aditivo). Dado guardado
   desde agora; de-para consome depois.
4. **CHIP-CATMAP (PÓS-DEMO):** árvore ML local via /sites/MLB/categories/all (cache MD5-diff,
   mecanismo oficial); tabela de-para grupo→categoria-folha-ML com sugestão automática (votos de
   EANs do grupo + domain_discovery limit=3) e confirmação humana em lote (padrão ACCEPT/REVIEW
   M-04). Cobre os ~0,8% sem EAN + robustez/novos marketplaces.
5. **NCM DESCARTADO como âncora de categoria** — nenhum recurso ML/mercado sanciona NCM→categoria
   (pesquisa Q3, resultado negativo documentado). NCM segue só como dado fiscal.
6. **Fee = cache TTL curto via listing_prices, NUNCA tabela estática longa** — sem endpoint bulk;
   mudança estrutural oficial de 02/03 (provisoes/comissao-por-vender) muda FORMA da cobrança, não
   só números. TTL 7d comissão (D-78) mantido. fee_schedules categoria permanece só degrau 3.5.
7. **Numeração de migration (AMENDA D-87, supersede §6.3/§8/reservas acima):** reserva de número
   para migration FUTURA está CANCELADA — reservar cria hazard de apply fora de ordem (fresh DB =
   ordem lexical; DB existente = ordem de chegada). Regra: número = próximo-livre concedido pelo
   hub no momento do GRANT, sequencial, sem buracos. Alocação corrente: 0068 = T1
   (pricing_tariff_defaults), 0069 = GRUPO-IMPORT (grupo/descrgrupo), T3 pega next-free ao existir.
   Menções a "0069 reservado T3" e "reserva hub 0070-0074" neste doc estão obsoletas.
