# Repository Governance & Context Health Rebaseline — Design

> **Status:** DESIGN CANDIDATE — direction operator-approved on 2026-08-25; written-spec review required before repository changes
> **Trigger:** D6-R2 frontend work became slow/context-heavy after temporary methodology/router experiments and accumulated rebaseline evidence remained in the active tree
> **Method:** DevelopmentConexus Engineering Method v1.0.0
> **Baseline:** `main@bdbbef43ed3a5e9d912e67ddac5173024352eaa3`
> **Product semantics:** OUT OF SCOPE — preserve current Product/OAD meaning, 106 operations, 31 ordinary Permissions, H/A/S, runtime NONE

## 1. Decision

**Outcome: `RESTRUCTURE NOW` — restore the Marketplace Central repository to a small active-authority graph, retire absorbed intermediate material to Git history, and make CI proportional to the changed claim.**

This is a repository-governance/context-health correction, not another methodology framework.

Target operating path:

```text
AGENTS.md
  → docs/roadmap.md                 current state / allowed work / next action
  → applicable local method         how to reason
  → docs/index.md when needed       where the current owner lives
  → smallest current owner(s)       Product/architecture/contract/frontend authority
  → evidence/history only when the concrete question requires it
```

The active tree answers **what is true now and where that truth lives**. Git history/evidence answers **how the repository arrived there**.

## 2. Root cause

The repository is effective at producing findings, reviews, ratifications, repairs and proofs, but it has not consistently performed the final lifecycle step after a decision is absorbed:

```text
investigate
→ decide
→ ratify/prove
→ rehome into canonical owner
→ RETIRE intermediate material from active tree
```

Instead, many intermediate artifacts remained routed/reachable as durable current material. This creates three related failure classes.

### 2.1 Context authority inflation

D6-R2 currently contains a large chain of `NOTIF-01`, finding, ratification, repair, feed-forward and review artifacts. `D6-R2-AUTHORITY-ROUTE.md` intentionally keeps the entire chain reachable even though much of its surviving meaning has already been rehomed.

Consequences:

- agents see multiple historical snapshots while trying to answer a current question;
- stale counts/status/method versions must be reconciled repeatedly;
- a task-specific read expands into a historical reconstruction rather than a current-owner lookup;
- context-window cost grows without improving current decision quality.

Concrete examples already observed:

- `D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md` still contains older 99-operation / 30-Permission and early-lock-state prose;
- `D6-R2-P8-BLOCK-LEDGER.md` still contains older B12/B10 status and Frontend Method v2.1 prose;
- `D6-R2-AUTHORITY-ROUTE.md` routes the complete historical chain despite `docs/roadmap.md` being the current-status authority.

### 2.2 Retirement drift

The ADR registry explicitly kept some legacy ADRs only until D7 adjudication. D7 is now accepted/closed, but the registry still lists those entries as reopened D7 residues.

This proves the repository has a missing or inconsistently executed retirement step rather than merely “too many files.”

### 2.3 Verification inflation

The repository exposes one required `npm run gate`, but the Product OAD verifier recursively reconstructs and executes historical proof generations (current 106 → historical 99 → pre-auth → baseline 95) even for documentation-only changes.

Recent docs-only PR evidence shows the aggregate gate spending roughly 85 seconds in the Product proof while no Product OAD file changed.

This contradicts the repository rule that historical proofs must not become permanent prerequisites for unrelated planning work.

## 3. Governing invariants

### G1 — Repository authority is current, not archival

The active tree retains a document only while at least one of these is true:

1. it owns current Product/architecture/contract/repository semantics;
2. it is current mutable status/routing/method authority;
3. it is evidence still required for a currently open claim, reopen trigger or unresolved proof obligation;
4. an active current owner cites it because its exact content remains necessary to interpret current meaning.

If none applies, Git history is sufficient archive.

### G2 — Rehome before retire

No intermediate artifact is removed while it contains still-required meaning that has no current owner.

Retirement procedure:

```text
identify surviving meaning
→ locate/create legitimate current owner
→ rehome surviving meaning once
→ update current routes/references
→ verify no active obligation depends on the intermediate artifact
→ remove the intermediate file from the active tree
```

