package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
)

// The three ports below compose degrau 3 (live COTAÇÃO) of the tariff ladder:
// resolve a product's identity, turn that identity into a provider category,
// then quote a commission for that category at a price. Each returns an
// explicit found bool so the composite resolver can fall silently to degrau 4
// on any miss (ADR-17: unknown is never a fabricated tariff). All I/O is
// neutral (strings / domain.Money) — no provider payload crosses these ports;
// that stays in the adapters.

// ProductIdentity is the per-product data degrau 3 needs to find a category:
// the barcode and title feed the catalog match (D-83: both are sent together),
// and the optional catalog price is the fallback quote basis when the request
// carries no PriceBasis.
type ProductIdentity struct {
	EAN    string
	Titulo string
	// Preco is the product's known catalog price, or nil when absent
	// (ADR-17). Used only as the commission-quote basis of last resort.
	Preco *domain.Money
}

// ProductIdentityReader resolves a pricing product id to its identity. The
// bool is false (not an error) when the product is unknown or carries no
// identity, so the caller degrades to degrau 4 rather than failing.
type ProductIdentityReader interface {
	IdentityForProduct(ctx context.Context, productID int) (ProductIdentity, bool, error)
}

// FonteCategoria marks how a category was resolved. It is degrau-3-internal
// provenance (distinct from the wire domain.Fonte of the final tariff, which
// is COTACAO): the catmap source lands later behind the SAME CategoryResolver
// port with its own fonte value.
type FonteCategoria string

const (
	// FonteCategoriaEANCatalogo: category came from an EAN-anchored catalog
	// match.
	FonteCategoriaEANCatalogo FonteCategoria = "EAN-CATALOGO"
	// FonteCategoriaTitulo: category came from a title-only catalog match
	// (no EAN hit).
	FonteCategoriaTitulo FonteCategoria = "TITULO"
)

// CategoryQuery asks for a provider category from a product's barcode + title.
// Both are passed together (D-83); the resolver decides which anchored the hit
// and stamps the fonte accordingly.
type CategoryQuery struct {
	EAN    string
	Titulo string
}

// CategoryResolution is a resolved provider category stamped with its
// provenance and the source timestamp of the read.
type CategoryResolution struct {
	CategoriaID string
	Fonte       FonteCategoria
	// Confianca is the match confidence when the source exposes one (e.g. a
	// prediction rank), or nil when it does not (ADR-17 — not a fabricated 0).
	Confianca *string
	// Data is the source timestamp (RFC3339) of the category read.
	Data string
}

// CategoryResolver turns a product identity into a provider category. The bool
// is false (not an error) when no category could be matched, so the caller
// degrades to degrau 4.
type CategoryResolver interface {
	ResolveCategory(ctx context.Context, q CategoryQuery) (CategoryResolution, bool, error)
}

// CommissionQuery asks for a commission percent for a category at a price under
// a listing type (gold_special = clássico, gold_pro = premium/full — mapped
// from modalidade in the composition adapter, kept out of this neutral port).
type CommissionQuery struct {
	CategoriaID   string
	Preco         domain.Money
	ListingTypeID string
}

// CommissionQuote is a live-quoted commission percent stamped with the source
// timestamp of the fee read.
type CommissionQuote struct {
	// ComissaoPct is the commission percent as a decimal string (e.g.
	// "13.00").
	ComissaoPct string
	// Data is the source timestamp (RFC3339) of the fee quote.
	Data string
}

// CommissionQuoter quotes a live commission for a category+price+listing type.
// The bool is false (not an error) when the provider returns no commission, so
// the caller degrades to degrau 4.
type CommissionQuoter interface {
	QuoteCommission(ctx context.Context, q CommissionQuery) (CommissionQuote, bool, error)
}
