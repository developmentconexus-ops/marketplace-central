package connectors

import (
	"context"
	"testing"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type fakeListingWriter struct {
	result   connectorsdomain.ListingWriteResult
	err      error
	requests []connectorsdomain.ListingWriteRequest
}

func (f *fakeListingWriter) UpdateListing(_ context.Context, request connectorsdomain.ListingWriteRequest) (connectorsdomain.ListingWriteResult, error) {
	f.requests = append(f.requests, request)
	return f.result, f.err
}

func TestListingWriterMapsPauseAndEditRejections(t *testing.T) {
	provider := &fakeListingWriter{result: connectorsdomain.ListingWriteResult{Result: connectorsdomain.WriteResultRejected, Message: "safe paused"}}
	writer := NewListingWriter(connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre", ListingWrites: provider}}), "tenant", "mercado_livre")
	out, err := writer.Apply(context.Background(), ports.WriteItem{ProtocolType: domain.ProtocolTypeListingPause, InstallationID: "inst", ListingID: "inst~MLB1~-", IdempotencyKey: "key", After: []byte(`{}`)})
	if err != nil || out.Failure == nil || out.Failure.Code != domain.FailureCodeListingPausedRemote || out.Failure.MessageProvider != "safe paused" {
		t.Fatalf("pause=%+v err=%v", out.Failure, err)
	}
	provider.result = connectorsdomain.ListingWriteResult{Result: connectorsdomain.WriteResultRejected, Message: "safe edit"}
	out, _ = writer.Apply(context.Background(), ports.WriteItem{ProtocolType: domain.ProtocolTypeListingEdit, InstallationID: "inst", ListingID: "inst~MLB1~-", After: []byte(`{"attributes":[{"id":"TITLE","value_name":"x"}]}`)})
	if out.Failure == nil || out.Failure.Code != domain.FailureCodeProviderValidation || provider.requests[1].Action != connectorsdomain.ListingWriteEdit || len(provider.requests[1].Attributes) != 1 {
		t.Fatalf("edit=%+v requests=%+v", out.Failure, provider.requests)
	}
	missing := NewListingWriter(connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre"}}), "tenant", "mercado_livre")
	out, _ = missing.Apply(context.Background(), ports.WriteItem{ProtocolType: domain.ProtocolTypeListingPause, ListingID: "inst~MLB1~-"})
	if out.Failure == nil || out.Failure.Code != domain.FailureCodeProviderValidation {
		t.Fatalf("missing=%+v", out.Failure)
	}
}
