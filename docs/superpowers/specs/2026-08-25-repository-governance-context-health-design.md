# Repository Governance & Context Health Rebaseline — Design

> **Status:** DESIGN APPROVED — operator-approved on 2026-08-25; implementation plan review required before repository changes
> **Trigger:** D6-R2 frontend work became slow/context-heavy after temporary methodology/router experiments and accumulated rebaseline evidence remained in the active tree
> **Method:** DevelopmentConexus Engineering Method v1.0.0
> **Baseline:** `main@bdbbef43ed3a5e9d912e67ddac5173024352eaa3`
> **Product semantics:** OUT OF SCOPE — preserve current Product/OAD meaning, 106 operations, 31 ordinary Permissions, H/A/S, runtime NONE

## 1. Decision

**Outcome: `RESTRUCTURE NOW` — restore Marketplace Central to a small active-authority graph, retire absorbed intermediate material to Git history, and make CI proportional to the changed claim.**

This is repository-governance/context-health correction, not another methodology framework.

Target operating path:

```text
AGENTS.md
→ docs/roadmap.md
→ applicable local method
→ docs/index.md when needed
→ smallest current owner(s)
→ evidence/Git history only when the concrete question requires it
```

The active tree answers **what is true now and where that truth lives**. Git history/evidence answers **how the repository arrived there**.

## 2. Root cause

Marketplace Central is effective at producing findings, reviews, ratifications, repairs and proofs, but has not consistently performed the final lifecycle step after a decision is absorbed:

```text
investigate
→ decide
→ ratify/prove
→ rehome into canonical owner
→ RETIRE intermediate material from active tree
```

This created three failure classes.

### 2.1 Context authority inflation

D6-R2 contains a large `NOTIF-01` / finding / ratification / repair / feed-forward / review chain. `D6-R2-AUTHORITY-ROUTE.md` intentionally keeps the chain reachable even when its surviving meaning has already been rehomed.

Observed consequences:

- agents encounter multiple historical snapshots while answering a current question;
- stale counts/status/method versions are repeatedly reconciled;
- task reads expand into project-history reconstruction;
- context grows without equivalent decision-quality gain.

Concrete examples:

- `D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md` still contains older 99-operation / 30-Permission and early-lock-state prose;
- `D6-R2-P8-BLOCK-LEDGER.md` still contains older B12/B10 status and Frontend Method v2.1 prose;
- `D6-R2-AUTHORITY-ROUTE.md` routes the complete historical chain despite roadmap owning current status.

### 2.2 Retirement drift

The ADR registry kept some legacy ADRs only until D7 adjudication. D7 is accepted/closed, but those entries still remain active residues.

### 2.3 Verification inflation

The single `npm run gate` currently invokes the heavy Product OAD proof for unrelated documentation changes. Recent docs-only evidence showed roughly 85 seconds spent in Product proof while Product files did not change.

Implementation-planning refinement found that the 95/99 replay is not pure history: it still protects valid current invariants. Therefore this health increment makes the heavy proof **diff-aware without weakening or splitting it**. Physical current-vs-history proof decomposition is deferred until a separate parity exercise can prove every current protection survives.

## 3. Governing invariants

### G1 — Active tree is current, not archival

Keep a document active only while it:

1. owns current Product/architecture/contract/repository semantics;
2. owns current mutable status/routing/method authority;
3. is evidence for a live unresolved/proof obligation; or
4. is explicitly required by a current owner to interpret current meaning.

Otherwise Git history is the archive.

### G2 — Rehome before retire

```text
identify surviving meaning
→ locate legitimate current owner
→ rehome surviving meaning once
→ update current routes/references
→ prove no live obligation depends on the intermediate file
→ retire it
```

### G3 — One meaning, one current owner

A finding, review, ratification, plan or proof does not remain parallel semantic authority after acceptance.

### G4 — Index routes current questions, not history

`docs/index.md` answers: **where is the smallest current authority likely to answer this question?** It is not a complete document catalog.

### G5 — Methods reason; repository owns state

Preserve local:

- `engineering-method.md`;
- `frontend-product-experience-planning-method.md`;
- `engineering-rules.md`.

Do not restore external methodology router/pin/profile machinery.

### G6 — Inquiry may expand; default context stays selective

There is no hard context/file budget. Start at current owners and expand only when extra evidence can materially change/falsify the conclusion.

### G7 — CI proves the changed claim proportionately

One required aggregate check remains. Cheap universal checks run always. The heavy existing Product proof runs only when Product/proof inputs change; when diff evidence is unavailable it runs fail-safe.

## 4. Target governance model

### Bootstrap

`AGENTS.md` remains small: revalidate repo/branch/PR, read roadmap, route to the applicable method, use index for current owners, preserve Marketplace/Git safety rails.

### Roadmap

`docs/roadmap.md` remains sole mutable current-stage/status/allowed-work/next-action authority and stays compact.

