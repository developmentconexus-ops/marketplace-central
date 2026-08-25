# Repository Governance & Context Health Rebaseline — Design

> **Status:** DESIGN APPROVED — operator-approved on 2026-08-25; implementation plan review required before repository changes
> **Trigger:** D6-R2 frontend work became slow/context-heavy after temporary methodology/router experiments and accumulated rebaseline evidence remained in the active tree
> **Method:** DevelopmentConexus Engineering Method v1.0.0
> **Baseline:** `main@bdbbef43ed3a5e9d912e67ddac5173024352eaa3`
> **Product semantics:** OUT OF SCOPE — preserve current Product/OAD meaning, 106 operations, 31 ordinary Permissions, H/A/S, runtime NONE

## 1. Decision

**Outcome: `RESTRUCTURE NOW` — restore Marketplace Central to a small active-authority graph, retire absorbed intermediate material to Git history, and make CI proportional to the changed claim.**

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

This created three failure classes: context-authority inflation, retirement drift, and verification inflation.

D6-R2 has a large `NOTIF-01`/finding/ratification/review chain; accumulated parent/ledger files still contain stale counts/status/method versions; and the single `npm run gate` runs heavy Product proof even for unrelated docs-only work.

Implementation planning found one important correction: the current 95/99 replay is not pure history; it still carries valid current Product protections. Therefore this health increment **does not split or weaken that proof**. It makes the heavy proof diff-aware so unrelated frontend/docs work skips it. A future split is a separate reopen only when Product-changing CI itself is a material bottleneck and assertion-level proof parity is demonstrated.

## 3. Governing invariants

### G1 — Active tree is current, not archival

Keep a document active only while it owns current semantics/status/routing/method authority, is current evidence for a live obligation, or is explicitly required by a current owner. Otherwise Git history is sufficient archive.

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

`docs/index.md` answers where the smallest current authority likely lives. It is not a complete document catalog.

### G5 — Methods reason; repository owns state

Preserve local Engineering Method, Frontend Method and repository engineering rules. Do not restore external methodology router/pin/profile machinery.

### G6 — Inquiry may expand; default context stays selective

No hard context/file budget. Start at current owners and expand only when extra evidence can materially change/falsify the conclusion.

### G7 — CI proves the changed claim proportionately

One required aggregate check remains. Cheap universal checks run always. The existing heavy Product proof runs only when Product/proof inputs change; when diff evidence is unavailable it runs fail-safe.

## 4. Target governance model

- `AGENTS.md`: small bootstrap/router.
- `docs/roadmap.md`: sole mutable current status/allowed-work/next-action authority.
- `docs/index.md`: direct current-owner map.
- D-stage/current contract owners: one current semantic home each.
- Evidence/history: read only when a concrete claim needs it.
- No new giant decision database or archive folder.

## 5. Active-tree lifecycle

```text
KEEP / CURRENT AUTHORITY
KEEP / CURRENT EVIDENCE
REHOME THEN RETIRE
RETIRE FROM ACTIVE TREE
```

Every removal cites its current replacement owner or states that no target meaning survives.

## 6. Health audit scope

Audit:

- every `D6-R2-NOTIF-01-*`, D6-R2 Fable/finding/ratification/review artifact;
- `D6-R2-AUTHORITY-ROUTE.md`;
- accumulated D6-R2 frontend parent/ledger material;
- P5/P8/P9 block-specific current evidence;
- ADRs 008/010/018/026/030 against satisfied D7 retirement conditions;
- ADR 017/034 against their separate Fact condition; ADR 035 remains;
- completed `docs/plans/**`;
- current locked HTML evidence;
- verification scripts only as needed to preserve current proof while removing unrelated CI cost.

Large canonical files such as `components.yaml`, D4 and W2 are not split merely because they are large. Reassess only after accidental history/routing cost is removed.

## 7. CI target

Keep exactly:

```text
required
→ npm run gate
```

Every change retains cheap checks such as required files, one Product OpenAPI entrypoint, unsafe workflow rejection, diff/conflict checks and implementation-block enforcement.

Diff-aware routing:

```text
unrelated frontend/docs change
→ cheap checks only

Product/proof-input change
→ existing full Product proof unchanged

no reliable diff base
→ full Product proof fail-safe
```

Do not physically split the current 95/99/pre-auth/current chain in this health increment because it still protects current invariants. Reopen decomposition only with assertion-level parity proof.

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
2. representative questions resolve through a small current-owner set;
3. every retired artifact has no surviving target meaning or an exact replacement owner;
4. no current links point to retired files;
5. index does not route stale status/count/method versions as current truth;
6. ADR retirement matches accepted stage state;
7. completed plans do not remain active without a live consumer;
8. unrelated docs/frontend CI skips heavy Product proof while cheap checks remain;
9. Product/proof-input changes retain the exact heavy Product proof currently green;
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

Absorbed intermediate material is retired only after rehome proof; satisfied ADR residues are retired; Product OAD is unchanged; unrelated frontend/docs work no longer pays for heavy Product proof; Product/proof changes retain full protection; #69/#70 remain resumable.

No arbitrary file-count target is frozen.

## 13. Reopen triggers

Reopen if a supposed historical file is the only current owner of an invariant; Git history cannot satisfy a named live evidence need; the diff predicate can skip a materially consumed proof input; Product-changing CI becomes a bottleneck with a provable safe decomposition opportunity; direct routing makes a decision undiscoverable; a ledger contains unique operator-locked meaning; or a large canonical owner still causes repeated context overflow after cleanup and has a proven responsibility split.

## 14. Non-goals

No Product redesign, no change to #70's accepted finding, no Product-count change, no runtime/frontend/backend implementation, no D0–D8 restart, no external methodology dependency, no archive/knowledge platform, no unaudited branch deletion, no automatic merge.

## 15. Current gate

The written design is operator-approved. **Implementation-plan review is now the only gate.** No cleanup or CI execution begins until that plan is explicitly approved.
