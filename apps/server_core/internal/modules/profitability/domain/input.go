package domain

import "time"

type InputScope string
type InputKind string
type InputQuality string
type ManualAdjustmentCategory string

const (
	InputScopeOrder InputScope = "order"
	InputScopeItem  InputScope = "item"

	InputKindRevenue    InputKind = "revenue"
	InputKindSaleFee    InputKind = "sale_fee"
	InputKindCost       InputKind = "cost"
	InputKindTaxICMS    InputKind = "tax_icms"
	InputKindTaxIPI     InputKind = "tax_ipi"
	InputKindTaxPIS     InputKind = "tax_pis"
	InputKindTaxCOFINS  InputKind = "tax_cofins"
	InputKindFreight    InputKind = "freight"
	InputKindCommission InputKind = "commission"

	InputQualityComplete       InputQuality = "complete"
	InputQualityMissing        InputQuality = "missing"
	InputQualityUnresolvedLink InputQuality = "unresolved_link"
	InputQualityRejectedLink   InputQuality = "rejected_link"
	InputQualityConflictLink   InputQuality = "conflict_link"
	InputQualityStale          InputQuality = "stale"
	InputQualityManual         InputQuality = "manual"
	InputQualityPartial        InputQuality = "partial"

	ManualAdjustmentFreight    ManualAdjustmentCategory = "freight"
	ManualAdjustmentCommission ManualAdjustmentCategory = "commission"
	ManualAdjustmentCost       ManualAdjustmentCategory = "cost"
	ManualAdjustmentGeneric    ManualAdjustmentCategory = "generic_adjustment"
)

type MarginInput struct {
	InstallationID      string       `json:"installation_id"`
	ProviderOrderID     string       `json:"provider_order_id"`
	ProviderItemID      string       `json:"provider_item_id,omitempty"`
	ProviderVariationID string       `json:"provider_variation_id,omitempty"`
	Scope               InputScope   `json:"scope"`
	Kind                InputKind    `json:"kind"`
	Amount              *float64     `json:"amount,omitempty"`
	Currency            string       `json:"currency"`
	SourceSystem        string       `json:"source_system"`
	SourceReference     string       `json:"source_reference,omitempty"`
	ObservedAt          *time.Time   `json:"observed_at,omitempty"`
	Quality             InputQuality `json:"quality"`
	QualityReason       string       `json:"quality_reason,omitempty"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
}

type ManualAdjustment struct {
	AdjustmentID        string                   `json:"adjustment_id"`
	IdempotencyKey      string                   `json:"idempotency_key"`
	InstallationID      string                   `json:"installation_id"`
	ProviderOrderID     string                   `json:"provider_order_id"`
	ProviderItemID      string                   `json:"provider_item_id,omitempty"`
	ProviderVariationID string                   `json:"provider_variation_id,omitempty"`
	Scope               InputScope               `json:"scope"`
	Category            ManualAdjustmentCategory `json:"category"`
	Amount              float64                  `json:"amount"`
	Currency            string                   `json:"currency"`
	Reason              string                   `json:"reason"`
	Actor               ActorMetadata            `json:"actor"`
	CreatedAt           time.Time                `json:"created_at"`
}

type ActorMetadata struct {
	ActorType string `json:"actor_type"`
	ActorID   string `json:"actor_id"`
	ActorName string `json:"actor_name,omitempty"`
}

type ImportInputsResult struct {
	InstallationID string        `json:"installation_id"`
	ImportedCount  int           `json:"imported_count"`
	Items          []MarginInput `json:"items"`
}

type ProfitSnapshotQuality string
type ProfitSnapshotFlag string

const (
	ProfitSnapshotComplete       ProfitSnapshotQuality = "complete"
	ProfitSnapshotIncomplete     ProfitSnapshotQuality = "incomplete"
	ProfitSnapshotNegativeMargin ProfitSnapshotQuality = "negative_margin"
	ProfitSnapshotNotRealized    ProfitSnapshotQuality = "not_realized"

	ProfitFlagMissingRevenue    ProfitSnapshotFlag = "missing_revenue"
	ProfitFlagMissingSaleFee    ProfitSnapshotFlag = "missing_sale_fee"
	ProfitFlagMissingCost       ProfitSnapshotFlag = "missing_cost"
	ProfitFlagMissingTax        ProfitSnapshotFlag = "missing_tax"
	ProfitFlagMissingFreight    ProfitSnapshotFlag = "missing_freight"
	ProfitFlagMissingCommission ProfitSnapshotFlag = "missing_commission"
	ProfitFlagMissingLink       ProfitSnapshotFlag = "missing_link"
	ProfitFlagNegativeMargin    ProfitSnapshotFlag = "negative_margin"
	ProfitFlagOrderCancelled    ProfitSnapshotFlag = "order_cancelled"
	ProfitFlagOrderStateUnknown ProfitSnapshotFlag = "order_state_unknown"
)

type ProfitSnapshot struct {
	InstallationID      string                `json:"installation_id"`
	ProviderOrderID     string                `json:"provider_order_id"`
	ProviderItemID      string                `json:"provider_item_id,omitempty"`
	ProviderVariationID string                `json:"provider_variation_id,omitempty"`
	Scope               InputScope            `json:"scope"`
	Currency            string                `json:"currency"`
	RealizationState    OrderRealizationState `json:"realization_state"`
	RevenueAmount       *float64              `json:"revenue_amount,omitempty"`
	SaleFeeAmount       *float64              `json:"sale_fee_amount,omitempty"`
	CostAmount          *float64              `json:"cost_amount,omitempty"`
	TaxAmount           *float64              `json:"tax_amount,omitempty"`
	FreightAmount       *float64              `json:"freight_amount,omitempty"`
	CommissionAmount    *float64              `json:"commission_amount,omitempty"`
	AdjustmentAmount    *float64              `json:"adjustment_amount,omitempty"`
	ContributionAmount  *float64              `json:"contribution_amount,omitempty"`
	MarginPercent       *float64              `json:"margin_percent,omitempty"`
	Quality             ProfitSnapshotQuality `json:"quality"`
	Flags               []ProfitSnapshotFlag  `json:"flags"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
}

type CalculateSnapshotsResult struct {
	InstallationID  string           `json:"installation_id"`
	CalculatedCount int              `json:"calculated_count"`
	Items           []ProfitSnapshot `json:"items"`
}
