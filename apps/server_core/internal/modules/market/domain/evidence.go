package domain

import (
	"fmt"
	"math/big"
	"strings"
	"time"
)

type MarketPriceScope string

const (
	MarketPriceScopeListing        MarketPriceScope = "listing"
	MarketPriceScopeCatalogProduct MarketPriceScope = "catalog_product"
)

type MarketPriceSource string

const (
	MarketPriceSourceMLSalePrice     MarketPriceSource = "ml_sale_price"
	MarketPriceSourceMLPriceToWin    MarketPriceSource = "ml_price_to_win"
	MarketPriceSourceMLCatalogOffers MarketPriceSource = "ml_catalog_offers"
)

type MarketPriceSnapshotStatus string

const (
	MarketPriceSnapshotStatusValid   MarketPriceSnapshotStatus = "VALID"
	MarketPriceSnapshotStatusFailed  MarketPriceSnapshotStatus = "FAILED"
	MarketPriceSnapshotStatusExpired MarketPriceSnapshotStatus = "EXPIRED"
)

type MatchStatus string

const (
	MatchStatusAccept      MatchStatus = "ACCEPT"
	MatchStatusReview      MatchStatus = "REVIEW"
	MatchStatusReject      MatchStatus = "REJECT"
	MatchStatusNoCandidate MatchStatus = "NO_CANDIDATE"
)

type PriceEvidenceStatus string

const (
	PriceEvidenceStatusOK                 PriceEvidenceStatus = "OK"
	PriceEvidenceStatusNoPriceEvidence    PriceEvidenceStatus = "NO_PRICE_EVIDENCE"
	PriceEvidenceStatusInsufficientMarket PriceEvidenceStatus = "INSUFFICIENT_MARKET"
)

type BlockingState string

const (
	BlockingStateNoCandidate        BlockingState = "NO_CANDIDATE"
	BlockingStateNoPriceEvidence    BlockingState = "NO_PRICE_EVIDENCE"
	BlockingStateInsufficientMarket BlockingState = "INSUFFICIENT_MARKET"
	BlockingStateSemCusto           BlockingState = "SEM_CUSTO"
)

type ConfidenceBand string

const (
	ConfidenceBandAlta  ConfidenceBand = "ALTA"
	ConfidenceBandMedia ConfidenceBand = "MEDIA"
	ConfidenceBandBaixa ConfidenceBand = "BAIXA"
)

type MarketAggregateStatus string

const (
	MarketAggregateStatusOK                 MarketAggregateStatus = "OK"
	MarketAggregateStatusInsufficientMarket MarketAggregateStatus = "INSUFFICIENT_MARKET"
	MarketAggregateStatusNoPriceEvidence    MarketAggregateStatus = "NO_PRICE_EVIDENCE"
)

type MarketPosition struct {
	Rank  int `json:"rank"`
	Total int `json:"total"`
}

type MarketPriceSnapshot struct {
	Scope         MarketPriceScope          `json:"scope"`
	RefID         string                    `json:"ref_id"`
	Source        MarketPriceSource         `json:"source"`
	Price         *Money                    `json:"price"`
	ObservedAt    time.Time                 `json:"observed_at"`
	FetchedAt     time.Time                 `json:"fetched_at"`
	ExpiresAt     time.Time                 `json:"expires_at"`
	Status        MarketPriceSnapshotStatus `json:"status"`
	FailureReason *string                   `json:"failure_reason,omitempty"`
	RequestID     string                    `json:"request_id"`
}

type MarketPriceSnapshotInput struct {
	Scope         MarketPriceScope
	RefID         string
	Source        MarketPriceSource
	Price         *Money
	ObservedAt    time.Time
	FetchedAt     time.Time
	ExpiresAt     time.Time
	Status        MarketPriceSnapshotStatus
	FailureReason *string
	RequestID     string
}

