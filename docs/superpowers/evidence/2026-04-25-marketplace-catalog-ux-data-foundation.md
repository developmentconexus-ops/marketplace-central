# Marketplace Catalog UX & Data Foundation Evidence

Date: 2026-04-25
Scope: `T-031` to `T-034`

## Commands

- `$env:GOCACHE='C:\Users\leandro.theodoro.MN-NTB-LEANDROT\Documents\marketplace-central\.gocache'; go test ./...` (workdir `apps/server_core`) - PASS
- `cmd /c npm test -- --runInBand` (workdir repo root) - PASS (`17` files, `147` tests)
- `cmd /c npm run build` (workdir repo root) - PASS
- `cmd /c npm exec --workspace @marketplace-central/web vitest run packages/sdk-runtime/src/index.test.ts packages/feature-marketplaces/src/MarketplaceSettingsPage.test.tsx packages/feature-marketplaces/src/components/MarketplaceIcon.test.tsx packages/feature-marketplaces/src/components/StatusBadge.test.tsx` - PASS (`4` files, `37` tests)

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

## Known Constraints

- Browser screenshot validation for this phase was not captured in this session.
- `MarketplaceSettingsPage` loading test still emits React `act(...)` warnings in stderr, but suite passes and behavior assertions are stable.
- Frontend test/build commands require elevated execution in this environment due sandbox `spawn EPERM` with Vite/esbuild.
