# F-02-market-contract-module

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-06
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contract: IC-04 `../../research/market-data-interface-contract.md` (entities, CollectorPort signature, endpoints, evidence_state). R-04 (G1 OAuth FAILED; scraping forbidden ML ToS 7.6; 6-signal separation). ADR-14.

## Milestone

M-06 corrigir-atributo-market-contracts. Independent of F-01 (parallel-eligible).

## Brief

Create the `market` Go module contract-only: migrations for `market_observations` + `market_references` per IC-04 DB shape, domain entities with 6-signal separation (our_sale_price/winner_price/competitive_target/catalog_offer_price/catalog_stats/manual as distinct signal kinds — never merged), `CollectorPort` interface exactly per IC-04 Go signature, `GET /market/observations` + `GET /market/references` (max 200 ids; unknown ids → `no_price_evidence` items, not errors), evidence_state computation (observed|insufficient_market|no_price_evidence). NO production adapter; test-double CollectorPort in `_test` packages only proves round-trip. NO seed data ever in production paths.

EARS:
- While no collector has ever run (this mission's permanent state), when observations are queried for any id, the API shall return items with `evidence_state: no_price_evidence` and null signals — 200, never 404, never fabricated values.
- While a test-double collector writes observations, when queried, signals shall round-trip per IC-04 canonical examples with source + captured_at mandatory.
- While >200 ids are requested, when querying, the API shall return 422 `too_many_ids`.
- While two signal kinds exist for one listing, when returned, they shall remain separate signal entries (no merged/derived price — 6-signal separation).

## Inputs

- IC-04 verbatim (shapes, port signature, examples, error matrix), M-01 module skeleton precedent, migration numbering, R-04 signal definitions.

## Expected Output

- Module + migrations + endpoints; OpenAPI + sdk-runtime (`listMarketObservations`, `listMarketReferences`) same commit; `/market` proxy row exists from M-02.
- CollectorPort + test double; contract tests: empty-state, round-trip, id-cap, signal separation.
- Grep-provable absence: no ML market/search API calls, no scraping deps, no seed inserts outside `_test`.

## Constraints

- Contract-only is a hard gate: any production CollectorPort implementation is a defect this mission (G1 failed; future mission wires collector).
- Signals never aggregated into synthetic "market price"; evidence_state is the only derived field.
- source + captured_at NOT NULL on stored rows (provenance mandatory, R-04; IC-04 — null `captured_at` exists only in synthetic `no_price_evidence` read items).
- Tenant scoping; read endpoints only.

## Negative Scenarios

- 201-id request → 422 `too_many_ids`.
- Malformed listing id in list → item-level `no_price_evidence` (per IC-04: unknown ids are not errors) — request still 200.
- Attempt to write observation without source → DB constraint rejection (test asserts).

## Validation Expectations

- `go test` output: contract tests green incl. empty-state and separation tests.
- Integration transcript: `GET /market/observations?installation_id=…&listing_ids=…` (IC-04 params) on virgin DB → 200 all `no_price_evidence`.
- Grep proof: no production CollectorPort impl; test double confined to `_test`; no seed data.
- Diff proof: OpenAPI + sdk-runtime same commit.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: compile context pack; read IC-04 + M-01 module precedent only.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
