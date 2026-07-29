# CHIP-VENDAVEL — dispatch ledger

Chip: CHIP-VENDAVEL (MIS-006, interim post-mission chip).
Worktree: `.claude/worktrees/chip-vendavel` · branch `worktree-chip-vendavel`.
BASE-SHA: `554788d576d04b21719f4a17e4702dd0f0aff4e1`.
Hub session: `local_99feb041-a5b3-4161-b6dc-bd38e65b6156`.
Contract: `.mnfs/MIS-006-integracao-fundacao/_pos-chip-vendavel/validation-contract.md` (VC-1..VC-7).

Ownership declared (write-set, single track — no concurrent chip in flight):
`apps/server_core/migrations/` block 0083–0084 only · `internal_read/adapters/oracle/{sync.go,catalog_page.go}` ·
`products_mirror` writers/readers (additive columns `usoprod`, `ad_ecommerce`) · xlsx parser optional columns ·
module `tenant_config` (3 boolean columns + handler) · `MirrorMatcher` path ·
FE `apps/web/src/pages/integracoes/IntegracoesPage.tsx`, `packages/feature-products` catalog page,
`packages/web-query` + `packages/sdk-runtime` domain blocks.
Pre-authorised additive grants: composition-root wiring region (`internal/composition/root.go`),
SDK domain block, migration-count fixtures.

## Operator ratifications received in this chip session

`RATIFIED-BY-OPERATOR:` 2026-07-29 — stock LOCATION rule, amending the context pack's
`only_em_estoque` (which named only `CODEMP IN (1,2)`). Operator, verbatim:

> "Uma coisa CODLOCAL é legal 10101 é estoque revenda, 10108 é show room não é pra contar, e
> 10102 é outlet esses que vendem"

Asked whether it binds the filter only or also the stock NUMBER already rendered on the catalog
screen (which reads only `10101` today), the operator chose "filtro e número, os dois". So
sellable stock = `CODEMP IN (1, 2) AND CODLOCAL IN (10101, 10102)` — show room `10108` excluded —
and the same definition serves the catalog's displayed stock and the mirror's `estoque_total`.

Open against this ratification (routed as `REQUEST db-consult` to the hub, 2026-07-29): whether
available = `ESTOQUE - RESERVADO` (what both existing queries already do) or physical `ESTOQUE`
(operator wrote "e temos os reserváveis também"); whether `CODPARC = 0` remains the own-stock
predicate; the live count under the amended rule (the pack's 3.822 was measured with no CODLOCAL
filter, so it is an upper bound); and whether any other CODLOCAL sells. Not blocking — the chip
implements `CODEMP IN (1,2) AND CODLOCAL IN (10101,10102)` with available = `ESTOQUE - RESERVADO`
and leaves the VC-2 number to the hub's live measurement post-merge.

Rows are written AT DISPATCH TIME (core §1 ledger rule) and completed with the result.

| # | phase | role | model / effort | path | prompt | log / output | result | SHA |
|---|---|---|---|---|---|---|---|---|
| 1 | pre-P2 | Investigator / repo map | Claude subagent (Explore, sonnet) | Agent tool, sync | inline brief: map tenant_config / products_mirror writers / sync Q1+Q4 / catalog_page / xlsx parser / MirrorMatcher / OpenAPI+SDK / FE cards + hooks + tests | returned to session context; facts quoted verbatim into the P2 planner prompt (`scratchpad/prompt-p2-planner.md`) | COMPLETED — 13-section `file:line` map | base 554788d5 |
| 2 | P2 | Feature planner (batch, whole chip) | `gpt-5.6-sol` / `medium` | OS-process codex, `--sandbox read-only`, stdin closed, cwd = chip worktree | `scratchpad/prompt-p2-planner.md` | `scratchpad/agent__p2-planner.log` · `scratchpad/agent__p2-planner.last.md` | COMPLETED — 14 slice cards (S0..S13), write-set DAG, contract-satisfiability table, verification map, seam-closure checklist, must-fail design. Copied verbatim to `BATCH-PLAN.md`; chip amendments appended there (A-1..A-6) because the plan predates the closed A7 ruling | base 554788d5 |
| 3 | P2 | `REQUEST db-consult` → hub → MNOS Oracle specialist | — (cross-session event) | hub `local_99feb041` | 4 questions: RESERVADO semantics, `CODPARC = 0`, live count under the amended rule, other selling CODLOCAL | hub reply quoted verbatim into `ADDENDUM-01-codlocal.md` | ANSWERED — A7 closed: subtract RESERVADO (354 products zero out by reservation alone), keep `CODPARC = 0` (defensive), **VC-2 = 2.923 distinct CODPROD**, named whitelist of 8 locations (only 10101 + 10102 sell). Provenance: db-consult MNOS 2026-07-29 via hub, read-only aggregation | main @ aee1a222 |
| 4 | P2 | `REQUEST DR-1` → hub | — (cross-session event) | hub `local_99feb041` | mirror-catalog parity for `xlsx`/`catalogo_cliente` tenants: apply the predicate + live counts on the upload path (~8 files / ~215 lines), or issue a named+dated deferral with a non-numeric UI state | hub ruling quoted into `BATCH-PLAN.md` §A-7 | **GRANTED** — parity approved, no deferral; the pack is corrected to THREE binding read paths (Oracle · mirror-serving-catalog · MirrorMatcher). Conditions: SQL filter before `LIMIT+1` with a >1-page fixture · measure absent-vs-zero on the xlsx writer · must-fail at the CALLER (`routing.Reader`, source `xlsx`). Mirror counter counts against the mirror. VC-2's 2.923 is the LIVE ORACLE number only; the mirror path discharges by agreement within its own population | main @ 66f32125 |
| 5 | P3 | Implement worker — slice S1-SCHEMA | `gpt-5.6-luna` / `high` | OS-process codex, `--sandbox workspace-write`, stdin closed, cwd = chip worktree | `scratchpad/prompt-s1-schema.md` | `scratchpad/agent__s1-schema.log` · `scratchpad/agent__s1-schema.last.md` | **ACCEPTED** at P4 — see the review below | `23b9d428` |
| 6 | P3 | Implement worker — slice S2-CONFIG-DB | `gpt-5.6-luna` / `high` | OS-process codex, `--sandbox workspace-write`, stdin closed, cwd = chip worktree | `scratchpad/prompt-s2-config-db.md` | `scratchpad/agent__s2-config-db.log` · `scratchpad/agent__s2-config-db.last.md` | **ACCEPTED** at P4, committed by the orchestrator (the worker could not) | `29e7cb69` |

## P4 adversarial review — S2-CONFIG-DB (`29e7cb69`)

**ACCEPTED**, with the same treatment S1 got: the worker's report is a claim, not evidence.

1. **`go build ./...` refuted a SECOND time.** The worker again reported the build "blocked by
   VCS stamping" with `-buildvcs=false` passing. Measured from `apps/server_core` with absolute
   caches: **EXIT=0**, unqualified. Two workers, two identical false alarms, both from the codex
   seat and neither reproducible at the orchestrator's. Recorded as an environment signature of
   that seat, NOT as an L0 result — and no `-buildvcs=false` is buried in any lane to silence it,
   because that trades an alarm for blindness.
2. **The commit failure was real, and it is a finding.** `git commit` died on
   `.git/worktrees/chip-vendavel/index.lock: permission denied`. No lock file exists now. In a
   git worktree the real git directory lives under the MAIN checkout
   (`.git/worktrees/<name>/`), outside the sandbox's writable root — so a worker under
   `--sandbox workspace-write` can be denied the index while every file write succeeds. S1's
   worker did commit, so it is not deterministic; the orchestrator commits when it happens.
   Worth reporting to the hub for the profile.
