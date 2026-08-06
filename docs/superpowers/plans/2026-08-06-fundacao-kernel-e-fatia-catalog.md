# Fundação: kernel, fronteira e a primeira fatia (catalog ← Sankhya)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Construir o kernel de seis membros e a primeira fatia vertical — observação de produto do Sankhya até `contexts/catalog` persistido — com a fronteira imposta pelo compilador e não por revisão.

**Architecture:** Nasce `internal/kernel/` e `internal/contexts/` **ao lado** de `internal/modules/`, sem ponte e sem adapter de transição. Nenhum ficheiro existente em `internal/modules/` é alterado por este plano; nada é apagado por este plano. O código novo é alcançável mas ainda não substitui nenhum caminho de produção — a substituição é dos planos seguintes, quando a fatia ratificar o protocolo. A fronteira entre contextos e a fronteira de vendor são a regra `internal` do Go (nível 1); o que o compilador não exprime é um teste de arquitetura em stdlib puro (nível 2). Zero dependências novas.

**Tech Stack:** Go 1.25.1 (toolchain 1.26.4 instalado), Postgres via `pgx/v5` (já em `go.mod`), Oracle via `godror` (já em `go.mod`), `math/big` da stdlib para aritmética exata, `go/parser` + `go/ast` da stdlib para os checadores.

**Spec:** [docs/superpowers/specs/2026-08-06-protocolo-de-codigo-design.md](../specs/2026-08-06-protocolo-de-codigo-design.md) @ `1162d6e0`

---

## Global Constraints

Estas valem para **todas** as tarefas. Não se repetem em cada uma.

- **Módulo Go:** `marketplace-central/apps/server_core`. Todo o caminho de import começa por aí.
- **Diretório de trabalho de todo o comando Go:** `apps/server_core`. Nunca a raiz do repositório — `GOMODCACHE` na raiz cai numa armadilha de `.gitignore`.
- **Comando de teste, sempre com cache absoluta:**
  ```bash
  cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./... 
  ```
- **Nenhuma dependência nova.** `go.mod` não muda neste plano. Se uma tarefa parecer exigir uma biblioteca, é sinal de erro de desenho: pára e reporta. Alteração de dependência é `REQUEST` ao hub, nunca ato unilateral.
- **Nada em `internal/modules/` é modificado, movido ou apagado por este plano.** Uma única exceção, na Tarefa 6, e é aditiva.
- **Nunca `git push`.** Commits locais só. Push exige permissão explícita do operador.
- **Nunca `git reset`, `revert`, `stash`, `clean`, nem apagar estado desconhecido.**
- **`float64` é proibido** em qualquer tipo criado por este plano. Sem exceção, sem "só neste teste".
- **Toda a chave de linha multi-tenant começa em `tenant_id`.** Sem exceção.
- **RED antes do código.** Todos os RED de uma tarefa correm e falham *antes* de qualquer implementação dessa tarefa. Um RED narrado sem output colado é uma tarefa por fazer.
- **Prova de falha é o output, não a alegação.** Cola-se o texto do erro. "Falhou como esperado" não é evidência.
- **Ficheiro novo entra no commit.** `git status --porcelain --untracked-files=all` vazio antes de fechar cada tarefa — `git diff --exit-code` **ignora ficheiro não rastreado** e passa verde por cima de um ficheiro esquecido.
- **Se uma tarefa bloquear**, escreve o motivo em `.mnfs/HARNESS-DEBTS.md` com `file:line`, faz commit disso, e **passa à tarefa seguinte que não dependa dela**. Não inventes um contorno. Não alteres o plano.

---

## Estrutura de ficheiros

Tudo criado por este plano:

```
apps/server_core/
  internal/
    kernel/
      tenant/tenant.go              T1  ID tipado, sem construtor a partir de string nua
      channel/channel.go            T1  Code + AccountRef
      exact/decimal.go              T2  Decimal sobre big.Rat, exato
      exact/money.go                T2  Money = Decimal + Currency, sem float64
      provenance/evidence.go        T3  Evidence: sistema, tipo, chave, momento, hash
      period/period.go              T4  EffectivePeriod semiaberto [from, until)
      fact/knowledge.go             T5  Knowledge, Fact[T], construtores fail-closed
    contexts/
      catalog/
        contracts/observation.go    T7  ProductObservation, IngestResult, Disposition
        contracts/identifier.go     T7  Identifier, IdentifierKind
        port/reader.go              T7  Reader — o que outros contextos podem perguntar
        internal/domain/id.go       T8  ProductID opaco, SourceProductKey
        internal/domain/product.go  T8  Product, NewProduct, Apply
        internal/application/ingest.go   T9  Service.Ingest
        internal/postgres/repo.go   T10 único escritor de catalog.*
        module.go                   T9  Module: fachada do contexto
    adapters/erp/sankhyaoracle/
      sankhyaoracle.go              T11 fachada: New() Bundle
      internal/oracle/rows.go       T11 ProductRow — forma Sankhya, inalcançável de fora
      catalogfeed/mapper.go         T11 ProductRow → catalog/contracts
  migrations/0097_catalog_context.sql   T10
  internal/arch/
    boundary_test.go                T6  fronteira de contexto e de vendor
    numbers_test.go                 T6  float64 e número nu em contracts
    vendor_test.go                  T6  token de vendor fora de adapters
  tests/integration/
    catalog_ingest_test.go          T12 fatia completa contra Postgres real
scripts/arch-gate.sh                T13
```

---

## Tarefa 1: `kernel/tenant` e `kernel/channel`

Fixa o padrão que as quatro tarefas seguintes repetem: identificador é tipo próprio, construtor valida, e não há literal de struct possível de fora do pacote.

**Files:**
- Create: `apps/server_core/internal/kernel/tenant/tenant.go`
- Create: `apps/server_core/internal/kernel/tenant/tenant_test.go`
- Create: `apps/server_core/internal/kernel/channel/channel.go`
- Create: `apps/server_core/internal/kernel/channel/channel_test.go`

**Interfaces:**
- Produces: `tenant.ID` com `tenant.Parse(string) (ID, error)`, `ID.String() string`, `ID.IsZero() bool`. `channel.Code` com `channel.ParseCode(string) (Code, error)`, `Code.String() string`. `channel.AccountRef` com `channel.NewAccountRef(Code, string) (AccountRef, error)`, `AccountRef.Channel() Code`, `AccountRef.External() string`.
- Consumes: nada.

- [ ] **Step 1: Escreve o teste que falha — `tenant`**

`apps/server_core/internal/kernel/tenant/tenant_test.go`:

```go
package tenant_test

import (
	"testing"

	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func TestParseRejectsEmpty(t *testing.T) {
	if _, err := tenant.Parse(""); err == nil {
		t.Fatal("Parse(\"\") returned no error; empty tenant must be rejected")
	}
}

func TestParseRejectsWhitespaceOnly(t *testing.T) {
	if _, err := tenant.Parse("   "); err == nil {
		t.Fatal("Parse(\"   \") returned no error; blank tenant must be rejected")
	}
}

func TestParseRoundTrips(t *testing.T) {
	id, err := tenant.Parse("tnt_7f3b2")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got := id.String(); got != "tnt_7f3b2" {
		t.Fatalf("String() = %q, want %q", got, "tnt_7f3b2")
	}
	if id.IsZero() {
		t.Fatal("IsZero() = true for a parsed id")
	}
}

func TestZeroValueIsZero(t *testing.T) {
	var id tenant.ID
	if !id.IsZero() {
		t.Fatal("zero ID.IsZero() = false")
	}
	if id.String() != "" {
		t.Fatalf("zero ID.String() = %q, want empty", id.String())
	}
}
```

- [ ] **Step 2: Corre e confirma que falha**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/tenant/... -run Test -v
```

Esperado: falha de compilação, `no required module provides package .../internal/kernel/tenant` ou `undefined: tenant.Parse`. Cola o output no relatório.

- [ ] **Step 3: Implementa `tenant`**

`apps/server_core/internal/kernel/tenant/tenant.go`:

```go
// Package tenant carries the tenant identifier. It is a kernel member because
// every context scopes every row and every query by it, with the same meaning
// and the same invariant everywhere: a tenant is never empty and never guessed.
package tenant

import (
	"errors"
	"strings"
)

// ErrEmpty is returned when a tenant identifier carries no content.
var ErrEmpty = errors.New("tenant: identifier is empty")

// ID is a tenant identifier. The field is unexported so no caller outside this
// package can build one by struct literal and skip validation.
type ID struct {
	value string
}

// Parse builds an ID, rejecting empty and blank input.
func Parse(s string) (ID, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return ID{}, ErrEmpty
	}
	return ID{value: trimmed}, nil
}

// String returns the identifier, or the empty string for the zero value.
func (i ID) String() string { return i.value }

// IsZero reports whether this is the zero value rather than a parsed tenant.
func (i ID) IsZero() bool { return i.value == "" }
```

- [ ] **Step 4: Corre e confirma que passa**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/tenant/... -v
```

Esperado: `ok  marketplace-central/apps/server_core/internal/kernel/tenant`, quatro testes PASS.

- [ ] **Step 5: Escreve o teste que falha — `channel`**

`apps/server_core/internal/kernel/channel/channel_test.go`:

```go
package channel_test

import (
	"testing"

	"marketplace-central/apps/server_core/internal/kernel/channel"
)

func TestParseCodeRejectsEmpty(t *testing.T) {
	if _, err := channel.ParseCode(""); err == nil {
		t.Fatal("ParseCode(\"\") returned no error")
	}
}

func TestParseCodeNormalisesCase(t *testing.T) {
	code, err := channel.ParseCode("MercadoLivre")
	if err != nil {
		t.Fatalf("ParseCode: %v", err)
	}
	if got := code.String(); got != "mercadolivre" {
		t.Fatalf("String() = %q, want %q", got, "mercadolivre")
	}
}

func TestNewAccountRefRejectsZeroCode(t *testing.T) {
	var zero channel.Code
	if _, err := channel.NewAccountRef(zero, "acc-1"); err == nil {
		t.Fatal("NewAccountRef with zero Code returned no error")
	}
}

func TestNewAccountRefRejectsEmptyExternal(t *testing.T) {
	code, err := channel.ParseCode("mercadolivre")
	if err != nil {
		t.Fatalf("ParseCode: %v", err)
	}
	if _, err := channel.NewAccountRef(code, ""); err == nil {
		t.Fatal("NewAccountRef with empty external id returned no error")
	}
}

func TestAccountRefRoundTrips(t *testing.T) {
	code, err := channel.ParseCode("mercadolivre")
	if err != nil {
		t.Fatalf("ParseCode: %v", err)
	}
	ref, err := channel.NewAccountRef(code, "123456789")
	if err != nil {
		t.Fatalf("NewAccountRef: %v", err)
	}
	if ref.Channel().String() != "mercadolivre" || ref.External() != "123456789" {
		t.Fatalf("round trip lost data: %q / %q", ref.Channel().String(), ref.External())
	}
}
```

- [ ] **Step 6: Corre e confirma que falha**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/channel/... -v
```

Esperado: `undefined: channel.ParseCode`. Cola o output.

- [ ] **Step 7: Implementa `channel`**

`apps/server_core/internal/kernel/channel/channel.go`:

```go
// Package channel names a sales channel and an account on it. The code is data,
// never a Go enum: adding a marketplace must not require editing this package.
package channel

import (
	"errors"
	"strings"
)

var (
	// ErrEmptyCode is returned when a channel code carries no content.
	ErrEmptyCode = errors.New("channel: code is empty")
	// ErrEmptyExternal is returned when an account has no external identifier.
	ErrEmptyExternal = errors.New("channel: external account identifier is empty")
)

// Code identifies a sales channel, lower-cased and trimmed.
type Code struct {
	value string
}

// ParseCode builds a Code from free text, rejecting blank input.
func ParseCode(s string) (Code, error) {
	trimmed := strings.ToLower(strings.TrimSpace(s))
	if trimmed == "" {
		return Code{}, ErrEmptyCode
	}
	return Code{value: trimmed}, nil
}

// String returns the code, or the empty string for the zero value.
func (c Code) String() string { return c.value }

// IsZero reports whether this is the zero value.
func (c Code) IsZero() bool { return c.value == "" }

// AccountRef points at one seller account on one channel.
type AccountRef struct {
	channel  Code
	external string
}

// NewAccountRef builds an AccountRef, rejecting a zero code or a blank external id.
func NewAccountRef(c Code, external string) (AccountRef, error) {
	if c.IsZero() {
		return AccountRef{}, ErrEmptyCode
	}
	trimmed := strings.TrimSpace(external)
	if trimmed == "" {
		return AccountRef{}, ErrEmptyExternal
	}
	return AccountRef{channel: c, external: trimmed}, nil
}

// Channel returns the channel this account belongs to.
func (a AccountRef) Channel() Code { return a.channel }

// External returns the channel's own identifier for this account.
func (a AccountRef) External() string { return a.external }
```

- [ ] **Step 8: Corre tudo e confirma que passa**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/... -v && go vet ./internal/kernel/...
```

Esperado: dois pacotes `ok`, nove testes PASS, `go vet` sem saída.

- [ ] **Step 9: Confirma que nada ficou por rastrear e faz commit**

```bash
git status --porcelain --untracked-files=all
git add apps/server_core/internal/kernel/tenant apps/server_core/internal/kernel/channel
git commit -m "feat(kernel): tenant.ID e channel.Code/AccountRef com construtor fail-closed"
```

O `git status` antes do `add` deve listar exatamente os quatro ficheiros novos e nada mais.

---

## Tarefa 2: `kernel/exact` — Decimal e Money sem `float64`

O defeito medido: quatro structs `Money` independentes (`listings/domain/read_model.go:83`, `connectors/domain/money.go:10`, `pricing/domain/decimal.go:12`, `market/domain/market.go:13`), mais um gémeo em `connectors/domain/capability.go:200`, mais `erp_import/domain/import.go:6` que é `string` — e **51 campos monetários em `float64`** num sistema fiscal.

A aritmética é sobre `math/big.Rat`, que é exata para `+ - * /`. Isto não é preciosismo: a §4.4 do protocolo resolve `0.495P = 88.45 + U`, uma divisão que em binário não fecha. O arredondamento existe uma vez só, à saída, e é meio-para-par.

**Files:**
- Create: `apps/server_core/internal/kernel/exact/decimal.go`
- Create: `apps/server_core/internal/kernel/exact/decimal_test.go`
- Create: `apps/server_core/internal/kernel/exact/money.go`
- Create: `apps/server_core/internal/kernel/exact/money_test.go`

**Interfaces:**
- Produces: `exact.Decimal` com `exact.ParseDecimal(string) (Decimal, error)`, `exact.FromInt(int64) Decimal`, e métodos `Add`, `Sub`, `Mul`, `Div(Decimal) (Decimal, error)`, `Cmp(Decimal) int`, `IsZero() bool`, `Neg() Decimal`, `StringFixed(scale int) string`, `Rat() *big.Rat`. `exact.Money` com `exact.NewMoney(Decimal, Currency) (Money, error)`, `exact.ParseMoney(string, Currency) (Money, error)`, `Money.Amount() Decimal`, `Money.Currency() Currency`, `Money.Add(Money) (Money, error)`, `Money.Sub(Money) (Money, error)`, `Money.MulDecimal(Decimal) Money`, `Money.String() string`. `exact.Currency` com `exact.ParseCurrency(string) (Currency, error)`.
- Consumes: nada.

- [ ] **Step 1: Escreve os testes que falham — `Decimal`**

`apps/server_core/internal/kernel/exact/decimal_test.go`:

```go
package exact_test

import (
	"testing"

	"marketplace-central/apps/server_core/internal/kernel/exact"
)

func TestParseDecimalRejectsGarbage(t *testing.T) {
	for _, in := range []string{"", " ", "abc", "1.2.3", "1,50", "NaN", "Inf"} {
		if _, err := exact.ParseDecimal(in); err == nil {
			t.Fatalf("ParseDecimal(%q) returned no error", in)
		}
	}
}

func TestParseDecimalRoundTrips(t *testing.T) {
	d, err := exact.ParseDecimal("82.45")
	if err != nil {
		t.Fatalf("ParseDecimal: %v", err)
	}
	if got := d.StringFixed(2); got != "82.45" {
		t.Fatalf("StringFixed(2) = %q, want %q", got, "82.45")
	}
}

// The whole reason this type exists: 0.1 + 0.2 is exactly 0.3, and a division
// that does not terminate in binary keeps its exact value until it is printed.
func TestArithmeticIsExact(t *testing.T) {
	a, _ := exact.ParseDecimal("0.1")
	b, _ := exact.ParseDecimal("0.2")
	c, _ := exact.ParseDecimal("0.3")
	if a.Add(b).Cmp(c) != 0 {
		t.Fatalf("0.1 + 0.2 != 0.3; got %s", a.Add(b).StringFixed(20))
	}
}

func TestDivIsExactAndReversible(t *testing.T) {
	num, _ := exact.ParseDecimal("88.45")
	den, _ := exact.ParseDecimal("0.495")
	q, err := num.Div(den)
	if err != nil {
		t.Fatalf("Div: %v", err)
	}
	if q.Mul(den).Cmp(num) != 0 {
		t.Fatalf("(88.45/0.495)*0.495 != 88.45; got %s", q.Mul(den).StringFixed(20))
	}
	// The published two-decimal answer from the spec's worked example.
	if got := q.StringFixed(2); got != "178.69" {
		t.Fatalf("StringFixed(2) = %q, want %q", got, "178.69")
	}
}

func TestDivByZeroIsAnError(t *testing.T) {
	num, _ := exact.ParseDecimal("1")
	zero := exact.FromInt(0)
	if _, err := num.Div(zero); err == nil {
		t.Fatal("Div by zero returned no error")
	}
}

func TestStringFixedRoundsHalfToEven(t *testing.T) {
	cases := map[string]string{
		"2.345": "2.34", // half, down to even
		"2.355": "2.36", // half, up to even
		"2.344": "2.34",
		"2.346": "2.35",
		"-2.345": "-2.34",
	}
	for in, want := range cases {
		d, err := exact.ParseDecimal(in)
		if err != nil {
			t.Fatalf("ParseDecimal(%q): %v", in, err)
		}
		if got := d.StringFixed(2); got != want {
			t.Fatalf("StringFixed(2) of %q = %q, want %q", in, got, want)
		}
	}
}

func TestZeroValueIsUsable(t *testing.T) {
	var d exact.Decimal
	if !d.IsZero() {
		t.Fatal("zero Decimal.IsZero() = false")
	}
	if got := d.StringFixed(2); got != "0.00" {
		t.Fatalf("zero StringFixed(2) = %q, want %q", got, "0.00")
	}
}
```

- [ ] **Step 2: Corre e confirma que falha**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/exact/... -v
```

Esperado: `undefined: exact.ParseDecimal`. Cola o output.

- [ ] **Step 3: Implementa `Decimal`**

`apps/server_core/internal/kernel/exact/decimal.go`:

```go
// Package exact carries the platform's numeric types. There is no constructor
// from float64 anywhere in this package, and that is the point: a binary float
// cannot represent 0.1, and a tax base built from one is wrong before any rule
// is applied. Arithmetic is exact rational arithmetic; rounding happens once,
// at presentation, and is half-to-even.
package exact

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var (
	// ErrNotANumber is returned when input is not a plain decimal literal.
	ErrNotANumber = errors.New("exact: not a decimal number")
	// ErrDivideByZero is returned by Div when the divisor is zero.
	ErrDivideByZero = errors.New("exact: divide by zero")
	// ErrNegativeScale is returned by StringFixed for a negative scale.
	ErrNegativeScale = errors.New("exact: negative scale")
)

// Decimal is an exact rational number. The zero value is zero and is usable.
type Decimal struct {
	// r is nil for the zero value. Every method must treat nil as zero.
	r *big.Rat
}

func (d Decimal) rat() *big.Rat {
	if d.r == nil {
		return new(big.Rat)
	}
	return d.r
}

// ParseDecimal accepts a plain decimal literal: optional sign, digits, an
// optional single dot, digits. It deliberately rejects exponent notation,
// thousands separators, comma decimal marks, NaN and Inf, because each of those
// is a sign the value came from somewhere that has not been read carefully.
func ParseDecimal(s string) (Decimal, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return Decimal{}, fmt.Errorf("%w: empty", ErrNotANumber)
	}
	body := t
	if body[0] == '+' || body[0] == '-' {
		body = body[1:]
	}
	if body == "" {
		return Decimal{}, fmt.Errorf("%w: %q", ErrNotANumber, s)
	}
	dots := 0
	digits := 0
	for _, c := range body {
		switch {
		case c == '.':
			dots++
			if dots > 1 {
				return Decimal{}, fmt.Errorf("%w: %q", ErrNotANumber, s)
			}
		case c >= '0' && c <= '9':
			digits++
		default:
			return Decimal{}, fmt.Errorf("%w: %q", ErrNotANumber, s)
		}
	}
	if digits == 0 {
		return Decimal{}, fmt.Errorf("%w: %q", ErrNotANumber, s)
	}
	r, ok := new(big.Rat).SetString(t)
	if !ok {
		return Decimal{}, fmt.Errorf("%w: %q", ErrNotANumber, s)
	}
	return Decimal{r: r}, nil
}

// MustParseDecimal is ParseDecimal for compile-time-known literals in tests and
// for constants. It panics on bad input, which is correct only because the
// input is a literal in the source and never data.
func MustParseDecimal(s string) Decimal {
	d, err := ParseDecimal(s)
	if err != nil {
		panic(err)
	}
	return d
}

// FromInt builds a Decimal from a whole number.
func FromInt(n int64) Decimal {
	return Decimal{r: new(big.Rat).SetInt64(n)}
}

// Add returns d + o.
func (d Decimal) Add(o Decimal) Decimal {
	return Decimal{r: new(big.Rat).Add(d.rat(), o.rat())}
}

// Sub returns d - o.
func (d Decimal) Sub(o Decimal) Decimal {
	return Decimal{r: new(big.Rat).Sub(d.rat(), o.rat())}
}

// Mul returns d * o.
func (d Decimal) Mul(o Decimal) Decimal {
	return Decimal{r: new(big.Rat).Mul(d.rat(), o.rat())}
}

// Div returns d / o, exactly, or ErrDivideByZero.
func (d Decimal) Div(o Decimal) (Decimal, error) {
	if o.IsZero() {
		return Decimal{}, ErrDivideByZero
	}
	return Decimal{r: new(big.Rat).Quo(d.rat(), o.rat())}, nil
}

// Neg returns -d.
func (d Decimal) Neg() Decimal {
	return Decimal{r: new(big.Rat).Neg(d.rat())}
}

