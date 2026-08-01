package migrations

import (
	"strings"
	"testing"
)

const divergencesMigration = "0087_divergences.sql"

func TestDivergencesMigrationMatchesIC02(t *testing.T) {
	sql := readMigration(t, divergencesMigration)
	normalizedSQL := normalizeMigrationSQL(sql)

	body := createTableBody(t, sql, "divergences")
	assertExactColumns(t, body, []string{
		"id", "tenant_id", "provider", "entity_type", "entity_id", "kind",
		"expected_value", "observed_value", "expected_source", "observed_source",
		"expected_observed_at", "observed_observed_at", "detected_at",
		"last_evaluated_at", "resolved_at", "created_at", "updated_at",
	})

	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS divergences",
		"id bigserial PRIMARY KEY",
		"expected_value numeric(14,4) NOT NULL",
		"observed_value numeric(14,4) NOT NULL",
		"expected_observed_at timestamptz NOT NULL",
		"observed_observed_at timestamptz NOT NULL",
		"detected_at timestamptz NOT NULL",
		"last_evaluated_at timestamptz NOT NULL",
		"resolved_at timestamptz",
		"CONSTRAINT divergences_entity_type_check CHECK (entity_type IN ('listing', 'order_line'))",
		"CONSTRAINT divergences_kind_check CHECK (kind IN ('estoque', 'tarifa'))",
		"CREATE UNIQUE INDEX IF NOT EXISTS divergences_one_open_row ON divergences (tenant_id, provider, entity_type, entity_id, kind) WHERE resolved_at IS NULL",
		"CREATE INDEX IF NOT EXISTS idx_divergences_open ON divergences (tenant_id, resolved_at) WHERE resolved_at IS NULL",
	} {
		if !strings.Contains(normalizedSQL, normalizeMigrationSQL(required)) {
			t.Errorf("divergences migration missing %q", required)
		}
	}

	assertNullableColumn(t, body, "resolved_at")
	for _, column := range []string{
		"tenant_id", "provider", "entity_type", "entity_id", "kind",
		"expected_value", "observed_value", "expected_source", "observed_source",
		"expected_observed_at", "observed_observed_at", "detected_at",
		"last_evaluated_at", "created_at", "updated_at",
	} {
		assertNotNullableColumn(t, body, column)
	}
	assertOnlyLifecycleDefaults(t, body)

	if strings.Contains(strings.ToUpper(sql), "FOREIGN KEY") {
		t.Error("divergences migration must not add a foreign key")
	}

	// The load-bearing constraint: a named partial unique index, not a
	// generic table constraint, so a must-fail test can assert its name.
	if strings.Count(normalizedSQL, "divergences_one_open_row") != 1 {
		t.Fatalf("divergences_one_open_row must be declared exactly once")
	}
}
