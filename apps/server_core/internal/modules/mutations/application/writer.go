package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type ProtocolWriteValidation struct {
	Actor      string
	Execute    *bool
	SourceAsOf *time.Time
	Now        time.Time
	MaxAge     time.Duration
}

type WriterRouter struct {
	price, stock, listing, resync ports.WriterPort
	linkage                       ports.LinkageWriter
	resolved, policy              PresenceReader
}

type RoutedWriteOutcome struct {
	State   domain.ItemState
	Failure *domain.Failure
}

func NewWriterRouter(price, stock, listing, resync ports.WriterPort, linkage ports.LinkageWriter, resolved, policy PresenceReader) *WriterRouter {
	return &WriterRouter{price: price, stock: stock, listing: listing, resync: resync, linkage: linkage, resolved: resolved, policy: policy}
}

func ValidateProtocol(v ProtocolWriteValidation) error {
	if strings.TrimSpace(v.Actor) == "" {
		return gateError(FailureCodeActorRequired, "ator é obrigatório")
	}
	if err := RequireExecute(v.Execute); err != nil {
		return err
	}
	return RequireFreshSource(v.SourceAsOf, v.Now, v.MaxAge)
}

func (r *WriterRouter) Apply(ctx context.Context, item ports.WriteItem) (out ports.WriteOutcome, err error) {
	routed, err := r.ApplyItem(ctx, item)
	return ports.WriteOutcome{Failure: routed.Failure}, err
}

func (r *WriterRouter) ApplyItem(ctx context.Context, item ports.WriteItem) (routed RoutedWriteOutcome, err error) {
	if item.ProtocolType == domain.ProtocolTypeListingCreate {
		if err := validateListingCreate(item.After); err != nil {
			return RoutedWriteOutcome{}, err
		}
		return routedFailure(domain.FailureCodeTypeNotEnabled, "Este tipo de alteração ainda não está habilitado."), nil
	}
	if r == nil {
		return RoutedWriteOutcome{}, errors.New("mutation writer router is not configured")
	}
	if failure := r.itemGate(ctx, item); failure != nil {
		return RoutedWriteOutcome{State: domain.ItemStateFailed, Failure: failure}, nil
	}
	out, state, err := r.dispatch(ctx, item)
	if err != nil {
		return RoutedWriteOutcome{}, err
	}
	routed.State = state
	if routed.State == "" {
		routed.State = domain.ItemStateApplied
	}
	routed.Failure = out.Failure
	if out.Failure != nil {
		routed.State = domain.ItemStateFailed
	}
	return routed, nil
}

func (r *WriterRouter) itemGate(ctx context.Context, item ports.WriteItem) *domain.Failure {
	if item.ProtocolType == domain.ProtocolTypeListingEdit {
		if mismatch, err := sellerSKUMismatch(item.After); err != nil {
			f := failed(domain.FailureCodeProviderValidation, "A edição do anúncio é inválida.").Failure
			return f
		} else if mismatch {
			f := failed(domain.FailureCodeSKUInvariantViolation, "SELLER_SKU deve ser igual ao CODPROD vinculado.").Failure
			return f
		}
	}
	if r.resolved != nil {
		if err := RequireResolvedLink(ctx, item.ProtocolType, item.ListingID, r.resolved); err != nil {
			return itemFailure(err)
		}
	}
	if r.policy != nil {
		if err := RequirePolicy(ctx, item.ProtocolType, item.ListingID, r.policy); err != nil {
			return itemFailure(err)
		}
	}
	return nil
}

func (r *WriterRouter) dispatch(ctx context.Context, item ports.WriteItem) (ports.WriteOutcome, domain.ItemState, error) {
	var writer ports.WriterPort
	switch item.ProtocolType {
	case domain.ProtocolTypePriceUpdate:
		writer = r.price
	case domain.ProtocolTypeStockCorrect:
		writer = r.stock
	case domain.ProtocolTypeListingPause, domain.ProtocolTypeListingEdit:
		writer = r.listing
	case domain.ProtocolTypeListingResync:
		writer = r.resync
	case domain.ProtocolTypeLinkApply:
		return r.applyLink(ctx, item)
	default:
		return failed(domain.FailureCodeTypeNotEnabled, "Este tipo de alteração não está habilitado."), domain.ItemStateFailed, nil
	}
	if writer == nil {
		return ports.WriteOutcome{}, "", errors.New("mutation intent writer is not configured")
	}
	out, err := writer.Apply(ctx, item)
	return out, "", err
}

