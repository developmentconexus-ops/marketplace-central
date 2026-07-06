package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"marketplace-central/apps/server_core/internal/composition"
	"marketplace-central/apps/server_core/internal/platform/pgdb"
)

func TestNewRootRouterBuildsWithoutLegacyConnectorCredentials(t *testing.T) {
	t.Setenv("ME_CLIENT_ID", "")
	t.Setenv("ME_CLIENT_SECRET", "")

	_, err := composition.NewRootRouter(nil, nil, pgdb.Config{
		DefaultTenantID: "tenant_default",
		EncryptionKey:   "0123456789abcdef0123456789abcdef",
	})
	if err != nil {
		t.Fatalf("expected NewRootRouter to build without legacy connector credentials, got %v", err)
	}
}

func TestNewRootRouterBuildsWhenMelhorEnvioCredentialsArePresent(t *testing.T) {
	t.Setenv("ME_CLIENT_ID", "test-client")
	t.Setenv("ME_CLIENT_SECRET", "test-secret")

	router, err := composition.NewRootRouter(nil, nil, pgdb.Config{
		DefaultTenantID: "tenant_default",
		EncryptionKey:   "0123456789abcdef0123456789abcdef",
	})
	if err != nil {
		t.Fatalf("expected NewRootRouter without error, got %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected /healthz to return 200, got %d", rec.Code)
	}

	startReq := httptest.NewRequest(http.MethodGet, "/connectors/melhor-envio/auth/start", nil)
	startRec := httptest.NewRecorder()
	router.ServeHTTP(startRec, startReq)
	if startRec.Code == http.StatusNotFound {
		t.Fatal("expected Melhor Envio auth route to be registered")
	}
}
