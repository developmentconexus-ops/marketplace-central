# Cross-Repository Engineering Alignment — Marketplace Central ↔ MetalDocs

> **TEMPORARY REVIEW EVIDENCE ONLY — NOT PRODUCT / ARCHITECTURE / ROADMAP AUTHORITY**
>
> This file exists only for the bounded cross-repository review requested by the operator. It must be deleted after convergence. Neither repository may inherit a decision from this file without its own owning-stage adjudication and explicit operator approval where required.

## 0. Exact review subjects

### Marketplace Central

```text
repository        developmentconexus-ops/marketplace-central
candidate branch  stage/d6-frontend
candidate SHA     cb55238c1908b087989825ff4d2ad9ce6f08527b
candidate PR      #54 — Draft / open / unmerged
program state     D0–D5 accepted; D6 open
D6-B1             OPERATOR-RATIFIED
active decision   D6-B2 frontend topology/dependency adjudication
D7–D9             BLOCKED
implementation    BLOCKED UNTIL D9
```

### MetalDocs

```text
repository        developmentconexus-ops/MetalDocs
main              0b4ef6ef891b01f907804cff4bd3c0022aebad80
candidate branch  arch/t8h-global-coherence
candidate SHA     8c2ae8515fecf513cfd699e9d0e53eb2551fd835
candidate PR      #148 — Draft / open / unmerged
program state     T1→T8-G integrated and operator-ratified
T8-F              OPERATOR-RATIFIED / INTEGRATED
T8-G              OPERATOR-RATIFIED / INTEGRATED
T8-H              OPERATOR-RATIFIED / INTEGRATION PENDING
T9→T12             NOT OPEN
implementation    BLOCKED
```

Current candidate CI was revalidated before opening this dialogue:

```text
Marketplace Central cb55238c...  required PASS + pr-title PASS
MetalDocs 8c2ae851...             required CI #1117 SUCCESS
```

## 1. Review purpose

The operator wants the two repositories to challenge each other before Marketplace Central freezes D6-B2.

The goal is **not** shared code, a shared monorepo, a platform repository, or forced uniformity.

The goal is:

> **When the two products have the same engineering property and failure class, prefer the same proven technology/pattern unless evidence proves a better repo-specific choice. When the properties differ materially, preserve the difference explicitly rather than forcing symmetry.**

This review applies DevelopmentConexus Root Cause / Global Maximum / YAGNI / falsification discipline to both repositories.

## 2. Non-negotiable review law

1. Repository current authority beats this dialogue and beats the other repository.
2. A choice already ratified in MetalDocs is not automatically correct for Marketplace Central.
3. A newer Marketplace idea is not automatically a reason to reopen MetalDocs.
4. Same-looking features are insufficient; compare the protected property and failure class.
5. Shared technology is desirable only when it reduces total complexity without weakening either product's authority model.
6. No generic internal platform/framework may be invented to "share" the repositories.
7. No code/package is shared by implication.
8. No stage jump is allowed:
   - Marketplace D6-B2 may freeze only frontend realization decisions.
   - Marketplace backend/auth/runtime/persistence mechanics remain D7 obligations/candidates, not D6 authority.
   - MetalDocs T8 authorities reopen only on a material falsifier, not preference.
9. Technology claims must be revalidated against current official documentation / maintained upstream repositories for the exact mechanism before final convergence.
10. Exact dependency versions belong to implementation manifests unless a version-specific semantic property requires architecture freeze.

Every reviewed item must end as one of:

```text
ALIGN
DIVERGE_JUSTIFIED
REOPEN_MARKETPLACE
REOPEN_METALDOCS
DEFER
STOP
```

## 3. Smallest authority packs

Do not recursively read either repository.

### Marketplace Central bootstrap

```text
AGENTS.md
→ docs/index.md
→ docs/roadmap.md
```

For frontend questions use only the routed D6 pack, primarily:

```text
docs/engineering/rebaseline/D6-FRONTEND.md
ARCHITECTURE.md
```

For a concrete auth/identity question switch to:

