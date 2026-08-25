# Marketplace Central — Engineering Rules

> **Scope:** repository-local execution, Git, CI, proof, and Marketplace safety specialization. This file does not own Product semantics, current program status, or the content of the adopted engineering/frontend methods.

## Adopted local methods

Use the repository-local copies without locally redefining them:

- [`engineering-method.md`](engineering-method.md) — DevelopmentConexus Engineering Method v1.0.0;
- [`frontend-product-experience-planning-method.md`](frontend-product-experience-planning-method.md) — Frontend Product Experience Planning Method v2.3.

`AGENTS.md` and `docs/index.md` route to these files directly. There is no external methodology router, pin, profile-selection step, file-count limit, or owner-count limit.

## Verification

```powershell
npm run gate
```

The required gate is intentionally one objective aggregate check. It protects only properties suited to mechanical verification: required operating files, one canonical Product OpenAPI entrypoint, unsafe workflow triggers, unresolved merge-conflict markers, the current implementation block on changed paths, and the canonical Product OAD proof.

Do not use CI to judge Global Maximum, evidence quality, architecture quality, UX quality, how many files were read, or whether a planning artifact used a preferred format. Those are engineering-review questions governed by the adopted methods.

For a claim that touches a specific wireframe, authorization slice, provider behavior, runtime property, or another protected surface, run the relevant targeted proof when the claim requires it. Do not make every historical proof a permanent prerequisite for unrelated planning work.

A mock/fake proves only the local mock boundary. Claims about Mercado Livre, Sankhya, runtime, persistence, browser behavior, or another real dependency require Evidence proportional to the real claim.

## Marketplace safety boundaries

- Sankhya and marketplace providers are external systems. Sankhya target access is through its sanctioned API Gateway; Direct Oracle/database is not a target fallback.
- Unknown, absent, partial, unavailable, or unproved facts never become plausible known defaults.
- Organization isolation fails closed across Organizations.
- Consequential external writes require explicit owner meaning, duplicate protection/idempotency where required, auditability, reconciliation, and no blind retry after an ambiguous potentially accepted write.
- Provider PII is minimized; secrets/PII never enter commits, logs, review material, fixtures, or durable documentation.
- Mercado Livre is the first proving marketplace. Provider DTO/protocol detail remains behind the consuming semantic boundary.
- Do not create an MPC Product/PIM master, generic Integration/Mutation/Workflow authority, provider-plugin platform, or AI-specific business bypass by convenience.
- Live Mercado Livre writes require explicit operator authorization.

## Dependency and target-architecture rule

A dependency or lockfile change must be inside the declared work scope, with a real current consumer and proof appropriate to the dependency claim.

Never bend accepted target architecture merely to preserve removed or superseded code/tests. Preserve current security, PII, Organization-isolation, data-integrity, Product-contract, and irreversible/external-effect protections.

## Git and publication

Preserve unowned state. Do not reset, clean, stash, force-update, force-push, or rewrite shared history by convenience.

Use the existing branch/PR when work is already in progress. Keep changes attributable and reviewable. Squash merge remains the normal integration mechanism where used by the repository, and merge requires explicit operator authorization.

Independent review is Evidence, not authority. Use it when material risk justifies it; do not create review ceremony merely because a template or prior workflow once required it.

## Documentation and context

`docs/roadmap.md` owns current stage/status/allowed-work/next-action. `docs/index.md` is a navigation aid.

There is no fixed task-pack size. Start where useful and expand into any Product, architecture, contract, Evidence, research, Git history, code, runtime, or external source that can materially change or falsify the conclusion. Whole-repository analysis is allowed when the task warrants it.

Temporary working material should not become a second authority. Durable accepted conclusions belong in their real owning document or Evidence home; do not create parallel roadmaps or permanent session handoff/dialogue trees by default.
