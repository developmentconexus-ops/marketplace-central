package domain

import "time"

type LinkCandidateState string
type LinkCandidateMatchInput string

const (
	LinkCandidateStateManual     LinkCandidateState = "manual"
	LinkCandidateStateExactSKU   LinkCandidateState = "exact_sku"
	LinkCandidateStateExactEAN   LinkCandidateState = "exact_ean"
	LinkCandidateStateTitleMatch LinkCandidateState = "title_match"
	LinkCandidateStateUnresolved LinkCandidateState = "unresolved"
	LinkCandidateStateConflict   LinkCandidateState = "conflict"
)

const (
	LinkCandidateMatchInputManual    LinkCandidateMatchInput = "manual"
	LinkCandidateMatchInputSellerSKU LinkCandidateMatchInput = "seller_sku"
	LinkCandidateMatchInputEAN       LinkCandidateMatchInput = "ean"
	LinkCandidateMatchInputTitle     LinkCandidateMatchInput = "title"
	LinkCandidateMatchInputNone      LinkCandidateMatchInput = "none"
)

type ListingIdentity struct {
	InstallationID      string `json:"installation_id"`
	ProviderItemID      string `json:"provider_item_id"`
	ProviderVariationID string `json:"provider_variation_id,omitempty"`
}

type LinkCandidate struct {
	CandidateID             string                  `json:"candidate_id"`
	InstallationID          string                  `json:"installation_id"`
	ProviderCode            string                  `json:"provider_code"`
	ProviderItemID          string                  `json:"provider_item_id"`
	ProviderVariationID     string                  `json:"provider_variation_id,omitempty"`
	InternalProductID       *int                    `json:"internal_product_id,omitempty"`
	InternalProductName     string                  `json:"internal_product_name,omitempty"`
	InternalReferenceCode   string                  `json:"internal_reference_code,omitempty"`
	State                   LinkCandidateState      `json:"state"`
	MatchInput              LinkCandidateMatchInput `json:"match_input"`
	MatchValue              string                  `json:"match_value,omitempty"`
	SourceSnapshotFetchedAt *time.Time              `json:"source_snapshot_fetched_at,omitempty"`
	CreatedAt               time.Time               `json:"created_at"`
	UpdatedAt               time.Time               `json:"updated_at"`
}

type LinkCandidateGenerationResult struct {
	InstallationID string          `json:"installation_id"`
	GeneratedCount int             `json:"generated_count"`
	Items          []LinkCandidate `json:"items"`
}
