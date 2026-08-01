package mercadolivre

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// ListListingsScanIDs is the ids-only counterpart to ListListingsScanPage
// (capability_adapter.go): it issues the SAME search_type=scan request (same
// query params, same access-token/error handling), but returns BEFORE the
// N+1 ReadListing hydration loop that ListListingsScanPage runs per id.
//
// This is Passo 1 of F-03 (backfill-cursor-ingest): IC-06 requires the
// enumerator to never hydrate and the hydrator to never enumerate — this
// method is the enumerator half. Hydration for a page of ids returned here
// happens separately, in batches of 20, via GetItemsMultiget
// (items_multiget_reader.go), not via ReadListing/ReadStock (which are 1
// HTTP call per item and belong to the refresh/webhook path, not backfill).
//
// ListListingsScanPage itself is untouched (capability_adapter.go is
// frozen for this feature) — this is a new, independent method on the same
// receiver, reusing its unexported helpers (normalizeAccountRef,
// a.accessToken, limitOrDefault, a.doJSON) exactly as ListListingsScanPage
// does.
func (a *CapabilityAdapter) ListListingsScanIDs(ctx context.Context, input domain.ListListingsInput) ([]string, string, error) {
	accountRef, err := normalizeAccountRef(input.AccountRef)
	if err != nil {
		return nil, "", err
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return nil, "", err
	}

	query := url.Values{}
	query.Set("search_type", "scan")
	query.Set("limit", strconv.Itoa(limitOrDefault(input.Limit, 50)))
	if cursor := strings.TrimSpace(input.Cursor); cursor != "" {
		query.Set("scroll_id", cursor)
	}

	var response struct {
		ScrollID string   `json:"scroll_id"`
		Results  []string `json:"results"`
	}
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, "/users/"+url.PathEscape(accountRef.ProviderAccountID)+"/items/search?"+query.Encode(), nil, &response); err != nil {
		return nil, "", err
	}

	ids := make([]string, 0, len(response.Results))
	for _, id := range response.Results {
		id = strings.TrimSpace(id)
		if id != "" {
			ids = append(ids, id)
		}
	}
	return ids, strings.TrimSpace(response.ScrollID), nil
}
