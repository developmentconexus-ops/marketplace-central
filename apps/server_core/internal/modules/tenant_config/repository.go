package tenant_config

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrInvalidActiveSource is returned by Repository.Set when cfg.Source is not
// one of the three ratified ActiveSource values.
var ErrInvalidActiveSource = errors.New("invalid_active_source")

// Repository is the Postgres-backed ActiveSourceLookup + writer for a
// tenant's active-source configuration.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository builds a Repository over pool.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

var _ ActiveSourceLookup = (*Repository)(nil)

// Get returns the tenant's active-source Config, or ErrUnknownActiveSource if
// the tenant has no row (fail-closed — never a silent default).
func (r *Repository) Get(ctx context.Context, tenantID string) (Config, error) {
	var (
		cfg   Config
		setBy *string
	)
	cfg.TenantID = tenantID
	err := r.pool.QueryRow(ctx, `
		SELECT active_source, source_kind, set_at, set_by
		FROM active_source
		WHERE tenant_id = $1
	`, tenantID).Scan(&cfg.Source, &cfg.Kind, &cfg.SetAt, &setBy)
	if errors.Is(err, pgx.ErrNoRows) {
		return Config{}, ErrUnknownActiveSource
	}
	if err != nil {
		return Config{}, fmt.Errorf("get active source: %w", err)
	}
	if setBy != nil {
		cfg.SetBy = *setBy
	}
	return cfg, nil
}

// Set validates cfg.Source and upserts the tenant's active-source row. If
// cfg.Kind is empty it is derived from cfg.Source via DefaultKind. Tenant-
// scoped by construction: tenant_id is the row key (AC-01).
func (r *Repository) Set(ctx context.Context, cfg Config) error {
	if !cfg.Source.Valid() {
		return ErrInvalidActiveSource
	}
	kind := cfg.Kind
	if kind == "" {
		kind = cfg.Source.DefaultKind()
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO active_source (tenant_id, active_source, source_kind, set_at, set_by)
		VALUES ($1, $2, $3, now(), $4)
		ON CONFLICT (tenant_id) DO UPDATE SET
			active_source = EXCLUDED.active_source,
			source_kind = EXCLUDED.source_kind,
			set_at = now(),
			set_by = EXCLUDED.set_by
	`, cfg.TenantID, cfg.Source, kind, nullableString(cfg.SetBy))
	if err != nil {
		return fmt.Errorf("set active source: %w", err)
	}
	return nil
}

func nullableString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
