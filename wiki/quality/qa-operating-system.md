# QA Operating System

Marketplace Central closes work with evidence, not confidence.

## Truth Layers

- Runtime truth: what actually runs locally or in the target environment
- Contract truth: OpenAPI, SDK runtime, generated types when present
- Architecture truth: `ARCHITECTURE.md` and ADRs
- Wiki truth: module docs, operations references, and roadmap notes
- Execution truth: tests, scripts, QA evidence, and bounded defers

When truth layers disagree, classify the mismatch before continuing.

## Mismatch Classes

| Class | Meaning | Default action |
|---|---|---|
| runtime prerequisite | app, DB, auth, config, or startup truth is not reliable | stop local feature work |
| shared contract prerequisite | OpenAPI/SDK/API shape is wrong or missing | fix contract first |
| module-local implementation | mismatch is inside one module boundary | fix in module |
| screen-local implementation | mismatch is isolated to one frontend surface | fix in feature package |
| wiki-memory drift | docs are stale but runtime/contract are clear | update docs with change |
| workflow/tooling gap | verification path is missing or unclear | create/repair the check |
| architecture contradiction | ownership, source of truth, or dependency direction conflicts | stop and write/revise ADR |
| defer | accepted limitation outside current scope | record owner, reason, and revisit trigger |

## Stop Gates

- Scope gate: name owner module, contract surface, data source, and side effects.
- Prerequisite gate: confirm runtime, auth/config, DB, route, contract, and migration assumptions.
- Implementation gate: tests/builds for the touched slice pass.
- Review gate: findings are listed by severity and dispositioned.
- QA gate: happy, error, empty, persistence, idempotency, and async paths are checked when relevant.
- Evidence gate: commands, outcomes, artifacts, and bounded defers are recorded before closure.

## Mercado Livre Safety Gates

Stock, price, and order actions can affect real revenue. Any Mercado Livre write must prove:

- internal product link is resolved and not ambiguous
- stock/price source and timestamp are visible
- safety policy was applied before the write
- action is idempotent or protected against duplicate execution
- audit record includes before/after values and provider response

## Closure Format

Every completed non-trivial task should report:

- changed scope
- verification commands and results
- QA/review findings and disposition
- explicit bounded defers
- commit or reason no commit was created
