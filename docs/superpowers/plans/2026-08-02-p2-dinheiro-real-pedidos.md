# P2 — Fios cortados em /pedidos — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fechar os seis campos que `/pedidos` entrega vazios porque o dado do provider nunca foi declarado ou porque o tipo torna `NULL` inalcançável — sem tocar em margem, que já funciona.

**Architecture:** Metade "fio cortado" do M-06, re-ancorada por medição ao vivo em 2026-08-02. `orders/**` continua superfície exclusiva do M-06 (ADR-024); este plano não abre trilha nova. Toda tarefa é a mesma classe de defeito que o P1 fechou em `/anuncios`: o provider manda a chave, o DTO não declara, `encoding/json` descarta em silêncio, a coluna nasce vazia.

**Tech Stack:** Go 1.x (`apps/server_core`), pgx + Postgres, OpenAPI 3 + `packages/sdk-runtime`, React/TS (`apps/web`).

---

## Comunicação com o hub — leia antes de tudo

**Você pode e deve falar comigo.** A sessão anterior terminou em silêncio e o custo foi perseguir erro que não era dela.

| situação | evento | quando |
|---|---|---|
| dúvida que muda o que você vai escrever | `REQUEST` | **antes** de escrever, não depois |
| premissa deste plano bateu de frente com o repo | `ESCALATION` | imediatamente, com `file:line` |
| dependência nova (pacote, migration, `.env`, servidor) | `REQUEST` | antes de criar ou instalar qualquer coisa |
| tarefa terminou e commitou | `COMMITTED` | a cada tarefa |
| plano inteiro terminou | `CLOSED` | uma vez, no fim, com a evidência |
| travou | `BLOCKED` | assim que travar, não depois de 6 rodadas |

**Regra dura:** se este plano afirma algo sobre o repositório e você mede o contrário, **pare e mande `ESCALATION`**. Não conserte o plano em silêncio. Este documento é a segunda versão justamente porque a primeira acreditou num design que estava errado — o design dizia "margem não existe" e a medição mostrou margem em 33/38.

**Nunca sem `REQUEST` aprovado:** subir servidor, tocar `:8080`, criar `.env`, instalar dependência, `git push`, escrever ao vivo no Mercado Livre. Leitura ao vivo do provider é permitida.

---

## O que já funciona — medido, não suposto

Rodei a API contra os 38 pedidos reais do dev stack. **Não reimplemente nada disto:**

| fato | medida | onde |
|---|---|---|
| margem por pedido | **33/38 não-nulo** | `enrich_service.go:391` → `domain.BuildProfitability` |
| os 5 nulos | são exatamente os 5 `SEM_VINCULO` | correto por ADR-17 — sem custo não há margem |
| comissão × quantidade | correto | `enrich_service.go:444` (`sale_fee` é POR UNIDADE) |
| frete do vendedor | 38/38 | `/shipments/{id}/costs` → `senders[].cost`, não `lead_time` |
| custo do ERP na data da venda | `TGFCUS.CUSSEMICM` com `DTATUAL <= data do pedido` | margem histórica não se reescreve |
| `taxa_fixa` / `tarifa_full` | desconhecidos 38/38 | honesto — estão dentro do `sale_fee`, separar exige engine |

**Existe um `ComputeDecomposition` neste plano? Não. Se você sentir vontade de escrever um, pare.** `BuildProfitability` já existe, já é chamado e já acerta. Um segundo calculador seria dois produtores para o mesmo número.

---

## O que está quebrado — os seis fios

| campo | medida | causa rastreada |
|---|---|---|
| `provider_status_detail` | `""` em **38/38** | tipo é `string` (`order.go:53`) — `NULL` inalcançável |
| motivo de cancelamento | vazio nos **7 cancelados** | DTO não declara `cancel_detail`; lê `status_detail`, que vem `null` |
| `logistic_type` | **0/38** no banco | DTO de envio não declara `logistic` |
| `tracking_method` | **0/38** no banco | DTO de envio não declara `tracking_method` |
| `currency` | **null 38/38** | **nenhum produtor** — sem coluna, sem projeção, sem atribuição |
| `fulfillment` | **null 38/38** | **nenhum produtor** |

---

## Fora de escopo — decidido pelo operador em 2026-08-02

### O imposto e o DIFAL estão errados. **Não conserte neste plano.**