// Cmp returns -1, 0 or 1 as d is less than, equal to, or greater than o.
func (d Decimal) Cmp(o Decimal) int { return d.rat().Cmp(o.rat()) }

// IsZero reports whether d is exactly zero.
func (d Decimal) IsZero() bool { return d.rat().Sign() == 0 }

// Sign returns -1, 0 or 1.
func (d Decimal) Sign() int { return d.rat().Sign() }

// Rat returns a copy of the underlying exact value. A copy, so no caller can
// mutate a Decimal that another value shares.
func (d Decimal) Rat() *big.Rat { return new(big.Rat).Set(d.rat()) }

// StringFixed renders d with exactly scale digits after the point, rounding
// half-to-even. Half-to-even and not half-up: half-up is biased away from zero,
// and a bias applied to every line of every order is a systematic error, not a
// rounding difference.
func (d Decimal) StringFixed(scale int) string {
	if scale < 0 {
		panic(ErrNegativeScale)
	}
	r := d.rat()
	neg := r.Sign() < 0
	abs := new(big.Rat).Abs(r)

	pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(scale)), nil)
	scaled := new(big.Rat).Mul(abs, new(big.Rat).SetInt(pow))

	quo, rem := new(big.Int).QuoRem(scaled.Num(), scaled.Denom(), new(big.Int))
	// Compare 2*rem against denom to find the half.
	twice := new(big.Int).Lsh(rem, 1)
	switch twice.Cmp(scaled.Denom()) {
	case 1:
		quo.Add(quo, big.NewInt(1))
	case 0:
		if quo.Bit(0) == 1 { // odd: round up to make it even
			quo.Add(quo, big.NewInt(1))
		}
	}

	digits := quo.String()
	var out string
	if scale == 0 {
		out = digits
	} else {
		for len(digits) <= scale {
			digits = "0" + digits
		}
		out = digits[:len(digits)-scale] + "." + digits[len(digits)-scale:]
	}
	if neg && strings.Trim(out, "0.") != "" {
		out = "-" + out
	}
	return out
}
```

- [ ] **Step 4: Corre e confirma que passa**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/exact/... -run TestParseDecimal -run Test -v
```

Esperado: sete testes PASS. Se `TestStringFixedRoundsHalfToEven` falhar em `-2.345`, o bug está no ramo do sinal e não no arredondamento — corrige-o antes de avançar, não relaxes o teste.

- [ ] **Step 5: Escreve os testes que falham — `Money`**

`apps/server_core/internal/kernel/exact/money_test.go`:

```go
package exact_test

import (
	"testing"

	"marketplace-central/apps/server_core/internal/kernel/exact"
)

func TestParseCurrencyRejectsGarbage(t *testing.T) {
	for _, in := range []string{"", "R$", "brlx", "12"} {
		if _, err := exact.ParseCurrency(in); err == nil {
			t.Fatalf("ParseCurrency(%q) returned no error", in)
		}
	}
}

func TestNewMoneyRejectsZeroCurrency(t *testing.T) {
	var zero exact.Currency
	if _, err := exact.NewMoney(exact.FromInt(10), zero); err == nil {
		t.Fatal("NewMoney with zero Currency returned no error")
	}
}

func TestMoneyStringCarriesCurrency(t *testing.T) {
	brl, err := exact.ParseCurrency("brl")
	if err != nil {
		t.Fatalf("ParseCurrency: %v", err)
	}
	m, err := exact.ParseMoney("82.45", brl)
	if err != nil {
		t.Fatalf("ParseMoney: %v", err)
	}
	if got := m.String(); got != "BRL 82.45" {
		t.Fatalf("String() = %q, want %q", got, "BRL 82.45")
	}
}

// Adding two currencies is not a rounding problem, it is a wrong answer.
func TestMoneyAddRejectsMismatchedCurrency(t *testing.T) {
	brl, _ := exact.ParseCurrency("BRL")
	usd, _ := exact.ParseCurrency("USD")
	a, _ := exact.ParseMoney("10.00", brl)
	b, _ := exact.ParseMoney("10.00", usd)
	if _, err := a.Add(b); err == nil {
		t.Fatal("Add across currencies returned no error")
	}
	if _, err := a.Sub(b); err == nil {
		t.Fatal("Sub across currencies returned no error")
	}
}

func TestMoneyMulDecimalKeepsExactness(t *testing.T) {
	brl, _ := exact.ParseCurrency("BRL")
	price, _ := exact.ParseMoney("178.69", brl)
	rate := exact.MustParseDecimal("0.16")
	fee := price.MulDecimal(rate)
	if got := fee.Amount().StringFixed(4); got != "28.5904" {
		t.Fatalf("StringFixed(4) = %q, want %q", got, "28.5904")
	}
}
```

- [ ] **Step 6: Corre e confirma que falha**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/exact/... -run TestMoney -v
```

Esperado: `undefined: exact.ParseCurrency`. Cola o output.

- [ ] **Step 7: Implementa `Money`**

`apps/server_core/internal/kernel/exact/money.go`:

```go
package exact

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrBadCurrency is returned for anything that is not a 3-letter code.
	ErrBadCurrency = errors.New("exact: currency must be a 3-letter code")
	// ErrCurrencyMismatch is returned when two Money values in different
	// currencies are combined.
	ErrCurrencyMismatch = errors.New("exact: currency mismatch")
)

// Currency is an ISO-4217-shaped code, upper-cased. This package validates the
// shape, not membership of the ISO list: the list changes and is not our truth.
type Currency struct {
	code string
}

// ParseCurrency builds a Currency from a three-letter code.
func ParseCurrency(s string) (Currency, error) {
	t := strings.ToUpper(strings.TrimSpace(s))
	if len(t) != 3 {
		return Currency{}, fmt.Errorf("%w: %q", ErrBadCurrency, s)
	}
	for _, c := range t {
		if c < 'A' || c > 'Z' {
			return Currency{}, fmt.Errorf("%w: %q", ErrBadCurrency, s)
		}
	}
	return Currency{code: t}, nil
}

// String returns the code, or the empty string for the zero value.
func (c Currency) String() string { return c.code }

// IsZero reports whether this is the zero value.
func (c Currency) IsZero() bool { return c.code == "" }

// Money is an amount in a currency. There is no exported field and no
// constructor from float64.
type Money struct {
	amount   Decimal
	currency Currency
}

// NewMoney pairs an amount with a currency, rejecting a zero currency. A Money
// without a currency is not an amount, it is a number that lost its meaning.
func NewMoney(amount Decimal, c Currency) (Money, error) {
	if c.IsZero() {
		return Money{}, fmt.Errorf("%w: currency is the zero value", ErrBadCurrency)
	}
	return Money{amount: amount, currency: c}, nil
}

// ParseMoney parses a decimal literal into Money.
func ParseMoney(s string, c Currency) (Money, error) {
	d, err := ParseDecimal(s)
	if err != nil {
		return Money{}, err
	}
	return NewMoney(d, c)
}

// Amount returns the exact amount.
func (m Money) Amount() Decimal { return m.amount }

// Currency returns the currency.
func (m Money) Currency() Currency { return m.currency }

// Add returns m + o, or ErrCurrencyMismatch.
func (m Money) Add(o Money) (Money, error) {
	if m.currency != o.currency {
		return Money{}, fmt.Errorf("%w: %s + %s", ErrCurrencyMismatch, m.currency, o.currency)
	}
	return Money{amount: m.amount.Add(o.amount), currency: m.currency}, nil
}

// Sub returns m - o, or ErrCurrencyMismatch.
func (m Money) Sub(o Money) (Money, error) {
	if m.currency != o.currency {
		return Money{}, fmt.Errorf("%w: %s - %s", ErrCurrencyMismatch, m.currency, o.currency)
	}
	return Money{amount: m.amount.Sub(o.amount), currency: m.currency}, nil
}

// MulDecimal scales an amount by a dimensionless factor — a commission rate, a
// tax rate, a quantity. The currency is unchanged because a rate has none.
func (m Money) MulDecimal(d Decimal) Money {
	return Money{amount: m.amount.Mul(d), currency: m.currency}
}

// String renders the amount at two decimal places with its currency. It is for
// humans and logs; persistence uses Amount().StringFixed at the column's scale.
func (m Money) String() string {
	if m.currency.IsZero() {
		return m.amount.StringFixed(2)
	}
	return m.currency.String() + " " + m.amount.StringFixed(2)
}
```

- [ ] **Step 8: Corre tudo e confirma que passa**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/... -v && go vet ./internal/kernel/...
```

Esperado: três pacotes `ok`, todos os testes PASS.

- [ ] **Step 9: Prova que nenhum `float64` entrou**

```bash
grep -rn "float64\|float32" apps/server_core/internal/kernel/
```

Esperado: **saída vazia**. Se houver uma linha, é um defeito da tarefa, não uma nota de rodapé.

- [ ] **Step 10: Commit**

```bash
git status --porcelain --untracked-files=all
git add apps/server_core/internal/kernel/exact
git commit -m "feat(kernel): exact.Decimal sobre big.Rat e exact.Money, sem construtor de float64"
```

---

## Tarefa 3: `kernel/provenance` — Evidence

Sem evidência, um facto é uma alegação. A §4.1 recusa qualquer facto sem ela, e é isso que impede o `Unknown` de ser um saco onde se enfia o que não se quis investigar.

**Files:**
- Create: `apps/server_core/internal/kernel/provenance/evidence.go`
- Create: `apps/server_core/internal/kernel/provenance/evidence_test.go`

**Interfaces:**
- Produces: `provenance.Evidence` com `provenance.NewEvidence(system, objectKind, externalKey string, observedAt time.Time, payloadHash string) (Evidence, error)`, e leitores `System()`, `ObjectKind()`, `ExternalKey()`, `ObservedAt()`, `PayloadHash()`, `Ref() string`, `IsZero() bool`.
- Consumes: nada.

- [ ] **Step 1: Escreve o teste que falha**

`apps/server_core/internal/kernel/provenance/evidence_test.go`:

```go
package provenance_test

import (
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

func at() time.Time { return time.Date(2026, 8, 6, 12, 15, 0, 0, time.UTC) }

func TestNewEvidenceRejectsMissingParts(t *testing.T) {
	cases := []struct {
		name       string
		system     string
		kind       string
		key        string
		observed   time.Time
		hash       string
	}{
		{"no system", "", "product", "10529", at(), "sha256:ab91"},
		{"no kind", "sankhya", "", "10529", at(), "sha256:ab91"},
		{"no key", "sankhya", "product", "", at(), "sha256:ab91"},
		{"zero time", "sankhya", "product", "10529", time.Time{}, "sha256:ab91"},
		{"no hash", "sankhya", "product", "10529", at(), ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := provenance.NewEvidence(c.system, c.kind, c.key, c.observed, c.hash); err == nil {
				t.Fatalf("NewEvidence accepted %s", c.name)
			}
		})
	}
}

func TestEvidenceRefIsStable(t *testing.T) {
	e, err := provenance.NewEvidence("sankhya", "product", "10529", at(), "sha256:ab91")
	if err != nil {
		t.Fatalf("NewEvidence: %v", err)
	}
	if got := e.Ref(); got != "sankhya/product:10529" {
		t.Fatalf("Ref() = %q, want %q", got, "sankhya/product:10529")
	}
	if !e.ObservedAt().Equal(at()) {
		t.Fatalf("ObservedAt() = %v, want %v", e.ObservedAt(), at())
	}
}

func TestObservedAtIsStoredInUTC(t *testing.T) {
	saoPaulo := time.FixedZone("-03", -3*60*60)
	local := time.Date(2026, 8, 6, 9, 15, 0, 0, saoPaulo)
	e, err := provenance.NewEvidence("sankhya", "product", "10529", local, "sha256:ab91")
	if err != nil {
		t.Fatalf("NewEvidence: %v", err)
	}
	if e.ObservedAt().Location() != time.UTC {
		t.Fatalf("ObservedAt().Location() = %v, want UTC", e.ObservedAt().Location())
	}
	if !e.ObservedAt().Equal(at()) {
		t.Fatalf("ObservedAt() = %v, want the same instant as %v", e.ObservedAt(), at())
	}
}

func TestZeroEvidenceIsZero(t *testing.T) {
	var e provenance.Evidence
	if !e.IsZero() {
		t.Fatal("zero Evidence.IsZero() = false")
	}
}
```

- [ ] **Step 2: Corre e confirma que falha**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/provenance/... -v
```

Esperado: `undefined: provenance.NewEvidence`. Cola o output.

- [ ] **Step 3: Implementa**

`apps/server_core/internal/kernel/provenance/evidence.go`:

```go
// Package provenance answers "how do we know this?" for every fact the platform
// holds. A fact without evidence is an assertion, and the platform does not
// carry assertions: kernel/fact refuses to build one.
package provenance

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrIncomplete is returned when any part of the evidence is missing.
var ErrIncomplete = errors.New("provenance: incomplete evidence")

// Evidence records where an observation came from and when we saw it.
//
// ObservedAt is when WE collected it. It is deliberately not the same thing as
// when the source says the fact started to hold, nor the legal period the fact
// is valid for — those are separate fields on the facts themselves (§6).
type Evidence struct {
	system      string
	objectKind  string
	externalKey string
	observedAt  time.Time
	payloadHash string
}

// NewEvidence builds Evidence, refusing anything partial.
func NewEvidence(system, objectKind, externalKey string, observedAt time.Time, payloadHash string) (Evidence, error) {
	system = strings.TrimSpace(system)
	objectKind = strings.TrimSpace(objectKind)
	externalKey = strings.TrimSpace(externalKey)
	payloadHash = strings.TrimSpace(payloadHash)

	switch {
	case system == "":
		return Evidence{}, fmt.Errorf("%w: system is empty", ErrIncomplete)
	case objectKind == "":
		return Evidence{}, fmt.Errorf("%w: object kind is empty", ErrIncomplete)
	case externalKey == "":
		return Evidence{}, fmt.Errorf("%w: external key is empty", ErrIncomplete)
	case observedAt.IsZero():
		return Evidence{}, fmt.Errorf("%w: observed_at is the zero time", ErrIncomplete)
	case payloadHash == "":
		return Evidence{}, fmt.Errorf("%w: payload hash is empty", ErrIncomplete)
	}

	return Evidence{
		system:      system,
		objectKind:  objectKind,
		externalKey: externalKey,
		observedAt:  observedAt.UTC(),
		payloadHash: payloadHash,
	}, nil
}

// System returns the source system, e.g. "sankhya" or "mercadolivre".
func (e Evidence) System() string { return e.system }

// ObjectKind returns what kind of thing was observed, e.g. "product".
func (e Evidence) ObjectKind() string { return e.objectKind }

// ExternalKey returns the source system's own key for the object.
func (e Evidence) ExternalKey() string { return e.externalKey }

// ObservedAt returns, in UTC, when we collected the observation.
func (e Evidence) ObservedAt() time.Time { return e.observedAt }

// PayloadHash returns the hash of the raw payload this was read from.
func (e Evidence) PayloadHash() string { return e.payloadHash }

// Ref is the short human-readable pointer, "system/kind:key".
func (e Evidence) Ref() string {
	if e.IsZero() {
		return ""
	}
	return e.system + "/" + e.objectKind + ":" + e.externalKey
}

// IsZero reports whether this is the zero value rather than built evidence.
func (e Evidence) IsZero() bool { return e.system == "" }
```

- [ ] **Step 4: Corre e confirma que passa**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/provenance/... -v
```

Esperado: quatro testes de topo PASS, cinco subtestes PASS.

- [ ] **Step 5: Commit**

```bash
git status --porcelain --untracked-files=all
git add apps/server_core/internal/kernel/provenance
git commit -m "feat(kernel): provenance.Evidence — facto sem proveniência não se constrói"
```

---

## Tarefa 4: `kernel/period` — EffectivePeriod

Custo, alíquota, comissão e classificação fiscal mudam de valor sem mudar de identidade. Retroajustar vigência depois de as colunas serem números nus é reconstrução de semântica, e seis meses depois já ninguém sabe o que era desconhecido.

**Files:**
- Create: `apps/server_core/internal/kernel/period/period.go`
- Create: `apps/server_core/internal/kernel/period/period_test.go`

**Interfaces:**
- Produces: `period.EffectivePeriod` com `period.New(from time.Time, until time.Time) (EffectivePeriod, error)`, `period.From(from time.Time) (EffectivePeriod, error)` para períodos ainda abertos, e `From() time.Time`, `Until() (time.Time, bool)`, `IsOpen() bool`, `Contains(time.Time) bool`, `Overlaps(EffectivePeriod) bool`, `IsZero() bool`.
- Consumes: nada.

- [ ] **Step 1: Escreve o teste que falha**

`apps/server_core/internal/kernel/period/period_test.go`:

```go
package period_test

import (
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/kernel/period"
)

func d(day int) time.Time { return time.Date(2026, 8, day, 0, 0, 0, 0, time.UTC) }

func TestNewRejectsZeroFrom(t *testing.T) {
	if _, err := period.New(time.Time{}, d(10)); err == nil {
		t.Fatal("New with zero from returned no error")
	}
}

func TestNewRejectsUntilBeforeFrom(t *testing.T) {
	if _, err := period.New(d(10), d(5)); err == nil {
		t.Fatal("New with until before from returned no error")
	}
}

func TestNewRejectsEmptyInterval(t *testing.T) {
	if _, err := period.New(d(10), d(10)); err == nil {
		t.Fatal("New with until == from returned no error; a half-open [t,t) holds nothing")
	}
}

// Half-open [from, until). The instant `until` belongs to the NEXT period, so
// two adjacent periods never both claim the same instant.
func TestContainsIsHalfOpen(t *testing.T) {
	p, err := period.New(d(5), d(10))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if !p.Contains(d(5)) {
		t.Fatal("Contains(from) = false; from is inside")
	}
	if p.Contains(d(10)) {
		t.Fatal("Contains(until) = true; until is outside")
	}
	if p.Contains(d(4)) || p.Contains(d(11)) {
		t.Fatal("Contains returned true outside the interval")
	}
}

func TestOpenPeriodHasNoUntil(t *testing.T) {
	p, err := period.From(d(5))
	if err != nil {
		t.Fatalf("From: %v", err)
	}
	if !p.IsOpen() {
		t.Fatal("IsOpen() = false for an open period")
	}
	if _, ok := p.Until(); ok {
		t.Fatal("Until() reported a bound for an open period")
	}
	if !p.Contains(d(9999 % 28)) && !p.Contains(d(28)) {
		t.Fatal("open period does not contain a later instant")
	}
}

func TestOverlapsIsSymmetric(t *testing.T) {
	a, _ := period.New(d(1), d(10))
	b, _ := period.New(d(5), d(15))
	c, _ := period.New(d(10), d(20)) // adjacent to a, not overlapping
	if !a.Overlaps(b) || !b.Overlaps(a) {
		t.Fatal("overlapping periods reported as disjoint")
	}
	if a.Overlaps(c) || c.Overlaps(a) {
		t.Fatal("adjacent half-open periods reported as overlapping")
	}
}

func TestZeroPeriodIsZero(t *testing.T) {
	var p period.EffectivePeriod
	if !p.IsZero() {
		t.Fatal("zero EffectivePeriod.IsZero() = false")
	}
	if p.Contains(d(5)) {
		t.Fatal("zero period contains an instant")
	}
}
```

- [ ] **Step 2: Corre e confirma que falha**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/period/... -v
```

Esperado: `undefined: period.New`. Cola o output.

- [ ] **Step 3: Implementa**

`apps/server_core/internal/kernel/period/period.go`:

```go
// Package period carries validity in time. Costs, tax rates, commissions and
// fiscal classifications change value without changing identity; a number
// without its period is a number whose meaning has to be reconstructed later,
// and by then nobody remembers.
//
// The interval is half-open, [from, until). Two adjacent periods therefore
// never both claim the same instant, which is the property that makes
// "the cost as of T" a question with exactly one answer.
package period

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrNoFrom is returned when the period has no start.
	ErrNoFrom = errors.New("period: from is the zero time")
	// ErrNotOrdered is returned when until is not strictly after from.
	ErrNotOrdered = errors.New("period: until must be strictly after from")
)

// EffectivePeriod is the half-open interval a fact is valid for.
type EffectivePeriod struct {
	from  time.Time
	until time.Time // zero means open-ended
}

// New builds a bounded period.
func New(from, until time.Time) (EffectivePeriod, error) {
	if from.IsZero() {
		return EffectivePeriod{}, ErrNoFrom
	}
	if until.IsZero() {
		return EffectivePeriod{}, fmt.Errorf("%w: use From for an open period", ErrNotOrdered)
	}
	if !until.After(from) {
		return EffectivePeriod{}, fmt.Errorf("%w: from=%s until=%s", ErrNotOrdered, from, until)
	}
	return EffectivePeriod{from: from.UTC(), until: until.UTC()}, nil
}

// From builds an open-ended period that starts at from and has no end yet.
func From(from time.Time) (EffectivePeriod, error) {
	if from.IsZero() {
		return EffectivePeriod{}, ErrNoFrom
	}
	return EffectivePeriod{from: from.UTC()}, nil
}

// Start returns the first instant the fact holds.
func (p EffectivePeriod) Start() time.Time { return p.from }

// Until returns the first instant the fact no longer holds, and whether the
// period is bounded at all.
func (p EffectivePeriod) Until() (time.Time, bool) {
	if p.until.IsZero() {
		return time.Time{}, false
	}
	return p.until, true
}

// IsOpen reports whether the period has no end yet.
func (p EffectivePeriod) IsOpen() bool { return !p.from.IsZero() && p.until.IsZero() }

// Contains reports whether t falls in [from, until).
func (p EffectivePeriod) Contains(t time.Time) bool {
	if p.IsZero() {
		return false
	}
	if t.Before(p.from) {
		return false
	}
	if p.until.IsZero() {
		return true
	}
	return t.Before(p.until)
}

// Overlaps reports whether the two periods share at least one instant.
func (p EffectivePeriod) Overlaps(o EffectivePeriod) bool {
	if p.IsZero() || o.IsZero() {
		return false
	}
	if !p.until.IsZero() && !o.from.Before(p.until) {
		return false
	}
	if !o.until.IsZero() && !p.from.Before(o.until) {
		return false
	}
	return true
}

// IsZero reports whether this is the zero value rather than a built period.
func (p EffectivePeriod) IsZero() bool { return p.from.IsZero() }
```

Nota sobre a assinatura: o leitor do início chama-se `Start()` e não `From()`, porque `From` já é o nome do construtor de período aberto. As tarefas seguintes usam `Start()`.

- [ ] **Step 4: Ajusta o teste à assinatura e corre**

O teste do Step 1 não chama `From()` como leitor, por isso compila sem alteração. Corre:

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/period/... -v
```

Esperado: sete testes PASS.

- [ ] **Step 5: Commit**

