# Marketplace Integration Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build slice 1 of the marketplace integration foundation: honest installation connection snapshots, runtime capabilities, credential-safe connector execution, Mercado Livre live read/probe paths, and aligned API/SDK/UI.

**Architecture:** `integrations` owns installation/auth/credential/capability/operation governance. `connectors` owns provider API execution and payload/error translation. Business modules consume normalized capability results and own business policy; no provider writes are exposed in this slice.

**Tech Stack:** Go `apps/server_core`, PostgreSQL migrations with `pgxpool`, OpenAPI YAML, TypeScript `packages/sdk-runtime`, React `packages/feature-integrations`, Vitest, Go tests with `GOCACHE=.gocache`.

## Global Constraints

- Mercado Livre is the first real adapter, not the architecture center.
- Runtime capabilities appear only when implemented, wired, composed, and executable.
- Future writes/messages/shipments/webhooks are modeled but not exposed as operational capabilities in slice 1.
- No provider writes in this slice.
- Do not expose raw tokens outside `integrations`.
- Do not show future capabilities as runtime available.
- Update OpenAPI and `packages/sdk-runtime` together for API changes.
- Never claim live integration validation from mocks/fakes/seams.
- Fake/mock/seam tests prove only local contract and deterministic behavior.
- Live Mercado Livre evidence is required before declaring real provider integration validated.
- Unknown fees, stock, freight, tax, product linkage, or account identity are explicit quality states, not zero/default values.

---

## File Structure

### Backend Integrations

- Modify `apps/server_core/internal/modules/integrations/domain/installation.go`: keep existing installation fields for compatibility and add nested connection/runtime capability fields.
- Create `apps/server_core/internal/modules/integrations/domain/connection_snapshot.go`: connection state, health, auth strategy, next action, reauth reason, and constructor helpers.
- Create `apps/server_core/internal/modules/integrations/domain/runtime_capability.go`: platform runtime capability enums and validation helpers.
- Modify `apps/server_core/internal/modules/integrations/domain/operation_run.go`: add provider evidence, translated error code, and duration fields when they are missing from the current operation run model.
- Modify `apps/server_core/internal/modules/integrations/ports/installation_repository.go`: replace narrow setters with `ApplyConnectionSnapshot`.
- Modify `apps/server_core/internal/modules/integrations/adapters/postgres/installation_repo.go`: scan/persist connection fields and runtime-safe projection.
- Create `apps/server_core/migrations/0018_integration_connection_snapshot.sql`: add snapshot columns needed beyond existing flattened columns.
- Modify `apps/server_core/internal/modules/integrations/application/auth_flow_service.go`: implement `ApplyAuthResult` semantics inside callback/API-key/refresh.
- Create `apps/server_core/internal/modules/integrations/application/credential_resolver.go`: resolve active access token by tenant/installation/provider account without leaking secrets.
- Create `apps/server_core/internal/modules/integrations/application/provider_operation_service.go`: validate runtime capability, open/finish operation runs, and call connector ports for read/probe operations.
- Modify `apps/server_core/internal/modules/integrations/application/installation_service.go`: expose installation list/get with connection snapshot and runtime capabilities.
- Modify `apps/server_core/internal/modules/integrations/transport/http_handler.go`: return the new installation shape and add explicit read/probe operation endpoints for account, listings, orders, and fee quote.
- Modify `apps/server_core/internal/modules/integrations/transport/auth_handler.go`: ensure auth status mirrors installation connection snapshot.

### Backend Connectors

- Modify `apps/server_core/internal/modules/connectors/domain/capability.go`: add account probe, fee quote snapshots, and typed provider error codes for auth/provider unavailability.
- Modify `apps/server_core/internal/modules/connectors/ports/marketplace_capability.go`: add `AccountProbe` and `FeeQuoteReader`; keep `StockWriter` modeled but not registered for slice 1 runtime.
- Modify `apps/server_core/internal/modules/connectors/application/marketplace_capability_service.go`: expose only registered read/probe capability getters.
- Modify `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go`: implement account probe and fee quote reads; ensure provider payloads remain private.
- Modify `apps/server_core/internal/modules/connectors/adapters/mercado_livre/fee_sync.go`: demote static seed behavior so it is not advertised as live runtime sync; live fee evidence for slice 1 comes from `fee_quote_read`.
- Modify `apps/server_core/internal/composition/root.go`: instantiate Mercado Livre capability adapter with `integrations` credential resolver; do not register stock write in slice 1.

### API, SDK, UI

- Modify `contracts/api/marketplace-central.openapi.yaml`: add `IntegrationConnectionSnapshot`, `IntegrationRuntimeCapability`, and operation evidence/result fields.
- Modify `packages/sdk-runtime/src/index.ts`: mirror OpenAPI installation/auth/operation models.
- Modify `packages/sdk-runtime/src/index.test.ts`: assert client paths and model compatibility where existing patterns allow.
- Modify `packages/feature-integrations/src/components/InstallationCard.tsx`: render connection snapshot instead of credential id.
- Modify `packages/feature-integrations/src/components/InstallationDrawer.tsx`: render connection, capabilities, and audit without raw credential details.
- Modify `packages/feature-integrations/src/components/AuthStatusPanel.tsx`: use snapshot/auth next action consistently.
- Modify `packages/feature-integrations/src/components/OperationsTimeline.tsx`: show translated operation evidence/errors.
- Modify `packages/feature-integrations/src/IntegrationsHubPage.tsx`: consume new SDK shape and refresh operation state.
- Modify `packages/feature-integrations/src/IntegrationsHubPage.test.tsx`: update fixtures and assert no misleading credential-disconnected UI.

### Validation And Evidence

- Modify the relevant `.mnfs/MIS-001-mercado-livre-operating-cockpit/.../validation.md` files if this implementation is executed under MNFS.
- Modify the relevant `.mnfs/.../validation-result.md` only after real verification.
- Update `docs/superpowers/specs/2026-07-08-marketplace-integration-foundation-design.md` only if implementation reveals architecture truth that changes the approved design.

---

## Task 1: Read-Only Audit And Slice Lock

**Files:**
- Read: `apps/server_core/internal/modules/integrations/**/*`
- Read: `apps/server_core/internal/modules/connectors/**/*`
- Read: `apps/server_core/internal/composition/root.go`
- Read: `contracts/api/marketplace-central.openapi.yaml`
- Read: `packages/sdk-runtime/src/index.ts`
- Read: `packages/feature-integrations/src/**/*`
- Create: `docs/superpowers/plans/2026-07-08-marketplace-integration-foundation-audit.md`

**Interfaces:**
- Consumes: approved spec `docs/superpowers/specs/2026-07-08-marketplace-integration-foundation-design.md`.
- Produces: an audit document listing exact current gaps and confirming slice 1 capability list.

- [ ] **Step 1: Capture current repo state**

Run:

```powershell
git status --short
```

Expected: may show pre-existing dirty work. Do not revert unrelated changes.

- [ ] **Step 2: Audit integrations files**

Run:

```powershell
rg "ExternalAccount|ActiveCredential|UpdateStatus|SetProviderAccountID|AuthStatus|OperationRun|Capability" apps/server_core/internal/modules/integrations -n
```

Expected: identify where installation projection, auth status, operation runs, and capability states currently live.

- [ ] **Step 3: Audit connectors files**

Run:

```powershell
rg "ProviderCapabilitySet|StockWrites|ListListings|ListOrders|fee_sync|accessTokenResolver" apps/server_core/internal/modules/connectors apps/server_core/internal/composition/root.go -n
```

Expected: confirm Mercado Livre adapter exists but is not instantiated in composition; confirm stock write exists but must not be exposed in slice 1.

- [ ] **Step 4: Audit API/SDK/UI files**

Run:

```powershell
rg "IntegrationInstallation|active_credential_id|Credential|last_verified|runtime_capabilities|declared_capabilities" contracts/api/marketplace-central.openapi.yaml packages/sdk-runtime/src/index.ts packages/feature-integrations/src -n
```

Expected: confirm UI currently renders credential-centric state and SDK/OpenAPI expose flattened installation fields.

- [ ] **Step 5: Write audit document**

Create `docs/superpowers/plans/2026-07-08-marketplace-integration-foundation-audit.md` with this exact structure:

```markdown
# Marketplace Integration Foundation Audit

Date: 2026-07-08

## Confirmed Current State

- `AuthFlowService.HandleCallback` rotates credential and updates auth session/status, but does not atomically project the full installation connection snapshot.
- `IntegrationInstallation` is flattened and exposes credential-oriented fields.
- Mercado Livre connector read paths exist but are not instantiated in composition as runtime marketplace capabilities.
- Mercado Livre stock write code exists but is outside slice 1 runtime exposure.
- `pricing_fee_sync` must be verified as live provider-backed or demoted from runtime-live semantics.

## Slice 1 Runtime Capabilities

- `account_probe`
- `listing_read`
- `order_read`
- `fee_quote_read` only after live provider-backed implementation
- `stock_read` only after composed and live validated

## Explicitly Modeled But Not Runtime Available

- `stock_write`
- `message_reply`
- `shipment_write`
- `webhook_receive`
- `listing_write`

## Stop Conditions

- Stop if OpenAPI, SDK, backend, or UI cannot agree on one installation shape.
- Stop if implementing live fee sync requires unsafe writes or speculative provider semantics.
- Stop if a test seam is the only evidence for a live provider claim.
```

- [ ] **Step 6: Commit audit**

Run:

```powershell
git add -- docs/superpowers/plans/2026-07-08-marketplace-integration-foundation-audit.md
git commit -m "docs(integrations): audit marketplace foundation slice"
```

Expected: commit only the audit doc.

---

## Task 2: Model Connection Snapshot And Runtime Capabilities

**Files:**
- Create: `apps/server_core/internal/modules/integrations/domain/connection_snapshot.go`
- Create: `apps/server_core/internal/modules/integrations/domain/runtime_capability.go`
- Modify: `apps/server_core/internal/modules/integrations/domain/installation.go`
- Create: `apps/server_core/internal/modules/integrations/domain/connection_snapshot_test.go`
- Create: `apps/server_core/internal/modules/integrations/domain/runtime_capability_test.go`

**Interfaces:**
- Consumes: no new code from previous tasks except audit decisions.
- Produces:
  - `domain.ConnectionSnapshot`
  - `domain.RuntimeCapability`
  - `domain.NewConnectedSnapshot(input domain.ConnectedSnapshotInput) domain.ConnectionSnapshot`
  - `domain.RuntimeCapability.Available() bool`

- [ ] **Step 1: Write failing connection snapshot tests**

Create `apps/server_core/internal/modules/integrations/domain/connection_snapshot_test.go`:

