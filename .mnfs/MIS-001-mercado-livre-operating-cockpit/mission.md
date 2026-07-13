# MIS-001-mercado-livre-operating-cockpit

```yaml
id: MIS-001
type: mission
status: ready
owner: Mission Strategist
parent: none
created: 2026-07-06
updated: 2026-07-13
validation_level: QA-0
lifecycle_scope: mission
planning_phase: ready
```

## Objective

Validate Marketplace Central as a coherent internal Mercado Livre cockpit using
real read-only Mercado Livre and Sankhya data, durable MPC-owned state, and an
operator-facing journey that connects product, listing, stock, sale, and margin
without requiring production authentication or provider writes.

## Outcome

A trusted local operator starts from a stock attention item, opens its exact
Mercado Livre listing, inspects the linked Sankhya Product 360, returns to that
listing, opens a related sale and its explainable margin, returns to the same
listing, and reviews a stock-correction simulation that is visibly not executed.
The browser, API, PostgreSQL, Mercado Livre reads, and Oracle reads are evidenced
as one vertical journey rather than as isolated module demos.

## Product Intent

Marketplace Central is an MVP for validating operating concepts, integrations,
and UX. It is not yet a production-complete ERP, marketplace write engine, or
multi-user control plane. The UI is organized around operator objects and
journeys; the backend remains a modular monolith with domain-owned rules.

## Scope

- Preserve the passed M-01 through M-05 and M-08 outcomes and all historical QA.
- Preserve M-06 as a failed historical production-grade gate while reusing its
  implemented order, profitability, Sankhya-linkage, and UI evidence.
- Establish canonical product identity from Sankhya before building Product 360.
- Reorganize the browser into Overview, Products, Listings, Sales, and Operations.
- Add deep links and a shared installation context across those workspaces.
- Use real Mercado Livre and Oracle reads plus real PostgreSQL persistence.
- Present stock changes as previews with current value, proposed value, reason,
  source timestamps, and payload; no live write is required for MVP acceptance.
- Validate a bounded real vertical sample in the browser, including unknown,
  stale, ambiguous, error, and empty states.

## Domain Scope

- **Actors and roles:** one trusted local internal operator; integration runtimes.
- **Core entities:** SankhyaProduct, ProviderInstallation, MarketplaceListing,
  ProductLink, StockSnapshot, StockRecommendation, MarketplaceOrder,
  MarginInput, ProfitSnapshot, AttentionItem.
- **Lifecycle and states:** listing link `unresolved|conflict|resolved|rejected`;
  workspace data quality `current|stale|unknown|conflict`; simulation
  `draft|reviewed`. Attention is derived from current domain facts and has no
  separately persisted lifecycle.
- **Collaboration:** none in MVP because there is one local operator.
- **Classification and metadata:** product identity, EAN, manufacturer reference,
  provider item/variation, installation, link state, stock risk, margin quality,
  source and observed timestamp.
- **Audit and history:** imported-source timestamps, link decisions, manual local
  adjustments, simulations, and integration runs remain inspectable; actor labels
  are explicitly `operator_supplied_unverified`.
- **Notifications and delivery:** an in-app attention queue only.
- **Search, filter, reporting:** product, listing, sale, integration health, data
  quality, stock risk, and margin filters; deep links preserve context.
- **Admin and configuration:** installation health, OAuth state, policies, probes,
  and integration run history in Operations.
- **Integration and external:** real Mercado Livre REST reads, real Oracle reads,
  and PostgreSQL persistence; provider writes are simulated for acceptance.

## Non-Scope

- Production authentication, RBAC, tenant/user authorization, and trusted-principal
  attribution are deferred because the approved MVP is a single-operator local tool.
- Provider writes are not an acceptance dependency; any retained experimental write
  path is disabled by default and cannot be represented as authenticated.
- Automatic repricing, automatic stock writes, listing/SKU mutation, and durable
  unknown-result reconciliation are deferred to a production phase.
- Multi-marketplace operation, fiscal issuance, purchasing, fulfillment, shipping,
  customer service, notifications, and multi-company workflows are excluded.
- Competitor-price monitoring, kits, promotions, and broad commercial intelligence
  remain post-MVP discovery because the operator selected only the stock workflow.
- Formal SLA, multi-region availability, and full WCAG certification are excluded.

## Current State

- M-01 through M-05 and M-08 are passed and remain immutable historical outcomes.
- M-06 fixed-SHA review passed at `1eb8831fb1d0d1b84f4d1325978bbc4f76c9ed0f`,
  but proportional QA failed because C03 required an unproved trusted-principal
  boundary and the registered Windows Go command was invalid.