type ValidatedOffer struct {
	CatalogProductID string      `json:"catalog_product_id"`
	SellerID         string      `json:"seller_id"`
	Price            *Money      `json:"price"`
	Condition        string      `json:"condition"`
	MatchStatus      MatchStatus `json:"match_status"`
	FetchedAt        time.Time   `json:"fetched_at"`
}

type ValidatedOfferInput struct {
	CatalogProductID string
	SellerID         string
	Price            *Money
	Condition        string
	MatchStatus      MatchStatus
	FetchedAt        time.Time
}

type MatchDecision struct {
	ProductID        string         `json:"product_id"`
	CatalogProductID *string        `json:"catalog_product_id,omitempty"`
	Status           MatchStatus    `json:"match_status"`
	Anchors          []string       `json:"anchors"`
	Contradictions   []string       `json:"contradictions"`
	ConfidenceBand   ConfidenceBand `json:"confidence_band"`
	DecidedAt        time.Time      `json:"decided_at"`
}

type MatchDecisionInput struct {
	ProductID        string
	CatalogProductID *string
	Status           MatchStatus
	Anchors          []string
	Contradictions   []string
	ConfidenceBand   ConfidenceBand
	DecidedAt        time.Time
}

type CompetitiveSignal struct {
	ListingID   string          `json:"listing_id"`
	OurPrice    *Money          `json:"our_price"`
	WinnerPrice *Money          `json:"winner_price,omitempty"`
	TargetPrice *Money          `json:"target_price,omitempty"`
	Position    *MarketPosition `json:"position,omitempty"`
	FetchedAt   time.Time       `json:"fetched_at"`
}

type CompetitiveSignalInput struct {
	ListingID   string
	OurPrice    *Money
	WinnerPrice *Money
	TargetPrice *Money
	Position    *MarketPosition
	FetchedAt   time.Time
}

type MarketAggregate struct {
	ProductID  string                `json:"product_id"`
	Median     *Money                `json:"median,omitempty"`
	MinValid   *Money                `json:"min_valid,omitempty"`
	MaxValid   *Money                `json:"max_valid,omitempty"`
	NOffers    int                   `json:"n_offers"`
	NSellers   int                   `json:"n_sellers"`
	Source     MarketPriceSource     `json:"source"`
	FetchedAt  time.Time             `json:"fetched_at"`
	ComputedAt time.Time             `json:"computed_at"`
	Status     MarketAggregateStatus `json:"status"`
}

type MarketAggregateInput struct {
	ProductID  string
	Median     *Money
	MinValid   *Money
	MaxValid   *Money
	NOffers    int
	NSellers   int
	Source     MarketPriceSource
	FetchedAt  time.Time
	ComputedAt time.Time
	Status     MarketAggregateStatus
}

func NewMarketPriceSnapshot(input MarketPriceSnapshotInput) (MarketPriceSnapshot, error) {
	if !input.Scope.IsValid() {
		return MarketPriceSnapshot{}, fmt.Errorf("scope is invalid")
	}
	if strings.TrimSpace(input.RefID) == "" {
		return MarketPriceSnapshot{}, fmt.Errorf("ref_id is required")
	}
	if !input.Source.IsValid() {
		return MarketPriceSnapshot{}, fmt.Errorf("source is invalid")
	}
	if !input.Status.IsValid() {
		return MarketPriceSnapshot{}, fmt.Errorf("status is invalid")
	}
	if input.ObservedAt.IsZero() || input.FetchedAt.IsZero() || input.ExpiresAt.IsZero() {
		return MarketPriceSnapshot{}, fmt.Errorf("snapshot timestamps are required")
	}
	if strings.TrimSpace(input.RequestID) == "" {
		return MarketPriceSnapshot{}, fmt.Errorf("request_id is required")
	}
	if err := validatePositiveMoney("price", input.Price); err != nil {
		return MarketPriceSnapshot{}, err
	}
	if input.Status == MarketPriceSnapshotStatusFailed {
		if input.Price != nil {
			return MarketPriceSnapshot{}, fmt.Errorf("price must be absent for FAILED snapshot")
		}
		if strings.TrimSpace(pointerValue(input.FailureReason)) == "" {
			return MarketPriceSnapshot{}, fmt.Errorf("failure_reason is required for FAILED snapshot")
		}
	} else {
		if input.Price == nil {
			return MarketPriceSnapshot{}, fmt.Errorf("price is required for VALID/EXPIRED snapshot")
		}
		if input.FailureReason != nil {
			return MarketPriceSnapshot{}, fmt.Errorf("failure_reason only valid for FAILED snapshot")
		}
	}
	return MarketPriceSnapshot{
		Scope: input.Scope, RefID: strings.TrimSpace(input.RefID), Source: input.Source,
		Price: input.Price, ObservedAt: input.ObservedAt, FetchedAt: input.FetchedAt,
		ExpiresAt: input.ExpiresAt, Status: input.Status, FailureReason: input.FailureReason,
		RequestID: strings.TrimSpace(input.RequestID),
	}, nil
}

