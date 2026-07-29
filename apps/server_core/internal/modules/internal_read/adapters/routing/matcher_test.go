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
		lookup := fakeLookup{cfg: tenant_config.Config{
			TenantID: "t1",
			Source:   source,
			SetAt:    time.Now(),
			SellableAssortment: tenant_config.SellableAssortment{
				OnlyRevenda:           true,
				OnlyEmEstoque:         false,
				OnlyEcommerceEligible: true,
			},
		}}

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
		policy, ok := erpinternalread.SellableAssortmentFromContext(mirror.gotCtx)
		wantPolicy := erpinternalread.SellableAssortmentPolicy{OnlyRevenda: true, OnlyEcommerceEligible: true}
		if !ok || policy != wantPolicy {
			t.Fatalf("source %q: pinned assortment policy = %+v ok=%v, want %+v", source, policy, ok, wantPolicy)
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

// An anchor that resolves to nothing is "sem correspondência" for that listing;
// it must not abort generation for the whole account.
func TestMirrorMatcherTreatsAnAbsentProductAsNoCandidates(t *testing.T) {
	mirror := &fakeReader{err: &erpinternalread.ERPProductNotFoundError{ProductID: 0}}
	lookup := fakeLookup{cfg: tenant_config.Config{TenantID: "t1", Source: tenant_config.SourceSankhya, SetAt: time.Now()}}

	got, err := NewMirrorMatcher(mirror, lookup, "t1").FindProductsForLinking(context.Background(), internalreadports.FindProductsInput{SellerSKU: ptrString("81")})
	if err != nil || len(got) != 0 {
		t.Fatalf("candidates=%+v err=%v, want an empty result", got, err)
	}
}

func ptrString(v string) *string { return &v }

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
