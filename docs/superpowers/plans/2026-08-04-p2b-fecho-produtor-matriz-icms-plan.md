# Fecho do P2.b — o produtor da matriz de ICMS Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Dar relógio ao par leitor/escritor da matriz de ICMS que o P2.b entregou órfão, para que `icms_matrix_mirror` deixe de ter 0 linhas e os 34 pedidos vinculados passem a mostrar margem em vez de `NULL`.

**Architecture:** O trabalho é de composição, não de cálculo. `internal_read/adapters/oracle.ICMSMatrixReader` (lê `METALPRD.TGFICM`, traduz para `domain.ICMSMatrixCell`) e `internal_read/adapters/mirror.ICMSMatrixWriter` (grava versionado em `icms_matrix_mirror`) já existem, já têm teste de unidade e de integração, e **só têm chamador em `_test.go`**. Esta fatia cria o corpo de job que liga um ao outro (`internal_read/application`, o mesmo lugar onde `market/application/collection_job.go` mora), registra esse job num `syncapp.Scheduler` no `root.go`, e declara as duas arestas de módulo que o P2.b deixou sem declarar. Junto vão duas correções de governança que a auditoria do P2.b nomeou: o `panic` em produção e o prefixo de migração colidido.

**Tech Stack:** Go 1.x (`apps/server_core`), pgx/v5, `database/sql` + godror (Oracle, somente leitura), PostgreSQL, `scripts/harness.ps1` (lanes).

---

## Fora de escopo, declarado

**A-1 (`TgficmRow` no `domain`) não entra nesta fatia.** `internal_read/domain/icms_matrix.go:16` carrega nomes de coluna do Sankhya (`TipRestricao`, `CodTrib`) e `ResolveCell:113-127` decide sobre os literais `"I"`/`"N"`/`"S"`, que são códigos de restrição do fornecedor — regra de negócio escrita no dicionário do ERP. É o desvio de padrão real, e é **refactor de tipo com teste existente em cima**. Fazer junto misturaria "a matriz nunca rodou" com "a matriz mudou de forma" no mesmo diff, e o live drive não saberia qual dos dois quebrou. Fatia própria, depois de os números baterem ao vivo.

**A-5 / D-27 (regex do checador cego a import de raiz de módulo)** já está registrada em `.mnfs/HARNESS-DEBTS.md`. É dívida de harness, não desta fatia.

---

## File Structure

| arquivo | responsabilidade | ação |
|---|---|---|
| `apps/server_core/internal/modules/sync/domain/sync_state.go` | enum de entidade de sync | **modificar** — nova entidade `icms_matrix` (seam compartilhado) |
| `apps/server_core/internal/modules/sync/domain/sync_state_test.go` | prova do enum | **criar** |
| `apps/server_core/internal/modules/internal_read/application/icms_matrix_job.go` | corpo do job: resolve → aplica → cursor | **criar** |
| `apps/server_core/internal/modules/internal_read/application/icms_matrix_job_test.go` | prova do corpo sem banco | **criar** |
| `contracts/governance/modules.json` | arestas de módulo declaradas | **modificar** — `internal_read → sync`, `orders → pricing` |
| `apps/server_core/internal/composition/root.go` | sítio de composição | **modificar** — scheduler + `RegisterJob` + primeira volta |
| `apps/server_core/internal/modules/orders/adapters/pricingtax/reader.go` | adapter de imposto do pedido | **modificar** — `panic` vira `error` |
| `apps/server_core/migrations/0093_icms_matrix.sql` | DDL da matriz | **renomear** para `0094_icms_matrix.sql` + seed idempotente |

**O sítio de composição desta fatia é `apps/server_core/internal/composition/root.go`, logo depois do bloco do scheduler de produtos que termina em `:711`.** Está declarado aqui, no plano, de propósito: o P2.b falhou exatamente por não declarar o seu.

---

## Task 1: Entidade de sync `icms_matrix`

A entidade é o que faz a matriz aparecer sozinha no card **Saúde do sync** da tela `/integracoes`: `sync/adapters/postgres/health_reader.go` faz `SELECT` de toda linha de `sync_state` do tenant sem lista fixa, e `apps/web/src/pages/integracoes/SyncHealthCard.tsx:25-28` gera o rótulo por title-case genérico, sem allowlist. Nenhuma mudança de frontend é necessária — `icms_matrix` renderiza como "Icms Matrix".

Não há migração: o comentário em `sync_state.go:12-14` registra que a validação da entidade é da camada de aplicação, não uma constraint de banco (migração 0075).

**Files:**
- Modify: `apps/server_core/internal/modules/sync/domain/sync_state.go:17-44`
- Create: `apps/server_core/internal/modules/sync/domain/sync_state_test.go`

- [ ] **Step 1: Escrever o teste que falha**

Criar `apps/server_core/internal/modules/sync/domain/sync_state_test.go`:

