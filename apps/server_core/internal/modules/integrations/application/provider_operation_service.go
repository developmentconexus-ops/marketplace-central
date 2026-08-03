package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	connectorsports "marketplace-central/apps/server_core/internal/modules/connectors/ports"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

const (
	providerOperationTypeAccountProbe = "account_probe"
	providerOperationTypeListingRead  = "listing_read"
	providerOperationTypeOrderRead    = "order_read"
	providerOperationTypeFeeQuoteRead = "fee_quote_read"
	providerOperationTypeStockRead    = "stock_read"
)

type providerOperationInstallationReader interface {
	Get(ctx context.Context, installationID string) (domain.Installation, bool, error)
}

type ProviderOperationService struct {
	tenantID      string
	installations providerOperationInstallationReader
	capabilities  *connectorsapp.MarketplaceCapabilityService
	catalogMatch  connectorsports.CatalogMatchReader
	stateWriter   providerOperationCapabilityStateWriter
	operations    *OperationService
	now           func() time.Time
}

type providerOperationCapabilityStateWriter interface {
	Upsert(ctx context.Context, states []domain.CapabilityState) error
}

type ProviderOperationServiceConfig struct {
	TenantID         string
	Installations    providerOperationInstallationReader
	Capabilities     *connectorsapp.MarketplaceCapabilityService
	CatalogMatch     connectorsports.CatalogMatchReader
	CapabilityStates providerOperationCapabilityStateWriter
	Operations       *OperationService
	Now              func() time.Time
}

func NewProviderOperationService(cfg ProviderOperationServiceConfig) *ProviderOperationService {
	tenantID := strings.TrimSpace(cfg.TenantID)
	if tenantID == "" {
		tenantID = "tenant_default"
	}
	now := cfg.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &ProviderOperationService{
		tenantID:      tenantID,
		installations: cfg.Installations,
		capabilities:  cfg.Capabilities,
		catalogMatch:  cfg.CatalogMatch,
		stateWriter:   cfg.CapabilityStates,
		operations:    cfg.Operations,
		now:           now,
	}
}

func (s *ProviderOperationService) ProbeAccount(ctx context.Context, installationID string) (connectorsdomain.AccountSnapshot, error) {
	inst, err := s.loadExecutableInstallation(ctx, installationID, domain.RuntimeCapabilityAccountProbe)
	if err != nil {
		return connectorsdomain.AccountSnapshot{}, err
	}
	prober, err := s.capabilities.AccountProber(inst.ProviderCode)
	if err != nil {
		return connectorsdomain.AccountSnapshot{}, err
	}
	ref := s.accountRef(inst)
	startedAt := s.now()
	result, execErr := prober.ProbeAccount(ctx, ref)
	recordErr := s.recordOperation(ctx, inst.InstallationID, providerOperationTypeAccountProbe, startedAt, execErr, map[string]any{
		"provider_account_id":   result.ProviderAccountID,
		"provider_account_name": result.ProviderAccountName,
		"site_id":               result.SiteID,
		"status":                result.Status,
	})
	if execErr != nil {
		return connectorsdomain.AccountSnapshot{}, execErr
	}
	if recordErr != nil {
		return connectorsdomain.AccountSnapshot{}, recordErr
	}
	return result, nil
}

func (s *ProviderOperationService) ListListings(ctx context.Context, installationID string, limit int) ([]connectorsdomain.ListingSnapshot, error) {
	inst, err := s.loadExecutableInstallation(ctx, installationID, domain.RuntimeCapabilityListingRead)
	if err != nil {
		return nil, err
	}
	reader, err := s.capabilities.ListingReader(inst.ProviderCode)
	if err != nil {
		return nil, err
	}
	startedAt := s.now()
	result, execErr := reader.ListListings(ctx, connectorsdomain.ListListingsInput{
		AccountRef: s.accountRef(inst),
		Limit:      limit,
	})
	recordErr := s.recordOperation(ctx, inst.InstallationID, providerOperationTypeListingRead, startedAt, execErr, map[string]any{
		"listing_count": len(result),
		"limit":         limit,
	})
	if execErr != nil {
		return nil, execErr
	}
	if recordErr != nil {
		return nil, recordErr
	}
	return result, nil
}

