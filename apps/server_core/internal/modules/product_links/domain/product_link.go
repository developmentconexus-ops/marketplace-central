package domain

import "time"

type ProductLinkState string
type ProductLinkAction string

const (
	ProductLinkStateNone       ProductLinkState = "none"
	ProductLinkStateUnresolved ProductLinkState = "unresolved"
	ProductLinkStateConflict   ProductLinkState = "conflict"
	ProductLinkStateResolved   ProductLinkState = "resolved"
	ProductLinkStateRejected   ProductLinkState = "rejected"
)

const (
	ProductLinkActionApproveCandidate ProductLinkAction = "approve_candidate"
	ProductLinkActionRejectListing    ProductLinkAction = "reject_listing"
	ProductLinkActionManualResolve    ProductLinkAction = "manual_resolve"
)

type ActorMetadata struct {
	ActorType string `json:"actor_type"`
	ActorID   string `json:"actor_id"`
	ActorName string `json:"actor_name,omitempty"`
}

type ProductLink struct {
	InstallationID        string           `json:"installation_id"`
	ProviderCode          string           `json:"provider_code"`
	ProviderItemID        string           `json:"provider_item_id"`
	ProviderVariationID   string           `json:"provider_variation_id,omitempty"`
	State                 ProductLinkState `json:"state"`
	SourceCandidateID     string           `json:"source_candidate_id,omitempty"`
	InternalProductID     *int             `json:"internal_product_id,omitempty"`
	InternalProductName   string           `json:"internal_product_name,omitempty"`
	InternalReferenceCode string           `json:"internal_reference_code,omitempty"`
	CreatedAt             time.Time        `json:"created_at"`
	UpdatedAt             time.Time        `json:"updated_at"`
}

type ProductLinkAuditEntry struct {
	AuditID                string            `json:"audit_id"`
	InstallationID         string            `json:"installation_id"`
	ProviderCode           string            `json:"provider_code"`
	ProviderItemID         string            `json:"provider_item_id"`
	ProviderVariationID    string            `json:"provider_variation_id,omitempty"`
	Action                 ProductLinkAction `json:"action"`
	Reason                 string            `json:"reason,omitempty"`
	SourceCandidateID      string            `json:"source_candidate_id,omitempty"`
	Actor                  ActorMetadata     `json:"actor"`
	PreviousState          ProductLinkState  `json:"previous_state"`
	NextState              ProductLinkState  `json:"next_state"`
	PreviousInternalProductID *int           `json:"previous_internal_product_id,omitempty"`
	NextInternalProductID  *int              `json:"next_internal_product_id,omitempty"`
	CreatedAt              time.Time         `json:"created_at"`
}

type ProductLinkTransition struct {
	Link  ProductLink
	Audit ProductLinkAuditEntry
}

type ProductLinkWorkflowItem struct {
	Identity    ListingIdentity        `json:"identity"`
	CurrentLink *ProductLink           `json:"current_link,omitempty"`
	Candidates  []LinkCandidate        `json:"candidates"`
	Audit       []ProductLinkAuditEntry `json:"audit"`
}

type ProductLinkResolutionResult struct {
	Link  ProductLink           `json:"link"`
	Audit ProductLinkAuditEntry `json:"audit"`
}