Medido: `pricing_calc_profiles` tem **0 linhas**. `calc_repository.go:49` responde `pgx.ErrNoRows` com `NewDefaultCalcProfile()` — SIMPLES 4%, `origem: "default"`, DIFAL desligado. Então:

- `imposto: 12,00` num pedido de R$ 299,90 é 4% de um regime que **ninguém configurou** (o operador ratificou Regime Normal, não SIMPLES).
- `difal: 0,00` é o zero explícito de um switch que nunca foi ligado — num pedido interestadual para a Bahia, com 27 alíquotas por UF já sentadas em `pricing_difal_rates`, nenhuma consultada.

A aritmética da margem está impecável e duas das seis entradas são invenção, então a margem exibida está **otimista e errada**.

**Isso é o P2.b**, por decisão do operador: perfil fiscal, regime, tela de configuração e a correção do default. Um plano só para isso.

**Você está proibido de tocar em `pricingtax/`, `calc_repository.go` ou no perfil fiscal nesta fatia.** Se sua mudança faz o `imposto` mudar de valor, ela está fora do escopo — reverta e mande `REQUEST`. Isso vale inclusive se parecer uma melhoria óbvia: o operador quer a correção fiscal isolada, para poder medi-la sozinha.

### Demais itens

| item | dono |
|---|---|
| perfil fiscal, regime, tela de configuração, default mentiroso | **P2.b** |
| `CODEMP` fixo em 1 no leitor de custo (`orders_adapters.go:53`), enquanto vendável é `CODEMP(1,2)` | **dívida D-17**, registrada na Task 1, sem conserto aqui |
| `nf_state` null 38/38 | **não é bug** — é `'linked'` só com evento Sankhya `evidence_state='exact'`, e não há nenhum. Reconciliação de NF-e, outra fatia |
| colunas `net_amount` / `margin_pct` / `decomposition` (0/38) | **ficam vazias.** O cálculo no caminho de leitura já está certo e o custo já é datado; persistir criaria um segundo lugar para o mesmo número, com risco de divergir. YAGNI |
| backfill 12 meses + scheduler 5min | M-06 F-01 — plano próprio |
| DIFAL como motor de regras | P3 |

---

## Decisão de máximo global — por que o `COALESCE` de orders FICA

O escritor de `listings` passou a **declarar as colunas que possui** (ADR-C2, P1). O de `orders` faz o oposto: `COALESCE(EXCLUDED.x, tabela.x)` em `buyer_nickname`, `pack_id`, `date_last_updated_ml` e nas colunas fiscais do comprador (`order_repo.go:667-680`).

**Não unifique.** Parecem o mesmo problema, não são:

- ADR-C2 protege coluna de **outro** escritor — o sweep de anúncios apagava `quality_score`, que era do dono.
- O `COALESCE` de orders protege contra **snapshot parcial do mesmo escritor** — o ML só expõe dado fiscal do comprador depois que existe etiqueta. Escrever `NULL` apagaria fato já aprendido.

Uniformidade não é máximo global; corretude é. Cada escritor declara qual regra usa e por quê.

**Dívida que isso cria, registrada e não escondida:** com `COALESCE` nós **retemos PII para sempre** se o ML remover o dado do comprador — o provider não distingue "ainda não te contei" de "apaguei". É política de expurgo, decisão do operador. **Não** conserte aqui.

---

## Estrutura de arquivos

| arquivo | responsabilidade nesta fatia |
|---|---|
| `internal/modules/orders/domain/order.go` | dois campos viram `*string` |
| `.../mercado_livre/order_ingest_reader.go` | declarar `cancel_detail`, `currency_id` |
| `.../mercado_livre/shipment_ingest_reader.go` | declarar `logistic`, `tracking_method` |
| `.../mercado_livre/testdata/*.json` | fixtures redigidas para o detector `rawkeys` |
| `internal/modules/orders/adapters/postgres/order_repo.go` | projeção e scan dos campos novos |
| `migrations/0093_orders_currency_and_shipment_wiring.sql` | coluna `currency` |

---

### Task 1: Verificar as premissas — e registrar a dívida D-17

Esta tarefa não produz feature. Produz o direito de confiar no resto do plano.