```text
docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md
```

For technology research use the routed derived guide only when the concrete question requires it:

```text
docs/development/evidence-grounded-production-engineering-for-llm-agents.md
```

### MetalDocs bootstrap

```text
AGENTS.md
→ docs/index.md
→ docs/roadmap.md
```

Switch bounded packs by question:

```text
frontend      docs/architecture/frontend.md + docs/decisions/t8f-ratification.md
backend       docs/architecture/backend.md
AuthZ/AuthN   docs/architecture/authorization-and-audit.md + owning identity/runtime material only as needed
persistence   docs/architecture/persistence.md
runtime       docs/architecture/runtime.md
wire          docs/architecture/wire-contract.md
```

The five-file default applies per concrete question/pack, not by recursively loading all authorities at once.

---

# 4. Marketplace Lead → MetalDocs review

## 4.1 Overall assessment

**PRELIMINARY VERDICT: MetalDocs is a strong reference implementation profile, not a template. No broad MetalDocs reopen is justified.**

The reviewed T8-F/T8-G/T8-B/T8-D decisions are unusually aligned with Marketplace Central's already accepted constraints:

```text
Go backend
React browser client
PostgreSQL canonical product state
OpenAPI wire SSOT
TanStack Query server state
external OIDC boundary
provider claims != Product Authorization
thin generated-contract transport
modular-monolith preference
mechanism != semantic authority
no ORM / BFF / Redis / microfrontends / generic event bus by default
OpenTelemetry / OTLP observability direction
```

The similarities are meaningful because they arise from comparable protected properties, not because the repositories were designed to copy each other.

The review also found at least one **materially justified divergence** already:

```text
MetalDocs Launch:
  one Company root
  RLS / pooled-tenancy substrate explicitly deferred

Marketplace Central:
  Organization is a real tenant/isolation root
  target may not rely only on developers remembering predicates
  exact RLS/isolation mechanism remains D7
```

Therefore shared stack does **not** mean identical persistence isolation.

## 4.2 Frontend contract — strong alignment candidate

MetalDocs T8-F ratified:

```text
one React SPA
TanStack Query = server-state authority
router/URL = navigation state
local form state = unaccepted draft
local React state = ephemeral UI
no Redux/Zustand/global server mirror
OpenAPI-generated wire shapes
openapi-typescript
one thin browser transport
feature query/command functions
React lenses
no frontend AuthZ engine
no generic BFF
no SSR/offline/realtime/microfrontend baseline
```

Marketplace D6-B1 independently ratified the same four state classes and the same "frontend is a client, never business authority" law.

### A1 — generated TypeScript wire

**Preliminary classification: ALIGN.**

Use `openapi-typescript` as the TypeScript wire-shape generator unless a current official/proof comparison finds a materially stronger fit.

Reason:

- both repositories already require OpenAPI as the single machine-readable wire authority;
- both reject hand-written DTO shadow authority;
- generated `paths/components` preserve exact operation/path/schema meaning without generating a second business abstraction;
- MetalDocs has independently ratified this approach;
- Marketplace's current D5 proof already proves TypeScript generation from the OAD.

### A2 — browser transport

**Preliminary classification: ALIGN toward thin native `fetch`; challenge required before freeze.**

MetalDocs deliberately ratified:

```text
openapi-typescript
→ one thin native-fetch transport
→ feature query/command functions
```

and T8-G explicitly says `openapi-fetch` is **not** a requirement and must not be introduced for symmetry.

Marketplace had begun considering a more opinionated generated/fetch SDK before this cross-review. That idea is now reopened by evidence.

The shared transport property is not "avoid libraries". It is:

```text
wire shapes generated from OAD
+
one small place for browser mechanics
-
parallel DTO/business authority
-
generic business retry/action framework
```

Marketplace-specific transport must still preserve its own requirements, including exact Organization scope, Idempotency-Key, conditional headers where applicable, Problem Details, raw/non-JSON responses where admitted, and no provider technical-ingress leakage into the Product SDK.

