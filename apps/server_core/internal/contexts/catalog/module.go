// Package catalog is the context's façade. The composition root builds a Module
// and hands it to consumers as a port; it never reaches past this file, because
// everything past this file is under internal/.
package catalog

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/application"
	"marketplace-central/apps/server_core/internal/contexts/catalog/port"
)

// Module is Catalog, assembled.
type Module struct {
	service *application.Service
	reader  port.Reader
}

// New assembles the context from its dependencies.
func New(store application.Store, ids application.IDFactory, reader port.Reader) *Module {
	return &Module{service: application.NewService(store, ids), reader: reader}
}

// IngestProduct folds one source observation into the catalogue.
func (m *Module) IngestProduct(ctx context.Context, o contracts.ProductObservation) (contracts.IngestResult, error) {
	return m.service.Ingest(ctx, o)
}

// Reader is Catalog's answer to identity questions from other contexts.
func (m *Module) Reader() port.Reader { return m.reader }
