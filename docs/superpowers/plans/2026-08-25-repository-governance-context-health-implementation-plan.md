# Repository Governance & Context Health Rebaseline Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restore Marketplace Central to a small, current-authority repository graph, retire absorbed intermediate planning/history from the active tree, and make the single required CI gate proportional without weakening current Product contract proof.

**Architecture:** Keep the proven repository-local operating path `AGENTS → roadmap → applicable method → index → smallest current owner`; use Git history as archive. Execution is staged so retirement is evidence-driven: first classify every candidate and identify surviving meaning/replacement owner, then rehome, then retire. CI keeps one required `npm run gate`, but current Product proof and historical migration replay become separate claims with diff-aware routing.

**Tech Stack:** Git/GitHub; Markdown repository authority; PowerShell `scripts/gate.ps1`; Node.js `>=26.3.0 <27`; OpenAPI 3.1.2; `@redocly/cli@2.45.0`; `openapi-typescript@7.13.0`; TypeScript `5.9.3`; Go current CI toolchain plus exact `go1.25.1` compatibility proof where the existing Product proof requires it; `oapi-codegen v2.8.0`.

**Spec:** `docs/superpowers/specs/2026-08-25-repository-governance-context-health-design.md`

## Global Constraints

- Preserve current Marketplace Central Product semantics and the current **106 Product operations / 31 ordinary Permissions / Principal kinds H / A / S**.
- Preserve active runtime baseline **NONE** and the D9 implementation block.
- Preserve repository-local `docs/development/engineering-method.md`, `docs/development/frontend-product-experience-planning-method.md`, and `docs/development/engineering-rules.md`; do not restore an external `conexus-methodology` pin/router/profile layer.
- Keep `AGENTS.md` as the small bootstrap/router, `docs/roadmap.md` as sole mutable current-stage/status/allowed-work/next-action authority, and `docs/index.md` as selective current-owner navigation.
- Rehome surviving meaning before retiring its intermediate artifact. Filename/age alone never proves safe deletion.
- Git history is the archive. Do not create `docs/archive/`, copy old files into a second history tree, rewrite history, force-push, reset, clean, or stash by convenience.
- Keep current evidence while a live claim/proof depends on it; especially preserve operator-locked frontend HTML and current B10 evidence needed by the active D6-R2 reopen.
- PR #69 / B20 remains **PAUSED / NO P8** during this health increment.
- PR #70 / human-operable read-projection prerequisite remains paused at its existing implementation-plan gate; this health increment does not alter its Product/OAD design.
- No Product/OAD semantic change belongs in this health PR. `contracts/api/product/**` must remain byte-identical to `main` unless a new material finding stops this plan.
- Keep exactly one required GitHub check: `required → npm run gate`.
- Cheap universal repository checks run on every change. Heavy Product proof runs only when its current contract/proof inputs change. Historical 95/99/pre-auth replay remains invokable evidence but is not a permanent prerequisite for unrelated frontend/docs work.
- Do not weaken current Product proof simply to make CI faster. Remove obsolete/repeated historical work, not current protection.
- Remote branch deletion is out of this PR; it happens only after merge in a separate branch-by-branch absorption audit.
- Merge requires explicit operator authorization.

---

## File Map

### Bootstrap / repository governance

- `AGENTS.md` — small bootstrap/router; add only the selective-start / expand-on-falsifier rule if not already explicit enough.
- `docs/roadmap.md` — current health gate while this increment is active, then integration-ready status.
- `docs/index.md` — direct current-owner map; remove stage-wide historical routing.
- `docs/development/engineering-rules.md` — durable active-tree lifecycle and proportional-CI repository rules.
- `ARCHITECTURE.md` — expected no change; already owns “Git history is history”. Modify only if the audit proves a missing stable cross-stage invariant.

### D6-R2 current frontend authority/evidence

- `docs/engineering/rebaseline/D6-FRONTEND.md`
- `docs/engineering/rebaseline/D6-R2-P5-SCREEN-SURFACE-INVENTORY.md`
- current block-specific P8/P9 documents retained by the Task 1 audit
- current operator-locked HTML in `qualification/d6-r2-wireframes/`
- `docs/engineering/rebaseline/D6-R2-AUTHORITY-ROUTE.md` — expected retirement after direct `docs/index.md` routes are sufficient
- `docs/engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md` — audit/rehome candidate, not blindly deleted
- `docs/engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md` — audit/rehome candidate because mutable status duplicates roadmap and is stale
- every `docs/engineering/rebaseline/D6-R2-NOTIF-01-*` plus D6-R2 Fable/finding/ratification/review artifacts — audit candidates

### ADR cleanup

- `docs/architecture/decisions/README.md`
- expected retirement candidates after replacement-owner verification:
  - `docs/architecture/decisions/008-production-deploy-topology.md`
  - `docs/architecture/decisions/010-mercado-livre-polling-visible-refresh.md`
  - `docs/architecture/decisions/018-mutation-envelope-table-and-poller.md`
  - `docs/architecture/decisions/026-scheduler-phase-vocabulary.md`
  - `docs/architecture/decisions/030-scheduler-second-instance-per-installation.md`
