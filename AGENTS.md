# Marketplace Central — Agent Bootstrap

## Current program

Marketplace Central is in **Architecture Rebaseline / Technical System Design**. Product implementation is intentionally blocked until the D0–D9 design program is completed and accepted.

### Immediate operating boundary

The **documentary / governance authority cleanup is closed** (PR #41). **D0 — Product / System Definition has not been opened with the operator yet.**

Until D0 is opened and produces an accepted decision:

- no software design change is authorized — not module/context layout, persistence, API, frontend, runtime, auth or integration architecture;
- do **not** redesign, refactor, delete or choose a target for legacy product code;
- product/source findings are evidence only and are carried into the relevant D-stage;
- legacy source disposition is adjudicated **stage by stage across D0–D9** and implemented only after the corresponding architecture/cutover decision is accepted.

Do not reopen the documentary cleanup to keep working. It ended on a measured condition: retired documentary authorities and their active consumers removed or retargeted, current governance self-contained, verification green without weakening controls or raising ratchets, and a fresh session able to identify one authority path and one exact next action. `docs/engineering/rebaseline/README.md` records how each of those was discharged.

Start every session in this order:

1. `AGENTS.md`
2. `docs/engineering/rebaseline/README.md` — current phase, status, exact next action
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md` — stable product/platform constraints only
5. `docs/architecture/decisions/README.md` — ADR registry and current/reopened status
6. the artifact for the active D-stage, once design begins
7. code/contracts/runtime evidence needed for the specific decision being made

Do **not** reconstruct the roadmap from Git history, deleted plans, old handoffs or memory. Git history is evidence only when the current rebaseline explicitly asks for historical evidence.

## Rebaseline gate

The governing sequence is:

`documentary authority cleanup → D0 product/system definition → D1 domains/boundaries → D2 identity/tenant/data ownership → D3 communication/events → D4 external integrations → D5 API → D6 frontend → D7 runtime/jobs/transactions → D8 golden flows → D9 adversarial architecture review → implementation DAG/plan → implementation`

Until D9 is accepted:

- do not start product-feature implementation merely because an old plan says it is next;
- do not create a new context/module, schema cutover, API redesign, frontend topology migration or legacy deletion unless the accepted design flow explicitly authorizes that change;
- documentary cleanup, measurement and proof tooling are allowed only when they directly support the current gate and do not smuggle in target architecture;
- an implementation plan for the product rebaseline is premature before D9.

The current status and exact next action live in one place only: `docs/engineering/rebaseline/README.md`.

## Binding engineering method

All non-trivial planning, implementation, refactoring and review follows `docs/engineering/standards/root-cause-global-maximum-method.md`.

Binding principle:

> Always simplify the code, never simplify correctness. Find the root cause, determine whether the proposed solution is only a local maximum, and prefer the global structure that makes the defect class unrepresentable or mechanically impossible at the strongest reasonable boundary.

Before choosing a solution, name:

- the observed symptom/evidence;
- the root cause;
- the target property/invariant;
- the authority and ownership boundary;
- credible local-maximum and global-maximum candidates;
- the strongest reasonable enforcement mechanism;
- the proof that will distinguish the fixed and broken worlds.

YAGNI removes speculative or redundant capability. It does not remove required invariants, fail-closed behavior or proof.

## Authority and conflict handling

For **current target design**, use this order:

1. operator-approved decisions recorded by the active rebaseline and current accepted ADRs;
2. `ARCHITECTURE.md` stable product/platform constraints;
3. accepted D-stage design outputs;
4. OpenAPI/governance/runtime code as current-state evidence for what exists today;
5. tests/builds/commits as execution evidence.

Current code shape is evidence about the present system, **not an argument that the target must preserve that shape**. A prior ADR that `docs/architecture/decisions/README.md` marks reopened is historical evidence, not target authority until the named D-stage re-adjudicates it.

When architecture, contract, runtime, ownership or verification evidence conflicts, stop the local conclusion and classify the conflict. Do not silently pick a side.

## How work lands

- One branch per change.
- Conventional-commit PR title; repository workflow enforces the allowed type set.
- **Never push without explicit operator permission and never merge without explicit operator authorization.**
- Work is tracked in GitHub issues/PRs when connector support permits it. If a connector blocks creation of a tracking issue, record that fact in the PR; do not invent an issue number.
- PR scope is declared in both directions: what changes and what does not.
- Before merge, perform a cold review against the declared target property and scope.

## Verification

`scripts/gate.ps1` is the single local/CI implementation of the gate lanes.

```powershell
npm run gate
npm run gate:full
```

Rules:

- CI and local verification invoke the same gate implementation.
- Red is a stop, not a note.
- Ratchet baselines may shrink and must never be raised merely to make a change green.
- Presence is not execution: a check must prove it actually ran the intended units.
- A mock/fake proves local contract behavior; it does not prove a real external integration.
- Claims about Oracle, Mercado Livre, Postgres deployment state or browser behavior require the appropriate real-environment proof when the claim depends on that environment.

Current gates may contain legacy ratchets anchored to current code layout. During the rebaseline they are transitional controls/evidence, not proof of target architecture. They are changed only when the D-stage that owns the affected software boundary makes that decision, except for references whose only purpose is retired documentary authority.

## Architecture safety rules that remain binding during the rebaseline

- Mercado Livre is the first operational marketplace control plane unless explicitly reopened by the accepted design process.
- Oracle/Sankhya is an external source behind MPC-owned adapter boundaries; raw Oracle/driver knowledge does not belong in business policy.
- Provider wire DTOs and provider-specific protocol knowledge remain at provider boundaries.
- Unknown business/operational facts never become plausible zero/default values.
- External writes are never blindly retried after an ambiguous outcome and must be auditable/reconcilable.
- Raw provider PII is not retained merely for convenience.
- Frontend is not a second business-logic authority.

Anything more specific than these may be under D0–D9 re-adjudication; check the ADR registry and rebaseline status before relying on an older rule.

## Operational rules

- One writer owns a checkout/shared seam at a time.
- Do not reset, revert, stash, clean or delete unknown working state.
- Do not expose secrets or PII in logs, transcripts, commits or documentation.
- Dependency changes are explicit scope, not incidental environment preparation.
- Go commands run from `apps/server_core`; when hand-running Go on Windows, `GOCACHE` and `GOMODCACHE` must resolve to absolute paths.
- Live Mercado Livre writes require explicit operator authorization.

## Documentation lifecycle

The active repository intentionally does **not** retain an `old/`, archive wiki, dated implementation-plan cemetery or permanent session-handoff tree.

- Git history is the archive.
- `docs/engineering/rebaseline/README.md` is the sole current progress/router document.
- Accepted durable decisions live in ADRs or `ARCHITECTURE.md`.
- Stage evidence/design lives in the current D-stage artifact.
- Supporting references must be explicitly labeled non-authoritative.
- A superseded roadmap/spec/handoff is deleted after any still-valid decision has been absorbed into the current authority.
- A temporary session-handoff file is allowed only while a cleanup is unfinished, and is deleted as soon as the canonical topology passes the fresh-session test on its own.

A new session should never need to decide which of several roadmaps is current.