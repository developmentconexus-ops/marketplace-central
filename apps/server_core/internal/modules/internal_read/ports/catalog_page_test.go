package ports

import (
	"errors"
	"testing"
)

func TestCatalogCursorRoundTripsLastProductID(t *testing.T) {
	want := Cursor{InternalProductID: 12345}
	encoded, err := want.Encode()
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	got, err := DecodeCursor(encoded)
	if err != nil {
		t.Fatalf("DecodeCursor() error = %v", err)
	}
	if got != want {
		t.Fatalf("DecodeCursor() = %+v, want %+v", got, want)
	}
}

func TestCatalogCursorRejectsInvalidValues(t *testing.T) {
	for _, encoded := range []string{"", "%%%garbage", "MA==", "LTU=", "MTAuMA=="} {
		if _, err := DecodeCursor(encoded); !errors.Is(err, ErrInvalidCursor) {
			t.Fatalf("DecodeCursor(%q) error = %v, want ErrInvalidCursor", encoded, err)
		}
	}
	if _, err := (Cursor{}).Encode(); !errors.Is(err, ErrInvalidCursor) {
		t.Fatalf("Cursor{}.Encode() error = %v, want ErrInvalidCursor", err)
	}
}
