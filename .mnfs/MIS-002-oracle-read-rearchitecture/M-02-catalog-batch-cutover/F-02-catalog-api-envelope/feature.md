# F-02-catalog-api-envelope

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-02
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-002. Contract of record: IC-01.

## Milestone

M-02. Depends on F-01 (port exists).

## Brief

Expose the catalog page port through transport with the IC-01 cursor envelope on the paginated catalog listing route, plus `GET /catalog/products/search` (`?q=&limit=<1..50>`, same envelope, `next_cursor` always null, per IC-01 search row); enforce the full error matrix; update OpenAPI spec and `sdk-runtime` typed client in the SAME commit (repo rule).

## Inputs

- F-01 port + typed errors.
- IC-01: envelope example, error matrix (400 `invalid_cursor`/`invalid_limit`, 503 `source_unavailable`, 504 `deadline_exceeded`), route namespace rule (existing `/catalog/*` prefix, no new prefixes), limit range 1..100, default 50 (ceiling 200 belongs to `ImportMarginInputs` only — M-03).
- Existing transport handlers + OpenAPI file + `packages/sdk-runtime/src/index.ts` patterns.
- Route class: interactive (15s), registered per M-01 route-class mechanism.

## Expected Output

- Handler: query params `cursor`, `limit`; envelope `{items, next_cursor, page_size, as_of}`; `Cache-Control: no-cache` passthrough plumbed to FreshnessPolicy MaxAge=0 (cache arrives M-04; wire the header→policy mapping now so M-04 is drop-in).
- OpenAPI paths+schemas for the new route; sdk-runtime method `listCatalogProductFacts({cursor?, limit?})` returning typed envelope.
- Old unpaginated listing flow removed/redirected per IC-01 deprecation note.
- One intentional commit.

## Constraints

- No business logic in handler — decode/validate/map only.
- Error bodies exactly IC-01 codes; no driver or SQL detail.
- OpenAPI + sdk-runtime same commit; `npm run build` green in sdk package.
- Do not add caching (M-04) or frontend changes (M-05).

## Inputs/Outputs

Per IC-01 examples verbatim (envelope JSON, item shape, error bodies).

## Negative Scenarios

- While `cursor=%%%garbage`, when GET runs, the system shall return 400 `{"error":"invalid_cursor"}` without an Oracle call.
- While `limit=101`, when GET runs, the system shall return 400 `{"error":"invalid_limit"}` (allowed range 1..100 per IC-01 catalog route row).
- While the port returns source-unavailable, when GET runs, the system shall return 503 `{"error":"source_unavailable"}`.

## Validation Expectations

- httptest transcripts: happy 3-page chain; all four error cases.
- `git show --stat` proving OpenAPI + sdk-runtime in same commit; sdk build green.
- `GOCACHE=.gocache go test ./...` green.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` at execution time.

## Handoff

- Current status: briefed
- Next owner: Feature Implementer (mpc-implementer, gpt-5.6-luna high)
- Next action: spec → plan → implement → evidence
- Required files/evidence: `validation.md`
- Blockers or open decisions: None
