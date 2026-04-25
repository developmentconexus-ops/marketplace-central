# Integration Catalog Plugin Framework Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make integrations behave like catalog plugins so future providers can be added with provider-owned packages instead of deep central wiring.

**Architecture:** Keep Integrations as the catalog owner: provider definitions are the available catalog, installations are tenant-specific connections, and capabilities gate behavior. Provider packages self-register definitions/auth factories; connector packages self-register optional fee syncers. `root.go` consumes registries instead of hardcoded lists.

**Tech Stack:** Go 1.25.x, pgxpool, modular monolith ports/adapters, Nexus Brain, Markdown wiki.

---

## Parallel Execution Map

- **Agent A - Registry Core:** Task 1. Owns `integrations/adapters/providers`.
- **Agent B - Provider Plugins:** Task 2. Owns ML/Magalu/Shopee integration adapter registration.
- **Agent C - Fee Sync Plugins:** Task 3. Owns connector fee syncer registration.
- **Agent D - Docs/Nexus:** Task 5 can run after Task 1 interfaces are known.
- **Main integrator:** Task 4 and Task 6 after Tasks 1-3 land.

Avoid overlapping writes: only one agent edits `root.go`, `roadmap.json`, or ADR index.

---

### Task 1: Add Registration-Backed Provider Registry

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/adapters/providers/registry.go`
- Test: `apps/server_core/internal/modules/integrations/adapters/providers/registry_test.go`

- [ ] **Step 1: Write failing registry tests**

Use package `providers` (not `providers_test`) so tests can exercise unexported build helpers without exposing reset hooks to production code.

```go
func TestBuildDefinitionsRejectsDuplicateProviderCode(t *testing.T) {
	t.Parallel()

	defs, err := buildDefinitions([]domain.ProviderDefinition{
		{ProviderCode: "mercado_livre", TenantID: "system"},
		{ProviderCode: "mercado_livre", TenantID: "system"},
	})

	if err == nil {
		t.Fatal("buildDefinitions() error = nil, want duplicate error")
	}
	if defs != nil {
		t.Fatalf("buildDefinitions() defs = %#v, want nil", defs)
	}
	if !strings.Contains(err.Error(), "INTEGRATIONS_PROVIDER_DUPLICATE") {
		t.Fatalf("error = %v, want INTEGRATIONS_PROVIDER_DUPLICATE", err)
	}
}
```

Run:

```powershell
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/adapters/providers
```

Expected: FAIL because `buildDefinitions` does not exist yet.

- [ ] **Step 2: Implement registry entries and validation**

Change `registry.go` to store registered entries instead of an inline hardcoded array. Keep public `NewRegistry() *Registry` for minimal call-site churn, and add `Validate() error` for root startup.

```go
type AuthFactory func() application.MarketplaceAuthAdapter
type FeeSyncerFactory func() marketplacesports.FeeScheduleSyncer

var (
	definitionEntries []domain.ProviderDefinition
	authFactories     []AuthFactory
	feeSyncerFactories []FeeSyncerFactory
)

func RegisterDefinition(def domain.ProviderDefinition) {
	definitionEntries = append(definitionEntries, cloneProviderDefinition(def))
}

func RegisterAuthFactory(factory AuthFactory) {
	if factory != nil {
		authFactories = append(authFactories, factory)
	}
}

func RegisterFeeSyncerFactory(factory FeeSyncerFactory) {
	if factory != nil {
		feeSyncerFactories = append(feeSyncerFactories, factory)
	}
}
```

Add helpers:

```go
func buildDefinitions(entries []domain.ProviderDefinition) ([]domain.ProviderDefinition, error) {
	seen := map[string]struct{}{}
	out := make([]domain.ProviderDefinition, 0, len(entries))
	for _, entry := range entries {
		code := strings.TrimSpace(entry.ProviderCode)
		if code == "" {
			return nil, errors.New("INTEGRATIONS_PROVIDER_CODE_REQUIRED")
		}
		if _, exists := seen[code]; exists {
			return nil, fmt.Errorf("INTEGRATIONS_PROVIDER_DUPLICATE: provider_code=%s", code)
		}
		seen[code] = struct{}{}
		out = append(out, cloneProviderDefinition(entry))
	}
	return out, nil
}
```

- [ ] **Step 3: Add construction methods**

Add:

```go
func NewRegistry() *Registry {
	defs, err := buildDefinitions(definitionEntries)
	return &Registry{definitions: defs, err: err}
}

func (r *Registry) Validate() error {
	if r == nil {
		return errors.New("INTEGRATIONS_PROVIDER_REGISTRY_NIL")
	}
	return r.err
}

func BuildAuthAdapters() []application.MarketplaceAuthAdapter {
	out := make([]application.MarketplaceAuthAdapter, 0, len(authFactories))
	for _, factory := range authFactories {
		if adapter := factory(); adapter != nil {
			out = append(out, adapter)
		}
	}
	return out
}

