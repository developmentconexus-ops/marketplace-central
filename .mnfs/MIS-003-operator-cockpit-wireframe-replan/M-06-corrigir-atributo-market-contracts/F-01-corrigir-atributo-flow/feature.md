# F-01-corrigir-atributo-flow

```yaml
id: F-01
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

MIS-003. Binding contracts: IC-03 (`listing_edit` attribute intent), IC-05 (states/copy), R-01 (corrigir-atributo mini-flow per P3 review — full wizard is Non-Scope), `GOV_API_SDK_SPLIT`.

## Milestone

M-06 corrigir-atributo-market-contracts.

## Brief

Server: `GET /listings/categories/{categoryId}/attributes` — the IC-02 operation row granted to this feature (single permitted extension of the M-01-owned `/listings` prefix) — ML category attributes via connectors capability, canonical shape (id, name pt, type, required, allowed values, constraints), cached (L2-style, 24h class — new class `category_meta`). Client: guided fix flow launched from listing pendência rows flagged with attribute gaps (Anúncios drawer + produto detail): pick flagged attribute → form field generated from attribute schema (enum→select, number→input with unit, boolean→toggle) with client-side validation from constraints → preview/confirm via M-03 modal carrying `listing_edit` attribute intent → protocolo result.

EARS:
- While a category's attributes are cached and fresh, when the flow opens, no provider call shall fire (cache proof).
- While an attribute is enum-typed, when the form renders, only allowed values shall be selectable; free text impossible.
- While client validation passes but provider rejects, when applying, the item shall fail `provider_validation` with provider detail behind "▸ técnico" (honest surface, no client mask).
- While the fix applies successfully, when terminal, `invalidateAfterMutation('listing_edit')` shall fire and the pendência flag shall clear after listings refresh.

## Inputs

- IC-02 `getCategoryAttributes` row (owned here by grant), IC-03 listing_edit intent schema, connectors capability for category attributes (verify existing capability or add read method — adapter-layer only), L2 cache class registry, M-03 modal, R-01 pendência flag inventory, M-04 F-01 detail surface.
- Early de-risk (readiness ★6 advisory): before building the form renderer, run one live read-lane verification of a real ML category-attributes payload (read-only, env-var credentials, no values in artifacts) to confirm the canonical mapping covers real attribute shapes; record source/type/time in `validation.md`. Mocked renderer tests do not substitute for this one honest read.

## Expected Output

- Endpoint + OpenAPI + sdk-runtime same commit; L2 class `category_meta` registered.
- Schema-driven form component (three type renderers minimum: enum/number/boolean; string fallback).
- Flow wired in both entry points; component tests: renderer per type, constraint validation, provider-rejection surface.
- Integration test: attributes endpoint against stubbed capability.

## Constraints

- No full wizard, no listing_create path, no category CHANGE (attribute values only) — P3 scope decision.
- Raw ML attribute payload stays in adapter; canonical DTO out.
- Writes only via M-03 envelope; pt-BR.

## Negative Scenarios

- Unknown categoryId → 404 `category_not_found`.
- Capability failure → 502-mapped honest error; UI ErrorState with retry (no stale-cache fabrication if cache empty).
- Required attribute submitted empty → client block pre-preview.
- Value violating numeric constraint → client block with constraint copy.

## Interaction Model

Flow state: select-attribute → edit(validate live) → M-03 modal takes over (its machine). Entry rows show flag source ("atributo obrigatório ausente: {name}"). Cache: TanStack Query key `['listings','category-attributes',categoryId]`, staleTime 24h client-side mirroring `category_meta`.

## Validation Expectations

- Vitest output: renderer/validation/rejection tests green.
- `go test` output: endpoint + cache integration tests.
- Browser proof: flow screenshots select→form→preview→protocolo result on stub; pendência flag cleared post-refresh.
- Diff proof: OpenAPI + sdk-runtime same commit.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: compile context pack; verify connectors category-attributes capability existence first (R-03 route gap allowed: one scoped investigation).
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: capability read method may need adding at adapter layer (bounded, adapter-only).
