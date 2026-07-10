//go:build cgo

package oracle

import (
	"context"
	"database/sql"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"

	_ "github.com/godror/godror"
)

func OpenDB(ctx context.Context, cfg Config) (*sql.DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	db, err := sql.Open("godror", buildDSN(cfg))
	if err != nil {
		return nil, domain.NewReadError(domain.ReadErrorSourceUnavailable, "oracle connection bootstrap failed", nil)
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, domain.NewReadError(domain.ReadErrorSourceUnavailable, "oracle connection ping failed", nil)
	}
	return db, nil
}
