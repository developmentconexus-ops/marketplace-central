# MIS-001-mercado-livre-operating-cockpit

```yaml
id: MIS-001
type: mission
status: in_progress
owner: Mission Strategist
parent: none
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: mission
planning_phase: ready
```

## Objective

Turn Marketplace Central into an internal Mercado Livre operating cockpit that prevents unsafe stock exposure on existing listings, preserves a future multi-marketplace hub architecture, and creates the foundation for order margin and commercial intelligence.

## Outcome

Operators can link existing Mercado Livre listings to internal Sankhya products, compare Mercado Livre announced stock with internal sellable stock, see blocked/unsafe cases before they cause cancellations, and apply audited manual stock corrections through Mercado Livre capabilities. The same business modules and capability ports remain usable by future marketplace adapters.

## Scope

- Remove VTEX target surfaces from code, contract, SDK, UI, tests, and docs.
- Define marketplace capability ports owned by MPC business modules, with Mercado Livre as the first real adapter.
- Define the MNOS/Sankhya read contract needed by MPC without copying the whole MNOS system.
- Build Product Links for existing Mercado Livre listings using EAN and `seller_sku` first, then title heuristics with operator approval for ambiguity.
- Build Stock Seguro using internal sellable stock = `SUM(ESTOQUE - RESERVADO)` for `CODEMP IN (1,2)` and `CODLOCAL = 10101`.
- Keep `CODLOCAL = 10108` showroom inventory excluded from sellable marketplace stock.
- Start with assisted/manual stock actions only; no automatic provider writes.
- Prepare Orders + Margin and Commercial Intelligence as later milestones under the same mission.

## Domain Scope

- Actors and roles: internal operator, future reviewer/approver, background scheduler.
- Core entities: ProviderInstallation, ProductLink, ListingSnapshot, StockSnapshot, StockPolicy, StockAction, MarketplaceOrder, ProfitSnapshot.
- Lifecycle and states: product link candidate -> resolved/conflict/rejected; stock state healthy/oversell/undersell/stale/blocked; stock action proposed/approved/applied/failed/skipped.
- Classification and metadata: provider code, item id, variation id, SKU, EAN, product group, company, location, margin quality.
- Audit and history: every provider write candidate and applied write records before/after, policy, source timestamps, provider response, operator.
- Search, filter, reporting: dashboard filters for risk, link state, product group, stock divergence, listing status, margin quality.
- Admin and config: sellable stock policy, buffer rules, ineligible product groups, marketplace capability status.
- Integration and external: Mercado Livre OAuth/installations, Mercado Livre item/order APIs, MNOS/Sankhya read-only Oracle access through explicit ports.

## Non-Scope

- Publishing brand-new Mercado Livre listings in the first Stock Seguro milestone; existing listings are the target.
- Automatic stock write decisions; all initial writes are assisted/manual and audit-first.
- Full webhooks as the first ingestion path; polling/scheduler comes first, notifications can refresh later.
- Full multi-marketplace operations before Mercado Livre is reliable; future providers must implement the capability contracts later.
- Copying all Sankhya data into MPC Postgres; MPC stores only its own state and snapshots needed for audit/operation.

## Current State

- `integrations` foundation, provider catalog, OAuth, OpenAPI, SDK runtime, and frontend packages exist.
- Mercado Livre provider registration and OAuth adapter exist.
- Mercado Livre fee sync adapter exists at foundation depth.
- Product links, inventory, orders, and profitability modules are planned, not implemented.
- VTEX routes, SDK methods, frontend pages, adapters, tests, and docs still exist and contradict ADR-005.
- MNOS contains the read-only Sankhya/Oracle knowledge needed for stock, product, price, cost, tax, and sales views.

## Clarified Decisions

- Resolved:
  - Existing Mercado Livre listings are the first operational target.
  - Polling/scheduler is acceptable before webhooks.
  - Stock writes start manual/assisted only.
  - Internal stock source comes from MNOS/Sankhya contracts, not from a generic MetalShopping Postgres assumption.
  - VTEX can be removed.
  - Orders + Margin follows Stock Seguro instead of blocking it.
  - Official Mercado Livre docs via Context7 are required for provider semantics.
- Accepted assumptions:
  - Buffer default starts at 1 unit because it is conservative, simple to explain, and reversible per product/group policy.
  - `CODEMP IN (1,2)` and `CODLOCAL = 10101` are the first sellable stock scope because the operator identified revenda as the stock that matters; this remains configurable.
  - Mercado Livre adapter should use direct HTTP, not the deprecated Go SDK, because the official SDK is archived and the docs recommend using current documentation.
- Owner decisions still open:
  - Exact excluded product groups, weight, size, and margin threshold rules will be finalized during Stock Policy feature execution because they require business review of real product examples.
- Blocked items:
  - Independent MNFS readiness review was not run in this session because no fresh `mission-reviewer` dispatch tool is available; mission remains `needs_revision` until that gate runs.

### Clarification Interview

