# Marketplace Central — Engineering Rules

> **Scope:** repository-local execution, CI/proof, safety rails, and justified Repository Standard specializations. This file owns no Product semantics, methodology kernel, or mutable program status.

## Canonical methodology

Canonical repository:

`developmentconexus-ops/conexus-methodology`

Accepted pin:

`9c7210d1504bef01c0d134a6c3ae8627deebb535`

Start at pinned `ROUTER.md`. Organizational reasoning, repository operation, frontend planning, realization, and independent review are owned by the applicable pinned Method files. MPC does not copy or redefine them here.

## Local aggregate verification

```powershell
npm run gate
```

`scripts/gate.ps1` is the single repository aggregate implementation and `.github/workflows/ci.yml` exposes one required check named `required` for pull requests and `main`.

The gate protects only durable repository and Product-contract properties that justify permanent automation:

- bootstrap/context budget and canonical routing;
- exact methodology pin;
- candidate/main hygiene for temporary review transport;
- current durable-document reachability;
- one canonical Product OpenAPI entrypoint;
- no runtime/implementation surface while the roadmap blocks implementation;
- canonical Product OAD proof, including its real contract/auth/security negative controls.

CI does **not** permanently parse exact P8/P9 wording, operator ratification prose, review dialogue, historical closure text, prototype labels, or arbitrary negative-control counts. Those remain method-stage Evidence and are reviewed/operated at the gate that consumes them. Do not turn Evidence into a permanent CI authority merely because a script can grep it.

A red current material control is a stop. Do not weaken a protected Product/security/repository property merely to make CI green; retire accidental controls when they no longer protect a live failure class.

## MPC safety and Evidence specialization

- Sankhya and marketplace providers are external systems. Sankhya target access is through its sanctioned API Gateway; Direct Oracle/database is not a target fallback.
- Unknown, absent, partial, unavailable, or unproved facts never become plausible known defaults.
- Organization isolation fails closed across Organizations.
- Consequential external writes require explicit owner meaning, duplicate protection/idempotency where required, auditability, reconciliation, and no blind retry after ambiguous possible acceptance.
- Provider PII is minimized; secrets/PII never enter commits, logs, review transport, fixtures, or durable docs.
- Mercado Livre is the first proving marketplace. Provider DTO/protocol detail remains behind the consuming semantic boundary.
- Do not create an MPC Product/PIM master, generic Integration/Mutation/Workflow authority, provider-plugin platform, or AI-specific business bypass by convenience.
- Live Mercado Livre writes require explicit operator authorization.
- Mocks/fakes prove local behavior only. Claims about Mercado Livre, Sankhya, runtime, persistence, or browser behavior require Evidence proportional to the real dependency.

## MPC dependency / target-architecture specialization

Dependency or lockfile changes must be explicitly inside the declared acceptance increment, including why the dependency is required now and the proof that exercises it.

Never bend accepted target architecture merely to preserve removed/superseded code or tests. Preserve current security, PII, Organization-isolation, data-integrity, Product-contract, and irreversible/external-effect protections.

## Publication specialization

Use conventional-commit PR titles and state both what changes and what deliberately does not change. Never force-push or rewrite shared history. Explicit operator merge authorization remains required where `docs/roadmap.md` or the active increment requires it.

Dependent work does not stack on an unmerged acceptance increment by default. A roadmap stage may remain OPEN across several integrated increments.

Independent review lifecycle/transport is owned by the pinned `ADVERSARIAL-REVIEW-METHOD.md`. Reviewer output is Evidence, not Product/status authority.

## Retained derived production guide

[`evidence-grounded-production-engineering-for-llm-agents.md`](evidence-grounded-production-engineering-for-llm-agents.md) is an optional historical/portable **non-authoritative reference**. Current production realization uses the pinned `ROUTER.md` → `METHOD.md` + `REALIZATION-METHOD.md` profile. Keep the guide only while it has a real reference consumer; the aggregate gate does not require its mere existence.

## Justified Repository Standard path specializations

These are path/layout deviations only. They do not change authority semantics.

- `ARCHITECTURE.md` remains at repository root while dense current citations make a cosmetic move higher-cost than useful.
- `docs/architecture/decisions/` remains the ADR registry while current citations consume it.
- `docs/engineering/rebaseline/` remains the accepted D-stage authority package while D0–D8/D6-R2 are densely cross-linked there.
- `docs/engineering/rebaseline/EVIDENCE-REGISTER.md` remains with that package until a useful post-D9 consolidation.

When a specialization is removed, repair live consumers in the same acceptance increment and remove the old path; do not create indefinite compatibility authority.
