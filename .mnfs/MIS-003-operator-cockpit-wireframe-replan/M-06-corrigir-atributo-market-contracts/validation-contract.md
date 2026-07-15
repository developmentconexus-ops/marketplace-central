# Milestone Validation Contract

```yaml
id: M-06
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: milestone
```

## Milestone ID

M-06-corrigir-atributo-market-contracts

## QA Level

QA-2 — integration lane (stub capability/adapter); browser walkthrough for the flow.

## Required Outcome

Corrigir-atributo mini-flow live end-to-end via envelope; market module contract-only, honestly empty.

## Criteria

## Criterion: Category attributes endpoint + cache
ID: M06-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: integration test against stubbed capability; repeated call within TTL
- Expected: `GET /listings/categories/{id}/attributes` returns canonical DTO (id, name pt, type, required, allowed values, constraints); second call within 24h serves from L2 class `category_meta` with zero capability invocations (spy assert); unknown category → 404 `category_not_found`
- Actual:
- Artifact: `F-01-corrigir-atributo-flow/validation.md`
Blocking failure: raw ML payload in response or cache miss on fresh entry
Blocking failure observed: No
Owner: QA Validator

## Criterion: Schema-driven form validation
ID: M06-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: vitest renderer tests per attribute type
- Expected: enum → select restricted to allowed values (free text impossible); number → input enforcing numeric constraints with constraint copy on violation; boolean → toggle; required empty → blocked pre-preview
- Actual:
- Artifact: `F-01.../validation.md`
Blocking failure: constraint bypass reaching preview
Blocking failure observed: No
Owner: QA Validator

## Criterion: Attribute fix flows through envelope
ID: M06-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser stub-lane: pendência row → select attribute → edit → M-03 modal → terminal
- Expected: `listing_edit` protocolo with attribute intent; terminal fires `invalidateAfterMutation('listing_edit')`; pendência flag cleared after listings refresh; provider rejection surfaces `provider_validation` with detail behind "▸ técnico"
- Actual:
- Artifact: `F-01.../validation.md` screenshots
Drive (UI — agent-browser; UI criteria only; omit for non-UI):
- Fixture: seeded listing with missing required attribute, stub adapter, installation inst_1
- Steps:
  - open http://localhost:5174/anuncios?installation=inst_1&tab=pendencia
  - click flagged row action "Corrigir atributo"
  - fill attribute field "<valid value>"
  - click confirm checkbox, click "Confirmar"
  - assert text "MP-"
- Expected: protocolo result view with applied item
Blocking failure: write path outside envelope or client masking provider rejection
Blocking failure observed: No
Owner: QA Validator

## Criterion: Market empty-state honesty
ID: M06-C04
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: integration test on virgin DB
- Expected: `GET /market/observations?installation_id=inst_1&listing_ids=a,b` (IC-04 params) → 200, every item `"evidence_state": "no_price_evidence"` with null signals; never 404, never fabricated values; malformed id in list → item-level no_price_evidence, request still 200
- Actual:
- Artifact: `F-02-market-contract-module/validation.md`
Blocking failure: error or invented value for unknown id
Blocking failure observed: No
Owner: QA Validator

## Criterion: Six-signal separation + provenance
ID: M06-C05
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: round-trip test via test-double CollectorPort writing two signal kinds for one listing; constraint test omitting source
- Expected: query returns two SEPARATE signal entries (no merged/derived price); write without `source`/`captured_at` rejected by NOT NULL constraint
- Actual:
- Artifact: `F-02.../validation.md`
Blocking failure: signal merge or provenance-less row accepted
Blocking failure observed: No
Owner: QA Validator

## Criterion: Contract-only enforced
ID: M06-C06
Level: Milestone
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: grep sweep: CollectorPort implementations, scraping libs in go.mod, seed inserts, ML market/search API calls
- Expected: test-double confined to `_test` packages; zero production CollectorPort implementations; zero scraping dependencies; zero market seed data outside `_test`; zero ML search/market endpoints called (ML ToS 7.6 / G1 FAILED posture)
- Actual:
- Artifact: `F-02.../validation.md` grep transcripts
Blocking failure: any production collector or scraped/fabricated data path
Blocking failure observed: No
Owner: QA Validator

## Criterion: id-cap and API/SDK discipline
ID: M06-C07
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: integration test with 201 ids; git diff of endpoint commits
- Expected: 422 `{"error":{"code":"too_many_ids"}}` at 201 ids, 200 at 200 ids; OpenAPI + sdk-runtime same commit for all new M-06 endpoints
- Actual:
- Artifact: `F-02.../validation.md` + diff excerpts
Blocking failure: cap unenforced or split commit
Blocking failure observed: No
Owner: QA Validator

## Criterion: Governance lanes green at milestone SHA
ID: M06-C08
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: full governance lane suite at fixed SHA (`GOCACHE` absolute)
- Expected: all lanes green incl. module-layering checks for both new code paths
- Actual:
- Artifact: milestone `validation-result.md` lane logs
Blocking failure: any red lane
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

Feature proofs `F-0*/validation.md`; rollup `validation-result.md` with fixed SHA + dual-gate records.

## Blocking Failures

Any criterion blocking failure; market data fabrication of any kind blocks the mission (escalate to MIS-003-C07).

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: none

## Handoff

- Current status: planned.
- Next owner: QA Validator (after F-01/F-02 accepted).
- Next action: execute criteria at fixed SHA.
- Required files/evidence: as above.
- Blockers or open decisions: none.
