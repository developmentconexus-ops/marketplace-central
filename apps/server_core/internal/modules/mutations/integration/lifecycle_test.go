//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	listingspostgres "marketplace-central/apps/server_core/internal/modules/listings/adapters/postgres"
	listingsapp "marketplace-central/apps/server_core/internal/modules/listings/application"
	mutationslistings "marketplace-central/apps/server_core/internal/modules/mutations/adapters/listings"
	mutationspostgres "marketplace-central/apps/server_core/internal/modules/mutations/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/mutations/adapters/stub"
	"marketplace-central/apps/server_core/internal/modules/mutations/application"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
	"marketplace-central/apps/server_core/internal/modules/mutations/transport"
	"marketplace-central/apps/server_core/internal/platform/httpx"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

type lane struct {
	t                    *testing.T
	pool                 *pgxpool.Pool
	tenant, installation string
	h                    http.Handler
	commands             transport.Handler
	repo                 *mutationspostgres.Repository
	writer               *stub.Writer
}

func newLane(t *testing.T) *lane {
	t.Helper()
	if os.Getenv("MPC_TEST_DATABASE_URL") == "" {
		t.Skip("MPC_TEST_DATABASE_URL is unset")
	}
	pool, _ := testpostgres.OpenPool(t, "mutation-lifecycle")
	now := time.Now().UTC()
	x := &lane{t: t, pool: pool, tenant: fmt.Sprintf("mutation-http-%d", now.UnixNano()), installation: "inst-http", writer: stub.NewWriter(map[string]stub.Result{})}
	x.repo = mutationspostgres.NewRepository(pool, x.tenant)
	listingRepo := listingspostgres.NewRepository(pool, x.tenant)
	reader := listingsapp.NewReadService(listingRepo, nil, nil, nil, time.Now)
	service := application.NewService(x.repo, mutationslistings.NewSelectionResolver(reader), reader, time.Now)
	mux := httpx.NewRouter()
	h := transport.NewHandler(service, application.NewApprovalService(x.repo, time.Now), application.NewRetryService(x.repo, time.Now))
	x.commands = h
	h.Register(mux)
	mux.HandleFunc("POST /mutations/{id}/retry-failures", h.HandleRetry)
	transport.NewQueryHandler(application.NewReadService(x.repo)).Register(mux)
	x.h = mux
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM listings WHERE tenant_id=$1`, x.tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM mutation_protocols WHERE tenant_id=$1`, x.tenant)
	})
	return x
}

func (x *lane) commandRequest(handler http.HandlerFunc, method, path, body string) (int, map[string]any) {
	x.t.Helper()
	r := httptest.NewRecorder()
	handler.ServeHTTP(r, httptest.NewRequest(method, path, bytes.NewBufferString(body)))
	var payload map[string]any
	_ = json.Unmarshal(r.Body.Bytes(), &payload)
	return r.Code, payload
}

func (x *lane) request(method, path, body string) (int, map[string]any) {
	x.t.Helper()
	r := httptest.NewRecorder()
	x.h.ServeHTTP(r, httptest.NewRequest(method, path, bytes.NewBufferString(body)))
	var payload map[string]any
	_ = json.Unmarshal(r.Body.Bytes(), &payload)
	return r.Code, payload
}

func (x *lane) seedListing(item, price string, fetched *time.Time) string {
	x.t.Helper()
	_, err := x.pool.Exec(context.Background(), `INSERT INTO listings
		(tenant_id,installation_id,provider,provider_listing_id,variation_id,title,status,price_amount,price_currency,sync_state,created_at,updated_at,fetched_at)
		VALUES ($1,$2,'mercado_livre',$3,'-',$3,'active',$4,'BRL','synced',now(),now(),$5)`, x.tenant, x.installation, item, price, fetched)
	if err != nil {
		x.t.Fatal(err)
	}
	return x.installation + "~" + item + "~-"
}

func (x *lane) create(ids []string) string {
	x.t.Helper()
	selection, _ := json.Marshal(map[string]any{"mode": "explicit", "listing_ids": ids})
	body := fmt.Sprintf(`{"installation_id":%q,"type":"price_update","selection":%s,"intent":{"new_price":{"amount":"49.90","currency":"BRL"}},"actor":"operator"}`, x.installation, selection)
	status, payload := x.request(http.MethodPost, "/mutations", body)
	if status != http.StatusCreated {
		x.t.Fatalf("create status=%d payload=%v", status, payload)
	}
	return payload["protocol_id"].(string)
}

