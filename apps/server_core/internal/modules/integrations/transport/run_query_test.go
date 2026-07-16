package transport

import (
	"errors"
	"net/url"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

func TestParseRunQueryAcceptsStrictFilters(t *testing.T) {
	t.Parallel()

	query, err := ParseRunQuery(url.Values{
		"installation_id": {"inst_001"},
		"limit":           {"25"},
		"filter.module":   {"listing_read"},
		"filter.status":   {"running"},
	})
	if err != nil {
		t.Fatalf("ParseRunQuery() error = %v", err)
	}
	if query.InstallationID != "inst_001" || query.Limit != 25 || query.Module != "listing_read" || query.Status != domain.OperationRunStatusRunning {
		t.Fatalf("query = %#v", query)
	}
}

func TestParseRunQueryRejectsInvalidStatus(t *testing.T) {
	t.Parallel()

	_, err := ParseRunQuery(url.Values{"installation_id": {"inst_001"}, "filter.status": {"not-a-status"}})
	if err == nil {
		t.Fatal("ParseRunQuery() error = nil, want invalid_filter")
	}
	var invalid *InvalidRunFilterError
	if !errors.As(err, &invalid) || invalid.Code() != "invalid_filter" {
		t.Fatalf("error = %#v, want InvalidRunFilterError", err)
	}
}

func TestParseRunQueryPassesThroughUnknownModule(t *testing.T) {
	t.Parallel()

	query, err := ParseRunQuery(url.Values{"installation_id": {"inst_001"}, "filter.module": {"future_operation"}})
	if err != nil {
		t.Fatalf("ParseRunQuery() error = %v, want pass-through module", err)
	}
	if query.Module != "future_operation" {
		t.Fatalf("module = %q, want future_operation", query.Module)
	}
}