| # | Taxon | Question | Proposed default | Operator answer |
| --- | --- | --- | --- | --- |
| 1 | actor model | Who applies stock writes initially? | Internal operator approves every write | Manual only; no auto-write |
| 2 | lifecycle/transitions | Should webhooks be required before stock/order sync? | Polling first, webhook later | Polling first is acceptable |
| 3 | persistence/reset | Should MPC mirror Sankhya? | No, read Sankhya live through MNOS-derived contracts | Use MNOS/Sankhya direct access, bring only what MPC needs |
| 4 | UI convergence | Should first UI target stock or full margin cockpit? | Stock Seguro first, margin next | Stock Seguro first |
| 5 | validation expectations | Which evidence sources are required? | Official Mercado Livre docs + local MNOS evidence + tests | Use Context7 docs and ask operator where needed |
| 6 | build/runtime conventions | Can VTEX be removed? | Quarantine then remove safely | Remove everything VTEX |

## Architecture Spine

### System Shape

MPC business modules own marketplace-independent operations. Provider adapters implement capability ports. Mercado Livre is the first adapter, not the business model. The web client consumes only `packages/sdk-runtime`; React renders operational state and never performs stock, margin, or provider API calculations.

### Runtime Topology

- Browser: `apps/web`, feature packages under `packages/feature-*`, shared UI under `packages/ui`.
- API server: Go `apps/server_core`, modular monolith.
- MPC store: PostgreSQL for MPC-owned links, policies, snapshots, action audit, orders, margin snapshots.
- Internal data source: Sankhya Oracle through MNOS-derived read contracts; read-only, no Sankhya writes.
- Provider source: Mercado Livre REST APIs through `connectors/adapters/mercado_livre`.
- Scheduler: initial polling jobs call application services, which call ports; synchronous UI reads use MPC read models and never block on provider availability.

### Runtime Contract

- Business decisions live in `product_links`, `inventory`, `orders`, `profitability`, and `pricing_strategy`.
- Provider HTTP details live only in connector adapters.
- Integrations owns auth, credential lifecycle, installation status, and capability health.
- Product links are required before stock writes, margin resolution, and repricing recommendations can be trusted.
- Unknown cost, freight, fee, tax, stock, or link is a data-quality state, never a zero/default.
- Every write to Mercado Livre must be idempotent or duplicate-safe, audited, and traceable to source values and policy.

### Cross-Cutting Decisions

| Decision | Status | Prevents | Must preserve | Validation impact |
| --- | --- | --- | --- | --- |
| Business modules own native operations; adapters implement provider capabilities | accepted | Mercado Livre rules leaking into the whole app | Future providers can implement the same ports with provider-specific mapping | Interface contract and module import checks |
| Stock Seguro uses `SUM(ESTOQUE - RESERVADO)` over `CODEMP IN (1,2)` and `CODLOCAL=10101` | accepted | Showroom or reserved stock being announced | Policy remains configurable and auditable | Unit tests for policy math and SQL contract tests |
| Initial buffer default is 1 unit | accepted | Race between store sale and ML listing stock | Override by product/group; policy visible in UI | Stock recommendation criteria |
| VTEX is removed, not extended | accepted | Optimizing a dead control plane | Provider catalog can still show future providers | Contract/SDK/router tests have no VTEX routes |
| Polling first, webhooks later | accepted | Blocking useful Stock Seguro on notification setup | Ingestion remains idempotent so webhooks can refresh later | Scheduler and idempotency tests |
| Direct HTTP Mercado Livre adapter, no deprecated Go SDK | accepted | Depending on archived SDK behavior | Adapter is thin, tested, and docs-driven | Adapter tests against documented shapes |

Accepted trade-offs:
- Polling can be stale between runs; this is mitigated by visible source timestamps and by later webhook refresh.
- Buffer 1 may undersell by one unit; this is preferable to overselling in the first assisted workflow.
- Direct HTTP costs more adapter code than SDK usage; it avoids deprecated SDK risk.

## Shared Contracts

| Contract | Boundary | Path | Why it exists |
| --- | --- | --- | --- |
| Marketplace capability interface | Business modules <-> connectors/adapters | `research/marketplace-capability-interface-contract.md` | Prevents provider endpoint shapes from becoming business rules |
| MNOS/Sankhya read interface | MPC application <-> Sankhya read edge | `research/mnos-sankhya-read-interface-contract.md` | Fixes stock, product, price, cost, and sales input semantics |

## Milestone Strategy

| ID | Name | System change | Why this order | Path |
| --- | --- | --- | --- | --- |
| M-01 | VTEX removal and architecture reset | Remove dead VTEX surfaces and align docs/contracts | Clears contradictions before new work | `M-01-vtex-removal-architecture-reset/` |
| M-02 | Marketplace capability framework | Add provider capability ports and Mercado Livre adapter spine | Prevents Mercado Livre hardcoding | `M-02-marketplace-capability-framework/` |
| M-03 | MNOS/Sankhya read contract | Bring required internal data semantics into MPC | Stock Seguro needs trusted internal stock/product data | `M-03-mnos-sankhya-read-contract/` |
| M-04 | Product Links ML | Resolve existing ML listing/variation to internal product | Stock write safety depends on unambiguous links | `M-04-product-links-ml/` |
| M-05 | Stock Seguro ML | Show divergence, risk, recommendation, and manual audited action | First business value: prevent cancellation/oversell | `M-05-stock-seguro-ml/` |
| M-06 | Orders + Margin ML | Ingest orders and calculate sale margin quality | Adds revenue/margin visibility after stock is safe | `M-06-orders-margin-ml/` |
| M-07 | Commercial Intelligence | Margin guardrails, aging stock, kits, promotions | Uses stock/link/order/margin foundation | `M-07-commercial-intelligence/` |

