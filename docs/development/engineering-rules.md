# Marketplace Central — Engineering Rules

> **Scope:** repository-local execution, Git, CI, proof, documentation lifecycle, and Marketplace safety specialization. This file does not own Product semantics, current program status, or the content of the adopted engineering/frontend methods.

## Adopted local methods

Use the repository-local copies without locally redefining them:

- [`engineering-method.md`](engineering-method.md) — DevelopmentConexus Engineering Method v1.0.0;
- [`frontend-product-experience-planning-method.md`](frontend-product-experience-planning-method.md) — Frontend Product Experience Planning Method v2.3.

`AGENTS.md` and `docs/index.md` route to these files directly. There is no external methodology router, pin, profile-selection step, file-count limit, or owner-count limit.

## Verification

```powershell
npm run gate
```

The required gate is intentionally one objective aggregate check. It protects only properties suited to mechanical verification: required operating files, one canonical Product OpenAPI entrypoint, unsafe workflow triggers, unresolved merge-conflict markers, the current implementation block on changed paths, and Product-contract proof when the changed surface can affect that claim.

Do not use CI to judge Global Maximum, evidence quality, architecture quality, UX quality, how many files were read, or whether a planning artifact used a preferred format. Those are engineering-review questions governed by the adopted methods.

Run proof proportionately to the claim. Cheap repository checks apply broadly; expensive Product generation/semantic proof belongs on Product/proof-input changes or when a targeted claim explicitly needs it. If reliable changed-surface detection is unavailable, fail safe by running the stronger proof.

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

`docs/roadmap.md` owns current stage/status/allowed-work/next-action. `docs/index.md` routes concrete questions to current owners.

Start with the smallest current owner set likely to answer the task. Expand into any Product, architecture, contract, Evidence, research, Git history, code, runtime, or external source that can materially change or falsify the conclusion. Whole-repository analysis is allowed when the task warrants it; permission to expand is not a requirement to preload historical chains.

### Active-tree lifecycle

Use this lifecycle after a decision/proof increment closes:

```text
current owner / current router / current method                 KEEP
current evidence with a named live consumer                    KEEP while live
accepted intermediate with unique surviving meaning            REHOME once, then retire
fully absorbed or superseded intermediate                      RETIRE
history                                                         Git
```

Rules:

- rehome surviving meaning into the legitimate current owner before retiring its source artifact;
- a plan, spec, finding, review, ratification or proof record does not remain a parallel semantic authority after its durable conclusion is accepted elsewhere;
- historical snapshots stay truthful; do not rewrite their old status merely to make them look current;
- Git history is the default archive. Do not create a `docs/archive/` tree merely to keep retired working material visible;
- `docs/index.md` and other routers point to current owners/current evidence, not comprehensive historical chains;
- if an intermediate artifact is still the only owner/evidence for a current obligation, keep it or rehome it first—filename or age is never deletion evidence.

Temporary working material should therefore disappear from the active tree after its accepted meaning, status and proof obligations have been rehomed or closed. Durable current conclusions belong in their real Product/architecture/contract/repository owner.
