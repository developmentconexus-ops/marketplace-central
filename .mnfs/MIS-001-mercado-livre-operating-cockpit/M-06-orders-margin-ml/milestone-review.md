# Milestone Review — M-06 (round 3)

reviewer: milestone-reviewer (cold)
reviewed_sha: `81b8a4b12c3fe32c011f3d362ede393dd7484381`
verdict: Fail

## Summary

The fixed SHA contains accepted F-06 quantity/cost correction evidence and
accepted F-07 exact Oracle sale-line provenance behavior. The proportional
read-only Go rerun passed. F-06 closes the round-2 quantity-extension defect;
F-07 correctly requires exact positive Oracle `NUNOTA`/`SEQUENCIA` and, because
the order fact has no owner-verified upstream mapping, deliberately supplies an
empty identity so tax remains missing rather than guessed or defaulted.

M-06 still cannot pass. The contract retains generic non-executable command
labels and rollup-file artifacts, historical load-bearing passes remain
unbound to durable concrete `Evidence type: ran` output artifacts, and the
manual-adjustment write still accepts caller-supplied actor identity without a
verified principal or installation/order/item authorization. The owner-deferred
C03 work is honestly recorded as deferred and failing; deferral is not a
security mitigation or criterion waiver.

## Must-meet (★)

| # | Category | Procedure run | Cited excerpt (file:line) | Verdict | Yes-if (condition to pass) | Defect locus (FAIL only: file:line + offending token) |
|---|---|---|---|---|---|---|
| ★1 | Criteria Coverage | Enumerated `M-06-C01` through `M-06-C03`, mapped all accepted feature and milestone evidence, and searched the milestone artifact scope for `TBD`/`TODO`; none were found. C03 has a recorded Fail rather than an omitted result. | `validation-contract.md:30`, `:45`, `:60`; `validation-result.md:171-173`; `corrections/correction-task.md:31-35` | **PASS** | — | — |
| ★2 | Evidence Honesty | Read every load-bearing proof, then repeated the adversarial pass. F-06/F-07 truthfully mark focused commands as ran and state their limits, but the authoritative historical C01/C02/C03 pass rows still cite prose/results rather than durable concrete command-output artifacts, and several new Feature artifacts point to ephemeral “task execution.” A new focused rerun corroborates mechanics but does not rewrite or silently upgrade deficient historical proof. | `validation-result.md:23-29`; `F-06-quantity-cost-semantics/validation.md:70-87`; `F-07-order-specific-tax-provenance/validation.md:61-75`; `validation-result.md:134` | **FAIL** | Bind every load-bearing milestone claim to a durable concrete path containing exact command/interaction output and `Evidence type: ran`; preserve prior-round history and add fresh QA-owned evidence rather than relabelling old prose. | `F-06-quantity-cost-semantics/validation.md:77` — offending value `Artifact: command output in this Feature execution task`; `F-07-order-specific-tax-provenance/validation.md:67` — offending value `this validation record and package test output in task execution`; `validation-result.md:23-29` — historical `PASS` rows lack concrete ran-output paths. |
| ★3 | Verifiability | Checked each contract criterion for an exact executable command/interaction, observable, blocking failure, and concrete evidence path. Expected behavior and blocking failures are meaningful, but all three contract commands remain generic suite labels and all three artifacts point to the QA rollup rather than captured output. | `validation-contract.md:36-40`, `:51-55`, `:66-70` | **FAIL** | Replace each generic suite label with exact executable command(s) or declared live interaction(s), exact observable values, and durable output-artifact paths. | `validation-contract.md:36` — offending token `orders service/repository tests`; `validation-contract.md:51` — `profitability service tests`; `validation-contract.md:66` — `manual adjustment API/repository tests`. |
| ★4 | Integration / Composition | Enumerated the cross-feature seams: quantity and monetary amount scope, quality flags, tax identity/provenance, Oracle predicates, profitability port/adapter calls, actor shape, and shared API behavior. Repeated the adversarial seam pass. F-06 extends only per-unit CUSSEMICM and retains fee/tax line amounts. F-07 consumes an exact identity contract, queries exact `NUNOTA`/`SEQUENCIA` plus consistency predicates, and does not infer identity from product/date. With upstream mapping absent, the empty identity composes into explicit missing tax. No consume-only seam divergence was found. | `F-06-quantity-cost-semantics/validation.md:17-22`, `:139-153`; `F-07-order-specific-tax-provenance/validation.md:21-26`, `:50-57`, `:131-141`; `apps/server_core/internal/modules/profitability/application/service.go:235-246`; `apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go:357-372` | **PASS** | — | — |
| ★5 | Traceability | Checked backward and forward closure. F-01/F-02/F-03/F-04/F-05 remain accepted evidence for C01-C03/C02/UI requirements; F-06 and F-07 explicitly trace to C02. Every contract criterion has a recorded milestone result and no accepted feature is orphaned. | `validation-result.md:171-173`; `F-06-quantity-cost-semantics/feature.md:33`; `F-06-quantity-cost-semantics/spec.md:64-68`; `F-07-order-specific-tax-provenance/spec.md:79-89` | **PASS** | — | — |
| ★6 | Correction Integrity | Read the retry policy, round-2 result, append-only ledger, and accepted correction evidence. Attempt 1 is within the cap of 2; round-2 Fail remains the baseline; F-06’s prior C02 defect is tied to a new ran evidence path and was rerun; C03 is not silently upgraded. This document performs the required whole-milestone seven-criterion round-3 re-gate. | `validation-contract.md:87-94`; `corrections/correction-task.md:48-57`, `:61-63`, `:71-72`; `F-06-quantity-cost-semantics/validation.md:70-87`, `:139-157` | **PASS** | — | — |
| ★7 | Security Posture | Inspected the manual-write auth boundary and its tests. The HTTP request body supplies `actor`; transport forwards it unchanged; application code trims it and checks only non-empty type/id. No verified principal, installation/order/item authorization, or upstream authenticated boundary is implemented or evidenced. The correction ledger explicitly defers this work and retains C03 as failing. | `apps/server_core/internal/modules/profitability/transport/http_handler.go:100-116`; `apps/server_core/internal/modules/profitability/application/service.go:264-289`; `corrections/correction-task.md:19-20`, `:31-35`, `:71-72` | **FAIL** | Bind manual-adjustment writes to a verified principal, derive/verify actor identity server-side, authorize tenant/installation/order/item scope, and add a Security-typed criterion with durable ran evidence; or record an explicit architecture decline that actually satisfies the rubric’s decline-with-reason rule. Owner deferral alone does not pass C03. | `apps/server_core/internal/modules/profitability/transport/http_handler.go:107` — offending value `Actor ... json:"actor"`; `apps/server_core/internal/modules/profitability/transport/http_handler.go:116` — `Actor: req.Actor`; `apps/server_core/internal/modules/profitability/application/service.go:288` — only empty-string validation protects actor provenance. |

