//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// Prova contra Postgres real as duas coisas que os fakes das Tasks 1-3 não
// alcançam: que reautorizar a MESMA installation não viola o índice parcial, e
// que uma installation 'disconnected' de fato libera a conta — que é o
// predicado exato que ensureProviderAccountUnlinked replica.
func TestReauthKeepsOneRowAndDisconnectedReleasesTheAccount(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_harness_reauth")
	testpostgres.SeedProvider(t, pool, testpostgres.ProviderFixture{Code: "mercado_livre", DisplayName: "Mercado Livre"})

	repo := integrationspostgres.NewInstallationRepository(pool, cfg.DefaultTenantID)
	ctx := context.Background()
	now := time.Now().UTC()
	account := fmt.Sprintf("seller-%d", now.UnixNano())
	first := fmt.Sprintf("inst-reauth-a-%d", now.UnixNano())
	second := fmt.Sprintf("inst-reauth-b-%d", now.UnixNano())

	t.Cleanup(func() {
		for _, id := range []string{first, second} {
			if _, err := pool.Exec(ctx, `DELETE FROM integration_auth_sessions WHERE installation_id = $1`, id); err != nil {
				t.Errorf("cleanup integration_auth_sessions %s: %v", id, err)
			}
			if _, err := pool.Exec(ctx, `DELETE FROM integration_installations WHERE installation_id = $1`, id); err != nil {
				t.Errorf("cleanup integration_installations %s: %v", id, err)
			}
		}
	})

	create := func(id string, status integrationsdomain.InstallationStatus) error {
		return repo.CreateInstallation(ctx, integrationsdomain.Installation{
			InstallationID: id,
			TenantID:       cfg.DefaultTenantID,
			ProviderCode:   "mercado_livre",
			Family:         integrationsdomain.IntegrationFamilyMarketplace,
			DisplayName:    "Mercado Livre (teste)",
			Status:         status,
			HealthStatus:   integrationsdomain.HealthStatusHealthy,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
	}

	connect := func(id string, status integrationsdomain.InstallationStatus) error {
		inst := integrationsdomain.Installation{
			InstallationID:    id,
			TenantID:          cfg.DefaultTenantID,
			ProviderCode:      "mercado_livre",
			ExternalAccountID: account,
			Status:            status,
			HealthStatus:      integrationsdomain.HealthStatusHealthy,
		}
		snapshot := integrationsdomain.ProjectConnectionSnapshot(
			inst, integrationsdomain.AuthStrategyOAuth2, nil, "",
		)
		return repo.ApplyConnectionSnapshot(ctx, id, snapshot, "")
	}

	if err := create(first, integrationsdomain.InstallationStatusPendingConnection); err != nil {
		t.Fatalf("CreateInstallation first: %v", err)
	}
	if err := connect(first, integrationsdomain.InstallationStatusConnected); err != nil {
		t.Fatalf("primeira conexao: %v", err)
	}

	// --- 1. Reautorizar a MESMA installation continua sendo uma linha só -----
	if err := connect(first, integrationsdomain.InstallationStatusConnected); err != nil {
		t.Fatalf("reautorizacao da mesma installation: %v (o indice parcial nao deveria disparar)", err)
	}
	if got := countActiveForAccount(t, pool, account); got != 1 {
		t.Fatalf("linhas ativas com a conta = %d, want 1", got)
	}

	// --- 2. Controle positivo: uma SEGUNDA installation ativa é recusada -----
	// Sem este passo, o passo 3 seria vacuoso: se o índice não existisse, tudo
	// passaria e pareceria "liberado".
	if err := create(second, integrationsdomain.InstallationStatusPendingConnection); err != nil {
		t.Fatalf("CreateInstallation second: %v", err)
	}
	if err := connect(second, integrationsdomain.InstallationStatusConnected); err == nil {
		t.Fatal("uma segunda installation ativa com a mesma conta foi aceita: o indice unico nao esta valendo")
	}

	// --- 3. Desconectar a primeira LIBERA a conta ---------------------------
	// É exatamente o predicado que ensureProviderAccountUnlinked replica; se o
	// banco e o código discordarem aqui, um dos dois está errado.
	if err := connect(first, integrationsdomain.InstallationStatusDisconnected); err != nil {
		t.Fatalf("desconectar a primeira: %v", err)
	}
	if err := connect(second, integrationsdomain.InstallationStatusConnected); err != nil {
		t.Fatalf("segunda installation apos desconectar a primeira: %v (o indice exclui 'disconnected')", err)
	}
}

func countActiveForAccount(t *testing.T, pool *pgxpool.Pool, account string) int {
	t.Helper()
	var count int
	err := pool.QueryRow(context.Background(), `
		SELECT count(*) FROM integration_installations
		 WHERE external_account_id = $1
		   AND status NOT IN ('disconnected', 'failed')
	`, account).Scan(&count)
	if err != nil {
		t.Fatalf("contagem: %v", err)
	}
	return count
}
