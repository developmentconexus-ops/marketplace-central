# Slice 11 brief — round-3 corrective + G1 ruling receipt

```yaml
base: 7f5a1b8c
scope: apps/server_core/internal/modules/listings/application/ ONLY (read_service.go + read_service_test.go)
out_of_scope: internal_read/ (cross-module seam, hub-owned), composition, transport, OpenAPI, SDK, migrations
gate_that_produced_this: dual-gate round 3 — FAIL both sides (_gate-evidence/round-3/dual-gate-round3-verdict.md)
```

## G1 — HUB RULING: **(C)** — adjudication receipt (VERBATIM-BINDING for round 4)

The hub verified the chip's receipts independently before ruling (`reader.go:508-513`,
`safeOracleCause:515-523`, `read_error.go:8-9`, `Unwrap:28-33` — all confirmed), and found the seam
**wider than the chip reported**: `sankhya_linkage_reader.go:424` and `open_cgo.go:62` share
`safeOracleCause`, and `catalog/transport/http_handler.go:289` additionally string-matches
`strings.Contains("source_unavailable")` on top of `IsReadErrorCode`. **4+ modules consume it.**
Option (A) inside M-01 = a chip writing into a 4-module seam, plus scope widening against the
operator directive. **(A) REJECTED.**

Ruling (C), two parts:

1. **In-module fix, slice 11:** `errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)`
   checked in `isSourceUnavailable` **BEFORE** `IsReadErrorCode`. Cancellation/timeout never degrades
   — it propagates and logs at ERROR. Slice 10's strict classification already covers the rest.
   Test MUST inject a **real `context.Canceled` cause travelling through the `ReadError` wrap** — NOT
   a fake carrying the right code. Fakes-with-the-right-code are exactly the hole that hid this.
2. **Residual deferred WITH A NAME:** an adapter/data defect still degrades, misclassified but
   logged at WARN (B2b landed). → **hub board task #2**: `internal_read` taxonomy refactor,
   hub-owned, post-M-01. Registered scope: new codes, reclassify the 3 sites, preserve the sentinel
   instead of flattening it at `:519-523`, update 4+ consumers including killing catalog's
   string-match.

### Round-4 prompt-pack MUST carry this receipt

> **B2 = resolved in-module (cancel/timeout) + residual deferred-with-name to hub task #2 per G1
> ruling (C).**

Without it a reviewer re-blocks on a residue that has already been adjudicated. The §9 resolution
check consumes the receipt rather than re-litigating it.

## G2 (important) — revert a regression the round-2 gate had certified

`read_service.go:337` guard reverts `unavailable != nil` → `ceilingErr != nil`.
`icmsWorstCaseByUF:563-567` fills `WorstCaseICMSPct` and `PriceNetBasis` (`priceNetBasis:574-576`)
from **price + ceiling only** — no cost dependency — and `belowMargin:520-536` already nils on a nil
cost. So a **cost-only** outage must keep the FULL matrix with only `below_margin_at_uf` null.
Nulling the whole matrix discards known-true facts because an unrelated fact is missing. ADR-17
licenses null for *unknown*, never for known.

**Owned honestly:** `slice10-review.md:29-31` (mine) blessed this as a *"BONUS fix… closing a latent
gap."* It was a regression. The §14 slice review is the milestone owner's, so the misclassification
is the milestone owner's. The cold gate caught what the slice review had already approved — the
hub is recording this in the ledger as evidence for the standard itself (implementer ≠ reviewer,
and the cold gate sits downstream of the slice review for exactly this reason).

## G3 (important) — pin the branch G2 left unpinned

`read_service_test.go:317` guards the matrix assertion behind `outage.ceilingErr &&`, so the
cost-outage case asserts nothing. **Bar: the test must FAIL if G2's fix is reverted.** A test that
still passes with the fix reverted is not coverage — round-2's own I2 criterion.

## G4 (blocking) — strip comment narration

`HARNESS.md:203-204`, AI-slop checklist, **any hit = REJECT the slice**: *"comment narration /
PR-voice comments"*. Slice 10 narrated adjudication history in production code ("Per the hub ruling
(B1, PURE…)", "pins B2"). Hub ACK: **decision context lives in the ledger/evidence, not in the
code.** Comments state only invariants a future reader needs; when a comment carries no constraint
the code cannot show, delete it.

## G5 (important) — telemetry must not state the opposite of the outcome

The ceiling WARN (`:117`, `:189`, `:329`) says "degrading fact to null" but fires **before** the
`needsBelowMarginScan` branch (`:153`, `:192`) decides — so a filtered request logs "degrading" and
then 503s. One approach, applied consistently at every site (List/Summary/ByProduct/Get + enrich's
cost path). Invariants that survive: ERROR on propagate, WARN on degrade, `err` + `op` (+ `fact`
where already present), at most once per request, never per row.

## Governance — HUB RULING

Do **NOT** catch the worktree's `docs/` up (merging master into a mid-flight branch = noise). Pin
reviewers to the **main checkout's ABSOLUTE paths**:
- `C:\Users\leandro.theodoro\Documents\marketplace-central\docs\REVIEW-STANDARD.md`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\docs\HARNESS.md`

Codex's cwd is already the main root, so it resolves. The worktree's stale `AGENTS.md` is harmless
to this flow — **do not cite `docs/superpowers/`**.

## Must not break (all pinned, all currently green)

- Any fact-dependent filter (`has_exception` both directions AND `exception=below_margin`, uniform
  via `needsBelowMarginScan`) during a source-unavailable outage → error, never a page built by
  dropping the filter.
- A read WITHOUT a fact-dependent filter → still degrades to null (ADR-17), not 503.
- `matchesDependentFilter` genuinely applied on the happy path.
- Only `ReadErrorSourceUnavailable` degrades; everything else propagates.
- Telemetry once per request, never per row.

## DoD

L0 (build+vet) → sonnet §14 review ∥ L1 → integration lane → COMMITTED (restart pre-armed) →
re-drive → gate round 4 (delta since `7f5a1b8c` + resolution check of EVERY round-3 finding, G1
consuming the ruling receipt above) → P8 CLOSED. **No push.**
