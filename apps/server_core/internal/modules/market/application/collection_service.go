package application

import (
	"context"
	"errors"
	"fmt"

	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

var ErrCollectionNotConfigured = errors.New("MARKET_COLLECTION_NOT_CONFIGURED")

// CollectionService performs one-shot collection and persistence of market
// evidence. Collector implementations belong outside this package; the
// service only depends on the collector and append-only store ports.
type CollectionService struct {
	collector        ports.CollectorPort
	observationStore ports.ObservationStore
	referenceStore   ports.ReferenceStore
}

func NewCollectionService(
	collector ports.CollectorPort,
	observationStore ports.ObservationStore,
	referenceStore ports.ReferenceStore,
) *CollectionService {
	return &CollectionService{
		collector:        collector,
		observationStore: observationStore,
		referenceStore:   referenceStore,
	}
}

func (s *CollectionService) CollectObservations(ctx context.Context, installationID string, listingIDs []string) error {
	if s == nil || s.collector == nil || s.observationStore == nil {
		return ErrCollectionNotConfigured
	}

	observations, err := s.collector.CollectObservations(ctx, installationID, listingIDs)
	if err != nil {
		return err
	}
	for index, observation := range observations {
		if err := observation.ValidateStored(); err != nil {
			return fmt.Errorf("observation %d: %w", index, err)
		}
	}
	if len(observations) == 0 {
		return nil
	}
	return s.observationStore.AppendObservations(ctx, observations)
}

func (s *CollectionService) CollectReferences(ctx context.Context, productIDs []string) error {
	if s == nil || s.collector == nil || s.referenceStore == nil {
		return ErrCollectionNotConfigured
	}

	references, err := s.collector.CollectReferences(ctx, productIDs)
	if err != nil {
		return err
	}
	for index, reference := range references {
		if err := reference.ValidateStored(); err != nil {
			return fmt.Errorf("reference %d: %w", index, err)
		}
	}
	if len(references) == 0 {
		return nil
	}
	return s.referenceStore.AppendReferences(ctx, references)
}
