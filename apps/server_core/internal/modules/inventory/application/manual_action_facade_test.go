package application

import (
	"context"
	"testing"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
	"marketplace-central/apps/server_core/internal/modules/inventory/ports"
)

type stubInstallationReader struct {
	ref   domain.ProviderStockRef
	found bool
}

func (s stubInstallationReader) GetInstallation(context.Context, string) (domain.ProviderStockRef, bool, error) {
	return s.ref, s.found, nil
}

type countingSnapshotStore struct{ calls int }

func (c *countingSnapshotStore) SaveListingSnapshot(context.Context, domain.ListingSnapshot) error {
	c.calls++
	return nil
}

func TestManualFacadeApprovedActionBlocksRiskWithoutOptimisticSnapshot(t *testing.T) {
	now := time.Now().UTC()
	risks := NewStockRiskService(
		stubSnapshotReader{items: []domain.ListingSnapshot{{
			Identity:          domain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"},
			ProviderCode:      "mercado_livre",
			AvailableQuantity: testIntPtr(9),
			ObservedAt:        &now,
		}}},
		stubLinkReader{items: []domain.LinkRecord{{
			Identity:          domain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"},
			State:             domain.LinkStateResolved,
			InternalProductID: testIntPtr(101),
		}}},
		&stubInternalStockReader{facts: map[int64]*ports.StockFact{101: {
			ProductID: 101,
			Quantity:  testFloatPtr(8),
			Source:    internalreaddomain.SourceMetadata{FetchedAt: now},
		}}},
	)
	envelope := &fakeStockEnvelope{protocolID: "prot-9"}
	actions := NewStockActionServiceWithEnvelope(&memoryActionStore{}, envelope)
	snapshots := &countingSnapshotStore{}
	facade := NewManualActionFacade(risks, actions, stubInstallationReader{ref: domain.ProviderStockRef{InstallationID: "inst-1"}, found: true}, snapshots)

	result, err := facade.Apply(context.Background(), ApplyManualActionInput{
		ActionID:          "act-facade-1",
		InstallationID:    "inst-1",
		ProviderItemID:    "MLB123",
		RequestedQuantity: 7,
		Reason:            "ajuste manual",
		Approval: domain.ManualApproval{
			Approved:   true,
			Operator:   domain.StockActionOperator{ActorID: "op-1"},
			ApprovedAt: now,
		},
	})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	if result.Action.State != domain.StockActionStateApproved {
		t.Fatalf("action state = %s, want approved (provider write is asynchronous)", result.Action.State)
	}
	if result.Action.MutationProtocolID == nil || *result.Action.MutationProtocolID != "prot-9" {
		t.Fatalf("mutation protocol id = %v, want prot-9", result.Action.MutationProtocolID)
	}
	if envelope.calls != 1 {
		t.Fatalf("envelope calls = %d, want 1", envelope.calls)
	}
	if snapshots.calls != 0 {
		t.Fatalf("snapshot writes = %d, want 0 — provider quantity is unknown until the poller applies the protocol", snapshots.calls)
	}
	if result.Risk.ProviderQuantity == nil || *result.Risk.ProviderQuantity != 9 {
		t.Fatalf("risk provider quantity = %v, want untouched 9", result.Risk.ProviderQuantity)
	}
	if result.Risk.Actionable || result.Risk.Actionability != domain.StockRiskActionabilityBlocked {
		t.Fatalf("risk actionable = %v / %s, want blocked while protocol pending", result.Risk.Actionable, result.Risk.Actionability)
	}
	if result.Risk.BlockingReason.Code != "mutation_protocol_pending" {
		t.Fatalf("blocking code = %q, want mutation_protocol_pending", result.Risk.BlockingReason.Code)
	}
}

func TestManualFacadeSkippedActionLeavesRiskUntouched(t *testing.T) {
	now := time.Now().UTC()
	risks := NewStockRiskService(
		stubSnapshotReader{items: []domain.ListingSnapshot{{
			Identity:          domain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"},
			AvailableQuantity: testIntPtr(9),
			ObservedAt:        &now,
		}}},
		stubLinkReader{items: []domain.LinkRecord{{
			Identity:          domain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"},
			State:             domain.LinkStateResolved,
			InternalProductID: testIntPtr(101),
		}}},
		&stubInternalStockReader{facts: map[int64]*ports.StockFact{101: {
			ProductID: 101,
			Quantity:  testFloatPtr(8),
			Source:    internalreaddomain.SourceMetadata{FetchedAt: now},
		}}},
	)
	envelope := &fakeStockEnvelope{protocolID: "prot-9"}
	actions := NewStockActionServiceWithEnvelope(&memoryActionStore{}, envelope)
	snapshots := &countingSnapshotStore{}
	facade := NewManualActionFacade(risks, actions, stubInstallationReader{ref: domain.ProviderStockRef{InstallationID: "inst-1"}, found: true}, snapshots)

	result, err := facade.Apply(context.Background(), ApplyManualActionInput{
		ActionID:          "act-facade-2",
		InstallationID:    "inst-1",
		ProviderItemID:    "MLB123",
		RequestedQuantity: 7,
		Approval:          domain.ManualApproval{Approved: false},
	})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if result.Action.State != domain.StockActionStateSkipped {
		t.Fatalf("action state = %s, want skipped", result.Action.State)
	}
	if envelope.calls != 0 || snapshots.calls != 0 {
		t.Fatalf("envelope calls = %d, snapshot writes = %d, want 0/0", envelope.calls, snapshots.calls)
	}
	if result.Risk.BlockingReason.Code == "mutation_protocol_pending" {
		t.Fatal("skipped action must not mark the risk as protocol-pending")
	}
}
