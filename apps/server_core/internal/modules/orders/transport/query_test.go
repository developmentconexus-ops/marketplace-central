package transport

import (
	"errors"
	"net/url"
	"testing"
)

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