Deletion never means “forget the decision”; Git retains history.

### G3 — One meaning, one current owner

A finding, review, ratification, plan or proof does not remain a parallel semantic authority after acceptance. Its durable conclusion belongs in the owner that actually governs that meaning.

### G4 — Index routes current questions, not history

`docs/index.md` is a selective current-authority map. It must not become a comprehensive catalog of every artifact ever produced.

Its job is:

> Given a concrete question, where is the smallest current authority likely to answer it?

### G5 — Methods reason; repository owns Product state

Keep the repository-local:

- `docs/development/engineering-method.md` — Global Maximum / engineering reasoning;
- `docs/development/frontend-product-experience-planning-method.md` — frontend planning;
- `docs/development/engineering-rules.md` — repo-specific Git/CI/safety specialization.

Do not restore the external `conexus-methodology` router/pin/profile experiment. Do not create another methodology layer.

### G6 — Inquiry may expand; default context stays selective

There is no hard file/context budget. A material finding may require whole-repository/history/external-source analysis.

But the default operating rule is:

```text
start from current route
→ read smallest likely owner set
→ expand only when evidence can materially change/falsify the conclusion
```

“Context may expand” does not mean “historical material must remain in the bootstrap graph.”

### G7 — CI proves the changed claim proportionately

There remains exactly one required aggregate check.

Every change runs cheap universal repository checks. Heavy Product OAD/generator/history proofs run only when the diff touches the surface whose claim they protect or when explicitly invoked as targeted evidence.

Historical proof remains available when material; it is not replayed by default on unrelated docs/frontend-planning changes.

## 4. Target repository governance model

### 4.1 Bootstrap

`AGENTS.md` remains intentionally small. It should do only the following:

1. revalidate repository / branch / HEAD / main / relevant PR;
2. read `docs/roadmap.md`;
3. route to Engineering Method for material engineering;
4. route to Frontend Method for frontend Product Experience planning;
5. route to `docs/index.md` to locate current owners when needed;
6. state that current repository authority outranks chat/history;
7. state hard Marketplace safety rails and Git/merge rules.

No task-pack, methodology profile, external pin, historical-artifact router or default whole-repository read list belongs in AGENTS.

### 4.2 Roadmap

`docs/roadmap.md` remains the sole mutable current-stage authority and should stay compact.

It answers only:

- current stage/checkpoint;
- currently accepted baseline;
- active increment/prerequisite;
- allowed/blocked work;
- active upstream finding if one exists;
- exact next action;
- Product surface/runtime gate when materially useful.

It does not reconstruct full project chronology.

### 4.3 Documentation index

`docs/index.md` becomes the primary selective current-authority map.

Target shape:

| Question / task | Current owner(s) |
| --- | --- |
| Product purpose / scope / actors | D0 |
| domain ownership / semantic edges | D1 |
| identity / Organization / provenance | D2 |
| Mercado Livre / Sankhya external contract | D4 |
| publication / Readiness / ListingIntent seam | D4-R1 |
| Product API semantics | D5 |
| schema/read-write grammar | W2 |
| collection/search grammar | W3 |
| access/Permission/client class | W4 |
| frontend architecture / interaction planning | D6 + Frontend Method |
| current D6-R2 screen/surface program | current P5/P8/P9 owner only when active |
| runtime | D7 |
| golden-flow proof | D8 |
| stable cross-stage invariants | `ARCHITECTURE.md` |
| evidence for a named decision | Evidence Register / Git history only when needed |

Do not route through a stage-wide “everything remains reachable” authority pack.

### 4.4 Decision/status registry

Do not create a new giant decision database.

Where a compact registry materially helps navigation, it contains only:

```text
decision/topic
current status
current owner
reopen trigger/reference when material
```

It does not copy the decision prose itself. Existing ADR registry remains the status authority for retained ADR files; D-stage owners remain semantic authority.

## 5. Active-tree artifact lifecycle

Every durable artifact belongs to one of four classifications.

### A — KEEP / CURRENT AUTHORITY

Current owner or current mutable router/method.

Examples:

