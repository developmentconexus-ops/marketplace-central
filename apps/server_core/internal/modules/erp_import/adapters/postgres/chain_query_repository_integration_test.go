//go:build integration

package postgres_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/modules/erp_import/application"
	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
	"marketplace-central/apps/server_core/internal/modules/erp_import/transport"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestGetImportChainCountsCurrentQueueAcrossInstallations(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	importID := domain.ImportID("61000000-0000-0000-0000-000000000001")

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_links WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_products WHERE tenant_id=$1`, tenant)
	})

	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_protocols
			(id, tenant_id, file_sha256, protocol, source, imported_at, status, accepted_count, rejected_count, warning_count)
		VALUES ($1, $2, 'chain-counts-hash', '#610-E', 'xlsx', now(), 'COMPLETED', 99, 0, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_products (tenant_id, protocol_id, codprod, descrprod, custo, stock_physical)
		VALUES
			($2, $1, '101', 'Product 101', 1, 0),
			($2, $1, '102', 'Product 102', 1, 0),
			($2, $1, '103', 'Product 103', 1, 0),
			($2, $1, '104', 'Product 104', 1, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	// codprod 101 is resolved TWICE — one listing per installation. Vinculados must
	// still count it once, so dropping the DISTINCT from the vinculados CTE turns
	// this fixture's 2 into 3 and fails the assertion below.
	if _, err := pool.Exec(ctx, `
		INSERT INTO product_links
			(tenant_id, installation_id, provider_code, provider_item_id, state, internal_product_id)
		VALUES
			($1, 'installation-a', 'mercadolivre', 'item-101', 'resolved', 101),
			($1, 'installation-b', 'mercadolivre', 'item-101-b', 'resolved', 101),
			($1, 'installation-a', 'mercadolivre', 'item-102', 'resolved', 102),
			($1, 'installation-a', 'mercadolivre', 'item-103', 'pending', 103);
	`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO sync_state (tenant_id, installation_id, entity, cursor)
		VALUES
			($1, 'installation-a', 'market_queue', '{"pending":["101","102","OUTSIDE"]}'::jsonb),
			($1, 'installation-b', 'market_queue', '{"pending":["101","103"]}'::jsonb)
	`, tenant); err != nil {
		t.Fatal(err)
	}

	// queue_read_at is statement_timestamp() read inside PostgreSQL, so the
	// bracket has to come from that same clock. Bracketing with time.Now()
	// compared the container's clock against the host's; the two drift under
	// Docker Desktop, which is what made this assertion flake.
	before := databaseNow(ctx, t, pool)
	got, err := repo.GetImportChain(ctx, tenant, importID)
	after := databaseNow(ctx, t, pool)
	if err != nil {
		t.Fatal(err)
	}
	if got.Protocol != "#610-E" || got.Importados != 4 || got.Vinculados != 2 || got.Enfileirados != 3 {
		t.Fatalf("chain=%#v", got)
	}
	if got.QueueReadAt.Before(before) || got.QueueReadAt.After(after) {
		t.Fatalf("queue_read_at=%s outside [%s,%s]", got.QueueReadAt, before, after)
	}
}