## Adversarial second pass

### ★2 Evidence honesty

No accepted proof was treated as stronger than its artifact. F-07 explicitly
states that no live Oracle execution or real linkage is claimed
(`F-07-order-specific-tax-provenance/validation.md:43-46`, `:111-117`). Its
deterministic evidence proves exact-selection and unknown-state behavior, not a
real resolved-order tax lookup. F-06 proves deterministic quantity semantics,
not integrated live margin. These honest limits do not repair the separate
historical/durable-artifact deficiencies cited under ★2.

### ★4 Integration / composition

The exact Oracle provenance seam is internally consistent: query predicates
bind `d.NUNOTA = :1`, `d.SEQUENCIA = :2`, `i.CODPROD = :3`, and
`d.CODINC = :4` (`reader.go:364-372`), while the profitability caller supplies
an empty `TaxSourceIdentity` because no owner-verified order-to-Oracle mapping
exists (`service.go:238-246`). The resulting missing tax is the required honest
unknown state, not zero/default and not product/date aggregation. The upstream
mapping remains an operational/product gap, but there is no seam contradiction
at this SHA.

## Proportional read-only corroboration

- Fixed SHA verification: `git rev-parse HEAD` returned
  `81b8a4b12c3fe32c011f3d362ede393dd7484381`; `git cat-file -t` returned
  `commit`.
- Command (from `apps/server_core`, with repository-local `GOCACHE`):
  `go test ./internal/modules/orders/... ./internal/modules/profitability/... ./internal/modules/internal_read/adapters/oracle ./internal/modules/internal_read/adapters/fake -count=1`.
- Result: exit 0; all test-bearing packages passed. This corroborates C01,
  corrected C02 quantity/unknown behavior, exact Oracle identity selection,
  profitability composition, and C03 field-validation mechanics.
- Security disposition: passing C03 persistence/validation tests do not prove
  trusted actor provenance or authorization and therefore cannot upgrade ★7.
- Live-drive disposition: not performed by this cold inspection reviewer.
  Round-3 execution QA remains responsible for required live/API/UI evidence
  and for folding any mismatch or `could-not-drive` result into the gate.

## Should-meet

| # | Category | Finding | Cited excerpt | Auto-fixable? |
|---|---|---|---|---|
| 8 | Artifact Integrity | Contract and milestone metadata still say `planned`, criterion statuses remain `Pending`, and the QA-owned validation result still truthfully records round-2 Fail. This is stale planning metadata, but this reviewer does not edit QA-owned state. | `validation-contract.md:6`, `:34`, `:49`, `:64`, `:98-100`; `validation-result.md:5`, `:9` | Yes, by the correct owner after gate disposition |
| 9 | Restart And Density | The correction ledger and accepted F-06/F-07 handoffs name status, owner, next action, SHA/evidence, and retained blockers sufficiently for a fresh QA session. | `corrections/correction-task.md:65-72`; `F-06-quantity-cost-semantics/validation.md:137-167`; `F-07-order-specific-tax-provenance/validation.md:119-141` | No |

## Risks

- Real resolved-order tax remains operationally unavailable until an
  owner-approved upstream process supplies exact Oracle `NUNOTA` and
  `SEQUENCIA` per marketplace order item. The implementation correctly keeps
  tax unknown meanwhile; no complete-margin claim is supported.
- Caller-controlled manual-adjustment identity can forge audit attribution and
  is not authorized to installation/order/item scope. This is the decisive
  remaining product/security defect.
- Required round-3 execution QA and live-driving evidence may add failures or a
  missing-live-evidence blocker; it cannot upgrade this cold-review Fail.

## Verdict computation

must_meet_total: 7 | must_meet_pass: 4 ⇒ **Fail** because ★2, ★3, and ★7 fail.

Milestone status: `correction_needed`.

## Recommendation

Do not pass M-06 at this SHA. Preserve the accepted F-06/F-07 behavior. The
next authorized correction must establish the trusted-principal and scoped
authorization boundary for C03 and produce durable exact ran-output artifacts;
the contract’s generic command labels should be made executable and
observable. Keep the absent Oracle upstream mapping explicit and keep tax
unknown until exact `NUNOTA`/`SEQUENCIA` is supplied.

## Next Handoff

Next owner: Milestone Orchestrator, then independent QA Validator. QA must run
the proportional round-3 execution/live gate and alone may update
`validation-result.md`. No implementation or QA-owned artifact was changed by
this review.
