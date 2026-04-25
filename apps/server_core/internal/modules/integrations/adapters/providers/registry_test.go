package providers_test

import (
	"testing"

	_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/magalu"
	_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
	_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/shopee"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/amazon"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/leroymerlin"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/magalu"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/madeiramadeira"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/mercadolivre"
	"marketplace-central/apps/server_core/internal/modules/integrations/adapters/providers"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/shopee"
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
		"shopee":          domain.AuthStrategyAPIKey,
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

	first[0].ProviderCode = "mutated"
	first[0].DeclaredCapabilities[0] = "mutated_capability"
	first[0].Metadata["rollout_stage"] = "mutated"
	first = append(first, first[0])

	second := registry.All()
	if second[0].ProviderCode == "mutated" {
		t.Fatalf("ProviderCode mutation leaked back into registry")
	}
	if second[0].DeclaredCapabilities[0] == "mutated_capability" {
		t.Fatalf("DeclaredCapabilities mutation leaked back into registry")
	}
	if got := second[0].Metadata["rollout_stage"]; got == "mutated" {
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
	if got := shopee.Metadata["execution_mode"]; got != "blocked" {
		t.Fatalf("shopee execution_mode = %v, want blocked", got)
	}
	if got := shopee.Metadata["unavailable_reason"]; got == nil || got == "" {
		t.Fatalf("shopee unavailable_reason must be present when execution_mode is blocked")
	}
}