```bash
git status --porcelain --untracked-files=all
git add apps/server_core/internal/kernel/period
git commit -m "feat(kernel): period.EffectivePeriod semiaberto — vigência não se reconstrói depois"
```

---

## Tarefa 5: `kernel/fact` — o fim do zero silencioso

Esta é a tarefa que substitui o ADR-017. A regra mais citada da casa — 1.378 citações — deixa de ser doutrina aplicada à mão e passa a construtor cujo mau uso não compila ou devolve erro.

O defeito medido que motiva: `sourcekind/sourcekind.go:47` classifica qualquer sistema desconhecido como `LiveReadThrough` no ramo `default`. O núcleo partilhado viola a regra que o resto do repositório cita 1.378 vezes.

E o número que a justifica: `P = 178,69 + 2,02·U`. Cada real ignorado move o preço dois reais.

**Files:**
- Create: `apps/server_core/internal/kernel/fact/knowledge.go`
- Create: `apps/server_core/internal/kernel/fact/knowledge_test.go`

**Interfaces:**
- Produces: `fact.Knowledge` com as constantes `fact.Known`, `fact.Estimated`, `fact.Unknown`, `fact.NotApplicable`. `fact.Fact[T]` com `fact.NewKnown[T](T, provenance.Evidence) (Fact[T], error)`, `fact.NewEstimated[T](T, reason string, provenance.Evidence) (Fact[T], error)`, `fact.NewUnknown[T](reason string, provenance.Evidence) (Fact[T], error)`, `fact.NewNotApplicable[T](reason string, provenance.Evidence) (Fact[T], error)`, e `Value() (T, bool)`, `State() Knowledge`, `Reason() string`, `Evidence() provenance.Evidence`, `IsUsable() bool`, `MustValue() T`.
- Consumes: `provenance.Evidence` da Tarefa 3.

- [ ] **Step 1: Escreve o teste que falha**

`apps/server_core/internal/kernel/fact/knowledge_test.go`:

```go
package fact_test

import (
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

func ev(t *testing.T) provenance.Evidence {
	t.Helper()
	e, err := provenance.NewEvidence("sankhya", "product", "10529",
		time.Date(2026, 8, 6, 12, 15, 0, 0, time.UTC), "sha256:ab91")
	if err != nil {
		t.Fatalf("NewEvidence: %v", err)
	}
	return e
}

func TestKnownCarriesItsValue(t *testing.T) {
	f, err := fact.NewKnown(exact.FromInt(18), ev(t))
	if err != nil {
		t.Fatalf("NewKnown: %v", err)
	}
	v, ok := f.Value()
	if !ok {
		t.Fatal("Value() ok = false for a Known fact")
	}
	if v.Cmp(exact.FromInt(18)) != 0 {
		t.Fatalf("Value() = %s, want 18", v.StringFixed(0))
	}
	if f.State() != fact.Known {
		t.Fatalf("State() = %v, want Known", f.State())
	}
}

// Known zero is a value. It is NOT the same thing as Unknown, and this test is
// the one that stops the two from ever collapsing into each other again.
func TestKnownZeroIsNotUnknown(t *testing.T) {
	known, err := fact.NewKnown(exact.FromInt(0), ev(t))
	if err != nil {
		t.Fatalf("NewKnown: %v", err)
	}
	v, ok := known.Value()
	if !ok {
		t.Fatal("Value() ok = false for Known zero")
	}
	if !v.IsZero() {
		t.Fatal("Known zero did not carry zero")
	}
	if !known.IsUsable() {
		t.Fatal("Known zero is not usable; a measured zero is a fact")
	}

	unknown, err := fact.NewUnknown[exact.Decimal]("marketplace did not return it", ev(t))
	if err != nil {
		t.Fatalf("NewUnknown: %v", err)
	}
	if _, ok := unknown.Value(); ok {
		t.Fatal("Unknown returned a value")
	}
	if unknown.IsUsable() {
		t.Fatal("Unknown is usable; nothing may compute with it")
	}
	if unknown.State() == known.State() {
		t.Fatal("Unknown and Known collapsed into the same state")
	}
}

func TestUnknownRequiresAReason(t *testing.T) {
	if _, err := fact.NewUnknown[exact.Decimal]("", ev(t)); err == nil {
		t.Fatal("NewUnknown with empty reason returned no error")
	}
	if _, err := fact.NewUnknown[exact.Decimal]("   ", ev(t)); err == nil {
		t.Fatal("NewUnknown with blank reason returned no error")
	}
}

func TestEstimatedRequiresAReason(t *testing.T) {
	if _, err := fact.NewEstimated(exact.FromInt(5), "", ev(t)); err == nil {
		t.Fatal("NewEstimated with empty reason returned no error")
	}
}

func TestEveryStateRequiresEvidence(t *testing.T) {
	var none provenance.Evidence
	if _, err := fact.NewKnown(exact.FromInt(1), none); err == nil {
		t.Fatal("NewKnown without evidence returned no error")
	}
	if _, err := fact.NewUnknown[exact.Decimal]("no data", none); err == nil {
		t.Fatal("NewUnknown without evidence returned no error")
	}
	if _, err := fact.NewEstimated(exact.FromInt(1), "modelled", none); err == nil {
		t.Fatal("NewEstimated without evidence returned no error")
	}
	if _, err := fact.NewNotApplicable[exact.Decimal]("CST 60 has no own ICMS", none); err == nil {
		t.Fatal("NewNotApplicable without evidence returned no error")
	}
}

func TestNotApplicableIsNotUnknown(t *testing.T) {
	na, err := fact.NewNotApplicable[exact.Decimal]("CST 60: no own ICMS on this operation", ev(t))
	if err != nil {
		t.Fatalf("NewNotApplicable: %v", err)
	}
	if na.State() != fact.NotApplicable {
		t.Fatalf("State() = %v, want NotApplicable", na.State())
	}
	if _, ok := na.Value(); ok {
		t.Fatal("NotApplicable returned a value")
	}
	if na.IsUsable() {
		t.Fatal("NotApplicable is usable")
	}
}

func TestZeroFactIsUnknownAndUnusable(t *testing.T) {
	var f fact.Fact[exact.Decimal]
	if f.State() != fact.Unknown {
		t.Fatalf("zero Fact.State() = %v, want Unknown; the zero value must be the safe one", f.State())
	}
	if _, ok := f.Value(); ok {
		t.Fatal("zero Fact returned a value")
	}
	if f.IsUsable() {
		t.Fatal("zero Fact is usable")
	}
}

func TestMustValuePanicsOnUnknown(t *testing.T) {
	f, err := fact.NewUnknown[exact.Decimal]("no data", ev(t))
	if err != nil {
		t.Fatalf("NewUnknown: %v", err)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("MustValue on Unknown did not panic")
		}
	}()
	_ = f.MustValue()
}

func TestEstimatedIsUsableButDistinct(t *testing.T) {
	f, err := fact.NewEstimated(exact.FromInt(7), "modelled from last 30 days", ev(t))
	if err != nil {
		t.Fatalf("NewEstimated: %v", err)
	}
	if !f.IsUsable() {
		t.Fatal("Estimated is not usable; an estimate is a value with a caveat, not an absence")
	}
	if f.State() == fact.Known {
		t.Fatal("Estimated collapsed into Known")
	}
	if f.Reason() == "" {
		t.Fatal("Estimated lost its reason")
	}
}
```

- [ ] **Step 2: Corre e confirma que falha**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/fact/... -v
```

Esperado: `undefined: fact.NewKnown`. Cola o output.

- [ ] **Step 3: Implementa**

`apps/server_core/internal/kernel/fact/knowledge.go`:

```go
// Package fact makes "unknown is never zero" structural instead of doctrinal.
//
// The repository cites that rule 1378 times and still broke it in its own
// shared core, because a rule enforced by remembering is a rule that gets
// forgotten. Here the invalid combinations are unbuildable: Unknown never
// carries a value, Known always does, and neither exists without evidence.
//
// The cost of getting this wrong is not cosmetic. Solving a target-margin
// price with one unknown component U gives P = 178.69 + 2.02*U: every unit of
// ignored cost moves the answer by two. A system that treats unknown as zero
// does not err slightly, it errs by twice what it ignored.
package fact

import (
	"errors"
	"fmt"
	"strings"

	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

// Knowledge is what we know about a quantity. The zero value is Unknown, on
// purpose: a Fact that nobody built must be the safe one, not the confident one.
type Knowledge uint8

const (
	// Unknown means we have no value. It is the zero value.
	Unknown Knowledge = iota
	// Known means we observed the value, including an observed zero.
	Known
	// Estimated means we derived the value and it is usable with a caveat.
	Estimated
	// NotApplicable means the quantity does not exist for this case — which is
	// a different statement from not knowing it.
	NotApplicable
)

// String renders the state for logs and error messages.
func (k Knowledge) String() string {
	switch k {
	case Known:
		return "known"
	case Estimated:
		return "estimated"
	case NotApplicable:
		return "not_applicable"
	default:
		return "unknown"
	}
}

var (
	// ErrReasonRequired is returned when a non-Known state carries no reason.
	ErrReasonRequired = errors.New("fact: state requires a reason code")
	// ErrEvidenceRequired is returned when a fact is built without evidence.
	ErrEvidenceRequired = errors.New("fact: evidence is required")
	// ErrNotUsable is the panic value of MustValue on a fact with no value.
	ErrNotUsable = errors.New("fact: no value in this state")
)

// Fact carries a value together with what we know about it and how we know it.
// Every field is unexported, so the only way to obtain one is a constructor,
// and every constructor validates.
type Fact[T any] struct {
	state    Knowledge
	value    *T
	reason   string
	evidence provenance.Evidence
}

// NewKnown records an observed value. A zero value here is a measured zero and
// is a fact.
func NewKnown[T any](v T, e provenance.Evidence) (Fact[T], error) {
	if e.IsZero() {
		return Fact[T]{}, fmt.Errorf("%w: known value %v", ErrEvidenceRequired, v)
	}
	return Fact[T]{state: Known, value: &v, evidence: e}, nil
}

// NewEstimated records a derived value with the reason it is an estimate.
func NewEstimated[T any](v T, reason string, e provenance.Evidence) (Fact[T], error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return Fact[T]{}, fmt.Errorf("%w: estimated", ErrReasonRequired)
	}
	if e.IsZero() {
		return Fact[T]{}, fmt.Errorf("%w: estimated %v", ErrEvidenceRequired, v)
	}
	return Fact[T]{state: Estimated, value: &v, reason: reason, evidence: e}, nil
}

// NewUnknown records an absence. It takes no value — there is no parameter to
// pass a zero into — and demands the reason we do not have one.
func NewUnknown[T any](reason string, e provenance.Evidence) (Fact[T], error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return Fact[T]{}, fmt.Errorf("%w: unknown", ErrReasonRequired)
	}
	if e.IsZero() {
		return Fact[T]{}, fmt.Errorf("%w: unknown", ErrEvidenceRequired)
	}
	return Fact[T]{state: Unknown, reason: reason, evidence: e}, nil
}

// NewNotApplicable records that the quantity does not exist for this case.
func NewNotApplicable[T any](reason string, e provenance.Evidence) (Fact[T], error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return Fact[T]{}, fmt.Errorf("%w: not applicable", ErrReasonRequired)
	}
	if e.IsZero() {
		return Fact[T]{}, fmt.Errorf("%w: not applicable", ErrEvidenceRequired)
	}
	return Fact[T]{state: NotApplicable, reason: reason, evidence: e}, nil
}

// Value returns the value and whether there is one. The bool is not optional
// to read: `v, _ := f.Value()` on an Unknown fact hands back the zero value of
// T, which is exactly the mistake this package exists to prevent. The vet-level
// check for that pattern lives in internal/arch.
func (f Fact[T]) Value() (T, bool) {
	if f.value == nil {
		var zero T
		return zero, false
	}
	return *f.value, true
}

// MustValue returns the value or panics. It is for call sites that have already
// checked IsUsable in the same function.
func (f Fact[T]) MustValue() T {
	v, ok := f.Value()
	if !ok {
		panic(fmt.Errorf("%w: state=%s reason=%q evidence=%s",
			ErrNotUsable, f.state, f.reason, f.evidence.Ref()))
	}
	return v
}

// State returns what we know.
func (f Fact[T]) State() Knowledge { return f.state }

// Reason returns why, for every state except Known.
func (f Fact[T]) Reason() string { return f.reason }

// Evidence returns how we know.
func (f Fact[T]) Evidence() provenance.Evidence { return f.evidence }

// IsUsable reports whether a calculation may consume this fact. Known and
// Estimated are usable; Unknown and NotApplicable are not.
func (f Fact[T]) IsUsable() bool { return f.value != nil }
```

- [ ] **Step 4: Corre e confirma que passa**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/fact/... -v
```

Esperado: nove testes PASS.

- [ ] **Step 5: Corre o kernel inteiro**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/... && go vet ./internal/kernel/...
```

Esperado: seis pacotes `ok`. `go vet` sem saída.

- [ ] **Step 6: Prova que o kernel tem exatamente seis membros**

```bash
ls apps/server_core/internal/kernel
```

Esperado, exatamente: `channel  exact  fact  period  provenance  tenant`. Um sétimo diretório é uma emenda ao protocolo, não uma decisão de tarefa — se aparecer, é defeito.

- [ ] **Step 7: Commit**

```bash
git status --porcelain --untracked-files=all
git add apps/server_core/internal/kernel/fact
git commit -m "feat(kernel): fact.Fact[T] substitui o ADR-017 — Unknown com valor não se constrói"
```

---

## Tarefa 6: os instrumentos de nível 2, provados contra fixture

As regras que o compilador **não** exprime precisam de um detetor. E um detetor precisa de prova de que dispara — senão é um teste que passa porque não havia nada para ver.

Isto não é hipotético nesta casa. O instrumento que abriu a Onda 1, `internal/composition/module_boundary_arch_test.go:89`, salta todo o ficheiro `_test.go`:

```go
if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
    return nil
}
```

Escondeu **100 imports cruzados em 67 ficheiros de teste**. O número que ele reportou não estava errado por pouco.

Por isso o detetor desta tarefa vive em código normal e não em `_test.go`, e tem dois testes: um contra `testdata/` fabricado, que prova que dispara e conta certo; outro contra a árvore real, que exige zero. Sem o primeiro, o segundo é verde vazio.

**Files:**
- Create: `apps/server_core/internal/arch/scan.go`
- Create: `apps/server_core/internal/arch/scan_test.go`
- Create: `apps/server_core/internal/arch/repo_test.go`
- Create: `apps/server_core/internal/arch/testdata/violations/contexts/alpha/internal/domain/leak.go.txt`
- Create: `apps/server_core/internal/arch/testdata/violations/contexts/alpha/contracts/money.go.txt`
- Create: `apps/server_core/internal/arch/testdata/violations/contexts/alpha/internal/application/vendor.go.txt`
- Create: `apps/server_core/internal/arch/testdata/clean/contexts/beta/contracts/ok.go.txt`
- Modify: `apps/server_core/internal/composition/module_boundary_arch_test.go` (a linha do salto e o comentário que a explica)

**Interfaces:**
- Produces: `arch.Finding` (`File`, `Line`, `Rule`, `Detail`), `arch.Findings`; os detetores `arch.ScanCrossContextInternal(root string) (Findings, error)`, `arch.ScanFloatInContracts(root string) (Findings, error)`, `arch.ScanVendorTokens(root string, tokens []string) (Findings, error)` e as variantes `...Suffix(root, ..., suffix string)` usadas pelas fixtures; a lista fechada `arch.VendorTokens`; as constantes `arch.RuleCrossContextInternal`, `arch.RuleFloatInContracts`, `arch.RuleVendorToken`.
- Consumes: só a stdlib.

- [ ] **Step 1: Escreve as fixtures**

Terminam em `.go.txt` para o toolchain não tentar compilá-las. O detetor faz parse de qualquer ficheiro que passe no sufixo que lhe é dado.

`apps/server_core/internal/arch/testdata/violations/contexts/alpha/internal/domain/leak.go.txt`:

```go
package domain

import (
	_ "marketplace-central/apps/server_core/internal/contexts/beta/internal/domain"
	_ "marketplace-central/apps/server_core/internal/contexts/gamma/contracts"
)
```

A primeira linha de import é violação — alpha entra nos internals de beta. A segunda é legal: `contracts` é superfície publicada. O detetor tem de separar as duas.

`apps/server_core/internal/arch/testdata/violations/contexts/alpha/contracts/money.go.txt`:

```go
package contracts

type Quote struct {
	Commission float64
	Freight    float32
	Reference  string
}
```

`apps/server_core/internal/arch/testdata/violations/contexts/alpha/internal/application/vendor.go.txt`:

```go
package application

const defaultChannel = "mercado_livre"

func channelIsDefault(code string) bool { return code == defaultChannel }
```

Um token de vendor numa literal de string. Nenhum em identificador — os identificadores aqui são deliberadamente neutros, para o teste do Step 2 poder distinguir os dois casos.

`apps/server_core/internal/arch/testdata/clean/contexts/beta/contracts/ok.go.txt`:

```go
package contracts

import (
	"marketplace-central/apps/server_core/internal/kernel/exact"
	_ "marketplace-central/apps/server_core/internal/contexts/gamma/contracts"
)

type Quote struct {
	Commission exact.Money
	Reference  string
}
```

- [ ] **Step 2: Escreve o teste que falha — o detetor contra as fixtures**

`apps/server_core/internal/arch/scan_test.go`:

```go
package arch_test

import (
	"path/filepath"
	"testing"

	"marketplace-central/apps/server_core/internal/arch"
)

const fixtureSuffix = ".go.txt"

func violations() string { return filepath.Join("testdata", "violations") }
func clean() string      { return filepath.Join("testdata", "clean") }

func TestCrossContextInternalFiresOnFixture(t *testing.T) {
	got, err := arch.ScanCrossContextInternalSuffix(violations(), fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("found %d findings, want exactly 1: %v", len(got), got)
	}
	if got[0].Rule != arch.RuleCrossContextInternal {
		t.Fatalf("rule = %q, want %q", got[0].Rule, arch.RuleCrossContextInternal)
	}
	if got[0].Line != 4 {
		t.Fatalf("line = %d, want 4 (the beta/internal/domain import)", got[0].Line)
	}
}

func TestCrossContextInternalIsSilentOnCleanFixture(t *testing.T) {
	got, err := arch.ScanCrossContextInternalSuffix(clean(), fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("found %d findings on the clean fixture, want 0: %v", len(got), got)
	}
}

func TestFloatInContractsFiresOnFixture(t *testing.T) {
	got, err := arch.ScanFloatInContractsSuffix(violations(), fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("found %d findings, want 2 (float64 and float32): %v", len(got), got)
	}
}

func TestFloatInContractsIsSilentOnCleanFixture(t *testing.T) {
	got, err := arch.ScanFloatInContractsSuffix(clean(), fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("found %d findings on the clean fixture, want 0: %v", len(got), got)
	}
}

// The detector must see the token in a string literal, not only in identifiers.
// The measured defect IS a literal: orders/application/ingest_service.go:334
// compares the string "mercado_livre" in the application layer. A detector that
// only reads identifiers walks straight past it.
func TestVendorTokenFiresOnStringLiteral(t *testing.T) {
	got, err := arch.ScanVendorTokensSuffix(violations(), arch.VendorTokens, fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("found %d findings, want exactly 1: %v", len(got), got)
	}
	if got[0].Rule != arch.RuleVendorToken {
		t.Fatalf("rule = %q, want %q", got[0].Rule, arch.RuleVendorToken)
	}
	if !strings.Contains(got[0].Detail, "string literal") {
		t.Fatalf("detail = %q, want it to name the string literal", got[0].Detail)
	}
}

func TestVendorTokenFiresOnIdentifier(t *testing.T) {
	got, err := arch.ScanVendorTokensSuffix(violations(), []string{"channelisdefault"}, fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("found %d findings, want exactly 1 identifier hit: %v", len(got), got)
	}
	if !strings.Contains(got[0].Detail, "identifier") {
		t.Fatalf("detail = %q, want it to name the identifier", got[0].Detail)
	}
}

func TestVendorTokensAreSilentOnCleanFixture(t *testing.T) {
	got, err := arch.ScanVendorTokensSuffix(clean(), arch.VendorTokens, fixtureSuffix)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("found %d findings on the clean fixture, want 0: %v", len(got), got)
	}
}
```

O ficheiro precisa de `"strings"` no bloco de imports do teste. Acrescenta-o.

- [ ] **Step 3: Corre e confirma que falha**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/arch/... -v
```

Esperado: `undefined: arch.ScanCrossContextInternalSuffix`. Cola o output.

- [ ] **Step 4: Implementa o detetor**

`apps/server_core/internal/arch/scan.go`:

