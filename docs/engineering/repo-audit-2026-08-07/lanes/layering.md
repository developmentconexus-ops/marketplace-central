# Lane: layering

> I have run something really deep into MetalDocs and I am changing the way I code there to move to
> something more professional towards issues, PRs, PR review, CodeRabbit mechanical full validation
> and so much more. For that I had to identify every error in my code, my platform, to improve it and
> to create this full validation. I want to run it here as well so we move on the same path, this way
> it gets so much harder to send bad PRs.

Calibration: solid professional level, not Google-tier. Success condition is a **mechanism**
that makes a bad change hard to land, not a cleaner codebase.

## Findings

| ID | class | finding | evidence | scale |
|---|---|---|---|---|
| L-01 | drift | `contracts/governance/modules.json` `temporary_exceptions` (`rule_id: module-target-layer`) is a second, disconnected ledger for the exact rule `TestModuleBoundaryADR023` enforces. Neither reads the other. `Policy.psm1:98-102` (`Test-GovernanceContracts`) validates only that `source_module`/`target_module` are known ids and `rule_id`/`path` are well-formed — it never opens a `.go` file. 3 of 5 declared exception paths do not exist on disk; the other 2 (`tenant_config/active_source.go:11`, `integrations/adapters/feesync/marketplace_executor.go:11`) are live and appear inside the 234 the Go test currently reports, uncounted against the exception. All 5 carry `removal_owner: M-10`, a milestone the mission ledger records as "deferred" since MIS-003 superseded MIS-001 on 2026-07-14 (`.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:262`). | 5 exceptions total, 3/5 stale paths, 2/5 unreconciled with the 234; commands below |
| L-02 | duplication | The layering rule "core must not import a provider-specific type" is hand-implemented twice, byte-identical apart from names: `apps/server_core/internal/modules/channelfees/boundary_test.go` and `.../divergences/boundary_test.go` (`diff` below, only symbol/package names differ). Neither reuses `internal/arch.ScanVendorTokens`, which already implements a superset (identifiers *and* string literals, not just import paths) but is never pointed at `internal/modules/`. 19 of 21 modules have no such guard at all. | 2 files, ~50 lines each, 0 reuse of the general scanner; 19/21 modules uncovered |
| L-03 | gap | `internal/composition` — the ADR-023 carve-out root — is structurally exempt from `TestModuleBoundaryADR023` (it doesn't live under `internal/modules`, so the walk never visits it) **and** is never passed to `arch.ScanVendorTokens`, which the test file (`internal/arch/repo_test.go:32-63`) only ever calls with `internal/contexts` or `internal/kernel` as root. Design spec Regra 2.5 (`docs/superpowers/specs/2026-08-06-protocolo-de-codigo-design.md` §2, "A raiz de composição não declara adapters") already names this as the target state and cites 10 of 16 composition-declared types as Mercado Livre adapters — self-acknowledged in prose, backed by zero mechanical detector. | 43 vendor-token occurrences across 5 files in `internal/composition`, 0 of them ever scanned; command below |
| L-04 | gap | `TestModuleBoundaryADR023` only recognizes import paths with the literal prefix `.../internal/modules/` (`modulesImportPrefix`, `module_boundary_arch_test.go:35`). It has nothing to say about a module importing `internal/contexts/<x>` or `internal/adapters/<x>` — measured today at 0 occurrences either direction, so the tree is clean, but that cleanliness is accidental (small `contexts/` surface, one context landed) rather than monitored: the prefix match means a future `internal/modules/orders` importing `internal/contexts/pricing/port` would silently pass this test even though it is exactly the kind of cross-tree coupling ADR-023 exists to prevent. | 0 occurrences today (verified), but 0 detector coverage of the direction |
| L-05 | layering | Of the 234 `TestModuleBoundaryADR023` violations, 146 (62%) originate in the `adapters/` layer — the layer the test's own author calls out in comment as deliberately excluded from any carve-out ("an adapter reaching into a sibling module's domain is coupling with an extra step, not translation"). `application` contributes another 42 (18%). Only 9 originate in `ports/` itself, but the port surface that exists is typed by the target module's `domain`: `internal/modules/connectors/ports/catalog_read.go:6,10,13-15` types `CatalogReader`/`CatalogOffersReader` with `domain.ProviderAccountRef`, `domain.CatalogOffer`, `domain.CatalogProduct` — any legitimate consumer of the published port must also import `connectors/domain` to name the values it passes, which is exactly ADR-023's own diagnosis ("the boundary is the leak"), now generalized past the module it was measured on. | 146/234 adapters-origin, 42/234 application-origin, 9/234 ports-origin; 1 sampled port file shown |
| L-06 | layering | Shape of the 234: 49 distinct directed module→module edges among 17 of 21 modules. 3 bidirectional pairs: `catalog<->erp_import` (3/2), `catalog<->internal_read` (4/4), `erp_import<->internal_read` (8/18 — the single heaviest directed edge, 26 combined). `connectors` is the heaviest single sink: 64 violations from 6 distinct source modules (orders 16, integrations 15, mutations 14, listings 12, product_links 5, inventory 2), never through `connectors/ports`. `connectors` + `internal_read` + `erp_import` together absorb 132/234 (56%). | 49 edges / 17 modules / 234 sites; commands below |
| L-07 | layering | `internal/platform/archguard` is a third, independent layering instrument — not `TestModuleBoundaryADR023`, not `internal/arch`. It allowlists every interactive-request-time code path permitted to reach the sole Mercado Livre HTTP client, verified against the real `composition/root.go` plus 4 must-fail fixture tests (5th-site injection, aliased-site injection, allowlist-shrink). All 5 tests pass today. Its scope is explicit and self-documented as excluding batch/background/scheduler paths and `market/*` wiring "by design" (`archguard_test.go:12-27`) — no other detector covers what it excludes, so that exclusion is asserted in a comment, not proven by a test. | 5/5 tests pass; scope note at `archguard_test.go:12-27`, verified by run |
| L-08 | gap | No layering detector — `TestModuleBoundaryADR023`, `internal/arch`, `archguard`, or the two `boundary_test.go` files — ever looks outside `.go` files. `apps/web/src` and `packages/sdk-runtime/src` are entirely outside every mechanism this lane found. Concretely: FE test/page files carry 77 occurrences of `mercado_livre`/`mercadolivre` as data (`provider_code: "mercado_livre"`), which is not a violation of any stated rule (Regra 2.3 is scoped to `contexts/`, and the SDK types the field `string`, not a closed enum — matching the rule's *intent* even though the rule's *text* never reaches the frontend). Recorded as a gap in coverage, not a violation, because no rule currently claims this surface. | 77 FE occurrences, 0 SDK occurrences, 0 detectors scoped to `apps/web` or `packages/sdk-runtime` |

## The five heaviest, with detail

**1. `connectors` is the boundary's biggest, unfixed sink (L-05, L-06).** 64 of 234 violations
(27%) target `connectors/domain` (40), `connectors/application` (19) or `connectors/adapters` (5),
from 6 different source modules, and none of them go through `connectors/ports` even though that
package exists and has 14 interfaces. The reason is visible in the port file itself:
`internal/modules/connectors/ports/catalog_read.go`
```go
import (
    "context"
    "marketplace-central/apps/server_core/internal/modules/connectors/domain"
)
type CatalogOffersReader interface {
    ListCatalogOffers(ctx context.Context, accountRef domain.ProviderAccountRef, catalogProductID string) ([]domain.CatalogOffer, error)
}
```
A caller that only wants the port must still import `connectors/domain` to name
`ProviderAccountRef`/`CatalogOffer`. ADR-023's Context section made this exact observation about 9
of the original 35 ("9 of those 35 originate in X/ports... the boundary is the leak") — this lane's
measurement shows the same mechanism at 234-scale, concentrated on one module.

**2. `erp_import <-> internal_read` is the heaviest single coupling, and it is mutual (L-06).**
18 violations run `internal_read -> erp_import`, 8 run the other way — 26 combined, more than any
other module pair, and bidirectional, meaning neither module can be said to own the relationship.
ADR-023's Context section already named this pair as "the maior aresta do grafo, nos dois sentidos"
at the 128-count measurement; it is still the heaviest edge at 234.

**3. The governance exception ledger for this exact rule has decayed into decoration (L-01).**
`contracts/governance/modules.json` carries a `temporary_exceptions` array keyed on
`rule_id: "module-target-layer"` — the same ADR-023 §2 rule — with 5 entries, each pointing at a
specific file:line and a `removal_owner`. None of that is checked against the codebase:
`Policy.psm1`'s `Test-GovernanceContracts` (lines 86-102) validates only that `source_module`,
`target_module` and `rule_id` reference real module ids — never that the `path` exists, never that
the import it describes is still there. Verified: `connectors/adapters/magalu/fee_seed.go`,
`connectors/adapters/mercado_livre/fee_sync.go` and `connectors/adapters/shopee/fee_seed.go` (3 of
5 entries) do not exist on disk at all. The other 2 do exist and their imports are real — but they
show up inside the Go test's 234-count with no cross-reference either way. All 5 share
`removal_owner: "M-10"`; `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:262` records
"M-10/11/12 stay deferred" as of 2026-07-14, under a mission (MIS-001) that mission itself has since
been superseded. The mechanism that is supposed to retire these exceptions has no owner currently
executing.

**4. The composition root is the one place Regra 2.5 names and the one place nothing scans (L-03).**
`docs/superpowers/specs/2026-08-06-protocolo-de-codigo-design.md` §2, Regra 2.5, states plainly "a
raiz de composição não declara adapters" and cites 10 of composition's 16 declared types as
Mercado Livre adapters as the defect the rule exists to catch. That citation is prose, not a test.
`internal/composition` is invisible to `TestModuleBoundaryADR023` by construction (outside
`internal/modules`) and is never passed to `arch.ScanVendorTokens` (only `internal/contexts` and
`internal/kernel` are, per `internal/arch/repo_test.go:32-63`). Measured directly: 43 vendor-token
occurrences across `market_adapters.go`, `orders_adapters.go`, `orders_ingest_adapters.go`,
`pricing_adapters.go`, `root.go`. This is consistent with the design doc's narrative but is
currently unmeasured by anything that runs — there is no ratchet tracking whether this number goes
up or down as new contexts land.

**5. The one rule that is enforced twice by hand is enforced nowhere generally (L-02).**
`channelfees/boundary_test.go` and `divergences/boundary_test.go` are the same ~50-line AST-import
scanner with only names changed (`diff` shows exactly that). Both exist because those two modules'
own feature specs stated a "Q6" constraint at build time. `internal/arch.ScanVendorTokens` already
implements a strictly more thorough version of the same check (identifiers and string literals, not
just import paths) and is already wired into a test file — but only ever against `internal/contexts`
and `internal/kernel`. Nobody pointed it at `internal/modules/channelfees` or
`internal/modules/divergences` to retire the two hand-written copies, and nobody pointed anything at
the other 19 modules, 17 of which are demonstrably not clean of provider-specific coupling (L-06).

## What is actually fine

- `TestModuleBoundaryADR023` itself does exactly what ADR-023 §2 says, including both ratified
  carve-outs (`internal/composition` invisible by construction; `sourcekind` explicit in
  `sharedCore`), and it walks `_test.go` files too (the D-49 fix from 2026-08-06 is real — confirmed
  by rerunning: the current run reports `_test.go` sites throughout the `sites:` list).
- `internal/arch`'s four scanners (`ScanCrossContextInternal`, `ScanFloatInContracts`,
  `ScanVendorTokens`, `ScanFactValueDiscard`) are well-built: each has must-fail fixture tests
  proving the detector actually fires (not just "green because it never ran"), and all 14 tests in
  `internal/arch` pass today (`go test ./internal/arch/...`, verified).
