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
	// ProductLinkActionUndo marks a reversal transition (S4): the audit row
	// it produces is a NEW append-only entry — never a delete/rewrite of the
	// entry it reverses — so the full resolution history stays intact.
	ProductLinkActionUndo ProductLinkAction = "undo"
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
	// BatchID links this per-item audit row to the product_link_batches
	// aggregate row it was written under (S3). Empty for non-batch
	// resolutions (approve/reject/manual outside a batch).
	BatchID                string            `json:"batch_id,omitempty"`
	CreatedAt              time.Time         `json:"created_at"`
}

type ProductLinkTransition struct {
	Link  ProductLink
	Audit ProductLinkAuditEntry
	// Decision is the E10 row this transition also writes, in the SAME
	// transaction as Link and Audit. A rejection and an undo carry one too:
	// they name no anchor, but they DO settle the link, and omitting them
	// would leave the decision they overrule reading as the one in force.
	// nil only where the caller stated no rule at all.
	Decision *ProductLinkDecision
}

// ProductLinkBatchStatus is the aggregate outcome of a batch-apply run (S3).
type ProductLinkBatchStatus string

const (
	ProductLinkBatchStatusCompleted ProductLinkBatchStatus = "COMPLETED"
	ProductLinkBatchStatusPartial   ProductLinkBatchStatus = "PARTIAL"
	ProductLinkBatchStatusFailed    ProductLinkBatchStatus = "FAILED"
)

// ProductLinkBatch is the module-owned batch audit row (C02 "tabela
// própria"): one row per ApplyBatch call, aggregating counts across all
// per-item transitions/audit entries tagged with the same BatchID.
type ProductLinkBatch struct {
	BatchID        string                 `json:"batch_id"`
	InstallationID string                 `json:"installation_id"`
	Actor          ActorMetadata          `json:"actor"`
	RequestedCount int                    `json:"requested_count"`
	AppliedCount   int                    `json:"applied_count"`
	FailedCount    int                    `json:"failed_count"`
	Status         ProductLinkBatchStatus `json:"status"`
	CreatedAt      time.Time              `json:"created_at"`
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