**Files:**
- Create: `.mnfs/MIS-007-ml-sync/M-06-orders-backfill-decomposition/evidence/p2-premise-check.md`
- Modify: `.mnfs/HARNESS-DEBTS.md`

- [ ] **Step 1: Confirmar que a margem NÃO é o problema**

```bash
docker exec marketplace-central-postgres-1 psql -U marketplace -d marketplace_central -tAc "select installation_id from orders_marketplace_orders limit 1;"
```

Com esse valor:

```bash
curl -s "http://localhost:8080/orders?limit=50&installation_id=<INST>" -o /tmp/p2-orders.json
```

Se o servidor não estiver de pé, **mande `REQUEST`** — chip não sobe servidor.

Conte quantos pedidos têm `margem_pct` não-nulo. Esperado: **33 de 38** (os números absolutos podem crescer; a proporção "todos menos os `SEM_VINCULO`" não pode mudar).

**Se margem vier nula em todos, PARE e mande `ESCALATION`.** Significa que algo regrediu entre 2026-08-02 e agora, e este plano inteiro está apontado para o alvo errado.

- [ ] **Step 2: Confirmar os seis buracos**

Do mesmo JSON, conte: `currency` null, `fulfillment` null, `provider_status_detail == ""`. Esperado: **38/38 nos três**.

```bash
docker exec marketplace-central-postgres-1 psql -U marketplace -d marketplace_central -tAc "select count(*)||'|'||count(logistic_type)||'|'||count(tracking_method) from order_shipments;"
```
Esperado: `N|0|0`.

- [ ] **Step 3: Confirmar que o perfil fiscal continua ausente**

```bash
docker exec marketplace-central-postgres-1 psql -U marketplace -d marketplace_central -tAc "select count(*) from pricing_calc_profiles;"
```
Esperado: `0`.

Isto **não é para você consertar** — é para registrar o estado em que o P2.b vai pegar o sistema. Se vier `1`, alguém configurou o perfil e o P2.b já começou; mande `REQUEST` antes de seguir.

- [ ] **Step 4: Registrar a dívida D-17**

Acrescente a `.mnfs/HARNESS-DEBTS.md`:

```markdown
## D-17 — CODEMP fixo em 1 no leitor de custo de pedidos

`internal/composition/orders_adapters.go:53` fixa `CompanyID: 1` na consulta
`TGFCUS`. O predicado de vendável ratificado cobre `CODEMP(1,2)`. Produto cujo
custo só exista na empresa 2 devolve linha nenhuma → custo desconhecido →
margem some, sem nada na tela explicando por quê.

Nenhum dos 38 pedidos medidos em 2026-08-02 caiu nisso. Bomba armada, não
disparada. Dono: M-06. Registrada por decisão do operador em 2026-08-02.
```

- [ ] **Step 5: Escrever a evidência e commitar**

Grave as saídas verbatim no arquivo de evidência, com data. Qualquer divergência ⇒ `PREMISSA REPROVADA` no topo e `ESCALATION` antes de seguir.

```bash
git add .mnfs/
git commit -m "evidence(M-06/P2): medicao de premissas + divida D-17 CODEMP"
```

---

### Task 2: `NULL` volta a ser alcançável nos dois campos de status

`ProviderStatusDetail` e `CancellationDetail` são `string`. Desconhecido vira `""`, e `""` é indistinguível de "o provider mandou vazio". É ADR-17 ferido **no tipo**, antes de qualquer código rodar — e é por isso que a coluna aparece "38/38 preenchida" contando strings vazias.

**Files:**
- Modify: `apps/server_core/internal/modules/orders/domain/order.go:53-54`
- Modify: `apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go` (upsert + `scanReadModel` em `:253`)
- Test: `apps/server_core/internal/modules/orders/adapters/postgres/order_repo_integration_test.go`

- [ ] **Step 1: Write the failing test**

