# Marketplace Central — Agent Bootstrap

## Current program

Marketplace Central is in **Architecture Rebaseline / Technical System Design**. Product implementation is intentionally blocked until the D0–D9 design program is completed and accepted.

### Immediate operating boundary

The repository is currently finishing **documentary / governance authority cleanup** on PR #41 before deeper software-design discussion resumes.

Until that cleanup is explicitly closed:

- change only documentation, documentation authority, governance metadata/schemas, gates/workflows/scripts where they consume retired documentary authority, and proof needed to verify that cleanup;
- do **not** redesign, refactor, delete or choose a target for legacy product code;
- do **not** use the cleanup as an excuse to settle module/context, persistence, API, frontend, runtime or integration architecture;
- product/source findings discovered while inspecting cleanup are evidence only and are carried into the relevant D-stage;
- legacy source disposition is adjudicated **stage by stage across D0–D9** and implemented only after the corresponding architecture/cutover decision is accepted. It is not part of the current documentary cleanup.

The current cleanup ends when retired documentary authorities and their active consumers are removed/retargeted, current governance is self-contained, verification is green without weakening controls or raising ratchets, and a fresh session can identify one authority path and one exact next action.

Start every session in this order:

1. `AGENTS.md`
2. `docs/engineering/rebaseline/README.md` — current stage, status, exact next action
3. `docs/engineering/rebaseline/TMP-SESSION-HANDOFF.md` **while the documentary cleanup remains open** — continuity only, never design authority
4. `docs/engineering/standards/root-cause-global-maximum-method.md`
5. `ARCHITECTURE.md` — only stable product-level constraints
6. `docs/architecture/decisions/README.md` — ADR registry and current/reopened status
7. the document for the active D-stage, if separate from the rebaseline README
8. code/contracts/runtime evidence needed for that stage

Do **not** reconstruct the roadmap from Git history, deleted plans, old handoffs or memory. Git history is evidence only when the current rebaseline explicitly asks for historical evidence.

## Rebaseline gate

The governing sequence is:

`documentary authority cleanup → D0 current state → D1 contexts → D2 identity/data → D3 communication/events → D4 external integrations → D5 API → D6 frontend → D7 runtime/transactions/outbox → D8 golden flows → D9 adversarial global-maximum review → implementation DAG → implementation plan → implementation`

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
2. `ARCHITECTURE.md` stable product-level constraints;
3. accepted D-stage design outputs;
4. OpenAPI/governance/runtime code as current-state evidence for what exists today;
5. tests/builds/commits as execution evidence.

Current code shape is evidence about the present system, **not an argument that the target must preserve that shape**. A prior ADR that `docs/architecture/decisions/README.md` marks `reopened by ADR-035` is historical evidence, not target authority until the named D-stage re-adjudicates it.

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

The current gate contains legacy ratchets anchored to `internal/modules`. During D0/D1 they are evidence and transitional controls; do not mistake their current scope for target architecture.

## Architecture safety rules that remain binding during the rebaseline

- Mercado Livre is the first operational marketplace control plane.
- Oracle/Sankhya is an external source behind MPC-owned adapter boundaries; raw Oracle/driver knowledge does not belong in business contexts.
- Marketplace provider integration enters through `internal/adapters/marketplace/<vendor>` and implements ports owned by consuming contexts (ADR-033).
- Provider wire DTOs and provider-specific protocol knowledge remain inside their adapter tree.
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
- Stage evidence/design lives in the current D-stage document.
- Supporting references must be explicitly labeled non-authoritative.
- A superseded roadmap/spec/handoff is deleted after any still-valid decision has been absorbed into the current authority.
- `TMP-SESSION-HANDOFF.md` is allowed only while the current documentary cleanup is unfinished; delete it after the canonical topology passes the fresh-session test.

A new session should never need to decide which of several roadmaps is current.