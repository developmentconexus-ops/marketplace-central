package transport

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/listings/domain"
)

type ListingQuery struct {
	Filter domain.ListingFilter
	Q      string
	Limit  int
}

type InvalidFilterError struct{ Key string }

func (e *InvalidFilterError) Error() string {
	return fmt.Sprintf("invalid_filter: invalid value for %s", e.Key)
}
func (e *InvalidFilterError) Code() string            { return "invalid_filter" }
func (e *InvalidFilterError) Details() map[string]any { return map[string]any{"key": e.Key} }

func invalidFilter(key string) error { return &InvalidFilterError{Key: key} }

func ParseListingQuery(values url.Values) (ListingQuery, error) {
	result := ListingQuery{Limit: 50}
	for key := range values {
		if strings.HasPrefix(key, "filter.") {
			if _, ok := allowedFilterKeys[strings.TrimPrefix(key, "filter.")]; !ok {
				return ListingQuery{}, invalidFilter(strings.TrimPrefix(key, "filter."))
			}
		}
	}
	q, err := scalar(values, "q")
	if err != nil {
		return ListingQuery{}, err
	}
	result.Q = strings.TrimSpace(q)
	limit, err := scalar(values, "limit")
	if err != nil {
		return ListingQuery{}, err
	}
	if limit != "" {
		parsed, parseErr := strconv.Atoi(limit)
		if parseErr != nil || parsed < 1 || parsed > 200 {
			return ListingQuery{}, invalidFilter("limit")
		}
		result.Limit = parsed
	}

	if value, err := filterScalar(values, "status"); err != nil {
		return ListingQuery{}, err
	} else if value != "" {
		result.Filter.Status = domain.ListingStatus(value)
		if !result.Filter.Status.IsValid() {
			return ListingQuery{}, invalidFilter("status")
		}
	}
	if value, err := filterScalar(values, "sync_state"); err != nil {
		return ListingQuery{}, err
	} else if value != "" {
		result.Filter.SyncState = domain.ListingSyncState(value)
		if !result.Filter.SyncState.IsValid() {
			return ListingQuery{}, invalidFilter("sync_state")
		}
	}
	if value, err := filterScalar(values, "link_state"); err != nil {
		return ListingQuery{}, err
	} else if value != "" {
		result.Filter.LinkState = domain.LinkState(value)
		if !validLinkState(result.Filter.LinkState) {
			return ListingQuery{}, invalidFilter("link_state")
		}
	}
	if value, err := filterScalar(values, "exception"); err != nil {
		return ListingQuery{}, err
	} else if value != "" {
		result.Filter.Exception = domain.ListingException(value)
		if !validException(result.Filter.Exception) {
			return ListingQuery{}, invalidFilter("exception")
		}
	}
	if value, err := filterScalar(values, "has_exception"); err != nil {
		return ListingQuery{}, err
	} else if value != "" {
		parsed, parseErr := strconv.ParseBool(value)
		if parseErr != nil || (value != "true" && value != "false") {
			return ListingQuery{}, invalidFilter("has_exception")
		}
		result.Filter.HasException = &parsed
	}
	if value, err := filterScalar(values, "listing_type_code"); err != nil {
		return ListingQuery{}, err
	} else if value != "" {
		if _, ok := domain.ListingTypeForCode(domain.ListingTypeCode(value)); !ok {
			return ListingQuery{}, invalidFilter("listing_type_code")
		}
		result.Filter.ListingTypeCode = domain.ListingTypeCode(value)
	}
	if value, err := filterScalar(values, "product_id"); err != nil {
		return ListingQuery{}, err
	} else if value != "" {
		if !positiveCanonicalInteger(value) {
			return ListingQuery{}, invalidFilter("product_id")
		}
		result.Filter.ProductID = value
	}
	return result, nil
}

var allowedFilterKeys = map[string]struct{}{"status": {}, "sync_state": {}, "link_state": {}, "exception": {}, "has_exception": {}, "listing_type_code": {}, "product_id": {}}

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
	full := "filter." + key
	value, err := scalar(values, full)
	if err != nil {
		return "", invalidFilter(key)
	}
	// A present-but-empty filter param (filter.status=) is malformed, not "absent":
	// aliasing it to the default would let a bad request pass as 200.
	if _, present := values[full]; present && value == "" {
		return "", invalidFilter(key)
	}
	if strings.TrimSpace(value) != value {
		return "", invalidFilter(key)
	}
	return value, nil
}

func validLinkState(value domain.LinkState) bool {
	return value == domain.LinkStateUnresolved || value == domain.LinkStateConflict || value == domain.LinkStateResolved || value == domain.LinkStateRejected
}
func validException(value domain.ListingException) bool {
	return value == domain.ListingExceptionSyncError || value == domain.ListingExceptionStale || value == domain.ListingExceptionUnlinked || value == domain.ListingExceptionBelowMargin
}
func positiveCanonicalInteger(value string) bool {
	parsed, err := strconv.ParseInt(value, 10, 64)
	return err == nil && parsed > 0 && strconv.FormatInt(parsed, 10) == value
}
