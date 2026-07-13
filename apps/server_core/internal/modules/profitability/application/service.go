package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	profitabilitydomain "marketplace-central/apps/server_core/internal/modules/profitability/domain"
	"marketplace-central/apps/server_core/internal/modules/profitability/ports"
)

type ImportMarginInputsInput struct {
	InstallationID string
	Limit          int
}

type CreateManualAdjustmentInput struct {
	InstallationID      string
	IdempotencyKey      string
	ProviderOrderID     string
	ProviderItemID      string
	ProviderVariationID string
	Scope               profitabilitydomain.InputScope
	Category            profitabilitydomain.ManualAdjustmentCategory
	Amount              float64
	Currency            string
	Reason              string
	Actor               profitabilitydomain.ActorMetadata
}

type Service struct {
	orders          ports.OrderReader
	internal        ports.InternalFactReader
	lineage         ports.SankhyaLineageReader
	inputs          ports.MarginInputStore
	adjustments     ports.ManualAdjustmentStore
	snapshots       ports.ProfitSnapshotStore
	now             func() time.Time
	newAdjustmentID func() (string, error)
}

type ServiceConfig struct {
	Orders          ports.OrderReader
	Internal        ports.InternalFactReader
	Lineage         ports.SankhyaLineageReader
	Inputs          ports.MarginInputStore
	Adjustments     ports.ManualAdjustmentStore
	Snapshots       ports.ProfitSnapshotStore
	Now             func() time.Time
	NewAdjustmentID func() (string, error)
}

const (
	costSourceReference    = "CUSSEMICM#per_unit_to_item_line_total"
	saleFeeLineTotalSuffix = "#sale_fee_line_total"
	taxLineTotalSuffix     = "#line_total"
)

func NewService(cfg ServiceConfig) *Service {
	now := cfg.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	newAdjustmentID := cfg.NewAdjustmentID
	if newAdjustmentID == nil {
		newAdjustmentID = generateAdjustmentID
	}
	return &Service{
		orders:          cfg.Orders,
		internal:        cfg.Internal,
		lineage:         cfg.Lineage,
		inputs:          cfg.Inputs,
		adjustments:     cfg.Adjustments,
		snapshots:       cfg.Snapshots,
		now:             now,
		newAdjustmentID: newAdjustmentID,
	}
}

func (s *Service) ImportMarginInputs(ctx context.Context, input ImportMarginInputsInput) (profitabilitydomain.ImportInputsResult, error) {
	if s.orders == nil || s.inputs == nil {
		return profitabilitydomain.ImportInputsResult{}, errors.New("PROFITABILITY_IMPORT_NOT_CONFIGURED")
	}
	installationID := strings.TrimSpace(input.InstallationID)
	if installationID == "" {
		return profitabilitydomain.ImportInputsResult{}, errors.New("PROFITABILITY_INSTALLATION_REQUIRED")
	}
	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}
	orders, err := s.orders.ListOrders(ctx, installationID, limit)
	if err != nil {
		return profitabilitydomain.ImportInputsResult{}, err
	}
	now := s.now()
	items := make([]profitabilitydomain.MarginInput, 0)
	orderIDs := make([]string, 0, len(orders))
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ProviderOrderID)
		items = append(items, s.buildOrderLevelInputs(order, now)...)
		for _, item := range order.Items {
			items = append(items, s.buildItemInputs(ctx, order, item, now)...)
		}
	}
	if err := s.inputs.ReplaceInputs(ctx, installationID, orderIDs, items); err != nil {
		return profitabilitydomain.ImportInputsResult{}, err
	}
	return profitabilitydomain.ImportInputsResult{
		InstallationID: installationID,
		ImportedCount:  len(items),
		Items:          items,
	}, nil
}