- AGENTS / roadmap / index / methods / engineering rules;
- ARCHITECTURE;
- accepted D-stage owners whose semantics are still current;
- canonical Product OAD and current bounded verification logic.

### B — KEEP / CURRENT EVIDENCE

Not authority, but still needed for an unresolved/current claim or a named live proof obligation.

Examples may include:

- currently operator-locked low-fi wireframes still used by the active frontend stage;
- live provider evidence still required by an unresolved capability claim.

### C — REHOME THEN RETIRE

Contains still-valid meaning, but that meaning belongs in a different current owner. Rehome once, then remove the intermediate artifact.

Likely examples discovered during health audit:

- accepted findings/ratifications whose surviving rule is not yet present in the canonical D-stage owner;
- stage closure/ledger material that still contains unique current frontend decisions.

### D — RETIRE FROM ACTIVE TREE

Fully absorbed, superseded or historical material with no current unresolved proof obligation.

Likely candidates include:

- old ratification wrappers after their exact decision is canonical elsewhere;
- review/adjudication artifacts whose accepted fixes have landed;
- old implementation plans after execution and rehome;
- superseded route packs;
- legacy ADR residues whose own retirement condition is already satisfied;
- historical verifier layers that are needed only as Git evidence rather than as a current mandatory gate.

Every removal must cite its current replacement owner or state explicitly that no target meaning survives.

## 6. First health audit scope

The health implementation must inventory, classify and adjudicate at least these groups.

### 6.1 D6-R2 / NOTIF-01 chain

Audit every `D6-R2-NOTIF-01-*`, Fable adjudication/ratification, finding, repair, feed-forward and proof document.

For each:

```text
KEEP AUTHORITY
KEEP CURRENT EVIDENCE
REHOME THEN RETIRE
RETIRE
```

The target is not “delete NOTIF-01.” The target is to keep only current meaning/evidence that still has a real consumer.

`D6-R2-AUTHORITY-ROUTE.md` is presumed unnecessary unless audit proves a current navigation job that `docs/index.md` cannot satisfy directly.

### 6.2 D6-R2 accumulated frontend documents

Audit at least:

- `D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md`;
- `D6-R2-P8-BLOCK-LEDGER.md`;
- P5/P8/P9 block-specific artifacts;
- currently locked HTML evidence.

Rules:

- do not rewrite historical facts merely to make old files look current;
- if a file is still a current owner, compact/update only its current-owned portion;
- if it is an accumulated execution ledger rather than current authority, rehome the surviving block truth and retire it;
- locked HTML remains current evidence only while the active frontend stage depends on the exact accepted structure.

### 6.3 ADR registry

Re-evaluate retained legacy ADR retirement conditions against the now-accepted D7/D8 state.

Priority candidates whose registry currently says “retire when D7 adjudicates”:

```text
008
010
018
026
030
```

`017/034` are adjudicated against their own Fact rehoming condition.

`035` remains while D0–D9 transition authority is still required.

Update the registry after each retirement; do not keep a stale status table.

### 6.4 Plans/specs

Audit both existing locations:

```text
docs/plans/**
docs/superpowers/plans/**
docs/superpowers/specs/**
```

Rule:

> A plan/spec is working material. After its accepted meaning is rehomed and its execution/proof is complete, remove it from the active tree unless a concrete current consumer still requires the exact document.

The current human-operable read-projection spec/plan remain current working material while PR #70 is paused. After that prerequisite eventually lands and meaning is rehomed, they are expected retirement candidates.

This repository-health spec itself is also temporary working material and must retire after its rules are rehomed into `AGENTS.md` / `docs/index.md` / `docs/development/engineering-rules.md` and the health increment closes.

### 6.5 Architecture decisions/citations

Keep only retained ADRs and citation files with a current retained ADR consumer. Do not create `docs/archive/`.

### 6.6 Product contract files

Large canonical files such as `components.yaml`, D4 and W2 are **not** split merely because they are large.

First remove accidental historical/context complexity. Reassess canonical file size afterward. Split only if one file demonstrably owns multiple independently navigable responsibilities and the split reduces total decision/context cost without creating parallel authority.

## 7. CI / verification target

### 7.1 One required check remains