```go
package domain

import (
	"testing"
	"time"
)

func TestNewConnectedSnapshotNormalizesAccountAndNextAction(t *testing.T) {
	now := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	expires := now.Add(6 * time.Hour)

	snapshot := NewConnectedSnapshot(ConnectedSnapshotInput{
		ProviderCode:        " mercado_livre ",
		ExternalAccountID:   " 691607102 ",
		ExternalAccountName: " METALNOBREACABAMENTOS ",
		AuthStrategy:        AuthStrategyOAuth2,
		VerifiedAt:          now,
		ExpiresAt:           &expires,
	})

	if snapshot.State != ConnectionStateConnected {
		t.Fatalf("state = %q, want %q", snapshot.State, ConnectionStateConnected)
	}
	if snapshot.Health != HealthStatusHealthy {
		t.Fatalf("health = %q, want %q", snapshot.Health, HealthStatusHealthy)
	}
	if snapshot.ProviderCode != "mercado_livre" {
		t.Fatalf("provider = %q", snapshot.ProviderCode)
	}
	if snapshot.ExternalAccountID != "691607102" {
		t.Fatalf("external account id = %q", snapshot.ExternalAccountID)
	}
	if snapshot.ExternalAccountName != "METALNOBREACABAMENTOS" {
		t.Fatalf("external account name = %q", snapshot.ExternalAccountName)
	}
	if snapshot.NextAction != ConnectionNextActionNone {
		t.Fatalf("next action = %q, want none", snapshot.NextAction)
	}
	if snapshot.ExpiresAt == nil || !snapshot.ExpiresAt.Equal(expires) {
		t.Fatalf("expires at was not preserved")
	}
}

func TestDisconnectedSnapshotRequiresAuthAction(t *testing.T) {
	snapshot := NewDisconnectedSnapshot("mercado_livre", "reauth required")

	if snapshot.State != ConnectionStateDisconnected {
		t.Fatalf("state = %q", snapshot.State)
	}
	if snapshot.Health != HealthStatusWarning {
		t.Fatalf("health = %q", snapshot.Health)
	}
	if snapshot.NextAction != ConnectionNextActionAuthorize {
		t.Fatalf("next action = %q", snapshot.NextAction)
	}
	if snapshot.ReauthReason != "reauth required" {
		t.Fatalf("reauth reason = %q", snapshot.ReauthReason)
	}
}
```

- [ ] **Step 2: Write failing runtime capability tests**

Create `apps/server_core/internal/modules/integrations/domain/runtime_capability_test.go`:

```go
package domain

import "testing"

func TestRuntimeCapabilityAvailableOnlyWhenExecutable(t *testing.T) {
	capability := RuntimeCapability{
		Code:       RuntimeCapabilityListingRead,
		State:      RuntimeCapabilityStateAvailable,
		Executable: true,
	}

	if !capability.Available() {
		t.Fatalf("expected capability to be available")
	}
}

func TestRuntimeCapabilityUnavailableWhenNotExecutable(t *testing.T) {
	capability := RuntimeCapability{
		Code:       RuntimeCapabilityStockWrite,
		State:      RuntimeCapabilityStateAvailable,
		Executable: false,
	}

	if capability.Available() {
		t.Fatalf("stock write must not be available without executable runtime path")
	}
}

func TestRuntimeCapabilityRejectsEmptyCode(t *testing.T) {
	err := ValidateRuntimeCapability(RuntimeCapability{State: RuntimeCapabilityStateAvailable})
	if err == nil {
		t.Fatalf("expected empty capability code to be invalid")
	}
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/integrations/domain -run "TestNewConnectedSnapshot|TestDisconnectedSnapshot|TestRuntimeCapability" -count=1
```

Expected: FAIL because `ConnectionSnapshot`, `RuntimeCapability`, and helper names are not defined.

- [ ] **Step 4: Implement connection snapshot domain**

Create `apps/server_core/internal/modules/integrations/domain/connection_snapshot.go`:

```go
package domain

import (
	"strings"
	"time"
)

type ConnectionState string
type AuthStrategy string
type ConnectionNextAction string

const (
	ConnectionStateDraft        ConnectionState = "draft"
	ConnectionStatePending      ConnectionState = "pending_connection"
	ConnectionStateConnected    ConnectionState = "connected"
	ConnectionStateDegraded     ConnectionState = "degraded"
	ConnectionStateNeedsReauth  ConnectionState = "needs_reauth"
	ConnectionStateDisconnected ConnectionState = "disconnected"

	AuthStrategyOAuth2 AuthStrategy = "oauth2"
	AuthStrategyAPIKey AuthStrategy = "api_key"
	AuthStrategyNone   AuthStrategy = "none"

	ConnectionNextActionNone      ConnectionNextAction = "none"
	ConnectionNextActionAuthorize ConnectionNextAction = "authorize"
	ConnectionNextActionReauth    ConnectionNextAction = "reauth"
	ConnectionNextActionConfigure ConnectionNextAction = "configure"
	ConnectionNextActionRetry     ConnectionNextAction = "retry"
)

type ConnectionSnapshot struct {
	State               ConnectionState       `json:"state"`
	Health              HealthStatus          `json:"health"`
	ProviderCode        string                `json:"provider_code"`
	ExternalAccountID   string                `json:"external_account_id"`
	ExternalAccountName string                `json:"external_account_name"`
	AuthStrategy        AuthStrategy          `json:"auth_strategy"`
	LastVerifiedAt      *time.Time            `json:"last_verified_at,omitempty"`
	ExpiresAt           *time.Time            `json:"expires_at,omitempty"`
	NextAction          ConnectionNextAction  `json:"next_action"`
	ReauthReason        string                `json:"reauth_reason,omitempty"`
}

type ConnectedSnapshotInput struct {
	ProviderCode        string
	ExternalAccountID   string
	ExternalAccountName string
	AuthStrategy        AuthStrategy
	VerifiedAt          time.Time
	ExpiresAt           *time.Time
}

func NewConnectedSnapshot(input ConnectedSnapshotInput) ConnectionSnapshot {
	verifiedAt := input.VerifiedAt.UTC()
	var expiresAt *time.Time
	if input.ExpiresAt != nil {
		value := input.ExpiresAt.UTC()
		expiresAt = &value
	}
	authStrategy := input.AuthStrategy
	if authStrategy == "" {
		authStrategy = AuthStrategyNone
	}
	return ConnectionSnapshot{
		State:               ConnectionStateConnected,
		Health:              HealthStatusHealthy,
		ProviderCode:        strings.TrimSpace(input.ProviderCode),
		ExternalAccountID:   strings.TrimSpace(input.ExternalAccountID),
		ExternalAccountName: strings.TrimSpace(input.ExternalAccountName),
		AuthStrategy:        authStrategy,
		LastVerifiedAt:      &verifiedAt,
		ExpiresAt:           expiresAt,
		NextAction:          ConnectionNextActionNone,
	}
}

func NewDisconnectedSnapshot(providerCode string, reason string) ConnectionSnapshot {
	return ConnectionSnapshot{
		State:        ConnectionStateDisconnected,
		Health:       HealthStatusWarning,
		ProviderCode: strings.TrimSpace(providerCode),
		AuthStrategy: AuthStrategyNone,
		NextAction:   ConnectionNextActionAuthorize,
		ReauthReason: strings.TrimSpace(reason),
	}
}
```

- [ ] **Step 5: Implement runtime capability domain**

Create `apps/server_core/internal/modules/integrations/domain/runtime_capability.go`:

```go
package domain

import (
	"errors"
	"strings"
	"time"
)

type RuntimeCapabilityCode string
type RuntimeCapabilityState string

const (
	RuntimeCapabilityAccountProbe  RuntimeCapabilityCode = "account_probe"
	RuntimeCapabilityListingRead   RuntimeCapabilityCode = "listing_read"
	RuntimeCapabilityOrderRead     RuntimeCapabilityCode = "order_read"
	RuntimeCapabilityFeeQuoteRead  RuntimeCapabilityCode = "fee_quote_read"
	RuntimeCapabilityStockRead     RuntimeCapabilityCode = "stock_read"
	RuntimeCapabilityStockWrite    RuntimeCapabilityCode = "stock_write"
	RuntimeCapabilityMessageReply  RuntimeCapabilityCode = "message_reply"
	RuntimeCapabilityShipmentRead  RuntimeCapabilityCode = "shipment_read"
	RuntimeCapabilityWebhookReceive RuntimeCapabilityCode = "webhook_receive"

	RuntimeCapabilityStateAvailable     RuntimeCapabilityState = "available"
	RuntimeCapabilityStateUnavailable   RuntimeCapabilityState = "unavailable"
	RuntimeCapabilityStateNeedsAuth     RuntimeCapabilityState = "needs_auth"
	RuntimeCapabilityStateDegraded      RuntimeCapabilityState = "degraded"
	RuntimeCapabilityStateNotConfigured RuntimeCapabilityState = "not_configured"
)

type RuntimeCapability struct {
	Code               RuntimeCapabilityCode  `json:"code"`
	State              RuntimeCapabilityState `json:"state"`
	Executable         bool                   `json:"executable"`
	LiveValidated      bool                   `json:"live_validated"`
	LocalValidated     bool                   `json:"local_validated"`
	UnavailableReason  string                 `json:"unavailable_reason,omitempty"`
	LastValidatedAt    *time.Time             `json:"last_validated_at,omitempty"`
}

func (c RuntimeCapability) Available() bool {
	return c.Code != "" && c.State == RuntimeCapabilityStateAvailable && c.Executable
}

func ValidateRuntimeCapability(capability RuntimeCapability) error {
	if strings.TrimSpace(string(capability.Code)) == "" {
		return errors.New("INTEGRATIONS_CAPABILITY_INVALID")
	}
	if capability.State == "" {
		return errors.New("INTEGRATIONS_CAPABILITY_INVALID")
	}
	return nil
}
```

- [ ] **Step 6: Update installation domain**

Modify `apps/server_core/internal/modules/integrations/domain/installation.go` by adding fields at the end of `Installation`:

```go
	ConnectionSnapshot ConnectionSnapshot   `json:"connection"`
	RuntimeCapabilities []RuntimeCapability `json:"runtime_capabilities"`
```

Keep existing flattened fields for migration compatibility during this task. Do not remove `ExternalAccountID`, `ExternalAccountName`, `ActiveCredentialID`, or `LastVerifiedAt` yet.

- [ ] **Step 7: Run domain tests**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/integrations/domain -count=1
```

Expected: PASS.

- [ ] **Step 8: Commit domain model**

Run:

```powershell
git add -- apps/server_core/internal/modules/integrations/domain/connection_snapshot.go apps/server_core/internal/modules/integrations/domain/runtime_capability.go apps/server_core/internal/modules/integrations/domain/installation.go apps/server_core/internal/modules/integrations/domain/connection_snapshot_test.go apps/server_core/internal/modules/integrations/domain/runtime_capability_test.go
git commit -m "feat(integrations): model connection snapshot"
```

Expected: commit only task 2 files.

---

## Task 3: Persist Installation Connection Snapshot

**Files:**
- Create: `apps/server_core/migrations/0018_integration_connection_snapshot.sql`
- Modify: `apps/server_core/internal/modules/integrations/ports/installation_repository.go`
- Modify: `apps/server_core/internal/modules/integrations/adapters/postgres/installation_repo.go`
- Modify: `apps/server_core/internal/modules/integrations/application/installation_service.go`
- Create or modify: `apps/server_core/internal/modules/integrations/application/installation_service_test.go`

**Interfaces:**
- Consumes: `domain.ConnectionSnapshot`, `domain.RuntimeCapability`.
- Produces:
  - `InstallationRepository.ApplyConnectionSnapshot(ctx, installationID string, snapshot domain.ConnectionSnapshot, activeCredentialID string) error`
  - `InstallationService.ApplyConnectionSnapshot(ctx, installationID string, snapshot domain.ConnectionSnapshot, activeCredentialID string) error`

- [ ] **Step 1: Write failing service test**

Add to `apps/server_core/internal/modules/integrations/application/installation_service_test.go`:

```go
func TestInstallationServiceApplyConnectionSnapshotRejectsEmptyInstallation(t *testing.T) {
	svc := NewInstallationService(&fakeInstallationRepo{}, "tenant_default")

	err := svc.ApplyConnectionSnapshot(context.Background(), "", domain.ConnectionSnapshot{}, "cred_1")

	if err == nil || err.Error() != "INTEGRATIONS_INSTALLATION_INVALID" {
		t.Fatalf("err = %v, want INTEGRATIONS_INSTALLATION_INVALID", err)
	}
}

