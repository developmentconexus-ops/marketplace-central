package domain_test

import (
	"testing"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
)

func TestValidFretePolicy(t *testing.T) {
	cases := []struct {
		name   string
		policy string
		want   bool
	}{
		{"estimativa is valid", domain.FretePolicyEstimativa, true},
		{"sem_dados is valid", domain.FretePolicySemDados, true},
		{"empty is invalid", "", false},
		{"unknown value is invalid", "outro", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := domain.ValidFretePolicy(tc.policy); got != tc.want {
				t.Fatalf("ValidFretePolicy(%q) = %v, want %v", tc.policy, got, tc.want)
			}
		})
	}
}
