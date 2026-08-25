# Repository Governance & Context Health Rebaseline Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restore Marketplace Central to a small current-authority graph, retire absorbed intermediate documentation safely, and make the single required CI gate proportional to the changed claim without weakening Product protection.

**Architecture:** Keep the proven local route `AGENTS → roadmap → applicable method → index → smallest current owner`, with Git history as archive. Retirement is evidence-driven: classify first, rehome surviving meaning second, delete only after replacement/current-consumer proof. For CI, implementation planning found that the existing 95/99 replay still carries valid current contract protections, so this plan does **not** split or weaken that verifier; it makes the heavy proof diff-aware so unrelated frontend/docs work skips it. A future split is a separate reopen only if Product-changing CI itself becomes a material bottleneck and proof parity can be demonstrated.

**Tech Stack:** Git/GitHub; Markdown authority; PowerShell `scripts/gate.ps1`; Node.js `>=26.3.0 <27`; existing OpenAPI/Redocly/TypeScript/Go Product proof toolchain.

**Spec:** `docs/superpowers/specs/2026-08-25-repository-governance-context-health-design.md`

## Global Constraints

- Preserve current Product semantics: **106 Product operations / 31 ordinary Permissions / Principal kinds H/A/S**.
- Preserve runtime baseline **NONE** and D9 implementation block.
- Preserve local `engineering-method.md`, `frontend-product-experience-planning-method.md`, and `engineering-rules.md`; do not restore an external methodology pin/router/profile layer.
- `AGENTS.md` stays the small bootstrap; `docs/roadmap.md` remains sole mutable current-status authority; `docs/index.md` routes current questions to current owners.
- Rehome surviving meaning before retirement. Filename, age, prefix, or “ratification” naming never proves safe deletion.
- Git history is the archive. Do not create `docs/archive/`, duplicate historical trees, rewrite history, force-push, reset, clean, or stash by convenience.
- Preserve current evidence with a live consumer, especially operator-locked D6-R2 HTML and current B10 evidence.
- PR #69 remains **PAUSED / NO P8**. PR #70 remains paused at its current implementation-plan gate.
- No Product/OAD semantic change belongs in this health PR. `contracts/api/product/**` must remain unchanged from `main`.
- Keep exactly one required GitHub check: `required → npm run gate`.
- The existing heavy Product proof remains unchanged semantically in this health PR. Optimization is routing, not weaker verification.
- Remote branch deletion is a separate post-merge audit.
- Merge requires explicit operator authorization.

---

## Task 1: Build the path-by-path active-tree retirement map before deleting anything

**Files:**
- Modify: `docs/superpowers/plans/2026-08-25-repository-governance-context-health-implementation-plan.md` only to append the execution audit table
- Read: repository docs/scripts candidates generated below

**Produces:** the exact deletion/rehoming input for Tasks 2–4.

- [ ] **Step 1: Reanchor and prove Product files are untouched**

```bash
git status --short
git rev-parse HEAD
git rev-parse main
git merge-base main HEAD
git diff --name-status main...HEAD
git diff --exit-code main...HEAD -- contracts/api/product
```

Expected: no Product contract diff.

- [ ] **Step 2: Generate the candidate inventory under `.git/`**

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

`.git/repository-health-candidates.txt` is diagnostic and is never committed.

- [ ] **Step 3: Gather incoming references and headers for every candidate**

```bash
while IFS= read -r path; do
  test -n "$path" || continue
  printf '\n=== %s ===\n' "$path"
  base=$(basename "$path")
  rg -n --fixed-strings "$base" AGENTS.md ARCHITECTURE.md docs contracts scripts qualification .github \
    --glob '!docs/superpowers/plans/2026-08-25-repository-governance-context-health-implementation-plan.md' || true
  sed -n '1,100p' "$path"
done < .git/repository-health-candidates.txt
```

For each candidate, read enough of the file and the claimed replacement/current consumer to answer:

```text
Does it still own current meaning?
Is it current evidence for a live obligation?
If meaning survives, which exact current owner already owns it or must receive it?
Would removal lose current truth, or only investigation/history?
```

- [ ] **Step 4: Append the execution audit table to this temporary plan**

Append `## Execution Audit — Active Tree` with one row per candidate:

```markdown
| Path | Class | Surviving meaning | Replacement owner / live consumer | Action |
| --- | --- | --- | --- | --- |
```

Allowed classes/actions:

```text
KEEP AUTHORITY          → keep
KEEP CURRENT EVIDENCE   → keep
REHOME THEN RETIRE      → rehome+rm
RETIRE                   → rm
```

Safety rules:

- `RETIRE` requires no surviving current meaning, or an exact current replacement owner already containing it.
- `REHOME THEN RETIRE` requires an exact existing owner. If none exists, STOP: cleanup must not invent Product architecture.
- Current operator-locked HTML is `KEEP CURRENT EVIDENCE` unless separately adjudicated by the frontend method.
- Each row states an exact reason; no filename-only or vague age-based rationale.

- [ ] **Step 5: Validate the audit and commit the no-delete checkpoint**

Review all `rm` / `rehome+rm` rows manually, then:

```bash
git diff --check
git add docs/superpowers/plans/2026-08-25-repository-governance-context-health-implementation-plan.md
git commit -m "docs: classify active repository authority and history"
```

No repository file is retired before this commit.

---

## Task 2: Restore the small bootstrap graph and current-owner index

**Files:**
- Modify: `AGENTS.md`
- Modify: `docs/index.md`
- Modify: `docs/development/engineering-rules.md`
- Modify: `docs/roadmap.md`
- Delete after direct routes are proven: `docs/engineering/rebaseline/D6-R2-AUTHORITY-ROUTE.md`

**Produces:** `AGENTS → roadmap → method → index → current owner` with no “everything remains reachable” route pack.

- [ ] **Step 1: Add one selective-context rule to AGENTS**

Preserve the current bootstrap and add:

```markdown
Start with the smallest current owner set likely to answer the task. Expand into additional owners, Evidence, Git history, runtime or external sources only when they can materially change or falsify the conclusion. “No fixed context budget” permits expansion; it does not make historical artifact packs part of the default read path.
```

Do not add profiles, task packs, external methodology refs, generated routing, or file-count limits.

- [ ] **Step 2: Make `docs/index.md` route directly to current owners**

Keep the existing D0/D1/D2/D3/D4/D5/D6/D7/D8 and method routes. Replace the D6-R2 authority-pack route with direct entries:

```text
frontend architecture/client state/topology
→ D6-FRONTEND.md + Frontend Method

complete current screen/surface inventory
→ D6-R2-P5-SCREEN-SURFACE-INVENTORY.md

B10 current acceptance evidence
→ D6-R2-P8-B10-PREPARATION-RATIFICATION.md
→ D6-R2-P9-B10-PREPARATION-SCREEN-CONTRACT.md
→ qualification/d6-r2-wireframes/b10-preparation.html

current block/program status
→ roadmap only
```

For approvals/notifications, route only to Task 1 artifacts classified as current authority/evidence; do not enumerate the full NOTIF chain.

- [ ] **Step 3: Add the durable artifact lifecycle to `engineering-rules.md`**

Add this compact repository rule:

```text
current owner/router/method                              KEEP
current evidence with a named live consumer             KEEP while live
accepted intermediate with unique surviving meaning     REHOME once, then retire
fully absorbed/superseded intermediate                  RETIRE
history                                                  Git
```

Also state:

- plans/specs/reviews/ratifications do not remain parallel semantic authority after acceptance;
- no `docs/archive/` baseline;
- rehome before retire;
- index/routers point to current owners rather than historical chains.

- [ ] **Step 4: Put the health increment into the branch-local roadmap**

Preserve Product and D6-R2 status, but set current sequencing to:

```text
Repository Governance & Context Health Rebaseline — EXECUTION ACTIVE
PR #69 — PAUSED / NO P8
PR #70 — PAUSED at implementation-plan gate
Product 106/31/H-A-S — unchanged
runtime NONE — unchanged
next action — finish health proof/integration, then reanchor #70 and #69
```

Do not advance D9 or implementation.

- [ ] **Step 5: Remove the stage-wide authority route pack only after references are gone**

```bash
rg -n --fixed-strings 'D6-R2-AUTHORITY-ROUTE.md' AGENTS.md ARCHITECTURE.md docs contracts scripts qualification .github || true
```

After `docs/index.md` no longer depends on it, expected current consumers are zero. Then:

```bash
git rm docs/engineering/rebaseline/D6-R2-AUTHORITY-ROUTE.md
```

- [ ] **Step 6: Verify and commit**

```bash
git diff --check
rg -n 'conexus-methodology|methodology profile|authority pack' AGENTS.md docs/index.md docs/development/engineering-rules.md || true
git add AGENTS.md docs/index.md docs/development/engineering-rules.md docs/roadmap.md
git add -u docs/engineering/rebaseline/D6-R2-AUTHORITY-ROUTE.md
git commit -m "docs: restore selective repository authority routing"
```

