package migrations

import (
	"strings"
	"testing"
)

const icmsMatrixMigration = "0094_icms_matrix.sql"

// TestICMSMatrixMigrationAddsProductsMirrorFiscalColumns proves the five
// ADR-17 honest-unknown columns 0094 adds to products_mirror arrive bare
// nullable (no DEFAULT), additive-only (ADD COLUMN IF NOT EXISTS).
func TestICMSMatrixMigrationAddsProductsMirrorFiscalColumns(t *testing.T) {
	sql := normalizeMigrationSQL(readMigration(t, icmsMatrixMigration))
	lower := strings.ToLower(sql)

	// Each declaration is asserted as a full statement (through the
	// terminating ";") so a bare-nullable ALTER cannot be confused with an
	// unrelated NOT NULL column of the same name elsewhere in the file
	// (icms_matrix_mirror also declares its own, legitimately NOT NULL,
	// grupo_icms column as part of its primary key).
	for _, statement := range []string{
		"alter table products_mirror add column if not exists grupo_icms integer;",
		"alter table products_mirror add column if not exists origprod smallint;",
		"alter table products_mirror add column if not exists st_retido_entrada numeric(14,4);",
		"alter table products_mirror add column if not exists restituicao_unit numeric(14,4);",
		"alter table products_mirror add column if not exists fiscal_dt_ref date;",
	} {
		if !strings.Contains(lower, statement) {
			t.Errorf("products_mirror fiscal-columns migration missing exact bare-nullable statement %q (found NOT NULL/DEFAULT tacked on, or the ALTER is missing entirely)", statement)
		}
	}
}

// TestICMSAliquotaInternaSchema proves the LEGAL rate table shape: one open
// row per UF via the partial unique index, never a DB CHECK on the UF
// literal (a future UF/DF split needs no migration to this table's shape).
func TestICMSAliquotaInternaSchema(t *testing.T) {
	sql := readMigration(t, icmsMatrixMigration)
	normalized := normalizeMigrationSQL(sql)
	body := createTableBody(t, sql, "icms_aliquota_interna")

	assertExactColumns(t, body, []string{
		"uf", "aliquota", "fcp_embutido", "fonte", "lei", "vigente_desde", "vigente_ate",
	})
	for _, column := range []string{"uf", "aliquota", "fcp_embutido", "fonte", "lei", "vigente_desde"} {
		assertNotNullableColumn(t, body, column)
	}
	assertNullableColumn(t, body, "vigente_ate")

	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS icms_aliquota_interna",
		"PRIMARY KEY (uf, vigente_desde)",
		"CREATE UNIQUE INDEX IF NOT EXISTS icms_aliquota_interna_vigente ON icms_aliquota_interna (uf) WHERE vigente_ate IS NULL",
	} {
		if !strings.Contains(normalized, normalizeMigrationSQL(required)) {
			t.Errorf("icms_aliquota_interna schema missing %q", required)
		}
	}
}

// TestICMSMatrixMirrorSchema proves the synced-cell mirror shape: tenant_id
// is part of the primary key (multi-tenant non-negotiable), codtrib is
// nullable (an unread tributary classification is a valid state, not an
// error), and the one-open-row guarantee is the named partial unique index,
// never table-constraint UNIQUE (Postgres has no partial UNIQUE table
// constraint) and never left to writer discipline.
func TestICMSMatrixMirrorSchema(t *testing.T) {
	sql := readMigration(t, icmsMatrixMigration)
	normalized := normalizeMigrationSQL(sql)
	body := createTableBody(t, sql, "icms_matrix_mirror")

	assertExactColumns(t, body, []string{
		"tenant_id", "uf_origem", "uf_destino", "grupo_icms", "codtrib", "aliquota",
		"red_base", "perc_fcp", "linhas_candidatas", "ambiguo", "vigente_desde",
		"vigente_ate", "synced_at",
	})
	for _, column := range []string{
		"tenant_id", "uf_origem", "uf_destino", "grupo_icms", "linhas_candidatas",
		"ambiguo", "vigente_desde", "synced_at",
	} {
		assertNotNullableColumn(t, body, column)
	}
	// codtrib is deliberately nullable per the plan: a cell can be read
	// before its tributary classification is known.
	for _, column := range []string{"codtrib", "aliquota", "red_base", "perc_fcp", "vigente_ate"} {
		assertNullableColumn(t, body, column)
	}

	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS icms_matrix_mirror",
		"PRIMARY KEY (tenant_id, uf_origem, uf_destino, grupo_icms, vigente_desde)",
		"CREATE UNIQUE INDEX IF NOT EXISTS icms_matrix_mirror_vigente ON icms_matrix_mirror (tenant_id, uf_origem, uf_destino, grupo_icms) WHERE vigente_ate IS NULL",
	} {
		if !strings.Contains(normalized, normalizeMigrationSQL(required)) {
			t.Errorf("icms_matrix_mirror schema missing %q", required)
		}
	}

	if strings.Contains(strings.ToUpper(sql), "FOREIGN KEY") {
		t.Error("icms_matrix_mirror must not add a foreign key (uf_origem/uf_destino/grupo_icms reference ERP-side facts, not rows this schema owns)")
	}
}