func NewValidatedOffer(input ValidatedOfferInput) (ValidatedOffer, error) {
	if strings.TrimSpace(input.CatalogProductID) == "" {
		return ValidatedOffer{}, fmt.Errorf("catalog_product_id is required")
	}
	if strings.TrimSpace(input.SellerID) == "" {
		return ValidatedOffer{}, fmt.Errorf("seller_id is required")
	}
	if strings.TrimSpace(input.Condition) == "" {
		return ValidatedOffer{}, fmt.Errorf("condition is required")
	}
	if !input.MatchStatus.IsValid() {
		return ValidatedOffer{}, fmt.Errorf("match_status is invalid")
	}
	if input.FetchedAt.IsZero() {
		return ValidatedOffer{}, fmt.Errorf("fetched_at is required")
	}
	if err := validatePositiveMoney("price", input.Price); err != nil {
		return ValidatedOffer{}, err
	}
	return ValidatedOffer{
		CatalogProductID: strings.TrimSpace(input.CatalogProductID), SellerID: strings.TrimSpace(input.SellerID),
		Price: input.Price, Condition: input.Condition, MatchStatus: input.MatchStatus, FetchedAt: input.FetchedAt,
	}, nil
}

func NewMatchDecision(input MatchDecisionInput) (MatchDecision, error) {
	if strings.TrimSpace(input.ProductID) == "" {
		return MatchDecision{}, fmt.Errorf("product_id is required")
	}
	if !input.Status.IsValid() {
		return MatchDecision{}, fmt.Errorf("match_status is invalid")
	}
	if !input.ConfidenceBand.IsValid() {
		return MatchDecision{}, fmt.Errorf("confidence_band is invalid")
	}
	if input.DecidedAt.IsZero() {
		return MatchDecision{}, fmt.Errorf("decided_at is required")
	}
	return MatchDecision{
		ProductID: strings.TrimSpace(input.ProductID), CatalogProductID: input.CatalogProductID,
		Status: input.Status, Anchors: input.Anchors, Contradictions: input.Contradictions,
		ConfidenceBand: input.ConfidenceBand, DecidedAt: input.DecidedAt,
	}, nil
}

func NewCompetitiveSignal(input CompetitiveSignalInput) (CompetitiveSignal, error) {
	if strings.TrimSpace(input.ListingID) == "" {
		return CompetitiveSignal{}, fmt.Errorf("listing_id is required")
	}
	if input.FetchedAt.IsZero() {
		return CompetitiveSignal{}, fmt.Errorf("fetched_at is required")
	}
	if input.OurPrice == nil {
		return CompetitiveSignal{}, fmt.Errorf("our_price is required")
	}
	for field, money := range map[string]*Money{
		"our_price": input.OurPrice, "winner_price": input.WinnerPrice, "target_price": input.TargetPrice,
	} {
		if err := validatePositiveMoney(field, money); err != nil {
			return CompetitiveSignal{}, err
		}
	}
	if input.Position != nil && (input.Position.Rank < 1 || input.Position.Total < 1 || input.Position.Rank > input.Position.Total) {
		return CompetitiveSignal{}, fmt.Errorf("position is invalid")
	}
	return CompetitiveSignal{
		ListingID: strings.TrimSpace(input.ListingID), OurPrice: input.OurPrice,
		WinnerPrice: input.WinnerPrice, TargetPrice: input.TargetPrice,
		Position: input.Position, FetchedAt: input.FetchedAt,
	}, nil
}