```go
//go:build integration

func TestUpsertOrderDistinguishesAbsentDetailFromEmptyDetail(t *testing.T) {
	testpostgres.SkipWithoutTarget(t)
	ctx := context.Background()
	repo, installationID := newOrderIntegrationRepository(t)

	empty := ""
	absent := newIntegrationOrder(t, installationID, "ORD-ABSENT")
	absent.CancellationDetail = nil
	present := newIntegrationOrder(t, installationID, "ORD-EMPTY")
	present.CancellationDetail = &empty

	for _, o := range []domain.MarketplaceOrder{absent, present} {
		if _, err := repo.UpsertOrder(ctx, o); err != nil {
			t.Fatalf("UpsertOrder(%s): %v", o.ProviderOrderID, err)
		}
	}

	var absentIsNull, presentIsNull bool
	if err := repo.pool.QueryRow(ctx,
		`SELECT
		   (SELECT cancellation_detail IS NULL FROM orders_marketplace_orders WHERE provider_order_id='ORD-ABSENT'),
		   (SELECT cancellation_detail IS NULL FROM orders_marketplace_orders WHERE provider_order_id='ORD-EMPTY')`,
	).Scan(&absentIsNull, &presentIsNull); err != nil {
		t.Fatalf("leitura: %v", err)
	}
	if !absentIsNull {
		t.Fatal("ausente gravou não-NULL: o tipo string tornou NULL inalcançável")
	}
	if presentIsNull {
		t.Fatal("string vazia virou NULL: perdemos a distinção no outro sentido")
	}
}
```

O segundo `if` existe para impedir a correção preguiçosa `NULLIF(x,'')`, que troca um colapso por outro.

- [ ] **Step 2: Run test to verify it fails**

A lane hermética faz **SKIP silencioso** sem `MPC_TEST_DATABASE_URL`, e SKIP é byte-idêntico a verde. Confirme `--- FAIL`, nunca `--- SKIP`:

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache MPC_TEST_DATABASE_URL='<dsn da lane>' go test -tags integration ./internal/modules/orders/adapters/postgres/ -count=1 -run TestUpsertOrderDistinguishesAbsentDetailFromEmptyDetail -v
```
Esperado: não compila (`cannot use nil as string`). Esse é o vermelho certo.

Sem o DSN da lane, mande `REQUEST`. Não invente banco e **não altere autenticação de container** para contornar.

- [ ] **Step 3: Write minimal implementation**

Em `order.go`:

```go
	// ProviderStatusDetail e CancellationDetail são *string, não string: o
	// provider distingue "não mandou o campo" de "mandou vazio", e string
	// colapsa os dois em "". Foi esse colapso que fez cancellation_detail
	// aparecer preenchido em 38/38 pedidos — inclusive nos 7 cancelados, onde
	// o motivo real do cancelamento é justamente o que faltava.
	ProviderStatusDetail *string `json:"provider_status_detail,omitempty"`
	CancellationDetail   *string `json:"cancellation_detail,omitempty"`
```

No `order_repo.go` os dois parâmetros vão direto (pgx converte `*string` nil em `NULL`). **Remova qualquer `NULLIF($n,'')`** aplicado a eles. Em `scanReadModel` (`:253`), troque o scan direto por `pgtype.Text` e só atribua quando `Valid`, no mesmo molde de `buyerNickname` logo abaixo.

O compilador aponta cada chamador. Em cada um: ausente ⇒ `nil`, nunca `&""`.

- [ ] **Step 4: Run test to verify it passes**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache MPC_TEST_DATABASE_URL='<dsn>' go test -tags integration ./internal/modules/orders/... -count=1 && GOCACHE=$(pwd)/.gocache go test ./internal/... -count=1
```

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/orders/
git commit -m "fix(orders): status_detail e cancellation_detail viram *string; NULL volta a ser alcancavel"
```

---

### Task 3: O DTO de pedido declara `cancel_detail` e `currency_id`

Não escreva detector novo. `internal/platform/rawkeys` (P1) existe para isto, e este é o **segundo adapter** a usá-lo — é o que transforma o mecanismo em regra em vez de conserto pontual. Assinatura: `rawkeys.Undeclared(raw json.RawMessage, dto any, ignore []string) ([]string, error)`.

**Files:**
- Create: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/testdata/order_body.json`
- Create: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/order_ingest_rawkeys_test.go`
- Modify: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/order_ingest_reader.go`

- [ ] **Step 1: Write the fixture**

Payload de pedido real com **toda PII substituída por literal redigido**. Comprador é pessoa real: nome, documento, endereço e CEP **não entram**, nem em fixture de teste.

