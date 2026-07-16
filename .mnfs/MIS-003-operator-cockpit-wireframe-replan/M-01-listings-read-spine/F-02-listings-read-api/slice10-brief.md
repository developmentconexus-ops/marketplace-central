# F-02 Slice 10 — corrective brief: fail-honest degraded filters + error classification + telemetry

```yaml
slice: 10 (corrective, M-01, same chip/worktree)
origin: dual-gate DELTA round 2 FAIL — both sides (cold Opus + GPT-5.6 Sol), independent, in agreement
hub_ruling: R1 = (C) PURE; R2 scope CONFIRMED. Adjudication 2026-07-16 (HUB-EVENT-BLOCKED-dual-gate-delta-fail.md answered)
base: 1f6b72d8 (c4e8ab91 slice 9 + gate evidence)
no_push. GOCACHE absolute. no self-boot. slice 9 untouched.
model_deviation: implementation on sonnet (operator directive, overrides HARNESS §1 Luna-high) — logged in DISPATCH-LEDGER
```

## R1 RULING — (C) PURE, uniform fail-honest

- **Rule:** a request using **ANY** fact-dependent filter (`filter.has_exception` **AND** `filter.exception=below_margin` — uniform, no special case) during a fact-source outage → **source-unavailable ERROR via the EXISTING error envelope**. No OpenAPI/SDK change: a 503-class response in this scenario was the pre-slice-8 behavior, so a subset of it does not widen the contract.
- **Reads without a fact-dependent filter:** keep slice 8's degrade-to-null (ADR-17 correct — null = not-evaluable, never a false fact).
- **(A) REJECTED** — applying only the SQL-computable component yields a silently **under-inclusive** result: an operator triaging exceptions would see an incomplete set with no signal. It also redefines the filter's semantics without ratification. A lie in the other direction.
- **(B) DEFERRED WITH A NAME (G3-traceable)** — an explicit `service_degraded` / "filter not applied" signal is the RIGHT pattern *once the cockpit UI milestone exists to consume it*. Today there is no second consumer → YAGNI + widens the contract at close time. **Declared follow-up of the cockpit UI milestone.** Recorded in the mission ledger, not silently dropped.

## R1 TELEMETRY (2nd blocking — mandatory in slice 10)

- Every degrade activation → `slog.Warn` with **reader identity + error class**, **once per request, never per row** (no spam).
- **Strict classification:** ONLY `ReadErrorSourceUnavailable` degrades. Cancellation / timeout / adapter defect **PROPAGATE** as errors (never a silent null) and log at **ERROR**.
- Parity with the unmapped-status WARN: a total outage is never silent again.

## R2 SCOPE — CONFIRMED (slice-8 degrade path ONLY)

`passThrough` + `passThroughGroups` (semantics above) + error classification + telemetry + removal of the dead `unavailablePolicyService` wrapper (reviewers' unreachability verification accepted) + re-pin of the tests that encode the defect (`TestReadServiceDependentFilterPassesThrough…` + degraded coverage for `has_exception`/`below_margin` in **both filter directions**). **Slice 9 untouched.**

## Loci (from the merged verdict's verified anchors)

| # | File | Locus | Change |
|---|---|---|---|
| 1 | `internal/modules/listings/application/read_service.go` | ~`:365-375` `passThrough`, ~`:228-239` `passThroughGroups` | fact-dependent filter + outage → source-unavailable error; no filter → degrade unchanged |
| 2 | same | ~`:444-446` `needsBelowMarginScan` | already true for `HasException != nil` and `below_margin` — this predicate IS the trigger condition |
| 3 | same | ~`:394` cost error, ~`:122` `ceilingErr` | strict classification; only source-unavailable degrades; others propagate + log ERROR |
| 4 | same | degrade sites | `slog.Warn` reader + error class, 1×/request |
| 5 | `internal/composition/root.go` | `:106-116` | drop dead `unavailablePolicyService`; pass `listingsmarketplaces.NewPolicyReader(marketSvc)` directly (~`:486`) |
| 6 | same | ~`:246-249` `enrichGroups` | nit: return real `factsUnavailable` instead of `false` on the error path |
| 7 | `application/read_service_test.go` | ~`:499-564` | re-pin `TestReadServiceDependentFilterPassesThrough…` to the correct behavior + rename; add both-direction degraded coverage |

## DoD (pre-commit)

- L0: `go build ./...` 0, `go vet ./...` 0 (whole repo, absolute GOCACHE).
- L1 unit: `go test -count=1 ./internal/modules/listings/... ./internal/composition/...` 0. Integration lane if any migration touched (none expected).
- Test-first RED→GREEN proof for the B1 block and the classification block.
- §14 sonnet cold review APPROVE (no blocking / unresolved important), run ∥ L1 per §15.

## Flow (slice-9 protocol, per hub)

plan → implement → L0 green → [§14 sonnet cold ∥ L1] → COMMITTED event → hub restart pre-armed → **gate round 3 = §9 DELTA since c4e8ab91 + explicit resolution check (resolved/unresolved) of every round-2 finding** — does NOT re-litigate what was adjudicated → P8 CLOSED. No push.
