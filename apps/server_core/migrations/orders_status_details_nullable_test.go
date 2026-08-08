package migrations

import (
	"strings"
	"testing"
)

const ordersStatusDetailsNullableMigration = "0093_orders_status_details_nullable.sql"

// TestOrdersStatusDetailsBecomeNullable proves 0093 relaxes provider_status_detail
// and cancellation_detail off their 0027 `NOT NULL DEFAULT ”` declaration —
// the DB-level half of making NULL reachable (domain.MarketplaceOrder's
// ProviderStatusDetail/CancellationDetail are *string; without dropping the
// constraint here, writing nil would still fail at the database).
func TestOrdersStatusDetailsBecomeNullable(t *testing.T) {
	sql := normalizeMigrationSQL(readMigration(t, ordersStatusDetailsNullableMigration))
	lowerSQL := strings.ToLower(sql)

	for _, statement := range []string{
		"alter table orders_marketplace_orders alter column provider_status_detail drop not null",
		"alter table orders_marketplace_orders alter column provider_status_detail drop default",
		"alter table orders_marketplace_orders alter column cancellation_detail drop not null",
		"alter table orders_marketplace_orders alter column cancellation_detail drop default",
	} {
		if !strings.Contains(lowerSQL, statement) {
			t.Errorf("0093 migration missing exact statement %q", statement)
		}
	}

	// No new NOT NULL / DEFAULT should sneak back in on either column.
	for _, forbidden := range []string{
		"provider_status_detail text not null",
		"cancellation_detail text not null",
	} {
		if strings.Contains(lowerSQL, forbidden) {
			t.Errorf("0093 migration must not re-add %q", forbidden)
		}
	}
}
