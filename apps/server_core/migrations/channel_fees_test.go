package migrations

import (
	"strings"
	"testing"
)

const channelFeesMigration = "0086_channel_fees.sql"

func TestChannelFeesMigrationMatchesIC01(t *testing.T) {
	sql := readMigration(t, channelFeesMigration)
	normalizedSQL := normalizeMigrationSQL(sql)

	body := createTableBody(t, sql, "channel_fees")
	assertExactColumns(t, body, []string{
		"id", "tenant_id", "provider", "installation_id", "layer", "fee_kind",
		"subject_type", "subject_id", "value_type", "value", "currency", "detail",
		"origem", "coletado_em", "source_time", "created_at", "updated_at",
	})

	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS channel_fees",
		"id bigserial PRIMARY KEY",
		"layer smallint NOT NULL",
		"value numeric(14,4) NOT NULL",
		"currency char(3)",
		"detail jsonb",
		"CONSTRAINT channel_fees_layer_check CHECK (layer IN (1, 2, 3))",
		"CONSTRAINT channel_fees_fee_kind_check CHECK (fee_kind IN ('commission', 'freight'))",
		"CONSTRAINT channel_fees_subject_type_check CHECK (subject_type IN ('category', 'listing', 'order', 'order_line'))",
		"CONSTRAINT channel_fees_value_type_check CHECK (value_type IN ('percent', 'amount'))",
		"CONSTRAINT channel_fees_currency_when_amount_check CHECK ((value_type = 'amount' AND currency IS NOT NULL) OR (value_type = 'percent' AND currency IS NULL))",
		"CONSTRAINT channel_fees_origem_check CHECK (origem IN ('api_listing_prices', 'api_shipping_options', 'api_order', 'api_shipment', 'config'))",
		"CONSTRAINT channel_fees_natural_key_key UNIQUE (tenant_id, provider, installation_id, subject_type, subject_id, layer, fee_kind)",
	} {
		if !strings.Contains(normalizedSQL, normalizeMigrationSQL(required)) {
			t.Errorf("channel_fees migration missing %q", required)
		}
	}

	for _, column := range []string{"currency", "detail", "source_time"} {
		assertNullableColumn(t, body, column)
	}
	for _, column := range []string{
		"tenant_id", "provider", "installation_id", "layer", "fee_kind",
		"subject_type", "subject_id", "value_type", "value", "origem",
		"coletado_em", "created_at", "updated_at",
	} {
		assertNotNullableColumn(t, body, column)
	}
	assertOnlyLifecycleDefaults(t, body)

	if strings.Contains(strings.ToUpper(sql), "FOREIGN KEY") {
		t.Error("channel_fees migration must not add a foreign key")
	}
	if strings.Contains(strings.ToLower(body), "raw") {
		t.Error("channel_fees table must not carry a raw column")
	}
}