func TestInstallationServiceApplyConnectionSnapshotDelegatesToRepository(t *testing.T) {
	repo := &fakeInstallationRepo{}
	svc := NewInstallationService(repo, "tenant_default")
	now := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	snapshot := domain.NewConnectedSnapshot(domain.ConnectedSnapshotInput{
		ProviderCode:        "mercado_livre",
		ExternalAccountID:   "691607102",
		ExternalAccountName: "METALNOBREACABAMENTOS",
		AuthStrategy:        domain.AuthStrategyOAuth2,
		VerifiedAt:          now,
	})

	err := svc.ApplyConnectionSnapshot(context.Background(), "inst-1", snapshot, "cred-1")

	if err != nil {
		t.Fatalf("ApplyConnectionSnapshot error = %v", err)
	}
	if repo.appliedInstallationID != "inst-1" {
		t.Fatalf("installation id = %q", repo.appliedInstallationID)
	}
	if repo.appliedSnapshot.ExternalAccountID != "691607102" {
		t.Fatalf("external account = %q", repo.appliedSnapshot.ExternalAccountID)
	}
	if repo.appliedCredentialID != "cred-1" {
		t.Fatalf("credential id = %q", repo.appliedCredentialID)
	}
}
```

If `fakeInstallationRepo` already exists in this test file, extend it with fields and method:

```go
type fakeInstallationRepo struct {
	appliedInstallationID string
	appliedSnapshot       domain.ConnectionSnapshot
	appliedCredentialID   string
}

func (f *fakeInstallationRepo) ApplyConnectionSnapshot(ctx context.Context, installationID string, snapshot domain.ConnectionSnapshot, activeCredentialID string) error {
	f.appliedInstallationID = installationID
	f.appliedSnapshot = snapshot
	f.appliedCredentialID = activeCredentialID
	return nil
}
```

- [ ] **Step 2: Run test to verify failure**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/integrations/application -run TestInstallationServiceApplyConnectionSnapshot -count=1
```

Expected: FAIL because `ApplyConnectionSnapshot` is not defined on service/repository.

- [ ] **Step 3: Add migration**

Create `apps/server_core/migrations/0018_integration_connection_snapshot.sql`:

```sql
ALTER TABLE integration_installations
  ADD COLUMN IF NOT EXISTS connection_state text NOT NULL DEFAULT 'draft',
  ADD COLUMN IF NOT EXISTS auth_strategy text NOT NULL DEFAULT 'none',
  ADD COLUMN IF NOT EXISTS credential_expires_at timestamptz,
  ADD COLUMN IF NOT EXISTS next_action text NOT NULL DEFAULT 'authorize',
  ADD COLUMN IF NOT EXISTS reauth_reason text NOT NULL DEFAULT '';

UPDATE integration_installations
SET connection_state = status
WHERE connection_state = 'draft'
  AND status IN ('draft', 'pending_connection', 'connected', 'degraded', 'disconnected');

UPDATE integration_installations
SET next_action = CASE
  WHEN status = 'connected' THEN 'none'
  WHEN status = 'requires_reauth' THEN 'reauth'
  WHEN status = 'pending_connection' THEN 'authorize'
  ELSE next_action
END;
```

- [ ] **Step 4: Extend installation repository port**

Modify `apps/server_core/internal/modules/integrations/ports/installation_repository.go`:

```go
ApplyConnectionSnapshot(ctx context.Context, installationID string, snapshot domain.ConnectionSnapshot, activeCredentialID string) error
```

Keep existing narrow methods temporarily to minimize blast radius. They can be removed in a later cleanup task if unused.

- [ ] **Step 5: Implement service method**

Add to `apps/server_core/internal/modules/integrations/application/installation_service.go`:

```go
func (s *InstallationService) ApplyConnectionSnapshot(ctx context.Context, installationID string, snapshot domain.ConnectionSnapshot, activeCredentialID string) error {
	installationID = strings.TrimSpace(installationID)
	if installationID == "" {
		return errors.New("INTEGRATIONS_INSTALLATION_INVALID")
	}
	if snapshot.State == "" || snapshot.Health == "" {
		return errors.New("INTEGRATIONS_INSTALLATION_INVALID")
	}
	return s.repo.ApplyConnectionSnapshot(ctx, installationID, snapshot, strings.TrimSpace(activeCredentialID))
}
```

- [ ] **Step 6: Implement PostgreSQL persistence**

Add to `apps/server_core/internal/modules/integrations/adapters/postgres/installation_repo.go`:

```go
func (r *InstallationRepository) ApplyConnectionSnapshot(ctx context.Context, installationID string, snapshot domain.ConnectionSnapshot, activeCredentialID string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE integration_installations
		SET status = $3,
		    health_status = $4,
		    external_account_id = $5,
		    external_account_name = $6,
		    active_credential_id = NULLIF($7, ''),
		    last_verified_at = $8,
		    connection_state = $9,
		    auth_strategy = $10,
		    credential_expires_at = $11,
		    next_action = $12,
		    reauth_reason = $13,
		    updated_at = now()
		WHERE tenant_id = $1
		  AND installation_id = $2
	`, r.tenantID, installationID,
		domain.InstallationStatus(snapshot.State),
		snapshot.Health,
		snapshot.ExternalAccountID,
		snapshot.ExternalAccountName,
		activeCredentialID,
		snapshot.LastVerifiedAt,
		snapshot.State,
		snapshot.AuthStrategy,
		snapshot.ExpiresAt,
		snapshot.NextAction,
		snapshot.ReauthReason,
	)
	return err
}
```

- [ ] **Step 7: Extend scan query and scanInstallation**

Modify all `SELECT` lists in `installation_repo.go` to include:

```sql
connection_state, auth_strategy, credential_expires_at, next_action, reauth_reason
```

Extend `scanInstallation` with `pgtype.Timestamptz credentialExpires` and string variables. After scan, set:

```go
inst.ConnectionSnapshot = domain.ConnectionSnapshot{
	State:               domain.ConnectionState(connectionState),
	Health:              inst.HealthStatus,
	ProviderCode:        inst.ProviderCode,
	ExternalAccountID:   inst.ExternalAccountID,
	ExternalAccountName: inst.ExternalAccountName,
	AuthStrategy:        domain.AuthStrategy(authStrategy),
	LastVerifiedAt:      inst.LastVerifiedAt,
	NextAction:          domain.ConnectionNextAction(nextAction),
	ReauthReason:        reauthReason,
}
if credentialExpires.Valid {
	ts := credentialExpires.Time.UTC()
	inst.ConnectionSnapshot.ExpiresAt = &ts
}
```

- [ ] **Step 8: Run integration application tests**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/integrations/application ./apps/server_core/internal/modules/integrations/adapters/postgres -count=1
```

Expected: PASS. If Postgres adapter tests require a real database and skip when unavailable, record that as local adapter test behavior.

- [ ] **Step 9: Commit persistence**

Run:

```powershell
git add -- apps/server_core/migrations/0018_integration_connection_snapshot.sql apps/server_core/internal/modules/integrations/ports/installation_repository.go apps/server_core/internal/modules/integrations/adapters/postgres/installation_repo.go apps/server_core/internal/modules/integrations/application/installation_service.go apps/server_core/internal/modules/integrations/application/installation_service_test.go
git commit -m "feat(integrations): persist connection snapshot"
```

Expected: commit only task 3 files.

---

## Task 4: Apply Auth Result Atomically

**Files:**
- Modify: `apps/server_core/internal/modules/integrations/application/auth_flow_service.go`
- Modify: `apps/server_core/internal/modules/integrations/application/auth_flow_service_test.go`
- Modify: `apps/server_core/internal/modules/integrations/application/auth_flow_service_security_test.go`
- Modify: `apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter.go`
- Modify: `apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter_test.go`

**Interfaces:**
- Consumes:
  - `InstallationService.ApplyConnectionSnapshot`
  - `domain.NewConnectedSnapshot`
- Produces:
  - `AuthFlowService.applyAuthResult(ctx, inst, payload) (AuthStatus, error)`
  - Credential payload includes `provider_account_name` in encrypted metadata.

- [ ] **Step 1: Add failing callback projection test**

Add to `apps/server_core/internal/modules/integrations/application/auth_flow_service_test.go`:

```go
func TestAuthFlowHandleCallbackProjectsConnectionSnapshot(t *testing.T) {
	clock := fixedAuthFlowClock{now: time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)}
	stores := newAuthFlowTestStores(t)
	stores.installations.items["inst-1"] = domain.Installation{
		InstallationID: "inst-1",
		TenantID:       "tenant_default",
		ProviderCode:   "mercado_livre",
		Family:         domain.IntegrationFamilyMarketplace,
		DisplayName:    "Mercado Livre",
		Status:         domain.InstallationStatusPendingConnection,
		HealthStatus:   domain.HealthStatusWarning,
	}
	adapter := &fakeMarketplaceAuthAdapter{
		providerCode: "mercado_livre",
		callbackPayload: CredentialPayload{
			SecretType:          "oauth2",
			AccessToken:         "access-token",
			RefreshToken:        "refresh-token",
			ProviderAccountID:   "691607102",
			ProviderAccountName: "METALNOBREACABAMENTOS",
			ExpiresAt:           ptrTime(clock.now.Add(time.Hour)),
		},
	}
	service := newTestAuthFlowService(t, stores, clock, adapter)
	state := stores.saveOAuthState(t, service, "inst-1")

	status, err := service.HandleCallback(context.Background(), HandleCallbackInput{
		InstallationID: "inst-1",
		State:          state,
		Code:           "code",
		RedirectURI:    "https://app.example/callback",
	})

	if err != nil {
		t.Fatalf("HandleCallback error = %v", err)
	}
	if status.Connection.ExternalAccountID != "691607102" {
		t.Fatalf("status connection account = %q", status.Connection.ExternalAccountID)
	}
	inst := stores.installations.items["inst-1"]
	if inst.ExternalAccountID != "691607102" {
		t.Fatalf("installation external account = %q", inst.ExternalAccountID)
	}
	if inst.ExternalAccountName != "METALNOBREACABAMENTOS" {
		t.Fatalf("installation external account name = %q", inst.ExternalAccountName)
	}
	if inst.ActiveCredentialID == "" {
		t.Fatalf("active credential id was not projected")
	}
	if inst.ConnectionSnapshot.NextAction != domain.ConnectionNextActionNone {
		t.Fatalf("next action = %q", inst.ConnectionSnapshot.NextAction)
	}
}
```

If the existing test helpers use different names, adapt only helper names while preserving assertions.

