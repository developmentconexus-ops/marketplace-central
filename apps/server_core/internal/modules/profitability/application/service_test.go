package application

import (
	"context"
	"errors"
	"testing"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	profitabilitydomain "marketplace-central/apps/server_core/internal/modules/profitability/domain"
)

type stubOrderReader struct {
	orders []profitabilitydomain.OrderFact
}

func (s stubOrderReader) ListOrders(context.Context, string, int) ([]profitabilitydomain.OrderFact, error) {
	return s.orders, nil
}

type stubInternalFactReader struct {
	cost internalreaddomain.CostAsOf
	tax  internalreaddomain.TaxInputs
}

func (s stubInternalFactReader) GetCostAsOf(context.Context, int, time.Time) (internalreaddomain.CostAsOf, error) {
	return s.cost, nil
}
func (s stubInternalFactReader) GetTaxInputs(context.Context, int, time.Time) (internalreaddomain.TaxInputs, error) {
	return s.tax, nil
}

type recordingInternalFactReader struct {
	costCalls int
	taxCalls  int
}

func (s *recordingInternalFactReader) GetCostAsOf(context.Context, int, time.Time) (internalreaddomain.CostAsOf, error) {
	s.costCalls++
	amount := 99.0
	return internalreaddomain.CostAsOf{Amount: &amount}, nil
}

func (s *recordingInternalFactReader) GetTaxInputs(context.Context, int, time.Time) (internalreaddomain.TaxInputs, error) {
	s.taxCalls++
	amount := 11.0
	return internalreaddomain.TaxInputs{ICMSAmount: &amount, IPIAmount: &amount, PISAmount: &amount, COFINSAmount: &amount}, nil
}

type stubInputStore struct {
	items []profitabilitydomain.MarginInput
}

func (s *stubInputStore) ReplaceInputs(context.Context, string, []string, []profitabilitydomain.MarginInput) error {
	return nil
}
func (s *stubInputStore) ListInputs(context.Context, string, int) ([]profitabilitydomain.MarginInput, error) {
	return s.items, nil
}

type stubAdjustmentStore struct {
	items []profitabilitydomain.ManualAdjustment
}

func (s *stubAdjustmentStore) CreateOrGetAdjustment(_ context.Context, adjustment profitabilitydomain.ManualAdjustment) (profitabilitydomain.ManualAdjustment, error) {
	for _, item := range s.items {
		if item.InstallationID == adjustment.InstallationID && item.IdempotencyKey == adjustment.IdempotencyKey {
			return item, nil
		}
	}
	s.items = append(s.items, adjustment)
	return adjustment, nil
}
func (s *stubAdjustmentStore) ListAdjustments(context.Context, string, int) ([]profitabilitydomain.ManualAdjustment, error) {
	return s.items, nil
}

type stubSnapshotStore struct {
	items            []profitabilitydomain.ProfitSnapshot
	replacedOrderIDs []string
}

func (s *stubSnapshotStore) ReplaceSnapshots(_ context.Context, _ string, orderIDs []string, snapshots []profitabilitydomain.ProfitSnapshot) error {
	s.replacedOrderIDs = append([]string(nil), orderIDs...)
	s.items = snapshots
	return nil
}
func (s *stubSnapshotStore) ListSnapshots(context.Context, string, int) ([]profitabilitydomain.ProfitSnapshot, error) {
	return s.items, nil
}

func TestImportMarginInputsCompleteAndMissingLink(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	price := 10.0
	costAmount := 6.0
	icms := 1.0
	service := NewService(ServiceConfig{
		Orders: stubOrderReader{orders: []profitabilitydomain.OrderFact{
			{InstallationID: "inst-1", ProviderOrderID: "o1", ProviderUpdatedAt: &now, FetchedAt: now, Items: []profitabilitydomain.OrderItemFact{
				{ProviderItemID: "MLB1", Quantity: 2, UnitPrice: &price, SaleFeeAmount: &price, LinkQuality: profitabilitydomain.OrderLinkResolved, InternalProductID: intPtr(1)},
			}},
			{InstallationID: "inst-1", ProviderOrderID: "o2", ProviderUpdatedAt: &now, FetchedAt: now, Items: []profitabilitydomain.OrderItemFact{
				{ProviderItemID: "MLB2", Quantity: 1, UnitPrice: &price, LinkQuality: profitabilitydomain.OrderLinkMissing},
			}},
		}},
		Internal: stubInternalFactReader{
			cost: internalreaddomain.CostAsOf{Amount: &costAmount, Source: internalreaddomain.SourceMetadata{System: "oracle"}, QualityFlags: []internalreaddomain.QualityFlag{internalreaddomain.QualityComplete}},
			tax:  internalreaddomain.TaxInputs{ICMSAmount: &icms, Source: internalreaddomain.SourceMetadata{System: "oracle"}, QualityFlags: []internalreaddomain.QualityFlag{internalreaddomain.QualityComplete}},
		},
		Inputs:      &stubInputStore{},
		Adjustments: &stubAdjustmentStore{},
		Now:         func() time.Time { return now },
	})

	result, err := service.ImportMarginInputs(context.Background(), ImportMarginInputsInput{InstallationID: "inst-1", Limit: 2})
	if err != nil {
		t.Fatalf("ImportMarginInputs() error = %v", err)
	}
	if result.ImportedCount == 0 {
		t.Fatalf("expected imported inputs")
	}
	foundRevenue := false
	foundMissingCost := false
	for _, item := range result.Items {
		if item.ProviderOrderID == "o1" && item.Kind == profitabilitydomain.InputKindRevenue && item.Amount != nil && *item.Amount == 20 {
			foundRevenue = true
		}
		if item.ProviderOrderID == "o2" && item.Kind == profitabilitydomain.InputKindCost && item.Quality == profitabilitydomain.InputQualityMissing {
			foundMissingCost = true
		}
	}
	if !foundRevenue {
		t.Fatal("expected revenue input for resolved item")
	}
	if !foundMissingCost {
		t.Fatal("expected missing cost input for item without resolved link")
	}
}

