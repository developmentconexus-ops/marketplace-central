package domain

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestReconcileOrderLineIdentitiesUsesProviderIdentityNotMutableSnapshotFields(t *testing.T) {
	existing := []MarketplaceOrderItem{{
		MPCLineID:           fixedLineID(1),
		ReconciliationState: LineReconciliationStable,
		ProviderItemID:      "MLB1",
		ProviderVariationID: "variation",
		SellerSKU:           "sku",
		Quantity:            1,
		UnitPrice:           floatPointer(10),
	}}
	incoming := []MarketplaceOrderItem{{
		ProviderItemID:      "MLB1",
		ProviderVariationID: "variation",
		SellerSKU:           "sku",
		Quantity:            4,
		UnitPrice:           floatPointer(25),
	}}

	got, err := ReconcileOrderLineIdentities(existing, incoming, failingAllocator)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if got[0].MPCLineID != existing[0].MPCLineID || got[0].ReconciliationState != LineReconciliationStable {
		t.Fatalf("reconciled item = %#v, want retained stable identity", got[0])
	}
}

func TestReconcileOrderLineIdentitiesPreservesDuplicateIDMultisetAndAllocatesOnlyExcess(t *testing.T) {
	existing := []MarketplaceOrderItem{
		{MPCLineID: fixedLineID(2), ReconciliationState: LineReconciliationAmbiguous, ProviderItemID: "MLB1", SellerSKU: "sku"},
		{MPCLineID: fixedLineID(1), ReconciliationState: LineReconciliationAmbiguous, ProviderItemID: "MLB1", SellerSKU: "sku"},
	}
	incoming := []MarketplaceOrderItem{
		{ProviderItemID: "MLB1", SellerSKU: "sku", Quantity: 9},
		{ProviderItemID: "MLB1", SellerSKU: "sku", Quantity: 1},
		{ProviderItemID: "MLB1", SellerSKU: "sku", Quantity: 3},
	}
	allocations := 0
	got, err := ReconcileOrderLineIdentities(existing, incoming, func() (MPCLineID, error) {
		allocations++
		return fixedLineID(3), nil
	})
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if allocations != 1 {
		t.Fatalf("allocations = %d, want one new excess identity", allocations)
	}
	wantIDs := []MPCLineID{fixedLineID(1), fixedLineID(2), fixedLineID(3)}
	for index, item := range got {
		if item.MPCLineID != wantIDs[index] || item.ReconciliationState != LineReconciliationAmbiguous {
			t.Fatalf("item %d = %#v, want id %q ambiguous", index, item, wantIDs[index])
		}
	}
}

func TestReconcileOrderLineIdentitiesKeepsLegacyAndPriorAmbiguityExplicit(t *testing.T) {
	tests := []struct {
		name  string
		state LineReconciliationState
	}{
		{name: "legacy", state: LineReconciliationLegacyUnresolved},
		{name: "ambiguity does not silently clear", state: LineReconciliationAmbiguous},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			existing := []MarketplaceOrderItem{{MPCLineID: fixedLineID(1), ReconciliationState: test.state, ProviderItemID: "MLB1"}}
			got, err := ReconcileOrderLineIdentities(existing, []MarketplaceOrderItem{{ProviderItemID: "MLB1"}}, failingAllocator)
			if err != nil {
				t.Fatalf("reconcile: %v", err)
			}
			if got[0].MPCLineID != existing[0].MPCLineID || got[0].ReconciliationState != test.state {
				t.Fatalf("reconciled item = %#v, want retained %s state", got[0], test.state)
			}
		})
	}
}

func TestSankhyaLinkageValidateAndSemanticEquality(t *testing.T) {
	base := validSankhyaLinkage()
	if err := base.Validate(); err != nil {
		t.Fatalf("valid linkage rejected: %v", err)
	}
	retry := base
	retry.Audit.EventID = "event-retry"
	retry.Audit.RecordedAt = retry.Audit.RecordedAt.Add(time.Minute)
	retry.Lines = []SankhyaLineMapping{base.Lines[1], base.Lines[0]}
	if !base.SemanticallyEqual(retry) {
		t.Fatal("semantically identical retry was not equal")
	}
	retry.Lines[0].Origin.LineNumber = 99
	if base.SemanticallyEqual(retry) {
		t.Fatal("conflicting line reuse was semantically equal")
	}
	conflictingKey := base
	conflictingKey.ExternalOrderKey += "-other"
	if base.SemanticallyEqual(conflictingKey) {
		t.Fatal("conflicting external order key was semantically equal")
	}
}

func TestSankhyaLinkageRejectsUnknownAndDuplicateFacts(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*SankhyaLinkage)
	}{
		{name: "tenant", mutate: func(value *SankhyaLinkage) { value.Scope.TenantID = "" }},
		{name: "external order key", mutate: func(value *SankhyaLinkage) { value.ExternalOrderKey = "ml:v1:wrong" }},
		{name: "zero document", mutate: func(value *SankhyaLinkage) { value.Header.DocumentID = 0 }},
		{name: "unknown evidence", mutate: func(value *SankhyaLinkage) { value.Audit.EvidenceState = "" }},
		{name: "duplicate line", mutate: func(value *SankhyaLinkage) { value.Lines[1].MPCLineID = value.Lines[0].MPCLineID }},
		{name: "duplicate origin", mutate: func(value *SankhyaLinkage) { value.Lines[1].Origin = value.Lines[0].Origin }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := validSankhyaLinkage()
			test.mutate(&value)
			if err := value.Validate(); !errors.Is(err, ErrInvalidSankhyaLinkage) {
				t.Fatalf("validate error = %v, want invalid linkage", err)
			}
		})
	}
}

func validSankhyaLinkage() SankhyaLinkage {
	now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	return SankhyaLinkage{
		Scope:            LinkageScope{TenantID: "tenant", InstallationID: "installation", ProviderOrderID: "order"},
		ExternalOrderKey: "ml:v1:installation:order",
		Header:           InternalDocumentIdentity{DocumentID: 313001},
		Lines: []SankhyaLineMapping{
			{MPCLineID: fixedLineID(1), Origin: InternalDocumentLineIdentity{DocumentID: 313001, LineNumber: 1}},
			{MPCLineID: fixedLineID(2), Origin: InternalDocumentLineIdentity{DocumentID: 313001, LineNumber: 2}},
		},
		Audit: LinkageAudit{
			EventID: "event", EventType: LinkageEventConfirmed,
			ActorType: "operator", ActorID: "actor", Reason: "confirmed exact lines",
			IdempotencyKey: "idempotency", SourceAt: now, RecordedAt: now,
			ConfigurationRevision: "config-1", EvidenceState: LinkageEvidenceExact,
			EvidenceReference: "safe-ref",
		},
	}
}

func fixedLineID(value byte) MPCLineID {
	encoded := make([]byte, 32)
	for index := range encoded {
		encoded[index] = "0123456789abcdef"[value%16]
	}
	return MPCLineID("mpl_" + string(encoded))
}

func floatPointer(value float64) *float64 { return &value }

func failingAllocator() (MPCLineID, error) {
	return "", errors.New("allocator should not be called")
}

func TestCanonicalLineMappingsDoesNotMutateInput(t *testing.T) {
	value := validSankhyaLinkage()
	want := append([]SankhyaLineMapping(nil), value.Lines...)
	_ = canonicalLineMappings(value.Lines)
	if !reflect.DeepEqual(value.Lines, want) {
		t.Fatal("canonicalization mutated caller order")
	}
}
