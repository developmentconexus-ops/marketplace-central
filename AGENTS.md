# Marketplace Central — Agent Instructions

## Start here

Before proposing or changing anything, read in this order:

1. `AGENTS.md`
2. `docs/engineering/rebaseline/README.md`
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/architecture/decisions/README.md`
6. accepted artifact(s) for the active D-stage
7. supporting evidence needed for the specific decision

`docs/engineering/rebaseline/README.md` is the **sole authority for current stage, status, allowed/blocked work and exact next action**. Never infer current authority from memory, Git history, retired plans or existing code shape.

## Engineering method

All non-trivial architecture, design, implementation, refactoring, debugging and review **must follow** `docs/engineering/standards/root-cause-global-maximum-method.md`.

Do not duplicate or locally reinterpret the method here.

## Target-design authority

For target design, use this order:

1. operator-approved decisions recorded in the active design program;
2. stable constraints in `ARCHITECTURE.md`;
3. accepted active/prior D-stage artifacts;
4. accepted ADRs not explicitly reopened;
5. current code, schemas, APIs, tests and runtime as current-state evidence.

Existing implementation is evidence, not target authority. A reopened ADR is historical evidence until its owning stage re-adjudicates it. Surface material conflicts explicitly; never silently choose one authority.

## MPC safety rails

Unless explicitly reopened by accepted architecture:

- Go backend is the canonical business execution authority; React is a client, not a second authority.
- Sankhya/Oracle and marketplace providers are external systems.
- Provider-specific DTO/protocol knowledge stays behind integration boundaries.
- Unknown/absent facts must not become plausible known defaults.
- External writes require explicit authority, auditability and reconciliation; ambiguous potentially accepted writes are not blindly retried.
- Provider PII is minimized.
- Tenant-ready isolation remains a real invariant.
- Mercado Livre is the first marketplace used to prove the operating loop.

More specific product semantics belong to accepted D-stage artifacts; do not infer them from legacy modules.

## Repository workflow

- One branch per change.
- Use conventional-commit PR titles accepted by repository governance.
- Never push without explicit operator permission.
- Never merge without explicit operator authorization.
- Declare PR scope in both directions: what changes and what does not.
- Perform a cold review against the intended property before merge.
- Do not reset, revert, stash, clean or delete working state you do not own.

## Verification

Repository gates:

```powershell
npm run gate
npm run gate:full
```

`scripts/gate.ps1` is the shared local/CI gate implementation.

- Red is a stop.
- Presence is not execution; a check must exercise the relevant subject.
- Mocks/fakes prove local contract behavior, not real integration behavior.
- Claims depending on Oracle, Mercado Livre, Postgres deployment or browser runtime require appropriate real-environment evidence.
- Legacy ratchets are transitional evidence unless accepted design explicitly retains them.

## Operational safety

- Never expose secrets or PII in logs, transcripts, commits or docs.
- Dependency changes require explicit scope.
- Run Go commands from `apps/server_core`.
- Live Mercado Livre writes require explicit operator authorization.
- Do not alter target architecture merely to preserve legacy code or tests.

## Documentation discipline

- Git history is the archive; do not create parallel roadmaps, permanent handoff trees or active `old/` copies.
- `docs/engineering/rebaseline/README.md` is the sole current program router.
- Durable architecture belongs in accepted D-stage artifacts, ADRs or `ARCHITECTURE.md`.
- Supporting evidence must remain distinguishable from target authority.
- Remove superseded active documentation after its durable content is absorbed into canonical authority.

A fresh session should need one authority path, not judgment about which roadmap is current.
