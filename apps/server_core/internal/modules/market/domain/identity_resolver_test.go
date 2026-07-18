package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestResolveIdentityTruthTable(t *testing.T) {
	decidedAt := time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	ptr := func(value string) *string { return &value }
	attrs := func(gtin, brand, model *string) []IdentityCandidateAttribute {
		return []IdentityCandidateAttribute{
			{ID: "GTIN", ValueName: gtin},
			{ID: "BRAND", ValueName: brand},
			{ID: "MODEL", ValueName: model},
		}
	}
	tests := []struct {
		name               string
		product            LocalProductIdentity
		candidates         []IdentityCandidate
		wantStatus         MatchStatus
		wantBand           ConfidenceBand
		wantCatalogID      *string
		wantAnchors        []string
		wantContradictions []string
	}{
		{"all anchors", LocalProductIdentity{"1", "Furadeira Doka 127V", ptr("7891"), ptr("XP1"), ptr("Dóka")}, []IdentityCandidate{{"MLA", "Furadeira Doka 127V", attrs(ptr("7891"), ptr("doka"), ptr("xp1"))}}, MatchStatusAccept, ConfidenceBandAlta, ptr("MLA"), []string{"ean", "marca", "refforn"}, []string{}},
		{"ean and brand", LocalProductIdentity{"2", "Serra Doka", ptr("7892"), nil, ptr("Doka")}, []IdentityCandidate{{"MLB", "Serra Doka", attrs(ptr("7892"), ptr("Doka"), nil)}}, MatchStatusAccept, ConfidenceBandMedia, ptr("MLB"), []string{"ean", "marca"}, []string{}},
		{"ean and model", LocalProductIdentity{"3", "Serra XP3", ptr("7893"), ptr("XP3"), nil}, []IdentityCandidate{{"MLC", "Serra XP3", attrs(ptr("7893"), nil, ptr("XP3"))}}, MatchStatusAccept, ConfidenceBandMedia, ptr("MLC"), []string{"ean", "refforn"}, []string{}},
		{"ean only", LocalProductIdentity{"4", "Produto", ptr("7894"), nil, nil}, []IdentityCandidate{{"MLD", "Outro", attrs(ptr("7894"), nil, nil)}}, MatchStatusReview, ConfidenceBandMedia, ptr("MLD"), []string{"ean"}, []string{}},
		{"ean absent two anchors capped", LocalProductIdentity{"5", "Produto", nil, ptr("M5"), ptr("Doka")}, []IdentityCandidate{{"MLE", "Outro", attrs(nil, ptr("Doka"), ptr("M5"))}}, MatchStatusReview, ConfidenceBandBaixa, ptr("MLE"), []string{"marca", "refforn"}, []string{}},
		{"ean absent weak", LocalProductIdentity{"6", "Produto", nil, nil, ptr("Doka")}, []IdentityCandidate{{"MLF", "Outro", attrs(nil, ptr("Doka"), nil)}}, MatchStatusReview, ConfidenceBandBaixa, ptr("MLF"), []string{"marca"}, []string{}},
		{"no candidates", LocalProductIdentity{"7", "Produto", nil, nil, nil}, nil, MatchStatusNoCandidate, ConfidenceBandBaixa, nil, []string{}, []string{}},
		{"voltage contradiction", LocalProductIdentity{"8", "Motor 127V", ptr("7898"), nil, nil}, []IdentityCandidate{{"MLG", "Motor 220V", attrs(ptr("7898"), nil, nil)}}, MatchStatusReject, ConfidenceBandMedia, ptr("MLG"), []string{"ean"}, []string{"voltagem"}},
		{"measure contradiction", LocalProductIdentity{"9", "Frasco 300ml", ptr("7899"), nil, nil}, []IdentityCandidate{{"MLH", "Frasco 500 ml", attrs(ptr("7899"), nil, nil)}}, MatchStatusReject, ConfidenceBandMedia, ptr("MLH"), []string{"ean"}, []string{"medida"}},
		{"color contradiction", LocalProductIdentity{"10", "Cadeira preta", ptr("7810"), nil, nil}, []IdentityCandidate{{"MLI", "Cadeira branca", attrs(ptr("7810"), nil, nil)}}, MatchStatusReject, ConfidenceBandMedia, ptr("MLI"), []string{"ean"}, []string{"cor"}},
		{"kit contradiction", LocalProductIdentity{"11", "Chave unidade", ptr("7811"), nil, nil}, []IdentityCandidate{{"MLJ", "Kit chave", attrs(ptr("7811"), nil, nil)}}, MatchStatusReject, ConfidenceBandMedia, ptr("MLJ"), []string{"ean"}, []string{"kit"}},
		{"doka menegotti collision", LocalProductIdentity{"12", "Betoneira Doka", ptr("7812"), nil, ptr("Doka")}, []IdentityCandidate{{"MLK", "Betoneira Menegotti", attrs(ptr("7812"), ptr("Menegotti"), nil)}}, MatchStatusReject, ConfidenceBandMedia, ptr("MLK"), []string{"ean"}, []string{"marca"}},
		{"doka volkswagen collision", LocalProductIdentity{"13", "Peça Doka", ptr("7813"), nil, ptr("Doka")}, []IdentityCandidate{{"MLL", "Peça VW", attrs(ptr("7813"), ptr("Volkswagen"), nil)}}, MatchStatusReject, ConfidenceBandMedia, ptr("MLL"), []string{"ean"}, []string{"marca"}},
		{"accept beats weak", LocalProductIdentity{"14", "Serra Doka", ptr("7814"), nil, ptr("Doka")}, []IdentityCandidate{{"weak", "Serra", attrs(nil, nil, nil)}, {"strong", "Serra Doka", attrs(ptr("7814"), ptr("Doka"), nil)}}, MatchStatusAccept, ConfidenceBandMedia, ptr("strong"), []string{"ean", "marca"}, []string{}},
		{"all contradicted stable score", LocalProductIdentity{"15", "Motor Doka 127V", ptr("7815"), nil, ptr("Doka")}, []IdentityCandidate{{"low", "Motor Menegotti 220V", attrs(nil, ptr("Menegotti"), nil)}, {"high", "Motor Menegotti 220V", attrs(ptr("7815"), ptr("Menegotti"), nil)}}, MatchStatusReject, ConfidenceBandMedia, ptr("high"), []string{"ean"}, []string{"voltagem", "marca"}},
		{"title similarity never anchors", LocalProductIdentity{"16", "Furadeira profissional Doka", nil, nil, nil}, []IdentityCandidate{{"MLM", "Furadeira profissional Doka", attrs(nil, nil, nil)}}, MatchStatusReview, ConfidenceBandBaixa, ptr("MLM"), []string{}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveIdentity(tt.product, tt.candidates, decidedAt)
			if err != nil {
				t.Fatalf("ResolveIdentity() error = %v", err)
			}
			if got.Status != tt.wantStatus || got.ConfidenceBand != tt.wantBand || !reflect.DeepEqual(got.CatalogProductID, tt.wantCatalogID) || !reflect.DeepEqual(got.Anchors, tt.wantAnchors) || !reflect.DeepEqual(got.Contradictions, tt.wantContradictions) {
				t.Fatalf("ResolveIdentity() = %#v, want status=%s band=%s catalog=%v anchors=%v contradictions=%v", got, tt.wantStatus, tt.wantBand, tt.wantCatalogID, tt.wantAnchors, tt.wantContradictions)
			}
			if got.Anchors == nil || got.Contradictions == nil {
				t.Fatalf("evidence slices must be non-nil: %#v", got)
			}
		})
	}
}

