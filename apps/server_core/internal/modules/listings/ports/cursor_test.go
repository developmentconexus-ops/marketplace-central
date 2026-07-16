package ports

import (
	"errors"
	"testing"
)

func TestListingCursorRoundTrip(t *testing.T) {
	want := ListingCursor{LastTitle: "Ação Premium 漢字", ListingID: "inst_1~MLB123~-"}
	encoded, err := want.Encode()
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecodeListingCursor(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("round trip = %#v; want %#v", got, want)
	}
}

func TestListingCursorRejectsInvalidEnvelopes(t *testing.T) {
	values := []string{"", "garbage", "e30=", "eyJ2IjoxLCJraW5kIjoibGlzdGluZyJ9", "eyJ2IjoxLCJraW5kIjoibGlzdGluZyIsInRpdGxlIjoieCIsImxpc3RpbmdfaWQiOiJpIiwiaW50cnVkZXIiOnRydWV9", "eyJ2IjoyLCJraW5kIjoibGlzdGluZyIsInRpdGxlIjoieCIsImxpc3RpbmdfaWQiOiJpfnB+LSJ9", "eyJ2IjoxLCJraW5kIjoiZ3JvdXAiLCJ0aXRsZSI6IngiLCJsaXN0aW5nX2lkIjoiaX5wfi0ifQ==", "eyJ2IjoxLCJraW5kIjoibGlzdGluZyIsInRpdGxlIjoieCIsImxpc3RpbmdfaWQiOiJub3QtY2Fub25pY2FsIn0="}
	for _, value := range values {
		if _, err := DecodeListingCursor(value); !errors.Is(err, ErrInvalidCursor) {
			t.Errorf("DecodeListingCursor(%q) error = %v; want invalid_cursor", value, err)
		}
	}
}

func TestGroupCursorRoundTripAndValidation(t *testing.T) {
	want := GroupCursor{NullLast: false, ProductTitle: "Árvore", ProductID: "42664"}
	encoded, err := want.Encode()
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecodeGroupCursor(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("round trip = %#v; want %#v", got, want)
	}
	if _, err := DecodeGroupCursor("eyJ2IjoxLCJraW5kIjoibGlzdGluZyIsInRpdGxlIjoieCIsImxpc3RpbmdfaWQiOiJpfnB+LSJ9"); !errors.Is(err, ErrInvalidCursor) {
		t.Fatalf("wrong kind error = %v", err)
	}
}

func TestZeroCursorsMeanFirstPageButAreNeverEncoded(t *testing.T) {
	if !((ListingCursor{}).IsFirstPage()) || !((GroupCursor{}).IsFirstPage()) {
		t.Fatal("zero cursors must mean first page")
	}
	if encoded, err := (ListingCursor{}).Encode(); encoded != "" || !errors.Is(err, ErrInvalidCursor) {
		t.Fatalf("listing zero Encode() = %q, %v", encoded, err)
	}
	if encoded, err := (GroupCursor{}).Encode(); encoded != "" || !errors.Is(err, ErrInvalidCursor) {
		t.Fatalf("group zero Encode() = %q, %v", encoded, err)
	}
}
