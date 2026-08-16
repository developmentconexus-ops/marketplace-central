# Marketplace Central — Agent Bootstrap

## Read order

Before proposing or changing anything, read:

1. `AGENTS.md`
2. `docs/engineering/rebaseline/README.md`
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/architecture/decisions/README.md`
6. accepted/current D-stage artifact(s) named by the router
7. supporting evidence needed for the specific decision

`docs/engineering/rebaseline/README.md` is the **sole authority for current stage, status, allowed/blocked work and exact next action**. Never infer current authority from memory, Git history, retired plans or existing code shape.

## DevelopmentConexus Engineering Method

Canonical source: `developmentconexus-ops/conexus-methodology/METHOD.md`  
Local consumed version: **1.0.0**  
Local availability copy: `docs/engineering/standards/root-cause-global-maximum-method.md`

The local file is a byte-for-byte context copy, **not a fork or second authority**. Replace it manually from the canonical source only when an operator-approved methodology update is adopted. Do not add automatic sync, submodules, packages, bots, registries or distribution machinery without a proven failure class.

This repo may specialize or operationalize the organizational method, but must never silently redefine or weaken it. Surface any conflict inside the method's scope. The D0–D9 Architecture Rebaseline lifecycle is **repo-specific specialization, not a second organizational engineering method**.

## Target-design authority

For target design:

1. operator-approved accepted D-stage decisions indicated by the router;
2. stable constraints in `ARCHITECTURE.md`;
3. ADR authority from `docs/architecture/decisions/README.md`;
4. current code, schemas, APIs, tests and runtime as current-state evidence only.

A reopened ADR is historical evidence until re-adjudicated. Surface material conflicts; never resolve them by silent reinterpretation.

## MPC safety rails

Unless explicitly reopened by accepted architecture:

- Go backend is canonical business execution; React is a client, not a second authority.
- Sankhya/Oracle and marketplace providers are external systems.
- Provider-specific DTO/protocol knowledge stays behind integration boundaries.
- Unknown/absent facts never become plausible known defaults.
- External writes require explicit authority, auditability and reconciliation; ambiguous potentially accepted writes are not blindly retried.
- Provider PII is minimized.
- Tenant-ready isolation remains a real invariant.
- Mercado Livre is the first marketplace used to prove the operating loop.

More specific product semantics belong to accepted D-stage artifacts, not legacy modules.

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

`scripts/gate.ps1` is the shared local/CI gate implementation. Do not weaken verification to make it pass. Claims depending on Oracle, Mercado Livre, Postgres deployment or browser runtime require appropriate real-environment evidence.

## Operational safety and documentation

- Never expose secrets or PII in logs, transcripts, commits or docs.
- Dependency changes require explicit scope.
- Run Go commands from `apps/server_core`.
- Live Mercado Livre writes require explicit operator authorization.
- Do not alter target architecture merely to preserve legacy code/tests.
- Git history is the archive; do not create parallel roadmaps, permanent handoff trees or active `old/` copies.
- Durable architecture belongs in accepted D-stage artifacts, ADRs or `ARCHITECTURE.md`; supporting evidence must remain distinguishable from target authority.
