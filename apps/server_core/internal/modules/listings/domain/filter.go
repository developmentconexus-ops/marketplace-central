package domain

import (
	"fmt"
	"strconv"
	"strings"
)

var FilterKeys = []string{"status", "sync_state", "link_state", "exception", "has_exception", "listing_type_code", "product_id"}

func IsFilterKey(key string) bool {
	for _, allowed := range FilterKeys {
		if key == allowed {
			return true
		}
	}
	return false
}

type InvalidFilterError struct{ Key string }

func (e *InvalidFilterError) Error() string { return fmt.Sprintf("invalid value for %s", e.Key) }

func SetFilterValue(filter *ListingFilter, key, value string) error {
	if value == "" || strings.TrimSpace(value) != value || strings.Contains(value, ":") {
		return &InvalidFilterError{Key: key}
	}
	switch key {
	case "status":
		filter.Status = ListingStatus(value)
		if !filter.Status.IsValid() {
			return &InvalidFilterError{Key: key}
		}
	case "sync_state":
		filter.SyncState = ListingSyncState(value)
		if !filter.SyncState.IsValid() {
			return &InvalidFilterError{Key: key}
		}
	case "link_state":
		filter.LinkState = LinkState(value)
		if filter.LinkState != LinkStateUnresolved && filter.LinkState != LinkStateConflict && filter.LinkState != LinkStateResolved && filter.LinkState != LinkStateRejected {
			return &InvalidFilterError{Key: key}
		}
	case "exception":
		filter.Exception = ListingException(value)
		if filter.Exception != ListingExceptionSyncError && filter.Exception != ListingExceptionStale && filter.Exception != ListingExceptionUnlinked && filter.Exception != ListingExceptionBelowMargin &&
			filter.Exception != ListingExceptionSemVinculo && filter.Exception != ListingExceptionAbaixoCusto && filter.Exception != ListingExceptionSemEvidencia {
			return &InvalidFilterError{Key: key}
		}
	case "has_exception":
		parsed, err := strconv.ParseBool(value)
		if err != nil || value != "true" && value != "false" {
			return &InvalidFilterError{Key: key}
		}
		filter.HasException = &parsed
	case "listing_type_code":
		filter.ListingTypeCode = ListingTypeCode(value)
		if _, ok := ListingTypeForCode(filter.ListingTypeCode); !ok {
			return &InvalidFilterError{Key: key}
		}
	case "product_id":
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed < 1 || strconv.FormatInt(parsed, 10) != value {
			return &InvalidFilterError{Key: key}
		}
		filter.ProductID = value
	default:
		return &InvalidFilterError{Key: key}
	}
	return nil
}
