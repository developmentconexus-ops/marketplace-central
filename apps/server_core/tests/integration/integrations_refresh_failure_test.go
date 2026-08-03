//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// Proves, against real Postgres, the two things the fakes in the unit tests
// cannot: that the refresh failure survives the round trip with the shape the
// domain wrote, and that ListExpiringSessions honours next_retry_at. The second
// is what makes the backoff exist at all — if that WHERE clause did not filter,
// the ticker would retry the dead account on every tick and consecutive_failures
// would climb forever.
func TestRefreshFailurePersistsAndSuppressesTheSweep(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_harness_refresh_failure")
	testpostgres.SeedProvider(t, pool, testpostgres.ProviderFixture{Code: "mercado_livre", DisplayName: "Mercado Livre"})

	installationRepo := integrationspostgres.NewInstallationRepository(pool, cfg.DefaultTenantID)
	sessionRepo := integrationspostgres.NewAuthSessionRepository(pool, cfg.DefaultTenantID)

	ctx := context.Background()
	installationID := fmt.Sprintf("inst-refresh-fail-%d", time.Now().UTC().UnixNano())
	now := time.Now().UTC()

	t.Cleanup(func() {
		// Report the errors: a silent `_, _ =` here would hide a wrong table
		// name and leak rows into every later run of the lane.
		if _, err := pool.Exec(ctx, `DELETE FROM integration_auth_sessions WHERE installation_id = $1`, installationID); err != nil {
			t.Errorf("cleanup integration_auth_sessions: %v", err)
		}
		if _, err := pool.Exec(ctx, `DELETE FROM integration_installations WHERE installation_id = $1`, installationID); err != nil {
			t.Errorf("cleanup integration_installations: %v", err)
		}
	})

	if err := installationRepo.CreateInstallation(ctx, integrationsdomain.Installation{
		InstallationID: installationID,
		TenantID:       cfg.DefaultTenantID,
		ProviderCode:   "mercado_livre",
		Family:         integrationsdomain.IntegrationFamilyMarketplace,
		DisplayName:    "Mercado Livre (teste)",
		Status:         integrationsdomain.InstallationStatusConnected,
		HealthStatus:   integrationsdomain.HealthStatusHealthy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}); err != nil {
		t.Fatalf("CreateInstallation: %v", err)
	}

	// --- 1. The failure persists with the shape the domain wrote -------------
	// access_token_expires_at in the past: without next_retry_at this session
	// QUALIFIES for the sweep. That qualification is what makes step 3 a proof.
	expiredAt := now.Add(-time.Minute)
	nextRetryAt := now.Add(time.Hour)
	failureCode := integrationsdomain.ErrRefreshTokenInvalid.Error()

	if err := sessionRepo.UpsertAuthSession(ctx, integrationsdomain.AuthSession{
		AuthSessionID:        "auth_" + installationID,
		TenantID:             cfg.DefaultTenantID,
		InstallationID:       installationID,
		State:                integrationsdomain.AuthStateRefreshFailed,
		ProviderAccountID:    "seller-1",
		AccessTokenExpiresAt: &expiredAt,
		NextRetryAt:          &nextRetryAt,
		RefreshFailureCode:   failureCode,
		ConsecutiveFailures:  1,
		CreatedAt:            now,
		UpdatedAt:            now,
	}); err != nil {
		t.Fatalf("UpsertAuthSession: %v", err)
	}

	stored, found, err := sessionRepo.GetAuthSession(ctx, installationID)
	if err != nil || !found {
		t.Fatalf("GetAuthSession found=%v err=%v", found, err)
	}
	if stored.State != integrationsdomain.AuthStateRefreshFailed {
		t.Fatalf("State = %q, want refresh_failed (does the 0016 CHECK accept that value?)", stored.State)
	}
	if stored.RefreshFailureCode != failureCode {
		t.Fatalf("RefreshFailureCode = %q, want %q", stored.RefreshFailureCode, failureCode)
	}
	if stored.ConsecutiveFailures != 1 {
		t.Fatalf("ConsecutiveFailures = %d, want 1", stored.ConsecutiveFailures)
	}
	if stored.NextRetryAt == nil {
		t.Fatal("NextRetryAt = nil after the round trip")
	}

	// --- 2. Positive control: with no backoff, the sweep DOES pick it up -----
	// Without this step, step 3 would be vacuous: an empty list for any other
	// reason (wrong tenant, wrong expiry) would look like a working backoff.
	clearRetry := integrationsdomain.AuthSession{
		AuthSessionID:        "auth_" + installationID,
		TenantID:             cfg.DefaultTenantID,
		InstallationID:       installationID,
		State:                integrationsdomain.AuthStateExpiring,
		ProviderAccountID:    "seller-1",
		AccessTokenExpiresAt: &expiredAt,
		NextRetryAt:          nil,
		RefreshFailureCode:   failureCode,
		ConsecutiveFailures:  1,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	if err := sessionRepo.UpsertAuthSession(ctx, clearRetry); err != nil {
		t.Fatalf("UpsertAuthSession (no backoff): %v", err)
	}
	if !sweepContains(t, sessionRepo, installationID) {
		t.Fatal("a varredura NAO pegou a sessao vencida sem next_retry_at: o teste do backoff seria vacuoso")
	}

	// --- 3. With next_retry_at in the future, the sweep skips it -------------
	withBackoff := clearRetry
	withBackoff.NextRetryAt = &nextRetryAt
	if err := sessionRepo.UpsertAuthSession(ctx, withBackoff); err != nil {
		t.Fatalf("UpsertAuthSession (with backoff): %v", err)
	}
	if sweepContains(t, sessionRepo, installationID) {
		t.Fatal("a varredura pegou a sessao com next_retry_at no futuro: o backoff nao existe")
	}

	// --- 4. The critical state reaches the database and comes back ----------
	snapshot := integrationsdomain.ProjectConnectionSnapshot(
		integrationsdomain.Installation{
			InstallationID: installationID,
			TenantID:       cfg.DefaultTenantID,
			ProviderCode:   "mercado_livre",
			Status:         integrationsdomain.InstallationStatusRequiresReauth,
			HealthStatus:   integrationsdomain.HealthStatusCritical,
		},
		integrationsdomain.AuthStrategyOAuth2,
		&expiredAt,
		failureCode,
	)
	if err := installationRepo.ApplyConnectionSnapshot(ctx, installationID, snapshot, ""); err != nil {
		t.Fatalf("ApplyConnectionSnapshot: %v", err)
	}

	reread, found, err := installationRepo.GetInstallation(ctx, installationID)
	if err != nil || !found {
		t.Fatalf("GetInstallation found=%v err=%v", found, err)
	}
	if reread.HealthStatus != integrationsdomain.HealthStatusCritical {
		t.Fatalf("HealthStatus = %q, want critical (first write of 'critical' in this schema)", reread.HealthStatus)
	}
	if reread.ConnectionSnapshot.NextAction != integrationsdomain.ConnectionNextActionReauth {
		t.Fatalf("NextAction = %q, want reauth", reread.ConnectionSnapshot.NextAction)
	}
	if reread.ConnectionSnapshot.ReauthReason != failureCode {
		t.Fatalf("ReauthReason = %q, want %q", reread.ConnectionSnapshot.ReauthReason, failureCode)
	}
}

func sweepContains(t *testing.T, repo *integrationspostgres.AuthSessionRepository, installationID string) bool {
	t.Helper()
	sessions, err := repo.ListExpiringSessions(context.Background(), 10*time.Minute)
	if err != nil {
		t.Fatalf("ListExpiringSessions: %v", err)
	}
	for _, s := range sessions {
		if s.InstallationID == installationID {
			return true
		}
	}
	return false
}
