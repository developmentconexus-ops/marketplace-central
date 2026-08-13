# Marketplace Central — Architecture Rebaseline Status

> **THIS IS THE SESSION START HERE DOCUMENT.**  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Baseline being re-adjudicated:** `main@de1dc88bcef5a6ed5515378e7c646682c0bc15d2`  
> **Current gate:** `D0 — Documentation authority and current-state rebaseline`  
> **Implementation authorization:** **BLOCKED until D9 is accepted**  
> **Current review surface:** draft PR `#40` on `docs/architecture-rebaseline`

## Why this file exists

Any fresh session must be able to answer, without reading old plans or chat history:

1. Where are we?
2. What is already decided?
3. What remains open?
4. What is forbidden right now?
5. What exact gate comes next?

This file is the **only status/progress authority** for the architecture rebaseline. Do not create another roadmap, handoff, mission tree, status file or parallel checklist.

## Mandatory read order for architecture work

1. `/AGENTS.md`
2. `docs/architecture/rebaseline/README.md` — this file; current gate and next action
3. `docs/architecture/rebaseline/PROGRAM.md` — approved D0–D9 program and gate definitions
4. `/ARCHITECTURE.md` — current architectural constitution and constraints
5. `docs/engineering/standards/root-cause-global-maximum-method.md` — decision method
6. the document for the **current D-stage only**, when it exists
7. code/contracts/tests needed to verify the question at hand

Historical documents are Git history, not active authority.

## Current state

### D0 — Documentation authority and current-state rebaseline

**Status:** ACTIVE / operator-approved direction; branch cleanup in progress.

D0 establishes:

- one documentation hierarchy;
- the Root-Cause / Global-Maximum method;
- a current-state evidence register;
- disposition of legacy ADRs/plans/wiki/handoffs;
- a stable D0–D9 workflow for future sessions;
- explicit prohibition on jumping from folder architecture directly to implementation planning.

**D0 does not decide the final context list, schema, event catalog, API or frontend topology.** Those belong to D1–D7.

### Exact next action

Finish and review the D0 authority-cleanup PR. After D0 is accepted and merged, begin:

> **D1 — Context Adjudication:** determine the final business contexts and non-context components from domain semantics, lifecycle, state ownership, commands, contracts and named failure modes — not from the current folder tree or the previous 13-context proposal.

Do **not** write an implementation plan after D0.

## Program gates

| Gate | Question answered | Required output | Status |
|---|---|---|---|
| **D0** | What is the authority hierarchy, current state and design method? | canonical docs + legacy disposition + evidence register | **ACTIVE** |
| **D1** | What business contexts actually exist? | context map + admission/rejection rationale + legacy ownership map | NOT STARTED |
| **D2** | What are the identities and data authorities? | identity model + table/schema ownership + recoverability/reset decision | BLOCKED BY D1 |
| **D3** | How do internal components communicate? | sync-call map + event catalog + projection rules + transaction/outbox semantics | BLOCKED BY D1–D2 |
| **D4** | How do Mercado Livre and Sankhya/Oracle integrate? | capability matrix + auth/pagination/rate-limit/retry/write semantics + port ownership | BLOCKED BY D1–D2 |
| **D5** | What is the external HTTP contract? | API operation inventory/disposition + target OpenAPI + generation/validation decision | BLOCKED BY D1–D4 |
| **D6** | How does the frontend consume the system? | route/screen/API map + cache/query/mutation/error-state topology | BLOCKED BY D5 |
| **D7** | How does the runtime operate? | serving/worker/scheduler/outbox/transaction/deploy topology | BLOCKED BY D2–D4 |
| **D8** | Do the designs compose in real product flows? | golden-flow simulations with IDs, tables, contracts, events, retries and failures | BLOCKED BY D1–D7 |
| **D9** | Is the whole design a coherent global maximum? | adversarial contradiction review + accepted residuals + implementation DAG | BLOCKED BY D8 |
| **IMPLEMENT** | Can workers execute without inventing architecture? | implementation plans derived from accepted D1–D9 | BLOCKED BY D9 |

## Gate completion rule

A D-stage closes only when:

- its target property is explicit;
- relevant current-state evidence was measured;
- local-vs-global alternatives were considered;
- authorities and boundaries are unambiguous at that stage;
- downstream consequences were checked;
- no material contradiction remains;
- operator accepts the written result.

If a later gate discovers a material contradiction, reopen the smallest earlier gate whose decision is invalidated and record why here.

## Binding decisions already made

These are not reopened without a material finding or changed constraint:

- Root-cause/global-maximum analysis governs structural work.
- Simplify accidental complexity; do not simplify correctness.
- No `old/`, `legacy/` or other source-code cemetery. Git is the archive.
- Hard cutover is allowed because there are no production users requiring compatibility.
- Hard cutover does **not** authorize an indefinitely red/unverifiable `main`.
- Compatibility layers require a positive measured reason.
- Legacy modules/plans/ADRs are current-state/history evidence, not target authority.
- External marketplace protocol knowledge belongs behind vendor adapters implementing consumer-owned ports.
- Do not introduce Kafka/NATS/RabbitMQ or another broker without a measured need that Postgres/outbox cannot adequately meet.
- A clean database baseline is allowed only after D2 proves recoverability of all state that would be discarded.
- No implementation plan is authorized before D9.

## Session handoff rule

At the end of every D-stage or material replan, update **this file in the same PR** with:

- gate status;
- accepted decisions;
- reopened decisions;
- exact next action;
- implementation authorization state.

A session that cannot reconstruct the next action from this file is evidence that the documentation topology has regressed.