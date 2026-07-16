package transport

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/ports"
)

var ErrInvalidFilter = errors.New("invalid_filter")
var ErrInvalidCursor = ports.ErrInvalidRunCursor
var ErrInstallationRequired = errors.New("installation_required")

type InvalidFilterError struct{ Key string }

func (e *InvalidFilterError) Error() string {
	return fmt.Sprintf("invalid_filter: invalid value for %s", e.Key)
}
func (e *InvalidFilterError) Code() string            { return "invalid_filter" }
func (e *InvalidFilterError) Details() map[string]any { return map[string]any{"key": e.Key} }
func (e *InvalidFilterError) Unwrap() error           { return ErrInvalidFilter }

type InvalidRunFilterError = InvalidFilterError

type InvalidCursorError struct{}

func (*InvalidCursorError) Error() string           { return ErrInvalidCursor.Error() }
func (*InvalidCursorError) Code() string            { return "invalid_cursor" }
func (*InvalidCursorError) Details() map[string]any { return map[string]any{} }
func (*InvalidCursorError) Unwrap() error           { return ErrInvalidCursor }

type InvalidRunCursorError = InvalidCursorError

type InstallationRequiredError struct{}

func (*InstallationRequiredError) Error() string { return ErrInstallationRequired.Error() }
func (*InstallationRequiredError) Code() string  { return "installation_required" }
func (*InstallationRequiredError) Details() map[string]any {
	return map[string]any{"key": "installation_id"}
}
func (*InstallationRequiredError) Unwrap() error { return ErrInstallationRequired }

func ParseRunQuery(values url.Values) (ports.RunListQuery, error) {
	allowed := map[string]struct{}{
		"installation_id": {}, "cursor": {}, "limit": {},
		"filter.module": {}, "filter.status": {},
	}
	for key := range values {
		if _, ok := allowed[key]; !ok {
			return ports.RunListQuery{}, invalidRunFilter(strings.TrimPrefix(key, "filter."))
		}
	}

	installationID, ok := runSingle(values, "installation_id")
	if !ok || strings.TrimSpace(installationID) == "" {
		return ports.RunListQuery{}, &InstallationRequiredError{}
	}
	if strings.TrimSpace(installationID) != installationID {
		return ports.RunListQuery{}, invalidRunFilter("installation_id")
	}

	query := ports.RunListQuery{InstallationID: installationID, Limit: 20}
	if raw, present := values["limit"]; present {
		if len(raw) != 1 {
			return ports.RunListQuery{}, invalidRunFilter("limit")
		}
		limit, err := strconv.Atoi(raw[0])
		if err != nil || limit < 1 || limit > 100 {
			return ports.RunListQuery{}, invalidRunFilter("limit")
		}
		query.Limit = limit
	}

	if raw, present := values["cursor"]; present {
		if len(raw) != 1 || raw[0] == "" {
			return ports.RunListQuery{}, &InvalidRunCursorError{}
		}
		cursor, err := ports.DecodeRunCursor(raw[0])
		if err != nil {
			return ports.RunListQuery{}, &InvalidRunCursorError{}
		}
		query.Cursor = cursor
	}

	module, err := runFilterScalar(values, "module")
	if err != nil {
		return ports.RunListQuery{}, err
	}
	query.Module = module

	status, err := runFilterScalar(values, "status")
	if err != nil {
		return ports.RunListQuery{}, err
	}
	if status != "" {
		query.Status = domain.OperationRunStatus(status)
		if !validRunStatus(query.Status) {
			return ports.RunListQuery{}, invalidRunFilter("status")
		}
	}

	return query, nil
}

func runSingle(values url.Values, key string) (string, bool) {
	items, ok := values[key]
	if !ok || len(items) != 1 {
		return "", false
	}
	return items[0], true
}

func runFilterScalar(values url.Values, key string) (string, error) {
	fullKey := "filter." + key
	items, present := values[fullKey]
	if !present {
		return "", nil
	}
	if len(items) != 1 || items[0] == "" || strings.TrimSpace(items[0]) != items[0] {
		return "", invalidRunFilter(key)
	}
	return items[0], nil
}

func invalidRunFilter(key string) error { return &InvalidFilterError{Key: key} }

func validRunStatus(status domain.OperationRunStatus) bool {
	switch status {
	case domain.OperationRunStatusQueued,
		domain.OperationRunStatusRunning,
		domain.OperationRunStatusSucceeded,
		domain.OperationRunStatusFailed,
		domain.OperationRunStatusCancelled:
		return true
	default:
		return false
	}
}
