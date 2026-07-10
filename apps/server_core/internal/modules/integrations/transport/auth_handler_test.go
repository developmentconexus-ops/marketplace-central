package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type stubAuthFlow struct {
	startInput    application.StartAuthorizeInput
	callbackInput application.HandleCallbackInput
	submitInput   application.SubmitAPIKeyInput
	statusInput   application.GetAuthStatusInput
	feeSyncInput  application.StartFeeSyncInput
	feeSyncAccept application.FeeSyncAccepted
	operations    []domain.OperationRun
	probeAccount  connectorsdomain.AccountSnapshot
	listings      []connectorsdomain.ListingSnapshot
	orders        []connectorsdomain.OrderSnapshot
	feeQuote      connectorsdomain.FeeQuoteSnapshot
	stock         connectorsdomain.StockSnapshot
	listingsLimit int
	ordersLimit   int
	feeQuoteInput connectorsdomain.FeeQuoteInput
	stockInput    connectorsdomain.ProviderListingRef
}

var _ AuthFlowReader = (*stubAuthFlow)(nil)

func (s *stubAuthFlow) StartAuthorize(ctx context.Context, input application.StartAuthorizeInput) (application.AuthorizeStart, error) {
	s.startInput = input
	return application.AuthorizeStart{InstallationID: input.InstallationID, ProviderCode: "mercado_livre", State: "state-1", AuthURL: "https://provider.test/auth"}, nil
}

func (s *stubAuthFlow) HandleCallback(ctx context.Context, input application.HandleCallbackInput) (application.AuthStatus, error) {
	s.callbackInput = input
	return application.AuthStatus{InstallationID: "inst-from-state", Status: domain.InstallationStatusConnected}, nil
}

func (s *stubAuthFlow) SubmitAPIKey(ctx context.Context, input application.SubmitAPIKeyInput) (application.AuthStatus, error) {
	s.submitInput = input
	return application.AuthStatus{InstallationID: input.InstallationID, Status: domain.InstallationStatusConnected}, nil
}

func (s *stubAuthFlow) RefreshCredential(ctx context.Context, input application.RefreshCredentialInput) (application.AuthStatus, error) {
	return application.AuthStatus{InstallationID: input.InstallationID, Status: domain.InstallationStatusConnected}, nil
}

func (s *stubAuthFlow) Disconnect(ctx context.Context, input application.DisconnectInput) (application.AuthStatus, error) {
	return application.AuthStatus{InstallationID: input.InstallationID, Status: domain.InstallationStatusDisconnected}, nil
}

func (s *stubAuthFlow) StartReauth(ctx context.Context, input application.StartReauthInput) (application.AuthorizeStart, error) {
	return application.AuthorizeStart{InstallationID: input.InstallationID, ProviderCode: "mercado_livre", State: "state-2", AuthURL: "https://provider.test/reauth"}, nil
}

func (s *stubAuthFlow) GetAuthStatus(ctx context.Context, input application.GetAuthStatusInput) (application.AuthStatus, error) {
	s.statusInput = input
	return application.AuthStatus{InstallationID: input.InstallationID, Status: domain.InstallationStatusConnected}, nil
}

func (s *stubAuthFlow) StartSync(ctx context.Context, input application.StartFeeSyncInput) (application.FeeSyncAccepted, error) {
	s.feeSyncInput = input
	if s.feeSyncAccept.OperationRunID == "" {
		s.feeSyncAccept = application.FeeSyncAccepted{
			InstallationID: input.InstallationID,
			OperationRunID: "op_001",
			Status:         domain.OperationRunStatusQueued,
		}
	}
	return s.feeSyncAccept, nil
}

func (s *stubAuthFlow) ListOperationRuns(ctx context.Context, installationID string) ([]domain.OperationRun, error) {
	if len(s.operations) == 0 {
		return []domain.OperationRun{
			{
				OperationRunID: "op_001",
				InstallationID: installationID,
				OperationType:  "pricing_fee_sync",
				Status:         domain.OperationRunStatusQueued,
				ResultCode:     "INTEGRATIONS_FEE_SYNC_QUEUED",
				AttemptCount:   1,
				ActorType:      "user",
			},
		}, nil
	}
	return s.operations, nil
}

