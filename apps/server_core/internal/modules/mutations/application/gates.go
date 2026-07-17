package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
)

type PresenceReader func(context.Context, string) (bool, error)

func InsertWithActorGate(ctx context.Context, actor string, insert func(context.Context) error) error {
	if strings.TrimSpace(actor) == "" {
		return gateError(FailureCodeActorRequired, "ator é obrigatório")
	}
	return insert(ctx)
}

func RequireExecute(execute *bool) error {
	if execute == nil || !*execute {
		return gateError(FailureCodeExecuteRequired, "aprovação exige execute=true explícito")
	}
	return nil
}

func RequireResolvedLink(ctx context.Context, protocolType domain.ProtocolType, listingID string, resolved PresenceReader) error {
	switch protocolType {
	case domain.ProtocolTypePriceUpdate, domain.ProtocolTypeStockCorrect, domain.ProtocolTypeListingEdit:
		ok, err := resolved(ctx, listingID)
		if err != nil {
			return fmt.Errorf("read mutation link resolution: %w", err)
		}
		if !ok {
			return gateError(domain.FailureCodeLinkUnresolved, "vínculo do anúncio não resolvido")
		}
	}
	return nil
}

func RequirePolicy(ctx context.Context, protocolType domain.ProtocolType, listingID string, present PresenceReader) error {
	switch protocolType {
	case domain.ProtocolTypePriceUpdate, domain.ProtocolTypeStockCorrect:
		ok, err := present(ctx, listingID)
		if err != nil {
			return fmt.Errorf("read mutation policy: %w", err)
		}
		if !ok {
			return gateError(domain.FailureCodePolicyMissing, "política obrigatória não encontrada")
		}
	}
	return nil
}

func RequireFreshSource(sourceAsOf *time.Time, now time.Time, threshold time.Duration) error {
	if sourceAsOf == nil {
		return gateError(FailureCodeSourceTimeUnavailable, "data da fonte não disponível")
	}
	if sourceAsOf.Before(now.Add(-threshold)) {
		return gateError(domain.FailureCodeStaleSource, "dados de origem desatualizados")
	}
	return nil
}

func WriteThroughGates(
	ctx context.Context,
	idempotencyKey string,
	alreadyApplied func(context.Context, string) (bool, error),
	persistAudit func(context.Context) error,
	write func(context.Context) error,
) (skipped bool, err error) {
	applied, err := alreadyApplied(ctx, idempotencyKey)
	if err != nil {
		return false, fmt.Errorf("read applied idempotency key: %w", err)
	}
	if applied {
		return true, nil
	}
	if err := persistAudit(ctx); err != nil {
		return false, fmt.Errorf("persist mutation before/after audit: %w", err)
	}
	if err := write(ctx); err != nil {
		return false, err
	}
	return false, nil
}
