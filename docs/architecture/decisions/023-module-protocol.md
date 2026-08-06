# ADR-023: The module protocol — private internals, one published consumable

**Date:** 2026-08-05
**Status:** accepted
**Measured against:** `d57ea44b`. Every count below was measured in the repository, not
estimated. The measurement lives in
`docs/superpowers/plans/2026-08-05-arquitetura-protocolo-de-modulo-plan.md` §1 and
`docs/architecture/decisions/_citations/layer-vocabulary-census.md`.

## Context

`apps/server_core/internal/modules/` holds 21 modules. They were built one at a time, each
by a different mission, and the rule for how one module may talk to another was never
written down. The result is not chaos — it is worse, because it looks orderly. Every module
has `domain/`, `application/`, `ports/`, `adapters/`, `transport/` directories, so a reader
concludes there is a protocol. There is not.

What was actually measured:

- **121 cross-module imports** in non-test code. 76 originate in `adapters/`, which is the
  sanctioned translation point and not a violation. **35 remain**, and they are not spread
  thin: 24 of them target just three modules — `connectors/domain` (12), `internal_read/domain`
  (8), `sync/application` (4).
- **9 of those 35 originate in `X/ports`.** A module's port — the thing that is supposed to
  BE the boundary — is typed with another module's `domain`. The boundary is the leak.
- **12 distinct layer-directory names** in use, of which 8 directories contain a single
  one-line `doc.go` and nothing else.
- **126 interfaces in `ports/`**, of which **7 have no non-test reference anywhere** outside
  their own declaration file. Nothing implements them by name, nothing takes them as a
  parameter, nothing holds them as a field.

The reading is not that there are 35 bugs. It is that **three modules never published a
consumable**, so fourteen callers reached into their internals — which is the only thing
left to do when a module offers no front door. Fixing 14 callers treats the symptom;
publishing 3 consumables removes the cause.

## Decision

### §1 — What a module is

A module is a unit with **private internal logic** and **one public consumable**. Other
modules speak to it only through the consumable, and always the same way.

### §2 — The import rule

> **A module is importable by another module ONLY at `X/ports`. `X/domain`,
> `X/application`, `X/adapters` and `X/transport` are private.**

One line, and a grep can check it. The 35 measured violations are exactly the violations of
this line.

**Carve-out, sanctioned and singular:** `apps/server_core/internal/composition` is the
composition root. Assembling the dependency graph is its entire job, so it may import any
layer of any module. It is the only package with this exemption, and the exemption does not
transit: a module may not reach another module's internals by way of the root.

### §3 — A port must stand on its own

> **An interface in `X/ports` may be typed only with: types declared in `X/ports` itself,
> standard-library types, or types from the shared core (§5). Never with `Y/domain`.**

Without §3, §2 is theatre: the import becomes legal and the coupling survives. This is the
clause that closes the 9 leaking ports.

### §4 — When a port exists

> **A port exists when there are two or more consumers, OR when the crossing is
> technological (database, provider HTTP, Oracle, clock, queue). Never for symmetry.**

This is a counter-rule, and it is here deliberately. Without it, "standardising" produces
more interfaces, not fewer. The 7 ports with no consumer exist because nobody had written
this sentence down.

### §5 — A deliberately small shared core

Packages with no module dependencies that any module may import. Today: `sourcekind`. It
will receive the money type (§6).

Entry rule: a package joins the shared core only when **two or more** modules must name its
types in a public signature, and the package itself imports no module. The list is closed in
this ADR; adding to it requires an amendment.

### §6 — One money

> **In the domain: a single `Money` type, in the shared core, with no `float64`.
> On the wire: decimal as `string`.**

Four representations of money currently cross the wire simultaneously: `string` (pricing's
decomposition), `float64`/`*float64` (17 absolute-value fields, including
`orders/transport/http_handler.go:451-466` — `Comissao`, `Frete`, `Imposto`, `Custo`,
`MargemValor`, `TarifaFull`), `domain.Money` (3 fields), and `int`/`*int` (4 fields). The
same fiscal decomposition leaves `pricing` as `string` and leaves `orders` as `float64`.

`string` rather than a JSON number because a JSON number is a double, and a `float64` on the
money path cannot express absence — which violates ADR-017 §1 structurally, in 17 places
today.

Accepted consequence: converting `orders` from `float64` to `string` **breaks the published
contract**. It ships as one commit — OpenAPI, SDK, handler, frontend — once.

### §7 — The layer vocabulary is a closed list

**Canonical:** `domain`, `application`, `ports`, `adapters`, `transport`, `background`,
`wiring`.

Each non-canonical name found in the census was judged on its measured content, not its
plausibility:

