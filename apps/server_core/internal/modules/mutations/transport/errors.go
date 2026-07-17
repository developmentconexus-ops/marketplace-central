package transport

import (
	"errors"
	"net/http"

	"marketplace-central/apps/server_core/internal/modules/mutations/application"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type requestError struct {
	code    string
	message string
}

func (e *requestError) Error() string { return e.code }

var (
	errInvalidBody          = &requestError{code: "invalid_body", message: "corpo da requisição inválido"}
	errInvalidIntent        = &requestError{code: "invalid_intent", message: "intenção inválida"}
	errInstallationRequired = &requestError{code: "installation_required", message: "installation_id é obrigatório"}
	errTypeNotEnabled       = &requestError{code: "type_not_enabled", message: "este tipo de alteração não está habilitado"}
	errMethodNotAllowed     = &requestError{code: "method_not_allowed", message: "method not allowed"}
)

type mutationErrorEnvelope struct {
	Error mutationError `json:"error"`
}

type mutationError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

func writeMutationError(w http.ResponseWriter, err error) {
	status, code, message := mapMutationError(err)
	httpx.WriteJSON(w, status, mutationErrorEnvelope{Error: mutationError{
		Code: code, Message: message, Details: map[string]any{},
	}})
}

func writeMutationMethodError(w http.ResponseWriter) {
	w.Header().Set("Allow", http.MethodPost)
	writeMutationError(w, errMethodNotAllowed)
}

func mapMutationError(err error) (int, string, string) {
	if err == nil {
		return http.StatusInternalServerError, "internal_error", "internal error"
	}
	var requestErr *requestError
	if errors.As(err, &requestErr) {
		switch requestErr.code {
		case "invalid_body":
			return http.StatusBadRequest, requestErr.code, requestErr.message
		case "invalid_intent":
			return http.StatusBadRequest, requestErr.code, requestErr.message
		case "installation_required":
			return http.StatusBadRequest, requestErr.code, requestErr.message
		case "type_not_enabled":
			return http.StatusUnprocessableEntity, requestErr.code, requestErr.message
		case "method_not_allowed":
			return http.StatusMethodNotAllowed, requestErr.code, requestErr.message
		}
	}
	if errors.Is(err, application.ErrProtocolNotFound) {
		return http.StatusNotFound, "protocol_not_found", "protocolo não encontrado"
	}
	if errors.Is(err, domain.ErrInvalidStateTransition) {
		return http.StatusConflict, domain.InvalidStateCode, "transição de estado não permitida"
	}
	var gate *application.GateError
	if errors.As(err, &gate) {
		switch gate.Code {
		case application.FailureCodeActorRequired,
			domain.FailureCodeTypeNotEnabled,
			application.FailureCodeEmptySelection,
			application.FailureCodeSelectionTooLarge,
			application.FailureCodeExecuteRequired,
			application.FailureCodeSourceTimeUnavailable:
			return http.StatusUnprocessableEntity, string(gate.Code), gate.MessagePT
		case application.FailureCodePreviewStale,
			application.FailureCodeNothingToRetry:
			return http.StatusConflict, string(gate.Code), gate.MessagePT
		}
	}
	return http.StatusInternalServerError, "internal_error", "internal error"
}