- [ ] **Step 2: Run test to verify failure**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/integrations/application -run TestAuthFlowHandleCallbackProjectsConnectionSnapshot -count=1
```

Expected: FAIL because `AuthStatus.Connection` and full projection do not exist.

- [ ] **Step 3: Extend AuthStatus**

Modify `AuthStatus` in `auth_flow_service.go`:

```go
type AuthStatus struct {
	InstallationID  string                    `json:"installation_id"`
	Status          domain.InstallationStatus `json:"status"`
	HealthStatus    domain.HealthStatus       `json:"health_status"`
	ProviderCode    string                    `json:"provider_code,omitempty"`
	ExternalAccount string                    `json:"external_account_id,omitempty"`
	Connection      domain.ConnectionSnapshot `json:"connection"`
}
```

- [ ] **Step 4: Return credential from saveCredential**

Change:

```go
func (s *AuthFlowService) saveCredential(ctx context.Context, installationID string, payload CredentialPayload) error
```

to:

```go
func (s *AuthFlowService) saveCredential(ctx context.Context, installationID string, payload CredentialPayload) (domain.Credential, error)
```

Add `provider_account_name` to encrypted payload:

```go
secret := map[string]any{
	"type":                  payload.SecretType,
	"access_token":          payload.AccessToken,
	"refresh_token":         payload.RefreshToken,
	"api_key":               payload.APIKey,
	"provider_account_id":   payload.ProviderAccountID,
	"provider_account_name": payload.ProviderAccountName,
}
```

Return the rotated credential:

```go
credential, err := s.credentials.Rotate(ctx, RotateCredentialInput{...})
return credential, err
```

- [ ] **Step 5: Add applyAuthResult helper**

Add to `auth_flow_service.go`:

```go
func (s *AuthFlowService) applyAuthResult(ctx context.Context, inst domain.Installation, payload CredentialPayload) (AuthStatus, error) {
	credential, err := s.saveCredential(ctx, inst.InstallationID, payload)
	if err != nil {
		return AuthStatus{}, err
	}

	now := s.clock.Now().UTC()
	if _, err := s.authSessions.Upsert(ctx, UpsertAuthSessionInput{
		AuthSessionID:        fmt.Sprintf("auth_%s", inst.InstallationID),
		InstallationID:       inst.InstallationID,
		ProviderAccountID:    payload.ProviderAccountID,
		State:                domain.AuthStateValid,
		AccessTokenExpiresAt: payload.ExpiresAt,
		LastVerifiedAt:       &now,
	}); err != nil {
		return AuthStatus{}, err
	}

	authStrategy := domain.AuthStrategyOAuth2
	if payload.SecretType == "api_key" {
		authStrategy = domain.AuthStrategyAPIKey
	}
	snapshot := domain.NewConnectedSnapshot(domain.ConnectedSnapshotInput{
		ProviderCode:        inst.ProviderCode,
		ExternalAccountID:   payload.ProviderAccountID,
		ExternalAccountName: payload.ProviderAccountName,
		AuthStrategy:        authStrategy,
		VerifiedAt:          now,
		ExpiresAt:           payload.ExpiresAt,
	})
	if err := s.installations.ApplyConnectionSnapshot(ctx, inst.InstallationID, snapshot, credential.CredentialID); err != nil {
		return AuthStatus{}, err
	}

	return AuthStatus{
		InstallationID:  inst.InstallationID,
		Status:          domain.InstallationStatusConnected,
		HealthStatus:    domain.HealthStatusHealthy,
		ProviderCode:    inst.ProviderCode,
		ExternalAccount: payload.ProviderAccountID,
		Connection:      snapshot,
	}, nil
}
```

- [ ] **Step 6: Replace callback/API-key/refresh projection**

In `HandleCallback`, replace separate `saveCredential`, `authSessions.Upsert`, and `UpdateStatus` calls with:

```go
return s.applyAuthResult(ctx, inst, payload)
```

In `SubmitAPIKey`, after adapter verification, call:

```go
return s.applyAuthResult(ctx, inst, payload)
```

In `RefreshCredential`, preserve refresh token fallback, then call:

```go
return s.applyAuthResult(ctx, inst, payload)
```

- [ ] **Step 7: Update auth status to mirror snapshot**

In `GetAuthStatus`, return:

```go
return AuthStatus{
	InstallationID:  input.InstallationID,
	Status:          inst.Status,
	HealthStatus:    inst.HealthStatus,
	ProviderCode:    inst.ProviderCode,
	ExternalAccount: inst.ExternalAccountID,
	Connection:      inst.ConnectionSnapshot,
}, nil
```

- [ ] **Step 8: Ensure Mercado Livre auth adapter carries account name**

In `apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter.go`, ensure the `/users/me` response maps account nickname/name into `CredentialPayload.ProviderAccountName`. If the adapter currently only sets account id, add:

```go
ProviderAccountName: firstNonEmpty(me.Nickname, me.FirstName, me.LastName, me.ID),
```

Keep provider response structs private to the adapter.

- [ ] **Step 9: Run auth flow tests**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/integrations/application ./apps/server_core/internal/modules/integrations/adapters/mercadolivre -count=1
```

Expected: PASS.

- [ ] **Step 10: Commit auth projection**

Run:

```powershell
git add -- apps/server_core/internal/modules/integrations/application/auth_flow_service.go apps/server_core/internal/modules/integrations/application/auth_flow_service_test.go apps/server_core/internal/modules/integrations/application/auth_flow_service_security_test.go apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter.go apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter_test.go
git commit -m "feat(integrations): project auth connection snapshot"
```

Expected: commit only task 4 files.

---

## Task 5: Add Credential Resolver Boundary

**Files:**
- Create: `apps/server_core/internal/modules/integrations/application/credential_resolver.go`
- Create: `apps/server_core/internal/modules/integrations/application/credential_resolver_test.go`
- Modify: `apps/server_core/internal/modules/integrations/application/auth_flow_service.go`

**Interfaces:**
- Consumes: active credentials and encryption service already used by auth flow.
- Produces:
  - `CredentialResolver.ResolveAccessToken(ctx context.Context, ref CredentialResolutionRef) (ResolvedCredential, error)`
  - `CredentialResolutionRef{TenantID, InstallationID, ProviderCode, ProviderAccountID string}`

- [ ] **Step 1: Write failing resolver tests**

Create `apps/server_core/internal/modules/integrations/application/credential_resolver_test.go`:

```go
package application

import (
	"context"
	"testing"
)

func TestCredentialResolverRejectsIncompleteReference(t *testing.T) {
	resolver := NewCredentialResolver(nil, nil)

	_, err := resolver.ResolveAccessToken(context.Background(), CredentialResolutionRef{
		TenantID:       "tenant_default",
		InstallationID: "inst-1",
		ProviderCode:   "mercado_livre",
	})

	if err == nil || err.Error() != "INTEGRATIONS_CREDENTIAL_INVALID_REFERENCE" {
		t.Fatalf("err = %v", err)
	}
}

func TestCredentialResolverDoesNotReturnRefreshToken(t *testing.T) {
	store := &fakeCredentialLookup{
		credential: encryptedCredentialForTest(t, map[string]any{
			"access_token":  "access-token",
			"refresh_token": "refresh-token",
		}),
		found: true,
	}
	encryptor := fakePayloadEncryptor{}
	resolver := NewCredentialResolver(store, encryptor)

	resolved, err := resolver.ResolveAccessToken(context.Background(), CredentialResolutionRef{
		TenantID:           "tenant_default",
		InstallationID:     "inst-1",
		ProviderCode:       "mercado_livre",
		ProviderAccountID:  "691607102",
	})

	if err != nil {
		t.Fatalf("ResolveAccessToken error = %v", err)
	}
	if resolved.AccessToken != "access-token" {
		t.Fatalf("access token = %q", resolved.AccessToken)
	}
	if resolved.ProviderAccountID != "691607102" {
		t.Fatalf("provider account id = %q", resolved.ProviderAccountID)
	}
}
```

Use existing fake encryptor/test helpers if present. The key assertion is that only an access token leaves the resolver.

- [ ] **Step 2: Run resolver test to verify failure**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/integrations/application -run TestCredentialResolver -count=1
```

Expected: FAIL because resolver types do not exist.

- [ ] **Step 3: Implement resolver**

Create `apps/server_core/internal/modules/integrations/application/credential_resolver.go`:

```go
package application

