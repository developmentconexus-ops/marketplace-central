package application

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/marketplaces/domain"
	"marketplace-central/apps/server_core/internal/modules/marketplaces/ports"
)

type fakeRepository struct {
	listPoliciesByIDsCalls               int
	listPoliciesByIDsIDs                 [][]string
	policiesByID                         map[string]domain.Policy
	getPricingPolicyForInstallationCalls int
	pricingPolicyForInstallation         domain.Policy
	pricingPolicyForInstallationFound    bool
	pricingPolicyForInstallationErr      error
}

func (f *fakeRepository) SaveAccount(context.Context, domain.Account) error      { return nil }
func (f *fakeRepository) SavePolicy(context.Context, domain.Policy) error        { return nil }
func (f *fakeRepository) ListAccounts(context.Context) ([]domain.Account, error) { return nil, nil }
func (f *fakeRepository) ListPolicies(context.Context) ([]domain.Policy, error)  { return nil, nil }
func (f *fakeRepository) ListPoliciesByIDs(_ context.Context, ids []string) ([]domain.Policy, error) {
	f.listPoliciesByIDsCalls++
	f.listPoliciesByIDsIDs = append(f.listPoliciesByIDsIDs, append([]string(nil), ids...))
	result := make([]domain.Policy, 0, len(ids))
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		if policy, ok := f.policiesByID[id]; ok {
			result = append(result, policy)
		}
	}
	return result, nil
}

func (f *fakeRepository) GetPricingPolicyForInstallation(context.Context, string) (domain.Policy, bool, error) {
	f.getPricingPolicyForInstallationCalls++
	return f.pricingPolicyForInstallation, f.pricingPolicyForInstallationFound, f.pricingPolicyForInstallationErr
}

var _ ports.Repository = (*fakeRepository)(nil)

func TestService_ListPoliciesByIDs_EmptyIDsReturnsEmptyAndSkipsRepository(t *testing.T) {
	repo := &fakeRepository{}
	svc := NewService(repo, "tenant-1")

	policies, err := svc.ListPoliciesByIDs(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListPoliciesByIDs returned error: %v", err)
	}
	if len(policies) != 0 {
		t.Fatalf("expected empty result, got %d policies", len(policies))
	}
	if repo.listPoliciesByIDsCalls != 0 {
		t.Fatalf("expected repository not to be called, got %d calls", repo.listPoliciesByIDsCalls)
	}
}

func TestService_ListPoliciesByIDs_PreservesMissingAndDeduplicates(t *testing.T) {
	repo := &fakeRepository{policiesByID: map[string]domain.Policy{
		"m1": {PolicyID: "m1", AccountID: "a1"},
		"m3": {PolicyID: "m3", AccountID: "a2"},
	}}
	svc := NewService(repo, "tenant-1")

	policies, err := svc.ListPoliciesByIDs(context.Background(), []string{"m1", "missing", "m3", "m1"})
	if err != nil {
		t.Fatalf("ListPoliciesByIDs returned error: %v", err)
	}
	want := []domain.Policy{{PolicyID: "m1", AccountID: "a1"}, {PolicyID: "m3", AccountID: "a2"}}
	if !reflect.DeepEqual(policies, want) {
		t.Fatalf("unexpected policies: got %#v want %#v", policies, want)
	}
	if got := repo.listPoliciesByIDsIDs; !reflect.DeepEqual(got, [][]string{{"m1", "missing", "m3", "m1"}}) {
		t.Fatalf("unexpected repository IDs: %#v", got)
	}
}

func TestService_GetPricingPolicyForInstallation(t *testing.T) {
	t.Run("blank installation ID skips repository", func(t *testing.T) {
		repo := &fakeRepository{}
		svc := NewService(repo, "tenant-1")

		policy, found, err := svc.GetPricingPolicyForInstallation(context.Background(), " \t")
		if err != nil {
			t.Fatalf("GetPricingPolicyForInstallation returned error: %v", err)
		}
		if found || policy != (domain.Policy{}) {
			t.Fatalf("expected zero policy and not found, got %#v, %v", policy, found)
		}
		if repo.getPricingPolicyForInstallationCalls != 0 {
			t.Fatalf("expected repository not to be called, got %d calls", repo.getPricingPolicyForInstallationCalls)
		}
	})

	policy := domain.Policy{PolicyID: "p1", TenantID: "tenant-1", AccountID: "a1", MarketplaceCode: "shop", MinMarginPercent: 12.5}
	t.Run("returns found policy", func(t *testing.T) {
		repo := &fakeRepository{pricingPolicyForInstallation: policy, pricingPolicyForInstallationFound: true}
		svc := NewService(repo, "tenant-1")

		got, found, err := svc.GetPricingPolicyForInstallation(context.Background(), "installation-1")
		if err != nil {
			t.Fatalf("GetPricingPolicyForInstallation returned error: %v", err)
		}
		if !found || got != policy {
			t.Fatalf("expected found policy %#v, got %#v, %v", policy, got, found)
		}
		if repo.getPricingPolicyForInstallationCalls != 1 {
			t.Fatalf("expected one repository call, got %d", repo.getPricingPolicyForInstallationCalls)
		}
	})

	t.Run("returns not found without policy", func(t *testing.T) {
		repo := &fakeRepository{}
		svc := NewService(repo, "tenant-1")

		got, found, err := svc.GetPricingPolicyForInstallation(context.Background(), "installation-1")
		if err != nil {
			t.Fatalf("GetPricingPolicyForInstallation returned error: %v", err)
		}
		if found || got != (domain.Policy{}) {
			t.Fatalf("expected zero policy and not found, got %#v, %v", got, found)
		}
	})

	repoErr := errors.New("repository unavailable")
	t.Run("propagates repository error", func(t *testing.T) {
		repo := &fakeRepository{pricingPolicyForInstallationErr: repoErr}
		svc := NewService(repo, "tenant-1")

		got, found, err := svc.GetPricingPolicyForInstallation(context.Background(), "installation-1")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected repository error %v, got %v", repoErr, err)
		}
		if found || got != (domain.Policy{}) {
			t.Fatalf("expected zero policy and not found, got %#v, %v", got, found)
		}
	})
}
