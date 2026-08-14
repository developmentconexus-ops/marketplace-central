# TEMP — Marketplace Central Rebaseline Session Handoff

> **Temporary continuity document. Not design authority.**
>
> `AGENTS.md` and `docs/engineering/rebaseline/README.md` always win. Delete this file as soon as the documentary/governance cleanup passes the fresh-session authority test.

## Current state

Marketplace Central is **not in product implementation** and **not yet in D0 design**.

The active review surface is **PR #41**, branch `docs/architecture-rebaseline-clean`.

PR #40 is superseded and must not be merged.

The current task is deliberately narrow:

> **finish documentary / governance authority cleanup, then stop cleaning.**

## What cleanup means

In scope:

- retired documents and dead references;
- stale ownership/milestone/mission language that still acts as authority;
- governance registries that route agents/tools to retired authority;
- gates/workflows/scripts only where they consume retired documentary authority;
- auxiliary tools only where they recreate retired documentation trees;
- proof that current governance is self-contained and fresh-session routing is unambiguous.

Out of scope:

- choosing `modules` vs `contexts`;
- deleting/refactoring legacy business code;
- deciding target schema, API, frontend, auth, runtime, jobs, events or integration architecture;
- performing a comprehensive codebase census before architecture discussion.

When cleanup inspection reveals product/runtime facts, record them as **supporting evidence** and leave the software untouched unless the only change required is to stop consumption/recreation of retired documentation authority.

## Cleanup line of finish

The cleanup is complete when:

1. no retired document competes as current architecture/program authority;
2. no active governance registry routes to retired authority;
3. gates/workflows/scripts no longer depend conceptually on retired documentary authority;
4. auxiliary tools no longer recreate retired documentation trees;
5. `AGENTS.md`, the canonical rebaseline README, `ARCHITECTURE.md`, ADR registry and governance contracts are mutually coherent;
6. repository verification is green without weakening controls or inflating ratchets;
7. no material dead reference remains on the active authority path;
8. a fresh session can identify one authority path and one exact next action without this handoff or chat memory.

Then delete this file and mark documentary cleanup DONE.

## What happens next

After cleanup, begin architecture discussion with the operator at:

```text
D0 — Product / System Definition
  ↓
D1 — Domains / Boundaries
  ↓
D2 — Identity / Tenant / Data Ownership
  ↓
D3 — Communication / Events
  ↓
D4 — External Integrations
  ↓
D5 — API
  ↓
D6 — Frontend
  ↓
D7 — Runtime / Jobs / Transactions
  ↓
D8 — Golden Flows
  ↓
D9 — Adversarial Architecture Review
  ↓
Implementation DAG / Plan
  ↓
Implementation
```

Product implementation remains blocked until D9 is accepted.

## How legacy enters D0–D9

Legacy is **evidence, not authority and not an audit prerequisite**.

For each design question:

1. state the decision/problem;
2. identify only the evidence needed to decide it;
3. inspect relevant current code/schema/runtime and external sources;
4. compare credible alternatives and trade-offs;
5. make an explicit operator-approved decision;
6. record the resulting contract/artifact;
7. classify affected legacy only when that stage has established its target owner/cutover.

Possible later dispositions include KEEP, KEEP AS REFERENCE, REFACTOR, MIGRATE/MOVE, REPLACE and DELETE.

## Known cleanup finding still open

`apps/server_core/cmd/mlprobe` still references retired `docs/design/ML-API-QUERY-CATALOG.md` and writes evidence into `docs/design/evidence/ml-api`.

Retargeting those documentary references/output is in cleanup scope. Redesigning the probe's Mercado Livre/database behavior is not.

## Exact continuation action

Continue PR #41 consumer-side authority cleanup:

1. align ADR-035 and any remaining active authority with the cleanup-first / D0-product-definition sequence;
2. retire or reclassify the old `D0-current-state-and-authority.md` surface so it cannot act as a prerequisite roadmap;
3. retarget `mlprobe` away from `docs/design/...`;
4. inspect active governance/gates/workflows for remaining documentary-authority references;
5. run verification;
6. perform the fresh-session test;
7. delete this handoff;
8. mark cleanup DONE and stop.
