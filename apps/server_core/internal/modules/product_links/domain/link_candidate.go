package domain

import "time"

type LinkCandidateState string
type LinkCandidateMatchInput string
type LinkCandidateConfidenceBand string
type LinkCandidateMatchStatus string
type LinkCandidateReasonDirection string

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

// IC-01 Amendment A2 — deterministic confidence/band/status vocabulary. Enum
// strings are byte-exact across migration values, Go consts, OpenAPI enum,
// SDK type, and UI; do not translate or reshape them.
const (
	LinkCandidateConfidenceBandAlta  LinkCandidateConfidenceBand = "ALTA"
	LinkCandidateConfidenceBandMedia LinkCandidateConfidenceBand = "MEDIA"
	LinkCandidateConfidenceBandBaixa LinkCandidateConfidenceBand = "BAIXA"
)

const (
	LinkCandidateMatchStatusAccept      LinkCandidateMatchStatus = "ACCEPT"
	LinkCandidateMatchStatusReview      LinkCandidateMatchStatus = "REVIEW"
	LinkCandidateMatchStatusReject      LinkCandidateMatchStatus = "REJECT"
	LinkCandidateMatchStatusNoCandidate LinkCandidateMatchStatus = "NO_CANDIDATE"
)

const (
	LinkCandidateReasonDirectionFor         LinkCandidateReasonDirection = "FOR"
	LinkCandidateReasonDirectionAgainst     LinkCandidateReasonDirection = "AGAINST"
	LinkCandidateReasonDirectionUnavailable LinkCandidateReasonDirection = "UNAVAILABLE"
)

// LinkCandidateReason records one anchor's contribution (or unavailability)
// to a candidate's match_status, so the decision's evidence is always
// visible (ADR-17) rather than silently dropped.
type LinkCandidateReason struct {
	Anchor    string                       `json:"anchor"`
	Direction LinkCandidateReasonDirection `json:"direction"`
	Detail    string                       `json:"detail,omitempty"`
}

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
	State                   LinkCandidateState          `json:"state"`
	MatchInput              LinkCandidateMatchInput     `json:"match_input"`
	MatchValue              string                      `json:"match_value,omitempty"`
	SourceSnapshotFetchedAt *time.Time                  `json:"source_snapshot_fetched_at,omitempty"`
	Confidence              int                         `json:"confidence"`
	ConfidenceBand          LinkCandidateConfidenceBand `json:"confidence_band"`
	MatchStatus             LinkCandidateMatchStatus    `json:"match_status"`
	Reasons                 []LinkCandidateReason       `json:"reasons"`
	CreatedAt               time.Time                   `json:"created_at"`
	UpdatedAt               time.Time                   `json:"updated_at"`
}

type LinkCandidateGenerationResult struct {
	InstallationID string          `json:"installation_id"`
	GeneratedCount int             `json:"generated_count"`
	Items          []LinkCandidate `json:"items"`
}