```go
// Package arch holds the level-2 instruments: the rules the Go compiler cannot
// express by itself. It is deliberately NOT a _test.go file.
//
// The instrument that opened Wave 1 lived in a test, skipped every _test.go it
// walked, and hid 100 cross-module imports across 67 files. Living in normal
// code means these detectors are themselves testable against fabricated
// violations, which is the only way to know that a green run means "nothing to
// find" and not "nothing was looked at".
package arch

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Rule names, used in findings and in the gate's output.
const (
	RuleCrossContextInternal = "context/internal-import"
	RuleFloatInContracts     = "numbers/float-in-contracts"
	RuleVendorToken          = "adapters/vendor-token-outside-adapters"
)

// VendorTokens is the closed list of marketplace names that may not appear
// outside adapters/. A list and not a regex, because a regex over an import
// path is blind in ways a list is not.
var VendorTokens = []string{
	"mercado_livre", "mercadolivre", "mercadolibre", "meli",
	"shopee", "amazon", "magalu", "americanas",
}

const modulePrefix = "marketplace-central/apps/server_core/"

// Finding is one violation, at a file and a line.
type Finding struct {
	File   string
	Line   int
	Rule   string
	Detail string
}

// Findings is a list of violations, sorted by file then line.
type Findings []Finding

func (f Findings) sortInPlace() {
	sort.Slice(f, func(i, j int) bool {
		if f[i].File != f[j].File {
			return f[i].File < f[j].File
		}
		return f[i].Line < f[j].Line
	})
}

// walk parses every file under root whose name ends in suffix. Nothing is
// skipped by name: a _test.go file is code and lives under the same rules.
func walk(root, suffix string, visit func(path string, fset *token.FileSet, file *ast.File) error) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case "node_modules", ".git", ".gocache":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), suffix) {
			return nil
		}
		src, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		fset := token.NewFileSet()
		parsed, parseErr := parser.ParseFile(fset, path, src, parser.ParseComments)
		if parseErr != nil {
			return parseErr
		}
		return visit(path, fset, parsed)
	})
}

// contextOf returns the context name for a path under .../contexts/NAME/... .
func contextOf(path string) (string, bool) {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for i, p := range parts {
		if p == "contexts" && i+1 < len(parts) {
			return parts[i+1], true
		}
	}
	return "", false
}

// importedContextInternal returns the context whose internals an import path
// reaches into, and whether it reaches into any.
func importedContextInternal(importPath string) (string, bool) {
	rest, ok := strings.CutPrefix(importPath, modulePrefix+"internal/contexts/")
	if !ok {
		return "", false
	}
	parts := strings.Split(rest, "/")
	if len(parts) < 2 {
		return "", false
	}
	for _, p := range parts[1:] {
		if p == "internal" {
			return parts[0], true
		}
	}
	return "", false
}

// ScanCrossContextInternalSuffix reports every file inside one context that
// imports another context's internal packages.
//
// The Go toolchain already refuses this at compile time. This detector exists
// for the window in which a context is being moved, and for the gate's report:
// a number a human can read is worth more than a build error nobody records.
// It must never be the only enforcement.
func ScanCrossContextInternalSuffix(root, suffix string) (Findings, error) {
	var out Findings
	err := walk(root, suffix, func(path string, fset *token.FileSet, file *ast.File) error {
		here, inContext := contextOf(path)
		if !inContext {
			return nil
		}
		for _, imp := range file.Imports {
			value, uErr := strconv.Unquote(imp.Path.Value)
			if uErr != nil {
				continue
			}
			target, reaches := importedContextInternal(value)
			if !reaches || target == here {
				continue
			}
			out = append(out, Finding{
				File:   filepath.ToSlash(path),
				Line:   fset.Position(imp.Pos()).Line,
				Rule:   RuleCrossContextInternal,
				Detail: here + " imports " + value,
			})
		}
		return nil
	})
	out.sortInPlace()
	return out, err
}

// ScanCrossContextInternal scans real Go files.
func ScanCrossContextInternal(root string) (Findings, error) {
	return ScanCrossContextInternalSuffix(root, ".go")
}

// ScanFloatInContractsSuffix reports every float type in a contracts package.
//
// A published contract carrying float64 promises the number survives the round
// trip, and it does not: 0.1 has no binary representation, and a tax base built
// from one is wrong before any rule is applied. 51 such fields were measured in
// the module tree this replaces.
func ScanFloatInContractsSuffix(root, suffix string) (Findings, error) {
	var out Findings
	err := walk(root, suffix, func(path string, fset *token.FileSet, file *ast.File) error {
		if !strings.Contains(filepath.ToSlash(path), "/contracts/") {
			return nil
		}
		ast.Inspect(file, func(n ast.Node) bool {
			ident, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if ident.Name != "float64" && ident.Name != "float32" {
				return true
			}
			out = append(out, Finding{
				File:   filepath.ToSlash(path),
				Line:   fset.Position(ident.Pos()).Line,
				Rule:   RuleFloatInContracts,
				Detail: ident.Name + " in a published contract",
			})
			return true
		})
		return nil
	})
	out.sortInPlace()
	return out, err
}

// ScanFloatInContracts scans real Go files.
func ScanFloatInContracts(root string) (Findings, error) {
	return ScanFloatInContractsSuffix(root, ".go")
}

// ScanVendorTokensSuffix reports every vendor name appearing in an identifier
// or a string literal under root. Both, because the measured defect is a
// literal and an identifier-only detector walks past it.
func ScanVendorTokensSuffix(root string, tokens []string, suffix string) (Findings, error) {
	lower := make([]string, len(tokens))
	for i, t := range tokens {
		lower[i] = strings.ToLower(t)
	}
	var out Findings
	err := walk(root, suffix, func(path string, fset *token.FileSet, file *ast.File) error {
		report := func(pos token.Pos, where, text string) {
			hay := strings.ToLower(text)
			for _, t := range lower {
				if strings.Contains(hay, t) {
					out = append(out, Finding{
						File:   filepath.ToSlash(path),
						Line:   fset.Position(pos).Line,
						Rule:   RuleVendorToken,
						Detail: t + " in " + where,
					})
					return
				}
			}
		}
		ast.Inspect(file, func(n ast.Node) bool {
			switch v := n.(type) {
			case *ast.Ident:
				report(v.Pos(), "identifier", v.Name)
			case *ast.BasicLit:
				if v.Kind == token.STRING {
					report(v.Pos(), "string literal", v.Value)
				}
			}
			return true
		})
		return nil
	})
	out.sortInPlace()
	return out, err
}

// ScanVendorTokens scans real Go files.
func ScanVendorTokens(root string, tokens []string) (Findings, error) {
	return ScanVendorTokensSuffix(root, tokens, ".go")
}
```

- [ ] **Step 5: Corre e confirma que as fixtures passam**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/arch/... -v
```

Esperado: sete testes PASS. Se `TestVendorTokenFiresOnStringLiteral` encontrar mais do que 1, o detetor está a contar o nome do pacote ou o caminho do ficheiro — corrige o detetor, nunca o número esperado.

- [ ] **Step 6: Escreve o teste contra a árvore real**

`apps/server_core/internal/arch/repo_test.go`:

```go
package arch_test

import (
	"os"
	"path/filepath"
	"testing"

	"marketplace-central/apps/server_core/internal/arch"
)

// contextsRoot is internal/contexts. Until Task 7 it does not exist, and a test
// that passes because it looked at nothing must say so out loud.
func contextsRoot(t *testing.T) string {
	t.Helper()
	root := filepath.Join("..", "contexts")
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		t.Skip("internal/contexts does not exist yet; nothing was scanned")
	}
	return root
}

func report(t *testing.T, got arch.Findings, what string) {
	t.Helper()
	for _, f := range got {
		t.Errorf("%s:%d %s: %s", f.File, f.Line, f.Rule, f.Detail)
	}
	if len(got) != 0 {
		t.Fatalf("%d %s", len(got), what)
	}
}

func TestNoContextImportsAnotherContextInternal(t *testing.T) {
	got, err := arch.ScanCrossContextInternal(contextsRoot(t))
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	report(t, got, "cross-context internal imports")
}

func TestNoFloatInAnyContract(t *testing.T) {
	got, err := arch.ScanFloatInContracts(contextsRoot(t))
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	report(t, got, "float fields in published contracts")
}

func TestNoVendorTokenInsideContexts(t *testing.T) {
	got, err := arch.ScanVendorTokens(contextsRoot(t), arch.VendorTokens)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	report(t, got, "vendor tokens inside contexts/")
}

// The kernel exists from Task 1, so this one never skips.
func TestNoVendorTokenInKernel(t *testing.T) {
	got, err := arch.ScanVendorTokens(filepath.Join("..", "kernel"), arch.VendorTokens)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	report(t, got, "vendor tokens in the kernel")
}
```

- [ ] **Step 7: Corre e confirma o estado honesto**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/arch/... -v
```

Esperado: sete PASS de fixture, **três SKIP** com `internal/contexts does not exist yet`, e `TestNoVendorTokenInKernel` PASS. Se os três passarem em vez de saltar, o `os.Stat` está a olhar para o sítio errado e o verde não vale nada.

- [ ] **Step 8: Corrige o ponto cego do instrumento antigo**

Em `apps/server_core/internal/composition/module_boundary_arch_test.go`, a linha:

```go
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
```

passa a:

```go
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
```

E o comentário de bloco no topo ganha, como último parágrafo antes da declaração de `modulesImportPrefix`:

```go
// Corrigido em 2026-08-06: esta caminhada saltava todo o ficheiro _test.go e por
// isso reportava um número que não era o do repositório — 100 imports cruzados
// em 67 ficheiros de teste ficavam de fora. Um teste é código e vive sob as
// mesmas regras. A contagem sobe porque o instrumento passou a ver, não porque
// o código piorou.
```

- [ ] **Step 9: Mede o número verdadeiro**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/composition/... -run TestModuleBoundary -v 2>&1 | tail -40
```

Esperado: continua VERMELHO. É o RED da Onda 1 e não é esta tarefa que o fecha; o que muda é a contagem. **Anota o número novo no relatório.** Se não subir acima de 128, o salto não era o único ponto cego, e isso é um achado por si só — escreve-o em `.mnfs/HARNESS-DEBTS.md`.

- [ ] **Step 10: Commit**

```bash
git status --porcelain --untracked-files=all
git add apps/server_core/internal/arch apps/server_core/internal/composition/module_boundary_arch_test.go
git commit -m "test(arch): detetores de nível 2 provados contra fixture; instrumento antigo deixa de saltar _test.go"
```

---

## Tarefa 7: `contexts/catalog` — a superfície publicada

O contexto nasce pela superfície: o que publica (`contracts`) e o que pergunta (`port`). Nada em `internal/` ainda. Isto obriga a decisão de fronteira a ser tomada **antes** de existir implementação para a justificar a posteriori.

O defeito que esta tarefa fecha é o maior de todos os medidos: `catalog/domain/canonical_product.go:12` declara `type InternalProductID int`, e esse `int` **é** o `CODPROD` do Sankhya, descrito como canónico. Enquanto isso for verdade, um segundo ERP não é aditivo — é colisão de identidade a atravessar todos os contratos.

**Files:**
- Create: `apps/server_core/internal/contexts/catalog/contracts/identifier.go`
- Create: `apps/server_core/internal/contexts/catalog/contracts/observation.go`
- Create: `apps/server_core/internal/contexts/catalog/contracts/contracts_test.go`
- Create: `apps/server_core/internal/contexts/catalog/port/reader.go`

**Interfaces:**
- Produces: `contracts.IdentifierKind` com `contracts.IdentifierEAN`, `IdentifierSKU`, `IdentifierSourceCode`; `contracts.Identifier` com `contracts.NewIdentifier(IdentifierKind, string) (Identifier, error)`, `Kind()`, `Value()`, `String()`; `contracts.SourceProductKey` com `contracts.NewSourceProductKey(tenant.ID, system, instance, objectKind, externalKey string) (SourceProductKey, error)`, `Tenant()`, `System()`, `Instance()`, `ObjectKind()`, `ExternalKey()`, `IsZero()`, `String()`; `contracts.ProductObservation` (campos `Key`, `Description fact.Fact[string]`, `Identifiers []Identifier`, `Evidence provenance.Evidence`) com `Validate() error`; `contracts.Disposition` com `DispositionCreated`, `DispositionChanged`, `DispositionIdempotent`; `contracts.IngestResult` (campos `ProductID string`, `Disposition`, `Version int`, `DuplicateIdentifiers []Identifier`); `port.Summary` e `port.Reader`.
- Consumes: `tenant.ID` (T1), `fact.Fact[string]` (T5), `provenance.Evidence` (T3).

- [ ] **Step 1: Escreve o teste que falha**

`apps/server_core/internal/contexts/catalog/contracts/contracts_test.go`:

```go
package contracts_test

import (
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func tid(t *testing.T) tenant.ID {
	t.Helper()
	id, err := tenant.Parse("tnt_7f3b2")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	return id
}

func ev(t *testing.T) provenance.Evidence {
	t.Helper()
	e, err := provenance.NewEvidence("sankhya", "product", "10529",
		time.Date(2026, 8, 6, 12, 15, 0, 0, time.UTC), "sha256:ab91")
	if err != nil {
		t.Fatalf("provenance.NewEvidence: %v", err)
	}
	return e
}

func TestNewIdentifierRejectsBlank(t *testing.T) {
	if _, err := contracts.NewIdentifier(contracts.IdentifierEAN, "   "); err == nil {
		t.Fatal("NewIdentifier with blank value returned no error")
	}
}

func TestNewIdentifierRejectsUnknownKind(t *testing.T) {
	if _, err := contracts.NewIdentifier(contracts.IdentifierKind("barcode"), "789"); err == nil {
		t.Fatal("NewIdentifier accepted a kind outside the closed set")
	}
}

func TestNewSourceProductKeyRequiresEveryPart(t *testing.T) {
	id := tid(t)
	cases := []struct{ name, system, instance, kind, key string }{
		{"no system", "", "sankhya-prod-01", "product", "10529"},
		{"no instance", "sankhya", "", "product", "10529"},
		{"no kind", "sankhya", "sankhya-prod-01", "", "10529"},
		{"no key", "sankhya", "sankhya-prod-01", "product", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := contracts.NewSourceProductKey(id, c.system, c.instance, c.kind, c.key); err == nil {
				t.Fatalf("NewSourceProductKey accepted %s", c.name)
			}
		})
	}
	var noTenant tenant.ID
	if _, err := contracts.NewSourceProductKey(noTenant, "sankhya", "sankhya-prod-01", "product", "10529"); err == nil {
		t.Fatal("NewSourceProductKey accepted a zero tenant")
	}
}

// The ERP code lives inside the source key and nowhere else. This is the test
// that stops CODPROD from becoming the canonical identifier a second time.
func TestSourceProductKeyCarriesTheErpCodeAndNothingElseDoes(t *testing.T) {
	k, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	want := "tnt_7f3b2/sankhya/sankhya-prod-01/product/10529"
	if got := k.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
	if k.ExternalKey() != "10529" {
		t.Fatalf("ExternalKey() = %q, want %q", k.ExternalKey(), "10529")
	}
}

// Two Sankhya installations are two sources even for the same product code.
func TestSourceInstanceDiscriminates(t *testing.T) {
	a, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	b, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-02", "product", "10529")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	if a.String() == b.String() {
		t.Fatal("two installations produced the same key; the instance is not discriminating")
	}
}

func TestObservationValidateRequiresKeyAndEvidence(t *testing.T) {
	desc, err := fact.NewKnown("Cafeteira Eletrica 30 Xicaras 220 V", ev(t))
	if err != nil {
		t.Fatalf("fact.NewKnown: %v", err)
	}
	key, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	ean, err := contracts.NewIdentifier(contracts.IdentifierEAN, "7891234567890")
	if err != nil {
		t.Fatalf("NewIdentifier: %v", err)
	}

	obs := contracts.ProductObservation{
		Key:         key,
		Description: desc,
		Identifiers: []contracts.Identifier{ean},
		Evidence:    ev(t),
	}
	if err := obs.Validate(); err != nil {
		t.Fatalf("Validate on a complete observation: %v", err)
	}

	noEvidence := obs
	noEvidence.Evidence = provenance.Evidence{}
	if err := noEvidence.Validate(); err == nil {
		t.Fatal("Validate accepted an observation with no evidence")
	}

	noKey := obs
	noKey.Key = contracts.SourceProductKey{}
	if err := noKey.Validate(); err == nil {
		t.Fatal("Validate accepted an observation with no source key")
	}
}

// A source that returns an empty description has not told us the product has no
// name. It has told us nothing, and those are different — which is why the
// field is a Fact and not a string.
func TestObservationAcceptsAnUnknownDescription(t *testing.T) {
	desc, err := fact.NewUnknown[string]("source returned no description", ev(t))
	if err != nil {
		t.Fatalf("fact.NewUnknown: %v", err)
	}
	key, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	obs := contracts.ProductObservation{Key: key, Description: desc, Evidence: ev(t)}
	if err := obs.Validate(); err != nil {
		t.Fatalf("Validate rejected a legitimately unknown description: %v", err)
	}
}
```

- [ ] **Step 2: Corre e confirma que falha**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/contexts/... -v
```

Esperado: `undefined: contracts.NewIdentifier`. Cola o output.

- [ ] **Step 3: Implementa `identifier.go`**

`apps/server_core/internal/contexts/catalog/contracts/identifier.go`:

```go
// Package contracts is Catalog's published surface: everything another context
// is allowed to know about a product. The rest lives under internal/ and the
// compiler is what keeps it there.
package contracts

import (
	"errors"
	"fmt"
	"strings"

	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// IdentifierKind is the sort of identifier, not its authority. An EAN can be
// missing, duplicated, recycled, or simply wrong at the source; it is evidence
// for a linking decision, never the identity of a product.
type IdentifierKind string

const (
	// IdentifierEAN is a barcode.
	IdentifierEAN IdentifierKind = "ean"
	// IdentifierSKU is the seller's own stock code.
	IdentifierSKU IdentifierKind = "sku"
	// IdentifierSourceCode is the source system's primary code.
	IdentifierSourceCode IdentifierKind = "source_code"
)

var (
	// ErrUnknownIdentifierKind is returned for a kind outside the closed set.
	ErrUnknownIdentifierKind = errors.New("catalog: unknown identifier kind")
	// ErrBlank is returned when a required part carries no content.
	ErrBlank = errors.New("catalog: value is empty")
)

// Identifier is one identifying string of a known kind.
type Identifier struct {
	kind  IdentifierKind
	value string
}

// NewIdentifier validates both the kind and the value.
func NewIdentifier(kind IdentifierKind, value string) (Identifier, error) {
	switch kind {
	case IdentifierEAN, IdentifierSKU, IdentifierSourceCode:
	default:
		return Identifier{}, fmt.Errorf("%w: %q", ErrUnknownIdentifierKind, kind)
	}
	v := strings.TrimSpace(value)
	if v == "" {
		return Identifier{}, fmt.Errorf("%w: identifier of kind %q", ErrBlank, kind)
	}
	return Identifier{kind: kind, value: v}, nil
}

// Kind returns the identifier kind.
func (i Identifier) Kind() IdentifierKind { return i.kind }

// Value returns the identifier string.
func (i Identifier) Value() string { return i.value }

// String renders "kind:value".
func (i Identifier) String() string { return string(i.kind) + ":" + i.value }

// SourceProductKey is the address of an object inside a source system, and the
// only place a source system's own code is allowed to live.
//
// The module tree this replaces declared `type InternalProductID int`, and that
// int WAS Sankhya's CODPROD, described as canonical. While that holds, adding a
// second ERP is not additive — it is an identity collision running through
// every contract. Instance is part of the key for the same reason: two Sankhya
// installations are two sources even when the product code matches.
type SourceProductKey struct {
	tenant      tenant.ID
	system      string
	instance    string
	objectKind  string
	externalKey string
}

// NewSourceProductKey validates every part. There is no partial key.
func NewSourceProductKey(t tenant.ID, system, instance, objectKind, externalKey string) (SourceProductKey, error) {
	if t.IsZero() {
		return SourceProductKey{}, fmt.Errorf("%w: tenant", ErrBlank)
	}
	system = strings.TrimSpace(system)
	instance = strings.TrimSpace(instance)
	objectKind = strings.TrimSpace(objectKind)
	externalKey = strings.TrimSpace(externalKey)
	switch {
	case system == "":
		return SourceProductKey{}, fmt.Errorf("%w: source system", ErrBlank)
	case instance == "":
		return SourceProductKey{}, fmt.Errorf("%w: source instance", ErrBlank)
	case objectKind == "":
		return SourceProductKey{}, fmt.Errorf("%w: object kind", ErrBlank)
	case externalKey == "":
		return SourceProductKey{}, fmt.Errorf("%w: external key", ErrBlank)
	}
	return SourceProductKey{
		tenant:      t,
		system:      system,
		instance:    instance,
		objectKind:  objectKind,
		externalKey: externalKey,
	}, nil
}

// Tenant returns the owning tenant.
func (k SourceProductKey) Tenant() tenant.ID { return k.tenant }

// System returns the source system name.
func (k SourceProductKey) System() string { return k.system }

// Instance returns which installation of that system.
func (k SourceProductKey) Instance() string { return k.instance }

// ObjectKind returns what kind of object the key addresses.
func (k SourceProductKey) ObjectKind() string { return k.objectKind }

// ExternalKey returns the source system's own code.
func (k SourceProductKey) ExternalKey() string { return k.externalKey }

// IsZero reports whether this is the zero value rather than a built key.
func (k SourceProductKey) IsZero() bool { return k.system == "" }

// String is the stable storage and log form.
func (k SourceProductKey) String() string {
	if k.IsZero() {
		return ""
	}
	return strings.Join([]string{
		k.tenant.String(), k.system, k.instance, k.objectKind, k.externalKey,
	}, "/")
}
```

- [ ] **Step 4: Implementa `observation.go`**

`apps/server_core/internal/contexts/catalog/contracts/observation.go`:

```go
package contracts

import (
	"fmt"

	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

// ProductObservation is what a source feed hands to Catalog: one product as one
// system saw it at one moment. It is not a product — Catalog decides that.
type ProductObservation struct {
	Key         SourceProductKey
	Description fact.Fact[string]
	Identifiers []Identifier
	Evidence    provenance.Evidence
}

// Validate rejects an observation that cannot be acted on. It deliberately does
// NOT require a known description: a source that returned nothing is a fact
// about the source, and Catalog records it as Unknown rather than refusing the
// product or inventing a name.
func (o ProductObservation) Validate() error {
	if o.Key.IsZero() {
		return fmt.Errorf("%w: source product key", ErrBlank)
	}
	if o.Evidence.IsZero() {
		return fmt.Errorf("%w: evidence", ErrBlank)
	}
	for i, id := range o.Identifiers {
		if id.Value() == "" {
			return fmt.Errorf("%w: identifier at index %d", ErrBlank, i)
		}
	}
	return nil
}

// Disposition is what ingesting an observation did.
type Disposition string

const (
	// DispositionCreated means a new canonical product was minted.
	DispositionCreated Disposition = "created"
	// DispositionChanged means an existing product gained a new version.
	DispositionChanged Disposition = "changed"
	// DispositionIdempotent means this exact payload was already recorded and
	// nothing moved. Re-polling must be free.
	DispositionIdempotent Disposition = "idempotent"
)

// IngestResult reports what happened, including conflicts a human must see.
//
// DuplicateIdentifiers is not an error. Catalog still creates a distinct
// product, because silently merging two ERP codes that share a bad EAN is a
// data loss that cannot be undone; surfacing the conflict is reversible.
type IngestResult struct {
	ProductID            string
	Disposition          Disposition
	Version              int
	DuplicateIdentifiers []Identifier
}
```

- [ ] **Step 5: Implementa `port/reader.go`**

`apps/server_core/internal/contexts/catalog/port/reader.go`:

```go
// Package port carries what other contexts may ASK Catalog. A question, never a
// table: a consumer asks "which product is this?" and does not join to
// catalog.products. That is the difference that lets Catalog change its storage
// without breaking anybody, and the reason there is no foreign key across
// context schemas.
package port

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// Summary answers "what is this product?", flattened for a consumer that has no
// business knowing Catalog's internal model.
type Summary struct {
	ProductID   string
	Description string
	Identifiers []contracts.Identifier
	Version     int
}

// Reader answers identity questions about products.
type Reader interface {
	// ByProductID returns the summary, and false when no such product exists
	// for this tenant. Not existing is an answer, not an error.
	ByProductID(ctx context.Context, t tenant.ID, productID string) (Summary, bool, error)

	// ByIdentifier returns every product carrying this identifier.
	//
	// The slice is the point. More than one is a real and expected answer,
	// because identifiers are evidence and not keys. A signature returning a
	// single product would force this method to pick, and picking silently is
	// how a duplicated EAN becomes a wrong link nobody can see.
	ByIdentifier(ctx context.Context, t tenant.ID, id contracts.Identifier) ([]Summary, error)
}
```