func TestImportMarginInputsExtendsOnlyUnitCostByQuantity(t *testing.T) {
	now := time.Date(2026, 7, 12, 12, 0, 0, 0, time.UTC)
	unitCost := 10.0
	saleFee := 3.0
	icms, ipi, pis, cofins := 1.0, 2.0, 3.0, 4.0
	service := NewService(ServiceConfig{
		Orders: stubOrderReader{orders: []profitabilitydomain.OrderFact{{
			InstallationID: "inst-1", ProviderOrderID: "order-1", ProviderUpdatedAt: &now, FetchedAt: now,
			Items: []profitabilitydomain.OrderItemFact{
				{ProviderItemID: "quantity-1", Quantity: 1, UnitPrice: floatPtr(100), SaleFeeAmount: &saleFee, LinkQuality: profitabilitydomain.OrderLinkResolved, InternalProductID: intPtr(1)},
				{ProviderItemID: "quantity-2", Quantity: 2, UnitPrice: floatPtr(100), SaleFeeAmount: &saleFee, LinkQuality: profitabilitydomain.OrderLinkResolved, InternalProductID: intPtr(2)},
				{ProviderItemID: "quantity-7", Quantity: 7, UnitPrice: floatPtr(100), SaleFeeAmount: &saleFee, LinkQuality: profitabilitydomain.OrderLinkResolved, InternalProductID: intPtr(7)},
			},
		}}},
		Internal: stubInternalFactReader{
			cost: internalreaddomain.CostAsOf{Amount: &unitCost, AmountScope: internalreaddomain.CostAmountScopePerUnit, Source: internalreaddomain.SourceMetadata{System: "oracle"}, QualityFlags: []internalreaddomain.QualityFlag{internalreaddomain.QualityComplete}},
			tax:  internalreaddomain.TaxInputs{ICMSAmount: &icms, IPIAmount: &ipi, PISAmount: &pis, COFINSAmount: &cofins, Source: internalreaddomain.SourceMetadata{System: "oracle"}, QualityFlags: []internalreaddomain.QualityFlag{internalreaddomain.QualityComplete}},
		},
		Inputs: &stubInputStore{},
		Now:    func() time.Time { return now },
	})

	result, err := service.ImportMarginInputs(context.Background(), ImportMarginInputsInput{InstallationID: "inst-1"})
	if err != nil {
		t.Fatalf("ImportMarginInputs() error = %v", err)
	}
	wantCosts := map[string]float64{"quantity-1": 10, "quantity-2": 20, "quantity-7": 70}
	wantTaxes := map[profitabilitydomain.InputKind]float64{
		profitabilitydomain.InputKindTaxICMS: 1, profitabilitydomain.InputKindTaxIPI: 2,
		profitabilitydomain.InputKindTaxPIS: 3, profitabilitydomain.InputKindTaxCOFINS: 4,
	}
	for providerItemID, wantCost := range wantCosts {
		seenCost, seenSaleFee := false, false
		seenTaxes := map[profitabilitydomain.InputKind]bool{}
		for _, input := range result.Items {
			if input.ProviderItemID != providerItemID {
				continue
			}
			switch input.Kind {
			case profitabilitydomain.InputKindCost:
				seenCost = amountEquals(input.Amount, wantCost) && input.SourceReference == costSourceReference
			case profitabilitydomain.InputKindSaleFee:
				seenSaleFee = amountEquals(input.Amount, saleFee) && input.SourceReference == "order-1"+saleFeeLineTotalSuffix
			case profitabilitydomain.InputKindTaxICMS, profitabilitydomain.InputKindTaxIPI, profitabilitydomain.InputKindTaxPIS, profitabilitydomain.InputKindTaxCOFINS:
				seenTaxes[input.Kind] = amountEquals(input.Amount, wantTaxes[input.Kind]) && input.SourceReference == string(input.Kind[4:])+taxLineTotalSuffix
			}
		}
		if !seenCost || !seenSaleFee {
			t.Fatalf("item %s cost/sale_fee scopes = cost:%v sale_fee:%v", providerItemID, seenCost, seenSaleFee)
		}
		for kind := range wantTaxes {
			if !seenTaxes[kind] {
				t.Fatalf("item %s tax %s was not preserved as a line total", providerItemID, kind)
			}
		}
	}
}

