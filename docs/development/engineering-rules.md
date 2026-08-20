# Marketplace Central — Engineering Rules

> **Scope:** repository-local execution, Git, CI, proof and safety specialization. This file does not own Product semantics or current program status.

## Canonical organizational authority

Use, without copying into this repository:

- `developmentconexus-ops/conexus-methodology/METHOD.md` — **DevelopmentConexus Engineering Method v1.0.0**;
- `developmentconexus-ops/conexus-methodology/REPOSITORY-STANDARD.md` — **DevelopmentConexus Repository Standard v1.0.0**.

The Method governs engineering reasoning. The Repository Standard governs the operating envelope. MPC rules may strengthen them for real local constraints and must not silently weaken them.

## Derived production engineering guide

[`evidence-grounded-production-engineering-for-llm-agents.md`](evidence-grounded-production-engineering-for-llm-agents.md) is a restored **derived / portable / non-authoritative** guide for production coding, technology research, dependency/framework evaluation, source hierarchy and proof. It operationalizes the canonical Method; it does not replace or amend it.

The restored file is preserved from historical blob `8de8ff4afbfcc2ee37a7db6ea0019e717740ebcf`. Its original relative parent-Method link referred to the retired local availability copy. The current parent authority is `developmentconexus-ops/conexus-methodology/METHOD.md` v1.0.0 as declared above.

## Local verification

```powershell
npm run gate
npm run gate:full
```

`scripts/gate.ps1` is the shared local/CI implementation. A red current control is a stop; do not raise a budget, weaken a guard, or add an exception merely to make it green.

For PR/diff claims, proof is over the intended `base...candidate` range. A clean checkout with an empty diff is not PR-diff proof.

## Safety and evidence boundaries

- Sankhya and marketplace providers are external systems. Sankhya target access is through its sanctioned API Gateway; Direct Oracle/database is not a target fallback.
- Unknown, absent, partial, unavailable, or unproved facts never become plausible known defaults.
- Organization isolation fails closed across Organizations.
- Consequential external writes require explicit owner meaning, duplicate protection/idempotency where required, auditability, reconciliation, and no blind retry of an ambiguous potentially accepted write.
- Provider PII is minimized; never expose secrets or PII in commits, logs, transcripts, fixtures, or documentation.
- Mercado Livre is the first proving marketplace. Provider DTO/protocol detail stays behind the consuming semantic boundary.
- Do not create an MPC Product/PIM master, generic Integration/Mutation/Workflow authority, provider plugin platform, or AI-specific business bypass by convenience.
- Live Mercado Livre writes require explicit operator authorization. Mocks/fakes prove local behavior only; claims about real external systems require evidence proportional to the real dependency.

## Dependency and target-architecture rule

Dependency or lockfile change must be explicitly inside the declared PR scope, including why the dependency is necessary now and the verification that exercises it.

Never bend accepted target architecture merely to preserve removed or superseded code/tests. A check that exists only for a retired subject may be removed only when attributable evidence proves either:

```text
subject population = 0
```

or

```text
full replacement coverage = proved
```

Narrative alone is not evidence. Preserve current security, PII, Organization-isolation, data-integrity, Product-contract and irreversible/external-effect protections.

## Branch, PR and review isolation

Normal lifecycle:

```text
main
→ one coherent candidate branch / one Draft PR
→ verification and bounded corrections
→ explicit operator merge authorization
→ squash merge
→ delete head branch
```

Do not stack a later stage on an unmerged earlier stage by default. Do not force-push or rewrite shared history.

Canonical future stage candidate naming:

```text
stage/<d-stage-or-gate>
```

For a material independent Fable review, create from the exact candidate:

```text
review/<stage>-fable
```

The review branch may add only:

```text
docs/work/current/ai-dialog.md
```

The candidate and `main` never contain `docs/work/**` or `ai-dialog.md`. The gate must prove the candidate ref exists and the review branch differs from that exact candidate by only the admitted review file. Reviewer output is Evidence, not authority; Round 2 occurs only for a surviving material contradiction.

## Documentation and context

Fresh route:

```text
AGENTS.md
→ docs/index.md
→ docs/roadmap.md
→ 1–2 owning documents
```

Default task pack is at most five files. Do not recursively read phase history, ADRs, Evidence, research or Git history before a concrete task requires them.

Durable documentation must not depend on `docs/work/**`. Temporary candidates, plans, handoffs and review dialogue are absorbed or deleted before merge. Do not create `docs/superpowers/`, permanent handoff/dialogue/round trees, active `old/` or archive trees, or parallel roadmaps.

## Justified Repository Standard path deviations

These are path/layout deviations only. They do not change authority or semantics.

### Stable architecture root

- **Standard preference:** `docs/architecture/index.md`.
- **Local deviation:** keep `ARCHITECTURE.md` at repository root for now.
- **Reason:** multiple accepted D-stage/ADR documents cite the established root path; moving it now creates broad authority-document churn with no task-pack reduction.
- **Removal trigger:** rehome when a material architecture edit already requires coordinated citation updates, or before the next structural architecture generation if the move becomes low-churn.

### ADR registry

- **Standard preference:** `docs/decisions/`.
- **Local deviation:** keep `docs/architecture/decisions/` while the surviving ADR residues retain explicit D7/transition retirement conditions.
- **Reason:** the registry and citation subtree are internally linked and referenced by accepted phase authorities; a cosmetic move would touch many durable records without improving current routing.
- **Removal trigger:** rehome when the surviving residue set is materially revised/retired or a new current decision generation requires registry restructuring.

### D-stage authorities

- **Standard preference:** `docs/phases/`.
- **Local deviation:** keep `docs/engineering/rebaseline/` through current D5.
- **Reason:** D0–D5 documents form a densely cross-linked accepted authority package; mass movement during an operating-envelope PR risks path churn without semantic/context benefit.
- **Removal trigger:** after D5 closes, reassess once the canonical OAD absorbs part of D5's machine-readable meaning; otherwise rehome before D6 only if the move is mechanical and low-risk.

### Evidence Register

- **Standard preference:** `docs/evidence/`.
- **Local deviation:** keep `docs/engineering/rebaseline/EVIDENCE-REGISTER.md` while it remains the baseline-attributed supporting record for the active rebaseline package.
- **Reason:** current phase documents cite it in place and it is opened only on demand; moving it now provides no default-context reduction.
- **Removal trigger:** when the rebaseline evidence package is consolidated after D5/D9 or a new current Evidence home is created for a real consumer.

No compatibility pointer is created for any future move. When a deviation is removed, repair live consumers in the same PR and delete the old path.
