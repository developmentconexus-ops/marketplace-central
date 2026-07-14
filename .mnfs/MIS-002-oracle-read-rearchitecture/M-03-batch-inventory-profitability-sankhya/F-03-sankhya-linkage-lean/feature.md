# F-03-sankhya-linkage-lean

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-002. Security criterion MIS-02-C04; cache exclusion MIS-02-C05.

## Milestone

M-03. Depends on F-01 helper (IN-list chunker).

## Brief

Lean the Sankhya linkage reader: (1) ValidateConfiguration (ping + 2 metadata queries) runs once at startup, never per request; (2) per-candidate line queries collapse to one IN-list query for all candidate ids; (3) raw driver causes at `sankhya_linkage_reader.go:367-369` routed through `safeOracleCause`; (4) assert linkage results are never cached (no FreshnessPolicy > 0 path).

## Inputs

- `apps/server_core/internal/modules/product_links/.../sankhya_linkage_reader.go:104,160` (per-call validation), `:150-155` (line N+1), `:367-369` (redaction gap).
- `internal_read/adapters/oracle/reader.go:508-524` — `wrapOracleError`/`safeOracleCause` reference implementation.
- F-01 `oraclebatch` chunker.

## Expected Output

- Startup-once validation (composition root; failure = startup error, not per-request latency).
- Candidate scoring: 1 candidates query + 1 chunked lines query.
- All error paths redacted; regression test with DSN-shaped fake error.
- One intentional commit.

## Constraints

- Linkage correctness untouched — same scoring inputs, fewer round trips. Recent commit `97fd4b58 fix(product-links): require canonical identity` behavior must be preserved (canonical identity requirement stays).
- Linkage data NEVER cached (staleness → wrong provider writes; mission ADR).
- No provider-write changes; read path only.
- `GOCACHE=.gocache`.

## Inputs/Outputs

Reader interface unchanged for callers except removed validation latency; internal query plan changes only. Candidates list ordering preserved: match score descending, tie-break `internal_product_id` ascending (IC-01 semantics row). Error type shape preserved (callers match on typed errors, not strings).

## Negative Scenarios

- While Oracle rejects credentials, when startup validation runs, the system shall fail startup with redacted cause (no username/DSN in message).
- While a candidate id set spans 700 ids, when lines load, the system shall issue exactly 2 line queries.
- While the reader errors mid-scoring, when the caller receives the error, the system shall see typed error + redacted cause only.

## Validation Expectations

- Fake-queryer count: k candidates → 2 queries total (candidates + lines), k>500 → 1+2.
- Redaction test asserting absence of user/host/service substrings.
- Grep/test evidence no per-call ValidateConfiguration remains.
- `GOCACHE=.gocache go test ./...` green.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` at execution time.

## Handoff

- Current status: briefed
- Next owner: Feature Implementer (mpc-implementer, gpt-5.6-luna high)
- Next action: spec → plan → implement → evidence
- Required files/evidence: `validation.md`
- Blockers or open decisions: None
