package port

import (
	"context"
	"testing"

	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// stubFeed is a minimal ProductFeed implementation, only good enough to make
// the interface assertion below compile. It carries no behavior of its own.
type stubFeed struct{}

func (stubFeed) NextPage(ctx context.Context, t tenant.ID, after Cursor, limit int) (Page, error) {
	return Page{}, nil
}

// TestFeedPortHidesTheSourceKeyShape is the whole point of the port: the cursor
// is a token, not an ERP row id, so an adapter whose paging key is a string or a
// timestamp can implement it. The legacy readports.Cursor exposes
// InternalProductID and cannot.
func TestFeedPortHidesTheSourceKeyShape(t *testing.T) {
	var _ ProductFeed = stubFeed{}
	c := NewCursor("opaque-token")
	if c.Token() != "opaque-token" || c.IsStart() {
		t.Fatalf("cursor = %+v", c)
	}
	if !(Cursor{}).IsStart() {
		t.Fatalf("the zero cursor must mean start")
	}
}
