# F-03-selection-preview-api

```yaml
id: F-03
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

MIS-003. Binding contracts: IC-03 (operations, selection grammar, caps, error matrix, preview TTL), IC-02 (filter grammar reused for by-filter selection). `GOV_API_SDK_SPLIT`.

## Milestone

M-03 mutation-envelope-writes. Depends on F-01 + F-02.

## Brief

Expose the envelope over HTTP: `POST /mutations` (create draft: type + intent + selection by explicit ids OR by IC-02 filter), `POST /mutations/{id}/preview` (resolve selection → snapshot per-item before-rows into mutation_items, return totals + per-item preview page), `POST /mutations/{id}/approve` (requires `execute: true`), `POST /mutations/{id}/cancel`, `POST /mutations/{id}/retry` (clones retryable failed items into NEW protocol), `GET /mutations` (list, filterable), `GET /mutations/{id}`, `GET /mutations/{id}/items` (cursor). Enforce caps: 2000 items, empty selection 422, preview TTL 15 min. OpenAPI + sdk-runtime same commit; `/mutations` dev-proxy row exists from M-02.

EARS:
- While a draft has by-filter selection, when previewed, the API shall resolve the filter against current listings, snapshot each matched row's before-values, and return IC-03 totals `{items, previewed, applied, failed, skipped}`; re-preview replaces the snapshot.
- While a preview is older than 15 min, when approve is called, the API shall return 409 `preview_stale` (operator must re-preview).
- While approve lacks `execute: true`, when called, the API shall return 400 `execute_required` (gate 3).
- While selection resolves >2000 items, when previewing, the API shall return 422 `selection_too_large` with the count.
- While a protocol has zero retryable failed items, when retry is called, the API shall return 422 `nothing_to_retry`.

## Inputs

- IC-03 operations + error matrix + caps (verbatim), IC-02 filter grammar (shared resolver), F-01 domain + F-02 writers, OpenAPI/sdk-runtime update pattern, existing transport middleware (actor/tenant).

## Expected Output

- Eight endpoints per IC-03; actor supplied by the client in the POST body and recorded verbatim (`operator_supplied_unverified`, ADR-009) — missing/empty → 400 `actor_required` (gate 1).
- Selection resolver shared with listings module read path (no duplicate filter parser).
- Retry: new protocol with `retried_from: MP-xxxxxx`, only `retryable: true` failed items cloned; original untouched.
- OpenAPI + sdk-runtime methods (`createMutation`, `previewMutation`, `approveMutation`, `cancelMutation`, `retryMutationFailures`, `listMutations`, `getMutation`, `listMutationItems`) same commit.
- Integration tests: full lifecycle transcript on stub adapter; every error-matrix row.

## Constraints

- No UI (F-04). Transport layer thin; lifecycle rules stay in domain (F-01).
- Approve is the ONLY endpoint that transitions to approved; no combined preview+approve shortcut.
- Filter re-resolution at preview time only — approved set is the snapshot, never re-resolved at apply.

## Negative Scenarios

- Error matrix complete: `execute_required` 400, `preview_stale` 409, `selection_too_large` 422, `empty_selection` 422, `nothing_to_retry` 422, unknown protocol 404, illegal transition 409.
- Approve after listings changed post-preview → applies snapshot values; remote-conflict surfaces per-item as `conflict_remote_changed` at apply (F-02), not approve-time error.

## Validation Expectations

- Integration transcript: create→preview→approve→poll to `applied` on stub; retry flow producing new MP id with `retried_from`.
- Error-matrix table test output: all seven rows status+code asserted.
- Diff proof: openapi.yaml + sdk-runtime same commit.
- SQL proof: re-preview replaces items (count + values), approved snapshot immutable after apply starts.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (after F-02 accepted).
- Next action: compile context pack; read IC-03 + F-01/F-02 module code only.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