3. **The RED is weaker than S1's and needed the same correction.**
   `repository_test.go:79:9: got.SellableAssortment undefined` is a COMPILE error — it proves the
   field did not exist, not that any assertion discriminates a value. Two mutations, both killed:
   - `Set`'s `ON CONFLICT DO UPDATE` given `only_ecommerce = false` (a source switch quietly
     resetting the operator's rule) → `repository_test.go:186: after source switch OnlyEcommerce
     = false, want true`.
   - `SetSellableAssortment`'s `WHERE tenant_id = $1` replaced by `WHERE $1::text IS NOT NULL`
     (tenant scoping dropped, the classic leak) → killed in THREE independent places: the first
     tenant re-read after the second write (`:160`, `:163`), the source-switch preservation
     (`:180`, `:183`), and the missing-tenant sentinel returning `<nil>` instead of
     `ErrUnknownActiveSource` (`:190`).
   A first attempt at that second mutation (`WHERE $1 IS NOT NULL`, no cast) is NOT counted: it
   died on `SQLSTATE 42P18 could not determine data type of parameter $1`, which is the mutation
   failing to compile rather than the test catching a leak. A mutation that kills itself proves
   nothing.
4. **Counts measured, not read off a tail:** 13 PASS, **0 SKIP**, 0 FAIL across two packages.
   The skip count matters here — these tests call `SkipWithoutTarget` and a fully skipped run
   prints `ok` too.
5. **Diff reviewed for scope:** `Get` reads the three columns in the SAME query, still
   `WHERE tenant_id = $1`, no second round trip and no default-if-missing branch; `Set` is
   untouched; `SellableAssortment` holds plain booleans, which is correct because the columns are
   `NOT NULL DEFAULT` — the value is always known, so a pointer would invent an unknown that
   cannot occur.

## P4 adversarial review — S1-SCHEMA (`23b9d428`), by the orchestrator

**ACCEPTED.** Four checks, three of them against the worker's own report:

1. **The worker's build failure is REFUTED by measurement.** It reported `go build ./...`
   failing with `error obtaining VCS status: exit status 128`. Re-run from `apps/server_core`
   with absolute `GOCACHE`/`GOMODCACHE`: `go build ./...` → **EXIT=0**, and
   `go build -buildvcs=false ./...` → EXIT=0. Environment false alarm at the worker's seat, not
   an L0 failure. A worker's report is never the evidence.
2. **The RED is real but weak on its own** — `sellable_assortment_test.go:9: open
   0083_sellable_assortment_config.sql: file does not exist` names the test and fails, but a
   file-absent RED proves only that the file was missing, NOT that the assertions discriminate
   the ratified VALUES. So the orchestrator ran a mutation: flipped `only_ecommerce` to
   `DEFAULT true` and made `usoprod` `NOT NULL` in the two migrations, re-ran, and both mutants
   were killed by name:
   `missing exact declaration "…only_ecommerce boolean not null default false"` and
   `products_mirror.usoprod must have no DEFAULT and no NOT NULL`. Tree restored, re-run green,
   `git status --porcelain migrations/` empty. The defaults and the nullability are now pinned
   by an instrument that has been shown to go red.
3. **Independent confirmation of the migration count:** the orchestrator's own
   `go run ./cmd/testdb migrate` against a fresh database printed `applied 71 migration(s)` —
   the same 71 both `runner_test.go` assertions were bumped to, measured from the DB side
   instead of the fixture side. Both assertion sites were bumped (`:25` and `:64`).
4. **Scope and hygiene:** no `.gomodcache` at the repo root (only `apps/server_core/.gomodcache`);
   `writer.go` gained the two struct fields and NOT the `upsertSQL` columns, which is correct —
   S4 owns that seam and has its own test.

Reservation recorded, not blocking: `mirror_test.go` asserts through `reflect.FieldByName`. That
is only justified because the test had to compile before the fields existed (a direct reference
would have made the RED a compile error rather than a named failing test). Later slices assert
these fields directly; if any of them keeps using reflection instead, that is test theater and
gets rejected there.

## Test-database provisioning (orchestrator glue, not a slice)

The hermetic lane is the chip's own instrument (profile §4) — distinct from the hub-owned dev
stack (`:8080`/`:5174`), which the chip never touches. Provisioned once, so no worker starts a
container or invents a DSN:

- `npm run harness:pg:up` → session container `mpc-pg-session-3eee515d`, port `55864`,
  `status=ready` (one per checkout, hashed from the path — hub and chip never collide).
- Slice database `mpc_test_deadbeefdeadbeefdeadbeefdeadbeef`. The name is not decorative:
  `testsupport/postgres/target.go:19` rejects anything that does not match
  `^mpc_test_[0-9a-f]{32}$` with `HPG_TARGET_INVALID` — a friendly name like `mpc_chip_vendavel`
  is refused, which is how the first attempt failed.
- `go run ./cmd/testdb migrate` → `applied 71 migration(s)`.
- Workers receive `scratchpad/testdb-env.ps1` (dot-source only). It derives the URL from
  `scripts/.runs/pg-session.json` at run time, so the container password is never written into
  a prompt, a log, or the evidence pack.

`SkipWithoutTarget` makes a skipped run and a green run byte-identical at the tail (profile §11,
vacuous green). Every DB slice brief therefore demands RUN/PASS/**SKIP** counts as a result.

---

## S2B-RENAME — `only_ecommerce` → `only_ecommerce_eligible`

| field | value |
|---|---|
| dispatched | worker `gpt-5.6-luna` / `high`, OS-process, `agent__s2b-rename.log` + `.last.md` |
| base | `2e943795` |
| commit | `f76c1325` (orchestrator committed; see the build claim below) |
| evidence | `evidence/S2B-red.txt`, `evidence/S2B-green.txt` (worker), `evidence/S2B-orchestrator.txt` |

Hub ruling. The name asserted the opposite of the clause it controls: `only_ecommerce` reads as
"only the published ones", while the ratified clause is `NVL(AD_ECOMMERCE,'X') <> 'N'`, which
removes only what the ERP explicitly marked as outside e-commerce and keeps everything undecided.
Measured live, the strict reading would cut 2.923 → 442. The rename lands INSIDE `0083`, which has
not reached main — no `ALTER RENAME`, no new file.

### P4 — what the orchestrator measured, not what the worker reported

The RED discriminates: the migration test failed naming the OLD declaration string verbatim, so
the assertion is bound to the exact declaration rather than to the column's existence.

The slice database was dropped, recreated and re-migrated after the rename — mandatory, because a
stale schema would make the worker's own package read as a code defect. `applied 71 migration(s)`,
the SAME count as before the rename: independent confirmation, by measurement rather than by
assertion, that the count fixture is untouched.

`go test ./internal/modules/tenant_config/... ./migrations -count=1 -v` → **33 RUN / 33 PASS /
0 SKIP / 0 FAIL**. `SKIP=0` carries the weight: without `MPC_TEST_DATABASE_URL` these tests skip
silently, and a fully skipped run is byte-identical at the tail to a fully green one (profile §11).

Sweep at the CALLER, anchored, asserting ABSENCE — `git grep -nE 'only_ecommerce($|[^_])|OnlyEcommerce($|[^E])'`:
zero in `apps/server_core`, `apps/web`, `packages/`, `contracts/`; zero `_eligible_eligible` (the
prefix trap did not fire); new-name count **19**, identical to the pre-image count of the old name
taken from `git HEAD` BEFORE the worker ran — one-to-one, nothing lost or doubled. Five survivors
are deliberate: `BATCH-PLAN:1066` names the rename and must quote the dead name, and
`DISPATCH-LEDGER:70,71,103,106` are verbatim history of the S1/S2 runs. Rewriting those would
falsify the trail.

### The rename destapou prosa falsa, corrigida in the same commit

`BATCH-PLAN:60` still carried `m.ad_ecommerce = 'S'` — the formula the DR-3 revision revoked — and
`CONTEXT-PACK:23` carried both the dead name and the revoked semantics. A worker reads the plan,
not the migration: that clause would have entered S7 as the strict reading and cut the assortment
to 442. The S3 brief carried the same falsehood and was corrected BEFORE dispatch.

### Finding — third false build claim from the same seat

The worker reported `go build ./...` failing with `error obtaining VCS status: exit status 128`
and **withheld its commit on that basis**. Measured in the same worktree, `cd apps/server_core`,
absolute caches: `EXIT=0`. Third identical occurrence (S1, S2, S2B) — a signature of the codex
seat, not a repo fact. No `-buildvcs=false` was buried in the lane: that trades an alarm for a
blindness (hub, binding). The real cost this time was a green slice left uncommitted, so the S3
brief now states that this observation must not withhold a commit.

### Finding — the sweep instrument is partial, cause NOT established

The Grep tool rooted at the worktree's ABSOLUTE path returned 5 matches and missed all 19 code
occurrences; `git grep` from inside the worktree returns all 19. The `.gitignore '.claude/*'`
hypothesis was tested and REFUTED (`rg` from inside the worktree finds them). Recorded as an
instrument constraint, not a diagnosis. Operationally: rename sweeps in this chip use `git grep`
from inside the worktree. A partial sweep is worse than none — it has the shape of a complete one,
and a gate seat using the wrong tool would report "zero occurrences" with full confidence.

---

## S3-XLSX — optional `USOPROD` / `AD_ECOMMERCE` on the upload path

| field | value |
|---|---|
| dispatched | worker `gpt-5.6-luna` / `high`, OS-process, `agent__s3-xlsx.log` + `.last.md` |
| base | `a94694c2` |
| commit | `8f825a7c` |
| evidence | `evidence/S3-red.txt`, `evidence/S3-green.txt` (worker), `evidence/S3-orchestrator.txt` |

Production diff is correct: both fields flow through `optionalCell`, the two required lists
(`Parse`, `ParseLenient`) are untouched, and `usoprod`/`ad_ecommerce` reach the stage table, the
`CopyFrom` column list, the insert and the `ON CONFLICT DO UPDATE SET`. `nullableStringValue` is a
justified thin wrapper — the existing `nullableString` takes a value, not a pointer.

### P4 — three mutations, three kills

The slice is entirely honest-unknown, so presence assertions would be worthless. Injected and
killed: (M1) parser emitting `""` for an absent column → caught at `parser_test.go:81`;
(M2) `nullableStringValue` persisting `""` instead of SQL NULL → caught at
`mirror_repository_integration_test.go:37`, which matters because the round-trip is what proves
SQL NULL rather than a Go nil that a nil-preserving scan would yield anyway; (M3) adding
`USOPROD` to the required list → four PRE-EXISTING tests fail with `MISSING_REQUIRED_COLUMN`, so
"legacy workbooks keep importing" is guarded rather than merely true by inspection. All reverted,
tree clean against `8f825a7c`, `erp_import/...` green.

### Finding — a pre-existing clock-skew flake, NOT caused by S3

The worker reported `TestGetImportChainCountsCurrentQueueAcrossInstallations` failing as
"unrelated timing ... reproduced in isolation". The second half is FALSE and was refuted: it
PASSES in isolation and PASSES the full 28-test package; it fails only in the `erp_import/...`
tree run, and non-deterministically (run 1 FAIL, run 2 PASS at the same tip).

Root cause: `query_repository.go` selects `statement_timestamp() AS queue_read_at` — the
POSTGRES clock, inside the container — while the test brackets it with host-side `time.Now()`
(`chain_query_repository_integration_test.go:74-85`). It asserts a container clock reading falls
inside a host clock window of ~10ms. The observed miss was ~77µs, and the precision tell is in
the output: 6 decimals (PG microsecond) against 7 (Windows 100ns).

The test is NOT in `8f825a7c` — S3 neither wrote nor touched it. Latent defect, surfaced by added
package load. Out of the S3 write-set, so reported to the hub and not fixed here. It will bite
the P5 ladder and any gate seat that reads a single red tree run as attributable.

### Plan rot corrected in the same commit (hub standing order)

Formula revocation is a write-set change. Sweeping the pack for every predicate DR-2/DR-3/A-14
touched found the S6 slice card still carrying ALL THREE dead clauses in its `done_criteria` —
`(p.USOPROD IS NULL OR = 'R')`, `(stock.sellable_qty IS NULL OR > 0)` and
`(p.AD_ECOMMERCE ... = 'S')` — superseded by amendments ~750 lines further down the same file.
S6's worker reads the card, not the amendment. Corrected to the binding live forms with the
supersession recorded inline and a pointer to the A-14 asymmetry table. `BATCH-PLAN:1086` also
printed the revoked formula under the heading "Original reasoning (the principle, unchanged)" —
the principle is unchanged, the formula is not; struck through and marked.

## S4-SANKHYA-SYNC — Q1 reads the two fields, Q4 is pinned to the sellable cut

| field | value |
|---|---|
| dispatched | worker `gpt-5.6-luna` / `high`, OS-process, `agent__s4-sankhya-sync.log` + `.last.md` |
| base | `4bd84e62` |
| commit | `84e87e63` — **created by the orchestrator on the worker's behalf** (see below) |
| evidence | `evidence/S4-red.txt`, `evidence/S4-green.txt` (worker), `evidence/S4-orchestrator.txt` |

Production diff is correct. Q1 selects `p.USOPROD, p.AD_ECOMMERCE`, scans them as
`sql.NullString` and maps them through the existing `nullStr`, so Oracle CHAR-padding and blank
cells arrive as nil rather than `""`. Q4 carries `CODEMP IN (1, 2)` and
`CODLOCAL IN (10101, 10102)` from one named constant that also carries both location names, and
sums `NVL(ESTOQUE,0) - NVL(RESERVADO,0)` — the NVL pair is load-bearing, a NULL `RESERVADO`
would otherwise null the row out of the SUM entirely. Whitelist, never blacklist: a location
created tomorrow must not start selling by itself.

The pin creates the defect it must not ship, and `applyStock` changes with it. It now iterates
the **Q1** population and writes `totals[codprod]`, which is `0` for a product the pinned Q4 did
not return. On this path stock is KNOWN — the sync read `TGFEST` — so absence inside the cut is
zero available where we sell, not unknown. It stops exactly there: a product absent from Q1 is
never in that map, so `absent_in_last_snapshot` keeps its last-known values (ADR-04). Both sides
of ADR-17, each on its own side.

### P4 — two mutations, two kills, plus one test-quality defect fixed

| mutation | result |
|---|---|
| `applyStock` reverted to the old Q4-only loop (A-9) | KILLED — `sync_test.go:88: 100.EstoqueTotal = <nil>, want known zero 0` |
| `AND CODEMP IN (1, 2)` deleted from `sankhyaStockSQL` | KILLED — `sync_test.go:25: sankhyaStockSQL missing sellable company predicate CODEMP IN (1, 2)` |

Both reverted; `git diff --stat` on `sync.go` back to `20 insertions(+), 11 deletions(-)`.

Lanes counted per-test, never tail-read (profile §11): untagged **RUN=186 PASS=185 SKIP=1
FAIL=0 EXIT=0**; `-tags=integration` **RUN=188 PASS=187 SKIP=1 FAIL=0 EXIT=0**. The one SKIP is
named — `TestOracleLiveBaseline`, which needs `MPC_ORACLE_CONNECT_STRING` this seat must not
have — and is explicitly NOT evidence: the live half of VC-5 is the hub's post-merge step. All
four S4 tests confirmed PASS **by name**, including `TestPgWriterPersistsSellableAssortmentFields`
against the real database.

**Defect found in the slice's own test, fixed by the orchestrator:** M1's first injection
produced a nil-deref PANIC at `sync_test.go:88`, not a diagnostic. The two known-zero assertions
used `t.Errorf` and then dereferenced `EstoqueTotal` unguarded, so the A-9 regression they exist
to catch arrived as a stack trace. Made `t.Fatalf`, reason recorded next to them, M1 re-injected
to get the clean line above. The teeth were always there; the failure MESSAGE was what got lost,
exactly when it was needed. A mutation that "kills" by panicking has not shown that the test
discriminates the VALUE.

A-11 discharged by value across three populations: `100` stocked only outside the cut → known
`0`; `200` sellable in 10101 → `7`; `300` absent from Q1 → `mw.absent["300"]` true AND the prior
`*EstoqueTotal == 9` preserved. A-12 discharged: the F-01 comment recorded the opposite rule and
was true only while Q4 read all of `TGFEST` — rewritten with the trail, not softened.

### Finding — fourth identical false build claim, and a second index.lock denial

`go build ./...` reported EXIT=1 with `error obtaining VCS status: exit status 128`, for the
fourth time from the codex seat (S1, S2, S2B, S4). Measured here, same command, same worktree:
EXIT=0. Under the ratified rule the worker filed it as an observation and did not withhold work;
`-buildvcs=false` was not added anywhere.

The worker then could not commit: `fatal: Unable to create
'.git/worktrees/chip-vendavel/index.lock': Permission denied` — the same denial S2B hit, while
S1 and S3 committed fine from the same seat, so it is non-deterministic and still unattributed.
`84e87e63` was created by the orchestrator after P4 completed. Recorded rather than smoothed
over: the authorship of that commit is not the authorship of the code in it.

---

## S5-MATCHER — round 2 (round 1 BLOCKED at the scope wall, correctly)

| field | value |
|---|---|
| slice | `S5-MATCHER` |
| worker | codex `gpt-5.6-luna` / `high` (implement, standard) |
| dispatch | stdin (`codex exec … -`), UTF-8 pinned both sides — profile §12 |
| brief | `scratchpad/prompt-s5-matcher.md` (round-2 RE-BRIEF appended) |
| log | `scratchpad/agent__s5-matcher.log` |
| final message | `scratchpad/agent__s5-matcher.last.md` |
| commit | `a62c56d68362d92e32f6abc17d3ceee18ef2939d` (worker committed; no index.lock denial this time) |
| evidence | `evidence/S5-red.txt`, `evidence/S5-green.txt`, `evidence/S5-orchestrator.txt` |
| lane (orchestrator re-run) | RUN 58 · PASS 58 · SKIP 0 · FAIL 0, counted per line |
| P4 | ACCEPTED as a slice; see §3 of the review — the chip's promise is NOT discharged for xlsx |

Round 1 stopped at the scope wall instead of shipping a rule that cuts nothing, and the hub
granted the two bounded edits in `mirror_query_repository.go`. The RED it left behind is what
makes this round's green non-vacuous, so it was preserved rather than regenerated.

### Grant honoured exactly

Two edits, no filtering: the column constant (`:13`) and the scanner (`:22`, `:39-40`). All four
mirror read sites share that constant and that scanner, so the SELECT list and the `Scan` order
move together by construction — checked, not assumed.

### The abstraction was measured before it was accepted

The worker added a context-carried `SellableAssortmentPolicy` and justified it with an import
cycle. Anti-slop §4 forbids speculative abstraction, so the cycle was measured:
`tenant_config/active_source.go:11` imports the reader; the reader imports `tenant_config` zero
times. Cycle real, two named consumers, carrier accepted.

### REQUEST open — the rule is inert for one of the two sources

The worker reported in two lines that the upstream selectors omit the columns. Measured here it
is one hop earlier and it disables the deliverable for xlsx: `import_repository.go:67` copies 14
columns into `erp_import_products` and neither of the two is among them, because migration 0084
added them to `products_mirror` ONLY. So `xlsx/parser.go:116` reads `USOPROD` off the sheet and
the value dies at the first write. Mirror `usoprod` stays NULL for xlsx, nil PASSES the ratified
`IS NULL OR = 'R'`, and the rule cuts nothing — the round-1 blocker's own sentence, one hop up.
Sankhya is unaffected (`mirror/writer.go:150` writes straight to `products_mirror`).

The S5 test cannot catch it: it seeds the mirror with a raw `UPDATE`, which was the only way to
reach the reader at all. That green proves the reader and the matcher, not the pipeline.

Not fixed and not decided here — the repair needs migration 0085 (grant was 0083–0084) plus two
files in no slice's write-set. REQUEST filed with three options; recommended a bounded S5B.

Second REQUEST in the same message: the comparisons are case-sensitive (`:231`, `:247`) and
nothing normalizes case on this path. Sankhya returns `R`; the spreadsheet is operator-authored.
A lowercase `n` passes the e-commerce rule and a lowercase `r` is cut from resale — wrong in both
directions, worse in the permissive one. Real ambiguity, so it went to the hub instead of being
settled locally.

## S5B-WRITE-CHAIN — gpt-5.6-luna / high — OS-process

| field | value |
|---|---|
| dispatched | D-122, base `5bcf0bef` |
| prompt | `scratchpad/prompt-s5b-write-chain.md` |
| log | `scratchpad/agent__s5b-write-chain.log` |
| last message | `scratchpad/agent__s5b-write-chain.last.md` |
| tail marker | returned verbatim — `S5B-BRIEF-RECEIVED … TAIL-MARKER-4d90c2` |
| commit | `fe0f14a1` |
| review | `evidence/S5B-orchestrator.txt` |
| lanes | `evidence/S5B-green.txt` (worker) + orchestrator correction (same file) |
| verdict | ACCEPTED |

### The lane the worker ran proved nothing, and the brief is why

The commands block I wrote for this slice omitted the `. testdb-env.ps1` dot-source line that
S1 through S5 all carried. Without it `MPC_TEST_DATABASE_URL` is unset, every DB test takes
`SkipWithoutTarget`, and the package still prints `ok` with exit 0 — profile §11 signature 3,
by the book. The worker's integration run was RUN 27 / PASS 1 / **SKIP 26**, and the whole-chain
round-trip test, the entire reason the slice exists, was one of the skips.

The worker reported the skip counts explicitly and wrote *"this is not a claim of live
integration GREEN"* in its own evidence file. That sentence is the only reason the vacuum was
visible instead of shipped. The instrument the profile mandates — count SKIP per line, never read
the tail — worked, and it worked in the worker's hands, against the worker's own result.

Re-measured at the orchestrator seat with the env sourced: **RUN 29 / PASS 29 / SKIP 0 / FAIL 0**
(lane A) and **20 / 20 / 0 / 0** (lane B), after applying the granted migration 0085 to the slice
database — the first run failed 13 tests with `column "usoprod" does not exist (SQLSTATE 42703)`,
because a migration in the tree is not a migration in the database. `applied 1 migration(s)`,
exactly one, which also proves the slice DB carried no other drift.

### Must-fail, four times, because a skipped test has never failed either

Each protection severed separately and reverted: hop 1 (`CopyFrom` value), hop 2 (`Scan`
assignment), the xlsx fold, the Sankhya fold. Four value-naming REDs, quoted in
`evidence/S5B-green.txt`. The two hops fail identically — `want "R", got <nil>` — which is the
shape of a test that asserts ARRIVAL rather than departure, exactly as the D-122 ruling demanded.

### A pointer where a value should have been

The worker's own RED read `canonical USOPROD = 0x34f73c6925b0, want R`: correct assertion, `%v`
on a `*string`. The test discriminated the value and then declined to say what it saw. Fixed with
a `show()` helper and **re-injected** to confirm the new message names it (`= "r", want R`) —
repairing the message without re-proving the RED would have been the same class of error.

### Glue applied at this seat, all pre-approved

F-3 (the fail-open comment carrying why + validity condition), F-4 (dead conjunct), and the
message helper. Declared in §4 of the review rather than folded silently into the worker's diff.

### Findings for the hub

- `query_repository.go:196` still drops both columns from its 12-column SELECT over
  `erp_import_products`. No production caller found. Reported, untouched — a third place the same
  fact can be lost.
- Slice briefs that touch a DB must carry the dot-source line, and must require SKIP to be
  reported explicitly.
- Nothing migrates the slice database automatically; `OpenPool` validates and connects, it does
  not migrate. After a slice adds a migration, the orchestrator migrates before reviewing.

## S6-ORACLE-CATALOG — the live query cuts, the counter counts, the outlet joins

Worker gpt-5.6-luna / high, OS-process. Commit `8303f7c6`; orchestrator glue on top.
Evidence: `evidence/S6-red.txt`, `evidence/S6-green.txt`, `evidence/S6-orchestrator.txt`.

Five RED assertions, each naming its own miss: the missing live predicate, page/count drift,
`missing_stock` on a fact the ERP answered, an empty `NOT IN ()`, and the location defaults.

### P4 — counted at my seat, because the worker's own numbers were wrong

internal_read RUN 191 · PASS 190 · SKIP 1 · FAIL 0. pricing RUN 158 · PASS 158 · SKIP 0 · FAIL 0.
`vet` 0, `build` 0. The worker reported PASS=130 and PASS=116 — an undercount on green lanes.
Harmless direction, still a number that did not come from counting. The lane is unit-only by
design (SQL asserted as generated text), so the absent DB env is correct here; A-15's dot-source
rule binds S7, whose lane reaches Postgres. The one SKIP is `TestOracleLiveBaseline`, pre-existing
at `f51cabca`. The worker's `build EXIT=1 / error obtaining VCS status` did not reproduce at my
seat and is recorded as a seat-local git quirk — it was right not to paper it over.

### The predicate cuts before the page, and the count reads the same rule

Appended ahead of cursor/search/`FETCH FIRST`, so the cut happens in SQL before pagination — a
predicate applied after the fetch would have paged the whole catalog and thinned each page, and
the count would have disagreed with the screen. Page and count share `catalogStockCTE()` and
`catalogAssortmentPredicate()`; drift is prevented by there being one source, not by a comment.
`CatalogProductFactsByIDs` passes `IncludeAll: true` — a caller naming ids is asking about THOSE
products. Both tests pinning the old location defaults were re-pointed BY VALUE, not loosened.

`missing_stock` is not retired: 8 non-test producers remain. S6 removed exactly one — the live
catalog page, the path DR-2/A-14 ruled on, where the ERP was asked and answered, so absence is a
known negative and `0` is the honest value. ADR-17's two sides, each on its own path.

### Finding — the new optional port is asserted at 1 of 5 seats, and the other 4 fail silently

`ports.CatalogAssortmentReader` has exactly one compile-time assert (`oracle/catalog_page.go:312`)
and exactly one implementer. Its sibling `CatalogPageReader` has five
(`oracle:311`, `routing:159`, `cache:312`, `timing:136`, `service:66`). The live chain composed at
`composition/root.go:449→453→479` is oracle → cache → timing → routing, and every seat forwards by
RUNTIME assertion: `routing/reader.go:151` fails into `ReadErrorSourceUnavailable`, which is a 503
on the screen. A decorator missing the new port therefore does not fail to build — it deletes the
capability and calls the source unavailable. That is the catalog-503 defect of CHIP-M02 rebuilt,
whose lesson was compile-time asserts for EVERY optional port.

Not S6's defect: the decorators are outside its write-set, and asserting a type that does not yet
implement the interface is a build failure — it would drag S9's work into S6. What was owed was
the condition written where the next writer stands. Applied as glue (F-3) on the interface
declaration: it names the four seats, names the 503 mechanism, demands an arrival test that reads
the count THROUGH the composed reader from `composition/root.go`, and forbids a runtime
`.(CatalogAssortmentReader)` with a fallback at the HTTP seam — the shape that turns a missing seat
into a quiet wrong answer instead of a build error.

Same family as S5B's two hops and A-16's third site: a value that leaves its source and is never
asserted to ARRIVE. Here the value is a capability, and the arrival test belongs to S9.

## S7-MIRROR-CATALOG — dispatched

gpt-5.6-sol / low (complex, per the card), OS-process, prompt over stdin.
Log `scratchpad/agent__s7-mirror-catalog.log`, last message `agent__s7-mirror-catalog.last.md`.
Row written at dispatch time; result rows follow after P4.

Brief carries, beyond the card: A-15 in full — the `. testdb-env.ps1` dot-source line, per-line
RUN/PASS/SKIP/FAIL counting with `PASS+SKIP==RUN` reconciled, the prohibition on claiming
anything about `//go:build integration` files the commands never compiled, and the `%v`-on-a-
pointer corollary with the required `%s` + nil-rendering helper. Also the A-14 asymmetry table,
so the worker does not "fix" the mirror's NULL tolerance into the live predicate; the D-113
archive note on `LatestCompletedSnapshot`; and the instruction to prove pagination with a
fixture larger than one page rather than a sampled `Limit=N`.

