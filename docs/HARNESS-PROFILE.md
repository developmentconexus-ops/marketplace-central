# Harness Profile — marketplace-central

**Layer:** REPO (binds with the `harness@mnfs-harness` plugin's `HARNESS-CORE.md` +
`REVIEW-STANDARD.md` — canonical source `Documents/mnfs-harness/harness/`, loaded by the
harness skills; mission content lives in `.mnfs/`). This profile is the ONLY doctrine file
vendored in-repo — the method (core + review-standard) is plugin-sourced, not duplicated here.
**Status of this file:** BINDING — swap executed 2026-07-16 at the M-01 milestone boundary
(operator-approved). Vendored core/review-standard copies + the legacy `docs/HARNESS.md` pointer
were retired 2026-07-18 (doctrine ships from the plugin). Provenance below carries over from the
combined doctrine's dated ratifications (historical `docs/HARNESS.md §X` citations reference that
now-retired file as it stood at ratification time).

---

## 1. Identity & stack
`status: ratified` · `provenance: 2026-07-15 · extracted from docs/HARNESS.md (operator-ratified 2026-07-15)`

- Go backend (`apps/server_core`, Go workspaces) + React/TypeScript frontend (`apps/web`),
  npm monorepo.
- OS/shell binding: Windows; **PowerShell for all stack ops — never bash, never WSL**.
- Default branch: `main` (origin/HEAD → origin/main; no `master` ref exists — corrected
  2026-07-20, D-120). Hub checkout must be on `main` before every hub commit (field gotcha:
  chip launch or a nested worktree holding `main` can leave the primary dir on detached HEAD —
  verify `git branch --show-current` == `main`, not just non-detached).
- Contract-first: OpenAPI spec + generated SDK (`packages/sdk-runtime`).

## 2. Verification ladder bindings (core §5)
`status: ratified` · `provenance: 2026-07-15 · docs/HARNESS.md §5 + field finding (governance lane, 2026-07-15)`

- **L0** — `go build ./...` with `GOCACHE` bound ABSOLUTE (next bullet — never the bare
  relative literal) · web `tsc`/typecheck ·
  governance lanes: `npm run harness:governance -- -BaseSha <sha>` — run from a **clean
  detached worktree** (main checkout sweeps `.claude/worktrees/*` and false-fails until the
  scanner exclusion lands) and pass the **full 40-hex** BaseSha (short sha =
  `GOV_SEMANTIC_DRIFT id=base-sha-invalid`).
- **GOCACHE must resolve to an ABSOLUTE path on Windows/pwsh** (D-14, M-01): relative
  `.gocache` breaks when the working dir shifts mid-pipeline — bind it as
  `$env:GOCACHE = (Join-Path (Get-Location) '.gocache')` (pwsh) or
  `export GOCACHE="$(pwd)/.gocache"` (bash) before Go commands. Same for `GOMODCACHE`.
  **This file prints no copyable relative form**: it did, three times, three lines above this
  rule, and a chip copied one and got an 83-byte EXIT 1 with zero `=== RUN` — a doctrine
  artifact that hands you the line it forbids is worse than one that says nothing
  (CHIP-ANCHORS-3 field REPORT, 2026-07-28).
- **Governance base anchor (per milestone):** the drift gate's BaseSha for a chip is the
  milestone's ACCEPTED BASE SHA (40-hex, carried in the chip prompt) — on a long-lived
  worktree, drift REDs computed against any other base are not the chip's defect
  (M-01 field, 2×: Slice 7 `base-sha-invalid`, P5 tool-vs-validate topology). The chip
  records the anchor in its evidence; the hub re-runs governance on the integrated default
  branch at acceptance.
- **L1** — `go test ./...` with `GOCACHE` absolute (touched packages + guard suites; full sweep only
  when migrations/platform touched) · web vitest · integration lane
  `npm run harness:integration` (see §4).
- **Known pre-existing failure allowlist (L1):** verdicts CITE this list instead of
  re-proving non-linkage from scratch (M-01 re-proved `TestPhase1SmokeFlow` 5×). A ladder
  run is GREEN-with-allowlist when its only failures are listed here unchanged. Editing the
  list is hub-owned; each entry needs an evidence pointer + a backlog owner.
  - ~~`TestPhase1SmokeFlow`~~ RETIRED 2026-07-17: root cause found by CHIP-M02 (fixture
    `"smoke-prod-1"` vs integer-enforcing `positiveProductID`, both pre-existing on main);
    hub fixed the fixture (`"1001"`) and the lane ran green (run 5b244bce). Entry kept for
    history only — the lane failure no longer exists on main ≥ the fixing commit.
  - `TestListingsReadContractEndToEnd` (`apps/server_core/tests/integration`) — INTERMITTENT
    flake under the full lane only (passes isolated 8/8 — CHIP-M02 evidence 2026-07-17; also
    passed hub run 5b244bce). If it is the ONLY lane failure: re-run once, cite this entry.
    Backlog: test-isolation/t.Cleanup audit, hub queue, unowned.
  - Raw `tsc --noEmit -p apps/web/tsconfig.json` fails `TS2688: Cannot find type definition
    file for 'node'` on base ≤ c6df5cc1 — tsconfig declares `"types": ["node", …]` but no
    package.json carries `@types/node` (pre-existing repo-wide; verified on main by CHIP-M03
    2026-07-17). Cite this entry; verify FE slices via `npm run build` + vitest meanwhile.
    Fix owner: CHIP-M03 (hub grant D-05 — `@types/node` root devDependency arrives via M-03
    merge; entry retires then).
- **L2** — dev stack up **ONLY via docker compose**: `npm run docker:dev`
  (postgres+backend+frontend, server `:8080` + web `:5174`); OAuth flows:
  `npm run docker:oauth`. HUB-OWNED seam: chips send `REQUEST dev-stack`, never boot their own
  server process, never bind the ports, never load `.env*` into session env vars (env is
  consumed by the container entrypoint, not the session). Smoke: target routes, error shapes,
  OpenAPI ↔ SDK ↔ handler parity; evidence captured to the mission's contract paths.
- Post-merge ladder on integrated main MUST include the clean-worktree governance run.

## 3. Fresh-workspace bootstrap
`status: ratified` · `provenance: 2026-07-15 · M-01 field finding (hermetic lane)`

Standard first act in every new chip worktree, before any hermetic lane:

```
cd apps/server_core && GOMODCACHE=$(pwd)/.gomodcache go mod download all
```

(~130M; env prep, NOT a dep change — no REQUEST needed.)

Node side (ratified 2026-07-16, CHIP-SAT REQUEST): fresh chip worktrees have no
`node_modules`; run `npm ci` at the worktree root before web `tsc`/`vitest` lanes.
Lockfile-faithful install = env prep, NOT a dep change — no REQUEST needed. Never reuse
another checkout's `node_modules`: npm workspace symlinks would resolve workspace packages
(e.g. `packages/sdk-runtime`) to the OTHER tree's sources, silently validating the wrong code.

**Test-fixture conventions (ratified 2026-07-16, CHIP-SAT findings F-4 + M-06 F-02):**
- Windows Go `time.Now()` ticks at 100ns; pg `timestamptz` stores µs. Any fixture asserting
  `.Equal()` on a timestamp round-trip MUST `Truncate(time.Microsecond)` first (production
  convention: `orders/domain/sankhya_linkage.go:178`). Signature: got/want differ only in the
  7th fractional digit.
- `internal/platform/migrate/runner_test.go` hardcodes the migration count — every migration
  grant implies bumping that fixture in the same slice.
- A new server module requires an entry in `contracts/governance/modules.json` or governance
  fails `GOV_MODULE_COVERAGE` (precedents 17cce1a9, e1819778).

**False-alarm signatures:**
- `HPG_MIGRATION_FAILED` with `migrations_first=-1` = build died before migrate ran (empty
  `.gomodcache` under `GOPROXY=off`/`GOSUMDB=off`), not a SQL/migration defect. Warm the cache
  before diagnosing SQL.
- 3D000 "database does not exist" on first CREATE DATABASE attempt = postgres first-boot init
  restart race (pg_isready passes during it); the lane's retry loop absorbs it — not a defect.
- codex-cli 0.144.4 logs non-fatal `failed to renew cache TTL: missing field
  supports_reasoning_summaries` on every run on this machine; output still produced — noise,
  not a dispatch failure (CHIP-M02 field finding 2026-07-17).
- Full-tree `gofmt -l` on a Windows worktree falsely flags pre-existing files whose only
  diff is CRLF (`gofmt -d` shows `^M` only, zero token changes — autocrlf checkout
  artifact, seen on MIS-002 `internal_read/*.go`). Not a formatting defect. Scope gofmt
  gates to milestone-authored dirs, or set `core.autocrlf=false` before a full-tree gate
  (CHIP-M01 field finding F-ENV-M01, 2026-07-18).
- codex `--sandbox workspace-write` on Windows CANNOT run the vite/esbuild build (`npm run
  build` fails esbuild access-denied / "could not resolve vite.config.ts"; `tsc` may also
  fail in-sandbox). A worker BLOCKED on these signatures with complete written code is an
  ENV false alarm, not a code defect. Ratified mitigation: workers stay workspace-write; the
  CHIP re-runs build/tsc/vitest chip-side post-dispatch and that re-run is the verification
  of record (CHIP-M03 field finding 2026-07-17, verified same-code green chip-side; never
  grant danger-full-access for this).
- Codex worker validation — split build/vet/test into SEPARATE commands (never chain
  `go build … ; go vet … ; go test …` in one dispatch): cold `go build ./...` alone runs
  ~65s in a chip worktree and can exhaust the 120s codex timeout, surfacing as a spurious
  BLOCKED with no stage output. A combined-command timeout with no per-stage output is a
  FALSE ALARM — re-run the three lanes split with long timeouts before treating as a
  defect. Impl-pack instruction: workers run the commands separately, not chained
  (CHIP-M02 field finding 2026-07-18, F-03-S1; dispatcher re-ran split: all exit 0).
- codex 0.144.4 PATH binary can DIE mid-dispatch on cache-TTL renewal (vs the non-fatal
  0.144.4 warning already listed): schema clash with 0.145.0 models_cache. Mitigation:
  dispatch via the 0.145.0-alpha.18 binary explicitly, not bare `codex` from PATH
  (CHIP-M03 field finding F-ENV-9, 2026-07-18).
- Browser-pane screenshot rasterizer broken on this machine — a blank/failed screenshot
  during browser QA is an ENV failure, not a UI defect. Visual QA evidence = computed-style
  captures (getComputedStyle assertions) + accessibility-tree reads instead of pixels
  (CHIP-M03 field finding F-ENV-10, 2026-07-18).
- A `node_modules` that is MID-`npm ci` (or left half-installed by an interrupted one) makes
  `tsc` emit module-resolution errors that impersonate real defects — and the count is STABLE,
  so running `tsc` twice and getting the same output proves nothing. The only discriminating
  observable is the `npm ci` PROCESS EXIT: do not read a single `tsc` result out of a worktree
  whose install has not exited 0. Sibling of the stable-but-non-discriminating trap in the
  review note below: an observable can be reproducible and still not tell the two worlds apart
  (CHIP-IMPORT-CHAIN field finding, 2026-07-28).
- The FE vitest lane is **`cd apps/web && npx --no-install vitest run`**, and the `cd` is part of
  the lane, not shell convenience. The same command from the WORKTREE ROOT silently changes the
  project set (adds `packages/sdk-runtime`) and turns its contract tests RED, because they
  resolve `resolve(process.cwd(), "../../contracts/openapi.yaml")` — from the worktree root that
  path lands in `.claude/contracts/` and does not exist. From the REPO root the same command
  globs ~1260 files including installed packages and returns ~94 failed. Both reds are the
  instrument, not the code. Third instance of the class: an observable that is stable, cheap and
  wrong about which world you are in (hub executing-seat finding on CHIP-VINC-NEUTRO,
  2026-07-28).