Any matching text must be an explicit prohibition/history statement, not a live dependency.

---

## Task 3: Rehome and retire absorbed D6-R2 intermediate material

**Files:**
- Modify only when required by Task 1 `REHOME THEN RETIRE` rows: existing owners among
  - `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
  - `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
  - `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
  - `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`
  - `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
  - `docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md`
  - `docs/engineering/rebaseline/D5-API.md`
  - `docs/engineering/rebaseline/D5-B2-WIRE-CONTRACT.md` (W1)
  - `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md`
  - `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md`
  - `docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md`
  - `docs/engineering/rebaseline/D6-FRONTEND.md`
  - `docs/engineering/rebaseline/D6-R2-P5-SCREEN-SURFACE-INVENTORY.md`
  - an exact retained block-specific P8/P9 owner named by the audit
  - `docs/engineering/rebaseline/D7-RUNTIME-JOBS-TRANSACTIONS.md` or its accepted D7-B/C/D/E owner when runtime meaning survives
  - `docs/engineering/rebaseline/D8-GOLDEN-FLOWS.md` when golden-flow meaning survives
- Delete: only Task 1 rows whose Action is `rm` or `rehome+rm`
- Preserve: Task 1 KEEP rows and locked/current HTML

**Produces:** current D6-R2 work no longer depends on intermediate findings/ratifications/review chains.

- [ ] **Step 1: Rehome each surviving conclusion once**

For every `REHOME THEN RETIRE` row, add only the accepted normative conclusion to its exact owner. Do not copy investigation chronology, reviewer dialogue, obsolete counts/status, or ratification ceremony.

If rehome would require a new Product owner, operation, Permission, runtime behavior, or material architecture change, STOP and surface a separate planning finding.

- [ ] **Step 2: Adjudicate the two accumulated/stale D6-R2 ledgers**

Apply the Task 1 result to:

```text
D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md
D6-R2-P8-BLOCK-LEDGER.md
```

If current unique meaning is already carried by D6/P5/block-specific P8/P9/roadmap, retire them. If unique current meaning survives, rehome it first. Do not rewrite old 99/30/v2.1 snapshots merely to keep these files alive.

- [ ] **Step 3: Build the exact delete list from the committed audit and retire only those paths**

Create an untracked file under `.git/` containing one exact `rm`/`rehome+rm` path per line:

```bash
: > .git/repository-health-retire.txt
cat .git/repository-health-retire.txt
while IFS= read -r path; do
  test -n "$path" || continue
  git rm -- "$path"
done < .git/repository-health-retire.txt
```

The executor fills `.git/repository-health-retire.txt` only from the committed audit table; no glob-based delete is allowed.

- [ ] **Step 4: Prove deleted files have no current incoming references**

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

Expected: zero current stale refs.

- [ ] **Step 5: Prove frontend evidence was not changed accidentally**

```bash
git diff --exit-code main -- qualification/d6-r2-wireframes
```

Expected: no HTML evidence changed.

- [ ] **Step 6: Commit**

```bash
git add docs/engineering/rebaseline docs/index.md docs/roadmap.md
git add -u docs/engineering/rebaseline
git commit -m "docs: retire absorbed D6-R2 intermediate material"
```

---

## Task 4: Retire satisfied ADR residue and completed old plans

**Files:**
- Modify: `docs/architecture/decisions/README.md`
- Candidate deletes after replacement verification:
  - `008-production-deploy-topology.md`
  - `010-mercado-livre-polling-visible-refresh.md`
  - `018-mutation-envelope-table-and-poller.md`
  - `026-scheduler-phase-vocabulary.md`
  - `030-scheduler-second-instance-per-installation.md`
  - `_citations/adr-009-citations.md`
  - `_citations/adr-013-citations.md`
  - `_citations/adr-07-twodigit-citations.md`
  - `_citations/adr-08-twodigit-citations.md`
- Candidate completed plans:
  - `docs/plans/2026-08-22-op-read-01-repair.md`
  - `docs/plans/2026-08-23-authorization-request-d5-wire.md`
  - `docs/plans/2026-08-23-notif-d6-r-b00-r2.md`
  - `docs/plans/2026-08-23-notif-d6-r-b11.md`
- Preserve: ADR 017/034/035 plus their still-consumed provenance

