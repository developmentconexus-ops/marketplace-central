package transport

import (
	"errors"
	"net/url"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/listings/domain"
)

func TestParseListingQueryAcceptsFixedGrammar(t *testing.T) {
	values := url.Values{"filter.status": {"active"}, "filter.sync_state": {"paused_sync"}, "filter.link_state": {"resolved"}, "filter.exception": {"below_margin"}, "filter.has_exception": {"true"}, "filter.listing_type_code": {"gold_pro"}, "filter.product_id": {"42664"}, "q": {"  ação  "}}
	got, err := ParseListingQuery(values)
	if err != nil {
		t.Fatal(err)
	}
	if got.Limit != 50 || got.Q != "ação" || got.Filter.Status != domain.ListingStatusActive || got.Filter.SyncState != domain.ListingSyncStatePausedSync || got.Filter.LinkState != domain.LinkStateResolved || got.Filter.Exception != domain.ListingExceptionBelowMargin || got.Filter.HasException == nil || !*got.Filter.HasException || got.Filter.ListingTypeCode != "gold_pro" || got.Filter.ProductID != "42664" {
		t.Fatalf("ParseListingQuery() = %#v", got)
	}
}

func TestParseListingQueryAcceptsAllFixedEnumValues(t *testing.T) {
	cases := map[string][]string{
		"status": {"active", "paused", "closed", "unknown", "under_review", "inactive", "payment_required", "not_yet_active"}, "sync_state": {"synced", "error", "stale", "queued", "syncing", "paused_sync"},
		"link_state": {"unresolved", "conflict", "resolved", "rejected"}, "exception": {"sync_error", "stale", "unlinked", "below_margin"},
	}
	for key, values := range cases {
		for _, value := range values {
			if _, err := ParseListingQuery(url.Values{"filter." + key: {value}}); err != nil {
				t.Errorf("%s=%s: %v", key, value, err)
			}
		}
	}
}

func TestParseListingQueryRejectsInvalidScalarsWithKeyDetails(t *testing.T) {
	cases := []url.Values{
		{"filter.bogus": {"x"}}, {"filter.status[]": {"active"}}, {"filter.status.ne": {"active"}}, {"filter.status": {"banana"}}, {"filter.sync_state": {"banana"}}, {"filter.link_state": {"banana"}}, {"filter.exception": {"banana"}}, {"filter.has_exception": {"1"}}, {"filter.status": {"active", "paused"}}, {"q": {"a", "b"}}, {"limit": {"0"}}, {"limit": {"201"}}, {"limit": {"abc"}}, {"limit": {"10", "20"}},
		// present-but-empty filter param is malformed, not "absent" (must not alias to default).
		{"filter.status": {""}}, {"filter.has_exception": {""}}, {"filter.product_id": {""}},
	}
	for _, values := range cases {
		_, err := ParseListingQuery(values)
		var typed *InvalidFilterError
		if !errors.As(err, &typed) || typed.Key == "" || typed.Details()["key"] != typed.Key {
			t.Errorf("ParseListingQuery(%v) error = %#v; want keyed InvalidFilterError", values, err)
			continue
		}
		if !contains(err.Error(), typed.Key) {
			t.Errorf("error %q does not name key %q", err, typed.Key)
		}
	}
}

func contains(value, part string) bool {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return true
		}
	}
	return false
}
