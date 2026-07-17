package transport

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var ErrInstallationRequired = errors.New("installation_required")
var ErrInvalidFilter = errors.New("invalid_filter")

type InstallationRequiredError struct{}

func (*InstallationRequiredError) Error() string { return ErrInstallationRequired.Error() }
func (*InstallationRequiredError) Code() string  { return ErrInstallationRequired.Error() }
func (*InstallationRequiredError) Details() map[string]any {
	return map[string]any{"key": "installation_id"}
}
func (*InstallationRequiredError) Unwrap() error { return ErrInstallationRequired }

type InvalidFilterError struct{ Key string }

func (e *InvalidFilterError) Error() string {
	return fmt.Sprintf("invalid_filter: invalid value for %s", e.Key)
}
func (*InvalidFilterError) Code() string { return ErrInvalidFilter.Error() }
func (e *InvalidFilterError) Details() map[string]any {
	return map[string]any{"key": e.Key}
}
func (*InvalidFilterError) Unwrap() error { return ErrInvalidFilter }

type ObservationQuery struct {
	InstallationID string
	ListingIDs     []string
}

type ReferenceQuery struct {
	ProductIDs []string
}

func ParseObservationQuery(values url.Values) (ObservationQuery, error) {
	if err := rejectUnknown(values, map[string]struct{}{
		"installation_id": {},
		"listing_ids":     {},
	}); err != nil {
		return ObservationQuery{}, err
	}

	installationID, ok := singleValue(values, "installation_id")
	if !ok || strings.TrimSpace(installationID) == "" {
		return ObservationQuery{}, &InstallationRequiredError{}
	}
	if strings.TrimSpace(installationID) != installationID {
		return ObservationQuery{}, &InvalidFilterError{Key: "installation_id"}
	}

	listingIDs, err := commaSeparated(values, "listing_ids")
	if err != nil {
		return ObservationQuery{}, err
	}
	return ObservationQuery{InstallationID: installationID, ListingIDs: listingIDs}, nil
}

func ParseReferenceQuery(values url.Values) (ReferenceQuery, error) {
	if err := rejectUnknown(values, map[string]struct{}{"product_ids": {}}); err != nil {
		return ReferenceQuery{}, err
	}
	productIDs, err := commaSeparated(values, "product_ids")
	if err != nil {
		return ReferenceQuery{}, err
	}
	return ReferenceQuery{ProductIDs: productIDs}, nil
}

func rejectUnknown(values url.Values, allowed map[string]struct{}) error {
	for key := range values {
		if _, ok := allowed[key]; !ok {
			return &InvalidFilterError{Key: strings.TrimPrefix(key, "filter.")}
		}
	}
	return nil
}

func singleValue(values url.Values, key string) (string, bool) {
	items, ok := values[key]
	if !ok || len(items) != 1 {
		return "", false
	}
	return items[0], true
}

func commaSeparated(values url.Values, key string) ([]string, error) {
	items, present := values[key]
	if !present {
		return []string{}, nil
	}
	if len(items) != 1 {
		return nil, &InvalidFilterError{Key: key}
	}
	if items[0] == "" {
		return []string{}, nil
	}
	parts := strings.Split(items[0], ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts, nil
}