### A-15 ratified into the profile by the operator

`@f1cba2a` on main, amendment log §11 + §3, `%v`/`*string` corollary included. What was chip-local
law since the S5B acceptance is now law for every future chip. The S7 brief was written under it
before the ratification arrived, so nothing in this chip changes.

### P4 — the worker shipped one red, so I injected the other four

Lanes at my seat with the env sourced: postgres integration RUN 32 · PASS 32 · **SKIP 0** · FAIL 0;
`erp_import/...` unit RUN 142 · PASS 142 · SKIP 0 · FAIL 0; `vet` 0, `build` 0. SKIP 0 is the
load-bearing number — this is the lane S5B faked by omission. The worker's `build EXIT=1 /
error obtaining VCS status` did not reproduce here, second slice running; neither worker reached
for `-buildvcs=false`.

The worker delivered ONE red and said plainly the rest were written during implementation. Honest
and insufficient: criterion 7 demanded both directions. Four mutations, four kills, each naming a
value — M1 collisions before the cut (`strict survivor quality = [complete ean_collision]`), M2
policy ignored (`relaxed twin 201 quality = [complete ambiguous_product]`), M3 stock clause
neutralised (`[a c d], want [a d]`), M4 page stops cutting in SQL (`first filtered fetch =
[301 302 303]` — the excluded 302 consumed a LIMIT slot). Criterion 8 closed by a seeded dirty
`usoprod='r'` row being CUT; criterion 9 by an upload with physical stock and no reservation
collapsing to a genuine unknown. No test deleted or loosened. S13-CLOCK flake fired 3× and passed
on re-run; recorded, not dropped. Residual recorded: the count query lacks the page's
`DISTINCT ON (codigo_produto::bigint)`, so a leading-zero code would over-report — narrow, and the
worker wrote the collapse test, so it was not blind.

## A-17 — the tenant's toggles never reached the catalog (hub RULING, pack `@67e1edc9`)

Found at S7's P4, escalated rather than decided. `catalogPage` and `GetCatalogAssortmentCounts`
took the rule from a hardcoded `defaultSellableAssortment()` instead of
`SellableAssortmentFromContext(ctx)` — the mechanism S5 built and `FindProductsForLinking` uses
eight lines above in the same file. The Oracle side (S6, accepted by me) had the same shape. Swept:
the ONLY non-test consumer of the stored policy was `routing/matcher.go:45-48`. VC-3 (badge with
`only_em_estoque` off) and VC-2 (counter matching the SAME rule) could not pass. The hub verified
the sites itself and found my sweep right and incomplete in our favour — the Oracle COUNT query
took no option at all.

Ruling: **S9 owns it**, resolving the tenant policy ONCE at the routing seam where `matcher.go`
already resolves it for linking — one producer, N consumers; a separate corrective slice would
write the same port signature S9 rewrites. Write-set extended: `internal_read/ports/catalog_page.go`,
`internal_read/adapters/oracle/catalog_page.go` (predicate AND count query),
`erp_import/adapters/internalread/reader.go`. **`IncludeAll` dies from the port** — a bool beside a
policy is two mechanisms that must agree, which is F-1 applied to ourselves; "ver todos" and
`CatalogProductFactsByIDs` pass an all-inclusive policy from a named domain constructor, and the
default for an absent row lives at the `tenant_config` load seam. **Oracle rides the same patch**;
S6 does not reopen. Must-fail at contract grade: through the COMPOSED reader from root.go, flipping
a toggle on the real `tenant_config` row moves page AND count together.

The acceptance error was mine and the hub's both: we each checked that `IncludeAll` ARRIVED at the
SQL and stopped. Arrival has two halves — it reaches the consumer AND it comes from the right
producer. Filed as a gate rule: a slice that plumbs an option or a config answers both.

## S8-CONTRACT-SDK — dispatched

gpt-5.6-luna / high (standard), OS-process, prompt over stdin.
Log `scratchpad/agent__s8-contract-sdk.log`, last message `agent__s8-contract-sdk.last.md`.

The card's lane was WRONG and I measured it before dispatching rather than shipping it. It ran
from `apps/web` with two sdk-runtime files as filters; `apps/web/vitest.config.ts` does not include
`sdk-runtime`, so the command prints `No test files found, exiting with code 1`. The SDK owns
`packages/sdk-runtime/vitest.config.ts` and its own `test` script — from there the base is 5 files
/ 77 tests green. Card corrected in the same commit as this row, with the reason written next to it.

The brief carries: the repo rule that OpenAPI and `sdk-runtime` land in ONE commit; per-line
counting of files/tests with every skip named; the instruction that `No test files found` is a
filter that matched nothing, never a green, and that the fix is the directory rather than a looser
filter; and the layer distinction A-17 created — `include_all` STAYS on the wire, because "ver
todos" is a per-request screen choice, while the Go port's `IncludeAll` bool dies in S9.

## A-18 — HUB RULING: the ladder's second FE lane, and where `include_all` stops

Escalated from the S8 dispatch note above. VC-7 amended on main @`ed1b4183`.

**1. The P5 ladder runs the sdk-runtime suite EXPLICITLY, from the package directory.** VC-7 now
names TWO vitest lanes, both counted per line: workspace `@marketplace-central/web` and
`packages/sdk-runtime`. The hub's words: "green at the root" with the current script would have
been vacuous by half — the contract slice's own suite could never enter it.

**2. The root `package.json` enters NO write-set.** It is a shared seam with no owner among the
slices, so it stays with the hub and the script fix is a POST-merge step: touching it now would
drift this tree mid-flight, and the ladder must use the explicit commands either way. Recorded in
the contract with the reason.

**3. The lane rule from my S8 catch is now law in VC-7:** `No test files found` with exit 1 is a
filter that matched nothing, never a green; the correction is the directory, never a looser filter.

**4. `include_all` layering — endorsed with ONE condition.** The wire parameter stays in the HTTP
contract (per-request screen choice, S8's to publish). The condition, pinned into S9's card: the
parameter RESOLVES at the transport seam, through the named domain constructor, and only the
POLICY crosses the port. `include_all` never travels ALONGSIDE the policy into the service or the
reader — otherwise the two-mechanisms-that-must-agree shape A-17 killed at the port comes back
through the back door one layer up. S9 builds that seam; S8 only writes the wire.

## S8-CONTRACT-SDK — result: REPROVED at P4, repaired at this seat under grant A-19

Worker `20dbc321` (gpt-5.6-luna / high, OS-process). Review + correction: this commit.
Evidence: `evidence/S8-orchestrator.txt` (sections 1-8 the review as written before the ruling,
section 9 the correction).

**The lane never ran at the worker's seat.** esbuild walks upward to resolve
`packages/sdk-runtime/vitest.config.ts` and the dispatch sandbox denies the traversal:
`Cannot read directory "../../../../../../..": Access is denied.` Vitest aborted before
collection — zero files, zero tests.

The worker wrote its own zero rather than a green: "tests passed 0 observed... no skip name
exists. A hub-side re-run outside this sandbox is REQUESTed before treating the Vitest lane as
green." It did not drop the filter, did not move to `apps/web`, did not loosen the config path.
The VC-7 lane rule held under a failure the brief had not anticipated. Blind, not dishonest —
and blind still shipped, because the lane it could not run is the lane that would have stopped it.
Filed by the hub as `FINDING-sandbox-blind-fe-lane.md` (A-19 (C)); chip-locally binding from now:
a brief touching `packages/*` declares the blindness up front and names this seat as the one that
measures, and a worker's zero-observed is BLIND, never green.

**F-4 — the contract described a body and a status the server cannot emit.** Both new operations
declared the nested `ErrorResponse`, and GET declared a 404. The handler that will serve them is
`tenant_config/transport/http_handler.go`, whose `writeError` emits `map[string]string{"error":
code}` plus optional `detail` — flat — and which has no not-found branch at all; under A-17 an
absent tenant row resolves to defaults at the load seam, so GET cannot 404. Three false statements
plus one omission (no 500 on either operation, though `store.Get` failing is a certain path).

`ErrorResponse` appears 109 times in this spec, so it is not deprecated in general. The narrower
truth is measurable: its last reference before `/erp/imports:` sits at line 3184, and the region
3186→3422 was deliberately free of it, because both families living there emit flat errors. An
alien guard caught a real defect by accident of file ordering.

**A-19 rulings applied.** (A) `SellableAssortmentError`, flat `{error, detail?}`,
`enum: [invalid_body, internal_error]`; GET 200+500, PUT 200+400+500; the 404 deleted, not
softened (R-25). The surface becomes a criterion on S10's card — the handler now owes exactly
this. (A2) Both guards that sliced to `\ncomponents:` re-pointed BY VALUE, assertions untouched,
and the move said out loud. The hub's warning was confirmed by measurement first: after the
contract fix the erpImport guard read SIX 500s against its four, because the two 500s added to a
FOREIGN path sat inside its window. Inflating that count to six would have been the same collision
under another name. The window moved; the number did not. (B) Option (i) — repaired here, own
commit, glue ceiling exceeded by explicit grant.

**Must-fail: five mutations, five kills, each naming the value.** N1 restores the deleted 404;
N2 points PUT's 400 back at the nested shape; N3 smuggles an unmeasured code into the enum. N4 and
N5 are the load-bearing pair — they break something INSIDE each re-pointed window, which is the
only thing that distinguishes re-pointing from loosening; a window that had gone vacuous would
have stayed green there. Applied against a scratchpad backup and restored from it, not via
`git checkout`, because the fix was uncommitted at the time.

**Lane, counted per line at every state:** base 5 files / 77 tests green; at `20dbc321` 76 passed
2 failed; after the contract fix alone still 2 failed (the ErrorResponse assertion goes green, the
500 count fails 6 vs 4); after this commit 5 files / 78 tests / 0 failed / **0 skipped**. No run
printed `No test files found`. `tsc --noEmit` exit 0.

**Residual, written down rather than rediscovered:** other guards in this suite still bound
themselves by something other than their own subject. No currently-passing guard was re-pointed.
The next slice that appends a path or a schema should expect to meet the same shape.

## S9-CATALOG-HTTP — dispatched

gpt-5.6-sol / low (complex), OS-process, prompt over stdin.
Log `scratchpad/agent__s9-catalog-http.log`, last message `agent__s9-catalog-http.last.md`.
Tail marker `9c41be` echoed by the worker — the brief arrived with its accents intact.

**Reclassified standard -> complex, with the reason, not ad hoc at implement time.** The card was
written `standard` for 12 files confined to routing, transport and composition. A-17 added five
more across the port, both catalog adapters and their tests, and turned the slice into a
policy-plumbing job spanning ports -> two adapters -> cache -> timing -> routing -> service ->
transport -> composition, with a must-fail that has to hold through the composed reader on two
lanes at once. S7 was dispatched complex on eight files; this is seventeen. The complexity flag is
supposed to come from the plan rather than from the implement moment — so it is amended in the
plan, in the same commit that dispatches.

The brief carries, as criteria rather than advice:
  - A-17: the tenant policy resolved ONCE at the routing seam where matcher.go already resolves it
    for linking, handed as a VALUE to page and count alike; `IncludeAll` removed from the port; an
    all-inclusive policy built by a NAMED domain constructor; the absent-row default at the
    `tenant_config` load seam; `defaultSellableAssortment()` dead.
  - A-18: `include_all` stops at the transport seam — it resolves there through the named
    constructor and only the POLICY crosses the port, so the two-mechanisms-that-must-agree shape
    cannot return one layer up.
  - The CHIP-M02 503: every seat in the chain composed at root.go carries a compile-time
    `var _ ports.CatalogAssortmentReader = ...`, and a runtime assertion with a fallback at the
    HTTP seam is forbidden by name. A missing seat must be a build error, never a quiet wrong
    answer.
  - The contract-grade must-fail: through the reader COMPOSED in root.go, flipping a toggle on the
    REAL tenant_config row moves page AND count together — mirror side on the integration lane,
    Oracle side by query-text assertion on the unit lane. Both sides, not one.
  - A-15: the env dot-source line for the integration lane, RUN/PASS/SKIP/FAIL counted per line,
    and the `%v`-on-a-`*string` corollary (a must-fail that prints an address is re-injected after
    the message is fixed, not accepted).
  - The published contract as a thing to OBEY: `/catalog/products/counts` and the `include_all`
    query parameter already exist in the spec at `73190f23`. Measured divergence between what the
    handler emits and what the YAML declares is a REQUEST, never a silent adjustment of either
    side — with S8's defect named as the reason that sentence is in the brief.

## RULING A-20 — STOP-THE-LINE by CLASS, and the error-surface disposition

Two operator ratifications, registered by the hub in the profile @`1889d0dd`.

**1. Profile rule — STOP-THE-LINE de classe.** A SECOND occurrence of the same defect PATTERN — a
rule copied in dialects, a positional window, a vacuous lane, a guard living at the caller — stops
the line before another point fix. The response is: name the class, root-cause it by measurement,
and take an explicit disposition — (a) the general fix as its own immediate unit, or (b) a
registered debt with an entry criterion. **My P6 gates now evaluate by CLASS, not by instance:**
recurrence of a known class is a finding even when the instance in front of the gate is green. This
travels in the gate briefs alongside the two-proof sweep rule.

Measured against what this chip already shipped: the positional-window class hit TWICE in S8
(`erpImport.test.ts` and `activeSource.test.ts`, both slicing to `indexOf("\ncomponents:")`). The
treatment taken there was already branch (a) — both windows were made position-INDEPENDENT by a
path-indent lookahead rather than re-anchored one at a time, and the new guard for the new surface
was BORN with a value window, so the class cannot recur on the next appended path. That is the
shape the rule now demands by name; §7 of `evidence/S8-orchestrator.txt` is where it is measured.

**2. Error surface — DECIDED, single pattern.** The operator ordered unification ("tem que ser um
padrão de Erro, não esse legacy"). Disposition under rule 1: branch (a) — a dedicated unification
chip, NEXT in the queue, immediately after this chip's merge. Aborting this chip in flight was
evaluated and rejected: it would collide with S8–S10.

Practical consequence here: **nothing changes in A-19.** Flat `SellableAssortmentError` stands —
aligning with the module's own neighbour is the correct move now. It is a TRANSITIONAL standard,
and the one-line note is written in two places so a gate does not reprove a duality that already
has a disposition, and so the unification chip finds the sites this chip added by the record:
  - `BATCH-PLAN.md`, S10-CONFIG-HTTP card, under the error-surface criterion.
  - `evidence/S8-orchestrator.txt` §9.8.

## RULING A-21 — the optional request, and the grant that let me apply it

Hub ruling @`77d546d7`, on the S9 reprove in `evidence/S9-orchestrator.txt`. The reprove was
endorsed and F-2's classification under A-20 rule 1 held: "dois mecanismos / valor mágico", second
of its class in this chip. The hub's own measurement went further than the finding did — the
semantics were INVERTED, not merely awkward: `requested == AllProducts()` produced include-all, and
any explicit non-zero policy was discarded in silence, so the only value a caller could actually
communicate was the zero one and the rest of the signature was theatre.

**Ratified design.** `requested *SellableAssortmentPolicy` — nil is "the tenant's stored rule",
non-nil is "exactly this rule", HONOURED. `AllProductsAssortment()` survives as the honest value
constructor for "no cut". `DefaultSellableAssortment()` is DELETED, with the zero-reference grep
counted in the evidence.

**Invariant added by the hub, by construction.** The nil resolves ONCE, at the routing seam — the
single producer of A-17. Below routing only a CONCRETE policy travels; a site that "would need a
default" is a programming error, not a fallback. That is what `ErrUnresolvedAssortmentPolicy` and
`RequireAssortmentPolicy()` encode, and it caught my own wiring bug on its first run (23 red tests
naming the unresolved policy) before any of it reached a lane as a silent wrong answer.

**GRANT (option (i)).** Explicit scope grant from the hub: port signature + transport + the eleven
sites — larger than the orchestrator's ≤10-line glue ceiling, which is exceeded HERE and only here,
by this grant. Own commit on top of `5adeeb56`; the worker's commit is left intact.

**Four must-fails demanded, all discharged** — `evidence/S9-A21-fix.txt`:
  1. M1 re-executed post-fix on BOTH lanes WITH the env (PASS=109 FAIL=2 SKIP=0). The reason it is
     re-run with the env is the `ok` printed byte-identical under mutation by an env-less
     composition run; the bare `return` that produced it is now a named `t.Skip`.
  2. NEW must-fail: a non-nil policy named by the caller is honoured by the COMPOSED chain, at both
     the routing and the composition level. M2 reproduces the shipped defect exactly and both go
     red, naming which of the two questions the code answered.
  3. `DefaultSellableAssortment` = 0 hits in code, counted; the 10 remaining are this pack's prose.
  4. A-15 per-line accounting on every run above.

---

## S10-CONFIG-HTTP — two rounds, and the swap the first round could not see

**Worker:** gpt-5.6-luna high, OS process. Round 1 `agent__s10.log` / `.last.md`;
round 2 `agent__s10r2.log` / `.last.md`. **Committed @`eb297b44`** on `595c15e3`.
**Evidence:** `evidence/S10-orchestrator.txt`. Checklists written before each diff was read:
`p4-checklist-s10.md`, `p4-checklist-s10-round2.md`.

**Round 1: handler accepted on the first read, tests reproved.** All six hard kills passed —
`*bool` presence per field so an absent boolean is a 400 and not a silent false, the 200 body as a
re-read rather than an echo, `ErrUnknownActiveSource` to 500 with no invented all-false default,
the flat two-code error surface with no 404, both routes registered interactive, and the write set
exactly the three granted files.

**F-1 — the symmetric fixture, ANCHORS-3's class, second occurrence in this chip.** Every fixture
set `only_revenda` and `only_ecommerce_eligible` to the SAME value, and two booleans carrying the
same value are indistinguishable under a transposition. With three booleans no single fixture can
make all three pairwise distinct, so the fix is to pin them one at a time: three cases, one toggle
true each, run through BOTH directions.

**F-2 — self-referential expectation.** `want := newSellableAssortmentResponse(...)` built the
expected value by calling the mapper under test. Same shape as measure-then-write from ANCHORS-3:
both sides move together, so the assertion cannot fail for the reason it exists.

**F-3 —** the trailing-garbage guard had no test.

**The measurement, run at my seat and not taken on report.** Two independent tag swaps, copy-based
backup/restore (never `git checkout --` — the files were uncommitted):
  - M-A, request struct: PASS=20 FAIL=3 SKIP=0, exactly the two PUT cases whose value differs.
  - M-B, response struct: PASS=18 FAIL=5 SKIP=0, four cases across both directions.
  - **Round 1's test passes under BOTH.** That is what makes F-1 a defect and not a style note:
    round 1 would have shipped a live JSON tag transposition under a green lane.
  - `TestHandlerGetSellableAssortment` also passes under M-B, because it unmarshals into the
    production struct whose matching tags cancel the swap. Left as-is (the map-based table covers
    that ground) and named in the evidence so its green is never read as a wire-name proof.

**Lesson to carry: compare the WIRE, not the struct.** A response test that unmarshals into the
same struct the handler marshalled is blind to every tag defect by construction. `map[string]bool`
sees the names; the struct sees only the shape.

**A-15 again, and it paid again.** The worker's lane reported SKIP=3 with the three database tests
never executed. Named, re-run with the env: RUN=39 PASS=39 FAIL=0 SKIP=0.

**Glue I wrote myself (within the ceiling).** The package doc block. Round 2 widened its subject to
"the tenant configuration" while its body still described active-source only — including
"fail-closed 400 when unset", which is FALSE for the assortment routes. R-24's case: a total claim
false over part of its scope is made to match, not softened. The block now names both surfaces and
names the 500 asymmetry as an open question rather than letting a reader assume symmetry.

**Governance — measured, not assumed.** The worker reported `status=failed` without attribution. I
isolated the delta by removing only the `tenant_config` entry and running the lane both ways: 43
violations either side, lane already red at base; the entry ADDS
`GOV_MODULE_LAYER tenant_config-erp_import-adapters` and REMOVES
`GOV_MODULE_COVERAGE tenant_config`. The edge is one symbol — `active_source.go:74` republishes
the `erp_import` error sentinel. **OPEN:** a `temporary_exceptions` entry needs a `removal_owner`,
which is an owner identifier and not mine to invent. REQUEST sent; no entry written.

**Still parked:** the honest status for a tenant with no config row (500 interim), REQUEST open.

**Process finding.** I dispatched round 2 without snapshotting the round-1 bytes of the accepted
files, so afterwards I could not prove "only the test file changed" — the diff base is HEAD
`595c15e3`, which predates both rounds. I re-reviewed the handler and `modules.json` diffs in full
instead, which is weaker than a byte comparison and is named as such in the evidence. Snapshot the
accepted files BEFORE dispatching a follow-up round.

## S10-COND — the fail-closed 400 A-22 ordered, and the branches round 1 left unwatched

Two rounds, both GPT-5.6 Luna high, OS-process path.
Round 1 prompt `scratchpad/prompt-s10c-a22.md` → `agent__s10c.log` / `agent__s10c.last.md`.
Round 2 prompt `scratchpad/prompt-s10c-r2.md` → `agent__s10cr2.log` / `agent__s10cr2.last.md`.

### A-22 is transport-only, and the measurement is why

The hub ruled that a tenant with no config row answers **400 `unknown_erp_source`** on BOTH verbs
of `/config/sellable-assortment`, and in the same ruling RETRACTED its own A-17 sentence about a
default resolved at the tenant_config load seam. Before dispatching I measured what the ruling
actually costs: `repository.go:91` already returns the sentinel —

```go
if commandTag.RowsAffected() == 0 {
	return ErrUnknownActiveSource
}
```

— and `repository_test.go:189` already asserts it. So A-22 is a mapping change at three transport
seats and nothing below. The brief said that, and the write set was five files, additive only.

### Round 1 — accepted on the handler, reproved on the tests

The handler was better than asked: the worker mapped the sentinel at THREE seats, not two — GET's
`store.Get`, PUT's `SetSellableAssortment`, and PUT's post-write re-read. Measured at this seat with
`grep -n "getErr\|setSellableAssortmentErr"`: three construction sites, lines 18/20/41/53/88/345/363,
all passing `tenant_config.ErrUnknownActiveSource`. That grep is also what reproved the round:

**R2-1 — the third seat had no test.** Delete that branch and every lane stayed green. It is the
branch I asked for by name, because the row can vanish between the write and the read-back, and an
untested branch is a branch the next reader deletes. Second time in this same file: S10 round 1
shipped the trailing-garbage guard with no test (F-3).

**R2-2 — the generic 500 lost its only coverage.** Renaming the two interim-500 tests into the new
400 tests was correct — they asserted 500 for exactly the error A-22 calls a 400, so they were
asserting a falsehood, and a false assertion is deleted, not softened (R-24/R-25). But nothing
replaced them. After round 1 no test in the package drove a NON-sentinel store error, so every
`writeError(500, "internal_error")` was invisible. A slice that makes one path fail closed must not
leave "everything else" unwatched — that is how 500 becomes the accidental default nobody tests.

### Round 2 — four tests, and both must-fails died naming the right code

`postSetErr` added to the fake, one field, documented by the failure it catches. Then the read-back
400 plus generic-500 coverage at all three seats. MF-6 (delete the sentinel branch at the read-back
seat) and MF-7 (delete GET's 500 fallback, return 200 with the zero config so it still compiles):

```
--- FAIL: TestHandlerPutSellableAssortmentUnknownSourceOnReadBackReturnsBadRequest
    http_handler_test.go:398: status = 500, want 400 for unknown_erp_source, body={"error":"internal_error"}
--- FAIL: TestHandlerGetSellableAssortmentStoreFailureReturnsInternalError
    http_handler_test.go:420: status = 200, want 500, body={"only_revenda":false,"only_em_estoque":false,"only_ecommerce_eligible":false}
```

Each mutation isolated one failure and three passes — the tests are not reaching each other's seats
by accident. `RESTORE_HASH` matched after every mutation.

### A-5 honoured, and the brief was wrong about its own arithmetic

I snapshotted the five accepted round-1 files with sha256 BEFORE dispatching round 2
(`scratchpad/snap-s10c-r1/`), so "only the test file changed" is a byte comparison this time rather
than a re-read. That rule was born from my own round-1 failure on S10 and it paid immediately.

The brief said "five new tests" and then specified four. The worker built four, said so plainly, and
was right. Counting in prose is not a spec; the enumerated list is.

### The lane — SKIP is 0 now, and that is how the dead endpoint became visible

Counted at this seat, not taken on report:

```
ok   .../tenant_config/transport   RUN=18 PASS=18 SKIP=0 FAIL=0
FAIL .../tenant_config             RUN=7  PASS=4  SKIP=0 FAIL=3
```

The three failures are the repository tests, all refused at `127.0.0.1:55864`. `docker ps` puts the
dev-stack postgres on **5435**; 55864 is an ephemeral port of a container that no longer exists.
`MPC_TEST_DATABASE_URL` is set (length 144, value never printed).

**Finding, and it is a correction to my own rule.** A-15 makes a DB-touching brief carry the
dot-source line. That line proves the VARIABLE, not the ENDPOINT. Round 1 skipped honestly; round 2
set the variable and turned an honest skip into an environmental red. A brief that demands SKIP=0
must demand a reachable endpoint, or it just trades a silent pass for a red that says nothing about
the code. REQUEST R-1 open with the hub — the endpoint is a stack seam and I do not re-point it.

Consequence stated as a gap, not as coverage: `repository_test.go:189` is the assertion that makes
A-22 transport-only, and **it has never executed in this chip**.

### R-1 discharged — the endpoint was re-created, and `repository_test.go:189` finally ran

A-25 granted `pg-session-up` in this worktree, on the reasoning that a dead endpoint is RE-CREATED,
not re-pointed, and that the prohibition on booting a server is about the DEV STACK, not this
per-checkout facility. 5435 was never touched. Full measurement in
`evidence/S10-COND-db-seat.txt`; the load-bearing lines:

```
container=mpc-pg-session-3eee515d port=55840 pw_len=48
CREATE DATABASE      attempt=1 exit=0
go run ./cmd/testdb migrate  ->  applied 72 migration(s)   (tree holds 72 *.sql — same inventory the lane computes)
go run ./cmd/testdb migrate  ->  applied 0 migration(s)
go test ./internal/modules/tenant_config/... -v   RUN=43 PASS=43 SKIP=0 FAIL=0
--- PASS: TestRepository_SetSellableAssortment_RoundTripPersistsPerTenant (0.01s)
```

The database name is not free-form: `internal/testsupport/postgres/target.go:19` pins
`^mpc_test_[0-9a-f]{32}$` and `LoadConfig` returns `HPG_TARGET_INVALID` for anything else. The
placeholder the scratchpad env line carried had 36 hex characters, so it failed that regex as well
as not existing.

Then, beyond the grant, two must-fails on the sentinel, because PASS is not evidence that an
assertion asserts. Dropping the `RowsAffected()==0` guard:

```
repository_test.go:190: SetSellableAssortment() missing tenant error = <nil>, want ErrUnknownActiveSource
```

Scanning the three assortment columns into throwaways so `Get` returns the zero policy — the one
worth naming, because it fails on the READ side with the write path untouched, which is what
separates "the columns round-trip through Postgres" from "the struct round-trips through Go":

```
repository_test.go:94:  first tenant OnlyRevenda = false, want true
... 8 lines ...
repository_test.go:186: after source switch OnlyEcommerceEligible = false, want true
```

Both reverted byte-for-byte; `git diff` on `repository.go` is empty and RESTORED-GREEN passes.

### FINDING for the hub — the sentinel is in NO harness lane

`harness.ps1 integration` was run FIRST and it passed: `status=passed`, `migrations_first=72`,
`migrations_second=0`, `port=55840`. It ran **zero** tenant_config tests.

`Get-HarnessIntegrationTestPackages` (`scripts/harness/Postgres.psm1:42-61`) discovers module
packages by scanning the FIRST FIVE LINES of every `internal/modules/**/*_test.go` for
`^//go:build\b.*\bintegration\b`. No file in `internal/modules/tenant_config` carries that tag.
So the integration lane never compiles the package; the unit lane (`./tests/unit/...`) never reaches
`internal/modules` at all; and a bare `go test ./...` without the env var reaches it and skips clean
by design (`target.go:47-52`). Three lanes, three silences — a DB-touching test with no build tag is
invisible to the lane that provisions a database and skip-clean everywhere else, so it can rot green
indefinitely.

