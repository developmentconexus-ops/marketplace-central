# P1 — Espinha de preço e estoque de /anuncios — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fazer `/anuncios` mostrar preço, moeda, tipo de anúncio e quantidade publicada corretos nos 34 anúncios reais, e instalar os mecanismos que impedem a mesma classe de defeito de voltar.

**Architecture:** O defeito nasce no primeiro salto — o DTO do adapter não declara chaves que o provider manda, e `encoding/json` as descarta em silêncio. Este plano ataca em ordem de causa: primeiro constrói o detector (reconciliação raw × DTO) e deixa que **ele nomeie** as quatro chaves perdidas; só então corrige o DTO, o mapper, a persistência do raw, e a regra de escrita que vinha apagando valor bom. Fecha selando o contrato para que um nulo nesses campos volte a reprovar.

**Tech Stack:** Go 1.x (`apps/server_core`), pgx + Postgres, OpenAPI 3 (`contracts/api/marketplace-central.openapi.yaml`) + `packages/sdk-runtime` escrito à mão, React/TS (`apps/web`, Vite).

---

## Contexto que o executor precisa antes de começar

**O sintoma.** Em `/anuncios`, 34 de 34 anúncios têm `price_amount`, `price_currency`,
`listing_type_code` e `published_quantity` nulos no banco. A tela mostra `—`.

**A causa, medida.** `mlMultigetItemBody`
(`internal/modules/connectors/adapters/mercado_livre/items_multiget_reader.go:136-161`) não
declara as chaves `price`, `currency_id`, `listing_type_id` e `initial_quantity`. O
Mercado Livre manda as quatro. `encoding/json` ignora chave não declarada por padrão — sem
erro, sem log. O que não entra no DTO não chega ao mapper, não chega ao banco.

**O agravante.** `UpsertPulledRows`
(`internal/modules/listings/adapters/postgres/repository.go:435-446`) faz
`SET x = EXCLUDED.x` cru em toda coluna. Como o mapper não preenche esses quatro campos,
**cada re-sync grava `NULL` por cima de valor que já esteve correto**. Os preços não estão
faltando: foram apagados.

