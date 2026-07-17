package postgres

import "github.com/jackc/pgx/v5/pgxpool"

// Repository persists tenant-scoped market evidence.
type Repository struct {
	pool     *pgxpool.Pool
	tenantID string
}

func NewRepository(pool *pgxpool.Pool, tenantID string) *Repository {
	return &Repository{pool: pool, tenantID: tenantID}
}
