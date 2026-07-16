package transport

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

var ErrInvalidFilter = errors.New("invalid_filter")
var ErrUnsupportedFilter = errors.New("unsupported_filter")
var ErrInvalidCursor = ports.ErrInvalidCursor

type InvalidFilterError struct{ Key string }

func (e *InvalidFilterError) Error() string {
	return fmt.Sprintf("invalid_filter: invalid value for %s", e.Key)
}
func (e *InvalidFilterError) Code() string            { return "invalid_filter" }
func (e *InvalidFilterError) Details() map[string]any { return map[string]any{"key": e.Key} }
func (e *InvalidFilterError) Unwrap() error           { return ErrInvalidFilter }

type UnsupportedFilterError struct{ Key string }

func (e *UnsupportedFilterError) Error() string {
	return fmt.Sprintf("unsupported_filter: %s", e.Key)
}
func (e *UnsupportedFilterError) Code() string            { return "unsupported_filter" }
func (e *UnsupportedFilterError) Details() map[string]any { return map[string]any{"key": e.Key} }
func (e *UnsupportedFilterError) Unwrap() error           { return ErrUnsupportedFilter }

type InvalidCursorError struct{}

func (*InvalidCursorError) Error() string { return ErrInvalidCursor.Error() }
func (*InvalidCursorError) Unwrap() error { return ErrInvalidCursor }
func (*InvalidCursorError) Code() string  { return "invalid_cursor" }
func (*InvalidCursorError) Details() map[string]any {
	return map[string]any{}
}

func ParseOrderQuery(values url.Values) (ports.OrderListQuery, error) {
	result := ports.OrderListQuery{Limit: 20}
	allowed := map[string]struct{}{
		"installation_id": {}, "limit": {}, "cursor": {}, "status": {},
		"date_from": {}, "date_to": {}, "q": {},
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if key == "fulfillment" || key == "filter.fulfillment" {
			return ports.OrderListQuery{}, unsupportedFilter("fulfillment")
		}
		if _, ok := allowed[key]; !ok {
			if strings.HasPrefix(key, "filter.") {
				return ports.OrderListQuery{}, invalidFilter(strings.TrimPrefix(key, "filter."))
			}
			return ports.OrderListQuery{}, invalidFilter(key)
		}
	}

	installationID, err := scalar(values, "installation_id")
	if err != nil {
		return ports.OrderListQuery{}, err
	}
	if installationID == "" || strings.TrimSpace(installationID) != installationID {
		return ports.OrderListQuery{}, invalidFilter("installation_id")
	}
	result.InstallationID = installationID

	limit, err := scalar(values, "limit")
	if err != nil {
		return ports.OrderListQuery{}, err
	}
	if limit != "" {
		parsed, parseErr := strconv.Atoi(limit)
		if parseErr != nil || parsed < 1 || parsed > 100 {
			return ports.OrderListQuery{}, invalidFilter("limit")
		}
		result.Limit = parsed
	}

	result.Q, err = scalar(values, "q")
	if err != nil {
		return ports.OrderListQuery{}, err
	}
	result.Q = strings.TrimSpace(result.Q)

	if result.Status, err = filterScalar(values, "status"); err != nil {
		return ports.OrderListQuery{}, err
	}
	if result.DateFrom, err = parseDate(values, "date_from"); err != nil {
		return ports.OrderListQuery{}, err
	}
	if result.DateTo, err = parseDate(values, "date_to"); err != nil {
		return ports.OrderListQuery{}, err
	}
	if result.DateFrom != nil && result.DateTo != nil && result.DateFrom.After(*result.DateTo) {
		return ports.OrderListQuery{}, invalidFilter("date_from")
	}

	cursor, err := scalar(values, "cursor")
	if err != nil {
		return ports.OrderListQuery{}, err
	}
	if cursor != "" {
		result.Cursor, err = ports.DecodeOrderCursor(cursor)
		if err != nil {
			return ports.OrderListQuery{}, &InvalidCursorError{}
		}
	}
	return result, nil
}

func invalidFilter(key string) error { return &InvalidFilterError{Key: key} }

func unsupportedFilter(key string) error { return &UnsupportedFilterError{Key: key} }

func scalar(values url.Values, key string) (string, error) {
	items, ok := values[key]
	if !ok {
		return "", nil
	}
	if len(items) != 1 {
		return "", invalidFilter(key)
	}
	return items[0], nil
}

func filterScalar(values url.Values, key string) (string, error) {
	value, err := scalar(values, key)
	if err != nil {
		return "", err
	}
	if _, present := values[key]; present && value == "" {
		return "", invalidFilter(key)
	}
	if strings.TrimSpace(value) != value {
		return "", invalidFilter(key)
	}
	return value, nil
}

func parseDate(values url.Values, key string) (*time.Time, error) {
	value, err := scalar(values, key)
	if err != nil {
		return nil, err
	}
	if value == "" {
		if _, present := values[key]; present {
			return nil, invalidFilter(key)
		}
		return nil, nil
	}
	var parsed time.Time
	if len(value) == len("2006-01-02") {
		parsed, err = time.Parse("2006-01-02", value)
	} else {
		parsed, err = time.Parse(time.RFC3339, value)
	}
	if err != nil {
		return nil, invalidFilter(key)
	}
	return &parsed, nil
}
