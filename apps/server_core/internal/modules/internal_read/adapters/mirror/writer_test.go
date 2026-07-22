package mirror

import (
	"context"
	"errors"
	"testing"
)

// The guards short-circuit before any pool access, so a nil pool is safe here and
// proves they fire fail-closed without a database.

func TestApplySnapshotRejectsEmptySnapshot(t *testing.T) {
	w := NewPgWriter(nil)
	if _, err := w.ApplySnapshot(context.Background(), "tenant_x", nil, nil); !errors.Is(err, ErrEmptySnapshot) {
		t.Fatalf("empty rows err = %v, want ErrEmptySnapshot", err)
	}
}

func TestApplySnapshotRejectsEmptyTenant(t *testing.T) {
	w := NewPgWriter(nil)
	_, err := w.ApplySnapshot(context.Background(), "", []Row{{CodigoProduto: "A"}}, nil)
	if err == nil {
		t.Fatal("empty tenantID: got nil error, want failure")
	}
	if errors.Is(err, ErrEmptySnapshot) {
		t.Fatalf("empty tenantID misclassified as ErrEmptySnapshot: %v", err)
	}
}
