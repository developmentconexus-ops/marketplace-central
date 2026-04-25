# Marketplace Catalog UX & Data Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Complete the six-provider marketplace catalog and redesign the marketplace UI around integration provider/installations as the canonical operational source.

**Architecture:** Provider metadata comes from `integrations/providers`; operational connection state comes from `integrations/installations`. Marketplace accounts and pricing policies remain pricing setup data, not the primary connection model. New provider work uses the plugin self-registration framework from ADR-004.

**Tech Stack:** Go 1.25, PostgreSQL via `pgxpool`, OpenAPI, TypeScript, React, Vite, Vitest, Testing Library, `sdk-runtime`.

---

## Parallel Execution Strategy

Use subagents with disjoint write ownership:

- **Agent A - Provider Backend:** owns integration provider packages, registry tests, auth enum, and `composition/root.go`.
- **Agent B - Contract/SDK:** owns OpenAPI and `packages/sdk-runtime`.
- **Agent C - Marketplace UI:** owns `packages/feature-marketplaces`.
- **Agent D - Verification:** owns test/build/browser evidence after A/B/C land.

Run A and B in parallel first. Start C after B defines the final SDK types. Run D last.

---

## File Structure Map

### Create

- `apps/server_core/internal/modules/integrations/adapters/amazon/auth_adapter.go`
- `apps/server_core/internal/modules/integrations/adapters/amazon/auth_adapter_test.go`
- `apps/server_core/internal/modules/integrations/adapters/leroymerlin/auth_adapter.go`
- `apps/server_core/internal/modules/integrations/adapters/leroymerlin/auth_adapter_test.go`
- `apps/server_core/internal/modules/integrations/adapters/madeiramadeira/auth_adapter.go`
- `apps/server_core/internal/modules/integrations/adapters/madeiramadeira/auth_adapter_test.go`
- `packages/feature-marketplaces/src/ProviderCatalogCard.tsx`
- `packages/feature-marketplaces/src/ProviderCatalogPanel.tsx`

### Modify

