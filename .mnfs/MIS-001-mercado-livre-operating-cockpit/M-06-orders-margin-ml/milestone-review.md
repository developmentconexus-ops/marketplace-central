# Milestone Review — M-06 assisted Sankhya linkage

reviewer: independent fixed-SHA reviewer

reviewed_sha: `ef4b08c78d30a5e2269e79b051a432c9dc12b58d`

accepted_resumed_base: `6fe6e2a056c0397d7c5ad45555581ab1175c7cef`

reviewed_range: `6fe6e2a056c0397d7c5ad45555581ab1175c7cef..ef4b08c78d30a5e2269e79b051a432c9dc12b58d`

verdict: **Fail — correction required before QA can consider a milestone pass**

## Summary

The frozen range establishes the intended tenant-scoped architecture: opaque
MPC order-line identity, an append-only/idempotent Postgres confirmation
ledger, explicit TOP 313 operator selection, exact TOP 306 descendants through
`TGFVAR`, fail-closed runtime configuration, and matching OpenAPI/SDK routes.
No architecture, tenant-boundary, provider-payload, or public-contract parity
defect was found in those seams.

The milestone cannot pass at this SHA. A partial lineage whose currently known
TOP 306 descendants have all four known tax components is imported as
`partial`, but snapshot aggregation ignores that quality and can calculate a
`complete` margin. This violates F14-AC04 and M-06-C02. The frozen Feature
validation overstates coverage of that exact case. Independently, the existing
trusted-principal/manual-adjustment requirement remains deliberately deferred
and failing; this review does not waive it.

## Findings

| Severity | Classification | Finding | Evidence | Required disposition |
| --- | --- | --- | --- | --- |
| **High** | **contract** | Partial or otherwise incomplete exact-line tax facts can become a complete profit snapshot. `aggregateExactTaxInput` correctly emits a non-nil amount with `InputQualityPartial` for partial lineage or incomplete facts, but `applyInput` records missing tax only when `Amount == nil`; `finalizeSnapshot` then sees non-nil tax and no missing flag and may calculate contribution/margin with `ProfitSnapshotComplete`. A counterexample is partial lineage with one known descendant and all four known tax components. | `apps/server_core/internal/modules/profitability/application/service.go:342-376`, `:846-864`, `:886-958`; F-14 spec `:55-66`, `:103-107` | Propagate partial/incomplete tax quality into item and order snapshot incompleteness while retaining known amounts; add an end-to-end import-to-snapshot test for all-known components under partial lineage and for complete lineage with incomplete source quality. |
| **High (known deferred constraint)** | **ownership** | Production write attribution/authorization is not established. Manual adjustments still accept caller-supplied actor data, and assisted confirmation explicitly records `operator_supplied_unverified`. The correction ledger correctly retains M-06-C03 as pending/failing; owner deferral is not a pass or security boundary. | `corrections/correction-task.md:19-20`, `:31-35`, `:71-74`; `F-13-assisted-linkage-api-runtime/validation.md:248-249`; `apps/server_core/internal/modules/orders/transport/sankhya_linkage_handler.go:101-125`, `:222-229` | Keep production writes disabled or otherwise deployment-restricted until an owner-approved trusted-principal and tenant/installation/order/item authorization boundary is implemented and evidenced. QA must not pass C03 meanwhile. |
| **Medium** | **verification** | F-14 validation claims that known partial-lineage sums keep item and order margin incomplete, but the focused test only asserts partial quality on imported inputs. The snapshot test uses a nil PIS amount, so it does not exercise the all-components-known partial-lineage counterexample above. | `F-14-profitability-sankhya-lineage/validation.md:50-52`, `:72`, `:109`; `apps/server_core/internal/modules/profitability/application/service_test.go:309-336`, `:921-946` | Replace the overclaim with a real regression test and fresh exact command evidence after correction. |
| **Medium** | **verification** | The QA-owned milestone result still describes the older round-2 SHA, and the validation contract still uses generic suite labels with the rollup itself as the artifact. Feature records provide useful deterministic evidence, but there is no fresh QA-owned fixed-SHA result for `ef4b08c…`. | `validation-contract.md:36-40`, `:51-55`, `:66-70`; `validation-result.md:1-10`; F-09 through F-14 `validation.md` files | QA must run exact registered commands against the frozen SHA, capture durable outputs, and issue the new fixed-SHA result without relabelling old evidence. |

## Reviewed controls without additional findings

- **Architecture and tenant scope:** migration 0033 carries tenant and
  installation scope through line identity, event, line, foreign-key,
  idempotency, and origin uniqueness constraints; append-only triggers reject
  update/delete (`0033_orders_sankhya_linkage.sql:17-18`, `:37-110`,
  `:122-128`). Repositories additionally enforce the configured tenant.
- **Stable identity and idempotency:** refresh reconciliation excludes mutable
  quantity, price, and position from identity; ambiguous/legacy lines cannot
  confirm. Confirmation is atomic, line-owned, duplicate protected, and a
  same-key retry succeeds only for semantically equal intent.
- **Exact Sankhya lineage:** candidate SQL binds the exact external key and TOP
  313; descendant SQL binds origin `NUNOTA`/`SEQUENCIA`, joins `TGFVAR`, and
  requires TOP 306 (`sankhya_linkage_reader.go:323-326`, `:349-358`). Nullable
  quantities and one-to-many descendants remain explicit.
- **Runtime activation and public contract:** all assisted-linkage settings are
  explicit and bounded; malformed/disabled configuration leaves the routes
  registered but unavailable. Requests cannot supply tenant, event ID,
  configuration revision, evidence reference, actor type, or external key.
  OpenAPI and `sdk-runtime` expose the same three routes and DTO shapes.
- **Committed-range hygiene:** both SHAs resolve, the accepted resumed base is
  an ancestor of the frozen SHA, the range contains seven intentional commits,
  and `git diff --check 6fe6e2a…..ef4b08c…` returned no errors.

## QA constraints

1. Do not pass M-06-C02 until the partial-quality snapshot defect is corrected
   and an all-known-components/partial-lineage regression proves item and order
   contribution/margin remain unrealized.
2. Do not pass M-06-C03 or describe production writes as authenticated while
   actor identity remains caller supplied or explicitly unverified.
3. No live Oracle deployment fact is proved by the frozen evidence. Before
   activation, obtain values-free evidence for the configured TGFCAB field and
   capacity, exact/no-transform behavior, a durable uniqueness mechanism plus
   duplicate probe, SELECT-only MPC permissions, and representative TOP
   313-to-306 `TGFVAR` lineage including partial/one-to-many behavior. A
   syntactically valid environment alone is insufficient.
4. Fake/unit evidence proves deterministic contract behavior only. QA must
   keep metadata, uniqueness, Oracle availability, customer configuration, and
   real lineage states `unknown` unless the live lane actually proves them.

## Blockers and next

Blockers: the High contract defect above; the retained owner-deferred C03
trusted-principal boundary for any full milestone pass. Live Oracle facts are
an activation/QA constraint, not evidence supplied by this review.

Next owner: Milestone Orchestrator for one bounded profitability correction,
then a new frozen SHA, independent review, and proportional QA. Only QA may
write/pass `validation-result.md`.

Review method: committed-object inspection only. No source, test, contract,
Feature artifact, runtime data, live system, dependency, Git ref/history, or
QA-owned result was modified by this reviewer.