func TestStubLaneLifecycleSnapshotRetryAndFailureCodes(t *testing.T) {
	x := newLane(t)
	now := time.Now().UTC()
	ids := []string{x.seedListing("MLB-1", "89.00", &now)}
	for i := range []domain.FailureCode{domain.FailureCodeProviderValidation, domain.FailureCodeProviderRateLimited, domain.FailureCodeProviderUnavailable, domain.FailureCodeProviderAuth, domain.FailureCodeListingPausedRemote, domain.FailureCodeLinkUnresolved, domain.FailureCodePolicyMissing, domain.FailureCodeSKUInvariantViolation, domain.FailureCodeStaleSource, domain.FailureCodeConflictRemoteChanged, domain.FailureCodeTypeNotEnabled, domain.FailureCodeInternal} {
		ids = append(ids, x.seedListing(fmt.Sprintf("MLB-F-%02d", i), "10.00", &now))
	}
	id := x.create(ids)
	if status, _ := x.request(http.MethodPost, "/mutations/"+id+"/preview", ""); status != http.StatusOK {
		t.Fatalf("preview status=%d", status)
	}
	_, err := x.pool.Exec(context.Background(), `UPDATE listings SET price_amount=999 WHERE tenant_id=$1 AND installation_id=$2`, x.tenant, x.installation)
	if err != nil {
		t.Fatal(err)
	}
	var before string
	if err := x.pool.QueryRow(context.Background(), `SELECT before->'price'->>'amount' FROM mutation_items WHERE tenant_id=$1 AND protocol_id=$2 AND seq=1`, x.tenant, id).Scan(&before); err != nil || before != "89.00" {
		t.Fatalf("immutable before=%q err=%v", before, err)
	}
	codes := []domain.FailureCode{domain.FailureCodeProviderValidation, domain.FailureCodeProviderRateLimited, domain.FailureCodeProviderUnavailable, domain.FailureCodeProviderAuth, domain.FailureCodeListingPausedRemote, domain.FailureCodeLinkUnresolved, domain.FailureCodePolicyMissing, domain.FailureCodeSKUInvariantViolation, domain.FailureCodeStaleSource, domain.FailureCodeConflictRemoteChanged, domain.FailureCodeTypeNotEnabled, domain.FailureCodeInternal}
	program := map[string]stub.Result{}
	for i, code := range codes {
		program[id+":"+ids[i+1]] = stub.Result{Failure: &domain.Failure{Code: code, MessagePT: string(code), Retryable: code.RetryableByDefault()}}
	}
	x.writer = stub.NewWriter(program)
	if status, _ := x.request(http.MethodPost, "/mutations/"+id+"/approve", `{"execute":true}`); status != http.StatusOK {
		t.Fatalf("approve status=%d", status)
	}
	if worked, err := application.NewPoller(x.repo, x.writer, time.Now).Pass(context.Background(), x.installation); err != nil || !worked {
		t.Fatalf("poll worked=%v err=%v", worked, err)
	}
	status, terminal := x.request(http.MethodGet, "/mutations/"+id, "")
	if status != http.StatusOK || terminal["state"] != string(domain.ProtocolStatePartiallyFailed) {
		t.Fatalf("terminal=%v status=%d", terminal, status)
	}
	status, itemPage := x.request(http.MethodGet, "/mutations/"+id+"/items?limit=200", "")
	items := itemPage["items"].([]any)
	if status != http.StatusOK || len(items) != len(ids) {
		t.Fatalf("items status=%d payload=%v", status, itemPage)
	}
	seen := map[string]bool{}
	for _, raw := range items[1:] {
		failure := raw.(map[string]any)["failure"].(map[string]any)
		seen[failure["code"].(string)] = true
	}
	for _, code := range codes {
		if !seen[string(code)] {
			t.Errorf("failure code %s absent from HTTP item transcript", code)
		}
	}
	status, retried := x.request(http.MethodPost, "/mutations/"+id+"/retry-failures", "")
	if status != http.StatusOK || retried["state"] != string(domain.ProtocolStateDraft) || retried["retried_from"] != id {
		t.Fatalf("retry status=%d payload=%v", status, retried)
	}
	var retryItems int
	if err := x.pool.QueryRow(context.Background(), `SELECT count(*) FROM mutation_items WHERE tenant_id=$1 AND protocol_id=$2`, x.tenant, retried["protocol_id"]).Scan(&retryItems); err != nil || retryItems != 0 {
		t.Fatalf("retry items=%d err=%v", retryItems, err)
	}

	foreign := mutationspostgres.NewRepository(x.pool, x.tenant+"-other")
	if _, err := foreign.ReplacePreview(context.Background(), id, []ports.ReplaceItemInput{{ListingID: ids[0], After: json.RawMessage(`{}`)}}, now, now); err == nil {
		t.Fatal("cross-tenant ReplacePreview succeeded")
	}
}
