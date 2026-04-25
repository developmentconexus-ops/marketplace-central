# Marketplace Catalog UX & Data Foundation Evidence

Date: 2026-04-25
Scope: `T-031` to `T-034`

## Commands

- `$env:GOCACHE='C:\Users\leandro.theodoro.MN-NTB-LEANDROT\Documents\marketplace-central\.gocache'; go test ./...` (workdir `apps/server_core`) - PASS
- `cmd /c npm test -- --runInBand` (workdir repo root) - PASS (`17` files, `147` tests)
- `cmd /c npm run build` (workdir repo root) - PASS
- `cmd /c npm exec --workspace @marketplace-central/web vitest run packages/sdk-runtime/src/index.test.ts packages/feature-marketplaces/src/MarketplaceSettingsPage.test.tsx packages/feature-marketplaces/src/components/MarketplaceIcon.test.tsx packages/feature-marketplaces/src/components/StatusBadge.test.tsx` - PASS (`4` files, `37` tests)
- `C:\Program Files\PostgreSQL\16\bin\psql.exe -d <MC_DATABASE_URL-without-query> -f apps/server_core/migrations/0019_integrations_provider_auth_strategy_lwa.sql` - PASS
- Browser Use validation pass on `http://localhost:5173/marketplaces` (in-app browser) - PASS

## Backend + Contract Evidence

- `/integrations/providers` catalog now has six providers registered through provider-owned self-registration packages:
  - `mercado_livre`, `magalu`, `shopee`, `amazon`, `leroy_merlin`, `madeira_madeira`
- `integrations` domain auth strategy now supports `lwa` for Amazon.
- OpenAPI and `sdk-runtime` now align on provider metadata used by catalog UX (`execution_mode`, `rollout_stage`, baseline fee fields, credential schema).

## Frontend Evidence

- `MarketplaceSettingsPage` now renders provider-first catalog cards from integrations provider/installations data.
- Shopee appears in catalog with blocked/unavailable state.
- Interactive providers trigger authorization flow via injected navigation callback.
- Manual providers expose credential submission path.
- Status badges and provider icons include new provider/state combinations required by this rollout.

## Browser Validation Evidence

- Provider grid shows all six providers: `Amazon Brasil`, `Leroy Merlin`, `Madeira Madeira`, `Magalu`, `Mercado Livre`, `Shopee`.
- Auth strategy labels verified in UI:
  - Amazon shows `lwa`.
  - Manual providers show non-interactive auth (`api_key`/`token`) with credential workflows.
- Status badges verified from live page:
  - `available`: 3 providers
  - `connected`: 2 providers
  - `blocked`: 1 provider (`Shopee`)
- Panel/action validation via Browser Use:
  - Amazon panel: `Create installation` enabled, `Run fee sync` disabled (no installation yet).
  - Shopee panel: `Unavailable` action disabled and blocked status shown.
  - Madeira Madeira panel: `Install mode: manual`, credential input visible, `Submit credentials` disabled until installation/credentials state allows submission.
- Screenshots captured:
  - [01-marketplaces-grid.png](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/docs/superpowers/evidence/screenshots/2026-04-25-marketplace-catalog/01-marketplaces-grid.png)
  - [02-amazon-panel-lwa.png](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/docs/superpowers/evidence/screenshots/2026-04-25-marketplace-catalog/02-amazon-panel-lwa.png)
  - [03-shopee-blocked-panel.png](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/docs/superpowers/evidence/screenshots/2026-04-25-marketplace-catalog/03-shopee-blocked-panel.png)
  - [04-madeira-manual-credentials-panel.png](/C:/Users/leandro.theodoro.MN-NTB-LEANDROT/Documents/marketplace-central/docs/superpowers/evidence/screenshots/2026-04-25-marketplace-catalog/04-madeira-manual-credentials-panel.png)

## Known Constraints

- `MarketplaceSettingsPage` loading test still emits React `act(...)` warnings in stderr, but suite passes and behavior assertions are stable.
- Frontend test/build commands require elevated execution in this environment due sandbox `spawn EPERM` with Vite/esbuild.
