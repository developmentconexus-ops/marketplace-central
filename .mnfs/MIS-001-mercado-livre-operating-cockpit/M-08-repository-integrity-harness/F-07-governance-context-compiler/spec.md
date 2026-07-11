# F-07 Governance Registry and Context Compiler — Specification

```yaml
id: F-07
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Feature ID

F-07-governance-context-compiler

## Problem

Workers currently receive manually written context. Module boundaries, runtime
keys, execution lanes, invariants, and shared seams are scattered across prose
and code. F-02 demonstrated the failure mode: a broad requirement such as
“block OAuth and database targets” omitted active keys because no executable
inventory constrained implementation.

F-07 must establish one machine-readable governance owner and compile bounded
worker context from current MNFS artifacts and Git state. It must expose known
legacy violations without blessing them and must reject any undeclared growth.

## Requirements

### Governance contracts

- Use strict JSON Schema draft 2020-12 with `additionalProperties: false`.
- Validate through PowerShell 7 `Test-Json -SchemaFile`; no new package or
  module dependency.
- Own only facts with deterministic verification in F-07:
  module roots/dependencies/composition/API prefixes, runtime keys/readers,
  execution-lane permissions, code invariants, and exclusive shared seams.
- Represent current violations as exact temporary exceptions with stable ID,
  path/key or import edge, reason, and removal owner. Exceptions cannot use
  wildcards and cannot authorize new violations.
- Keep PostgreSQL driver rules scoped to `adapters/postgres`; Oracle adapters
  may legitimately use their Oracle driver boundary.
- Do not claim semantic enforcement for tenant SQL, unknown-to-zero, provider
  write policy, or evidence honesty until targeted tests/schemas own those
  behaviors in later features.

### Context packs

- Compile from a pinned MNFS feature path and full 40-character base SHA.
- Derive objective, observable done state, criterion IDs, and proof commands
  from `feature.md`, `spec.md`, a schema-valid machine work-contract block in
  `plan.md`, and the parent validation contract. Do not infer command mappings
  from free-form prose.
- Record ordered source paths with SHA-256 hashes; never copy entire source
  documents into the pack.
- Record deterministic risk classification, review policy, advisory model,
  allowed/forbidden paths, shared seams, side effects, commands, target labels,
  stop conditions, retry budget, and compact handoff fields.
- Estimate input size as `ceil(serialized UTF-8 character count / 4)`, excluding
  the estimate field itself. Name it an estimate and reject values above 2,000.
- Use repo-relative normalized `/` paths. Reject rooted paths and `..`.
- A currentness check must reject base-SHA mismatch, missing source, source-hash
  mismatch, path overlap/out-of-scope, undeclared seam, dangling criterion
  proof, side-effect conflict, and target/evidence inflation.
- Validation must recompile expected objective, done state, sources, criteria,
  commands, risk, seams, side effects, and estimate from canonical files. Pack
  fields are output, never trusted validation inputs.

### Public surface

- Keep `scripts/harness.ps1` as stable dispatcher.
- Add commands `governance-validate`, `governance-drift`, `context-compile`, and
  `context-validate`; existing F-02 commands remain behaviorally unchanged.
- Root commands emit only `status`, stable `error_code`, safe IDs/paths, and
  artifact path. No environment value or secret enters output.

## Acceptance Evidence

- Positive registries and context fixture validate with exit 0.
- Schema, reference, semantic, currentness, criterion, path/seam, side-effect,
  target-label, and size negative fixtures exit 1 with pinned reason codes.
- Runtime inventory includes every active production/dev reader and exact
  temporary direct-reader exceptions; a new reader fixture fails.
- Module inventory covers exactly current module directories; a new module,
  import edge, forbidden layer, duplicate migration prefix, production panic,
  OpenAPI/SDK split, or frontend direct fetch fixture fails unless an exact
  non-expired exception exists.
- Context compilation for F-07 produces a pack at or below the declared
  estimate ceiling and validates against the current base and source hashes.

## Non-Goals

- Do not refactor environment readers, cross-module imports, migrations,
  production panic sites, or F-02 child execution; later milestones/features
  remove those declared exceptions.
- Do not implement active lease storage/conflict resolution; F-05 owns it.
- Do not implement evidence, run-state, or eval-result schemas; F-04/F-05/F-09
  own them.
- Do not add hooks, agents, skill, worktree, PostgreSQL lifecycle, or cold gate.
- Do not change product code, OpenAPI, SDK, migrations, or provider behavior.

## Design

### Registry ownership

`contracts/governance/` becomes authoritative for its accepted machine facts.
`AGENTS.md` and `ARCHITECTURE.md` retain rationale and link to registry IDs.
Every registry has `schema_version: "1.0"`, strict schema, unique IDs, and a
semantic validator.

`feature-work-contract.schema.json` owns the structured block embedded under
`## Machine Work Contract` in each executable feature plan. The block declares
required extra sources, allowed/forbidden paths, side effects, lane-bound
command templates, criterion-to-command mappings, stop conditions, retry budget, and
handoff fields. Objective/done, source hashes, target/evidence labels, seams,
risk, review, model advice, and token estimate are derived by the compiler.
Only `{base_sha}`, `{context_path}`, and `{feature_path}` template variables are
supported; compilation expands all variables so emitted commands are exact.

`modules.json` records current roots, code owner path, composition requirement,
OpenAPI prefixes, and declared module dependencies. Target layer rules live in
invariants; exact current forbidden-layer edges live in temporary exceptions.

