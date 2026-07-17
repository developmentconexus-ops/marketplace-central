package transport

import (
	"errors"
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