func TestImportMarginInputsPreservesUnknownCost(t *testing.T) {
	now := time.Date(2026, 7, 12, 12, 0, 0, 0, time.UTC)
	service := NewService(ServiceConfig{
		Orders: stubOrderReader{orders: []profitabilitydomain.OrderFact{{
			InstallationID: "inst-1", ProviderOrderID: "order-1", ProviderUpdatedAt: &now, FetchedAt: now,
			Items: []profitabilitydomain.OrderItemFact{{ProviderItemID: "quantity-7", Quantity: 7, LinkQuality: profitabilitydomain.OrderLinkResolved, InternalProductID: intPtr(7)}},
		}}},
		Internal: stubInternalFactReader{
			cost: internalreaddomain.CostAsOf{AmountScope: internalreaddomain.CostAmountScopePerUnit, QualityFlags: []internalreaddomain.QualityFlag{internalreaddomain.QualityMissingCost}},
		},
		Inputs: &stubInputStore{},
		Now:    func() time.Time { return now },
	})

	result, err := service.ImportMarginInputs(context.Background(), ImportMarginInputsInput{InstallationID: "inst-1"})
	if err != nil {
		t.Fatalf("ImportMarginInputs() error = %v", err)
	}
	for _, input := range result.Items {
		if input.Kind == profitabilitydomain.InputKindCost {
			if input.Amount != nil || input.Quality != profitabilitydomain.InputQualityMissing || input.QualityReason != "missing_cost" || input.SourceReference != costSourceReference {
				t.Fatalf("unknown cost = amount:%v quality:%s reason:%q reference:%q", input.Amount, input.Quality, input.QualityReason, input.SourceReference)
			}
			return
		}
	}
	t.Fatal("expected a cost input")
}

func TestImportMarginInputsDoesNotReadInternalFactsForNonResolvedLinkWithProductID(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	productID := 20303
	tests := []struct {
		name        string
		linkQuality profitabilitydomain.OrderLinkQuality
		wantQuality profitabilitydomain.InputQuality
	}{
		{name: "rejected", linkQuality: profitabilitydomain.OrderLinkRejected, wantQuality: profitabilitydomain.InputQualityRejectedLink},
		{name: "conflict", linkQuality: profitabilitydomain.OrderLinkConflict, wantQuality: profitabilitydomain.InputQualityConflictLink},
		{name: "unresolved", linkQuality: profitabilitydomain.OrderLinkUnresolved, wantQuality: profitabilitydomain.InputQualityUnresolvedLink},
		{name: "missing", linkQuality: profitabilitydomain.OrderLinkMissing, wantQuality: profitabilitydomain.InputQualityMissing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			internal := &recordingInternalFactReader{}
			service := NewService(ServiceConfig{
				Orders: stubOrderReader{orders: []profitabilitydomain.OrderFact{{
					InstallationID:  "inst-1",
					ProviderOrderID: "order-1",
					FetchedAt:       now,
					Items: []profitabilitydomain.OrderItemFact{{
						ProviderItemID:    "MLB4834373620",
						LinkQuality:       tt.linkQuality,
						InternalProductID: &productID,
					}},
				}}},
				Internal: internal,
				Inputs:   &stubInputStore{},
				Now:      func() time.Time { return now },
			})

			result, err := service.ImportMarginInputs(context.Background(), ImportMarginInputsInput{InstallationID: "inst-1"})
			if err != nil {
				t.Fatalf("ImportMarginInputs() error = %v", err)
			}
			if internal.costCalls != 0 || internal.taxCalls != 0 {
				t.Fatalf("internal_read calls = cost:%d tax:%d, want 0 for %s link", internal.costCalls, internal.taxCalls, tt.linkQuality)
			}

			blockedKinds := map[profitabilitydomain.InputKind]bool{
				profitabilitydomain.InputKindCost:      false,
				profitabilitydomain.InputKindTaxICMS:   false,
				profitabilitydomain.InputKindTaxIPI:    false,
				profitabilitydomain.InputKindTaxPIS:    false,
				profitabilitydomain.InputKindTaxCOFINS: false,
			}
			for _, input := range result.Items {
				if _, blocked := blockedKinds[input.Kind]; !blocked {
					continue
				}
				blockedKinds[input.Kind] = true
				if input.Amount != nil || input.Quality != tt.wantQuality || input.QualityReason != "internal product link is not resolved" {
					t.Fatalf("%s input = amount:%v quality:%s reason:%q, want nil/%s/link-blocked", input.Kind, input.Amount, input.Quality, input.QualityReason, tt.wantQuality)
				}
			}
			for kind, found := range blockedKinds {
				if !found {
					t.Fatalf("missing blocked %s input", kind)
				}
			}
		})
	}
}

func TestCreateManualAdjustmentRequiresReasonAndPersistsAudit(t *testing.T) {
	store := &stubAdjustmentStore{}
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	service := NewService(ServiceConfig{
		Inputs:      storeInputNoop{},
		Adjustments: store,
		Snapshots:   storeSnapshotNoop{},
		Now:         func() time.Time { return now },
	})

	_, err := service.CreateManualAdjustment(context.Background(), CreateManualAdjustmentInput{
		InstallationID:  "inst-1",
		ProviderOrderID: "o1",
		Amount:          12.5,
	})
	if err == nil {
		t.Fatal("expected reason validation error")
	}

	adj, err := service.CreateManualAdjustment(context.Background(), CreateManualAdjustmentInput{
		InstallationID:  "inst-1",
		ProviderOrderID: "o1",
		Scope:           profitabilitydomain.InputScopeOrder,
		Category:        profitabilitydomain.ManualAdjustmentFreight,
		Amount:          12.5,
		Currency:        "BRL",
		Reason:          "manual freight correction",
		IdempotencyKey:  "initial-adjustment",
		Actor:           profitabilitydomain.ActorMetadata{ActorType: "operator", ActorID: "leandro"},
	})
	if err != nil {
		t.Fatalf("CreateManualAdjustment() error = %v", err)
	}
	if adj.Category != profitabilitydomain.ManualAdjustmentFreight || len(store.items) != 1 {
		t.Fatalf("adjustment=%+v stored=%d", adj, len(store.items))
	}
}