```json
{
  "id": 2000012659424976,
  "status": "cancelled",
  "status_detail": null,
  "date_created": "2026-06-14T09:12:03.000-04:00",
  "date_closed": "2026-06-14T09:40:11.000-04:00",
  "date_last_updated": "2026-06-15T11:02:44.000-04:00",
  "last_updated": "2026-06-15T11:02:44.000-04:00",
  "pack_id": null,
  "currency_id": "BRL",
  "cancel_detail": {
    "group": "respondent",
    "code": "cancel_purchase_by_buyer",
    "description": "Comprador se arrependeu",
    "requested_by": "buyer",
    "date": "2026-06-14T09:39:02.000-04:00",
    "application_id": null
  },
  "buyer": { "id": 111111111, "nickname": "REDACTED" },
  "seller": { "id": 222222222 },
  "order_items": [
    {
      "item": {
        "id": "MLB4735324525",
        "title": "Produto de exemplo",
        "seller_sku": "26909",
        "variation_id": null,
        "variation_attributes": []
      },
      "quantity": 2,
      "unit_price": 69.9,
      "full_unit_price": 69.9,
      "sale_fee": 11.53,
      "listing_type_id": "gold_special"
    }
  ],
  "payments": [
    { "id": 333333333, "status": "refunded", "transaction_amount": 139.8, "total_paid_amount": 139.8 }
  ],
  "shipping": { "id": 44444444444 },
  "tags": ["not_delivered", "test_order"],
  "total_amount": 139.8,
  "paid_amount": 139.8,
  "expiration_date": "2026-06-21T09:12:03.000-04:00",
  "feedback": { "sale": null, "purchase": null },
  "context": { "channel": "marketplace", "site": "MLB", "flows": [] },
  "taxes": { "amount": null, "currency_id": null },
  "coupon": { "id": null, "amount": 0 }
}
```

- [ ] **Step 2: Write the failing test**

```go
package mercadolivre

import (
	"encoding/json"
	"os"
	"testing"

	"marketplace-central/apps/server_core/internal/platform/rawkeys"
)

// orderBodyIgnoredKeys: cada entrada é uma DECISÃO, não um esquecimento.
// Entrada sem motivo escrito devolve exatamente o silêncio que o detector
// existe para quebrar.
var orderBodyIgnoredKeys = []string{
	"context",         // canal/site; já conhecido pela instalação
	"coupon",          // cupom do ML; entra na decomposição só se o operador ratificar
	"expiration_date", // prazo de pagamento; sem consumidor
	"feedback",        // reputação; fora de escopo
	"paid_amount",     // já coberto por payments[].total_paid_amount
	"seller",          // é a nossa própria conta
	"taxes",           // imposto do ML (medido null); o nosso é P2.b
	"total_amount",    // derivado de order_items; não é fonte
}

func TestOrderBodyDeclaresEveryConsumedKey(t *testing.T) {
	raw, err := os.ReadFile("testdata/order_body.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	missing, err := rawkeys.Undeclared(json.RawMessage(raw), mlOrderBody{}, orderBodyIgnoredKeys)
	if err != nil {
		t.Fatalf("Undeclared() error = %v", err)
	}
	if len(missing) != 0 {
		t.Fatalf("mlOrderBody nao declara chaves que o ML manda: %v", missing)
	}
}
```

`mlOrderBody` é um placeholder — **leia `order_ingest_reader.go` e use o nome real do struct de fio.** Não adivinhe.

- [ ] **Step 3: Run test to verify it fails — e LEIA a lista**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/connectors/adapters/mercado_livre/ -count=1 -run TestOrderBodyDeclaresEveryConsumedKey -v
```
Esperado: FAIL nomeando **pelo menos** `cancel_detail` e `currency_id`.

**Se aparecer chave que você não esperava, não a jogue em `orderBodyIgnoredKeys` para ficar verde.** Cada chave é uma decisão: entra no DTO, ou entra na lista **com motivo escrito**. Silenciar sem motivo é reconstruir o defeito dentro do teste que existe para pegá-lo.

- [ ] **Step 4: Write minimal implementation**

```go
	// cancel_detail é a fonte REAL do motivo de cancelamento. O código lia
	// status_detail, que veio null em 7/7 pedidos cancelados enquanto
	// cancel_detail estava presente. Não é dado duplicado: é dado zero.
	CancelDetail *mlOrderCancelDetail `json:"cancel_detail"`
	// CurrencyID alimenta orders.currency, que estava null em 38/38 por não
	// ter produtor nenhum — nem coluna, nem projeção, nem atribuição.
	CurrencyID string `json:"currency_id"`
