// Package tarifflive assembles degrau 3 (live COTAÇÃO) of the tariff ladder:
// it chains the identity -> category -> commission ports into a single
// commission ComponentResolution stamped {COTACAO, degrau 3, source
// timestamp}. Every step can miss; a miss returns found=false so the composite
// resolver falls to degrau 4 (ADR-17: never a fabricated tariff). Read errors
// propagate — the composite owns the fallback-on-error policy, not this
// assembler.
package tarifflive

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// Mercado Livre listing types the commission is quoted under. Clássico is
// gold_special; premium and full both settle as gold_pro.
const (
	listingTypeClassico = "gold_special"
	listingTypePremium  = "gold_pro"
)

// priceBasisCurrency is the settlement currency for the request price basis
// (the decompose price carries no explicit currency).
const priceBasisCurrency = "BRL"

// Resolver chains the three degrau-3 ports.
type Resolver struct {
	identities  ports.ProductIdentityReader
	categories  ports.CategoryResolver
	commissions ports.CommissionQuoter
}

func NewResolver(identities ports.ProductIdentityReader, categories ports.CategoryResolver, commissions ports.CommissionQuoter) *Resolver {
	return &Resolver{identities: identities, categories: categories, commissions: commissions}
}

// ResolveCommission resolves the degrau-3 commission component for req. It
// returns found=false (no error) whenever degrau 3 cannot honestly answer:
// no product id, unknown identity, no category match, no quotable price, an
// unmappable modalidade, or no commission quoted.
func (r *Resolver) ResolveCommission(ctx context.Context, req ports.TariffRequest) (domain.ComponentResolution, bool, error) {
	if req.ProductID == nil {
		return domain.ComponentResolution{}, false, nil
	}

	id, ok, err := r.identities.IdentityForProduct(ctx, *req.ProductID)
	if err != nil || !ok {
		return domain.ComponentResolution{}, false, err
	}

	cat, ok, err := r.categories.ResolveCategory(ctx, ports.CategoryQuery{EAN: id.EAN, Titulo: id.Titulo})
	if err != nil || !ok {
		return domain.ComponentResolution{}, false, err
	}

	price, ok := pickPrice(req.PriceBasis, id.Preco)
	if !ok {
		return domain.ComponentResolution{}, false, nil
	}

	listingType, ok := listingTypeFor(req.Modalidade)
	if !ok {
		return domain.ComponentResolution{}, false, nil
	}

	quote, ok, err := r.commissions.QuoteCommission(ctx, ports.CommissionQuery{
		CategoriaID:   cat.CategoriaID,
		Preco:         price,
		ListingTypeID: listingType,
	})
	if err != nil || !ok {
		return domain.ComponentResolution{}, false, err
	}

	valor := quote.ComissaoPct
	data := quote.Data
	return domain.ComponentResolution{
		Valor:  &valor,
		Fonte:  domain.FonteCotacao,
		Degrau: 3,
		Data:   &data,
	}, true, nil
}

// pickPrice chooses the commission quote basis: the request price basis first,
// then the product's catalog price. Absent both, degrau 3 cannot quote.
func pickPrice(basis *string, catalog *domain.Money) (domain.Money, bool) {
	if basis != nil {
		return domain.Money{Amount: *basis, Currency: priceBasisCurrency}, true
	}
	if catalog != nil {
		return *catalog, true
	}
	return domain.Money{}, false
}

// listingTypeFor maps a pricing modalidade to its ML listing type. An unknown
// modalidade is unmappable (degrau 3 skips rather than guessing).
func listingTypeFor(m domain.Modalidade) (string, bool) {
	switch m {
	case domain.ModalidadeClassico:
		return listingTypeClassico, true
	case domain.ModalidadePremium, domain.ModalidadeFull:
		return listingTypePremium, true
	default:
		return "", false
	}
}