// ListOrders enumera pedidos da instalação. O input do provider chega inteiro:
// o service resolve a conta (AccountRef) e não toca em mais nada, para que uma
// capacidade nova do provider (janela, status, ordenação) não exija outra
// mudança de assinatura aqui.
func (s *ProviderOperationService) ListOrders(ctx context.Context, input connectorsdomain.ListOrdersInput, installationID string) ([]connectorsdomain.OrderSnapshot, error) {
	inst, err := s.loadExecutableInstallation(ctx, installationID, domain.RuntimeCapabilityOrderRead)
	if err != nil {
		return nil, err
	}
	reader, err := s.capabilities.OrderReader(inst.ProviderCode)
	if err != nil {
		return nil, err
	}
	input.AccountRef = s.accountRef(inst)
	startedAt := s.now()
	result, execErr := reader.ListOrders(ctx, input)
	recordErr := s.recordOperation(ctx, inst.InstallationID, providerOperationTypeOrderRead, startedAt, execErr, map[string]any{
		"order_count": len(result),
		"limit":       input.Limit,
	})
	if execErr != nil {
		return nil, execErr
	}
	if recordErr != nil {
		return nil, recordErr
	}
	return result, nil
}

func (s *ProviderOperationService) ReadFeeQuote(ctx context.Context, installationID string, input connectorsdomain.FeeQuoteInput) (connectorsdomain.FeeQuoteSnapshot, error) {
	inst, err := s.loadExecutableInstallation(ctx, installationID, domain.RuntimeCapabilityFeeQuoteRead)
	if err != nil {
		return connectorsdomain.FeeQuoteSnapshot{}, err
	}
	reader, err := s.capabilities.FeeQuoteReader(inst.ProviderCode)
	if err != nil {
		return connectorsdomain.FeeQuoteSnapshot{}, err
	}
	input.AccountRef = s.accountRef(inst)
	startedAt := s.now()
	result, execErr := reader.ReadFeeQuote(ctx, input)
	recordErr := s.recordOperation(ctx, inst.InstallationID, providerOperationTypeFeeQuoteRead, startedAt, execErr, map[string]any{
		"site_id":         result.SiteID,
		"category_id":     result.CategoryID,
		"listing_type_id": result.ListingTypeID,
		"price_amount":    result.PriceAmount,
		"currency_id":     result.CurrencyID,
	})
	if execErr != nil {
		return connectorsdomain.FeeQuoteSnapshot{}, execErr
	}
	if recordErr != nil {
		return connectorsdomain.FeeQuoteSnapshot{}, recordErr
	}
	return result, nil
}

func (s *ProviderOperationService) ReadStock(ctx context.Context, installationID string, ref connectorsdomain.ProviderListingRef) (connectorsdomain.StockSnapshot, error) {
	inst, err := s.loadExecutableInstallation(ctx, installationID, domain.RuntimeCapabilityStockRead)
	if err != nil {
		return connectorsdomain.StockSnapshot{}, err
	}
	reader, err := s.capabilities.StockReader(inst.ProviderCode)
	if err != nil {
		return connectorsdomain.StockSnapshot{}, err
	}
	ref.AccountRef = s.accountRef(inst)
	startedAt := s.now()
	result, execErr := reader.ReadStock(ctx, ref)
	recordErr := s.recordOperation(ctx, inst.InstallationID, providerOperationTypeStockRead, startedAt, execErr, map[string]any{
		"provider_item_id":      result.ProviderItemID,
		"provider_variation_id": result.ProviderVariationID,
		"available_quantity":    result.AvailableQuantity,
		"provider_status":       result.ProviderStatus,
		"scope":                 result.Scope,
	})
	if execErr != nil {
		return connectorsdomain.StockSnapshot{}, execErr
	}
	if recordErr != nil {
		return connectorsdomain.StockSnapshot{}, recordErr
	}
	return result, nil
}

func (s *ProviderOperationService) accountRef(inst domain.Installation) connectorsdomain.ProviderAccountRef {
	return connectorsdomain.ProviderAccountRef{
		TenantID:          s.tenantID,
		InstallationID:    inst.InstallationID,
		ProviderCode:      inst.ProviderCode,
		ProviderAccountID: inst.ExternalAccountID,
	}
}