- `apps/server_core/internal/modules/integrations/domain/lifecycle.go`
- `apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter.go`
- `apps/server_core/internal/modules/integrations/adapters/magalu/auth_adapter.go`
- `apps/server_core/internal/modules/integrations/adapters/shopee/auth_adapter.go`
- `apps/server_core/internal/modules/integrations/adapters/providers/registry_test.go`
- `apps/server_core/internal/composition/root.go`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.ts`
- `packages/sdk-runtime/src/index.test.ts`
- `packages/feature-marketplaces/src/MarketplaceSettingsPage.tsx`
- `packages/feature-marketplaces/src/MarketplaceSettingsPage.test.tsx`
- `packages/feature-marketplaces/src/components/MarketplaceIcon.tsx`
- `packages/feature-marketplaces/src/components/MarketplaceIcon.test.tsx`
- `packages/feature-marketplaces/src/components/StatusBadge.tsx`
- `packages/feature-marketplaces/src/components/StatusBadge.test.tsx`

---

## Task 1: Backend Provider Catalog Coverage

**Files:**
- Create: `apps/server_core/internal/modules/integrations/adapters/amazon/auth_adapter.go`
- Create: `apps/server_core/internal/modules/integrations/adapters/leroymerlin/auth_adapter.go`
- Create: `apps/server_core/internal/modules/integrations/adapters/madeiramadeira/auth_adapter.go`
- Modify: `apps/server_core/internal/modules/integrations/domain/lifecycle.go`
- Modify: `apps/server_core/internal/modules/integrations/adapters/providers/registry_test.go`
- Modify: `apps/server_core/internal/composition/root.go`

- [ ] **Step 1: Write failing registry coverage**

Update `registry_test.go` so `TestRegistryIncludesCoreProviders` expects exactly six providers:

```go
wantCodes := map[string]domain.AuthStrategy{
	"mercado_livre":   domain.AuthStrategyOAuth2,
	"magalu":          domain.AuthStrategyOAuth2,
	"shopee":          domain.AuthStrategyAPIKey,
	"amazon":          domain.AuthStrategyLWA,
	"leroy_merlin":    domain.AuthStrategyAPIKey,
	"madeira_madeira": domain.AuthStrategyToken,
}
```

Expected now: compile fails because `AuthStrategyLWA` and new packages do not exist.

- [ ] **Step 2: Add first-class LWA auth strategy**

Add to `apps/server_core/internal/modules/integrations/domain/lifecycle.go`:

```go
AuthStrategyLWA AuthStrategy = "lwa"
```

- [ ] **Step 3: Add provider self-registrations**

Each new package must call `integrationsproviders.RegisterDefinition(...)` in `init()`.

Required provider metadata:

```go
Metadata: map[string]any{
	"country": "BR",
	"rollout_stage": "v1",          // amazon
	"execution_mode": "available",  // amazon, leroy, madeira
	"fee_source": "seed",
	"baseline_commission_percent": 0.12,
	"baseline_fixed_fee_amount": 0.0,
	"credential_schema": []map[string]any{
		{"key": "seller_id", "label": "Seller ID", "secret": false},
	},
}
```

Use these provider specifics:

- `amazon`: `AuthStrategyLWA`, `InstallModeInteractive`, baseline commission `0.12`.
- `leroy_merlin`: `AuthStrategyAPIKey`, `InstallModeManual`, baseline commission `0.18`, credential fields `api_key`, `shop_id`.
- `madeira_madeira`: `AuthStrategyToken`, `InstallModeManual`, baseline commission `0.15`, credential fields `api_token`.

- [ ] **Step 4: Normalize existing definitions**

Keep ML/Magalu active and stable. Mark Shopee visible but blocked by setting:

```go
Metadata: map[string]any{
	"country": "BR",
	"rollout_stage": "blocked",
	"execution_mode": "blocked",
	"unavailable_reason": "Shopee access is blocked until partner credentials are approved.",
	"fee_source": "seed",
}
```

Do not expose a normal connect CTA for Shopee in the UI.

- [ ] **Step 5: Register packages in composition root**

Add side-effect imports:

```go
_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/amazon"
_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/leroymerlin"
_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/madeiramadeira"
```

- [ ] **Step 6: Verify backend slice**

Run:

```powershell
$env:GOCACHE=".gocache"; go test ./internal/modules/integrations/... ./internal/composition/...
```

Expected: PASS.

- [ ] **Step 7: Commit**

```powershell
git add apps/server_core/internal/modules/integrations apps/server_core/internal/composition/root.go
git commit -m "feat(integrations): complete marketplace provider catalog"
```

---

## Task 2: Contract and SDK Alignment

**Files:**
- Modify: `contracts/api/marketplace-central.openapi.yaml`
- Modify: `packages/sdk-runtime/src/index.ts`
- Modify: `packages/sdk-runtime/src/index.test.ts`

- [ ] **Step 1: Write SDK tests for catalog metadata**

Add tests proving `listIntegrationProviders()` accepts:

- `auth_strategy: "lwa"`
- `metadata.execution_mode`
- `metadata.rollout_stage`
- `metadata.baseline_commission_percent`
- `metadata.credential_schema`

Use an Amazon fixture and a blocked Shopee fixture.

- [ ] **Step 2: Update OpenAPI**

In `IntegrationProviderDefinition.auth_strategy`, include:

```yaml
enum: [oauth2, lwa, api_key, token, none, unknown]
```

Keep `metadata` as object, but document these stable keys in schema description:

```yaml
metadata:
  type: object
  description: Provider UI metadata. Stable keys include country, rollout_stage, execution_mode, unavailable_reason, fee_source, baseline_commission_percent, baseline_fixed_fee_amount, credential_schema, docs_url.
  additionalProperties: true
