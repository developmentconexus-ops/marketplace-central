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

func DefaultSellableAssortment() SellableAssortmentPolicy {
	return SellableAssortmentPolicy{OnlyRevenda: true, OnlyEmEstoque: true, OnlyEcommerceEligible: true}
}

func AllProductsAssortment() SellableAssortmentPolicy {
	return SellableAssortmentPolicy{}
}

type CatalogAssortmentCounts struct {
	SellableCount int
	TotalCount    int
}

type CatalogPageReader interface {
	ListCatalogProductFacts(context.Context, Cursor, int) (CatalogFactPage, error)
	// SearchCatalogProductFacts pages the matches for a query the same way
	// ListCatalogProductFacts pages the whole catalog. It takes a cursor because
	// a search that only ever returned the first page could not distinguish
	// "these are all the matches" from "these are the first 50 of many" — the
	// caller got a full page with no next cursor and no way to ask for the rest.
	SearchCatalogProductFacts(context.Context, string, Cursor, int) (CatalogFactPage, error)
	// CatalogProductFactsByIDs answers "the facts for exactly these products".
	// A screen that works on a known set — the products linked to listings, say —
	// cannot express that as a keyset page: the ids are scattered across the whole
	// catalog, so paging to them means reading everything in between. An id with
	// no row in the active source is simply absent from the result; the caller
	// learns which products the source does not know rather than being handed a
	// fabricated blank.
	CatalogProductFactsByIDs(context.Context, []int64) (CatalogFactPage, error)
}

// CatalogAssortmentReader is OPTIONAL, and only the Oracle reader implements it
// today. Every seat between that reader and a caller re-asserts the port it
// forwards at RUNTIME — routing/reader.go resolveCatalogPager turns a failed
// assertion into ReadErrorSourceUnavailable, which is a 503 on the screen. So a
// decorator that does not implement this interface does not fail to compile; it
// deletes the capability and reports the source as unavailable. That is how the
// catalog-503 defect was built once already.
//
// VALIDITY CONDITION: whoever wires a consumer to this port owes a compile-time
// `var _ ports.CatalogAssortmentReader = ...` on EVERY seat in the live chain —
// cache.CatalogPageReader, observability.TimingReader, routing.Reader,
// application.Service — and an arrival test that reads the count through the
// composed reader from composition/root.go, not off the Oracle reader directly.
// A runtime `.(CatalogAssortmentReader)` with a fallback at the HTTP seam is the
// shape this condition exists to forbid: it turns a missing seat into a quiet
// wrong answer instead of a build error.
type CatalogAssortmentReader interface {
	ListCatalogProductFactsWithPolicy(context.Context, Cursor, int, SellableAssortmentPolicy) (CatalogFactPage, error)
	SearchCatalogProductFactsWithPolicy(context.Context, string, Cursor, int, SellableAssortmentPolicy) (CatalogFactPage, error)
	GetCatalogAssortmentCounts(context.Context, SellableAssortmentPolicy) (CatalogAssortmentCounts, error)
}
