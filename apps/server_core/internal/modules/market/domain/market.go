package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Money is a known monetary value. Unknown money is represented by a nil
// *Money, never by a zero-valued Money.
type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// CatalogStats contains per-seller-deduplicated catalog price statistics.
// Values are decimal strings because catalog statistics share the BRL money
// contract without repeating a currency on every percentile.
type CatalogStats struct {
	Min         string `json:"min"`
	P25         string `json:"p25"`
	Median      string `json:"median"`
	TrimmedMean string `json:"trimmed_mean"`
	P75         string `json:"p75"`
	Max         string `json:"max"`
	NOffers     int    `json:"n_offers"`
	NSellers    int    `json:"n_sellers"`
}

type CompetitiveStatus string

const (
	CompetitiveStatusWinning   CompetitiveStatus = "winning"
	CompetitiveStatusCompeting CompetitiveStatus = "competing"
	CompetitiveStatusNotListed CompetitiveStatus = "not_listed"
)

type EvidenceState string

const (
	EvidenceStateObserved           EvidenceState = "observed"
	EvidenceStateInsufficientMarket EvidenceState = "insufficient_market"
	EvidenceStateNoPriceEvidence    EvidenceState = "no_price_evidence"
)

type Source string

const (
	SourceOfficialAPI Source = "official_api"
	SourceVendor      Source = "vendor"
	SourceManual      Source = "manual"
)

type MatchState string

const (
	MatchStateAccept      MatchState = "accept"
	MatchStateReview      MatchState = "review"
	MatchStateReject      MatchState = "reject"
	MatchStateNoCandidate MatchState = "no_candidate"
)

// MarketObservation is the independent market evidence captured for one own
// listing. CompetitiveTarget is a target, never a promise that the listing
// will win. The signal fields deliberately remain separate; Source: manual is
// the manual signal's provenance and does not create an aggregate price.
type MarketObservation struct {
	ListingID         string             `json:"listing_id"`
	OurSalePrice      *Money             `json:"our_sale_price"`
	CompetitiveStatus *CompetitiveStatus `json:"competitive_status"`
	WinnerPrice       *Money             `json:"winner_price"`
	CompetitiveTarget *Money             `json:"competitive_target"`
	CatalogOfferPrice *Money             `json:"catalog_offer_price"`
	CatalogStats      *CatalogStats      `json:"catalog_stats"`
	EvidenceState     EvidenceState      `json:"evidence_state"`
	CapturedAt        *time.Time         `json:"captured_at"`
	Source            *Source            `json:"source"`
}

// MarketReference is the independent catalog evidence captured for one ERP
// product. It carries no derived market price.
type MarketReference struct {
	ProductID        string        `json:"product_id"`
	CatalogProductID *string       `json:"catalog_product_id"`
	MatchState       MatchState    `json:"match_state"`
	MatchMethod      *string       `json:"match_method"`
	CatalogStats     *CatalogStats `json:"catalog_stats"`
	EvidenceState    EvidenceState `json:"evidence_state"`
	CapturedAt       *time.Time    `json:"captured_at"`
	Source           *Source       `json:"source"`
}

// MarketObservationInput is the constructor input for a stored observation.
type MarketObservationInput struct {
	ListingID         string
	OurSalePrice      *Money
	CompetitiveStatus *CompetitiveStatus
	WinnerPrice       *Money
	CompetitiveTarget *Money
	CatalogOfferPrice *Money
	CatalogStats      *CatalogStats
	EvidenceState     EvidenceState
	CapturedAt        *time.Time
	Source            *Source
}

// MarketReferenceInput is the constructor input for a stored reference.
type MarketReferenceInput struct {
	ProductID        string
	CatalogProductID *string
	MatchState       MatchState
	MatchMethod      *string
	CatalogStats     *CatalogStats
	EvidenceState    EvidenceState
	CapturedAt       *time.Time
	Source           *Source
}

var (
	ErrListingIDRequired        = errors.New("listing_id is required")
	ErrProductIDRequired        = errors.New("product_id is required")
	ErrSourceRequired           = errors.New("source is required")
	ErrCapturedAtRequired       = errors.New("captured_at is required")
	ErrEvidenceStateInvalid     = errors.New("evidence_state is invalid")
	ErrSourceInvalid            = errors.New("source is invalid")
	ErrCompetitiveStatusInvalid = errors.New("competitive_status is invalid")
	ErrMatchStateInvalid        = errors.New("match_state is invalid")
	ErrCatalogStatsSellerCount  = errors.New("catalog_stats requires at least 5 sellers")
)