## Quality Attributes

| Attribute | Target (concrete) | Owner (ADR/seam) | Validation criterion |
| --- | --- | --- | --- |
| Performance | Stock dashboard API returns 100 linked listings with risk fields in p95 < 500ms using seeded repository tests | Inventory read model | MIS-001-C07 |
| Security | No access/refresh tokens or Sankhya secrets appear in logs, audit payloads, UI responses, or validation artifacts | Integrations + Sankhya read seam | MIS-001-C05 |
| Reliability | Listing, stock, and order ingestion are idempotent by provider id and provider updated timestamp/resource | Capability framework + module repositories | MIS-001-C03 |
| Observability | Every provider action logs `action`, `result`, `duration_ms` and persists audit before/after/policy/provider response | Inventory action audit | MIS-001-C04 |
| Usability | Operator-facing stock/link/margin screens have loading, error, empty, blocked, conflict, and stale states | Feature packages | MIS-001-C06 |
| Maintainability | Business modules cannot import provider HTTP packages, `net/http`, or another module's internals | Module boundary checks | MIS-001-C02 |
| Compatibility | Go uses `pgxpool.Pool`, React uses SDK runtime only, and Sankhya access remains read-only | Project conventions | MIS-001-C01 |

## Non-Functional Scope

| Declined attribute | Reason |
| --- | --- |
| High availability/SLA | Internal cockpit first; reliability is targeted, but no multi-region/runtime HA promise in this mission |
| Full accessibility certification | UI must be usable and have states, but formal WCAG audit is outside this mission |

## Validation Strategy

Validation is milestone-gated. Each feature creates `spec.md`, `plan.md`, and `validation.md`. Each milestone creates `validation-result.md`. Final mission validation rolls up milestone evidence into `validation-result.md`.

Evidence path convention:
- Feature execution evidence: `<feature-root>/validation.md`
- Milestone QA rollup: `<milestone-root>/validation-result.md`
- Mission QA rollup: `.mnfs/MIS-001-mercado-livre-operating-cockpit/validation-result.md`

## Risks

| id | risk | likelihood (L/M/H) | impact (L/M/H) | mitigation | trigger | owner |
| --- | --- | --- | --- | --- | --- | --- |
| R-001 | Ambiguous EAN/SKU/title links could map ML listing to wrong product | M | H | Link states include conflict/unresolved and block writes until approved | Candidate maps multiple products/listings | Product Links owner |
| R-002 | Sankhya stock rules may miss real sellable exceptions | M | H | Policy is explicit, visible, and configurable by company/location/group/product | Operator identifies stock mismatch with real item | Inventory owner |
| R-003 | Mercado Livre API semantics differ by variation or distributed stock setup | M | H | Use official docs, adapter tests, and source timestamps; block unsupported write modes | Adapter sees unknown item shape | Connectors owner |
| R-004 | Removing VTEX may break old tests/routes unexpectedly | M | M | Inventory first, remove with OpenAPI/SDK/router test evidence | Contract or router test fails | Architecture reset owner |
| R-005 | Margin quality can be overstated if freight/manual adjustments are missing | H | M | Missing freight/fee/tax/manual input creates `missing_*` quality states | Profit snapshot lacks required input | Profitability owner |
| R-006 | Context drift between MNOS and MPC data meanings | M | H | Copy only contract semantics and cite MNOS source files; no ad hoc SQL | New query bypasses contract | Sankhya read owner |

## Research Links

- `research/mercado-livre-and-hub-evidence.md`
- `research/marketplace-capability-interface-contract.md`
- `research/mnos-sankhya-read-interface-contract.md`
- Existing local research: `docs/research/2026-07-06-mercado-livre-operating-cockpit.md`

## Handoff

- Current status: In progress; M-01 has passed and M-02 is starting.
- Current owner: Milestone Orchestrator.
- Next owner: Feature Implementer for M-02/F-01.
- Next action: Execute M-02/F-01 capability port contract, then route returned `spec.md`, `plan.md`, changed paths, and `validation.md` through milestone acceptance review.
- Required artifact paths:
  - `.mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md`
  - `.mnfs/MIS-001-mercado-livre-operating-cockpit/validation-contract.md`
  - `.mnfs/MIS-001-mercado-livre-operating-cockpit/architecture-map.md`
  - `.mnfs/MIS-001-mercado-livre-operating-cockpit/research/`
- Required evidence paths:
  - `.mnfs/MIS-001-mercado-livre-operating-cockpit/validation-result.md`
  - `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-*/validation-result.md`
- Blocked decisions: None for M-02/F-01; product exclusion rules and margin threshold require business examples before M-05/M-07 execution.
