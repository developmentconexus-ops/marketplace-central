package application

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// BatchRunRequest is the input for RunBatch.
type BatchRunRequest struct {
	ProductIDs     []string
	PolicyIDs      []string
	OriginCEP      string
	DestCEP        string
	PriceSource    string
	PriceOverrides map[string]float64
	Modalidade     domain.Modalidade
}

// BatchSimulationItem is one product x policy result row.
type BatchSimulationItem struct {
	ProductID        string  `json:"product_id"`
	PolicyID         string  `json:"policy_id"`
	SellingPrice     float64 `json:"selling_price"`
	CostAmount       float64 `json:"cost_amount"`
	CommissionAmount float64 `json:"commission_amount"`
	FreightAmount    float64 `json:"freight_amount"`
	FixedFeeAmount   float64 `json:"fixed_fee_amount"`
	MarginAmount     float64 `json:"margin_amount"`
	MarginPercent    float64 `json:"margin_percent"`
	Status           string  `json:"status"`
	CommissionSource string  `json:"commission_source"`
	FreightSource    string  `json:"freight_source"`
}

// BatchRunResult holds all simulation rows.
type BatchRunResult struct {
	Items []BatchSimulationItem
}

// BatchOrchestrator runs batch simulations across all products x policies.
type BatchOrchestrator struct {
	products       ports.ProductProvider
	policies       ports.PolicyProvider
	freight        ports.FreightQuoter
	tariffResolver ports.TariffResolver // nil-safe: skipped when not wired, see WithTariffResolver
	tenantID       string
}

// NewBatchOrchestrator creates a BatchOrchestrator with its dependencies.
// The tariff resolver is attached separately via WithTariffResolver; when
// unset the orchestrator falls back to pol.CommissionPercent.
func NewBatchOrchestrator(
	products ports.ProductProvider,
	policies ports.PolicyProvider,
	freight ports.FreightQuoter,
	tenantID string,
) *BatchOrchestrator {
	return &BatchOrchestrator{products: products, policies: policies, freight: freight, tenantID: tenantID}
}

// WithTariffResolver attaches the shared ports.TariffResolver so batch
// commission is resolved by the ML listing category (the same seam the
// single-product decompose path uses), not the ERP taxonomy. Nil-safe:
// without it the orchestrator uses the policy's commission_percent.
func (o *BatchOrchestrator) WithTariffResolver(r ports.TariffResolver) *BatchOrchestrator {
	o.tariffResolver = r
	return o
}