- expected citation retirement with their last retained consumer:
  - `docs/architecture/decisions/_citations/adr-009-citations.md`
  - `docs/architecture/decisions/_citations/adr-013-citations.md`
  - `docs/architecture/decisions/_citations/adr-07-twodigit-citations.md`
  - `docs/architecture/decisions/_citations/adr-08-twodigit-citations.md`
- retain pending their separate condition:
  - ADR 017 / ADR 034 / `_citations/adr-017-citations.md`
  - ADR 035
  - `_citations/RENUMBERING-REGISTRY.md` while retained reconstructed ADRs still consume it

### Planning material

- `docs/plans/2026-08-22-op-read-01-repair.md`
- `docs/plans/2026-08-23-authorization-request-d5-wire.md`
- `docs/plans/2026-08-23-notif-d6-r-b00-r2.md`
- `docs/plans/2026-08-23-notif-d6-r-b11.md`
- current health spec/plan under `docs/superpowers/` — temporary; retire in Task 7 after their durable rules are rehomed

### CI / Product proof

- `.github/workflows/ci.yml` — expected no structural change; one required job remains
- `scripts/gate.ps1` — diff-aware proof routing
- `scripts/verify-product-oad.mjs` — current Product proof only after split
- create `scripts/verify-product-oad-history.mjs` — targeted historical 95/99/pre-auth replay, copied from the current full replay before simplifying the current proof
- retained current targeted Product proofs:
  - `scripts/verify-oad-source-reachability.mjs`
  - `scripts/lib/publication-requirements-oad-proof.mjs`
  - `scripts/verify-operational-read-contract.mjs`
  - `scripts/verify-performance-evidence-knowledge.mjs`
  - `scripts/verify-notification-oad.mjs`
  - `scripts/verify-authorization-request-oad.mjs`
- historical replay dependencies retained only for targeted history proof:
  - `scripts/verify-product-oad-current99.mjs`
  - `scripts/verify-product-oad-pre-auth.mjs`
  - `scripts/verify-product-oad-baseline.mjs`

---

## Task 1: Produce the evidence-backed active-tree retirement map before deleting anything

**Files:**
- Modify: `docs/superpowers/plans/2026-08-25-repository-governance-context-health-implementation-plan.md` (append the execution audit table only; this plan is temporary and is retired in Task 7)
- Read only: current repository docs/scripts listed below

**Interfaces:**
- Consumes: approved artifact lifecycle `KEEP AUTHORITY | KEEP CURRENT EVIDENCE | REHOME THEN RETIRE | RETIRE`.
- Produces: exact path-by-path retirement map used as the only deletion input for Tasks 2–4.

- [ ] **Step 1: Reanchor the execution branch and prove the Product contract baseline is untouched**

Run:

```bash
git status --short
git rev-parse HEAD
git rev-parse main
git merge-base main HEAD
git diff --name-status main...HEAD
git diff --exit-code main...HEAD -- contracts/api/product
```

Expected before implementation: the health branch differs from `main` only by approved health planning artifacts/status updates; `contracts/api/product/**` has no diff.

- [ ] **Step 2: Generate the exact candidate inventory into an untracked Git-internal file**

Run:

```bash
{
  git ls-files 'docs/engineering/rebaseline/D6-R2-NOTIF-01-*'
  git ls-files 'docs/engineering/rebaseline/D6-R2-*FABLE*'
  git ls-files 'docs/engineering/rebaseline/D6-R2-*RATIFICATION*'
  git ls-files 'docs/engineering/rebaseline/D6-R2-*FINDING*'
  printf '%s\n' \
    docs/engineering/rebaseline/D6-R2-AUTHORITY-ROUTE.md \
    docs/engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md \
    docs/engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md \
    docs/engineering/rebaseline/D6-R2-P5-B110-AUTHORIZATION-REQUEST-SUPERSESSION.md \
    docs/engineering/rebaseline/D6-R2-P8-B110-APPROVALS-CANDIDATE.md \
    docs/engineering/rebaseline/D6-R2-P8-B110-APPROVALS-RATIFICATION.md \
    docs/engineering/rebaseline/D6-R2-P9-AUTHORIZATION-REQUEST-BIDIRECTIONAL-SCREEN-CONTRACTS.md \
    docs/engineering/rebaseline/D6-R2-P6-B10-PREPARATION-REFERENCE-STUDY.md \
    docs/engineering/rebaseline/D6-R2-P8-B10-PREPARATION-RATIFICATION.md \
    docs/engineering/rebaseline/D6-R2-P9-B10-PREPARATION-SCREEN-CONTRACT.md
  git ls-files 'docs/plans/**'
  git ls-files 'docs/architecture/decisions/*.md'
  git ls-files 'docs/architecture/decisions/_citations/*.md'
  git ls-files 'scripts/verify-*'
} | sort -u > .git/repository-health-candidates.txt
wc -l .git/repository-health-candidates.txt
```

