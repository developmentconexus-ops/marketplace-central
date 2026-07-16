package connectors

import (
	"context"
	"testing"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type fakeStockWriter struct {
	result  connectorsdomain.StockWriteResult
	err     error
	request connectorsdomain.StockWriteRequest
}

func (f *fakeStockWriter) UpdateAvailableQuantity(_ context.Context, request connectorsdomain.StockWriteRequest) (connectorsdomain.StockWriteResult, error) {
	f.request = request
	return f.result, f.err
}

func TestStockWriterMapsProviderResultsAndPreservesReferences(t *testing.T) {
	provider := &fakeStockWriter{result: connectorsdomain.StockWriteResult{Result: connectorsdomain.StockWriteResultApplied}}
	writer := NewStockWriter(connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre", StockWrites: provider}}), "tenant", "mercado_livre", "account")
	out, err := writer.Apply(context.Background(), ports.WriteItem{InstallationID: "inst", ListingID: "inst~MLB1~VAR1", IdempotencyKey: "key", After: []byte(`{"publish_quantity":3,"reason":"operator"}`)})
	if err != nil || out.Failure != nil || provider.request.ProviderItemID != "MLB1" || provider.request.ProviderVariationID != "VAR1" || provider.request.RequestedQuantity != 3 {
		t.Fatalf("out=%+v request=%+v err=%v", out, provider.request, err)
	}
	provider.result = connectorsdomain.StockWriteResult{Result: connectorsdomain.StockWriteResultUnsupportedShape, Message: "safe shape"}
	out, _ = writer.Apply(context.Background(), ports.WriteItem{InstallationID: "inst", ListingID: "inst~MLB1~VAR1", After: []byte(`{"publish_quantity":3}`)})
	if out.Failure == nil || out.Failure.Code != domain.FailureCodeProviderValidation || out.Failure.MessageProvider != "safe shape" {
		t.Fatalf("unsupported=%+v", out.Failure)
	}
	missing := NewStockWriter(connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre"}}), "tenant", "mercado_livre", "account")
	out, _ = missing.Apply(context.Background(), ports.WriteItem{ListingID: "inst~MLB1~VAR1", After: []byte(`{"publish_quantity":3}`)})
	if out.Failure == nil || out.Failure.Code != domain.FailureCodeProviderValidation {
		t.Fatalf("missing=%+v", out.Failure)
	}
}
