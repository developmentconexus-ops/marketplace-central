//go:build !cgo

package oracle

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

func OpenDB(context.Context, Config) (Database, error) {
	return nil, domain.NewReadError(domain.ReadErrorSourceUnavailable, "oracle driver requires cgo and a configured C toolchain", nil)
}
