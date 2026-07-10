package productlinks

import (
	"context"

	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
	productlinksports "marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

type LinkReader struct {
	workflows  productlinksports.ProductLinkWorkflowStore
	candidates productlinksports.LinkCandidateStore
}

func NewLinkReader(workflows productlinksports.ProductLinkWorkflowStore, candidates productlinksports.LinkCandidateStore) LinkReader {
	return LinkReader{workflows: workflows, candidates: candidates}
}

func (r LinkReader) ResolveLinks(ctx context.Context, installationID string, identities []ordersdomain.ListingIdentity) (map[string]ordersdomain.ItemLink, error) {
	links, err := r.workflows.ListProductLinks(ctx, installationID, max(50, len(identities)*2))
	if err != nil {
		return nil, err
	}
	candidates, err := r.candidates.ListLinkCandidates(ctx, installationID, max(50, len(identities)*2))
	if err != nil {
		return nil, err
	}
	result := map[string]ordersdomain.ItemLink{}
	for _, candidate := range candidates {
		key := linkKey(candidate.InstallationID, candidate.ProviderItemID, candidate.ProviderVariationID)
		current := result[key]
		if current.Quality == ordersdomain.LinkQualityResolved || current.Quality == ordersdomain.LinkQualityRejected {
			continue
		}
		result[key] = ordersdomain.ItemLink{
			Identity: ordersdomain.ListingIdentity{
				InstallationID:      candidate.InstallationID,
				ProviderItemID:      candidate.ProviderItemID,
				ProviderVariationID: candidate.ProviderVariationID,
			},
			Quality: mapCandidateState(candidate.State),
		}
	}
	for _, link := range links {
		key := linkKey(link.InstallationID, link.ProviderItemID, link.ProviderVariationID)
		result[key] = ordersdomain.ItemLink{
			Identity: ordersdomain.ListingIdentity{
				InstallationID:      link.InstallationID,
				ProviderItemID:      link.ProviderItemID,
				ProviderVariationID: link.ProviderVariationID,
			},
			Quality:           mapProductLinkState(link.State),
			InternalProductID: resolvedInternalProductID(link),
		}
	}
	return result, nil
}

func resolvedInternalProductID(link productlinksdomain.ProductLink) *int {
	if link.State != productlinksdomain.ProductLinkStateResolved {
		return nil
	}
	return link.InternalProductID
}

func mapCandidateState(state productlinksdomain.LinkCandidateState) ordersdomain.LinkQuality {
	switch state {
	case productlinksdomain.LinkCandidateStateConflict:
		return ordersdomain.LinkQualityConflict
	case productlinksdomain.LinkCandidateStateExactEAN,
		productlinksdomain.LinkCandidateStateExactSKU,
		productlinksdomain.LinkCandidateStateTitleMatch,
		productlinksdomain.LinkCandidateStateManual,
		productlinksdomain.LinkCandidateStateUnresolved:
		return ordersdomain.LinkQualityUnresolved
	default:
		return ordersdomain.LinkQualityMissing
	}
}

func mapProductLinkState(state productlinksdomain.ProductLinkState) ordersdomain.LinkQuality {
	switch state {
	case productlinksdomain.ProductLinkStateResolved:
		return ordersdomain.LinkQualityResolved
	case productlinksdomain.ProductLinkStateRejected:
		return ordersdomain.LinkQualityRejected
	case productlinksdomain.ProductLinkStateConflict:
		return ordersdomain.LinkQualityConflict
	case productlinksdomain.ProductLinkStateUnresolved:
		return ordersdomain.LinkQualityUnresolved
	default:
		return ordersdomain.LinkQualityMissing
	}
}

func linkKey(installationID, providerItemID, providerVariationID string) string {
	return installationID + "|" + providerItemID + "|" + providerVariationID
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
