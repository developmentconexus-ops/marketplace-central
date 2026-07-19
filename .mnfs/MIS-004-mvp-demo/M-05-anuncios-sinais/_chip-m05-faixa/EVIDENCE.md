# CHIP-M05-FAIXA — Evidence Pack

**Branch:** `chip/m05-faixa` · **Base:** `ffd8d2af` · **Head:** `d44ad66c`
**Commits (4):**
- `693dc32c` feat(market): competitor faixa (min—median—max) + own-seller exclusion
- `487a1cc0` fix(market): propagate OwnSellerID error, don't fabricate NO_PRICE_EVIDENCE
- `050d93da` feat(web): faixa de mercado in listing drawer + honest ML target label
- `d44ad66c` harden(market): trim both sides of own-seller comparison

Full diff: [`chip-diff.patch`](./chip-diff.patch) (656 lines, 18 files).

## Deliverables

| # | Mandate | Where | Status |
|---|---------|-------|--------|
| 1 | Faixa de mercado (min—mediana—máx + concorrentes) in drawer | `ListingDetailPanel.tsx` `FaixaCard` + grid | ✅ |
| 2 | Own-seller exclusion at market **application** seam, dynamic id, own offer kept in raw evidence | `collection_pipeline_service.go` `excludeOwnSeller`; `OwnSellerID` port + adapter (`ExternalAccountID`) | ✅ |
| 3 | `max_valid` — migration 0070 + pipeline compute + runner count 57→58 | `0070_*.sql`, `aggregation.go`, `aggregate_repository.go`, `runner_test.go` | ✅ |
| 4 | Label honesty: "Preço p/ vencer" → "Alvo buybox (ML)" | `ListingDetailPanel.tsx` | ✅ |
| 5 | ADR-17: absent stats render "—", never 0/fabricated | nil-safe at every seam; no `omitempty` on signal faixa tags; FE `formatMoney`→`UnknownValue` | ✅ |
| — | IC-03 read chain (EvidenceAggregate + MarketSignal carry median/min/max) | `listings/ports/evidence.go`, `listings/domain/signal.go`, `read_service.go`, composition adapter | ✅ |
| — | OpenAPI + sdk-runtime updated together | `marketplace-central.openapi.yaml` `ListingReadMarketSignal`; `sdk-runtime/src/index.ts` `ListingMarketSignal` | ✅ |

## Architecture

- Market domain `ComputeMarketAggregate` stays **pure** — exclusion lives in the application layer only.
- `listings` never imports `modules/market`; the D-21 composition adapter (`listingsEvidenceAdapter`) is the sole bridge, doing the anti-corruption map.
- Own seller id resolved dynamically: `OwnSellerID` → `accountRefForTenant` → `inst.ExternalAccountID` (ML `/users/me` id, same field ML offers carry as `SellerID`). **No `691607102` (or any literal seller id) in production code** — verified by both reviewers.
- Migration 0070 CHECK is symmetric with the existing `min_valid` constraint (nil pair XOR positive-BRL pair).

## Gates

**Go ladder** (`apps/server_core`, hermetic `GOCACHE`/`GOMODCACHE`):
- `go build ./...` — clean
- `go vet` market/listings/composition/migrate — clean
- `go test -count=1` market/listings/composition/migrate — **all green** (incl. migrate count = 58)

**Web vitest** (chip junction to main `node_modules` + throwaway `vitest.chip.config.ts`, deleted pre-commit):
- `src/pages/ListingDetailPanel.test.tsx` — **14/14 pass** (faixa render, label rename, ADR-17 "—" for null bounds)
- Full suite — **348 pass**, 44/45 suites. See FINDING below for the 1 failing suite.

**P6 dual gate — Claude side (cold Opus + adversarial reviewer, agreement required):**
- Cold Opus (`general-purpose`, opus): **VERDICT PASS**, blockers: none.
- Adversarial (`general-purpose`, sonnet): **VERDICT PASS**, blockers: none.
- **AGREEMENT: PASS.** Both independently flagged one non-blocking point (untrimmed own-seller comparison) — applied in `d44ad66c`.

## FINDING (out of scope, pre-existing) — for hub ratification

`apps/web/src/app/AppRouter.test.tsx` fails to load: `Failed to resolve entry for package "@marketplace-central/feature-simulator"`. Cause is **not** this chip — `packages/feature-simulator/` is a committed **stub on base `ffd8d2af`** (only `package.json`, no `main`/`module`/`exports`/`src`). The suite fails identically before any M05-FAIXA change. Untouched by this chip; reported so the milestone owner can decide whether feature-simulator needs its entry/build before demo.

## Live-drive QA notes (from reviewers, for the P7 browser gate)

1. Runtime own-seller match relies on `ValidatedOffer.SellerID == inst.ExternalAccountID` (both the ML numeric user id). Confirm on the live 90008 case that our MLB offer drops out of the competitor min/median/max/n and the faixa reflects competitors only.
2. `OwnSellerID` error aborts catalog-evidence collection (fail-honest, not a regression — `ListCatalogOffers` already resolved the same ref moments earlier).

## Hard-rule compliance

Zero ML writes. No server booted / `:8080` / `.env` / docker touched. No push/reset/revert/stash/clean/`-D`/WSL/cold-clone/cache-purge. Only dep-adjacent action = migration 0070 (pre-granted). `GOCACHE`/`GOMODCACHE` absolute inside `apps/server_core`; `gomodcache`/`node_modules`/chip-config all absent from the tracked tree (verified). Throwaway vitest config deleted pre-commit.
