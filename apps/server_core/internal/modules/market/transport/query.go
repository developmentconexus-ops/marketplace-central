package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

type CodprodQuery struct {
	Codprods []string
}

type SignalQuery struct {
	ListingIDs []string
}

type CollectionRequest struct {
	Codprod string
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

func ParseCodprodQuery(values url.Values) (CodprodQuery, error) {
	if err := rejectUnknown(values, map[string]struct{}{"codprod": {}}); err != nil {
		return CodprodQuery{}, err
	}
	codprods, err := commaSeparated(values, "codprod")
	if err != nil {
		return CodprodQuery{}, err
	}
	return CodprodQuery{Codprods: codprods}, nil
}

func ParseSignalQuery(values url.Values) (SignalQuery, error) {
	if err := rejectUnknown(values, map[string]struct{}{"listing_ids": {}}); err != nil {
		return SignalQuery{}, err
	}
	listingIDs, err := commaSeparated(values, "listing_ids")
	if err != nil {
		return SignalQuery{}, err
	}
	return SignalQuery{ListingIDs: listingIDs}, nil
}

// ParseCollectionRequest decodes the POST /market/collections body. A
// missing or blank codprod is an invalid_filter, not a silent no-op — the
// caller must name exactly one product (D-F4-p: single-product synchronous).
func ParseCollectionRequest(r *http.Request) (CollectionRequest, error) {
	defer r.Body.Close()
	var body struct {
		Codprod string `json:"codprod"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return CollectionRequest{}, &InvalidFilterError{Key: "codprod"}
	}
	codprod := strings.TrimSpace(body.Codprod)
	if codprod == "" {
		return CollectionRequest{}, &InvalidFilterError{Key: "codprod"}
	}
	return CollectionRequest{Codprod: codprod}, nil
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