func (s *ProviderOperationService) loadExecutableInstallation(ctx context.Context, installationID string, capability domain.RuntimeCapabilityCode) (domain.Installation, error) {
	if s.installations == nil || s.capabilities == nil || s.operations == nil {
		return domain.Installation{}, errors.New("INTEGRATIONS_PROVIDER_OPERATION_NOT_CONFIGURED")
	}
	inst, found, err := s.installations.Get(ctx, strings.TrimSpace(installationID))
	if err != nil {
		return domain.Installation{}, err
	}
	if !found {
		return domain.Installation{}, domain.ErrInstallationNotFound
	}
	for _, runtimeCapability := range inst.RuntimeCapabilities {
		if runtimeCapability.Code == capability && runtimeCapability.Runnable() {
			return inst, nil
		}
	}
	return domain.Installation{}, errors.New("INTEGRATIONS_CAPABILITY_UNAVAILABLE")
}

func (s *ProviderOperationService) recordOperation(ctx context.Context, installationID, operationType string, startedAt time.Time, execErr error, evidence map[string]any) error {
	if s.operations == nil {
		return nil
	}
	completedAt := s.now()
	record := RecordOperationInput{
		OperationRunID:   operationRunID(operationType, completedAt),
		InstallationID:   installationID,
		OperationType:    operationType,
		Status:           domain.OperationRunStatusSucceeded,
		ResultCode:       "INTEGRATIONS_OPERATION_SUCCEEDED",
		AttemptCount:     1,
		ActorType:        "system",
		ActorID:          "provider_operation_service",
		ProviderEvidence: evidence,
		DurationMs:       completedAt.Sub(startedAt).Milliseconds(),
		StartedAt:        ptrTime(startedAt),
		CompletedAt:      ptrTime(completedAt),
	}
	if execErr != nil {
		record.Status = domain.OperationRunStatusFailed
		record.ResultCode = "INTEGRATIONS_OPERATION_FAILED"
		record.FailureCode = execErr.Error()
		record.TranslatedErrorCode = string(connectorsdomain.ErrorCodeOf(execErr))
	}
	if _, err := s.operations.Record(ctx, record); err != nil {
		return err
	}
	return s.persistCapabilityState(ctx, installationID, runtimeCapabilityCodeForOperation(operationType), completedAt, execErr)
}

func operationRunID(prefix string, now time.Time) string {
	return fmt.Sprintf("%s_%d", strings.TrimSpace(prefix), now.UTC().UnixNano())
}

func (s *ProviderOperationService) persistCapabilityState(ctx context.Context, installationID string, capability domain.RuntimeCapabilityCode, evaluatedAt time.Time, execErr error) error {
	if s.stateWriter == nil || strings.TrimSpace(installationID) == "" || capability == "" {
		return nil
	}

	state := domain.CapabilityState{
		InstallationID:  installationID,
		CapabilityCode:  string(capability),
		Status:          domain.CapabilityStatusEnabled,
		ReasonCode:      "INTEGRATIONS_OPERATION_SUCCEEDED",
		LastEvaluatedAt: ptrTime(evaluatedAt.UTC()),
	}

	if execErr != nil {
		state = domain.CapabilityState{
			InstallationID:  installationID,
			CapabilityCode:  string(capability),
			Status:          capabilityStatusFromError(execErr),
			ReasonCode:      firstNonEmpty(string(connectorsdomain.ErrorCodeOf(execErr)), execErr.Error(), "INTEGRATIONS_OPERATION_FAILED"),
			LastEvaluatedAt: ptrTime(evaluatedAt.UTC()),
		}
	}

	return s.stateWriter.Upsert(ctx, []domain.CapabilityState{state})
}

func runtimeCapabilityCodeForOperation(operationType string) domain.RuntimeCapabilityCode {
	switch strings.TrimSpace(operationType) {
	case providerOperationTypeAccountProbe:
		return domain.RuntimeCapabilityAccountProbe
	case providerOperationTypeListingRead:
		return domain.RuntimeCapabilityListingRead
	case providerOperationTypeOrderRead:
		return domain.RuntimeCapabilityOrderRead
	case providerOperationTypeFeeQuoteRead:
		return domain.RuntimeCapabilityFeeQuoteRead
	case providerOperationTypeStockRead:
		return domain.RuntimeCapabilityStockRead
	default:
		return ""
	}
}

func capabilityStatusFromError(err error) domain.CapabilityStatus {
	switch {
	case errors.Is(err, domain.ErrReauthAccountMismatch),
		errors.Is(err, domain.ErrRefreshTokenInvalid):
		return domain.CapabilityStatusRequiresReauth
	case connectorsdomain.ErrorCodeOf(err) == connectorsdomain.ErrCodeProviderUnsupportedShape:
		return domain.CapabilityStatusUnsupported
	default:
		return domain.CapabilityStatusDegraded
	}
}
