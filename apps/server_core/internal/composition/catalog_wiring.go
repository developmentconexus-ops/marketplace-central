package composition

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/adapters/erp/sankhyaoracle"
	"marketplace-central/apps/server_core/internal/contexts/catalog"
)

// CatalogWiring is the assembled catalog slice.
//
// Note what this file CANNOT name: nothing under any context's or adapter's
// internal/. That is not a convention — the import above was tried and the
// compiler refused it. This struct is what is left when the refusal is obeyed.
type CatalogWiring struct {
	Module *catalog.Module
	Feed   sankhyaoracle.Bundle
}

// WireCatalog assembles the slice from its two edges.
func WireCatalog(pool *pgxpool.Pool, oracleDB *sql.DB, sankhyaInstance string) (CatalogWiring, error) {
	bundle, err := sankhyaoracle.New(oracleDB, sankhyaInstance, time.Now)
	if err != nil {
		return CatalogWiring{}, err
	}
	return CatalogWiring{Module: catalog.New(pool), Feed: bundle}, nil
}