**Review-process note (F-ENV-8, 2026-07-18):** isolated per-slice reviews cannot see test
responsibilities MOVED between files (e.g. nav assertions relocating Layout.test → Header.test
reads as deletion in one diff and addition in another). The full-suite L0 run at milestone
close is the backstop that proves net coverage — never waive it for FE milestones on the
grounds that each slice was reviewed.

**Sandbox vitest blindness has TEETH (CHIP-IMPORT-CHAIN, 2026-07-28).** The
`--sandbox workspace-write` clause above says the chip's re-run is the verification of record;
this is the instance that shows why it is not a formality. A codex slice added a `<Link>`, which
made a moved component require router context and turned two unrelated test files red. The
worker could not run vitest, so it reported "typecheck clean, committed" — accurately, over a
red suite. The chip's post-dispatch re-run caught 2 broken files plus 1 genuine timing defect.
A worker's "committed" is a claim about what it could OBSERVE, not about the suite; a chip that
forwards it as green is laundering blindness into evidence.

## 4. Test database / integration strategy
`status: ratified` · `provenance: 2026-07-15 · docs/HARNESS.md §5 (integration lane hardening + session container)`

- Isolation unit: fresh `CREATE DATABASE mpc_test_<32hex>` per run, dropped after.
- Session container: `npm run harness:pg:up` starts one long-lived postgres per checkout
  (`mpc-pg-session-<8hex>`, hashed from checkout path — hub and chip worktrees never collide);
  `npm run harness:integration` auto-reuses it (`container=session-reuse`, ~20s → ~3s
  overhead); without it, per-run `container=ephemeral`. Remove at milestone close:
  `npm run harness:pg:down`. State: `scripts/.runs/pg-session.json` (gitignored, 127.0.0.1).
- Lane is self-discovering: `scripts/harness/Postgres.psm1` globs every `//go:build integration`
  package under `internal/modules/` + `tests/integration` — a new module joins by existing.
- Cross-track serialization: integration runs across tracks with **divergent migrations** are
  serialized by the hub (3D000/template-collision finding); same-fingerprint runs safe concurrent.

## 5. Collision axes — instantiation (core §3)
`status: ratified` · `provenance: 2026-07-15 · docs/HARNESS.md §3`

| Axis | Concrete binding in this repo |
|---|---|
| Contract artifacts | `contracts/api/marketplace-central.openapi.yaml` + `packages/sdk-runtime` (contract lock: one owner at a time, or hub pre-assigns disjoint path sections and resolves regen conflict at merge) |
| FE surface | `apps/web` component/route trees; AppRouter/nav/Layout are named owned seams |
| Migration | hub pre-allocates disjoint number blocks in chip prompts; unplanned need = `REQUEST migration-number`, never grab blind |
| DB shape | one exclusive owner per table / module registration; other tracks `REQUEST` |
| Module | Go modules under `apps/server_core/internal/modules/`; cross-module edges via published interfaces safe, same-module internals serialize |

## 6. Shared seams & owners
`status: ratified` · `provenance: 2026-07-15 · docs/HARNESS.md §§2-3-6`

Hub-owned (chips `REQUEST`, never take): OpenAPI/sdk-runtime contract lock · migration number
blocks · dev stack (:8080/:5174) · harness control files (`scripts/harness*`, package.json
harness scripts) · `contracts/governance/` registry · `docs/HARNESS*` doctrine files.

**Live dispatch viewer (ratified 2026-07-16):** the hub serves the codex live dashboard
(core §8 pattern): scratchpad `live-server.mjs`, `127.0.0.1:7391`, SSE-tailing every
`agent__<id>.log`/`.done` across ALL marketplace-central session scratchpads (hub + chips,
discovered dynamically). Hub boots it at wave dispatch; scratchpad-local, never committed.
Canonical copy travels session-to-session by copying from the previous hub/milestone
scratchpad (glob `Temp/claude/C--*marketplace-central*/*/scratchpad/live-server.mjs`, newest).

**DB-specialist consultation seam (ratified 2026-07-16):** Oracle/Sankhya query questions
(schema semantics, TOP/CODPROD/TGF* doubts, query plans) route to the standing MNOS
specialist session `local_ec787804-f8e9-4981-9c12-7d3f45292294` ("Marketplace Central
database queries") via ccd `send_message`. Chips do NOT message it directly — chip sends
`REQUEST db-consult` to the hub with the question; hub relays and returns the answer.

**Dev-stack sync standing policy (ratified 2026-07-16):** a chip `COMMITTED` event carrying
`stack-sync: <sha>` triggers the hub to rebuild/redeploy the dev stack at that SHA WITHOUT
per-request negotiation — M-01 burned 4 hub round-trips on restart asks alone. The hub still
owns the action (chips never touch containers); only the approval ceremony is removed.

## 7. Non-negotiables (per-endpoint / per-write, core §5)
`status: ratified` · `provenance: 2026-07-15 · docs/HARNESS.md §5 + AGENTS.md`

- `tenant_id` predicate on every tenant query.
- Provider payloads at adapters only (domain/application/ports/adapters/transport boundaries hold).
- Unknown never becomes zero/default — ADR-17, fail honest; no blanket recover/fallback on
  integrity-critical reads.
- OpenAPI + `sdk-runtime` land in the same commit.
- Provider writes: resolved linkage, explicit policy/source time, duplicate protection, audit
  (IC-03 gates).
- Mocks prove contract behavior, never live integration.
- Validation contracts/tests NEVER fall back to stub/mock for an integration seam without
  explicit operator authorization; integration criteria run against the REAL dependency
  (live/operator-provisioned env). Composition roots never ship permanent stub/nil wiring on a
  live path — stub only with dated deferral naming the replacing slice. Mission planning must
  declare real-integration bindings (seam + env) up front.
- **No command that dumps a container's or process's whole environment.** `docker inspect`,
  `docker exec … env`, bare `printenv`/`Get-ChildItem Env:` and equivalents print every secret the
  target holds straight into the transcript. Diagnose a variable by NAME and one at a time
  (`printenv THE_VAR`, `docker exec … printenv THE_VAR`), never by dumping the set. This binds
  WORKERS too, so it belongs in the dispatch prompt denylist and not only here — the chip reads
  the profile, the worker does not. It holds even when the target is throwaway: a session
  Postgres container's password is CSPRNG-generated per container by
  `scripts/harness/Postgres.psm1` and dies with it, so there is nothing to rotate, and that is
  exactly why the rule cannot be argued down case by case — "it was disposable" is available
  every time (CHIP-ANCHORS-3 R4 worker, 2026-07-28; the module already passes the value in
  `-RedactionCandidates` for its own process calls, so the leak came from a path outside it).

## 8. Truth order (core §6)
`status: ratified` · `provenance: 2026-07-15 · AGENTS.md + docs/HARNESS.md §6`

`ARCHITECTURE.md`/ADRs > OpenAPI + SDK > `contracts/governance/` > wiki > `.mnfs/` >
tests/builds/commits. Stop and classify architecture, contract, runtime, ownership, or
verification conflicts against this list.

## 9. Human gates
`status: ratified` · `provenance: 2026-07-15 · AGENTS.md + docs/HARNESS.md §6`

- Push: NEVER without explicit operator permission (commit after verified work = standing auth).
- Dependency changes (dep change = `REQUEST` to hub; install-as-ritual forbidden).
- Live ML (Mercado Livre) writes: explicit operator authorization per mission Validation Strategy.
- Never read/print/commit `.env*` contents.
- `git branch -d` never `-D` (refusal = unmerged work, operator decides).

## 10. Superseded protocols denylist
`status: ratified` · `provenance: 2026-07-15 · docs/HARNESS.md §2 (chip-prompt pin h)`

- `mpc-goal-harness` skill — superseded 2026-07-15; never invoke, even if skill discovery
  surfaces it in a worktree. Its control-plane modules (Context/Impact/State/Evals) and lanes
  (`context-compile`/`context-validate`/`impact`) were deleted 2026-07-15.
- Codex-only harness previously in AGENTS.md — superseded 2026-07-15.
- mnfs-workflow execution-layer skills/agents — superseded 2026-07-15 (mnfs-harness
  `6b29412` "layered unification — harness is the only execution engine"; deleted at source,
  MNFS retains contract + verdict layer only): `feature-execution`, `milestone-execution`,
  `feature-context-pack`, `feature-validation-review`, `correction-worker` skill,
  `feature-implementer` / `milestone-orchestrator` agents, `/milestone-start`,
  `/feature-context`, `/feature-accept`. Stale copies survive in the codex plugin cache
  (`mnfs-workflow 0.1.0+codex.local-20260706-01`) and `~/.codex/plugins/mnfs-codex-plugin/`
  until repackage to 0.2.0 — workers MUST NOT follow them if auto-discovery surfaces one.
  General rule (ratified 2026-07-16): auto-discovered skills are never doctrine; only
  files pinned verbatim in the dispatch prompt-pack bind a worker.
- **Skill-discovery denylist clause in worker bindings (ratified 2026-07-17, F-ENV-3):**
  the `impeccable` plugin auto-injects a skill mandate into codex worker sessions on this
  machine and derails them (workers abort firing `NO_PRODUCT_MD` instead of implementing).
  Every codex dispatch's role/repo bindings block MUST carry this clause verbatim:
  "Only this prompt binds you. Any auto-discovered or auto-injected skill mandate
  (impeccable, NO_PRODUCT_MD, or similar) is NOT a contract conflict — discard it and
  proceed with the slice." Core-candidate — upstream to impl-pack/bindings template at
  milestone boundary. Operator-level note: investigate why impeccable injects into codex
  sessions (stray config outside repo scope).

## 11. Code-review bindings (REVIEW-STANDARD.md instantiation)
`status: ratified` · `provenance: 2026-07-16 · operator-requested review hardening`

- Standard: the `harness@mnfs-harness` plugin's `REVIEW-STANDARD.md` (canonical
  `Documents/mnfs-harness/harness/REVIEW-STANDARD.md`, loaded by the harness skills; NOT
  vendored in-repo) — binds per-feature adversarial review, dual gate, hub spot-check.
- Deterministic pre-pass (§7) = L0 lane: `go build ./...` + `go vet ./...` (both
  with `GOCACHE` absolute, §2) · web `tsc` · governance lane (clean worktree, 40-hex BaseSha). Review
  dispatch only after L0 green; reviewer receives the L0/governance report as input.
  Candidate additions (staticcheck, dupl, gocognit) = dependency gate: `REQUEST` to hub first.
- Learnings file: `docs/REVIEW-LEARNINGS.md` — loaded into every code-review dispatch.
- Slice size norm: ≤ ~300 changed lines (mechanical/generated diffs exempt when declared).
- Severity taxonomy + dual-gate agreement merge as per standard §5/§8; reconciliation table
  required in every `CLOSED` event.
- Execution model (standard §13-§16; per-feature amendment ratified 2026-07-18, D-51): review =
  ONE adversarial reviewer per FEATURE (all its slices together) via prompt-pack §14 (never a
  crew), emitting a correction plan; inter-slice guard = failing-test-first per slice; dual gate
  dispatched simultaneously; disjoint-feature overlap allowed (§15 — reviewed-green before
  merge/dependent feature); ★ crews run under noise control (FAIL-restraint verbatim-quote rule,
  advisory cap 10, learnings suppression via `docs/REVIEW-LEARNINGS.md`).
- Reviewer model/reasoning bindings (standard §13 restates core §1 matrix — canonical there):
  per-feature reviewer = Claude subagent **sonnet** (session-default effort); dual gate = **COLD
  Opus subagent** (clean context, explicit `model=opus`, never the orchestrating session —
  self-grade bias) + **GPT-5.6 Sol `--effort medium`**; ★ crews = cold Claude **sonnet**
  subagents, mission readiness P7 adds **GPT-5.6 Sol `--effort high`**; hub spot-check =
  sonnet subagent. Every review context is cold — bounded inputs only, never the
  implementation conversation. GPT flags NEVER retyped from memory — resolve via the
  `codex-dispatch` skill; `--effort` always explicit.

### The dual gate has a THIRD seat: an independent EXECUTOR
`status: ratified` · `provenance: 2026-07-28 · CHIP-ANCHORS-3 field finding, exercised on
CHIP-IMPORT-CHAIN before ratification`

