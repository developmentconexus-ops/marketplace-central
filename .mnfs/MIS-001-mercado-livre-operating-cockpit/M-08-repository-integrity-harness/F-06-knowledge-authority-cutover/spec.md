# F-06 Knowledge Authority Cutover — Specification

```yaml
id: F-06
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-08
created: 2026-07-10
updated: 2026-07-10
validation_level: QA-3
lifecycle_scope: feature
```

## Feature ID

F-06-knowledge-authority-cutover

## Problem

Marketplace Central has current execution state in `.mnfs`, but root guidance
still requires `.brain`, whose roadmap and pulse are stale and contradict
accepted M04/M05 evidence and blocked M06 evidence. Current ADRs 004-006 also
live under `.brain`, preventing deletion without an authority migration.

## Requirements

- Rehome only current ADRs 004, 005, and 006 under
  `docs/architecture/decisions/`, preserving their decision semantics and IDs.
- Remove superseded ADRs 001-003, stale roadmap/pulse/session state, and the
  completed legacy plan from the current tree; Git history is sufficient.
- Update active startup, truth-order, planning, and architecture references so
  `.brain` is neither read nor written.
- Reconcile architecture status from accepted M04/M05 and blocked M06 MNFS
  validation, without promoting M06 to passed.
- Keep historical evidence unchanged when an old `.brain` path is part of the
  truthful record of what happened.
- Delete every tracked `.brain` path in the same intentional cutover commit.

## Acceptance Evidence

- Tracked-file scan reports zero `.brain` paths.
- Active-reference scan reports zero `.brain` or `Nexus Brain` matches.
- New ADR index contains exactly 004-006 and the three files preserve context,
  decision, rationale, consequences, and rejected alternatives.
- Architecture labels product links and inventory as active validated
  foundations, and orders/profitability as implemented but milestone-blocked.
- Git history still resolves deleted `.brain` content.

## Non-Goals

- Do not create governance registries or schemas; F-07 owns them.
- Do not rewrite historical plans, completed MNFS evidence, or old baseline
  inventories merely to remove old path strings.
- Do not alter product code, OpenAPI, SDK, migrations, dependencies, or runtime.
- Do not repair M06 blockers or restart M09+.

## Design

Perform one atomic documentation/authority migration. The new ADR directory
becomes architecture-decision truth when `AGENTS.md` changes. `.mnfs` remains
the sole execution-status owner. Wiki remains human operating knowledge and
links to `.mnfs` for sequencing. Architecture keeps durable topology and
accurately summarizes current module maturity from MNFS verdicts.

No tombstone directory, redirect file, symlink, or compatibility reader is
created for `.brain`.

## Edge Cases

- Historical documents mentioning `.brain` remain unchanged only when dated or
  completed evidence makes the reference historically accurate.
- ADR numbering remains 004-006; deletion of retired IDs does not renumber
  accepted decisions.
- M06 has real implementation and partial live evidence, but status remains
  blocked until its explicit paid resolved-link scenario and cold gate pass.
- `contracts/governance/README.md` remains reserved until F-07 establishes its
  schemas and authority; F-06 must not prematurely call it runtime truth.

## Acceptance Criteria

### F06-AC01 — Current ADR authority migrated

- Traces to milestone criterion ID: `M-08-C08`.
- Proven by: ADR file/index inspection and `rg -n 'ADR-00[4-6]'` across current
  architecture sources.

### F06-AC02 — Active guidance has one execution truth

- Traces to milestone criterion ID: `M-08-C08`.
- Proven by: scoped active-reference scan over `AGENTS.md`, architecture, wiki,
  current handoff, and the M08 execution guide.

### F06-AC03 — Stale brain tree retired without knowledge loss

- Traces to milestone criterion ID: `M-08-C08`.
- Proven by: `git ls-files -- .brain` returning empty and
  `git log --all -- .brain` returning prior history.

### F06-AC04 — Module maturity is evidence-honest

- Traces to milestone criterion ID: `M-08-C08`.
- Proven by: direct comparison of `ARCHITECTURE.md` wording with M04, M05, and
  M06 `validation-result.md` verdicts.

## Handoff

- Current status: `spec_ready`.
- Next owner: Feature Implementer.
- Next action: Execute the split plan in a fresh build session.
- Required files/evidence: feature brief, this spec, plan, M04-M06 validation
  results, ADR sources, and active-reference inventory.
- Blockers or open decisions: None; `.brain` retirement is operator-approved.
