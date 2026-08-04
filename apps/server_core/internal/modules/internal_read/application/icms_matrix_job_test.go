package application_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/application"
	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type fakeResolver struct {
	cells []domain.ICMSMatrixCell
	err   error
	calls int
}

func (f *fakeResolver) ResolveCells(ctx context.Context) ([]domain.ICMSMatrixCell, error) {
	f.calls++
	return f.cells, f.err
}

type fakeApplier struct {
	written  int
	err      error
	gotCells []domain.ICMSMatrixCell
	gotTenant string
	calls    int
}

func (f *fakeApplier) ApplyCells(ctx context.Context, tenantID string, cells []domain.ICMSMatrixCell) (int, error) {
	f.calls++
	f.gotTenant = tenantID
	f.gotCells = cells
	return f.written, f.err
}

func cell(uf string, grupo int, ambiguo bool) domain.ICMSMatrixCell {
	return domain.ICMSMatrixCell{
		UFOrigem:         "MG",
		UFDestino:        uf,
		GrupoICMS:        grupo,
		LinhasCandidatas: 1,
		Ambiguo:          ambiguo,
	}
}

func fixedNow() time.Time {
	return time.Date(2026, 8, 4, 12, 0, 0, 0, time.UTC)
}

func TestICMSMatrixJobPersistsCountsInCursor(t *testing.T) {
	resolver := &fakeResolver{cells: []domain.ICMSMatrixCell{
		cell("BA", 122, false),
		cell("RJ", 311, true),
		cell("SP", 122, false),
	}}
	applier := &fakeApplier{written: 3}

	job := application.NewICMSMatrixJob("tenant-1", resolver, applier, fixedNow)
	next, err := job(context.Background(), nil)
	if err != nil {
		t.Fatalf("job: %v", err)
	}
	if applier.gotTenant != "tenant-1" {
		t.Errorf("tenant repassado = %q, quero \"tenant-1\"", applier.gotTenant)
	}
	if len(applier.gotCells) != 3 {
		t.Errorf("células repassadas = %d, quero 3", len(applier.gotCells))
	}

	var cursor application.ICMSMatrixCursor
	if err := json.Unmarshal(next, &cursor); err != nil {
		t.Fatalf("unmarshal cursor: %v", err)
	}
	if cursor.Cells != 3 {
		t.Errorf("cursor.Cells = %d, quero 3", cursor.Cells)
	}
	if cursor.Ambiguos != 1 {
		t.Errorf("cursor.Ambiguos = %d, quero 1", cursor.Ambiguos)
	}
	if cursor.Written != 3 {
		t.Errorf("cursor.Written = %d, quero 3", cursor.Written)
	}
	if !cursor.CompletedAt.Equal(fixedNow()) {
		t.Errorf("cursor.CompletedAt = %v, quero %v", cursor.CompletedAt, fixedNow())
	}
}

func TestICMSMatrixJobResolveErrorKeepsCursorAndFails(t *testing.T) {
	boom := errors.New("oracle fora do ar")
	resolver := &fakeResolver{err: boom}
	applier := &fakeApplier{}

	job := application.NewICMSMatrixJob("tenant-1", resolver, applier, fixedNow)
	prev := json.RawMessage(`{"cells":9}`)
	next, err := job(context.Background(), prev)
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, quero envolver %v", err, boom)
	}
	if string(next) != string(prev) {
		t.Errorf("cursor = %s, quero o anterior %s intacto", next, prev)
	}
	if applier.calls != 0 {
		t.Errorf("applier chamado %d vez(es); leitura falhou, nada pode ser aplicado", applier.calls)
	}
}

func TestICMSMatrixJobApplyErrorSurfaces(t *testing.T) {
	boom := errors.New("mirror: refusing to apply empty icms matrix")
	resolver := &fakeResolver{cells: nil}
	applier := &fakeApplier{err: boom}

	job := application.NewICMSMatrixJob("tenant-1", resolver, applier, fixedNow)
	prev := json.RawMessage(`{"cells":9}`)
	next, err := job(context.Background(), prev)
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, quero envolver %v", err, boom)
	}
	if string(next) != string(prev) {
		t.Errorf("cursor = %s, quero o anterior %s intacto", next, prev)
	}
}