Reviewer must falsify whether native `fetch` remains the smallest common solution or whether a maintained library removes meaningful complexity without hiding these semantics.

### A3 — TanStack Query

**Preliminary classification: ALIGN.**

Same property in both repositories:

```text
server truth cached as query state
!= URL/navigation
!= form draft
!= ephemeral UI
```

Use operation semantic inputs as query identity; do not normalize into a universal client entity store. Mutation invalidation remains owner/lens-specific rather than a generic action framework.

### A4 — router library

**Preliminary classification: DEFER pending reciprocal review.**

MetalDocs explicitly did **not** freeze a router library. Its durable contract is the stable route meanings, not implementation technology.

Marketplace has heavier contextual URL state (Organization + optional/exact Marketplace Installation + performance periods/comparison + filters) and may have a stronger need for typed search-param validation.

Do not force one of these outcomes merely for uniformity:

```text
TanStack Router everywhere
React Router everywhere
hand-written browser routing everywhere
```

Question is property-first:

> Does a router dependency materially reduce invalid URL/context states, route/search duplication and agent-generated drift in both repositories at lower total complexity than a small explicit route layer?

If yes, align. If only Marketplace has the property, `DIVERGE_JUSTIFIED` is acceptable.

### A5 — feature/package topology is the main unresolved frontend tension

**Preliminary classification: REOPEN_MARKETPLACE candidate design / reciprocal challenge required; MetalDocs authority itself remains coherent.**

MetalDocs ratified feature folders by stable semantic **lenses / human flows**:

```text
features/library
features/document-official
features/document-work
features/governance-work
features/history
features/audit
features/admin/...
```

It explicitly says frontend topology is not inferred one-for-one from Go semantic owners.

Marketplace's first unapproved D6-B2 hypothesis leaned toward feature folders named by Product semantic owners:

```text
features/performance
features/market
features/economics
features/offering
...
```

That may still be correct because Marketplace's task-oriented routes often reuse owner-native capability across several screens, but the cross-repo review must attack whether this becomes domain-package mirroring rather than frontend feature architecture.

Required reciprocal question:

> What is the smallest common topology law: lens-first, owner-first, or a bounded combination where route/lens composition consumes owner-scoped query modules without making every D1 owner a top-level frontend feature?

No topology is approved until this is resolved.

### A6 — forms, UI library and design system

**Preliminary classification: DEFER.**

Neither repository has evidence that a broad form framework or universal design-system platform is architectural necessity.

Shared implementation dependencies should be selected only after concrete form/accessibility/component needs are enumerated. Avoid adopting React Hook Form, Zod-as-business-authority, shadcn/Radix/MUI, Storybook, or a design-token platform merely because they are common.

A small shared choice may still be desirable later if both products prove the same accessibility/form/component property.

## 4.3 Authentication / access — strong semantic alignment

### A7 — external OIDC boundary

**Preliminary classification: ALIGN.**

Both repositories already reject local end-user credential authority and provider claims as Product Authorization truth.

Marketplace D2 already states:

```text
interactive human AuthN -> external standards-based OIDC provider
stable binding          -> (issuer, subject)
Keycloak                -> preferred first self-hosted candidate
provider/deployment     -> later realization
```

MetalDocs ratified:

```text
external OIDC Provider boundary
provider-neutral identity adapter
provider roles/groups != MetalDocs Authorization
coreos/go-oidc/v3 + golang.org/x/oauth2 reference mechanism
```

This is a natural shared platform profile candidate:

```text
Architecture depends on OIDC, not Keycloak internals.
Keycloak may be the first concrete provider for both products if deployment/security/operations proof is satisfactory.
```

No shared realm, shared database, shared user table or shared application session is implied.

### A8 — browser session / CSRF realization

**Preliminary classification: DEFER to Marketplace D7; strong alignment candidate.**

MetalDocs has ratified a same-origin server-side ApplicationSession with:

```text
HttpOnly Secure cookie
SameSite=Lax
OIDC callback outside application OpenAPI
synchronizer CSRF token for unsafe requests
```

