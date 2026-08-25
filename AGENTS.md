# Marketplace Central — Agent Bootstrap

## Fresh-session route

Before relying on chat, handoff, or remembered state:

1. revalidate repository identity, current branch/HEAD, remote `main`, relevant PR base/head, and aggregate CI;
2. read `AGENTS.md` → [`docs/index.md`](docs/index.md) → [`docs/roadmap.md`](docs/roadmap.md);
3. open the pinned methodology `ROUTER.md`, select only the Method profile required by the task;
4. load one or two repository task owners. The repository-local task authority pack defaults to **5 files or fewer**.

Do not recursively read phase history, ADRs, Evidence, Git history, research, removed runtime, or review dialogue before a current claim requires it.

## Canonical DevelopmentConexus methodology

Canonical repository:

`developmentconexus-ops/conexus-methodology`

Accepted immutable pin:

`9c7210d1504bef01c0d134a6c3ae8627deebb535`

Start methodology selection at:

`ROUTER.md`

Never consume floating methodology `main`. Do not copy the Method suite into this repository as independent authority. A methodology upgrade is a separate acceptance increment that explicitly moves this pin.

## Local authority

- Repository current authority outranks handoff/chat/history.
- [`docs/roadmap.md`](docs/roadmap.md) is the **sole mutable current stage/status/allowed-work/next-action authority**.
- [`docs/index.md`](docs/index.md) routes tasks only.
- `ARCHITECTURE.md`, accepted D-stage owners, and the ADR registry own their stated durable scopes.
- Review output, Evidence, code, tests, research, and Git history support decisions; they do not become Product/status authority merely by existing.
- Repository-local engineering/Git/CI/proof specialization lives in [`docs/development/engineering-rules.md`](docs/development/engineering-rules.md).

## MPC hard rails

Unless accepted authority explicitly reopens them:

- Sankhya and marketplace providers are external systems; Sankhya target access uses its sanctioned API Gateway, never Direct Oracle/database fallback.
- Unknown/absent/partial/unavailable facts never become plausible known defaults.
- Organization isolation fails closed.
- Consequential external writes require explicit owner meaning, duplicate protection where required, auditability/reconciliation, and no blind retry after ambiguous possible acceptance.
- Provider PII is minimized; secrets/PII do not enter commits, logs, transcripts, fixtures, or docs.
- Mercado Livre is the first proving marketplace.
- Do not invent an MPC Product/PIM master, generic Integration/Mutation/Workflow authority, provider-plugin platform, or AI-specific business bypass by convenience.
- Live Mercado Livre writes require explicit operator authorization.
- Product implementation begins only when [`docs/roadmap.md`](docs/roadmap.md) permits it.

## Repository lifecycle

Default to one branch + Draft PR per **acceptance increment**, not per whole roadmap stage. A stage may remain open after an increment is integrated. Do not stack dependent work on an unmerged prior increment by default. Never force-push or rewrite shared history. Squash merge is normal integration and requires explicit operator authorization where current authority requires it.

Temporary/review transport is branch-only and must not enter candidate/main. Independent review follows the pinned `ADVERSARIAL-REVIEW-METHOD.md`; local review mechanics are in `docs/development/engineering-rules.md`.

## Verification

```powershell
npm run gate
```

`scripts/gate.ps1` is the single repository aggregate gate. CI protects structural repository invariants and the canonical Product contract; planning prose, P8/P9 wording, walkthrough records, and historical review evidence are not permanent string-matched CI surfaces. A red current material control is still a stop, and claims about real providers/runtime/browser behavior require Evidence proportional to the real dependency.
