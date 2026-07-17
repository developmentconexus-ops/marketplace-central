package domain

import "time"

type StockActionState string
type StockActionTrigger string
type StockWriteResultStatus string

const (
	StockActionStateProposed StockActionState = "proposed"
	StockActionStateBlocked  StockActionState = "blocked"
	StockActionStateApproved StockActionState = "approved"
	StockActionStateApplied  StockActionState = "applied"
	StockActionStateFailed   StockActionState = "failed"
	StockActionStateSkipped  StockActionState = "skipped"

	StockActionTriggerManual StockActionTrigger = "manual"

	StockWriteResultApplied          StockWriteResultStatus = "applied"
	StockWriteResultRejected         StockWriteResultStatus = "rejected"
	StockWriteResultTransientFailure StockWriteResultStatus = "transient_failure"
	StockWriteResultUnsupportedShape StockWriteResultStatus = "unsupported_shape"
)

type StockActionOperator struct {
	ActorType string `json:"actor_type"`
	ActorID   string `json:"actor_id"`
	ActorName string `json:"actor_name,omitempty"`
}

type ManualApproval struct {
	Approved   bool                `json:"approved"`
	Operator   StockActionOperator `json:"operator"`
	ApprovedAt time.Time           `json:"approved_at"`
}

type ProviderStockRef struct {
	TenantID            string `json:"tenant_id"`
	InstallationID      string `json:"installation_id"`
	ProviderCode        string `json:"provider_code"`
	ProviderAccountID   string `json:"provider_account_id"`
	ProviderItemID      string `json:"provider_item_id"`
	ProviderVariationID string `json:"provider_variation_id,omitempty"`
}

type StockWriteRequest struct {
	ProviderRef       ProviderStockRef
	IdempotencyKey    string
	RequestedQuantity int
	Reason            string
}

type StockWriteResult struct {
	Status              StockWriteResultStatus `json:"status,omitempty"`
	IdempotencyKey      string                 `json:"idempotency_key,omitempty"`
	ProviderCode        string                 `json:"provider_code,omitempty"`
	ProviderItemID      string                 `json:"provider_item_id,omitempty"`
	ProviderVariationID string                 `json:"provider_variation_id,omitempty"`
	Message             string                 `json:"message,omitempty"`
	ResponseSummary     string                 `json:"response_summary,omitempty"`
}

type StockAction struct {
	ActionID            string                  `json:"action_id"`
	MutationProtocolID  *string                 `json:"mutation_protocol_id,omitempty"`
	State               StockActionState        `json:"state"`
	Trigger             StockActionTrigger      `json:"trigger"`
	ProviderRef         ProviderStockRef        `json:"provider_ref"`
	BeforeQuantity      *int                    `json:"before_quantity,omitempty"`
	RequestedQuantity   int                     `json:"requested_quantity"`
	RecommendedQuantity *int                    `json:"recommended_quantity,omitempty"`
	PolicyID            string                  `json:"policy_id"`
	InternalObservedAt  *time.Time              `json:"internal_observed_at,omitempty"`
	ProviderObservedAt  *time.Time              `json:"provider_observed_at,omitempty"`
	Operator            StockActionOperator     `json:"operator"`
	Reason              string                  `json:"reason,omitempty"`
	IdempotencyKey      string                  `json:"idempotency_key"`
	BlockingReason      BlockingReason          `json:"blocking_reason,omitempty"`
	FailureReason       BlockingReason          `json:"failure_reason,omitempty"`
	ProviderResult      StockWriteResult        `json:"provider_result,omitempty"`
	AuditEvents         []StockActionAuditEvent `json:"audit_events"`
	CreatedAt           time.Time               `json:"created_at"`
	UpdatedAt           time.Time               `json:"updated_at"`
}

type StockActionAuditEvent struct {
	State          StockActionState       `json:"state"`
	Reason         string                 `json:"reason"`
	ProviderStatus StockWriteResultStatus `json:"provider_status,omitempty"`
	OccurredAt     time.Time              `json:"occurred_at"`
}
