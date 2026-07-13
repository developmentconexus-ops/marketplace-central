# F-15 Partial Tax Snapshot Quality Validation

## Verdict

The bounded profitability correction passed its focused fake/unit proof.
Known non-complete tax amounts remain present, but item and order snapshots now
carry `missing_tax`, remain incomplete, and suppress contribution/margin. The
fully complete control still calculates complete item and order margin.

The registered repository-wide Go command did not exit zero because the
untouched inventory test `TestStockRiskServiceClassifiesOversellAndFilters`
failed and reproduced in isolation. Profitability passed within that run and
again as a complete module. No inventory edit is authorized or included, so
this outside-scope verification conflict remains for Milestone/QA disposition.

## Context evidence

- Accepted/frozen base SHA:
  `ef4b08c78d30a5e2269e79b051a432c9dc12b58d`.
- The L1 post-plan context was compiled from the Feature contracts and the
  committed F-14 specification, then validated with `-RequireCurrentBase`.
- Estimated input is 1,138 tokens.
- Artifact: `context.json` in this Feature directory.
- The reviewer-owned uncommitted `milestone-review.md` was read as dispatched
  finding evidence but was not made a source-hash dependency or modified.

## Implementation evidence

- `applyInput` continues adding every known tax amount to `TaxAmount`.
- The same tax branch now records `missing_tax` when the amount is nil **or**
  the input quality is not `complete`.
- Existing item finalization suppresses contribution/margin when that flag is
  present, and existing item-to-order rollup propagates the flag and retained
  tax total to the order snapshot.
- The import-to-snapshot regression covers:
  - partial lineage with ICMS/IPI/PIS/COFINS all known and complete at source;
  - complete lineage with all components known but stale/incomplete source;
  - complete lineage with all components and source quality complete.
- Complete synthetic snapshot fixtures now declare their tax input quality
  explicitly instead of relying on an empty quality value as implicit success.
- F-14 validation has an append-only correction notice identifying the old
  input-only/nil-PIS coverage gap. The correction ledger appends attempt 2,
  keeps the last QA result failing/pending re-gate, and retains C03 as
  explicitly deferred/failing.

## go-profitability-partial-tax

From `apps/server_core`, using the repository-local `.gocache` as `GOCACHE`:

```text
go test -count=1 ./internal/modules/profitability/...
```

Result: exit 0. Profitability application, adapters, domain, ports, and
transport compiled; all profitability tests passed. The focused regression
`TestImportMarginInputsPropagatesTaxQualityIntoSnapshots` also passed alone
with `-count=1`.

## go-repository

From `apps/server_core`, using the repository-local `.gocache` as `GOCACHE`:

```text
go test -count=1 ./...
```

Result: exit 1. The only reported failure was outside F-15:

```text
--- FAIL: TestStockRiskServiceClassifiesOversellAndFilters
    stock_risk_service_test.go:74: len(items)=0, want 1
FAIL marketplace-central/apps/server_core/internal/modules/inventory/application
```

The profitability application package passed in the repository run. A fresh
isolated rerun of exactly the inventory test reproduced the same failure.
`git diff` contains no inventory path. F-15 does not claim a green repository
gate and does not broaden into unrelated inventory business logic.

## git-diff-check

- `git diff --check`: exit 0; Git emitted only line-ending conversion warnings.
- `git diff --cached --check`: exit 0 after staging only the five dispatched
  implementation/evidence paths plus the F-15 Feature directory.
- The parent-owned modified milestone review and untracked fixed-SHA review,
  QA, research, probe, and output artifacts remain unmodified and unstaged.
- Current-base context validation passed again before commit.

## Acceptance mapping

- F15-AC01: partial lineage with all four known tax components retains tax 8
  at item/order while both stay incomplete with nil contribution/margin.
- F15-AC02: complete lineage with stale/incomplete exact-line source quality
  produces the same retained-amount and suppressed-margin result.
- F15-AC03: complete lineage plus complete source quality produces complete
  item contribution/margin 72/72% and order contribution/margin 67/67%.
- F15-AC04: F-14 overclaim correction and attempt-2 ledger entry are present;
  C03 remains deferred/failing.

## Scope and handoff

- No lineage, Oracle, orders, auth/manual-adjustment, API/SDK, migration, UI,
  runtime configuration, database/provider, dependency, secret, or PII change
  or live operation occurred.
- Feature correction status: implemented and focused proof passed.
- Verification blocker: repository-wide Go exit 1 from the reproducible,
  untouched inventory test described above.
- Next: Milestone Orchestrator freezes the returned commit for replacement
  independent review and proportional QA; only QA may change milestone status.
