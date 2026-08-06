package domain

import (
	"errors"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

var (
	// ErrForeignTenant is returned when an observation belongs to another tenant.
	ErrForeignTenant = errors.New("catalog: observation belongs to another tenant")
	// ErrNoProductID is returned when a product is built without an identifier.
	ErrNoProductID = errors.New("catalog: product id is required")
)

// Product is one canonical product. Every field is unexported and every method
// that changes state returns a NEW Product: a caller that persists inside a
// transaction which then fails must not be holding a mutated aggregate.
type Product struct {
	id          ProductID
	tenant      tenant.ID
	version     int
	description fact.Fact[string]
	identifiers []contracts.Identifier
	sourceKeys  []contracts.SourceProductKey
	lastHash    string
}

// NewProduct mints version 1 from a first observation.
func NewProduct(id ProductID, o contracts.ProductObservation) (Product, error) {
	if id.IsZero() {
		return Product{}, ErrNoProductID
	}
	if err := o.Validate(); err != nil {
		return Product{}, err
	}
	return Product{
		id:          id,
		tenant:      o.Key.Tenant(),
		version:     1,
		description: o.Description,
		identifiers: append([]contracts.Identifier(nil), o.Identifiers...),
		sourceKeys:  []contracts.SourceProductKey{o.Key},
		lastHash:    o.Evidence.PayloadHash(),
	}, nil
}

// ID returns the canonical identifier.
func (p Product) ID() ProductID { return p.id }

// Tenant returns the owning tenant.
func (p Product) Tenant() tenant.ID { return p.tenant }

// Version returns the catalog version, which counts substantive changes only.
func (p Product) Version() int { return p.version }

// Description returns the description as a fact, which may be Unknown.
func (p Product) Description() fact.Fact[string] { return p.description }

// Identifiers returns a copy of the identifier list.
func (p Product) Identifiers() []contracts.Identifier {
	return append([]contracts.Identifier(nil), p.identifiers...)
}

// SourceKeys returns a copy of every source address this product answers to.
func (p Product) SourceKeys() []contracts.SourceProductKey {
	return append([]contracts.SourceProductKey(nil), p.sourceKeys...)
}

// LastPayloadHash returns the hash of the last payload that changed this product.
func (p Product) LastPayloadHash() string { return p.lastHash }

// Apply folds a new observation in and reports what it did.
//
// Sameness is decided by the raw payload hash and not by comparing fields.
// Comparing fields means choosing which fields count, and that choice rots:
// the day a field is added, every product silently stops changing on it.
func (p Product) Apply(o contracts.ProductObservation) (Product, contracts.Disposition, error) {
	if err := o.Validate(); err != nil {
		return p, "", err
	}
	if o.Key.Tenant() != p.tenant {
		return p, "", fmt.Errorf("%w: product %s is %s, observation is %s",
			ErrForeignTenant, p.id, p.tenant, o.Key.Tenant())
	}
	if o.Evidence.PayloadHash() == p.lastHash {
		return p, contracts.DispositionIdempotent, nil
	}

	next := Product{
		id:          p.id,
		tenant:      p.tenant,
		version:     p.version + 1,
		description: o.Description,
		identifiers: mergeIdentifiers(p.identifiers, o.Identifiers),
		sourceKeys:  mergeSourceKeys(p.sourceKeys, o.Key),
		lastHash:    o.Evidence.PayloadHash(),
	}
	return next, contracts.DispositionChanged, nil
}

// mergeIdentifiers unions by kind+value, keeping the existing order stable.
func mergeIdentifiers(existing, incoming []contracts.Identifier) []contracts.Identifier {
	seen := make(map[string]bool, len(existing)+len(incoming))
	out := make([]contracts.Identifier, 0, len(existing)+len(incoming))
	for _, list := range [][]contracts.Identifier{existing, incoming} {
		for _, id := range list {
			k := id.String()
			if seen[k] {
				continue
			}
			seen[k] = true
			out = append(out, id)
		}
	}
	return out
}

// mergeSourceKeys appends the key unless this product already answers to it. A
// second source is added, never substituted: one canonical product legitimately
// has an ERP address and a spreadsheet address at the same time.
func mergeSourceKeys(existing []contracts.SourceProductKey, k contracts.SourceProductKey) []contracts.SourceProductKey {
	for _, have := range existing {
		if have.String() == k.String() {
			return append([]contracts.SourceProductKey(nil), existing...)
		}
	}
	return append(append([]contracts.SourceProductKey(nil), existing...), k)
}
