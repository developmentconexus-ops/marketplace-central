package internalread

import (
	"context"
	"errors"
	"strings"

	internaldomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
)

type sankhyaLinkageSource interface {
	ValidateConfiguration(ctx context.Context) error
	FindCandidates(ctx context.Context, input internalports.SankhyaCandidateInput) ([]internaldomain.SankhyaDocumentCandidate, error)
	ListDescendants(ctx context.Context, input internalports.SankhyaDescendantInput) (internaldomain.SankhyaLineage, error)
}

type SankhyaLinkageReader struct {
	source                sankhyaLinkageSource
	configurationRevision string
}

func NewSankhyaLinkageReader(source sankhyaLinkageSource, configurationRevision string) (*SankhyaLinkageReader, error) {
	configurationRevision = strings.TrimSpace(configurationRevision)
	if configurationRevision == "" {
		return nil, &ordersdomain.AssistedSankhyaReadError{Kind: ordersdomain.AssistedSankhyaReadConfigurationInvalid}
	}
	return &SankhyaLinkageReader{source: source, configurationRevision: configurationRevision}, nil
}

func (r *SankhyaLinkageReader) ValidateConfiguration(ctx context.Context) error {
	if r == nil || r.source == nil || r.configurationRevision == "" {
		return &ordersdomain.AssistedSankhyaReadError{Kind: ordersdomain.AssistedSankhyaReadConfigurationInvalid}
	}
	return translateReadError(r.source.ValidateConfiguration(ctx))
}

func (r *SankhyaLinkageReader) ConfigurationRevision() string {
	if r == nil {
		return ""
	}
	return r.configurationRevision
}

func (r *SankhyaLinkageReader) FindCandidates(ctx context.Context, externalOrderKey string) ([]ordersdomain.AssistedSankhyaCandidate, error) {
	if r == nil || r.source == nil {
		return nil, &ordersdomain.AssistedSankhyaReadError{Kind: ordersdomain.AssistedSankhyaReadConfigurationInvalid}
	}
	candidates, err := r.source.FindCandidates(ctx, internalports.SankhyaCandidateInput{ExternalOrderKey: externalOrderKey})
	if err != nil {
		return nil, translateReadError(err)
	}
	result := make([]ordersdomain.AssistedSankhyaCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		converted := ordersdomain.AssistedSankhyaCandidate{
			Header:         ordersdomain.InternalDocumentIdentity{DocumentID: candidate.DocumentID},
			DocumentNumber: candidate.DocumentNumber,
			OperationCode:  candidate.OperationCode,
			NegotiatedAt:   candidate.NegotiatedAt,
			Amount:         candidate.Amount,
		}
		for _, line := range candidate.Lines {
			converted.Lines = append(converted.Lines, ordersdomain.AssistedSankhyaCandidateLine{
				Identity:  ordersdomain.InternalDocumentLineIdentity{DocumentID: line.Identity.DocumentID, LineNumber: line.Identity.LineNumber},
				ProductID: line.ProductID,
				Quantity:  line.Quantity,
				Amount:    line.Amount,
			})
		}
		result = append(result, converted)
	}
	return result, nil
}

func (r *SankhyaLinkageReader) ListDescendants(ctx context.Context, origin ordersdomain.InternalDocumentLineIdentity, expectedOriginQuantity *float64) (ordersdomain.AssistedSankhyaLineage, error) {
	if r == nil || r.source == nil {
		return ordersdomain.AssistedSankhyaLineage{}, &ordersdomain.AssistedSankhyaReadError{Kind: ordersdomain.AssistedSankhyaReadConfigurationInvalid}
	}
	lineage, err := r.source.ListDescendants(ctx, internalports.SankhyaDescendantInput{
		Origin:                 internaldomain.InternalDocumentLine{DocumentID: origin.DocumentID, LineNumber: origin.LineNumber},
		ExpectedOriginQuantity: expectedOriginQuantity,
	})
	if err != nil {
		return ordersdomain.AssistedSankhyaLineage{}, translateReadError(err)
	}
	result := ordersdomain.AssistedSankhyaLineage{
		Origin: ordersdomain.InternalDocumentLineIdentity{DocumentID: lineage.Origin.DocumentID, LineNumber: lineage.Origin.LineNumber},
		State:  convertLineageState(lineage.State),
	}
	for _, descendant := range lineage.Descendants {
		result.Descendants = append(result.Descendants, ordersdomain.AssistedSankhyaDescendant{
			Identity: ordersdomain.InternalDocumentLineIdentity{
				DocumentID: descendant.Identity.DocumentID,
				LineNumber: descendant.Identity.LineNumber,
			},
			AttendedQuantity: descendant.AttendedQuantity,
		})
	}
	return result, nil
}

func translateReadError(err error) error {
	if err == nil {
		return nil
	}
	var readErr *internaldomain.ReadError
	if !errors.As(err, &readErr) {
		return &ordersdomain.AssistedSankhyaReadError{Kind: ordersdomain.AssistedSankhyaReadUnavailable}
	}
	kind := ordersdomain.AssistedSankhyaReadUnavailable
	switch readErr.Code {
	case internaldomain.ReadErrorConfigurationInvalid, internaldomain.ReadErrorMetadataMismatch, internaldomain.ReadErrorUniquenessUnproved:
		kind = ordersdomain.AssistedSankhyaReadConfigurationInvalid
	case internaldomain.ReadErrorCandidateAmbiguous, internaldomain.ReadErrorTOPMismatch:
		kind = ordersdomain.AssistedSankhyaReadCandidateConflict
	case internaldomain.ReadErrorLineageConflict:
		kind = ordersdomain.AssistedSankhyaReadLineageConflict
	}
	return &ordersdomain.AssistedSankhyaReadError{Kind: kind}
}

func convertLineageState(state internaldomain.SankhyaLineageState) ordersdomain.AssistedSankhyaLineageState {
	switch state {
	case internaldomain.SankhyaLineageNone:
		return ordersdomain.AssistedSankhyaLineageNone
	case internaldomain.SankhyaLineagePartial:
		return ordersdomain.AssistedSankhyaLineagePartial
	case internaldomain.SankhyaLineageComplete:
		return ordersdomain.AssistedSankhyaLineageComplete
	case internaldomain.SankhyaLineageConflict:
		return ordersdomain.AssistedSankhyaLineageConflict
	default:
		return ordersdomain.AssistedSankhyaLineageUnavailable
	}
}

var _ ordersports.SankhyaLinkageReader = (*SankhyaLinkageReader)(nil)
