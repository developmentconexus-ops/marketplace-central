package migrations

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

func TestListingsMigrationMatchesIC02(t *testing.T) {
	raw, err := fs.ReadFile(Source(), "0036_listings.sql")
	if err != nil {
		t.Fatal(err)
	}
	sql := string(raw)
	normalizedSQL := normalizeMigrationSQL(sql)

	listingBody := createTableBody(t, sql, "listings")
	assertExactColumns(t, listingBody, []string{
		"tenant_id", "installation_id", "provider", "provider_listing_id", "variation_id", "title",
		"listing_type_code", "status", "price_amount", "price_currency", "published_quantity",
		"sync_state", "sync_error", "quality_score", "sales_30d", "fetched_at", "created_at", "updated_at",
	})
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS listings",
		"literal '-' at the application boundary",
		"variation_id TEXT NOT NULL",
		"price_amount NUMERIC",
		"price_currency TEXT",
		"published_quantity INT",
		"sync_error JSONB",
		"quality_score INT",
		"sales_30d INT",
		"fetched_at TIMESTAMPTZ",
		"PRIMARY KEY (tenant_id, installation_id, provider_listing_id, variation_id)",
		"CONSTRAINT listings_status_check CHECK (status IN ('active', 'paused', 'closed', 'unknown'))",
		"CONSTRAINT listings_sync_state_check CHECK (sync_state IN ('synced', 'error', 'stale', 'queued', 'syncing', 'paused_sync'))",
		"CONSTRAINT listings_quality_score_check CHECK (quality_score IS NULL OR quality_score BETWEEN 0 AND 100)",
		"CONSTRAINT listings_price_pair_check CHECK ((price_amount IS NULL AND price_currency IS NULL) OR (price_amount IS NOT NULL AND price_currency IS NOT NULL))",
		"TIMESTAMPTZ NOT NULL DEFAULT now()",
		"CREATE INDEX IF NOT EXISTS idx_listings_f02_title_key ON listings (tenant_id, installation_id, title, provider_listing_id, variation_id)",
	} {
		if !strings.Contains(normalizedSQL, normalizeMigrationSQL(required)) {
			t.Fatalf("listings migration missing %q", required)
		}
	}
	for _, column := range []string{
		"listing_type_code", "price_amount", "price_currency", "published_quantity",
		"sync_error", "quality_score", "sales_30d", "fetched_at",
	} {
		assertNullableColumn(t, listingBody, column)
	}
	for _, column := range []string{
		"tenant_id", "installation_id", "provider", "provider_listing_id", "variation_id",
		"title", "status", "sync_state", "created_at", "updated_at",
	} {
		assertNotNullableColumn(t, listingBody, column)
	}
	assertOnlyLifecycleDefaults(t, listingBody)

	for _, forbidden := range []string{"link", "product", "seller_sku", "cost", "below_margin", "pending_issue"} {
		if strings.Contains(strings.ToLower(listingBody), forbidden) {
			t.Errorf("listings table contains forbidden field text %q", forbidden)
		}
	}
	if strings.Contains(strings.ToLower(sql), "foreign key") {
		t.Error("listings migration must not add a foreign key")
	}
	if got := strings.Count(sql, "ON listings ("); got != 1 {
		t.Fatalf("listings migration has %d listings indexes, want exactly one F-02-ready index", got)
	}

	eventsBody := createTableBody(t, sql, "listing_sync_events")
	assertExactColumns(t, eventsBody, []string{
		"tenant_id", "installation_id", "provider_listing_id", "variation_id", "event_id", "at", "kind", "message_pt",
	})
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS listing_sync_events",
		"PRIMARY KEY (tenant_id, installation_id, provider_listing_id, variation_id, event_id)",
		"CONSTRAINT listing_sync_events_kind_check CHECK (kind IN ('synced', 'sync_error', 'closed', 'paused', 'refreshed'))",
		"CREATE INDEX IF NOT EXISTS idx_listing_sync_events_timeline ON listing_sync_events (tenant_id, installation_id, provider_listing_id, variation_id, at DESC, event_id DESC)",
		"at TIMESTAMPTZ NOT NULL",
	} {
		if !strings.Contains(normalizedSQL, normalizeMigrationSQL(required)) {
			t.Fatalf("listing sync events schema missing %q", required)
		}
	}
	assertOnlyLifecycleDefaults(t, eventsBody)
}

func normalizeMigrationSQL(sql string) string {
	sql = strings.Join(strings.Fields(sql), " ")
	sql = strings.ReplaceAll(sql, "( ", "(")
	sql = strings.ReplaceAll(sql, " )", ")")
	sql = strings.ReplaceAll(sql, " ,", ",")
	return sql
}

func createTableBody(t *testing.T, sql, table string) string {
	t.Helper()
	re := regexp.MustCompile(`(?is)CREATE TABLE IF NOT EXISTS ` + regexp.QuoteMeta(table) + `\s*\((.*?)\n\);`)
	match := re.FindStringSubmatch(sql)
	if len(match) != 2 {
		t.Fatalf("migration missing table %q", table)
	}
	return match[1]
}

func assertExactColumns(t *testing.T, body string, want []string) {
	t.Helper()
	got := make([]string, 0, len(want))
	for _, line := range strings.Split(body, "\n") {
		fields := strings.Fields(strings.TrimSpace(strings.TrimSuffix(line, ",")))
		if len(fields) == 0 || strings.HasPrefix(fields[0], "CONSTRAINT") || strings.HasPrefix(fields[0], "PRIMARY") ||
			fields[0] == "CHECK" || fields[0] == "OR" || fields[0] == "AND" {
			continue
		}
		got = append(got, fields[0])
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("columns = %v, want %v", got, want)
	}
}

func assertNullableColumn(t *testing.T, body, column string) {
	t.Helper()
	declaration := columnDeclaration(t, body, column)
	if strings.Contains(strings.ToUpper(declaration), "NOT NULL") {
		t.Errorf("%s must permit SQL NULL: %s", column, declaration)
	}
}

func assertNotNullableColumn(t *testing.T, body, column string) {
	t.Helper()
	declaration := columnDeclaration(t, body, column)
	if !strings.Contains(strings.ToUpper(declaration), "NOT NULL") {
		t.Errorf("%s must be NOT NULL: %s", column, declaration)
	}
}

func columnDeclaration(t *testing.T, body, column string) string {
	t.Helper()
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, column+" ") {
			return strings.TrimSuffix(trimmed, ",")
		}
	}
	t.Fatalf("column %q not found", column)
	return ""
}

func assertOnlyLifecycleDefaults(t *testing.T, body string) {
	t.Helper()
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.Contains(strings.ToUpper(trimmed), "DEFAULT") {
			continue
		}
		fields := strings.Fields(trimmed)
		if len(fields) == 0 || (fields[0] != "created_at" && fields[0] != "updated_at") {
			t.Errorf("only created_at and updated_at may have defaults: %s", trimmed)
		}
	}
}