import (
	"context"
	"errors"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type CredentialResolutionRef struct {
	TenantID          string
	InstallationID    string
	ProviderCode      string
	ProviderAccountID string
}

type ResolvedCredential struct {
	InstallationID    string
	ProviderCode      string
	ProviderAccountID string
	AccessToken       string
}

type credentialLookup interface {
	GetActiveCredential(ctx context.Context, installationID string) (domain.Credential, bool, error)
}

type payloadDecryptor interface {
	DecryptJSON(encoded []byte) (map[string]any, string, error)
}

type CredentialResolver struct {
	credentials credentialLookup
	decryptor    payloadDecryptor
}

func NewCredentialResolver(credentials credentialLookup, decryptor payloadDecryptor) *CredentialResolver {
	return &CredentialResolver{credentials: credentials, decryptor: decryptor}
}

func (r *CredentialResolver) ResolveAccessToken(ctx context.Context, ref CredentialResolutionRef) (ResolvedCredential, error) {
	ref.TenantID = strings.TrimSpace(ref.TenantID)
	ref.InstallationID = strings.TrimSpace(ref.InstallationID)
	ref.ProviderCode = strings.TrimSpace(ref.ProviderCode)
	ref.ProviderAccountID = strings.TrimSpace(ref.ProviderAccountID)
	if ref.TenantID == "" || ref.InstallationID == "" || ref.ProviderCode == "" || ref.ProviderAccountID == "" {
		return ResolvedCredential{}, errors.New("INTEGRATIONS_CREDENTIAL_INVALID_REFERENCE")
	}
	if r.credentials == nil || r.decryptor == nil {
		return ResolvedCredential{}, errors.New("INTEGRATIONS_CREDENTIAL_RESOLVER_NOT_CONFIGURED")
	}

	credential, found, err := r.credentials.GetActiveCredential(ctx, ref.InstallationID)
	if err != nil {
		return ResolvedCredential{}, err
	}
	if !found {
		return ResolvedCredential{}, domain.ErrCredentialNotFound
	}

	payload, _, err := r.decryptor.DecryptJSON(credential.EncryptedPayload)
	if err != nil {
		return ResolvedCredential{}, domain.ErrCredentialDecryptionFailed
	}
	accessToken, ok := credentialPayloadString(payload, "access_token")
	if !ok {
		return ResolvedCredential{}, errors.New("INTEGRATIONS_CREDENTIAL_ACCESS_TOKEN_MISSING")
	}

	return ResolvedCredential{
		InstallationID:    ref.InstallationID,
		ProviderCode:      ref.ProviderCode,
		ProviderAccountID: ref.ProviderAccountID,
		AccessToken:       accessToken,
	}, nil
}
```

- [ ] **Step 4: Run resolver tests**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/integrations/application -run TestCredentialResolver -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit resolver**

Run:

```powershell
git add -- apps/server_core/internal/modules/integrations/application/credential_resolver.go apps/server_core/internal/modules/integrations/application/credential_resolver_test.go
git commit -m "feat(integrations): add credential resolver"
```

Expected: commit only task 5 files.

---

## Task 6: Add Connector Account Probe And Fee Quote Ports

**Files:**
- Modify: `apps/server_core/internal/modules/connectors/domain/capability.go`
- Modify: `apps/server_core/internal/modules/connectors/ports/marketplace_capability.go`
- Modify: `apps/server_core/internal/modules/connectors/application/marketplace_capability_service.go`
- Modify: `apps/server_core/internal/modules/connectors/application/marketplace_capability_service_test.go`

**Interfaces:**
- Produces:
  - `ports.AccountProber`
  - `ports.FeeQuoteReader`
  - `domain.AccountSnapshot`
  - `domain.FeeQuoteInput`
  - `domain.FeeQuoteSnapshot`

- [ ] **Step 1: Write failing service tests**

Add to `apps/server_core/internal/modules/connectors/application/marketplace_capability_service_test.go`:

```go
func TestMarketplaceCapabilityServiceAccountProber(t *testing.T) {
	prober := fakeAccountProber{}
	service := NewMarketplaceCapabilityService([]ProviderCapabilitySet{{
		ProviderCode:  "mercado_livre",
		AccountProbes: prober,
	}})

	got, err := service.AccountProber("mercado_livre")

	if err != nil {
		t.Fatalf("AccountProber error = %v", err)
	}
	if got == nil {
		t.Fatalf("expected account prober")
	}
}

func TestMarketplaceCapabilityServiceFeeQuoteReaderUnsupported(t *testing.T) {
	service := NewMarketplaceCapabilityService([]ProviderCapabilitySet{{ProviderCode: "mercado_livre"}})

	_, err := service.FeeQuoteReader("mercado_livre")

	if err == nil {
		t.Fatalf("expected unsupported fee quote reader")
	}
}
```

Add fakes:

```go
type fakeAccountProber struct{}

func (fakeAccountProber) ProbeAccount(ctx context.Context, ref domain.ProviderAccountRef) (domain.AccountSnapshot, error) {
	return domain.AccountSnapshot{ProviderCode: ref.ProviderCode, ProviderAccountID: ref.ProviderAccountID}, nil
}
```

- [ ] **Step 2: Run tests to verify failure**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/connectors/application -run "AccountProber|FeeQuoteReader" -count=1
```

Expected: FAIL because ports and service methods do not exist.

- [ ] **Step 3: Add connector domain snapshots**

Modify `apps/server_core/internal/modules/connectors/domain/capability.go`:

```go
type AccountSnapshot struct {
	ProviderCode        string
	ProviderAccountID   string
	ProviderAccountName string
	SiteID              string
	Status              string
	FetchedAt           time.Time
	RawProviderRef      any
}

type FeeQuoteInput struct {
	AccountRef    ProviderAccountRef
	SiteID        string
	CategoryID    string
	ListingTypeID string
	PriceAmount   float64
	CurrencyID    string
}

type FeeQuoteSnapshot struct {
	ProviderCode      string
	SiteID            string
	CategoryID        string
	ListingTypeID     string
	PriceAmount       float64
	CurrencyID        string
	CommissionPercent *float64
	FixedFeeAmount    *float64
	SourceUpdatedAt   *time.Time
	FetchedAt         time.Time
	RawProviderRef    any
}
```

- [ ] **Step 4: Add connector ports and constants**

Modify `apps/server_core/internal/modules/connectors/ports/marketplace_capability.go`:

```go
const (
	CapabilityAccountProbe = "account_probe"
	CapabilityFeeQuoteRead = "fee_quote_read"
)

type AccountProber interface {
	ProbeAccount(ctx context.Context, ref domain.ProviderAccountRef) (domain.AccountSnapshot, error)
}

type FeeQuoteReader interface {
	ReadFeeQuote(ctx context.Context, input domain.FeeQuoteInput) (domain.FeeQuoteSnapshot, error)
}
```

- [ ] **Step 5: Add service fields and getters**

Modify `ProviderCapabilitySet` in `marketplace_capability_service.go`:

```go
AccountProbes ports.AccountProber
FeeQuotes     ports.FeeQuoteReader
```

Add methods:

```go
func (s *MarketplaceCapabilityService) AccountProber(providerCode string) (ports.AccountProber, error) {
	capability, err := s.provider(providerCode)
	if err != nil {
		return nil, err
	}
	if capability.AccountProbes == nil {
		return nil, unsupported(providerCode, ports.CapabilityAccountProbe)
	}
	return capability.AccountProbes, nil
}

func (s *MarketplaceCapabilityService) FeeQuoteReader(providerCode string) (ports.FeeQuoteReader, error) {
	capability, err := s.provider(providerCode)
	if err != nil {
		return nil, err
	}
	if capability.FeeQuotes == nil {
		return nil, unsupported(providerCode, ports.CapabilityFeeQuoteRead)
	}
	return capability.FeeQuotes, nil
}
```

- [ ] **Step 6: Run connector application tests**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/connectors/application ./apps/server_core/internal/modules/connectors/domain -count=1
```

Expected: PASS.

- [ ] **Step 7: Commit connector ports**

Run:

```powershell
git add -- apps/server_core/internal/modules/connectors/domain/capability.go apps/server_core/internal/modules/connectors/ports/marketplace_capability.go apps/server_core/internal/modules/connectors/application/marketplace_capability_service.go apps/server_core/internal/modules/connectors/application/marketplace_capability_service_test.go
git commit -m "feat(connectors): add read probe capability ports"
```

Expected: commit only task 6 files.

---

## Task 7: Implement Mercado Livre Read Probe Capabilities

**Files:**
- Modify: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go`
- Modify: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter_test.go`
- Modify: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/fee_sync.go`

**Interfaces:**
- Consumes: `ports.AccountProber`, `ports.FeeQuoteReader`.
- Produces:
  - `CapabilityAdapter.ProbeAccount`
  - `CapabilityAdapter.ReadFeeQuote`
  - `ProviderCapabilitySet` without `StockWrites` for slice 1 runtime.

- [ ] **Step 1: Write failing account probe test**

Add to `capability_adapter_test.go`:

```go
func TestCapabilityAdapterProbeAccountReadsUsersMe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/me" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token-1" {
			t.Fatalf("authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":691607102,"nickname":"METALNOBREACABAMENTOS","site_id":"MLB","status":{"site_status":"active"}}`))
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL: server.URL,
		AccessTokenResolver: func(ctx context.Context, ref domain.ProviderAccountRef) (string, error) {
			return "token-1", nil
		},
		Now: func() time.Time { return time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC) },
	})

	snapshot, err := adapter.ProbeAccount(context.Background(), domain.ProviderAccountRef{
		TenantID:          "tenant_default",
		InstallationID:    "inst-1",
		ProviderCode:      "mercado_livre",
		ProviderAccountID: "691607102",
	})

	if err != nil {
		t.Fatalf("ProbeAccount error = %v", err)
	}
	if snapshot.ProviderAccountID != "691607102" {
		t.Fatalf("account id = %q", snapshot.ProviderAccountID)
	}
	if snapshot.ProviderAccountName != "METALNOBREACABAMENTOS" {
		t.Fatalf("account name = %q", snapshot.ProviderAccountName)
	}
	if snapshot.Status != "active" {
		t.Fatalf("status = %q", snapshot.Status)
	}
}
```

- [ ] **Step 2: Write failing fee quote test**

Add to `capability_adapter_test.go`:

```go
func TestCapabilityAdapterReadFeeQuoteReadsListingPrices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sites/MLB/listing_prices" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		query := r.URL.Query()
		if query.Get("price") != "100" || query.Get("listing_type_id") != "gold_special" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"listing_type_id":"gold_special","sale_fee_details":{"percentage_fee":11,"fixed_fee":6}}]`))
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL: server.URL,
		SiteID:  "MLB",
		AccessTokenResolver: func(ctx context.Context, ref domain.ProviderAccountRef) (string, error) {
			return "token-1", nil
		},
		Now: func() time.Time { return time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC) },
	})

	quote, err := adapter.ReadFeeQuote(context.Background(), domain.FeeQuoteInput{
		AccountRef: domain.ProviderAccountRef{
			TenantID:          "tenant_default",
			InstallationID:    "inst-1",
			ProviderCode:      "mercado_livre",
			ProviderAccountID: "691607102",
		},
		SiteID:        "MLB",
		ListingTypeID: "gold_special",
		PriceAmount:   100,
		CurrencyID:    "BRL",
	})

	if err != nil {
		t.Fatalf("ReadFeeQuote error = %v", err)
	}
	if quote.CommissionPercent == nil || *quote.CommissionPercent != 11 {
		t.Fatalf("commission percent = %v", quote.CommissionPercent)
	}
	if quote.FixedFeeAmount == nil || *quote.FixedFeeAmount != 6 {
		t.Fatalf("fixed fee = %v", quote.FixedFeeAmount)
	}
}
```

