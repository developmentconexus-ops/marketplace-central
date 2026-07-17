//go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

func TestMutationTransportErrorMatrix(t *testing.T) {
	x := newLane(t)
	assert := func(name, method, path, body string, status int, code string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			gotStatus, payload := x.request(method, path, body)
			errorBody, _ := payload["error"].(map[string]any)
			if gotStatus != status || errorBody["code"] != code {
				t.Fatalf("status=%d payload=%v", gotStatus, payload)
			}
		})
	}
	assert("invalid_body", http.MethodPost, "/mutations", "[", 400, "invalid_body")
	assert("invalid_intent", http.MethodPost, "/mutations", `{"installation_id":"i","type":"price_update","selection":{"mode":"explicit","listing_ids":["i~x~-"]},"intent":null,"actor":"a"}`, 400, "invalid_intent")
	assert("installation_required", http.MethodPost, "/mutations", `{"type":"price_update","selection":{},"intent":{},"actor":"a"}`, 400, "installation_required")
	assert("type_not_enabled", http.MethodPost, "/mutations", `{"installation_id":"i","type":"disabled","selection":{},"intent":{},"actor":"a"}`, 422, "type_not_enabled")
	t.Run("method_not_allowed", func(t *testing.T) {
		status, payload := x.commandRequest(x.commands.HandleCreate, http.MethodPut, "/mutations", "")
		if status != 405 || payload["error"].(map[string]any)["code"] != "method_not_allowed" { t.Fatalf("status=%d payload=%v", status, payload) }
	})
	assert("invalid_filter", http.MethodGet, "/mutations?installation_id=i&bogus=x", "", 400, "invalid_filter")
	assert("invalid_cursor", http.MethodGet, "/mutations?installation_id=i&cursor=bad", "", 400, "invalid_cursor")
	assert("protocol_not_found", http.MethodGet, "/mutations/missing", "", 404, "protocol_not_found")
	assert("actor_required", http.MethodPost, "/mutations", `{"installation_id":"i","type":"price_update","selection":{"mode":"explicit","listing_ids":["i~x~-"]},"intent":{},"actor":""}`, 422, "actor_required")
	assert("empty_selection", http.MethodPost, "/mutations/"+createFilterProtocol(t, x)+"/preview", "", 422, "empty_selection")

	ids := make([]string, 2001)
	for i := range ids {
		ids[i] = fmt.Sprintf("%s~MLB-%d~-", x.installation, i)
	}
	large := x.create(ids)
	assert("selection_too_large", http.MethodPost, "/mutations/"+large+"/preview", "", 422, "selection_too_large")

	now := time.Now().UTC()
	missingTimeID := x.seedListing("MLB-NO-TIME", "1.00", nil)
	assert("source_time_unavailable", http.MethodPost, "/mutations/"+x.create([]string{missingTimeID})+"/preview", "", 422, "source_time_unavailable")
	assert("internal_error", http.MethodPost, "/mutations/"+x.create([]string{"not-a-listing-id"})+"/preview", "", 500, "internal_error")

	previewed := x.create([]string{x.seedListing("MLB-STALE", "1.00", &now)})
	if status, _ := x.request(http.MethodPost, "/mutations/"+previewed+"/preview", ""); status != 200 {
		t.Fatal(status)
	}
	assert("execute_required", http.MethodPost, "/mutations/"+previewed+"/approve", `{"execute":false}`, 422, "execute_required")
	_, _ = x.pool.Exec(context.Background(), `UPDATE mutation_protocols SET previewed_at=now()-interval '16 minutes' WHERE tenant_id=$1 AND protocol_id=$2`, x.tenant, previewed)
	assert("preview_stale", http.MethodPost, "/mutations/"+previewed+"/approve", `{"execute":true}`, 409, "preview_stale")

	draft := x.create([]string{x.seedListing("MLB-DRAFT", "1.00", &now)})
	assert("invalid_state", http.MethodPost, "/mutations/"+draft+"/approve", `{"execute":true}`, 409, domain.InvalidStateCode)
	assert("nothing_to_retry", http.MethodPost, "/mutations/"+seedTerminalWithoutRetry(t, x)+"/retry-failures", "", 409, "nothing_to_retry")
}

func createFilterProtocol(t *testing.T, x *lane) string {
	t.Helper()
	status, payload := x.request(http.MethodPost, "/mutations", fmt.Sprintf(`{"installation_id":%q,"type":"price_update","selection":{"mode":"filter"},"intent":{},"actor":"operator"}`, x.installation))
	if status != 201 {
		t.Fatalf("filter create=%d %v", status, payload)
	}
	return payload["protocol_id"].(string)
}

func seedTerminalWithoutRetry(t *testing.T, x *lane) string {
	t.Helper()
	p, err := x.repo.CreateProtocol(context.Background(), ports.CreateProtocolInput{InstallationID: x.installation, Type: domain.ProtocolTypePriceUpdate, Actor: "operator", Intent: json.RawMessage(`{}`), Selection: json.RawMessage(`{"mode":"explicit","listing_ids":["x"]}`), CreatedAt: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}
	_, err = x.pool.Exec(context.Background(), `UPDATE mutation_protocols SET state='failed_preserved',finished_at=now() WHERE tenant_id=$1 AND protocol_id=$2`, x.tenant, p.ProtocolID)
	if err != nil {
		t.Fatal(err)
	}
	return p.ProtocolID
}
