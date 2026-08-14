# Governance Contracts

This directory is the machine-readable boundary for **current repository governance** during the architecture rebaseline.

These contracts describe facts that the current gates can enforce mechanically. They are **not target-system architecture** and do not decide whether the present module, runtime, API or integration topology survives D0–D9.

The active registries are:

- `modules.json` — current module/context roots, composition requirements, API prefixes, declared dependencies, and exact temporary boundary exceptions;
- `runtime-config.json` — current runtime keys, aliases, sensitivity, lifecycle, lanes, readers, and exact temporary direct-reader exceptions;
- `execution-lanes.json` — current harness lane inheritance, targets, network/database/side-effect behavior, gates and evidence classes;
- `invariants.json` — checks enforced mechanically and their exact current exceptions;
- `shared-seams.json` — narrow path sets that require exclusive ownership while they exist;
- `knowledge-routes.json` — bounded routes into current authority and supporting evidence.

`schemas/context-pack.schema.json` defines a generic bounded execution-context shape for tooling. It deliberately does not encode historical Milestone/Feature ownership.

`harness-evals.json` is a supporting synthetic evaluation corpus for harness invariants. It is not a roadmap, status document or target-design authority.

Validate the governance contracts with PowerShell 7.6 or later:

```powershell
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1
```

Temporary exceptions are baseline records, not permissions and not promises of a specific future refactor. Each exception identifies the exact current path/key/edge and the D-stage in which that fact must be re-adjudicated if it remains relevant. New violations remain invalid.

Current program status and the exact next action live only in `docs/engineering/rebaseline/README.md`.
