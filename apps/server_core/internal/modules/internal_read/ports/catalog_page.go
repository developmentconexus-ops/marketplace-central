package ports

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"
)

// ErrInvalidCursor is the sentinel wrapped by InvalidCursorError.
var ErrInvalidCursor = errors.New("invalid_cursor")

// InvalidCursorError is returned when an opaque catalog cursor cannot be
// decoded to a positive internal product id.
type InvalidCursorError struct{}

func (*InvalidCursorError) Error() string { return ErrInvalidCursor.Error() }
func (*InvalidCursorError) Unwrap() error { return ErrInvalidCursor }

func invalidCursor() error { return &InvalidCursorError{} }

func NewInvalidCursorError() error { return invalidCursor() }

// Cursor is the keyset position for a catalog page. The zero value means the
// first page; non-zero values contain the last CODPROD returned to the caller.
type Cursor struct {
	InternalProductID int64
}

// DecodeCursor validates and decodes the base64 cursor envelope used by the
// catalog transport.
func DecodeCursor(encoded string) (Cursor, error) {
	encoded = strings.TrimSpace(encoded)
	if encoded == "" {
		return Cursor{}, invalidCursor()
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return Cursor{}, invalidCursor()
	}
	id, err := strconv.ParseInt(string(decoded), 10, 64)
	if err != nil || id <= 0 {
		return Cursor{}, invalidCursor()
	}
	return Cursor{InternalProductID: id}, nil
}

// Encode returns the opaque cursor representation for a positive product id.
func (c Cursor) Encode() (string, error) {
	if c.InternalProductID <= 0 {
		return "", invalidCursor()
	}
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(c.InternalProductID, 10))), nil
}

type CatalogFactPage struct {
	Items      []CatalogProductFact
	NextCursor *Cursor
	AsOf       time.Time
}

type CatalogProductFact struct {
	InternalProductID     int64
	Reference             *string
	ManufacturerReference *string
	Description           *string
	EAN                   *string
	BrandName             *string
	NCM                   *string
	QualityFlags          []string
	Active                bool
	SellableStock         CatalogQuantityFact
	CurrentPrice          CatalogMoneyFact
	Cost                  CatalogMoneyFact
}

type CatalogQuantityFact struct {
	Quantity *float64
	Quality  []string
}

type CatalogMoneyFact struct {
	Amount   *string
	Currency string
	Quality  []string
}

type SellableAssortmentPolicy struct {
	OnlyRevenda           bool
	OnlyEmEstoque         bool
	OnlyEcommerceEligible bool
}

// AllProductsAssortment is the honest "no cut" policy: every toggle off, so the
// caller sees the whole catalog. It is a real policy value, not a signal — a
// caller asking for it gets exactly it, the same as any other explicit policy.
func AllProductsAssortment() SellableAssortmentPolicy {
	return SellableAssortmentPolicy{}
}

// ErrUnresolvedAssortmentPolicy reports a nil policy reaching a seat below the
// routing reader. Nil means "resolve the tenant's stored policy", and routing is
// the one seat that can do that; past it, the policy is always concrete. So this
// is a wiring bug, and it is returned rather than papered over with a default —
// a default here would be exactly the invented rule this port stopped carrying.
var ErrUnresolvedAssortmentPolicy = errors.New("sellable assortment policy is unresolved below the routing seam")

// RequireAssortmentPolicy dereferences a policy that the routing invariant says
// is already concrete, and fails loudly when it is not.
func RequireAssortmentPolicy(policy *SellableAssortmentPolicy) (SellableAssortmentPolicy, error) {
	if policy == nil {
		return SellableAssortmentPolicy{}, ErrUnresolvedAssortmentPolicy
	}
	return *policy, nil
}

type CatalogAssortmentCounts struct {
	SellableCount int
	TotalCount    int
}

// CatalogPageReader is the whole contract a catalog source owes. Every read that
// pages the catalog carries the assortment policy, because every such read is
// answering "which products" and the cut is part of that answer.
//
// THE POLICY IS OPTIONAL AT THE REQUEST, NEVER AT THE READ. The pointer says
// which of two questions the caller is asking:
//
//	nil      -> "apply the tenant's stored assortment rule"
//	non-nil  -> "apply exactly this rule", and it IS applied, not inspected
//
// Nil resolves in ONE place, routing.Reader.resolveCatalogAssortment, which is
// the same seam that resolves the rule for linking (A-17: one producer, N
// consumers). Below it every seat receives a concrete policy, so a nil there is
// a wiring bug and returns ErrUnresolvedAssortmentPolicy.
//
// The shape this replaced is worth naming, because it reads as harmless: the
// parameter was a plain value, and a caller asking for the tenant's rule sent a
// hardcoded all-true constant that routing recognised by EQUALITY and threw
// away. Only one value was actually communicable — the zero one — and any
// explicit policy a caller passed was silently discarded. A wire boolean was
// crossing the port dressed as a policy, which is the two-mechanisms shape this
// port exists to keep out.
//
// The assortment methods were a separate, nominally optional port until F4. It
// was not optional in fact — seven of the nine implementations had it — and the
// optionality was maintained by hand, at the price of a runtime type assertion
// at every seat that forwarded it. Those assertions turned a mis-wiring into a
// 503 or a panic on a live request instead of a build error, which is how the
// catalog-503 defect was built. Folding the two ports is what lets the compiler
// certify the chain end to end; nothing here needs a hand-written
// `var _ ports.CatalogPageReader = ...` to hold it up any more.
type CatalogPageReader interface {
	ListCatalogProductFacts(context.Context, Cursor, int, *SellableAssortmentPolicy) (CatalogFactPage, error)
	// SearchCatalogProductFacts pages the matches for a query the same way
	// ListCatalogProductFacts pages the whole catalog. It takes a cursor because
	// a search that only ever returned the first page could not distinguish
	// "these are all the matches" from "these are the first 50 of many" — the
	// caller got a full page with no next cursor and no way to ask for the rest.
	SearchCatalogProductFacts(context.Context, string, Cursor, int, *SellableAssortmentPolicy) (CatalogFactPage, error)
	// CatalogProductFactsByIDs answers "the facts for exactly these products".
	// A screen that works on a known set — the products linked to listings, say —
	// cannot express that as a keyset page: the ids are scattered across the whole
	// catalog, so paging to them means reading everything in between. An id with
	// no row in the active source is simply absent from the result; the caller
	// learns which products the source does not know rather than being handed a
	// fabricated blank.
	CatalogProductFactsByIDs(context.Context, []int64) (CatalogFactPage, error)
	// GetCatalogAssortmentCounts reports the tenant's cut against the whole
	// catalog — the N and the M behind "Vendáveis N de M". It is the same
	// question the paged reads ask, asked once and without paging, so it takes
	// the policy on the same terms they do.
	GetCatalogAssortmentCounts(context.Context, *SellableAssortmentPolicy) (CatalogAssortmentCounts, error)
}
