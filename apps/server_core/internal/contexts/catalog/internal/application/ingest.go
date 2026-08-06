package application

import (
	"context"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/domain"
)

// Service is Catalog's ingestion use case.
type Service struct {
	store Store
	ids   IDFactory
}

// NewService wires the use case.
func NewService(store Store, ids IDFactory) *Service {
	return &Service{store: store, ids: ids}
}

// Ingest folds one source observation into the catalogue.
//
// Resolution is by source key and by nothing else. Identifiers are looked up
// only to REPORT a conflict, never to decide identity: an EAN that two ERP
// codes share is bad master data, and folding the second product into the first
// destroys the distinction with no way back.
func (s *Service) Ingest(ctx context.Context, o contracts.ProductObservation) (contracts.IngestResult, error) {
	if err := o.Validate(); err != nil {
		return contracts.IngestResult{}, err
	}

	existing, found, err := s.store.BySourceKey(ctx, o.Key)
	if err != nil {
		return contracts.IngestResult{}, fmt.Errorf("catalog: resolve source key: %w", err)
	}

	if found {
		next, disposition, applyErr := existing.Apply(o)
		if applyErr != nil {
			return contracts.IngestResult{}, applyErr
		}
		if disposition == contracts.DispositionChanged {
			if err := s.store.Update(ctx, next); err != nil {
				return contracts.IngestResult{}, fmt.Errorf("catalog: update product: %w", err)
			}
		}
		return contracts.IngestResult{
			ProductID:   next.ID().String(),
			Disposition: disposition,
			Version:     next.Version(),
		}, nil
	}

	duplicates, err := s.duplicateIdentifiers(ctx, o)
	if err != nil {
		return contracts.IngestResult{}, err
	}

	id, err := s.ids.NewProductID()
	if err != nil {
		return contracts.IngestResult{}, fmt.Errorf("catalog: mint product id: %w", err)
	}
	product, err := domain.NewProduct(id, o)
	if err != nil {
		return contracts.IngestResult{}, err
	}
	if err := s.store.Insert(ctx, product); err != nil {
		return contracts.IngestResult{}, fmt.Errorf("catalog: insert product: %w", err)
	}

	return contracts.IngestResult{
		ProductID:            product.ID().String(),
		Disposition:          contracts.DispositionCreated,
		Version:              product.Version(),
		DuplicateIdentifiers: duplicates,
	}, nil
}

// duplicateIdentifiers reports which of the observation's identifiers already
// belong to another product OF THE SAME TENANT.
func (s *Service) duplicateIdentifiers(ctx context.Context, o contracts.ProductObservation) ([]contracts.Identifier, error) {
	var out []contracts.Identifier
	for _, id := range o.Identifiers {
		matches, err := s.store.ByIdentifier(ctx, o.Key.Tenant(), id)
		if err != nil {
			return nil, fmt.Errorf("catalog: look up identifier %s: %w", id, err)
		}
		if len(matches) > 0 {
			out = append(out, id)
		}
	}
	return out, nil
}