func BuildFeeSyncers() []marketplacesports.FeeScheduleSyncer {
	out := make([]marketplacesports.FeeScheduleSyncer, 0, len(feeSyncerFactories))
	for _, factory := range feeSyncerFactories {
		if syncer := factory(); syncer != nil && strings.TrimSpace(syncer.MarketplaceCode()) != "" {
			out = append(out, syncer)
		}
	}
	return out
}
```

- [ ] **Step 4: Run provider registry tests**

Run:

```powershell
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/adapters/providers
```

Expected: PASS.

---

### Task 2: Register Existing Provider Definitions and Auth Factories

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter.go`
- Modify: `apps/server_core/internal/modules/integrations/adapters/magalu/auth_adapter.go`
- Modify: `apps/server_core/internal/modules/integrations/adapters/shopee/auth_adapter.go`
- Test: existing adapter tests in each package

- [ ] **Step 1: Add provider registration to Mercado Livre**

In `mercadolivre/auth_adapter.go`, add imports for `os`, `providers`, and register:

```go
func init() {
	providers.RegisterDefinition(domain.ProviderDefinition{
		ProviderCode: "mercado_livre",
		TenantID:     "system",
		Family:       domain.IntegrationFamilyMarketplace,
		DisplayName:  "Mercado Livre",
		AuthStrategy: domain.AuthStrategyOAuth2,
		InstallMode:  domain.InstallModeInteractive,
		Metadata: map[string]any{
			"country":       "BR",
			"release_stage": "stable",
			"fee_source":    "api_sync",
		},
		DeclaredCapabilities: []string{
			"catalog_publish", "pricing_fee_sync", "inventory_sync", "order_read",
			"message_read", "message_reply", "shipment_tracking", "webhook_receive",
		},
		IsActive: true,
	})
	providers.RegisterAuthFactory(func() application.MarketplaceAuthAdapter {
		return NewAdapter(Config{
			ClientID:     strings.TrimSpace(os.Getenv("MPC_PROVIDER_MERCADOLIVRE_CLIENT_ID")),
			ClientSecret: strings.TrimSpace(os.Getenv("MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET")),
			AuthorizeURL: "https://auth.mercadolivre.com.br/authorization",
			TokenURL:     "https://api.mercadolibre.com/oauth/token",
		})
	})
}
```

- [ ] **Step 2: Add provider registration to Magalu**

Use the same pattern with:

```go
ProviderCode: "magalu"
DisplayName: "Magalu"
AuthStrategy: domain.AuthStrategyOAuth2
InstallMode: domain.InstallModeInteractive
ClientID env: "MPC_PROVIDER_MAGALU_CLIENT_ID"
ClientSecret env: "MPC_PROVIDER_MAGALU_CLIENT_SECRET"
AuthorizeURL: "https://id.magalu.com/login"
TokenURL: "https://id.magalu.com/oauth/token"
Metadata fee_source: "api_sync"
```

- [ ] **Step 3: Add provider registration to Shopee**

Use:

```go
ProviderCode: "shopee"
DisplayName: "Shopee"
AuthStrategy: domain.AuthStrategyAPIKey
InstallMode: domain.InstallModeManual
Metadata fee_source: "seed"
providers.RegisterAuthFactory(func() application.MarketplaceAuthAdapter {
	return NewAdapter(Config{})
})
```

- [ ] **Step 4: Run adapter package tests**

Run:

```powershell
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/adapters/mercadolivre ./internal/modules/integrations/adapters/magalu ./internal/modules/integrations/adapters/shopee
```

Expected: PASS.

---

### Task 3: Register Existing Fee Syncers