```go
package domain

import "testing"

func TestEntityValid(t *testing.T) {
	valid := []Entity{
		EntityProducts,
		EntityListings,
		EntityOrders,
		EntityMarket,
		EntityMarketQueue,
		EntityTariffs,
		EntityICMSMatrix,
	}
	for _, e := range valid {
		if !e.Valid() {
			t.Errorf("Entity(%q).Valid() = false, quero true", e)
		}
	}

	invalid := []Entity{"", "icms", "ICMS_MATRIX", "matrix", "produtos"}
	for _, e := range invalid {
		if e.Valid() {
			t.Errorf("Entity(%q).Valid() = true, quero false", e)
		}
	}
}

func TestEntityICMSMatrixWireValue(t *testing.T) {
	// O valor de fio é a chave de sync_state e o rótulo que a tela /integracoes
	// gera por title-case. Mudá-lo abandona a linha já gravada no banco em vez
	// de continuá-la, então o literal é asserção, não detalhe.
	if got := string(EntityICMSMatrix); got != "icms_matrix" {
		t.Fatalf("EntityICMSMatrix = %q, quero \"icms_matrix\"", got)
	}
}
```

- [ ] **Step 2: Rodar e ver falhar**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/sync/domain/ -run TestEntity -v
```

Esperado: FAIL de compilação — `undefined: EntityICMSMatrix`.

- [ ] **Step 3: Adicionar a entidade**

Em `apps/server_core/internal/modules/sync/domain/sync_state.go`, dentro do bloco `const` (depois de `EntityTariffs`):

```go
	EntityTariffs     Entity = "tariffs"
	// EntityICMSMatrix é o espelho de (uf_origem, uf_destino, grupo_icms) lido
	// do TGFICM (icms_matrix_mirror). Entidade própria, e não um estágio do
	// stream de produtos, porque a matriz é fiscal e muda por vigência de
	// legislação — a cadência dela não é a do catálogo, e a falha de uma não
	// pode se esconder no sucesso da outra.
	EntityICMSMatrix Entity = "icms_matrix"
```

E no `switch` de `Valid()`:

```go
	case EntityProducts, EntityListings, EntityOrders, EntityMarket, EntityMarketQueue, EntityTariffs, EntityICMSMatrix:
		return true
```

- [ ] **Step 4: Rodar e ver passar**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/sync/domain/ -run TestEntity -v
```

Esperado: `PASS` — `TestEntityValid` e `TestEntityICMSMatrixWireValue` ok.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/sync/domain/sync_state.go apps/server_core/internal/modules/sync/domain/sync_state_test.go
git commit -m "feat(sync): entidade icms_matrix no enum de sync_state"
```

---

## Task 2: Corpo do job de sync da matriz

O job mora em `internal_read/application` porque `internal_read` é dono das duas pontas (o leitor Oracle e o escritor do espelho). É o mesmo lugar e a mesma forma de `market/application/collection_job.go:38-44`: uma função que devolve `syncapp.JobFunc`, com as dependências declaradas como interfaces **não exportadas definidas pelo consumidor**, para a camada de aplicação nunca importar `adapters/oracle` nem `adapters/mirror`.

Sobre a matriz vazia: `ApplyCells` já devolve `ErrEmptyICMSMatrix` quando `cells` é vazio (`icms_matrix_writer.go:23-28`), porque aplicar vazio fecharia toda célula aberta do tenant. O job **não** trata esse caso à parte — deixa o erro subir, o scheduler grava em `sync_state.last_error`, e a tela `/integracoes` mostra a falha. Silenciar seria dizer "sincronizou" sobre um Oracle inalcançável.

**Files:**
- Create: `apps/server_core/internal/modules/internal_read/application/icms_matrix_job.go`
- Test: `apps/server_core/internal/modules/internal_read/application/icms_matrix_job_test.go`

- [ ] **Step 1: Escrever o teste que falha**

Criar `apps/server_core/internal/modules/internal_read/application/icms_matrix_job_test.go`:

```go
package application_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/application"
	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type fakeResolver struct {
	cells []domain.ICMSMatrixCell
	err   error
	calls int
}

func (f *fakeResolver) ResolveCells(ctx context.Context) ([]domain.ICMSMatrixCell, error) {
	f.calls++
	return f.cells, f.err
}

type fakeApplier struct {
	written  int
	err      error
	gotCells []domain.ICMSMatrixCell
	gotTenant string
	calls    int
}

func (f *fakeApplier) ApplyCells(ctx context.Context, tenantID string, cells []domain.ICMSMatrixCell) (int, error) {
	f.calls++
	f.gotTenant = tenantID
	f.gotCells = cells
	return f.written, f.err
}

func cell(uf string, grupo int, ambiguo bool) domain.ICMSMatrixCell {
	return domain.ICMSMatrixCell{
		UFOrigem:         "MG",
		UFDestino:        uf,
		GrupoICMS:        grupo,
		LinhasCandidatas: 1,
		Ambiguo:          ambiguo,
	}
}

