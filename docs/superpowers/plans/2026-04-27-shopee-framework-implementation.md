# Shopee Framework Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement Shopee as a real interactive provider in MPC using signed partner auth flow, production-grade token lifecycle, and aligned contract/UI/docs.

**Architecture:** Keep Shopee integration inside existing `integrations` + `connectors` boundaries. Implement signed partner auth in the Shopee auth adapter, keep all secrets server-side, and expose only stable provider metadata/status through current integration routes. Preserve existing module seams; add only what is required for Shopee auth correctness and operability.

**Tech Stack:** Go (`apps/server_core`), OpenAPI (`contracts/api/marketplace-central.openapi.yaml`), TypeScript SDK/UI (`packages/sdk-runtime`, `packages/feature-marketplaces`), Wiki docs (`wiki/framework/vendors/shopee`).

---

## Assumptions

- We will use a Shopee-specific auth strategy string (`shopee_partner`) instead of overloading `oauth2` or `api_key`.
- First rollout scope is auth + installation lifecycle + provider catalog correctness; deep operational connectors (orders/catalog/logistics execution) stay for subsequent tasks.
- Environment variables will be read in adapter package (consistent with existing provider adapters).

## Parallel Subagent Plan (disjoint ownership)

1. **Worker A (Backend Auth Core)**  
   Ownership: `apps/server_core/internal/modules/integrations/adapters/shopee/**`  
   Tasks: 2, 5

2. **Worker B (Contract/SDK/UI Surface)**  
   Ownership: `contracts/api/marketplace-central.openapi.yaml`, `packages/sdk-runtime/src/**`, `packages/feature-marketplaces/src/**`, `apps/server_core/internal/modules/integrations/transport/**`, `apps/server_core/internal/modules/integrations/domain/lifecycle.go`  
   Tasks: 1, 3, 4

3. **Worker C (Wiki Keeper, already created)**  
   Ownership: `wiki/framework/vendors/shopee/**`, `wiki/modules/integrations.md`  
   Tasks: 6

---

### Task 1: Add Shopee Auth Strategy to Domain + Contract + SDK

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/domain/lifecycle.go`
- Modify: `contracts/api/marketplace-central.openapi.yaml`
- Modify: `packages/sdk-runtime/src/index.ts`
- Test: `apps/server_core/internal/modules/integrations/domain/lifecycle_test.go`
- Test: `packages/sdk-runtime/src/index.test.ts`

- [ ] **Step 1: Write failing tests**

```go
func TestAuthStrategyShopeePartnerConstant(t *testing.T) {
	if AuthStrategy("shopee_partner") != AuthStrategyShopeePartner {
		t.Fatalf("missing shopee auth strategy constant")
	}
}
```

```ts
type _ShopeeAuthStrategyCompiles = Extract<
  IntegrationProviderDefinition["auth_strategy"],
  "shopee_partner"
>;
```

- [ ] **Step 2: Run tests to verify failure**

Run:
```bash
cd apps/server_core
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/domain -run ShopeePartner -count=1
```

Run:
```bash
npm run -w packages/sdk-runtime test -- --run
```

Expected: both fail because `shopee_partner` is not in current enums/types.

- [ ] **Step 3: Implement minimal enum/type changes**

```go
const (
    AuthStrategyShopeePartner AuthStrategy = "shopee_partner"
)
```

```yaml
auth_strategy:
  type: string
  enum: [oauth2, lwa, api_key, token, none, shopee_partner, unknown]
```

```ts
auth_strategy: "oauth2" | "lwa" | "api_key" | "token" | "none" | "shopee_partner" | "unknown";
```

- [ ] **Step 4: Run tests/build to verify pass**

Run:
```bash
cd apps/server_core
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/domain
```

Run:
```bash
npm run -w packages/sdk-runtime build
npm run -w packages/sdk-runtime test -- --run
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/integrations/domain/lifecycle.go apps/server_core/internal/modules/integrations/domain/lifecycle_test.go contracts/api/marketplace-central.openapi.yaml packages/sdk-runtime/src/index.ts packages/sdk-runtime/src/index.test.ts
git commit -m "feat(integrations): add shopee_partner auth strategy contract"
```

---

### Task 2: Implement Shopee Signed Partner Auth Adapter

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/adapters/shopee/auth_adapter.go`
- Create: `apps/server_core/internal/modules/integrations/adapters/shopee/signer.go`
- Create: `apps/server_core/internal/modules/integrations/adapters/shopee/http_types.go`
- Test: `apps/server_core/internal/modules/integrations/adapters/shopee/auth_adapter_test.go`
- Test: `apps/server_core/internal/modules/integrations/adapters/shopee/signer_test.go`

- [ ] **Step 1: Write failing tests for auth URL/sign/token flow**

```go
func TestStartAuthorizeBuildsSignedPartnerURL(t *testing.T) { /* expects /api/v2/shop/auth_partner with partner_id,timestamp,sign,redirect */ }
func TestExchangeCallbackRequiresShopIDAndReturnsTokens(t *testing.T) { /* expects access+refresh+shop_id */ }
func TestRefreshCallsAccessTokenGetEndpoint(t *testing.T) { /* expects /api/v2/auth/access_token/get */ }
```

