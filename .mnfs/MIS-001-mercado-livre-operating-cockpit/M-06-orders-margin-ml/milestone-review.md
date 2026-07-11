# M-06 cold milestone review — 2026-07-10

## Independent review result

**BLOCKED, not failed.** The final independent quality reviewer found no
actionable code-quality defect in the corrected F-03 profitability path. It
confirmed safe partial-tax propagation, correct cost-adjustment scope,
idempotent adjustment persistence, required OpenAPI/SDK request keys, the UI
retry fingerprint lifecycle, and the profitability/orders boundary.

This is not an approval of M-06 because the live, approved resolved-link
scenario required by the F-03 design and plan has not occurred. A targeted
buyer-PII scope review has now passed: normalized orders intentionally exclude
buyer/contact/address data and retain only the operational shipping identifier
and safe provider endpoint reference.

## Review trail

| Review | Result | Disposition |
| --- | --- | --- |
| Initial cold spec review | BLOCKED | Required result artifact and live resolved-link scenario absent. |
| Initial cold quality review | FAIL | Found partial-tax unknown-to-realized risk, wrong order-scope cost mapping, and non-idempotent adjustments. |
| Correction review | CONDITIONAL PASS | Found a UI pending-key lifecycle defect after a failed request and material form change. |
| Focused UI review | PASS | Fingerprint-scoped retry key regression passed. |
| Final independent quality review | PASS | No actionable code-quality findings. |
| Unapproved-link correction review | PASS | SPEC and QUALITY approved: candidates/non-resolved links expose no product ID; profitability independently requires resolved quality. |
| Buyer-PII scope review | PASS | No buyer/contact/address fields in normalized contract, API, SDK, or Orders UI; raw reference is a safe `/orders/{id}` path. |
| Milestone outcome | BLOCKED | Awaiting explicit link approval/actor and live realized evidence. |

## Evidence integrity

The validation result distinguishes contract/unit evidence from real targets:
Mercado Livre and Oracle are live-provider evidence; PostgreSQL evidence uses
the real Docker database; browser evidence uses the built-in browser. No mock
or compile-only result is being used as proof of an external integration.

See `validation-result.md` for criterion-level commands, observations, and
the exact unblock sequence.

# Milestone Review — M-06 (round 2)

reviewer: milestone-reviewer (cold)

reviewed_sha: `5548ae406cb26d0703c111236d703281bb227d3e`

