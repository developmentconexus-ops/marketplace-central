package providers

import (
	"errors"
	"fmt"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	marketplacesports "marketplace-central/apps/server_core/internal/modules/marketplaces/ports"
)

type AuthFactory func() application.MarketplaceAuthAdapter
type FeeSyncerFactory func() marketplacesports.FeeScheduleSyncer

var (
	definitionEntries  []domain.ProviderDefinition
	authFactories      []AuthFactory
	feeSyncerFactories []FeeSyncerFactory
)

type Registry struct {
	definitions []domain.ProviderDefinition
	err         error
}

func RegisterDefinition(def domain.ProviderDefinition) {
	definitionEntries = append(definitionEntries, cloneProviderDefinition(def))
}

func RegisterAuthFactory(factory AuthFactory) {
	if factory == nil {
		return
	}
	authFactories = append(authFactories, factory)
}

func RegisterFeeSyncerFactory(factory FeeSyncerFactory) {
	if factory == nil {
		return
	}
	feeSyncerFactories = append(feeSyncerFactories, factory)
}

func NewRegistry() *Registry {
	defs, err := buildDefinitions(definitionEntries)
	return &Registry{
		definitions: defs,
		err:         err,
	}
}

func (r *Registry) Validate() error {
	if r == nil {
		return errors.New("INTEGRATIONS_PROVIDER_REGISTRY_NIL")
	}
	return r.err
}

func (r *Registry) All() []domain.ProviderDefinition {
	if r == nil || len(r.definitions) == 0 {
		return nil
	}

	out := make([]domain.ProviderDefinition, len(r.definitions))
	for i := range r.definitions {
		out[i] = cloneProviderDefinition(r.definitions[i])
	}
	return out
}

func BuildAuthAdapters() []application.MarketplaceAuthAdapter {
	out := make([]application.MarketplaceAuthAdapter, 0, len(authFactories))
	for _, factory := range authFactories {
		adapter := factory()
		if adapter == nil {
			continue
		}
		out = append(out, adapter)
	}
	return out
}

func BuildFeeSyncers() []marketplacesports.FeeScheduleSyncer {
	out := make([]marketplacesports.FeeScheduleSyncer, 0, len(feeSyncerFactories))
	for _, factory := range feeSyncerFactories {
		syncer := factory()
		if syncer == nil {
			continue
		}
		code := strings.TrimSpace(syncer.MarketplaceCode())
		if code == "" {
			continue
		}
		out = append(out, syncer)
	}
	return out
}

func buildDefinitions(entries []domain.ProviderDefinition) ([]domain.ProviderDefinition, error) {
	out := make([]domain.ProviderDefinition, 0, len(entries))
	seenCodes := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		code := strings.TrimSpace(entry.ProviderCode)
		if code == "" {
			return nil, errors.New("INTEGRATIONS_PROVIDER_CODE_REQUIRED")
		}
		if _, exists := seenCodes[code]; exists {
			return nil, fmt.Errorf("INTEGRATIONS_PROVIDER_DUPLICATE: provider_code=%s", code)
		}
		seenCodes[code] = struct{}{}
		out = append(out, cloneProviderDefinition(entry))
	}
	return out, nil
}

func cloneProviderDefinition(def domain.ProviderDefinition) domain.ProviderDefinition {
	cloned := def

	if def.DeclaredCapabilities != nil {
		cloned.DeclaredCapabilities = append([]string(nil), def.DeclaredCapabilities...)
	}

	if def.Metadata != nil {
		cloned.Metadata = make(map[string]any, len(def.Metadata))
		for key, value := range def.Metadata {
			cloned.Metadata[key] = cloneAny(value)
		}
	}

	return cloned
}

func cloneAny(value any) any {
	switch v := value.(type) {
	case map[string]any:
		cloned := make(map[string]any, len(v))
		for key, nested := range v {
			cloned[key] = cloneAny(nested)
		}
		return cloned
	case []any:
		cloned := make([]any, len(v))
		for i := range v {
			cloned[i] = cloneAny(v[i])
		}
		return cloned
	case []string:
		return append([]string(nil), v...)
	default:
		return value
	}
}