func fixedNow() time.Time {
	return time.Date(2026, 8, 4, 12, 0, 0, 0, time.UTC)
}

func TestICMSMatrixJobPersistsCountsInCursor(t *testing.T) {
	resolver := &fakeResolver{cells: []domain.ICMSMatrixCell{
		cell("BA", 122, false),
		cell("RJ", 311, true),
		cell("SP", 122, false),
	}}
	applier := &fakeApplier{written: 3}

	job := application.NewICMSMatrixJob("tenant-1", resolver, applier, fixedNow)
	next, err := job(context.Background(), nil)
	if err != nil {
		t.Fatalf("job: %v", err)
	}
	if applier.gotTenant != "tenant-1" {
		t.Errorf("tenant repassado = %q, quero \"tenant-1\"", applier.gotTenant)
	}
	if len(applier.gotCells) != 3 {
		t.Errorf("células repassadas = %d, quero 3", len(applier.gotCells))
	}

	var cursor application.ICMSMatrixCursor
	if err := json.Unmarshal(next, &cursor); err != nil {
		t.Fatalf("unmarshal cursor: %v", err)
	}
	if cursor.Cells != 3 {
		t.Errorf("cursor.Cells = %d, quero 3", cursor.Cells)
	}
	if cursor.Ambiguos != 1 {
		t.Errorf("cursor.Ambiguos = %d, quero 1", cursor.Ambiguos)
	}
	if cursor.Written != 3 {
		t.Errorf("cursor.Written = %d, quero 3", cursor.Written)
	}
	if !cursor.CompletedAt.Equal(fixedNow()) {
		t.Errorf("cursor.CompletedAt = %v, quero %v", cursor.CompletedAt, fixedNow())
	}
}

func TestICMSMatrixJobResolveErrorKeepsCursorAndFails(t *testing.T) {
	boom := errors.New("oracle fora do ar")
	resolver := &fakeResolver{err: boom}
	applier := &fakeApplier{}

	job := application.NewICMSMatrixJob("tenant-1", resolver, applier, fixedNow)
	prev := json.RawMessage(`{"cells":9}`)
	next, err := job(context.Background(), prev)
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, quero envolver %v", err, boom)
	}
	if string(next) != string(prev) {
		t.Errorf("cursor = %s, quero o anterior %s intacto", next, prev)
	}
	if applier.calls != 0 {
		t.Errorf("applier chamado %d vez(es); leitura falhou, nada pode ser aplicado", applier.calls)
	}
}

func TestICMSMatrixJobApplyErrorSurfaces(t *testing.T) {
	boom := errors.New("mirror: refusing to apply empty icms matrix")
	resolver := &fakeResolver{cells: nil}
	applier := &fakeApplier{err: boom}

	job := application.NewICMSMatrixJob("tenant-1", resolver, applier, fixedNow)
	prev := json.RawMessage(`{"cells":9}`)
	next, err := job(context.Background(), prev)
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, quero envolver %v", err, boom)
	}
	if string(next) != string(prev) {
		t.Errorf("cursor = %s, quero o anterior %s intacto", next, prev)
	}
}
```

- [ ] **Step 2: Rodar e ver falhar**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/internal_read/application/ -run TestICMSMatrixJob -v
```

Esperado: FAIL de compilação — `undefined: application.NewICMSMatrixJob`.

- [ ] **Step 3: Escrever o job**

Criar `apps/server_core/internal/modules/internal_read/application/icms_matrix_job.go`:

```go
package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
)

// icmsMatrixResolver é a metade de leitura da sincronização da matriz: lê o
// TGFICM da origem fixa e devolve uma célula resolvida por (uf_destino,
// grupo). Interface declarada aqui, no consumidor, para esta camada nunca
// importar adapters/oracle — quem satisfaz hoje é
// internal_read/adapters/oracle.ICMSMatrixReader.
type icmsMatrixResolver interface {
	ResolveCells(ctx context.Context) ([]domain.ICMSMatrixCell, error)
}

// icmsMatrixApplier é a metade de escrita: aplica as células ao espelho
// versionado icms_matrix_mirror e devolve quantas linhas escreveu. Quem
// satisfaz hoje é internal_read/adapters/mirror.ICMSMatrixWriter.
type icmsMatrixApplier interface {
	ApplyCells(ctx context.Context, tenantID string, cells []domain.ICMSMatrixCell) (int, error)
}

// ICMSMatrixCursor é o cursor persistido em sync_state.cursor. Registra o que
// o ciclo fez, não que ele ocorreu: um scheduler que tica com Cells=0 é
// distinguível de um que sincronizou de fato. Ambiguos é contado à parte
// porque uma célula ambígua é gravada sem alíquota (ADR-17) e portanto é uma
// pendência fiscal silenciosa no cálculo do pedido — o operador precisa ver o
// número subir.
type ICMSMatrixCursor struct {
	Cells       int       `json:"cells"`
	Ambiguos    int       `json:"ambiguos"`
	Written     int       `json:"written"`
	CompletedAt time.Time `json:"completed_at"`
}

// NewICMSMatrixJob monta o corpo do ciclo de sincronização da matriz de ICMS.
//
// Existe porque o P2.b entregou o leitor e o escritor sem nenhum relógio: os
// dois só tinham chamador em _test.go, icms_matrix_mirror ficou com 0 linhas, e
// o consumidor inteiro (root.go -> orders/adapters/pricingtax ->
// pricing/adapters/postgres.MatrixReader) leu a tabela vazia e devolveu
// Found:false, deixando margem NULL em 39 de 39 pedidos. O ADR-17 funcionou
// perfeitamente e escondeu que o dado nunca existiu.
//
// Falha fechada em ambas as pontas: erro de leitura não aplica nada (aplicar
// uma lista vazia fecharia toda célula aberta do tenant), e erro de escrita
// sobe para o scheduler, que grava sync_state.last_error e faz a falha
// aparecer na tela /integracoes. Nenhum dos dois inventa um cursor de sucesso.
func NewICMSMatrixJob(
	tenantID string,
	resolver icmsMatrixResolver,
	applier icmsMatrixApplier,
	now func() time.Time,
) syncapp.JobFunc {
	if now == nil {
		now = time.Now
	}
	return func(ctx context.Context, cursor json.RawMessage) (json.RawMessage, error) {
		cells, err := resolver.ResolveCells(ctx)
		if err != nil {
			return cursor, fmt.Errorf("icms matrix sync: resolver células: %w", err)
		}
		written, err := applier.ApplyCells(ctx, tenantID, cells)
		if err != nil {
			return cursor, fmt.Errorf("icms matrix sync: aplicar células: %w", err)
		}
		ambiguos := 0
		for _, c := range cells {
			if c.Ambiguo {
				ambiguos++
			}
		}
		next, err := json.Marshal(ICMSMatrixCursor{
			Cells:       len(cells),
			Ambiguos:    ambiguos,
			Written:     written,
			CompletedAt: now().UTC(),
		})
		if err != nil {
			// A escrita já aconteceu. Reportar o ciclo como falho diria que a
			// matriz não chegou, o que é falso — mantém o cursor anterior.
			return cursor, nil
		}
		return next, nil
	}
}
```

- [ ] **Step 4: Rodar e ver passar**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/internal_read/application/ -run TestICMSMatrixJob -v
```

Esperado: `PASS` nos três testes.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/internal_read/application/icms_matrix_job.go apps/server_core/internal/modules/internal_read/application/icms_matrix_job_test.go
git commit -m "feat(internal_read): corpo do job de sync da matriz de ICMS"
```

---

## Task 3: Declarar as arestas de módulo

Duas arestas de código existem sem estar declaradas em `contracts/governance/modules.json`:

1. **`internal_read → sync`** — nova, criada pela Task 2 (`icms_matrix_job.go` importa `sync/application` pelo tipo `syncapp.JobFunc`). Camada `application` não é `adapters`/`transport`/`registry`, então só a dependência precisa ser declarada, não uma exceção de camada.
2. **`orders → pricing`** — **anterior a esta fatia**, viva desde `9d56b9a7`: `orders/adapters/pricingtax/reader.go:14-15` importa `pricing/domain` e `pricing/ports`. A T5 do P2.b reescreveu o arquivo e não corrigiu a linha de uma palavra em `modules.json`. Fica junto porque é a mesma linha do mesmo arquivo e o gate vai ler as duas.

**Files:**
- Modify: `contracts/governance/modules.json:13` e `:19`

- [ ] **Step 1: Medir a lane ANTES (base)**

```bash
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command governance -BaseSha $(git rev-parse HEAD)
```

Anotar o conjunto de `(error_code, id, path)` da saída. Sem os 40 hex de `-BaseSha` a lane devolve `GOV_SEMANTIC_DRIFT`/`base-sha-invalid` — o veredito só é sólido como diff de conjunto entre base e head, mesmo instrumento nos dois lados.

- [ ] **Step 2: Declarar `internal_read → sync`**

`contracts/governance/modules.json:13` — trocar

```json
"dependencies": ["catalog", "inventory"]
```

por

```json
"dependencies": ["catalog", "inventory", "sync"]
```

na linha do módulo `internal_read`.

- [ ] **Step 3: Declarar `orders → pricing`**

`contracts/governance/modules.json:19` — trocar

```json
"dependencies": ["connectors", "integrations", "internal_read", "product_links", "sync"]
```

por

```json
"dependencies": ["connectors", "integrations", "internal_read", "pricing", "product_links", "sync"]
```

na linha do módulo `orders`.

- [ ] **Step 4: Medir a lane DEPOIS e diferenciar**