- The current UI exposes separate technical routes for Products, Classifications,
  Marketplaces, Integrations, Product Links, Stock Seguro, Orders, and Simulator.
- Orders already combine order, inputs, adjustments, and margin; other journeys
  remain fragmented and lack Product 360, Listing detail, global installation
  context, and a cross-domain attention queue.
- Current SDK methods already cover most read/import/link/stock/order/profitability
  capabilities needed for the MVP; M-13 must prefer those methods over new backend
  breadth.
- Legacy MSDB product identity and provider/runtime exceptions remain documented;
  only identity defects that block the MVP are owned by M-09.

## Clarified Decisions

- Resolved:
  - The product is an MVP for validating real integrations, concepts, and UX.
  - The actor model is one trusted local operator with no login.
  - Audit actor labels are not authenticated and must say so.
  - Provider writes are not required to pass; stock change is a simulation.
  - The UI is organized by Overview, Products, Listings, Sales, and Operations.
  - Existing technical modules stay separate behind the operator workspaces.
  - Validation uses deterministic tests, real read-only ML/Oracle smoke,
    PostgreSQL, and browser-driven vertical evidence.
  - Official commands must run in Windows PowerShell; Go cache paths are absolute.
- Accepted assumptions:
  - The app remains single-tenant while business tables retain `tenant_id`; this is
    reversible because no multi-user authorization behavior is exposed.
  - Desktop is the primary operating surface and one 390x844 responsive pass is
    retained; this is reversible because no mobile-only workflow is introduced.
  - M-09 removes only legacy identity paths needed by the MVP; broader runtime
    consolidation can follow later without changing operator contracts.
  - Existing SDK operations are sufficient for the first workspace composition;
    if a required field is absent, the owning feature must stop and amend IC-003
    plus OpenAPI/SDK together instead of inventing a client-only calculation.
  - Composite listing route segments use `installation~item~variation` and the
    literal `-` for a null variation. Provenance: IC-003 route decision approved
    with the joined-workspace UX; reversible through legacy redirects plus one
    coordinated IC/OpenAPI/SDK amendment before dispatch.
  - Deterministic UI evidence uses installation `inst-mvp-ml`, products `1001`
    and `1002`, listing `MLB-MVP-1`, and order `ORDER-MVP-1`; reset is limited to
    fixture-owned rows. Provenance: the accepted no-auth local seed in IC-003;
    reversible because these are test identities, not production mappings.
- Validation boundary approved: proportional MVP proof keeps deterministic tests,
  real read-only Mercado Livre and Oracle checks, PostgreSQL persistence,
  secret/PII protection, no provider writes, and the vertical browser journey.
- The round-3 requests for atomic evidence wrappers, hash manifests, OCR automation,
  dependency-claims ledgers, and hostile-origin certification are explicitly not
  readiness gates for this trusted-local MVP. See `readiness-review.md`.

### Clarification Interview

| # | Taxon | Question | Proposed default | Operator answer |
| --- | --- | --- | --- | --- |
| 1 | actor model | Identity depth for the MVP | Trusted local operator, no login | Approved option A |
| 2 | lifecycle/transitions | Is a provider write required to pass? | No; simulation is sufficient | Approved option A |
| 3 | persistence/reset | How should state behave? | Durable PostgreSQL with resettable MVP data | Approved recommendation |
| 4 | UI convergence | One screen per module or joined journeys? | Joined operator workspaces | Replan from hub references and UX as a whole |
| 5 | validation expectations | What proves the MVP? | Tests + real reads + PostgreSQL + browser | Approved recommendation |
| 6 | build/runtime conventions | What is the command environment? | Windows PowerShell-safe commands | Approved recommendation |

## Architecture Spine

### System Shape

The React client presents object-centered workspaces and consumes only
`packages/sdk-runtime`. Go modules keep business rules and adapter isolation.
PostgreSQL stores MPC-owned links, snapshots, adjustments, and run history.
Oracle and Mercado Livre remain external sources read through owned adapters.

### Runtime Topology

- Browser: `apps/web`, route shell, and `packages/feature-*` workspaces.
- API: `apps/server_core`, modular Go monolith, no `/v1` prefix.
- Store: PostgreSQL for MPC-owned state; reset is an explicit local operation.
- Internal source: read-only Sankhya Oracle through `internal_read` adapters.
- Provider source: Mercado Livre REST through connector/integration adapters.
- Transport: `packages/sdk-runtime`; web never performs provider/Oracle calls.
- Authentication: no application session/cookie; existing OAuth remains integration
  credential lifecycle, not operator identity.

### Runtime Contract

