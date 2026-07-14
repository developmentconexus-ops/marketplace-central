# M-06 validation result

```yaml
milestone: M-06
status: correction_needed
validator: QA Validator
validated_at: 2026-07-11
validation_round: 2
verdict: Fail
reviewed_sha: 5548ae406cb26d0703c111236d703281bb227d3e
tested_head: 2e8c9250303c9e9300055dc8030e9ae7fb62093c
scope: Orders + Margin
```

## Historical round-1 verdict — 2026-07-10

**BLOCKED — do not mark M-06 passed.** The required idempotence, honest-quality,
and adjustment-audit criteria have supporting evidence below. The mandatory
real resolved-product-link realization scenario remains unproven: a candidate
must be explicitly approved with a truthful audit actor, then re-imported and
calculated against live Oracle inputs.

## Historical round-1 criterion evidence

| ID | Result | Evidence | Target type |
| --- | --- | --- | --- |
| M-06-C01 Idempotent Order Ingestion | PASS | Two live Mercado Livre imports (`limit=50`) kept PostgreSQL counts at orders/items/payments `30/30/30` before, after first, and after second import. Provider responses remained 30 orders; no duplicate rows were introduced. Orders service/repository tests also pass. | Live Mercado Livre + real Docker PostgreSQL; tests supplementary |
| M-06-C02 Margin Quality Honesty | PASS | Live pipeline recalculated 60 snapshots: `missing_tax_count=60`, `missing_tax_with_realized_math=0`, `complete_count=0`, `negative_count=0`. Live inputs left missing freight/commission/tax/link as unknown. Profitability Go suite passes, including partial-tax preservation. | Real Docker PostgreSQL + live imported data; tests supplementary |
| M-06-C03 Manual Adjustment Audit | PASS | `TestManualAdjustmentsAppendOnlyReadbackAndConstraints` passed against real Docker PostgreSQL. It proves required audit fields/readback and atomic duplicate idempotency behavior: same key returns the immutable original, while distinct keys create distinct rows. OpenAPI, SDK, and UI require an idempotency key. | Real Docker PostgreSQL; UI/SDK contract tests supplementary |

## Integration and user-interface evidence

- Docker Compose PostgreSQL, backend, and frontend were healthy; backend
  `/healthz` returned HTTP 200.
- Mercado Livre account probe and listings/orders probes returned live data;
  30 orders were imported.
- Live Oracle/Sankhya smoke passed product, stock, current price, as-of
  CUSSEMICM cost, sales history, and tax queries. The observed Godror timezone
  warning did not fail a query and remains an operational follow-up.
- Built-in browser QA at `/orders` passed desktop and mobile rendering with no
  console errors. Evidence: `../../../../.superpowers/sdd/evidence/m06-f03-total-validation-desktop.png`
  and `../../../../.superpowers/sdd/evidence/m06-f03-total-validation-mobile.png`.
- The correction cold review initially found and the implementation corrected:
  partial unknown tax becoming realized math, order-scope cost being applied to
  cost instead of adjustment, and non-idempotent manual adjustments. The final
  independent quality review returned PASS; its local Go run did not itself
  have `MC_DATABASE_URL`, so the cited PostgreSQL evidence is the separate
  real-Docker run.
- A later preflight found that exact-match candidates leaked an internal
  product ID before approval. The adapter and profitability boundary were
  corrected with genuine RED/GREEN evidence and independent SPEC/QUALITY
  approval. Fresh real reimport now records zero non-resolved items with an
  internal product ID and zero known cost/tax values on non-resolved items.
  `MLB4834373620` is `unresolved` with a null product ID; its cost and four tax
  inputs are null with `unresolved_link`. This supersedes the earlier 29-cost
  readback as resolved-link evidence.

## Buyer-PII scope review

PASS for the M-06 orders read model. The review searched the order transport,
domain, OpenAPI, SDK, and Orders UI for buyer, recipient, email, phone,
document, and address fields. The normalized ingestion contract contains no
buyer object or contact/address data. `RawProviderRef` is constructed by
`safeOrderProviderReference` as the provider endpoint path `/orders/{id}`, and
the adapter test explicitly asserts the provider raw reference is not exposed
to the normalized ingestion snapshot. The remaining `shipping_id` is an
operational shipment identifier, not buyer contact data. No buyer PII was
found exposed beyond the stated operational scope.

## Remaining blockers and exact next action

1. Approve an unambiguous generated product-link candidate with explicit
   `actor_type`, `actor_id`, and `actor_name` (no invented actor).
2. Re-import its live Mercado Livre order and run the profitability pipeline.
3. Capture a snapshot with a resolved link and real Oracle CUSSEMICM/tax inputs;
   it may only be `complete` when every required value is known.
4. Re-run the cold gate after the live resolved-link evidence is recorded.

## Fresh Docker revalidation addendum — 2026-07-10