- [ ] **Step 6: Corre e confirma que passa**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/contexts/... -v && go vet ./internal/contexts/...
```

Esperado: sete testes de topo PASS, quatro subtestes PASS, `go vet` sem saída.

- [ ] **Step 7: Os detetores deixam de saltar**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/arch/... -v
```

Esperado: os três `SKIP` da Tarefa 6 passam a **PASS**, porque `internal/contexts` já existe. Confirma essa transição explicitamente — se continuarem a saltar, o detetor não está a olhar para nada e o verde da Tarefa 6 nunca valeu.

- [ ] **Step 8: Commit**

```bash
git status --porcelain --untracked-files=all
git add apps/server_core/internal/contexts/catalog
git commit -m "feat(catalog): superficie publicada — o codigo do ERP vive na SourceProductKey e em mais lado nenhum"
```

---

## Tarefa 8: `catalog/internal/domain` — identidade opaca e versão

O identificador canónico é gerado por nós e não deriva de nada da fonte. A versão sobe quando o payload muda e não sobe quando não muda — e "não muda" é uma comparação de hash, não de campos, porque comparar campos é decidir quais é que contam e essa decisão apodrece.

**Files:**
- Create: `apps/server_core/internal/contexts/catalog/internal/domain/id.go`
- Create: `apps/server_core/internal/contexts/catalog/internal/domain/product.go`
- Create: `apps/server_core/internal/contexts/catalog/internal/domain/product_test.go`

**Interfaces:**
- Produces: `domain.ProductID` com `domain.NewProductID(string) (ProductID, error)`, `ProductID.String()`, `ProductID.IsZero()`; `domain.Product` com `domain.NewProduct(ProductID, contracts.ProductObservation) (Product, error)`, e os métodos `ID()`, `Tenant()`, `Version()`, `Description() fact.Fact[string]`, `Identifiers() []contracts.Identifier`, `SourceKeys() []contracts.SourceProductKey`, `LastPayloadHash() string`, `Apply(contracts.ProductObservation) (Product, contracts.Disposition, error)`.
- Consumes: `contracts.*` (T7), `fact.Fact[string]` (T5), `tenant.ID` (T1).

- [ ] **Step 1: Escreve o teste que falha**

`apps/server_core/internal/contexts/catalog/internal/domain/product_test.go`:

```go
package domain_test

import (
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/domain"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func tid(t *testing.T) tenant.ID {
	t.Helper()
	id, err := tenant.Parse("tnt_7f3b2")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	return id
}

// obs builds an observation with the given description and payload hash.
func obs(t *testing.T, description, hash string) contracts.ProductObservation {
	t.Helper()
	e, err := provenance.NewEvidence("sankhya", "product", "10529",
		time.Date(2026, 8, 6, 12, 15, 0, 0, time.UTC), hash)
	if err != nil {
		t.Fatalf("provenance.NewEvidence: %v", err)
	}
	desc, err := fact.NewKnown(description, e)
	if err != nil {
		t.Fatalf("fact.NewKnown: %v", err)
	}
	key, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("contracts.NewSourceProductKey: %v", err)
	}
	ean, err := contracts.NewIdentifier(contracts.IdentifierEAN, "7891234567890")
	if err != nil {
		t.Fatalf("contracts.NewIdentifier: %v", err)
	}
	return contracts.ProductObservation{
		Key:         key,
		Description: desc,
		Identifiers: []contracts.Identifier{ean},
		Evidence:    e,
	}
}

func pid(t *testing.T) domain.ProductID {
	t.Helper()
	id, err := domain.NewProductID("prd_0198f65fc4477d4f9e9710c2b8ff922a")
	if err != nil {
		t.Fatalf("domain.NewProductID: %v", err)
	}
	return id
}

func TestNewProductIDRejectsBlank(t *testing.T) {
	if _, err := domain.NewProductID("  "); err == nil {
		t.Fatal("NewProductID accepted a blank id")
	}
}

// The canonical id must carry no source semantics. A caller that can pass the
// ERP code as the product id has already lost the property this type exists for.
func TestNewProductIDRejectsABareSourceCode(t *testing.T) {
	if _, err := domain.NewProductID("10529"); err == nil {
		t.Fatal("NewProductID accepted a bare numeric source code as a canonical id")
	}
}

func TestNewProductStartsAtVersionOne(t *testing.T) {
	p, err := domain.NewProduct(pid(t), obs(t, "Cafeteira Eletrica 30 Xicaras 220 V", "sha256:ab91"))
	if err != nil {
		t.Fatalf("NewProduct: %v", err)
	}
	if p.Version() != 1 {
		t.Fatalf("Version() = %d, want 1", p.Version())
	}
	if p.ID().String() != "prd_0198f65fc4477d4f9e9710c2b8ff922a" {
		t.Fatalf("ID() = %q", p.ID().String())
	}
	if p.Tenant().String() != "tnt_7f3b2" {
		t.Fatalf("Tenant() = %q, want tnt_7f3b2", p.Tenant().String())
	}
	if p.LastPayloadHash() != "sha256:ab91" {
		t.Fatalf("LastPayloadHash() = %q", p.LastPayloadHash())
	}
	if len(p.SourceKeys()) != 1 {
		t.Fatalf("SourceKeys() has %d entries, want 1", len(p.SourceKeys()))
	}
}

func TestNewProductRejectsAnInvalidObservation(t *testing.T) {
	bad := obs(t, "x", "sha256:ab91")
	bad.Evidence = provenance.Evidence{}
	if _, err := domain.NewProduct(pid(t), bad); err == nil {
		t.Fatal("NewProduct accepted an observation with no evidence")
	}
}

// Re-polling must be free. Same payload hash, same everything: no new version.
func TestApplySamePayloadIsIdempotent(t *testing.T) {
	p, err := domain.NewProduct(pid(t), obs(t, "Cafeteira Eletrica 30 Xicaras 220 V", "sha256:ab91"))
	if err != nil {
		t.Fatalf("NewProduct: %v", err)
	}
	next, disp, err := p.Apply(obs(t, "Cafeteira Eletrica 30 Xicaras 220 V", "sha256:ab91"))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if disp != contracts.DispositionIdempotent {
		t.Fatalf("Disposition = %q, want %q", disp, contracts.DispositionIdempotent)
	}
	if next.Version() != 1 {
		t.Fatalf("Version() = %d after an idempotent apply, want 1", next.Version())
	}
}

func TestApplyChangedPayloadBumpsVersion(t *testing.T) {
	p, err := domain.NewProduct(pid(t), obs(t, "Cafeteira Eletrica 30 Xicaras 220 V", "sha256:ab91"))
	if err != nil {
		t.Fatalf("NewProduct: %v", err)
	}
	next, disp, err := p.Apply(obs(t, "Cafeteira Eletrica Inox 30 Xicaras 220 V", "sha256:cf02"))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if disp != contracts.DispositionChanged {
		t.Fatalf("Disposition = %q, want %q", disp, contracts.DispositionChanged)
	}
	if next.Version() != 2 {
		t.Fatalf("Version() = %d, want 2", next.Version())
	}
	desc, ok := next.Description().Value()
	if !ok || !strings.Contains(desc, "Inox") {
		t.Fatalf("Description() = %q, ok=%v; want the new description", desc, ok)
	}
}

// Apply returns a new value. The receiver must not move, or a caller that
// persists the old one after a failed transaction writes a version that the
// database never saw.
func TestApplyDoesNotMutateTheReceiver(t *testing.T) {
	p, err := domain.NewProduct(pid(t), obs(t, "original", "sha256:ab91"))
	if err != nil {
		t.Fatalf("NewProduct: %v", err)
	}
	if _, _, err := p.Apply(obs(t, "changed", "sha256:cf02")); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if p.Version() != 1 {
		t.Fatalf("receiver Version() = %d after Apply, want 1", p.Version())
	}
	desc, _ := p.Description().Value()
	if desc != "original" {
		t.Fatalf("receiver Description() = %q, want %q", desc, "original")
	}
}

// A second source key for the same canonical product is added, not replaced.
func TestApplyAccumulatesSourceKeys(t *testing.T) {
	p, err := domain.NewProduct(pid(t), obs(t, "original", "sha256:ab91"))
	if err != nil {
		t.Fatalf("NewProduct: %v", err)
	}
	second := obs(t, "original", "sha256:cf02")
	second.Key, err = contracts.NewSourceProductKey(tid(t), "spreadsheet", "import-2026-08", "product", "ROW-42")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	next, _, err := p.Apply(second)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(next.SourceKeys()) != 2 {
		t.Fatalf("SourceKeys() has %d entries, want 2", len(next.SourceKeys()))
	}
}

// An observation belonging to another tenant is not a version, it is a bug.
func TestApplyRejectsAForeignTenant(t *testing.T) {
	p, err := domain.NewProduct(pid(t), obs(t, "original", "sha256:ab91"))
	if err != nil {
		t.Fatalf("NewProduct: %v", err)
	}
	other, err := tenant.Parse("tnt_other")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	foreign := obs(t, "original", "sha256:cf02")
	foreign.Key, err = contracts.NewSourceProductKey(other, "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("NewSourceProductKey: %v", err)
	}
	if _, _, err := p.Apply(foreign); err == nil {
		t.Fatal("Apply accepted an observation from another tenant")
	}
}
```

- [ ] **Step 2: Corre e confirma que falha**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/contexts/catalog/... -v
```

Esperado: `undefined: domain.NewProductID`. Cola o output.

- [ ] **Step 3: Implementa `id.go`**

`apps/server_core/internal/contexts/catalog/internal/domain/id.go`:

```go
// Package domain holds Catalog's model and its invariants. It is under
// internal/, so no other context can import it — not by convention, by the
// compiler refusing to link.
package domain

import (
	"errors"
	"fmt"
	"strings"
)

// ProductIDPrefix marks a canonical product identifier. The prefix is not
// decoration: it makes "someone passed the ERP code as the product id"
// detectable at the constructor instead of three joins downstream.
const ProductIDPrefix = "prd_"

var (
	// ErrBlankProductID is returned for an empty identifier.
	ErrBlankProductID = errors.New("catalog: product id is empty")
	// ErrNotCanonicalProductID is returned when the identifier does not carry
	// the canonical prefix.
	ErrNotCanonicalProductID = errors.New("catalog: product id is not canonical")
)

// ProductID is the platform's own identifier for a product. It is opaque and it
// derives from nothing at the source: no ERP code, no EAN, no SKU.
type ProductID struct {
	value string
}

// NewProductID validates the shape of a canonical identifier.
func NewProductID(s string) (ProductID, error) {
	v := strings.TrimSpace(s)
	if v == "" {
		return ProductID{}, ErrBlankProductID
	}
	if !strings.HasPrefix(v, ProductIDPrefix) || len(v) <= len(ProductIDPrefix) {
		return ProductID{}, fmt.Errorf("%w: %q has no %q prefix", ErrNotCanonicalProductID, v, ProductIDPrefix)
	}
	return ProductID{value: v}, nil
}

// String returns the identifier, or the empty string for the zero value.
func (p ProductID) String() string { return p.value }

// IsZero reports whether this is the zero value.
func (p ProductID) IsZero() bool { return p.value == "" }
```

- [ ] **Step 4: Implementa `product.go`**

`apps/server_core/internal/contexts/catalog/internal/domain/product.go`:

```go
package domain

import (
	"errors"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

var (
	// ErrForeignTenant is returned when an observation belongs to another tenant.
	ErrForeignTenant = errors.New("catalog: observation belongs to another tenant")
	// ErrNoProductID is returned when a product is built without an identifier.
	ErrNoProductID = errors.New("catalog: product id is required")
)

// Product is one canonical product. Every field is unexported and every method
// that changes state returns a NEW Product: a caller that persists inside a
// transaction which then fails must not be holding a mutated aggregate.
type Product struct {
	id          ProductID
	tenant      tenant.ID
	version     int
	description fact.Fact[string]
	identifiers []contracts.Identifier
	sourceKeys  []contracts.SourceProductKey
	lastHash    string
}

// NewProduct mints version 1 from a first observation.
func NewProduct(id ProductID, o contracts.ProductObservation) (Product, error) {
	if id.IsZero() {
		return Product{}, ErrNoProductID
	}
	if err := o.Validate(); err != nil {
		return Product{}, err
	}
	return Product{
		id:          id,
		tenant:      o.Key.Tenant(),
		version:     1,
		description: o.Description,
		identifiers: append([]contracts.Identifier(nil), o.Identifiers...),
		sourceKeys:  []contracts.SourceProductKey{o.Key},
		lastHash:    o.Evidence.PayloadHash(),
	}, nil
}

// ID returns the canonical identifier.
func (p Product) ID() ProductID { return p.id }

// Tenant returns the owning tenant.
func (p Product) Tenant() tenant.ID { return p.tenant }

// Version returns the catalog version, which counts substantive changes only.
func (p Product) Version() int { return p.version }

// Description returns the description as a fact, which may be Unknown.
func (p Product) Description() fact.Fact[string] { return p.description }

// Identifiers returns a copy of the identifier list.
func (p Product) Identifiers() []contracts.Identifier {
	return append([]contracts.Identifier(nil), p.identifiers...)
}

// SourceKeys returns a copy of every source address this product answers to.
func (p Product) SourceKeys() []contracts.SourceProductKey {
	return append([]contracts.SourceProductKey(nil), p.sourceKeys...)
}

// LastPayloadHash returns the hash of the last payload that changed this product.
func (p Product) LastPayloadHash() string { return p.lastHash }

// Apply folds a new observation in and reports what it did.
//
// Sameness is decided by the raw payload hash and not by comparing fields.
// Comparing fields means choosing which fields count, and that choice rots:
// the day a field is added, every product silently stops changing on it.
func (p Product) Apply(o contracts.ProductObservation) (Product, contracts.Disposition, error) {
	if err := o.Validate(); err != nil {
		return p, "", err
	}
	if o.Key.Tenant() != p.tenant {
		return p, "", fmt.Errorf("%w: product %s is %s, observation is %s",
			ErrForeignTenant, p.id, p.tenant, o.Key.Tenant())
	}
	if o.Evidence.PayloadHash() == p.lastHash {
		return p, contracts.DispositionIdempotent, nil
	}

	next := Product{
		id:          p.id,
		tenant:      p.tenant,
		version:     p.version + 1,
		description: o.Description,
		identifiers: mergeIdentifiers(p.identifiers, o.Identifiers),
		sourceKeys:  mergeSourceKeys(p.sourceKeys, o.Key),
		lastHash:    o.Evidence.PayloadHash(),
	}
	return next, contracts.DispositionChanged, nil
}

// mergeIdentifiers unions by kind+value, keeping the existing order stable.
func mergeIdentifiers(existing, incoming []contracts.Identifier) []contracts.Identifier {
	seen := make(map[string]bool, len(existing)+len(incoming))
	out := make([]contracts.Identifier, 0, len(existing)+len(incoming))
	for _, list := range [][]contracts.Identifier{existing, incoming} {
		for _, id := range list {
			k := id.String()
			if seen[k] {
				continue
			}
			seen[k] = true
			out = append(out, id)
		}
	}
	return out
}

// mergeSourceKeys appends the key unless this product already answers to it. A
// second source is added, never substituted: one canonical product legitimately
// has an ERP address and a spreadsheet address at the same time.
func mergeSourceKeys(existing []contracts.SourceProductKey, k contracts.SourceProductKey) []contracts.SourceProductKey {
	for _, have := range existing {
		if have.String() == k.String() {
			return append([]contracts.SourceProductKey(nil), existing...)
		}
	}
	return append(append([]contracts.SourceProductKey(nil), existing...), k)
}
```

`Product.id` é usado num `%s` do `fmt.Errorf`; `ProductID` tem `String()`, por isso formata. `tenant.ID` também.

- [ ] **Step 5: Corre e confirma que passa**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/contexts/catalog/... -v
```

Esperado: dez testes PASS entre `contracts` e `domain`.

- [ ] **Step 6: Prova a fronteira ao vivo, uma vez**

Este passo mede a regra 1.1 no nosso repositório em vez de a acreditar. Cria o ficheiro:

`apps/server_core/internal/contexts/catalog/contracts/zz_probe.go`

```go
package contracts

import _ "marketplace-central/apps/server_core/internal/contexts/catalog/internal/domain"
```

Isto é legal — `contracts` está na árvore de `catalog`. Agora corre:

```bash
cd apps/server_core && go build ./internal/contexts/...
```

Esperado: **exit 0**. Depois muda o import para um contexto que não existe ainda e confirma que a mensagem é de pacote em falta e não de `internal` — a prova a sério da fronteira entre contextos vem na Tarefa 11, quando existir um segundo contexto para a violar. Apaga a sonda:

```bash
rm apps/server_core/internal/contexts/catalog/contracts/zz_probe.go
go build ./internal/contexts/...
```

Esperado: exit 0 e `git status --porcelain --untracked-files=all` sem a sonda.

- [ ] **Step 7: Commit**

```bash
git status --porcelain --untracked-files=all
git add apps/server_core/internal/contexts/catalog/internal
git commit -m "feat(catalog): ProductID opaco e Product com versao por hash de payload"
```

---

## Tarefa 9: `catalog/internal/application` e `module.go` — o caso de uso

A aplicação decide: chave de fonte conhecida ou nova, produto novo ou versão nova, e o que fazer quando um identificador já pertence a outro produto. A resposta a esta última é a que custa se estiver errada — fundir dois códigos de ERP em silêncio porque partilham um EAN mau é perda de dados irreversível; criar dois produtos e reportar o conflito é reversível.

**Files:**
- Create: `apps/server_core/internal/contexts/catalog/internal/application/ports.go`
- Create: `apps/server_core/internal/contexts/catalog/internal/application/ingest.go`
- Create: `apps/server_core/internal/contexts/catalog/internal/application/ingest_test.go`
- Create: `apps/server_core/internal/contexts/catalog/internal/application/memstore_test.go`
- Create: `apps/server_core/internal/contexts/catalog/module.go`

**Interfaces:**
- Produces: `application.Store` (interface: `BySourceKey`, `ByIdentifier`, `Insert`, `Update`), `application.IDFactory` (interface: `NewProductID`), `application.Service` com `application.NewService(Store, IDFactory) *Service` e `Service.Ingest(ctx, contracts.ProductObservation) (contracts.IngestResult, error)`; `catalog.Module` com `catalog.New(application.Store, application.IDFactory) *Module`, `Module.IngestProduct(ctx, contracts.ProductObservation) (contracts.IngestResult, error)` e `Module.Reader() port.Reader`.
- Consumes: `domain.*` (T8), `contracts.*` (T7).

- [ ] **Step 1: Escreve o duplo de teste**

`apps/server_core/internal/contexts/catalog/internal/application/memstore_test.go`:

```go
package application_test

import (
	"context"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/application"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/domain"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// memStore is an in-memory Store. It is a test double for the decision logic
// only: it proves nothing about persistence, atomicity or isolation, and the
// integration lane in Task 10 is what covers those.
type memStore struct {
	byID  map[string]domain.Product
	order []string
}

func newMemStore() *memStore {
	return &memStore{byID: map[string]domain.Product{}}
}

func (m *memStore) BySourceKey(_ context.Context, k contracts.SourceProductKey) (domain.Product, bool, error) {
	for _, id := range m.order {
		p := m.byID[id]
		for _, have := range p.SourceKeys() {
			if have.String() == k.String() {
				return p, true, nil
			}
		}
	}
	return domain.Product{}, false, nil
}

func (m *memStore) ByIdentifier(_ context.Context, t tenant.ID, id contracts.Identifier) ([]domain.Product, error) {
	var out []domain.Product
	for _, key := range m.order {
		p := m.byID[key]
		if p.Tenant() != t {
			continue
		}
		for _, have := range p.Identifiers() {
			if have.String() == id.String() {
				out = append(out, p)
				break
			}
		}
	}
	return out, nil
}

func (m *memStore) Insert(_ context.Context, p domain.Product) error {
	if _, exists := m.byID[p.ID().String()]; exists {
		return fmt.Errorf("memstore: duplicate insert of %s", p.ID())
	}
	m.byID[p.ID().String()] = p
	m.order = append(m.order, p.ID().String())
	return nil
}

func (m *memStore) Update(_ context.Context, p domain.Product) error {
	if _, exists := m.byID[p.ID().String()]; !exists {
		return fmt.Errorf("memstore: update of unknown %s", p.ID())
	}
	m.byID[p.ID().String()] = p
	return nil
}

// seqIDs hands out predictable canonical ids so a test can name them.
type seqIDs struct{ n int }

func (s *seqIDs) NewProductID() (domain.ProductID, error) {
	s.n++
	return domain.NewProductID(fmt.Sprintf("%s%032d", domain.ProductIDPrefix, s.n))
}

var _ application.Store = (*memStore)(nil)
var _ application.IDFactory = (*seqIDs)(nil)
```

- [ ] **Step 2: Escreve o teste que falha**

`apps/server_core/internal/contexts/catalog/internal/application/ingest_test.go`:

```go
package application_test

import (
	"context"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/application"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func tid(t *testing.T) tenant.ID {
	t.Helper()
	id, err := tenant.Parse("tnt_7f3b2")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	return id
}

// observation builds one product observation.
func observation(t *testing.T, externalKey, description, ean, hash string) contracts.ProductObservation {
	t.Helper()
	e, err := provenance.NewEvidence("sankhya", "product", externalKey,
		time.Date(2026, 8, 6, 12, 15, 0, 0, time.UTC), hash)
	if err != nil {
		t.Fatalf("provenance.NewEvidence: %v", err)
	}
	desc, err := fact.NewKnown(description, e)
	if err != nil {
		t.Fatalf("fact.NewKnown: %v", err)
	}
	key, err := contracts.NewSourceProductKey(tid(t), "sankhya", "sankhya-prod-01", "product", externalKey)
	if err != nil {
		t.Fatalf("contracts.NewSourceProductKey: %v", err)
	}
	id, err := contracts.NewIdentifier(contracts.IdentifierEAN, ean)
	if err != nil {
		t.Fatalf("contracts.NewIdentifier: %v", err)
	}
	return contracts.ProductObservation{
		Key: key, Description: desc,
		Identifiers: []contracts.Identifier{id}, Evidence: e,
	}
}

func newService() (*application.Service, *memStore) {
	store := newMemStore()
	return application.NewService(store, &seqIDs{}), store
}

func TestIngestFirstObservationCreates(t *testing.T) {
	svc, _ := newService()
	got, err := svc.Ingest(context.Background(),
		observation(t, "10529", "Cafeteira Eletrica 30 Xicaras 220 V", "7891234567890", "sha256:ab91"))
	if err != nil {
		t.Fatalf("Ingest: %v", err)
	}
	if got.Disposition != contracts.DispositionCreated {
		t.Fatalf("Disposition = %q, want %q", got.Disposition, contracts.DispositionCreated)
	}
	if got.Version != 1 {
		t.Fatalf("Version = %d, want 1", got.Version)
	}
	if got.ProductID == "10529" {
		t.Fatal("the canonical id is the ERP code; identity leaked from the source")
	}
	if len(got.DuplicateIdentifiers) != 0 {
		t.Fatalf("DuplicateIdentifiers = %v, want none", got.DuplicateIdentifiers)
	}
}

func TestIngestSamePayloadIsIdempotent(t *testing.T) {
	svc, _ := newService()
	obs := observation(t, "10529", "Cafeteira", "7891234567890", "sha256:ab91")
	first, err := svc.Ingest(context.Background(), obs)
	if err != nil {
		t.Fatalf("first Ingest: %v", err)
	}
	second, err := svc.Ingest(context.Background(), obs)
	if err != nil {
		t.Fatalf("second Ingest: %v", err)
	}
	if second.Disposition != contracts.DispositionIdempotent {
		t.Fatalf("Disposition = %q, want %q", second.Disposition, contracts.DispositionIdempotent)
	}
	if second.ProductID != first.ProductID {
		t.Fatalf("ProductID moved between polls: %q then %q", first.ProductID, second.ProductID)
	}
	if second.Version != 1 {
		t.Fatalf("Version = %d after an idempotent poll, want 1", second.Version)
	}
}

func TestIngestChangedPayloadBumpsVersionAndKeepsIdentity(t *testing.T) {
	svc, _ := newService()
	first, err := svc.Ingest(context.Background(),
		observation(t, "10529", "Cafeteira Eletrica 30 Xicaras 220 V", "7891234567890", "sha256:ab91"))
	if err != nil {
		t.Fatalf("first Ingest: %v", err)
	}
	second, err := svc.Ingest(context.Background(),
		observation(t, "10529", "Cafeteira Eletrica Inox 30 Xicaras 220 V", "7891234567890", "sha256:cf02"))
	if err != nil {
		t.Fatalf("second Ingest: %v", err)
	}
	if second.Disposition != contracts.DispositionChanged {
		t.Fatalf("Disposition = %q, want %q", second.Disposition, contracts.DispositionChanged)
	}
	if second.Version != 2 {
		t.Fatalf("Version = %d, want 2", second.Version)
	}
	if second.ProductID != first.ProductID {
		t.Fatalf("a changed description minted a new identity: %q then %q", first.ProductID, second.ProductID)
	}
}

// The load-bearing test. ERP code 10530 reports the SAME EAN as 10529 because
// the master data is bad. Two source keys are two products. Merging them
// silently is irreversible; reporting the duplicate is not.
func TestIngestDuplicateEanDoesNotMerge(t *testing.T) {
	svc, _ := newService()
	first, err := svc.Ingest(context.Background(),
		observation(t, "10529", "Cafeteira A", "7891234567890", "sha256:ab91"))
	if err != nil {
		t.Fatalf("first Ingest: %v", err)
	}
	second, err := svc.Ingest(context.Background(),
		observation(t, "10530", "Cafeteira B", "7891234567890", "sha256:dd77"))
	if err != nil {
		t.Fatalf("second Ingest: %v", err)
	}
	if second.Disposition != contracts.DispositionCreated {
		t.Fatalf("Disposition = %q, want %q; a duplicate EAN must not fold into the existing product",
			second.Disposition, contracts.DispositionCreated)
	}
	if second.ProductID == first.ProductID {
		t.Fatal("two ERP codes sharing a bad EAN were merged into one product")
	}
	if len(second.DuplicateIdentifiers) != 1 {
		t.Fatalf("DuplicateIdentifiers = %v, want exactly the shared EAN", second.DuplicateIdentifiers)
	}
	if second.DuplicateIdentifiers[0].Value() != "7891234567890" {
		t.Fatalf("DuplicateIdentifiers[0] = %q", second.DuplicateIdentifiers[0].Value())
	}
}

// Another tenant reporting the same EAN is not a duplicate. It is a different
// company's catalogue, and a conflict reported across tenants is a data leak.
func TestIngestDoesNotReportDuplicatesAcrossTenants(t *testing.T) {
	svc, _ := newService()
	if _, err := svc.Ingest(context.Background(),
		observation(t, "10529", "Cafeteira A", "7891234567890", "sha256:ab91")); err != nil {
		t.Fatalf("first Ingest: %v", err)
	}

	other, err := tenant.Parse("tnt_other")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	obs := observation(t, "10529", "Cafeteira A", "7891234567890", "sha256:ee88")
	obs.Key, err = contracts.NewSourceProductKey(other, "sankhya", "sankhya-prod-01", "product", "10529")
	if err != nil {
		t.Fatalf("contracts.NewSourceProductKey: %v", err)
	}

	got, err := svc.Ingest(context.Background(), obs)
	if err != nil {
		t.Fatalf("second Ingest: %v", err)
	}
	if len(got.DuplicateIdentifiers) != 0 {
		t.Fatalf("DuplicateIdentifiers = %v across tenants; that is a leak", got.DuplicateIdentifiers)
	}
}

func TestIngestRejectsAnInvalidObservation(t *testing.T) {
	svc, _ := newService()
	bad := observation(t, "10529", "Cafeteira", "7891234567890", "sha256:ab91")
	bad.Evidence = provenance.Evidence{}
	if _, err := svc.Ingest(context.Background(), bad); err == nil {
		t.Fatal("Ingest accepted an observation with no evidence")
	}
}
```

- [ ] **Step 3: Corre e confirma que falha**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/contexts/catalog/internal/application/... -v
```

Esperado: `undefined: application.NewService`. Cola o output.

- [ ] **Step 4: Implementa `ports.go`**

`apps/server_core/internal/contexts/catalog/internal/application/ports.go`:

```go
// Package application holds Catalog's use cases. It depends on interfaces
// declared here and implemented by internal/postgres, so the decision logic is
// testable without a database and the database is replaceable without touching
// a decision.
package application

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/domain"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// Store persists products. It is the ONLY writer of the catalog schema.
type Store interface {
	// BySourceKey resolves a source address to the product it already maps to.
	BySourceKey(ctx context.Context, k contracts.SourceProductKey) (domain.Product, bool, error)

	// ByIdentifier returns every product of this tenant carrying the identifier.
	// Scoped by tenant in the signature and not only in the query, because a
	// cross-tenant answer here is a data leak, not a bug in a WHERE clause.
	ByIdentifier(ctx context.Context, t tenant.ID, id contracts.Identifier) ([]domain.Product, error)

	// Insert writes a new product, its identifiers, its source keys and the
	// observation, atomically.
	Insert(ctx context.Context, p domain.Product) error

	// Update writes a new version of an existing product, atomically.
	Update(ctx context.Context, p domain.Product) error
}

// IDFactory mints canonical product identifiers. It is an interface so a test
// can make them predictable; nothing about the value may depend on the source.
type IDFactory interface {
	NewProductID() (domain.ProductID, error)
}
```

- [ ] **Step 5: Implementa `ingest.go`**

`apps/server_core/internal/contexts/catalog/internal/application/ingest.go`:

```go
package application

import (
	"context"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/domain"
)

// Service is Catalog's ingestion use case.
type Service struct {
	store Store
	ids   IDFactory
}

// NewService wires the use case.
func NewService(store Store, ids IDFactory) *Service {
	return &Service{store: store, ids: ids}
}

// Ingest folds one source observation into the catalogue.
//
// Resolution is by source key and by nothing else. Identifiers are looked up
// only to REPORT a conflict, never to decide identity: an EAN that two ERP
// codes share is bad master data, and folding the second product into the first
// destroys the distinction with no way back.
func (s *Service) Ingest(ctx context.Context, o contracts.ProductObservation) (contracts.IngestResult, error) {
	if err := o.Validate(); err != nil {
		return contracts.IngestResult{}, err
	}

	existing, found, err := s.store.BySourceKey(ctx, o.Key)
	if err != nil {
		return contracts.IngestResult{}, fmt.Errorf("catalog: resolve source key: %w", err)
	}

	if found {
		next, disposition, applyErr := existing.Apply(o)
		if applyErr != nil {
			return contracts.IngestResult{}, applyErr
		}
		if disposition == contracts.DispositionChanged {
			if err := s.store.Update(ctx, next); err != nil {
				return contracts.IngestResult{}, fmt.Errorf("catalog: update product: %w", err)
			}
		}
		return contracts.IngestResult{
			ProductID:   next.ID().String(),
			Disposition: disposition,
			Version:     next.Version(),
		}, nil
	}

	duplicates, err := s.duplicateIdentifiers(ctx, o)
	if err != nil {
		return contracts.IngestResult{}, err
	}

	id, err := s.ids.NewProductID()
	if err != nil {
		return contracts.IngestResult{}, fmt.Errorf("catalog: mint product id: %w", err)
	}
	product, err := domain.NewProduct(id, o)
	if err != nil {
		return contracts.IngestResult{}, err
	}
	if err := s.store.Insert(ctx, product); err != nil {
		return contracts.IngestResult{}, fmt.Errorf("catalog: insert product: %w", err)
	}

	return contracts.IngestResult{
		ProductID:            product.ID().String(),
		Disposition:          contracts.DispositionCreated,
		Version:              product.Version(),
		DuplicateIdentifiers: duplicates,
	}, nil
}

// duplicateIdentifiers reports which of the observation's identifiers already
// belong to another product OF THE SAME TENANT.
func (s *Service) duplicateIdentifiers(ctx context.Context, o contracts.ProductObservation) ([]contracts.Identifier, error) {
	var out []contracts.Identifier
	for _, id := range o.Identifiers {
		matches, err := s.store.ByIdentifier(ctx, o.Key.Tenant(), id)
		if err != nil {
			return nil, fmt.Errorf("catalog: look up identifier %s: %w", id, err)
		}
		if len(matches) > 0 {
			out = append(out, id)
		}
	}
	return out, nil
}
```

- [ ] **Step 6: Implementa `module.go`**

`apps/server_core/internal/contexts/catalog/module.go`:

```go
// Package catalog is the context's façade. The composition root builds a Module
// and hands it to consumers as a port; it never reaches past this file, because
// everything past this file is under internal/.
package catalog

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/application"
	"marketplace-central/apps/server_core/internal/contexts/catalog/port"
)

