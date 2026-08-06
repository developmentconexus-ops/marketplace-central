// Package catalog is the context's façade. The composition root builds a Module
// and hands it to consumers as a port; it never reaches past this file, because
// everything past this file is under internal/.
package catalog

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/application"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres"
	"marketplace-central/apps/server_core/internal/contexts/catalog/port"
)

// Module is Catalog, assembled.
type Module struct {
	service *application.Service
	reader  port.Reader
}

// New assembles the context from the ONLY thing an outsider may legitimately
// name: a connection pool.
//
// Every other collaborator — the store, the id factory, the reader — is chosen
// here, inside the context, because their types live under internal/ and a
// parameter typed by one of them would have no legal caller. That was the state
// this constructor was in: the composition root could only satisfy it by
// importing catalog/internal/postgres, which the compiler refused. The refusal
// was correct; the signature was the defect.
func New(pool *pgxpool.Pool) *Module {
	repo := postgres.NewRepository(pool)
	return &Module{
		service: application.NewService(repo, postgres.NewULIDFactory()),
		reader:  postgres.NewSummaryReader(repo),
	}
}

// IngestProduct folds one source observation into the catalogue.
func (m *Module) IngestProduct(ctx context.Context, o contracts.ProductObservation) (contracts.IngestResult, error) {
	return m.service.Ingest(ctx, o)
}

// Reader is Catalog's answer to identity questions from other contexts.
func (m *Module) Reader() port.Reader { return m.reader }
