# F-01-adapter-hardening-pool-deadlines

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-01
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-002 (`../../mission.md`), IC-01 (`../../research/catalog-read-interface-contract.md`).

## Milestone

M-01 foundation-observability.

## Brief

Land the pending (uncommitted) Oracle adapter hardening refactor as one intentional commit, raise the pool default to 12, and add HTTP timeout discipline: `http.Server` Read/Write/ReadHeader timeouts plus per-route-class context deadlines (interactive 15s, batch 120s) per IC-01's route-class table.

## Inputs

- Working tree already contains the refactor: `internal_read/adapters/oracle/{config.go,config_test.go,open_cgo.go,open_nocgo.go,reader.go,reader_test.go,reader_live_test.go,database.go(new)}` + `composition/root.go` type change. Review, keep, finish — do not rewrite.
- `cmd/server/main.go:35` uses bare `http.ListenAndServe` — replace with configured `http.Server`.
- IC-01 `## Batch & Route-Class Rules` lists which routes are batch-class.
- ADR-04 in mission.md.

## Expected Output

- One commit: adapter refactor + `defaultPoolMaxSessions = 12` + server timeouts + route-class middleware applying `context.WithTimeout` (15s/120s) with 504 `deadline_exceeded` response on expiry.
- Route class declared where routes are registered (composition/transport), not guessed per handler.

## Constraints

- Do not add retry, shutdown handling, or /readyz (Resilience declined — Non-Scope).
- Do not touch query shapes (M-02/M-03 own that).
- No SQL/driver types outside the adapter package.
- `GOCACHE=.gocache` for tests. No reset/stash/clean of the existing working tree.

## Inputs/Outputs

504 body: `{"error":"deadline_exceeded"}` (IC-01 error matrix). Env knob for pool stays `MPC_ORACLE_POOL_MAX_SESSIONS`; new default 12; validation rules from the refactor unchanged.

## Negative Scenarios

- While a handler exceeds its interactive 15s budget, when the deadline fires, the system shall return 504 `deadline_exceeded` and cancel the in-flight Oracle call (context propagation).
- While `MPC_ORACLE_POOL_MAX_SESSIONS=0` is set, when config loads, the system shall reject startup config with the existing validation error (refactor behavior preserved).

## Validation Expectations

- `GOCACHE=.gocache go test ./...` green.
- httptest: stalled interactive route → 504 at ~15s; same stall on a batch-class route → no 504 (120s budget).
- Unit: default config yields MaxSessions 12.
- `git log` shows exactly one new intentional commit.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution, not planning.

## Handoff

- Current status: briefed
- Next owner: Feature Implementer (mpc-implementer, gpt-5.6-luna high)
- Next action: create `spec.md` then `plan.md`, execute, return validation evidence
- Required files/evidence: `validation.md` with command transcripts
- Blockers or open decisions: None