**Files:**
- Modify: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/*fee*`
- Modify: `apps/server_core/internal/modules/connectors/adapters/magalu/*fee*`
- Modify: `apps/server_core/internal/modules/connectors/adapters/shopee/*fee*`
- Test: `apps/server_core/internal/modules/integrations/adapters/providers/registry_test.go`

- [ ] **Step 1: Add fee syncer registration in each connector package**

In each connector fee syncer file, import `integrations/adapters/providers` and register:

```go
func init() {
	providers.RegisterFeeSyncerFactory(func() marketplacesports.FeeScheduleSyncer {
		return NewFeeSyncer()
	})
}
```

Use the existing marketplace ports import alias:

```go
marketplacesports "marketplace-central/apps/server_core/internal/modules/marketplaces/ports"
```

- [ ] **Step 2: Add registry test for fee syncers**

In `registry_test.go`, assert registered syncers include:

```go
wantCodes := map[string]bool{
	"mercado_livre": true,
	"magalu":        true,
	"shopee":        true,
}
for _, syncer := range BuildFeeSyncers() {
	delete(wantCodes, syncer.MarketplaceCode())
}
if len(wantCodes) != 0 {
	t.Fatalf("missing fee syncers: %v", wantCodes)
}
```

- [ ] **Step 3: Run registry and connector tests**

Run:

```powershell
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/adapters/providers ./internal/modules/connectors/adapters/...
```

Expected: PASS.

---

### Task 4: Simplify Composition Root

**Files:**
- Modify: `apps/server_core/internal/composition/root.go`

- [ ] **Step 1: Replace named adapter imports with side-effect imports**

Keep provider registration packages imported for init registration:

```go
_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/magalu"
_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/shopee"
_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/magalu"
_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/mercadolivre"
_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/shopee"
```

Remove now-unused `os`, `strings`, named `connml`, `connmagalu`, `connshopee`, `integrationsml`, `integrationsmagalu`, and `integrationsshopee` imports.

- [ ] **Step 2: Validate provider registry before seeding**

Replace provider setup with:

```go
providerRegistry := integrationsproviders.NewRegistry()
if err := providerRegistry.Validate(); err != nil {
	return nil, fmt.Errorf("integration provider registry: %w", err)
}
```

- [ ] **Step 3: Use registered syncers and auth adapters**

Replace both manual fee syncer lists with:

```go
registeredFeeSyncers := integrationsproviders.BuildFeeSyncers()
feeSyncExecutor := integrationsfeesync.NewMarketplaceExecutor(feeRepo, registeredFeeSyncers)
connectorsFeeSyncSvc := connectorsapp.NewFeeSyncService(feeRepo, registeredFeeSyncers...)
```

Replace manual auth adapter construction with:

```go
Adapters: integrationsproviders.BuildAuthAdapters(),
```

- [ ] **Step 4: Run root-level backend tests**

Run:

```powershell
$env:GOCACHE='.gocache'; go test ./internal/composition ./...
```

Expected: PASS.

---

### Task 5: Update Wiki and Nexus

**Files:**
- Modify: `wiki/architecture/marketplace-extensibility.md`
- Modify: `wiki/modules/integrations.md`
- Modify: `.brain/roadmap.json`
- Create: `.brain/decisions/004-integration-catalog-plugin-framework.md`
- Modify: `.brain/decisions/_index.md`
- Modify: `.brain/system-pulse.md`

- [ ] **Step 1: Update Wiki wording**

Document:

```text
ProviderDefinition = available catalog item
Installation = tenant-connected instance
Capabilities = feature gates
Provider package = owns its definition and auth factory
Connector package = owns optional operational syncers
```

- [ ] **Step 2: Add Nexus ADR**

Create ADR-004 with:

```markdown
# ADR-004: Integration Catalog Plugin Framework

**Date:** 2026-04-25
**Status:** accepted

## Context
Marketplace definitions already use plugin-style registration, but integrations still required central edits in provider registry and composition root. This made each new provider harder to add and increased root wiring churn.

## Decision
Integration providers self-register catalog definitions, auth adapter factories, and optional fee syncers. Composition root consumes registries instead of owning provider-specific construction.

## Rationale
This matches the product model of installable integrations while keeping provider-specific config close to provider code. Registration reduces repeated central edits without introducing runtime plugin loading.

## Consequences
Future providers are added through provider-owned packages plus side-effect imports where Go requires registration. Duplicate provider codes are treated as startup configuration errors.

## Alternatives Considered
- Keep root wiring: rejected because it scales poorly with each provider.
- Dynamic runtime plugins: rejected as unnecessary complexity for the current modular monolith.
```

- [ ] **Step 3: Add roadmap task**

Add one planned task under Phase 7:

```json
{
  "id": "T-030",
  "title": "Integration catalog plugin framework",
  "description": "Refactor provider definitions, auth adapters, and fee syncers into self-registration patterns so new integrations can be added with minimal central wiring.",
  "status": "planned",
  "priority": "high",
  "depends_on": ["T-027"],
  "started_at": null,
  "completed_at": null,
  "notes": null,
  "subtasks": [
    "Add provider definition and auth factory registration",
    "Register existing ML/Magalu/Shopee providers",
    "Register existing fee syncers",
    "Update composition root and documentation"
  ]
}
```

- [ ] **Step 4: Verify docs do not contain secrets**

Run:

```powershell
Select-String -Path wiki/**/*.md,.brain/**/*.md -Pattern "CLIENT_SECRET=|APP_TOKEN=|DATABASE_URL=postgres"
```

Expected: no matches.

---

### Task 6: Final Verification and Commit

**Files:**
- All files changed by Tasks 1-5

- [ ] **Step 1: Full backend test**

Run:

```powershell
$env:GOCACHE='.gocache'; go test ./...
```

Expected: PASS.

- [ ] **Step 2: Inspect root wiring**

Run:

```powershell
Select-String -Path apps/server_core/internal/composition/root.go -Pattern "NewFeeSyncer|MPC_PROVIDER_|NewAdapter\\("
```

Expected: no manual ML/Magalu/Shopee auth or fee sync construction remains in `root.go`.

- [ ] **Step 3: Commit**

Run:

```powershell
git add apps/server_core/internal/modules/integrations apps/server_core/internal/modules/connectors apps/server_core/internal/composition/root.go wiki .brain
git commit -m "refactor(integrations): add catalog plugin registration"
```

Expected: commit succeeds.
