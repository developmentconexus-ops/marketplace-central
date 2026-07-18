package domain

import (
	"fmt"
	"strings"
	"time"
)

type VerdictInputRef struct {
	Source string
	AsOf   time.Time
}
type VerdictMarketRange struct {
	MinValid, Median  *Money
	Currency          string
	NOffers, NSellers int
}
type SnapshotRef struct {
	Source MarketPriceSource
	AsOf   time.Time
}
type CostRef struct{ AsOf time.Time }

type VerdictInput struct {
	ProductID   string
	Decision    MatchDecision
	Aggregate   *MarketAggregate
	OurPrice    *SnapshotRef
	TargetPrice *SnapshotRef
	Cost        *CostRef
}

type Verdict struct {
	ProductID           string
	MatchStatus         MatchStatus
	PriceEvidenceStatus PriceEvidenceStatus
	VerdictLabel        *string
	BlockingState       *BlockingState
	InputsUsed          map[string]VerdictInputRef
	MarketRange         *VerdictMarketRange
}

func AssembleVerdict(input VerdictInput) (Verdict, error) {
	if strings.TrimSpace(input.ProductID) == "" {
		return Verdict{}, fmt.Errorf("product_id is required")
	}
	if !input.Decision.Status.IsValid() {
		return Verdict{}, fmt.Errorf("match_status is invalid")
	}
	if err := validateVerdictInputRefs(input); err != nil {
		return Verdict{}, err
	}

	verdict := Verdict{ProductID: input.ProductID, MatchStatus: input.Decision.Status, InputsUsed: make(map[string]VerdictInputRef)}
	addVerdictRef(verdict.InputsUsed, "our_price", input.OurPrice)
	addVerdictRef(verdict.InputsUsed, "target_price", input.TargetPrice)
	if input.Cost != nil {
		verdict.InputsUsed["erp_cost"] = VerdictInputRef{Source: "erp_cost", AsOf: input.Cost.AsOf}
	}

	if input.Decision.Status != MatchStatusAccept {
		verdict.PriceEvidenceStatus = PriceEvidenceStatusNoPriceEvidence
		verdict.BlockingState = blockingState(BlockingStateNoCandidate)
		return verdict, nil
	}
	if input.Aggregate != nil && (input.Aggregate.Source != MarketPriceSourceMLCatalogOffers || !input.Aggregate.Status.IsValid() || input.Aggregate.FetchedAt.IsZero()) {
		return Verdict{}, fmt.Errorf("market aggregate reference is invalid")
	}
	if input.Aggregate == nil || input.Aggregate.Status == MarketAggregateStatusNoPriceEvidence {
		verdict.PriceEvidenceStatus = PriceEvidenceStatusNoPriceEvidence
		verdict.BlockingState = blockingState(BlockingStateNoPriceEvidence)
		return verdict, nil
	}

	verdict.MarketRange = &VerdictMarketRange{MinValid: input.Aggregate.MinValid, Median: input.Aggregate.Median, Currency: "BRL", NOffers: input.Aggregate.NOffers, NSellers: input.Aggregate.NSellers}
	verdict.InputsUsed["market"] = VerdictInputRef{Source: string(input.Aggregate.Source), AsOf: input.Aggregate.FetchedAt}
	if input.Aggregate.Status == MarketAggregateStatusInsufficientMarket {
		verdict.PriceEvidenceStatus = PriceEvidenceStatusInsufficientMarket
		verdict.BlockingState = blockingState(BlockingStateInsufficientMarket)
		return verdict, nil
	}
	verdict.PriceEvidenceStatus = PriceEvidenceStatusOK
	if input.Cost == nil {
		verdict.BlockingState = blockingState(BlockingStateSemCusto)
	}
	return verdict, nil
}

func validateVerdictInputRefs(input VerdictInput) error {
	for name, ref := range map[string]*SnapshotRef{"our_price": input.OurPrice, "target_price": input.TargetPrice} {
		if ref != nil && ((ref.Source != MarketPriceSourceMLSalePrice && ref.Source != MarketPriceSourceMLPriceToWin) || ref.AsOf.IsZero()) {
			return fmt.Errorf("%s reference is invalid", name)
		}
	}
	if input.Cost != nil && input.Cost.AsOf.IsZero() {
		return fmt.Errorf("erp_cost reference is invalid")
	}
	return nil
}

func addVerdictRef(inputs map[string]VerdictInputRef, key string, ref *SnapshotRef) {
	if ref != nil {
		inputs[key] = VerdictInputRef{Source: string(ref.Source), AsOf: ref.AsOf}
	}
}

func blockingState(value BlockingState) *BlockingState { return &value }
