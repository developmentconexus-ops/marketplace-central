package ports

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"marketplace-central/apps/server_core/internal/modules/listings/domain"
)

var ErrInvalidCursor = errors.New("invalid_cursor")

type InvalidCursorError struct{}

func (*InvalidCursorError) Error() string { return ErrInvalidCursor.Error() }
func (*InvalidCursorError) Unwrap() error { return ErrInvalidCursor }

type ListingCursor struct {
	LastTitle string
	ListingID string
}

type GroupCursor struct {
	NullLast     bool
	ProductTitle string
	ProductID    string
}

type listingCursorEnvelope struct {
	Version   int    `json:"v"`
	Kind      string `json:"kind"`
	Title     string `json:"title"`
	ListingID string `json:"listing_id"`
}

type groupCursorEnvelope struct {
	Version      int    `json:"v"`
	Kind         string `json:"kind"`
	NullLast     bool   `json:"null_last"`
	ProductTitle string `json:"product_title"`
	ProductID    string `json:"product_id"`
}

func (c ListingCursor) IsFirstPage() bool { return c == ListingCursor{} }
func (c GroupCursor) IsFirstPage() bool   { return c == GroupCursor{} }

func (c ListingCursor) Encode() (string, error) {
	if c.IsFirstPage() || c.LastTitle == "" {
		return "", &InvalidCursorError{}
	}
	if _, err := domain.ParseListingID(c.ListingID); err != nil {
		return "", &InvalidCursorError{}
	}
	return encodeCursor(listingCursorEnvelope{Version: 1, Kind: "listing", Title: c.LastTitle, ListingID: c.ListingID})
}

func DecodeListingCursor(encoded string) (ListingCursor, error) {
	var envelope listingCursorEnvelope
	if err := decodeCursor(encoded, &envelope); err != nil || envelope.Version != 1 || envelope.Kind != "listing" || envelope.Title == "" {
		return ListingCursor{}, &InvalidCursorError{}
	}
	if _, err := domain.ParseListingID(envelope.ListingID); err != nil {
		return ListingCursor{}, &InvalidCursorError{}
	}
	return ListingCursor{LastTitle: envelope.Title, ListingID: envelope.ListingID}, nil
}

func (c GroupCursor) Encode() (string, error) {
	if c.IsFirstPage() || (!c.NullLast && (c.ProductTitle == "" || !canonicalPositiveInteger(c.ProductID))) || (c.NullLast && (c.ProductTitle != "" || c.ProductID != "")) {
		return "", &InvalidCursorError{}
	}
	return encodeCursor(groupCursorEnvelope{Version: 1, Kind: "group", NullLast: c.NullLast, ProductTitle: c.ProductTitle, ProductID: c.ProductID})
}

func DecodeGroupCursor(encoded string) (GroupCursor, error) {
	var envelope groupCursorEnvelope
	if err := decodeCursor(encoded, &envelope); err != nil || envelope.Version != 1 || envelope.Kind != "group" {
		return GroupCursor{}, &InvalidCursorError{}
	}
	cursor := GroupCursor{NullLast: envelope.NullLast, ProductTitle: envelope.ProductTitle, ProductID: envelope.ProductID}
	if cursor.IsFirstPage() || (!cursor.NullLast && (cursor.ProductTitle == "" || !canonicalPositiveInteger(cursor.ProductID))) || (cursor.NullLast && (cursor.ProductTitle != "" || cursor.ProductID != "")) {
		return GroupCursor{}, &InvalidCursorError{}
	}
	return cursor, nil
}

func canonicalPositiveInteger(value string) bool {
	id, err := strconv.ParseInt(value, 10, 64)
	return err == nil && id > 0 && strconv.FormatInt(id, 10) == value
}

func encodeCursor(value any) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", &InvalidCursorError{}
	}
	return base64.StdEncoding.EncodeToString(payload), nil
}

func decodeCursor(encoded string, target any) error {
	if encoded == "" {
		return &InvalidCursorError{}
	}
	payload, err := base64.StdEncoding.Strict().DecodeString(encoded)
	if err != nil || base64.StdEncoding.EncodeToString(payload) != encoded {
		return &InvalidCursorError{}
	}
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return &InvalidCursorError{}
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return &InvalidCursorError{}
	}
	return nil
}
