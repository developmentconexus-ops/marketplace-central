package domain

import "testing"

func TestConnectorsMoney(t *testing.T) {
	t.Parallel()

	valid := &Money{Amount: "89.90", Currency: "BRL"}
	if err := ValidateMoney("price", valid); err != nil {
		t.Fatalf("valid money error = %v", err)
	}
	if err := ValidateMoney("price", nil); err != nil {
		t.Fatalf("nil money error = %v", err)
	}

	for name, money := range map[string]*Money{
		"non-decimal amount": {Amount: "89,90", Currency: "BRL"},
		"empty currency":     {Amount: "89.90", Currency: ""},
		"lowercase currency": {Amount: "89.90", Currency: "brl"},
	} {
		t.Run(name, func(t *testing.T) {
			if err := ValidateMoney("price", money); err == nil {
				t.Fatalf("ValidateMoney() error = nil, want validation error")
			}
		})
	}
}
