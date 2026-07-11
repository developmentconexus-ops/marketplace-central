# F-08 Hermetic Child Runtime — Specification

```yaml
id: F-08
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Problem

F-02 rejects a hand-maintained set of ambient variables, then executes Go and
npm inside the inherited process environment. Unknown variables still leak,
configured values can escape redaction, and unit execution temporarily mutates
the parent CWD and `GOCACHE`. Adding more denied names cannot make this hermetic.

F-08 must construct every child environment from an explicit allowlist and run
subprocesses through a typed, shell-free boundary. It preserves the public
harness surface while leaving database lifecycle, rich evidence, leases, and
provider behavior to their owning features/modules.

## Requirements

### Environment construction

- `Environment.psm1` constructs a fresh case-insensitive dictionary; it never
  copies the parent environment and then removes keys.
- A lane must exist in `execution-lanes.json` and declare
  `inherit_parent: false`. Runtime values may enter only when the key is allowed
  by that lane and consistently declared in `runtime-config.json`.
- Unit inherits only the fixed host/tool keys `SystemRoot`, `WINDIR`, `ComSpec`,
  `PATH`, `PATHEXT`, `TEMP`, and `TMP` when present. It generates canonical
  repository-local `GOCACHE` and `GOMODCACHE`, plus `GOPROXY=off` and
  `GOSUMDB=off`; missing modules fail deterministically instead of reaching the
  network. It receives no application runtime key.
- Do not inherit `HOME`, `USERPROFILE`, `APPDATA`, `LOCALAPPDATA`,
  `PSModulePath`, tool caches, proxy/certificate/cloud/CI variables, or unknown
  future variables. They are hidden inputs, not required unit configuration.
- Unit ignores `-EnvFile` entirely. Live preflight retains declared Oracle
  legacy aliases but normalizes them to canonical keys before use.
- Environment construction never mutates the parent environment or location.

### Process execution

- Resolve executable paths before child construction. Use
  `System.Diagnostics.ProcessStartInfo` with `UseShellExecute=false`, an
  absolute executable, explicit canonical working directory,
  `Environment.Clear()`, and one `ArgumentList.Add()` call per argument.
- Capture stdout and stderr separately. Begin both asynchronous drains before
  waiting so full pipes cannot deadlock. Preserve the actual child exit code;
  timeout kills the process tree and returns a stable typed reason.
- Redact before output enters a returned object, exception, console, or run
  artifact. Replace exact configured candidate values longest-first and apply
  assignment/URI defenses. Values classified `secret` or `internal` are
  candidates, except short public gate literals such as `1`.
- Parent CWD and environment remain byte-for-byte unchanged after success,
  failure, start error, or timeout.

### Modular surface

- `scripts/harness.ps1` remains the stable parameter-compatible dispatcher.
- `Environment.psm1` owns env-file parsing, canonical alias resolution, lane
  validation, and child dictionary construction.
- `Execution.psm1` owns typed process requests/results, safe invocation,
  capture, timeout, and redaction.
- Existing `Policy.psm1` and `Context.psm1` public behavior remains unchanged.
- Do not create placeholder `Evidence.psm1` or `State.psm1`; F-04 and F-05 own
  those concepts when behavior exists.
- Business rules, provider payloads, database lifecycle, migrations, worktree
  leases, and browser automation remain outside F-08.

### Compatibility

- Preserve direct commands `unit`, `integration`, `live`, `browser`,
  `provider-write`, `governance-validate`, `governance-drift`, `governance`,
  `context-compile`, and `context-validate` with their parameters.
- Exercise all eight root npm aliases as behavior. Expected blocked lanes must
  retain their pinned pre-contact/pre-network reason.
- `unit -PreflightOnly` reports `target=fake`; normal unit runs Go from
  `apps/server_core` and web tests from the repository root.
- Integration remains blocked before contact until F-03 supplies the
  ephemeral PostgreSQL lifecycle. Browser non-preflight and provider write
  remain fail-closed. Live remains preflight-only.
- Absolute dispatcher invocation and `npm --prefix <repo>` work from an
  unrelated directory; repository-relative public paths resolve from
  `$PSScriptRoot`, never caller CWD.

## Stable Reason Codes

- Environment: `HENV_REPOSITORY_ROOT_INVALID`, `HENV_LANE_UNSUPPORTED`,
  `HENV_PARENT_INHERITANCE_FORBIDDEN`, `HENV_RUNTIME_KEY_FORBIDDEN`,
  `HENV_RUNTIME_CONTRACT_DRIFT`, `HENV_TOOL_KEY_MISSING`.
- Execution: `HEXEC_FILE_NOT_FOUND`, `HEXEC_WORKING_DIRECTORY_INVALID`,
  `HEXEC_START_FAILED`, `HEXEC_TIMEOUT`, `HEXEC_EXIT_NONZERO`.
- Errors expose codes, safe key names, lane IDs, and safe paths only.

## Acceptance Criteria

### F08-AC01 — Unit child environment is structurally hermetic

- Traces to `M-08-C02`.
- A contaminated parent containing known runtime keys and an arbitrary future
  sentinel does not block unit; a child probe proves none reached the child.
- The observed child key set is a subset of the fixed safe set plus generated
  `GOCACHE`, `GOMODCACHE`, `GOPROXY`, and `GOSUMDB`; both caches are
  repository-local and parent values remain unchanged.

### F08-AC02 — Execution is shell-free, CWD-stable, and lossless

- Traces to `M-08-C02`.
- Argument round-trip, explicit CWD, separate large stdout/stderr drains,
  nonzero exit preservation, timeout, and parent immutability pass fixtures.

### F08-AC03 — Secret candidates are redacted before propagation

- Traces to `M-08-C04`.
- Raw, assigned, embedded, URI, stdout, and stderr candidate values are absent
  from returned output, dispatcher output, and minimal run summary.

### F08-AC04 — Public commands remain compatible and fail closed

- Traces to `M-08-C02` and `M-08-C04`.
- All eight npm aliases execute their expected behavior from a foreign CWD;
  unsupported lane/target fails before process start and before run creation.

### F08-AC05 — Governance and context behavior do not regress

- Traces to `M-08-C09`.
- Existing governance/context suites and current repository commands pass with
  no broadened scanner exemption or evidence inflation.

## Non-Goals

- F-03 owns database creation, migrations, fixtures, idempotence, and cleanup.
- F-04 owns rich manifests, aggregation, cold gate, and reproducibility.
- F-05 owns leases, run/session state, recovery, worktrees, and Codex surfaces.
- F-09 owns the complete regression corpus/dogfood closure.
- No product, OpenAPI, SDK, migration, Docker, provider, or UI behavior changes.
- No live integration evidence is claimed by F-08 fixtures.

## Stop Conditions

- The implementation needs a broader lane allowlist than the accepted
  registries declare: stop and amend both contracts intentionally.
- A public alias requires database/network/provider behavior owned by a later
  feature: preserve its fail-closed blocker.
- A test can pass only by mutating parent environment/CWD, invoking through a
  shell-built command, weakening redaction, or exempting the new module from
  governance: stop and redesign.
- Unrelated dirty paths or a shared-seam owner conflict appears: stop without
  reset, revert, stash, clean, checkout, or restore.

## Handoff

- Current status: `spec_ready`.
- Next owner: Fresh Feature Implementer.
- Next action: Execute Phase 1 RED fixtures, preserve intended failures, then
  implement environment/process modules through GREEN.
- Evidence class: fake/contract only.
- Blockers: None.
