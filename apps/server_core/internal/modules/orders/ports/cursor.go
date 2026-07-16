package ports

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"
)

var ErrInvalidCursor = errors.New("invalid_cursor")

type InvalidCursorError struct{}

func (*InvalidCursorError) Error() string { return ErrInvalidCursor.Error() }
func (*InvalidCursorError) Unwrap() error { return ErrInvalidCursor }

// OrderCursor identifies the last row in newest-first order. Timestamp is
// the canonical order timestamp; ProviderOrderID is its deterministic tie-break.
type OrderCursor struct {
	Timestamp       time.Time
	ProviderOrderID string
}

type orderCursorEnvelope struct {
	Version         int       `json:"v"`
	Kind            string    `json:"kind"`
	Timestamp       time.Time `json:"timestamp"`
	ProviderOrderID string    `json:"provider_order_id"`
}

func (c OrderCursor) IsFirstPage() bool {
	return c.Timestamp.IsZero() && c.ProviderOrderID == ""
}

func (c OrderCursor) Encode() (string, error) {
	if c.IsFirstPage() || c.Timestamp.IsZero() || strings.TrimSpace(c.ProviderOrderID) != c.ProviderOrderID || c.ProviderOrderID == "" {
		return "", &InvalidCursorError{}
	}
	return encodeOrderCursor(orderCursorEnvelope{
		Version:         1,
		Kind:            "order",
		Timestamp:       c.Timestamp,
		ProviderOrderID: c.ProviderOrderID,
	})
}

func DecodeOrderCursor(encoded string) (OrderCursor, error) {
	if encoded == "" {
		return OrderCursor{}, &InvalidCursorError{}
	}
	payload, err := base64.StdEncoding.Strict().DecodeString(encoded)
	if err != nil || base64.StdEncoding.EncodeToString(payload) != encoded {
		return OrderCursor{}, &InvalidCursorError{}
	}

	var envelope orderCursorEnvelope
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&envelope); err != nil {
		return OrderCursor{}, &InvalidCursorError{}
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return OrderCursor{}, &InvalidCursorError{}
	}
	if envelope.Version != 1 || envelope.Kind != "order" || envelope.Timestamp.IsZero() || envelope.ProviderOrderID == "" || strings.TrimSpace(envelope.ProviderOrderID) != envelope.ProviderOrderID {
		return OrderCursor{}, &InvalidCursorError{}
	}
	return OrderCursor{Timestamp: envelope.Timestamp, ProviderOrderID: envelope.ProviderOrderID}, nil
}

func encodeOrderCursor(value orderCursorEnvelope) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", &InvalidCursorError{}
	}
	return base64.StdEncoding.EncodeToString(payload), nil
}
