package composition

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	mutationsstub "marketplace-central/apps/server_core/internal/modules/mutations/adapters/stub"
	mutationsapp "marketplace-central/apps/server_core/internal/modules/mutations/application"
	mutationsbg "marketplace-central/apps/server_core/internal/modules/mutations/background"
	"marketplace-central/apps/server_core/internal/platform/httpx"
	"marketplace-central/apps/server_core/internal/platform/pgdb"
)

func TestRefreshListingsOpenAPIContractParity(t *testing.T) {
	raw, err := os.ReadFile("../../../../contracts/api/marketplace-central.openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	spec := string(raw)
	start := strings.Index(spec, "  /listings/refresh:")
	end := strings.Index(spec[start+1:], "\n  /")
	if start < 0 || end < 0 {
		t.Fatal("/listings/refresh path not found")
	}
	path := spec[start : start+1+end]
	for _, want := range []string{"operationId: refreshListings", "$ref: '#/components/schemas/RefreshListingsRequest'", "$ref: '#/components/schemas/RefreshListingsAccepted'", `"202":`, `"400":`, `"404":`, `"409":`} {
		if !strings.Contains(path, want) {
			t.Fatalf("refresh contract missing %q", want)
		}
	}
	for _, unwanted := range []string{`"200":`, `"405":`, `"500":`, `"503":`} {
		if strings.Contains(path, unwanted) {
			t.Fatalf("refresh contract unexpectedly contains %q", unwanted)
		}
	}
}

func TestRootRuntimeWiresMutationPoller(t *testing.T) {
	runtime, err := NewRootRuntime(nil, pgdb.Config{DefaultTenantID: "tenant_default", EncryptionKey: "0123456789abcdef0123456789abcdef"})
	if err != nil {
		t.Fatalf("NewRootRuntime() error = %v", err)
	}
	if runtime.MutationPoller == nil {
		t.Fatal("MutationPoller is nil")
	}
	var _ *mutationsbg.Poller = runtime.MutationPoller
}

func TestRootRuntimeSelectsMutationWriterLane(t *testing.T) {
	for _, tc := range []struct {
		name, value string
		real        bool
	}{
		{name: "unset"},
		{name: "false", value: "false"},
		{name: "yes", value: "yes"},
		{name: "one", value: "1"},
		{name: "true", value: "true", real: true},
		{name: "trimmed uppercase true", value: "  TRUE  ", real: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MPC_PROVIDER_WRITES_ENABLED", tc.value)
			runtime, err := NewRootRuntime(nil, pgdb.Config{DefaultTenantID: "tenant_default", EncryptionKey: "0123456789abcdef0123456789abcdef"})
			if err != nil {
				t.Fatalf("NewRootRuntime() error = %v", err)
			}
			if runtime.mutationLane.envelope == nil {
				t.Fatal("inventory mutation envelope is nil")
			}
			if tc.real {
				if _, ok := runtime.mutationLane.writer.(*mutationsapp.WriterRouter); !ok {
					t.Fatalf("writer = %T, want real WriterRouter", runtime.mutationLane.writer)
				}
				if !runtime.mutationLane.real || !runtime.mutationLane.price || !runtime.mutationLane.stock || !runtime.mutationLane.listing || !runtime.mutationLane.linkage || !runtime.mutationLane.resync {
					t.Fatalf("real lane surfaces = %+v, want all configured", runtime.mutationLane)
				}
				return
			}
			if _, ok := runtime.mutationLane.writer.(*mutationsstub.Writer); !ok {
				t.Fatalf("writer = %T, want fallback stub", runtime.mutationLane.writer)
			}
			if runtime.mutationLane.real || runtime.mutationLane.price || runtime.mutationLane.stock || runtime.mutationLane.listing || runtime.mutationLane.linkage || runtime.mutationLane.resync {
				t.Fatalf("fallback lane unexpectedly exposes real surfaces: %+v", runtime.mutationLane)
			}
		})
	}
}

func TestFeeScheduleRoutesUseBatchDeadline(t *testing.T) {
	mux := httpx.NewRouteClassMux()
	registerBatchRoutes(mux)

	for _, path := range []string{"/admin/fee-schedules/sync", "/admin/fee-schedules/seed"} {
		path := path
		t.Run(path, func(t *testing.T) {
			var deadline time.Time
			var hasDeadline bool
			mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				deadline, hasDeadline = r.Context().Deadline()
				w.WriteHeader(http.StatusNoContent)
			})

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, path, nil))
			if recorder.Code != http.StatusNoContent {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
			}
			if !hasDeadline {
				t.Fatal("batch route has no context deadline")
			}
			if remaining := time.Until(deadline); remaining < 119*time.Second || remaining > 120*time.Second {
				t.Fatalf("batch deadline remaining = %s, want approximately 120s", remaining)
			}
		})
	}
}

func TestRootRuntimeRegistersListingsReadRoutes(t *testing.T) {
	runtime, err := NewRootRuntime(nil, pgdb.Config{
		DefaultTenantID: "tenant_default",
		EncryptionKey:   "0123456789abcdef0123456789abcdef",
	})
	if err != nil {
		t.Fatalf("NewRootRuntime() error = %v", err)
	}

	for _, tc := range []struct {
		path   string
		method string
		status int
	}{
		{path: "/listings", method: http.MethodGet, status: http.StatusBadRequest},
		{path: "/listings/by-product", method: http.MethodGet, status: http.StatusBadRequest},
		{path: "/listings/summary", method: http.MethodGet, status: http.StatusBadRequest},
		{path: "/listings/not-a-listing-id", method: http.MethodGet, status: http.StatusNotFound},
		{path: "/listings/refresh", method: http.MethodPost, status: http.StatusBadRequest},
	} {
		t.Run(tc.path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			runtime.Handler.ServeHTTP(recorder, httptest.NewRequest(tc.method, tc.path, nil))
			if recorder.Code != tc.status {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, tc.status, recorder.Body.String())
			}
		})
	}
}