// Module is Catalog, assembled.
type Module struct {
	service *application.Service
	reader  port.Reader
}

// New assembles the context from its dependencies.
func New(store application.Store, ids application.IDFactory, reader port.Reader) *Module {
	return &Module{service: application.NewService(store, ids), reader: reader}
}

// IngestProduct folds one source observation into the catalogue.
func (m *Module) IngestProduct(ctx context.Context, o contracts.ProductObservation) (contracts.IngestResult, error) {
	return m.service.Ingest(ctx, o)
}

// Reader is Catalog's answer to identity questions from other contexts.
func (m *Module) Reader() port.Reader { return m.reader }
```

- [ ] **Step 7: Corre e confirma que passa**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/contexts/... -v && go vet ./internal/contexts/...
```

Esperado: seis testes de `application` PASS, e os de `contracts` e `domain` continuam PASS.

- [ ] **Step 8: Confirma que os detetores continuam a zero**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/arch/... 
```

Esperado: `ok`. Sem SKIP nos três testes de contexto.

- [ ] **Step 9: Commit**

```bash
git status --porcelain --untracked-files=all
git add apps/server_core/internal/contexts/catalog
git commit -m "feat(catalog): Ingest resolve por chave de fonte; EAN duplicado reporta conflito e nunca funde"
```

---

## Tarefa 10: esquema `catalog`, repositório e a lane hermética

Um escritor por tabela, chave começada em `tenant_id`, e RLS ligada com `FORCE` — porque `WHERE tenant_id = ?` na aplicação é disciplina e a disciplina falha uma vez.

O estado de conhecimento vai para a tabela. Uma coluna `description` anulável não distingue "a fonte disse que não tem nome" de "a fonte não disse nada", e essa distinção é o ADR-017 inteiro.

**Files:**
- Create: `apps/server_core/migrations/0097_catalog_context.sql`
- Create: `apps/server_core/internal/contexts/catalog/internal/postgres/repository.go`
- Create: `apps/server_core/tests/integration/catalog_ingest_test.go`

**Interfaces:**
- Produces: `postgres.NewRepository(*pgxpool.Pool) *Repository` satisfazendo `application.Store`; `postgres.NewSummaryReader(*Repository) *SummaryReader` satisfazendo `port.Reader`; `postgres.NewULIDFactory() *ULIDFactory` satisfazendo `application.IDFactory`. São três porque `Store.ByIdentifier` devolve `[]domain.Product` e `port.Reader.ByIdentifier` devolve `[]port.Summary`, e um tipo não pode ter dois métodos com o mesmo nome.
- Consumes: `application.Store`/`IDFactory` (T9), `port.Reader` (T7), `domain.Product` (T8).

- [ ] **Step 1: Escreve a migração**

`apps/server_core/migrations/0097_catalog_context.sql`:

```sql
-- Catalog context: its own schema, its own writer, no foreign key leaving it.
--
-- Every key leads with tenant_id and RLS is FORCEd, so a query that forgets the
-- tenant returns nothing instead of returning somebody else's catalogue.
CREATE SCHEMA IF NOT EXISTS catalog;

CREATE TABLE IF NOT EXISTS catalog.products (
    tenant_id           text        NOT NULL,
    product_id          text        NOT NULL,
    version             integer     NOT NULL,
    -- The knowledge state is a column because "the source said nothing" and
    -- "the source said it has no name" are different facts. A nullable
    -- description alone cannot tell them apart.
    description_state   text        NOT NULL,
    description_value   text        NULL,
    description_reason  text        NULL,
    last_payload_hash   text        NOT NULL,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT products_pkey PRIMARY KEY (tenant_id, product_id),
    CONSTRAINT products_version_positive CHECK (version >= 1),
    CONSTRAINT products_state_known CHECK (
        description_state IN ('known', 'estimated', 'unknown', 'not_applicable')
    ),
    -- known carries a value; unknown carries a reason and no value.
    CONSTRAINT products_state_consistent CHECK (
        (description_state = 'known'      AND description_value IS NOT NULL)
     OR (description_state = 'estimated'  AND description_value IS NOT NULL AND description_reason IS NOT NULL)
     OR (description_state IN ('unknown', 'not_applicable')
                                          AND description_value IS NULL     AND description_reason IS NOT NULL)
    )
);

CREATE TABLE IF NOT EXISTS catalog.product_identifiers (
    tenant_id  text NOT NULL,
    product_id text NOT NULL,
    kind       text NOT NULL,
    value      text NOT NULL,
    CONSTRAINT product_identifiers_pkey PRIMARY KEY (tenant_id, product_id, kind, value),
    CONSTRAINT product_identifiers_product_fkey
        FOREIGN KEY (tenant_id, product_id)
        REFERENCES catalog.products (tenant_id, product_id) ON DELETE CASCADE
);

-- Deliberately NOT unique on (tenant_id, kind, value): two products sharing a
-- bad EAN is real master data, and a unique index would turn it into an insert
-- failure at 03:00 instead of a conflict a human can look at.
CREATE INDEX IF NOT EXISTS product_identifiers_lookup
    ON catalog.product_identifiers (tenant_id, kind, value);

CREATE TABLE IF NOT EXISTS catalog.source_product_keys (
    tenant_id       text NOT NULL,
    source_system   text NOT NULL,
    source_instance text NOT NULL,
    object_kind     text NOT NULL,
    external_key    text NOT NULL,
    product_id      text NOT NULL,
    CONSTRAINT source_product_keys_pkey
        PRIMARY KEY (tenant_id, source_system, source_instance, object_kind, external_key),
    CONSTRAINT source_product_keys_product_fkey
        FOREIGN KEY (tenant_id, product_id)
        REFERENCES catalog.products (tenant_id, product_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS catalog.source_observations (
    tenant_id       text        NOT NULL,
    product_id      text        NOT NULL,
    payload_hash    text        NOT NULL,
    source_system   text        NOT NULL,
    object_kind     text        NOT NULL,
    external_key    text        NOT NULL,
    observed_at     timestamptz NOT NULL,
    recorded_at     timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT source_observations_pkey PRIMARY KEY (tenant_id, product_id, payload_hash),
    CONSTRAINT source_observations_product_fkey
        FOREIGN KEY (tenant_id, product_id)
        REFERENCES catalog.products (tenant_id, product_id) ON DELETE CASCADE
);

ALTER TABLE catalog.products             ENABLE ROW LEVEL SECURITY;
ALTER TABLE catalog.product_identifiers  ENABLE ROW LEVEL SECURITY;
ALTER TABLE catalog.source_product_keys  ENABLE ROW LEVEL SECURITY;
ALTER TABLE catalog.source_observations  ENABLE ROW LEVEL SECURITY;

-- FORCE, because without it the table owner silently bypasses every policy and
-- the whole mechanism is decorative in exactly the environment we test in.
ALTER TABLE catalog.products             FORCE ROW LEVEL SECURITY;
ALTER TABLE catalog.product_identifiers  FORCE ROW LEVEL SECURITY;
ALTER TABLE catalog.source_product_keys  FORCE ROW LEVEL SECURITY;
ALTER TABLE catalog.source_observations  FORCE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_isolation ON catalog.products;
CREATE POLICY tenant_isolation ON catalog.products
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation ON catalog.product_identifiers;
CREATE POLICY tenant_isolation ON catalog.product_identifiers
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation ON catalog.source_product_keys;
CREATE POLICY tenant_isolation ON catalog.source_product_keys
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));

