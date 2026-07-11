# M-06 Milestone Orchestrator reconciliation — 2026-07-11

```yaml
milestone: M-06
status: blocked_pending_operator_approval
role: Milestone Orchestrator
base_sha: 775317c01505f75f08faf2ab176a5716a14dbe64
scope: live paid-order and product-link candidate reconciliation
```

This is an operational checkpoint artifact. It does not replace the QA-owned
`validation-result.md` and does not claim that M-06 passed.

## Safety boundary

- The backend and PostgreSQL targets were healthy.
- Mercado Livre operations were reads only: orders and listing snapshots were
  re-imported into MPC's local database.
- Candidate generation read the configured live Oracle product contract and
  wrote only MPC candidate rows.
- No product-link candidate was approved, no workflow link or approval audit
  was created, no actor was invented, and no provider write was issued.
- No buyer PII, credential material, or provider payload was recorded.

## Fresh live evidence

Target installation:
`inst-mercado_livre-d373dc64-577f-4950-b11b-b90244f30cb2` under
`tenant_default`.

1. A fresh Mercado Livre order import returned and idempotently persisted 30
   orders: 24 `paid` orders with `approved` payments and 6 `cancelled` orders
   with `refunded` payments. The latest order fetch was
   `2026-07-11T21:59:50.164663Z`.
2. A fresh listing import returned 34 live listing snapshots.
3. Candidate generation produced 36 MPC candidate rows. Candidate source
   snapshot times and updates were between `2026-07-11T22:00:51Z` and
   `2026-07-11T22:00:55Z`.
4. Tenant- and installation-scoped PostgreSQL reconciliation found four
   single-candidate, single-product, `exact_ean` identities overlapping 16
   paid/approved orders:

| Provider item | Candidate internal product | Paid/approved orders |
| --- | ---: | ---: |
| `MLB4735328201` | `15956` | 6 |
| `MLB6896003832` | `42194` | 6 |
| `MLB4735364085` | `41912` | 3 |
| `MLB6896003262` | `39587` | 1 |

5. The same scoped readback found `0` persisted product links and `0` product
   link audit entries. The order rows therefore remain `missing` or otherwise
   non-resolved; candidate generation is not approval.

## Exact operational blocker

M-06 now has paid live orders that can pair with unambiguous candidates, but
it has no truthful human approval. An actual authorized operator must select
one of the exact-EAN candidates and provide the real `actor_type`, `actor_id`,
`actor_name`, approval reason, and explicit approval action. The Milestone
Orchestrator cannot infer or invent those facts.

Until that approval exists, the required resolved-link realized-margin proof
cannot be produced and M-06 remains blocked.

## Exact continuation after real approval

1. Verify the persisted link and approval audit for the selected candidate.
2. Re-import the affected live Mercado Livre orders.
3. Import profitability inputs through the live Oracle reader and capture
   as-of `CUSSEMICM` plus tax quality without converting unknowns to zero.
4. Recalculate and record the resolved-link realized snapshot; mark it
   complete only if every required input is actually known.
5. Freeze the resulting SHA, request one fixed-SHA review, and run
   proportional QA. Only QA may update/pass `validation-result.md`.