func TestCreateManualAdjustmentRequiresActorIdentity(t *testing.T) {
	tests := []struct {
		name  string
		actor profitabilitydomain.ActorMetadata
	}{
		{name: "empty actor", actor: profitabilitydomain.ActorMetadata{}},
		{name: "empty actor type", actor: profitabilitydomain.ActorMetadata{ActorID: "operator-1"}},
		{name: "empty actor id", actor: profitabilitydomain.ActorMetadata{ActorType: "operator"}},
		{name: "whitespace actor type", actor: profitabilitydomain.ActorMetadata{ActorType: "  ", ActorID: "operator-1"}},
		{name: "whitespace actor id", actor: profitabilitydomain.ActorMetadata{ActorType: "operator", ActorID: "  "}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &stubAdjustmentStore{}
			service := NewService(ServiceConfig{Adjustments: store})

			_, err := service.CreateManualAdjustment(context.Background(), validAdjustmentInput(tt.actor))
			if err == nil || err.Error() != "PROFITABILITY_ADJUSTMENT_ACTOR_REQUIRED" {
				t.Fatalf("CreateManualAdjustment() error = %v, want actor required", err)
			}
			if len(store.items) != 0 {
				t.Fatalf("stored %d adjustments after validation failure", len(store.items))
			}
		})
	}
}

func TestCreateManualAdjustmentRejectsInvalidScopeAndCategory(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*CreateManualAdjustmentInput)
		wantError string
	}{
		{
			name: "invalid scope",
			mutate: func(input *CreateManualAdjustmentInput) {
				input.Scope = profitabilitydomain.InputScope("listing")
			},
			wantError: "PROFITABILITY_ADJUSTMENT_SCOPE_INVALID",
		},
		{
			name: "invalid category",
			mutate: func(input *CreateManualAdjustmentInput) {
				input.Category = profitabilitydomain.ManualAdjustmentCategory("rebate")
			},
			wantError: "PROFITABILITY_ADJUSTMENT_CATEGORY_INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &stubAdjustmentStore{}
			service := NewService(ServiceConfig{Adjustments: store})
			input := validAdjustmentInput(profitabilitydomain.ActorMetadata{ActorType: "operator", ActorID: "operator-1"})
			tt.mutate(&input)

			_, err := service.CreateManualAdjustment(context.Background(), input)
			if err == nil || err.Error() != tt.wantError {
				t.Fatalf("CreateManualAdjustment() error = %v, want %s", err, tt.wantError)
			}
			if len(store.items) != 0 {
				t.Fatalf("stored %d adjustments after validation failure", len(store.items))
			}
		})
	}
}

func TestCreateManualAdjustmentRejectsIncoherentItemIdentity(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*CreateManualAdjustmentInput)
		wantError string
	}{
		{
			name: "item scope without provider item",
			mutate: func(input *CreateManualAdjustmentInput) {
				input.Scope = profitabilitydomain.InputScopeItem
			},
			wantError: "PROFITABILITY_ADJUSTMENT_ITEM_REQUIRED",
		},
		{
			name: "order scope with provider item",
			mutate: func(input *CreateManualAdjustmentInput) {
				input.ProviderItemID = "MLB-1"
			},
			wantError: "PROFITABILITY_ADJUSTMENT_ORDER_ITEM_FORBIDDEN",
		},
		{
			name: "order scope with provider variation",
			mutate: func(input *CreateManualAdjustmentInput) {
				input.ProviderVariationID = "variation-1"
			},
			wantError: "PROFITABILITY_ADJUSTMENT_ORDER_ITEM_FORBIDDEN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &stubAdjustmentStore{}
			service := NewService(ServiceConfig{Adjustments: store})
			input := validAdjustmentInput(profitabilitydomain.ActorMetadata{ActorType: "operator", ActorID: "operator-1"})
			tt.mutate(&input)

			_, err := service.CreateManualAdjustment(context.Background(), input)
			if err == nil || err.Error() != tt.wantError {
				t.Fatalf("CreateManualAdjustment() error = %v, want %s", err, tt.wantError)
			}
			if len(store.items) != 0 {
				t.Fatalf("stored %d adjustments after validation failure", len(store.items))
			}
		})
	}
}

