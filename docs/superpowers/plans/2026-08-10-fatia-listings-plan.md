# Fatia vertical — perna `listings` (contexto + adapter Mercado Livre)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development
> (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use
> checkbox (`- [ ]`) syntax for tracking.

**Goal:** Aterrar a segunda perna da fatia vertical do protocolo
(`docs/superpowers/specs/2026-08-06-protocolo-de-codigo-design.md` §15.3):
`contexts/listings` (anúncios observados no canal) alimentado por
`adapters/marketplace/mercadolivre` no molde §2, com caminho de produção real
(`cmd/listingsingest`) e payloads reais gravados.

**Architecture:** Réplica do molde ratificado pela Onda 0 (catalog): contexto com
`contracts`/`port` públicos e `internal/` privado, adapter de marketplace com DTO de fio
confinado em `internal/api`, composição na raiz, CLI de operador. A perna testa a metade do
molde que a Onda 0 NÃO testou: o formato `adapters/marketplace/<vendor>/` (§2), incluindo
auth por token e paginação por scroll do ML.

**Tech Stack:** Go (pgx/v5, net/http), Postgres (schema próprio + RLS FORCE), PowerShell
gate lanes.

**Branch:** `fatia-listings` (este plano está commitado nele). Aterrissagem: PR → check
`required` → merge pelo operador.

---

## Restrições globais (herdadas, não opcionais)

1. Comandos Go SEMPRE de `apps/server_core` com caches absolutos:
   `$env:GOCACHE = (Join-Path (Get-Location) '.gocache')`; idem `GOMODCACHE`. O gate faz
   isso sozinho; execução manual precisa fazer igual.
2. PowerShell do operador não aceita `&&` — separar com `;` em comandos colados.
3. Nunca `git push` de main; push de branch de trabalho tem autorização em pé
   (`docs/HARNESS-PROFILE.md` §9). Merge é do operador, por evento.
4. RED antes do código, contra árvore limpa do arquivo sob teste.
5. Desconhecido nunca vira zero/default — é `fact.Fact` com estado explícito (§4.1 do
   protocolo; o kernel já impõe).
6. Nenhum mock em seam de integração; a lane hermética usa Postgres real.
7. Escrita viva em provider: NENHUMA — este plano só faz GET no ML.

## Medição (mc-planning Fase 1 — tudo aberto em 2026-08-10, árvore `23989a29`)

1. **O que já existe que faz parte disto?**
   - Molde do contexto: `contexts/catalog` completo — contracts
     (`internal/contexts/catalog/contracts/observation.go:12-61`), port com cursor opaco
     (`port/feed.go:19-46`), fachada `New(pool)` (`module.go:32-38`), repositório com RLS
     por transação (`internal/postgres/repository.go:49`
     `set_config('app.tenant_id',$1,true)`).
   - Molde da composição+CLI: `composition/catalog_wiring.go:24-30`,
     `composition/catalog_ingest.go:28-67` (`RunCatalogIngest` com report por disposição),
     `cmd/catalogingest/main.go:43-96` (fail-closed de tenant em `:118-123`).
   - Mecânica ML já provada no legado: enumeração scan
     (`modules/connectors/adapters/mercado_livre/items_scan_ids_reader.go:30-63` —
     `GET /users/{id}/items/search?search_type=scan&scroll_id=...`, resposta
     `{scroll_id, results[]}`), hidratação multiget
     (`items_multiget_reader.go:209-260` — `GET /items?ids=`, lote de 20, elementos
     `{code, body}` na ordem do pedido), campos do body medidos em
     `items_multiget_reader.go:143-181` (`id,title,seller_sku,attributes,status,price
     (json.Number),currency_id,listing_type_id,available_quantity,variations[...]`),
     EAN = atributo `GTIN` com fallback `EAN`, `value_name`
     (`modules/listings/adapters/connectors/multiget_mapper.go:197,210`).
   - Resolução de credencial ML fora do módulo integrations:
     `cmd/mlprobe/main.go:164-190` — SQL sobre `integration_installations` +
     `integration_credentials` (provider `mercado_livre`, status `connected`, credencial
     ativa não revogada, `ORDER BY version DESC LIMIT 1`) + decrypt AES-GCM via
     `modules/integrations/adapters/crypto.NewLocalKeyService(key,"local-key-v1").DecryptJSON`;
     chaves de ambiente `MC_DATABASE_URL` + `MPC_ENCRYPTION_KEY` (`main.go:76-79`).
   - Kernel: `channel.Code/AccountRef` (`kernel/channel/channel.go:18-59`),
     `exact.Money/ParseMoney/Currency` (`kernel/exact/money.go:19-72`),
     `fact.Fact[T]` com `NewKnown/NewUnknown/NewEstimated/NewNotApplicable`
     (`kernel/fact/knowledge.go:64-151`), `provenance.NewEvidence(system, objectKind,
     externalKey, observedAt, payloadHash)` (`kernel/provenance/evidence.go:30`).
2. **Onde vive o gap?** Não existe `contexts/listings` nem `adapters/marketplace/`
   (universo: `internal/contexts/` tem só `catalog`; `internal/adapters/` tem só
   `erp/sankhyaoracle` — `find` desta sessão). O legado `modules/listings` (44 arquivos)
   está fora do molde e acoplado por observer a `product_links`
   (`composition/root.go:824`).
3. **Quem consome os caminhos alterados?** Ninguém ainda — tudo é árvore nova. O legado
   continua intacto e rodando (scheduler em `root.go:841-843`); nada deste plano o toca.
   Sítio de composição do código novo: `cmd/listingsingest` (Task 8), único chamador de
   produção.
4. **Contrato existente?** Nenhum endpoint HTTP novo — sem OpenAPI/SDK nesta perna (CLI de
   operador, como `catalogingest`). Registro de governança: entrada `kind:"context"` em
   `contracts/governance/modules.json` (molde do catalog em `:188-195`).
5. **Estado vivo?** Instalação ML conectada existe (mlprobe rodou contra ela; tabela
   `listings` legada tem anúncios — `cmd/mlprobe/main.go:192-196` consulta
   `status='active'`). A contagem exata é medida no live drive (Task 9), não recordada.
6. **O que prova quebrado hoje / consertado depois?** Hoje:
   `SELECT count(*) FROM listings.listings` falha — o schema não existe. Depois:
   `count(*)` > 0 com o mesmo universo de anúncios que o legado vê para a mesma
   instalação (reconciliação na Task 9).
7. **Custo (bucket de rate ML compartilhado):** enumeração = ⌈N/50⌉ GETs de scan + 1 GET
   terminal; hidratação = ⌈N/20⌉ GETs multiget. Para N≈34 anúncios: ~2+2+1 = 5 GETs por
   corrida completa. Sem POST/PUT. CLI é disparo manual do operador — zero chamadas
   recorrentes.
8. **O que falha silencioso às 3h?** Nada novo agendado — a CLI é manual e morre com
   stderr + exit≠0 (molde `catalogingest/main.go:34-40`). Item com `code!=200` no multiget
   FALHA a página nomeando o id (nunca é pulado em silêncio). A dívida "sem agendamento e
   sem visibilidade de staleness na tela" fica registrada na Task 10 — é da perna do
   scheduler (`platform/scheduler`), não desta.

## O que já existe (Fase 2 — por que não serve)

| Artefato próximo | Por que não serve |
|---|---|
| `modules/listings` (44 arquivos) | É o legado a substituir: DTO de vendor atravessa camadas (`listings/adapters/connectors/backfill.go:9` importa `ItemMultigetDTO` de `connectors`), acoplado a `product_links` por observer. Não se estende legado (Fase 3 pergunta 4) |
| `modules/connectors/adapters/mercado_livre` | Mecânica correta, formato errado: DTOs exportados para fora do vendor (o que §2.2 proíbe); depende de `integrations` para token por dentro do módulo. Vira REFERÊNCIA de mecânica, nunca import |
| `contexts/catalog` | É o molde a replicar, não a estender — listings é outro contexto com outro dono de dados |
| `cmd/mlprobe` | Prova a resolução de credencial e as respostas reais do ML; é diagnóstico descartável, não caminho de produção |

**Local vs global (Fase 3):** cópias do conceito "observação de anúncio ML" existem 2× no
legado (`listings` + snapshots de `product_links`, alimentadas pelo MESMO pull). O global
fix é o contexto novo único; o acoplamento gêmeo se desfaz quando `linking` aterrar (perna
seguinte). Extensão de legado: zero — o plano não toca `internal/modules/`.

**Seams (Fase 4):** árvore nova inteira (sem colisão; `git worktree list` = só main nesta
sessão, nenhum branch de código em voo). Migrations é seam compartilhado: prefixos novos
`0099`/`0100` únicos (último existente: `0098`; há prefixo duplicado histórico `0093` — o
fixture conta por nome completo, `migrate/runner_test.go:39-46`). Sem OpenAPI/SDK. Tenant:
todas as tabelas novas com `tenant_id` na PK + RLS FORCE (molde `0097`).

---

### Task 1: `contexts/listings/contracts` — a observação e seu resultado

**Files:**
- Create: `apps/server_core/internal/contexts/listings/contracts/key.go`
- Create: `apps/server_core/internal/contexts/listings/contracts/observation.go`
- Create: `apps/server_core/internal/contexts/listings/contracts/contracts_test.go`

- [ ] **Step 1: RED — testes primeiro**

`contracts_test.go`:

```go
package contracts_test

import (
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/kernel/channel"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func testEvidence(t *testing.T) provenance.Evidence {
	t.Helper()
	e, err := provenance.NewEvidence("mercado_livre", "item", "MLB123", time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC), "hash-1")
	if err != nil {
		t.Fatalf("evidence: %v", err)
	}
	return e
}

func testKey(t *testing.T) contracts.SourceListingKey {
	t.Helper()
	tid, err := tenant.Parse("tenant_default")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	code, err := channel.ParseCode("mercado_livre")
	if err != nil {
		t.Fatalf("channel: %v", err)
	}
	account, err := channel.NewAccountRef(code, "179571326")
	if err != nil {
		t.Fatalf("account: %v", err)
	}
	key, err := contracts.NewSourceListingKey(tid, account, "MLB123")
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	return key
}

func TestNewSourceListingKeyRejectsBlankListingID(t *testing.T) {
	tid, _ := tenant.Parse("tenant_default")
	code, _ := channel.ParseCode("mercado_livre")
	account, _ := channel.NewAccountRef(code, "179571326")
	if _, err := contracts.NewSourceListingKey(tid, account, "  "); err == nil {
		t.Fatal("blank listing id accepted")
	}
}

func TestValidateRejectsZeroKeyAndZeroEvidence(t *testing.T) {
	title, _ := fact.NewUnknown[string]("ml omitted title", testEvidence(t))
	obs := contracts.ListingObservation{Title: title}
	if err := obs.Validate(); err == nil || !strings.Contains(err.Error(), "key") {
		t.Fatalf("zero key accepted or wrong error: %v", err)
	}
	obs.Key = testKey(t)
	if err := obs.Validate(); err == nil || !strings.Contains(err.Error(), "evidence") {
		t.Fatalf("zero evidence accepted or wrong error: %v", err)
	}
}

func TestValidateRejectsVariationWithBlankID(t *testing.T) {
	obs := contracts.ListingObservation{
		Key:      testKey(t),
		Evidence: testEvidence(t),
		Variations: []contracts.VariationObservation{{VariationID: ""}},
	}
	if err := obs.Validate(); err == nil || !strings.Contains(err.Error(), "variation") {
		t.Fatalf("blank variation id accepted or wrong error: %v", err)
	}
}

func TestValidateAcceptsAllFactsUnknown(t *testing.T) {
	// Um anúncio de que o ML só devolveu o id é um fato sobre o ML; Listings
	// grava Unknown, nunca recusa nem inventa (protocolo §4.1).
	obs := contracts.ListingObservation{Key: testKey(t), Evidence: testEvidence(t)}
	if err := obs.Validate(); err != nil {
		t.Fatalf("all-unknown observation rejected: %v", err)
	}
}
```

- [ ] **Step 2: Rodar e ver a falha de compilação nomeada**

De `apps/server_core` (caches absolutos, restrição global 1):

```bash
go test ./internal/contexts/listings/contracts/ -count=1
```

Esperado: FAIL — `no required module provides package` / `undefined: contracts.…`.

- [ ] **Step 3: Implementar**

`key.go`:

```go
// Package contracts is Listings' published vocabulary: what a channel adapter
// hands in, and what ingesting it did. Nothing here names a vendor type.
package contracts

import (
	"errors"
	"fmt"
	"strings"

	"marketplace-central/apps/server_core/internal/kernel/channel"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// ErrBlank marks a required field that arrived empty.
var ErrBlank = errors.New("listings contracts: required field is blank")

// SourceListingKey is the identity of one listing as one channel account
// publishes it. Listings does not mint ids: the provider's listing id IS the
// identity, scoped by tenant and channel account.
type SourceListingKey struct {
	tenant    tenant.ID
	account   channel.AccountRef
	listingID string
}

// NewSourceListingKey validates every component; a key with a blank part would
// silently collide rows across accounts.
func NewSourceListingKey(t tenant.ID, account channel.AccountRef, listingID string) (SourceListingKey, error) {
	if t.IsZero() {
		return SourceListingKey{}, fmt.Errorf("%w: tenant", ErrBlank)
	}
	if account.Channel().IsZero() {
		return SourceListingKey{}, fmt.Errorf("%w: channel account", ErrBlank)
	}
	listingID = strings.TrimSpace(listingID)
	if listingID == "" {
		return SourceListingKey{}, fmt.Errorf("%w: listing id", ErrBlank)
	}
	return SourceListingKey{tenant: t, account: account, listingID: listingID}, nil
}

func (k SourceListingKey) Tenant() tenant.ID           { return k.tenant }
func (k SourceListingKey) Account() channel.AccountRef { return k.account }
func (k SourceListingKey) ListingID() string           { return k.listingID }
func (k SourceListingKey) IsZero() bool                { return k.listingID == "" }
```

`observation.go`:

```go
package contracts

import (
	"fmt"

	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

// VariationObservation is one variation as the channel reported it. SellerSKU
// and GTIN travel here because the linking context (next leg) anchors on them;
// dropping them now would force a second pull later.
type VariationObservation struct {
	VariationID       string
	Price             fact.Fact[exact.Money]
	AvailableQuantity fact.Fact[int]
	SellerSKU         fact.Fact[string]
	GTIN              fact.Fact[string]
}

// ListingObservation is what a channel adapter hands to Listings: one listing
// as the channel saw it at one moment. Every field the channel may omit is a
// Fact, never a zero value (protocolo §4.1).
//
// RawPayload is the channel's own bytes for this observation, opaque to this
// context: stored for audit and reconciliation (§15.3 "payloads reais
// gravados"), never parsed past the adapter.
type ListingObservation struct {
	Key               SourceListingKey
	Title             fact.Fact[string]
	Status            fact.Fact[string]
	ListingType       fact.Fact[string]
	Price             fact.Fact[exact.Money]
	AvailableQuantity fact.Fact[int]
	SellerSKU         fact.Fact[string]
	GTIN              fact.Fact[string]
	Variations        []VariationObservation
	RawPayload        []byte
	Evidence          provenance.Evidence
}

// Validate rejects an observation that cannot be recorded. It deliberately
// accepts every fact as Unknown: a channel that said nothing is a fact about
// the channel.
func (o ListingObservation) Validate() error {
	if o.Key.IsZero() {
		return fmt.Errorf("%w: key", ErrBlank)
	}
	if o.Evidence.IsZero() {
		return fmt.Errorf("%w: evidence", ErrBlank)
	}
	for i, v := range o.Variations {
		if v.VariationID == "" {
			return fmt.Errorf("%w: variation id at index %d", ErrBlank, i)
		}
	}
	return nil
}

// Disposition is what ingesting an observation did — same closed set the
// catalog leg ratified (contexts/catalog/contracts/observation.go:39-49).
type Disposition string

const (
	DispositionCreated    Disposition = "created"
	DispositionChanged    Disposition = "changed"
	DispositionIdempotent Disposition = "idempotent"
)

// IngestResult reports what happened to one observation.
type IngestResult struct {
	Disposition Disposition
	Version     int
}
```

- [ ] **Step 4: Verde**

```bash
go test ./internal/contexts/listings/contracts/ -count=1 -v
```

Esperado: PASS nos 4 testes, cada um com `--- PASS:` nomeado.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/contexts/listings/contracts
git commit -m "feat(listings): contracts -- observation vocabulary with explicit knowledge states"
```

---

### Task 2: `contexts/listings/port` — o feed que o adapter implementa

**Files:**
- Create: `apps/server_core/internal/contexts/listings/port/feed.go`
- Create: `apps/server_core/internal/contexts/listings/port/feed_test.go`

Molde: `contexts/catalog/port/feed.go` inteiro — cursor opaco pela MESMA razão medida lá
(o scroll_id do ML é exatamente o token que o contexto não pode conhecer). Sem port de
leitura nesta perna: nenhum consumidor existe ainda, e superfície sem chamador é o
anti-padrão "orphan contract operation" (Fase 2). `linking` cria o Reader quando precisar.

- [ ] **Step 1: RED**

`feed_test.go`:

```go
package port_test

import (
	"testing"

	"marketplace-central/apps/server_core/internal/contexts/listings/port"
)

func TestZeroCursorIsStart(t *testing.T) {
	var c port.Cursor
	if !c.IsStart() {
		t.Fatal("zero cursor must be the start of the feed")
	}
}

func TestCursorRoundTripsToken(t *testing.T) {
	c := port.NewCursor("scroll-abc")
	if c.IsStart() || c.Token() != "scroll-abc" {
		t.Fatalf("cursor lost its token: start=%v token=%q", c.IsStart(), c.Token())
	}
}
```

```bash
go test ./internal/contexts/listings/port/ -count=1
```

Esperado: FAIL de compilação (`undefined: port.Cursor`).

- [ ] **Step 2: Implementar**

`feed.go`:

```go
// Package port carries what Listings asks of a channel source. The cursor is
// opaque on purpose: ML pages by scroll_id, another channel may page by
// timestamp, and neither shape belongs in this contract (the legacy port that
// typed one source's row id into the cursor is the measured counter-example,
// contexts/catalog/port/feed.go:14-18).
package port

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// Cursor is the source's own position marker, unreadable by Listings.
type Cursor struct{ token string }

func NewCursor(token string) Cursor { return Cursor{token: token} }
func (c Cursor) Token() string      { return c.token }
func (c Cursor) IsStart() bool      { return c.token == "" }

// Page is one batch of observations plus where to continue. Done is explicit:
// an empty page mid-feed is legal and must not stop the walk.
type Page struct {
	Observations []contracts.ListingObservation
	Next         Cursor
	Done         bool
}

// ListingFeed is a source of listing observations. Listings asks; the adapter
// decides how to page and how to authenticate.
type ListingFeed interface {
	NextPage(ctx context.Context, t tenant.ID, after Cursor, limit int) (Page, error)
}
```

- [ ] **Step 3: Verde + commit**

```bash
go test ./internal/contexts/listings/port/ -count=1 -v
```

Esperado: `--- PASS: TestZeroCursorIsStart`, `--- PASS: TestCursorRoundTripsToken`.

```bash
git add apps/server_core/internal/contexts/listings/port
git commit -m "feat(listings): port -- opaque-cursor listing feed"
```

---

### Task 3: `internal/application` — o caso de uso de ingestão

**Files:**
- Create: `apps/server_core/internal/contexts/listings/internal/application/ports.go`
- Create: `apps/server_core/internal/contexts/listings/internal/application/ingest.go`
- Create: `apps/server_core/internal/contexts/listings/internal/application/ingest_test.go`
- Create: `apps/server_core/internal/contexts/listings/internal/application/memstore_test.go`

Sem pacote `domain` separado nesta perna: a decisão inteira é o fold
criado/mudado/idempotente por hash — 3 ramos. O catalog precisou de `domain` porque cunha
identidade (ULID + resolução de identificadores); Listings não cunha nada. Se a perna
`linking` exigir invariantes de anúncio, o pacote nasce lá com elas (YAGNI, Fase 3).

- [ ] **Step 1: RED**

`ingest_test.go`:

```go
package application_test

import (
	"context"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/application"
	"marketplace-central/apps/server_core/internal/kernel/channel"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func obs(t *testing.T, listingID, payloadHash string) contracts.ListingObservation {
	t.Helper()
	tid, _ := tenant.Parse("tenant_default")
	code, _ := channel.ParseCode("mercado_livre")
	account, _ := channel.NewAccountRef(code, "179571326")
	key, err := contracts.NewSourceListingKey(tid, account, listingID)
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	ev, err := provenance.NewEvidence("mercado_livre", "item", listingID, time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC), payloadHash)
	if err != nil {
		t.Fatalf("evidence: %v", err)
	}
	return contracts.ListingObservation{Key: key, Evidence: ev, RawPayload: []byte(`{"id":"` + listingID + `"}`)}
}

func TestFirstObservationCreatesVersionOne(t *testing.T) {
	svc := application.NewService(newMemStore())
	got, err := svc.Ingest(context.Background(), obs(t, "MLB1", "h1"))
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if got.Disposition != contracts.DispositionCreated || got.Version != 1 {
		t.Fatalf("got %+v, want created v1", got)
	}
}

func TestSamePayloadHashIsIdempotent(t *testing.T) {
	svc := application.NewService(newMemStore())
	if _, err := svc.Ingest(context.Background(), obs(t, "MLB1", "h1")); err != nil {
		t.Fatalf("first ingest: %v", err)
	}
	got, err := svc.Ingest(context.Background(), obs(t, "MLB1", "h1"))
	if err != nil {
		t.Fatalf("second ingest: %v", err)
	}
	if got.Disposition != contracts.DispositionIdempotent || got.Version != 1 {
		t.Fatalf("got %+v, want idempotent v1", got)
	}
}

func TestChangedPayloadMintsNewVersion(t *testing.T) {
	svc := application.NewService(newMemStore())
	if _, err := svc.Ingest(context.Background(), obs(t, "MLB1", "h1")); err != nil {
		t.Fatalf("first ingest: %v", err)
	}
	got, err := svc.Ingest(context.Background(), obs(t, "MLB1", "h2"))
	if err != nil {
		t.Fatalf("second ingest: %v", err)
	}
	if got.Disposition != contracts.DispositionChanged || got.Version != 2 {
		t.Fatalf("got %+v, want changed v2", got)
	}
}

func TestInvalidObservationNeverReachesTheStore(t *testing.T) {
	store := newMemStore()
	svc := application.NewService(store)
	bad := contracts.ListingObservation{} // zero key, zero evidence
	if _, err := svc.Ingest(context.Background(), bad); err == nil {
		t.Fatal("invalid observation accepted")
	}
	if store.saves != 0 {
		t.Fatalf("store touched %d times by an invalid observation", store.saves)
	}
}
```

`memstore_test.go` (dublê em memória do port INTERNO — prova o contrato do caso de uso,
nunca a integração; a integração é a Task 6):

```go
package application_test

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/application"
)