- `catalog` owns canonical internal product identity; `product_links` owns the
  listing-to-product decision; `inventory` owns stock risk/recommendation;
  `orders` owns normalized sales; `profitability` owns margin quality;
  `integrations` owns connection and run health.
- Workspaces may compose read models but may not recreate domain calculations.
- Unknown cost, tax, freight, fee, stock, link, or source time remains null plus
  an explicit quality state; it never becomes zero.
- Every browser request goes through SDK runtime. API changes update OpenAPI and
  SDK runtime in the same commit.
- Stock preview has no execution transition in the MVP contract.

### Cross-Cutting Decisions

| Decision | Status | Prevents | Must preserve | Validation impact |
| --- | --- | --- | --- | --- |
| ADR-006 MVP evidence boundary | accepted | Production hardening blocking concept validation | Real reads, honest labels, no accidental writes | MIS-001-C09, M-14-C01 |
| ADR-007 Object-centered workspaces | accepted | Navigation mirroring technical packages | Backend module ownership and deep links | MIS-001-C10, M-13-C01 |
| ADR-008 Product and listing separation | accepted | Provider listing identity becoming internal SKU | One product to many listings and explicit link states | MIS-001-C11, M-09-C01 |
| ADR-009 Proportional security | accepted | Unauthenticated labels being represented as trusted | Secret/PII controls and default-disabled writes | MIS-001-C05, M-13-C04 |
| ADR-010 Vertical validation | accepted | Isolated feature tests being reported as an MVP journey | Real sources, browser actions, evidence provenance | MIS-001-C12, M-14-C02 |
| ADR-011 Historical gate preservation | accepted | Rewriting M-06 history to manufacture a pass | Fixed SHA and failed QA remain intact | MIS-001-C13 |

Accepted trade-offs:
- No login means actor attribution is not trustworthy; the UI and evidence label it
  `operator_supplied_unverified` and provider execution is outside acceptance.
- Composed workspaces may issue several SDK reads; no performance SLA is set until
  the bounded vertical sample identifies a measured bottleneck.
- Deferring runtime consolidation retains known technical debt, but avoids changing
  provider infrastructure before UX evidence identifies a need.

## Shared Contracts

| Contract | Boundary | Path | Why it exists |
| --- | --- | --- | --- |
| Marketplace capability interface | modules ↔ provider adapters | `research/marketplace-capability-interface-contract.md` | Keeps provider payloads out of business modules |
| Oracle internal-read interface | application ↔ Oracle | `research/mnos-sankhya-read-interface-contract.md` | Owns ERP read semantics and unknown handling |
| MVP operator workspace contract | web routes ↔ SDK/domain reads | `research/mvp-operator-workspace-interface-contract.md` | Fixes object identities, navigation, attention states, and simulation semantics |
| MVP validation lane contract | QA commands ↔ real/local runtimes | `research/mvp-validation-lane-interface-contract.md` | Fixes PowerShell commands, schemas, evidence paths, ordering, and no-write gates before dispatch |

## Milestone Strategy

| ID | Name | System change | Why this order | Path |
| --- | --- | --- | --- | --- |
| M-01–M-05 | Accepted product foundations | VTEX removal, capabilities, Oracle reads, links, Stock Seguro | Preserved historical prerequisites | Existing milestone paths |
| M-06 | Orders + Margin historical gate | Implemented order/margin foundation; QA remains failed | Evidence is reused but the verdict is not rewritten | `M-06-orders-margin-ml/` |
| M-08 | Development harness | Passed bounded session/checkpoint control plane | Preserved support prerequisite | `M-08-repository-integrity-harness/` |
| M-09 | Canonical Product Foundation | Make `CODPROD` the MVP product anchor and remove blocking MSDB identity paths | Product 360 cannot use ambiguous identity | `M-09-canonical-product-foundation/` |
| M-13 | Integrated Operator Workspaces | Replace module-shaped navigation with joined Overview, Product, Listing, Sales, and Operations journeys | Reuses implemented capabilities before adding backend breadth | `M-13-integrated-operator-workspaces/` |
| M-14 | Real Vertical MVP Validation | Prove one bounded real read-only journey plus simulated stock correction | Final MVP evidence after identity and UX converge | `M-14-real-vertical-mvp-validation/` |
| M-07 | Commercial Intelligence | Reassess recommendations from validated real data | Starts only after M-14 | `M-07-commercial-intelligence/` |

Deferred milestones M-10, M-11, and M-12 retain their IDs and original intent but
move after MVP validation. M-10 runs only when runtime debt blocks evolution; M-11
owns production writes/auth; M-12 owns real listing/SKU mutation.

## Quality Attributes