```

```go
type mlOrderCancelDetail struct {
	Group       string `json:"group"`
	Code        string `json:"code"`
	Description string `json:"description"`
	RequestedBy string `json:"requested_by"`
}
```

No mapeamento para o domínio:

```go
	// Forma canônica "<requested_by>:<code>" — requested_by ∈ {buyer, meli,
	// seller} responde "quem cancelou", que é a pergunta operacional. A
	// description do ML é texto de UI dele e muda sem aviso: não é vocabulário
	// nosso e não entra no domínio (ADR-C4).
	if body.CancelDetail != nil {
		detail := strings.TrimSpace(body.CancelDetail.RequestedBy) + ":" + strings.TrimSpace(body.CancelDetail.Code)
		if detail != ":" {
			order.CancellationDetail = &detail
		}
	}
```

- [ ] **Step 5: Run test to verify it passes**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/connectors/... -count=1
```

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/internal/modules/connectors/adapters/mercado_livre/
git commit -m "fix(mercado_livre): DTO de pedido declara cancel_detail e currency_id"
```

---

### Task 4: O DTO de envio declara `logistic` e `tracking_method`

Mesma classe, mesmo detector, terceiro sítio. As colunas `logistic_type` e `tracking_method` já existem em `order_shipments` e estão **0/38**. `logistic.type` também é a fonte de `fulfillment` na Task 5.

**Files:**
- Create: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/testdata/shipment_body.json`
- Create: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/shipment_ingest_rawkeys_test.go`
- Modify: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/shipment_ingest_reader.go`

- [ ] **Step 1: Fixture + teste vermelho**

Fixture de `GET /shipments/{id}` com `destination` **redigido** (endereço do comprador é PII). Teste no molde exato da Task 3, apontando `rawkeys.Undeclared` para o struct de fio do envio, com sua própria lista `shipmentBodyIgnoredKeys` — **cada entrada com motivo escrito**.

Inclua na fixture, no mínimo:

```json
{
  "id": 44444444444,
  "status": "delivered",
  "substatus": null,
  "tracking_number": "REDACTED",
  "tracking_method": "Normal",
  "logistic": { "type": "drop_off", "mode": "me2", "direction": "forward" },
  "lead_time": { "shipping_method": { "id": 1 }, "cost": 0 },
  "destination": { "shipping_address": { "city": { "name": "REDACTED" }, "state": { "id": "BR-BA" }, "zip_code": "REDACTED" } },
  "site_id": "MLB",
  "date_created": "2026-06-14T09:20:00.000-04:00"
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/connectors/adapters/mercado_livre/ -count=1 -run TestShipmentBodyDeclaresEveryConsumedKey -v
```
Esperado: FAIL nomeando `logistic` e `tracking_method`.

- [ ] **Step 3: Write minimal implementation**

```go
	// logistic.type é o modal (self_service, drop_off, xd_drop_off,
	// fulfillment…) e governa quem guarda o estoque e quem paga o quê.
	// tracking_method é a transportadora. Nenhum dos dois era declarado: 0/38.
	Logistic       *mlShipmentLogistic `json:"logistic"`
	TrackingMethod string              `json:"tracking_method"`
```

```go
type mlShipmentLogistic struct {
	Type string `json:"type"`
}
```

Propague para `LogisticType`/`TrackingMethod` no DTO de envio e no upsert de `order_shipments`, com `nil` quando ausente — nunca `""` (ADR-17).

- [ ] **Step 4: Run test to verify it passes**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/connectors/... -count=1
```

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/connectors/adapters/mercado_livre/
git commit -m "fix(mercado_livre): DTO de envio declara logistic.type e tracking_method"
```

---

### Task 5: `currency` e `fulfillment` ganham produtor

Os dois estão `null` em 38/38 porque **não existe produtor nenhum** — não é NULL vindo de coluna vazia, é campo de contrato que ninguém nunca preencheu. Mesma classe de `quality_score` no P1. Lá a decisão foi deletar do contrato porque não havia fonte. Aqui **há fonte** (`currency_id` na Task 3, `logistic.type` na Task 4), então a decisão do operador é **produzir**.

