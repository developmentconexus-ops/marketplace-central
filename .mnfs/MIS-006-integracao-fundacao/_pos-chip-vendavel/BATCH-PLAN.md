# CHIP-VENDAVEL — P2 batch plan (verbatim worker output)

Source: gpt-5.6-sol / medium, read-only, cwd = chip worktree. Log+last-message in the
session scratchpad; ledger row 2 points at them. Everything below is the planner verbatim;
the chip amendments live in the section appended at the END of this file and WIN over it.

## Design decisions

### Storage choice

Put `only_revenda`, `only_em_estoque`, and `only_ecommerce_eligible` directly on `active_source`.

Evidence:

- `active_source` already has exactly one row per tenant, keyed by `tenant_id` ([`0076_products_mirror_active_source.sql:57`](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-vendavel/apps/server_core/migrations/0076_products_mirror_active_source.sql:57)).
- `tenant_config.Repository.Get` already performs the required tenant-scoped lookup at [`repository.go:37`](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-vendavel/apps/server_core/internal/modules/tenant_config/repository.go:37).
- Adding a sibling table would introduce a second one-to-one lifetime, another join, and missing-row semantics without supporting a second consumer or cardinality.
- Migration defaults can represent the operator-fixed initial configuration exactly:
  `BOOLEAN NOT NULL DEFAULT true`, `true`, `false`.
- `Repository.Set` must continue updating only active-source fields, so changing the ERP source cannot reset assortment configuration.

Migration allocation:

- `0083_sellable_assortment_config.sql`: three columns on `active_source`.
- `0084_products_mirror_sellable_fields.sql`: nullable `usoprod TEXT`, `ad_ecommerce TEXT`.
- Both additions raise both migration-count fixtures from `69` to `71` at [`runner_test.go:25`](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-vendavel/apps/server_core/internal/platform/migrate/runner_test.go:25) and [`runner_test.go:64`](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-vendavel/apps/server_core/internal/platform/migrate/runner_test.go:64).

### Catalog HTTP shape

Use these additive surfaces:

- `GET /config/sellable-assortment`
- `PUT /config/sellable-assortment`
- `GET /catalog/products/counts` → `{sellable_count, total_count}`
- Add optional `include_all: boolean = false` to:
  - `GET /catalog/products`
  - `GET /catalog/products/search`

Absent `include_all` means filtered. `include_all=true` is request-local context only and never calls the configuration mutation. Explicit `?ids=` reads remain unfiltered because they serve already-selected products outside the `/catalogo` browse/search flow.

The count endpoint always returns configured-rule `N` and unfiltered `M`, independently of `include_all`.

### Decision request DR-1 — upload-source catalog parity

Current behavior:

1. `routing.Reader.resolve` sends `xlsx` and `catalogo_cliente` to the upload reader at [`routing/reader.go:51-54`](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-vendavel/apps/server_core/internal/modules/internal_read/adapters/routing/reader.go:51), while Sankhya goes to the live reader at lines 55–59.
2. Upload catalog paging reaches `MirrorCatalogPage`, whose SQL at [`mirror_query_repository.go:130-143`](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-vendavel/apps/server_core/internal/modules/erp_import/adapters/postgres/mirror_query_repository.go:130) currently filters only tenant, source, presence, numeric identity, cursor, and search.
3. Therefore, if the rule is added only to Oracle:
   - an XLSX tenant sees the full unfiltered mirror on `/catalogo`;
   - a source-routed count either returns `503 source_unavailable` because the upload reader lacks the count capability, or becomes false if it is accidentally computed against Oracle instead of the screen’s mirror source.

Recommendation to HUB: approve identical honest-unknown predicates and live counts on the mirror catalog path.

Predicate:

```sql
AND (NOT only_revenda OR m.usoprod IS NULL OR m.usoprod = 'R')
AND (NOT only_em_estoque OR m.estoque_total IS NULL OR m.estoque_total > 0)
AND (NOT only_ecommerce_eligible OR m.ad_ecommerce IS NULL OR m.ad_ecommerce <> 'N')
```

Estimated cost: four production files and four existing test files, about 95 production plus 120 test lines:

- `erp_import/ports/repository.go`
- `erp_import/adapters/postgres/mirror_query_repository.go`
- `erp_import/adapters/internalread/reader.go`
- `erp_import/domain/mirror.go`
- their existing reader, source-contract, search, and integration tests

This is the cheapest option that keeps page and counter truthful for every already-supported source. Slice S7 remains non-dispatchable until HUB answers DR-1.

---

## Slice cards

### S0 — Source-parity decision gate

- `id`: `S0-DR1`
- `goal`: Obtain HUB approval for applying the rule and counts to the upload/mirror catalog path.
- `write_set`: `[]`
- `failing_test_first`: `N/A — decision gate`
- `done_criteria`: HUB records either “approve mirror parity” or a named, dated deferral with an explicit non-numeric XLSX UI state; no implementation proceeds by silently serving a false count.
- `validation_kind`: `manual`
- `commands`: `[]`
- `expected_artifacts`:
  - `.mnfs/MIS-006-integracao-fundacao/_pos-chip-vendavel/DR-1-upload-catalog-parity.md`
- `complexity`: `standard`
- `open_questions`:
  - `DR-1: Approve the recommended mirror catalog predicate/count implementation (~8 files/~215 lines)?`

### S1 — Database and nullable domain shape

- `id`: `S1-SCHEMA`
- `goal`: Add persistent toggle defaults and nullable mirror assortment facts without converting unknowns to empty/zero.
- `write_set`:
  - `apps/server_core/migrations/0083_sellable_assortment_config.sql`
  - `apps/server_core/migrations/0084_products_mirror_sellable_fields.sql`
  - `apps/server_core/migrations/sellable_assortment_test.go`
  - `apps/server_core/internal/platform/migrate/runner_test.go`
  - `apps/server_core/internal/modules/erp_import/domain/import.go`
  - `apps/server_core/internal/modules/erp_import/domain/mirror.go`
  - `apps/server_core/internal/modules/erp_import/domain/mirror_test.go`
  - `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go`
- `failing_test_first`: `TestSellableAssortmentMigrationsDeclareDefaultsAndNullableMirrorFields` in `apps/server_core/migrations/sellable_assortment_test.go`
- `done_criteria`:
  - `active_source` has three non-null booleans with defaults `true`, `true`, `false`.
  - `products_mirror.usoprod` and `.ad_ecommerce` are nullable text with no default.
  - `NormalizedRow`, `MirrorProduct`, and `mirror.Row` carry pointer-valued fields.
  - Both migration inventory assertions expect exactly `71`.
- `validation_kind`: `unit`
- `commands`:

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\apps\server_core'
$env:GOCACHE = (Join-Path (Get-Location) '.gocache')
$env:GOMODCACHE = (Join-Path (Get-Location) '.gomodcache')
go test ./migrations ./internal/platform/migrate -run 'TestSellableAssortmentMigrationsDeclareDefaultsAndNullableMirrorFields|TestCanonicalSource' -count=1 -v
```

- `expected_artifacts`:
  - `S1-red.txt` naming `TestSellableAssortmentMigrationsDeclareDefaultsAndNullableMirrorFields`
  - `S1-green.txt`
  - migration inventory showing `71`
- `complexity`: `standard`
- `open_questions`: `[]`

### S2 — Tenant configuration persistence

- `id`: `S2-CONFIG-DB`
- `goal`: Persist and retrieve all three toggles per tenant through the existing tenant-config repository.
- `write_set`:
  - `apps/server_core/internal/modules/tenant_config/active_source.go`
  - `apps/server_core/internal/modules/tenant_config/active_source_test.go`
  - `apps/server_core/internal/modules/tenant_config/context_test.go`
  - `apps/server_core/internal/modules/tenant_config/repository.go`
  - `apps/server_core/internal/modules/tenant_config/repository_test.go`
- `failing_test_first`: `TestRepository_SetSellableAssortment_RoundTripPersistsPerTenant` in `apps/server_core/internal/modules/tenant_config/repository_test.go`
- `done_criteria`:
  - `Config` contains an explicit `SellableAssortment` value.
  - `Get` selects all three columns with `WHERE tenant_id=$1`.
  - `SetSellableAssortment` performs a tenant-predicated update and returns `ErrUnknownActiveSource` when no row exists.
  - `Set` for `active_source` does not overwrite the toggles.
  - Two tenants can hold different rules without leakage.
- `validation_kind`: `integration`
- `commands`:

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\apps\server_core'
$env:GOCACHE = (Join-Path (Get-Location) '.gocache')
$env:GOMODCACHE = (Join-Path (Get-Location) '.gomodcache')
go test ./internal/modules/tenant_config/... -run '^TestRepository_SetSellableAssortment_RoundTripPersistsPerTenant$' -count=1 -v
```

- `expected_artifacts`:
  - Red and green outputs with `=== RUN`, not a skip
  - SQL result showing two distinct tenant rows
- `complexity`: `standard`
- `open_questions`: `[]`

### S3 — Optional XLSX columns and upload writer

- `id`: `S3-XLSX`
- `goal`: Accept XLSX files with or without `USOPROD` and `AD_ECOMMERCE`, preserving missing values as NULL.
- `write_set`:
  - `apps/server_core/internal/modules/erp_import/adapters/xlsx/parser.go`
  - `apps/server_core/internal/modules/erp_import/adapters/xlsx/parser_test.go`
  - `apps/server_core/internal/modules/erp_import/adapters/postgres/mirror_repository.go`
  - `apps/server_core/internal/modules/erp_import/adapters/postgres/mirror_repository_test.go`
  - `apps/server_core/internal/modules/erp_import/adapters/postgres/mirror_repository_integration_test.go`
- `failing_test_first`: `TestParserSellableColumnsAreOptionalAndHonestUnknown` in `apps/server_core/internal/modules/erp_import/adapters/xlsx/parser_test.go`
- `done_criteria`:
  - Neither column joins the required-column lists at parser lines 22–35.
  - Header folding recognizes `USOPROD` and `AD_ECOMMERCE`.
  - A sheet without them yields nil pointers and remains accepted.
  - A sheet with them maps exact values.
  - The stage table, `CopyFrom`, insert, and conflict update all carry both fields.
  - Post-import SQL returns NULL for the first fixture and exact `R`/`S` for the second.
