# Governance Contracts

This directory is the machine-readable contract boundary for repository
governance introduced by M-08/F-07. Every registry uses schema version `1.0`
and validates against a strict JSON Schema draft 2020-12 document in
`schemas/`.

The five registries own deterministic facts only:

- `modules.json`: module roots, composition requirements, API prefixes,
  declared dependencies, and exact temporary forbidden-layer edges.
- `runtime-config.json`: runtime keys, aliases, sensitivity, lifecycle, lanes,
  readers, and exact temporary direct-reader exceptions.
- `execution-lanes.json`: inheritance, target, network, database, side-effect,
  gate, and evidence contracts for each harness lane.
- `invariants.json`: checks that can be enforced mechanically and their exact
  current exceptions.
- `shared-seams.json`: narrow path sets that require exclusive ownership.

`schemas/context-pack.schema.json` defines the compiled Feature execution
context shape; it intentionally has no registry document. F-05 owns knowledge
routes, run-state/lease, checkpoints, and Portfolio/Milestone/Feature handoffs.
The preserved F-09 synthetic evaluation is not an active governance surface.

Validate the Phase 1 contracts with PowerShell 7.6 or later:

```powershell
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1
```

Temporary exceptions are baseline records, not permissions. Each exception
must identify one exact path, key, or import edge, state its removal owner, and
must not contain wildcard paths. New violations remain invalid.
