# AGENTS - Marketplace Central

## Start Here

On every session start, read:

1. `ARCHITECTURE.md`
2. `wiki/README.md`
3. `.mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md`
4. The active milestone `execution-guide.md` named by the current goal, when applicable

Use `IMPLEMENTATION_PLAN.md` only for historical reconciliation. Current planning and execution truth live in `.mnfs/`.

## Truth Order

- Architecture truth: `ARCHITECTURE.md` and ADRs in `docs/architecture/decisions/`
- Contract truth: `contracts/api/marketplace-central.openapi.yaml`
- Wiki truth: `wiki/README.md` and module/quality pages
- Execution truth: tests, builds, QA evidence, commits
- Historical truth: `IMPLEMENTATION_PLAN.md`

When these disagree, classify the mismatch before coding. Stop on architecture, contract, runtime, module ownership, or verification contradictions.

## Project Principles

- YAGNI: build the smallest durable solution that solves the real problem.
- Global maximum: do not optimize legacy paths, workarounds, or wrong abstractions.
- Evidence first: decisions and close-out claims need facts, verification, and references.
- Industry-grade boundaries: domain, application, ports, adapters, and transport stay separate.
- Maintainability is a feature: clarity, naming, auditability, and safe change velocity matter.

## Senior Code Style

Write code that makes the domain obvious and repetition hard to reintroduce:

- Prefer meaningful structs/value objects over loose parameter lists or map-shaped data.
- Name concepts from the business: `StockPolicy`, `ProductLink`, `MarginQuality`, `SafeStockAction`.
- Keep functions small and semantic: validate, normalize, calculate, map, persist, dispatch.
- Extract shared helpers only when they remove real duplication or clarify a repeated pattern.
- Keep business rules in application/domain code, never in React, transport handlers, or provider adapters.
- Keep provider payloads at the adapter boundary; translate them into domain types before use.
- Use typed enums/constants for statuses, modes, provider codes, and error reasons.
- Avoid magic strings, magic numbers, boolean flag soup, and functions with many related parameters.
- Prefer module-owned policies/services for repeated rules such as stock safety, margin quality, and link resolution.
- Match existing local patterns before inventing a new abstraction.

Good code should read like a careful model of the operation, not a sequence of patches that happen to pass.

## Module Boundaries

Every backend module under `apps/server_core/internal/modules/*` follows:

```text
domain/       pure entities, value objects, enums
application/  use cases and orchestration
ports/        interfaces owned by the module
adapters/     Postgres, provider clients, external dependencies
transport/    HTTP decode/encode only
events/       async contracts when needed
readmodel/    query-optimized projections when needed
```

Rules:

- Transport delegates to application services.
- Application code does not import `net/http`, `pgx`, provider SDKs, or another module's internals.
- Adapters implement ports; connectors do not own tenant business state.
- New modules must be registered in `composition/root.go`.
- Cross-module access goes through ports/application APIs, not repositories or SQL from another module.

## Backend Rules

- Every tenant-owned query scopes by `tenant_id`.
- Shared/global tables must be explicit, documented, and protected.
- Every handler validates method/input and returns structured JSON errors.
- Error codes use `MODULE_ENTITY_REASON`.
- Every handler logs `action`, `result`, and `duration_ms`.
- Every write is idempotent or explicitly protected against duplicate execution.
- Use `pgxpool.Pool`; do not introduce raw `sql.DB`.
- No `panic()` in production code.
- Money stays `float64` in domain and `numeric(14,2)` in Postgres until an ADR changes it.
- Use `GOCACHE=.gocache` for local Go test/build commands.

## Validation Policy

- Never declare an external integration, internal read source, or environment-dependent flow "validated" using only mocks, fakes, seams, fixtures, or compile-only coverage.
- Mocks/fakes/seams are allowed to validate contract shape, deterministic business rules, and failure handling, but they do not prove real runtime integration.
- Any work that touches real providers, real databases, real credentials, real queues, or real environment configuration must distinguish clearly between:
  - contract validation
  - integration validation
  - production-like/live validation
- A milestone or feature that depends on real integration behavior is not `passed` unless the validation artifact records real-environment evidence or explicitly stays in a non-passed status.
- If real validation is impossible in the current session, record the exact blocker, the missing environment/prerequisite, and the remaining validation step. Do not silently substitute mock evidence for live evidence.
- Validation artifacts must say whether each command ran against fake/test doubles or against a real target system.

## Frontend Rules

- Data comes only through `packages/sdk-runtime`.
- React components render state; business math stays in Go.
- No direct backend `fetch()` from feature packages.
- No localStorage/SQLite/browser persistence as source of truth.
- Every data-fetching surface has loading, error, and empty states.
- Page-level UI lives in `packages/feature-*`; shared primitives live in `packages/ui`.

## Mercado Livre Safety

Any provider write that can affect stock, price, order handling, or customer communication must prove:

- the internal product link is resolved and unambiguous
- the source data and timestamp are visible
- the policy/rule used for the action is explicit
- the write is idempotent or duplicate-safe
- the audit record stores before/after values and provider response

Unknown cost, freight, fee, tax, or product linkage is a data-quality state, not a zero/default value.

## Workflow

- For new work, name the owning module, contract surface, data source, side effects, and verification path.
- Substantial planning and implementation follows the Claude MNFS workflow in `.claude/plugins/mnfs-workflow`: Mission -> Milestone -> Feature, with `.mnfs/` artifacts as execution truth.
- Use MNFS commands/roles for large work: `mission-init`, `milestone-start`, `feature-context`, `feature-accept`, `milestone-validate`, `mission-validate`, `correction-create`, `mission-closeout`, and `status`.
- MNFS is dry-run first. Write `.mnfs/` artifacts only with explicit apply/write/create approval.
- A feature is not accepted without `spec.md`, `plan.md`, changed paths, and `validation.md` evidence. A milestone is not passed without `validation-result.md` and the integrity gate.
- For any integration-facing work, `validation.md` and `validation-result.md` must explicitly separate fake/mock evidence from real-environment evidence.
- For Sankhya, Mercado Livre, OAuth, DB, queue, or network-dependent behavior, fake-only validation can at most prove contract/business logic readiness; it cannot prove end-to-end readiness.
- For API changes, update OpenAPI and `sdk-runtime` together.
- For architecture changes, update `ARCHITECTURE.md`, ADRs in `docs/architecture/decisions/`, relevant wiki pages, and `.mnfs/` execution artifacts when scope or sequencing changes.
- For code changes, run impacted tests/builds before claiming done.
- For docs-only changes, proofread/diff-review before commit.
- One completed task should end with one intentional commit.
- Do not leave uncommitted work at session end unless the user explicitly asks.

## Commit Format

`<type>(<scope>): <what>`

Types: `feat`, `fix`, `docs`, `chore`, `refactor`, `test`

Examples:

- `feat(inventory): add stock divergence policy`
- `fix(connectors): refresh mercado livre token on 401`
- `docs(architecture): record mercado livre first scope`