func TestRootMountsDashboardSummary(t *testing.T) {
	runtime, err := NewRootRuntime(nil, pgdb.Config{
		DefaultTenantID: "tenant_default",
		EncryptionKey:   "0123456789abcdef0123456789abcdef",
	})
	if err != nil {
		t.Fatalf("NewRootRuntime() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/dashboard/summary", nil))
	if recorder.Code == http.StatusNotFound || recorder.Code == http.StatusMethodNotAllowed {
		t.Fatalf("dashboard summary route is not mounted: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestRootMountsSyncRuns(t *testing.T) {
	runtime, err := NewRootRuntime(nil, pgdb.Config{
		DefaultTenantID: "tenant_default",
		EncryptionKey:   "0123456789abcdef0123456789abcdef",
	})
	if err != nil {
		t.Fatalf("NewRootRuntime() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/sync/runs", nil))
	if recorder.Code == http.StatusNotFound || recorder.Code == http.StatusMethodNotAllowed {
		t.Fatalf("sync runs route is not mounted: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestRootMountsMarketObservations(t *testing.T) {
	runtime, err := NewRootRuntime(nil, pgdb.Config{
		DefaultTenantID: "tenant_default",
		EncryptionKey:   "0123456789abcdef0123456789abcdef",
	})
	if err != nil {
		t.Fatalf("NewRootRuntime() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/market/observations", nil))
	if recorder.Code == http.StatusNotFound || recorder.Code == http.StatusMethodNotAllowed {
		t.Fatalf("market observations route is not mounted: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestRootMountsMarketReferences(t *testing.T) {
	runtime, err := NewRootRuntime(nil, pgdb.Config{
		DefaultTenantID: "tenant_default",
		EncryptionKey:   "0123456789abcdef0123456789abcdef",
	})
	if err != nil {
		t.Fatalf("NewRootRuntime() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/market/references", nil))
	if recorder.Code == http.StatusNotFound || recorder.Code == http.StatusMethodNotAllowed {
		t.Fatalf("market references route is not mounted: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestRootDoesNotWireMarketCollectorOrScheduler(t *testing.T) {
	raw, err := os.ReadFile("root.go")
	if err != nil {
		t.Fatal(err)
	}
	source := string(raw)
	for _, forbidden := range []string{
		"marketapp.NewCollectionService(",
		"marketbackground",
		"marketCollector",
		"marketScheduler",
	} {
		if strings.Contains(source, forbidden) {
			t.Fatalf("root.go wires forbidden market runtime component %q", forbidden)
		}
	}
}

func TestRootMountsOrdersOnce(t *testing.T) {
	runtime, err := NewRootRuntime(nil, pgdb.Config{
		DefaultTenantID: "tenant_default",
		EncryptionKey:   "0123456789abcdef0123456789abcdef",
	})
	if err != nil {
		t.Fatalf("NewRootRuntime() error = %v", err)
	}

	listRecorder := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(listRecorder, httptest.NewRequest(http.MethodGet, "/orders?installation_id=installation-1", nil))
	if listRecorder.Code == http.StatusNotFound || listRecorder.Code == http.StatusMethodNotAllowed {
		t.Fatalf("orders list route is not mounted: status=%d body=%s", listRecorder.Code, listRecorder.Body.String())
	}

	detailRecorder := httptest.NewRecorder()
	runtime.Handler.ServeHTTP(detailRecorder, httptest.NewRequest(http.MethodGet, "/orders/provider-order-1?installation_id=installation-1", nil))
	if detailRecorder.Code != http.StatusNotFound || !strings.Contains(detailRecorder.Body.String(), "order_not_found") {
		t.Fatalf("orders detail route was not handled by the orders transport: status=%d body=%s", detailRecorder.Code, detailRecorder.Body.String())
	}
}

func TestRootRuntimeRegistersMutationRoutes(t *testing.T) {
	runtime, err := NewRootRuntime(nil, pgdb.Config{DefaultTenantID: "tenant_default", EncryptionKey: "0123456789abcdef0123456789abcdef"})
	if err != nil {
		t.Fatal(err)
	}
	for _, route := range []struct{ method, path, body string }{
		{http.MethodPost, "/mutations", "["},
		{http.MethodPost, "/mutations/missing/preview", "["},
		{http.MethodPost, "/mutations/missing/approve", "["},
		{http.MethodPost, "/mutations/missing/cancel", "["},
		{http.MethodPost, "/mutations/missing/retry-failures", "["},
		{http.MethodGet, "/mutations", ""},
		{http.MethodPost, "/mutations/missing", ""},
		{http.MethodPost, "/mutations/missing/items", ""},
	} {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			runtime.Handler.ServeHTTP(recorder, httptest.NewRequest(route.method, route.path, strings.NewReader(route.body)))
			if recorder.Code == http.StatusNotFound && !strings.Contains(recorder.Body.String(), "protocol_not_found") {
				t.Fatalf("route was not registered: status=%d body=%s", recorder.Code, recorder.Body.String())
			}
		})
	}
}
