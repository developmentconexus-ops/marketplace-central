# F-02 write-types-adapters — validation record

Chip: CHIP-M03 · branch `chip/m-03-mutation-envelope-writes` · closed 2026-07-16.

## Slices → commits → reviews

| Slice | Commits | Review (ledger) | Verdict |
| --- | --- | --- | --- |
| S1 seven provider-write gates (WriteThroughGates) | cca00867 | D-18 | ACCEPT |
| S2 price/listing write capability contracts (explicit-only exposure) | 2395bc0f → REJECT → 6dd51ff0 (StockReads→StockWriter auto-promotion removed) | D-20/D-23 | ACCEPT after fix |
| S3 ML PriceWriter (absolute price only) | f3592993 | D-22 | ACCEPT |
| S4 ML listing pause/edit + resync writer | 1bbf1bd5 → REJECT → 5635b494 (honest state mapping + parameterized sanitizer) | D-25/D-27 | ACCEPT after fix |
| S5 product-link + policy adapters | 75e4bb8c → REJECT → caf44cc8 (resolved-link gate scoped to gated listing) | D-28/D-29 | ACCEPT after fix |
| S6 six-intent writer router + IC-03 failure mapping | ae727dc1 + d0248d74 (G2 note) | D-31 | ACCEPT-WITH-CONDITIONS, all closed (G2 note; MapFailure precedence confirmed; production callers landed S8a) |
| S7 stock_correct via envelope + StockActionService fold (migration 0039) | d14a94d7 → REJECT → 8ae27443 (canonical no-variation sentinel; facade honest pending semantics) | D-33/D-34 | ACCEPT after fix |
| S8 (original) | — | D-35 | BLOCKED honest (plan-shape conflict) → replan D-36 → S8a/S8b (plan.md 262d5c67) |
| S8a poller-owned gates + ML writer bridges | 7752b64f | D-38 | ACCEPT (incl. focused S6 delta re-review per D-36) |
| S8b gated composition wiring + capability exposure + runner guard | b4fa8281 | D-40 | ACCEPT (gate = blocking dimension, proven) |

## Contract evidence (M03-C* served: C04 write types, C05 failure model, C06 adapters, C07 gates)

- Six enabled intents dispatched by WriterRouter: price_update (absolute only), stock_correct,
  listing_pause, listing_edit, listing_resync, link actions (approve_candidate /
  manual_resolve / reject_listing via dedicated `LinkageWriter.ApplyLink` — `skipped`
  survives; WriteOutcome cannot carry it). listing_create → typed
  `FailureCodeTypeNotEnabled`.
- Gate ownership (replan D-36, option 3): poller EXCLUSIVELY owns idempotency + audit —
  WriteThroughGates at poller layer on claim's `AppliedIdempotencyKeys` + `MarkItemApplying`,
  exactly once per item; router carries no gate callbacks (reviewer traced byte-equivalent
  behavior vs pre-change inline logic).
- IC-03 failure mapping: MapFailure (network check first, typed `ErrorCodeOf` override after
  — typed code carries retryability) + MapRejected; retryable ONLY provider_rate_limited /
  provider_unavailable; sanitized message_provider passthrough only; production callers
  `failedWrite`/`rejectedWrite` reused by all three bridges (D-31 carry closed).
- sku_invariant_violation raised BEFORE any adapter call (zero-adapter-call proven S6).
- Idempotency key `{protocol_id}:{listing_id}` forwarded end-to-end to provider writes
  (X-Idempotency-Key via doRawWithIdempotency — F-01 load-bearing carry, re-verified D-38).
- Capability exposure explicit-only (F02-S2 doctrine): PriceWrites (S3) + StockWrites +
  ListingWrites (S8b) as explicit ProviderCapabilitySet fields; auto-promotion removed
  6dd51ff0; absent capability → typed unsupported failure, never a silent stub.
- StockActionService fold (S7): ApplyManual delegates to mutation envelope
  (CreateStockCorrection → protocol), FindByID idempotency short-circuit, honest-null
  `mutation_protocol_id` (migration 0039, nullable TEXT, FK absence deliberate — D-36);
  facade approved+protocol → risk blocked `mutation_protocol_pending`, zero optimistic
  snapshot writes (ADR-17).
- Composite listing id canonical round-trip: `installation~item~variation` with
  `NoVariationID = "-"` sentinel normalized both directions (S7 fix, envelope round-trip
  tests).
- Per-item installation resolution in composition (S8b): stock/resync writers resolve
  installation per mutation item; not-found → honest error → FailureCodeInternal
  non-retryable; no manufactured account facts.
- Run() re-entrancy guard: sync.Mutex reset-on-return; concurrent re-entry no-op + restart
  after return, both proven with real goroutines (F-01 carry closed).

## Lanes

- Unit: green — composition + mutations + inventory + connectors sweeps `-count=1`,
  re-run cold by each reviewer AND by the orchestrator pre-pass per slice.
- Build/vet: green every slice.
- `-race`: unavailable on this machine (no cgo/gcc) — flagged D-40 for hub confirmation
  elsewhere; not a slice defect.

## Live-write posture (security-critical)

Real provider write lane is wired ONLY when `MPC_PROVIDER_WRITES_ENABLED` is trimmed
case-insensitive `true` — proven the single switch (every root.go path traced, D-40);
`"1"`/`"yes"`/unset/empty/`"false"` keep the F-01 stub. DEFAULT OFF. Flipping the flag is an
explicit operator deployment action; nothing in code auto-enables. No live Mercado Livre
write was executed at any point during F-02 — all provider behavior proven via contract
mocks/stubs. Live ML execution remains gated on explicit operator authorization via hub
ESCALATION.

## Carries / hub-queue items out of F-02

- Provider-code asymmetry (D-40 question): price/listing bridges constructed with hardcoded
  `"mercado_livre"` while stock/resync resolve per installation — honest today
  (single provider registered); revisit when a second write-capable provider lands.
- Helper placement suggestion (D-38): `failedWrite`/`rejectedWrite` live in price_writer.go,
  could move to failure_mapping.go — cosmetic.
- Named constant for inline `domain.ErrorCode("CONNECTORS_INTERNAL")` (S6 carry,
  non-blocking).
- Edit-on-paused listing maps to honest Rejected — M-06 UI must surface this (cross-milestone
  note for hub).
- F-03 must expose the envelope over HTTP (8 endpoints) — OpenAPI + sdk-runtime same commit.
