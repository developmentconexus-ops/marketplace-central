# Milestone Review — M-01-listings-read-spine — round 1

```yaml
milestone: M-01-listings-read-spine
mission: MIS-003-operator-cockpit-wireframe-replan
round: 1
review_sha: e2cde3648a6bdd534afc0ae076a08a93d06d7c7a
diff_base: d0d30d68
reviewer_model: cold independent milestone-reviewer crew (5 subagents) + qa-validator execution pass
folded_verdict: Blocked
folded_status: blocked (missing required live-provider evidence — C10 operator-deferred; NOT a code defect)
```

## Crew composition

Ordering per command race-fix: `qa-validator` execution + live pass ran FIRST and wrote `_gate-evidence/round-1/`; the cold `milestone-reviewer` crew then read the fully-populated evidence.

- **qa-validator** (execution-grounded re-run + API live-drive) — fresh ephemeral-Postgres lanes, `TEST_EXIT=0`.
- **Reviewer A** — ★1 (coverage) + ★5 (traceability).
- **Reviewer B** — ★2 (evidence honesty) + ★3 (verifiability).
- **Reviewer C** — ★4 (integration/composition) + ★7 (security).
- **Adversarial ★2** — second independent evidence-honesty pass.
- **Adversarial ★4** — second independent composition pass.
- Reviewer D (★6 correction integrity) — NOT dispatched: no correction round ran (round 1).

## Folded per-★ result

UNION of all sub-reviewer findings + qa-validator re-run + structural precondition. No sub-reviewer FAIL downgraded to PASS.

| ★ | Criterion | Folded | Cited locus / basis |
|---|-----------|--------|---------------------|
| ★1 | Criteria coverage | **FAIL → routes to Blocked** | `validation-contract.md:164-177` M01-C10 `Required: Yes` unmet; only recorded result is `could-not-drive` at `_gate-evidence/round-1/qa-validator-report.md:39` (operator-deferred, honest — NOT a fabricated pass). C01–C09 all covered with recorded PASS. Per rubric live-runtime rule, `could-not-drive` on a runnable (live-provider) surface routes to **Blocked**, not correction-Fail. |
| ★2 | Evidence honesty | **PASS** | Reviewer B + adversarial ★2 both PASS. Every C01–C09 green cites a concrete `ran` artifact reproduced by `TEST_EXIT=0` (`rerun-lane.log:85-87`, `rerun-lane-c02.log:68-72`). C10 honestly `could-not-drive`/Pending. p95 fresh-run 3.1384 ms recorded as observed (not re-asserting prior 3.2563 ms). C08 governance-lane caveat disclosed, substance re-derived via `git show --stat`. |
| ★3 | Verifiability | **PASS** | Reviewer B PASS. Each contract criterion carries concrete command + observable + named blocking-failure + artifact; each PASS independently re-derivable from its named test (`qa-validator-report.md:30-38`). |
| ★4 | Integration/composition | **PASS** | Reviewer C + adversarial ★4 both PASS across 5 attack vectors: zero `integrations/adapters/postgres` reach-through from listings (`gateway.go:6-8` published boundary only); atomic `pg_advisory_xact_lock` before SELECT/INSERT in one tx (`operation_run_repo.go:31`); 5 routes registered once, no sibling collision (`root.go:506,514`; `http_handler.go:29,137-140`); cross-module coupling via published ports + anti-corruption mapping only. |
| ★5 | Traceability | **PASS** | Reviewer A PASS. Bidirectional closure: F-01→C01/C02/C03/C08(+C10 lane), F-02→C04–C09; every criterion maps to a concrete test name / commit / recorded evidence. No orphan feature. |
| ★6 | Correction integrity | **N/A** | No correction round ran (round 1). Vacuously satisfied. |
| ★7 | Security / tenancy | **PASS** | Reviewer C PASS (HIGH blast radius). tenant_id predicate on EVERY new query enumerated file:line (`repository.go:79,254,292,318,349,379,418,461`; `operation_run_repo.go:39,42,59,118`) — Q2 blocker NOT tripped. Provider payloads confined to adapters (only `translated_error_code` persists); no token/PII leak (`root.go:511` logs err not body). |

## Verdict computation

Fold rule: ALL seven PASS → Pass; any FAIL → Fail; required artifacts absent OR a runnable milestone with no live-driven evidence (`could-not-drive`) → **Blocked**.

Six ★ PASS (★2,★3,★4,★5,★7 + ★6 N/A). ★1 is not a defect-Fail — it is a coverage gap caused solely by an unexecuted **Required** live-provider criterion (C10) honestly recorded `could-not-drive`. That is the textbook Blocked condition: missing required live evidence, never a silent Pass, never a correction-Fail (there is nothing to correct in code).

**FOLDED VERDICT: Blocked** — pending the C10 live-provider-read lane (real Mercado Livre connected installation + OAuth), which the operator explicitly deferred this session, OR an operator-ratified downgrade of C10's `Required` status in the milestone contract.

## Advisory reconciliations (not verdict-changing)

1. `validation-contract.md:127` C07 literal `below_margin` is **stale** pre-D-22 (2026-07-14) wording. Operator-ratified **D-22** (`DECISIONS.md:101`, 2026-07-15) renamed the honest response field to `below_margin_worst_case`, reflected in the BINDING `contracts/api/marketplace-central.openapi.yaml` (`ListingReadModel.below_margin_worst_case`, `ListingSummaryExceptions.below_margin_worst_case`+`margin_unknown`). Under harness truth order OpenAPI outranks `.mnfs`. Code + tests already correct; the `below_margin` token surviving in the filter enum is the distinct query-filter surface, correctly present. Recommend reconciling the `.mnfs` contract wording (doc-only; not a code defect).
2. `validation-contract.md:168` still reads C10 `Required: Yes` while the P6 ledger asserts deferral. Reconcile via the live lane OR an operator-ratified `Required`→deferred amendment.
3. C08/C09/C10 name `validation-result.md` as their evidence artifact; this review authors it now, reconciling the forward-reference.
```