Keep:

```text
GitHub required check: required
→ npm run gate
```

No new CI jobs/check names are introduced.

### 7.2 Universal cheap checks

Every PR/push keeps cheap objective checks such as:

- required bootstrap files;
- exactly one Product OpenAPI entrypoint;
- unsafe workflow trigger rejection;
- diff/conflict-marker check;
- implementation-block changed-path enforcement.

### 7.3 Diff-aware Product OAD proof

Heavy current Product OAD proof runs when changed files include the canonical Product contract or proof/authority inputs whose meaning can alter the current Product wire.

The implementation plan must define an explicit allowlist/predicate from current repository evidence rather than a broad “all docs” trigger.

At minimum, candidates include:

```text
contracts/api/product/**
scripts/verify-product-oad*.mjs
scripts/lib/*oad*.mjs
D5/W1/W2/W3/W4 authority files when the verifier derives exact expectations from them
D4/D4-R1 only where a current mechanical Product proof explicitly reads them
```

A docs-only frontend planning change that cannot change Product OAD semantics must not pay for full Product generation.

### 7.4 Current proof vs historical replay

Separate two claims:

```text
CURRENT CONTRACT PROOF
→ does the current canonical 106/31/H-A-S OAD satisfy its current contract/tooling invariants?

HISTORICAL MIGRATION/NON-REGRESSION REPLAY
→ can an older 95/99/pre-auth transition be reconstructed and re-proved?
```

The required gate protects the first when relevant.

The second becomes a targeted proof invoked only when changing the historical replay machinery or a claim that specifically depends on that transition.

Do not recursively regenerate 95 → 99 → 106 on every Product proof unless the implementation audit demonstrates a current correctness property that cannot be expressed directly against current authority.

### 7.5 Generation determinism

Current OAD changes still receive proportional generation proof (Redocly / TypeScript / Go) sufficient for the current contract claim. The health change must not weaken current schema/access/source-reachability/route generation properties merely to make CI faster.

The optimization target is removing **obsolete repeated proof**, not current protection.

## 8. Git branch hygiene

Branch deletion is a separate post-integration operational step, not mixed into the content cleanup PR.

After the health PR lands:

1. list every remote branch;
2. identify associated PR/merge status;
3. prove whether its commits are absorbed/superseded or still unique;
4. preserve current active branches;
5. delete only branches objectively safe to retire.

Never infer safety from branch age/name alone. No force push/history rewrite.

## 9. Interaction with current frontend work

During the health rebaseline:

- PR #69 / B20 remains **PAUSED / NO P8**;
- PR #70 / human-operable read-projection prerequisite remains **PAUSED at implementation-plan gate**;
- do not modify Product/OAD semantics as part of health cleanup;
- do not render B20;
- do not silently discard the B10/B20 finding.

After health integration:

```text
reanchor from cleaned main
→ revalidate PR #70 design/finding against current owners
→ if still valid, resume its approved gate sequence
→ bounded B10 correspondence revalidation
→ resume B20
→ continue D6-R2 block-by-block
```

The expected benefit is that each frontend block can locate its current Product/architecture authority directly and reopen only the smallest owner when frontend evidence falsifies an older planning assumption.

## 10. Rejected approaches

### 10.1 Add another methodology/profile/router layer — REJECTED

This recreates the failure class. Methods remain local reasoning standards; repository authority remains local/current.

### 10.2 Move old material to `docs/archive/` — REJECTED

An archive folder preserves the same duplicate active-tree/history cost. Git is the archive.

### 10.3 Delete everything called finding/ratification — REJECTED

Filename is not evidence of redundancy. Surviving unique meaning/evidence must be rehomed or retained until its consumer closes.

### 10.4 Rewrite old historical documents to current status — REJECTED

This destroys truthful historical snapshots. Prefer rehome + retirement when a historical execution artifact no longer owns current state.

### 10.5 Split every large canonical document — REJECTED

Large current owner files may be legitimate. Split only after accidental history/routing cost is removed and a responsibility boundary is proven.

### 10.6 Disable Product verification to speed CI — REJECTED

Current Product proof remains mandatory for Product-affecting changes. Only unrelated/historical replay cost is removed from the default path.

