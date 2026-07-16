package marketplaces

import (
	"context"
	"errors"
	"testing"

	marketplacesdomain "marketplace-central/apps/server_core/internal/modules/marketplaces/domain"
	mutationsapp "marketplace-central/apps/server_core/internal/modules/mutations/application"
	mutationsdomain "marketplace-central/apps/server_core/internal/modules/mutations/domain"
	mutationsports "marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type fakePolicyService struct {
	calls  int
	policy marketplacesdomain.Policy
	found  bool
	err    error
}

func (f *fakePolicyService) GetPricingPolicyForInstallation(context.Context, string) (marketplacesdomain.Policy, bool, error) {
	f.calls++
	return f.policy, f.found, f.err
}

func TestReaderDelegatesAndPreservesPolicyAbsence(t *testing.T) {
	service := &fakePolicyService{}
	reader := NewPolicyReader(service)

	found, err := reader.HasPricingPolicy(context.Background(), "inst-1")
	if err != nil || found {
		t.Fatalf("found=%v err=%v, want absent without error", found, err)
	}
	if service.calls != 1 {
		t.Fatalf("service calls = %d, want 1", service.calls)
	}
	var _ mutationsports.PolicyReader = reader
}

func TestReaderReturnsFoundPolicyAndPropagatesErrors(t *testing.T) {
	service := &fakePolicyService{policy: marketplacesdomain.Policy{PolicyID: "policy-1"}, found: true}
	reader := NewPolicyReader(service)
	policy, found, err := reader.GetPricingPolicyForInstallation(context.Background(), "inst-1")
	if err != nil || !found || policy.PolicyID != "policy-1" {
		t.Fatalf("policy=%#v found=%v err=%v", policy, found, err)
	}

	want := errors.New("policy service unavailable")
	service.err = want
	_, found, err = reader.GetPricingPolicyForInstallation(context.Background(), "inst-1")
	if !errors.Is(err, want) || found {
		t.Fatalf("found=%v err=%v, want propagated error and absence", found, err)
	}
}

func TestMissingPolicyGateReturnsIC03CodeFromAdapter(t *testing.T) {
	service := &fakePolicyService{}
	reader := NewPolicyReader(service)
	err := mutationsapp.RequirePolicy(context.Background(), mutationsdomain.ProtocolTypePriceUpdate, "inst-1", reader.HasPricingPolicy)
	var gateErr *mutationsapp.GateError
	if !errors.As(err, &gateErr) || gateErr.Code != mutationsdomain.FailureCodePolicyMissing {
		t.Fatalf("error = %v, want policy_missing", err)
	}
	if service.calls != 1 {
		t.Fatalf("policy service calls = %d, want 1", service.calls)
	}
}