type memStore struct {
	current map[string]application.CurrentListing
	saves   int
}

func newMemStore() *memStore {
	return &memStore{current: map[string]application.CurrentListing{}}
}

func memKey(k contracts.SourceListingKey) string {
	return k.Tenant().String() + "|" + k.Account().Channel().String() + "|" + k.Account().External() + "|" + k.ListingID()
}

func (m *memStore) Current(_ context.Context, k contracts.SourceListingKey) (application.CurrentListing, bool, error) {
	c, ok := m.current[memKey(k)]
	return c, ok, nil
}

func (m *memStore) SaveVersion(_ context.Context, o contracts.ListingObservation, version int) error {
	m.saves++
	m.current[memKey(o.Key)] = application.CurrentListing{Version: version, PayloadHash: o.Evidence.PayloadHash()}
	return nil
}
```

```bash
go test ./internal/contexts/listings/internal/application/ -count=1
```

Esperado: FAIL de compilação (`undefined: application.NewService`).

- [ ] **Step 2: Implementar**

`ports.go`:

```go
// Package application owns the ingest decision. Its Store port is internal:
// only the postgres adapter inside this context implements it.
package application

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
)

// CurrentListing is what the store knows about a listing before an ingest.
type CurrentListing struct {
	Version     int
	PayloadHash string
}

