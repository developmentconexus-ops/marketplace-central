package application

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

type ImportOrdersInput struct {
	InstallationID string
	Limit          int
	Offset         int
	UpdatedAfter   *time.Time
}

// ImportService enumerates provider order ids for an installation and ingests each one through
// the single write path (F-02, IC-06/ADR-024). It no longer builds domain.MarketplaceOrder
// itself — that shape only exists inside IngestService now, built from the richer per-order
// OrderDetail fetch, not from this list-endpoint snapshot. The snapshot's only remaining job is
// discovery: which provider_order_ids exist for this installation.
type ImportService struct {
	source   ports.OrderSource
	ingestor ports.OrderIngestor
	now      func() time.Time
}

type ImportServiceConfig struct {
	Source   ports.OrderSource
	Ingestor ports.OrderIngestor
	Now      func() time.Time
}

func NewImportService(cfg ImportServiceConfig) *ImportService {
	now := cfg.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &ImportService{
		source:   cfg.Source,
		ingestor: cfg.Ingestor,
		now:      now,
	}
}

// Import enumerates ListOrders' snapshot for provider_order_ids, then calls IngestOrder once per
// id. A per-order failure — whether the typed domain.ErrOrderUnavailable (403/404 on the order
// fetch, no writes) or any other real error from a downstream fetch/write — is counted as a
// skip and logged; the run itself never aborts because ONE order failed (IC-06: "erro por item
// no ingest: registra, segue o batch"). Only a failure to enumerate at all (ListOrders itself)
// or a misconfiguration fails Import outright.
func (s *ImportService) Import(ctx context.Context, input ImportOrdersInput) (domain.ImportResult, error) {
	if s.source == nil || s.ingestor == nil {
		return domain.ImportResult{}, errors.New("ORDERS_IMPORT_NOT_CONFIGURED")
	}
	installationID := strings.TrimSpace(input.InstallationID)
	if installationID == "" {
		return domain.ImportResult{}, errors.New("ORDERS_INSTALLATION_REQUIRED")
	}
	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}

	refs, err := s.source.ListOrders(ctx, ports.ListOrdersInput{
		InstallationID: installationID,
		Limit:          limit,
		Offset:         input.Offset,
		UpdatedAfter:   input.UpdatedAfter,
	})
	if err != nil {
		return domain.ImportResult{}, err
	}

	result := domain.ImportResult{InstallationID: installationID}
	for _, ref := range refs {
		result.EnumeratedCount++
		if ref.ProviderUpdatedAt != nil && !ref.ProviderUpdatedAt.IsZero() {
			if result.MaxProviderUpdatedAt == nil || ref.ProviderUpdatedAt.After(*result.MaxProviderUpdatedAt) {
				updatedAt := *ref.ProviderUpdatedAt
				result.MaxProviderUpdatedAt = &updatedAt
			}
		}
		providerOrderID := strings.TrimSpace(ref.ProviderOrderID)
		if providerOrderID == "" {
			continue
		}
		if err := s.ingestor.IngestOrder(ctx, installationID, providerOrderID); err != nil {
			result.SkippedCount++
			slog.Warn("orders.import.skip", "provider_order_id", providerOrderID, "unavailable", errors.Is(err, domain.ErrOrderUnavailable), "err", err)
			continue
		}
		result.ImportedCount++
	}
	return result, nil
}