func TestGetImportChainHTTPIntegration(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	importID := domain.ImportID("63000000-0000-0000-0000-000000000001")
	emptyImportID := domain.ImportID("63000000-0000-0000-0000-000000000002")

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_links WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_products WHERE tenant_id=$1`, tenant)
	})

	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_protocols
			(id, tenant_id, file_sha256, protocol, source, imported_at, status, accepted_count, rejected_count, warning_count)
		VALUES
			($1, $2, 'chain-http-hash', '#630-E', 'xlsx', now(), 'COMPLETED', 4, 0, 0),
			($3, $2, 'chain-http-empty-hash', '#631-E', 'xlsx', now(), 'COMPLETED', 1, 0, 0);
	`, importID, tenant, emptyImportID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_products (tenant_id, protocol_id, codprod, descrprod, custo, stock_physical)
		VALUES
			($2, $1, '301', 'Product 301', 1, 0),
			($2, $1, '302', 'Product 302', 1, 0),
			($2, $1, '303', 'Product 303', 1, 0),
			($2, $1, '304', 'Product 304', 1, 0),
			($2, $3, '399', 'Product 399', 1, 0);
	`, importID, tenant, emptyImportID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO product_links
			(tenant_id, installation_id, provider_code, provider_item_id, state, internal_product_id)
		VALUES
			($1, 'installation-a', 'mercadolivre', 'item-301', 'resolved', 301),
			($1, 'installation-a', 'mercadolivre', 'item-302', 'resolved', 302);
	`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO sync_state (tenant_id, installation_id, entity, cursor)
		VALUES
			($1, 'installation-a', 'market_queue', '{"pending":["301","302"]}'::jsonb),
			($1, 'installation-b', 'market_queue', '{"pending":["301","303"]}'::jsonb);
	`, tenant); err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	queries := application.NewQueryService(repo, tenant)
	transport.NewHandler(nil, queries).Register(mux)

	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/erp/imports/"+string(importID)+"/chain", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	body := response.Body.String()
	for _, want := range []string{`"protocol":"#630-E"`, `"importados":4`, `"vinculados":2`, `"enfileirados":3`} {
		if !strings.Contains(body, want) {
			t.Fatalf("body=%s missing %s", body, want)
		}
	}
	if !strings.Contains(body, `"queue_read_at":"`) {
		t.Fatalf("body=%s missing queue_read_at", body)
	}
	var decoded struct {
		QueueReadAt string `json:"queue_read_at"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if _, err := time.Parse(time.RFC3339, decoded.QueueReadAt); err != nil {
		t.Fatalf("queue_read_at=%q is not RFC3339: %v", decoded.QueueReadAt, err)
	}

	missing := httptest.NewRecorder()
	mux.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/erp/imports/63000000-0000-0000-0000-000000000099/chain", nil))
	if missing.Code != http.StatusNotFound {
		t.Fatalf("missing status=%d body=%s", missing.Code, missing.Body.String())
	}

	empty := httptest.NewRecorder()
	mux.ServeHTTP(empty, httptest.NewRequest(http.MethodGet, "/erp/imports/"+string(emptyImportID)+"/chain", nil))
	if empty.Code != http.StatusOK || !strings.Contains(empty.Body.String(), `"enfileirados":0`) {
		t.Fatalf("empty status=%d body=%s, want numeric enfileirados zero", empty.Code, empty.Body.String())
	}
}

func TestGetImportChainMissingAndProtocolWithoutSyncState(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	importID := domain.ImportID("62000000-0000-0000-0000-000000000001")
	otherImportID := domain.ImportID("62000000-0000-0000-0000-000000000002")
	otherTenant := tenant + "-other"

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, otherTenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_links WHERE tenant_id=$1`, otherTenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_products WHERE tenant_id IN ($1, $2)`, tenant, otherTenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_protocols WHERE tenant_id=$1`, otherTenant)
	})

	if _, err := repo.GetImportChain(ctx, tenant, "62000000-0000-0000-0000-000000000099"); !errors.Is(err, ports.ErrImportNotFound) {
		t.Fatalf("missing err=%v", err)
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_protocols
			(id, tenant_id, file_sha256, protocol, source, imported_at, status, accepted_count, rejected_count, warning_count)
		VALUES
			($1, $2, 'chain-empty-hash', '#620-E', 'xlsx', now(), 'COMPLETED', 1, 0, 0),
			($4, $3, 'chain-other-hash', '#620-E', 'xlsx', now(), 'COMPLETED', 1, 0, 0);
	`, importID, tenant, otherTenant, otherImportID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_products (tenant_id, protocol_id, codprod, descrprod, custo, stock_physical)
		VALUES
			($2, $1, '201', 'Main product', 1, 0),
			($3, $4, '201', 'Other product', 1, 0);
	`, importID, tenant, otherTenant, otherImportID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO product_links
			(tenant_id, installation_id, provider_code, provider_item_id, state, internal_product_id)
		VALUES ($1, 'installation-other', 'mercadolivre', 'item-other-201', 'resolved', 201);
	`, otherTenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO sync_state (tenant_id, installation_id, entity, cursor)
		VALUES ($1, 'installation-other', 'market_queue', '{"pending":["201"]}'::jsonb)
	`, otherTenant); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetImportChain(ctx, tenant, importID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Protocol != "#620-E" || got.Importados != 1 || got.Vinculados != 0 || got.Enfileirados != 0 {
		t.Fatalf("chain=%#v", got)
	}
}

