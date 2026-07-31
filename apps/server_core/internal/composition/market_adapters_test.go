package composition

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	erpinternalread "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	internalreadapp "marketplace-central/apps/server_core/internal/modules/internal_read/application"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	marketports "marketplace-central/apps/server_core/internal/modules/market/ports"
)

// stubInternalReader is a ports.Reader whose FindProductsForLinking returns a
// caller-supplied (candidates, err); the other methods are unused by the
// identity adapter and panic if reached.
type stubInternalReader struct {
	candidates []internalreaddomain.ProductCandidate
	err        error
}

func (s stubInternalReader) FindProductsForLinking(context.Context, internalreadports.FindProductsInput) ([]internalreaddomain.ProductCandidate, error) {
	return s.candidates, s.err
}

func (stubInternalReader) GetSellableStock(context.Context, internalreadports.SellableStockInput) (internalreaddomain.SellableStock, error) {
	panic("unused")
}

func (stubInternalReader) GetCurrentPrice(context.Context, internalreadports.CurrentPriceInput) (internalreaddomain.CurrentPrice, error) {
	panic("unused")
}

func (stubInternalReader) GetCostAsOf(context.Context, internalreadports.CostAsOfInput) (internalreaddomain.CostAsOf, error) {
	panic("unused")
}

func (stubInternalReader) GetSalesHistory(context.Context, internalreadports.SalesHistoryInput) (internalreaddomain.SalesHistory, error) {
	panic("unused")
}

func (stubInternalReader) GetTaxInputs(context.Context, internalreadports.TaxInput) (internalreaddomain.TaxInputs, error) {
	panic("unused")
}

func (stubInternalReader) ListCatalogProductFacts(context.Context, internalreadports.Cursor, int, *internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogFactPage, error) {
	panic("unused")
}

func (stubInternalReader) SearchCatalogProductFacts(context.Context, string, internalreadports.Cursor, int, *internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogFactPage, error) {
	panic("unused")
}

func (stubInternalReader) CatalogProductFactsByIDs(context.Context, []int64) (internalreadports.CatalogFactPage, error) {
	panic("unused")
}

func (stubInternalReader) GetCatalogAssortmentCounts(context.Context, *internalreadports.SellableAssortmentPolicy) (internalreadports.CatalogAssortmentCounts, error) {
	panic("unused")
}

// TestGetLocalIdentityMapsReaderNotFoundToSentinel guards FINDING-M02-COLLECT-4XX
// item 4 (D-86): a codprod absent from the ERP source makes the internal_read
// reader return *erp_import.ERPProductNotFoundError. GetLocalIdentity MUST
// translate that to the market port sentinel marketports.ErrProductNotFound so
// the collection handler answers 404, not a raw 500 internal_error.
func TestGetLocalIdentityMapsReaderNotFoundToSentinel(t *testing.T) {
	svc := internalreadapp.NewService(stubInternalReader{
		err: &erpinternalread.ERPProductNotFoundError{ProductID: 90001},
	})
	adapter := newMarketProductIdentityReaderAdapter(svc, true, nil, nil)

	_, err := adapter.GetLocalIdentity(context.Background(), "90001")
	if !errors.Is(err, marketports.ErrProductNotFound) {
		t.Fatalf("GetLocalIdentity() error = %v, want ErrProductNotFound", err)
	}
}