contract: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-contract.md`

verdict: **Fail**

## Round 2 summary

The quantity-one live evidence closes the previous paid/resolved/realized
evidence gap without claiming a complete margin. M-06 still fails because
quantity-extended revenue composes with unextended product cost,
manual-adjustment actor identity is caller-forgeable, and the
evidence/correction artifacts do not satisfy the milestone rubric.

## Must-meet (★)

| # | Category | Procedure run | Cited excerpt (file:line) | Verdict | Yes-if | Defect locus |
| --- | --- | --- | --- | --- | --- | --- |
| ★1 | Criteria Coverage | Enumerated `M-06-C01` through `M-06-C03`, mapped each to a recorded result, and searched the fixed-SHA milestone scope for `TBD`/`TODO`; none were present. | `validation-contract.md:30`, `:45`, `:60`; `validation-result.md:23`, `:24`, `:25` | **PASS** | — | — |
| ★2 | Evidence Honesty | Read every load-bearing `PASS` and performed the adversarial evidence pass. The quantity-one reconciliation is truthfully limited to an incomplete result, but the three criterion passes are not bound to concrete output artifacts marked `Evidence type: ran`. | `validation-result.md:23`, `:24`, `:25`; `orchestrator-reconciliation-2026-07-11.md:141`, `:145`, `:149`, `:151` | **FAIL** | Mark every load-bearing proof `Evidence type: ran` and bind each claimed pass to its exact command or interaction output artifact. | `validation-result.md:23`, `:24`, `:25` — offending token `PASS` without a cited `ran` artifact. |
| ★3 | Verifiability | Checked each criterion for an exact command or interaction, observable, blocking failure, and concrete evidence path. Expected outcomes and blocking failures are concrete, but all command fields are generic suite labels and their artifacts point back to the rollup. | `validation-contract.md:36`, `:39`, `:51`, `:54`, `:66`, `:69` | **FAIL** | Replace generic command labels with exact executable commands or live interactions, expected observable values, and concrete output artifacts. | `validation-contract.md:36` — `orders service/repository tests`; `:51` — `profitability service tests`; `:66` — `manual adjustment API/repository tests`. |
| ★4 | Integration / Composition | Enumerated quantity, revenue, cost, calculation, route, enum, SDK, ID, and time seams and performed the adversarial integration pass. Revenue is extended by quantity, while Oracle product cost is persisted unchanged and subtracted once. Live quantity-two and quantity-seven rows confirm the mismatch. | `F-02-margin-input-model/spec.md:47`; `apps/server_core/internal/modules/profitability/application/service.go:229`, `:230`, `:429`, `:508`, `:818`, `:820`; `orchestrator-reconciliation-2026-07-11.md:134`, `:139`, `:153`, `:158` | **FAIL** | Define unit versus extended `CUSSEMICM` semantics, compose cost at the same item-line scope as revenue, adjudicate fee/tax quantity semantics, add quantity-one/two/seven tests, and refresh live evidence. | `apps/server_core/internal/modules/profitability/application/service.go:508` — `Amount: cost.Amount`; `:230` calls `mapCostInput` without item quantity. |
| ★5 | Traceability | Traced every accepted feature backward to at least one criterion and every criterion forward to accepted feature or milestone verification. | `F-01-order-ingestion/feature.md:6` and `validation.md:70`; `F-02-margin-input-model/feature.md:6` and `validation.md:47`, `:76`; `F-03-profit-snapshot-calculation/feature.md:6` and `validation.md:71`; `F-04-orders-margin-ui/feature.md:6` and `validation.md:19`; `validation-result.md:21` | **PASS** | — | — |
| ★6 | Correction Integrity | Read the retry policy, prior result, and correction trail. Corrections ran, but attempt count and prior-result fields are blank, no append-only numbered log proves attempts stayed within the cap, and prior fail-to-pass transitions are not tied to new `ran` evidence paths. | `validation-contract.md:87`, `:88`, `:89`; prior review trail above; `validation-result.md:39`, `:45` | **FAIL** | Reconcile an append-only numbered correction ledger against `max_correction_attempts: 2`, bind every prior fail upgrade to new ran evidence, and obtain owner disposition if the cap was exceeded. | `validation-contract.md:87` — blank `correction_attempts:`; `:89` — blank `last_validation_result:`. |
| ★7 | Security Posture | Inspected the PII and manual-write surfaces. Tenant scoping and buyer-PII minimization are evidenced, but the manual-adjustment endpoint accepts caller-supplied actor identity, forwards it unchanged, and validates only that type and ID are non-empty. The reviewed composition wraps the registered mux only with CORS. | `apps/server_core/internal/modules/profitability/transport/http_handler.go:30`, `:95`, `:107`, `:113`, `:116`; `apps/server_core/internal/modules/profitability/application/service.go:265`, `:279`, `:321`; `apps/server_core/internal/composition/root.go:345`, `:404`; `validation-contract.md:59` | **FAIL** | Bind the write to an authenticated and authorized principal or prove an upstream authenticated boundary, derive or verify actor identity server-side, authorize installation/order/item scope, and add a Security criterion with ran evidence. | `apps/server_core/internal/modules/profitability/transport/http_handler.go:116` — `Actor: req.Actor`; `apps/server_core/internal/modules/profitability/application/service.go:279` — only an empty-string check protects audit identity. |

## Required adjudications

### Prior live-evidence blocker

**Closed as an evidence gap.** Candidate A is reconciled as one resolved link
and one approval audit; six affected orders are paid/approved and resolved.
The quantity-one order has real Oracle `CUSSEMICM`, partial known tax,
`realized/incomplete` quality, `missing_tax`, and null contribution/margin.

Evidence:
`orchestrator-reconciliation-2026-07-11.md:88`, `:103`, `:111`, `:124`,
`:128`, `:145`, `:149`, `:151`.

This satisfies the prior requirement for a live paid resolved-link realization
scenario while respecting the rule that completeness may be claimed only when
every required input is known. It does not prove a complete margin or make
M-06 pass.

### Quantity-extension risk

**Confirmed blocking integration defect.** The quantity-one proof is
unaffected by this specific defect. Quantity-two and quantity-seven rows
cannot be accepted as correct contribution evidence: revenue is line-extended,
while product cost remains `91.57` and would be subtracted only once if the
remaining inputs became complete.

## Should-meet

| # | Category | Finding | Cited excerpt | Auto-fixable? |
| --- | --- | --- | --- | --- |
| 8 | Artifact Integrity | Milestone and contract metadata remain `planned`, criteria remain `Pending`, and correction fields are blank despite accepted features and a newer reconciliation checkpoint. | `milestone.md:6`; `validation-contract.md:6`, `:34`, `:49`, `:64`, `:87`; `orchestrator-reconciliation-2026-07-11.md:5` | Yes, after owner reconciliation |
| 9 | Restart And Density | The reconciliation names status, base SHA, checkpoint, evidence boundary, blockers, and next action sufficiently for a fresh session. | `orchestrator-reconciliation-2026-07-11.md:3`, `:7`, `:79`, `:81`, `:147`, `:160` | No |

## Verdict computation

`must_meet_total: 7 | must_meet_pass: 2 => Fail`

Milestone status: `correction_needed`.

## Round 2 recommendation

Do not pass M-06 at this SHA. Reconcile correction-attempt authority, then
apply only authorized bounded corrections for quantity-scoped cost composition
and the manual-adjustment authorization boundary, repair the evidence
contract, and request a new full fixed-SHA gate.