## 11. Proof strategy

The health implementation plan must prove at least:

1. **Bootstrap:** a fresh session can identify current stage and next action from AGENTS + roadmap without a historical route pack.
2. **Selective index:** representative questions (Product identity, publication, B20 Listing read, runtime) each resolve to a small current-owner set.
3. **No current-authority loss:** every retired artifact either has zero surviving target meaning or a cited current replacement owner containing that meaning.
4. **No broken current refs:** repository links/references to retired active files are removed/repointed where current consumers exist.
5. **No stale current route:** index/router does not direct a fresh agent to a known superseded status/method/count as current truth.
6. **ADR retirement consistency:** ADR registry retirement conditions match accepted D-stage state.
7. **Plans lifecycle:** completed plans/specs do not remain active without a named current consumer.
8. **CI docs-only proportionality:** a documentation-only change outside Product mechanical authority executes cheap gate checks without full OAD generation/replay.
9. **CI Product proportionality:** a Product OAD change still executes current lint/bundle/generation/semantic proof and fails on a controlled Product-contract regression.
10. **Historical proof availability:** historical replay, if retained, is invokable as a targeted proof but is not a default prerequisite for unrelated work.
11. **One required check:** workflow still exposes one required aggregate gate.
12. **Product invariance:** health PR changes no Product operation/Permission/Principal count or Product semantics.
13. **Frontend preservation:** locked current HTML evidence needed by D6-R2 remains untouched unless separately adjudicated.
14. **Git history preservation:** no force update/history rewrite; retired files remain recoverable through Git history.

Do not add a CI test that lints prose status/version strings across historical docs. Staleness/authority quality is an engineering-review concern; mechanical proof should protect routing/reachability/current machine contracts only where objectively expressible.

## 12. Success criteria

The health increment is successful when:

- the normal fresh-session route is understandable from `AGENTS.md`, `docs/roadmap.md`, the applicable method and `docs/index.md`;
- `docs/index.md` routes directly to current owners rather than a comprehensive D6-R2 history pack;
- absorbed intermediate D6-R2/NOTIF/review/ratification/plan material is removed from active tree after surviving meaning is rehomed;
- legacy ADRs whose explicit retirement conditions are satisfied are removed and the registry is current;
- Product OAD canonical semantics remain unchanged;
- docs-only planning CI no longer executes recursive Product OAD historical replay;
- Product-affecting changes retain strong current-contract proof;
- Git history remains the historical archive;
- PR #69/#70 remain safely resumable from the cleaned governance baseline.

A numerical file-count target is deliberately not frozen. The goal is **minimal sufficient active authority/evidence**, not arbitrary deletion volume.

## 13. Reopen triggers

Reopen this design if implementation evidence shows:

- a supposedly historical artifact is the only surviving owner of a current invariant;
- Git history alone cannot satisfy a named current evidence/proof obligation;
- current Product contract correctness genuinely depends on replaying an older 95/99 transition on every Product change;
- direct owner routing makes a material current decision undiscoverable without a bounded supplemental registry;
- removing a route/ledger would lose operator-locked frontend meaning that has no current replacement owner;
- a canonical large file, after cleanup, still creates repeated task-level context overflow and has a provable responsibility split.

## 14. Explicit non-goals

This health rebaseline does **not**:

- redesign Marketplace Central Product semantics;
- change the accepted Global Maximum/read-projection finding in PR #70;
- change Product OAD operation/Permission/Principal counts;
- implement backend/frontend/runtime code;
- restart D0–D8;
- create a new methodology repository dependency;
- create an archive folder;
- introduce a document database, knowledge graph, generated documentation portal or semantic search system;
- delete remote branches before branch-by-branch absorption audit;
- merge PR #69, #70 or the health PR without explicit operator authorization.

## 15. Written-spec gate

Before implementation planning:

- operator reviews this written design;
- retirement/rehoming rules are accepted;
- CI proportionality boundary is accepted;
- no required current authority/evidence class is missing;
- scope remains repository health/governance only.

Only after written-spec approval should an implementation plan enumerate exact files to keep/rehome/retire and exact gate predicates.