| Name | Measured | Disposition |
|---|---|---|
| `readmodel` | 4 dirs, one line `package readmodel` each, **0 importers** | **Delete.** Scaffolding that makes the architecture read richer than it is. |
| `events` | 4 dirs, one line `package events` each, **0 importers** | **Delete.** Identical shape. |
| `composition` (in-module) | 4 dirs, 11 files, 1.082 LOC, 6 cross-module importers | **Canonical, renamed to `wiring`.** The code is load-bearing — it exists so other modules can obtain a ready scheduler without importing `adapters/postgres` directly. The name is not: `internal/composition` is the composition root and owns that word. Two different things named `composition` at two scopes is the ambiguity this ADR exists to remove. |
| `background` | 2 dirs, 8 files, 582 LOC, wired from the root | **Canonical.** A timer is an entry point, exactly as HTTP is. It is the scheduled-entry counterpart to `transport`, and calling it `application` would erase that. |
| `registry` | 1 dir, 180 LOC, 3 importers | **Fold.** The `MarketplaceConnector` interface moves to `marketplaces/ports`; the plugin `init()` registrations move to `adapters`. A single-occurrence pattern does not earn a vocabulary word. |
| `observability` | 1 dir, 467 LOC, 1 importer | **Fold into `adapters`.** It wraps `internal_read/ports` and logs; that is an adapter. |
| `integration` | 1 dir, 2 files, both `_test.go`, `//go:build integration` | **Not a layer.** It is a hermetic test lane directory. The lane discovers it correctly by build tag (`scripts/harness/Postgres.psm1:50-58`), so it runs — it simply must not be counted as architecture. |

After this ADR the list is closed and a new layer name is a checker failure, not a judgment
call.

### §8 — Modules that do not fit the mould

Two modules are outside it, and only one is legitimate:

- **`sourcekind`** — no layer directories at all. Correct: it is the shared type with no
  dependencies (§5).
- **`tenant_config`** — has `transport/` and nothing else, with domain-shaped code
  (`active_source.go`, `context.go`) and adapter-shaped code (`repository.go`, raw SQL
  against `*pgxpool.Pool` at lines 37-43, 74-82, 92-98) sitting loose at the module root.
  This is not a sanctioned shape. It is brought into the mould, or its exemption is written
  down with a reason and an owner.

## Rationale

Every clause here is the shortest rule that makes a measured defect impossible, and nothing
more.

§2 is one line because a rule that needs a paragraph needs a judgment, and a judgment cannot
be a gate. §3 exists because §2 alone is satisfiable while staying fully coupled — the
9 leaking ports are the proof, and they were written by people who believed they were
following the rule. §4 exists because the natural failure mode of a standardisation effort
is more structure: the 7 empty ports were each written in good faith, by symmetry with a
neighbour. §7 is a closed list rather than a guideline because 8 empty directories survived
years of review — an open vocabulary has no failure state, so nothing ever fails.

The shared core is capped at "two or more modules must name it in a public signature"
because a shared core that grows is a second, undeclared dependency graph.

## Consequences

- Three modules — `connectors`, `internal_read`, `sync` — must publish a consumable before
  their 24 invasive imports can close. That is the largest piece of work this ADR implies,
  and it lands before the money migration, because the money type has to cross those
  consumables.
- The composition root keeps its exemption, and therefore stays the one place where a change
  can quietly recouple everything. It is the file to read first in review.
- The layer rename (`composition` → `wiring`) touches 4 directories and 6 importing sites.
  Mechanical, and it is the compiler's job to find them.
- `float64` disappears from the money path, which makes several DTO fields pointer/nullable
  and changes the published shape of `orders`. Frontend ships in the same commit.
- Deletion of `readmodel`/`events` removes 8 directories and zero behaviour.
- **This ADR is inert until it is checked.** `GOV_MODULE_LAYER` today tests only
  `adapters`, `transport`, `registry` (`scripts/harness/Policy.psm1:328`), so `domain`,
  `application` and `ports` are structurally invisible to it — the 35 violations cannot be
  detected by the checker that exists to detect them. Extending it, and proving each check
  by a must-fail, is the closing wave of this work. A rule without a checker is a comment.

## Alternatives Considered

**Fix the 14 callers, leave the modules as they are.** Rejected: it treats 35 symptoms of
3 causes, and the next module to need `connectors` data has the same two options as before —
reach in, or write the consumable. The condition regenerates.

**Make everything public and rely on review.** Rejected on evidence: review is what we had.
It produced 8 empty directories, 7 unused ports, 9 leaking ports, and 4 simultaneous money
representations, over the same period in which the rules were being cited in code comments
1.365 times.

**Enforce the boundary with separate Go modules per module.** Rejected as disproportionate:
`go.mod` per module makes every internal refactor a version negotiation, for a boundary a
grep can enforce. Revisit only if the grep proves insufficient.

**A larger shared core (money, IDs, errors, time, tenant).** Rejected: each addition is a
type every module is then entitled to depend on, and a shared core is invisible coupling
precisely because importing it is always legal. Two-consumer minimum, closed list, amendment
to grow.
