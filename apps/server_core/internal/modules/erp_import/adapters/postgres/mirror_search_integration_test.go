//go:build integration

package postgres_test

import (
	"context"
	"reflect"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// The operator types what is on the keyboard ("valvula", "mictorio") while the
// ERP stores accented text, and they search by whatever identifier is at hand —
// reference code, product code or EAN, not only the description.
func TestMirrorCatalogSearchIsAccentInsensitiveAcrossIdentifiers(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror WHERE tenant_id=$1`, tenant)
	})

	_, err := pool.Exec(ctx, `INSERT INTO products_mirror
		(tenant_id,source,codigo_produto,descricao,referencia,ean,absent_in_last_snapshot)
		VALUES ($1,$2,'90033','VÁLVULA DE MICTÓRIO COMPACT PRESSMATIC','MLB6896003832','7899999000011',false),
		       ($1,$2,'24864','CHAVE P/ TUBO 36" GEDORE RED','R27160030','4060833012096',false)`, tenant, domain.SourceXLSX)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		query string
		want  []string
	}{
		{"valvula", []string{"90033"}},
		{"VALVULA", []string{"90033"}},
		{"válvula", []string{"90033"}},
		{"mictorio", []string{"90033"}},
		{"R27160030", []string{"24864"}},
		{"r27160030", []string{"24864"}},
		{"4060833012096", []string{"24864"}},
		{"24864", []string{"24864"}},
		{"", []string{"24864", "90033"}},
	} {
		page, err := repo.MirrorCatalogPage(ctx, tenant, domain.SourceXLSX, tc.query, 0, 10)
		if err != nil {
			t.Fatalf("query %q: %v", tc.query, err)
		}
		if got := mirrorCodes(page); !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("query %q matched %v, want %v", tc.query, got, tc.want)
		}
	}
}