Marketplace has not opened D7 and therefore may not import this as current authority. D7 should explicitly compare the same pattern before inventing a different auth-session stack.

A divergence requires a real Marketplace-specific security/runtime property.

## 4.4 Backend topology / Go technology

### A9 — modular monolith

**Preliminary classification: ALIGN as a future Marketplace D7 candidate, not a D6 decision.**

MetalDocs T8-B ratified one Go module + owner-first modular monolith + application orchestration + transport + platform + composition + closed-world dependency direction.

Marketplace stable architecture already points toward a modular monolith and separates business authorities, adapters, platform mechanisms and composition. However, its current abstract vocabulary (`contexts/adapters/kernel/platform/composition/views`) is not yet the D7 realization.

Marketplace D7 should compare whether adopting the same **class vocabulary and dependency law** as MetalDocs reduces cross-repository cognitive/tooling drift without distorting Marketplace's 13 business boundaries.

Do not copy MetalDocs owner names or application leaves.

### A10 — HTTP stack

**Preliminary classification: DEFER to Marketplace D7, likely ALIGN.**

MetalDocs chooses Go `net/http` and explicitly rejects Gin/Echo/Fiber without a named need.

Marketplace should not introduce a web framework merely because implementation examples are easier. D7 should identify the exact routing/middleware properties first and compare standard-library + narrow router/mux needs against frameworks.

### A11 — PostgreSQL runtime and SQL posture

**Preliminary classification: ALIGN on base technologies; isolation may DIVERGE_JUSTIFIED.**

Strong common candidates:

```text
PostgreSQL
pgx/v5 + pgxpool
owner-private SQL
sqlc where generated SQL access materially reduces hand-written mapping risk
no ORM by default
narrow explicit transactions/OCC/constraints
```

MetalDocs uses one DB because one local ACID transition is a current property. Marketplace may reach the same conclusion, but must adjudicate its own transaction boundaries in D7.

Isolation divergence:

```text
MetalDocs single-company Launch -> RLS deferred
Marketplace tenant-ready        -> D7 owes fail-closed isolation beyond developer discipline
```

Do not weaken Marketplace tenancy to align the stacks.

### A12 — durable jobs / River

**Preliminary classification: DEFER / do not align by symmetry.**

MetalDocs has named durable-work consumers and therefore ratified River on PostgreSQL.

Marketplace almost certainly has asynchronous marketplace/provider work, but D7 has not yet derived its exact durability/transaction/reconciliation requirements. River is a credible shared candidate only if those requirements match.

No "we already use River in MetalDocs" shortcut is allowed.

### A13 — observability

**Preliminary classification: DEFER to Marketplace D7, strong ALIGN candidate.**

MetalDocs ratified:

```text
OpenTelemetry Go metrics/traces
OTLP vendor-neutral export
log/slog JSON
otelhttp
otelpgx
bounded cardinality/redaction
```

These are generic mechanism concerns and likely satisfy the same Marketplace properties. Marketplace D7 should use this as the default candidate to falsify rather than researching from zero, while remaining free to reject it on evidence.

### A14 — configuration / migration / test infrastructure

**Preliminary classification: DEFER to owning Marketplace stages, prefer common candidates when properties match.**

MetalDocs current reference candidates include:

```text
sethvargo/go-envconfig
tern/v2
sqlc
Testcontainers-Go
OSV-Scanner / govulncheck
```

These are valuable pre-vetted candidates, not automatic Marketplace decisions.

## 4.5 Runtime/deployment shape

### A15 — same-origin SPA + API

**Preliminary classification: DEFER to Marketplace D7; strong alignment candidate.**

MetalDocs ratified one application origin serving SPA + application API + OIDC callback, with no BFF/SSR/service-worker/realtime requirement.

Marketplace has a stable Product origin and React client but D7 owns concrete serving/deployment. The same-origin profile should be tested first because it can reduce CORS/auth/session complexity, but it is not a D6-B2 decision.

### A16 — no speculative infrastructure

**Preliminary classification: ALIGN as engineering law.**

