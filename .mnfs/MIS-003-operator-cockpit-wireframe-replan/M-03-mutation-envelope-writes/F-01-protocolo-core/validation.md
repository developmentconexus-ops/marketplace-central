# F-01 protocolo-core — validation record

Chip: CHIP-M03 · branch `chip/m-03-mutation-envelope-writes` · closed 2026-07-16.

## Slices → commits → reviews

| Slice | Commits | Review (ledger) | Verdict |
| --- | --- | --- | --- |
| S1 schema 0038 (mutation_protocols + mutation_items + claim index) | 6fd596e4 | D-05 | ACCEPT |
| S2 lifecycle domain (8 protocol / 6 item states, 12 failure codes) | 2f3231ea | D-04 | ACCEPT |
| S3 MP-%06d tenant-scoped IDs + draft create/read repo | 22b1f0bb | D-07 | ACCEPT |
| S4 item snapshot + SKIP LOCKED claim repo | 3f176541, 40e5a52 | D-09 | ACCEPT (applied_at honest-null fix folded in) |
| S5 WriterPort + programmable stub + poller pass | 36745c9b → REJECT → 733a9a90 (durable outcomes) + bbd2b603 (applying mark) | D-11/D-13 | ACCEPT after fix slices; all conditions closed |
| S6 crash-resume integration + background runner + composition stub wiring | 8795a063 → REJECT → 72197bf5 (runner unit tests) | D-15/D-16 | ACCEPT |

## Contract evidence (M03-C* served: C01 schema/lifecycle, C02 poller/resume, C03 idempotency, C05 failure model)

- Lifecycle: 64/64 transition cases, 6/6 terminal, 12/12 retryability table-driven (domain).
- IDs: tenant-scoped monotonic MP-%06d via `pg_advisory_xact_lock` + MAX+1, race/crash-safe
  (precedented idiom), cross-tenant reads not-found.
- Claim: FOR UPDATE SKIP LOCKED, one protocol per (tenant, installation); second concurrent
  claim gets nothing; different installations claim concurrently (integration-proven).
- Crash safety (reviewed adversarially, REJECT→fix cycle): item outcomes + applying marks are
  pool-durable OUTSIDE the claim tx; crash reverts protocol to approved while outcomes
  survive; resume never resends a durably-applied idempotency key (unit + real-postgres e2e:
  `TestPollerPassCrashResumeDoesNotResendDurableAppliedItems` + poller_integration_test.go).
- Item audit immutability: terminal outcome UPDATE rejected in SQL predicate + affected-rows
  check; applied_at NULL unless state=applied (honest-null).
- Items reach IC-03 `applying` state durably before each provider send (re-review condition).
- Background runner: tick-driven, cancel-clean, log-and-continue — unit-proven with injected
  ticks; inert in composition (nothing starts it; stub writer carries dated deferral
  `2026-07-16, replaced by F02-S8`).

## Lanes

- Unit: green (all mutations packages + composition).
- Integration (real postgres, session container, 38 migrations): green — repeated per slice;
  full-repo lane GREEN-with-allowlist (TestPhase1SmokeFlow cited; orders flake F-B once,
  passed on re-run).

## Load-bearing carry-forwards

- F-02 adapters MUST forward idempotency key `{protocol_id}:{listing_id}` on every provider
  write — the send↔outcome-commit crash window closes only via provider-side dedup (plan.md
  F-02 header note).
- F-02 (S8 live wiring): add Run() re-entrancy guard; context.Canceled log cosmetics.

## Live-write posture

No live Mercado Livre writes anywhere in F-01. WriterPort implementor is a programmable stub
proving contract behavior. Live ML write lane remains gated on explicit operator
authorization via hub ESCALATION.