This file is diagnostic only and MUST NOT be committed.

- [ ] **Step 3: For every candidate, gather current consumers and semantic replacement evidence**

For each path in `.git/repository-health-candidates.txt`, run the equivalent of:

```bash
path='<candidate-path>'
base=$(basename "$path")
printf '\n=== %s ===\n' "$path"
rg -n --fixed-strings "$base" AGENTS.md ARCHITECTURE.md docs contracts scripts qualification .github \
  --glob '!docs/superpowers/plans/2026-08-25-repository-governance-context-health-implementation-plan.md' || true
sed -n '1,80p' "$path"
```

Then read enough of the file and its claimed replacement owner(s) to answer all four questions:

```text
1. Does this file still own current meaning?
2. Is it current evidence for a live unresolved/proof obligation?
3. If it contains surviving meaning, where is that exact meaning owned now?
4. Would removing it make a fresh agent lose a current decision rather than only its history?
```

- [ ] **Step 4: Append the exact execution audit table to this plan**

Append a section named `## Execution Audit — Active Tree` with one row per candidate using exactly these columns:

```markdown
| Path | Class | Surviving meaning | Current replacement owner / current consumer | Action |
| --- | --- | --- | --- | --- |
| `...` | KEEP AUTHORITY / KEEP CURRENT EVIDENCE / REHOME THEN RETIRE / RETIRE | concise exact statement or `none` | exact path/section or `none` | keep / rehome+rm / rm |
```

Rules:

- `RETIRE` requires `Surviving meaning = none` OR an exact current replacement owner already containing it.
- `REHOME THEN RETIRE` requires an exact existing canonical owner; if no legitimate owner exists, STOP as a material finding rather than inventing one inside cleanup.
- operator-locked HTML needed by D6-R2 is `KEEP CURRENT EVIDENCE` unless a separate frontend adjudication says otherwise.
- no row may use vague text such as “probably”, “old”, “looks superseded”, or “misc”.

- [ ] **Step 5: Validate the audit before any delete**

Run:

```bash
rg -n 'probably|maybe|TBD|TODO|unknown replacement|archive/' docs/superpowers/plans/2026-08-25-repository-governance-context-health-implementation-plan.md
```

Expected: no audit row uses an unresolved deletion rationale and no archive plan exists.

- [ ] **Step 6: Commit the audit checkpoint**

```bash
git add docs/superpowers/plans/2026-08-25-repository-governance-context-health-implementation-plan.md
git commit -m "docs: classify active repository authority and history"
```

This commit is the reviewable safety boundary before any retirement.

---

## Task 2: Restore the small bootstrap graph and retire the stage-wide historical route pack

**Files:**
- Modify: `AGENTS.md`
- Modify: `docs/index.md`
- Modify: `docs/development/engineering-rules.md`
- Modify: `docs/roadmap.md`
- Delete only after direct routes are proven: `docs/engineering/rebaseline/D6-R2-AUTHORITY-ROUTE.md`

**Interfaces:**
- Consumes: Task 1 audit rows for bootstrap/router files.
- Produces: `AGENTS → roadmap → method → index → current owner` with no comprehensive D6-R2 history pack in the normal route.

- [ ] **Step 1: Make selective context expansion explicit in AGENTS without adding a new methodology layer**

Keep the existing bootstrap order and add only this bounded operating rule after the start list:

```markdown
Start with the smallest current owner set likely to answer the task. Expand into additional owners, Evidence, Git history, runtime or external sources only when they can materially change or falsify the conclusion. “No fixed context budget” permits expansion; it does not make historical artifact packs part of the default read path.
```

Do not add profiles, task packs, external methodology refs, generated routing or file-count limits.

- [ ] **Step 2: Rewrite `docs/index.md` as a direct current-owner map**

Preserve the method and core-owner routes, but remove the route through `D6-R2-AUTHORITY-ROUTE.md`. The frontend/current program portion must route directly, proportionately:

```text
frontend architecture / client-state / topology
→ D6-FRONTEND.md + Frontend Method

current complete screen/surface inventory
→ D6-R2-P5-SCREEN-SURFACE-INVENTORY.md

current B10 preparation acceptance evidence
→ D6-R2-P8-B10-PREPARATION-RATIFICATION.md
→ D6-R2-P9-B10-PREPARATION-SCREEN-CONTRACT.md
→ qualification/d6-r2-wireframes/b10-preparation.html

current approvals acceptance evidence when needed
→ retained B110 ratification/P9 artifact from Task 1 audit

current mutable block/program status
→ roadmap only
```

Do not list every NOTIF/finding/ratification file.

- [ ] **Step 3: Add the durable active-tree lifecycle to `engineering-rules.md`**

Add a compact repository-local section with this exact policy:

```text
current owner/current router/current method                 KEEP
current evidence with a named live consumer                KEEP while live
accepted intermediate with unique surviving meaning        REHOME once, then retire
fully absorbed/superseded intermediate                     RETIRE
history                                                     Git
```

Also state:

- a plan/spec/review/ratification is not a second semantic authority after acceptance;
- no `docs/archive/` baseline;
- rehome before retirement;
- task routing points to current owners, not historical chains.

- [ ] **Step 4: Update roadmap to make repository health the current prerequisite during execution**

Keep Product/frontend status intact but set the current sequencing to:

```text
Repository Governance & Context Health Rebaseline — EXECUTION ACTIVE
PR #69 B20 — PAUSED / NO P8
PR #70 human-operable read projection — PAUSED at implementation-plan gate
Product/OAD semantics — unchanged
exact next action — complete health rehome/retirement/CI proof, integrate only after operator authorization, then reanchor #70/#69
```

Do not move D6-R2, D9 or runtime implementation forward.

- [ ] **Step 5: Prove the route pack has no remaining current consumer, then delete it**

Run before deletion:

```bash
rg -n --fixed-strings 'D6-R2-AUTHORITY-ROUTE.md' AGENTS.md ARCHITECTURE.md docs contracts scripts qualification .github || true
```

Expected after the index edit: no current route requires it. Then:

```bash
git rm docs/engineering/rebaseline/D6-R2-AUTHORITY-ROUTE.md
```

- [ ] **Step 6: Verify bootstrap/navigation diff**

```bash
git diff --check
git diff --name-status HEAD~1..HEAD 2>/dev/null || true
rg -n 'conexus-methodology|methodology profile|authority pack' AGENTS.md docs/index.md docs/development/engineering-rules.md
```

Expected: no external methodology router/profile is introduced; direct current owner routing is visible.

- [ ] **Step 7: Commit**

```bash
git add AGENTS.md docs/index.md docs/development/engineering-rules.md docs/roadmap.md
git add -u docs/engineering/rebaseline/D6-R2-AUTHORITY-ROUTE.md
git commit -m "docs: restore selective repository authority routing"
```

---

## Task 3: Rehome surviving D6-R2 meaning and retire absorbed frontend/NOTIF/review intermediates

**Files:**
- Modify only when Task 1 marks `REHOME THEN RETIRE`: existing canonical owners among
  - `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
  - `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
  - `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
  - `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`
  - `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
  - `docs/engineering/rebaseline/D5-API.md`
  - `docs/engineering/rebaseline/D5-B2-W1/W2/W3/W4` exact existing owner files as named by the audit
  - `docs/engineering/rebaseline/D6-FRONTEND.md`
  - `docs/engineering/rebaseline/D6-R2-P5-SCREEN-SURFACE-INVENTORY.md`
  - current block-specific P8/P9 owner selected by the audit
  - `docs/engineering/rebaseline/D7-RUNTIME-JOBS-TRANSACTIONS.md` / D7-B/C/D/E owner when runtime meaning survives
  - `docs/engineering/rebaseline/D8-GOLDEN-FLOWS.md` when golden-flow meaning survives
- Delete: only paths whose Task 1 audit action is `rm` or `rehome+rm`
- Preserve: all Task 1 `KEEP AUTHORITY` / `KEEP CURRENT EVIDENCE` paths and locked HTML

**Interfaces:**
- Consumes: exact Task 1 execution audit.
- Produces: no accepted current meaning that depends on a historical D6-R2 intermediate file.

- [ ] **Step 1: Rehome every `REHOME THEN RETIRE` row into its named existing owner**

For each audit row, copy only the surviving normative conclusion—not investigation chronology, reviewer dialogue, ratification ceremony or obsolete status—into the exact replacement owner.

Required form in the owner is proportional to the meaning, for example:

```markdown
- AuthorizationRequest presentation snapshots are purpose-specific historical presentation and never current source identity.
```

Not:

```markdown
- In D6-R2-NOTIF-01-D2-R4 we first considered A, then Fable said B, then ratification...
```

If the surviving meaning would require a new Product owner, new operation/Permission, or a material architecture change, STOP this health task and surface the finding; do not smuggle product redesign into cleanup.

- [ ] **Step 2: Resolve accumulated D6-R2 parent/ledger status duplication**

Use the audit outcome for:

```text
D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md
D6-R2-P8-BLOCK-LEDGER.md
```

Preferred outcome when their current unique meaning is already carried by D6/P5/block-specific P8/P9/roadmap:

```bash
git rm docs/engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md
git rm docs/engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md
```

If the audit found unique current meaning, rehome exactly that meaning first, then remove. Do not “refresh” historical 99/30/v2.1 snapshots merely to keep the ledgers alive.

- [ ] **Step 3: Retire the absorbed D6-R2/NOTIF/review rows marked safe**

Generate the delete list from the committed execution audit, review it, then delete only those exact paths. A safe shell pattern is:

```bash
# Build this file manually from audit rows whose Action is rm/rehome+rm; do not infer from filename.
: > .git/repository-health-retire.txt
# one exact repository path per line
cat .git/repository-health-retire.txt
while IFS= read -r path; do
  test -n "$path" || continue
  git rm -- "$path"
done < .git/repository-health-retire.txt
```

The delete-list file remains untracked under `.git/`.