```bash
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command governance -BaseSha $(git rev-parse HEAD)
```

Esperado: o conjunto do head é **subconjunto estrito** do da base — some `GOV_MODULE_DEPENDENCY orders-pricing`, e nenhum código novo aparece. Se aparecer qualquer `(error_code, id, path)` que não estava na base, é defeito desta fatia e para aqui.

- [ ] **Step 5: Commit**

```bash
git add contracts/governance/modules.json
git commit -m "chore(governance): declara arestas internal_read->sync e orders->pricing"
```

---

## Task 4: Registrar o job no `root.go` — o sítio de composição

Esta é a task que o P2.b não tinha. Sem ela, tudo que veio antes continua órfão.

O scheduler é próprio (instância separada), como o de mercado em `root.go:738-748`, e não uma segunda entrada no de produtos: cadência fiscal ≠ cadência de catálogo. Escopo `InstallationScopeERP` porque a matriz é fato do ERP, não de uma instalação de marketplace.

`Scheduler.Start` é ticker puro, sem volta de boot (`sync/application/scheduler.go:105-119`, D-16). Com intervalo de 24 h isso deixaria a matriz vazia por 24 h depois de cada deploy — e margem `NULL` em todo pedido nesse intervalo. `RunOnce` é exportado (`:124`), então a primeira volta custa duas linhas.

**Files:**
- Modify: `apps/server_core/internal/composition/root.go` — inserir logo depois do bloco do scheduler de produtos que fecha em `:711`

- [ ] **Step 1: Inserir o bloco de composição**

Depois da chave que fecha o `if activeSourceLookup != nil { ... }` em `root.go:711`, e antes de `marketModuleRepo := ...` em `:713`:

```go
	// P2.b, fecho: a matriz de ICMS tinha leitor (Oracle) e escritor (espelho) e
	// nenhum relógio. icms_matrix_mirror ficou com 0 linhas, o consumidor abaixo
	// (ordersTaxReader, :625) leu tabela vazia, devolveu Found:false, e todo
	// pedido saiu com margem NULL — ADR-17 correto escondendo dado inexistente.
	// Este é o produtor que faltava.
	//
	// Scheduler próprio, e não uma segunda entrada no de produtos: a matriz muda
	// por vigência de legislação, não por movimento de catálogo, e a falha de uma
	// não pode se esconder no sucesso da outra. Escopo ERP porque a matriz é fato
	// do ERP, não de uma instalação de marketplace.
	if oracleDB != nil {
		icmsMatrixScheduler := syncapp.NewScheduler(
			syncpg.NewSyncStateRepository(pool, cfg.DefaultTenantID),
			synccomposition.InstallationScopeERP, 24*time.Hour, time.Now,
		)
		if err := icmsMatrixScheduler.RegisterJob(
			syncdomain.EntityICMSMatrix,
			internalreadapp.NewICMSMatrixJob(
				cfg.DefaultTenantID,
				internalreadoracle.NewICMSMatrixReader(oracleDB),
				mirror.NewICMSMatrixWriter(pool),
				time.Now,
			),
		); err != nil {
			return nil, fmt.Errorf("icms matrix job registration: %w", err)
		}
		// Start é ticker puro (D-16): sem esta primeira volta a matriz ficaria 24h
		// vazia depois de cada deploy e a margem de todo pedido sairia NULL nesse
		// intervalo. RunOnce isola falha por entidade, então uma matriz
		// inalcançável no boot vira sync_state.last_error, não um boot travado.
		go func() {
			icmsMatrixScheduler.RunOnce(context.Background())
			icmsMatrixScheduler.Start(context.Background())
		}()
	}
```

Todos os identificadores usados já estão importados em `root.go`: `syncapp`, `syncpg`, `syncdomain`, `synccomposition`, `internalreadapp` (`:51`), `internalreadoracle`, `mirror` (`:47`), `fmt`, `time`, `context`. Nenhum import novo.

- [ ] **Step 2: Compilar**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go build ./...
```

Esperado: sem saída (sucesso).

- [ ] **Step 3: `go vet` no pacote de composição**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go vet ./internal/composition/...
```

Esperado: sem saída.

- [ ] **Step 4: Provar que o registro existe no binário, não só no diff**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/sync/... ./internal/modules/internal_read/... -count=1
```

Esperado: `ok` em todos os pacotes. A prova de verdade é a contagem no banco, na Task 7 — este passo só garante que nada quebrou.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/composition/root.go
git commit -m "feat(composition): registra o job de sync da matriz de ICMS no scheduler"
```

---

## Task 5: `panic` de produção vira `error`

`orders/adapters/pricingtax/reader.go:238-244` tem um `mustRat` que dá `panic`. Existem dois `mustRat` equivalentes já baselineados em `pricing` (exceção `production-panic-pricing-decompose-mustrat`), mas **uma exceção de baseline não se herda por cópia para arquivo novo** — este é `GOV_PRODUCTION_PANIC` novo, provado novo: `git show 0f4d1e43^:apps/server_core/internal/modules/orders/adapters/pricingtax/reader.go | grep panic` vem vazio.

