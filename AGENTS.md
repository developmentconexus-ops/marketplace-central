# Marketplace Central — Agent Bootstrap

## Fresh actor route

**DEFAULT READ: `AGENTS.md` → `docs/index.md` → `docs/roadmap.md` → 1–2 task owners. STOP.**

Normal work fits in **5 files or fewer**. Do not recursively read phase history, ADRs, Evidence, Git history, removed runtime, or review dialogue before a concrete task requires them.

Before work:

1. revalidate repository, branch and remote HEAD;
2. read the route above;
3. use `docs/index.md` to select the smallest owning authority pack;
4. surface material authority conflicts instead of silently choosing one.

## Local authority model

- Operator-ratified decisions recorded in accepted repository authority outrank drafts, chat and reviewer claims.
- `docs/roadmap.md` is the **sole mutable current-program stage/status/allowed-work/next-action authority**.
- `docs/index.md` routes tasks; it owns no mutable program status.
- `ARCHITECTURE.md` owns stable cross-stage constraints under the documented local path deviation.
- Accepted D-stage artifacts own detailed semantics in their stage scope.
- `docs/architecture/decisions/README.md` owns ADR disposition/retirement state under the documented local path deviation.
- Evidence and Git history support decisions; they do not become Product authority by existence.

## Organizational method and repository standard

Canonical organizational authorities:

- `developmentconexus-ops/conexus-methodology/METHOD.md` — **DevelopmentConexus Engineering Method v1.0.0**;
- `developmentconexus-ops/conexus-methodology/REPOSITORY-STANDARD.md` — **DevelopmentConexus Repository Standard v1.0.0**.

Do not copy either into this repository as an independent authority. MPC-specific execution, Git, CI and proof rules live in [`docs/development/engineering-rules.md`](docs/development/engineering-rules.md).

## MPC safety rails

Unless explicitly reopened by accepted authority:

- Sankhya and marketplace providers are external systems; Sankhya target integration uses its sanctioned API Gateway, not Direct Oracle/database fallback.
- Provider DTO/protocol knowledge stays behind the consuming semantic boundary.
- Unknown, absent, partial or unavailable facts never become plausible known defaults.
- Organization isolation fails closed across Organizations.
- Consequential external writes require explicit owner meaning, duplicate protection/idempotency where required, auditability and reconciliation; ambiguous potentially accepted writes are not blindly retried.
- Provider PII is minimized; secrets/PII never enter commits, logs, transcripts or docs.
- Mercado Livre is the first marketplace used to prove the operating loop.
- Do not recreate an MPC Product/PIM master, generic Integration/Mutation/Workflow authority, provider plugin platform or AI-specific business bypass by convenience.
- Live Mercado Livre writes require explicit operator authorization.

## Repository-specific publication rules

- One branch/one Draft PR per coherent stage or gate by default; do not stack later stages on an unmerged earlier stage.
- Use a conventional-commit PR title and state both what changes and what deliberately does not change.
- Dependency/lockfile changes require explicit declared scope.
- Never push or merge outside operator authorization. Never force-push or rewrite shared history.
- Temporary plans, candidates, handoffs and review dialogue must be absorbed or deleted before merge.
- A material Fable review uses `review/<stage>-fable` and may differ from the exact candidate only by `docs/work/current/ai-dialog.md`; candidate and `main` contain no `docs/work/**`.
- Do not add `docs/superpowers/`, active `old/`/archive trees, parallel roadmaps or permanent session handoffs/dialogue.
- Do not begin Product implementation unless `docs/roadmap.md` permits it.

## Verification

```powershell
npm run gate
npm run gate:full
```

`scripts/gate.ps1` is the shared local/CI gate.

- Red is a stop for every current control; never weaken a property merely to make it green.
- Retire a check only with attributable `subject population = 0` or proved full replacement coverage.
- Do not warp target architecture to preserve removed/legacy code or tests.
- Presence is not execution; material guards require a deterministic falsifier/negative control.
- Mocks/fakes prove local behavior, not real integration behavior.
- Claims about Mercado Livre, Sankhya, runtime, persistence or browser behavior require evidence proportional to the real dependency.

## Material stop conditions

STOP/SPLIT only when repository-standard alignment would require changing accepted Product semantics, ownership/trust boundaries, a current safety invariant, the admitted Product API surface, or another operator-ratified architecture decision. Path aesthetics alone are never a reason to reopen D0–D5.