- [ ] **Step 2: Run tests to verify failure**

Run:
```bash
cd apps/server_core
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/adapters/shopee -count=1
```

Expected: FAIL on unsupported behavior / missing signer.

- [ ] **Step 3: Implement minimal signed auth behavior**

```go
// Config additions
PartnerID string
PartnerKey string
BaseURL string
Now func() time.Time
```

```go
func (a *Adapter) AuthStrategy() domain.AuthStrategy { return domain.AuthStrategyShopeePartner }

func (a *Adapter) StartAuthorize(...) (application.AuthorizeStart, error) {
  // build /api/v2/shop/auth_partner signed URL
}

func (a *Adapter) ExchangeCallback(...) (application.CredentialPayload, error) {
  // require code + shop_id, call /api/v2/auth/token/get
}

func (a *Adapter) Refresh(...) (application.CredentialPayload, error) {
  // call /api/v2/auth/access_token/get
}
```

```go
func Sign(partnerKey string, base string) string {
  mac := hmac.New(sha256.New, []byte(partnerKey))
  mac.Write([]byte(base))
  return hex.EncodeToString(mac.Sum(nil))
}
```

- [ ] **Step 4: Run adapter tests**

Run:
```bash
cd apps/server_core
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/adapters/shopee -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/integrations/adapters/shopee/auth_adapter.go apps/server_core/internal/modules/integrations/adapters/shopee/signer.go apps/server_core/internal/modules/integrations/adapters/shopee/http_types.go apps/server_core/internal/modules/integrations/adapters/shopee/auth_adapter_test.go apps/server_core/internal/modules/integrations/adapters/shopee/signer_test.go
git commit -m "feat(integrations): implement shopee signed partner auth adapter"
```

---

### Task 3: Wire Callback Query Mapping for Shopee (`shop_id`)

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/transport/auth_handler.go`
- Test: `apps/server_core/internal/modules/integrations/transport/auth_handler_test.go`
- Modify: `contracts/api/marketplace-central.openapi.yaml`

- [ ] **Step 1: Write failing transport tests**

```go
func TestCallbackAcceptsShopIDAsProviderAccountID(t *testing.T) {
  // GET /integrations/auth/callback?code=x&state=y&shop_id=123
  // expect HandleCallbackInput.ProviderAccountID == "123"
}
```

- [ ] **Step 2: Run tests to verify failure**

Run:
```bash
cd apps/server_core
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/transport -run Callback -count=1
```

Expected: FAIL (current mapping only uses Amazon-style query key).

- [ ] **Step 3: Implement query fallback order + contract update**

```go
ProviderAccountID: firstNonEmptyQuery(r.URL.Query(), "shop_id", "merchant_id", "selling_partner_id"),
```

```yaml
/integrations/auth/callback:
  get:
    parameters:
      - name: shop_id
        in: query
        required: false
        schema: { type: string }
```

- [ ] **Step 4: Run tests**

Run:
```bash
cd apps/server_core
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/transport -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/integrations/transport/auth_handler.go apps/server_core/internal/modules/integrations/transport/auth_handler_test.go contracts/api/marketplace-central.openapi.yaml
git commit -m "fix(integrations): map shopee callback shop_id in auth callback"
```

---

### Task 4: Promote Shopee Provider Definition from Blocked to Interactive

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/adapters/shopee/auth_adapter.go`
- Test: `apps/server_core/internal/modules/integrations/adapters/providers/registry_test.go`
- Modify: `packages/feature-marketplaces/src/ProviderCatalogPanel.tsx`
- Test: `packages/feature-marketplaces/src/MarketplaceSettingsPage.test.tsx`

- [ ] **Step 1: Write failing metadata/UX tests**

```go
func TestShopeeDefinitionUsesInteractiveSignedAuth(t *testing.T) {
  // expect auth_strategy=shopee_partner, install_mode=interactive, execution_mode=available
}
```

```ts
it("renders authorize flow for shopee interactive installation", async () => {
  // provider.auth_strategy = "shopee_partner"
  // expect primary action = Authorize/Reauthorize by installation status
});
```

- [ ] **Step 2: Run tests to verify failure**

Run:
```bash
cd apps/server_core
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/adapters/providers -count=1
```

Run:
```bash
npm run -w packages/feature-marketplaces test -- --run
```

Expected: FAIL while still blocked/manual/api_key assumptions remain.

- [ ] **Step 3: Implement provider metadata and minimal UI adjustment**

```go
AuthStrategy: domain.AuthStrategyShopeePartner,
InstallMode:  domain.InstallModeInteractive,
Metadata: map[string]any{
  "rollout_stage": "v1",
  "execution_mode": "available",
  "docs_url": "https://open.shopee.com/developer-guide/20",
},
DeclaredCapabilities: []string{"pricing_fee_sync"},
```

```tsx
<span className="font-medium text-slate-900">{provider.auth_strategy}</span>
```

(No new UI business logic; keep existing install-mode driven flow.)

