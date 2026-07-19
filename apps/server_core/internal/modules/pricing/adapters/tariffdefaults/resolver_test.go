package tariffdefaults

import (
	"context"
	"errors"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// fakeStore is an in-test stub of ports.TariffDefaultsStore (no DB).
type fakeStore struct {
	defaults domain.TariffDefaults
	err      error
}

func (f *fakeStore) GetTariffDefaults(ctx context.Context, tenantID, installationID string) (domain.TariffDefaults, error) {
	return f.defaults, f.err
}

func (f *fakeStore) UpsertTariffDefaults(ctx context.Context, tenantID, installationID string, in domain.TariffDefaults) (domain.TariffDefaults, error) {
	return f.defaults, f.err
}

func amountPtr(s string) *string { return &s }

func TestResolver_Resolve_Comissao(t *testing.T) {
	cases := []struct {
		name       string
		modalidade domain.Modalidade
		wantPct    string
	}{
		{"classico", domain.ModalidadeClassico, "13.00"},
		{"premium", domain.ModalidadePremium, "16.00"},
		{"full", domain.ModalidadeFull, "16.00"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeStore{defaults: domain.TariffDefaults{
				ComissaoClassicoPct: "13.00",
				ComissaoPremiumPct:  "16.00",
				FretePolicy:         domain.FretePolicySemDados,
			}}
			r := NewResolver(store, "tenant-1", "")

			got, err := r.Resolve(context.Background(), ports.TariffRequest{Modalidade: tc.modalidade})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Comissao.Valor == nil || *got.Comissao.Valor != tc.wantPct {
				t.Fatalf("Comissao.Valor = %v, want %q", got.Comissao.Valor, tc.wantPct)
			}
			if got.Comissao.Fonte != domain.FontePadrao {
				t.Fatalf("Comissao.Fonte = %v, want %v", got.Comissao.Fonte, domain.FontePadrao)
			}
			if got.Comissao.Degrau != 4 {
				t.Fatalf("Comissao.Degrau = %d, want 4", got.Comissao.Degrau)
			}
			if !got.Comissao.Estimativa {
				t.Fatalf("Comissao.Estimativa = false, want true")
			}
		})
	}
}

func TestResolver_Resolve_Frete(t *testing.T) {
	cases := []struct {
		name        string
		policy      string
		amount      *string
		wantValor   *string
		wantEstimat bool
	}{
		{"estimativa com valor", domain.FretePolicyEstimativa, amountPtr("21.50"), amountPtr("21.50"), true},
		{"estimativa sem valor", domain.FretePolicyEstimativa, nil, nil, true},
		{"sem_dados com valor presente ainda assim nil", domain.FretePolicySemDados, amountPtr("21.50"), nil, true},
		{"sem_dados sem valor", domain.FretePolicySemDados, nil, nil, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeStore{defaults: domain.TariffDefaults{
				ComissaoClassicoPct:   "13.00",
				ComissaoPremiumPct:    "16.00",
				FretePolicy:           tc.policy,
				FreteEstimativaAmount: tc.amount,
			}}
			r := NewResolver(store, "tenant-1", "")

			got, err := r.Resolve(context.Background(), ports.TariffRequest{Modalidade: domain.ModalidadeClassico})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantValor == nil {
				if got.Frete.Valor != nil {
					t.Fatalf("Frete.Valor = %v, want nil (NO-DATA)", *got.Frete.Valor)
				}
			} else {
				if got.Frete.Valor == nil || *got.Frete.Valor != *tc.wantValor {
					t.Fatalf("Frete.Valor = %v, want %q", got.Frete.Valor, *tc.wantValor)
				}
			}
			if got.Frete.Fonte != domain.FontePadrao {
				t.Fatalf("Frete.Fonte = %v, want %v", got.Frete.Fonte, domain.FontePadrao)
			}
			if got.Frete.Degrau != 4 {
				t.Fatalf("Frete.Degrau = %d, want 4", got.Frete.Degrau)
			}
			if got.Frete.Estimativa != tc.wantEstimat {
				t.Fatalf("Frete.Estimativa = %v, want %v", got.Frete.Estimativa, tc.wantEstimat)
			}
		})
	}
}

func TestResolver_Resolve_StoreErrorPropagates(t *testing.T) {
	wantErr := errors.New("boom")
	store := &fakeStore{err: wantErr}
	r := NewResolver(store, "tenant-1", "")

	_, err := r.Resolve(context.Background(), ports.TariffRequest{Modalidade: domain.ModalidadeClassico})
	if !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want %v", err, wantErr)
	}
}
