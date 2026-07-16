package productlinks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	mutationsdomain "marketplace-central/apps/server_core/internal/modules/mutations/domain"
	mutationsports "marketplace-central/apps/server_core/internal/modules/mutations/ports"
	productlinksapp "marketplace-central/apps/server_core/internal/modules/product_links/application"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type resolutionAPI interface {
	ApproveCandidate(context.Context, productlinksapp.ApproveCandidateInput) (productlinksdomain.ProductLinkResolutionResult, error)
	RejectListing(context.Context, productlinksapp.RejectListingInput) (productlinksdomain.ProductLinkResolutionResult, error)
	ManualResolve(context.Context, productlinksapp.ManualResolveInput) (productlinksdomain.ProductLinkResolutionResult, error)
}

type workflowReader interface {
	ListLinkWorkflows(context.Context, productlinksapp.ListLinkWorkflowsInput) ([]productlinksdomain.ProductLinkWorkflowItem, error)
}

type Writer struct {
	resolver resolutionAPI
}

func NewWriter(resolver resolutionAPI) *Writer {
	return &Writer{resolver: resolver}
}

var _ mutationsports.LinkageWriter = (*Writer)(nil)

func (w *Writer) ApplyLink(ctx context.Context, input mutationsports.LinkageInput) (mutationsports.LinkageOutcome, error) {
	if w == nil || w.resolver == nil {
		return mutationsports.LinkageOutcome{}, errors.New("product link writer is not configured")
	}
	if strings.TrimSpace(input.InstallationID) == "" {
		return mutationsports.LinkageOutcome{}, errors.New("product link installation is required")
	}
	preflight, err := w.preflight(ctx, input)
	if err != nil {
		return mutationsports.LinkageOutcome{}, err
	}
	if preflight != nil {
		return *preflight, nil
	}

	var result productlinksdomain.ProductLinkResolutionResult
	actor := productlinksdomain.ActorMetadata{ActorType: input.ActorType, ActorID: input.ActorID, ActorName: input.ActorName}
	switch input.Action {
	case mutationsports.LinkageActionApproveCandidate:
		if strings.TrimSpace(input.CandidateID) == "" {
			return mutationsports.LinkageOutcome{}, errors.New("product link candidate is required")
		}
		result, err = w.resolver.ApproveCandidate(ctx, productlinksapp.ApproveCandidateInput{
			CandidateID: input.CandidateID,
			Actor:       actor,
			Reason:      input.Reason,
		})
	case mutationsports.LinkageActionManualResolve:
		result, err = w.resolver.ManualResolve(ctx, productlinksapp.ManualResolveInput{
			InstallationID:        input.InstallationID,
			ProviderCode:          input.ProviderCode,
			ProviderItemID:        input.ProviderItemID,
			ProviderVariationID:   input.ProviderVariationID,
			InternalProductID:     input.InternalProductID,
			InternalProductName:   input.InternalProductName,
			InternalReferenceCode: input.InternalReferenceCode,
			Actor:                 actor,
			Reason:                input.Reason,
		})
	case mutationsports.LinkageActionRejectListing:
		result, err = w.resolver.RejectListing(ctx, productlinksapp.RejectListingInput{
			InstallationID:      input.InstallationID,
			ProviderCode:        input.ProviderCode,
			ProviderItemID:      input.ProviderItemID,
			ProviderVariationID: input.ProviderVariationID,
			Actor:               actor,
			Reason:              input.Reason,
		})
	default:
		return mutationsports.LinkageOutcome{}, fmt.Errorf("unsupported product link action %q", input.Action)
	}
	if err != nil {
		return mutationsports.LinkageOutcome{}, err
	}

	if identicalResolution(result.Audit) {
		return mutationsports.LinkageOutcome{State: mutationsdomain.ItemStateSkipped}, nil
	}
	if conflictingResolution(result.Audit) {
		failure := mutationsdomain.Failure{
			Code:      mutationsdomain.FailureCodeConflictRemoteChanged,
			MessagePT: "O vínculo remoto mudou desde a prévia.",
			Retryable: false,
		}
		return mutationsports.LinkageOutcome{Failure: &failure}, nil
	}
	return mutationsports.LinkageOutcome{State: mutationsdomain.ItemStateApplied}, nil
}

