# F-10 — Pragmatic Harness Cutover Specification

## Problem

The active harness still implements the superseded cold-clone, isolated-cache,
provisioning, and snapshot workflow. It cannot be retained as an alias or
compatibility path because it turns ordinary feature validation into a
clean-machine simulation.

## Requirements

1. Retire every active cold command, alias, lane, snapshot/provisioning helper,
   schema field, and cold-only test while preserving historical MNFS evidence.
2. Add an `impact` harness command that accepts a context pack, validates its
   base and source hashes, discovers changed paths in the normal checkout,
   reports impacted declared seams, and blocks before tests for stale or
   out-of-scope inputs.
3. Resolve each selected command ID through a versioned, closed registry to a
   typed executable and argv. The pack and plan never supply executable shell
   text.
4. Persist a schema-valid, redacted outcome with base SHA, commit SHA, changed
   paths, risk, ordered command records, target labels, exits, artifacts,
   remaining risks, and aggregate classification.

## Non-goals

- Product, OpenAPI/SDK, migration, frontend, provider, Docker, network,
  browser, Oracle, or provider-write work.
- Replacing existing environment, process-execution, policy, context, or
  PostgreSQL behavior beyond the cold cutover.
- Rewriting historical F-04 evidence.

## Design

`scripts/harness/Impact.psm1` owns a closed command registry and impact-gate
orchestration. The façade delegates `-Command impact` to it. Context contracts
carry a task command ID plus derived target metadata, while the registry owns
the executable/argv shape. Evidence owns redaction and outcome serialization;
its schema becomes current-checkout oriented rather than cold-candidate
oriented.

## Edge cases

- A pack with a stale base/source hash, an undeclared command, unknown target,
  or changed path outside `paths.allowed` fails closed with a stable reason.
- A failing registered command records its result and makes the aggregate fail;
  later commands do not execute.
- Absolute paths and secret-like output cannot persist in outcome evidence.

## Acceptance criteria

### F10-AC01

M-08-C17: Boundary test/search proves no active cold command, alias, lane,
snapshot/provisioning code, or cold-only tests remain.

### F10-AC02

M-08-C15: Fixture pack runs exactly two registered IDs in order; an undeclared
ID is rejected without execution.

### F10-AC03

M-08-C15: Stale pack, out-of-scope changed path, and unknown target fixtures
block with pinned reason codes.

### F10-AC04

M-08-C15: Outcome fixture validates redacted base/commit SHAs, changed paths,
risk, records, and aggregate semantics.
