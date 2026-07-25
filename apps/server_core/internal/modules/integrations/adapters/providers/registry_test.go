package providers_test

import (
	"testing"

	_ "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
	_ "marketplace-central/apps/server_core/internal/modules/integrations/adapters/mercadolivre"
	"marketplace-central/apps/server_core/internal/modules/integrations/adapters/providers"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

// Mercado Livre is the only marketplace this platform integrates with. The
// registry used to carry five more (Magalu, Shopee, Amazon, Leroy Merlin,
// Madeira Madeira) whose adapters had no working read, write or auth path —
// they only published themselves as connectable providers and seeded invented
// commission rates, so the product offered accounts it could never sync.

func TestRegistryCarriesMercadoLivreOnly(t *testing.T) {
	t.Parallel()

	registry := providers.NewRegistry()
	if err := registry.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	defs := registry.All()
	if got, want := len(defs), 1; got != want {
		t.Fatalf("len(All()) = %d, want %d", got, want)
	}
	if got, want := defs[0].ProviderCode, "mercado_livre"; got != want {
		t.Fatalf("ProviderCode = %q, want %q", got, want)
	}
	if got, want := defs[0].AuthStrategy, domain.AuthStrategyOAuth2; got != want {
		t.Fatalf("AuthStrategy = %q, want %q", got, want)
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
	if got, want := len(second), 1; got != want {
		t.Fatalf("len(All()) after mutation = %d, want %d", got, want)
	}
}

func TestBuildAuthAdaptersCarriesMercadoLivreOnly(t *testing.T) {
	t.Parallel()

	adapters := providers.BuildAuthAdapters()
	if got, want := len(adapters), 1; got != want {
		t.Fatalf("len(BuildAuthAdapters()) = %d, want %d", got, want)
	}
	if got, want := adapters[0].ProviderCode(), "mercado_livre"; got != want {
		t.Fatalf("ProviderCode() = %q, want %q", got, want)
	}
}

// No marketplace seeds fee schedules any more. The ML seeder wrote a flat 16%
// Clássico / 22% Premium with a zero fixed fee, which is not how fees work here:
// the commission arrives per listing from the ML Fees API (sale_fee, per unit),
// the fixed fee is 50% of the price for items up to R$12,50, and the simulator
// reads its own tenant rates (pricing_tariff_defaults) — which the seeded rows
// contradicted (13%/16%). An empty schedule is honest; an invented one prices
// against rates nobody quoted.
func TestBuildFeeSyncersSeedsNothing(t *testing.T) {
	t.Parallel()

	if got := providers.BuildFeeSyncers(); len(got) != 0 {
		t.Fatalf("len(BuildFeeSyncers()) = %d, want 0", len(got))
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

	def, ok := byCode["mercado_livre"]
	if !ok {
		t.Fatal("provider \"mercado_livre\" missing from registry")
	}
	for _, key := range []string{"country", "rollout_stage", "execution_mode", "fee_source"} {
		if def.Metadata[key] == nil {
			t.Fatalf("provider \"mercado_livre\" missing metadata.%s", key)
		}
	}
	if got := def.InstallMode; got != domain.InstallModeInteractive {
		t.Fatalf("mercado_livre InstallMode = %v, want interactive", got)
	}

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
	if len(def.DeclaredCapabilities) != len(wantMercadoLivre) {
		t.Fatalf("mercado_livre DeclaredCapabilities len = %d, want %d", len(def.DeclaredCapabilities), len(wantMercadoLivre))
	}
	for i, want := range wantMercadoLivre {
		if got := def.DeclaredCapabilities[i]; got != want {
			t.Fatalf("mercado_livre DeclaredCapabilities[%d] = %q, want %q", i, got, want)
		}
	}
}