// TestICMSAliquotaInternaSeedHasAll27UFsExactlyOnce proves the seed carries
// every Brazilian UF (26 states + DF) exactly once, and that the four UFs
// where the ERP's own TGFICM rate is stale/wrong (D-42/D-43) carry a
// specific, cited vigente_desde instead of the generic historical fallback.
func TestICMSAliquotaInternaSeedHasAll27UFsExactlyOnce(t *testing.T) {
	sql := readMigration(t, icmsMatrixMigration)
	upper := strings.ToUpper(sql)
	if !strings.Contains(upper, "INSERT INTO ICMS_ALIQUOTA_INTERNA") {
		t.Fatal("migration must seed icms_aliquota_interna")
	}

	allUFs := []string{
		"AC", "AL", "AM", "AP", "BA", "CE", "DF", "ES", "GO", "MA", "MT", "MS", "MG",
		"PA", "PB", "PR", "PE", "PI", "RJ", "RN", "RO", "RR", "RS", "SC", "SP", "SE", "TO",
	}
	if len(allUFs) != 27 {
		t.Fatalf("test fixture itself must list 27 UFs, has %d", len(allUFs))
	}
	for _, uf := range allUFs {
		literal := "('" + uf + "',"
		if got := strings.Count(sql, literal); got != 1 {
			t.Errorf("UF %s must appear exactly once as a seed row value, found %d", uf, got)
		}
	}

	// The four UFs whose ERP rate the mission refuted against the legislation
	// (D-42/D-43): each must carry its own cited law, not the generic
	// historical fonte/lei pair.
	citedLaw := map[string]string{
		"AL": "Lei 9.776/2025",
		"DF": "Lei 7.326/2023",
		"PR": "Lei 21.850/2023",
		"RJ": "Lei 10.253/2023 + LC 210/2023",
	}
	for uf, lei := range citedLaw {
		if !strings.Contains(sql, "'"+lei+"'") {
			t.Errorf("UF %s must cite %q as its lei, verbatim in the seed", uf, lei)
		}
	}
	// Column-aligned padding in the seed (extra spaces after short numeric
	// literals) is cosmetic; normalize before comparing so alignment
	// whitespace can never masquerade as a missing value.
	normalized := normalizeMigrationSQL(sql)
	if !strings.Contains(normalized, normalizeMigrationSQL("'AC', 19.0, 0, 'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'")) {
		t.Error("a representative historical UF (AC) must carry the generic fonte/lei pair and the 2000-01-01 stand-in date")
	}

	// vigente_desde is a column, not a literal seed date freely floating
	// around comments; the historical fallback date must appear at least 23
	// times (27 UFs minus the 4 specifically dated ones).
	if got := strings.Count(sql, "'2000-01-01'"); got != 23 {
		t.Errorf("expected exactly 23 UFs seeded with the historical 2000-01-01 vigente_desde, got %d", got)
	}
}

// TestICMSMatrixMigrationIsAdditiveOnly proves 0094 never touches an
// existing table's data or an unrelated table's shape — additive columns and
// brand-new tables only.
func TestICMSMatrixMigrationIsAdditiveOnly(t *testing.T) {
	sql := readMigration(t, icmsMatrixMigration)
	lower := strings.ToLower(sql)
	for _, forbidden := range []string{"alter column", "drop column", "drop table", "alter table listings", "alter table orders"} {
		if strings.Contains(lower, forbidden) {
			t.Errorf("icms_matrix migration must be additive only, found %q", forbidden)
		}
	}
}
