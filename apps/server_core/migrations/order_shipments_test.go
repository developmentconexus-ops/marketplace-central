package migrations

import (
	"strings"
	"testing"
)

const orderShipmentsMigration = "0088_order_shipments.sql"

func TestOrderShipmentsMigrationMatchesIC03(t *testing.T) {
	sql := readMigration(t, orderShipmentsMigration)
	normalizedSQL := normalizeMigrationSQL(sql)

	body := createTableBody(t, sql, "order_shipments")
	assertExactColumns(t, body, []string{
		"tenant_id", "provider", "provider_shipment_id", "provider_order_id",
		"status", "substatus", "logistic_type", "tracking_number", "tracking_method",
		"sla_status", "sla_limit_at", "cost_gross", "cost_seller", "currency",
		"receiver_name", "dest_city", "dest_state", "dest_zip", "source_time",
		"fetched_at", "created_at", "updated_at",
	})

	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS order_shipments",
		"status text NOT NULL",
		"fetched_at timestamptz NOT NULL",
		"PRIMARY KEY (tenant_id, provider, provider_shipment_id)",
		"CREATE INDEX IF NOT EXISTS idx_order_shipments_provider_order_id ON order_shipments (provider_order_id)",
	} {
		if !strings.Contains(normalizedSQL, normalizeMigrationSQL(required)) {
			t.Errorf("order_shipments migration missing %q", required)
		}
	}

	for _, column := range []string{
		"tenant_id", "provider", "provider_shipment_id", "provider_order_id",
		"status", "fetched_at", "created_at", "updated_at",
	} {
		assertNotNullableColumn(t, body, column)
	}
	for _, column := range []string{
		"substatus", "logistic_type", "tracking_number", "tracking_method",
		"sla_status", "sla_limit_at", "cost_gross", "cost_seller", "currency",
		"receiver_name", "dest_city", "dest_state", "dest_zip", "source_time",
	} {
		assertNullableColumn(t, body, column)
	}
	assertOnlyLifecycleDefaults(t, body)

	// status/substatus carry no CHECK: the vocabulary is the provider's own
	// and can change without our say (IC-03 Enums And Statuses).
	if strings.Contains(normalizedSQL, "status_check") {
		t.Error("order_shipments status/substatus must not be constrained by a CHECK")
	}
	if strings.Contains(strings.ToUpper(sql), "FOREIGN KEY") {
		t.Error("order_shipments migration must not add a foreign key")
	}

	// Security criterion (P7 r01 B-7, ADR-025 amended): zero raw/payload-of-PII
	// columns anywhere on this table. Every declared column must be one of the
	// typed facts above, never a jsonb catch-all.
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.Contains(strings.ToLower(trimmed), "raw") {
			t.Errorf("order_shipments must not carry a raw/PII-payload column: %q", trimmed)
		}
		if strings.Contains(strings.ToLower(trimmed), "jsonb") {
			t.Errorf("order_shipments must not carry a jsonb column: %q", trimmed)
		}
	}
}
