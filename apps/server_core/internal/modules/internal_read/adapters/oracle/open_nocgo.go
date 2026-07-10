//go:build !cgo

package oracle

import (
	"context"
	"database/sql"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

func OpenDB(context.Context, Config) (*sql.DB, error) {
	return nil, domain.NewReadError(domain.ReadErrorSourceUnavailable, "oracle driver requires cgo and a configured C toolchain", nil)
}
