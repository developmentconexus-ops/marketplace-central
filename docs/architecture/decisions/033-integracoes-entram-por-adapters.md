# ADR-033: External marketplace integrations enter through `adapters/marketplace/<vendor>`, not `connectors`

**Date:** 2026-08-07
**Status:** accepted, amends `ARCHITECTURE.md` frozen decision 7

## Context

`ARCHITECTURE.md` frozen decision 7 (`ARCHITECTURE.md:39`) reads: "External marketplace
integrations enter only through the `connectors` module via port interfaces"
`ARCHITECTURE.md:5` states that frozen decisions "must not be rediscussed without an
explicit ADR." The repository's order of truth puts `ARCHITECTURE.md`/ADRs above the
design spec, so until this document exists, the ratified design's routing of marketplace
integrations to `adapters/marketplace/<vendor>` has not actually won, no matter how
measured its case is.

`connectors` is a module under `apps/server_core/internal/modules/`, of the exact kind
this plan's parent mission (the module-protocol work, ADR-023) is decomposing: private
internals plus one published consumable, replaced context by context by
`apps/server_core/internal/contexts/`. `docs/superpowers/specs/2026-08-06-protocolo-de-codigo-design.md`
§1.2 lists `connectors` among the eight current modules that "deixam de ser módulos,"
migrating to `adapters/marketplace/{mercadolivre,shopee,amazon}/`, and §2 puts the vendor
boundary at the Go compiler's `internal/` rule (`adapters/marketplace/<vendor>/internal/api`),
not at a module-naming convention.

This plan's own tasks already demonstrate the mechanism the design relies on, applied to a
context rather than a vendor: `apps/server_core/internal/contexts/catalog/module.go:23-38`
now reads

```go
// New assembles the context from the ONLY thing an outsider may legitimately
// name: a connection pool.
//
// Every other collaborator — the store, the id factory, the reader — is chosen
// here, inside the context, because their types live under internal/ and a
// parameter typed by one of them would have no legal caller.
func New(pool *pgxpool.Pool) *Module {
```

— a single-parameter facade whose private collaborators live under `internal/` and are
unnameable from outside the tree. That is the same shape the design spec assigns to
`adapters/marketplace/<vendor>`, and it is a shape the `connectors` module, as a legacy
`internal/modules/` module rather than a boundary enforced by `internal/`, does not itself
enforce structurally.

## Decision

**External marketplace integrations enter through `adapters/marketplace/<vendor>`,
implementing context ports declared in `contexts/<name>/port`, not through the
`connectors` module.** This supersedes the routing stated in frozen decision 7 of
`ARCHITECTURE.md`, which is amended to cite this ADR for the current routing while the
decision itself — that integrations are boundaried behind port interfaces, never called
ad hoc from application code — is preserved.

**§1 — Transition, not a rewrite.** `apps/server_core/internal/modules/connectors/`
continues to serve Mercado Livre unchanged for as long as the context tree that replaces
it has not landed. It is inventoried, never grown — `ARCHITECTURE.md` frozen decision 11
("Global-maximum design beats local patches") and this plan's global constraint 1 both
already forbid extending it.

**§2 — No new marketplace code enters `connectors`.** From this ADR forward, new
marketplace-integration code is written under `adapters/marketplace/<vendor>/`, per the
design spec §2 (Regras 2.1–2.5): the port belongs to the consuming context
(Regra 2.1), the vendor's wire DTOs stay unreachable outside the vendor's own tree
(Regra 2.2), the vendor root exposes a single constructor typed only by consumer ports
(Regra 2.2-a), no vendor name appears inside `contexts/` (Regra 2.3), and there is no
single `Marketplace` interface (Regra 2.4).

**§3 — `connectors` is deleted only when its replacement lands.** Per the design spec §15
step 5, `internal/` closes behind each context as it lands, not behind all contexts at
once. `connectors` is removed the same way, not on a fixed date.

**§4 — `ARCHITECTURE.md` is updated.** Frozen decision 7 is amended to read: "External
marketplace integrations enter through `adapters/marketplace/<vendor>`, implementing
ports owned by the consuming context — see ADR-033. `connectors`
(`apps/server_core/internal/modules/connectors/`) continues to serve Mercado Livre during
the migration and receives no new marketplace code."

## Rationale

Decision 7 named a module, `connectors`, as the sole entry point for marketplace
integrations, at a time when "module under `internal/modules/`" was the only pattern the
repository had. The module-protocol mission (ADR-023) changed what a boundary is made of:
under `internal/contexts/` and `adapters/`, the boundary is the Go compiler's `internal/`
rule, not a name a reviewer has to remember and enforce. Naming `connectors` specifically,
rather than the general shape "one adapter tree per vendor, entered through ports,"
freezes the repository into the module it is actively replacing. The design's routing
does not contradict the intent behind decision 7 — integrations still enter only through
port interfaces, never ad hoc — it relocates the enforcement to where this plan's other
completed tasks already put it for contexts.

## Consequences

- `ARCHITECTURE.md` frozen decision 7 cites this ADR instead of naming `connectors`
  directly.
- No task after this one may add marketplace-integration code to
  `apps/server_core/internal/modules/connectors/`; new code goes to
  `adapters/marketplace/<vendor>/`.
- `connectors` remains live and unchanged until the context(s) that consume Mercado Livre
  data have their own `adapters/marketplace/mercadolivre` tree ready to replace it — this
  ADR does not schedule that migration, only the entry point for new code.
- The design spec's own migration table (§1.2) and adapter layout (§2) are now backed by
  an ADR and are no longer merely proposed.

## Alternatives Considered

**Keep frozen decision 7 verbatim and treat `adapters/marketplace/<vendor>` as a rename
of `connectors`.** Rejected: `connectors` is a flat `internal/modules/` module whose
layers (`domain`, `application`, `ports`, `adapters`, `transport`) sit in ordinary,
compiler-unenforced directories (ADR-023 §2, §7). A rename would carry the same
enforcement gap forward under a new name; the point of the migration is the change in
what backs the boundary, not the path.

**Amend decision 7 to say "the current context or the `connectors` module, whichever
applies," and skip a dedicated ADR.** Rejected: `ARCHITECTURE.md:5` requires an explicit
ADR for any frozen-decision change, and a two-target decision has no single term the
compiler or a detector can key on — it reintroduces the judgment call this migration
exists to remove.
