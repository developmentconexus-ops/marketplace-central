package port_test

import (
	"testing"

	"marketplace-central/apps/server_core/internal/contexts/listings/port"
)

func TestZeroCursorIsStart(t *testing.T) {
	var c port.Cursor
	if !c.IsStart() {
		t.Fatal("zero cursor must be the start of the feed")
	}
}

func TestCursorRoundTripsToken(t *testing.T) {
	c := port.NewCursor("scroll-abc")
	if c.IsStart() || c.Token() != "scroll-abc" {
		t.Fatalf("cursor lost its token: start=%v token=%q", c.IsStart(), c.Token())
	}
}
