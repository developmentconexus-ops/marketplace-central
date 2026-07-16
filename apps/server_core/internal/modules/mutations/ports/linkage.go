package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
)

type LinkageReader interface {
	HasResolvedLink(context.Context, string) (bool, error)
}

type LinkageAction string

const (
	LinkageActionApproveCandidate LinkageAction = "approve_candidate"
	LinkageActionManualResolve    LinkageAction = "manual_resolve"
	LinkageActionRejectListing    LinkageAction = "reject_listing"
)

type LinkageInput struct {
	InstallationID        string
	Action                LinkageAction
	CandidateID           string
	ProviderCode          string
	ProviderItemID        string
	ProviderVariationID   string
	InternalProductID     int
	InternalProductName   string
	InternalReferenceCode string
	ActorType             string
	ActorID               string
	ActorName             string
	Reason                string
}

type LinkageOutcome struct {
	State   domain.ItemState
	Failure *domain.Failure
}

type LinkageWriter interface {
	ApplyLink(context.Context, LinkageInput) (LinkageOutcome, error)
}
