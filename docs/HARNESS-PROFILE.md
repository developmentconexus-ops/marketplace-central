# Harness Profile — marketplace-central

**Layer:** REPO (binds with `docs/HARNESS-CORE.md` — vendored from plugin `harness@mnfs-harness`;
mission content lives in `.mnfs/`).
**Status of this file:** BINDING — swap executed 2026-07-16 at the M-01 milestone boundary
(operator-approved). `docs/HARNESS.md` is now a pointer. Provenance below carries over from the
combined doctrine's dated ratifications.

---

## 1. Identity & stack
`status: ratified` · `provenance: 2026-07-15 · extracted from docs/HARNESS.md (operator-ratified 2026-07-15)`

- Go backend (`apps/server_core`, Go workspaces) + React/TypeScript frontend (`apps/web`),
  npm monorepo.
- OS/shell binding: Windows; **PowerShell for all stack ops — never bash, never WSL**.
- Default branch: `master`. Hub checkout must be on `master` before every hub commit
  (field gotcha: chip launch can switch the hub's working dir onto the scaffold branch).
- Contract-first: OpenAPI spec + generated SDK (`packages/sdk-runtime`).

## 2. Verification ladder bindings (core §5)
`status: ratified` · `provenance: 2026-07-15 · docs/HARNESS.md §5 + field finding (governance lane, 2026-07-15)`

- **L0** — `go build ./...` with `GOCACHE=.gocache` · web `tsc`/typecheck ·
  governance lanes: `npm run harness:governance -- -BaseSha <sha>` — run from a **clean
  detached worktree** (main checkout sweeps `.claude/worktrees/*` and false-fails until the
  scanner exclusion lands) and pass the **full 40-hex** BaseSha (short sha =
  `GOV_SEMANTIC_DRIFT id=base-sha-invalid`).
- **GOCACHE must resolve to an ABSOLUTE path on Windows/pwsh** (D-14, M-01): relative
  `.gocache` breaks when the working dir shifts mid-pipeline — bind it as
  `$env:GOCACHE = (Join-Path (Get-Location) '.gocache')` (or equivalent) before Go commands.
- **Governance base anchor (per milestone):** the drift gate's BaseSha for a chip is the
  milestone's ACCEPTED BASE SHA (40-hex, carried in the chip prompt) — on a long-lived
  worktree, drift REDs computed against any other base are not the chip's defect
  (M-01 field, 2×: Slice 7 `base-sha-invalid`, P5 tool-vs-validate topology). The chip
  records the anchor in its evidence; the hub re-runs governance on the integrated default
  branch at acceptance.
- **L1** — `GOCACHE=.gocache go test ./...` (touched packages + guard suites; full sweep only
  when migrations/platform touched) · web vitest · integration lane
  `npm run harness:integration` (see §4).
- **Known pre-existing failure allowlist (L1):** verdicts CITE this list instead of
  re-proving non-linkage from scratch (M-01 re-proved `TestPhase1SmokeFlow` 5×). A ladder
  run is GREEN-with-allowlist when its only failures are listed here unchanged. Editing the
  list is hub-owned; each entry needs an evidence pointer + a backlog owner.
  - `TestPhase1SmokeFlow` (`apps/server_core/tests/integration`) — pre-existing sibling
    breakage, predates M-01 (evidence: M-01 DISPATCH-LEDGER.md hub ruling B); backlog: hub
    queue, unowned.
- **L2** — dev stack up **ONLY via docker compose**: `npm run docker:dev`
  (postgres+backend+frontend, server `:8080` + web `:5174`); OAuth flows:
  `npm run docker:oauth`. HUB-OWNED seam: chips send `REQUEST dev-stack`, never boot their own
  server process, never bind the ports, never load `.env*` into session env vars (env is
  consumed by the container entrypoint, not the session). Smoke: target routes, error shapes,
  OpenAPI ↔ SDK ↔ handler parity; evidence captured to the mission's contract paths.
- Post-merge ladder on integrated master MUST include the clean-worktree governance run.

## 3. Fresh-workspace bootstrap
`status: ratified` · `provenance: 2026-07-15 · M-01 field finding (hermetic lane)`

Standard first act in every new chip worktree, before any hermetic lane:

```
cd apps/server_core && GOMODCACHE=$(pwd)/.gomodcache go mod download all
```

(~130M; env prep, NOT a dep change — no REQUEST needed.)

**False-alarm signatures:**
- `HPG_MIGRATION_FAILED` with `migrations_first=-1` = build died before migrate ran (empty
  `.gomodcache` under `GOPROXY=off`/`GOSUMDB=off`), not a SQL/migration defect. Warm the cache
  before diagnosing SQL.
- 3D000 "database does not exist" on first CREATE DATABASE attempt = postgres first-boot init
  restart race (pg_isready passes during it); the lane's retry loop absorbs it — not a defect.

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

## 11. Code-review bindings (REVIEW-STANDARD.md instantiation)
`status: ratified` · `provenance: 2026-07-16 · operator-requested review hardening`

- Standard: `docs/REVIEW-STANDARD.md` (method copy; canonical in mnfs-harness
  `harness/REVIEW-STANDARD.md`) — binds slice review, dual gate, hub spot-check.
- Deterministic pre-pass (§7) = L0 lane: `go build ./...` + `go vet ./...` (both
  `GOCACHE=.gocache`) · web `tsc` · governance lane (clean worktree, 40-hex BaseSha). Review
  dispatch only after L0 green; reviewer receives the L0/governance report as input.
  Candidate additions (staticcheck, dupl, gocognit) = dependency gate: `REQUEST` to hub first.
- Learnings file: `docs/REVIEW-LEARNINGS.md` — loaded into every code-review dispatch.
- Slice size norm: ≤ ~300 changed lines (mechanical/generated diffs exempt when declared).
- Severity taxonomy + dual-gate agreement merge as per standard §5/§8; reconciliation table
  required in every `CLOSED` event.
- Execution model (standard §13-§16, ratified 2026-07-16): slice review = ONE reviewer via
  prompt-pack §14 (never a crew); dual gate dispatched simultaneously; disjoint-slice overlap
  allowed (§15 — reviewed-green before merge/dependent slice, not before every next slice);
  ★ crews run under noise control (FAIL-restraint verbatim-quote rule, advisory cap 10,
  learnings suppression via `docs/REVIEW-LEARNINGS.md`).
- Reviewer model/reasoning bindings (standard §13 restates core §1 matrix — canonical there):
  slice reviewer = Claude subagent **sonnet** (session-default effort); dual gate = **COLD
  Opus subagent** (clean context, explicit `model=opus`, never the orchestrating session —
  self-grade bias) + **GPT-5.6 Sol `--effort medium`**; ★ crews = cold Claude **sonnet**
  subagents, mission readiness P7 adds **GPT-5.6 Sol `--effort high`**; hub spot-check =
  sonnet subagent. Every review context is cold — bounded inputs only, never the
  implementation conversation. GPT flags NEVER retyped from memory — resolve via the
  `codex-dispatch` skill; `--effort` always explicit.

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
2026-07-16 · §10 · ratified · mnfs-workflow execution-layer skills denylisted (deleted at source in mnfs-harness 6b29412 layered unification; stale codex cache 0.1.0 + ~/.codex/plugins/mnfs-codex-plugin still ship them — operator field finding: CHIP-SAT worker auto-loaded feature-execution). General rule: auto-discovered skills never bind; only prompt-pack pins are doctrine. Cache repackage to 0.2.0 deferred to W1 close (no tooling swap under running workers).
```