- [ ] **Step 3: Run adapter tests to verify failure**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/connectors/adapters/mercado_livre -run "ProbeAccount|ReadFeeQuote" -count=1
```

Expected: FAIL because methods do not exist.

- [ ] **Step 4: Remove stock write from runtime set**

Modify `ProviderCapabilitySet()` in `capability_adapter.go`:

```go
func (a *CapabilityAdapter) ProviderCapabilitySet() connectorsapp.ProviderCapabilitySet {
	return connectorsapp.ProviderCapabilitySet{
		ProviderCode:   "mercado_livre",
		AccountProbes:  a,
		Listings:       a,
		StockReads:     a,
		Orders:         a,
		FeeQuotes:      a,
	}
}
```

Do not delete `UpdateAvailableQuantity`; it remains modeled for future slice but not registered as runtime capability.

- [ ] **Step 5: Implement ProbeAccount**

Add to `capability_adapter.go`:

```go
func (a *CapabilityAdapter) ProbeAccount(ctx context.Context, ref domain.ProviderAccountRef) (domain.AccountSnapshot, error) {
	accountRef, err := normalizeAccountRef(ref)
	if err != nil {
		return domain.AccountSnapshot{}, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.AccountSnapshot{}, err
	}

	var response mlAccountResponse
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/users/me", nil, &response); err != nil {
		return domain.AccountSnapshot{}, err
	}

	return domain.AccountSnapshot{
		ProviderCode:        "mercado_livre",
		ProviderAccountID:   normalizeAnyID(response.ID),
		ProviderAccountName: firstNonEmpty(response.Nickname, response.FirstName+" "+response.LastName),
		SiteID:              response.SiteID,
		Status:              firstNonEmpty(response.Status.SiteStatus, response.Status.List.Status),
		FetchedAt:           a.now(),
		RawProviderRef:      "/users/me",
	}, nil
}
```

Add private response structs:

```go
type mlAccountResponse struct {
	ID        any    `json:"id"`
	Nickname  string `json:"nickname"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	SiteID    string `json:"site_id"`
	Status    struct {
		SiteStatus string `json:"site_status"`
		List       struct {
			Status string `json:"status"`
		} `json:"list"`
	} `json:"status"`
}
```

- [ ] **Step 6: Implement ReadFeeQuote**

Add to `capability_adapter.go`:

```go
func (a *CapabilityAdapter) ReadFeeQuote(ctx context.Context, input domain.FeeQuoteInput) (domain.FeeQuoteSnapshot, error) {
	accountRef, err := normalizeAccountRef(input.AccountRef)
	if err != nil {
		return domain.FeeQuoteSnapshot{}, err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.FeeQuoteSnapshot{}, err
	}
	if input.PriceAmount <= 0 {
		return domain.FeeQuoteSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderValidation, "price amount must be positive")
	}
	listingTypeID := strings.TrimSpace(input.ListingTypeID)
	if listingTypeID == "" {
		return domain.FeeQuoteSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "listing type id is required")
	}
	siteID := firstNonEmpty(input.SiteID, a.siteID)

	query := url.Values{}
	query.Set("price", strconv.FormatFloat(input.PriceAmount, 'f', -1, 64))
	query.Set("listing_type_id", listingTypeID)
	if categoryID := strings.TrimSpace(input.CategoryID); categoryID != "" {
		query.Set("category_id", categoryID)
	}

	var response []mlListingPriceResponse
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/sites/"+url.PathEscape(siteID)+"/listing_prices?"+query.Encode(), nil, &response); err != nil {
		return domain.FeeQuoteSnapshot{}, err
	}
	if len(response) == 0 {
		return domain.FeeQuoteSnapshot{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider returned no fee quote")
	}
	price := response[0]
	return domain.FeeQuoteSnapshot{
		ProviderCode:      "mercado_livre",
		SiteID:            siteID,
		CategoryID:        strings.TrimSpace(input.CategoryID),
		ListingTypeID:     firstNonEmpty(price.ListingTypeID, listingTypeID),
		PriceAmount:       input.PriceAmount,
		CurrencyID:        firstNonEmpty(input.CurrencyID, "BRL"),
		CommissionPercent: price.SaleFeeDetails.PercentageFee,
		FixedFeeAmount:    price.SaleFeeDetails.FixedFee,
		FetchedAt:         a.now(),
		RawProviderRef:    "/sites/" + siteID + "/listing_prices",
	}, nil
}
```

Add private response struct:

```go
type mlListingPriceResponse struct {
	ListingTypeID  string `json:"listing_type_id"`
	SaleFeeDetails struct {
		PercentageFee *float64 `json:"percentage_fee"`
		FixedFee      *float64 `json:"fixed_fee"`
	} `json:"sale_fee_details"`
}
```

- [ ] **Step 7: Add interface assertions**

At the bottom of `capability_adapter.go`, add:

```go
_ ports.AccountProber  = (*CapabilityAdapter)(nil)
_ ports.FeeQuoteReader = (*CapabilityAdapter)(nil)
```

- [ ] **Step 8: Demote static fee sync**

Inspect `apps/server_core/internal/modules/connectors/adapters/mercado_livre/fee_sync.go`.

If it only seeds static rows, rename metadata/comment text that claims API sync to seed semantics and ensure provider runtime capability does not expose `pricing_fee_sync` as live. Do not claim live fee sync in slice 1; live fee evidence comes from `ReadFeeQuote`.

- [ ] **Step 9: Run Mercado Livre adapter tests**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/connectors/adapters/mercado_livre -count=1
```

Expected: PASS.

- [ ] **Step 10: Commit Mercado Livre reads**

Run:

```powershell
git add -- apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter_test.go apps/server_core/internal/modules/connectors/adapters/mercado_livre/fee_sync.go
git commit -m "feat(connectors): add mercado livre read probes"
```

Expected: commit only task 7 files.

---

## Task 8: Wire Connector Runtime Through Composition

**Files:**
- Modify: `apps/server_core/internal/composition/root.go`
- Create: `apps/server_core/internal/modules/integrations/application/provider_operation_service.go`
- Create: `apps/server_core/internal/modules/integrations/application/provider_operation_service_test.go`
- Modify: `apps/server_core/internal/modules/integrations/domain/operation_run.go`
- Modify: `apps/server_core/internal/modules/integrations/application/operation_service.go`

**Interfaces:**
- Consumes:
  - `CredentialResolver.ResolveAccessToken`
  - `connectors/application.MarketplaceCapabilityService`
- Produces:
  - provider operation application service for account/listing/order/fee reads
  - operation run begin/finish around connector execution

- [ ] **Step 1: Write failing provider operation test**

Create `apps/server_core/internal/modules/integrations/application/provider_operation_service_test.go`:

```go
package application

import (
	"context"
	"testing"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

func TestProviderOperationServiceRejectsUnavailableCapability(t *testing.T) {
	service := ProviderOperationService{
		tenantID: "tenant_default",
		installations: fakeProviderOperationInstallations{
			inst: domain.Installation{
				InstallationID: "inst-1",
				ProviderCode:   "mercado_livre",
				Status:         domain.InstallationStatusConnected,
				HealthStatus:   domain.HealthStatusHealthy,
				ExternalAccountID: "691607102",
				RuntimeCapabilities: []domain.RuntimeCapability{{
					Code:       domain.RuntimeCapabilityStockWrite,
					State:      domain.RuntimeCapabilityStateUnavailable,
					Executable: false,
				}},
			},
			found: true,
		},
	}

	_, err := service.ProbeAccount(context.Background(), "inst-1")

	if err == nil || err.Error() != "INTEGRATIONS_CAPABILITY_UNAVAILABLE" {
		t.Fatalf("err = %v", err)
	}
}

func TestProviderOperationServiceProbesAccount(t *testing.T) {
	now := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	prober := fakeProviderOperationAccountProber{snapshot: connectorsdomain.AccountSnapshot{
		ProviderCode:        "mercado_livre",
		ProviderAccountID:   "691607102",
		ProviderAccountName: "METALNOBREACABAMENTOS",
		FetchedAt:           now,
	}}
	service := newProviderOperationServiceForTest(prober)

	result, err := service.ProbeAccount(context.Background(), "inst-1")

	if err != nil {
		t.Fatalf("ProbeAccount error = %v", err)
	}
	if result.ProviderAccountName != "METALNOBREACABAMENTOS" {
		t.Fatalf("account name = %q", result.ProviderAccountName)
	}
}
```

Use compact fakes in the same file. The service must reject missing runtime capability before calling connectors.

- [ ] **Step 2: Run test to verify failure**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/integrations/application -run TestProviderOperationService -count=1
```

Expected: FAIL because service does not exist.

- [ ] **Step 3: Implement provider operation service**

Create `apps/server_core/internal/modules/integrations/application/provider_operation_service.go`:

```go
package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type providerOperationInstallationReader interface {
	Get(ctx context.Context, installationID string) (domain.Installation, bool, error)
}

type ProviderOperationService struct {
	tenantID       string
	installations providerOperationInstallationReader
	capabilities  *connectorsapp.MarketplaceCapabilityService
	operations    *OperationService
}

type ProviderOperationServiceConfig struct {
	TenantID       string
	Installations providerOperationInstallationReader
	Capabilities  *connectorsapp.MarketplaceCapabilityService
	Operations    *OperationService
}

func NewProviderOperationService(cfg ProviderOperationServiceConfig) *ProviderOperationService {
	tenantID := strings.TrimSpace(cfg.TenantID)
	if tenantID == "" {
		tenantID = "tenant_default"
	}
	return &ProviderOperationService{
		tenantID:       tenantID,
		installations: cfg.Installations,
		capabilities:  cfg.Capabilities,
		operations:    cfg.Operations,
	}
}

func (s *ProviderOperationService) ProbeAccount(ctx context.Context, installationID string) (connectorsdomain.AccountSnapshot, error) {
	inst, err := s.loadExecutableInstallation(ctx, installationID, domain.RuntimeCapabilityAccountProbe)
	if err != nil {
		return connectorsdomain.AccountSnapshot{}, err
	}
	prober, err := s.capabilities.AccountProber(inst.ProviderCode)
	if err != nil {
		return connectorsdomain.AccountSnapshot{}, err
	}
	ref := connectorsdomain.ProviderAccountRef{
		TenantID:          s.tenantID,
		InstallationID:    inst.InstallationID,
		ProviderCode:      inst.ProviderCode,
		ProviderAccountID: inst.ExternalAccountID,
	}
	return prober.ProbeAccount(ctx, ref)
}

func (s *ProviderOperationService) loadExecutableInstallation(ctx context.Context, installationID string, capability domain.RuntimeCapabilityCode) (domain.Installation, error) {
	if s.installations == nil || s.capabilities == nil {
		return domain.Installation{}, errors.New("INTEGRATIONS_PROVIDER_OPERATION_NOT_CONFIGURED")
	}
	inst, found, err := s.installations.Get(ctx, strings.TrimSpace(installationID))
	if err != nil {
		return domain.Installation{}, err
	}
	if !found {
		return domain.Installation{}, domain.ErrInstallationNotFound
	}
	for _, runtimeCapability := range inst.RuntimeCapabilities {
		if runtimeCapability.Code == capability && runtimeCapability.Available() {
			return inst, nil
		}
	}
	return domain.Installation{}, errors.New("INTEGRATIONS_CAPABILITY_UNAVAILABLE")
}

func operationRunID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UTC().UnixNano())
}
```

- [ ] **Step 4: Add listing/order/fee operation methods**

Add these methods to `provider_operation_service.go`:

```go
func (s *ProviderOperationService) ListListings(ctx context.Context, installationID string, limit int) ([]connectorsdomain.ListingSnapshot, error) {
	inst, err := s.loadExecutableInstallation(ctx, installationID, domain.RuntimeCapabilityListingRead)
	if err != nil {
		return nil, err
	}
	reader, err := s.capabilities.ListingReader(inst.ProviderCode)
	if err != nil {
		return nil, err
	}
	return reader.ListListings(ctx, connectorsdomain.ListListingsInput{
		AccountRef: connectorsdomain.ProviderAccountRef{
			TenantID:          s.tenantID,
			InstallationID:    inst.InstallationID,
			ProviderCode:      inst.ProviderCode,
			ProviderAccountID: inst.ExternalAccountID,
		},
		Limit: limit,
	})
}

func (s *ProviderOperationService) ListOrders(ctx context.Context, installationID string, limit int) ([]connectorsdomain.OrderSnapshot, error) {
	inst, err := s.loadExecutableInstallation(ctx, installationID, domain.RuntimeCapabilityOrderRead)
	if err != nil {
		return nil, err
	}
	reader, err := s.capabilities.OrderReader(inst.ProviderCode)
	if err != nil {
		return nil, err
	}
	return reader.ListOrders(ctx, connectorsdomain.ListOrdersInput{
		AccountRef: connectorsdomain.ProviderAccountRef{
			TenantID:          s.tenantID,
			InstallationID:    inst.InstallationID,
			ProviderCode:      inst.ProviderCode,
			ProviderAccountID: inst.ExternalAccountID,
		},
		Limit: limit,
	})
}

func (s *ProviderOperationService) ReadFeeQuote(ctx context.Context, installationID string, input connectorsdomain.FeeQuoteInput) (connectorsdomain.FeeQuoteSnapshot, error) {
	inst, err := s.loadExecutableInstallation(ctx, installationID, domain.RuntimeCapabilityFeeQuoteRead)
	if err != nil {
		return connectorsdomain.FeeQuoteSnapshot{}, err
	}
	reader, err := s.capabilities.FeeQuoteReader(inst.ProviderCode)
	if err != nil {
		return connectorsdomain.FeeQuoteSnapshot{}, err
	}
	input.AccountRef = connectorsdomain.ProviderAccountRef{
		TenantID:          s.tenantID,
		InstallationID:    inst.InstallationID,
		ProviderCode:      inst.ProviderCode,
		ProviderAccountID: inst.ExternalAccountID,
	}
	return reader.ReadFeeQuote(ctx, input)
}
```

- [ ] **Step 5: Wire Mercado Livre adapter in composition**

Modify imports in `apps/server_core/internal/composition/root.go`:

```go
mercadolivreconnector "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
```

Remove or keep the blank import only if it is still needed for fee seed side effects. Prefer explicit construction for capability execution.

After `credentialSvc`, `encryptionSvc`, and `authFlowSvc` are available, create resolver and adapter:

```go
credentialResolver := integrationsapp.NewCredentialResolver(credentialSvc, encryptionSvc)
mercadoLivreCapabilities := mercadolivreconnector.NewCapabilityAdapter(mercadolivreconnector.CapabilityAdapterConfig{
	AccessTokenResolver: func(ctx context.Context, ref connectorsdomain.ProviderAccountRef) (string, error) {
		resolved, err := credentialResolver.ResolveAccessToken(ctx, integrationsapp.CredentialResolutionRef{
			TenantID:          ref.TenantID,
			InstallationID:    ref.InstallationID,
			ProviderCode:      ref.ProviderCode,
			ProviderAccountID: ref.ProviderAccountID,
		})
		if err != nil {
			return "", err
		}
		return resolved.AccessToken, nil
	},
})
marketplaceCapabilities := connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{
	mercadoLivreCapabilities.ProviderCapabilitySet(),
})
```

Add `connectorsdomain` import:

```go
connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
```

Pass `marketplaceCapabilities` into `NewProviderOperationService`.

- [ ] **Step 6: Run composition build**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/composition ./apps/server_core/internal/modules/integrations/application -count=1
```

Expected: PASS.

- [ ] **Step 7: Commit runtime wiring**

Run:

```powershell
git add -- apps/server_core/internal/composition/root.go apps/server_core/internal/modules/integrations/application/provider_operation_service.go apps/server_core/internal/modules/integrations/application/provider_operation_service_test.go apps/server_core/internal/modules/integrations/domain/operation_run.go apps/server_core/internal/modules/integrations/application/operation_service.go
git commit -m "feat(integrations): wire provider read operations"
```

Expected: commit only task 8 files.

---

## Task 9: Align OpenAPI And SDK Runtime

**Files:**
- Modify: `contracts/api/marketplace-central.openapi.yaml`
- Modify: `packages/sdk-runtime/src/index.ts`
- Modify: `packages/sdk-runtime/src/index.test.ts`

**Interfaces:**
- Consumes: backend installation/auth/operation JSON shape.
- Produces:
  - `IntegrationConnectionSnapshot`
  - `IntegrationRuntimeCapability`
  - updated `IntegrationInstallation`
  - updated `IntegrationAuthStatusResponse`

- [ ] **Step 1: Update OpenAPI schemas**

In `contracts/api/marketplace-central.openapi.yaml`, add schemas:

```yaml
    IntegrationConnectionSnapshot:
      type: object
      required:
        - state
        - health
        - provider_code
        - external_account_id
        - external_account_name
        - auth_strategy
        - next_action
      properties:
        state:
          type: string
          enum: [draft, pending_connection, connected, degraded, needs_reauth, disconnected]
        health:
          type: string
          enum: [healthy, warning, critical]
        provider_code:
          type: string
        external_account_id:
          type: string
        external_account_name:
          type: string
        auth_strategy:
          type: string
          enum: [oauth2, api_key, none]
        last_verified_at:
          type: string
          format: date-time
          nullable: true
        expires_at:
          type: string
          format: date-time
          nullable: true
        next_action:
          type: string
          enum: [none, authorize, reauth, configure, retry]
        reauth_reason:
          type: string

    IntegrationRuntimeCapability:
      type: object
      required:
        - code
        - state
        - executable
        - live_validated
        - local_validated
      properties:
        code:
          type: string
        state:
          type: string
          enum: [available, unavailable, needs_auth, degraded, not_configured]
        executable:
          type: boolean
        live_validated:
          type: boolean
        local_validated:
          type: boolean
        unavailable_reason:
          type: string
        last_validated_at:
          type: string
          format: date-time
          nullable: true
```

Update `IntegrationInstallation` to include:

```yaml
        connection:
          $ref: '#/components/schemas/IntegrationConnectionSnapshot'
        runtime_capabilities:
          type: array
          items:
            $ref: '#/components/schemas/IntegrationRuntimeCapability'
```

Update `IntegrationAuthStatusResponse` to include:

```yaml
        connection:
          $ref: '#/components/schemas/IntegrationConnectionSnapshot'
```

- [ ] **Step 2: Update SDK types**

In `packages/sdk-runtime/src/index.ts`, add:

```ts
export interface IntegrationConnectionSnapshot {
  state: "draft" | "pending_connection" | "connected" | "degraded" | "needs_reauth" | "disconnected";
  health: "healthy" | "warning" | "critical";
  provider_code: string;
  external_account_id: string;
  external_account_name: string;
  auth_strategy: "oauth2" | "api_key" | "none";
  last_verified_at?: string | null;
  expires_at?: string | null;
  next_action: "none" | "authorize" | "reauth" | "configure" | "retry";
  reauth_reason?: string;
}

export interface IntegrationRuntimeCapability {
  code: string;
  state: "available" | "unavailable" | "needs_auth" | "degraded" | "not_configured";
  executable: boolean;
  live_validated: boolean;
  local_validated: boolean;
  unavailable_reason?: string;
  last_validated_at?: string | null;
}
```

Update `IntegrationInstallation`:

```ts
  connection: IntegrationConnectionSnapshot;
  runtime_capabilities: IntegrationRuntimeCapability[];
```

Update `IntegrationAuthStatusResponse`:

```ts
  connection: IntegrationConnectionSnapshot;
```

Keep existing flattened fields during this slice to avoid breaking older UI/tests. Mark them as compatibility in comments if the file already uses comments for deprecation.

- [ ] **Step 3: Add SDK test fixture**

Add to `packages/sdk-runtime/src/index.test.ts`:

```ts
it("lists integration installations with connection snapshot", async () => {
  const fetchImpl = vi.fn().mockResolvedValue({
    ok: true,
    json: async () => ({
      items: [{
        installation_id: "inst-1",
        tenant_id: "tenant_default",
        provider_code: "mercado_livre",
        family: "marketplace",
        display_name: "Mercado Livre",
        status: "connected",
        health_status: "healthy",
        external_account_id: "691607102",
        external_account_name: "METALNOBREACABAMENTOS",
        connection: {
          state: "connected",
          health: "healthy",
          provider_code: "mercado_livre",
          external_account_id: "691607102",
          external_account_name: "METALNOBREACABAMENTOS",
          auth_strategy: "oauth2",
          last_verified_at: "2026-07-08T12:00:00Z",
          expires_at: "2026-07-08T18:00:00Z",
          next_action: "none",
        },
        runtime_capabilities: [{
          code: "account_probe",
          state: "available",
          executable: true,
          live_validated: true,
          local_validated: true,
        }],
        created_at: "2026-07-08T12:00:00Z",
        updated_at: "2026-07-08T12:00:00Z",
      }],
    }),
  });
  const client = createMarketplaceCentralClient({ baseUrl: "http://api.test", fetchImpl: fetchImpl as unknown as typeof fetch });

  const result = await client.listIntegrationInstallations();

  expect(result.items[0].connection.external_account_name).toBe("METALNOBREACABAMENTOS");
  expect(result.items[0].runtime_capabilities[0].code).toBe("account_probe");
});
```

- [ ] **Step 4: Run SDK tests**

Run:

```powershell
npm test -- --run packages/sdk-runtime/src/index.test.ts
```

Expected: PASS.

- [ ] **Step 5: Commit API and SDK**

Run:

```powershell
git add -- contracts/api/marketplace-central.openapi.yaml packages/sdk-runtime/src/index.ts packages/sdk-runtime/src/index.test.ts
git commit -m "feat(api): expose integration connection snapshot"
```

Expected: commit only task 9 files.

---

## Task 10: Update Integrations UI To Connection-First

**Files:**
- Modify: `packages/feature-integrations/src/components/InstallationCard.tsx`
- Modify: `packages/feature-integrations/src/components/InstallationDrawer.tsx`
- Modify: `packages/feature-integrations/src/components/AuthStatusPanel.tsx`
- Modify: `packages/feature-integrations/src/components/OperationsTimeline.tsx`
- Modify: `packages/feature-integrations/src/IntegrationsHubPage.tsx`
- Modify: `packages/feature-integrations/src/IntegrationsHubPage.test.tsx`

**Interfaces:**
- Consumes: SDK `IntegrationConnectionSnapshot` and `IntegrationRuntimeCapability`.
- Produces: UI that never labels a connected installation as credential-disconnected because credential id is missing.

- [ ] **Step 1: Update test fixtures first**

In `IntegrationsHubPage.test.tsx`, add a helper:

```ts
function connectedSnapshot(overrides: Partial<IntegrationInstallation["connection"]> = {}): IntegrationInstallation["connection"] {
  return {
    state: "connected",
    health: "healthy",
    provider_code: "mercado_livre",
    external_account_id: "691607102",
    external_account_name: "METALNOBREACABAMENTOS",
    auth_strategy: "oauth2",
    last_verified_at: "2026-07-08T12:00:00Z",
    expires_at: "2026-07-08T18:00:00Z",
    next_action: "none",
    ...overrides,
  };
}

function runtimeCapabilities(): IntegrationInstallation["runtime_capabilities"] {
  return [
    { code: "account_probe", state: "available", executable: true, live_validated: true, local_validated: true },
    { code: "listing_read", state: "available", executable: true, live_validated: true, local_validated: true },
    { code: "order_read", state: "available", executable: true, live_validated: true, local_validated: true },
  ];
}
```

Update installation fixtures to include:

```ts
connection: connectedSnapshot(),
runtime_capabilities: runtimeCapabilities(),
```

- [ ] **Step 2: Add failing UI assertion**

Add a test:

```ts
it("does not render missing credential as disconnected when connection is healthy", async () => {
  mockListInstallations.mockResolvedValue({
    items: [{
      ...sampleInstallation,
      active_credential_id: undefined,
      connection: connectedSnapshot(),
      runtime_capabilities: runtimeCapabilities(),
    }],
  });

  render(<IntegrationsHubPage client={mockClient} />);

  expect(await screen.findByText("METALNOBREACABAMENTOS")).toBeInTheDocument();
  expect(screen.queryByText("Not connected")).not.toBeInTheDocument();
});
```

- [ ] **Step 3: Run UI test to verify failure**

Run:

```powershell
npm test -- --run packages/feature-integrations/src/IntegrationsHubPage.test.tsx
```

Expected: FAIL because `InstallationCard` and drawer still render credential state.

- [ ] **Step 4: Update InstallationCard**

Replace credential display with connection display:

```tsx
const connection = installation.connection;
const externalAccountLabel =
  connection.external_account_name || connection.external_account_id || "Account not verified";
```

Render:

```tsx
<dt className="text-xs uppercase tracking-wide text-slate-400">External account</dt>
<dd className="mt-1 text-slate-700">{externalAccountLabel}</dd>
```

Replace credential section with:

```tsx
<dt className="text-xs uppercase tracking-wide text-slate-400">Next action</dt>
<dd className="mt-1 text-slate-700">{connection.next_action}</dd>
```

Replace last verified with:

```tsx
Last verified {connection.last_verified_at ?? "not yet verified"}
```

- [ ] **Step 5: Update InstallationDrawer connection section**

Use:

```tsx
const connection = installation.connection;
```

Replace credential field with:

```tsx
<div>
  <dt className="text-xs uppercase tracking-wide text-slate-400">Next action</dt>
  <dd className="mt-1 text-slate-700">{connection.next_action}</dd>
</div>
<div>
  <dt className="text-xs uppercase tracking-wide text-slate-400">Auth</dt>
  <dd className="mt-1 text-slate-700">{connection.auth_strategy}</dd>
</div>
```

Add a capabilities section:

```tsx
<section className="rounded-2xl border border-slate-200 bg-white p-4">
  <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">Runtime capabilities</p>
  <ul className="mt-3 space-y-2">
    {installation.runtime_capabilities.map((capability) => (
      <li key={capability.code} className="flex items-center justify-between rounded-xl bg-slate-50 px-3 py-2 text-sm">
        <span className="font-medium text-slate-700">{capability.code}</span>
        <span className="text-xs text-slate-500">{capability.state}</span>
      </li>
    ))}
  </ul>
</section>
```

- [ ] **Step 6: Align AuthStatusPanel**

If `AuthStatusPanel` reads `authStatus.external_account_id`, prefer:

```tsx
const connection = authStatus?.connection;
```

Display account from:

```tsx
connection?.external_account_name || connection?.external_account_id || "Not verified yet"
```

Display next action from:

```tsx
connection?.next_action ?? "authorize"
```

- [ ] **Step 7: Run UI tests**

Run:

```powershell
npm test -- --run packages/feature-integrations/src/IntegrationsHubPage.test.tsx
```

Expected: PASS.

- [ ] **Step 8: Commit UI**

Run:

```powershell
git add -- packages/feature-integrations/src/components/InstallationCard.tsx packages/feature-integrations/src/components/InstallationDrawer.tsx packages/feature-integrations/src/components/AuthStatusPanel.tsx packages/feature-integrations/src/components/OperationsTimeline.tsx packages/feature-integrations/src/IntegrationsHubPage.tsx packages/feature-integrations/src/IntegrationsHubPage.test.tsx
git commit -m "feat(integrations): render connection-first UI"
```

Expected: commit only task 10 files.

---

## Task 11: End-To-End Local Verification

**Files:**
- No source files unless a test exposes a real defect.
- Modify only validation artifacts after commands complete.

**Interfaces:**
- Consumes: all previous tasks.
- Produces: local contract validation evidence.

- [ ] **Step 1: Run backend integration module tests**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/integrations/... -count=1
```

Expected: PASS.

- [ ] **Step 2: Run connector tests**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/connectors/... -count=1
```

Expected: PASS.

- [ ] **Step 3: Run composition build tests**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/composition -count=1
```

Expected: PASS.

- [ ] **Step 4: Run frontend package tests**

Run:

```powershell
npm test -- --run packages/sdk-runtime/src/index.test.ts packages/feature-integrations/src/IntegrationsHubPage.test.tsx
```

Expected: PASS.

- [ ] **Step 5: Run app builds if package scripts exist**

Run:

```powershell
npm run build --if-present
```

Expected: PASS or document exact package/script that does not exist. Do not claim build validation if the command does not run.

- [ ] **Step 6: Commit only fixes if verification exposed defects**

If code changes were required:

```powershell
git add -- <exact changed files>
git commit -m "fix(integrations): stabilize marketplace foundation"
```

Expected: no commit if no fixes were needed.

---

## Task 12: Live Mercado Livre Read Validation

**Files:**
- Create or modify: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02*/validation.md` if M2 is the active evidence target.
- Create or modify: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02*/validation-result.md` if M2 is the active evidence target.
- If M2 artifact path differs, use the actual M2 path discovered with `Get-ChildItem .mnfs -Recurse -Filter validation.md`.

**Interfaces:**
- Consumes: live ngrok/local app, active Mercado Livre credential, no write endpoints.
- Produces: live read evidence with timestamp, endpoint, account id, and result summary.

- [ ] **Step 1: Confirm server and ngrok URL**

Run:

```powershell
docker compose ps
```

Expected: frontend/backend/ngrok containers running. If not running, start with the repo's documented docker compose command.

- [ ] **Step 2: Confirm installation API shape**

Run:

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/integrations/installations" -Method Get | ConvertTo-Json -Depth 10
```

Expected: Mercado Livre installation includes `connection.external_account_id`, `connection.external_account_name`, `connection.next_action = "none"`, and runtime read/probe capabilities.

- [ ] **Step 3: Validate auth status mirrors connection**

Run with the real installation id:

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/integrations/installations/inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98/auth/status" -Method Get | ConvertTo-Json -Depth 10
```

Expected: status `connected`, health `healthy`, and the same connection snapshot as installation list.

- [ ] **Step 4: Run live provider read probes through platform endpoints**

Use the read/probe endpoints implemented in task 8. Example shape if endpoints are exposed under integrations:

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/integrations/installations/inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98/probes/account" -Method Post | ConvertTo-Json -Depth 10
Invoke-RestMethod -Uri "http://localhost:8080/integrations/installations/inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98/probes/listings?limit=3" -Method Get | ConvertTo-Json -Depth 10
Invoke-RestMethod -Uri "http://localhost:8080/integrations/installations/inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98/probes/orders?limit=1" -Method Get | ConvertTo-Json -Depth 10
Invoke-RestMethod -Uri "http://localhost:8080/integrations/installations/inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98/probes/fee-quote?listing_type_id=gold_special&price=100" -Method Get | ConvertTo-Json -Depth 10
```

Expected:

- account probe returns provider account `691607102` or the current connected account id.
- listings returns at least a successful provider response summary; zero listings is acceptable only if provider says zero.
- orders returns a successful provider response summary; zero orders is acceptable only if provider says zero.
- fee quote returns commission/fixed fee values or a typed provider unsupported/error state. Unknown is not converted to zero.

- [ ] **Step 5: Validate UI through ngrok**

Open:

```text
https://multiradial-unironically-nieves.ngrok-free.dev/integrations?installation=inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98
```

Expected:

- Connection shows account name/id.
- No healthy connected card says `Credential Not connected`.
- Runtime capabilities show only executable read/probe capabilities.
- Operation timeline records read/probe operations or existing fee sync operations with honest result/error status.

- [ ] **Step 6: Record live evidence**

In the active MNFS validation artifact, add:

```markdown
## Live Mercado Livre Read Validation

Date: 2026-07-08
Environment: local backend through active development credentials
Installation: inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98

Evidence:
- Installation connection snapshot returned connected/healthy with external account id/name.
- Auth status mirrored the same connection snapshot.
- Account probe returned provider account id/name from live Mercado Livre.
- Listings read returned live provider response summary.
- Orders read returned live provider response summary.
- Fee quote returned live provider response or typed provider error without defaulting unknown values to zero.

Boundary:
- No provider writes were executed.
- Mock/fake tests were used only for local contract behavior.
- Live validation covers read/probe behavior only.
```

- [ ] **Step 7: Commit validation evidence**

Run:

```powershell
git add -- .mnfs
git commit -m "test(integrations): record mercado livre live reads"
```

Expected: commit only validation artifact changes.

---

## Task 13: Final Closeout And Handoff

**Files:**
- Modify: `.brain/system-pulse.md` if project status changed.
- Modify: `.brain/roadmap.json` if milestone status changed.
- Modify: `wiki/README.md` or module wiki only if implementation revealed new architecture truth.
- Modify: `docs/superpowers/specs/2026-07-08-marketplace-integration-foundation-design.md` only if design truth changed.

**Interfaces:**
- Consumes: all implementation and validation evidence.
- Produces: final status that distinguishes local, integration, and live validation.

- [ ] **Step 1: Check final status**

Run:

```powershell
git status --short
```

Expected: only intentional closeout files are modified.

- [ ] **Step 2: Run final verification**

Run:

```powershell
$env:GOCACHE="$PWD\.gocache"; go test ./apps/server_core/internal/modules/integrations/... ./apps/server_core/internal/modules/connectors/... ./apps/server_core/internal/composition -count=1
npm test -- --run packages/sdk-runtime/src/index.test.ts packages/feature-integrations/src/IntegrationsHubPage.test.tsx
```

Expected: PASS. If any live-only step cannot run, document exact blocker and do not mark live validation passed.

- [ ] **Step 3: Update closeout notes**

Add a concise note to `.brain/system-pulse.md`:

```markdown
- Marketplace integration foundation slice 1 implemented: installation connection snapshots, runtime read/probe capability governance, credential resolver boundary, Mercado Livre read/probe wiring, API/SDK/UI alignment. Validation separates local tests from live Mercado Livre read evidence. Provider writes remain deferred.
```

Only add this after implementation and validation are actually complete.

- [ ] **Step 4: Commit closeout**

Run:

```powershell
git add -- .brain/system-pulse.md .brain/roadmap.json wiki/README.md docs/superpowers/specs/2026-07-08-marketplace-integration-foundation-design.md
git commit -m "docs(integrations): close marketplace foundation slice"
```

Expected: commit only files that truly changed. If only `.brain/system-pulse.md` changed, stage only that file.

- [ ] **Step 5: Final response checklist**

Final response must include:

```text
Implemented:
- Connection snapshot projection
- Runtime capability honesty
- Credential resolver boundary
- Mercado Livre read/probe wiring
- API/SDK/UI alignment

Validated:
- Local Go tests: command + pass/fail
- Local frontend tests: command + pass/fail
- Live Mercado Livre reads: endpoint summaries + pass/fail

Not validated:
- Provider writes, because slice 1 intentionally excludes writes.
```

Do not say M2/Mercado Livre is fully validated for writes.

---

## Self-Review

Spec coverage:

- Target architecture is covered by tasks 2, 5, 6, and 8.
- Auth projection is covered by tasks 3 and 4.
- Credential safety is covered by task 5.
- Mercado Livre read/probe execution is covered by tasks 6, 7, 8, and 12.
- API/SDK/UI alignment is covered by tasks 9 and 10.
- Validation boundaries are covered by tasks 11, 12, and 13.
- Deferred writes are explicitly excluded in tasks 7, 8, 10, 12, and 13.

Placeholder scan:

- No `TBD`, `TODO`, or unspecified "add tests" steps remain.
- Where helper names may differ in existing tests, the expected assertions and implementation shape are explicit.

Type consistency:

- `ConnectionSnapshot`, `RuntimeCapability`, `CredentialResolver`, `AccountProber`, and `FeeQuoteReader` are introduced before consumers use them.
- Runtime capability names match the approved spec: `account_probe`, `listing_read`, `order_read`, `fee_quote_read`, `stock_read`, and deferred `stock_write`.
