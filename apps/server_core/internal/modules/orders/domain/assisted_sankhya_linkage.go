package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidAssistedSankhyaLinkage = errors.New("invalid assisted sankhya linkage")
	ErrAssistedSankhyaOrderNotFound  = errors.New("assisted sankhya order not found")
	ErrAssistedSankhyaReadFailed     = errors.New("assisted sankhya read failed")
)

const AssistedSankhyaActorOperatorSuppliedUnverified = "operator_supplied_unverified"

type AssistedSankhyaReadErrorKind string

const (
	AssistedSankhyaReadConfigurationInvalid AssistedSankhyaReadErrorKind = "configuration_invalid"
	AssistedSankhyaReadCandidateConflict    AssistedSankhyaReadErrorKind = "candidate_conflict"
	AssistedSankhyaReadLineageConflict      AssistedSankhyaReadErrorKind = "lineage_conflict"
	AssistedSankhyaReadUnavailable          AssistedSankhyaReadErrorKind = "unavailable"
)

type AssistedSankhyaReadError struct {
	Kind AssistedSankhyaReadErrorKind
}

func (e *AssistedSankhyaReadError) Error() string {
	return fmt.Sprintf("%s: %s", ErrAssistedSankhyaReadFailed, e.Kind)
}

func (e *AssistedSankhyaReadError) Unwrap() error { return ErrAssistedSankhyaReadFailed }

type AssistedSankhyaCandidateLine struct {
	Identity  InternalDocumentLineIdentity
	ProductID *int64
	Quantity  *float64
	Amount    *float64
}

type AssistedSankhyaCandidate struct {
	Header         InternalDocumentIdentity
	DocumentNumber *int64
	OperationCode  int64
	NegotiatedAt   *time.Time
	Amount         *float64
	Lines          []AssistedSankhyaCandidateLine
}

type AssistedSankhyaCandidateResult struct {
	Scope            LinkageScope
	ExternalOrderKey string
	Candidates       []AssistedSankhyaCandidate
}

type AssistedSankhyaLineSelection struct {
	MPCLineID     MPCLineID
	CandidateLine InternalDocumentLineIdentity
}

type AssistedSankhyaLineageState string

const (
	AssistedSankhyaLineageNone        AssistedSankhyaLineageState = "none"
	AssistedSankhyaLineagePartial     AssistedSankhyaLineageState = "partial"
	AssistedSankhyaLineageComplete    AssistedSankhyaLineageState = "complete"
	AssistedSankhyaLineageConflict    AssistedSankhyaLineageState = "conflict"
	AssistedSankhyaLineageUnavailable AssistedSankhyaLineageState = "unavailable"
)

type AssistedSankhyaDescendant struct {
	Identity         InternalDocumentLineIdentity
	AttendedQuantity *float64
}

type AssistedSankhyaLineage struct {
	MPCLineID   MPCLineID
	Origin      InternalDocumentLineIdentity
	Descendants []AssistedSankhyaDescendant
	State       AssistedSankhyaLineageState
}

type AssistedSankhyaConfirmationResult struct {
	Linkage SankhyaLinkage
	Lines   []AssistedSankhyaLineage
}

// ValidateAssistedSankhyaSelection proves a complete bijection between the
// stable persisted MPC lines and the exact lines of one selected TOP 313
// candidate. It deliberately ignores mutable product, quantity, date, and
// amount evidence as identity proof.
func ValidateAssistedSankhyaSelection(order MarketplaceOrder, candidate AssistedSankhyaCandidate, selections []AssistedSankhyaLineSelection) ([]SankhyaLineMapping, error) {
	invalid := func(field string) error {
		return fmt.Errorf("%w: %s", ErrInvalidAssistedSankhyaLinkage, field)
	}
	if candidate.Header.DocumentID <= 0 || candidate.OperationCode != 313 || len(candidate.Lines) == 0 {
		return nil, invalid("selected_candidate")
	}
	if len(order.Items) == 0 || len(order.Items) != len(candidate.Lines) || len(selections) != len(order.Items) {
		return nil, invalid("line_bijection")
	}

	orderLines := make(map[MPCLineID]struct{}, len(order.Items))
	for _, item := range order.Items {
		if !item.MPCLineID.Valid() || item.ReconciliationState != LineReconciliationStable {
			return nil, invalid("stable_mpc_line")
		}
		if _, duplicate := orderLines[item.MPCLineID]; duplicate {
			return nil, invalid("duplicate_order_mpc_line")
		}
		orderLines[item.MPCLineID] = struct{}{}
	}

	candidateLines := make(map[InternalDocumentLineIdentity]struct{}, len(candidate.Lines))
	for _, line := range candidate.Lines {
		if line.Identity.DocumentID != candidate.Header.DocumentID || line.Identity.LineNumber <= 0 {
			return nil, invalid("candidate_line")
		}
		if _, duplicate := candidateLines[line.Identity]; duplicate {
			return nil, invalid("duplicate_candidate_line")
		}
		candidateLines[line.Identity] = struct{}{}
	}

	seenMPC := make(map[MPCLineID]struct{}, len(selections))
	seenCandidate := make(map[InternalDocumentLineIdentity]struct{}, len(selections))
	mappings := make([]SankhyaLineMapping, 0, len(selections))
	for _, selection := range selections {
		if _, exists := orderLines[selection.MPCLineID]; !exists {
			return nil, invalid("unknown_mpc_line")
		}
		if _, duplicate := seenMPC[selection.MPCLineID]; duplicate {
			return nil, invalid("duplicate_mpc_line")
		}
		if _, exists := candidateLines[selection.CandidateLine]; !exists {
			return nil, invalid("unknown_candidate_line")
		}
		if _, duplicate := seenCandidate[selection.CandidateLine]; duplicate {
			return nil, invalid("duplicate_selected_candidate_line")
		}
		seenMPC[selection.MPCLineID] = struct{}{}
		seenCandidate[selection.CandidateLine] = struct{}{}
		mappings = append(mappings, SankhyaLineMapping{MPCLineID: selection.MPCLineID, Origin: selection.CandidateLine})
	}
	if len(seenMPC) != len(orderLines) || len(seenCandidate) != len(candidateLines) {
		return nil, invalid("line_bijection")
	}
	return mappings, nil
}
