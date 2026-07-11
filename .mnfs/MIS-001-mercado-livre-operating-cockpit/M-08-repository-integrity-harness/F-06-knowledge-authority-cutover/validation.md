# F-06 Knowledge Authority Cutover — Validation

```yaml
id: F-06
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
evidence_type: ran
```

## Result

- Result: `quick_validation_passed`
- Blocking failures: none
- Scope: documentation and repository authority cutover only
- External integration/runtime validation: not applicable; no provider, database, queue, credential, or environment-dependent behavior changed

## Changed Paths

- Modified: `AGENTS.md`, `ARCHITECTURE.md`, `wiki/README.md`, `wiki/vision/system-vision.md`, and `docs/superpowers/handoffs/2026-07-09-m06-orders-margin-session-handoff.md`.
- Created: `docs/architecture/decisions/README.md` and ADR files 004-006.
- Deleted: all 11 tracked `.brain` files.
- Created/updated: this validation artifact and the F-06 feature handoff.

## Acceptance Evidence

### F06-AC01 — Current ADR authority migrated

- Evidence type: `ran`
- Command: `rg -n 'ADR-00[4-6]' ARCHITECTURE.md docs/architecture/decisions wiki .mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md`
- Expected: current IDs remain 004-006 and appear in the new decision owner/index or explicit current references.
- Actual: exit `0`; ADR-004, ADR-005, and ADR-006 each appeared in its new file and exactly once in the new index. Historical mission/wiki references to ADR-005 remained truthful.
- Manual inspection: each migrated ADR preserves its original ID, date, accepted status, context, decision, rationale, consequences, and rejected alternatives.

### F06-AC02 — Active guidance has one execution truth

- Evidence type: `ran`
- Command: `rg -n --hidden --glob '!.git/**' --glob '!**/.env*' --glob '!**/secrets/**' '(\.brain[/\\]|Nexus Brain)' AGENTS.md ARCHITECTURE.md wiki docs/superpowers/handoffs/2026-07-09-m06-orders-margin-session-handoff.md .mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/execution-guide.md`
- Expected: exit `1`, no matches.
- Actual: exit `1`, no matches; the verification wrapper emitted `NO_ACTIVE_MATCHES` and exited `0`.
- Ownership: architecture truth is `ARCHITECTURE.md` plus `docs/architecture/decisions/`; human operating knowledge is the wiki; execution truth is `.mnfs/`; `contracts/governance/` remains reserved for F-07.

### F06-AC03 — Stale brain tree retired without knowledge loss

- Evidence type: `ran`
- Command: `git ls-files -- .brain`
- Expected: no output after the approved deletions are staged.
- Actual: exit `0`, no output; zero tracked `.brain` paths remain in the cutover index.
- Command: `git log --all --oneline -- .brain`
- Expected: historical commits remain resolvable.
- Actual: exit `0`; history returned commits from `a4fd306` through `02dd3f1`, including prior architecture and roadmap updates.
- Historical-reference classification: dated/completed evidence outside the plan's active-source list remains truthful historical evidence and was intentionally not rewritten. The scoped active-source scan has zero current references.

### F06-AC04 — Module maturity is evidence-honest

- Evidence type: `ran`
- Command: `rg -n '(product_links|inventory|orders|profitability|listing|stock|order)' ARCHITECTURE.md`
- Expected: M04/M05 foundations are active/validated; M06 modules are implemented but milestone-blocked; no live provider mutation is claimed.
- Actual: exit `0`; `product_links` and `inventory` are `active validated foundation`; `orders` and `profitability` are `implemented — M-06 milestone blocked`; connector wording explicitly says no live provider stock mutation is claimed.
- Command: ``rg -n 'status: passed|Result: `passed`' .mnfs/MIS-001-mercado-livre-operating-cockpit/M-04-product-links-ml/validation-result.md .mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/validation-result.md``
- Actual: exit `0`; both M04 and M05 report `status: passed` and `Result: passed`.
- Command: `rg -n 'status: blocked|BLOCKED' .mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-result.md`
- Actual: exit `0`; M06 reports `status: blocked` and retains the explicit `BLOCKED — do not mark M-06 passed` verdict.

## Integrity Evidence

- `git diff --check`: exit `0`; no whitespace errors. Git emitted only line-ending conversion warnings for existing Windows working-copy policy.
- `git status --short`: only allowed F-06 paths changed.
- `git diff --stat`: before staging the new files, reported the five active-source edits and 11 `.brain` deletions; staged review includes the four new ADR paths and this feature's two execution artifacts.
- No `.env`, secrets, credentials, provider payloads, or buyer PII were read.

## Risks And Boundaries

- This feature changes knowledge authority, not runtime behavior.
- Deleted content remains available through Git history; no tombstone, redirect, symlink, compatibility reader, or legacy roadmap copy was added.
- M06 remains blocked on its required paid resolved-link realization evidence and cold gate.
- F-07 still owns machine-governance schemas and registries.

## Handoff

- Current status: `quick_validation_passed`.
- Next owner: Milestone Orchestrator / independent QA-3 review.
- Next action: review the fixed cutover commit against F06-AC01 through F06-AC04 and accept or return one scoped correction batch.
- Blockers: none.
