package connectors

import (
	"context"
	"testing"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type fakePriceWriter struct {
	result  connectorsdomain.PriceWriteResult
	err     error
	request connectorsdomain.PriceWriteRequest
}

func (f *fakePriceWriter) UpdatePrice(_ context.Context, request connectorsdomain.PriceWriteRequest) (connectorsdomain.PriceWriteResult, error) {
	f.request = request
	return f.result, f.err
}

func TestPriceWriterResolvesExplicitCapabilityAndMapsOutcomes(t *testing.T) {
	missing := NewPriceWriter(connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre"}}), "tenant", "mercado_livre")
	out, err := missing.Apply(context.Background(), ports.WriteItem{ListingID: "inst~MLB1~-", After: []byte(`{"new_price":{"amount":"10.00","currency":"BRL"}}`)})
	if err != nil || out.Failure == nil || out.Failure.Code != domain.FailureCodeProviderValidation {
		t.Fatalf("missing outcome=%+v err=%v", out, err)
	}

	provider := &fakePriceWriter{err: connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderRateLimited, "safe provider message")}
	writer := NewPriceWriter(connectorsapp.NewMarketplaceCapabilityService([]connectorsapp.ProviderCapabilitySet{{ProviderCode: "mercado_livre", PriceWrites: provider}}), "tenant", "mercado_livre")
	out, _ = writer.Apply(context.Background(), ports.WriteItem{InstallationID: "inst", ListingID: "inst~MLB1~-", IdempotencyKey: "key", After: []byte(`{"new_price":{"amount":"10.00","currency":"BRL"}}`)})
	if out.Failure == nil || out.Failure.Code != domain.FailureCodeProviderRateLimited || out.Failure.MessageProvider != "safe provider message" || !out.Failure.Retryable {
		t.Fatalf("mapped=%+v", out.Failure)
	}

	provider.err, provider.result = nil, connectorsdomain.PriceWriteResult{Result: connectorsdomain.WriteResultRejected, Message: "safe rejection"}
	out, _ = writer.Apply(context.Background(), ports.WriteItem{InstallationID: "inst", ListingID: "inst~MLB1~-", IdempotencyKey: "key", After: []byte(`{"new_price":{"amount":"10.00","currency":"BRL"}}`)})
	if out.Failure == nil || out.Failure.Code != domain.FailureCodeProviderValidation || out.Failure.MessageProvider != "safe rejection" {
		t.Fatalf("rejected=%+v", out.Failure)
	}
	provider.result.Result, provider.result.Message = connectorsdomain.WriteResultStatus("unknown"), "must not pass"
	out, _ = writer.Apply(context.Background(), ports.WriteItem{InstallationID: "inst", ListingID: "inst~MLB1~-", IdempotencyKey: "key", After: []byte(`{"new_price":{"amount":"10.00","currency":"BRL"}}`)})
	if out.Failure == nil || out.Failure.Code != domain.FailureCodeInternal || out.Failure.MessageProvider != "unknown provider price outcome" {
		t.Fatalf("unknown=%+v", out.Failure)
	}
}