func (s *stubAuthFlow) ProbeAccount(ctx context.Context, installationID string) (connectorsdomain.AccountSnapshot, error) {
	if s.probeAccount.ProviderAccountID == "" {
		s.probeAccount = connectorsdomain.AccountSnapshot{
			ProviderCode:        "mercado_livre",
			ProviderAccountID:   "691607102",
			ProviderAccountName: "METALNOBREACABAMENTOS",
			SiteID:              "MLB",
			Status:              "active",
		}
	}
	return s.probeAccount, nil
}

func (s *stubAuthFlow) ListListings(ctx context.Context, installationID string, limit int) ([]connectorsdomain.ListingSnapshot, error) {
	s.listingsLimit = limit
	if len(s.listings) == 0 {
		s.listings = []connectorsdomain.ListingSnapshot{{
			ProviderCode:   "mercado_livre",
			ProviderItemID: "MLB123",
			Title:          "Produto teste",
		}}
	}
	return s.listings, nil
}

func (s *stubAuthFlow) ListOrders(ctx context.Context, installationID string, limit int) ([]connectorsdomain.OrderSnapshot, error) {
	s.ordersLimit = limit
	if len(s.orders) == 0 {
		s.orders = []connectorsdomain.OrderSnapshot{{
			ProviderCode:    "mercado_livre",
			ProviderOrderID: "2001",
			ProviderStatus:  "paid",
		}}
	}
	return s.orders, nil
}

func (s *stubAuthFlow) ReadFeeQuote(ctx context.Context, installationID string, input connectorsdomain.FeeQuoteInput) (connectorsdomain.FeeQuoteSnapshot, error) {
	s.feeQuoteInput = input
	if s.feeQuote.ListingTypeID == "" {
		percent := 11.0
		fixed := 6.0
		s.feeQuote = connectorsdomain.FeeQuoteSnapshot{
			ProviderCode:      "mercado_livre",
			SiteID:            "MLB",
			ListingTypeID:     "gold_special",
			PriceAmount:       100,
			CurrencyID:        "BRL",
			CommissionPercent: &percent,
			FixedFeeAmount:    &fixed,
		}
	}
	return s.feeQuote, nil
}

func (s *stubAuthFlow) ReadStock(ctx context.Context, installationID string, input connectorsdomain.ProviderListingRef) (connectorsdomain.StockSnapshot, error) {
	s.stockInput = input
	if s.stock.ProviderItemID == "" {
		s.stock = connectorsdomain.StockSnapshot{
			ProviderCode:      "mercado_livre",
			ProviderItemID:    input.ProviderItemID,
			AvailableQuantity: 11,
			ProviderStatus:    "active",
			Scope:             connectorsdomain.StockScopeItem,
		}
	}
	return s.stock, nil
}

func TestAuthHandlerCompileTimeFeeSyncContract(t *testing.T) {
	t.Parallel()

	assertFeeSyncTransportContract(t, &stubAuthFlow{})
}

func assertFeeSyncTransportContract(t *testing.T, flow AuthFlowReader) {
	t.Helper()

	_, _ = flow.StartSync(context.Background(), application.StartFeeSyncInput{InstallationID: "inst_001"})
	_, _ = flow.ListOperationRuns(context.Background(), "inst_001")
}

func TestAuthHandlerStartAuthorizeDelegatesToService(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	handler := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodPost, "/integrations/installations/inst-1/auth/authorize", bytes.NewReader([]byte(`{"redirect_uri":"https://app.test/callback","scopes":["read"]}`)))
	rr := httptest.NewRecorder()

	handler.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rr.Code, rr.Body.String())
	}
	if flow.startInput.InstallationID != "inst-1" || flow.startInput.RedirectURI != "https://app.test/callback" {
		t.Fatalf("start input = %#v, want installation and redirect", flow.startInput)
	}

	var response application.AuthorizeStart
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.AuthURL == "" || response.State == "" {
		t.Fatalf("response = %#v, want auth URL and state", response)
	}
}

