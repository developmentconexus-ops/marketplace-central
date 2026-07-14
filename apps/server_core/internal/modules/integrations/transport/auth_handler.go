package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type AuthFlowReader interface {
	StartAuthorize(ctx context.Context, input application.StartAuthorizeInput) (application.AuthorizeStart, error)
	HandleCallback(ctx context.Context, input application.HandleCallbackInput) (application.AuthStatus, error)
	SubmitAPIKey(ctx context.Context, input application.SubmitAPIKeyInput) (application.AuthStatus, error)
	RefreshCredential(ctx context.Context, input application.RefreshCredentialInput) (application.AuthStatus, error)
	Disconnect(ctx context.Context, input application.DisconnectInput) (application.AuthStatus, error)
	StartReauth(ctx context.Context, input application.StartReauthInput) (application.AuthorizeStart, error)
	GetAuthStatus(ctx context.Context, input application.GetAuthStatusInput) (application.AuthStatus, error)
	StartSync(ctx context.Context, input application.StartFeeSyncInput) (application.FeeSyncAccepted, error)
	ListOperationRuns(ctx context.Context, installationID string) ([]domain.OperationRun, error)
	ProbeAccount(ctx context.Context, installationID string) (connectorsdomain.AccountSnapshot, error)
	ListListings(ctx context.Context, installationID string, limit int) ([]connectorsdomain.ListingSnapshot, error)
	ListOrders(ctx context.Context, installationID string, limit int) ([]connectorsdomain.OrderSnapshot, error)
	ReadFeeQuote(ctx context.Context, installationID string, input connectorsdomain.FeeQuoteInput) (connectorsdomain.FeeQuoteSnapshot, error)
	ReadStock(ctx context.Context, installationID string, ref connectorsdomain.ProviderListingRef) (connectorsdomain.StockSnapshot, error)
}

type AuthHandler struct {
	flow AuthFlowReader
}

func NewAuthHandler(flow AuthFlowReader) AuthHandler {
	return AuthHandler{flow: flow}
}

func (h AuthHandler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("/integrations/installations/", h.handleInstallationAuth)
	mux.HandleFunc("/integrations/auth/callback", h.handleCallback)
}

func (h AuthHandler) handleCallback(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeIntegrationError(w, http.StatusMethodNotAllowed, "INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED", "method not allowed")
		slog.Info("integrations.auth.callback", "action", "handle_callback", "result", "405", "duration_ms", time.Since(start).Milliseconds())
		return
	}
	code := firstNonEmptyQuery(r.URL.Query(), "code", "spapi_oauth_code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		writeIntegrationError(w, http.StatusBadRequest, "INTEGRATIONS_AUTH_STATE_INVALID", "missing callback params")
		slog.Info("integrations.auth.callback", "action", "handle_callback", "result", "400", "duration_ms", time.Since(start).Milliseconds())
		return
	}
	redirectURI := strings.TrimSpace(os.Getenv("MPC_OAUTH_REDIRECT_URI"))
	result, err := h.flow.HandleCallback(r.Context(), application.HandleCallbackInput{
		Code:              code,
		State:             state,
		RedirectURI:       redirectURI,
		ProviderAccountID: firstNonEmptyQuery(r.URL.Query(), "shop_id", "merchant_id", "selling_partner_id"),
	})
	if err != nil {
		slog.Warn("integrations.auth.callback", "action", "handle_callback", "result", "302", "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		http.Redirect(w, r, buildOAuthCallbackRedirectURL(buildIntegrationsRedirectPath("failed", "")), http.StatusFound)
		return
	}
	slog.Info("integrations.auth.callback", "action", "handle_callback", "result", "302", "duration_ms", time.Since(start).Milliseconds())
	http.Redirect(w, r, buildOAuthCallbackRedirectURL(buildIntegrationsRedirectPath("connected", result.InstallationID)), http.StatusFound)
}

