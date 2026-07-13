package domain

import (
	"errors"
	"math"
	"time"
)

const (
	ReadErrorConfigurationInvalid ReadErrorCode = "configuration_invalid"
	ReadErrorMetadataMismatch     ReadErrorCode = "field_metadata_mismatch"
	ReadErrorUniquenessUnproved   ReadErrorCode = "uniqueness_unproved"
	ReadErrorCandidateAmbiguous   ReadErrorCode = "candidate_ambiguous"
	ReadErrorLineageConflict      ReadErrorCode = "lineage_conflict"
	ReadErrorTOPMismatch          ReadErrorCode = "top_mismatch"
)

type SankhyaLineageState string

const (
	SankhyaLineageNone     SankhyaLineageState = "none"
	SankhyaLineagePartial  SankhyaLineageState = "partial"
	SankhyaLineageComplete SankhyaLineageState = "complete"
	SankhyaLineageConflict SankhyaLineageState = "conflict"
)

type InternalDocumentLine struct {
	DocumentID int64
	LineNumber int
}

func (l InternalDocumentLine) Valid() bool {
	return l.DocumentID > 0 && l.LineNumber > 0
}

type SankhyaLineEvidence struct {
	Identity  InternalDocumentLine
	ProductID *int64
	Quantity  *float64
	Amount    *float64
}

type SankhyaDocumentCandidate struct {
	DocumentID     int64
	DocumentNumber *int64
	OperationCode  int64
	NegotiatedAt   *time.Time
	Amount         *float64
	Lines          []SankhyaLineEvidence
}

type SankhyaDescendant struct {
	Identity         InternalDocumentLine
	AttendedQuantity *float64
}

type SankhyaLineage struct {
	Origin      InternalDocumentLine
	Descendants []SankhyaDescendant
	State       SankhyaLineageState
}

func NewSankhyaLineage(origin InternalDocumentLine, expectedOriginQuantity *float64, descendants []SankhyaDescendant) (SankhyaLineage, error) {
	result := SankhyaLineage{
		Origin:      origin,
		Descendants: append([]SankhyaDescendant(nil), descendants...),
	}
	if !origin.Valid() {
		return result, NewReadError(ReadErrorLineageConflict, "sankhya lineage origin is invalid", nil)
	}
	if expectedOriginQuantity != nil && (*expectedOriginQuantity < 0 || math.IsNaN(*expectedOriginQuantity) || math.IsInf(*expectedOriginQuantity, 0)) {
		return result, NewReadError(ReadErrorLineageConflict, "sankhya expected origin quantity is invalid", nil)
	}
	if len(descendants) == 0 {
		result.State = SankhyaLineageNone
		return result, nil
	}

	seen := make(map[InternalDocumentLine]struct{}, len(descendants))
	allKnown := true
	total := 0.0
	for _, descendant := range descendants {
		if !descendant.Identity.Valid() {
			result.State = SankhyaLineageConflict
			return result, NewReadError(ReadErrorLineageConflict, "sankhya descendant identity is invalid", nil)
		}
		if _, ok := seen[descendant.Identity]; ok {
			result.State = SankhyaLineageConflict
			return result, NewReadError(ReadErrorLineageConflict, "sankhya descendant identity is duplicated", nil)
		}
		seen[descendant.Identity] = struct{}{}
		if descendant.AttendedQuantity == nil {
			allKnown = false
			continue
		}
		quantity := *descendant.AttendedQuantity
		if quantity < 0 || math.IsNaN(quantity) || math.IsInf(quantity, 0) {
			result.State = SankhyaLineageConflict
			return result, NewReadError(ReadErrorLineageConflict, "sankhya attended quantity is invalid", nil)
		}
		total += quantity
	}

	if expectedOriginQuantity != nil && total > *expectedOriginQuantity {
		result.State = SankhyaLineageConflict
		return result, NewReadError(ReadErrorLineageConflict, "sankhya attended quantity exceeds the expected origin quantity", nil)
	}
	if expectedOriginQuantity == nil || !allKnown || total < *expectedOriginQuantity {
		result.State = SankhyaLineagePartial
		return result, nil
	}
	result.State = SankhyaLineageComplete
	return result, nil
}

func IsSankhyaLinkageError(err error, code ReadErrorCode) bool {
	var readErr *ReadError
	return errors.As(err, &readErr) && readErr.Code == code
}
