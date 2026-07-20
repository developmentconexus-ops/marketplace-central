package tenant_config

import (
	"context"
	"testing"
	"time"
)

func TestWithActiveSource_FromContext_RoundTrip(t *testing.T) {
	cfg := Config{TenantID: "t1", Source: SourceXLSX, SetAt: time.Now(), SetBy: "operator"}
	ctx := WithActiveSource(context.Background(), cfg)

	got, ok := FromContext(ctx)
	if !ok {
		t.Fatal("FromContext ok = false, want true")
	}
	if got != cfg {
		t.Errorf("FromContext() = %+v, want %+v", got, cfg)
	}
}

func TestFromContext_Absent(t *testing.T) {
	_, ok := FromContext(context.Background())
	if ok {
		t.Error("FromContext ok = true on bare context, want false")
	}
}