func TestMapMarketConnectorError(t *testing.T) {
	statusOf := func(err error) (int, bool) {
		var statusErr *marketports.ProviderStatusError
		if errors.As(err, &statusErr) {
			return statusErr.StatusCode, true
		}
		return 0, false
	}

	t.Run("nil passes through", func(t *testing.T) {
		if got := mapMarketConnectorError(nil); got != nil {
			t.Fatalf("mapMarketConnectorError(nil) = %v, want nil", got)
		}
	})

	t.Run("deadline exceeded passes through as timeout", func(t *testing.T) {
		// Raw context error.
		got := mapMarketConnectorError(context.DeadlineExceeded)
		if !errors.Is(got, context.DeadlineExceeded) {
			t.Fatalf("timeout not preserved: %v", got)
		}
		// The connectors chain shape after the capability_adapter.go fix:
		// CapabilityError{Transient, Cause: transport-error-wrapping-deadline}.
		chained := &connectorsdomain.CapabilityError{
			Code:    connectorsdomain.ErrCodeProviderTransient,
			Message: "provider request failed",
			Cause:   fmt.Errorf("Get \"https://api\": %w", context.DeadlineExceeded),
		}
		got = mapMarketConnectorError(chained)
		if !errors.Is(got, context.DeadlineExceeded) {
			t.Fatalf("timeout not preserved through CapabilityError chain: %v", got)
		}
		if _, isStatus := statusOf(got); isStatus {
			t.Fatalf("timeout must not be flattened to ProviderStatusError: %v", got)
		}
	})

	t.Run("catalog offers unavailable maps to disabled sentinel", func(t *testing.T) {
		got := mapMarketConnectorError(connectorsdomain.ErrCatalogOffersUnavailable)
		if !errors.Is(got, marketports.ErrCatalogOffersDisabled) {
			t.Fatalf("got %v, want ErrCatalogOffersDisabled", got)
		}
	})

	t.Run("rate limited sentinel and capability code map to ErrRateLimited", func(t *testing.T) {
		for _, err := range []error{
			fmt.Errorf("wrapped: %w", connectorsdomain.ErrRateLimited),
			&connectorsdomain.CapabilityError{Code: connectorsdomain.ErrCodeProviderRateLimited},
		} {
			got := mapMarketConnectorError(err)
			if !errors.Is(got, marketports.ErrRateLimited) {
				t.Fatalf("mapMarketConnectorError(%v) = %v, want ErrRateLimited", err, got)
			}
		}
	})

	t.Run("status classes", func(t *testing.T) {
		cases := []struct {
			name string
			err  error
			want int
		}{
			{"sentinel unauthorized", fmt.Errorf("wrapped: %w", connectorsdomain.ErrUnauthorized), http.StatusUnauthorized},
			{"sentinel not found", fmt.Errorf("wrapped: %w", connectorsdomain.ErrNotFound), http.StatusNotFound},
			{"sentinel provider unavailable", fmt.Errorf("wrapped: %w", connectorsdomain.ErrProviderUnavailable), http.StatusBadGateway},
			{"capability auth", &connectorsdomain.CapabilityError{Code: connectorsdomain.ErrCodeProviderAuth}, http.StatusUnauthorized},
			{"capability invalid reference", &connectorsdomain.CapabilityError{Code: connectorsdomain.ErrCodeProviderInvalidReference}, http.StatusNotFound},
			{"capability transient", &connectorsdomain.CapabilityError{Code: connectorsdomain.ErrCodeProviderTransient}, http.StatusBadGateway},
			{"capability validation", &connectorsdomain.CapabilityError{Code: connectorsdomain.ErrCodeProviderValidation}, http.StatusBadRequest},
			{"capability payload invalid", &connectorsdomain.CapabilityError{Code: connectorsdomain.ErrCodeProviderPayloadInvalid}, http.StatusBadRequest},
			{"capability unsupported shape", &connectorsdomain.CapabilityError{Code: connectorsdomain.ErrCodeProviderUnsupportedShape}, http.StatusBadRequest},
			{"unknown error", errors.New("mystery"), http.StatusBadGateway},
		}
		for _, tc := range cases {
			got := mapMarketConnectorError(tc.err)
			status, ok := statusOf(got)
			if !ok {
				t.Fatalf("%s: got %v, want ProviderStatusError", tc.name, got)
			}
			if status != tc.want {
				t.Fatalf("%s: status = %d, want %d", tc.name, status, tc.want)
			}
		}
	})
}

func TestMapMarketProviderErrorContextDeadline(t *testing.T) {
	expired, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	// Even an already-flattened connector error becomes a TIMEOUT when the
	// caller's context deadline expired (IC-03 Amendment A1).
	got := mapMarketProviderError(expired, connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderTransient, "provider request failed"))
	if !errors.Is(got, context.DeadlineExceeded) {
		t.Fatalf("expired ctx: got %v, want wrap of context.DeadlineExceeded", got)
	}

	// A live context defers to the regular connector mapping.
	got = mapMarketProviderError(context.Background(), fmt.Errorf("wrapped: %w", connectorsdomain.ErrRateLimited))
	if !errors.Is(got, marketports.ErrRateLimited) {
		t.Fatalf("live ctx: got %v, want ErrRateLimited", got)
	}
}