func (h AuthHandler) handleInstallationAuth(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	path := strings.TrimPrefix(r.URL.Path, "/integrations/installations/")
	segments := strings.Split(path, "/")
	if len(segments) < 2 {
		slog.Info("integrations.auth.installation", "action", "route_not_found", "result", "404", "duration_ms", time.Since(start).Milliseconds())
		http.NotFound(w, r)
		return
	}
	installationID := segments[0]
	suffix := strings.Join(segments[1:], "/")

	switch suffix {
	case "auth/authorize", "auth/start":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeIntegrationError(w, http.StatusMethodNotAllowed, "INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED", "method not allowed")
			slog.Info("integrations.auth.start", "action", "start_authorize", "result", "405", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		var body struct {
			RedirectURI string   `json:"redirect_uri"`
			Scopes      []string `json:"scopes"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		redirectURI := strings.TrimSpace(body.RedirectURI)
		if redirectURI == "" {
			redirectURI = strings.TrimSpace(os.Getenv("MPC_OAUTH_REDIRECT_URI"))
		}
		result, err := h.flow.StartAuthorize(r.Context(), application.StartAuthorizeInput{
			InstallationID: installationID,
			RedirectURI:    redirectURI,
			Scopes:         body.Scopes,
		})
		if err != nil {
			status, code, message := mapIntegrationError(err)
			writeIntegrationError(w, status, code, message)
			slog.Error("integrations.auth.start", "action", "start_authorize", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			return
		}
		slog.Info("integrations.auth.start", "action", "start_authorize", "result", "200", "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, result)
	case "auth/credentials", "auth/api-key":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeIntegrationError(w, http.StatusMethodNotAllowed, "INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED", "method not allowed")
			slog.Info("integrations.auth.credentials", "action", "submit_credentials", "result", "405", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		var body struct {
			APIKey      string            `json:"api_key"`
			Credentials map[string]string `json:"credentials"`
			Metadata    map[string]string `json:"metadata"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeIntegrationError(w, http.StatusBadRequest, "INTEGRATIONS_APIKEY_MISSING_FIELDS", "malformed request body")
			slog.Info("integrations.auth.credentials", "action", "submit_credentials", "result", "400", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		apiKey := body.APIKey
		if apiKey == "" {
			apiKey = body.Credentials["api_key"]
		}
		meta := body.Metadata
		if len(meta) == 0 {
			meta = body.Credentials
		}
		result, err := h.flow.SubmitAPIKey(r.Context(), application.SubmitAPIKeyInput{
			InstallationID: installationID,
			APIKey:         apiKey,
			Metadata:       meta,
		})
		if err != nil {
			status, code, message := mapIntegrationError(err)
			writeIntegrationError(w, status, code, message)
			slog.Error("integrations.auth.credentials", "action", "submit_credentials", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			return
		}
		slog.Info("integrations.auth.credentials", "action", "submit_credentials", "result", "200", "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, result)
	case "disconnect":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeIntegrationError(w, http.StatusMethodNotAllowed, "INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED", "method not allowed")
			slog.Info("integrations.auth.disconnect", "action", "disconnect", "result", "405", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		result, err := h.flow.Disconnect(r.Context(), application.DisconnectInput{InstallationID: installationID})
		if err != nil {
			status, code, message := mapIntegrationError(err)
			writeIntegrationError(w, status, code, message)
			slog.Error("integrations.auth.disconnect", "action", "disconnect", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			return
		}
		slog.Info("integrations.auth.disconnect", "action", "disconnect", "result", "200", "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, result)
	case "reauth/authorize":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeIntegrationError(w, http.StatusMethodNotAllowed, "INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED", "method not allowed")
			slog.Info("integrations.auth.reauth", "action", "start_reauth", "result", "405", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		var body struct {
			RedirectURI string   `json:"redirect_uri"`
			Scopes      []string `json:"scopes"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		redirectURI := strings.TrimSpace(body.RedirectURI)
		if redirectURI == "" {
			redirectURI = strings.TrimSpace(os.Getenv("MPC_OAUTH_REDIRECT_URI"))
		}
		result, err := h.flow.StartReauth(r.Context(), application.StartReauthInput{
			InstallationID: installationID,
			RedirectURI:    redirectURI,
			Scopes:         body.Scopes,
		})
		if err != nil {
			status, code, message := mapIntegrationError(err)
			writeIntegrationError(w, status, code, message)
			slog.Error("integrations.auth.reauth", "action", "start_reauth", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			return
		}
		slog.Info("integrations.auth.reauth", "action", "start_reauth", "result", "200", "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, result)
	case "auth/status":
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			writeIntegrationError(w, http.StatusMethodNotAllowed, "INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED", "method not allowed")
			slog.Info("integrations.auth.status", "action", "get_status", "result", "405", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		result, err := h.flow.GetAuthStatus(r.Context(), application.GetAuthStatusInput{
			InstallationID: installationID,
		})
		if err != nil {
			status, code, message := mapIntegrationError(err)
			writeIntegrationError(w, status, code, message)
			slog.Error("integrations.auth.status", "action", "get_status", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			return
		}
		slog.Info("integrations.auth.status", "action", "get_status", "result", "200", "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, result)
	case "fee-sync":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeIntegrationError(w, http.StatusMethodNotAllowed, "INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED", "method not allowed")
			slog.Info("integrations.fee_sync.start", "action", "start_sync", "result", "405", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		result, err := h.flow.StartSync(r.Context(), application.StartFeeSyncInput{
			InstallationID: installationID,
			ActorType:      "user",
		})
		if err != nil {
			status, code, message := mapIntegrationError(err)
			writeIntegrationError(w, status, code, message)
			slog.Error("integrations.fee_sync.start", "action", "start_sync", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			return
		}
		slog.Info("integrations.fee_sync.start", "action", "start_sync", "result", "202", "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusAccepted, result)
	case "operations":
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			writeIntegrationError(w, http.StatusMethodNotAllowed, "INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED", "method not allowed")
			slog.Info("integrations.operations.list", "action", "list_operation_runs", "result", "405", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		items, err := h.flow.ListOperationRuns(r.Context(), installationID)
		if err != nil {
			status, code, message := mapIntegrationError(err)
			writeIntegrationError(w, status, code, message)
			slog.Error("integrations.operations.list", "action", "list_operation_runs", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			return
		}
		slog.Info("integrations.operations.list", "action", "list_operation_runs", "result", "200", "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
	case "probes/account":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeIntegrationError(w, http.StatusMethodNotAllowed, "INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED", "method not allowed")
			slog.Info("integrations.probes.account", "action", "probe_account", "result", "405", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		result, err := h.flow.ProbeAccount(r.Context(), installationID)
		if err != nil {
			status, code, message := mapIntegrationError(err)
			writeIntegrationError(w, status, code, message)
			slog.Error("integrations.probes.account", "action", "probe_account", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			return
		}
		slog.Info("integrations.probes.account", "action", "probe_account", "result", "200", "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, result)
	case "probes/listings":
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			writeIntegrationError(w, http.StatusMethodNotAllowed, "INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED", "method not allowed")
			slog.Info("integrations.probes.listings", "action", "list_listings", "result", "405", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		result, err := h.flow.ListListings(r.Context(), installationID, parsePositiveInt(r.URL.Query().Get("limit"), 20))
		if err != nil {
			status, code, message := mapIntegrationError(err)
			writeIntegrationError(w, status, code, message)
			slog.Error("integrations.probes.listings", "action", "list_listings", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			return
		}
		slog.Info("integrations.probes.listings", "action", "list_listings", "result", "200", "count", len(result), "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": result})
	case "probes/orders":
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			writeIntegrationError(w, http.StatusMethodNotAllowed, "INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED", "method not allowed")
			slog.Info("integrations.probes.orders", "action", "list_orders", "result", "405", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		result, err := h.flow.ListOrders(r.Context(), installationID, parsePositiveInt(r.URL.Query().Get("limit"), 20))
		if err != nil {
			status, code, message := mapIntegrationError(err)
			writeIntegrationError(w, status, code, message)
			slog.Error("integrations.probes.orders", "action", "list_orders", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			return
		}
		slog.Info("integrations.probes.orders", "action", "list_orders", "result", "200", "count", len(result), "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": result})
	case "probes/fee-quote":
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			writeIntegrationError(w, http.StatusMethodNotAllowed, "INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED", "method not allowed")
			slog.Info("integrations.probes.fee_quote", "action", "read_fee_quote", "result", "405", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		priceAmount, err := parseRequiredFloat(r.URL.Query().Get("price"))
		if err != nil {
			writeIntegrationError(w, http.StatusBadRequest, "INTEGRATIONS_OPERATION_INVALID", "price query parameter must be a positive number")
			slog.Info("integrations.probes.fee_quote", "action", "read_fee_quote", "result", "400", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		result, err := h.flow.ReadFeeQuote(r.Context(), installationID, connectorsdomain.FeeQuoteInput{
			SiteID:        strings.TrimSpace(r.URL.Query().Get("site_id")),
			CategoryID:    strings.TrimSpace(r.URL.Query().Get("category_id")),
			ListingTypeID: strings.TrimSpace(r.URL.Query().Get("listing_type_id")),
			PriceAmount:   priceAmount,
			CurrencyID:    strings.TrimSpace(r.URL.Query().Get("currency_id")),
		})
		if err != nil {
			status, code, message := mapIntegrationError(err)
			writeIntegrationError(w, status, code, message)
			slog.Error("integrations.probes.fee_quote", "action", "read_fee_quote", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			return
		}
		slog.Info("integrations.probes.fee_quote", "action", "read_fee_quote", "result", "200", "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, result)
	case "probes/stock":
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			writeIntegrationError(w, http.StatusMethodNotAllowed, "INTEGRATIONS_AUTH_METHOD_NOT_ALLOWED", "method not allowed")
			slog.Info("integrations.probes.stock", "action", "read_stock", "result", "405", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		providerItemID := strings.TrimSpace(r.URL.Query().Get("provider_item_id"))
		if providerItemID == "" {
			writeIntegrationError(w, http.StatusBadRequest, "INTEGRATIONS_OPERATION_INVALID", "provider_item_id query parameter is required")
			slog.Info("integrations.probes.stock", "action", "read_stock", "result", "400", "duration_ms", time.Since(start).Milliseconds())
			return
		}
		result, err := h.flow.ReadStock(r.Context(), installationID, connectorsdomain.ProviderListingRef{
			ProviderItemID:      providerItemID,
			ProviderVariationID: strings.TrimSpace(r.URL.Query().Get("provider_variation_id")),
		})
		if err != nil {
			status, code, message := mapIntegrationError(err)
			writeIntegrationError(w, status, code, message)
			slog.Error("integrations.probes.stock", "action", "read_stock", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			return
		}
		slog.Info("integrations.probes.stock", "action", "read_stock", "result", "200", "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, result)
	default:
		slog.Info("integrations.auth.installation", "action", "route_not_found", "result", "404", "duration_ms", time.Since(start).Milliseconds())
		http.NotFound(w, r)
	}
}

func buildWebRedirectURL(path string) string {
	webOrigin := strings.TrimSpace(os.Getenv("MPC_WEB_ORIGIN"))
	if webOrigin == "" {
		return path
	}
	return strings.TrimRight(webOrigin, "/") + path
}

func buildOAuthCallbackRedirectURL(path string) string {
	redirectURI := strings.TrimSpace(os.Getenv("MPC_OAUTH_REDIRECT_URI"))
	if redirectURI != "" {
		if parsed, err := url.Parse(redirectURI); err == nil && parsed.Scheme != "" && parsed.Host != "" {
			return parsed.Scheme + "://" + parsed.Host + path
		}
	}
	return buildWebRedirectURL(path)
}

func firstNonEmptyQuery(values url.Values, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(values.Get(key)); value != "" {
			return value
		}
	}
	return ""
}

func buildIntegrationsRedirectPath(authStatus, installationID string) string {
	query := url.Values{}
	query.Set("auth", strings.TrimSpace(authStatus))
	if strings.TrimSpace(installationID) != "" {
		query.Set("installation", strings.TrimSpace(installationID))
	}
	return "/integrations?" + query.Encode()
}

func parsePositiveInt(raw string, fallback int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func parseRequiredFloat(raw string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil || value <= 0 {
		return 0, errInvalidFloat
	}
	return value, nil
}

var errInvalidFloat = strconv.ErrSyntax