```

- [ ] **Step 3: Update SDK types**

Add `lwa` to `IntegrationProviderDefinition.auth_strategy`.

Add a typed helper shape:

```ts
export interface IntegrationProviderMetadata {
  country?: string;
  rollout_stage?: "v1" | "wave_2" | "blocked";
  execution_mode?: "available" | "blocked" | "planned";
  unavailable_reason?: string;
  fee_source?: "api_sync" | "seed" | "manual";
  baseline_commission_percent?: number;
  baseline_fixed_fee_amount?: number;
  credential_schema?: CredentialField[];
  docs_url?: string;
}
```

Set `metadata?: IntegrationProviderMetadata & Record<string, unknown>` on `IntegrationProviderDefinition`.

- [ ] **Step 4: Verify SDK slice**

Run:

```powershell
npm exec --workspace @marketplace-central/web vitest run packages/sdk-runtime/src/index.test.ts
```

Expected: PASS.

- [ ] **Step 5: Commit**

```powershell
git add contracts/api/marketplace-central.openapi.yaml packages/sdk-runtime/src/index.ts packages/sdk-runtime/src/index.test.ts
git commit -m "feat(sdk-runtime): expose integration provider catalog metadata"
```

---

## Task 3: Marketplace Catalog UI Redesign

**Files:**
- Create: `packages/feature-marketplaces/src/ProviderCatalogCard.tsx`
- Create: `packages/feature-marketplaces/src/ProviderCatalogPanel.tsx`
- Modify: `packages/feature-marketplaces/src/MarketplaceSettingsPage.tsx`
- Modify: `packages/feature-marketplaces/src/MarketplaceSettingsPage.test.tsx`
- Modify: `packages/feature-marketplaces/src/components/MarketplaceIcon.tsx`
- Modify: `packages/feature-marketplaces/src/components/StatusBadge.tsx`

- [ ] **Step 1: Write failing UI tests**

Add tests to `MarketplaceSettingsPage.test.tsx`:

- renders all six providers from `listIntegrationProviders()`
- maps connected installation status onto the provider card
- disables the Shopee action when `metadata.execution_mode === "blocked"`
- shows loading, error, and empty provider states
- never calls direct backend APIs outside the injected SDK client

- [ ] **Step 2: Expand client interface**

`MarketplaceSettingsPage` must use these SDK methods:

```ts
listIntegrationProviders()
listIntegrationInstallations()
createIntegrationInstallation(req)
startIntegrationAuthorization(installationId)
startIntegrationReauthorization(installationId)
submitIntegrationCredentials(installationId, req)
disconnectIntegrationInstallation(installationId)
startIntegrationFeeSync(installationId)
```

Keep `listMarketplaceAccounts()` and `listMarketplacePolicies()` only for pricing policy summaries.

- [ ] **Step 3: Replace account-first grid with provider catalog grid**

Render one card per provider. Card fields:

- provider name
- auth strategy
- rollout/execution state
- connected installation status when present
- baseline commission/fixed fee
- capability summary from `declared_capabilities`
- primary action

Primary action rules:

- blocked provider: disabled button with unavailable reason
- no installation: create draft installation, then start auth or open manual credentials
- pending OAuth/LWA: authorize
- requires reauth: reauthorize
- connected/degraded: open details panel

- [ ] **Step 4: Implement provider panel**

Panel shows:

- provider metadata
- installation status
- credential fields for manual providers
- fee sync action for connected providers
- pricing setup summary if a marketplace policy exists

Do not perform pricing or commission calculations in React.

- [ ] **Step 5: Update icons and badges**

Add icon/color handling for:

```ts
amazon
leroy_merlin
madeira_madeira
```

Add badge handling for:

```ts
connected
disconnected
degraded
requires_reauth
blocked
planned
available
```

- [ ] **Step 6: Verify frontend slice**

Run:

```powershell
npm exec --workspace @marketplace-central/web vitest run packages/feature-marketplaces/src/MarketplaceSettingsPage.test.tsx packages/feature-marketplaces/src/components/MarketplaceIcon.test.tsx packages/feature-marketplaces/src/components/StatusBadge.test.tsx
```

Expected: PASS.

- [ ] **Step 7: Commit**

```powershell
git add packages/feature-marketplaces
git commit -m "feat(marketplaces): render integration-backed provider catalog"
```

---

## Task 4: Cross-Slice Verification and Evidence

**Files:**
- Create: `docs/superpowers/evidence/2026-04-25-marketplace-catalog-ux-data-foundation.md`
- Modify: `.brain/session-log.md`
- Modify: `.brain/system-pulse.md`
- Modify: `.brain/roadmap.json`
- Modify (if changed): `wiki/modules/integrations.md`
- Modify (if changed): `wiki/modules/marketplaces.md`

- [ ] **Step 1: Run backend tests**

```powershell
$env:GOCACHE=".gocache"; go test ./...
```

Expected: PASS.

- [ ] **Step 2: Run frontend tests**

```powershell
npm test -- --runInBand
```

Expected: PASS.

- [ ] **Step 3: Run frontend build**

```powershell
npm run build
```

Expected: PASS.

- [ ] **Step 4: Capture browser evidence**

Start local apps, open `/marketplaces`, and capture evidence that:

- all six providers render
- ML/Magalu show OAuth-capable states
- Amazon shows `lwa`
- Leroy shows API-key/manual setup
- Madeira shows token/manual setup
- Shopee is visible but disabled/blocked
- connected/degraded/requires-reauth statuses display correctly when fixtures or local data exist

- [ ] **Step 5: Write evidence ledger**

Create `docs/superpowers/evidence/2026-04-25-marketplace-catalog-ux-data-foundation.md` with:

```markdown
# Marketplace Catalog UX & Data Foundation Evidence