After the Docker data volume was recreated, the current stack re-imported 20
real Mercado Livre orders, imported real margin inputs through the configured
internal reader, and recalculated 40 real PostgreSQL snapshots. The browser
showed 15 incomplete paid orders and 5 not-realized cancelled orders without
rendering unknown contribution or margin as zero. Desktop and 390x844 mobile
browser QA passed with no console errors.

The focused F-03 Go suite, SDK suite, Orders web tests, and web build all
passed in the current Docker stack. The broad `go test ./apps/server_core/...`
command remains a separate cold-gate failure: F-03 profitability passed, but
unrelated OAuth-environment, marketplace-fixture, migration-path, and
inventory-expectation tests fail. No out-of-scope test or production change
was made in this shared worktree.

The fresh database currently has no generated product-link candidate. A
truthful actor approval can only be requested after a candidate exists; a
resolved-link realized-margin scenario therefore remains unproven. This
addendum preserves the verdict: **BLOCKED — do not mark M-06 passed.**

## Candidate follow-up — 2026-07-10

Twenty real Mercado Livre listing identities generated twenty exact-EAN
candidate records after live Oracle smoke passed and the backend was restarted
to recover from a transient startup ping failure. The only two candidates
that overlap the reimported 30 orders belong to cancelled orders. They can
only yield `not_realized`, even after approval, and therefore cannot supply
the mandatory paid resolved-link margin evidence. No candidate was approved
and no actor was invented. The full gate remains blocked by this missing
realized scenario and the separately recorded unrelated full-server test
failures.

## Sources

- `.superpowers/sdd/m06-f03-live-validation-2026-07-10.md`
- `.superpowers/sdd/m06-f03-cold-gate-correction-report.md`
- `.superpowers/sdd/m06-unapproved-link-correction-report.md`
- `docs/superpowers/specs/2026-07-09-m06-f03-order-realization-design.md`
- `docs/superpowers/plans/2026-07-09-m06-f03-order-realization.md`

## Round-2 formal gate — 2026-07-11

- Verdict: **Fail** (formal QA verdict; exact value: `Fail`)
- Reviewer: QA Validator, folding the independent round-2 cold review with execution QA.
- Contract checked: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-contract.md`.
- Reviewed freeze: `5548ae406cb26d0703c111236d703281bb227d3e`.
- Tested HEAD: `2e8c9250303c9e9300055dc8030e9ae7fb62093c`; the only freeze-to-HEAD change was the persisted `milestone-review.md` artifact.
- Execution evidence: `_gate-evidence/round-2/gate-results.md`, `_gate-evidence/round-2/ui/drive-log.txt`, and `_gate-evidence/round-2/ui/flows.json`.

### Crew review fold

| Criterion | Folded result | Blocking reason |
| --- | --- | --- |
| ★1 Criteria coverage | Pass | None. |
| ★2 Evidence honesty | Fail | The cold review found the load-bearing historical PASS rows were not bound to concrete `ran` artifacts; new QA evidence cannot downgrade that independent failure. |
| ★3 Verifiability | Fail | Contract commands remain generic, and required live visual/interactive UI evidence is `could-not-drive`. |
| ★4 Integration/composition | Fail | Quantity-extended revenue composes with unextended CUSSEMICM cost. |
| ★5 Traceability | Pass | Feature-to-criterion and criterion-to-evidence closure remains intact. |
| ★6 Correction integrity | Fail | Correction-attempt and prior-result fields remain blank and no append-only authorized correction ledger closes the history. |
| ★7 Security posture | Fail | The manual-adjustment actor is caller-supplied, only checked for non-empty fields, and not bound to an authenticated/authorized principal. |

`must_meet_pass: 2 / 7`; any must-meet failure computes **Fail**. The fold never upgrades a cold-review failure. The browser `could-not-drive` condition is an additional evidence limitation and a blocker for any future Pass, but concrete defects already establish **Fail**, so the formal verdict is not changed to Blocked.

### Re-run corroboration sample

| Sample | Observed result | Gate effect |
| --- | --- | --- |
| Exact Go orders plus composition suite; targeted C01 tests | Pass | Corroborates C01. |
| Exact Go profitability plus composition suite with `GOCACHE=.gocache`; targeted C02 tests | Pass | Unknown propagation is reproduced, but cannot override the quantity-scoped live defect. |
| Targeted C03 tests, including `TestManualAdjustmentsAppendOnlyReadbackAndConstraints` | Pass | Append-only persistence/constraints are reproduced, but cannot prove trusted actor provenance or authorization. |
| SDK runtime suite | Pass, 35 tests | Supplementary contract corroboration. |
| Orders UI suite | Pass, 13 tests | Supplementary component corroboration; not live visual evidence. |
| Web route/context/proxy tests | Pass | Supplementary route corroboration. |
| Web build | Pass | Build corroboration only. |

No import, calculation, adjustment, approval, POST, or provider call was issued by this finalization pass.

### Live runtime validation

- Read-only GETs returned HTTP 200 for `/healthz`, `/orders`, `/profitability/margin-inputs`, `/profitability/profit-snapshots`, and `/profitability/manual-adjustments`.
- Persisted readback contained 30 unique orders, 30 unique items, 30 unique payments, and 60 snapshots; `missing_tax_with_math=0`.
- Six Candidate A orders are paid, approved, and resolved. The quantity-two and quantity-seven live rows retain cost `91.57` while revenue is `339.98` and `1189.93`; expected quantity-extended costs are `183.14` and `640.99`. This confirms the blocking ★4 quantity-cost defect.
- The quantity-one evidence truthfully closes the old paid/resolved-link evidence gap, but remains incomplete because tax is missing.
- The manual-adjustment actor remains caller-supplied and is validated only as nonempty. No POST was issued. This confirms the blocking ★7 security defect.
- In-app browser discovery found no controllable browser, so UI live drive is `could-not-drive`.
- The finalizer's single permitted reachability check of `http://localhost:5174/orders` returned HTTP 200. This is reachability only and is not visual or interactive evidence.

