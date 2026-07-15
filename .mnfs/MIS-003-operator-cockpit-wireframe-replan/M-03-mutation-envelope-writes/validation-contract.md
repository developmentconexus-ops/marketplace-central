# Milestone Validation Contract

```yaml
id: M-03
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

M-03-mutation-envelope-writes

## QA Level

QA-2 — lifecycle + gates proven on stub adapter in integration lane (ephemeral-postgres). LIVE ML write optional, only under governed provider-write lane with explicit operator authorization at execution time (RK-03); absence of live run does NOT block pass.

## Required Outcome

IC-03 envelope end-to-end: tables, lifecycle, poller, six write types through 7 gates, HTTP surface, preview/confirm UI + protocolo detail.

## Criteria

## Criterion: Lifecycle transition table enforced
ID: M03-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: `go test` transition-table test enumerating every legal IC-03 arrow and every illegal pair
- Expected: legal transitions succeed; illegal (e.g. draft→approved, applying→cancelled) rejected with 409 lifecycle error and state unchanged in DB
- Actual:
- Artifact: `F-01-protocolo-core/validation.md`
Blocking failure: any illegal transition accepted
Blocking failure observed: No
Owner: QA Validator

## Criterion: Poller exclusivity and crash resume (Q3)
ID: M03-C02
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: concurrency test (two pollers, one approved protocol); kill-mid-chunk test then restart
- Expected: exactly one poller claims (FOR UPDATE SKIP LOCKED); after restart, pre-crash applied items untouched (idempotency key check skips), remaining items complete, terminal state computed from full item set
- Actual:
- Artifact: `F-01.../validation.md` transcripts + SQL asserts
Blocking failure: double claim, double-send, or lost items after restart
Blocking failure observed: No
Owner: QA Validator

## Criterion: Terminal states honest
ID: M03-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: matrix test on stub: all-ok / mixed / all-fail item outcomes
- Expected: `applied` / `partially_failed` / `failed_preserved` respectively; failed_preserved protocol retains all items + failure payloads unmutated on later polls
- Actual:
- Artifact: `F-01.../validation.md`
Blocking failure: wrong terminal state or post-terminal mutation
Blocking failure observed: No
Owner: QA Validator

## Criterion: Seven write gates each rejected (Q2)
ID: M03-C04
Level: Milestone
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: gate matrix tests — one negative per gate
- Expected: missing actor → protocol creation rejected; duplicate idempotency key → item skipped not re-sent; approve without `execute:true` → 400 `execute_required`; unresolved link → item `link_unresolved` with zero provider calls (spy assert); missing policy → item failed `policy_missing`, no default applied; missing source_timestamp → 422 `source_time_unavailable` pre-apply; audit-write failure → item aborts before provider call
- Actual:
- Artifact: `F-02-write-types-adapters/validation.md`
Blocking failure: any gate bypassable
Blocking failure observed: No
Owner: QA Validator

## Criterion: Failure taxonomy mapping
ID: M03-C05
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: mapping table test (stub provider returns 429/401/validation/paused/timeout/unknown)
- Expected: `provider_rate_limited` retryable=true; `provider_auth`; `provider_validation`; `listing_paused_remote`; `provider_unavailable` retryable=true; unknown → `internal` retryable=false with `message_provider` preserved
- Actual:
- Artifact: `F-02.../validation.md`
Blocking failure: unmapped error or lost `message_provider`
Blocking failure observed: No
Owner: QA Validator

## Criterion: ADR-16 SKU invariant
ID: M03-C06
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: `listing_edit` intent test changing SELLER_SKU to value ≠ linked CODPROD
- Expected: rejected at validation pre-apply with item failure `sku_invariant_violation` (IC-03 taxonomy); provider never called
- Actual:
- Artifact: `F-02.../validation.md`
Blocking failure: SKU-divergent edit reaches provider
Blocking failure observed: No
Owner: QA Validator

## Criterion: listing_create contract-only
ID: M03-C07
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: `POST /mutations` with type `listing_create`
- Expected: 422 `{"error":{"code":"type_not_enabled"}}`; zero protocol rows created
- Actual:
- Artifact: `F-03-selection-preview-api/validation.md`
Blocking failure: draft created or different code
Blocking failure observed: No
Owner: QA Validator

## Criterion: Preview snapshot semantics
ID: M03-C08
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: integration: by-filter selection, preview, mutate underlying listing, approve, apply on stub
- Expected: preview snapshots before-values into mutation_items; re-preview replaces snapshot; approve applies the SNAPSHOT values (not re-resolved); preview older than 15 min → approve 409 `preview_stale`; >2000 selection → 422 `selection_too_large`; empty → 422 `empty_selection`
- Actual:
- Artifact: `F-03.../validation.md`
Blocking failure: apply re-resolves selection or caps unenforced
Blocking failure observed: No
Owner: QA Validator

## Criterion: Retry clones to new protocol
ID: M03-C09
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: integration: partially_failed protocol (mix retryable/non-retryable failures) → `POST /mutations/{id}/retry`
- Expected: NEW protocol id `MP-…` with `retried_from` set; only retryable-failed items cloned; original protocol rows byte-identical after retry; zero retryable → 422 `nothing_to_retry`
- Actual:
- Artifact: `F-03.../validation.md`
Blocking failure: original mutated or non-retryable cloned
Blocking failure observed: No
Owner: QA Validator

## Criterion: Preview/confirm UI gating
ID: M03-C10
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: vitest: modal flow tests; browser stub-lane walkthrough
- Expected: confirm control disabled until preview rows render; explicit checkbox + button required; 409 `preview_stale` returns modal to preview step with copy "Prévia expirada. Gere novamente."; double-click confirm issues single approve call; terminal state fires `invalidateAfterMutation` exactly once
- Actual:
- Artifact: `F-04-preview-confirm-ui/validation.md` + screenshots
Drive (UI — agent-browser; UI criteria only; omit for non-UI):
- Fixture: IC-02 seed + stub adapter, installation inst_1
- Steps:
  - open http://localhost:5174/anuncios?installation=inst_1
  - click "Atualizar preço" (with 2 rows selected)
  - assert text "itens" (preview totals)
  - click confirm checkbox, click "Confirmar"
  - assert text "MP-"
- Expected: protocolo id visible in result view with per-item before/after
Blocking failure: one-click apply possible or silent re-approve
Blocking failure observed: No
Owner: QA Validator

## Criterion: Protocolo detail page (Q4)
ID: M03-C11
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser: F5 on `/protocolos/MP-…` mid-apply; failed-item inspection
- Expected: page polls (2s) to terminal after reload; each failed item shows failure code + failureCopy pt-BR string + provider detail behind "▸ técnico"; "Repetir itens com falha" navigates to NEW protocol with `retried_from` backlink
- Actual:
- Artifact: `F-04.../validation.md` screenshots
Blocking failure: reload loses protocol state or raw provider text shown outside técnico
Blocking failure observed: No
Owner: QA Validator

## Criterion: No secrets/PII in envelope surfaces (Q2)
ID: M03-C12
Level: Milestone
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: grep/audit of protocolo API responses, mutation_items rows, and logs from full stub run
- Expected: zero access tokens, zero Authorization headers, zero buyer PII in protocol payloads/audit/logs; `message_provider` sanitized (message text only)
- Actual:
- Artifact: milestone `validation-result.md` audit transcript
Blocking failure: any token/PII in any envelope artifact
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

Feature proofs `F-0*/validation.md`; rollup `validation-result.md` with fixed SHA, dual-gate records, StockActionService regression transcript.

## Blocking Failures

Any criterion blocking failure; StockActionService fold regression (existing flow broken) blocks regardless.

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: none

## Handoff

- Current status: planned.
- Next owner: QA Validator (after F-01..F-04 accepted).
- Next action: execute criteria at fixed SHA; ask operator re live-write lane (optional).
- Required files/evidence: as above.
- Blockers or open decisions: live-write authorization (operator, execution time).
