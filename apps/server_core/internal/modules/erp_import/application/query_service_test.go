package application

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
)

func TestQueryServiceListImportsForwardsTenant(t *testing.T) {
	want := []domain.ImportReport{{ID: "import-1", Protocol: "#001-E"}}
	repo := &fakeImportRepository{list: want}
	got, err := NewQueryService(repo, "tenant-query").ListImports(context.Background())
	if err != nil {
		t.Fatalf("ListImports() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) || repo.listTenant != "tenant-query" {
		t.Fatalf("result/tenant = %#v/%q, want %#v/tenant-query", got, repo.listTenant, want)
	}
}

func TestQueryServiceGetImportForwardsIDAndNotFound(t *testing.T) {
	repo := &fakeImportRepository{getErr: ports.ErrImportNotFound}
	_, err := NewQueryService(repo, "tenant-query").GetImport(context.Background(), "missing-id")
	if !errors.Is(err, ports.ErrImportNotFound) {
		t.Fatalf("GetImport() error = %v, want ErrImportNotFound", err)
	}
	if repo.getTenant != "tenant-query" || repo.getID != "missing-id" {
		t.Fatalf("tenant/id = %q/%q", repo.getTenant, repo.getID)
	}
}
