# CHIP-T2-MIN — degrau 3 (COTAÇÃO) — Evidence Pack

Branch: `chip/t2-min-degrau3`
BaseSha (fork): `40f7f4deaaece5be198789cf27cad28241a63add`
Tip SHA: `b3b1a23520ef91174c3344fa20fbe91e9702c89b`
Lane: Claude-only (Codex quota exhausted til 2026-07-25).

## Scope delivered
Degrau 3 = live COTAÇÃO commission tier of the pricing tariff ladder. Thin adapter
binding the already-merged ML catalog-match + fee-quote read machine to the T1
`TariffResolver` port. Miss/error in the live chain falls SILENTLY to degrau 4
(config default); never fabricates, never panics.

OUT OF SCOPE (untouched): persistence history/migration 0069, frete-by-dims,
buy_box F1 investigation, encoding fix (retracted D-85).

## Commits (7 TDD slices, green per slice)
- S1 `be103b6` domain FonteCotacao const + ComponentResolution.Data + toTarifaDTO maps Data
- S2 `d3d0054` 3 ports (ProductIdentityReader/CategoryResolver/CommissionQuoter) + TariffRequest.ProductID/PriceBasis
- S3 `bce00e2` ProductIdentityReader adapter (catalog svc)
- S4 `198745f` tarifflive assembler (silent-miss matrix; D-83 EAN+titulo together; modalidade→listing_type inside assembler)
- S5 `2627b02` tariffcomposite chain (degrau3→degrau4; frete stays degrau4; base error propagates)
- S6 `1212f0e` calc_service threading + manual-commission gate
- S7 `b3b1a23` composition wiring (pricing_adapters.go raw readers + root.go composite behind FeeQuoteReader availability)

## Build / vet / test (final, GOCACHE absolute, no GOFLAGS)
- `go build ./apps/server_core/...` → clean
- `go vet ./apps/server_core/internal/modules/pricing/... ./apps/server_core/internal/composition/` → clean
- `go test -count=1` pricing (12 pkgs) + composition → all `ok`:
  catalog, costread, marketplace, postgres, tariffcomposite, tariffdefaults,
  tarifflive, application, domain, ports, transport, composition.
- IC-04 freeze: TestDecompositionShapeFrozen / TestDifalForUFResultShapeFrozen /
  TestDecomposeInputShapeFrozen / TestDifalForUFPortReflectsOverrideAndScopesTenant → PASS.

## Contract / SDK
ZERO change. Degrau 3 fills the SAME `PricingTarifa` fields (`data` was already a
nullable wire field, previously hardcoded nil; `fonte` has no OpenAPI enum). No
openapi/sdk-runtime file in the diff. calc_ports_contract_test.go untouched.

## Key invariants verified
- **Unit match (R$79-class guard):** ML `percentage_fee` is a percent (11 == 11%);
  domain `ComissaoPct` string is a percent; `domain.pctOfPrice` divides by 100 exactly
  once downstream. No ÷100/×100 mismatch. `FormatFloat('f',-1,64)`.
- **ADR-17 (unknown ≠ zero):** nil CommissionPercent → miss; ≤0/unparseable price → miss;
  absent catalog price → Preco nil; no category → miss. Never a fabricated 0/default.
- **Fallback ownership:** composite is the SOLE owner of degrau3-error→degrau4 silent
  fallback. Assemblers + composition adapters PROPAGATE real read errors (found=false only
  for honest "no data"). No blanket recover/catch-all.
- **Manual gate:** explicit user commission ⇒ ProductID/PriceBasis not threaded ⇒ degrau 3
  skipped before any I/O (no ML call).
- **Frete:** composite only ever mutates `base.Comissao`; Frete stays degrau 4.
- **ZERO ML writes:** both provider calls read-only probes (ReadCatalogMatch / ReadFeeQuote),
  wired to raw capabilities (NOT the integrations probe wrapper → no provider_operations
  audit spam per decompose).
- **Boundaries:** ML payload confined to composition adapters; pricing domain/ports carry
  only strings + domain.Money. No import cycle.
- **D-83:** catalog-match probe always sends EAN + Query(titulo) together.

## P6 dual-gate — AGREEMENT: PASS
- Cold Opus reviewer (no prior context): **VERDICT PASS**, zero blockers. Verified all 7
  gates (correctness/unit, architecture, fallback discipline, tenancy, contract, test
  quality, zero writes).
- Sonnet adversarial reviewer (break-it): **VERDICT PASS**, no exploit found after probing
  nil-panic paths, fabrication, price-basis confusion, listing_type mapping, manual gate,
  frete isolation, contract drift. Ran the suites green.

Both AGREE → gate satisfied.

## Non-blocking observations (relayed to hub, out-of-scope for MIN)
1. `firstConnectedMLAccountRef` called independently by BOTH category adapter and commission
   quoter ⇒ two `installations.List(ctx)` per decompose (each runs reconcileState +
   projectRuntimeCapabilities). Hot-path optimization candidate: resolve account ref once.
2. `CategoryResolution.Data`/`Confianca` populated by composition adapter but discarded by
   tarifflive assembler (final Data uses commission fee timestamp). Harmless dead data.
3. Composite swallows degrau-3 errors with no log/metric ⇒ a persistently failing live probe
   is invisible in prod. A debug counter would help operability.

## Close status
Chip gate (Go + governance + P6) satisfied. Chip close ≠ M-07 close — only fresh browser QA
passes the milestone.
