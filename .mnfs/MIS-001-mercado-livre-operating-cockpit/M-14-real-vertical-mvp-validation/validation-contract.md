# M-14 Validation Contract

```yaml
id: M-14
type: milestone-validation-contract
status: ready
owner: Mission Strategist
parent: MIS-001
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-4
lifecycle_scope: milestone
```

## Required Outcome

At one reviewed SHA, a bounded values-minimized sample proves real read-only Mercado
Livre and Sankhya Oracle access, PostgreSQL application state, and the object-centered
browser journey through a visibly non-executing stock simulation.

## Criteria

### M-14-C01 — Real source provenance

- Required: Yes.
- Proof: governed Oracle and Mercado Livre read commands.
- Expected: sanitized evidence names source, observed time, read-only status, command,
  and reviewed SHA for both sources.
- Blocks: mock/fixture reported as real, missing provenance, or any live write.

### M-14-C02 — Vertical browser journey

- Required: Yes.
- Proof: IC-004 MVP-BD-01 at desktop and 390x844.
- Expected: Overview → Listing → Product → Listing → Sale/Margin → Listing → Stock
  Simulation retains installation/entity identity and displays source/quality/time.
- Blocks: lost context, manual identity reconstruction, or provider-execution copy.

### M-14-C03 — PowerShell command validity

- Required: Yes.
- Proof: execute the actual Windows PowerShell commands used for targeted/broad tests,
  build, PostgreSQL, and real reads; record cwd, exit code, and reviewed SHA.
- Expected: Go commands use an absolute repository-local `GOCACHE` and correct module cwd.
- Blocks: invalid host syntax, skipped test represented as pass, or SHA drift.

### M-14-C04 — PostgreSQL persistence and idempotency

- Required: Yes.
- Proof: existing integration lane plus two imports of the selected listing/order.
- Expected: application state survives reload and natural-key counts do not grow on
  the second import except for explicitly append-only audit records.
- Blocks: duplicate listing/order/link/snapshot identities or ephemeral-only state.

### M-14-C05 — Evidence security and honesty

- Required: Yes.
- Proof: QA inspects command output, notes, screenshots, and network summaries.
- Expected: no credential, buyer email/phone/document/address, raw payload, or full
  Oracle row; every artifact is labeled deterministic, PostgreSQL, real read, or browser.
- Blocks: secret/PII leak or fake evidence labeled real.

### M-14-C06 — No provider mutation

- Required: Yes.
- Proof: inspect reachable UI controls and sanitized browser/provider method traces.
- Expected: simulation shows `executed=false`; provider traffic for this drive is
  read-only; no execute control or write result is exposed.
- Blocks: any Mercado Livre mutation or a UI claim that one occurred.

## Proportional QA Order

1. Freeze SHA and review changed paths.
2. Run impacted deterministic server, SDK, and web tests plus web build.
3. Run PostgreSQL integration/idempotency proof.
4. Run one bounded real Oracle read and one bounded real Mercado Livre read.
5. Complete IC-004 MVP-BD-01.
6. Fold secret/PII/no-write findings and issue `validation-result.md`.

QA may reuse existing harness commands and keep concise sanitized outputs under
`_fixed-sha-qa/`. A custom 14-command runner, evidence manifest, OCR pipeline, hash
correlation, or atomic-output protocol is not required for acceptance.

## Evidence Requirements

- Three feature `validation.md` files.
- Fixed-SHA review.
- Sanitized deterministic/PostgreSQL/real-read command results.
- Browser interaction note plus desktop/mobile screenshots.
- Final M-14 `validation-result.md` from QA.

## Retry Policy

Maximum two scoped correction attempts. If a coherent real sample or credential is
unavailable, persist `externally_blocked`; do not substitute fixtures.

## Handoff

- Current status: Ready after M-09 and M-13 pass.
- Next owner: M-14 Milestone Orchestrator, then `mpc-verifier`.
- Blocking failures: any write, PII/secret leak, fake-as-real evidence, invalid
  PowerShell proof, duplicate state, or broken browser journey.