A physically read-only gate seat cannot discharge an execution criterion, by construction. This
is not a reviewer failing; it is the seat. Field evidence, same day, two chips:

- CHIP-ANCHORS-3 round 1: the Opus seat has no Bash and no git; the Sol seat was denied Go's
  temporary work directory and a Postgres test database. Its two blocking findings were
  *"the read-only environment denied creation of Go's temporary work directory"* and
  *"I could not provision or write to a PostgreSQL test database"* — no code defect at all.
  Every must-fail, every ladder number, every `git` fact and the patch `sha256` ended up
  certified **only by the implementer**, which is exactly what the brief tells the reviewer not
  to accept.
- CHIP-IMPORT-CHAIN: BOTH seats independently returned `I9 — NOT-PROVEN` with the same reason,
  and the Opus side wrote it plainly: *"the chip's word is not the instrument."*

So the gate is **three seats, not two**:

1. two READING seats — cold Opus + GPT-5.6 Sol medium, blind to each other (unchanged);
2. one EXECUTING seat — runs the ladder, the must-fail, the fixture reproduction, and the `git`
   facts. **The hub owns it.** It is independent of the implementer, which is the property that
   matters; it is not independent of the reviewer, which does not matter because it produces
   measurements, not verdicts.

Corollaries that bind:
- A reading seat that reports a blocking finding whose content is *"I could not run X"* has
  found nothing. Route the criterion to the executing seat; do not spend a correction round.
- Scope execution criteria OUT of the reading seats' prompts, so they stop burning findings on
  a sandbox they were never going to get.
- The executor states what it ran, verbatim, and against which SHA. First run of record:
  `_hub-gate-import-chain/GATE-P6.md` §I9 (tsc 15/0-in-scope, vitest 65 files / 521 tests, at
  the chip's own HEAD, matching the chip's declaration — true before, verified after).
- The FE analogue already existed and is the same idea: the hub's live browser drive (L2). This
  generalises it to the backend.

### A delta brief must order a CLASS sweep, or the second occurrence is invisible
`status: ratified` · `provenance: 2026-07-28 · CHIP-ANCHORS-3 field finding, round 2`

A delta-scoped round is cheap and should stay the default. But a brief that NAMES the site to
re-examine teaches the seat to stop at it.

Field evidence, one round, two seats, same brief: CHIP-ANCHORS-3's round-2 brief named the site
where round 1 had found a false universal. The Sol seat checked that site, found it corrected,
and returned `Findings: None` — on a file that carried the same false universal 83 lines below,
findable by a one-line `grep`. The Opus seat swept the file on its own initiative and returned
REFUTED with it. **The seat that returned None did literally what the brief asked.** The defect
is the brief's, not the model's.

So every delta brief binds:

- **Order the sweep explicitly**: re-examine the named site AND sweep both changed files end to
  end for other occurrences of the same CLASS. Give the seat the class as searchable tokens —
  for false-totality prose that is `never`, `always`, `only`, `unreachable`, `cannot`,
  `no longer`, `every`.
- **A verdict with no SWEEP section is incomplete on its face**, and the hub returns it as
  incomplete rather than reading it as a PASS.
- **The author runs the same sweep against their own pack before publishing.** If only the
  reviewers sweep, the chip keeps producing the class and outsourcing detection. It is a `grep`.
- Corollary for the hub: **an approving verdict corroborates only within the scope it declares
  having swept.** `Findings: None` with no declared scope does not contradict a REFUTED that was
  verified by string — a smaller reading does not out-vote a larger one, it sits underneath it.
  Splits of this shape are not ties and are not adjudicated as ties.

### A sweep is only as wide as its pattern — reconcile the extraction against the population
`status: ratified` · `provenance: 2026-07-28 · CHIP-VINC-NEUTRO round 5 field finding, verified by string at the hub`

The section above orders the sweep. This one is about the sweep lying to the person who ran it.

A sweep proves a class closed by extracting every member and checking each one. If the extraction
pattern cannot MATCH part of the population, the members it drops are not reported as unchecked —
they are not reported at all, and the sweep reads as complete. The instrument's blind spot is
invisible in its own output.

Field evidence, measured at the hub on `3915f33b`: a document titled `EXHAUSTIVE FIXTURE SWEEP`
proved its own exhaustiveness with `grep -oh 'anchor: "[a-z_]*"'`. The population is
`grep -c 'anchor: "'` = **23**; the character class `[a-z_]` cannot match a capital or an accent,
so four sites were invisible — `"SKU idêntico"` (×2), `"Título parcial"`, `"EAN"`. The sweep
reported ONE violation of the class where there were five, and a machine re-sweep of the same file
returned **13 failed / 5 passed** against the document's **two** findings. The extraction was
narrower than the population and nothing in the output said so.

So any sweep offered as evidence that a class is closed binds:

- **Reconcile two counts, and print both.** Population (the loose anchor: `grep -c` on the
  bare marker) against extraction (what the pattern actually yielded). Unequal without a stated
  reason = the sweep reports its own blind spot instead of a verdict. This is one extra `grep`.
- **Must-fail the pattern against a known member**, the same obligation a test carries. A pattern
  that has never been shown to match something is not known to match anything.
- **Prefer the checker the language already has.** `wireFixtures.ts` replaced this whole class by
  making the impossible fixture UNWRITEABLE — a throwing constructor plus a guard that reads the
  Go declaration — rather than detectable by a sixth `grep`. When a sweep has failed twice, the
  next artifact is a mechanism, not a wider regex (third-round rule, below).
- Corollary, and the reason this is not merely a `grep` tip: **a sweep run by the same faculty
  that produced the defect inherits the defect.** Here the narrow reading appeared in three
  layers — the fixtures, the sweep of the fixtures, and the proof that the sweep was exhaustive.
  The count reconciliation is cheap precisely because it does not depend on that faculty.
- **The population is set by the FACT, never by the edit's footprint.** The most common way the
  pattern comes out narrow is that nobody chose one: the sweep silently runs over *the files that
  were already open*. CHIP-VINC-NEUTRO r7 — a merge moved `ImportacaoSection` out of
  `pages/vinculos/`, killing a `listErpImports` mock in TWO test files; the chip deleted the one
  in the file it was editing and the second survived a full round, caught by the hub with
  `grep -rn` over the directory. The cause was a moved component, so the population is *every
  file that mocked the port*, tree-wide — the edit's footprint has no bearing on it. Ask what
  KILLED the thing, and let that answer pick the search root; a sweep whose root is the diff is
  reporting the diff back to you.

### Vacuous green — an instrument that passes for a reason unrelated to the code
`status: ratified` · `provenance: 2026-07-28 · hub executing seat, CHIP-ANCHORS-3 · CHIP-IMPORT-CHAIN field finding #1`

Sibling of stability ≠ discrimination (§3). There the instrument was stable across both worlds;
here it never looked at either. **Exit 0 is not evidence that anything ran.**

Four instances, one afternoon:

1. `go test -run 'TestX'` where `TestX` is the CONTEXT function of a patch hunk header, not the
   added test → `no tests to run`, **PASS**, exit 0.
2. The target file carries `//go:build integration`, so a green `go test ./...` — 153 packages,
   107 `ok` — **never compiled it**. Package-count green is not lane coverage.
3. `-tags integration` without `MPC_TEST_DATABASE_URL` → every DB test `SkipWithoutTarget` →
   `ok` in the summary.