func (s *Service) buildOrderLevelInputs(order profitabilitydomain.OrderFact, now time.Time) []profitabilitydomain.MarginInput {
	observedAt := order.ProviderUpdatedAt
	if observedAt == nil {
		observedAt = &order.FetchedAt
	}
	base := profitabilitydomain.MarginInput{
		InstallationID:  order.InstallationID,
		ProviderOrderID: order.ProviderOrderID,
		Scope:           profitabilitydomain.InputScopeOrder,
		Currency:        "BRL",
		ObservedAt:      observedAt,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	return []profitabilitydomain.MarginInput{
		{
			InstallationID:  base.InstallationID,
			ProviderOrderID: base.ProviderOrderID,
			Scope:           base.Scope,
			Kind:            profitabilitydomain.InputKindFreight,
			Currency:        base.Currency,
			SourceSystem:    "profitability",
			SourceReference: "manual_or_future_provider_enrichment",
			ObservedAt:      base.ObservedAt,
			Quality:         profitabilitydomain.InputQualityMissing,
			QualityReason:   "freight not provided by F-02 base model",
			CreatedAt:       base.CreatedAt,
			UpdatedAt:       base.UpdatedAt,
		},
		{
			InstallationID:  base.InstallationID,
			ProviderOrderID: base.ProviderOrderID,
			Scope:           base.Scope,
			Kind:            profitabilitydomain.InputKindCommission,
			Currency:        base.Currency,
			SourceSystem:    "profitability",
			SourceReference: "manual_or_future_provider_enrichment",
			ObservedAt:      base.ObservedAt,
			Quality:         profitabilitydomain.InputQualityMissing,
			QualityReason:   "commission adjustment not provided by default",
			CreatedAt:       base.CreatedAt,
			UpdatedAt:       base.UpdatedAt,
		},
	}
}

func (s *Service) buildItemInputs(ctx context.Context, order profitabilitydomain.OrderFact, item profitabilitydomain.OrderItemFact, now time.Time) []profitabilitydomain.MarginInput {
	observedAt := order.ProviderUpdatedAt
	if observedAt == nil {
		observedAt = &order.FetchedAt
	}
	base := profitabilitydomain.MarginInput{
		InstallationID:      order.InstallationID,
		ProviderOrderID:     order.ProviderOrderID,
		ProviderItemID:      item.ProviderItemID,
		ProviderVariationID: item.ProviderVariationID,
		Scope:               profitabilitydomain.InputScopeItem,
		Currency:            "BRL",
		ObservedAt:          observedAt,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	inputs := []profitabilitydomain.MarginInput{
		{
			InstallationID:      base.InstallationID,
			ProviderOrderID:     base.ProviderOrderID,
			ProviderItemID:      base.ProviderItemID,
			ProviderVariationID: base.ProviderVariationID,
			Scope:               base.Scope,
			Kind:                profitabilitydomain.InputKindRevenue,
			Amount:              revenueAmount(item),
			Currency:            base.Currency,
			SourceSystem:        "orders",
			SourceReference:     order.ProviderOrderID,
			ObservedAt:          base.ObservedAt,
			Quality:             qualityForKnownAmount(revenueAmount(item)),
			QualityReason:       qualityReasonForKnownAmount(revenueAmount(item), "order item unit_price is missing"),
			CreatedAt:           base.CreatedAt,
			UpdatedAt:           base.UpdatedAt,
		},
		{
			InstallationID:      base.InstallationID,
			ProviderOrderID:     base.ProviderOrderID,
			ProviderItemID:      base.ProviderItemID,
			ProviderVariationID: base.ProviderVariationID,
			Scope:               base.Scope,
			Kind:                profitabilitydomain.InputKindSaleFee,
			Amount:              item.SaleFeeAmount,
			Currency:            base.Currency,
			SourceSystem:        "orders",
			SourceReference:     order.ProviderOrderID + saleFeeLineTotalSuffix,
			ObservedAt:          base.ObservedAt,
			Quality:             qualityForKnownAmount(item.SaleFeeAmount),
			QualityReason:       qualityReasonForKnownAmount(item.SaleFeeAmount, "sale fee is missing from order snapshot"),
			CreatedAt:           base.CreatedAt,
			UpdatedAt:           base.UpdatedAt,
		},
	}
	if item.LinkQuality != profitabilitydomain.OrderLinkResolved || item.InternalProductID == nil {
		return append(inputs,
			linkBlockedInput(base, profitabilitydomain.InputKindCost, item.LinkQuality),
			linkBlockedInput(base, profitabilitydomain.InputKindTaxICMS, item.LinkQuality),
			linkBlockedInput(base, profitabilitydomain.InputKindTaxIPI, item.LinkQuality),
			linkBlockedInput(base, profitabilitydomain.InputKindTaxPIS, item.LinkQuality),
			linkBlockedInput(base, profitabilitydomain.InputKindTaxCOFINS, item.LinkQuality),
		)
	}
	if s.internal == nil {
		return append(inputs,
			missingInput(base, profitabilitydomain.InputKindCost, "internal_read is unavailable"),
			missingInput(base, profitabilitydomain.InputKindTaxICMS, "internal_read is unavailable"),
			missingInput(base, profitabilitydomain.InputKindTaxIPI, "internal_read is unavailable"),
			missingInput(base, profitabilitydomain.InputKindTaxPIS, "internal_read is unavailable"),
			missingInput(base, profitabilitydomain.InputKindTaxCOFINS, "internal_read is unavailable"),
		)
	}

	cost, costErr := s.internal.GetCostAsOf(ctx, *item.InternalProductID, effectiveAt(order))
	inputs = append(inputs, mapCostInput(base, cost, item.Quantity, costErr))

	inputs = append(inputs, s.buildExactLineageTaxInputs(ctx, order, item, base)...)
	return inputs
}

type exactTaxRead struct {
	source internalreaddomain.TaxSourceIdentity
	value  internalreaddomain.TaxInputs
	valid  bool
}

func (s *Service) buildExactLineageTaxInputs(ctx context.Context, order profitabilitydomain.OrderFact, item profitabilitydomain.OrderItemFact, base profitabilitydomain.MarginInput) []profitabilitydomain.MarginInput {
	if !hasUniqueStableMPCLine(order, item) {
		return missingTaxInputs(base, "stable MPC line identity is unavailable")
	}
	if s.lineage == nil {
		return missingTaxInputs(base, "confirmed Sankhya lineage is unavailable")
	}
	lineage, err := s.lineage.ResolveCurrentLineage(ctx, order.InstallationID, order.ProviderOrderID, item.MPCLineID)
	if err != nil {
		return missingTaxInputs(base, "confirmed Sankhya lineage could not be read")
	}
	if lineage.State != ports.SankhyaLineagePartial && lineage.State != ports.SankhyaLineageComplete {
		return missingTaxInputs(base, "confirmed Sankhya lineage is "+firstNonEmpty(string(lineage.State), "unknown"))
	}
	sources, valid := canonicalTaxSources(lineage.Descendants)
	if !valid || len(sources) == 0 {
		return missingTaxInputs(base, "confirmed Sankhya descendants are missing or invalid")
	}

	reads := make([]exactTaxRead, 0, len(sources))
	for _, source := range sources {
		value, readErr := s.internal.GetTaxInputs(ctx, *item.InternalProductID, effectiveAt(order), source)
		reads = append(reads, exactTaxRead{
			source: source,
			value:  value,
			valid:  readErr == nil && value.SourceIdentity == source,
		})
	}
	reference := exactTaxReference(sources)
	return []profitabilitydomain.MarginInput{
		aggregateExactTaxInput(base, profitabilitydomain.InputKindTaxICMS, reference+":icms"+taxLineTotalSuffix, lineage.State, reads, func(value internalreaddomain.TaxInputs) *float64 { return value.ICMSAmount }),
		aggregateExactTaxInput(base, profitabilitydomain.InputKindTaxIPI, reference+":ipi"+taxLineTotalSuffix, lineage.State, reads, func(value internalreaddomain.TaxInputs) *float64 { return value.IPIAmount }),
		aggregateExactTaxInput(base, profitabilitydomain.InputKindTaxPIS, reference+":pis"+taxLineTotalSuffix, lineage.State, reads, func(value internalreaddomain.TaxInputs) *float64 { return value.PISAmount }),
		aggregateExactTaxInput(base, profitabilitydomain.InputKindTaxCOFINS, reference+":cofins"+taxLineTotalSuffix, lineage.State, reads, func(value internalreaddomain.TaxInputs) *float64 { return value.COFINSAmount }),
	}
}

func hasUniqueStableMPCLine(order profitabilitydomain.OrderFact, item profitabilitydomain.OrderItemFact) bool {
	if item.ReconciliationState != profitabilitydomain.OrderLineStable || !validMPCLineID(item.MPCLineID) {
		return false
	}
	matches := 0
	for _, candidate := range order.Items {
		if candidate.MPCLineID == item.MPCLineID {
			matches++
		}
	}
	return matches == 1
}

func validMPCLineID(value string) bool {
	if value != strings.TrimSpace(value) || len(value) != 36 || !strings.HasPrefix(value, "mpl_") {
		return false
	}
	_, err := hex.DecodeString(value[4:])
	return err == nil
}

func canonicalTaxSources(descendants []ports.SankhyaTaxSource) ([]internalreaddomain.TaxSourceIdentity, bool) {
	sources := make([]internalreaddomain.TaxSourceIdentity, 0, len(descendants))
	seen := make(map[internalreaddomain.TaxSourceIdentity]struct{}, len(descendants))
	for _, descendant := range descendants {
		documentID := int(descendant.DocumentID)
		if descendant.DocumentID <= 0 || int64(documentID) != descendant.DocumentID || descendant.LineNumber <= 0 {
			return nil, false
		}
		source := internalreaddomain.TaxSourceIdentity{DocumentID: documentID, LineNumber: descendant.LineNumber}
		if _, duplicate := seen[source]; duplicate {
			return nil, false
		}
		seen[source] = struct{}{}
		sources = append(sources, source)
	}
	sort.Slice(sources, func(left, right int) bool {
		if sources[left].DocumentID != sources[right].DocumentID {
			return sources[left].DocumentID < sources[right].DocumentID
		}
		return sources[left].LineNumber < sources[right].LineNumber
	})
	return sources, true
}

func exactTaxReference(sources []internalreaddomain.TaxSourceIdentity) string {
	parts := make([]string, len(sources))
	for index, source := range sources {
		parts[index] = fmt.Sprintf("%d/%d", source.DocumentID, source.LineNumber)
	}
	return "sankhya:top306:" + strings.Join(parts, ",")
}

func aggregateExactTaxInput(base profitabilitydomain.MarginInput, kind profitabilitydomain.InputKind, reference string, lineageState ports.SankhyaLineageState, reads []exactTaxRead, amountOf func(internalreaddomain.TaxInputs) *float64) profitabilitydomain.MarginInput {
	result := missingInput(base, kind, "tax component is unknown for at least one exact TOP 306 descendant")
	result.SourceReference = reference
	known := len(reads) > 0
	completeFacts := true
	total := 0.0
	for _, read := range reads {
		amount := amountOf(read.value)
		if !read.valid || amount == nil {
			known = false
			completeFacts = false
			continue
		}
		total += *amount
		if mapInternalQuality(amount, read.value.QualityFlags) != profitabilitydomain.InputQualityComplete {
			completeFacts = false
		}
		result.SourceSystem = firstNonEmpty(result.SourceSystem, read.value.Source.System, "internal_read")
		result.ObservedAt = firstTime(result.ObservedAt, read.value.Source.ObservedAt)
	}
	if !known {
		return result
	}
	result.Amount = &total
	if lineageState == ports.SankhyaLineageComplete && completeFacts {
		result.Quality = profitabilitydomain.InputQualityComplete
		result.QualityReason = ""
		return result
	}
	result.Quality = profitabilitydomain.InputQualityPartial
	if lineageState == ports.SankhyaLineagePartial {
		result.QualityReason = "confirmed Sankhya lineage is partial"
	} else {
		result.QualityReason = "one or more exact TOP 306 tax facts are incomplete"
	}
	return result
}

func missingTaxInputs(base profitabilitydomain.MarginInput, reason string) []profitabilitydomain.MarginInput {
	return []profitabilitydomain.MarginInput{
		missingInput(base, profitabilitydomain.InputKindTaxICMS, reason),
		missingInput(base, profitabilitydomain.InputKindTaxIPI, reason),
		missingInput(base, profitabilitydomain.InputKindTaxPIS, reason),
		missingInput(base, profitabilitydomain.InputKindTaxCOFINS, reason),
	}
}

func (s *Service) ListInputs(ctx context.Context, installationID string, limit int) ([]profitabilitydomain.MarginInput, error) {
	if s.inputs == nil {
		return nil, errors.New("PROFITABILITY_LIST_NOT_CONFIGURED")
	}
	if strings.TrimSpace(installationID) == "" {
		return nil, errors.New("PROFITABILITY_INSTALLATION_REQUIRED")
	}
	if limit <= 0 {
		limit = 50
	}
	return s.inputs.ListInputs(ctx, installationID, limit)
}

func (s *Service) CreateManualAdjustment(ctx context.Context, input CreateManualAdjustmentInput) (profitabilitydomain.ManualAdjustment, error) {
	if s.adjustments == nil {
		return profitabilitydomain.ManualAdjustment{}, errors.New("PROFITABILITY_ADJUSTMENT_NOT_CONFIGURED")
	}
	installationID := strings.TrimSpace(input.InstallationID)
	providerOrderID := strings.TrimSpace(input.ProviderOrderID)
	providerItemID := strings.TrimSpace(input.ProviderItemID)
	providerVariationID := strings.TrimSpace(input.ProviderVariationID)
	idempotencyKey := strings.TrimSpace(input.IdempotencyKey)
	reason := strings.TrimSpace(input.Reason)
	actor := profitabilitydomain.ActorMetadata{
		ActorType: strings.TrimSpace(input.Actor.ActorType),
		ActorID:   strings.TrimSpace(input.Actor.ActorID),
		ActorName: strings.TrimSpace(input.Actor.ActorName),
	}
	if installationID == "" || providerOrderID == "" {
		return profitabilitydomain.ManualAdjustment{}, errors.New("PROFITABILITY_ADJUSTMENT_SCOPE_REQUIRED")
	}
	if idempotencyKey == "" {
		return profitabilitydomain.ManualAdjustment{}, errors.New("PROFITABILITY_ADJUSTMENT_IDEMPOTENCY_KEY_REQUIRED")
	}
	if reason == "" {
		return profitabilitydomain.ManualAdjustment{}, errors.New("PROFITABILITY_ADJUSTMENT_REASON_REQUIRED")
	}
	if actor.ActorType == "" || actor.ActorID == "" {
		return profitabilitydomain.ManualAdjustment{}, errors.New("PROFITABILITY_ADJUSTMENT_ACTOR_REQUIRED")
	}
	scope := input.Scope
	if scope == "" {
		scope = profitabilitydomain.InputScopeOrder
	}
	if scope != profitabilitydomain.InputScopeOrder && scope != profitabilitydomain.InputScopeItem {
		return profitabilitydomain.ManualAdjustment{}, errors.New("PROFITABILITY_ADJUSTMENT_SCOPE_INVALID")
	}
	category := input.Category
	if category == "" {
		category = profitabilitydomain.ManualAdjustmentGeneric
	}
	if category != profitabilitydomain.ManualAdjustmentFreight &&
		category != profitabilitydomain.ManualAdjustmentCommission &&
		category != profitabilitydomain.ManualAdjustmentCost &&
		category != profitabilitydomain.ManualAdjustmentGeneric {
		return profitabilitydomain.ManualAdjustment{}, errors.New("PROFITABILITY_ADJUSTMENT_CATEGORY_INVALID")
	}
	if scope == profitabilitydomain.InputScopeItem && providerItemID == "" {
		return profitabilitydomain.ManualAdjustment{}, errors.New("PROFITABILITY_ADJUSTMENT_ITEM_REQUIRED")
	}
	if scope == profitabilitydomain.InputScopeOrder && (providerItemID != "" || providerVariationID != "") {
		return profitabilitydomain.ManualAdjustment{}, errors.New("PROFITABILITY_ADJUSTMENT_ORDER_ITEM_FORBIDDEN")
	}
	adjustmentID, err := s.newAdjustmentID()
	if err != nil {
		return profitabilitydomain.ManualAdjustment{}, err
	}
	adjustment := profitabilitydomain.ManualAdjustment{
		AdjustmentID:        adjustmentID,
		IdempotencyKey:      idempotencyKey,
		InstallationID:      installationID,
		ProviderOrderID:     providerOrderID,
		ProviderItemID:      providerItemID,
		ProviderVariationID: providerVariationID,
		Scope:               scope,
		Category:            category,
		Amount:              input.Amount,
		Currency:            firstNonEmpty(strings.TrimSpace(input.Currency), "BRL"),
		Reason:              reason,
		Actor:               actor,
		CreatedAt:           s.now(),
	}
	return s.adjustments.CreateOrGetAdjustment(ctx, adjustment)
}

func generateAdjustmentID() (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return "adj_" + hex.EncodeToString(randomBytes), nil
}

func (s *Service) ListAdjustments(ctx context.Context, installationID string, limit int) ([]profitabilitydomain.ManualAdjustment, error) {
	if s.adjustments == nil {
		return nil, errors.New("PROFITABILITY_ADJUSTMENT_NOT_CONFIGURED")
	}
	if strings.TrimSpace(installationID) == "" {
		return nil, errors.New("PROFITABILITY_INSTALLATION_REQUIRED")
	}
	if limit <= 0 {
		limit = 50
	}
	return s.adjustments.ListAdjustments(ctx, installationID, limit)
}

func (s *Service) CalculateSnapshots(ctx context.Context, installationID string, limit int) (profitabilitydomain.CalculateSnapshotsResult, error) {
	if s.orders == nil || s.inputs == nil || s.adjustments == nil || s.snapshots == nil {
		return profitabilitydomain.CalculateSnapshotsResult{}, errors.New("PROFITABILITY_SNAPSHOT_NOT_CONFIGURED")
	}
	installationID = strings.TrimSpace(installationID)
	if installationID == "" {
		return profitabilitydomain.CalculateSnapshotsResult{}, errors.New("PROFITABILITY_INSTALLATION_REQUIRED")
	}
	if limit <= 0 {
		limit = 50
	}
	orders, err := s.orders.ListOrders(ctx, installationID, 1000)
	if err != nil {
		return profitabilitydomain.CalculateSnapshotsResult{}, err
	}
	realizationByOrderID := make(map[string]profitabilitydomain.OrderRealizationState, len(orders))
	for _, order := range orders {
		realizationByOrderID[order.ProviderOrderID] = normalizeRealizationState(order.RealizationState)
	}
	inputs, err := s.inputs.ListInputs(ctx, installationID, 1000)
	if err != nil {
		return profitabilitydomain.CalculateSnapshotsResult{}, err
	}
	adjustments, err := s.adjustments.ListAdjustments(ctx, installationID, 1000)
	if err != nil {
		return profitabilitydomain.CalculateSnapshotsResult{}, err
	}
	snapshots := buildProfitSnapshots(inputs, adjustments, realizationByOrderID, s.now())
	sort.Slice(snapshots, func(i, j int) bool {
		left, right := snapshots[i], snapshots[j]
		if left.ProviderOrderID != right.ProviderOrderID {
			return left.ProviderOrderID < right.ProviderOrderID
		}
		if left.Scope != right.Scope {
			return left.Scope == profitabilitydomain.InputScopeOrder
		}
		if left.ProviderItemID != right.ProviderItemID {
			return left.ProviderItemID < right.ProviderItemID
		}
		return left.ProviderVariationID < right.ProviderVariationID
	})
	orderSet := map[string]struct{}{}
	orderIDs := make([]string, 0)
	for _, snapshot := range snapshots {
		if _, exists := orderSet[snapshot.ProviderOrderID]; !exists {
			orderSet[snapshot.ProviderOrderID] = struct{}{}
			orderIDs = append(orderIDs, snapshot.ProviderOrderID)
		}
	}
	if err := s.snapshots.ReplaceSnapshots(ctx, installationID, orderIDs, snapshots); err != nil {
		return profitabilitydomain.CalculateSnapshotsResult{}, err
	}
	items := snapshots
	if limit < len(items) {
		items = items[:limit]
	}
	return profitabilitydomain.CalculateSnapshotsResult{
		InstallationID:  installationID,
		CalculatedCount: len(snapshots),
		Items:           items,
	}, nil
}

func (s *Service) ListSnapshots(ctx context.Context, installationID string, limit int) ([]profitabilitydomain.ProfitSnapshot, error) {
	if s.snapshots == nil {
		return nil, errors.New("PROFITABILITY_SNAPSHOT_NOT_CONFIGURED")
	}
	installationID = strings.TrimSpace(installationID)
	if installationID == "" {
		return nil, errors.New("PROFITABILITY_INSTALLATION_REQUIRED")
	}
	if limit <= 0 {
		limit = 50
	}
	return s.snapshots.ListSnapshots(ctx, installationID, limit)
}

func revenueAmount(item profitabilitydomain.OrderItemFact) *float64 {
	if item.UnitPrice == nil {
		return nil
	}
	value := *item.UnitPrice * float64(item.Quantity)
	return &value
}

func qualityForKnownAmount(amount *float64) profitabilitydomain.InputQuality {
	if amount == nil {
		return profitabilitydomain.InputQualityMissing
	}
	return profitabilitydomain.InputQualityComplete
}

func qualityReasonForKnownAmount(amount *float64, missingReason string) string {
	if amount == nil {
		return missingReason
	}
	return ""
}

func linkBlockedInput(base profitabilitydomain.MarginInput, kind profitabilitydomain.InputKind, linkQuality profitabilitydomain.OrderLinkQuality) profitabilitydomain.MarginInput {
	return profitabilitydomain.MarginInput{
		InstallationID:      base.InstallationID,
		ProviderOrderID:     base.ProviderOrderID,
		ProviderItemID:      base.ProviderItemID,
		ProviderVariationID: base.ProviderVariationID,
		Scope:               base.Scope,
		Kind:                kind,
		Currency:            base.Currency,
		SourceSystem:        "internal_read",
		ObservedAt:          base.ObservedAt,
		Quality:             mapLinkQuality(linkQuality),
		QualityReason:       "internal product link is not resolved",
		CreatedAt:           base.CreatedAt,
		UpdatedAt:           base.UpdatedAt,
	}
}

func missingInput(base profitabilitydomain.MarginInput, kind profitabilitydomain.InputKind, reason string) profitabilitydomain.MarginInput {
	return profitabilitydomain.MarginInput{
		InstallationID:      base.InstallationID,
		ProviderOrderID:     base.ProviderOrderID,
		ProviderItemID:      base.ProviderItemID,
		ProviderVariationID: base.ProviderVariationID,
		Scope:               base.Scope,
		Kind:                kind,
		Currency:            base.Currency,
		SourceSystem:        "internal_read",
		ObservedAt:          base.ObservedAt,
		Quality:             profitabilitydomain.InputQualityMissing,
		QualityReason:       reason,
		CreatedAt:           base.CreatedAt,
		UpdatedAt:           base.UpdatedAt,
	}
}

func mapCostInput(base profitabilitydomain.MarginInput, cost internalreaddomain.CostAsOf, quantity int, err error) profitabilitydomain.MarginInput {
	if err != nil {
		return profitabilitydomain.MarginInput{
			InstallationID:      base.InstallationID,
			ProviderOrderID:     base.ProviderOrderID,
			ProviderItemID:      base.ProviderItemID,
			ProviderVariationID: base.ProviderVariationID,
			Scope:               base.Scope,
			Kind:                profitabilitydomain.InputKindCost,
			Currency:            base.Currency,
			SourceSystem:        "internal_read",
			SourceReference:     costSourceReference,
			ObservedAt:          base.ObservedAt,
			Quality:             profitabilitydomain.InputQualityMissing,
			QualityReason:       err.Error(),
			CreatedAt:           base.CreatedAt,
			UpdatedAt:           base.UpdatedAt,
		}
	}
	return profitabilitydomain.MarginInput{
		InstallationID:      base.InstallationID,
		ProviderOrderID:     base.ProviderOrderID,
		ProviderItemID:      base.ProviderItemID,
		ProviderVariationID: base.ProviderVariationID,
		Scope:               base.Scope,
		Kind:                profitabilitydomain.InputKindCost,
		Amount:              extendUnitCost(cost.Amount, quantity),
		Currency:            base.Currency,
		SourceSystem:        firstNonEmpty(cost.Source.System, "internal_read"),
		SourceReference:     costSourceReference,
		ObservedAt:          firstTime(cost.Source.ObservedAt, &cost.EffectiveAt),
		Quality:             mapInternalQuality(cost.Amount, cost.QualityFlags),
		QualityReason:       qualityReason(cost.QualityFlags),
		CreatedAt:           base.CreatedAt,
		UpdatedAt:           base.UpdatedAt,
	}
}

func extendUnitCost(amount *float64, quantity int) *float64 {
	if amount == nil {
		return nil
	}
	lineTotal := *amount * float64(quantity)
	return &lineTotal
}

func effectiveAt(order profitabilitydomain.OrderFact) time.Time {
	if order.ProviderClosedAt != nil {
		return order.ProviderClosedAt.UTC()
	}
	if order.ProviderCreatedAt != nil {
		return order.ProviderCreatedAt.UTC()
	}
	if order.ProviderUpdatedAt != nil {
		return order.ProviderUpdatedAt.UTC()
	}
	return order.FetchedAt.UTC()
}

func mapLinkQuality(linkQuality profitabilitydomain.OrderLinkQuality) profitabilitydomain.InputQuality {
	switch linkQuality {
	case profitabilitydomain.OrderLinkRejected:
		return profitabilitydomain.InputQualityRejectedLink
	case profitabilitydomain.OrderLinkConflict:
		return profitabilitydomain.InputQualityConflictLink
	case profitabilitydomain.OrderLinkUnresolved:
		return profitabilitydomain.InputQualityUnresolvedLink
	default:
		return profitabilitydomain.InputQualityMissing
	}
}

func mapInternalQuality(amount *float64, flags []internalreaddomain.QualityFlag) profitabilitydomain.InputQuality {
	if internalreaddomain.HasQualityFlag(flags, internalreaddomain.QualityStaleSource) {
		return profitabilitydomain.InputQualityStale
	}
	if amount == nil || len(flags) > 0 && !internalreaddomain.HasQualityFlag(flags, internalreaddomain.QualityComplete) {
		return profitabilitydomain.InputQualityMissing
	}
	return profitabilitydomain.InputQualityComplete
}

func qualityReason(flags []internalreaddomain.QualityFlag) string {
	if len(flags) == 0 {
		return ""
	}
	parts := make([]string, 0, len(flags))
	for _, flag := range flags {
		if flag == internalreaddomain.QualityComplete {
			continue
		}
		parts = append(parts, string(flag))
	}
	return strings.Join(parts, ",")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func firstTime(values ...*time.Time) *time.Time {
	for _, value := range values {
		if value != nil {
			t := value.UTC()
			return &t
		}
	}
	return nil
}

type snapshotAccumulator struct {
	snapshot profitabilitydomain.ProfitSnapshot
	flags    map[profitabilitydomain.ProfitSnapshotFlag]struct{}
}

func buildProfitSnapshots(inputs []profitabilitydomain.MarginInput, adjustments []profitabilitydomain.ManualAdjustment, realizationByOrderID map[string]profitabilitydomain.OrderRealizationState, now time.Time) []profitabilitydomain.ProfitSnapshot {
	orderAcc := map[string]*snapshotAccumulator{}
	itemAcc := map[string]*snapshotAccumulator{}

	for _, input := range inputs {
		if input.Scope == profitabilitydomain.InputScopeOrder {
			acc := ensureAccumulator(orderAcc, snapshotKey(input.ProviderOrderID, "", ""), profitabilitydomain.ProfitSnapshot{
				InstallationID:   input.InstallationID,
				ProviderOrderID:  input.ProviderOrderID,
				Scope:            profitabilitydomain.InputScopeOrder,
				Currency:         firstNonEmpty(input.Currency, "BRL"),
				RealizationState: realizationStateForOrder(realizationByOrderID, input.ProviderOrderID),
				CreatedAt:        now,
				UpdatedAt:        now,
			})
			applyInput(acc, input)
			continue
		}
		acc := ensureAccumulator(itemAcc, snapshotKey(input.ProviderOrderID, input.ProviderItemID, input.ProviderVariationID), profitabilitydomain.ProfitSnapshot{
			InstallationID:      input.InstallationID,
			ProviderOrderID:     input.ProviderOrderID,
			ProviderItemID:      input.ProviderItemID,
			ProviderVariationID: input.ProviderVariationID,
			Scope:               profitabilitydomain.InputScopeItem,
			Currency:            firstNonEmpty(input.Currency, "BRL"),
			RealizationState:    realizationStateForOrder(realizationByOrderID, input.ProviderOrderID),
			CreatedAt:           now,
			UpdatedAt:           now,
		})
		applyInput(acc, input)
	}

	for _, adjustment := range adjustments {
		if adjustment.Scope == profitabilitydomain.InputScopeOrder {
			acc := ensureAccumulator(orderAcc, snapshotKey(adjustment.ProviderOrderID, "", ""), profitabilitydomain.ProfitSnapshot{
				InstallationID:   adjustment.InstallationID,
				ProviderOrderID:  adjustment.ProviderOrderID,
				Scope:            profitabilitydomain.InputScopeOrder,
				Currency:         firstNonEmpty(adjustment.Currency, "BRL"),
				RealizationState: realizationStateForOrder(realizationByOrderID, adjustment.ProviderOrderID),
				CreatedAt:        now,
				UpdatedAt:        now,
			})
			applyAdjustment(acc, adjustment)
			continue
		}
		acc := ensureAccumulator(itemAcc, snapshotKey(adjustment.ProviderOrderID, adjustment.ProviderItemID, adjustment.ProviderVariationID), profitabilitydomain.ProfitSnapshot{
			InstallationID:      adjustment.InstallationID,
			ProviderOrderID:     adjustment.ProviderOrderID,
			ProviderItemID:      adjustment.ProviderItemID,
			ProviderVariationID: adjustment.ProviderVariationID,
			Scope:               profitabilitydomain.InputScopeItem,
			Currency:            firstNonEmpty(adjustment.Currency, "BRL"),
			RealizationState:    realizationStateForOrder(realizationByOrderID, adjustment.ProviderOrderID),
			CreatedAt:           now,
			UpdatedAt:           now,
		})
		applyAdjustment(acc, adjustment)
	}

	for _, itemAccValue := range itemAcc {
		finalizeSnapshot(itemAccValue)
		orderKey := snapshotKey(itemAccValue.snapshot.ProviderOrderID, "", "")
		orderAccumulator := ensureAccumulator(orderAcc, orderKey, profitabilitydomain.ProfitSnapshot{
			InstallationID:   itemAccValue.snapshot.InstallationID,
			ProviderOrderID:  itemAccValue.snapshot.ProviderOrderID,
			Scope:            profitabilitydomain.InputScopeOrder,
			Currency:         firstNonEmpty(itemAccValue.snapshot.Currency, "BRL"),
			RealizationState: realizationStateForOrder(realizationByOrderID, itemAccValue.snapshot.ProviderOrderID),
			CreatedAt:        now,
			UpdatedAt:        now,
		})
		aggregateItemIntoOrder(orderAccumulator, itemAccValue.snapshot)
	}

	results := make([]profitabilitydomain.ProfitSnapshot, 0, len(orderAcc)+len(itemAcc))
	for _, acc := range orderAcc {
		finalizeSnapshot(acc)
		results = append(results, acc.snapshot)
	}
	for _, acc := range itemAcc {
		results = append(results, acc.snapshot)
	}
	return results
}

func ensureAccumulator(target map[string]*snapshotAccumulator, key string, snapshot profitabilitydomain.ProfitSnapshot) *snapshotAccumulator {
	if acc, ok := target[key]; ok {
		return acc
	}
	acc := &snapshotAccumulator{snapshot: snapshot, flags: map[profitabilitydomain.ProfitSnapshotFlag]struct{}{}}
	target[key] = acc
	return acc
}

func applyInput(acc *snapshotAccumulator, input profitabilitydomain.MarginInput) {
	switch input.Kind {
	case profitabilitydomain.InputKindRevenue:
		acc.snapshot.RevenueAmount = input.Amount
	case profitabilitydomain.InputKindSaleFee:
		acc.snapshot.SaleFeeAmount = input.Amount
	case profitabilitydomain.InputKindCost:
		acc.snapshot.CostAmount = input.Amount
	case profitabilitydomain.InputKindFreight:
		acc.snapshot.FreightAmount = input.Amount
	case profitabilitydomain.InputKindCommission:
		acc.snapshot.CommissionAmount = input.Amount
	case profitabilitydomain.InputKindTaxICMS, profitabilitydomain.InputKindTaxIPI, profitabilitydomain.InputKindTaxPIS, profitabilitydomain.InputKindTaxCOFINS:
		acc.snapshot.TaxAmount = addOptional(acc.snapshot.TaxAmount, input.Amount)
		if input.Amount == nil {
			acc.flags[profitabilitydomain.ProfitFlagMissingTax] = struct{}{}
		}
	}
	if strings.Contains(input.QualityReason, "link is not resolved") || input.Quality == profitabilitydomain.InputQualityUnresolvedLink || input.Quality == profitabilitydomain.InputQualityRejectedLink || input.Quality == profitabilitydomain.InputQualityConflictLink {
		acc.flags[profitabilitydomain.ProfitFlagMissingLink] = struct{}{}
	}
}

func applyAdjustment(acc *snapshotAccumulator, adjustment profitabilitydomain.ManualAdjustment) {
	switch adjustment.Category {
	case profitabilitydomain.ManualAdjustmentFreight:
		acc.snapshot.FreightAmount = addOptional(acc.snapshot.FreightAmount, &adjustment.Amount)
	case profitabilitydomain.ManualAdjustmentCommission:
		acc.snapshot.CommissionAmount = addOptional(acc.snapshot.CommissionAmount, &adjustment.Amount)
	case profitabilitydomain.ManualAdjustmentCost:
		if adjustment.Scope == profitabilitydomain.InputScopeItem {
			acc.snapshot.CostAmount = addOptional(acc.snapshot.CostAmount, &adjustment.Amount)
			return
		}
		acc.snapshot.AdjustmentAmount = addOptional(acc.snapshot.AdjustmentAmount, &adjustment.Amount)
	default:
		acc.snapshot.AdjustmentAmount = addOptional(acc.snapshot.AdjustmentAmount, &adjustment.Amount)
	}
}

func finalizeSnapshot(acc *snapshotAccumulator) {
	s := &acc.snapshot
	flags := map[profitabilitydomain.ProfitSnapshotFlag]struct{}{}
	if _, ok := acc.flags[profitabilitydomain.ProfitFlagMissingLink]; ok {
		flags[profitabilitydomain.ProfitFlagMissingLink] = struct{}{}
	}
	if _, ok := acc.flags[profitabilitydomain.ProfitFlagMissingTax]; ok {
		flags[profitabilitydomain.ProfitFlagMissingTax] = struct{}{}
	}
	s.ContributionAmount = nil
	s.MarginPercent = nil
	switch s.RealizationState {
	case profitabilitydomain.OrderRealizationNotRealized:
		flags[profitabilitydomain.ProfitFlagOrderCancelled] = struct{}{}
		s.Flags = flattenFlags(flags)
		s.Quality = profitabilitydomain.ProfitSnapshotNotRealized
		return
	case profitabilitydomain.OrderRealizationUnknown:
		flags[profitabilitydomain.ProfitFlagOrderStateUnknown] = struct{}{}
		s.Flags = flattenFlags(flags)
		s.Quality = profitabilitydomain.ProfitSnapshotIncomplete
		return
	}
	requiredMissing := false
	if s.RevenueAmount == nil {
		flags[profitabilitydomain.ProfitFlagMissingRevenue] = struct{}{}
		requiredMissing = true
	}
	if s.SaleFeeAmount == nil {
		flags[profitabilitydomain.ProfitFlagMissingSaleFee] = struct{}{}
		requiredMissing = true
	}
	if s.Scope == profitabilitydomain.InputScopeItem && s.CostAmount == nil {
		flags[profitabilitydomain.ProfitFlagMissingCost] = struct{}{}
		requiredMissing = true
	}
	if (s.Scope == profitabilitydomain.InputScopeItem && s.TaxAmount == nil) || hasAccumulatorFlag(acc, profitabilitydomain.ProfitFlagMissingTax) {
		flags[profitabilitydomain.ProfitFlagMissingTax] = struct{}{}
		requiredMissing = true
	}
	if s.Scope == profitabilitydomain.InputScopeOrder && s.FreightAmount == nil {
		flags[profitabilitydomain.ProfitFlagMissingFreight] = struct{}{}
		requiredMissing = true
	}
	if s.Scope == profitabilitydomain.InputScopeOrder && s.CommissionAmount == nil {
		flags[profitabilitydomain.ProfitFlagMissingCommission] = struct{}{}
		requiredMissing = true
	}

	s.Flags = flattenFlags(flags)
	if !requiredMissing {
		contribution := *s.RevenueAmount - *s.SaleFeeAmount
		if s.CostAmount != nil {
			contribution -= *s.CostAmount
		}
		if s.TaxAmount != nil {
			contribution -= *s.TaxAmount
		}
		if s.FreightAmount != nil {
			contribution -= *s.FreightAmount
		}
		if s.CommissionAmount != nil {
			contribution -= *s.CommissionAmount
		}
		if s.AdjustmentAmount != nil {
			contribution += *s.AdjustmentAmount
		}
		s.ContributionAmount = &contribution
		if *s.RevenueAmount != 0 {
			margin := contribution / *s.RevenueAmount * 100
			s.MarginPercent = &margin
		}
		s.Quality = profitabilitydomain.ProfitSnapshotComplete
		if contribution < 0 {
			flags[profitabilitydomain.ProfitFlagNegativeMargin] = struct{}{}
			s.Flags = flattenFlags(flags)
			s.Quality = profitabilitydomain.ProfitSnapshotNegativeMargin
		}
		return
	}
	s.Quality = profitabilitydomain.ProfitSnapshotIncomplete
}

func hasAccumulatorFlag(acc *snapshotAccumulator, flag profitabilitydomain.ProfitSnapshotFlag) bool {
	_, exists := acc.flags[flag]
	return exists
}

func normalizeRealizationState(state profitabilitydomain.OrderRealizationState) profitabilitydomain.OrderRealizationState {
	if state == profitabilitydomain.OrderRealizationRealized || state == profitabilitydomain.OrderRealizationNotRealized {
		return state
	}
	return profitabilitydomain.OrderRealizationUnknown
}

func realizationStateForOrder(states map[string]profitabilitydomain.OrderRealizationState, providerOrderID string) profitabilitydomain.OrderRealizationState {
	return normalizeRealizationState(states[providerOrderID])
}

func snapshotKey(orderID, itemID, variationID string) string {
	return orderID + "|" + itemID + "|" + variationID
}

func addOptional(current, extra *float64) *float64 {
	if current == nil && extra == nil {
		return nil
	}
	var sum float64
	if current != nil {
		sum += *current
	}
	if extra != nil {
		sum += *extra
	}
	return &sum
}

func flattenFlags(flags map[profitabilitydomain.ProfitSnapshotFlag]struct{}) []profitabilitydomain.ProfitSnapshotFlag {
	if len(flags) == 0 {
		return nil
	}
	items := make([]profitabilitydomain.ProfitSnapshotFlag, 0, len(flags))
	for flag := range flags {
		items = append(items, flag)
	}
	return items
}

func aggregateItemIntoOrder(orderAcc *snapshotAccumulator, item profitabilitydomain.ProfitSnapshot) {
	orderAcc.snapshot.RevenueAmount = addOptional(orderAcc.snapshot.RevenueAmount, item.RevenueAmount)
	orderAcc.snapshot.SaleFeeAmount = addOptional(orderAcc.snapshot.SaleFeeAmount, item.SaleFeeAmount)
	orderAcc.snapshot.CostAmount = addOptional(orderAcc.snapshot.CostAmount, item.CostAmount)
	orderAcc.snapshot.TaxAmount = addOptional(orderAcc.snapshot.TaxAmount, item.TaxAmount)
	orderAcc.snapshot.AdjustmentAmount = addOptional(orderAcc.snapshot.AdjustmentAmount, item.AdjustmentAmount)
	for _, flag := range item.Flags {
		orderAcc.flags[flag] = struct{}{}
	}
}