Os quatro sítios de chamada (`:87`, `:92`, `:97`, `:102`) estão dentro de um laço que **já** devolve `error`, então propagar é mecânico e não muda nenhuma assinatura pública.

**Files:**
- Modify: `apps/server_core/internal/modules/orders/adapters/pricingtax/reader.go:8-16`, `:84-103`, `:235-244`

- [ ] **Step 1: Escrever o teste que falha**

Acrescentar em `apps/server_core/internal/modules/orders/adapters/pricingtax/reader_test.go` (o arquivo já existe e é `package pricingtax` — teste interno, então a função não exportada é alcançável direto):

```go
func TestParseTaxComponentRejectsGarbageWithoutPanic(t *testing.T) {
	// Uma string de imposto que não parseia é erro de programação em
	// pricing/domain, mas em produção ela não pode derrubar o processo do
	// servidor inteiro: o pedido que não calcula é UM pedido, não o serviço.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic em vez de erro: %v", r)
		}
	}()
	if _, err := parseTaxComponent("não é número", "ICMS saída"); err == nil {
		t.Fatal("err = nil, quero erro nomeando o componente")
	}
}
```

- [ ] **Step 2: Rodar e ver falhar**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/orders/adapters/pricingtax/ -run TestParseTaxComponent -v
```

Esperado: FAIL de compilação — `undefined: parseTaxComponent`.

- [ ] **Step 3: Trocar `mustRat` por `parseTaxComponent`**

Em `reader.go`, acrescentar `"fmt"` ao bloco de import:

```go
import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
	pricingdomain "marketplace-central/apps/server_core/internal/modules/pricing/domain"
	pricingports "marketplace-central/apps/server_core/internal/modules/pricing/ports"
)
```

Substituir o bloco `:84-103` inteiro por:

```go
		if tax.ICMSSaida == nil {
			icmsKnown = false
		} else if icmsKnown {
			v, err := parseTaxComponent(*tax.ICMSSaida, "ICMS saída")
			if err != nil {
				return ordersports.OrderTaxes{}, err
			}
			icmsSum.Add(icmsSum, v)
		}
		if tax.Difal == nil {
			difalKnown = false
		} else if difalKnown {
			v, err := parseTaxComponent(*tax.Difal, "DIFAL")
			if err != nil {
				return ordersports.OrderTaxes{}, err
			}
			difalSum.Add(difalSum, v)
		}
		if tax.PisCofins == nil {
			pisKnown = false
		} else if pisKnown {
			v, err := parseTaxComponent(*tax.PisCofins, "PIS/COFINS")
			if err != nil {
				return ordersports.OrderTaxes{}, err
			}
			pisSum.Add(pisSum, v)
		}
		if tax.RestituicaoST == nil {
			restKnown = false
		} else if restKnown {
			v, err := parseTaxComponent(*tax.RestituicaoST, "restituição de ST")
			if err != nil {
				return ordersports.OrderTaxes{}, err
			}
			restSum.Add(restSum, v)
		}
```

Substituir `mustRat` (`:235-244`) por:

```go
// parseTaxComponent lê uma saída de TaxesForItem de volta para big.Rat. Essas
// strings são sempre produto do próprio FormatRatHalfUp, então uma falha aqui é
// erro de programação em pricing/domain, não dado ruim do operador — mas ela
// vira erro, não panic: o que não calcula é UM pedido, e derrubar o processo do
// servidor por isso troca um número faltando por um serviço fora do ar. O nome
// do componente entra na mensagem porque um erro que diz apenas "parse falhou"
// não diz qual dos quatro impostos quebrou.
func parseTaxComponent(s, componente string) (*big.Rat, error) {
	r, err := pricingdomain.ParseRat(s)
	if err != nil {
		return nil, fmt.Errorf("orders pricingtax: %s inválido (%q): %w", componente, s, err)
	}
	return r, nil
}
```

- [ ] **Step 4: Rodar e ver passar**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/orders/... -count=1
```

Esperado: `ok` em todos os pacotes de `orders`, incluindo `TestParseTaxComponentRejectsGarbageWithoutPanic`.

- [ ] **Step 5: Provar que não sobrou `panic` no arquivo**

```bash
grep -n "panic" apps/server_core/internal/modules/orders/adapters/pricingtax/reader.go
```