### Round-2 criteria results

| Criterion | Status | Evidence and blocking disposition |
| --- | --- | --- |
| M-06-C01 Idempotent Order Ingestion | **Pass** | Exact integrated/targeted tests passed and readback contains 30 unique orders/items/payments. |
| M-06-C02 Margin Quality Honesty | **Fail** | Unknown propagation tests and `missing_tax_with_math=0` pass, but quantity-extended revenue composes with unextended `91.57` cost for quantity two/seven rows. |
| M-06-C03 Manual Adjustment Audit | **Fail** | Append-only persistence tests pass, but caller-controlled actor identity is neither derived from a verified principal nor shown authorized for installation/order/item scope. |

### Exact correction scope

1. First reconcile `correction_attempts`, `last_validation_result`, the append-only correction log, and authority to proceed under `max_correction_attempts: 2`; do not rewrite round-1 or round-2 history.
2. If authorized, apply a bounded correction for quantity extension, with quantity-one/two/seven plus nil-cost tests and explicit unit-versus-line CUSSEMICM, fee, and tax semantics.
3. Obtain the architecture/owner decision for the trusted-principal boundary, then apply a bounded authorization correction that derives actor identity from a verified principal and authorizes installation/order/item scope. Update OpenAPI and `sdk-runtime` together if the request contract changes.
4. Bind every load-bearing claim to the exact command or interaction output artifact marked `Evidence type: ran`, and replace generic contract command labels with exact executable commands and observables.
5. Refresh live evidence after authorized corrections, then request a complete new fixed-SHA cold review and execution QA, including in-app visual/interactive `/orders` validation.

No correction was implemented, no prior approval was reissued, and QA performed no provider write or runtime-data mutation. Next owner: Milestone Orchestrator for correction-history/cap reconciliation, owner decisions, bounded correction authorization, and a later fresh fixed-SHA gate.


## Proportional fixed-SHA QA — 2026-07-13

- Frozen SHA: `1eb8831fb1d0d1b84f4d1325978bbc4f76c9ed0f`
- Verdict: **failed**
- QA mode: proportional fixed-SHA review

### Criterion dispositions

| Criterion | Disposition | Evidence / reason |
| --- | --- | --- |
| M-06-C01 | not independently re-verified | F-11 validation evidence may support linkage-read contract safety, but QA was restricted to the registered commands and that evidence does not waive requirements outside F-11's scope. |
| M-06-C02 | not independently re-verified | The supplied F-07 and F-15 evidence remains scoped feature evidence; no registered command established this milestone criterion. |
| M-06-C03 | failed / deliberately deferred | No trusted-principal and authorization boundary is proven. Production writes must not be described as authenticated. This explicitly prevents an M-06 pass. |

### Registered command evidence

1. `git rev-parse HEAD` exited `0` and returned `1eb8831fb1d0d1b84f4d1325978bbc4f76c9ed0f`.
2. `git diff --check 0cfa801b7f9cfe57c0cd81f7c953e81a8a706cbf..1eb8831fb1d0d1b84f4d1325978bbc4f76c9ed0f` exited `0` with no output.
3. `GOCACHE=.gocache go test ./internal/modules/internal_read/domain ./internal/modules/internal_read/application ./internal/modules/internal_read/adapters/oracle -count=1` exited `1` in the required Windows PowerShell execution environment. Go reported that the three paths do not exist and that `GOCACHE` must be absolute.

### Blockers and next

M-06-C03 is the milestone blocker: establish and prove a trusted-principal plus authorization boundary before representing production writes as authenticated, then rerun proportional QA at a new frozen SHA. Correct the registered Go-test invocation/package routes as part of the validation contract before using it as passing evidence.

