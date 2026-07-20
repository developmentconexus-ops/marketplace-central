package migrations

import (
	"io/fs"
	"strings"
	"testing"
)

// TestSyncStateMigrationShapeE8 asserts migration 0075 matches the E8 sync_state
// shape (M01-C1) and stays cadence-agnostic (M01-C6 at the migration layer)
// without needing a live database.
func TestSyncStateMigrationShapeE8(t *testing.T) {
	raw, err := fs.ReadFile(Source(), "0075_sync_sync_state.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	sql := string(raw)
	body := createTableBody(t, sql, "sync_state")

	assertExactColumns(t, body, []string{
		"tenant_id", "installation_id", "entity",
		"cursor", "schedule", "last_full_sync_at", "last_incremental_at",
		"last_error", "consecutive_failures",
	})

	for _, col := range []string{"tenant_id", "installation_id", "entity", "consecutive_failures"} {
		assertNotNullableColumn(t, body, col)
	}
	for _, col := range []string{"cursor", "schedule", "last_full_sync_at", "last_incremental_at", "last_error"} {
		assertNullableColumn(t, body, col)
	}

	if !strings.Contains(body, "PRIMARY KEY (tenant_id, installation_id, entity)") {
		t.Errorf("sync_state must have composite PRIMARY KEY (tenant_id, installation_id, entity): %s", body)
	}

	// entity is a semantic enum validated in the application layer (F-01), never a
	// DB CHECK — so a future entity needs no migration.
	if strings.Contains(strings.ToUpper(body), "CHECK") {
		t.Errorf("sync_state must not constrain entity with a DB CHECK: %s", body)
	}

	// Cadence-agnostic (D6): no fixed business cadence baked into the DDL.
	low := strings.ToLower(sql)
	for _, banned := range []string{"daily", "diário", "diario", "hourly", "weekly"} {
		if strings.Contains(low, banned) {
			t.Errorf("migration must stay cadence-agnostic; found %q", banned)
		}
	}

	// Additive, no seed.
	if strings.Contains(strings.ToUpper(sql), "INSERT") {
		t.Error("sync_state migration must not seed rows")
	}
	if strings.Contains(strings.ToUpper(sql), "ALTER TABLE") {
		t.Error("sync_state migration must be a pure additive CREATE, no ALTER")
	}
}
