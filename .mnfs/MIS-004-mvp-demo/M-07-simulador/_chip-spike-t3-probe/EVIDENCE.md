# CHIP-SPIKE-T3 — catalog-match-probe — EVIDENCE

**Purpose:** read-only probe surface so the HUB can live-validate the tier-3
COTAÇÃO-MATCH flow (EAN → ML catalog → category → listing_prices) with real xlsx
products BEFORE the full tier-3 implementation chip (CHIP-T2). De-risking spike,
not demo-critical.

- **Branch:** `chip/spike-t3-probe`
- **Base:** `3d158885bd8864d945ec7aad279f7d758e01d03c`
- **Tip:** `64adc8cc20ad56c277bdd5705ec2fdc5ae413b9f`
- **Commits:**
  - `72be3aa9` feat: read-only tier-3 catalog-match probe surface
  - `64adc8cc` refactor: surface catalog-match probe fetched_at (P6 nit)

## What shipped (read-only, GET only, no writes / no persistence / no OpenAPI/SDK change)

| Layer | File | Role |
|---|---|---|
| connectors/domain | `catalog_match.go` | normalized snapshot types (CatalogHit, BuyBoxSnapshot, DomainPrediction, CatalogMatchSnapshot/Input) |
| connectors/ports | `catalog_match.go` | `CatalogMatchReader` port |
| connectors/adapters/mercado_livre | `catalog_match_reader.go` (+test) | `CapabilityAdapter.ReadCatalogMatch` — search → product/buy-box → domain-discovery |
| integrations/application | `catalog_match_probe.go` (+test) | `ProviderOperationService.ProbeCatalogMatch` — composes two-tier fee quote + flags |
| integrations/application | `provider_operation_service.go` | additive `CatalogMatch` config + struct field |
| integrations/transport | `auth_handler.go` (+test) | `GET .../probes/catalog-match?ean&q&price` |
| composition | `root.go` | facade delegate + inject `CatalogMatch: mercadoLivreCapabilities` |

## Probe response shape (normalized, honest)

`{catalog_hits:[{product_id,name,domain_id,status}], buy_box:{category_id|null,price|null,listing_type}|null,
domain_discovery:[{category_id,category_name,domain_id}] (≤3), fee_quote:{category_id,price_amount,currency_id,
classico:{listing_type_id,percentage_fee|null,fixed_fee|null},premium:{...}}|null,
flags:{category_predita,buy_box_null,no_catalog_hit}, fetched_at}`

- Category resolution: `buy_box.category_id` else `domain_discovery[0].category_id` (`category_predita=true`).
- Price: caller `price` override (>0) else buy-box price.
- `fee_quote` = null when category unresolved OR price ≤ 0 (flagged, never a zero-quote).
- fee tiers via the **existing** `FeeQuoteReader`: classico=`gold_special`, premium=`gold_pro`.
- Every unknown = null + flag, never zero-defaulted (ADR-17).

## API facts (doc-grounded, Context7 /websites/developers_mercadolivre_br_pt_br)

- `GET /products/search?site_id=MLB&product_identifier={EAN}` → results carry `domain_id` (NOT category_id).
- `GET /products/{id}` → `buy_box_winner` (nullable) → `category_id`, `price`, `listing_type_id`; `name`.
- `GET /sites/MLB/domain_discovery/search?q={title}&limit=3` → ranked category predictions (no confidence score).

## Verification ladder

- **L0 build:** `go build ./...` → exit 0 (GOCACHE=.gocache).
- **L0 vet:** `go vet ./internal/modules/integrations/...` → clean.
- **L1 tests:** touched packages green —
  - `connectors/adapters/mercado_livre` ok (3 catalog-match tests: buy-box winner, nil-buy-box fallback, requires ean|query)
  - `integrations/application` ok (3 probe tests: buy-box category quote, predicted category on nil buy-box, fee_quote omitted when unresolved)
  - `integrations/transport` ok (2 route tests: snapshot 200 + snake_case, rejects POST)
- **Governance lane:** `harness:governance -BaseSha 72be3aa9…` from **clean detached worktree** →
  `status=passed`; all `baseline_exception=*` pre-existing (none introduced by this chip).

## P6 dual gate (both non-blocking)

- **Cold Opus subagent (fresh, no context):** VERDICT **PASS-WITH-NITS**. Confirmed: no
  flag value-copy bug, ADR-17 clean, no dead decoded fields, boundaries clean (ml* structs
  confined to adapter), fixtures faithful (domain_id not category_id; nullable buy_box_winner),
  errors propagate + recordOperation marks failed on execErr.
- **Adversarial sonnet:** VERDICT **SHIP**. Independently traced flag path (no drop),
  no nil-panic path, no fixture lies, no boundary leak, listing-type map not swapped,
  no error swallowing; `go build`/`go vet`/tests confirmed green.

### Nits accepted / dispositioned
- **fetched_at discarded** (Opus nit) → **FIXED** in `64adc8cc` (surfaced + asserted).
- **currency_id hardcoded "BRL"** (sonnet) → **NOT a defect**: the existing fee-quote
  adapter itself sets `firstNonEmpty(input.CurrencyID, "BRL")` and does not parse currency
  from the provider response; BRL is a codebase-wide MLB site-invariant. Probe is MLB-scoped.
  Left as-is (Opus concurred it is not an ADR-17 violation).
- **`runtimeCapabilityCodeForOperation` has no `catalog_match_probe` case** → intentional:
  the probe is a composite diagnostic surface, not a runtime capability; the empty-capability
  guard makes `persistCapabilityState` a correct no-op. Operation-run audit trail still records.

## FINDINGS (reported to hub for ratification / CHIP-T2 hand-off)

1. **buy_box_winner.category_id unverified live.** The prompt asserts `/products/{id}`
   exposes `buy_box_winner.category_id`; the repo's own live capture
   (`docs/research/probes/ml_official_durable_price_probe_test.go`) recorded item_id/price/
   currency_id/listing_type_id on the winner but did **not** capture `category_id`. Decoder
   reads `winner.category_id` defensively (empty → falls through to domain-discovery
   prediction, `category_predita=true`). **CHIP-T2 must confirm against a real payload
   whether category_id sits inside buy_box_winner or at product top-level** — if the latter,
   the buy-box category always reads empty and every quote silently becomes predicted.
2. **MIS-004 mission artifacts absent at base SHA.** `mission.md` / `DESIGN-TARIFAS-ML.md`
   for MIS-004-mvp-demo are not present on `master` at the fork point (hub-tree only, not
   merged). Non-blocking here (dispatch prompt carried inline doc-grounded API facts + shape),
   but flagged so the hub knows chip-visible master lacks the mission plan.
3. **`go.work.sum` drifted during bootstrap** (drops pgx/go-spew sums, adds protobuf) from
   `go mod` activity while warming caches. **NOT committed** (dep seams are hub-owned; revert
   is forbidden by doctrine). Left dirty in the chip worktree; governance ran from a clean
   detached worktree so it was unaffected. Hub may want to reconcile go.work.sum on master.
4. **currency site-invariant** (see nit above) — worth making `currency_id` provider-sourced
   or caller-supplied in the real CHIP-T2 implementation if any non-MLB site is ever in scope.

## HARD-RULE compliance

Never pushed. Never booted a server / bound :8080 / read .env. No deps installed. No
reset/revert/stash/clean. Zero ML writes (all GET). Tokens never logged or returned
(slog lines carry action/result/duration only; token resolved server-side).
