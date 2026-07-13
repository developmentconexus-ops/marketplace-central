# F-08 Sankhya Linkage Discovery

Date: 2026-07-13

Target: repository plus governed `live-oracle` lane

Safety: read-only; no Oracle/provider/application write, DDL, secret, PII, or
operational row output

## Result

Repository discovery is complete. Live SELECT discovery did not run because
`live-oracle-linkage-discovery` has no executable definition in repository
truth. The Docker runner proves only a generic Oracle smoke-test route and is
not a linkage-discovery substitute. Customer field names, deployed uniqueness,
copy behavior, and live 313-to-306 lineage therefore remain `unknown`.

## Observed Repository Facts

| Fact | Evidence |
| --- | --- |
| MPC orders are keyed by (`tenant_id`, `installation_id`, `provider_order_id`). | migration `0027_orders_marketplace_orders.sql` |
| Items are keyed by positional `line_no`; each import deletes and reinserts item rows. It is not yet immutable cross-system line identity. | `orders/adapters/postgres/order_repo.go` |
| Profitability can attribute tax only with positive exact `DocumentID` and `LineNumber`; absent identity remains missing. | `internal_read/domain/internal_tax.go`, F-07 validation |
| Repository research establishes `TGFCAB.NUNOTA`, `TGFITE` (`NUNOTA`,`SEQUENCIA`), and `TGFVAR` origin/destination lineage with `QTDATENDIDA`. | `research/order-sankhya-linkage-architecture-2026-07-12.md` |
| No repository application or Oracle adapter currently implements ML-order-to-313 linkage. | routed ports plus targeted linkage search |
| OpenAPI/SDK currently expose order import/list and profitability operations, but no Sankhya linkage workflow. | OpenAPI orders/profitability selectors; SDK orders/profitability selectors |

## Observed Runtime Facts

None. No live Oracle SELECT was executed. In particular, this document does
not claim a deployed custom field, constraint, TOP behavior, copy rule, or
observed `TGFVAR` row.

## Supported Contract Facts From Accepted Research

- TOP 313 is the sales-order origin and TOP 306 is the invoice destination for
  the intended workflow.
- Exact lineage predicate is origin (`NUNOTAORIG`,`SEQUENCIAORIG`) to
  destination (`NUNOTA`,`SEQUENCIA`) in `TGFVAR`; `QTDATENDIDA` expresses the
  realized descendant quantity.
- Header/item custom-field names in prior research are conceptual only.
- Partner, buyer, date, product, quantity, price, and value may narrow
  candidates but cannot prove linkage.

## Inferences Adopted By This Feature

- One configured TOP 313 header text field carrying
  `ml:v1:<installation-or-account-id>:<provider-order-id>` is required as an
  external reconciliation guard. Its deployed column name is configuration,
  not source code.
- A TOP 313 item field is **not required for the current assisted-only flow**.
  Explicit operator confirmation of every exact 313
  (`NUNOTA`,`SEQUENCIA`) into an append-only tenant/installation-scoped MPC
  ledger is sufficient proof. An item field can be later hardening, but cannot
  replace ledger proof or `TGFVAR` validation.
- This decision avoids depending on unproved item-field copy behavior and
  supports existing manually created orders. It also handles duplicate ML item
  attributes by assigning immutable `mpc_line_id` rather than treating
  `line_no` or product attributes as identity.

## Exact Read Predicates For The Next Worker

These are contract predicates, not evidence that a live query ran:

1. Header validation: `TGFCAB.NUNOTA = :origin_nunota`, TOP effective at the
   document operation timestamp equals 313, and the safely quoted configured
   header column equals `:provider_order_key`.
2. Origin line validation: `TGFITE.NUNOTA = :origin_nunota` and
   `TGFITE.SEQUENCIA = :origin_sequencia`; optional product/quantity values are
   consistency checks only.
3. Descendant lineage: `TGFVAR.NUNOTAORIG = :origin_nunota` and
   `TGFVAR.SEQUENCIAORIG = :origin_sequencia`; join destination `TGFCAB` and
   require its effective TOP to equal 306; return destination
   (`NUNOTA`,`SEQUENCIA`) plus nullable/explicit `QTDATENDIDA`.

TOP must be resolved with the repository's effective-operation semantics, not
by assuming that a mutable header code alone proves historical TOP. All binds
are exact; no product/date aggregate is permitted.

## Unknowns And Activation Consequence

| Unknown | Consequence |
| --- | --- |
| Deployed TOP 313 header field name and metadata | Linkage feature disabled until admin config validates it. |
| Whether supported customization can enforce non-null-field uniqueness | Confirmation disabled until a supported unique mechanism is attested and probed. |
| Whether 313 metadata copies to 306 | No dependency; 306 identity is accepted only through `TGFVAR`. |
| Live 313/306 TOP effective semantics and representative lineage cardinality | Runtime activation requires bounded aggregate/metadata SELECT evidence. |
| Cancellation/devolution TOP policy | Reversal remains an explicit unsupported/unknown state; mappings are never deleted. |
| Sanctioned field-entry surface and administrator owner | Deployment cannot enable confirmation until named. |

## Candidate Versus Proof

Candidate output is bounded and non-authoritative. Proof requires: tenant and
installation match; immutable provider order/MPC line identity; configured
header key exact match; explicit actor/reason/time; exact 313 header and every
selected line; no conflicting active ledger row; and an idempotent append. A
306 line is proof only when an exact `TGFVAR` descendant of a confirmed 313
line. Ambiguous, partial, or missing evidence remains `unlinked`/`unknown`, and
tax remains missing.
