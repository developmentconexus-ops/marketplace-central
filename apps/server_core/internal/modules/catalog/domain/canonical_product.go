package domain

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidInternalProductID = errors.New("invalid_identity")
var ErrInvalidSourceFact = errors.New("invalid_source_fact")

// InternalProductID is the canonical MPC product identity: Sankhya CODPROD.
type InternalProductID int

func NewInternalProductID(codprod int) (InternalProductID, error) {
	if codprod <= 0 {
		return 0, ErrInvalidInternalProductID
	}
	return InternalProductID(codprod), nil
}

type SourceFactQuality string

const (
	SourceFactQualityCurrent  SourceFactQuality = "current"
	SourceFactQualityStale    SourceFactQuality = "stale"
	SourceFactQualityUnknown  SourceFactQuality = "unknown"
	SourceFactQualityConflict SourceFactQuality = "conflict"
)

// NumericSourceFact preserves the distinction between an observed zero and an
// unknown value. Quality is assigned by the owning server-side domain.
type NumericSourceFact struct {
	Source        string            `json:"source"`
	Value         *float64          `json:"value"`
	Quality       SourceFactQuality `json:"quality"`
	ObservedAt    *time.Time        `json:"observed_at"`
	QualityReason *string           `json:"quality_reason"`
}

func NewNumericSourceFact(source string, value *float64, quality SourceFactQuality, observedAt *time.Time, qualityReason *string) (NumericSourceFact, error) {
	if strings.TrimSpace(source) == "" {
		return NumericSourceFact{}, ErrInvalidSourceFact
	}
	if quality != SourceFactQualityCurrent && quality != SourceFactQualityStale && quality != SourceFactQualityUnknown && quality != SourceFactQualityConflict {
		return NumericSourceFact{}, ErrInvalidSourceFact
	}
	needsReason := quality == SourceFactQualityUnknown || quality == SourceFactQualityStale
	if needsReason && (qualityReason == nil || strings.TrimSpace(*qualityReason) == "") {
		return NumericSourceFact{}, ErrInvalidSourceFact
	}
	if (quality == SourceFactQualityUnknown || quality == SourceFactQualityConflict) && value != nil {
		return NumericSourceFact{}, ErrInvalidSourceFact
	}
	if quality == SourceFactQualityUnknown && observedAt != nil {
		return NumericSourceFact{}, ErrInvalidSourceFact
	}
	if quality == SourceFactQualityCurrent && (value == nil || observedAt == nil) {
		return NumericSourceFact{}, ErrInvalidSourceFact
	}
	if quality == SourceFactQualityStale && (value == nil || observedAt == nil) {
		return NumericSourceFact{}, ErrInvalidSourceFact
	}
	return NumericSourceFact{Source: strings.TrimSpace(source), Value: value, Quality: quality, ObservedAt: observedAt, QualityReason: qualityReason}, nil
}

// CanonicalProduct is the contract populated by the Sankhya catalog read path.
// It intentionally has no legacy string product identifier.
type CanonicalProduct struct {
	InternalProductID     InternalProductID `json:"internal_product_id"`
	Name                  string            `json:"name"`
	EAN                   *string           `json:"ean"`
	ManufacturerReference *string           `json:"manufacturer_reference"`
	SellerSKU             *string           `json:"seller_sku"`
	BrandName             *string           `json:"brand_name"`
	NCM                   *string           `json:"ncm"`
	QualityFlags          []string          `json:"quality_flags"`
	ProductGroupName      *string           `json:"product_group_name"`
	CostAmount            NumericSourceFact `json:"cost_amount"`
	PriceAmount           NumericSourceFact `json:"price_amount"`
	StockQuantity         NumericSourceFact `json:"stock_quantity"`
}
