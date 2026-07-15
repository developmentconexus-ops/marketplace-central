# F-04-integracoes-sync

```yaml
id: F-04
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contracts: IC-05 (route `/integracoes`, sync ns 30s), IC-03 (`listMutations`), R-01 screen 1k, F-01 sync-runs API, R-02 (legacy Integrations page facts).

## Milestone

M-05. Depends on F-01 + F-03 (route seam order). Last M-05 feature.

## Brief

Build Integrações & Sync (wireframe 1k) at `/integracoes`, replacing legacy Integrations page: installation cards (conta ML, health/status, connect/reconnect via existing OAuth flow endpoints — parity from legacy page), sync runs history table (`listSyncRuns`: module, status, started/finished, counts, error detail behind "▸ técnico"), protocolos list section (`listMutations` filterable by type/state, each row → `/protocolos/:id`), per-module manual sync triggers (existing endpoints + listings `refreshListings`), links to `/classifications` + `/marketplaces` config surfaces (IC-05 secondary placement).

EARS:
- While a run is `running`, when history renders, the row shall show blu tag "sincronizando" and poll (sync ns 30s) until terminal.
- While an installation token is invalid, when the card renders, health shall show err tag with reconnect action (existing OAuth flow — parity, no new auth surface).
- While protocolos list renders, when a row is clicked, navigation shall land `/protocolos/:id` (M-03 page).
- While `/integrations?…` legacy URL is visited, redirect shall land `/integracoes?…`.

## Inputs

- R-01 §1k inventory, R-02 legacy Integrations facts (OAuth/parity checklist), F-01 runs API, IC-03 listMutations, IC-05 keys/components, listings refresh op (IC-02).

## Expected Output

- `/integracoes` page; legacy Integrations deleted; OAuth connect/reconnect parity preserved.
- Component tests: running-poll, health states, protocolo navigation, redirect.

## Constraints

- OAuth flow untouched server-side; UI parity only (auth surface is ★7-sensitive — no new token handling, no token display).
- Manual sync triggers reuse existing endpoints only; no new sync orchestration.
- pt-BR; sync_state labels per IC-02 map.

## Negative Scenarios

- Runs API error → history section ErrorState; installation cards independent (separate queries).
- Trigger sync while running → 409 attach behavior (same as Anúncios refresh pattern).
- Zero installations → EmptyState with connect CTA (entry point of the M-02 empty-context flow).

## Validation Expectations

- Vitest output: listed tests green.
- Browser proof: 1k screenshot with run history + protocolos section; redirect proof; reconnect affordance visible on seeded invalid-token installation.
- Grep proof: legacy Integrations component gone; no token values rendered anywhere (audit).

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (after F-03 accepted).
- Next action: compile context pack; read R-01 §1k + R-02 Integrations section + F-01 output + IC-03/IC-05.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