func TestCreateManualAdjustmentAppendsFreightAndCommissionWithCompleteAudit(t *testing.T) {
	store := &stubAdjustmentStore{}
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	service := NewService(ServiceConfig{
		Adjustments: store,
		Now:         func() time.Time { return now },
	})

	freightInput := validAdjustmentInput(profitabilitydomain.ActorMetadata{
		ActorType: " operator ",
		ActorID:   " operator-1 ",
		ActorName: " Leandro ",
	})
	freightInput.Category = profitabilitydomain.ManualAdjustmentFreight
	freightInput.Amount = -12.50
	freightInput.Reason = " freight correction "
	commissionInput := freightInput
	commissionInput.Category = profitabilitydomain.ManualAdjustmentCommission
	commissionInput.IdempotencyKey = "commission-adjustment"
	commissionInput.Amount = 3.25
	commissionInput.Reason = " commission correction "

	freight, err := service.CreateManualAdjustment(context.Background(), freightInput)
	if err != nil {
		t.Fatalf("CreateManualAdjustment(freight) error = %v", err)
	}
	commission, err := service.CreateManualAdjustment(context.Background(), commissionInput)
	if err != nil {
		t.Fatalf("CreateManualAdjustment(commission) error = %v", err)
	}

	if freight.AdjustmentID == commission.AdjustmentID {
		t.Fatalf("adjustments at the same timestamp reused id %q", freight.AdjustmentID)
	}
	if len(store.items) != 2 {
		t.Fatalf("stored %d adjustments, want 2", len(store.items))
	}
	if freight.Amount != -12.50 || commission.Amount != 3.25 {
		t.Fatalf("signed amounts changed: freight=%v commission=%v", freight.Amount, commission.Amount)
	}
	for _, adjustment := range store.items {
		if adjustment.Actor.ActorType != "operator" || adjustment.Actor.ActorID != "operator-1" || adjustment.Actor.ActorName != "Leandro" {
			t.Fatalf("actor was not trimmed: %+v", adjustment.Actor)
		}
		if adjustment.Reason == "" || adjustment.CreatedAt != now {
			t.Fatalf("incomplete audit: %+v", adjustment)
		}
	}
}

func TestCreateManualAdjustmentPropagatesIDGenerationError(t *testing.T) {
	idErr := errors.New("random source unavailable")
	store := &stubAdjustmentStore{}
	service := NewService(ServiceConfig{
		Adjustments: store,
		NewAdjustmentID: func() (string, error) {
			return "", idErr
		},
	})

	_, err := service.CreateManualAdjustment(context.Background(), validAdjustmentInput(
		profitabilitydomain.ActorMetadata{ActorType: "operator", ActorID: "operator-1"},
	))
	if !errors.Is(err, idErr) {
		t.Fatalf("CreateManualAdjustment() error = %v, want wrapped %v", err, idErr)
	}
	if len(store.items) != 0 {
		t.Fatalf("stored %d adjustments after id generation failure", len(store.items))
	}
}

func validAdjustmentInput(actor profitabilitydomain.ActorMetadata) CreateManualAdjustmentInput {
	return CreateManualAdjustmentInput{
		InstallationID:  "inst-1",
		ProviderOrderID: "order-1",
		Scope:           profitabilitydomain.InputScopeOrder,
		Category:        profitabilitydomain.ManualAdjustmentGeneric,
		Amount:          1,
		Currency:        "BRL",
		Reason:          "manual correction",
		IdempotencyKey:  "adjustment-test-key",
		Actor:           actor,
	}
}

type storeInputNoop struct{}

func (storeInputNoop) ReplaceInputs(context.Context, string, []string, []profitabilitydomain.MarginInput) error {
	return nil
}
func (storeInputNoop) ListInputs(context.Context, string, int) ([]profitabilitydomain.MarginInput, error) {
	return nil, nil
}

type storeSnapshotNoop struct{}

func (storeSnapshotNoop) ReplaceSnapshots(context.Context, string, []string, []profitabilitydomain.ProfitSnapshot) error {
	return nil
}
func (storeSnapshotNoop) ListSnapshots(context.Context, string, int) ([]profitabilitydomain.ProfitSnapshot, error) {
	return nil, nil
}

func intPtr(v int) *int           { return &v }
func floatPtr(v float64) *float64 { return &v }

