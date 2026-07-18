package transport

import (
	"net/http"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

// handleSignals lists the latest competitive signal per requested listing.
// ADR-17: nullable Money/Position fields never omit their key.
func (h Handler) handleSignals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeMarketError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", map[string]any{})
		return
	}
	query, err := ParseSignalQuery(r.URL.Query())
	if err != nil {
		writeMarketQueryError(w, err)
		return
	}
	signals, err := h.evidence.Signals(r.Context(), query.ListingIDs)
	if err != nil {
		writeMarketServiceError(w, err)
		return
	}
	items := make([]signalWire, 0, len(signals))
	for _, s := range signals {
		items = append(items, newSignalWire(s))
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}

// handleAggregates lists the latest product-keyed market aggregate per
// requested codprod.
func (h Handler) handleAggregates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeMarketError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", map[string]any{})
		return
	}
	query, err := ParseCodprodQuery(r.URL.Query())
	if err != nil {
		writeMarketQueryError(w, err)
		return
	}
	aggregates, err := h.evidence.Aggregates(r.Context(), query.Codprods)
	if err != nil {
		writeMarketServiceError(w, err)
		return
	}
	items := make([]aggregateWire, 0, len(aggregates))
	for _, a := range aggregates {
		items = append(items, newAggregateWire(a))
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}

// handleVerdicts lists the assembled price-intel verdict per requested
// codprod. A never-collected codprod still yields a 200 NO_PRICE_EVIDENCE
// verdict (D-F4-q): it is never a 404.
func (h Handler) handleVerdicts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeMarketError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", map[string]any{})
		return
	}
	query, err := ParseCodprodQuery(r.URL.Query())
	if err != nil {
		writeMarketQueryError(w, err)
		return
	}
	verdicts, err := h.evidence.Verdicts(r.Context(), query.Codprods)
	if err != nil {
		writeMarketServiceError(w, err)
		return
	}
	items := make([]verdictWire, 0, len(verdicts))
	for _, v := range verdicts {
		items = append(items, newVerdictWire(v))
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}

type signalWire struct {
	ListingID   string                 `json:"listing_id"`
	OurPrice    *domain.Money          `json:"our_price"`
	WinnerPrice *domain.Money          `json:"winner_price"`
	TargetPrice *domain.Money          `json:"target_price"`
	Position    *domain.MarketPosition `json:"position"`
	FetchedAt   time.Time              `json:"fetched_at"`
}

func newSignalWire(s domain.CompetitiveSignal) signalWire {
	return signalWire{
		ListingID:   s.ListingID,
		OurPrice:    s.OurPrice,
		WinnerPrice: s.WinnerPrice,
		TargetPrice: s.TargetPrice,
		Position:    s.Position,
		FetchedAt:   s.FetchedAt.UTC(),
	}
}

type aggregateWire struct {
	ProductID  string                       `json:"product_id"`
	Median     *domain.Money                `json:"median"`
	MinValid   *domain.Money                `json:"min_valid"`
	NOffers    int                          `json:"n_offers"`
	NSellers   int                          `json:"n_sellers"`
	Source     domain.MarketPriceSource     `json:"source"`
	FetchedAt  time.Time                    `json:"fetched_at"`
	ComputedAt time.Time                    `json:"computed_at"`
	Status     domain.MarketAggregateStatus `json:"status"`
}

func newAggregateWire(a domain.MarketAggregate) aggregateWire {
	return aggregateWire{
		ProductID:  a.ProductID,
		Median:     a.Median,
		MinValid:   a.MinValid,
		NOffers:    a.NOffers,
		NSellers:   a.NSellers,
		Source:     a.Source,
		FetchedAt:  a.FetchedAt.UTC(),
		ComputedAt: a.ComputedAt.UTC(),
		Status:     a.Status,
	}
}

type verdictInputRefWire struct {
	Source string    `json:"source"`
	AsOf   time.Time `json:"as_of"`
}

type verdictMarketRangeWire struct {
	MinValid *domain.Money `json:"min_valid"`
	Median   *domain.Money `json:"median"`
	Currency string        `json:"currency"`
	NOffers  int           `json:"n_offers"`
	NSellers int           `json:"n_sellers"`
}

type verdictWire struct {
	MatchStatus         domain.MatchStatus             `json:"match_status"`
	PriceEvidenceStatus domain.PriceEvidenceStatus     `json:"price_evidence_status"`
	VerdictLabel        *string                        `json:"verdict_label"`
	BlockingState       *domain.BlockingState          `json:"blocking_state"`
	InputsUsed          map[string]verdictInputRefWire `json:"inputs_used"`
	MarketRange         *verdictMarketRangeWire        `json:"market_range"`
}

func newVerdictWire(v domain.Verdict) verdictWire {
	inputs := make(map[string]verdictInputRefWire, len(v.InputsUsed))
	for key, ref := range v.InputsUsed {
		inputs[key] = verdictInputRefWire{Source: ref.Source, AsOf: ref.AsOf.UTC()}
	}
	var marketRange *verdictMarketRangeWire
	if v.MarketRange != nil {
		marketRange = &verdictMarketRangeWire{
			MinValid: v.MarketRange.MinValid,
			Median:   v.MarketRange.Median,
			Currency: v.MarketRange.Currency,
			NOffers:  v.MarketRange.NOffers,
			NSellers: v.MarketRange.NSellers,
		}
	}
	return verdictWire{
		MatchStatus:         v.MatchStatus,
		PriceEvidenceStatus: v.PriceEvidenceStatus,
		VerdictLabel:        v.VerdictLabel,
		BlockingState:       v.BlockingState,
		InputsUsed:          inputs,
		MarketRange:         marketRange,
	}
}