// TestGetImportChainCountsLeadingZeroCodprod pins the resolved_products join
// fix: erp_import_products.codprod keeps the raw spreadsheet string
// ('00101'), while product_links.internal_product_id is a ParseInt'd integer
// column (101). Compared as text, '00101' = '101' is false and this fixture's
// linked product silently dropped out of vinculados. The comparison is
// numeric, so the padding is irrelevant without any key being folded.
func TestGetImportChainCountsLeadingZeroCodprod(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	importID := domain.ImportID("64000000-0000-0000-0000-000000000001")

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_links WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_products WHERE tenant_id=$1`, tenant)
	})

	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_protocols
			(id, tenant_id, file_sha256, protocol, source, imported_at, status, accepted_count, rejected_count, warning_count)
		VALUES ($1, $2, 'chain-leading-zero-hash', '#640-E', 'xlsx', now(), 'COMPLETED', 1, 0, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_products (tenant_id, protocol_id, codprod, descrprod, custo, stock_physical)
		VALUES ($2, $1, '00101', 'Product 00101', 1, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO product_links
			(tenant_id, installation_id, provider_code, provider_item_id, state, internal_product_id)
		VALUES ($1, 'installation-a', 'mercadolivre', 'item-00101', 'resolved', 101);
	`, tenant); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetImportChain(ctx, tenant, importID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Protocol != "#640-E" || got.Importados != 1 || got.Vinculados != 1 || got.Enfileirados != 0 {
		t.Fatalf("chain=%#v", got)
	}
}

// TestGetImportChainCountsCollidingCodprodsOnce pins the resolved_products
// OVERCOUNT that canonicalizing BOTH sides created. An ERP file can carry
// '101' and '00101' as two rows — erp_import_products.codprod is TEXT and its
// primary key treats them as distinct — while product_links carries ONE link,
// to internal product 101. Under ltrim-both-sides the two rows folded onto one
// join key and both matched: vinculados = 2 for a single linked product. The
// truth is 1, and it is 1 only because the counted key is the link's
// internal_product_id; exact numeric comparison alone still matches both rows.
func TestGetImportChainCountsCollidingCodprodsOnce(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	importID := domain.ImportID("67000000-0000-0000-0000-000000000001")

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_links WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_products WHERE tenant_id=$1`, tenant)
	})

	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_protocols
			(id, tenant_id, file_sha256, protocol, source, imported_at, status, accepted_count, rejected_count, warning_count)
		VALUES ($1, $2, 'chain-colliding-codprod-hash', '#670-E', 'xlsx', now(), 'COMPLETED', 2, 0, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_products (tenant_id, protocol_id, codprod, descrprod, custo, stock_physical)
		VALUES
			($2, $1, '101', 'Product 101', 1, 0),
			($2, $1, '00101', 'Product 00101', 1, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO product_links
			(tenant_id, installation_id, provider_code, provider_item_id, state, internal_product_id)
		VALUES ($1, 'installation-a', 'mercadolivre', 'item-101', 'resolved', 101);
	`, tenant); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetImportChain(ctx, tenant, importID)
	if err != nil {
		t.Fatal(err)
	}
	// importados counts rows and vinculados counts linked products, so 2 and 1
	// are both exact here — the numbers are allowed to differ, but only when
	// the file really did name one product twice.
	if got.Protocol != "#670-E" || got.Importados != 2 || got.Vinculados != 1 {
		t.Fatalf("chain=%#v", got)
	}
}

// TestGetImportChainQueueMatchesCodprodExactly pins the queued_products join
// as EXACT string equality, and the overcount that canonicalizing it produced.
//
// Both sides of that join are the same raw string out of the same slice:
// import_service builds `codes` from the accepted rows' Codprod and passes it
// to MarketEnqueuer.EnqueueMarketProducts → AppendPendingCodigos, which writes
// it verbatim into cursor->'pending'; those same accepted rows are what
// CopyFrom writes as erp_import_products.codprod. No other writer exists —
// AppendPendingCodigo (singular) has zero callers, and ProductsCursor has no
// pending field. Raw equality is therefore already the exact comparison, and
// ltrim only bought collisions: 'ABC' and '00ABC' are two import rows, and the
// single pending "ABC" counted both.
//
// The fixture holds all three cases at once: a padded codprod queued under its
// own padded value ('00101'), the ltrim collision pair, and a pending entry
// that belongs to no row of this import ("OUTSIDE").
func TestGetImportChainQueueMatchesCodprodExactly(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	importID := domain.ImportID("66000000-0000-0000-0000-000000000001")

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_links WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_products WHERE tenant_id=$1`, tenant)
	})

	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_protocols
			(id, tenant_id, file_sha256, protocol, source, imported_at, status, accepted_count, rejected_count, warning_count)
		VALUES ($1, $2, 'chain-queue-exact-hash', '#660-E', 'xlsx', now(), 'COMPLETED', 3, 0, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_products (tenant_id, protocol_id, codprod, descrprod, custo, stock_physical)
		VALUES
			($2, $1, '00101', 'Product 00101', 1, 0),
			($2, $1, 'ABC', 'Product ABC', 1, 0),
			($2, $1, '00ABC', 'Product 00ABC', 1, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	// '00101' is linked as internal product 101, which vinculados reaches by
	// numeric comparison. 'ABC' and '00ABC' are non-numeric and can never be
	// linked, which is correct: internal_product_id is an integer column.
	if _, err := pool.Exec(ctx, `
		INSERT INTO product_links
			(tenant_id, installation_id, provider_code, provider_item_id, state, internal_product_id)
		VALUES ($1, 'installation-a', 'mercadolivre', 'item-00101', 'resolved', 101);
	`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO sync_state (tenant_id, installation_id, entity, cursor)
		VALUES ($1, 'installation-a', 'market_queue', '{"pending":["00101","ABC","OUTSIDE"]}'::jsonb);
	`, tenant); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetImportChain(ctx, tenant, importID)
	if err != nil {
		t.Fatal(err)
	}
	// enfileirados = 2: '00101' and 'ABC'. '00ABC' is a different row that was
	// never queued (ltrim made it a third), "OUTSIDE" belongs to no row here.
	if got.Protocol != "#660-E" || got.Importados != 3 || got.Vinculados != 1 || got.Enfileirados != 2 {
		t.Fatalf("chain=%#v", got)
	}
}

