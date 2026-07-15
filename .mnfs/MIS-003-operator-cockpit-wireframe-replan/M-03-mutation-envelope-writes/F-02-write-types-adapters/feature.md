# F-02-write-types-adapters

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contracts: IC-03 (per-type intents, failure taxonomy), governance `contracts/governance/execution-lanes.json` 7 provider-write gates (actor, idempotency, execute, resolved-link, policy, source-timestamp, before-after-audit). ADR-16 (SELLER_SKU=CODPROD invariant), ADR-008 link states.

## Milestone

M-03 mutation-envelope-writes. Depends on F-01 (WriterPort seam exists).

## Brief

Implement the six write types behind WriterPort: `price_update` (new PriceWriter in connectors ML adapter), `stock_correct` (wire existing StockWriter live + fold StockActionService flows into envelope), `link_apply` (resolve link via product_links module API), `listing_pause`, `listing_resync` (triggers listings refresh for affected rows), `listing_edit` (title/attributes; SKU field enforced per ADR-16). `listing_create` intent validates then returns `type_not_enabled` (contract-only this mission). Every write passes the 7 gates; gate failures map to IC-03 failure codes (e.g. unresolved link → `link_unresolved`, missing policy → item failed, never silent default).

EARS:
- While an item's listing has link_state ≠ resolved, when `price_update`/`stock_correct` applies, the item shall fail `link_unresolved` without calling the provider.
- While provider returns 429, when writing, the item shall fail `provider_rate_limited` retryable=true and the poller shall back off the chunk.
- While a listing was paused remotely since preview, when writing, the item shall fail `listing_paused_remote` (no forced write).
- While `listing_edit` changes SELLER_SKU to a value ≠ linked CODPROD, when validating, the intent shall be rejected pre-apply (ADR-16 invariant).
- While source data (price/stock) is older than the intent's declared source_timestamp tolerance, when applying, the item shall fail `stale_source`.

## Inputs

- IC-03 per-type intent schemas, execution-lanes.json gate definitions, existing StockWriter (`connectors/adapters/mercado_livre/`), StockActionService + `inventory_stock_actions` flows, product_links resolve API, pricing policy module read API, F-01 WriterPort signature.

## Expected Output

- ML capability adapter writers per type; raw provider payloads stay in adapter.
- Stub adapter (test double) implementing same port with programmable per-type outcomes — used by integration lane.
- StockActionService fold: existing stock-action entry points create envelope protocols; old table read path preserved for history; regression tests for existing flows green.
- Gate enforcement in application layer (not adapter): each gate has a dedicated check + negative test.
- Failure mapping table ML error → IC-03 code, unit-tested.

## Constraints

- No endpoint changes (F-03); no UI.
- Provider payloads never leave adapters; before/after audit stores canonical values.
- No live ML calls in tests — stub only; live lane is milestone-QA, operator-authorized.
- Existing StockWriter behavior contracts preserved (it was implemented, never live-run — treat as code precedent, not proven).

## Negative Scenarios

- All seven gates: one negative test each (missing actor `actor_required`, duplicate idempotency key, execute=false, unresolved link `link_unresolved`, missing policy `policy_missing`, missing source_timestamp `source_time_unavailable`, audit write failure aborts item); ADR-16 rejection → `sku_invariant_violation`.
- ML 401 mid-chunk → `provider_auth`, remaining items in chunk still attempted (auth failure not assumed global).
- `listing_create` intent → 422 `type_not_enabled` at validation, no protocol items created.

## Validation Expectations

- `go test` output: gate matrix (7 negatives), failure-mapping table, ADR-16 rejection, per-type apply tests against stub.
- StockActionService regression transcript: pre-existing flow produces envelope protocol + legacy history intact.
- Grep proof: no `mercadolibre.com` / provider DTO imports outside `adapters/`.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (after F-01 accepted).
- Next action: compile context pack; read IC-03 + execution-lanes.json + named adapter/service paths.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
