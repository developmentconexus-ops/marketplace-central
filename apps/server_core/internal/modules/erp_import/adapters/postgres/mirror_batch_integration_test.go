//go:build integration

package postgres_test

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// Stock Seguro asks for the stock of a whole page of listings at once. The batch
// read must answer from the ACTIVE source only, must ignore rows the last
// snapshot no longer carries, and must simply omit a code it has no row for —
// never invent one.
func TestMirrorProductsByCodesReadsOneSourceInOneRoundTrip(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror WHERE tenant_id=$1`, tenant)
	})

	_, err := pool.Exec(ctx, `INSERT INTO products_mirror
		(tenant_id,source,codigo_produto,descricao,estoque_total,absent_in_last_snapshot)
		VALUES ($1,'xlsx','90033','VÁLVULA DE MICTÓRIO','12',false),
		       ($1,'xlsx','90034','TORNEIRA AUTOMÁTICA',NULL,false),
		       ($1,'xlsx','90035','PRODUTO FORA DO ÚLTIMO SNAPSHOT','5',true),
		       ($1,'catalogo_cliente','90033','OUTRA FONTE','999',false)`, tenant)
	if err != nil {
		t.Fatal(err)
	}

	rows, err := repo.MirrorProductsByCodes(ctx, tenant, domain.SourceXLSX, []string{"90033", "90034", "90035", "90099"})
	if err != nil {
		t.Fatal(err)
	}

	codes := mirrorCodes(rows)
	sort.Strings(codes)
	if want := []string{"90033", "90034"}; !reflect.DeepEqual(codes, want) {
		t.Fatalf("got %v, want %v", codes, want)
	}
	for _, row := range rows {
		switch row.CodigoProduto {
		case "90033":
			if row.EstoqueTotal == nil || *row.EstoqueTotal != "12" {
				t.Fatalf("90033 must carry the active source's quantity, got %+v", row.EstoqueTotal)
			}
		case "90034":
			// Honest unknown: the source reported no stock for this product.
			if row.EstoqueTotal != nil {
				t.Fatalf("90034 has no reported stock, got %q", *row.EstoqueTotal)
			}
		}
	}

	empty, err := repo.MirrorProductsByCodes(ctx, tenant, domain.SourceXLSX, nil)
	if err != nil || len(empty) != 0 {
		t.Fatalf("an empty ask makes no query: rows=%v err=%v", empty, err)
	}
}
