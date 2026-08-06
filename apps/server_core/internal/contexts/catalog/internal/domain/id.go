// Package domain holds Catalog's model and its invariants. It is under
// internal/, so no other context can import it — not by convention, by the
// compiler refusing to link.
package domain

import (
	"errors"
	"fmt"
	"strings"
)

// ProductIDPrefix marks a canonical product identifier. The prefix is not
// decoration: it makes "someone passed the ERP code as the product id"
// detectable at the constructor instead of three joins downstream.
const ProductIDPrefix = "prd_"

var (
	// ErrBlankProductID is returned for an empty identifier.
	ErrBlankProductID = errors.New("catalog: product id is empty")
	// ErrNotCanonicalProductID is returned when the identifier does not carry
	// the canonical prefix.
	ErrNotCanonicalProductID = errors.New("catalog: product id is not canonical")
)

// ProductID is the platform's own identifier for a product. It is opaque and it
// derives from nothing at the source: no ERP code, no EAN, no SKU.
type ProductID struct {
	value string
}

// NewProductID validates the shape of a canonical identifier.
func NewProductID(s string) (ProductID, error) {
	v := strings.TrimSpace(s)
	if v == "" {
		return ProductID{}, ErrBlankProductID
	}
	if !strings.HasPrefix(v, ProductIDPrefix) || len(v) <= len(ProductIDPrefix) {
		return ProductID{}, fmt.Errorf("%w: %q has no %q prefix", ErrNotCanonicalProductID, v, ProductIDPrefix)
	}
	return ProductID{value: v}, nil
}

// String returns the identifier, or the empty string for the zero value.
func (p ProductID) String() string { return p.value }

// IsZero reports whether this is the zero value.
func (p ProductID) IsZero() bool { return p.value == "" }