func TestAuthHandlerStartAuthorizeUsesEnvRedirectWhenRequestMissingRedirect(t *testing.T) {
	t.Setenv("MPC_OAUTH_REDIRECT_URI", "https://public.example/integrations/auth/callback")

	flow := &stubAuthFlow{}
	handler := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodPost, "/integrations/installations/inst-1/auth/authorize", bytes.NewReader([]byte(`{}`)))
	rr := httptest.NewRecorder()

	handler.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rr.Code, rr.Body.String())
	}
	if flow.startInput.RedirectURI != "https://public.example/integrations/auth/callback" {
		t.Fatalf("redirect_uri = %q, want env fallback redirect", flow.startInput.RedirectURI)
	}
}

func TestAuthHandlerRejectsWrongMethod(t *testing.T) {
	t.Parallel()

	handler := NewAuthHandler(&stubAuthFlow{})
	req := httptest.NewRequest(http.MethodGet, "/integrations/installations/inst-1/auth/authorize", nil)
	rr := httptest.NewRecorder()

	handler.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rr.Code)
	}
	if rr.Header().Get("Allow") != http.MethodPost {
		t.Fatalf("Allow = %q, want POST", rr.Header().Get("Allow"))
	}
}

func TestAuthHandlerCallbackAcceptsProviderCodeAndStateOnly(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	handler := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/auth/callback?code=provider-code&state=signed-state", nil)
	rr := httptest.NewRecorder()

	handler.handleCallback(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302; body=%s", rr.Code, rr.Body.String())
	}
	if flow.callbackInput.InstallationID != "" || flow.callbackInput.Code != "provider-code" || flow.callbackInput.State != "signed-state" {
		t.Fatalf("callback input = %#v, want code/state without installation id", flow.callbackInput)
	}
	if location := rr.Header().Get("Location"); location != "/integrations?auth=connected&installation=inst-from-state" {
		t.Fatalf("Location = %q, want service installation redirect", location)
	}
}

func TestAuthHandlerCallbackUsesEnvRedirectURI(t *testing.T) {
	t.Setenv("MPC_OAUTH_REDIRECT_URI", "https://public.example/integrations/auth/callback")

	flow := &stubAuthFlow{}
	handler := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/auth/callback?code=provider-code&state=signed-state", nil)
	rr := httptest.NewRecorder()

	handler.handleCallback(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302; body=%s", rr.Code, rr.Body.String())
	}
	if flow.callbackInput.RedirectURI != "https://public.example/integrations/auth/callback" {
		t.Fatalf("callback redirect_uri = %q, want env fallback redirect", flow.callbackInput.RedirectURI)
	}
	if location := rr.Header().Get("Location"); location != "https://public.example/integrations?auth=connected&installation=inst-from-state" {
		t.Fatalf("Location = %q, want callback origin redirect", location)
	}
}

func TestAuthHandlerCallbackAcceptsAmazonSPAPICode(t *testing.T) {
	t.Setenv("MPC_OAUTH_REDIRECT_URI", "https://public.example/integrations/auth/callback")

	flow := &stubAuthFlow{}
	handler := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/auth/callback?spapi_oauth_code=spapi-code&selling_partner_id=A1SELLER&state=signed-state", nil)
	rr := httptest.NewRecorder()

	handler.handleCallback(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302; body=%s", rr.Code, rr.Body.String())
	}
	if flow.callbackInput.Code != "spapi-code" || flow.callbackInput.State != "signed-state" {
		t.Fatalf("callback input = %#v, want SP-API code and state", flow.callbackInput)
	}
	if flow.callbackInput.ProviderAccountID != "A1SELLER" {
		t.Fatalf("provider account id = %q, want selling partner id", flow.callbackInput.ProviderAccountID)
	}
}

func TestAuthHandlerCallbackPrefersShopeeShopID(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	handler := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/auth/callback?code=provider-code&state=signed-state&shop_id=shop-1&merchant_id=merchant-1&selling_partner_id=seller-1", nil)
	rr := httptest.NewRecorder()

	handler.handleCallback(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302; body=%s", rr.Code, rr.Body.String())
	}
	if flow.callbackInput.ProviderAccountID != "shop-1" {
		t.Fatalf("provider account id = %q, want shop id", flow.callbackInput.ProviderAccountID)
	}
}