- [ ] **Step 4: Prove every deleted filename has no current incoming reference**

For each path in `.git/repository-health-retire.txt`:

```bash
while IFS= read -r path; do
  base=$(basename "$path")
  hits=$(rg -n --fixed-strings "$base" AGENTS.md ARCHITECTURE.md docs contracts scripts qualification .github \
    --glob '!docs/superpowers/plans/2026-08-25-repository-governance-context-health-implementation-plan.md' || true)
  if [ -n "$hits" ]; then
    printf 'STALE REF %s\n%s\n' "$path" "$hits"
    exit 1
  fi
done < .git/repository-health-retire.txt
```

Expected: zero stale current references.

- [ ] **Step 5: Prove current frontend evidence is preserved**

Run:

```bash
git diff --exit-code main -- qualification/d6-r2-wireframes/b00-app-shell.html \
  qualification/d6-r2-wireframes/b00-r2-notifications.html \
  qualification/d6-r2-wireframes/b01-overview.html \
  qualification/d6-r2-wireframes/b10-preparation.html \
  qualification/d6-r2-wireframes/b11-notifications-inbox.html \
  qualification/d6-r2-wireframes/b110-approvals.html \
  qualification/d6-r2-wireframes/b12-notification-routing-settings.html
```

Expected: no locked/current HTML changed by repository cleanup.

- [ ] **Step 6: Commit D6-R2 canonicalization/retirement**

```bash
git add docs/engineering/rebaseline docs/index.md docs/roadmap.md
git add -u docs/engineering/rebaseline
git commit -m "docs: retire absorbed D6-R2 intermediate authority"
```

---

## Task 4: Retire satisfied legacy ADR residue and completed planning material

**Files:**
- Modify: `docs/architecture/decisions/README.md`
- Delete after replacement verification: ADR 008, 010, 018, 026, 030
- Delete after last-consumer verification: four citation files listed in File Map
- Delete: completed `docs/plans/*.md` rows classified RETIRE by Task 1
- Preserve: ADR 017/034/035 and required citation provenance

**Interfaces:**
- Consumes: accepted D7/D8 state plus Task 1 audit.
- Produces: ADR registry whose active residues have real unsatisfied retirement conditions only.

- [ ] **Step 1: Verify replacement owners for each D7-retirement ADR**

Read the retained ADR and corresponding accepted current owner:

```text
ADR 008 deployment topology residue   → D7-E / current D7 deployment-operability authority
ADR 010 polling/acquisition residue   → D7 runtime + D4 external acquisition authority
ADR 018 durable external effect       → D7-C + D3/D4 safety semantics
ADR 026 scheduler/phase runtime       → D7 runtime/jobs authority
ADR 030 scheduler instance topology   → D7 runtime/jobs authority
```

For each, confirm its registry retirement condition says D7 adjudication and that accepted D7 owns the surviving target meaning. If any ADR still contains unique current target meaning absent from the named owner, rehome that exact meaning before removal; if it materially changes current architecture, STOP.

- [ ] **Step 2: Delete the five satisfied ADR residues and update the registry**

```bash
git rm \
  docs/architecture/decisions/008-production-deploy-topology.md \
  docs/architecture/decisions/010-mercado-livre-polling-visible-refresh.md \
  docs/architecture/decisions/018-mutation-envelope-table-and-poller.md \
  docs/architecture/decisions/026-scheduler-phase-vocabulary.md \
  docs/architecture/decisions/030-scheduler-second-instance-per-installation.md
```

Update `docs/architecture/decisions/README.md` so these are no longer described as active residues; record that their retirement condition was satisfied by accepted D7 and that Git history retains them.

- [ ] **Step 3: Remove citation files whose last retained consumer disappeared**

First verify current refs:

```bash
for f in \
  adr-009-citations.md \
  adr-013-citations.md \
  adr-07-twodigit-citations.md \
  adr-08-twodigit-citations.md; do
  rg -n --fixed-strings "$f" docs/architecture/decisions || true
done
```

After retiring ADR 010/018/026/030, expected current retained consumers are zero. Then:

```bash
git rm \
  docs/architecture/decisions/_citations/adr-009-citations.md \
  docs/architecture/decisions/_citations/adr-013-citations.md \
  docs/architecture/decisions/_citations/adr-07-twodigit-citations.md \
  docs/architecture/decisions/_citations/adr-08-twodigit-citations.md
```

Keep `_citations/adr-017-citations.md` and `RENUMBERING-REGISTRY.md` while retained ADR 017/034 still need them.

- [ ] **Step 4: Retire completed old plans with no live consumer**

For each of the four current `docs/plans/*.md` files, use its Task 1 audit row. Expected completed candidates are:

```text
docs/plans/2026-08-22-op-read-01-repair.md
docs/plans/2026-08-23-authorization-request-d5-wire.md
docs/plans/2026-08-23-notif-d6-r-b00-r2.md
docs/plans/2026-08-23-notif-d6-r-b11.md
```

If their accepted outputs are already in current owners/evidence and no active current doc requires the exact plan, retire them:

```bash
git rm docs/plans/2026-08-22-op-read-01-repair.md \
       docs/plans/2026-08-23-authorization-request-d5-wire.md \
       docs/plans/2026-08-23-notif-d6-r-b00-r2.md \
       docs/plans/2026-08-23-notif-d6-r-b11.md
```

If Task 1 classified any KEEP/REHOME, follow that row instead of deleting blindly.

- [ ] **Step 5: Verify ADR/planning reachability**

```bash
rg -n '008-production-deploy-topology|010-mercado-livre-polling-visible-refresh|018-mutation-envelope-table-and-poller|026-scheduler-phase-vocabulary|030-scheduler-second-instance-per-installation' AGENTS.md ARCHITECTURE.md docs || true
rg -n '2026-08-22-op-read-01-repair|2026-08-23-authorization-request-d5-wire|2026-08-23-notif-d6-r-b00-r2|2026-08-23-notif-d6-r-b11' AGENTS.md ARCHITECTURE.md docs || true
```

Any remaining current reference must be repointed to the replacement owner or prove the file should not have been retired.

- [ ] **Step 6: Commit**

```bash
git add docs/architecture/decisions docs/plans docs/engineering/rebaseline
git add -u docs/architecture/decisions docs/plans
git commit -m "docs: retire satisfied ADR and planning residue"
```

---

## Task 5: Make the single aggregate gate proportional and separate current Product proof from historical replay

**Files:**
- Modify: `scripts/gate.ps1`
- Modify: `scripts/verify-product-oad.mjs`
- Create: `scripts/verify-product-oad-history.mjs`
- Preserve current targeted scripts listed in File Map
- Expected no change: `.github/workflows/ci.yml`, `contracts/api/product/**`

**Interfaces:**
- Consumes: existing `GATE_BASE_SHA/GATE_HEAD_SHA` changed-file logic in `gate.ps1` and current full replay in `verify-product-oad.mjs`.
- Produces: one required gate with three internal proof classes: universal cheap checks, current Product proof when affected, historical replay only when historical machinery/authority changes.

- [ ] **Step 1: Capture the existing full Product replay as the targeted history verifier before simplifying it**

Run from the current health branch:

```bash
cp scripts/verify-product-oad.mjs scripts/verify-product-oad-history.mjs
```

In `verify-product-oad-history.mjs`, change only user-visible proof labels needed to make its purpose explicit; retain the existing 106 → historical 99 → pre-auth → baseline 95 replay behavior. Do not change the historical semantic assertions in this task.

- [ ] **Step 2: Convert `verify-product-oad.mjs` to current-contract proof only**

Remove from the current verifier:

```text
historicalVerifier/preAuthVerifier/baselineVerifier fixture-copy dependencies
NOTIF path stripping used only to reconstruct old surfaces
rewindAuthorizationRequestSurface()
historical99Proof()
try { historical99Proof(); currentAuthProof(); }
```

Keep and strengthen the current proof path:

```text
Redocly lint
2 deterministic current bundles
current 106-operation uniqueness/presence
split H vs A/S authentication profile
exact current 31 ordinary Permission vocabulary consistency
2 deterministic TypeScript generations + strict compilation
2 deterministic Go generations + Go test
current Product source/generation claims already enforced by the current script
```

Add direct current Permission checks instead of deriving confidence from old 95/99 replay:

```js
const ordinary = normalize(all
  .map((entry) => entry.operation['x-mpc-required-permission'])
  .filter((value) => value && value !== 'authenticated'));
assert(ordinary.length === 31, `current ordinary Permission vocabulary must contain 31 values, found ${ordinary.length}`);
```

Resolve the effective AccessContext and AccessRole Permission enums and assert they equal the same 31-value set. Keep explicit presence checks for current NOTIF and AuthorizationRequest operation IDs already in the script.

The current script must end by running **only** current proof.

- [ ] **Step 3: Have the current Product proof invoke current semantic slices, not historical migration reconstruction**

After the current core lint/bundle/generation proof succeeds, invoke these current-only semantic proofs:

```text
scripts/verify-oad-source-reachability.mjs
scripts/verify-operational-read-contract.mjs
scripts/verify-performance-evidence-knowledge.mjs
scripts/verify-notification-oad.mjs
scripts/verify-authorization-request-oad.mjs
```

Use `process.execPath` child execution and fail on non-zero exit. These proofs all evaluate the current 106-operation OAD; none reconstructs the 95/99 migration chain.

Do not add wireframe/prose verifiers to Product proof.

- [ ] **Step 4: Add explicit diff predicates to `scripts/gate.ps1`**

Keep all current universal checks. After `$changedFiles` is known, add two functions/booleans.

Current Product proof is affected when any changed path matches one of:

```text
^contracts/api/product/
^scripts/gate\.ps1$
^scripts/verify-product-oad\.mjs$
^scripts/verify-oad-source-reachability\.mjs$
^scripts/lib/publication-requirements-oad-proof\.mjs$
^scripts/verify-operational-read-contract\.mjs$
^scripts/verify-performance-evidence-knowledge\.mjs$
^scripts/verify-notification-oad\.mjs$
^scripts/verify-authorization-request-oad\.mjs$
^\.github/workflows/ci\.yml$
```

