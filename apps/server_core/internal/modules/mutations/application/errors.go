package application

import (
	"fmt"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
)

const (
	FailureCodeActorRequired         domain.FailureCode = "actor_required"
	FailureCodeExecuteRequired       domain.FailureCode = "execute_required"
	FailureCodeSourceTimeUnavailable domain.FailureCode = "source_time_unavailable"
)

type GateError struct {
	Code      domain.FailureCode
	MessagePT string
}

func (e *GateError) Error() string {
	return fmt.Sprintf("mutation gate %s: %s", e.Code, e.MessagePT)
}

func (e *GateError) Failure() domain.Failure {
	return domain.Failure{
		Code:      e.Code,
		MessagePT: e.MessagePT,
		Retryable: e.Code.RetryableByDefault(),
	}
}

func gateError(code domain.FailureCode, message string) error {
	return &GateError{Code: code, MessagePT: message}
}