| Attribute | Target (concrete) | Owner (ADR/seam) | Validation criterion |
| --- | --- | --- | --- |
| Security | No secret/PII in UI, logs, or evidence; all provider writes disabled by default; actor label is `operator_supplied_unverified` | ADR-009 | MIS-001-C05, M-13-C04 |
| Reliability | Repeating the same listing/order import yields one durable identity and no duplicate adjustment/snapshot | Module repositories | MIS-001-C03 |
| Observability | Operations shows source, observed time, last run status, and actionable error for every MVP integration | IC-003 Operations | M-13-C03 |
| Usability | Five workspaces expose loading, error, empty, stale/unknown, blocked/conflict, and success where applicable; all attention links open filtered context | ADR-007 | MIS-001-C10, M-13-C01 |
| Maintainability | Browser uses SDK only; modules do not import provider adapters or another module's internals | Architecture boundaries | MIS-001-C02 |
| Compatibility | Registered Go commands use an absolute `GOCACHE` and valid package cwd in PowerShell; web commands run from repository root | ADR-010 | M-14-C03 |

## Non-Functional Scope

| Declined attribute | Reason |
| --- | --- |
| Performance SLA | Bounded internal MVP sample; measure first and create a target only if the browser journey shows a bottleneck |
| Production availability | Local validation cockpit; no uptime or multi-region promise |
| Production authentication/authorization | Single trusted local operator; security targets secrets, PII, and accidental writes instead |
| Full accessibility certification | Responsive keyboard-usable workflows are required, but a formal WCAG audit is post-MVP |

## Validation Strategy

- Feature execution evidence: `<feature-root>/validation.md`.
- Milestone QA rollup: `<milestone-root>/validation-result.md`.
- Mission QA rollup: `.mnfs/MIS-001-mercado-livre-operating-cockpit/validation-result.md`.
- Deterministic tests prove business rules and UI states.
- Real-environment evidence separately proves Mercado Livre reads, Oracle reads,
  PostgreSQL persistence, and browser behavior; mocks never claim integration.
- The final drive starts at `/`, follows stock attention to Listing, opens Product,
  returns to Listing, opens Sale, returns to Listing, reviews stock simulation,
  reloads, and verifies state and source timestamps.
- IC-004 records a small PowerShell-safe validation ladder, expected evidence, and
  stop conditions. QA may use existing harness commands or direct repository
  commands as long as it records the frozen SHA and does not widen side effects.

## Risks

| id | risk | likelihood (L/M/H) | impact (L/M/H) | mitigation | trigger | owner |
| --- | --- | --- | --- | --- | --- | --- |
| R-007 | Product identity cutover breaks existing links | M | H | M-09 migration/readback contract and immutable provider identities | CODPROD differs from legacy SKU/reference | M-09 owner |
| R-008 | Product 360 becomes a client-side business engine | M | H | SDK-only composition; domain calculations remain server-owned | UI derives stock/margin values | M-13 owner |
| R-009 | Simulation appears to have executed | M | H | `simulation` badge, no execute transition, payload marked preview | UI success copy implies provider mutation | M-13 owner |
| R-010 | Real-source evidence leaks secret or PII | L | H | Redaction checks and values-minimized evidence | artifact contains credential/buyer data | QA Validator |
| R-011 | Historical M-06 fail is silently overwritten | M | H | ADR-011 and immutable fixed-SHA references | artifact claims M-06 passed | Portfolio owner |
| R-012 | Workspace convergence expands into ERP breadth | M | M | Explicit non-scope and existing SDK-first constraint | fiscal/fulfillment/purchasing work requested | Mission Strategist |
| R-013 | Windows QA commands fail again | M | M | Preflight every registered command with exact cwd and absolute cache | command uses Unix assignment or wrong package path | M-14 owner |

## Research Links

- `research/marketplace-hub-ux-benchmark-20260713.md`
- `research/current-ui-journey-inventory-20260713.md`
- `research/mvp-operator-workspace-interface-contract.md`
- `research/mvp-validation-lane-interface-contract.md`

## Handoff

- Current status: Ready under the operator-approved proportional validation boundary.
- Current owner: Portfolio Hub.
- Next owner: M-09 Milestone Orchestrator in a clean user-started task.
- Next action: Execute M-09 from the accepted planning SHA using the harness packet.
- Required artifact paths: `mission.md`, `architecture-map.md`, `validation-contract.md`, `M-09-canonical-product-foundation/`, `M-13-integrated-operator-workspaces/`, `M-14-real-vertical-mvp-validation/`.
- Required evidence paths: each planned feature `validation.md`, each planned milestone `validation-result.md`, and mission `validation-result.md`.
- Blocked decisions: None for M-09 dispatch. A genuine identity/contract conflict
  discovered during execution must be asked in the visible Milestone.