Historical replay is affected when any changed path matches one of:

```text
^scripts/verify-product-oad-history\.mjs$
^scripts/verify-product-oad-current99\.mjs$
^scripts/verify-product-oad-pre-auth\.mjs$
^scripts/verify-product-oad-baseline\.mjs$
^docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR\.md$
^docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT\.md$
^docs/engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX\.md$
^docs/engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE\.md$
```

If no diff base is available, fail safe by running the current Product proof rather than silently skipping it.

Execution behavior:

```powershell
if ($productProofAffected) {
    & node 'scripts/verify-product-oad.mjs'
    # fail on non-zero
    Write-Host 'product_oad_proof: PASS'
} else {
    Write-Host 'product_oad_proof: SKIPPED_NOT_AFFECTED'
}

if ($historicalReplayAffected) {
    & node 'scripts/verify-product-oad-history.mjs'
    # fail on non-zero
    Write-Host 'product_oad_history_replay: PASS'
} else {
    Write-Host 'product_oad_history_replay: SKIPPED_NOT_AFFECTED'
}
```

Do not add another GitHub job/check.

- [ ] **Step 5: Prove current and historical verifiers independently**

Run:

```powershell
node scripts/verify-product-oad.mjs
node scripts/verify-product-oad-history.mjs
```

Expected:

```text
current verifier → current 106/31/H-A-S + current semantic/generator proofs PASS
history verifier → historical 95/99/pre-auth replay PASS
```

- [ ] **Step 6: Prove docs-only routing skips Product generation**

Use the known pre-CI health commit as a synthetic docs-only diff. Resolve it by commit message from Task 4:

```bash
DOCS_ONLY_HEAD=$(git log -1 --format=%H --grep='docs: retire satisfied ADR and planning residue')
BASE=$(git merge-base main "$DOCS_ONLY_HEAD")
```

Run the **current working-tree gate script** with those diff endpoints:

```powershell
$env:GATE_BASE_SHA = '<BASE from command above>'
$env:GATE_HEAD_SHA = '<DOCS_ONLY_HEAD from command above>'
npm run gate
```

Expected output includes:

```text
product_oad_proof: SKIPPED_NOT_AFFECTED
product_oad_history_replay: SKIPPED_NOT_AFFECTED
```

and the cheap universal gate still passes.

- [ ] **Step 7: Prove the health PR itself exercises the changed proof machinery**

Unset the synthetic endpoints and run against the actual health diff:

```powershell
Remove-Item Env:GATE_BASE_SHA -ErrorAction SilentlyContinue
Remove-Item Env:GATE_HEAD_SHA -ErrorAction SilentlyContinue
npm run gate
```

Because the health PR changes gate/current/history proof machinery, expected output includes current Product proof PASS and historical replay PASS at least once during this PR.

- [ ] **Step 8: Prove Product contract bytes/meaning were not changed by health work**

```bash
git diff --exit-code main...HEAD -- contracts/api/product
git diff --name-only main...HEAD -- contracts/api/product
```

Expected: no Product contract files changed.

- [ ] **Step 9: Commit CI proportionality**

```bash
git add scripts/gate.ps1 scripts/verify-product-oad.mjs scripts/verify-product-oad-history.mjs
git commit -m "ci: make Product proof proportional to changed claims"
```

---

## Task 6: Run repository-health coherence checks and independent review

**Files:**
- Modify only to fix validated findings: files changed in Tasks 2–5
- Modify: `docs/roadmap.md` when the health candidate is proven

**Interfaces:**
- Consumes: cleaned active tree + proportional gate.
- Produces: integration-ready health candidate with Product/frontend evidence preserved.

- [ ] **Step 1: Verify active-tree size/routing improvement without setting an arbitrary quota**

Run:

```bash
printf 'rebaseline_files='; git ls-files 'docs/engineering/rebaseline/**' | wc -l
printf 'notif_chain_files='; git ls-files 'docs/engineering/rebaseline/D6-R2-NOTIF-01-*' | wc -l
printf 'planning_files='; { git ls-files 'docs/plans/**'; git ls-files 'docs/superpowers/**'; } | wc -l
printf 'remote_history_pack_refs='; rg -l 'D6-R2-AUTHORITY-ROUTE' AGENTS.md docs ARCHITECTURE.md | wc -l
```

Do not assert a target number. Review whether every surviving file has a current authority/evidence reason from Task 1.

- [ ] **Step 2: Run representative fresh-session navigation checks manually**

Using only `AGENTS.md → docs/roadmap.md → method/index`, prove these questions resolve without a historical route pack:

```text
Who owns Product identity?                         → D2
Where are publication/ListingIntent semantics?     → D4-R1 + relevant D5/OAD owner
Where is B20 current screen inventory?             → P5 + roadmap
Where is current runtime authority?                → D7
Where is current B10 acceptance evidence?          → block-specific P8/P9 + HTML
```