func TestAuthHandlerCallbackRedirectsToWebOriginWhenConfigured(t *testing.T) {
	t.Setenv("MPC_WEB_ORIGIN", "https://app.example")

	flow := &stubAuthFlow{}
	handler := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/auth/callback?code=provider-code&state=signed-state", nil)
	rr := httptest.NewRecorder()

	handler.handleCallback(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302; body=%s", rr.Code, rr.Body.String())
	}
	if location := rr.Header().Get("Location"); location != "https://app.example/integrations?auth=connected&installation=inst-from-state" {
		t.Fatalf("Location = %q, want frontend origin redirect", location)
	}
}

func TestAuthHandlerCallbackPrefersCallbackOriginOverWebOrigin(t *testing.T) {
	t.Setenv("MPC_WEB_ORIGIN", "http://localhost:5174")
	t.Setenv("MPC_OAUTH_REDIRECT_URI", "https://public.example/integrations/auth/callback")

	flow := &stubAuthFlow{}
	handler := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/auth/callback?code=provider-code&state=signed-state", nil)
	rr := httptest.NewRecorder()

	handler.handleCallback(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302; body=%s", rr.Code, rr.Body.String())
	}
	if location := rr.Header().Get("Location"); location != "https://public.example/integrations?auth=connected&installation=inst-from-state" {
		t.Fatalf("Location = %q, want callback origin to win over localhost web origin", location)
	}
}

func TestAuthHandlerSubmitAPIKeyDelegatesToService(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	handler := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodPost, "/integrations/installations/inst-shopee/auth/credentials", bytes.NewReader([]byte(`{"api_key":"secret","metadata":{"shop_id":"shop-1"}}`)))
	rr := httptest.NewRecorder()

	handler.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rr.Code, rr.Body.String())
	}
	if flow.submitInput.InstallationID != "inst-shopee" || flow.submitInput.APIKey != "secret" || flow.submitInput.Metadata["shop_id"] != "shop-1" {
		t.Fatalf("submit input = %#v, want api key payload", flow.submitInput)
	}
}

func TestHandleInstallationFeeSyncReturnsAccepted(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{
		feeSyncAccept: application.FeeSyncAccepted{
			InstallationID: "inst_001",
			OperationRunID: "op_001",
			Status:         domain.OperationRunStatusQueued,
		},
	}
	h := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodPost, "/integrations/installations/inst_001/fee-sync", bytes.NewReader([]byte(`{}`)))
	rr := httptest.NewRecorder()

	h.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if flow.feeSyncInput.InstallationID != "inst_001" || flow.feeSyncInput.ActorType != "user" {
		t.Fatalf("fee sync input = %#v", flow.feeSyncInput)
	}
}

