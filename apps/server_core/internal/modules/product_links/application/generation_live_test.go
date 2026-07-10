package application

import (
	"context"
	"os"
	"testing"
	"time"

	internalreadoracle "marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle"
	internalreadapp "marketplace-central/apps/server_core/internal/modules/internal_read/application"
	productlinkspostgres "marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestGenerateLinkCandidatesLive(t *testing.T) {
	if os.Getenv("MPC_PRODUCT_LINKS_LIVE_TEST") != "1" {
		t.Skip("set MPC_PRODUCT_LINKS_LIVE_TEST=1 to run live product links validation")
	}

	pgURL := os.Getenv("MPC_PRODUCT_LINKS_POSTGRES_URL")
	if pgURL == "" {
		pgURL = "postgres://marketplace:marketplace@127.0.0.1:5435/marketplace_central?sslmode=disable"
	}
	installationID := os.Getenv("MPC_PRODUCT_LINKS_INSTALLATION_ID")
	if installationID == "" {
		installationID = "inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98"
	}

	oracleCfg, err := internalreadoracle.LoadConfigFromEnv(os.Getenv)
	if err != nil {
		t.Fatalf("load Oracle config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, pgURL)
	if err != nil {
		t.Fatalf("open Postgres pool: %v", err)
	}
	defer pool.Close()

	oracleDB, err := internalreadoracle.OpenDB(ctx, oracleCfg)
	if err != nil {
		t.Fatalf("open Oracle db: %v", err)
	}
	defer oracleDB.Close()

	snapshotRepo := productlinkspostgres.NewListingSnapshotRepository(pool, "tenant_default")
	candidateRepo := productlinkspostgres.NewLinkCandidateRepository(pool, "tenant_default")
	matcher := internalreadapp.NewService(internalreadoracle.NewReader(oracleDB))
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: snapshotRepo,
		Matcher:   matcher,
		Store:     candidateRepo,
	})

	result, err := svc.GenerateLinkCandidates(ctx, GenerateLinkCandidatesInput{
		InstallationID: installationID,
		Limit:          20,
	})
	if err != nil {
		t.Fatalf("generate live link candidates: %v", err)
	}
	if result.GeneratedCount == 0 {
		t.Fatal("expected at least one generated candidate")
	}

	persisted, err := svc.ListLinkCandidates(ctx, ListLinkCandidatesInput{
		InstallationID: installationID,
		Limit:          20,
	})
	if err != nil {
		t.Fatalf("list persisted live candidates: %v", err)
	}
	if len(persisted) == 0 {
		t.Fatal("expected persisted candidates after live generation")
	}

	stateCounts := map[string]int{}
	for _, candidate := range persisted {
		stateCounts[string(candidate.State)]++
		if candidate.InstallationID == "" || candidate.ProviderItemID == "" {
			t.Fatalf("candidate missing required identity: %#v", candidate)
		}
	}
	t.Logf("generated_count=%d persisted_count=%d states=%v", result.GeneratedCount, len(persisted), stateCounts)
}
