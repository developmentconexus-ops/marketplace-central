package main

import (
	"strings"
	"testing"
)

// TestRequireTenantConfigured_MissingEnv is the RED for D-39: catalogingest
// must fail closed when MC_DEFAULT_TENANT_ID is unset, instead of silently
// inheriting pgdb.LoadConfig's "tenant_default" substitution and writing a
// live batch of catalog rows under a fabricated tenant.
func TestRequireTenantConfigured_MissingEnv(t *testing.T) {
	getenv := func(string) string { return "" }

	err := requireTenantConfigured(getenv)
	if err == nil {
		t.Fatal("requireTenantConfigured() error = nil, want error naming MC_DEFAULT_TENANT_ID")
	}
	if !strings.Contains(err.Error(), "MC_DEFAULT_TENANT_ID") {
		t.Fatalf("requireTenantConfigured() error = %q, want it to name MC_DEFAULT_TENANT_ID", err.Error())
	}
}

// TestRequireTenantConfigured_TypoedVarNameStillFails proves the guard keys
// on the exact variable name pgdb.LoadConfig reads, not on the resolved
// value: a typo'd variable name must still fail closed.
func TestRequireTenantConfigured_TypoedVarNameStillFails(t *testing.T) {
	env := map[string]string{"MC_DEFAUTL_TENANT_ID": "tenant_default"} // typo, deliberate
	getenv := func(k string) string { return env[k] }

	if err := requireTenantConfigured(getenv); err == nil {
		t.Fatal("requireTenantConfigured() error = nil, want error: typo'd variable name must not satisfy the guard")
	}
}

// TestRequireTenantConfigured_LegitimateTenantDefaultValuePasses proves the
// guard checks presence of MC_DEFAULT_TENANT_ID, never its value: a
// deployment that legitimately names its tenant "tenant_default" must still
// be allowed to run.
func TestRequireTenantConfigured_LegitimateTenantDefaultValuePasses(t *testing.T) {
	env := map[string]string{"MC_DEFAULT_TENANT_ID": "tenant_default"}
	getenv := func(k string) string { return env[k] }

	if err := requireTenantConfigured(getenv); err != nil {
		t.Fatalf("requireTenantConfigured() error = %v, want nil: a legitimately-named tenant_default must pass", err)
	}
}