func TestHandleInstallationOperationsReturnsItems(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	h := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/installations/inst_001/operations", nil)
	rr := httptest.NewRecorder()

	h.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var payload struct {
		Items []domain.OperationRun `json:"items"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode operations response: %v", err)
	}
	if len(payload.Items) != 1 || payload.Items[0].InstallationID != "inst_001" {
		t.Fatalf("items=%#v", payload.Items)
	}
}

func TestHandleInstallationAccountProbeReturnsSnapshot(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	h := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodPost, "/integrations/installations/inst_001/probes/account", bytes.NewReader([]byte(`{}`)))
	rr := httptest.NewRecorder()

	h.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var payload connectorsdomain.AccountSnapshot
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode account probe response: %v", err)
	}
	if payload.ProviderAccountName != "METALNOBREACABAMENTOS" {
		t.Fatalf("payload=%#v", payload)
	}
}

func TestHandleInstallationAccountProbeUsesSnakeCaseJSON(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	h := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodPost, "/integrations/installations/inst_001/probes/account", bytes.NewReader([]byte(`{}`)))
	rr := httptest.NewRecorder()

	h.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode account probe response: %v", err)
	}
	if _, ok := payload["provider_code"]; !ok {
		t.Fatalf("payload=%v, want snake_case provider_code", payload)
	}
	if _, ok := payload["ProviderCode"]; ok {
		t.Fatalf("payload=%v, want no PascalCase ProviderCode", payload)
	}
}

func TestHandleInstallationListingsProbeReturnsItems(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	h := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/installations/inst_001/probes/listings?limit=3", nil)
	rr := httptest.NewRecorder()

	h.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if flow.listingsLimit != 3 {
		t.Fatalf("limit=%d", flow.listingsLimit)
	}
}

func TestHandleInstallationListingsProbeUsesSnakeCaseJSON(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	h := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/installations/inst_001/probes/listings?limit=1", nil)
	rr := httptest.NewRecorder()

	h.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var payload struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode listings probe response: %v", err)
	}
	if len(payload.Items) != 1 {
		t.Fatalf("items=%v", payload.Items)
	}
	if _, ok := payload.Items[0]["provider_item_id"]; !ok {
		t.Fatalf("payload=%v, want snake_case provider_item_id", payload.Items[0])
	}
	if _, ok := payload.Items[0]["ProviderItemID"]; ok {
		t.Fatalf("payload=%v, want no PascalCase ProviderItemID", payload.Items[0])
	}
}

func TestHandleInstallationOrdersProbeReturnsItems(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	h := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/installations/inst_001/probes/orders?limit=1", nil)
	rr := httptest.NewRecorder()

	h.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if flow.ordersLimit != 1 {
		t.Fatalf("limit=%d", flow.ordersLimit)
	}
}

func TestHandleInstallationFeeQuoteProbeReturnsSnapshot(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	h := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/installations/inst_001/probes/fee-quote?listing_type_id=gold_special&price=100&currency_id=BRL&site_id=MLB", nil)
	rr := httptest.NewRecorder()

	h.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if flow.feeQuoteInput.ListingTypeID != "gold_special" || flow.feeQuoteInput.PriceAmount != 100 {
		t.Fatalf("fee quote input=%#v", flow.feeQuoteInput)
	}
}

func TestHandleInstallationFeeQuoteProbeUsesSnakeCaseJSON(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	h := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/installations/inst_001/probes/fee-quote?listing_type_id=gold_special&price=100&currency_id=BRL&site_id=MLB", nil)
	rr := httptest.NewRecorder()

	h.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode fee quote response: %v", err)
	}
	if _, ok := payload["listing_type_id"]; !ok {
		t.Fatalf("payload=%v, want snake_case listing_type_id", payload)
	}
	if _, ok := payload["ListingTypeID"]; ok {
		t.Fatalf("payload=%v, want no PascalCase ListingTypeID", payload)
	}
}

func TestHandleInstallationStockProbeReturnsSnapshot(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	h := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/installations/inst_001/probes/stock?provider_item_id=MLB123&provider_variation_id=VAR1", nil)
	rr := httptest.NewRecorder()

	h.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if flow.stockInput.ProviderItemID != "MLB123" || flow.stockInput.ProviderVariationID != "VAR1" {
		t.Fatalf("stock input=%#v", flow.stockInput)
	}
	var payload connectorsdomain.StockSnapshot
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode stock response: %v", err)
	}
	if payload.AvailableQuantity != 11 {
		t.Fatalf("payload=%#v", payload)
	}
}

func TestHandleInstallationStockProbeRequiresProviderItemID(t *testing.T) {
	t.Parallel()

	h := NewAuthHandler(&stubAuthFlow{})
	req := httptest.NewRequest(http.MethodGet, "/integrations/installations/inst_001/probes/stock", nil)
	rr := httptest.NewRecorder()

	h.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestHandleInstallationStockProbeUsesSnakeCaseJSON(t *testing.T) {
	t.Parallel()

	flow := &stubAuthFlow{}
	h := NewAuthHandler(flow)
	req := httptest.NewRequest(http.MethodGet, "/integrations/installations/inst_001/probes/stock?provider_item_id=MLB123", nil)
	rr := httptest.NewRecorder()

	h.handleInstallationAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode stock response: %v", err)
	}
	if _, ok := payload["provider_item_id"]; !ok {
		t.Fatalf("payload=%v, want snake_case provider_item_id", payload)
	}
	if _, ok := payload["available_quantity"]; !ok {
		t.Fatalf("payload=%v, want snake_case available_quantity", payload)
	}
	if _, ok := payload["ProviderItemID"]; ok {
		t.Fatalf("payload=%v, want no PascalCase ProviderItemID", payload)
	}
}