func TestCalculateSnapshotsHandlesCompleteIncompleteAndNegative(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	snapshotStore := &stubSnapshotStore{}
	service := NewService(ServiceConfig{
		Orders: stubOrderReader{orders: []profitabilitydomain.OrderFact{
			{InstallationID: "inst-1", ProviderOrderID: "o1", RealizationState: profitabilitydomain.OrderRealizationRealized},
			{InstallationID: "inst-1", ProviderOrderID: "o2", RealizationState: profitabilitydomain.OrderRealizationRealized},
			{InstallationID: "inst-1", ProviderOrderID: "o3", RealizationState: profitabilitydomain.OrderRealizationRealized},
			{InstallationID: "inst-1", ProviderOrderID: "o4", RealizationState: profitabilitydomain.OrderRealizationNotRealized},
			{InstallationID: "inst-1", ProviderOrderID: "o5", RealizationState: profitabilitydomain.OrderRealizationUnknown},
		}},
		Inputs: &stubInputStore{items: []profitabilitydomain.MarginInput{
			{InstallationID: "inst-1", ProviderOrderID: "o1", Scope: profitabilitydomain.InputScopeOrder, Kind: profitabilitydomain.InputKindFreight, Amount: floatPtr(1), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o1", Scope: profitabilitydomain.InputScopeOrder, Kind: profitabilitydomain.InputKindCommission, Amount: floatPtr(1), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o1", ProviderItemID: "item-1", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindRevenue, Amount: floatPtr(20), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o1", ProviderItemID: "item-1", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindSaleFee, Amount: floatPtr(2), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o1", ProviderItemID: "item-1", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindCost, Amount: floatPtr(10), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o1", ProviderItemID: "item-1", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindTaxICMS, Amount: floatPtr(1), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o2", ProviderItemID: "item-2", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindRevenue, Amount: floatPtr(10), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o2", ProviderItemID: "item-2", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindSaleFee, Amount: floatPtr(2), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o2", ProviderItemID: "item-2", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindCost, Currency: "BRL", QualityReason: "internal product link is not resolved"},
			{InstallationID: "inst-1", ProviderOrderID: "o2", ProviderItemID: "item-2", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindTaxICMS, Currency: "BRL", QualityReason: "internal product link is not resolved"},
			{InstallationID: "inst-1", ProviderOrderID: "o3", ProviderItemID: "item-3", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindRevenue, Amount: floatPtr(5), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o3", ProviderItemID: "item-3", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindSaleFee, Amount: floatPtr(2), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o3", ProviderItemID: "item-3", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindCost, Amount: floatPtr(6), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o3", ProviderItemID: "item-3", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindTaxICMS, Amount: floatPtr(1), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o4", Scope: profitabilitydomain.InputScopeOrder, Kind: profitabilitydomain.InputKindFreight, Amount: floatPtr(1), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o4", Scope: profitabilitydomain.InputScopeOrder, Kind: profitabilitydomain.InputKindCommission, Amount: floatPtr(2), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o4", ProviderItemID: "item-4", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindRevenue, Amount: floatPtr(20), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o4", ProviderItemID: "item-4", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindSaleFee, Amount: floatPtr(3), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o4", ProviderItemID: "item-4", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindCost, Amount: floatPtr(10), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o4", ProviderItemID: "item-4", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindTaxICMS, Amount: floatPtr(4), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o5", Scope: profitabilitydomain.InputScopeOrder, Kind: profitabilitydomain.InputKindFreight, Amount: floatPtr(5), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o5", Scope: profitabilitydomain.InputScopeOrder, Kind: profitabilitydomain.InputKindCommission, Amount: floatPtr(6), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o5", ProviderItemID: "item-5", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindRevenue, Amount: floatPtr(30), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o5", ProviderItemID: "item-5", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindSaleFee, Amount: floatPtr(7), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o5", ProviderItemID: "item-5", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindCost, Amount: floatPtr(8), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o5", ProviderItemID: "item-5", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindTaxICMS, Amount: floatPtr(9), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o6", ProviderItemID: "item-6", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindRevenue, Amount: floatPtr(11), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o6", ProviderItemID: "item-6", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindSaleFee, Amount: floatPtr(1), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o6", ProviderItemID: "item-6", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindCost, Amount: floatPtr(2), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "o6", ProviderItemID: "item-6", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindTaxICMS, Amount: floatPtr(3), Currency: "BRL"},
		}},
		Adjustments: &stubAdjustmentStore{items: []profitabilitydomain.ManualAdjustment{
			{InstallationID: "inst-1", ProviderOrderID: "o1", Scope: profitabilitydomain.InputScopeItem, ProviderItemID: "item-1", Category: profitabilitydomain.ManualAdjustmentGeneric, Amount: 1},
			{InstallationID: "inst-1", ProviderOrderID: "o4", Scope: profitabilitydomain.InputScopeOrder, Category: profitabilitydomain.ManualAdjustmentGeneric, Amount: -2},
			{InstallationID: "inst-1", ProviderOrderID: "o5", Scope: profitabilitydomain.InputScopeOrder, Category: profitabilitydomain.ManualAdjustmentGeneric, Amount: -3},
		}},
		Snapshots: snapshotStore,
		Now:       func() time.Time { return now },
	})

	result, err := service.CalculateSnapshots(context.Background(), "inst-1", 20)
	if err != nil {
		t.Fatalf("CalculateSnapshots() error = %v", err)
	}
	if result.CalculatedCount == 0 || len(snapshotStore.items) == 0 {
		t.Fatal("expected calculated snapshots")
	}
	var completeFound, incompleteFound, negativeFound, cancelledFound, unknownFound, missingOrderFound bool
	for _, item := range snapshotStore.items {
		if item.ProviderOrderID == "o1" && item.ProviderItemID == "item-1" && item.RealizationState == profitabilitydomain.OrderRealizationRealized && item.Quality == profitabilitydomain.ProfitSnapshotComplete && item.ContributionAmount != nil && *item.ContributionAmount == 8 {
			completeFound = true
		}
		if item.ProviderOrderID == "o2" && item.ProviderItemID == "item-2" && item.Quality == profitabilitydomain.ProfitSnapshotIncomplete {
			incompleteFound = true
		}
		if item.ProviderOrderID == "o3" && item.ProviderItemID == "item-3" && item.Quality == profitabilitydomain.ProfitSnapshotNegativeMargin {
			negativeFound = true
		}
		if item.ProviderOrderID == "o4" && item.Scope == profitabilitydomain.InputScopeOrder {
			cancelledFound = item.RealizationState == profitabilitydomain.OrderRealizationNotRealized &&
				item.Quality == profitabilitydomain.ProfitSnapshotNotRealized &&
				item.ContributionAmount == nil && item.MarginPercent == nil &&
				amountEquals(item.RevenueAmount, 20) && amountEquals(item.SaleFeeAmount, 3) &&
				amountEquals(item.CostAmount, 10) && amountEquals(item.TaxAmount, 4) &&
				amountEquals(item.FreightAmount, 1) && amountEquals(item.CommissionAmount, 2) &&
				amountEquals(item.AdjustmentAmount, -2) &&
				hasSnapshotFlag(item.Flags, profitabilitydomain.ProfitFlagOrderCancelled) &&
				!hasSnapshotFlag(item.Flags, profitabilitydomain.ProfitFlagNegativeMargin)
		}
		if item.ProviderOrderID == "o5" && item.Scope == profitabilitydomain.InputScopeOrder {
			unknownFound = item.RealizationState == profitabilitydomain.OrderRealizationUnknown &&
				item.Quality == profitabilitydomain.ProfitSnapshotIncomplete &&
				item.ContributionAmount == nil && item.MarginPercent == nil &&
				amountEquals(item.RevenueAmount, 30) && amountEquals(item.SaleFeeAmount, 7) &&
				amountEquals(item.CostAmount, 8) && amountEquals(item.TaxAmount, 9) &&
				amountEquals(item.FreightAmount, 5) && amountEquals(item.CommissionAmount, 6) &&
				amountEquals(item.AdjustmentAmount, -3) &&
				hasSnapshotFlag(item.Flags, profitabilitydomain.ProfitFlagOrderStateUnknown)
		}
		if item.ProviderOrderID == "o6" && item.ProviderItemID == "item-6" {
			missingOrderFound = item.RealizationState == profitabilitydomain.OrderRealizationUnknown &&
				item.Quality == profitabilitydomain.ProfitSnapshotIncomplete &&
				item.ContributionAmount == nil && item.MarginPercent == nil &&
				hasSnapshotFlag(item.Flags, profitabilitydomain.ProfitFlagOrderStateUnknown)
		}
	}
	if !completeFound || !incompleteFound || !negativeFound || !cancelledFound || !unknownFound || !missingOrderFound {
		t.Fatalf("expected complete=%v incomplete=%v negative=%v cancelled=%v unknown=%v missing_order=%v", completeFound, incompleteFound, negativeFound, cancelledFound, unknownFound, missingOrderFound)
	}
}

