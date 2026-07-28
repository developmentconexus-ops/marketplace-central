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

type fakeChainRepository struct {
	*fakeImportRepository
	chain       domain.ImportChain
	chainTenant string
	chainID     domain.ImportID
}

func (f *fakeChainRepository) GetImportChain(_ context.Context, tenantID string, importID domain.ImportID) (domain.ImportChain, error) {
	f.chainTenant = tenantID
	f.chainID = importID
	return f.chain, nil
}

func TestQueryServiceGetImportChainForwardsTenantAndID(t *testing.T) {
	repo := &fakeChainRepository{
		fakeImportRepository: &fakeImportRepository{},
		chain:                domain.ImportChain{Protocol: "#001-E", Importados: 4, Vinculados: 2, Enfileirados: 3},
	}
	got, err := NewQueryService(repo, "tenant-query").GetImportChain(context.Background(), "import-1")
	if err != nil {
		t.Fatalf("GetImportChain() error = %v", err)
	}
	if !reflect.DeepEqual(got, repo.chain) {
		t.Fatalf("GetImportChain() = %#v, want %#v", got, repo.chain)
	}
	if repo.chainTenant != "tenant-query" || repo.chainID != "import-1" {
		t.Fatalf("tenant/id = %q/%q", repo.chainTenant, repo.chainID)
	}
}

func TestQueryServiceGetImportChainFailsHonestlyWithoutTheCapability(t *testing.T) {
	// A repository that reads imports but not chains — the shape a decorator
	// leaves behind when it wraps the concrete repo and drops the optional port.
	got, err := NewQueryService(&fakeImportRepository{}, "tenant-query").GetImportChain(context.Background(), "import-1")
	if !errors.Is(err, ErrImportChainUnavailable) {
		t.Fatalf("GetImportChain() error = %v, want ErrImportChainUnavailable", err)
	}
	if !reflect.DeepEqual(got, domain.ImportChain{}) {
		t.Fatalf("GetImportChain() = %#v, want zero chain", got)
	}
}