4. The integration lane runs `go test -tags=integration` **without `-v`** and its artifact
   records only `target`/`status`/`run_id`, so a fully skipped run and a fully green run are
   **byte-identical** in `summary.txt` (CHIP-IMPORT-CHAIN field finding #1, independent).
5. `grep -E "^ *(Test Files|Tests)"` over vitest output returns **empty, exit 0** — vitest emits
   ANSI colour escapes BEFORE the leading whitespace, so `^ *` never anchors. The command ran,
   the tests ran, and the *measuring instrument* saw nothing. Hit independently by the hub
   executing seat and by CHIP-VINC-NEUTRO on the same day, 2026-07-28.

Binding on the executing seat:

- **Count, never tail.** Report packages/tests/assertions actually executed — `ok=N`,
  `no test files=N`, `FAIL=N` — not the last lines of the output. A tail can be empty and read
  as clean.
- **A must-fail that does not go red is not a must-fail.** Before believing a green, prove the
  command can go red at all: wrong-name, wrong-tag and skipped-target all produce the same
  green as a correct run.
- **Name the test by grepping `^func Test` in the file**, never from a hunk header — the `@@`
  context line names the PRECEDING function.
- **Skips are a result, not a footnote.** A lane that can skip must report the skip count, or
  its artifact cannot distinguish ran-and-passed from never-ran.
- **An empty filter is a failed measurement, not a clean result.** Before grepping any captured
  output, print its byte and line count; a non-empty file whose filter yields zero lines means
  the PATTERN failed, not the run. On any coloured tool (vitest, vite, npm), strip escapes first
  — `sed 's/\x1b\[[0-9;]*m//g'` — because `^`-anchored patterns cannot see past them.

### A must-fail arm proves only what its mutation ISOLATES — and never catches an over-strict guard
`status: ratified` · `provenance: 2026-07-28 · CHIP-VINC-NEUTRO round 8, both findings chip-side`

Two limits of the must-fail instrument, found in the same round, both by the chip against its own
arms. They bound what a wall of red is allowed to be reported as proving.

- **A mutation that breaks several rules at once is red but UNATTRIBUTABLE.** The chip's first
  attempt at its third arm deleted a whole `if` line and turned FOUR arms red: without the split,
  signals also reached a check meant only for absences and threw the wrong message there. Red, and
  worthless — the arm you are trying to certify is the one you can no longer see. Narrow the
  mutation until exactly ONE arm moves, then report THAT. The chip discarded the wide mutation,
  replaced it with a two-line reorder, and did not report the four-red result as the stronger
  evidence — which is the behaviour this rule ratifies, not merely the conclusion.
- **Must-fail arms are blind to over-strictness by construction.** A guard that rejects too much
  makes every must-fail arm GREENER, never redder; no arm in the wall can catch it. The only
  instrument that can is a MUST-PASS: a demonstrably PRODUCIBLE input asserted NOT to throw. A
  must-fail wall without one certifies a guard that could be rejecting every real candidate.

Corollary the same round supplied twice: a guard wider than the fact is the same defect as a guard
narrower than it, pointing the other way. The chip's round-6 rule `anchor appears twice` would have
rejected real candidates, because `TitleMatch` seeds `title` FOR and the hard-negative block appends
`title` AGAINST to the same candidate (`generation_service.go:551` + `:560-562`), and the finalizer
dedupes only absences (`:657-661`). Verified at the hub by string before ratification.

### PRODUCIBLE is closed over the generator's SITES; REACHABLE is closed over the DECLARATIONS that exist
`status: ratified` · `provenance: 2026-07-28 · CHIP-VINC-NEUTRO round 8, discriminator supplied by the chip against its own seven-round argument`

Two different questions, and a guard that answers one must not be written up as answering the other.

- **Producibility** asks whether SOME site in the generator emits this shape. It is closed over the
  code paths, and a fixture-vs-generator mechanism can decide it.
- **Reachability** asks whether the shape occurs given the DECLARATIONS the tree actually carries —
  configured providers, capability tables, seeded rows. It is closed over data, and no amount of
  reading the generator decides it.

CHIP-VINC-NEUTRO defended defect V-1 for seven rounds: a queue row whose motivos are ALL
`INCOMPARABLE`, rendering an overflow button with zero chips. Producible — the generator has sites
that emit it. Unreachable — `mercado_livre` is the tree's only capability declaration and does not
supply `marca`, so every candidate carries a `marca UNAVAILABLE`, which the old expression DID
enumerate. The empty cell needs a provider supplying all four anchors; none exists. Seven rounds of
dual gates, all reading the generator, could not see it. One live drive did, immediately.

The sharpest part is where the falsehood sat. The chip's own fixture already carried the words
*"No declaration emits it … unreachable"*. The instrument knew. The prose beside it — where the
VALUE CLAIM lives — still said the screen was broken. **A mechanism being right does not make the
text that sells it right**, and only the text reaches the operator. Restate the value claim in the
terms the mechanism actually decided, or the mechanism is laundering a claim it never tested.

Binding consequences:

- A must-fail arm running on an UNREACHABLE fixture proves the guard DISCRIMINATES. It does not
  prove the live screen was broken. Wording it as "exactly the screen the contract forbids" lends
  the unreachable artifact the authority of the reachable one.
- A gate criterion asserting a screen state must be discharged against DATA, not against the
  generator. The hub owes the chip that seam; withholding it is what let this survive.
- When the true criterion is found, the false premise is REMOVED from the gate artifact, not
  annotated beside it (R-25) — and the replacement must still fail the old code on the same
  observable, or the criterion was not corrected but abandoned.

### The gate's artifacts are the orchestrator's to persist, and the pack must be IN GIT
`status: ratified` · `provenance: 2026-07-28 · CHIP-ANCHORS-3 round 3, findings 6 and 8`

Two failures of custody, same round, same pack. Neither is about review quality — both are about the
verdict surviving long enough to be read.

**Persistence is a step of the ORCHESTRATOR, never delegated to the seat.** A brief that tells the
seat to write its own verdict to a path persists nothing, and it fails differently on each side:
the cold Opus seat has **no Write tool by construction** (§11 third seat — that is the same property
that makes it a reading seat), and the Sol seat's sandbox refused the write outright —
`patch rejected: writing is blocked by read-only sandbox`. So the brief asks for something one seat
cannot do and the other is forbidden to do, and the failure is silent in both cases: the verdict
arrives in the completion notification and nowhere else.

- The orchestrator writes the artifact **in the same act in which the verdict arrives**, before any
  analysis. Verbatim paste first; reading second.
- A refused `apply_patch` is **recoverable**: its payload is in the rollout, and stripping the `+`
  prefix yields the seat's literal intended file. State that provenance in the artifact head, plus
  the seat's tool-set. `208 lines, zero unprefixed lines` is a checkable claim; "transcribed" is not.
- The residual risk is OMISSION, not fabrication, and it is unfalsifiable from the artifact. That is
  why the paste comes first.

**"It arrived in the notification" is not custody, and the hub broke this rule itself**
(2026-07-28, CHIP-VINC-NEUTRO round 7). A seat return is delivered in two shapes depending on its
SIZE: a short verdict is inlined in the completion notification and survives in the transcript; a
long one is written to `tasks/<id>.output` and the notification carries only a summary. Measured
after the fact: **four seat output files, all `0 bytes`**, and no block of the seat's text anywhere
in the transcript. The Sol verdict survived at 5.625 chars because it was inlined; the Opus verdict
of the same round is **gone** — the only trace is the hub's own paraphrase, which is not verbatim
and does not become one by being the last copy standing. The window that killed it was between
ARRIVING and BEING PASTED, and it was crossed by a context compaction, so nothing signalled the
loss. Corollaries, both binding:

- The bigger and more valuable the verdict, the more likely it is the one that vanishes. Paste on
  the turn it lands — not on the turn the ruling is written.
- A lost verdict is declared in the pack, with what was measured to establish the loss, and the
  findings that existed only in the paraphrase **do not become orders**. Re-deriving from memory
  what a seat "would have said" is fabrication with extra steps.
- Recovering a verdict by size (`grep` for the largest matching block) picks up the WRONG round:
  the hub did exactly this and pasted a different chip's gate — `CHIP-ANCHORS-2`'s `refforn` /
  `RouteClassMux` verdict — under this round's heading before checking the header. Match on the
  seat's task id, never on size.

**`git status --porcelain .mnfs/<pack>` clean belongs in the gate brief.** CHIP-ANCHORS-3 ran six
commits with `EVIDENCE.md`, `dispatches/` and six `p6-*.patch` files **untracked** — `git check-ignore`
exits 1, so no rule ignored them; it was omission. The seat reads the pack **from disk** and cannot
tell the difference, so nothing in the review surfaces it. For six commits the two verdicts that cost
the most to obtain existed only in the filesystem of a **disposable worktree** — and the hub destroys
those routinely on close.

The hub checks the same thing before it destroys anything: a worktree teardown after an untracked
pack is a silent, permanent loss of the gate record, and the merge is not proof (a merge carries
tracked files only).

**But `--porcelain` alone is the wrong instrument, and it is the hub's to run, not the chip's.**
Two corrections found the same day the rule was ratified:

- **It false-positives on this platform.** With `core.autocrlf=true`, `git status --porcelain`
  reports ` M` for a file whose content is byte-identical — the stat cache saw a touched mtime and
  reports possibly-modified without re-hashing. Measured on CHIP-ANCHORS-3's worktree at
  `13a09177`: `generation_service.go` shows ` M`, while `git diff --quiet` exits 0, the blob hash is
  the same on both sides (`45438316`), and the byte counts match (41615 = 41615). A freeze-violation
  alarm from ` M` alone is unfounded. Worse than the false alarm is the habituation: once ` M` is
  routine noise, a real modification hides inside it. **Ask the two questions with two instruments** —
  `git diff --quiet HEAD` for content drift, `git status --porcelain --untracked-files=all` read for
  `??` lines only, for unversioned files. `--porcelain` conflates them.
- **A pack cannot assert its own tracked-ness.** The claim is self-referential by construction: the
  file carrying the assertion must be committed *after* the sentence is written, so the state the
  sentence describes is not the state that ends up in the commit. The chip may report what it saw;
  only the hub's executing seat, reading the tip from outside the worktree, discharges it.
- **`git ls-tree` without `-r` counts tree entries, not files.** On CHIP-VINC-NEUTRO's pack the
  non-recursive form returns `4` (two blobs and a subtree at the root) where `-r` returns `40`. `4`
  is not a wrong file count — it is a right count of something else, which is why a seat would fill
  the mandatory section with it and be telling the truth. The custody clause that exists *because* a
  pack went untracked shipped an instrument that could not tell 40 files from 4.
- **A count without its tip rots.** The hub published `38` for this pack; the tree says `38` at
  `ea856c32`, `39` at `7b5c18eb`, `40` at the dispatch tip. The measurement was right where it was
  taken and became false the moment it was quoted against a later tip. Same shape as `base_sha` is a
  FLOOR: every count carries the SHA it was measured at, or it is unfalsifiable.

### A scripted edit on an authority artifact must show its `git diff` before the commit
`status: ratified` · `provenance: 2026-07-28 · CHIP-VINC-NEUTRO, round 6 dispatch`

Re-filing two ledger rows was scripted on the anchor `line.startswith("| 9 |")`. `EVIDENCE.md`
carries **two** tables numbered that way, so the script rewrote both and **deleted** rows 9 and 10 of
the tsc error inventory:

```
- | 9  | ProdutoPage.partialFailure.test.tsx(40,45) | TS2322 same |
- | 10 | ProdutoPage.partialFailure.test.tsx(41,46) | TS2322 same |
```

The script reported success. What exposed it was `grep -c "DISPATCHED"` returning `4` for two
rows — a count disagreeing with the fact it was checking — and only then did `git diff` name the
victims.

This is the vacuous-green family with the damage inverted: **"edited correctly" and "edited
correctly AND destroyed something else" produce the same exit 0.** And it is worse than a
miscount, because a wrong count leaves residue to reconcile against, while a silent deletion leaves
none — the deleted rows belong to no population any sweep anchors, so nothing downstream could
miss them. R-25: honest-unknown is for a gap; deleting is falsity.

- Every scripted or programmatic edit of an authority artifact (pack, ledger, doctrine, contract)
  runs `git diff` and the diff is **read** before the commit. The diff's file and line counts are
  the receipt; "the script said OK" is not.
- **Prove the anchor is unique before using it.** `grep -c '<anchor>'` over the target; more than one
  hit stops the edit. An anchor chosen from the row you are looking at has never been tested against
  the rows you are not.
- The check belongs at edit time, not sweep time. There is no later instrument that finds this.

### Third-round rule — a third defect of the same shape stops the patching
`status: ratified` · `provenance: 2026-07-25 · operator ruling, D-121 · field evidence CHIP-M05 (6 dual-gate rounds)`

**Trigger (either one fires it):**
1. A gate round finds a defect of the **same shape** as a defect already corrected in two
   earlier rounds of the same chip — same failure mechanism, even if a different file,
   layer or artifact; or
2. a chip reaches its **third correction round** on any one criterion.

Shape, not location. CHIP-M05's three instances were a code comment, a scope argument in
prose, and an OpenAPI description — three artifacts, one mechanism: *a claim about
observable behaviour written from the layer that NAMES the value instead of the layer that
PRODUCES it.* Each was point-fixed, each time correctly, and nobody swept for the class
until the sixth round. Two of the three had contradicting evidence already sitting in the
repository.

**What the chip must do instead of emitting a fourth point-fix:**

- **Name the shape in one sentence** — the mechanism, not the instance. If the sentence
  cannot be written, the root cause is not understood yet and the fix is premature.
- **Sweep the whole class, exhaustively and by tool.** Enumerate every site in the write-set
  where that mechanism could occur, `grep`-anchored, with `file:line` for each, and record
  the full list in the evidence pack — including the sites found CLEAN. An unbounded
  "I looked and it's fine" is not a sweep; the enumeration is the artifact.
- **Classify each site** as correct / defective / not-applicable-because-<reason>. Fix every
  defective one in the same round.
- **Prove the class closed, not the instance.** Must-fail on a representative site, or an
  enumeration a reviewer can re-run.
- **Ask whether the remedy is structural or informational** — a guard, a funnel, a type, a
  lint rule, or just a corrected sentence. **Before adding any abstraction, get an
  independent adversarial judgement briefed to argue AGAINST it** (YAGNI/DRY rule-of-three,
  §11 standard). In CHIP-M05 that judgement killed the proposed domain constructor: the
  database CHECK was already the single authority, and a Go copy of the same policy would
  have been a second source of truth free to drift. The real defect was a guard nobody
  called, plus missing funnel discipline, plus missing coverage — not a missing type.
- **Record it as a `ROUND-N FULL ANALYSIS` section** in the evidence pack, and name it in the
  `CLOSED` event.

**Binding on the reviewers and the hub:** a third-round point-fix arriving WITHOUT the class
sweep is not acceptable evidence — the gate returns it, and the hub does not merge on it. The
cost of the sweep is one round; the cost of skipping it was three extra rounds and a wrong
description shipped into a published contract.

**Why this is not just "review harder":** rounds 1-3 of CHIP-M05 each PASSED their own gate on
the fix in front of them. Every individual verdict was correct. The failure is that a gate
scoped to the reported defect cannot see a class, and neither can a chip that keeps answering
the question it was asked. Only an explicit trip-wire on repetition catches it — this is the
local-maximum guard (§11 G1-G3) applied to the correction loop itself.

### A finding BLOCKS only when it names a wrong observable; everything else is a REPORT
`status: ratified` · `provenance: 2026-07-28 · operator ruling ("não nos deixar levar pelo que já temos, mas o que deve ser"), D-122 · field evidence CHIP-ANCHORS-3 r4 + CHIP-VINC-NEUTRO r6`

The gate drifted off the code. CHIP-ANCHORS-3 round 4 returned **both seats REFUTED with three
blocking findings, over a diff the reading seat itself measured as "40 insertions, 12 deletions,
all `//` text; no behavior changed"** (`p6-sol-gate-r4.md`, `a4709d43`). Zero of the three
findings touched behaviour. Meanwhile the wave-2 packs had grown to 9,730 / 11,502 / 21,788 lines
of evidence against 605 / 1,811 / 1,318 lines of code. The gate had stopped reviewing the code and
started reviewing its own paperwork — and a gate that reviews paperwork produces more paperwork,
which is the loop, not the exit.

**The rule.** A finding is `BLOCKING` only if it names a wrong **observable**: behaviour,
security, data, or a published contract. Everything else — prose, counts, metadata, citation
drift, formatting, pack hygiene, self-consistency of the evidence document — files as `REPORT`:
recorded in the pack, corrected in the same round, **does not hold the merge and does not open a
new round**.

**The discriminator is one question, and it is answered before the class is assigned:** *leave
this finding exactly as it is, ship the code — what does a user, an operator, a caller, or a
stored row do differently?* No answer is a REPORT. This is a question about the world, not about
the document, which is precisely why a document-shaped finding cannot pass it.

**Reclassification is the HUB's, never the seat's.** A seat reports everything it finds, at the
severity it believes; suppressing or softening a finding at the seat is the falsity R-25 already
bans. The hub assigns the class, in writing, with the string-level verification it used — the
same verification a merge would need. Reclassifying without re-reading the code is the same
defect one layer up.

**The escape hatch is trip-wired.** Three REPORTs of the same shape fire the third-round rule
above: the shape gets named, the class gets swept, and the class becomes blocking as a class. A
pattern cannot be shipped one waived instance at a time.

**What this preserves.** The instrument works — it is the aim that drifted. The same dual gate
caught `seller_sku` reading `refforn` in ANCHORS-2 round 1, and the executing seat caught a
substring-blind guard by RUNNING a mutation: both executable, both binary, both cheap.
CHIP-VINC-NEUTRO round 6 returned both seats REFUTED with three findings that were all code, all
in the chip's own write-set, and all confirmed by the hub by string. That is the gate on target.
Rigor stays on behaviour, mutation, and execution; it leaves prose-about-prose.

### The reading seats are given the DIFF, not the pack
`status: ratified` · `provenance: 2026-07-28 · D-122 · field evidence CHIP-ANCHORS-3 r1-r4 (4 rounds, 0 AGREEMENT, last one comment-only)`

The subsection above asks the seats for discipline. This one removes the need for it. A rule that
depends on a reviewer choosing not to report what is in front of them is informational; the
structural form is to stop putting it in front of them. Four rounds were spent amending the rules
of the gate while feeding it the same input.

**Reading-seat input, exhaustive:**
1. the code diff **against the merge target's CURRENT tip** — `git diff main <chip-tip> -- <code paths>`,
   `.mnfs/` excluded. Not the chip's dispatch base: a diff against a stale base cannot show a
   REVERT, and reverting a shipped feature is the most expensive thing a merge can do. (This is
   the same two-questions-two-instruments split as pack custody: the governance DRIFT gate keeps
   using the milestone's accepted 40-hex base — `base_sha` is a floor — while the MERGE decision
   reads current-tip. Different questions.) Measured on CHIP-ANCHORS-3 `a02be7f2`: against its
   dispatch base the delta reads `+629/-43` across 11 files and looks clean; against current
   `main` the same branch DELETES 10 files and 462 lines of the `/importacoes` feature that
   CHIP-IMPORT-CHAIN merged at `45b887b3`, plus `cmd/mlprobe`. Four gate rounds never saw it,
   because none of them were looking at that diff;
2. the milestone/chip validation contract, criteria verbatim;
3. the executing seat's RAW lane outputs (ladder, must-fail, fixtures) — measurements, not prose
   about measurements.

**The evidence pack is not reviewer input.** It is hub custody: it records what happened, for the
hub and for whoever reads this in six months. The hub reads it once, for acceptance. Two seats
reading 10-20k lines of it at review depth is how a 605-line change bought four rounds.

A seat MAY request a **named, bounded** pack section when a CODE finding turns on the chip's
stated rationale ("§4 of the pack, the derivation for `:379`"). Named and bounded — never the
tree, never "read the pack".

**What this buys:** a seat cannot file a pack-prose blocking finding, because it cannot see the
pack. The BLOCKING/REPORT split stops being a judgement call and becomes a property of the setup.

**The cost, stated honestly:** a seat reading only the diff will sometimes re-derive something the
pack already settled. That is accepted and it is the right trade — re-deriving from the code is
cheap and is exactly the work we want done twice, while re-reading prose is the work that was
expensive. It also means the pack must stop being where an argument LIVES: if a derivation matters
to the code, it belongs in the code as a comment, where the next reader is standing.

### A total guarantee in a docstring is a claim about EVERY input
`status: ratified` · `provenance: 2026-07-28 · D-122 · field evidence CHIP-VINC-NEUTRO r6 F1/F2 (both seats, blind to each other)`

R-24 ("a claim is TOTAL or it is a report") was ratified against pack prose. It fires in CODE too,
and the code instance is worse, because a docstring compiles, passes the lane, and is what the
next author greps.

`wireFixtures.ts:202` reads `/** A candidate the backend can actually emit. Throws if it is not
one. */` — a total guarantee. The implementation is partial in two independent axes:
`assertProducibleScore:153-170` validates the score triple but imposes SHAPE for exactly one
status (`NO_CANDIDATE`, `:164-169`), and `assertProducibleReasons:110-151` gates its capability
check behind `if (providerCode === "mercado_livre")` at `:142`. Nine fixtures were written
trusting the sentence; three were wrong.

**The tell, and it is general: both seats found this by reading the implementation, neither by
reading the docstring. A false total guarantee does not fail the reader — it fails the believer.**
No amount of review attention on the sentence finds it, because the sentence is not where the
defect is.

**Binding.** A docstring stating a total guarantee (`always`, `never`, `every`, `throws if it is
not`) must be true for every input the function accepts, or state its scope in the same sentence.
When the guard is table-driven, the guarantee points at the TABLE and the table carries the full
tuple — then a partial table is visible as a hole instead of hiding behind a total sentence.
Verify a guarantee by finding the guard's most restrictive gate (`if (x === "literal")`, a shape
check on one branch, a `find` over a partial tuple), never by reading the sentence.