Same shape as HARNESS-DEBTS B-5 one level up: there the variable was set and the database was
absent; here the lane is green and the test was never compiled. Both make a green a statement about
something other than what the reader thinks. Reported, not fixed — the lane's package list is a
seam with an owner.

Also measured: the lane DROPS its per-run database on the way out. After `status=passed` the session
container held only `postgres/template1/template0`, which is why this seat had to create and migrate
its own.

### The same silence hid something bigger, in `internal/composition`

`root_test.go:47` is `TestRootComposedCatalogReaderMovesPageAndCountWithStoredPolicy` — the test
that proves the STORED assortment policy moves both the page and the count through the composed
reader, which is the centre of this chip. It had been SKIPPING in every composition run of this
chip:

```
root_test.go:49: MPC_TEST_DATABASE_URL is unset; the composed catalog reader needs the integration database
--- SKIP: TestRootComposedCatalogReaderMovesPageAndCountWithStoredPolicy (0.00s)
```

`internal/composition` is in no lane's package list either — the integration lane adds
`./tests/integration` plus tagged `internal/modules` packages, and composition is neither. Nobody
had ever given it a database. Against the live target it is clean and the policy test EXECUTED for
the first time:

```
go test ./internal/composition/ -v   RUN=72 PASS=72 SKIP=0 FAIL=0
--- PASS: TestRootComposedCatalogReaderMovesPageAndCountWithStoredPolicy (0.09s)
```

