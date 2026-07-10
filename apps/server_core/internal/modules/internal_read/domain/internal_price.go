package domain

import "time"

const (
	DefaultPriceTableID = 0
	DefaultPriceLocalID = 10101
)

type CurrentPricePolicy struct {
	PriceTableID int
	LocationID   int
	EffectiveAt  time.Time
}

func DefaultCurrentPricePolicy(now time.Time) CurrentPricePolicy {
	return CurrentPricePolicy{
		PriceTableID: DefaultPriceTableID,
		LocationID:   DefaultPriceLocalID,
		EffectiveAt:  now,
	}
}

type CurrentPrice struct {
	ProductID     int
	PriceTableID  int
	LocationID    int
	Amount        *float64
	VersionID     *int64
	EffectiveFrom *time.Time
	EffectiveTo   *time.Time
	Source        SourceMetadata
	QualityFlags  []QualityFlag
}
