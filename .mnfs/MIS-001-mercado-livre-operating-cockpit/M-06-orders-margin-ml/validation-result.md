# M-06 validation result — 2026-07-10

```yaml
milestone: M-06
status: blocked
validator: QA Validator
validated_at: 2026-07-10
scope: Orders + Margin
```

## Verdict

**BLOCKED — do not mark M-06 passed.** The required idempotence, honest-quality,
and adjustment-audit criteria have supporting evidence below. The mandatory
real resolved-product-link realization scenario remains unproven: a candidate
must be explicitly approved with a truthful audit actor, then re-imported and
calculated against live Oracle inputs.

## Criterion evidence

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