// TestObserveProviderFailureWarnsOnceSanitized guards FINDING-M02-LIVE-2
// observability (D-86): a provider validation body (e.g. the ML 400 that today
// dies silently) must surface a single warn per operation+status, with the body
// whitespace-collapsed, token-redacted, and length-bounded — never the raw
// bearer token or an unbounded payload.
func TestObserveProviderFailureWarnsOnceSanitized(t *testing.T) {
	var buf bytes.Buffer
	adapter := &marketPriceIntelCollectorAdapter{
		log: slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})),
	}
	err := &connectorsdomain.CapabilityError{
		Code:    connectorsdomain.ErrCodeProviderValidation,
		Message: "  {\"message\":\"invalid item\",\n\t\"auth\":\"Bearer secret-abc123\"}  " + strings.Repeat("x", 400),
	}
	mapped := &marketports.ProviderStatusError{StatusCode: http.StatusBadRequest}

	adapter.observeProviderFailure("ListCatalogOffers", err, mapped)
	adapter.observeProviderFailure("ListCatalogOffers", err, mapped) // same op+status: deduped

	logged := buf.String()
	if n := strings.Count(logged, "market collection provider failure"); n != 1 {
		t.Fatalf("warn count = %d, want 1 (deduped); log =\n%s", n, logged)
	}
	if strings.Contains(logged, "secret-abc123") {
		t.Fatalf("provider token leaked into log:\n%s", logged)
	}
	if !strings.Contains(logged, "ListCatalogOffers") || !strings.Contains(logged, "400") {
		t.Fatalf("log missing operation/status:\n%s", logged)
	}
	if strings.Contains(logged, "\n\t") {
		t.Fatalf("provider body not whitespace-collapsed:\n%s", logged)
	}
	// A different status is a distinct signal and must warn again.
	adapter.observeProviderFailure("ListCatalogOffers", err, &marketports.ProviderStatusError{StatusCode: http.StatusBadGateway})
	if n := strings.Count(buf.String(), "market collection provider failure"); n != 2 {
		t.Fatalf("warn count after new status = %d, want 2", n)
	}
}

func TestSanitizeProviderBodyRedactsNonBearerCredentials(t *testing.T) {
	// The provider body is echoed straight from the HTTP response, so a leaked
	// credential can take any shape, not just "Bearer X" (FINDING-M02-LIVE-2
	// observability hardening — adversarial P6).
	secrets := []string{"APP_USR-1234567890123456-071912-abcdef", "TG-abc123def456", "refresh-XYZ789"}
	cases := []string{
		`{"error":"invalid_token","message":"Invalid access_token: APP_USR-1234567890123456-071912-abcdef"}`,
		`{"refresh_token":"refresh-XYZ789","status":401}`,
		"a valid TG-abc123def456 token was rejected",
		"Authorization=Bearer APP_USR-1234567890123456-071912-abcdef",
	}
	for _, body := range cases {
		got := sanitizeProviderBody(body)
		for _, secret := range secrets {
			if strings.Contains(got, secret) {
				t.Fatalf("secret %q leaked through sanitize:\n in:  %s\n out: %s", secret, body, got)
			}
		}
		if !strings.Contains(got, "[REDACTED]") {
			t.Fatalf("expected redaction marker in %q", got)
		}
	}

	// A benign body with no credential must survive intact (no over-redaction of
	// the diagnostic signal, e.g. an ML "status" code or error message).
	benign := sanitizeProviderBody(`{"message":"item MLB123 not found","status":404}`)
	if strings.Contains(benign, "[REDACTED]") {
		t.Fatalf("over-redacted a benign body: %s", benign)
	}
}

func TestObserveProviderFailureIgnoresBodilessErrors(t *testing.T) {
	var buf bytes.Buffer
	adapter := &marketPriceIntelCollectorAdapter{
		log: slog.New(slog.NewTextHandler(&buf, nil)),
	}
	// Timeouts and bare sentinels carry no provider body — nothing to observe.
	adapter.observeProviderFailure("GetOwnItemPricing", context.DeadlineExceeded, &marketports.ProviderStatusError{StatusCode: http.StatusGatewayTimeout})
	adapter.observeProviderFailure("GetOwnItemPricing", marketports.ErrRateLimited, marketports.ErrRateLimited)
	if buf.Len() != 0 {
		t.Fatalf("expected no log for bodiless errors, got:\n%s", buf.String())
	}
}