If any route requires reading a ratification chain just to discover the current owner, fix `docs/index.md`; do not recreate an everything-router.

- [ ] **Step 3: Run link/reference safety checks for retired files**

Use the committed Task 1 audit rows with `Action=rm/rehome+rm` and verify every deleted basename has zero current references outside this temporary plan. Also run:

```bash
git diff --check main...HEAD
git status --short
```

- [ ] **Step 4: Run the final required gate and targeted Product history proof**

```powershell
npm run gate
node scripts/verify-product-oad-history.mjs
```

Expected: PASS. The targeted second command proves history remains available independent of required CI routing.

- [ ] **Step 5: Request a fresh full PR review**

On PR #71, request:

```text
@coderabbitai full review
```

Review against:

```text
approved repository-health spec
no Product/OAD semantic change
no current-authority loss
no broken current links
rehome-before-retire discipline
one required CI check
current Product proof preserved
historical replay available but not default for unrelated docs
```

- [ ] **Step 6: Adjudicate review findings**

- fix every Critical/Important finding that is valid;
- reject incorrect findings with concrete repository evidence;
- if a finding proves a retired file was the only current owner of material semantics, restore/rehome it before proceeding;
- rerun `npm run gate` after fixes.

- [ ] **Step 7: Update roadmap to integration-ready health state**

Set current health result to:

```text
Repository Governance & Context Health Rebaseline — PROVED / INTEGRATION CANDIDATE
PR #69 — PAUSED / NO P8
PR #70 — PAUSED pending reanchor after health integration
Product 106/31/H-A-S — unchanged
Runtime NONE — unchanged
exact next action — operator reviews/authorizes health PR merge; after integration reanchor #70, then #69
```

- [ ] **Step 8: Commit review/closure fixes**

```bash
git add AGENTS.md ARCHITECTURE.md docs scripts .github
# Do not add unrelated/untracked diagnostics under .git; Git ignores them by location.
git commit -m "docs: close repository health rebaseline"
```

If there is nothing to commit after review, do not create an empty ceremony commit.

---

## Task 7: Retire the health working spec/plan themselves before final merge review

**Files:**
- Delete: `docs/superpowers/specs/2026-08-25-repository-governance-context-health-design.md`
- Delete: `docs/superpowers/plans/2026-08-25-repository-governance-context-health-implementation-plan.md`
- Modify if needed: `docs/index.md`, `docs/roadmap.md`, `docs/development/engineering-rules.md` to ensure no current link depends on the temporary files

**Interfaces:**
- Consumes: durable rules already rehomed into AGENTS/index/engineering-rules/roadmap and the accepted PR diff.
- Produces: a health PR that does not leave its own temporary planning artifacts in the active tree.

- [ ] **Step 1: Prove every durable health rule has a current owner**

Verify:

```text
selective bootstrap/context expansion  → AGENTS.md
current status/next action             → docs/roadmap.md
current-owner navigation               → docs/index.md
artifact lifecycle + Git-as-history    → docs/development/engineering-rules.md (+ existing ARCHITECTURE invariant)
CI proportionality                     → engineering-rules + scripts/gate.ps1
current vs historical Product proof    → scripts/verify-product-oad.mjs / scripts/verify-product-oad-history.mjs
```

If any durable rule exists only in the health spec/plan, rehome it before deletion.

- [ ] **Step 2: Remove current references to the temporary health spec/plan**

```bash
rg -n '2026-08-25-repository-governance-context-health-(design|implementation-plan)' AGENTS.md ARCHITECTURE.md docs scripts .github || true
```

Update roadmap/index so they refer to PR #71 / current health state rather than requiring the temporary files after close.

- [ ] **Step 3: Delete the temporary health artifacts**

```bash
git rm docs/superpowers/specs/2026-08-25-repository-governance-context-health-design.md
git rm docs/superpowers/plans/2026-08-25-repository-governance-context-health-implementation-plan.md
```

If the directories become empty, remove the empty directories naturally; Git tracks files, not empty directories.

- [ ] **Step 4: Final verification-before-completion**

Run:

```bash
git diff --check main...HEAD
git diff --exit-code main...HEAD -- contracts/api/product
git status --short
```

and:

```powershell
npm run gate
node scripts/verify-product-oad-history.mjs
```

Expected:

- current Product contract diff: empty;
- one required gate path remains;
- current proof passes when affected by proof machinery;
- historical replay passes targeted;
- no health spec/plan remains in active tree;
- no retired current link remains.

- [ ] **Step 5: Commit working-artifact retirement**

```bash
git add docs/roadmap.md docs/index.md docs/development/engineering-rules.md
git add -u docs/superpowers
git commit -m "docs: retire repository health working artifacts"
```

- [ ] **Step 6: Stop at merge authorization**

Do **not** merge PR #71 automatically. Report final diff, CI/review state, retired/current owner summary and ask for explicit operator merge authorization.

After merge, perform remote branch cleanup as a separate audited operation and then reanchor PR #70 / PR #69 from the cleaned `main`.
