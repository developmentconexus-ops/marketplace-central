# Marketplace Central — Agent Bootstrap

## Start here

**DEFAULT READ: `AGENTS.md` → `docs/README.md`. STOP.**

Do not recursively read documentation, all D-stages, Git history, legacy code, review dialogue or evidence before a concrete task requires them. `docs/README.md` is the **sole authority for current program status, allowed/blocked work, exact next action and selective read routes**.

Before reasoning from or writing to the repository:

1. revalidate the repository, branch and remote HEAD;
2. read the two-file default bootstrap above;
3. open only the additional authority/evidence named by `docs/README.md` for the specific question;
4. surface material conflicts instead of silently choosing an authority.

## Authority by scope

- Operator-ratified decisions recorded in accepted repository authority outrank drafts, chat and reviewer claims.
- `docs/README.md` owns current program status and routes readers; it does not replace detailed semantics.
- `ARCHITECTURE.md` owns stable cross-stage constraints.
- Accepted D-stage artifacts own detailed semantics in their stage scope.
- `docs/architecture/decisions/README.md` alone owns ADR file status/disposition.
- Code, schemas, APIs, tests and runtime are current-state evidence unless accepted authority explicitly carries their meaning.
- Git history is the archive. Retired plans, candidates, handoffs and review dialogue do not regain authority because they remain discoverable.

## Engineering method and staged delivery

For material architecture, design, implementation, refactoring, debugging or review, use `docs/engineering/standards/root-cause-global-maximum-method.md`. Apply it proportionately; trivial work remains trivial.

**Target architecture is not initial delivery scope.** After D9, classify implementation work as:

```text
BUILD NOW   current consumer/golden flow or correctness requirement
SEAM NOW    preserve the right boundary now; build the future capability later
PROVE FIRST run the smallest bounded spike before committing a mechanism
DEFER       no current consumer/failure class; record the reopen trigger
```

Preserve from the first real slice any property whose retrofit would threaten correctness, security, data meaning, ownership or irreversible external effects. Do not build speculative scale/platform machinery. Prefer proven end-to-end vertical slices over broad horizontal construction.

## Independent Fable review

Fable is an operator-run independent challenger, never architecture authority. Use Fable only after a D-stage candidate is consolidated, or for an operator-explicit exceptional cross-cutting package. Follow the canonical **Standard Fable review workflow** in `developmentconexus-ops/conexus-methodology/README.md`.

`AI-DIALOG.md` is temporary: create it only on the review branch for an active bounded cycle, absorb ratified outcomes into canonical authority, then delete it before merge. Reviewer severity creates no requirement; Round 2 occurs only for a surviving material contradiction.

## MPC safety rails

Unless explicitly reopened by accepted authority:

- Go backend is canonical business execution; React is a client, not a second authority.
- Sankhya and marketplace providers are external systems.
- Sankhya target integration uses its sanctioned API Gateway; Direct Oracle/database is not target fallback.
- Provider DTO/protocol knowledge stays behind integration boundaries owned by consuming semantics.
- Unknown, absent, partial or unavailable facts never become plausible known defaults.
- External writes require explicit owner meaning, duplicate protection, auditability and reconciliation; ambiguous potentially accepted writes are not blindly retried.
- Provider PII is minimized.
- Organization isolation is a real invariant and must fail closed across Organizations.
- Mercado Livre is the first marketplace used to prove the operating loop.
- Do not recreate an MPC Product/PIM master, generic Integration/Mutation/Workflow authority, provider plugin platform or AI-specific business bypass without an accepted reopen.

## Repository workflow

- One branch per coherent change.
- After the current checkpoint, use one PR per D-stage; do not open the next D-stage until the previous one is accepted and merged to `main`.
- Use a conventional-commit PR title and declare both what changes and what deliberately does not change.
- Dependency or lockfile changes require explicit declared scope.
- Temporary candidates, plans, review channels and handoffs must be absorbed or deleted before merge. Do not add `docs/superpowers/`, active `old/` trees, parallel roadmaps or permanent session handoffs.
- Cold-review the final diff against the intended property and scope before merge.
- Never push without explicit operator permission. Never merge without explicit operator authorization.
- Do not reset, revert, stash, clean, force-update or delete working state you do not own.
- Live Mercado Livre writes require explicit operator authorization.
- Product implementation remains blocked until `docs/README.md` records D9 as accepted.

## Verification

```powershell
npm run gate
npm run gate:full
```

`scripts/gate.ps1` is the shared local/CI gate implementation.

- Red is a stop for every current control; never raise a baseline merely to make it green.
- A legacy-only check may retire only after attributable current-tree evidence proves its subject population is zero or full replacement coverage; prefer ratchet-to-zero before deletion. Narrative alone is insufficient.
- Never silence current security, Organization-isolation, data-integrity, Product-contract, PII or irreversible external-effect controls merely to accelerate the rebaseline.
- Do not alter accepted target architecture merely to preserve legacy code or tests.
- Presence is not execution; a check must exercise the relevant subject and demonstrate its negative control where required.
- Mocks/fakes prove local contract behavior, not real integration behavior.
- Claims depending on Mercado Livre, Sankhya, PostgreSQL deployment or browser runtime require appropriate real-environment evidence.
- Run Go commands from `apps/server_core`.
- Never expose secrets or PII in logs, transcripts, commits or docs.
