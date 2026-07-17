package application

import (
	"context"
	"errors"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

const MaxReadIDs = 200

var ErrTooManyIDs = errors.New("too_many_ids")

// TooManyIDsError reports a read request that exceeds the batch cap.
// Transport layers can map Code() to the public 422 error envelope.
type TooManyIDsError struct {
	Count int
	Limit int
}

func (e *TooManyIDsError) Error() string { return ErrTooManyIDs.Error() }

func (e *TooManyIDsError) Code() string { return ErrTooManyIDs.Error() }

func (e *TooManyIDsError) Unwrap() error { return ErrTooManyIDs }

type ObservationReadResult struct {
	Items []domain.MarketObservation
	AsOf  time.Time
}

type ReferenceReadResult struct {
	Items []domain.MarketReference
	AsOf  time.Time
}

// ReadService serves stored latest market rows and honest synthetic empty
// items for requested IDs without stored evidence.
type ReadService struct {
	observationStore ports.ObservationStore
	referenceStore   ports.ReferenceStore
	now              func() time.Time
}

func NewReadService(
	observationStore ports.ObservationStore,
	referenceStore ports.ReferenceStore,
	now func() time.Time,
) *ReadService {
	if now == nil {
		now = time.Now
	}
	return &ReadService{
		observationStore: observationStore,
		referenceStore:   referenceStore,
		now:              now,
	}
}

func (s *ReadService) ListObservations(ctx context.Context, listingIDs []string) (ObservationReadResult, error) {
	result := ObservationReadResult{AsOf: s.now().UTC()}
	if len(listingIDs) > MaxReadIDs {
		return ObservationReadResult{}, &TooManyIDsError{Count: len(listingIDs), Limit: MaxReadIDs}
	}
	result.Items = make([]domain.MarketObservation, 0, len(listingIDs))
	if len(listingIDs) == 0 {
		return result, nil
	}

	rows, err := s.observationStore.LatestObservations(ctx, listingIDs)
	if err != nil {
		return ObservationReadResult{}, err
	}
	byListingID := make(map[string]domain.MarketObservation, len(rows))
	for _, row := range rows {
		byListingID[row.ListingID] = row
	}
	for _, listingID := range listingIDs {
		if row, ok := byListingID[listingID]; ok {
			result.Items = append(result.Items, row)
			continue
		}
		result.Items = append(result.Items, noPriceEvidenceObservation(listingID))
	}
	return result, nil
}

func (s *ReadService) ListReferences(ctx context.Context, productIDs []string) (ReferenceReadResult, error) {
	result := ReferenceReadResult{AsOf: s.now().UTC()}
	if len(productIDs) > MaxReadIDs {
		return ReferenceReadResult{}, &TooManyIDsError{Count: len(productIDs), Limit: MaxReadIDs}
	}
	result.Items = make([]domain.MarketReference, 0, len(productIDs))
	if len(productIDs) == 0 {
		return result, nil
	}

	rows, err := s.referenceStore.LatestReferences(ctx, productIDs)
	if err != nil {
		return ReferenceReadResult{}, err
	}
	byProductID := make(map[string]domain.MarketReference, len(rows))
	for _, row := range rows {
		byProductID[row.ProductID] = row
	}
	for _, productID := range productIDs {
		if row, ok := byProductID[productID]; ok {
			result.Items = append(result.Items, row)
			continue
		}
		result.Items = append(result.Items, noPriceEvidenceReference(productID))
	}
	return result, nil
}

func noPriceEvidenceObservation(listingID string) domain.MarketObservation {
	return domain.MarketObservation{
		ListingID:     listingID,
		EvidenceState: domain.EvidenceStateNoPriceEvidence,
	}
}

func noPriceEvidenceReference(productID string) domain.MarketReference {
	return domain.MarketReference{
		ProductID:     productID,
		MatchState:    domain.MatchStateNoCandidate,
		EvidenceState: domain.EvidenceStateNoPriceEvidence,
	}
}