### Index

`docs/index.md` routes directly to current owners such as D0/D1/D2/D4/D4-R1/D5/W2/W3/W4/D6/D7/D8 and current block-specific frontend evidence when needed. It does not route through a stage-wide historical pack.

### Decision/status registry

Do not create a new decision database. Existing ADR registry remains ADR-file status authority. D-stage owners remain semantic authority.

## 5. Active-tree lifecycle

```text
KEEP / CURRENT AUTHORITY
KEEP / CURRENT EVIDENCE
REHOME THEN RETIRE
RETIRE FROM ACTIVE TREE
```

Every removal cites its current replacement owner or states that no target meaning survives.

## 6. Health audit scope

The implementation audits:

- every `D6-R2-NOTIF-01-*`, D6-R2 Fable/finding/ratification/review artifact;
- `D6-R2-AUTHORITY-ROUTE.md`;
- accumulated D6-R2 frontend parent/ledger material;
- P5/P8/P9 block-specific current evidence;
- ADRs 008/010/018/026/030 against satisfied D7 retirement conditions;
- 017/034 against their separate Fact condition; 035 remains;
- completed `docs/plans/**`;
- current locked HTML evidence;
- verification scripts only as needed to distinguish current proof from historical/targeted proof.

Large canonical files such as `components.yaml`, D4 and W2 are not split merely because they are large. Reassess only after accidental history/routing cost is removed.

## 7. CI target

Keep:

```text
required
→ npm run gate
```

Every change retains cheap checks such as required files, one Product OpenAPI entrypoint, unsafe workflow rejection, diff/conflict checks and implementation-block enforcement.

The heavy Product proof runs when changed paths can affect Product contract/proof semantics. Unrelated frontend/docs planning skips it.

The current 95/99/pre-auth/current proof chain remains intact in this health increment because it still carries current protections. If Product-changing CI later becomes a real bottleneck, reopen a bounded proof-decomposition task and require assertion-level parity before removing replay.

## 8. Branch hygiene

Remote branch deletion is a separate post-health-merge operation. Each branch requires PR/merge/unique-commit absorption proof. No age/name heuristics and no history rewrite.

## 9. Current frontend sequencing

During health:

- PR #69 stays **PAUSED / NO P8**;
- PR #70 stays paused at its implementation-plan gate;
- Product/OAD semantics stay unchanged;
- B10/B20 finding is preserved.

After health integration:

```text
reanchor cleaned main
→ revalidate #70 finding against current owners
→ resume accepted #70 gate if still valid
→ bounded B10 correspondence revalidation
→ resume B20
→ continue D6-R2 block-by-block
```

## 10. Rejected approaches

Reject another methodology/profile/router layer, `docs/archive/`, filename-based mass deletion, rewriting history files to fake current status, splitting all large canonical files, or disabling/weaking Product proof for speed.

## 11. Proof strategy

Health proves:

1. fresh bootstrap reaches current owners without historical route pack;
2. representative questions resolve through a small owner set;
3. every retired artifact has no surviving target meaning or an exact replacement owner;
4. no current links point to retired files;
5. index does not route stale status/count/method versions as current truth;
6. ADR retirement matches accepted stage state;
7. completed plans do not remain active without a live consumer;
8. unrelated docs/frontend CI skips heavy Product proof while cheap checks remain;
9. Product/proof-input changes retain the exact heavy proof currently green;
10. one required check remains;
11. Product semantics/counts remain unchanged;
12. locked frontend evidence remains unchanged;
13. Git history remains recoverable history.

No prose-status lint is added.

## 12. Success criteria

A fresh session should normally need only:

```text
AGENTS
→ roadmap
→ applicable method
→ index
→ task owner(s)
```

Absorbed intermediate D6-R2/review/ratification/plan material is retired only after rehome proof; satisfied ADR residues are retired; Product OAD is unchanged; unrelated frontend/docs work no longer pays for heavy Product proof; Product/proof changes retain full protection; #69/#70 remain resumable.

No arbitrary file-count target is frozen.

## 13. Reopen triggers

Reopen if a supposed historical file is the only current owner of an invariant; Git history cannot satisfy a named live evidence need; the diff predicate can skip a materially consumed proof input; Product-changing CI becomes a bottleneck with a provable safe decomposition opportunity; direct routing makes a decision undiscoverable; a ledger contains unique operator-locked meaning; or a large canonical owner still causes repeated context overflow after cleanup and has a proven responsibility split.

## 14. Non-goals

No Product redesign, no change to #70's accepted finding, no Product-count change, no runtime/frontend/backend implementation, no D0–D8 restart, no external methodology dependency, no archive/knowledge platform, no unaudited branch deletion, no automatic merge.

## 15. Current gate

The written design is operator-approved. **Implementation-plan review is now the only gate.** No cleanup or CI execution begins until that plan is explicitly approved.