// Apply keeps the generic mutation writer shape for the poller. F02-S6 can use
// ApplyLink when it needs the link-specific skipped terminal state.
func (w *Writer) Apply(ctx context.Context, item mutationsports.WriteItem) (mutationsports.WriteOutcome, error) {
	var payload struct {
		Action                mutationsports.LinkageAction `json:"action"`
		CandidateID           string                       `json:"candidate_id"`
		ProviderCode          string                       `json:"provider_code"`
		ProviderItemID        string                       `json:"provider_item_id"`
		ProviderVariationID   string                       `json:"provider_variation_id"`
		InternalProductID     int                          `json:"product_id"`
		InternalProductName   string                       `json:"product_name"`
		InternalReferenceCode string                       `json:"reference_code"`
		ActorType             string                       `json:"actor_type"`
		ActorID               string                       `json:"actor_id"`
		ActorName             string                       `json:"actor_name"`
		Reason                string                       `json:"reason"`
	}
	if err := json.Unmarshal(item.After, &payload); err != nil {
		return mutationsports.WriteOutcome{}, fmt.Errorf("decode product link intent: %w", err)
	}
	if payload.ProviderItemID == "" {
		payload.ProviderItemID = item.ListingID
	}
	outcome, err := w.ApplyLink(ctx, mutationsports.LinkageInput{
		InstallationID:        item.InstallationID,
		Action:                payload.Action,
		CandidateID:           payload.CandidateID,
		ProviderCode:          payload.ProviderCode,
		ProviderItemID:        payload.ProviderItemID,
		ProviderVariationID:   payload.ProviderVariationID,
		InternalProductID:     payload.InternalProductID,
		InternalProductName:   payload.InternalProductName,
		InternalReferenceCode: payload.InternalReferenceCode,
		ActorType:             payload.ActorType,
		ActorID:               payload.ActorID,
		ActorName:             payload.ActorName,
		Reason:                payload.Reason,
	})
	if err != nil {
		return mutationsports.WriteOutcome{}, err
	}
	return mutationsports.WriteOutcome{Failure: outcome.Failure}, nil
}

func (w *Writer) HasResolvedLink(ctx context.Context, installationID string) (bool, error) {
	reader, ok := w.resolver.(workflowReader)
	if !ok {
		return false, errors.New("product link workflow reader is not configured")
	}
	items, err := reader.ListLinkWorkflows(ctx, productlinksapp.ListLinkWorkflowsInput{InstallationID: installationID, Limit: 2000})
	if err != nil {
		return false, err
	}
	for _, item := range items {
		if item.CurrentLink != nil && item.CurrentLink.State == productlinksdomain.ProductLinkStateResolved {
			return true, nil
		}
	}
	return false, nil
}

func (w *Writer) preflight(ctx context.Context, input mutationsports.LinkageInput) (*mutationsports.LinkageOutcome, error) {
	reader, ok := w.resolver.(workflowReader)
	if !ok {
		return nil, nil
	}
	items, err := reader.ListLinkWorkflows(ctx, productlinksapp.ListLinkWorkflowsInput{InstallationID: input.InstallationID, Limit: 2000})
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if !workflowMatches(item, input) {
			continue
		}
		if item.CurrentLink == nil {
			return nil, nil
		}
		targetState := productlinksdomain.ProductLinkStateResolved
		var targetProductID *int
		if input.Action == mutationsports.LinkageActionRejectListing {
			targetState = productlinksdomain.ProductLinkStateRejected
		} else if input.Action == mutationsports.LinkageActionApproveCandidate {
			for _, candidate := range item.Candidates {
				if candidate.CandidateID == input.CandidateID {
					targetProductID = candidate.InternalProductID
					break
				}
			}
			if targetProductID == nil {
				return nil, nil
			}
		} else if input.InternalProductID > 0 {
			targetProductID = &input.InternalProductID
		} else {
			return nil, nil
		}
		if item.CurrentLink.State != targetState {
			return nil, nil
		}
		if targetState == productlinksdomain.ProductLinkStateResolved && !sameProduct(item.CurrentLink.InternalProductID, targetProductID) {
			failure := mutationsdomain.Failure{Code: mutationsdomain.FailureCodeConflictRemoteChanged, MessagePT: "O vínculo remoto mudou desde a prévia."}
			return &mutationsports.LinkageOutcome{Failure: &failure}, nil
		}
		return &mutationsports.LinkageOutcome{State: mutationsdomain.ItemStateSkipped}, nil
	}
	return nil, nil
}

func workflowMatches(item productlinksdomain.ProductLinkWorkflowItem, input mutationsports.LinkageInput) bool {
	if input.Action == mutationsports.LinkageActionApproveCandidate {
		for _, candidate := range item.Candidates {
			if candidate.CandidateID == input.CandidateID {
				return true
			}
		}
		return false
	}
	identity := item.Identity
	return identity.InstallationID == input.InstallationID &&
		identity.ProviderItemID == input.ProviderItemID &&
		identity.ProviderVariationID == input.ProviderVariationID
}

func identicalResolution(audit productlinksdomain.ProductLinkAuditEntry) bool {
	if audit.PreviousState != audit.NextState {
		return false
	}
	switch audit.NextState {
	case productlinksdomain.ProductLinkStateResolved:
		return sameProduct(audit.PreviousInternalProductID, audit.NextInternalProductID)
	case productlinksdomain.ProductLinkStateRejected:
		return audit.PreviousInternalProductID == nil && audit.NextInternalProductID == nil
	default:
		return false
	}
}

func conflictingResolution(audit productlinksdomain.ProductLinkAuditEntry) bool {
	return audit.PreviousState == productlinksdomain.ProductLinkStateResolved &&
		audit.NextState == productlinksdomain.ProductLinkStateResolved &&
		!sameProduct(audit.PreviousInternalProductID, audit.NextInternalProductID)
}

func sameProduct(left, right *int) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}