Two packages, one cause. The finding above is not about tenant_config; it is about the seam.

### Committed

`ba91ad2` — handler (three seats), OpenAPI + SDK contract test in the same commit, the
`module-edge-tenant-config-erp-import-adapter` registry entry with `removal_owner: CHIP-ERROR-UNIFY`,
and the BATCH-PLAN retraction of the A-17 load-seam sentence in both places it appeared.

FE lanes re-run at this seat after the commit: `@marketplace-central/sdk-runtime` 78 passed / 78,
0 skipped — the contract test carrying the `enum: [unknown_erp_source, invalid_body, internal_error]`
and the "400 on GET and PUT" assertion is green against the shipped spec.


## S11-INTEGRACOES-FE — two rounds; round 1's green could not see its own missing refetch

Dispatched to GPT-5.6 Luna high (standard implement), prompt at `scratchpad/prompt-s11-r2.md`,
log/last-message at `scratchpad/agent__s11r2.*`. Owned files: `IntegracoesPage.tsx`,
`IntegracoesPage.test.tsx`, `packages/web-query/src/activeSource.ts`, `packages/web-query/src/index.ts`.

### Round 1 — REPROVED at this seat, three reasons

R2-1: the counts fake was a CONSTANT. A card that never refetches after a toggle and a card that
refetches correctly return the same constant, so the test could not tell them apart — the same
both-worlds shape as CHIP-IMPORT-CHAIN. Round 2 makes the fake state-derived
(`IntegracoesPage.test.tsx:142-144`):

