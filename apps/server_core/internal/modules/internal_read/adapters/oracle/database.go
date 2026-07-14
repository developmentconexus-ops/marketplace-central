package oracle

import "database/sql"

// Database is the bounded Oracle runtime used by the adapter. It deliberately
// exposes only read operations, health probing, pool statistics, and shutdown.
type Database interface {
	queryer
	Close() error
	Stats() sql.DBStats
}
