package composition

import (
	"database/sql"
	"time"

	"marketplace-central/apps/server_core/internal/adapters/erp/sankhyaoracle"
	"marketplace-central/apps/server_core/internal/contexts/catalog"
	catalogpostgres "marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CatalogWiring is the assembled catalog slice.
//
// Note what this struct CANNOT say: there is no field typed by anything under
// an adapter's internal/, because no such type can be named here. That is the
// property, and it is enforced by the compiler rather than by this comment.
type CatalogWiring struct {
	Module *catalog.Module
	Feed   sankhyaoracle.Bundle
}

// WireCatalog assembles the slice from its two edges.
func WireCatalog(pool *pgxpool.Pool, oracleDB *sql.DB, sankhyaInstance string) (CatalogWiring, error) {
	repo := catalogpostgres.NewRepository(pool)
	module := catalog.New(repo, catalogpostgres.NewULIDFactory(), catalogpostgres.NewSummaryReader(repo))

	bundle, err := sankhyaoracle.New(oracleDB, sankhyaInstance, time.Now)
	if err != nil {
		return CatalogWiring{}, err
	}
	return CatalogWiring{Module: module, Feed: bundle}, nil
}
