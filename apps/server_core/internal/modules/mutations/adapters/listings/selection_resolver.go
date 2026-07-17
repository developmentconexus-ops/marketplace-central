package listings

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"

	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
)

const selectionPageSize = 200
const selectionSignalLimit = 2001

type Selection struct {
	Mode       string
	ListingIDs []string
	Filter     listingsdomain.ListingFilter
	Q          string
}

type selectionJSON struct {
	Mode       string                     `json:"mode"`
	ListingIDs []string                   `json:"listing_ids"`
	Filter     map[string]json.RawMessage `json:"filter"`
	Q          string                     `json:"q"`
}

func ParseSelectionJSON(payload []byte) (Selection, error) {
	var input selectionJSON
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		return Selection{}, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return Selection{}, errors.New("invalid selection JSON")
	}
	result := Selection{Mode: input.Mode, ListingIDs: input.ListingIDs, Q: strings.TrimSpace(input.Q)}
	for key, raw := range input.Filter {
		var value string
		if key == "has_exception" {
			var boolean bool
			if err := json.Unmarshal(raw, &boolean); err != nil {
				return Selection{}, &listingsdomain.InvalidFilterError{Key: key}
			}
			if boolean {
				value = "true"
			} else {
				value = "false"
			}
		} else if err := json.Unmarshal(raw, &value); err != nil {
			return Selection{}, &listingsdomain.InvalidFilterError{Key: key}
		}
		if err := listingsdomain.SetFilterValue(&result.Filter, key, value); err != nil {
			return Selection{}, err
		}
	}
	if result.Mode != "explicit" && result.Mode != "filter" {
		return Selection{}, errors.New("invalid selection mode")
	}
	if result.Mode == "explicit" && (len(result.ListingIDs) == 0 || input.Filter != nil || result.Q != "") {
		return Selection{}, errors.New("invalid explicit selection")
	}
	if result.Mode == "filter" && len(result.ListingIDs) != 0 {
		return Selection{}, errors.New("invalid filter selection")
	}
	return result, nil
}

type listingReader interface {
	List(context.Context, listingsports.ListingQuery) (listingsports.ListingRowPage, error)
}

type SelectionResolver struct{ reader listingReader }

func NewSelectionResolver(reader listingReader) SelectionResolver {
	return SelectionResolver{reader: reader}
}

func (r SelectionResolver) Resolve(ctx context.Context, installationID string, selection Selection) ([]string, error) {
	if selection.Mode == "explicit" {
		return append([]string(nil), selection.ListingIDs...), nil
	}
	ids := make([]string, 0, selectionSignalLimit)
	var cursor listingsports.ListingCursor
	for {
		page, err := r.reader.List(ctx, listingsports.ListingQuery{InstallationID: installationID, Filter: selection.Filter, Q: selection.Q, Cursor: cursor, Limit: selectionPageSize})
		if err != nil {
			return nil, err
		}
		for _, item := range page.Items {
			ids = append(ids, item.ListingID)
			if len(ids) == selectionSignalLimit {
				return ids, nil
			}
		}
		if page.NextCursor == nil {
			return ids, nil
		}
		cursor = *page.NextCursor
	}
}