// Store persists listing versions. SaveVersion must be atomic: the listing
// row, its variations and the observation land together or not at all.
type Store interface {
	Current(ctx context.Context, key contracts.SourceListingKey) (CurrentListing, bool, error)
	SaveVersion(ctx context.Context, o contracts.ListingObservation, version int) error
}
```

`ingest.go`:

```go
package application

import (
	"context"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
)

// Service is the one use case this leg ships: fold a channel observation in.
type Service struct {
	store Store
}

func NewService(store Store) *Service { return &Service{store: store} }

// Ingest decides created/changed/idempotent by payload hash. Re-polling an
// unchanged channel must be free (same rule the catalog leg ratified).
func (s *Service) Ingest(ctx context.Context, o contracts.ListingObservation) (contracts.IngestResult, error) {
	if err := o.Validate(); err != nil {
		return contracts.IngestResult{}, err
	}
	current, exists, err := s.store.Current(ctx, o.Key)
	if err != nil {
		return contracts.IngestResult{}, fmt.Errorf("listings: read current: %w", err)
	}
	switch {
	case !exists:
		if err := s.store.SaveVersion(ctx, o, 1); err != nil {
			return contracts.IngestResult{}, fmt.Errorf("listings: save v1: %w", err)
		}
		return contracts.IngestResult{Disposition: contracts.DispositionCreated, Version: 1}, nil
	case current.PayloadHash == o.Evidence.PayloadHash():
		return contracts.IngestResult{Disposition: contracts.DispositionIdempotent, Version: current.Version}, nil
	default:
		next := current.Version + 1
		if err := s.store.SaveVersion(ctx, o, next); err != nil {
			return contracts.IngestResult{}, fmt.Errorf("listings: save v%d: %w", next, err)
		}
		return contracts.IngestResult{Disposition: contracts.DispositionChanged, Version: next}, nil
	}
}
```

- [ ] **Step 3: Verde + commit**

```bash
go test ./internal/contexts/listings/internal/application/ -count=1 -v
```

Esperado: 4× `--- PASS:`.

```bash
git add apps/server_core/internal/contexts/listings/internal
git commit -m "feat(listings): ingest use case -- created/changed/idempotent by payload hash"
```

---

### Task 4: migrations `0099`/`0100` + fixture de contagem

**Files:**
- Create: `apps/server_core/migrations/0099_listings_context.sql`
- Create: `apps/server_core/migrations/0100_listings_app_role.sql`
- Modify: `apps/server_core/internal/platform/migrate/runner_test.go` (contagem 85→87 nos
  DOIS pontos: `:25` e `:64`)

- [ ] **Step 1: RED — bump da contagem primeiro**

Em `runner_test.go`, trocar `85` por `87` em `:25` e `:64` (o fixture conta nomes de
arquivo completos — `runner_test.go:39-46`). Rodar:

```bash
go test ./internal/platform/migrate/ -count=1
```

Esperado: FAIL `fixture inventory drift: got 85 canonical migrations, want 87`.

- [ ] **Step 2: `0099_listings_context.sql`** (molde `0097`: schema próprio, tenant na
  frente de toda PK, RLS FORCE, tripla estado/valor/razão por fato):

```sql
-- Listings context: its own schema, its own writer, no foreign key leaving it.
-- Identity is the provider's listing id scoped by tenant + channel account;
-- Listings mints nothing. Every fact column is a state/value/reason triple
-- because "ML said nothing" and "ML said zero" are different facts (0097 is
-- the ratified precedent).
CREATE SCHEMA IF NOT EXISTS listings;

CREATE TABLE IF NOT EXISTS listings.listings (
    tenant_id           text        NOT NULL,
    channel             text        NOT NULL,
    account_external_id text        NOT NULL,
    listing_id          text        NOT NULL,
    version             integer     NOT NULL,
    title_state         text        NOT NULL,
    title_value         text        NULL,
    title_reason        text        NULL,
    status_state        text        NOT NULL,
    status_value        text        NULL,
    status_reason       text        NULL,
    listing_type_state  text        NOT NULL,
    listing_type_value  text        NULL,
    listing_type_reason text        NULL,
    price_state         text        NOT NULL,
    price_amount        text        NULL,
    price_currency      text        NULL,
    price_reason        text        NULL,
    available_qty_state  text       NOT NULL,
    available_qty_value  integer    NULL,
    available_qty_reason text       NULL,
    seller_sku_state    text        NOT NULL,
    seller_sku_value    text        NULL,
    seller_sku_reason   text        NULL,
    gtin_state          text        NOT NULL,
    gtin_value          text        NULL,
    gtin_reason         text        NULL,
    last_payload_hash   text        NOT NULL,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT listings_pkey PRIMARY KEY (tenant_id, channel, account_external_id, listing_id),
    CONSTRAINT listings_version_positive CHECK (version >= 1),
    CONSTRAINT listings_title_state CHECK (title_state IN ('known','estimated','unknown','not_applicable')),
    CONSTRAINT listings_status_state CHECK (status_state IN ('known','estimated','unknown','not_applicable')),
    CONSTRAINT listings_listing_type_state CHECK (listing_type_state IN ('known','estimated','unknown','not_applicable')),
    CONSTRAINT listings_price_state CHECK (price_state IN ('known','estimated','unknown','not_applicable')),
    CONSTRAINT listings_available_qty_state CHECK (available_qty_state IN ('known','estimated','unknown','not_applicable')),
    CONSTRAINT listings_seller_sku_state CHECK (seller_sku_state IN ('known','estimated','unknown','not_applicable')),
    CONSTRAINT listings_gtin_state CHECK (gtin_state IN ('known','estimated','unknown','not_applicable')),
    -- known carries a value; unknown carries a reason and no value (0097:26-31).
    CONSTRAINT listings_title_consistent CHECK (
        (title_state = 'known' AND title_value IS NOT NULL)
     OR (title_state = 'estimated' AND title_value IS NOT NULL AND title_reason IS NOT NULL)
     OR (title_state IN ('unknown','not_applicable') AND title_value IS NULL AND title_reason IS NOT NULL)),
    CONSTRAINT listings_price_consistent CHECK (
        (price_state = 'known' AND price_amount IS NOT NULL AND price_currency IS NOT NULL)
     OR (price_state = 'estimated' AND price_amount IS NOT NULL AND price_currency IS NOT NULL AND price_reason IS NOT NULL)
     OR (price_state IN ('unknown','not_applicable') AND price_amount IS NULL AND price_currency IS NULL AND price_reason IS NOT NULL))
);

