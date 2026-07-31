package migrations

import (
	"strings"
	"testing"
)

func TestSellableAssortmentMigrationsDeclareDefaultsAndNullableMirrorFields(t *testing.T) {
	configSQL := strings.ToLower(normalizeMigrationSQL(readMigration(t, "0083_sellable_assortment_config.sql")))
	for _, declaration := range []string{
		"alter table active_source add column only_revenda boolean not null default true",
		"alter table active_source add column only_em_estoque boolean not null default true",
		"alter table active_source add column only_ecommerce_eligible boolean not null default false",
	} {
		if !strings.Contains(configSQL, declaration) {
			t.Errorf("active_source migration missing exact declaration %q", declaration)
		}
	}

	mirrorSQL := strings.ToLower(normalizeMigrationSQL(readMigration(t, "0084_products_mirror_sellable_fields.sql")))
	for _, declaration := range []string{
		"alter table products_mirror add column usoprod text",
		"alter table products_mirror add column ad_ecommerce text",
	} {
		if !strings.Contains(mirrorSQL, declaration) {
			t.Errorf("products_mirror migration missing nullable declaration %q", declaration)
		}
	}
	for _, column := range []string{"usoprod", "ad_ecommerce"} {
		if strings.Contains(mirrorSQL, column+" text default") || strings.Contains(mirrorSQL, column+" text not null") {
			t.Errorf("products_mirror.%s must have no DEFAULT and no NOT NULL", column)
		}
	}
}