`runtime-config.json` records canonical keys, aliases, owner, sensitivity,
allowed lanes, exact readers, lifecycle, and optional removal owner. Edge
readers such as Vite/Docker/harness are distinct from typed Go config owners.
Current direct readers are exact temporary exceptions; new copies fail.

`execution-lanes.json` records environment inheritance, target label, network,
database, side effects, gates, and evidence class. It declares the target
contract (`inherit_parent: false`) even though F-08 implements process
isolation.

`invariants.json` contains only deterministic checks. Current violations use
exact temporary exceptions; the target rule remains visible and blocking for
new code.

`shared-seams.json` contains narrow exclusive path sets for API/SDK,
migrations, composition root, dependency graph, architecture decisions, and
provider capability contract.

### PowerShell modules

`Policy.psm1` exports:

```powershell
Import-GovernanceDocument -Path <path> -SchemaPath <path>
Test-GovernanceContracts -RepositoryRoot <path>
Test-GovernanceDrift -RepositoryRoot <path> -BaseSha <sha>
Get-ApplicableInvariants -Registry <object> -Paths <string[]>
Get-SharedSeams -Registry <object> -Paths <string[]>
```

`Context.psm1` exports:

```powershell
New-HarnessContextPack -FeaturePath <path> -BaseSha <sha> -AllowedPath <string[]> -OutputPath <path>
Test-HarnessContextPack -Path <path> -RepositoryRoot <path> [-RequireCurrentBase]
```

Functions return typed result objects. Only dispatcher formats stdout and exit
code.

### Stable reason codes

Governance uses `GOV_DOCUMENT_MISSING`, `GOV_SCHEMA_INVALID`,
`GOV_REFERENCE_INVALID`, `GOV_SEMANTIC_DRIFT`, plus precise `GOV_*` invariant
codes. Runtime inventory uses `RCFG_UNDECLARED_READ`,
`RCFG_UNAPPROVED_READER`, `RCFG_DIRECT_READ_FORBIDDEN`,
`RCFG_DYNAMIC_READER_UNBOUNDED`, `RCFG_ALIAS_UNDECLARED`,
`RCFG_ALIAS_COLLISION`, `RCFG_LANE_VIOLATION`,
`RCFG_SECRET_CLASS_MISMATCH`, and `RCFG_READER_MISSING`.

Context uses `CTX_SCHEMA_INVALID`, `CTX_FEATURE_INVALID`,
`CTX_BASE_SHA_INVALID`, `CTX_BASE_SHA_MISMATCH`, `CTX_SOURCE_MISSING`,
`CTX_SOURCE_HASH_MISMATCH`, `CTX_PATH_SCOPE_CONFLICT`,
`CTX_PATH_OUTSIDE_SCOPE`, `CTX_SHARED_SEAM_UNDECLARED`,
`CTX_CRITERION_PROOF_MISSING`, `CTX_PROOF_REFERENCE_INVALID`,
`CTX_SIDE_EFFECT_CONFLICT`, `CTX_TARGET_EVIDENCE_INFLATION`, and
`CTX_TOKEN_BUDGET_EXCEEDED`.

## Edge Cases

- `API_PORT` has precedence over `SERVER_ADDR`; they are separate canonical
  keys, not aliases.
- Oracle canonical keys may resolve declared `SANKHYA_ORACLE_*` legacy aliases;
  one alias cannot map to two canonical owners.
- `MS_DATABASE_URL` is a legacy-current key scheduled for M-09 removal, not an
  alias of `MC_DATABASE_URL`.
- `MPC_TEST_DATABASE_URL` is a reserved explicit harness input, not an approved
  ambient reader.
- Tool/host keys (`PATH`, temp directories, `GOCACHE`) belong to F-08 child
  environment policy, not application runtime configuration.
- Vite built-in `import.meta.env.DEV` is not a repository-owned key.
- Context validation with `-RequireCurrentBase` is a pre-dispatch check; a
  worker commit intentionally changes HEAD afterward.
- A requested allowed path must be equal to or narrower than a declared work
  contract path. Mentioning one descendant never authorizes its ancestor.
- Token estimate uses UTF-8 byte count, not .NET UTF-16 string length.
- Dispatcher output converts absolute internal paths to safe repo-relative
  paths or omits them.

## Acceptance Criteria

### F07-AC01 — Registry schemas and references are deterministic

- Traces to milestone criterion ID: `M-08-C09`.
- Proven by: governance contract tests and `governance-validate` command.

### F07-AC02 — Current code drift is explicit and cannot grow

- Traces to milestone criterion ID: `M-08-C09`.
- Proven by: semantic drift fixtures plus current-repository
  `governance-drift` command.

### F07-AC03 — Context is criterion-complete and source-current

- Traces to milestone criterion ID: `M-08-C09`.
- Proven by: compile/validate positive fixture and base/source/proof negative
  fixtures.

### F07-AC04 — Paths, seams, side effects, and evidence labels fail closed

- Traces to milestone criterion ID: `M-08-C09`.
- Proven by: path overlap/out-of-scope, undeclared seam, side-effect conflict,
  and target/evidence inflation fixtures.

### F07-AC05 — Context remains bounded

- Traces to milestone criterion ID: `M-08-C09`.
- Proven by: deterministic estimate positive case and 2,001-token rejection.

## Handoff

- Current status: `spec_ready`.
- Next owner: Feature Implementer.
- Next action: Execute the split plan serially with one writer per phase.
- Required files/evidence: registries/schemas, PowerShell modules, fixtures,
  RED/GREEN output, current drift output, context pack, and validation artifact.
- Blockers or open decisions: None; known violations must be exceptions, not
  silently allowed dependencies/readers.