```ts
getCatalogAssortmentCounts.mockImplementation(() =>
  Promise.resolve({ sellable_count: storedAssortment.only_revenda ? 2 : 3, total_count: 4 }),
);
```

R2-2: the two binding tests did not name which toggle bound to which field, so a swap of two
booleans passed. Round 2 names them (`:160-162`), e.g.
`expect(revenda.checked, "Somente produtos de revenda must bind to only_revenda").toBe(true)`.

R2-3: the no-active-source state was unhandled — the card rendered as though the server had chosen
an empty rule.

### Round 2 must-fails, verbatim from the worker and re-run at this seat

```
FAIL ... recalculates the result line after a toggle instead of showing the pre-change count
TestingLibraryElementError: Unable to find an element with the text: Resultado: 3 de 4 produtos.
Tests 1 failed | 17 passed (18)
```

```
AssertionError: Somente produtos de revenda must bind to only_revenda: expected false to be true
```

The first RED is the one that matters: it is the missing `invalidateQueries` that round 1's constant
fake could not see.

### Verified at this seat, not on report

- The stateful fake is genuinely state-derived (read at `:131-144`), not a two-value script.
- The no-source state is discriminated EXACTLY, not by a status class
  (`IntegracoesPage.tsx:302-306`): `status === 400 && code === "unknown_erp_source"`.
