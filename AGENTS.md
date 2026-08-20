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

## Target architecture versus staged delivery

**Target architecture is not the first-release implementation scope.** D0→D9 defines durable meaning, boundaries, invariants and evolution constraints; it does not require every admitted Product capability, infrastructure mechanism or future-scale concern to be implemented before the first useful internal slice can run.

For implementation planning after D9, classify material work explicitly as:

```text
BUILD NOW
→ required by the current proven consumer, current golden flow or a safety/correctness invariant

SEAM NOW
→ preserve the right ownership/boundary/extension point now, but do not build the future capability behind it yet

PROVE FIRST
→ run the smallest bounded spike needed to decide whether a mechanism or dependency is actually required

DEFER
→ no current consumer, failure class or irreversible constraint justifies implementation yet; record the reopen trigger
```

This classification is delivery guidance, not a second architecture authority or roadmap.

Preserve from the first real implementation any property whose later retrofit would threaten correctness, security, data meaning or ownership, including proportionately:

- canonical owner and data meaning;
- Organization isolation and cross-Organization fail-closed behavior;
- source-qualified external identity;
- authentication/authorization boundaries;
- migrations for persistent state;
- secrets outside source, logs and Product schemas;
- consequential idempotency/reconciliation where the accepted contract requires it;
- audit/recovery evidence required by accepted authority;
- provider/business-system adapter boundaries;
- one machine-readable Product wire authority.

Do **not** build speculative scale or platform machinery merely because the target architecture leaves room for it. A single current Organization, a small internal user set, one marketplace/provider, one deployment node or local execution does not justify self-service tenancy, generic connector/plugin frameworks, Kafka, Kubernetes, service mesh, multi-region, sharding, generic workflow engines or equivalent machinery without a real consumer/failure class.

Internal use also does not justify a shortcut that corrupts data, weakens tenant/security boundaries, leaks secrets, duplicates irreversible external effects or destroys future ownership. Production readiness is proportional to the current blast radius and consumer, while accepted correctness invariants remain binding.

When implementation opens, prefer **vertical slices** that make one real end-to-end flow work and become provable before broad horizontal construction. Avoid implementing every repository/service/endpoint layer first and discovering only later whether a useful flow works.

Reclassify BUILD NOW / SEAM NOW / PROVE FIRST / DEFER only on real evidence such as a new consumer, changed requirement, measured scale, newly reachable failure mode, provider change or implementation proof that invalidates an assumption.

This section operationalizes the Method's smallest-sustainable-solution, YAGNI and “prepare the seam, not the entire future capability” rules; it does not replace or amend the organizational Method.

## Fable independent review

### Operational identity

In this repository, **Fable is the role/name for the operator-run Claude Code session used as an independent second-model reviewer**. Fable is not a repository bot, background service, architecture authority or continuation of GPT's private chat context. The concrete Claude model/version may change; the reviewer role and authority contract do not.

GPT/lead and Fable communicate durably through the repository, normally through `AI-DIALOG.md`. The operator manually transfers a compact bootstrap prompt between the GPT session and Claude Code. Chat text is never authority; both sessions reconstruct repository authority independently.

### Purpose and authority

Fable exists to challenge material work adversarially from another model and perspective: find omitted failure classes, invalid assumptions, hidden or duplicate authority, security/concurrency/recovery gaps, YAGNI or overengineering, foreseeable retrofit traps, credible alternatives and a better Global Maximum where one exists.

Fable is not used for agreement theater, routine restatement, implementation assistance by default or ceremonial micro-review. Its output is **non-authoritative evidence**. GPT adjudicates every material finding against current repository authority; reviewer severity creates no requirement. Round 2 occurs only when a real material contradiction survives. Operator ratification and canonical filing are required before review conclusions become authority or program status.

After reconstructing this repository's authority/read order, follow the canonical **Standard Fable review workflow** in `developmentconexus-ops/conexus-methodology/README.md`. That canonical workflow remains the process authority; this section records only the Marketplace Central operating identity and handoff discipline.

### Handoff and input-token discipline

Before composing a Fable handoff, the lead must inspect the current router, candidate and active `AI-DIALOG.md` cycle and send only the minimum context Fable cannot safely discover from the repository.

A normal operator prompt contains only:

- repository, branch and expected HEAD as untrusted bootstrap references;
- an instruction to fetch and revalidate remote HEAD before reasoning or writing;
- an instruction to read `AGENTS.md`, follow the router and execute the active review cycle in `AI-DIALOG.md`;
- the target/review objective in one bounded sentence when it is not already unmistakable in the active cycle;
- exact authorized write/publish scope;
- the required final `HANDOFF → GPT` result.

Prefer repository paths and refs over pasted contents. Do **not** paste the read order, authority tree, candidate body, canonical artifacts, prior review dialogue, long architecture recap or detailed attack checklist when those already exist in `AGENTS.md`, the router, the candidate or `AI-DIALOG.md`. Include extra context only when it is not repository-discoverable, is materially ambiguous, or is required for safety/authorization.

The durable task-specific detail belongs in the bounded active cycle inside `AI-DIALOG.md`; the operator's Claude Code prompt should normally point to that cycle rather than duplicate it. One coherent package is reviewed once. Do not split related findings into repeated prompts merely to obtain more review turns.

Compact prompt shape:

```text
Fable/Claude Code: use developmentconexus-ops/marketplace-central,
branch <branch>; expected HEAD <sha> is reference only.
Fetch and revalidate remote HEAD. Read AGENTS.md and follow the router.
Execute the active independent-review cycle in AI-DIALOG.md against the
named candidate/authorities. Modify only the explicitly authorized review
channel, commit + push when authorized, verify remote, and end HANDOFF → GPT.
```

### Expected reviewer behavior

Fable must:

- reconstruct authority independently before reading the lead conclusion as persuasive context;
- review the smallest coherent package as one system;
- attack root cause, invariants, boundaries, authority, alternatives, YAGNI, security, concurrency, recovery and hardest foreseeable change;
- inspect current primary specifications, provider/framework documentation, source or reference implementations only when they materially improve the claim;
- report material findings only, with evidence, root cause, corrected invariant/direction, credible alternatives, Global versus Local Maximum, essential versus accidental complexity, reopen trigger and `APPROVE / REVISE / REJECT` disposition;
- distinguish what survived attack from what was merely not examined, without restating the entire candidate;
- avoid selecting later-stage implementation mechanics unless the current semantic contract is already incorrect without that decision;
- modify only the explicitly authorized review channel, preserve the other agent's turns, and verify the published remote commit when publication is authorized.

When Fable returns, GPT revalidates the branch/HEAD, confronts every material finding technically, records explicit agreements/disagreements, requests Round 2 only for a surviving material contradiction, and asks the operator for ratification only after the package has genuinely converged.

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