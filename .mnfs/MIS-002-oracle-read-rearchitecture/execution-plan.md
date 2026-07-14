# MIS-002 — Parallel Execution Plan (Hub-Orchestrated)

```yaml
type: execution-plan
mission: MIS-002
owner: Portfolio Hub (session local_00112a95)
created: 2026-07-14
model_workers: mpc-implementer / mpc-verifier (gpt-5.6-luna, high reasoning)
orchestration: hub spawn_task per wave; standalone Milestone sessions in fresh worktrees; terminal callback via cross-session messaging
```

## Dependency Analysis (what actually blocks what)

Hard dependencies (cannot parallelize):
- M-01-C04 (live scale + EXPLAIN PLAN verdict) → gates all new SQL (M-02 F-01, M-03 ports).
- M-02 F-01 port → M-02 F-02 envelope/OpenAPI/SDK (consumes the port).
- M-03 F-01 chunker → M-03 F-02 and F-03 (both consume `oraclebatch`).
- M-02 F-02 SDK method → M-05 (hooks call the typed SDK client).
- M-02/M-03 final port shapes → M-04 decorator wiring.

NOT dependencies (parallelizable):
- M-02 ⟂ M-03: disjoint modules (catalog vs inventory/profitability/product_links); both only ADD new files under `internal_read` (page port vs `oraclebatch` pkg). Shared touchpoint: `composition/root.go` wiring — trivial merge.
- M-04 ⟂ M-05: server-only (cache adapter + composition) vs web-only (`apps/web`). Zero shared files.
- M-03 F-02 ⟂ F-03: after F-01 chunker lands, cost/tax batch and Sankhya linkage touch different modules.

## Waves (3 instead of 5 serial windows)

| Wave | Runs in parallel | Gate to enter | Merge order into main |
| --- | --- | --- | --- |
| W0 | Commit current working tree (adapter refactor + MNFS artifacts) | operator approval | one commit, direct |
| W1 | **M-01** (single session; F-01 mostly already in tree → fast) | W0 done | `merge --no-ff` after QA pass |
| W2 | **M-02** ∥ **M-03** (two sessions, two worktrees) | M-01 QA pass + M-01-C04 verdict recorded | M-02 merges FIRST (defines envelope pattern); M-03 rebases on merged main, resolves `composition/root.go`, then merges |
| W3 | **M-04** ∥ **M-05** (two sessions, two worktrees) | W2 both merged | independent; M-04 first if simultaneous (server before client staleTime claims) |

Intra-milestone parallelism (declared in each `milestone.md` under "Feature Parallelization" — orchestrator MUST use it):
- M-01: F-01 code ∥ F-02 baseline-evidence lane (live docker lane, writes only evidence md); F-02 instrumentation after F-01.
- M-02: interface-first handshake — F-01 lands port interface/types as step 1; F-02 (transport/OpenAPI/SDK vs fake port) dispatches immediately after that step, not after F-01 completion.
- M-03: F-01 ∥ F-02 ∥ F-03 from minute one — chunker API `oraclebatch.Chunks(ids []int64, max int) [][]int64` contract-fixed in milestone.md; integration order A→B→C.
- M-04: single feature, no internal split.
- M-05: F-01 QueryClient wiring as step 1; F-02 invalidation parallel after it (namespaces contract-fixed in IC-01).

Worker fleet peak: W2 runs up to 5 gpt-5.6-luna workers simultaneously (M-02 ×2 + M-03 ×3), each with its own verifier pass.

## Seam Ownership Matrix (one writer per seam per wave)

| Seam | W1 | W2 | W3 |
| --- | --- | --- | --- |
| `internal_read/adapters/oracle` existing files | M-01 | reference-only (both add NEW files) | — |
| `internal_read/ports` | M-01 | M-02 (catalog port) + M-03 (batch ports) — new files each | — |
| `composition/root.go` | M-01 | BOTH (declared conflict point; M-03 rebases) | M-04 |
| catalog module + transport routes | — | M-02 | — |
| inventory / profitability / product_links | — | M-03 | — |
| OpenAPI + `packages/sdk-runtime` | — | M-02 only | — |
| cache adapter pkg + write-service invalidation calls | — | — | M-04 |
| `apps/web` | — | M-02 (minimal envelope adjust only) | M-05 |

Rule: a wave-2 session finding it must edit a seam owned by the other lane STOPS and asks the operator (visible session) instead of editing.

## Hub Protocol

1. Hub spawns one chip per Milestone session (spawn_task). Fresh worktree per session — hence W0 commit prerequisite.
2. Each Milestone session: MNFS + hub-style acceptance (verifier per feature, fixed-SHA review, QA gate), then persists `checkpoint.md` and sends verdict + path to hub session `local_00112a95` via cross-session messaging.
3. Hub folds callbacks; when a wave's gate is satisfied, spawns the next wave's chips (W2: two chips simultaneously; W3: two chips simultaneously).
4. Merges to main happen in the declared order; the later lane rebases before its merge. Hub verifies post-merge test ladder (`GOCACHE=.gocache go test ./...`, `npm run build`) per merge.

## Quality Guarantees Preserved

- QA-2 gates unchanged: every milestone still passes its validation contract; only QA passes a milestone.
- M-01-C04 gate unchanged — no new SQL before plan evidence.
- Parallelism never shares a seam writer; conflict surface reduced to `composition/root.go` (mechanical wiring).
- staleTime/TTL coupling (M-05 mirrors IC-01) is contract-fixed, not runtime-coupled — safe to build M-05 while M-04 in flight.

## Wall-clock Estimate

Serial: 5 milestone windows. Parallel: 3 waves (W1 short — adapter work mostly landed). Critical path: W0 → M-01 → M-02 → (M-04 | M-05). M-03 fully off critical path.