## Commands
- `$env:GOCACHE=".gocache"; go test ./...` - PASS
- `npm test -- --runInBand` - PASS
- `npm run build` - PASS

## Browser Evidence
- `/marketplaces` renders all six providers.
- Shopee is visible and disabled with a blocked reason.
- Amazon displays `lwa`.
- Manual providers expose credential setup affordances.

## Known Constraints
- Amazon LWA live token exchange is not required for this phase.
- Baseline fees are catalog metadata; full category schedules remain provider-specific sync work.
```

- [ ] **Step 6: Commit**

```powershell
git add docs/superpowers/evidence/2026-04-25-marketplace-catalog-ux-data-foundation.md
git commit -m "test(marketplaces): capture provider catalog validation evidence"
```

- [ ] **Step 7: Update Nexus and Wiki**

Run Nexus checkpoint after implementation/testing:

```text
Use skill: nexus-checkpoint
```

Update wiki only if contract, provider capability semantics, or operational flows changed from current docs:

```text
wiki/modules/integrations.md
wiki/modules/marketplaces.md
```

- [ ] **Step 8: Commit knowledge sync**

```powershell
git add .brain/session-log.md .brain/system-pulse.md .brain/roadmap.json wiki/modules/integrations.md wiki/modules/marketplaces.md
git commit -m "docs(nexus): checkpoint roadmap and update module docs after catalog rollout"
```

---

## Done Gate

This plan is complete only when:

- `/integrations/providers` returns six marketplace providers.
- `IntegrationProviderDefinition.auth_strategy` supports `lwa`.
- `/marketplaces` renders the dense integration-backed provider catalog.
- Shopee is visible but non-connectable.
- Amazon, Leroy Merlin, and MadeiraMadeira are visible with honest auth/setup metadata.
- All frontend data access goes through `sdk-runtime`.
- Go tests, frontend tests, and frontend build pass.
- Evidence file exists with command results and browser validation notes.
- Nexus brain is checkpointed to reflect completed tasks and current phase.
- Wiki modules are updated when behavior/contracts changed.

## Self-Review Notes

Spec coverage:

- `T-031` is covered by Task 1.
- `T-033` is covered by Task 2.
- `T-032` is covered by Task 3.
- `T-034` is covered by Task 4.

Type consistency:

- Provider codes use `mercado_livre`, `magalu`, `shopee`, `amazon`, `leroy_merlin`, `madeira_madeira`.
- Auth strategies use `oauth2`, `lwa`, `api_key`, `token`, `none`, `unknown`.
- UI reads `IntegrationProviderDefinition.metadata`, not `MarketplaceDefinition.metadata`, for the provider catalog.