func (r *WriterRouter) applyLink(ctx context.Context, item ports.WriteItem) (ports.WriteOutcome, domain.ItemState, error) {
	if r.linkage == nil {
		return ports.WriteOutcome{}, "", errors.New("mutation linkage writer is not configured")
	}
	var payload struct {
		Action                ports.LinkageAction `json:"action"`
		CandidateID           string              `json:"candidate_id"`
		ProviderCode          string              `json:"provider_code"`
		ProviderItemID        string              `json:"provider_item_id"`
		ProviderVariationID   string              `json:"provider_variation_id"`
		InternalProductID     int                 `json:"product_id"`
		InternalProductName   string              `json:"product_name"`
		InternalReferenceCode string              `json:"reference_code"`
		ActorType             string              `json:"actor_type"`
		ActorID               string              `json:"actor_id"`
		ActorName             string              `json:"actor_name"`
		Reason                string              `json:"reason"`
	}
	if err := json.Unmarshal(item.After, &payload); err != nil {
		return ports.WriteOutcome{}, "", fmt.Errorf("decode link mutation: %w", err)
	}
	input := ports.LinkageInput{InstallationID: item.InstallationID, Action: payload.Action, CandidateID: payload.CandidateID, ProviderCode: payload.ProviderCode, ProviderItemID: payload.ProviderItemID, ProviderVariationID: payload.ProviderVariationID, InternalProductID: payload.InternalProductID, InternalProductName: payload.InternalProductName, InternalReferenceCode: payload.InternalReferenceCode, ActorType: payload.ActorType, ActorID: payload.ActorID, ActorName: payload.ActorName, Reason: payload.Reason}
	out, err := r.linkage.ApplyLink(ctx, input)
	if err != nil {
		return ports.WriteOutcome{}, "", err
	}
	return ports.WriteOutcome{Failure: out.Failure}, out.State, nil
}

func itemFailure(err error) *domain.Failure {
	var gate *GateError
	if errors.As(err, &gate) {
		f := gate.Failure()
		return &f
	}
	f := failed(domain.FailureCodeInternal, "Falha interna ao validar alteração.").Failure
	return f
}

func failed(code domain.FailureCode, message string) ports.WriteOutcome {
	return ports.WriteOutcome{Failure: &domain.Failure{Code: code, MessagePT: message, Retryable: code.RetryableByDefault()}}
}

func routedFailure(code domain.FailureCode, message string) RoutedWriteOutcome {
	return RoutedWriteOutcome{State: domain.ItemStateFailed, Failure: failed(code, message).Failure}
}

func sellerSKUMismatch(raw json.RawMessage) (bool, error) {
	var intent struct {
		ProductID  json.Number                      `json:"product_id"`
		Attributes []struct{ ID, ValueName string } `json:"attributes"`
	}
	if err := json.Unmarshal(raw, &intent); err != nil {
		return false, err
	}
	codprod := string(intent.ProductID)
	for _, attribute := range intent.Attributes {
		if strings.EqualFold(strings.TrimSpace(attribute.ID), "SELLER_SKU") {
			return strings.TrimSpace(attribute.ValueName) != codprod, nil
		}
	}
	return false, nil
}

func validateListingCreate(raw json.RawMessage) error {
	var intent struct {
		ProductID       json.Number       `json:"product_id"`
		ListingTypeCode string            `json:"listing_type_code"`
		Price           json.RawMessage   `json:"price"`
		CategoryID      string            `json:"category_id"`
		Attributes      []json.RawMessage `json:"attributes"`
	}
	if err := json.Unmarshal(raw, &intent); err != nil {
		return gateError(domain.FailureCodeProviderValidation, "intenção listing_create inválida")
	}
	id, err := strconv.Atoi(string(intent.ProductID))
	if err != nil || id <= 0 || strings.TrimSpace(intent.ListingTypeCode) == "" || len(intent.Price) == 0 || strings.TrimSpace(intent.CategoryID) == "" || intent.Attributes == nil {
		return gateError(domain.FailureCodeProviderValidation, "intenção listing_create incompleta")
	}
	return nil
}
