# Mission Validation Contract

```yaml
id: MIS-002
type: mission-validation-contract
status: draft
owner: Mission Strategist
parent: MIS-002
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: mission
```

## Mission ID

MIS-002

## QA Level

QA-2 — verdicts from commands + logs + governed live-Oracle lane evidence; no UI-drive automation required (data-layer mission).

## Required Final State

All five milestones passed; Oracle-backed pages cost O(1) Oracle queries; cache layers active with `as_of` visible; no regression in layer boundaries or error redaction.

## Criteria

## Criterion: Catalog page is one Oracle query
ID: MIS-02-C01
Level: Mission
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: `GOCACHE=.gocache go test ./internal/modules/catalog/... ./internal/modules/internal_read/... -run 'PageQueryCount|CatalogPage' -v` (fake queryer counts calls) + M-01 latency log excerpt before/after
- Expected: exactly 1 recorded Oracle query per `ListCatalogProductFacts` page call (any page size 1..100); baseline log shows 1+3N pattern eliminated
- Actual:
- Artifact: `validation-result.md` + `M-02-catalog-batch-cutover/validation-result.md`
Blocking failure: any page call issuing >1 Oracle query
Blocking failure observed: No
Owner: QA Validator

## Criterion: Unknown facts never become zero
ID: MIS-02-C02
Level: Mission
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: `GOCACHE=.gocache go test ./internal/modules/internal_read/... -run 'Nullable|Quality' -v`
- Expected: page items with missing stock/price/cost return JSON `null` + quality flag (`missing_cost` etc.); no `0` substitution in any new port
- Actual:
- Artifact: `validation-result.md`
Blocking failure: any new batch/page port emitting 0 for an unknown fact
Blocking failure observed: No
Owner: QA Validator

## Criterion: Layer boundaries and build health
ID: MIS-02-C03
Level: Mission
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `go build ./...` + `GOCACHE=.gocache go test ./...` in `apps/server_core`; `npm run build` in `apps/web` (post M-05); grep proves no godror/SQL import outside `internal_read/adapters/oracle`
- Expected: all green; grep returns only adapter-package matches
- Actual:
- Artifact: `validation-result.md`
Blocking failure: SQL/driver types leaked outside adapter; failing build/tests
Blocking failure observed: No
Owner: QA Validator

## Criterion: No credential or raw driver detail leaks
ID: MIS-02-C04
Level: Mission
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: `GOCACHE=.gocache go test ./internal/modules/internal_read/... -run 'SafeCause|Redact|Leak' -v` + review of new log statements in M-01 instrumentation
- Expected: 503 bodies contain only `source_unavailable`; log lines contain method name + duration + oracle numeric code at most — never username/password/connect string/raw driver text (Sankhya reader's raw-cause gap fixed by M-03)
- Actual:
- Artifact: `validation-result.md`
Blocking failure: any log/body carrying credentials, DSN, or raw driver error text
Blocking failure observed: No
Owner: QA Validator

## Criterion: as_of and force-refresh work end to end
ID: MIS-02-C05
Level: Mission
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: httptest transcript: GET (miss) → GET (hit, same `as_of`) → GET with `Cache-Control: no-cache` (new `as_of`)
- Expected: second response `as_of` equals first; third response `as_of` strictly newer; linkage candidates endpoint shows no cache hit ever (log proves miss on repeat)
- Actual:
- Artifact: `validation-result.md` + `M-04-server-cache/validation-result.md`
Blocking failure: cache serving linkage data; no-cache not bypassing
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

Command transcripts + structured-log excerpts per criterion; live-Oracle facts only via `scripts/run-live-oracle-docker.ps1` outputs referenced from M-01/M-02 results. Milestone rollups: `M-0*/validation-result.md`.

## Blocking Failures

Any criterion's blocking-failure line; plus: mock-only evidence presented as live-Oracle proof.

## Retry Policy

Per milestone contract (max 2 correction attempts each); mission validation runs after M-05 passes.

## Handoff

- Current status: draft (criteria pending)
- Next owner: QA Validator (after M-05)
- Next action: milestone execution per mission strategy
- Required files/evidence: milestone validation-result files, this contract
- Blockers or open decisions: None - contract complete
