package migrations

import (
	"io/fs"
	"strings"
	"testing"
)

func TestMutationEnvelopeMigrationMatchesIC03(t *testing.T) {
	raw, err := fs.ReadFile(Source(), "0038_mutation_protocols.sql")
	if err != nil {
		t.Fatal(err)
	}
	sql := string(raw)
	normalizedSQL := normalizeMigrationSQL(sql)

	protocolsBody := createTableBody(t, sql, "mutation_protocols")
	assertExactColumns(t, protocolsBody, []string{
		"tenant_id", "protocol_id", "installation_id", "type", "state", "actor", "intent",
		"selection", "totals", "source_as_of", "retried_from", "created_at", "previewed_at",
		"approved_at", "finished_at",
	})
	for _, required := range []string{
		"PRIMARY KEY (tenant_id, protocol_id)",
		"CONSTRAINT mutation_protocols_type_check CHECK (type IN ('price_update', 'stock_correct', 'link_apply', 'listing_pause', 'listing_resync', 'listing_edit', 'listing_create'))",
		"CONSTRAINT mutation_protocols_state_check CHECK (state IN ('draft', 'previewed', 'approved', 'applying', 'applied', 'partially_failed', 'failed_preserved', 'cancelled'))",
		"CREATE INDEX IF NOT EXISTS mutation_protocols_claim_idx ON mutation_protocols (tenant_id, installation_id, created_at) WHERE state IN ('approved', 'applying')",
	} {
		if !strings.Contains(normalizedSQL, normalizeMigrationSQL(required)) {
			t.Errorf("mutation protocols schema missing %q", required)
		}
	}
	for _, column := range []string{
		"tenant_id", "protocol_id", "installation_id", "type", "state", "actor", "intent",
		"selection", "totals", "created_at",
	} {
		assertNotNullableColumn(t, protocolsBody, column)
	}
	for _, column := range []string{
		"source_as_of", "retried_from", "previewed_at", "approved_at", "finished_at",
	} {
		assertNullableColumn(t, protocolsBody, column)
	}

	itemsBody := createTableBody(t, sql, "mutation_items")
	assertExactColumns(t, itemsBody, []string{
		"tenant_id", "protocol_id", "seq", "item_id", "listing_id", "idempotency_key",
		"before", "after", "state", "failure", "applied_at",
	})
	for _, required := range []string{
		"PRIMARY KEY (tenant_id, protocol_id, seq)",
		"CONSTRAINT mutation_items_idempotency_key_unique UNIQUE (tenant_id, idempotency_key)",
		"CONSTRAINT mutation_items_state_check CHECK (state IN ('previewed', 'approved', 'applying', 'applied', 'failed', 'skipped'))",
	} {
		if !strings.Contains(normalizedSQL, normalizeMigrationSQL(required)) {
			t.Errorf("mutation items schema missing %q", required)
		}
	}
	for _, column := range []string{
		"tenant_id", "protocol_id", "seq", "item_id", "listing_id", "idempotency_key", "after", "state",
	} {
		assertNotNullableColumn(t, itemsBody, column)
	}
	for _, column := range []string{"before", "failure", "applied_at"} {
		assertNullableColumn(t, itemsBody, column)
	}

	if strings.Contains(strings.ToUpper(protocolsBody+itemsBody), "DEFAULT") {
		t.Error("mutation fact columns must not have database defaults")
	}
	if strings.Contains(strings.ToUpper(sql), "FOREIGN KEY") {
		t.Error("mutation migration must match the existing no-foreign-key migration style")
	}
}
