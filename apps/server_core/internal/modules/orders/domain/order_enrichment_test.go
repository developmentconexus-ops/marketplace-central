package domain

import "testing"

func TestMaskBuyer(t *testing.T) {
	cases := []struct {
		name     string
		nickname string
		want     MaskedBuyer
	}{
		{"single token", "JOAOSILVA", MaskedBuyer{Display: "JOAOSILVA"}},
		{"two tokens", "João Silva", MaskedBuyer{Display: "João S."}},
		{"three tokens uses only second", "João Silva Santos", MaskedBuyer{Display: "João S."}},
		{"blank", "   ", MaskedBuyer{}},
		{"empty", "", MaskedBuyer{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := MaskBuyer(tc.nickname)
			if got.Display != tc.want.Display {
				t.Fatalf("Display = %q, want %q", got.Display, tc.want.Display)
			}
			if got.City != nil {
				t.Fatalf("City = %v, want nil", got.City)
			}
			if got.UF != nil {
				t.Fatalf("UF = %v, want nil", got.UF)
			}
		})
	}
}

func TestDeriveVinculoStatus(t *testing.T) {
	cases := []struct {
		name  string
		items []MarketplaceOrderItem
		want  VinculoStatus
	}{
		{
			name:  "no items",
			items: nil,
			want:  VinculoStatusSemVinculo,
		},
		{
			name: "one resolved item",
			items: []MarketplaceOrderItem{
				{LinkQuality: LinkQualityUnresolved},
				{LinkQuality: LinkQualityResolved},
			},
			want: VinculoStatusOK,
		},
		{
			name: "no resolved item",
			items: []MarketplaceOrderItem{
				{LinkQuality: LinkQualityRejected},
				{LinkQuality: LinkQualityMissing},
				{LinkQuality: LinkQualityConflict},
			},
			want: VinculoStatusSemVinculo,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := DeriveVinculoStatus(tc.items)
			if got != tc.want {
				t.Fatalf("DeriveVinculoStatus() = %q, want %q", got, tc.want)
			}
		})
	}
}