**Comandos.** Tudo em `apps/server_core`. `GOCACHE` precisa ser caminho absoluto e ir inline
no comando — a forma `export X && ...` é recusada pelo classificador de permissão:

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/... -count=1
```

**Regras que valem para todas as tarefas:**
- Commit a cada tarefa. Nunca `git push` — exige autorização explícita do operador.
- Nenhum arquivo de teste pode conter PII real: nome, documento, endereço ou CEP de comprador
  **ou do vendedor** (o payload de item do ML traz o endereço do vendedor). Use os valores
  redigidos deste plano, literalmente.
- Não altere `contracts/api/*.yaml` sem alterar `packages/sdk-runtime` no **mesmo commit**.

---

### Task 1: Detector de chave não declarada — o mecanismo

Um helper compartilhado que compara as chaves de um payload cru com as que um DTO declara.
Vai em `platform` (não no adapter do ML) porque o ADR-C6 o exige para todo adapter — Sankhya
e xlsx reusam sem reimplementar.

**Files:**
- Create: `apps/server_core/internal/platform/rawkeys/rawkeys.go`
- Test: `apps/server_core/internal/platform/rawkeys/rawkeys_test.go`

- [ ] **Step 1: Write the failing test**

```go
package rawkeys

import (
	"encoding/json"
	"testing"
)

type sampleDTO struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Nested struct {
		Mode string `json:"mode"`
	} `json:"shipping"`
	NoTag    string
	Excluded string `json:"-"`
}

func TestUndeclaredNamesOnlyKeysTheDTOOmits(t *testing.T) {
	raw := json.RawMessage(`{"id":"MLB1","title":"x","shipping":{"mode":"me2"},"price":10,"currency_id":"BRL"}`)

	got, err := Undeclared(raw, sampleDTO{}, nil)
	if err != nil {
		t.Fatalf("Undeclared() error = %v", err)
	}

	want := []string{"currency_id", "price"}
	if len(got) != len(want) {
		t.Fatalf("Undeclared() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Undeclared() = %v, want %v", got, want)
		}
	}
}

func TestUndeclaredRespectsIgnoreList(t *testing.T) {
	raw := json.RawMessage(`{"id":"MLB1","title":"x","price":10}`)

	got, err := Undeclared(raw, sampleDTO{}, []string{"price"})
	if err != nil {
		t.Fatalf("Undeclared() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("Undeclared() = %v, want empty", got)
	}
}

func TestUndeclaredRejectsNonObjectPayload(t *testing.T) {
	if _, err := Undeclared(json.RawMessage(`[1,2]`), sampleDTO{}, nil); err == nil {
		t.Fatal("Undeclared() error = nil, want error for non-object payload")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:
```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/platform/rawkeys/ -count=1 -v
```
Expected: FAIL — o pacote nem compila (`undefined: Undeclared`).

- [ ] **Step 3: Write minimal implementation**

```go
// Package rawkeys detects provider payload keys that no DTO declares.
//
// It exists because encoding/json silently DISCARDS undeclared keys: an
// adapter DTO that forgets a field the provider sends produces no error, no
// log and no symptom — the value simply never reaches the domain. That exact
// failure emptied price/currency/listing_type/published_quantity on every
// listing (ADR-C6).
//
// It is deliberately in platform, not in one adapter: ADR-C6 applies to every
// adapter (Mercado Livre, Sankhya, xlsx), and a per-adapter copy would rebuild
// the asymmetry that caused the defect.
package rawkeys

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Undeclared returns the TOP-LEVEL keys present in raw that dto does not
// declare via a json tag, sorted, minus anything in ignore.
//
// Top-level only: nested shapes belong to their own DTO and are checked by
// passing that nested value in its own call. Reporting nested keys here would
// mix two different DTOs' responsibilities in one list.
func Undeclared(raw json.RawMessage, dto any, ignore []string) ([]string, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil, fmt.Errorf("rawkeys: payload is not a JSON object: %w", err)
	}

	declared := declaredKeys(reflect.TypeOf(dto))
	skip := make(map[string]struct{}, len(ignore))
	for _, key := range ignore {
		skip[key] = struct{}{}
	}

	var missing []string
	for key := range fields {
		if _, ok := declared[key]; ok {
			continue
		}
		if _, ok := skip[key]; ok {
			continue
		}
		missing = append(missing, key)
	}
	sort.Strings(missing)
	return missing, nil
}

func declaredKeys(t reflect.Type) map[string]struct{} {
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	keys := map[string]struct{}{}
	if t == nil || t.Kind() != reflect.Struct {
		return keys
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]
		if name == "" {
			// No json tag: encoding/json matches the Go field name
			// case-insensitively. Declaring the field name keeps the
			// detector aligned with the decoder's real behavior.
			name = field.Name
		}
		keys[name] = struct{}{}
	}
	return keys
}
```

- [ ] **Step 4: Run test to verify it passes**

Run:
```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/platform/rawkeys/ -count=1 -v
```
Expected: PASS — `TestUndeclaredNamesOnlyKeysTheDTOOmits`,
`TestUndeclaredRespectsIgnoreList`, `TestUndeclaredRejectsNonObjectPayload`.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/platform/rawkeys/
git commit -m "feat(platform): detector de chave de payload nao declarada pelo DTO"
```

---

### Task 2: Apontar o detector para o DTO real — o teste que nomeia o defeito

Este é o teste diagnóstico. Ele roda contra o DTO **como está hoje** e deve listar exatamente
as quatro chaves perdidas. É a prova de que o instrumento enxerga o defeito antes de
consertá-lo — sem isso, o verde da Task 3 não distingue "consertado" de "instrumento cego".

**Files:**
- Create: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/testdata/item_multiget_body.json`
- Create: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/items_multiget_rawkeys_test.go`

- [ ] **Step 1: Write the fixture**

Payload de item real do ML, com **todo valor de PII substituído por literal redigido**. As
chaves ficam; os valores sensíveis, não.

```json
{
  "id": "MLB4834219830",
  "title": "Produto de exemplo",
  "seller_id": 111111111,
  "category_id": "MLB270310",
  "price": 729.9,
  "base_price": 729.9,
  "original_price": null,
  "currency_id": "BRL",
  "listing_type_id": "gold_special",
  "available_quantity": 5,
  "initial_quantity": 12,
  "sold_quantity": 7,
  "condition": "new",
  "permalink": "https://produto.mercadolivre.com.br/MLB-4834219830",
  "thumbnail": "https://http2.mlstatic.com/D_NQ_NP_1.jpg",
  "status": "active",
  "sub_status": [],
  "tags": ["good_quality_thumbnail"],
  "date_created": "2025-11-04T13:22:31.000-04:00",
  "last_updated": "2026-07-30T09:11:02.000-04:00",
  "catalog_product_id": null,
  "catalog_listing": false,
  "user_product_id": "MLBU111111111",
  "inventory_id": null,
  "domain_id": "MLB-CELLPHONES",
  "seller_sku": "SKU-EXEMPLO-1",
  "seller_custom_field": null,
  "attributes": [
    {"id": "SELLER_SKU", "value_id": null, "value_name": "SKU-EXEMPLO-1"},
    {"id": "GTIN", "value_id": null, "value_name": "7891234567895"}
  ],
  "sale_terms": [],
  "pictures": [{"id": "111111-MLB", "url": "https://http2.mlstatic.com/D_111111-MLB.jpg"}],
  "shipping": {
    "mode": "me2",
    "logistic_type": "drop_off",
    "free_shipping": true,
    "local_pick_up": false,
    "store_pick_up": false,
    "tags": ["self_service_in"],
    "methods": [],
    "dimensions": null
  },
  "seller_address": {
    "city": {"name": "REDACTED"},
    "state": {"id": "BR-MG", "name": "REDACTED"},
    "country": {"id": "BR"},
    "address_line": "REDACTED",
    "zip_code": "REDACTED"
  },
  "variations": [],
  "deal_ids": [],
  "health": null,
  "warranty": null,
  "accepts_mercadopago": true,
  "automatic_relist": false
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

// multigetBodyIgnoredKeys are keys the multiget body DTO deliberately does NOT
// type. Every entry is a decision, not an oversight: the value is either
// already carried elsewhere, or out of scope for this seam.
//
// ADR-C6: a key may only be added here with a reason. An unexplained entry
// turns the detector back into the silence it exists to break.
var multigetBodyIgnoredKeys = []string{
	"accepts_mercadopago", // não usado em nenhuma tela
	"automatic_relist",    // não usado em nenhuma tela
	"base_price",          // preço de lista antes de promoção; fora do escopo do P1
	"catalog_listing",     // competição de catálogo = PLANEJADO, bloqueada por opt-in
	"deal_ids",            // promoções; fora do escopo do P1
	"domain_id",           // domínio da categoria; fora do escopo do P1
	"health",              // descontinuado pelo ML; a fonte é /item/{id}/performance
	"inventory_id",        // estoque Full; a conta não usa
	"last_updated",        // sincronização incremental é do M-06
	"original_price",      // preço de lista antes de promoção; fora do escopo do P1
	"pictures",            // thumbnail já é tipado
	"sale_terms",          // garantia/termos; fora do escopo do P1
	"seller_address",      // PII do vendedor: nunca sai do adapter
	"seller_id",           // já conhecido pela conta da instalação
	"sub_status",          // motivo de bloqueio = PLANEJADO (P4)
	"user_product_id",     // estoque Full; a conta não usa
	"warranty",            // fora do escopo do P1
}

func TestMultigetItemBodyDeclaresEveryConsumedKey(t *testing.T) {
	raw, err := os.ReadFile("testdata/item_multiget_body.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	missing, err := rawkeys.Undeclared(json.RawMessage(raw), mlMultigetItemBody{}, multigetBodyIgnoredKeys)
	if err != nil {
		t.Fatalf("Undeclared() error = %v", err)
	}
	if len(missing) != 0 {
		t.Fatalf("mlMultigetItemBody nao declara chaves que o ML manda: %v", missing)
	}
}
```

- [ ] **Step 3: Run test to verify it fails — e confira O QUE ele diz**

Run:
```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/connectors/adapters/mercado_livre/ -count=1 -run TestMultigetItemBodyDeclaresEveryConsumedKey -v
```
Expected: FAIL com exatamente esta lista:
```
mlMultigetItemBody nao declara chaves que o ML manda: [currency_id initial_quantity listing_type_id price]
```

**Se a lista vier diferente das quatro, pare e reporte.** Chave a mais significa que a lista
de ignorados está incompleta; chave a menos significa que o instrumento está cego, e nesse
caso o verde da próxima tarefa não provaria nada.

- [ ] **Step 4: Commit o teste vermelho**

```bash
git add apps/server_core/internal/modules/connectors/adapters/mercado_livre/testdata/ apps/server_core/internal/modules/connectors/adapters/mercado_livre/items_multiget_rawkeys_test.go
git commit -m "test(mercado_livre): teste vermelho nomeia as 4 chaves que o DTO do multiget perde"
```

---

### Task 3: O DTO declara as quatro chaves

**Files:**
- Modify: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/items_multiget_reader.go:56-85` (`ItemMultigetDTO`)
- Modify: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/items_multiget_reader.go:136-161` (`mlMultigetItemBody`)

- [ ] **Step 1: Declare as chaves no shape de fio**

Em `mlMultigetItemBody`, logo depois de `Status`, insira:

```go
	// Price/CurrencyID/ListingTypeID/InitialQuantity: o ML manda os quatro em
	// TODO item, e este struct não os declarava — encoding/json os descartava
	// em silêncio, deixando price_amount/price_currency/listing_type_code/
	// published_quantity NULL em 34 de 34 anúncios. Price é json.Number (e não
	// float64) pela mesma convenção de dinheiro do resto do pacote: evita
	// deriva de precisão binária em valor monetário.
	Price           *json.Number `json:"price"`
	CurrencyID      string       `json:"currency_id"`
	ListingTypeID   string       `json:"listing_type_id"`
	InitialQuantity *int         `json:"initial_quantity"`
```

- [ ] **Step 2: Exponha no DTO**

Em `ItemMultigetDTO`, depois de `Status`:

```go
	// Price é a forma decimal em string (mesma convenção de
	// ItemMultigetVariationDTO.Price) — nunca float.
	Price           *string
	CurrencyID      string
	ListingTypeID   string
	// InitialQuantity é a quantidade PUBLICADA. Não é
	// AvailableQuantity + SoldQuantity: reposição de estoque incrementa a
	// disponível sem tocar a inicial. São dois números distintos.
	InitialQuantity *int
```

- [ ] **Step 3: Preencha na conversão de fio para DTO**

Encontre onde `mlMultigetItemBody` vira `ItemMultigetDTO` (mesmo arquivo, na função que trata
cada elemento do multiget — procure por `Thumbnail:` para achar o literal do struct) e
acrescente ao literal:

```go
		Price:           multigetPriceString(body.Price),
		CurrencyID:      strings.TrimSpace(body.CurrencyID),
		ListingTypeID:   strings.TrimSpace(body.ListingTypeID),
		InitialQuantity: body.InitialQuantity,
```

E adicione o helper no fim do arquivo:

```go
// multigetPriceString converte o número de fio para a forma decimal em string.
// json.Number já preserva os dígitos exatos recebidos; converter para float
// aqui reintroduziria a deriva binária que a convenção de dinheiro do pacote
// existe para evitar.
func multigetPriceString(n *json.Number) *string {
	if n == nil {
		return nil
	}
	s := strings.TrimSpace(n.String())
	if s == "" {
		return nil
	}
	return &s
}
```

- [ ] **Step 4: Run test to verify it passes**

Run:
```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/connectors/adapters/mercado_livre/ -count=1 -run TestMultigetItemBodyDeclaresEveryConsumedKey -v
```
Expected: PASS.

Rode também o pacote inteiro para garantir que nada quebrou:
```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/connectors/... -count=1
```
Expected: ok, sem FAIL.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/connectors/adapters/mercado_livre/items_multiget_reader.go
git commit -m "fix(mercado_livre): DTO do multiget declara price, currency_id, listing_type_id e initial_quantity"
```

---

### Task 4: O mapper leva os quatro campos ao domínio

Declarar no DTO não basta: o mapper de produção precisa atribuí-los. Esta é a segunda metade
do salto que estava quebrado.

**Files:**
- Modify: `apps/server_core/internal/modules/listings/adapters/connectors/multiget_mapper.go:99-123`
- Test: `apps/server_core/internal/modules/listings/adapters/connectors/multiget_mapper_test.go`

- [ ] **Step 1: Write the failing test**

Acrescente ao arquivo de teste existente:

```go
func TestMapMultigetItemToListingCarriesPriceAndQuantities(t *testing.T) {
	price := "729.90"
	initial := 12
	fetchedAt := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)

	row, err := MapMultigetItemToListing("t1", "inst-1", "mercado_livre", mercadolivre.ItemMultigetDTO{
		ProviderItemID:  "MLB4834219830",
		Title:           "Produto de exemplo",
		Status:          "active",
		Price:           &price,
		CurrencyID:      "BRL",
		ListingTypeID:   "gold_special",
		InitialQuantity: &initial,
	}, fetchedAt)
	if err != nil {
		t.Fatalf("MapMultigetItemToListing() error = %v", err)
	}

	if row.PriceAmount == nil || string(*row.PriceAmount) != "729.90" {
		t.Fatalf("PriceAmount = %v, want 729.90", row.PriceAmount)
	}
	if row.PriceCurrency == nil || string(*row.PriceCurrency) != "BRL" {
		t.Fatalf("PriceCurrency = %v, want BRL", row.PriceCurrency)
	}
	if row.ListingTypeCode == nil || string(*row.ListingTypeCode) != "gold_special" {
		t.Fatalf("ListingTypeCode = %v, want gold_special", row.ListingTypeCode)
	}
	if row.PublishedQuantity == nil || *row.PublishedQuantity != 12 {
		t.Fatalf("PublishedQuantity = %v, want 12", row.PublishedQuantity)
	}
}

func TestMapMultigetItemToListingKeepsAbsentPriceNil(t *testing.T) {
	fetchedAt := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)

	row, err := MapMultigetItemToListing("t1", "inst-1", "mercado_livre", mercadolivre.ItemMultigetDTO{
		ProviderItemID: "MLB4834219830",
		Title:          "Produto de exemplo",
		Status:         "active",
	}, fetchedAt)
	if err != nil {
		t.Fatalf("MapMultigetItemToListing() error = %v", err)
	}

	// ADR-17: desconhecido é nil, nunca zero nem string vazia.
	if row.PriceAmount != nil {
		t.Fatalf("PriceAmount = %v, want nil", row.PriceAmount)
	}
	if row.PriceCurrency != nil {
		t.Fatalf("PriceCurrency = %v, want nil", row.PriceCurrency)
	}
	if row.ListingTypeCode != nil {
		t.Fatalf("ListingTypeCode = %v, want nil", row.ListingTypeCode)
	}
	if row.PublishedQuantity != nil {
		t.Fatalf("PublishedQuantity = %v, want nil", row.PublishedQuantity)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:
```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/listings/adapters/connectors/ -count=1 -run TestMapMultigetItemToListingCarries -v
```
Expected: FAIL — `PriceAmount = <nil>, want 729.90`.

- [ ] **Step 3: Write minimal implementation**

Em `MapMultigetItemToListing`, antes da chamada a `listingsdomain.NewListing`, acrescente:

```go
	// Os quatro campos que o DTO passou a declarar. ADR-17: ausência do
	// provider vira nil, nunca zero — quem não tem preço não tem preço zero.
	var priceAmount *listingsdomain.PriceAmount
	if item.Price != nil {
		amount := listingsdomain.PriceAmount(strings.TrimSpace(*item.Price))
		if amount != "" {
			priceAmount = &amount
		}
	}
	var priceCurrency *listingsdomain.PriceCurrency
	if currency := strings.TrimSpace(item.CurrencyID); currency != "" {
		value := listingsdomain.PriceCurrency(currency)
		priceCurrency = &value
	}
	var listingTypeCode *listingsdomain.ListingTypeCode
	if code := strings.TrimSpace(item.ListingTypeID); code != "" {
		value := listingsdomain.ListingTypeCode(code)
		listingTypeCode = &value
	}
```

E dentro do literal `listingsdomain.ListingInput{...}`, logo depois de `Title:`, acrescente:

```go
		ListingTypeCode:   listingTypeCode,
		PriceAmount:       priceAmount,
		PriceCurrency:     priceCurrency,
		// PublishedQuantity é a quantidade PUBLICADA (initial_quantity), que o
		// ML mantém distinta de available_quantity: reposição sobe a
		// disponível sem tocar a inicial.
		PublishedQuantity: item.InitialQuantity,
```

- [ ] **Step 4: Run test to verify it passes**

Run:
```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/listings/... -count=1
```
Expected: ok, sem FAIL — inclusive
`TestMapMultigetItemToListingKeepsAbsentPriceNil`.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/listings/adapters/connectors/
git commit -m "fix(listings): mapper de multiget leva preco, moeda, tipo e quantidade publicada ao dominio"
```

---

### Task 5: Persistir o raw — pré-requisito da reconciliação contínua

As colunas `listings.raw` e `listings.raw_truncated` **já existem** e nunca são escritas.
O `ItemMultigetDTO` já captura `Raw` e `RawTruncated`. Falta só o caminho até o banco.
Sem isso, o detector da Task 1 só roda contra fixture — e fixture é cego a campo novo que o
provider lançar depois.

**Files:**
- Modify: `apps/server_core/internal/modules/listings/domain/listing.go:139-172` (`ListingInput`) e o struct `Listing`
- Modify: `apps/server_core/internal/modules/listings/adapters/connectors/multiget_mapper.go`
- Modify: `apps/server_core/internal/modules/listings/adapters/postgres/repository.go:424-453`
- Test: `apps/server_core/internal/modules/listings/adapters/connectors/multiget_mapper_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestMapMultigetItemToListingCarriesRaw(t *testing.T) {
	fetchedAt := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	raw := json.RawMessage(`{"id":"MLB4834219830","price":729.9}`)

	row, err := MapMultigetItemToListing("t1", "inst-1", "mercado_livre", mercadolivre.ItemMultigetDTO{
		ProviderItemID: "MLB4834219830",
		Title:          "Produto de exemplo",
		Status:         "active",
		Raw:            raw,
		RawTruncated:   true,
	}, fetchedAt)
	if err != nil {
		t.Fatalf("MapMultigetItemToListing() error = %v", err)
	}

	if string(row.Raw) != string(raw) {
		t.Fatalf("Raw = %s, want %s", row.Raw, raw)
	}
	if !row.RawTruncated {
		t.Fatal("RawTruncated = false, want true")
	}
}
```

Acrescente `"encoding/json"` aos imports do arquivo de teste se ainda não estiver lá.

- [ ] **Step 2: Run test to verify it fails**

Run:
```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/listings/adapters/connectors/ -count=1 -run TestMapMultigetItemToListingCarriesRaw -v
```
Expected: FAIL — não compila, `row.Raw undefined`.

- [ ] **Step 3: Write minimal implementation**

**3a.** Em `listing.go`, acrescente ao struct `Listing` e ao struct `ListingInput` (nos dois,
no fim, antes do fecha-chaves):

```go
	// Raw é o payload cru do provider, já com PII redigida na CAPTURA
	// (ADR-C6). Existe para que a reconciliação de chaves detecte campo que o
	// provider manda e nenhum DTO declara. RawTruncated marca que o payload
	// excedeu o teto de captura do adapter e NÃO é a resposta inteira — sem
	// esse marcador, uma chave ausente por truncamento seria indistinguível
	// de uma chave que o provider não mandou.
	Raw          json.RawMessage
	RawTruncated bool
```

Acrescente `"encoding/json"` aos imports de `listing.go`.

Em `NewListing`, antes do `return`, copie os dois campos do input para o valor construído
(siga o padrão que a função já usa para os demais campos).

**3b.** Em `MapMultigetItemToListing`, dentro do literal `ListingInput{...}`:

```go
		Raw:          item.Raw,
		RawTruncated: item.RawTruncated,
```

**3c.** Em `UpsertPulledRows`, acrescente as duas colunas. Na lista do `INSERT`, depois de
`absent_since`, e nos parâmetros. O `SELECT` passa a ter mais dois placeholders — renumere os
dois últimos (`$32`, `$33` do `WHERE`) para `$34` e `$35`:

```sql
			INSERT INTO listings (tenant_id, installation_id, provider, provider_listing_id, variation_id, title,
				listing_type_code, status, price_amount, price_currency, published_quantity, sync_state,
				sync_error, quality_score, sales_30d, fetched_at,
				sold_quantity, category_id, condition, permalink, thumbnail, date_created_ml, tags,
				catalog_product_id, shipping_mode, free_shipping, logistic_type, available_quantity,
				last_seen_at, absent_since, raw, raw_truncated, created_at, updated_at)
			SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,
				$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,
				$29,NULL,$30,$31,$32,$33
			WHERE $1 = $34 AND $2 = $35
```

E nos argumentos, entre `seenAt` e `row.CreatedAt.UTC()`:

```go
			seenAt, rawOrNil(row.Raw), row.RawTruncated, row.CreatedAt.UTC(), seenAt,
```

Acrescente `"encoding/json"` aos imports de `repository.go` — o helper abaixo precisa dele.
Helper no fim de `repository.go`:

```go
// rawOrNil evita gravar o literal JSON "null" em vez de SQL NULL quando o
// adapter não capturou payload. Os dois se parecem numa consulta e significam
// coisas diferentes: "o provider mandou null" e "não temos payload".
func rawOrNil(raw json.RawMessage) any {
	if len(raw) == 0 {
		return nil
	}
	return []byte(raw)
}
```

O `DO UPDATE SET` **não** ganha `raw` ainda — a Task 6 reescreve esse bloco inteiro.

- [ ] **Step 4: Run test to verify it passes**

Run:
```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/listings/... -count=1 && GOCACHE=$(pwd)/.gocache go build ./...
```
Expected: ok, sem FAIL; build sem erro.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/listings/
git commit -m "feat(listings): persistir payload cru e marcador de truncamento (pre-requisito ADR-C6)"
```

---

### Task 6: Escritor dono das colunas — a regra que para de apagar valor bom

`SET x = EXCLUDED.x` cru foi o que apagou os preços. Trocar por `COALESCE` em tudo seria
máximo local: nunca apagaria, nem quando deve — anúncio que perde o SKU no provider ficaria
com o SKU velho para sempre, e dado zumbi é pior de auditar que dado ausente.

A regra é a do ADR-C2: **o escritor declara as colunas que possui, e o `UPDATE` toca só
essas.**

**Files:**
- Modify: `apps/server_core/internal/modules/listings/adapters/postgres/repository.go:435-446`
- Test: `apps/server_core/internal/modules/listings/adapters/postgres/repository_integration_test.go`

- [ ] **Step 1: Write the failing test**

```go
//go:build integration

func TestUpsertPulledRowsDoesNotEraseColumnsOfAnotherWriter(t *testing.T) {
	testpostgres.SkipWithoutTarget(t)
	ctx := context.Background()
	repo, installationID := newIntegrationRepository(t)

	price := domain.PriceAmount("729.90")
	currency := domain.PriceCurrency("BRL")
	at := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)

	full := newIntegrationListing(t, installationID, "MLB1")
	full.PriceAmount = &price
	full.PriceCurrency = &currency
	if err := repo.UpsertPulledRows(ctx, installationID, []domain.Listing{full}, at); err != nil {
		t.Fatalf("primeira escrita: %v", err)
	}

	// Um produtor que legitimamente não tem preço reescreve a MESMA linha.
	// Antes do ADR-C2 isso gravava NULL por cima do preço correto.
	blind := newIntegrationListing(t, installationID, "MLB1")
	blind.PriceAmount = nil
	blind.PriceCurrency = nil
	if err := repo.UpsertPulledRows(ctx, installationID, []domain.Listing{blind}, at.Add(time.Minute)); err != nil {
		t.Fatalf("segunda escrita: %v", err)
	}

	var gotAmount, gotCurrency *string
	if err := repo.pool.QueryRow(ctx,
		`SELECT price_amount::text, price_currency FROM listings WHERE provider_listing_id = $1`,
		"MLB1").Scan(&gotAmount, &gotCurrency); err != nil {
		t.Fatalf("leitura: %v", err)
	}
	if gotAmount == nil || *gotAmount != "729.90" {
		t.Fatalf("price_amount = %v, want 729.90 preservado", gotAmount)
	}
	if gotCurrency == nil || *gotCurrency != "BRL" {
		t.Fatalf("price_currency = %v, want BRL preservado", gotCurrency)
	}
}
```

Se `newIntegrationRepository` ou `newIntegrationListing` não existirem com esses nomes no
arquivo, use os helpers que o arquivo já tiver — leia o topo dele antes de escrever.

- [ ] **Step 2: Run test to verify it fails**

A lane hermética exige nome de banco casando `^mpc_test_[0-9a-f]{32}$` **e** host em IP de
loopback (`127.0.0.1`, nunca `localhost` — `net.ParseIP` falha em nome). Sem
`MPC_TEST_DATABASE_URL` o teste faz SKIP silencioso, que é indistinguível de verde.

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache MPC_TEST_DATABASE_URL='postgres://marketplace:<senha>@127.0.0.1:5435/mpc_test_<32hex>' go run ./cmd/testdb migrate
cd apps/server_core && GOCACHE=$(pwd)/.gocache MPC_TEST_DATABASE_URL='postgres://marketplace:<senha>@127.0.0.1:5435/mpc_test_<32hex>' go test -tags integration ./internal/modules/listings/adapters/postgres/ -count=1 -run TestUpsertPulledRowsDoesNotEraseColumnsOfAnotherWriter -v
```
Expected: FAIL com `price_amount = <nil>, want 729.90 preservado`. **Confirme que a saída traz
`--- FAIL`, não `--- SKIP`** — SKIP aqui é instrumento morto, não aprovação.

- [ ] **Step 3: Write minimal implementation**

Substitua o bloco `ON CONFLICT ... DO UPDATE SET` inteiro por:

```sql
			ON CONFLICT (tenant_id, installation_id, provider_listing_id, variation_id) DO UPDATE SET
				provider=EXCLUDED.provider, title=EXCLUDED.title, listing_type_code=EXCLUDED.listing_type_code,
				status=EXCLUDED.status, price_amount=EXCLUDED.price_amount, price_currency=EXCLUDED.price_currency,
				published_quantity=EXCLUDED.published_quantity, sync_state=EXCLUDED.sync_state,
				sync_error=EXCLUDED.sync_error, fetched_at=EXCLUDED.fetched_at,
				sold_quantity=EXCLUDED.sold_quantity, category_id=EXCLUDED.category_id, condition=EXCLUDED.condition,
				permalink=EXCLUDED.permalink, thumbnail=EXCLUDED.thumbnail, date_created_ml=EXCLUDED.date_created_ml,
				tags=EXCLUDED.tags, catalog_product_id=EXCLUDED.catalog_product_id, shipping_mode=EXCLUDED.shipping_mode,
				free_shipping=EXCLUDED.free_shipping, logistic_type=EXCLUDED.logistic_type,
				available_quantity=EXCLUDED.available_quantity,
				raw=EXCLUDED.raw, raw_truncated=EXCLUDED.raw_truncated,
				last_seen_at=EXCLUDED.last_seen_at, absent_since=NULL, updated_at=EXCLUDED.updated_at
```

Note o que **saiu** da lista: `quality_score` e `sales_30d`. Documente logo acima da query:

```go
	// ADR-C2 — conjunto de colunas deste escritor.
	//
	// Este escritor é o backfill/sweep do catálogo do provider. Ele possui
	// tudo que vem do payload do item, e SÓ isso. As colunas de fora ficam
	// intocadas mesmo quando a linha é reescrita:
	//
	//   quality_score — dono é o leitor de saúde do anúncio (P4). Este
	//                   escritor nunca teve o valor; escrevê-lo aqui apagaria
	//                   o do dono a cada sweep.
	//   sales_30d     — campo DERIVADO de pedidos por janela, não lido do
	//                   provider (a API não tem endpoint de vendas por
	//                   janela). Dono é o agregador de pedidos.
	//
	// Não use COALESCE para "proteger" coluna alheia: COALESCE também impede
	// o DONO de apagar quando o provider REMOVE o valor, trocando perda de
	// dado por dado zumbi.
```

- [ ] **Step 4: Run test to verify it passes**

Run (com o mesmo `MPC_TEST_DATABASE_URL` da Step 2):
```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache MPC_TEST_DATABASE_URL='<mesmo dsn>' go test -tags integration ./internal/modules/listings/adapters/postgres/ -count=1 -v
```
Expected: PASS, com `--- PASS: TestUpsertPulledRowsDoesNotEraseColumnsOfAnotherWriter` nomeado
na saída.

E a suíte comum:
```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/... -count=1
```
Expected: ok, sem FAIL.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/listings/adapters/postgres/
git commit -m "fix(listings): escritor declara suas colunas; sweep para de apagar quality_score e sales_30d"
```

---

### Task 7: Matar o mapper morto

`MapListingSnapshotToCanonicalRows` (`mapper.go`) preenche corretamente os quatro campos —
**e não tem nenhum chamador de produção**, só `mapper_test.go` e um comentário em
`multiget_mapper.go:30`. Enquanto isso o mapper vivo não os preenchia. Dois mappers para o
mesmo destino, um certo e morto, outro vivo e incompleto: é essa duplicação que permite ligar
o errado sem ninguém notar.

**Files:**
- Delete: `apps/server_core/internal/modules/listings/adapters/connectors/mapper.go`
- Delete: `apps/server_core/internal/modules/listings/adapters/connectors/mapper_test.go`
- Modify: `apps/server_core/internal/modules/listings/adapters/connectors/multiget_mapper.go:30-45` (comentário que o cita)

- [ ] **Step 1: Confirme a ausência de chamador antes de apagar**

```bash
cd apps/server_core && grep -rn "MapListingSnapshotToCanonicalRows" --include=*.go internal/ cmd/
```
Expected: apenas ocorrências em `mapper.go`, `mapper_test.go` e o comentário em
`multiget_mapper.go`. **Se aparecer qualquer outro arquivo, PARE** — há chamador de produção
e esta tarefa muda de natureza; reporte em vez de apagar.

- [ ] **Step 2: Verifique se `mapper.go` tem helper usado pelo mapper vivo**

```bash
cd apps/server_core && grep -n "^func " internal/modules/listings/adapters/connectors/mapper.go
```

Se algum helper listado for usado por `multiget_mapper.go` (procure cada nome com `grep -rn`),
**mova esse helper para `multiget_mapper.go` antes de apagar o arquivo**. Só a função
`MapListingSnapshotToCanonicalRows` e o que for exclusivo dela devem sumir.

- [ ] **Step 3: Apague e corrija o comentário**

```bash
cd apps/server_core && git rm internal/modules/listings/adapters/connectors/mapper.go internal/modules/listings/adapters/connectors/mapper_test.go
```

Em `multiget_mapper.go`, o bloco de comentário `// Unlike MapListingSnapshotToCanonicalRows
(mapper.go, ...)` passa a citar arquivo que não existe. Substitua as duas referências por uma
frase que se sustenta sozinha:

```go
// MapMultigetItemsToListings é o Passo 2 do F-03 (backfill-cursor-ingest):
// converte um lote hidratado de mercadolivre.ItemMultigetDTO (Passo 1,
// GetItemsMultiget) em linhas canônicas domain.Listing.
//
// Duas regras deste mapper, ambas exigidas pelo IC-07:
//   - status VERBATIM, sem remapeamento para um conjunto fechado — o
//     vocabulário é do provider, e um remap silenciosamente descartaria
//     status legítimos que a tabela já aceita.
//   - UMA linha por item (VariationID = NoVariationID, ADR-019), com os dados
//     de variação como filhos em Variations, e não uma linha achatada por
//     variação.
//
// Erro em um item (item.Err != nil) é registrado e pulado, nunca aborta o
// lote — mesma garantia do IC-06.
```

- [ ] **Step 4: Run tests to verify nothing broke**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go build ./... && GOCACHE=$(pwd)/.gocache go vet ./internal/... && GOCACHE=$(pwd)/.gocache go test ./internal/... -count=1
```
Expected: build ok, vet limpo, testes ok sem FAIL.

- [ ] **Step 5: Commit**

```bash
git add -A apps/server_core/internal/modules/listings/adapters/connectors/
git commit -m "refactor(listings): remove mapper sem chamador de producao que duplicava o destino canonico"
```

---

### Task 8: Selar o contrato — nulo nesses campos volta a reprovar

Hoje `ListingReadModel` declara `price`, `published_quantity`, `quality_score` e `sales_30d`
como `required` **e** `nullable`. Obrigatório-que-pode-ser-nulo é forma que não recusa payload
nenhum: foi por ela que 34 nulos chegaram à tela sem nenhum gate acusar (ADR-C3).

Agora que o dado existe, o selo pode apertar. Faça esta tarefa **depois** da Task 6 e de um
sweep real — apertar antes deixaria a API respondendo 500 sobre dado que ainda não foi
preenchido.

**Files:**
- Modify: `contracts/api/marketplace-central.openapi.yaml` (schema `ListingReadModel`, ~L3588-3661)
- Modify: `packages/sdk-runtime/src/index.ts` (~L380-436)

- [ ] **Step 1: Rode um sweep e confirme que o dado chegou**

```bash
curl -s -X POST http://localhost:8080/listings/refresh -H 'Content-Type: application/json' -d '{"installation_id":"inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0"}'
```

```bash
docker exec marketplace-central-postgres-1 psql -U marketplace -d marketplace_central -tAc "select count(*) total, count(price_amount) preco, count(price_currency) moeda, count(listing_type_code) tipo, count(published_quantity) publicada from listings;"
```
Expected: `34|34|34|34|34`.

**Se qualquer coluna vier abaixo de 34, PARE e reporte** — apertar o contrato sobre dado
incompleto quebra a tela sem entregar nada.

- [ ] **Step 2: Aperte o schema**

Em `ListingReadModel`, nos campos `price`, `price_currency`, `listing_type_code` e
`published_quantity`, **remova `nullable: true`** e mantenha os quatro em `required`.

Nos campos `quality_score` e `sales_30d`, faça o oposto — **remova-os de `required` e do
bloco `properties`**, e registre por quê:

```yaml
    # quality_score e sales_30d saíram do contrato: selo PLANEJADO (ADR-C3).
    # quality_score tem fonte confirmada (GET /item/{id}/performance) e dono no
    # plano P4. sales_30d é DERIVADO de pedidos por janela — a API do provider
    # não tem endpoint de vendas por período. Campo não implementado não fica
    # nulável para sempre: sai do contrato até ter produtor.
```

- [ ] **Step 3: Espelhe no SDK — mesmo commit**

Em `packages/sdk-runtime/src/index.ts`, na interface `ListingReadModel`:

```ts
  price: string;
  price_currency: string;
  listing_type_code: string;
  published_quantity: number;
  // quality_score e sales_30d removidos: selo PLANEJADO (ADR-C3). Voltam com
  // o produtor — quality_score no P4, sales_30d quando houver agregador de
  // pedidos por janela.
```

Remova as declarações de `quality_score` e `sales_30d` da interface.

- [ ] **Step 4: Verifique tipo e build**

```bash
cd apps/web && npx --no-install tsc --noEmit -p tsconfig.json
```
Expected: **erro** em todo ponto do front que ainda lê `quality_score` ou `sales_30d`. Isso é
o resultado desejado: o selo transformou um branco silencioso em erro de compilação. Remova
esses usos da tela (a coluna sai; não vira `—`, conforme ADR-C1) e rode de novo até limpar.

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/... -count=1
```
Expected: ok. Se algum teste de contrato falhar por causa do schema, ajuste o teste ao novo
selo — não afrouxe o schema de volta.

- [ ] **Step 5: Commit**

```bash
git add contracts/api/marketplace-central.openapi.yaml packages/sdk-runtime/src/index.ts apps/web/src
git commit -m "contract(listings): selo governa a forma — preco nao-nulavel, quality_score e sales_30d saem ate ter produtor"
```

---

### Task 9: Tirar a string chumbada da tela

O cabeçalho de grupo de `/anuncios` renderiza o literal `"ERP est. —"`. Não é dado ausente
renderizado como traço: é texto escrito no componente, que mostraria o mesmo traço se o
estoque do ERP estivesse cheio.

**Files:**
- Modify: `apps/web/src/pages/AnunciosTable.tsx:223-266`

- [ ] **Step 1: Localize o literal**

```bash
cd apps/web && grep -rn "ERP est" src/
```
Expected: uma ocorrência em `src/pages/AnunciosTable.tsx`.

- [ ] **Step 2: Troque o literal pelo dado, com os três estados do ADR-C1**

`ListingGroup` não carrega campo de estoque do ERP. Enquanto o produtor não existe, o correto
não é traço — é a coluna **não existir** (ADR-C1: capability ausente esconde a coluna, só
desconhecido vira `—`).

Remova o trecho que renderiza `ERP est. —` do cabeçalho de grupo. Não coloque nada no lugar.

- [ ] **Step 3: Registre a dívida com dono, não com traço**

Acrescente logo acima do cabeçalho de grupo:

```tsx
{/*
  Estoque do ERP saiu do cabeçalho: era o literal "ERP est. —", texto escrito
  no componente que mostraria o mesmo traço com o estoque cheio. ADR-C1 — sem
  produtor, a coluna não existe; traço é reservado a desconhecido de verdade.
  Volta quando o espelho do ERP expuser o saldo por anúncio vinculado.
*/}
```

- [ ] **Step 4: Verifique tipo e tela**

```bash
cd apps/web && npx --no-install tsc --noEmit -p tsconfig.json
```
Expected: sem erro.

```bash
cd apps/web && npx --no-install vitest run src/pages/AnunciosTable.test.tsx
```
Expected: PASS. Se algum teste afirmava a presença de `ERP est.`, remova essa asserção — ela
certificava o literal, não o dado.

- [ ] **Step 5: Commit**

```bash
git add apps/web/src/pages/AnunciosTable.tsx apps/web/src/pages/AnunciosTable.test.tsx
git commit -m "fix(web): remove literal 'ERP est. —' do cabecalho de grupo de /anuncios"
```

---

## Verificação final — dirigida no navegador

O plano só está feito quando a tela mostra o dado. Testes verdes com a tela errada já
aconteceram nesta missão.

- [ ] **Sweep completo**

```bash
curl -s -X POST http://localhost:8080/listings/refresh -H 'Content-Type: application/json' -d '{"installation_id":"inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0"}'
```

- [ ] **Preenchimento no banco**

```bash
docker exec marketplace-central-postgres-1 psql -U marketplace -d marketplace_central -tAc "select count(*) total, count(price_amount) preco, count(listing_type_code) tipo, count(published_quantity) publicada, count(raw) raw from listings;"
```
Expected: `34|34|34|34|34`.

- [ ] **Não-regressão: o segundo sweep não apaga**

Rode o mesmo `curl` uma segunda vez e repita a consulta. Expected: os mesmos `34` em todas as
colunas. Este é o controle que prova o ADR-C2 no caminho real — antes dele, a segunda corrida
zerava as quatro colunas.

- [ ] **Drive na tela**

Abra `http://localhost:5174/anuncios`. Confirme:
- preço com valor real em toda linha, nenhum `—` na coluna de preço;
- tipo de anúncio preenchido;
- nenhuma ocorrência do texto `ERP est.` no DOM;
- console sem erro;
- `performance.getEntriesByType('resource')` filtrado por `/mercadolibre|mercadolivre/i`
  devolve **0** — a tela lê do nosso banco, nunca do provider.

---

## Fora do escopo deste plano

Cada item abaixo tem plano dono. Nenhum é dívida sem endereço:

| item | plano |
|---|---|
| comissão, frete do vendedor, margem, `cancel_detail`, nulabilidade de tipos de pedido | P2 |
| motor de DIFAL | P3 |
| `quality_score` via `/performance`, coluna PENDÊNCIA via `/infractions` | P4 |
| `logistic_type`, `tracking_method`, `/sla`, `/history`, `/invoice_data` | P5 |
| KPIs agregados no backend | P6 |
| competição de catálogo (`price_to_win`, `boosts[]`) | bloqueado por opt-in de catálogo no ML — ação do operador |
| reconciliação de raw rodando como job contínuo sobre `listings.raw` | P4 reusa `rawkeys`; o job periódico entra junto do primeiro consumidor |
