# Evidence-Grounded Production Engineering for LLM Agents

> **Status:** DERIVED / PORTABLE GUIDE — NON-AUTHORITATIVE BY DEFAULT  
> **Audience:** human, LLM and hybrid software-engineering work  
> **Parent method:** [DevelopmentConexus Engineering Method](https://github.com/developmentconexus-ops/conexus-methodology/blob/9c7210d1504bef01c0d134a6c3ae8627deebb535/METHOD.md), canonical organizational method at the repository's accepted pin  
> **Purpose:** guide production-grade technology research, reuse, implementation and proof without becoming product architecture, repository status authority or a second engineering method  
> **External references last reviewed:** 2026-08-19

## 1. Adoption and authority contract

This guide operationalizes the DevelopmentConexus Engineering Method for research and production implementation. It does not redefine that method.

Before using this guide, an agent MUST:

1. read the repository bootstrap/instructions;
2. locate the repository's status/router and accepted architecture/decision authorities;
3. revalidate the exact repository, branch/ref and current versions involved;
4. distinguish current target authority from code, tests, history, examples and review evidence;
5. treat this document as guidance unless the repository explicitly adopts it.

Authority order is repository-specific, but the following separation always applies:

```text
repository authority        → what this product/system must mean and preserve
normative specifications    → how protocols and standards actually behave
official technology docs    → supported behavior of the selected version
validated implementations   → proven realization patterns and operational lessons
experiments/tests            → whether the claim holds in this concrete system
```

A reference implementation, popular framework or large-company practice never overrides product authority. Existing code never becomes target authority merely because it exists.

A repository adopting this guide SHOULD link it from its own bootstrap or engineering index. Copying the file alone does not create authority.

---

## 2. Objective

For every material engineering decision, find the smallest sustainable production solution that:

- preserves the repository's accepted invariants and ownership boundaries;
- uses established standards and mature components where they genuinely fit;
- removes accidental complexity instead of moving it elsewhere;
- makes failures, uncertainty and recovery explicit;
- is testable against the real claim;
- remains operable, observable, secure and reversible enough for its blast radius;
- avoids speculative platforms, generic frameworks and compatibility tax without a real consumer.

The target is not maximum novelty and not maximum reuse. It is the best evidenced fit.

---

## 3. Mandatory epistemic discipline

Every material conclusion MUST classify its basis:

| Class | Meaning |
|---|---|
| **Known** | Directly established by current repository authority, primary documentation, executable evidence or real runtime observation. |
| **Inferred** | Reasoned from cited evidence; the inference and assumptions are explicit. |
| **Unknown** | Evidence is insufficient or contradictory. No convenient default is allowed. |
| **Deferred** | The question belongs to a later stage or real consumer; deferral, safety basis and reopen trigger are explicit. |

Never convert:

```text
unknown → zero / false / empty / unsupported / safe
provider response → business success
mock success → live integration proof
example code → supported production contract
current implementation → target architecture
popularity → fitness
latest version → compatible version
```

When external facts may have changed, re-check current official sources and the exact adopted version. Do not rely on model memory for current APIs, security guidance, version support or provider behavior.

---

## 4. Research protocol

### 4.1 Formulate the property before searching

Do not search for “best architecture”, “best Go framework” or “best database pattern”. State the concrete property and failure class first.

Examples:

```text
How do we prevent two concurrent consumers from applying the same consequential intake?

How do we preserve one externally accepted but locally ambiguous effect without blind replay?

How do we prevent a stale client from overwriting a newer owner decision?

How do we prove a query remains tenant-isolated across every access path?
```

A technology search without a target invariant is solution shopping.

### 4.2 Source hierarchy

Use the strongest claim-relative sources available:

1. repository authority and accepted decision artifacts;
2. normative standards and specifications;
3. official documentation for the exact technology/provider version;
4. official repositories, release notes, security advisories and reference implementations;
5. mature open-source implementations with comparable failure properties;
6. peer-reviewed research or well-documented production engineering reports;
7. community articles, discussions and examples only as discovery leads or corroboration.

For a material choice, one blog post, benchmark screenshot, search snippet or generated answer is insufficient.

### 4.3 Research record

Record proportionately:

```text
question / protected property
repository authority consulted
current versions and environment
primary sources consulted
implementation references inspected
known / inferred / unknown / deferred
credible alternatives
selected direction and why it fits
proof strategy
residual risks and reopen triggers
```

Research artifacts are evidence, not parallel product authority.

### 4.4 Stop condition

Stop researching when:

- the required property and constraints are clear;
- primary sources establish the relevant behavior;
- credible alternatives have been compared;
- the selected option's fit and failure modes are understood;
- a proportional proof can falsify the claim;
- no material contradiction remains.

Do not continue collecting references after the decision signal has saturated.

---

## 5. Adopt, adapt, build or defer

Use this decision order:

| Outcome | Select when |
|---|---|
| **ADOPT** | A standard/native capability or mature component satisfies the invariant without distorting ownership or adding disproportionate operations. |
| **ADAPT** | A proven pattern is close, but a bounded adapter or owner-specific specialization is required. |
| **BUILD** | No existing option satisfies a real material property, and the smallest custom component can be clearly bounded and proven. |
| **DEFER** | The consumer, failure class or requirement is not real yet and the seam can be added later without structural rework. |
| **STOP** | A prerequisite or external capability required for correctness is not proven. |

Before building custom machinery, answer:

1. What exact gap prevents use of a standard/native capability?
2. Which defect class would remain reachable if we only adapted an existing option?
3. What is the smallest custom surface?
4. Who owns its meaning and lifecycle?
5. How will it be replaced or removed?
6. Which negative test proves it is necessary and correct?

“Not invented here” and “invented here” are both biases. Evidence decides.

---

## 6. Dependency and repository validation

A dependency is production architecture, even when introduced by one import.

Before adopting a library, service, generator, framework or repository, evaluate:

| Dimension | Required questions |
|---|---|
| **Authority** | Is it official, a recognized reference implementation or a maintained independent project? |
| **Problem fit** | Does it protect our actual invariant, or only resemble the use case? |
| **Version fit** | Is the exact version compatible with our language, runtime, database, browser and tooling versions? |
| **Maintenance** | Are releases, issues, security fixes and maintainers active enough for the risk? |
| **Correctness evidence** | Are there tests, CI, conformance suites, fuzzing or real production references? |
| **Security** | Is there a security policy, advisory history and responsible update path? |
| **Operations** | What state, processes, upgrades, backups, metrics and failure modes does it add? |
| **Dependency graph** | What transitive dependencies, binaries, build steps and supply-chain exposure enter? |
| **License** | Is the license compatible with intended distribution and modification? |
| **Lock-in** | Can the component remain behind a bounded interface? What data/protocol migration would replacement require? |
| **Reversibility** | Can we remove it without rewriting business authority or persistent meaning? |
| **Total complexity** | Does it reduce system complexity, or only move complexity into configuration and operations? |

Stars, download counts and brand reputation are discovery signals, not proof.

Prefer:

- standard-library or platform-native behavior when sufficient;
- explicit, narrow dependencies over broad frameworks;
- components whose failure behavior can be tested locally;
- stable interfaces around external dependencies;
- pinned, reviewable versions and reproducible dependency resolution;
- upgrades as deliberate changes with release-note, compatibility and regression review.

Do not fork a dependency casually. A fork creates ownership of security updates, compatibility and release engineering.

---

## 7. Production engineering workflow

### 7.1 Before code

A production change begins with:

```text
current repo/ref revalidated
→ authority and scope reconstructed
→ materiality classified
→ root cause and target invariant stated
→ failure modes attacked
→ standards/reference research completed
→ alternatives compared
→ proof strategy defined
→ write scope authorized
```

For material behavior, define how the claim could be proven false before implementation.

### 7.2 During code

An agent MUST:

- preserve one authority for each material meaning;
- keep provider/framework DTOs behind adapters;
- make illegal states structurally difficult or impossible where reasonable;
- prefer explicit types and bounded contracts over maps/property bags;
- use database/schema constraints for database-owned invariants where appropriate;
- fail closed on ambiguity at security, tenant, identity and external-effect boundaries;
- keep errors contextual without leaking secrets or PII;
- propagate cancellation/deadlines across blocking or remote work;
- bound concurrency, memory, payloads, retries, queues and external calls;
- make retry behavior explicit and idempotency claim-relative;
- preserve historical evidence needed for explanation/recovery without creating payload archives;
- keep files/components small enough to understand and test independently;
- avoid unrelated refactoring and speculative extensibility;
- update contract, implementation and tests together when one contract changes.

Never introduce silent fallback, fake success, hardcoded production answers or tests that only imitate the dependency whose behavior is claimed.

### 7.3 After code

Verification MUST report exact commands and outcomes, not “should pass”.

Run the repository-prescribed gate plus the smallest additional proof required by the claim. A skipped or unavailable proof remains unproven and is stated explicitly.

Before claiming completion:

- inspect the final diff;
- confirm only intended files changed;
- run relevant verification from a clean enough environment;
- exercise at least one negative fixture for each new control;
- revalidate remote HEAD after publish;
- confirm CI/check status rather than assuming it;
- preserve Unknowns and residual risks in the handoff.

---

## 8. Proof matrix

Use proportional proof; not every change requires every row.

| Claim | Suitable proof examples |
|---|---|
| Pure function/value semantics | table tests, property tests, fuzzing, boundary cases |
| Type/API contract | compile/type failure, schema validation, contract diff, generated-client conformance |
| Database invariant | constraint violation test, transactional integration test, migration test |
| Concurrency | race detector, deterministic concurrent test, stress test, linearization/invariant check |
| Idempotency | duplicate and concurrent intake tests, changed-request same-key negative test, restart recovery |
| External effect | sandbox/live controlled proof, ambiguity injection, authoritative reread/reconciliation |
| Tenant/security boundary | cross-tenant negative tests, permission matrix, threat-model abuse cases |
| Recovery | crash/restart, timeout, lost-response, redelivery, expired credential/cursor tests |
| Performance | representative dataset, `EXPLAIN (ANALYZE, BUFFERS)` where safe, load profile, resource limits |
| Frontend behavior | component/integration tests, browser E2E, accessibility checks, network/error states |
| Observability | prove the alert/metric/trace/log fires and identifies the failure without leaking data |
| Deployment/migration | forward migration, rollback/roll-forward strategy, backup/restore or recovery rehearsal |

A mock proves only the behavior of the mock boundary. It does not prove a provider, database engine, OAuth server, browser or network.

---

## 9. Stack-specific research lenses

These are evaluation lenses, not universal stack mandates. Apply them only when the repository has selected the technology.

### 9.1 Go

Prefer the language and standard library before adding framework machinery, while recognizing real gaps.

For production Go:

- follow the repository's adopted Go version and official documentation for that version;
- keep packages cohesive, dependency direction explicit and interfaces consumer-owned;
- use `context.Context` for request-scoped cancellation/deadlines; do not store it as ambient global state;
- wrap errors with useful operation context while preserving inspectable causes where needed;
- avoid goroutine leaks, unbounded fan-out and shared mutable state without explicit synchronization;
- test concurrent code with the race detector where materially applicable;
- fuzz parsers, wire boundaries and invariant-heavy value logic where arbitrary input can expose edge cases;
- scan reachable vulnerabilities with `govulncheck` or the repository-approved equivalent;
- manage dependencies through Go modules, review `go.mod`/`go.sum`, and treat `replace` directives as explicit temporary risk;
- benchmark only representative hot paths and compare results under controlled conditions;
- prefer generated code only when its source contract is singular and drift is mechanically detected.

Official anchors:

- https://go.dev/doc/
- https://go.dev/doc/security/best-practices
- https://go.dev/doc/security/fuzz/
- https://go.dev/doc/modules/managing-dependencies

### 9.2 PostgreSQL

Database choices must derive from persistent meaning, invariants, query patterns and concurrency—not ORM defaults.

For production PostgreSQL:

- model exact data meaning before physical tables;
- use `NOT NULL`, `CHECK`, `UNIQUE`, `FOREIGN KEY` and `EXCLUDE` where they correctly own the invariant;
- understand the selected isolation level and re-check rules for concurrent updates;
- make transaction boundaries match one consistency claim rather than one HTTP handler by habit;
- avoid application-only uniqueness or “check then insert” races when a database constraint can own the property;
- add indexes for measured query/access patterns, accounting for write/storage cost;
- inspect real plans and representative cardinalities with `EXPLAIN`; use `EXPLAIN ANALYZE` only where executing the statement is safe;
- keep migrations reviewable, forward-safe and tested against representative state;
- choose expand/contract or a hard cutover based on actual compatibility consumers, not ritual;
- test lock contention, deadlocks, retries, large transactions and restart/recovery where material;
- keep backups, restore evidence, connection limits and vacuum/statistics behavior within the operational design.

Official anchors:

- https://www.postgresql.org/docs/current/
- https://www.postgresql.org/docs/current/transaction-iso.html
- https://www.postgresql.org/docs/current/ddl-constraints.html
- https://www.postgresql.org/docs/current/using-explain.html

### 9.3 HTTP, REST, OpenAPI and OAuth

Use standard semantics instead of private conventions.

- derive resource/method/status meaning from HTTP semantics;
- use Problem Details for machine-readable API failures when adopted, while preserving valid business outcomes as domain semantics;
- keep idempotency, optimistic concurrency, authorization and business disposition separate;
- use one machine-readable wire authority when the repository adopts OpenAPI;
- derive or mechanically prove supported clients and server behavior against that authority;
- select an OpenAPI version from required expressiveness plus actual generator/tool compatibility—not recency alone;
- prove conformance controls with deliberate negative drift;
- keep OAuth/OIDC protocol identity separate from product membership/permission/business authority;
- apply current OAuth Security BCP guidance and exact provider documentation to the selected flow;
- never expose tokens, authorization codes, client secrets or raw provider failures in normal logs/problems/history.

Normative anchors:

- https://www.rfc-editor.org/rfc/rfc9110.html
- https://www.rfc-editor.org/rfc/rfc9457.html
- https://www.rfc-editor.org/rfc/rfc9700.html
- https://spec.openapis.org/oas/latest.html

### 9.4 React, TypeScript and server-state clients

The frontend is a client of accepted server/domain authority unless the repository explicitly decides otherwise.

- use TypeScript strict checking unless a documented compatibility constraint prevents it;
- distinguish server state, URL/navigation state, form draft state and local UI state;
- do not duplicate server authority into context/local state by convenience;
- avoid contradictory, redundant and deeply duplicated React state;
- understand the exact query library defaults for staleness, retry, refetch, cache lifetime and structural sharing;
- classify mutations by consequential semantics; a UI retry button never makes an ambiguous effect safe to repeat;
- render known-empty, unknown, unavailable, stale, partial, pending and error states honestly where the contract distinguishes them;
- keep access control on the server; route/button visibility is usability, not authorization;
- test keyboard, focus, names/labels, error messaging, contrast and other applicable accessibility criteria;
- use browser E2E evidence for browser behavior; component mocks alone do not prove routing, storage, cookies, redirects or network interaction.

Official anchors:

- https://react.dev/learn/choosing-the-state-structure
- https://www.typescriptlang.org/tsconfig/strict.html
- https://www.typescriptlang.org/tsconfig/strictNullChecks.html
- https://tanstack.com/query/latest/docs/framework/react/guides/important-defaults
- https://www.w3.org/TR/WCAG22/

### 9.5 External integrations and data pipelines

For every external read/write, establish:

```text
correct Organization/tenant
+ source/installation namespace
+ authentication/credential capability
+ operation authority and coverage
+ provenance/freshness
+ failure and absence semantics
+ duplicate/order/replay behavior
+ authoritative reread or reconciliation
```

Rules:

- consumer owns meaning; adapter owns protocol;
- provider notifications are triggers/evidence unless the contract proves occurrence authority;
- never fetch arbitrary callback-supplied URLs;
- a completed page walk proves only the provider-defined traversed scope;
- partial/unavailable acquisition never becomes deletion or complete absence;
- provider 2xx never becomes convergence without the claimed authoritative proof;
- ambiguous potentially accepted writes are not blindly retried;
- PII and raw payload retention are minimized;
- real integration claims require controlled real-dependency evidence.

### 9.6 Observability and reliability

Observability is part of the protected property, not decorative telemetry.

- define service/user-level indicators from real objectives before creating dashboards;
- instrument traces, metrics and logs only where they support diagnosis, capacity, security or objective measurement;
- propagate correlation context without turning it into business identity;
- use bounded-cardinality attributes and avoid secrets/PII;
- make consequential async work, retries, backlog, dead-letter/refusal and recovery observable;
- alert on actionable conditions with ownership and response guidance;
- test that telemetry and alerts fire under injected failure;
- define timeouts, retries, load shedding, backpressure and graceful shutdown at external/async boundaries;
- size reliability controls to real risk and objectives; do not copy hyperscale topology without the scale or failure class.

Reference anchors:

- https://opentelemetry.io/docs/
- https://opentelemetry.io/docs/concepts/signals/
- https://sre.google/books/

### 9.7 Security and supply chain

Security requirements derive from assets, trust boundaries, threats and exposure.

- threat-model material flows before choosing controls;
- apply least privilege and current authorization at every relevant boundary;
- validate untrusted input at the owning boundary and encode output for its context;
- use parameterized database access and avoid command/shell construction from untrusted input;
- keep secrets outside source, logs, browser storage, problems and business history;
- rotate/revoke credentials without rewriting business identity;
- scan dependencies and reachable vulnerabilities;
- pin and review build/deployment dependencies proportionately;
- require review for authentication, authorization, cryptography, deserialization, file upload and tenant-boundary changes;
- derive security test cases from a maintained verification standard such as OWASP ASVS when it fits the application class;
- never claim compliance from checklist presence alone—prove the controls that support the claim.

Reference anchor:

- https://owasp.org/www-project-application-security-verification-standard/

---

## 10. File upload and untrusted binary guidance

When a product admits uploads, an agent MUST separately decide:

- real consumer and owner of the uploaded meaning;
- accepted media/content types and maximum sizes;
- streaming and memory/disk bounds;
- server-side content verification rather than trusting filename or declared MIME alone;
- filename normalization and path traversal prevention;
- malware/content scanning requirement based on exposure and use;
- storage isolation and authorization for retrieval;
- immutable content identity/digest where idempotency or history requires it;
- retention, deletion, legal/privacy and historical-reference rules;
- image/document transformation failure semantics;
- safe download headers and active-content handling;
- observability without payload/PII logging.

Do not create a generic Asset/Media platform merely because one operation accepts a file. Conversely, do not omit the retrieval, security and lifecycle properties required by the real consumer.

---

## 11. LLM execution contract

A production-oriented LLM SHOULD use this compact protocol.

### 11.1 Orientation

```text
1. Identify repo, branch/ref and requested outcome.
2. Revalidate HEAD and working state.
3. Read repository instructions/router/accepted authorities.
4. State allowed scope and blocked work.
5. Identify current versions and real external dependencies.
```

### 11.2 Decision

```text
6. Classify materiality.
7. State evidence as Known / Inferred / Unknown / Deferred.
8. Identify root cause, target invariant and failure class.
9. Research primary sources and validated implementations.
10. Compare credible alternatives, including adopt/adapt/build/defer.
11. Select the smallest sustainable option.
12. Define proof and reopen triggers before implementation.
```

### 11.3 Implementation

```text
13. Confirm write scope/authorization.
14. Make the smallest coherent change.
15. Preserve boundaries, explicit failure and data meaning.
16. Add positive and negative tests that exercise the protected subject.
17. Run repository gates and claim-specific verification.
18. Inspect diff and publish only intended files.
19. Revalidate remote commit and CI/checks.
20. Report proven, unproven, residual and next-authoritative action separately.
```

### 11.4 Mandatory refusal to invent

An agent MUST stop, defer or keep Unknown rather than:

- fabricate an API, field, library behavior or provider capability;
- cite a repository without inspecting the relevant version/content;
- silently choose a framework because it is familiar;
- introduce a new business authority through technical reuse;
- generate mocks/fallbacks that conceal missing production capability;
- weaken a test/gate to make the change pass;
- claim verification that was not executed;
- claim “best practice” without naming the property, source and context fit;
- promise background work or future evidence not produced in the current execution.

---

## 12. Material technology decision template

Use only when the choice is material; trivial work should remain trivial.

```markdown
# <Decision>

## Authority and scope
- Repository/ref:
- Responsible authority/stage:
- Allowed / excluded:

## Protected property
- Root cause:
- Target invariant:
- Failure class:

## Evidence
- Known:
- Inferred:
- Unknown:
- Deferred:
- Official/normative sources:
- Implementations/repositories inspected:
- Exact versions/environment:

## Alternatives
1. <option> — fit, gaps, operational cost
2. <option> — fit, gaps, operational cost
3. <option> — fit, gaps, operational cost

## Decision
- Outcome: ADOPT | ADAPT | BUILD | DEFER | STOP
- Why this is the Global Maximum for current constraints:
- Essential complexity preserved:
- Accidental complexity rejected:

## Enforcement and proof
- Structural/runtime controls:
- Positive proof:
- Negative proof:
- Real-dependency proof when applicable:

## Residuals
- Risks/Unknowns:
- Reopen triggers:
- Later owner/stage:
```

---

## 13. Production-ready completion test

A change is not production-ready merely because it compiles or has a happy-path test.

The completion claim is proportional, but a material production change should be able to answer:

### Meaning and boundaries

- What authority owns the changed meaning?
- Which invariants are preserved?
- Did the change introduce duplicate state/authority or provider/framework leakage?

### Correctness

- What happens on duplicate, stale, concurrent, partial, out-of-order and invalid input?
- Are unknown, absent and failure states honest?
- Do database/type/schema constraints cover all admitted paths?

### Security and privacy

- What is untrusted?
- How are tenant scope, authorization, secrets and PII protected?
- Which negative abuse cases were exercised?

### Reliability and operations

- What happens on timeout, lost response, dependency outage, restart and overload?
- Is retry safe and bounded?
- Can the system recover and can operators see the condition?

### Data and change management

- Are migrations, retention, historical explanation, rollback/roll-forward and backup implications addressed?
- Is the exact source of truth preserved?

### Verification

- Which commands/tests ran?
- Which real dependency was exercised?
- Which control was deliberately proven to fire?
- What remains unproven?

### Delivery

- Is the final diff scoped?
- Are dependency/license/version changes explicit?
- Is remote HEAD verified?
- Are CI/check results known?
- Is the next action taken only from repository authority?

---

## 14. Portable adoption by MNFS or another software factory

For future reuse:

1. keep the DevelopmentConexus Engineering Method as the canonical reasoning method;
2. adopt this file as a derived production-research/implementation guide, not a second method;
3. require every generated workspace to expose its repository authority path and current status router;
4. inject stack-specific sections only after the project selects those technologies;
5. resolve official documentation and exact versions at execution time;
6. make evidence records, negative proof and verification output first-class build artifacts where material;
7. let reusable agents centralize research, verification and runtime mechanics without becoming product/business authority;
8. version this guide deliberately when repeated evidence shows a rule is missing or misclassifying work;
9. never distribute stale copied variants as competing authorities—keep one adopted source/version per factory and explicit repository specialization.

A software factory should automate repeatable proof and context acquisition, not automate unsupported assumptions.

---

## 15. Reference map

This list is a starting map, not a frozen version catalog. Always resolve current official guidance and the exact selected versions.

### Organizational method

- DevelopmentConexus Engineering Method: https://github.com/developmentconexus-ops/conexus-methodology/blob/9c7210d1504bef01c0d134a6c3ae8627deebb535/METHOD.md

### Protocols and API contracts

- HTTP Semantics — RFC 9110: https://www.rfc-editor.org/rfc/rfc9110.html
- Problem Details for HTTP APIs — RFC 9457: https://www.rfc-editor.org/rfc/rfc9457.html
- OAuth 2.0 Security Best Current Practice — RFC 9700: https://www.rfc-editor.org/rfc/rfc9700.html
- OpenAPI Specification: https://spec.openapis.org/oas/latest.html

### Backend and data

- Go documentation: https://go.dev/doc/
- Go security practices: https://go.dev/doc/security/best-practices
- Go fuzzing: https://go.dev/doc/security/fuzz/
- Go dependency management: https://go.dev/doc/modules/managing-dependencies
- PostgreSQL current documentation: https://www.postgresql.org/docs/current/

### Frontend and accessibility

- React state structure: https://react.dev/learn/choosing-the-state-structure
- TypeScript strict mode: https://www.typescriptlang.org/tsconfig/strict.html
- TanStack Query important defaults: https://tanstack.com/query/latest/docs/framework/react/guides/important-defaults
- WCAG 2.2: https://www.w3.org/TR/WCAG22/

### Security, observability and reliability

- OWASP ASVS: https://owasp.org/www-project-application-security-verification-standard/
- OpenTelemetry documentation: https://opentelemetry.io/docs/
- Google SRE books: https://sre.google/books/

---

## 16. Final invariant

> **An LLM produces production engineering only when its decisions are derived from current repository authority, grounded in primary evidence and validated implementations, realized through the smallest fitting structure, and proven against the real failure properties it claims to protect. Reuse is preferred when it preserves meaning; custom work is justified only by an evidenced gap; Unknown never becomes convenient code.**