Both products currently reject, absent a named consumer:

```text
microservices
BFF
SSR correctness dependency
offline-first/service worker
WebSocket/realtime
Redis
external Search
service mesh
generic event bus/outbox
generic workflow platform
microfrontends
```

That shared subtraction posture should remain aligned.

## 4.6 Tooling / dependency governance

### A17 — common versions when the same dependency is selected

**Preliminary classification: ALIGN as repository-development policy candidate.**

If both repositories independently select the same dependency for the same property, prefer the same supported major/minor baseline and similar update/advisory cadence unless one repo has a proven compatibility constraint.

This reduces:

```text
LLM/human context switching
duplicate upgrade research
security advisory divergence
inconsistent examples/patterns
```

It does **not** imply a shared lockfile, shared package, shared release cadence or coupled deployment.

### A18 — mechanical boundaries

**Preliminary classification: ALIGN on property; exact tooling DEFER.**

MetalDocs backend ratifies closed-world/default-deny first-party dependency enforcement.

Marketplace's architecture also requires mechanical protection against private/cross-owner leakage. The same proof philosophy should apply. Frontend import-boundary tooling should be selected only after the D6-B2 topology itself is agreed; do not choose an ESLint plugin before knowing the actual allowed-edge graph.

---

# 5. Initial cross-repo alignment matrix

| Area | MetalDocs ratified direction | Marketplace current authority/candidate | Marketplace Lead preliminary result |
|---|---|---|---|
| React browser client | one React SPA | React client | ALIGN |
| Server-state client | TanStack Query | TanStack Query | ALIGN |
| TS wire generation | openapi-typescript | generated TS required; exact D6-B2 tool open | ALIGN candidate |
| Browser transport | thin native fetch | exact client technology still open | ALIGN candidate / challenge |
| Router library | intentionally unfrozen | D6-B2 open | DEFER / challenge |
| Global client store | absent | excluded | ALIGN |
| Feature topology | semantic lenses | first unapproved idea leaned owner-first | REOPEN_MARKETPLACE candidate design |
| Form/UI framework | unfrozen | unfrozen | DEFER |
| Auth protocol | external OIDC | external OIDC | ALIGN |
| First IdP | provider-neutral architecture | Keycloak preferred self-hosted candidate | ALIGN candidate at realization |
| Product AuthZ from IdP claims | rejected | rejected | ALIGN |
| Browser app session | HttpOnly cookie + CSRF | D7 not open | DEFER / likely align candidate |
| Go backend | one-module modular monolith | Go canonical; realization later | DEFER to D7 / likely align |
| HTTP server | net/http | D7 open later | DEFER / likely align |
| PostgreSQL | canonical state | canonical state | ALIGN base technology |
| PostgreSQL driver | pgx/v5 | D7 not open | DEFER / likely align |
| ORM | rejected baseline | no current need | ALIGN subtraction |
| RLS | deferred single-Company | tenant isolation mechanism owed D7 | DIVERGE_JUSTIFIED |
| Durable jobs | River for named consumers | D7 not open | DEFER |
| Observability | OTel + OTLP + slog | D7 not open | DEFER / strong align candidate |
| Config | go-envconfig reference | D7 not open | DEFER |
| SQL generation | sqlc candidate | D7 not open | DEFER |
| Migrations | tern/v2 candidate | D7 not open | DEFER |
| Test env | Testcontainers-Go candidate | later proof stages | DEFER |
| BFF/SSR/offline/realtime | absent | excluded without evidence | ALIGN |
| Shared code/platform repo | absent | not requested | ALIGN: do not invent |

---

# 6. Required MetalDocs reciprocal review of Marketplace Central

MetalDocs reviewer: start fresh and independently verify both exact refs. Do not optimize for agreement with the Marketplace Lead.

Your task is to **attack Marketplace D6-B2 and also attack MetalDocs T8-F/T8-G where Marketplace evidence may expose a stronger Global Maximum**.

## 6.1 Questions you must answer