type fakeInstallationRepo struct {
	items []integrationsdomain.Installation
}

func (f fakeInstallationRepo) CreateInstallation(context.Context, integrationsdomain.Installation) error {
	return errors.New("not implemented")
}

func (f fakeInstallationRepo) GetInstallation(_ context.Context, installationID string) (integrationsdomain.Installation, bool, error) {
	for _, inst := range f.items {
		if inst.InstallationID == installationID {
			return inst, true, nil
		}
	}
	return integrationsdomain.Installation{}, false, nil
}

func (f fakeInstallationRepo) ListInstallations(context.Context) ([]integrationsdomain.Installation, error) {
	return append([]integrationsdomain.Installation(nil), f.items...), nil
}

func (f fakeInstallationRepo) UpdateInstallationStatus(context.Context, string, integrationsdomain.InstallationStatus, integrationsdomain.HealthStatus) error {
	return errors.New("not implemented")
}

func (f fakeInstallationRepo) ApplyConnectionSnapshot(context.Context, string, integrationsdomain.ConnectionSnapshot, string) error {
	return errors.New("not implemented")
}

func TestAccountRefForTenant(t *testing.T) {
	newAdapter := func(items ...integrationsdomain.Installation) *marketPriceIntelCollectorAdapter {
		svc := integrationsapp.NewInstallationService(fakeInstallationRepo{items: items}, "tenant_default")
		return newMarketPriceIntelCollectorAdapter(nil, svc, "tenant_default")
	}
	mlInstallation := func(id string, status integrationsdomain.InstallationStatus) integrationsdomain.Installation {
		return integrationsdomain.Installation{
			InstallationID:    id,
			TenantID:          "tenant_default",
			ProviderCode:      mercadoLivreProviderCode,
			Status:            status,
			ExternalAccountID: "acct-" + id,
		}
	}
	want503 := func(t *testing.T, err error) {
		t.Helper()
		var statusErr *marketports.ProviderStatusError
		if !errors.As(err, &statusErr) || statusErr.StatusCode != http.StatusServiceUnavailable {
			t.Fatalf("got %v, want ProviderStatusError{503}", err)
		}
	}

	t.Run("connected mercado_livre installation preferred", func(t *testing.T) {
		adapter := newAdapter(
			mlInstallation("inst-degraded", integrationsdomain.InstallationStatusDegraded),
			mlInstallation("inst-connected", integrationsdomain.InstallationStatusConnected),
		)
		ref, err := adapter.accountRefForTenant(context.Background())
		if err != nil {
			t.Fatalf("accountRefForTenant() error = %v", err)
		}
		if ref.InstallationID != "inst-connected" || ref.ProviderAccountID != "acct-inst-connected" {
			t.Fatalf("got %+v, want the Connected installation", ref)
		}
		if ref.TenantID != "tenant_default" || ref.ProviderCode != mercadoLivreProviderCode {
			t.Fatalf("got %+v, want tenant/provider populated", ref)
		}
	})

	t.Run("no connected installation is honest 503, never a fallback", func(t *testing.T) {
		adapter := newAdapter(
			mlInstallation("inst-degraded", integrationsdomain.InstallationStatusDegraded),
			mlInstallation("inst-reauth", integrationsdomain.InstallationStatusRequiresReauth),
		)
		_, err := adapter.accountRefForTenant(context.Background())
		want503(t, err)
	})

	t.Run("no mercado_livre installation at all is 503", func(t *testing.T) {
		other := integrationsdomain.Installation{
			InstallationID:    "inst-other",
			ProviderCode:      "magalu",
			Status:            integrationsdomain.InstallationStatusConnected,
			ExternalAccountID: "acct-other",
		}
		adapter := newAdapter(other)
		_, err := adapter.accountRefForTenant(context.Background())
		want503(t, err)

		empty := newAdapter()
		_, err = empty.accountRefForTenant(context.Background())
		want503(t, err)
	})
}
