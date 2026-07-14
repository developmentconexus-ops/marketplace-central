package oracle

import (
	"database/sql"

	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

// Database is the bounded Oracle runtime used by the adapter. It deliberately
// exposes only read operations, health probing, pool statistics, and shutdown.
type Database interface {
	queryer
	Close() error
	Stats() ports.PoolStats
}

func poolStatsFromDBStats(stats sql.DBStats) ports.PoolStats {
	return ports.PoolStats{
		Open:      stats.OpenConnections,
		InUse:     stats.InUse,
		WaitCount: int(stats.WaitCount),
	}
}
