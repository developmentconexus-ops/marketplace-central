# Milestone Validation Contract

```yaml
id: M-05
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

M-05-visao-geral-pedidos-sync-central

## QA Level

QA-2 — integration lane for aggregates; vitest + browser walkthrough for workspaces.

## Required Outcome

Visão geral at `/`, Pedidos, Integrações & Sync live over new aggregate/orders/runs read APIs; legacy Dashboard + Integrations replaced.

## Criteria

## Criterion: Summary degrades honestly
ID: M05-C01
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: integration test with linkage sub-read stubbed to fail
- Expected: 200 with `"pending_links": null` (SQL/JSON null, not 0) and `"degraded": ["linkage"]`; all-sources-fail case → 200 all-null counters + full degraded list, never 500
- Actual:
- Artifact: `F-01-aggregate-sync-endpoints/validation.md`
Blocking failure: 0 substituted for failed counter or 500 on partial failure
Blocking failure observed: No
Owner: QA Validator

## Criterion: Orders read contract
ID: M05-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: integration tests over seeded orders (0027/0033 tables)
- Expected: cursor list newest-first; `filter.status`/fulfillment/date filters honored; `GET /orders/{unknown}` → 404 `order_not_found`; malformed date → 400 `invalid_filter`; canonical fields only (no raw provider payload in response)
- Actual:
- Artifact: `F-01.../validation.md`
Blocking failure: provider payload leak or filter row failing
Blocking failure observed: No
Owner: QA Validator

## Criterion: Sync runs listing
ID: M05-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: integration test over seeded `integration_operation_runs` incl. one running and one run seeded with `started_at` older than 90 days
- Expected: cursor list filterable by module/status, sorted `started_at DESC` (newest first, asserted on seeded rows); running row has `"status":"running"` and `"finished_at": null`; the >90d-old seeded run is absent from the listing (90d window per mission "sync central 90d view")
- Actual:
- Artifact: `F-01.../validation.md`
Blocking failure: running run misreported as finished
Blocking failure observed: No
Owner: QA Validator

## Criterion: Module-boundary composition (Q6)
ID: M05-C04
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: code review + grep of summary implementation
- Expected: dashboard summary calls other modules' application services only; zero SQL touching another module's tables from the dashboard transport; OpenAPI+SDK same-commit for all new endpoints
- Actual:
- Artifact: milestone `validation-result.md` review note + diff excerpts
Blocking failure: cross-module SQL join or split API/SDK commit
Blocking failure observed: No
Owner: QA Validator

## Criterion: Visão geral cards + deep links
ID: M05-C05
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser on seeded data; vitest deep-link map test
- Expected: `/` renders 1e cards with real counts; null counter card shows "—" + tooltip "fonte indisponível: {source}"; clicking sync-error pendência lands `/anuncios?tab=pendencia&filter.exception=sync_error` pre-filtered; degraded banner lists sources while healthy cards stay live
- Actual:
- Artifact: `F-02-visao-geral/validation.md` screenshots
Drive (UI — agent-browser; UI criteria only; omit for non-UI):
- Fixture: seeded counters + one degraded source, installation inst_1
- Steps:
  - open http://localhost:5174/?installation=inst_1
  - assert text "Pendências"
  - click pendência row "sincronização"
  - assert url ~ /anuncios\?.*filter.exception=sync_error
- Expected: pre-filtered workspace opens
Blocking failure: 0 rendered for null counter or dead deep link
Blocking failure observed: No
Owner: QA Validator

## Criterion: Dashboard direct-fetch retired
ID: M05-C06
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: grep legacy Dashboard component + network trace of `/`
- Expected: legacy Dashboard deleted; `/` issues only `getDashboardSummary` (+context) through TanStack Query; zero direct fetches
- Actual:
- Artifact: `F-02.../validation.md`
Blocking failure: legacy component alive or direct fetch on `/`
Blocking failure observed: No
Owner: QA Validator

## Criterion: Pedidos read-only workspace
ID: M05-C07
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: vitest + browser; grep of page imports
- Expected: 1j table with status/fulfillment (blu) tags per fixed label map; unknown NF renders "—" with "NF: ERP" hint (not "pendente"); URL round-trip on filters; zero mutation SDK methods imported; visiting `/orders?installation=X&filter.status=paid` lands `/pedidos?installation=X&filter.status=paid` (query preserved, no 404)
- Actual:
- Artifact: `F-03-pedidos-workspace/validation.md`
Blocking failure: any mutation control present or NF default fabricated
Blocking failure observed: No
Owner: QA Validator

## Criterion: Integrações & Sync central (Q4)
ID: M05-C08
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser on seeded runs + protocolos; vitest poll test
- Expected: run history with running row tag "sincronizando" (blu) polling at 30s to terminal; protocolos section lists mutations with type/state, row click lands `/protocolos/:id`; error detail behind "▸ técnico" with pt-BR translation visible; installation card with invalid token shows err tag + "Reconectar"; clicking "Reconectar" issues the existing auth-session start request (asserted in the network trace) and navigates to the provider authorize URL — no new auth surface
- Actual:
- Artifact: `F-04-integracoes-sync/validation.md` screenshots
Blocking failure: raw error outside técnico, or "Reconectar" click fails to issue the auth-session start request
Blocking failure observed: No
Owner: QA Validator

## Criterion: No token exposure (Q2)
ID: M05-C09
Level: Milestone
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: grep rendered output + API responses for token material during Integrações walkthrough
- Expected: zero access/refresh token values in any DOM, API response, or evidence artifact; OAuth flow server-side untouched (git diff scope proof)
- Actual:
- Artifact: milestone `validation-result.md` audit section
Blocking failure: any token value surfaced
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

Feature proofs `F-0*/validation.md`; rollup `validation-result.md` with fixed SHA + dual-gate records.

## Blocking Failures

Any criterion blocking failure; any counter fabricating 0 for unknown (ADR-17) blocks regardless.

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: none

## Handoff

- Current status: planned.
- Next owner: QA Validator (after F-01..F-04 accepted).
- Next action: execute criteria at fixed SHA.
- Required files/evidence: as above.
- Blockers or open decisions: none.
