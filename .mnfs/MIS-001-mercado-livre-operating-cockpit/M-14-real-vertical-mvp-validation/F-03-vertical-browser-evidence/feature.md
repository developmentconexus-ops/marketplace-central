# F-03-vertical-browser-evidence

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-14
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-3
lifecycle_scope: feature
```

## Mission

MIS-001 MVP operator journey.

## Milestone

M-14 Real Vertical MVP Validation.

## Brief

Make the selected real scenario addressable through accepted IC-003 routes and define
exact browser observables/evidence targets without issuing the final QA verdict.

## Inputs

- F-01 real scenario support, F-02 accepted validation lanes, and IC-004 MVP-BD-01.
- Passed M-13 routes/components and M-14-C02 drive.

## Inputs/Outputs

- Input: validated local selection, frozen SHA, addressable IC-003 routes, and the
  accepted F-01/F-02 `validation.md` evidence.
- Output: stable test IDs/addressability only where proved necessary; QA later
  writes the exact IC-004 browser interaction, screenshots, and method/path files.

## Interaction Model

- Start `/`, open stock attention to Listing, navigate Listing → Product → Listing
  → Sale → Listing, review simulation, reload, and inspect identity, installation,
  source, and quality at each step.
- Browser capture records URL, visible assertion, timestamp, screenshot path, and
  network method/path summary; it excludes response bodies with PII.

## State Model

The successful current path is accompanied by at least one unknown/stale/conflict
negative sample. Browser state follows IC-003 and no test-only state name is added.

## Negative Scenarios

- Real sample not addressable: persist IC-004 terminal `externally_blocked` with
  `missing_edge=route_addressability`; do not invent an API error/UI state.
- Screenshot contains buyer PII: discard/redact and repeat before acceptance.
- Network observes provider mutation: stop and mark blocking failure.
- Deep-link reload loses context: feature remains incomplete.

## Expected Output

QA can drive the real scenario with exact visible assertions and safe evidence paths;
no transcript reconstruction or manual identity search is needed.

## Constraints

- Do not manufacture UI fixtures for the successful real lane.
- Do not change domain calculations to satisfy browser assertions.
- Owned paths: minimal browser-addressability/test IDs/evidence capture support and this root.
- Forbidden paths: provider writes, auth/RBAC, final `validation-result.md` verdict.

## Criteria IDs

- M-14-C02 Vertical browser journey.
- M-14-C05 Evidence security and honesty.
- M-14-C06 No provider mutation.

## Validation Expectations

- Dry browser drive reaches every M-14-C02 URL and visible assertion.
- Reload retains identity/context.
- Method/path summary proves no provider mutation.
- Evidence redaction scan finds zero secret/PII patterns.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer after F-01/F-02.
- Next action: Create spec/plan and implement only browser-evidence readiness gaps.
- Required files/evidence: F-01/F-02 `validation.md`, IC-003, IC-004, and this
  feature's `validation.md`.
- Blockers or open decisions: Missing real data is reported only through the IC-004
  terminal checkpoint; fixtures cannot fill the successful lane.