**Produces:** ADR registry and planning tree whose active files still have real consumers.

- [ ] **Step 1: Verify each D7-conditioned ADR has a current replacement owner**

Use these owner checks:

```text
008 → D7-E / accepted deployment-operability authority
010 → D7 runtime + D4 acquisition authority
018 → D7-C + accepted D3/D4 external-effect safety
026 → D7 runtime/jobs authority
030 → D7 runtime/jobs authority
```

If an ADR contains unique current target meaning absent from its replacement owner, rehome exactly that meaning first. A material architecture contradiction stops this cleanup.

- [ ] **Step 2: Retire satisfied ADRs and update the registry**

```bash
git rm \
  docs/architecture/decisions/008-production-deploy-topology.md \
  docs/architecture/decisions/010-mercado-livre-polling-visible-refresh.md \
  docs/architecture/decisions/018-mutation-envelope-table-and-poller.md \
  docs/architecture/decisions/026-scheduler-phase-vocabulary.md \
  docs/architecture/decisions/030-scheduler-second-instance-per-installation.md
```

Update `docs/architecture/decisions/README.md` so these are no longer active residues; state that accepted D7 satisfied their retirement condition and Git history retains them.

- [ ] **Step 3: Remove citation files whose last retained consumer disappeared**

```bash
for f in adr-009-citations.md adr-013-citations.md adr-07-twodigit-citations.md adr-08-twodigit-citations.md; do
  rg -n --fixed-strings "$f" docs/architecture/decisions || true
done
```

After the five ADR removals, zero retained consumers are expected for these four citation files. Then:

```bash
git rm \
  docs/architecture/decisions/_citations/adr-009-citations.md \
  docs/architecture/decisions/_citations/adr-013-citations.md \
  docs/architecture/decisions/_citations/adr-07-twodigit-citations.md \
  docs/architecture/decisions/_citations/adr-08-twodigit-citations.md
```

Keep `adr-017-citations.md` and `RENUMBERING-REGISTRY.md` while retained 017/034 consume them.

- [ ] **Step 4: Apply the Task 1 audit to the four old `docs/plans` files**

For rows classified RETIRE after output rehome/proof, remove the exact files. If any row is KEEP/REHOME, follow that audit row instead of deleting by convention.

- [ ] **Step 5: Verify references and commit**

```bash
rg -n '008-production-deploy-topology|010-mercado-livre-polling-visible-refresh|018-mutation-envelope-table-and-poller|026-scheduler-phase-vocabulary|030-scheduler-second-instance-per-installation' AGENTS.md ARCHITECTURE.md docs || true
git diff --check
git add docs/architecture/decisions docs/plans docs/engineering/rebaseline
git add -u docs/architecture/decisions docs/plans
git commit -m "docs: retire satisfied ADR and planning residue"
```

Any surviving current reference must be repointed to the replacement owner before commit.

---

## Task 5: Make the existing heavy Product proof diff-aware without weakening it

**Files:**
- Modify: `scripts/gate.ps1`
- Keep semantically unchanged: `scripts/verify-product-oad.mjs` and its current 95/99/pre-auth/current proof chain
- Expected no change: `.github/workflows/ci.yml`, `contracts/api/product/**`

**Produces:** docs/frontend planning changes run only cheap universal checks; Product/proof-input changes still run the exact heavy proof that is green today.

### Implementation-planning finding

The existing 95/99 replay is not merely historical ceremony: it currently protects still-valid baseline invariants such as operation/Permission mapping, idempotency, ETag carriers, Product Problems and generated route behavior. Therefore this health increment does **not** split it until a future dedicated proof refactor can demonstrate current-contract parity. This preserves the approved health invariant “do not weaken current Product proof”.

- [ ] **Step 1: Add an exact Product-proof changed-path predicate to `gate.ps1`**

Keep all existing cheap checks. Replace the unconditional final Product proof call with this shape:

```powershell
function Test-ChangedPathMatches([string[]]$patterns) {
    foreach ($file in $changedFiles) {
        $normalized = $file.Replace('\', '/')
        foreach ($pattern in $patterns) {
            if ($normalized -match $pattern) { return $true }
        }
    }
    return $false
}

$productProofPatterns = @(
    '^contracts/api/product/',
    '^scripts/gate\.ps1$',
    '^scripts/verify-product-oad\.mjs$',
    '^scripts/verify-product-oad-current99\.mjs$',
    '^scripts/verify-product-oad-pre-auth\.mjs$',
    '^scripts/verify-product-oad-baseline\.mjs$',
    '^scripts/verify-oad-source-reachability\.mjs$',
    '^scripts/lib/publication-requirements-oad-proof\.mjs$',
    '^scripts/verify-operational-read-contract\.mjs$',
    '^scripts/verify-performance-evidence-knowledge\.mjs$',
    '^scripts/verify-notification-oad\.mjs$',
    '^scripts/verify-authorization-request-oad\.mjs$',
    '^docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR\.md$',
    '^docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT\.md$',
    '^docs/engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX\.md$',
    '^docs/engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE\.md$',
    '^\.github/workflows/ci\.yml$'
)

$productProofAffected = if (-not $base) { $true } else { Test-ChangedPathMatches $productProofPatterns }

if ($productProofAffected) {
    $productProof = & node 'scripts/verify-product-oad.mjs' 2>&1
    $productProofExit = $LASTEXITCODE
    $productProof | ForEach-Object { Write-Host $_ }
    if ($productProofExit -ne 0) { Fail 'Product OAD proof failed' }
    Write-Host 'product_oad_proof: PASS'
} else {
    Write-Host 'product_oad_proof: SKIPPED_NOT_AFFECTED'
}
```

Fail-safe law: when a diff base is unavailable, run the Product proof.

- [ ] **Step 2: Prove a synthetic docs-only range skips the heavy proof**

Use the Task 4 commit as a known docs-only endpoint:

```bash
git log -1 --format=%H --grep='docs: retire satisfied ADR and planning residue'
git merge-base main "$(git log -1 --format=%H --grep='docs: retire satisfied ADR and planning residue')"
```

Copy the printed merge-base SHA and Task 4 SHA into:

```powershell
$env:GATE_BASE_SHA = 'printed-merge-base-sha'
$env:GATE_HEAD_SHA = 'printed-task4-sha'
npm run gate
```

Expected:

```text
product_oad_proof: SKIPPED_NOT_AFFECTED
gate: PASS
```

The universal required-file/workflow/diff/implementation-block checks still execute.

- [ ] **Step 3: Prove the health diff still runs the full existing proof**

Clear the synthetic endpoints and run:

```powershell
Remove-Item Env:GATE_BASE_SHA -ErrorAction SilentlyContinue
Remove-Item Env:GATE_HEAD_SHA -ErrorAction SilentlyContinue
npm run gate
```

Because `scripts/gate.ps1` changed and is in the predicate, expected output still contains the existing 95/99/106 proof evidence and ends with:

```text
product_oad_proof: PASS
gate: PASS
```

- [ ] **Step 4: Prove Product bytes stayed unchanged**

```bash
git diff --exit-code main...HEAD -- contracts/api/product
```

Expected: empty diff.

- [ ] **Step 5: Commit**

```bash
git add scripts/gate.ps1
git commit -m "ci: make Product proof proportional to changed paths"
```

Reopen proof decomposition only if Product-changing CI itself becomes a material bottleneck; that future work must map every still-current 95/99 assertion to equivalent current-authority proof before removing replay.

---

## Task 6: Run global repository-health coherence and independent review

**Files:**
- Modify only to fix validated findings from Tasks 2–5
- Modify: `docs/roadmap.md` for integration-candidate status

**Produces:** a proved health candidate with current Product/frontend evidence intact.

- [ ] **Step 1: Measure the active tree after cleanup without imposing a quota**

```bash
printf 'rebaseline_files='; git ls-files 'docs/engineering/rebaseline/**' | wc -l
printf 'notif_chain_files='; git ls-files 'docs/engineering/rebaseline/D6-R2-NOTIF-01-*' | wc -l
printf 'old_plan_files='; git ls-files 'docs/plans/**' | wc -l
printf 'authority_route_refs='; rg -l 'D6-R2-AUTHORITY-ROUTE' AGENTS.md docs ARCHITECTURE.md | wc -l
```

Review every survivor against the Task 1 audit. No numerical target is enforced.

List the largest remaining current docs/authority files:

```bash
for f in $(git ls-files '*.md' '*.yaml' '*.yml'); do
  test -f "$f" && printf '%10s %s\n' "$(wc -c < "$f")" "$f"
done | sort -nr | head -20
```

Do not split a large canonical owner in this PR. If accidental history is gone but a canonical file still repeatedly causes task-level overflow because it owns independently navigable responsibilities, record a separate reopen after health integration.