- `validation_kind`: `integration`
- `commands`:

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\apps\server_core'
$env:GOCACHE = (Join-Path (Get-Location) '.gocache')
$env:GOMODCACHE = (Join-Path (Get-Location) '.gomodcache')
go test ./internal/modules/erp_import/adapters/xlsx -run '^TestParserSellableColumnsAreOptionalAndHonestUnknown$' -count=1 -v
go test -tags=integration ./internal/modules/erp_import/adapters/postgres -run '^TestMergeSnapshotPersistsOptionalSellableColumns$' -count=1 -v
```

- `expected_artifacts`:
  - RED/GREEN outputs for both spreadsheets
  - SQL rows proving NULL versus exact values
- `complexity`: `standard`
- `open_questions`: `[]`

### S4 — Sankhya sync fields and stock-company pin

- `id`: `S4-SANKHYA-SYNC`
- `goal`: Read both Sankhya flags, write them to the mirror, and restrict Q4 to companies 1 and 2.
- `write_set`:
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go`
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/sync_test.go`
  - `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go`
  - `apps/server_core/internal/modules/internal_read/adapters/mirror/writer_test.go`
  - `apps/server_core/internal/modules/internal_read/adapters/mirror/writer_integration_test.go`
- `failing_test_first`: `TestSankhyaStockSQLPinsSellableCompanies` in `apps/server_core/internal/modules/internal_read/adapters/oracle/sync_test.go`
- `done_criteria`:
  - Q1 selects `p.USOPROD, p.AD_ECOMMERCE`, scans `sql.NullString`, and maps through the existing `nullStr` helper at [`sync.go:298`](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-vendavel/apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go:298).
  - Q4 contains exact `AND CODEMP IN (1, 2)` after `CODPARC = 0`.
  - Sankhya `upsertSQL` includes both fields in insert and conflict update.
  - Sync fixture asserts exact `R`, `V`, `S`, and nil values—not mere presence.
- `validation_kind`: `unit`
- `commands`:

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\apps\server_core'
$env:GOCACHE = (Join-Path (Get-Location) '.gocache')
$env:GOMODCACHE = (Join-Path (Get-Location) '.gomodcache')
go test ./internal/modules/internal_read/adapters/oracle -run '^TestSankhyaStockSQLPinsSellableCompanies$|^TestSankhyaSyncMapsSellableAssortmentFields$' -count=1 -v
go test ./internal/modules/internal_read/adapters/mirror -run '^TestPgWriterPersistsSellableAssortmentFields$' -count=1 -v
```

- `expected_artifacts`:
  - `failure_token=test=TestSankhyaStockSQLPinsSellableCompanies`
  - Green Q1/Q4 mapping output
  - HUB-owned post-sync SQL counts
- `complexity`: `complex`
- `open_questions`: `[]`

### S5 — Link-candidate filtering

- `id`: `S5-MATCHER`
- `goal`: Prevent non-sellable mirror rows from becoming link candidates while preserving honest-unknown and positive cases.
- `write_set`:
  - `apps/server_core/internal/modules/erp_import/domain/mirror.go`
  - `apps/server_core/internal/modules/erp_import/domain/mirror_test.go`
  - `apps/server_core/internal/modules/erp_import/adapters/internalread/reader.go`
  - `apps/server_core/internal/modules/erp_import/adapters/internalread/reader_test.go`
  - `apps/server_core/internal/modules/internal_read/adapters/routing/matcher.go`
  - `apps/server_core/internal/modules/internal_read/adapters/routing/matcher_test.go`
  - `apps/server_core/internal/modules/internal_read/adapters/routing/matcher_integration_test.go`
  - `apps/server_core/internal/modules/erp_import/adapters/postgres/mirror_query_repository.go` —
    **added D-122 by hub GRANT** (ruling 1 of 4), bounded to exactly two edits: (a)
    `mirrorProductColumns:13` gains `,m.usoprod,m.ad_ecommerce`; (b) `scanMirrorProduct:19-44`
    gains two `sql.NullString` read through the existing `nullStringPointer`. **Zero filtering in
    this file** — mirror-side SQL filtering remains S7's seam. The hub accepted the ordering
    argument: the first CONSUMER of a column brings it, rather than a preparatory slice landing a
    column nobody reads yet.
