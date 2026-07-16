package ports

import (
	"encoding/base64"
	"errors"
	"testing"
	"time"
)

func TestOrderCursorRoundTrip(t *testing.T) {
	want := OrderCursor{
		Timestamp:       time.Date(2026, 7, 16, 12, 30, 0, 123456789, time.UTC),
		ProviderOrderID: "order-42",
	}

	encoded, err := want.Encode()
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecodeOrderCursor(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Timestamp.Equal(want.Timestamp) || got.ProviderOrderID != want.ProviderOrderID {
		t.Fatalf("round trip = %#v; want %#v", got, want)
	}
}

func TestDecodeOrderCursorRejectsMalformedWrongVersionAndTrailingJSON(t *testing.T) {
	valid := `{"v":1,"kind":"order","timestamp":"2026-07-16T12:30:00Z","provider_order_id":"order-42"}`
	cases := []string{
		"",
		"not-base64",
		encodeOrderCursorPayload(`{"v":2,"kind":"order","timestamp":"2026-07-16T12:30:00Z","provider_order_id":"order-42"}`),
		encodeOrderCursorPayload(valid + `{"extra":true}`),
		encodeOrderCursorPayload(`{"v":1,"kind":"order","timestamp":"not-a-time","provider_order_id":"order-42"}`),
	}

	for _, encoded := range cases {
		if _, err := DecodeOrderCursor(encoded); !errors.Is(err, ErrInvalidCursor) {
			t.Errorf("DecodeOrderCursor(%q) error = %v; want invalid_cursor", encoded, err)
		}
	}
}

func encodeOrderCursorPayload(payload string) string {
	return base64.StdEncoding.EncodeToString([]byte(payload))
}
