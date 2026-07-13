package domain

import (
	"errors"
	"testing"
)

const (
	testMPCLineOne MPCLineID = "mpl_11111111111111111111111111111111"
	testMPCLineTwo MPCLineID = "mpl_22222222222222222222222222222222"
)

func TestValidateAssistedSankhyaSelectionAcceptsOnlyCompleteBijection(t *testing.T) {
	order, candidate, selections := assistedSelectionFixture()

	mappings, err := ValidateAssistedSankhyaSelection(order, candidate, selections)
	if err != nil {
		t.Fatalf("ValidateAssistedSankhyaSelection() error = %v", err)
	}
	if len(mappings) != 2 || mappings[0].MPCLineID != testMPCLineOne || mappings[0].Origin.LineNumber != 2 || mappings[1].Origin.LineNumber != 1 {
		t.Fatalf("mappings = %#v, want exact explicit selection", mappings)
	}
}

func TestValidateAssistedSankhyaSelectionRejectsInvalidOrIncompleteProof(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*MarketplaceOrder, *AssistedSankhyaCandidate, *[]AssistedSankhyaLineSelection)
	}{
		{name: "legacy MPC line", mutate: func(order *MarketplaceOrder, _ *AssistedSankhyaCandidate, _ *[]AssistedSankhyaLineSelection) {
			order.Items[0].ReconciliationState = LineReconciliationLegacyUnresolved
		}},
		{name: "ambiguous MPC line", mutate: func(order *MarketplaceOrder, _ *AssistedSankhyaCandidate, _ *[]AssistedSankhyaLineSelection) {
			order.Items[0].ReconciliationState = LineReconciliationAmbiguous
		}},
		{name: "missing MPC mapping", mutate: func(_ *MarketplaceOrder, _ *AssistedSankhyaCandidate, selections *[]AssistedSankhyaLineSelection) {
			*selections = (*selections)[:1]
		}},
		{name: "extra candidate line", mutate: func(_ *MarketplaceOrder, candidate *AssistedSankhyaCandidate, _ *[]AssistedSankhyaLineSelection) {
			candidate.Lines = append(candidate.Lines, AssistedSankhyaCandidateLine{Identity: InternalDocumentLineIdentity{DocumentID: 31301, LineNumber: 3}})
		}},
		{name: "duplicate MPC mapping", mutate: func(_ *MarketplaceOrder, _ *AssistedSankhyaCandidate, selections *[]AssistedSankhyaLineSelection) {
			(*selections)[1].MPCLineID = (*selections)[0].MPCLineID
		}},
		{name: "duplicate selected candidate line", mutate: func(_ *MarketplaceOrder, _ *AssistedSankhyaCandidate, selections *[]AssistedSankhyaLineSelection) {
			(*selections)[1].CandidateLine = (*selections)[0].CandidateLine
		}},
		{name: "different candidate document", mutate: func(_ *MarketplaceOrder, _ *AssistedSankhyaCandidate, selections *[]AssistedSankhyaLineSelection) {
			(*selections)[0].CandidateLine.DocumentID = 31302
		}},
		{name: "duplicate exact candidate line", mutate: func(_ *MarketplaceOrder, candidate *AssistedSankhyaCandidate, _ *[]AssistedSankhyaLineSelection) {
			candidate.Lines[1].Identity = candidate.Lines[0].Identity
		}},
		{name: "wrong TOP", mutate: func(_ *MarketplaceOrder, candidate *AssistedSankhyaCandidate, _ *[]AssistedSankhyaLineSelection) {
			candidate.OperationCode = 306
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			order, candidate, selections := assistedSelectionFixture()
			test.mutate(&order, &candidate, &selections)
			if _, err := ValidateAssistedSankhyaSelection(order, candidate, selections); !errors.Is(err, ErrInvalidAssistedSankhyaLinkage) {
				t.Fatalf("error = %v, want invalid assisted linkage", err)
			}
		})
	}
}

func assistedSelectionFixture() (MarketplaceOrder, AssistedSankhyaCandidate, []AssistedSankhyaLineSelection) {
	order := MarketplaceOrder{Items: []MarketplaceOrderItem{
		{MPCLineID: testMPCLineOne, ReconciliationState: LineReconciliationStable},
		{MPCLineID: testMPCLineTwo, ReconciliationState: LineReconciliationStable},
	}}
	candidate := AssistedSankhyaCandidate{
		Header:        InternalDocumentIdentity{DocumentID: 31301},
		OperationCode: 313,
		Lines: []AssistedSankhyaCandidateLine{
			{Identity: InternalDocumentLineIdentity{DocumentID: 31301, LineNumber: 1}},
			{Identity: InternalDocumentLineIdentity{DocumentID: 31301, LineNumber: 2}},
		},
	}
	selections := []AssistedSankhyaLineSelection{
		{MPCLineID: testMPCLineOne, CandidateLine: candidate.Lines[1].Identity},
		{MPCLineID: testMPCLineTwo, CandidateLine: candidate.Lines[0].Identity},
	}
	return order, candidate, selections
}
