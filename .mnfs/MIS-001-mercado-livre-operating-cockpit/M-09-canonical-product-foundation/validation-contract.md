# M-09 Validation Contract

```yaml
id: M-09
type: milestone-validation-contract
status: ready
owner: Mission Strategist
parent: MIS-001
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: milestone
```

## Required Outcome

One canonical positive Sankhya CODPROD crosses Oracle, catalog, API, SDK, product
links, and UI inputs without active legacy MSDB dependency or zero-filled unknowns.

QA records the frozen SHA once before review and commands. SHA drift stops the run;
per-command atomic wrappers are not required.

## Criteria

### M-09-C01 — Canonical identity

- Required: Yes.
- Proof: targeted catalog, internal_read, product_links, API/OpenAPI, and SDK tests.
- Expected: `internal_product_id` is a positive CODPROD; EAN, manufacturer reference,
  seller SKU, item ID, and variation ID remain separate nullable identifiers.
- Blocks: any provider/manufacturer identifier occupies the internal product field.

### M-09-C02 — Honest nullable facts

- Required: Yes.
- Proof: catalog/internal_read domain and serialization tests.
- Expected: missing cost, price, or stock is null with `unknown` and a reason; a
  known zero remains numeric zero with current source metadata.
- Blocks: missing data and known zero serialize identically.

### M-09-C03 — Legacy runtime cutover

- Required: Yes.
- Proof: targeted residue scan plus composition/catalog tests.
- Expected: active server composition and catalog reads no longer require
  `platform/msdb`, `MS_DATABASE_URL`, or `MS_TENANT_ID`.
- Blocks: server startup or the MVP catalog path still depends on MSDB.

### M-09-C04 — Deterministic compatibility

- Required: Yes.
- Proof: migration/readback tests for mapped, unmapped, and ambiguous records.
- Expected: only proven CODPROD equality preserves a link; zero candidates remain
  `not_found`, multiple candidates remain `identity_conflict`, and neither is guessed.
- Blocks: ambiguous or absent identity is attached to a CODPROD.

### M-09-C05 — Real Oracle product read

- Required: Yes for milestone Pass; unavailable credentials/source yield
  `externally_blocked`.
- Proof: governed read-only Oracle smoke for at least one positive CODPROD.
- Evidence: sanitized command result with frozen SHA, source, observed time, and
  `read_only=true`; no raw row, credential, or PII.
- Blocks: only fixtures prove the Oracle cutover or the lane mutates Oracle.

## Proportional QA

1. Review the frozen SHA and changed paths.
2. Run the targeted Go and SDK tests named by the feature validations.
3. Run broader server tests when shared composition/API seams changed.
4. Run the real Oracle read once after deterministic proof passes.
5. Verify OpenAPI and SDK changed together when the public shape changed.

## Evidence Requirements

- Each feature writes `validation.md` with changed paths, commands, outcomes, and
  honest limitations.
- QA writes `validation-result.md` and may keep concise sanitized outputs under
  `_fixed-sha-qa/`.
- No manifest schema, OCR, or evidence-wrapper implementation is required.

## Retry Policy

Maximum two scoped correction attempts. Stop immediately on an architecture,
contract, ownership, live-write, or nondeterministic identity conflict.

The ordinary budget was consumed at SHA
`230dc78306d3775894a00b5424238529382cc9b0`. Portfolio authorizes exactly one
additional attempt under `correction-contract-c01-final.md`, limited to the two
remaining M-09-C01 findings. That contract does not reset the ordinary budget and
permits no further retry or scope expansion.

After C01 review passed, proportional QA at
`97fd4b58d55a7d14a2b45f0c3bae15b2e374822a` exposed a historical test-clock defect
outside the M-09 implementation diff. Portfolio separately authorizes one test-only
QA unblock under `correction-contract-qa-inventory-clock.md`. It does not reopen any
M-09 criterion, implementation path, or retry budget.

## Handoff

- Current status: Ready.
- Next owner: M-09 Milestone Orchestrator, then `mpc-verifier`.
- First work item: F-01 canonical product identity.
- Open decisions: none; an unmappable legacy identity is a visible Milestone stop.