- [ ] **Step 4: Run backend + frontend tests**

Run:
```bash
cd apps/server_core
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/adapters/providers -count=1
```

Run:
```bash
npm run -w packages/feature-marketplaces test -- --run
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/integrations/adapters/shopee/auth_adapter.go apps/server_core/internal/modules/integrations/adapters/providers/registry_test.go packages/feature-marketplaces/src/ProviderCatalogPanel.tsx packages/feature-marketplaces/src/MarketplaceSettingsPage.test.tsx
git commit -m "feat(shopee): enable interactive provider metadata and rollout"
```

---

### Task 5: Auth Flow Service Integration Regression Coverage

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/application/auth_flow_service_test.go`
- Modify: `apps/server_core/internal/modules/integrations/application/auth_flow_service_security_test.go`

- [ ] **Step 1: Write failing integration tests**

```go
func TestStartAuthorizeShopeePersistsPendingConnection(t *testing.T) { /* adapter StartAuthorize + status transition */ }
func TestHandleCallbackShopeeRotatesCredentialAndStoresShopID(t *testing.T) { /* ProviderAccountID via shop_id */ }
func TestRefreshCredentialShopeePreservesShopIDWhenProviderOmits(t *testing.T) { /* safety fallback */ }
```

- [ ] **Step 2: Run tests to verify failure**

Run:
```bash
cd apps/server_core
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/application -run Shopee -count=1
```

Expected: FAIL until Shopee flow assertions are covered.

- [ ] **Step 3: Add minimal fixtures/mocks for Shopee adapter path**

```go
mockAdapter := &stubAdapter{
  providerCode: "shopee",
  startAuthorize: func(...) ...,
  exchangeCallback: func(...) application.CredentialPayload{ SecretType: "shopee_partner", ProviderAccountID: "123456" },
}
```

- [ ] **Step 4: Run package tests**

Run:
```bash
cd apps/server_core
$env:GOCACHE='.gocache'; go test ./internal/modules/integrations/application -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/integrations/application/auth_flow_service_test.go apps/server_core/internal/modules/integrations/application/auth_flow_service_security_test.go
git commit -m "test(integrations): add shopee auth flow regression coverage"
```

---

### Task 6: Wiki Keeper Continuous Sync During Implementation

**Files:**
- Modify: `wiki/framework/vendors/shopee/README.md`
- Modify: `wiki/framework/vendors/shopee/getting-started.md`
- Modify: `wiki/framework/vendors/shopee/api-best-practices.md`
- Modify: `wiki/framework/vendors/shopee/capability-matrix.md`
- Modify/Create: `wiki/framework/vendors/shopee/implementation-sync.md`
- Modify: `wiki/modules/integrations.md`

- [ ] **Step 1: Start wiki update lane from implementation kickoff**

Use this runbook:
```md
wiki/framework/vendors/shopee/implementation-sync.md
```

- [ ] **Step 2: For each merged Shopee task, append changelog row + update impacted sections**

```md
| 2026-04-27 | auth | switched to shopee_partner signed auth | go test ... PASS | @worker | <commit> |
```

- [ ] **Step 3: Verify doc/code alignment**

Run:
```bash
Select-String -Path wiki/framework/vendors/shopee/*.md -Pattern "api_key|blocked" -SimpleMatch
```

Expected: no stale statements if provider is now interactive.

- [ ] **Step 4: Commit docs update after each implementation wave**

```bash
git add wiki/framework/vendors/shopee wiki/modules/integrations.md
git commit -m "docs(shopee): sync implementation status and runbook evidence"
```

---

## Final Verification Gate (before marking Shopee done)

Run:
```bash
cd apps/server_core
$env:GOCACHE='.gocache'; go test ./...
```

Run:
```bash
npm run -w packages/sdk-runtime build
npm run -w packages/feature-marketplaces test -- --run
npm run -w apps/web build
```

Manual checks:
- Provider catalog shows Shopee as interactive with correct auth strategy.
- Start authorize returns a valid Shopee auth URL.
- Callback with `shop_id` finalizes installation as `connected`.
- Refresh path rotates credentials without dropping provider account linkage.

---

## Self-Review

### 1. Spec coverage
- Shopee implementation via framework: covered (Tasks 1-5).
- Parallel subagent execution: covered (Parallel Subagent Plan).
- Wiki continuously updated: covered (Task 6, plus `implementation-sync.md` runbook).
- Professional-grade verification: covered (Final Verification Gate).

### 2. Placeholder scan
- No TODO/TBD placeholders left.
- Every task includes explicit file targets, commands, and commit actions.

### 3. Type consistency
- `shopee_partner` propagated in domain, OpenAPI, and SDK type unions.
- Callback provider account mapping uses shared `ProviderAccountID` semantics across auth flow.

---

Plan complete and saved to `docs/superpowers/plans/2026-04-27-shopee-framework-implementation.md`. Two execution options:

1. Subagent-Driven (recommended) - I dispatch a fresh subagent per task, review between tasks, fast iteration
2. Inline Execution - Execute tasks in this session using executing-plans, batch execution with checkpoints

Which approach?
