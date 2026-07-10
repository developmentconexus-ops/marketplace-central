package productlinks

import (
	"context"

	inventorydomain "marketplace-central/apps/server_core/internal/modules/inventory/domain"
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

func (r LinkReader) ListLinks(ctx context.Context, installationID string, limit int) ([]inventorydomain.LinkRecord, error) {
	links, err := r.workflows.ListProductLinks(ctx, installationID, limit)
	if err != nil {
		return nil, err
	}
	candidates, err := r.candidates.ListLinkCandidates(ctx, installationID, limit)
	if err != nil {
		return nil, err
	}
	records := map[string]inventorydomain.LinkRecord{}
	for _, candidate := range candidates {
		key := linkKey(candidate.InstallationID, candidate.ProviderItemID, candidate.ProviderVariationID)
		current := records[key]
		if current.Identity.InstallationID == "" {
			current.Identity = inventorydomain.ListingIdentity{
				InstallationID:      candidate.InstallationID,
				ProviderItemID:      candidate.ProviderItemID,
				ProviderVariationID: candidate.ProviderVariationID,
			}
			current.ProviderCode = candidate.ProviderCode
			current.State = inferCandidateLinkState(candidate.State)
			current.InternalProductID = candidate.InternalProductID
			current.InternalProductName = candidate.InternalProductName
			current.InternalReferenceCode = candidate.InternalReferenceCode
			records[key] = current
			continue
		}
		if inferCandidateLinkState(candidate.State) == inventorydomain.LinkStateConflict {
			current.State = inventorydomain.LinkStateConflict
			records[key] = current
		}
	}
	for _, link := range links {
		key := linkKey(link.InstallationID, link.ProviderItemID, link.ProviderVariationID)
		records[key] = inventorydomain.LinkRecord{
			Identity: inventorydomain.ListingIdentity{
				InstallationID:      link.InstallationID,
				ProviderItemID:      link.ProviderItemID,
				ProviderVariationID: link.ProviderVariationID,
			},
			ProviderCode:          link.ProviderCode,
			State:                 mapProductLinkState(link.State),
			InternalProductID:     link.InternalProductID,
			InternalProductName:   link.InternalProductName,
			InternalReferenceCode: link.InternalReferenceCode,
		}
	}
	result := make([]inventorydomain.LinkRecord, 0, len(records))
	for _, record := range records {
		result = append(result, record)
	}
	return result, nil
}

func inferCandidateLinkState(state productlinksdomain.LinkCandidateState) inventorydomain.LinkState {
	switch state {
	case productlinksdomain.LinkCandidateStateConflict:
		return inventorydomain.LinkStateConflict
	case productlinksdomain.LinkCandidateStateExactEAN,
		productlinksdomain.LinkCandidateStateExactSKU,
		productlinksdomain.LinkCandidateStateTitleMatch,
		productlinksdomain.LinkCandidateStateManual,
		productlinksdomain.LinkCandidateStateUnresolved:
		return inventorydomain.LinkStateUnresolved
	default:
		return inventorydomain.LinkStateNone
	}
}

func mapProductLinkState(state productlinksdomain.ProductLinkState) inventorydomain.LinkState {
	switch state {
	case productlinksdomain.ProductLinkStateResolved:
		return inventorydomain.LinkStateResolved
	case productlinksdomain.ProductLinkStateRejected:
		return inventorydomain.LinkStateRejected
	case productlinksdomain.ProductLinkStateConflict:
		return inventorydomain.LinkStateConflict
	case productlinksdomain.ProductLinkStateUnresolved:
		return inventorydomain.LinkStateUnresolved
	default:
		return inventorydomain.LinkStateNone
	}
}

func linkKey(installationID, providerItemID, providerVariationID string) string {
	return installationID + "|" + providerItemID + "|" + providerVariationID
}
