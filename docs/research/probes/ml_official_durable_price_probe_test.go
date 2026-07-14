package mercadolivre

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
	integrationscrypto "marketplace-central/apps/server_core/internal/modules/integrations/adapters/crypto"
	melauth "marketplace-central/apps/server_core/internal/modules/integrations/adapters/mercadolivre"
	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/platform/pgdb"
)

type durableBatch struct {
	Results []struct {
		SKU       string `json:"SKU"`
		CatalogID string `json:"CatalogID"`
	} `json:"results"`
}

func TestMLOfficialDurablePriceProbe(t *testing.T) {
	inPath, outPath := os.Getenv("ML_DURABLE_INPUT"), os.Getenv("ML_DURABLE_OUTPUT")
	if inPath == "" || outPath == "" {
		t.Fatal("ML_DURABLE_INPUT and ML_DURABLE_OUTPUT are required")
	}
	raw, err := os.ReadFile(inPath)
	if err != nil {
		t.Fatal(err)
	}
	var batch durableBatch
	if err := json.Unmarshal(raw, &batch); err != nil {
		t.Fatal(err)
	}

	cfg, err := pgdb.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Fatal(err)
	}
	installRepo := integrationspostgres.NewInstallationRepository(pool, cfg.DefaultTenantID)
	installationSvc := integrationsapp.NewInstallationService(installRepo, cfg.DefaultTenantID)
	installs, err := installRepo.ListInstallations(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var installationID, accountID string
	for _, x := range installs {
		if x.ProviderCode == "mercado_livre" && x.ExternalAccountID != "" && x.ActiveCredentialID != "" {
			installationID, accountID = x.InstallationID, x.ExternalAccountID
			break
		}
	}
	if installationID == "" {
		t.Fatal("no Mercado Livre installation with active credential")
	}
	credRepo := integrationspostgres.NewCredentialRepository(pool, cfg.DefaultTenantID)
	credSvc := integrationsapp.NewCredentialService(credRepo, cfg.DefaultTenantID)
	cryptoSvc, err := integrationscrypto.NewLocalKeyService(cfg.EncryptionKey, "local-key-v1")
	if err != nil {
		t.Fatal(err)
	}
	authSessionRepo := integrationspostgres.NewAuthSessionRepository(pool, cfg.DefaultTenantID)
	authSvc := integrationsapp.NewAuthService(authSessionRepo, cfg.DefaultTenantID)
	oauthStateRepo := integrationspostgres.NewOAuthStateRepository(pool, cfg.DefaultTenantID)
	oauthCodec, err := integrationsapp.NewHMACOAuthStateCodec(cfg.EncryptionKey)
	if err != nil {
		t.Fatal(err)
	}
	authAdapter := melauth.NewAdapter(melauth.Config{
		ClientID:     strings.TrimSpace(os.Getenv("MPC_PROVIDER_MERCADOLIVRE_CLIENT_ID")),
		ClientSecret: strings.TrimSpace(os.Getenv("MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET")),
		TokenURL:     "https://api.mercadolibre.com/oauth/token",
		APIBaseURL:   "https://api.mercadolibre.com",
		HTTPClient:   &http.Client{Timeout: 25 * time.Second},
	})
	authFlow, err := integrationsapp.NewAuthFlowService(integrationsapp.AuthFlowConfig{
		TenantID:        cfg.DefaultTenantID,
		Installations:   installationSvc,
		Credentials:     credSvc,
		AuthSessions:    authSvc,
		OAuthStates:     oauthStateRepo,
		OAuthStateCodec: oauthCodec,
		Encryptor:       cryptoSvc,
		Adapters:        []integrationsapp.MarketplaceAuthAdapter{authAdapter},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := authFlow.RefreshCredential(ctx, integrationsapp.RefreshCredentialInput{InstallationID: installationID}); err != nil {
		t.Fatalf("refresh Mercado Livre credential: %v", err)
	}
	resolver := integrationsapp.NewCredentialResolver(credSvc, cryptoSvc)
	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		HTTPClient: &http.Client{Timeout: 25 * time.Second},
		AccessTokenResolver: func(ctx context.Context, ref domain.ProviderAccountRef) (string, error) {
			x, err := resolver.ResolveAccessToken(ctx, integrationsapp.CredentialResolutionRef{
				TenantID: ref.TenantID, InstallationID: ref.InstallationID,
				ProviderCode: ref.ProviderCode, ProviderAccountID: ref.ProviderAccountID,
			})
			return x.AccessToken, err
		},
	})
	ref := domain.ProviderAccountRef{TenantID: cfg.DefaultTenantID, InstallationID: installationID, ProviderCode: "mercado_livre", ProviderAccountID: accountID}
	token, err := adapter.accessToken(ctx, ref)
	if err != nil {
		t.Fatal(err)
	}

	seen := map[string]bool{}
	products := make([]map[string]any, 0, 22)
	for _, row := range batch.Results {
		if row.CatalogID == "" || seen[row.CatalogID] {
			continue
		}
		seen[row.CatalogID] = true
		resp, body, err := adapter.doRaw(ctx, ref, token, http.MethodGet, "/products/"+url.PathEscape(row.CatalogID), nil)
		entry := map[string]any{"sku": row.SKU, "catalog_id": row.CatalogID}
		if err != nil {
			entry["transport_error"] = durableSafe(err.Error())
			products = append(products, entry)
			continue
		}
		entry["http_status"] = resp.StatusCode
		var product map[string]any
		if json.Unmarshal(body, &product) == nil {
			entry["status"] = product["status"]
			entry["domain_id"] = product["domain_id"]
			entry["buy_box_winner_price_range"] = product["buy_box_winner_price_range"]
			if winner, ok := product["buy_box_winner"].(map[string]any); ok {
				entry["buy_box_winner"] = map[string]any{
					"item_id": winner["item_id"], "price": winner["price"],
					"currency_id": winner["currency_id"], "condition": winner["condition"],
					"listing_type_id": winner["listing_type_id"],
				}
			}
		}
		products = append(products, entry)
	}
	sort.Slice(products, func(i, j int) bool {
		return fmt.Sprint(products[i]["catalog_id"]) < fmt.Sprint(products[j]["catalog_id"])
	})

	itemID := "MLB3474723449"
	itemEvidence := map[string]any{"item_id": itemID}
	resp, body, err := adapter.doRaw(ctx, ref, token, http.MethodGet, "/items/"+itemID, nil)
	if err != nil {
		itemEvidence["item_transport_error"] = durableSafe(err.Error())
	} else {
		itemEvidence["item_http_status"] = resp.StatusCode
		var item map[string]any
		if json.Unmarshal(body, &item) == nil {
			if resp.StatusCode == http.StatusOK {
				itemEvidence["belongs_to_connected_account"] = fmt.Sprint(item["seller_id"]) == accountID
				itemEvidence["status"] = item["status"]
				itemEvidence["sub_status"] = item["sub_status"]
				itemEvidence["catalog_product_id"] = item["catalog_product_id"]
				itemEvidence["price"] = item["price"]
				itemEvidence["available_quantity"] = item["available_quantity"]
			} else {
				itemEvidence["api_error"] = map[string]any{"code": item["error"], "message": item["message"]}
			}
		}
	}

	suggestionEvidence := map[string]any{"item_id": itemID}
	resp, body, err = adapter.doRaw(ctx, ref, token, http.MethodGet, "/suggestions/items/"+itemID+"/details", nil)
	if err != nil {
		suggestionEvidence["transport_error"] = durableSafe(err.Error())
	} else {
		suggestionEvidence["http_status"] = resp.StatusCode
		var s map[string]any
		if json.Unmarshal(body, &s) == nil {
			if resp.StatusCode != http.StatusOK {
				suggestionEvidence["api_error"] = map[string]any{"code": s["error"], "message": s["message"]}
			}
			for _, key := range []string{"item_id", "suggested_price", "lowest_price", "internal_price", "currency_id", "reference_type", "status", "last_updated"} {
				if value, ok := s[key]; ok {
					suggestionEvidence[key] = value
				}
			}
			if metadata, ok := s["metadata"].(map[string]any); ok {
				sanitized := map[string]any{}
				for _, key := range []string{"graph", "compared_values", "reference", "price_reference"} {
					if value, ok := metadata[key]; ok {
						sanitized[key] = value
					}
				}
				suggestionEvidence["metadata"] = sanitized
			}
		}
	}

	withWinner, withRange := 0, 0
	for _, p := range products {
		if p["buy_box_winner"] != nil {
			withWinner++
		}
		if p["buy_box_winner_price_range"] != nil {
			withRange++
		}
	}
	payload := map[string]any{
		"captured_at":                time.Now().UTC().Format(time.RFC3339),
		"method":                     "official_api_get_only",
		"products":                   products,
		"summary":                    map[string]any{"catalog_products": len(products), "with_buy_box_winner": withWinner, "with_price_range": withRange},
		"shown_item":                 itemEvidence,
		"shown_item_price_reference": suggestionEvidence,
	}
	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(outPath, out, 0600); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("ML_DURABLE_SUMMARY products=%d winner=%d range=%d item=%v suggestion=%v\n", len(products), withWinner, withRange, itemEvidence["item_http_status"], suggestionEvidence["http_status"])
}

func durableSafe(s string) string {
	if len(s) > 180 {
		return s[:180]
	}
	return s
}
