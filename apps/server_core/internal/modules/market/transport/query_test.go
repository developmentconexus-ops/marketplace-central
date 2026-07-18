package transport

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestParseObservationQueryPreservesOrderedDuplicateIDs(t *testing.T) {
	t.Parallel()

	query, err := ParseObservationQuery(url.Values{
		"installation_id": {"inst_1"},
		"listing_ids":     {"a,b,a"},
	})
	if err != nil {
		t.Fatalf("ParseObservationQuery() error = %v", err)
	}
	if query.InstallationID != "inst_1" {
		t.Fatalf("installation_id = %q, want inst_1", query.InstallationID)
	}
	if got, want := query.ListingIDs, []string{"a", "b", "a"}; !equalStrings(got, want) {
		t.Fatalf("listing_ids = %#v, want %#v", got, want)
	}
}

func TestParseObservationQueryTrimsWhitespaceInIDs(t *testing.T) {
	t.Parallel()

	query, err := ParseObservationQuery(url.Values{
		"installation_id": {"inst_1"},
		"listing_ids":     {"a, b ,a"},
	})
	if err != nil {
		t.Fatalf("ParseObservationQuery() error = %v", err)
	}
	if got, want := query.ListingIDs, []string{"a", "b", "a"}; !equalStrings(got, want) {
		t.Fatalf("listing_ids = %#v, want %#v", got, want)
	}
}

func TestParseReferenceQueryRejectsUnknownParams(t *testing.T) {
	t.Parallel()

	_, err := ParseReferenceQuery(url.Values{"product_ids": {"p1"}, "unexpected": {"x"}})
	var invalid *InvalidFilterError
	if !errors.As(err, &invalid) || invalid.Key != "unexpected" {
		t.Fatalf("error = %#v, want InvalidFilterError for unexpected", err)
	}
}

func TestParseObservationQueryEmptyIDsProducesEmptyList(t *testing.T) {
	t.Parallel()

	query, err := ParseObservationQuery(url.Values{"installation_id": {"inst_1"}, "listing_ids": {""}})
	if err != nil {
		t.Fatalf("ParseObservationQuery() error = %v", err)
	}
	if query.ListingIDs == nil || len(query.ListingIDs) != 0 {
		t.Fatalf("listing_ids = %#v, want non-nil empty list", query.ListingIDs)
	}
}

func TestParseCodprodQueryPreservesOrder(t *testing.T) {
	t.Parallel()

	query, err := ParseCodprodQuery(url.Values{"codprod": {"a,b,a"}})
	if err != nil {
		t.Fatalf("ParseCodprodQuery() error = %v", err)
	}
	if got, want := query.Codprods, []string{"a", "b", "a"}; !equalStrings(got, want) {
		t.Fatalf("codprods = %#v, want %#v", got, want)
	}
}

func TestParseCodprodQueryRejectsUnknownParams(t *testing.T) {
	t.Parallel()

	_, err := ParseCodprodQuery(url.Values{"codprod": {"p1"}, "unexpected": {"x"}})
	var invalid *InvalidFilterError
	if !errors.As(err, &invalid) || invalid.Key != "unexpected" {
		t.Fatalf("error = %#v, want InvalidFilterError for unexpected", err)
	}
}

func TestParseCodprodQueryEmptyProducesEmptyList(t *testing.T) {
	t.Parallel()

	query, err := ParseCodprodQuery(url.Values{"codprod": {""}})
	if err != nil {
		t.Fatalf("ParseCodprodQuery() error = %v", err)
	}
	if query.Codprods == nil || len(query.Codprods) != 0 {
		t.Fatalf("codprods = %#v, want non-nil empty list", query.Codprods)
	}
}

func TestParseSignalQueryPreservesOrder(t *testing.T) {
	t.Parallel()

	query, err := ParseSignalQuery(url.Values{"listing_ids": {"l1, l2 ,l1"}})
	if err != nil {
		t.Fatalf("ParseSignalQuery() error = %v", err)
	}
	if got, want := query.ListingIDs, []string{"l1", "l2", "l1"}; !equalStrings(got, want) {
		t.Fatalf("listing_ids = %#v, want %#v", got, want)
	}
}

func TestParseSignalQueryRejectsUnknownParams(t *testing.T) {
	t.Parallel()

	_, err := ParseSignalQuery(url.Values{"listing_ids": {"l1"}, "unexpected": {"x"}})
	var invalid *InvalidFilterError
	if !errors.As(err, &invalid) || invalid.Key != "unexpected" {
		t.Fatalf("error = %#v, want InvalidFilterError for unexpected", err)
	}
}

func TestParseCollectionRequestTrimsAndValidatesCodprod(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/market/collections", bytes.NewBufferString(`{"codprod":" 123 "}`))
	got, err := ParseCollectionRequest(req)
	if err != nil {
		t.Fatalf("ParseCollectionRequest() error = %v", err)
	}
	if got.Codprod != "123" {
		t.Fatalf("codprod = %q, want 123", got.Codprod)
	}
}

func TestParseCollectionRequestBlankCodprodIsInvalidFilter(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/market/collections", bytes.NewBufferString(`{"codprod":"  "}`))
	_, err := ParseCollectionRequest(req)
	var invalid *InvalidFilterError
	if !errors.As(err, &invalid) || invalid.Key != "codprod" {
		t.Fatalf("error = %#v, want InvalidFilterError for codprod", err)
	}
}

func TestParseCollectionRequestMalformedBodyIsInvalidFilter(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/market/collections", bytes.NewBufferString(`not json`))
	_, err := ParseCollectionRequest(req)
	var invalid *InvalidFilterError
	if !errors.As(err, &invalid) || invalid.Key != "codprod" {
		t.Fatalf("error = %#v, want InvalidFilterError for codprod", err)
	}
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