- Zero cross-imports exist today between `internal/modules` and `internal/contexts` in either
  direction, and zero between `internal/modules` and `internal/composition` (modules never reach
  into the composition root). Verified by direct grep, both directions. The tree that exists today
  is clean; see L-04 for why that isn't the same as monitored.
- `internal/adapters/erp/sankhyaoracle` already follows the "molde" the design spec describes for
  vendor/context roots (Regra 2.2-a): one exported package, one constructor (`New(db *sql.DB,
  instance string, now func() time.Time)`), a `Bundle` typed by `port.ProductFeed`, wire types
  genuinely unreachable under its own `internal/`.
- `internal/contexts/catalog/module.go`'s `New(pool *pgxpool.Pool) *Module` is the same shape,
  cited directly by ADR-033 as the model — confirmed matching on read.
- `internal/platform/archguard` is a well-built, narrowly-scoped, currently-green guard with 5
  must-fail fixture tests that actually inject a 5th site, an aliased site, and a shrunk allowlist
  and confirm the detector reacts correctly to each (verified by run, not by reading the test names).
- SDK and frontend type `provider_code` as `string`, not a closed union/enum — matching Regra 2.3's
  intent ("código de canal é dado em runtime, nunca enum fechado") even though the rule's text
  doesn't formally extend past `contexts/`.
