package application

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

type ImportOrdersInput struct {
	InstallationID string
	Limit          int
}

type ImportService struct {
	source ports.OrderSource
	links  ports.LinkReader
	store  ports.OrderStore
	now    func() time.Time
}

type ImportServiceConfig struct {
	Source ports.OrderSource
	Links  ports.LinkReader
	Store  ports.OrderStore
	Now    func() time.Time
}

func NewImportService(cfg ImportServiceConfig) *ImportService {
	now := cfg.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &ImportService{
		source: cfg.Source,
		links:  cfg.Links,
		store:  cfg.Store,
		now:    now,
	}
}

func (s *ImportService) Import(ctx context.Context, input ImportOrdersInput) (domain.ImportResult, error) {
	if s.source == nil || s.links == nil || s.store == nil {
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

	snapshots, err := s.source.ListOrders(ctx, installationID, limit)
	if err != nil {
		return domain.ImportResult{}, err
	}

	identities := collectIdentities(installationID, snapshots)
	links, err := s.links.ResolveLinks(ctx, installationID, identities)
	if err != nil {
		return domain.ImportResult{}, err
	}

	orders := normalizeOrders(installationID, snapshots, links, s.now())
	imported, skipped, err := s.store.UpsertOrders(ctx, orders)
	if err != nil {
		return domain.ImportResult{}, err
	}

	return domain.ImportResult{
		InstallationID: installationID,
		ImportedCount:  imported,
		SkippedCount:   skipped,
		Items:          orders,
	}, nil
}

func collectIdentities(installationID string, snapshots []domain.OrderIngestionSnapshot) []domain.ListingIdentity {
	seen := map[string]struct{}{}
	identities := make([]domain.ListingIdentity, 0)
	for _, order := range snapshots {
		for _, item := range order.Items {
			identity := domain.ListingIdentity{
				InstallationID:      installationID,
				ProviderItemID:      strings.TrimSpace(item.ProviderItemID),
				ProviderVariationID: strings.TrimSpace(item.ProviderVariationID),
			}
			key := keyOf(identity)
			if identity.ProviderItemID == "" {
				continue
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			identities = append(identities, identity)
		}
	}
	return identities
}

func normalizeOrders(installationID string, snapshots []domain.OrderIngestionSnapshot, links map[string]domain.ItemLink, now time.Time) []domain.MarketplaceOrder {
	items := make([]domain.MarketplaceOrder, 0, len(snapshots))
	for _, snapshot := range snapshots {
		order := domain.MarketplaceOrder{
			InstallationID:       installationID,
			ProviderCode:         strings.TrimSpace(snapshot.ProviderCode),
			ProviderOrderID:      strings.TrimSpace(snapshot.ProviderOrderID),
			ProviderStatus:       strings.TrimSpace(snapshot.ProviderStatus),
			ProviderStatusDetail: strings.TrimSpace(snapshot.ProviderStatusDetail),
			ProviderCreatedAt:    snapshot.ProviderCreatedAt,
			ProviderClosedAt:     snapshot.ProviderClosedAt,
			ProviderUpdatedAt:    snapshot.ProviderUpdatedAt,
			FetchedAt:            snapshot.FetchedAt.UTC(),
			ShippingID:           strings.TrimSpace(snapshot.ShippingID),
			CancellationDetail:   strings.TrimSpace(snapshot.CancellationDetail),
			Tags:                 trimNonEmpty(snapshot.Tags),
			RawProviderRef:       safeOrderProviderReference(snapshot.ProviderCode, snapshot.ProviderOrderID),
			CreatedAt:            now,
			UpdatedAt:            now,
		}
		for _, payment := range snapshot.Payments {
			order.Payments = append(order.Payments, domain.MarketplaceOrderPayment{
				ProviderPaymentID: strings.TrimSpace(payment.PaymentID),
				ProviderStatus:    strings.TrimSpace(payment.Status),
				TransactionAmount: payment.TransactionAmount,
				TotalPaidAmount:   payment.TotalPaidAmount,
			})
		}
		for _, item := range snapshot.Items {
			identity := domain.ListingIdentity{
				InstallationID:      installationID,
				ProviderItemID:      strings.TrimSpace(item.ProviderItemID),
				ProviderVariationID: strings.TrimSpace(item.ProviderVariationID),
			}
			link := links[keyOf(identity)]
			if link.Quality == "" {
				link.Quality = domain.LinkQualityMissing
			}
			order.Items = append(order.Items, domain.MarketplaceOrderItem{
				ProviderItemID:      identity.ProviderItemID,
				ProviderVariationID: identity.ProviderVariationID,
				SellerSKU:           strings.TrimSpace(item.SellerSKU),
				Title:               strings.TrimSpace(item.Title),
				Quantity:            item.Quantity,
				UnitPrice:           item.UnitPrice,
				SaleFeeAmount:       item.SaleFeeAmount,
				LinkQuality:         link.Quality,
				InternalProductID:   link.InternalProductID,
			})
		}
		items = append(items, order)
	}
	return items
}

func safeOrderProviderReference(providerCode, providerOrderID string) string {
	if strings.TrimSpace(providerCode) != "mercado_livre" {
		return ""
	}
	providerOrderID = strings.TrimSpace(providerOrderID)
	if providerOrderID == "" {
		return ""
	}
	return "/orders/" + url.PathEscape(providerOrderID)
}

func keyOf(identity domain.ListingIdentity) string {
	return identity.InstallationID + "|" + identity.ProviderItemID + "|" + identity.ProviderVariationID
}

func trimNonEmpty(values []string) []string {
	trimmed := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			trimmed = append(trimmed, value)
		}
	}
	return trimmed
}