CREATE TABLE IF NOT EXISTS listings.listing_variations (
    tenant_id           text    NOT NULL,
    channel             text    NOT NULL,
    account_external_id text    NOT NULL,
    listing_id          text    NOT NULL,
    variation_id        text    NOT NULL,
    price_state         text    NOT NULL,
    price_amount        text    NULL,
    price_currency      text    NULL,
    price_reason        text    NULL,
    available_qty_state  text   NOT NULL,
    available_qty_value  integer NULL,
    available_qty_reason text   NULL,
    seller_sku_state    text    NOT NULL,
    seller_sku_value    text    NULL,
    seller_sku_reason   text    NULL,
    gtin_state          text    NOT NULL,
    gtin_value          text    NULL,
    gtin_reason         text    NULL,
    CONSTRAINT listing_variations_pkey
        PRIMARY KEY (tenant_id, channel, account_external_id, listing_id, variation_id),
    CONSTRAINT listing_variations_listing_fkey
        FOREIGN KEY (tenant_id, channel, account_external_id, listing_id)
        REFERENCES listings.listings (tenant_id, channel, account_external_id, listing_id)
        ON DELETE CASCADE
);

-- One row per distinct payload the channel ever showed us: the real bytes,
-- kept for reconciliation and for the linking leg (§15.3).
CREATE TABLE IF NOT EXISTS listings.source_observations (
    tenant_id           text        NOT NULL,
    channel             text        NOT NULL,
    account_external_id text        NOT NULL,
    listing_id          text        NOT NULL,
    payload_hash        text        NOT NULL,
    payload             jsonb       NOT NULL,
    observed_at         timestamptz NOT NULL,
    recorded_at         timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT source_observations_pkey
        PRIMARY KEY (tenant_id, channel, account_external_id, listing_id, payload_hash),
    CONSTRAINT source_observations_listing_fkey
        FOREIGN KEY (tenant_id, channel, account_external_id, listing_id)
        REFERENCES listings.listings (tenant_id, channel, account_external_id, listing_id)
        ON DELETE CASCADE
);

ALTER TABLE listings.listings            ENABLE ROW LEVEL SECURITY;
ALTER TABLE listings.listing_variations  ENABLE ROW LEVEL SECURITY;
ALTER TABLE listings.source_observations ENABLE ROW LEVEL SECURITY;
ALTER TABLE listings.listings            FORCE ROW LEVEL SECURITY;
ALTER TABLE listings.listing_variations  FORCE ROW LEVEL SECURITY;
ALTER TABLE listings.source_observations FORCE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation ON listings.listings;
CREATE POLICY tenant_isolation ON listings.listings
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS tenant_isolation ON listings.listing_variations;
CREATE POLICY tenant_isolation ON listings.listing_variations
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
DROP POLICY IF EXISTS tenant_isolation ON listings.source_observations;
CREATE POLICY tenant_isolation ON listings.source_observations
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
```

Nota de escopo, deliberada: `sold_quantity`, `permalink`, `condition`, `category_id`,
`tags`, `shipping`, `thumbnail`, `catalog_product_id` ficam FORA das colunas — o payload
bruto os retém em `source_observations.payload`, e a perna que precisar deles (pricing,
intelligence) os promove a coluna com dono. Coluna sem consumidor é superfície órfã.

- [ ] **Step 3: `0100_listings_app_role.sql`** (molde `0098` — o role `mpc_app` já existe;
  só USAGE + grants no schema novo):

```sql
-- Same reasoning as 0098: RLS FORCEd in 0099 is decorative until the
-- application role that cannot bypass it is granted the schema. mpc_app
-- already exists (0098); this only extends it to the listings schema.
GRANT USAGE ON SCHEMA listings TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON listings.listings            TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON listings.listing_variations  TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON listings.source_observations TO mpc_app;
-- No GRANT on any other schema (0098:24-25).
```

- [ ] **Step 4: Verde no fixture + commit**

```bash
go test ./internal/platform/migrate/ -count=1 -v
```

Esperado: PASS, contagem 87.

```bash
git add apps/server_core/migrations apps/server_core/internal/platform/migrate/runner_test.go
git commit -m "feat(listings): schema with RLS-forced tenant isolation and raw observation retention"
```

---

### Task 5: repositório Postgres + fachada `module.go` + lane hermética

**Files:**
- Create: `apps/server_core/internal/contexts/listings/internal/postgres/repository.go`
- Create: `apps/server_core/internal/contexts/listings/module.go`
- Create: `apps/server_core/tests/integration/listings_ingest_test.go`

- [ ] **Step 1: RED — teste de integração primeiro** (molde
  `tests/integration/catalog_ingest_composition_test.go`: build tag `integration`, Postgres
  real da lane hermética; copiar o cabeçalho de setup DAQUELE arquivo — pool, migrations,
  tenant — que é infraestrutura da lane, e trocar o corpo):

```go
//go:build integration

package integration

// Segue o setup de catalog_ingest_composition_test.go (pool + migrations + tenant).
// Corpo do teste:

func TestListingsIngestPersistsVersionAndVariations(t *testing.T) {
	// setup: pool, tid — idêntico ao teste de catalog da lane.
	module := listings.New(pool)

	first := listingObservation(t, tid, "MLB-IT-1", "hash-a") // helper local: observação completa com 1 variação conhecida
	res, err := module.IngestListing(ctx, first)
	if err != nil {
		t.Fatalf("ingest v1: %v", err)
	}
	if res.Disposition != contracts.DispositionCreated || res.Version != 1 {
		t.Fatalf("got %+v, want created v1", res)
	}

	// A aceitação é o BANCO, não o retorno (mc-planning: aceite = observável).
	var rows int
	if err := scanOne(pool, tid, &rows,
		`SELECT count(*) FROM listings.listings WHERE tenant_id=$1 AND listing_id='MLB-IT-1'`); err != nil {
		t.Fatalf("count listings: %v", err)
	}
	if rows != 1 {
		t.Fatalf("listings rows = %d, want 1", rows)
	}
	if err := scanOne(pool, tid, &rows,
		`SELECT count(*) FROM listings.listing_variations WHERE tenant_id=$1 AND listing_id='MLB-IT-1'`); err != nil {
		t.Fatalf("count variations: %v", err)
	}
	if rows != 1 {
		t.Fatalf("variation rows = %d, want 1", rows)
	}

	second := listingObservation(t, tid, "MLB-IT-1", "hash-b")
	res, err = module.IngestListing(ctx, second)
	if err != nil {
		t.Fatalf("ingest v2: %v", err)
	}
	if res.Disposition != contracts.DispositionChanged || res.Version != 2 {
		t.Fatalf("got %+v, want changed v2", res)
	}
	if err := scanOne(pool, tid, &rows,
		`SELECT count(*) FROM listings.source_observations WHERE tenant_id=$1 AND listing_id='MLB-IT-1'`); err != nil {
		t.Fatalf("count observations: %v", err)
	}
	if rows != 2 {
		t.Fatalf("observation rows = %d, want 2 (one per distinct payload)", rows)
	}
}
```

O helper `listingObservation` constrói a observação com `fact.NewKnown` para
title/price/qty/sku/gtin e uma `VariationObservation` conhecida; `scanOne` executa a query
dentro de uma transação com `SELECT set_config('app.tenant_id',$1,true)` — RLS FORÇADO
significa que um count sem o setting devolve 0 e o teste ficaria cego (padrão de
`repository.go:49` do catalog). Escrever os dois helpers completos no arquivo.

A lane hermética só compila `./tests/integration` com tag `integration`
(memória: lane não é superconjunto da de unidade) — o teste roda na Task 9 junto com o
resto, ou agora se o Docker local estiver de pé:

```bash
pwsh -NoProfile -File ./scripts/gate.ps1 -Lane integration
```

Esperado AGORA: FAIL de compilação nomeando `listings.New` (o RED da task).

- [ ] **Step 2: `repository.go`** — implementa `application.Store`:

```go
// Package postgres is Listings' own writer. Every statement runs inside a
// transaction that first pins app.tenant_id, because RLS is FORCEd and a
// query without the setting sees an empty world (catalog repository.go:49 is
// the ratified precedent).
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/application"
	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
)

type Repository struct{ pool *pgxpool.Pool }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// factColumns flattens a Fact[string] into its state/value/reason triple.
func factColumns(f fact.Fact[string]) (state string, value, reason *string) {
	state = f.State().String()
	if v, ok := f.Value(); ok {
		value = &v
	}
	if r := f.Reason(); r != "" {
		reason = &r
	}
	return state, value, reason
}

func intFactColumns(f fact.Fact[int]) (state string, value *int, reason *string) {
	state = f.State().String()
	if v, ok := f.Value(); ok {
		value = &v
	}
	if r := f.Reason(); r != "" {
		reason = &r
	}
	return state, value, reason
}

func moneyFactColumns(f fact.Fact[exact.Money]) (state string, amount, currency, reason *string) {
	state = f.State().String()
	if v, ok := f.Value(); ok {
		// StringFixed, not String: Decimal has no String(), and persistence
		// renders at the column's scale (kernel/exact/money.go:97).
		a := v.Amount().StringFixed(2)
		c := v.Currency().String()
		amount, currency = &a, &c
	}
	if r := f.Reason(); r != "" {
		reason = &r
	}
	return state, amount, currency, reason
}

func (r *Repository) withTenantTx(ctx context.Context, tenantID string, fn func(pgx.Tx) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("listings postgres: begin: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", tenantID); err != nil {
		return fmt.Errorf("listings postgres: set tenant: %w", err)
	}
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repository) Current(ctx context.Context, key contracts.SourceListingKey) (application.CurrentListing, bool, error) {
	var current application.CurrentListing
	found := false
	err := r.withTenantTx(ctx, key.Tenant().String(), func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx,
			`SELECT version, last_payload_hash FROM listings.listings
			 WHERE tenant_id=$1 AND channel=$2 AND account_external_id=$3 AND listing_id=$4`,
			key.Tenant().String(), key.Account().Channel().String(), key.Account().External(), key.ListingID())
		switch err := row.Scan(&current.Version, &current.PayloadHash); err {
		case nil:
			found = true
			return nil
		case pgx.ErrNoRows:
			return nil
		default:
			return fmt.Errorf("listings postgres: read current: %w", err)
		}
	})
	return current, found, err
}

