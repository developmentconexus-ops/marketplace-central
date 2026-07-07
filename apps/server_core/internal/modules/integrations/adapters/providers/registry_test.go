package providers_test

import (
	"testing"

	_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/magalu"
	_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
	_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/shopee"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/amazon"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/leroymerlin"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/madeiramadeira"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/magalu"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/mercadolivre"
	"marketplace-central/apps/server_core/internal/modules/integrations/adapters/providers"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/shopee"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

func TestRegistryIncludesCoreProviders(t *testing.T) {
	t.Parallel()

	registry := providers.NewRegistry()
	if err := registry.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	defs := registry.All()
	if got, want := len(defs), 6; got != want {
		t.Fatalf("len(All()) = %d, want %d", got, want)
	}

	wantCodes := map[string]domain.AuthStrategy{
		"mercado_livre":   domain.AuthStrategyOAuth2,
		"magalu":          domain.AuthStrategyOAuth2,
		"shopee":          domain.AuthStrategyShopeePartner,
		"amazon":          domain.AuthStrategyLWA,
		"leroy_merlin":    domain.AuthStrategyAPIKey,
		"madeira_madeira": domain.AuthStrategyToken,
	}
	for _, def := range defs {
		wantStrategy, ok := wantCodes[def.ProviderCode]
		if !ok {
			t.Fatalf("unexpected provider code %q in registry", def.ProviderCode)
		}
		if def.AuthStrategy != wantStrategy {
			t.Fatalf("provider %q AuthStrategy = %q, want %q", def.ProviderCode, def.AuthStrategy, wantStrategy)
		}
		delete(wantCodes, def.ProviderCode)
	}
	if len(wantCodes) != 0 {
		t.Fatalf("missing expected providers: %v", wantCodes)
	}
}

func TestRegistryAllReturnsDefensiveCopy(t *testing.T) {
	t.Parallel()

	registry := providers.NewRegistry()
	first := registry.All()
	if len(first) == 0 {
		t.Fatal("All() returned no provider definitions")
	}

	mutableIndex := -1
	for i, def := range first {
		if len(def.DeclaredCapabilities) > 0 {
			mutableIndex = i
			break
		}
	}
	if mutableIndex == -1 {
		t.Fatal("All() returned no provider with declared capabilities")
	}

	first[mutableIndex].ProviderCode = "mutated"
	first[mutableIndex].DeclaredCapabilities[0] = "mutated_capability"
	first[mutableIndex].Metadata["rollout_stage"] = "mutated"
	first = append(first, first[mutableIndex])

	second := registry.All()
	if second[mutableIndex].ProviderCode == "mutated" {
		t.Fatalf("ProviderCode mutation leaked back into registry")
	}
	if second[mutableIndex].DeclaredCapabilities[0] == "mutated_capability" {
		t.Fatalf("DeclaredCapabilities mutation leaked back into registry")
	}
	if got := second[mutableIndex].Metadata["rollout_stage"]; got == "mutated" {
		t.Fatalf("Metadata mutation leaked back into registry")
	}
	if got, want := len(second), 6; got != want {
		t.Fatalf("len(All()) after mutation = %d, want %d", got, want)
	}
}

func TestBuildAuthAdaptersIncludesCoreAdapters(t *testing.T) {
	t.Parallel()

	adapters := providers.BuildAuthAdapters()
	if got, want := len(adapters), 6; got != want {
		t.Fatalf("len(BuildAuthAdapters()) = %d, want %d", got, want)
	}

	wantCodes := map[string]bool{
		"mercado_livre":   true,
		"magalu":          true,
		"shopee":          true,
		"amazon":          true,
		"leroy_merlin":    true,
		"madeira_madeira": true,
	}
	for _, adapter := range adapters {
		delete(wantCodes, adapter.ProviderCode())
	}
	if len(wantCodes) != 0 {
		t.Fatalf("missing expected auth adapters: %v", wantCodes)
	}
}

func TestBuildFeeSyncersIncludesCoreSyncers(t *testing.T) {
	t.Parallel()

	syncers := providers.BuildFeeSyncers()
	if got, want := len(syncers), 3; got != want {
		t.Fatalf("len(BuildFeeSyncers()) = %d, want %d", got, want)
	}

	wantCodes := map[string]bool{
		"mercado_livre": true,
		"magalu":        true,
		"shopee":        true,
	}
	for _, syncer := range syncers {
		delete(wantCodes, syncer.MarketplaceCode())
	}
	if len(wantCodes) != 0 {
		t.Fatalf("missing expected fee syncers: %v", wantCodes)
	}
}

func TestRegistryProviderMetadataCoverage(t *testing.T) {
	t.Parallel()

	registry := providers.NewRegistry()
	defs := registry.All()

	byCode := make(map[string]domain.ProviderDefinition, len(defs))
	for _, def := range defs {
		byCode[def.ProviderCode] = def
	}

	for _, code := range []string{
		"mercado_livre",
		"magalu",
		"shopee",
		"amazon",
		"leroy_merlin",
		"madeira_madeira",
	} {
		def, ok := byCode[code]
		if !ok {
			t.Fatalf("provider %q missing from registry", code)
		}
		if def.Metadata["country"] == nil {
			t.Fatalf("provider %q missing metadata.country", code)
		}
		if def.Metadata["rollout_stage"] == nil {
			t.Fatalf("provider %q missing metadata.rollout_stage", code)
		}
		if def.Metadata["execution_mode"] == nil {
			t.Fatalf("provider %q missing metadata.execution_mode", code)
		}
		if def.Metadata["fee_source"] == nil {
			t.Fatalf("provider %q missing metadata.fee_source", code)
		}
	}

	shopee := byCode["shopee"]
	if got := shopee.InstallMode; got != domain.InstallModeInteractive {
		t.Fatalf("shopee InstallMode = %v, want interactive", got)
	}
	if got := shopee.Metadata["rollout_stage"]; got != "v1" {
		t.Fatalf("shopee rollout_stage = %v, want v1", got)
	}
	if got := shopee.Metadata["execution_mode"]; got != "available" {
		t.Fatalf("shopee execution_mode = %v, want available", got)
	}
	if _, ok := shopee.Metadata["unavailable_reason"]; ok {
		t.Fatalf("shopee unavailable_reason should not be present")
	}
	if want := []string{"pricing_fee_sync"}; len(shopee.DeclaredCapabilities) != len(want) || shopee.DeclaredCapabilities[0] != want[0] {
		t.Fatalf("shopee DeclaredCapabilities = %#v, want %#v", shopee.DeclaredCapabilities, want)
	}

	mercadoLivre := byCode["mercado_livre"]
	wantMercadoLivre := []string{
		"listing_read",
		"pricing_fee_sync",
		"stock_read",
		"stock_write",
		"order_read",
		"message_read",
		"message_reply",
		"shipment_tracking",
		"webhook_receive",
	}
	if len(mercadoLivre.DeclaredCapabilities) != len(wantMercadoLivre) {
		t.Fatalf("mercado_livre DeclaredCapabilities len = %d, want %d", len(mercadoLivre.DeclaredCapabilities), len(wantMercadoLivre))
	}
	for i, want := range wantMercadoLivre {
		if got := mercadoLivre.DeclaredCapabilities[i]; got != want {
			t.Fatalf("mercado_livre DeclaredCapabilities[%d] = %q, want %q", i, got, want)
		}
	}

	for _, code := range []string{"amazon", "leroy_merlin", "madeira_madeira"} {
		if got := len(byCode[code].DeclaredCapabilities); got != 0 {
			t.Fatalf("%s DeclaredCapabilities len = %d, want 0", code, got)
		}
	}
}