- 4 of 21 modules — `channelfees`, `classifications`, `divergences`, `sourcekind` — appear in zero
  of the 234 boundary violations, either as source or target.

## Unverified / needs judgment

- Whether the `internal/composition` vendor-token count (43) is growing, shrinking, or flat over
  time — no historical snapshot was available in-lane; only the point-in-time count was measured.
- Whether `M-10`'s "deferred" status (as of the 2026-07-14 MIS-003 replan note) still holds today —
  the mission ledger entry is the only evidence found; no more recent status update for M-10
  specifically was located.
- Whether other `internal/modules/<x>/boundary_test.go`-style ad hoc guards exist for concerns other
  than provider-token leakage — this lane searched specifically for `boundary_test.go` naming and
  for `VendorToken`/`Guard`/`Boundary`/`Layer` keyword co-occurrence; a differently-named hand-guard
  would not have been found by that search.
- `contracts/governance/invariants.json` carries 4 unrelated `removal_owner: HARNESS-D-N` exceptions
  (all `rule_id: production-panic`) that resemble L-01's shape but are not a layering-rule concern —
  noted only to avoid double-counting them against L-01; PHASE-0 already flags this file's D-2/D-50
  question as a separate, cicd/governance-lane matter.

## Commands run

```
# ADR-023 boundary test, full run with breakdown
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/composition/... -run TestModuleBoundaryADR023 -v

# internal/arch scanners, full suite incl. must-fail fixtures
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/arch/... -v

# archguard, full suite incl. must-fail fixtures
cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/platform/archguard/... -v

# module->module edge aggregation from the 234 sites (awk over saved test output)
awk '/^          internal\/modules\// { ... split "from -> to", count by module pair }' boundary_out.txt | sort -rn

# mutual-edge and distinct-module-count check (same awk script, extended)

# governance exception paths: existence check
ls apps/server_core/internal/modules/connectors/adapters/magalu/fee_seed.go \
   apps/server_core/internal/modules/connectors/adapters/mercado_livre/fee_sync.go \
   apps/server_core/internal/modules/connectors/adapters/shopee/fee_seed.go
# -> all 3: No such file or directory

# governance exception paths that DO exist, cross-checked against the 234 sites list
grep -n "feesync/marketplace_executor\|tenant_config/active_source" boundary_out.txt
# -> both present in the 234, unreconciled with the exception entry

# M-10 status
grep -rn "M-10" .mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md

# composition-root vendor token count (unscanned by any wired test)
grep -rniE "mercado_livre|mercadolivre|mercadolibre|\bmeli\b|shopee|amazon|magalu|americanas" \
  apps/server_core/internal/composition --include=*.go | grep -v _test.go | wc -l
# -> 43

# byte-diff of the two duplicated boundary tests
diff apps/server_core/internal/modules/channelfees/boundary_test.go \
     apps/server_core/internal/modules/divergences/boundary_test.go

# cross-tree import check, both directions, modules<->contexts and modules<->composition
grep -rn "internal/contexts" apps/server_core/internal/modules --include=*.go
grep -rln "internal/modules" apps/server_core/internal/contexts --include=*.go
grep -rln "internal/composition" apps/server_core/internal/modules --include=*.go

# frontend/SDK vendor-token scan (outside every detector's reach)
grep -rniE "mercado_livre|mercadolivre|mercadolibre|\bmeli\b" apps/web/src --include=*.ts --include=*.tsx | wc -l
grep -rn "mercado_livre|mercadolivre" packages/sdk-runtime/src/index.ts
```