DROP POLICY IF EXISTS tenant_isolation ON catalog.source_observations;
CREATE POLICY tenant_isolation ON catalog.source_observations
    USING (tenant_id = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.tenant_id', true));
```

`current_setting('app.tenant_id', true)` devolve `NULL` quando não está definido, e `tenant_id = NULL` é `NULL`, logo falso. **Sessão sem tenant vê zero linhas.** É o comportamento que queremos, e o teste do Step 6 é que o prova.

- [ ] **Step 2: Aplica a migração e confirma**

```bash
docker compose -f docker-compose.dev.yml ps
```

Se a stack não estiver de pé, não a levantes por iniciativa própria — escreve o bloqueio em `.mnfs/HARNESS-DEBTS.md`, faz commit da migração e do repositório, e salta para a Tarefa 11, que não depende do Postgres.

Com a stack de pé, aplica pelo caminho que o repositório já usa para migrações e confirma:

```bash
psql "$DATABASE_URL" -c "\dt catalog.*"
```

Esperado: as quatro tabelas.

- [ ] **Step 3: Escreve o repositório**

`apps/server_core/internal/contexts/catalog/internal/postgres/repository.go`:

```go
// Package postgres is Catalog's only writer. No other package in the platform
// issues DML against the catalog schema, and no other schema holds a foreign
// key into it.
package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/domain"
	"marketplace-central/apps/server_core/internal/contexts/catalog/port"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// Repository implements application.Store. The context's published reader is
// SummaryReader at the bottom of this file, which wraps it.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository wires the repository to a pool.
func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// withTenant runs fn in a transaction whose session is scoped to t.
//
// SET LOCAL and not SET: the scope dies with the transaction, so a pooled
// connection handed to the next request never carries the previous tenant.
func (r *Repository) withTenant(ctx context.Context, t tenant.ID, fn func(pgx.Tx) error) error {
	if t.IsZero() {
		return errors.New("catalog/postgres: tenant is required")
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", t.String()); err != nil {
		return fmt.Errorf("catalog/postgres: scope tenant: %w", err)
	}
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// factColumns renders a Fact[string] as its three columns.
func factColumns(f fact.Fact[string]) (state string, value *string, reason *string) {
	state = f.State().String()
	if v, ok := f.Value(); ok {
		value = &v
	}
	if r := f.Reason(); r != "" {
		reason = &r
	}
	return state, value, reason
}

// Insert writes a new product and everything that hangs off it, atomically.
func (r *Repository) Insert(ctx context.Context, p domain.Product) error {
	return r.withTenant(ctx, p.Tenant(), func(tx pgx.Tx) error {
		return writeProduct(ctx, tx, p, true)
	})
}

// Update writes a new version of an existing product, atomically.
func (r *Repository) Update(ctx context.Context, p domain.Product) error {
	return r.withTenant(ctx, p.Tenant(), func(tx pgx.Tx) error {
		return writeProduct(ctx, tx, p, false)
	})
}

func writeProduct(ctx context.Context, tx pgx.Tx, p domain.Product, insert bool) error {
	state, value, reason := factColumns(p.Description())
	if insert {
		_, err := tx.Exec(ctx, `
			INSERT INTO catalog.products
				(tenant_id, product_id, version, description_state, description_value,
				 description_reason, last_payload_hash)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			p.Tenant().String(), p.ID().String(), p.Version(),
			state, value, reason, p.LastPayloadHash())
		if err != nil {
			return fmt.Errorf("catalog/postgres: insert product: %w", err)
		}
	} else {
		tag, err := tx.Exec(ctx, `
			UPDATE catalog.products
			   SET version = $3, description_state = $4, description_value = $5,
			       description_reason = $6, last_payload_hash = $7, updated_at = now()
			 WHERE tenant_id = $1 AND product_id = $2`,
			p.Tenant().String(), p.ID().String(), p.Version(),
			state, value, reason, p.LastPayloadHash())
		if err != nil {
			return fmt.Errorf("catalog/postgres: update product: %w", err)
		}
		if tag.RowsAffected() != 1 {
			return fmt.Errorf("catalog/postgres: update touched %d rows for %s",
				tag.RowsAffected(), p.ID())
		}
	}

	for _, id := range p.Identifiers() {
		if _, err := tx.Exec(ctx, `
			INSERT INTO catalog.product_identifiers (tenant_id, product_id, kind, value)
			VALUES ($1,$2,$3,$4)
			ON CONFLICT (tenant_id, product_id, kind, value) DO NOTHING`,
			p.Tenant().String(), p.ID().String(), string(id.Kind()), id.Value()); err != nil {
			return fmt.Errorf("catalog/postgres: insert identifier %s: %w", id, err)
		}
	}

	for _, k := range p.SourceKeys() {
		if _, err := tx.Exec(ctx, `
			INSERT INTO catalog.source_product_keys
				(tenant_id, source_system, source_instance, object_kind, external_key, product_id)
			VALUES ($1,$2,$3,$4,$5,$6)
			ON CONFLICT (tenant_id, source_system, source_instance, object_kind, external_key)
			DO UPDATE SET product_id = EXCLUDED.product_id`,
			k.Tenant().String(), k.System(), k.Instance(), k.ObjectKind(), k.ExternalKey(),
			p.ID().String()); err != nil {
			return fmt.Errorf("catalog/postgres: insert source key %s: %w", k, err)
		}
	}
	return nil
}

// BySourceKey resolves a source address to the product it maps to.
func (r *Repository) BySourceKey(ctx context.Context, k contracts.SourceProductKey) (domain.Product, bool, error) {
	var out domain.Product
	var found bool
	err := r.withTenant(ctx, k.Tenant(), func(tx pgx.Tx) error {
		var productID string
		row := tx.QueryRow(ctx, `
			SELECT product_id FROM catalog.source_product_keys
			 WHERE tenant_id = $1 AND source_system = $2 AND source_instance = $3
			   AND object_kind = $4 AND external_key = $5`,
			k.Tenant().String(), k.System(), k.Instance(), k.ObjectKind(), k.ExternalKey())
		switch err := row.Scan(&productID); {
		case errors.Is(err, pgx.ErrNoRows):
			return nil
		case err != nil:
			return fmt.Errorf("catalog/postgres: resolve source key: %w", err)
		}
		p, ok, err := loadProduct(ctx, tx, k.Tenant(), productID)
		if err != nil {
			return err
		}
		out, found = p, ok
		return nil
	})
	return out, found, err
}

// loadProduct rebuilds an aggregate from its rows.
func loadProduct(ctx context.Context, tx pgx.Tx, t tenant.ID, productID string) (domain.Product, bool, error) {
	var version int
	var state, hash string
	var value, reason *string
	var recordedAt time.Time
	row := tx.QueryRow(ctx, `
		SELECT version, description_state, description_value, description_reason,
		       last_payload_hash, updated_at
		  FROM catalog.products WHERE tenant_id = $1 AND product_id = $2`,
		t.String(), productID)
	switch err := row.Scan(&version, &state, &value, &reason, &hash, &recordedAt); {
	case errors.Is(err, pgx.ErrNoRows):
		return domain.Product{}, false, nil
	case err != nil:
		return domain.Product{}, false, fmt.Errorf("catalog/postgres: load product: %w", err)
	}

	// Rehydration goes back through the same constructors that guard the write
	// path. A row that cannot be turned back into a valid aggregate is a defect
	// we want to hear about here, not three layers later.
	id, err := domain.NewProductID(productID)
	if err != nil {
		return domain.Product{}, false, err
	}
	ident, err := loadIdentifiers(ctx, tx, t, productID)
	if err != nil {
		return domain.Product{}, false, err
	}
	keys, err := loadSourceKeys(ctx, tx, t, productID)
	if err != nil {
		return domain.Product{}, false, err
	}
	if len(keys) == 0 {
		return domain.Product{}, false, fmt.Errorf("catalog/postgres: product %s has no source key", productID)
	}

	e, err := provenance.NewEvidence("catalog", "product", productID, recordedAt.UTC(), hash)
	if err != nil {
		return domain.Product{}, false, err
	}
	desc, err := rehydrateFact(state, value, reason, e)
	if err != nil {
		return domain.Product{}, false, err
	}

	obs := contracts.ProductObservation{
		Key: keys[0], Description: desc, Identifiers: ident, Evidence: e,
	}
	p, err := domain.NewProduct(id, obs)
	if err != nil {
		return domain.Product{}, false, err
	}
	for _, k := range keys[1:] {
		extra := obs
		extra.Key = k
		p, _, err = p.Apply(extra)
		if err != nil {
			return domain.Product{}, false, err
		}
	}
	return domain.ReconstituteVersion(p, version, hash), true, nil
}
```

**Nota de desenho, e é uma restrição real:** `loadProduct` reconstrói o agregado pelos construtores, o que é o que queremos, mas `Apply` faz a versão subir. Por isso o domínio precisa de uma porta de reidratação explícita — `domain.ReconstituteVersion(p Product, version int, hash string) Product` — que existe só para o repositório e diz no nome que está a repor um estado persistido e não a decidir um novo. Acrescenta-a a `product.go` da Tarefa 8:

```go
// ReconstituteVersion restores a persisted version and payload hash onto a
// rebuilt aggregate. It exists for the repository alone: rehydration must not
// look like a decision, and Apply must never be reachable with a version the
// database did not choose.
func ReconstituteVersion(p Product, version int, hash string) Product {
	p.version = version
	p.lastHash = hash
	return p
}
```

A evidência reidratada usa `updated_at` como instante de observação, já lido no `SELECT` acima. `time` entra nos imports do ficheiro.

- [ ] **Step 4: Completa os auxiliares do repositório**

Ainda em `repository.go`:

```go
func loadIdentifiers(ctx context.Context, tx pgx.Tx, t tenant.ID, productID string) ([]contracts.Identifier, error) {
	rows, err := tx.Query(ctx, `
		SELECT kind, value FROM catalog.product_identifiers
		 WHERE tenant_id = $1 AND product_id = $2 ORDER BY kind, value`,
		t.String(), productID)
	if err != nil {
		return nil, fmt.Errorf("catalog/postgres: load identifiers: %w", err)
	}
	defer rows.Close()
	var out []contracts.Identifier
	for rows.Next() {
		var kind, value string
		if err := rows.Scan(&kind, &value); err != nil {
			return nil, err
		}
		id, err := contracts.NewIdentifier(contracts.IdentifierKind(kind), value)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func loadSourceKeys(ctx context.Context, tx pgx.Tx, t tenant.ID, productID string) ([]contracts.SourceProductKey, error) {
	rows, err := tx.Query(ctx, `
		SELECT source_system, source_instance, object_kind, external_key
		  FROM catalog.source_product_keys
		 WHERE tenant_id = $1 AND product_id = $2
		 ORDER BY source_system, source_instance, object_kind, external_key`,
		t.String(), productID)
	if err != nil {
		return nil, fmt.Errorf("catalog/postgres: load source keys: %w", err)
	}
	defer rows.Close()
	var out []contracts.SourceProductKey
	for rows.Next() {
		var system, instance, kind, key string
		if err := rows.Scan(&system, &instance, &kind, &key); err != nil {
			return nil, err
		}
		k, err := contracts.NewSourceProductKey(t, system, instance, kind, key)
		if err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

// rehydrateFact turns three columns back into a Fact, through the same
// constructors the write path uses.
func rehydrateFact(state string, value, reason *string, e provenance.Evidence) (fact.Fact[string], error) {
	text := ""
	if value != nil {
		text = *value
	}
	why := ""
	if reason != nil {
		why = *reason
	}
	switch state {
	case fact.Known.String():
		return fact.NewKnown(text, e)
	case fact.Estimated.String():
		return fact.NewEstimated(text, why, e)
	case fact.NotApplicable.String():
		return fact.NewNotApplicable[string](why, e)
	default:
		return fact.NewUnknown[string](why, e)
	}
}

// ByIdentifier returns every product of this tenant carrying the identifier.
func (r *Repository) ByIdentifier(ctx context.Context, t tenant.ID, id contracts.Identifier) ([]domain.Product, error) {
	var out []domain.Product
	err := r.withTenant(ctx, t, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT product_id FROM catalog.product_identifiers
			 WHERE tenant_id = $1 AND kind = $2 AND value = $3 ORDER BY product_id`,
			t.String(), string(id.Kind()), id.Value())
		if err != nil {
			return fmt.Errorf("catalog/postgres: lookup identifier: %w", err)
		}
		var ids []string
		for rows.Next() {
			var productID string
			if err := rows.Scan(&productID); err != nil {
				rows.Close()
				return err
			}
			ids = append(ids, productID)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return err
		}
		for _, productID := range ids {
			p, ok, err := loadProduct(ctx, tx, t, productID)
			if err != nil {
				return err
			}
			if ok {
				out = append(out, p)
			}
		}
		return nil
	})
	return out, err
}

// ByProductID reads one product as a summary. SummaryReader forwards to it.
func (r *Repository) ByProductID(ctx context.Context, t tenant.ID, productID string) (port.Summary, bool, error) {
	var out port.Summary
	var found bool
	err := r.withTenant(ctx, t, func(tx pgx.Tx) error {
		p, ok, err := loadProduct(ctx, tx, t, productID)
		if err != nil || !ok {
			return err
		}
		out, found = summarise(p), true
		return nil
	})
	return out, found, err
}

// summarise flattens an aggregate for a consumer. An Unknown description
// becomes the empty string ONLY here, at the edge, and the caller is told the
// product exists rather than being handed a fabricated name.
func summarise(p domain.Product) port.Summary {
	desc, _ := p.Description().Value()
	return port.Summary{
		ProductID:   p.ID().String(),
		Description: desc,
		Identifiers: p.Identifiers(),
		Version:     p.Version(),
	}
}

// ULIDFactory mints canonical identifiers. The value is random and carries no
// source semantics whatsoever, which is the whole property.
type ULIDFactory struct{}

// NewULIDFactory builds the factory.
func NewULIDFactory() *ULIDFactory { return &ULIDFactory{} }

// NewProductID mints prd_ + 32 hex characters from crypto/rand.
func (f *ULIDFactory) NewProductID() (domain.ProductID, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return domain.ProductID{}, fmt.Errorf("catalog: mint id: %w", err)
	}
	return domain.NewProductID(domain.ProductIDPrefix + hex.EncodeToString(b[:]))
}
```

`ByIdentifier` do `port.Reader` recebe o mesmo nome do da `Store` mas devolve `[]port.Summary`. Como um tipo não pode ter dois métodos com o mesmo nome, o `Repository` implementa a versão da `Store` (`[]domain.Product`) e o `Module` da Tarefa 9 recebe um adaptador fino. Acrescenta ao fim de `repository.go`:

```go
// SummaryReader adapts the repository to port.Reader without a second method of
// the same name on the same type.
type SummaryReader struct{ repo *Repository }

// NewSummaryReader wraps the repository as the context's published reader.
func NewSummaryReader(r *Repository) *SummaryReader { return &SummaryReader{repo: r} }

// ByProductID implements port.Reader.
func (s *SummaryReader) ByProductID(ctx context.Context, t tenant.ID, productID string) (port.Summary, bool, error) {
	return s.repo.ByProductID(ctx, t, productID)
}

// ByIdentifier implements port.Reader.
func (s *SummaryReader) ByIdentifier(ctx context.Context, t tenant.ID, id contracts.Identifier) ([]port.Summary, error) {
	products, err := s.repo.ByIdentifier(ctx, t, id)
	if err != nil {
		return nil, err
	}
	out := make([]port.Summary, 0, len(products))
	for _, p := range products {
		out = append(out, summarise(p))
	}
	return out, nil
}
```

Fecha o ficheiro com as asserções de compilação, que é onde uma porta em falta é apanhada em segundos em vez de no arranque:

```go
var (
	_ application.Store = (*Repository)(nil)
	_ application.IDFactory = (*ULIDFactory)(nil)
	_ port.Reader = (*SummaryReader)(nil)
)
```

`application` entra nos imports.

- [ ] **Step 5: Compila**

```bash
cd apps/server_core && go build ./internal/contexts/... && go vet ./internal/contexts/...
```

Esperado: exit 0. Se o compilador reclamar de `time` em falta, acrescenta o import.

- [ ] **Step 6: Escreve o teste de integração**

`apps/server_core/tests/integration/catalog_ingest_test.go` — **as duas primeiras linhas têm de ser exatamente** a tag e a linha em branco, senão a lane não o vê:

```go
//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/catalog"
	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	catalogpostgres "marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func catalogObservation(t *testing.T, tid tenant.ID, externalKey, description, ean, hash string) contracts.ProductObservation {
	t.Helper()
	e, err := provenance.NewEvidence("sankhya", "product", externalKey, time.Now().UTC(), hash)
	if err != nil {
		t.Fatalf("provenance.NewEvidence: %v", err)
	}
	desc, err := fact.NewKnown(description, e)
	if err != nil {
		t.Fatalf("fact.NewKnown: %v", err)
	}
	key, err := contracts.NewSourceProductKey(tid, "sankhya", "sankhya-it-01", "product", externalKey)
	if err != nil {
		t.Fatalf("contracts.NewSourceProductKey: %v", err)
	}
	id, err := contracts.NewIdentifier(contracts.IdentifierEAN, ean)
	if err != nil {
		t.Fatalf("contracts.NewIdentifier: %v", err)
	}
	return contracts.ProductObservation{Key: key, Description: desc,
		Identifiers: []contracts.Identifier{id}, Evidence: e}
}

func TestCatalogIngestPersistsAndIsIdempotent(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_catalog_ingest")
	tid, err := tenant.Parse(cfg.DefaultTenantID)
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	repo := catalogpostgres.NewRepository(pool)
	module := catalog.New(repo, catalogpostgres.NewULIDFactory(), catalogpostgres.NewSummaryReader(repo))

	ctx := context.Background()
	stamp := time.Now().UTC().UnixNano()
	externalKey := "it-" + time.Now().UTC().Format("150405.000000000")

	first, err := module.IngestProduct(ctx,
		catalogObservation(t, tid, externalKey, "Cafeteira Eletrica", "789ean"+itoa(stamp), "sha256:one"+itoa(stamp)))
	if err != nil {
		t.Fatalf("first IngestProduct: %v", err)
	}
	if first.Disposition != contracts.DispositionCreated || first.Version != 1 {
		t.Fatalf("first = %+v, want created/1", first)
	}

	// Same payload hash: nothing moves, and the id does not.
	second, err := module.IngestProduct(ctx,
		catalogObservation(t, tid, externalKey, "Cafeteira Eletrica", "789ean"+itoa(stamp), "sha256:one"+itoa(stamp)))
	if err != nil {
		t.Fatalf("second IngestProduct: %v", err)
	}
	if second.Disposition != contracts.DispositionIdempotent {
		t.Fatalf("second Disposition = %q, want idempotent", second.Disposition)
	}
	if second.ProductID != first.ProductID {
		t.Fatalf("ProductID moved: %q then %q", first.ProductID, second.ProductID)
	}

	// Changed payload: version 2, same identity.
	third, err := module.IngestProduct(ctx,
		catalogObservation(t, tid, externalKey, "Cafeteira Eletrica Inox", "789ean"+itoa(stamp), "sha256:two"+itoa(stamp)))
	if err != nil {
		t.Fatalf("third IngestProduct: %v", err)
	}
	if third.Version != 2 || third.ProductID != first.ProductID {
		t.Fatalf("third = %+v, want version 2 on %s", third, first.ProductID)
	}

	got, found, err := module.Reader().ByProductID(ctx, tid, first.ProductID)
	if err != nil {
		t.Fatalf("ByProductID: %v", err)
	}
	if !found {
		t.Fatalf("ByProductID did not find %s", first.ProductID)
	}
	if got.Description != "Cafeteira Eletrica Inox" || got.Version != 2 {
		t.Fatalf("summary = %+v, want the version 2 description", got)
	}
}

// The RLS proof. Reading with a different tenant scope must return nothing —
// not filtered by the application, filtered by the database.
func TestCatalogIsInvisibleToAnotherTenant(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_catalog_rls")
	tid, err := tenant.Parse(cfg.DefaultTenantID)
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	repo := catalogpostgres.NewRepository(pool)
	module := catalog.New(repo, catalogpostgres.NewULIDFactory(), catalogpostgres.NewSummaryReader(repo))

	ctx := context.Background()
	stamp := time.Now().UTC().UnixNano()
	created, err := module.IngestProduct(ctx,
		catalogObservation(t, tid, "rls-"+itoa(stamp), "Cafeteira", "789rls"+itoa(stamp), "sha256:rls"+itoa(stamp)))
	if err != nil {
		t.Fatalf("IngestProduct: %v", err)
	}

	intruder, err := tenant.Parse("tnt_intruder")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	_, found, err := module.Reader().ByProductID(ctx, intruder, created.ProductID)
	if err != nil {
		t.Fatalf("cross-tenant ByProductID: %v", err)
	}
	if found {
		t.Fatalf("tenant %s read product %s belonging to %s", intruder, created.ProductID, tid)
	}
}

// itoa keeps the fixture strings unique per run, so a re-run does not collide
// with the rows the previous run left behind.
func itoa(n int64) string { return strconv.FormatInt(n, 10) }
```

`strconv` entra nos imports do ficheiro.

- [ ] **Step 7: Corre a lane hermética**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test -tags integration ./tests/integration/ -run TestCatalog -v
```

Esperado: dois testes PASS. **Confirma que aparecem os nomes dos dois testes no output** — `ok` sem nomes é indistinguível de pulado.

Se `TestCatalogIsInvisibleToAnotherTenant` falhar com `found = true`, a RLS não está ativa: verifica se o papel da ligação é dono das tabelas e se o `FORCE` foi aplicado. Não relaxes o teste.

- [ ] **Step 8: Commit**

```bash
git status --porcelain --untracked-files=all
git add apps/server_core/migrations/0097_catalog_context.sql \
        apps/server_core/internal/contexts/catalog \
        apps/server_core/tests/integration/catalog_ingest_test.go
git commit -m "feat(catalog): esquema proprio, escritor unico e RLS FORCE provada por leitura cruzada"
```

---

## Tarefa 11: o adaptador Sankhya e a prova da fronteira na nossa árvore

O adaptador é o único sítio onde `CODPROD`, `TGFPRO` e `METALPRD` podem aparecer. A prova de que isso é imposto pelo compilador — e não por revisão — faz-se aqui, injetando a violação e colando a rejeição.

`internal/oracle` usa `database/sql` e mais nada. **Não importa `godror`**: o driver regista-se no composition root. Assim este pacote compila sem `cgo` e o teste unitário corre no host sem falso verde.

**Files:**
- Create: `apps/server_core/internal/adapters/erp/sankhyaoracle/sankhyaoracle.go`
- Create: `apps/server_core/internal/adapters/erp/sankhyaoracle/internal/oracle/rows.go`
- Create: `apps/server_core/internal/adapters/erp/sankhyaoracle/catalogfeed/mapper.go`
- Create: `apps/server_core/internal/adapters/erp/sankhyaoracle/catalogfeed/mapper_test.go`
- Create: `apps/server_core/internal/composition/catalog_wiring.go`

**Interfaces:**
- Produces: `sankhyaoracle.New(db *sql.DB, instance string) (Bundle, error)` com `Bundle{Instance string; CatalogFeed catalogfeed.Feed}`.
- Consumes: `contracts.ProductObservation` (T7), `kernel/fact`, `kernel/provenance`, `kernel/tenant`.

- [ ] **Step 1: Escreve a linha crua, dentro de `internal/`**

`internal/adapters/erp/sankhyaoracle/internal/oracle/rows.go`:

```go
// Package oracle holds the wire shape of Sankhya's tables. Nothing outside
// adapters/erp/sankhyaoracle can import it — the compiler says so, not a
// reviewer — which is why CODPROD may exist here and nowhere else.
package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ProductRow is TGFPRO as it comes off the wire, nulls and all.
//
// Every optional column is a sql.Null*, because "the column was NULL" and "the
// column was zero" are different facts and collapsing them here would make the
// lie unrecoverable downstream.
type ProductRow struct {
	Codprod    int64
	Descrprod  sql.NullString
	Referencia sql.NullString
	Marca      sql.NullString
	NCM        sql.NullString
	Ativo      sql.NullString
	ReadAt     time.Time
}

// SelectActiveProducts is the only statement this package issues.
const SelectActiveProducts = `
SELECT p.CODPROD, p.DESCRPROD, p.REFERENCIA, p.MARCA, p.NCM, p.ATIVO
  FROM METALPRD.TGFPRO p
 WHERE p.ATIVO = 'S'
   AND p.CODPROD > :1
 ORDER BY p.CODPROD
 FETCH FIRST :2 ROWS ONLY`

// FetchActiveProducts pages TGFPRO by CODPROD.
//
// The cursor is the last CODPROD seen, not an offset: OFFSET re-reads the whole
// prefix and silently skips rows when the table changes under the paging.
func FetchActiveProducts(ctx context.Context, db *sql.DB, afterCodprod int64, limit int, now time.Time) ([]ProductRow, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("sankhyaoracle: limit must be positive, got %d", limit)
	}
	rows, err := db.QueryContext(ctx, SelectActiveProducts, afterCodprod, limit)
	if err != nil {
		return nil, fmt.Errorf("sankhyaoracle: query TGFPRO: %w", err)
	}
	defer rows.Close()

	out := make([]ProductRow, 0, limit)
	for rows.Next() {
		var r ProductRow
		if err := rows.Scan(&r.Codprod, &r.Descrprod, &r.Referencia, &r.Marca, &r.NCM, &r.Ativo); err != nil {
			return nil, fmt.Errorf("sankhyaoracle: scan TGFPRO: %w", err)
		}
		r.ReadAt = now.UTC()
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sankhyaoracle: iterate TGFPRO: %w", err)
	}
	return out, nil
}
```

- [ ] **Step 2: Escreve o teste do mapeador, antes do mapeador**

`internal/adapters/erp/sankhyaoracle/catalogfeed/mapper_test.go`:

```go
package catalogfeed

import (
	"database/sql"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/adapters/erp/sankhyaoracle/internal/oracle"
	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func testTenant(t *testing.T) tenant.ID {
	t.Helper()
	id, err := tenant.Parse("tnt_metal")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	return id
}

func TestMapProductCarriesDescriptionAsKnown(t *testing.T) {
	row := oracle.ProductRow{
		Codprod:    4711,
		Descrprod:  sql.NullString{String: "Cafeteira Eletrica", Valid: true},
		Referencia: sql.NullString{String: "7891234567895", Valid: true},
		ReadAt:     time.Date(2026, 8, 6, 3, 0, 0, 0, time.UTC),
	}
	obs, err := MapProduct(testTenant(t), "sankhya-it-01", row)
	if err != nil {
		t.Fatalf("MapProduct: %v", err)
	}
	if obs.Description.State() != fact.Known {
		t.Fatalf("Description state = %v, want Known", obs.Description.State())
	}
	if v, _ := obs.Description.Value(); v != "Cafeteira Eletrica" {
		t.Fatalf("Description = %q", v)
	}
	if obs.Key.ExternalKey() != "4711" {
		t.Fatalf("ExternalKey = %q, want 4711", obs.Key.ExternalKey())
	}
	if len(obs.Identifiers) != 1 || obs.Identifiers[0].Value() != "7891234567895" {
		t.Fatalf("Identifiers = %v", obs.Identifiers)
	}
}

// The whole point of the kernel: a NULL column must not arrive downstream as "".
func TestMapProductTurnsNullDescriptionIntoUnknownNotEmptyString(t *testing.T) {
	row := oracle.ProductRow{
		Codprod:   4712,
		Descrprod: sql.NullString{},
		ReadAt:    time.Date(2026, 8, 6, 3, 0, 0, 0, time.UTC),
	}
	obs, err := MapProduct(testTenant(t), "sankhya-it-01", row)
	if err != nil {
		t.Fatalf("MapProduct: %v", err)
	}
	if obs.Description.State() != fact.Unknown {
		t.Fatalf("Description state = %v, want Unknown", obs.Description.State())
	}
	if _, ok := obs.Description.Value(); ok {
		t.Fatal("Unknown description handed out a value")
	}
	if obs.Description.Reason() == "" {
		t.Fatal("Unknown description carries no reason")
	}
}

// A blank REFERENCIA is not an identifier. Emitting "" would make every
// product with a blank EAN collide with every other one.
func TestMapProductDropsBlankIdentifier(t *testing.T) {
	row := oracle.ProductRow{
		Codprod:    4713,
		Descrprod:  sql.NullString{String: "Chaleira", Valid: true},
		Referencia: sql.NullString{String: "   ", Valid: true},
		ReadAt:     time.Date(2026, 8, 6, 3, 0, 0, 0, time.UTC),
	}
	obs, err := MapProduct(testTenant(t), "sankhya-it-01", row)
	if err != nil {
		t.Fatalf("MapProduct: %v", err)
	}
	if len(obs.Identifiers) != 0 {
		t.Fatalf("Identifiers = %v, want none", obs.Identifiers)
	}
}

// Two reads of the same unchanged row must hash the same, or every sync run
// bumps every version and idempotence means nothing.
func TestMapProductHashIsStableAcrossReadTimes(t *testing.T) {
	base := oracle.ProductRow{
		Codprod:   4714,
		Descrprod: sql.NullString{String: "Torradeira", Valid: true},
		ReadAt:    time.Date(2026, 8, 6, 3, 0, 0, 0, time.UTC),
	}
	later := base
	later.ReadAt = base.ReadAt.Add(6 * time.Hour)

	a, err := MapProduct(testTenant(t), "sankhya-it-01", base)
	if err != nil {
		t.Fatalf("MapProduct: %v", err)
	}
	b, err := MapProduct(testTenant(t), "sankhya-it-01", later)
	if err != nil {
		t.Fatalf("MapProduct: %v", err)
	}
	if a.Evidence.PayloadHash() != b.Evidence.PayloadHash() {
		t.Fatalf("hash moved with read time: %q vs %q",
			a.Evidence.PayloadHash(), b.Evidence.PayloadHash())
	}
	if a.Evidence.ObservedAt().Equal(b.Evidence.ObservedAt()) {
		t.Fatal("ObservedAt did not move, so the test proves nothing")
	}
}

func TestMapProductRejectsZeroCodprod(t *testing.T) {
	_, err := MapProduct(testTenant(t), "sankhya-it-01", oracle.ProductRow{Codprod: 0})
	if err == nil {
		t.Fatal("MapProduct accepted CODPROD 0")
	}
}

var _ = contracts.ProductObservation{}
```

- [ ] **Step 3: Corre e vê falhar**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/adapters/erp/sankhyaoracle/...
```

Esperado: FAIL de compilação, `undefined: MapProduct`. Cola o output.

- [ ] **Step 4: Escreve o mapeador**

`internal/adapters/erp/sankhyaoracle/catalogfeed/mapper.go`:

```go
// Package catalogfeed translates Sankhya rows into Catalog observations. It is
// the only place in the platform that knows both vocabularies at once.
package catalogfeed

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/adapters/erp/sankhyaoracle/internal/oracle"
	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// SourceSystem names Sankhya in every SourceProductKey this package mints.
const SourceSystem = "sankhya"

// ObjectKind names the TGFPRO row kind.
const ObjectKind = "product"

// Feed pages products out of Sankhya as Catalog observations.
type Feed struct {
	db       *sql.DB
	instance string
	now      func() time.Time
}

// NewFeed builds the feed. instance identifies WHICH Sankhya installation this
// is, because two installations both call their first product CODPROD 1.
func NewFeed(db *sql.DB, instance string, now func() time.Time) (Feed, error) {
	if db == nil {
		return Feed{}, fmt.Errorf("catalogfeed: db is required")
	}
	if strings.TrimSpace(instance) == "" {
		return Feed{}, fmt.Errorf("catalogfeed: instance is required")
	}
	if now == nil {
		now = time.Now
	}
	return Feed{db: db, instance: strings.TrimSpace(instance), now: now}, nil
}

// Page returns up to limit observations with CODPROD greater than after, and
// the cursor to pass next. A zero cursor means the page was the last one.
func (f Feed) Page(ctx context.Context, t tenant.ID, after int64, limit int) ([]contracts.ProductObservation, int64, error) {
	rows, err := oracle.FetchActiveProducts(ctx, f.db, after, limit, f.now())
	if err != nil {
		return nil, 0, err
	}
	out := make([]contracts.ProductObservation, 0, len(rows))
	var next int64
	for _, r := range rows {
		obs, err := MapProduct(t, f.instance, r)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, obs)
		next = r.Codprod
	}
	if len(rows) < limit {
		next = 0
	}
	return out, next, nil
}

// MapProduct turns one TGFPRO row into one observation.
func MapProduct(t tenant.ID, instance string, row oracle.ProductRow) (contracts.ProductObservation, error) {
	if row.Codprod <= 0 {
		return contracts.ProductObservation{},
			fmt.Errorf("catalogfeed: CODPROD must be positive, got %d", row.Codprod)
	}
	externalKey := strconv.FormatInt(row.Codprod, 10)

	key, err := contracts.NewSourceProductKey(t, SourceSystem, instance, ObjectKind, externalKey)
	if err != nil {
		return contracts.ProductObservation{}, err
	}

	hash := payloadHash(row)
	e, err := provenance.NewEvidence(SourceSystem, ObjectKind, externalKey, row.ReadAt.UTC(), hash)
	if err != nil {
		return contracts.ProductObservation{}, err
	}

	desc, err := mapDescription(row.Descrprod, e)
	if err != nil {
		return contracts.ProductObservation{}, err
	}

	var identifiers []contracts.Identifier
	if ean := strings.TrimSpace(row.Referencia.String); row.Referencia.Valid && ean != "" {
		id, err := contracts.NewIdentifier(contracts.IdentifierEAN, ean)
		if err != nil {
			return contracts.ProductObservation{}, err
		}
		identifiers = append(identifiers, id)
	}

	obs := contracts.ProductObservation{
		Key: key, Description: desc, Identifiers: identifiers, Evidence: e,
	}
	if err := obs.Validate(); err != nil {
		return contracts.ProductObservation{}, err
	}
	return obs, nil
}

// mapDescription is where ADR-017 stops being a document. A NULL column becomes
// Unknown with a reason, never the empty string.
func mapDescription(col sql.NullString, e provenance.Evidence) (fact.Fact[string], error) {
	if !col.Valid {
		return fact.NewUnknown[string]("TGFPRO.DESCRPROD is NULL", e)
	}
	if v := strings.TrimSpace(col.String); v != "" {
		return fact.NewKnown(v, e)
	}
	return fact.NewUnknown[string]("TGFPRO.DESCRPROD is blank", e)
}

// payloadHash digests the business columns and NOTHING else.
//
// ReadAt is deliberately excluded: including it would make every sync run
// produce a new hash, which would defeat idempotence and bump every version on
// every pass.
func payloadHash(row oracle.ProductRow) string {
	h := sha256.New()
	fmt.Fprintf(h, "codprod=%d\n", row.Codprod)
	for _, f := range []struct {
		name string
		col  sql.NullString
	}{
		{"descrprod", row.Descrprod},
		{"referencia", row.Referencia},
		{"marca", row.Marca},
		{"ncm", row.NCM},
		{"ativo", row.Ativo},
	} {
		if f.col.Valid {
			fmt.Fprintf(h, "%s=%s\n", f.name, f.col.String)
		} else {
			// A distinct marker: NULL must not hash the same as "".
			fmt.Fprintf(h, "%s=<null>\n", f.name)
		}
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}
```

- [ ] **Step 5: Corre e vê passar**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/adapters/erp/sankhyaoracle/... -v
```

Esperado: os cinco testes PASS, pelo nome.

- [ ] **Step 6: Escreve a fachada**

`internal/adapters/erp/sankhyaoracle/sankhyaoracle.go`:

```go
// Package sankhyaoracle is the only importable surface of the Sankhya adapter.
// Wire rows and SQL live under internal/oracle and cannot be reached from any
// other package in the module, including the composition root.
//
// The driver is NOT imported here. Registering godror is the composition
// root's job; this package speaks database/sql and therefore builds without
// cgo, which is what keeps its tests from passing vacuously on a host with no
// Oracle client.
package sankhyaoracle

import (
	"database/sql"
	"time"

	"marketplace-central/apps/server_core/internal/adapters/erp/sankhyaoracle/catalogfeed"
)

// Bundle is everything this adapter offers.
type Bundle struct {
	Instance    string
	CatalogFeed catalogfeed.Feed
}

// New builds the bundle for one Sankhya installation.
func New(db *sql.DB, instance string, now func() time.Time) (Bundle, error) {
	feed, err := catalogfeed.NewFeed(db, instance, now)
	if err != nil {
		return Bundle{}, err
	}
	return Bundle{Instance: instance, CatalogFeed: feed}, nil
}
```

- [ ] **Step 7: A prova da fronteira, na nossa árvore**

Não basta ter provado no esqueleto. Prova aqui, agora, e cola os dois outputs.

Cria o ficheiro de violação:

```bash
cd apps/server_core && mkdir -p internal/contexts/catalog/tmpprobe && cat > internal/contexts/catalog/tmpprobe/probe.go <<'EOF'
package tmpprobe

import "marketplace-central/apps/server_core/internal/adapters/erp/sankhyaoracle/internal/oracle"

var _ = oracle.ProductRow{}
EOF
GOCACHE="$(pwd)/.gocache" go build ./internal/contexts/catalog/tmpprobe/
```

Esperado: **exit diferente de zero**, com uma mensagem da forma:

```
use of internal package marketplace-central/apps/server_core/internal/adapters/erp/sankhyaoracle/internal/oracle not allowed
```

Cola o output verbatim no relatório. **Se compilar, para tudo**: a fronteira não existe, e todo o desenho da §2 do protocolo está errado. Escreve isso em `.mnfs/HARNESS-DEBTS.md` como bloqueio de classe e não continues a Tarefa 12.

Remove a sonda:

```bash
cd apps/server_core && rm -rf internal/contexts/catalog/tmpprobe && git status --porcelain --untracked-files=all
```

Esperado: sem qualquer linha `tmpprobe`.

- [ ] **Step 8: Liga no composition root**

`internal/composition/catalog_wiring.go` — ficheiro **novo**, nada em `internal/modules/` é tocado:

```go
package composition

import (
	"database/sql"
	"time"

	"marketplace-central/apps/server_core/internal/adapters/erp/sankhyaoracle"
	"marketplace-central/apps/server_core/internal/contexts/catalog"
	catalogpostgres "marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CatalogWiring is the assembled catalog slice.
//
// Note what this struct CANNOT say: there is no field typed by anything under
// an adapter's internal/, because no such type can be named here. That is the
// property, and it is enforced by the compiler rather than by this comment.
type CatalogWiring struct {
	Module *catalog.Module
	Feed   sankhyaoracle.Bundle
}

// WireCatalog assembles the slice from its two edges.
func WireCatalog(pool *pgxpool.Pool, oracleDB *sql.DB, sankhyaInstance string) (CatalogWiring, error) {
	repo := catalogpostgres.NewRepository(pool)
	module := catalog.New(repo, catalogpostgres.NewULIDFactory(), catalogpostgres.NewSummaryReader(repo))

	bundle, err := sankhyaoracle.New(oracleDB, sankhyaInstance, time.Now)
	if err != nil {
		return CatalogWiring{}, err
	}
	return CatalogWiring{Module: module, Feed: bundle}, nil
}
```

- [ ] **Step 9: Compila e vet a árvore toda**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go build ./... && GOCACHE="$(pwd)/.gocache" go vet ./...
```

Esperado: exit 0 em ambos. `WireCatalog` não é chamado por ninguém ainda — é intencional; o `go vet` não reclama de função exportada por usar.

- [ ] **Step 10: Commit**

```bash
git status --porcelain --untracked-files=all
git add apps/server_core/internal/adapters apps/server_core/internal/composition/catalog_wiring.go
git commit -m "feat(sankhya): adaptador com fachada e fronteira provada pelo compilador"
```

---

## Tarefa 12: o portão que corre sozinho

Sem portão, tudo o que este plano construiu é convenção outra vez. O portão corre os detetores da Tarefa 6, o `gofmt`, o `vet` e os testes, e o que ele mede tem de ser exatamente o que o protocolo diz.

**Files:**
- Create: `scripts/arch-gate.sh`
- Create: `apps/server_core/internal/arch/cmd/archscan/main.go`
- Modify: `apps/server_core/internal/arch/repo_test.go` (remove os `t.Skip`)

**Interfaces:**
- Consumes: `arch.ScanCrossContextInternal`, `arch.ScanFloatInContracts`, `arch.ScanVendorTokens` (T6).
- Produces: `scripts/arch-gate.sh`, exit 0 limpo / não-zero com `file:line` por achado.

- [ ] **Step 1: Escreve o executável do scanner**

O portão é um script de shell mas a lógica é Go, porque os detetores já existem em Go e reescrevê-los em `grep` seria a terceira cópia da mesma regra.

`apps/server_core/internal/arch/cmd/archscan/main.go`:

```go
// Command archscan runs the architecture detectors over a tree and prints every
// finding as file:line. Exit 0 means zero findings.
package main

import (
	"flag"
	"fmt"
	"os"

	"marketplace-central/apps/server_core/internal/arch"
)

func main() {
	root := flag.String("root", "internal", "directory to scan")
	flag.Parse()

	var findings []arch.Finding
	for _, scan := range []func(string) ([]arch.Finding, error){
		arch.ScanCrossContextInternal,
		arch.ScanFloatInContracts,
		arch.ScanVendorTokens,
	} {
		got, err := scan(*root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "archscan: %v\n", err)
			os.Exit(2)
		}
		findings = append(findings, got...)
	}

	for _, f := range findings {
		fmt.Printf("%s:%d: %s: %s\n", f.File, f.Line, f.Rule, f.Detail)
	}
	if len(findings) > 0 {
		fmt.Fprintf(os.Stderr, "archscan: %d finding(s)\n", len(findings))
		os.Exit(1)
	}
}
```

- [ ] **Step 2: Prova que o scanner dispara**

Um scanner que só alguma vez imprimiu zero é indistinguível de um scanner que não olha para nada. Corre-o contra o diretório de fixtures da Tarefa 6:

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go run ./internal/arch/cmd/archscan -root internal/arch/testdata; echo "EXIT=$?"
```

Esperado: **`EXIT=1`** e pelo menos uma linha `file:line: <regra>: <detalhe>` por cada regra. Cola o output.

Se der `EXIT=0`, o scanner está cego — os fixtures da Tarefa 6 terminam em `.go.txt` e o `archscan` só lê `.go`. Nesse caso o `-root` de fixtures não serve, e o que provas em vez disso é o teste da Tarefa 6 (`go test ./internal/arch/ -run TestScan -v`), que carrega os fixtures pelo caminho certo. Regista qual dos dois usaste; não declares a prova sem um output.

- [ ] **Step 3: Remove os `t.Skip` de `repo_test.go`**

`internal/contexts` já existe desde a Tarefa 7. Os três `t.Skip` da Tarefa 6 têm de sair — um skip permanente é um teste apagado que ainda aparece como verde.

```bash
cd apps/server_core && grep -n "t.Skip" internal/arch/repo_test.go
```

Apaga cada linha `t.Skip(...)` e a condição que a guarda. Depois:

```bash
cd apps/server_core && grep -c "t.Skip" internal/arch/repo_test.go
```

Esperado: `0`.

- [ ] **Step 4: Corre os testes de arquitetura**

```bash
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/arch/... -v
```

Esperado: todos PASS, **com os nomes visíveis** e sem nenhum `SKIP`.

- [ ] **Step 5: Escreve o portão**

`scripts/arch-gate.sh`:

```bash
#!/usr/bin/env bash
# Architecture gate. Run from the repository root.
#
# Exit 0 means: formatted, vetted, tested, and zero architecture findings.
# Every check prints what it measured, because a gate that only prints "ok" is
# indistinguishable from a gate that skipped.
set -euo pipefail

cd "$(dirname "$0")/.."
SERVER="apps/server_core"
export GOCACHE="$PWD/$SERVER/.gocache"

fail=0
step() { printf '\n== %s ==\n' "$1"; }

step "gofmt"
unformatted="$(gofmt -l "$SERVER/internal" || true)"
if [ -n "$unformatted" ]; then
  echo "$unformatted"
  echo "gofmt: files above are not formatted"
  fail=1
else
  echo "gofmt: clean"
fi

step "go vet"
if ! (cd "$SERVER" && go vet ./...); then fail=1; fi

step "architecture detectors"
if (cd "$SERVER" && go run ./internal/arch/cmd/archscan -root internal); then
  echo "archscan: zero findings"
else
  fail=1
fi

step "no float in the kernel"
floats="$(grep -rn 'float64\|float32' "$SERVER/internal/kernel" || true)"
if [ -n "$floats" ]; then
  echo "$floats"
  echo "kernel: float is forbidden"
  fail=1
else
  echo "kernel: no float"
fi

step "unit tests"
if ! (cd "$SERVER" && go test ./internal/...); then fail=1; fi

step "working tree"
dirty="$(git status --porcelain --untracked-files=all)"
if [ -n "$dirty" ]; then
  echo "$dirty"
  echo "tree: dirty — a gate cannot certify a tree it did not see"
  fail=1
else
  echo "tree: clean"
fi

printf '\n'
if [ "$fail" -ne 0 ]; then
  echo "ARCH GATE: FAIL"
  exit 1
fi
echo "ARCH GATE: PASS"
```

```bash
chmod +x scripts/arch-gate.sh
```

- [ ] **Step 6: Prova que o portão reprova**

Um portão que nunca reprovou não é um portão. Injeta uma violação, corre, confirma o `FAIL`, e desfaz **por edição, nunca por `git checkout`/`reset`**:

```bash
cd apps/server_core && printf '\npackage kernelprobe\n\nvar Broken float64\n' > internal/kernel/exact/probe_float.go
cd "$(git rev-parse --show-toplevel)" && ./scripts/arch-gate.sh; echo "EXIT=$?"
```

Esperado: `ARCH GATE: FAIL` e `EXIT=1`, com a linha do `probe_float.go` no bloco `no float in the kernel`. Cola o output.

```bash
rm apps/server_core/internal/kernel/exact/probe_float.go
./scripts/arch-gate.sh; echo "EXIT=$?"
```

Esperado: `ARCH GATE: PASS` e `EXIT=0`.

- [ ] **Step 7: Commit**

```bash
git status --porcelain --untracked-files=all
git add scripts/arch-gate.sh apps/server_core/internal/arch
git commit -m "feat(arch): portao de arquitetura provado a reprovar e a aprovar"
```

- [ ] **Step 8: A medição de fecho**

Escreve `.mnfs/MEASUREMENTS/2026-08-06-fundacao-kernel.md` com estes números, cada um obtido pelo comando ao lado, nenhum estimado:

```bash
cd apps/server_core
echo "kernel packages:      $(ls internal/kernel | wc -l)"
echo "context packages:     $(ls internal/contexts | wc -l)"
echo "float in kernel:      $(grep -rn 'float64\|float32' internal/kernel | wc -l)"
echo "float in contracts:   $(grep -rn 'float64\|float32' internal/contexts/*/contracts | wc -l)"
echo "modules touched:      $(git diff --name-only HEAD~12..HEAD -- internal/modules | wc -l)"
echo "new dependencies:     $(git diff HEAD~12..HEAD -- go.mod | grep -c '^+\s' || true)"
```

Esperado, e cada desvio é um achado a registar em `.mnfs/HARNESS-DEBTS.md`:

| medida | esperado |
|---|---|
| kernel packages | 6 |
| context packages | 1 |
| float in kernel | 0 |
| float in contracts | 0 |
| modules touched | 0 |
| new dependencies | 0 |

E corre o portão uma última vez, do zero:

```bash
./scripts/arch-gate.sh; echo "EXIT=$?"
```

Esperado: `ARCH GATE: PASS`, `EXIT=0`.

```bash
git add .mnfs/MEASUREMENTS/2026-08-06-fundacao-kernel.md
git commit -m "docs(arch): medicao de fecho da fundacao — 6 pacotes de kernel, 0 float, 0 modulos tocados"
```

---

## O que este plano NÃO faz

Escrito por extenso porque um plano que não diz onde acaba convida a que se lhe acrescente âmbito às três da manhã.

- **Não apaga nem altera nada em `internal/modules/`.** As 128 violações de fronteira medidas em `9555a96c` continuam lá depois deste plano. Este plano constrói o sítio para onde elas vão migrar e o instrumento que impede novas; a migração é a Onda 2.
- **Não migra dados.** Nenhum `INSERT ... SELECT` de tabela antiga para o esquema `catalog`.
- **Não expõe HTTP.** Não há rota, não há handler, não há OpenAPI, não há SDK. Uma superfície publicada obriga a atualizar OpenAPI e `sdk-runtime` no mesmo commit, e isso é uma tarefa própria com o seu próprio portão.
- **Não liga o `WireCatalog` ao arranque.** O composition root ganha a função; ninguém a chama. Ligá-la é uma decisão de operação, não de fundação.
- **Não toca `sync_state`.** O `Entity` enum que atravessa seis domínios fica intacto; decompô-lo é Onda 2.
- **Não acrescenta dependências.** Se alguma tarefa parecer precisar de uma, isso é `REQUEST` ao hub e paragem da tarefa — nunca um `go get`.

## Se algo bloquear durante a noite

A regra é uma só: **não inventes contorno e não alteres o plano.**

1. Escreve o bloqueio em `.mnfs/HARNESS-DEBTS.md` com `file:line` e o output verbatim que o motivou.
2. Faz commit do que já está feito e verde.
3. Passa à tarefa seguinte que **não dependa** da bloqueada. As dependências reais são: T2–T5 dependem de T1 só por convenção (são independentes entre si); T6 é independente de tudo; T7 depende de T2/T3/T5; T8 depende de T7; T9 depende de T8; T10 depende de T9 **e do Postgres de pé**; T11 depende de T7; T12 depende de T6 e T7.
4. Se o Postgres não estiver de pé, **T10 salta e T11 e T12 correm à mesma** — T11 não toca em base de dados nenhuma e T12 só precisa de T6 e T7.