// RunBatch calculates margins for every product x policy combination.
func (o *BatchOrchestrator) RunBatch(ctx context.Context, req BatchRunRequest) (BatchRunResult, error) {
	prods, err := o.products.GetProductsForBatch(ctx, req.ProductIDs)
	if err != nil {
		return BatchRunResult{}, fmt.Errorf("PRICING_BATCH_LOAD_PRODUCTS: %w", err)
	}

	pols, err := o.policies.GetPoliciesForBatch(ctx, req.PolicyIDs)
	if err != nil {
		return BatchRunResult{}, fmt.Errorf("PRICING_BATCH_LOAD_POLICIES: %w", err)
	}

	needsME := false
	for _, pol := range pols {
		if pol.ShippingProvider == "melhor_envio" {
			needsME = true
			break
		}
	}

	meConnected := false
	if needsME {
		connConnected, connErr := o.freight.IsConnected(ctx)
		if connErr == nil {
			meConnected = connConnected
		}
	}

	type freightState struct {
		amount    float64
		source    string
		available bool
	}

	freightResults := make(map[string]freightState)
	if needsME && meConnected {
		for _, p := range prods {
			if p.HeightCM == nil || p.WidthCM == nil || p.LengthCM == nil || p.WeightG == nil {
				continue
			}

			quoted, err := o.freight.QuoteFreight(ctx, ports.FreightRequest{
				OriginCEP: req.OriginCEP,
				DestCEP:   req.DestCEP,
				Products: []ports.FreightProduct{
					{
						ProductID: p.ProductID,
						HeightCM:  *p.HeightCM,
						WidthCM:   *p.WidthCM,
						LengthCM:  *p.LengthCM,
						WeightKg:  *p.WeightG / 1000,
						Value:     p.PriceAmount,
					},
				},
			})
			if err != nil {
				freightResults[p.ProductID] = freightState{source: "me_error"}
				continue
			}

			result, ok := quoted[p.ProductID]
			if !ok {
				freightResults[p.ProductID] = freightState{source: "me_error"}
				continue
			}

			source := result.Source
			if source == "" {
				source = "me_error"
			}
			freightResults[p.ProductID] = freightState{
				amount:    result.Amount,
				source:    source,
				available: source == "melhor_envio",
			}
		}
	}

	items := make([]BatchSimulationItem, 0, len(prods)*len(pols))
	for _, pol := range pols {
		for _, prod := range prods {
			sellingPrice := prod.PriceAmount
			if req.PriceSource == "suggested_price" && prod.SuggestedPrice != nil {
				sellingPrice = *prod.SuggestedPrice
			}
			if req.PriceOverrides != nil {
				if override, ok := req.PriceOverrides[prod.ProductID+"::"+pol.PolicyID]; ok && override > 0 {
					sellingPrice = override
				}
			}

			freightAmt := pol.DefaultShipping
			freightSource := "fixed"
			freightAvailable := true
			switch pol.ShippingProvider {
			case "melhor_envio":
				freightAmt = 0
				freightAvailable = false
				if prod.HeightCM == nil || prod.WidthCM == nil || prod.LengthCM == nil || prod.WeightG == nil {
					freightSource = "no_dimensions"
				} else if !meConnected {
					freightSource = "me_not_connected"
				} else if fr, ok := freightResults[prod.ProductID]; ok {
					freightAmt = fr.amount
					freightSource = fr.source
					freightAvailable = fr.available
				} else {
					freightSource = "me_error"
				}
			case "marketplace":
				freightAmt = pol.DefaultShipping
				freightSource = "marketplace"
				freightAvailable = true
			default:
				freightAmt = pol.DefaultShipping
				freightSource = "fixed"
				freightAvailable = true
			}

			var commissionAmt float64
			commissionSource := "policy"
			commissionKnown := true
			switch {
			case pol.CommissionOverride != nil:
				commissionAmt = round2(sellingPrice * *pol.CommissionOverride)
				commissionSource = "override"
			case o.tariffResolver != nil:
				var productID *int
				if id, perr := strconv.Atoi(prod.ProductID); perr == nil {
					productID = &id
				}
				priceBasis := strconv.FormatFloat(sellingPrice, 'f', 2, 64)
				res, rerr := o.tariffResolver.Resolve(ctx, ports.TariffRequest{
					Modalidade: req.Modalidade,
					ProductID:  productID,
					PriceBasis: &priceBasis,
				})
				if rerr != nil {
					return BatchRunResult{}, fmt.Errorf("PRICING_BATCH_RESOLVE_TARIFF: %w", rerr)
				}
				if valor := res.Comissao.Valor; valor != nil && strings.TrimSpace(*valor) != "" {
					amt, cerr := commissionDecimal(sellingPrice, strings.TrimSpace(*valor))
					if cerr != nil {
						return BatchRunResult{}, fmt.Errorf("PRICING_BATCH_COMMISSION_DECIMAL: %w", cerr)
					}
					commissionAmt = amt
					commissionSource = "resolver"
				} else {
					// ADR-17: resolved but no commission datum → honest unknown, never a default %.
					commissionKnown = false
					commissionSource = "unknown"
				}
			default:
				commissionAmt = round2(sellingPrice * pol.CommissionPercent)
			}

			var marginAmt, marginPct float64
			status := ""
			if !commissionKnown {
				status = statusCritical
			} else {
				marginAmt = sellingPrice - prod.CostAmount - commissionAmt - pol.FixedFeeAmount - freightAmt
				if sellingPrice > 0 {
					marginPct = marginAmt / sellingPrice
				}
				status = simulationStatusForBatch(marginPct, freightAvailable)
			}

			items = append(items, BatchSimulationItem{
				ProductID:        prod.ProductID,
				PolicyID:         pol.PolicyID,
				SellingPrice:     sellingPrice,
				CostAmount:       prod.CostAmount,
				CommissionAmount: commissionAmt,
				FreightAmount:    freightAmt,
				FixedFeeAmount:   pol.FixedFeeAmount,
				MarginAmount:     marginAmt,
				MarginPercent:    marginPct,
				Status:           status,
				CommissionSource: commissionSource,
				FreightSource:    freightSource,
			})
		}
	}

	return BatchRunResult{Items: items}, nil
}

// round2 rounds a monetary value half-up to 2 decimals, matching the
// domain decompose engine so batch commission equals the single-product
// decompose commission for the same resolver value (FIX-4 parity).
func round2(v float64) float64 { return math.Round(v*100) / 100 }

// commissionDecimal computes commission = price × pct% via the exact-decimal
// engine (domain.ParseRat + FormatRatHalfUp), the same routine the single-
// product decompose path uses, so the batch commission is bit-for-bit equal
// to decompose for the same resolved percent. pctStr is a percent decimal
// string ("15", "13.00"); price is formatted to 2dp before the exact multiply.
func commissionDecimal(price float64, pctStr string) (float64, error) {
	pct, err := domain.ParseRat(pctStr)
	if err != nil {
		return 0, err
	}
	priceRat, err := domain.ParseRat(fmt.Sprintf("%.2f", price))
	if err != nil {
		return 0, err
	}
	frac := new(big.Rat).Quo(pct, big.NewRat(100, 1))
	amount := new(big.Rat).Mul(priceRat, frac)
	f, err := strconv.ParseFloat(domain.FormatRatHalfUp(amount, 2), 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}