Esperado: **nenhuma saída**.

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/internal/modules/orders/adapters/pricingtax/reader.go apps/server_core/internal/modules/orders/adapters/pricingtax/reader_test.go
git commit -m "fix(orders): panic de parse de imposto vira erro nomeado por componente"
```

---

## Task 6: Renumerar `0093_icms_matrix.sql` → `0094`

Três fatias paralelas pegaram o prefixo `0093`. As outras duas (`0093_orders_status_details_nullable.sql`, aplicada em 02/08, e `0093_sync_state_market_queue_entity_split.sql`, aplicada em 03/08) colidem entre si desde antes desta fatia — essa colisão está no baseline. A que **esta** linha de trabalho acrescentou é `0093_icms_matrix.sql` (aplicada em 04/08 00:28:38), e é ela que sai.

**Atenção, e é o ponto todo desta task:** `schema_migrations` é chaveado por **nome de arquivo** (`internal/platform/migrate/runner.go:98`). Renomear um arquivo já aplicado faz o runner tratá-lo como não aplicado e rodar de novo. O DDL é todo `IF NOT EXISTS`, mas o `INSERT` de seed das 27 UFs **não é idempotente** — reaplicar viola a PK `(uf, vigente_desde)` e derruba o boot. Por isso o renome vem junto com `ON CONFLICT DO NOTHING`. Não editar `schema_migrations` à mão: a linha velha fica lá como registro histórico honesto de que aquele arquivo rodou.

**Files:**
- Rename: `apps/server_core/migrations/0093_icms_matrix.sql` → `apps/server_core/migrations/0094_icms_matrix.sql`
- Modify: o `INSERT` de seed dentro do arquivo renomeado

- [ ] **Step 1: Confirmar que `0094` está livre**

```bash
git ls-files apps/server_core/migrations | grep -E "^apps/server_core/migrations/009[3-6]"
```

Esperado: `0093_icms_matrix.sql`, `0093_orders_status_details_nullable.sql`, `0093_sync_state_market_queue_entity_split.sql`, `0096_orders_currency.sql`. `0094` e `0095` não aparecem.

- [ ] **Step 2: Renomear preservando histórico**

```bash
git mv apps/server_core/migrations/0093_icms_matrix.sql apps/server_core/migrations/0094_icms_matrix.sql
```

- [ ] **Step 3: Tornar o seed idempotente**

Em `apps/server_core/migrations/0094_icms_matrix.sql`, no fim do `INSERT INTO icms_aliquota_interna ... VALUES (...)`, trocar o `;` final da última linha (`('TO', 20.0, ...)`) por:

```sql
    ('TO', 20.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01')
ON CONFLICT (uf, vigente_desde) DO NOTHING;
```

E acrescentar acima do `INSERT`, junto do comentário de seed existente:

```sql
-- ON CONFLICT DO NOTHING porque schema_migrations é chaveado por NOME DE
-- ARQUIVO (platform/migrate/runner.go): esta migração já rodou como
-- 0093_icms_matrix.sql em ambientes vivos, e o renome para 0094 (colisão de
-- prefixo entre três fatias paralelas) faz o runner reaplicá-la. DO NOTHING
-- mantém o seed já semeado intacto — nunca sobrescreve uma alíquota que uma
-- migração posterior tenha corrigido.
```

- [ ] **Step 4: Provar a idempotência contra o banco de dev**

```bash
docker exec marketplace-central-postgres-1 sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "SELECT count(*) AS ufs FROM icms_aliquota_interna;"'
```

Esperado: `27`. Guardar o número; depois do próximo boot do server ele tem que continuar `27`, não `54` nem erro.

- [ ] **Step 5: Rodar a lane de governança e ver o prefixo sumir**

```bash
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command governance -BaseSha $(git rev-parse HEAD)
```

Esperado: `GOV_MIGRATION_PREFIX` com id envolvendo `icms_matrix` desaparece do conjunto. A colisão remanescente entre `0093_orders_status_details_nullable` e `0093_sync_state_market_queue_entity_split` **permanece** — é de outras fatias, está no baseline, e não é desta fatia consertar.

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/migrations/
git commit -m "chore(migrations): renumera icms_matrix para 0094 e torna o seed idempotente"
```

---

## Task 7: Verificação ao vivo — a Task 7 original do P2.b

**Chip não sobe servidor.** `REQUEST` ao hub para o live drive, com rebuild do binário: um container de pé com código velho faz o live drive mentir — comparar `Created` do container com o horário do primeiro commit desta fatia antes de acreditar em qualquer número.

### O critério de aceite desta fatia é uma contagem no banco, não um teste verde

```sql
SELECT count(*) FILTER (WHERE vigente_ate IS NULL)                  AS celulas_abertas,
       count(*) FILTER (WHERE vigente_ate IS NULL AND ambiguo)      AS ambiguas,
       count(DISTINCT uf_destino) FILTER (WHERE vigente_ate IS NULL) AS ufs_com_celula
FROM icms_matrix_mirror;
```

`celulas_abertas = 0` reprova a fatia na hora, independentemente de qualquer lane verde. Foi exatamente esse `SELECT` que teria reprovado o P2.b em vez de deixá-lo passar sete tasks.

E a linha de sync tem que existir:

```sql
SELECT entity, last_full_sync_at, consecutive_failures, cursor
FROM sync_state
WHERE entity = 'icms_matrix';
```

Esperado: uma linha, `last_full_sync_at` posterior ao boot, `consecutive_failures = 0`, e `cursor` com `cells`/`ambiguos`/`written` não-nulos.

### Alvos numéricos — o dry-run do P2.b já rodou contra os 38 pedidos reais

Os números abaixo saíram de SQL contra o dev stack **antes de existir código** (`.mnfs/MIS-007-ml-sync/M-06-orders-backfill-decomposition/_evidence/p2b-dryrun/calc2.sql`). São alvo de aceite, não expectativa:

```
pedidos | nao_calculavel | positivos_hoje | positivos_novo | viram_negativo | soma_hoje | soma_nova
     38 |              7 |             33 |             28 |              3 |   4560.40 |   1927.01
```

Controle positivo — pedido `2000017515486360` (BA, produto `41912`, `origprod = 0`):

| componente | valor |
|---|---|
| `a_ef` | 23,98% |
| ICMS total | 71,92 |
| PIS/COFINS | 19,00 |
| restituição | +37,69 |
| custo | 154,53 |
| comissão | 40,49 |
| frete | 23,65 |
| **margem nova** | **+28,00** |

Os **7 não-calculáveis** são 5 pedidos sem vínculo de produto e 2 porque o grupo `311` não tem célula de `TGFICM` para RJ. Os dois últimos **têm que aparecer como pendência nomeada** — se saírem com número, a fatia falhou o ADR-17.

**Steps:**
- [ ] Rebuild do binário e boot. Provar o rebuild: `Created` do container posterior ao primeiro commit desta fatia.
- [ ] Rodar os dois `SELECT` acima. Números no pack. `celulas_abertas = 0` = fatia reprovada.
- [ ] Confirmar que `icms_aliquota_interna` continua com **27** linhas depois do boot (prova da idempotência da Task 6).
- [ ] Rodar a agregação dos sete números **contra o código** (a API de pedidos), não contra o SQL do dry-run, e bater os sete. Divergência é defeito ou é achado — nos dois casos escrita no pack antes do `CLOSED`.
- [ ] Drawer do pedido `2000017515486360`: `+28,00`, com as três linhas novas visíveis. Screenshot.
- [ ] Drawer de um pedido **interno** (MG): ICMS de saída `0` com motivo ST **ou** 18%, DIFAL `0` explícito, restituição `0` com motivo. Screenshot.
- [ ] Um dos 2 pedidos de grupo `311`/RJ: pendência nomeada (**D-37**), não alíquota interna.
- [ ] Um item do produto `15956`: aviso de idade do bloco de ST (`fiscal_dt_ref = 2022-02-04`) visível. Se sair silencioso, é **D-44** aberta, não `CLOSED`.
- [ ] Tela `/integracoes`: linha **Icms Matrix** presente no card Saúde do sync, com data de último sucesso. Screenshot.
- [ ] `CLOSED` ao hub com SHA final e os números medidos.

---

## Dívidas que esta fatia NÃO fecha

| id | dívida | por quê fica |
|---|---|---|
| A-1 | `TgficmRow`/`ResolveCell` escrevem a regra no dicionário do Sankhya (`internal_read/domain/icms_matrix.go:16,113-127`) | refactor de tipo; fatia própria depois de os números baterem ao vivo |
| D-27 | checagem de dependência cega a import da raiz de módulo (`Policy.psm1:322`) | dívida de harness, já registrada em `.mnfs/HARNESS-DEBTS.md` |
| D-17 | `CODEMP = 1` fixo no caminho de custo | anterior ao P2.b |
| D-28 | histórico da matriz começa no primeiro sync; a defasagem da BA (17,0% → 20,5% em 20/07/2026) não é recuperável do espelho (está em `TGFHICM`) | escopo do P2.c |
| D-38 | `Imposto` (alíquota de regime) continua no struct e nos sítios do simulador | morte é P4 |
| D-44 | bloco de ST envelhece sem sinal na origem (produto `15956`, `DTMOV = 2022-02-04`) | esta fatia expõe a idade; política de expiração é outra decisão |
| D-45 | `icms_aliquota_interna` é semeada à mão, sem fonte automática | nenhum órgão publica tabela consultável |
| D-46 | `a_inter` derivado de `ORIGPROD ∈ {1,2,3,8}` → 4%, sem FCI | correto para o catálogo medido |
| D-47 | scheduler da matriz de ICMS não tem retry/backoff: se o `RunOnce` de boot falhar (Oracle indisponível no exato momento do deploy), a matriz fica vazia por 24h até o próximo tick — a mesma classe de janela que o `RunOnce` de boot fecha para "todo deploy", mas não para "boot com falha transitória" | achado do review final da fatia (2026-08-04); o intervalo de 24h é decisão do plano (vigência fiscal não muda por hora) e mudar exige decisão do operador, não conserto silencioso |
