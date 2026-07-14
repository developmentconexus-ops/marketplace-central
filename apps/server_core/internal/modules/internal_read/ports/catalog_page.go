package ports

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"
)

// ErrInvalidCursor is returned when an opaque catalog cursor cannot be
// decoded to a positive internal product id.
var ErrInvalidCursor = errors.New("invalid_cursor")

// Cursor is the keyset position for a catalog page. The zero value means the
// first page; non-zero values contain the last CODPROD returned to the caller.
type Cursor struct {
	InternalProductID int64
}

// DecodeCursor validates and decodes the base64 cursor envelope used by the
// catalog transport.
func DecodeCursor(encoded string) (Cursor, error) {
	if strings.TrimSpace(encoded) == "" {
		return Cursor{}, ErrInvalidCursor
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return Cursor{}, ErrInvalidCursor
	}
	id, err := strconv.ParseInt(string(decoded), 10, 64)
	if err != nil || id <= 0 {
		return Cursor{}, ErrInvalidCursor
	}
	return Cursor{InternalProductID: id}, nil
}

// Encode returns the opaque cursor representation for a positive product id.
func (c Cursor) Encode() (string, error) {
	if c.InternalProductID <= 0 {
		return "", ErrInvalidCursor
	}
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(c.InternalProductID, 10))), nil
}

type CatalogFactPage struct {
	Items      []CatalogProductFact
	NextCursor *Cursor
	AsOf       time.Time
}

type CatalogProductFact struct {
	InternalProductID int64
	Reference          *string
	Description        *string
	EAN                *string
	Active             bool
	SellableStock      CatalogQuantityFact
	CurrentPrice       CatalogMoneyFact
	Cost               CatalogMoneyFact
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

type CatalogPageReader interface {
	ListCatalogProductFacts(context.Context, Cursor, int) (CatalogFactPage, error)
	SearchCatalogProductFacts(context.Context, string, int) (CatalogFactPage, error)
}