func NewMarketAggregate(input MarketAggregateInput) (MarketAggregate, error) {
	if strings.TrimSpace(input.ProductID) == "" {
		return MarketAggregate{}, fmt.Errorf("product_id is required")
	}
	if !input.Source.IsValid() {
		return MarketAggregate{}, fmt.Errorf("source is invalid")
	}
	if !input.Status.IsValid() {
		return MarketAggregate{}, fmt.Errorf("status is invalid")
	}
	if input.FetchedAt.IsZero() || input.ComputedAt.IsZero() {
		return MarketAggregate{}, fmt.Errorf("aggregate timestamps are required")
	}
	if input.NOffers < 0 || input.NSellers < 0 {
		return MarketAggregate{}, fmt.Errorf("aggregate counts must not be negative")
	}
	for field, money := range map[string]*Money{"median": input.Median, "min_valid": input.MinValid, "max_valid": input.MaxValid} {
		if err := validatePositiveMoney(field, money); err != nil {
			return MarketAggregate{}, err
		}
	}
	return MarketAggregate{
		ProductID: strings.TrimSpace(input.ProductID), Median: input.Median, MinValid: input.MinValid, MaxValid: input.MaxValid,
		NOffers: input.NOffers, NSellers: input.NSellers, Source: input.Source,
		FetchedAt: input.FetchedAt, ComputedAt: input.ComputedAt, Status: input.Status,
	}, nil
}

func (s MarketPriceScope) IsValid() bool {
	return s == MarketPriceScopeListing || s == MarketPriceScopeCatalogProduct
}

func (s MarketPriceSource) IsValid() bool {
	return s == MarketPriceSourceMLSalePrice || s == MarketPriceSourceMLPriceToWin || s == MarketPriceSourceMLCatalogOffers
}

func (s MarketPriceSnapshotStatus) IsValid() bool {
	return s == MarketPriceSnapshotStatusValid || s == MarketPriceSnapshotStatusFailed || s == MarketPriceSnapshotStatusExpired
}

func (s MatchStatus) IsValid() bool {
	return s == MatchStatusAccept || s == MatchStatusReview || s == MatchStatusReject || s == MatchStatusNoCandidate
}

func (s PriceEvidenceStatus) IsValid() bool {
	return s == PriceEvidenceStatusOK || s == PriceEvidenceStatusNoPriceEvidence || s == PriceEvidenceStatusInsufficientMarket
}

func (s BlockingState) IsValid() bool {
	return s == BlockingStateNoCandidate || s == BlockingStateNoPriceEvidence || s == BlockingStateInsufficientMarket || s == BlockingStateSemCusto
}

func (s ConfidenceBand) IsValid() bool {
	return s == ConfidenceBandAlta || s == ConfidenceBandMedia || s == ConfidenceBandBaixa
}

func (s MarketAggregateStatus) IsValid() bool {
	return s == MarketAggregateStatusOK || s == MarketAggregateStatusInsufficientMarket || s == MarketAggregateStatusNoPriceEvidence
}

func validatePositiveMoney(field string, money *Money) error {
	if money == nil {
		return nil
	}
	if err := validateMoney(field, money); err != nil {
		return err
	}
	rational, ok := new(big.Rat).SetString(money.Amount)
	if !ok || rational.Sign() <= 0 {
		return fmt.Errorf("%s.amount must be greater than zero", field)
	}
	return nil
}

func pointerValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