func (r *Repository) SaveVersion(ctx context.Context, o contracts.ListingObservation, version int) error {
	k := o.Key
	return r.withTenantTx(ctx, k.Tenant().String(), func(tx pgx.Tx) error {
		titleS, titleV, titleR := factColumns(o.Title)
		statusS, statusV, statusR := factColumns(o.Status)
		typeS, typeV, typeR := factColumns(o.ListingType)
		priceS, priceA, priceC, priceR := moneyFactColumns(o.Price)
		qtyS, qtyV, qtyR := intFactColumns(o.AvailableQuantity)
		skuS, skuV, skuR := factColumns(o.SellerSKU)
		gtinS, gtinV, gtinR := factColumns(o.GTIN)
		if _, err := tx.Exec(ctx,
			`INSERT INTO listings.listings (
			   tenant_id, channel, account_external_id, listing_id, version,
			   title_state, title_value, title_reason,
			   status_state, status_value, status_reason,
			   listing_type_state, listing_type_value, listing_type_reason,
			   price_state, price_amount, price_currency, price_reason,
			   available_qty_state, available_qty_value, available_qty_reason,
			   seller_sku_state, seller_sku_value, seller_sku_reason,
			   gtin_state, gtin_value, gtin_reason,
			   last_payload_hash)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28)
			 ON CONFLICT (tenant_id, channel, account_external_id, listing_id) DO UPDATE SET
			   version=EXCLUDED.version,
			   title_state=EXCLUDED.title_state, title_value=EXCLUDED.title_value, title_reason=EXCLUDED.title_reason,
			   status_state=EXCLUDED.status_state, status_value=EXCLUDED.status_value, status_reason=EXCLUDED.status_reason,
			   listing_type_state=EXCLUDED.listing_type_state, listing_type_value=EXCLUDED.listing_type_value, listing_type_reason=EXCLUDED.listing_type_reason,
			   price_state=EXCLUDED.price_state, price_amount=EXCLUDED.price_amount, price_currency=EXCLUDED.price_currency, price_reason=EXCLUDED.price_reason,
			   available_qty_state=EXCLUDED.available_qty_state, available_qty_value=EXCLUDED.available_qty_value, available_qty_reason=EXCLUDED.available_qty_reason,
			   seller_sku_state=EXCLUDED.seller_sku_state, seller_sku_value=EXCLUDED.seller_sku_value, seller_sku_reason=EXCLUDED.seller_sku_reason,
			   gtin_state=EXCLUDED.gtin_state, gtin_value=EXCLUDED.gtin_value, gtin_reason=EXCLUDED.gtin_reason,
			   last_payload_hash=EXCLUDED.last_payload_hash,
			   updated_at=now()`,
			k.Tenant().String(), k.Account().Channel().String(), k.Account().External(), k.ListingID(), version,
			titleS, titleV, titleR, statusS, statusV, statusR, typeS, typeV, typeR,
			priceS, priceA, priceC, priceR, qtyS, qtyV, qtyR, skuS, skuV, skuR, gtinS, gtinV, gtinR,
			o.Evidence.PayloadHash()); err != nil {
			return fmt.Errorf("listings postgres: upsert listing: %w", err)
		}
		if _, err := tx.Exec(ctx,
			`DELETE FROM listings.listing_variations
			 WHERE tenant_id=$1 AND channel=$2 AND account_external_id=$3 AND listing_id=$4`,
			k.Tenant().String(), k.Account().Channel().String(), k.Account().External(), k.ListingID()); err != nil {
			return fmt.Errorf("listings postgres: clear variations: %w", err)
		}
		for _, v := range o.Variations {
			vPriceS, vPriceA, vPriceC, vPriceR := moneyFactColumns(v.Price)
			vQtyS, vQtyV, vQtyR := intFactColumns(v.AvailableQuantity)
			vSkuS, vSkuV, vSkuR := factColumns(v.SellerSKU)
			vGtinS, vGtinV, vGtinR := factColumns(v.GTIN)
			if _, err := tx.Exec(ctx,
				`INSERT INTO listings.listing_variations (
				   tenant_id, channel, account_external_id, listing_id, variation_id,
				   price_state, price_amount, price_currency, price_reason,
				   available_qty_state, available_qty_value, available_qty_reason,
				   seller_sku_state, seller_sku_value, seller_sku_reason,
				   gtin_state, gtin_value, gtin_reason)
				 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`,
				k.Tenant().String(), k.Account().Channel().String(), k.Account().External(), k.ListingID(), v.VariationID,
				vPriceS, vPriceA, vPriceC, vPriceR, vQtyS, vQtyV, vQtyR, vSkuS, vSkuV, vSkuR, vGtinS, vGtinV, vGtinR); err != nil {
				return fmt.Errorf("listings postgres: insert variation %s: %w", v.VariationID, err)
			}
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO listings.source_observations (
			   tenant_id, channel, account_external_id, listing_id, payload_hash, payload, observed_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)
			 ON CONFLICT (tenant_id, channel, account_external_id, listing_id, payload_hash) DO NOTHING`,
			k.Tenant().String(), k.Account().Channel().String(), k.Account().External(), k.ListingID(),
			o.Evidence.PayloadHash(), o.RawPayload, o.Evidence.ObservedAt()); err != nil {
			return fmt.Errorf("listings postgres: record observation: %w", err)
		}
		return nil
	})
}
```

Pré-verificação obrigatória antes de fechar o arquivo: `fact.Knowledge.String()`
(`kernel/fact/knowledge.go:39`) — confirmar que devolve exatamente
`known/estimated/unknown/not_applicable` (os literais das CHECK constraints). Se divergir,
os literais do SQL seguem o kernel, nunca o contrário.

- [ ] **Step 3: `module.go`** (molde `contexts/catalog/module.go:32-38` — construtor de
  raiz única, §2.2-a):

```go
// Package listings is the context's façade. Everything past this file is
// under internal/; the composition root names a pool and nothing else.
package listings

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/application"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/postgres"
)

type Module struct{ service *application.Service }

func New(pool *pgxpool.Pool) *Module {
	return &Module{service: application.NewService(postgres.NewRepository(pool))}
}

// IngestListing folds one channel observation into Listings.
func (m *Module) IngestListing(ctx context.Context, o contracts.ListingObservation) (contracts.IngestResult, error) {
	return m.service.Ingest(ctx, o)
}
```

- [ ] **Step 4: Verde local no que compila sem Docker + commit**

```bash
go build ./...
go test ./internal/contexts/listings/... -count=1 -v
```

Esperado: build limpo; testes de unidade PASS (a lane hermética prova a integração na
Task 9).

```bash
git add apps/server_core/internal/contexts/listings apps/server_core/tests/integration/listings_ingest_test.go
git commit -m "feat(listings): postgres writer under forced RLS + module facade"
```

---

### Task 6: registro de governança

**Files:**
- Modify: `contracts/governance/modules.json` (entrada nova ao lado da do catalog, `:188-195`)

- [ ] **Step 1: Adicionar a entrada** (kind `context` — unicidade do registro é por
  (kind, id), medida em `scripts/harness/Policy.psm1:361-363`, então coexiste com o módulo
  legado `listings`):

```json
    {
      "id": "listings",
      "kind": "context",
      "root": "apps/server_core/internal/contexts/listings",
      "code_owner_path": "apps/server_core/internal/contexts/listings",
      "composition_required": false,
      "openapi_prefixes": [],
      "dependencies": []
    }
```

- [ ] **Step 2: Lane de governança + commit**

```bash
pwsh -NoProfile -File ./scripts/gate.ps1 -Lane governance
```

Esperado: PASS (sem `GOV_CONTEXT_UNREGISTERED` para `listings` — o detector que caminha a
árvore acusaria o diretório órfão sem esta entrada).

```bash
git add contracts/governance/modules.json
git commit -m "chore(governance): register the listings context"
```

---

### Task 7: `adapters/marketplace/mercadolivre` — o adapter no molde §2

**Files:**
- Create: `apps/server_core/internal/adapters/marketplace/mercadolivre/mercadolivre.go`
- Create: `apps/server_core/internal/adapters/marketplace/mercadolivre/internal/api/client.go`
- Create: `apps/server_core/internal/adapters/marketplace/mercadolivre/internal/api/items.go`
- Create: `apps/server_core/internal/adapters/marketplace/mercadolivre/listings/feed.go`
- Create: `apps/server_core/internal/adapters/marketplace/mercadolivre/listings/feed_test.go`

O DTO de fio vive em `internal/api` e é inalcançável de fora da árvore do vendor (§2.2 —
regra do `internal` do Go na raiz do vendor). A mecânica replica o legado MEDIDO
(nunca importado): scan em `items_scan_ids_reader.go:40-53`, multiget em
`items_multiget_reader.go:243-260`, EAN = atributo `GTIN`/`EAN` `value_name`
(`multiget_mapper.go:197,210`).

- [ ] **Step 1: RED — teste do feed contra servidor HTTP de fixture**

`feed_test.go` (usa `net/http/httptest` com payloads ML reais REDUZIDOS; fixture de >1
página no scan — memória CHIP-MERCADO: truncamento de página 1 é invisível sem fixture
multi-página):

```go
package listings_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre"
	"marketplace-central/apps/server_core/internal/contexts/listings/port"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// fixture: scan devolve 2 ids na página 1 (scroll_id avança) e 1 id na página
// 2; multiget devolve bodies com title/price/status/sku/GTIN e uma variação.
func fixtureServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/users/179571326/items/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok-test" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.URL.Query().Get("search_type") != "scan" {
			t.Errorf("search_type = %q, want scan", r.URL.Query().Get("search_type"))
		}
		switch r.URL.Query().Get("scroll_id") {
		case "":
			_ = json.NewEncoder(w).Encode(map[string]any{"scroll_id": "s2", "results": []string{"MLB1", "MLB2"}})
		case "s2":
			_ = json.NewEncoder(w).Encode(map[string]any{"scroll_id": "s3", "results": []string{"MLB3"}})
		default:
			_ = json.NewEncoder(w).Encode(map[string]any{"scroll_id": "", "results": []string{}})
		}
	})
	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query().Get("ids")
		var out []map[string]any
		for _, id := range splitIDs(ids) {
			body := map[string]any{
				"id": id, "title": "Produto " + id, "status": "active",
				"price": json.Number("199.90"), "currency_id": "BRL",
				"listing_type_id": "gold_special", "available_quantity": 5,
				"seller_sku": "SKU-" + id,
				"attributes": []map[string]any{{"id": "GTIN", "value_name": "7891234567890"}},
				"variations": []map[string]any{{
					"id": 111, "price": json.Number("199.90"), "available_quantity": 3,
					"seller_sku": "VSKU-" + id,
					"attributes": []map[string]any{{"id": "GTIN", "value_name": "7899999999999"}},
				}},
			}
			out = append(out, map[string]any{"code": 200, "body": body})
		}
		_ = json.NewEncoder(w).Encode(out)
	})
	return httptest.NewServer(mux)
}

func splitIDs(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func newFeed(t *testing.T, baseURL string) port.ListingFeed {
	t.Helper()
	bundle, err := mercadolivre.New(mercadolivre.Config{
		BaseURL:   baseURL,
		UserID:    "179571326",
		Channel:   "mercado_livre",
		AccountID: "179571326",
		Token:     func(context.Context) (string, error) { return "tok-test", nil },
	})
	if err != nil {
		t.Fatalf("bundle: %v", err)
	}
	return bundle.ListingFeed
}

func TestFeedWalksEveryScanPageAndMapsFacts(t *testing.T) {
	srv := fixtureServer(t)
	defer srv.Close()
	feed := newFeed(t, srv.URL)
	tid, _ := tenant.Parse("tenant_default")

	page1, err := feed.NextPage(context.Background(), tid, port.Cursor{}, 50)
	if err != nil {
		t.Fatalf("page 1: %v", err)
	}
	if len(page1.Observations) != 2 || page1.Done {
		t.Fatalf("page 1: %d obs done=%v, want 2 obs not done", len(page1.Observations), page1.Done)
	}
	obs := page1.Observations[0]
	if title, ok := obs.Title.Value(); !ok || title != "Produto MLB1" {
		t.Fatalf("title = %v known=%v", title, ok)
	}
	if price, ok := obs.Price.Value(); !ok || price.Amount().StringFixed(2) != "199.90" || price.Currency().String() != "BRL" {
		t.Fatalf("price fact wrong: %v known=%v", price, ok)
	}
	if gtin, ok := obs.GTIN.Value(); !ok || gtin != "7891234567890" {
		t.Fatalf("gtin = %q known=%v", gtin, ok)
	}
	if len(obs.Variations) != 1 || obs.Variations[0].VariationID != "111" {
		t.Fatalf("variations = %+v", obs.Variations)
	}
	if len(obs.RawPayload) == 0 {
		t.Fatal("raw payload not retained")
	}

	page2, err := feed.NextPage(context.Background(), tid, page1.Next, 50)
	if err != nil {
		t.Fatalf("page 2: %v", err)
	}
	if len(page2.Observations) != 1 {
		t.Fatalf("page 2: %d obs, want 1 (page-1 truncation blindness)", len(page2.Observations))
	}

	page3, err := feed.NextPage(context.Background(), tid, page2.Next, 50)
	if err != nil {
		t.Fatalf("page 3: %v", err)
	}
	if !page3.Done {
		t.Fatal("empty scan page must report Done")
	}
}

func TestFeedFailsLoudOnPerItemError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/179571326/items/search", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"scroll_id": "s2", "results": []string{"MLB9"}})
	})
	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]any{{"code": 404, "body": map[string]any{}}})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	feed := newFeed(t, srv.URL)
	tid, _ := tenant.Parse("tenant_default")
	_, err := feed.NextPage(context.Background(), tid, port.Cursor{}, 50)
	if err == nil || !strings.Contains(err.Error(), "MLB9") {
		t.Fatalf("per-item failure must fail the page naming the id, got: %v", err)
	}
}

func TestFeedMapsAbsentFieldsToUnknownNeverZero(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/179571326/items/search", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"scroll_id": "s2", "results": []string{"MLB7"}})
	})
	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		// body só com id: todo o resto ausente.
		_ = json.NewEncoder(w).Encode([]map[string]any{{"code": 200, "body": map[string]any{"id": "MLB7"}}})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	feed := newFeed(t, srv.URL)
	tid, _ := tenant.Parse("tenant_default")
	page, err := feed.NextPage(context.Background(), tid, port.Cursor{}, 50)
	if err != nil {
		t.Fatalf("page: %v", err)
	}
	o := page.Observations[0]
	if _, ok := o.Price.Value(); ok {
		t.Fatal("absent price came back Known — the zero-fabrication this kernel exists to end")
	}
	if o.Price.State().String() != "unknown" {
		t.Fatalf("price state = %s, want unknown", o.Price.State())
	}
}
```

(Imports `strings` onde usado.) Rodar:

```bash
go test ./internal/adapters/marketplace/mercadolivre/... -count=1
```

Esperado: FAIL de compilação (`undefined: mercadolivre.New`).

- [ ] **Step 2: `internal/api/client.go`** — HTTP + auth + decode, o único lugar que conhece
  o fio:

```go
// Package api is Mercado Livre's wire: HTTP, auth header, pagination tokens
// and raw DTOs. Nothing outside adapters/marketplace/mercadolivre can import
// it — that is the Go internal rule at the vendor root (§2.2), and it is the
// boundary the legacy connectors module never had.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TokenSource supplies a live access token. The adapter does not know where
// tokens come from (env, DB credential, future account context) — the
// composition root decides.
type TokenSource func(ctx context.Context) (string, error)

type Client struct {
	base   string
	userID string
	token  TokenSource
	http   *http.Client
}

func NewClient(baseURL, userID string, token TokenSource) (*Client, error) {
	if baseURL == "" || userID == "" || token == nil {
		return nil, fmt.Errorf("mercadolivre api: base url, user id and token source are all required")
	}
	return &Client{base: baseURL, userID: userID, token: token, http: &http.Client{Timeout: 30 * time.Second}}, nil
}

func (c *Client) getJSON(ctx context.Context, path string, out any) error {
	tok, err := c.token(ctx)
	if err != nil {
		return fmt.Errorf("mercadolivre api: token: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+path, nil)
	if err != nil {
		return fmt.Errorf("mercadolivre api: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("mercadolivre api: GET %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("mercadolivre api: read %s: %w", path, err)
	}
	if resp.StatusCode != http.StatusOK {
		// O status e o corpo dizem o quê; o token NUNCA aparece em erro.
		return fmt.Errorf("mercadolivre api: GET %s: status %d: %s", path, resp.StatusCode, truncate(body, 300))
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("mercadolivre api: decode %s: %w", path, err)
	}
	return nil
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "…"
}
```

- [ ] **Step 3: `internal/api/items.go`** — scan + multiget + DTOs de fio (formas medidas
  em `items_scan_ids_reader.go:47-50` e `items_multiget_reader.go:135-181`):

```go
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const multigetBatchSize = 20

// ScanPage enumerates ids via search_type=scan. An empty results slice means
// the walk is over (measured behavior of the legacy reader).
type ScanPage struct {
	ScrollID string   `json:"scroll_id"`
	Results  []string `json:"results"`
}

func (c *Client) ScanIDs(ctx context.Context, scrollID string, limit int) (ScanPage, error) {
	q := url.Values{}
	q.Set("search_type", "scan")
	q.Set("limit", strconv.Itoa(limit))
	if s := strings.TrimSpace(scrollID); s != "" {
		q.Set("scroll_id", s)
	}
	var page ScanPage
	err := c.getJSON(ctx, "/users/"+url.PathEscape(c.userID)+"/items/search?"+q.Encode(), &page)
	return page, err
}

// Item is the multiget body wire shape — the fields this slice consumes plus
// Raw, which retains the channel's full bytes for the observation record.
type Item struct {
	ID                string          `json:"id"`
	Title             string          `json:"title"`
	Status            string          `json:"status"`
	ListingTypeID     string          `json:"listing_type_id"`
	Price             *json.Number    `json:"price"`
	CurrencyID        string          `json:"currency_id"`
	AvailableQuantity *int            `json:"available_quantity"`
	SellerSKU         string          `json:"seller_sku"`
	Attributes        []Attribute     `json:"attributes"`
	Variations        []Variation     `json:"variations"`
	Raw               json.RawMessage `json:"-"`
}

type Attribute struct {
	ID        string `json:"id"`
	ValueName string `json:"value_name"`
}

type Variation struct {
	ID                json.Number  `json:"id"`
	Price             *json.Number `json:"price"`
	AvailableQuantity *int         `json:"available_quantity"`
	SellerSKU         string       `json:"seller_sku"`
	Attributes        []Attribute  `json:"attributes"`
}

// GTIN resolves the EAN attribute: id GTIN with EAN as fallback, value_name —
// the exact rule the legacy mapper measured (multiget_mapper.go:197).
func GTIN(attrs []Attribute) string {
	for _, want := range []string{"GTIN", "EAN"} {
		for _, a := range attrs {
			if a.ID == want && strings.TrimSpace(a.ValueName) != "" {
				return strings.TrimSpace(a.ValueName)
			}
		}
	}
	return ""
}

type multigetElement struct {
	Code int             `json:"code"`
	Body json.RawMessage `json:"body"`
}

// ItemsMultiget hydrates ids in request order, batching at 20 (measured ML
// cap). A per-item code!=200 fails the WHOLE call naming the id: a silently
// skipped listing would read as "absent from the channel", which is a
// different fact.
func (c *Client) ItemsMultiget(ctx context.Context, ids []string) ([]Item, error) {
	items := make([]Item, 0, len(ids))
	for start := 0; start < len(ids); start += multigetBatchSize {
		end := start + multigetBatchSize
		if end > len(ids) {
			end = len(ids)
		}
		batch := ids[start:end]
		var elements []multigetElement
		if err := c.getJSON(ctx, "/items?ids="+url.QueryEscape(strings.Join(batch, ",")), &elements); err != nil {
			return nil, err
		}
		if len(elements) != len(batch) {
			return nil, fmt.Errorf("mercadolivre api: multiget asked %d ids, got %d elements", len(batch), len(elements))
		}
		for i, el := range elements {
			if el.Code != 200 {
				return nil, fmt.Errorf("mercadolivre api: multiget item %s returned code %d", batch[i], el.Code)
			}
			var item Item
			if err := json.Unmarshal(el.Body, &item); err != nil {
				return nil, fmt.Errorf("mercadolivre api: decode item %s: %w", batch[i], err)
			}
			item.Raw = el.Body
			items = append(items, item)
		}
	}
	return items, nil
}
```

- [ ] **Step 4: `listings/feed.go`** — o mapper fio→observação (implementa
  `contexts/listings/port.ListingFeed`):

```go
// Package listings implements contexts/listings/port against Mercado Livre.
// It is the ONLY translator: facts come out with explicit knowledge states,
// and an absent wire field becomes Unknown with a reason, never a zero.
package listings

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre/internal/api"
	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/port"
	"marketplace-central/apps/server_core/internal/kernel/channel"
	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

const systemName = "mercado_livre"

type Feed struct {
	client  *api.Client
	account channel.AccountRef
	now     func() time.Time
}

func NewFeed(client *api.Client, account channel.AccountRef, now func() time.Time) *Feed {
	return &Feed{client: client, account: account, now: now}
}

// NextPage: one scan page of ids, hydrated by multiget, mapped to
// observations. The cursor token IS ML's scroll_id — opaque to the context.
func (f *Feed) NextPage(ctx context.Context, t tenant.ID, after port.Cursor, limit int) (port.Page, error) {
	scan, err := f.client.ScanIDs(ctx, after.Token(), limit)
	if err != nil {
		return port.Page{}, err
	}
	if len(scan.Results) == 0 {
		return port.Page{Done: true}, nil
	}
	items, err := f.client.ItemsMultiget(ctx, scan.Results)
	if err != nil {
		return port.Page{}, err
	}
	observations := make([]contracts.ListingObservation, 0, len(items))
	observedAt := f.now().UTC()
	for _, item := range items {
		obs, err := f.mapItem(t, item, observedAt)
		if err != nil {
			return port.Page{}, fmt.Errorf("mercadolivre listings: map %s: %w", item.ID, err)
		}
		observations = append(observations, obs)
	}
	return port.Page{Observations: observations, Next: port.NewCursor(scan.ScrollID)}, nil
}

func (f *Feed) mapItem(t tenant.ID, item api.Item, observedAt time.Time) (contracts.ListingObservation, error) {
	key, err := contracts.NewSourceListingKey(t, f.account, item.ID)
	if err != nil {
		return contracts.ListingObservation{}, err
	}
	sum := sha256.Sum256(item.Raw)
	evidence, err := provenance.NewEvidence(systemName, "item", item.ID, observedAt, hex.EncodeToString(sum[:]))
	if err != nil {
		return contracts.ListingObservation{}, err
	}

	obs := contracts.ListingObservation{Key: key, Evidence: evidence, RawPayload: item.Raw}
	if obs.Title, err = stringFact(item.Title, "title", evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	if obs.Status, err = stringFact(item.Status, "status", evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	if obs.ListingType, err = stringFact(item.ListingTypeID, "listing_type_id", evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	if obs.Price, err = moneyFact(item.Price, item.CurrencyID, evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	if obs.AvailableQuantity, err = intFact(item.AvailableQuantity, "available_quantity", evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	if obs.SellerSKU, err = stringFact(item.SellerSKU, "seller_sku", evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	if obs.GTIN, err = stringFact(api.GTIN(item.Attributes), "gtin attribute", evidence); err != nil {
		return contracts.ListingObservation{}, err
	}
	for _, v := range item.Variations {
		mapped := contracts.VariationObservation{VariationID: v.ID.String()}
		if mapped.Price, err = moneyFact(v.Price, item.CurrencyID, evidence); err != nil {
			return contracts.ListingObservation{}, err
		}
		if mapped.AvailableQuantity, err = intFact(v.AvailableQuantity, "variation available_quantity", evidence); err != nil {
			return contracts.ListingObservation{}, err
		}
		if mapped.SellerSKU, err = stringFact(v.SellerSKU, "variation seller_sku", evidence); err != nil {
			return contracts.ListingObservation{}, err
		}
		if mapped.GTIN, err = stringFact(api.GTIN(v.Attributes), "variation gtin attribute", evidence); err != nil {
			return contracts.ListingObservation{}, err
		}
		obs.Variations = append(obs.Variations, mapped)
	}
	return obs, nil
}

func stringFact(value, field string, e provenance.Evidence) (fact.Fact[string], error) {
	if value == "" {
		return fact.NewUnknown[string]("ml omitted "+field, e)
	}
	return fact.NewKnown(value, e)
}

func intFact(value *int, field string, e provenance.Evidence) (fact.Fact[int], error) {
	if value == nil {
		return fact.NewUnknown[int]("ml omitted "+field, e)
	}
	return fact.NewKnown(*value, e)
}

func moneyFact(price *json.Number, currencyID string, e provenance.Evidence) (fact.Fact[exact.Money], error) {
	if price == nil || currencyID == "" {
		return fact.NewUnknown[exact.Money]("ml omitted price or currency", e)
	}
	currency, err := exact.ParseCurrency(currencyID)
	if err != nil {
		return fact.Fact[exact.Money]{}, err
	}
	money, err := exact.ParseMoney(price.String(), currency)
	if err != nil {
		return fact.Fact[exact.Money]{}, err
	}
	return fact.NewKnown(money, e)
}
```

(Import `encoding/json` para `*json.Number` no mapper.)

- [ ] **Step 5: `mercadolivre.go`** — a fachada única do vendor:

```go
// Package mercadolivre is the vendor root: New() is the only importable
// surface (§2). The wire lives under internal/api and cannot leave.
package mercadolivre

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre/internal/api"
	mllistings "marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre/listings"
	listingsport "marketplace-central/apps/server_core/internal/contexts/listings/port"
	"marketplace-central/apps/server_core/internal/kernel/channel"
)

// Config carries what the composition root decides: where ML is, which
// account, and how tokens are obtained. Token is a plain func type, NOT
// api.TokenSource: an internal/ type on the façade would not compile for
// callers outside the vendor tree (§2.2-a). Go assigns this func value to the
// named api.TokenSource implicitly — identical underlying type.
type Config struct {
	BaseURL   string
	UserID    string
	Channel   string
	AccountID string
	Token     func(ctx context.Context) (string, error)
}

type Bundle struct {
	ListingFeed listingsport.ListingFeed
}

func New(cfg Config) (Bundle, error) {
	client, err := api.NewClient(cfg.BaseURL, cfg.UserID, cfg.Token)
	if err != nil {
		return Bundle{}, err
	}
	code, err := channel.ParseCode(cfg.Channel)
	if err != nil {
		return Bundle{}, err
	}
	account, err := channel.NewAccountRef(code, cfg.AccountID)
	if err != nil {
		return Bundle{}, err
	}
	return Bundle{ListingFeed: mllistings.NewFeed(client, account, time.Now)}, nil
}
```

Nota: a assinatura de `api.NewClient(base, userID string, token api.TokenSource)` aceita
o func literal do Config sem conversão (tipo subjacente idêntico) — o compilador é o
árbitro; `go build ./...` na Task 8 prova.

- [ ] **Step 6: Verde + commit**

```bash
go build ./...
go test ./internal/adapters/marketplace/mercadolivre/... -count=1 -v
```

Esperado: 3× `--- PASS:` (`TestFeedWalksEveryScanPage…`, `TestFeedFailsLoudOnPerItemError`,
`TestFeedMapsAbsentFieldsToUnknownNeverZero`).

```bash
git add apps/server_core/internal/adapters/marketplace
git commit -m "feat(adapters): mercadolivre vendor in the S2 mold -- wire confined to internal/api"
```

---

### Task 8: composição + `cmd/listingsingest`

**Files:**
- Create: `apps/server_core/internal/composition/listings_wiring.go`
- Create: `apps/server_core/internal/composition/listings_ingest.go`
- Create: `apps/server_core/cmd/listingsingest/main.go`
- Create: `apps/server_core/tests/integration/listings_ingest_composition_test.go`

- [ ] **Step 1: RED — teste de composição** (molde
  `catalog_ingest_composition_test.go:43-94`: caminho real `RunListingsIngest` + Postgres
  da lane + feed de fixture httptest local; asserção final é `count(*)` nas 3 tabelas):

Corpo: sobe o `fixtureServer` da Task 7 (extrair o helper para um arquivo compartilhado do
pacote de teste de integração OU duplicar 40 linhas — duplicar é aceitável em `_test.go`
de pacotes distintos, anotar qual foi feito), monta `mercadolivre.New` com token fixo,
`listings.New(pool)`, roda `composition.RunListingsIngest(ctx, module, feed, tid, 50)` e
assevera `report.Pages==3`, `report.Observed==3`, `report.Created==3`, e depois um segundo
walk com `report.Idempotent==3`; fecha com `count(*)` de `listings.listings` == 3.

```bash
pwsh -NoProfile -File ./scripts/gate.ps1 -Lane integration
```

Esperado agora: FAIL de compilação (`undefined: composition.RunListingsIngest`).

- [ ] **Step 2: `listings_ingest.go`** (molde `catalog_ingest.go:16-67`, mesmo report):

```go
package composition

import (
	"context"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/listings"
	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/port"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// ListingsIngestReport counts every disposition separately: "nothing changed"
// and "nothing was read" are different outcomes (catalog_ingest.go:13-15).
type ListingsIngestReport struct {
	Pages      int
	Observed   int
	Created    int
	Changed    int
	Idempotent int
}

// RunListingsIngest walks a listing feed to exhaustion and folds every
// observation into Listings. Production path: the root owns it, the context
// decides, the adapter speaks wire.
func RunListingsIngest(ctx context.Context, module *listings.Module, feed port.ListingFeed, t tenant.ID, pageSize int) (ListingsIngestReport, error) {
	if pageSize <= 0 {
		return ListingsIngestReport{}, fmt.Errorf("composition: page size must be positive, got %d", pageSize)
	}
	var report ListingsIngestReport
	cursor := port.Cursor{}
	for {
		page, err := feed.NextPage(ctx, t, cursor, pageSize)
		if err != nil {
			return report, fmt.Errorf("composition: read listings feed page %d: %w", report.Pages+1, err)
		}
		report.Pages++
		for _, obs := range page.Observations {
			result, err := module.IngestListing(ctx, obs)
			if err != nil {
				return report, fmt.Errorf("composition: ingest %s: %w", obs.Key.ListingID(), err)
			}
			report.Observed++
			switch result.Disposition {
			case contracts.DispositionCreated:
				report.Created++
			case contracts.DispositionChanged:
				report.Changed++
			case contracts.DispositionIdempotent:
				report.Idempotent++
			}
		}
		if page.Done {
			return report, nil
		}
		if page.Next.IsStart() {
			return report, fmt.Errorf("composition: feed reported more pages but did not advance the cursor")
		}
		cursor = page.Next
	}
}
```

- [ ] **Step 3: `listings_wiring.go`**:

```go
package composition

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre"
	"marketplace-central/apps/server_core/internal/contexts/listings"
)

// ListingsWiring is the assembled listings slice. Like CatalogWiring, this
// file cannot name anything under internal/ — the compiler enforces it.
type ListingsWiring struct {
	Module *listings.Module
	Feed   mercadolivre.Bundle
}

// WireListings assembles the slice. The token source and account identity are
// the root's decision, passed in from the operator entry point.
func WireListings(pool *pgxpool.Pool, mlBaseURL, mlUserID, accountID string, token func(context.Context) (string, error)) (ListingsWiring, error) {
	bundle, err := mercadolivre.New(mercadolivre.Config{
		BaseURL:   mlBaseURL,
		UserID:    mlUserID,
		Channel:   "mercado_livre",
		AccountID: accountID,
		Token:     token,
	})
	if err != nil {
		return ListingsWiring{}, err
	}
	return ListingsWiring{Module: listings.New(pool), Feed: bundle}, nil
}
```

- [ ] **Step 4: `cmd/listingsingest/main.go`** — molde `catalogingest/main.go` (estrutura
  run/exit, fail-closed de tenant `:118-123`) + resolução de credencial do molde
  `mlprobe/main.go:164-190` (SQL verbatim: `integration_installations` ⋈
  `integration_credentials`, `provider_code='mercado_livre'`, `status='connected'`,
  `is_active AND revoked_at IS NULL`, `ORDER BY version DESC LIMIT 1`; decrypt
  `crypto.NewLocalKeyService(key, "local-key-v1").DecryptJSON`; campos `access_token` e
  `user_id` do payload decodificado — verificar no payload real se o user id vem como
  `user_id`; mlprobe deriva o seller de `/users/me`, e se o payload não carregar o id, o
  CLI faz o MESMO GET `/users/me` uma vez e usa o `id` da resposta):

```go
// Command listingsingest drives the Mercado Livre listing feed into the
// listings context. Read-only against ML (GETs only); writes only to the
// listings schema. Token comes from the connected installation's credential
// (same resolution mlprobe uses); the token itself is never printed.
package main

// Estrutura completa espelhando catalogingest/main.go:
//   main() { ctx; if err := run(ctx); err != nil { stderr; exit 1 } }
//   run():
//     1. pgdb.LoadConfig() + requireTenantConfigured (copiar a função e o
//        comentário D-39 de catalogingest/main.go:118-123 — mesmo fail-closed)
//     2. key := os.Getenv("MPC_ENCRYPTION_KEY"); vazio -> erro nomeando a var
//     3. pool := pgdb.NewPool(ctx, dbCfg)
//     4. token, accountID, credTenant := resolveMLCredential(ctx, pool, key)
//        (molde mlprobe/main.go:164-190; accountID = user id do seller,
//        via payload ou GET /users/me)
//     5. tenantID := tenant.Parse(dbCfg.DefaultTenantID); se credTenant !=
//        tenantID.String() -> erro: credencial de outro tenant NUNCA é usada
//     6. wiring := composition.WireListings(pool, "https://api.mercadolibre.com",
//        accountID, accountID, func(context.Context) (string, error) { return token, nil })
//     7. report := composition.RunListingsIngest(ctx, wiring.Module,
//        wiring.Feed.ListingFeed, tenantID, pageSize)  // pageSize: env
//        MPC_LISTINGS_INGEST_PAGE_SIZE, default 50, mesmo loadPageSize de
//        catalogingest/main.go:129-139
//     8. fmt.Printf("listings ingest report: pages=%d observed=%d created=%d changed=%d idempotent=%d\n", ...)
```

> **Superseded na execução, pontos 4/5/7** (commits `b0c24e77`, e a correção de escopo
> de tenant depois da revisão do PR #31). O esqueleto acima copiava o molde do
> `catalogingest`, e o molde era ele próprio a dívida: `MC_DEFAULT_TENANT_ID` e
> `MPC_LISTINGS_INGEST_PAGE_SIZE` viajavam como ambiente para coisas que quem invoca
> decide por corrida. Passaram a flags — `-tenant` (obrigatória, sem default) e
> `-page-size` (default 50) — e `loadPageSize`/`requireTenantConfigured` deixaram de
> existir nos dois comandos. O ponto 4 ganhou também o `tenant_id` como **predicado** da
> consulta de credencial, não como verificação do resultado: escolher a credencial mais
> recente entre todos os tenants e recusá-la depois decifra a de outro antes de perguntar
> de quem é, e deixa todos os tenants menos um sem conseguir ingerir.

Escrever o arquivo completo (não o esqueleto acima — o esqueleto fixa a ordem e as
decisões; o código segue os dois moldes citados linha a linha). O import do pacote
`crypto` de `modules/integrations/adapters/crypto` em `cmd/` é o precedente medido do
mlprobe — `cmd/` está fora de `internal/modules/`, então nenhum predicado GOV o acusa
(Policy.psm1 só varre `internal/modules/`); ainda assim, anotar a dívida na Task 10: a
resolução de credencial pertence ao contexto `account` quando ele nascer.

- [ ] **Step 5: Verde + commit**

```bash
go build ./...
pwsh -NoProfile -File ./scripts/gate.ps1 -Lane integration
```

Esperado: build limpo; lane integration PASS com `tests_run` incluindo os testes novos e
`failure_token` ausente (se Docker local indisponível, marcar e delegar ao runner — a
Task 9 fecha).

```bash
git add apps/server_core/internal/composition apps/server_core/cmd/listingsingest apps/server_core/tests/integration
git commit -m "feat(listings): composition + operator CLI -- the production path"
```

---

### Task 9: gate completo + live drive (operador)

- [ ] **Step 1: Gate**

```bash
npm run gate
```

Esperado: `gate: PASS` em todos os lanes (inclui `boundary` — a árvore nova não está sob
`internal/modules/`, então a contagem legada não muda; e `arch` — os detectores de vendor
token têm `adapters/` isento por regra medida, `internal/arch/scan.go:239-250` — a árvore
nova em `adapters/marketplace/` cai na isenção `/adapters/` por desenho).

- [ ] **Step 2: Push + PR**

```bash
git push -u origin fatia-listings
```

```bash
gh pr create --title "feat(listings): vertical slice leg -- listings context + mercadolivre adapter" --body "Protocolo §15.3, segunda perna da fatia. contexts/listings + adapters/marketplace/mercadolivre (molde §2, DTO confinado em internal/api) + cmd/listingsingest. Só GETs no ML. Evidência do live drive segue em comentário após o merge do operador.

🤖 Generated with [Claude Code](https://claude.com/claude-code)"
```

Checks verdes → merge é do operador.

- [ ] **Step 3: Live drive — DEPOIS do merge, binário novo** (precondição: container
  rebuildado de commit ≥ o do merge — binário velho faz o live drive mentir, memória).
  Operador roda no container:

```bash
docker exec -w /workspace/apps/server_core marketplace-central-backend-1 go run ./cmd/listingsingest -tenant tenant_default
```

`-tenant` é obrigatório e não tem default (commit `b0c24e77`): o tenant é parâmetro
desta corrida, não configuração do container. Trocar `tenant_default` pelo tenant real
da instalação ML se for outro — sem a flag o comando recusa correr, que é o
comportamento pretendido. `-page-size` é opcional, default 50.

Esperado: `listings ingest report: pages=P observed=N created=N changed=0 idempotent=0`
com N > 0. Segunda corrida imediata: `idempotent=N, created=0` — re-poll é grátis.

- [ ] **Step 4: Reconciliação (aceite = observável, dois instrumentos)**

```sql
-- instrumento 1: o contexto novo
SELECT count(*) FROM listings.listings;
-- instrumento 2: o legado, mesmo universo (mesma instalação ML)
SELECT count(DISTINCT provider_listing_id) FROM listings WHERE status='active';
```

As duas contagens conferem para o mesmo universo (anúncios ativos da instalação
conectada), com a diferença explicada linha a linha se houver (anúncio pausado/fechado que
o scan devolve e o legado filtra, etc. — cada divergência nomeada, nunca "aproximadamente
igual"). Colar números + explicação no PR/issue como evidência. Amostra dirigida: para 2
listing_ids, comparar `title_value`, `price_amount`, `seller_sku_value`, `gtin_value`
contra o anúncio real no site do ML.

---

### Task 10: registro final — emendas, dívidas, ratificação parcial

**Files:**
- Modify: `docs/superpowers/specs/2026-08-06-protocolo-de-codigo-design.md` (§17, log de emendas)

- [ ] **Step 1: Linha no §17** com o que a perna confirmou/refutou do molde §2 (fachada de
  vendor, DTO confinado, porta no contexto) — cada afirmação com o `file:line` novo que a
  paga. Se nada foi refutado, a linha diz isso explicitamente ("perna listings: molde §2
  confirmado sem emenda; fachada+internal compilaram na primeira composição" ou o que for
  verdade).

- [ ] **Step 2: Dívidas nomeadas** (comentário no PR e, se o operador mantiver o arquivo,
  em `.mnfs/HARNESS-DEBTS.md` — verificar antes se ainda é o lugar; a doutrina diz que
  `.mnfs/` é arquivo morto, então o default é registrar como ISSUES GitHub):
  - Resolução de credencial ML no `cmd/` (molde mlprobe) → pertence ao contexto `account`.
  - Sem refresh de token: corrida com token expirado morre com 401 alto e claro — rerun
    após refresh pelo caminho legado. Fecha quando `account` nascer.
  - Sem agendamento: CLI manual até `platform/scheduler` existir (§1.2 destino do `sync`).
  - Campos do payload não promovidos a coluna (sold_quantity, permalink, …) — promover
    quando um consumidor nomear a necessidade.
  - Colunas `channel`/`account_external_id` denormalizadas por chave natural — revisitar
    quando `account` der identidade estável a instalações.

- [ ] **Step 3: Commit + push final no branch do PR** (antes do merge) ou comentário no
  PR, conforme o timing do operador.

---

## Fora de escopo (deliberado)

- Contexto `linking` e o desligamento do observer legado listings→product_links — perna
  seguinte da fatia.
- Apagar `internal/modules/listings` — só quando `linking` aterrar (o snapshot observer
  legado depende dele; §15.5: `internal/` fecha atrás de cada contexto que aterra).
- Webhooks/refresh incremental do ML — o legado continua dono do refresh; esta perna é o
  backfill de observação no molde novo.
- Escrita no ML (preço, estoque) — é `changecontrol`, muito adiante na cadeia.
