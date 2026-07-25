package routing

import (
	"context"
	"errors"
	"testing"
	"time"

	erpinternalread "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread"
	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/modules/tenant_config"
)

func TestMirrorMatcherPinsEveryActiveSourceIncludingTheLiveOne(t *testing.T) {
	for _, source := range []tenant_config.ActiveSource{
		tenant_config.SourceSankhya,
		tenant_config.SourceXLSX,
		tenant_config.SourceCatalogoCliente,
	} {
		mirror := &fakeReader{}
		lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: source, SetAt: time.Now()}}

		if _, err := NewMirrorMatcher(mirror, lookup, "t1").FindProductsForLinking(context.Background(), internalreadports.FindProductsInput{}); err != nil {
			t.Fatalf("source %q: FindProductsForLinking() error = %v", source, err)
		}
		if !mirror.called {
			t.Fatalf("source %q: the mirror reader was not consulted", source)
		}
		got, ok := erpinternalread.ActiveSourceFromContext(mirror.gotCtx)
		if !ok || got != erpdomain.ImportSource(source) {
			t.Fatalf("source %q: pinned erp source = %q ok=%v", source, got, ok)
		}
		if cfg, ok := tenant_config.FromContext(mirror.gotCtx); !ok || cfg.Source != source {
			t.Fatalf("source %q: pinned tenant config = %+v ok=%v", source, cfg, ok)
		}
	}
}

// No fallback between sources: an unknown source is an error, never another
// source's catalog.
func TestMirrorMatcherRefusesAnUnknownActiveSource(t *testing.T) {
	mirror := &fakeReader{}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.ActiveSource("erp_novo")}}

	_, err := NewMirrorMatcher(mirror, lookup, "t1").FindProductsForLinking(context.Background(), internalreadports.FindProductsInput{})
	if !errors.Is(err, tenant_config.ErrUnknownActiveSource) {
		t.Fatalf("error = %v, want ErrUnknownActiveSource", err)
	}
	if mirror.called {
		t.Fatal("an unknown source must not reach the mirror reader")
	}
}

func TestMirrorMatcherPropagatesLookupFailure(t *testing.T) {
	mirror := &fakeReader{}
	wantErr := errors.New("config unavailable")
	lookup := fakeLookup{err: wantErr}

	_, err := NewMirrorMatcher(mirror, lookup, "t1").FindProductsForLinking(context.Background(), internalreadports.FindProductsInput{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want the lookup failure", err)
	}
	if mirror.called {
		t.Fatal("a failed lookup must not reach the mirror reader")
	}
}