func TestCalculateSnapshotsLimitsOnlyResponseAndReplacesCompleteDeterministicSet(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	store := &stubSnapshotStore{items: []profitabilitydomain.ProfitSnapshot{
		{
			InstallationID:     "inst-1",
			ProviderOrderID:    "order-b",
			Scope:              profitabilitydomain.InputScopeOrder,
			RealizationState:   profitabilitydomain.OrderRealizationRealized,
			Quality:            profitabilitydomain.ProfitSnapshotComplete,
			ContributionAmount: floatPtr(9),
		},
	}}
	service := NewService(ServiceConfig{
		Orders: stubOrderReader{orders: []profitabilitydomain.OrderFact{
			{InstallationID: "inst-1", ProviderOrderID: "order-b", RealizationState: profitabilitydomain.OrderRealizationNotRealized},
			{InstallationID: "inst-1", ProviderOrderID: "order-a", RealizationState: profitabilitydomain.OrderRealizationRealized},
		}},
		Inputs: &stubInputStore{items: []profitabilitydomain.MarginInput{
			{InstallationID: "inst-1", ProviderOrderID: "order-b", ProviderItemID: "item-b", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindRevenue, Amount: floatPtr(20), Currency: "BRL"},
			{InstallationID: "inst-1", ProviderOrderID: "order-a", ProviderItemID: "item-a", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindRevenue, Amount: floatPtr(10), Currency: "BRL"},
		}},
		Adjustments: &stubAdjustmentStore{},
		Snapshots:   store,
		Now:         func() time.Time { return now },
	})

	result, err := service.CalculateSnapshots(context.Background(), "inst-1", 1)
	if err != nil {
		t.Fatalf("CalculateSnapshots() error = %v", err)
	}
	if result.CalculatedCount != 4 || len(store.items) != 4 || len(result.Items) != 1 {
		t.Fatalf("counts = calculated:%d persisted:%d returned:%d, want 4, 4, 1", result.CalculatedCount, len(store.items), len(result.Items))
	}
	if got := store.replacedOrderIDs; len(got) != 2 || got[0] != "order-a" || got[1] != "order-b" {
		t.Fatalf("replaced order IDs = %v, want [order-a order-b]", got)
	}
	wantOrder := []string{"order-a", "order-a", "order-b", "order-b"}
	wantScope := []profitabilitydomain.InputScope{profitabilitydomain.InputScopeOrder, profitabilitydomain.InputScopeItem, profitabilitydomain.InputScopeOrder, profitabilitydomain.InputScopeItem}
	for index, snapshot := range store.items {
		if snapshot.ProviderOrderID != wantOrder[index] || snapshot.Scope != wantScope[index] {
			t.Fatalf("snapshot[%d] = %s/%s, want %s/%s", index, snapshot.ProviderOrderID, snapshot.Scope, wantOrder[index], wantScope[index])
		}
	}
	cancelled := store.items[2]
	if cancelled.RealizationState != profitabilitydomain.OrderRealizationNotRealized || cancelled.Quality != profitabilitydomain.ProfitSnapshotNotRealized {
		t.Fatalf("cancelled snapshot = realization:%s quality:%s, want not_realized", cancelled.RealizationState, cancelled.Quality)
	}
}

