package listings

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	listingstransport "marketplace-central/apps/server_core/internal/modules/listings/transport"
)

type fakeListingReader struct {
	queries []listingsports.ListingQuery
	pages   []listingsports.ListingRowPage
}

func (f *fakeListingReader) List(_ context.Context, query listingsports.ListingQuery) (listingsports.ListingRowPage, error) {
	f.queries = append(f.queries, query)
	page := f.pages[0]
	f.pages = f.pages[1:]
	return page, nil
}

func TestParseSelectionJSONUsesIC02Grammar(t *testing.T) {
	cases := []struct {
		key, scalar, jsonValue string
		valid                  bool
	}{
		{"status", "active", `"active"`, true}, {"status", "banana", `"banana"`, false},
		{"sync_state", "paused_sync", `"paused_sync"`, true}, {"link_state", "resolved", `"resolved"`, true},
		{"exception", "below_margin", `"below_margin"`, true}, {"has_exception", "true", `true`, true},
		{"has_exception", "1", `1`, false}, {"listing_type_code", "gold_pro", `"gold_pro"`, true},
		{"product_id", "42664", `"42664"`, true}, {"product_id", "042664", `"042664"`, false},
		{"bogus", "x", `"x"`, false}, {"status[]", "active", `"active"`, false}, {"status", "gte:active", `"gte:active"`, false},
	}
	for _, tc := range cases {
		payload := []byte(fmt.Sprintf(`{"mode":"filter","filter":{"%s":%s}}`, tc.key, tc.jsonValue))
		_, jsonErr := ParseSelectionJSON(payload)
		_, urlErr := listingstransport.ParseListingQuery(url.Values{"filter." + tc.key: {tc.scalar}})
		if (jsonErr == nil) != tc.valid || (urlErr == nil) != tc.valid {
			t.Errorf("%s=%s JSON error=%v URL error=%v; valid=%v", tc.key, tc.jsonValue, jsonErr, urlErr, tc.valid)
		}
	}
}

func TestSelectionResolverStopsAt2001AndPreservesScopeAndCursor(t *testing.T) {
	pages := make([]listingsports.ListingRowPage, 0, 12)
	for pageNumber := 0; pageNumber < 10; pageNumber++ {
		items := make([]listingsdomain.ListingReadModel, selectionPageSize)
		for i := range items {
			items[i].ListingID = fmt.Sprintf("inst-1~MLB-%d~-", pageNumber*selectionPageSize+i)
		}
		next := listingsports.ListingCursor{LastTitle: fmt.Sprintf("page-%d", pageNumber), ListingID: items[len(items)-1].ListingID}
		pages = append(pages, listingsports.ListingRowPage{Items: items, NextCursor: &next})
	}
	pages = append(pages, listingsports.ListingRowPage{Items: []listingsdomain.ListingReadModel{{ListingID: "inst-1~MLB-2000~-"}, {ListingID: "must-not-be-returned"}}})
	pages = append(pages, listingsports.ListingRowPage{Items: []listingsdomain.ListingReadModel{{ListingID: "must-not-be-fetched"}}})
	reader := &fakeListingReader{pages: pages}
	filter := listingsdomain.ListingFilter{Status: listingsdomain.ListingStatusActive}
	ids, err := NewSelectionResolver(reader).Resolve(context.Background(), "inst-1", Selection{Mode: "filter", Filter: filter, Q: "camisa"})
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2001 || len(reader.queries) != 11 {
		t.Fatalf("ids=%d queries=%d", len(ids), len(reader.queries))
	}
	for i, query := range reader.queries {
		if query.InstallationID != "inst-1" || query.Filter != filter || query.Q != "camisa" || query.Limit != selectionPageSize {
			t.Errorf("query[%d]=%#v", i, query)
		}
	}
	if reader.queries[1].Cursor != *pages[0].NextCursor {
		t.Errorf("second cursor=%#v", reader.queries[1].Cursor)
	}
}

func TestSelectionResolverReturnsHonestEmptyResult(t *testing.T) {
	reader := &fakeListingReader{pages: []listingsports.ListingRowPage{{}}}
	ids, err := NewSelectionResolver(reader).Resolve(context.Background(), "inst-1", Selection{Mode: "filter"})
	if err != nil || len(ids) != 0 {
		t.Fatalf("ids=%v err=%v", ids, err)
	}
}
