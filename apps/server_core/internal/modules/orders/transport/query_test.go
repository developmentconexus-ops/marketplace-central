package transport

import (
	"errors"
	"net/url"
	"testing"
)

func TestParseOrderQueryRejectsUncoveredInvalidBranches(t *testing.T) {
	cases := []struct {
		name   string
		values url.Values
		key    string
	}{
		{"limit below range", url.Values{"installation_id": {"inst-1"}, "limit": {"0"}}, "limit"},
		{"limit above range", url.Values{"installation_id": {"inst-1"}, "limit": {"101"}}, "limit"},
		{"missing installation", url.Values{}, "installation_id"},
		{"inverted date range", url.Values{"installation_id": {"inst-1"}, "date_from": {"2026-07-17"}, "date_to": {"2026-07-16"}}, "date_from"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseOrderQuery(tc.values)
			var typed *InvalidFilterError
			if !errors.As(err, &typed) || typed.Key != tc.key {
				t.Fatalf("error = %#v; want invalid_filter for %s", err, tc.key)
			}
		})
	}
}

func TestParseOrderQueryClassifiesBadKeysDeterministically(t *testing.T) {
	values := url.Values{"z_bad": {"x"}, "a_bad": {"x"}, "installation_id": {"inst-1"}}
	for range 20 {
		_, err := ParseOrderQuery(values)
		var typed *InvalidFilterError
		if !errors.As(err, &typed) || typed.Key != "a_bad" {
			t.Fatalf("error = %#v; want stable first key a_bad", err)
		}
	}
}

func TestParseOrderQueryRepeatedCursorIsInvalidFilter(t *testing.T) {
	_, err := ParseOrderQuery(url.Values{"installation_id": {"inst-1"}, "cursor": {"one", "two"}})
	var typed *InvalidFilterError
	if !errors.As(err, &typed) || typed.Key != "cursor" {
		t.Fatalf("error = %#v; want invalid_filter for cursor", err)
	}
}

func TestParseOrderQueryRejectsMalformedDateAsInvalidFilter(t *testing.T) {
	_, err := ParseOrderQuery(url.Values{
		"installation_id": {"inst-1"},
		"date_from":       {"2026-07-16T25:00:00Z"},
	})

	var typed *InvalidFilterError
	if !errors.As(err, &typed) || typed.Key != "date_from" || typed.Code() != "invalid_filter" {
		t.Fatalf("ParseOrderQuery() error = %#v; want invalid_filter for date_from", err)
	}
}

func TestParseOrderQueryRejectsUnknownOrRepeatedFilter(t *testing.T) {
	cases := []url.Values{
		{"installation_id": {"inst-1"}, "filter.xyz": {"value"}},
		{"installation_id": {"inst-1"}, "status": {"paid", "cancelled"}},
	}

	for _, values := range cases {
		_, err := ParseOrderQuery(values)
		var typed *InvalidFilterError
		if !errors.As(err, &typed) || typed.Key == "" {
			t.Errorf("ParseOrderQuery(%v) error = %#v; want keyed invalid_filter", values, err)
		}
	}
}

func TestParseOrderQueryFulfillmentReturnsUnsupportedFilter(t *testing.T) {
	_, err := ParseOrderQuery(url.Values{
		"installation_id": {"inst-1"},
		"fulfillment":     {"shipped"},
	})

	var typed *UnsupportedFilterError
	if !errors.As(err, &typed) || typed.Key != "fulfillment" || typed.Code() != "unsupported_filter" {
		t.Fatalf("ParseOrderQuery() error = %#v; want unsupported_filter for fulfillment", err)
	}
}

func TestParseOrderQueryPreservesExistingLimitDefault(t *testing.T) {
	got, err := ParseOrderQuery(url.Values{"installation_id": {"inst-1"}})
	if err != nil {
		t.Fatal(err)
	}
	if got.InstallationID != "inst-1" || got.Limit != 20 {
		t.Fatalf("ParseOrderQuery() = %#v; want installation_id inst-1 and limit 20", got)
	}
}