**Files:**
- Create: `apps/server_core/migrations/0093_orders_currency.sql`
- Modify: `apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go` (upsert, projeção `:52` e `:221`, `scanReadModel` `:253`)
- Test: `apps/server_core/internal/modules/orders/adapters/postgres/order_repo_integration_test.go`

- [ ] **Step 1: A migration**

Confirme antes que `0093` está livre:

```bash
ls apps/server_core/migrations/ | tail -5
```
Se `0093` já existir, use o próximo livre e **evite `0094`/`0095`, reservados pelo M-06**. Colisão de numeração ⇒ `ESCALATION`.

```sql
-- 0093: orders.currency existia no contrato de leitura sem produtor nenhum —
-- null em 38/38. A fonte é o currency_id que o ML já manda em cada pedido e
-- que o DTO passou a declarar. Nullable: pedido antigo não tem a informação e
-- desconhecido nunca vira 'BRL' por conveniência (ADR-17).
ALTER TABLE orders_marketplace_orders ADD COLUMN IF NOT EXISTS currency text;
```

`fulfillment` **não ganha coluna**: é derivado de `order_shipments.logistic_type = 'fulfillment'`, que a Task 4 passou a preencher. Guardar uma cópia criaria dois lugares para o mesmo fato.

- [ ] **Step 2: Write the failing test**

```go
//go:build integration

func TestReadModelExposesCurrencyAndFulfillment(t *testing.T) {
	testpostgres.SkipWithoutTarget(t)
	ctx := context.Background()
	repo, installationID := newOrderIntegrationRepository(t)

	brl := "BRL"
	full := newIntegrationOrder(t, installationID, "ORD-FULL")
	full.Currency = &brl
	if _, err := repo.UpsertOrder(ctx, full); err != nil {
		t.Fatalf("UpsertOrder: %v", err)
	}
	insertShipmentWithLogisticType(t, repo, installationID, "ORD-FULL", "fulfillment")

	drop := newIntegrationOrder(t, installationID, "ORD-DROP")
	drop.Currency = &brl
	if _, err := repo.UpsertOrder(ctx, drop); err != nil {
		t.Fatalf("UpsertOrder: %v", err)
	}
	insertShipmentWithLogisticType(t, repo, installationID, "ORD-DROP", "drop_off")

	models := listReadModels(t, repo, installationID)

	got := models["ORD-FULL"]
	if got.Currency == nil || *got.Currency != "BRL" {
		t.Fatalf("Currency = %v, want BRL — o campo continua sem produtor", got.Currency)
	}
	if got.Fulfillment == nil || *got.Fulfillment != "fulfillment" {
		t.Fatalf("Fulfillment = %v, want fulfillment", got.Fulfillment)
	}

	// O modal que NÃO é fulfillment não pode virar fulfillment, e também não
	// pode virar nil: drop_off é fato conhecido, não ausência de fato.
	other := models["ORD-DROP"]
	if other.Fulfillment == nil || *other.Fulfillment != "drop_off" {
		t.Fatalf("Fulfillment do drop_off = %v, want drop_off", other.Fulfillment)
	}
}
```

O segundo bloco impede a implementação preguiçosa que devolve `nil` para tudo que não é `fulfillment` — o que faria a coluna parecer honesta enquanto apaga um fato conhecido.

- [ ] **Step 3: Run test to verify it fails**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache MPC_TEST_DATABASE_URL='<dsn>' go test -tags integration ./internal/modules/orders/adapters/postgres/ -count=1 -run TestReadModelExposesCurrencyAndFulfillment -v
```
Esperado: FAIL — `Currency = <nil>`. Confirme `FAIL`, não `SKIP`.

- [ ] **Step 4: Write minimal implementation**

Adicione `Currency *string` ao `MarketplaceOrder` do domínio, alimentado pelo `CurrencyID` da Task 3 (vazio ⇒ `nil`). Grave na coluna nova.

Nas **duas** projeções (`:52` e `:221` — as duas existem e as duas precisam) acrescente `o.currency` e o modal do envio:

```sql
       o.currency,
       (SELECT s.logistic_type
          FROM order_shipments s
         WHERE s.tenant_id = $1 AND s.installation_id = $2
           AND s.provider_shipment_id = o.provider_shipment_id) AS fulfillment,
