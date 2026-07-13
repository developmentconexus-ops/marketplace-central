package orders

import (
	"context"
	"reflect"
	"testing"

	ordersapp "marketplace-central/apps/server_core/internal/modules/orders/application"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	profitabilityports "marketplace-central/apps/server_core/internal/modules/profitability/ports"
)

type stubAssistedLineageResolver struct {
	input   ordersapp.ResolveCurrentSankhyaLineageInput
	lineage ordersdomain.AssistedSankhyaLineage
}

func (s *stubAssistedLineageResolver) ResolveCurrentLineage(_ context.Context, input ordersapp.ResolveCurrentSankhyaLineageInput) (ordersdomain.AssistedSankhyaLineage, error) {
	s.input = input
	return s.lineage, nil
}

func TestSankhyaLineageReaderMapsExactScopeStateAndDescendants(t *testing.T) {
	resolver := &stubAssistedLineageResolver{lineage: ordersdomain.AssistedSankhyaLineage{
		State: ordersdomain.AssistedSankhyaLineagePartial,
		Descendants: []ordersdomain.AssistedSankhyaDescendant{
			{Identity: ordersdomain.InternalDocumentLineIdentity{DocumentID: 30602, LineNumber: 2}},
			{Identity: ordersdomain.InternalDocumentLineIdentity{DocumentID: 30601, LineNumber: 1}},
		},
	}}
	reader := NewSankhyaLineageReader(resolver, " tenant-1 ")

	got, err := reader.ResolveCurrentLineage(context.Background(), " install-1 ", " order-1 ", "mpl_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err != nil {
		t.Fatalf("ResolveCurrentLineage() error = %v", err)
	}
	wantInput := ordersapp.ResolveCurrentSankhyaLineageInput{
		TenantID: "tenant-1", InstallationID: "install-1", ProviderOrderID: "order-1",
		MPCLineID: ordersdomain.MPCLineID("mpl_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
	}
	if resolver.input != wantInput {
		t.Fatalf("input = %#v, want %#v", resolver.input, wantInput)
	}
	want := profitabilityports.SankhyaLineage{State: profitabilityports.SankhyaLineagePartial, Descendants: []profitabilityports.SankhyaTaxSource{
		{DocumentID: 30602, LineNumber: 2}, {DocumentID: 30601, LineNumber: 1},
	}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("lineage = %#v, want %#v", got, want)
	}
}