// TestGetImportChainNonArrayPendingCursorDoesNotFailQuery pins the
// queued_products CASE guard: sync_state.cursor->'pending' is normally a
// JSON array, but a malformed or partial cursor write can leave it as an
// object or a scalar. Before the jsonb_typeof guard,
// jsonb_array_elements_text raised at query time (not a valid JSON array)
// and took the whole chain endpoint down for the tenant instead of reading
// as an empty queue.
func TestGetImportChainNonArrayPendingCursorDoesNotFailQuery(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_products WHERE tenant_id=$1`, tenant)
	})

	cases := []struct {
		name           string
		cursor         string
		importID       domain.ImportID
		protocol       string
		installationID string
	}{
		{name: "pending as object", cursor: `{"pending":{"foo":"bar"}}`, importID: domain.ImportID("65000000-0000-0000-0000-000000000001"), protocol: "#651-E", installationID: "installation-non-array-object"},
		{name: "pending as scalar", cursor: `{"pending":"501"}`, importID: domain.ImportID("65000000-0000-0000-0000-000000000002"), protocol: "#652-E", installationID: "installation-non-array-scalar"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			importID := tc.importID
			if _, err := pool.Exec(ctx, `
				INSERT INTO erp_import_protocols
					(id, tenant_id, file_sha256, protocol, source, imported_at, status, accepted_count, rejected_count, warning_count)
				VALUES ($1, $2, $3, $4, 'xlsx', now(), 'COMPLETED', 1, 0, 0);
			`, importID, tenant, "chain-non-array-hash-"+tc.name, tc.protocol); err != nil {
				t.Fatal(err)
			}
			if _, err := pool.Exec(ctx, `
				INSERT INTO erp_import_products (tenant_id, protocol_id, codprod, descrprod, custo, stock_physical)
				VALUES ($2, $1, '501', 'Product 501', 1, 0);
			`, importID, tenant); err != nil {
				t.Fatal(err)
			}
			if _, err := pool.Exec(ctx, `
				INSERT INTO sync_state (tenant_id, installation_id, entity, cursor)
				VALUES ($1, $3, 'market_queue', $2::jsonb);
			`, tenant, tc.cursor, tc.installationID); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_protocols WHERE id=$1 AND tenant_id=$2`, importID, tenant)
			})

			got, err := repo.GetImportChain(ctx, tenant, importID)
			if err != nil {
				t.Fatalf("expected a valid response, got query error: %v", err)
			}
			if got.Protocol != domain.Protocol(tc.protocol) || got.Importados != 1 || got.Enfileirados != 0 {
				t.Fatalf("chain=%#v", got)
			}
		})
	}
}

// databaseNow reads the same clock the queried value comes from, which is what
// lets the bracket stay exact instead of needing a tolerance window. A window
// wide enough to absorb host-to-container drift is also wide enough to accept a
// frozen or hardcoded timestamp.
func databaseNow(ctx context.Context, t *testing.T, pool *pgxpool.Pool) time.Time {
	t.Helper()
	var now time.Time
	if err := pool.QueryRow(ctx, `SELECT statement_timestamp()`).Scan(&now); err != nil {
		t.Fatalf("read database clock: %v", err)
	}
	return now
}
