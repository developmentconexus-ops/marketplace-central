# Planning checklist — by area

Read the sections your change touches. Each item is phrased as something the **plan** must
answer, because an unanswered item becomes the executor's improvisation.

## Contents

1. [Architecture & boundaries](#1-architecture--boundaries)
2. [Contract seam — OpenAPI ↔ SDK ↔ handler](#2-contract-seam--openapi--sdk--handler)
3. [Data & migrations](#3-data--migrations)
4. [Integration & provider](#4-integration--provider)
5. [Sync, scheduling & failure visibility](#5-sync-scheduling--failure-visibility)
6. [Frontend](#6-frontend)
7. [Verification lanes](#7-verification-lanes)
8. [Test design](#8-test-design)
9. [Honest values — no hardcode, no fabricated fact](#9-honest-values--no-hardcode-no-fabricated-fact)
10. [Legacy, redundancy & the global maximum](#10-legacy-redundancy--the-global-maximum)
11. [Collision adjudication](#11-collision-adjudication)
12. [Evidence, debts & operator gates](#12-evidence-debts--operator-gates)

---

## 1. Architecture & boundaries

- Which module owns this behaviour? If two could, the one that owns the **data** wins; a
  second module reads it through a published interface, never through the first module's
  `adapters/`.
- Does the change stay inside `domain → ports → application → adapters → transport`? Name
  the layer of every file the plan creates.
- Does any new cross-module import exist in that module's `dependencies` in
  `contracts/governance/modules.json`? If not, the plan either adds the edge deliberately or
  routes through an existing published interface.
- Provider and driver types (`godror`, ML JSON payloads, `pgx` rows) stay in `adapters/`.
  If a provider field name appears in `application/` or `domain/`, the boundary already broke.
- New module → registry entry in `modules.json`, and the entry lands **with the merge**, not
  before it.
- Is `composition/root.go` touched? It is an exclusive seam and ~940 lines; one owner, one
  task, and the plan states the wiring verbatim rather than "wire it up".
- Would this be a new module? Prefer a new *layer file* in an existing module. New modules
  cost registry, composition, governance and review surface; justify with the boundary they
  create, not with tidiness.

## 2. Contract seam — OpenAPI ↔ SDK ↔ handler

- Does the operation **already exist**? Grep the spec for the concept and for zero-call-site
  `operationId`s. Fully-wired-but-unused operations are common here; the work is often a
  button, not an endpoint.
- OpenAPI + `packages/sdk-runtime` in the **same task and same commit**
  (`GOV_API_SDK_SPLIT`). Never a task "add the endpoint" followed by a task "add the SDK method".
- Where does the type belong in the SDK? Grep all of `packages/sdk-runtime/src/` for the type
  name; domain types live in sibling files (`market.ts`, `dashboard.ts`, …), not always
  `index.ts`.
- Does the plan extend the transport-side `openapi_contract_test.go` for the new path block —
  `operationId`, request/response `$ref`s, and each status code the handler can return?
- Error envelope: does the handler return an existing `apierror` code, or a new one? A new
  code is a contract change and gets asserted in the SDK error-contract tests.
- Status codes: does the plan enumerate every non-200 the handler produces, and does the SDK
  narrow them? A contract that lists only the happy path is how a 409 becomes a blank screen.
- Does the web reach this only through the SDK? Direct `fetch` in `apps/web/src` or
  `packages/` fails `GOV_FRONTEND_FETCH`.
- Backwards compatibility: is any existing field's meaning changing? Renaming a field is a
  break for every consumer; the plan lists them.

## 3. Data & migrations

- **Column and table names are measured, not recalled.** `\d <table>` against the dev
  database, and the output pasted into the plan's measurement section. Inventing
  `last_success_at` costs a round.
- New migration → unique numeric prefix (`GOV_MIGRATION_PREFIX`; duplicate prefixes already
  needed an exception once) **and** a bump of the hardcoded count in
  `internal/platform/migrate/runner_test.go` in the same task.
- Forward-only. There are no down migrations; a mistake is corrected by a new migration.
- `tenant_id` on every business table, and a `tenant_id` predicate on **every** query the
  plan writes. Not "the service scopes it" — the predicate is in the SQL.
- Is the read cheap? State the index the query uses. A per-row query inside a page render is
  a defect the plan should catch, not the profiler.
- Does the plan add a table where a column would do, or a column where an existing table
  already carries the fact? Cite the existing shape before extending it.
- Migration numbers are a hub-allocated seam in parallel work — the plan states its number
  and does not grab blind.

## 4. Integration & provider

- Which capability/port does this call go through? Provider access is only via `connectors`
  adapters behind ports — never a direct HTTP call from a business module.
- **Call budget** against the shared 900/min token bucket: calls per item × items per cycle ×
  cycles per hour, and what starves if it saturates.
- What happens on `ErrCodeProviderRateLimited`, on 401/403, on a token that failed to
  refresh? Each needs a stated behaviour, and each needs to be **visible** (see §5).
- Provider writes require, per non-negotiables: resolved linkage, explicit policy/source
  time, duplicate protection, and audit. A plan proposing a write states all four.
- Live provider writes need explicit operator authorization named in the plan. Never assume
  standing authorization from a previous mission.
- Raw provider payloads: PII is scrubbed at capture, and no payload is persisted verbatim
  without saying why.
- Oracle is read-only. Oracle SQL stays in `internal_read/adapters/oracle`; a downstream
  module never writes ad-hoc SQL against it.
- Pagination: does the plan prove it handles >1 page? Silent page-1 truncation is invisible
  to a live drive against a small account — it needs a fixture with more than one page.

## 5. Sync, scheduling & failure visibility

- Is the entity already declared in `sync/domain/entity.go`? Declared ≠ registered — check
  whether a job exists for it, and check `sync_state` for a live row.
- Register a `JobFunc` on `syncapp.Scheduler` rather than building a second ticker, unless
  the plan writes down why the seam does not serve.
- Cursor semantics: what does the job persist, and what does it read on the next run? A job
  that ignores its cursor re-does work every cycle.
- Interval: is it short enough that the missing boot-run is harmless? A long interval starves
  after every restart until the boot-catch-up debt is paid.
- **Failure visibility is part of the feature.** Trace the path from the failure to a pixel:
  which call records it, which read exposes it, which component renders it, and what the
  operator sees. If the answer is "the log", the operator never sees it, and the screen stays
  green over a dead job. This is the single most repeated defect class in this repo.
- Does a *partial* failure look different from success? A job that processes 50 items and
  fails 49 must not report the same thing as one that succeeded.
- Does the plan avoid a blind sweep? Prefer renewing what is already stale over discovering
  everything, and keep discovery on the operator's explicit click where that was ratified.

## 6. Frontend

- Where does the defect actually live — the page, or a shared package? A formatter in
  `packages/web-query/src/index.ts` is on every screen that imports it. Count the call sites
  before choosing where the fix goes.
- Is there already a component for this in `packages/ui/` or `packages/web-query/`? Two
  components with the same name and different `aria-label`s have shipped here, which silently
  splits every label-based test.
- Is the `packages/feature-*` surface you are editing actually mounted? Grep from
  `apps/web/src`. Unmounted packages are dead code; editing them changes nothing on screen.
- Data through the SDK only. No `fetch`, no direct URL construction.
- Loading, empty, error and stale states: each named, each with its exact rendered text.
  "Handle the error case" is a placeholder, not a plan.
- Accessibility labels are test contracts. Changing an `aria-label` breaks tests in files you
  are not editing — the plan lists them.
- Does the screen ever show a value whose **age** matters? Then the age is rendered in a form
  that can express the domain's staleness threshold. A time-of-day stamp cannot distinguish
  fifteen minutes from fifteen days.

## 7. Verification lanes

- Every command in the plan is copy-pasteable **with its working directory**.
  - FE: `cd apps/web && npx --no-install vitest run`
  - Go: from `apps/server_core`, with an absolute `GOCACHE`
  - Governance: clean worktree, full 40-hex base SHA
- Integration tests: `//go:build integration` in the **first lines** of the file, package
  under `internal/modules/**` or `tests/integration`. The lane self-discovers; transport
  packages do not compile in it.
- Does the plan distinguish *skipped* from *green*? A lane that ran zero tests and a lane
  that passed look byte-identical in the summary — assert on a failure token or a named test,
  not on `status=passed`.
- Known pre-existing reds are cited from HARNESS-PROFILE §2, not re-proved.
- A live drive names: the URL, the click path, the exact expected rendered text, and its
  **precondition** — the running container was built at or after the slice's commit. A stale
  binary makes absence-of-defect and absence-of-code indistinguishable.
- Dev stack: `docker compose` only, hub-owned. The plan never tells an executor to boot a
  server, bind `:8080`/`:5174`, or load `.env` into the session.

## 8. Test design

- **Red before green, with the failure text.** The plan shows the exact expected failure
  message. "Verify it fails" without the message is not a control.
- **Presence is not value.** `toBeInTheDocument()` passes with any string. If the defect is a
  wrong value, the assertion is on the value. This class survived an entire test suite here.
- **A must-fail arm proves only what its mutation isolates.** State which line the injected
  defect changes and which assertion catches it.
- Someone else's test that this change breaks is **restored or restated**, never deleted. If
  the responsibility moved to another file, the plan says which file now carries it.
- No mock or stub stands in for an integration seam. Mocks prove contract behaviour;
  integration criteria run against the real dependency.
- Composition roots ship no permanent stub/nil wiring on a live path. A stub is allowed only
  with a dated deferral naming the slice that replaces it.
- Timestamps: Windows `time.Now()` is 100ns, Postgres `timestamptz` is µs. Any round-trip
  equality fixture truncates to microseconds first, or it fails on the 7th fractional digit.
- Does the test assert something that would pass in **both** worlds — with and without the
  fix? Then it is not evidence. Name the discriminating observable.
- Fixtures that are symmetric in the property under test assert nothing. If the fixture would
  pass with the inputs swapped, it does not test order.

## 9. Honest values — no hardcode, no fabricated fact

- **ADR-17:** unknown never becomes `0`, `""`, `false`, or a plausible default. It becomes an
  explicit state that reaches the screen as unknown.
- No business constant as a literal: tax rates, tariffs, fees, TTLs, company codes, warehouse
  codes, thresholds. They come from configuration, the database, or the provider. A literal
  labelled "temporary" has never been temporary here.
- An empty table is not a fact. If a configuration table has zero rows, the code says
  "unconfigured", it does not substitute a default and present it as measured.
- A default that is *plausible* is worse than one that is obviously wrong, because nobody
  checks it. Prefer failing honest over guessing well.
- No blanket `recover`/fallback on an integrity-critical read.
- If the plan computes a number the operator will act on, name its source per input, and name
  what happens when one input is missing.

## 10. Legacy, redundancy & the global maximum

- Count the copies of the concept before writing the fix. ≥2 → the fix belongs where they
  converge.
- Is the symptom's cause one layer down? Fixing the screen when the package is wrong leaves
  every other call site broken and makes the next occurrence look new.
- Does a seam already exist (scheduler, port, formatter, error envelope)? Reuse inherits its
  reconciliation, isolation and failure visibility; a parallel implementation inherits none.
- Is this extending VTEX or another legacy abstraction to solve a Mercado Livre problem?
  Frozen decision 11 forbids it. Legacy is inventoried, or removed in its own slice.
- Is the code you are about to extend actually reachable? Dead code that compiles is common;
  grep for its callers before planning against it.
- When the global fix is bigger: take it, scope it honestly, list every call site, and put
  the refactor in its own committed task **ahead** of the feature. A deferred refactor here
  has never happened later.
- When it is genuinely out of scope: register it in `.mnfs/HARNESS-DEBTS.md` with the
  measurement that proves it, in the format *measured fact → cost paid → candidate fix*.

## 11. Collision adjudication

- List the in-flight branches: `git worktree list`.
- For each, `git diff --name-only main...<branch>`, and show the intersection with this
  plan's file set. Empty intersection is a measurement; "they're different areas" is a guess.
- Exclusive seams need a named owner for the duration: the api-sdk pair, `migrations/`,
  `composition/root.go`, the dependency graph, ADRs, harness control files, the
  provider-capability contract.
- Migration numbers and contract sections are pre-allocated, never grabbed blind.
- Base SHA: a governance verdict is only attributable if the tree and tip it measured are
  named. Cite both.

## 12. Evidence, debts & operator gates

- Evidence lives in `.mnfs` artifacts. Unwritten = didn't happen; the plan says which
  artifact each verification writes to.
- Dependency change is a `REQUEST`, never an install-as-ritual step in a task.
- Push requires explicit operator permission. Committing verified work does not.
- The plan never instructs: reset, revert, stash, clean, delete unknown state, use WSL,
  cold-clone, purge caches, read or print `.env*`, dump an environment, `git branch -D`, or
  write to Oracle.
- Live Mercado Livre writes name the operator authorization they require, per mission.
- The closing task lists the debt IDs discharged and those left open, each with its
  measurement.
