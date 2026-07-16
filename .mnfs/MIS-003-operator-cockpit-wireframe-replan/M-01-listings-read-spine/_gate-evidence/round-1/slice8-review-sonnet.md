# Slice 8 — conforming slice review (REVIEW-STANDARD §13/§14)

```yaml
standard: docs/REVIEW-STANDARD.md (ratified 2026-07-16), prompt-pack §14 verbatim
layer: slice review — ONE independent Claude subagent, model=sonnet, session-default effort (never crew)
reviewed_sha: a595f36c
files: read_service.go + read_service_test.go + composition/root.go (code only; .mnfs excluded)
inputs_bounded: diff@a595f36c + slice8-corrective-brief.md + slice8-plan.md + L0 report + docs/REVIEW-LEARNINGS.md
l0_precondition: go build ./... =0, go vet ./... =0, go test -count=1 ./internal/modules/listings/... =0 (GOCACHE absolute)
verdict: APPROVE (LGTM, zero surviving findings)
```

## Verdict: APPROVE

- **G1 (right globally):** yes — mirrors existing `Summary` degrade (ADR-17), no new pattern; NO-STUB wiring removes real stub vs deferring. Hub-approved brief/ADDENDUM authorizes the root.go touch.
- **G2 (alternatives considered):** yes — `slice8-plan.md:20` documents the filter-path decision (pass-through vs empty vs 503) + mid-scan restart-from-original-cursor (`:18`).
- **G3 (local-max trap):** no — in-scope M-01 C10 corrective, hub-dispatched.

## Findings: none survived verification.

## Receipts (recorded, not findings)

- **Zero-value safety confirmed:** `oracleDB` is a nil `internalreadoracle.Database` interface (database.go:11) when Oracle absent → passed to `NewBatchReader` whose `queryer` param (reader.go:15) receives a true nil interface (no typed-nil trap) → `ensureBatchAvailable` (batch_reader.go:205) `r.db==nil` → clean `SourceUnavailable`, no panic.
- **Dead-type removal clean:** `unavailablePolicyService` + `marketplacesdomain` import gone, zero remaining refs (consistent with green build/vet).
- **Degrade wrapper kept:** `unavailableListingPolicyReader` now wraps real `marketSvc`, degrades only on `ReadErrorSourceUnavailable`; DB/logic errors stay hard.
- **Cursor restart pinned:** `scan`/`scanGroups` restart from original cursor (read_service.go:337-338/192-193) matches plan:18; test `TestReadServiceDependentFilterPassesThroughWhenOptionalFactsUnavailable` asserts exact call counts + cursor identity (ceiling=1 call, cost=2 calls second at original cursor).
- **Get ICMS asymmetry:** ceiling fail → skip cost + nil ICMSWorstCaseByUF (costCalls==0 verified); cost-only fail → matrix from known ceilings, margin null.
- **Hard errors preserved:** installation/policy/repo/timeline unchanged; `TestReadServiceInstallationAndPolicyFailuresRemainHard` asserts zero downstream calls on early failure.
- **Noted (below threshold, not flagged):** `enrich` error return now structurally always nil (all failure→`factsUnavailable=true`); harmless, matches plan's contract change, callers handle it.

## Notes

- §11 size (~355 code lines): reviewer did NOT raise a size finding; fix+wiring serial by hub direction.
- Supersedes the earlier non-standard cavecrew review (which also returned CLEAN — corroborating).
