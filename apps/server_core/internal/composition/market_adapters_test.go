package composition

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	marketports "marketplace-central/apps/server_core/internal/modules/market/ports"
)

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