1. Does Marketplace's ratified D6-B1 create any frontend property that invalidates MetalDocs' `openapi-typescript + thin native fetch + TanStack Query` profile as the best shared baseline?
2. Should Marketplace use lens-first feature topology like MetalDocs, owner-first topology, or a bounded hybrid? Give concrete counterexamples from both products.
3. Is there a material reason for either repository to select a router library now? If yes, compare current official TanStack Router and React Router behavior against the exact route/search/context properties. If no, state the reopen trigger.
4. Would standardizing on `openapi-typescript + native fetch` create hidden manual transport complexity in Marketplace (Idempotency-Key, ETag, Problem Details, exact responses, Technical Ingress separation), or is it still the smallest solution?
5. Is any more opinionated OpenAPI client generator materially superior for both repositories today? If proposing one, prove it does not create a second query/cache/business abstraction.
6. Does MetalDocs' frontend lens topology contain a lesson that should alter Marketplace's D6-B2 feature/package direction?
7. Conversely, does Marketplace's stronger owner/composition separation expose a weakness in MetalDocs T8-F feature topology that merits a bounded reopen?
8. Should both products standardize on Keycloak as the first concrete OIDC provider while keeping OIDC/provider-neutral architecture? Identify real reasons not to.
9. Should Marketplace D7 begin from MetalDocs' backend reference profile (`net/http`, pgx, owner-first modular monolith, platform mechanisms, composition root), or are Marketplace's provider-heavy/tenant-heavy failure classes materially different?
10. Which MetalDocs backend/runtime technologies are genuinely reusable **decisions** for Marketplace, and which are only convenient familiarity?
11. Is River likely to be the Global Maximum for Marketplace durable work, or must it remain unselected until D7 derives exact transaction/retry/reconciliation properties?
12. Is Marketplace's tenant isolation a sufficient property difference to justify RLS or another stronger mechanism while MetalDocs remains without RLS?
13. Should both repos align observability on OpenTelemetry/OTLP + slog when implementation opens?
14. Which frontend build/test/UI dependencies should intentionally remain unfrozen now?
15. Is there any technology already ratified in MetalDocs that current 2026 official evidence now makes obsolete, risky or materially inferior?
16. Is there any Marketplace accepted decision that should be reopened because MetalDocs proves a better architecture class?
17. What shared technology profile can be written **without** creating a shared platform, coupled releases or speculative framework work?

## 6.2 Required output format

Write only under `## MetalDocs reciprocal response` below.

For every material point use:

```text
ID
PROPERTY / FAILURE CLASS
METALDOCS CURRENT DECISION
MARKETPLACE CURRENT DECISION
CURRENT PRIMARY EVIDENCE
CLASSIFICATION = ALIGN | DIVERGE_JUSTIFIED | REOPEN_MARKETPLACE | REOPEN_METALDOCS | DEFER | STOP
RATIONALE
SMALLEST NEXT ACTION
REOPEN TRIGGER
```

Then provide:

```text
SHARED PROFILE — safe to align now
REPO-SPECIFIC DIFFERENCES — must remain different
DEFERRED PROFILE — align later only after owning stage proves need
RECONSTRUCTION / REOPEN DECISION for each repo
```

Do not write to either candidate branch. Do not open implementation work. Do not begin Marketplace D7 or MetalDocs T9.

---

## MetalDocs reciprocal response

<!-- MetalDocs reviewer writes here only. -->

---

# 7. Convergence rule

After the MetalDocs response:

1. Marketplace Lead independently adjudicates every response item against both repositories and current upstream evidence.
2. If a material contradiction survives, reopen only the smallest owning authority in the affected repository.
3. If no material contradiction survives, produce a bounded cross-repo technology alignment result separating:
   - shared frontend decisions that Marketplace D6-B2 may actually freeze;
   - future Marketplace D7 candidates only;
   - MetalDocs decisions that remain unchanged;
   - justified divergences.
4. Operator approves/rejects the bounded result.
5. Delete this temporary dialogue file and close review PR(s) unmerged.
6. Neither repository merges this Evidence file into its candidate/main tree.
