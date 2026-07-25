package migrations

import (
	"io/fs"
	"strings"
	"testing"
)

const productLinkDecisionsMigration = "0082_product_link_decisions.sql"

func readMigration(t *testing.T, name string) string {
	t.Helper()
	raw, err := fs.ReadFile(Source(), name)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func TestProductLinkDecisionsMatchesE10(t *testing.T) {
	sql := normalizeMigrationSQL(readMigration(t, productLinkDecisionsMigration))
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS product_link_decisions",
		"tenant_id text NOT NULL",
		"decision_id text NOT NULL",
		"rule_matched text NOT NULL",
		"actor text NOT NULL",
		"created_at timestamptz NOT NULL",
		"superseded_by text",
		"PRIMARY KEY (tenant_id, decision_id)",
		"CHECK (actor IN ('system', 'operator'))",
		"CHECK (rule_matched IN ('exact_codprod_unique', 'exact_ean_unique', 'concordant_codprod_ean', 'manual'))",
		"FOREIGN KEY (tenant_id, superseded_by) REFERENCES product_link_decisions (tenant_id, decision_id)",
	} {
		if !strings.Contains(sql, normalizeMigrationSQL(required)) {
			t.Errorf("E10 audit schema missing %q", required)
		}
	}
	if !strings.Contains(sql, normalizeMigrationSQL("link_id text GENERATED ALWAYS AS (")) {
		t.Errorf("E10 audit schema must carry the contract-named link_id")
	}
}

// ADR-17: the collision count is a fact that was either read or not. NULL says
// "not read"; 0 would say "the anchor matched no product", which no decision
// could have been taken on.
func TestCollisionsAtDecisionIsHonestlyUnknownNeverZero(t *testing.T) {
	sql := readMigration(t, productLinkDecisionsMigration)
	body := createTableBody(t, sql, "product_link_decisions")
	assertNullableColumn(t, body, "collisions_at_decision")
	if strings.Contains(normalizeMigrationSQL(body), "collisions_at_decision integer NOT NULL") ||
		strings.Contains(normalizeMigrationSQL(body), "collisions_at_decision integer DEFAULT") {
		t.Errorf("collisions_at_decision must never carry a default")
	}
	if !strings.Contains(normalizeMigrationSQL(sql), normalizeMigrationSQL("CHECK (collisions_at_decision IS NULL OR collisions_at_decision >= 1)")) {
		t.Errorf("a decision taken on zero matched products is not a possible reading")
	}
}

// D-121-2: corroboration is the only rule the system may apply by itself. The
// DB refuses an actor='system' row on any single-anchor rule, so a future
// caller cannot re-open the auto-approve path the operator closed.
func TestOnlyCorroboratedDecisionsMayCarryTheSystemActor(t *testing.T) {
	sql := normalizeMigrationSQL(readMigration(t, productLinkDecisionsMigration))
	if !strings.Contains(sql, normalizeMigrationSQL("CHECK (actor <> 'system' OR rule_matched = 'concordant_codprod_ean')")) {
		t.Fatalf("migration must bind actor='system' to the concordant rule")
	}
}

// A8 is already satisfied by the PK product_links has carried since 0025. The
// milestone proves idempotency; it does not add a key. A unique on the product
// would have capped a product at one listing and collapsed two variations of
// the same listing into one link.
func TestE10MigrationDoesNotTouchProductLinks(t *testing.T) {
	sql := strings.ToUpper(readMigration(t, productLinkDecisionsMigration))
	for _, forbidden := range []string{"ALTER TABLE PRODUCT_LINKS", "ON PRODUCT_LINKS", "CREATE TABLE IF NOT EXISTS PRODUCT_LINKS "} {
		if strings.Contains(sql, forbidden) {
			t.Errorf("E10 migration must be additive only, found %q", forbidden)
		}
	}
}