// NewStoredObservation constructs an observation that is safe for append-only
// storage. Synthetic no-price read items should be built directly as a
// MarketObservation so their provenance can remain nil.
func NewStoredObservation(input MarketObservationInput) (MarketObservation, error) {
	observation := MarketObservation{
		ListingID:         input.ListingID,
		OurSalePrice:      input.OurSalePrice,
		CompetitiveStatus: input.CompetitiveStatus,
		WinnerPrice:       input.WinnerPrice,
		CompetitiveTarget: input.CompetitiveTarget,
		CatalogOfferPrice: input.CatalogOfferPrice,
		CatalogStats:      input.CatalogStats,
		EvidenceState:     input.EvidenceState,
		CapturedAt:        input.CapturedAt,
		Source:            input.Source,
	}
	if err := observation.ValidateStored(); err != nil {
		return MarketObservation{}, err
	}
	return observation, nil
}

// NewStoredReference constructs a reference that is safe for append-only
// storage.
func NewStoredReference(input MarketReferenceInput) (MarketReference, error) {
	reference := MarketReference{
		ProductID:        input.ProductID,
		CatalogProductID: input.CatalogProductID,
		MatchState:       input.MatchState,
		MatchMethod:      input.MatchMethod,
		CatalogStats:     input.CatalogStats,
		EvidenceState:    input.EvidenceState,
		CapturedAt:       input.CapturedAt,
		Source:           input.Source,
	}
	if err := reference.ValidateStored(); err != nil {
		return MarketReference{}, err
	}
	return reference, nil
}

// ValidateStored enforces the provenance and value invariants for a persisted
// observation.
func (o MarketObservation) ValidateStored() error {
	if strings.TrimSpace(o.ListingID) == "" {
		return ErrListingIDRequired
	}
	if o.CapturedAt == nil {
		return ErrCapturedAtRequired
	}
	if o.Source == nil {
		return ErrSourceRequired
	}
	if !validEvidenceState(o.EvidenceState) {
		return ErrEvidenceStateInvalid
	}
	if !validSource(*o.Source) {
		return ErrSourceInvalid
	}
	if o.CompetitiveStatus != nil && !validCompetitiveStatus(*o.CompetitiveStatus) {
		return ErrCompetitiveStatusInvalid
	}
	for name, money := range map[string]*Money{
		"our_sale_price":      o.OurSalePrice,
		"winner_price":        o.WinnerPrice,
		"competitive_target":  o.CompetitiveTarget,
		"catalog_offer_price": o.CatalogOfferPrice,
	} {
		if err := validateMoney(name, money); err != nil {
			return err
		}
	}
	return validateCatalogStats(o.CatalogStats)
}

// ValidateStored enforces the provenance and value invariants for a persisted
// reference.
func (r MarketReference) ValidateStored() error {
	if strings.TrimSpace(r.ProductID) == "" {
		return ErrProductIDRequired
	}
	if r.CapturedAt == nil {
		return ErrCapturedAtRequired
	}
	if r.Source == nil {
		return ErrSourceRequired
	}
	if !validEvidenceState(r.EvidenceState) {
		return ErrEvidenceStateInvalid
	}
	if !validSource(*r.Source) {
		return ErrSourceInvalid
	}
	if !validMatchState(r.MatchState) {
		return ErrMatchStateInvalid
	}
	return validateCatalogStats(r.CatalogStats)
}

func validCompetitiveStatus(status CompetitiveStatus) bool {
	switch status {
	case CompetitiveStatusWinning, CompetitiveStatusCompeting, CompetitiveStatusNotListed:
		return true
	default:
		return false
	}
}

func validEvidenceState(state EvidenceState) bool {
	switch state {
	case EvidenceStateObserved, EvidenceStateInsufficientMarket, EvidenceStateNoPriceEvidence:
		return true
	default:
		return false
	}
}

func validSource(source Source) bool {
	switch source {
	case SourceOfficialAPI, SourceVendor, SourceManual:
		return true
	default:
		return false
	}
}

func validMatchState(state MatchState) bool {
	switch state {
	case MatchStateAccept, MatchStateReview, MatchStateReject, MatchStateNoCandidate:
		return true
	default:
		return false
	}
}

var decimalPattern = regexp.MustCompile(`^[+-]?(?:\d+(?:\.\d+)?|\.\d+)$`)

func validateMoney(field string, money *Money) error {
	if money == nil {
		return nil
	}
	if !decimalPattern.MatchString(money.Amount) {
		return fmt.Errorf("%s.amount must be a non-empty decimal string", field)
	}
	if money.Currency != "BRL" {
		return fmt.Errorf("%s.currency must be BRL", field)
	}
	return nil
}

func validateCatalogStats(stats *CatalogStats) error {
	if stats == nil {
		return nil
	}
	if stats.NSellers < 5 {
		return ErrCatalogStatsSellerCount
	}
	if stats.NOffers < 0 {
		return errors.New("catalog_stats.n_offers must not be negative")
	}
	for name, value := range map[string]string{
		"min": stats.Min, "p25": stats.P25, "median": stats.Median,
		"trimmed_mean": stats.TrimmedMean, "p75": stats.P75, "max": stats.Max,
	} {
		if value != "" && !decimalPattern.MatchString(value) {
			return fmt.Errorf("catalog_stats.%s must be a decimal string", name)
		}
	}
	return nil
}