The structural remedy is the same one the chip proposed and the hub ratified: widen the table,
delete the exemptions. `PRODUCIBLE_SCORES` carrying `(state, match_input)` per producer site
dissolves BOTH the `NO_CANDIDATE` special case and the `mercado_livre` gate — one mechanism,
zero hand-written exceptions, and the sentence becomes true.

## 12. Implementer dispatch bindings (core §4 instantiation)
`status: ratified` · `provenance: 2026-07-16 · operator-ratified alt-D+ design session (4 adversarial Opus×Sol rounds + MIS-003 field evidence: hub + CHIP-SAT/CHIP-M02/CHIP-M03)`

- **Prompt-pack:** every codex implementer dispatch carries the canonical implementer
  prompt-pack from the plugin's `HARNESS-CORE.md` §4 ("Implementer prompt-pack — canonical,
  versioned"), stamped `impl-pack v<semver> · milestone <id> · body-sha256 <hash>`. Block
  order fixed: pack → role/repo bindings → slice card (variable, always last); mid-milestone
  changes = dated addenda appended after the card, folded at milestone boundary. The ledger
  records `pack_version` + body hash per dispatch.
- **Deterministic lane tooling:** vendored at `scripts/harness/dispatch/`
  (`New-DispatchPrompt.ps1` assembler · `Invoke-CodexDispatch.ps1` dispatcher ·
  `Test-HarnessPreflight.ps1` · `HarnessDispatch.psm1` · `roles.psd1` · Pester tests),
  hash-locked by `scripts.lock.json` — every entrypoint fail-closes on hash mismatch
  (LOCK-MISMATCH = re-vendor from the harness plugin, NEVER re-lock locally;
  `Update-ScriptLock.ps1` refuses to run outside plugin source). Scripts validate FORM only
  (max verdict `FORM-COMPLETE`); merit stays with the review ladder.
- **Dispatch registry:** append-only JSONL, lifecycle `assembled → started →
  completed/failed/cancelled`; assembly alone never yields a complete row — only the
  dispatcher promotes to `started`. Location: `output/harness/dispatch-registry.jsonl`
  (`assumed` — hub may re-bind to a `.mnfs/` path at first field use; whichever path the hub
  uses first becomes the binding). Dispatcher writes atomic `agent__<id>.result.json`
  receipts (dispatcherPid, timestamps, exit code, prompt/log/output hashes) replacing the
  bare `.done` sentinel, and NEVER auto-re-issues a writing worker (timeout/exception →
  `failed` + stop for diagnosis).
- **Role → flags:** `roles.psd1` mirrors core §1 (canonical there; edits land in CORE first).
  Per-role default sandbox: writers `workspace-write`, judges/auditors `read-only`;
  `--effort` always explicit.
- **F-A commit-denial clause (index.lock):** if a worker's `git commit` is denied by an
  existing `.git/index.lock` (or sandbox git-write denial), the worker must ATTEMPT the
  commit once, and on denial LEAVE THE FILES IN PLACE and REPORT the denial verbatim in its
  final report (evidence type `could-not-run`) — never delete work, never retry-loop, never
  remove the lock file itself. The chip/hub owns lock diagnosis (CHIP-M03 field finding,
  MIS-003 W1). `Test-HarnessPreflight.ps1` surfaces `git-index-lock` as an advisory check
  before writer dispatch.

### The hub ANNOUNCES the cut — a freeze the chip cannot observe is not a freeze

The gate's input is frozen the moment the hub generates the diff, and until that moment the chip
is working normally and correctly. A chip cannot honour a freeze whose start it was never told.
Every round where a chip "violated the freeze" by committing after dispatch is first a round
where the hub cut the patch in silence.

- The hub sends the chip a CUT message naming the tip the patch was generated from, in the same
  act as generating it. From that tip the chip writes nothing until the verdict returns.
- A chip commit that arrives after the cut does not automatically void the round. The hub
  measures the delta and applies one discriminator: **does it touch production code, or anything
  a returned verdict rests on?** No to both → the round stands, and the delta is recorded as
  HUB-verified rather than seat-verified, stated in exactly those words. Yes to either → re-cut
  and re-dispatch.
- The residual is named, not hidden: for a standing round, some lines were verified by one
  reader instead of two. Burning two fresh seats over a dead-mock deletion is the point-fix the
  third-round rule exists to prevent; pretending the patch was current is the falsity R-25 bans.

Field evidence (CHIP-VINC-NEUTRO round 7): the hub cut a 2,662-line patch at `293c1485` and told
nobody. The chip then landed `2b956e19`, deleting a genuinely dead mock and the comment that had
predicted its own falsity — a correction the cold seat independently filed as a REPORT on the
same locus. Delta measured: one test file, 18/20, zero production lines. Round stood. The chip
reported the collision itself and refused to decide the hub's question for it, which is the
behaviour the rule is written to preserve.

### An automated gate must name the tree it measured, or its verdict is unattributable

A hook, script or lane that reports on "the repo" without naming WHICH checkout it read produces
a verdict that cannot be believed OR disbelieved. In a worktree-per-chip topology the hub session's
shell cwd is not the hub's checkout: the two differ by an arbitrary amount of history, and a gate
that resolves paths relative to cwd is measuring a tree nobody is working in.

- Any automated verdict about repository state PRINTS the absolute path it measured and the tip
  SHA it measured at, in the same message as the verdict. A verdict without both is `unknown`,
  not `fail` — and specifically it does not license a human-visible accusation.
- The reader's obligation on receiving such a verdict is to re-measure in the NAMED checkout
  before acting. Two instruments, not one: the claim and the tree it was made about.
- Corollary for the hub: a `cd` in a hub command is not optional hygiene. Every hub command
  states its checkout in the same invocation, because PowerShell and the Bash tool both reset
  cwd between calls and the reset destination is the stale worktree.

