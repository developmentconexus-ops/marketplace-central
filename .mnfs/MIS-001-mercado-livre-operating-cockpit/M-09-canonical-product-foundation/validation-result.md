# M-09 Validation Result

```yaml
milestone: M-09
mode: proportional_qa
verdict: failed
frozen_sha: 32b32f6de00875589468c71eb70c6eb3e5d49278
observed_at: 2026-07-13T22:41:48.6614485Z
sha_drift: false
```

## Verdict

**FAIL.** Deterministic Go, SDK, and runner-contract lanes passed, but the
required exact active-residue scan failed at the frozen SHA. QA stopped before
the live Oracle lane as required.

## Criteria

| Criterion | Assessment | Evidence |
| --- | --- | --- |
| M-09-C01 | PASS | Both prerequisite fixed-SHA reviews passed; targeted catalog, internal-read, product-links, transport, and migration-related Go packages passed; SDK/OpenAPI parity passed 40/40. |
| M-09-C02 | PASS | Targeted catalog/internal-read domain and adapter tests passed, including the named nullable-source-fact coverage distinguishing unknown from known zero. |
| M-09-C03 | FAIL | The exact active-residue scan found legacy `MetalShopping` wording in active catalog files and `MS_DATABASE_URL` in `docker/dev/README.md`; the registered no-match expectation was not met. |
| M-09-C04 | PASS | Targeted classifications, pricing, migration runner, and embedded migration contract tests passed; named feature evidence covers mapped, `not_found`, and `identity_conflict` outcomes without guessing. |
| M-09-C05 | NOT COMPLETED | The live read-only Oracle lane was not run after the C03 stop condition. Existing Oracle evidence names an older SHA and was not refreshed or accepted for this SHA. |

## Completed command evidence

- Targeted Go lane — PASS.
- Full `apps/server_core` Go lane — PASS, including inventory application.
- SDK/OpenAPI lane — PASS (40/40).
- Governed runner contract — PASS (14/14).
- Active residue scan — FAIL with four matches across three catalog/dev files.
- Live Oracle runner — NOT RUN due mandatory stop.

Concise sanitized evidence is recorded in
`_fixed-sha-qa/deterministic-qa.md`. No source modification, database/provider/
Oracle write, sensitive-data exposure, or SHA drift occurred.

## Terminal disposition

- status: `failed`
- blocker: M-09-C03 active legacy residue and registered residue-scan failure
- remaining required proof: corrected no-match active-residue scan, then a
  frozen-SHA governed read-only Oracle product read
- retry: none authorized by the named correction contracts