func TestResolveIdentityIsDeterministic(t *testing.T) {
	value := func(s string) *string { return &s }
	product := LocalProductIdentity{Codprod: "42", Descr: "Kit furadeira Dóka azul 127V", EAN: value("78942"), RefForn: value("XP42"), Marca: value("Dóka")}
	candidates := []IdentityCandidate{
		{CatalogProductID: "weak", Title: "Furadeira Doka azul 127V", Attrs: []IdentityCandidateAttribute{{ID: "GTIN", ValueName: value("78942")}}},
		{CatalogProductID: "strong", Title: "Kit furadeira Doka azul 127V", Attrs: []IdentityCandidateAttribute{{ID: "GTIN", ValueName: value("78942")}, {ID: "BRAND", ValueName: value("doka")}, {ID: "MODEL", ValueName: value("xp42")}}},
	}
	decidedAt := time.Date(2026, 7, 18, 13, 0, 0, 0, time.UTC)
	first, err := ResolveIdentity(product, candidates, decidedAt)
	if err != nil {
		t.Fatal(err)
	}
	second, err := ResolveIdentity(product, candidates, decidedAt)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("decisions differ:\nfirst:  %#v\nsecond: %#v", first, second)
	}
}

func TestResolveIdentityReturnsConstructorErrors(t *testing.T) {
	_, err := ResolveIdentity(LocalProductIdentity{}, nil, time.Time{})
	if err == nil {
		t.Fatal("ResolveIdentity() error = nil, want constructor validation error")
	}
}