- [ ] **Step 2: Prove representative fresh-session navigation**

Using only `AGENTS.md → roadmap → applicable method/index`, verify:

```text
Product identity                              → D2
publication / ListingIntent semantics         → D4-R1 + relevant D5/OAD owner
B20 screen inventory                          → P5 + roadmap
runtime                                       → D7
B10 current acceptance evidence               → retained block-specific P8/P9 + HTML
```

If discovery requires a ratification chain merely to find the owner, fix `docs/index.md`; do not recreate a historical router.

- [ ] **Step 3: Recheck protected Product/frontend evidence**

```bash
git diff --check main...HEAD
git diff --exit-code main...HEAD -- contracts/api/product
git diff --exit-code main -- qualification/d6-r2-wireframes
git status --short
```

Expected: Product and current wireframe evidence unchanged.

- [ ] **Step 4: Run final CI**

```powershell
npm run gate
```

Expected: PASS. Because `scripts/gate.ps1` changed in this PR, the full Product proof runs here once.

- [ ] **Step 5: Request a fresh full review on PR #71**

Post:

```text
@coderabbitai full review
```

Review criteria:

```text
no current-authority loss
rehome-before-retire
no Product/OAD semantic change
no locked frontend evidence change
selective bootstrap/index
one required CI check
unrelated docs skip heavy Product proof
Product/proof-input changes retain full existing proof
```

- [ ] **Step 6: Adjudicate findings and update roadmap**

Fix valid Critical/Important findings. If review proves a retired file was the only current owner of material meaning, restore/rehome before proceeding. Rerun `npm run gate` after fixes.

Set branch-local roadmap result to:

```text
Repository Governance & Context Health Rebaseline — PROVED / INTEGRATION CANDIDATE
PR #69 — PAUSED / NO P8
PR #70 — PAUSED pending reanchor after health integration
Product 106/31/H-A-S — unchanged
runtime NONE — unchanged
next action — operator reviews/authorizes PR #71 merge; after integration reanchor #70, then #69
```

Commit only if review/roadmap changes exist:

```bash
git add AGENTS.md ARCHITECTURE.md docs scripts .github
git commit -m "docs: close repository health rebaseline"
```

Do not create an empty closure commit.

---

## Task 7: Retire the health spec/plan themselves and stop at merge authorization

**Files:**
- Delete: `docs/superpowers/specs/2026-08-25-repository-governance-context-health-design.md`
- Delete: `docs/superpowers/plans/2026-08-25-repository-governance-context-health-implementation-plan.md`
- Modify only if needed to remove current links: `docs/roadmap.md`, `docs/index.md`, `docs/development/engineering-rules.md`

**Produces:** the health PR does not leave its own temporary planning machinery in the active tree.

- [ ] **Step 1: Prove every durable health rule has been rehomed**

```text
selective bootstrap/context expansion  → AGENTS.md
current status                         → docs/roadmap.md
current-owner navigation               → docs/index.md
artifact lifecycle / Git as history    → docs/development/engineering-rules.md + existing ARCHITECTURE invariant
CI proportionality                     → engineering-rules + scripts/gate.ps1
```

If any durable health rule exists only in the temporary spec/plan, move it to the named current owner before deletion.

- [ ] **Step 2: Remove current references to the temporary health files, then delete them**

```bash
rg -n '2026-08-25-repository-governance-context-health' AGENTS.md ARCHITECTURE.md docs scripts .github || true
```

After roadmap/index no longer require them:

```bash
git rm docs/superpowers/specs/2026-08-25-repository-governance-context-health-design.md
git rm docs/superpowers/plans/2026-08-25-repository-governance-context-health-implementation-plan.md
```

- [ ] **Step 3: Verification before completion**

```bash
git diff --check main...HEAD
git diff --exit-code main...HEAD -- contracts/api/product
git diff --exit-code main -- qualification/d6-r2-wireframes
git status --short
```

```powershell
npm run gate
```

Expected: PASS, no Product/OAD diff, no current wireframe diff, no temporary health planning files, no broken current routes.

- [ ] **Step 4: Commit and stop**

```bash
git add docs/roadmap.md docs/index.md docs/development/engineering-rules.md
git add -u docs/superpowers
git commit -m "docs: retire repository health working artifacts"
```

Do **not** merge PR #71. Report final diff, CI/review state, retained/retired owner summary, and request explicit operator merge authorization.

After merge, audit remote branches separately and then reanchor PR #70 / PR #69 from the cleaned `main`.
