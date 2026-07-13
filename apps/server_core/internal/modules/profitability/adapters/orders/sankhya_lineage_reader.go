package orders

import (
	"context"
	"strings"

	ordersapp "marketplace-central/apps/server_core/internal/modules/orders/application"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	profitabilityports "marketplace-central/apps/server_core/internal/modules/profitability/ports"
)

type assistedLineageResolver interface {
	ResolveCurrentLineage(context.Context, ordersapp.ResolveCurrentSankhyaLineageInput) (ordersdomain.AssistedSankhyaLineage, error)
}

type SankhyaLineageReader struct {
	resolver assistedLineageResolver
	tenantID string
}

func NewSankhyaLineageReader(resolver assistedLineageResolver, tenantID string) SankhyaLineageReader {
	return SankhyaLineageReader{resolver: resolver, tenantID: strings.TrimSpace(tenantID)}
}

func (r SankhyaLineageReader) ResolveCurrentLineage(ctx context.Context, installationID, providerOrderID, mpcLineID string) (profitabilityports.SankhyaLineage, error) {
	lineage, err := r.resolver.ResolveCurrentLineage(ctx, ordersapp.ResolveCurrentSankhyaLineageInput{
		TenantID:        r.tenantID,
		InstallationID:  strings.TrimSpace(installationID),
		ProviderOrderID: strings.TrimSpace(providerOrderID),
		MPCLineID:       ordersdomain.MPCLineID(mpcLineID),
	})
	if err != nil {
		return profitabilityports.SankhyaLineage{}, err
	}
	result := profitabilityports.SankhyaLineage{State: mapSankhyaLineageState(lineage.State)}
	result.Descendants = make([]profitabilityports.SankhyaTaxSource, len(lineage.Descendants))
	for index, descendant := range lineage.Descendants {
		result.Descendants[index] = profitabilityports.SankhyaTaxSource{
			DocumentID: descendant.Identity.DocumentID,
			LineNumber: descendant.Identity.LineNumber,
		}
	}
	return result, nil
}

func mapSankhyaLineageState(state ordersdomain.AssistedSankhyaLineageState) profitabilityports.SankhyaLineageState {
	switch state {
	case ordersdomain.AssistedSankhyaLineageNone:
		return profitabilityports.SankhyaLineageNone
	case ordersdomain.AssistedSankhyaLineagePartial:
		return profitabilityports.SankhyaLineagePartial
	case ordersdomain.AssistedSankhyaLineageComplete:
		return profitabilityports.SankhyaLineageComplete
	case ordersdomain.AssistedSankhyaLineageConflict:
		return profitabilityports.SankhyaLineageConflict
	case ordersdomain.AssistedSankhyaLineageUnavailable:
		return profitabilityports.SankhyaLineageUnavailable
	default:
		return profitabilityports.SankhyaLineageConflict
	}
}
