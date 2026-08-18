# Marketplace Central — Agent Bootstrap

## Read order

Before proposing or changing anything, read:

1. `AGENTS.md`
2. `docs/engineering/rebaseline/README.md`
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/engineering/rebaseline/DECISION-RECONCILIATION-BASELINE.md`
6. `docs/architecture/decisions/README.md`
7. accepted/current D-stage artifact(s) named by the router
8. supporting evidence needed for the specific decision

`docs/engineering/rebaseline/README.md` is the **sole authority for current stage, status, allowed/blocked work and exact next action**. Never infer current authority from memory, Git history, retired plans, review dialogue, stale candidates or existing code shape.

The Decision Reconciliation Baseline is an always-read routing map for which decision generation is current. It is not a second semantic architecture. Detailed meaning remains in `ARCHITECTURE.md` and accepted D-stage artifacts. ADR file status/disposition belongs only to the ADR registry.

## DevelopmentConexus Engineering Method

Canonical source: `developmentconexus-ops/conexus-methodology/METHOD.md`  
Local consumed version: **1.0.0**  
Local availability copy: `docs/engineering/standards/root-cause-global-maximum-method.md`

The local file is a byte-for-byte context copy, **not a fork or second authority**. Replace it manually from the canonical source only when an operator-approved methodology update is adopted. Do not add automatic sync, submodules, packages, bots, registries or distribution machinery without a proven failure class.

This repo may specialize or operationalize the organizational method, but must never silently redefine or weaken it. Surface any conflict inside the method's scope. The D0–D9 Architecture Rebaseline lifecycle is repo-specific specialization, not a second organizational engineering method.

## Fable independent review

For material independent review with Fable, after reconstructing this repository's authority/read order, follow the canonical Standard Fable review workflow in `developmentconexus-ops/conexus-methodology/README.md`.

Repository authority remains local; Fable review is non-authoritative input until operator ratification and canonical filing.

## Target-design authority by scope

- Router: current program status and next action.
- `ARCHITECTURE.md`: stable cross-stage constraints.
- Decision Reconciliation Baseline: current decision-generation routing.
- ADR registry: sole ADR file status/disposition authority.
- Accepted D-stage artifacts: detailed semantics in stage scope.
- Code, schemas, APIs, tests and runtime: current-state evidence only unless accepted authority explicitly rehomes meaning.

Surface material conflicts; never resolve them by silent reinterpretation.

## MPC safety rails

Unless explicitly reopened by accepted architecture:

- Go backend is canonical business execution; React is a client, not a second authority.
- Sankhya and marketplace providers are external systems.
- Sankhya target integration uses its sanctioned API Gateway; Direct Oracle/database is not target fallback.
- Provider-specific DTO/protocol knowledge stays behind integration boundaries.
- Unknown/absent facts never become plausible known defaults.
- External writes require explicit owner meaning, duplicate protection, auditability and reconciliation; ambiguous potentially accepted writes are not blindly retried.
- Provider PII is minimized.
- Tenant-ready Organization isolation remains a real invariant.
- Mercado Livre is the first marketplace used to prove the operating loop.
- Implementation may not recreate an MPC Product/PIM master, generic Integration/Mutation/Workflow authority or AI-specific business bypass.

More specific product semantics belong to accepted D-stage artifacts.

## Repository workflow

- One branch per change.
- Use conventional-commit PR titles accepted by repository governance.
- Never push without explicit operator permission; never merge without explicit operator authorization.
- Declare what a PR changes and what it deliberately does not change.
- Cold-review the intended property before merge.
- Do not reset, revert, stash, clean or delete working state you do not own.

## Verification

```powershell
npm run gate
npm run gate:full
```

`scripts/gate.ps1` is the shared local/CI gate implementation. Do not weaken verification to make it pass. Claims depending on external systems require the appropriate real-environment evidence.

## Operational safety and documentation

- Never expose secrets or PII in logs, transcripts, commits or docs.
- Dependency changes require explicit scope.
- Run Go commands from `apps/server_core`.
- Live Mercado Livre writes require explicit operator authorization.
- Do not alter target architecture merely to preserve legacy code/tests.
- Git history is the archive; do not create parallel roadmaps, permanent handoff trees or active `old/` copies.
- Durable architecture belongs in accepted D-stage artifacts, target ADRs, `ARCHITECTURE.md` or the reconciliation routing baseline; supporting evidence must remain distinguishable from target authority.