```

Em `scanReadModel`, escaneie os dois como `pgtype.Text` e só atribua quando `Valid` — no molde de `buyerNickname`.

**Cuidado:** `scanReadModel` é compartilhado pelas duas projeções. Acrescentar coluna a uma só faz a outra estourar em runtime, não em compile time. Mude as duas na mesma edição.

- [ ] **Step 5: Run test to verify it passes**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache MPC_TEST_DATABASE_URL='<dsn>' go test -tags integration ./internal/modules/orders/... -count=1 && GOCACHE=$(pwd)/.gocache go test ./internal/... -count=1
```

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/migrations/ apps/server_core/internal/modules/orders/
git commit -m "feat(orders): currency e fulfillment ganham produtor real"
```

---

### Task 6: Contrato, SDK e tela

**Files:**
- Modify: `contracts/api/marketplace-central.openapi.yaml` (`OrderReadModel`)
- Modify: `packages/sdk-runtime/src/index.ts`
- Modify: `apps/web/src/pages/PedidosPage.tsx` + drawer

- [ ] **Step 1: Medir a cobertura real antes de selar**

Reingira os pedidos e conte, no mesmo JSON da Task 1: `currency` não-nulo, `fulfillment` não-nulo, `provider_status_detail` não-nulo, e quantos dos **7 cancelados** têm `cancellation_detail`.

**Sele o contrato conforme o que você mediu, não conforme o que esperava.** Se `cancellation_detail` continuar vazio em algum cancelado, investigue antes de selar — pode ser cancelamento antigo que o ML não expõe mais, e aí o selo correto é `CONDICIONAL`, não `SEMPRE`.

- [ ] **Step 2: Selar**

| campo | selo | forma |
|---|---|---|
| `currency` | `SEMPRE` (o ML manda em todo pedido) | `required`, sem `nullable` |
| `fulfillment` | `CONDICIONAL(envio)` — pedido sem envio não tem modal | `nullable: true`, condição **na descrição** |
| `provider_status_detail` | `CONDICIONAL(provider informa)` | `nullable: true` |
| `cancellation_detail` | `CONDICIONAL(cancelado)` | `nullable: true` |

**Nunca `required` + `nullable` no mesmo campo.** É a forma que não recusa payload nenhum, e foi por ela que 34 preços nulos passaram por todos os gates no P1.

- [ ] **Step 3: Espelhar no SDK — mesmo commit**

`packages/sdk-runtime/src/index.ts` na mesma edição. Sem isso o gate de contrato reprova e o front mente sobre o tipo.

- [ ] **Step 4: Tela**

- Motivo do cancelamento visível nos cancelados. Traduza o `code` para rótulo em português **na camada de apresentação** — o domínio guarda o código estável.
- Modal logístico e rastreio no drawer.
- Ausente continua `—`. **Não invente rótulo para o que é nulo.**

```bash
cd apps/web && npx --no-install vitest run src/pages/PedidosPage.test.tsx
```

**Sobre o `tsc`:** o front já está vermelho com 13 erros herdados de 16–19/jul, blame confirmado, nenhum deste trabalho. Rode `npx --no-install tsc --noEmit -p tsconfig.json`, compare com essa linha de base, garanta que **você não somou nenhum**. Não persiga os herdados.

- [ ] **Step 5: Commit**

```bash
git add contracts/api/marketplace-central.openapi.yaml packages/sdk-runtime/src/index.ts apps/web/src
git commit -m "contract(orders): currency, fulfillment e motivo de cancelamento no read model"
```

---

## Verificação final — dirigida no navegador

- [ ] `currency` e `fulfillment` saíram de 0/38.
- [ ] `logistic_type` e `tracking_method` saíram de 0 no banco.
- [ ] Os 7 cancelados mostram motivo; pedido não cancelado tem `NULL` — **não** `""`. Prove com `IS NULL` no banco, não com igualdade a string vazia.
- [ ] **Margem continua em 33/38.** Se mudou, você tocou em algo fora do escopo — investigue antes de fechar.
- [ ] **`imposto` continua 4% e `difal` continua 0.** Estão errados de propósito; consertá-los é o P2.b. Se mudaram, reverta.
- [ ] `/pedidos` no navegador: console limpo, **0 chamadas ao provider** a partir da tela.
- [ ] Mandar `CLOSED` com essas medições.