- `failing_test_first`: `TestMirrorMatcher_ActiveRevendaRuleControlsCandidateBirth` in `apps/server_core/internal/modules/internal_read/adapters/routing/matcher_integration_test.go`
- `blocker` (D-122, **RESOLVED by the grant above** — kept in full because the RED it produced is
  this slice's evidence that the green which follows is not vacuous): the mirror READ path never selected the two
  columns. `mirror_query_repository.go:13` defines `mirrorProductColumns` without `usoprod` /
  `ad_ecommerce`, and `scanMirrorProduct:19-44` does not read them — 0 occurrences in the file —
  even though the WRITE path populates both (S3, S4) and `domain.MirrorProduct` carries them (S1).
  Every mirror row therefore arrives `Usoprod=nil`, and the ratified mirror clause is
  `IS NULL OR = 'R'`, where nil PASSES. Implementing the slice without that file would ship a rule
  that compiles, has a green test, and cuts nothing: honest-unknown against a genuine gap (a
  spreadsheet without the column) is byte-identical to honest-unknown against a gap we created by
  not doing the SELECT. Requested extension is additive and bounded to two edits — the constant
  gains `,m.usoprod,m.ad_ecommerce`, the scanner gains the two `sql.NullString` through the
  existing `nullStringPointer`. No filtering there; mirror SQL filtering stays S7's. Blast radius
  checked: all four read sites (`MirrorRows:54`, `MirrorProductByCode:79`, `MirrorProductsByCodes:100`,
  `MirrorCatalogPage:131`) interpolate the SAME constant and go through the SAME scanner, so the
  column list and the `Scan` order move together by construction.
- `done_criteria`:
  - `MirrorMatcher` pins the fetched rule alongside the active source before calling the mirror reader.
  - Filtering happens in `FindProductsForLinking` after loading the cached index, not while building the TTL cache; a database toggle takes effect without waiting for cache expiry.
  - `USOPROD='V'` yields no candidate when enabled, then yields one after `only_revenda=false` is persisted.
  - A known sellable row is asserted born.
  - A nil `USOPROD` row is asserted born.
  - **The columns are asserted to ARRIVE, not merely to be selectable (D-122).** A row written with
    `USOPROD='V'` must be read back with `Usoprod` equal to `"V"`. Without this the three criteria
    above all pass against `Usoprod=nil`: nil PASSES the ratified mirror clause `IS NULL OR = 'R'`,
    so "cut" and "not cut" are indistinguishable when the SELECT is the thing that is broken. This
    criterion is what the blocker below was blocking.
- `validation_kind`: `integration`
- `commands`:

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\apps\server_core'
$env:GOCACHE = (Join-Path (Get-Location) '.gocache')
$env:GOMODCACHE = (Join-Path (Get-Location) '.gomodcache')
go test -tags=integration ./internal/modules/internal_read/adapters/routing -run '^TestMirrorMatcher_ActiveRevendaRuleControlsCandidateBirth$' -count=1 -v
go test ./internal/modules/erp_import/adapters/internalread ./internal/modules/internal_read/adapters/routing -count=1
```

- `expected_artifacts`:
  - `failure_token=test=TestMirrorMatcher_ActiveRevendaRuleControlsCandidateBirth`
  - SQL before/after toggle
  - Candidate IDs before/after
- `complexity`: `standard`
- `open_questions`: `[]`

### S5B — The xlsx write chain, and canonical case at the write seam

Added D-122 by hub ruling after S5's P4. Two rulings live here.

**Why it exists.** S5 landed green and the rule was still inert for `source=xlsx`. Measured chain:
`xlsx/parser.go:116` reads `USOPROD` off the sheet → `import_repository.go:67` copies 14 columns
into `erp_import_products` and **neither of the two is among them**, because migration 0084 added
them to `products_mirror` ONLY → `mirror_repository.go:30` syncs 12 columns. So mirror `usoprod`
is NULL forever for xlsx, nil PASSES the ratified `IS NULL OR = 'R'`, and the rule cuts nothing.
Sankhya is unaffected: `mirror/writer.go:150` writes straight to `products_mirror`.

**Authority is VC-6, not preference.** The contract already says the xlsx path WITH the columns
populates them — populates them in the MIRROR, or the sentence asserts nothing. The round-trip was
always contracted; the plan simply never gave it an owner. Leaving it out of the chip would
violate a ratified contract, so "accept out-of-chip" was never an option.

**The lesson, registered because it is the reusable part.** S3 was marked `completed` having
delivered a value that could not reach a consumer — a parser with no destination column is half a
chain. **A write-slice is verified by asking "does the value ARRIVE at the reader?", never "does
the value LEAVE the source?"** The same question kills the S5 blocker, this one, and the next.

- `id`: `S5B-WRITE-CHAIN`
- `goal`: Close the sheet→mirror round-trip for both sellable columns, and canonicalize their case
  once at the write seam so no reader has to fold.
- `migration_grant`: **0085 GRANTED** — the chip's block becomes 0083–0085. Adds `usoprod` and
  `ad_ecommerce` to `erp_import_products`, nullable, no default (ADR-17: absent ≠ a value).
- `write_set`:
  - `apps/server_core/migrations/0085_erp_import_products_sellable_fields.sql`
  - `apps/server_core/internal/modules/erp_import/adapters/postgres/import_repository.go` — the
    `CopyFrom` column list at `:67` and its row builder.
  - `apps/server_core/internal/modules/erp_import/adapters/postgres/mirror_repository.go` — the
    sync `SELECT` at `:30` and the mirror upsert it feeds.
  - `apps/server_core/internal/modules/erp_import/adapters/xlsx/parser.go` — normalization at the
    point `NormalizedRow` is born.
  - `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go` — the SAME
    normalization on the Sankhya write at `:150`.
  - `apps/server_core/internal/modules/erp_import/adapters/internalread/reader.go` — comment ONLY
    (see the case criterion); no behaviour change.
  - the corresponding `_test.go` files, plus
    `internal_read/adapters/routing/matcher_integration_test.go` (drops its raw seeding `UPDATE`).
- `failing_test_first`: a round-trip test that fails today because the value never arrives.
- `done_criteria`:
  - **The full round trip is proven by one must-fail:** a sheet carrying `USOPROD` → cell → row →
    `erp_import_products` → sync → `products_mirror` → the value READ BACK equals what the sheet
    said. Assert the VALUE. A test that stops at `erp_import_products` re-creates the defect one
    hop later.
  - **S5's raw seeding `UPDATE` is DELETED** from `matcher_integration_test.go`. It was the
    symptom: the only way to reach the reader while the chain was broken. Once the chain exists,
    seeding the mirror by hand is a fixture that lies — it would keep passing if S5B regressed.
  - **The rule is CASE-INSENSITIVE, implemented by canonicalizing on WRITE.** `TrimSpace + ToUpper`
    where `NormalizedRow` is born, and the same at the Sankhya writer, for BOTH columns. One slice
    owns write-canonicalization for both sources.
    - Rationale, so nobody re-litigates it: folding at each comparison is guard-at-the-caller, the
      pattern that has already failed twice in this chip. Canonicalizing once at the write seam is
      by-construction. And Oracle is not immune by nature — `AD_ECOMMERCE` is `VARCHAR2` with no
      `CHECK` (measured in DR-3); "everything is uppercase today" is state, not a guarantee.
  - **Canonicalizing case does NOT collapse the tri-state.** An out-of-domain value (`"SIM"`) still
    flows through as honest-unknown. Upper-casing is about `r` vs `R`, never about what the value
    means. Assert this explicitly — it is the criterion most likely to be quietly widened.
  - `matchesSellableAssortment` keeps comparing against canonical values and does NOT fold; one
    comment line says why (the write canonicalizes).
- `validation_kind`: `integration`
- `complexity`: `standard`
- `open_questions`: `[]`

### S6 — Live Oracle catalog predicate and live count

- `id`: `S6-ORACLE-CATALOG`
- `goal`: Apply the configured rule to live list/search queries and calculate live N/M from the same predicate.
- `write_set`:
  - `apps/server_core/internal/modules/internal_read/ports/catalog_page.go`
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/catalog_page.go`
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/catalog_page_test.go`
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go` — **added D-122 by hub
    ruling.** Decision (b) was ratified AFTER this card's write-set froze, so it created work with
    no owner. `reader.go:156-158` raises `QualityMissingStock` from the same LEFT JOIN NULL that
    `catalog_page.go:233-235` does. Shipping only the second makes `/catalogo` render `0` while
    `EstoqueTab` renders `—` for the same product of the same ERP — the two-readers divergence
    that decision (b) exists to kill (:997-999). A slice may not ship half a decision whose other
    half creates the defect the decision names. Binding condition 4 (QA sees both screens) already
    presupposed this file; now the write-set carries what the condition presupposes.
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/reader_test.go`
  - `apps/server_core/internal/modules/internal_read/domain/internal_stock.go` — **added D-122 by
    hub ruling 2 of 4. The default policy CHANGES.** Authority is the operator's ratified verbatim
    in EMENDA A7: asked whether the outlet cut applied to the filter or also to the number on
    screen, the answer was **"os dois"**. `DefaultSellableStockPolicy():18-26` becomes
    `LocationIDs: []int{10101, 10102}` and `ExcludedLocationIDs: nil`. The `ExcludedLocationIDs`
    FIELD stays (S7 and future policies use it); only this default's value goes, because 10108 was
    already outside a whitelist of `{10101}` — excluding it was a dead predicate, and a dead
    predicate reads as protection that is not there.
  - `apps/server_core/internal/modules/internal_read/domain/contract_test.go` — **added D-122**,
    the free must-fail for the line above: `:21-25` asserts `LocationIDs == []int{10101}` and
    `ExcludedLocationIDs == []int{10108}` BY VALUE, so it goes RED naming both values the moment
    the policy moves. Update the expectation; do not weaken the assertion into presence.
  - `apps/server_core/internal/modules/internal_read/adapters/fake/reader_test.go` — **added
    D-122**, one assertion re-pointed. `:38` loops `ExcludedLocationIDs` looking for 10108 and
    fails with "expected showroom location 10108 to stay excluded"; against an empty list it
    fails. The intent (10108 does not sell) survives the change and must keep a test: re-point it
    at the WHITELIST — 10108 ∉ `LocationIDs` — which is the stronger statement anyway, since a
    whitelist excludes 10108 and every location invented tomorrow. **Do not delete the loop.**
  - ONE test file in the `profitability` package (new, or a table extension of an existing one) —
    **added D-122 by hub ruling**, additive and TEST-ONLY. `profitability` PRODUCTION files stay
    out of the write-set, which is what makes binding condition 3 (do not redesign the gate)
    verifiable from the write-set itself rather than by promise.
- `failing_test_first`: `TestCatalogPageSellableAssortmentDefaultsAndScreenEscape` in `apps/server_core/internal/modules/internal_read/adapters/oracle/catalog_page_test.go`
- `done_criteria`:
  - List and search use one shared predicate builder.
  - Enabled clauses are (LIVE side — see the A-14 asymmetry table below; the mirror side in S7
    is deliberately DIFFERENT and a worker who "aligns" the two breaks one of them):
    - `p.USOPROD = 'R'`
    - `NVL(stock.sellable_qty, 0) > 0`
    - `NVL(p.AD_ECOMMERCE, 'X') <> 'N'`

    These three superseded the clauses this card carried until D-122. The dead forms were
    `(p.USOPROD IS NULL OR ...)`, `(stock.sellable_qty IS NULL OR ...)` and
    `(p.AD_ECOMMERCE ... = 'S')`. They are revoked by A-14 (DR-2: the live query read all of
    TGFEST, so absence inside the cut is KNOWN ZERO, not unknown) and by the DR-3 revision
    (only an explicit `'N'` is an assertion; strict `= 'S'` would cut 2.923 → 442).
  - `include_all` omits only those rule predicates; active/cursor/search constraints remain.
  - Count query uses the same stock CTE and predicate builder and returns exact N/M.
  - **The live predicate FOLDS CASE in SQL (D-122, case ruling).** Both comparisons wrap the column:
    `UPPER(TRIM(p.USOPROD)) = 'R'` and `NVL(UPPER(TRIM(p.AD_ECOMMERCE)), 'X') <> 'N'`. This is the
    ONE place a fold is correct, and the reason is that it is the one place we do not own the
    write: S5B canonicalizes what WE store, and nothing canonicalizes what Oracle holds
    (`AD_ECOMMERCE` is `VARCHAR2` with no `CHECK` — measured in DR-3, so uppercase today is state,
    not a guarantee). Without the fold the live side would cut exactly what the mirror accepts, and
    DR-1 forbids the rule depending on the active source. Index cost is irrelevant: these queries
    scan ~10k rows either way.
  - **The outlet location is in the cut, in BOTH readers (D-122, hub ruling 2).** The page's stock
    CTE (`catalog_page.go:131-137`) goes from `est.CODLOCAL = 10101` to `IN (10101, 10102)`, and
    `DefaultSellableStockPolicy` moves with it. Comment both codes by NAME — 10101 = `1_REVENDA`,
    10102 = `2_OUTLET`. It stays a WHITELIST: a location created tomorrow must not start selling by
    itself. Measured cut `CODEMP(1,2) ∧ CODLOCAL(10101,10102)` = 2.923 products.
    The price CTE (`:146`, `e.CODLOCAL = 10101`) and the cost CTE (`:159`, `c.CODEMP = 1`) are a
    DIFFERENT fact with a different rule and are NOT touched — report only.
  - **Binding condition 2 gains the outlet-only case (D-122):** a product whose stock lives only in
    10102 was invisible before the pin and must now be born sellable. Assert it by value.
  - **`buildNotIntListClause` returns `""` for an empty list (D-122, hub ruling 3 of 4).**
    `ExcludedLocationIDs: nil` makes every caller pass an empty slice, and the helper currently
    emits ` AND est.CODLOCAL NOT IN ()` — ORA-00936, invisible to every test at this seat because
    no seat here reaches Oracle. Two callers exist: `reader.go:503` guards the call with
    `if len(policy.ExcludedLocationIDs) > 0`, `stock_batch_reader.go:115` does NOT. The fix goes in
    the HELPER, not the caller — guard-at-the-caller is exactly the pattern that just failed here,
    two callers and one guard. By-construction beats by-vigilance. Consequence: `stock_batch_reader.go`
    needs NO production edit and stays out of the write-set; the `reader.go:503` guard STAYS (a
    redundant guard is not a defect); a test names empty-list → clause ABSENT, and if it is written
    in the batch reader's `_test.go` that TEST file is an additive extension, production untouched.
  - **Do NOT "fix" `buildIntListClause` for symmetry.** Measured at this seat, both of its callers
    are non-empty by construction: `buildSellableStockQuery:489-491` returns
    `ReadErrorUnsupportedQuery` for an empty company or location list, and
    `buildStockBatchQuery:103` builds from `DefaultSellableStockPolicy()`, a constant that is never
    empty. A silently-permitted `IN ()` there would HIDE an empty policy that should blow up — an
    empty whitelist means "sell nothing", and an inclusion list that quietly disappears means "sell
    everything". The asymmetry is the correct semantics and goes in the helper's comment, together
    with what a future caller must preserve: the batch reader's protection is a property of a
    VALUE, not a guard, so a `buildStockBatchQuery` that ever accepts a caller-supplied policy has
    to acquire the edge check `buildSellableStockQuery` already has.
  - Pagination caller tests still prove a gapless limit+1 chain after filtering.
  - **The emitted fact follows the predicate (decision (b), D-122).** Both Oracle sites stop
    raising `QualityMissingStock` for a NULL that the pinned query PROVES is zero: after the pin
    the query read all of TGFEST, so a LEFT JOIN NULL is known-zero, and the flag would be lying
    about itself. `catalog_page.go:233-235` and `reader.go:156-158` emit quantity `0` with no
    flag. `missing_stock` is NOT deleted — it stays reachable on the genuine-unknown path, which
    is the producer side and belongs to S7 (see below).
  - **Binding condition 2 — the must-fail lives in the `profitability` package** and proves BOTH
    sides of the boundary, because the outcome decision (b) names is the GATE's behaviour, and
    observable behaviour is read on the path that PRODUCES it (`service.go:456`). Proving only at
    the fact level would leave the fact→gate link resting on code reading, not on a test:
    - fact with quantity `0` and quality `Complete` → the service COMPUTES (today it blocks);
    - fact with `missing_stock` → the service still BLOCKS. This is condition 1's reachability
      seen from the CONSUMER, obtained free in the same test.
    Assert the VALUES (`0`, `Complete`) and the gate's outcome, never presence. The must-fail
    re-injects the old behaviour (NULL + `missing_stock` for the 10108-only product) and must die
    NAMING the value — a mutation that kills by panicking has not proved the test discriminates
    the value (S4 rule).
  - Sweep the surrounding copy (condition 4): any screen string saying "sem dado" / "incompleto"
    for this case is now false and goes with it.
- `validation_kind`: `unit`
- `commands`:

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\apps\server_core'
$env:GOCACHE = (Join-Path (Get-Location) '.gocache')
$env:GOMODCACHE = (Join-Path (Get-Location) '.gomodcache')
go test ./internal/modules/internal_read/adapters/oracle -run '^TestCatalogPageSellableAssortmentDefaultsAndScreenEscape$|^TestCatalogAssortmentCountUsesThePagePredicate$|^TestCatalogPageCursorChainIsGaplessAndNonOverlapping$' -count=1 -v
```

- `expected_artifacts`:
  - RED/GREEN output
  - Captured SQL and bind lists
  - Exact count assertion
- `complexity`: `complex`
- `open_questions`: `[]`

### S7 — Mirror catalog parity and golden fixture

- `id`: `S7-MIRROR-CATALOG`
- `goal`: Make upload-source catalog paging and counts obey the same rule as link generation and Oracle.
- `write_set`:
  - `apps/server_core/internal/modules/erp_import/ports/repository.go`
  - `apps/server_core/internal/modules/erp_import/application/import_service_test.go`
  - `apps/server_core/internal/modules/erp_import/adapters/postgres/mirror_query_repository.go`
  - `apps/server_core/internal/modules/erp_import/adapters/postgres/mirror_repository_integration_test.go`
  - `apps/server_core/internal/modules/erp_import/adapters/postgres/mirror_search_integration_test.go`
  - `apps/server_core/internal/modules/erp_import/adapters/internalread/reader.go`
  - `apps/server_core/internal/modules/erp_import/adapters/internalread/reader_test.go`
  - `apps/server_core/internal/modules/erp_import/adapters/internalread/source_contract_test.go`
- `failing_test_first`: `TestMirrorCatalogSellableAssortmentGoldenDefaults` in `apps/server_core/internal/modules/erp_import/adapters/postgres/mirror_repository_integration_test.go`
- `done_criteria`:
  - The four-product fixture returns exactly `a,d` and counts exactly `2,4`.
  - Disabling stock returns `a,c,d`; `c` retains numeric stock `0`.
  - Nil `usoprod` passes.
  - `include_all` returns all four.
  - **EAN collisions are counted AFTER the cut (D-122, VC-4 amended @54ff204 on main).** A twin that
    the rule excludes must not leave the survivor flagged ambiguous. Today the count comes from
    `MirrorEANCollisionCounts` (`mirror_query_repository.go:166`), a SQL aggregate over the whole
    mirror that knows nothing about the policy, and it is read at
    `erp_import/adapters/internalread/reader.go:184` → `:203` → `:568`, all BEFORE the assortment
    filter runs. **The two ambiguity mechanisms inside that one function must agree**: the other
    one (`:219`, `len(results) > 1`) is already computed after the cut, so the function currently
    contradicts itself.
    Why it is contract and not a nit: under D-121 an ambiguous candidate goes to manual
    confirmation instead of auto-approving, so cutting junk makes auto-linking WORSE for the
    survivor — the inverse of this chip's purpose.
    Must-fail in BOTH directions: rule active → the survivor auto-approves under D-121; rule
    disabled → it returns to ambiguous.
    Write-set check done at registration time, so it is not re-argued at dispatch: this criterion
    needs `mirror_query_repository.go` (the `MirrorEANCollisionCounts` aggregate) and
    `internalread/reader.go` (the three read sites), and **both are already in S7's write-set
    above**. No extension is required — unlike the three D-122 rulings before it, this one landed
    on a card that could already carry it.
  - **The mirror predicate does NOT fold case (D-122, case ruling).** S5B canonicalizes on write,
    so mirror storage is already `TrimSpace + ToUpper`; a fold here would be the second guard on a
    value that is canonical by construction, and it would hide a regression in S5B rather than
    surface it. Rows written BEFORE S5B resolve themselves on the next import/sync (the live
    Sankhya sync is a hub step post-merge). If this slice's P4 MEASURES a surviving dirty row,
    cleanup gets decided then — on the measurement, not in advance.
  - Search and pagination apply the rule at SQL level before `LIMIT`, avoiding short or skipped pages.
  - Counts are tenant- and source-predicated and use no cache.
  - **Binding condition 1 of decision (b), PRODUCER side (D-122 hub ruling).** Prove by test that
    `missing_stock` stays reachable on the path with a genuine unknown — an upload whose
    spreadsheet lacks the column (A-8) — at `erp_import/adapters/internalread/reader.go:211-214`.
    If no path raises it, report that instead: it then becomes dead code, the next reviewer
    deletes it with good reason, and the signal is lost for good. This site sits in S5's write-set
    and was NOT pulled into S6, which would have manufactured the collision the matrix exists to
    prevent. S6 proves reachability from the CONSUMER side; this proves it from the PRODUCER
    side; neither half closes condition 1 alone.
- `validation_kind`: `integration`
- `commands`:

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\apps\server_core'
$env:GOCACHE = (Join-Path (Get-Location) '.gocache')
$env:GOMODCACHE = (Join-Path (Get-Location) '.gomodcache')
go test -tags=integration ./internal/modules/erp_import/adapters/postgres -run '^TestMirrorCatalogSellableAssortmentGoldenDefaults$' -count=1 -v
go test ./internal/modules/erp_import/adapters/internalread -run 'Catalog|SellableAssortment' -count=1 -v
```

- `expected_artifacts`:
  - Golden fixture SQL
  - `Vendáveis 2 de 4` expected-value record
  - RED/GREEN output with no skips
- `complexity`: `complex`
- `open_questions`:
  - `Blocked on S0/DR-1 HUB approval for mirror catalog parity.`

### S8 — OpenAPI and SDK in one commit

- `id`: `S8-CONTRACT-SDK`
- `goal`: Land all additive HTTP contract and SDK client changes atomically.
- `write_set`:
  - `contracts/api/marketplace-central.openapi.yaml`
  - `packages/sdk-runtime/src/activeSource.ts`
  - `packages/sdk-runtime/src/activeSource.test.ts`
  - `packages/sdk-runtime/src/index.ts`
  - `packages/sdk-runtime/src/index.test.ts`
- `failing_test_first`: `exposes the sellable-assortment paths, schemas, options, and exact client URLs` in `packages/sdk-runtime/src/activeSource.test.ts`
- `done_criteria`:
  - OpenAPI defines both config operations, counts operation, response schemas, and `include_all` on list/search.
  - SDK exposes:
    - `SellableAssortmentConfig`
    - `SetSellableAssortmentRequest`
    - `CatalogAssortmentCounts`
    - `getSellableAssortment`
    - `setSellableAssortment`
    - `getCatalogAssortmentCounts`
    - `include_all` in catalog list/search options
  - URL tests assert exact paths and query strings.
  - Spec and SDK are one commit.
- `validation_kind`: `fe-unit`
- `commands`:

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\packages\sdk-runtime'
npx --no-install vitest run
```

Card command corrected before dispatch (measured, not guessed). The original ran from `apps/web`
and passed the two sdk-runtime files as filters. `apps/web/vitest.config.ts` includes `src/**`,
`feature-products/src/CatalogPage.test.tsx`, `web-query`, and `ui` — **not** `sdk-runtime`, so the
command prints `No test files found, exiting with code 1`. The SDK owns its own
`packages/sdk-runtime/vitest.config.ts` (`include: src/**/*.test.ts`) and its own `test` script;
run from there it is 5 files / 77 tests green at the base. A filter that matches nothing is the
FE twin of the fully-skipped Go lane: it exits non-zero here, but any worker tempted to "fix" it
by dropping the filter would have silently run a different suite.

Recorded for P5, not for S8: the root `package.json` `test` script is
`npm run test --workspace @marketplace-central/web`, so the sdk-runtime suite is NOT part of the
root ladder. Escalated; **hub ruling, VC-7 amended @`ed1b4183`**: the P5 ladder runs the
sdk-runtime suite EXPLICITLY from the package directory, and VC-7 now names TWO vitest lanes,
each counted per line — workspace `@marketplace-central/web` AND `packages/sdk-runtime`. "Green at
the root" under the current script would have been vacuous by half: the contract slice's own suite
could never enter it. The root `package.json` belongs to NO slice's write-set — it is an ownerless
shared seam, so the hub keeps it and fixes the script POST-merge; editing it mid-flight would drift
this tree, and the ladder must use the explicit commands regardless. The lane rule is now law:
`No test files found` with exit 1 is a filter that matched nothing, never a green, and the
correction is the directory, never a looser filter.

- `expected_artifacts`:
  - Contract/SDK commit SHA
  - RED/GREEN Vitest output
  - OpenAPI-to-SDK symbol checklist
- `complexity`: `standard`
- `open_questions`: `[]`

### S9 — Catalog routing, count transport, and composition

- `id`: `S9-CATALOG-HTTP`
- `goal`: Route source-specific counts and expose filtered/list/count behavior through catalog HTTP.
- `write_set`:
  - `apps/server_core/internal/modules/internal_read/adapters/routing/reader.go`
  - `apps/server_core/internal/modules/internal_read/adapters/routing/reader_test.go`
  - `apps/server_core/internal/modules/internal_read/application/service.go`
  - `apps/server_core/internal/modules/internal_read/application/service_test.go`
  - `apps/server_core/internal/modules/internal_read/adapters/cache/cache.go`
  - `apps/server_core/internal/modules/internal_read/adapters/cache/cache_test.go`
  - `apps/server_core/internal/modules/internal_read/observability/timing.go`
  - `apps/server_core/internal/modules/internal_read/observability/timing_test.go`
  - `apps/server_core/internal/modules/catalog/transport/http_handler.go`
  - `apps/server_core/internal/modules/catalog/transport/http_handler_test.go`
  - `apps/server_core/internal/composition/root.go`
  - `apps/server_core/internal/composition/root_test.go`
  - `apps/server_core/internal/modules/internal_read/ports/catalog_page.go` (A-17 extension)
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/catalog_page.go` (A-17 extension: predicate AND count query)
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/catalog_page_test.go` (A-17 extension)
  - `apps/server_core/internal/modules/erp_import/adapters/internalread/reader.go` (A-17 extension: the mirror `catalogPage`)
  - `apps/server_core/internal/modules/erp_import/adapters/internalread/reader_test.go` (A-17 extension)
- `failing_test_first`: `TestHandlerCatalogDefaultsFilteredAndIncludeAllIsRequestLocal` in `apps/server_core/internal/modules/catalog/transport/http_handler_test.go`
- `done_criteria`:
  - `routing.Reader` resolves source once and delegates count to the same reader as the page.
  - Service forwards the count capability without fallback.
  - `GET /catalog/products/counts` is an interactive route.
  - Missing `include_all` pins filtered mode; exact `true` pins all; malformed values return 400.
  - **A-17 — the tenant's three toggles reach the catalog, resolved ONCE at the routing seam.**
    Today nothing does: the only non-test consumer of the stored policy is `routing/matcher.go:45-48`,
    and both catalog readers hardcode the rule (`defaultSellableAssortment()` on the mirror,
    `catalogAssortmentPredicate(includeAll bool)` on Oracle, whose COUNT query takes no option at
    all). VC-3 (badge with `only_em_estoque` off) and VC-2 (counter running the SAME rule) cannot
    pass until this lands. Resolve the policy where `matcher.go` already resolves it — one
    producer, N consumers — and hand the VALUE to page and count alike.
  - **`IncludeAll` is removed from the port.** A bool beside a policy is two mechanisms that must
    agree, which is F-1 applied to ourselves. "Ver todos" and `CatalogProductFactsByIDs` pass an
    all-inclusive policy built by a NAMED domain constructor — a caller never assembles a
    zero-value by hand. The default for a tenant row that is absent lives at the `tenant_config`
    load seam, one place; `defaultSellableAssortment()` dies. Page and count receive the same
    value. Plumbing shape (context pin like linking, or an options field) is S9's call, defended
    in its evidence — but ONE mechanism, resolved at the routing seam.
  - **`include_all` stops at the transport seam (hub ruling on the A-17 layering).** The wire
    parameter STAYS in the HTTP contract — "ver todos" is a per-request choice of the screen, and
    S8 already published it. The condition: S9 builds that seam so the parameter RESOLVES there,
    through the named domain constructor, into the all-inclusive policy — and only the POLICY
    crosses the port. `include_all` never travels ALONGSIDE the policy into the service or the
    reader. Otherwise the two mechanisms that must agree come back through the back door one layer
    up, which is the exact shape criterion above kills at the port. Wire and port are different
    layers; the transport seam is where the wire becomes a value.
  - **Must-fail at contract grade:** through the reader COMPOSED in `root.go`, flipping a toggle on
    the REAL `tenant_config` row moves page AND count together; reverting the threading makes the
    test fail naming the value. Mirror side on the integration lane; Oracle side by query-text
    assertion on the unit lane (the S6 form).
  - `?ids=` remains all-products.
  - Root wires the same routed service into both page and count handler fields.
  - Root test asserts all literal paths and non-nil count wiring.
  - EVERY seat in the live chain composed at `root.go:449→453→479` — `cache.CatalogPageReader`,
    `observability.TimingReader`, `routing.Reader`, `application.Service` — implements
    `ports.CatalogAssortmentReader` AND carries a compile-time
    `var _ ports.CatalogAssortmentReader = ...`. Measured at S6 close: the port had 1 assert
    against the sibling port's 5. A seat that skips it does not fail to build — `routing/reader.go:151`
    turns the failed runtime assertion into `ReadErrorSourceUnavailable`, a 503 on the screen.
    This is the CHIP-M02 catalog-503 defect, and the condition is written at
    `internal_read/ports/catalog_page.go:117`.
  - The arrival test reads the count THROUGH the reader composed in `composition/root.go`, not off
    the Oracle reader directly. Asserting the value leaves its source has already lost the same
    class of fact three times in this chip (S5B's two hops, A-16's third site).
  - No runtime `.(ports.CatalogAssortmentReader)` with a fallback at the HTTP seam. A missing seat
    must be a build error, never a quiet wrong answer.
- `validation_kind`: `unit`
- `commands`:

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\apps\server_core'
$env:GOCACHE = (Join-Path (Get-Location) '.gocache')
$env:GOMODCACHE = (Join-Path (Get-Location) '.gomodcache')
go test ./internal/modules/internal_read/application ./internal/modules/internal_read/adapters/routing ./internal/modules/catalog/transport ./internal/composition -run 'Catalog|Assortment|Root' -count=1 -v
```

- `expected_artifacts`:
  - Handler RED/GREEN output
  - Router source-delegation output
  - Composition route inventory
- `complexity`: `standard`
- `open_questions`: `[]`

### S10 — Configuration HTTP and governance registration

- `id`: `S10-CONFIG-HTTP`
- `goal`: Expose database-backed toggle GET/PUT and register tenant_config governance ownership.
- `write_set`:
  - `apps/server_core/internal/modules/tenant_config/transport/http_handler.go`
  - `apps/server_core/internal/modules/tenant_config/transport/http_handler_test.go`
  - `contracts/governance/modules.json`
- `failing_test_first`: `TestHandlerSellableAssortmentPutThenGetReturnsStoredValues` in `apps/server_core/internal/modules/tenant_config/transport/http_handler_test.go`
- `done_criteria`:
  - Handler registers GET/PUT `/config/sellable-assortment` as interactive.
  - PUT requires three booleans and returns stored exact values.
  - No browser-storage concept appears in server code.
  - `modules.json` adds `tenant_config`, root `apps/server_core/internal/modules/tenant_config`, OpenAPI prefix `/config`, dependency `erp_import`.
- `validation_kind`: `unit`
- `commands`:

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\apps\server_core'
$env:GOCACHE = (Join-Path (Get-Location) '.gocache')
$env:GOMODCACHE = (Join-Path (Get-Location) '.gomodcache')
go test ./internal/modules/tenant_config/transport -run '^TestHandlerSellableAssortmentPutThenGetReturnsStoredValues$' -count=1 -v
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel'
npm run harness:governance -- -BaseSha 554788d576d04b21719f4a17e4702dd0f0aff4e1
```

- `expected_artifacts`:
  - Handler RED/GREEN output
  - Governance output with no unowned `/config` prefix
- `complexity`: `standard`
- `open_questions`: `[]`

### S11 — `/integracoes` card

- `id`: `S11-INTEGRACOES-FE`
- `goal`: Render and persist the three toggles with a live source-truthful “Resultado: N de M produtos”.
- `write_set`:
  - `packages/web-query/src/activeSource.ts`
  - `packages/web-query/src/index.ts`
  - `apps/web/src/pages/integracoes/IntegracoesPage.tsx`
  - `apps/web/src/pages/integracoes/IntegracoesPage.test.tsx`
- `failing_test_first`: `renders Sortimento vendável, persists all toggles, and shows the exact live count` in `apps/web/src/pages/integracoes/IntegracoesPage.test.tsx`
- `done_criteria`:
  - Sibling card uses the existing card and option tokens from [`IntegracoesPage.tsx:307-339`](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-vendavel/apps/web/src/pages/integracoes/IntegracoesPage.tsx:307).
  - Labels and errors are pt-BR.
  - All three controlled toggles initialize from GET and mutate the database through PUT.
  - Success updates config cache and invalidates catalog/count queries.
  - Exact text is `Resultado: 2 de 4 produtos` in the fixture.
  - Tests assert `localStorage` and `sessionStorage` are never written.
- `validation_kind`: `fe-unit`
- `commands`:

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\apps\web'
npx --no-install vitest run src/pages/integracoes/IntegracoesPage.test.tsx
```

- `expected_artifacts`:
  - RED/GREEN DOM output
  - Exact mocked GET/PUT payloads
  - Zero storage-write assertion
- `complexity`: `standard`
- `open_questions`: `[]`

### S12 — `/catalogo` filtered view and escape

- `id`: `S12-CATALOGO-FE`
- `goal`: Open filtered, show “Vendáveis N de M”, and provide a screen-only all-products escape with honest stock badges.
- `write_set`:
  - `packages/feature-products/src/catalogQueries.ts`
  - `packages/feature-products/src/CatalogPage.tsx`
  - `packages/feature-products/src/CatalogPage.test.tsx`
  - `packages/web-query/src/index.ts`
  - `packages/web-query/src/index.test.ts`
  - `apps/web/src/app/AppRouter.tsx`
  - `apps/web/src/app/AppRouter.test.tsx`
- `failing_test_first`: `opens with Vendáveis 2 de 4 and ver todos never mutates tenant config` in `packages/feature-products/src/CatalogPage.test.tsx`
- `done_criteria`:
  - Initial list/search requests omit `include_all` or send false.
  - Counter chip is exactly `Vendáveis 2 de 4`.
  - `Ver todos` refetches with `include_all=true` and never calls `setSellableAssortment`.
  - Returning to filtered mode uses a distinct query key.
  - A zero-stock product returned after `only_em_estoque` is disabled remains in the table with badge `Sem estoque`.
  - AppRouter test proves `/catalogo` mounts the real component with the client and source cache partition.
- `validation_kind`: `fe-unit`
- `commands`:

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\apps\web'
npx --no-install vitest run ../../packages/feature-products/src/CatalogPage.test.tsx src/app/AppRouter.test.tsx ../../packages/web-query/src/index.test.ts
```

- `expected_artifacts`:
  - RED/GREEN DOM output
  - Exact request sequence for filtered/all/filtered
  - Zero config-mutation assertion
- `complexity`: `standard`
- `open_questions`: `[]`

### S13 — Full ladder and HUB-only live proof

- `id`: `S13-VERIFY`
- `goal`: Run the full static/test ladder and hand live Oracle/browser validation to HUB.
- `write_set`: `[]`
- `failing_test_first`: `N/A — lane aggregation; all named RED tokens must already exist`
- `done_criteria`:
  - Go build/test ladder green.
  - BOTH required FE Vitest lanes green and counted per line (VC-7 @`ed1b4183`): `apps/web`
    (workspace `@marketplace-central/web`) AND `packages/sdk-runtime`. Neither may print
    `No test files found` — that is a filter matching nothing, not a green.
  - Root typecheck green.
  - Governance green against the full base SHA.
  - Integration report proves named tests ran and did not skip.
  - HUB supplies live sync SQL and browser/network evidence.
- `validation_kind`: `lane`
- `commands`:

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\apps\server_core'
$env:GOCACHE = (Join-Path (Get-Location) '.gocache')
$env:GOMODCACHE = (Join-Path (Get-Location) '.gomodcache')
go build ./...
go test ./... -count=1

Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel'
npx --no-install tsc --noEmit -p apps/web/tsconfig.json

Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\apps\web'
npx --no-install vitest run

Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel\packages\sdk-runtime'
npx --no-install vitest run

Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel'
npm run harness:integration
npm run harness:governance -- -BaseSha 554788d576d04b21719f4a17e4702dd0f0aff4e1
```

- `expected_artifacts`:
  - Go build/test outputs
  - TSC and Vitest outputs
  - Integration `=== RUN`/skip counts
  - Governance report
  - HUB live-sync and browser-drive bundle
- `complexity`: `standard`
- `open_questions`:
  - `HUB must run the live Oracle sync and browser drive; the chip must not boot the stack, bind ports, or read .env files.`

---

## Per-slice write-set DAG

```text
S0 decision ───────────────────────────────┐
                                          v
S1 schema/domain ──> S2 config DB ───────> S5 matcher
       │                  │
       ├──> S3 XLSX       └──────────────> S6 Oracle catalog
       ├──> S4 sync
       └─────────────────────────────────> S7 mirror catalog (after S0 approval)

S6 + S7 ──> S8 contract+SDK ──> S9 catalog HTTP/routing
S2 ───────────────────────────> S10 config HTTP
S8 + S10 ─────────────────────> S11 Integracoes FE
S8 + S9 + S11 ────────────────> S12 Catalog FE
S3 + S4 + S5 + S6 + S7 + S9 + S10 + S11 + S12 ──> S13
```

Shared-file ordering constraints:

- `domain/mirror.go`: S1 before S5/S7.
- `erp_import/adapters/internalread/reader.go`: S5 before S7.
- `packages/web-query/src/index.ts`: S11 before S12.
- `root.go`: only S9 owns it.
- `cache/cache.go`, `observability/timing.go`: only S9 owns them. Added to its write-set at S6
  close — they are two of the four decorator seats the new optional port must reach.
- `internal_read/ports/catalog_page.go`, `oracle/catalog_page.go`, `erp_import/adapters/internalread/reader.go`:
  only S9 owns them from A-17 onward. S6 and S7 are closed and do not reopen — the hub ruled one
  patch, one owner, because a corrective slice would write the same port signature S9 rewrites.
- OpenAPI and SDK: only S8 owns them and they land together.
- `runner_test.go`: only S1 owns both count edits.
- No two slices execute concurrently; this chip has one writer.

---

## Contract-satisfiability pass

| Surface | Current state | Planned change |
|---|---|---|
| `/catalog/products` | Occupied at OpenAPI lines 338–411; operation `listCatalogProductFacts` already exists at line 341 | Add optional `include_all`; do not replace the operation or response |
| `/catalog/products/search` | Occupied at lines 435–505; operation `searchCatalogProductFacts` at line 438 | Add optional `include_all` |
| `/catalog/products/counts` | Free; no count endpoint exists | Add `GET`, operation `getCatalogAssortmentCounts`, literal route registered before/alongside `{id}` |
| `/config/active-source` | Occupied at lines 3302–3342 | Leave behavior and schemas intact |
| `/config/sellable-assortment` | Free | Add GET/PUT in the existing tenant-config handler family |
| `ActiveSourceConfig` | Occupied at schema line 8166 | Leave source fields intact; use a separate additive assortment schema |
| `SetActiveSourceRequest` | Occupied at line 8185 | Leave intact |
| `SellableAssortmentConfig` | Free | Add required three booleans |
| `SetSellableAssortmentRequest` | Free | Add required three booleans |
| `CatalogAssortmentCounts` | Free | Add required integer `sellable_count`, `total_count`, both `minimum: 0` |
| SDK `getActiveSource` / `setActiveSource` | Occupied at [`index.ts:1896-1898`](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-vendavel/packages/sdk-runtime/src/index.ts:1896) | Leave intact |
| SDK `getSellableAssortment` | Free | Add adjacent to active-source methods |
| SDK `setSellableAssortment` | Free | Add adjacent to active-source methods |
| SDK `getCatalogAssortmentCounts` | Free | Add in the catalog/config domain block |
| `CatalogPageOptions` | Occupied at [`index.ts:239`](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-vendavel/packages/sdk-runtime/src/index.ts:239) | Add optional `include_all?: boolean` |
| SDK list/search methods | Occupied at lines 1851 and 1859 | Forward `include_all` through the existing `catalogQuery` helper |

---

## Verification map

| Criterion | Slice | Named evidence | HUB-only portion |
|---|---|---|---|
| VC-1 | S2, S10, S11, S13 | `TestRepository_SetSellableAssortment_RoundTripPersistsPerTenant`; `TestHandlerSellableAssortmentPutThenGetReturnsStoredValues`; Integracoes persistence test; SQL before/after | Browser clears both storages, reloads, captures DOM and `localStorage.length === 0` |
| VC-2 | S6, S7, S9, S11, S12 | `TestCatalogAssortmentCountUsesThePagePredicate`; golden `2/4`; exact FE text assertions | Live Oracle SQL and DOM comparison, allowing source drift from approximately 3,822 |
| VC-3 | S6, S7, S12 | `TestCatalogPageSellableAssortmentDefaultsAndScreenEscape`; golden fixture; CatalogPage filtered/all/zero-stock test | Browser DOM before/after `Ver todos` and after disabling stock rule |
| VC-4 | S5 | `TestMirrorMatcher_ActiveRevendaRuleControlsCandidateBirth`; `failure_token=test=...`; positive sellable and nil assertions | None beyond optional browser `/vinculos` smoke |
| VC-5 | S4 | `TestSankhyaStockSQLPinsSellableCompanies`; `TestSankhyaSyncMapsSellableAssortmentFields`; mirror writer test | Live sync and post-sync mirror SQL are HUB-only |
| VC-6 | S3 | `TestParserSellableColumnsAreOptionalAndHonestUnknown`; `TestMergeSnapshotPersistsOptionalSellableColumns` | None |
| VC-7 | S13 | `go build ./...`; `go test ./...`; root TSC; BOTH required vitest lanes, counted per line — `cd apps/web && npx --no-install vitest run` AND `cd packages/sdk-runtime && npx --no-install vitest run` (VC-7 amended @`ed1b4183`); integration and governance lanes | Zero-console-error and zero-4xx/5xx browser/network drive on `/integracoes`, `/catalogo`, `/vinculos` |
| Day-1 golden | S7, S12 | `TestMirrorCatalogSellableAssortmentGoldenDefaults`; CatalogPage exact `Vendáveis 2 de 4` | Optional live screenshot, not needed for fixture discharge |

HUB live SQL after sync:

```sql
SELECT
  count(*) FILTER (
    WHERE (usoprod IS NULL OR usoprod = 'R')
      AND (estoque_total IS NULL OR estoque_total > 0)
  ) AS sellable_count,
  count(*) AS total_count,
  count(*) FILTER (WHERE usoprod IS NOT NULL) AS usoprod_known,
  count(*) FILTER (WHERE ad_ecommerce IS NOT NULL) AS ad_ecommerce_known
FROM products_mirror
WHERE tenant_id = :tenant_id
  AND source = 'sankhya'
  AND absent_in_last_snapshot = false;
```

VC-1 SQL:

```sql
SELECT tenant_id, only_revenda, only_em_estoque, only_ecommerce_eligible
FROM active_source
WHERE tenant_id = :tenant_id;
```

---

## Seam-closure checklist

| Seam | Owner |
|---|---|
| Composition root routes page and count through the same `routing.Reader` | S9 `root.go` + `root_test.go` |
| Existing tenant-config repository/handler instance serves both active source and assortment | S10 handler registration; existing root construction at `root.go:478-481` |
| Router tests assert source delegation and no fallback | S9 `routing/reader_test.go` |
| Oracle SQL tests assert bind placeholders and predicate placement | S6 `catalog_page_test.go` |
| Q4 exact company predicate | S4 `sync_test.go` |
| Catalog literal `/counts` route does not fall into `{id}` | S9 handler/root tests |
| App shell mounts real `/catalogo` client/source props | S12 `AppRouter.test.tsx` |
| SDK client methods, not types only | S8 `index.ts` + exact URL tests |
| OpenAPI and SDK same commit | S8 |
| Governance `tenant_config` entry and `/config` ownership | S10 |
| Both migration inventory assertions | S1 |
| Both mirror writers | S3 upload writer; S4 Sankhya writer |
| XLSX without optional columns | S3 |
| Search and list both filtered | S6/S7 caller tests |
| Screen escape does not mutate config | S12 |
| Live Oracle sync | HUB post-merge |
| Live browser/network drive | HUB post-merge |

---

## Must-fail design

### VC-4

Named test:

```text
TestMirrorMatcher_ActiveRevendaRuleControlsCandidateBirth
```

Mutation:

- Remove only the statement in `MirrorMatcher.FindProductsForLinking` that pins `SellableAssortment` into the mirror reader context.
- Keep active-source pinning intact.

Expected isolated RED:

- With database `only_revenda=true`, the `USOPROD='V'` row becomes a candidate, failing its exact empty-candidate assertion.
- The source-routing tests remain green because active-source context still exists.
- Direct reader rule tests remain green because they pin the rule themselves.
- Thus only the MirrorMatcher arm is isolated.

MUST-PASS assertions in the same test:

- `USOPROD='R'`, stock `5` candidate is born.
- nil `USOPROD`, stock `3` candidate is born.
- After database `only_revenda=false`, the `V` candidate is born.

### VC-5

Named test:

```text
TestSankhyaStockSQLPinsSellableCompanies
```

Mutation:

- Delete only `AND CODEMP IN (1, 2)` from `sankhyaStockSQL`.
- Do not alter Q1 or fake stock rows.

Expected isolated RED:

- The exact Q4 SQL assertion fails.
- Sync field-mapping tests remain green because the fake driver does not model extra companies.
- Q1 field tests remain green.
- Therefore the mutation isolates the company-scope defect instead of reddening several mapping arms.

Separate MUST-PASS:

- `TestSankhyaSyncMapsSellableAssortmentFields` asserts a known `USOPROD='R'` row and `AD_ECOMMERCE='S'` reach `mirror.Row` exactly.
- Existing real-zero stock assertion remains `0`, not nil.

---

## Risks and traps

- Both migration counts at `runner_test.go:25` and `:64` must become `71`.
- `0083` and `0084` are both required; a third migration is an escalation.
- `Repository.Set` must not overwrite assortment toggles when changing active source.
- Both mirror writers must carry both new fields:
  - Sankhya `mirror/writer.go`
  - upload `postgres/mirror_repository.go`
- Both mirror readers/scanners must update select order and `Scan` destinations together.
- `NULL <> 'R'` evaluates to NULL, not true. Use `IS NULL OR = 'R'`; never a bare inequality.
- The same applies to `AD_ECOMMERCE` and missing stock.
- Never put nullable-field predicates in a form that silently turns the `LEFT JOIN stock` into an inner exclusion.
- Filtering mirror rows before building the linking TTL cache would make a database toggle stale. Filter after cache lookup at the caller.
- Catalog filters must execute in SQL before `LIMIT + 1`; filtering returned rows in Go breaks pagination and can create false end-of-list states.
- Counts must run against the source selected by `routing.Reader`, not always Oracle and not a constant.
- Search and list are separate callers and both require validation.
- `CatalogProductFactsByIDs` should remain unfiltered to avoid hiding already-linked product detail reads.
- A `Sem estoque` badge is presentation only; it must not become another filter.
- `include_all` belongs to request context/query key only. It must never call PUT or alter `active_source`.
- All three UI toggles must send the complete expected configuration; no absent-boolean coercion.
- No localStorage/sessionStorage fallback or optimistic persistence.
- Existing `erp_source` query documentation is stale relative to database routing; do not reuse it as the assortment switch.
- Integration tests with no `MPC_TEST_DATABASE_URL` skip. A green artifact is valid only with `=== RUN` evidence and recorded skip count.
- Go commands must run from `apps/server_core` with separate absolute `GOCACHE` and `GOMODCACHE`.
- The required FE lane is exactly `Set-Location apps/web` followed by `npx --no-install vitest run`.
- The chip must not boot the stack, bind ports, or inspect `.env*`; live sync and browser validation belong to HUB.

---

# CHIP AMENDMENTS TO THE PLAN (2026-07-29)

The plan above was authored BEFORE the A7 ruling closed (`ADDENDUM-01-codlocal.md`, hub
`local_99feb041`, pack/contract amended on main @ `aee1a222`). Where the two disagree, this
section and ADDENDUM-01 win. Workers read the slice card AND this section.

## A-1 — S4: the Q4 pin is company AND location AND net of reservation

The plan's S4 pins only `AND CODEMP IN (1, 2)`. Ratified `sankhyaStockSQL` becomes:

```sql
SELECT CODPROD, CODLOCAL, SUM(NVL(ESTOQUE, 0) - NVL(RESERVADO, 0)) AS DISPONIVEL
FROM METALPRD.TGFEST
WHERE CODPARC = 0
  AND CODEMP IN (1, 2)
  AND CODLOCAL IN (10101, 10102)
GROUP BY CODPROD, CODLOCAL
```

`GROUP BY CODPROD, CODLOCAL` stays. `CODPARC = 0` stays (defensive, not dead — see ADDENDUM-01).
The selling locations are a named constant carrying the TGFLOC names, never bare literals
sprinkled across queries; the arithmetic `ESTOQUE - RESERVADO` sits at ONE point per site so the
specialist's answer would be a one-line change.

`TestSankhyaStockSQLPinsSellableCompanies` is renamed `TestSankhyaStockSQLPinsSellableCompaniesAndLocations`
and must go RED if **either** predicate leaves the query — deleting `CODEMP IN (1,2)` alone and
deleting `CODLOCAL IN (10101,10102)` alone are TWO separate must-fail runs, both recorded.

Mirror semantics: `estoque_total` now means AVAILABLE SELLABLE stock. The comment at
`sync.go:247-251` asserts `total == sum of children`; it stays true under this filter and gets
rewritten to say what the numbers now mean. If any mirror consumer depends on the old sense
(gross, all companies/locations), that is a `REQUEST` to the hub BEFORE the change.

## A-2 — S6: the live catalog stock CTE gains outlet

`catalog_page.go:131-137` reads `est.CODLOCAL = 10101`; becomes `IN (10101, 10102)` with the same
named constant. This moves the number on `/catalogo` by itself (outlet enters) — the operator
wants that. `price_candidates` at `:146` also pins `e.CODLOCAL = 10101`: price selection is OUT
of scope, do not "align" it while passing through.

## A-3 — VC-2 acceptance number is 2.923, and the delta is an assertion

The plan's VC-2 row says "allowing source drift from approximately 3,822". That number is DEAD
(gross stock, no reservation, no location cut) — delete it, do not soften it. Locked acceptance:
**2.923 distinct sellable CODPROD**. Neighbours, so a wrong screen is diagnosable: 3.508 = forgot
the location cut · 3.277 = forgot the reservation · 3.822 = forgot both.

Count DISTINCT CODPROD, never sum per location (10101 = 2.795, 10102 = 187, overlap is real).

The 585-product delta between the correct filter and the filter-without-CODLOCAL is an
ASSERTION in the fixture, not a note: both counts must come out DIFFERENT, otherwise the test
does not sustain the pin.

## A-4 — consistency is a contract, not a nicety

Mirror, live catalog and the N/M counter use the SAME definition of available. If the three
sites diverge, VC-2 fails by design. That is what makes DR-1 (below) load-bearing rather than
cosmetic.

## A-5 — provenance line for the evidence pack

Every figure in ADDENDUM-01 is cited as: **db-consult MNOS 2026-07-29 via hub** (read-only
COUNT/aggregation, zero raw rows).

## A-6 — DR-1 stays a hub decision, S7 stays undispatchable

S0/DR-1 (mirror-catalog parity for xlsx/catalogo_cliente tenants) went to the hub as a `REQUEST`
with the planner's predicate and cost estimate. S7 is not dispatched until the hub rules. No
slice ships a count that would be false for an upload-source tenant.

## A-7 — DR-1 GRANTED: the mirror catalog path is binding (hub, 2026-07-29, main @ 66f32125)

The hub granted parity with no deferral, and corrected the pack: the rule binds **three** read
paths, not two — live Oracle · **the mirror serving the catalog** · MirrorMatcher. Reason of
record: filtering only Oracle would make a business rule depend on the active SOURCE, which is
the coupling MIS-006 spent seven milestones removing; the mirror exists so that source is an
adapter detail, not semantics.

Three binding conditions on S7:

1. **Filter in SQL before `LIMIT+1`**, and prove pagination with a fixture of MORE THAN ONE
   page carrying a sellable item AFTER the first page's cut. A fixture too small produces a
   vacuous green (CHIP-MERCADO: page-1 truncation invisible to the live drive because the
   worktree DB held < 50 rows).
2. **Measure what the xlsx writer stores for an ABSENT datum — NULL or 0** (see A-8: measured).
   Two distinct fixture cases, never one: known-zero stock is CUT by `only_em_estoque`, unknown
   stock PASSES.
3. **The must-fail lives at the CALLER** — `routing.Reader` with source `xlsx`, the path the
   screen actually uses — and must fail if the predicate leaves ANY of the three sites. A green
   repository test with an unfiltered router is the empty green that cost us B-08.

The mirror-path counter counts **against the mirror**, with the same predicate. A counter
reading one store while the screen renders another is a false number — worse than a 503,
because it does not look like a defect.

VC-2 scoping correction: **2.923 is the LIVE ORACLE number only.** The mirror path has its own
population (sankhya mirror 10.538 rows, xlsx 55). VC-2 discharges by AGREEMENT within each
source — screen == SQL of the same rule against the same store — and only the Oracle path
carries 2.923 as an absolute reference.

## A-8 — measured: the upload writer is already honest about absent stock

Condition 2 answered by measurement, not assumption, at
`erp_import/adapters/postgres/mirror_repository.go:93-99` and `:113-121`:

- `nullableString(row.StockPhysical)` maps an empty cell to `nil`
  (`import_repository.go:111-116`), so an absent `ESTOQUE_FISICO` reaches the stage table as
  NULL, not `0`. The lenient parser is the only path that can omit the column
  (`xlsx/parser.go:33-35`); the strict path still requires it.
- The insert derives `estoque_total` as
  `CASE WHEN stock_physical IS NOT NULL THEN stock_physical::numeric - COALESCE(stock_reserved::numeric, 0) END`
  — physical unknown leaves `estoque_total` NULL; a reported `0` becomes numeric `0`.

So **no fix is required in S3** — both ADR-17 sides already hold on the upload writer. What S3
gains instead is the test that PINS the distinction, because nothing today would catch a
regression that coerced absent to zero: one fixture row with known-zero stock (cut by
`only_em_estoque`) and one with unknown stock (passes). The predicate
`m.estoque_total IS NULL OR m.estoque_total > 0` is only safe because of this measured
behaviour, and the test is what keeps it true.

## A-9 — measured: the Sankhya writer WOULD cancel the pin. The trap is real, and it is in S4.

The hub named a symmetric trap on the Sankhya side. Measured, and **CONFIRMED in the current
code** — this is not a hypothetical:

- `mirror.Row.EstoqueTotal` is `*float64`, and the type comment at `mirror/writer.go:29-30`
  states nil round-trips to SQL NULL — never 0.
- `applyStock` (`oracle/sync.go:258-296`) builds `totals` only from rows the Q4 result returned,
  and writes `rows[codprod].row.EstoqueTotal = &v` **only for products present in `totals`**
  (`:291-294`). A product with no Q4 row keeps nil.
- Therefore, once A-1 pins Q4 to `CODEMP IN (1,2) AND CODLOCAL IN (10101,10102)`, every product
  whose stock exists ONLY in show room / almoxarifado / another company stops producing a Q4
  row and lands in the mirror with `estoque_total = NULL`.
- The sellable predicate is `m.estoque_total IS NULL OR m.estoque_total > 0`. NULL passes.
  **The 585 products the pin just removed are readmitted, and the cut returns to 3.508** while
  every site looks correct in isolation. The only symptom is the number on screen.

### Ruling applied

On the Sankhya path stock is **KNOWN**. The sync read TGFEST; absence of a row inside the
sellable cut means "zero available where we sell", not "unknown". **Known-zero → `0`.** The NULL
of honest-unknown belongs to the spreadsheet that did not carry the column, not to an ERP that
answered. Both sides of ADR-17, on the correct side each.

So `applyStock` sets `EstoqueTotal` to `0` for every product in the base snapshot that the
pinned Q4 did not return — and the comment at `sync.go:247-251` ("A product absent from TGFEST
keeps NULL estoque") becomes FALSE under the pin and gets rewritten, not softened.

### The test that ties S4 to the VC-2 number

`TestSankhyaSyncStoresKnownZeroForProductsOutsideTheSellableCut`: fixture with a product whose
stock exists ONLY in `10108`; sync; assert `estoque_total == 0` — **the value, not presence** —
and assert the product is OUT of the sellable set. Without this test the pin goes green and the
screen lies.

Same class as A-8, mirrored: there the risk was absent becoming zero and emptying the catalog;
here it is zero becoming absent and re-inflating it.

## A-10 — sweep of `estoque_total` consumers (hub condition on A-9), done at the CALLER

Every consumer that branches on `EstoqueTotal == nil`, swept by pattern across `apps/server_core`
and the FE packages, non-test files:

| site | NULL branch today | after `NULL → 0` |
|---|---|---|
| `erp_import/adapters/internalread/reader.go:211-214` (`SellableStock`) | flags `QualityMissingStock`, `Quantity` stays nil, row still returned | real quantity `0`, no flag |
| `erp_import/adapters/internalread/reader.go:428-436` (`catalogFact`) | same — flag on the fact, row still returned | quantity `0`, no flag |
| `erp_import/adapters/internalread/stock_batch_reader.go:77-81` | same — fact returned with the flag | quantity `0`, no flag |
| `packages/feature-products/src/CatalogPage.tsx:30-31` (`factQuantity`) | renders `— (missing_stock)` | renders `0` |

**No consumer treats NULL as "skip this row".** Every one keeps the row and flags the fact, so
the change cannot silently include or exclude products — it changes the CELL from
`— (missing_stock)` to `0`, which is the intended, truer statement. Nothing to REQUEST.

The only other NULL producer feeding these readers is the xlsx path, whose NULL is genuine
unknown (A-8) and is unaffected: the change lands in `applyStock`, on the Sankhya path only.

## A-11 — the two populations must not merge (hub condition on A-9)

Known-zero applies to a product the **base snapshot (Q1) returned** and the pinned Q4 did not.
It does NOT apply to a product missing from the snapshot entirely — that is
`absent_in_last_snapshot`, which has its own ADR-04 semantics (flagged, never deleted,
last-known values preserved). Merging the two would stamp `0` on products the ERP never
reported, manufacturing a fact out of an absence.

The fixture carries THREE stock cases, not two:
1. product in Q1, stock only in `10108` → `estoque_total = 0`, OUT of the sellable set;
2. product in Q1 with available stock in `10101` → real quantity, IN the set;
3. product NOT in the current snapshot → `absent_in_last_snapshot = true`, last-known value
   preserved, **not** rewritten to `0`.

## A-12 — F-01 is superseded, with its trail intact

The F-01 decision recorded at `sync.go:247-251` ("a product absent from TGFEST keeps NULL
estoque; a real 0 balance is distinct and legitimate") was CORRECT while Q4 read all of TGFEST.
It is superseded by the A7 cut on 2026-07-29, because absence now means "zero available where we
sell" rather than "the ERP said nothing". The rewritten comment states the new meaning AND
records that F-01 held until the pin existed — the false prose goes, the trail stays.

## A-13 — sweep of the flag's CONSUMERS (not its producers), and a measured absence

`QualityMissingStock` stops being raised on the Sankhya path under A-9, so the question is who
READS it. Non-test consumers:

| consumer | what it does with the flag |
|---|---|
| `apps/web/src/pages/produto/EstoqueTab.tsx:48` | `hasFisico = stock.value !== null && !quality_flags.includes("missing_stock")` — display gate: value+freshness vs the "importar planilha" CTA |
| `apps/web/src/pages/produto/ProdutoHeader.tsx:97` | renders the flags as badges (passive) |
| `profitability/application/service.go:331,338,456` | `mapInternalQuality(amount, fact.QualityFlags)`, gating at `:456` on `!= InputQualityComplete` — a REAL business gate |

A gate exists, so the question was worth asking. It does not lose its signal under A-9: all three
consume facts through `routing.Reader`, and for a `sankhya` tenant routing serves the LIVE Oracle
reader, not the mirror. The mirror's `source='sankhya'` rows are consumed today by MirrorMatcher,
which does not read quality flags. **Measured absence: no gate stops firing.** Recorded
explicitly — an absence that was measured is worth as much as a presence.

## A-14 — CLOSED: RULING DR-2 = **(b), the fact follows** (hub `local_99feb041`, 2026-07-29)

The hub ruled (b). Reasons, in the weight the hub gave them:

- **(a) keeps a lie ON THE SCREEN, not merely a divergence.** `EstoqueTab` renders the
  "importar planilha" CTA while the flag is up. Under (a), a Sankhya tenant is told to import a
  spreadsheet to repair a datum that is not missing — the ERP answered, the product has zero
  available where we sell, and no spreadsheet fixes that. False copy on an operator surface is
  deleted here, not carried as debt. That alone decides it.
- **Divergence between two readers of the SAME source** is the defect class this mission spent
  seven milestones killing. Same ERP, same product, `—` in one reader and `0` in the other is
  reader-semantics coupling coming back through the side door.
- **The flag would be lying about itself.** `missing_stock` means "I don't know". After the pin
  the query READ all of TGFEST; a NULL from the LEFT JOIN is known-zero — for a product outside
  the cut and for a product with no line at all.

On pricing: a product that starts computing with stock `0` is, by construction, a product with no
stock in the selling locations, which the sellable rule already cuts from the assortment by
default. Its contribution is zero. Nothing new is released; we stop labelling as "incomplete
data" a datum we hold. The hub reports the change to the operator and returns if vetoed — the
chip does not stop for it.

### Four binding conditions

1. **Do not kill the flag — kill it only where the zero is known.** `missing_stock` stays
   reachable on the path with genuine unknown (upload without the column, A-8). Prove it stays
   reachable with a test, or report that it does not. If no path raises it, it becomes dead code
   and the next reviewer deletes it with good reason — and the signal is lost for good.
2. **Must-fail that NAMES the gate's behaviour change:** a product stocked only in `10108` →
   fact with quantity `0` and quality `Complete`, and `profitability` COMPUTES instead of
   blocking. Assert the VALUES (`0`, `Complete`) and the gate's outcome — never presence.
3. **Scope locked:** change what the fact EMITS and the quality mapping that follows from it.
   No redesign of `profitability`. Any need to touch its logic beyond that is a `REQUEST`.
4. **QA sees both screens:** `EstoqueTab` (the "importar planilha" CTA → `0`) and the catalog.
   Sweep the surrounding copy: any string saying "sem dado" / "incompleto" for this case is now
   false and goes with it.

S6 is unblocked. The predicate half stands as recorded below: on Oracle the cut uses
`NVL(stock.sellable_qty, 0) > 0`, never honest-unknown — so S6's planned
`(stock.sellable_qty IS NULL OR stock.sellable_qty > 0)` clause is REPLACED, and the plan's
`(p.USOPROD IS NULL OR p.USOPROD = 'R')` gets the same treatment for the same reason: Q1 read
`TGFPRO`, so `USOPROD` is known there too.

**The mirror path does NOT inherit this.** On the mirror, `usoprod IS NULL OR usoprod = 'R'` and
`estoque_total IS NULL OR estoque_total > 0` STAY, because an xlsx tenant whose file lacks the
column has genuine unknown (A-8) and must not be cut by a fact nobody stated. Same predicate
name, two populations, two correct forms — S6 (live) is strict, S7 (mirror) keeps the NULL
branch. A worker who "aligns" them has broken one of the two.

### `AD_ECOMMERCE` — RULING DR-3: honest-unknown on BOTH sides (hub, 2026-07-29)

The chip asked whether the live path should read `p.AD_ECOMMERCE = 'S'` strictly, by the same
reasoning as (b). **The hub ruled no. Do not extend the strict form to this column.**

**Superseded in its FORM by the live measurement (hub, pack @ `0e7145ee`) — the principle below
stands, the SQL is the specialist's.** Measured in METALPRD: NULL 6.939 · `'N'` 2.993 ·
`'S'` 606. Inside the 2.923 sellable: `'N'` 1.595 · NULL 887 · `'S'` 442. The column is
`VARCHAR2(10)` nullable, **no default and no CHECK** — S/N/NULL is a convention, not a database
guarantee. Binding form:

```sql
AND NVL(AD_ECOMMERCE, 'X') <> 'N'                        -- live
AND (m.ad_ecommerce IS NULL OR m.ad_ecommerce <> 'N')    -- mirror
```

Two things the measurement decides:

1. **The principle holds with room to spare.** Strict `= 'S'` would drop 2.923 → 442 — **85% of
   the assortment** — and the operator would read that as "the feature broke". The tri-state is
   real: 2.993 products carry an EXPLICIT `'N'`, so the "no" is asserted, not inferred from blank.
2. **`IS NULL OR = 'S'` is wrong in the tail; `<> 'N'` is not.** With no CHECK, a new value can
   appear tomorrow. Under `= 'S'` it would be CUT — a value nobody knows how to read would come
   to mean "outside e-commerce", which is precisely the inference honest-unknown forbids. Only an
   explicit `'N'` cuts. And the `NVL(..., 'X')` is not decoration: bare `AD_ECOMMERCE <> 'N'`
   evaluates to UNKNOWN for NULL and would silently drop the 887 undecided rows — SQL
   three-valued logic reinstating the exact bug being avoided.

Today the two forms differ by ZERO rows. The choice is made on the tail, not on the number, and
the difference shows up on the day nobody is watching.

**The mirror stores all THREE states as text — never collapse to a boolean.** Published /
refused / undecided are three facts and a boolean loses the third irrecoverably. Confirmed
against S1: `0084_products_mirror_sellable_fields.sql:9` declares `ad_ecommerce TEXT`, nullable,
no default.

**The toggle is renamed `only_ecommerce` → `only_ecommerce_eligible`** (hub). The old name
asserts "only the published ones", which is exactly what the clause deliberately does NOT do — an
option name is a contract for the next developer as much as the copy is for the operator. The
rename lands INSIDE `0083`, which has not reached main, rather than spending `0084` on an
`ALTER`; the migration-count fixtures are unaffected (the file count does not change) but every
in-flight reference to the old name does.

**Screen copy (pt-BR, operator's words):** label **"Somente elegíveis ao e-commerce"**, support
text **"Remove os produtos que o ERP marcou como fora do e-commerce. Produtos ainda sem definição
continuam no sortimento."** Expected result with the clause ON: 442 + 887 = **~1.329** — recorded
next to the 2.923 as the clause-on reference number.

**The fixture carries all three cases and asserts by VALUE:** `'S'` passes · NULL passes ·
`'N'` cuts.

Drift observed by the specialist: 2.924 vs 2.923 — one product moving between two queries against
the live database. That is what the VC-2 drift tolerance is for; it is not generosity.

### Original reasoning (the PRINCIPLE is unchanged; the FORMULA below is SUPERSEDED)

~~`ad_ecommerce IS NULL OR ad_ecommerce = 'S'` — on the live Oracle path AND on the mirror.~~
REVOKED by the DR-3 revision: the binding forms are `NVL(AD_ECOMMERCE,'X') <> 'N'` (live) and
`IS NULL OR <> 'N'` (mirror), per the A-14 table. The reasoning below is why blank must PASS,
and that reasoning is what survived; the formula that expressed it did not.

The deciding fact is not in the database; it is what the operator said when ratifying the rule,
verbatim: *"Por agora não é muito confiável mas vai ser."* A field the client ADMITS not
maintaining has an unambiguous NULL: **nobody filled it in** — not "they decided no".

**This is exactly where the (b) argument does NOT transfer, and the distinction is the whole
point.** (b) rests on "the query read ALL of TGFEST, so a missing row IS the answer: zero
available where we sell". For `AD_ECOMMERCE` the query reads the COLUMN and finds it empty —
reading a blank field is not receiving an answer, it is observing that no statement was made.
Three states, not two: `'S'` = asserted yes · `'N'` = asserted no (cut) · NULL/blank = nobody
said (passes). The strict form would empty the client's catalog the instant they enabled the
option, and the symptom would read as "the feature is broken".

**So the predicate is deliberately mixed, and the mix is not a bug:**

| clause | live (S6) | mirror (S7) | why |
|---|---|---|---|
| stock | `NVL(sellable_qty,0) > 0` | `IS NULL OR > 0` | live read all of TGFEST (absence = known zero); the xlsx tenant may have no column at all |
| `usoprod` | `= 'R'` | `IS NULL OR = 'R'` | Q1 reads the ERP column; the mirror may not carry it |
| `ad_ecommerce` | `NVL(AD_ECOMMERCE,'X') <> 'N'` | `IS NULL OR <> 'N'` | only an explicit `'N'` is an assertion; blank and any future unknown value are not, and no CHECK constrains the column |

The asymmetry is a property of **the stock and `usoprod` facts**, not a property of "the live
side". A worker who notices that the three clauses of one predicate disagree and "fixes" the
inconsistency breaks the rule — partial guard under a total sentence has already cost a whole
chip here.

The hub relayed the value distribution (`S` / `N` / NULL, overall and inside the 2.923) to
db-consult, to confirm whether an explicit `'N'` exists in volume. If it does, the cut gets
sharper and **no code changes** — it only gains a test case for `'N'`. Do not wait for it.

### The original REQUEST, kept for the trail

`catalog_page.go` `LEFT JOIN stock`, so a product with no row in the CTE gets
`stock.sellable_qty = NULL` → `QualityMissingStock` (`catalog_page.go:233-235`); `oracle/reader.go:157`
does the same for `SellableStock`. Once A-2 pins the CTE to `CODLOCAL IN (10101,10102)`, a product
stocked only in show room loses its row, goes NULL, and `(stock.sellable_qty IS NULL OR > 0)`
**readmits the 585 on the live path** — the Oracle counter would read 3.508 instead of 2.923.

The predicate half is settled: on Oracle the stock is KNOWN (the query read TGFEST), so the cut
uses `NVL(stock.sellable_qty, 0) > 0`, not honest-unknown.

The open half, routed to the hub: does the EMITTED FACT follow?
- **(a) predicate/counter only** — the fact stays NULL + `missing_stock`. Zero blast radius, but
  the same ERP renders `—` live and `0` from the mirror for the same product.
- **(b) the fact becomes `0`** with no flag, as already ratified for the mirror. Coherent across
  readers and truer. Measured blast radius: `profitability/service.go:456` starts treating those
  products as `InputQualityComplete` (today blocked as incomplete, then priced with stock 0), and
  `EstoqueTab` swaps the "importar planilha" CTA for `0`.

Chip recommends (b) — coherent with A7 and the true statement — while naming that it changes
pricing behaviour for products with no stock in the selling locations, which is a product
decision, not an implementation one. If the hub rules (a), the screen divergence is recorded as
named debt in the CLOSED event. **S6 does not ship until this is answered.**


## A-15 — lane rules for every remaining brief (hub RATIFIED, main `@ca72344`)

Ratified after S5B shipped a fully-skipped integration lane that looked green. Binding on every
brief I write from here (S7 onward). Filed by the hub as `FINDING-slice-db-lane-rules.md`
`@ca72344`, then **ratified into `docs/HARNESS-PROFILE.md` by the operator `@f1cba2a`** on main
(amendment log §11 + §3), the `%v`/`*string` corollary on the same line. Law for every future
chip, not only this one. Nothing changes in this chip's flow — S7's brief already carries it.

1. **A brief whose lane touches a database carries the env dot-source line.** Mine omitted it and
   the worker's 26 DB tests reported `ok` by skipping. The line is not boilerplate — it is the
   difference between a lane and a lane-shaped silence.
2. **The brief requires RUN / PASS / SKIP / FAIL counted per line**, not read off the tail. This
   half is what saved S5B: the worker reported SKIP 26 and declined to claim integration green.
3. **A slice that adds a migration → the orchestrator migrates the slice database BEFORE
   reviewing**, and checks `applied N` against the number of new migrations. `applied 1` for one
   new migration both applied it and proved the base carried no other drift.
4. **Corollary, same family:** a failure message must print the VALUE. `%v` on a `*string` prints
   an address, so the test discriminates correctly and then refuses to say what it saw — the next
   reader cannot tell a wrong value from a nil. Fix the message, then RE-INJECT to prove the new
   message names it.

Slices still to dispatch whose lane touches a DB: **S7** (mirror SQL + golden fixture) and any
later slice that reaches Postgres. S6 is unit-only by design — SQL is asserted as generated text —
so its brief correctly carries NO env, and instead carries the counting instrument plus an
instruction to claim nothing about the `//go:build integration` files its commands never compile.

## A-16 — `query_repository.go` keeps its 12 columns, with the condition written (hub ruling)

The S5B review found a THIRD site dropping `usoprod`/`ad_ecommerce`:
`query_repository.go:196`, backing `LatestCompletedSnapshot`. No production caller was found —
only the interface and tests.

Ruling: **comment with condition, no extension.** Widening a SELECT for a reader nobody reads is
the imagined-problem class. Applied with the F-3 technique — the comment says the columns are
deliberately absent because no consumer asks, and states that whoever wires a real consumer
inherits BOTH halves: the two columns AND a test asserting they ARRIVE at that consumer, since
the same fact was already lost twice upstream by testing departure instead of arrival.

If S7's P4 finds a real caller the hub did not see, it becomes S7's work with this finding as the
authority. Archive note carried from the hub: there is a D-113 memory of a blind
`LatestCompletedSnapshot` displacing Sankhya — if S7 touches that path, measure before writing.
