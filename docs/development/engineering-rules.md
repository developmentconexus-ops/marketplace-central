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
npm run gate:full
```

`scripts/gate.ps1` is the repository-local aggregate implementation. It must prove current repository properties over the intended base...candidate range, including bootstrap/context budget, routing/reachability, temporary/review hygiene, current Product proof, and deterministic negative controls where a reusable guard exists.

A red current control is a stop. Do not raise a budget, weaken a guard, or add an exception merely to make it green. Retire a subject-specific control only with attributable `subject population = 0` or proved replacement coverage.

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

Independent review lifecycle/transport is owned by the pinned `ADVERSARIAL-REVIEW-METHOD.md`; the repository aggregate gate enforces candidate/main hygiene and exact review isolation. Reviewer output is Evidence, not Product/status authority.

## Retained derived production guide

[`evidence-grounded-production-engineering-for-llm-agents.md`](evidence-grounded-production-engineering-for-llm-agents.md) is retained as historical/portable **non-authoritative reference** from blob `8de8ff4afbfcc2ee37a7db6ea0019e717740ebcf`.

It is not a current Method selector and its historical parent-link prose is not normative. Current production realization uses the pinned `ROUTER.md` → `METHOD.md` + `REALIZATION-METHOD.md` profile. Keep the guide only while it has a real reference consumer; do not update it into a second Method copy.

## Justified Repository Standard path specializations

These are path/layout deviations only. They do not change authority semantics.

### Stable architecture root

- **Organizational default:** `docs/architecture/index.md`.
- **MPC specialization:** retain `ARCHITECTURE.md` at repository root.
- **Consumer:** dense accepted D-stage/ADR citations to the established root path.
- **Why default is insufficient now:** moving it would create broad authority-document churn without reducing active context or changing correctness.
- **Removal trigger:** a material architecture-generation rewrite already touching those citations, or another concrete consumer making the move low-churn and valuable.

### ADR registry

- **Organizational default:** `docs/decisions/`.
- **MPC specialization:** retain `docs/architecture/decisions/` while current ADR residue remains linked there.
- **Consumer:** accepted ADR registry/citation subtree.
- **Removal trigger:** material ADR generation restructuring or retirement of the surviving residue set.

### Accepted D-stage authority package

- **Organizational default:** durable phase/result authority under `docs/phases/` when useful.
- **MPC specialization:** retain the densely cross-linked accepted package under `docs/engineering/rebaseline/`.
- **Consumer:** current D0–D8/D6-R2 authority, proofs, and routing.
- **Why default is insufficient now:** a cosmetic mass move would create high citation churn with no Product/context gain.
- **Removal trigger:** post-D9 consolidation or another material authority-generation change that makes rehoming part of useful work rather than path aesthetics.

### Evidence Register

- **Organizational default:** `docs/evidence/`.
- **MPC specialization:** retain `docs/engineering/rebaseline/EVIDENCE-REGISTER.md` while it supports the accepted rebaseline package.
- **Consumer:** current accepted phase authorities citing the register.
- **Removal trigger:** post-D9 evidence consolidation or creation of a new current Evidence home for a real consumer.

When a specialization is removed, repair live consumers in the same acceptance increment and remove the old path; do not create indefinite compatibility authority.
