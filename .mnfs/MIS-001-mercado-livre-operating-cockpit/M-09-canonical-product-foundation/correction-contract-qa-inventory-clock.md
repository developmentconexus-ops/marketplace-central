# M-09 Inventory Clock QA-Unblock Contract

```yaml
id: M-09-QA-CORR-01
type: portfolio-qa-unblock-contract
status: authorized
owner: M-09 Milestone Orchestrator
parent: M-09
authorized_by: Portfolio Hub
base_sha: 97fd4b58d55a7d14a2b45f0c3bae15b2e374822a
attempts_allowed: 1
created: 2026-07-13
```

## Classification

This is not a retry of M-09-CORR-03. Fixed-SHA review passed C01, and inventory
has no source changes between the accepted planning base and the QA SHA. The full
Go lane exposed a historical time-dependent test: it fixes both observations at
`2026-07-09T12:00:00Z`, while production evaluates freshness against the real clock
with a 30-minute maximum age. The test therefore becomes stale as calendar time
advances and can no longer exercise its intended oversell branch.

## Authorized Correction

Make `TestStockRiskServiceClassifiesOversellAndFilters` use a fresh runtime-relative
observation time so it continues to test oversell/filter behavior. Do not change
production clock behavior, freshness policy, risk classification, quantities, or
the assertion being proved.

## Allowed Paths

- `apps/server_core/internal/modules/inventory/application/stock_risk_service_test.go`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-09-canonical-product-foundation/qa-inventory-clock/validation.md`
- M-09 checkpoint, fixed-SHA review/QA evidence, and `validation-result.md`

All production files, M-09 implementation paths, Oracle runner, contracts, SDK,
migrations, M-06, mission, and roadmap paths are forbidden.

## Proof And Flow

1. Run the exact failing test and the inventory application package with
   repository-local absolute `GOCACHE` and `-count=1`.
2. Run the full Go lane once as correction proof.
3. Create one intentional test-only commit and freeze its SHA.
4. Request `mpc-verifier` fixed-SHA review to confirm the change is test-only and
   does not weaken freshness or oversell assertions.
5. After review Pass, rerun proportional QA from the beginning at that SHA,
   including the residue, SDK, runner-contract, and live Oracle lanes previously
   skipped. Only QA may pass M-09.

## Terminal Rule

No retry or path expansion is authorized. Any further deterministic failure or
required production change returns terminal `failed` to Portfolio.