- **C-5 (compare the WIRE, not the struct):** the discriminator was checked against the SDK's real
  producer rather than against the neighbouring mock. `getSellableAssortment` (`index.ts:1909`) and
  `getCatalogAssortmentCounts` (`:1871`) both delegate to `getJson` (`:1719`), which throws the FLAT
  object at `:1723`:
  `throw { status: response.status, error: (data as ErrorResponse).error } satisfies MarketplaceCentralClientError`.
  Not an `Error` subclass, not a nested body — exactly what the page destructures.
- A-23 honoured on both counts: invalidation is the targeted `queryKeyNamespaces.catalog` prefix
  (`activeSource.ts`), never active-source's blanket `invalidateQueries()`; and `erp_source` is in
  the counts key (`index.ts` `catalogQueryKeys.counts`), so a source flip cannot reuse the previous
  source's count under a new label.
- Copy is the operator's pt-BR with no dev marker: `Nenhuma fonte definida ainda — escolha a fonte
  que o app vai ler.` Provisional per A-23.

### Lanes, counted at this seat

```
@marketplace-central/web  src/pages/integracoes   18 passed / 18   0 skipped
@marketplace-central/sdk-runtime                  78 passed / 78   0 skipped
```

`@marketplace-central/web-query` has **no `test` script** (`npm error Missing script: "test"`), and
my first reading of that was wrong: it is not a missing lane. The package has three test files of
its own and `apps/web/vitest.config.ts` includes
`"../../packages/web-query/src/**/*.test.{ts,tsx}"`, so they run in the
`@marketplace-central/web` lane. Only the per-package script is absent. The filtered run above
excluded them on purpose (`-- --run src/pages/integracoes`); they are counted in S13's full lane.

Flag for S12 while I was in that config: `feature-products` is included by EXACT FILENAME
(`../../packages/feature-products/src/CatalogPage.test.tsx`), not by a glob. A second test file
added anywhere in that package would run in NO lane and report a silent zero. S12's write set
touches only `CatalogPage.test.tsx`, so the slice is unaffected — but its brief must say so, because
"add a test file" is the obvious thing a worker would reach for.

### Contamination cleared under A-25 R-2(a) before any of the above was trusted

The round-1 worker claimed the standard vitest lane was blocked and that no emitted `.js` remained.
Both refuted by measurement here: the lane works, and **186** untracked `.js`/`.js.map` files with
same-basename `.ts`/`.tsx` siblings were present, shadowing sources under extensionless resolution.
Root cause was my own `tsc -b` during the S11 review. Deleted per-file under the grant
(`DELETED_COUNT=186`, the single skip being my own untracked `catalog_routes_test.go`, which is not
a `.js`), tracked files untouched, no `git add -A`. The S11 lane was then RE-RUN post-purge and
reproduced 18/18 with 0 skipped — so the green is not read off shadow JS. `.gitignore`/`tsconfig`
remain the hub's seam per R-2(b); not proposed here.
