# F-06 Knowledge Authority Cutover — Plan

```yaml
id: F-06
type: feature-plan
status: planned
owner: Feature Implementer
parent: M-08
created: 2026-07-10
updated: 2026-07-10
validation_level: QA-3
lifecycle_scope: feature
split_decision: split
split_reason: more than four owned paths and deletion of an authority tree require a fresh build session
```

## Feature ID

F-06-knowledge-authority-cutover

## Allowed Paths

- `AGENTS.md`
- `ARCHITECTURE.md`
- `docs/architecture/decisions/README.md`
- `docs/architecture/decisions/004-integration-catalog-plugin-framework.md`
- `docs/architecture/decisions/005-mercado-livre-first-control-plane.md`
- `docs/architecture/decisions/006-oracle-internal-read-owned-by-mpc.md`
- `wiki/README.md`
- `wiki/vision/system-vision.md`
- `docs/superpowers/handoffs/2026-07-09-m06-orders-margin-session-handoff.md`
- `.brain/**` for deletion only
- `F-06-knowledge-authority-cutover/validation.md`
- `F-06-knowledge-authority-cutover/feature.md` for final status/handoff only

All other paths are forbidden. Historical `.brain` references outside the
listed active sources are evidence and must not be rewritten.

## Steps

1. Capture pre-cutover tracked files and scoped active references. Record only
   file paths/lines, never environment values.
2. Create `docs/architecture/decisions/README.md` indexing exactly ADR-004,
   ADR-005, and ADR-006. Recreate the three decision files with the same IDs,
   status, context, decision, rationale, consequences, and alternatives from
   their `.brain/decisions/` sources.
3. Update `AGENTS.md`: remove `.brain` from startup, truth order, planning, and
   architecture-change requirements; point ADR truth to
   `docs/architecture/decisions/`; keep `.mnfs` as execution truth.
4. Update `ARCHITECTURE.md`: link the new ADR directory; mark `product_links`
   and `inventory` active/validated; mark `orders` and `profitability`
   implemented with M06 milestone blocked; update connector wording without
   claiming a live provider write.
5. Update wiki planning statements and the current M06 handoff startup list to
   use `.mnfs` and current validation artifacts. Mark the dated M06 handoff as
   superseded for portfolio sequencing while preserving its historical RED and
   live-evidence record.
6. Delete the complete tracked `.brain` directory. Do not add a tombstone,
   redirect, or copied legacy roadmap.
7. Run all verification commands. If any active reference or current decision
   is missing, restore only the missing current fact through `apply_patch`; do
   not restore `.brain`.
8. Write `validation.md` with exact commands, expected/actual outputs, changed
   paths, evidence type `ran`, historical-reference classification, risks, and
   handoff. Update feature status to `quick_validation_passed` only when every
   command passes.
9. Stage only allowed paths and create one intentional commit:
   `refactor(governance): retire brain knowledge store`.

## Exact Content Contracts

`docs/architecture/decisions/README.md` must use this index shape:

```markdown
# Architecture Decision Records

Current accepted decisions:

- [ADR-004: Integration Catalog Plugin Framework](004-integration-catalog-plugin-framework.md)
- [ADR-005: Mercado Livre First Control Plane](005-mercado-livre-first-control-plane.md)
- [ADR-006: MPC-Owned Oracle Internal Reads](006-oracle-internal-read-owned-by-mpc.md)

Superseded and delivery-only decisions remain available through Git history.
```

`AGENTS.md` startup must resolve to this order:

```markdown
1. `ARCHITECTURE.md`
2. `wiki/README.md`
3. `.mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md`
4. The active milestone `execution-guide.md` named by the current goal, when applicable
```

Planning ownership statements must say `.mnfs/` is execution truth and
`IMPLEMENTATION_PLAN.md` is historical reconciliation only.

## Files Expected To Change

- Create the four `docs/architecture/decisions/*` paths listed above.
- Modify `AGENTS.md`, `ARCHITECTURE.md`, `wiki/README.md`,
  `wiki/vision/system-vision.md`, and the dated M06 handoff.
- Delete all 11 tracked `.brain` files.
- Create feature `validation.md` and update feature handoff/status.

## Verification Commands

### F06-AC01

```powershell
rg -n 'ADR-00[4-6]' ARCHITECTURE.md docs/architecture/decisions wiki .mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md
```

Expected: each current ID appears in the new ADR index/file set or an explicit
current architecture reference; no ID is renumbered.

### F06-AC02

```powershell
rg -n --hidden --glob '!.git/**' --glob '!**/.env*' --glob '!**/secrets/**' '(\.brain[/\\]|Nexus Brain)' AGENTS.md ARCHITECTURE.md wiki docs/superpowers/handoffs/2026-07-09-m06-orders-margin-session-handoff.md .mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/execution-guide.md
```

Expected: exit 1 with no matches.

### F06-AC03

```powershell
git ls-files -- .brain
git log --all --oneline -- .brain
```

Expected: first command emits nothing; second command emits historical commits.

### F06-AC04

```powershell
rg -n '(product_links|inventory|orders|profitability|listing|stock|order)' ARCHITECTURE.md
rg -n 'status: passed|Result: `passed`' .mnfs/MIS-001-mercado-livre-operating-cockpit/M-04-product-links-ml/validation-result.md .mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/validation-result.md
rg -n 'status: blocked|BLOCKED' .mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-result.md
```

Expected: architecture wording matches passed M04/M05 and blocked M06 without
claiming live provider mutation beyond recorded evidence.

### Integrity

```powershell
git diff --check
git status --short
git diff --stat
```

Expected: no whitespace errors; only allowed paths changed before staging.

## QA Steps

1. Read each new ADR beside its source and confirm no decision semantic changed.
2. Classify any remaining repository-wide `.brain` match as dated historical
   evidence. Any active instruction or current runbook match blocks completion.
3. Confirm `git show HEAD^:.brain/system-pulse.md` or another historical commit
   can still retrieve deleted material after commit.

## Rollback/Risk Notes

- Do not use `git reset`, `git checkout`, `git clean`, or `git restore`.
- If a current fact was omitted, add it to its new canonical owner with
  `apply_patch` before commit.
- If unrelated changes appear, stop and report ownership conflict; do not stage
  or delete them.
- Deleted material remains recoverable through Git history, but current code and
  guidance must never depend on that recovery path.

## Handoff

- Current status: `planned`.
- Next owner: Fresh Feature Implementer in `build` mode.
- Next action: Execute steps 1-9 without re-planning or broad discovery.
- Required files/evidence: spec, plan, ADR sources, M04-M06 validation results,
  scoped reference scans, diff, validation artifact, and commit SHA.
- Blockers or open decisions: None.
