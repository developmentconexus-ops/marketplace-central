package providers_test

import (
	"testing"

	_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/magalu"
	_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
	_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/shopee"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/magalu"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/mercadolivre"
	"marketplace-central/apps/server_core/internal/modules/integrations/adapters/providers"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/shopee"
)

func TestRegistryIncludesCoreProviders(t *testing.T) {
	t.Parallel()

	registry := providers.NewRegistry()
	if err := registry.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	defs := registry.All()
	if got, want := len(defs), 3; got != want {
		t.Fatalf("len(All()) = %d, want %d", got, want)
	}

	wantCodes := map[string]bool{
		"mercado_livre": true,
		"magalu":        true,
		"shopee":        true,
	}
	for _, def := range defs {
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
	first[0].Metadata["release_stage"] = "mutated"
	first = append(first, first[0])

	second := registry.All()
	if second[0].ProviderCode == "mutated" {
		t.Fatalf("ProviderCode mutation leaked back into registry")
	}
	if second[0].DeclaredCapabilities[0] == "mutated_capability" {
		t.Fatalf("DeclaredCapabilities mutation leaked back into registry")
	}
	if got := second[0].Metadata["release_stage"]; got == "mutated" {
		t.Fatalf("Metadata mutation leaked back into registry")
	}
	if got, want := len(second), 3; got != want {
		t.Fatalf("len(All()) after mutation = %d, want %d", got, want)
	}
}

func TestBuildAuthAdaptersIncludesCoreAdapters(t *testing.T) {
	t.Parallel()

	adapters := providers.BuildAuthAdapters()
	if got, want := len(adapters), 3; got != want {
		t.Fatalf("len(BuildAuthAdapters()) = %d, want %d", got, want)
	}

	wantCodes := map[string]bool{
		"mercado_livre": true,
		"magalu":        true,
		"shopee":        true,
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
