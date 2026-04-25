package providers

import (
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

func TestBuildDefinitionsRejectsDuplicateProviderCode(t *testing.T) {
	t.Parallel()

	_, err := buildDefinitions([]domain.ProviderDefinition{
		{ProviderCode: "mercado_livre", TenantID: "system"},
		{ProviderCode: "mercado_livre", TenantID: "system"},
	})
	if err == nil {
		t.Fatal("buildDefinitions() error = nil, want duplicate error")
	}
	if !strings.Contains(err.Error(), "INTEGRATIONS_PROVIDER_DUPLICATE") {
		t.Fatalf("buildDefinitions() error = %v, want INTEGRATIONS_PROVIDER_DUPLICATE", err)
	}
}

func TestBuildDefinitionsRejectsEmptyProviderCode(t *testing.T) {
	t.Parallel()

	_, err := buildDefinitions([]domain.ProviderDefinition{
		{ProviderCode: " ", TenantID: "system"},
	})
	if err == nil {
		t.Fatal("buildDefinitions() error = nil, want INTEGRATIONS_PROVIDER_CODE_REQUIRED")
	}
	if !strings.Contains(err.Error(), "INTEGRATIONS_PROVIDER_CODE_REQUIRED") {
		t.Fatalf("buildDefinitions() error = %v, want INTEGRATIONS_PROVIDER_CODE_REQUIRED", err)
	}
}