func TestBuildProfitSnapshotsKeepsPartialTaxIncompleteAtItemAndOrder(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	snapshots := buildProfitSnapshots([]profitabilitydomain.MarginInput{
		{InstallationID: "inst-1", ProviderOrderID: "order-1", ProviderItemID: "item-1", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindRevenue, Amount: floatPtr(100), Currency: "BRL"},
		{InstallationID: "inst-1", ProviderOrderID: "order-1", ProviderItemID: "item-1", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindSaleFee, Amount: floatPtr(10), Currency: "BRL"},
		{InstallationID: "inst-1", ProviderOrderID: "order-1", ProviderItemID: "item-1", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindCost, Amount: floatPtr(40), Currency: "BRL"},
		{InstallationID: "inst-1", ProviderOrderID: "order-1", ProviderItemID: "item-1", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindTaxICMS, Amount: floatPtr(8), Currency: "BRL", Quality: profitabilitydomain.InputQualityComplete},
		{InstallationID: "inst-1", ProviderOrderID: "order-1", ProviderItemID: "item-1", Scope: profitabilitydomain.InputScopeItem, Kind: profitabilitydomain.InputKindTaxPIS, Currency: "BRL", Quality: profitabilitydomain.InputQualityMissing},
	}, nil, map[string]profitabilitydomain.OrderRealizationState{"order-1": profitabilitydomain.OrderRealizationRealized}, now)

	var item, order *profitabilitydomain.ProfitSnapshot
	for index := range snapshots {
		snapshot := &snapshots[index]
		if snapshot.Scope == profitabilitydomain.InputScopeItem {
			item = snapshot
		} else {
			order = snapshot
		}
	}
	for _, snapshot := range []*profitabilitydomain.ProfitSnapshot{item, order} {
		if snapshot == nil || snapshot.Quality != profitabilitydomain.ProfitSnapshotIncomplete || snapshot.ContributionAmount != nil || snapshot.MarginPercent != nil || !hasSnapshotFlag(snapshot.Flags, profitabilitydomain.ProfitFlagMissingTax) {
			t.Fatalf("partial-tax snapshot = %+v, want incomplete missing_tax with no realized math", snapshot)
		}
	}
	if !amountEquals(item.TaxAmount, 8) || !amountEquals(order.TaxAmount, 8) {
		t.Fatalf("known tax component was not retained: item=%v order=%v", item.TaxAmount, order.TaxAmount)
	}
}

func TestBuildProfitSnapshotsRoutesCostAdjustmentByScope(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	snapshots := buildProfitSnapshots(nil, []profitabilitydomain.ManualAdjustment{
		{InstallationID: "inst-1", ProviderOrderID: "order-1", Scope: profitabilitydomain.InputScopeOrder, Category: profitabilitydomain.ManualAdjustmentCost, Amount: 7, Currency: "BRL"},
		{InstallationID: "inst-1", ProviderOrderID: "order-1", ProviderItemID: "item-1", Scope: profitabilitydomain.InputScopeItem, Category: profitabilitydomain.ManualAdjustmentCost, Amount: 3, Currency: "BRL"},
	}, map[string]profitabilitydomain.OrderRealizationState{"order-1": profitabilitydomain.OrderRealizationRealized}, now)

	for _, snapshot := range snapshots {
		if snapshot.Scope == profitabilitydomain.InputScopeOrder && (!amountEquals(snapshot.AdjustmentAmount, 7) || !amountEquals(snapshot.CostAmount, 3)) {
			t.Fatalf("order cost adjustment = %+v, want order adjustment plus aggregated item cost", snapshot)
		}
		if snapshot.Scope == profitabilitydomain.InputScopeItem && (!amountEquals(snapshot.CostAmount, 3) || snapshot.AdjustmentAmount != nil) {
			t.Fatalf("item cost adjustment = %+v, want cost only", snapshot)
		}
	}
}

func TestCreateManualAdjustmentRejectsMissingIdempotencyKey(t *testing.T) {
	store := &stubAdjustmentStore{}
	service := NewService(ServiceConfig{Adjustments: store})
	input := validAdjustmentInput(profitabilitydomain.ActorMetadata{ActorType: "operator", ActorID: "operator-1"})
	input.IdempotencyKey = ""
	_, err := service.CreateManualAdjustment(context.Background(), input)
	if err == nil || err.Error() != "PROFITABILITY_ADJUSTMENT_IDEMPOTENCY_KEY_REQUIRED" {
		t.Fatalf("CreateManualAdjustment() error = %v, want idempotency key required", err)
	}
	if len(store.items) != 0 {
		t.Fatalf("stored %d adjustments after missing idempotency key", len(store.items))
	}
}

func TestCreateManualAdjustmentReturnsOriginalForDuplicateSubmission(t *testing.T) {
	store := &stubAdjustmentStore{}
	service := NewService(ServiceConfig{Adjustments: store, NewAdjustmentID: func() (string, error) { return "adj-original", nil }})
	input := validAdjustmentInput(profitabilitydomain.ActorMetadata{ActorType: "operator", ActorID: "operator-1"})
	input.IdempotencyKey = "submission-1"
	first, err := service.CreateManualAdjustment(context.Background(), input)
	if err != nil {
		t.Fatalf("first CreateManualAdjustment() error = %v", err)
	}
	input.Amount = 99
	second, err := service.CreateManualAdjustment(context.Background(), input)
	if err != nil {
		t.Fatalf("retry CreateManualAdjustment() error = %v", err)
	}
	if len(store.items) != 1 || second != first || store.items[0] != first {
		t.Fatalf("duplicate submission persisted or changed audit: first=%+v second=%+v stored=%+v", first, second, store.items)
	}
}

func amountEquals(amount *float64, want float64) bool {
	return amount != nil && *amount == want
}

func hasSnapshotFlag(flags []profitabilitydomain.ProfitSnapshotFlag, want profitabilitydomain.ProfitSnapshotFlag) bool {
	for _, flag := range flags {
		if flag == want {
			return true
		}
	}
	return false
}
