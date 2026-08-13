# Documentation Map

Marketplace Central intentionally keeps a **small active documentation surface**. Git history is the archive; old roadmaps/specs/handoffs are not mirrored inside the live tree.

## Read first

1. [`../AGENTS.md`](../AGENTS.md) — agent routing, process, current prohibitions.
2. [`engineering/rebaseline/README.md`](engineering/rebaseline/README.md) — **sole current program status and exact next action**.
3. [`engineering/standards/root-cause-global-maximum-method.md`](engineering/standards/root-cause-global-maximum-method.md) — canonical engineering decision method.
4. [`../ARCHITECTURE.md`](../ARCHITECTURE.md) — stable product-level constraints during D0–D9.
5. [`architecture/decisions/README.md`](architecture/decisions/README.md) — ADR registry and current/reopened status.
6. Current D-stage document(s), as linked by the rebaseline README.

## Machine/runtime authorities

These are current-state/runtime artifacts, not automatically target architecture:

- `contracts/api/marketplace-central.openapi.yaml` — current HTTP contract;
- `contracts/governance/` — current machine-enforced governance/registries;
- `contracts/gate/` — current verification ratchets;
- `scripts/gate.ps1` — shared local/CI gate implementation;
- application code/migrations — evidence of what runs today.

D0–D9 may replace current contract/topology mechanisms after explicit adjudication.

## Supporting references

Supporting references are useful but **non-authoritative** for target structure:

- `engineering/defect-class-catalog.md` — observed defect classes and proof failures from earlier work;
- `operations/live-oracle-docker.md` — current read-only live Oracle validation runner/procedure.

A supporting reference cannot override the current rebaseline, `ARCHITECTURE.md` or ADR status.

## Intentionally removed from the active tree

The documentation hygiene rebaseline removes historical surfaces including:

- `.mnfs/` mission/debt archives;
- `wiki/` architecture/framework/module pages;
- dated `docs/superpowers/` plans/specs/handoffs/evidence/runbooks;
- old repository audit/synthesis reports;
- old system blueprints/storage schemas/handoffs;
- old research/marketplace-provider planning files;
- old wave/delivery evidence;
- root legacy `IMPLEMENTATION_PLAN.md` and `EVIDENCE.md`;
- retired harness/wave doctrine files;
- obsolete production-deploy instructions while D7 is re-adjudicating runtime/deployment topology.

They remain recoverable from Git history when a current stage needs historical evidence.

## Rule for new documentation

Do not create another roadmap.

- Current progress belongs in `engineering/rebaseline/README.md`.
- Durable architecture decisions belong in an ADR or `ARCHITECTURE.md`.
- Stage-specific reasoning/evidence belongs in the active D-stage file.
- Operational instructions live under `operations/` only when they describe a current supported procedure.
- External research is gathered for the D-stage that needs it and either absorbed into the decision or clearly labeled as time-bounded reference.
- A superseded temporary design/handoff is deleted once its surviving decision/evidence is absorbed.