Field evidence, twice in one session: the Stop hook fired
`CLOSED claimed but no evidence pack exists in this worktree (.mnfs/**/_chip-*/EVIDENCE.md)`
while `C:\…\marketplace-central\.mnfs\` carried 15 `_chip-*/EVIDENCE.md` files and each in-flight
chip carried its own on its branch. Both halves of the accusation were false: no `CLOSED` had been
sent (the token appeared inside ORDERS to chips), and the packs existed. The hook had resolved
`.mnfs/**` against `.claude/worktrees/epic-lehmann-4ffbad`, whose `.mnfs/` holds MIS-001..004 and
no MIS-006 pack at all. The cost is not the wasted check — it is that an alarm which is wrong
twice trains its reader to skip the third.

### Contingency lane — codex quota outage (TEMPORARY, ratified 2026-07-18, D-23)
`status: ratified` · `expiry: codex quota return (2026-07-25) or MIS-004 demo close, whichever first`

Machine-wide codex quota exhausted 2026-07-18 until 2026-07-25 (post-demo; evidence
CHIP-M04 F-ENV-M04-1). All GPT roles rebind to Claude for MIS-004 pre-demo work:
- P2 planner: **cold Opus subagent** (fresh context, plan-only prompt; sonnet fallback if
  Opus limits bite). Plan outputs unchanged (core §4.1 — write-sets, contract-satisfiability,
  verification map).
- Implementers: **sonnet subagents** (already sanctioned core §1 fallback, field-proven
  M-01 s8–s12). Anti-slop contract + slice cards unchanged.
- Slice reviewer: independent Claude (sonnet), implementer ≠ reviewer unchanged.
- P6 dual gate: **cold Opus + INDEPENDENT second Claude reviewer (sonnet, separate
  dispatch, adversarial-refute prompt)** replacing the GPT side. Cross-vendor diversity is
  lost — compensate: the second gate prompt MUST be refutation-framed, and the hub
  spot-check right stays. Agreement/reconciliation rules unchanged.
- Investigators/bulk reads: sonnet (haiku for trivial greps).
On expiry the matrix reverts to core §1; milestones closed under this lane MAY receive a
retroactive GPT-5.6 Sol medium review at mission closeout (operator's call).
- **Evidence types bind here too (core §5):** worker reports classify each verification as
  `ran` / `assumed` / `could-not-run`; a Pass is only recordable on `ran` with artifact.

## Amendment log

```
2026-07-15 · all sections · ratified · profile extracted from combined docs/HARNESS.md (operator-ratified doctrine) at A′ restructure; swap to CORE+profile binding scheduled for M-01 close
2026-07-15 · §2 · ratified · governance lane: run from clean detached worktree + full 40-hex BaseSha (field finding: main-checkout sweep of .claude/worktrees false-fails; short sha = GOV_SEMANTIC_DRIFT base-sha-invalid) — memory/governance-lane-clean-worktree.md
2026-07-15 · §3 · ratified · fresh-worktree GOMODCACHE warm + HPG_MIGRATION_FAILED/migrations_first=-1 false-alarm signature (M-01 field finding)
2026-07-15 · §4 · ratified · session postgres container harness:pg:up/down (mpc-pg-session-<8hex>) + createdb first-boot retry absorbing 3D000
2026-07-16 · §11 + rubrics · ratified · REVIEW-STANDARD execution model (operator study, method-level): prompt-pack vehicle, one-reviewer slices, simultaneous dual gate, disjoint-slice overlap cadence, ★ crew noise control (FAIL-restraint + advisory cap 10 + learnings) — plugin rubrics amended upstream
2026-07-16 · §11 · ratified · reviewer model/reasoning bindings pinned in standard §13 (operator correction — execution must carry the ratified matrix, evidence-based): sonnet slice reviewer, Opus+Sol-medium dual gate, sonnet ★ crews + Sol-high P7, codex-dispatch resolver mandatory for GPT flags
2026-07-16 · §11 · ratified · dual-gate Claude side = COLD Opus SUBAGENT (clean context, model=opus explícito), never the orchestrating session (operator ruling — self-grade anchoring bias; quota ruling forbids fable subagents, not opus); every review context cold, bounded inputs only
2026-07-16 · §11 (new) + HARNESS.md §4 · ratified · REVIEW-STANDARD adopted (operator-requested hardening): fixed review order, global-vs-local-maximum G1-G3, YAGNI+DRY rule-of-three, two-axis severity, anchor-or-abstain+receipts, deterministic pre-pass (go vet added to L0), dual-gate agreement merge, delta re-review, learnings memory docs/REVIEW-LEARNINGS.md, ≤300-line slices
2026-07-15 · §7 · ratified · no-stub doctrine: validation contracts/tests never fallback stub/mock on integration seams without operator authorization; real dependency required; planning declares real bindings up front (operator ruling after M-01 C10 — root.go nil-DB cost reader + permanent-unavailable policy reader passed hermetic gates, blocked live validation)
2026-07-15 · §2 L2 + §6 · corrected · dev stack ONLY via docker compose (npm run docker:dev / docker:oauth), hub-owned; chips never boot own server, never bind :8080/:5174, never load .env* into session env (field violation: M-01 chip ran bare worktree server with real .env for C10 — 42P01s were self-inflicted bypass of compose postgres+entrypoint)
2026-07-16 · header · ratified · binding swap executed at M-01 boundary: CORE (vendored docs/HARNESS-CORE.md @ mnfs-harness cd114e6) + this profile + mission Parallel Execution Plan; docs/HARNESS.md reduced to pointer; harness-control seam extended to the three doctrine files
2026-07-16 · §2 · ratified · GOCACHE absolute-path binding on Windows/pwsh (D-14, M-01 ledger practice applied ad hoc)
2026-07-16 · §2 · ratified · per-milestone governance base anchor: chip drift gate runs against the milestone's accepted 40-hex base SHA from the chip prompt; other-base REDs on long-lived worktrees are not chip defects (M-01 field 2×)
2026-07-16 · §2 · ratified · known pre-existing failure allowlist for L1 (cite, don't re-prove; hub-owned edits; entries carry evidence + backlog owner) — seeded with TestPhase1SmokeFlow (M-01 re-proved it 5×)
2026-07-16 · §6 · ratified · dev-stack sync standing policy: COMMITTED event with stack-sync <sha> → hub rebuilds without negotiation (M-01: 4 round-trips were restart asks)
2026-07-16 · §6 · ratified · live dispatch viewer (hub-served live-server.mjs :7391, dynamic multi-session scratchpad discovery) + DB-specialist consultation seam (MNOS session local_ec787804, hub-relayed via REQUEST db-consult) — operator-requested at W1 dispatch
2026-07-16 · (upstream) · ratified · core amendments landed in mnfs-harness cd114e6: sonnet fallback implementer row, COMMITTED event grammar, lean close (★ crew superseded at close by P6 dual gate), additive contract-lock named mechanism, P2 required plan outputs (write-DAG + contract satisfiability + lock pre-identification), Claude-side dispatch visibility accepted limitation
2026-07-16 · §3 · ratified · test-fixture conventions from CHIP-SAT W1 close: µs truncation on timestamp round-trip fixtures (F-4), migration-count fixture bump per migration grant, modules.json entry per new module (GOV_MODULE_COVERAGE); F-2 orders positional flake fixed on main d9a36a0a (set containment)
2026-07-16 · §3 · ratified · node bootstrap clause for fresh worktrees: npm ci at worktree root = env prep (mirror of gomodcache clause); never symlink-reuse another checkout's node_modules (CHIP-SAT field finding + REQUEST)
2026-07-16 · §10 · ratified · mnfs-workflow execution-layer skills denylisted (deleted at source in mnfs-harness 6b29412 layered unification; stale codex cache 0.1.0 + ~/.codex/plugins/mnfs-codex-plugin still ship them — operator field finding: CHIP-SAT worker auto-loaded feature-execution). General rule: auto-discovered skills never bind; only prompt-pack pins are doctrine. Cache repackage to 0.2.0 deferred to W1 close (no tooling swap under running workers).
2026-07-16 · header + §12 (new) · ratified · alt-D+ implementation method adopted (operator-ratified after 4 adversarial Opus×Sol rounds + MIS-003 field evidence): docs/HARNESS-CORE.md + docs/REVIEW-STANDARD.md re-vendored @ mnfs-harness 6206cc1 (implementer prompt-pack v1.0.0 in CORE §4, canonical dispatch-prompt architecture, deterministic lane, evidence types ran/assumed/could-not-run, reproduce+1-fixup→BLOCKED; REVIEW §9 remedy re-review resumes same reviewer, §13 reviewer reads worker prompt-file, §14 slim read-mandate pack); deterministic dispatch tooling vendored scripts/harness/dispatch/ with fail-closed scripts.lock.json (28 Pester green at source); F-A index.lock commit-denial clause (attempt once, leave files, report verbatim — CHIP-M03); harness plugin 0.3.0 synced to local Claude Code cache. Field-test milestone next — chip converses with design session; pack v1.1.0 fed by its retro.
2026-07-17 · §3 · ratified · codex-cli 0.144.4 non-fatal cache-TTL warning (`failed to renew cache TTL: missing field supports_reasoning_summaries`) added to false-alarm signatures — CHIP-M02 field finding, MIS-004 wave A
2026-07-18 · §3 · ratified · full-tree gofmt CRLF false-alarm on Windows worktrees (CHIP-M01 F-ENV-M01): pre-existing files flagged with ^M-only diffs = autocrlf artifact, not defect; scope gofmt to authored dirs
2026-07-17 · §2 · ratified · L1 allowlist maintenance: TestPhase1SmokeFlow RETIRED (CHIP-M02 root cause: fixture "smoke-prod-1" vs positiveProductID integer validator, both pre-existing on base 59d0e62f; hub fixture fix "1001", lane green run 5b244bce) · +TestListingsReadContractEndToEnd intermittent full-lane flake entry (passes isolated 8/8; backlog test-isolation audit, hub queue)
2026-07-17 · §10 · ratified · F-ENV-3: impeccable plugin auto-injects skill mandate into codex worker sessions and derails them (NO_PRODUCT_MD abort, CHIP-M03 field 2×); mandatory skill-discovery denylist clause in every codex dispatch bindings block; core-candidate, upstream at milestone boundary
2026-07-17 · §3 + §2 · ratified · codex workspace-write sandbox cannot run vite/esbuild build on Windows — false-alarm signature + mitigation (a): chip-side build/tsc/vitest re-run is verification of record, workers stay workspace-write (CHIP-M03 field finding; core-candidate, upstream at milestone boundary) · L1 allowlist +TS2688 @types/node pre-existing base break (fix owner CHIP-M03 via grant D-05)
2026-07-18 · §3 · ratified · codex combined build+vet+test single-command timeout = false alarm (cold go build ~65s vs 120s codex cap); split lanes, re-run split before treating as defect (CHIP-M02 field finding F-03-S1; core-candidate for impl-pack template, upstream at milestone boundary)
2026-07-18 · §3 · ratified · CHIP-M03 close findings F-ENV-8/9/10: full-suite L0 at close = backstop for moved-test blind spot (never waive); codex dispatch via 0.145.0-alpha.18 binary (0.144.4 PATH binary dies on cache-TTL renewal); browser screenshot rasterizer broken → visual QA via computed-style + a11y-tree captures
2026-07-18 · §12 (contingency block) · ratified · TEMPORARY Claude-only lane for codex quota outage (machine-wide, resets 2026-07-25 post-demo; CHIP-M04 finding F-ENV-M04-1, ledger D-23): planner=cold Opus subagent, implementers=sonnet, dual gate=cold Opus + independent adversarial sonnet, investigators=sonnet; expiry=quota return or demo close; optional retroactive Sol-medium review at mission closeout
2026-07-18 · (upstream) · ratified · efficiency amendments landed in mnfs-harness 15001de (harness 0.3.1), re-vendored to docs/: REVIEW §8 targeted refuter, REVIEW §9 dual-gate delta re-verdict (D-34 field), CORE §4.1 seam-closure checklist + CORE §2(g) pre-authorized seam grants (M-04 D-29/D-32/D-35 field), CORE §5 selective re-verify PILOT (MIS-004 = pilot mission); operator-ratified 2026-07-18; plugin cache synced 0.3.1
2026-07-18 · (upstream) · ratified · hub support crew landed in mnfs-harness c0fd334 (harness 0.3.2), CORE §1 + harness-hub skill, re-vendored to docs/: fixed cheap-subagent crew hub-ops (ladder/stack/governance/housekeeping) + hub-scribe (files hub-authored ledger/status/commits) + hub-analyst (read-only evidence checks); judgment (rulings, event replies, acceptance, collision calls) never delegated; crew never pushes/merges/authors doctrine; operator-ratified 2026-07-18; plugin cache synced 0.3.2
2026-07-18 · (upstream) · ratified · Opus adversarial-gate remediation landed in mnfs-harness 8fa7ad8 (harness 0.3.3), re-vendored to docs/: crew ships as plugin agent definitions harness/agents/hub-{ops,scribe,analyst}.md (spawn via Agent tool, persistent via SendMessage, mandatory respawn at milestone boundary); CORE §2(g) grant does not waive §3 collision test; CORE §5 PILOT spot-check risk-weighted (integrity-critical/live-integration first, rotating); REVIEW §9 delta re-verdict = declared bounded exception to §8/§13 cold mandate (round-1 gates stay cold); REVIEW §8 refuter may escalate to general breadth on thin first-family pass (must say so); scribe fail-closed branch-guard; analyst base-SHA statement; worker SKILL pin 10 (grants bind as written, P5 follows §5 PILOT). Gate verdict PASS-WITH-CONDITIONS, all 12 findings remediated; plugin cache synced 0.3.3
2026-07-20 · §1 · ratified · default branch corrected `master`→`main` (operator directive D-120: "be on default main"). No `master` ref ever existed (origin/HEAD → origin/main); the stale `master` binding masked a detached-HEAD boot anomaly where nested worktree `hub-erp-main` held `main` while the primary dir sat detached. Consolidated: freed `main` (detached hub-erp-main HEAD), checked out `main` in primary; hub-erp-main dir cleanup deferred (Windows "Function not implemented" on node_modules junction / .gomodcache read-only — `git worktree prune` after stack re-point)
2026-07-25 · §11 (new subsection) · ratified · third-round rule: a third defect of the SAME SHAPE (or a third correction round on one criterion) stops point-fixing and requires a named mechanism + tool-anchored exhaustive class sweep (clean sites listed too) + class-level must-fail + independent anti-abstraction judgement before any new type, filed as ROUND-N FULL ANALYSIS; a third-round point-fix without the sweep is not acceptable evidence and the hub does not merge on it (operator ruling D-121; field evidence CHIP-M05 — same mechanism point-fixed in a code comment, a scope argument and an OpenAPI description across 6 rounds, each round passing its own gate)
2026-07-28 · §3 · ratified · mid-`npm ci` `node_modules` fabricates `tsc` module-resolution errors that impersonate real defects, and the error count is STABLE across repeated runs — reproducibility is not discrimination; the only discriminating observable is the `npm ci` process exit (CHIP-IMPORT-CHAIN field finding; sibling of the 15-error trap, where the expected error COMPOSITION also failed to prove which tree was read)
2026-07-28 · §3 (review-process note) · ratified · sandbox vitest blindness has teeth: a codex worker's "typecheck clean, committed" is a claim about what it could OBSERVE, not about the suite — one dispatch added a `<Link>`, requiring router context in a moved component, and turned 2 unrelated test files red invisibly; the chip's post-dispatch vitest re-run (already the verification of record under the workspace-write clause) caught them plus a genuine timing defect. A chip forwarding a worker's "committed" as green is laundering blindness into evidence (CHIP-IMPORT-CHAIN field finding)
2026-07-28 · §11 (new subsection) · ratified · the dual gate gains a THIRD seat, an independent EXECUTOR owned by the hub: a physically read-only reading seat cannot discharge an execution criterion by construction, so ladder/must-fail/fixture/`git` facts move to a seat that is independent of the IMPLEMENTER (independence from the reviewer is irrelevant — it produces measurements, not verdicts). A reading seat whose blocking finding reads "I could not run X" has found nothing; route it, do not spend a correction round, and scope execution criteria out of reading prompts. Field evidence CHIP-ANCHORS-3 round 1 (Sol side refuted on a denied Go tmpdir and a denied Postgres, zero code findings) + CHIP-IMPORT-CHAIN (both seats returned I9 NOT-PROVEN independently). Exercised BEFORE ratification: `_hub-gate-import-chain/GATE-P6.md` §I9
2026-07-28 · §11 · noted · a gate seat can FABRICATE a file:line citation — the Sol side of the CHIP-IMPORT-CHAIN gate reported `apps/web/vitest.chip.config.ts` as existing on disk and quoted its `:4` and `:16`, for a file present in no tree and tracked nowhere; probable origin is `chip.md:113-114`, which ORDERS the file deleted, read as an observation. Existence claims from a reading seat are verifiable by the hub in seconds and should be, before they enter a pack
2026-07-28 · §3 · ratified · the FE vitest lane's `cd apps/web` is PART OF THE LANE: the same command from the worktree root adds `packages/sdk-runtime`, whose contract tests resolve `../../contracts/openapi.yaml` from `process.cwd()` and go red on a path that does not exist there; from the repo root it globs ~1260 files and returns ~94 failed. Both reds are the instrument. Third instance of stable-but-non-discriminating (hub executing-seat finding, first backend-independent run of the §11 third seat, on CHIP-VINC-NEUTRO @`ebb309ac`)
2026-07-28 · §11 · noted · a physically read-only Opus seat has NO Write, so its task `.output` is 0 bytes BY CONSTRUCTION and the verdict exists only in the completion notification — observed twice (CHIP-ANCHORS-2, CHIP-VINC-NEUTRO). "Transcribed, not captured" with declared provenance is therefore the honest and the ONLY available form for that seat; demanding on-disk capture demands what the instrument does not offer. The residual risk is not fabrication but OMISSION (a finding emitted and not carried over), and it is unfalsifiable from the artifact — mitigation is verbatim paste as the FIRST act after the notification arrives, before any analysis, plus the seat's tool-set named in the header
2026-07-28 · §7 · ratified · no command that dumps a whole environment (`docker inspect`, `docker exec … env`, bare `printenv`) — diagnose by variable NAME, one at a time; binds WORKERS, so it goes in the dispatch-prompt denylist, not only in the profile the chip reads. Holds for throwaway targets: the session Postgres password is CSPRNG-generated per container and dies with it, so there is nothing to rotate — which is precisely why "it was disposable" cannot be a case-by-case exemption (CHIP-ANCHORS-3 R4 worker field finding)
2026-07-28 · (upstream) · accepted · CHIP-IMPORT-CHAIN F-3 (`New-DispatchPrompt.ps1` assembles an unknown role string cleanly, failing only later at `Invoke-CodexDispatch.ps1` with `ROLE-UNKNOWN` — validate against `roles.psd1` at assembly time) and F-4 (vitest in a junction-only worktree needs an absolute `setupFiles` path plus `server.fs.strict: false`; the junction realpath resolves outside the vite root) belong to `mnfs-harness`, not this profile — routed upstream, not filed here
2026-07-28 · §11 (vacuous green) · ratified · fifth instance, and the first where the MEASURING FILTER was itself vacuous: `go test` writes ANSI colour escapes BEFORE the leading whitespace, so a `^ *` anchor matches nothing and the count comes back 0 — reported as "clean". An empty filter is a failed measurement, not a clean result; strip escapes first (`sed 's/\x1b\[[0-9;]*m//g'`) and assert the population is non-empty before reading the filtered count (hub executing-seat finding) — commit `a9dc6caa`
2026-07-28 · §11 (pack custody) · ratified · custody is TWO questions and needs TWO instruments: content drift = `git diff --quiet HEAD` (exit code), untracked = `git status --porcelain --untracked-files=all` read for `??` lines ONLY. `git status --porcelain` alone false-positives with `core.autocrlf=true` — the stat cache reports ` M` for a byte-identical file (blob equal both sides, 41615 = 41615 bytes), and the hub nearly reported a freeze violation from it. Corollary: a pack cannot assert its own tracked-ness — the file carrying the sentence is committed after the sentence exists — commit `9f8a6ec1`
2026-07-28 · §11 (new subsection) · ratified · a scripted edit on an authority artifact must show its `git diff` before the commit: "edited correctly" and "edited correctly AND destroyed something else" share exit 0. Same row: `git ls-tree` WITHOUT `-r` counts tree entries, not files (4 vs 40 on the same pack) — a right count of the wrong thing, which is why a reader reports it in good faith; and a count without its tip SHA rots the moment it is quoted (the same pack measured 38/39/40 files at three different SHAs, none of them a miscount) — commit `25716bdb`
2026-07-28 · §11 + `scripts/harness/pack-measure.sh` (new) · ratified · the self-reference class gets a mechanical corrector: numbers about a pack are GENERATED from OUTSIDE at a named SHA, never typed into the file they describe (four defects in one day came from a document measuring itself — the fixed point does not exist). Every figure carries its SHA and its UNIT, because bytes ≠ characters (a 598-figure "discrepancy" between two seats was neither side's arithmetic: 20737 bytes vs 20536 chars, neither naming the unit) and LF-terminated lines ≠ `split('\n')` parts. The corrector shipped the defect it exists to kill and was caught by validation against a known answer: `wc -m` without a UTF-8 locale silently returns BYTES, so both columns agreed on accented prose; fixed with `tr -d '\200-\277'` — commit `c66ea7c7`
2026-07-28 · §11 (new subsection) · ratified · BLOCKING is reserved for a wrong OBSERVABLE (behaviour, security, data, published contract); prose, counts, metadata, citation drift and pack hygiene file as REPORT — corrected in the same round, never holding the merge or opening a new one. Discriminator: ship the code with the finding untouched, and name what a user, operator, caller or stored row does differently; no answer = REPORT. Reclassification is the HUB's with string-level verification, never the seat's (a seat still reports everything — softening at the seat is the falsity R-25 bans), and three REPORTs of one shape fire the third-round rule. Field evidence: CHIP-ANCHORS-3 r4 returned both seats REFUTED with 3 blocking findings over a diff its own reading seat measured as "40 insertions, 12 deletions, all `//` text; no behavior changed", while the wave-2 packs stood at 9,730 / 11,502 / 21,788 evidence lines against 605 / 1,811 / 1,318 code lines — the gate had begun reviewing its own paperwork (operator ruling: "não nos deixar levar pelo que já temos, mas o que deve ser")
2026-07-28 · §11 (new subsection) · ratified · STRUCTURAL half of the row above: the reading seats are given the DIFF, not the pack. Input is exhaustively (1) `<base>..<tip>` with `.mnfs/` excluded, (2) the validation contract verbatim, (3) the executing seat's raw lane outputs; a seat may request a NAMED, BOUNDED pack section when a code finding turns on the chip's rationale, never the tree. The pack is hub custody, read once at acceptance — two seats reading it at review depth is how a 605-line change bought four rounds with zero AGREEMENT. A seat then cannot file a pack-prose blocking finding because it cannot see the pack: the BLOCKING/REPORT split stops being a judgement call. Accepted cost, stated: a diff-only seat sometimes re-derives what the pack settled — cheap, and the right work to do twice; corollary is that a derivation which matters to the code belongs IN the code
2026-07-28 · §11 (new subsection) · ratified · R-24 fires in CODE: a docstring stating a total guarantee (`always`/`never`/`every`/`throws if it is not`) is a claim about every input the function accepts, or it states its scope in the same sentence. `wireFixtures.ts:202` promises "Throws if it is not one" while `assertProducibleScore:164-169` imposes shape for exactly ONE status and `assertProducibleReasons:142` gates its capability check behind `if (providerCode === "mercado_livre")` — nine fixtures were written trusting the sentence, three were wrong. Both gate seats found it by reading the IMPLEMENTATION, neither by reading the docstring: a false total guarantee does not fail the reader, it fails the believer, so review attention on the sentence can never find it. Verify by locating the guard's most restrictive gate; when table-driven, the guarantee points at the table and the table carries the FULL tuple, so a partial table is a visible hole instead of a total sentence. Trigger named by the hub on 2026-07-27 ("the day this shows up in CODE it is an amendment on the spot") and satisfied by CHIP-VINC-NEUTRO r6 F1/F2
2026-07-28 · §11 (new subsection) · ratified · a DELTA brief must order a CLASS sweep, or the second occurrence is structurally invisible: a brief that NAMES the site to re-examine teaches the seat to stop at it. CHIP-ANCHORS-3 round 2 — the Sol seat checked the named site, found it corrected, returned `Findings: None` on a file carrying the same false universal 83 lines below, findable by a one-line `grep`; the Opus seat swept on its own initiative and returned REFUTED with it. The seat that returned None did literally what the brief asked, so the defect is the BRIEF's. Binds: give the class as searchable tokens (`never`/`always`/`only`/`unreachable`/`cannot`/`no longer`/`every` for false-totality prose), a verdict with no SWEEP section is incomplete on its face, and the AUTHOR runs the same sweep against their own pack before publishing. Hub corollary: an approving verdict corroborates only within the scope it declares having swept — `Findings: None` with no declared scope does not out-vote a REFUTED verified by string, so splits of this shape are not ties (chip-authored field finding, ratified verbatim)
2026-07-28 · §11 (new subsection) · ratified · VACUOUS GREEN — an instrument that passes for a reason unrelated to the code; sibling of stable-but-non-discriminating (§3), except it never looked at either world. Exit 0 is not evidence that anything ran. Four instances in one afternoon of the hub's executing seat on CHIP-ANCHORS-3: (1) `-run 'TestX'` naming the CONTEXT function of a patch hunk header instead of the added test → `no tests to run`, PASS; (2) the target file carries `//go:build integration`, so a green `go test ./...` (153 packages, 107 `ok`) never compiled it; (3) `-tags integration` without `MPC_TEST_DATABASE_URL` → every DB test skips → `ok`; (4) the integration lane runs without `-v` and records only target/status/run_id, so a fully skipped run and a fully green run are byte-identical in `summary.txt` (CHIP-IMPORT-CHAIN field finding #1, independent). Binds the executing seat: COUNT never tail (`ok=N`, `no test files=N`, `FAIL=N`), prove the command can go red before believing a green, name tests by grepping `^func Test` (the `@@` context line names the PRECEDING function), and report skip counts as a result rather than a footnote
2026-07-28 · §11 (new subsection) · ratified · A SWEEP IS ONLY AS WIDE AS ITS PATTERN — the members an extraction cannot MATCH are not reported as unchecked, they are not reported at all, so the instrument's blind spot is invisible in its own output. Verified by string at the hub on `3915f33b`: a document titled `EXHAUSTIVE FIXTURE SWEEP` proved its exhaustiveness with `grep -oh 'anchor: "[a-z_]*"'` against a population of `grep -c 'anchor: "'` = 23; the character class cannot match a capital or an accent, so `"SKU idêntico"` (×2), `"Título parcial"` and `"EAN"` were invisible — ONE violation reported where there were five, and a machine re-sweep of the same file returned 13 failed / 5 passed against the document's two findings. Binds any sweep offered as class-closure evidence: reconcile population count against extraction count and PRINT BOTH (unequal without a stated reason = the sweep reports its own blind spot, one extra `grep`); must-fail the pattern against a known member; after a sweep has failed twice the next artifact is a MECHANISM, not a wider regex. Corollary: a sweep run by the same faculty that produced the defect inherits the defect — here the narrow reading appeared in three layers (the fixtures, the sweep, and the proof the sweep was exhaustive), and the count reconciliation is cheap precisely because it does not depend on that faculty (CHIP-VINC-NEUTRO round 5 field finding)
2026-07-28 · §11 (new subsection) · ratified · GATE CUSTODY, two failures in one round (CHIP-ANCHORS-3 round 3, findings 6+8). (a) PERSISTENCE IS THE ORCHESTRATOR'S STEP, never delegated to the seat: a brief telling the seat to write its own verdict persists nothing and fails differently per side — the cold Opus seat has no Write BY CONSTRUCTION (same property that makes it a reading seat) and the Sol sandbox refuses outright (`patch rejected: writing is blocked by read-only sandbox`); the orchestrator pastes verbatim in the same act the verdict arrives, before analysis, and a refused `apply_patch` is RECOVERABLE from the rollout (strip the `+` prefix — "208 lines, zero unprefixed lines" is checkable, "transcribed" is not). Residual risk is OMISSION, unfalsifiable from the artifact, which is why the paste comes first. (b) `git status --porcelain .mnfs/<pack>` CLEAN belongs in the gate brief: the chip ran six commits with EVIDENCE.md, dispatches/ and six p6-*.patch UNTRACKED (`git check-ignore` exits 1 — omission, not a rule), and the seat reads the pack FROM DISK so nothing in the review surfaces it; for six commits both gate verdicts existed only inside a DISPOSABLE worktree. The hub checks the same before teardown — a merge is not proof, it carries tracked files only
2026-07-28 · §11 (new subsection) · ratified · AN AUTOMATED GATE MUST NAME THE TREE IT MEASURED. Second occurrence in one session of the Stop hook reporting `CLOSED claimed but no evidence pack exists in this worktree (.mnfs/**/_chip-*/EVIDENCE.md)` against a tree that is not the hub's checkout: the primary carries 15 `_chip-*/EVIDENCE.md` and each in-flight chip carries its own on its branch, while `.claude/worktrees/epic-lehmann-4ffbad/.mnfs/` holds MIS-001..004 and no MIS-006 pack at all. BOTH halves of the accusation were false — no `CLOSED` was sent either (the token appeared inside ORDERS to chips). Binds: an automated verdict about repo state prints the absolute path and tip SHA it measured, in the same message as the verdict; without both it is `unknown`, not `fail`, and does not license a human-visible accusation. Reader re-measures in the NAMED checkout before acting. The cost is not the wasted check but that an alarm wrong twice trains its reader to skip the third — the hub named this trigger on the first occurrence and this is the second
2026-07-28 · §2 · corrected · the profile printed the COPYABLE line `GOCACHE=.gocache` three times (L0, L1, and the §7 pre-pass), the first of them three lines above the rule requiring an ABSOLUTE path. A chip copied it and got an 83-byte EXIT 1 with zero `=== RUN` — the sixth instance of vacuous green, and the only one the doctrine itself handed over. Relative literals removed from all three sites and replaced with a pointer to the rule; the rule now carries both the pwsh and the bash binding, extends to `GOMODCACHE`, and states that this file prints no copyable relative form. A doctrine artifact that hands you the line it forbids is worse than one that says nothing (CHIP-ANCHORS-3 field REPORT on a hub-owned file — the chip measured it instead of working around it in silence)
2026-07-28 · §11 (new subsection) · ratified · THE HUB ANNOUNCES THE CUT. A chip cannot honour a freeze whose start it was never told, so every "the chip wrote during the gate" round is first a round where the hub generated the diff in silence. Binds: the hub sends a CUT message naming the tip the patch came from, in the same act as generating it; a later chip commit does not automatically void the round — the hub measures the delta and asks whether it touches production code or anything a returned verdict rests on, standing the round with the delta recorded as HUB-verified rather than seat-verified (in those words) when the answer is no to both, re-cutting when it is yes to either. CHIP-VINC-NEUTRO r7: hub cut 2,662 lines at `293c1485` and told nobody; the chip then landed `2b956e19` deleting a genuinely dead mock and the comment that had predicted its own falsity — a correction the cold seat independently filed as a REPORT on the same locus. Delta: one test file, 18/20, zero production lines; round stood. Second hub-side instance of the same shape as the mid-wave merge (`678a6d51`): the hub changed shared state and did not tell the party it bound
2026-07-28 · §11 (existing subsection extended) · corrected · "IT ARRIVED IN THE NOTIFICATION" IS NOT CUSTODY. The rule already said the orchestrator writes the artifact in the same act in which the verdict arrives — the hub broke its own rule and lost a verdict. A seat return is delivered in two shapes by SIZE: short verdicts are inlined in the completion notification and survive in the transcript; long ones go to `tasks/<id>.output` with only a summary inlined. Measured after the fact: four seat output files, all `0 bytes`, and no block of the seat's text anywhere in the transcript. CHIP-VINC-NEUTRO r7 — Sol survived at 5,625 chars because it was inlined; the Opus verdict of the same round is GONE, leaving only the hub's paraphrase, which is not verbatim and does not become one by being the last copy standing. The window was between ARRIVING and BEING PASTED and it was crossed by a context compaction, so nothing signalled the loss. Binds: the bigger and more valuable the verdict, the likelier it is the one that vanishes — paste on the turn it lands, not on the turn the ruling is written; a lost verdict is DECLARED in the pack with the measurement that established the loss, and findings that existed only in the paraphrase do not become orders; and recovery by SIZE picks the wrong round — the hub grepped for the largest matching block and pasted CHIP-ANCHORS-2's `refforn`/`RouteClassMux` gate under this round's heading before checking the header. Match on the seat's task id, never on size
2026-07-28 · §11 (existing subsection extended) · ratified · THE POPULATION IS SET BY THE FACT, NEVER BY THE EDIT FOOTPRINT. The commonest way a sweep pattern comes out narrow is that nobody chose one — it silently runs over the files already open. CHIP-VINC-NEUTRO r7: a merge moved `ImportacaoSection` out of `pages/vinculos/` and killed a `listErpImports` mock in TWO test files; the chip deleted the one in the file it was editing, the second survived a full round, and the hub caught it with `grep -rn` over the directory (`VinculosDesign.golden.test.tsx:25,31,120,121`). The cause was a moved component, so the population is every file that mocked the port, tree-wide; the edit footprint has no bearing on it. Binds: ask what KILLED the thing and let that answer pick the search root — a sweep whose root is the diff is reporting the diff back to you. Chip-side field finding, named by the chip against its own work (`1fcf7f1a`) and verified at the hub by string before ratification
2026-07-28 · §11 (new subsection) · ratified · A MUST-FAIL ARM PROVES ONLY WHAT ITS MUTATION ISOLATES, AND NEVER CATCHES AN OVER-STRICT GUARD. Two limits of the instrument, both found chip-side by CHIP-VINC-NEUTRO r8 against its own arms. (1) A mutation that breaks several rules at once goes red without telling you which rule the arm measures: deleting a whole `if` line turned FOUR arms red because signals then reached an absence-only check and threw the wrong message there; the chip discarded it for a two-line reorder that moves exactly one arm, and declined to report the four-red result as stronger evidence. Narrow until one arm moves. (2) A guard that rejects too much makes EVERY must-fail greener, never redder — no arm in the wall can catch over-strictness, only a MUST-PASS on a demonstrably producible input asserted not to throw. Corollary supplied the same round: a guard wider than the fact is the same defect as one narrower, pointing the other way (`anchor appears twice` would have rejected real candidates — `generation_service.go:551` seeds `title` FOR and `:560-562` appends `title` AGAINST to the same candidate, and the finalizer at `:657-661` dedupes only absences). All three loci verified at the hub by string before ratification
2026-07-28 · §11 (new subsection) · ratified · PRODUCIBLE IS CLOSED OVER THE GENERATOR SITES; REACHABLE IS CLOSED OVER THE DECLARATIONS THAT EXIST. CHIP-VINC-NEUTRO defended V-1 (a row whose motivos are all INCOMPARABLE, rendering an overflow button with zero chips) for seven rounds. Producible — sites emit the shape. Unreachable — `mercado_livre` is the tree's only capability declaration and does not supply `marca`, so every candidate carries a `marca UNAVAILABLE`, which the old enumeration DID cover; the empty cell needs a provider supplying all four anchors and none exists. Seven dual gates reading the generator could not see it; one live drive did, immediately. The falsehood sat in the prose, not the instrument: the chip's own fixture already read "No declaration emits it … unreachable" while the value claim beside it still said the screen was broken — a mechanism being right does not make the text that sells it right, and only the text reaches the operator. Binds: a must-fail on an unreachable fixture proves the guard DISCRIMINATES, never that the live screen was broken; a gate criterion asserting a screen state is discharged against DATA, and the hub owes the chip that seam; the false premise is REMOVED from the gate artifact (R-25), and the replacement must still fail the old code on the same observable or the criterion was abandoned rather than corrected. Discriminator supplied by the chip against its own seven-round argument; corroborated at the hub by string (`QueueRow.tsx:159` on main omits INCOMPARABLE from the enumeration entirely, so it is unshowable at ANY limit, and `hidden` counts it) and by the live drive arithmetic (`04983aab`)
