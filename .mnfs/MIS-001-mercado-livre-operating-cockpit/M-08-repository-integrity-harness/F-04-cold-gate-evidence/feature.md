# F-04-cold-gate-evidence

```yaml
id: F-04
type: feature-brief
status: superseded
owner: Mission Strategist
parent: M-08
created: 2026-07-10
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-08 Repository Integrity Harness.

## Brief

Historical experiment that attempted to prove a candidate through a detached
local clone, empty caches, dependency provisioning, Docker image preparation,
and two identical runs.

## Supersession Decision

The experiment is not accepted and is no longer a V1 requirement. It repeatedly
blocked before product validation on host-specific Git/sandbox behavior while
the accepted unit environment and ephemeral PostgreSQL lanes already protected
the meaningful safety boundaries. The operator explicitly rejected local
clean-machine simulation as contrary to the harness objective.

Its spec, plan, validation, commits, and ignored run artifacts remain historical
evidence. They must not be rewritten to imply success. F-10 removes active
cold-only code and preserves reusable evidence/redaction primitives.

## Constraints

- Do not resume, repair, retry, or accept the cold pair.
- Do not delete historical evidence or hide the blocked result.
- No future criterion may depend on cold clone, clean caches, or local dependency reprovisioning.

## Handoff

- Current status: Superseded, not passed.
- Next owner: F-10 Feature Implementer.
- Next action: Retire active cold-only surfaces and introduce the impact gate.
- Required files/evidence: Existing F-04 artifacts plus F-10 cutover validation.
- Blockers or open decisions: None; cold diagnosis is closed.

