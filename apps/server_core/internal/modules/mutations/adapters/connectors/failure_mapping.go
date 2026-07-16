package connectors

import (
	"context"
	"errors"
	"net"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
)

func MapFailure(err error) domain.Failure {
	code := domain.FailureCodeInternal
	var networkError net.Error
	if errors.Is(err, context.DeadlineExceeded) || errors.As(err, &networkError) {
		code = domain.FailureCodeProviderUnavailable
	}
	switch connectorsdomain.ErrorCodeOf(err) {
	case connectorsdomain.ErrCodeProviderRateLimited:
		code = domain.FailureCodeProviderRateLimited
	case connectorsdomain.ErrCodeProviderAuth:
		code = domain.FailureCodeProviderAuth
	case connectorsdomain.ErrCodeProviderValidation, connectorsdomain.ErrCodeProviderUnsupportedShape, connectorsdomain.ErrCodeProviderInvalidReference:
		code = domain.FailureCodeProviderValidation
	case connectorsdomain.ErrCodeProviderTransient:
		code = domain.FailureCodeProviderUnavailable
	}
	message := ""
	var provider *connectorsdomain.CapabilityError
	if errors.As(err, &provider) {
		message = provider.Message
	}
	return domain.Failure{Code: code, MessagePT: messagePT(code), MessageProvider: message, Retryable: code == domain.FailureCodeProviderRateLimited || code == domain.FailureCodeProviderUnavailable}
}

func MapRejected(message string, paused bool) domain.Failure {
	code := domain.FailureCodeProviderValidation
	if paused {
		code = domain.FailureCodeListingPausedRemote
	}
	return domain.Failure{Code: code, MessagePT: messagePT(code), MessageProvider: message}
}

func messagePT(code domain.FailureCode) string {
	switch code {
	case domain.FailureCodeProviderRateLimited:
		return "O provedor limitou temporariamente as solicitações."
	case domain.FailureCodeProviderAuth:
		return "A autenticação do provedor foi recusada."
	case domain.FailureCodeProviderUnavailable:
		return "O provedor está temporariamente indisponível."
	case domain.FailureCodeListingPausedRemote:
		return "O anúncio está pausado no provedor."
	case domain.FailureCodeProviderValidation:
		return "O provedor recusou os dados enviados."
	default:
		return "Falha interna ao aplicar alteração."
	}
}
