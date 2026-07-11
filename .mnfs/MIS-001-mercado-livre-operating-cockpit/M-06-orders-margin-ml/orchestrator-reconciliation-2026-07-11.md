# M-06 Milestone Orchestrator reconciliation — 2026-07-11

```yaml
milestone: M-06
status: evidence_recorded_pending_review_qa
role: Milestone Orchestrator
base_sha: 1465f4245475d9ce8480f6920493a2d5468caa8f
scope: live paid-order and product-link candidate reconciliation
```

This is an operational checkpoint artifact. It does not replace the QA-owned
`validation-result.md` and does not claim that M-06 passed.

The pre-approval checkpoint below is retained as history. The continuation
checkpoint at the end supersedes its operational blocker.

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

## Continuation checkpoint — approved paid resolved-link evidence

Checkpoint ID: `M06-CONT-003`.

### Approval persistence reconciliation

The operator correction reported that Candidate A had already been approved
by the stopped prior session. The approval was **not** reissued.

A single PostgreSQL `REPEATABLE READ`, `READ ONLY` transaction reconciled the
exact tenant, installation, listing, empty variation, candidate, link, and
audit state:

- candidate rows: `1`; exact Candidate A identity, exact-EAN match, and
  internal product `15956`: true;
- link rows: `1`; `resolved`, source Candidate A, and internal product `15956`:
  true;
- audit rows for the listing identity: `1`; exact `approve_candidate`,
  canonical actor/reason, `none -> resolved`, and `NULL -> 15956`: true;
- link `updated_at` and audit `created_at` are both
  `2026-07-11T22:14:45.277233Z`.

The repository implementation persists the link upsert and audit insert in
one database transaction. The readback therefore confirms the reported
atomic transition and also confirms that no duplicate approval audit exists.

### Fresh live order and Oracle reconciliation

The public operational surfaces are installation-scoped, so the refresh was
bounded to the real installation with `limit=50` and its readback was filtered
to Candidate A. No Mercado Livre provider write was invoked.

1. Live Mercado Livre order import returned `30` orders with `0` skipped.
2. Six Candidate A orders are paid, have approved payment state, and now
   persist the empty-variation item as `resolved` to internal product `15956`.
   The latest affected order fetch is `2026-07-11T22:22:01.212695Z`.
3. Profitability input import persisted `270` inputs. Candidate A has `42`
   item inputs across the six orders.
4. The Oracle reader was called with the order effective date and the approved
   internal product. It returned six complete `CUSSEMICM` cost inputs of
   `91.57`; the selected Oracle source row was observed at
   `2026-04-06T03:00:00Z`.
5. The six orders have `24` explicit tax component rows: `22` are null with
   missing quality and `2` are known with complete quality. There are `0`
   null/complete rows and `0` valued/missing rows, so an unknown tax was not
   converted to zero and a known value was not relabelled unknown.

### Persisted resolved-link snapshots

Recalculation persisted `60` snapshots. The six Candidate A item snapshots
are all `realized` and `incomplete`; each retains `missing_tax`, a null
contribution, and a null margin.

| Provider order | Quantity | Revenue | Fee | CUSSEMICM cost | Tax total | Result |
| --- | ---: | ---: | ---: | ---: | ---: | --- |
| `2000016845051264` | 7 | 1189.93 | 22.95 | 91.57 | unknown | realized, incomplete |
| `2000017090094824` | 1 | 169.99 | 22.95 | 91.57 | unknown | realized, incomplete |
| `2000017140396460` | 1 | 169.99 | 22.95 | 91.57 | unknown | realized, incomplete |
| `2000017236767026` | 1 | 169.99 | 22.95 | 91.57 | unknown | realized, incomplete |
| `2000017276984774` | 1 | 169.99 | 22.95 | 91.57 | 21.18 partial | realized, incomplete |
| `2000017336572246` | 2 | 339.98 | 22.95 | 91.57 | unknown | realized, incomplete |

The load-bearing quantity-one proof is order `2000017276984774`: it is a live
paid/approved order with an approved resolved link, real Oracle as-of cost,
and a partial tax result. Oracle supplied COFINS `21.18` and a known PIS zero;
ICMS and IPI remain unknown. The snapshot preserves the known partial sum but
does not calculate contribution or margin while tax is incomplete.

### Review risk and gate boundary

This checkpoint proves the required paid/resolved/realized state, Oracle
`CUSSEMICM` lookup, tax lookup, and honest unknown propagation. It does not
claim a complete margin or M-06 pass.

Live reconciliation also exposed a review risk outside the quantity-one proof:
revenue is extended by order quantity, while the current cost input retains
the per-product Oracle `CUSSEMICM` amount. The quantity-two and quantity-seven
rows above must not be accepted as correct contribution evidence unless the
fixed-SHA review and QA establish the intended cost semantics or require a
scoped correction.

No buyer PII, credential material, or provider payload was recorded. The next
step is to freeze this evidence commit, request one fixed-SHA review, and run
proportional QA. Only QA may change the milestone verdict.
